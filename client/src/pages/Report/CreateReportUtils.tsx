import {NotificationChannel} from './NotificationSection/NotificationSection';
import {ApplicationDataRow} from './ApplicationSection/ApplicationSection';
import {NodeDataRow} from './NodesSection/NodesSection';
import {
    AccuracyLevel,
    ClusterUpdateData,
    ManagmentServiceApiInstance,
    NotificationChannelKind,
    ReportPost,
    ReportType,
} from 'api/managment-service.ts';
import {periodToMilliseconds} from './SchedulePeriod/SchedulePeriod.tsx';
import {dateFromTimestampMs} from 'lib/date.ts';

export const fetchClusterData = async (
    clusterId: string
): Promise<{
    notificationChannels: NotificationChannel[];
    applications: ApplicationDataRow[];
    nodes: NodeDataRow[];
}> => {
    try {
        const [clusterDetails, runningApplications, runningNodes] = await Promise.all([
            ManagmentServiceApiInstance.getClusterDetails(clusterId),
            ManagmentServiceApiInstance.getApplications(clusterId),
            ManagmentServiceApiInstance.getNodes(clusterId),
        ]);

        const runningApplicationsMap = runningApplications.reduce((acc, app) => {
            acc[`${app.name}-${app.kind}`] = app.running;
            return acc;
        }, {} as Record<string, boolean>);

        const runningNodesMap = runningNodes.reduce((acc, node) => {
            acc[node.name] = node.running;
            return acc;
        }, {} as Record<string, boolean>);

        const notificationChannels = [
            ...clusterDetails.slackReceivers.map(receiver => ({
                id: receiver.id.toString(),
                name: receiver.receiverName,
                details: receiver.webhookUrl,
                service: 'SLACK' as NotificationChannelKind,
                added: dateFromTimestampMs(receiver.createdAt),
                updated: dateFromTimestampMs(receiver.updatedAt),
            })),
            ...clusterDetails.discordReceivers.map(receiver => ({
                id: receiver.id.toString(),
                name: receiver.receiverName,
                details: receiver.webhookUrl,
                service: 'DISCORD' as NotificationChannelKind,
                added: dateFromTimestampMs(receiver.createdAt),
                updated: dateFromTimestampMs(receiver.updatedAt),
            })),
            ...clusterDetails.emailReceivers.map(receiver => ({
                id: receiver.id.toString(),
                name: receiver.receiverName,
                details: receiver.receiverEmail,
                service: 'EMAIL' as NotificationChannelKind,
                added: dateFromTimestampMs(receiver.createdAt),
                updated: dateFromTimestampMs(receiver.updatedAt),
            })),
        ];

        const applications = clusterDetails.applicationConfigurations.map(config => ({
            name: config.name,
            kind: config.kind,
            accuracy: config.accuracy as AccuracyLevel,
            customPrompt: config.customPrompt,
            running: runningApplicationsMap[`${config.name}-${config.kind}`] ?? false,
        }));

        const nodes = clusterDetails.nodeConfigurations.map(config => ({
            name: config.name,
            accuracy: config.accuracy,
            customPrompt: config.customPrompt,
            running: runningNodesMap[config.name] ?? false,
        }));

        return {notificationChannels, applications, nodes};
    } catch (error) {
        console.error('Failed to fetch cluster data:', error);
        return {notificationChannels: [], applications: [], nodes: []};
    }
};

export const filterNotificationChannels = (channels: NotificationChannel[]) => {
    const slackReceiverIds: number[] = [];
    const discordReceiverIds: number[] = [];
    const mailReceiverIds: number[] = [];

    channels.forEach(channel => {
        const channelId = parseInt(channel.id, 10);
        if (!isNaN(channelId)) {
            switch (channel.service) {
                case 'SLACK':
                    slackReceiverIds.push(channelId);
                    break;
                case 'DISCORD':
                    discordReceiverIds.push(channelId);
                    break;
                case 'EMAIL':
                    mailReceiverIds.push(channelId);
                    break;
                default:
                    console.warn(`Unknown service: ${channel.service}`);
            }
        } else {
            console.warn(`Invalid channel id: ${channel.id}`);
        }
    });

    return {slackReceiverIds, discordReceiverIds, mailReceiverIds};
};


export const generateReport = ({
                                   id,
                                   notificationChannels,
                                   applications,
                                   nodes,
                                   generationType,
                                   accuracy,
                                   generationPeriod,
                                   startDateMs,
                                   endDateMs,
                               }: {
    id: string | undefined;
    notificationChannels: NotificationChannel[];
    applications: ApplicationDataRow[];
    nodes: NodeDataRow[];
    generationType: ReportType;
    accuracy: AccuracyLevel;
    generationPeriod: string;
    startDateMs: number;
    endDateMs: number;
}) => {
    const {slackReceiverIds, discordReceiverIds, mailReceiverIds} =
        filterNotificationChannels(notificationChannels);

    const schedulePeriodMs = periodToMilliseconds[generationPeriod] || 0;

    if (generationType === 'ON-DEMAND') {
        const report: ReportPost = {
            clusterId: id ?? '',
            accuracy: 'HIGH',
            sinceMs: startDateMs,
            toMs: endDateMs,
            slackReceiverIds,
            discordReceiverIds,
            emailReceiverIds: mailReceiverIds,
            applicationConfigurations: applications.map(app => ({
                applicationName: app.name,
                accuracy: app.accuracy,
                customPrompt: app.customPrompt,
            })),
            nodeConfigurations: nodes.map(node => ({
                nodeName: node.name,
                accuracy: node.accuracy,
                customPrompt: node.customPrompt,
            })),
        };
        ManagmentServiceApiInstance.generateOnDemandReport(report);
    } else if (generationType === 'SCHEDULED' && id) {
        const clusterUpdateData: ClusterUpdateData = {
            id,
            accuracy,
            isEnabled: true,
            generatedEveryMillis: schedulePeriodMs,
            slackReceiverIds,
            discordReceiverIds,
            emailReceiverIds: mailReceiverIds,
            applicationConfigurations: applications.map(app => ({
                name: app.name,
                kind: app.kind,
                accuracy: app.accuracy,
                customPrompt: app.customPrompt,
            })),
            nodeConfigurations: nodes.map(node => ({
                name: node.name,
                accuracy: node.accuracy,
                customPrompt: node.customPrompt,
            })),
        };
        ManagmentServiceApiInstance.updateCluster(clusterUpdateData);
        ManagmentServiceApiInstance.scheduleReport(clusterUpdateData.id,
            clusterUpdateData.generatedEveryMillis);
    }
};
