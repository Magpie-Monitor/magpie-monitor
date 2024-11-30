import axios, {AxiosInstance} from 'axios';

export interface EmailChannelForm {
  name: string;
  email: string;
}

export interface DiscordChannelForm {
  name: string;
  webhookUrl: string;
}

export interface SlackChannelForm {
  name: string;
  webhookUrl: string;
}

export interface NotificationChannel {
  id: string;
  receiverName: string;
  createdAt: number;
  updatedAt: number;
}

export interface SlackNotificationChannel extends NotificationChannel {
  webhookUrl: string;
}

export interface EmailNotificationChannel extends NotificationChannel {
  receiverEmail: string;
}

export interface DiscordNotificationChannel extends NotificationChannel {
  webhookUrl: string;
}

interface UserInfo {
  nickname: string;
  email: string;
}

interface TokenInfo {
  expTime: number;
}

export type AccuracyLevel = 'HIGH' | 'MEDIUM' | 'LOW';
export type UrgencyLevel = 'HIGH' | 'MEDIUM' | 'LOW';
export type ReportType = 'ON_DEMAND' | 'SCHEDULED';

export interface ReportAwaitingGeneration {
  clusterId: string;
  reportType: ReportType;
  sinceMs: number;
  toMs: number;
  [key: string]: string | number;
}

export interface ReportSummary {
  id: string;
  clusterId: string;
  title: string;
  urgency: UrgencyLevel | null;
  requestedAtMs: number;
  sinceMs: number;
  toMs: number;
  [key: string]: string | number | null;
}

export interface ReportDetails {
  id: string;
  clusterId: string;
  title: string;
  urgency: UrgencyLevel;
  totalApplicationEntries: number;
  totalNodeEntries: number;
  analyzedApplications: number;
  analyzedNodes: number;
  sinceMs: number;
  toMs: number;
}

export interface ClusterSummary {
  clusterId: string;
  running: boolean;
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

export type NotificationChannelKind = 'SLACK' | 'DISCORD' | 'EMAIL';

export interface NotificationChannel {
  id: string;
  name: string;
  service: NotificationChannelKind;
  details: string;
  updated: number;
  added: number;
}

export interface Application {
  name: string;
  running: boolean;
  kind: string;
}

export interface Node {
  name: string;
  running: boolean;
}

export interface Slack {
  id: string;
  receiverName: string;
  webhookUrl: string;
  createdAt: number;
  updatedAt: number;
}

export interface Discord {
  id: string;
  receiverName: string;
  webhookUrl: string;
  createdAt: number;
  updatedAt: number;
}

export interface Email {
  id: string;
  receiverName: string;
  receiverEmail: string;
  createdAt: number;
  updatedAt: number;
}

const MANAGMENT_SERVICE_URL = import.meta.env.VITE_BACKEND_URL;
const VALID_URGENCY_VALUES: ReportSummary['urgency'][] = [
  'HIGH',
  'MEDIUM',
  'LOW',
];

export interface ApplicationIncident {
  id: string;
  clusterId: string;
  title: string;
  category: string;
  applicationName: string;
  summary: string;
  customPrompt: string;
  urgency: UrgencyLevel;
  accuracy: AccuracyLevel;
  recommendation: string;
  sources: ApplicationIncidentSource[];
}

export interface ApplicationIncidentSource {
  podName: string;
  containerName: string;
  image: string;
  content: string;
  timestamp: number;
}

export interface NodeIncident {
  id: string;
  clusterId: string;
  title: string;
  nodeName: string;
  category: string;
  urgency: UrgencyLevel;
  customPrompt: string;
  accuracy: AccuracyLevel;
  summary: string;
  recommendation: string;
  sources: NodeIncidentSource[];
}

export interface AllIncidentsFromReport {
  applicationIncidents: ApplicationIncident[];
  nodeIncidents: NodeIncident[];
}

export interface NodeIncidentSource {
  nodeName: string;
  filename: string;
  content: string;
  timestamp: number;
}

export interface ReportPost {
  clusterId: string;
  accuracy: AccuracyLevel;
  sinceMs: number;
  toMs: number;
  slackReceiverIds: number[];
  discordReceiverIds: number[];
  emailReceiverIds: number[];
  applicationConfigurations: {
    applicationName: string;
    customPrompt: string;
    accuracy: AccuracyLevel;
  }[];
  nodeConfigurations: {
    nodeName: string;
    customPrompt: string;
    accuracy: AccuracyLevel;
  }[];
}

export interface ClusterUpdateData {
  id: string;
  accuracy: AccuracyLevel;
  isEnabled: boolean;
  generatedEveryMillis: number;
  slackReceiverIds: number[];
  discordReceiverIds: number[];
  emailReceiverIds: number[];
  applicationConfigurations: {
    name: string;
    kind: string;
    accuracy: AccuracyLevel;
    customPrompt: string;
  }[];
  nodeConfigurations: {
    name: string;
    accuracy: AccuracyLevel;
    customPrompt: string;
  }[];
}

export interface ApplicationConfiguration {
  name: string;
  kind: string;
  accuracy: AccuracyLevel;
  customPrompt: string;
}

export interface NodeConfiguration {
  name: string;
  accuracy: AccuracyLevel;
  customPrompt: string;
}

export interface ClusterDetails {
  id: string;
  accuracy: AccuracyLevel;
  isEnabled: boolean;
  running: boolean;
  generatedEveryMillis: number;
  slackReceivers: Slack[];
  discordReceivers: Discord[];
  emailReceivers: Email[];
  applicationConfigurations: ApplicationConfiguration[];
  nodeConfigurations: NodeConfiguration[];
}

export interface EditDiscordChannelForm {
  id: string;
  name: string;
  webhookUrl: string;
}

export interface EditSlackChannelForm {
  id: string;
  name: string;
  webhookUrl: string;
}

export interface EditEmailChannelForm {
  id: string;
  name: string;
  email: string;
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
    return this.getTokenInfoWithoutRefreshing();
  }

