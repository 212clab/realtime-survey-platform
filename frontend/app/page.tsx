import Link from "next/link";

// Survey 데이터의 타입을 정의합니다.
interface Survey {
  id: number;
  title: string;
}

// 임시 데이터 배열에 Survey[] 타입을 적용합니다.
const dummySurveys: Survey[] = [
  { id: 1, title: "가장 좋아하는 개발 언어는?" },
  { id: 2, title: "점심 메뉴로 뭐가 좋을까요?" },
  { id: 3, title: "DevOps 공부, 가장 어려운 점은?" },
];

export default function HomePage() {
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
        <Link
          href="/surveys/new"
          style={{
            textDecoration: "none",
            background: "#0070f3",
            color: "white",
            padding: "0.75rem 1.5rem",
            borderRadius: "8px",
          }}
        >
          + 새 설문 만들기
        </Link>
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
          {dummySurveys.map((survey) => (
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
          ))}
        </ul>
      </section>
    </main>
  );
}
