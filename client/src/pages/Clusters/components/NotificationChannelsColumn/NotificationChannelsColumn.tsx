import { NotificationChannelColumn } from 'pages/Clusters/Clusters';
import './NotificationChannelsColumns.scss';
import NotificationChannelTag from 'components/NotificationChannelTag/NotificationChannelTag';
import slackLogo from 'assets/slack-icon.png';
import discordLogo from 'assets/discord-icon.png';
import mailLogo from 'assets/mail-icon.png';

const NotificationChannelTagFactory = ({
  channel,
}: {
  channel: NotificationChannelColumn;
}) => {
  switch (channel.kind) {
    case 'SLACK':
      return (
        <NotificationChannelTag logoImgSrc={slackLogo} name={channel.name} />
      );
    case 'DISCORD':
      return (
        <NotificationChannelTag logoImgSrc={discordLogo} name={channel.name} />
      );
    case 'EMAIL':
      return (
        <NotificationChannelTag logoImgSrc={mailLogo} name={channel.name} />
      );
  }
};

const NotificationChannelsColumn = ({
  channels,
}: {
  channels: NotificationChannelColumn[];
}) => {
  return (
    <div className="notification-channels-column">
      {channels.map((channel, index) => (
        <NotificationChannelTagFactory channel={channel} key={index} />
      ))}
    </div>
  );
};

export default NotificationChannelsColumn;