  public async getUserInfo(): Promise<UserInfo> {
    await this.refreshTokenIfExpired();
    const user = await this.axiosInstance.get('/api/v1/auth/user-details');
    return user.data;
  }

  public async getReport(id: string): Promise<ReportDetails> {
    await this.refreshTokenIfExpired();
    const report = await this.axiosInstance.get(`/api/v1/reports/${id}`);
    return report.data;
    // return {
    //   analyzedApplications: 123,
    //   analyzedNodes: 20,
    //   id: 'report-1',
    //   clusterId: 'cluster-1',
    //   title: 'Report for this time',
    //   urgency: 'HIGH',
    //   totalApplicationEntries: 1000,
    //   totalNodeEntries: 1000,
    //   sinceMs: 212414123,
    //   toMs: 212414123,
    // };
  }

  public async getIncidentsFromReport(
    reportId: string,
  ): Promise<AllIncidentsFromReport> {
    await this.refreshTokenIfExpired();
    const report = await this.axiosInstance.get(
      `/api/v1/reports/${reportId}/incidents`,
    );
    return report.data;
    // return {
    //   applicationIncidents: [
    //     {
    //       clusterId: 'Cluster 1',
    //       id: 'udd',
    //       title: 'Something wrong with app-1',
    //       urgency: 'HIGH',
    //       category: 'Serious category',
    //       customPrompt: 'Custom prompt',
    //       accuracy: 'HIGH',
    //       applicationName: 'application-1',
    //       summary: 'This is an summar of the incident',
    //       recommendation: 'This is an recommendation of the incident',
    //       sources: [
    //         {
    //           container: 'container-1',
    //           timestamp: 213213124,
    //           pod: 'pod-1',
    //           image: 'image-1',
    //           content: 'LOGS',
    //         },
    //       ],
    //     },
    //   ],
    //   nodeIncidents: [
    //     {
    //       id: 'uid',
    //       clusterId: 'cluster-1',
    //       category: 'Serious category',
    //       title: 'Something wrong with lke-123',
    //       nodeName: 'lke-123213213',
    //       urgency: 'HIGH',
    //       accuracy: 'HIGH',
    //       customPrompt: 'Custom prompt',
    //       summary: 'Node incident summary',
    //       recommendation: 'Node incident recommendation',
    //       sources: [
    //         {
    //           timestamp: 123123123,
    //           nodeName: 'lke-1231231',
    //           content: 'LOGSGLOGSLGOS',
    //           filename: 'file1/tmp/file2',
    //         },
    //       ],
    //     },
    //   ],
    // };
  }

  public async getReports(reportType: ReportType): Promise<ReportSummary[]> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get('/api/v1/reports', {
      params: {
        reportType,
      },
    });
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

