import { DemoContainer } from '@mui/x-date-pickers/internals/demo';
import { LocalizationProvider } from '@mui/x-date-pickers-pro/LocalizationProvider';
import { AdapterDayjs } from '@mui/x-date-pickers-pro/AdapterDayjs';
import { DateRangePicker } from '@mui/x-date-pickers-pro/DateRangePicker';

export default function BasicDateRangePicker() {
    return (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
            <DemoContainer components={['DateRangePicker']}>
                <DateRangePicker localeText={{ start: 'Check-in', end: 'Check-out' }} />
            </DemoContainer>
        </LocalizationProvider>
    );
}

// import { useState } from 'react';
// import './DateRangeSection.scss';
// import SectionComponent from 'components/SectionComponent/SectionComponent';
// import { LocalizationProvider } from '@mui/x-date-pickers-pro/LocalizationProvider';
// import { AdapterDayjs } from '@mui/x-date-pickers-pro/AdapterDayjs';
// import { DateRangePicker } from '@mui/x-date-pickers-pro/DateRangePicker';
//
// const DateRangeSection = () => {
//     const [selectedStartDate, setSelectedStartDate] = useState<Date | null>(null);
//     const [selectedEndDate, setSelectedEndDate] = useState<Date | null>(null);
//     const [currentMonth, setCurrentMonth] = useState(new Date().getMonth());
//     const [currentYear, setCurrentYear] = useState(new Date().getFullYear());
//
//     const handleDateRangeChange = (newDateRange) => {
//         const [start, end] = newDateRange;
//         setSelectedStartDate(start?.toDate() || null);
//         setSelectedEndDate(end?.toDate() || null);
//     };
//
//     return (
//         <SectionComponent icon={'date-range-icon'} title={'Date Range'}>
//             <LocalizationProvider dateAdapter={AdapterDayjs}>
//                 <DateRangePicker
//                     localeText={{ start: 'Check-in', end: 'Check-out' }}
//                     value={[selectedStartDate, selectedEndDate]}
//                     onChange={handleDateRangeChange}
//                 />
//             </LocalizationProvider>
//             <Calendar
//                 currentMonth={currentMonth}
//                 currentYear={currentYear}
//                 selectedStartDate={selectedStartDate}
//                 selectedEndDate={selectedEndDate}
//                 onMonthChange={setCurrentMonth}
//                 onYearChange={setCurrentYear}
//             />
//         </SectionComponent>
//     );
// };
//
// export default DateRangeSection;
