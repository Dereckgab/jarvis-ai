"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { apiClient } from "@/lib/api-client";

interface User {
  id: string;
  email: string;
  name: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<boolean>;
  register: (email: string, password: string, name: string) => Promise<boolean>;
  logout: () => void;
  refreshAccessToken: () => Promise<boolean>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [refreshToken, setRefreshToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Initialize auth from localStorage
  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    const storedRefreshToken = localStorage.getItem("refreshToken");
    const storedUser = localStorage.getItem("user");

    if (storedToken && storedUser) {
      setToken(storedToken);
      setRefreshToken(storedRefreshToken);
      setUser(JSON.parse(storedUser));
    }

    setIsLoading(false);
  }, []);

  const login = async (email: string, password: string): Promise<boolean> => {
    try {
      const response = await apiClient.login(email, password);

      if (!response.success || !response.data) {
        return false;
      }

      const { access_token: newToken, refresh_token: newRefreshToken } = response.data;

      // Store tokens
      localStorage.setItem("token", newToken);
      localStorage.setItem("refreshToken", newRefreshToken);

      // Decode JWT to get user info (simplified - in production use jwt-decode library)
      const decodedUser: User = {
        id: "user_id",
        email,
        name: email.split("@")[0],
      };

      localStorage.setItem("user", JSON.stringify(decodedUser));

      setToken(newToken);
      setRefreshToken(newRefreshToken);
      setUser(decodedUser);

      return true;
    } catch (error) {
      console.error("Login error:", error);
      return false;
    }
  };

  const register = async (email: string, password: string, name: string): Promise<boolean> => {
    try {
      const response = await apiClient.register(email, password, name);

      if (!response.success || !response.data) {
        return false;
      }

      const { access_token: newToken, refresh_token: newRefreshToken } = response.data;

      // Store tokens
      localStorage.setItem("token", newToken);
      localStorage.setItem("refreshToken", newRefreshToken);

      const newUser: User = {
        id: "user_id",
        email,
        name,
      };

      localStorage.setItem("user", JSON.stringify(newUser));

      setToken(newToken);
      setRefreshToken(newRefreshToken);
      setUser(newUser);

      return true;
    } catch (error) {
      console.error("Registration error:", error);
      return false;
    }
  };

  const logout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("refreshToken");
    localStorage.removeItem("user");

    setToken(null);
    setRefreshToken(null);
    setUser(null);
  };

  const refreshAccessToken = async (): Promise<boolean> => {
    if (!refreshToken) {
      logout();
      return false;
    }

    try {
      const response = await apiClient.refreshToken(refreshToken);

      if (!response.success || !response.data) {
        logout();
        return false;
      }

      const { access_token: newToken, refresh_token: newRefreshToken } = response.data;
      localStorage.setItem("token", newToken);
      localStorage.setItem("refreshToken", newRefreshToken ?? refreshToken);
      setToken(newToken);
      setRefreshToken(newRefreshToken ?? refreshToken);

      return true;
    } catch (error) {
      console.error("Token refresh error:", error);
      logout();
      return false;
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        refreshToken,
        isLoading,
        isAuthenticated: !!token,
        login,
        register,
        logout,
        refreshAccessToken,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);

  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
}
