package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexNoDuplicateColumn reports the keys top-level statements define that
// list a column twice: the table-level PRIMARY KEY, UNIQUE, INDEX and FOREIGN
// KEY definitions of a CREATE TABLE, those an ALTER TABLE adds, and the index
// a CREATE INDEX of any kind creates, FULLTEXT, SPATIAL and VECTOR included.
// A FULLTEXT, SPATIAL or VECTOR definition, a key inside ALTER TABLE's
// parenthesized ADD form, and statements nested in a routine body are not
// judged.
//
// A key's columns are its key parts that are plain columns, a parenthesized
// one included, compared as written, so `a` and `A` differ and an expression
// key part is never a duplicate. A key gets one finding, for the first
// column that repeats an earlier one, held in Other.
//
// Name is the key's name: the CONSTRAINT symbol, or else the name after
// INDEX, KEY, UNIQUE or FOREIGN KEY, or else PRIMARY for a primary key, which
// the AST keeps no other name for; it is empty for any other key with no
// name. Table is the table's name without its database qualifier, and Value
// the kind of key, KeyIndex for every CREATE INDEX, UNIQUE included.
//
// A CREATE TABLE finding anchors the key definition. Every other finding has
// no range: V1 anchored it at the first character of the statement's text
// that is not white space, which a leading comment moves off the statement.
// Findings are in statement order. Ported from Bytebase's V1 rule
// INDEX_NO_DUPLICATE_COLUMN.
func IndexNoDuplicateColumn(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	forEachKey(stmts, func(key keyDefinition) {
		var columns []string
		switch {
		case key.index != nil:
			columns = keyPartColumns(key.index.Columns)
		case key.kind == v1review.KeyForeign:
			columns = key.constraint.Columns
		case key.kind != "":
			columns = keyPartColumns(key.constraint.IndexColumns)
		default:
		}
		seen := make(map[string]bool, len(columns))
		for _, column := range columns {
			if seen[column] {
				findings = append(findings, v1review.Finding{Code: v1review.IndexColumnDuplicated, Name: key.name, Table: key.table, Other: column, Value: key.kind, Start: key.start, End: key.end})
				return
			}
			seen[column] = true
		}
	})
	return findings
}
