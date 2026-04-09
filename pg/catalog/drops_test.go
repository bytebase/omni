package catalog

import (
	"strings"
	"testing"
)

func TestDropsFromDiff_NoChanges(t *testing.T) {
	from, err := LoadSQL("CREATE TABLE t (id integer);")
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL("CREATE TABLE t (id integer);")
	if err != nil {
		t.Fatal(err)
	}
	diff := Diff(from, to)
	plan := DropsFromDiff(from, to, diff)
	if plan == nil {
		t.Fatal("expected non-nil plan")
	}
	if len(plan.Ops) != 0 {
		t.Errorf("expected 0 ops for unchanged schemas, got %d: %+v", len(plan.Ops), plan.Ops)
	}
	if plan.SQL() != "" {
		t.Errorf("expected empty SQL for empty plan, got %q", plan.SQL())
	}
}

func TestDropsFromDiff_NilDiff(t *testing.T) {
	plan := DropsFromDiff(nil, nil, nil)
	if plan == nil {
		t.Fatal("expected non-nil plan even for nil diff")
	}
	if len(plan.Ops) != 0 {
		t.Errorf("expected 0 ops for nil diff, got %d", len(plan.Ops))
	}
}

func TestDropsFromDiff_DropSchema(t *testing.T) {
	from, err := LoadSQL("CREATE SCHEMA reporting; CREATE TABLE reporting.r (id integer);")
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL("")
	if err != nil {
		t.Fatal(err)
	}
	diff := Diff(from, to)
	plan := DropsFromDiff(from, to, diff)

	var schemaDrops []MigrationOp
	for _, op := range plan.Ops {
		if op.Type == OpDropSchema {
			schemaDrops = append(schemaDrops, op)
		}
	}
	if len(schemaDrops) != 1 {
		t.Fatalf("expected 1 OpDropSchema, got %d; ops: %+v", len(schemaDrops), plan.Ops)
	}
	op := schemaDrops[0]
	if op.ObjectName != "reporting" {
		t.Errorf("expected ObjectName=reporting, got %q", op.ObjectName)
	}
	if op.SchemaName != "" {
		t.Errorf("expected empty SchemaName for top-level schema op, got %q", op.SchemaName)
	}
	if op.SQL != "" {
		t.Errorf("expected empty SQL, got %q", op.SQL)
	}
	if op.Phase != PhasePre {
		t.Errorf("expected PhasePre, got %v", op.Phase)
	}
	if op.Priority != PrioritySchema {
		t.Errorf("expected PrioritySchema, got %d", op.Priority)
	}
	if op.ObjType != 'n' {
		t.Errorf("expected ObjType='n', got %c", op.ObjType)
	}
	if !op.Transactional {
		t.Errorf("expected Transactional=true")
	}
	if op.ObjOID == 0 {
		t.Errorf("expected non-zero ObjOID")
	}
}

// dropOpExpect is the metadata that every "simple drop" test asserts on the
// resulting MigrationOp. It deliberately omits the SQL and ObjOID fields
// (SQL is always "" by DropsFromDiff contract; ObjOID is verified to be
// nonzero rather than equal to a fixed value, since it depends on catalog
// allocation order).
type dropOpExpect struct {
	opType     MigrationOpType
	schemaName string
	objectName string
	objType    byte
	priority   int
}

// assertSingleDropOp finds the single op in plan matching want.opType and
// want.objectName, then verifies every metadata field that DropsFromDiff
// must populate identically to GenerateMigration: empty SQL, PhasePre,
// Transactional, nonzero ObjOID, plus the per-category fields in want.
func assertSingleDropOp(t *testing.T, plan *MigrationPlan, want dropOpExpect) {
	t.Helper()
	matches := plan.Filter(func(op MigrationOp) bool {
		return op.Type == want.opType && op.ObjectName == want.objectName
	}).Ops
	if len(matches) != 1 {
		t.Fatalf("expected 1 %s for %s, got %d; all ops: %+v", want.opType, want.objectName, len(matches), plan.Ops)
	}
	op := matches[0]
	if op.SchemaName != want.schemaName {
		t.Errorf("SchemaName: got %q, want %q", op.SchemaName, want.schemaName)
	}
	if op.SQL != "" {
		t.Errorf("SQL: got %q, want empty (DropsFromDiff contract)", op.SQL)
	}
	if op.Phase != PhasePre {
		t.Errorf("Phase: got %v, want PhasePre", op.Phase)
	}
	if op.Priority != want.priority {
		t.Errorf("Priority: got %d, want %d", op.Priority, want.priority)
	}
	if op.ObjType != want.objType {
		t.Errorf("ObjType: got %c, want %c", op.ObjType, want.objType)
	}
	if !op.Transactional {
		t.Errorf("Transactional: got false, want true")
	}
	if op.ObjOID == 0 {
		t.Errorf("ObjOID: got 0, want nonzero")
	}
}

