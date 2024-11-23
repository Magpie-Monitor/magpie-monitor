package repositories

import (
	"context"
	"encoding/json"
	"github.com/IBM/fp-go/array"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"time"
)

// TODO: To be clarified once the contract from agent is agreed upon
type ApplicationLogs struct {
	ClusterId     string     `json:"clusterId"`
	Kind          string     `json:"kind"`
	CollectedAtMs int64      `json:"collectedAtMs"`
	Name          string     `json:"name"`
	Namespace     string     `json:"namespace"`
	Pods          []*PodLogs `json:"pods"`
}

type PodLogs struct {
	Name       string           `json:"name"`
	Containers []*ContainerLogs `json:"containers"`
}

type ContainerLogs struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Content string `json:"content"`
}

func (l *ApplicationLogsDocument) GetContent() *string {
	return &l.Content
}

type ApplicationLogsDocument struct {
	Id              string `json:"_id,omitempty" bson:"id,omitempty"`
	ClusterId       string `json:"clusterId" bson:"clusterId"`
	Kind            string `json:"kind" bson:"kind"`
	CollectedAtMs   int64  `json:"collectedAtMs"`
	ApplicationName string `json:"applicationName" bson:"applicationName"`
	Namespace       string `json:"namespace" bson:"namespace"`
	PodName         string `json:"podName" bson:"podName"`
	ContainerName   string `json:"containerName" bson:"containerName"`
	Image           string `json:"image" bson:"image"`
	Content         string `json:"content" bson:"content"`
}

func (l *ApplicationLogs) Flatten() []*ApplicationLogsDocument {
	var documents []*ApplicationLogsDocument
	for _, pod := range l.Pods {
		for _, container := range pod.Containers {
			documents = append(documents, &ApplicationLogsDocument{
				ClusterId:       l.ClusterId,
				Kind:            l.Kind,
				CollectedAtMs:   l.CollectedAtMs,
				Namespace:       l.Namespace,
				PodName:         pod.Name,
				ContainerName:   container.Name,
				Image:           container.Image,
				Content:         container.Content,
				ApplicationName: l.Name,
			})
		}
	}

	return documents
}

type ApplicationLogsRepository interface {
	CreateIndex(ctx context.Context, indexName string) error
	GetLogs(ctx context.Context, cluster string, startDate time.Time, endDate time.Time) ([]*ApplicationLogsDocument, error)
	InsertLogs(ctx context.Context, logs *ApplicationLogs) error
	RemoveIndex(ctx context.Context, indexName string) error
	GetBatchedLogs(ctx context.Context, cluster string, startDate time.Time, endDate time.Time) (ApplicationLogsBatchRetriever, error)
	GetLogsByIds(ctx context.Context, clusterId string, startDate time.Time, endDate time.Time, ids []string) ([]*ApplicationLogsDocument, error)
}

func ProvideAsApplicationLogsRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ApplicationLogsRepository)),
	)

}

type ElasticSearchApplicationLogsRepository struct {
	esClient *es.TypedClient
	logger   *zap.Logger
	indices  map[string]bool
}

func (r *ElasticSearchApplicationLogsRepository) doesIndexExists(index string) bool {
	return r.indices[index]
}

func (r *ElasticSearchApplicationLogsRepository) GetLogsByIds(
	ctx context.Context,
	clusterId string,
	startDate time.Time,
	endDate time.Time,
	ids []string,
) ([]*ApplicationLogsDocument, error) {

	indices := r.getIndiciesWithClusterAndDateRange(clusterId, startDate, endDate)

	applicationLogsById, err :=
		elasticsearch.GetAndMapDocumentsByIds[*ApplicationLogsDocument](ctx, r.esClient, indices, ids, r.logger)

	if err != nil {
		r.logger.Error("Failed to fetch and map document ids", zap.Error(err), zap.Any("ids", ids))
		return nil, err
	}

	applicationLogs := make([]*ApplicationLogsDocument, 0, len(applicationLogsById))
	for id, log := range applicationLogsById {
		log.Id = id
		applicationLogs = append(applicationLogs, log)
	}

	return applicationLogs, nil
}

func getApplicationLogsIndexName(applicationLogs *ApplicationLogs) string {

	val := elasticsearch.GetIndexName(applicationLogs.ClusterId, "applications", applicationLogs.CollectedAtMs)
	return val
}

func (r *ElasticSearchApplicationLogsRepository) getIndiciesWithClusterAndDateRange(cluster string,
	fromDate time.Time,
	toDate time.Time) []string {

	filter := array.Filter(
		elasticsearch.FilterIndicesByClusterAndDateRange(cluster, "applications", fromDate, toDate))

	r.updateIndices()

	return filter(maps.Keys(r.indices))
}

