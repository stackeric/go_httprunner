package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const (
	//GET is Http Get
	GET string = "GET"
	//POST is Http post
	POST string = "POST"
)

// CaseConfig contain config info in a case
type CaseConfig struct {
	Name      string   `yaml:"name"`
	BaseURL   string   `yaml:"base_url"`
	Variables Variable `yaml:"variables"`
}

//IsValid check
func (c *CaseConfig) IsValid() bool {
	if c.Name == "" || c.BaseURL == "" {
		return false
	}
	return true
}

// TestCase is the minimal test steps
type TestCase struct {
	CaseCtx   Context
	Config    CaseConfig `yaml:"config"`
	TestSteps []TestStep `yaml:"teststeps"`
}

//IsValid check
func (c *TestCase) IsValid() bool {
	if ok := c.Config.IsValid(); !ok {
		return false
	}
	for _, step := range c.TestSteps {
		if ok := step.IsValid(); !ok {
			return false
		}
	}
	return true
}

//NewTestCase read an yaml file ,parse it to TestCase
func NewTestCase(yamlPath string) (c TestCase, err error) {
	yamlFilePath := yamlPath
	err = c.LoadCaseFromYaml(yamlFilePath)
	if err != nil {
		return c, err
	}
	c.CaseCtx = Context{}
	err = c.Config.Variables.Parse(c.CaseCtx)
	if err != nil {
		return
	}
	return c, nil
}

// LoadCaseFromYaml read yaml file to TestCase
func (c *TestCase) LoadCaseFromYaml(yamlPath string) (err error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return
	}
	if !c.IsValid() {
		return ErrKeyNotFound
	}
	return
}

//Run run an test case
func (c *TestCase) Run() (r TestCaseReport, err error) {

	for n, step := range c.TestSteps {
		stepReport, _ := step.Run(n, c.Config.BaseURL, c.CaseCtx)
		r.AddProgress(n, stepReport)
	}
	r.SetResult(true)
	for _, v := range r.StepReports {
		if !v.result {
			r.SetResult(false)
			break
		}
	}
	return
}
