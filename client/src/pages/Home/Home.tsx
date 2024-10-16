import ImportantFindings from './components/ImportantFindings';
import StatItem from './components/StatItem';
import './Home.scss';

const Home = () => {
    return (
        <div className="home-container">
            <div className="main-content">
                <p className="title">Dashboard</p>
                <div className="last-report">
                    <p>Last report from <a href="#">production-services</a>
                        (19.04.2023 - 25.04.2023)</p>
                </div>
                <h3>Scan stats</h3>
                <div className="scan-stats">
                    <div className="dashboard-stats">
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
                    <div className="chart-placeholder">
                        <p>Chart Placeholder</p>
                    </div>
                </div>
                <ImportantFindings/>
            </div>
        </div>
    );
};

export default Home;
