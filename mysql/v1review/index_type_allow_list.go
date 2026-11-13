package v1review

import (
	"slices"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexTypeAllowList reports the indexes top-level statements define whose
// type is not in allow: the table-level PRIMARY KEY, UNIQUE, INDEX, FULLTEXT
// and SPATIAL definitions of a CREATE TABLE, those an ALTER TABLE adds, and
// the index a CREATE INDEX creates. A FOREIGN KEY, a key a column definition
// declares, a key inside ALTER TABLE's parenthesized ADD form, and
// statements nested in a routine body are not judged.
//
// An index's type is the one a USING clause before its key parts names, as
// parsed. A USING clause after the key parts is an index option and is not
// read, so `KEY k (a) USING HASH` is judged by the default: FULLTEXT or
// SPATIAL for those indexes, BTREE for every other, whatever the table's
// storage engine. The type is compared with allow's
// entries as they are, so an entry in another letter case allows nothing.
//
// Type holds the index's type; Name and Table are empty. A CREATE TABLE
// finding anchors the index definition; every other finding anchors the
// whole statement. Findings are in statement order. Ported from Bytebase's
// V1 rule INDEX_TYPE_ALLOW_LIST.
func IndexTypeAllowList(stmts []ast.Node, allow []string) []v1review.Finding {
	var findings []v1review.Finding
	forEachKey(stmts, func(key keyDefinition) {
		using, fulltext, spatial := "", false, false
		if key.index != nil {
			using, fulltext, spatial = key.index.IndexType, key.index.Fulltext, key.index.Spatial
		} else {
			switch key.constraint.Type {
			case ast.ConstrPrimaryKey, ast.ConstrUnique, ast.ConstrIndex:
			case ast.ConstrFulltextIndex:
				fulltext = true
			case ast.ConstrSpatialIndex:
				spatial = true
			default:
				return
			}
			using = key.constraint.IndexType
		}
		indexType := "BTREE"
		switch {
		case using != "":
			indexType = using
		case fulltext:
			indexType = "FULLTEXT"
		case spatial:
			indexType = "SPATIAL"
		}
		if slices.Contains(allow, indexType) {
			return
		}
		start, end := span(key.loc)
		findings = append(findings, v1review.Finding{Code: v1review.IndexTypeNotAllowed, Type: indexType, Start: start, End: end})
	})
	return findings
}
