import './OnDemandReport.scss';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon.tsx';
import PageTemplate from 'components/PageTemplate/PageTemplate.tsx';
import AccuracySection from './AccuracySection/AccuracySection.tsx';
import DateRangeSection from './DateRangeSection/DateRangeSection.tsx';
import NotificationSection from './NotificationSection/NotificationSection.tsx';
import ApplicationSection from './ApplicationSection/ApplicationSection.tsx';
import NodesSection from './NodesSection/NodesSection.tsx';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import {useNavigate, useParams} from 'react-router-dom';
import { useState } from 'react';
import { NotificationChannel } from './NotificationSection/NotificationSection';
import { ApplicationDataRow } from './ApplicationSection/ApplicationSection';
import { NodeDataRow } from './NodesSection/NodesSection';
import {AccuracyLevel, ManagmentServiceApiInstance, ReportPost} from 'api/managment-service.ts';
import GeneratedInfoPopup from './GeneratedInfoPopup/GeneratedInfoPopup.tsx';

const OnDemandReport = () => {
    const { id } = useParams<{ id: string }>();
    const [notificationChannels, setNotificationChannels] = useState<NotificationChannel[]>([]);
    const [applications, setApplications] = useState<ApplicationDataRow[]>([]);
    const [nodes, setNodes] = useState<NodeDataRow[]>([]);
    const [accuracy, setAccuracy] = useState<AccuracyLevel>('HIGH');
    const navigate = useNavigate();
    const [startDateMs, setStartDateMs] = useState<number>(Date.now());
    const [endDateMs, setEndDateMs] = useState<number>(Date.now());
    const [showInfoPopup, setShowInfoPopup] = useState(false);

    const handleDateRangeChange = (startMs: number, endMs: number) => {
        setStartDateMs(startMs);
        setEndDateMs(endMs);
    };

    const filterNotificationChannels = (channels: NotificationChannel[]) => {
        const slackReceiverIds: number[] = [];
        const discordReceiverIds: number[] = [];
        const mailReceiverIds: number[] = [];

        channels.forEach((channel) => {
            const channelId = parseInt(channel.id, 10);
            if (!isNaN(channelId)) {
                switch (channel.service.toLowerCase()) {
                    case 'slack':
                        slackReceiverIds.push(channelId);
                        break;
                    case 'discord':
                        discordReceiverIds.push(channelId);
                        break;
                    case 'email':
                        mailReceiverIds.push(channelId);
                        break;
                    default:
                        console.warn(`Unknown service: ${channel.service}`);
                }
            } else {
                console.warn(`Invalid channel id: ${channel.id}`);
            }
        });

        return { slackReceiverIds, discordReceiverIds, mailReceiverIds };
    };


    const handleGenerateReport = () => {
        const { slackReceiverIds, discordReceiverIds, mailReceiverIds } =
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
        // console.log(report);
        ManagmentServiceApiInstance.generateOnDemandReport(report);
        setShowInfoPopup(true);
    };

    const handleCancelReport = () => {
        navigate('/dashboard');
    };

    return (
    <PageTemplate header={<HeaderWithIcon title={`Generate report on demand for ${id}`} />}>
        <div className="on-demand-report">
            <div className="on-demand-report__wrapper">
                <div className="on-demand-report__row">
                    <AccuracySection setParentAccuracy={setAccuracy} />
                    <DateRangeSection onDateChange={handleDateRangeChange} />
                </div>
            </div>
                <NotificationSection setNotificationChannels={setNotificationChannels}/>
                <ApplicationSection setApplications={setApplications}
                                    clusterId={id ?? ''} defaultAccuracy={accuracy}/>
                <NodesSection setNodes={setNodes} clusterId={id ?? ''} defaultAccuracy={accuracy}/>
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

export default OnDemandReport;