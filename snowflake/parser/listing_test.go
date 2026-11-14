package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// wantObjectClauses parses in as an ObjectStmt and checks its clauses'
// keywords.
func wantObjectClauses(t *testing.T, in string, want ...string) *ast.ObjectStmt {
	t.Helper()
	s := mustObjectStmt(t, in)
	got := clauseKeywords(s)
	if len(got) != len(want) {
		t.Errorf("%q: clauses = %q, want %q", in, got, want)
		return s
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%q: clauses = %q, want %q", in, got, want)
			break
		}
	}
	return s
}

func TestParseListings(t *testing.T) {
	wantObjectClauses(t, "CREATE EXTERNAL LISTING IF NOT EXISTS l SHARE s AS 'm' PUBLISH = FALSE REVIEW = TRUE COMMENT = 'c'",
		"SHARE", "AS", "")
	wantObjectClauses(t, "CREATE EXTERNAL LISTING l APPLICATION PACKAGE p FROM '@db.s.st/dir' REVIEW = FALSE",
		"APPLICATION PACKAGE", "FROM", "")
	wantObjectClauses(t, "CREATE EXTERNAL LISTING l AS $$title: x$$", "AS")
	wantObjectClauses(t, "CREATE ORGANIZATION LISTING l SHARE s AS 'm' PUBLISH = TRUE", "SHARE", "AS", "")
	wantObjectClauses(t, "CREATE ORGANIZATION LISTING IF NOT EXISTS l FROM '@st'", "FROM")
	wantObjectClauses(t, "ALTER LISTING l")
	wantObjectClauses(t, "ALTER LISTING IF EXISTS l PUBLISH", "PUBLISH")
	wantObjectClauses(t, "ALTER LISTING l UNPUBLISH", "UNPUBLISH")
	wantObjectClauses(t, "ALTER LISTING l REVIEW", "REVIEW")
	wantObjectClauses(t, "ALTER LISTING l AS 'm' PUBLISH = TRUE COMMENT = 'c'", "AS", "")
	wantObjectClauses(t, "ALTER LISTING l ADD VERSION v3 FROM @db.public.st/manifests COMMENT = 'c'", "ADD VERSION", "FROM", "")
	wantObjectClauses(t, "ALTER LISTING l ADD VERSION IF NOT EXISTS v3 FROM @st", "ADD VERSION IF NOT EXISTS", "FROM")
	wantObjectClauses(t, "ALTER LISTING l ADD VERSION FROM @st", "ADD VERSION", "FROM")
	s := wantObjectClauses(t, "ALTER LISTING l ADD VERSION v FROM @repo/branches/my-branch/ COMMENT = 'c'", "ADD VERSION", "FROM", "")
	if o := s.Clauses[1].Options[0]; o.Words != "@repo/branches/my-branch/" {
		t.Errorf("FROM = %q", o.Words)
	}
	wantObjectClauses(t, "ALTER LISTING IF EXISTS l ADD TARGETS $$targets: x$$", "ADD TARGETS")
	wantObjectClauses(t, "ALTER LISTING l REMOVE TARGETS 'm'", "REMOVE TARGETS")
	wantObjectClauses(t, "ALTER LISTING l RENAME TO m", "RENAME TO")
	wantObjectClauses(t, "ALTER LISTING l SET COMMENT = 'c'", "SET")

	d, ok := mustParseOne(t, "DROP LISTING l").(*ast.DropStmt)
	if !ok || d.Object != "LISTING" {
		t.Errorf("DROP LISTING: got %#v", d)
	}
	desc, ok := mustParseOne(t, "DESC LISTING l REVISION = DRAFT").(*ast.DescribeStmt)
	if !ok || desc.ObjectType != "LISTING" || desc.Name.String() != "l" || desc.Revision != "DRAFT" {
		t.Errorf("DESC LISTING: got %+v", desc)
	}
	show, ok := mustParseOne(t, "SHOW TERSE AVAILABLE LISTINGS LIMIT 5 IS_IMPORTED = TRUE IS_SHARED_WITH_ME = TRUE").(*ast.ShowStmt)
	if !ok || show.ObjectClass != "AVAILABLE LISTINGS" || len(show.TrueFilters) != 2 || show.TrueFilters[1] != "IS_SHARED_WITH_ME" {
		t.Errorf("SHOW AVAILABLE LISTINGS: got %+v", show)
	}
}

func TestParseListingRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE LISTING l AS 'm'",
		"CREATE OR REPLACE EXTERNAL LISTING l AS 'm'",
		"CREATE EXTERNAL LISTING l",
		"CREATE EXTERNAL LISTING db.l AS 'm'",
		"CREATE EXTERNAL LISTING l SHARE db.s AS 'm'",
		"CREATE EXTERNAL LISTING l AS m",
		"CREATE EXTERNAL LISTING l FROM @st",
		"CREATE EXTERNAL LISTING l FROM 'x' COMMENT = 'c'", // the FROM form documents no COMMENT
		"CREATE EXTERNAL LISTING l AS 'm' PUBLISH = 'yes'",
		"CREATE EXTERNAL LISTING l AS 'm' DISTRIBUTION = EXTERNAL",
		"CREATE EXTERNAL LISTING l AS 'm' WITH TAG (t = 'v')",
		"CREATE ORGANIZATION LISTING l AS 'm' REVIEW = TRUE",
		"CREATE ORGANIZATION LISTING l AS 'm' COMMENT = 'c'",
		"ALTER LISTING l PUBLISH UNPUBLISH",
		"ALTER LISTING l AS 'm' FOO = 1",
		"ALTER LISTING IF EXISTS l ADD VERSION v FROM @st",
		"ALTER LISTING l ADD VERSION v FROM st", // the page writes a stage location
		"ALTER LISTING l ADD VERSION v GARBAGE @st",
		"ALTER LISTING l ADD VERSION v FROM @",
		"ALTER LISTING l ADD VERSION v FROM @st REVIEW = TRUE",
		"ALTER LISTING l ADD VERSION IF NOT EXISTS FROM @st",
		"ALTER LISTING l ADD TARGETS t", // the manifest is a string
		"ALTER LISTING l RENAME TO db.m",
		"ALTER LISTING l SET PUBLISH = TRUE",
		"ALTER LISTING l SET COMMENT = 'c' COMMENT = 'd'",
		"ALTER LISTING l UNSET COMMENT",
		"DROP LISTING IF EXISTS l",
		"DESC LISTING l REVISION = LATEST",
		"DESC LISTING l REVISION =",
		"SHOW AVAILABLE LISTINGS IS_IMPORTED = FALSE",
		"SHOW AVAILABLE LISTINGS IS_SHARED_WITH_ME = TRUE IS_IMPORTED = TRUE",
		"SHOW AVAILABLE LISTINGS IS_IMPORTED = TRUE IS_IMPORTED = TRUE",
	} {
		mustRejectObject(t, in)
	}
}
