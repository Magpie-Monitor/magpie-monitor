import ActionButton, {
  ActionButtonColor,
} from 'components/ActionButton/ActionButton';
import './NotificationButtons.scss';
import { NotificationsChannel, useNotification } from 'pages/Notification/NotificationContext';

export interface NotificationButtonsProps {
  channel: NotificationsChannel;
  adress: string;
  linkName: string;
  destination: string;
  createdAt: string;
  updateAt: string;
}

const NotificationButtons = ({
  channel,
  adress,
  linkName,
  destination,
  createdAt,
  updateAt,
}: NotificationButtonsProps) => {
  const updater = useNotification();

  return (
    <div className="notification-buttons">
      <ActionButton
        onClick={() => updater(channel, adress, linkName, destination, createdAt, updateAt)}
        description="UPDATE"
        color={ActionButtonColor.GREEN}
      />
      <ActionButton
        onClick={() => {}}
        description="TEST"
        color={ActionButtonColor.OLIVE}
      />
      <ActionButton
        onClick={() => {}}
        description="DELETE"
        color={ActionButtonColor.RED}
      />
    </div>
  );
};

export default NotificationButtons;
