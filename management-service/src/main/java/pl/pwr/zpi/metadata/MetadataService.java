package pl.pwr.zpi.metadata;

import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.dto.ApplicationMetadata;
import pl.pwr.zpi.metadata.dto.Cluster;
import pl.pwr.zpi.metadata.dto.NodeMetadata;
import pl.pwr.zpi.utils.client.HttpClient;

import java.util.List;
import java.util.Map;

@Service
@RequiredArgsConstructor
public class MetadataService {

    @Value("${metadata.base.url}")
    private String METADATA_SERVICE_BASE_URL;
    private final HttpClient httpClient;

    public List<ApplicationMetadata> getApplicationMetadata(String clusterId, Long sinceMillis, Long toMillis) {
        String url = String.format("%s/v1/metadata/clusters/%s/applications", METADATA_SERVICE_BASE_URL, clusterId);
        return httpClient.getList(
                url,
                Map.of(
                        "sinceMillis", sinceMillis.toString(),
                        "toMillis", toMillis.toString()
                ),
                ApplicationMetadata.class
        );
    }

    public List<NodeMetadata> getNodeMetadata(String clusterId, Long sinceMillis, Long toMillis) {
        String url = String.format("%s/v1/metadata/clusters/%s/nodes", METADATA_SERVICE_BASE_URL, clusterId);
        return httpClient.getList(
                url,
                Map.of(
                        "sinceMillis", sinceMillis.toString(),
                        "toMillis", toMillis.toString()
                ),
                NodeMetadata.class
        );
    }

    public List<Cluster> getClusters() {
        String url = String.format("%s/v1/metadata/clusters", METADATA_SERVICE_BASE_URL);
        return httpClient.getList(
                url,
                Map.of(),
                Cluster.class
        );
    }

}
