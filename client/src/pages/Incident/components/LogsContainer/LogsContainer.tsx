import './LogsContainer.scss';
import { useState } from 'react';

interface LogsBoxParams {
  logs: string;
}

const LogsBox = ({ logs }: LogsBoxParams) => {
  const [isLogBoxExpanded, setIsLogBoxExpanded] = useState(false);

  return (
    <>
      <div
        className={isLogBoxExpanded ? 'logs-box--expanded' : 'logs-box'}
        onClick={() => setIsLogBoxExpanded((prev) => !prev)}
      >
        {logs.trim()}
      </div>
    </>
  );
};

export default LogsBox;
