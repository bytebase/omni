package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func testParseCreateView(input string) (*ast.CreateViewStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.CreateViewStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a CreateViewStmt"})
	}
	return stmt, result.Errors
}

func testParseCreateMaterializedView(input string) (*ast.CreateMaterializedViewStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.CreateMaterializedViewStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a CreateMaterializedViewStmt"})
	}
	return stmt, result.Errors
}

func testParseAlterView(input string) (*ast.AlterViewStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.AlterViewStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not an AlterViewStmt"})
	}
	return stmt, result.Errors
}

func testParseAlterMaterializedView(input string) (*ast.AlterMaterializedViewStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.AlterMaterializedViewStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not an AlterMaterializedViewStmt"})
	}
	return stmt, result.Errors
}

// ---------------------------------------------------------------------------
// CREATE VIEW tests
// ---------------------------------------------------------------------------

func TestCreateView_Basic(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW v AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt == nil {
		t.Fatal("expected CreateViewStmt, got nil")
	}
	if stmt.Name.Normalize() != "V" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "V")
	}
	if stmt.OrReplace || stmt.Secure || stmt.Recursive || stmt.IfNotExists {
		t.Error("unexpected flags set")
	}
	if stmt.Query == nil {
		t.Error("Query is nil")
	}
}

func TestCreateView_OrReplace(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE OR REPLACE VIEW v AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
}

func TestCreateView_Secure(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE SECURE VIEW v AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Secure {
		t.Error("expected Secure=true")
	}
}

func TestCreateView_Recursive(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE RECURSIVE VIEW v AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Recursive {
		t.Error("expected Recursive=true")
	}
}

func TestCreateView_OrReplaceSecure(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE OR REPLACE SECURE VIEW v AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
	if !stmt.Secure {
		t.Error("expected Secure=true")
	}
}

func TestCreateView_IfNotExists(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW IF NOT EXISTS v AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists=true")
	}
}

func TestCreateView_QualifiedName(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW mydb.myschema.v AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Name.Normalize() != "MYDB.MYSCHEMA.V" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "MYDB.MYSCHEMA.V")
	}
}

func TestCreateView_ColumnList(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW v (col1, col2 COMMENT 'a column') AS SELECT 1, 2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(stmt.Columns))
	}
	if stmt.Columns[0].Name.Normalize() != "COL1" {
		t.Errorf("col[0] = %q, want COL1", stmt.Columns[0].Name.Normalize())
	}
	if stmt.Columns[1].Name.Normalize() != "COL2" {
		t.Errorf("col[1] = %q, want COL2", stmt.Columns[1].Name.Normalize())
	}
	if stmt.Columns[1].Comment == nil || *stmt.Columns[1].Comment != "a column" {
		t.Error("col[1] comment missing or wrong")
	}
}

func TestCreateView_CopyGrants(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW v COPY GRANTS AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.CopyGrants {
		t.Error("expected CopyGrants=true")
	}
}

func TestCreateView_Comment(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW v COMMENT = 'my view' AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Comment == nil || *stmt.Comment != "my view" {
		t.Errorf("comment = %v, want 'my view'", stmt.Comment)
	}
}

func TestCreateView_WithTag(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW v WITH TAG (env = 'prod') AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Tags) != 1 {
		t.Fatalf("tags = %d, want 1", len(stmt.Tags))
	}
	if stmt.Tags[0].Value != "prod" {
		t.Errorf("tag value = %q, want 'prod'", stmt.Tags[0].Value)
	}
}

func TestCreateView_WithRowAccessPolicy(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW v WITH ROW ACCESS POLICY my_policy ON (col1, col2) AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.RowPolicy == nil {
		t.Fatal("RowPolicy is nil")
	}
	if stmt.RowPolicy.PolicyName.Normalize() != "MY_POLICY" {
		t.Errorf("policy = %q, want MY_POLICY", stmt.RowPolicy.PolicyName.Normalize())
	}
	if len(stmt.RowPolicy.Columns) != 2 {
		t.Fatalf("policy cols = %d, want 2", len(stmt.RowPolicy.Columns))
	}
}

func TestCreateView_WithCTE(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW v AS WITH cte AS (SELECT 1) SELECT * FROM cte")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Query == nil {
		t.Error("Query is nil")
	}
}

