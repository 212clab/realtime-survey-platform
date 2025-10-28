import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const code = request.nextUrl.searchParams.get("code");

  if (!code) {
    return NextResponse.redirect(
      "http://localhost:3000/login?error=code-not-found"
    );
  }

  try {
    const res = await fetch(
      `http://user-service:8080/auth/github/callback?code=${code}`
    );
    if (!res.ok) {
      throw new Error("Failed to login with github on backend");
    }

    const data = await res.json();
    const token = data.token;

    // 1. 리디렉션할 URL을 준비합니다.
    const url = request.nextUrl.clone();
    url.pathname = "/"; // 홈페이지로 이동

    // 2. 리디렉션 응답을 먼저 생성합니다.
    const response = NextResponse.redirect(url);

    // 3. 생성된 응답에 쿠키를 설정합니다.
    response.cookies.set({
      name: "auth_token",
      value: token,
      httpOnly: true,
      secure: process.env.NODE_ENV !== "development",
      path: "/",
      maxAge: 60 * 60 * 24, // 1 day
    });

    return response;
  } catch (error) {
    console.error("Callback handler error:", error);
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    url.searchParams.set("error", "true");
    return NextResponse.redirect(url);
  }
}
