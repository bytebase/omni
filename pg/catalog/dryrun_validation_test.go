package catalog

import (
	"errors"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// =============================================================================
// Helpers
// =============================================================================

// parseStmts parses SQL and returns a slice of nodes.
func parseStmts(t *testing.T, sql string) []nodes.Node {
	t.Helper()
	list, err := pgparser.Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v\nSQL: %s", err, sql)
	}
	if list == nil || len(list.Items) == 0 {
		t.Fatalf("parse returned no statements\nSQL: %s", sql)
	}
	out := make([]nodes.Node, len(list.Items))
	for i, item := range list.Items {
		if raw, ok := item.(*nodes.RawStmt); ok {
			out[i] = raw.Stmt
		} else {
			out[i] = item
		}
	}
	return out
}

// execTestSQL sets up a catalog with setupSQL, then runs testSQL via Exec.
func execTestSQL(t *testing.T, setupSQL, testSQL string) []ExecResult {
	t.Helper()
	c := New()
	if setupSQL != "" {
		results, err := c.Exec(setupSQL, nil)
		if err != nil {
			t.Fatalf("setup parse error: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatalf("setup error: %v", r.Error)
			}
		}
	}
	results, err := c.Exec(testSQL, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("test parse error: %v", err)
	}
	return results
}

// dryRunExpectError expects the last non-skipped result to have an error with the given SQLSTATE code.
func dryRunExpectError(t *testing.T, setupSQL, testSQL string, code string, msgContains string) {
	t.Helper()
	results := execTestSQL(t, setupSQL, testSQL)
	if len(results) == 0 {
		t.Fatal("no results")
	}
	last := results[len(results)-1]
	if last.Error == nil {
		t.Fatalf("expected error %s, got nil", code)
	}
	var catErr *Error
	if !errors.As(last.Error, &catErr) {
		t.Fatalf("expected *Error, got %T: %v", last.Error, last.Error)
	}
	if catErr.Code != code {
		t.Errorf("error code: got %s, want %s (msg: %s)", catErr.Code, code, catErr.Message)
	}
	if msgContains != "" && !strings.Contains(strings.ToLower(catErr.Message), strings.ToLower(msgContains)) {
		t.Errorf("error message %q does not contain %q", catErr.Message, msgContains)
	}
}

// dryRunExpectOK expects all results to have no errors.
func dryRunExpectOK(t *testing.T, sql string) {
	t.Helper()
	results := execTestSQL(t, "", sql)
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("stmt %d: unexpected error: %v", i, r.Error)
		}
	}
}

