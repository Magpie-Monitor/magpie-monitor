import React from 'react';
import './UrgencyBadge.scss';

interface UrgencyBadgeProps {
    label: 'HIGH' | 'MEDIUM' | 'LOW';
}

const UrgencyBadge: React.FC<UrgencyBadgeProps> = ({ label }) => {
    return (
        <span className={`urgency-badge urgency-badge--${label.toLowerCase()}`}>
            {label}
        </span>
    );
};

export default UrgencyBadge;
