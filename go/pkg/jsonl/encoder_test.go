package jsonl_test

import (
	// "encoding/json"
	// "fmt"
	"bytes"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/jsonl"
	// "reflect"
	// "strings"
	"testing"
)

type encodingTestCase struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestEncode(t *testing.T) {

	testsCases := []struct {
		description string
		encoded     string
		structSlice []encodingTestCase
	}{
		{
			description: "Given list of structs encode them into a string",
			encoded:     "{\"name\":\"John\",\"age\":23}\n{\"name\":\"Anna\",\"age\":32}\n",
			structSlice: []encodingTestCase{
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

			// testCaseReader := strings.NewReader(tc.rawString)
			results := bytes.NewBufferString("")
			sliceWithInterface := make([]interface{}, 0, len(tc.structSlice))
			for _, strct := range tc.structSlice {
				sliceWithInterface = append(sliceWithInterface, strct)
			}

			err := jsonl.NewJsonLinesEncoder(results).Encode(sliceWithInterface)
			if err != nil {
				t.Fatalf("Failed to decode an rawString into jsonl")
			}
			if tc.encoded != results.String() {
				t.Fatalf("Wanted %+v, got %+v", tc.encoded, results)
			}

		})
	}

}
