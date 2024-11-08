import SectionComponent from 'components/SectionComponent/SectionComponent';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import Table, { TableColumn } from 'components/Table/Table';
import './Clusters.scss';
import Channels from './components/NotificationChannelsColumn/NotificationChannelsColumn';
import { useEffect, useState } from 'react';
import {
  ClusterSummary,
  ManagmentServiceApiInstance,
  AccuracyLevel,
} from 'api/managment-service';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import Spinner from 'components/Spinner/Spinner.tsx';
import ReportActionsCell from './ReportActionsCell';
import AccuracyBadge from 'components/AccuracyBadge/AccuracyBadge.tsx';

interface ClusterDataRow {
  name: string;
  state: 'ONLINE' | 'OFFLINE';
  accuracy: AccuracyLevel;
  notificationChannels: NotificationChannelColumn[];
  updatedAt: string;
  [key: string]: string | NotificationChannelColumn[];
}

type NotificationChannelKind = 'SLACK' | 'DISCORD' | 'EMAIL';

export interface NotificationChannelColumn {
  kind: NotificationChannelKind;
  name: string;
}

const transformNotificationChannelsToColumns = (
  cluster: ClusterSummary,
): NotificationChannelColumn[] => {
  return cluster.slackChannels
    .map(
      (channel): NotificationChannelColumn => ({
        kind: 'SLACK',
        name: channel.name,
      }),
    )
    .concat(
      cluster.discordChannels.map((channel) => ({
        kind: 'DISCORD',
        name: channel.name,
      })),
    )
    .concat(
      cluster.mailChannels.map((channel) => ({
        kind: 'EMAIL',
        name: channel.email,
      })),
    );
};
const transformIsRunningLabel = (
  cluster: ClusterSummary,
): ClusterDataRow['state'] => {
  return cluster.running? 'ONLINE' : 'OFFLINE';
};

const transformUpdatedAtDate = (cluster: ClusterSummary) => {
  const date = new Date(cluster.updatedAt);
  return date.toLocaleString();
};

const columns: Array<TableColumn<ClusterDataRow>> = [
  {
    header: 'Name',
    columnKey: 'name',
    customComponent: (row: ClusterDataRow) => (
      <LinkComponent
        to={`/reports/${row.name}/scheduled`}
        isRunning={row.state === 'ONLINE'}
      >
        {row.name}
      </LinkComponent>
    ),
  },
  {
    header: 'Notification',
    columnKey: 'notificationChannels',
    customComponent: ({ notificationChannels }) => (
      <Channels channels={notificationChannels} />
    ),
  },
  {
    header: 'Accuracy',
    columnKey: 'accuracy',
    customComponent: ({ accuracy }) => <AccuracyBadge label={accuracy} />,
  },
  {
    header: 'Updated at',
    columnKey: 'updatedAt',
  },
  {
    header: 'Reports',
    columnKey: 'actions',
    customComponent: (row: ClusterDataRow) => (
      <ReportActionsCell clusterId={row.name} />
    ),
  },
];

const Clusters = () => {
  const [clusters, setClusters] = useState<ClusterDataRow[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchClusters = async () => {
    try {
      const clustersData = await ManagmentServiceApiInstance.getClusters();

      const clusterRows = clustersData.map(
        (cluster): ClusterDataRow => ({
          name: cluster.clusterId,
          accuracy: cluster.accuracy,
          state: transformIsRunningLabel(cluster),
          notificationChannels: transformNotificationChannelsToColumns(cluster),
          updatedAt: transformUpdatedAtDate(cluster),
        }),
      );

      setClusters(clusterRows);
      setIsLoading(false);
    } catch (e: unknown) {
      console.error('Failed to fetch clusters', e);
    }
  };

  useEffect(() => {
    fetchClusters();
  }, []);

  const header = <HeaderWithIcon title={'Clusters'} />;

  return (
    <PageTemplate header={header}>
      <SectionComponent
        title={'Clusters'}
        icon={<SVGIcon iconName="clusters-icon" />}
      >
        {isLoading && <Spinner />}
        {!isLoading && clusters.length > 0 && (
          <Table columns={columns} rows={clusters} />
        )}
        {!isLoading && clusters.length === 0 && (
          <div>No registered clusters yet</div>
        )}
      </SectionComponent>
    </PageTemplate>
  );
};

export default Clusters;
