import { Dropdown, DropdownButton, DropdownMenu, DropdownItem, DropdownDivider } from './Dropdown';
import { ChevronDownIcon } from '@heroicons/react/solid';
import { Link } from 'react-router-dom';

const UserNavbar = ({ onToggleSidebar }) => {
  const currentUser = {
    firstname: "nakล่าสมบัติ",
    lastname: "Doefdfggdgdf",
    role: "user",
  };

  const signOut = () => {
    console.log("User signed out");
  };

  return (
    <nav className="bg-neutral-600 shadow-lg sticky top-0 z-50">
      <div className="w-full px-8">
        <div className="flex justify-between items-center h-18">

          {/* ฝั่งซ้าย: ปุ่มเมนู + โลโก้ */}
          <div className="flex items-center gap-4 -ml-5">
            {/* ปุ่มเมนู */}
            <button
              onClick={onToggleSidebar}
              className="p-2 hover:bg-neutral-500 rounded-lg transition-colors active:scale-95"
              aria-label="Toggle sidebar"
            >
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                viewBox="0 0 24 24" 
                fill="white" 
                className="w-6 h-6"
              >
                <path 
                  fillRule="evenodd" 
                  d="M3 6.75A.75.75 0 0 1 3.75 6h16.5a.75.75 0 0 1 0 1.5H3.75A.75.75 0 0 1 3 6.75ZM3 12a.75.75 0 0 1 .75-.75h16.5a.75.75 0 0 1 0 1.5H3.75A.75.75 0 0 1 3 12Zm0 5.25a.75.75 0 0 1 .75-.75h16.5a.75.75 0 0 1 0 1.5H3.75a.75.75 0 0 1-.75-.75Z" 
                  clipRule="evenodd" 
                />
              </svg>
            </button>

            {/* Logo */}
            <Link 
              to="/user" 
              className="flex items-center group"
            >
              <img
                src="/images/nakla.svg"
                alt="Manpower"
                className="h-13 w-auto group-hover:scale-110 transition-transform duration-200"
              />
            </Link>
          </div>

          {/* ฝั่งขวา: Notification and Profile */}
          <div className="flex items-center gap-3">
            {/* Notification Button */}
            <button 
              className="p-2 hover:bg-neutral-500 rounded-lg transition-colors active:scale-95 relative"
              aria-label="Notifications"
            >
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                viewBox="0 0 24 24" 
                fill="white" 
                className="w-6 h-6"
              >
                <path 
                  fillRule="evenodd" 
                  d="M5.25 9a6.75 6.75 0 0 1 13.5 0v.75c0 2.123.8 4.057 2.118 5.52a.75.75 0 0 1-.297 1.206c-1.544.57-3.16.99-4.831 1.243a3.75 3.75 0 1 1-7.48 0 24.585 24.585 0 0 1-4.831-1.244.75.75 0 0 1-.298-1.205A8.217 8.217 0 0 0 5.25 9.75V9Zm4.502 8.9a2.25 2.25 0 1 0 4.496 0 25.057 25.057 0 0 1-4.496 0Z" 
                  clipRule="evenodd" 
                />
              </svg>
            </button>

            {/* User Profile with Dropdown */}
            <Dropdown className="hover:bg-neutral-500 rounded-xl">
              <DropdownButton className="flex items-center gap-3 px-3 py-2 rounded-xl cursor-pointer transition-colors active:scale-95">
                {/* Avatar */}
                <div className="w-9 h-9 rounded-lg bg-amber-700 flex items-center justify-center text-white font-medium text-lg">
                  {currentUser.firstname.charAt(0)}
                </div>

                {/* User Info */}
                <div className="text-left">
                  <div className="text-sm font-medium text-white truncate">
                    {currentUser.firstname}
                  </div>
                  <div className="text-xs text-gray-400 truncate">
                    {currentUser.role}
                  </div>
                </div>

                {/* Chevron Icon */}
                <ChevronDownIcon className="w-5 h-5 text-gray-400" />
              </DropdownButton>

              <DropdownMenu className="w-56">
                <DropdownItem>
                  <Link to="/profile" className="flex items-center gap-3 w-full py-1">
                    <svg 
                      xmlns="http://www.w3.org/2000/svg" 
                      viewBox="0 0 24 24" 
                      fill="gray" 
                      className="w-6 h-6"
                    >
                      <path 
                        fillRule="evenodd" 
                        d="M18.685 19.097A9.723 9.723 0 0 0 21.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 0 0 3.065 7.097A9.716 9.716 0 0 0 12 21.75a9.716 9.716 0 0 0 6.685-2.653Zm-12.54-1.285A7.486 7.486 0 0 1 12 15a7.486 7.486 0 0 1 5.855 2.812A8.224 8.224 0 0 1 12 20.25a8.224 8.224 0 0 1-5.855-2.438ZM15.75 9a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z" 
                        clipRule="evenodd" 
                      />
                    </svg>
                    <span>My Profile</span>
                  </Link>
                </DropdownItem>

                <DropdownDivider />

                <DropdownItem onClick={signOut}>
                  <div className="flex items-center gap-3 py-1">
                    <svg 
                      xmlns="http://www.w3.org/2000/svg" 
                      viewBox="0 0 24 24" 
                      fill="gray" 
                      className="w-6 h-6"
                    >
                      <path 
                        fillRule="evenodd" 
                        d="M7.5 3.75A1.5 1.5 0 0 0 6 5.25v13.5a1.5 1.5 0 0 0 1.5 1.5h6a1.5 1.5 0 0 0 1.5-1.5V15a.75.75 0 0 1 1.5 0v3.75a3 3 0 0 1-3 3h-6a3 3 0 0 1-3-3V5.25a3 3 0 0 1 3-3h6a3 3 0 0 1 3 3V9A.75.75 0 0 1 15 9V5.25a1.5 1.5 0 0 0-1.5-1.5h-6Zm10.72 4.72a.75.75 0 0 1 1.06 0l3 3a.75.75 0 0 1 0 1.06l-3 3a.75.75 0 1 1-1.06-1.06l1.72-1.72H9a.75.75 0 0 1 0-1.5h10.94l-1.72-1.72a.75.75 0 0 1 0-1.06Z" 
                        clipRule="evenodd" 
                      />
                    </svg>
                    <span className="text-red-600 font-medium">ออกจากระบบ</span>
                  </div>
                </DropdownItem>
              </DropdownMenu>
            </Dropdown>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default UserNavbar;