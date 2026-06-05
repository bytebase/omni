package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-ddl-bigquery node: CREATE ROW ACCESS POLICY, DROP ROW
// ACCESS POLICY, DROP ALL ROW ACCESS POLICIES (DDL-021/052). BigQuery-only at the
// union level (Spanner rejects ROW ACCESS POLICY outright, probed 2026-06-05);
// verdicts triangulated against the legacy GoogleSQLParser.g4 + BigQuery truth1.

func rapOf(t *testing.T, sql string) *ast.CreateRowAccessPolicyStmt {
	t.Helper()
	n := parseDDL(t, sql)
	r, ok := n.(*ast.CreateRowAccessPolicyStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateRowAccessPolicyStmt", sql, n)
	}
	return r
}

func TestCreateRowAccessPolicy_Full(t *testing.T) {
	// DDL-021.
	r := rapOf(t, "CREATE ROW ACCESS POLICY my_policy ON mydataset.mytable GRANT TO (\"user:alice@example.com\") FILTER USING (department = SESSION_USER())")
	if r.Name != "my_policy" {
		t.Errorf("Name = %q, want my_policy", r.Name)
	}
	if r.Table.String() != "mydataset.mytable" {
		t.Errorf("Table = %q", r.Table.String())
	}
	if !r.HasGrantTo || len(r.Grantees) != 1 {
		t.Errorf("HasGrantTo=%v Grantees=%d, want true/1", r.HasGrantTo, len(r.Grantees))
	}
	if r.Filter == nil {
		t.Error("Filter = nil, want the USING expression")
	}
}

func TestCreateRowAccessPolicy_NoName_NoGrant(t *testing.T) {
	// identifier? is optional; the grant-to clause is optional; FILTER is optional
	// (USING required).
	r := rapOf(t, "CREATE ROW ACCESS POLICY ON ds.t USING (TRUE)")
	if r.Name != "" {
		t.Errorf("Name = %q, want empty", r.Name)
	}
	if r.HasGrantTo {
		t.Error("HasGrantTo = true, want false")
	}
	if r.Filter == nil {
		t.Error("Filter = nil")
	}
}

func TestCreateRowAccessPolicy_OrReplaceToGrantees(t *testing.T) {
	// `TO grantee_list` (without the GRANT keyword) is the alternate grant-to form.
	r := rapOf(t, "CREATE OR REPLACE ROW ACCESS POLICY p ON ds.t GRANT TO ('user:a@b.com', 'group:g@b.com') FILTER USING (x > 0)")
	if !r.OrReplace {
		t.Error("OrReplace = false, want true")
	}
	if len(r.Grantees) != 2 {
		t.Errorf("Grantees = %d, want 2", len(r.Grantees))
	}
}

func TestCreateRowAccessPolicy_IfNotExists(t *testing.T) {
	r := rapOf(t, "CREATE ROW ACCESS POLICY IF NOT EXISTS p ON ds.t FILTER USING (TRUE)")
	if !r.IfNotExists {
		t.Error("IfNotExists = false, want true")
	}
}

func TestCreateRowAccessPolicy_Rejects(t *testing.T) {
	cases := []string{
		"CREATE ROW ACCESS POLICY p ON ds.t",             // missing FILTER USING
		"CREATE ROW ACCESS POLICY p FILTER USING (TRUE)", // missing ON table
		"CREATE ROW ACCESS POLICY p ON ds.t USING TRUE",  // USING without parens
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- DROP ROW ACCESS POLICY / DROP ALL ROW ACCESS POLICIES ---

func TestDropRowAccessPolicy(t *testing.T) {
	// DDL-052.
	d := bqDropOf(t, "DROP ROW ACCESS POLICY IF EXISTS my_policy ON mydataset.mytable")
	if d.Object != ast.BQDropRowAccessPolicy {
		t.Errorf("Object = %v, want ROW ACCESS POLICY", d.Object)
	}
	if !d.IfExists {
		t.Error("IfExists = false, want true")
	}
	if d.Name.String() != "my_policy" || d.OnTable.String() != "mydataset.mytable" {
		t.Errorf("Name=%q OnTable=%q", d.Name.String(), d.OnTable.String())
	}
}

func TestDropAllRowAccessPolicies(t *testing.T) {
	// DDL-052.
	n := parseDDL(t, "DROP ALL ROW ACCESS POLICIES ON mydataset.mytable")
	d, ok := n.(*ast.DropAllRowAccessPoliciesStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.DropAllRowAccessPoliciesStmt", n)
	}
	if d.Table.String() != "mydataset.mytable" {
		t.Errorf("Table = %q", d.Table.String())
	}
}

func TestDropRowAccessPolicy_Rejects(t *testing.T) {
	cases := []string{
		"DROP ROW ACCESS POLICY p",     // missing ON table
		"DROP ROW POLICY p ON t",       // missing ACCESS
		"DROP ALL ROW ACCESS POLICIES", // missing ON table
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}
