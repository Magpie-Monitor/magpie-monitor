package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func GetQueryByTimestamps(fromDate time.Time, toDate time.Time) *types.Query {

	fromDateFilter := types.Float64(fromDate.UnixNano())
	toDateFilter := types.Float64(toDate.UnixNano())

	return &types.Query{
		Range: map[string]types.RangeQuery{
			"timestamp": types.NumberRangeQuery{
				Gte: &fromDateFilter,
				Lte: &toDateFilter,
			},
		},
	}
}

func SearchIndices(ctx context.Context, esClient *elasticsearch.TypedClient, indices []string, query *types.Query) (*search.Response, error) {
	size := 1000
	return esClient.Search().
		Index(strings.Join(indices, ",")).
		Request(&search.Request{
			Query: query,

			Size: &size,
		}).Do(ctx)
}

func getYYYYMM(date time.Time) string {

	fmt.Printf("%d-%d", date.Year(), date.Month())

	return fmt.Sprintf("%d-%d", date.Year(), date.Month())
}

func GetIndexName(cluster string, sourceName string, timestamp int64) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		cluster,
		sourceName,
		getYYYYMM(time.Unix(0, timestamp)))
}

func GetIndexParams(index string) (cluster string, source string, year int, month int, err error) {
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

	cluster = strings.Join(elements[0:len(elements)-3], "-")

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

		existingCluster, existingKind, year, month, err := GetIndexParams(index)
		if err != nil {
			return false
		}

		existingDate := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, &time.Location{}).Unix()

		return cluster == existingCluster &&
			time.Date(fromDate.Year(), fromDate.Month(), 0, 0, 0, 0, 0, &time.Location{}).Unix() <= existingDate &&
			toDate.Unix() >= existingDate &&
			existingKind == kind
	}

}
