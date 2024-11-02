import React from 'react';
import './StateBadge.scss';

interface StateBadgeProps {
    label: 'ONLINE' | 'OFFLINE';
}

const StateBadge: React.FC<StateBadgeProps> = ({ label }) => {
    return (
        <span className={`state-badge state-badge--${label.toLowerCase()}`}>
            {label}
        </span>
    );
};

export default StateBadge;
