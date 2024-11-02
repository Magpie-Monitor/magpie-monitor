package pl.pwr.zpi.reports.entity.report.node;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class NodeIncidentSource {
    private Long timestamp;
    private String content;
    private String filename;
}