// TestDropsFromDiff_SimpleTypeAndSequenceDrops covers the trivial-mirror
// drop categories (enum, domain, range, composite, sequence) where the
// helper just copies fields from a single DiffDrop entry. Each subtest
// loads SDL into from, empty into to, runs DropsFromDiff, and verifies
// every metadata field.
//
// Categories with non-trivial logic (schemas — empty SchemaName; sequence
// identity-skip filter; constraints/triggers with ParentObject; dependent
// view cascade) have their own dedicated tests below.
func TestDropsFromDiff_SimpleTypeAndSequenceDrops(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		want    dropOpExpect
	}{
		{
			name:    "enum",
			fromSQL: "CREATE TYPE color AS ENUM ('red', 'green');",
			want: dropOpExpect{
				opType: OpDropType, schemaName: "public", objectName: "color",
				objType: 't', priority: PriorityType,
			},
		},
		{
			name:    "domain",
			fromSQL: "CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);",
			want: dropOpExpect{
				opType: OpDropType, schemaName: "public", objectName: "positive_int",
				objType: 't', priority: PriorityType,
			},
		},
		{
			name:    "range",
			fromSQL: "CREATE TYPE intrange AS RANGE (subtype = integer);",
			want: dropOpExpect{
				opType: OpDropType, schemaName: "public", objectName: "intrange",
				objType: 't', priority: PriorityType,
			},
		},
		{
			// Composites are stored as relations in PG, so ObjType is 'r'
			// not 't' — easy to get wrong, so it's pinned here.
			name:    "composite",
			fromSQL: "CREATE TYPE point2d AS (x integer, y integer);",
			want: dropOpExpect{
				opType: OpDropType, schemaName: "public", objectName: "point2d",
				objType: 'r', priority: PriorityType,
			},
		},
		{
			name:    "sequence",
			fromSQL: "CREATE SEQUENCE my_seq;",
			want: dropOpExpect{
				opType: OpDropSequence, schemaName: "public", objectName: "my_seq",
				objType: 'S', priority: PrioritySequence,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			from, err := LoadSQL(tc.fromSQL)
			if err != nil {
				t.Fatal(err)
			}
			to, err := LoadSQL("")
			if err != nil {
				t.Fatal(err)
			}
			diff := Diff(from, to)
			plan := DropsFromDiff(from, to, diff)
			assertSingleDropOp(t, plan, tc.want)
		})
	}
}

