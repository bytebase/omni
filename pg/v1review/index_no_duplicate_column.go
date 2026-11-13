package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexNoDuplicateColumn reports the keys top-level statements define that
// list a column twice: the index a CREATE INDEX creates, and the table-level
// PRIMARY KEY, UNIQUE and FOREIGN KEY constraints a CREATE TABLE defines or
// an ALTER TABLE ... ADD adds. A key a column definition declares, an
// EXCLUDE constraint, and a statement in a CREATE SCHEMA or in a routine
// body are not judged.
//
// A CREATE INDEX's columns are its plain column elements, so an expression
// element is never a duplicate; a constraint's are the columns it lists, a
// FOREIGN KEY's its referencing columns. Columns are compared as the parser
// stores them, unquoted ones folded to lower case. A key gets one finding,
// for the first column that repeats an earlier one, held in Other.
//
// Name is the index's or the constraint's name as the parser stores it, the
// empty name when it has none, and Table the table's name without its
// schema. Value is the kind of key: a CREATE INDEX is KeyIndex, unique or
// not, and a UNIQUE constraint KeyUnique.
//
// Every finding anchors the whole statement. Bytebase's V1 advisor reported
// a CREATE INDEX on the statement text's first non-blank line, and a
// constraint on the line of the first occurrence in the statement text of
// "PRIMARY KEY", "UNIQUE" or "FOREIGN KEY", by its kind, matched case
// sensitively and so not necessarily the constraint's own: that text is
// Search. Findings are in statement order. Ported from Bytebase's V1 rule
// INDEX_NO_DUPLICATE_COLUMN.
func IndexNoDuplicateColumn(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	forEachKey(stmts, func(key keyDefinition) {
		seen := make(map[string]bool, len(key.columns))
		for _, column := range key.columns {
			if seen[column] {
				start, end := span(key.loc)
				findings = append(findings, v1review.Finding{Code: v1review.IndexColumnDuplicated, Name: key.name, Table: key.table, Other: column, Value: key.kind, Search: keyKeyword[key.kind], Start: start, End: end})
				return
			}
			seen[column] = true
		}
	})
	return findings
}

// keyKeyword is the text the V1 advisor searched a statement for to place a
// constraint's finding, by the constraint's kind.
var keyKeyword = map[string]string{
	v1review.KeyPrimary: "PRIMARY KEY",
	v1review.KeyUnique:  "UNIQUE",
	v1review.KeyForeign: "FOREIGN KEY",
}
