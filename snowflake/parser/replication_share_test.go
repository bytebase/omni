package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func mustCreateGroup(t *testing.T, input string) *ast.CreateReplicationGroupStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateReplicationGroupStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateReplicationGroupStmt", input, node)
	}
	return stmt
}

func mustAlterGroup(t *testing.T, input string) *ast.AlterReplicationGroupStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterReplicationGroupStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterReplicationGroupStmt", input, node)
	}
	return stmt
}

func mustCreateAccount(t *testing.T, input string) *ast.CreateAccountStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateAccountStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateAccountStmt", input, node)
	}
	return stmt
}

func mustAlterAccount(t *testing.T, input string) *ast.AlterAccountStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterAccountStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterAccountStmt", input, node)
	}
	return stmt
}

func mustCreateShare(t *testing.T, input string) *ast.CreateShareStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateShareStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateShareStmt", input, node)
	}
	return stmt
}

func mustAlterShare(t *testing.T, input string) *ast.AlterShareStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterShareStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterShareStmt", input, node)
	}
	return stmt
}

// findGroupOption returns the named option, or nil.
func findGroupOption(opts []*ast.GroupOption, name string) *ast.GroupOption {
	for _, o := range opts {
		if o.Name == name {
			return o
		}
	}
	return nil
}

// names joins a list of object names with commas for assertions.
func names(ns []*ast.ObjectName) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = n.String()
	}
	return strings.Join(parts, ",")
}

// ---------------------------------------------------------------------------
// CREATE FAILOVER GROUP / REPLICATION GROUP — primary form
// ---------------------------------------------------------------------------

