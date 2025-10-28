"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import Link from "next/link";

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await fetch("/api/login", {
        // Next.js 서버를 통해 로그인 요청
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (res.ok) {
        const data = await res.json();
        login(data.token); // 성공 시 AuthContext에 토큰 저장
        router.push("/"); // 홈으로 이동
      } else {
        alert("로그인에 실패했습니다. 아이디와 비밀번호를 확인해주세요.");
      }
    } catch (error) {
      alert("로그인 중 오류가 발생했습니다.");
    }
  };

  return (
    <main
      style={{
        fontFamily: "sans-serif",
        padding: "2rem",
        maxWidth: "400px",
        margin: "100px auto",
      }}
    >
      <h1 style={{ textAlign: "center" }}>로그인</h1>
      <form
        onSubmit={handleSubmit}
        style={{
          marginTop: "2rem",
          display: "flex",
          flexDirection: "column",
          gap: "1rem",
        }}
      >
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="Username"
          required
          style={{
            padding: "0.75rem",
            borderRadius: "8px",
            border: "1px solid #ddd",
          }}
        />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Password"
          required
          style={{
            padding: "0.75rem",
            borderRadius: "8px",
            border: "1px solid #ddd",
          }}
        />
        <button
          type="submit"
          style={{
            cursor: "pointer",
            background: "#0070f3",
            color: "white",
            padding: "0.75rem",
            borderRadius: "8px",
            border: "none",
            fontSize: "1rem",
          }}
        >
          로그인
        </button>
      </form>
      <p style={{ textAlign: "center", marginTop: "1rem" }}>
        계정이 없으신가요?{" "}
        <Link href="/signup" style={{ color: "#0070f3" }}>
          회원가입
        </Link>
      </p>
    </main>
  );
}
