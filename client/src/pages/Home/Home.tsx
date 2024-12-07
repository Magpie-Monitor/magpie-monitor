import {ManagmentServiceApiInstance, ReportWithDetails} from 'api/managment-service';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import { useEffect, useState } from 'react';
import ReportDetailsSection from './components/ReportDetailsSection/ReportDetailsSection';

const Home = () => {
  const [lastReport, setLastReport] = useState<ReportWithDetails | null>(null);
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

  return (
    <PageTemplate header={<HeaderWithIcon title={'Dashboard'} />}>
      <ReportDetailsSection
        lastReport={lastReport}
        isReportLoading={isReportLoading}
      />
    </PageTemplate>
  );
};

export default Home;
