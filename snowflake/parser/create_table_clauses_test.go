package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// The object-level clauses and forms of CREATE TABLE and its variants, and of
// CREATE [MATERIALIZED] VIEW, as their pages write them.

func TestStartsObjectClause(t *testing.T) {
	all := clauseRowAccess | clauseAggregation | clauseJoin | clauseStorageLifecycle | clauseContact
	cases := []struct {
		sql     string
		allowed objectClause
		want    objectClause
	}{
		{"ROW ACCESS POLICY p ON (a)", all, clauseRowAccess},
		{"WITH ROW ACCESS POLICY p ON (a)", all, clauseRowAccess},
		{"AGGREGATION POLICY p", all, clauseAggregation},
		{"WITH AGGREGATION POLICY p", all, clauseAggregation},
		{"JOIN POLICY p", all, clauseJoin},
		{"STORAGE LIFECYCLE POLICY p ON (a)", all, clauseStorageLifecycle},
		{"WITH CONTACT (steward = u)", all, clauseContact},
		// CONTACT needs its WITH.
		{"CONTACT (steward = u)", all, 0},
		// A clause the page does not list is not one.
		{"JOIN POLICY p", clauseRowAccess, 0},
		{"WITH TAG (t = 'v')", all, 0},
		{"ROW_TIMESTAMP = TRUE", all, 0},
	}
	for _, c := range cases {
		p := &Parser{lexer: NewLexer(c.sql), input: c.sql}
		p.advance()
		if got := p.startsObjectClause(c.allowed); got != c.want {
			t.Errorf("startsObjectClause(%q) = %d, want %d", c.sql, got, c.want)
		}
	}
}

func TestCreateTable_ObjectClauses(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (a INT) ROW ACCESS POLICY p ON (a)",
		"CREATE TABLE t (a INT) AGGREGATION POLICY p ENTITY KEY (a)",
		"CREATE TABLE t (a INT) WITH JOIN POLICY p ALLOWED JOIN KEYS (a)",
		"CREATE TABLE t (a INT) STORAGE LIFECYCLE POLICY p ON (a)",
		"CREATE TABLE t (a INT) WITH CONTACT (steward = u, support = v)",
		"CREATE TABLE t (a INT) ENABLE_SCHEMA_EVOLUTION = TRUE MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 ICEBERG_DEFAULT_DDL_COLLATION = 'en' ERROR_LOGGING = FALSE ROW_TIMESTAMP = TRUE",
		"CREATE TABLE t LIKE s COPY = (COMMENTS, TAGS)",
		"CREATE TABLE t CLONE s COPY TAGS",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

func TestCreateTable_CTASClausesBeforeAs(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TABLE t CLUSTER BY (a) COPY = (GRANTS) AS SELECT 1 AS a")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	// COPY = ( ... ) is kept in Options, as CREATE ICEBERG TABLE keeps it.
	if len(stmt.ClusterBy) != 1 || len(stmt.Options) != 1 || stmt.Options[0].Name != "COPY" || stmt.AsSelect == nil {
		t.Errorf("ClusterBy = %v, Options = %v, AsSelect = %v", stmt.ClusterBy, stmt.Options, stmt.AsSelect)
	}
}

func TestCreateTable_UsingTemplate(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE OR REPLACE TABLE t COPY GRANTS USING TEMPLATE (SELECT 1)")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.UsingTemplate == nil || !stmt.CopyGrants {
		t.Errorf("UsingTemplate = %v, CopyGrants = %v", stmt.UsingTemplate, stmt.CopyGrants)
	}
}

