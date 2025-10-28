import { NextResponse } from "next/server";

export async function GET() {
  try {
    const res = await fetch("http://survey-service:8080/surveys", {
      cache: "no-store",
    });

    if (!res.ok) {
      throw new Error("Backend API call failed");
    }

    const data = await res.json();
    return NextResponse.json(data);
  } catch (error) {
    return NextResponse.json(
      { message: "Internal Server Error" },
      { status: 500 }
    );
  }
}
