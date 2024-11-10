package pl.pwr.zpi.cluster.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import pl.pwr.zpi.cluster.entity.Cluster;

public interface ClusterRepository extends JpaRepository<Cluster, String> {
}
