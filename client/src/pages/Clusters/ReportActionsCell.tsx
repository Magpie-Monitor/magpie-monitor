import React from 'react';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import { useNavigate } from 'react-router-dom';

interface ReportActionsCellProps {
    clusterId: string;
}

const ReportActionsCell: React.FC<ReportActionsCellProps> = ({ clusterId }) => {
    const navigate = useNavigate();

    const handleNewReportOnDemand = () => {
        navigate(`/reports/${clusterId}/on-demand`);
    };

    const handleScheduledReport = () => {
        navigate(`/reports/${clusterId}/scheduled`);
    };

    return (
        <div className='clusters--button'>
            <ActionButton
                onClick={handleScheduledReport}
                description="Scheduled"
                color={ActionButtonColor.GREEN}
            />
            <ActionButton
                onClick={handleNewReportOnDemand}
                description="On Demand"
                color={ActionButtonColor.GREEN}
            />
        </div>
    );
};

export default ReportActionsCell;
