import './Notification.scss';
import { NotificationContext } from 'pages/Notification/NotificationContext';
import SlackTable from './NotificationTable/SlackTable';
import DiscordTable from './NotificationTable/DiscordTable';
import EmailTable from './NotificationTable/EmailTable';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import { useState } from 'react';

const Notification = () => {
  const [isPopupShowed, setIsPopupShowed] = useState<boolean>(false);
  const [currentPopup, setCurrentPopup] = useState<React.ReactNode>();

  const addNewChannel = (popup: React.ReactNode) => {
    if (isPopupShowed) return;

    setIsPopupShowed(true);
    setCurrentPopup(popup);
  };

  return (
    <NotificationContext.Provider
      value={{
        isPopupDisplayed: isPopupShowed,
        hidePopup: () => setIsPopupShowed(false),
        updater: () => {},
        createNewChannel: (popup: React.ReactNode) => addNewChannel(popup),
      }}
    >
      <PageTemplate
        header={
          <HeaderWithIcon
            icon={<SVGIcon iconName="notification-icon" />}
            title="Notification channels"
          />
        }
      >
        <div className="notification__body">
          <SlackTable />
          <DiscordTable />
          <EmailTable />
        </div>
        {isPopupShowed && currentPopup}
      </PageTemplate>
    </NotificationContext.Provider>
  );
};

export default Notification;
