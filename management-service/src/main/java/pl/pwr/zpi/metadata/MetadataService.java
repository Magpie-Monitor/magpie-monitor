package pl.pwr.zpi.metadata;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.utils.client.Client;
import pl.pwr.zpi.utils.client.HttpClient;

@Service
@RequiredArgsConstructor
public class MetadataService {

    private final String METADATA_SERVICE_URL = "http://localhost:9090/v1/metadata";
    private final HttpClient httpClient;

    public void getClusters() {

    }

    public void getNodes() {

    }

    public void getApplications() {

    }
}
