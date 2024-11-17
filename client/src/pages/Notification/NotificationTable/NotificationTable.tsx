import './NotificationTable.scss';

export interface NotificationTableRowProps {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
  action?: string;
  [key: string]: string | number | undefined;
}

export interface WebhookTableRowProps extends NotificationTableRowProps {
  webhookUrl: string;
}

export interface EmailTableRowProps extends NotificationTableRowProps {
  email: string;
}
