package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	strbody := `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, strbody)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	response, err := NewResponse(resp)
	assert.Nil(t, err)
	body := response.StrBody
	t.Run("validator ok", func(t *testing.T) {
		vs := []Validator{
			{
				Compare:  "eq",
				Key:      "name.first",
				Expected: "Janet",
			},
		}
		for _, v := range vs {
			_, err := v.Run(body)
			assert.Nil(t, err)
		}
	})
	t.Run("validator key not exist", func(t *testing.T) {

		vs := []Validator{
			{
				Compare:  "eq",
				Key:      "body.args.sum_v",
				Expected: "3",
			},
		}
		for _, v := range vs {
			_, err := v.Run(body)
			assert.Equal(t, err, ErrKeyNotFound)
		}
	})

	t.Run("validator value error", func(t *testing.T) {
		vs := []Validator{
			{
				Compare:  "eq",
				Key:      "name.first",
				Expected: "3",
			},
		}
		for _, v := range vs {
			_, err := v.Run(body)
			assert.Equal(t, err, ErrValidate)
		}
	})

	t.Run("validator http status code", func(t *testing.T) {
		vs := []Validator{
			{
				Compare:  "eq",
				Key:      "status_code",
				Expected: 200,
			},
		}
		for _, v := range vs {
			_, err := v.Run(body)
			assert.Nil(t, err)
		}
	})
}
