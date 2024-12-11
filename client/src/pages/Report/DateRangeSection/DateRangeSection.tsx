import { useState } from 'react';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import './DateRangeSection.scss';
import { DateRangePicker } from '@nextui-org/date-picker';
import './DateRangeSection.scss';
import {
  AnyCalendarDate,
  getLocalTimeZone,
  today,
} from '@internationalized/date';

interface DateRangeSectionProps {
  onDateChange: (startMs: number, endMs: number) => void;
}

const startOfDayFromCalendarDate = (calendarDate: AnyCalendarDate): number => {
  return new Date(
    calendarDate.year,
    calendarDate.month - 1,
    calendarDate.day,
    0,
    0,
    0,
    0,
  ).getTime();
};

const endOfDayFromCalendarDate = (calendarDate: AnyCalendarDate): number => {
  return new Date(
    calendarDate.year,
    calendarDate.month - 1,
    calendarDate.day,
    23,
    59,
    59,
    999,
  ).getTime();
};

const DateRangeSection = ({ onDateChange }: DateRangeSectionProps) => {
  const [startDate, setStartDate] = useState(today(getLocalTimeZone()));
  const [endDate, setEndDate] = useState(today(getLocalTimeZone()));

  return (
    <SectionComponent
      icon={<SVGIcon iconName="date-range-icon" />}
      title={'Date Range'}
    >
      <DateRangePicker
        selectorIcon={<SVGIcon iconName="date-range-selector-icon" />}
        className="date-range"
        classNames={{
          popoverContent: 'radius-large date-range__popover',
        }}
        value={{
          start: startDate,
          end: endDate,
        }}
        maxValue={today(getLocalTimeZone())}
        onChange={(dates) => {
          if (!dates) {
            return;
          }
          setStartDate(dates.start);
          setEndDate(dates.end);
          onDateChange(
            startOfDayFromCalendarDate(dates.start),
            endOfDayFromCalendarDate(dates.end),
          );
        }}
      />
    </SectionComponent>
  );
};
export default DateRangeSection;
