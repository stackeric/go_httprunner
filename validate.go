package main

import (
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
)

// ConstValidator contain value validator from http response
var ConstValidator = map[string]string{
	"status_code": "",
}

//Comparer is assert name
type Comparer string

// Do parse assert type ,and do actual assert
func (c Comparer) Do(raw interface{}, expected interface{}) bool {
	switch c {
	case "eq":
		return reflect.DeepEqual(raw, expected)
	default:
		return false
	}
}

//Validator read json, extract actual value, compare expected value
type Validator struct {
	Compare  Comparer    `yaml:"compare"`
	Key      string      `yaml:"key"`
	Expected interface{} `yaml:"expected"`
}

//Run to compare result
func (v *Validator) Run(body []byte) (string, error) {
	value := gjson.GetBytes(body, v.Key)
	if !value.Exists() {
		return fmt.Sprintf("%s not found", v.Key), ErrKeyNotFound
	}
	var ok bool
	switch v.Expected.(type) {
	case int, int32:
		ok = v.Compare.Do(value.Int(), int64(v.Expected.(int)))
	case float32, float64:
		ok = v.Compare.Do(value.Float(), v.Expected.(float64))
	case bool:
		ok = v.Compare.Do(value.Bool(), v.Expected.(bool))
	default:
		ok = v.Compare.Do(value.String(), v.Expected.(string))
	}
	if !ok {
		return fmt.Sprintf("got %v want %v", value.Raw, v.Expected), ErrValidate
	}
	return "", nil
}
