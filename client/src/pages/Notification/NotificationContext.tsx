import { createContext, useContext } from 'react';
import { NotificationChannelKind } from 'api/managment-service.ts';

export interface NotificationContextProps {
  isPopupDisplayed: boolean;
  hidePopup: () => void;
  updater: (
    channel: NotificationsChannelKind,
    adress: string,
    linkName: string,
    createdAt: string,
    updateAt: string,
  ) => void;
  createNewChannel: (Popup: React.ReactNode) => void;
}

const defaultNotificationContextProps: NotificationContextProps = {
  isPopupDisplayed: false,
  hidePopup: (): void => {
    throw new Error('Function not implemented.');
  },
  updater: (
    channel: NotificationsChannel,
    adress: string,
    linkName: string,
    createdAt: string,
    updateAt: string,
  ): void => {
    throw new Error('Function not implemented.');
  },
  createNewChannel: (Popup: React.ReactNode): void => {
    throw new Error('Function not implemented.');
  },
};

export const NotificationContext = createContext<NotificationContextProps>(
  defaultNotificationContextProps,
);

export const useNotification = () => {
  return useContext(NotificationContext);
};
