package vgen

import (
	"encoding/json"
)

type ErrorMap map[string][]error

// Add single error
func (e *ErrorMap) AddError(key string, err error) {
	if *e == nil {
		*e = make(ErrorMap)
	}
	(*e)[key] = append((*e)[key], err)
}

// Transfer all errors from one map to another
func (e *ErrorMap) AddErrors(errors ErrorMap) {
	if errors == nil {
		return
	}
	if *e == nil {
		*e = make(ErrorMap)
	}
	for k, v := range errors {
		(*e)[k] = append((*e)[k], v...)
	}
}

// Create ErrorMap with multiple errors
func NewErrorMap(prefix string, errs ...error) ErrorMap {
	var errors ErrorMap
	for _, err := range errs {
		errors.AddError(prefix, err)
	}
	return errors
}

// Json format a ErrorMap
func (e ErrorMap) Debug() string {
	b, _ := json.MarshalIndent(e, "", "  ")
	return string(b)
}

func (e ErrorMap) MarshalJSON() ([]byte, error) {
	result := make(map[string][]string)
	for key, errors := range e {
		for _, err := range errors {
			result[key] = append(result[key], err.Error())
		}
	}
	return json.Marshal(result)
}
