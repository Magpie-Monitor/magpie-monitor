import LabelField from 'components/LabelField/LabelField';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import LogsBox from 'pages/Incident/components/LogsContainer/LogsContainer';
import './ApplicationSourceSection.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon';

interface ApplicationSourceParams {
  pod: string;
  container: string;
  image: string;
  content: string;
  timestamp: number;
}

const ApplicationSourceSection = ({
  pod,
  container,
  image,
  content,
  timestamp,
}: ApplicationSourceParams) => {
  return (
    <SectionComponent
      title={'Source'}
      icon={<SVGIcon iconName="incident-source-icon" />}
    >
      <div className="application-incident-source">
        <div className="application-incident-source__metadata">
          <LabelField field={pod} label="Pod" />
          <LabelField field={container} label="Container" />
          <LabelField field={image} label="Image" />
          <LabelField
            field={new Date(timestamp).toLocaleString()}
            label="Timestamp"
          />
        </div>
        <LogsBox logs={content} />
      </div>
    </SectionComponent>
  );
};

export default ApplicationSourceSection;
