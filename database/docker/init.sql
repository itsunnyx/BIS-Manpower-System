-- =========================================
-- ENUMs (มาตรฐานค่าในระบบ)
-- =========================================
DO $$ BEGIN
  CREATE TYPE request_stage AS ENUM (
    'origin',         -- ผู้ยื่นต้นสังกัด/หัวหน้างาน
    'manager',        -- ผู้จัดการฝ่าย (ตามฝ่ายที่ร้องขอ)
    'hr',             -- ผู้จัดการฝ่าย HR
    'management'      -- ผู้บริหาร (ถ้าต้องมีชั้นนี้)
  );
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE decision_status AS ENUM (
    'pending',        -- รอดำเนินการ
    'submitted',      -- ส่งคำร้องแล้ว (ต้นสังกัด)
    'forwarded',      -- ส่งต่อชั้นถัดไป
    'approved',       -- อนุมัติ
    'rejected',       -- ไม่อนุมัติ
    'returned',       -- ส่งกลับแก้ไข
    'cancelled'       -- ยกเลิกคำร้อง
  );
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE sync_status_enum AS ENUM ('success','failed');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- =========================================
-- Departments (หน่วยงาน)
-- =========================================
CREATE TABLE IF NOT EXISTS departments (
  department_id   SERIAL PRIMARY KEY,
  department_name VARCHAR(255) NOT NULL UNIQUE,
  created_at      TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);

-- =========================================
-- Users (ผู้ใช้งานระบบ)
-- role ครอบคลุม actor ใน DFD: manager, hr, approver, recruiter(เจ้าหน้าที่สรรหา), admin
-- =========================================
CREATE TABLE IF NOT EXISTS users (
  user_id        SERIAL PRIMARY KEY,
  name           VARCHAR(255) NOT NULL,
  email          VARCHAR(255) NOT NULL UNIQUE,
  password_hash  TEXT NOT NULL,
  role           VARCHAR(32) NOT NULL CHECK (role IN ('manager','hr','approver','recruiter','admin')),
  department_id  INT REFERENCES departments(department_id) ON UPDATE CASCADE,
  created_at     TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
  updated_at     TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);

