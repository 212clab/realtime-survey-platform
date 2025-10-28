import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/contexts/AuthContext"; // 👈 AuthProvider import

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Realtime Survey",
  description: "A survey platform built with Next.js and Go",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <AuthProvider>
          {" "}
          {/* 👈 앱 전체를 감싸줍니다 */}
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}
