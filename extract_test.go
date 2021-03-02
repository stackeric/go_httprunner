package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	extractor := Extractor{
		"session_foo2": "name.first",
	}
	ctx := Context{}
	responseBody := []byte(`{"name":{"first":"Janet","last":"Prichard"},"age":47}`)

	err := extractor.Run(ctx, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "Janet", ctx["session_foo2"])
}
