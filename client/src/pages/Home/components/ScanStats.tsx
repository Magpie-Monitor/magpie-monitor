import StatItem from './StatItem';
import './ScanStats.scss';

const ScanStats = () => {
    return (
        <div className="scan-stats">
            <h3 className="scan-stats__heading">Scan stats</h3>
            <div className="scan-stats__content">
                <div className="scan-stats__items">
                    <StatItem title="Analyzed apps"
                              value={368} unit="applications" valueColor="#5CD060"/>
                    <StatItem title="Analyzed hosts"
                              value={24} unit="hosts" valueColor="#5CD060"/>
                    <StatItem title="Kamil Nowak counter"
                              value={3} unit="Kamil Nowak’s" valueColor="#5CD060"/>
                    <StatItem title="Critical incidents"
                              value={145} unit="incidents" valueColor="#E01300"/>
                    <StatItem title="Application entries"
                              value={38721} unit="entries" valueColor="#5CD060"/>
                    <StatItem title="Node entries"
                              value={12938} unit="entries" valueColor="#5CD060"/>
                </div>
                <div className="scan-stats__chart">
                    <p>Chart Placeholder</p>
                </div>
            </div>
        </div>
    );
};

export default ScanStats;
