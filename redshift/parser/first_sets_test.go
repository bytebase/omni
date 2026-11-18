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