// dryRunExpectWarning expects the last result to have no error and at least one warning.
func dryRunExpectWarning(t *testing.T, setupSQL, testSQL string, code string) {
	t.Helper()
	results := execTestSQL(t, setupSQL, testSQL)
	if len(results) == 0 {
		t.Fatal("no results")
	}
	last := results[len(results)-1]
	if last.Error != nil {
		t.Fatalf("expected no error, got: %v", last.Error)
	}
	found := false
	for _, w := range last.Warnings {
		if w.Code == code {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected warning with code %s, got warnings: %v", code, last.Warnings)
	}
}

// =============================================================================
// Phase 1: CREATE TABLE + ALTER TABLE
// =============================================================================

func TestDryRunValidation_CreateTable(t *testing.T) {
	t.Run("DuplicateTable", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"CREATE TABLE t(id int)",
			CodeDuplicateTable, "already exists")
	})

	t.Run("IfNotExistsWarning", func(t *testing.T) {
		dryRunExpectWarning(t,
			"CREATE TABLE t(id int)",
			"CREATE TABLE IF NOT EXISTS t(id int)",
			CodeWarningSkip)
	})

	t.Run("DuplicateColumn", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TABLE t(id int, id text)",
			CodeDuplicateColumn, "specified more than once")
	})

	t.Run("DefaultDoubleQuotedIdentifierRejected", func(t *testing.T) {
		dryRunExpectError(t, "",
			`CREATE TABLE t (a text DEFAULT "OTHER")`,
			CodeFeatureNotSupported, "cannot use column reference in DEFAULT expression")
	})

	t.Run("DefaultColumnReferenceRejected", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TABLE t (a int, b int DEFAULT a)",
			CodeFeatureNotSupported, "cannot use column reference in DEFAULT expression")
	})

	t.Run("DefaultStringLiteralOK", func(t *testing.T) {
		dryRunExpectOK(t, "CREATE TABLE t (a text DEFAULT 'OTHER')")
	})

	t.Run("DefaultTypeMismatchRejected", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TABLE t (a int DEFAULT true)",
			CodeDatatypeMismatch, "")
	})

	for _, tc := range []struct {
		name        string
		expr        string
		code        string
		msgContains string
	}{
		{
			name:        "DefaultAggregateRejected",
			expr:        "count(*)",
			code:        CodeGroupingError,
			msgContains: "aggregate functions are not allowed in DEFAULT expressions",
		},
		{
			name:        "DefaultWindowFunctionRejected",
			expr:        "row_number() OVER ()",
			code:        CodeWindowingError,
			msgContains: "window functions are not allowed in DEFAULT expressions",
		},
		{
			name:        "DefaultSetReturningFunctionRejected",
			expr:        "generate_series(1, 2)",
			code:        CodeFeatureNotSupported,
			msgContains: "set-returning functions are not allowed in DEFAULT expressions",
		},
		{
			name:        "DefaultGroupingOperationRejected",
			expr:        "GROUPING(1)",
			code:        CodeGroupingError,
			msgContains: "grouping operations are not allowed in DEFAULT expressions",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dryRunExpectError(t, "",
				"CREATE TABLE t (a int DEFAULT "+tc.expr+")",
				tc.code, tc.msgContains)
		})
	}

	t.Run("UndefinedColumnType", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TABLE t(id nosuchtype)",
			CodeUndefinedObject, "")
	})

	t.Run("SetofColumn", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TABLE t(id SETOF int)",
			CodeInvalidTableDefinition, "")
	})

	t.Run("OnCommitNonTemp", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TABLE t(id int) ON COMMIT DROP",
			CodeInvalidTableDefinition, "")
	})

	t.Run("ArrayDimensionLimit", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TABLE t(a int[][][][][][][][])",
			CodeProgramLimitExceeded, "")
	})

	t.Run("PartitionByTooManyColumns", func(t *testing.T) {
		// PARTITION_MAX_KEYS = 32 in PG. Build a CREATE TABLE with 33 columns in PARTITION BY.
		cols := ""
		parts := ""
		for i := 0; i < 33; i++ {
			if i > 0 {
				cols += ", "
				parts += ", "
			}
			cols += "c" + itoa(i) + " int"
			parts += "c" + itoa(i)
		}
		sql := "CREATE TABLE t(" + cols + ") PARTITION BY RANGE (" + parts + ")"
		dryRunExpectError(t, "", sql, CodeProgramLimitExceeded, "")
	})

	t.Run("ListPartitionMultiColumn", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TABLE t(a int, b int) PARTITION BY LIST (a, b)",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("InheritFromPartition", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE parent(id int) PARTITION BY RANGE (id); CREATE TABLE child PARTITION OF parent FOR VALUES FROM (1) TO (10)",
			"CREATE TABLE grandchild(id int) INHERITS (child)",
			CodeWrongObjectType, "")
	})

	t.Run("InheritFromPartitionedTable", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE parent(id int) PARTITION BY RANGE (id)",
			"CREATE TABLE child(id int) INHERITS (parent)",
			CodeWrongObjectType, "")
	})

	t.Run("TempTableInheritsPermanent", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE perm(id int)",
			"CREATE TEMP TABLE child(id int) INHERITS (perm)",
			CodeWrongObjectType, "")
	})

	t.Run("DuplicateParent", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE parent(id int)",
			"CREATE TABLE child() INHERITS (parent, parent)",
			CodeDuplicateObject, "")
	})

	t.Run("InheritColumnTypeMismatch", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE parent(id int)",
			"CREATE TABLE child(id text) INHERITS (parent)",
			CodeDatatypeMismatch, "")
	})
}

