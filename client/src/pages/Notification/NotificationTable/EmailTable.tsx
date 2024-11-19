import Table, { TableColumn } from 'components/Table/Table';
import NotificationButtons from 'pages/Notification/NotificationButtons/NotificationButtons';
import NotificationNameLink from 'pages/Notification/NotificationNameLink/NotificationNameLink';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import emailIcon from 'assets/mail-icon.svg';
import { NotificationTableRowProps } from './NotificationTable';
import { useEffect, useState } from 'react';
import { ManagmentServiceApiInstance } from 'api/managment-service';
import LoadingTable from './LoadingTable';
import EmailColumn from 'pages/Notification/EmailCell/EmailCell';
import { NotificationContextProps, useNotification } from '../NotificationContext';

export interface EmailTableRowProps extends NotificationTableRowProps {
  email: string;
}

const emailColumns: Array<TableColumn<EmailTableRowProps>> = [
  {
    header: 'Name',
    columnKey: 'receiverName',
    customComponent: ({ receiverName, destination }: EmailTableRowProps) => (
      <NotificationNameLink linkName={receiverName} destination={destination} />
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
      webhookUrl,
      linkName,
      createdAt,
      updateAt,
    }: EmailTableRowProps) => (
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

const EmailTable = () => {
  const [rows, setRows] = useState<EmailTableRowProps[]>([]);
  const [isLoading, setLoading] = useState<boolean>(true);
  const {
    hidePopup,
    createNewChannel,
  }: NotificationContextProps = useNotification();

  const fetchEmailChannels = async () => {
    try {
      const channels = await ManagmentServiceApiInstance.getEmailChannels();
      setRows(channels);
    } catch (error) {
      console.error('Error fetching email channels: ', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEmailChannels();
  }, []);

  return (
    <SectionComponent
      icon={<img src={emailIcon} />}
      title="Email"
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
          <Table columns={emailColumns} rows={rows} alignLeft={false} />
        ) : (
          <p>No Email channels was yet configured</p>
        )}
      </LoadingTable>
    </SectionComponent>
  );
};

export default EmailTable;
