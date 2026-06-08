package catalog

import "testing"

func TestDiffTrigger(t *testing.T) {
	const trigFuncSQL = `CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;`
	const trigFunc2SQL = `CREATE FUNCTION trig_fn2() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;`

	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, d *SchemaDiff)
	}{
		{
			name:    "trigger added on table",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", tr.Action)
				}
				if tr.Name != "my_trig" {
					t.Fatalf("expected my_trig, got %s", tr.Name)
				}
				if tr.From != nil {
					t.Fatalf("expected From to be nil")
				}
				if tr.To == nil {
					t.Fatalf("expected To to be non-nil")
				}
			},
		},
		{
			name: "trigger dropped from table",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", tr.Action)
				}
				if tr.Name != "my_trig" {
					t.Fatalf("expected my_trig, got %s", tr.Name)
				}
				if tr.From == nil {
					t.Fatalf("expected From to be non-nil")
				}
				if tr.To != nil {
					t.Fatalf("expected To to be nil")
				}
			},
		},
		{
			name: "trigger timing changed BEFORE to AFTER",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig AFTER INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", tr.Action)
				}
				if tr.From.Timing != TriggerBefore {
					t.Fatalf("expected From timing BEFORE, got %v", tr.From.Timing)
				}
				if tr.To.Timing != TriggerAfter {
					t.Fatalf("expected To timing AFTER, got %v", tr.To.Timing)
				}
			},
		},
		{
			name: "trigger events changed INSERT to INSERT OR UPDATE",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT OR UPDATE ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", tr.Action)
				}
				if tr.From.Events != TriggerEventInsert {
					t.Fatalf("expected From events INSERT, got %v", tr.From.Events)
				}
				if tr.To.Events != TriggerEventInsert|TriggerEventUpdate {
					t.Fatalf("expected To events INSERT|UPDATE, got %v", tr.To.Events)
				}
			},
		},
		{
			name: "trigger level changed ROW to STATEMENT",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH STATEMENT EXECUTE FUNCTION trig_fn();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", tr.Action)
				}
				if !tr.From.ForEachRow {
					t.Fatalf("expected From.ForEachRow to be true")
				}
				if tr.To.ForEachRow {
					t.Fatalf("expected To.ForEachRow to be false")
				}
			},
		},
		{
			name: "trigger WHEN clause changed",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW WHEN (NEW.id > 0) EXECUTE FUNCTION trig_fn();`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW WHEN (NEW.id > 10) EXECUTE FUNCTION trig_fn();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", tr.Action)
				}
				if tr.From.WhenExpr == tr.To.WhenExpr {
					t.Fatalf("expected WHEN clauses to differ")
				}
			},
		},
		{
			name: "trigger function changed",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + trigFunc2SQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + trigFunc2SQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn2();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", tr.Action)
				}
				if tr.From == nil || tr.To == nil {
					t.Fatalf("expected both From and To to be non-nil")
				}
			},
		},
		{
			name: "trigger enabled state changed",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();
				ALTER TABLE t1 DISABLE TRIGGER my_trig;`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", tr.Action)
				}
				if tr.From.Enabled != 'D' {
					t.Fatalf("expected From.Enabled='D', got %c", tr.From.Enabled)
				}
				if tr.To.Enabled != 'O' {
					t.Fatalf("expected To.Enabled='O', got %c", tr.To.Enabled)
				}
			},
		},
		{
			name: "trigger UPDATE OF columns changed",
			fromSQL: `CREATE TABLE t1 (a int, b int, c int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE UPDATE OF a ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			toSQL: `CREATE TABLE t1 (a int, b int, c int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE UPDATE OF b ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", tr.Action)
				}
				if int16SliceEqual(tr.From.Columns, tr.To.Columns) {
					t.Fatalf("expected Columns to differ")
				}
			},
		},
		{
			name: "trigger transition tables changed",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig AFTER INSERT ON t1 REFERENCING NEW TABLE AS new_rows FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig AFTER INSERT ON t1 REFERENCING NEW TABLE AS new_data FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", tr.Action)
				}
				if tr.From.NewTransitionName != "new_rows" {
					t.Fatalf("expected From.NewTransitionName=new_rows, got %s", tr.From.NewTransitionName)
				}
				if tr.To.NewTransitionName != "new_data" {
					t.Fatalf("expected To.NewTransitionName=new_data, got %s", tr.To.NewTransitionName)
				}
			},
		},
		{
			name: "trigger arguments changed",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn('arg1');`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn('arg2');`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Triggers) != 1 {
					t.Fatalf("expected 1 trigger diff, got %d", len(e.Triggers))
				}
				tr := e.Triggers[0]
				if tr.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", tr.Action)
				}
				if stringSliceEqual(tr.From.Args, tr.To.Args) {
					t.Fatalf("expected Args to differ")
				}
			},
		},
		{
			name: "trigger unchanged produces no entry",
			fromSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			toSQL: `CREATE TABLE t1 (id int);` + trigFuncSQL + `
				CREATE TRIGGER my_trig BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 0 {
					t.Fatalf("expected 0 relation diffs (unchanged), got %d", len(d.Relations))
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
