import './NotificationChannelColumn.scss';
import { NotificationChannelColumn } from '@pages/Report/NotificationSection/NotificationUtils.tsx';
import NotificationChannelTag from 'components/NotificationChannelTag/NotificationChannelTag';
import slackLogo from 'assets/slack-icon.png';
import discordLogo from 'assets/discord-icon.png';
import mailLogo from 'assets/mail-icon.svg';

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
    default:
      return null;
  }
};

const NotificationChannelDisplay = ({
                                      channel,
                                    }: {
  channel: NotificationChannelColumn;
}) => {
  return (
      <div className="notification-channels-column">
        <NotificationChannelTagFactory channel={channel} />
      </div>
  );
};

export default NotificationChannelDisplay;
