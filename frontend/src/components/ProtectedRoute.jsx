// src/components/ProtectedRoute.jsx
import { Navigate, Outlet, useLocation } from 'react-router-dom'
import Cookies from 'js-cookie'

export default function ProtectedRoute() {
  const token = Cookies.get('auth_token')
  const location = useLocation()

  // 允許訪問 /confirm（不需要登入）
  if (location.pathname.startsWith('/confirm')) {
    return <Outlet />
  }

  if (!token) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />
  }

  return <Outlet />
}
