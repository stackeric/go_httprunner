package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestSuite(t *testing.T) {

	suiteYamlFile := "examples/suite.yaml"
	t.Run("test suite parser", func(t *testing.T) {
		suite, err := NewTestSuite(suiteYamlFile)
		assert.Nil(t, err)
		assert.True(t, suite.IsValid())
	})
}
