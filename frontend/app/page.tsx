"use client"; // 👈 클라이언트 컴포넌트로 전환

import Link from "next/link";
import { useEffect, useState } from "react"; // 👈 useEffect, useState import
import { useAuth } from "@/contexts/AuthContext"; // 👈 useAuth import

interface Survey {
  id: number;
  title: string;
}

export default function HomePage() {
  const { token, logout, isLoading } = useAuth(); // 👈 AuthContext에서 정보 가져오기
  const [surveys, setSurveys] = useState<Survey[]>([]);

  // 클라이언트 사이드에서 API를 호출하도록 변경
  useEffect(() => {
    async function getSurveys() {
      try {
        // Next.js 서버의 대리인(Route Handler)을 통해 API 호출
        const res = await fetch("/api/surveys/list"); // 새로운 Route Handler 경로
        if (!res.ok) throw new Error("Failed to fetch surveys");
        const data = await res.json();
        setSurveys(data);
      } catch (error) {
        console.error("Error fetching surveys:", error);
      }
    }
    getSurveys();
  }, []);

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
          {!isLoading && // 로딩이 끝난 후에 버튼을 보여줌
            (token ? (
              <>
                <Link href="/surveys/new" style={buttonStyle}>
                  + 새 설문 만들기
                </Link>
                <button
                  onClick={logout}
                  style={{
                    ...buttonStyle,
                    marginLeft: "1rem",
                    background: "#dc3545",
                  }}
                >
                  로그아웃
                </button>
              </>
            ) : (
              <Link href="/login" style={buttonStyle}>
                로그인 / 회원가입
              </Link>
            ))}
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

// 버튼 스타일 재사용을 위해 변수로 추출
const buttonStyle: React.CSSProperties = {
  textDecoration: "none",
  background: "#0070f3",
  color: "white",
  padding: "0.75rem 1.5rem",
  borderRadius: "8px",
  border: "none",
  cursor: "pointer",
};
