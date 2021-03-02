package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert.Equal(t, ErrKeyNotFound.Error(), "you miss some required key")
}