func TestParseCreateGroup_Primary(t *testing.T) {
	t.Run("failover minimal", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE FAILOVER GROUP fg OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = org.acct")
		if !stmt.Failover {
			t.Error("Failover not set")
		}
		if stmt.Name.String() != "fg" {
			t.Errorf("Name = %q, want fg", stmt.Name.String())
		}
		if stmt.Replica != nil {
			t.Errorf("unexpected Replica: %v", stmt.Replica)
		}
		ot := findGroupOption(stmt.Options, "OBJECT_TYPES")
		if ot == nil || len(ot.Values) != 1 || ot.Values[0] != "ROLES" {
			t.Errorf("OBJECT_TYPES = %v, want [ROLES]", ot)
		}
		aa := findGroupOption(stmt.Options, "ALLOWED_ACCOUNTS")
		if aa == nil || len(aa.Values) != 1 || aa.Values[0] != "ORG.ACCT" {
			t.Errorf("ALLOWED_ACCOUNTS = %v, want [ORG.ACCT]", aa)
		}
	})

	t.Run("replication minimal", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE REPLICATION GROUP rg OBJECT_TYPES = DATABASES ALLOWED_ACCOUNTS = org.acct")
		if stmt.Failover {
			t.Error("Failover should be false for REPLICATION GROUP")
		}
	})

	t.Run("or replace + if not exists", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE OR REPLACE FAILOVER GROUP IF NOT EXISTS fg OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = org.acct")
		if !stmt.OrReplace || !stmt.IfNotExists {
			t.Errorf("OrReplace=%v IfNotExists=%v, want both true", stmt.OrReplace, stmt.IfNotExists)
		}
	})

	t.Run("multi-word object types", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE FAILOVER GROUP fg OBJECT_TYPES = ACCOUNT PARAMETERS, RESOURCE MONITORS, NETWORK POLICIES ALLOWED_ACCOUNTS = org.acct")
		ot := findGroupOption(stmt.Options, "OBJECT_TYPES")
		want := []string{"ACCOUNT PARAMETERS", "RESOURCE MONITORS", "NETWORK POLICIES"}
		if ot == nil || strings.Join(ot.Values, "|") != strings.Join(want, "|") {
			t.Errorf("OBJECT_TYPES = %v, want %v", ot.Values, want)
		}
	})

	// Structural correctness for the trickiest behavior: an unparenthesized
	// comma list of multi-word object types must NOT absorb the following
	// space-separated option name (ALLOWED_ACCOUNTS) as a list element, and the
	// two options must come out distinct.
	t.Run("list boundary stops at next option", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE FAILOVER GROUP fg OBJECT_TYPES = DATABASES, ROLES ALLOWED_ACCOUNTS = org.a, org.b")
		ot := findGroupOption(stmt.Options, "OBJECT_TYPES")
		if ot == nil || strings.Join(ot.Values, "|") != "DATABASES|ROLES" {
			t.Errorf("OBJECT_TYPES = %v, want [DATABASES ROLES]", ot)
		}
		aa := findGroupOption(stmt.Options, "ALLOWED_ACCOUNTS")
		if aa == nil || strings.Join(aa.Values, "|") != "ORG.A|ORG.B" {
			t.Errorf("ALLOWED_ACCOUNTS = %v, want [ORG.A ORG.B]", aa)
		}
		if len(stmt.Options) != 2 {
			t.Errorf("expected exactly 2 options, got %d: %+v", len(stmt.Options), stmt.Options)
		}
	})

	t.Run("all optional list options", func(t *testing.T) {
		in := "CREATE FAILOVER GROUP fg " +
			"OBJECT_TYPES = DATABASES, INTEGRATIONS " +
			"ALLOWED_DATABASES = db1, db2 " +
			"ALLOWED_EXTERNAL_VOLUMES = ev1 " +
			"ALLOWED_SHARES = sh1, sh2 " +
			"ALLOWED_INTEGRATION_TYPES = SECURITY INTEGRATIONS, API INTEGRATIONS " +
			"ALLOWED_ACCOUNTS = org.a1, org.a2 " +
			"IGNORE EDITION CHECK " +
			"REPLICATION_SCHEDULE = '10 MINUTE' " +
			"OPTIMIZED_REFRESH = TRUE " +
			"ERROR_INTEGRATION = my_int"
		stmt := mustCreateGroup(t, in)
		if !stmt.IgnoreEditionCheck {
			t.Error("IgnoreEditionCheck not set")
		}
		if ad := findGroupOption(stmt.Options, "ALLOWED_DATABASES"); ad == nil || len(ad.Values) != 2 {
			t.Errorf("ALLOWED_DATABASES = %v, want 2 values", ad)
		}
		if sh := findGroupOption(stmt.Options, "ALLOWED_SHARES"); sh == nil || len(sh.Values) != 2 {
			t.Errorf("ALLOWED_SHARES = %v, want 2 values", sh)
		}
		rs := findGroupOption(stmt.Options, "REPLICATION_SCHEDULE")
		if rs == nil || rs.Lit == nil || rs.Lit.Value != "10 MINUTE" {
			t.Errorf("REPLICATION_SCHEDULE = %v, want literal '10 MINUTE'", rs)
		}
		opt := findGroupOption(stmt.Options, "OPTIMIZED_REFRESH")
		if opt == nil || len(opt.Values) != 1 || opt.Values[0] != "TRUE" {
			t.Errorf("OPTIMIZED_REFRESH = %v, want [TRUE]", opt)
		}
		ei := findGroupOption(stmt.Options, "ERROR_INTEGRATION")
		if ei == nil || len(ei.Values) != 1 || ei.Values[0] != "MY_INT" {
			t.Errorf("ERROR_INTEGRATION = %v, want [MY_INT]", ei)
		}
	})

	t.Run("with tag clause", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE FAILOVER GROUP fg OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = org.acct WITH TAG (cost_center = 'sales')")
		if len(stmt.Tags) != 1 || stmt.Tags[0].Name.String() != "cost_center" || stmt.Tags[0].Value != "sales" {
			t.Errorf("Tags = %+v, want one cost_center=sales", stmt.Tags)
		}
	})

	t.Run("bare tag clause (no WITH)", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE FAILOVER GROUP fg OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = org.acct TAG (t = 'v')")
		if len(stmt.Tags) != 1 {
			t.Errorf("Tags = %+v, want one", stmt.Tags)
		}
	})

	t.Run("replication_schedule cron string", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE FAILOVER GROUP fg OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = org.acct REPLICATION_SCHEDULE = 'USING CRON 0 0 10-20 * TUE,THU UTC'")
		rs := findGroupOption(stmt.Options, "REPLICATION_SCHEDULE")
		if rs == nil || rs.Lit == nil || !strings.HasPrefix(rs.Lit.Value, "USING CRON") {
			t.Errorf("REPLICATION_SCHEDULE = %v", rs)
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE FAILOVER GROUP / REPLICATION GROUP — secondary AS REPLICA OF form
// ---------------------------------------------------------------------------

func TestParseCreateGroup_Secondary(t *testing.T) {
	t.Run("failover replica", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE FAILOVER GROUP myfg2 AS REPLICA OF myorg.myaccount.myfg")
		if stmt.Replica == nil {
			t.Fatal("Replica not set")
		}
		if stmt.Replica.String() != "myorg.myaccount.myfg" {
			t.Errorf("Replica = %q, want myorg.myaccount.myfg", stmt.Replica.String())
		}
		if len(stmt.Options) != 0 {
			t.Errorf("secondary form should have no options, got %v", stmt.Options)
		}
	})

	t.Run("replication replica", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE REPLICATION GROUP myrg2 AS REPLICA OF myorg.myaccount.myrg")
		if stmt.Failover {
			t.Error("Failover should be false")
		}
		if stmt.Replica == nil || stmt.Replica.String() != "myorg.myaccount.myrg" {
			t.Errorf("Replica = %v", stmt.Replica)
		}
	})

	t.Run("or replace + if not exists replica", func(t *testing.T) {
		stmt := mustCreateGroup(t, "CREATE OR REPLACE FAILOVER GROUP IF NOT EXISTS fg2 AS REPLICA OF a.b.c")
		if !stmt.OrReplace || !stmt.IfNotExists || stmt.Replica == nil {
			t.Errorf("OrReplace=%v IfNotExists=%v Replica=%v", stmt.OrReplace, stmt.IfNotExists, stmt.Replica)
		}
	})
}

// ---------------------------------------------------------------------------
// ALTER FAILOVER GROUP / REPLICATION GROUP — all action variants
// ---------------------------------------------------------------------------

