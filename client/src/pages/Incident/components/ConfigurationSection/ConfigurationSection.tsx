import SVGIcon from 'components/SVGIcon/SVGIcon';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import './ConfigurationSection.scss';
import { AccuracyLevel } from 'api/managment-service';
import LabelField from 'components/LabelField/LabelField';

interface ConfigurationSectionParams {
  accuracy: AccuracyLevel;
  customPrompt: string;
}

const ConfigurationSection = ({
  accuracy,
  customPrompt,
}: ConfigurationSectionParams) => {
  return (
    <SectionComponent
      title={'Configuration'}
      icon={<SVGIcon iconName="incident-configuration-icon" />}
    >
      <div className="incident-configuration">
        <LabelField label="Accuracy" field={accuracy} />
        <LabelField label="Custom prompt" field={customPrompt} />
      </div>
    </SectionComponent>
  );
};

export default ConfigurationSection;
