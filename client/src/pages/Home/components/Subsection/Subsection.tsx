import 'pages/ReportDetails/ReportDetails.scss';
import './Subsection.scss';

interface ReportDetailsSubsectionProps {
  title: string;
  children: React.ReactNode;
}

const ReportDetailsSubsection = ({
  children,
  title,
}: ReportDetailsSubsectionProps) => {
  return (
    <div className="report-details__subsection">
      <div className="report-details__subsection__title">{title}</div>
      {children}
    </div>
  );
};

export default ReportDetailsSubsection;