func TestCreateView_AllModifiers(t *testing.T) {
	sql := `CREATE OR REPLACE SECURE VIEW v
		COMMENT = 'test'
		COPY GRANTS
		WITH TAG (env = 'prod')
		AS SELECT 1`
	stmt, errs := testParseCreateView(sql)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
	if !stmt.Secure {
		t.Error("expected Secure=true")
	}
	if !stmt.CopyGrants {
		t.Error("expected CopyGrants=true")
	}
	if stmt.Comment == nil || *stmt.Comment != "test" {
		t.Error("comment missing or wrong")
	}
	if len(stmt.Tags) != 1 {
		t.Errorf("tags = %d, want 1", len(stmt.Tags))
	}
}

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED VIEW tests
// ---------------------------------------------------------------------------

func TestCreateMaterializedView_Basic(t *testing.T) {
	stmt, errs := testParseCreateMaterializedView("CREATE MATERIALIZED VIEW mv AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt == nil {
		t.Fatal("expected CreateMaterializedViewStmt, got nil")
	}
	if stmt.Name.Normalize() != "MV" {
		t.Errorf("name = %q, want MV", stmt.Name.Normalize())
	}
	if stmt.Query == nil {
		t.Error("Query is nil")
	}
}

func TestCreateMaterializedView_OrReplaceSecure(t *testing.T) {
	stmt, errs := testParseCreateMaterializedView("CREATE OR REPLACE SECURE MATERIALIZED VIEW mv AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace=true")
	}
	if !stmt.Secure {
		t.Error("expected Secure=true")
	}
}

func TestCreateMaterializedView_IfNotExists(t *testing.T) {
	stmt, errs := testParseCreateMaterializedView("CREATE MATERIALIZED VIEW IF NOT EXISTS mv AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists=true")
	}
}

func TestCreateMaterializedView_ClusterBy(t *testing.T) {
	stmt, errs := testParseCreateMaterializedView("CREATE MATERIALIZED VIEW mv CLUSTER BY (col1, col2) AS SELECT col1, col2 FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.ClusterBy) != 2 {
		t.Fatalf("ClusterBy = %d, want 2", len(stmt.ClusterBy))
	}
}

func TestCreateMaterializedView_ClusterByLinear(t *testing.T) {
	stmt, errs := testParseCreateMaterializedView("CREATE MATERIALIZED VIEW mv CLUSTER BY LINEAR (col1) AS SELECT col1 FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Linear {
		t.Error("expected Linear=true")
	}
	if len(stmt.ClusterBy) != 1 {
		t.Fatalf("ClusterBy = %d, want 1", len(stmt.ClusterBy))
	}
}

func TestCreateMaterializedView_Comment(t *testing.T) {
	stmt, errs := testParseCreateMaterializedView("CREATE MATERIALIZED VIEW mv COMMENT = 'mv comment' AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Comment == nil || *stmt.Comment != "mv comment" {
		t.Errorf("comment = %v, want 'mv comment'", stmt.Comment)
	}
}

func TestCreateMaterializedView_ColumnList(t *testing.T) {
	stmt, errs := testParseCreateMaterializedView("CREATE MATERIALIZED VIEW mv (a, b) AS SELECT 1, 2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(stmt.Columns))
	}
}

func TestCreateMaterializedView_WithTag(t *testing.T) {
	stmt, errs := testParseCreateMaterializedView("CREATE MATERIALIZED VIEW mv WITH TAG (k = 'v') AS SELECT 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Tags) != 1 {
		t.Fatalf("tags = %d, want 1", len(stmt.Tags))
	}
}

func TestCreateMaterializedView_RowAccessPolicy(t *testing.T) {
	stmt, errs := testParseCreateMaterializedView("CREATE MATERIALIZED VIEW mv WITH ROW ACCESS POLICY pol ON (id) AS SELECT id FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.RowPolicy == nil {
		t.Fatal("RowPolicy is nil")
	}
	if stmt.RowPolicy.PolicyName.Normalize() != "POL" {
		t.Errorf("policy = %q, want POL", stmt.RowPolicy.PolicyName.Normalize())
	}
}

// ---------------------------------------------------------------------------
// ALTER VIEW tests
// ---------------------------------------------------------------------------

func TestAlterView_RenameTo(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v RENAME TO v2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewRename {
		t.Errorf("action = %v, want AlterViewRename", stmt.Action)
	}
	if stmt.NewName.Normalize() != "V2" {
		t.Errorf("new name = %q, want V2", stmt.NewName.Normalize())
	}
}

func TestAlterView_IfExists(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW IF EXISTS v RENAME TO v2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.IfExists {
		t.Error("expected IfExists=true")
	}
}

