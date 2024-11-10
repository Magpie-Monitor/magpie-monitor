package pl.pwr.zpi.cluster.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import pl.pwr.zpi.cluster.entity.ClusterConfiguration;

public interface ClusterRepository extends JpaRepository<ClusterConfiguration, String> {
}
