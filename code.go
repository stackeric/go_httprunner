package main

import "errors"

var (
	//ErrKeyNotFound occur when parse yaml file
	ErrKeyNotFound = errors.New("you miss some required key")
	//ErrValidate occur validator failed
	ErrValidate = errors.New("oh, got unexpected value")
	// ErrVariableNotFound occur when variable not exist in context
	ErrVariableNotFound = errors.New("variable not exist")
)
