package pl.pwr.zpi.reports.entity.report.application;

import lombok.Data;

import java.util.List;

@Data
public class ApplicationReport {
    private String applicationName;
    private String precision;
    private String customPrompt;
    private List<ApplicationIncident> incidents;
}
