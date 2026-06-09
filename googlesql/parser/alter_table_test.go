package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the ALTER statement family (parser-ddl node). Structural-gate
// assertions; oracle-verified accept/reject in alter_table_oracle_test.go.

// alterOf parses sql and asserts the single statement is an *AlterStmt.
func alterOf(t *testing.T, sql string) *ast.AlterStmt {
	t.Helper()
	n := parseDDL(t, sql)
	a, ok := n.(*ast.AlterStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.AlterStmt", sql, n)
	}
	return a
}

// firstAction returns the single alter action, failing if there is not exactly one.
func firstAction(t *testing.T, a *ast.AlterStmt) *ast.AlterAction {
	t.Helper()
	if len(a.Actions) != 1 {
		t.Fatalf("got %d actions, want 1", len(a.Actions))
	}
	return a.Actions[0]
}

func TestAlterTable_AddColumn(t *testing.T) {
	a := alterOf(t, "ALTER TABLE Singers ADD COLUMN Nickname STRING(100)")
	if a.Object != ast.AlterTable {
		t.Errorf("Object = %v, want TABLE", a.Object)
	}
	if a.Name.String() != "Singers" {
		t.Errorf("Name = %q, want Singers", a.Name.String())
	}
	act := firstAction(t, a)
	if act.Kind != ast.AlterAddColumn {
		t.Fatalf("action kind = %v, want ADD COLUMN", act.Kind)
	}
	if act.Column == nil || act.Column.Name != "Nickname" {
		t.Errorf("ADD COLUMN column = %+v, want Nickname", act.Column)
	}
	if act.Column.Type.Text != "STRING(100)" {
		t.Errorf("column type = %q, want STRING(100)", act.Column.Type.Text)
	}
}

func TestAlterTable_AddColumnIfNotExists(t *testing.T) {
	a := alterOf(t, "ALTER TABLE T ADD COLUMN IF NOT EXISTS Bio STRING(MAX)")
	act := firstAction(t, a)
	if !act.IfNotExists {
		t.Error("ADD COLUMN: IfNotExists = false, want true")
	}
}

func TestAlterTable_DropColumn(t *testing.T) {
	a := alterOf(t, "ALTER TABLE T DROP COLUMN Nickname")
	act := firstAction(t, a)
	if act.Kind != ast.AlterDropColumn || act.ColumnName != "Nickname" {
		t.Errorf("action = %+v, want DROP COLUMN Nickname", act)
	}
}

func TestAlterTable_SpannerAlterColumnType(t *testing.T) {
	// Spanner in-place column redefinition: ALTER COLUMN <name> <type>.
	a := alterOf(t, "ALTER TABLE Singers ALTER COLUMN FirstName STRING(2048)")
	act := firstAction(t, a)
	if act.Kind != ast.AlterSpannerAlterColumn {
		t.Fatalf("action kind = %v, want Spanner ALTER COLUMN", act.Kind)
	}
	if act.ColumnName != "FirstName" || act.NewType.Text != "STRING(2048)" {
		t.Errorf("action = %+v, want FirstName STRING(2048)", act)
	}
}

func TestAlterTable_SpannerAlterColumnNotNull(t *testing.T) {
	a := alterOf(t, "ALTER TABLE T ALTER COLUMN c STRING(10) NOT NULL")
	act := firstAction(t, a)
	if act.Kind != ast.AlterSpannerAlterColumn || !act.NotNull {
		t.Errorf("action = %+v, want Spanner ALTER COLUMN with NOT NULL", act)
	}
}

func TestAlterTable_AlterColumnSetOptions(t *testing.T) {
	a := alterOf(t, "ALTER TABLE Events ALTER COLUMN LastMod SET OPTIONS (allow_commit_timestamp = true)")
	act := firstAction(t, a)
	if act.Kind != ast.AlterAlterColumnOptions {
		t.Fatalf("action kind = %v, want ALTER COLUMN SET OPTIONS", act.Kind)
	}
	if act.ColumnName != "LastMod" || len(act.Options) != 1 {
		t.Errorf("action = %+v, want LastMod with 1 option", act)
	}
}

func TestAlterTable_AlterColumnSetDefault(t *testing.T) {
	a := alterOf(t, "ALTER TABLE Orders ALTER COLUMN Status SET DEFAULT ('new')")
	act := firstAction(t, a)
	if act.Kind != ast.AlterAlterColumnSetDefault || act.Default == nil {
		t.Errorf("action = %+v, want ALTER COLUMN SET DEFAULT", act)
	}
}

