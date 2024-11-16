import { useState } from 'react';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import './DateRangeSection.scss';
import { getDateFromTimestamps } from 'lib/date.ts';

interface DateRangeSectionProps {
    onDateChange: (startMs: number, endMs: number) => void;
}

const DateRangeSection = ({ onDateChange }: DateRangeSectionProps) => {
    const [startDate, setStartDate] = useState(Date.now());
    const [endDate, setEndDate] = useState(Date.now());

    const getStartOfDay = (dateMs: number) => {
        const date = new Date(dateMs);
        date.setHours(0, 0, 0, 0);
        return date.getTime();
    };

    const getEndOfDay = (dateMs: number) => {
        const date = new Date(dateMs);
        date.setHours(23, 56, 56, 999);
        return date.getTime();
    };

    const handleStartDateChange = (event: { target: { value: string } }) => {
        const startMs = getStartOfDay(Date.parse(event.target.value));
        setStartDate(startMs);
        onDateChange(startMs, getEndOfDay(endDate));
    };

    const handleEndDateChange = (event: { target: { value: string } }) => {
        const endMs = getEndOfDay(Date.parse(event.target.value));
        setEndDate(endMs);
        onDateChange(getStartOfDay(startDate), endMs);
    };

    return (
        <SectionComponent icon={<SVGIcon iconName="date-range-icon" />} title={'Date Range'}>
            <div className="date-range">
                <label>
                    Start Date:
                    <input
                        type="date"
                        value={getDateFromTimestamps(startDate)}
                        onChange={handleStartDateChange}
                    />
                </label>
                <label>
                    End Date:
                    <input
                        type="date"
                        value={getDateFromTimestamps(endDate)}
                        onChange={handleEndDateChange}
                    />
                </label>
            </div>
        </SectionComponent>
    );
};

export default DateRangeSection;
