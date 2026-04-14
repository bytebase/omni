package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type firstSetOracle struct {
	db  *sql.DB
	ctx context.Context
}

// oracleSetupError remembers a setup failure so every subsequent test
// observes it and can fail or skip uniformly. sync.Once runs exactly
// once, so errors must be captured here rather than returned from Do.
var oracleSetupError error

var (
	firstSetOracleOnce sync.Once
	firstSetOracleInst *firstSetOracle
)

func startFirstSetOracle(t *testing.T) *firstSetOracle {
	t.Helper()
	firstSetOracleOnce.Do(func() {
		ctx := context.Background()
		container, err := tcpg.Run(ctx, "postgres:17-alpine",
			tcpg.WithDatabase("omni_fs"),
			tcpg.WithUsername("postgres"),
			tcpg.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2)),
		)
		if err != nil {
			oracleSetupError = fmt.Errorf("container start: %w", err)
			return
		}
		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			oracleSetupError = fmt.Errorf("conn string: %w", err)
			return
		}
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			oracleSetupError = fmt.Errorf("db open: %w", err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			oracleSetupError = fmt.Errorf("ping: %w", err)
			return
		}
		firstSetOracleInst = &firstSetOracle{db: db, ctx: ctx}
	})

	if oracleSetupError != nil {
		// CI must fail loudly — otherwise a misconfigured CI silently
		// skips every FIRST-set test and disables the entire guardrail.
		// Local dev without docker is allowed to skip.
		if isCI() {
			t.Fatalf("first-set oracle unavailable in CI: %v", oracleSetupError)
		}
		t.Skipf("first-set oracle unavailable (local dev): %v", oracleSetupError)
	}
	return firstSetOracleInst
}

// isCI reports whether we're running under continuous integration.
// Matches omni's existing conventions and major CI providers.
func isCI() bool {
	// Respect the standard CI env var set by GitHub Actions, CircleCI,
	// GitLab CI, etc. If omni has a project-specific variable, add it here.
	return os.Getenv("CI") == "true" || os.Getenv("CI") == "1"
}

// probeResult classifies a PG error into accept/reject for FIRST-set purposes.
type probeResult int

const (
	// probeAccept means PG did NOT report a syntax error. This includes
	// outright success AND semantic errors like 42704 undefined_object,
	// because semantic errors prove the parser accepted the leading
	// tokens before reaching name resolution. For FIRST-set purposes,
	// "syntactically valid" is the question, not "semantically valid".
	probeAccept probeResult = iota
	probeReject                    // syntax_error (42601)
)

func classifyProbeErr(err error) probeResult {
	if err == nil {
		return probeAccept
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "42601" {
			return probeReject
		}
	}
	return probeAccept
}

// probe runs one SQL statement inside a transaction with a per-call
// timeout, then rolls back so side effects don't leak. The timeout
// prevents a pathological probe from hanging the entire test binary
// until the Go test timeout expires.
func (o *firstSetOracle) probe(t *testing.T, sqlStr string) probeResult {
	t.Helper()
	ctx, cancel := context.WithTimeout(o.ctx, 5*time.Second)
	defer cancel()
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx, sqlStr)
	return classifyProbeErr(err)
}

// candidateKind categorizes a probe token.
type candidateKind int

const (
	kindKeyword candidateKind = iota // From the Keywords table
	kindIdent                        // Bare identifier, not a keyword
	kindLiteral                      // Numeric / string / bit / hex literal
	kindPunct                        // Operator or punctuation opener
)

// candidateToken represents one token to probe in a FIRST-set test.
//
// For keyword candidates, `name` is both the SQL surface form (lowercase
// keyword) and the string used to construct Parser.cur.Str, and `token`
// is the omni token constant from keywords.go.
//
// For non-keyword candidates, `probeSQL` is the SQL fragment to
// interpolate into the probe template (may include trailing disambiguators
// like `"(1)"`, `"+1"`), and `token` is the omni token constant for the
// LEADING token of that fragment — e.g. the `'('` rune for `"(1)"`, the
// PARAM constant for `"$1"`. When the predicate is called, the test
// constructs Parser.cur = Token{Type: token, Str: leadStr} so the
// lead-token check reflects what the first lexed token actually is.
type candidateToken struct {
	name     string          // SQL keyword string (keyword kind) OR leading-token literal (non-keyword)
	probeSQL string          // SQL fragment interpolated into the probe template
	display  string          // Human-readable label for test output
	token    int             // omni token type constant — MUST be set for every candidate
	category KeywordCategory // Valid only when kind == kindKeyword
	kind     candidateKind
}

