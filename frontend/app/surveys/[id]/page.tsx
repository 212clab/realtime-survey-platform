import Link from "next/link";

// 임시 데이터. 실제로는 이 ID를 사용해 API로 데이터를 가져옵니다.
const dummySurveyDetail = {
  id: 1,
  title: "가장 좋아하는 개발 언어는?",
  options: [
    { id: 101, text: "JavaScript" },
    { id: 102, text: "Python" },
    { id: 103, text: "Go" },
    { id: 104, text: "Rust" },
  ],
};

// params 객체를 통해 URL의 동적 세그먼트([id]) 값을 받을 수 있습니다.
export default function SurveyDetailPage({
  params,
}: {
  params: { id: string };
}) {
  const { id } = params; // URL에서 설문조사 ID를 추출합니다.

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

      {/* URL에서 받은 ID를 화면에 표시 */}
      <p style={{ color: "#888", marginTop: "1rem" }}>Survey ID: {id}</p>

      <h1 style={{ fontSize: "2rem", marginTop: "0.5rem" }}>
        {dummySurveyDetail.title}
      </h1>

      <div
        style={{
          marginTop: "2rem",
          display: "flex",
          flexDirection: "column",
          gap: "1rem",
        }}
      >
        {dummySurveyDetail.options.map((option) => (
          <button
            key={option.id}
            style={{
              padding: "1rem",
              border: "1px solid #ddd",
              borderRadius: "8px",
              textAlign: "left",
              fontSize: "1rem",
              cursor: "pointer",
              background: "white",
            }}
          >
            {option.text}
          </button>
        ))}
      </div>
    </main>
  );
}
