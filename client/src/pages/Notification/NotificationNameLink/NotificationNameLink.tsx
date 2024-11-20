import { Link } from 'react-router-dom';
import './NotificationNameLink.scss';

export interface NotificationNameLinkProps {
  linkName: string;
}

const NotificationNameLink = ({ linkName }: NotificationNameLinkProps) => {
  return (
    <Link className="notification-link-name" to={''}>
      {linkName}
    </Link>
  );
};

export default NotificationNameLink;
