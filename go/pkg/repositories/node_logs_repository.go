package repositories

import (
	"context"
	"encoding/json"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NodeLogs struct {
	Timestamp uint
	Host      string
	Message   string
}

type NodeLogsRepository interface {
	CreateIndex(ctx context.Context, indexName string) error
	GetAllLogs(ctx context.Context) ([]NodeLogs, error)
	InsertLogs(ctx context.Context, logs NodeLogs) error
}

func ProvideAsNodeLogsRepository(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(NodeLogsRepository)),
	)
}

type ElasticSearchNodeLogsRepository struct {
	esClient *es.TypedClient
	logger   *zap.Logger
}

func (r *ElasticSearchNodeLogsRepository) GetAllLogs(ctx context.Context) ([]NodeLogs, error) {
	r.logger.Info("Called a Get all NodeLogs!")
	res, err := r.esClient.Search().
		Index("test_index").
		Request(&search.Request{
			Query: &types.Query{
				MatchAll: &types.MatchAllQuery{},
			},
		}).Do(ctx)

	if err != nil {
		r.logger.Error("Failed to get all node logs", zap.Error(err))
		return nil, err
	}

	var nodeLogs []NodeLogs
	for _, value := range res.Hits.Hits {
		var log NodeLogs
		err := json.Unmarshal(value.Source_, &log)
		if err != nil {
			r.logger.Error("Failed to decode node logs", zap.Error(err))
			return nil, err
		}
		nodeLogs = append(nodeLogs, log)

	}

	return nodeLogs, nil
}

func (r *ElasticSearchNodeLogsRepository) CreateIndex(ctx context.Context, indexName string) error {

	_, err := r.esClient.Indices.Create(indexName).
		Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"timestamp": types.NewIntegerNumberProperty(),
					"host":      types.NewTextProperty(),
					"message":   types.NewTextProperty(),
				},
			},
		}).Do(ctx)

	return err
}

func (r *ElasticSearchNodeLogsRepository) InsertLogs(ctx context.Context, logs NodeLogs) error {

	_, err := r.esClient.Index("test_index").Request(logs).Do(ctx)
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

	return &ElasticSearchNodeLogsRepository{
		esClient: p.ElasticSearchClient,
		logger:   p.Logger,
	}
}

var _ NodeLogsRepository = &ElasticSearchNodeLogsRepository{}
