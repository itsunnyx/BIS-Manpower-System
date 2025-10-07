import { createBrowserRouter, RouterProvider } from "react-router-dom";

import AMainpage from "../pages/approver/AMainpage";
import Approve from "../pages/approver/Approve";
import ADashboard from "../pages/approver/Dashboard"

import Usermainpage from "../pages/user/Usermainpage";


import Usermanage from "../pages/admin/Usermanage";

import Userlayout from "../layouts/Userlayout";
import Approverlayout from "../layouts/Approverlayout";
import Adminlayout from "../layouts/Adminlayout";

import Notfound from "../pages/Notfound";
import UserRForm from "../pages/user/UserRForm";

import Login from "../pages/auth/Login";
import Loginlayout from "../layouts/Loginlayout";



const router = createBrowserRouter([
     // {
     //      path: "/",
     //      element: <Loginlayout />,
     //      children: [
     //           { index: true, element: <Login /> },
     //           { path: "*", element: <Notfound /> },
     //      ]
     // },
     {
          path: '/user',
          element: <Userlayout />,
          children: [
               { index: true, element: <Usermainpage /> },
               { path: "requestform", element: <UserRForm /> },
               { path: "*", element: <Notfound /> },
          ]

     },
     {
          path: '/approver',
          element: <Approverlayout />,
          children: [
               { index: true, element: <AMainpage /> },
               { path: "dashboard", element: <ADashboard /> },
               { path: "approve", element: <Approve /> },
               { path: "*", element: <Notfound /> },
          ]

     },
     {
          path: '/admin',
          element: <Adminlayout />,
          children: [
               { index: true, element: <Usermanage /> },
               { path: "*", element: <Notfound /> },
          ]

     },
     {
          path: '/login',
          element: <Loginlayout />,
          children: [
               { index: true, element: <Login /> },
               { path: "*", element: <Notfound /> },
          ]

     }

]);

const AppRoutes = () => {
     return <RouterProvider router={router} />;
};

export default AppRoutes;
