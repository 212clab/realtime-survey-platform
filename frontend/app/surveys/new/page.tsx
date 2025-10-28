"use client"; // 클라이언트 컴포넌트로 전환 (useState, form 이벤트 사용)

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";

export default function NewSurveyPage() {
  const [title, setTitle] = useState("");
  const [options, setOptions] = useState(["", ""]); // 기본 선택지 2개
  const router = useRouter();

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

    const API_URL = "/api/surveys"; // 우리 Next.js 서버의 API Route

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
