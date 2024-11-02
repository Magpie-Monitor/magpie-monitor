import SVGIcon from 'components/SVGIcon/SVGIcon';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import './SummarySection.scss';

interface SummarySectionParams {
  summary: string;
}

const SummarySection = ({ summary }: SummarySectionParams) => {
  return (
    <SectionComponent
      title={'Summary'}
      icon={<SVGIcon iconName="incident-summary-icon" />}
    >
      <div className="incident-summary">{summary}</div>
    </SectionComponent>
  );
};

export default SummarySection;
