import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import ReportTitle from 'pages/Home/components/ReportTitle/ReportTitle';
import Spinner from 'components/Spinner/Spinner';
import { IncidentStats, ReportStats } from 'hooks/useReportStats';
import { ReportDetails } from 'api/managment-service';
import StatisticsDisplay, {
  StatItemData,
} from 'components/StatisticsDisplay/StatisticsDisplay';
import colors from 'global/colors';
import ReportDetailsSubsection from 'pages/Home/components/Subsection/Subsection';
import IncidentList from 'components/IncidentList/IncidentList';
import {
  GenericIncident,
  genericIncidentsFromApplicationIncidents,
  genericIncidentsFromNodeIncidents,
  urgencyIncidentCount,
} from 'types/incident';
import './ReportDetailsSection.scss';
import { useNavigate } from 'react-router-dom';

const statItems = (
  report: ReportDetails,
  stats: IncidentStats,
): StatItemData[] => [
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
    title: 'Medium incidents',
    value: stats.mediumUrgencyIncidents,
    unit: 'incidents',
    valueColor: colors.urgency.medium,
  },
  {
    title: 'Low incidents',
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

const ReportDetailsSection = ({
  report,
  incidentStats,
  areIncidentsLoading,
  isReportLoading,
  incidents,
}: ReportStats) => {
  const navigate = useNavigate();

  const handleNodeIncidentNavigation = (incident: GenericIncident) => {
    navigate(`/node-incidents/${incident.id}`);
  };

  const handleApplicationIncidentNavigation = (incident: GenericIncident) => {
    navigate(`/application-incidents/${incident.id}`);
  };

  if (isReportLoading || !report) {
    return <Spinner />;
  }
  return (
    <SectionComponent
      icon={<SVGIcon iconName="chart-icon" />}
      title={
        <ReportTitle
          source={report.clusterId}
          startTime={report.sinceMs}
          endTime={report.toMs}
        />
      }
    >
      <div className="dashboard-report-details-section">
        {areIncidentsLoading && <Spinner />}
        {incidents && incidentStats && (
          <div className="dashboard-report-details-section__incidents">
            <ReportDetailsSubsection title={'Statistics'}>
              <StatisticsDisplay
                statItems={statItems(report, incidentStats)}
                urgencyIncidentCount={urgencyIncidentCount(incidentStats)}
              />
            </ReportDetailsSubsection>

            <ReportDetailsSubsection title="Application incidents">
              <IncidentList
                incidents={genericIncidentsFromApplicationIncidents(
                  incidents.applicationIncidents,
                )}
                onClick={handleApplicationIncidentNavigation}
              />
            </ReportDetailsSubsection>

            <ReportDetailsSubsection title="Node incidents">
              <IncidentList
                incidents={genericIncidentsFromNodeIncidents(
                  incidents.nodeIncidents,
                )}
                onClick={handleNodeIncidentNavigation}
              />
            </ReportDetailsSubsection>
          </div>
        )}
      </div>
    </SectionComponent>
  );
};
export default ReportDetailsSection;
