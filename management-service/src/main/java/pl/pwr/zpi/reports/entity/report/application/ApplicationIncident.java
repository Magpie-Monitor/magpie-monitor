package pl.pwr.zpi.reports.entity.report.application;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.data.mongodb.core.mapping.MongoId;
import pl.pwr.zpi.reports.enums.Accuracy;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ApplicationIncident {
    @MongoId
    private String id;
    private String reportId;
    private String title;
    private Accuracy accuracy;
    private String customPrompt;
    private String clusterId;
    private String applicationName;
    private String category;
    private String summary;
    private String recommendation;
    private Urgency urgency;
    private List<ApplicationIncidentSource> sources;

    public void extendSourcesWithIncidentId() {
        sources.forEach(source -> source.setIncidentId(this.id));
    }
}
