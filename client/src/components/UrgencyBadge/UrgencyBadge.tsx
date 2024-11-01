import React from 'react';
import './UrgencyBadge.scss';
import { AccuracyLevel } from 'api/managment-service.ts';

interface UrgencyBadgeProps {
    label: AccuracyLevel;
}

const UrgencyBadge: React.FC<UrgencyBadgeProps> = ({ label }) => {
    return (
        <span className={`urgency-badge urgency-badge--${label.toLowerCase()}`}>
            {label}
        </span>
    );
};

export default UrgencyBadge;
