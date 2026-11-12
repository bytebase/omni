package compat

import "fmt"

// ProbeTail is what PrepareVerdict's probe appends to a statement. It starts
// on a new line so a trailing comment in the statement cannot swallow it.
const ProbeTail = "\n; SELECT FROM"

// ProbeTailInStatement is the probe's tail on a server that refuses to
// prepare a text of two statements before it parses either, as OceanBase
// does (1235): ProbeTail would get that one answer for every statement.
// This tail is no statement of its own, and a syntax error wherever it
// stands. It can miss a refusal that an action makes on the statement's
// last reduction, when the parser reads the tail before it reduces.
const ProbeTailInStatement = "\n @@@ ((("

// lookupErrorCodes are errors a MySQL-family server raises while parsing,
// for a name it looks up in what the server has: a system variable (1193),
// a character set (1115) and a collation (1273, 1253). They are about that
// list and not about the shape of the text, so the statement counts as
// read: `SELECT @@nosuch` answers 1193 while parsing, and omni, which keeps
// no such list, reads it.
var lookupErrorCodes = map[uint16]bool{
	1115: true,
	1193: true,
	1253: true,
	1273: true,
}

// PrepareVerdict reads a MySQL-family server's answer to PREPARE of a
// statement: whether the server accepts its syntax. err is what PREPARE
// returned. code returns the server's error number, and false for an error
// the server did not send (a dropped connection, a refused local file).
// probe prepares the statement with ProbeTail (or ProbeTailInStatement)
// after it, and is called only when err leaves the verdict open.
//
// 1064 is the parser's rejection, and not the only one: a grammar action
// refuses text its production read, with a code it chooses, while the parse
// is still running. Other errors come after the parse (1146 for a missing
// table, 1295 for a statement PREPARE does not carry), and a code alone
// cannot tell the two apart: 1105 is both. The server reads all of the text
// before it looks at anything, so a statement refused after parsing answers
// the tail's 1064 to the probe, and one an action refused answers its own
// code again, because the action ended the parse before the tail.
//
// An error the server did not send is no verdict, and comes back as an
// error so the caller recovers or fails rather than records agreement.
func PrepareVerdict(err error, code func(error) (uint16, bool), probe func() error) (bool, error) {
	if err == nil {
		return true, nil
	}
	n, ok := code(err)
	if !ok {
		return false, err
	}
	if n == 1064 {
		return false, nil
	}
	if lookupErrorCodes[n] {
		return true, nil
	}
	tailErr := probe()
	if tailErr == nil {
		return false, fmt.Errorf("the server took the probe's syntax error after answering %w", err)
	}
	tailed, ok := code(tailErr)
	if !ok {
		return false, tailErr
	}
	return tailed != n, nil
}
