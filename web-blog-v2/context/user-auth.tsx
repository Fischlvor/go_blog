'use client';

import { createContext, useContext, useEffect, useState, useCallback } from 'react';
import type { UserProfile } from '@/lib/client-api/types';
import { getUserInfo, logout as apiLogout } from '@/lib/client-api/user/user';

interface UserAuthContextValue {
  user: UserProfile | null;
  isLoggedIn: boolean;
  isLoading: boolean;
  token: string | null;
  setToken: (token: string) => void;
  refreshUser: () => Promise<void>;
  logout: () => Promise<void>;
}

const UserAuthContext = createContext<UserAuthContextValue | null>(null);

export function UserAuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [token, setTokenState] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const setToken = useCallback((newToken: string) => {
    localStorage.setItem('access_token', newToken);
    setTokenState(newToken);
  }, []);

  const refreshUser = useCallback(async () => {
    const storedToken = localStorage.getItem('access_token');
    if (!storedToken) {
      setUser(null);
      setTokenState(null);
      setIsLoading(false);
      return;
    }
    setTokenState(storedToken);
    try {
      const userInfo = await getUserInfo();
      setUser(userInfo);
    } catch {
      // token 无效，清除
      localStorage.removeItem('access_token');
      setUser(null);
      setTokenState(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      await apiLogout();
    } catch {
      // ignore
    } finally {
      localStorage.removeItem('access_token');
      setUser(null);
      setTokenState(null);
    }
  }, []);

  useEffect(() => {
    refreshUser();
  }, [refreshUser]);

  return (
    <UserAuthContext.Provider
      value={{
        user,
        isLoggedIn: !!user,
        isLoading,
        token,
        setToken,
        refreshUser,
        logout,
      }}
    >
      {children}
    </UserAuthContext.Provider>
  );
}

export function useUserAuth(): UserAuthContextValue {
  const ctx = useContext(UserAuthContext);
  if (!ctx) throw new Error('useUserAuth must be used within UserAuthProvider');
  return ctx;
}
