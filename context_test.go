package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	t.Run("Parse string Variable", func(t *testing.T) {
		ctx := Context{}
		v := Variable{
			"foo1": "foo1",
		}
		err := v.Parse(ctx)
		assert.Nil(t, err)
		assert.Equal(t, "foo1", v["foo1"])
		assert.Equal(t, "foo1", ctx["foo1"])
	})

	t.Run("Parse int Variable", func(t *testing.T) {
		ctx := Context{}
		v := Variable{
			"foo2": 2,
		}
		err := v.Parse(ctx)
		assert.Nil(t, err)

		assert.Equal(t, 2, v["foo2"])
	})

	t.Run("Parse Bool Variable", func(t *testing.T) {
		ctx := Context{}
		v := Variable{
			"foo3": false,
		}
		err := v.Parse(ctx)
		assert.Nil(t, err)

		assert.Equal(t, false, v["foo3"])
	})

	t.Run("Parse dynamic Variable", func(t *testing.T) {
		ctx := Context{
			"foo1": "foo1",
		}
		v := Variable{
			"foo1": "$foo1",
		}
		err := v.Parse(ctx)
		assert.Nil(t, err)

		assert.Equal(t, "foo1", v["foo1"])
	})

	t.Run("Parse Func Variable", func(t *testing.T) {
		ctx := Context{}
		v := Variable{
			"sum_v": "${sum_two(1, 2)}",
		}
		err := v.Parse(ctx)
		assert.Nil(t, err)
		assert.Equal(t, 3, v["sum_v"])
	})
}

func TestIsFunc(t *testing.T) {
	t.Run("With Args", func(t *testing.T) {
		v := "${sum_two(1, 2)}"
		got, ok := IsFunc(v)
		assert.Equal(t, true, ok)
		assert.Equal(t, 3, got)
	})
}
