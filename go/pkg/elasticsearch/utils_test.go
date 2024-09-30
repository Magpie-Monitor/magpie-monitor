package elasticsearch_test

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/elasticsearch"
	"testing"
	"time"
)

func TestFilterIndicesByClusterAndDateRange(t *testing.T) {
	testCases := []struct {
		cluster        string
		kind           string
		fromDate       time.Time
		toDate         time.Time
		index          string
		wantUnfiltered bool
		description    string
	}{
		{
			cluster:        "test_cluster",
			kind:           "applications",
			fromDate:       time.Date(2024, 9, 14, 30, 20, 10, 100, time.Local),
			toDate:         time.Date(2024, 10, 14, 30, 20, 10, 100, time.Local),
			index:          "test_cluster-applications-2024-09",
			wantUnfiltered: true,
			description:    "Match index for a two month difference between fromDate and toDate for a first month",
		},
		{
			cluster:        "test_cluster",
			kind:           "applications",
			fromDate:       time.Date(2024, 9, 14, 30, 20, 10, 100, time.Local),
			toDate:         time.Date(2024, 10, 14, 30, 20, 10, 100, time.Local),
			index:          "test_cluster-applications-2024-10",
			wantUnfiltered: true,
			description:    "Match index for a two month difference between fromDate and toDate for a second month",
		},
		{
			cluster:        "test_cluster",
			kind:           "applications",
			fromDate:       time.Date(2024, 9, 14, 30, 20, 10, 100, time.Local),
			toDate:         time.Date(2024, 9, 15, 30, 20, 10, 100, time.Local),
			index:          "test_cluster-applications-2024-9",
			wantUnfiltered: true,
			description:    "Match index for an one month difference between fromDate and toDate",
		},
		{
			cluster:        "test_cluster",
			kind:           "applications",
			fromDate:       time.Date(2024, 9, 14, 30, 20, 10, 100, time.Local),
			toDate:         time.Date(2024, 9, 2, 30, 20, 10, 100, time.Local),
			index:          "test_cluster-applications-2024-9",
			wantUnfiltered: false,
			description:    "Don't match a date with fromDate bigger than toDate",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			isFiltered := elasticsearch.FilterIndicesByClusterAndDateRange(
				tc.cluster,
				tc.kind,
				tc.fromDate,
				tc.toDate,
			)(tc.index)

			if isFiltered != tc.wantUnfiltered {
				t.Fatalf("Want %t, got %t for %+v", tc.wantUnfiltered, isFiltered, tc)
			}
		})
	}

}
