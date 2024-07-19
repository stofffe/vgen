package cli

import (
	"errors"
	"fmt"
	"strings"
)

type InfoMessage struct {
	path string
	info string
}

func (g InfoMessage) Format() string {
	return fmt.Sprintf("[INFO] %s: %s", g.path, g.info)
}

type WarningMessage struct {
	warning string
	path    string
}

func (g WarningMessage) Format() string {
	return fmt.Sprintf("[WARNING] %s: %s", g.path, g.warning)
}

type ErrorMessage struct {
	err  error
	path string
}

func (g ErrorMessage) Format(detailed bool) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[ERROR] %s", g.path))

	var derr DetailedError
	if errors.As(g.err, &derr) {
		if derr.pos != nil {
			builder.WriteString(fmt.Sprintf(" line %d", *derr.pos))
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
	msg string
	err error
	pos *int
}

func (e DetailedError) Error() string {
	return e.err.Error()
}
