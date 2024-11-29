import { ToastMessage } from 'providers/ToastProvider/ToastProvider';
import './Toast.scss';

interface ToastProps {
  message: ToastMessage;
}

const Toast = ({ message }: ToastProps) => {
  return (
    <div className={`toast--${message.type.toLowerCase()}`}>
      <div className="toast__header">{message.type}</div>
      <div className="toast__message">{message.message}</div>
    </div>
  );
};

export default Toast;
