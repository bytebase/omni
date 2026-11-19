package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// The statements this file covers are read by parser/resource_group.go. Every
// accepted form is engine-verified: the pinned StarRocks 3.4.10 takes it, and
// every refused form is one it answers with a syntax error.

func TestCreateResourceGroup(t *testing.T) {
	stmt := parseOne(t, `CREATE RESOURCE GROUP rg TO(user = 'a') WITH('cpu_core_limit' = '10')`).(*ast.CreateResourceGroupStmt)
	if stmt.Name.String() != "rg" {
		t.Errorf("Name = %q, want rg", stmt.Name.String())
	}
	if len(stmt.Classifiers) != 1 || len(stmt.Classifiers[0]) != 1 {
		t.Fatalf("Classifiers = %v, want one classifier", stmt.Classifiers)
	}
	if len(stmt.Properties) != 1 || stmt.Properties[0].Key != "cpu_core_limit" {
		t.Errorf("Properties = %v, want cpu_core_limit", stmt.Properties)
	}

	stmt = parseOne(t, `CREATE RESOURCE GROUP IF NOT EXISTS rg WITH('a' = 'a')`).(*ast.CreateResourceGroupStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists = false, want true")
	}
	stmt = parseOne(t, `CREATE RESOURCE GROUP OR REPLACE rg WITH('a' = 'a')`).(*ast.CreateResourceGroupStmt)
	if !stmt.OrReplace {
		t.Error("OrReplace = false, want true")
	}

	// Each classifier is a Predicate: the AST builder casts it.
	for _, sql := range []string{
		`CREATE RESOURCE GROUP rg TO(1) WITH('a' = 'a')`,
		`CREATE RESOURCE GROUP rg TO('a') WITH('a' = 'a')`,
		`CREATE RESOURCE GROUP rg TO(a) WITH('a' = 'a')`,
		`CREATE RESOURCE GROUP rg`,
		`CREATE RESOURCE GROUP rg WITH`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestAlterResourceGroup(t *testing.T) {
	stmt := parseOne(t, `ALTER RESOURCE GROUP rg ADD(user = 'a'), (role = 'b')`).(*ast.AlterResourceGroupStmt)
	if stmt.Action != "ADD" || len(stmt.Classifiers) != 2 {
		t.Errorf("Action = %q Classifiers = %v, want ADD with two", stmt.Action, stmt.Classifiers)
	}

	stmt = parseOne(t, `ALTER RESOURCE GROUP rg DROP(1, 2)`).(*ast.AlterResourceGroupStmt)
	if stmt.Action != "DROP" || len(stmt.ClassIDs) != 2 {
		t.Errorf("Action = %q ClassIDs = %v, want DROP with two", stmt.Action, stmt.ClassIDs)
	}

	stmt = parseOne(t, `ALTER RESOURCE GROUP rg DROP ALL`).(*ast.AlterResourceGroupStmt)
	if stmt.Action != "DROP ALL" {
		t.Errorf("Action = %q, want DROP ALL", stmt.Action)
	}

	stmt = parseOne(t, `ALTER RESOURCE GROUP rg WITH('a' = 'a')`).(*ast.AlterResourceGroupStmt)
	if stmt.Action != "WITH" || len(stmt.Properties) != 1 {
		t.Errorf("Action = %q Properties = %v, want WITH with one", stmt.Action, stmt.Properties)
	}

	for _, sql := range []string{
		`ALTER RESOURCE GROUP rg DROP`,
		`ALTER RESOURCE GROUP rg DROP(ALL)`,
		`ALTER RESOURCE GROUP rg ADD`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestDropResourceGroup(t *testing.T) {
	stmt := parseOne(t, `DROP RESOURCE GROUP rg`).(*ast.DropResourceGroupStmt)
	if stmt.Name.String() != "rg" {
		t.Errorf("Name = %q, want rg", stmt.Name.String())
	}
	if _, errs := parseForTest(`DROP RESOURCE GROUP IF EXISTS rg`); len(errs) == 0 {
		t.Error("DROP RESOURCE GROUP IF EXISTS parsed, want a syntax error")
	}
}

func TestResourceGroupIsNotResource(t *testing.T) {
	// CREATE RESOURCE still reads the resource statement: GROUP is what tells
	// the two apart.
	if _, errs := parseForTest(`CREATE RESOURCE "rg" PROPERTIES("type" = "spark")`); len(errs) != 0 {
		t.Errorf("CREATE RESOURCE broke: %v", errs)
	}
}

func TestRefreshExternalTable(t *testing.T) {
	stmt := parseOne(t, `REFRESH EXTERNAL TABLE hive.t`).(*ast.RefreshExternalTableStmt)
	if stmt.Name.String() != "hive.t" || len(stmt.Partitions) != 0 {
		t.Errorf("Name = %q Partitions = %v", stmt.Name.String(), stmt.Partitions)
	}
	stmt = parseOne(t, `REFRESH EXTERNAL TABLE hive.t PARTITION('p1', 'p2')`).(*ast.RefreshExternalTableStmt)
	if len(stmt.Partitions) != 2 || stmt.Partitions[1] != "p2" {
		t.Errorf("Partitions = %v, want [p1 p2]", stmt.Partitions)
	}
	for _, sql := range []string{
		`REFRESH EXTERNAL TABLE`,
		`REFRESH EXTERNAL TABLE hive.t PARTITION(p1)`,
		`REFRESH EXTERNAL TABLE hive.t PARTITION('p1',)`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}

func TestDataCacheRules(t *testing.T) {
	stmt := parseOne(t, `CREATE DATACACHE RULE *.*.* PRIORITY = 1`).(*ast.CreateDataCacheRuleStmt)
	if stmt.Target != "*.*.*" || stmt.Priority != "1" || stmt.Negative {
		t.Errorf("Target = %q Priority = %q Negative = %v", stmt.Target, stmt.Priority, stmt.Negative)
	}

	stmt = parseOne(t, `CREATE DATACACHE RULE hive.db.t WHERE a = 1 PRIORITY = - 2 PROPERTIES('a' = 'a')`).(*ast.CreateDataCacheRuleStmt)
	if stmt.Target != "hive.db.t" || stmt.Where == nil || stmt.Priority != "2" || !stmt.Negative {
		t.Errorf("stmt = %+v, want hive.db.t with a WHERE and -2", stmt)
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("Properties = %v, want one", stmt.Properties)
	}

	if stmt := parseOne(t, `DROP DATACACHE RULE 3`).(*ast.DropDataCacheRuleStmt); stmt.ID != "3" {
		t.Errorf("ID = %q, want 3", stmt.ID)
	}
	if _, ok := parseOne(t, `CLEAR DATACACHE RULES`).(*ast.ClearDataCacheRulesStmt); !ok {
		t.Error("CLEAR DATACACHE RULES is not a ClearDataCacheRulesStmt")
	}

	for _, sql := range []string{
		`CREATE DATACACHE RULE *.* PRIORITY = 1`,
		`CREATE DATACACHE RULE *.*.*`,
		`CREATE DATACACHE RULE *.*.* PRIORITY`,
		`DROP DATACACHE RULE a`,
		`CLEAR DATACACHE RULE`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want a syntax error", sql)
		}
	}
}
