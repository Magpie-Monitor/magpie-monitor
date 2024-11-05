package pl.pwr.zpi.reports.entity.report.node;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class NodeIncidentSource {
    private Long timestamp;
    private String content;
    private String filename;
}
