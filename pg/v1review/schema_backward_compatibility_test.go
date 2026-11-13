package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestSchemaBackwardCompatibility(t *testing.T) {
	at := func(code v1review.Code, stmt int, text string) systemFinding {
		return systemFinding{Code: code, Statement: stmt, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []systemFinding
	}{
		{name: "additive changes", sql: "CREATE TABLE t (a int); ALTER TABLE u ADD COLUMN b int; CREATE INDEX i ON u (a); ALTER TABLE u ALTER COLUMN a SET DEFAULT 1;"},
		{name: "drop database", sql: "DROP DATABASE d;", want: []systemFinding{at(v1review.IncompatibleDropDatabase, 0, "DROP DATABASE d")}},
		{
			name: "drop table and materialized view",
			sql:  "DROP TABLE a, b; DROP MATERIALIZED VIEW m;",
			want: []systemFinding{
				at(v1review.IncompatibleDropTable, 0, "DROP TABLE a, b"),
				at(v1review.IncompatibleDropTable, 1, "DROP MATERIALIZED VIEW m"),
			},
		},
		{name: "other drops", sql: "DROP VIEW v; DROP INDEX i; DROP SCHEMA s; DROP SEQUENCE q; DROP FUNCTION f();"},
		{
			name: "renames",
			sql:  "ALTER TABLE a RENAME TO b; ALTER INDEX i RENAME TO j; ALTER SCHEMA s RENAME TO r; ALTER TABLE t RENAME CONSTRAINT c TO d; ALTER TABLE t RENAME COLUMN a TO b;",
			want: []systemFinding{
				at(v1review.IncompatibleRenameTable, 0, "ALTER TABLE a RENAME TO b"),
				at(v1review.IncompatibleRenameTable, 1, "ALTER INDEX i RENAME TO j"),
				// The parser keeps no range for a schema rename.
				at(v1review.IncompatibleRenameTable, 2, ""),
				at(v1review.IncompatibleRenameTable, 3, "ALTER TABLE t RENAME CONSTRAINT c TO d"),
				at(v1review.IncompatibleRenameColumn, 4, "ALTER TABLE t RENAME COLUMN a TO b"),
			},
		},
		{
			name: "alter table reports its first incompatible action",
			sql:  "ALTER TABLE t ADD COLUMN c int, ALTER COLUMN a TYPE bigint, DROP COLUMN b;",
			want: []systemFinding{at(v1review.IncompatibleAlterColumn, 0, "ALTER TABLE t ADD COLUMN c int, ALTER COLUMN a TYPE bigint, DROP COLUMN b")},
		},
		{
			name: "alter table constraints",
			sql: "ALTER TABLE t DROP COLUMN a; ALTER TABLE t ADD PRIMARY KEY (a); ALTER TABLE t ADD CONSTRAINT u UNIQUE (a); " +
				"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u (a); ALTER TABLE t ADD CHECK (a > 0); ALTER TABLE t ADD CHECK (a > 0) NOT VALID; " +
				"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u (a) NOT VALID; ALTER TABLE t ADD EXCLUDE USING gist (a WITH =);",
			want: []systemFinding{
				at(v1review.IncompatibleDropColumn, 0, "ALTER TABLE t DROP COLUMN a"),
				at(v1review.IncompatibleAddPrimaryKey, 1, "ALTER TABLE t ADD PRIMARY KEY (a)"),
				at(v1review.IncompatibleAddUniqueKey, 2, "ALTER TABLE t ADD CONSTRAINT u UNIQUE (a)"),
				at(v1review.IncompatibleAddForeignKey, 3, "ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u (a)"),
				at(v1review.IncompatibleAddCheck, 4, "ALTER TABLE t ADD CHECK (a > 0)"),
				at(v1review.IncompatibleAddForeignKey, 6, "ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u (a) NOT VALID"),
			},
		},
		{
			name: "a check added not valid does not stop the search",
			sql:  "ALTER TABLE t ADD CHECK (a > 0) NOT VALID, DROP COLUMN b;",
			want: []systemFinding{at(v1review.IncompatibleDropColumn, 0, "ALTER TABLE t ADD CHECK (a > 0) NOT VALID, DROP COLUMN b")},
		},
		{
			name: "create unique index",
			sql:  "CREATE UNIQUE INDEX i ON t (a); CREATE INDEX j ON t (a);",
			want: []systemFinding{at(v1review.IncompatibleAddUniqueKey, 0, "CREATE UNIQUE INDEX i ON t (a)")},
		},
		{
			name: "the last created table may change freely",
			sql:  "CREATE TABLE t (a int); ALTER TABLE t DROP COLUMN a; CREATE UNIQUE INDEX i ON t (a); ALTER TABLE t RENAME COLUMN a TO b; ALTER TABLE t RENAME TO u; DROP TABLE t;",
			want: []systemFinding{
				at(v1review.IncompatibleRenameTable, 4, "ALTER TABLE t RENAME TO u"),
				at(v1review.IncompatibleDropTable, 5, "DROP TABLE t"),
			},
		},
		{
			name: "only the last created table is remembered",
			sql:  "CREATE TABLE t (a int); CREATE TABLE u (a int); ALTER TABLE t DROP COLUMN a; ALTER TABLE u DROP COLUMN a;",
			want: []systemFinding{at(v1review.IncompatibleDropColumn, 2, "ALTER TABLE t DROP COLUMN a")},
		},
		{
			name: "the remembered table is named without its schema",
			sql:  `CREATE TABLE s1.t (a int); ALTER TABLE s2.t DROP COLUMN a; ALTER TABLE "T" DROP COLUMN a; ALTER TABLE T DROP COLUMN a;`,
			want: []systemFinding{at(v1review.IncompatibleDropColumn, 2, `ALTER TABLE "T" DROP COLUMN a`)},
		},
		{
			name: "create table as is not remembered",
			sql:  "CREATE TABLE t AS SELECT 1 AS a; ALTER TABLE t DROP COLUMN a;",
			want: []systemFinding{at(v1review.IncompatibleDropColumn, 1, "ALTER TABLE t DROP COLUMN a")},
		},
		{
			name: "an alter of another relation kind is judged as an alter table",
			sql:  "ALTER FOREIGN TABLE f DROP COLUMN a; ALTER VIEW v ALTER COLUMN a SET DEFAULT 1; ALTER INDEX i SET (fillfactor = 70);",
			want: []systemFinding{at(v1review.IncompatibleDropColumn, 0, "ALTER FOREIGN TABLE f DROP COLUMN a")},
		},
		{
			name: "a routine body and explain are not judged",
			sql:  "CREATE FUNCTION f() RETURNS void LANGUAGE sql AS $$ DROP TABLE t $$; EXPLAIN SELECT 1;",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, SchemaBackwardCompatibility)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
