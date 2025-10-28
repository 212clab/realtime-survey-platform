import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const code = request.nextUrl.searchParams.get("code");

  if (!code) {
    return NextResponse.json({ message: "Code not found" }, { status: 400 });
  }

  try {
    // 1. 백엔드(user-service)에 허가증(code)을 전달하고 최종 JWT를 요청합니다.
    const res = await fetch(
      `http://user-service:8080/auth/github/callback?code=${code}`
    );

    if (!res.ok) {
      throw new Error("Failed to login with github");
    }

    const data = await res.json(); // 백엔드가 JWT를 반환한다고 가정

    // 2. JWT를 쿠키에 저장하고 사용자를 홈으로 리디렉션합니다.
    const response = NextResponse.redirect("http://localhost:3000/");
    response.cookies.set("auth_token", data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV !== "development",
      path: "/",
      maxAge: 60 * 60 * 24, // 1 day
    });

    return response;
  } catch (error) {
    return NextResponse.redirect("http://localhost:3000/login?error=true");
  }
}
