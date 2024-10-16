import React from 'react';
import './StatItem.scss';

interface StatItemProps {
    title: string;
    value: number;
    unit: string;
    valueColor: string;
}

const StatItem: React.FC<StatItemProps> = ({ title, value, unit, valueColor }) => {
    return (
        <div className="stat-item">
            <div className="stat-title">{title}</div>
            <div className="stat-value">
                <span className="value" style={{ color: valueColor }}>{value}</span>
                <span className="unit">{unit}</span>
            </div>
        </div>
    );
};

export default StatItem;
