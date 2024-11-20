import './Notification.scss';
import SlackTable from './NotificationTable/SlackTable';
import DiscordTable from './NotificationTable/DiscordTable';
import EmailTable from './NotificationTable/EmailTable';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import SVGIcon from 'components/SVGIcon/SVGIcon';

const Notification = () => {
  return (
    <PageTemplate
      header={
        <HeaderWithIcon
          icon={<SVGIcon iconName="big-notification-icon" />}
          title="Notification channels"
        />
      }
    >
      <div className="notification__body">
        <SlackTable />
        <DiscordTable />
        <EmailTable />
      </div>
    </PageTemplate>
  );
};

export default Notification;
