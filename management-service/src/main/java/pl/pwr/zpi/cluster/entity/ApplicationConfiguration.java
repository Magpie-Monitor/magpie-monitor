package pl.pwr.zpi.cluster.entity;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import pl.pwr.zpi.cluster.dto.ApplicationConfigurationDTO;
import pl.pwr.zpi.reports.enums.Accuracy;

@Data
@Entity
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ApplicationConfiguration {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;
    private String name;
    private String kind;
    private Accuracy accuracy;
    private String customPrompt;

    public static ApplicationConfiguration fromApplicationConfigurationDTO(ApplicationConfigurationDTO applicationConfigurationDTO) {
        return ApplicationConfiguration.builder()
                .name(applicationConfigurationDTO.name())
                .kind(applicationConfigurationDTO.kind())
                .accuracy(applicationConfigurationDTO.accuracy())
                .customPrompt(applicationConfigurationDTO.customPrompt())
                .build();
    }
}
