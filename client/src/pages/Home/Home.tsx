import {ManagmentServiceApiInstance} from 'api/managment-service';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import { useEffect, useState } from 'react';
import ReportDetailsSection from './components/ReportDetailsSection/ReportDetailsSection';
import useReportDetails from 'hooks/useReportStats';

const Home = () => {
  const [lastReportId, setLastReportId] = useState<string | null>(null);

  useEffect(() => {
    const getLatestReport = async () => {
      try {
        const latestReport = await ManagmentServiceApiInstance.getLatestReport();
        if (latestReport) {
          setLastReportId(latestReport.id);
        }
      } catch (e: unknown) {
        console.error('Failed to fetch the latest report', e);
      }
    };

    getLatestReport();
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
