import LabelField from 'components/LabelField/LabelField';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import './ApplicationMetadataSection.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon';

interface ApplicationMetadataSectionParams {
  clusterId: string;
  applicationName: string;
  startDateMs: number;
  endDateMs: number;
}

const ApplicationMetadataSection = ({
  clusterId,
  applicationName,
  startDateMs,
  endDateMs,
}: ApplicationMetadataSectionParams) => {
  return (
    <SectionComponent
      title={'Application'}
      icon={<SVGIcon iconName={'application-incident-metadata-icon'} />}
    >
      <div className="application-incident-metadata">
        <div className="application-incident-metadata__column">
          <LabelField label={'Cluster'} field={clusterId} />
          <LabelField label={'Application'} field={applicationName} />
        </div>

        <div className="application-incident-metadata__column">
          <LabelField
            label={'Start Date'}
            field={new Date(startDateMs).toLocaleString()}
          />
          <LabelField
            label={'End Date'}
            field={new Date(endDateMs).toLocaleString()}
          />
        </div>
      </div>
    </SectionComponent>
  );
};

export default ApplicationMetadataSection;
