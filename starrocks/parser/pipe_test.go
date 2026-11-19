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
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}
