import { NotificationChannel } from 'api/managment-service';

export type NotificationChannelKind = 'SLACK' | 'DISCORD' | 'EMAIL';

export interface NotificationChannelColumn {
    kind: NotificationChannelKind;
    name: string;
}

export const transformNotificationChannelToServiceColumn = (
    notificationChannel: NotificationChannel,
): NotificationChannelColumn => ({
    kind: notificationChannel.service.toUpperCase() as NotificationChannelKind,
    name: notificationChannel.service,
});

export const transformNotificationChannelToDetailsColumn = (
    notificationChannel: NotificationChannel,
): NotificationChannelColumn => ({
    kind: notificationChannel.service.toUpperCase() as NotificationChannelKind,
    name: notificationChannel.details,
});