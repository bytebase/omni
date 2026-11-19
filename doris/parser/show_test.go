package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// SHOW TABLES
// ---------------------------------------------------------------------------

func TestShowTables(t *testing.T) {
	file, errs := parseForTest("SHOW TABLES")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.ShowStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowStmt, got %T", file.Stmts[0])
	}
	if n.Type != "TABLES" {
		t.Errorf("Type = %q, want TABLES", n.Type)
	}
	if n.Full {
		t.Error("Full should be false")
	}
}

func TestShowTablesLike(t *testing.T) {
	file, errs := parseForTest("SHOW TABLES LIKE '%cm%'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "TABLES" {
		t.Errorf("Type = %q, want TABLES", n.Type)
	}
	if n.Like != "%cm%" {
		t.Errorf("Like = %q, want %%cm%%", n.Like)
	}
}

// TestShowLikeNeedsPattern: LIKE takes a pattern, which the end of the input
// does not supply (container-verified reject). Best-effort parsing still
// yields the statement.
func TestShowLikeNeedsPattern(t *testing.T) {
	for _, sql := range []string{"SHOW TABLES LIKE", "SHOW DATABASES LIKE"} {
		_, errs := parseForTest(sql)
		if len(errs) == 0 || errs[0].Position != len(sql) {
			t.Errorf("%q: errs = %v, want one at the end of the input", sql, errs)
		}
		if r := ParseBestEffort(sql); len(r.File.Stmts) != 1 {
			t.Errorf("ParseBestEffort(%q): stmts = %d, want 1", sql, len(r.File.Stmts))
		}
	}
}

func TestShowFullTables(t *testing.T) {
	file, errs := parseForTest("SHOW FULL TABLES")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if !n.Full {
		t.Error("Full should be true")
	}
	if n.Type != "TABLES" {
		t.Errorf("Type = %q, want TABLES", n.Type)
	}
}

func TestShowTablesFromDB(t *testing.T) {
	file, errs := parseForTest("SHOW TABLES FROM mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.From != "mydb" {
		t.Errorf("From = %q, want mydb", n.From)
	}
}

// ---------------------------------------------------------------------------
// SHOW DATABASES
// ---------------------------------------------------------------------------

func TestShowDatabases(t *testing.T) {
	file, errs := parseForTest("SHOW DATABASES")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "DATABASES" {
		t.Errorf("Type = %q, want DATABASES", n.Type)
	}
}

func TestShowDatabasesLike(t *testing.T) {
	file, errs := parseForTest("SHOW DATABASES LIKE 'test%'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Like != "test%" {
		t.Errorf("Like = %q, want test%%", n.Like)
	}
}

// ---------------------------------------------------------------------------
// SHOW COLUMNS
// ---------------------------------------------------------------------------

func TestShowColumns(t *testing.T) {
	file, errs := parseForTest("SHOW COLUMNS FROM t_agg")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "COLUMNS" {
		t.Errorf("Type = %q, want COLUMNS", n.Type)
	}
	if n.Target == nil || len(n.Target.Parts) == 0 || n.Target.Parts[0] != "t_agg" {
		t.Errorf("Target = %v, want t_agg", n.Target)
	}
}

func TestShowFullColumns(t *testing.T) {
	file, errs := parseForTest("SHOW FULL COLUMNS FROM t_agg")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if !n.Full {
		t.Error("Full should be true")
	}
	if n.Type != "COLUMNS" {
		t.Errorf("Type = %q, want COLUMNS", n.Type)
	}
}

// ---------------------------------------------------------------------------
// SHOW CREATE TABLE / VIEW / DATABASE / CATALOG
// ---------------------------------------------------------------------------

