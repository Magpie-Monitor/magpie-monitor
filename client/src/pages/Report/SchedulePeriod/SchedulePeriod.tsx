import React, {useEffect, useState} from 'react';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import TagButton from 'components/TagButton/TagButton';
import { dateFromTimestampMs } from 'lib/date';
import './SchedulePeriod.scss';

export interface SchedulePeriodProps {
    setGenerationPeriod: (generationPeriod: string) => void;
}

export interface SchedulePeriodOptions {
    periods: string[];
}

export const schedulePeriodOptions: SchedulePeriodOptions = {
    periods: ['3 days', '5 days', '1 week', '2 weeks', '1 month'],
};

export const periodToMilliseconds: Record<string, number> = {
    '3 days': 3 * 24 * 60 * 60 * 1000,
    '5 days': 5 * 24 * 60 * 60 * 1000,
    '1 week': 7 * 24 * 60 * 60 * 1000,
    '2 weeks': 14 * 24 * 60 * 60 * 1000,
    '1 month': 30 * 24 * 60 * 60 * 1000,
};

const SchedulePeriod: React.FC<SchedulePeriodProps> = ({ setGenerationPeriod }) => {
    const [period, setPeriod] = useState<string>('1 week');
    const [nextReportDate, setNextReportDate] = useState<string>('');

    useEffect(() => {
        const calculateNextReportDate = () => {
            const now = Date.now();
            const periodMs = periodToMilliseconds[period] || 0;
            const nextTimestamp = now + periodMs;
            setNextReportDate(dateFromTimestampMs(nextTimestamp));
        };

        calculateNextReportDate();
    }, [period]);

    const handlePeriodChange = (newPeriod: string) => {
        setPeriod(newPeriod);
        setGenerationPeriod(newPeriod);
    };

    return (
        <SectionComponent
            icon={<SVGIcon iconName="precision-icon" />}
            title="Schedule period"
        >
            <div className="schedule-period">
                <TagButton
                    listItems={schedulePeriodOptions.periods}
                    chosenItem={period}
                    onSelect={handlePeriodChange}
                />
                <div className="schedule-period__text">
                    <p className="schedule-period__next-report">
                        *Next report (generation may take up to 24h): {nextReportDate}
                    </p>
                </div>
            </div>
        </SectionComponent>
    );
};

export default SchedulePeriod;
