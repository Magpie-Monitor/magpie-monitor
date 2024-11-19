import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, {TableColumn} from 'components/Table/Table.tsx';
import {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {
    GeneratingReport,
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
import CustomTag from 'components/CustomTag/CustomTag.tsx';

const Reports = () => {
    const [rowsOnDemand, setRowsOnDemand] = useState<ReportSummary[]>([]);
    const [rowsScheduled, setRowsScheduled] = useState<ReportSummary[]>([]);
    const [rowsGenerating, setRowsGenerating] = useState<GeneratingReport[]>([]);
    const [loadingOnDemand, setLoadingOnDemand] = useState<boolean>(true);
    const [loadingScheduled, setLoadingScheduled] = useState<boolean>(true);
    const [loadingGenerating, setLoadingGenerating] = useState<boolean>(true);
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
                <UrgencyBadge label={row.urgency}/>
            ),
        },
        {header: 'Start date', columnKey: 'startDate'},
        {header: 'End date', columnKey: 'endDate'},
    ];

    const columnsGenerating: Array<TableColumn<GeneratingReport>> = [
        {
            header: 'Cluster',
            columnKey: 'clusterId',
            customComponent: (row: GeneratingReport) => (
                <LinkComponent to="">
                    {row.clusterId}
                </LinkComponent>
            ),
        },
        {
            header: 'Report Type',
            columnKey: 'reportType',
            customComponent: (row: GeneratingReport) => <CustomTag name={row.reportType}/>,
        },
        {header: 'Start date', columnKey: 'startDate'},
        {header: 'End date', columnKey: 'endDate'},
    ];

    const fetchReportsOnDemand = async () => {
        try {
            const reports = await ManagmentServiceApiInstance.getReportsOnDemand();
            const mappedReports = reports.map((report: ReportSummary) => ({
                ...report,
                startDate: dateFromTimestampMs(report.sinceMs),
                endDate: dateFromTimestampMs(report.toMs),
            }));
            setRowsOnDemand(mappedReports);
        } catch (error) {
            console.error('Error fetching on-demand reports:', error);
        } finally {
            setLoadingOnDemand(false);
        }
    };

    const fetchReportsScheduled = async () => {
        try {
            const reports = await ManagmentServiceApiInstance.getReportsScheduled();
            const mappedReports = reports.map((report: ReportSummary) => ({
                ...report,
                startDate: dateFromTimestampMs(report.sinceMs),
                endDate: dateFromTimestampMs(report.toMs),
            }));
            setRowsScheduled(mappedReports);
        } catch (error) {
            console.error('Error fetching scheduled reports:', error);
        } finally {
            setLoadingScheduled(false);
        }
    };

    const fetchGeneratingReports = async () => {
        try {
            const reports = await ManagmentServiceApiInstance.getGeneratingReports();
            const mappedReports = reports.map((report: GeneratingReport) => ({
                ...report,
                startDate: dateFromTimestampMs(report.sinceMs),
                endDate: dateFromTimestampMs(report.toMs),
            }));
            setRowsGenerating(mappedReports);
        } catch (error) {
            console.error('Error fetching generating reports:', error);
        } finally {
            setLoadingGenerating(false);
        }
    };

    useEffect(() => {
        fetchReportsOnDemand();
        fetchReportsScheduled();
        fetchGeneratingReports();
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
            <div className="reports__body">
                {rowsGenerating.length > 0 && (
                    <SectionComponent
                        icon={<SVGIcon iconName="chart-icon"/>}
                        title={'Reports currently generating'}
                    >
                        {loadingGenerating ? (
                            <Spinner/>
                        ) : (
                            <Table columns={columnsGenerating} rows={rowsGenerating}/>
                        )}
                    </SectionComponent>
                )}

                <SectionComponent
                    icon={<SVGIcon iconName="chart-icon"/>}
                    title={'Generated reports scheduled'}
                >
                    {loadingScheduled ? (
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
                    {loadingOnDemand ? (
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