func TestDryRunValidation_AlterTable(t *testing.T) {
	t.Run("TableNotExists", func(t *testing.T) {
		dryRunExpectError(t, "",
			"ALTER TABLE nosuch ADD COLUMN x int",
			CodeUndefinedTable, "")
	})

	t.Run("AddColumnDuplicate", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"ALTER TABLE t ADD COLUMN id int",
			CodeDuplicateColumn, "")
	})

	t.Run("AddColumnUndefinedType", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"ALTER TABLE t ADD COLUMN x nosuchtype",
			CodeUndefinedObject, "")
	})

	t.Run("AddColumnDefaultDoubleQuotedIdentifierRejected", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			`ALTER TABLE t ADD COLUMN a text DEFAULT "OTHER"`,
			CodeFeatureNotSupported, "cannot use column reference in DEFAULT expression")
	})

	t.Run("AddColumnDefaultTypeMismatchRejected", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"ALTER TABLE t ADD COLUMN a int DEFAULT true",
			CodeDatatypeMismatch, "")
	})

	for _, tc := range []struct {
		name        string
		expr        string
		code        string
		msgContains string
	}{
		{
			name:        "AddColumnDefaultAggregateRejected",
			expr:        "count(*)",
			code:        CodeGroupingError,
			msgContains: "aggregate functions are not allowed in DEFAULT expressions",
		},
		{
			name:        "AddColumnDefaultWindowFunctionRejected",
			expr:        "row_number() OVER ()",
			code:        CodeWindowingError,
			msgContains: "window functions are not allowed in DEFAULT expressions",
		},
		{
			name:        "AddColumnDefaultSetReturningFunctionRejected",
			expr:        "generate_series(1, 2)",
			code:        CodeFeatureNotSupported,
			msgContains: "set-returning functions are not allowed in DEFAULT expressions",
		},
		{
			name:        "AddColumnDefaultGroupingOperationRejected",
			expr:        "GROUPING(1)",
			code:        CodeGroupingError,
			msgContains: "grouping operations are not allowed in DEFAULT expressions",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dryRunExpectError(t,
				"CREATE TABLE t(id int)",
				"ALTER TABLE t ADD COLUMN a int DEFAULT "+tc.expr,
				tc.code, tc.msgContains)
		})
	}

	t.Run("AddColumnFKToMissing", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"ALTER TABLE t ADD COLUMN x int REFERENCES nosuch(id)",
			CodeUndefinedTable, "")
	})

	t.Run("DropColumnNotExists", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"ALTER TABLE t DROP COLUMN nosuch",
			CodeUndefinedColumn, "")
	})

	t.Run("DropColumnWithDependentView", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int, name text); CREATE VIEW v AS SELECT id, name FROM t",
			"ALTER TABLE t DROP COLUMN name",
			CodeDependentObjects, "")
	})

	t.Run("SetDefaultColumnNotExists", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"ALTER TABLE t ALTER COLUMN nosuch SET DEFAULT 1",
			CodeUndefinedColumn, "")
	})

	t.Run("SetDefaultColumnNotExistsBeforeDefaultExpr", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			`ALTER TABLE t ALTER COLUMN nosuch SET DEFAULT "OTHER"`,
			CodeUndefinedColumn, "")
	})

	t.Run("SetDefaultDoubleQuotedIdentifierRejected", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(a text)",
			`ALTER TABLE t ALTER COLUMN a SET DEFAULT "OTHER"`,
			CodeFeatureNotSupported, "cannot use column reference in DEFAULT expression")
	})

	t.Run("SetDefaultTypeMismatchRejected", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(a int)",
			"ALTER TABLE t ALTER COLUMN a SET DEFAULT true",
			CodeDatatypeMismatch, "")
	})

	for _, tc := range []struct {
		name        string
		expr        string
		code        string
		msgContains string
	}{
		{
			name:        "SetDefaultAggregateRejected",
			expr:        "count(*)",
			code:        CodeGroupingError,
			msgContains: "aggregate functions are not allowed in DEFAULT expressions",
		},
		{
			name:        "SetDefaultWindowFunctionRejected",
			expr:        "row_number() OVER ()",
			code:        CodeWindowingError,
			msgContains: "window functions are not allowed in DEFAULT expressions",
		},
		{
			name:        "SetDefaultSetReturningFunctionRejected",
			expr:        "generate_series(1, 2)",
			code:        CodeFeatureNotSupported,
			msgContains: "set-returning functions are not allowed in DEFAULT expressions",
		},
		{
			name:        "SetDefaultGroupingOperationRejected",
			expr:        "GROUPING(1)",
			code:        CodeGroupingError,
			msgContains: "grouping operations are not allowed in DEFAULT expressions",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dryRunExpectError(t,
				"CREATE TABLE t(a int)",
				"ALTER TABLE t ALTER COLUMN a SET DEFAULT "+tc.expr,
				tc.code, tc.msgContains)
		})
	}

	t.Run("SetDefaultIdentityColumnBeforeDefaultExpr", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int GENERATED ALWAYS AS IDENTITY, name text)",
			"ALTER TABLE t ALTER COLUMN id SET DEFAULT name",
			CodeSyntaxError, "identity column")
	})

	t.Run("SetDefaultGeneratedColumnBeforeDefaultExpr", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(a int, b int GENERATED ALWAYS AS (a + 1) STORED)",
			"ALTER TABLE t ALTER COLUMN b SET DEFAULT a",
			CodeSyntaxError, "generated column")
	})

	t.Run("SetNotNullColumnNotExists", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"ALTER TABLE t ALTER COLUMN nosuch SET NOT NULL",
			CodeUndefinedColumn, "")
	})

	t.Run("DuplicatePrimaryKey", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int PRIMARY KEY)",
			"ALTER TABLE t ADD CONSTRAINT pk2 PRIMARY KEY (id)",
			CodeDuplicatePKey, "")
	})

	t.Run("FKToNonexistentTable", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"ALTER TABLE t ADD CONSTRAINT fk1 FOREIGN KEY (id) REFERENCES nosuch(id)",
			CodeUndefinedTable, "")
	})

	t.Run("FKTypeMismatch", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE ref(id int PRIMARY KEY); CREATE TABLE t(name text)",
			"ALTER TABLE t ADD CONSTRAINT fk1 FOREIGN KEY (name) REFERENCES ref(id)",
			CodeDatatypeMismatch, "")
	})

	t.Run("FKPermanentReferenceTemp", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TEMP TABLE tmp(id int PRIMARY KEY); CREATE TABLE perm(tid int)",
			"ALTER TABLE perm ADD CONSTRAINT fk1 FOREIGN KEY (tid) REFERENCES tmp(id)",
			CodeInvalidFK, "")
	})

	t.Run("RenameColumnNotExists", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"ALTER TABLE t RENAME COLUMN nosuch TO x",
			CodeUndefinedColumn, "")
	})

	t.Run("AlterTypeWithDependentView", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int, name text); CREATE VIEW v AS SELECT name FROM t",
			"ALTER TABLE t ALTER COLUMN name TYPE varchar(100)",
			CodeFeatureNotSupported, "")
	})

	t.Run("AlterTableOnView", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE VIEW v AS SELECT id FROM t",
			"ALTER TABLE v ADD COLUMN x int",
			CodeWrongObjectType, "")
	})
}

