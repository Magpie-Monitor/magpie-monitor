package pl.pwr.zpi.metadata.dto.application;

import java.util.Objects;

public record Application(String name, String kind, boolean running) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Application that = (Application) o;
        return Objects.equals(name, that.name) && Objects.equals(kind, that.kind);
    }

    @Override
    public int hashCode() {
        return Objects.hash(name, kind);
    }
}
