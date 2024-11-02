package pl.pwr.zpi.reports.entity.report.node;

import lombok.Data;

@Data
public class NodeIncidentSource {
    private Long timestamp;
    private String content;
    private String filename;
}
