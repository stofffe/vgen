package vgen

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

//
// Rule types
//

type Rule[T any] func(field_name string, input T) ErrorMap

type Rules[T any] struct {
	required bool
	rules    []Rule[T]
}

func NewRules[T any](required bool, rules ...Rule[T]) Rules[T] {
	return Rules[T]{
		required: required,
		rules:    rules,
	}
}

func (r Rules[T]) Validate(field_name string, value *T) ErrorMap {
	errors := make(ErrorMap)
	if r.required && value == nil {
		errors.AddErrors(SingleError(field_name, fmt.Errorf("required")))
		return errors
	}
	deref_value := *value

	for _, rule := range r.rules {
		if err := rule(field_name, deref_value); err != nil {
			errors.AddErrors(err)
		}
	}
	return errors
}

//
// Rule implementations
//

func List[T any](rules ...Rule[T]) Rule[[]T] {
	return func(field_name string, input []T) ErrorMap {
		errors := make(ErrorMap)
		for i, element := range input {
			for _, rule := range rules {
				err := rule(fmt.Sprintf("%s[%d]", field_name, i), element)
				errors.AddErrors(err)
			}
		}
		if len(errors) > 0 {
			return errors
		}
		return nil
	}
}

func MapValue[V any](value_rules ...Rule[V]) Rule[map[string]V] {
	return func(field_name string, input map[string]V) ErrorMap {
		errors := make(ErrorMap)
		for key, value := range input {
			for _, rule := range value_rules {
				value_err := rule(fmt.Sprintf("%s.%s", field_name, key), value)
				errors.AddErrors(value_err)
			}
		}
		if len(errors) > 0 {
			return errors
		}
		return nil
	}
}

func CustomMessage[T any](message string, rule Rule[T]) Rule[T] {
	return func(field_name string, input T) ErrorMap {
		err := rule(field_name, input)

		if len(err) > 0 {
			return SingleError(field_name, errors.New(message))
		}

		return nil
	}
}

func Eq[T comparable](value T) Rule[T] {
	return func(field_name string, input T) ErrorMap {
		if input != value {
			return SingleError(field_name, fmt.Errorf("not equal to %v", value))
		}
		return nil
	}
}

func Gt[T cmp.Ordered](value T) Rule[T] {
	return func(field_name string, input T) ErrorMap {
		if !(input > value) {
			return SingleError(field_name, fmt.Errorf("not greater than %v", value))
		}
		return nil
	}
}

func Gte[T cmp.Ordered](value T) Rule[T] {
	return func(field_name string, input T) ErrorMap {
		if !(input >= value) {
			return SingleError(field_name, fmt.Errorf("not greater than or equal to %v", value))
		}
		return nil
	}
}

func Lt[T cmp.Ordered](value T) Rule[T] {
	return func(field_name string, input T) ErrorMap {
		if !(input < value) {
			return SingleError(field_name, fmt.Errorf("not less than %v", value))
		}
		return nil
	}
}

func Lte[T cmp.Ordered](value T) Rule[T] {
	return func(field_name string, input T) ErrorMap {
		if !(input <= value) {
			return SingleError(field_name, fmt.Errorf("not less than or equal to %v", value))
		}
		return nil
	}
}
func LenEq[T any](value int) Rule[[]T] {
	return func(field_name string, input []T) ErrorMap {
		if !(len(input) == value) {
			return SingleError(field_name, fmt.Errorf("len not equal to %v", value))
		}
		return nil
	}
}
func LenGt[T any](value int) Rule[[]T] {
	return func(field_name string, input []T) ErrorMap {
		if !(len(input) > value) {
			return SingleError(field_name, fmt.Errorf("len not greater than %v", value))
		}
		return nil
	}
}
func LenGte[T any](value int) Rule[[]T] {
	return func(field_name string, input []T) ErrorMap {
		if !(len(input) >= value) {
			return SingleError(field_name, fmt.Errorf("len not greater than or equal to %v", value))
		}
		return nil
	}
}
func LenLt[T any](value int) Rule[[]T] {
	return func(field_name string, input []T) ErrorMap {
		if !(len(input) < value) {
			return SingleError(field_name, fmt.Errorf("len not less than %v", value))
		}
		return nil
	}
}
func LenLte[T any](value int) Rule[[]T] {
	return func(field_name string, input []T) ErrorMap {
		if !(len(input) <= value) {
			return SingleError(field_name, fmt.Errorf("len not less than or equal to %v", value))
		}
		return nil
	}
}

func OneOf[T comparable](values ...T) Rule[T] {
	return func(field_name string, input T) ErrorMap {
		if !slices.Contains(values, input) {
			return SingleError(field_name, fmt.Errorf("not one of %v", values))
		}
		return nil
	}
}

func NotOneOf[T comparable](values ...T) Rule[T] {
	return func(field_name string, input T) ErrorMap {
		if slices.Contains(values, input) {
			return SingleError(field_name, fmt.Errorf("is one of %v", values))
		}
		return nil
	}
}

// Currently no distinction between key and value errors

// func MapKey[V any](key_rules ...Rule[string]) Rule[map[string]V] {
// 	return func(field_name string, input map[string]V) ErrorMap {
// 		errors := make(ErrorMap)
// 		for key := range input {
// 			for _, rule := range key_rules {
// 				key_err := rule(fmt.Sprintf("%s.%s", field_name, key), key)
// 				errors.AddErrors(key_err)
// 			}
// 		}
// 		if len(errors) > 0 {
// 			return errors
// 		}
// 		return nil
// 	}
// }
//
// func MapKeyValue[V any](key_rules []Rule[string], value_rules []Rule[V]) Rule[map[string]V] {
// 	return func(field_name string, input map[string]V) ErrorMap {
// 		errors := make(ErrorMap)
// 		for key, value := range input {
// 			for _, rule := range key_rules {
// 				key_err := rule(fmt.Sprintf("%s.%s", field_name, key), key)
// 				errors.AddErrors(key_err)
// 			}
// 			for _, rule := range value_rules {
// 				value_err := rule(fmt.Sprintf("%s.%s", field_name, key), value)
// 				errors.AddErrors(value_err)
// 			}
// 		}
// 		if len(errors) > 0 {
// 			return errors
// 		}
// 		return nil
// 	}
// }
