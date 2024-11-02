package pl.pwr.zpi.reports.entity.report.application;

import lombok.Builder;
import lombok.Data;
import pl.pwr.zpi.reports.enums.Precision;

import java.util.List;

@Data
@Builder
public class ApplicationReport {
    private String applicationName;
    private Precision precision;
    private String customPrompt;
    private List<ApplicationIncident> incidents;
}
