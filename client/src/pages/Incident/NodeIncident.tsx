import './Incident.scss';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  ManagmentServiceApiInstance,
  NodeIncident,
  NodeIncidentSource,
} from 'api/managment-service';
import SummarySection from './components/SummarySection/SummarySection';
import RecommendationSection from './components/RecommendationSection/RecommendationSection';
import IncidentHeader from './components/IncidentHeader/IncidentHeader';
import NodeMetadataSection from './components/NodeMetadataSection/NodeMetadataSection';
import NodeSourceSection from './components/NodeSourceSection/NodeSourceSection';
import { getFirstAndLastDateFromTimestamps } from 'lib/date';
import ConfigurationSection from './components/ConfigurationSection/ConfigurationSection';
import { animated, useTransition } from '@react-spring/web';
import usePaginatedContent from 'hooks/usePaginatedContent';
import { FadeInTransition } from 'hooks/TransitionParams';
import useInfiniteScroll from 'hooks/useInfiniteScroll';
import CenteredSpinner from 'components/CenteredSpinner/CenteredSpinner';

const NODE_SOURCE_PAGE_SIZE = 5;

const NodeIncidentPage = () => {
  const [incident, setIncident] = useState<NodeIncident>();
  const [isLoading, setIsLoading] = useState(true);
  const { id } = useParams();
  const {
    content,
    contentPage,
    setTotalContentCount,
    addContent,
    isAllContentFetched,
  } = usePaginatedContent<NodeIncidentSource>();

  const [isFetchingSources, setIsFetchingSources] = useState(true);

  const transitions = useTransition(content, FadeInTransition);
  const pageTemplateRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchNodeIncident = async () => {
      try {
        const fetchedIncident =
          await ManagmentServiceApiInstance.getNodeIncident(id!);

        setIsLoading(false);
        setIncident(fetchedIncident);
      } catch (err: unknown) {
        // eslint-disable-next-line
        console.error('Failed to fetch application incident');
      }
    };
    fetchNodeIncident();
  }, [id]);

  const handleScroll = async () => {
    if (isAllContentFetched()) {
      return;
    }

    if (!incident) {
      return;
    }

    setIsFetchingSources(true);
    const newSources = await ManagmentServiceApiInstance.getNodeIncidentSources(
      incident!.id,
      contentPage,
      NODE_SOURCE_PAGE_SIZE,
    );

    addContent(newSources.data);
    setTotalContentCount(newSources.totalEntries);
    setIsFetchingSources(false);
  };

  useInfiniteScroll({ handleScroll, scrollTargetRef: pageTemplateRef });

  useEffect(() => {
    const fetchSources = async () => {
      if (!incident) {
        return;
      }
      try {
        setIsFetchingSources(true);
        const newSources =
          await ManagmentServiceApiInstance.getNodeIncidentSources(
            incident!.id,
            contentPage,
            NODE_SOURCE_PAGE_SIZE,
          );

        setIsFetchingSources(false);

        addContent(newSources.data);
        setTotalContentCount(newSources.totalEntries);
      } catch (err) {
        // eslint-disable-next-line
        console.error('Failed to fetch sources');
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
            <NodeMetadataSection
              nodeName={incident.nodeName}
              startDateMs={startDate}
              endDateMs={endDate}
            />

            <ConfigurationSection
              accuracy={incident.accuracy}
              customPrompt={incident.customPrompt}
            />
          </div>
        </div>
        <div className="incident__row--two-columns">
          <SummarySection summary={incident.summary} />
          <RecommendationSection recommendation={incident.recommendation} />
        </div>

        {transitions((style, source) => (
          <animated.div style={style}>
            <NodeSourceSection
              content={source.content}
              filename={source.filename}
              timestamp={source.timestamp}
            />
          </animated.div>
        ))}

        {isFetchingSources && <CenteredSpinner />}
      </div>
    </PageTemplate>
  );
};

export default NodeIncidentPage;
