package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

// NamingIdentifierNoKeyword judges the names top-level statements introduce
// against the parser's keyword list (parser.IsKeyword: reserved and
// non-reserved words alike, matched in any case):
//
//   - CREATE TABLE, in every form: the table, then each written column,
//     then each table-level constraint and index;
//   - ALTER TABLE: the columns of ADD COLUMN in both its forms, the new
//     name of RENAME TO, RENAME COLUMN, CHANGE COLUMN and RENAME INDEX, and
//     a constraint or index added with ADD;
//   - CREATE INDEX: the index.
//
// A constraint or index is judged by the name the AST keeps for it and, for
// a FOREIGN KEY, also by the name written after FOREIGN KEY. So a UNIQUE
// given both a CONSTRAINT symbol and an index name is judged by the index
// name alone, and the name written after PRIMARY KEY is not judged.
//
// Not judged: MODIFY COLUMN, a column's own CONSTRAINT symbol, a constraint
// inside ALTER TABLE's parenthesized ADD form, RENAME TABLE, CREATE VIEW,
// CREATE DATABASE, aliases, and statements nested in a routine body. A name
// is judged as parsed without its database qualifier, less any backticks it
// begins or ends with, and that is Finding.Name, in its source case.
//
// The V1 advisor placed a CREATE TABLE column's finding by searching the
// statement text for the column name as parsed, which Finding.Search
// carries: the line is that of the first place those bytes occur, be it a
// comment, a string or part of a longer name. Such a finding's range is the
// statement's. Every other finding has no Search: a CREATE TABLE constraint
// anchors the constraint, and the table name and every ALTER TABLE and
// CREATE INDEX finding anchor the whole statement, whichever line the name
// is on. Findings are in statement order. NameIsKeyword is the only code.
// Ported from Bytebase's V1 rule NAMING_IDENTIFIER_NO_KEYWORD.
func NamingIdentifierNoKeyword(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	// check judges name; search is set for a finding placed by text search.
	check := func(loc ast.Loc, name, search string) {
		name = strings.Trim(name, "`")
		if !parser.IsKeyword(name) {
			return
		}
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: v1review.NameIsKeyword, Name: name, Search: search, Start: start, End: end})
	}
	checkConstraint := func(loc ast.Loc, c *ast.Constraint) {
		for _, name := range []string{c.Name, c.IndexName} {
			if name != "" {
				check(loc, name, "")
			}
		}
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table != nil {
				check(n.Loc, n.Table.Name, "")
			}
			for _, col := range n.Columns {
				if col != nil {
					check(n.Loc, col.Name, col.Name)
				}
			}
			for _, constraint := range n.Constraints {
				if constraint != nil {
					checkConstraint(constraint.Loc, constraint)
				}
			}
		case *ast.AlterTableStmt:
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATAddColumn:
					for _, col := range commandColumns(cmd) {
						check(n.Loc, col.Name, "")
					}
				case ast.ATRenameTable:
					name := cmd.NewName
					if cmd.NewTable != nil {
						name = cmd.NewTable.Name
					}
					if name != "" {
						check(n.Loc, name, "")
					}
				case ast.ATRenameColumn, ast.ATRenameIndex:
					if cmd.NewName != "" {
						check(n.Loc, cmd.NewName, "")
					}
				case ast.ATChangeColumn:
					if cmd.Column != nil {
						check(n.Loc, cmd.Column.Name, "")
					}
				case ast.ATAddConstraint, ast.ATAddIndex:
					if cmd.Constraint != nil {
						checkConstraint(n.Loc, cmd.Constraint)
					}
				default:
				}
			}
		case *ast.CreateIndexStmt:
			if n.IndexName != "" {
				check(n.Loc, n.IndexName, "")
			}
		}
	}
	return findings
}
