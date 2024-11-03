import Table, { TableColumn } from 'components/Table/Table';
import NotificationButtons from 'pages/Notification/NotificationButtons/NotificationButtons';
import HiddenWebhook from 'pages/Notification/HiddenWebhook/HiddenWebhook';
import NotificationNameLink from 'pages/Notification/NotificationNameLink/NotificationNameLink';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import slackIcon from 'assets/slack-icon.png';
import { NotificationTableRowProps } from './NotificationTable';
import { useEffect, useState } from 'react';
import { ManagmentServiceApiInstance } from 'api/managment-service';
import LoadingTable from './LoadingTable';
import {
  NotificationContextProps,
  useNotification,
} from 'pages/Notification/NotificationContext';
import AddSlackChannelPopup from 'pages/Notification/AddNewChannelPopup/AddSlackChannelPopup';

export interface SlackTableRowProps extends NotificationTableRowProps {
  webhookUrl: string;
}

const slackColumns: Array<TableColumn<SlackTableRowProps>> = [
  {
    header: 'Name',
    columnKey: 'receiverName',
    customComponent: ({ receiverName, destination }: SlackTableRowProps) => (
      <NotificationNameLink linkName={receiverName} destination={destination} />
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
    header: 'Update at',
    columnKey: 'updateAt',
  },
  {
    header: 'Actions',
    columnKey: 'action',
    customComponent: ({
      webhookUrl,
      linkName,
      createdAt,
      updateAt,
    }: SlackTableRowProps) => (
      <NotificationButtons
        channel={'SLACK'}
        adress={webhookUrl}
        linkName={linkName}
        createdAt={createdAt}
        updateAt={updateAt}
      />
    ),
  },
];

const SlackTable = () => {
  const [rows, setRows] = useState<SlackTableRowProps[]>([]);
  const [isLoading, setLoading] = useState<boolean>(true);
  const {
    hidePopup,
    createNewChannel,
  }: NotificationContextProps = useNotification();

  const fetchSlackChannels = async () => {
    try {
      const channels = await ManagmentServiceApiInstance.getSlackChannels();
      setRows(channels);
    } catch (error) {
      console.error('Error fetching slack channels: ', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSlackChannels();
  }, []);

  return (
    <SectionComponent
      icon={<img src={slackIcon} />}
      title="Slack"
      callback={() => {
        createNewChannel(
          <AddSlackChannelPopup
            isDisplayed={true}
            setIsDisplayed={hidePopup}
          />,
        );
      }}
    >
      <LoadingTable isLoading={isLoading}>
        {rows.length > 0 ? (
          <Table columns={slackColumns} rows={rows} alignLeft={false} />
        ) : (
          <p>No Slack channels was yet configured</p>
        )}
      </LoadingTable>
    </SectionComponent>
  );
};

export default SlackTable;
