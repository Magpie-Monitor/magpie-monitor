import axios, { AxiosInstance } from 'axios';

interface UserInfo {
  nickname: string;
  email: string;
}

interface TokenInfo {
  expTime: number;
}

export type AccuracyLevel = 'HIGH' | 'MEDIUM' | 'LOW';
export type UrgencyLevel = 'HIGH' | 'MEDIUM' | 'LOW';

export interface ReportSummary {
  id: string;
  clusterId: string;
  title: string;
  urgency: UrgencyLevel;
  sinceMs: number;
  toMs: number;
  [key: string]: string | number;
}

export interface ClusterSummary {
  id: string;
  isRunning: boolean;
  accuracy: AccuracyLevel;
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

export interface NotificationChannel {
  id: string;
  name: string;
  service: string;
  details: string;
  updated: string;
  added: string;
}

export interface Application {
  name: string;
  running: boolean;
  accuracy: AccuracyLevel;
  customPrompt: string;
  updated: string;
  added: string;
}

export interface Node {
  name: string;
  running: boolean;
  accuracy: AccuracyLevel;
  customPrompt: string;
  updated: string;
  added: string;
}

const MANAGMENT_SERVICE_URL = import.meta.env.VITE_BACKEND_URL;
const VALID_URGENCY_VALUES: ReportSummary['urgency'][] = [
  'HIGH',
  'MEDIUM',
  'LOW',
];

export interface ApplicationIncident {
  clusterId: string;
  applicationName: string;
  summary: string;
  recommendation: string;
  sources: ApplicationIncidentSource[];
}

export interface ApplicationIncidentSource {
  pod: string;
  container: string;
  image: string;
  content: string;
  timestamp: number;
}

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

  public async getApplicationIncident(
    _id: string,
  ): Promise<ApplicationIncident> {
    return {
      clusterId: 'Cluster 1',
      applicationName: 'application-1',
      summary: 'This is an summar of the incident',
      recommendation: 'This is an recommendation of the incident',
      sources: [
        {
          container: 'container-1',
          timestamp: 213213124,
          pod: 'pod-1',
          image: 'image-1',
          content:
            'LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
            LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
            LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
             LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS',
        },
        {
          container: 'container-1',
          timestamp: 213213124,
          pod: 'pod-1',
          image: 'image-1',
          content:
            'LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
            LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
            LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
             LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS',
        },

        {
          container: 'container-1',
          timestamp: 213213124,
          pod: 'pod-1',
          image: 'image-1',
          content:
            'LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
            LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
          LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
            LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS\
             LOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOSLOGSLOGSLOGSLGOSLGOSGLOSGLOSGLOSGLOS',
        },
      ],
    };
  }

  public async getClusters(): Promise<ClusterSummary[]> {
    const mockClusters: Array<ClusterSummary> = [
      {
        id: 'cluster-1-abcd',
        isRunning: true,
        accuracy: 'LOW',
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
        accuracy: 'HIGH',
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

  public async getNotificationChannels(): Promise<NotificationChannel[]> {
    const mockNotificatoinChannels: Array<NotificationChannel> = [
      {
        id: '1',
        name: 'Infra team slack',
        service: 'SLACK',
        details: 'wms_dev/#infra-alerts',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
      },
      {
        id: '2',
        name: 'Infra team discord',
        service: 'DISCORD',
        details: 'wms_dev/#dev-infra-alerts',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
      },
      {
        id: '3',
        name: 'Kontakt wms',
        service: 'EMAIL',
        details: 'kontakt@wmsdev.pl',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 21:37',
      },
    ];
    return mockNotificatoinChannels;
  }

  public async getApplications(): Promise<Application[]> {
    const mockApplications: Array<Application> = [
      {
        name: 'alerts-api-database',
        running: true,
        accuracy: 'HIGH',
        customPrompt: 'ignore s3 logs...',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
      },
      {
        name: 'alerts-api-backend',
        running: false,
        accuracy: 'LOW',
        customPrompt: '',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
      },
      {
        name: 'is-jsos-down',
        running: true,
        accuracy: 'MEDIUM',
        customPrompt: 'dont ignore s3 logs...',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
      },
    ];
    return mockApplications;
  }

  public async getNodes(): Promise<Node[]> {
    const mockNodes: Array<Node> = [
      {
        name: 'node 1',
        running: true,
        accuracy: 'HIGH',
        customPrompt: 'ignore s3 logs...',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
      },
      {
        name: 'node 2',
        running: true,
        accuracy: 'LOW',
        customPrompt: 'ignore s3 logs...',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
      },
      {
        name: 'node 3',
        running: false,
        accuracy: 'MEDIUM',
        customPrompt: '',
        updated: '07.03.2024 15:32',
        added: '07.03.2024 15:32',
      },
    ];
    return mockNodes;
  }
}

export const ManagmentServiceApiInstance = new ManagmentServiceApi();
