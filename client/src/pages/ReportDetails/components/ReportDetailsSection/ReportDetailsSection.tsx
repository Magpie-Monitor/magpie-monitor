import SectionComponent from '@components/SectionComponent/SectionComponent';
import SVGIcon from '@components/SVGIcon/SVGIcon';
import ReportTitle from '../ReportTitle/ReportTitle';

const ReportDetailsSection = () => {
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
  );
};
