import { useState } from 'react'

export default function MemberSearchPage() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState([])

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
    console.log('新增服務', member)
  }

  return (
    <div className="member-search-page">
      <div style={{ marginBottom: '1rem' }}>
        <input
          type="text"
          placeholder="輸入會員關鍵字"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          style={{ marginRight: '0.5rem' }}
        />
        <button onClick={handleSearch}>搜尋</button>
      </div>
      <ul>
        {results.map((m) => (
          <li key={m.id} style={{ marginBottom: '0.5rem' }}>
            <span style={{ marginRight: '1rem' }}>{m.name || m.id}</span>
            <button onClick={() => handleAddService(m)}>新增服務</button>
          </li>
        ))}
      </ul>
    </div>
  )
}