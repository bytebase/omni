package parser

import "testing"

// The FIRST-set predicates of first_sets.go, held to the token slices they
// are built from. The oracle tests that held those slices to PostgreSQL 17's
// grammar were deleted: Redshift runs no PostgreSQL 17, and they asserted
// PostgreSQL's FIRST sets, not Redshift's.

// TestIsConstTypenameStartMatchesSimpleLeadSet asserts the invariant that
// isConstTypenameStart accepts every token in simpleTypenameLeadTokens.
// This is the positive direction of the predicate's contract — by
// reusing simpleTypenameLeadSet, it must accept the entire slice.
func TestIsConstTypenameStartMatchesSimpleLeadSet(t *testing.T) {
	for _, tok := range simpleTypenameLeadTokens {
		name := ""
		for _, kw := range Keywords {
			if kw.Token == tok {
				name = kw.Name
				break
			}
		}
		if name == "" {
			t.Fatalf("token %d in simpleTypenameLeadTokens has no Keywords entry", tok)
		}
		p := &Parser{cur: Token{Type: tok, Str: name}}
		if !p.isConstTypenameStart() {
			t.Errorf("isConstTypenameStart must accept %q (in simpleTypenameLeadTokens)", name)
		}
	}
}

// TestIsConstTypenameStartRejectsAllOtherKeywords is the exhaustive
// negative test: for every keyword in the omni Keywords table that is
// NOT in simpleTypenameLeadSet, the predicate must reject it. This
// catches the "ColNameKeyword silently leaks into the type FIRST set"
// regression class that the spot-check version (renamed below) cannot.
//
// Why we need this: PG's AexprConst grammar position is ambiguous with
// `func_name Sconst`, so PG cannot tell us via oracle whether a given
// keyword is a real ConstTypename lead. The exhaustive in-Go sweep
// against the explicit slice is the only way to lock down the negative
// coverage.
func TestIsConstTypenameStartRejectsAllOtherKeywords(t *testing.T) {
	for _, kw := range Keywords {
		if simpleTypenameLeadSet[kw.Token] {
			continue
		}
		p := &Parser{cur: Token{Type: kw.Token, Str: kw.Name}}
		if p.isConstTypenameStart() {
			t.Errorf("isConstTypenameStart wrongly accepts %q (tok %d, category %d) — not in simpleTypenameLeadSet",
				kw.Name, kw.Token, kw.Category)
		}
	}
}

// TestIsConstTypenameStartRejectsNonTypeStarters is the documented
// spot-check kept alongside the exhaustive sweep above. Its value is
// readability — a future reader scanning this file sees concrete
// examples (IDENT/BETWEEN/EXISTS/...) of tokens that PG would route
// through `func_name Sconst` and that the predicate rejects.
func TestIsConstTypenameStartRejectsNonTypeStarters(t *testing.T) {
	rejects := []struct {
		tok  int
		name string
	}{
		{IDENT, "foo"},       // bare IDENT → func_name path in PG
		{BETWEEN, "between"}, // ColNameKeyword, not a type
		{EXISTS, "exists"},
		{COALESCE, "coalesce"},
		{GREATEST, "greatest"},
		{ROW, "row"},
		{VALUES, "values"},
		{POSITION, "position"},
		{TRIM, "trim"},
	}
	for _, r := range rejects {
		p := &Parser{cur: Token{Type: r.tok, Str: r.name}}
		if p.isConstTypenameStart() {
			t.Errorf("isConstTypenameStart must reject %q (tok %d)", r.name, r.tok)
		}
	}
}

