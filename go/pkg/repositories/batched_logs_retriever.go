package repositories

import (
	"context"
	"encoding/json"
	elasticutils "github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type BatchedDocumentsRetriever[T any] interface {
	GetNextBatch() ([]T, error)
	HasNextBatch() bool
}

type ElasticBatchedDocumentsRetriever struct {
	scrollId string
	hits     []types.Hit
	client   *elasticsearch.TypedClient
}

type ElasticApplicationBatchedLogsRetriever struct {
	documentRetriever BatchedDocumentsRetriever[types.Hit]
}

type ElasticNodeBatchedLogsRetriever struct {
	documentRetriever BatchedDocumentsRetriever[types.Hit]
}

type ApplicationLogsBatchRetriever = BatchedDocumentsRetriever[*ApplicationLogsDocument]
type NodeLogsBatchRetriever = BatchedDocumentsRetriever[*NodeLogsDocument]

func NewElasticBatchedDocumentsRetriever(client *elasticsearch.TypedClient, scrollQuery *search.Search) (*ElasticBatchedDocumentsRetriever, error) {

	initialScroll, err := scrollQuery.
		Scroll(elasticutils.DEFAULT_ELASTIC_SCROLL_EXPIRY_PERIOD).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	return &ElasticBatchedDocumentsRetriever{
		scrollId: *initialScroll.ScrollId_,
		client:   client,
		hits:     initialScroll.Hits.Hits,
	}, nil
}

func (r *ElasticBatchedDocumentsRetriever) GetNextBatch() ([]types.Hit, error) {
	results, err := elasticutils.GetNextScrollPage(context.Background(), r.client, r.scrollId)
	if err != nil {
		return nil, err
	}

	previousHits := r.hits

	r.hits = results.Hits.Hits
	return previousHits, nil
}

func (r *ElasticBatchedDocumentsRetriever) HasNextBatch() bool {
	return len(r.hits) > 0
}

func NewElasticBatchedApplicationLogsRetriever(documentRetriever ElasticBatchedDocumentsRetriever) (*ElasticApplicationBatchedLogsRetriever, error) {

	return &ElasticApplicationBatchedLogsRetriever{
		documentRetriever: &documentRetriever,
	}, nil
}

func (r *ElasticApplicationBatchedLogsRetriever) GetNextBatch() ([]*ApplicationLogsDocument, error) {

	hits, err := r.documentRetriever.GetNextBatch()
	if err != nil {
		return nil, err
	}

	var applicationLogs []*ApplicationLogsDocument

	for _, value := range hits {
		var log ApplicationLogsDocument
		err := json.Unmarshal(value.Source_, &log)
		if err != nil {
			return nil, err
		}

		log.Id = *value.Id_

		if log.Content != "" {
			applicationLogs = append(applicationLogs, &log)
		}
	}

	return applicationLogs, nil
}

func (r *ElasticApplicationBatchedLogsRetriever) HasNextBatch() bool {
	return r.documentRetriever.HasNextBatch()
}

func NewElasticBatchedNodeLogsRetriever(documentRetriever ElasticBatchedDocumentsRetriever) (*ElasticNodeBatchedLogsRetriever, error) {

	return &ElasticNodeBatchedLogsRetriever{
		documentRetriever: &documentRetriever,
	}, nil
}

func (r *ElasticNodeBatchedLogsRetriever) GetNextBatch() ([]*NodeLogsDocument, error) {

	hits, err := r.documentRetriever.GetNextBatch()
	if err != nil {
		return nil, err
	}

	var nodeLogs []*NodeLogsDocument

	for _, value := range hits {
		var log NodeLogsDocument
		err := json.Unmarshal(value.Source_, &log)
		if err != nil {
			return nil, err
		}

		log.Id = *value.Id_

		if log.Content != "" {
			nodeLogs = append(nodeLogs, &log)
		}
	}

	return nodeLogs, nil
}

func (r *ElasticNodeBatchedLogsRetriever) HasNextBatch() bool {
	return r.documentRetriever.HasNextBatch()
}

var _ ApplicationLogsBatchRetriever = &ElasticApplicationBatchedLogsRetriever{}
var _ NodeLogsBatchRetriever = &ElasticNodeBatchedLogsRetriever{}
