-- =====================================================================
-- Manpower Planning - init.sql (compact production-ready, idempotent)
-- =====================================================================

-- ================= ENUMS =================
DO $$ BEGIN
  CREATE TYPE request_status AS ENUM ('pending','approved','rejected');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE user_role AS ENUM ('requester','hr','approver','manager','admin');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- =============== LOOKUP TABLES =======================
CREATE TABLE IF NOT EXISTS divisions(
  division_id SERIAL PRIMARY KEY,
  name_th     VARCHAR(200) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS departments(
  department_id SERIAL PRIMARY KEY,
  division_id   INT REFERENCES divisions(division_id) ON DELETE RESTRICT,
  name_th       VARCHAR(200) NOT NULL,
  UNIQUE(division_id, name_th)
);

CREATE TABLE IF NOT EXISTS employment_types(
  employment_type_id SERIAL PRIMARY KEY,
  code VARCHAR(50) UNIQUE,
  name_th VARCHAR(200) NOT NULL
);

CREATE TABLE IF NOT EXISTS contract_types(
  contract_type_id SERIAL PRIMARY KEY,
  code VARCHAR(50) UNIQUE,
  name_th VARCHAR(200) NOT NULL
);

CREATE TABLE IF NOT EXISTS request_reasons(
  reason_id SERIAL PRIMARY KEY,
  code VARCHAR(50) UNIQUE,
  name_th VARCHAR(200) NOT NULL
);

CREATE TABLE IF NOT EXISTS genders(
  gender_id SERIAL PRIMARY KEY,
  code VARCHAR(20) UNIQUE,
  name_th VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS nationalities(
  nationality_id SERIAL PRIMARY KEY,
  name_th VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS education_levels(
  education_level_id SERIAL PRIMARY KEY,
  code VARCHAR(50) UNIQUE,
  name_th VARCHAR(200) NOT NULL
);

-- =============== USERS =========================
CREATE TABLE IF NOT EXISTS users (
  user_id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  full_name VARCHAR(255) NOT NULL,
  password_user VARCHAR(255) NOT NULL,  -- dev seed เท่านั้น; โปรดใช้แฮชในจริง
  role user_role NOT NULL,              -- requester/hr/approver/management/admin
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  last_login TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ในกรณีที่ users.role เคยเป็น VARCHAR มาก่อน ให้พยายามแปลงเป็น enum
DO $migrate$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name='users' AND column_name='role' AND udt_name <> 'user_role'
  ) THEN
    BEGIN
      ALTER TABLE users
        ALTER COLUMN role TYPE user_role USING role::user_role;
    EXCEPTION WHEN others THEN
      RAISE NOTICE 'Role migration to enum skipped (manual fix may be required).';
    END;
  END IF;
END
$migrate$;

-- =============== COMMON TRIGGERS/FUNCTIONS ==========================
-- updated_at helper
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_users_updated ON users;
CREATE TRIGGER trg_users_updated
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

-- =====================================================================
-- Doc number generator: PQYYMM####
-- =====================================================================
CREATE TABLE IF NOT EXISTS doc_no_counters (
  prefix   VARCHAR(10) NOT NULL,
  yymm     CHAR(4)     NOT NULL,
  last_seq INT         NOT NULL DEFAULT 0,
  PRIMARY KEY (prefix, yymm)
);

CREATE OR REPLACE FUNCTION next_doc_no(p_prefix TEXT DEFAULT 'PQ')
RETURNS TEXT
LANGUAGE plpgsql
AS $$
DECLARE
  v_yymm CHAR(4) := to_char(CURRENT_DATE, 'YYMM');
  v_seq  INT;
BEGIN
  INSERT INTO doc_no_counters(prefix, yymm, last_seq)
  VALUES (p_prefix, v_yymm, 0)
  ON CONFLICT (prefix, yymm) DO NOTHING;

  UPDATE doc_no_counters
  SET last_seq = last_seq + 1
  WHERE prefix = p_prefix AND yymm = v_yymm
  RETURNING last_seq INTO v_seq;

  RETURN p_prefix || v_yymm || lpad(v_seq::text, 4, '0');
END;
$$;

-- =============== JOB POSITIONS (Auto-fill) ===================
DROP TABLE IF EXISTS job_positions CASCADE;
CREATE TABLE job_positions (
  position_id     SERIAL PRIMARY KEY,
  position_code   VARCHAR(50)  NOT NULL UNIQUE,        -- เช่น HUN-P1, SAL-ADM
  position_title  VARCHAR(255) NOT NULL,
  division_id     INT REFERENCES divisions(division_id)     ON DELETE RESTRICT,
  department_id   INT REFERENCES departments(department_id) ON DELETE RESTRICT,
  position_level  VARCHAR(100),
  is_active       BOOLEAN NOT NULL DEFAULT TRUE,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_job_pos_dept ON job_positions(department_id, is_active);
CREATE INDEX IF NOT EXISTS idx_job_pos_div  ON job_positions(division_id, is_active);

CREATE OR REPLACE FUNCTION trg_job_positions_touch()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at := now(); RETURN NEW;
END $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_job_positions_touch ON job_positions;
CREATE TRIGGER trg_job_positions_touch
BEFORE UPDATE ON job_positions
FOR EACH ROW EXECUTE FUNCTION trg_job_positions_touch();

-- =============== MANPOWER REQUEST (หัวเอกสาร) =======
DROP TABLE IF EXISTS manpower_request_requirements CASCADE;
DROP TABLE IF EXISTS approval_logs CASCADE;
DROP TABLE IF EXISTS manpower_requests CASCADE;

CREATE TABLE manpower_requests(
  request_id       SERIAL PRIMARY KEY,
  doc_no           VARCHAR(50) UNIQUE NOT NULL,
  doc_date         DATE NOT NULL DEFAULT CURRENT_DATE,

  division_id      INT REFERENCES divisions(division_id),
  department_id    INT REFERENCES departments(department_id),

  requested_by     INT REFERENCES users(user_id) ON DELETE RESTRICT,
  contact_name     VARCHAR(255),

  position_code    VARCHAR(50),                      -- FK -> job_positions.position_code
  position_title   VARCHAR(255) NOT NULL,            -- snapshot จาก job_positions (เวอร์ชันวันที่สร้าง)
  position_level   VARCHAR(100),

  employment_type_id INT REFERENCES employment_types(employment_type_id),
  contract_type_id   INT REFERENCES contract_types(contract_type_id),
  reason_id          INT REFERENCES request_reasons(reason_id),
  reason_note        TEXT,

  num_required     INT NOT NULL DEFAULT 1,

  origin_status     request_status NOT NULL DEFAULT 'pending',
  hr_status         request_status NOT NULL DEFAULT 'pending',
  manager_status request_status NOT NULL DEFAULT 'pending',

  overall_status   request_status NOT NULL DEFAULT 'pending',
  due_date         DATE,
  print_count      INT NOT NULL DEFAULT 0,
  last_printed_at  TIMESTAMPTZ,

  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT chk_mpr_num_required CHECK (num_required > 0)
);

-- ผูก FK position_code -> job_positions
ALTER TABLE manpower_requests
  DROP CONSTRAINT IF EXISTS fk_mpr_position_code,
  ADD CONSTRAINT fk_mpr_position_code
  FOREIGN KEY (position_code) REFERENCES job_positions(position_code);

-- ถ้าไม่ส่ง doc_no → ออกเลขอัตโนมัติ
CREATE OR REPLACE FUNCTION set_doc_no_if_null()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.doc_no IS NULL OR NEW.doc_no = '' THEN
    NEW.doc_no := next_doc_no('PQ');
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_mpr_autonumber ON manpower_requests;
CREATE TRIGGER trg_mpr_autonumber
BEFORE INSERT ON manpower_requests
FOR EACH ROW EXECUTE FUNCTION set_doc_no_if_null();

-- updated_at
DROP TRIGGER IF EXISTS trg_mpr_updated ON manpower_requests;
CREATE TRIGGER trg_mpr_updated
BEFORE UPDATE ON manpower_requests
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

-- รวมสถานะอัตโนมัติ
CREATE OR REPLACE FUNCTION set_overall_status()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.origin_status='rejected' OR NEW.hr_status='rejected' OR NEW.manager_status='rejected' THEN
    NEW.overall_status := 'rejected';
  ELSIF NEW.origin_status='approved' AND NEW.hr_status='approved' AND NEW.manager_status='approved' THEN
    NEW.overall_status := 'approved';
  ELSE
    NEW.overall_status := 'pending';
  END IF;
  RETURN NEW;
END $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_mpr_overall ON manpower_requests;
CREATE TRIGGER trg_mpr_overall
BEFORE INSERT OR UPDATE ON manpower_requests
FOR EACH ROW EXECUTE FUNCTION set_overall_status();

-- ตรวจความสอดคล้อง division_id กับ department_id
CREATE OR REPLACE FUNCTION chk_dept_belongs_to_div()
RETURNS TRIGGER AS $$
DECLARE v_div INT;
BEGIN
  IF NEW.department_id IS NULL THEN
    RETURN NEW;
  END IF;
  SELECT division_id INTO v_div FROM departments WHERE department_id = NEW.department_id;
  IF NEW.division_id IS DISTINCT FROM v_div THEN
    RAISE EXCEPTION 'department % does not belong to division %', NEW.department_id, NEW.division_id;
  END IF;
  RETURN NEW;
END $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_mpr_div_dept ON manpower_requests;
CREATE TRIGGER trg_mpr_div_dept
BEFORE INSERT OR UPDATE ON manpower_requests
FOR EACH ROW EXECUTE FUNCTION chk_dept_belongs_to_div();

-- =============== REQUIREMENTS (คุณสมบัติ) =============
CREATE TABLE manpower_request_requirements(
  request_id INT PRIMARY KEY REFERENCES manpower_requests(request_id) ON DELETE CASCADE,
  age_min INT, age_max INT,
  gender_id INT REFERENCES genders(gender_id),
  nationality_id INT REFERENCES nationalities(nationality_id),
  experience_years_min INT,
  education_level_id INT REFERENCES education_levels(education_level_id),
  special_qualification TEXT,
  CONSTRAINT chk_age_nonneg CHECK (COALESCE(age_min,0) >= 0 AND COALESCE(age_max,0) >= 0),
  CONSTRAINT chk_age_order CHECK (
    age_min IS NULL OR age_max IS NULL OR age_min <= age_max
  ),
  CONSTRAINT chk_exp_nonneg CHECK (COALESCE(experience_years_min,0) >= 0)
);

-- =============== APPROVAL LOG ========================
CREATE TABLE approval_logs(
  log_id SERIAL PRIMARY KEY,
  request_id INT REFERENCES manpower_requests(request_id) ON DELETE CASCADE,
  action_by  INT REFERENCES users(user_id),
  actor_role user_role,               -- จำกัดให้เป็นหนึ่งใน enum เดิม
  stage      VARCHAR(50),             -- origin/hr/management (ตัว workflow)
  action     VARCHAR(50),             -- approve/reject/print/update
  from_status request_status,
  to_status   request_status,
  comment    TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =============== PROCEDURES / HELPERS =================
-- กดพิมพ์: เพิ่มตัวนับ + บันทึกเวลา + เขียน log
CREATE OR REPLACE FUNCTION bump_print_count(p_request_id INT, p_user_id INT DEFAULT NULL)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
  UPDATE manpower_requests
  SET print_count = print_count + 1,
      last_printed_at = now()
  WHERE request_id = p_request_id;

  INSERT INTO approval_logs(request_id, action_by, actor_role, stage, action, comment)
  VALUES (p_request_id, p_user_id, 'requester', NULL, 'print', 'printed form');
END;
$$;

-- อัปเดตสถานะรายขั้น + audit
-- p_stage: 'origin' | 'hr' | 'management'
-- p_status: 'approved' | 'rejected'
CREATE OR REPLACE FUNCTION set_stage_status(
  p_request_id INT,
  p_stage TEXT,
  p_status request_status,
  p_user_id INT,
  p_actor_role user_role,
  p_comment TEXT DEFAULT NULL
) RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
  v_from request_status;
BEGIN
  IF p_stage NOT IN ('origin','hr','manager') THEN
    RAISE EXCEPTION 'invalid stage: %', p_stage;
  END IF;
  IF p_status NOT IN ('approved','rejected') THEN
    RAISE EXCEPTION 'invalid status for stage update: %', p_status;
  END IF;

  IF p_stage = 'origin' THEN
    SELECT origin_status INTO v_from FROM manpower_requests WHERE request_id = p_request_id FOR UPDATE;
    UPDATE manpower_requests SET origin_status = p_status WHERE request_id = p_request_id;
  ELSIF p_stage = 'hr' THEN
    SELECT hr_status INTO v_from FROM manpower_requests WHERE request_id = p_request_id FOR UPDATE;
    UPDATE manpower_requests SET hr_status = p_status WHERE request_id = p_request_id;
  ELSE
    SELECT manager_status INTO v_from FROM manpower_requests WHERE request_id = p_request_id FOR UPDATE;
    UPDATE manpower_requests SET manager_status = p_status WHERE request_id = p_request_id;
  END IF;

  INSERT INTO approval_logs(request_id, action_by, actor_role, stage, action, from_status, to_status, comment)
  VALUES (p_request_id, p_user_id, p_actor_role, p_stage,
          CASE WHEN p_status='approved' THEN 'approve' ELSE 'reject' END,
          v_from, p_status, p_comment);
END;
$$;

-- =============== INDEXES =============================
CREATE INDEX IF NOT EXISTS idx_mpr_doc_date        ON manpower_requests(doc_date);
CREATE INDEX IF NOT EXISTS idx_mpr_dept            ON manpower_requests(department_id);
CREATE INDEX IF NOT EXISTS idx_mpr_statuses        ON manpower_requests(overall_status, origin_status, hr_status, manager_status);
CREATE INDEX IF NOT EXISTS idx_mpr_position_title  ON manpower_requests(position_title);
CREATE INDEX IF NOT EXISTS idx_mpr_position_code   ON manpower_requests(position_code);

-- =============== SEED LOOKUPS (ย่อ) ==================
INSERT INTO divisions(name_th) VALUES
  ('สายงานผลิต'),
  ('สายงานพาณิชย์'),
  ('สายงานสำนักงาน'),
  ('สายงานแม่พิมพ์')
ON CONFLICT (name_th) DO NOTHING;

-- แผนกตัวอย่าง
WITH dv AS ( SELECT division_id, name_th FROM divisions )
INSERT INTO departments(division_id, name_th)
SELECT d.division_id, x.dept_name
FROM (
  SELECT 'สายงานสำนักงาน' AS div_name, 'ฝ่ายทรัพยากรบุคคลและธุรการ' AS dept_name UNION ALL
  SELECT 'สายงานผลิต', 'ฝ่ายซ่อมบำรุง' UNION ALL
  SELECT 'สายงานพาณิชย์', 'ฝ่ายธุรการขาย' UNION ALL
  SELECT 'สายงานผลิต', 'ฝ่ายวางแผนและจัดการวัตถุดิบ'
) x
JOIN dv d ON d.name_th = x.div_name
ON CONFLICT (division_id, name_th) DO NOTHING;

INSERT INTO employment_types(code,name_th) VALUES
('MONTHLY','รายเดือน'),('DAILY','รายวัน'),('INTERN','นักศึกษาฝึกงาน')
ON CONFLICT DO NOTHING;

INSERT INTO contract_types(code,name_th) VALUES
('PERMANENT','สัญญาจ้างแบบไม่มีกำหนดระยะเวลา'),
('FIXED_TERM','สัญญาจ้างมีกำหนดระยะเวลา')
ON CONFLICT DO NOTHING;

INSERT INTO request_reasons(code,name_th) VALUES
('INCREASE_HEADCOUNT','เพิ่มอัตรากำลังคน'),
('REPLACE','ทดแทนพนักงานลาออก')
ON CONFLICT DO NOTHING;

INSERT INTO genders(code,name_th) VALUES
('ANY','ไม่จำกัดเพศ'),('MALE','ชาย'),('FEMALE','หญิง')
ON CONFLICT DO NOTHING;

INSERT INTO education_levels(code,name_th) VALUES
('HS','มัธยมปลาย/ปวช.'),('BA','ปริญญาตรี'),('MA','ปริญญาโท')
ON CONFLICT DO NOTHING;

INSERT INTO nationalities(name_th) VALUES ('ไทย') ON CONFLICT DO NOTHING;

-- ผู้ใช้ตัวอย่าง (dev)
INSERT INTO users(email, full_name, password_user, role)
VALUES
  ('requester@naklasombut.com','หัวหน้างาน','1234','requester'),
  ('hr@naklasombut.com','ฝ่ายบุคคล','1234','hr'),
  ('approver@naklasombut.com','ต้นสังกัด','1234','approver'),
  ('manager@naklasombut.com','ผู้บริหาร','1234','manager'),
  ('admin@naklasombut.com','แอดมิน','1234','admin')
ON CONFLICT (email) DO NOTHING;

-- Seed ตำแหน่งงานหลักสำหรับ Auto-fill
WITH x AS (
  SELECT d.department_id, v.division_id, d.name_th AS dept_name
  FROM departments d JOIN divisions v ON v.division_id = d.division_id
  WHERE d.name_th IN ('ฝ่ายทรัพยากรบุคคลและธุรการ','ฝ่ายซ่อมบำรุง','ฝ่ายธุรการขาย','ฝ่ายวางแผนและจัดการวัตถุดิบ')
)
INSERT INTO job_positions(position_code, position_title, division_id, department_id, position_level)
SELECT * FROM (
  SELECT 'HUN-P1','พนักงานทำความสะอาด', v.division_id, v.department_id, 'Staff'       FROM x v WHERE v.dept_name='ฝ่ายทรัพยากรบุคคลและธุรการ' UNION ALL
  SELECT 'HUN-ADM','เจ้าหน้าที่ธุรการ',   v.division_id, v.department_id, 'Staff'       FROM x v WHERE v.dept_name='ฝ่ายทรัพยากรบุคคลและธุรการ' UNION ALL
  SELECT 'MTN-TECH','ช่างซ่อมบำรุง',      v.division_id, v.department_id, 'Technician'  FROM x v WHERE v.dept_name='ฝ่ายซ่อมบำรุง' UNION ALL
  SELECT 'MTN-SUP','หัวหน้าช่างซ่อมบำรุง',v.division_id, v.department_id, 'Supervisor' FROM x v WHERE v.dept_name='ฝ่ายซ่อมบำรุง' UNION ALL
  SELECT 'SAL-ADM','เจ้าหน้าที่ธุรการขาย', v.division_id, v.department_id, 'Staff'       FROM x v WHERE v.dept_name='ฝ่ายธุรการขาย' UNION ALL
  SELECT 'SAL-CS','เจ้าหน้าที่บริการลูกค้า',v.division_id, v.department_id, 'Staff'      FROM x v WHERE v.dept_name='ฝ่ายธุรการขาย' UNION ALL
  SELECT 'SCM-PLN','เจ้าหน้าที่วางแผนวัสดุ',v.division_id, v.department_id, 'Staff'      FROM x v WHERE v.dept_name='ฝ่ายวางแผนและจัดการวัตถุดิบ'
) s
ON CONFLICT (position_code) DO NOTHING;


-- ================== SAMPLE MANPOWER REQUESTS (3 รายการ) ==================
-- หมายเหตุ:
-- - เลือกแผนก: ฝ่ายซ่อมบำรุง / ฝ่ายธุรการขาย / ฝ่ายวางแผนและจัดการวัตถุดิบ
-- - requested_by ใช้ user ตัวอย่าง 'requester@naklasombut.com'
-- - ใช้ position_code จากตาราง job_positions ที่ seed ไปแล้ว

BEGIN;

-- 1) ฝ่ายซ่อมบำรุง : ช่างซ่อมบำรุง 2 อัตรา
WITH dept AS (
  SELECT d.department_id, v.division_id
  FROM departments d JOIN divisions v ON v.division_id = d.division_id
  WHERE d.name_th = 'ฝ่ายซ่อมบำรุง' LIMIT 1
),
ins AS (
  INSERT INTO manpower_requests
  (doc_no, doc_date, division_id, department_id, requested_by, contact_name,
   position_code, position_title, position_level,
   employment_type_id, contract_type_id, reason_id, reason_note,
   num_required, origin_status, hr_status, manager_status, due_date)
  SELECT
    NULL, CURRENT_DATE,
    dept.division_id, dept.department_id,
    (SELECT user_id FROM users WHERE email='requester@naklasombut.com'),
    'หัวหน้างาน',
    'MTN-TECH', 'ช่างซ่อมบำรุง', 'Technician',
    (SELECT employment_type_id FROM employment_types WHERE code='MONTHLY'),
    (SELECT contract_type_id   FROM contract_types   WHERE code='PERMANENT'),
    (SELECT reason_id          FROM request_reasons  WHERE code='INCREASE_HEADCOUNT'),
    'รองรับแผน PM เพิ่ม',
    2, 'pending','pending','pending',
    CURRENT_DATE + INTERVAL '14 day'
  FROM dept
  RETURNING request_id
)
INSERT INTO manpower_request_requirements
(request_id, age_min, age_max, gender_id, nationality_id, experience_years_min, education_level_id, special_qualification)
SELECT
  i.request_id, 22, 35,
  (SELECT gender_id FROM genders WHERE code='ANY'),
  (SELECT nationality_id FROM nationalities WHERE name_th='ไทย'),
  1,
  (SELECT education_level_id FROM education_levels WHERE code='HS'),
  'พร้อมทำงานเป็นกะ/โอที'
FROM ins i;

-- 2) ฝ่ายธุรการขาย : เจ้าหน้าที่ธุรการขาย 1 อัตรา (ทดแทน)
WITH dept AS (
  SELECT d.department_id, v.division_id
  FROM departments d JOIN divisions v ON v.division_id = d.division_id
  WHERE d.name_th = 'ฝ่ายธุรการขาย' LIMIT 1
),
ins AS (
  INSERT INTO manpower_requests
  (doc_no, doc_date, division_id, department_id, requested_by, contact_name,
   position_code, position_title, position_level,
   employment_type_id, contract_type_id, reason_id, reason_note,
   num_required, origin_status, hr_status, manager_status, due_date)
  SELECT
    NULL, CURRENT_DATE,
    dept.division_id, dept.department_id,
    (SELECT user_id FROM users WHERE email='requester@naklasombut.com'),
    'หัวหน้างาน',
    'SAL-ADM', 'เจ้าหน้าที่ธุรการขาย', 'Staff',
    (SELECT employment_type_id FROM employment_types WHERE code='MONTHLY'),
    (SELECT contract_type_id   FROM contract_types   WHERE code='PERMANENT'),
    (SELECT reason_id          FROM request_reasons  WHERE code='REPLACE'),
    'ทดแทนพนักงานลาออก',
    1, 'pending','pending','pending',
    CURRENT_DATE + INTERVAL '10 day'
  FROM dept
  RETURNING request_id
)
INSERT INTO manpower_request_requirements
(request_id, age_min, age_max, gender_id, nationality_id, experience_years_min, education_level_id, special_qualification)
SELECT
  i.request_id, 21, 30,
  (SELECT gender_id FROM genders WHERE code='ANY'),
  (SELECT nationality_id FROM nationalities WHERE name_th='ไทย'),
  0,
  (SELECT education_level_id FROM education_levels WHERE code='BA'),
  'สื่อสารดี ใช้ Excel ได้'
FROM ins i;

-- 3) ฝ่ายวางแผนและจัดการวัตถุดิบ : เจ้าหน้าที่วางแผนวัสดุ 1 อัตรา (สัญญาจ้างมีกำหนด)
WITH dept AS (
  SELECT d.department_id, v.division_id
  FROM departments d JOIN divisions v ON v.division_id = d.division_id
  WHERE d.name_th = 'ฝ่ายวางแผนและจัดการวัตถุดิบ' LIMIT 1
),
ins AS (
  INSERT INTO manpower_requests
  (doc_no, doc_date, division_id, department_id, requested_by, contact_name,
   position_code, position_title, position_level,
   employment_type_id, contract_type_id, reason_id, reason_note,
   num_required, origin_status, hr_status, manager_status, due_date)
  SELECT
    NULL, CURRENT_DATE,
    dept.division_id, dept.department_id,
    (SELECT user_id FROM users WHERE email='requester@naklasombut.com'),
    'หัวหน้างาน',
    'SCM-PLN', 'เจ้าหน้าที่วางแผนวัสดุ', 'Staff',
    (SELECT employment_type_id FROM employment_types WHERE code='MONTHLY'),
    (SELECT contract_type_id   FROM contract_types   WHERE code='FIXED_TERM'),
    (SELECT reason_id          FROM request_reasons  WHERE code='INCREASE_HEADCOUNT'),
    'รองรับโครงการใหม่ 6 เดือน',
    1, 'pending','pending','pending',
    CURRENT_DATE + INTERVAL '30 day'
  FROM dept
  RETURNING request_id
)
INSERT INTO manpower_request_requirements
(request_id, age_min, age_max, gender_id, nationality_id, experience_years_min, education_level_id, special_qualification)
SELECT
  i.request_id, 23, 32,
  (SELECT gender_id FROM genders WHERE code='ANY'),
  (SELECT nationality_id FROM nationalities WHERE name_th='ไทย'),
  1,
  (SELECT education_level_id FROM education_levels WHERE code='BA'),
  'วิเคราะห์ข้อมูล/สต็อกได้ดี'
FROM ins i;

COMMIT;