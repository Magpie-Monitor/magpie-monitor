import axios, { AxiosInstance } from 'axios';

interface UserInfo {
  nickname: string;
  email: string;
}

interface TokenInfo {
  expTime: number;
}

export interface ReportSummary {
  id: string;
  clusterId: string;
  title: string;
  urgency: 'HIGH' | 'MEDIUM' | 'LOW';
  sinceMs: number;
  toMs: number;
  [key: string]: string | number;
}

const MANAGMENT_SERVICE_URL = import.meta.env.VITE_BACKEND_URL;
const VALID_URGENCY_VALUES: ReportSummary['urgency'][] = ['HIGH', 'MEDIUM', 'LOW'];


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
    await this.axiosInstance.get('/api/v1/auth/logout');
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

  public async getReports(): Promise<ReportSummary[]> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get('/api/v1/reports');
    const reports: ReportSummary[] = response.data;
    reports.forEach((report) => {
      if (!VALID_URGENCY_VALUES.includes(report.urgency)) {
        throw new Error(`Invalid urgency value "${report.urgency}" for report ID ${report.id}. Allowed values are: ${VALID_URGENCY_VALUES.join(', ')}`);
      }
    });

    return reports;
  }
}

export const ManagmentServiceApiInstance = new ManagmentServiceApi();