func TestDropsFromDiff_DropFunction(t *testing.T) {
	t.Run("simple drop emits with signature in ObjectName", func(t *testing.T) {
		from, err := LoadSQL(`CREATE FUNCTION add_one(x integer) RETURNS integer AS $$ SELECT x+1 $$ LANGUAGE sql;`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		// Use plan.Filter to find the single drop.
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropFunction
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropFunction, got %d; ops: %+v", len(drops), plan.Ops)
		}
		op := drops[0]
		// ObjectName must contain both the function name AND the argument type.
		if !strings.Contains(op.ObjectName, "add_one") {
			t.Errorf("expected ObjectName to contain 'add_one', got %q", op.ObjectName)
		}
		if !strings.Contains(op.ObjectName, "integer") {
			t.Errorf("expected ObjectName to contain 'integer' (signature), got %q", op.ObjectName)
		}
		if op.SchemaName != "public" {
			t.Errorf("expected SchemaName=public, got %q", op.SchemaName)
		}
		if op.ObjType != 'f' {
			t.Errorf("expected ObjType='f', got %c", op.ObjType)
		}
		if op.Priority != PriorityFunction {
			t.Errorf("expected PriorityFunction, got %d", op.Priority)
		}
		if op.Phase != PhasePre {
			t.Errorf("expected PhasePre")
		}
		if op.SQL != "" {
			t.Errorf("expected empty SQL, got %q", op.SQL)
		}
		if !op.Transactional {
			t.Errorf("expected Transactional=true")
		}
		if op.ObjOID == 0 {
			t.Errorf("expected non-zero ObjOID")
		}
	})

	t.Run("overloads are distinguishable by signature in ObjectName", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE FUNCTION add_one(x integer) RETURNS integer AS $$ SELECT x+1 $$ LANGUAGE sql;
			CREATE FUNCTION add_one(x bigint) RETURNS bigint AS $$ SELECT x+1 $$ LANGUAGE sql;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropFunction
		}).Ops
		if len(drops) != 2 {
			t.Fatalf("expected 2 OpDropFunction (one per overload), got %d", len(drops))
		}
		// The two ObjectNames must differ — they encode different signatures.
		if drops[0].ObjectName == drops[1].ObjectName {
			t.Errorf("expected distinct ObjectNames for overloads, both got %q", drops[0].ObjectName)
		}
	})

	t.Run("return type change exercises DiffModify+signatureChanged path", func(t *testing.T) {
		// CRITICAL: arg types must stay identical, otherwise funcIdentity
		// (which includes arg types) changes and the diff produces
		// DiffDrop+DiffAdd instead of DiffModify — bypassing the
		// signatureChanged code path entirely. Return type is NOT in
		// funcIdentity but IS in signatureChanged, so changing only the
		// return type is the simplest way to actually exercise the
		// DiffModify branch in dropsForFunctions.
		from, err := LoadSQL(`CREATE FUNCTION foo(x integer) RETURNS integer AS $$ SELECT x $$ LANGUAGE sql;`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE FUNCTION foo(x integer) RETURNS bigint AS $$ SELECT x::bigint $$ LANGUAGE sql;`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)

		// Precondition: confirm we're actually exercising DiffModify, not
		// DiffDrop+DiffAdd. If this fails, the test passes vacuously — see
		// commit log for the original bug.
		modifies := 0
		for _, e := range diff.Functions {
			if e.Action == DiffModify {
				modifies++
			}
		}
		if modifies != 1 {
			t.Fatalf("precondition: expected exactly 1 DiffModify entry for return-type change, got %d; test would pass vacuously. diff.Functions: %+v", modifies, diff.Functions)
		}

		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropFunction
		}).Ops
		if len(drops) != 1 {
			t.Errorf("expected 1 OpDropFunction for return-type change via DiffModify+signatureChanged, got %d: %+v", len(drops), drops)
		}
	})

	t.Run("body-only change emits no drop", func(t *testing.T) {
		from, err := LoadSQL(`CREATE FUNCTION foo(x integer) RETURNS integer AS $$ SELECT x $$ LANGUAGE sql;`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE FUNCTION foo(x integer) RETURNS integer AS $$ SELECT x + 0 $$ LANGUAGE sql;`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)

		// Precondition: confirm we're actually exercising DiffModify, not
		// "no diff at all". Without this, the test passes vacuously if a
		// future refactor stops detecting the body change.
		modifies := 0
		for _, e := range diff.Functions {
			if e.Action == DiffModify {
				modifies++
			}
		}
		if modifies != 1 {
			t.Fatalf("precondition: expected exactly 1 DiffModify entry for body-only change, got %d; test would pass vacuously. diff.Functions: %+v", modifies, diff.Functions)
		}

		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropFunction
		}).Ops
		if len(drops) != 0 {
			t.Errorf("expected 0 OpDropFunction for body-only change (CREATE OR REPLACE handles it), got %d: %+v", len(drops), drops)
		}
	})
}

func TestDropsFromDiff_DropTable(t *testing.T) {
	t.Run("regular table emits OpDropTable", func(t *testing.T) {
		from, err := LoadSQL("CREATE TABLE t (id integer);")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)
		assertSingleDropOp(t, plan, dropOpExpect{
			opType:     OpDropTable,
			schemaName: "public",
			objectName: "t",
			objType:    'r',
			priority:   PriorityTable,
		})
	})

	t.Run("partitioned table emits OpDropTable", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE measurements (id integer, ts timestamp) PARTITION BY RANGE (ts);
		`)
		if err != nil {
			t.Fatal(err)
		}
		// Precondition: confirm the parser tagged this as RelKind 'p'.
		// Without this check the subtest would pass vacuously via the 'r'
		// branch if a future parser refactor changed the behavior.
		if rel := from.GetRelation("public", "measurements"); rel == nil || rel.RelKind != 'p' {
			relKind := byte(0)
			if rel != nil {
				relKind = rel.RelKind
			}
			t.Fatalf("precondition: expected RelKind 'p' for partitioned parent, got %q", relKind)
		}
		to, err := LoadSQL("")
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)
		assertSingleDropOp(t, plan, dropOpExpect{
			opType:     OpDropTable,
			schemaName: "public",
			objectName: "measurements",
			objType:    'r',
			priority:   PriorityTable,
		})
	})

	t.Run("view does not emit OpDropTable", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE VIEW v AS SELECT id FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE TABLE t (id integer);`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)
		// Filter for OpDropTable; expect zero (the view drop will be emitted
		// by dropsForViews in a later task, not here).
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropTable
		}).Ops
		if len(drops) != 0 {
			t.Errorf("expected 0 OpDropTable for a dropped view, got %d: %+v", len(drops), drops)
		}
	})
}

