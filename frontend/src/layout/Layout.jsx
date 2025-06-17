import Header from "../components/Header/Header.jsx";

export default function Layout({ children }) {
  return (
    <>
      <Header />
      <main style={{ padding: "2rem" }}>{children}</main>
    </>
  );
}
