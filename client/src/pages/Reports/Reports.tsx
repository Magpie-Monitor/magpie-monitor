import './Reports.scss';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import Table, {TableColumn} from 'components/Table/Table.tsx';
import {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {ManagmentServiceApiInstance, ReportSummary} from 'api/managment-service';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import UrgencyBadge from 'components/UrgencyBadge/UrgencyBadge.tsx';

const Reports = () => {
    const [rows, setRows] = useState<ReportSummary[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const navigate = useNavigate();

    const handleRowClick = (id: string) => {
        navigate(`/reports/${id}`);
    };

    const columns: Array<TableColumn<ReportSummary>>  = [
        {
            header: 'Cluster',
            columnKey: 'clusterId',
            customComponent: (row: ReportSummary) => (
                <a className='reports__content__link' href="#" onClick={() => handleRowClick(row.id)}>
                    {row.clusterId}
                </a>
            )
        },
        {header: 'Title', columnKey: 'title'},
        {
            header: 'Urgency',
            columnKey: 'urgency',
            customComponent: (row: ReportSummary) => <UrgencyBadge label={row.urgency}/>
        },
        {header: 'Start date', columnKey: 'startDate'},
        {header: 'End date', columnKey: 'endDate'}
    ];

    const fetchReports = async () => {
        try {
            const reports = await ManagmentServiceApiInstance.getReports();
            const mappedReports = reports.map((report: ReportSummary) => ({
                ...report,
                startDate: new Date(report.sinceMs / 1e6).toLocaleString(),
                endDate: new Date(report.toMs / 1e6).toLocaleString(),
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
        <div className="reports">
            <div className="reports__content">
                <div>
                    <div className='reports__content__heading'>
                        <SVGIcon iconName='reports-list-icon'/>
                        <p className="reports__content__heading__paragraph">Reports</p>
                    </div>
                    <div className="reports__content__dashboard">
                        <SectionComponent
                            icon={'setting-icon'}
                            title={'Weekly reports'}
                        >
                            <div className="reports__content__dashboard__content">
                                {loading ? (
                                    <p className="reports__content__dashboard__content__paragraph">Loading...</p>
                                ) : rows.length === 0 ? (
                                    <p className="reports__content__dashboard__content__paragraph">No reports. Generate new report (TBA: link)</p>
                                ) : (
                                    <Table
                                        columns={columns}
                                        rows={rows}
                                    />
                                )}</div>
                        </SectionComponent>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Reports;
