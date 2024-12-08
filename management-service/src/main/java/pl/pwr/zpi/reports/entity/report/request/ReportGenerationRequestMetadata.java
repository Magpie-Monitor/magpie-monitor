package pl.pwr.zpi.reports.entity.report.request;

import jakarta.persistence.Column;
import lombok.Builder;
import lombok.Data;
import org.springframework.data.mongodb.core.mapping.MongoId;
import pl.pwr.zpi.reports.dto.event.ReportRequestFailed;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;
import pl.pwr.zpi.reports.enums.ReportType;

import java.util.List;

@Data
@Builder
public class ReportGenerationRequestMetadata {
    @MongoId
    @Column(name = "correlationId")
    private String correlationId;
    private ReportGenerationStatus status;
    private ReportRequestFailed error;
    private CreateReportRequest createReportRequest;
    private ReportType reportType;
    private long requestedAt;

    public static ReportGenerationRequestMetadata fromCreateReportRequest(
            String correlationId,
            CreateReportRequest createReportRequest,
            ReportType reportType
    ) {
        return ReportGenerationRequestMetadata.builder()
                .correlationId(correlationId)
                .status(ReportGenerationStatus.GENERATING)
                .createReportRequest(createReportRequest)
                .reportType(reportType)
                .requestedAt(System.currentTimeMillis())
                .build();
    }

    public String getClusterId() {
        return createReportRequest.clusterId();
    }

    public List<Long> getSlackReceiverIds() {
        return createReportRequest.slackReceiverIds();
    }

    public List<Long> getDiscordReceiverIds() {
        return createReportRequest.discordReceiverIds();
    }

    public List<Long> getMailReceiverIds() {
        return createReportRequest.emailReceiverIds();
    }
}
