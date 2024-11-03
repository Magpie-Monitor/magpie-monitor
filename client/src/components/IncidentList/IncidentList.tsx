import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import './IncidentList.scss';
import { dateFromTimestampMs } from 'lib/date';
import { GenericIncident } from 'types/incident';

interface IncidentListProps {
  incidents: GenericIncident[];
}

const IncidentList = ({ incidents }: IncidentListProps) => {
  return (
    <div className="incident-list">
      <div className="incident-list__headers">
        <div className="incident-list__header">Source</div>
        <div className="incident-list__header">Category</div>
        <div className="incident-list__header">Title</div>
        <div className="incident-list__header">Date</div>
      </div>
      {incidents.map((incident, index) => (
        <div className="incident-list__entry" key={index}>
          <div className="incident-list__entry__source">{incident.source}</div>
          <div className="incident-list__entry__category">
            <SVGIcon iconName={'fire-icon'} /> {incident.category}
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
