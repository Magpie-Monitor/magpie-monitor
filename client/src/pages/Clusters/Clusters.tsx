import SectionComponent from 'components/SectionComponent/SectionComponent';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import Table, { TableColumn } from 'components/Table/Table';
import './Clusters.scss';
import Channels from './components/NotificationChannelsColumn/NotificationChannelsColumn';
import UrgencyBadge from 'components/UrgencyBadge/UrgencyBadge';
import ClusterLink from './components/ClusterLink/ClusterLink';
import { useEffect, useState } from 'react';
import {
  ClusterSummary,
  ManagmentServiceApiInstance,
} from 'api/managment-service';

interface ClusterDataRow {
  name: string;
  state: 'Online' | 'Offline';
  precision: 'HIGH' | 'MEDIUM' | 'LOW';
  notificationChannels: NotificationChannelColumn[];
  updatedAt: string;
  [key: string]: string | NotificationChannelColumn[];
}

type NotificationChannelKind = 'SLACK' | 'DISCORD' | 'EMAIL';

export interface NotificationChannelColumn {
  kind: NotificationChannelKind;
  name: string;
}

const columns: Array<TableColumn<ClusterDataRow>> = [
  {
    header: 'Name',
    columnKey: 'name',
    customComponent: ({ name }) => <ClusterLink name={name} />,
  },
  {
    header: 'State',
    columnKey: 'state',
  },
  {
    header: 'Precision',
    columnKey: 'precision',
    customComponent: ({ precision }) => <UrgencyBadge label={precision} />,
  },
  {
    header: 'Notification',
    columnKey: 'notificationChannels',
    customComponent: ({ notificationChannels }) => (
      <Channels channels={notificationChannels} />
    ),
  },
  {
    header: 'Updated at',
    columnKey: 'updatedAt',
  },
];

const tranformNotificationChannelsToColumns = (
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
  return cluster.isRunning ? 'Online' : 'Offline';
};

const transformUpdatedAtDate = (cluster: ClusterSummary) => {
  const date = new Date(cluster.updatedAt);
  return date.toLocaleString();
};

const Clusters = () => {
  const [clusters, setClusters] = useState<ClusterDataRow[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchClusters = async () => {
    try {
      const clustersData = await ManagmentServiceApiInstance.getClusters();

      const clusterRows = clustersData.map(
        (cluster): ClusterDataRow => ({
          name: cluster.id,
          precision: cluster.precision,
          state: transformIsRunningLabel(cluster),
          notificationChannels: tranformNotificationChannelsToColumns(cluster),
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
      <SectionComponent title={'Clusters'} icon={'clusters-icon'}>
        {isLoading && <div>Loading...</div>}
        {!isLoading && clusters.length > 0 && (
          <Table columns={columns} rows={clusters} />
        )}
        {!isLoading && clusters.length == 0 && (
          <div>No registered clusters yet</div>
        )}
      </SectionComponent>
    </PageTemplate>
  );
};

export default Clusters;
