package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestIndexCreateConcurrently(t *testing.T) {
	sqls := []string{
		"CREATE INDEX i1 ON t (a)",
		"CREATE INDEX CONCURRENTLY i2 ON t (a)",
		"CREATE TABLE s1.t (a int)",
		"CREATE INDEX i3 ON s2.t (a)",
		"CREATE TABLE u AS SELECT 1 a",
		"CREATE UNIQUE INDEX ON u (a)",
		"DROP INDEX i1, i2",
		"DROP INDEX CONCURRENTLY i3",
		"DROP TABLE t",
		"CREATE FUNCTION f() RETURNS void AS $$ CREATE INDEX i ON v (a) $$ LANGUAGE sql",
	}
	created := func(name, table, text string) indexFinding {
		return indexFinding{Code: v1review.IndexCreatedNotConcurrently, Name: name, Table: table, Text: text}
	}
	want := []indexFinding{
		created("i1", "t", sqls[0]),
		// t counts as created whatever its schema.
		created("", "u", sqls[5]),
		{Code: v1review.IndexDroppedNotConcurrently, Text: sqls[6]},
	}
	got := indexFindings(sqls, IndexCreateConcurrently(parseEach(t, sqls)))
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
	var statements []int
	for _, f := range IndexCreateConcurrently(parseEach(t, sqls)) {
		statements = append(statements, f.Statement)
	}
	if diff := cmp.Diff([]int{0, 5, 6}, statements); diff != "" {
		t.Errorf("statements (-want +got):\n%s", diff)
	}
}
