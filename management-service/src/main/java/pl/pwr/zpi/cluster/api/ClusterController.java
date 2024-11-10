package pl.pwr.zpi.cluster.api;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.cluster.dto.ClusterConfigurationDTO;
import pl.pwr.zpi.cluster.dto.UpdateClusterConfigurationRequest;
import pl.pwr.zpi.cluster.dto.UpdateClusterConfigurationResponse;
import pl.pwr.zpi.cluster.service.ClusterService;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadataDTO;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadataDTO;
import pl.pwr.zpi.metadata.dto.node.NodeMetadataDTO;
import pl.pwr.zpi.metadata.service.MetadataService;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/clusters")
public class ClusterController {

    private final MetadataService metadataService;
    private final ClusterService clusterService;

    @GetMapping
    public ResponseEntity<List<ClusterMetadataDTO>> getClusters() {
        return ResponseEntity.ok(metadataService.getAllClusters());
    }

    @GetMapping("/{id}/nodes")
    public ResponseEntity<List<NodeMetadataDTO>> getClusterNodes(@PathVariable String id) {
        return ResponseEntity.ok(metadataService.getClusterNodes(id));
    }

    @GetMapping("/{id}/applications")
    public ResponseEntity<List<ApplicationMetadataDTO>> getClusterApplications(@PathVariable String id) {
        return ResponseEntity.ok(metadataService.getClusterApplications(id));
    }

    @GetMapping("/{id}")
    public ResponseEntity<ClusterConfigurationDTO> getClusterById(@PathVariable String id) {
        return ResponseEntity.of(clusterService.getClusterById(id));
    }

    @PutMapping
    public ResponseEntity<UpdateClusterConfigurationResponse> updateClusterConfiguration(
            @RequestBody UpdateClusterConfigurationRequest configurationRequest) {
        return ResponseEntity.ok(clusterService.updateClusterConfiguration(configurationRequest));
    }
}
