package pl.pwr.zpi.metadata.messaging.event.node;

import java.util.List;

public record NodeMetadata(String name, boolean running, List<String> files) {
}
