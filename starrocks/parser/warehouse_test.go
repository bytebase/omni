package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// WAREHOUSE
// ---------------------------------------------------------------------------

// TestCreateWarehouse: createWarehouseStatement (StarRocks.g4:2038).
func TestCreateWarehouse(t *testing.T) {
	stmt := parseOne(t, `CREATE WAREHOUSE IF NOT EXISTS wh COMMENT 'shared' PROPERTIES("k" = "v")`).(*ast.CreateWarehouseStmt)

	if stmt.Name != "wh" || !stmt.IfNotExists {
		t.Errorf("Name = %q IfNotExists = %v, want wh and true", stmt.Name, stmt.IfNotExists)
	}
	if stmt.Comment != "shared" {
		t.Errorf("Comment = %q, want shared", stmt.Comment)
	}
	if len(stmt.Properties) != 1 || stmt.Properties[0].Key != "k" || stmt.Properties[0].Value != "v" {
		t.Errorf("Properties = %+v, want k=v", stmt.Properties)
	}
	if stmt.Tag() != ast.T_CreateWarehouseStmt {
		t.Errorf("Tag() = %v, want T_CreateWarehouseStmt", stmt.Tag())
	}

	// The shortest form is a bare name; it ends at that name.
	short := parseOne(t, `CREATE WAREHOUSE wh`).(*ast.CreateWarehouseStmt)
	if short.Comment != "" || short.Properties != nil {
		t.Errorf("Comment = %q Properties = %v, want neither", short.Comment, short.Properties)
	}
	if end, want := short.Loc.End, len(`CREATE WAREHOUSE wh`); end != want {
		t.Errorf("Loc.End = %d, want %d", end, want)
	}

	// identifierOrString takes a string too, as it does in the grammar.
	if got := parseOne(t, `CREATE WAREHOUSE 'wh'`).(*ast.CreateWarehouseStmt).Name; got != "wh" {
		t.Errorf("Name = %q, want wh", got)
	}
}

// TestDropSuspendResumeWarehouse: dropWarehouseStatement (StarRocks.g4:2043),
// suspendWarehouseStatement (:2047) and resumeWarehouseStatement (:2051). The
// last two take a plain identifier, not identifierOrString.
func TestDropSuspendResumeWarehouse(t *testing.T) {
	drop := parseOne(t, `DROP WAREHOUSE IF EXISTS wh`).(*ast.DropWarehouseStmt)
	if drop.Name != "wh" || !drop.IfExists {
		t.Errorf("DROP: Name = %q IfExists = %v, want wh and true", drop.Name, drop.IfExists)
	}

	suspend := parseOne(t, `SUSPEND WAREHOUSE IF EXISTS wh`).(*ast.SuspendWarehouseStmt)
	if suspend.Name != "wh" || !suspend.IfExists {
		t.Errorf("SUSPEND: Name = %q IfExists = %v, want wh and true", suspend.Name, suspend.IfExists)
	}

	resume := parseOne(t, `RESUME WAREHOUSE wh`).(*ast.ResumeWarehouseStmt)
	if resume.Name != "wh" || resume.IfExists {
		t.Errorf("RESUME: Name = %q IfExists = %v, want wh and false", resume.Name, resume.IfExists)
	}

	tags := map[ast.NodeTag]ast.Node{
		ast.T_DropWarehouseStmt:    drop,
		ast.T_SuspendWarehouseStmt: suspend,
		ast.T_ResumeWarehouseStmt:  resume,
	}
	for want, n := range tags {
		if got := n.Tag(); got != want {
			t.Errorf("%T.Tag() = %v, want %v", n, got, want)
		}
	}
}

// TestWarehouseRejected: every warehouse statement names its warehouse, and
// SUSPEND/RESUME take the identifier the grammar states.
func TestWarehouseRejected(t *testing.T) {
	for _, sql := range []string{
		`CREATE WAREHOUSE`,
		`CREATE WAREHOUSE wh zzz`,
		`CREATE WAREHOUSE wh COMMENT`,
		`CREATE WAREHOUSE wh PROPERTIES()`,
		`DROP WAREHOUSE`,
		`DROP WAREHOUSE wh zzz`,
		`SUSPEND WAREHOUSE`,
		`SUSPEND WAREHOUSE 'wh'`,
		`SUSPEND WAREHOUSE wh zzz`,
		`RESUME WAREHOUSE`,
		`RESUME WAREHOUSE 'wh'`,
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed; the engine rejects it", sql)
		}
	}
}