func TestDropsFromDiff_TableRecreate(t *testing.T) {
	t.Run("RelKind change regular to partitioned emits OpDropTable", func(t *testing.T) {
		from, err := LoadSQL(`CREATE TABLE t (id integer, ts timestamp);`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE TABLE t (id integer, ts timestamp) PARTITION BY RANGE (ts);`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		// Precondition: confirm this produces a DiffModify, not DiffDrop+DiffAdd.
		modifies := 0
		for _, e := range diff.Relations {
			if e.Action == DiffModify && e.Name == "t" {
				modifies++
			}
		}
		if modifies != 1 {
			t.Fatalf("precondition: expected 1 DiffModify for RelKind change, got %d", modifies)
		}

		plan := DropsFromDiff(from, to, diff)
		assertSingleDropOp(t, plan, dropOpExpect{
			opType:     OpDropTable,
			schemaName: "public",
			objectName: "t",
			objType:    'r',
			priority:   PriorityTable,
		})
	})

	t.Run("no recreate for column-only modification", func(t *testing.T) {
		from, err := LoadSQL(`CREATE TABLE t (id integer, name text);`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE TABLE t (id integer, name text, age integer);`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropTable
		}).Ops
		if len(drops) != 0 {
			t.Errorf("expected 0 OpDropTable for column-only change, got %d: %+v", len(drops), drops)
		}
	})
}

// ---------------------------------------------------------------------------
// Task 11: dropsForColumns
// ---------------------------------------------------------------------------

func TestDropsFromDiff_DropColumn(t *testing.T) {
	t.Run("drop column emits OpDropColumn with table name in ObjectName", func(t *testing.T) {
		from, err := LoadSQL("CREATE TABLE t (id integer, name text);")
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL("CREATE TABLE t (id integer);")
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)

		// Precondition: this should be a DiffModify on relation "t" with a DiffDrop column.
		var foundColDrop bool
		for _, rel := range diff.Relations {
			if rel.Action == DiffModify && rel.Name == "t" {
				for _, col := range rel.Columns {
					if col.Action == DiffDrop && col.Name == "name" {
						foundColDrop = true
					}
				}
			}
		}
		if !foundColDrop {
			t.Fatal("precondition: expected DiffModify on t with DiffDrop for column 'name'")
		}

		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropColumn
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropColumn, got %d: %+v", len(drops), drops)
		}
		op := drops[0]
		// CRITICAL: ObjectName is the TABLE name, not the column name.
		if op.ObjectName != "t" {
			t.Errorf("ObjectName: got %q, want %q (table name, not column name)", op.ObjectName, "t")
		}
		if op.SchemaName != "public" {
			t.Errorf("SchemaName: got %q, want %q", op.SchemaName, "public")
		}
		if op.SQL != "" {
			t.Errorf("SQL: got %q, want empty", op.SQL)
		}
		if op.Phase != PhasePre {
			t.Errorf("Phase: got %v, want PhasePre", op.Phase)
		}
		if op.ObjType != 'r' {
			t.Errorf("ObjType: got %c, want 'r'", op.ObjType)
		}
		if op.Priority != PriorityColumn {
			t.Errorf("Priority: got %d, want PriorityColumn(%d)", op.Priority, PriorityColumn)
		}
		if !op.Transactional {
			t.Errorf("Transactional: got false, want true")
		}
		if op.ObjOID == 0 {
			t.Errorf("ObjOID: got 0, want nonzero")
		}
	})

	t.Run("view column drop does not emit OpDropColumn", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer, name text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer, name text);
			CREATE VIEW v AS SELECT id FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropColumn
		}).Ops
		if len(drops) != 0 {
			t.Errorf("expected 0 OpDropColumn for view column changes, got %d: %+v", len(drops), drops)
		}
	})
}

// ---------------------------------------------------------------------------
// Task 12: dropsForCheckCascades
// ---------------------------------------------------------------------------

