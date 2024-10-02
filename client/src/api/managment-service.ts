import axios from 'axios';

const MANAGMENT_SERVICE_URL = import.meta.env.VITE_BACKEND_URL;

export const managmentServiceAxios = axios.create({
  baseURL: MANAGMENT_SERVICE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

export const login = async () => {
  // const response = await managmentServiceAxios.get('/auth/login', { code: authCode });
  window.location.href = `${MANAGMENT_SERVICE_URL}/`;
};

export const getTokenInfo = async () => {
  const response = await managmentServiceAxios.get('/auth/profile');
  return response.data;
};

export const logout = async () => {
  document.cookie = 'id_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
  document.location.href = '/';
};
