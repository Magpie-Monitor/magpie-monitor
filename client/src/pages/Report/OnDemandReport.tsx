import './OnDemandReport.scss';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon.tsx';
import PageTemplate from 'components/PageTemplate/PageTemplate.tsx';
import AccuracySection from './AccuracySection/AccuracySection.tsx';
import DateRangeSection from './DateRangeSection/DateRangeSection.tsx';
import NotificationSection from './NotificationSection/NotificationSection.tsx';
import ApplicationSection from './ApplicationSection/ApplicationSection.tsx';
import NodesSection from './NodesSection/NodesSection.tsx';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import { useParams } from 'react-router-dom';
import { useState } from 'react';
import { NotificationChannel } from './NotificationSection/NotificationSection';
import { ApplicationDataRow } from './ApplicationSection/ApplicationSection';
import { NodeDataRow } from './NodesSection/NodesSection';
import {AccuracyLevel} from 'api/managment-service.ts';

const OnDemandReport = () => {
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
    <PageTemplate header={<HeaderWithIcon title={`Generate report on demand for ${id}`} />}>
        <div className="on-demand-report__section">
                <div className="on-demand-report__row">
                    <AccuracySection setParentAccuracy={setAccuracy} />
                    <DateRangeSection/>
                </div>

                <NotificationSection setNotificationChannels={setNotificationChannels}/>
                <ApplicationSection setApplications={setApplications} clusterId={id ?? ''} defaultAccuracy={accuracy}/>
                <NodesSection setNodes={setNodes} clusterId={id ?? ''} defaultAccuracy={accuracy}/>
            </div>

            <div className="on-demand-report__actions">
                <ActionButton onClick={handleGenerateReport}
                              description="Generate" color={ActionButtonColor.GREEN}/>
                <ActionButton onClick={() => {
                }} description="Cancel" color={ActionButtonColor.RED}/>
            </div>
    </PageTemplate>
    );
};

export default OnDemandReport;