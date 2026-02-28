import axios from "axios";
import { refreshAccessToken } from "./auth";
import { getAccessToken, saveAccessToken, removeAccessToken } from "../utils/token";

const api = axios.create({
  baseURL: "/api/v1/",
  withCredentials: true, // required for refresh cookie
});

// Attach access token
api.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auto refresh on 401
api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config;

    if (error.response?.status === 401 && !original._retry) {
      // Prevent infinite loop: don't retry if the failed request was already a refresh attempt
      if (original.url?.includes("/auth/refresh")) {
        return Promise.reject(error);
      }

      original._retry = true;

      try {
        const data = await refreshAccessToken();
        saveAccessToken(data.access_token);

        original.headers.Authorization = `Bearer ${data.access_token}`;
        return api(original);
      } catch (err) {
        removeAccessToken();
        window.location.href = "/login";
        return Promise.reject(err);
      }
    }

    return Promise.reject(error);
  }
);

export default api;