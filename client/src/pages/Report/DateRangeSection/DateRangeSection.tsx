import { useState } from 'react';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';

interface DateRangeSectionProps {
    onDateChange: (startMs: number, endMs: number) => void;
}

const DateRangeSection = ({ onDateChange }: DateRangeSectionProps) => {
    const [startDate, setStartDate] = useState(Date.now());
    const [endDate, setEndDate] = useState(Date.now());

    const handleStartDateChange = (event: { target: { value: string; }; }) => {
        const startMs = Date.parse(event.target.value);
        setStartDate(startMs);
        onDateChange(startMs, endDate);
    };

    const handleEndDateChange = (event: { target: { value: string; }; }) => {
        const endMs = Date.parse(event.target.value);
        setEndDate(endMs);
        onDateChange(startDate, endMs);
    };

    return (
        <SectionComponent icon={<SVGIcon iconName='date-range-icon' />} title={'Date Range'}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <label>
                    Start Date:
                    <input
                        type="date"
                        value={new Date(startDate).toISOString().split('T')[0]}
                        onChange={handleStartDateChange}
                    />
                </label>
                <label>
                    End Date:
                    <input
                        type="date"
                        value={new Date(endDate).toISOString().split('T')[0]}
                        onChange={handleEndDateChange}
                    />
                </label>
            </div>
        </SectionComponent>
    );
};

export default DateRangeSection;
