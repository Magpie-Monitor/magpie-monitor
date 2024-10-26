package pl.pwr.zpi.metadata.dto.node;

import java.util.List;

public record NodeMetadata(String name, boolean running, List<String> files) {
}
