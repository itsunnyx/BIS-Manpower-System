import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);
  const [error, setError] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();

  const togglePasswordVisibility = () => {
    setShowPassword(prev => !prev);
  };

  const handleLogin = async (e) => {
  e.preventDefault();
  setError('');

  try {
    // เรียก API Backend
    const response = await fetch('http://localhost:8080/api/v1/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();

    if (!response.ok) {
      const errorMessage = data.message || 'Authentication Fail: Please check user or Password';
      setError(errorMessage);
      return;
    }

    // Login สำเร็จ
    const { token, role } = data;
    localStorage.setItem('jwt_token', token);
    localStorage.setItem('user_role', role);

    // Navigate ตาม role
    switch (role.toLowerCase()) {
      case 'admin':
        navigate('/admin');             
        break;
      case 'approve':
        navigate('/approver');          
        break;
      case 'user':
        navigate('/user');              
        break;
      default:
        navigate('/');                 
        break;
    }

  } catch (err) {
    console.error('Login error:', err);
    setError('An unexpected error occurred. Please try again later.');
  }
};

  return (
    <div className="flex items-center justify-center min-h-screen bg-white p-4">
      <div className="w-full max-w-sm p-8 space-y-6">
        <div className="text-center">
          {/* Logo - ใช้รูป nakla.svg */}
          <div className="mx-auto w-40 mb-4">
            <img
              src="/images/nakla.svg"
              alt="Nakla Manpower Logo"
              className="w-full h-auto"
            />
          </div>
          <h2 className="text-3xl font-light text-gray-800 mt-4 mb-8">
            Log in
          </h2>
        </div>

        <form className="space-y-4" onSubmit={handleLogin}>
          {/* Email Field */}
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700">Email</label>
            <input
              id="email"
              type="text"
              required
              placeholder="Enter your email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-gray-500 focus:border-gray-500 sm:text-sm"
            />
          </div>

          {/* Password Field */}
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700">Password</label>
            <div className="mt-1 relative">
              <input
                id="password"
                type={showPassword ? "text" : "password"}
                required
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="block w-full pr-10 pl-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-gray-500 focus:border-gray-500 sm:text-sm"
              />

              <span
                onClick={togglePasswordVisibility}
                className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 cursor-pointer hover:text-gray-600 transition duration-150"
              >
                {showPassword ? (
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.477 0 8.268 2.943 9.542 7-1.274 4.057-5.065 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                ) : (
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.477 0 8.268 2.943 9.542 7-1.274 4.057-5.065 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                    <line x1="3" y1="3" x2="21" y2="21" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"></line>
                  </svg>
                )}
              </span>
            </div>
          </div>

          {/* Remember Me */}
          <div className="flex items-center">
            <input
              id="remember-me"
              name="remember-me"
              type="checkbox"
              checked={rememberMe}
              onChange={(e) => setRememberMe(e.target.checked)}
              className="h-4 w-4 text-gray-600 border-gray-300 rounded focus:ring-gray-500"
            />
            <label htmlFor="remember-me" className="ml-2 block text-sm text-gray-900">
              Remeber me
            </label>
          </div>

          {/* Sign In Button */}
          <div>
            <button
              type="submit"
              className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition duration-150"
            >
              Sign in
            </button>
          </div>

          {/* Error Message */}
          {error && (
            <p className="text-red-500 text-xs font-semibold text-center mt-2">
              {error}
            </p>
          )}
        </form>
      </div>
    </div>
  );
};

export default Login;
