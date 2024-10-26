package pl.pwr.zpi.metadata.dto.node;

public record Node(
        String clusterId,
        String name,
        boolean running) {

    public static Node of(NodeMetadata nodeMetadata) {
        return new Node("", nodeMetadata.name(), true);
    }
}
