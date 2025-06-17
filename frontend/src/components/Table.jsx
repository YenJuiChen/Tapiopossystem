import React, { useState } from "react";
import {
  Table as MuiTable,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TableSortLabel,
  Paper,
} from "@mui/material";

function Table({
  data = [],
  columns = [],
  onSortingChange,
  isPrint = false,
  onRowClick,
  renderExpandedRow,
}) {
  const [sortConfig, setSortConfig] = useState({
    column: "created_at",
    direction: "desc",
  });

  const handleSort = (column) => {
    let direction = "asc";
    if (sortConfig.column === column && sortConfig.direction === "asc") {
      direction = "desc";
    }
    setSortConfig({ column, direction });
    if (onSortingChange) {
      onSortingChange(column, direction);
    }
  };

  const safeData = Array.isArray(data) ? data : [];

  if (isPrint) {
    // ✅ 純 HTML table 供列印用
    return (
      <div className="print-table-wrapper">
        <table className="print-table">
          <thead>
            <tr>
              {columns.map((col) => (
                <th key={col.accessorKey}>{col.header}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {safeData.length === 0 ? (
              <tr>
                <td colSpan={columns.length} style={{ textAlign: "center" }}>
                  無資料
                </td>
              </tr>
            ) : (
              safeData.map((row, idx) => (
                <tr key={idx}>
                  {columns.map((col) => (
                    <td key={col.accessorKey}>
                      {col.cell
                        ? col.cell(row[col.accessorKey], row)
                        : row[col.accessorKey]}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    );
  }

  // ✅ 原本 MUI Table（網頁版用）
  return (
    <TableContainer component={Paper}>
      <MuiTable>
        <TableHead>
          <TableRow>
            {columns.map((col) => (
              <TableCell key={col.accessorKey}>
                {col.enableSorting ? (
                  <TableSortLabel
                    active={sortConfig.column === col.accessorKey}
                    direction={sortConfig.direction}
                    onClick={() => handleSort(col.accessorKey)}
                  >
                    {col.header}
                  </TableSortLabel>
                ) : (
                  col.header
                )}
              </TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {safeData.length === 0 ? (
            <TableRow>
              <TableCell colSpan={columns.length} align="center">
                無資料
              </TableCell>
            </TableRow>
          ) : (
            safeData.map((row, idx) => (
              <React.Fragment key={idx}>
                <TableRow
                  onClick={onRowClick ? () => onRowClick(row) : undefined}
                >
                  {columns.map((col) => (
                    <TableCell key={col.accessorKey}>
                      {col.cell
                        ? col.cell(row[col.accessorKey], row)
                        : row[col.accessorKey]}
                    </TableCell>
                  ))}
                </TableRow>
                {renderExpandedRow ? renderExpandedRow(row) : null}
              </React.Fragment>
            ))
          )}
        </TableBody>
      </MuiTable>
    </TableContainer>
  );
}

export default Table;