package pl.pwr.zpi.metadata.event.dto.application;

import java.io.Serializable;

public record ApplicationMetadata(String name, String kind) implements Serializable {
}
