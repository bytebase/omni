package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseCreateDataProtection(t *testing.T) {
	for _, tc := range []struct {
		in, kind string
		clauses  []string
	}{
		{"CREATE BACKUP SET IF NOT EXISTS s FOR DYNAMIC ICEBERG TABLE db.sc.t WITH TAG (a = 's') WITH BACKUP POLICY p COMMENT = 'c'",
			"BACKUP SET", []string{"FOR DYNAMIC ICEBERG TABLE", "", "WITH BACKUP POLICY", ""}},
		{"CREATE OR REPLACE BACKUP SET s FOR SCHEMA db.sc", "BACKUP SET", []string{"FOR SCHEMA"}},
		{"CREATE OR ALTER BACKUP SET s FOR DATABASE db WITH BACKUP POLICY p", "BACKUP SET", []string{"FOR DATABASE", "WITH BACKUP POLICY"}},
		{"CREATE BACKUP SET s FOR TABLE t", "BACKUP SET", []string{"FOR TABLE"}},
		{"CREATE SNAPSHOT SET s FOR DYNAMIC TABLE t TAG (a = 's', b = 's') WITH SNAPSHOT POLICY p", "SNAPSHOT SET",
			[]string{"FOR DYNAMIC TABLE", "", "WITH SNAPSHOT POLICY"}},
		{"CREATE BACKUP POLICY p WITH TAG (a = 's') WITH RETENTION LOCK SCHEDULE = '60 MINUTE' EXPIRE_AFTER_DAYS = 7 COMMENT = 'c'",
			"BACKUP POLICY", []string{"", "WITH RETENTION LOCK", ""}},
		{"CREATE OR ALTER BACKUP POLICY p SCHEDULE = 'USING CRON 0 0 * * * UTC'", "BACKUP POLICY", []string{""}},
		{"CREATE BACKUP POLICY p", "BACKUP POLICY", nil},
		{"CREATE OR REPLACE SNAPSHOT POLICY p WITH RETENTION LOCK", "SNAPSHOT POLICY", []string{"WITH RETENTION LOCK"}},
		{`CREATE SNAPSHOT IF NOT EXISTS s FROM SERVICE db.sc.svc VOLUME "data" INSTANCE 0 COMMENT = 'c' WITH TAG (a = 's')`,
			"SNAPSHOT", []string{"FROM SERVICE", "VOLUME", "INSTANCE", ""}},
		{"CREATE DATA MOVEMENT POLICY IF NOT EXISTS p ENFORCE_RULES = (r1, r2) ALERT_RULES = (r3) COMMENT = 'c'",
			"DATA MOVEMENT POLICY", []string{""}},
		{"CREATE OR REPLACE DATA MOVEMENT POLICY p", "DATA MOVEMENT POLICY", nil},
	} {
		s := mustObjectStmt(t, tc.in)
		if s.Verb != "CREATE" || s.Kind != tc.kind || strings.Join(clauseKeywords(s), "|") != strings.Join(tc.clauses, "|") {
			t.Errorf("%q: kind %q, clauses %q; want %q, %q", tc.in, s.Kind, clauseKeywords(s), tc.kind, tc.clauses)
		}
	}
	s := mustObjectStmt(t, `CREATE SNAPSHOT s FROM SERVICE svc VOLUME "data" INSTANCE 1`)
	if c := s.Clauses[2]; len(c.Options) != 1 || c.Options[0].Lit.Value != "1" || s.Clauses[1].Names[0].Name.Name != "data" {
		t.Errorf("clauses = %+v", s.Clauses)
	}
	s = mustObjectStmt(t, "CREATE DATA MOVEMENT POLICY p ENFORCE_RULES = (r1, r2)")
	if g := s.Clauses[0].Options[0].Group; len(g) != 2 || !g[0].Bare || g[1].Name != "R2" {
		t.Errorf("ENFORCE_RULES = %+v", g)
	}
}

