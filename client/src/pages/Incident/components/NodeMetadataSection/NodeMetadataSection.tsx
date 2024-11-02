import LabelField from 'components/LabelField/LabelField';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import './NodeMetadataSection.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import { dateFromTimestampMs } from 'lib/date';

interface NodeMetadataSectionParams {
  nodeName: string;
  startDateMs: number;
  endDateMs: number;
}

const NodeMetadataSection = ({
  nodeName,
  startDateMs,
  endDateMs,
}: NodeMetadataSectionParams) => {
  return (
    <SectionComponent
      title={'Node'}
      icon={<SVGIcon iconName={'node-incident-metadata-icon'} />}
    >
      <div className="node-incident-metadata">
        <div className="node-incident-metadata__column">
          <LabelField label={'Node'} field={nodeName} />
        </div>

        <div className="node-incident-metadata__column">
          <LabelField
            label={'Start Date'}
            field={dateFromTimestampMs(startDateMs)}
          />
          <LabelField
            label={'End Date'}
            field={dateFromTimestampMs(endDateMs)}
          />
        </div>
      </div>
    </SectionComponent>
  );
};

export default NodeMetadataSection;