func TestDropsFromDiff_CheckCascade(t *testing.T) {
	t.Run("type change drops referencing CHECK constraint", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (val integer CHECK (val > 0));
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (val text);
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)

		// Precondition: DiffModify on relation with a column type change.
		var foundColModify bool
		for _, rel := range diff.Relations {
			if rel.Action == DiffModify && rel.Name == "t" {
				for _, col := range rel.Columns {
					if col.Action == DiffModify && col.Name == "val" {
						foundColModify = true
					}
				}
			}
		}
		if !foundColModify {
			t.Fatal("precondition: expected DiffModify on t with DiffModify for column 'val'")
		}

		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropConstraint
		}).Ops
		if len(drops) < 1 {
			t.Fatalf("expected at least 1 OpDropConstraint for CHECK cascade, got %d: all ops: %+v", len(drops), plan.Ops)
		}
		// The cascade CHECK drop should have PhaseMain (from outer loop override),
		// NOT PhasePre (which is for direct constraint drops).
		op := drops[0]
		if op.Phase != PhaseMain {
			t.Errorf("Phase: got %v, want PhaseMain (outer loop override)", op.Phase)
		}
		if op.ObjType != 'r' {
			t.Errorf("ObjType: got %c, want 'r' (from outer loop override)", op.ObjType)
		}
		if op.Priority != PriorityColumn {
			t.Errorf("Priority: got %d, want PriorityColumn(%d) (from outer loop override)", op.Priority, PriorityColumn)
		}
		if op.ObjectName != "t" {
			t.Errorf("ObjectName: got %q, want %q (table name)", op.ObjectName, "t")
		}
		if !op.Transactional {
			t.Errorf("Transactional: got false, want true")
		}
	})

	t.Run("no type change emits no check cascade drop", func(t *testing.T) {
		from, err := LoadSQL(`CREATE TABLE t (val integer CHECK (val > 0));`)
		if err != nil {
			t.Fatal(err)
		}
		// Change default only, no type change.
		to, err := LoadSQL(`CREATE TABLE t (val integer DEFAULT 1 CHECK (val > 0));`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropConstraint
		}).Ops
		if len(drops) != 0 {
			t.Errorf("expected 0 OpDropConstraint when no type change, got %d: %+v", len(drops), drops)
		}
	})
}

// ---------------------------------------------------------------------------
// Task 13: dropsForConstraints
// ---------------------------------------------------------------------------

func TestDropsFromDiff_DropConstraint(t *testing.T) {
	t.Run("drop CHECK constraint emits OpDropConstraint", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer, val integer, CONSTRAINT val_positive CHECK (val > 0));
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer, val integer);
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropConstraint
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropConstraint, got %d: %+v", len(drops), plan.Ops)
		}
		op := drops[0]
		if op.ObjectName != "val_positive" {
			t.Errorf("ObjectName: got %q, want %q (constraint name)", op.ObjectName, "val_positive")
		}
		if op.ParentObject != "t" {
			t.Errorf("ParentObject: got %q, want %q (table name)", op.ParentObject, "t")
		}
		if op.SchemaName != "public" {
			t.Errorf("SchemaName: got %q, want %q", op.SchemaName, "public")
		}
		if op.Phase != PhasePre {
			t.Errorf("Phase: got %v, want PhasePre", op.Phase)
		}
		if op.ObjType != 'c' {
			t.Errorf("ObjType: got %c, want 'c'", op.ObjType)
		}
		if op.Priority != PriorityConstraint {
			t.Errorf("Priority: got %d, want PriorityConstraint(%d)", op.Priority, PriorityConstraint)
		}
		if op.SQL != "" {
			t.Errorf("SQL: got %q, want empty", op.SQL)
		}
		// ObjOID should NOT be set (zero) for constraint drops per source.
		if op.ObjOID != 0 {
			t.Errorf("ObjOID: got %d, want 0 (not set by buildDropConstraintOp)", op.ObjOID)
		}
	})

	t.Run("modify CHECK constraint emits OpDropConstraint for old", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (val integer, CONSTRAINT val_check CHECK (val > 0));
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (val integer, CONSTRAINT val_check CHECK (val > 10));
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)

		// Precondition: confirm DiffModify on the constraint.
		var foundConModify bool
		for _, rel := range diff.Relations {
			if rel.Action == DiffModify {
				for _, ce := range rel.Constraints {
					if ce.Action == DiffModify && ce.Name == "val_check" {
						foundConModify = true
					}
				}
			}
		}
		if !foundConModify {
			t.Fatal("precondition: expected DiffModify on constraint 'val_check'")
		}

		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropConstraint
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropConstraint for modified constraint, got %d: %+v", len(drops), plan.Ops)
		}
		if drops[0].ObjectName != "val_check" {
			t.Errorf("ObjectName: got %q, want %q", drops[0].ObjectName, "val_check")
		}
	})
}

// ---------------------------------------------------------------------------
// Tasks 14-15: dropsForViews
// ---------------------------------------------------------------------------