// =============================================================================
// Phase 2: CREATE VIEW + CREATE INDEX
// =============================================================================

func TestDryRunValidation_CreateView(t *testing.T) {
	t.Run("SelectInto", func(t *testing.T) {
		// SELECT INTO is rejected in view definitions
		// pgparser may not support this syntax directly in CREATE VIEW context,
		// so we test with IntoClause if needed
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"CREATE VIEW v AS SELECT id INTO newtbl FROM t",
			CodeFeatureNotSupported, "")
	})

	t.Run("UnloggedView", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"CREATE UNLOGGED VIEW v AS SELECT id FROM t",
			CodeFeatureNotSupported, "")
	})

	t.Run("ColumnCountMismatch", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int, name text)",
			"CREATE VIEW v(a, b, c) AS SELECT id, name FROM t",
			CodeDatatypeMismatch, "")
	})

	t.Run("OrReplaceNonView", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"CREATE OR REPLACE VIEW t AS SELECT 1 AS id",
			CodeWrongObjectType, "")
	})

	t.Run("OrReplaceReduceColumns", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int, name text); CREATE VIEW v AS SELECT id, name FROM t",
			"CREATE OR REPLACE VIEW v AS SELECT id FROM t",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("OrReplaceRenameColumn", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE VIEW v AS SELECT id FROM t",
			"CREATE OR REPLACE VIEW v(new_name) AS SELECT id FROM t",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("OrReplaceTypeChange", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int, name text); CREATE VIEW v AS SELECT id, name FROM t",
			"CREATE OR REPLACE VIEW v AS SELECT id, name::int FROM t",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("CheckOptionOnSetOp", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t1(id int); CREATE TABLE t2(id int)",
			"CREATE VIEW v WITH (check_option=cascaded) AS SELECT id FROM t1 UNION SELECT id FROM t2",
			CodeFeatureNotSupported, "")
	})

	t.Run("ModifyingCTE", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"CREATE VIEW v AS WITH x AS (DELETE FROM t RETURNING id) SELECT id FROM x",
			CodeFeatureNotSupported, "")
	})
}

