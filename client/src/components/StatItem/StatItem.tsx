import React from 'react';
import CountUp from 'react-countup';
import './StatItem.scss';

interface StatItemProps {
  title: string;
  value: string | number;
  unit: string;
  valueColor: string;
}

const StatItem: React.FC<StatItemProps> = ({
  title,
  value,
  unit,
  valueColor,
}) => {
  return (
    <div className="stat">
      <div className="stat__title">{title}</div>
      <div className="stat__value">
        <span className="stat__number" style={{ color: valueColor }}>
          {typeof value === 'string' ? (
            value
          ) : (
            <CountUp start={0} end={value} />
          )}
        </span>
        <span className="stat__unit">{unit}</span>
      </div>
    </div>
  );
};

export default StatItem;
