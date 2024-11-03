import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import IncidentList from 'components/IncidentList/IncidentList';
import Spinner from 'components/Spinner/Spinner';
import useReportDetails, { IncidentStats } from 'hooks/useReportStats';
import StatisticsDisplay, {
  StatItemData,
} from 'components/StatisticsDisplay/StatisticsDisplay';
import {
  ApplicationIncident,
  NodeIncident,
  ReportDetails,
  UrgencyLevel,
} from 'api/managment-service';
import colors from 'global/colors';
import ReportTitle from './components/ReportTitle/ReportTitle';
import { useParams } from 'react-router-dom';
import './ReportDetails.scss';
import ReportDetailsSubsection from './components/Subsection/Subsection';

const urgencyIncidentCount = (
  stats: IncidentStats,
): Record<UrgencyLevel, number> => ({
  LOW: stats.lowUrgencyIncidents,
  MEDIUM: stats.mediumUrgencyIncidents,
  HIGH: stats.highUrgencyIncidents,
});

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

const genericIncidentsFromApplicationIncidents = (
  incidents: ApplicationIncident[],
) =>
  incidents.map((incident) => ({
    source: incident.applicationName,
    category: incident.category,
    urgency: incident.urgency,
    title: incident.title,
    timestamp: incident.sources[0].timestamp,
  }));

const genericIncidentsFromNodeIncidents = (incidents: NodeIncident[]) =>
  incidents.map((incident) => ({
    source: incident.nodeName,
    category: incident.category,
    urgency: incident.urgency,
    title: incident.title,
    timestamp: incident.sources[0].timestamp,
  }));

const ReportDetailsPage = () => {
  const { id } = useParams();
  const {
    incidents,
    report,
    incidentStats,
    areIncidentsLoading,
    isReportLoading,
  } = useReportDetails(id!);

  if (isReportLoading || !report) {
    return <Spinner />;
  }

  return (
    <PageTemplate header={<HeaderWithIcon title={'Report'} />}>
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
        <div className="report-details">
          {areIncidentsLoading && <Spinner />}
          {incidents && incidentStats && (
            <>
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
                />
              </ReportDetailsSubsection>

              <ReportDetailsSubsection title="Node incidents">
                <IncidentList
                  incidents={genericIncidentsFromNodeIncidents(
                    incidents.nodeIncidents,
                  )}
                />
              </ReportDetailsSubsection>
            </>
          )}
        </div>
      </SectionComponent>
    </PageTemplate>
  );
};

export default ReportDetailsPage;
