import { createContext, useContext } from 'react';

export type NotificationsChannel = 'SLACK' | 'DISCORD' | 'EMAIL';

export type NotificationContextProps = (
    channel: NotificationsChannel,
    adress: string,
    linkName: string,
    destination: string,
    createdAt: string,
    updateAt: string,
  ) => void;

export const NotificationContext = createContext<NotificationContextProps>(
  ()=>{}
);

export const useNotification = () => {
  return useContext(NotificationContext);
};
