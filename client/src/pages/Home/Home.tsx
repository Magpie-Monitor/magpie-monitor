import {ManagmentServiceApiInstance, ReportDetailedWithIncidents} from 'api/managment-service';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import { useEffect, useState } from 'react';
import ReportDetailsSection from './components/ReportDetailsSection/ReportDetailsSection';

const Home = () => {
  const [lastReport, setLastReport] = useState<ReportDetailedWithIncidents | null>(null);
  const [isReportLoading, setIsReportLoading] = useState(true);

  useEffect(() => {
    const fetchNewestReport = async () => {
      try {
        const newestReport = await ManagmentServiceApiInstance.getNewestReport();
        if (newestReport) {
          setLastReport(newestReport);
        }
      } catch (e: unknown) {
        console.error('Failed to fetch the newest report', e);
      }
    };

    fetchNewestReport();
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
