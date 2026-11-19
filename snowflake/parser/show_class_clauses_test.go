package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func parseOneShow(t *testing.T, sql string) *ast.ShowStmt {
	t.Helper()
	file, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): %d statements", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.ShowStmt)
	if !ok {
		t.Fatalf("Parse(%q): %T, want *ast.ShowStmt", sql, file.Stmts[0])
	}
	if stmt.Loc.End != len(sql) {
		t.Errorf("Parse(%q): Loc.End = %d, want %d", sql, stmt.Loc.End, len(sql))
	}
	return stmt
}

// TestShowClassScopes covers the IN scopes only one class's page writes.
func TestShowClassScopes(t *testing.T) {
	for _, tt := range []struct {
		sql, class, scope, name string
	}{
		{"SHOW FAILOVER GROUPS IN ACCOUNT myaccount1", "FAILOVER GROUPS", "ACCOUNT", "MYACCOUNT1"},
		{"SHOW REPLICATION GROUPS IN ACCOUNT myaccount1", "REPLICATION GROUPS", "ACCOUNT", "MYACCOUNT1"},
		{"SHOW TERSE ROLES LIKE 'a%' IN CLASS snowflake.ml.anomaly_detection LIMIT 5", "ROLES", "CLASS", "SNOWFLAKE.ML.ANOMALY_DETECTION"},
		{"SHOW FUNCTIONS IN CLASS snowflake.core.budget", "FUNCTIONS", "CLASS", "SNOWFLAKE.CORE.BUDGET"},
		{"SHOW PROCEDURES IN CLASS budget", "PROCEDURES", "CLASS", "BUDGET"},
		{"SHOW GIT BRANCHES LIKE 'm%' IN GIT REPOSITORY db.s.repo", "GIT BRANCHES", "GIT REPOSITORY", "DB.S.REPO"},
		{"SHOW GIT TAGS IN GIT REPOSITORY repo", "GIT TAGS", "GIT REPOSITORY", "REPO"},
		{"SHOW VERSIONS IN OPENFLOW CONNECTOR my_db.my_schema.my_connector", "VERSIONS", "OPENFLOW CONNECTOR", "MY_DB.MY_SCHEMA.MY_CONNECTOR"},
		{"SHOW RULES IN EVENT ROUTING TABLE (db.s.ert)", "RULES", "EVENT ROUTING TABLE", "DB.S.ERT"},
		{"SHOW CHANNELS IN PIPE db.s.p", "CHANNELS", "PIPE", "DB.S.P"},
		{"SHOW CHANNELS IN PIPE", "CHANNELS", "PIPE", ""},
		{"SHOW PARAMETERS FOR PIPE p", "PARAMETERS", "PIPE", "P"},
	} {
		stmt := parseOneShow(t, tt.sql)
		if stmt.ObjectClass != tt.class || stmt.Scope != ast.ShowScopeOther || stmt.ScopeText != tt.scope {
			t.Errorf("%q: class=%q scope=%v %q", tt.sql, stmt.ObjectClass, stmt.Scope, stmt.ScopeText)
		}
		got := ""
		if stmt.ScopeName != nil {
			got = stmt.ScopeName.Normalize()
		}
		if got != tt.name {
			t.Errorf("%q: scope name = %q, want %q", tt.sql, got, tt.name)
		}
	}

	// Without a name after the keywords, the word is a schema as before.
	if s := parseOneShow(t, "SHOW ROLES IN CLASS"); s.Scope != ast.ShowScopeSchema {
		t.Errorf("IN CLASS alone: scope = %v", s.Scope)
	}

	for _, sql := range []string{
		// The account takes one part.
		"SHOW FAILOVER GROUPS IN ACCOUNT org.acct",
		// The scopes belong to their class.
		"SHOW TABLES IN GIT REPOSITORY repo",
		"SHOW TABLES IN OPENFLOW CONNECTOR c",
		"SHOW TABLES IN EVENT ROUTING TABLE (t)",
		"SHOW DATABASES IN ACCOUNT a",
		// The event routing table's name is parenthesized.
		"SHOW RULES IN EVENT ROUTING TABLE (t",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}

// TestShowClassClauses covers the clauses one class's page writes after its
// scope, and SHOW SERVICES' OF TYPE list.
func TestShowClassClauses(t *testing.T) {
	for _, tt := range []struct {
		sql, qualifier, name string
	}{
		{"SHOW SEMANTIC DIMENSIONS LIKE 'c%' IN tpch_rev_analysis FOR METRIC order_average_value STARTS WITH 'C' LIMIT 3", "FOR METRIC", "ORDER_AVERAGE_VALUE"},
		{"SHOW SEMANTIC DIMENSIONS IN SCHEMA.a FOR METRIC m", "FOR METRIC", "M"},
		{"SHOW SEMANTIC DIMENSIONS IN DATABASE.a FOR METRIC m", "FOR METRIC", "M"},
		{"SHOW VERSIONS LIKE 'VERSION$%' IN ARTIFACT REPOSITORY my_db.my_schema.my_app_repo FOR PACKAGE my_package LIMIT 10", "FOR PACKAGE", "MY_PACKAGE"},
		{"SHOW RELEASE DIRECTIVES IN APPLICATION PACKAGE pkg FOR RELEASE CHANNEL qa", "FOR RELEASE CHANNEL", "QA"},
		{"SHOW SHARED CONTENT IN APPLICATION PACKAGE decl_share_app_pkg FOR VERSION VERSION$2", "FOR VERSION", "VERSION$2"},
		{"SHOW FUNCTIONS IN MODEL db.s.m VERSION v1", "VERSION", "V1"},
		{"SHOW RUN METRICS IN EXPERIMENT e RUN r LIMIT 2", "RUN", "R"},
		{"SHOW RUN PARAMETERS LIKE 'x%' IN EXPERIMENT e RUN r", "RUN", "R"},
		{"SHOW REPLICATION DATABASES WITH PRIMARY myorg.account1.mydb1", "WITH PRIMARY", "MYORG.ACCOUNT1.MYDB1"},
		{"SHOW REPLICATION DATABASES LIKE 'mydb%' WITH PRIMARY account1.mydb1", "WITH PRIMARY", "ACCOUNT1.MYDB1"},
	} {
		stmt := parseOneShow(t, tt.sql)
		if stmt.Qualifier != tt.qualifier || stmt.QualifierName == nil || stmt.QualifierName.Normalize() != tt.name {
			t.Errorf("%q: qualifier = %q %v", tt.sql, stmt.Qualifier, stmt.QualifierName)
		}
	}

	// The semantic pages' IN choice takes a <semantic_view_name>, a QNAME,
	// so SCHEMA.a and DATABASE.a are names rather than scope keywords.
	for _, tt := range []struct{ sql, name string }{
		{"SHOW SEMANTIC FACTS IN DATABASE.a", "DATABASE.A"},
		{"SHOW SEMANTIC METRICS IN SCHEMA.a", "SCHEMA.A"},
		{"SHOW SEMANTIC METRICS IN SCHEMA.a STARTS WITH 'C' LIMIT 2", "SCHEMA.A"},
	} {
		stmt := parseOneShow(t, tt.sql)
		if stmt.Scope != ast.ShowScopeSchema || stmt.ScopeName == nil || stmt.ScopeName.Normalize() != tt.name {
			t.Errorf("%q: scope=%v name=%v", tt.sql, stmt.Scope, stmt.ScopeName)
		}
	}

	s := parseOneShow(t, "SHOW JOB SERVICES EXCLUDE JOBS IN COMPUTE POOL p LIMIT 5 OF TYPE a, b")
	if len(s.OfTypes) != 2 || s.OfTypes[1].Normalize() != "B" || s.ObjectClass != "JOB SERVICES EXCLUDE JOBS" {
		t.Errorf("OF TYPE: class=%q types=%v", s.ObjectClass, s.OfTypes)
	}
	s = parseOneShow(t, "SHOW GRANTS TO SHARE s IN APPLICATION PACKAGE pkg LIMIT 3")
	if s.ScopeText != "APPLICATION PACKAGE" || s.ScopeName.Normalize() != "PKG" || !s.HasLimit {
		t.Errorf("TO SHARE IN APPLICATION PACKAGE: %+v", s)
	}

	for _, sql := range []string{
		// FOR PACKAGE and FOR VERSION are not optional.
		"SHOW VERSIONS IN ARTIFACT REPOSITORY r",
		"SHOW SHARED CONTENT IN APPLICATION PACKAGE p",
		// FOR METRIC follows IN <semantic_view_name> only.
		"SHOW SEMANTIC DIMENSIONS IN SCHEMA s FOR METRIC m",
		// Only the semantic pages take a QNAME as their IN choice; the
		// others write the SCHEMA scope and a name.
		"SHOW AGENTS IN SCHEMA.a",
		"SHOW SEMANTIC DIMENSIONS FOR METRIC m",
		// Each clause belongs to its class and scope.
		"SHOW TABLES IN APPLICATION PACKAGE p FOR VERSION v",
		"SHOW FUNCTIONS IN SCHEMA s VERSION v",
		"SHOW RUN METRICS IN SCHEMA s RUN r",
		"SHOW DATABASES WITH PRIMARY a.b",
		"SHOW TABLES OF TYPE a",
		"SHOW GRANTS TO ROLE r IN APPLICATION PACKAGE p",
		// Names: one part where the grammar writes ID; WITH PRIMARY has two or three.
		"SHOW RUN METRICS IN EXPERIMENT e RUN s.r",
		"SHOW REPLICATION DATABASES WITH PRIMARY mydb1",
		// OF TYPE takes a list without a trailing comma.
		"SHOW SERVICES OF TYPE a,",
		"SHOW SERVICES OF TYPE",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}
