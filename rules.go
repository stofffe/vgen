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
	Validate(fieldName string, input T) ErrorMap
}

type Rules[T any] struct {
	Required bool
	Rules    []Rule[T]
}

//
// type Rules[T any] []Rule[T]
//

func (r Rules[T]) Validate(fieldName string, input *T) ErrorMap {
	errors := EmptyErrorMap()

	// required
	if input == nil {
		if r.Required {
			errors.AddError(fieldName, fmt.Errorf("required"))
		}
		return errors
	}
	value := *input

	for _, rule := range r.Rules {
		errors.AddErrors(rule.Validate(fieldName, value))
	}

	return errors
}

func RulesOptional[T any](rules ...Rule[T]) Rules[T] {
	return Rules[T]{
		Required: false,
		Rules:    rules,
	}
}
func RulesRequired[T any](rules ...Rule[T]) Rules[T] {
	return Rules[T]{
		Required: true,
		Rules:    rules,
	}
}

// func RulesDefault[T any](defaultValue T, rules ...Rule[T]) Rules[T] {
// 	return Rules[T]{
// 		Required: false,
// 		Rules:    rules,
// 	}
// }

// Implements Rule interface
type RuleFunc[T any] func(fieldName string, input T) ErrorMap

func (rule RuleFunc[T]) Validate(fieldName string, input T) ErrorMap {
	return rule(fieldName, input)
}

//
// Rule implementations
//

// Validate rules to for internal value of pointer
func Deref[T any](rules ...Rule[T]) RuleFunc[*T] {
	return func(fieldName string, input *T) ErrorMap {
		errors := EmptyErrorMap()
		if input == nil {
			return errors
		}
		for _, rule := range rules {
			errors.AddErrors(rule.Validate(fieldName, *input))
		}
		return errors
	}
}

func Nested[T any, R Rule[T]](rules R) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		errors.AddErrors(rules.Validate(fieldName, input))
		return errors
	}
}

func List[T any](rules ...Rule[T]) RuleFunc[[]T] {
	return func(fieldName string, input []T) ErrorMap {
		errors := EmptyErrorMap()
		for i, element := range input {
			for _, rule := range rules {
				errors.AddErrors(rule.Validate(fmt.Sprintf("%s[%d]", fieldName, i), element))
			}
		}
		return errors
	}
}

func MapValue[V any](value_rules ...Rule[V]) RuleFunc[map[string]V] {
	return func(fieldName string, input map[string]V) ErrorMap {
		errors := EmptyErrorMap()
		for key, value := range input {
			for _, rule := range value_rules {
				errors.AddErrors(rule.Validate(fmt.Sprintf("%s.%s", fieldName, key), value))
			}
		}
		return errors
	}
}

func MapHasKey[T map[string]V, V any](key string) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if _, ok := (input)[key]; !ok {
			errors.AddError(fieldName, fmt.Errorf("map must contain key %v", key))
		}
		return errors
	}
}

func PrefixMessage[T any](prefix string, rule Rule[T]) RuleFunc[T] { // TODO change rulefunc
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()

		innerErrors := rule.Validate(fieldName, input)

		for key, errs := range innerErrors {
			for _, err := range errs {
				errors[key] = append(errors[key], fmt.Errorf("%s%s", prefix, err.Error()))
			}
		}
		return innerErrors
	}
}

func CustomMessage[T any](message string, rule Rule[T]) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		innerErrors := rule.Validate(fieldName, input)

		if innerErrors != nil {
			errors.AddError(fieldName, fmt.Errorf("%s", message))
		}

		return errors
	}
}

func Eq[T comparable](value T) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if input != value {
			errors.AddError(fieldName, fmt.Errorf("must be equal to %v", value))
		}
		return errors
	}
}

func Neq[T comparable](value T) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if input == value {
			errors.AddError(fieldName, fmt.Errorf("must not be equal to %v", value))
		}
		return errors
	}
}

func Gt[T cmp.Ordered](value T) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !(input > value) {
			errors.AddError(fieldName, fmt.Errorf("must be greater than %v", value))
		}
		return errors
	}
}

func Gte[T cmp.Ordered](value T) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !(input >= value) {
			errors.AddError(fieldName, fmt.Errorf("must be greater than or equal to %v", value))
		}
		return errors
	}
}

func Lt[T cmp.Ordered](value T) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !(input < value) {
			errors.AddError(fieldName, fmt.Errorf("must be less than %v", value))
		}
		return errors
	}
}

