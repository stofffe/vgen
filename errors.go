package vgen

import (
	"encoding/json"
)

type ErrorList []error
type ErrorMap map[string]ErrorList

func SingleError(prefix string, err error) ErrorMap {
	errors := make(ErrorMap)
	errors[prefix] = ErrorList{err}
	return errors
}

func (e ErrorMap) AddErrors(errors ErrorMap) {
	for k, v := range errors {
		e[k] = append(e[k], v...)
	}
}

func (e ErrorMap) Debug() string {
	b, _ := json.MarshalIndent(e, "", "  ")
	return string(b)
}

func (e ErrorList) MarshalJSON() ([]byte, error) {
	list := []string{}
	for _, err := range e {
		list = append(list, err.Error())
	}
	return json.Marshal(list)
}
