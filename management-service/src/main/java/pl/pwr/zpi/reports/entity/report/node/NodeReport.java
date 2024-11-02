package pl.pwr.zpi.reports.entity.report.node;

import lombok.Data;

import java.util.List;

@Data
public class NodeReport {
    private String node;
    private String precision;
    private String customPrompt;
    private List<NodeReport> nodeIncidents;
}
