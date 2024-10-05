package jsonl

import (
	"encoding/json"
	"io"
	"reflect"
)

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
		panic("JsonLinesEncoder accepts only slices of structs")
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
