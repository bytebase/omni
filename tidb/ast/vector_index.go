package ast

import "strings"

// indexTypeHNSW is the only index type a vector index has.
const indexTypeHNSW = "HNSW"

// SyntaxRefused reports whether TiDB's grammar action refuses this CREATE INDEX
// after the grammar itself accepted the text.
//
// CreateIndexStmt's action in v8.5.5 pkg/parser/parser.y answers ErrSyntax, which
// is error 1149 and not the 1064 of a yacc failure, for a vector index that is not
// exactly one functional key part of type HNSW, and for HNSW on any other kind of
// index. The server answers it before it looks at any table: `CREATE VECTOR INDEX i
// ON nosuch (a)` is 1149 where `CREATE VECTOR INDEX i ON nosuch ((a))` is 1146.
//
// The parser reads these as parsed, because its corpus takes only 1064 as a
// rejection, so the review and the catalog ask here instead.
func (s *CreateIndexStmt) SyntaxRefused() bool {
	// The action reads the option list's type first and the one before ON only
	// when the list has none, so a type among the options wins: USING BTREE ON
	// t ((a)) USING HNSW runs, and USING HNSW ON t ((a)) USING BTREE is 1149.
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

// SyntaxRefused reports whether TiDB's grammar action refuses this table
// constraint after the grammar itself accepted the text: a VECTOR INDEX that is
// not exactly one functional key part of type HNSW, which ConstraintVectorIndex in
// v8.5.5 pkg/parser/parser.y answers with ErrSyntax, error 1149.
//
// Every other kind of constraint takes what the grammar takes, HNSW included:
// `ADD INDEX (a) USING HNSW` and `CREATE TABLE t (a INT, INDEX (a) USING HNSW)` run.
func (c *Constraint) SyntaxRefused() bool {
	if c.Type != ConstrVector {
		return false
	}
	// This action reads the other way round from CreateIndexStmt's: the type
	// before the key parts overwrites the option list's.
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
func lastUsingOption(opts []*IndexOption) string {
	indexType := ""
	for _, opt := range opts {
		if strings.EqualFold(opt.Name, "USING") {
			if s, ok := opt.Value.(*StringLit); ok {
				indexType = strings.ToUpper(s.Value)
			}
		}
	}
	return indexType
}

func oneFunctionalKeyPart(parts []*IndexColumn) bool {
	return len(parts) == 1 && parts[0].Functional
}
