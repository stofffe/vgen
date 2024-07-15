package vgen

import (
	"encoding/json"
)

type ErrorMap map[string][]error

func (e *ErrorMap) AddError(key string, err error) {
	if *e == nil {
		*e = make(ErrorMap)
	}
	(*e)[key] = append((*e)[key], err)
}

func (e *ErrorMap) AddErrors(errors ErrorMap) {
	if *e == nil {
		*e = make(ErrorMap)
	}
	for k, v := range errors {
		(*e)[k] = append((*e)[k], v...)
	}
}

func EmptyErrorMap() ErrorMap {
	return nil
}

func NewErrorMap(prefix string, errs ...error) ErrorMap {
	var errors ErrorMap
	for _, err := range errs {
		errors.AddError(prefix, err)
	}
	return errors
}

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
