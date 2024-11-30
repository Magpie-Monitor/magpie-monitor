import './Incident.scss';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import { useCallback, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  ApplicationIncident,
  ApplicationIncidentSource,
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
import InfiniteScroll from 'react-infinite-scroll-component';
import { useTransition, animated } from '@react-spring/web';

const APPLICATION_SOURCE_PAGE_SIZE = 5;
const PAGE_TEMPLATE_INFINITE_SCROLL_ID = 'application-incident-page';

const ApplicationIncidentPage = () => {
  const [incident, setIncident] = useState<ApplicationIncident>();
  const [isLoading, setIsLoading] = useState(true);
  const [sourcesPage, setSourcesPage] = useState(0);
  const [sources, setSources] = useState<ApplicationIncidentSource[]>([]);
  const [allSourcesCount, setAllSourcesCount] = useState(-1);
  const { id } = useParams();

  const transitions = useTransition(sources, {
    from: { opacity: 0, transform: 'translateY(20px)' },
    enter: { opacity: 1, transform: 'translateY(0)' },
    leave: { opacity: 0, transform: 'translateY(20px)' },
    config: { duration: 200 },
    trail: 400,
  });

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

  const fetchSources = useCallback(async () => {
    const newSources =
      await ManagmentServiceApiInstance.getApplicationIncidentSources(
        incident!.id,
        sourcesPage,
        APPLICATION_SOURCE_PAGE_SIZE,
      );

    if (newSources.data.length === 0) {
      return;
    }

    setAllSourcesCount(newSources.totalEntries);
    setSources((prev) => [...prev, ...newSources.data]);
    setSourcesPage((page) => page + 1);
  }, [incident, sourcesPage]);

  useEffect(() => {
    fetchSources();
  }, [incident, fetchSources]);

  if (isLoading || !incident) {
    return (
      <PageTemplate header={''}>
        <Spinner />
      </PageTemplate>
    );
  }

  const [startDate, endDate] = getFirstAndLastDateFromTimestamps(
    incident.sources.map(({ timestamp }) => timestamp),
  );

  return (
    <PageTemplate
      header={
        <IncidentHeader id={id!} name={incident.title} timestamp={startDate} />
      }
      id={PAGE_TEMPLATE_INFINITE_SCROLL_ID}
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
            <SummarySection summary={incident.summary} />
            <RecommendationSection recommendation={incident.recommendation} />
          </div>
        </div>

        <InfiniteScroll
          dataLength={sources.length}
          next={fetchSources}
          hasMore={sources.length < allSourcesCount}
          loader={<Spinner />}
          scrollableTarget={PAGE_TEMPLATE_INFINITE_SCROLL_ID}
          className="incident"
        >
          {transitions((style, source) => (
            <animated.div style={style}>
              <ApplicationSourceSection
                content={source.content}
                container={source.containerName}
                pod={source.podName}
                image={source.image}
                timestamp={source.timestamp}
              />
            </animated.div>
          ))}
        </InfiniteScroll>
      </div>
    </PageTemplate>
  );
};

export default ApplicationIncidentPage;
