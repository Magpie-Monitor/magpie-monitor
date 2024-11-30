import './Incident.scss';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import { useCallback, useEffect, useState } from 'react';
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
import Spinner from 'components/Spinner/Spinner';
import ConfigurationSection from './components/ConfigurationSection/ConfigurationSection';
import InfiniteScroll from 'react-infinite-scroll-component';
import { animated, useTransition } from '@react-spring/web';

const NODE_SOURCE_PAGE_SIZE = 5;
const PAGE_TEMPLATE_INFINITE_SCROLL_ID = 'node-incident-page';

const NodeIncidentPage = () => {
  const [incident, setIncident] = useState<NodeIncident>();
  const [isLoading, setIsLoading] = useState(true);
  const [sourcesPage, setSourcesPage] = useState(0);
  const [sources, setSources] = useState<NodeIncidentSource[]>([]);
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

  const fetchSources = useCallback(async () => {
    const newSources = await ManagmentServiceApiInstance.getNodeIncidentSources(
      incident!.id,
      sourcesPage,
      NODE_SOURCE_PAGE_SIZE,
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
        {' '}
        <Spinner />{' '}
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
        <div className="incident__sources">
          <div className="incident__row--two-columns">
            <SummarySection summary={incident.summary} />
            <RecommendationSection recommendation={incident.recommendation} />
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
                <NodeSourceSection
                  content={source.content}
                  timestamp={source.timestamp}
                  filename={source.filename}
                />
              </animated.div>
            ))}
          </InfiniteScroll>
        </div>
      </div>
    </PageTemplate>
  );
};

export default NodeIncidentPage;
