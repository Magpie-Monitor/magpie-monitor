package pl.pwr.zpi.metadata.broker.dto.node;

import java.util.List;

public record NodeMetadata(String name, List<String> files) {
}