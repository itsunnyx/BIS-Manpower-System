import React, { useState, useEffect, useRef } from 'react';

// --- Main Dropdown Component ---
export const Dropdown = ({ children, className }) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);

  const toggleDropdown = () => setIsOpen(!isOpen);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // children[0] คือ DropdownButton
  // children[1] คือ DropdownMenu

  return (
    <div ref={dropdownRef} className={`relative inline-block text-left ${className || ''}`}>
      {/* 1. สร้าง div ครอบ Button เพื่อรับ onClick แทนการใช้ React.cloneElement */}
      <div onClick={toggleDropdown}>
        {children[0]}
      </div>

      {/* 2. ใช้ Transition เพื่อให้มี Animation (Optional แต่แนะนำ) */}
      {isOpen && (
        <div
          className="absolute right-0 z-10 mt-2 w-full origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
          role="menu"
          aria-orientation="vertical"
          aria-labelledby="menu-button"
          tabIndex="-1"
        >
          <div className="py-1" role="none">
            {children[1]}
          </div>
        </div>
      )}
    </div>
  );
};

// --- Sub-components ---

// ไม่จำเป็นต้องรับ onClick อีกต่อไป เพราะ div ด้านบนจัดการแล้ว
export const DropdownButton = ({ children, className, ...props }) => (
  <button
    type="button"
    {...props}
    className={`inline-flex w-full justify-center items-center gap-x-1.5 rounded-md px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ${className || ''}`}
    id="menu-button"
    aria-expanded="true"
    aria-haspopup="true"
  >
    {children}
  </button>
);

// แก้ไข DropdownMenu ให้นิยามแค่ children
export const DropdownMenu = ({ children }) => {
    // ใช้ React.Fragment เพื่อไม่ให้มี div ครอบที่ไม่จำเป็น
    return <>{children}</>;
};

// แก้ไข DropdownItem ให้เป็น component ที่สมบูรณ์มากขึ้น
export const DropdownItem = ({ children, onClick, href }) => {
  // ถ้ามี href ให้ render เป็น <a> แต่ถ้ามี onClick ให้ render เป็น <button>
  const Tag = href ? 'a' : 'button';

  return (
    <Tag
      href={href}
      onClick={onClick}
      className="text-neutral-900 block w-full text-left px-4 py-2  hover:bg-gray-100 hover:text-neutral-900"
      role="menuitem"
      tabIndex="-1"
    >
      {children}
    </Tag>
  );
};

export const DropdownDivider = () => (
  <div className="border-t border-gray-200 my-1"></div>
);