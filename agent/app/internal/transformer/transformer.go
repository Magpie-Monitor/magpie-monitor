package transformer

import "strings"

type Transformer interface {
	Transform(logs string) string
}

type DummyTransformer struct {
}

func (d DummyTransformer) Transform(logs string) string {
	return strings.ToUpper(logs)
}
