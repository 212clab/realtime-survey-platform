"use client";

export default function LoginPage() {
  const handleSocialLogin = (provider: "google") => {
    // 소셜 로그인 시작 API로 리디렉션
    window.location.href = `/api/auth/login/${provider}`;
  };

  return (
    <main
      style={{
        fontFamily: "sans-serif",
        padding: "2rem",
        maxWidth: "400px",
        margin: "100px auto",
        textAlign: "center",
      }}
    >
      <h1>로그인</h1>
      <p style={{ color: "#888", marginTop: "1rem" }}>
        설문조사를 만들려면 로그인이 필요해요. <br />
        소셜 계정으로 간편하게 시작하세요!
      </p>

      <div
        style={{
          marginTop: "2rem",
          display: "flex",
          flexDirection: "column",
          gap: "1rem",
        }}
      >
        <button onClick={() => handleSocialLogin("google")} style={buttonStyle}>
          <img
            src="/icons/googlechrome.svg"
            width={24}
            height={24}
            alt="google"
          />
          Google 계정으로 로그인
        </button>
      </div>
    </main>
  );
}

const buttonStyle: React.CSSProperties = {
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  gap: "0.5rem",
  cursor: "pointer",
  padding: "0.75rem",
  borderRadius: "8px",
  border: "1px solid #ddd",
  fontSize: "1rem",
  background: "white",
};
