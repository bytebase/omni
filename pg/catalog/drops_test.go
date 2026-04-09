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
