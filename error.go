package docopt

import (
	"fmt"
)

type errorType int

const (
	errorUser errorType = iota
	errorLanguage
)

func (e errorType) String() string {
	switch e {
	case errorUser:
		return "errorUser"
	case errorLanguage:
		return "errorLanguage"
	}
	return ""
}

// UserError records an error with program arguments.
type UserError struct {
	msg   string
	Usage string
}

func (e UserError) Error() string {
	return e.msg
}
func newUserError(msg string, f ...interface{}) error {
	return &UserError{fmt.Sprintf(msg, f...), ""}
}

// LanguageError records an error with the doc string.
type LanguageError struct {
	msg string
}

func (e LanguageError) Error() string {
	return e.msg
}
func newLanguageError(msg string, f ...interface{}) error {
	return &LanguageError{fmt.Sprintf(msg, f...)}
}

var newError = fmt.Errorf

// newUnmatchedError builds a descriptive error when pattern matching fails.
// It inspects leftover patterns to identify unknown options or arguments.
func newUnmatchedError(matched bool, left *patternList) error {
	if left != nil && len(*left) > 0 {
		var unknown []string
		for _, p := range *left {
			if p.t&patternOption != 0 {
				if p.long != "" {
					unknown = append(unknown, p.long)
				} else if p.short != "" {
					unknown = append(unknown, p.short)
				}
			} else if p.t&patternArgument != 0 {
				if val, ok := p.value.(string); ok && val != "" {
					unknown = append(unknown, val)
				} else if p.name != "" {
					unknown = append(unknown, p.name)
				}
			}
		}
		if len(unknown) == 1 {
			return newUserError("unknown option or argument: %s", unknown[0])
		}
		if len(unknown) > 1 {
			return newUserError("unknown options or arguments: %s", fmt.Sprintf("%v", unknown))
		}
	}
	return newUserError("")
}