func TestDryRunValidation_CreateIndex(t *testing.T) {
	t.Run("IndexOnView", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE VIEW v AS SELECT id FROM t",
			"CREATE INDEX idx ON v(id)",
			CodeWrongObjectType, "")
	})

	t.Run("UniqueNonBtreeHash", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"CREATE UNIQUE INDEX idx ON t USING gin(id)",
			CodeFeatureNotSupported, "")
	})

	t.Run("HashMultiColumn", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(a int, b int)",
			"CREATE INDEX idx ON t USING hash(a, b)",
			CodeFeatureNotSupported, "")
	})

	t.Run("IncludeNonBtreeGistSpgist", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int, val int)",
			"CREATE INDEX idx ON t USING hash(id) INCLUDE (val)",
			CodeFeatureNotSupported, "")
	})

	t.Run("ColumnNotExists", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"CREATE INDEX idx ON t(nosuch)",
			CodeUndefinedColumn, "")
	})

	t.Run("SystemColumn", func(t *testing.T) {
		// pgddl doesn't track system columns, so ctid appears as undefined column.
		// PG would return 0A000, pgddl returns 42703.
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"CREATE INDEX idx ON t(ctid)",
			CodeUndefinedColumn, "")
	})

	t.Run("UniquePartitionedMissingKey", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(a int, b int) PARTITION BY RANGE (a)",
			"CREATE UNIQUE INDEX idx ON t(b)",
			CodeFeatureNotSupported, "")
	})

	t.Run("DuplicateIndexName", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE INDEX idx ON t(id)",
			"CREATE INDEX idx ON t(id)",
			CodeDuplicateObject, "")
	})

	// NoColumns: parser rejects CREATE INDEX idx ON t() at parse time — not testable via DryRun.
}

