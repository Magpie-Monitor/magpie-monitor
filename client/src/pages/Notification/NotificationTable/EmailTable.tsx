import Table, { TableColumn } from 'components/Table/Table';
import NotificationButtons from 'pages/Notification/NotificationButtons/NotificationButtons';
import NotificationNameLink from 'pages/Notification/NotificationNameLink/NotificationNameLink';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import emailIcon from 'assets/mail-icon.svg';
import { NotificationTableRowProps } from './NotificationTable';
import { useEffect, useState } from 'react';
import {
  EmailNotificationChannel,
  ManagmentServiceApiInstance,
} from 'api/managment-service';
import LoadingTable from './LoadingTable';
import EmailColumn from 'pages/Notification/EmailCell/EmailCell';
import NewEmailChannelPopup from 'pages/Notification/NewChannelPopup/NewEmailChannelPopup';
import { dateFromTimestampMs } from 'lib/date';
import EditEmailChannelPopup from 'pages/Notification/EditChannelPopup/EditEmailChannelPopup';
import './NotificationTable.scss';
import { useToast } from 'providers/ToastProvider/ToastProvider';

interface EmailTableRowProps extends NotificationTableRowProps {
  email: string;
}

const getEmailChannelTableRow = ({
  id,
  receiverName,
  updatedAt,
  createdAt,
  receiverEmail,
}: EmailNotificationChannel): EmailTableRowProps => ({
  name: receiverName,
  updatedAt: dateFromTimestampMs(updatedAt),
  createdAt: dateFromTimestampMs(createdAt),
  email: receiverEmail,
  id,
});

const EmailTable = () => {
  const [rows, setRows] = useState<EmailTableRowProps[]>([]);
  const [isLoading, setLoading] = useState<boolean>(true);
  const [isNewChannelPopupDisplayed, setIsNewChannelPopupDisplayed] =
    useState<boolean>(false);
  const [isEditChannelPopupDisplayed, setIsEditChannelPopupDisplayed] =
    useState<boolean>(false);
  const [editChannelPopupData, setEditChannelPopupData] =
    useState<EmailTableRowProps | null>(null);

  const { showMessage } = useToast();

  const fetchEmailChannels = async () => {
    try {
      const channels = await ManagmentServiceApiInstance.getEmailChannels();
      setRows(channels.map(getEmailChannelTableRow));
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Error fetching email channels: ', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEmailChannels();
  }, [isLoading]);

  const emailColumns: Array<TableColumn<EmailTableRowProps>> = [
    {
      header: 'Name',
      columnKey: 'receiverName',
      customComponent: ({ name }: EmailTableRowProps) => (
        <NotificationNameLink linkName={name} />
      ),
    },
    {
      header: 'Email',
      columnKey: 'receiverEmail',
      customComponent: ({ email }: EmailTableRowProps) => (
        <EmailColumn email={email} />
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
      customComponent: (props) => (
        <NotificationButtons
          onUpdate={() => {
            setIsEditChannelPopupDisplayed(true);
            setEditChannelPopupData(props);
          }}
          onTest={async () => {
            try {
              await ManagmentServiceApiInstance.testEmailChannel(props.id);
              showMessage({
                message: 'Test notification was sent',
                type: 'INFO',
              });
            } catch (e: unknown) {
              showMessage({
                message: 'Failed to send test notification',
                type: 'ERROR',
              });
            }
          }}
          onDelete={async () => {
            try {
              await ManagmentServiceApiInstance.deleteEmailChannel(props.id);
              showMessage({
                message: 'Email channel was deleted',
                type: 'INFO',
              });
            } catch (e: unknown) {
              showMessage({
                message: `Failed to delete email channel: ${e}`,
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
      icon={<img src={emailIcon} className="notification-table__icon" />}
      title="Email"
      callback={() => {
        setIsNewChannelPopupDisplayed(true);
      }}
    >
      <LoadingTable isLoading={isLoading}>
        {rows.length > 0 ? (
          <Table columns={emailColumns} rows={rows} alignLeft={false} />
        ) : (
          <p>No Email channels was yet configured</p>
        )}
      </LoadingTable>

      {isEditChannelPopupDisplayed && editChannelPopupData && (
        <EditEmailChannelPopup
          id={editChannelPopupData.id}
          name={editChannelPopupData.name}
          email={editChannelPopupData.email}
          isDisplayed={isEditChannelPopupDisplayed}
          setIsDisplayed={setIsEditChannelPopupDisplayed}
          onSubmit={() => setLoading(true)}
        />
      )}
      {isNewChannelPopupDisplayed && (
        <NewEmailChannelPopup
          setIsDisplayed={setIsNewChannelPopupDisplayed}
          isDisplayed={isNewChannelPopupDisplayed}
          onSubmit={() => setLoading(true)}
        />
      )}
    </SectionComponent>
  );
};

export default EmailTable;
