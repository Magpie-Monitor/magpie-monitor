package pl.pwr.zpi.metadata;


import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/clusters")
public class MetadataController {

    @GetMapping
    public void getClusters() {

    }

    @GetMapping("/{id}")
    public void getClusterById(@PathVariable String id) {

    }
}
