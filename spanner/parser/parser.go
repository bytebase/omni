// Package parser provides Cloud Spanner SQL parsing for both of Cloud
// Spanner's SQL dialects:
//   - GoogleSQL via Parse / ParseGoogleSQL / ParseBestEffort, delegating to
//     github.com/bytebase/omni/googlesql/parser with DialectSpanner (backed by
//     testdata/compat.tsv).
//   - PostgreSQL via ParsePostgres, delegating to
//     github.com/bytebase/omni/pg/parser (backed by testdata/pg/compat.tsv).
//
// BigQuery and Cloud Spanner share the GoogleSQL language, while Spanner's
// GoogleSQL dialect rejects BigQuery-only procedural scripting, QUALIFY,
// IS [NOT] DISTINCT FROM, parenthesized outer query aliases, and BigQuery-only
// DDL forms. This package gives Cloud Spanner its own dialect entry points and
// live Spanner-emulator-judged syntax compatibility suites, independent of the
// union mode in googlesql/parser and DialectBigQuery in bigquery/parser.
package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
	googlesqlparser "github.com/bytebase/omni/googlesql/parser"
	pgast "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
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

// ParsePostgres parses Cloud Spanner PostgreSQL-dialect SQL input by delegating
// to github.com/bytebase/omni/pg/parser.Parse.
func ParsePostgres(sql string) (*pgast.List, error) {
	return pgparser.Parse(sql)
}
