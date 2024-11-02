import React from 'react';
import './AccuracyBadge.scss';
import { AccuracyLevel } from 'api/managment-service.ts';

interface AccuracyBadgeProps {
    label: AccuracyLevel;
}

const AccuracyBadge: React.FC<AccuracyBadgeProps> = ({ label }) => {
    return (
        <span className={`accuracy-badge accuracy-badge--${label.toLowerCase()}`}>
            {label}
        </span>
    );
};

export default AccuracyBadge;
