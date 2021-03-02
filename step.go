package main

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

var (
	// TimeOut for http request
	TimeOut = 10 * time.Second
)

// Request contain an api info
type Request struct {
	Method string     `yaml:"method"`
	URL    string     `yaml:"url"`
	Params Parameters `yaml:"params"`
}

// TestStep is TestCase steps
type TestStep struct {
	StepCtx    Context
	Name       string      `yaml:"name"`
	Variables  Variable    `yaml:"variables"`
	Request    Request     `yaml:"request"`
	Extractors Extractor   `yaml:"extract"`
	Validators []Validator `yaml:"validate"`
}

//IsValid check
func (c *TestStep) IsValid() bool {
	if c.Name == "" || c.Request.Method == "" {
		return false
	}
	return true
}

//Run every single step
func (c *TestStep) Run(index int, base string, CaseCtx Context) (r StepReport, err error) {
	// Init Variable and Context
	c.StepCtx = CaseCtx
	err = c.Variables.Parse(c.StepCtx)
	if err != nil {
		return
	}
	err = c.Request.Params.Parse(c.StepCtx)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()

	originURL := base + c.Request.URL

	r = StepReport{
		url:    originURL,
		method: c.Request.Method,
	}

	req, err := http.NewRequestWithContext(ctx, c.Request.Method, originURL, nil)
	q := req.URL.Query()
	for k, v := range c.Request.Params {
		switch v.(type) {
		case int, float64:
			q.Add(k, strconv.Itoa(v.(int)))
		case bool:
			q.Add(k, strconv.FormatBool(v.(bool)))
		default:
			q.Add(k, v.(string))
		}
	}
	req.URL.RawQuery = q.Encode()
	r.req = req
	if err != nil {
		r.result = false
		r.err = err
		return
	}
	client := &http.Client{
		Timeout: TimeOut,
	}

	resp, err := client.Do(req)
	if err != nil {
		r.result = false
		r.err = err
		return
	}
	response, err := NewResponse(resp)
	if err != nil {
		r.result = false
		r.err = err
		return
	}
	err = c.Extractors.Run(c.StepCtx, response.StrBody)
	if err != nil {
		r.result = false
		r.err = err
		return
	}
	r.response = response
	r.result = true
	return
}
