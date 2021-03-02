package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"name":{"first":"Janet","last":"Prichard"},"age":47}`)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	response, err := NewResponse(resp)
	assert.Nil(t, err)
	value := gjson.GetBytes(response.StrBody, "status_code")
	assert.True(t, value.Exists())
	assert.Equal(t, value.Int(), int64(200))
}
