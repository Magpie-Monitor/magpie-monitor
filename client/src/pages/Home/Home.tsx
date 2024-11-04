import { ManagmentServiceApiInstance } from 'api/managment-service';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import Spinner from 'components/Spinner/Spinner';
import useReportDetails from 'hooks/useReportStats';

import { useEffect, useState } from 'react';
import ReportDetailsSection from './components/ReportDetailsSection/ReportDetailsSection';

const Home = () => {
  const [lastReportId, setLastReportId] = useState<string | null>(null);

  const {
    incidents,
    report,
    incidentStats,
    areIncidentsLoading,
    isReportLoading,
  } = useReportDetails(lastReportId!);

  useEffect(() => {
    const fetchReports = async () => {
      try {
        const reports = await ManagmentServiceApiInstance.getReports();
        if (reports.length > 0) {
          setLastReportId(reports[reports.length - 1].id);
        }
      } catch (e: unknown) {
        console.error('Failed to fetch reports');
      }
    };

    fetchReports();
  }, []);

  if (isReportLoading || !report) {
    return <Spinner />;
  }

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
