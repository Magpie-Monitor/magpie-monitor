package pl.pwr.zpi.metadata;

import com.fasterxml.jackson.core.type.TypeReference;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.dto.ApplicationMetadata;
import pl.pwr.zpi.metadata.dto.Cluster;
import pl.pwr.zpi.metadata.dto.NodeMetadata;
import pl.pwr.zpi.utils.client.HttpClient;

import java.util.List;
import java.util.Map;

@Slf4j
@Service
@RequiredArgsConstructor
public class MetadataService {

    @Value("${metadata.base.url}")
    private String METADATA_SERVICE_BASE_URL;
    private final HttpClient httpClient;

    public List<ApplicationMetadata> getApplicationMetadata(String clusterName, Long sinceMillis, Long toMillis) {
        String url = String.format("%s/v1/metadata/clusters/%s/applications", METADATA_SERVICE_BASE_URL, clusterName);
        return httpClient.getList(
                url,
                Map.of(
                        "sinceMillis", sinceMillis.toString(),
                        "toMillis", toMillis.toString()
                ),
                new TypeReference<>() {
                }
        );
    }

    public List<NodeMetadata> getNodeMetadata(String clusterName, Long sinceMillis, Long toMillis) {
        String url = String.format("%s/v1/metadata/clusters/%s/nodes", METADATA_SERVICE_BASE_URL, clusterName);
        return httpClient.getList(
                url,
                Map.of(
                        "sinceMillis", sinceMillis.toString(),
                        "toMillis", toMillis.toString()
                ),
                new TypeReference<>() {
                }
        );
    }

    public List<Cluster> getClusters() {
        String url = String.format("%s/v1/metadata/clusters", METADATA_SERVICE_BASE_URL);
        return httpClient.getList(
                url,
                Map.of(),
                new TypeReference<>() {
                }
        );
    }

}
