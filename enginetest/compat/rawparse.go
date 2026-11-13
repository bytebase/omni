package compat

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// RawParseHead is what a raw-parse probe sends before the statement, in one
// simple query. PostgreSQL parses every statement of a simple query before
// it analyzes or runs the first (exec_simple_query calls pg_parse_query on
// the whole string), so the division by zero is reached only when the
// statement after it parsed, and nothing of that statement ever runs. It
// ends in a newline so a statement that starts with a comment or a quote
// cannot join it.
const RawParseHead = "SELECT 1/0;\n"

// divisionByZero is the SQLSTATE of RawParseHead's own failure.
const divisionByZero = "22012"

// RawParseVerdict reads a PostgreSQL server's answer to the simple query
// RawParseHead + statement: whether its raw parser, the lexer and gram.y,
// accepts the statement. err is what the query returned. state returns the
// server's SQLSTATE and the position it reported, in characters from 1 into
// the text sent and 0 for none, and false for an error the server did not
// send (a dropped connection, a timeout).
//
// The head's division by zero means all of the text parsed: accept. Any
// other server error was raised before the head ran, so it is the raw
// parser's, whatever its SQLSTATE: 42601 from the lexer and the grammar, and
// what a grammar action chooses, as 0A000 for a CHECK constraint marked
// DEFERRABLE. The verdict is reject@<byte>, the byte offset into statement
// where the server put the error, and a bare reject where it reported no
// position, as that action does.
//
// An answer that is no verdict comes back as an error, so the caller fails
// rather than records agreement: no error at all, an error the server did
// not send, and a position outside the statement.
func RawParseVerdict(statement string, err error, state func(error) (sqlstate string, position int, ok bool)) (string, error) {
	if err == nil {
		return "", fmt.Errorf("the server ran the probe's division by zero without an error")
	}
	code, position, ok := state(err)
	if !ok {
		return "", err
	}
	if code == divisionByZero {
		return Accept, nil
	}
	if position == 0 {
		return Reject, nil
	}
	head := utf8.RuneCountInString(RawParseHead)
	chars := position - 1 - head
	if chars < 0 || chars > utf8.RuneCountInString(statement) {
		return "", fmt.Errorf("the server put its error at character %d of the probe, outside the statement: %w", position, err)
	}
	return RejectAt(byteOffset(statement, chars)), nil
}

// RejectAt is the verdict of a rejection at a byte offset into the
// statement: reject@<byte>.
func RejectAt(offset int) string { return Reject + "@" + strconv.Itoa(offset) }

// byteOffset returns the byte offset of the character at index chars of s.
// PostgreSQL counts positions in characters of the client encoding.
func byteOffset(s string, chars int) int {
	off := 0
	for i := 0; i < chars && off < len(s); i++ {
		_, size := utf8.DecodeRuneInString(s[off:])
		off += size
	}
	return off
}
