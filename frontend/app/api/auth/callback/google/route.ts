import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const code = request.nextUrl.searchParams.get("code");

  if (!code) {
    return NextResponse.redirect(
      "http://localhost:3000/login?error=code-not-found"
    );
  }

  try {
    // 백엔드(user-service)에 허가증(code)을 전달하고 최종 JWT를 요청합니다.
    const res = await fetch(
      `http://user-service:8080/auth/google/callback?code=${code}`
    );

    if (!res.ok) {
      const errorBody = await res.text();
      console.error("Backend error:", errorBody);
      throw new Error("Failed to login with google");
    }

    const data = await res.json(); // 백엔드가 JWT를 반환합니다.

    // JWT를 쿠키에 저장하고 사용자를 홈으로 리디렉션합니다.
    const response = NextResponse.redirect("http://localhost:3000/");
    response.cookies.set("auth_token", data.token, {
      httpOnly: true, // 브라우저 JS에서 접근 불가
      secure: process.env.NODE_ENV !== "development", // https에서만 전송
      path: "/",
      maxAge: 60 * 60 * 24, // 1일
    });

    return response;
  } catch (error) {
    console.error("Callback handler error:", error);
    return NextResponse.redirect("http://localhost:3000/login?error=true");
  }
}
