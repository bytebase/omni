package analysis

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

// numberedList joins format for i = 1..n with sep, every verb given i.
func numberedList(n int, format, sep string) string {
	verbs := strings.Count(format, "%d")
	parts := make([]string, n)
	for i := 1; i <= n; i++ {
		args := make([]any, verbs)
		for j := range args {
			args[j] = i
		}
		parts[i-1] = fmt.Sprintf(format, args...)
	}
	return strings.Join(parts, sep)
}

// querySpanWithin fails the test when GetQuerySpan(sql) errors or takes
// longer than limit.
func querySpanWithin(t *testing.T, sql string, limit time.Duration) *QuerySpan {
	t.Helper()
	type result struct {
		span *QuerySpan
		err  error
	}
	done := make(chan result, 1)
	go func() {
		span, err := GetQuerySpan(sql, DialectBigQuery)
		done <- result{span, err}
	}()
	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("GetQuerySpan: %v", r.err)
		}
		return r.span
	case <-time.After(limit):
		t.Fatalf("GetQuerySpan took over %v", limit)
	}
	return nil
}

// TestQuerySpanWideListsStayLinear guards the wide lists GetQuerySpan
// resolved by scanning a list once per entry: a USING list, a select list
// over a wide FROM list or a wide derived table, and the arms of a BY NAME
// set operation. Names are compared case-insensitively and the first match
// wins, which an index on the folded name now answers. It also guards two
// chains: a BY NAME chain whose arms each add a column, which checked the
// whole projection so far at every arm, and a USING chain on one key over
// different tables, whose key gains a table at every join and was rebuilt at
// every join. At these sizes each shape took 5-14s; it takes well under a
// second now.
func TestQuerySpanWideListsStayLinear(t *testing.T) {
	for _, tc := range []struct {
		name string
		sql  string
	}{
		{name: "USING list", sql: "SELECT * FROM t JOIN u USING (" + numberedList(40000, "k%d", ", ") + ")"},
		{name: "USING list over derived tables", sql: "SELECT * FROM (SELECT " + numberedList(16000, "a%d AS k%d", ", ") + " FROM t) JOIN (SELECT " +
			numberedList(16000, "b%d AS k%d", ", ") + " FROM u) USING (" + numberedList(16000, "k%d", ", ") + ")"},
		{name: "USING list over a USING join", sql: "SELECT * FROM (SELECT * FROM t JOIN u USING (" + numberedList(16000, "k%d", ", ") +
			")) JOIN v USING (" + numberedList(16000, "k%d", ", ") + ")"},
		{name: "qualified columns over many tables", sql: "SELECT " + numberedList(32000, "t%d.a", ", ") + " FROM " + numberedList(32000, "t%d", ", ")},
		{name: "dataset-qualified columns over many tables", sql: "SELECT " + numberedList(26000, "ds.t%d.a", ", ") + " FROM " + numberedList(26000, "ds.t%d", ", ")},
		{name: "bare columns over many derived tables", sql: "SELECT " + numberedList(22000, "c%d", ", ") + " FROM " + numberedList(22000, "(SELECT a AS c%d FROM t) x%d", ", ")},
		{name: "qualified columns of a wide derived table", sql: "SELECT " + numberedList(44000, "x.c%d", ", ") + " FROM (SELECT " + numberedList(44000, "a%d AS c%d", ", ") + " FROM t) x"},
		{name: "columns behind the stars of a USING join", sql: "SELECT " + numberedList(22000, "x.a%d", ", ") + " FROM (SELECT * FROM t JOIN u USING (" + numberedList(22000, "k%d", ", ") + ")) x"},
		{name: "BY NAME", sql: "SELECT " + numberedList(28000, "a%d", ", ") + " FROM t UNION ALL BY NAME SELECT " + numberedList(28000, "a%d", ", ") + " FROM u"},
		{name: "BY NAME ON", sql: "SELECT " + numberedList(22000, "a%d", ", ") + " FROM t UNION ALL BY NAME ON (" + numberedList(22000, "a%d", ", ") +
			") SELECT " + numberedList(22000, "a%d", ", ") + " FROM u"},
		{name: "BY NAME arms each adding a column", sql: "SELECT a0 FROM t" + numberedList(26000, " UNION ALL BY NAME SELECT a%d FROM t", "")},
		{name: "USING chain on one key over different tables", sql: "SELECT * FROM t0" + numberedList(6000, " JOIN t%d USING (k)", "")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			querySpanWithin(t, tc.sql, 3*time.Second)
		})
	}
}

// TestRecursiveCTERotationIsNotCubic guards the recursive-CTE fixpoint on an
// arm that rotates its columns. A source moves one column per pass, so the
// fixpoint takes a pass per column, and each pass resolved the arm and
// unioned every column's sources again: 400 columns took 15-20s. The result
// itself holds every source in every column, so it is quadratic.
func TestRecursiveCTERotationIsNotCubic(t *testing.T) {
	// The arm reads c2, c3, ..., cn, c1 into c1, c2, ..., cn.
	rotation := func(n int) string {
		arm := make([]string, n)
		for i := range arm {
			arm[i] = fmt.Sprintf("c%d", (i+1)%n+1)
		}
		return "WITH RECURSIVE r AS (SELECT " + numberedList(n, "a%d AS c%d", ", ") + " FROM t UNION ALL SELECT " +
			strings.Join(arm, ", ") + " FROM r WHERE c1 < 10) SELECT * FROM r"
	}

	// Each column gains the next column's source at every pass, in that order.
	span := querySpanWithin(t, rotation(3), time.Second)
	ref := func(c string) ColumnRef { return ColumnRef{Table: "t", Column: c} }
	want := []ColumnInfo{
		{Name: "c1", SourceColumns: []ColumnRef{ref("a1"), ref("a2"), ref("a3")}},
		{Name: "c2", SourceColumns: []ColumnRef{ref("a2"), ref("a3"), ref("a1")}},
		{Name: "c3", SourceColumns: []ColumnRef{ref("a3"), ref("a1"), ref("a2")}},
	}
	if !reflect.DeepEqual(span.Results, want) {
		t.Fatalf("Results = %+v, want %+v", span.Results, want)
	}

	const n = 400
	span = querySpanWithin(t, rotation(n), 3*time.Second)
	if len(span.Results) != n {
		t.Fatalf("got %d result columns, want %d", len(span.Results), n)
	}
	for _, col := range span.Results {
		if len(col.SourceColumns) != n {
			t.Fatalf("column %s has %d sources, want %d", col.Name, len(col.SourceColumns), n)
		}
	}
}
