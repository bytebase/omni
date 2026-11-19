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

	stmt = parseOne(t, `CREATE OR REPLACE PIPE db.p PROPERTIES('a' = 'a') AS INSERT INTO t VALUES(1)`).(*ast.CreatePipeStmt)
	if stmt.Name.String() != "db.p" || !stmt.OrReplace || stmt.IfNotExists {
		t.Errorf("stmt = %+v, want db.p with OR REPLACE", stmt)
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("Properties = %v, want one", stmt.Properties)
	}
	stmt = parseOne(t, `CREATE PIPE IF NOT EXISTS p AS INSERT INTO t SELECT 1`).(*ast.CreatePipeStmt)
	if !stmt.IfNotExists || stmt.OrReplace {
		t.Errorf("stmt = %+v, want IF NOT EXISTS alone", stmt)
	}

	for _, sql := range []string{
		`CREATE PIPE p AS INSERT a SELECT 1`,
		`CREATE PIPE p SELECT 1`,
		`CREATE PIPE AS INSERT INTO t VALUES(1)`,
		// The AST builder refuses the pair and any name but a database and a
		// name.
		`CREATE OR REPLACE PIPE IF NOT EXISTS p AS INSERT INTO t SELECT 1`,
		`CREATE PIPE a.b.c AS INSERT INTO t SELECT 1`,
		`DROP PIPE a.b.c`,
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

func TestFileCatalogKeepsItsParts(t *testing.T) {
	// createFileStatement takes a qualifiedName, so its parts survive: `a.b`
	// is two parts and `a.b` in backticks is one, as the server reads them.
	stmt := parseOne(t, `CREATE FILE 'x' IN a.b PROPERTIES('k' = 'v')`).(*ast.CreateFileStmt)
	if stmt.Catalog == nil || len(stmt.Catalog.Parts) != 2 || stmt.Catalog.Parts[1] != "b" {
		t.Errorf("Catalog = %v, want the two parts a and b", stmt.Catalog)
	}
	stmt = parseOne(t, "CREATE FILE 'x' IN `a.b` PROPERTIES('k' = 'v')").(*ast.CreateFileStmt)
	if stmt.Catalog == nil || len(stmt.Catalog.Parts) != 1 || stmt.Catalog.Parts[0] != "a.b" {
		t.Errorf("Catalog = %v, want the one part a.b", stmt.Catalog)
	}
	dropped := parseOne(t, `DROP FILE 'x' FROM a.b PROPERTIES('k' = 'v')`).(*ast.DropFileStmt)
	if dropped.Catalog == nil || len(dropped.Catalog.Parts) != 2 {
		t.Errorf("Catalog = %v, want the two parts a and b", dropped.Catalog)
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

func TestQualifySelectItemTakesAnAlias(t *testing.T) {
	// querySpecification's QUALIFY operand is a selectItem, so the alias the
	// rule allows is read: the server takes all four.
	for _, sql := range []string{
		`SELECT a FROM t QUALIFY ROW_NUMBER() OVER() AS rn = 1`,
		`SELECT a FROM t QUALIFY a AS rn = 1`,
		`SELECT a FROM t QUALIFY a rn = 1`,
		`SELECT a FROM t QUALIFY a AS 'rn' <= 2`,
	} {
		file, errs := parseForTest(sql)
		if len(errs) != 0 {
			t.Fatalf("parseForTest(%q): %v", sql, errs)
		}
		sel := file.Stmts[0].(*ast.SelectStmt)
		bin, ok := sel.Qualify.(*ast.BinaryExpr)
		if !ok {
			t.Fatalf("%q: Qualify = %T, want a comparison", sql, sel.Qualify)
		}
		if lit, ok := bin.Right.(*ast.Literal); !ok || lit.Kind != ast.LitInt {
			t.Errorf("%q: right = %T, want an integer bound", sql, bin.Right)
		}
	}
}

func TestStatementSpanEndsAtTheClosingParen(t *testing.T) {
	for sql, want := range map[string]int{
		`CREATE RESOURCE GROUP g WITH ('a'='b')`:                       len(`CREATE RESOURCE GROUP g WITH ('a'='b')`),
		`ALTER RESOURCE GROUP g WITH ('a'='b')`:                        len(`ALTER RESOURCE GROUP g WITH ('a'='b')`),
		`ALTER PIPE p SET('a'='b')`:                                    len(`ALTER PIPE p SET('a'='b')`),
		`CREATE PIPE p PROPERTIES('a'='b') AS INSERT INTO t VALUES(1)`: len(`CREATE PIPE p PROPERTIES('a'='b') AS INSERT INTO t VALUES(1)`),
		`CREATE USER 'u' PROPERTIES('a'='b')`:                          len(`CREATE USER 'u' PROPERTIES('a'='b')`),
		`ALTER USER 'u' SET PROPERTIES('a'='b')`:                       len(`ALTER USER 'u' SET PROPERTIES('a'='b')`),
	} {
		if got := ast.NodeLoc(parseOne(t, sql)).End; got != want {
			t.Errorf("%q: Loc.End = %d, want %d", sql, got, want)
		}
	}
}

func TestExternalCatalogsAreWalked(t *testing.T) {
	// A visitor must reach each catalog of BACKUP's and RESTORE's scope,
	// aliases included.
	for _, sql := range []string{
		`BACKUP EXTERNAL CATALOG(a, b) SNAPSHOT s TO r`,
		`RESTORE SNAPSHOT s FROM r EXTERNAL CATALOG(a AS b)`,
	} {
		var seen []string
		ast.Inspect(parseOne(t, sql), func(n ast.Node) bool {
			if ref, ok := n.(*ast.ExternalCatalogRef); ok {
				seen = append(seen, ref.Name)
			}
			return true
		})
		if len(seen) == 0 {
			t.Errorf("%q: Walk saw no ExternalCatalogRef", sql)
		}
	}
}

func TestQuotedFunctionNameTakesAnyArity(t *testing.T) {
	// A backticked name is an identifier, so the generic function production
	// reads it and LEFT's two-argument rule does not apply; the server takes
	// both `LEFT`(1) and `CHAR`().
	for _, sql := range []string{"SELECT `LEFT`(1)", "SELECT `CHAR`()"} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): %v", sql, errs)
		}
	}
	if _, errs := parseForTest(`SELECT LEFT(1)`); len(errs) == 0 {
		t.Error("LEFT(1) parsed, want a syntax error")
	}
}

