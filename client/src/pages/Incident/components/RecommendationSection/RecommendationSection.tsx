import SVGIcon from 'components/SVGIcon/SVGIcon';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import './RecommendationSection.scss';

interface RecommendationSectionParams {
  recommendation: string;
}

const RecommendationSection = ({
  recommendation,
}: RecommendationSectionParams) => {
  return (
    <SectionComponent
      title={'Recommendation'}
      icon={<SVGIcon iconName="incident-recommendation-icon" />}
    >
      <div className="incident-recommendation">{recommendation}</div>
    </SectionComponent>
  );
};

export default RecommendationSection;
