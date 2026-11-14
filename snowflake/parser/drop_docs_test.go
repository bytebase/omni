package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseDocDrops(t *testing.T) {
	for in, want := range map[string]string{
		"DROP SECRET s":                            "SECRET",
		"DROP SECRET IF EXISTS db.sch.s":           "SECRET",
		"DROP PASSWORD POLICY IF EXISTS db.s.p":    "PASSWORD POLICY",
		"DROP AUTHENTICATION POLICY p":             "AUTHENTICATION POLICY",
		"DROP GIT REPOSITORY IF EXISTS db.s.r":     "GIT REPOSITORY",
		"DROP EXTERNAL VOLUME IF EXISTS v":         "EXTERNAL VOLUME",
		"DROP SEMANTIC VIEW IF EXISTS db.s.v":      "SEMANTIC VIEW",
		"DROP ACCOUNT a GRACE_PERIOD_IN_DAYS = 14": "ACCOUNT",
	} {
		d, ok := mustParseOne(t, in).(*ast.DropStmt)
		if !ok || d.Kind != ast.DropObject || d.Object != want {
			t.Errorf("%q: got %+v", in, d)
			continue
		}
		if d.Loc.Start != 0 || d.Loc.End != len(in) {
			t.Errorf("%q: Loc = %+v", in, d.Loc)
		}
	}
	d := mustParseOne(t, "DROP ACCOUNT IF EXISTS a GRACE_PERIOD_IN_DAYS = 3").(*ast.DropStmt)
	if !d.IfExists || d.GracePeriodInDays == nil || d.GracePeriodInDays.Ival != 3 {
		t.Errorf("DROP ACCOUNT = %+v", d)
	}
	d = mustParseOne(t, "DROP SERVICE IF EXISTS s FORCE").(*ast.DropStmt)
	if d.Object != "SERVICE" || !d.IfExists || !d.Force || d.Loc.End != len("DROP SERVICE IF EXISTS s FORCE") {
		t.Errorf("DROP SERVICE FORCE = %+v", d)
	}
	d = mustParseOne(t, "DROP DYNAMIC TABLE IF EXISTS db.s.dt").(*ast.DropStmt)
	if d.Kind != ast.DropDynamicTable || !d.IfExists {
		t.Errorf("DROP DYNAMIC TABLE = %+v", d)
	}

	for _, in := range []string{
		"DROP ACCOUNT a", // GRACE_PERIOD_IN_DAYS is required
		"DROP ACCOUNT a GRACE_PERIOD_IN_DAYS = 'x'", // an integer
		"DROP ACCOUNT a GRACE_PERIOD_IN_DAYS 14",
		"DROP ACCOUNT db.a GRACE_PERIOD_IN_DAYS = 14", // one part
		"DROP EXTERNAL VOLUME db.v",                   // one part
		"DROP SECRET s CASCADE",
		"DROP SEMANTIC VIEW v RESTRICT",
		"DROP PASSWORD POLICY",
		"DROP 'SECRET' s", // a string is not a keyword
		"DROP SERVICE s FORCE FORCE",
		"DROP IMAGE REPOSITORY r FORCE", // FORCE is SERVICE's
	} {
		mustRejectObject(t, in)
	}
}

func TestStartsDocDrop(t *testing.T) {
	for in, want := range map[string]bool{
		"SECRET s":          true,
		"GIT REPOSITORY r":  true,
		"EXTERNAL VOLUME v": true,
		"EXTERNAL TABLE t":  false,
		"SEMANTIC v":        false,
		`"SECRET" s`:        false,
		"AUTHENTICATION p":  false,
		"PASSWORD POLICY p": true,
		"ACCOUNT a":         true,
		"MANAGED ACCOUNT a": false,
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		if got := p.startsDocDrop() != nil; got != want {
			t.Errorf("%q: got %v, want %v", in, got, want)
		}
	}
}