func TestRestoreAliasNeedsDatabase(t *testing.T) {
	// The alias is DATABASE's own clause.
	if _, errs := parseForTest(`RESTORE SNAPSHOT s FROM r AS x`); len(errs) == 0 {
		t.Error("RESTORE ... AS x parsed, want a syntax error")
	}
	if _, errs := parseForTest(`RESTORE SNAPSHOT s FROM r DATABASE d AS x`); len(errs) != 0 {
		t.Errorf("RESTORE ... DATABASE d AS x: %v", errs)
	}
}

func TestDropResourceGroupSpanEndsAtTheParen(t *testing.T) {
	sql := `ALTER RESOURCE GROUP rg DROP(1, 2)`
	if got := ast.NodeLoc(parseOne(t, sql)).End; got != len(sql) {
		t.Errorf("Loc.End = %d, want %d", got, len(sql))
	}
}

func TestUserHostBracketsMarkADomain(t *testing.T) {
	// `'u'@['example.com']` is userWithHostAndBlanket (:2896): the brackets
	// name a domain, which StarRocks keeps apart from a plain host.
	stmt := parseOne(t, `CREATE USER u@['example.com']`).(*ast.CreateUserStmt)
	if !stmt.Name.Domain || stmt.Name.Host != "example.com" {
		t.Errorf("identity = %+v, want a domain host", stmt.Name)
	}
	stmt = parseOne(t, `CREATE USER u@'example.com'`).(*ast.CreateUserStmt)
	if stmt.Name.Domain || stmt.Name.Host != "example.com" {
		t.Errorf("identity = %+v, want a plain host", stmt.Name)
	}
}

func TestBlacklistRefusals(t *testing.T) {
	// The AST builder refuses an empty pattern and an id outside a signed
	// 64-bit integer.
	for _, sql := range []string{
		`ADD SQLBLACKLIST ''`,
		`ADD BACKEND BLACKLIST 9223372036854775808`,
		`DELETE SQLBLACKLIST 9223372036854775808`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestExternalCatalogSnapshotTakesNoOn(t *testing.T) {
	// "Can not sepcify `ON` clause for external catalog": the ON clause needs
	// a database, and an external-catalog snapshot has none.
	for _, sql := range []string{
		`BACKUP EXTERNAL CATALOG(a) SNAPSHOT s TO r ON(t)`,
		`RESTORE SNAPSHOT s FROM r EXTERNAL CATALOG(a) ON(t)`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
	for _, sql := range []string{
		`BACKUP DATABASE db SNAPSHOT s TO r ON(t)`,
		`RESTORE SNAPSHOT s FROM r DATABASE d ON(t)`,
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): %v", sql, errs)
		}
	}
}

func TestShowFiltersTakeCompoundPredicates(t *testing.T) {
	// The WHERE of these SHOW statements is a Predicate, compound forms and
	// NOT included, while a bare expression is not one.
	for _, sql := range []string{
		`SHOW STATS META WHERE a = 'x' AND b = 'y'`,
		`SHOW STATS META WHERE NOT a = 1`,
		`SHOW ANALYZE WHERE a = 1 OR b = 2`,
		`SHOW HISTOGRAM META WHERE a = 1`,
		`SHOW STATS META WHERE EXISTS (SELECT 1)`,
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): %v", sql, errs)
		}
	}
	for _, sql := range []string{
		`SHOW ANALYZE WHERE CURRENT_USER`,
		`SHOW STATS META WHERE a`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestDescFullFilesRefused(t *testing.T) {
	// descTableStatement (:561) is a qualified name or FILES propertyList;
	// FULL precedes neither.
	if _, errs := parseForTest(`DESC FULL FILES('a' = 'b')`); len(errs) == 0 {
		t.Error("DESC FULL FILES parsed, want a syntax error")
	}
}

func TestTranslateTerminatorIsInTheSpan(t *testing.T) {
	// `TRANSLATE presto;` is accepted because the terminator is the token
	// translateSQL matched, so the span covers it and the text is empty.
	stmt := parseOne(t, `TRANSLATE presto;`).(*ast.TranslateStmt)
	if stmt.SQL != "" {
		t.Errorf("SQL = %q, want empty", stmt.SQL)
	}
	if stmt.Loc.End != len(`TRANSLATE presto;`) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(`TRANSLATE presto;`))
	}
}

