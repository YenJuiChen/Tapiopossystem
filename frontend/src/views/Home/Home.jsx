import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import './Home.css'

export default function Home() {
  const navigate = useNavigate()
  const [categories, setCategories] = useState([])

  useEffect(() => {
    const fetchCategories = async () => {
      const res = await fetch('/api/category-items')
      const data = await res.json()
      // 加上依照 id 升序排序
      data.sort((a, b) => a.id - b.id)
      setCategories(data)
    }
    fetchCategories()
  }, [])

  const handleClick = (category) => {
    if (category.id === 4) {
      navigate('/confirm?item_id=20')
    } else{
      navigate(`/items?category_id=${category.id}&category_name=${encodeURIComponent(category.name)}`)
    }
  }

  return (
    <div className="home-container">
      <h1 className="title">請選擇分類</h1>
      <div className="category-grid">
        {categories.map((cat) => (
          <button key={cat.id} className="category-button" onClick={() => handleClick(cat)}>
            {cat.name}
          </button>
        ))}
      </div>
    </div>
  )
}
