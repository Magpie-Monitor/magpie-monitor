package pl.pwr.zpi.reports.entity.report.application;

import lombok.Data;

@Data
public class ApplicationIncidentSource {
    private Long timestamp;
    private String podName;
    private String containerName;
    private String image;
    private String content;
}