func TestDropsFromDiff_DropView(t *testing.T) {
	t.Run("drop view emits OpDropView", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE VIEW v AS SELECT id FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE TABLE t (id integer);`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropView && op.ObjectName == "v"
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropView for view, got %d: %+v", len(drops), plan.Ops)
		}
		op := drops[0]
		if op.SchemaName != "public" {
			t.Errorf("SchemaName: got %q, want %q", op.SchemaName, "public")
		}
		if op.Phase != PhasePre {
			t.Errorf("Phase: got %v, want PhasePre", op.Phase)
		}
		if op.ObjType != 'r' {
			t.Errorf("ObjType: got %c, want 'r'", op.ObjType)
		}
		if op.Priority != PriorityView {
			t.Errorf("Priority: got %d, want PriorityView(%d)", op.Priority, PriorityView)
		}
		if !op.Transactional {
			t.Errorf("Transactional: got false, want true")
		}
		if op.ObjOID == 0 {
			t.Errorf("ObjOID: got 0, want nonzero")
		}
	})

	t.Run("drop matview emits OpDropView", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE MATERIALIZED VIEW mv AS SELECT id FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`CREATE TABLE t (id integer);`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropView && op.ObjectName == "mv"
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropView for matview, got %d: %+v", len(drops), plan.Ops)
		}
	})

	t.Run("view to matview flip emits OpDropView", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE VIEW v AS SELECT id FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE MATERIALIZED VIEW v AS SELECT id FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)

		// Precondition: confirm DiffModify with relkind flip.
		var foundFlip bool
		for _, rel := range diff.Relations {
			if rel.Action == DiffModify && rel.Name == "v" {
				if rel.From != nil && rel.To != nil && rel.From.RelKind != rel.To.RelKind {
					foundFlip = true
				}
			}
		}
		if !foundFlip {
			t.Fatal("precondition: expected DiffModify with RelKind flip on 'v'")
		}

		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropView && op.ObjectName == "v"
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropView for view→matview flip, got %d: %+v", len(drops), plan.Ops)
		}
	})

	t.Run("matview query modification emits OpDropView", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer, name text);
			CREATE MATERIALIZED VIEW mv AS SELECT id FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer, name text);
			CREATE MATERIALIZED VIEW mv AS SELECT id, name FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)

		// Precondition: DiffModify on the matview.
		var foundMatviewModify bool
		for _, rel := range diff.Relations {
			if rel.Action == DiffModify && rel.Name == "mv" {
				foundMatviewModify = true
			}
		}
		if !foundMatviewModify {
			t.Fatal("precondition: expected DiffModify on matview 'mv'")
		}

		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropView && op.ObjectName == "mv"
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropView for matview query modification, got %d: %+v", len(drops), plan.Ops)
		}
	})

	t.Run("add-only diff emits no OpDropView", func(t *testing.T) {
		from, err := LoadSQL(`CREATE TABLE t (id integer);`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE VIEW v AS SELECT id FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropView
		}).Ops
		if len(drops) != 0 {
			t.Errorf("expected 0 OpDropView for add-only diff, got %d: %+v", len(drops), drops)
		}
	})
}

// ---------------------------------------------------------------------------
// Task 16: dropsForTriggers
// ---------------------------------------------------------------------------

func TestDropsFromDiff_DropTrigger(t *testing.T) {
	t.Run("drop trigger emits OpDropTrigger", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE FUNCTION trg_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END $$ LANGUAGE plpgsql;
			CREATE TRIGGER my_trg BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trg_fn();
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE FUNCTION trg_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END $$ LANGUAGE plpgsql;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropTrigger
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropTrigger, got %d: %+v", len(drops), plan.Ops)
		}
		op := drops[0]
		if op.ObjectName != "my_trg" {
			t.Errorf("ObjectName: got %q, want %q (trigger name)", op.ObjectName, "my_trg")
		}
		if op.ParentObject != "t" {
			t.Errorf("ParentObject: got %q, want %q (table name)", op.ParentObject, "t")
		}
		if op.SchemaName != "public" {
			t.Errorf("SchemaName: got %q, want %q", op.SchemaName, "public")
		}
		if op.Phase != PhasePre {
			t.Errorf("Phase: got %v, want PhasePre", op.Phase)
		}
		if op.ObjType != 'T' {
			t.Errorf("ObjType: got %c, want 'T'", op.ObjType)
		}
		if op.Priority != PriorityTrigger {
			t.Errorf("Priority: got %d, want PriorityTrigger(%d)", op.Priority, PriorityTrigger)
		}
		if !op.Transactional {
			t.Errorf("Transactional: got false, want true")
		}
		if op.ObjOID == 0 {
			t.Errorf("ObjOID: got 0, want nonzero")
		}
	})

	t.Run("modify trigger BEFORE to AFTER emits OpDropTrigger", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE FUNCTION trg_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END $$ LANGUAGE plpgsql;
			CREATE TRIGGER my_trg BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trg_fn();
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE FUNCTION trg_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END $$ LANGUAGE plpgsql;
			CREATE TRIGGER my_trg AFTER INSERT ON t FOR EACH ROW EXECUTE FUNCTION trg_fn();
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)

		// Precondition: DiffModify on trigger.
		var foundTrigModify bool
		for _, rel := range diff.Relations {
			if rel.Action == DiffModify {
				for _, te := range rel.Triggers {
					if te.Action == DiffModify && te.Name == "my_trg" {
						foundTrigModify = true
					}
				}
			}
		}
		if !foundTrigModify {
			t.Fatal("precondition: expected DiffModify on trigger 'my_trg'")
		}

		plan := DropsFromDiff(from, to, diff)
		drops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropTrigger
		}).Ops
		if len(drops) != 1 {
			t.Fatalf("expected 1 OpDropTrigger for modified trigger (BEFORE→AFTER), got %d: %+v", len(drops), plan.Ops)
		}
		if drops[0].ObjectName != "my_trg" {
			t.Errorf("ObjectName: got %q, want %q", drops[0].ObjectName, "my_trg")
		}
	})
}

// ---------------------------------------------------------------------------
// Task 17: dropsForDependentViews
// ---------------------------------------------------------------------------

func TestDropsFromDiff_DependentViewCascade(t *testing.T) {
	t.Run("column type change cascades to dependent view", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer, val integer);
			CREATE VIEW v AS SELECT id, val FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer, val text);
			CREATE VIEW v AS SELECT id, val FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		hasViewDrop := false
		for _, op := range plan.Ops {
			if op.Type == OpDropView && op.ObjectName == "v" {
				hasViewDrop = true
			}
		}
		if !hasViewDrop {
			t.Errorf("expected dependent OpDropView for v cascaded by column type change, got none. ops: %+v", plan.Ops)
		}
	})

	t.Run("transitive view cascade (view depends on view depends on table)", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer, val integer);
			CREATE VIEW v1 AS SELECT id, val FROM t;
			CREATE VIEW v2 AS SELECT id FROM v1;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer, val text);
			CREATE VIEW v1 AS SELECT id, val FROM t;
			CREATE VIEW v2 AS SELECT id FROM v1;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		dropped := map[string]bool{}
		for _, op := range plan.Ops {
			if op.Type == OpDropView {
				dropped[op.ObjectName] = true
			}
		}
		if !dropped["v1"] {
			t.Errorf("expected OpDropView for v1, got none. ops: %+v", plan.Ops)
		}
		if !dropped["v2"] {
			t.Errorf("expected transitive OpDropView for v2 (depends on v1), got none. ops: %+v", plan.Ops)
		}
	})

	t.Run("column drop cascades to dependent view", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer, name text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer);
			CREATE VIEW v AS SELECT id FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		hasViewDrop := false
		for _, op := range plan.Ops {
			if op.Type == OpDropView && op.ObjectName == "v" {
				hasViewDrop = true
			}
		}
		if !hasViewDrop {
			t.Errorf("expected dependent OpDropView for v cascaded by column drop, got none. ops: %+v", plan.Ops)
		}
	})

	t.Run("no cascade when columns unchanged", func(t *testing.T) {
		from, err := LoadSQL(`
			CREATE TABLE t (id integer, val integer);
			CREATE VIEW v AS SELECT id, val FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSQL(`
			CREATE TABLE t (id integer, val integer, extra text);
			CREATE VIEW v AS SELECT id, val FROM t;
		`)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := DropsFromDiff(from, to, diff)

		// Adding a column should NOT cascade a view drop.
		viewDrops := plan.Filter(func(op MigrationOp) bool {
			return op.Type == OpDropView && op.ObjectName == "v"
		}).Ops
		if len(viewDrops) != 0 {
			t.Errorf("expected no dependent view drop for column-add-only change, got %d: %+v", len(viewDrops), viewDrops)
		}
	})
}

