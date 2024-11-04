import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import './IncidentList.scss';
import { dateFromTimestampMs } from 'lib/date';
import { GenericIncident } from 'types/incident';
import { UrgencyLevel } from '@api/managment-service';

interface IncidentListProps {
  incidents: GenericIncident[];
  onClick?: (incident: GenericIncident) => void;
}

const categoryUrgencyClass: Record<UrgencyLevel, string> = {
  LOW: 'low-urgency',
  MEDIUM: 'medium-urgency',
  HIGH: 'high-urgency',
};

const IncidentList = ({ incidents, onClick }: IncidentListProps) => {
  if (incidents.length === 0) {
    return <div className="incident-list--no-incidents">No incidents</div>;
  }
  return (
    <div className="incident-list">
      <div className="incident-list__headers">
        <div className="incident-list__header">Source</div>
        <div className="incident-list__header">Category</div>
        <div className="incident-list__header">Title</div>
        <div className="incident-list__header">Date</div>
      </div>
      {incidents.map((incident, index) => (
        <div
          className="incident-list__entry"
          key={index}
          onClick={onClick ? () => onClick(incident) : () => { }}
        >
          <div className="incident-list__entry__source">{incident.source}</div>
          <div
            className={`incident-list__entry__category--${categoryUrgencyClass[incident.urgency]}`}
          >
            <SVGIcon
              iconName={`incident-category-icon--${categoryUrgencyClass[incident.urgency]}`}
            />{' '}
            {incident.category}
          </div>
          <div className="incident-list__entry__summary">{incident.title}</div>
          <div className="incident-list__entry__date">
            {dateFromTimestampMs(incident.timestamp)}
          </div>
        </div>
      ))}
    </div>
  );
};

export default IncidentList;