var nonKeywordCandidates = []candidateToken{
	{name: "omni_probe_ident", probeSQL: "omni_probe_ident", display: "IDENT",
		token: IDENT, kind: kindIdent},
	{name: "1", probeSQL: "1", display: "ICONST",
		token: ICONST, kind: kindLiteral},
	{name: "1.5", probeSQL: "1.5", display: "FCONST",
		token: FCONST, kind: kindLiteral},
	{name: "'x'", probeSQL: "'x'", display: "SCONST",
		token: SCONST, kind: kindLiteral},
	{name: "B'1'", probeSQL: "B'1'", display: "BCONST",
		token: BCONST, kind: kindLiteral},
	{name: "X'AB'", probeSQL: "X'AB'", display: "XCONST",
		token: XCONST, kind: kindLiteral},
	{name: "$1", probeSQL: "$1", display: "PARAM",
		token: PARAM, kind: kindLiteral},
	{name: "(", probeSQL: "(1)", display: "LPAREN",
		token: int('('), kind: kindPunct},
	{name: "+", probeSQL: "+1", display: "PLUS",
		token: int('+'), kind: kindPunct},
	{name: "-", probeSQL: "-1", display: "MINUS",
		token: int('-'), kind: kindPunct},
}

func allCandidates() []candidateToken {
	out := make([]candidateToken, 0, len(Keywords)+len(nonKeywordCandidates))
	for _, kw := range Keywords {
		out = append(out, candidateToken{
			name:     kw.Name,
			probeSQL: kw.Name,
			display:  kw.Name,
			token:    kw.Token,
			category: kw.Category,
			kind:     kindKeyword,
		})
	}
	out = append(out, nonKeywordCandidates...)
	return out
}

// renderFn produces one or more candidate renderings for a token.
type renderFn func(c candidateToken) []string

// renderBare returns only the candidate's configured probeSQL.
func renderBare(c candidateToken) []string { return []string{c.probeSQL} }

// renderTypeCandidate expands a token into the combinations that real
// PG type grammar requires: bare, parenthesized type modifiers, and the
// multi-word forms (PRECISION, CHARACTER, VARYING) that the PG lexer
// emits as separate tokens. Tokens that REQUIRE a trailing element
// (like SETOF) get dedicated renderings that supply a valid continuation.
func renderTypeCandidate(c candidateToken) []string {
	if c.kind != kindKeyword {
		return []string{c.probeSQL}
	}
	base := c.name
	if base == "setof" {
		return []string{"setof int", "setof varchar(10)", "setof json"}
	}
	return []string{
		base,
		base + "(1)",
		base + "(1, 1)",
		base + " precision",
		base + " character",
		base + " character(1)",
		base + " varying",
		base + " varying(1)",
		base + " character varying",
		base + " character varying(1)",
	}
}

// probeOutcome records the result for a single candidate.
type probeOutcome struct {
	candidate candidateToken
	accepted  bool
}

// runFirstSetProbe probes each filtered candidate via every rendering
// returned by `render`. A candidate is marked accepted if any rendering
// parses under PG. Each probe runs inside the savepoint isolation
// provided by (*firstSetOracle).probe.
func runFirstSetProbe(
	t *testing.T,
	o *firstSetOracle,
	probeTemplate string, // fmt template with a single %s slot
	render renderFn,
	filter func(candidateToken) bool,
) []probeOutcome {
	t.Helper()
	var outcomes []probeOutcome
	for _, c := range allCandidates() {
		if !filter(c) {
			continue
		}
		accepted := false
		for _, form := range render(c) {
			sqlStr := fmt.Sprintf(probeTemplate, form)
			if o.probe(t, sqlStr) == probeAccept {
				accepted = true
				break
			}
		}
		outcomes = append(outcomes, probeOutcome{candidate: c, accepted: accepted})
	}
	return outcomes
}

// Filters for common cases.
func filterKeywordsOnly(c candidateToken) bool { return c.kind == kindKeyword }
func filterAll(c candidateToken) bool          { return true }

