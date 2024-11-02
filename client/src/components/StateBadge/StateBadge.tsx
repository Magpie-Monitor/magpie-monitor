import React from 'react';
import './StateBadge.scss';

interface StateBadgeProps {
    label: 'UP' | 'DOWN';
}

const StateBadge: React.FC<StateBadgeProps> = ({ label }) => {
    return (
        <span className={`state-badge state-badge--${label.toLowerCase()}`}>
            {label}
        </span>
    );
};

export default StateBadge;
