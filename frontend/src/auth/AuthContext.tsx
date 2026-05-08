import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';
import type { UserResponse } from '../types';
import { getToken, setToken } from '../api/client';
import { getUser } from '../api/users';
import { decodeJwt, isExpired } from './jwt';

interface AuthContextValue {
  token: string | null;
  user: UserResponse | null;
  loading: boolean;
  signIn(token: string): Promise<void>;
  signOut(): void;
  refreshUser(): Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string | null>(() => {
    const stored = getToken();
    if (!stored) return null;
    if (isExpired(decodeJwt(stored))) {
      setToken(null);
      return null;
    }
    return stored;
  });
  const [user, setUser] = useState<UserResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(!!token);

  const loadUser = useCallback(async (rawToken: string) => {
    const payload = decodeJwt(rawToken);
    if (!payload || isExpired(payload)) {
      setToken(null);
      setTokenState(null);
      setUser(null);
      return;
    }
    try {
      const fetched = await getUser(payload.UserID);
      setUser(fetched);
    } catch (err) {
      // If the token is rejected the user is signed out by the API client.
      console.error('failed to load user', err);
      setToken(null);
      setTokenState(null);
      setUser(null);
    }
  }, []);

  useEffect(() => {
    if (!token) {
      setLoading(false);
      return;
    }
    setLoading(true);
    loadUser(token).finally(() => setLoading(false));
  }, [token, loadUser]);

  const signIn = useCallback(
    async (newToken: string) => {
      setToken(newToken);
      setTokenState(newToken);
      await loadUser(newToken);
    },
    [loadUser],
  );

  const signOut = useCallback(() => {
    setToken(null);
    setTokenState(null);
    setUser(null);
  }, []);

  const refreshUser = useCallback(async () => {
    if (!token) return;
    await loadUser(token);
  }, [token, loadUser]);

  const value = useMemo<AuthContextValue>(
    () => ({ token, user, loading, signIn, signOut, refreshUser }),
    [token, user, loading, signIn, signOut, refreshUser],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
