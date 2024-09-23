import axios from "axios";

export const proxyAxios = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

export const login = async (authCode: string) => {
  const response = await proxyAxios.post("/auth/login", { code: authCode });
  return response.data;
};

export const getTokenInfo = async () => {
  const response = await proxyAxios.get("/auth/profile");
  return response.data;
};

export const logout = async () => {
  document.cookie = "id_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
  document.location.href = "/";
};
