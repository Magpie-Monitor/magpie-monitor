import './NotificationButtons.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon';

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
      <div onClick={onUpdate}>
        <SVGIcon iconName="edit-icon" />
      </div>
      <div onClick={onTest}>
        <SVGIcon iconName="send-notification-icon" />
      </div>
      <div onClick={onDelete}>
        <SVGIcon iconName="delete-icon" />
      </div>
    </div>
  );
};

export default NotificationButtons;
