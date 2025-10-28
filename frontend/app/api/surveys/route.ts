import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    // 1. 브라우저로부터 받은 요청 본문을 가져옵니다.
    const body = await request.json();

    // 2. 백엔드 서비스(Go)로 요청을 전달합니다. (내부 통신)
    const res = await fetch("http://survey-service:8080/surveys", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });

    // 3. 백엔드의 응답이 정상이 아니면 에러를 반환합니다.
    if (!res.ok) {
      throw new Error(`API call failed with status: ${res.status}`);
    }

    // 4. 백엔드로부터 받은 응답을 그대로 브라우저에 전달합니다.
    const data = await res.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("API Route Error:", error);
    return NextResponse.json(
      { message: "Internal Server Error" },
      { status: 500 }
    );
  }
}
