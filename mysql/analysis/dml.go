package analysis

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
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
