package catalog

import (
	"testing"
)

func TestDiffConstraint(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, d *SchemaDiff)
	}{
		{
			name:    "PRIMARY KEY added",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "CREATE TABLE t1 (id int PRIMARY KEY);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				// Find constraint diff entries (may also have column diffs for NOT NULL).
				var pkDiffs []ConstraintDiffEntry
				for _, cd := range e.Constraints {
					if cd.Action == DiffAdd && cd.To != nil && cd.To.Type == ConstraintPK {
						pkDiffs = append(pkDiffs, cd)
					}
				}
				if len(pkDiffs) != 1 {
					t.Fatalf("expected 1 PK add, got %d (total constraint diffs: %d)", len(pkDiffs), len(e.Constraints))
				}
			},
		},
		{
			name:    "PRIMARY KEY dropped",
			fromSQL: "CREATE TABLE t1 (id int PRIMARY KEY);",
			toSQL:   "CREATE TABLE t1 (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				var pkDiffs []ConstraintDiffEntry
				for _, cd := range e.Constraints {
					if cd.Action == DiffDrop && cd.From != nil && cd.From.Type == ConstraintPK {
						pkDiffs = append(pkDiffs, cd)
					}
				}
				if len(pkDiffs) != 1 {
					t.Fatalf("expected 1 PK drop, got %d (total constraint diffs: %d)", len(pkDiffs), len(e.Constraints))
				}
			},
		},
		{
			name: "UNIQUE constraint added",
			fromSQL: `CREATE TABLE t1 (id int, email text);`,
			toSQL:   `CREATE TABLE t1 (id int, email text, CONSTRAINT uq_email UNIQUE (email));`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				found := false
				for _, cd := range e.Constraints {
					if cd.Action == DiffAdd && cd.To != nil && cd.To.Type == ConstraintUnique {
						found = true
						if cd.Name != "uq_email" {
							t.Fatalf("expected constraint name uq_email, got %s", cd.Name)
						}
					}
				}
				if !found {
					t.Fatalf("expected UNIQUE add constraint diff")
				}
			},
		},
		{
			name: "UNIQUE constraint dropped",
			fromSQL: `CREATE TABLE t1 (id int, email text, CONSTRAINT uq_email UNIQUE (email));`,
			toSQL:   `CREATE TABLE t1 (id int, email text);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				found := false
				for _, cd := range e.Constraints {
					if cd.Action == DiffDrop && cd.From != nil && cd.From.Type == ConstraintUnique {
						found = true
					}
				}
				if !found {
					t.Fatalf("expected UNIQUE drop constraint diff")
				}
			},
		},
		{
			name: "CHECK constraint added",
			fromSQL: `CREATE TABLE t1 (id int, age int);`,
			toSQL:   `CREATE TABLE t1 (id int, age int, CONSTRAINT chk_age CHECK (age > 0));`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				found := false
				for _, cd := range e.Constraints {
					if cd.Action == DiffAdd && cd.To != nil && cd.To.Type == ConstraintCheck {
						found = true
						if cd.Name != "chk_age" {
							t.Fatalf("expected constraint name chk_age, got %s", cd.Name)
						}
					}
				}
				if !found {
					t.Fatalf("expected CHECK add constraint diff")
				}
			},
		},
		{
			name: "CHECK constraint dropped",
			fromSQL: `CREATE TABLE t1 (id int, age int, CONSTRAINT chk_age CHECK (age > 0));`,
			toSQL:   `CREATE TABLE t1 (id int, age int);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				found := false
				for _, cd := range e.Constraints {
					if cd.Action == DiffDrop && cd.From != nil && cd.From.Type == ConstraintCheck {
						found = true
					}
				}
				if !found {
					t.Fatalf("expected CHECK drop constraint diff")
				}
			},
		},
		{
			name: "CHECK constraint expression changed",
			fromSQL: `CREATE TABLE t1 (id int, age int, CONSTRAINT chk_age CHECK (age > 0));`,
			toSQL:   `CREATE TABLE t1 (id int, age int, CONSTRAINT chk_age CHECK (age > 18));`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				found := false
				for _, cd := range e.Constraints {
					if cd.Action == DiffModify && cd.From != nil && cd.From.Type == ConstraintCheck {
						found = true
						if cd.From.CheckExpr == cd.To.CheckExpr {
							t.Fatalf("expected CHECK expressions to differ")
						}
					}
				}
				if !found {
					t.Fatalf("expected CHECK modify constraint diff")
				}
			},
		},
		{
			name: "Foreign key added",
			fromSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int);
			`,
			toSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int, CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id));
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				// Find the child table diff.
				var childDiff *RelationDiffEntry
				for i := range d.Relations {
					if d.Relations[i].Name == "child" {
						childDiff = &d.Relations[i]
						break
					}
				}
				if childDiff == nil {
					t.Fatalf("expected child table diff, got %d relations", len(d.Relations))
				}
				found := false
				for _, cd := range childDiff.Constraints {
					if cd.Action == DiffAdd && cd.To != nil && cd.To.Type == ConstraintFK {
						found = true
						if cd.Name != "fk_parent" {
							t.Fatalf("expected constraint name fk_parent, got %s", cd.Name)
						}
					}
				}
				if !found {
					t.Fatalf("expected FK add constraint diff")
				}
			},
		},
		{
			name: "Foreign key dropped",
			fromSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int, CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id));
			`,
			toSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int);
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				var childDiff *RelationDiffEntry
				for i := range d.Relations {
					if d.Relations[i].Name == "child" {
						childDiff = &d.Relations[i]
						break
					}
				}
				if childDiff == nil {
					t.Fatalf("expected child table diff, got %d relations", len(d.Relations))
				}
				found := false
				for _, cd := range childDiff.Constraints {
					if cd.Action == DiffDrop && cd.From != nil && cd.From.Type == ConstraintFK {
						found = true
					}
				}
				if !found {
					t.Fatalf("expected FK drop constraint diff")
				}
			},
		},
		{
			name: "Foreign key actions changed",
			fromSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int,
					CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE CASCADE);
			`,
			toSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int,
					CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE SET NULL);
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				var childDiff *RelationDiffEntry
				for i := range d.Relations {
					if d.Relations[i].Name == "child" {
						childDiff = &d.Relations[i]
						break
					}
				}
				if childDiff == nil {
					t.Fatalf("expected child table diff, got %d relations", len(d.Relations))
				}
				found := false
				for _, cd := range childDiff.Constraints {
					if cd.Action == DiffModify && cd.From != nil && cd.From.Type == ConstraintFK {
						found = true
						if cd.From.FKDelAction == cd.To.FKDelAction {
							t.Fatalf("expected FK delete action to differ")
						}
					}
				}
				if !found {
					t.Fatalf("expected FK modify constraint diff")
				}
			},
		},
		{
			name: "Foreign key match type changed",
			fromSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int,
					CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id));
			`,
			toSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int,
					CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id) MATCH FULL);
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				var childDiff *RelationDiffEntry
				for i := range d.Relations {
					if d.Relations[i].Name == "child" {
						childDiff = &d.Relations[i]
						break
					}
				}
				if childDiff == nil {
					t.Fatalf("expected child table diff, got %d relations", len(d.Relations))
				}
				found := false
				for _, cd := range childDiff.Constraints {
					if cd.Action == DiffModify && cd.From != nil && cd.From.Type == ConstraintFK {
						found = true
						if cd.From.FKMatchType == cd.To.FKMatchType {
							t.Fatalf("expected FK match type to differ")
						}
					}
				}
				if !found {
					t.Fatalf("expected FK modify constraint diff for match type change")
				}
			},
		},
		{
			name: "Constraint deferrable/deferred changed",
			fromSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int,
					CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id));
			`,
			toSQL: `
				CREATE TABLE parent (id int PRIMARY KEY);
				CREATE TABLE child (id int, parent_id int,
					CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent(id) DEFERRABLE INITIALLY DEFERRED);
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				var childDiff *RelationDiffEntry
				for i := range d.Relations {
					if d.Relations[i].Name == "child" {
						childDiff = &d.Relations[i]
						break
					}
				}
				if childDiff == nil {
					t.Fatalf("expected child table diff, got %d relations", len(d.Relations))
				}
				found := false
				for _, cd := range childDiff.Constraints {
					if cd.Action == DiffModify && cd.Name == "fk_parent" {
						found = true
						if cd.From.Deferrable == cd.To.Deferrable {
							t.Fatalf("expected Deferrable to differ")
						}
						if cd.From.Deferred == cd.To.Deferred {
							t.Fatalf("expected Deferred to differ")
						}
					}
				}
				if !found {
					t.Fatalf("expected FK modify constraint diff for deferrable change")
				}
			},
		},
		{
			name: "EXCLUDE constraint added and dropped",
			fromSQL: `
				CREATE TABLE t1 (id int, val int, CONSTRAINT excl_val EXCLUDE USING gist (val WITH =));
			`,
			toSQL: `
				CREATE TABLE t1 (id int, val int);
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				found := false
				for _, cd := range e.Constraints {
					if cd.Action == DiffDrop && cd.From != nil && cd.From.Type == ConstraintExclude {
						found = true
					}
				}
				if !found {
					t.Fatalf("expected EXCLUDE drop constraint diff")
				}
			},
		},
		{
			name: "Constraint identity by name",
			fromSQL: `
				CREATE TABLE t1 (id int, age int, CONSTRAINT chk_age CHECK (age > 0));
			`,
			toSQL: `
				CREATE TABLE t1 (id int, age int, CONSTRAINT chk_age CHECK (age > 0));
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				// Same constraint name, same expression — no diff expected.
				if len(d.Relations) != 0 {
					t.Fatalf("expected 0 relation diffs, got %d", len(d.Relations))
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
