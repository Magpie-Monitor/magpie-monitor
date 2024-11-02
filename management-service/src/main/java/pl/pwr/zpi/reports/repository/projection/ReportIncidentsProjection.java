package pl.pwr.zpi.reports.repository.projection;

import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;

import java.util.List;

public interface ReportIncidentsProjection {

    interface ApplicationReportProjection {
        List<ApplicationIncident> getIncidents();
    }

    interface NodeReportProjection {
        List<NodeIncident> getIncidents();
    }

    List<ApplicationReportProjection> getApplicationReports();

    List<NodeReportProjection> getNodeReports();

    default List<ApplicationIncident> getApplicationIncidents() {
        return getApplicationReports().stream()
                .map(ApplicationReportProjection::getIncidents)
                .flatMap(List::stream)
                .toList();
    }

    default List<NodeIncident> getNodeIncidents() {
        return getNodeReports().stream()
                .map(NodeReportProjection::getIncidents)
                .flatMap(List::stream)
                .toList();
    }
}
