package cli

import (
	"errors"
	"fmt"
	"strings"
)

type InfoMessage struct {
	msg  string
	path string
}

func (g InfoMessage) Format() string {
	var builder strings.Builder
	builder.WriteString("[INFO]")
	if g.path != "" {
		builder.WriteString(" ")
		builder.WriteString(g.path)
	}
	builder.WriteString(": ")
	builder.WriteString(g.msg)
	return builder.String()
}

type WarningMessage struct {
	msg  string
	path string
}

func (g WarningMessage) Format() string {
	var builder strings.Builder
	builder.WriteString("[WARNING]")
	if g.path != "" {
		builder.WriteString(" ")
		builder.WriteString(g.path)
	}
	builder.WriteString(": ")
	builder.WriteString(g.msg)
	return builder.String()
}

type ErrorMessage struct {
	err  error
	path string
}

func (g ErrorMessage) Format(detailed bool) string {
	var builder strings.Builder
	builder.WriteString("[ERROR]")
	if g.path != "" {
		builder.WriteString(" ")
		builder.WriteString(g.path)
	}

	var derr DetailedError
	if errors.As(g.err, &derr) {
		if derr.line != nil {
			builder.WriteString(" ")
			builder.WriteString(fmt.Sprintf("line %d", *derr.line))
		}
		builder.WriteString(": ")
		if detailed {
			builder.WriteString(g.err.Error())
		} else {
			builder.WriteString(derr.msg)
		}
	} else {
		builder.WriteString(": ")
		if detailed {
			builder.WriteString(g.err.Error())
		} else {
			builder.WriteString("internal error (use -v flag for more information)")
		}
	}
	return builder.String()
}

type DetailedError struct {
	msg  string
	err  error
	line *int
}

func (e DetailedError) Error() string {
	return e.err.Error()
}
