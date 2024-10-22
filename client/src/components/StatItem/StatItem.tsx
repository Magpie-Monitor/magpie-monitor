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
        <div className="stat">
            <div className="stat__title">{title}</div>
            <div className="stat__value">
                <span className="stat__number" style={{ color: valueColor }}>
                    {value}
                </span>
                <span className="stat__unit">{unit}</span>
            </div>
        </div>
    );
};

export default StatItem;
