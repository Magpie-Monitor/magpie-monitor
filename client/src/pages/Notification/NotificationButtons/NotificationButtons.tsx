import ActionButton, {
  ActionButtonColor,
} from 'components/ActionButton/ActionButton';
import './NotificationButtons.scss';

export interface NotificationButtonsProps {
  onUpdate: () => void;
  onTest: () => void;
  onDelete: () => void;
}

const NotificationButtons = ({
  onTest,
  onUpdate,
  onDelete,
}: NotificationButtonsProps) => {
  return (
    <div className="notification-buttons">
      <ActionButton
        onClick={onUpdate}
        description="UPDATE"
        color={ActionButtonColor.GREEN}
      />
      <ActionButton
        onClick={onTest}
        description="TEST"
        color={ActionButtonColor.OLIVE}
      />
      <ActionButton
        onClick={onDelete}
        description="DELETE"
        color={ActionButtonColor.RED}
      />
    </div>
  );
};

export default NotificationButtons;