func TestParseAlterDataProtection(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER BACKUP SET s ADD BACKUP":                                       "ADD BACKUP",
		"ALTER BACKUP SET s APPLY BACKUP POLICY p FORCE":                      "APPLY BACKUP POLICY|FORCE",
		"ALTER BACKUP SET s SUSPEND BACKUP POLICY":                            "SUSPEND BACKUP POLICY",
		"ALTER BACKUP SET s resume backup expiration policy":                  "RESUME BACKUP EXPIRATION POLICY",
		"ALTER BACKUP SET s DELETE BACKUP IDENTIFIER 'id'":                    "DELETE BACKUP IDENTIFIER",
		"ALTER BACKUP SET s MODIFY BACKUP IDENTIFIER 'id' ADD LEGAL HOLD":     "MODIFY BACKUP IDENTIFIER|ADD LEGAL HOLD",
		"ALTER BACKUP SET s MODIFY BACKUP IDENTIFIER 'id' REMOVE LEGAL HOLD":  "MODIFY BACKUP IDENTIFIER|REMOVE LEGAL HOLD",
		"ALTER BACKUP SET s MODIFY BACKUP IDENTIFIER 'id' SET COMMENT = 'c'":  "MODIFY BACKUP IDENTIFIER|SET",
		"ALTER BACKUP SET s MODIFY BACKUP IDENTIFIER 'id' UNSET COMMENT":      "MODIFY BACKUP IDENTIFIER|UNSET",
		"ALTER BACKUP SET s RENAME TO t":                                      "RENAME TO",
		"ALTER BACKUP SET s SET TAG a = 's', b = 's'":                         "SET TAG",
		"ALTER BACKUP SET s UNSET TAG a, b":                                   "UNSET TAG",
		"ALTER SNAPSHOT SET s APPLY SNAPSHOT POLICY p":                        "APPLY SNAPSHOT POLICY",
		"ALTER SNAPSHOT SET s SUSPEND SNAPSHOT CREATION POLICY":               "SUSPEND SNAPSHOT CREATION POLICY",
		"ALTER SNAPSHOT SET s MODIFY SNAPSHOT IDENTIFIER 'id' ADD LEGAL HOLD": "MODIFY SNAPSHOT IDENTIFIER|ADD LEGAL HOLD",
		"ALTER SNAPSHOT SET s SET COMMENT = 'c'":                              "SET",
		"ALTER BACKUP POLICY p SET COMMENT = 'c' SCHEDULE = '1 MINUTE'":       "SET",
		"ALTER BACKUP POLICY p UNSET EXPIRE_AFTER_DAYS":                       "UNSET",
		"ALTER SNAPSHOT POLICY p RENAME TO q":                                 "RENAME TO",
		"ALTER SNAPSHOT IF EXISTS s SET COMMENT = 'c'":                        "SET",
		"ALTER STORAGE LIFECYCLE POLICY IF EXISTS p SET ARCHIVE_TIER = COOL":  "SET",
		"ALTER STORAGE LIFECYCLE POLICY p UNSET ARCHIVE_FOR_DAYS":             "UNSET",
		"ALTER STORAGE LIFECYCLE POLICY p SET TAG a = 's'":                    "SET TAG",
		"ALTER DATA MOVEMENT POLICY p ADD ENFORCE_RULES = (r1, r2)":           "ADD",
		"ALTER DATA MOVEMENT POLICY p REMOVE ALERT_RULES = (r1)":              "REMOVE",
		"ALTER DATA MOVEMENT POLICY IF EXISTS p SET ALERT_RULES = (r1)":       "SET",
		"ALTER DATA MOVEMENT POLICY p RENAME TO q":                            "RENAME TO",
		"ALTER DATA MOVEMENT RULE IF EXISTS r UNSET COMMENT":                  "UNSET",
		"ALTER DATA MOVEMENT RULE r SET COMMENT = 'c'":                        "SET",
	} {
		s := mustObjectStmt(t, in)
		if got := strings.Join(clauseKeywords(s), "|"); s.Verb != "ALTER" || got != want {
			t.Errorf("%q: clauses = %q, want %q", in, got, want)
		}
	}
	s := mustObjectStmt(t, "ALTER BACKUP SET s DELETE BACKUP IDENTIFIER 'abc'")
	if o := s.Clauses[0].Options; len(o) != 1 || o[0].Name != "IDENTIFIER" || o[0].Lit.Value != "abc" {
		t.Errorf("options = %+v", o)
	}
}