// TestIsTableConstraintStartCoverage validates the FIRST set against
// the canonical 6 tokens. There is no PG oracle test because the
// production is only reachable through CREATE/ALTER TABLE where
// surrounding column-definition rules create grammar ambiguity.
func TestIsTableConstraintStartCoverage(t *testing.T) {
	accepts := []struct {
		tok  int
		name string
	}{
		{CONSTRAINT, "constraint"},
		{CHECK, "check"},
		{UNIQUE, "unique"},
		{PRIMARY, "primary"},
		{FOREIGN, "foreign"},
		{EXCLUDE, "exclude"},
	}
	for _, a := range accepts {
		p := &Parser{cur: Token{Type: a.tok, Str: a.name}}
		if !p.isTableConstraintStart() {
			t.Errorf("isTableConstraintStart must accept %q", a.name)
		}
	}

	// Negative spot-checks: tokens that look type/column-ish but are NOT
	// table-constraint starters.
	rejects := []struct {
		tok  int
		name string
	}{
		{IDENT, "my_column"},
		{INT_P, "int"},
		{SELECT, "select"},
		{LIKE, "like"}, // LIKE is a CREATE TABLE element but NOT a constraint
	}
	for _, r := range rejects {
		p := &Parser{cur: Token{Type: r.tok, Str: r.name}}
		if p.isTableConstraintStart() {
			t.Errorf("isTableConstraintStart must reject %q", r.name)
		}
	}
}

func TestIsRedshiftDocShowStart(t *testing.T) {
	newShowParser := func(sql string) *Parser {
		p := &Parser{lexer: NewLexer(sql), source: sql}
		p.advance() // SHOW
		p.advance() // token after SHOW
		return p
	}
	for _, sql := range []string{
		"SHOW USER LOCKOUT",
		"SHOW COLUMN GRANTS ON TABLE s.t",
		"SHOW CONSTRAINTS PRIMARY KEYS FROM TABLE s.t",
		"SHOW FUNCTIONS FROM SCHEMA s",
		"SHOW PARAMETERS OF FUNCTION s.f(INT)",
		"SHOW RLS POLICIES",
		"SHOW MASKING POLICIES",
		"SHOW PROCEDURES FROM SCHEMA s",
		"SHOW TEMPLATE t",
		"SHOW TEMPLATES FROM SCHEMA s",
	} {
		if !newShowParser(sql).isRedshiftDocShowStart() {
			t.Errorf("isRedshiftDocShowStart(%q) = false, want true", sql)
		}
	}
	for _, sql := range []string{
		"SHOW ALL",
		"SHOW search_path",
		"SHOW rls",
		"SHOW template",
		"SHOW template;",
		"SHOW parameters.foo",
		`SHOW "RLS" POLICIES`,
		`SHOW USER "LOCKOUT"`,
		`SHOW "TEMPLATE" t`,
		"SHOW DATABASES",
	} {
		if newShowParser(sql).isRedshiftDocShowStart() {
			t.Errorf("isRedshiftDocShowStart(%q) = true, want false", sql)
		}
	}
}

func TestIsRedshiftTemplateLead(t *testing.T) {
	for _, kw := range Keywords {
		p := &Parser{}
		tok := Token{Type: kw.Token, Str: kw.Name}
		if got, want := p.isRedshiftTemplateLead(tok), kw.Token == TEMPLATE; got != want {
			t.Errorf("isRedshiftTemplateLead(%q) = %v, want %v", kw.Name, got, want)
		}
	}
	p := &Parser{lexer: NewLexer(`"template"`), source: `"template"`}
	p.advance()
	if p.isRedshiftTemplateLead(p.cur) {
		t.Error(`isRedshiftTemplateLead("template") = true, want false`)
	}
}

// TestIsRedshiftMultiwordTypenameStart holds the predicate to the multiword
// type map: every lead accepts each of its words, and nothing else is a
// multiword type.
func TestIsRedshiftMultiwordTypenameStart(t *testing.T) {
	for lead, words := range redshiftTypenameContinuations {
		for _, word := range words {
			p := &Parser{cur: Token{Type: lead}, nextBuf: Token{Type: word}, hasNext: true}
			if !p.isRedshiftMultiwordTypenameStart() {
				t.Errorf("isRedshiftMultiwordTypenameStart(%d, %d) = false, want true", lead, word)
			}
		}
	}
	for _, tc := range []struct {
		cur, next int
	}{
		{CHARACTER, IDENT},  // character a: one type, then extra input
		{IDENT, VARYING},    // x varying: an IDENT is no multiword type lead
		{IDENT, IDENT},      // a b
		{INT_P, PRECISION},  // INT takes no second word
		{CHARACTER, ';'},    // character;
		{BIT, PRECISION},    // BIT takes VARYING only
		{BINARY, PRECISION}, // BINARY takes VARYING only
	} {
		p := &Parser{cur: Token{Type: tc.cur}, nextBuf: Token{Type: tc.next}, hasNext: true}
		if p.isRedshiftMultiwordTypenameStart() {
			t.Errorf("isRedshiftMultiwordTypenameStart(%d, %d) = true, want false", tc.cur, tc.next)
		}
	}
}

