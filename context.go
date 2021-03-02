package main

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//Context contain all dynamic variable in test case
type Context map[string]interface{}

//Variable contain variable in test case  life
type Variable map[string]interface{}

// Parse return Variable actual value
func (v Variable) Parse(context Context) error {
	for key, value := range v {
		switch value.(type) {
		case string:
			if actualValue, ok := IsVariable(value.(string)); ok {
				rv, ok := context[actualValue]
				if !ok {
					return ErrVariableNotFound
				}
				v[key] = rv
			} else if actualValue, ok := IsFunc(value.(string)); ok {
				v[key] = actualValue
			}
		}
		context[key] = v[key]
	}
	return nil
}

// Parameters contain http api request parameters
type Parameters map[string]interface{}

// Parse return Variable actual value
func (v Parameters) Parse(context Context) error {
	for key, value := range v {
		switch value.(type) {
		case string:
			if actualValue, ok := IsVariable(value.(string)); ok {
				rv, ok := context[actualValue]
				if !ok {
					return ErrVariableNotFound
				}
				v[key] = rv
			} else if actualValue, ok := IsFunc(value.(string)); ok {
				v[key] = actualValue
			}
		}
	}
	return nil
}

// IsVariable check if string has variable
func IsVariable(value string) (string, bool) {
	// TODO Match ALL
	variablePatten := regexp.MustCompile(`\$\{(\w+)\}|\$(\w+)`)
	matched := variablePatten.FindAllStringIndex(value, 1)
	if matched == nil {
		return "", false
	}
	return value[matched[0][0]+1 : matched[0][1]], true
}

// IsFunc check if string has variable
func IsFunc(value string) (interface{}, bool) {
	funcPatten := regexp.MustCompile(`\$\{(\w+)\(([\$\w\.\-/\s=,]*)\)\}`)
	matched := funcPatten.FindAllStringIndex(value, 1)
	if matched == nil {
		return "", false
	}
	funcStringWithArgs := value[matched[0][0]+2 : matched[0][1]-1]

	funcStrPatten := regexp.MustCompile(`^(.*?)\((.*?)\)$`)
	funcMatch := funcStrPatten.FindAllStringSubmatchIndex(funcStringWithArgs, -1)
	i := funcMatch[0]
	// funcstr := funcStringWithArgs[i[0]:i[1]]
	funcname := funcStringWithArgs[i[2]:i[3]]
	args := funcStringWithArgs[i[4]:i[5]]
	splitArgs := strings.Split(args, ",")
	in := make([]interface{}, len(splitArgs))
	for i, v := range splitArgs {
		in[i] = strings.TrimSpace(v)
	}
	res, err := Call(funcname, in...)
	if err != nil {
		return "", false
	}
	return res, true
}

type stubMapping map[string]interface{}

//StubStorage contain all self define functions
var StubStorage = stubMapping{
	"sum_two": SumTwo,
}

// SumTwo sum two
func SumTwo(x, y int) int {
	return x + y
}

// Call call function by function name
func Call(funcName string, params ...interface{}) (result interface{}, err error) {
	f := reflect.ValueOf(StubStorage[funcName])
	t := f.Type()
	if len(params) != t.NumIn() {
		err = errors.New("params is out of index")
		return
	}
	// in := make([]reflect.Value, len(params))
	var callArgs []reflect.Value
	for i := 0; i < t.NumIn(); i++ {
		t := t.In(i)
		v := reflect.New(t).Elem()
		if i < len(params) {
			// Convert arg to type of v and set.
			arg := params[i]
			switch t.Kind() {
			case reflect.String:
				switch arg := arg.(type) {
				case string:
					v.SetString(arg)
				case []byte:
					v.SetString(string(arg))
				default:
					panic("not supported")
				}
			case reflect.Slice:
				if t.Elem() != reflect.TypeOf(byte(0)) {
					panic("not supported")
				}
				switch arg := arg.(type) {
				case string:
					v.SetBytes([]byte(arg))
				case []byte:
					v.SetBytes(arg)
				default:
					panic("not supported")
				}

			case reflect.Int:
				switch arg := arg.(type) {
				case int:
					v.SetInt(int64(arg))
				case string:
					i, err := strconv.ParseInt(arg, 10, 0)
					if err != nil {
						panic("bad int")
					}
					v.SetInt(i)
				default:
					panic("not supported")
				}
			default:
				panic("not supported")
			}
		}
		// Collect arguments for the call below.
		callArgs = append(callArgs, v)
	}
	// for k, param := range params {
	// 	in[k] = reflect.ValueOf(param)
	// }
	var res []reflect.Value
	res = f.Call(callArgs)
	result = res[0].Interface()
	return
}
