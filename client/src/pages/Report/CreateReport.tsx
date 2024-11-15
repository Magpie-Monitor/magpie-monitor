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
import {AccuracyLevel, ManagmentServiceApiInstance, ReportPost, ReportType}
    from 'api/managment-service.ts';
import GeneratedInfoPopup from './GeneratedInfoPopup/GeneratedInfoPopup.tsx';
import ReportGenerationType from './StateSection/ReportGenerationType.tsx';
import {dateFromTimestampMs} from 'lib/date.ts';

const CreateReport = () => {
    const {id} = useParams<{ id: string }>();
    const [notificationChannels, setNotificationChannels] = useState<NotificationChannel[]>([]);
    const [applications, setApplications] = useState<ApplicationDataRow[]>([]);
    const [nodes, setNodes] = useState<NodeDataRow[]>([]);
    const [accuracy, setAccuracy] = useState<AccuracyLevel>('HIGH');
    const [generationType, setGenerationType] = useState<ReportType>('ON DEMAND');
    const navigate = useNavigate();
    const [startDateMs, setStartDateMs] = useState<number>(Date.now());
    const [endDateMs, setEndDateMs] = useState<number>(Date.now());
    const [showInfoPopup, setShowInfoPopup] = useState(false);

    const handleDateRangeChange = (startMs: number, endMs: number) => {
        setStartDateMs(startMs);
        setEndDateMs(endMs);
    };

    useEffect(() => {
        if (generationType === 'SCHEDULED' && id) {
            const fetchClusterDetails = async () => {
                try {
                    const clusterDetails = await ManagmentServiceApiInstance.getClusterDetails(id);
                    const mappedNotificationChannels = [
                        ...clusterDetails.slackReceivers.map(receiver => ({
                            id: receiver.id.toString(),
                            name: receiver.receiverName,
                            details: receiver.webhookUrl,
                            service: 'SLACK',
                            added: dateFromTimestampMs(receiver.createdAt),
                            updated: dateFromTimestampMs(receiver.updatedAt),

                        })),
                        ...clusterDetails.discordReceivers.map(receiver => ({
                            id: receiver.id.toString(),
                            name: receiver.receiverName,
                            details: receiver.webhookUrl,
                            service: 'DISCORD',
                            added: dateFromTimestampMs(receiver.createdAt),
                            updated: dateFromTimestampMs(receiver.updatedAt),
                        })),
                        ...clusterDetails.emailReceivers.map(receiver => ({
                            id: receiver.id.toString(),
                            name: receiver.receiverName,
                            details: receiver.receiverEmail,
                            service: 'EMAIL',
                            added: dateFromTimestampMs(receiver.createdAt),
                            updated: dateFromTimestampMs(receiver.updatedAt),
                        })),
                    ];
                    setNotificationChannels(mappedNotificationChannels);

                    const mappedApplications = 
                        clusterDetails.applicationConfigurations.map(config => ({
                        name: config.name,
                        kind: config.kind,
                        accuracy: config.accuracy,
                        customPrompt: config.customPrompt,
                    }));
                    setApplications(mappedApplications);

                    const mappedNodes = clusterDetails.nodeConfigurations.map(config => ({
                        name: config.name,
                        accuracy: config.accuracy,
                        customPrompt: config.customPrompt,
                    }));
                    setNodes(mappedNodes);

                    setAccuracy(clusterDetails.accuracy);
                } catch (error) {
                    console.error('Failed to fetch cluster details:', error);
                }
            };

            fetchClusterDetails();
        }
    }, [generationType, id]);

    const filterNotificationChannels = (channels: NotificationChannel[]) => {
        const slackReceiverIds: number[] = [];
        const discordReceiverIds: number[] = [];
        const mailReceiverIds: number[] = [];

        channels.forEach((channel) => {
            const channelId = parseInt(channel.id, 10);
            if (!isNaN(channelId)) {
                switch (channel.service) {
                    case 'SLACK':
                        slackReceiverIds.push(channelId);
                        break;
                    case 'DISCORD':
                        discordReceiverIds.push(channelId);
                        break;
                    case 'EMAIL':
                        mailReceiverIds.push(channelId);
                        break;
                    default:
                        console.warn(`Unknown service: ${channel.service}`);
                }
            } else {
                console.warn(`Invalid channel id: ${channel.id}`);
            }
        });

        return {slackReceiverIds, discordReceiverIds, mailReceiverIds};
    };


    const handleGenerateReport = () => {
        const {slackReceiverIds, discordReceiverIds, mailReceiverIds} =
            filterNotificationChannels(notificationChannels);

        const report: ReportPost = {
            clusterId: id ?? '',
            accuracy: 'HIGH',
            sinceMs: startDateMs,
            toMs: endDateMs,
            slackReceiverIds: slackReceiverIds,
            discordReceiverIds: discordReceiverIds,
            mailReceiverIds: mailReceiverIds,
            applicationConfigurations: applications.map((app) => ({
                applicationName: app.name,
                accuracy: app.accuracy,
                customPrompt: app.customPrompt,
            })),
            nodeConfigurations: nodes.map((node) => ({
                nodeName: node.name,
                accuracy: node.accuracy,
                customPrompt: node.customPrompt,
            })),
        };
        // ManagmentServiceApiInstance.generateOnDemandReport(report);
        console.log(report);
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
                        <DateRangeSection onDateChange={handleDateRangeChange}/>
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