func TestLimitStaysWithinInt64(t *testing.T) {
	// LIMIT and OFFSET are signed 64-bit on the server; a longer integer is
	// "For input string".
	if _, errs := parseForTest(`SELECT 1 LIMIT 9223372036854775807`); len(errs) != 0 {
		t.Errorf("LIMIT 9223372036854775807: %v", errs)
	}
	for _, sql := range []string{
		`SELECT 1 LIMIT 9223372036854775808`,
		`SELECT 1 LIMIT 1, 9223372036854775808`,
		`SELECT 1 LIMIT 1 OFFSET 9223372036854775808`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestRefreshExternalTableNameParts(t *testing.T) {
	// qualifiedNameToTableName takes one, two or three parts.
	for _, sql := range []string{
		`REFRESH EXTERNAL TABLE t`,
		`REFRESH EXTERNAL TABLE db.t`,
		`REFRESH EXTERNAL TABLE cat.db.t`,
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): %v", sql, errs)
		}
	}
	if _, errs := parseForTest(`REFRESH EXTERNAL TABLE a.b.c.d`); len(errs) == 0 {
		t.Error("a four-part refresh parsed, want a syntax error")
	}
}

func TestDataCachePriorityIsAnInt32Magnitude(t *testing.T) {
	// Integer.parseInt reads the magnitude, so -2147483648 is refused while
	// -2147483647 is taken.
	if _, errs := parseForTest(`CREATE DATACACHE RULE *.*.* PRIORITY = - 2147483647`); len(errs) != 0 {
		t.Errorf("PRIORITY = - 2147483647: %v", errs)
	}
	if _, errs := parseForTest(`CREATE DATACACHE RULE *.*.* PRIORITY = - 2147483648`); len(errs) == 0 {
		t.Error("PRIORITY = - 2147483648 parsed, want a syntax error")
	}
}

func TestExportTemporaryPartition(t *testing.T) {
	// tableDesc (:2250) takes partitionNames (:2268), whose TEMPORARY modifier
	// names another partition namespace; the server takes the form.
	stmt := parseOne(t, `EXPORT TABLE t TEMPORARY PARTITION p TO '/x'`).(*ast.ExportStmt)
	if !stmt.TempPartition || len(stmt.Partition) != 1 {
		t.Errorf("stmt = %+v, want a temporary partition", stmt)
	}
	stmt = parseOne(t, `EXPORT TABLE t PARTITION p TO '/x'`).(*ast.ExportStmt)
	if stmt.TempPartition {
		t.Error("TempPartition = true, want false")
	}
}

func TestGroupByQuotedRollup(t *testing.T) {
	// A backticked ROLLUP is an identifier, so an empty argument list is an
	// ordinary call; the server takes `GROUP BY `ROLLUP`()`.
	for _, sql := range []string{"SELECT a FROM t GROUP BY `ROLLUP`()", "SELECT a FROM t GROUP BY `CUBE`()"} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): %v", sql, errs)
		}
	}
	for _, sql := range []string{`SELECT a FROM t GROUP BY ROLLUP()`, `SELECT a FROM t GROUP BY CUBE()`} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestSnapshotScopesAreExclusive(t *testing.T) {
	// A DATABASE clause and an external-catalog scope exclude each other, and
	// so do a DATABASE clause and a database-qualified label.
	for _, sql := range []string{
		`BACKUP EXTERNAL CATALOG(a) DATABASE db SNAPSHOT s TO r`,
		`BACKUP DATABASE db SNAPSHOT db.s TO r`,
		`RESTORE SNAPSHOT db.s FROM r DATABASE d`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
	for _, sql := range []string{
		`BACKUP DATABASE db SNAPSHOT s TO r`,
		`RESTORE SNAPSHOT s FROM r DATABASE d`,
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): %v", sql, errs)
		}
	}
}
