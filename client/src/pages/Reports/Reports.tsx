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
import {dateFromTimestampMs, dateTimeWithoutSecondsFromTimestampMs} from 'lib/date.ts';
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
      customComponent: (row: ReportSummary) =>
        <LinkComponent>
          {row.clusterId}
        </LinkComponent>
    },
    {
      header: 'Title',
      columnKey: 'title',
      customComponent: (row: ReportSummary) => (
        <div className="reports__title-with-icon">
          {row.urgency === null && (
            <div className="reports__spinner">
              <Spinner size="17px"/>
            </div>
          )}
          <span
            className={`reports__title ${
              row.urgency === null ? 'reports__title--inactive' : ''
            }`}
          >
          {row.title}
        </span>
        </div>
      ),
    },
    {
      header: 'Urgency',
      columnKey: 'urgency',
      customComponent: (row: ReportSummary) =>
        row.urgency ? <UrgencyBadge label={row.urgency}/> : null,
    },
    {
      header: 'Date Range',
      columnKey: 'dateRange',
      customComponent: (row: ReportSummary) => (
        <span>
        {row.startDate} - {row.endDate}
      </span>
      ),
    },
    {
      header: 'Requested at',
      columnKey: 'requestedAtDate',
    },
    {
      header: 'Actions',
      columnKey: '',
      customComponent: (row: ReportSummary) => (
        <button
          className={`reports__action-button ${
            row.urgency === null ? 'reports__action-button--inactive' : ''
          }`}
          onClick={() => row.urgency !== null && handleRowClick(row.id)}
          disabled={row.urgency === null}
        >
          <SVGIcon iconName="open-icon" />
        </button>
      ),
    }
    ,
  ];

  const fetchReportsOnDemand = async () => {
    try {
      const reports = await ManagmentServiceApiInstance.getReports('ON_DEMAND');
      const mappedReports = reports.map((report: ReportSummary) => ({
        ...report,
        startDate: dateFromTimestampMs(report.sinceMs),
        endDate: dateFromTimestampMs(report.toMs),
        requestedAtDate: dateTimeWithoutSecondsFromTimestampMs(report.requestedAtMs),
      }));

      setRowsOnDemand(prev => [
        ...mappedReports,
        ...prev,
      ]);
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
        requestedAtDate: dateTimeWithoutSecondsFromTimestampMs(report.requestedAtMs),
      }));

      setRowsScheduled(prev => [
        ...mappedReports,
        ...prev,
      ]);
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
        requestedAtDate: dateTimeWithoutSecondsFromTimestampMs(Date.now()),
      }));

      const updateRows =
        (filterType: string, setRows: React.Dispatch<React.SetStateAction<ReportSummary[]>>) => {
          const filteredReports = mappedReports.filter(report => report.reportType === filterType);
          setRows(prev => [...filteredReports, ...prev]);
        };

      updateRows('ON_DEMAND', setRowsOnDemand);
      updateRows('SCHEDULED', setRowsScheduled);
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

      setRowsOnDemand(prev => [...prev].sort((a, b) => b.requestedAtMs - a.requestedAtMs));
      setRowsScheduled(prev => [...prev].sort((a, b) => b.requestedAtMs - a.requestedAtMs));

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
          title={'Scheduled reports'}
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
          title={'Reports on demand'}
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