func TestDropsFromDiff_DropSequence_SkipsIdentity(t *testing.T) {
	from, err := LoadSQL(`
		CREATE TABLE t (id integer GENERATED ALWAYS AS IDENTITY);
		CREATE SEQUENCE my_seq;
	`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE t (id integer GENERATED ALWAYS AS IDENTITY);`)
	if err != nil {
		t.Fatal(err)
	}
	diff := Diff(from, to)
	plan := DropsFromDiff(from, to, diff)

	var seqDrops []MigrationOp
	for _, op := range plan.Ops {
		if op.Type == OpDropSequence {
			seqDrops = append(seqDrops, op)
		}
	}
	if len(seqDrops) != 1 {
		t.Fatalf("expected 1 OpDropSequence (my_seq, NOT the identity-backed sequence), got %d: %+v", len(seqDrops), seqDrops)
	}
	if seqDrops[0].ObjectName != "my_seq" {
		t.Errorf("expected ObjectName=my_seq, got %q", seqDrops[0].ObjectName)
	}
}

// ---------------------------------------------------------------------------
// Task 18: Differential test — DropsFromDiff vs GenerateMigration
// ---------------------------------------------------------------------------

// TestDropsFromDiff_DifferentialAgainstGenerateMigration is the correctness
// gate ensuring all dropsForX helpers produce the same drop ops as the
// canonical GenerateMigration path. It builds a non-trivial fixture
// exercising 8+ drop categories, runs both code paths on the same input,
// normalizes the outputs (zeroing SQL/Warning, re-sorting with sortDropOps),
// and compares field-by-field.
func TestDropsFromDiff_DifferentialAgainstGenerateMigration(t *testing.T) {
	// The "from" catalog has objects across many categories. The "to" catalog
	// drops the entire reporting schema (covering OpDropSchema, OpDropTable,
	// OpDropView, OpDropFunction, OpDropType, OpDropSequence), changes
	// public.t1.val from integer to text (covering OpDropConstraint via CHECK
	// cascade, and OpDropView via dependent-view cascade), modifies the matview
	// query (covering matview OpDropView), and drops the trigger (covering
	// OpDropTrigger).
	from, err := LoadSQL(`
		CREATE SCHEMA reporting;
		CREATE TABLE reporting.users (id integer, name text, val integer);
		ALTER TABLE reporting.users ADD CONSTRAINT users_val_check CHECK (val > 0);
		CREATE TABLE reporting.orders (id integer);
		CREATE VIEW reporting.user_view AS SELECT id, val FROM reporting.users;
		CREATE FUNCTION reporting.add_one(x integer) RETURNS integer AS $$ SELECT x+1 $$ LANGUAGE sql;
		CREATE TYPE reporting.color AS ENUM ('red', 'green');
		CREATE SEQUENCE reporting.seq;
		CREATE TABLE public.t1 (id integer, val integer);
		ALTER TABLE public.t1 ADD CONSTRAINT t1_val_positive CHECK (val > 0);
		CREATE VIEW public.v1 AS SELECT id, val FROM public.t1;
		CREATE MATERIALIZED VIEW public.mv1 AS SELECT id FROM public.t1;
		CREATE FUNCTION public.trig_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END $$ LANGUAGE plpgsql;
		CREATE TRIGGER trg BEFORE INSERT ON public.t1 FOR EACH ROW EXECUTE FUNCTION public.trig_fn();
	`)
	if err != nil {
		t.Fatal(err)
	}

	to, err := LoadSQL(`
		CREATE TABLE public.t1 (id integer, val text);
		CREATE VIEW public.v1 AS SELECT id, val FROM public.t1;
		CREATE MATERIALIZED VIEW public.mv1 AS SELECT id, val FROM public.t1;
		CREATE FUNCTION public.trig_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END $$ LANGUAGE plpgsql;
	`)
	if err != nil {
		t.Fatal(err)
	}

	diff := Diff(from, to)

	// --- GenerateMigration path (canonical) ---
	gmPlan := GenerateMigration(from, to, diff)
	isDrop := func(op MigrationOp) bool {
		return strings.HasPrefix(string(op.Type), "Drop")
	}
	gmDrops := gmPlan.Filter(isDrop).Ops
	// Normalize: zero out SQL and Warning (DropsFromDiff doesn't populate these).
	for i := range gmDrops {
		gmDrops[i].SQL = ""
		gmDrops[i].Warning = ""
	}
	sortDropOps(gmDrops)

	// --- DropsFromDiff path (under test) ---
	dfdPlan := DropsFromDiff(from, to, diff)
	dfdDrops := append([]MigrationOp(nil), dfdPlan.Ops...)
	sortDropOps(dfdDrops)

	// --- Compare ---
	if len(gmDrops) != len(dfdDrops) {
		t.Logf("GenerateMigration drops (%d):", len(gmDrops))
		for i, op := range gmDrops {
			t.Logf("  [%d] %+v", i, op)
		}
		t.Logf("DropsFromDiff drops (%d):", len(dfdDrops))
		for i, op := range dfdDrops {
			t.Logf("  [%d] %+v", i, op)
		}
		t.Fatalf("count mismatch: GenerateMigration produced %d drops, DropsFromDiff produced %d",
			len(gmDrops), len(dfdDrops))
	}

	for i := range gmDrops {
		gm := gmDrops[i]
		dfd := dfdDrops[i]
		if gm != dfd {
			t.Errorf("mismatch at index %d:\n  GM:  %+v\n  DFD: %+v", i, gm, dfd)
		}
	}

	// Sanity: we should have at least 5 different drop categories.
	categories := make(map[MigrationOpType]bool)
	for _, op := range gmDrops {
		categories[op.Type] = true
	}
	if len(categories) < 5 {
		t.Errorf("expected at least 5 drop categories exercised, got %d: %v", len(categories), categories)
	}
	t.Logf("differential test passed: %d drop ops compared across %d categories", len(gmDrops), len(categories))
}
