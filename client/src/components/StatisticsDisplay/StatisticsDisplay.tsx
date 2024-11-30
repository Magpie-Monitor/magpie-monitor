import StatItem from 'components/StatItem/StatItem.tsx';
import './StatisticsDisplay.scss';
import UrgencyPolarChart from 'components/UrgencyPolarChart/UrgencyPolarChart';
import { UrgencyLevel } from 'api/managment-service';
import colors from 'global/colors';

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
    return (
        <div className="statistics-display">
            <div className="statistics-display__content">
                <div className="statistics-display__items">
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
                {statItems.length > 0 && statItems[0].value && (
                    <div className="statistics-display__chart">
                        <UrgencyPolarChart urgencyIncidentCount={urgencyIncidentCount} />
                    </div>
                )}
            </div>
        </div>
    );
};

export default StatisticsDisplay;
