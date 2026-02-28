import api from "./http";
import { LoginResponse } from "../types/auth";

// LOGIN
export async function login(loginId: string, password: string): Promise<LoginResponse> {
  const res = await api.post("/auth/login", { login: loginId, password });
  return res.data;
}

// LOGOUT
export async function logout(): Promise<void> {
  await api.post("/auth/logout", {});
}

// REFRESH TOKEN
export async function refreshAccessToken(): Promise<any> {
  const res = await api.post("/auth/refresh", {});
  return res.data;
  // backend returns: { access_token, expires_in, token_type }
}