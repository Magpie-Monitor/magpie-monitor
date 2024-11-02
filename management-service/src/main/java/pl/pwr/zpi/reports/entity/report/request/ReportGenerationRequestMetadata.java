package pl.pwr.zpi.reports.entity.report.request;

import lombok.Builder;
import lombok.Data;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;

import java.util.List;

@Data
@Builder
public class ReportGenerationRequestMetadata {
    private String correlationId;
    private ReportGenerationStatus status;
    private CreateReportRequest createReportRequest;

    public static ReportGenerationRequestMetadata fromCreateReportRequest(
            String correlationId,
            CreateReportRequest createReportRequest
    ) {
        return ReportGenerationRequestMetadata.builder()
                .correlationId(correlationId)
                .status(ReportGenerationStatus.GENERATED)
                .createReportRequest(createReportRequest)
                .build();
    }

    public void markAsFailed() {
        this.status = ReportGenerationStatus.ERROR;
    }

    public void markAsGenerated() {
        this.status = ReportGenerationStatus.GENERATED;
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
