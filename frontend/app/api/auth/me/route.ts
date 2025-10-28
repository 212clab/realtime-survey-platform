import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";

export async function GET(request: NextRequest) {
  // 서버 사이드에서는 request 객체나 next/headers를 통해 쿠키를 안전하게 읽을 수 있습니다.
  const cookieStore = await cookies();
  const tokenCookie = cookieStore.get("auth_token");

  if (!tokenCookie) {
    // 토큰이 없으면 로그인하지 않은 상태
    return NextResponse.json({ isLoggedIn: false }, { status: 401 });
  }

  // 토큰이 있으면 로그인한 상태 (실제로는 여기서 토큰 유효성 검증도 해야 함)
  return NextResponse.json({ isLoggedIn: true, token: tokenCookie.value });
}
