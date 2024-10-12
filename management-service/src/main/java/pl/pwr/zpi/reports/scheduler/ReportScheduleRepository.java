package pl.pwr.zpi.reports.scheduler;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface ReportScheduleRepository extends JpaRepository<ReportSchedule, Long> {

    @Query(value = "SELECT * FROM report_schedule LIMIT 1", nativeQuery = true)
    Optional<ReportSchedule> getFirst();
}