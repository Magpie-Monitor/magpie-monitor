package pl.pwr.zpi.metadata.dto.node;

import java.util.Objects;

public record Node(String name, boolean running) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Node node = (Node) o;
        return Objects.equals(name, node.name);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(name);
    }
}
