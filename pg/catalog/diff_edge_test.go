package catalog

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDiffEdge(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, d *SchemaDiff)
	}{
		{
			name: "deterministic output — same input always produces same order",
			fromSQL: `
				CREATE TABLE t1 (id int);
				CREATE TABLE t2 (id int);
				CREATE TABLE t3 (id int);
			`,
			toSQL: `
				CREATE TABLE t2 (id int, name text);
				CREATE TABLE t3 (id int);
				CREATE TABLE t4 (id int);
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				// We already have the first diff result. Run Diff 9 more times
				// and verify all produce the exact same result.
				from, err := LoadSQL(`
					CREATE TABLE t1 (id int);
					CREATE TABLE t2 (id int);
					CREATE TABLE t3 (id int);
				`)
				if err != nil {
					t.Fatal(err)
				}
				to, err := LoadSQL(`
					CREATE TABLE t2 (id int, name text);
					CREATE TABLE t3 (id int);
					CREATE TABLE t4 (id int);
				`)
				if err != nil {
					t.Fatal(err)
				}

				// Capture a reference result.
				ref := Diff(from, to)

				for i := 0; i < 9; i++ {
					got := Diff(from, to)
					if !reflect.DeepEqual(summarizeDiff(ref), summarizeDiff(got)) {
						t.Fatalf("iteration %d: diff output is not deterministic", i+1)
					}
				}
			},
		},
		{
			name: "many objects of all types simultaneously",
			fromSQL: `
				CREATE TABLE t1 (id int PRIMARY KEY);
				CREATE TABLE t2 (id int);
				CREATE SEQUENCE seq1;
				CREATE TYPE color AS ENUM ('red', 'green');
				CREATE FUNCTION f1() RETURNS void LANGUAGE sql AS 'SELECT 1';
				CREATE VIEW v1 AS SELECT 1 AS x;
			`,
			toSQL: `
				CREATE TABLE t1 (id int PRIMARY KEY, name text);
				CREATE TABLE t3 (id int);
				CREATE SEQUENCE seq2;
				CREATE TYPE color AS ENUM ('red', 'green', 'blue');
				CREATE FUNCTION f2() RETURNS void LANGUAGE sql AS 'SELECT 2';
				CREATE VIEW v2 AS SELECT 2 AS y;
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				// Relations: t1 modified, t2 dropped, t3 added, v1 dropped, v2 added
				if len(d.Relations) < 4 {
					t.Fatalf("expected at least 4 relation diffs, got %d", len(d.Relations))
				}
				// Sequences: seq1 dropped, seq2 added
				if len(d.Sequences) < 2 {
					t.Fatalf("expected at least 2 sequence diffs, got %d", len(d.Sequences))
				}
				// Enums: color modified (new value)
				if len(d.Enums) < 1 {
					t.Fatalf("expected at least 1 enum diff, got %d", len(d.Enums))
				}
				// Functions: f1 dropped, f2 added
				if len(d.Functions) < 2 {
					t.Fatalf("expected at least 2 function diffs, got %d", len(d.Functions))
				}
			},
		},
		{
			name:    "empty schema vs non-existent schema",
			fromSQL: "",
			toSQL:   "CREATE SCHEMA s1;",
			check: func(t *testing.T, d *SchemaDiff) {
				// Schema s1 should appear as added.
				found := false
				for _, s := range d.Schemas {
					if s.Name == "s1" && s.Action == DiffAdd {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected schema s1 to be added, got schemas: %v", d.Schemas)
				}
				// No relation diffs since the schema is empty.
				for _, r := range d.Relations {
					if r.SchemaName == "s1" {
						t.Fatalf("expected no relation diffs in s1, got %+v", r)
					}
				}
			},
		},
		{
			name:    "object with reserved word name",
			fromSQL: "",
			toSQL: `CREATE TABLE "order" (id int, "select" text);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "order" {
					t.Fatalf("expected table name 'order', got %q", e.Name)
				}
			},
		},
		{
			name: "FK referencing table in different schema resolves correctly",
			fromSQL: `
				CREATE SCHEMA s1;
				CREATE SCHEMA s2;
				CREATE TABLE s1.parent (id int PRIMARY KEY);
				CREATE TABLE s2.child (id int, parent_id int REFERENCES s1.parent(id));
			`,
			toSQL: `
				CREATE SCHEMA s1;
				CREATE SCHEMA s2;
				CREATE TABLE s1.parent (id int PRIMARY KEY);
				CREATE TABLE s2.child (id int, parent_id int REFERENCES s1.parent(id));
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				// Identical schemas — no diff expected.
				if len(d.Relations) != 0 {
					t.Fatalf("expected 0 relation diffs for identical cross-schema FK, got %d", len(d.Relations))
				}
			},
		},
		{
			name: "FK cross-schema added",
			fromSQL: `
				CREATE SCHEMA s1;
				CREATE SCHEMA s2;
				CREATE TABLE s1.parent (id int PRIMARY KEY);
				CREATE TABLE s2.child (id int, parent_id int);
			`,
			toSQL: `
				CREATE SCHEMA s1;
				CREATE SCHEMA s2;
				CREATE TABLE s1.parent (id int PRIMARY KEY);
				CREATE TABLE s2.child (id int, parent_id int REFERENCES s1.parent(id));
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				// child should be modified with a new FK constraint.
				found := false
				for _, r := range d.Relations {
					if r.SchemaName == "s2" && r.Name == "child" && r.Action == DiffModify {
						for _, c := range r.Constraints {
							if c.Action == DiffAdd && c.To != nil && c.To.Type == ConstraintFK {
								found = true
							}
						}
					}
				}
				if !found {
					t.Fatalf("expected FK constraint added on s2.child referencing s1.parent")
				}
			},
		},
		{
			name: "diff is symmetric — Diff(A,B) adds == Diff(B,A) drops",
			fromSQL: `
				CREATE TABLE t1 (id int);
				CREATE TABLE t2 (id int);
				CREATE SEQUENCE seq1;
			`,
			toSQL: `
				CREATE TABLE t2 (id int);
				CREATE TABLE t3 (id int);
				CREATE SEQUENCE seq2;
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				// Load both catalogs to compute the reverse diff.
				fromCat, err := LoadSQL(`
					CREATE TABLE t1 (id int);
					CREATE TABLE t2 (id int);
					CREATE SEQUENCE seq1;
				`)
				if err != nil {
					t.Fatal(err)
				}
				toCat, err := LoadSQL(`
					CREATE TABLE t2 (id int);
					CREATE TABLE t3 (id int);
					CREATE SEQUENCE seq2;
				`)
				if err != nil {
					t.Fatal(err)
				}

				forward := Diff(fromCat, toCat)
				reverse := Diff(toCat, fromCat)

				// Count adds in forward, drops in reverse.
				fwdAdds := countRelActions(forward, DiffAdd)
				revDrops := countRelActions(reverse, DiffDrop)
				if fwdAdds != revDrops {
					t.Fatalf("forward adds (%d) != reverse drops (%d)", fwdAdds, revDrops)
				}

				// Count drops in forward, adds in reverse.
				fwdDrops := countRelActions(forward, DiffDrop)
				revAdds := countRelActions(reverse, DiffAdd)
				if fwdDrops != revAdds {
					t.Fatalf("forward drops (%d) != reverse adds (%d)", fwdDrops, revAdds)
				}

				// Also check sequences.
				fwdSeqAdds := countSeqActions(forward, DiffAdd)
				revSeqDrops := countSeqActions(reverse, DiffDrop)
				if fwdSeqAdds != revSeqDrops {
					t.Fatalf("forward seq adds (%d) != reverse seq drops (%d)", fwdSeqAdds, revSeqDrops)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, err := LoadSQL(tt.fromSQL)
			if err != nil {
				t.Fatal(err)
			}
			to, err := LoadSQL(tt.toSQL)
			if err != nil {
				t.Fatal(err)
			}
			d := Diff(from, to)
			tt.check(t, d)
		})
	}
}