func (r *ElasticSearchApplicationLogsRepository) GetLogs(ctx context.Context, cluster string, startDate time.Time, endDate time.Time) ([]*ApplicationLogsDocument, error) {

	applicationLogs := make([]*ApplicationLogsDocument, 0)

	batch, err := r.GetBatchedLogs(ctx, cluster, startDate, endDate)
	if err != nil {
		r.logger.Error("Failed to fetch application logs", zap.Error(err))
		return nil, err
	}

	for {
		if batch.HasNextBatch() {
			nextBatch, err := batch.GetNextBatch()
			if err != nil {
				r.logger.Error("Failed to get next batch of application logs")
				return nil, err
			}

			applicationLogs = append(applicationLogs, nextBatch...)
		} else {
			return applicationLogs, nil
		}
	}
}

func (r *ElasticSearchApplicationLogsRepository) GetBatchedLogs(
	ctx context.Context,
	cluster string,
	startDate time.Time,
	endDate time.Time) (ApplicationLogsBatchRetriever, error) {

	indices := r.getIndiciesWithClusterAndDateRange(cluster, startDate, endDate)
	query := elasticsearch.GetQueryByTimestamps(startDate, endDate)

	elasticDocumentRetriever, err := NewElasticBatchedDocumentsRetriever(
		r.esClient,
		elasticsearch.GetSearchQuery(ctx, r.esClient, indices, query, 10000),
	)
	if err != nil {
		r.logger.Error("Failed to get application logs", zap.Error(err))
		return nil, err
	}

	return &ElasticApplicationBatchedLogsRetriever{
		documentRetriever: elasticDocumentRetriever,
	}, nil
}

func (r *ElasticSearchApplicationLogsRepository) CreateIndex(ctx context.Context, indexName string) error {

	_, err := r.esClient.Indices.Create(indexName).Do(ctx)

	if err != nil {
		r.logger.Error("Failed to create an applicationLogs index", zap.Error(err))
		return err
	}

	err = r.updateIndices()
	if err != nil {
		r.logger.Error("Failed to update indices list", zap.Error(err))
	}

	return nil
}

func (r *ElasticSearchApplicationLogsRepository) RemoveIndex(ctx context.Context, indexName string) error {

	_, err := r.esClient.Indices.Delete(indexName).Do(ctx)

	if err != nil {
		r.logger.Error("Failed to delete an applicationLogs index", zap.Error(err))
		return err
	}

	err = r.updateIndices()
	if err != nil {
		r.logger.Error("Failed to update indices list", zap.Error(err))
	}

	return nil
}

func (r *ElasticSearchApplicationLogsRepository) InsertLogs(ctx context.Context, logs *ApplicationLogs) error {

	index := getApplicationLogsIndexName(logs)

	if !r.doesIndexExists(index) {
		r.CreateIndex(ctx, index)
	}

	bulk := r.esClient.Bulk().Index(index)

	for _, log := range logs.Flatten() {

		if log.Content == "" {

			//Skip inserting logs without content
			continue
		}

		jsonLog, err := json.Marshal(log)
		r.logger.Sugar().Infof("%s", jsonLog)
		if err != nil {
			r.logger.Error("Failed to convert log to json", zap.Error(err))
		}
		bulk.IndexOp(*types.NewIndexOperation(), jsonLog)
	}

	_, err := bulk.Do(ctx)

	if err != nil {
		r.logger.Error("Failed to insert application logs", zap.Error(err))
	}
	return nil
}

func GroupApplicationLogsByName(logs []*ApplicationLogsDocument) map[string][]*ApplicationLogsDocument {
	groupedLogs := make(map[string][]*ApplicationLogsDocument)
	for _, log := range logs {
		groupedLogs[log.ApplicationName] = append(groupedLogs[log.ApplicationName], log)
	}

	return groupedLogs
}

type ApplicationLogsParams struct {
	fx.In
	ElasticSearchClient *es.TypedClient
	Logger              *zap.Logger
}

func NewElasticSearchApplicationLogsRepository(p ApplicationLogsParams) *ElasticSearchApplicationLogsRepository {

	repository := ElasticSearchApplicationLogsRepository{
		esClient: p.ElasticSearchClient,
		logger:   p.Logger,
	}

	err := repository.updateIndices()
	if err != nil {
		repository.logger.Error("Failed to create fetch initial logsdb indices", zap.Error(err))
	}

	return &repository
}

func (r *ElasticSearchApplicationLogsRepository) updateIndices() error {

	indices, err := elasticsearch.GetAllIndicesSet(r.esClient)
	if err != nil {
		r.logger.Error("Failed to fetch logsdb indices", zap.Error(err))
		return err
	}

	r.indices = indices
	r.logger.Info("Fetched logsdb indices: ", zap.Any("indices", indices))
	return nil
}

// Compile-time check if ElasticSearchApplicationLogsRepository implements
// the ApplicationLogsRepository interface
var _ ApplicationLogsRepository = &ElasticSearchApplicationLogsRepository{}
