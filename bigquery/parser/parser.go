// Package parser provides BigQuery SQL parsing by delegating to
// github.com/bytebase/omni/googlesql/parser.
//
// BigQuery and Cloud Spanner share the GoogleSQL language, and Bytebase
// dispatches BIGQUERY directly to googlesql/parser. This package exists so
// BigQuery has its own live-BigQuery-judged syntax compatibility suite
// (testdata/compat.tsv), independent of the Cloud Spanner emulator verdicts in
// googlesql/parser/testdata/compat.tsv.
package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
	googlesqlparser "github.com/bytebase/omni/googlesql/parser"
)

// Parse parses BigQuery SQL input using googlesql/parser with DialectBigQuery.
func Parse(sql string) (*ast.File, error) {
	return googlesqlparser.ParseWithDialect(sql, googlesqlparser.DialectBigQuery)
}
