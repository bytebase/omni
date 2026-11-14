package compat

import (
	"fmt"
	"regexp"
	"strings"
)

// A SQLServerError is one error SQL Server reported for a batch: its number
// and its message.
type SQLServerError struct {
	Number  int32
	Message string
}

// SQL Server's errors for a variable a batch uses and does not declare, and
// for one declared as another kind than its use needs.
const (
	undeclaredScalar = 137   // Must declare the scalar variable "@a".
	undeclaredTable  = 1087  // Must declare the table variable "@a".
	notCursor        = 16948 // The variable '@a' is not a cursor variable, ...
)

// sqlServerLookups are errors SQL Server's parser raises for a name it looks
// up in what the server holds: a database (911), an object (208, 15151), a
// collation (448), a language (7678), and a component the server has not
// installed, as full-text search (7609). They are about the server and not
// about the text, so the statement counts as read, as a MySQL-family server's
// lookups do (PrepareVerdict).
var sqlServerLookups = map[int32]bool{
	208:   true,
	448:   true,
	911:   true,
	7609:  true,
	7678:  true,
	15151: true,
}

// variableName finds the variable an error names: "@a" or '@a'.
var variableName = regexp.MustCompile(`["'](@[^"']+)["']`)

// maxDeclared caps the variables ParseOnlyVerdict declares for one
// statement. Each re-ask names one more, so a statement that names more is
// not one a generated row holds.
const maxDeclared = 16

// ParseOnlyVerdict reads SQL Server's answer to a statement sent as a batch
// of its own on a session under SET PARSEONLY ON: whether its parser accepts
// the statement. err is what the batch returned. read returns the errors the
// server reported, in order, and false for an error the server did not send
// (a dropped connection, a timeout). declared sends the statement again
// after the given declarations, which end in a newline, and returns what
// that batch returned.
//
// SET PARSEONLY ON makes the server parse a batch and neither compile nor
// run it, so every error it reports is the parser's. That is more than a
// grammar: the parser also refuses a BREAK outside a WHILE (135), a column
// with no type (173), an option it does not know (155, 195). Each is a
// rejection. Two kinds of error are not about the text:
//
//   - A lookup of what the server holds (sqlServerLookups) is an
//     acceptance.
//   - A variable the statement does not declare (137, 1087, and 16948 for
//     a cursor) is the batch's to declare, and a row is one statement. The
//     verdict asks again with the variable declared before the statement, of
//     the kind the error names: INT, TABLE (a INT) or CURSOR. The parser
//     stops at the first such variable, so each re-ask can name one more.
//     A statement that must start its batch, as CREATE PROCEDURE must (111),
//     cannot be declared for, and keeps its rejection.
//
// The first error decides: the parser stops reading a statement at a
// rejection, and at an undeclared variable it reads on only to the next
// statement's. An answer that is no verdict comes back as an error, so the
// caller fails rather than records agreement: an error the server did not
// send, and an error with no error in it.
func ParseOnlyVerdict(err error, read func(error) ([]SQLServerError, bool), declared func(declarations string) error) (string, error) {
	var names []string // declared, in order
	kinds := map[string]string{}
	for {
		if err == nil {
			return Accept, nil
		}
		errs, ok := read(err)
		if !ok {
			return "", err
		}
		if len(errs) == 0 {
			return "", fmt.Errorf("SQL Server reported an error with no error in it: %w", err)
		}
		first := errs[0]
		if sqlServerLookups[first.Number] {
			return Accept, nil
		}
		var kind string
		switch first.Number {
		case undeclaredScalar:
			kind = "INT"
		case undeclaredTable:
			kind = "TABLE (a INT)"
		case notCursor:
			kind = "CURSOR"
		default:
			return Reject, nil
		}
		m := variableName.FindStringSubmatch(first.Message)
		if m == nil {
			return "", fmt.Errorf("SQL Server's error %d names no variable: %q", first.Number, first.Message)
		}
		name := m[1]
		if prev, ok := kinds[name]; ok && (prev == kind || first.Number != notCursor) {
			// Declared, and still not the variable the statement needs:
			// the statement uses it as two kinds, or must start its batch.
			return Reject, nil
		}
		if len(names) == maxDeclared {
			return "", fmt.Errorf("SQL Server still names an undeclared variable after %d declarations: %w", maxDeclared, err)
		}
		if _, ok := kinds[name]; !ok {
			names = append(names, name)
		}
		kinds[name] = kind
		var b strings.Builder
		for _, n := range names {
			fmt.Fprintf(&b, "DECLARE %s %s;\n", n, kinds[n])
		}
		err = declared(b.String())
	}
}
