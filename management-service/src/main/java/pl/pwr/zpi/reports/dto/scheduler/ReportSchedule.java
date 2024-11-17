package pl.pwr.zpi.reports.dto.scheduler;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import lombok.*;
import pl.pwr.zpi.reports.dto.request.CreateReportScheduleRequest;

@Data
@Entity
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ReportSchedule {
    @Id
    private String clusterId;
    @NonNull
    private Long periodMs;
    @NonNull
    private Long lastGenerationMs;

    public static ReportSchedule fromCreateScheduleRequest(CreateReportScheduleRequest scheduleRequest) {
        return ReportSchedule.builder()
                .clusterId(scheduleRequest.clusterId())
                .periodMs(scheduleRequest.periodMs())
                .lastGenerationMs(System.currentTimeMillis())
                .build();
    }
}