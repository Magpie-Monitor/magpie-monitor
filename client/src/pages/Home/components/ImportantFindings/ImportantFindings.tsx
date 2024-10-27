import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import './ImportantFindings.scss';

const ImportantFindings = () => {
    return (
        <div className="important-findings">
            <h3 className="important-findings__title">Important findings</h3>
            <div className="important-findings__table-label"></div>
            <div className="important-findings__columns">
                <div>Source</div>
                <div>Category</div>
                <div>Summary</div>
                <div>Date</div>
            </div>
            <div className="important-findings__finding">
                <div className="important-findings__finding__source">
                    Postgres:12.19-bullseye
                </div>
                <div className="important-findings__finding__category">
                    <SVGIcon iconName={'fire-icon'} /> Database outage
                </div>
                <div className="important-findings__finding__summary">
                    Multiple failed authentication attempts after an update
                </div>
                <div className="important-findings__finding__date">
                    07.03.2024 15:32
                </div>
            </div>
            <div className="important-findings__finding">
                <div className="important-findings__finding__source">
                    OpenJDK:24-ea-8-jdk
                </div>
                <div className="important-findings__finding__category">
                    <SVGIcon iconName={'fire-icon'}/> Internal service issue
                </div>
                <div className="important-findings__finding__summary">
                    Abnormal amount of failed requests due to CORS
                </div>
                <div className="important-findings__finding__date">
                    07.07.2024 15:32
                </div>
            </div>
            <div className="important-findings__finding">
                <div className="important-findings__finding__source">
                    Kubernetes CNI
                </div>
                <div className="important-findings__finding__category">
                    <SVGIcon iconName={'fire-icon'}/> Internal network error
                </div>
                <div className="important-findings__finding__summary">
                    Couldn’t start internal DNS server due to configuration error
                </div>
                <div className="important-findings__finding__date">
                    07.07.2024 15:32
                </div>
            </div>
        </div>
    );
};

export default ImportantFindings;
