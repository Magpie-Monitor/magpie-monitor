import { IncidentStats } from '@hooks/useReportStats';
import {
  ApplicationIncidentSimplified,
  NodeIncidentSimplified,
  ApplicationIncident,
  NodeIncident,
  UrgencyLevel,
} from 'api/managment-service';

export interface GenericIncident {
  id: string;
  source: string;
  category: string;
  title: string;
  timestamp: number;
  urgency: UrgencyLevel;
}

export const urgencyIncidentCount = (
  stats: IncidentStats,
): Record<UrgencyLevel, number> => ({
  LOW: stats.lowUrgencyIncidents,
  MEDIUM: stats.mediumUrgencyIncidents,
  HIGH: stats.highUrgencyIncidents,
});

export const mapApplicationIncidentsToGenericFormat = (
  incidents: ApplicationIncident[],
): GenericIncident[] =>
  incidents.map((incident) => ({
    id: incident.id,
    source: incident.applicationName,
    category: incident.category,
    urgency: incident.urgency,
    title: incident.title,
    timestamp: incident.sources[0].timestamp,
  }));

export const mapNodeIncidentsToGenericFormat = (
  incidents: NodeIncident[],
): GenericIncident[] =>
  incidents.map((incident) => ({
    id: incident.id,
    source: incident.nodeName,
    category: incident.category,
    urgency: incident.urgency,
    title: incident.title,
    timestamp: incident.sources[0].timestamp,
  }));

export const mapSimplifiedApplicationIncidentsToGenericFormat = (
  incidents: ApplicationIncidentSimplified[],
): GenericIncident[] =>
  incidents.map((incident) => ({
    id: incident.id,
    source: incident.applicationName,
    category: incident.category,
    urgency: incident.urgency,
    title: incident.title,
    timestamp: incident.sinceMs,
  }));

export const mapSimplifiedNodeIncidentsToGenericFormat = (
  incidents: NodeIncidentSimplified[],
): GenericIncident[] =>
  incidents.map((incident) => ({
    id: incident.id,
    source: incident.nodeName,
    category: incident.category,
    urgency: incident.urgency,
    title: incident.title,
    timestamp: incident.sinceMs,
  }));