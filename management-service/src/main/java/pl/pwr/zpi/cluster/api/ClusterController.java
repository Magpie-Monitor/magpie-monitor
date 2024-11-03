package pl.pwr.zpi.cluster.api;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
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

    @GetMapping
    public ResponseEntity<List<Cluster>> getClusters() {
        return ResponseEntity.ok(metadataService.getAllClusters());
    }

    @GetMapping("/{id}")
    public ResponseEntity<Cluster> getClusterById(@PathVariable String id) {
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
}
