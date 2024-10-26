package pl.pwr.zpi.metadata;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadata;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadata;
import pl.pwr.zpi.metadata.dto.node.Node;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1")
public class TestController {

    private final MetadataService metadataService;

    @GetMapping("/apps")
    public ResponseEntity<List<ApplicationMetadata>> getReportById() {
        return ResponseEntity.ok().body(metadataService.getClusterApplications("local-docker"));
    }

    @GetMapping("/nodes")
    public ResponseEntity<List<Node>> get() {
        return ResponseEntity.ok().body(metadataService.getClusterNodes("local-docker"));
    }

    @GetMapping("/clusters")
    public ResponseEntity<List<ClusterMetadata>> getClusters() {
        return ResponseEntity.ok(metadataService.getClusters());
    }


    @GetMapping("/cluster")
    public ResponseEntity<ClusterMetadata> getCluster() {
        return ResponseEntity.of(metadataService.getClusterById("bccec-1aaa"));
    }
}
