import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import Cookies from "js-cookie";
import "./LoginPage.css";

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const location = useLocation();

  const handleSubmit = (e) => {
    e.preventDefault();
    if (username === "admin" && password === "hcjfdg666") {
      Cookies.set("auth_token", "secret-token", { expires: 1 });
      const redirectPath = location.state?.from || "/";
      navigate(redirectPath, { replace: true });
    } else {
      setError("帳號或密碼錯誤");
    }
  };

  return (
    <div className="login-container">
      <h2>登入</h2>
      <form onSubmit={handleSubmit} className="login-form">
        <input
          type="text"
          placeholder="帳號"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          type="password"
          placeholder="密碼"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        {error && <p className="error-message">{error}</p>}
        <button type="submit">登入</button>
      </form>
    </div>
  );
}
