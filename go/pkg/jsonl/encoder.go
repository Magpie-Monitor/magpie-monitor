package jsonl

import (
	"encoding/json"
	"io"
	"reflect"
)

// JSONL is a encoding format which is an equivalent of list of JSONs
// where each object is separated by a newline. This format is used in the OpenAI Batch API.
type JsonLinesEncoder struct {
	writer io.Writer
}

func NewJsonLinesEncoder(writer io.Writer) *JsonLinesEncoder {
	return &JsonLinesEncoder{
		writer: writer,
	}
}

func (e *JsonLinesEncoder) Encode(v any) error {

	if reflect.TypeOf(v).Kind() != reflect.Slice {
		panic("JsonLinesEncoder accepts only slices")
	}
	rv := reflect.ValueOf(v)

	for i := 0; i < rv.Len(); i++ {

		jsonEncoding, err := json.Marshal(rv.Index(i).Interface())
		if err != nil {
			return err
		}

		e.writer.Write(jsonEncoding)
		e.writer.Write([]byte("\n"))
	}

	return nil
}
