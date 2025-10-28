"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation"; // useRouter import

interface Survey {
  id: number;
  title: string;
}

export default function HomePage() {
  const { token, logout, isLoading } = useAuth();
  const [surveys, setSurveys] = useState<Survey[]>([]);
  const router = useRouter(); // useRouter 훅 사용

  useEffect(() => {
    async function getSurveys() {
      try {
        const res = await fetch("/api/surveys/list");
        if (!res.ok) throw new Error("Failed to fetch surveys");
        const data = await res.json();
        setSurveys(data);
      } catch (error) {
        console.error("Error fetching surveys:", error);
      }
    }
    getSurveys();
  }, []);

  const handleRegisterClick = () => {
    if (token) {
      // 로그인 상태이면 설문 생성 페이지로 이동
      router.push("/surveys/new");
    } else {
      // 로그아웃 상태이면 로그인 페이지로 이동
      alert("설문조사를 등록하려면 로그인이 필요합니다.");
      router.push("/login");
    }
  };

  return (
    <main style={{ fontFamily: "sans-serif", padding: "2rem" }}>
      <header
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
        }}
      >
        <h1 style={{ fontSize: "2rem" }}>📝 실시간 설문조사</h1>
        <div>
          <button onClick={handleRegisterClick} style={buttonStyle}>
            설문조사 등록
          </button>
          {/* 로그인 상태일 때만 로그아웃 버튼 표시 */}
          {!isLoading && token && (
            <button
              onClick={logout}
              style={{
                ...buttonStyle,
                marginLeft: "1rem",
                background: "gray",
              }}
            >
              로그아웃
            </button>
          )}
        </div>
      </header>

      <section style={{ marginTop: "2rem" }}>
        <h2
          style={{
            fontSize: "1.5rem",
            borderBottom: "2px solid #eee",
            paddingBottom: "0.5rem",
          }}
        >
          진행중인 설문조사
        </h2>
        <ul style={{ listStyle: "none", padding: 0 }}>
          {surveys.length > 0 ? (
            surveys.map((survey) => (
              <Link
                key={survey.id}
                href={`/surveys/${survey.id}`}
                style={{ textDecoration: "none", color: "inherit" }}
              >
                <li
                  style={{
                    padding: "1rem",
                    border: "1px solid #ddd",
                    borderRadius: "8px",
                    marginTop: "1rem",
                    cursor: "pointer",
                  }}
                >
                  {survey.title}
                </li>
              </Link>
            ))
          ) : (
            <p style={{ marginTop: "1rem", color: "#888" }}>
              생성된 설문조사가 없습니다. 첫 설문을 만들어보세요!
            </p>
          )}
        </ul>
      </section>
    </main>
  );
}

// 버튼 스타일 재사용
const buttonStyle: React.CSSProperties = {
  textDecoration: "none",
  background: "#0070f3",
  color: "white",
  padding: "0.75rem 1.5rem",
  borderRadius: "8px",
  border: "none",
  cursor: "pointer",
  fontSize: "1rem",
};
