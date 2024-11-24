package cluster

import pl.pwr.zpi.cluster.dto.UpdateClusterConfigurationRequest
import pl.pwr.zpi.cluster.entity.ClusterConfiguration
import pl.pwr.zpi.cluster.repository.ClusterRepository
import pl.pwr.zpi.cluster.service.ClusterService
import pl.pwr.zpi.notifications.ReceiverService
import pl.pwr.zpi.metadata.service.MetadataService
import pl.pwr.zpi.notifications.discord.entity.DiscordReceiver
import pl.pwr.zpi.notifications.email.entity.EmailReceiver
import pl.pwr.zpi.notifications.slack.entity.SlackReceiver
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadataDTO
import pl.pwr.zpi.reports.enums.Accuracy
import spock.lang.Specification
import spock.lang.Subject

class ClusterServiceTest extends Specification {

    @Subject
    ClusterService clusterService

    ClusterRepository clusterRepository = Mock()
    ReceiverService receiverService = Mock()
    MetadataService metadataService = Mock()

    def setup() {
        clusterService = new ClusterService(clusterRepository, receiverService, metadataService)
    }

    def "should update cluster configuration and save it to repository"() {
        given:
        def request = new UpdateClusterConfigurationRequest(
                "test-cluster-id", Accuracy.HIGH, true, 1000L,
                [1L, 2L], [3L, 4L], [5L, 6L], [], []
        )
        def clusterConfiguration = new ClusterConfiguration(id: "test-cluster-id")

        mockReceiverService()

        clusterRepository.save(_) >> clusterConfiguration

        when:
        def response = clusterService.updateClusterConfiguration(request)

        then:
        1 * clusterRepository.save(_)
        response.clusterId != null
    }

    def "should return cluster configuration by id"() {
        given:
        def clusterId = "test-cluster-id"
        def clusterConfiguration = new ClusterConfiguration(id: clusterId)
        def metadata = new ClusterMetadataDTO(clusterId, true)

        clusterConfiguration.nodeConfigurations = []
        clusterConfiguration.applicationConfigurations = []

        clusterRepository.findById(clusterId) >> Optional.of(clusterConfiguration)
        metadataService.getClusterById(clusterId) >> Optional.of(metadata)

        when:
        def result = clusterService.getClusterById(clusterId)

        then:
        result.isPresent()
        result.get().id == clusterConfiguration.id
        result.get().running == true
    }

    def "should return cluster configuration with false running status if metadata is empty"() {
        given:
        def clusterId = "test-cluster-id"
        def clusterConfiguration = new ClusterConfiguration(id: clusterId)

        clusterConfiguration.nodeConfigurations = []
        clusterConfiguration.applicationConfigurations = []

        clusterRepository.findById(clusterId) >> Optional.of(clusterConfiguration)
        metadataService.getClusterById(clusterId) >> Optional.empty()

        when:
        def result = clusterService.getClusterById(clusterId)

        then:
        result.isPresent()
        result.get().id == clusterConfiguration.id
        result.get().running == false
    }

    private void mockReceiverService() {
        receiverService.getSlackReceiverById(1L) >> new SlackReceiver()
        receiverService.getSlackReceiverById(2L) >> new SlackReceiver()
        receiverService.getDiscordReceiverById(3L) >> new DiscordReceiver()
        receiverService.getDiscordReceiverById(4L) >> new DiscordReceiver()
        receiverService.getEmailReceiverById(5L) >> new EmailReceiver()
        receiverService.getEmailReceiverById(6L) >> new EmailReceiver()
    }
}
