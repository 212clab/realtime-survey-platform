import { NextResponse } from "next/server";
import { cookies } from "next/headers";

export async function GET() {
  // 쿠키를 삭제합니다.
  const cookieStore = await cookies();
  cookieStore.delete("auth_token");
  return NextResponse.json({ message: "Logged out successfully" });
}
