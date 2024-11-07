import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, { TableColumn } from 'components/Table/Table.tsx';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  ManagmentServiceApiInstance,
  ReportSummary,
} from 'api/managment-service';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import UrgencyBadge from 'components/UrgencyBadge/UrgencyBadge.tsx';
import PageTemplate from 'components/PageTemplate/PageTemplate';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon';
import LinkComponent from 'components/LinkComponent/LinkComponent.tsx';
import Spinner from 'components/Spinner/Spinner.tsx';

const Reports = () => {
  const [rows, setRows] = useState<ReportSummary[]>([]);
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
    { header: 'Title', columnKey: 'title' },
    {
      header: 'Urgency',
      columnKey: 'urgency',
      customComponent: (row: ReportSummary) => (
        <UrgencyBadge label={row.urgency} />
      ),
    },
    { header: 'Start date', columnKey: 'startDate' },
    { header: 'End date', columnKey: 'endDate' },
  ];

  const fetchReports = async () => {
    try {
      const reports = await ManagmentServiceApiInstance.getReports();
      const mappedReports = reports.map((report: ReportSummary) => ({
        ...report,
        startDate: new Date(report.sinceMs).toLocaleString(),
        endDate: new Date(report.toMs).toLocaleString(),
      }));
      setRows(mappedReports);
    } catch (error) {
      console.error('Error fetching reports:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReports();
  }, []);

  return (
    <PageTemplate
      header={
        <HeaderWithIcon
          title={'Reports'}
          icon={<SVGIcon iconName="reports-list-icon" />}
        />
      }
    >
      <SectionComponent
        icon={<SVGIcon iconName="chart-icon" />}
        title={'Weekly reports'}
      >
        {loading ? (
          <Spinner />
        ) : rows.length === 0 ? (
          <p>No reports. Generate new report (TBA: link)</p>
        ) : (
          <Table columns={columns} rows={rows} />
        )}
      </SectionComponent>
    </PageTemplate>
  );
};

export default Reports;
