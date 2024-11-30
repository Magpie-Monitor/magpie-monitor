import { ManagmentServiceApiInstance } from 'api/managment-service';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import useReportDetails from 'hooks/useReportStats';

import { useEffect, useState } from 'react';
import ReportDetailsSection from './components/ReportDetailsSection/ReportDetailsSection';

const Home = () => {
  const [lastReportId, setLastReportId] = useState<string | null>(null);

  useEffect(() => {
    const fetchReports = async () => {
      try {
        const [onDemandReports, scheduledReports] = await Promise.all([
          ManagmentServiceApiInstance.getReports('ON_DEMAND'),
          ManagmentServiceApiInstance.getReports('SCHEDULED'),
        ]);
        const reports = [...onDemandReports, ...scheduledReports];
        if (reports.length > 0) {
          reports.sort((a, b) => b.requestedAtMs - a.requestedAtMs);
          setLastReportId(reports[0].id);
        }
      } catch (e: unknown) {
        console.error('Failed to fetch reports');
      }
    };

    fetchReports();
  }, []);

  const {
    incidents,
    report,
    incidentStats,
    areIncidentsLoading,
    isReportLoading,
  } = useReportDetails(lastReportId);

  return (
    <PageTemplate header={<HeaderWithIcon title={'Dashboard'} />}>
      <ReportDetailsSection
        report={report}
        incidents={incidents}
        incidentStats={incidentStats}
        areIncidentsLoading={areIncidentsLoading}
        isReportLoading={isReportLoading}
      />
    </PageTemplate>
  );
};

export default Home;
