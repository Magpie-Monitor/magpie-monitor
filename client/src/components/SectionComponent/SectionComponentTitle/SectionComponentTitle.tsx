import React from 'react';
import './SectionComponentTitle.scss';

interface PlaceholderComponentProps {
    source: string;
    startTime: string;
    endTime: string;
}

const SectionComponentTitle: React.FC<PlaceholderComponentProps> =
    ({source, startTime, endTime}) => {
        return (
            <div className="section-title">
                <p className="section-title__title">
                    Last report from
                </p>
                <a href="#" className="section-title__source">{source}</a>
                <p className="section-title__date-range">
                    ({startTime} - {endTime})
                </p>
            </div>
        );
    };

export default SectionComponentTitle;
