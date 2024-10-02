import axios, { AxiosInstance } from 'axios';

interface UserInfo {
  nickname: string;
  email: string;
}

interface TokenInfo {
  expTime: number;
}

const MANAGMENT_SERVICE_URL = import.meta.env.VITE_BACKEND_URL;

class ManagmentServiceApi {
  private axiosInstance: AxiosInstance;

  private expiresIn: number;

  constructor() {
    this.axiosInstance = axios.create({
      baseURL: MANAGMENT_SERVICE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
      withCredentials: true,
    });
    this.expiresIn = 0;
  }

  public async login() {
    window.location.href = `${MANAGMENT_SERVICE_URL}/oauth2/authorization/google`;
  }

  public async logout() {
    document.cookie = 'authToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    document.cookie = 'refreshToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    document.location.href = '/';
  }

  private async refreshToken() {
    await this.axiosInstance.get('/api/v1/auth/refresh-token');
  }

  public isTokenExpired(): boolean {
    return this.expiresIn <= 0;
  }

  private async refreshTokenIfExpired() {
    if (this.isTokenExpired()) {
      await this.refreshToken();
    }
  }

  public async getTokenInfo(): Promise<TokenInfo> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get('/api/v1/auth/auth-token/validation-time');
    return response.data;
  }

  public async getUserInfo(): Promise<UserInfo> {
    await this.refreshTokenIfExpired();
    const user = await this.axiosInstance.get('/api/v1/auth/user-details');
    return user.data;
  }
}

export const ManagmentServiceApiInstance = new ManagmentServiceApi();
