package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/mget"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/scroll"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_ELASTIC_SCROLL_EXPIRY_PERIOD = "1d"
const DEFAULT_ELASTIC_PAGE_SIZE = 10000

func GetQueryByTimestamps(fromDate time.Time, toDate time.Time) *types.Query {

	fromDateFilter := types.Float64(fromDate.UnixMilli())
	toDateFilter := types.Float64(toDate.UnixMilli())

	return &types.Query{
		Range: map[string]types.RangeQuery{
			"collectedAtMs": types.NumberRangeQuery{
				Gte: &fromDateFilter,
				Lte: &toDateFilter,
			},
		},
	}
}

func GetQueryByTerms(terms map[string]types.TermsQueryField) *types.Query {

	var boost float32 = 1
	return &types.Query{
		Terms: &types.TermsQuery{
			Boost:      &boost,
			TermsQuery: terms,
		},
	}
}

func SearchIndices(ctx context.Context, esClient *elasticsearch.TypedClient, indices []string, query *types.Query) (*search.Response, error) {
	size := DEFAULT_ELASTIC_PAGE_SIZE
	return esClient.Search().
		Index(strings.Join(indices, ",")).
		Request(&search.Request{
			Query: query,

			Size: &size,
		}).
		Do(ctx)
}

func GetDocumentsByIds(ctx context.Context, esClient *elasticsearch.TypedClient, indicies []string, ids []string) (*mget.Response, error) {
	return esClient.
		Mget().
		Index(strings.Join(indicies, ",")).
		Request(
			&mget.Request{
				Ids: ids,
			},
		).Do(ctx)
}

func RequestSearchScroll(
	ctx context.Context,
	esClient *elasticsearch.TypedClient,
	indices []string,
	query *types.Query,
	size int) (*search.Response, error) {
	return esClient.Search().
		Index(strings.Join(indices, ",")).
		Request(&search.Request{
			Query: query,
			Size:  &size,
		}).
		Scroll(DEFAULT_ELASTIC_SCROLL_EXPIRY_PERIOD).
		Do(ctx)
}

func GetSearchQuery(ctx context.Context, esClient *elasticsearch.TypedClient, indices []string, query *types.Query, pageSize int) *search.Search {
	return esClient.Search().
		Index(strings.Join(indices, ",")).
		Request(&search.Request{
			Query: query,
			Size:  &pageSize,
		})
}

func GetNextScrollPage(ctx context.Context, esClient *elasticsearch.TypedClient, scrollId string) (*scroll.Response, error) {
	return esClient.Scroll().
		Scroll(DEFAULT_ELASTIC_SCROLL_EXPIRY_PERIOD).
		ScrollId(scrollId).
		Do(ctx)
}

func getYYYYMM(date time.Time) string {

	fmt.Printf("%d-%d", date.Year(), date.Month())

	return fmt.Sprintf("%d-%d", date.Year(), date.Month())
}

func GetIndexName(cluster string, sourceName string, collectedAtMs int64) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		cluster,
		sourceName,
		getYYYYMM(time.UnixMilli(collectedAtMs)))
}

func GetIndexParams(index string) (clusterId string, source string, year int, month int, err error) {
	elements := strings.Split(index, "-")

	if len(elements) < 3 {
		err = errors.New("Invalid cluster name")
		return
	}

	year, err = strconv.Atoi(elements[len(elements)-2])
	if err != nil {
		err = errors.New("Invalid cluster name")
		return
	}

	month, err = strconv.Atoi(elements[len(elements)-1])
	if err != nil {
		err = errors.New("Invalid cluster name")
		return
	}

	source = elements[len(elements)-3]

	clusterId = strings.Join(elements[0:len(elements)-3], "-")

	return

}

func GetAllIndicesSet(esClient *elasticsearch.TypedClient) (map[string]bool, error) {

	indicesResponse, err := esClient.Cat.Indices().Do(context.TODO())
	if err != nil {
		return nil, err
	}

	indices := make(map[string]bool, len(indicesResponse))

	for _, index := range indicesResponse {
		indices[*index.Index] = true
	}

	return indices, nil
}
func FilterIndicesByClusterAndDateRange(cluster string, kind string, fromDate time.Time, toDate time.Time) func(index string) bool {

	return func(index string) bool {

		// Don't match indices by incorrect dates.
		if fromDate.Unix() > toDate.Unix() {
			return false
		}

		clusterFromIndex, kindFromIndex, yearFromIndex, monthFromIndex, err := GetIndexParams(index)
		if err != nil {
			return false
		}

		dateFromIndex := time.Date(yearFromIndex, time.Month(monthFromIndex), 0, 0, 0, 0, 0, &time.Location{}).Unix()

		return cluster == clusterFromIndex &&
			// Match indecies from the month of the fromDate.
			time.Date(fromDate.Year(), fromDate.Month(), 0, 0, 0, 0, 0, &time.Location{}).Unix() <= dateFromIndex &&
			toDate.Unix() >= dateFromIndex &&
			kindFromIndex == kind
	}
}

func GetAndMapDocumentsByIds[T any](ctx context.Context,
	client *elasticsearch.TypedClient,
	indicies []string,
	ids []string,
	logger *zap.Logger,
) (map[string]T, error) {
	logsByIdsQuery, err := GetDocumentsByIds(ctx,
		client,
		indicies,
		ids)
	if err != nil {
		logger.Error("Failed to get documents by ids query",
			zap.Error(err),
			zap.Any("ids", ids),
		)
		return nil, err
	}

	mappedDocuments := make(map[string]T, 0)

	for _, value := range logsByIdsQuery.Docs {
		var log T

		result, ok := value.(*types.GetResult)
		if !ok || result.Source_ == nil {
			logger.Error("Failed to get document from id", zap.Any("document", value))

			// Skipping in case of made up id
			continue
		}

		err := json.Unmarshal(result.Source_, &log)
		logger.Info("Debug fetched source", zap.Any("source", result.Source_))
		if err != nil {
			logger.Error("Failed to decode documents", zap.Error(err), zap.Any("logs", result.Source_))
			return nil, err
		}

		mappedDocuments[result.Id_] = log
	}

	return mappedDocuments, nil
}