// =============================================================================
// Phase 3: Sequence + Type System
// =============================================================================

func TestDryRunValidation_Sequence(t *testing.T) {
	t.Run("IncrementZero", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE SEQUENCE s INCREMENT 0",
			CodeInvalidParameterValue, "")
	})

	t.Run("MinGreaterThanMax", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE SEQUENCE s MINVALUE 100 MAXVALUE 10",
			CodeInvalidParameterValue, "")
	})

	t.Run("StartOutOfRange", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE SEQUENCE s MINVALUE 10 MAXVALUE 20 START 30",
			CodeInvalidParameterValue, "")
	})

	t.Run("CacheZero", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE SEQUENCE s CACHE 0",
			CodeInvalidParameterValue, "")
	})

	t.Run("DuplicateOption", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE SEQUENCE s INCREMENT 1 INCREMENT 2",
			CodeSyntaxError, "")
	})

	t.Run("SmallintBounds", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE SEQUENCE s AS smallint MINVALUE -100000",
			CodeInvalidParameterValue, "")
	})
}

func TestDryRunValidation_Enum(t *testing.T) {
	t.Run("DuplicateTypeName", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TYPE color AS ENUM ('red', 'blue')",
			"CREATE TYPE color AS ENUM ('a', 'b')",
			CodeDuplicateObject, "")
	})

	t.Run("DuplicateLabel", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TYPE color AS ENUM ('red', 'blue', 'red')",
			CodeInvalidParameterValue, "")
	})

	t.Run("AddValueBeforeNonexistent", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TYPE color AS ENUM ('red', 'blue')",
			"ALTER TYPE color ADD VALUE 'green' BEFORE 'nosuch'",
			CodeInvalidParameterValue, "")
	})

	t.Run("AddValueAlreadyExists", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TYPE color AS ENUM ('red', 'blue')",
			"ALTER TYPE color ADD VALUE 'red'",
			CodeDuplicateObject, "")
	})
}

func TestDryRunValidation_Domain(t *testing.T) {
	t.Run("PseudoTypeBase", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE DOMAIN d AS void",
			CodeDatatypeMismatch, "")
	})

	t.Run("NullNotNullConflict", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE DOMAIN d AS int NULL NOT NULL",
			CodeSyntaxError, "")
	})

	t.Run("MultipleDefault", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE DOMAIN d AS int DEFAULT 1 DEFAULT 2",
			CodeSyntaxError, "")
	})

	t.Run("CheckNoInherit", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE DOMAIN d AS int CHECK (VALUE > 0) NO INHERIT",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("UniqueConstraint", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE DOMAIN d AS int CONSTRAINT uniq UNIQUE",
			CodeSyntaxError, "")
	})

	t.Run("CollateOnNonCollatable", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE DOMAIN d AS int COLLATE \"C\"",
			CodeDatatypeMismatch, "")
	})
}

func TestDryRunValidation_Range(t *testing.T) {
	t.Run("MissingSubtype", func(t *testing.T) {
		// Parser rejects CREATE TYPE r AS RANGE () — use a valid but subtypeless def.
		// PG requires SUBTYPE; we pass a non-subtype param to test the check.
		dryRunExpectError(t, "",
			`CREATE TYPE r AS RANGE (COLLATION = "C")`,
			CodeInvalidParameterValue, "")
	})

	t.Run("PseudoTypeSubtype", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TYPE r AS RANGE (SUBTYPE = void)",
			CodeDatatypeMismatch, "")
	})

	t.Run("DuplicateParam", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE TYPE r AS RANGE (SUBTYPE = int4, SUBTYPE = int4)",
			CodeInvalidParameterValue, "")
	})

	t.Run("CollationOnNonCollatable", func(t *testing.T) {
		dryRunExpectError(t, "",
			`CREATE TYPE r AS RANGE (SUBTYPE = int4, COLLATION = "C")`,
			CodeDatatypeMismatch, "")
	})
}

