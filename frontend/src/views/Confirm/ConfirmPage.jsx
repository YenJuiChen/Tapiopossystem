import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import './ConfirmPage.css'

function ConfirmPage() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const itemId = params.get('item_id')
  const code = params.get('code')

  const [product, setProduct] = useState(null)
  const [isCustomAmountItem, setIsCustomAmountItem] = useState(false)

  const [form, setForm] = useState({
    name: '',
    gender: '',
    phone: '',
    address: '',
    quantity: 1,
    customAmount: '',
    need_certificate: true,
    offering_taken_back: true,
  })

  const [step, setStep] = useState(1)
  const [paymentMethod, setPaymentMethod] = useState('')

  // 讀取從會員搜尋頁帶入的會員資料
  useEffect(() => {
    const stored = sessionStorage.getItem('selectedMember')
    if (stored) {
      try {
        const member = JSON.parse(stored)
        setForm(prev => ({
          ...prev,
          name: member.name || '',
          gender: member.gender || '',
          phone: member.phone || '',
          address: member.address || '',
        }))
      } catch (err) {
        console.error('failed to parse selectedMember', err)
      } finally {
        sessionStorage.removeItem('selectedMember')
      }
    }
  }, [])

  useEffect(() => {
    const fetchItem = async () => {
      try {
        const res = await fetch(`/api/item-detail?item_id=${itemId}`)
        const json = await res.json()
        setProduct(json.data)
        setForm(prev => ({ ...prev, need_certificate: json.data.is_print }))
      } catch (err) {
        console.error('載入商品資訊失敗', err)
      }
    }

    if (itemId) {
      fetchItem()
      if (itemId === '20') {
        setIsCustomAmountItem(true)
      }
    }
  }, [itemId])

  const formatNumber = (num) => Number(num).toLocaleString()

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target
    if (name === 'phone' || name === 'customAmount') {
      const onlyNumbers = value.replace(/\D/g, '')
      setForm(prev => ({ ...prev, [name]: onlyNumbers }))
    } else {
      setForm(prev => ({
        ...prev,
        [name]: type === 'checkbox' ? checked : value,
      }))
    }
  }

  const handleGenderSelect = (gender) => {
    setForm(prev => ({ ...prev, gender }))
  }

  const increaseQuantity = () => {
    setForm(prev => ({ ...prev, quantity: Math.max(1, prev.quantity + 1) }))
  }

  const decreaseQuantity = () => {
    setForm(prev => ({ ...prev, quantity: Math.max(1, prev.quantity - 1) }))
  }

  const handleQuantityInput = (e) => {
    const value = e.target.value.replace(/\D/g, '')
    setForm(prev => ({ ...prev, quantity: Math.max(1, parseInt(value || 1)) }))
  }

  const handleNext = () => {
    if (product.is_print) {
      if (!form.name || !form.gender || !form.phone) {
        alert('請填寫姓名、性別、電話')
        return
      }
    }
    if (isCustomAmountItem && !form.customAmount) {
      alert('請輸入金額')
      return
    }

    setStep(2)
  }

  const handleBack = () => setStep(1)

  const handleSubmit = async () => {
    if (!paymentMethod) {
      alert('請選擇付款方式')
      return
    }

    const quantity = parseInt(form.quantity, 10) || 1
    const amount = isCustomAmountItem
      ? parseInt(form.customAmount || '0', 10)
      : product.price * quantity

    const payload = {
      name: form.name,
      gender: form.gender,
      phone: form.phone,
      address: form.address,
      quantity: isCustomAmountItem ? 1 : quantity,
      amount,
      need_certificate: form.need_certificate,
      info: product.category_id === 5
        ? (form.offering_taken_back ? '供品領回' : '供品捐贈')
        : '',
      category: product.category_name,
      product_id: product.id,
      product_name: product.name,
      price: product.price,
      payment_method: paymentMethod,
      code: code || '',
    }

    try {
      const res = await fetch('/api/record', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })

      const result = await res.json()
      if (result.status === 'success') {
        alert('✅ 訂單已送出')
        navigate('/')
      } else {
        alert('❌ 訂單送出失敗')
      }
    } catch (err) {
      console.error(err)
      alert('發生錯誤')
    }
  }

  if (!product) return <p>載入中...</p>

  const total = isCustomAmountItem
    ? parseInt(form.customAmount || '0', 10)
    : product.price * form.quantity

  return (
    <div className="confirm-page">
      <h2>確認資料</h2>

      <div className="product-info">
        <p><strong>分類：</strong>{product.category_name}</p>
        <p><strong>商品：</strong>{product.name}</p>
        {!isCustomAmountItem && (
          <p><strong>單價：</strong>${formatNumber(product.price)}</p>
        )}
        {!isCustomAmountItem && (
          <p><strong>數量：</strong>{form.quantity}</p>
        )}
        <p><strong>總金額：</strong><strong>${formatNumber(total)}</strong></p>
      </div>

      {isCustomAmountItem && (
        <div className="form-group">
          <label>金額（手動輸入）</label>
          <input
            name="customAmount"
            value={form.customAmount}
            onChange={handleChange}
            className="form-input"
            placeholder="請輸入金額"
          />
        </div>
      )}

      {step === 1 && (
        <>
          {product.is_print && (
            <>
              <div className="form-group">
                <label>姓名</label>
                <input name="name" value={form.name} onChange={handleChange} className="form-input" />
              </div>

              <div className="form-group">
                <label>性別</label>
                <div className="pos-buttons">
                  <button className={form.gender === '男' ? 'active' : ''} onClick={() => handleGenderSelect('男')}>男</button>
                  <button className={form.gender === '女' ? 'active' : ''} onClick={() => handleGenderSelect('女')}>女</button>
                </div>
              </div>

              <div className="form-group">
                <label>電話</label>
                <input
                  name="phone"
                  type="tel"
                  pattern="[0-9]*"
                  inputMode="numeric"
                  value={form.phone}
                  onChange={handleChange}
                  className="form-input"
                />
              </div>

              <div className="form-group">
                <label>地址</label>
                <input name="address" value={form.address} onChange={handleChange} className="form-input" />
              </div>
            </>
          )}

          {!isCustomAmountItem && (
            <div className="form-group">
              <label>數量</label>
              <div className="quantity-control">
                <button onClick={decreaseQuantity}>-</button>
                <input
                  type="text"
                  name="quantity"
                  value={form.quantity}
                  onChange={handleQuantityInput}
                  className="quantity-input"
                />
                <button onClick={increaseQuantity}>+</button>
              </div>
            </div>
          )}

          {product.is_print && (
            <div className="form-group checkbox-group">
              <label className="pretty-checkbox">
                <input
                  type="checkbox"
                  name="need_certificate"
                  checked={form.need_certificate}
                  onChange={handleChange}
                />
                <span className="custom-checkmark" />
                <span className="checkbox-label">是否需要感謝狀</span>
              </label>
            </div>
          )}

          {product.is_print && product.category_id === 5 && (
            <div className="form-group checkbox-group">
              <label className="pretty-checkbox">
                <input
                  type="checkbox"
                  name="offering_taken_back"
                  checked={form.offering_taken_back}
                  onChange={handleChange}
                />
                <span className="custom-checkmark" />
                <span className="checkbox-label">供品是否領回</span>
              </label>
            </div>
          )}

          <button className="next-button" onClick={handleNext}>下一步</button>
        </>
      )}

      {step === 2 && (
        <>
            <div className="form-group">
              <label>付款方式</label>
              <div className="pos-buttons">
                <button
                  className={paymentMethod === 'cash' ? 'active' : ''}
                  onClick={() => setPaymentMethod('cash')}
                >
                  現金
                </button>
                <button className="disabled" disabled>
                  LINE Pay（尚未開放）
                </button>
              </div>
            </div>

          <div className="button-row">
            <button className="back-button" onClick={handleBack}>← 上一步</button>
            <button className="submit-button" onClick={handleSubmit}>送出</button>
          </div>
        </>
      )}
    </div>
  )
}

export default ConfirmPage
