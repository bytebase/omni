package v1review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingIdentifierNoKeyword judges the names statements introduce against
// the V1 keyword table (v1ReservedKeywords): a column definition's, a named
// constraint's, and the name a CREATE SCHEMA, CREATE DATABASE, CREATE INDEX
// or plain CREATE TABLE gives. Each statement is walked whole, so these are
// judged wherever the AST holds them:
//
//   - inside a routine or trigger body, an IF, WHILE, BEGIN or TRY block,
//     and a CREATE SCHEMA's own statements;
//   - in the column list of a table variable, a table type, a function's
//     returned table, and a CREATE EXTERNAL TABLE;
//   - in an ALTER TABLE: the columns of ADD and ALTER COLUMN, and the
//     constraint named by ADD CONSTRAINT and by DROP CONSTRAINT.
//
// Not judged: an inline INDEX of a CREATE TABLE; an XML, spatial or
// full-text index; a view, routine, trigger, type, sequence, synonym or
// variable name; the table of a CREATE TABLE AS SELECT, a CREATE EXTERNAL
// TABLE or a SELECT INTO; a rename through sp_rename; and the column list
// of an OPENJSON or an EXEC ... WITH RESULT SETS.
//
// A name is judged without its qualifier and without its brackets or
// quotes, lower-cased and then upper-cased for the lookup, so a bracketed
// name matches whatever its case though SQL Server accepts it as an
// ordinary identifier. Both foldings are Unicode's, so a name spelled with
// a letter that folds into ASCII, such as the Kelvin sign in KEY, matches
// too. Finding.Name is the lower-cased name.
//
// The V1 advisor placed each finding by searching the statement text for
// the name in its source spelling, which Finding.Search carries: the line
// is that of the first place those bytes occur, be it a comment, a string
// or part of a longer name. The range is the top-level statement's.
// Findings are in the order ast.Inspect visits the nodes. NameIsKeyword is
// the only code. Ported from Bytebase's V1 rule
// NAMING_IDENTIFIER_NO_KEYWORD.
func NamingIdentifierNoKeyword(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		start, end := span(ast.NodeLoc(stmt))
		ast.Inspect(stmt, func(n ast.Node) bool {
			var name string
			switch v := n.(type) {
			case *ast.ColumnDef:
				name = v.Name
			case *ast.ConstraintDef:
				name = v.Name
			case *ast.CreateSchemaStmt:
				name = v.Name
			case *ast.CreateDatabaseStmt:
				name = v.Name
			case *ast.CreateIndexStmt:
				name = v.Name
			case *ast.CreateTableStmt:
				if v.Name != nil {
					name = v.Name.Object
				}
			default:
				return true
			}
			if name == "" {
				return true
			}
			if lower := strings.ToLower(name); isV1ReservedKeyword(lower) {
				findings = append(findings, v1review.Finding{Code: v1review.NameIsKeyword, Name: lower, Search: name, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
