import {ManagmentServiceApiInstance, ReportDetails} from 'api/managment-service';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import { useEffect, useState } from 'react';
import ReportDetailsSection from './components/ReportDetailsSection/ReportDetailsSection';
import useReportDetails from 'hooks/useReportStats';

const Home = () => {
  const [lastReport, setLastReport] = useState<ReportDetails | null>(null);
  const [isReportLoading, setIsReportLoading] = useState(true);

  useEffect(() => {
    const getLatestReport = async () => {
      try {
        const latestReport = await ManagmentServiceApiInstance.getLatestReport();
        if (latestReport) {
          setLastReport(latestReport);
        }
      } catch (e: unknown) {
        console.error('Failed to fetch the latest report', e);
      }
    };

    getLatestReport();
    setIsReportLoading(false);
  }, []);

  const {
    incidents,
    incidentStats,
    areIncidentsLoading,
  } = useReportDetails(lastReport);

  return (
    <PageTemplate header={<HeaderWithIcon title={'Dashboard'} />}>
      <ReportDetailsSection
        report={lastReport}
        incidents={incidents}
        incidentStats={incidentStats}
        areIncidentsLoading={areIncidentsLoading}
        isReportLoading={isReportLoading}
      />
    </PageTemplate>
  );
};

export default Home;
