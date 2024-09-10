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

// TODO: To be clarified once the contract from agent is agreed upon
type ApplicationLogs struct {
	Timestamp       uint
	ApplicationName string
	Message         string
	ImageName       string
}

type ApplicationLogsRepository interface {
	CreateIndex(ctx context.Context, indexName string) error
	//TODO: To be removed after proper implementation of log-ingestion
	GetAllLogs(ctx context.Context) ([]*ApplicationLogs, error)
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
}

func (r *ElasticSearchApplicationLogsRepository) GetAllLogs(ctx context.Context) ([]*ApplicationLogs, error) {

	//TODO: Replace with index based on query params
	res, err := r.esClient.Search().
		Index("test_index").
		Request(&search.Request{
			Query: &types.Query{
				MatchAll: &types.MatchAllQuery{},
			},
		}).Do(ctx)

	if err != nil {
		r.logger.Error("Failed to get all application logs", zap.Error(err))
		return nil, err
	}

	var applicationLogs []*ApplicationLogs
	for _, value := range res.Hits.Hits {
		var log ApplicationLogs
		err := json.Unmarshal(value.Source_, &log)
		if err != nil {
			r.logger.Error("Failed to decode application logs", zap.Error(err))
			return nil, err
		}
		applicationLogs = append(applicationLogs, &log)

	}

	return applicationLogs, nil
}

func (r *ElasticSearchApplicationLogsRepository) CreateIndex(ctx context.Context, indexName string) error {

	_, err := r.esClient.Indices.Create(indexName).
		Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"timestamp":       types.NewTextProperty(),
					"applicationName": types.NewTextProperty(),
					"message":         types.NewTextProperty(),
					"imageName":       types.NewTextProperty(),
				},
			},
		}).Do(ctx)

	return err
}

func (r *ElasticSearchApplicationLogsRepository) InsertLogs(ctx context.Context, logs ApplicationLogs) error {

	//TODO: Replace with index based on query params
	_, err := r.esClient.Index("test_index").Request(logs).Do(ctx)
	if err != nil {
		r.logger.Error("Failed to insert application logs", zap.Error(err))
	}

	return err
}

type ApplicationLogsParams struct {
	fx.In
	ElasticSearchClient *es.TypedClient
}

func NewElasticSearchApplicationLogsRepository(p ApplicationLogsParams) *ElasticSearchApplicationLogsRepository {

	return &ElasticSearchApplicationLogsRepository{
		esClient: p.ElasticSearchClient,
	}
}

// Compile-time check if ElasticSearchApplicationLogsRepository implements
// the ApplicationLogsRepository interface
var _ ApplicationLogsRepository = &ElasticSearchApplicationLogsRepository{}
