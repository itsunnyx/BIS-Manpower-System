import React, { useEffect, useState } from "react";

const Usermainpage = () => {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    // เรียก API จาก backend โดยตรง (ไม่มี auth)
    const fetchRequests = async () => {
      try {
        const res = await fetch("/api/v1/requester/requests");
        if (!res.ok) {
          throw new Error(`HTTP error! status: ${res.status}`);
        }
        const data = await res.json();
        setRequests(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchRequests();
  }, []);

  if (loading) return <p>กำลังโหลดข้อมูล...</p>;
  if (error) return <p>เกิดข้อผิดพลาด: {error}</p>;

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-3">รายการคำขอ Manpower</h2>
      {requests.length === 0 ? (
        <p>ไม่มีข้อมูลคำขอ</p>
      ) : (
        <table className="border-collapse border border-gray-400 w-full text-sm">
          <thead>
            <tr className="bg-gray-100">
              <th className="border border-gray-400 px-2 py-1">Doc No</th>
              <th className="border border-gray-400 px-2 py-1">วันที่</th>
              <th className="border border-gray-400 px-2 py-1">แผนก</th>
              <th className="border border-gray-400 px-2 py-1">สถานะต้นสังกัด</th>
              <th className="border border-gray-400 px-2 py-1">สถานะ HR</th>
              <th className="border border-gray-400 px-2 py-1">สถานะบริหาร</th>
            </tr>
          </thead>
          <tbody>
            {requests.map((req, i) => (
              <tr key={i}>
                <td className="border border-gray-400 px-2 py-1">{req.doc_no}</td>
                <td className="border border-gray-400 px-2 py-1">
                  {new Date(req.doc_date).toLocaleDateString("th-TH")}
                </td>
                <td className="border border-gray-400 px-2 py-1">{req.department_name}</td>
                <td className="border border-gray-400 px-2 py-1">{req.origin_status}</td>
                <td className="border border-gray-400 px-2 py-1">{req.hr_status}</td>
                <td className="border border-gray-400 px-2 py-1">{req.manager_status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default Usermainpage;