// =============================================================================
// Phase 4: Function/Procedure + Trigger
// =============================================================================

func TestDryRunValidation_Function(t *testing.T) {
	t.Run("SetofParameter", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE FUNCTION f(x SETOF int) RETURNS void AS $$ $$ LANGUAGE sql",
			CodeInvalidFunctionDefinition, "")
	})

	t.Run("VariadicNotLast", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE FUNCTION f(VARIADIC x int[], y int) RETURNS void AS $$ $$ LANGUAGE sql",
			CodeInvalidFunctionDefinition, "")
	})

	t.Run("VariadicNotArray", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE FUNCTION f(VARIADIC x int) RETURNS void AS $$ $$ LANGUAGE sql",
			CodeInvalidFunctionDefinition, "")
	})

	t.Run("DuplicateParamName", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE FUNCTION f(x int, x int) RETURNS void AS $$ $$ LANGUAGE sql",
			CodeInvalidFunctionDefinition, "")
	})

	t.Run("DefaultParamOrder", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE FUNCTION f(x int DEFAULT 1, y int) RETURNS void AS $$ $$ LANGUAGE sql",
			CodeInvalidFunctionDefinition, "")
	})

	t.Run("ProcedureWindow", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE PROCEDURE p() WINDOW LANGUAGE sql AS $$ $$",
			CodeInvalidFunctionDefinition, "")
	})

	t.Run("CostNegative", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE FUNCTION f() RETURNS void COST -1 AS $$ $$ LANGUAGE sql",
			CodeInvalidParameterValue, "")
	})

	t.Run("LanguageNotExists", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE FUNCTION f() RETURNS void AS $$ $$ LANGUAGE nosuchlang",
			CodeUndefinedObject, "")
	})

	t.Run("NoBody", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE FUNCTION f() RETURNS void LANGUAGE sql",
			CodeInvalidFunctionDefinition, "")
	})

	t.Run("RowsWithoutSetof", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE FUNCTION f() RETURNS int ROWS 100 AS $$ SELECT 1 $$ LANGUAGE sql",
			CodeInvalidParameterValue, "")
	})
}

func TestDryRunValidation_Trigger(t *testing.T) {
	t.Run("InsteadOfOnTable", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE FUNCTION trig() RETURNS trigger AS $$ $$ LANGUAGE plpgsql",
			"CREATE TRIGGER tr INSTEAD OF INSERT ON t FOR EACH ROW EXECUTE FUNCTION trig()",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("RowTruncate", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE FUNCTION trig() RETURNS trigger AS $$ $$ LANGUAGE plpgsql",
			"CREATE TRIGGER tr BEFORE TRUNCATE ON t FOR EACH ROW EXECUTE FUNCTION trig()",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("InsteadOfNotForEachRow", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE VIEW v AS SELECT id FROM t; CREATE FUNCTION trig() RETURNS trigger AS $$ $$ LANGUAGE plpgsql",
			"CREATE TRIGGER tr INSTEAD OF INSERT ON v FOR EACH STATEMENT EXECUTE FUNCTION trig()",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("OldTableOnInsert", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE FUNCTION trig() RETURNS trigger AS $$ $$ LANGUAGE plpgsql",
			"CREATE TRIGGER tr AFTER INSERT ON t REFERENCING OLD TABLE AS old_t FOR EACH STATEMENT EXECUTE FUNCTION trig()",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("NewTableOnDelete", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE FUNCTION trig() RETURNS trigger AS $$ $$ LANGUAGE plpgsql",
			"CREATE TRIGGER tr AFTER DELETE ON t REFERENCING NEW TABLE AS new_t FOR EACH STATEMENT EXECUTE FUNCTION trig()",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("DuplicateTransitionTableName", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE FUNCTION trig() RETURNS trigger AS $$ $$ LANGUAGE plpgsql",
			"CREATE TRIGGER tr AFTER UPDATE ON t REFERENCING OLD TABLE AS x NEW TABLE AS x FOR EACH STATEMENT EXECUTE FUNCTION trig()",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("TriggerFunctionNotExists", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION nosuchfunc()",
			CodeUndefinedFunction, "")
	})

	t.Run("TriggerFunctionWrongReturnType", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE FUNCTION badfunc() RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql",
			"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION badfunc()",
			CodeInvalidObjectDefinition, "")
	})

	t.Run("DuplicateTriggerName", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE FUNCTION trig() RETURNS trigger AS $$ $$ LANGUAGE plpgsql; CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trig()",
			"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trig()",
			CodeDuplicateObject, "")
	})

	t.Run("TruncateOnView", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE VIEW v AS SELECT id FROM t; CREATE FUNCTION trig() RETURNS trigger AS $$ $$ LANGUAGE plpgsql",
			"CREATE TRIGGER tr BEFORE TRUNCATE ON v FOR EACH STATEMENT EXECUTE FUNCTION trig()",
			CodeInvalidObjectDefinition, "")
	})
}