func Lte[T cmp.Ordered](value T) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !(input <= value) {
			errors.AddError(fieldName, fmt.Errorf("must be less than or equal to %v", value))
		}
		return errors
	}
}
func LenEq[T []V, V any](value int) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !(len(input) == value) {
			errors.AddError(fieldName, fmt.Errorf("len must be equal to %v", value))
		}
		return errors
	}
}
func LenGt[T []V, V any](value int) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !(len(input) > value) {
			errors.AddError(fieldName, fmt.Errorf("len must be greater than %v", value))
		}
		return errors
	}
}
func LenGte[T []V, V any](value int) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !(len(input) >= value) {
			errors.AddError(fieldName, fmt.Errorf("len must be greater than or equal to %v", value))
		}
		return errors
	}
}
func LenLt[T []V, V any](value int) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !(len(input) < value) {
			errors.AddError(fieldName, fmt.Errorf("len must be less than %v", value))
		}
		return errors
	}
}
func LenLte[T []V, V any](value int) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !(len(input) <= value) {
			errors.AddError(fieldName, fmt.Errorf("len must be less than or equal to %v", value))
		}
		return errors
	}
}

func OneOf[T comparable](values ...T) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if !slices.Contains(values, input) {
			errors.AddError(fieldName, fmt.Errorf("must be one of %v", values))
		}
		return errors
	}
}

func NotOneOf[T comparable](values ...T) RuleFunc[T] {
	return func(fieldName string, input T) ErrorMap {
		errors := EmptyErrorMap()
		if slices.Contains(values, input) {
			errors.AddError(fieldName, fmt.Errorf("must not be one of %v", values))
		}
		return errors
	}
}

// // Currently no distinction between key and value errors
// func MapKey[T map[string]V, V any](key_rules ...Rule[string]) Rule[map[string]V] {
// 	return func(fieldName string, input map[string]V) ErrorMap {
// 		errors := make(ErrorMap)
// 		for key := range input {
// 			for _, rule := range key_rules {
// 				key_err := rule(fmt.Sprintf("%s.%s", fieldName, key), key).Prefix("key: ")
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
// 	return func(fieldName string, input map[string]V) ErrorMap {
// 		errors := make(ErrorMap)
// 		for key, value := range input {
// 			for _, rule := range key_rules {
// 				key_err := rule(fmt.Sprintf("%s.%s", fieldName, key), key)
// 				errors.AddErrors(key_err)
// 			}
// 			for _, rule := range value_rules {
// 				value_err := rule(fmt.Sprintf("%s.%s", fieldName, key), value)
// 				errors.AddErrors(value_err)
// 			}
// 		}
// 		if len(errors) > 0 {
// 			return errors
// 		}
// 		return nil
// 	}
// }
// func PointerNotNil[T *I, I any]() RuleFunc[T] {
// 	return func(fieldName string, input *T) ErrorMap {
// 		errors := EmptyErrorMap()
// 		if input == nil || *input == nil {
// 			errors.AddError(fieldName, fmt.Errorf("must not be null"))
// 		}
// 		return errors
// 	}
// }

// List of rules along with required check
// type Rules[T any] struct {
// 	required bool
// 	rules    []Rule[T]
// }
//
// func NewRules[T any](required bool, rules ...Rule[T]) Rules[T] {
// 	return Rules[T]{
// 		required: required,
// 		rules:    rules,
// 	}
// }
// func RequiredRules[T any](rules ...Rule[T]) Rules[T] {
// 	return Rules[T]{
// 		required: true,
// 		rules:    rules,
// 	}
// }
// func OptionalRules[T any](rules ...Rule[T]) Rules[T] {
// 	return Rules[T]{
// 		required: false,
// 		rules:    rules,
// 	}
// }
//
// func (r Rules[T]) Validate(fieldName string, value *T) ErrorMap {
// 	errors := make(ErrorMap)
// 	if value == nil {
// 		if r.required {
// 			errors.AddErrors(NewErrorMap(fieldName, fmt.Errorf("required")))
// 		}
// 		return errors
// 	}
// 	deref_value := *value
//
// 	for _, rule := range r.rules {
// 		if err := rule.Validate(fieldName, deref_value); err != nil {
// 			errors.AddErrors(err)
// 		}
// 	}
// 	return errors
// }
