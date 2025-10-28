import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const res = await fetch("http://user-service:8080/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      return NextResponse.json(
        { message: "Login failed" },
        { status: res.status }
      );
    }

    const data = await res.json();

    // 👇 여기서부터 수정: JSON을 반환하는 대신 쿠키를 설정합니다.
    const response = NextResponse.json({ success: true });
    response.cookies.set("auth_token", data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV !== "development",
      path: "/",
      maxAge: 60 * 60 * 24, // 1 day
    });
    return response;
    // 👆 여기까지 수정
  } catch (error) {
    return NextResponse.json(
      { message: "Internal Server Error" },
      { status: 500 }
    );
  }
}