func TestCreateTable_FromArchiveAndBackup(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TRANSIENT TABLE t FROM ARCHIVE OF s AS x WHERE x.a > 1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if a := stmt.FromArchive; a == nil || a.Table.Normalize() != "S" || a.Alias.Normalize() != "X" || a.Where == nil {
		t.Errorf("FromArchive = %+v", stmt.FromArchive)
	}
	if _, errs := testParseCreateTable("CREATE TRANSIENT TABLE IF NOT EXISTS t FROM ARCHIVE OF s WHERE 1"); len(errs) > 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
	stmt, errs = testParseCreateTable("CREATE TABLE t FROM BACKUP SET bs IDENTIFIER 'b1'")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if b := stmt.FromBackup; b == nil || b.BackupSet.Normalize() != "BS" || b.ID != "b1" {
		t.Errorf("FromBackup = %+v", stmt.FromBackup)
	}
}

func TestCreateTable_ReadOnlyClone(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE TEMP READ ONLY TABLE t CLONE s")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Temporary || !stmt.ReadOnly || stmt.Clone == nil {
		t.Errorf("Temporary = %v, ReadOnly = %v, Clone = %v", stmt.Temporary, stmt.ReadOnly, stmt.Clone)
	}
}

func TestCreateInteractiveTable(t *testing.T) {
	stmt, errs := testParseCreateTable("CREATE OR REPLACE INTERACTIVE TABLE t (a INT WITH MASKING POLICY m) CLUSTER BY (a) TARGET_LAG = '1 minute' WAREHOUSE = wh INITIALIZATION_WAREHOUSE = wh2 COMMENT = 'c' AS SELECT a FROM s")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Interactive || !stmt.OrReplace || len(stmt.Options) != 3 || stmt.AsSelect == nil {
		t.Errorf("Interactive = %v, OrReplace = %v, Options = %d, AsSelect = %v", stmt.Interactive, stmt.OrReplace, len(stmt.Options), stmt.AsSelect)
	}
	if _, errs := testParseCreateTable("CREATE INTERACTIVE TABLE t CLUSTER BY (a) AS SELECT a FROM s"); len(errs) > 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

// A warehouse setting is an identifier: a quoted name keeps its quotes.
func TestRefreshSetting_QuotedWarehouse(t *testing.T) {
	stmt, errs := testParseCreateTable(`CREATE INTERACTIVE TABLE t CLUSTER BY (a) WAREHOUSE = "my warehouse" INITIALIZATION_WAREHOUSE = wh AS SELECT 1 AS a`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(stmt.Options) != 2 || stmt.Options[0].Words != `"my warehouse"` || stmt.Options[1].Words != "wh" {
		t.Errorf("Options = %+v, want WAREHOUSE \"my warehouse\" and INITIALIZATION_WAREHOUSE wh", stmt.Options)
	}
	sql := `CREATE DYNAMIC TABLE t CLONE s WAREHOUSE = "my warehouse"`
	file, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	dt := file.Stmts[0].(*ast.CreateDynamicTableStmt)
	if len(dt.Options) != 1 || dt.Options[0].Words != `"my warehouse"` {
		t.Errorf("Options = %+v, want WAREHOUSE \"my warehouse\"", dt.Options)
	}
}

func TestCreateView_ClausesAndColumnPolicies(t *testing.T) {
	stmt, errs := testParseCreateView("CREATE VIEW v (a COMMENT 'x' MASKING POLICY m USING (a, b) PROJECTION POLICY pp, b) ROW ACCESS POLICY r ON (a) JOIN POLICY j ALLOWED JOIN KEYS (a) COPY = (TAGS, GRANTS) WITH CONTACT (steward = u) AS SELECT 1 AS a, 2 AS b")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	col := stmt.Columns[0]
	if col.Comment == nil || col.MaskingPolicy == nil || len(col.MaskingUsing) != 2 || col.ProjectionPolicy == nil {
		t.Errorf("column = %+v", col)
	}
	if stmt.RowPolicy == nil || !stmt.CopyGrants {
		t.Errorf("RowPolicy = %v, CopyGrants = %v", stmt.RowPolicy, stmt.CopyGrants)
	}
}

func TestCreateDynamicTable_Clauses(t *testing.T) {
	sql := "CREATE DYNAMIC ICEBERG TABLE t (a INT) TARGET_LAG = '1 minute' WAREHOUSE = wh PARTITION BY (a, BUCKET(4, a), YEAR(a)) AS SELECT a FROM s"
	stmt, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	dt := stmt.Stmts[0].(*ast.CreateDynamicTableStmt)
	if len(dt.PartitionBy) != 3 {
		t.Fatalf("PartitionBy = %d items, want 3", len(dt.PartitionBy))
	}
	if f, ok := dt.PartitionBy[1].(*ast.FuncCallExpr); !ok || f.Name.Name.Name != "BUCKET" || len(f.Args) != 2 {
		t.Errorf("PartitionBy[1] = %+v, want BUCKET(4, a)", dt.PartitionBy[1])
	}

	sql = "CREATE DYNAMIC TABLE t TARGET_LAG = DOWNSTREAM WAREHOUSE = wh START AT (OFFSET => -60) EXECUTE AS USER u USE SECONDARY ROLES r1, r2 WITH CONTACT (steward = u) AS SELECT 1"
	stmt, err = Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	dt = stmt.Stmts[0].(*ast.CreateDynamicTableStmt)
	if lit, ok := dt.StartAt.(*ast.Literal); !ok || lit.Ival != -60 {
		t.Errorf("StartAt = %+v, want the literal -60", dt.StartAt)
	}
	if dt.StartAtKind != "OFFSET" || dt.StartAt == nil || dt.ExecuteAsUser.Normalize() != "U" || len(dt.SecondaryRoles) != 2 {
		t.Errorf("StartAt = %s %v, ExecuteAsUser = %v, SecondaryRoles = %v", dt.StartAtKind, dt.StartAt, dt.ExecuteAsUser, dt.SecondaryRoles)
	}

	for _, sql := range []string{
		"CREATE DYNAMIC TABLE t CLONE s COPY GRANTS TARGET_LAG = '1 minute' WAREHOUSE = wh EXECUTE AS USER u USE SECONDARY ROLES ALL",
		"CREATE DYNAMIC ICEBERG TABLE t FROM BACKUP SET bs IDENTIFIER 'b1'",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

func TestCreateEventAndExternalTable_Clauses(t *testing.T) {
	for _, sql := range []string{
		"CREATE EVENT TABLE t CLONE s AT (OFFSET => -60) COPY GRANTS",
		"CREATE EVENT TABLE t ROW ACCESS POLICY r ON (a) COPY = (GRANTS) WITH CONTACT (steward = u)",
		"CREATE EXTERNAL TABLE t LOCATION = @s FILE_FORMAT = (FORMAT_NAME = 'f') ROW ACCESS POLICY r ON (VALUE) WITH CONTACT (steward = u)",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

func TestCreateTableVariants_Rejected(t *testing.T) {
	for _, sql := range []string{
		// READ ONLY takes only the CLONE form.
		"CREATE TEMP READ ONLY TABLE t (a INT)",
		// An interactive table needs CLUSTER BY and AS.
		"CREATE INTERACTIVE TABLE t (a INT)",
		"CREATE INTERACTIVE TABLE t AS SELECT 1",
		"CREATE INTERACTIVE TABLE t (a INT) AS SELECT 1 AS a",
		// A warehouse is a name, not a string.
		"CREATE INTERACTIVE TABLE t CLUSTER BY (a) WAREHOUSE = 'wh' AS SELECT 1 AS a",
		// TARGET_LAG is an interactive table's, not a table's.
		"CREATE TABLE t (a INT) TARGET_LAG = '1 minute'",
		// A materialized view copies only grants.
		"CREATE MATERIALIZED VIEW v COPY TAGS AS SELECT 1",
		// An event table copies only grants.
		"CREATE EVENT TABLE t COPY = (TAGS)",
		// An external table's row access policy is on VALUE.
		"CREATE EXTERNAL TABLE t LOCATION = @s FILE_FORMAT = (FORMAT_NAME = 'f') ROW ACCESS POLICY r ON (a)",
		// The storage lifecycle policy names its columns.
		"CREATE TABLE t (a INT) STORAGE LIFECYCLE POLICY p",
		// A row access policy names at least one column; a masking policy's USING too.
		"CREATE VIEW v ROW ACCESS POLICY r ON () AS SELECT 1",
		"CREATE VIEW v (a MASKING POLICY m USING (a,)) AS SELECT 1",
		"CREATE TABLE t (a INT MASKING POLICY m USING (a a))",
		// TARGET_LAG is a string, or DOWNSTREAM for a dynamic table only.
		"CREATE DYNAMIC TABLE t CLONE s TARGET_LAG = DRAFT",
		"CREATE INTERACTIVE TABLE t CLUSTER BY (a) TARGET_LAG = DOWNSTREAM AS SELECT 1",
		// START AT takes a string, or a negative number for OFFSET.
		"CREATE DYNAMIC TABLE t TARGET_LAG = DOWNSTREAM WAREHOUSE = wh START AT (STREAM => 123) AS SELECT 1",
		"CREATE DYNAMIC TABLE t TARGET_LAG = DOWNSTREAM WAREHOUSE = wh START AT (OFFSET => 'x') AS SELECT 1",
		"CREATE DYNAMIC TABLE t TARGET_LAG = DOWNSTREAM WAREHOUSE = wh START AT (OFFSET => 60) AS SELECT 1",
		// FROM BACKUP SET takes no modifier; FROM ARCHIVE OF only TRANSIENT and IF NOT EXISTS.
		"CREATE OR REPLACE TABLE t FROM BACKUP SET bs IDENTIFIER 'b1'",
		"CREATE TEMP TABLE t FROM BACKUP SET bs IDENTIFIER 'b1'",
		"CREATE TRANSIENT TABLE t FROM BACKUP SET bs IDENTIFIER 'b1'",
		"CREATE TABLE IF NOT EXISTS t FROM BACKUP SET bs IDENTIFIER 'b1'",
		"CREATE OR REPLACE TABLE t FROM ARCHIVE OF s WHERE 1",
		"CREATE TEMP TABLE t FROM ARCHIVE OF s WHERE 1",
		"CREATE HYBRID TABLE t FROM ARCHIVE OF s WHERE 1",
		// USING TEMPLATE takes only OR REPLACE and COPY GRANTS.
		"CREATE TEMP TABLE t USING TEMPLATE (SELECT 1)",
		"CREATE TABLE IF NOT EXISTS t USING TEMPLATE (SELECT 1)",
		"CREATE TABLE t CLUSTER BY (a) USING TEMPLATE (SELECT 1)",
		"CREATE TABLE t COMMENT = 'c' USING TEMPLATE (SELECT 1)",
		// An interactive table takes only its page's properties.
		"CREATE INTERACTIVE TABLE t (a INT) CLUSTER BY (a) CHANGE_TRACKING = TRUE AS SELECT 1 AS a",
		"CREATE INTERACTIVE TABLE t (a INT) CLUSTER BY (a) ENABLE_SCHEMA_EVOLUTION = TRUE AS SELECT 1 AS a",
		"CREATE INTERACTIVE TABLE t (a INT) CLUSTER BY (a) COPY GRANTS AS SELECT 1 AS a",
		"CREATE INTERACTIVE TABLE t (a INT) CLUSTER BY (a) COMMENT 'c' AS SELECT 1 AS a",
		"CREATE INTERACTIVE TABLE t (a INT) CLUSTER BY (a) WITH CONTACT (steward = u) AS SELECT 1 AS a",
		// An interactive table's column is a name and a data type.
		"CREATE INTERACTIVE TABLE t (a) CLUSTER BY (a) AS SELECT 1",
		"CREATE INTERACTIVE TABLE t (a INT NOT NULL) CLUSTER BY (a) AS SELECT 1",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted, want an error", sql)
		}
	}
}
