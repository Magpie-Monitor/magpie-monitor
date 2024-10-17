import './ImportantFindings.scss';

const ImportantFindings = () => {
    return (
        <div className="important-findings">
            <h3>Important findings</h3>
            <div className="table-label"></div>
            <div className="finding">
                <div className="finding-source">Postgres:12.19-bullseye</div>
                <div className="finding-category">Database outage</div>
                <div className="finding-summary">
                    Multiple failed authentication attempts after an update</div>
                <div className="finding-date">07.03.2024 15:32</div>
            </div>
            <div className="finding">
                <div className="finding-source">OpenJDK:24-ea-8-jdk</div>
                <div className="finding-category">Internal service issue</div>
                <div className="finding-summary">

                    Abnormal amount of failed requests due to CORS</div>
                <div className="finding-date">07.07.2024 15:32</div>
            </div>
            <div className="finding">
                <div className="finding-source">Kubernetes CNI</div>
                <div className="finding-category">Internal network error</div>
                <div className="finding-summary">
                    Couldn’t start internal DNS server due to configuration error</div>
                <div className="finding-date">07.07.2024 15:32</div>
            </div>
        </div>
    );
};

export default ImportantFindings;
