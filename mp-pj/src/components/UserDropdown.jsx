import React, { useState, useEffect, useRef } from 'react';

// สร้าง Component ชื่อ UserDropdown
const UserDropdown = () => {
  // State สำหรับจัดการการแสดงผลของ dropdown
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);

  // Function สำหรับปิด dropdown เมื่อคลิกนอกพื้นที่ของ component
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    // เพิ่ม event listener เมื่อ component ถูก mount
    document.addEventListener('mousedown', handleClickOutside);
    
    // Cleanup function: ลบ event listener ออกเมื่อ component ถูก unmount
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    // ใช้ div ครอบทั้งหมดและใส่ ref เพื่อตรวจจับการคลิก
    <div className="relative" ref={dropdownRef}>
      {/* --- Button สำหรับเปิด/ปิด Dropdown --- */}
      <button
        id="dropdownAvatarNameButton"
        onClick={() => setIsOpen(!isOpen)} // Toggle state เมื่อคลิก
        className="flex items-center text-sm pe-1 font-medium text-gray-900 rounded-full hover:text-blue-600 dark:hover:text-blue-500 md:me-0 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:text-white"
        type="button"
      >
        <span className="sr-only">Open user menu</span>
        <img
          className="w-8 h-8 me-2 rounded-full"
          src="https://flowbite.com/docs/images/people/profile-picture-3.jpg" // ใช้ URL เต็มเพื่อให้แสดงผลได้
          alt="user photo"
        />
        Bonnie Green
        <svg
          className="w-2.5 h-2.5 ms-3"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 10 6"
        >
          <path
            stroke="currentColor"
            strokeLinecap="round" // แก้ไข attribute stroke-linecap เป็น strokeLinecap
            strokeLinejoin="round" // แก้ไข attribute stroke-linejoin เป็น strokeLinejoin
            strokeWidth="2" // แก้ไข attribute stroke-width เป็น strokeWidth
            d="m1 1 4 4 4-4"
          />
        </svg>
      </button>

      {/* --- Dropdown Menu --- */}
      {/* ใช้เงื่อนไข isOpen เพื่อแสดง/ซ่อนเมนู */}
      {isOpen && (
        <div
          id="dropdownAvatarName"
          className="z-10 absolute top-full right-0 mt-2 bg-white divide-y divide-gray-100 rounded-lg shadow-sm w-44 dark:bg-gray-700 dark:divide-gray-600"
        >
          <div className="px-4 py-3 text-sm text-gray-900 dark:text-white">
            <div className="font-medium">Pro User</div>
            <div className="truncate">name@flowbite.com</div>
          </div>
          <ul
            className="py-2 text-sm text-gray-700 dark:text-gray-200"
            aria-labelledby="dropdownAvatarNameButton"
          >
            <li>
              <a href="#" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">
                Dashboard
              </a>
            </li>
            <li>
              <a href="#" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">
                Settings
              </a>
            </li>
            <li>
              <a href="#" className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white">
                Earnings
              </a>
            </li>
          </ul>
          <div className="py-2">
            <a href="#" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 dark:text-gray-200 dark:hover:text-white">
              Sign out
            </a>
          </div>
        </div>
      )}
    </div>
  );
};

export default UserDropdown;