package pl.pwr.zpi.cluster.api;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.cluster.dto.ClusterConfigurationDTO;
import pl.pwr.zpi.cluster.dto.ClusterConfigurationRequest;
import pl.pwr.zpi.cluster.dto.ClusterIdResponse;
import pl.pwr.zpi.cluster.service.ClusterService;
import pl.pwr.zpi.metadata.dto.application.Application;
import pl.pwr.zpi.metadata.dto.cluster.Cluster;
import pl.pwr.zpi.metadata.dto.node.Node;
import pl.pwr.zpi.metadata.service.MetadataService;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/clusters")
public class ClusterController {

    private final MetadataService metadataService;
    private final ClusterService clusterService;

    // TODO - return receivers
    @GetMapping
    public ResponseEntity<List<Cluster>> getClusters() {
        return ResponseEntity.ok(metadataService.getAllClusters());
    }

//    @GetMapping("/{id}")
//    public ResponseEntity<Cluster> getClusterById(@PathVariable String id) {
//        return ResponseEntity.of(metadataService.getClusterById(id));
//    }

    @GetMapping("/{id}/summary")
    public ResponseEntity<Cluster> getClusterSummaryById(@PathVariable String id) {
        return ResponseEntity.of(metadataService.getClusterById(id));
    }

    @GetMapping("/{id}/nodes")
    public ResponseEntity<List<Node>> getClusterNodes(@PathVariable String id) {
        return ResponseEntity.ok(metadataService.getClusterNodes(id));
    }

    @GetMapping("/{id}/applications")
    public ResponseEntity<List<Application>> getClusterApplications(@PathVariable String id) {
        return ResponseEntity.ok(metadataService.getClusterApplications(id));
    }

    @GetMapping("/{id}")
    public ResponseEntity<ClusterConfigurationDTO> getClusterById(@PathVariable String id) {
        return ResponseEntity.of(clusterService.getClusterById(id));
    }

    @PutMapping
    public ResponseEntity<ClusterIdResponse> updateClusterConfiguration(@RequestBody ClusterConfigurationRequest configurationRequest) {
        return ResponseEntity.ok(clusterService.updateClusterConfiguration(configurationRequest));
    }
}