func TestShowCreateTable(t *testing.T) {
	file, errs := parseForTest("SHOW CREATE TABLE demo.test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE TABLE" {
		t.Errorf("Type = %q, want CREATE TABLE", n.Type)
	}
	if n.Target == nil {
		t.Fatal("Target is nil")
	}
	if len(n.Target.Parts) != 2 {
		t.Errorf("Target.Parts = %v, want [demo test_table]", n.Target.Parts)
	}
}

func TestShowBriefCreateTable(t *testing.T) {
	file, errs := parseForTest("SHOW BRIEF CREATE TABLE demo.test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "BRIEF CREATE TABLE" {
		t.Errorf("Type = %q, want BRIEF CREATE TABLE", n.Type)
	}
}

func TestShowCreateView(t *testing.T) {
	file, errs := parseForTest("SHOW CREATE VIEW my_view")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE VIEW" {
		t.Errorf("Type = %q, want CREATE VIEW", n.Type)
	}
}

func TestShowCreateDatabase(t *testing.T) {
	file, errs := parseForTest("SHOW CREATE DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE DATABASE" {
		t.Errorf("Type = %q, want CREATE DATABASE", n.Type)
	}
}

func TestShowCreateCatalog(t *testing.T) {
	file, errs := parseForTest("SHOW CREATE CATALOG oracle")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CREATE CATALOG" {
		t.Errorf("Type = %q, want CREATE CATALOG", n.Type)
	}
	if n.Target == nil || n.Target.Parts[0] != "oracle" {
		t.Errorf("Target = %v, want oracle", n.Target)
	}
}

// ---------------------------------------------------------------------------
// SHOW VARIABLES
// ---------------------------------------------------------------------------

func TestShowVariablesLike(t *testing.T) {
	file, errs := parseForTest("SHOW VARIABLES LIKE 'max_connections'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "VARIABLES" {
		t.Errorf("Type = %q, want VARIABLES", n.Type)
	}
	if n.Like != "max_connections" {
		t.Errorf("Like = %q, want max_connections", n.Like)
	}
}

func TestShowVariablesWhere(t *testing.T) {
	file, errs := parseForTest("SHOW VARIABLES WHERE variable_name = 'version'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Where == nil {
		t.Error("Where should not be nil")
	}
}

// ---------------------------------------------------------------------------
// SHOW PARTITIONS
// ---------------------------------------------------------------------------

func TestShowPartitions(t *testing.T) {
	file, errs := parseForTest("SHOW PARTITIONS FROM t_agg")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "PARTITIONS" {
		t.Errorf("Type = %q, want PARTITIONS", n.Type)
	}
	if n.Target == nil || n.Target.Parts[0] != "t_agg" {
		t.Errorf("Target = %v, want t_agg", n.Target)
	}
}

func TestShowTemporaryPartitions(t *testing.T) {
	file, errs := parseForTest("SHOW TEMPORARY PARTITIONS FROM t_temp")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "PARTITIONS" {
		t.Errorf("Type = %q, want PARTITIONS", n.Type)
	}
	if n.Args != "TEMPORARY" {
		t.Errorf("Args = %q, want TEMPORARY", n.Args)
	}
}

