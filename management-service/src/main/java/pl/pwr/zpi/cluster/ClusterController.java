package pl.pwr.zpi.cluster;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/clusters/asda")
public class ClusterController {

    private final ClusterService clusterService;

    @GetMapping
    public void getClusters() {

    }


    @GetMapping("/{id}")
    public ResponseEntity<ClusterConfiguration> getClusterById(@PathVariable String id) {
        return null;
    }
}
