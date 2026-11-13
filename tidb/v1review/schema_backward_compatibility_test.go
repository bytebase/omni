package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestSchemaBackwardCompatibility(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []systemFinding
	}{
		{name: "drop database", sql: "DROP SCHEMA d;", want: []systemFinding{{Code: v1review.IncompatibleDropDatabase, Text: "DROP SCHEMA d"}}},
		{name: "drop table", sql: "DROP TABLE a, b;", want: []systemFinding{{Code: v1review.IncompatibleDropTable, Text: "DROP TABLE a, b"}}},
		{name: "rename table", sql: "RENAME TABLE a TO b;", want: []systemFinding{{Code: v1review.IncompatibleRenameTable, Text: "RENAME TABLE a TO b"}}},
		{name: "drop view and index are not reported", sql: "DROP VIEW v; DROP INDEX i ON t;"},
		{name: "alter rename table", sql: "ALTER TABLE a RENAME TO b;", want: []systemFinding{{Code: v1review.IncompatibleRenameTable, Text: "ALTER TABLE a RENAME TO b"}}},
		{name: "rename column", sql: "ALTER TABLE t RENAME COLUMN a TO b;", want: []systemFinding{{Code: v1review.IncompatibleRenameColumn, Text: "ALTER TABLE t RENAME COLUMN a TO b"}}},
		{name: "drop column", sql: "ALTER TABLE t DROP COLUMN a;", want: []systemFinding{{Code: v1review.IncompatibleDropColumn, Text: "ALTER TABLE t DROP COLUMN a"}}},
		{name: "add primary key", sql: "ALTER TABLE t ADD PRIMARY KEY (a);", want: []systemFinding{{Code: v1review.IncompatibleAddPrimaryKey, Text: "ALTER TABLE t ADD PRIMARY KEY (a)"}}},
		{name: "add unique key", sql: "ALTER TABLE t ADD UNIQUE KEY (a);", want: []systemFinding{{Code: v1review.IncompatibleAddUniqueKey, Text: "ALTER TABLE t ADD UNIQUE KEY (a)"}}},
		{name: "add foreign key", sql: "ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u (a);", want: []systemFinding{{Code: v1review.IncompatibleAddForeignKey, Text: "ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u (a)"}}},
		{name: "add check", sql: "ALTER TABLE t ADD CHECK (a > 0);", want: []systemFinding{{Code: v1review.IncompatibleAddCheck, Text: "ALTER TABLE t ADD CHECK (a > 0)"}}},
		{name: "add check not enforced is not reported", sql: "ALTER TABLE t ADD CHECK (a > 0) NOT ENFORCED;"},
		{name: "alter check enforced", sql: "ALTER TABLE t ALTER CONSTRAINT c ENFORCED;", want: []systemFinding{{Code: v1review.IncompatibleAlterCheck, Text: "ALTER TABLE t ALTER CONSTRAINT c ENFORCED"}}},
		{name: "alter check not enforced is not reported", sql: "ALTER TABLE t ALTER CHECK c NOT ENFORCED;"},
		{name: "modify column", sql: "ALTER TABLE t MODIFY a BIGINT;", want: []systemFinding{{Code: v1review.IncompatibleAlterColumn, Text: "ALTER TABLE t MODIFY a BIGINT"}}},
		{name: "change column", sql: "ALTER TABLE t CHANGE a b BIGINT;", want: []systemFinding{{Code: v1review.IncompatibleAlterColumn, Text: "ALTER TABLE t CHANGE a b BIGINT"}}},
		{name: "add column and index are not reported", sql: "ALTER TABLE t ADD COLUMN a INT, ADD INDEX (a); CREATE INDEX i ON t (a);"},
		{
			name: "first reported action only",
			sql:  "ALTER TABLE t ADD COLUMN c INT, ADD CHECK (a > 0) NOT ENFORCED, DROP COLUMN a, MODIFY b INT;",
			want: []systemFinding{{Code: v1review.IncompatibleDropColumn, Text: "ALTER TABLE t ADD COLUMN c INT, ADD CHECK (a > 0) NOT ENFORCED, DROP COLUMN a, MODIFY b INT"}},
		},
		{name: "unique index", sql: "CREATE UNIQUE INDEX i ON t (a);", want: []systemFinding{{Code: v1review.IncompatibleAddUniqueKey, Text: "CREATE UNIQUE INDEX i ON t (a)"}}},
		{
			name: "the table just created is free, a drop of it is not",
			sql:  "CREATE TABLE t (a INT); ALTER TABLE t DROP COLUMN a; CREATE UNIQUE INDEX i ON t (a); DROP TABLE t;",
			want: []systemFinding{{Code: v1review.IncompatibleDropTable, Statement: 3, Text: "DROP TABLE t"}},
		},
		{
			name: "only the last created table is remembered",
			sql:  "CREATE TABLE t (a INT); CREATE TABLE u (a INT); ALTER TABLE t DROP COLUMN a; ALTER TABLE u DROP COLUMN a;",
			want: []systemFinding{{Code: v1review.IncompatibleDropColumn, Statement: 2, Text: "ALTER TABLE t DROP COLUMN a"}},
		},
		{
			name: "a table created by LIKE or AS SELECT is remembered",
			sql:  "CREATE TABLE t LIKE u; ALTER TABLE t DROP COLUMN a; CREATE TABLE v AS SELECT 1 AS a; CREATE UNIQUE INDEX i ON v (a);",
		},
		{
			name: "remembered by name alone, as written",
			sql:  "CREATE TABLE d1.t (a INT); ALTER TABLE d2.t DROP COLUMN a; ALTER TABLE T DROP COLUMN a;",
			want: []systemFinding{{Code: v1review.IncompatibleDropColumn, Statement: 2, Text: "ALTER TABLE T DROP COLUMN a"}},
		},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN DROP TABLE t; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, SchemaBackwardCompatibility)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
