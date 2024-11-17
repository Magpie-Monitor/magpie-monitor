package pl.pwr.zpi.reports.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import pl.pwr.zpi.reports.dto.scheduler.ReportSchedule;

public interface ReportScheduleRepository extends JpaRepository<ReportSchedule, String> {
}
