package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	YamlFile      = "examples/example.yaml"
	ErrorYamlFile = "examples/error.yaml"
)

func TestTestCase(t *testing.T) {
	t.Run("Load Test Case From Yaml File", func(t *testing.T) {
		yamlFilePath := YamlFile
		var c TestCase
		err := c.LoadCaseFromYaml(yamlFilePath)
		assert.Nil(t, err, fmt.Sprintf("got %v, want nil", err))
	})
	t.Run("Read Value From TestCase ", func(t *testing.T) {
		testCase, err := NewTestCase(YamlFile)
		assert.Equal(t, err, nil)
		assert.True(t, testCase.IsValid())

		wantConfig := "request methods testcase: validate with functions"

		wantStepsLen1 := 1
		assert.Equal(t, wantConfig, testCase.Config.Name)
		assert.Equal(t, wantStepsLen1, len(testCase.TestSteps))

		wantStepName := "get with params"
		assert.Equal(t, wantStepName, testCase.TestSteps[0].Name)
	})

	t.Run("Check Required Field ", func(t *testing.T) {
		yamlFilePath := ErrorYamlFile
		testCase, err := NewTestCase(yamlFilePath)
		assert.Equal(t, err, ErrKeyNotFound)
		assert.False(t, testCase.IsValid())
	})

	t.Run("Run Test Case", func(t *testing.T) {

		testCase, err := NewTestCase(YamlFile)
		assert.Equal(t, err, nil)
		assert.True(t, testCase.IsValid())

		_, err = testCase.Run()
		assert.Nil(t, err)
	})
	t.Run("TestCase Report", func(t *testing.T) {
		testCase, err := NewTestCase(YamlFile)
		assert.Equal(t, err, nil)
		assert.True(t, testCase.IsValid())
		report, err := testCase.Run()
		assert.Nil(t, err)
		assert.True(t, report.Result())
		assert.Len(t, report.StepReports, len(testCase.TestSteps))
		for index, rep := range report.StepReports {
			assert.Equal(t, rep.index, index)
			assert.True(t, rep.result)
		}
	})

}
