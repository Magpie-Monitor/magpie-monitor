import SectionComponent from 'components/SectionComponent/SectionComponent';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import Table, {TableColumn} from 'components/Table/Table';
import './Clusters.scss';
import Channels from './components/NotificationChannelsColumn/NotificationChannelsColumn';
import {useEffect, useState} from 'react';
import {
  ClusterSummary,
  ManagmentServiceApiInstance,
  AccuracyLevel,
  NotificationChannelKind,
} from 'api/managment-service';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import ReportActionsCell from './ReportActionsCell';
import AccuracyBadge from 'components/AccuracyBadge/AccuracyBadge.tsx';
import {dateTimeFromTimestampMs} from 'lib/date.ts';
import CenteredSpinner from 'components/CenteredSpinner/CenteredSpinner';

interface ClusterDataRow {
  name: string;
  state: 'ONLINE' | 'OFFLINE';
  accuracy: AccuracyLevel;
  notificationChannels: NotificationChannelColumn[];
  updatedAt: string;

  [key: string]: string | NotificationChannelColumn[];
}

export interface NotificationChannelColumn {
  kind: NotificationChannelKind;
  name: string;
}

const transformNotificationChannelsToColumns = (
  cluster: ClusterSummary,
): NotificationChannelColumn[] => {
  return cluster.slackReceivers
    .map(
      (channel): NotificationChannelColumn => ({
        kind: 'SLACK',
        name: channel.receiverName,
      }),
    )
    .concat(
      cluster.discordReceivers.map((channel) => ({
        kind: 'DISCORD',
        name: channel.receiverName,
      })),
    )
    .concat(
      cluster.emailReceivers.map((channel) => ({
        kind: 'EMAIL',
        name: channel.receiverName,
      })),
    );
};
const transformIsRunningLabel = (
  cluster: ClusterSummary,
): ClusterDataRow['state'] => {
  return cluster.running ? 'ONLINE' : 'OFFLINE';
};

const columns: Array<TableColumn<ClusterDataRow>> = [
  {
    header: 'Name',
    columnKey: 'name',
    customComponent: (row: ClusterDataRow) => (
      <LinkComponent
        to={`/clusters/${row.name}/report`}
        isRunning={row.state === 'ONLINE'}
      >
        {row.name}
      </LinkComponent>
    ),
  },
  {
    header: 'Notification',
    columnKey: 'notificationChannels',
    customComponent: ({
                        notificationChannels,
                      }: {
      notificationChannels: NotificationChannelColumn[];
    }) => {
      return <Channels channels={notificationChannels}/>;
    },
  },
  {
    header: 'Accuracy',
    columnKey: 'accuracy',
    customComponent: (row: ClusterDataRow) => {
      return <AccuracyBadge label={row.accuracy}/>;
    },
  },
  {
    header: 'Updated at',
    columnKey: 'updatedAt',
  },
  {
    header: 'Reports',
    columnKey: 'actions',
    customComponent: (row: ClusterDataRow) => {
      return <ReportActionsCell clusterId={row.name}/>;
    },
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
          updatedAt: dateTimeFromTimestampMs(cluster.updatedAtMillis),
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

  const header = <HeaderWithIcon
                                title={'Clusters'}
                                icon={<SVGIcon iconName="clusters-icon"/>}
                              />;

  return (
    <PageTemplate header={header}>
      <SectionComponent
        title={'Clusters'}
        icon={<SVGIcon iconName="hive-icon"/>}
      >
        {isLoading && <CenteredSpinner/>}
        {!isLoading && clusters.length > 0 && (
          <Table columns={columns} rows={clusters}/>
        )}
        {!isLoading && clusters.length === 0 && (
          <div>No registered clusters yet</div>
        )}
      </SectionComponent>
    </PageTemplate>
  );
};

export default Clusters;
