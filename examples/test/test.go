package main

import (
	"fmt"
	"log"
)

// vgen(i)
type (
	Test1 struct {
		field1 string
	}
	Test2 struct {
		field1 string
	}
	Test3 string
)

type ParseError struct {
	detailed error
	inner    error
}

func (e ParseError) Error() string {
	return e.inner.Error()
}

func test() error {
	err := nested()
	if err != nil {
		return ParseError{
			detailed: fmt.Errorf("could not do nested: %v", err),
			inner:    err,
		}
	}
	return nil
}

func nested() error {
	return fmt.Errorf("this is the inner error")
}

func main() {
	verbose := false
	err := test()

	if err != nil {
		switch err := err.(type) {
		case ParseError:
			if verbose {
				log.Fatal(err.detailed)
			} else {
				log.Fatal(err.inner)
			}
		default:
			log.Fatal(err)
		}
	}

}
