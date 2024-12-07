import { IncidentStats } from '@hooks/useReportStats';
import {
  ApplicationIncidentSummary,
  NodeIncidentSummary,
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
    timestamp: incident.sinceMs,
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
    timestamp: incident.sinceMs,
  }));

export const mapSimplifiedApplicationIncidentsToGenericFormat = (
  incidents: ApplicationIncidentSummary[],
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
  incidents: NodeIncidentSummary[],
): GenericIncident[] =>
  incidents.map((incident) => ({
    id: incident.id,
    source: incident.nodeName,
    category: incident.category,
    urgency: incident.urgency,
    title: incident.title,
    timestamp: incident.sinceMs,
  }));