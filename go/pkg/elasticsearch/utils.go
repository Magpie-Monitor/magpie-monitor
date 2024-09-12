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

func SearchIndices(ctx context.Context, esClient *elasticsearch.TypedClient, indices []string) (*search.Response, error) {
	return esClient.Search().
		Index(strings.Join(indices, ",")).
		Request(&search.Request{
			Query: &types.Query{
				MatchAll: &types.MatchAllQuery{},
			},
		}).Do(ctx)
}

func getYYYYMM(date time.Time) string {
	return fmt.Sprintf("%d-%d", date.Year(), date.Month())
}

func GetIndexName(cluster string, sourceName string, timestamp int64) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		cluster,
		sourceName,
		getYYYYMM(time.Unix(timestamp, 0)))
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
func FilterIndicesByClusterAndDateRange(cluster string, fromDate time.Time, toDate time.Time) func(index string) bool {

	return func(index string) bool {

		existingCluster, _, year, month, err := GetIndexParams(index)
		if err != nil {
			return false
		}

		return cluster == existingCluster &&
			fromDate.Unix() <= time.Date(year, time.Month(month), 0, 0, 0, 0, 0, &time.Location{}).Unix() &&
			toDate.Unix() >= time.Date(year, time.Month(month), 0, 0, 0, 0, 0, &time.Location{}).Unix()
	}

}
