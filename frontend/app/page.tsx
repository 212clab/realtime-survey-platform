import Link from "next/link";

interface Survey {
  id: number;
  title: string;
}

// 서버 컴포넌트에서 직접 백엔드 API를 호출하는 함수
async function getSurveys(): Promise<Survey[]> {
  try {
    // 서버 컴포넌트는 Docker 내부 네트워크에 있으므로, 서비스 이름으로 직접 호출
    const res = await fetch("http://survey-service:8080/surveys", {
      cache: "no-store", // 항상 최신 데이터를 가져오기 위해 캐시 비활성화
    });

    if (!res.ok) {
      throw new Error("Failed to fetch surveys");
    }
    return res.json();
  } catch (error) {
    console.error("Error fetching surveys:", error);
    return []; // 에러 발생 시 빈 배열 반환
  }
}

// 페이지 컴포넌트를 async 함수로 변경
export default async function HomePage() {
  // 페이지 렌더링 전에 서버에서 데이터를 미리 가져옵니다.
  const surveys = await getSurveys();

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