func TestParseDropDataProtection(t *testing.T) {
	for in, want := range map[string]string{
		"DROP BACKUP SET s":                         "BACKUP SET",
		"DROP BACKUP POLICY p":                      "BACKUP POLICY",
		"DROP SNAPSHOT SET s":                       "SNAPSHOT SET",
		"DROP SNAPSHOT POLICY p":                    "SNAPSHOT POLICY",
		"DROP SNAPSHOT IF EXISTS s":                 "SNAPSHOT",
		"DROP STORAGE LIFECYCLE POLICY IF EXISTS p": "STORAGE LIFECYCLE POLICY",
		"DROP DATA MOVEMENT POLICY p":               "DATA MOVEMENT POLICY",
		"DROP DATA MOVEMENT RULE IF EXISTS r":       "DATA MOVEMENT RULE",
	} {
		d, ok := mustParseOne(t, in).(*ast.DropStmt)
		if !ok || d.Kind != ast.DropObject || d.Object != want {
			t.Errorf("%q: got %+v, want DropObject %q", in, d, want)
		}
	}
}

// The registered BACKUP / SNAPSHOT POLICY kinds win over the policy dispatch
// (startsPolicyKeyword), which still reads its own kinds.
func TestDataProtectionPolicyDispatch(t *testing.T) {
	for _, in := range []string{"CREATE BACKUP POLICY p", "ALTER SNAPSHOT POLICY p UNSET COMMENT"} {
		mustObjectStmt(t, in)
	}
	if _, ok := mustParseOne(t, "CREATE MASKING POLICY m AS (v STRING) RETURNS STRING -> v").(*ast.CreatePolicyStmt); !ok {
		t.Error("CREATE MASKING POLICY is not a CreatePolicyStmt")
	}
	if _, ok := mustParseOne(t, "ALTER SESSION POLICY p UNSET COMMENT").(*ast.AlterPolicyStmt); !ok {
		t.Error("ALTER SESSION POLICY is not an AlterPolicyStmt")
	}
}