func TestAlterView_SetComment(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v SET COMMENT = 'hello'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewSetComment {
		t.Errorf("action = %v, want AlterViewSetComment", stmt.Action)
	}
	if stmt.Comment == nil || *stmt.Comment != "hello" {
		t.Errorf("comment = %v, want 'hello'", stmt.Comment)
	}
}

func TestAlterView_UnsetComment(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v UNSET COMMENT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewUnsetComment {
		t.Errorf("action = %v, want AlterViewUnsetComment", stmt.Action)
	}
}

func TestAlterView_SetSecure(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v SET SECURE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewSetSecure {
		t.Errorf("action = %v, want AlterViewSetSecure", stmt.Action)
	}
}

func TestAlterView_UnsetSecure(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v UNSET SECURE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewUnsetSecure {
		t.Errorf("action = %v, want AlterViewUnsetSecure", stmt.Action)
	}
}

func TestAlterView_SetTag(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v SET TAG (env = 'prod')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewSetTag {
		t.Errorf("action = %v, want AlterViewSetTag", stmt.Action)
	}
	if len(stmt.Tags) != 1 {
		t.Fatalf("tags = %d, want 1", len(stmt.Tags))
	}
}

func TestAlterView_UnsetTag(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v UNSET TAG (env)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewUnsetTag {
		t.Errorf("action = %v, want AlterViewUnsetTag", stmt.Action)
	}
	if len(stmt.UnsetTags) != 1 {
		t.Fatalf("unset tags = %d, want 1", len(stmt.UnsetTags))
	}
}

func TestAlterView_AddRowAccessPolicy(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v ADD ROW ACCESS POLICY pol ON (col1, col2)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewAddRowAccessPolicy {
		t.Errorf("action = %v, want AlterViewAddRowAccessPolicy", stmt.Action)
	}
	if stmt.PolicyName.Normalize() != "POL" {
		t.Errorf("policy = %q, want POL", stmt.PolicyName.Normalize())
	}
	if len(stmt.PolicyCols) != 2 {
		t.Fatalf("policy cols = %d, want 2", len(stmt.PolicyCols))
	}
}

func TestAlterView_DropRowAccessPolicy(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v DROP ROW ACCESS POLICY pol")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewDropRowAccessPolicy {
		t.Errorf("action = %v, want AlterViewDropRowAccessPolicy", stmt.Action)
	}
	if stmt.PolicyName.Normalize() != "POL" {
		t.Errorf("policy = %q, want POL", stmt.PolicyName.Normalize())
	}
}

func TestAlterView_DropAllRowAccessPolicies(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v DROP ALL ROW ACCESS POLICIES")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewDropAllRowAccessPolicies {
		t.Errorf("action = %v, want AlterViewDropAllRowAccessPolicies", stmt.Action)
	}
}

func TestAlterView_AlterColumnSetMaskingPolicy(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v ALTER COLUMN col1 SET MASKING POLICY pol")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewColumnSetMaskingPolicy {
		t.Errorf("action = %v, want AlterViewColumnSetMaskingPolicy", stmt.Action)
	}
	if stmt.Column.Normalize() != "COL1" {
		t.Errorf("column = %q, want COL1", stmt.Column.Normalize())
	}
	if stmt.MaskingPolicy.Normalize() != "POL" {
		t.Errorf("policy = %q, want POL", stmt.MaskingPolicy.Normalize())
	}
}

func TestAlterView_AlterColumnSetMaskingPolicyUsing(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v ALTER COLUMN col1 SET MASKING POLICY pol USING (col1, col2)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewColumnSetMaskingPolicy {
		t.Errorf("action = %v, want AlterViewColumnSetMaskingPolicy", stmt.Action)
	}
	if len(stmt.MaskingUsing) != 2 {
		t.Fatalf("masking using = %d, want 2", len(stmt.MaskingUsing))
	}
}

func TestAlterView_AlterColumnUnsetMaskingPolicy(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v ALTER COLUMN col1 UNSET MASKING POLICY")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewColumnUnsetMaskingPolicy {
		t.Errorf("action = %v, want AlterViewColumnUnsetMaskingPolicy", stmt.Action)
	}
}

func TestAlterView_AlterColumnSetTag(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v ALTER COLUMN col1 SET TAG (env = 'prod')")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewColumnSetTag {
		t.Errorf("action = %v, want AlterViewColumnSetTag", stmt.Action)
	}
	if stmt.Column.Normalize() != "COL1" {
		t.Errorf("column = %q, want COL1", stmt.Column.Normalize())
	}
}

