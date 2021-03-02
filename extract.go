package main

import (
	"github.com/tidwall/gjson"
)

// Extractor contain value will extract from response
type Extractor map[string]string

// Run extract value from body ,set env in context
func (e Extractor) Run(ctx Context, body []byte) error {
	for k, v := range e {
		value := gjson.GetBytes(body, v)
		if !value.Exists() {
			return ErrKeyNotFound
		}
		ctx[k] = value.Value()
	}
	return nil
}
