import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import './ReportHeader.scss';
import { dateOnlyFromTimestampMs } from 'lib/date';

interface ReportHeaderProps {
  name: string;
  sinceMs: number;
  toMs: number;
}

const ReportHeader = ({ name, sinceMs, toMs }: ReportHeaderProps) => {
  const title = (
    <div className="report-header">
      <div className="report-header__name">{name}</div>
      <div className="report-header__timestamp">
        ({dateOnlyFromTimestampMs(sinceMs)} - {dateOnlyFromTimestampMs(toMs)})
      </div>
    </div>
  );

  return (
    <HeaderWithIcon
      title={title}
      icon={<SVGIcon iconName="report-details-icon" />}
    />
  );
};

export default ReportHeader;
