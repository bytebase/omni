package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// parser/pipe.go and parser/blacklist.go, engine-verified against the pinned
// StarRocks 3.4.10.

func TestCreatePipe(t *testing.T) {
	stmt := parseOne(t, `CREATE PIPE p AS INSERT INTO t SELECT a FROM s`).(*ast.CreatePipeStmt)
	if stmt.Name.String() != "p" || stmt.OrReplace || stmt.IfNotExists {
		t.Errorf("stmt = %+v, want pipe p", stmt)
	}
	if _, ok := stmt.Insert.(*ast.InsertStmt); !ok {
		t.Errorf("Insert = %T, want *ast.InsertStmt", stmt.Insert)
	}

	stmt = parseOne(t, `CREATE OR REPLACE PIPE IF NOT EXISTS db.p PROPERTIES('a' = 'a') AS INSERT INTO t VALUES(1)`).(*ast.CreatePipeStmt)
	if stmt.Name.String() != "db.p" || !stmt.OrReplace || !stmt.IfNotExists {
		t.Errorf("stmt = %+v, want db.p with OR REPLACE and IF NOT EXISTS", stmt)
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("Properties = %v, want one", stmt.Properties)
	}

	for _, sql := range []string{
		`CREATE PIPE p AS INSERT a SELECT 1`,
		`CREATE PIPE p SELECT 1`,
		`CREATE PIPE AS INSERT INTO t VALUES(1)`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestAlterAndDropPipe(t *testing.T) {
	for sql, want := range map[string]string{
		`ALTER PIPE p SUSPEND`:            "SUSPEND",
		`ALTER PIPE p RESUME`:             "RESUME",
		`ALTER PIPE p RETRY ALL`:          "RETRY ALL",
		`ALTER PIPE p RETRY FILE 'f.csv'`: "RETRY FILE",
		`ALTER PIPE p SET('a' = 'a')`:     "SET",
	} {
		stmt := parseOne(t, sql).(*ast.AlterPipeStmt)
		if stmt.Action != want {
			t.Errorf("%q: Action = %q, want %q", sql, stmt.Action, want)
		}
	}
	if stmt := parseOne(t, `ALTER PIPE p RETRY FILE 'f.csv'`).(*ast.AlterPipeStmt); stmt.File != "f.csv" {
		t.Errorf("File = %q, want f.csv", stmt.File)
	}

	stmt := parseOne(t, `DROP PIPE IF EXISTS db.p`).(*ast.DropPipeStmt)
	if stmt.Name.String() != "db.p" || !stmt.IfExists {
		t.Errorf("stmt = %+v, want db.p IF EXISTS", stmt)
	}

	for _, sql := range []string{
		`ALTER PIPE p`,
		`ALTER PIPE p RETRY`,
		`ALTER PIPE p RETRY FILE f`,
		`DROP PIPE IF p`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestDescPipe(t *testing.T) {
	for _, sql := range []string{`DESC PIPE db.p`, `DESCRIBE PIPE db.p`} {
		stmt := parseOne(t, sql).(*ast.DescPipeStmt)
		if stmt.Name.String() != "db.p" {
			t.Errorf("%q: Name = %q, want db.p", sql, stmt.Name.String())
		}
	}
	// PIPE is nonReserved, so without a name the word is a table's.
	if _, ok := parseOne(t, `DESC pipe`).(*ast.DescribeStmt); !ok {
		t.Error(`DESC pipe is not a DescribeStmt of the table named pipe`)
	}
	if _, errs := parseForTest(`DESC PIPE`); len(errs) != 0 {
		t.Errorf("DESC PIPE: %v", errs)
	}
}

func TestSqlBlacklist(t *testing.T) {
	stmt := parseOne(t, `ADD SQLBLACKLIST 'select 1'`).(*ast.BlacklistStmt)
	if !stmt.Add || stmt.Backend || stmt.SQL != "select 1" {
		t.Errorf("stmt = %+v, want an added SQL blacklist entry", stmt)
	}
	stmt = parseOne(t, `DELETE SQLBLACKLIST 1, 2`).(*ast.BlacklistStmt)
	if stmt.Add || stmt.Backend || len(stmt.IDs) != 2 {
		t.Errorf("stmt = %+v, want a deleted pair of ids", stmt)
	}
	if _, errs := parseForTest(`DELETE SQLBLACKLIST 'a'`); len(errs) == 0 {
		t.Error("DELETE SQLBLACKLIST 'a' parsed, want a syntax error")
	}
}

func TestBackendBlacklist(t *testing.T) {
	stmt := parseOne(t, `ADD BACKEND BLACKLIST 1, 2`).(*ast.BlacklistStmt)
	if !stmt.Add || !stmt.Backend || len(stmt.IDs) != 2 {
		t.Errorf("stmt = %+v, want an added backend pair", stmt)
	}
	stmt = parseOne(t, `DELETE BACKEND BLACKLIST 3`).(*ast.BlacklistStmt)
	if stmt.Add || !stmt.Backend || len(stmt.IDs) != 1 {
		t.Errorf("stmt = %+v, want a deleted backend id", stmt)
	}
	// ADD BACKEND is still the comment-free statement it was: the blacklist
	// form needs BLACKLIST after BACKEND.
	for _, sql := range []string{`ADD SQLBLACKLIST 1`, `ADD BACKEND BLACKLIST 'a'`} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestPlanAdvisor(t *testing.T) {
	stmt := parseOne(t, `ALTER PLAN ADVISOR ADD SELECT a FROM t`).(*ast.AlterPlanAdvisorStmt)
	if stmt.Action != "ADD" {
		t.Errorf("Action = %q, want ADD", stmt.Action)
	}
	if _, ok := stmt.Query.(*ast.SelectStmt); !ok {
		t.Errorf("Query = %T, want *ast.SelectStmt", stmt.Query)
	}
	stmt = parseOne(t, `ALTER PLAN ADVISOR DROP 'a'`).(*ast.AlterPlanAdvisorStmt)
	if stmt.Action != "DROP" || stmt.ID != "a" {
		t.Errorf("stmt = %+v, want DROP of a", stmt)
	}
	if _, ok := parseOne(t, `TRUNCATE PLAN ADVISOR`).(*ast.TruncatePlanAdvisorStmt); !ok {
		t.Error("TRUNCATE PLAN ADVISOR is not a TruncatePlanAdvisorStmt")
	}

	for _, sql := range []string{
		`ALTER PLAN ADVISOR`,
		`ALTER PLAN ADVISOR ADD`,
		`ALTER PLAN ADVISOR DROP a`,
		`TRUNCATE PLAN`,
		// queryStatement's explainDesc prefix is refused by the server's AST
		// builder: "query should not be a explain stmt".
		`ALTER PLAN ADVISOR ADD EXPLAIN SELECT 1`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestBackupCatalogListTakesNoAlias(t *testing.T) {
	// backupStatement's list is identifierList; restoreStatement's is
	// identifierWithAliasList, and the server refuses an alias in the former.
	for _, sql := range []string{
		`BACKUP EXTERNAL CATALOG(a) SNAPSHOT s TO r`,
		`BACKUP ALL EXTERNAL CATALOGS SNAPSHOT s TO r`,
		`BACKUP DATABASE db SNAPSHOT s TO r`,
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): %v", sql, errs)
		}
	}
	if _, errs := parseForTest(`BACKUP EXTERNAL CATALOG(a AS b) SNAPSHOT s TO r`); len(errs) == 0 {
		t.Error("BACKUP with an aliased catalog parsed, want a syntax error")
	}
}

func TestFileCatalogKeepsItsQualifiedName(t *testing.T) {
	// createFileStatement takes a qualifiedName, so a catalog and a database
	// both survive: `IN catalog.db` is not truncated to db.
	stmt := parseOne(t, `CREATE FILE 'x' IN catalog.db PROPERTIES('k' = 'v')`).(*ast.CreateFileStmt)
	if stmt.Database != "catalog.db" {
		t.Errorf("Database = %q, want catalog.db", stmt.Database)
	}
	dropped := parseOne(t, `DROP FILE 'x' FROM catalog.db PROPERTIES('k' = 'v')`).(*ast.DropFileStmt)
	if dropped.Database != "catalog.db" {
		t.Errorf("Database = %q, want catalog.db", dropped.Database)
	}
}

func TestNumericPrefixesStayNumbers(t *testing.T) {
	// '.' and '+' are no DIGIT_IDENTIFIER characters, so '1.2a' is a decimal
	// and its label rather than a column named 1.2a, and '1e+2a' alike; '.1E2a'
	// is one DOT_IDENTIFIER because the identifier is the longer match. The
	// server takes all three.
	for _, sql := range []string{`SELECT 1.2a`, `SELECT 1e+2a`} {
		file, errs := parseForTest(sql)
		if len(errs) != 0 {
			t.Fatalf("parseForTest(%q): %v", sql, errs)
		}
		sel := file.Stmts[0].(*ast.SelectStmt)
		if lit, ok := sel.Items[0].Expr.(*ast.Literal); !ok || lit.Kind != ast.LitFloat {
			t.Errorf("%q: item = %T, want a float literal", sql, sel.Items[0].Expr)
		}
	}

	file, errs := parseForTest("SELECT named_struct('x', 1).1E2a")
	if len(errs) != 0 {
		t.Fatalf("dereference: %v", errs)
	}
	sel := file.Stmts[0].(*ast.SelectStmt)
	deref, ok := sel.Items[0].Expr.(*ast.DereferenceExpr)
	if !ok || deref.Field != "1E2a" {
		t.Errorf("item = %#v, want the field 1E2a", sel.Items[0].Expr)
	}
}

func TestDigitIdentifierSwallowsTheExponent(t *testing.T) {
	// DIGIT_IDENTIFIER (StarRocksLex.g4:564) outruns DOUBLE_VALUE, so `1e2a`
	// and `1E2a` are names rather than a float and its label; the server
	// takes both.
	for _, sql := range []string{`SELECT 1e2a`, `SELECT 1E2a`} {
		file, errs := parseForTest(sql)
		if len(errs) != 0 {
			t.Fatalf("parseForTest(%q): %v", sql, errs)
		}
		sel := file.Stmts[0].(*ast.SelectStmt)
		if col, ok := sel.Items[0].Expr.(*ast.ColumnRef); !ok || col.Name.Parts[0] != sql[len("SELECT "):] {
			t.Errorf("%q: item = %T, want the column named %q", sql, sel.Items[0].Expr, sql[len("SELECT "):])
		}
	}
}

func TestRestoreCatalogAliasesAreStructural(t *testing.T) {
	// identifierWithAliasList (:2858) keeps the catalog and the name it is
	// restored as apart, so a catalog literally named `a AS b` is not the
	// same input as renaming a to b.
	stmt := parseOne(t, "RESTORE SNAPSHOT s FROM r EXTERNAL CATALOG(`a AS b`)").(*ast.RestoreStmt)
	if len(stmt.ExternalCatalogs) != 1 {
		t.Fatalf("ExternalCatalogs = %v, want one", stmt.ExternalCatalogs)
	}
	if ref := stmt.ExternalCatalogs[0]; ref.Name != "a AS b" || ref.Alias != "" {
		t.Errorf("ref = %+v, want the name `a AS b` and no alias", ref)
	}

	stmt = parseOne(t, `RESTORE SNAPSHOT s FROM r EXTERNAL CATALOG(a AS b, c)`).(*ast.RestoreStmt)
	if len(stmt.ExternalCatalogs) != 2 {
		t.Fatalf("ExternalCatalogs = %v, want two", stmt.ExternalCatalogs)
	}
	if ref := stmt.ExternalCatalogs[0]; ref.Name != "a" || ref.Alias != "b" {
		t.Errorf("ref = %+v, want a renamed to b", ref)
	}
	if ref := stmt.ExternalCatalogs[1]; ref.Name != "c" || ref.Alias != "" {
		t.Errorf("ref = %+v, want c with no alias", ref)
	}
}

func TestRefreshExternalTableSpanEndsAtTheParen(t *testing.T) {
	sql := `REFRESH EXTERNAL TABLE hive.t PARTITION('p1')`
	stmt := parseOne(t, sql).(*ast.RefreshExternalTableStmt)
	if stmt.Loc.End != len(sql) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(sql))
	}
}

func TestDescFilesNamesNoTarget(t *testing.T) {
	// DESC FILES names no table: it describes a files() table function, which
	// analysis reads as a read-only statement by Files.
	stmt := parseOne(t, `DESC FILES('a' = 'b')`).(*ast.DescribeStmt)
	if !stmt.Files || stmt.Target != nil || len(stmt.Properties) != 1 {
		t.Errorf("stmt = %+v, want the files form with one property", stmt)
	}
}