// TestIsRedshiftDatashareNameStart holds the predicate to the tokens it reads:
// the unquoted word DATASHARE and a name after it.
func TestIsRedshiftDatashareNameStart(t *testing.T) {
	newGranteeParser := func(sql string) *Parser {
		p := &Parser{lexer: NewLexer(sql), source: sql}
		p.advance()
		return p
	}
	for _, sql := range []string{"datashare d", "DATASHARE sales_share"} {
		if !newGranteeParser(sql).isRedshiftDatashareNameStart() {
			t.Errorf("isRedshiftDatashareNameStart(%q) = false, want true", sql)
		}
	}
	for _, sql := range []string{"datashare", "datashare, d", `"datashare" d`, "role d", "a.b d"} {
		if newGranteeParser(sql).isRedshiftDatashareNameStart() {
			t.Errorf("isRedshiftDatashareNameStart(%q) = true, want false", sql)
		}
	}
}

// TestIsRedshiftFederatedPermissionStart holds the CONNECT and DEBUG leads to
// the tokens they read: the unquoted words, and what may follow them.
func TestIsRedshiftFederatedPermissionStart(t *testing.T) {
	newStmtParser := func(sql string) *Parser {
		p := &Parser{lexer: NewLexer(sql), source: sql}
		p.advance() // GRANT or REVOKE
		p.advance() // the permission
		return p
	}
	for _, sql := range []string{
		"GRANT CONNECT TO PUBLIC",
		"GRANT CONNECT ON WORKGROUP TO PUBLIC",
		"GRANT DEBUG ON DATABASE d TO a:b",
	} {
		p := newStmtParser(sql)
		if sql[6] == 'C' && !p.isRedshiftConnectStart() {
			t.Errorf("isRedshiftConnectStart(%q) = false, want true", sql)
		}
		if sql[6] == 'D' && !p.isRedshiftDebugStart() {
			t.Errorf("isRedshiftDebugStart(%q) = false, want true", sql)
		}
	}
	for _, sql := range []string{"GRANT CONNECT FROM a", "GRANT DEBUG TO a", `GRANT "CONNECT" TO a`} {
		p := newStmtParser(sql)
		if p.isRedshiftConnectStart() || p.isRedshiftDebugStart() {
			t.Errorf("%q: a federated permission lead, want none", sql)
		}
	}
}

