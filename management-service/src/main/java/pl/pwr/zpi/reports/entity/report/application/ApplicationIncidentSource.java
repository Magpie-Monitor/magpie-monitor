package pl.pwr.zpi.reports.entity.report.application;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ApplicationIncidentSource {
    private String incidentId;
    private Long timestamp;
    private String podName;
    private String containerName;
    private String image;
    private String content;
}
