import './Incident.scss';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import { useEffect, useRef, useState } from 'react';
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
import ConfigurationSection from './components/ConfigurationSection/ConfigurationSection';
import { useTransition, animated } from '@react-spring/web';
import useInfiniteScroll from 'hooks/useInfiniteScroll';
import usePaginatedContent from 'hooks/usePaginatedContent';
import { FadeInTransition } from 'hooks/TransitionParams';
import CenteredSpinner from 'components/CenteredSpinner/CenteredSpinner';

const APPLICATION_SOURCE_PAGE_SIZE = 5;

const ApplicationIncidentPage = () => {
  const [incident, setIncident] = useState<ApplicationIncident>();
  const [isLoading, setIsLoading] = useState(true);
  const { id } = useParams();
  const {
    content,
    contentPage,
    setTotalContentCount,
    addContent,
    isAllContentFetched,
  } = usePaginatedContent<ApplicationIncidentSource>();
  const [isFetchingSources, setIsFetchingSources] = useState(true);

  const pageTemplateRef = useRef<HTMLDivElement>(null);

  const transitions = useTransition(content, FadeInTransition);

  const handleScroll = async () => {
    if (isAllContentFetched()) {
      return;
    }

    if (!incident) {
      return;
    }

    setIsFetchingSources(true);

    const newSources =
      await ManagmentServiceApiInstance.getApplicationIncidentSources(
        incident!.id,
        contentPage,
        APPLICATION_SOURCE_PAGE_SIZE,
      );

    setIsFetchingSources(false);

    addContent(newSources.data);
    setTotalContentCount(newSources.totalEntries);
  };

  useInfiniteScroll({ handleScroll, scrollTargetRef: pageTemplateRef });

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

  useEffect(() => {
    const fetchSources = async () => {
      if (!incident) {
        return;
      }
      try {
        setIsFetchingSources(true);
        const newSources =
          await ManagmentServiceApiInstance.getApplicationIncidentSources(
            incident!.id,
            contentPage,
            APPLICATION_SOURCE_PAGE_SIZE,
          );

        addContent(newSources.data);
        setTotalContentCount(newSources.totalEntries);

        setIsFetchingSources(false);
      } catch (err) {
        // eslint-disable-next-line
        console.error(`Failed to fetch sources ${err}`);
      }
    };
    fetchSources();
    // eslint-disable-next-line
  }, [incident]);

  if (isLoading || !incident) {
    return (
      <PageTemplate header={''}>
        <CenteredSpinner />
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
      scrollRef={pageTemplateRef}
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
        {isFetchingSources && <CenteredSpinner />}
      </div>
    </PageTemplate>
  );
};

export default ApplicationIncidentPage;
