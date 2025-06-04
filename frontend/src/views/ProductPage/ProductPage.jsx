import { useSearchParams } from 'react-router-dom'
import ProductSection from './ProductSection'

export default function ProductPage() {
  const [params] = useSearchParams()
  const categoryId = params.get('category_id')
  const categoryName = params.get('category_name')

  return (
    <div className="product-page">
      <ProductSection categoryId={categoryId} categoryName={categoryName} />
    </div>
  )
}
