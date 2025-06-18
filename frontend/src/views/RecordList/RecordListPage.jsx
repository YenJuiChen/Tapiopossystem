import { useEffect, useState, useMemo } from "react";
import FilterBar from "../../components/FilterBar";
import Pagination from "@mui/material/Pagination";
import Table from "../../components/Table";
import "./RecordListPage.css";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";

dayjs.extend(utc);
dayjs.extend(timezone);

function RecordListPage() {
  const today = dayjs().format("YYYY-MM-DD");
  const [filters, setFilters] = useState({
    category: "",
    item: "",
    startDate: today,
    endDate: today,
    startTime: "06:00",
    endTime: "21:00",
  });

  const [sortBy, setSortBy] = useState("created_at");
  const [order, setOrder] = useState("desc");
  const [data, setData] = useState([]);
  const [total, setTotal] = useState(0);
  const [options, setOptions] = useState([]);
  const [page, setPage] = useState(1);

  const pageSize = 20;

  const buildQuery = (override = {}) => {
    const params = new URLSearchParams({
      category: filters.category,
      item: filters.item,
      startDate: `${filters.startDate} ${filters.startTime}`,
      endDate: `${filters.endDate} ${filters.endTime}`,
      sortBy,
      order,
      limit: override.limit ?? pageSize,
      offset: override.offset ?? (page - 1) * pageSize,
    });
    return params.toString();
  };

  const fetchOptions = async () => {
    try {
      const res = await fetch("/api/category-items");
      const json = await res.json();
      setOptions(json || []);
    } catch (err) {
      console.error("Failed to fetch /api/category-items", err);
    }
  };

  const fetchData = async () => {
    try {
      const res = await fetch(`/api/records?${buildQuery()}`);
      const json = await res.json();
      setData(Array.isArray(json.data) ? json.data : []);
      setTotal(json.total || 0);
    } catch (err) {
      console.error("Failed to fetch /api/records", err);
      setData([]);
    }
  };

  const generatePrintableHTML = (data) => {
    const headers = columns.map((col) => `<th>${col.header}</th>`).join("");
    const rows = data
      .map(
        (row) => `
      <tr>
        ${columns
          .map((col) => {
            const value = col.cell
              ? col.cell(row[col.accessorKey], row)
              : (row[col.accessorKey] ?? "");
            return `<td>${value}</td>`;
          })
          .join("")}
      </tr>
    `,
      )
      .join("");

    return `
      <html>
        <head>
          <meta charset="UTF-8">
          <title>列印報表</title>
          <style>
            body {
              font-family: sans-serif;
              margin: 20px;
            }
            h2 {
              text-align: center;
              margin-bottom: 20px;
            }
            table {
              width: 100%;
              border-collapse: collapse;
              font-size: 12px;
            }
            th, td {
              border: 1px solid #000;
              padding: 6px;
              text-align: left;
              word-break: break-word;
            }
            th {
              background-color: #f0f0f0;
            }
            @page {
              size: A4;
              margin: 20mm;
            }
          </style>
        </head>
        <body>
          <h2>報表列印</h2>
          <table>
            <thead><tr>${headers}</tr></thead>
            <tbody>${rows}</tbody>
          </table>
          <script>window.onload = () => setTimeout(() => window.print(), 300)</script>
        </body>
      </html>
    `;
  };

  const handlePrint = async () => {
    try {
      const res = await fetch(
        `/api/records?${buildQuery({ limit: 9999, offset: 0 })}`,
      );
      const json = await res.json();
      const all = Array.isArray(json.data) ? json.data : [];
      const html = generatePrintableHTML(all);

      const win = window.open("", "_blank");
      if (win) {
        win.document.write(html);
        win.document.close();
      } else {
        alert("無法開啟列印視窗，請確認瀏覽器未封鎖彈出視窗");
      }
    } catch (err) {
      console.error("列印時發生錯誤", err);
    }
  };

  useEffect(() => {
    fetchOptions();
  }, []);

  useEffect(() => {
    fetchData();
  }, [filters, sortBy, order, page]);

  const columns = useMemo(
    () => [
      {
        accessorKey: "created_at",
        header: "日期",
        enableSorting: true,
        cell: (value) =>
          value
            ? dayjs(value).tz("Asia/Taipei").format("YYYY-MM-DD HH:mm:ss")
            : "",
      },
      { accessorKey: "name", header: "姓名", cell: (v) => v ?? "" },
      { accessorKey: "gender", header: "性別", cell: (v) => v ?? "" },
      { accessorKey: "phone", header: "電話", cell: (v) => v ?? "" },
      {
        accessorKey: "amount",
        header: "金額",
        cell: (v) => (typeof v === "number" ? v.toLocaleString("en-US") : ""),
      },
      { accessorKey: "product_name", header: "項目", cell: (v) => v ?? "" },
      { accessorKey: "category", header: "種類", cell: (v) => v ?? "" },
      { accessorKey: "address", header: "地址", cell: (v) => v ?? "" },
      {
        accessorKey: "payment_method",
        header: "付款方式",
        cell: (v) => v ?? "",
      },
      {
        accessorKey: "info",
        header: "資訊",
        cell: (v) =>
          v == null
            ? ""
            : typeof v === "object"
              ? JSON.stringify(v)
              : String(v),
      },
    ],
    [],
  );

  return (
    <div className="record-list-container">
      <h2>報表記錄</h2>

      <div className="record-header">
        <FilterBar
          filters={filters}
          setFilters={setFilters}
          options={options}
        />
        <button className="print-button" onClick={handlePrint}>
          列印報表
        </button>
      </div>

      <div className="pagination-top">
        <Pagination
          count={Math.ceil(total / pageSize)}
          page={page}
          onChange={(e, val) => setPage(val)}
        />
      </div>

      <Table
        data={data}
        columns={columns}
        onSortingChange={(column, direction) => {
          setSortBy(column);
          setOrder(direction);
        }}
      />
    </div>
  );
}

export default RecordListPage;