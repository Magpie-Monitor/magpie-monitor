package pl.pwr.zpi.reports.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import pl.pwr.zpi.reports.dto.scheduler.ClusterSchedule;

public interface SchedulerRepository extends JpaRepository<ClusterSchedule, String> {
}
