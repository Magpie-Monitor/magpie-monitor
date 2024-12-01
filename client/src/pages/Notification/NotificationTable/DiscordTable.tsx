import Table, { TableColumn } from 'components/Table/Table';
import NotificationButtons from 'pages/Notification/NotificationButtons/NotificationButtons';
import HiddenWebhook from 'pages/Notification/HiddenWebhook/HiddenWebhook';
import NotificationNameLink from 'pages/Notification/NotificationNameLink/NotificationNameLink';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import discordIcon from 'assets/discord-icon.png';
import { NotificationTableRowProps } from './NotificationTable';
import { useEffect, useState } from 'react';
import {
  DiscordNotificationChannel,
  ManagmentServiceApiInstance,
} from 'api/managment-service';
import LoadingTable from './LoadingTable';
import NewDiscordChannelPopup from 'pages/Notification/NewChannelPopup/NewDiscordChannelPopup';
import { dateTimeFromTimestampMs } from 'lib/date';
import EditDiscordChannelPopup from 'pages/Notification/EditChannelPopup/EditDiscordChannelPopup';
import './NotificationTable.scss';
import { useToast } from 'providers/ToastProvider/ToastProvider';

interface DiscordTableRowProps extends NotificationTableRowProps {
  webhookUrl: string;
}

const getDiscordChannelTableRow = ({
  id,
  receiverName,
  updatedAt,
  createdAt,
  webhookUrl,
}: DiscordNotificationChannel): DiscordTableRowProps => ({
  name: receiverName,
  updatedAt: dateTimeFromTimestampMs(updatedAt),
  createdAt: dateTimeFromTimestampMs(createdAt),
  webhookUrl: webhookUrl,
  id,
});

const DiscordTable = () => {
  const [rows, setRows] = useState<DiscordTableRowProps[]>([]);
  const [isLoading, setLoading] = useState<boolean>(true);
  const [isNewChannelPopupDisplayed, setIsNewChannelPopupDisplayed] =
    useState<boolean>(false);
  const [isEditChannelPopupDisplayed, setIsEditChannelPopupDisplayed] =
    useState<boolean>(false);
  const [editChannelPopupData, setEditChannelPopupData] =
    useState<DiscordTableRowProps | null>(null);

  const { showMessage } = useToast();

  const fetchDiscordChannels = async () => {
    try {
      const channels = await ManagmentServiceApiInstance.getDiscordChannels();
      setRows(channels.map(getDiscordChannelTableRow));
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Error fetching discord channels: ', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDiscordChannels();
  }, [isLoading]);

  const discordColumns: Array<TableColumn<DiscordTableRowProps>> = [
    {
      header: 'Name',
      columnKey: 'receiverName',
      customComponent: ({ name }: DiscordTableRowProps) => (
        <NotificationNameLink linkName={name} />
      ),
    },
    {
      header: 'Webhook url',
      columnKey: 'webhookUrl',
      customComponent: ({ webhookUrl }: DiscordTableRowProps) => (
        <HiddenWebhook url={webhookUrl} />
      ),
    },
    {
      header: 'Created at',
      columnKey: 'createdAt',
    },
    {
      header: 'Updated at',
      columnKey: 'updatedAt',
    },
    {
      header: 'Actions',
      columnKey: 'action',
      customComponent: (props: DiscordTableRowProps) => (
        <NotificationButtons
          onUpdate={() => {
            setEditChannelPopupData(props);
            setIsEditChannelPopupDisplayed(true);
          }}
          onTest={async () => {
            try {
              await ManagmentServiceApiInstance.testDiscordChannel(props.id);
              showMessage({
                message: 'Successfully sent a test notification',
                type: 'INFO',
              });
            } catch (e: unknown) {
              showMessage({
                message: `Failed to send test notification: ${e}`,
                type: 'ERROR',
              });
            }
          }}
          onDelete={async () => {
            try {
              await ManagmentServiceApiInstance.deleteDiscordChannel(props.id);
              showMessage({
                message: 'Discord channel was deleted',
                type: 'WARNING',
              });
            } catch (e: unknown) {
              showMessage({
                message: `Failed to delete notification channel: ${e}`,
                type: 'ERROR',
              });
            }
            setLoading(true);
          }}
        />
      ),
    },
  ];

  return (
    <SectionComponent
      icon={<img src={discordIcon} className="notification-table__icon" />}
      title="Discord"
      callback={() => {
        setIsNewChannelPopupDisplayed(true);
      }}
    >
      <LoadingTable isLoading={isLoading}>
        {rows.length > 0 ? (
          <Table columns={discordColumns} rows={rows} alignLeft={false} />
        ) : (
          <p>No Discord channels was yet configureddd</p>
        )}
      </LoadingTable>

      {isEditChannelPopupDisplayed && editChannelPopupData && (
        <EditDiscordChannelPopup
          id={editChannelPopupData.id}
          name={editChannelPopupData.name}
          webhookUrl={editChannelPopupData.webhookUrl}
          isDisplayed={isEditChannelPopupDisplayed}
          setIsDisplayed={setIsEditChannelPopupDisplayed}
          onSubmit={() => setLoading(true)}
        />
      )}
      {isNewChannelPopupDisplayed && (
        <NewDiscordChannelPopup
          setIsDisplayed={setIsNewChannelPopupDisplayed}
          isDisplayed={isNewChannelPopupDisplayed}
          onSubmit={() => setLoading(true)}
        />
      )}
    </SectionComponent>
  );
};

export default DiscordTable;
