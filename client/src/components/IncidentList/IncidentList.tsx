import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import './IncidentList.scss';
import { dateTimeFromTimestampMs } from 'lib/date';
import { GenericIncident } from 'types/incident';
import { UrgencyLevel } from '@api/managment-service';
import { useTransition, animated } from '@react-spring/web';

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
  const transitions = useTransition(incidents, {
    from: { opacity: 0, transform: 'translateY(20px)' },
    enter: { opacity: 1, transform: 'translateY(0)' },
    leave: { opacity: 0, transform: 'translateY(20px)' },
    config: { duration: 200 }, // Animation duration
    keys: incidents.map((incident) => incident.id), // Ensure unique keys
    trail: 400,
  });

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
      {transitions((style, incident) => (
        <animated.div
          style={style}
          className="incident-list__entry"
          key={incident.id}
          onClick={onClick ? () => onClick(incident) : () => {}}
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
            {dateTimeFromTimestampMs(incident.timestamp)}
          </div>
        </animated.div>
      ))}
    </div>
  );
};

export default IncidentList;