func TestParseAlterGroup_Actions(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg RENAME TO fg2")
		if stmt.Action != ast.AlterGroupRename || stmt.NewName.String() != "fg2" {
			t.Errorf("Action=%d NewName=%v", stmt.Action, stmt.NewName)
		}
	})

	t.Run("if exists rename", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP IF EXISTS fg RENAME TO fg2")
		if !stmt.IfExists {
			t.Error("IfExists not set")
		}
	})

	t.Run("set object types + schedule", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg SET OBJECT_TYPES = DATABASES, ROLES REPLICATION_SCHEDULE = '20 MINUTE'")
		if stmt.Action != ast.AlterGroupSet {
			t.Fatalf("Action=%d, want AlterGroupSet", stmt.Action)
		}
		if ot := findGroupOption(stmt.Options, "OBJECT_TYPES"); ot == nil || len(ot.Values) != 2 {
			t.Errorf("OBJECT_TYPES = %v", ot)
		}
	})

	t.Run("set integration types", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER REPLICATION GROUP rg SET OBJECT_TYPES = INTEGRATIONS ALLOWED_INTEGRATION_TYPES = SECURITY INTEGRATIONS")
		if stmt.Action != ast.AlterGroupSet {
			t.Errorf("Action=%d", stmt.Action)
		}
	})

	t.Run("set tag", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg SET TAG t1 = 'v1', t2 = 'v2'")
		if stmt.Action != ast.AlterGroupSetTag || len(stmt.Tags) != 2 {
			t.Errorf("Action=%d Tags=%+v", stmt.Action, stmt.Tags)
		}
	})

	t.Run("unset property", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER REPLICATION GROUP rg UNSET REPLICATION_SCHEDULE")
		if stmt.Action != ast.AlterGroupUnset || len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "REPLICATION_SCHEDULE" {
			t.Errorf("Action=%d UnsetProps=%v", stmt.Action, stmt.UnsetProps)
		}
	})

	t.Run("unset tag", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg UNSET TAG t1, t2")
		if stmt.Action != ast.AlterGroupUnsetTag || len(stmt.UnsetTags) != 2 {
			t.Errorf("Action=%d UnsetTags=%v", stmt.Action, stmt.UnsetTags)
		}
	})

	t.Run("add to allowed databases", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg ADD db1, db2 TO ALLOWED_DATABASES")
		if stmt.Action != ast.AlterGroupAdd || stmt.ListTarget != "ALLOWED_DATABASES" || names(stmt.Names) != "db1,db2" {
			t.Errorf("Action=%d Target=%q Names=%q", stmt.Action, stmt.ListTarget, names(stmt.Names))
		}
	})

	t.Run("add to allowed shares", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg ADD sh1 TO ALLOWED_SHARES")
		if stmt.ListTarget != "ALLOWED_SHARES" {
			t.Errorf("Target=%q", stmt.ListTarget)
		}
	})

	t.Run("add to allowed accounts ignore edition check", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg ADD org.a1, org.a2 TO ALLOWED_ACCOUNTS IGNORE EDITION CHECK")
		if stmt.ListTarget != "ALLOWED_ACCOUNTS" || !stmt.IgnoreEditionCheck || names(stmt.Names) != "org.a1,org.a2" {
			t.Errorf("Target=%q Ignore=%v Names=%q", stmt.ListTarget, stmt.IgnoreEditionCheck, names(stmt.Names))
		}
	})

	t.Run("remove from allowed databases", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg REMOVE db1 FROM ALLOWED_DATABASES")
		if stmt.Action != ast.AlterGroupRemove || stmt.ListTarget != "ALLOWED_DATABASES" {
			t.Errorf("Action=%d Target=%q", stmt.Action, stmt.ListTarget)
		}
	})

	t.Run("move databases to failover group", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg MOVE DATABASES db1, db2 TO FAILOVER GROUP fg2")
		if stmt.Action != ast.AlterGroupMove || stmt.MoveKind != "DATABASES" || stmt.MoveTo.String() != "fg2" {
			t.Errorf("Action=%d MoveKind=%q MoveTo=%v", stmt.Action, stmt.MoveKind, stmt.MoveTo)
		}
		if names(stmt.Names) != "db1,db2" {
			t.Errorf("Names=%q", names(stmt.Names))
		}
	})

	t.Run("move shares to replication group", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER REPLICATION GROUP rg MOVE SHARES sh1 TO REPLICATION GROUP rg2")
		if stmt.MoveKind != "SHARES" || stmt.MoveTo.String() != "rg2" {
			t.Errorf("MoveKind=%q MoveTo=%v", stmt.MoveKind, stmt.MoveTo)
		}
	})

	t.Run("refresh", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER REPLICATION GROUP rg REFRESH")
		if stmt.Action != ast.AlterGroupRefresh {
			t.Errorf("Action=%d", stmt.Action)
		}
	})

	t.Run("primary", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg PRIMARY")
		if stmt.Action != ast.AlterGroupPrimary {
			t.Errorf("Action=%d", stmt.Action)
		}
	})

	t.Run("suspend", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg SUSPEND")
		if stmt.Action != ast.AlterGroupSuspend || stmt.Immediate {
			t.Errorf("Action=%d Immediate=%v", stmt.Action, stmt.Immediate)
		}
	})

	t.Run("suspend immediate", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER FAILOVER GROUP fg SUSPEND IMMEDIATE")
		if !stmt.Immediate {
			t.Error("Immediate not set")
		}
	})

	t.Run("resume", func(t *testing.T) {
		stmt := mustAlterGroup(t, "ALTER REPLICATION GROUP rg RESUME")
		if stmt.Action != ast.AlterGroupResume {
			t.Errorf("Action=%d", stmt.Action)
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE ACCOUNT / MANAGED ACCOUNT
// ---------------------------------------------------------------------------

func TestParseCreateAccount(t *testing.T) {
	t.Run("account full", func(t *testing.T) {
		in := "CREATE ACCOUNT myaccount1 " +
			"ADMIN_NAME = admin ADMIN_PASSWORD = 'TestPassword1' " +
			"FIRST_NAME = Jane LAST_NAME = Smith EMAIL = 'myemail@myorg.org' " +
			"MUST_CHANGE_PASSWORD = TRUE EDITION = ENTERPRISE"
		stmt := mustCreateAccount(t, in)
		if stmt.Managed {
			t.Error("Managed should be false")
		}
		if stmt.Name.String() != "myaccount1" {
			t.Errorf("Name=%q", stmt.Name.String())
		}
		an := findGroupOption(stmt.Options, "ADMIN_NAME")
		if an == nil || len(an.Values) != 1 || an.Values[0] != "ADMIN" {
			t.Errorf("ADMIN_NAME=%v", an)
		}
		ap := findGroupOption(stmt.Options, "ADMIN_PASSWORD")
		if ap == nil || ap.Lit == nil || ap.Lit.Value != "TestPassword1" {
			t.Errorf("ADMIN_PASSWORD=%v", ap)
		}
		ed := findGroupOption(stmt.Options, "EDITION")
		if ed == nil || len(ed.Values) != 1 || ed.Values[0] != "ENTERPRISE" {
			t.Errorf("EDITION=%v", ed)
		}
	})

	t.Run("account with rsa public key", func(t *testing.T) {
		stmt := mustCreateAccount(t, "CREATE ACCOUNT a ADMIN_NAME = admin ADMIN_RSA_PUBLIC_KEY = 'key' EMAIL = 'e@x.com' EDITION = STANDARD")
		if findGroupOption(stmt.Options, "ADMIN_RSA_PUBLIC_KEY") == nil {
			t.Error("ADMIN_RSA_PUBLIC_KEY not captured")
		}
	})

	t.Run("account region + comment", func(t *testing.T) {
		stmt := mustCreateAccount(t, "CREATE ACCOUNT a ADMIN_NAME = admin ADMIN_PASSWORD = 'p' EMAIL = 'e@x.com' EDITION = BUSINESS_CRITICAL REGION = aws_us_west_2 COMMENT = 'hi'")
		if findGroupOption(stmt.Options, "REGION") == nil {
			t.Error("REGION not captured")
		}
		c := findGroupOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit == nil || c.Lit.Value != "hi" {
			t.Errorf("COMMENT=%v", c)
		}
	})

	t.Run("managed account comma-separated", func(t *testing.T) {
		stmt := mustCreateAccount(t, "CREATE MANAGED ACCOUNT ra1 ADMIN_NAME = admin, ADMIN_PASSWORD = 'pw', TYPE = READER")
		if !stmt.Managed {
			t.Error("Managed not set")
		}
		ty := findGroupOption(stmt.Options, "TYPE")
		if ty == nil || len(ty.Values) != 1 || ty.Values[0] != "READER" {
			t.Errorf("TYPE=%v", ty)
		}
		if len(stmt.Options) != 3 {
			t.Errorf("expected 3 options, got %d: %+v", len(stmt.Options), stmt.Options)
		}
	})

	t.Run("managed account with comment", func(t *testing.T) {
		stmt := mustCreateAccount(t, "CREATE MANAGED ACCOUNT ra1 ADMIN_NAME = admin, ADMIN_PASSWORD = 'pw', TYPE = READER, COMMENT = 'reader'")
		if len(stmt.Options) != 4 {
			t.Errorf("expected 4 options, got %d", len(stmt.Options))
		}
	})

	t.Run("or replace managed account", func(t *testing.T) {
		stmt := mustCreateAccount(t, "CREATE OR REPLACE MANAGED ACCOUNT ra1 ADMIN_NAME = admin, ADMIN_PASSWORD = 'pw', TYPE = READER")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
	})
}

// ---------------------------------------------------------------------------
// ALTER ACCOUNT — current-account + cross-account forms
// ---------------------------------------------------------------------------

func TestParseAlterAccount(t *testing.T) {
	t.Run("set param", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET NETWORK_POLICY = mypolicy")
		if stmt.Name != nil {
			t.Errorf("Name should be nil for current-account, got %v", stmt.Name)
		}
		if stmt.Action != ast.AlterAccountSet {
			t.Fatalf("Action=%d", stmt.Action)
		}
		np := findGroupOption(stmt.Options, "NETWORK_POLICY")
		if np == nil || len(np.Values) != 1 || np.Values[0] != "MYPOLICY" {
			t.Errorf("NETWORK_POLICY=%v", np)
		}
	})

	t.Run("set boolean param", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET DISABLE_USER_PRIVILEGE_GRANTS = TRUE")
		o := findGroupOption(stmt.Options, "DISABLE_USER_PRIVILEGE_GRANTS")
		if o == nil || len(o.Values) != 1 || o.Values[0] != "TRUE" {
			t.Errorf("opt=%v", o)
		}
	})

	t.Run("set multiple params", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET TIMEZONE = 'UTC', STATEMENT_TIMEOUT_IN_SECONDS = 3600")
		if len(stmt.Options) != 2 {
			t.Errorf("expected 2 options, got %d: %+v", len(stmt.Options), stmt.Options)
		}
	})

	t.Run("set resource monitor", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET RESOURCE_MONITOR = my_monitor")
		o := findGroupOption(stmt.Options, "RESOURCE_MONITOR")
		if o == nil || o.Values[0] != "MY_MONITOR" {
			t.Errorf("RESOURCE_MONITOR=%v", o)
		}
	})

	t.Run("unset param", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT UNSET NETWORK_POLICY")
		if stmt.Action != ast.AlterAccountUnset || len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "NETWORK_POLICY" {
			t.Errorf("Action=%d UnsetProps=%v", stmt.Action, stmt.UnsetProps)
		}
	})

	t.Run("unset multiple params", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT UNSET TIMEZONE, NETWORK_POLICY")
		if len(stmt.UnsetProps) != 2 {
			t.Errorf("UnsetProps=%v", stmt.UnsetProps)
		}
	})

	t.Run("set tag", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET TAG cost_center = 'sales'")
		if stmt.Action != ast.AlterAccountSetTag || len(stmt.Tags) != 1 {
			t.Errorf("Action=%d Tags=%+v", stmt.Action, stmt.Tags)
		}
	})

	t.Run("unset tag", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT UNSET TAG cost_center")
		if stmt.Action != ast.AlterAccountUnsetTag || len(stmt.UnsetTags) != 1 {
			t.Errorf("Action=%d UnsetTags=%v", stmt.Action, stmt.UnsetTags)
		}
	})

	t.Run("set packages policy force", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET PACKAGES POLICY packages_policy_prod_1 FORCE")
		if stmt.Action != ast.AlterAccountSetPolicy {
			t.Fatalf("Action=%d, want SetPolicy", stmt.Action)
		}
		if stmt.PolicyKind != "PACKAGES POLICY" || stmt.PolicyName.String() != "packages_policy_prod_1" || !stmt.Force {
			t.Errorf("Kind=%q Name=%v Force=%v", stmt.PolicyKind, stmt.PolicyName, stmt.Force)
		}
	})

	t.Run("set password policy", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET PASSWORD POLICY my_pp")
		if stmt.PolicyKind != "PASSWORD POLICY" || stmt.Force {
			t.Errorf("Kind=%q Force=%v", stmt.PolicyKind, stmt.Force)
		}
	})

	t.Run("set authentication policy with scope + force", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET AUTHENTICATION POLICY my_ap FOR ALL PERSON USERS FORCE")
		if stmt.PolicyKind != "AUTHENTICATION POLICY" || stmt.PolicyScope != "FOR ALL PERSON USERS" || !stmt.Force {
			t.Errorf("Kind=%q Scope=%q Force=%v", stmt.PolicyKind, stmt.PolicyScope, stmt.Force)
		}
	})

	t.Run("set session policy", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET SESSION POLICY my_sp")
		if stmt.PolicyKind != "SESSION POLICY" {
			t.Errorf("Kind=%q", stmt.PolicyKind)
		}
	})

	t.Run("set feature policy for all applications", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT SET FEATURE POLICY my_fp FOR ALL APPLICATIONS FORCE")
		if stmt.PolicyKind != "FEATURE POLICY" || stmt.PolicyScope != "FOR ALL APPLICATIONS" || !stmt.Force {
			t.Errorf("Kind=%q Scope=%q Force=%v", stmt.PolicyKind, stmt.PolicyScope, stmt.Force)
		}
	})

	t.Run("unset packages policy", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT UNSET PACKAGES POLICY")
		if stmt.Action != ast.AlterAccountUnsetPolicy || stmt.PolicyKind != "PACKAGES POLICY" {
			t.Errorf("Action=%d Kind=%q", stmt.Action, stmt.PolicyKind)
		}
	})

	t.Run("unset authentication policy with scope", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT UNSET AUTHENTICATION POLICY FOR ALL SERVICE USERS")
		if stmt.Action != ast.AlterAccountUnsetPolicy || stmt.PolicyScope != "FOR ALL SERVICE USERS" {
			t.Errorf("Action=%d Scope=%q", stmt.Action, stmt.PolicyScope)
		}
	})

	t.Run("cross-account set edition", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT myaccount SET EDITION = 'ENTERPRISE'")
		if stmt.Name == nil || stmt.Name.String() != "myaccount" {
			t.Errorf("Name=%v", stmt.Name)
		}
		if stmt.Action != ast.AlterAccountSet {
			t.Errorf("Action=%d", stmt.Action)
		}
	})

	t.Run("cross-account set is_org_admin", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT myaccount SET IS_ORG_ADMIN = TRUE")
		if stmt.Name == nil {
			t.Error("Name not set")
		}
	})

	t.Run("rename", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT myaccount RENAME TO myaccount2")
		if stmt.Action != ast.AlterAccountRename || stmt.Name.String() != "myaccount" || stmt.NewName.String() != "myaccount2" {
			t.Errorf("Action=%d Name=%v NewName=%v", stmt.Action, stmt.Name, stmt.NewName)
		}
		if stmt.SaveOldURL != nil {
			t.Errorf("SaveOldURL should be nil, got %v", *stmt.SaveOldURL)
		}
	})

	t.Run("rename save old url", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT myaccount RENAME TO myaccount2 SAVE_OLD_URL = FALSE")
		if stmt.SaveOldURL == nil || *stmt.SaveOldURL {
			t.Errorf("SaveOldURL=%v, want false", stmt.SaveOldURL)
		}
	})

	t.Run("drop old url", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT myaccount DROP OLD URL")
		if stmt.Action != ast.AlterAccountDropURL || stmt.Organization {
			t.Errorf("Action=%d Organization=%v", stmt.Action, stmt.Organization)
		}
	})

	t.Run("drop old organization url", func(t *testing.T) {
		stmt := mustAlterAccount(t, "ALTER ACCOUNT myaccount DROP OLD ORGANIZATION URL")
		if !stmt.Organization {
			t.Error("Organization not set")
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE SHARE
// ---------------------------------------------------------------------------

func TestParseCreateShare(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		stmt := mustCreateShare(t, "CREATE SHARE sales_s")
		if stmt.Name.String() != "sales_s" {
			t.Errorf("Name=%q", stmt.Name.String())
		}
		if stmt.OrReplace || stmt.IfNotExists || len(stmt.Options) != 0 {
			t.Errorf("unexpected: %+v", stmt)
		}
	})

	t.Run("or replace", func(t *testing.T) {
		stmt := mustCreateShare(t, "CREATE OR REPLACE SHARE s")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := mustCreateShare(t, "CREATE SHARE IF NOT EXISTS s")
		if !stmt.IfNotExists {
			t.Error("IfNotExists not set")
		}
	})

	t.Run("comment", func(t *testing.T) {
		stmt := mustCreateShare(t, "CREATE SHARE s COMMENT = 'my share'")
		c := findGroupOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit == nil || c.Lit.Value != "my share" {
			t.Errorf("COMMENT=%v", c)
		}
	})
}

// ---------------------------------------------------------------------------
// ALTER SHARE
// ---------------------------------------------------------------------------

func TestParseAlterShare(t *testing.T) {
	t.Run("add accounts", func(t *testing.T) {
		stmt := mustAlterShare(t, "ALTER SHARE s ADD ACCOUNTS = org.acct1, org.acct2")
		if stmt.Action != ast.AlterShareAdd || names(stmt.Accounts) != "org.acct1,org.acct2" {
			t.Errorf("Action=%d Accounts=%q", stmt.Action, names(stmt.Accounts))
		}
	})

	t.Run("add accounts share restrictions", func(t *testing.T) {
		stmt := mustAlterShare(t, "ALTER SHARE s ADD ACCOUNTS = org.a SHARE_RESTRICTIONS = FALSE")
		if stmt.ShareRestrictions == nil || *stmt.ShareRestrictions {
			t.Errorf("ShareRestrictions=%v, want false", stmt.ShareRestrictions)
		}
	})

	t.Run("remove accounts", func(t *testing.T) {
		stmt := mustAlterShare(t, "ALTER SHARE s REMOVE ACCOUNTS = org.a")
		if stmt.Action != ast.AlterShareRemove {
			t.Errorf("Action=%d", stmt.Action)
		}
	})

	t.Run("if exists add accounts", func(t *testing.T) {
		stmt := mustAlterShare(t, "ALTER SHARE IF EXISTS s ADD ACCOUNTS = org.a")
		if !stmt.IfExists {
			t.Error("IfExists not set")
		}
	})

	t.Run("set accounts + comment", func(t *testing.T) {
		stmt := mustAlterShare(t, "ALTER SHARE s SET ACCOUNTS = org.a1, org.a2 COMMENT = 'updated'")
		if stmt.Action != ast.AlterShareSet || names(stmt.Accounts) != "org.a1,org.a2" {
			t.Errorf("Action=%d Accounts=%q", stmt.Action, names(stmt.Accounts))
		}
		c := findGroupOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit.Value != "updated" {
			t.Errorf("COMMENT=%v", c)
		}
	})

	t.Run("set comment only", func(t *testing.T) {
		stmt := mustAlterShare(t, "ALTER SHARE s SET COMMENT = 'x'")
		if stmt.Action != ast.AlterShareSet || len(stmt.Accounts) != 0 {
			t.Errorf("Action=%d Accounts=%v", stmt.Action, stmt.Accounts)
		}
	})

	t.Run("set tag", func(t *testing.T) {
		stmt := mustAlterShare(t, "ALTER SHARE s SET TAG t = 'v'")
		if stmt.Action != ast.AlterShareSetTag || len(stmt.Tags) != 1 {
			t.Errorf("Action=%d Tags=%+v", stmt.Action, stmt.Tags)
		}
	})

	t.Run("unset tag", func(t *testing.T) {
		stmt := mustAlterShare(t, "ALTER SHARE s UNSET TAG t1, t2")
		if stmt.Action != ast.AlterShareUnsetTag || len(stmt.UnsetTags) != 2 {
			t.Errorf("Action=%d UnsetTags=%v", stmt.Action, stmt.UnsetTags)
		}
	})

	t.Run("unset comment", func(t *testing.T) {
		stmt := mustAlterShare(t, "ALTER SHARE s UNSET COMMENT")
		if stmt.Action != ast.AlterShareUnset || len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "COMMENT" {
			t.Errorf("Action=%d UnsetProps=%v", stmt.Action, stmt.UnsetProps)
		}
	})
}

// ---------------------------------------------------------------------------
// Negative tests — malformed statements must be rejected
// ---------------------------------------------------------------------------

func TestParseReplicationShare_Negative(t *testing.T) {
	bad := []string{
		// FAILOVER / REPLICATION GROUP
		"CREATE FAILOVER GROUP", // missing GROUP? no — missing name
		"CREATE FAILOVER GROUP fg OBJECT_TYPES = (DATABASES) ALLOWED_ACCOUNTS = o.a", // parenthesized list not valid (docs+grammar use bare comma lists)
		"CREATE FAILOVER fg OBJECT_TYPES = ROLES",                                    // missing GROUP keyword
		"CREATE FAILOVER GROUP fg OBJECT_TYPES =",                                    // option missing value
		"CREATE FAILOVER GROUP fg AS REPLICA",                                        // missing OF
		"CREATE FAILOVER GROUP fg AS REPLICA OF",                                     // missing source name
		"CREATE FAILOVER GROUP fg AS OF a.b.c",                                       // missing REPLICA
		"CREATE REPLICATION fg OBJECT_TYPES = ROLES",                                 // missing GROUP keyword
		"ALTER FAILOVER GROUP",                                                       // missing name
		"ALTER FAILOVER GROUP fg",                                                    // missing action
		"ALTER FAILOVER GROUP fg RENAME",                                             // missing TO
		"ALTER FAILOVER GROUP fg RENAME TO",                                          // missing new name
		"ALTER FAILOVER GROUP fg SET",                                                // nothing to set
		"ALTER FAILOVER GROUP fg SET TAG",                                            // missing tag assignment
		"ALTER FAILOVER GROUP fg SET TAG t =",                                        // tag missing value
		"ALTER FAILOVER GROUP fg UNSET",                                              // missing property
		"ALTER FAILOVER GROUP fg UNSET TAG",                                          // missing tag name
		"ALTER FAILOVER GROUP fg ADD db1 TO",                                         // missing target
		"ALTER FAILOVER GROUP fg ADD db1 TO FROBNICATE",                              // bad target
		"ALTER FAILOVER GROUP fg ADD TO ALLOWED_DATABASES",                           // missing name list
		"ALTER FAILOVER GROUP fg REMOVE db1 ALLOWED_DATABASES",                       // missing FROM
		"ALTER FAILOVER GROUP fg MOVE db1 TO FAILOVER GROUP fg2",                     // missing DATABASES/SHARES
		"ALTER FAILOVER GROUP fg MOVE DATABASES db1 TO fg2",                          // missing FAILOVER/REPLICATION GROUP
		"ALTER FAILOVER GROUP fg MOVE DATABASES db1 TO FAILOVER fg2",                 // missing GROUP
		"ALTER FAILOVER GROUP fg FROBNICATE",                                         // unknown action
		// ACCOUNT
		"CREATE ACCOUNT",                    // missing name
		"CREATE ACCOUNT a ADMIN_NAME =",     // option missing value
		"CREATE MANAGED a ADMIN_NAME = x",   // missing ACCOUNT keyword
		"ALTER ACCOUNT",                     // missing action
		"ALTER ACCOUNT SET",                 // nothing to set
		"ALTER ACCOUNT SET TAG",             // missing tag assignment
		"ALTER ACCOUNT UNSET",               // missing property
		"ALTER ACCOUNT a RENAME",            // missing TO
		"ALTER ACCOUNT a RENAME TO",         // missing new name
		"ALTER ACCOUNT a DROP OLD",          // missing URL
		"ALTER ACCOUNT a DROP URL",          // missing OLD
		"ALTER ACCOUNT a FROBNICATE",        // unknown action
		"ALTER ACCOUNT SET PACKAGES POLICY", // missing policy name
		// SHARE
		"CREATE SHARE",                 // missing name
		"CREATE SHARE IF NOT s",        // malformed IF NOT EXISTS
		"ALTER SHARE",                  // missing name
		"ALTER SHARE s",                // missing action
		"ALTER SHARE s ADD",            // missing ACCOUNTS
		"ALTER SHARE s ADD ACCOUNTS",   // missing =
		"ALTER SHARE s ADD ACCOUNTS =", // missing account list
		"ALTER SHARE s SET",            // nothing to set
		"ALTER SHARE s SET TAG",        // missing tag assignment
		"ALTER SHARE s UNSET",          // missing property
		"ALTER SHARE s UNSET TAG",      // missing tag name
		"ALTER SHARE s FROBNICATE",     // unknown action
	}
	for _, in := range bad {
		t.Run(in, func(t *testing.T) {
			result := ParseBestEffort(in)
			if len(result.Errors) == 0 {
				t.Errorf("expected parse error for %q, got none (stmts=%d)", in, len(result.File.Stmts))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Loc accuracy
// ---------------------------------------------------------------------------

func TestParseReplicationShare_Loc(t *testing.T) {
	// CREATE Loc starts at the CREATE keyword (offset 0).
	in := "CREATE SHARE s"
	stmt := mustCreateShare(t, in)
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(in) {
		t.Errorf("CREATE SHARE Loc = %+v, want {0, %d}", stmt.Loc, len(in))
	}

	in2 := "CREATE FAILOVER GROUP fg OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = org.acct"
	g := mustCreateGroup(t, in2)
	if g.Loc.Start != 0 || g.Loc.End != len(in2) {
		t.Errorf("CREATE GROUP Loc = %+v, want {0, %d}", g.Loc, len(in2))
	}

	// ALTER sub-parsers anchor Loc.Start at the object-type keyword, not ALTER
	// (matching the established ALTER TABLE/STAGE/integration convention).
	const alterPrefix = len("ALTER ")
	ai := "ALTER SHARE s UNSET COMMENT"
	as := mustAlterShare(t, ai)
	if as.Loc.Start != alterPrefix || as.Loc.End != len(ai) {
		t.Errorf("ALTER SHARE Loc = %+v, want {%d, %d}", as.Loc, alterPrefix, len(ai))
	}

	ag := "ALTER FAILOVER GROUP fg PRIMARY"
	gs := mustAlterGroup(t, ag)
	if gs.Loc.Start != alterPrefix || gs.Loc.End != len(ag) {
		t.Errorf("ALTER GROUP Loc = %+v, want {%d, %d}", gs.Loc, alterPrefix, len(ag))
	}

	// ALTER ACCOUNT anchors at the ACCOUNT keyword.
	aa := "ALTER ACCOUNT SET NETWORK_POLICY = mypolicy"
	aas := mustAlterAccount(t, aa)
	if aas.Loc.Start != alterPrefix || aas.Loc.End != len(aa) {
		t.Errorf("ALTER ACCOUNT Loc = %+v, want {%d, %d}", aas.Loc, alterPrefix, len(aa))
	}

	// Second statement gets a non-zero base offset; Loc must still be absolute.
	multi := "SELECT 1;\nCREATE SHARE s COMMENT = 'x'"
	res := ParseBestEffort(multi)
	if len(res.Errors) != 0 {
		t.Fatalf("multi parse errors: %v", res.Errors)
	}
	if len(res.File.Stmts) != 2 {
		t.Fatalf("expected 2 stmts, got %d", len(res.File.Stmts))
	}
	cs := res.File.Stmts[1].(*ast.CreateShareStmt)
	base := strings.Index(multi, "CREATE SHARE")
	if cs.Loc.Start != base {
		t.Errorf("second-stmt Loc.Start = %d, want %d", cs.Loc.Start, base)
	}
	c := findGroupOption(cs.Options, "COMMENT")
	if c == nil || c.Lit == nil || c.Lit.Value != "x" {
		t.Errorf("COMMENT literal not correct after base offset: %v", c)
	}
}

// ---------------------------------------------------------------------------
// Walk integration — the new statements must be walkable without panic, and the
// walker must reach their ObjectName children.
// ---------------------------------------------------------------------------

func TestParseReplicationShare_Walk(t *testing.T) {
	inputs := []string{
		"CREATE FAILOVER GROUP fg OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = org.acct",
		"CREATE FAILOVER GROUP fg2 AS REPLICA OF a.b.c",
		"ALTER FAILOVER GROUP fg ADD db1 TO ALLOWED_DATABASES",
		"ALTER FAILOVER GROUP fg MOVE DATABASES db1 TO FAILOVER GROUP fg2",
		"CREATE ACCOUNT a ADMIN_NAME = x ADMIN_PASSWORD = 'p' EMAIL = 'e@x.com' EDITION = STANDARD",
		"ALTER ACCOUNT a RENAME TO b",
		"CREATE SHARE s COMMENT = 'x'",
		"ALTER SHARE s SET ACCOUNTS = org.a COMMENT = 'x'",
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			node := mustParseOne(t, in)
			count := 0
			ast.Inspect(node, func(n ast.Node) bool {
				if n != nil {
					count++
				}
				return true
			})
			if count == 0 {
				t.Errorf("walk visited no nodes for %q", in)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Official docs corpus — every CREATE SHARE / ALTER ACCOUNT statement in the
// corresponding corpus directory must parse with zero errors to its expected AST
// type. The official docs are the authoritative oracle (truth1).
// ---------------------------------------------------------------------------

func TestReplicationShare_OfficialCorpus(t *testing.T) {
	t.Run("create-share", func(t *testing.T) {
		runReplShareCorpus(t, "testdata/official/create-share", "CREATE", "SHARE",
			func(n ast.Node) bool { _, ok := n.(*ast.CreateShareStmt); return ok })
	})
	t.Run("alter-account", func(t *testing.T) {
		runReplShareCorpus(t, "testdata/official/alter-account", "ALTER", "ACCOUNT",
			func(n ast.Node) bool { _, ok := n.(*ast.AlterAccountStmt); return ok })
	})
}

// runReplShareCorpus parses every .sql file in dir and asserts each statement
// whose first two words are <verb> <obj> parses cleanly to the wanted type.
func runReplShareCorpus(t *testing.T, dir, verb, obj string, wantFn func(ast.Node) bool) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read corpus dir %s: %v", dir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			for _, seg := range Split(string(data)) {
				text := strings.TrimSpace(seg.Text)
				if text == "" {
					continue
				}
				fields := strings.Fields(strings.ToUpper(text))
				if len(fields) < 2 || fields[0] != verb || fields[1] != obj {
					continue
				}
				node, errs := parseSingle(seg.Text, seg.ByteStart)
				if len(errs) > 0 {
					t.Errorf("statement %q produced %d error(s): %v", text, len(errs), errs)
					continue
				}
				if !wantFn(node) {
					t.Errorf("statement %q parsed to unexpected type %T", text, node)
				}
			}
		})
	}
}
