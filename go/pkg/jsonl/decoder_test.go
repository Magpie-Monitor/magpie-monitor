package jsonl_test

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/jsonl"
	"reflect"
	"strings"
	"testing"
)

type decodingTestCase struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecode(t *testing.T) {

	testsCases := []struct {
		description string
		rawString   string
		expect      []decodingTestCase
	}{
		{
			description: "Given list of jsons decode them into a slice of structs",
			rawString:   "{\"name\": \"John\", \"age\":23}\n{\"name\": \"Anna\", \"age\":32}\n",
			expect: []decodingTestCase{
				{
					Name: "John",
					Age:  23,
				},
				{
					Name: "Anna",
					Age:  32,
				},
			},
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.description, func(t *testing.T) {

			testCaseReader := strings.NewReader(tc.rawString)
			results := make([]decodingTestCase, 0)

			err := jsonl.NewJsonLinesDecoder(testCaseReader).Decode(&results)
			if err != nil {
				t.Fatalf("Failed to decode an rawString into jsonl")
			}
			if !reflect.DeepEqual(tc.expect, results) {
				t.Fatalf("Wanted %+v, got %+v", tc.expect, results)
			}

		})
	}

}
