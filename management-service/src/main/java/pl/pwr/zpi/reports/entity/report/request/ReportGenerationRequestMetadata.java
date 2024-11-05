package pl.pwr.zpi.reports.entity.report.request;

import jakarta.persistence.Id;
import lombok.Builder;
import lombok.Data;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;

import java.util.List;

@Data
@Builder
public class ReportGenerationRequestMetadata {
    @Id
    private String correlationId;
    private ReportGenerationStatus status;
    private CreateReportRequest createReportRequest;

    public static ReportGenerationRequestMetadata fromCreateReportRequest(
            String correlationId,
            CreateReportRequest createReportRequest
    ) {
        return ReportGenerationRequestMetadata.builder()
                .correlationId(correlationId)
                .status(ReportGenerationStatus.GENERATING)
                .createReportRequest(createReportRequest)
                .build();
    }

    public List<Long> getSlackReceiverIds() {
        return createReportRequest.slackReceiverIds();
    }

    public List<Long> getDiscordReceiverIds() {
        return createReportRequest.discordReceiverIds();
    }

    public List<Long> getMailReceiverIds() {
        return createReportRequest.mailReceiverIds();
    }
}
