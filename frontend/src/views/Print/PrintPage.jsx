import { useLocation, useEffect } from 'react'
import './PrintPage.css'

export default function PrintPage() {
  const location = useLocation()
  const data = location.state || {}

  useEffect(() => {
    // 自動列印
    const printAndRedirect = async () => {
      await new Promise((res) => setTimeout(res, 300)) // 等待畫面渲染
      window.print()
      // 回首頁或紀錄頁
      window.location.href = '/'
    }

    printAndRedirect()
  }, [])

  return (
    <div className="print-page">
      <h1>收據</h1>
      <p>分類：{data.category}</p>
      <p>商品：{data.product_name}</p>
      <p>單價：${data.price}</p>
      <p>數量：{data.quantity}</p>
      <p>姓名：{data.name}</p>
      <p>備註：{data.note}</p>
      <p>總金額：${data.total}</p>
    </div>
  )
}
