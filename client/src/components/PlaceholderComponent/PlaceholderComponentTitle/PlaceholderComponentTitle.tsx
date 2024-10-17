import React from 'react';
import './PlaceholderComponentTitle.scss';

interface PlaceholderComponentProps {
    source: string;
    startTime: string;
    endTime: string;
}

const PlaceholderComponentTitle: React.FC<PlaceholderComponentProps> =
    ({source, startTime, endTime}) => {
        return (
            <div className="placeholder-title-div">
                <p className="placeholder-title">
                    Last report from
                </p>
                <a>{source}</a>
                <p className="placeholder-title">
                    ({startTime} - {endTime})
                </p>
            </div>
        );
    };

export default PlaceholderComponentTitle;