func TestParseShowInDataProtectionSet(t *testing.T) {
	for in, want := range map[string]string{
		"SHOW BACKUPS IN BACKUP SET db.sc.s LIMIT 10": "BACKUP SET",
		"SHOW SNAPSHOTS IN SNAPSHOT SET s":            "SNAPSHOT SET",
	} {
		s, ok := mustParseOne(t, in).(*ast.ShowStmt)
		if !ok || s.Scope != ast.ShowScopeOther || s.ScopeText != want || s.ScopeName == nil {
			t.Errorf("%q: got %+v", in, s)
		}
	}
	for _, in := range []string{
		"SHOW BACKUPS IN BACKUP SET",
		"SHOW TABLES IN BACKUP SET s",
		"SHOW BACKUPS IN SNAPSHOT SET s",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseAlterAccountDataMovementPolicy(t *testing.T) {
	s := mustAlterAccount(t, "ALTER ACCOUNT SET DATA MOVEMENT POLICY p FORCE")
	if s.Action != ast.AlterAccountSetPolicy || s.PolicyKind != "DATA MOVEMENT POLICY" || s.PolicyName.String() != "p" || !s.Force {
		t.Errorf("got %+v", s)
	}
	s = mustAlterAccount(t, "ALTER ACCOUNT UNSET DATA MOVEMENT POLICY")
	if s.Action != ast.AlterAccountUnsetPolicy || s.PolicyKind != "DATA MOVEMENT POLICY" {
		t.Errorf("got %+v", s)
	}
}

func TestParseDataProtectionRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE BACKUP SET s",                                    // FOR is required
		"CREATE BACKUP SET s FOR VIEW v",                         // not documented
		"CREATE SNAPSHOT SET s FOR ICEBERG TABLE t",              // ICEBERG is BACKUP SET's
		"CREATE OR ALTER SNAPSHOT SET s FOR TABLE t",             // not documented
		"CREATE OR ALTER BACKUP SET s FOR TABLE t TAG (a = 's')", // no tags on CREATE OR ALTER
		"CREATE BACKUP SET s FOR TABLE t WITH BACKUP POLICY p WITH BACKUP POLICY q",
		"CREATE BACKUP SET s FOR TABLE t WITH SNAPSHOT POLICY p",
		"CREATE BACKUP POLICY p WITH RETENTION",
		"CREATE BACKUP POLICY p WITH RETENTION LOCK WITH RETENTION LOCK",
		"CREATE OR ALTER BACKUP POLICY p WITH TAG (a = 's')",
		"CREATE SNAPSHOT s FROM SERVICE svc VOLUME v INSTANCE a", // the instance is an index
		"CREATE SNAPSHOT s FROM SERVICE svc VOLUME a.v INSTANCE 0",
		"CREATE SNAPSHOT s FROM SERVICE svc INSTANCE 0",
		"CREATE SNAPSHOT s VOLUME v INSTANCE 0",
		"CREATE DATA MOVEMENT POLICY p WITH TAG (a = 's')",
		"CREATE DATA MOVEMENT POLICY p ENFORCE_RULES = ()",
		"CREATE DATA MOVEMENT POLICY p ENFORCE_RULES = ('r')",
		"ALTER BACKUP SET s",
		"ALTER BACKUP SET s ADD SNAPSHOT",
		"ALTER BACKUP SET s SUSPEND BACKUP",
		"ALTER BACKUP SET s DELETE BACKUP IDENTIFIER id",
		"ALTER BACKUP SET s MODIFY BACKUP IDENTIFIER 'id'",
		"ALTER BACKUP SET s MODIFY BACKUP IDENTIFIER 'id' ADD LEGAL",
		"ALTER BACKUP SET s MODIFY BACKUP IDENTIFIER 'id' SET TAG a = 's'",
		"ALTER SNAPSHOT SET s RENAME TO t", // not documented
		"ALTER SNAPSHOT s RENAME TO t",
		"ALTER SNAPSHOT s SET TAG a = 's'",
		"ALTER SNAPSHOT s SET COMMENT = 'c' TAG a = 's'",
		"ALTER STORAGE LIFECYCLE POLICY p SET BODY -> 1", // not read yet
		"ALTER DATA MOVEMENT POLICY p ADD COMMENT = 'c'",
		"ALTER DATA MOVEMENT POLICY p SET TAG a = 's'",
		"ALTER DATA MOVEMENT POLICY p SET ENFORCE_RULES = ()",
		"ALTER DATA MOVEMENT RULE r SET TAG a = 's'",
		"CREATE STORAGE LIFECYCLE POLICY p AS (a INT) RETURNS BOOLEAN -> a > 1", // not read yet
		"DROP BACKUP SET s CASCADE",
		// The pages write no IF EXISTS (dropNoIfExists).
		"DROP BACKUP SET IF EXISTS s",
		"DROP BACKUP POLICY IF EXISTS p",
		"DROP SNAPSHOT SET IF EXISTS s",
		"DROP SNAPSHOT POLICY IF EXISTS p",
		// Property names and value shapes are the pages'.
		"CREATE BACKUP SET s FOR TABLE t COMMENTS = 'c'",
		"CREATE BACKUP SET s FOR TABLE t COMMENT = 1",
		"CREATE BACKUP POLICY p SCHEDULER = 's'",
		"CREATE BACKUP POLICY p EXPIRE_AFTER_DAYS = 'x'",
		"CREATE SNAPSHOT POLICY p SCHEDULE = 60",
		"CREATE SNAPSHOT s FROM SERVICE svc VOLUME v INSTANCE 0 COMMENT = c",
		"CREATE DATA MOVEMENT POLICY p ALGORITHM = (a)",
		"ALTER BACKUP POLICY p SET EXPIRE_ROTATED_KEY_PAIR_AFTER_HOURS = 1",
		"ALTER BACKUP POLICY p UNSET SCHEDULER",
		"ALTER BACKUP SET s SET a = 's'",
		"ALTER BACKUP SET s UNSET SCHEDULE",
		"ALTER BACKUP SET s MODIFY BACKUP IDENTIFIER 'id' SET EXPIRE_AFTER_DAYS = 1",
		"ALTER SNAPSHOT SET s UNSET a",
		"ALTER SNAPSHOT s SET COMMENT = 1",
		"ALTER STORAGE LIFECYCLE POLICY p SET ARCHIVE_TIER = COPY",
		"ALTER STORAGE LIFECYCLE POLICY p SET ARCHIVE_TIER = 'COOL'",
		"ALTER STORAGE LIFECYCLE POLICY p UNSET ARCHIVE_TIER",
		"ALTER DATA MOVEMENT POLICY p UNSET ENFORCE_RULES",
		"ALTER DATA MOVEMENT POLICY p ADD ALERT_RULES = ('r')",
		"ALTER DATA MOVEMENT RULE r SET MAX_ROWS = 1",
	} {
		mustRejectObject(t, in)
	}
}
