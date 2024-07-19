package vgen

import (
	"cmp"
	"fmt"

	"github.com/stofffe/vgen/internal/utils"
)

//
// Rule types
//

// Implement in order to use as validation rule
type Rule[T any] interface {
	Validate(key string, input T) ErrorMap
}

type Rules[T any] struct {
	Required bool
	Rules    []Rule[T]
}

// Validate rules with required check
func (r Rules[T]) Validate(key string, input *T) ErrorMap {
	var errors ErrorMap

	if input == nil {
		if r.Required {
			errors.AddError(key, fmt.Errorf("required"))
		}
		return errors
	}

	value := *input
	for _, rule := range r.Rules {
		errors.AddErrors(rule.Validate(key, value))
	}

	return errors
}

// New rules without required check
func RulesOptional[T any](rules ...Rule[T]) Rules[T] {
	return Rules[T]{
		Required: false,
		Rules:    rules,
	}
}

// New rules with required check
func RulesRequired[T any](rules ...Rule[T]) Rules[T] {
	return Rules[T]{
		Required: true,
		Rules:    rules,
	}
}

// Implements Rule interface
type RuleFunc[T any] func(key string, input T) ErrorMap

func (rule RuleFunc[T]) Validate(key string, input T) ErrorMap {
	return rule(key, input)
}

//
// Rule implementations
//

// Validation rules for inner type of pointer
func Deref[T any](rules ...Rule[T]) RuleFunc[*T] {
	return func(key string, input *T) ErrorMap {
		var errors ErrorMap
		if input == nil {
			return errors
		}
		for _, rule := range rules {
			errors.AddErrors(rule.Validate(key, *input))
		}
		return errors
	}
}

// Validate nested Vgen struct
func Nested[T any, R Rule[T]](rules R) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		errors.AddErrors(rules.Validate(key, input))
		return errors
	}
}

// Prefixes the error
func PrefixMessage[T any](prefix string, rule Rule[T]) RuleFunc[T] { // TODO change rulefunc
	return func(key string, input T) ErrorMap {
		innerErrors := rule.Validate(key, input)

		if innerErrors == nil {
			return nil
		}

		var errors ErrorMap
		for key, errs := range innerErrors {
			for _, err := range errs {
				errors.AddError(key, fmt.Errorf("%s%s", prefix, err.Error()))
			}
		}
		return innerErrors
	}
}

// Replaces any errors with a custom message
func CustomMessage[T any](message string, rule Rule[T]) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		innerErrors := rule.Validate(key, input)

		if innerErrors != nil {
			errors.AddError(key, fmt.Errorf("%s", message))
		}

		return errors
	}
}

// Validates each element of list
func List[T any](rules ...Rule[T]) RuleFunc[[]T] {
	return func(key string, input []T) ErrorMap {
		var errors ErrorMap
		for i, element := range input {
			for _, rule := range rules {
				errors.AddErrors(rule.Validate(fmt.Sprintf("%s[%d]", key, i), element))
			}
		}
		return errors
	}
}

// Validates each value element of map
func MapValue[V any](value_rules ...Rule[V]) RuleFunc[map[string]V] {
	return func(key string, input map[string]V) ErrorMap {
		var errors ErrorMap
		for key, value := range input {
			for _, rule := range value_rules {
				errors.AddErrors(rule.Validate(fmt.Sprintf("%s.%s", key, key), value))
			}
		}
		return errors
	}
}

// Checks if map includes specific key
func MapHasKey[T map[string]V, V any](key string) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if _, ok := (input)[key]; !ok {
			errors.AddError(key, fmt.Errorf("map must contain key %v", key))
		}
		return errors
	}
}

// Equals
func Eq[T comparable](value T) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if input != value {
			errors.AddError(key, fmt.Errorf("must be equal to %v", value))
		}
		return errors
	}
}

// Not equals
func Neq[T comparable](value T) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if input == value {
			errors.AddError(key, fmt.Errorf("must not be equal to %v", value))
		}
		return errors
	}
}

// Greater than
func Gt[T cmp.Ordered](value T) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !(input > value) {
			errors.AddError(key, fmt.Errorf("must be greater than %v", value))
		}
		return errors
	}
}

// Greater than or equal to
func Gte[T cmp.Ordered](value T) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !(input >= value) {
			errors.AddError(key, fmt.Errorf("must be greater than or equal to %v", value))
		}
		return errors
	}
}

// Less than
func Lt[T cmp.Ordered](value T) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !(input < value) {
			errors.AddError(key, fmt.Errorf("must be less than %v", value))
		}
		return errors
	}
}

// Less than or equal to
func Lte[T cmp.Ordered](value T) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !(input <= value) {
			errors.AddError(key, fmt.Errorf("must be less than or equal to %v", value))
		}
		return errors
	}
}

// Len Equals
func LenEq[T []V, V any](value int) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !(len(input) == value) {
			errors.AddError(key, fmt.Errorf("len must be equal to %v", value))
		}
		return errors
	}
}

// Len greater than
func LenGt[T []V, V any](value int) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !(len(input) > value) {
			errors.AddError(key, fmt.Errorf("len must be greater than %v", value))
		}
		return errors
	}
}

// Len greater than or equal to
func LenGte[T []V, V any](value int) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !(len(input) >= value) {
			errors.AddError(key, fmt.Errorf("len must be greater than or equal to %v", value))
		}
		return errors
	}
}

// Len less than
func LenLt[T []V, V any](value int) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !(len(input) < value) {
			errors.AddError(key, fmt.Errorf("len must be less than %v", value))
		}
		return errors
	}
}

// Len less than or equal to
func LenLte[T []V, V any](value int) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !(len(input) <= value) {
			errors.AddError(key, fmt.Errorf("len must be less than or equal to %v", value))
		}
		return errors
	}
}

// One of the values
func OneOf[T comparable](values ...T) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if !utils.Contains(values, input) {
			errors.AddError(key, fmt.Errorf("must be one of %v", values))
		}
		return errors
	}
}

// Not one of the values
func NotOneOf[T comparable](values ...T) RuleFunc[T] {
	return func(key string, input T) ErrorMap {
		var errors ErrorMap
		if utils.Contains(values, input) {
			errors.AddError(key, fmt.Errorf("must not be one of %v", values))
		}
		return errors
	}
}
