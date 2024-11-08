import React from 'react';
import './StateBadge.scss';

interface StateBadgeProps {
    label: 'UP' | 'DOWN';
    className?: string;
}

const StateBadge: React.FC<StateBadgeProps> = ({ label, className = '' }) => {
    return (
        <span
            className={`state-badge state-badge--${label.toLowerCase()} ${className}`}
        >
            {label}
        </span>
    );
};

export default StateBadge;
