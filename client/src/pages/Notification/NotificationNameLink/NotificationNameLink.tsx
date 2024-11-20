import { Link } from 'react-router-dom';
import './NotificationNameLink.scss';

export interface NotificationNameLinkProps {
  linkName: string;
  destination: string;
}

const NotificationNameLink = ({
  linkName,
  destination,
}: NotificationNameLinkProps) => {
  return (
    <Link className="notification-link-name" to={destination}>
      {linkName}
    </Link>
  );
};

export default NotificationNameLink;