func TestAlterView_AlterColumnUnsetTag(t *testing.T) {
	stmt, errs := testParseAlterView("ALTER VIEW v ALTER COLUMN col1 UNSET TAG (env)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewColumnUnsetTag {
		t.Errorf("action = %v, want AlterViewColumnUnsetTag", stmt.Action)
	}
}

func TestAlterView_ModifyColumn(t *testing.T) {
	// MODIFY is an alias for ALTER per legacy grammar
	stmt, errs := testParseAlterView("ALTER VIEW v MODIFY COLUMN col1 SET MASKING POLICY pol")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterViewColumnSetMaskingPolicy {
		t.Errorf("action = %v, want AlterViewColumnSetMaskingPolicy", stmt.Action)
	}
}

// ---------------------------------------------------------------------------
// ALTER MATERIALIZED VIEW tests
// ---------------------------------------------------------------------------

func TestAlterMaterializedView_RenameTo(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv RENAME TO mv2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVRename {
		t.Errorf("action = %v, want AlterMVRename", stmt.Action)
	}
	if stmt.NewName.Normalize() != "MV2" {
		t.Errorf("new name = %q, want MV2", stmt.NewName.Normalize())
	}
}

func TestAlterMaterializedView_ClusterBy(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv CLUSTER BY (col1, col2)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVClusterBy {
		t.Errorf("action = %v, want AlterMVClusterBy", stmt.Action)
	}
	if len(stmt.ClusterBy) != 2 {
		t.Fatalf("ClusterBy = %d, want 2", len(stmt.ClusterBy))
	}
}

func TestAlterMaterializedView_DropClusteringKey(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv DROP CLUSTERING KEY")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVDropClusteringKey {
		t.Errorf("action = %v, want AlterMVDropClusteringKey", stmt.Action)
	}
}

func TestAlterMaterializedView_Suspend(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv SUSPEND")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVSuspend {
		t.Errorf("action = %v, want AlterMVSuspend", stmt.Action)
	}
}

func TestAlterMaterializedView_Resume(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv RESUME")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVResume {
		t.Errorf("action = %v, want AlterMVResume", stmt.Action)
	}
}

func TestAlterMaterializedView_SuspendRecluster(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv SUSPEND RECLUSTER")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVSuspendRecluster {
		t.Errorf("action = %v, want AlterMVSuspendRecluster", stmt.Action)
	}
}

func TestAlterMaterializedView_ResumeRecluster(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv RESUME RECLUSTER")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVResumeRecluster {
		t.Errorf("action = %v, want AlterMVResumeRecluster", stmt.Action)
	}
}

func TestAlterMaterializedView_SetSecure(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv SET SECURE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVSetSecure {
		t.Errorf("action = %v, want AlterMVSetSecure", stmt.Action)
	}
}

func TestAlterMaterializedView_UnsetSecure(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv UNSET SECURE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVUnsetSecure {
		t.Errorf("action = %v, want AlterMVUnsetSecure", stmt.Action)
	}
}

func TestAlterMaterializedView_SetComment(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv SET COMMENT = 'test mv'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVSetComment {
		t.Errorf("action = %v, want AlterMVSetComment", stmt.Action)
	}
	if stmt.Comment == nil || *stmt.Comment != "test mv" {
		t.Errorf("comment = %v, want 'test mv'", stmt.Comment)
	}
}

func TestAlterMaterializedView_UnsetComment(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv UNSET COMMENT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Action != ast.AlterMVUnsetComment {
		t.Errorf("action = %v, want AlterMVUnsetComment", stmt.Action)
	}
}

func TestAlterMaterializedView_SetSecureAndComment(t *testing.T) {
	stmt, errs := testParseAlterMaterializedView("ALTER MATERIALIZED VIEW mv SET SECURE COMMENT = 'secured mv'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	// SET SECURE with COMMENT is captured as AlterMVSetSecure with Comment populated
	if stmt.Action != ast.AlterMVSetSecure {
		t.Errorf("action = %v, want AlterMVSetSecure", stmt.Action)
	}
	if stmt.Comment == nil || *stmt.Comment != "secured mv" {
		t.Errorf("comment = %v, want 'secured mv'", stmt.Comment)
	}
}

// ---------------------------------------------------------------------------
// Walker smoke test
// ---------------------------------------------------------------------------

func TestViewWalker_CreateView(t *testing.T) {
	result := ParseBestEffort("CREATE VIEW v AS SELECT 1")
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatal("no stmts")
	}

	visited := 0
	ast.Inspect(result.File.Stmts[0], func(n ast.Node) bool {
		if n != nil {
			visited++
		}
		return true
	})
	if visited == 0 {
		t.Error("walker visited 0 nodes")
	}
}
