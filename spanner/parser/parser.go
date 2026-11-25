// Package parser provides Cloud Spanner GoogleSQL parsing through Parse /
// ParseGoogleSQL / ParseBestEffort, delegating to
// github.com/bytebase/omni/googlesql/parser with DialectSpanner (backed by
// testdata/compat.tsv).
//
// BigQuery and Cloud Spanner share the GoogleSQL language, while Spanner's
// GoogleSQL dialect rejects BigQuery-only procedural scripting, QUALIFY,
// IS [NOT] DISTINCT FROM, parenthesized outer query aliases, and BigQuery-only
// DDL forms. This package gives Cloud Spanner its own dialect entry points and
// a live Spanner-emulator-judged syntax compatibility suite, independent of the
// union mode in googlesql/parser and DialectBigQuery in bigquery/parser.
//
// Spanner's PostgreSQL dialect is not supported.
package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
	googlesqlparser "github.com/bytebase/omni/googlesql/parser"
)

// Parse parses Cloud Spanner GoogleSQL input using googlesql/parser with
// DialectSpanner.
func Parse(sql string) (*ast.File, error) {
	return googlesqlparser.ParseWithDialect(sql, googlesqlparser.DialectSpanner)
}

// ParseGoogleSQL is an explicit alias for Parse (Cloud Spanner's GoogleSQL
// dialect).
func ParseGoogleSQL(sql string) (*ast.File, error) {
	return Parse(sql)
}

// ParseBestEffort parses Cloud Spanner GoogleSQL input in tolerant mode using
// googlesql/parser with DialectSpanner.
func ParseBestEffort(sql string) *googlesqlparser.ParseResult {
	return googlesqlparser.ParseBestEffortWithDialect(sql, googlesqlparser.DialectSpanner)
}
