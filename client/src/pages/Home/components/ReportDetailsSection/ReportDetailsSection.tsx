import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import ReportTitle from 'pages/Home/components/ReportTitle/ReportTitle';
import {useReportStats, IncidentStats} from 'hooks/useReportStats';
import {ReportWithDetails} from 'api/managment-service';
import StatisticsDisplay, {
  StatItemData,
} from 'components/StatisticsDisplay/StatisticsDisplay';
import colors from 'global/colors';
import ReportDetailsSubsection from 'pages/Home/components/Subsection/Subsection';
import IncidentList from 'components/IncidentList/IncidentList';
import {
  GenericIncident,
  mapSimplifiedApplicationIncidentsToGenericFormat,
  mapSimplifiedNodeIncidentsToGenericFormat,
  urgencyIncidentCount,
} from 'types/incident';
import './ReportDetailsSection.scss';
import {useNavigate} from 'react-router-dom';
import CenteredSpinner from 'components/CenteredSpinner/CenteredSpinner';

export interface ReportStats {
  lastReport: ReportWithDetails | null;
  isReportLoading: boolean;
}

const statItems = (
  lastReport: ReportWithDetails,
  stats: IncidentStats,
): StatItemData[] => {
  const defaultStats: StatItemData[] = [
    {
      title: 'Analyzed apps',
      value: lastReport.reportDetailedSummary.analyzedApplications,
      unit: 'applications',
      valueColor: colors.urgency.low,
    },
    {
      title: 'Analyzed hosts',
      value: lastReport.reportDetailedSummary.analyzedNodes,
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
      value: lastReport.reportDetailedSummary.totalApplicationEntries,
      unit: 'entries',
      valueColor: colors.urgency.low,
    },
    {
      title: 'Node entries',
      value: lastReport.reportDetailedSummary.totalNodeEntries,
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

const ReportDetailsSection = ({
                                lastReport,
                                isReportLoading,
                              }: ReportStats) => {
  const navigate = useNavigate();

  if (isReportLoading || !lastReport) {
    return <CenteredSpinner/>;
  }
  const incidentStats = useReportStats(lastReport);

  const handleNodeIncidentNavigation = (incident: GenericIncident) => {
    navigate(`/node-incidents/${incident.id}`);
  };

  const handleApplicationIncidentNavigation = (incident: GenericIncident) => {
    navigate(`/application-incidents/${incident.id}`);
  };

  return (
    <SectionComponent
      icon={<SVGIcon iconName="chart-icon"/>}
      title={
        <ReportTitle
          source={lastReport.reportDetailedSummary.clusterId}
          startTime={lastReport.reportDetailedSummary.sinceMs}
          endTime={lastReport.reportDetailedSummary.toMs}
        />
      }
    >
      <div className="dashboard-report-details-section">
          <div className="dashboard-report-details-section__incidents">
            <ReportDetailsSubsection title={'Statistics'}>
              <StatisticsDisplay
                statItems={statItems(lastReport, incidentStats)}
                urgencyIncidentCount={urgencyIncidentCount(incidentStats)}
              />
            </ReportDetailsSubsection>
            {lastReport.applicationIncidents.length > 0 && (
              <ReportDetailsSubsection title="Application incidents">
                <IncidentList
                  incidents={mapSimplifiedApplicationIncidentsToGenericFormat(
                    lastReport.applicationIncidents,
                  )}
                  onClick={handleApplicationIncidentNavigation}
                />
              </ReportDetailsSubsection>
            )}
            {lastReport.nodeIncidents.length > 0 && (
              <ReportDetailsSubsection title="Node incidents">
                <IncidentList
                  incidents={mapSimplifiedNodeIncidentsToGenericFormat(
                    lastReport.nodeIncidents,
                  )}
                  onClick={handleNodeIncidentNavigation}
                />
              </ReportDetailsSubsection>
            )}
          </div>
      </div>
    </SectionComponent>
  );
};
export default ReportDetailsSection;
