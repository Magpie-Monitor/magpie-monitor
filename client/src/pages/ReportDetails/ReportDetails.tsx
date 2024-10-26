import './ReportDetails.scss';
import { useParams } from 'react-router-dom';

const ReportDetails = () => {
    const { id } = useParams();

    return (
        <div className="reports-details">
            <div className="reports-details__content">
                <p className="reports-details__content__heading">Report ID: {id}</p>
            </div>
        </div>
    );
};

export default ReportDetails;
