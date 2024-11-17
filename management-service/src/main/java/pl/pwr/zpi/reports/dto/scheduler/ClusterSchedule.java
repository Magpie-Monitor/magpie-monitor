package pl.pwr.zpi.reports.dto.scheduler;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import lombok.*;
import pl.pwr.zpi.reports.dto.request.CreateScheduleRequest;

@Data
@Entity
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ClusterSchedule {
    @Id
    private String clusterId;
    @NonNull
    private Long periodMs;
    @NonNull
    private Long lastGenerationMs;

    public static ClusterSchedule fromCreateScheduleRequest(CreateScheduleRequest scheduleRequest) {
        return ClusterSchedule.builder()
                .clusterId(scheduleRequest.clusterId())
                .periodMs(scheduleRequest.periodMs())
                .lastGenerationMs(System.currentTimeMillis())
                .build();
    }
}