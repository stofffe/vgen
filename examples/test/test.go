package main

import (
	"errors"
	"fmt"
)

// vgen(i)
type CustomError struct {
	msg string
}

func (e CustomError) Error() string {
	return e.msg
}

func main() {
	err := a()
	if err != nil {
		var customErr CustomError
		if errors.As(err, &customErr) {
			fmt.Println("custom", customErr)
		}
		innerErr := errors.Unwrap(errors.Unwrap(errors.Unwrap(err)))
		fmt.Println("inner", innerErr)
	}
}

func a() error {
	err := b()
	if err != nil {
		return fmt.Errorf("error in b: %w", err)
	}
	return nil
}
func b() error {
	err := c()
	if err != nil {
		return fmt.Errorf("error in c: %w", err)
	}
	return nil
}
func c() error {
	return CustomError{
		msg: "my error",
	}

}
