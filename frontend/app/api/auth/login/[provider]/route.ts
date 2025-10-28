import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  // 👇 request 객체의 URL에서 직접 경로를 분석합니다.
  const url = new URL(request.url);
  const pathSegments = url.pathname.split("/"); // 예: ['', 'api', 'auth', 'login', 'github']
  const provider = pathSegments[pathSegments.length - 1]; // 배열의 마지막 요소('github')를 가져옵니다.

  if (!provider) {
    return NextResponse.json(
      { message: "Provider parameter is missing" },
      { status: 400 }
    );
  }

  if (provider === "github") {
    const GITHUB_CLIENT_ID = process.env.GITHUB_CLIENT_ID;
    if (!GITHUB_CLIENT_ID) {
      console.error("GITHUB_CLIENT_ID is not set in the environment.");
      return NextResponse.json(
        { message: "Server configuration error" },
        { status: 500 }
      );
    }
    const GITHUB_SCOPE = "read:user user:email";
    return NextResponse.redirect(
      `https://github.com/login/oauth/authorize?client_id=${GITHUB_CLIENT_ID}&scope=${GITHUB_SCOPE}`
    );
  }

  if (provider === "google") {
    const GOOGLE_CLIENT_ID = process.env.GOOGLE_CLIENT_ID;
    if (!GOOGLE_CLIENT_ID) {
      console.error("GOOGLE_CLIENT_ID is not set in the environment.");
      return NextResponse.json(
        { message: "Server configuration error" },
        { status: 500 }
      );
    }
    const GOOGLE_REDIRECT_URI =
      "http://localhost:3000/api/auth/callback/google";
    const GOOGLE_SCOPE =
      "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile";
    return NextResponse.redirect(
      `https://accounts.google.com/o/oauth2/v2/auth?client_id=${GOOGLE_CLIENT_ID}&redirect_uri=${GOOGLE_REDIRECT_URI}&scope=${GOOGLE_SCOPE}&response_type=code`
    );
  }

  return NextResponse.json(
    { message: `Invalid provider: ${provider}` },
    { status: 400 }
  );
}
