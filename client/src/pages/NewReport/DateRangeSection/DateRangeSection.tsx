import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';

const DateRangeSection = () => {

    return (
        <SectionComponent icon={<SVGIcon iconName='date-range-icon' />} title={'Date Range'}>
            <p>Date range picker :D</p>
        </SectionComponent>
    );
};

export default DateRangeSection;