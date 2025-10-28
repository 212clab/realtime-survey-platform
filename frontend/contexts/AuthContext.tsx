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
  login: (token: string) => void;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true); // 페이지 로드 시 토큰 확인 중인지 여부
  const router = useRouter();

  useEffect(() => {
    // 앱이 처음 로드될 때 localStorage에서 토큰을 확인
    const storedToken = localStorage.getItem("jwt_token");
    if (storedToken) {
      setToken(storedToken);
    }
    setIsLoading(false);
  }, []);

  const login = (newToken: string) => {
    setToken(newToken);
    localStorage.setItem("jwt_token", newToken); // 토큰을 브라우저 저장소에 보관
  };

  const logout = () => {
    setToken(null);
    localStorage.removeItem("jwt_token"); // 저장소에서 토큰 삭제
    router.push("/login"); // 로그아웃 후 로그인 페이지로 이동
  };

  return (
    <AuthContext.Provider value={{ token, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
};

// 다른 컴포넌트에서 쉽게 AuthContext를 사용할 수 있도록 하는 커스텀 훅
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
