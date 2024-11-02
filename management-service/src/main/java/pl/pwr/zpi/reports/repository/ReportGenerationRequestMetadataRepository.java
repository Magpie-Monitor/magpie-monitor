package pl.pwr.zpi.reports.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata;
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;

import java.util.List;
import java.util.Optional;

public interface ReportGenerationRequestMetadataRepository extends MongoRepository<ReportGenerationRequestMetadata, String> {
    Optional<ReportGenerationRequestMetadata> findByCorrelationId(String correlationId);

    List<ReportGenerationRequestMetadata> findByStatus(ReportGenerationStatus generationStatus);
}
