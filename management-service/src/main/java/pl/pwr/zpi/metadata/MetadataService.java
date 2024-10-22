package pl.pwr.zpi.metadata;

import com.fasterxml.jackson.core.type.TypeReference;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.metadata.dto.ApplicationMetadata;
import pl.pwr.zpi.metadata.dto.ClusterMetadata;
import pl.pwr.zpi.metadata.dto.NodeMetadata;
import pl.pwr.zpi.metadata.messaging.event.application.AggregatedApplicationMetadata;
import pl.pwr.zpi.metadata.messaging.event.node.AggregatedNodeMetadata;
import pl.pwr.zpi.metadata.repository.AggregatedApplicationMetadataRepository;
import pl.pwr.zpi.metadata.repository.AggregatedNodeMetadataRepository;
import pl.pwr.zpi.utils.client.HttpClient;

import java.util.Collections;
import java.util.List;
import java.util.Map;

@Slf4j
@Service
@RequiredArgsConstructor
public class MetadataService {

    @Value("${metadata.base.url}")
    private String METADATA_SERVICE_BASE_URL;
    private final HttpClient httpClient;

    private final AggregatedApplicationMetadataRepository applicationMetadataRepository;
    private final AggregatedNodeMetadataRepository nodeMetadataRepository;

    public void saveApplicationMetadata(AggregatedApplicationMetadata applicationMetadata) {
        applicationMetadataRepository.save(applicationMetadata);
    }

    public void saveNodeMetadata(AggregatedNodeMetadata nodeMetadata) {
        nodeMetadataRepository.save(nodeMetadata);
    }

    @Deprecated
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

    @Deprecated
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

    @Deprecated
    public List<ClusterMetadata> getClusters() {
        String url = String.format("%s/v1/metadata/clusters", METADATA_SERVICE_BASE_URL);
        return httpClient.getList(
                url,
                Collections.emptyMap(),
                new TypeReference<>() {
                }
        );
    }
}
