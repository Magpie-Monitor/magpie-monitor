package jsonl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type JsonLinesDecoder struct {
	reader io.Reader
}

func NewJsonLinesDecoder(reader io.Reader) *JsonLinesDecoder {
	return &JsonLinesDecoder{
		reader: reader,
	}
}

func (e *JsonLinesDecoder) Decode(v any) error {

	if reflect.TypeOf(v).Kind() != reflect.Pointer {
		panic("JsonLines decoder only accepts pointer to slice")
	}

	if reflect.TypeOf(v).Elem().Kind() != reflect.Slice {
		panic("JsonLines decoder only accepts pointer to slice!")
	}

	elementType := reflect.TypeOf(v).Elem().Elem()
	content := bytes.NewBufferString("")
	_, err := io.Copy(content, e.reader)
	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't read raw jsonl content %s", err.Error()))
	}

	rv := reflect.ValueOf(v).Elem()

	trimmedContent := strings.Trim(content.String(), "\n")

	lines := strings.Split(trimmedContent, "\n")
	for i, line := range lines {
		elementInstance := reflect.New(elementType).Interface()
		err := json.Unmarshal([]byte(line), &elementInstance)

		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't decode jsonl line %s", err.Error()))
		}

		if i >= rv.Cap() {
			rv.Grow(1)
		}
		if i >= rv.Len() {
			rv.SetLen(i + 1)
		}

		rv.Index(i).Set(reflect.ValueOf(elementInstance).Elem())
	}

	return nil
}
