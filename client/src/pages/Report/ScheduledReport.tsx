import './ScheduledReport.scss';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon.tsx';
import PageTemplate from 'components/PageTemplate/PageTemplate.tsx';
import AccuracySection from './AccuracySection/AccuracySection.tsx';
import DateRangeSection from './DateRangeSection/DateRangeSection.tsx';
import NotificationSection, {NotificationChannel}
    from './NotificationSection/NotificationSection.tsx';
import ApplicationSection, {ApplicationDataRow} from './ApplicationSection/ApplicationSection.tsx';
import NodesSection, {NodeDataRow} from './NodesSection/NodesSection.tsx';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';
import StateSection from './StateSection/StateSection.tsx';
import {useParams} from 'react-router-dom';
import {useState} from 'react';
import {AccuracyLevel} from 'api/managment-service.ts';

const ScheduledReport = () => {
    const { id } = useParams<{ id: string }>();
    const [notificationChannels, setNotificationChannels] = useState<NotificationChannel[]>([]);
    const [applications, setApplications] = useState<ApplicationDataRow[]>([]);
    const [nodes, setNodes] = useState<NodeDataRow[]>([]);
    const [accuracy, setAccuracy] = useState<AccuracyLevel>('HIGH');

    const handleGenerateReport = () => {
        console.log('Notification Channels:', notificationChannels);
        console.log('Applications:', applications);
        console.log('Nodes:', nodes);
    };

    return (
        <PageTemplate header={<HeaderWithIcon title={`Configure scheduled report for ${id}`}/>}>
            <div className="scheduled-report__section">
                <div className="scheduled-report__row">
                    <div className="scheduled-report__row">
                        <StateSection/>
                        <AccuracySection setParentAccuracy={setAccuracy} />
                    </div>
                    <DateRangeSection/>
                </div>
                <NotificationSection setNotificationChannels={setNotificationChannels}/>
                <ApplicationSection setApplications={setApplications}
                                    clusterId={id ?? ''} defaultAccuracy={accuracy}/>
                <NodesSection setNodes={setNodes} clusterId={id ?? ''} defaultAccuracy={accuracy}/>
            </div>

            <div className="scheduled-report__actions">
                <ActionButton onClick={handleGenerateReport}
                              description="Generate" color={ActionButtonColor.GREEN}/>
                <ActionButton onClick={() => {
                }} description="Cancel" color={ActionButtonColor.RED}/>
            </div>
        </PageTemplate>
    );
};

export default ScheduledReport;
