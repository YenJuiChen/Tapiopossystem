import { useMemo, useState } from 'react'
import Table from '../../components/Table'
import './MemberSearchPage.css'
import { useNavigate } from 'react-router-dom'

export default function MemberSearchPage() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState([])
  const navigate = useNavigate()

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
          <button className="print-button" onClick={() => handleAddService(row)}>
            +
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

  const handleAddService = (member) => {
    sessionStorage.setItem('selectedMember', JSON.stringify(member))
    navigate('/confirm')
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
      <Table data={results} columns={columns} />
    </div>
  )
}