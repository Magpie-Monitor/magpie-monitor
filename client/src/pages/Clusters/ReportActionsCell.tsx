import React from 'react';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import { useNavigate } from 'react-router-dom';

interface ReportActionsCellProps {
    clusterId: string;
}

const ReportActionsCell: React.FC<ReportActionsCellProps> = ({ clusterId }) => {
    const navigate = useNavigate();

    const handleReportConfiguration = () => {
        navigate(`/clusters/${clusterId}/report`);
    };


    return (
            <ActionButton
                onClick={handleReportConfiguration}
                description="Report configuration"
                color={ActionButtonColor.GREEN}
            />
    );
};

export default ReportActionsCell;
