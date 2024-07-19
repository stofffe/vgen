package cli

import (
	"errors"
	"fmt"
)

type DetailedError struct {
	msg string
	err error
}

func (e DetailedError) Error() string {
	return e.err.Error()
}

type InfoMessage struct {
	path string
	info string
}
type WarningMessage struct {
	warning string
	path    string
}
type ErrorMessage struct {
	err  error
	path string
}

func (g InfoMessage) Format() string {
	return fmt.Sprintf("[INFO] %s: %s", g.path, g.info)
}
func (g WarningMessage) Format() string {
	return fmt.Sprintf("[WARNING] %s: %s", g.path, g.warning)
}
func (g ErrorMessage) Format(detailed bool) string {
	var derr DetailedError
	if errors.As(g.err, &derr) {
		if detailed {
			return fmt.Sprintf("[ERROR] %s: %s", g.path, g.err.Error())
		} else {
			return fmt.Sprintf("[ERROR] %s: %s", g.path, derr.msg)
		}
	} else {
		if detailed {
			return fmt.Sprintf("[ERROR] %s", g.err.Error())
		} else {
			return fmt.Sprintf("[ERROR] internal error (use -v flag for more information)")
		}
	}
}
