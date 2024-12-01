package pl.pwr.zpi.metadata

import pl.pwr.zpi.metadata.broker.dto.application.AggregatedApplicationMetadata
import pl.pwr.zpi.metadata.broker.dto.application.ApplicationMetadata
import pl.pwr.zpi.metadata.broker.dto.cluster.AggregatedClusterMetadata
import pl.pwr.zpi.metadata.broker.dto.cluster.ClusterMetadata
import pl.pwr.zpi.metadata.broker.dto.node.AggregatedNodeMetadata
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadataDTO
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadataDTO
import pl.pwr.zpi.metadata.dto.node.NodeMetadataDTO
import pl.pwr.zpi.metadata.entity.ClusterHistory
import pl.pwr.zpi.metadata.repository.AggregatedApplicationMetadataRepository
import pl.pwr.zpi.metadata.repository.AggregatedClusterMetadataRepository
import pl.pwr.zpi.metadata.repository.AggregatedNodeMetadataRepository
import pl.pwr.zpi.metadata.service.MetadataHistoryService
import pl.pwr.zpi.metadata.service.MetadataService
import spock.lang.Specification
import spock.lang.Subject

class MetadataServiceTest extends Specification {

    AggregatedApplicationMetadataRepository applicationMetadataRepository
    AggregatedNodeMetadataRepository nodeMetadataRepository
    AggregatedClusterMetadataRepository clusterMetadataRepository
    MetadataHistoryService metadataHistoryService

    @Subject
    MetadataService metadataService

    def setup() {
        applicationMetadataRepository = Mock()
        nodeMetadataRepository = Mock()
        clusterMetadataRepository = Mock()
        metadataHistoryService = Mock()
        metadataService = new MetadataService(applicationMetadataRepository, nodeMetadataRepository, clusterMetadataRepository, metadataHistoryService)
    }

    private AggregatedClusterMetadata createAggregatedClusterMetadata(Long collectedAt) {
        return new AggregatedClusterMetadata(collectedAt, List.of())
    }

    private AggregatedNodeMetadata createAggregatedNodeMetadata(Long collectedAt, String nodeName) {
        return new AggregatedNodeMetadata(collectedAt, nodeName, List.of())
    }

    private AggregatedApplicationMetadata createAggregatedApplicationMetadata(Long collectedAt, String clusterId) {
        return new AggregatedApplicationMetadata(collectedAt, clusterId, List.of())
    }

    private ClusterMetadataDTO createClusterMetadataDTO(String clusterId) {
        return new ClusterMetadataDTO(clusterId, true)
    }

    private ClusterMetadataDTO createInactiveClusterMetadataDTO(String clusterId) {
        return new ClusterMetadataDTO(clusterId, false)
    }

    private NodeMetadataDTO createNodeMetadataDTO(String nodeName) {
        return new NodeMetadataDTO(nodeName, true)
    }

    private ApplicationMetadataDTO createApplicationMetadataDTO(String appName) {
        return new ApplicationMetadataDTO(appName, "deployment", true)
    }

    private ApplicationMetadataDTO createInactiveApplicationMetadataDTO(String appName) {
        return new ApplicationMetadataDTO(appName, "deployment", false)
    }

    def "should return all clusters including inactive ones"() {
        given:
        clusterMetadataRepository.findFirstByOrderByCollectedAtMsDesc() >> Optional.of(new AggregatedClusterMetadata(1732395048724L, [
                new ClusterMetadata("cluster1")
        ]))

        metadataHistoryService.getClustersHistory() >> List.of(
                new ClusterHistory("cluster2", Set.of(), Set.of())
        )

        when:
        def result = metadataService.getAllClusters()

        then:
        result.size() == 2
        result.contains(createClusterMetadataDTO("cluster1"))
        result.contains(createInactiveClusterMetadataDTO("cluster2"))
    }

    def "should return pl.pwr.zpi.cluster by id"() {
        given:
        def clusterId = "cluster1"
        def clusterDTO = createInactiveClusterMetadataDTO(clusterId)

        clusterMetadataRepository.existsByMetadataClusterId(clusterId) >> true
        clusterMetadataRepository.findFirstByOrderByCollectedAtMsDesc() >> Optional.of(new AggregatedClusterMetadata(1732395048724L, List.of(new ClusterMetadata(clusterId))))
        metadataService.getActiveClusters() >> List.of(clusterDTO)

        when:
        def result = metadataService.getClusterById(clusterId)

        then:
        result.isPresent()
        result.get() == clusterDTO
    }

    def "should return empty if pl.pwr.zpi.cluster does not exist"() {
        given:
        def clusterId = "cluster1"

        clusterMetadataRepository.existsByMetadataClusterId(clusterId) >> false

        when:
        def result = metadataService.getClusterById(clusterId)

        then:
        result.isEmpty()
    }

    def "should return all applications for a pl.pwr.zpi.cluster including inactive ones"() {
        given:
        def clusterId = "cluster1"
        def activeApp = createApplicationMetadataDTO("app1")
        def inactiveApp = createInactiveApplicationMetadataDTO("app2")

        applicationMetadataRepository.findFirstByClusterIdOrderByCollectedAtMsDesc(clusterId) >> Optional.of(
                new AggregatedApplicationMetadata(1732395048724L, "test123", List.of(new ApplicationMetadata("app1", "deployment")))
        )

        metadataHistoryService.getApplicationHistory(clusterId) >> List.of(
                new ApplicationMetadataDTO("app2", "deployment", false)
        )

        when:
        def result = metadataService.getClusterApplications(clusterId)

        then:
        result.size() == 2
        result.contains(activeApp)
        result.contains(inactiveApp)
    }



    def "should save pl.pwr.zpi.cluster pl.pwr.zpi.metadata"() {
        given:
        def clusterMetadata = createAggregatedClusterMetadata(1732395048724L)

        when:
        metadataService.saveClusterMetadata(clusterMetadata)

        then:
        1 * clusterMetadataRepository.save(clusterMetadata)
    }

    def "should save application pl.pwr.zpi.metadata"() {
        given:
        def applicationMetadata = createAggregatedApplicationMetadata(1732395048724L, "test123")

        when:
        metadataService.saveApplicationMetadata(applicationMetadata)

        then:
        1 * applicationMetadataRepository.save(applicationMetadata)
    }

    def "should save node pl.pwr.zpi.metadata"() {
        given:
        def nodeMetadata = createAggregatedNodeMetadata(1732395048724L, "test123")

        when:
        metadataService.saveNodeMetadata(nodeMetadata)

        then:
        1 * nodeMetadataRepository.save(nodeMetadata)
    }
}