func TestSimpleTypenameLeadTokensMatchPG(t *testing.T) {
	o := startFirstSetOracle(t)

	// Probe each keyword via CAST. CAST's grammar is
	//   CAST '(' a_expr AS Typename ')'
	// (gram.y:14130), which exercises `Typename`, a SUPERSET of
	// SimpleTypename: Typename = SimpleTypename opt_array_bounds
	//                          | SETOF SimpleTypename opt_array_bounds | ...
	//
	// For FIRST-set purposes the only difference is SETOF, since
	// opt_array_bounds is a suffix. We therefore subtract SETOF from
	// PG's accepted set below — the SimpleTypename production proper
	// does NOT include SETOF, so isSimpleTypenameStart must reject it,
	// but the CAST probe will accept `CAST(NULL AS setof int)` because
	// renderTypeCandidate produces that form.
	//
	// Result classification:
	//   42601 syntax_error → not a SimpleTypename lead
	//   42704 undefined_object / success → valid lead (unknown or built-in)
	//
	// Multi-token type syntax (DOUBLE PRECISION, NATIONAL CHARACTER,
	// CHARACTER VARYING, etc.) is covered by renderTypeCandidate trying
	// every rendering; the probe is accepted if any rendering parses.

	// Scope the probe to keyword + IDENT candidates. Including literals,
	// PARAM, and punctuation would only produce "both reject" matches
	// (PG rejects `CAST(NULL AS 1)`, predicate rejects int-typed tokens)
	// which is noise that could mask a real regression in the IDENT path.
	filter := func(c candidateToken) bool {
		return c.kind == kindKeyword || c.kind == kindIdent
	}
	outcomes := runFirstSetProbe(t, o,
		"SELECT CAST(NULL AS %s)",
		renderTypeCandidate,
		filter,
	)

	// Subtract SETOF: CAST probes Typename, which includes SETOF, but
	// SimpleTypename does not. Without this subtraction the test would
	// report "PG accepts but predicate rejects: setof" on every run.
	// The Typename test in Task 1.7 handles SETOF explicitly via its
	// own predicate isTypenameStart.
	for i := range outcomes {
		if outcomes[i].candidate.token == SETOF {
			outcomes[i].accepted = false
		}
	}

	assertPredicateMatchesPG(t, "SimpleTypename", outcomes, predicateProbe{
		onKeyword: func(c candidateToken) bool {
			// Str MUST be set — isTypeFunctionName's category check reads
			// p.cur.Str and cross-verifies via LookupKeyword.
			p := &Parser{cur: Token{Type: c.token, Str: c.name}}
			return p.isSimpleTypenameStart()
		},
		onNonKeyword: func(c candidateToken) bool {
			// Only IDENT reaches SimpleTypename (via GenericType →
			// isTypeFunctionName). Literals, punctuation, PARAM are
			// rejected by the predicate.
			p := &Parser{cur: Token{Type: c.token, Str: c.name}}
			return p.isSimpleTypenameStart()
		},
	})
}

// predicateProbe dispatches a FIRST-set predicate against a single
// candidate. Both callbacks receive the full candidateToken because
// omni's category predicates (isTypeFunctionName, isColId, isColLabel)
// consult BOTH the token type AND the literal string: `LookupKeyword(
// p.cur.Str)` must return a Keyword whose Token matches p.cur.Type. A
// predicateProbe that ignored candidateToken.name would make every
// UnreservedKeyword / TypeFuncNameKeyword probe fail the category check
// and report massive false drift.
//
// onNonKeyword may be nil to signal the production doesn't care about
// non-keyword starters — those candidates are skipped entirely in that
// case.
type predicateProbe struct {
	onKeyword    func(c candidateToken) bool
	onNonKeyword func(c candidateToken) bool // may be nil
}

// assertPredicateMatchesPG performs the full bidirectional FIRST-set check:
//
//	(1) PG accepts ∧ predicate rejects → missingFromOmni
//	(2) predicate accepts ∧ PG rejects → extraInOmni
//
// Both directions iterate ALL outcomes returned by runFirstSetProbe.
func assertPredicateMatchesPG(
	t *testing.T,
	production string,
	outcomes []probeOutcome,
	probe predicateProbe,
) {
	t.Helper()
	var missingFromOmni, extraInOmni []string

	for _, out := range outcomes {
		var predicateOK bool
		switch out.candidate.kind {
		case kindKeyword:
			predicateOK = probe.onKeyword(out.candidate)
		default:
			if probe.onNonKeyword == nil {
				continue
			}
			predicateOK = probe.onNonKeyword(out.candidate)
		}
		switch {
		case out.accepted && !predicateOK:
			missingFromOmni = append(missingFromOmni, out.candidate.display)
		case !out.accepted && predicateOK:
			extraInOmni = append(extraInOmni, out.candidate.display)
		}
	}

	if len(missingFromOmni) > 0 || len(extraInOmni) > 0 {
		t.Errorf("%s FIRST-set drift:\n"+
			"  PG accepts but predicate rejects: %v\n"+
			"  predicate accepts but PG rejects: %v",
			production, missingFromOmni, extraInOmni)
	}
}

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

