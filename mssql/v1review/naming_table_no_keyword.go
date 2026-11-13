package v1review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTableNoKeyword judges the table a top-level plain CREATE TABLE
// creates against the V1 keyword table (v1ReservedKeywords). Statements
// nested in a routine body or an IF or BEGIN block are not judged, nor are
// CREATE TABLE AS SELECT, CREATE TABLE AS CLONE OF, CREATE EXTERNAL TABLE,
// SELECT INTO, and a rename through sp_rename.
//
// The name is judged without its database or schema qualifier and without
// its brackets or quotes, upper-cased for the lookup, so a bracketed name
// matches whatever its case though SQL Server accepts it as an ordinary
// identifier. The upper-casing is Unicode's, so a name spelled with a
// letter that upper-cases into ASCII, such as the dotless ı of ınsert,
// matches too. A temporary table's # prefix is part of the name, so it
// never matches. Finding.Name is the name before that upper-casing.
//
// Every finding anchors the whole statement. NameIsKeyword is the only
// code. Ported from Bytebase's V1 rule NAMING_TABLE_NO_KEYWORD.
func NamingTableNoKeyword(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.CreateTableStmt)
		if !ok || n.Name == nil {
			continue
		}
		name := n.Name.Object
		if name == "" || !isV1ReservedKeyword(name) {
			continue
		}
		start, end := span(n.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.NameIsKeyword, Name: name, Start: start, End: end})
	}
	return findings
}
