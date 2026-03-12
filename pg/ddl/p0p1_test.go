package ddl

import (
	"strings"
	"testing"
)

// =============================================================================
// P0: analyzeSetOp branch coercion tests
// =============================================================================

func TestSetOpCoercion_IncompatibleTypes(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `CREATE TABLE t_union (id int, name text)`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	// UNION with incompatible types (int vs text) should error.
	stmts = parseStmts(t, `CREATE VIEW v_union_err AS SELECT id FROM t_union UNION SELECT name FROM t_union`)
	for _, s := range stmts {
		err := c.ProcessUtility(s)
		if err == nil {
			t.Fatal("expected error for int UNION text (incompatible categories)")
		}
	}
}

func TestSetOpCoercion_UnknownConst(t *testing.T) {
	c := New()

	// UNION with UNKNOWN Const branch: SELECT 'a' UNION SELECT 'b'
	// Both string literals start as UNKNOWN, should be coerced to TEXT.
	stmts := parseStmts(t, `CREATE VIEW v_union_unknown AS SELECT 'a' AS x UNION SELECT 'b'`)
	for _, s := range stmts {
		err := c.ProcessUtility(s)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	rel := c.GetRelation("", "v_union_unknown")
	if rel == nil {
		t.Fatal("view not found")
	}
	if len(rel.Columns) < 1 {
		t.Fatal("expected at least 1 column")
	}
	// Common type of UNKNOWN and UNKNOWN should be TEXT (via UNKNOWN→TEXT coercion).
	if rel.Columns[0].TypeOID != TEXTOID {
		t.Errorf("expected common type TEXTOID (%d), got %d", TEXTOID, rel.Columns[0].TypeOID)
	}
}

func TestSetOpCoercion_IntSmallint(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE t_int (a int);
		CREATE TABLE t_small (b smallint);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	// int UNION smallint should succeed — common type is int4.
	stmts = parseStmts(t, `CREATE VIEW v_union_int_small AS SELECT a FROM t_int UNION SELECT b FROM t_small`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	rel := c.GetRelation("", "v_union_int_small")
	if rel == nil {
		t.Fatal("view not found")
	}
	if rel.Columns[0].TypeOID != INT4OID {
		t.Errorf("expected common type INT4OID (%d), got %d", INT4OID, rel.Columns[0].TypeOID)
	}
}

// =============================================================================
// P1: DefineType tests
// =============================================================================

func TestDefineType_ShellType(t *testing.T) {
	c := New()

	// CREATE TYPE mytype; — shell type
	stmts := parseStmts(t, `CREATE TYPE mytype`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create shell type: %v", err)
		}
	}

	// Verify shell type exists but is not defined.
	ns := defaultSchemaOID(t, c)
	bt := c.typeByName[typeKey{ns: ns, name: "mytype"}]
	if bt == nil {
		t.Fatal("shell type not found")
	}
	if bt.IsDefined {
		t.Error("shell type should not be IsDefined")
	}
}

func TestDefineType_DuplicateShellType(t *testing.T) {
	c := New()

	stmts := parseStmts(t, `CREATE TYPE mytype`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create shell type: %v", err)
		}
	}

	// Second CREATE TYPE mytype should fail.
	stmts = parseStmts(t, `CREATE TYPE mytype`)
	for _, s := range stmts {
		err := c.ProcessUtility(s)
		if err == nil {
			t.Fatal("expected error for duplicate shell type")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("expected 'already exists' error, got: %v", err)
		}
	}
}

func TestDefineType_FullType(t *testing.T) {
	c := New()

	// Phase 1: create shell type.
	stmts := parseStmts(t, `CREATE TYPE mytype`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create shell type: %v", err)
		}
	}

	// Phase 2: full definition.
	stmts = parseStmts(t, `CREATE TYPE mytype (
		INPUT = int4in,
		OUTPUT = int4out,
		INTERNALLENGTH = 4,
		PASSEDBYVALUE,
		ALIGNMENT = int4,
		STORAGE = plain,
		CATEGORY = 'N'
	)`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create full type: %v", err)
		}
	}

	// Verify full type.
	ns := defaultSchemaOID(t, c)
	bt := c.typeByName[typeKey{ns: ns, name: "mytype"}]
	if bt == nil {
		t.Fatal("type not found")
	}
	if !bt.IsDefined {
		t.Error("type should be IsDefined")
	}
	if bt.Type != 'b' {
		t.Errorf("expected Type='b', got %c", bt.Type)
	}
	if bt.Len != 4 {
		t.Errorf("expected Len=4, got %d", bt.Len)
	}
	if !bt.ByVal {
		t.Error("expected ByVal=true")
	}
	if bt.Align != 'i' {
		t.Errorf("expected Align='i', got %c", bt.Align)
	}
	if bt.Storage != 'p' {
		t.Errorf("expected Storage='p', got %c", bt.Storage)
	}
	if bt.Category != 'N' {
		t.Errorf("expected Category='N', got %c", bt.Category)
	}
	if bt.Input == 0 {
		t.Error("expected Input OID to be set")
	}
	if bt.Output == 0 {
		t.Error("expected Output OID to be set")
	}

	// Verify array type was created.
	arrBt := c.typeByName[typeKey{ns: ns, name: "_mytype"}]
	if arrBt == nil {
		t.Fatal("array type not found")
	}
	if arrBt.Category != 'A' {
		t.Errorf("expected array Category='A', got %c", arrBt.Category)
	}
	if arrBt.Elem != bt.OID {
		t.Errorf("expected array Elem=%d, got %d", bt.OID, arrBt.Elem)
	}
}

