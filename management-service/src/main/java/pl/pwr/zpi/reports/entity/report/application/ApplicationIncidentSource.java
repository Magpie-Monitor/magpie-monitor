package pl.pwr.zpi.reports.entity.report.application;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class ApplicationIncidentSource {
    private Long timestamp;
    private String podName;
    private String containerName;
    private String image;
    private String content;
}
