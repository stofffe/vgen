package vgen

import (
	"cmp"
	"fmt"
	"slices"
)

//
// Rule types
//

// Implement in order to use as validation rule
type Rule[T any] interface {
	Validate(field_name string, input T) ErrorMap
}

// Implements Rule interface
type RuleFunc[T any] func(field_name string, input T) ErrorMap

func (rule RuleFunc[T]) Validate(field_name string, input T) ErrorMap {
	return rule(field_name, input)
}

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
	if value == nil {
		if r.required {
			errors.AddErrors(NewErrorMap(field_name, fmt.Errorf("required")))
		}
		return errors
	}
	deref_value := *value

	for _, rule := range r.rules {
		if err := rule.Validate(field_name, deref_value); err != nil {
			errors.AddErrors(err)
		}
	}
	return errors
}

//
// Rule implementations
//

func Nested[T any, R Rule[T]](rules R) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		return rules.Validate(field_name, input)
	}
}

func List[T any](rules ...RuleFunc[T]) RuleFunc[[]T] {
	return func(field_name string, input []T) ErrorMap {
		errors := make(ErrorMap)
		for i, element := range input {
			for _, rule := range rules {
				err := rule(fmt.Sprintf("%s[%d]", field_name, i), element)
				errors.AddErrors(err)
			}
		}
		if errors.HasError() {
			return errors
		}
		return nil
	}
}

func MapValue[V any](value_rules ...RuleFunc[V]) RuleFunc[map[string]V] {
	return func(field_name string, input map[string]V) ErrorMap {
		errors := make(ErrorMap)
		for key, value := range input {
			for _, rule := range value_rules {
				value_err := rule(fmt.Sprintf("%s.%s", field_name, key), value)
				errors.AddErrors(value_err)
			}
		}
		if errors.HasError() {
			return errors
		}
		return nil
	}
}

func MapHasKey[T map[string]V, V any](key string) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if _, ok := input[key]; !ok {
			return NewErrorMap(field_name, fmt.Errorf("map must contain key %v", key))
		}
		return nil
	}
}

func PrefixMessage[T any](prefix string, rule RuleFunc[T]) RuleFunc[T] { // TODO change rulefunc
	return func(field_name string, input T) ErrorMap {
		prefixedErrors := make(ErrorMap)
		errors := rule(field_name, input)

		if errors.HasError() {
			for key, errs := range errors {
				for _, err := range errs {
					prefixedErrors[key] = append(prefixedErrors[key], fmt.Errorf("%s%s", prefix, err.Error()))
				}
			}
		}
		if errors.HasError() {
			return errors
		}
		return nil
	}
}

func CustomMessage[T any](message string, rule RuleFunc[T]) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		errs := rule(field_name, input)

		if errs.HasError() {
			return NewErrorMap(field_name, fmt.Errorf("%s", message))
		}

		return nil
	}
}

func Eq[T comparable](value T) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if input != value {
			return NewErrorMap(field_name, fmt.Errorf("must be equal to %v", value))
		}
		return nil
	}
}

func Neq[T comparable](value T) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if input == value {
			return NewErrorMap(field_name, fmt.Errorf("must not be equal to %v", value))
		}
		return nil
	}
}

func Gt[T cmp.Ordered](value T) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !(input > value) {
			return NewErrorMap(field_name, fmt.Errorf("must be greater than %v", value))
		}
		return nil
	}
}

func Gte[T cmp.Ordered](value T) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !(input >= value) {
			return NewErrorMap(field_name, fmt.Errorf("must be greater than or equal to %v", value))
		}
		return nil
	}
}

func Lt[T cmp.Ordered](value T) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !(input < value) {
			return NewErrorMap(field_name, fmt.Errorf("must be less than %v", value))
		}
		return nil
	}
}

func Lte[T cmp.Ordered](value T) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !(input <= value) {
			return NewErrorMap(field_name, fmt.Errorf("must be less than or equal to %v", value))
		}
		return nil
	}
}
func LenEq[T []V, V any](value int) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !(len(input) == value) {
			return NewErrorMap(field_name, fmt.Errorf("len must be equal to %v", value))
		}
		return nil
	}
}
func LenGt[T []V, V any](value int) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !(len(input) > value) {
			return NewErrorMap(field_name, fmt.Errorf("len must be greater than %v", value))
		}
		return nil
	}
}
func LenGte[T []V, V any](value int) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !(len(input) >= value) {
			return NewErrorMap(field_name, fmt.Errorf("len must be greater than or equal to %v", value))
		}
		return nil
	}
}
func LenLt[T []V, V any](value int) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !(len(input) < value) {
			return NewErrorMap(field_name, fmt.Errorf("len must be less than %v", value))
		}
		return nil
	}
}
func LenLte[T []V, V any](value int) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !(len(input) <= value) {
			return NewErrorMap(field_name, fmt.Errorf("len must be less than or equal to %v", value))
		}
		return nil
	}
}

func OneOf[T comparable](values ...T) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if !slices.Contains(values, input) {
			return NewErrorMap(field_name, fmt.Errorf("must be one of %v", values))
		}
		return nil
	}
}

func NotOneOf[T comparable](values ...T) RuleFunc[T] {
	return func(field_name string, input T) ErrorMap {
		if slices.Contains(values, input) {
			return NewErrorMap(field_name, fmt.Errorf("must not be one of %v", values))
		}
		return nil
	}
}

// // Currently no distinction between key and value errors
// func MapKey[T map[string]V, V any](key_rules ...Rule[string]) Rule[map[string]V] {
// 	return func(field_name string, input map[string]V) ErrorMap {
// 		errors := make(ErrorMap)
// 		for key := range input {
// 			for _, rule := range key_rules {
// 				key_err := rule(fmt.Sprintf("%s.%s", field_name, key), key).Prefix("key: ")
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
