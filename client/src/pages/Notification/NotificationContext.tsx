import { createContext, useContext } from 'react';
import {NotificationChannelKind} from 'api/managment-service.ts';

export type NotificationContextProps = (
    channel: NotificationChannelKind,
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