-- =========================================
-- D1: Manpower Requests (คำร้องขอกำลังคน)
-- สอดคล้องคอลัมน์หน้ารายงาน: เลขที่เอกสาร/วันที่เอกสาร/แผนก/สถานะรายชั้น/วันครบกำหนด
-- =========================================
CREATE TABLE IF NOT EXISTS manpower_requests (
  request_id          SERIAL PRIMARY KEY,

  -- เอกสาร/อ้างอิง
  doc_no              VARCHAR(32) NOT NULL UNIQUE,  -- เช่น PQ24110012
  doc_date            DATE NOT NULL,
  due_date            DATE,                         -- วันที่ครบกำหนด (ในภาพรายงาน)

  -- ผู้ขอและแผนกที่เกี่ยวข้อง
  department_id       INT NOT NULL REFERENCES departments(department_id) ON UPDATE CASCADE,
  requested_by        INT NOT NULL REFERENCES users(user_id) ON UPDATE CASCADE,

  -- รายละเอียดความต้องการ
  position_title      VARCHAR(255) NOT NULL,
  position_level      VARCHAR(100),
  num_required        INT NOT NULL CHECK (num_required > 0),
  reason              TEXT,

  -- สถานะตามชั้นที่ DFD ระบุ (ต้นสังกัด/ผู้จัดการฝ่าย → HR → ผู้บริหาร)
  origin_status       decision_status NOT NULL DEFAULT 'submitted',
  manager_status      decision_status NOT NULL DEFAULT 'pending',
  hr_status           decision_status NOT NULL DEFAULT 'pending',
  management_status   decision_status NOT NULL DEFAULT 'pending',

  -- สถานะรวม (เพื่อ query ง่ายบน dashboard)
  overall_status decision_status GENERATED ALWAYS AS (
    CASE
      WHEN management_status IN ('approved','rejected') THEN management_status
      WHEN hr_status IN ('returned','rejected')         THEN hr_status
      WHEN manager_status IN ('returned','rejected')    THEN manager_status
      WHEN origin_status IN ('cancelled','returned')    THEN origin_status
      ELSE 'pending'
    END
  ) STORED,

  created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =========================================
-- D2: Approval Logs (ประวัติการอนุมัติ/การดำเนินการ)
-- เก็บทุก step ตาม DFD: ผลการอนุมัติในใบคำร้อง, ดำเนินการตามฝ่ายที่ร้องขอ/ฝ่ายบุคคล ฯลฯ
-- =========================================
CREATE TABLE IF NOT EXISTS approval_logs (
  log_id      BIGSERIAL PRIMARY KEY,
  request_id  INT NOT NULL REFERENCES manpower_requests(request_id) ON DELETE CASCADE,
  actor_id    INT NOT NULL REFERENCES users(user_id) ON UPDATE CASCADE,
  stage       request_stage NOT NULL,   -- origin/manager/hr/management
  action      decision_status NOT NULL, -- submitted/forwarded/approved/rejected/returned/cancelled
  comment     TEXT,
  created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =========================================
-- Notifications (รองรับ Process 4: แจ้งเตือนสถานะ)
-- =========================================
CREATE TABLE IF NOT EXISTS notifications (
  notif_id     BIGSERIAL PRIMARY KEY,
  request_id   INT NOT NULL REFERENCES manpower_requests(request_id) ON DELETE CASCADE,
  to_user_id   INT NOT NULL REFERENCES users(user_id) ON UPDATE CASCADE,
  title        VARCHAR(255) NOT NULL,
  message      TEXT NOT NULL,
  is_read      BOOLEAN NOT NULL DEFAULT FALSE,
  created_at   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =========================================
-- Recruitment Sync (เชื่อม “เจ้าหน้าที่สรรหา/ระบบสรรหา” หลังอนุมัติ)
-- =========================================
CREATE TABLE IF NOT EXISTS recruitment_sync (
  sync_id          BIGSERIAL PRIMARY KEY,
  request_id       INT NOT NULL REFERENCES manpower_requests(request_id) ON DELETE CASCADE,
  sync_status      sync_status_enum NOT NULL,
  response_message TEXT,
  sync_time        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =========================================
-- Trigger function: อัปเดต updated_at อัตโนมัติ
-- =========================================
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_departments_updated_at
BEFORE UPDATE ON departments
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER trg_requests_updated_at
BEFORE UPDATE ON manpower_requests
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

-- =========================================
-- Indexes สำหรับงานรายงาน/แดชบอร์ด
-- =========================================
CREATE INDEX IF NOT EXISTS idx_users_email           ON users (lower(email));
CREATE INDEX IF NOT EXISTS idx_req_doc_date          ON manpower_requests (doc_date);
CREATE INDEX IF NOT EXISTS idx_req_department        ON manpower_requests (department_id);
CREATE INDEX IF NOT EXISTS idx_req_statuses          ON manpower_requests (overall_status, management_status, hr_status, manager_status);
CREATE INDEX IF NOT EXISTS idx_logs_request_stage    ON approval_logs (request_id, stage, action, created_at);
CREATE INDEX IF NOT EXISTS idx_notif_to_user         ON notifications (to_user_id, is_read, created_at);

-- =========================================
-- Sample Data (ให้ทดลอง flow ตาม DFD)
-- =========================================
INSERT INTO departments (department_name) VALUES
  ('ฝ่ายการตลาด'),
  ('ฝ่ายทรัพยากรบุคคล'),
  ('ฝ่ายไอที')
ON CONFLICT DO NOTHING;

INSERT INTO users (name,email,password_hash,role,department_id) VALUES
  ('สมชาย (หัวหน้างาน)','somchai@company.com','hashed_pw1','manager',
    (SELECT department_id FROM departments WHERE department_name='ฝ่ายการตลาด')),
  ('ศุภกร (ผจก.ฝ่ายการตลาด)','supakorn_mgr@company.com','hashed_pw2','approver',
    (SELECT department_id FROM departments WHERE department_name='ฝ่ายการตลาด')),
  ('สุภาพร (ผจก.ฝ่าย HR)','supaporn_hr@company.com','hashed_pw3','hr',
    (SELECT department_id FROM departments WHERE department_name='ฝ่ายทรัพยากรบุคคล')),
  ('อนันต์ (ผู้บริหาร)','anan@company.com','hashed_pw4','approver',
    (SELECT department_id FROM departments WHERE department_name='ฝ่ายไอที')),
  ('นีรนุช (เจ้าหน้าที่สรรหา)','neeranuch_rec@company.com','hashed_pw5','recruiter',
    (SELECT department_id FROM departments WHERE department_name='ฝ่ายทรัพยากรบุคคล'))
ON CONFLICT (email) DO NOTHING;

-- ตัวอย่างคำร้อง 1 ใบ (เดินครบทุกสเตจ)
INSERT INTO manpower_requests
  (doc_no, doc_date, due_date, department_id, requested_by,
   position_title, position_level, num_required, reason,
   origin_status, manager_status, hr_status, management_status)
VALUES
  ('PQ24110012', DATE '2024-11-23', DATE '2024-12-29',
   (SELECT department_id FROM departments WHERE department_name='ฝ่ายการตลาด'),
   (SELECT user_id FROM users WHERE email='somchai@company.com'),
   'Sales Executive','Staff', 1, 'ขยายทีมลูกค้าใหม่',
   'submitted','approved','approved','approved')
ON CONFLICT (doc_no) DO NOTHING;

-- บันทึกประวัติการดำเนินการ (D2)
INSERT INTO approval_logs (request_id, actor_id, stage, action, comment)
SELECT r.request_id, u.user_id, 'origin', 'submitted', 'ยื่นคำร้อง'
FROM manpower_requests r
JOIN users u ON u.email='somchai@company.com'
WHERE r.doc_no='PQ24110012';

INSERT INTO approval_logs (request_id, actor_id, stage, action, comment)
SELECT r.request_id, u.user_id, 'manager', 'approved', 'ผู้จัดการฝ่ายอนุมัติ'
FROM manpower_requests r
JOIN users u ON u.email='supakorn_mgr@company.com'
WHERE r.doc_no='PQ24110012';

INSERT INTO approval_logs (request_id, actor_id, stage, action, comment)
SELECT r.request_id, u.user_id, 'hr', 'approved', 'HR ตรวจสอบครบ ส่งผู้บริหาร'
FROM manpower_requests r
JOIN users u ON u.email='supaporn_hr@company.com'
WHERE r.doc_no='PQ24110012';

INSERT INTO approval_logs (request_id, actor_id, stage, action, comment)
SELECT r.request_id, u.user_id, 'management', 'approved', 'ผู้บริหารอนุมัติ 1 อัตรา'
FROM manpower_requests r
JOIN users u ON u.email='anan@company.com'
WHERE r.doc_no='PQ24110012';

-- แจ้งเตือนสถานะ (Process 4)
INSERT INTO notifications (request_id, to_user_id, title, message)
SELECT r.request_id, u.user_id, 'คำร้องได้รับอนุมัติ', 'คำร้อง PQ24110012 ได้รับอนุมัติ'
FROM manpower_requests r
JOIN users u ON u.email='somchai@company.com'
WHERE r.doc_no='PQ24110012';

-- ตัวอย่างซิงก์ไประบบสรรหา (หลังอนุมัติ)
INSERT INTO recruitment_sync (request_id, sync_status, response_message)
SELECT request_id, 'success', 'ส่งข้อมูลไประบบสรรหาแล้ว'
FROM manpower_requests
WHERE doc_no='PQ24110012';
