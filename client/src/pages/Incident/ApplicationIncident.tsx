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
import { getFirstAndLastDateFromTimestamps } from 'lib/date';
import Spinner from 'components/Spinner/Spinner';
import ConfigurationSection from './components/ConfigurationSection/ConfigurationSection';

const ApplicationIncidentPage = () => {
  const [incident, setIncident] = useState<ApplicationIncident>();
  const [isLoading, setIsLoading] = useState(true);

  const { id } = useParams();

  useEffect(() => {
    const fetchApplicationIncident = async () => {
      try {
        const fetchedIncident =
          await ManagmentServiceApiInstance.getApplicationIncident(id!);
        setIsLoading(false);
        setIncident(fetchedIncident);
      } catch (err: unknown) {
        // eslint-disable-next-line
        console.error('Failed to fetch application incident');
      }
    };
    fetchApplicationIncident();
  }, [id]);

  if (isLoading || !incident) {
    return <PageTemplate header={''}> <Spinner /> </PageTemplate>;
  }

  const [startDate, endDate] = getFirstAndLastDateFromTimestamps(
    incident.sources.map(({ timestamp }) => timestamp),
  );

  return (
    <PageTemplate
      header={
        <IncidentHeader id={id!} name={incident.title} timestamp={startDate} />
      }
    >
      <div className="incident">
        <div>
          <div className="incident__row--two-columns">
            <ApplicationMetadataSection
                clusterId={incident.clusterId}
                applicationName={incident.applicationName}
                startDateMs={startDate}
                endDateMs={endDate}
            />

            <ConfigurationSection
                accuracy={incident.accuracy}
                customPrompt={incident.customPrompt}
            />
          </div>
        </div>
        <div>
          <div className="incident__row--two-columns">
            <SummarySection summary={incident.summary}/>
            <RecommendationSection recommendation={incident.recommendation}/>
          </div>
        </div>
        {incident.sources.map((source, index) => (
          <ApplicationSourceSection
            content={source.content}
            key={index}
            container={source.containerName}
            pod={source.podName}
            image={source.image}
            timestamp={source.timestamp}
          />
        ))}
      </div>
    </PageTemplate>
);
};

export default ApplicationIncidentPage;
