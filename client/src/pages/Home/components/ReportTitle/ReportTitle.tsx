import React from 'react';
import './ReportTitle.scss';
import { dateFromTimestampMs } from 'lib/date';

interface ReportTitleProps {
    source: string;
    startTime: number;
    endTime: number;
}

const ReportTitle: React.FC<ReportTitleProps> = ({
    source,
    startTime,
    endTime,
}) => {
    return (
        <div className="report-title">
            <p className="report-title__title">Last report from</p>
            <a href="#" className="report-title__source">
                {source}
            </a>
            <p className="report-title__date-range">
                ({dateFromTimestampMs(startTime)} - {dateFromTimestampMs(endTime)})
            </p>
        </div>
    );
};

export default ReportTitle;
