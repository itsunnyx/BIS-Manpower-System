import React, { useState } from 'react'; 
import { Outlet } from 'react-router-dom';
import UserNavbar from '../components/UserNavbar';
import UserSidebar from '../components/UserSidebar';

const Userlayout = () => {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);

  const toggleSidebar = () => {
    setIsSidebarCollapsed(!isSidebarCollapsed);
  };

  return (
    <div className="flex flex-col h-screen w-full overflow-hidden">
      {/* Navbar - Fixed at top */}
      <UserNavbar onToggleSidebar={toggleSidebar} />

      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar */}
        <UserSidebar isCollapsed={isSidebarCollapsed} />

        {/* Content Area - Scrollable */}
        <main className="flex-1 overflow-y-auto overflow-x-hidden bg-gray-50 w-full">
          <div className="min-h-full">
            <div className="w-full max-w-[100vw] mx-auto px-10 py-8">
              <Outlet />
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default Userlayout;