// =============================================================================
// Phase 5: Schema + Drop + Comment
// =============================================================================

func TestDryRunValidation_Schema(t *testing.T) {
	t.Run("PgPrefix", func(t *testing.T) {
		dryRunExpectError(t, "",
			"CREATE SCHEMA pg_test",
			CodeReservedName, "")
	})

	t.Run("DuplicateSchema", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE SCHEMA myschema",
			"CREATE SCHEMA myschema",
			CodeDuplicateSchema, "")
	})

	t.Run("DropNonexistent", func(t *testing.T) {
		dryRunExpectError(t, "",
			"DROP SCHEMA nosuch",
			CodeUndefinedSchema, "")
	})

	t.Run("DropNonEmptyNoCascade", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE SCHEMA myschema; CREATE TABLE myschema.t(id int)",
			"DROP SCHEMA myschema",
			CodeSchemaNotEmpty, "")
	})
}

func TestDryRunValidation_Drop(t *testing.T) {
	t.Run("DropTableNotExists", func(t *testing.T) {
		dryRunExpectError(t, "",
			"DROP TABLE nosuch",
			CodeUndefinedTable, "")
	})

	t.Run("DropTableIfExistsWarning", func(t *testing.T) {
		dryRunExpectWarning(t, "",
			"DROP TABLE IF EXISTS nosuch",
			CodeWarningSkip)
	})

	t.Run("DropTableWithDependentViewNoCascade", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int); CREATE VIEW v AS SELECT id FROM t",
			"DROP TABLE t",
			CodeDependentObjects, "")
	})

	t.Run("DropFunctionNotExists", func(t *testing.T) {
		dryRunExpectError(t, "",
			"DROP FUNCTION nosuch()",
			CodeUndefinedFunction, "")
	})

	t.Run("DropFunctionOnProcedure", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE PROCEDURE p() LANGUAGE sql AS $$ $$",
			"DROP FUNCTION p()",
			CodeWrongObjectType, "")
	})

	t.Run("DropProcedureOnFunction", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE FUNCTION f() RETURNS void AS $$ $$ LANGUAGE sql",
			"DROP PROCEDURE f()",
			CodeWrongObjectType, "")
	})
}

func TestDryRunValidation_Comment(t *testing.T) {
	t.Run("CommentOnMissingTable", func(t *testing.T) {
		dryRunExpectError(t, "",
			"COMMENT ON TABLE nosuch IS 'hello'",
			CodeUndefinedTable, "")
	})

	t.Run("CommentOnMissingColumn", func(t *testing.T) {
		dryRunExpectError(t,
			"CREATE TABLE t(id int)",
			"COMMENT ON COLUMN t.nosuch IS 'hello'",
			CodeUndefinedColumn, "")
	})
}

// itoa is a minimal int-to-string helper for test code.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
