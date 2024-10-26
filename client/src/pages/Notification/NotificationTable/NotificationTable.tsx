import Table, { TableColumn } from 'components/Table/Table';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import './NotificationTable.scss';
import NotificationButtons from 'pages/Notification/NotificationButtons/NotificationButtons';
import HiddenWebhook from 'pages/Notification/HiddenWebhook/HiddenWebhook';
import NotificationNameLink from 'pages/Notification/NotificationNameLink/NotificationNameLink';
import { NotificationsChannel } from 'pages/Notification/NotificationContext';
import EmailColumn from 'pages/Notification/EmailCell/EmailCell';

interface NotificationTableRowProps {
  linkName: string;
  destination: string;
  createdAt: string;
  updateAt: string;
  action: string;
  [key: string]: string;
}

export interface WebhookTableRowProps extends NotificationTableRowProps {
  webhookUrl: string;
}

export interface EmailTableRowProps extends NotificationTableRowProps {
  email: string;
}

const isWebhookTableRowProps = (
  notificationTableRowProps: NotificationTableRowProps,
): notificationTableRowProps is WebhookTableRowProps => {
  return (
    (notificationTableRowProps as WebhookTableRowProps).webhookUrl !== undefined
  );
};

const isEmailTableRowProps = (
  notificationTableRowProps: NotificationTableRowProps,
): notificationTableRowProps is EmailTableRowProps => {
  return (notificationTableRowProps as EmailTableRowProps).email !== undefined;
};

interface NotificationTableProps {
  data: NotificationTableRowProps[];
  imageName: string;
  header: string;
  channel: NotificationsChannel;
}

const getTableColumnsForWebhookNotificationChannel = (channel: NotificationsChannel): 
Array<TableColumn<NotificationTableRowProps>> => [
  {
    header: 'Name',
    columnKey: 'name',
    customComponent: ({ linkName, destination }: NotificationTableRowProps) => (
      <NotificationNameLink linkName={linkName} destination={destination} />
    ),
  },
  {
    header: 'Webhook url',
    columnKey: 'webhookUrl',
    customComponent: ({ webhookUrl }: NotificationTableRowProps) => (
      <HiddenWebhook url={webhookUrl} />
    ),
  },
  {
    header: 'Created at',
    columnKey: 'createdAt',
  },
  {
    header: 'Update at',
    columnKey: 'updateAt',
  },
  {
    header: 'Actions',
    columnKey: 'action',
    customComponent: ({
      webhookUrl,
      linkName,
      destination,
      createdAt,
      updateAt,
    }: NotificationTableRowProps) => (
      <NotificationButtons
        channel={channel}
        adress={webhookUrl}
        linkName={linkName}
        destination={destination}
        createdAt={createdAt}
        updateAt={updateAt}
      />
    ),
  },
];

const getTableColumnsForEmailNotificationChannel = (channel: NotificationsChannel): Array<
  TableColumn<NotificationTableRowProps>
> => [
  {
    header: 'Name',
    columnKey: 'name',
    customComponent: ({ linkName, destination }: NotificationTableRowProps) => (
      <NotificationNameLink linkName={linkName} destination={destination} />
    ),
  },
  {
    header: 'Email',
    columnKey: 'email',
    customComponent: ({ email }: NotificationTableRowProps) => (
      <EmailColumn email={email} />
    ),
  },
  {
    header: 'Created at',
    columnKey: 'createdAt',
  },
  {
    header: 'Update at',
    columnKey: 'updateAt',
  },
  {
    header: 'Actions',
    columnKey: 'action',
    customComponent: ({
      email,
      linkName,
      destination,
      createdAt,
      updateAt,
    }: NotificationTableRowProps) => (
      <NotificationButtons
        channel={channel}
        adress={email}
        linkName={linkName}
        destination={destination}
        createdAt={createdAt}
        updateAt={updateAt}
      />
    ),
  },
];

const getTableColumns = (
  channel: NotificationsChannel,
  notificationTableRowProps: NotificationTableRowProps[],
): Array<TableColumn<NotificationTableRowProps>> => {
  if (notificationTableRowProps.length < 1)
    throw new Error('NotificationTableRowProps is empty array');

  switch (true) {
    case isWebhookTableRowProps(notificationTableRowProps[0]):
      return getTableColumnsForWebhookNotificationChannel(channel);
    case isEmailTableRowProps(notificationTableRowProps[0]):
      return getTableColumnsForEmailNotificationChannel(channel);
    default:
      throw new Error(
        'NotificationTableRowProps is neither WebhookTableRowProps nor EmailTableRowProps',
      );
  }
};

const NotificationTable = ({
  data,
  imageName,
  header,
  channel,
}: NotificationTableProps) => {
  const convertImageNameToSourcePath = (image: string) =>
    `/src/assets/${image}`;

  return (
    <div className="notification-table">
      <div className="notification-table__heading">
        <div className="notification-table__heading">
          <img src={convertImageNameToSourcePath(imageName)} />
          <p className="notification-table__heading__text">{header}</p>
        </div>
        <button className="notification-table__heading__button">
          <SVGIcon iconName="plus-icon" />
        </button>
      </div>
      <div className="notification-table__line" />
      <Table
        columns={getTableColumns(channel, data)}
        rows={data}
        alignLeft={false}
      />
    </div>
  );
};

export default NotificationTable;
