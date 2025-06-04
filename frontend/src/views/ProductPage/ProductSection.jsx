import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import './ProductSection.css'

export default function ProductSection({ categoryId, categoryName }) {
  const [products, setProducts] = useState([])
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const res = await fetch(`/api/items?category_id=${categoryId}`)
        const data = await res.json()
        setProducts(data)
      } catch (err) {
        console.error('載入商品失敗', err)
      } finally {
        setLoading(false)
      }
    }

    fetchProducts()
  }, [categoryId])

  const handleClick = (product) => {
    navigate(`/confirm?item_id=${product.id}`)
  }

  return (
    <section className="product-section">
      <h2 className="section-title">{categoryName}</h2>
      {loading ? (
        <p className="loading-text">載入中...</p>
      ) : (
        <div className="product-grid">
          {products.map((p) => (
            <div key={p.id} className="product-card" onClick={() => handleClick(p)}>
              <div className="product-name">{p.name}</div>
              <div className="product-price">${p.price}</div>
            </div>
          ))}
        </div>
      )}
    </section>
  )
}
