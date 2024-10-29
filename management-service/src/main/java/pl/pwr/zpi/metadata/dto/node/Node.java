package pl.pwr.zpi.metadata.dto.node;

import java.io.Serializable;

public record Node(
        String name,
        boolean running) implements Serializable {
}
