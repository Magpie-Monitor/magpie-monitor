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
            <div className="placeholder-title">
                <p className="title">
                    Last report from
                </p>
                <a href="#" className="source">{source}</a>
                <p className="date-range">
                    ({startTime} - {endTime})
                </p>
            </div>
        );
    };

export default SectionComponentTitle;