func TestAlterTable_AlterColumnDropDefault(t *testing.T) {
	a := alterOf(t, "ALTER TABLE Orders ALTER COLUMN Status DROP DEFAULT")
	act := firstAction(t, a)
	if act.Kind != ast.AlterAlterColumnDropDefault {
		t.Errorf("action kind = %v, want ALTER COLUMN DROP DEFAULT", act.Kind)
	}
}

func TestAlterTable_AddDropConstraint(t *testing.T) {
	a := alterOf(t, "ALTER TABLE Orders ADD CONSTRAINT FK_Cust FOREIGN KEY (CustomerId) REFERENCES Customers (CustomerId)")
	act := firstAction(t, a)
	if act.Kind != ast.AlterAddConstraint {
		t.Fatalf("action kind = %v, want ADD CONSTRAINT", act.Kind)
	}
	if act.ConstraintName != "FK_Cust" || act.Constraint == nil || act.Constraint.Kind != ast.ConstraintForeignKey {
		t.Errorf("action = %+v, want named FK constraint", act)
	}

	a2 := alterOf(t, "ALTER TABLE Orders DROP CONSTRAINT FK_Cust")
	act2 := firstAction(t, a2)
	if act2.Kind != ast.AlterDropConstraint || act2.ConstraintName != "FK_Cust" {
		t.Errorf("action = %+v, want DROP CONSTRAINT FK_Cust", act2)
	}
}

func TestAlterTable_SetOnDelete(t *testing.T) {
	a := alterOf(t, "ALTER TABLE T SET ON DELETE CASCADE")
	act := firstAction(t, a)
	if act.Kind != ast.AlterSetOnDelete || act.OnDelete != ast.FKActionCascade {
		t.Errorf("action = %+v, want SET ON DELETE CASCADE", act)
	}
}

func TestAlterTable_RenameTo(t *testing.T) {
	a := alterOf(t, "ALTER TABLE Singers RENAME TO Artists")
	act := firstAction(t, a)
	if act.Kind != ast.AlterRenameTo || act.RenameTo.String() != "Artists" {
		t.Errorf("action = %+v, want RENAME TO Artists", act)
	}
}

func TestAlterTable_RowDeletionPolicyActions(t *testing.T) {
	add := firstAction(t, alterOf(t, "ALTER TABLE T ADD ROW DELETION POLICY (OLDER_THAN(c, INTERVAL 30 DAY))"))
	if add.Kind != ast.AlterAddRowDeletionPolicy || add.RowDeletion == nil {
		t.Errorf("ADD policy action = %+v", add)
	}
	repl := firstAction(t, alterOf(t, "ALTER TABLE T REPLACE ROW DELETION POLICY (OLDER_THAN(c, INTERVAL 60 DAY))"))
	if repl.Kind != ast.AlterReplaceRowDeletionPolicy || repl.RowDeletion == nil {
		t.Errorf("REPLACE policy action = %+v", repl)
	}
	drop := firstAction(t, alterOf(t, "ALTER TABLE T DROP ROW DELETION POLICY"))
	if drop.Kind != ast.AlterDropRowDeletionPolicy {
		t.Errorf("DROP policy action kind = %v", drop.Kind)
	}
}

// --- BigQuery-only ALTER forms (Spanner-emulator non-authoritative) ---

func TestAlterTable_SetOptions(t *testing.T) {
	a := alterOf(t, "ALTER TABLE mydataset.mytable SET OPTIONS (description = 'd')")
	act := firstAction(t, a)
	if act.Kind != ast.AlterSetOptions || len(act.Options) != 1 {
		t.Errorf("action = %+v, want SET OPTIONS", act)
	}
}

func TestAlterTable_MultiAction(t *testing.T) {
	// BigQuery allows comma-separated actions (Spanner rejects multi-action).
	a := alterOf(t, "ALTER TABLE t ADD COLUMN a INT64, ADD COLUMN b STRING")
	if len(a.Actions) != 2 {
		t.Fatalf("got %d actions, want 2", len(a.Actions))
	}
	if a.Actions[0].Column.Name != "a" || a.Actions[1].Column.Name != "b" {
		t.Errorf("actions = %+v / %+v, want columns a, b", a.Actions[0].Column, a.Actions[1].Column)
	}
}

func TestAlterTable_SetDataType(t *testing.T) {
	a := alterOf(t, "ALTER TABLE t ALTER COLUMN c SET DATA TYPE STRING")
	act := firstAction(t, a)
	if act.Kind != ast.AlterAlterColumnType || act.NewType.Text != "STRING" {
		t.Errorf("action = %+v, want SET DATA TYPE STRING", act)
	}
}

