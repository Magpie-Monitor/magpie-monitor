package pl.pwr.zpi.metadata.dto.application;

import java.io.Serializable;

public record ApplicationMetadata(
        String name,
        String kind,
        boolean running) implements Serializable {
}