// summarizeDiff produces a deterministic string summary of a SchemaDiff for comparison.
func summarizeDiff(d *SchemaDiff) string {
	var parts []string
	for _, s := range d.Schemas {
		parts = append(parts, fmt.Sprintf("schema:%d:%s", s.Action, s.Name))
	}
	for _, r := range d.Relations {
		parts = append(parts, fmt.Sprintf("rel:%d:%s.%s", r.Action, r.SchemaName, r.Name))
		for _, c := range r.Columns {
			parts = append(parts, fmt.Sprintf("  col:%d:%s", c.Action, c.Name))
		}
		for _, c := range r.Constraints {
			parts = append(parts, fmt.Sprintf("  con:%d:%s", c.Action, c.Name))
		}
		for _, idx := range r.Indexes {
			parts = append(parts, fmt.Sprintf("  idx:%d:%s", idx.Action, idx.Name))
		}
		for _, tr := range r.Triggers {
			parts = append(parts, fmt.Sprintf("  trig:%d:%s", tr.Action, tr.Name))
		}
	}
	for _, s := range d.Sequences {
		parts = append(parts, fmt.Sprintf("seq:%d:%s.%s", s.Action, s.SchemaName, s.Name))
	}
	for _, f := range d.Functions {
		parts = append(parts, fmt.Sprintf("func:%d:%s", f.Action, f.Identity))
	}
	for _, e := range d.Enums {
		parts = append(parts, fmt.Sprintf("enum:%d:%s.%s", e.Action, e.SchemaName, e.Name))
	}
	for _, dm := range d.Domains {
		parts = append(parts, fmt.Sprintf("domain:%d:%s.%s", dm.Action, dm.SchemaName, dm.Name))
	}
	for _, rg := range d.Ranges {
		parts = append(parts, fmt.Sprintf("range:%d:%s.%s", rg.Action, rg.SchemaName, rg.Name))
	}
	return fmt.Sprintf("%v", parts)
}

func countRelActions(d *SchemaDiff, action DiffAction) int {
	count := 0
	for _, r := range d.Relations {
		if r.Action == action {
			count++
		}
	}
	return count
}

func countSeqActions(d *SchemaDiff, action DiffAction) int {
	count := 0
	for _, s := range d.Sequences {
		if s.Action == action {
			count++
		}
	}
	return count
}
