package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/snowflake/parser"
	"github.com/bytebase/omni/v1review"
)

func TestSchemaBackwardCompatibility(t *testing.T) {
	// compatFinding is a finding as the test compares it: its code, name and
	// statement, and the text its range anchors.
	type compatFinding struct {
		Code      v1review.Code
		Name      string
		Statement int
		Text      string
	}
	f := func(code v1review.Code, name string, stmt int, text string) compatFinding {
		return compatFinding{code, name, stmt, text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		db   string
		want []compatFinding
	}{
		{
			name: "drop table keyed with defaults, unquoted upper-cased",
			sql:  `DROP TABLE foo; DROP TABLE s.Bar; DROP TABLE d.s."Baz"`,
			db:   "MyDb",
			want: []compatFinding{
				f(v1review.IncompatibleDropTable, `"MyDb.PUBLIC.FOO"`, 0, "DROP TABLE foo"),
				f(v1review.IncompatibleDropTable, `"MyDb.S.BAR"`, 1, "DROP TABLE s.Bar"),
				f(v1review.IncompatibleDropTable, `"D.S.Baz"`, 2, `DROP TABLE d.s."Baz"`),
			},
		},
		{
			name: "empty current database kept as an empty part",
			sql:  "DROP TABLE t; DROP SCHEMA s",
			want: []compatFinding{
				f(v1review.IncompatibleDropTable, `".PUBLIC.T"`, 0, "DROP TABLE t"),
				f(v1review.IncompatibleDropSchema, `".S"`, 1, "SCHEMA s"),
			},
		},
		{
			name: "name printed Go-quoted",
			sql:  `DROP TABLE "a""b\c"`,
			db:   "DB",
			want: []compatFinding{f(v1review.IncompatibleDropTable, `"DB.PUBLIC.a\"b\\c"`, 0, `DROP TABLE "a""b\c"`)},
		},
		{
			name: "created without IF NOT EXISTS dropped freely, and never forgotten",
			sql:  "CREATE TABLE t(a INT); DROP TABLE t; DROP TABLE DB.PUBLIC.T",
			db:   "DB",
		},
		{
			name: "created with IF NOT EXISTS reported as maybe new",
			sql:  "CREATE TABLE IF NOT EXISTS t(a INT); DROP TABLE t; CREATE SCHEMA IF NOT EXISTS s; DROP SCHEMA s; CREATE DATABASE IF NOT EXISTS d; DROP DATABASE d",
			db:   "DB",
			want: []compatFinding{
				f(v1review.IncompatibleDropMaybeNewTable, `"DB.PUBLIC.T"`, 1, "DROP TABLE t"),
				f(v1review.IncompatibleDropMaybeNewSchema, `"DB.S"`, 3, "SCHEMA s"),
				f(v1review.IncompatibleDropMaybeNewDatabase, `"D"`, 5, "DATABASE d"),
			},
		},
		{
			name: "the last create of a name wins",
			sql:  "CREATE TABLE t(a INT); CREATE TABLE IF NOT EXISTS t(a INT); DROP TABLE t; CREATE OR REPLACE TABLE t AS SELECT 1 AS a; DROP TABLE t",
			db:   "DB",
			want: []compatFinding{f(v1review.IncompatibleDropMaybeNewTable, `"DB.PUBLIC.T"`, 2, "DROP TABLE t")},
		},
		{
			name: "schema name with its database",
			sql:  "CREATE SCHEMA d1.s; DROP SCHEMA D1.S; DROP SCHEMA d2.s",
			db:   "DB",
			want: []compatFinding{f(v1review.IncompatibleDropSchema, `"D2.S"`, 2, "SCHEMA d2.s")},
		},
		{
			name: "drop database",
			sql:  "DROP DATABASE d",
			db:   "DB",
			want: []compatFinding{f(v1review.IncompatibleDropDatabase, `"D"`, 0, "DATABASE d")},
		},
		{
			name: "drop column of any table, one finding an action",
			sql:  "CREATE TABLE t(a INT, b INT); ALTER TABLE t DROP COLUMN \"b\", a; ALTER TABLE t DROP c",
			db:   "DB",
			want: []compatFinding{
				f(v1review.IncompatibleDropColumn, `"b","A"`, 1, `TABLE t DROP COLUMN "b", a`),
				f(v1review.IncompatibleDropColumn, `"C"`, 2, `TABLE t DROP c`),
			},
		},
		{
			name: "statement on a later line",
			sql:  "SELECT 1;\n\n-- note\nDROP TABLE\n  t",
			db:   "DB",
			want: []compatFinding{f(v1review.IncompatibleDropTable, `"DB.PUBLIC.T"`, 1, "DROP TABLE\n  t")},
		},
		{
			name: "drop schema and database anchored after DROP",
			sql:  "DROP\nSCHEMA s; DROP\nDATABASE d",
			db:   "DB",
			want: []compatFinding{
				f(v1review.IncompatibleDropSchema, `"DB.S"`, 0, "SCHEMA s"),
				f(v1review.IncompatibleDropDatabase, `"D"`, 1, "DATABASE d"),
			},
		},
		{name: "drop view not judged", sql: "DROP VIEW v", db: "DB"},
		{name: "other alter actions not judged", sql: "ALTER TABLE t ADD COLUMN c INT; ALTER TABLE t RENAME TO u", db: "DB"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []compatFinding
			for _, finding := range SchemaBackwardCompatibility(file.Stmts, tt.db) {
				got = append(got, compatFinding{finding.Code, finding.Name, finding.Statement, tt.sql[finding.Start:finding.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
