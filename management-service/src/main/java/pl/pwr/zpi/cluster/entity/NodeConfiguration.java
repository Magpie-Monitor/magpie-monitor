package pl.pwr.zpi.cluster.entity;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import pl.pwr.zpi.cluster.dto.NodeConfigurationDTO;
import pl.pwr.zpi.reports.enums.Accuracy;

@Data
@Entity
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class NodeConfiguration {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;
    private String name;
    private Accuracy accuracy;
    private String customPrompt;

    public static NodeConfiguration fromNodeConfigurationDTO(NodeConfigurationDTO nodeConfigurationDTO) {
        return NodeConfiguration.builder()
                .name(nodeConfigurationDTO.name())
                .accuracy(nodeConfigurationDTO.accuracy())
                .customPrompt(nodeConfigurationDTO.customPrompt())
                .build();
    }
}
