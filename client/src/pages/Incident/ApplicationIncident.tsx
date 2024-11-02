import './Incident.scss';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  ApplicationIncident,
  ManagmentServiceApiInstance,
} from 'api/managment-service';

// eslint-disable-next-line
import ApplicationMetadataSection from './components/ApplicationMetadataSection/ApplicationMetadataSection';
import SummarySection from './components/SummarySection/SummarySection';
import RecommendationSection from './components/RecommendationSection/RecommendationSection';

// eslint-disable-next-line
import ApplicationSourceSection from './components/ApplicationSourceSection/ApplicationSourceSection';
import IncidentHeader from './components/IncidentHeader/IncidentHeader';

const Incident = () => {
  const [incident, setIncident] = useState<ApplicationIncident>();

  const { id } = useParams();

  useEffect(() => {
    const fetchApplicationIncident = async () => {
      try {
        const fetchedIncident =
          await ManagmentServiceApiInstance.getApplicationIncident(id!);
        setIncident(fetchedIncident);
      } catch (err: unknown) {
        // eslint-disable-next-line
        console.error('Failed to fetch application incident');
      }
    };
    fetchApplicationIncident();
  }, [id]);
  if (!incident) {
    return <div></div>;
  }

  return (
    <PageTemplate
      header={
        <IncidentHeader
          id={id!}
          name={incident.applicationName}
          timestamp={incident.sources[0].timestamp}
        />
      }
    >
      <div className="incident">
        <div className="incident__row--two-columns">
          <ApplicationMetadataSection
            clusterId={incident.clusterId}
            applicationName={incident.applicationName}
            startDateMs={1000000}
            endDateMs={1000000}
          />

          <SummarySection summary={incident.summary} />
        </div>
        <div className="incident__row--single-column">
          <RecommendationSection recommendation={incident.recommendation} />
        </div>
        {incident.sources.map((source, index) => (
          <ApplicationSourceSection
            key={index}
            content={source.content}
            container={source.container}
            pod={source.pod}
            image={source.image}
            timestamp={source.timestamp}
          />
        ))}
      </div>
    </PageTemplate>
  );
};

export default Incident;
