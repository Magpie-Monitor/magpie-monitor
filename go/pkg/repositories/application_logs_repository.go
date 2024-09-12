package repositories

import (
	"context"
	"encoding/json"
	"github.com/IBM/fp-go/array"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"time"
)

// TODO: To be clarified once the contract from agent is agreed upon
type ApplicationLogs struct {
	Cluster   string
	Kind      string
	Timestamp int64
	Name      string
	Namespace string
	Pods      []*PodLogs
}

type PodLogs struct {
	Name       string
	Containers []*ContainerLogs
}

type ContainerLogs struct {
	Name    string
	Image   string
	Content string
}

type ApplicationLogsDocument struct {
	Cluster         string
	Kind            string
	Timestamp       int64
	ApplicationName string
	Namespace       string
	PodName         string
	ContainerName   string
	Image           string
	Content         string
}

func (l *ApplicationLogs) Flatten() []*ApplicationLogsDocument {
	var documents []*ApplicationLogsDocument
	for _, pod := range l.Pods {
		for _, container := range pod.Containers {
			documents = append(documents, &ApplicationLogsDocument{
				Cluster:         l.Cluster,
				Kind:            l.Kind,
				Timestamp:       l.Timestamp,
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
	InsertLogs(ctx context.Context, logs ApplicationLogs) error
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

func getApplicationLogsIndexName(applicationLogs *ApplicationLogs) string {
	return elasticsearch.GetIndexName(applicationLogs.Cluster, "applications", applicationLogs.Timestamp)
}

func (r *ElasticSearchApplicationLogsRepository) getIndiciesWithClusterAndDateRange(cluster string,
	fromDate time.Time,
	toDate time.Time) []string {

	filter := array.Filter(
		elasticsearch.FilterIndicesByClusterAndDateRange(cluster, fromDate, toDate))

	r.updateIndices()

	return filter(maps.Keys(r.indices))
}

func (r *ElasticSearchApplicationLogsRepository) GetLogs(ctx context.Context, cluster string, startDate time.Time, endDate time.Time) ([]*ApplicationLogsDocument, error) {

	indices := r.getIndiciesWithClusterAndDateRange(cluster, startDate, endDate)
	if len(indices) == 0 {
		return []*ApplicationLogsDocument{}, nil
	}
	res, err := elasticsearch.SearchIndices(ctx, r.esClient, indices)

	if err != nil {
		r.logger.Error("Failed to get node logs", zap.Error(err))
		return nil, err
	}

	var nodeLogs []*ApplicationLogsDocument
	for _, value := range res.Hits.Hits {
		var log ApplicationLogsDocument
		err := json.Unmarshal(value.Source_, &log)
		if err != nil {
			r.logger.Error("Failed to decode node logs", zap.Error(err))
			return nil, err
		}
		nodeLogs = append(nodeLogs, &log)

	}

	return nodeLogs, nil
}

func (r *ElasticSearchApplicationLogsRepository) CreateIndex(ctx context.Context, indexName string) error {

	_, err := r.esClient.Indices.Create(indexName).
		Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"cluster":         types.NewTextProperty(),
					"kind":            types.NewTextProperty(),
					"timestamp":       types.NewPointInTimeReference(),
					"applicationName": types.NewTextProperty(),
					"namespace":       types.NewTextProperty(),
					"podName":         types.NewTextProperty(),
					"containerName":   types.NewTextProperty(),
					"image":           types.NewTextProperty(),
					"content":         types.NewTextProperty(),
				},
			},
		}).Do(ctx)

	if err != nil {
		r.logger.Error("Failed to create an index", zap.Error(err))
		return err
	}

	err = r.updateIndices()
	if err != nil {
		r.logger.Error("Failed to update indices list", zap.Error(err))
	}

	return nil
}

func (r *ElasticSearchApplicationLogsRepository) InsertLogs(ctx context.Context, logs ApplicationLogs) error {

	index := getApplicationLogsIndexName(&logs)

	if !r.doesIndexExists(index) {
		r.CreateIndex(ctx, index)
	}

	flattenedLogs := make([]interface{}, 0, len(logs.Flatten()))
	for _, log := range logs.Flatten() {
		flattenedLogs = append(flattenedLogs, log)
	}

	//TODO: Replace with index based on query params
	_, err := r.esClient.Bulk().Index(index).Request(&flattenedLogs).Do(ctx)
	if err != nil {
		r.logger.Error("Failed to insert application logs", zap.Error(err))
	}

	return err
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
