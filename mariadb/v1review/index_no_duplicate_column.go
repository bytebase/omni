package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexNoDuplicateColumn reports the keys top-level statements define that
// list a column twice: the table-level PRIMARY KEY, UNIQUE, INDEX and FOREIGN
// KEY definitions of a CREATE TABLE, those an ALTER TABLE adds, and the index
// a CREATE INDEX creates unless it is FULLTEXT or SPATIAL, so CREATE VECTOR
// INDEX is judged as an INDEX. A FULLTEXT, SPATIAL or VECTOR definition, a
// key inside ALTER TABLE's parenthesized ADD form, and statements nested in a
// routine body are not judged.
//
// A definition's columns are its key parts, compared as written, so `a` and
// `A` differ and a prefix length is ignored. A key gets one finding, for the
// first column that repeats an earlier one, held in Other.
//
// Name is the key's name as parsed, the empty name when it has none, and
// Table the table's name without its database qualifier. Value is the kind
// of key: a CREATE INDEX is KeyIndex or KeyUniqueIndex, a UNIQUE definition
// KeyUnique. The AST does not keep a name written after PRIMARY KEY, so that
// key's Name is empty unless CONSTRAINT names it; the V1 advisor then read a
// name out of the statement text, which the caller keeps.
//
// A CREATE TABLE finding anchors the key definition; every other finding
// anchors the whole statement. Findings are in statement order. Ported from
// Bytebase's V1 rule INDEX_NO_DUPLICATE_COLUMN.
func IndexNoDuplicateColumn(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	forEachKey(stmts, func(key keyDefinition) {
		if key.kind == "" {
			return
		}
		seen := make(map[string]bool, len(key.columns))
		for _, column := range key.columns {
			if seen[column] {
				start, end := span(key.loc)
				findings = append(findings, v1review.Finding{Code: v1review.IndexColumnDuplicated, Name: key.name, Table: key.table, Other: column, Value: key.kind, Start: start, End: end})
				return
			}
			seen[column] = true
		}
	})
	return findings
}
