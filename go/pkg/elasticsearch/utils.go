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
