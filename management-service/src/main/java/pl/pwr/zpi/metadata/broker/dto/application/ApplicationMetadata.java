package pl.pwr.zpi.metadata.broker.dto.application;

import java.io.Serializable;

public record ApplicationMetadata(String name, String kind) implements Serializable {
}
