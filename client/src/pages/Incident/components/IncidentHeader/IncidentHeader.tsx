import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import './IncidentHeader.scss';
import { dateFromTimestampMs } from 'lib/date';

interface IncidentHeaderProps {
  id: string;
  name: string;
  timestamp: number;
}

const IncidentHeader = ({ name, timestamp }: IncidentHeaderProps) => {
  const title = (
    <div className="incident-header">
      <div className="incident-header__name">{name}</div>
      <div className="incident-header__timestamp">
        {dateFromTimestampMs(timestamp)}
      </div>
    </div>
  );

  return (
    <HeaderWithIcon title={title} icon={<SVGIcon iconName="incident-icon" />} />
  );
};

export default IncidentHeader;
