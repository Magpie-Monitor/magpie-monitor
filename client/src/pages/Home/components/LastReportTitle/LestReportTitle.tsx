import React from 'react';
import './LastReportTitle.scss';

interface SectionComponentTitleProps {
    source: string;
    startTime: string;
    endTime: string;
}

const SectionComponentTitle: React.FC<SectionComponentTitleProps> = ({
    source,
    startTime,
    endTime,
}) => {
    return (
        <div className="section-title">
            <p className="section-title__title">Last report from</p>
            <a href="#" className="section-title__source">
                {source}
            </a>
            <p className="section-title__date-range">
                ({startTime} - {endTime})
            </p>
        </div>
    );
};

export default SectionComponentTitle;
