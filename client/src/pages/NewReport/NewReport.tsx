import './NewReport.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import PrecisionSection from './PrecisionSection/PrecisionSection';
import DateRangeSection from './DateRangeSection/DateRangeSection';
import NotificationSection from './NotificationSection/NotificationSection';
import ApplicationSection from './ApplicationSection/ApplicationSection';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton';

const NewReport = () => {
    return (
        <div className="new-report">
                <div className="new-report__header">
                    <SVGIcon iconName="reports-list-icon" />
                    <p className="new-report__title">
                        Generate report for production-services
                    </p>
                </div>

                <div className="new-report__section">
                    <div className="new-report__row">
                        <PrecisionSection />
                        <DateRangeSection />
                    </div>

                    <NotificationSection />
                    <ApplicationSection />
                </div>

                <div className="new-report__actions">
                    <ActionButton onClick={() => {}} description="Confirm" color={ActionButtonColor.GREEN} />
                    <ActionButton onClick={() => {}} description="Cancel" color={ActionButtonColor.RED} />
                </div>
        </div>
    );
};

export default NewReport;
