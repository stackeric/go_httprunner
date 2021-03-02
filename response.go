package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//Response is http response and more
type Response struct {
	Resp    *http.Response
	StrBody []byte
}

// NewResponse reconstruct http response
func NewResponse(r *http.Response) (Response, error) {
	rr := Response{
		Resp: r,
	}
	body, _ := ioutil.ReadAll(rr.Resp.Body)
	defer r.Body.Close()
	m := make(map[string]interface{})
	err := json.Unmarshal(body, &m)
	if err != nil {
		return rr, err
	}
	m["status_code"] = r.StatusCode

	b, err := json.Marshal(m)
	if err != nil {
		return rr, err
	}

	rr.StrBody = b
	return rr, nil
}