func TestAlterTable_RenameColumn(t *testing.T) {
	a := alterOf(t, "ALTER TABLE t RENAME COLUMN a TO b")
	act := firstAction(t, a)
	if act.Kind != ast.AlterRenameColumn || act.ColumnName != "a" || act.NewName != "b" {
		t.Errorf("action = %+v, want RENAME COLUMN a TO b", act)
	}
}

func TestAlterTable_AlterColumnDropNotNull(t *testing.T) {
	a := alterOf(t, "ALTER TABLE t ALTER COLUMN c DROP NOT NULL")
	act := firstAction(t, a)
	if act.Kind != ast.AlterAlterColumnDropNotNull {
		t.Errorf("action kind = %v, want ALTER COLUMN DROP NOT NULL", act.Kind)
	}
}

func TestAlterTable_AddPrimaryKeyNotEnforced(t *testing.T) {
	a := alterOf(t, "ALTER TABLE t ADD PRIMARY KEY (id) NOT ENFORCED")
	act := firstAction(t, a)
	if act.Kind != ast.AlterAddConstraint || act.Constraint.Kind != ast.ConstraintPrimaryKey {
		t.Fatalf("action = %+v, want ADD PRIMARY KEY", act)
	}
	if act.Constraint.Enforced != "NOT ENFORCED" {
		t.Errorf("PK enforced = %q, want NOT ENFORCED", act.Constraint.Enforced)
	}
}

func TestAlterTable_DropPrimaryKey(t *testing.T) {
	a := alterOf(t, "ALTER TABLE t DROP PRIMARY KEY")
	act := firstAction(t, a)
	if act.Kind != ast.AlterDropPrimaryKey {
		t.Errorf("action kind = %v, want DROP PRIMARY KEY", act.Kind)
	}
}

func TestAlterTable_SetDefaultCollate(t *testing.T) {
	a := alterOf(t, "ALTER TABLE t SET DEFAULT COLLATE 'und:ci'")
	act := firstAction(t, a)
	if act.Kind != ast.AlterSetDefaultCollate {
		t.Errorf("action kind = %v, want SET DEFAULT COLLATE", act.Kind)
	}
}

// --- ALTER SCHEMA / DATABASE / VIEW / INDEX ---

func TestAlterSchema_SetOptions(t *testing.T) {
	a := alterOf(t, "ALTER SCHEMA mydataset SET OPTIONS (default_table_expiration_days = 7)")
	if a.Object != ast.AlterSchema {
		t.Errorf("Object = %v, want SCHEMA", a.Object)
	}
	if firstAction(t, a).Kind != ast.AlterSetOptions {
		t.Error("want SET OPTIONS action")
	}
}

func TestAlterSchema_IfExists(t *testing.T) {
	a := alterOf(t, "ALTER SCHEMA IF EXISTS mydataset SET DEFAULT COLLATE 'und:ci'")
	if !a.IfExists {
		t.Error("IfExists = false, want true")
	}
}

func TestAlterDatabase_SetOptions(t *testing.T) {
	a := alterOf(t, "ALTER DATABASE db SET OPTIONS (version_retention_period = '7d', optimizer_version = 3)")
	if a.Object != ast.AlterDatabase {
		t.Errorf("Object = %v, want DATABASE", a.Object)
	}
	act := firstAction(t, a)
	if act.Kind != ast.AlterSetOptions || len(act.Options) != 2 {
		t.Errorf("action = %+v, want SET OPTIONS with 2 entries", act)
	}
}

func TestAlterView_SetOptions(t *testing.T) {
	a := alterOf(t, "ALTER VIEW mydataset.myview SET OPTIONS (description = 'd')")
	if a.Object != ast.AlterView {
		t.Errorf("Object = %v, want VIEW", a.Object)
	}
}

func TestAlterIndex_AddDropStoredColumn(t *testing.T) {
	add := alterOf(t, "ALTER INDEX SingersByName ADD STORED COLUMN BirthDate")
	if add.Object != ast.AlterIndex {
		t.Errorf("Object = %v, want INDEX", add.Object)
	}
	addAct := firstAction(t, add)
	if addAct.Kind != ast.AlterAddStoredColumn || addAct.StoredColumn != "BirthDate" {
		t.Errorf("action = %+v, want ADD STORED COLUMN BirthDate", addAct)
	}
	drop := firstAction(t, alterOf(t, "ALTER INDEX SingersByName DROP STORED COLUMN BirthDate"))
	if drop.Kind != ast.AlterDropStoredColumn || drop.StoredColumn != "BirthDate" {
		t.Errorf("action = %+v, want DROP STORED COLUMN BirthDate", drop)
	}
}

