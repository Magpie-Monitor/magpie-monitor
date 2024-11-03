import Table, { TableColumn } from 'components/Table/Table';
import NotificationButtons from 'pages/Notification/NotificationButtons/NotificationButtons';
import HiddenWebhook from 'pages/Notification/HiddenWebhook/HiddenWebhook';
import NotificationNameLink from 'pages/Notification/NotificationNameLink/NotificationNameLink';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import discordIcon from 'assets/discord-icon.png';
import { NotificationTableRowProps } from './NotificationTable';
import { useEffect, useState } from 'react';
import { ManagmentServiceApiInstance } from 'api/managment-service';
import LoadingTable from './LoadingTable';

export interface DiscordTableRowProps extends NotificationTableRowProps {
  webhookUrl: string;
}

const discordColumns: Array<TableColumn<DiscordTableRowProps>> = [
  {
    header: 'Name',
    columnKey: 'receiverName',
    customComponent: ({ receiverName, destination }: DiscordTableRowProps) => (
      <NotificationNameLink linkName={receiverName} destination={destination} />
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
    }: DiscordTableRowProps) => (
      <NotificationButtons
        channel={'DISCORD'}
        adress={webhookUrl}
        linkName={linkName}
        createdAt={createdAt}
        updateAt={updateAt}
      />
    ),
  },
];

const DiscordTable = () => {
  const [rows, setRows] = useState<DiscordTableRowProps[]>([]);
  const [isLoading, setLoading] = useState<boolean>(true);

  const fetchDiscordChannels = async () => {
    try {
      const channels = await ManagmentServiceApiInstance.getDiscordChannels();
      setRows(channels);
    } catch (error) {
      console.error('Error fetching discord channels: ', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDiscordChannels();
  }, []);

  return (
    <SectionComponent icon={<img src={discordIcon} />} title="Discord">
      <LoadingTable isLoading={isLoading}>
        {rows.length > 0 ? (
          <Table columns={discordColumns} rows={rows} alignLeft={false} />
        ) : (
          <p>No Discord channels was yet configureddd</p>
        )}
      </LoadingTable>
    </SectionComponent>
  );
};

export default DiscordTable;
