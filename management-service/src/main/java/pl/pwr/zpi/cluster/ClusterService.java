package pl.pwr.zpi.cluster;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.MetadataService;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadata;

import java.util.List;
import java.util.Optional;

@RequiredArgsConstructor
@Service
public class ClusterService {

    private final MetadataService metadataService;

    public ClusterConfiguration getClusterConfigurationById(String clusterId) {
        return ClusterConfiguration
                .builder()
                .id(clusterId)
                .precision("")
                .running(metadataService.getClusterById(clusterId).isPresent())
                .build();
    }

//    public List<ClusterConfiguration> getClusterConfigurations() {
//    }
}
