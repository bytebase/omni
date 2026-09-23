package parser

import (
	"errors"
	"strconv"
)

// ParseError describes a single parse error with its source location.
//
// Position and End are byte offsets into the parsed text: the shape every
// omni engine's ParseError shares, so a consumer converts them to a line and
// column the same way for every engine (review.Index). Error returns just
// the message; ParseError is a pure data carrier.
type ParseError struct {
	Message string
	// Position is the byte offset where the error starts; End is the
	// exclusive end of the offending token.
	Position int
	End      int
}

// Error implements the error interface, returning the message only.
func (e *ParseError) Error() string {
	return e.Message
}

// ParseErrors is the error Parse returns when the input does not parse:
// every ParseError found, in source order, so a diagnostics consumer can
// report a multi-statement script's failures at once. errors.As with a
// *ParseError target yields the first; AllErrors yields the list.
type ParseErrors []ParseError

// Error implements the error interface with the first message and a count.
func (e ParseErrors) Error() string {
	switch len(e) {
	case 0:
		return "no parse errors"
	case 1:
		return e[0].Message
	default:
		return e[0].Message + " (and " + itoa(len(e)-1) + " more)"
	}
}

// Unwrap exposes each ParseError to errors.As and errors.Is.
func (e ParseErrors) Unwrap() []error {
	out := make([]error, len(e))
	for i := range e {
		out[i] = &e[i]
	}
	return out
}

// AllErrors returns every ParseError carried by err: the whole list behind
// Parse's ParseErrors, a single *ParseError, or nil for nil or any other
// error.
func AllErrors(err error) []ParseError {
	var all ParseErrors
	if errors.As(err, &all) {
		return all
	}
	var one *ParseError
	if errors.As(err, &one) {
		return []ParseError{*one}
	}
	return nil
}

// FirstError reduces Parse's error to the first ParseError it carries, for a
// caller that reports one position. Any other error is returned as is.
func FirstError(err error) error {
	if errs := AllErrors(err); len(errs) > 0 {
		return &errs[0]
	}
	return err
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
