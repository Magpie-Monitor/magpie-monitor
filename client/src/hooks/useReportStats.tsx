import {
  UrgencyLevel,
  AllIncidentsFromReport,
  ManagmentServiceApiInstance,
  ReportDetails,
} from 'api/managment-service';
import { groupBy } from 'lib/arrays';
import { useEffect, useState } from 'react';

export interface ReportStats {
  incidents: AllIncidentsFromReport | null;
  report: ReportDetails | null;
  incidentStats: IncidentStats | null;
  areIncidentsLoading: boolean;
  isReportLoading: boolean;
}

export interface IncidentStats {
  totalApplicationIncidents: number;
  totalNodeIncidents: number;
  highUrgencyIncidents: number;
  mediumUrgencyIncidents: number;
  lowUrgencyIncidents: number;
  categoryWithMostIncidents: {
    categoryName: string;
    numberOfIncidents: number;
  };

  applicationWithMostIncidents: {
    applicationName: string;
    numberOfIncidents: number;
  };

  nodeWithMostIncidents: {
    nodeName: string;
    numberOfIncidents: number;
  };
}

const getRecordWithLongestValueArray = <T,>(
  entries: Record<string, Array<T>>,
): [string, number] => {
  let highestValue = 0;
  let highestValueKey = '';
  for (const [key, value] of Object.entries(entries)) {
    if (value.length > highestValue) {
      highestValue = value.length;
      highestValueKey = key;
    }
  }

  return [highestValueKey, highestValue];
};

const groupIncidentsByCategory = <T extends { category: string }>(
  incidents: T[],
): Record<string, T[]> => {
  return groupBy(incidents, (incident) => incident.category);
};

const groupIncidentsByUrgency = <T extends { urgency: UrgencyLevel }>(
  incidents: T[],
): Record<UrgencyLevel, T[]> => {
  return groupBy(incidents, (incident) => incident.urgency);
};

const groupIncidentsByApplication = <T extends { applicationName: string }>(
  incidents: T[],
): Record<string, T[]> => {
  return groupBy(incidents, (incident) => incident.applicationName);
};

const groupIncidentsByNode = <T extends { nodeName: string }>(
  incidents: T[],
): Record<string, T[]> => {
  return groupBy(incidents, (incident) => incident.nodeName);
};

export const useReportDetails = (reportId: string | null) => {
  const [report, setReport] = useState<ReportDetails | null>(null);
  const [incidents, setIncidents] = useState<AllIncidentsFromReport | null>(
    null,
  );
  const [isReportLoading, setIsReportLoading] = useState(true);
  const [areIncidentsLoading, setAreIncidentsLoading] = useState(true);
  const [incidentStats, setIncidentStats] = useState<IncidentStats | null>(
    null,
  );

  useEffect(() => {
    if (!reportId) return;

    const fetchReport = async (id: string) => {
      try {
        const reportData = await ManagmentServiceApiInstance.getReport(id);
        setReport(reportData);
        setIsReportLoading(false);
      } catch (e: unknown) {
        console.error('Failed to fetch report by id', id);
      }
    };

    fetchReport(reportId);
  }, [reportId]);

  useEffect(() => {
    if (!report || !reportId) return;
    const fetchIncidents = async () => {
      try {
        const incidentsData =
          await ManagmentServiceApiInstance.getIncidentsFromReport(reportId);
        setIncidents(incidentsData);
        setAreIncidentsLoading(false);
      } catch (e: unknown) {
        console.error('Failed to fetch report by id', reportId);
      }
    };

    if (!report) {
      return;
    }

    fetchIncidents();
  }, [reportId, report]);

  useEffect(() => {
    if (!incidents) return;

    const incidentsByUrgency = groupIncidentsByUrgency([
      ...incidents.applicationIncidents,
      ...incidents.nodeIncidents,
    ]);
    const incidentsByCategory = groupIncidentsByCategory([
      ...incidents.applicationIncidents,
      ...incidents.nodeIncidents,
    ]);
    const [mostPopularCategory, mostPopularCategoryIncidents] =
      getRecordWithLongestValueArray(incidentsByCategory);

    const incidentsByApp = groupIncidentsByApplication(
      incidents.applicationIncidents,
    );
    const [mostPopularApp, mostPopularAppIncidents] =
      getRecordWithLongestValueArray(incidentsByApp);

    const incidentsByNode = groupIncidentsByNode(incidents.nodeIncidents);
    const [mostPopularNode, mostPopularNodeIncidents] =
      getRecordWithLongestValueArray(incidentsByNode);

    setIncidentStats({
      totalApplicationIncidents: incidents.applicationIncidents.length,
      totalNodeIncidents: incidents.nodeIncidents.length,
      highUrgencyIncidents: incidentsByUrgency.HIGH
        ? incidentsByUrgency.HIGH.length
        : 0,
      mediumUrgencyIncidents: incidentsByUrgency.MEDIUM
        ? incidentsByUrgency.MEDIUM.length
        : 0,
      lowUrgencyIncidents: incidentsByUrgency.LOW
        ? incidentsByUrgency.LOW.length
        : 0,

      categoryWithMostIncidents: {
        categoryName: mostPopularCategory,
        numberOfIncidents: mostPopularCategoryIncidents,
      },

      applicationWithMostIncidents: {
        numberOfIncidents: mostPopularAppIncidents,
        applicationName: mostPopularApp,
      },

      nodeWithMostIncidents: {
        numberOfIncidents: mostPopularNodeIncidents,
        nodeName: mostPopularNode,
      },
    });
  }, [incidents]);

  return {
    incidents,
    report: report,
    incidentStats,
    areIncidentsLoading,
    isReportLoading,
  };
};

export default useReportDetails;