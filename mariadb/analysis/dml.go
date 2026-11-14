package analysis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
)

// dmlKeywordStart returns the offset of the UPDATE or DELETE keyword of the
// statement at loc. omni starts the node at its WITH clause, if any, and keeps
// no CTE on it, so the clause is the text between the two.
func dmlKeywordStart(sql string, loc ast.Loc) int {
	if loc.Start < 0 || loc.End > len(sql) || loc.Start >= loc.End {
		return loc.Start
	}
	tokens := parser.Tokenize(sql[loc.Start:loc.End])
	if len(tokens) == 0 || parser.TokenName(tokens[0].Type) != "WITH" {
		return loc.Start
	}
	depth := 0
	for _, token := range tokens[1:] {
		switch token.Type {
		case '(':
			depth++
		case ')':
			depth--
		default:
			if depth != 0 {
				continue
			}
			if name := parser.TokenName(token.Type); name == "UPDATE" || name == "DELETE" {
				return loc.Start + token.Loc
			}
		}
	}
	return loc.Start
}

// endBeforeReturning clamps end to just before a DML statement's RETURNING
// clause, or returns end unchanged when there is none.
//
// RETURNING is MariaDB's and is the last clause of the statements that take it.
// Every rewrite here that turns a DML into a SELECT by copying its text through
// to the statement end must stop short of it: RETURNING means nothing in a
// SELECT and the server rejects the result. Route a new rewrite through this
// function rather than repeating the search.
//
// An error means the clause is present but its keyword could not be located.
// Callers decline to transform rather than emit SQL that fails when it runs.
func endBeforeReturning(returning []ast.Node, sql string, end int) (int, error) {
	if end <= 0 || end > len(sql) {
		end = len(sql)
	}
	if len(returning) == 0 {
		return end, nil
	}
	var first int
	switch e := returning[0].(type) {
	case *ast.ColumnRef:
		first = e.Loc.Start
	case *ast.ResTarget:
		first = e.Loc.Start
	case *ast.StarExpr:
		first = e.Loc.Start
	default:
		return 0, fmt.Errorf("cannot locate RETURNING from %T", returning[0])
	}
	if first <= 0 || first > end {
		return 0, errors.New("cannot locate RETURNING")
	}
	// The keyword is the token before the clause's first item. A text search will
	// not do: a comment or string literal in or beside the clause may carry the
	// word, as in `RETURNING /* RETURNING */ b1`. The lexer skips comments and
	// gives a quoted RETURNING a token of its own, so neither is mistaken for it.
	tokens := parser.Tokenize(sql)
	for i, tok := range tokens {
		if tok.Loc != first {
			continue
		}
		if i == 0 {
			break
		}
		if prev := tokens[i-1]; strings.EqualFold(prev.Str, "RETURNING") {
			return prev.Loc, nil
		}
		break
	}
	return 0, errors.New("cannot locate the RETURNING keyword")
}
