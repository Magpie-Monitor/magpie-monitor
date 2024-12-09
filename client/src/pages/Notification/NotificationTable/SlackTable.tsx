import Table, { TableColumn } from 'components/Table/Table';
import NotificationButtons from 'pages/Notification/NotificationButtons/NotificationButtons';
import HiddenWebhook from 'pages/Notification/HiddenWebhook/HiddenWebhook';
import NotificationNameLink from 'pages/Notification/NotificationNameLink/NotificationNameLink';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import slackIcon from 'assets/slack-icon.png';
import { NotificationTableRowProps } from './NotificationTable';
import { useEffect, useState } from 'react';
import {
  ManagmentServiceApiInstance,
  SlackNotificationChannel,
} from 'api/managment-service';
import LoadingTable from './LoadingTable';
import NewSlackChannelPopup from 'pages/Notification/NewChannelPopup/NewSlackChannelPopup';
import { dateTimeWithoutSecondsFromTimestampMs } from 'lib/date';
import EditSlackChannelPopup from 'pages/Notification/EditChannelPopup/EditSlackChannelPopup';
import './NotificationTable.scss';
import { useToast } from 'providers/ToastProvider/ToastProvider';
interface SlackTableRowProps extends NotificationTableRowProps {
  webhookUrl: string;
}

const getSlackChannelTableRow = ({
  id,
  receiverName,
  updatedAt,
  createdAt,
  webhookUrl,
}: SlackNotificationChannel): SlackTableRowProps => ({
  name: receiverName,
  updatedAt: dateTimeWithoutSecondsFromTimestampMs(updatedAt),
  createdAt: dateTimeWithoutSecondsFromTimestampMs(createdAt),
  webhookUrl,
  id,
});

const SlackTable = () => {
  const [rows, setRows] = useState<SlackTableRowProps[]>([]);
  const [isLoading, setLoading] = useState<boolean>(true);
  const [isNewChannelPopupDisplayed, setIsNewChannelPopupDisplayed] =
    useState<boolean>(false);
  const [isEditChannelPopupDisplayed, setIsEditChannelPopupDisplayed] =
    useState<boolean>(false);
  const [editChannelPopupData, setEditChannelPopupData] =
    useState<SlackTableRowProps | null>(null);

  const { showMessage } = useToast();

  const fetchSlackChannels = async () => {
    try {
      const channels = await ManagmentServiceApiInstance.getSlackChannels();
      setRows(channels.map(getSlackChannelTableRow));
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Error fetching slack channels: ', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSlackChannels();
  }, [isLoading]);

  const slackColumns: Array<TableColumn<SlackTableRowProps>> = [
    {
      header: 'Name',
      columnKey: 'receiverName',
      customComponent: ({ name }: SlackTableRowProps) => (
        <NotificationNameLink linkName={name} />
      ),
    },
    {
      header: 'Webhook url',
      columnKey: 'webhookUrl',
      customComponent: ({ webhookUrl }: SlackTableRowProps) => (
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
      customComponent: (props: SlackTableRowProps) => (
        <NotificationButtons
          onUpdate={() => {
            setIsEditChannelPopupDisplayed(true);
            setEditChannelPopupData(props);
          }}
          onTest={async () => {
            try {
              await ManagmentServiceApiInstance.testSlackChannel(props.id);
              showMessage({
                message: 'Successfully sent a test notification',
                type: 'INFO',
              });
            } catch (e: unknown) {
              showMessage({
                message: `Failed to send notification: ${e}`,
                type: 'ERROR',
              });
            }
          }}
          onDelete={async () => {
            try {
              await ManagmentServiceApiInstance.deleteSlackChannel(props.id);
              showMessage({
                message: 'Slack channel was deleted',
                type: 'WARNING',
              });
            } catch (e: unknown) {
              showMessage({
                message: `Failed to delete slack channel: ${e}`,
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
      icon={<img src={slackIcon} className="notification-table__icon" />}
      title="Slack"
      callback={() => {
        setIsNewChannelPopupDisplayed(true);
      }}
    >
      <LoadingTable isLoading={isLoading}>
        {rows.length > 0 ? (
          <Table columns={slackColumns} rows={rows} alignLeft={false} />
        ) : (
          <p>No Slack channels was yet configured</p>
        )}
      </LoadingTable>
      {isEditChannelPopupDisplayed && editChannelPopupData && (
        <EditSlackChannelPopup
          id={editChannelPopupData.id}
          name={editChannelPopupData.name}
          webhookUrl={editChannelPopupData.webhookUrl}
          isDisplayed={isEditChannelPopupDisplayed}
          setIsDisplayed={setIsEditChannelPopupDisplayed}
          onSubmit={() => setLoading(true)}
        />
      )}
      {isNewChannelPopupDisplayed && (
        <NewSlackChannelPopup
          setIsDisplayed={setIsNewChannelPopupDisplayed}
          isDisplayed={isNewChannelPopupDisplayed}
          onSubmit={() => setLoading(true)}
        />
      )}
    </SectionComponent>
  );
};

export default SlackTable;
