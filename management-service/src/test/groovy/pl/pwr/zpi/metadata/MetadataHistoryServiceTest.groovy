package pl.pwr.zpi.metadata

import pl.pwr.zpi.metadata.dto.application.ApplicationMetadataDTO
import pl.pwr.zpi.metadata.dto.node.NodeMetadataDTO
import pl.pwr.zpi.metadata.entity.ClusterHistory
import pl.pwr.zpi.metadata.broker.dto.application.ApplicationMetadata
import pl.pwr.zpi.metadata.broker.dto.node.NodeMetadata
import pl.pwr.zpi.metadata.broker.dto.cluster.ClusterMetadata
import pl.pwr.zpi.metadata.repository.ClusterHistoryRepository
import pl.pwr.zpi.metadata.service.MetadataHistoryService
import spock.lang.Specification
import spock.lang.Subject
import spock.lang.Unroll

class MetadataHistoryServiceTest extends Specification {

    ClusterHistoryRepository clusterHistoryRepository
    @Subject
    MetadataHistoryService metadataHistoryService

    def setup() {
        clusterHistoryRepository = Mock()
        metadataHistoryService = new MetadataHistoryService(clusterHistoryRepository)
    }

    private Set<NodeMetadataDTO> createNodeHistory(String... nodeIds) {
        return nodeIds.collect { new NodeMetadataDTO(it, false) }.toSet()
    }

    private Set<ApplicationMetadataDTO> createAppHistory(String... appNames) {
        return appNames.collect { new ApplicationMetadataDTO(it, "type", false) }.toSet()
    }

    private ClusterHistory createClusterHistory(String clusterId, Set<NodeMetadataDTO> nodeHistory = [] as Set, Set<ApplicationMetadataDTO> appHistory = [] as Set) {
        return new ClusterHistory(clusterId, appHistory, nodeHistory)
    }

    @Unroll
    def "should return clusters history for #clusterId"() {
        given:
        def clusterHistoryList = [createClusterHistory(clusterId)]

        clusterHistoryRepository.findAll() >> clusterHistoryList

        when:
        def result = metadataHistoryService.getClustersHistory()

        then:
        result == clusterHistoryList

        where:
        clusterId  | _
        "cluster1" | _
        "cluster2" | _
    }


    def "should return node history for a pl.pwr.zpi.cluster"() {
        given:
        def clusterId = "cluster1"
        def nodeHistory = createNodeHistory("node1", "node2")
        def clusterHistory = createClusterHistory(clusterId, nodeHistory)

        clusterHistoryRepository.findById(clusterId) >> Optional.of(clusterHistory)

        when:
        def result = metadataHistoryService.getNodeHistory(clusterId)

        then:
        result == nodeHistory
    }

    def "should return empty node history if no history found for pl.pwr.zpi.cluster"() {
        given:
        def clusterId = "cluster1"
        clusterHistoryRepository.findById(clusterId) >> Optional.empty()

        when:
        def result = metadataHistoryService.getNodeHistory(clusterId)

        then:
        result.isEmpty()
    }

    def "should return application history for a pl.pwr.zpi.cluster"() {
        given:
        def clusterId = "cluster1"
        def appHistory = createAppHistory("app1", "app2")
        def clusterHistory = createClusterHistory(clusterId, [] as Set, appHistory)

        clusterHistoryRepository.findById(clusterId) >> Optional.of(clusterHistory)

        when:
        def result = metadataHistoryService.getApplicationHistory(clusterId)

        then:
        result == appHistory
    }

    def "should return empty application history if no history found for pl.pwr.zpi.cluster"() {
        given:
        def clusterId = "cluster1"
        clusterHistoryRepository.findById(clusterId) >> Optional.empty()

        when:
        def result = metadataHistoryService.getApplicationHistory(clusterId)

        then:
        result.isEmpty()
    }

    @Unroll
    def "should update clusters history if not already present for #clusterId"() {
        given:
        def clusterMetadata = new ClusterMetadata(clusterId)
        def clusterHistory = createClusterHistory(clusterId)

        clusterHistoryRepository.existsById(clusterId) >> false
        clusterHistoryRepository.save(_) >> clusterHistory

        when:
        metadataHistoryService.updateClustersHistory([clusterMetadata])

        then:
        1 * clusterHistoryRepository.save(_)

        where:
        clusterId << ["cluster1", "cluster2"]
    }

    @Unroll
    def "should not update clusters history if already present for #clusterId"() {
        given:
        def clusterMetadata = new ClusterMetadata(clusterId)
        clusterHistoryRepository.existsById(clusterId) >> true

        when:
        metadataHistoryService.updateClustersHistory([clusterMetadata])

        then:
        0 * clusterHistoryRepository.save(_)

        where:
        clusterId << ["cluster1"]
    }

    def "should update node history for a pl.pwr.zpi.cluster"() {
        given:
        def clusterId = "cluster1"
        def nodeMetadataList = List.of(new NodeMetadata("node1", List.of()), new NodeMetadata("node2", List.of()))
        def clusterHistory = createClusterHistory(clusterId)

        clusterHistoryRepository.findById(clusterId) >> Optional.of(clusterHistory)
        clusterHistoryRepository.save(_) >> clusterHistory

        when:
        metadataHistoryService.updateNodeHistory(clusterId, nodeMetadataList)

        then:
        1 * clusterHistoryRepository.save(_)
        clusterHistory.nodes().size() == 2
    }

    def "should add node history if pl.pwr.zpi.cluster does not exist"() {
        given:
        def clusterId = "cluster1"
        def nodeMetadataList = List.of(new NodeMetadata("node1", List.of()), new NodeMetadata("node2", List.of()))
        def clusterHistory = createClusterHistory(clusterId)

        clusterHistoryRepository.findById(clusterId) >> Optional.empty()
        clusterHistoryRepository.save(_) >> clusterHistory

        when:
        metadataHistoryService.updateNodeHistory(clusterId, nodeMetadataList)

        then:
        1 * clusterHistoryRepository.save(_)
    }

    def "should update application history for a pl.pwr.zpi.cluster"() {
        given:
        def clusterId = "cluster1"
        def appMetadataList = [new ApplicationMetadata("app1", "type1"), new ApplicationMetadata("app2", "type2")]
        def clusterHistory = createClusterHistory(clusterId)

        clusterHistoryRepository.findById(clusterId) >> Optional.of(clusterHistory)
        clusterHistoryRepository.save(_) >> clusterHistory

        when:
        metadataHistoryService.updateApplicationHistory(clusterId, appMetadataList)

        then:
        1 * clusterHistoryRepository.save(_)
        clusterHistory.applications().size() == 2
    }

    def "should add application history if pl.pwr.zpi.cluster does not exist"() {
        given:
        def clusterId = "cluster1"
        def appMetadataList = [new ApplicationMetadata("app1", "type1"), new ApplicationMetadata("app2", "type2")]
        def clusterHistory = createClusterHistory(clusterId)

        clusterHistoryRepository.findById(clusterId) >> Optional.empty()
        clusterHistoryRepository.save(_) >> clusterHistory

        when:
        metadataHistoryService.updateApplicationHistory(clusterId, appMetadataList)

        then:
        1 * clusterHistoryRepository.save(_)
    }
}