func TestTypenameLeadTokensMatchPG(t *testing.T) {
	o := startFirstSetOracle(t)
	// Probe Typename via the RETURNS clause of CREATE FUNCTION. This
	// position is grammatically `func_type`, and at the FIRST-set level
	// `func_type` equals `Typename`:
	//
	//   func_type: Typename
	//            | type_function_name attrs '%' TYPE_P
	//            | SETOF type_function_name attrs '%' TYPE_P
	//
	// The %TYPE alternatives start with `type_function_name`, which is
	// already in FIRST(Typename) via SimpleTypename → GenericType. So
	// FIRST(func_type) = FIRST(Typename) ∪ {SETOF}, and {SETOF} is already
	// in FIRST(Typename) (Typename: SETOF SimpleTypename | SimpleTypename
	// ...). Therefore the RETURNS position is a sound oracle for Typename's
	// FIRST set — and it correctly accepts SETOF, unlike CAST.
	//
	// Cross-check: if a future PG version introduces a token accepted in
	// one position but not the other, TestSimpleTypenameLeadTokensMatchPG
	// (CAST-based) and this test will disagree, which is a useful alert.

	// Scope to keyword + IDENT, matching the SimpleTypename test — see
	// TestSimpleTypenameLeadTokensMatchPG rationale.
	filter := func(c candidateToken) bool {
		return c.kind == kindKeyword || c.kind == kindIdent
	}
	outcomes := runFirstSetProbe(t, o,
		"CREATE FUNCTION __omni_probe() RETURNS %s AS $$ SELECT 1 $$ LANGUAGE sql",
		renderTypeCandidate,
		filter,
	)

	assertPredicateMatchesPG(t, "Typename via RETURNS", outcomes, predicateProbe{
		onKeyword: func(c candidateToken) bool {
			p := &Parser{cur: Token{Type: c.token, Str: c.name}}
			return p.isTypenameStart()
		},
		onNonKeyword: func(c candidateToken) bool {
			p := &Parser{cur: Token{Type: c.token, Str: c.name}}
			return p.isTypenameStart()
		},
	})
}

func TestFuncTypeLeadTokensMatchPG(t *testing.T) {
	o := startFirstSetOracle(t)
	// Probe func_type via the parameter type position of CREATE FUNCTION.
	// Use a dummy_name prefix so PG parses the candidate as func_type, not
	// arg_class. (Bare `CREATE FUNCTION __omni_probe(X)` would let PG
	// interpret X as an arg_class like IN/OUT.)
	//
	// This test is intentionally a sibling of TestTypenameLeadTokensMatchPG
	// (which uses RETURNS). They MUST agree on the keyword set —
	// FIRST(func_type) == FIRST(Typename) by grammar — and any future
	// disagreement would be a useful drift alert.

	filter := func(c candidateToken) bool {
		return c.kind == kindKeyword || c.kind == kindIdent
	}
	outcomes := runFirstSetProbe(t, o,
		"CREATE FUNCTION __omni_probe(dummy_name %s) RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql",
		renderTypeCandidate,
		filter,
	)

	// Subtract probe-contamination keywords. The `dummy_name %s` template
	// cannot fully isolate func_type's FIRST set because PG's func_arg
	// rule has multiple alternatives that let the surrounding syntax
	// "absorb" tokens that aren't func_type leads at all:
	//
	//   - IN, OUT, INOUT, VARIADIC: arg_class keywords. PG re-parses
	//     `dummy_name IN int` as `param_name arg_class func_type`, so IN
	//     is consumed by arg_class, not func_type. Not a func_type lead.
	//   - DEFAULT: after `dummy_name` is parsed as a GenericType, DEFAULT
	//     introduces the default-value clause of func_arg_with_default.
	//     Not a func_type lead.
	//   - ARRAY: the renderTypeCandidate form `array` lets PG parse
	//     `dummy_name` as the GenericType and ARRAY as opt_array_bounds
	//     attached to that type, so ARRAY is a type SUFFIX, never a
	//     func_type start.
	//
	// These match PG's semantics exactly (none of them can start a
	// func_type in isolation), so the subtraction aligns the oracle with
	// the predicate's actual contract rather than masking a real bug.
	funcTypeContamination := map[int]bool{
		IN_P: true, OUT_P: true, INOUT: true, VARIADIC: true,
		DEFAULT: true, ARRAY: true,
	}
	for i := range outcomes {
		if funcTypeContamination[outcomes[i].candidate.token] {
			outcomes[i].accepted = false
		}
	}

	assertPredicateMatchesPG(t, "func_type via func arg", outcomes, predicateProbe{
		onKeyword: func(c candidateToken) bool {
			p := &Parser{cur: Token{Type: c.token, Str: c.name}}
			return p.isFuncTypeStart()
		},
		onNonKeyword: func(c candidateToken) bool {
			p := &Parser{cur: Token{Type: c.token, Str: c.name}}
			return p.isFuncTypeStart()
		},
	})
}
