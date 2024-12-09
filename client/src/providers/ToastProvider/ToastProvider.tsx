import Toast from 'components/Toast/Toast';
import { createContext, useContext, useState } from 'react';
import './ToastProvider.scss';

export interface ToastMessage {
  message: string;
  type: 'INFO' | 'WARNING' | 'ERROR';
}

interface ToastMessageWithId extends ToastMessage {
  id: string;
}

interface ToastContextProps {
  messages: ToastMessage[];
  showMessage: (msg: ToastMessage) => void;
}

export const ToastContext = createContext<ToastContextProps>({
  messages: [],
  showMessage: () => {},
});

const TOAST_MESSAGE_TIMEOUT_MS = 3000;

export const useToast = () => {
  return useContext(ToastContext);
};

export const generateMessageId = () => {
  return Date.now().toString();
};

export const ToastProvider = (props: { children: React.ReactNode }) => {
  const [messages, setMessages] = useState<ToastMessageWithId[]>([]);

  const showMessage = (msg: ToastMessage) => {
    const messageId = generateMessageId();
    setMessages((prev) => [...prev, { ...msg, id: messageId }]);
    setTimeout(() => {
      setMessages((prev) => {
        return [...prev.filter((prevMessage) => prevMessage.id != messageId)];
      });
    }, TOAST_MESSAGE_TIMEOUT_MS);
  };

  return (
    <ToastContext.Provider value={{ messages, showMessage }}>
      {props.children}
      <div className="toasts">
        {messages.map((msg, index) => (
          <div key={index}>
            <Toast message={msg} />
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
};
