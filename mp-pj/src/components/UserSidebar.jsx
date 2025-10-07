import React from 'react';
import { NavLink } from 'react-router-dom';

const UserSidebar = ({ isCollapsed }) => {
  const menuItems = [
    {
      to: '/user',
      label: 'หน้าแรก',
      end: true,
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6 flex-shrink-0">
          <path d="M11.47 3.841a.75.75 0 0 1 1.06 0l8.69 8.69a.75.75 0 1 0 1.06-1.061l-8.689-8.69a2.25 2.25 0 0 0-3.182 0l-8.69 8.69a.75.75 0 1 0 1.061 1.06l8.69-8.689Z" />
          <path d="m12 5.432 8.159 8.159c.03.03.06.058.091.086v6.198c0 1.035-.84 1.875-1.875 1.875H15a.75.75 0 0 1-.75-.75v-4.5a.75.75 0 0 0-.75-.75h-3a.75.75 0 0 0-.75.75V21a.75.75 0 0 1-.75.75H5.625a1.875 1.875 0 0 1-1.875-1.875v-6.198a2.29 2.29 0 0 0 .091-.086L12 5.432Z" />
        </svg>
      ),
    },
    {
      to: '/user/requestform',
      label: 'ใบขออัตรากำลังคน',
      end: false,
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6 flex-shrink-0">
          <path fillRule="evenodd" d="M5.625 1.5c-1.036 0-1.875.84-1.875 1.875v17.25c0 1.035.84 1.875 1.875 1.875h12.75c1.035 0 1.875-.84 1.875-1.875V12.75A3.75 3.75 0 0 0 16.5 9h-1.875a1.875 1.875 0 0 1-1.875-1.875V5.25A3.75 3.75 0 0 0 9 1.5H5.625ZM7.5 15a.75.75 0 0 1 .75-.75h7.5a.75.75 0 0 1 0 1.5h-7.5A.75.75 0 0 1 7.5 15Zm.75 2.25a.75.75 0 0 0 0 1.5H12a.75.75 0 0 0 0-1.5H8.25Z" clipRule="evenodd" />
          <path d="M12.971 1.816A5.23 5.23 0 0 1 14.25 5.25v1.875c0 .207.168.375.375.375H16.5a5.23 5.23 0 0 1 3.434 1.279 9.768 9.768 0 0 0-6.963-6.963Z" />
        </svg>
      ),
    },
  ];

  return (
    <aside
      className={`
        flex flex-col
        bg-neutral-600 text-white 
        transition-all duration-300 ease-in-out
        flex-shrink-0 h-full
        ${isCollapsed ? 'w-16' : 'w-58'}
        shadow-lg
      `}
    >
      {/* Menu Items */}
      <nav className={`flex-1 overflow-y-auto overflow-x-hidden ${isCollapsed ? 'p-2' : 'p-4'}`}>
        <div className="space-y-2">
          {menuItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive }) => `
                flex items-center gap-3
                py-3
                ${isCollapsed ? 'px-2 justify-center' : 'px-4'}
                rounded-lg
                text-slate-50 
                hover:bg-neutral-500 hover:text-white 
                transition-all duration-200
                text-base
                font-medium
                ${isActive ? 'bg-neutral-500 shadow-lg ring-2 ring-neutral-400' : ''}
              `}
              title={isCollapsed ? item.label : ''}
            >
              <span className="flex-shrink-0">{item.icon}</span>
              {!isCollapsed && (
                <span className="whitespace-nowrap overflow-hidden text-ellipsis">{item.label}</span>
              )}
            </NavLink>
          ))}
        </div>
      </nav>
    </aside>
  );
};

export default UserSidebar;