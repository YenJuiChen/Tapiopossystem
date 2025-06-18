import React, { useEffect, useState, useRef } from "react";
import QRCode from "react-qr-code";
import { v4 as uuidv4 } from "uuid";
import html2canvas from "html2canvas";
import "./QRCodePage.css";

export default function QRCodePage() {
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState("");
  const [selectedItemId, setSelectedItemId] = useState("");
  const [items, setItems] = useState([]);
  const [qrcodes, setQrcodes] = useState([]);
  const previewRefs = useRef([]);

  const selectedCategoryName =
    categories.find((c) => String(c.id) === selectedCategory)?.name || "";
  const selectedItemName =
    items.find((i) => String(i.id) === selectedItemId)?.name || "";

  useEffect(() => {
    fetchCategories();
  }, []);

  const fetchCategories = async () => {
    try {
      const res = await fetch("/api/category-items?qrcodeOnly=1");
      const data = await res.json();

      const filteredCategories = data
        .map((cat) => {
          const validItems = (cat.items || []).filter((item) => item.is_qrcode);
          if (validItems.length === 0) return null;
          return {
            ...cat,
            items: validItems,
          };
        })
        .filter(Boolean);

      setCategories(filteredCategories);
    } catch (err) {
      console.error("載入分類失敗", err);
    }
  };

  const handleCategoryChange = (e) => {
    const catId = e.target.value;
    setSelectedCategory(catId);
    setSelectedItemId("");
    const found = categories.find((cat) => String(cat.id) === catId);
    setItems(found?.items || []);
  };

  const handleItemChange = (e) => {
    setSelectedItemId(e.target.value);
  };

  const generateQRCodes = () => {
    if (!selectedItemId) return;
    const baseUrl = window.location.origin;
    const newQRCodes = Array.from({ length: 50 }, () => {
      const uuid = uuidv4();
      return {
        uuid,
        url: `${baseUrl}/confirm?item_id=${selectedItemId}&code=${uuid}`,
      };
    });
    setQrcodes(newQRCodes);
    previewRefs.current = new Array(50).fill(null);
  };

  const handlePrint = async () => {
    if (!qrcodes.length) return;

    const label = `${selectedCategoryName} - ${selectedItemName}`;
    const images = await Promise.all(
      previewRefs.current.map(async (el) => {
        if (!el) return "";
        const canvas = await html2canvas(el);
        return canvas.toDataURL("image/png");
      }),
    );

    const html = `
      <html>
        <head>
          <meta charset="UTF-8">
          <title>列印 QRCode</title>
          <style>
            body {
              margin: 0;
              padding: 2mm;
              font-family: Arial, sans-serif;
              display: grid;
              grid-template-columns: repeat(5, 1fr);
              gap: 1mm;
              box-sizing: border-box;
              justify-items: center;
              align-items: center;
            }
            img {
              width: 130px;
              height: auto;
              display: block;
            }
            @page {
              size: A4 portrait;
              margin: 0;
            }
            @media print {
              html, body {
                width: 210mm;
                height: 297mm;
                overflow: hidden;
              }
            }
          </style>
        </head>
        <body>
          ${images.map((src) => (src ? `<img src="${src}" />` : "")).join("")}
          <script>window.onload = function() { setTimeout(() => window.print(), 300); }</script>
        </body>
      </html>
    `;

    const printWindow = window.open("", "", "width=900,height=1200");
    if (printWindow) {
      printWindow.document.write(html);
      printWindow.document.close();
    } else {
      alert("無法開啟列印視窗，請確認瀏覽器未封鎖彈出視窗");
    }
  };

  return (
    <div className="qrcode-page">
      <h2>產生 QRCODE</h2>

      <div className="qrcode-selects">
        <select value={selectedCategory} onChange={handleCategoryChange}>
          <option value="">請選擇分類</option>
          {categories.map((cat) => (
            <option key={cat.id} value={cat.id}>
              {cat.name}
            </option>
          ))}
        </select>

        <select
          value={selectedItemId}
          onChange={handleItemChange}
          disabled={!selectedCategory}
        >
          <option value="">請選擇項目</option>
          {items.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
            </option>
          ))}
        </select>
      </div>

      {selectedItemId && (
        <div className="qrcode-buttons">
          <button onClick={generateQRCodes}>重新產生 50 組 QRCODE</button>
          <button onClick={handlePrint}>列印</button>
        </div>
      )}

      {qrcodes.length > 0 && (
        <div className="qrcode-grid-preview">
          {qrcodes.map((q, idx) => (
            <div
              className="qrcode-box"
              key={idx}
              ref={(el) => (previewRefs.current[idx] = el)}
              style={{ padding: "2px", width: "135px", textAlign: "center" }}
            >
              <QRCode value={q.url} size={125} />
              <div style={{ fontSize: 9, marginTop: 2 }}>
                {selectedCategoryName} - {selectedItemName}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
