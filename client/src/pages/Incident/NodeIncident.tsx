import './Incident.scss';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  ManagmentServiceApiInstance,
  NodeIncident,
} from 'api/managment-service';

// eslint-disable-next-line
import SummarySection from './components/SummarySection/SummarySection';
import RecommendationSection from './components/RecommendationSection/RecommendationSection';

// eslint-disable-next-line
import IncidentHeader from './components/IncidentHeader/IncidentHeader';
import NodeMetadataSection from './components/NodeMetadataSection/NodeMetadataSection';
import NodeSourceSection from './components/NodeSourceSection/NodeSourceSection';
import { getFirstAndLastDateFromTimestamps } from 'lib/date';

const NodeIncidentPage = () => {
  const [incident, setIncident] = useState<NodeIncident>();

  const { id } = useParams();

  useEffect(() => {
    const fetchNodeIncident = async () => {
      try {
        const fetchedIncident =
          await ManagmentServiceApiInstance.getNodeIncident(id!);
        setIncident(fetchedIncident);
      } catch (err: unknown) {
        // eslint-disable-next-line
        console.error('Failed to fetch application incident');
      }
    };
    fetchNodeIncident();
  }, [id]);
  if (!incident) {
    return <div></div>;
  }

  const [startDate, endDate] = getFirstAndLastDateFromTimestamps(
    incident.sources.map(({ timestamp }) => timestamp),
  );

  return (
    <PageTemplate
      header={
        <IncidentHeader
          id={id!}
          name={incident.nodeName}
          timestamp={incident.sources[0].timestamp}
        />
      }
    >
      <div className="incident">
        <div className="incident__row--two-columns">
          <NodeMetadataSection
            nodeName={incident.nodeName}
            startDateMs={startDate}
            endDateMs={endDate}
          />

          <SummarySection summary={incident.summary} />
        </div>
        <div className="incident__row--single-column">
          <RecommendationSection recommendation={incident.recommendation} />
        </div>
        {incident.sources.map((source, index) => (
          <NodeSourceSection
            content={source.content}
            key={index}
            nodeName={source.nodeName}
            timestamp={source.timestamp}
          />
        ))}
      </div>
    </PageTemplate>
  );
};

export default NodeIncidentPage;
