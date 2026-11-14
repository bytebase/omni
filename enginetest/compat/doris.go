package compat

import (
	"regexp"
	"strings"
)

// DorisVerdict reads a Doris server's answer to a statement: whether Doris
// takes its syntax. err is what sending the statement returned, and message
// returns the server's error message, or false for an error the server did
// not send (a dropped connection, a timeout).
//
// The statement is to be sent inside a transaction (BEGIN), where Doris
// parses it and then refuses to run anything but INSERT, UPDATE, DELETE and
// the transaction's own statements: "This is in a transaction, only insert,
// update, delete, commit, rollback is acceptable". That refusal comes after
// both of Doris's parsers have read the statement, so it is an acceptance:
// the ANTLR parser (DorisParser.g4), and the legacy CUP parser
// (sql_parser.cup) for the statements the ANTLR parser hands back. The
// refusal comes from the first as "errCode = 2, detailMessage = This is
// ..." and from the second as "TException, msg: ... This is ...".
//
// A rejection is the parser's own exception. The ANTLR parser's
// ParseException says where it stopped, "\n<message>(line 1, pos 7)\n",
// for a syntax error ("mismatched input", "no viable alternative") and for
// what its AST builder refuses ("Every derived table must have its own
// alias"). The CUP parser's says "Syntax error in line 1:", or "sql
// parsing error, please check your sql" where it failed outright. Every
// other error comes after the parse and is an acceptance: a lookup or check
// of the AST builder, and the execution of a DML statement a transaction
// runs. So is a statement that ran.
//
// An error the server did not send is no verdict, and comes back as an
// error, so the caller fails rather than records agreement.
func DorisVerdict(err error, message func(error) (string, bool)) (string, error) {
	if err == nil {
		return Accept, nil
	}
	msg, ok := message(err)
	if !ok {
		return "", err
	}
	// The parsers' exceptions first: the ANTLR parser's echoes the
	// statement after "== SQL ==", which may quote the refusal below.
	switch {
	case dorisParseException.MatchString(msg),
		strings.Contains(msg, "Syntax error in line "),
		strings.Contains(msg, "sql parsing error, please check your sql"):
		return Reject, nil
	case strings.Contains(msg, dorisInTransaction):
		return Accept, nil
	}
	return Accept, nil
}

// dorisInTransaction is how Doris refuses to run a parsed statement inside
// a transaction.
const dorisInTransaction = "This is in a transaction, only insert, update, delete, commit, rollback is acceptable."

// dorisParseException matches the message of the ANTLR parser's
// ParseException: after "detailMessage = ", a newline, the message, and
// where the parser stopped.
var dorisParseException = regexp.MustCompile(`detailMessage = \n(?s:.*)\(line \d+, pos \d+\)\n`)
