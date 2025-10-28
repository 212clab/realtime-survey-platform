"use client";

import React, {
  createContext,
  useState,
  useContext,
  useEffect,
  ReactNode,
} from "react";
import { useRouter } from "next/navigation";

interface AuthContextType {
  token: string | null;
  isLoading: boolean;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true); // 항상 true로 시작
  const router = useRouter();

  useEffect(() => {
    // 컴포넌트가 마운트될 때 서버에 로그인 상태를 확인합니다.
    const checkLoginStatus = async () => {
      try {
        const res = await fetch("/api/auth/me");
        if (res.ok) {
          const data = await res.json();
          setToken(data.token);
        } else {
          setToken(null);
        }
      } catch (error) {
        setToken(null);
        console.error("Failed to check login status", error);
      } finally {
        setIsLoading(false);
      }
    };

    checkLoginStatus();
  }, []);

  const logout = async () => {
    // 로그아웃도 API를 통해 서버에 요청해서 쿠키를 삭제해야 합니다.
    await fetch("/api/auth/logout");
    setToken(null);
    router.push("/login");
  };

  return (
    <AuthContext.Provider value={{ token, isLoading, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
