package repositories

import (
	"context"
	"encoding/json"
	"github.com/IBM/fp-go/array"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	es "github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"time"
)

// TODO: To be clarified once the contract from agent is agreed upon
type NodeLogs struct {
	Id            string `json:"_id,omitempty"`
	ClusterId     string `json:"clusterId"`
	Kind          string `json:"kind"`
	CollectedAtMs int64  `json:"collectedAtMs"`
	Name          string `json:"name"`
	Filename      string `json:"filename"`
	Content       string `json:"content"`
}

type NodeLogsDocument = NodeLogs

type NodeLogsRepository interface {
	CreateIndex(ctx context.Context, indexName string) error
	GetLogs(ctx context.Context, cluster string, startDate time.Time, endDate time.Time) ([]*NodeLogsDocument, error)
	InsertLogs(ctx context.Context, logs *NodeLogs) error
}

func ProvideAsNodeLogsRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(NodeLogsRepository)),
	)
}

func (l *NodeLogsDocument) GetContent() string {
	return l.Content
}

type ElasticSearchNodeLogsRepository struct {
	esClient *es.TypedClient
	logger   *zap.Logger
	indices  map[string]bool
}

func (r *ElasticSearchNodeLogsRepository) doesIndexExists(index string) bool {
	return r.indices[index]
}

func getNodeLogsIndexName(nodeLogs *NodeLogs) string {
	return elasticsearch.GetIndexName(nodeLogs.ClusterId, "nodes", nodeLogs.CollectedAtMs)
}

func (r *ElasticSearchNodeLogsRepository) getIndiciesWithClusterAndDateRange(cluster string,
	fromDate time.Time,
	toDate time.Time) []string {

	filter := array.Filter(
		elasticsearch.FilterIndicesByClusterAndDateRange(cluster, "nodes", fromDate, toDate))

	r.updateIndices()

	return filter(maps.Keys(r.indices))
}

func (r *ElasticSearchNodeLogsRepository) GetLogs(ctx context.Context, cluster string, startDate time.Time, endDate time.Time) ([]*NodeLogsDocument, error) {

	indices := r.getIndiciesWithClusterAndDateRange(cluster, startDate, endDate)
	if len(indices) == 0 {
		return []*NodeLogsDocument{}, nil
	}

	query := elasticsearch.GetQueryByTimestamps(startDate, endDate)
	res, err := elasticsearch.SearchIndices(ctx, r.esClient, indices, query)

	if err != nil {
		r.logger.Error("Failed to get node logs", zap.Error(err))
		return nil, err
	}

	var nodeLogs []*NodeLogsDocument
	for _, value := range res.Hits.Hits {
		var log NodeLogsDocument
		err := json.Unmarshal(value.Source_, &log)
		if err != nil {
			r.logger.Error("Failed to decode node logs", zap.Error(err))
			return nil, err
		}

		log.Id = *value.Id_

		if log.Content != "" {
			nodeLogs = append(nodeLogs, &log)
		}

	}

	return nodeLogs, nil
}

func (r *ElasticSearchNodeLogsRepository) CreateIndex(ctx context.Context, indexName string) error {

	_, err := r.esClient.Indices.Create(indexName).Do(ctx)

	if err != nil {
		r.logger.Error("Failed to create an index for logsdb", zap.Error(err))
		return err
	}

	err = r.updateIndices()
	if err != nil {
		r.logger.Error("Failed to update logsdb indices list", zap.Error(err))
	}

	return nil
}

func (r *ElasticSearchNodeLogsRepository) InsertLogs(ctx context.Context, logs *NodeLogs) error {

	index := getNodeLogsIndexName(logs)

	if !r.doesIndexExists(index) {
		r.CreateIndex(ctx, index)
	}

	_, err := r.esClient.Index(index).Request(logs).Do(ctx)
	if err != nil {
		r.logger.Error("Failed to insert node logs", zap.Error(err))
	}

	return err
}

type NodeLogsParams struct {
	fx.In
	ElasticSearchClient *es.TypedClient
	Logger              *zap.Logger
}

func NewElasticSearchNodeLogsRepository(p NodeLogsParams) *ElasticSearchNodeLogsRepository {

	repository := ElasticSearchNodeLogsRepository{
		esClient: p.ElasticSearchClient,
		logger:   p.Logger,
	}

	err := repository.updateIndices()
	if err != nil {
		repository.logger.Error("Failed to create fetch initial logsdb indices", zap.Error(err))
	}

	return &repository
}

func (r *ElasticSearchNodeLogsRepository) updateIndices() error {

	indices, err := elasticsearch.GetAllIndicesSet(r.esClient)
	if err != nil {
		r.logger.Error("Failed to fetch logsdb indices", zap.Error(err))
		return err
	}

	r.indices = indices
	r.logger.Info("Fetched logsdb indices: ", zap.Any("indices", indices))
	return nil
}

// Compile-time check if ElasticSearchNodeLogsRepository implements
// the NodeLogsRepository interface
var _ NodeLogsRepository = &ElasticSearchNodeLogsRepository{}
