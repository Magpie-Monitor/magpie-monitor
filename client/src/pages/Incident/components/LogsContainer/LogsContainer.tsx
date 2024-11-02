import './LogsContainer.scss';

interface LogsBoxParams {
  logs: string;
}

const LogsBox = ({ logs }: LogsBoxParams) => {
  return <div className="logs-box">{logs}</div>;
};

export default LogsBox;
