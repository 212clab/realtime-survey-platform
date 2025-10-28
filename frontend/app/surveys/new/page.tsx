"use client"; // 클라이언트 컴포넌트로 전환 (useState, form 이벤트 사용)

import Link from "next/link";
import { useState, useEffect } from "react"; // useEffect 추가
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext"; // useAuth 훅 import

export default function NewSurveyPage() {
  const { token, isLoading } = useAuth(); // AuthContext에서 토큰과 로딩 상태 가져오기
  const router = useRouter();

  useEffect(() => {
    // 로딩이 끝났는데 토큰이 없으면 로그인 페이지로 쫓아냄
    if (!isLoading && !token) {
      router.push("/login");
    }
  }, [token, isLoading, router]);

  const [title, setTitle] = useState("");
  const [options, setOptions] = useState(["", ""]); // 기본 선택지 2개

  const handleOptionChange = (index: number, value: string) => {
    const newOptions = [...options];
    newOptions[index] = value;
    setOptions(newOptions);
  };

  const addOption = () => {
    setOptions([...options, ""]);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) {
      // 토큰이 없으면 제출 자체를 막음
      alert("로그인이 필요합니다.");
      return;
    }
    const API_URL = "/api/surveys"; // 우리 Next.js 서버의 API Route
    // 로딩 중이거나 토큰이 없으면 페이지 내용을 보여주지 않음
    if (isLoading || !token) {
      return <p>Loading...</p>;
    }
    const surveyData = {
      title,
      options: options
        .filter((opt) => opt.trim() !== "")
        .map((opt) => ({ text: opt })),
    };

    try {
      const res = await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(surveyData),
      });

      if (res.ok) {
        alert("설문이 성공적으로 생성되었습니다!");
        router.push("/"); // 성공 시 홈으로 이동
      } else {
        alert("설문 생성에 실패했습니다.");
      }
    } catch (error) {
      console.error("API 호출 오류:", error);
      alert("서버와 통신할 수 없습니다.");
    }
  };

  return (
    <main
      style={{
        fontFamily: "sans-serif",
        padding: "2rem",
        maxWidth: "600px",
        margin: "auto",
      }}
    >
      {" "}
      <Link href="/" style={{ textDecoration: "none", color: "#0070f3" }}>
        &larr; 목록으로 돌아가기
      </Link>
      <h1 style={{ fontSize: "2rem", marginTop: "1rem" }}>새 설문 만들기</h1>
      <form
        onSubmit={handleSubmit}
        style={{
          marginTop: "2rem",
          display: "flex",
          flexDirection: "column",
          gap: "1rem",
        }}
      >
        <div>
          <label
            htmlFor="title"
            style={{ display: "block", marginBottom: "0.5rem" }}
          >
            설문 제목
          </label>
          <input
            type="text"
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            style={{
              width: "100%",
              padding: "0.75rem",
              borderRadius: "8px",
              border: "1px solid #ddd",
            }}
            required
          />
        </div>

        <div>
          <label style={{ display: "block", marginBottom: "0.5rem" }}>
            선택지
          </label>
          {options.map((option, index) => (
            <input
              key={index}
              type="text"
              value={option}
              onChange={(e) => handleOptionChange(index, e.target.value)}
              placeholder={`선택지 ${index + 1}`}
              style={{
                width: "100%",
                padding: "0.75rem",
                borderRadius: "8px",
                border: "1px solid #ddd",
                marginBottom: "0.5rem",
              }}
            />
          ))}
          <button
            type="button"
            onClick={addOption}
            style={
              {
                /* ... 스타일링 ... */
              }
            }
          >
            선택지 추가
          </button>
        </div>

        <button
          type="submit"
          style={{
            cursor: "pointer",
            background: "#0070f3",
            color: "white",
            padding: "0.75rem 1.5rem",
            borderRadius: "8px",
            border: "none",
            fontSize: "1rem",
          }}
        >
          설문 생성하기
        </button>
      </form>
    </main>
  );
}
