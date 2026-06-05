package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-ddl-bigquery node generic-entity path: CREATE / ALTER
// / DROP CAPACITY / RESERVATION / ASSIGNMENT (DDL-024/025/026/053). The legacy
// GoogleSQLParser.g4 has NO dedicated rule for these — they parse via the generic
// entity mechanism (generic_entity_type: IDENTIFIER | PROJECT), so CAPACITY etc.
// must lex as bare identifiers (they are not reserved keywords). BigQuery-only at
// the union level (Spanner has no generic-entity mechanism; CREATE CAPACITY
// rejects, probed 2026-06-05).

func entityOf(t *testing.T, sql string) *ast.CreateEntityStmt {
	t.Helper()
	n := parseDDL(t, sql)
	e, ok := n.(*ast.CreateEntityStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateEntityStmt", sql, n)
	}
	return e
}

func TestCreateCapacity(t *testing.T) {
	// DDL-024.
	e := entityOf(t, "CREATE CAPACITY `admin_project.region-us.my-commitment` OPTIONS (slot_count = 100, plan = 'ANNUAL')")
	if e.EntityType != "CAPACITY" {
		t.Errorf("EntityType = %q, want CAPACITY", e.EntityType)
	}
	if e.Name.String() != "admin_project.region-us.my-commitment" {
		t.Errorf("Name = %q", e.Name.String())
	}
	if len(e.Options) != 2 {
		t.Errorf("Options = %d, want 2", len(e.Options))
	}
}

func TestCreateReservation(t *testing.T) {
	// DDL-025.
	e := entityOf(t, "CREATE RESERVATION `admin_project.region-us.prod` OPTIONS (slot_capacity = 100)")
	if e.EntityType != "RESERVATION" {
		t.Errorf("EntityType = %q, want RESERVATION", e.EntityType)
	}
}

func TestCreateAssignment(t *testing.T) {
	// DDL-026.
	e := entityOf(t, "CREATE ASSIGNMENT `admin_project.region-us.prod.my-assignment` OPTIONS (assignee = 'projects/my-project', job_type = 'QUERY')")
	if e.EntityType != "ASSIGNMENT" {
		t.Errorf("EntityType = %q, want ASSIGNMENT", e.EntityType)
	}
}

func TestCreateEntity_OrReplaceIfNotExistsBody(t *testing.T) {
	e := entityOf(t, "CREATE OR REPLACE RESERVATION IF NOT EXISTS `p.r` OPTIONS(x=1) AS JSON '{\"a\":1}'")
	if !e.OrReplace || !e.IfNotExists {
		t.Errorf("OrReplace=%v IfNotExists=%v, want both true", e.OrReplace, e.IfNotExists)
	}
	if e.BodyText == "" {
		t.Error("BodyText = empty, want the AS JSON body text")
	}
}

func TestCreateEntity_Rejects(t *testing.T) {
	cases := []string{
		"CREATE CAPACITY",                 // missing name
		"CREATE RESERVATION OPTIONS(x=1)", // OPTIONS where a name is required (OPTIONS is a keyword, not a path)
		// Regression (review finding): create_entity_statement has NO
		// opt_create_scope; a leading TEMP must reject (TEMP is not the entity type).
		"CREATE TEMP RESERVATION `p.r` OPTIONS(x=1)",
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- ALTER <generic-entity> ---

func TestAlterReservation_SetOptions(t *testing.T) {
	a := bqAlterOf(t, "ALTER RESERVATION `p.r` SET OPTIONS(slot_capacity=200)")
	if a.Object != ast.BQAlterEntity || a.EntityType != "RESERVATION" {
		t.Errorf("Object=%v EntityType=%q", a.Object, a.EntityType)
	}
	if len(a.SetOptions) != 1 {
		t.Errorf("SetOptions = %d, want 1", len(a.SetOptions))
	}
}

func TestAlterEntity_SetAs(t *testing.T) {
	a := bqAlterOf(t, "ALTER CAPACITY `p.c` SET AS JSON '{\"x\":1}'")
	if a.SetAsBody == "" {
		t.Error("SetAsBody = empty, want the JSON body")
	}
}

// --- DROP <generic-entity> ---

func TestDropCapacity(t *testing.T) {
	// DDL-053.
	d := bqDropOf(t, "DROP CAPACITY IF EXISTS `admin_project.region-us.my-commitment`")
	if d.Object != ast.BQDropEntity || d.EntityType != "CAPACITY" {
		t.Errorf("Object=%v EntityType=%q", d.Object, d.EntityType)
	}
	if !d.IfExists {
		t.Error("IfExists = false, want true")
	}
}

func TestDropReservationAssignment(t *testing.T) {
	// DDL-053.
	for _, tc := range []struct {
		sql  string
		want string
	}{
		{"DROP RESERVATION IF EXISTS `admin_project.region-us.prod`", "RESERVATION"},
		{"DROP ASSIGNMENT IF EXISTS `admin_project.region-us.prod.my-assignment`", "ASSIGNMENT"},
	} {
		d := bqDropOf(t, tc.sql)
		if d.EntityType != tc.want {
			t.Errorf("Parse(%q): EntityType = %q, want %q", tc.sql, d.EntityType, tc.want)
		}
	}
}
