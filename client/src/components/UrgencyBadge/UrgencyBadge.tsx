import React from 'react';
import './UrgencyBadge.scss';

interface UrgencyBadgeProps {
    label: string;
}

const UrgencyBadge: React.FC<UrgencyBadgeProps> = ({ label }) => {
    return (
        <span className={`urgency-badge urgency-${label.toLowerCase()}`}>
            {label}
        </span>
    );
};

export default UrgencyBadge;
