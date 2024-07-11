package vgen

import (
	"encoding/json"
)

type ErrorMap map[string]ErrorList

func (e ErrorMap) AddError(key string, err error) {
	e[key] = append(e[key], err)
}

func (e ErrorMap) AddErrors(errors ErrorMap) {
	for k, v := range errors {
		e[k] = append(e[k], v...)
	}
}

func (e ErrorMap) HasError() bool {
	if e == nil {
		return false
	}
	return len(e) > 0
}

func EmptyErrorMap() ErrorMap {
	return make(ErrorMap)
}

// func SingleError(key string, err error) ErrorMap {
// 	errors := EmptyErrorMap()
// 	errors[key] = ErrorList{err}
// 	return errors
// }

func NewErrorMap(key string, errs ...error) ErrorMap {
	errors := EmptyErrorMap()
	errors[key] = append(errors[key], errs...)
	return errors
}

// func (e ErrorMap) Prefix(prefix string) ErrorMap {
// 	errors := make(ErrorMap)
// 	for k, v := range e {
// 		errors[prefix+k] = v
// 	}
// 	return errors
// }

func (e ErrorMap) Debug() string {
	b, _ := json.MarshalIndent(e, "", "  ")
	return string(b)
}

type ErrorList []error

func (e ErrorList) MarshalJSON() ([]byte, error) {
	list := []string{}
	for _, err := range e {
		list = append(list, err.Error())
	}
	return json.Marshal(list)
}