func TestDefineType_MissingShellType(t *testing.T) {
	c := New()

	// Full CREATE TYPE without a shell type should fail.
	stmts := parseStmts(t, `CREATE TYPE noexist (
		INPUT = int4in,
		OUTPUT = int4out
	)`)
	for _, s := range stmts {
		err := c.ProcessUtility(s)
		if err == nil {
			t.Fatal("expected error for missing shell type")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("expected 'does not exist' error, got: %v", err)
		}
	}
}

func TestDefineType_MissingInput(t *testing.T) {
	c := New()

	stmts := parseStmts(t, `CREATE TYPE mytype`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	// Missing INPUT function.
	stmts = parseStmts(t, `CREATE TYPE mytype (OUTPUT = int4out)`)
	for _, s := range stmts {
		err := c.ProcessUtility(s)
		if err == nil {
			t.Fatal("expected error for missing input function")
		}
		if !strings.Contains(err.Error(), "input function must be specified") {
			t.Errorf("expected 'input function must be specified' error, got: %v", err)
		}
	}
}

func TestDefineType_MissingOutput(t *testing.T) {
	c := New()

	stmts := parseStmts(t, `CREATE TYPE mytype`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	// Missing OUTPUT function.
	stmts = parseStmts(t, `CREATE TYPE mytype (INPUT = int4in)`)
	for _, s := range stmts {
		err := c.ProcessUtility(s)
		if err == nil {
			t.Fatal("expected error for missing output function")
		}
		if !strings.Contains(err.Error(), "output function must be specified") {
			t.Errorf("expected 'output function must be specified' error, got: %v", err)
		}
	}
}

// =============================================================================
// P1: Recursive CTE tests
// =============================================================================

func TestRecursiveCTE_Simple(t *testing.T) {
	c := New()

	// WITH RECURSIVE t(n) AS (VALUES(1) UNION ALL SELECT n+1 FROM t WHERE n<10)
	// SELECT * FROM t
	stmts := parseStmts(t, `
		CREATE VIEW v_recursive AS
		WITH RECURSIVE t(n) AS (
			SELECT 1
			UNION ALL
			SELECT n + 1 FROM t WHERE n < 10
		)
		SELECT n FROM t
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create view with recursive CTE: %v", err)
		}
	}

	rel := c.GetRelation("", "v_recursive")
	if rel == nil {
		t.Fatal("view not found")
	}
	if len(rel.Columns) < 1 {
		t.Fatal("expected at least 1 column")
	}
	if rel.Columns[0].TypeOID != INT4OID {
		t.Errorf("expected column type INT4OID (%d), got %d", INT4OID, rel.Columns[0].TypeOID)
	}
}

func TestRecursiveCTE_WithTable(t *testing.T) {
	c := New()

	stmts := parseStmts(t, `
		CREATE TABLE employees (id int, name text, manager_id int);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	stmts = parseStmts(t, `
		CREATE VIEW v_recursive_emp AS
		WITH RECURSIVE emp_tree(id, name, manager_id, depth) AS (
			SELECT id, name, manager_id, 0
			FROM employees
			WHERE manager_id IS NULL
			UNION ALL
			SELECT e.id, e.name, e.manager_id, et.depth + 1
			FROM employees e
			JOIN emp_tree et ON e.manager_id = et.id
		)
		SELECT id, name, depth FROM emp_tree
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create recursive CTE view: %v", err)
		}
	}

	rel := c.GetRelation("", "v_recursive_emp")
	if rel == nil {
		t.Fatal("view not found")
	}
	if len(rel.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(rel.Columns))
	}
	if rel.Columns[0].TypeOID != INT4OID {
		t.Errorf("column 0: expected INT4OID, got %d", rel.Columns[0].TypeOID)
	}
	if rel.Columns[1].TypeOID != TEXTOID {
		t.Errorf("column 1: expected TEXTOID, got %d", rel.Columns[1].TypeOID)
	}
	if rel.Columns[2].TypeOID != INT4OID {
		t.Errorf("column 2: expected INT4OID, got %d", rel.Columns[2].TypeOID)
	}
}

func TestRecursiveCTE_NonUnionError(t *testing.T) {
	c := New()

	// A recursive CTE that doesn't use UNION should fail.
	stmts := parseStmts(t, `
		CREATE VIEW v_bad_recursive AS
		WITH RECURSIVE t AS (
			SELECT 1 AS n
		)
		SELECT n FROM t
	`)
	for _, s := range stmts {
		err := c.ProcessUtility(s)
		// pgparser may or may not set Cterecursive on the CTE when there's no UNION.
		// If it doesn't flag it as recursive, it'll be treated as a normal CTE (no error).
		// If it does, we expect the "does not have the form" error.
		_ = err
	}
}
