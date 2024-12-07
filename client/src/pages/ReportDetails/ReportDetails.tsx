import PageTemplate from 'components/PageTemplate/PageTemplate';
import {useReportDetails,  IncidentStats } from 'hooks/useReportStats';
import { useNavigate, useParams } from 'react-router-dom';
import './ReportDetails.scss';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import StatisticsDisplay, {
  StatItemData,
} from 'components/StatisticsDisplay/StatisticsDisplay';
import { ReportDetails } from 'api/managment-service';
import colors from 'global/colors';
import IncidentList from 'components/IncidentList/IncidentList';
import ReportHeader from './components/ReportHeader/ReportHeader';
import {
  mapApplicationIncidentsToGenericFormat, mapNodeIncidentsToGenericFormat,
  urgencyIncidentCount,
} from 'types/incident';
import CenteredSpinner from 'components/CenteredSpinner/CenteredSpinner';

const statItems = (
  report: ReportDetails,
  stats: IncidentStats,
): StatItemData[] => {
  const defaultStats: StatItemData[] = [
    {
      title: 'Analyzed apps',
      value: report.analyzedApplications,
      unit: 'applications',
      valueColor: colors.urgency.low,
    },
    {
      title: 'Analyzed hosts',
      value: report.analyzedNodes,
      unit: 'hosts',
      valueColor: colors.urgency.low,
    },
    {
      title: 'Critical incidents',
      value: stats.highUrgencyIncidents,
      unit: 'incidents',
      valueColor: colors.urgency.high,
    },
    {
      title: 'Medium urgency incidents',
      value: stats.mediumUrgencyIncidents,
      unit: 'incidents',
      valueColor: colors.urgency.medium,
    },
    {
      title: 'Low urgency incidents',
      value: stats.lowUrgencyIncidents,
      unit: 'incidents',
      valueColor: colors.urgency.low,
    },
    {
      title: 'Application entries',
      value: report.totalApplicationEntries,
      unit: 'entries',
      valueColor: colors.urgency.low,
    },
    {
      title: 'Node entries',
      value: report.totalNodeEntries,
      unit: 'entries',
      valueColor: colors.urgency.low,
    },
  ];

  if (stats.nodeWithMostIncidents.nodeName) {
    defaultStats.push({
      title: 'Node with highest number of incidents',
      value: stats.nodeWithMostIncidents.nodeName,
      unit: '',
      valueColor: colors.urgency.high,
    });

    defaultStats.push({
      title: `Incidents from ${stats.nodeWithMostIncidents.nodeName}`,
      value: stats.nodeWithMostIncidents.numberOfIncidents,
      unit: 'incidents',
      valueColor: colors.urgency.high,
    });
  }

  if (stats.applicationWithMostIncidents.applicationName) {
    defaultStats.push({
      title: 'Application with highest number of incidents',
      value: stats.applicationWithMostIncidents.applicationName,
      unit: '',
      valueColor: colors.urgency.high,
    });

    defaultStats.push({
      title: `Incidents from ${stats.applicationWithMostIncidents.applicationName}`,
      value: stats.applicationWithMostIncidents.numberOfIncidents,
      unit: 'incidents',
      valueColor: colors.urgency.high,
    });
  }

  return defaultStats;
};

const ReportDetailsPage = () => {
  const { id } = useParams();
  const { incidents, report, incidentStats, isReportLoading } =
    useReportDetails(id!);

  const navigate = useNavigate();

  if (isReportLoading || !report || !incidents || !incidentStats) {
    return (
      <PageTemplate header={''}>
        <CenteredSpinner />
      </PageTemplate>
    );
  }

  return (
    <PageTemplate
      header={
        <ReportHeader
          name={report.clusterId}
          sinceMs={report.sinceMs}
          toMs={report.toMs}
        />
      }
    >
      <div className="report-details">
        <SectionComponent
          icon={<SVGIcon iconName={'report-stats-icon'} />}
          title={'Statistics'}
        >
          <StatisticsDisplay
            statItems={statItems(report, incidentStats)}
            urgencyIncidentCount={urgencyIncidentCount(incidentStats)}
          />
        </SectionComponent>
        <SectionComponent
          icon={<SVGIcon iconName={'application-incident-metadata-icon'} />}
          title={'Application incidents'}
        >
          <IncidentList
            incidents={mapApplicationIncidentsToGenericFormat(
              incidents.applicationIncidents,
            )}
            onClick={({ id: incidentId }) =>
              navigate(`/application-incidents/${incidentId}`)
            }
          />
        </SectionComponent>

        <SectionComponent
          icon={<SVGIcon iconName={'node-incident-metadata-icon'} />}
          title={'Node incidents'}
        >
          <IncidentList
            incidents={mapNodeIncidentsToGenericFormat(
              incidents.nodeIncidents,
            )}
            onClick={({ id: incidentId }) =>
              navigate(`/node-incidents/${incidentId}`)
            }
          />
        </SectionComponent>
      </div>
    </PageTemplate>
  );
};

export default ReportDetailsPage;
