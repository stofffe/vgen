package cli

import (
	"fmt"
	"strings"
)

type DetailedError struct {
	inner    error
	detailed error
}

func NewInternalError(err error) DetailedError {
	return DetailedError{
		inner:    fmt.Errorf("internal error"),
		detailed: err,
	}
}

func NewDetailedError(err error, trace string) DetailedError {
	if a, ok := err.(DetailedError); ok {
		return DetailedError{
			inner:    a.inner,
			detailed: fmt.Errorf("%s: %s", trace, a.detailed),
		}
	}
	return DetailedError{
		inner:    err,
		detailed: fmt.Errorf("%s: %s", trace, err),
	}
}

func (e DetailedError) Error() string {
	var builder strings.Builder
	builder.WriteString(e.inner.Error())
	return builder.String()
}
