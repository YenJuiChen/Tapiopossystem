import { useMemo, useState } from 'react'
import Table from '../../components/Table'
import { TableRow, TableCell } from '@mui/material'
import './MemberSearchPage.css'

export default function MemberSearchPage() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState([])
  const [expandedRows, setExpandedRows] = useState([])
  const [ordersMap, setOrdersMap] = useState({})

  const handleClone = async (e, id) => {
    e.stopPropagation()
    try {
      const res = await fetch(`/api/records/${id}/clone`, { method: 'POST' })
      const json = await res.json()
      if (json.status === 'success') {
        alert('訂單成功送出')
      } else {
        alert('訂單送出失敗')
      }
    } catch (err) {
      console.error('clone failed', err)
    }
  }

  const columns = useMemo(
    () => [
      { accessorKey: 'id', header: 'ID' },
      { accessorKey: 'name', header: '姓名' },
      { accessorKey: 'gender', header: '性別' },
      { accessorKey: 'phone', header: '電話' },
      { accessorKey: 'address', header: '地址' },
      {
        accessorKey: 'action',
        header: '',
        cell: (_, row) => (
          <button className="print-button" onClick={() => handleToggleOrders(row)}>
            ☰
          </button>
        ),
      },
    ],
    []
  )

  const handleSearch = async () => {
    if (!query) return
    try {
      const res = await fetch(`/api/members?query=${encodeURIComponent(query)}`)
      const json = await res.json()
      setResults(Array.isArray(json) ? json : json.data || [])
    } catch (err) {
      console.error('Failed to fetch members', err)
    }
  }


  const handleToggleOrders = async (member) => {
    const key = member.id
    if (expandedRows.includes(key)) {
      setExpandedRows((prev) => prev.filter((id) => id !== key))
      return
    }
    const res = await fetch(
      `/api/member-orders?phone=${encodeURIComponent(member.phone)}&name=${encodeURIComponent(member.name)}`
    )
    const json = await res.json()
    setOrdersMap((prev) => ({
      ...prev,
      [key]: Array.isArray(json) ? json : json.data || [],
    }))
    setExpandedRows((prev) => [...prev, key])
  }

  return (
    <div className="record-list-container">
      <h2>會員搜尋</h2>
      <div className="record-header">
        <div className="search-bar">
          <input
            type="text"
            placeholder="輸入會員關鍵字"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <button onClick={handleSearch}>搜尋</button>
        </div>
      </div>
      <Table
        data={results}
        columns={columns}
        renderExpandedRow={(row) => {
          if (!expandedRows.includes(row.id)) return null
          const orders = ordersMap[row.id] || []
          const content = orders.length === 0 ? (
            <p>無訂單</p>
          ) : (
            <table className="sub-table">
              <thead>
                <tr>
                  <th>日期</th>
                  <th>項目</th>
                  <th>金額</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {orders.map((o) => (
                  <tr key={o.id}>
                    <td>{o.created_at}</td>
                    <td>{o.product_name}</td>
                    <td>{o.amount}</td>
                    <td>
                      <button
                        className="print-button"
                        onClick={(e) => handleClone(e, o.id)}
                      >
                        +
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )
          return (
            <TableRow>
              <TableCell colSpan={columns.length}>{content}</TableCell>
            </TableRow>
          )
        }}
      />
    </div>
  )
}