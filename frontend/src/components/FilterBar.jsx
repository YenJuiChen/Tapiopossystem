import React from 'react'
import { TextField, MenuItem, Grid } from '@mui/material'

const FilterBar = ({ filters, setFilters, options }) => {
  const handleChange = (e) => {
    const { name, value } = e.target
    setFilters((prev) => ({ ...prev, [name]: value }))
  }

  const selectedCategory = options.find(cat => cat.name === filters.category)
  const availableItems = selectedCategory ? selectedCategory.items : []

  return (
    <Grid container spacing={2} marginBottom={2}>
      <Grid item xs={12} sm={2}>
        <TextField
          select
          label="分類"
          name="category"
          value={filters.category}
          onChange={handleChange}
          fullWidth
        >
          <MenuItem value="">全部</MenuItem>
          {options.map((cat) => (
            <MenuItem key={cat.id} value={cat.name}>
              {cat.name}
            </MenuItem>
          ))}
        </TextField>
      </Grid>

      <Grid item xs={12} sm={2}>
        <TextField
          select
          label="項目"
          name="item"
          value={filters.item}
          onChange={handleChange}
          fullWidth
        >
          <MenuItem value="">全部</MenuItem>
          {availableItems.map((item) => (
            <MenuItem key={item.id} value={item.name}>
              {item.name}
            </MenuItem>
          ))}
        </TextField>
      </Grid>

      <Grid item xs={12} sm={2}>
        <TextField
          label="開始日期"
          type="date"
          name="startDate"
          value={filters.startDate}
          onChange={handleChange}
          fullWidth
          InputLabelProps={{ shrink: true }}
        />
      </Grid>

      <Grid item xs={12} sm={2}>
        <TextField
          label="開始時間"
          type="time"
          name="startTime"
          value={filters.startTime}
          onChange={handleChange}
          fullWidth
          InputLabelProps={{ shrink: true }}
        />
      </Grid>

      <Grid item xs={12} sm={2}>
        <TextField
          label="結束日期"
          type="date"
          name="endDate"
          value={filters.endDate}
          onChange={handleChange}
          fullWidth
          InputLabelProps={{ shrink: true }}
        />
      </Grid>

      <Grid item xs={12} sm={2}>
        <TextField
          label="結束時間"
          type="time"
          name="endTime"
          value={filters.endTime}
          onChange={handleChange}
          fullWidth
          InputLabelProps={{ shrink: true }}
        />
      </Grid>
    </Grid>
  )
}

export default FilterBar
