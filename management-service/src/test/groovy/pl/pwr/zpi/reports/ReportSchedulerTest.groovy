package pl.pwr.zpi.reports

import pl.pwr.zpi.reports.enums.Accuracy
import pl.pwr.zpi.reports.scheduler.ReportScheduler
import spock.lang.Specification
import pl.pwr.zpi.reports.dto.request.CreateReportRequest
import pl.pwr.zpi.reports.dto.scheduler.ReportSchedule
import pl.pwr.zpi.reports.enums.ReportType
import pl.pwr.zpi.reports.repository.ReportScheduleRepository
import pl.pwr.zpi.reports.service.ReportGenerationService
import pl.pwr.zpi.cluster.entity.ClusterConfiguration
import pl.pwr.zpi.cluster.repository.ClusterRepository

class ReportSchedulerTest extends Specification {

    def reportScheduleRepository = Mock(ReportScheduleRepository)
    def clusterRepository = Mock(ClusterRepository)
    def reportGenerationService = Mock(ReportGenerationService)

    def reportScheduler = new ReportScheduler(reportScheduleRepository, clusterRepository, reportGenerationService)

    def "should generate reports for all schedules"() {
        given:
        def schedule1 =createReportSchedule("cluster1", 1000L, 1000L)
        def schedule2 = createReportSchedule("cluster2", 2000L, 2000L)

        def clusterConfig1 = createClusterConfiguration("cluster1")
        def clusterConfig2 = createClusterConfiguration("cluster2")

        reportScheduleRepository.findAll() >> [schedule1, schedule2]
        clusterRepository.findById("cluster1") >> Optional.of(clusterConfig1)
        clusterRepository.findById("cluster2") >> Optional.of(clusterConfig2)

        when:
        reportScheduler.generateReports()

        then:
        2 * reportGenerationService.createReport(_ as CreateReportRequest, ReportType.SCHEDULED)
        2 * reportScheduleRepository.save(_ as ReportSchedule)
    }

    def "should process a schedule and generate report"() {
        given:
        def schedule = createReportSchedule("cluster1", 1000L, 1000L)
        def clusterConfig = createClusterConfiguration("cluster1")

        long nextGenerationTime = schedule.lastGenerationMs + schedule.periodMs
        CreateReportRequest reportRequest = CreateReportRequest.fromClusterConfiguration(clusterConfig, schedule.lastGenerationMs, nextGenerationTime)

        reportScheduleRepository.findAll() >> [schedule]
        clusterRepository.findById("cluster1") >> Optional.of(clusterConfig)

        when:
        reportScheduler.generateReports()

        then:
        1 * reportGenerationService.createReport(reportRequest, ReportType.SCHEDULED)
        1 * reportScheduleRepository.save(schedule)
        schedule.lastGenerationMs == nextGenerationTime
    }

    def "should throw exception when pl.pwr.zpi.cluster is not found"() {
        given:
        def schedule = createReportSchedule("cluster1", 1000L, 1000L)

        reportScheduleRepository.findAll() >> [schedule]
        clusterRepository.findById("cluster1") >> Optional.empty()

        when:
        reportScheduler.processSchedule(schedule)

        then:
        thrown(IllegalStateException)
    }

    def "should not generate report if next generation time is in the future"() {
        given:
        def futureTime = System.currentTimeMillis() + 10000L
        def schedule = createReportSchedule("cluster1", futureTime, 1000L)

        reportScheduleRepository.findAll() >> [schedule]

        when:
        reportScheduler.generateReports()

        then:
        0 * reportGenerationService.createReport(_, _)
    }

    private ReportSchedule createReportSchedule(String clusterId, long lastGenerationMs, long periodMs) {
        return ReportSchedule.builder()
                .clusterId(clusterId)
                .lastGenerationMs(lastGenerationMs)
                .periodMs(periodMs)
                .build()
    }

    private ClusterConfiguration createClusterConfiguration(String clusterId) {
        return ClusterConfiguration.builder()
                .id(clusterId)
                .accuracy(Accuracy.HIGH)
                .isEnabled(true)
                .generatedEveryMillis(2300000L)
                .slackReceivers([])
                .discordReceivers([])
                .emailReceivers([])
                .applicationConfigurations([])
                .nodeConfigurations([])
                .build()
    }
}
