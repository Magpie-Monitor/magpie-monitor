import React from 'react';
import EntriesSelector from 'components/EntriesSelector/EntriesSelector';
import { NotificationChannel } from 'pages/Report/NotificationSection/NotificationSection.tsx';

import Channels
  from 'pages/Report/NotificationSection/NotificationChannelsColumn/NotificationChannelColumn.tsx';
import {
  transformNotificationChannelToDetailsColumn,
  transformNotificationChannelToServiceColumn,
} from 'pages/Report/NotificationSection/NotificationUtils.tsx';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import { TableColumn } from 'components/Table/Table.tsx';

interface NotificationsEntriesSelectorProps {
  selectedChannels: NotificationChannel[];
  setSelectedChannels: React.Dispatch<
    React.SetStateAction<NotificationChannel[]>
  >;
  channelsToExclude: NotificationChannel[];
  onAdd: () => void;
  onClose: () => void;
  availableChannels: NotificationChannel[];
}

const NotificationsEntriesSelector: React.FC<
  NotificationsEntriesSelectorProps
> = ({
  selectedChannels,
  setSelectedChannels,
  channelsToExclude,
  onAdd,
  onClose,
  availableChannels
}) => {
  const getUniqueKey = (channel: NotificationChannel) =>
    `${channel.id}-${channel.service}`;

  const columns: TableColumn<NotificationChannel>[] = [
    {
      header: 'Name',
      columnKey: 'name',
      customComponent: (row) => <LinkComponent>{row.name}</LinkComponent>,
    },
    {
      header: 'Service',
      columnKey: 'service',
      customComponent: (row) => (
        <Channels channel={transformNotificationChannelToServiceColumn(row)} />
      ),
    },
    {
      header: 'Details',
      columnKey: 'details',
      customComponent: (row) => (
        <Channels channel={transformNotificationChannelToDetailsColumn(row)} />
      ),
    },
  ];

  return (
    <EntriesSelector<NotificationChannel>
      selectedItems={selectedChannels}
      setSelectedItems={setSelectedChannels}
      itemsToExclude={channelsToExclude}
      onAdd={onAdd}
      onClose={onClose}
      items={availableChannels}
      columns={columns}
      getKey={getUniqueKey}
      entityLabel="notification-channel"
      noEntriesMessage={
        <>
          <p>There is no notification channel to add.</p>
          <LinkComponent to="/notifications">
            You can create a new one here.
          </LinkComponent>
        </>
      }
      title="Select Notifications"
    />
  );
};

export default NotificationsEntriesSelector;
