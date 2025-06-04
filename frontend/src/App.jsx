import { Routes, Route, Navigate, Outlet, useLocation } from 'react-router-dom'
import Home from './views/Home/Home'
import ProductPage from './views/ProductPage/ProductPage'
import ConfirmPage from './views/Confirm/ConfirmPage'
import RecordListPage from './views/RecordList/RecordListPage'
import PrintPage from './views/Print/PrintPage'
import QRCodePage from './views/QRCode/QRCodePage'
import Layout from './layout/Layout'
import LoginPage from './views/LoginPage/LoginPage'
import Cookies from 'js-cookie'

function ProtectedRoute() {
  const token = Cookies.get('auth_token')
  const location = useLocation()

  if (location.pathname.startsWith('/confirm')) {
    return <Outlet />
  }

  if (!token) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />
  }

  return <Outlet />
}

function App() {
  return (
    <Routes>
      {/* 🔐 登入頁面 */}
      <Route path="/login" element={<LoginPage />} />

      {/* 🔐 受保護頁面 */}
      <Route element={<ProtectedRoute />}>
        {/* 🏠 首頁 */}
        <Route
          path="/"
          element={
            <Layout>
              <Home />
            </Layout>
          }
        />

        {/* 📦 商品列表 */}
        <Route
          path="/items"
          element={
            <Layout>
              <ProductPage />
            </Layout>
          }
        />

        {/* 📊 紀錄報表 */}
        <Route
          path="/records"
          element={
            <Layout>
              <RecordListPage
                defaultSortBy="created_at"
                defaultSortOrder="desc"
                defaultMonth={new Date().toISOString().slice(0, 7)}
              />
            </Layout>
          }
        />

        {/* 📱 自助結帳用 QRCode 產生頁面 */}
        <Route
          path="/generate-qrcode"
          element={
            <Layout>
              <QRCodePage />
            </Layout>
          }
        />
      </Route>

      {/* ✅ 確認頁面（不需登入） */}
      <Route path="/confirm" element={<ConfirmPage />} />

      {/* 🖨️ 列印頁（獨立） */}
      <Route path="/print" element={<PrintPage />} />
    </Routes>
  )
}

export default App
