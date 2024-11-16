import './CreateReport.scss';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon.tsx';
import PageTemplate from 'components/PageTemplate/PageTemplate.tsx';
import AccuracySection from './AccuracySection/AccuracySection.tsx';
import DateRangeSection from './DateRangeSection/DateRangeSection.tsx';
import NotificationSection from './NotificationSection/NotificationSection.tsx';
import ApplicationSection from './ApplicationSection/ApplicationSection.tsx';
import NodesSection from './NodesSection/NodesSection.tsx';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';
import {useNavigate, useParams} from 'react-router-dom';
import {useState, useEffect} from 'react';
import {NotificationChannel} from './NotificationSection/NotificationSection';
import {ApplicationDataRow} from './ApplicationSection/ApplicationSection';
import {NodeDataRow} from './NodesSection/NodesSection';
import GeneratedInfoPopup from './GeneratedInfoPopup/GeneratedInfoPopup.tsx';
import ReportGenerationType from './StateSection/ReportGenerationType.tsx';
import SchedulePeriod, {schedulePeriodOptions} from './SchedulePeriod/SchedulePeriod.tsx';
import {fetchClusterData, generateReport} from './CreateReportUtils.tsx';
import {AccuracyLevel, ReportType} from 'api/managment-service.ts';

const CreateReport = () => {
    const {id} = useParams<{ id: string }>();
    const [notificationChannels, setNotificationChannels] = useState<NotificationChannel[]>([]);
    const [applications, setApplications] = useState<ApplicationDataRow[]>([]);
    const [nodes, setNodes] = useState<NodeDataRow[]>([]);
    const [accuracy, setAccuracy] = useState<AccuracyLevel>('HIGH');
    const [generationType, setGenerationType] = useState<ReportType>('ON DEMAND');
    const [generationPeriod, setGenerationPeriod] =
        useState<string>(schedulePeriodOptions.periods[2]);
    const navigate = useNavigate();
    const [startDateMs, setStartDateMs] = useState<number>(Date.now());
    const [endDateMs, setEndDateMs] = useState<number>(Date.now());
    const [showInfoPopup, setShowInfoPopup] = useState(false);

    const handleDateRangeChange = (startMs: number, endMs: number) => {
        setStartDateMs(startMs);
        setEndDateMs(endMs);
    };

    useEffect(() => {
        const fetchData = async () => {
            const {notificationChannels, applications, nodes} = await fetchClusterData(id || '');
            setNotificationChannels(notificationChannels);
            setApplications(applications);
            setNodes(nodes);
        };
        fetchData();
    }, [generationType, id]);

    const handleGenerateReport = () => {
        generateReport({
            id,
            notificationChannels,
            applications,
            nodes,
            generationType,
            accuracy,
            generationPeriod,
            startDateMs,
            endDateMs,
        });
        setShowInfoPopup(true);
    };

    const handleCancelReport = () => {
        navigate('/dashboard');
    };

    return (
        <PageTemplate header={<HeaderWithIcon title={`Generate report on demand for ${id}`}/>}>
            <div className="on-demand-report">
                <div className="on-demand-report__wrapper">
                    <div className="on-demand-report__row">
                        <div className="on-demand-report__row">
                            <AccuracySection setParentAccuracy={setAccuracy}/>
                            <ReportGenerationType setParentGenerationType={setGenerationType}/>
                        </div>
                        {generationType === 'ON DEMAND' ? (
                            <DateRangeSection onDateChange={handleDateRangeChange}/>
                        ) : (
                            <SchedulePeriod setGenerationPeriod={setGenerationPeriod}/>
                        )}
                    </div>
                </div>
                <NotificationSection notificationChannels={notificationChannels}
                                     setNotificationChannels={setNotificationChannels}/>
                <ApplicationSection applications={applications} setApplications={setApplications}
                                    clusterId={id ?? ''} defaultAccuracy={accuracy}/>
                <NodesSection nodes={nodes} setNodes={setNodes}
                              clusterId={id ?? ''} defaultAccuracy={accuracy}/>
            </div>

            <div className="on-demand-report__actions">
                <ActionButton onClick={handleGenerateReport}
                              description="Generate" color={ActionButtonColor.GREEN}/>
                <ActionButton onClick={handleCancelReport}
                              description="Cancel" color={ActionButtonColor.RED}/>
            </div>
            <GeneratedInfoPopup
                isDisplayed={showInfoPopup}
                onClose={() => setShowInfoPopup(false)}
            />
        </PageTemplate>
    );
};

export default CreateReport;