func TestShowPartitionsWhere(t *testing.T) {
	file, errs := parseForTest(`SHOW PARTITIONS FROM t_agg WHERE PartitionName = "p2024"`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Target == nil {
		t.Fatal("Target is nil")
	}
	if n.Args == "" {
		t.Error("Args should contain WHERE clause")
	}
}

// ---------------------------------------------------------------------------
// SHOW GRANTS / ROLES
// ---------------------------------------------------------------------------

func TestShowGrants(t *testing.T) {
	file, errs := parseForTest("SHOW GRANTS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "GRANTS" {
		t.Errorf("Type = %q, want GRANTS", n.Type)
	}
}

func TestShowGrantsForUser(t *testing.T) {
	file, errs := parseForTest("SHOW GRANTS FOR jack@'%'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "GRANTS" {
		t.Errorf("Type = %q, want GRANTS", n.Type)
	}
	if n.Args == "" {
		t.Error("Args should contain FOR user")
	}
}

func TestShowAllGrants(t *testing.T) {
	file, errs := parseForTest("SHOW ALL GRANTS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "GRANTS" {
		t.Errorf("Type = %q, want GRANTS", n.Type)
	}
	if n.Args != "ALL" {
		t.Errorf("Args = %q, want ALL", n.Args)
	}
}

func TestShowRoles(t *testing.T) {
	file, errs := parseForTest("SHOW ROLES")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "ROLES" {
		t.Errorf("Type = %q, want ROLES", n.Type)
	}
}

// ---------------------------------------------------------------------------
// SHOW CATALOGS
// ---------------------------------------------------------------------------

func TestShowCatalogs(t *testing.T) {
	file, errs := parseForTest("SHOW CATALOGS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "CATALOGS" {
		t.Errorf("Type = %q, want CATALOGS", n.Type)
	}
}

func TestShowCatalogsLike(t *testing.T) {
	file, errs := parseForTest("SHOW CATALOGS LIKE 'hi%'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Like != "hi%" {
		t.Errorf("Like = %q, want hi%%", n.Like)
	}
}

// ---------------------------------------------------------------------------
// SHOW TABLE STATUS
// ---------------------------------------------------------------------------

func TestShowTableStatus(t *testing.T) {
	file, errs := parseForTest("SHOW TABLE STATUS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "TABLE STATUS" {
		t.Errorf("Type = %q, want TABLE STATUS", n.Type)
	}
}

func TestShowTableStatusFromDB(t *testing.T) {
	file, errs := parseForTest(`SHOW TABLE STATUS FROM db LIKE "%test%"`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "TABLE STATUS" {
		t.Errorf("Type = %q, want TABLE STATUS", n.Type)
	}
	if n.From != "db" {
		t.Errorf("From = %q, want db", n.From)
	}
	if n.Like != "%test%" {
		t.Errorf("Like = %q, want %%test%%", n.Like)
	}
}

// ---------------------------------------------------------------------------
// SHOW ALTER TABLE
// ---------------------------------------------------------------------------

func TestShowAlterTableColumn(t *testing.T) {
	file, errs := parseForTest("SHOW ALTER TABLE COLUMN")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "ALTER TABLE COLUMN" {
		t.Errorf("Type = %q, want ALTER TABLE COLUMN", n.Type)
	}
}

func TestShowAlterTableColumnWhere(t *testing.T) {
	file, errs := parseForTest(`SHOW ALTER TABLE COLUMN WHERE TableName = "table1" ORDER BY CreateTime DESC LIMIT 1`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "ALTER TABLE COLUMN" {
		t.Errorf("Type = %q, want ALTER TABLE COLUMN", n.Type)
	}
	if n.Args == "" {
		t.Error("Args should contain WHERE/ORDER BY/LIMIT clause")
	}
}

func TestShowAlterTableRollup(t *testing.T) {
	file, errs := parseForTest("SHOW ALTER TABLE ROLLUP FROM example_db")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type != "ALTER TABLE ROLLUP" {
		t.Errorf("Type = %q, want ALTER TABLE ROLLUP", n.Type)
	}
	if n.From != "example_db" {
		t.Errorf("From = %q, want example_db", n.From)
	}
}

// ---------------------------------------------------------------------------
// SHOW generic (PROCESSLIST, PROC, WARNINGS, ERRORS, etc.)
// ---------------------------------------------------------------------------

func TestShowProcesslist(t *testing.T) {
	file, errs := parseForTest("SHOW PROCESSLIST")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type == "" {
		t.Error("Type should not be empty")
	}
}

func TestShowWarnings(t *testing.T) {
	file, errs := parseForTest("SHOW WARNINGS")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if n.Type == "" {
		t.Error("Type should not be empty")
	}
}

// ---------------------------------------------------------------------------
// SHOW — Loc validity
// ---------------------------------------------------------------------------

func TestShowLocValid(t *testing.T) {
	input := "SHOW TABLES"
	file, errs := parseForTest(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ShowStmt)
	if !n.Loc.IsValid() {
		t.Error("ShowStmt.Loc should be valid")
	}
	if n.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", n.Loc.Start)
	}
}

func TestShowDatabasesFromCatalog(t *testing.T) {
	// Doris takes FROM alone (DorisParser.g4:294); best effort reads IN too.
	for _, sql := range []string{"SHOW DATABASES FROM hms_catalog", "SHOW DATABASES IN hms_catalog"} {
		r := ParseBestEffort(sql)
		file, errs := r.File, r.Errors
		if len(errs) != 0 {
			t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
		}
		n := file.Stmts[0].(*ast.ShowStmt)
		if n.Type != "DATABASES" || n.From != "hms_catalog" {
			t.Errorf("parseForTest(%q) = Type %q From %q, want DATABASES hms_catalog", sql, n.Type, n.From)
		}
	}
}

// TestShowStrictRejects holds Parse to the SHOW forms of DorisParser.g4 and
// the legacy show_stmt: what neither takes, or a form cut short, is a syntax
// error even where best effort reads it.
func TestShowStrictRejects(t *testing.T) {
	for _, sql := range []string{
		"SHOW",
		"SHOW a",
		"SHOW JOB",
		"SHOW JOB TASK FOR my_job",
		"SHOW STATS mytable",
		"SHOW ALL ANALYZE",
		"SHOW QUEUED ANALYZE",
		"SHOW ANALYZE FOR mydb.t1",
		"SHOW ROUTINE LOAD FROM mydb",
		"SHOW DATABASES IN hms_catalog",
		"SHOW PROCESSLIST zzz",
		"SHOW CREATE TABLE",
		"SHOW DATA PROPERTIES(a = a)",
		"SHOW TABLETS FROM a PARTITIONS a",
		"SHOW COUNT(*)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want a syntax error", sql)
		}
	}
	for _, sql := range []string{
		"SHOW FULL TABLES FROM db LIKE 'a%'",
		"SHOW LOAD WHERE Label COLLATE utf8_general_ci[1] = 'a'",
		"SHOW LOAD WHERE EXTRACT(YEAR FROM DATE CURRENT_DATE) > 0",
		"SHOW LOAD WHERE named_struct('name', Label).name = 'a'",
		"SHOW LOAD WHERE LENGTH(ALL Label) > 0",
		"SHOW LOAD WHERE CreateTime > DATE '2024-01-01'",
		"SHOW LOAD WHERE !(Label = 'job')",
		"SHOW LOAD WHERE Label MATCH_ANY 'job' AND 5 MOD 2 = 1",
		"SHOW LOAD WHERE Label COLLATE utf8_general_ci = 'job'",
		"SHOW DATA SKEW FROM t FOR VERSION AS OF 1",
		"SHOW REPLICA STATUS FROM t TABLESAMPLE (10 PERCENT)",
		"SELECT 1 MATCH_ANY",
		"SHOW DATA ALL FROM db.t ORDER BY a",
		"SHOW DATA PROPERTIES('a' = 'b')",
		"SHOW TABLETS FROM t PARTITION (p1, p2) WHERE a = 1 ORDER BY a LIMIT 1, 2",
		"SHOW ROW POLICY FOR 'u'@'%'",
		"SHOW REPLICA STATUS FROM t PARTITION (p) WHERE a = 'b'",
		"SHOW COUNT(*) WARNINGS",
		"SHOW INVERTED INDEX CHAR_FILTER",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

// TestSetGrammar holds SET to what both of Doris's parsers take;
// testdata/compat.tsv judged the near misses.
func TestSetGrammar(t *testing.T) {
	cases := []struct {
		sql string
		ok  bool
	}{
		{"SET CHARSET a, CHAR SET b, NAMES utf8 COLLATE c, GLOBAL x = DEFAULT, @@LOCAL.y = ON, @u = 1, z = ALL", true},
		{"SET GLOBAL = 1", true},
		{"SET PROPERTY = 1", true},
		{"SET PROPERTY FOR 'u' 'k' = 'v', 'k2' = NULL", true},
		{"SET PROPERTY 'k' = 1", false},
		{"SET LDAP_ADMIN_PASSWORD = PASSWORD('p')", true},
		{"SET NAMES = 1", false},
		{"SET NAMES a DEFAULT", false},
		{"SET a := 1", false},
		{"SET `names` = 1", true},
		{"SET a", false},
		{"SET v AS DEFAULT STORAGE VAULT", true},
	}
	for _, c := range cases {
		_, err := Parse(c.sql)
		if (err == nil) != c.ok {
			t.Errorf("Parse(%q): err = %v, want accepted=%v", c.sql, err, c.ok)
		}
	}
}

// TestSetTypeFromOption: a variable named NAMES or CHARSET keeps the
// VARIABLE type; the NAMES and CHARSET options keep theirs.
func TestSetTypeFromOption(t *testing.T) {
	for sql, want := range map[string]string{
		"SET `names` = 1":  "VARIABLE",
		"SET @charset = 1": "VARIABLE",
		"SET NAMES utf8":   "NAMES",
		"SET CHARSET utf8": "CHARSET",
	} {
		stmt := parseOne(t, sql).(*ast.SetStmt)
		if stmt.Type != want {
			t.Errorf("%s: Type = %q, want %q", sql, stmt.Type, want)
		}
	}
}

// TestParseUseCloudCluster: Doris' useCloudCluster form with no database,
// USE @cluster (DorisParser.g4:930), switches the compute cluster alone.
func TestParseUseCloudCluster(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		catalog string
		db      string
		cluster string
	}{
		{sql: "USE @cluster", cluster: "cluster"},
		{sql: "USE db", db: "db"},
		{sql: "USE cat.db", catalog: "cat", db: "db"},
		{sql: "USE db@cluster", db: "db", cluster: "cluster"},
		{sql: "USE cat.db@cluster", catalog: "cat", db: "db", cluster: "cluster"},
	} {
		file, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		stmt, ok := file.Stmts[0].(*ast.UseStmt)
		if !ok {
			t.Fatalf("Parse(%q): got %T, want *ast.UseStmt", tc.sql, file.Stmts[0])
		}
		if stmt.Catalog != tc.catalog || stmt.Database != tc.db || stmt.Cluster != tc.cluster {
			t.Errorf("Parse(%q): catalog=%q database=%q cluster=%q, want %q/%q/%q",
				tc.sql, stmt.Catalog, stmt.Database, stmt.Cluster, tc.catalog, tc.db, tc.cluster)
		}
	}

	for _, sql := range []string{"USE", "USE @", "USE @"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestShowTableStatusFromQualified: show_param's `'TABLE' 'STATUS' from_or_in
// ident '.' ident` (DorisLegacyParser.g4:1386), which the ANTLR parser's
// multipartIdentifier takes as well (DorisParser.g4:275).
func TestShowTableStatusFromQualified(t *testing.T) {
	stmt := parseOne(t, "SHOW TABLE STATUS FROM a.a").(*ast.ShowStmt)
	if stmt.Type != "TABLE STATUS" {
		t.Errorf("Type = %q, want TABLE STATUS", stmt.Type)
	}
	if stmt.From != "a.a" {
		t.Errorf("From = %q, want a.a", stmt.From)
	}
}

func TestShowTablesFromQualified(t *testing.T) {
	stmt := parseOne(t, "SHOW TABLES FROM a.a").(*ast.ShowStmt)
	if stmt.Type != "TABLES" || stmt.From != "a.a" {
		t.Errorf("Type = %q From = %q, want TABLES from a.a", stmt.Type, stmt.From)
	}
}

// TestShowTableID: SHOW TABLE INTEGER_VALUE names a table by id
// (DorisParser.g4:278, show_param at DorisLegacyParser.g4:1389).
func TestShowTableID(t *testing.T) {
	stmt := parseOne(t, "SHOW TABLE 1").(*ast.ShowStmt)
	if stmt.Type != "TABLE" || stmt.Args != "1" {
		t.Errorf("Type = %q Args = %q, want TABLE with id 1", stmt.Type, stmt.Args)
	}
}

// TestShowTableCreation: SHOW TABLE CREATION (DorisParser.g4:351, show_param
// at DorisLegacyParser.g4:1386).
func TestShowTableCreation(t *testing.T) {
	stmt := parseOne(t, "SHOW TABLE CREATION").(*ast.ShowStmt)
	if stmt.Type != "TABLE CREATION" {
		t.Errorf("Type = %q, want TABLE CREATION", stmt.Type)
	}
}

// TestShowCreateLoad: SHOW CREATE LOAD FOR label (DorisParser.g4:391,
// DorisLegacyParser.g4:878).
func TestShowCreateLoad(t *testing.T) {
	stmt := parseOne(t, "SHOW CREATE LOAD FOR a").(*ast.ShowStmt)
	if stmt.Type != "CREATE LOAD" {
		t.Errorf("Type = %q, want CREATE LOAD", stmt.Type)
	}
	if stmt.Target == nil || stmt.Target.Parts[0] != "a" {
		t.Errorf("Target = %+v, want a", stmt.Target)
	}
}

// TestShowCreateMaterializedViewOn: SHOW CREATE MATERIALIZED VIEW name ON
// table (DorisParser.g4:355, show_param at DorisLegacyParser.g4:1386).
func TestShowCreateMaterializedViewOn(t *testing.T) {
	stmt := parseOne(t, "SHOW CREATE MATERIALIZED VIEW a ON a").(*ast.ShowStmt)
	if stmt.Type != "CREATE MATERIALIZED VIEW" {
		t.Errorf("Type = %q, want CREATE MATERIALIZED VIEW", stmt.Type)
	}
	if stmt.Target == nil || stmt.Target.Parts[0] != "a" {
		t.Errorf("Target = %+v, want a", stmt.Target)
	}
	if stmt.On == nil || stmt.On.Parts[0] != "a" {
		t.Errorf("On = %+v, want a", stmt.On)
	}
}

// TestShowCreateProcedure: SHOW CREATE PROCEDURE name (DorisParser.g4:47),
// a statement of its own rather than one handed to the CUP parser. The name is
// required, so SHOW CREATE PROCEDURE is a syntax error.
func TestShowCreateProcedure(t *testing.T) {
	stmt := parseOne(t, "SHOW CREATE PROCEDURE a").(*ast.ShowStmt)
	if stmt.Type != "CREATE PROCEDURE" {
		t.Errorf("Type = %q, want CREATE PROCEDURE", stmt.Type)
	}
	if stmt.Target == nil || stmt.Target.Parts[0] != "a" {
		t.Errorf("Target = %+v, want a", stmt.Target)
	}
	for _, sql := range []string{"SHOW CREATE PROCEDURE", "SHOW CREATE PROCEDURE a zzz"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestShowIndexStats: SHOW INDEX STATS table index (DorisParser.g4:769,
// show_param at DorisLegacyParser.g4:1386).
func TestShowIndexStats(t *testing.T) {
	stmt := parseOne(t, "SHOW INDEX STATS a a").(*ast.ShowStatsStmt)
	if stmt.Type != "INDEX" {
		t.Errorf("Type = %q, want INDEX", stmt.Type)
	}
	if stmt.Target == nil || stmt.Target.Parts[0] != "a" || stmt.Index != "a" {
		t.Errorf("Target = %+v Index = %q, want table a and index a", stmt.Target, stmt.Index)
	}
}
