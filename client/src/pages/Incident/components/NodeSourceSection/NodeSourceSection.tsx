import LabelField from 'components/LabelField/LabelField';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import LogsBox from 'pages/Incident/components/LogsContainer/LogsContainer';
import './NodeSourceSection.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon';

interface NodeSourceParams {
  content: string;
  filename: string;
  timestamp: number;
}

const NodeSourceSection = ({
  content,
  filename,
  timestamp,
}: NodeSourceParams) => {
  return (
    <SectionComponent
      title={'Source'}
      icon={<SVGIcon iconName="incident-source-icon" />}
    >
      <div className="node-incident-source">
        <div className="node-incident-source__metadata">
          <LabelField field={filename} label="Filename" />
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

export default NodeSourceSection;
