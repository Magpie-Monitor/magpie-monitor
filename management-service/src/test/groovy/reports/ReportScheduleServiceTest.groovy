package reports

import lombok.extern.slf4j.Slf4j
import pl.pwr.zpi.reports.dto.request.CreateReportScheduleRequest
import pl.pwr.zpi.reports.dto.scheduler.ReportSchedule
import pl.pwr.zpi.reports.repository.ReportScheduleRepository
import pl.pwr.zpi.cluster.repository.ClusterRepository
import pl.pwr.zpi.reports.service.ReportScheduleService
import spock.lang.Ignore
import spock.lang.Specification

class ReportScheduleServiceTest extends Specification {

    def clusterRepository = Mock(ClusterRepository)
    def reportScheduleRepository = Mock(ReportScheduleRepository)
    def reportScheduleService = new ReportScheduleService(reportScheduleRepository, clusterRepository)

    @Ignore
    def "should schedule report when cluster exists"() {
        given:
        def clusterId = "cluster123"
        def scheduleRequest = new CreateReportScheduleRequest(clusterId, 86400000L)
        clusterRepository.existsById(clusterId) >> true
        1 * reportScheduleRepository.save(_ as ReportSchedule)

        when:
        reportScheduleService.scheduleReport(scheduleRequest)

        then:
        1 * clusterRepository.existsById(clusterId)
        noExceptionThrown()
    }

    def "should throw exception when cluster does not exist"() {
        given:
        def clusterId = "cluster123"
        def scheduleRequest = new CreateReportScheduleRequest(clusterId, 86400000L)
        clusterRepository.existsById(clusterId) >> false

        when:
        reportScheduleService.scheduleReport(scheduleRequest)

        then:
        1 * clusterRepository.existsById(clusterId)
        thrown(IllegalArgumentException)
    }
}