import axios from "axios";

const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID;
const GOOGLE_AUTH_URL = "https://accounts.google.com/o/oauth2";

export const googleAxios = axios.create({
  baseURL: GOOGLE_AUTH_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

export const navigateToGoogleAuth = () => {
  const redirectUri = window.location.origin;
  const url =
    `${GOOGLE_AUTH_URL}/auth?` +
    `client_id=${GOOGLE_CLIENT_ID}` +
    `&redirect_uri=${redirectUri}` +
    `&response_type=code` +
    `&scope=email` +
    `&access_type=offline` +
    `& approval_prompt=force`;
  window.location.href = url;
};
