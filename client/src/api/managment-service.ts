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

export interface ClusterSummary {
  id: string;
  isRunning: boolean;
  precision: 'HIGH' | 'MEDIUM' | 'LOW';
  updatedAt: number;
  slackChannels: {
    name: string;
    webhookUrl: string;
    updatedAt: number;
  }[];
  discordChannels: {
    name: string;
    webhookUrl: string;
    updatedAt: number;
  }[];
  mailChannels: {
    name: string;
    email: string;
    updatedAt: number;
  }[];
}

const MANAGMENT_SERVICE_URL = import.meta.env.VITE_BACKEND_URL;
const VALID_URGENCY_VALUES: ReportSummary['urgency'][] = [
  'HIGH',
  'MEDIUM',
  'LOW',
];

class ManagmentServiceApi {
  private axiosInstance: AxiosInstance;

  private expiresIn: number | null;

  constructor() {
    this.axiosInstance = axios.create({
      baseURL: MANAGMENT_SERVICE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
      withCredentials: true,
    });

    this.expiresIn = null;
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

  public async isTokenExpired(): Promise<boolean> {
    if (this.expiresIn == null) {
      try {
        const tokenInfo = await this.getTokenInfoWithoutRefreshing();
        this.expiresIn = tokenInfo.expTime;
      } catch {
        this.expiresIn = 0;
      }
    }

    return this.expiresIn <= 0;
  }

  private async refreshTokenIfExpired() {
    if (await this.isTokenExpired()) {
      await this.refreshToken();
    }
  }

  public async getTokenInfoWithoutRefreshing(): Promise<TokenInfo> {
    const response = await this.axiosInstance.get(
      '/api/v1/auth/auth-token/validation-time',
    );
    return response.data;
  }

  public async getTokenInfo(): Promise<TokenInfo> {
    await this.refreshTokenIfExpired();
    const data = this.getTokenInfoWithoutRefreshing();
    return data;
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
        throw new Error(
          `Invalid urgency value "${report.urgency}" for report ID 
            ${report.id}. Allowed values are: ${VALID_URGENCY_VALUES.join(', ')}`,
        );
      }
    });

    return reports;
  }

  public async getClusters(): Promise<ClusterSummary[]> {
    const mockClusters: Array<ClusterSummary> = [
      {
        id: 'cluster-1-abcd',
        isRunning: true,
        precision: 'HIGH',
        updatedAt: 1730233614763,
        slackChannels: [
          {
            webhookUrl: 'mywebhookurl',
            name: 'wms-dev/infra',
            updatedAt: 1730233614763,
          },
        ],
        discordChannels: [
          {
            webhookUrl: 'mywebhookurl',
            name: 'wms-dev/infra',
            updatedAt: 1730233614763,
          },
        ],
        mailChannels: [
          {
            email: 'mail',
            name: 'wms-dev/infra',
            updatedAt: 1730233614763,
          },
        ],
      },
      {
        id: 'cluster-2-abcd',
        isRunning: false,
        precision: 'HIGH',
        updatedAt: 1730233614763,
        slackChannels: [
          {
            webhookUrl: 'mywebhookurl',
            name: 'wms-dev/infra',
            updatedAt: 1730233614763,
          },
        ],
        discordChannels: [
          {
            webhookUrl: 'mywebhookurl',
            name: 'wms-dev/infra',
            updatedAt: 1730233614763,
          },

          {
            webhookUrl: 'mywebhookurl',
            name: 'wms-dev/infra-2',
            updatedAt: 1730233614763,
          },
        ],
        mailChannels: [
          {
            email: 'mail',
            name: 'wms-dev/infra',
            updatedAt: 1730233614763,
          },
        ],
      },
    ];
    return mockClusters;
  }
}

export const ManagmentServiceApiInstance = new ManagmentServiceApi();
