package jsonl

import (
	"encoding/json"
	"io"
)

type JsonLinesEncoder struct {
	writer io.Writer
}

func NewJsonLinesEncoder(writer io.Writer) *JsonLinesEncoder {
	return &JsonLinesEncoder{
		writer: writer,
	}
}

func (e *JsonLinesEncoder) Encode(v []interface{}) error {

	for _, elem := range v {
		jsonEncoding, err := json.Marshal(elem)
		if err != nil {
			return err
		}

		e.writer.Write(jsonEncoding)
		e.writer.Write([]byte("\n"))
	}

	return nil
}
