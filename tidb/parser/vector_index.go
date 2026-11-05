package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// indexTypeHNSW is the only index type a vector index has.
const indexTypeHNSW = "HNSW"

// refuseVectorIndexShape is what TiDB's grammar actions answer for text the
// productions themselves read: ErrSyntax, error 1149, with the wording of a
// syntax error. The server raises it while it parses, so
// `CREATE VECTOR INDEX i ON nosuch (a)` answers 1149 and
// `CREATE VECTOR INDEX i ON nosuch ((a))` answers 1146, and a syntax error later
// in the same text does not replace it.
func refuseVectorIndexShape(start int) error {
	return &ParseError{
		Message:  "a vector index takes one functional key part and the index type HNSW, and HNSW takes a vector index",
		Position: start,
	}
}

// createIndexRefused is CreateIndexStmt's action in v8.5.5 pkg/parser/parser.y: a
// vector index is exactly one functional key part of type HNSW, and no other
// index CREATE INDEX makes is of type HNSW.
//
// The action reads the option list's type first and the one before ON only when
// the list has none, so a type among the options wins: USING BTREE ON t ((a))
// USING HNSW runs, and USING HNSW ON t ((a)) USING BTREE is 1149.
func createIndexRefused(s *nodes.CreateIndexStmt) bool {
	indexType := lastUsingOption(s.Options)
	if indexType == "" {
		indexType = strings.ToUpper(s.IndexType)
	}
	if !s.Vector {
		return indexType == indexTypeHNSW
	}
	if indexType != "" && indexType != indexTypeHNSW {
		return true
	}
	return !oneFunctionalKeyPart(s.Columns)
}

// vectorConstraintRefused is ConstraintVectorIndex's action in v8.5.5
// pkg/parser/parser.y. It reads the type the other way round from
// CreateIndexStmt's: the one before the key parts overwrites the option list's.
//
// Every other kind of constraint takes what the grammar takes, HNSW included:
// `ADD INDEX (a) USING HNSW` and `CREATE TABLE t (a INT, INDEX (a) USING HNSW)` run.
func vectorConstraintRefused(c *nodes.Constraint) bool {
	indexType := strings.ToUpper(c.IndexType)
	if indexType == "" {
		indexType = lastUsingOption(c.IndexOptions)
	}
	if indexType != "" && indexType != indexTypeHNSW {
		return true
	}
	return !oneFunctionalKeyPart(c.IndexColumns)
}

// lastUsingOption returns the index type an option list settles on, in upper case.
func lastUsingOption(opts []*nodes.IndexOption) string {
	indexType := ""
	for _, opt := range opts {
		if strings.EqualFold(opt.Name, "USING") {
			if s, ok := opt.Value.(*nodes.StringLit); ok {
				indexType = strings.ToUpper(s.Value)
			}
		}
	}
	return indexType
}

func oneFunctionalKeyPart(parts []*nodes.IndexColumn) bool {
	return len(parts) == 1 && parts[0].Functional
}
