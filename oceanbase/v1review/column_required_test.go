package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnRequired(t *testing.T) {
	required := []string{"id", "created_ts", "updated_ts"}
	type missing struct {
		Name, Table, Anchor string
		Statement           int
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []missing
	}{
		{name: "create with every column", sql: "CREATE TABLE t(id INT, created_ts INT, updated_ts INT);"},
		{
			name: "tables in name order, columns in name order",
			sql:  "CREATE TABLE t1(a INT);\nCREATE TABLE book(id INT, updated_ts INT);",
			want: []missing{
				{"created_ts", "book", "CREATE TABLE book(id INT, updated_ts INT)", 1},
				{"created_ts", "t1", "CREATE TABLE t1(a INT)", 0},
				{"id", "t1", "CREATE TABLE t1(a INT)", 0},
				{"updated_ts", "t1", "CREATE TABLE t1(a INT)", 0},
			},
		},
		{
			name: "drop column anchors the action",
			sql:  "CREATE TABLE t(id INT, created_ts INT, updated_ts INT);\nALTER TABLE t ADD COLUMN x INT, DROP COLUMN id;",
			want: []missing{{"id", "t", "DROP COLUMN id", 1}},
		},
		{
			name: "rename and change drop the old name",
			sql:  "ALTER TABLE t RENAME COLUMN id TO pk;\nALTER TABLE u CHANGE COLUMN updated_ts u INT;",
			want: []missing{
				{"id", "t", "RENAME COLUMN id TO pk", 0},
				{"updated_ts", "u", "CHANGE COLUMN updated_ts u INT", 1},
			},
		},
		{
			name: "an unknown table is taken to be complete",
			sql:  "ALTER TABLE t ADD COLUMN x INT; ALTER TABLE t RENAME COLUMN x TO id;",
		},
		{
			name: "add restores a dropped column",
			sql:  "ALTER TABLE t DROP COLUMN id; ALTER TABLE t ADD COLUMN (id INT);",
		},
		{
			name: "drop table forgets the table",
			sql:  "CREATE TABLE t(a INT); DROP TABLE IF EXISTS t;",
		},
		{
			name: "create starts the table afresh",
			sql:  "ALTER TABLE t ADD COLUMN x INT; CREATE TABLE t(id INT, created_ts INT);",
			want: []missing{{"updated_ts", "t", "CREATE TABLE t(id INT, created_ts INT)", 1}},
		},
		{
			name: "create like has no columns, qualifier dropped",
			sql:  "CREATE TABLE db.t LIKE u;",
			want: []missing{
				{"created_ts", "t", "CREATE TABLE db.t LIKE u", 0},
				{"id", "t", "CREATE TABLE db.t LIKE u", 0},
				{"updated_ts", "t", "CREATE TABLE db.t LIKE u", 0},
			},
		},
		{
			name: "modify column is not judged",
			sql:  "CREATE TABLE t(id INT, created_ts INT, updated_ts INT); ALTER TABLE t MODIFY COLUMN id BIGINT;",
		},
		{
			name: "a routine body is not judged",
			sql:  "CREATE PROCEDURE p() BEGIN CREATE TABLE t(a INT); END;",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stmts := parse(t, tt.sql)
			var got []missing
			for _, f := range ColumnRequired(stmts, required) {
				if f.Code != v1review.ColumnMissing {
					t.Errorf("code %q", f.Code)
				}
				got = append(got, missing{f.Name, f.Table, tt.sql[f.Start:f.End], f.Statement})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ColumnRequired mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestColumnRequiredNoneRequired(t *testing.T) {
	if got := ColumnRequired(parse(t, "CREATE TABLE t(a INT);"), nil); got != nil {
		t.Errorf("got %v", got)
	}
}
