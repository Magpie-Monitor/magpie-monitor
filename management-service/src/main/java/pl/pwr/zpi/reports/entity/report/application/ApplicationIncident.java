package pl.pwr.zpi.reports.entity.report.application;

import lombok.Data;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Data
public class ApplicationIncident {
    private String id;
    private String clusterId;
    private String applicationName;
    private String category;
    private String summary;
    private String recommendation;
    private Urgency urgency;
    private List<ApplicationIncidentSource> sources;
}
