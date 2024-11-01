import './NewReport.scss';
import HeaderWithIcon from 'components/PageTemplate/components/HeaderWithIcon/HeaderWithIcon.tsx';
import PageTemplate from 'components/PageTemplate/PageTemplate.tsx';
import AccuracySection from './AccuracySection/AccuracySection.tsx';
import DateRangeSection from './DateRangeSection/DateRangeSection';
import NotificationSection from './NotificationSection/NotificationSection';
import ApplicationSection from './ApplicationSection/ApplicationSection';
import NodesSection from './NodesSection/NodesSection';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton';
import StateSection from './StateSection/StateSection.tsx';

const NewReport = () => {
    return (
    <PageTemplate header={<HeaderWithIcon title={'Generate report for production-services'} />}>
            <div className="new-report__section">
                <div className="new-report__row">
                    <div className="new-report__row">
                    <StateSection/>
                    <AccuracySection/>
                    </div>
                    <DateRangeSection/>
                </div>

                <NotificationSection/>
                <ApplicationSection/>
                <NodesSection/>
            </div>

            <div className="new-report__actions">
                <ActionButton onClick={() => {
                }} description="Generate" color={ActionButtonColor.GREEN}/>
                <ActionButton onClick={() => {
                }} description="Cancel" color={ActionButtonColor.RED}/>
            </div>
    </PageTemplate>
    );
};

export default NewReport;