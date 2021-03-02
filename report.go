package main

import "net/http"

// Report is an report interface for test case or suite
type Report interface {
	Result() bool
}

// TestCaseReport record test case progress
type TestCaseReport struct {
	result      bool
	StepReports []StepReport
}

// StepReport record every step info
type StepReport struct {
	index      int
	url        string
	method     string
	result     bool
	statusCode int
	req        *http.Request
	response   Response
	err        error
}

// Result return all result info
func (r *TestCaseReport) Result() bool {
	return r.result
}

// SetResult return all result info
func (r *TestCaseReport) SetResult(nr bool) {
	r.result = nr
}

// AddProgress add step info
func (r *TestCaseReport) AddProgress(index int, report StepReport) {
	r.StepReports = append(r.StepReports, report)
}
