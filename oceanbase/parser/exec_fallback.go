package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// OceanBase reads a statement that does not parse with its executable
// comments' bodies in it a second time, with every executable comment read as
// a plain comment, and runs it if that parses. ObParser::parse_sql retries a
// failed parse with enable_compatible_comment_ off, and the lexer then opens
// an ordinary comment at "/*!" (BEGIN(in_c_comment)), which ends at the first
// "*/". The retry reads the whole statement again, so every executable comment
// in it drops, not only the one that failed: "SELECT 1 /*! + 1 */ /*! FROM */"
// returns 1. A statement that parses but leaves an executable comment open is
// refused without a retry (the check after it), so "SELECT 1 /*! + 'a*/ -- '"
// is 1064 while "SELECT 1 /*! FROM 'a*/ , 2 -- '" returns two columns. A
// statement left with no token is an empty query, which OceanBase answers with
// no error: "/*! SELECT 1 FROM */" runs and returns nothing.
//
// The retry belongs to the SQL parser (retriesAsSQL, sqlParserReads). A
// stored-program or event statement goes to the PL parser, which has no retry,
// so "CALL p() /*! FROM */" and "CREATE PROCEDURE p() /*! FROM */ BEGIN SELECT
// 1; END" are 1064.
//
// This is not tolerance in Parse: it reads what the engine reads. Measured
// against OceanBase CE 4.3.5.6 (oceanbase/oceanbase-ce:4.3.5-lts) on 2026-11-05,
// through PREPARE and as a query, which agree (execFallbackCases).
//
// Ref: OceanBase CE v4.3.5_CE_BP6, the tag that builds 4.3.5.6
// (CMakeLists.txt:6, OB_BUILD_VERSION): src/sql/parser/ob_parser.cpp:966-1052 (ObParser::parse_sql, the retry at
// :973-981 and the open-comment check at :1006-1010), :1172-1219 (the routing
// in ObParser::parse), :39-122 and :342-522 (ObParser::is_pl_stmt,
// transform_normal); src/sql/parser/sql_parser_mysql_mode.l:1412-1421 and
// :1844-1854. v4.3.5_CE_BP1 has the same retry but routes no EVENT statement
// to the PL parser, which 4.3.5.6 does.

// retriesAsSQL reports whether OceanBase gives stmt, a statement's text from
// its first byte, the SQL parser's retry. ObParser::parse sends a statement to
// the SQL parser unless ObParser::is_pl_stmt finds a stored-program word at
// its start; is_pl_stmt reads words separated by whitespace, after skipping
// comments, executable ones included, and decides on the first word or two.
// Two kinds of statement go to the PL parser first and to the SQL parser when
// it fails: one that starts with "/*!", and CREATE FUNCTION, which the SQL
// parser reads as a loadable function.
func retriesAsSQL(stmt string) bool {
	i := 0
	for i < len(stmt) && (stmt[i] == ' ' || stmt[i] == '\t' || stmt[i] == '\r' || stmt[i] == '\f' || stmt[i] == '\n') {
		i++
	}
	if strings.HasPrefix(stmt[i:], "/*!") {
		return true
	}
	state := ""
	for i < len(stmt) {
		switch c := stmt[i]; {
		case isPLSpace(c):
			i++
		case c == '#' || c == '-' && i+1 < len(stmt) && stmt[i+1] == '-':
			for i < len(stmt) && stmt[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < len(stmt) && stmt[i+1] == '*':
			end := strings.Index(stmt[i+2:], "*/")
			if end < 0 {
				return true
			}
			i += 2 + end + 2
		default:
			j := i
			for j < len(stmt) && !isPLSpace(stmt[j]) {
				j++
			}
			word := strings.ToLower(stmt[i:j])
			i = j
			next, sql, decided := plWord(state, word)
			if decided {
				return sql
			}
			state = next
		}
	}
	return true
}

// plWord is ObParser::transform_normal: given the state the words before it
// left and the next word, it returns the next state, or, once that is decided,
// whether the statement gets the SQL parser's retry.
func plWord(state, word string) (next string, sql, decided bool) {
	switch state {
	case "":
		switch {
		case word == "do", word == "declare", word == "procedure", word == "function",
			word == "package", word == "trigger", word == "event", word == "type",
			word == "signal", word == "resignal", word == "call", strings.HasPrefix(word, "<<"):
			return "", false, true
		case word == "create", word == "dd", word == "begin", word == "drop",
			word == "alter", word == "update", word == "submit", word == "cancel":
			return word, false, false
		}
		return "", true, true
	case "create":
		switch {
		case word == "function":
			// The PL parser first, then the SQL parser for a single statement.
			return "", true, true
		case word == "procedure", word == "package", word == "trigger", word == "event",
			word == "type", strings.HasPrefix(word, "definer"):
			return "", false, true
		case word == "or", word == "replace", word == "editionable", word == "noneditionable", word == "force":
			return state, false, false
		}
		return "", true, true
	case "dd":
		return "", word != "begin", true
	case "begin":
		return "", strings.HasPrefix(word, ";") || strings.HasPrefix(word, "work"), true
	case "drop", "alter":
		switch {
		case word == "procedure", word == "function", word == "package", word == "trigger",
			word == "event", word == "type", strings.HasPrefix(word, "definer"):
			return "", false, true
		}
		return "", true, true
	case "update":
		return "", word != "of", true
	default: // "submit", "cancel"
		return "", word != "job", true
	}
}

// isPLSpace is ISSPACE in ob_parser.cpp.
func isPLSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// sqlParserReads reports whether OceanBase's SQL parser, which runs the retry,
// has the statement stmt is. Its grammar has no stored-program or event
// statement: the PL parser alone reads DO, CALL, SIGNAL, RESIGNAL, a stored
// routine, a trigger and an event, so a retry that comes to one of them is
// 1064, as "/*! x */ CALL p()" and "DO(1) /*! FROM */" are. DROP FUNCTION and
// CREATE FUNCTION ... SONAME are the loadable-function statements it does
// have.
//
// Ref: OceanBase CE v4.3.5_CE_BP6 src/sql/parser/sql_parser_mysql_mode.y:624
// (stmt), :5866 (create_function_stmt), :5874 (drop_function_stmt).
func sqlParserReads(stmt nodes.Node) bool {
	switch s := stmt.(type) {
	case *nodes.DoStmt, *nodes.CallStmt, *nodes.SignalStmt, *nodes.ResignalStmt,
		*nodes.CreateTriggerStmt, *nodes.DropTriggerStmt, *nodes.AlterRoutineStmt,
		*nodes.CreateEventStmt, *nodes.AlterEventStmt, *nodes.DropEventStmt:
		return false
	case *nodes.CreateFunctionStmt:
		return s.Soname != ""
	case *nodes.DropRoutineStmt:
		return !s.IsProcedure
	}
	return true
}
