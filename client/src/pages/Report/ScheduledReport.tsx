import './ScheduledReport.scss';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon.tsx';
import PageTemplate from 'components/PageTemplate/PageTemplate.tsx';
import AccuracySection from './AccuracySection/AccuracySection.tsx';
import DateRangeSection from './DateRangeSection/DateRangeSection.tsx';
import NotificationSection from './NotificationSection/NotificationSection.tsx';
import ApplicationSection from './ApplicationSection/ApplicationSection.tsx';
import NodesSection from './NodesSection/NodesSection.tsx';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';
import StateSection from './StateSection/StateSection.tsx';
import {useParams} from 'react-router-dom';

const ScheduledReport = () => {
    const {id} = useParams<{ id: string }>();

    return (
        <PageTemplate header={<HeaderWithIcon title={`Configure scheduled report for ${id}`}/>}>
            <div className="scheduled-report__section">
                <div className="scheduled-report__row">
                    <div className="scheduled-report__row">
                        <StateSection/>
                        <AccuracySection/>
                    </div>
                    <DateRangeSection/>
                </div>
                <NotificationSection/>
                <ApplicationSection/>
                <NodesSection/>
            </div>

            <div className="scheduled-report__actions">
                <ActionButton onClick={() => {
                }} description="Generate" color={ActionButtonColor.GREEN}/>
                <ActionButton onClick={() => {
                }} description="Cancel" color={ActionButtonColor.RED}/>
            </div>
        </PageTemplate>
    );
};

export default ScheduledReport;
