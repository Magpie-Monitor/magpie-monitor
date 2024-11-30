import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, {TableColumn} from 'components/Table/Table.tsx';
import {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {
  ReportAwaitingGeneration,
  ManagmentServiceApiInstance,
  ReportSummary,
} from 'api/managment-service';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import UrgencyBadge from 'components/UrgencyBadge/UrgencyBadge.tsx';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import Spinner from 'components/Spinner/Spinner.tsx';
import {dateFromTimestampMs} from 'lib/date.ts';
import './Reports.scss';

const Reports = () => {
  const [rowsOnDemand, setRowsOnDemand] = useState<ReportSummary[]>([]);
  const [rowsScheduled, setRowsScheduled] = useState<ReportSummary[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const navigate = useNavigate();

  const handleRowClick = (id: string) => {
    navigate(`/reports/${id}`);
  };

  const columns: Array<TableColumn<ReportSummary>> = [
    {
      header: 'Cluster',
      columnKey: 'clusterId',
      customComponent: (row: ReportSummary) => (
        <LinkComponent to="" onClick={() => handleRowClick(row.id)}>
          {row.clusterId}
        </LinkComponent>
      ),
    },
    {header: 'Title', columnKey: 'title'},
    {
      header: 'Urgency',
      columnKey: 'urgency',
      customComponent: (row: ReportSummary) => (
        row.urgency ? <UrgencyBadge label={row.urgency}/> : null
      ),
    },
    {header: 'Start date', columnKey: 'startDate'},
    {header: 'End date', columnKey: 'endDate'},
  ];

  const fetchReportsOnDemand = async () => {
    try {
      const reports = await ManagmentServiceApiInstance.getReports('ON_DEMAND');
      const mappedReports = reports.map((report: ReportSummary) => ({
        ...report,
        startDate: dateFromTimestampMs(report.sinceMs),
        endDate: dateFromTimestampMs(report.toMs),
      }))
        .sort((a, b) => b.requestedAtMs - a.requestedAtMs);
      setRowsOnDemand(mappedReports);
    } catch (error) {
      console.error('Error fetching on-demand reports:', error);
    }
  };

  const fetchReportsScheduled = async () => {
    try {
      const reports = await ManagmentServiceApiInstance.getReports('SCHEDULED');
      const mappedReports = reports.map((report: ReportSummary) => ({
        ...report,
        startDate: dateFromTimestampMs(report.sinceMs),
        endDate: dateFromTimestampMs(report.toMs),
      }))
        .sort((a, b) => b.requestedAtMs - a.requestedAtMs);
      setRowsScheduled(mappedReports);
    } catch (error) {
      console.error('Error fetching scheduled reports:', error);
    }
  };

  const fetchReportAwaitingGenerations = async () => {
    try {
      const reports = await ManagmentServiceApiInstance.getAwaitingGenerationReports();
      const mappedReports = reports.map((report: ReportAwaitingGeneration) => ({
        ...report,
        id: `${report.clusterId}-${report.sinceMs}`,
        title: 'Awaiting generation...',
        startDate: dateFromTimestampMs(report.sinceMs),
        endDate: dateFromTimestampMs(report.toMs),
        urgency: null,
        requestedAtMs: Date.now(),
      }));

      const onDemandReports = mappedReports.filter(report => report.reportType === 'ON_DEMAND');
      const scheduledReports = mappedReports.filter(report => report.reportType === 'SCHEDULED');

      setRowsOnDemand(prev => [
        ...onDemandReports.map(report => ({
          ...report
        })).sort((a, b) => b.requestedAtMs - a.requestedAtMs),
        ...prev,
      ]);

      setRowsScheduled(prev => [
        ...scheduledReports.map(report => ({
          ...report
        })).sort((a, b) => b.requestedAtMs - a.requestedAtMs),
        ...prev,
      ]);
    } catch (error) {
      console.error('Error fetching generating reports:', error);
    }
  };

  useEffect(() => {
    const fetchAllReports = async () => {
      setLoading(true);
      await Promise.all([
        fetchReportsOnDemand(),
        fetchReportsScheduled(),
        fetchReportAwaitingGenerations(),
      ]);
      setLoading(false);
    };
    fetchAllReports();
  }, []);

  return (
    <PageTemplate
      header={
        <HeaderWithIcon
          title={'Reports'}
          icon={<SVGIcon iconName="reports-list-icon"/>}
        />
      }
    >
      <div className="reports">
        <SectionComponent
          icon={<SVGIcon iconName="chart-icon"/>}
          title={'Generated reports scheduled'}
        >
          {loading ? (
            <Spinner/>
          ) : rowsScheduled.length === 0 ? (
            <>
              <p>No reports. &nbsp;</p>
              <LinkComponent to="/clusters">Generate new report</LinkComponent>
            </>
          ) : (
            <Table columns={columns} rows={rowsScheduled}/>
          )}
        </SectionComponent>

        <SectionComponent
          icon={<SVGIcon iconName="chart-icon"/>}
          title={'Generated reports on demand'}
        >
          {loading ? (
            <Spinner/>
          ) : rowsOnDemand.length === 0 ? (
            <>
              <p>No reports. &nbsp;</p>
              <LinkComponent to="/clusters">Generate new report</LinkComponent>
            </>
          ) : (
            <Table columns={columns} rows={rowsOnDemand}/>
          )}
        </SectionComponent>
      </div>
    </PageTemplate>
  );
};

export default Reports;
