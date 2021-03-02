package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEachStep(t *testing.T) {

	t.Run("Check Step Variable Parse", func(t *testing.T) {
		yamlFilePath := YamlFile
		testCase, err := NewTestCase(yamlFilePath)
		assert.Equal(t, err, nil)
		for _, step := range testCase.TestSteps {
			variables := step.Variables
			assert.Equal(t, "session_bar2", variables["foo2"])
			assert.Equal(t, "${sum_two(1, 2)}", variables["sum_v"])

			ctx := Context{}
			err = variables.Parse(ctx)

			assert.Nil(t, err)

			assert.Equal(t, "session_bar2", variables["foo2"])
			assert.Equal(t, 3, variables["sum_v"])
		}
	})

	t.Run("Check Request Params Parse", func(t *testing.T) {
		testCase, err := NewTestCase(YamlFile)
		assert.Equal(t, err, nil)
		ctx := testCase.CaseCtx
		for _, step := range testCase.TestSteps {
			err = step.Variables.Parse(ctx)
			assert.Nil(t, err)

			params := step.Request.Params
			err := params.Parse(ctx)
			assert.Nil(t, err)
			assert.Equal(t, "bar1", params["foo1"])
			assert.Equal(t, "session_bar2", params["foo2"])
			assert.Equal(t, 3, params["sum_v"])
		}
	})
	t.Run("Step Req And Validator", func(t *testing.T) {
		yamlFilePath := YamlFile
		testCase, err := NewTestCase(yamlFilePath)
		assert.Equal(t, err, nil)
		assert.True(t, testCase.IsValid())

		for index, step := range testCase.TestSteps {

			rp, err := step.Run(index, testCase.Config.BaseURL, testCase.CaseCtx)
			assert.Nil(t, err)
			response := rp.response
			assert.Nil(t, err)
			got := step.Validators
			wantLen := 2
			assert.Len(t, got, wantLen)
			for _, v := range got {
				_, err := v.Run(response.StrBody)
				assert.Nil(t, err)
			}
		}

	})
}
