package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

// columnFinding is a column finding as the tests compare it: the rest of the
// line its range starts on stands in for the range, since the column rules
// mean only the start.
type columnFinding struct {
	Code  v1review.Code
	Name  string
	Table string
	At    string
}

func columnFindings(sql string, findings []v1review.Finding) []columnFinding {
	var got []columnFinding
	for _, f := range findings {
		at, _, _ := strings.Cut(sql[f.Start:], "\n")
		got = append(got, columnFinding{f.Code, f.Name, f.Table, at})
	}
	return got
}

func TestColumnMaximumVarcharLength(t *testing.T) {
	long := func(at string) columnFinding { return columnFinding{Code: v1review.ColumnVarcharTooLong, At: at} }
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{name: "within the limit", sql: "CREATE TABLE t(a VARCHAR(100));"},
		{
			name: "any case, no length is the default, synonyms not judged",
			sql:  "CREATE TABLE t(\n  a varchar(101),\n  b VARCHAR,\n  c STRING,\n  d TEXT,\n  e NVARCHAR(500)\n);",
			want: []columnFinding{long("varchar(101),"), long("VARCHAR,")},
		},
		{
			name: "nested type names are judged",
			sql:  "SELECT a::VARCHAR(500) FROM t;",
			want: []columnFinding{long("VARCHAR(500) FROM t;")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnMaximumVarcharLength(parse(t, tt.sql), 100))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestColumnNoNull(t *testing.T) {
	nullable := func(name, at string) columnFinding {
		return columnFinding{Code: v1review.ColumnNullable, Name: name, At: at}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{name: "not null and primary keys", sql: "CREATE TABLE t(a INT NOT NULL, b INT PRIMARY KEY, c INT, PRIMARY KEY (C));"},
		{
			name: "names keyed as resolved, in key order, last definition wins",
			sql:  "CREATE TABLE t(\n  z INT,\n  \"q\" INT,\n  a INT NOT NULL,\n  a INT\n);",
			want: []columnFinding{nullable("A", "a INT"), nullable("Z", "z INT,"), nullable("q", "\"q\" INT,")},
		},
		{
			name: "alter: add, drop not null, add primary key",
			sql:  "ALTER TABLE t ADD COLUMN x INT, y INT;\nALTER TABLE t ALTER COLUMN id DROP NOT NULL;\nALTER TABLE t ADD PRIMARY KEY (y);",
			want: []columnFinding{nullable("X", "x INT, y INT;"), nullable("Y", "y INT;"), nullable("ID", "id DROP NOT NULL;")},
		},
		{name: "create as select, like and clone are not judged", sql: "CREATE TABLE t AS SELECT 1 a;\nCREATE TABLE u LIKE t;\nCREATE TABLE v CLONE t;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnNoNull(parse(t, tt.sql)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestColumnRequired(t *testing.T) {
	missing := func(name, table, at string) columnFinding {
		return columnFinding{Code: v1review.ColumnMissing, Name: name, Table: table, At: at}
	}
	required := []string{"ID", "id", "CREATED_TS", "ID"}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{name: "create with every entry", sql: "CREATE TABLE t(id INT, \"id\" INT, created_ts INT);"},
		{
			name: "create lacks entries, anchored at the last item's last byte",
			sql:  "CREATE TABLE db.\"S\".t(\n  a INT,\n  CONSTRAINT pk PRIMARY KEY (a)\n);",
			want: []columnFinding{
				missing("CREATED_TS", "db.\"S\".t", ")"),
				missing("ID", "db.\"S\".t", ")"),
				missing("id", "db.\"S\".t", ")"),
			},
		},
		{
			name: "drop, once per entry, anchored at the first drop action",
			sql:  "ALTER TABLE t DROP COLUMN x, id, ID;\nALTER TABLE t DROP COLUMN \"id\";",
			want: []columnFinding{
				missing("ID", "t", "DROP COLUMN x, id, ID;"),
				missing("id", "t", "DROP COLUMN \"id\";"),
			},
		},
		{name: "add and rename are not judged", sql: "ALTER TABLE t ADD COLUMN id INT;\nALTER TABLE t RENAME COLUMN id TO pk;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnRequired(parse(t, tt.sql), required))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
