package catalog

import "testing"

func TestDiffPolicy(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, d *SchemaDiff)
	}{
		{
			name: "policy added on table",
			fromSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;`,
			toSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 USING (owner_id = 1);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.Policies) != 1 {
					t.Fatalf("expected 1 policy diff, got %d", len(e.Policies))
				}
				p := e.Policies[0]
				if p.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", p.Action)
				}
				if p.Name != "my_pol" {
					t.Fatalf("expected my_pol, got %s", p.Name)
				}
				if p.From != nil {
					t.Fatalf("expected From to be nil")
				}
				if p.To == nil {
					t.Fatalf("expected To to be non-nil")
				}
			},
		},
		{
			name: "policy dropped from table",
			fromSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 USING (owner_id = 1);`,
			toSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Policies) != 1 {
					t.Fatalf("expected 1 policy diff, got %d", len(e.Policies))
				}
				p := e.Policies[0]
				if p.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", p.Action)
				}
				if p.Name != "my_pol" {
					t.Fatalf("expected my_pol, got %s", p.Name)
				}
				if p.From == nil {
					t.Fatalf("expected From to be non-nil")
				}
				if p.To != nil {
					t.Fatalf("expected To to be nil")
				}
			},
		},
		{
			name: "policy command type changed",
			fromSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 FOR SELECT USING (owner_id = 1);`,
			toSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 FOR UPDATE USING (owner_id = 1);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Policies) != 1 {
					t.Fatalf("expected 1 policy diff, got %d", len(e.Policies))
				}
				p := e.Policies[0]
				if p.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", p.Action)
				}
				if p.From.CmdType != "select" {
					t.Fatalf("expected from cmd select, got %s", p.From.CmdType)
				}
				if p.To.CmdType != "update" {
					t.Fatalf("expected to cmd update, got %s", p.To.CmdType)
				}
			},
		},
		{
			name: "policy permissive/restrictive changed",
			fromSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 AS PERMISSIVE USING (owner_id = 1);`,
			toSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 AS RESTRICTIVE USING (owner_id = 1);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Policies) != 1 {
					t.Fatalf("expected 1 policy diff, got %d", len(e.Policies))
				}
				p := e.Policies[0]
				if p.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", p.Action)
				}
				if !p.From.Permissive {
					t.Fatalf("expected from permissive=true")
				}
				if p.To.Permissive {
					t.Fatalf("expected to permissive=false")
				}
			},
		},
		{
			name: "policy roles changed",
			fromSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 TO PUBLIC USING (owner_id = 1);`,
			toSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 TO current_user USING (owner_id = 1);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Policies) != 1 {
					t.Fatalf("expected 1 policy diff, got %d", len(e.Policies))
				}
				p := e.Policies[0]
				if p.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", p.Action)
				}
			},
		},
		{
			name: "policy USING expression changed",
			fromSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 USING (owner_id = 1);`,
			toSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 USING (owner_id = 2);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Policies) != 1 {
					t.Fatalf("expected 1 policy diff, got %d", len(e.Policies))
				}
				p := e.Policies[0]
				if p.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", p.Action)
				}
				if p.From.UsingExpr == p.To.UsingExpr {
					t.Fatalf("expected USING expressions to differ")
				}
			},
		},
		{
			name: "policy WITH CHECK expression changed",
			fromSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 FOR INSERT WITH CHECK (owner_id = 1);`,
			toSQL: `CREATE TABLE t1 (id int, owner_id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
				CREATE POLICY my_pol ON t1 FOR INSERT WITH CHECK (owner_id = 2);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Policies) != 1 {
					t.Fatalf("expected 1 policy diff, got %d", len(e.Policies))
				}
				p := e.Policies[0]
				if p.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", p.Action)
				}
				if p.From.CheckExpr == p.To.CheckExpr {
					t.Fatalf("expected WITH CHECK expressions to differ")
				}
			},
		},
		{
			name:    "table RLS enabled/disabled changed",
			fromSQL: `CREATE TABLE t1 (id int);`,
			toSQL: `CREATE TABLE t1 (id int);
				ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if !e.RLSChanged {
					t.Fatalf("expected RLSChanged to be true")
				}
				if !e.RLSEnabled {
					t.Fatalf("expected RLSEnabled to be true")
				}
			},
		},
		{
			name:    "table FORCE RLS changed",
			fromSQL: `CREATE TABLE t1 (id int);`,
			toSQL: `CREATE TABLE t1 (id int);
				ALTER TABLE t1 FORCE ROW LEVEL SECURITY;`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if !e.RLSChanged {
					t.Fatalf("expected RLSChanged to be true")
				}
				if !e.ForceRLSEnabled {
					t.Fatalf("expected ForceRLSEnabled to be true")
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