  public async getAwaitingGenerationReports(): Promise<ReportAwaitingGeneration[]> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get('/api/v1/reports/await-generation');
    return response.data;
  }

  public async getApplicationIncident(
    id: string,
  ): Promise<ApplicationIncident> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get(
      `/api/v1/reports/application-incidents/${id}`,
    );

    return response.data;
  }

  public async getNodeIncident(id: string): Promise<NodeIncident> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get(
      `/api/v1/reports/node-incidents/${id}`,
    );

    return response.data;
  }

  public async getClusters(): Promise<ClusterSummary[]> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get('api/v1/clusters');
    const clusters = response.data;

    return clusters.map((cluster: ClusterSummary) => {
      return {
        ...cluster,
        updatedAt: 0,
        accuracy: 'LOW',
        slackChannels: [],
        discordChannels: [],
        mailChannels: [],
      };
    });
  }

  public async getNotificationChannels(): Promise<NotificationChannel[]> {
    await this.refreshTokenIfExpired();
    const [slack, discord, mail] = await Promise.all([
      this.axiosInstance.get('/api/v1/notification-channels/slack'),
      this.axiosInstance.get('/api/v1/notification-channels/discord'),
      this.axiosInstance.get('/api/v1/notification-channels/mails'),
    ]);
    const slackChannels = slack.data.map((channel: Slack) => ({
      id: channel.id,
      name: channel.receiverName,
      service: 'SLACK',
      details: channel.webhookUrl,
      updated: channel.updatedAt,
      added: channel.createdAt,
    }));

    const discordChannels = discord.data.map((channel: Discord) => ({
      id: channel.id,
      name: channel.receiverName,
      service: 'DISCORD',
      details: channel.webhookUrl,
      updated: channel.updatedAt,
      added: channel.createdAt,
    }));

    const mailChannels = mail.data.map((channel: Email) => ({
      id: channel.id.toString(),
      name: channel.receiverName,
      service: 'EMAIL',
      details: channel.receiverEmail,
      updated: channel.updatedAt,
      added: channel.createdAt,
    }));

    return [...slackChannels, ...discordChannels, ...mailChannels];

    // const mockNotificatoinChannels: Array<NotificationChannel> = [
    //   {
    //     id: '1',
    //     name: 'Infra team slack',
    //     service: 'SLACK',
    //     details: 'wms_dev/#infra-alerts',
    //     updated: '07.03.2024 15:32',
    //     added: '07.03.2024 15:32',
    //   },
    //   {
    //     id: '1',
    //     name: 'Infra team discord',
    //     service: 'DISCORD',
    //     details: 'wms_dev/#dev-infra-alerts',
    //     updated: '07.03.2024 15:32',
    //     added: '07.03.2024 15:32',
    //   },
    //   {
    //     id: '1',
    //     name: 'Kontakt wms',
    //     service: 'EMAIL',
    //     details: 'kontakt@wmsdev.pl',
    //     updated: '07.03.2024 15:32',
    //     added: '07.03.2024 21:37',
    //   },
    // ];
    // return mockNotificatoinChannels;
  }

  public async getApplications(clusterId: string): Promise<Application[]> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get(
      `/api/v1/clusters/${clusterId}/applications`,
    );

    return response.data;
    // console.log(clusterId);
    // const mockApplications: Array<Application> = [
    //   {
    //     name: 'alerts-api-database',
    //     running: true,
    //     kind: 'Deployment',
    //   },
    //   {
    //     name: 'alerts-api-backend',
    //     running: false,
    //     kind: 'Deployment',
    //   },
    //   {
    //     name: 'is-jsos-down',
    //     running: true,
    //     kind: 'Deployment',
    //   },
    // ];
    // return mockApplications;
  }

  public async getNodes(clusterId: string): Promise<Node[]> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get(
      `/api/v1/clusters/${clusterId}/nodes`,
    );

    return response.data;
    // console.log(clusterId);
    // const mockNodes: Array<Node> = [
    //   {
    //     name: 'node 1',
    //     running: true,
    //   },
    //   {
    //     name: 'node 2',
    //     running: true,
    //   },
    //   {
    //     name: 'node 3',
    //     running: false,
    //   },
    // ];
    // return mockNodes;
  }

  public async generateOnDemandReport(reportData: ReportPost): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.post('/api/v1/reports', reportData);
  }

  public async getClusterDetails(clusterId: string): Promise<ClusterDetails> {
    await this.refreshTokenIfExpired();
    const response = await this.axiosInstance.get(
      `/api/v1/clusters/${clusterId}`,
    );

    const clusterData = response.data;

    return {
      id: clusterData.id,
      accuracy: clusterData.accuracy,
      isEnabled: clusterData.isEnabled,
      running: clusterData.running,
      generatedEveryMillis: clusterData.generatedEveryMillis,
      slackReceivers: clusterData.slackReceivers.map((receiver: Slack) => ({
        id: receiver.id,
        receiverName: receiver.receiverName,
        webhookUrl: receiver.webhookUrl,
        createdAt: receiver.createdAt,
        updatedAt: receiver.updatedAt,
      })),
      discordReceivers: clusterData.discordReceivers.map(
        (receiver: Discord) => ({
          id: receiver.id,
          receiverName: receiver.receiverName,
          webhookUrl: receiver.webhookUrl,
          createdAt: receiver.createdAt,
          updatedAt: receiver.updatedAt,
        }),
      ),
      emailReceivers: clusterData.emailReceivers.map((receiver: Email) => ({
        id: receiver.id,
        receiverName: receiver.receiverName,
        receiverEmail: receiver.receiverEmail,
        createdAt: receiver.createdAt,
        updatedAt: receiver.updatedAt,
      })),
      applicationConfigurations: clusterData.applicationConfigurations.map(
        (config: ApplicationConfiguration) => ({
          name: config.name,
          kind: config.kind,
          accuracy: config.accuracy,
          customPrompt: config.customPrompt,
        }),
      ),
      nodeConfigurations: clusterData.nodeConfigurations.map(
        (config: NodeConfiguration) => ({
          name: config.name,
          accuracy: config.accuracy,
          customPrompt: config.customPrompt,
        }),
      ),
    };
  }

  public async updateCluster(clusterData: ClusterUpdateData): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.put('/api/v1/clusters', clusterData);
  }

  public async scheduleReport(
    clusterId: string,
    periodMs: number,
  ): Promise<void> {
    await this.refreshTokenIfExpired();
    const requestPayload = {
      clusterId,
      periodMs,
    };
    await this.axiosInstance.post('/api/v1/reports/schedule', requestPayload);
  }

  public async getSlackChannels(): Promise<SlackNotificationChannel[]> {
    await this.refreshTokenIfExpired();
    const slackChannels = await this.axiosInstance.get(
      '/api/v1/notification-channels/slack',
    );
    return slackChannels.data;
  }

  public async getDiscordChannels(): Promise<DiscordNotificationChannel[]> {
    await this.refreshTokenIfExpired();
    const discordChannels = await this.axiosInstance.get(
      '/api/v1/notification-channels/discord',
    );
    return discordChannels.data;
  }

  public async getEmailChannels(): Promise<EmailNotificationChannel[]> {
    await this.refreshTokenIfExpired();
    const emailChannels = await this.axiosInstance.get(
      '/api/v1/notification-channels/mails',
    );
    return emailChannels.data;
  }

  public async postSlackChannel(slackChannel: SlackChannelForm): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.post(
      '/api/v1/notification-channels/slack',
      slackChannel,
    );
  }

  public async postDiscordChannel(
    discordChannel: DiscordChannelForm,
  ): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.post(
      '/api/v1/notification-channels/discord',
      discordChannel,
    );
  }

  public async postEmailChannel(emailChannel: EmailChannelForm): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.post(
      '/api/v1/notification-channels/mails',
      emailChannel,
    );
  }

  public async editDiscordChannel(data: EditDiscordChannelForm): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.patch(
      `/api/v1/notification-channels/discord/${data.id}`,
      {
        name: data.name,
        webhookUrl: data.webhookUrl === '' ? undefined : data.webhookUrl,
      },
    );
  }

  public async editSlackChannel(data: EditSlackChannelForm): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.patch(
      `/api/v1/notification-channels/slack/${data.id}`,
      {
        name: data.name,
        webhookUrl: data.webhookUrl === '' ? undefined : data.webhookUrl,
      },
    );
  }

  public async editEmailChannel(data: EditEmailChannelForm): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.patch(
      `/api/v1/notification-channels/mails/${data.id}`,
      {
        name: data.name,
        email: data.email,
      },
    );
  }

  public async testDiscordChannel(id: string): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.post(
      `/api/v1/notification-channels/discord/${id}/test-notification`,
    );
  }

  public async testSlackChannel(id: string): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.post(
      `/api/v1/notification-channels/slack/${id}/test-notification`,
    );
  }

  public async testEmailChannel(id: string): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.post(
      `/api/v1/notification-channels/mails/${id}/test-notification`,
    );
  }

  public async deleteDiscordChannel(id: string): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.delete(
      `/api/v1/notification-channels/discord/${id}`,
    );
  }

  public async deleteSlackChannel(id: string): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.delete(
      `/api/v1/notification-channels/slack/${id}`,
    );
  }

  public async deleteEmailChannel(id: string): Promise<void> {
    await this.refreshTokenIfExpired();
    await this.axiosInstance.delete(
      `/api/v1/notification-channels/mails/${id}`,
    );
  }
}

export const ManagmentServiceApiInstance = new ManagmentServiceApi();
