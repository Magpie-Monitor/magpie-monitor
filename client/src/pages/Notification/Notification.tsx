import NotificationTable, {
  EmailTableRowProps,
  WebhookTableRowProps,
} from 'pages/Notification/NotificationTable/NotificationTable';
import './Notification.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import {
  NotificationContext,
} from 'pages/Notification/NotificationContext';

const slackNotificationTableRow: WebhookTableRowProps[] = [
  {
    linkName: 'infra_team_slack',
    webhookUrl: 'https://slack.com/api/webhooks/example',
    createdAt: '07.03.2024 21:37',
    updateAt: '14.03.2024 00:00',
    action: 'here will be button',
    destination: '/notification',
  },
  {
    linkName: 'dev team slack',
    webhookUrl: 'https://slack.com/api/webhooks/pdoasds',
    createdAt: '14.03.2024 21:37',
    updateAt: '21.03.2024 00:00',
    action: 'here will be button',
    destination: '/notification',
  },
];

const mailNotificationTableRow: EmailTableRowProps[] = [
  {
    linkName: 'infra_team_slack',
    createdAt: '07.03.2024 21:37',
    updateAt: '14.03.2024 00:00',
    action: 'here will be button',
    destination: '/notification',
    email: 'kontakt@wmsdev.pl',
  },
  {
    linkName: 'dev team slack',
    createdAt: '14.03.2024 21:37',
    updateAt: '21.03.2024 00:00',
    action: 'here will be button',
    destination: '/notification',
    email: 'kontakt@wmsdev.pl',
  },
];

const Notification = () => {
  const showUpdate = () => {};

  return (
    <NotificationContext.Provider value={showUpdate}>
      <div className="notification">
        <div>
          <div className="notification__header">
            <SVGIcon iconName="notification-icon" />
            <p className="notification__header__description">
              Notification channels
            </p>
          </div>
          <div className="notification__body">
            <NotificationTable
              data={slackNotificationTableRow}
              imageName="slack-icon.png"
              header="Slack"
              channel={'SLACK'}
            />
            <NotificationTable
              data={slackNotificationTableRow}
              imageName="discord-icon.png"
              header="Discord"
              channel={'DISCORD'}
            />
            <NotificationTable
              data={mailNotificationTableRow}
              imageName="mail-icon.png"
              header="Email"
              channel={'EMAIL'}
            />
          </div>
        </div>
      </div>
    </NotificationContext.Provider>
  );
};

export default Notification;
