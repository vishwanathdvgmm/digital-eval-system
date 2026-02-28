import { useEffect, useState, useRef } from "react";
import { AuthContext } from "../context/AuthContext";
import { login as apiLogin, logout as apiLogout, refreshAccessToken } from "../api/auth";
import { saveAccessToken, removeAccessToken, parseJwt, getAccessToken } from "../utils/token";
import { User } from "../types/auth";

interface Props {
  children: React.ReactNode;
}

export function AuthProvider({ children }: Props) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // LOGIN
  async function handleLogin(loginId: string, password: string) {
    const data = await apiLogin(loginId, password);
    saveAccessToken(data.access_token);
    setUser(data.user);
    return data.user;
  }

  // LOGOUT
  async function handleLogout() {
    await apiLogout();
    removeAccessToken();
    setUser(null);
  }

  // AUTO REFRESH ON PAGE LOAD
  async function attemptRefresh() {
    try {
      // Optimization: Check if we already have a valid token
      const token = getAccessToken();
      if (token) {
        const decoded = parseJwt(token);
        if (decoded) {
          setUser({
            id: decoded.uid,
            user_id: decoded.user_id,
            name: decoded.name,
            role: decoded.role,
            email: decoded.email,
            created_at: "", // Not available in token, placeholder
            updated_at: ""  // Not available in token, placeholder
          });
          setLoading(false);
          return;
        }
      }

      const data = await refreshAccessToken();
      saveAccessToken(data.access_token);

      // Decode token to get user details
      const decoded = parseJwt(data.access_token);
      if (decoded) {
        setUser({
          id: decoded.uid,
          user_id: decoded.user_id,
          name: decoded.name,
          role: decoded.role,
          email: decoded.email,
          created_at: "", // Not available in token, placeholder
          updated_at: ""  // Not available in token, placeholder
        });
      } else {
        setUser((prev) => prev ?? null);
      }
    } catch {
      removeAccessToken();
      setUser(null);
    } finally {
      setLoading(false);
    }
  }

  const initialized = useRef(false);

  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true;
      attemptRefresh();
    }
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login: handleLogin,
        logout: handleLogout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}