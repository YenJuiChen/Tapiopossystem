import { Link, useLocation } from "react-router-dom";
import "./Header.css";

export default function Header() {
  const location = useLocation();

  return (
    <header className="header">
      <nav className="nav">
        <Link to="/" className={location.pathname === "/" ? "active" : ""}>
          首頁
        </Link>
        <Link
          to="/records"
          className={location.pathname === "/records" ? "active" : ""}
        >
          報表
        </Link>
        <Link
          to="/certificate-records"
          className={
            location.pathname === "/certificate-records" ? "active" : ""
          }
        >
          感謝狀紀錄
        </Link>
        <Link
          to="/generate-qrcode"
          className={location.pathname === "/generate-qrcode" ? "active" : ""}
        >
          產商品條碼
        </Link>
      </nav>
    </header>
  );
}