// TestIsRedshiftIdentityGranteeStart holds the CONNECT and DEBUG grantee leads:
// ADMINS OF ACCOUNT, PUBLIC, and the ROLE marker of ROLE <prefix>:<name>.
func TestIsRedshiftIdentityGranteeStart(t *testing.T) {
	newGranteeParser := func(sql string) *Parser {
		p := &Parser{lexer: NewLexer(sql), source: sql}
		p.advance()
		return p
	}
	for _, sql := range []string{"ADMINS OF ACCOUNT '1'", "admins of account '1'"} {
		if !newGranteeParser(sql).isRedshiftAdminsOfAccountStart() {
			t.Errorf("isRedshiftAdminsOfAccountStart(%q) = false, want true", sql)
		}
	}
	for _, sql := range []string{"PUBLIC", "public"} {
		if !newGranteeParser(sql).isRedshiftPublicGranteeStart() {
			t.Errorf("isRedshiftPublicGranteeStart(%q) = false, want true", sql)
		}
	}
	for _, sql := range []string{"ROLE AWSIDC:analysts", "role a:b"} {
		if !newGranteeParser(sql).isRedshiftRoleIdentityStart() {
			t.Errorf("isRedshiftRoleIdentityStart(%q) = false, want true", sql)
		}
	}
	for _, tc := range []struct {
		sql      string
		admins   bool
		public   bool
		roleName bool
	}{
		{"ADMINS ACCOUNT '1'", false, false, false},
		{"PUBLIC:alice", false, false, false},
		{"ROLE:analysts", false, false, false},
		{`"PUBLIC"`, false, false, false},
		{"alice", false, false, false},
	} {
		p := newGranteeParser(tc.sql)
		if got := p.isRedshiftAdminsOfAccountStart(); got != tc.admins {
			t.Errorf("isRedshiftAdminsOfAccountStart(%q) = %v, want %v", tc.sql, got, tc.admins)
		}
		if got := p.isRedshiftPublicGranteeStart(); got != tc.public {
			t.Errorf("isRedshiftPublicGranteeStart(%q) = %v, want %v", tc.sql, got, tc.public)
		}
		if got := p.isRedshiftRoleIdentityStart(); got != tc.roleName {
			t.Errorf("isRedshiftRoleIdentityStart(%q) = %v, want %v", tc.sql, got, tc.roleName)
		}
	}
}

// TestIsRedshiftNamedGranteeStart holds the ROLE and GROUP grantee leads: the
// keyword and a name after it, never the keyword alone.
func TestIsRedshiftNamedGranteeStart(t *testing.T) {
	newGranteeParser := func(sql string) *Parser {
		p := &Parser{lexer: NewLexer(sql), source: sql}
		p.advance()
		return p
	}
	for _, sql := range []string{"ROLE analyst", "role analyst"} {
		if !newGranteeParser(sql).isRedshiftRoleGranteeStart() {
			t.Errorf("isRedshiftRoleGranteeStart(%q) = false, want true", sql)
		}
	}
	for _, sql := range []string{"GROUP analysts", "group analysts"} {
		if !newGranteeParser(sql).isRedshiftGroupGranteeStart() {
			t.Errorf("isRedshiftGroupGranteeStart(%q) = false, want true", sql)
		}
	}
	for _, tc := range []struct {
		sql   string
		role  bool
		group bool
	}{
		{"ROLE", false, false},
		{"ROLE, analyst", false, false},
		{"GROUP", false, false},
		{"analyst", false, false},
		{`"ROLE" analyst`, false, false},
	} {
		p := newGranteeParser(tc.sql)
		if got := p.isRedshiftRoleGranteeStart(); got != tc.role {
			t.Errorf("isRedshiftRoleGranteeStart(%q) = %v, want %v", tc.sql, got, tc.role)
		}
		if got := p.isRedshiftGroupGranteeStart(); got != tc.group {
			t.Errorf("isRedshiftGroupGranteeStart(%q) = %v, want %v", tc.sql, got, tc.group)
		}
	}
}

// TestIsRedshiftScopedPermissionStart holds the scoped-permission lead to the
// clause it reads: FOR straight after the first privilege.
func TestIsRedshiftScopedPermissionStart(t *testing.T) {
	newStmtParser := func(sql string) *Parser {
		p := &Parser{lexer: NewLexer(sql), source: sql}
		p.advance() // GRANT
		p.advance() // the privilege
		return p
	}
	for _, sql := range []string{"GRANT CREATE FOR SCHEMAS IN DATABASE d TO a", "GRANT SELECT FOR TABLES IN SCHEMA s TO a"} {
		if !newStmtParser(sql).isRedshiftScopedPermissionStart() {
			t.Errorf("isRedshiftScopedPermissionStart(%q) = false, want true", sql)
		}
	}
	for _, sql := range []string{"GRANT CREATE ON DATABASE d TO a", "GRANT SELECT ON t TO a", "GRANT CREATE, USAGE FOR SCHEMAS IN DATABASE d TO a"} {
		if newStmtParser(sql).isRedshiftScopedPermissionStart() {
			t.Errorf("isRedshiftScopedPermissionStart(%q) = true, want false", sql)
		}
	}
}
