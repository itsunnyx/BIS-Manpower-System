import React from 'react';
import { NavLink } from 'react-router-dom';

const Sidebar = ({ isCollapsed, isMobileOpen, onCloseMobile }) => {
  const menuItems = [
    {
      to: '/',
      label: 'หน้าแรก',
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5 sm:w-6 sm:h-6 flex-shrink-0">
          <path d="M11.47 3.841a.75.75 0 0 1 1.06 0l8.69 8.69a.75.75 0 1 0 1.06-1.061l-8.689-8.69a2.25 2.25 0 0 0-3.182 0l-8.69 8.69a.75.75 0 1 0 1.061 1.06l8.69-8.689Z" />
          <path d="m12 5.432 8.159 8.159c.03.03.06.058.091.086v6.198c0 1.035-.84 1.875-1.875 1.875H15a.75.75 0 0 1-.75-.75v-4.5a.75.75 0 0 0-.75-.75h-3a.75.75 0 0 0-.75.75V21a.75.75 0 0 1-.75.75H5.625a1.875 1.875 0 0 1-1.875-1.875v-6.198a2.29 2.29 0 0 0 .091-.086L12 5.432Z" />
        </svg>
      )
    },
    {
      to: '/requestform',
      label: 'ใบขออัตรากำลังคน',
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5 sm:w-6 sm:h-6 flex-shrink-0">
          <path fillRule="evenodd" d="M5.625 1.5c-1.036 0-1.875.84-1.875 1.875v17.25c0 1.035.84 1.875 1.875 1.875h12.75c1.035 0 1.875-.84 1.875-1.875V12.75A3.75 3.75 0 0 0 16.5 9h-1.875a1.875 1.875 0 0 1-1.875-1.875V5.25A3.75 3.75 0 0 0 9 1.5H5.625ZM7.5 15a.75.75 0 0 1 .75-.75h7.5a.75.75 0 0 1 0 1.5h-7.5A.75.75 0 0 1 7.5 15Zm.75 2.25a.75.75 0 0 0 0 1.5H12a.75.75 0 0 0 0-1.5H8.25Z" clipRule="evenodd" />
          <path d="M12.971 1.816A5.23 5.23 0 0 1 14.25 5.25v1.875c0 .207.168.375.375.375H16.5a5.23 5.23 0 0 1 3.434 1.279 9.768 9.768 0 0 0-6.963-6.963Z" />
        </svg>
      )
    },
    {
      to: '/approve',
      label: 'อนุมัติกำลังคน',
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5 sm:w-6 sm:h-6 flex-shrink-0">
          <path fillRule="evenodd" d="M9 1.5H5.625c-1.036 0-1.875.84-1.875 1.875v17.25c0 1.035.84 1.875 1.875 1.875h12.75c1.035 0 1.875-.84 1.875-1.875V12.75A3.75 3.75 0 0 0 16.5 9h-1.875a1.875 1.875 0 0 1-1.875-1.875V5.25A3.75 3.75 0 0 0 9 1.5Zm6.61 10.936a.75.75 0 1 0-1.22-.872l-3.236 4.53L9.53 14.47a.75.75 0 0 0-1.06 1.06l2.25 2.25a.75.75 0 0 0 1.14-.094l3.75-5.25Z" clipRule="evenodd" />
          <path d="M12.971 1.816A5.23 5.23 0 0 1 14.25 5.25v1.875c0 .207.168.375.375.375H16.5a5.23 5.23 0 0 1 3.434 1.279 9.768 9.768 0 0 0-6.963-6.963Z" />
        </svg>
      )
    },
    {
      to: '/dashboard',
      label: 'Dashboard',
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5 sm:w-6 sm:h-6 flex-shrink-0">
          <path fillRule="evenodd" d="M2.25 13.5a8.25 8.25 0 0 1 8.25-8.25.75.75 0 0 1 .75.75v6.75H18a.75.75 0 0 1 .75.75 8.25 8.25 0 0 1-16.5 0Z" clipRule="evenodd" />
          <path fillRule="evenodd" d="M12.75 3a.75.75 0 0 1 .75-.75 8.25 8.25 0 0 1 8.25 8.25.75.75 0 0 1-.75.75h-7.5a.75.75 0 0 1-.75-.75V3Z" clipRule="evenodd" />
        </svg>
      )
    },
    {
      to: '/usermanagement',
      label: 'จัดการผู้ใช้',
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="size-6">
          <path fillRule="evenodd" d="M8.25 6.75a3.75 3.75 0 1 1 7.5 0 3.75 3.75 0 0 1-7.5 0ZM15.75 9.75a3 3 0 1 1 6 0 3 3 0 0 1-6 0ZM2.25 9.75a3 3 0 1 1 6 0 3 3 0 0 1-6 0ZM6.31 15.117A6.745 6.745 0 0 1 12 12a6.745 6.745 0 0 1 6.709 7.498.75.75 0 0 1-.372.568A12.696 12.696 0 0 1 12 21.75c-2.305 0-4.47-.612-6.337-1.684a.75.75 0 0 1-.372-.568 6.787 6.787 0 0 1 1.019-4.38Z" clipRule="evenodd" />
          <path d="M5.082 14.254a8.287 8.287 0 0 0-1.308 5.135 9.687 9.687 0 0 1-1.764-.44l-.115-.04a.563.563 0 0 1-.373-.487l-.01-.121a3.75 3.75 0 0 1 3.57-4.047ZM20.226 19.389a8.287 8.287 0 0 0-1.308-5.135 3.75 3.75 0 0 1 3.57 4.047l-.01.121a.563.563 0 0 1-.373.486l-.115.04c-.567.2-1.156.349-1.764.441Z" />
        </svg>

      )
    }
  ];

  return (
    <>
      {/* Desktop Sidebar - แสดงตั้งแต่ md: ขึ้นไป */}
      <aside
        className={`
          hidden md:flex md:flex-col
          bg-neutral-600 text-white 
          transition-all duration-300 ease-in-out
          flex-shrink-0 h-full
          ${isCollapsed
            ? 'w-16 md:w-16 lg:w-16'
            : 'w-58 md:w-58 lg:w-58 xl:w-58'
          }
          shadow-lg
        `}
      >


        {/* Menu Items */}
        <nav className={`
          flex-1 overflow-y-auto overflow-x-hidden
          ${isCollapsed ? 'p-2' : 'p-3 md:p-4 lg:p-4'}
        `}>
          <div className="space-y-1 md:space-y-2">
            {menuItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) => `
                  flex items-center gap-2 md:gap-3
                  py-2.5 md:py-3
                  ${isCollapsed ? 'px-2 justify-center' : 'px-3 md:px-4'}
                  rounded-lg
                  text-slate-50 
                  hover:bg-neutral-500 hover:text-white 
                  active:bg-neutral-600
                  transition-all duration-200
                  text-sm md:text-base
                  font-medium
                  ${isActive
                    ? 'bg-neutral-500 shadow-lg ring-2 ring-neutral-400'
                    : 'hover:shadow-md'
                  }
                `}
                title={isCollapsed ? item.label : ''}
              >
                <span className="flex-shrink-0">
                  {item.icon}
                </span>
                {!isCollapsed && (
                  <span className="whitespace-nowrap overflow-hidden text-ellipsis">
                    {item.label}
                  </span>
                )}
              </NavLink>
            ))}
          </div>
        </nav>

      </aside>

      {/* Mobile Sidebar - แสดงเฉพาะ mobile/tablet */}
      <aside
        className={`
          md:hidden
          fixed top-0 left-0 bottom-0
          w-64 sm:w-72
          bg-neutral-600 text-white 
          shadow-2xl
          transform transition-transform duration-300 ease-in-out
          z-30
          ${isMobileOpen ? 'translate-x-0' : '-translate-x-full'}
          flex flex-col
        `}
      >
        {/* Header with Close button */}
        <div className="flex items-center justify-between p-4 sm:p-5 border-b border-neutral-500 flex-shrink-0">
          <h2 className="text-lg sm:text-xl font-bold">เมนู</h2>
          <button
            onClick={onCloseMobile}
            className="p-2 rounded-lg hover:bg-neutral-500 active:bg-neutral-600 transition-colors touch-manipulation"
            aria-label="ปิดเมนู"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="currentColor"
              className="w-6 h-6 sm:w-7 sm:h-7"
            >
              <path
                fillRule="evenodd"
                d="M5.47 5.47a.75.75 0 0 1 1.06 0L12 10.94l5.47-5.47a.75.75 0 1 1 1.06 1.06L13.06 12l5.47 5.47a.75.75 0 1 1-1.06 1.06L12 13.06l-5.47 5.47a.75.75 0 0 1-1.06-1.06L10.94 12 5.47 6.53a.75.75 0 0 1 0-1.06Z"
                clipRule="evenodd"
              />
            </svg>
          </button>
        </div>

        {/* Menu Items - Scrollable */}
        <nav className="flex-1 overflow-y-auto overflow-x-hidden p-4 sm:p-5">
          <div className="space-y-2">
            {menuItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                onClick={onCloseMobile}
                className={({ isActive }) => `
                  flex items-center gap-3
                  py-3 px-4
                  rounded-lg
                  text-slate-50 
                  hover:bg-neutral-500 hover:text-white 
                  active:bg-neutral-600
                  transition-all duration-200
                  text-sm sm:text-base
                  font-medium
                  touch-manipulation
                  ${isActive
                    ? 'bg-neutral-500 shadow-lg ring-2 ring-neutral-400'
                    : 'hover:shadow-md'
                  }
                `}
              >
                {item.icon}
                <span className="whitespace-nowrap">{item.label}</span>
              </NavLink>
            ))}
          </div>
        </nav>
      </aside>
    </>
  );
};

export default Sidebar;