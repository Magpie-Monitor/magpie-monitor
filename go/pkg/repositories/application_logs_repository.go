package repositories

import (
	"context"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ApplicationLogs struct {
	Timestamp        uint
	Application_name string
	Message          string
	Image_name       string
}

type ApplicationLogsRepository interface {
	CreateIndex(ctx context.Context, indexName string) error
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
	r.logger.Info("Called a Get all ApplicationLogs!")
	return nil, nil
}

func (r *ElasticSearchApplicationLogsRepository) CreateIndex(ctx context.Context, indexName string) error {

	_, err := r.esClient.Indices.Create(indexName).
		Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"source": types.NewTextProperty(),
					"data":   types.NewTextProperty(),
				},
			},
		}).Do(ctx)

	return err
}

func (r *ElasticSearchApplicationLogsRepository) InsertLogs(ctx context.Context, logs ApplicationLogs) error {

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

var _ ApplicationLogsRepository = &ElasticSearchApplicationLogsRepository{}
