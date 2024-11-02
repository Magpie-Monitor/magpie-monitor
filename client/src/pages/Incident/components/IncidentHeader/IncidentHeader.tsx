import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import './IncidentHeader.scss';

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
        {new Date(timestamp).toLocaleString()}
      </div>
    </div>
  );

  return (
    <HeaderWithIcon title={title} icon={<SVGIcon iconName="incident-icon" />} />
  );
};

export default IncidentHeader;
