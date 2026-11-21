// Package parser parses CockroachDB SQL statements using CockroachDB's own
// parser (github.com/cockroachdb/cockroachdb-parser).
package parser

import (
	crrawparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser/statements"

	// Registers parser.ParseDoBlockFn; without it the SQL parser panics on a DO block.
	_ "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
)

// Parse parses one or more CockroachDB SQL statements into statements carrying
// the AST and the SQL text CockroachDB's parser read.
func Parse(sql string) (statements.Statements, error) {
	return crrawparser.Parse(sql)
}
