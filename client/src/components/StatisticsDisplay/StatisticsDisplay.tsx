import StatItem from 'components/StatItem/StatItem.tsx';
import './StatisticsDisplay.scss';
import UrgencyPolarChart from 'components/UrgencyPolarChart/UrgencyPolarChart';
import { UrgencyLevel } from 'api/managment-service';
import colors from 'global/colors';
import { useEffect, useRef } from 'react';

export interface StatItemData {
  title: string;
  value: number | string;
  unit: string;
  valueColor?: string;
}

interface StatisticsDisplayProps {
  statItems: StatItemData[];
  urgencyIncidentCount: Record<UrgencyLevel, number>;
}

const StatisticsDisplay = ({
  statItems,
  urgencyIncidentCount,
}: StatisticsDisplayProps) => {
  // Handle dynamic resizing of chartjs component
  const statsRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<HTMLDivElement>(null);

  const showChart = () => {
    const indexOfData = Object.values(urgencyIncidentCount).findIndex(
      (item) => {
        return item > 0;
      },
    );

    return indexOfData != -1;
  };

  useEffect(() => {
    const adjustHeight = () => {
      if (statsRef.current && chartRef.current) {
        const statsHeight = statsRef.current.offsetHeight;
        chartRef.current.style.height = `${statsHeight}px`;
      }
    };

    adjustHeight();

    window.addEventListener('resize', adjustHeight);

    return () => {
      window.removeEventListener('resize', adjustHeight);
    };
  }, []);

  return (
    <div className="statistics-display">
      <div className="statistics-display__content">
        <div className="statistics-display__items" ref={statsRef}>
          {statItems.map((item, index) => (
            <StatItem
              key={index}
              title={item.title}
              value={item.value}
              unit={item.unit}
              valueColor={item.valueColor || colors.urgency.low}
            />
          ))}
        </div>
        {showChart() && (
          <div className="statistics-display__chart" ref={chartRef}>
            <UrgencyPolarChart urgencyIncidentCount={urgencyIncidentCount} />
          </div>
        )}
      </div>
    </div>
  );
};

export default StatisticsDisplay;