func TestAlterSearchIndex_AddDropStoredColumn(t *testing.T) {
	// Spanner DDL-049 (ALTER SEARCH INDEX … {ADD|DROP} STORED COLUMN). This is a
	// documented Spanner GoogleSQL form the union parser must accept; the live
	// emulator (image sha256:caf1bd24) non-authoritatively REJECTS it at the
	// syntax layer ("Encountered 'SEARCH' while parsing: alter_statement") — its
	// grammar is a subset that lags the docs (same situation as ALTER VECTOR
	// INDEX REBUILD, divergence #113). Triangulated against the Spanner DDL
	// reference + the legacy GoogleSQLParser.g4 (whose schema_object_kind lists
	// INDEX but not SEARCH INDEX, and whose generic-entity alt matches only
	// IDENTIFIER|PROJECT — SEARCH is a keyword token, so legacy rejects it too).
	// See the bq_ddl_oracle_test.go triangulation fixture for the live verdict.
	add := alterOf(t, "ALTER SEARCH INDEX SingerNameIndex ADD STORED COLUMN FirstName")
	if add.Object != ast.AlterSearchIndex {
		t.Errorf("Object = %v, want SEARCH INDEX", add.Object)
	}
	if add.Name.String() != "SingerNameIndex" {
		t.Errorf("Name = %q, want SingerNameIndex", add.Name.String())
	}
	addAct := firstAction(t, add)
	if addAct.Kind != ast.AlterAddStoredColumn || addAct.StoredColumn != "FirstName" {
		t.Errorf("action = %+v, want ADD STORED COLUMN FirstName", addAct)
	}

	drop := alterOf(t, "ALTER SEARCH INDEX SingerNameIndex DROP STORED COLUMN FirstName")
	if drop.Object != ast.AlterSearchIndex {
		t.Errorf("Object = %v, want SEARCH INDEX", drop.Object)
	}
	dropAct := firstAction(t, drop)
	if dropAct.Kind != ast.AlterDropStoredColumn || dropAct.StoredColumn != "FirstName" {
		t.Errorf("action = %+v, want DROP STORED COLUMN FirstName", dropAct)
	}
}

func TestAlterSearchIndex_IfExists(t *testing.T) {
	// DDL-049 allows IF EXISTS before the index name.
	a := alterOf(t, "ALTER SEARCH INDEX IF EXISTS idx ADD STORED COLUMN c")
	if a.Object != ast.AlterSearchIndex {
		t.Errorf("Object = %v, want SEARCH INDEX", a.Object)
	}
	if !a.IfExists {
		t.Error("IfExists = false, want true")
	}
	if a.Name.String() != "idx" {
		t.Errorf("Name = %q, want idx", a.Name.String())
	}
}

func TestAlterSearchIndex_QualifiedName(t *testing.T) {
	// The index name is a path_expression (may be schema-qualified).
	a := alterOf(t, "ALTER SEARCH INDEX myschema.idx DROP STORED COLUMN c")
	if a.Name.String() != "myschema.idx" {
		t.Errorf("Name = %q, want myschema.idx", a.Name.String())
	}
}

func TestAlterSearchIndex_String(t *testing.T) {
	if got := ast.AlterSearchIndex.String(); got != "SEARCH INDEX" {
		t.Errorf("AlterSearchIndex.String() = %q, want %q", got, "SEARCH INDEX")
	}
}

func TestAlterSearchIndex_Rejects(t *testing.T) {
	cases := []string{
		"ALTER SEARCH",                                // SEARCH without INDEX
		"ALTER SEARCH FOO idx ADD STORED COLUMN c",    // SEARCH not followed by INDEX
		"ALTER SEARCH INDEX",                          // missing name + action
		"ALTER SEARCH INDEX idx",                      // missing action list
		"ALTER SEARCH INDEX idx ADD",                  // dangling ADD
		"ALTER SEARCH INDEX idx ADD STORED",           // STORED without COLUMN
		"ALTER SEARCH INDEX idx ADD STORED COLUMN",    // STORED COLUMN without name
		"ALTER SEARCH INDEX idx DROP STORED",          // DROP STORED without COLUMN
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

func TestAlter_Rejects(t *testing.T) {
	cases := []string{
		"ALTER TABLE",                      // missing name + actions
		"ALTER TABLE T",                    // missing action list
		"ALTER TABLE T ADD",                // dangling ADD
		"ALTER TABLE T DROP COLUMN",        // missing column name
		"ALTER TABLE T RENAME",             // dangling RENAME
		"ALTER TABLE T ALTER COLUMN c SET", // dangling SET
		"ALTER TABLE T SET ON DELETE",      // missing action
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}
