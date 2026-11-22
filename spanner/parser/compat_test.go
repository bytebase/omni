package parser

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	googlesqlparser "github.com/bytebase/omni/googlesql/parser"
	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

// TestCompatibility holds spanner/parser (DialectSpanner) to
// testdata/compat.tsv: a statement for each alternative of GoogleSQL's grammar
// (../../googlesql/testdata/GoogleSQL.g4, converted from googlesql.tm by
// googlesql/parser/cmd/tm2g4) and its near-misses, with omni's verdict and the
// pinned Cloud Spanner emulator's (../testdata/pins.tsv). The test reads the
// committed file and starts no server in offline replay. OMNI_RECORD=1 judges
// rows on the pinned emulator; OMNI_RECORD=verify re-judges every row read-only.
func TestCompatibility(t *testing.T) {
	gtd := filepath.Join("..", "..", "googlesql", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		PinRole: "parser",
		Grammar: filepath.Join(gtd, "GoogleSQL.g4"),
		Tokens:  filepath.Join(gtd, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict reads the row as one statement, as the judge sends it.
func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// compatJudge asks the pinned Spanner emulator whether it accepts each
// statement, through googlesql/spanner-oracle, which runs a query, runs DML in
// a transaction it rolls back, and applies DDL to a database of its own. It
// runs only in a recording, outside any tape.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	// testcontainer.Judge starts the shared emulator under the default role;
	// when a dedicated pin role (such as parser) shares that release, verify it
	// pins the exact same image before reusing it.
	if pin.Role != "default" {
		serverPin, err := enginetest.LoadPin(testcontainer.Pins, "default")
		if err != nil {
			t.Fatalf("%v", err)
		}
		if serverPin.Image != pin.Image || serverPin.Digest != pin.Digest || serverPin.Tag != pin.Tag {
			t.Fatalf("%s pin (%s %s %s) differs from default server pin (%s %s %s)",
				pin.Role, pin.Image, pin.Digest, pin.Tag, serverPin.Image, serverPin.Digest, serverPin.Tag)
		}
	}
	o := testcontainer.Judge(t)
	ask := func(sql string) (testcontainer.Verdict, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return o.Verdict(ctx, sql)
	}
	// Startup self-checks: pin the emulator error shapes that judge rev 2 relies
	// on so an emulator upgrade that rewords any of them fails fast instead of
	// silently flipping verdicts.
	for _, probe := range []struct {
		sql  string
		want string
	}{
		{"MERGE a USING a ON NULL WHEN MATCHED THEN DELETE", compat.Accept},
		{"MERGE a USING a ON NULL WHEN MATCHED THEN DELETION", compat.Reject},
		{"SELECT EXCEPT", compat.Reject},
		{"CREATE PROPERTY GRAPH a NODE TABLES(a DYNAMIC LABEL(NULL) DYNAMIC LABEL(NULL))", compat.Reject},
		{"CREATE VIEW a AS SELECT NULL", compat.Accept},
		{"CREATE VIEW a AS NULL", compat.Reject},
	} {
		got, err := compatJudgeVerdict(probe.sql, ask)
		if err != nil || got != probe.want {
			t.Fatalf("compatJudge startup probe %q: got (%q, %v), want (%q, nil)", probe.sql, got, err, probe.want)
		}
	}
	return func(sql string) (string, error) {
		return compatJudgeVerdict(sql, ask)
	}
}

// compatJudgeVerdict evaluates sql with ask and maps the emulator's answer to
// compat.Accept or compat.Reject.
//
// For CREATE [OR REPLACE] VIEW statements without a SQL SECURITY clause, the
// Spanner emulator fails schema validation ("View `...` is missing the SQL
// SECURITY clause.") before parsing the view query body. Re-asking once with a
// SQL SECURITY INVOKER stand-in lets the emulator parse the view query body so
// broken bodies (e.g. CREATE VIEW a AS NULL) are rejected rather than masked.
func compatJudgeVerdict(sql string, ask func(string) (testcontainer.Verdict, error)) (string, error) {
	v, err := ask(sql)
	if err != nil {
		return "", err
	}
	if isMissingViewSQLSecurity(v) {
		rewritten, ok := injectViewSQLSecurity(sql)
		if !ok {
			return "", fmt.Errorf("no verdict: cannot inject SQL SECURITY into view statement %q: %s", sql, v.Message)
		}
		v2, err := ask(rewritten)
		if err != nil {
			return "", err
		}
		if isMissingViewSQLSecurity(v2) {
			return "", fmt.Errorf("no verdict: rewritten view still missing SQL SECURITY: %s", v2.Message)
		}
		return spannerVerdict(v2)
	}
	return spannerVerdict(v)
}

// spannerVerdict reads the harness's answer as a verdict. The harness rejects a
// statement on the parser's own error ("Syntax error:" for a query or DML, a
// grammar refusal, or "Error parsing Spanner DDL statement:" for DDL). Judge
// rev 2 also rejects three parser/lexer refusals that GoogleSQL emits as
// InvalidArgument without a "Syntax error:" prefix (see isSpannerParserRefusal):
//   - "Statement not supported: <Kind>Stmt": GoogleSQL's parse-error fallback
//     path (GetNextStatementKind) replaces a ParseStatement failure on an
//     unsupported statement kind with the ResolvedNodeKind name ending in
//     "Stmt", whereas a statement that parsed cleanly reports its AST node kind
//     ending in "Statement" or "Assignment" (accepted here as post-parse).
//   - "EXCEPT must be followed by ALL, DISTINCT, or \"(\"": GoogleSQL's
//     lookahead tokenizer refusal.
//   - "DYNAMIC LABEL|PROPERTIES cannot be used more than once in a single
//     element table": dynamic MakeSyntaxError in googlesql.tm's
//     dynamic_label_and_properties action.
//
// Any other harness accept, or Unimplemented (e.g. returning a STRUCT column),
// is accepted; every other error is fail-closed ("no verdict").
func spannerVerdict(v testcontainer.Verdict) (string, error) {
	switch {
	case v.Verdict == "accept":
		if v.Code == "InvalidArgument" && isSpannerParserRefusal(v.Message) {
			return compat.Reject, nil
		}
		return compat.Accept, nil
	case v.Verdict == "reject":
		return compat.Reject, nil
	case v.Code == "Unimplemented":
		return compat.Accept, nil
	}
	return "", fmt.Errorf("no verdict: %s %s (%s): %s", v.Verdict, v.Reason, v.Code, v.Message)
}

// spannerParserRefusalSubstrings enumerates GoogleSQL lexer/parser refusals
// emitted without a "Syntax error:" prefix and not captured by static literals
// in googlesql/spanner-oracle's grammar_refusals.go.
var spannerParserRefusalSubstrings = []string{
	`EXCEPT must be followed by ALL, DISTINCT, or "("`,
	`DYNAMIC LABEL cannot be used more than once in a single element table`,
	`DYNAMIC PROPERTIES cannot be used more than once in a single element table`,
}

func isSpannerParserRefusal(msg string) bool {
	if isUnsupportedStmtFallback(msg) {
		return true
	}
	for _, sub := range spannerParserRefusalSubstrings {
		if strings.HasPrefix(msg, sub) || strings.Contains(msg, ": "+sub) {
			return true
		}
	}
	return false
}

// isUnsupportedStmtFallback reports whether msg is GoogleSQL's parse-error
// fallback "Statement not supported: <ResolvedKind>Stmt ...", emitted only when
// ParseStatement failed before analyzer dispatch (as opposed to post-parse
// "Statement not supported: <ASTKind>Statement" / "...Assignment").
func isUnsupportedStmtFallback(msg string) bool {
	rest, ok := strings.CutPrefix(msg, "Statement not supported: ")
	if !ok {
		return false
	}
	kind := rest
	if i := strings.IndexAny(rest, " \t\r\n["); i >= 0 {
		kind = rest[:i]
	}
	if len(kind) <= len("Stmt") || !strings.HasSuffix(kind, "Stmt") {
		return false
	}
	for i := 0; i < len(kind); i++ {
		c := kind[i]
		isLetter := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_'
		isDigit := c >= '0' && c <= '9'
		if !isLetter && !(i > 0 && isDigit) {
			return false
		}
	}
	return true
}

func isMissingViewSQLSecurity(v testcontainer.Verdict) bool {
	return v.Verdict == "accept" && v.Kind == "ddl" && v.Code == "InvalidArgument" &&
		strings.HasPrefix(v.Message, "View `") &&
		strings.HasSuffix(v.Message, "` is missing the SQL SECURITY clause.")
}

// injectViewSQLSecurity rewrites "CREATE [OR REPLACE] VIEW <path> AS <body>" to
// insert "SQL SECURITY INVOKER" immediately before the top-level AS keyword,
// using GoogleSQL's lexer so comments, strings, and backtick-quoted identifiers
// (such as `AS`) are never misread.
func injectViewSQLSecurity(sql string) (string, bool) {
	toks, _ := googlesqlparser.Tokenize(sql)
	if len(toks) < 4 {
		return sql, false
	}
	i := 0
	if googlesqlparser.TokenName(toks[i].Type) != "CREATE" {
		return sql, false
	}
	i++
	if i+1 < len(toks) &&
		googlesqlparser.TokenName(toks[i].Type) == "OR" &&
		googlesqlparser.TokenName(toks[i+1].Type) == "REPLACE" {
		i += 2
	}
	if i >= len(toks) || googlesqlparser.TokenName(toks[i].Type) != "VIEW" {
		return sql, false
	}
	i++
	// View name: one or more identifier/keyword path segments separated by '.'.
	if i >= len(toks) || !isViewPathSegment(toks[i]) {
		return sql, false
	}
	i++
	for i < len(toks) && googlesqlparser.TokenName(toks[i].Type) == "'.'" {
		i++
		if i >= len(toks) || !isViewPathSegment(toks[i]) {
			return sql, false
		}
		i++
	}
	if i >= len(toks) || googlesqlparser.TokenName(toks[i].Type) != "AS" {
		return sql, false
	}
	pos := toks[i].Loc.Start
	if pos <= 0 || pos > len(sql) {
		return sql, false
	}
	insert := "SQL SECURITY INVOKER "
	prev := sql[pos-1]
	if prev != ' ' && prev != '\t' && prev != '\n' && prev != '\r' {
		insert = " " + insert
	}
	return sql[:pos] + insert + sql[pos:], true
}

func isViewPathSegment(tok googlesqlparser.Token) bool {
	name := googlesqlparser.TokenName(tok.Type)
	if name == "IDENTIFIER" {
		return true
	}
	// Non-reserved keywords can appear unquoted as a view path segment, but
	// delimiters, literals, EOF, INVALID, and AS itself are not path segments.
	if name == "AS" || name == "EOF" || name == "INVALID" || name == "UNKNOWN" ||
		name == "INTEGER" || name == "FLOAT" || name == "STRING" || name == "BYTES" {
		return false
	}
	// Keyword token names are uppercase ASCII words (A-Z, 0-9, _).
	for i := 0; i < len(name); i++ {
		c := name[i]
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return len(name) > 0
}

func TestSpannerVerdict(t *testing.T) {
	cases := []struct {
		name    string
		v       testcontainer.Verdict
		want    string
		wantErr bool
	}{
		// Direct accept / reject / Unimplemented
		{name: "ok", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "ok", Code: "OK"}, want: compat.Accept},
		{name: "table not found", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Table not found: a"}, want: compat.Accept},
		{name: "syntax reject", v: testcontainer.Verdict{Verdict: "reject", Kind: "query", Reason: "syntax", Code: "InvalidArgument", Message: "Syntax error: Unexpected end of input"}, want: compat.Reject},
		{name: "unimplemented struct", v: testcontainer.Verdict{Verdict: "error", Kind: "query", Reason: "infra", Code: "Unimplemented", Message: "Unsupported query shape: A struct value cannot be returned as a column value."}, want: compat.Accept},

		// Post-parse Statement not supported (<ASTKind>Statement / Assignment) -> accept
		{name: "parsed subpipeline", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: SubpipelineStatement [at 1:1] |> AGGREGATE ^"}, want: compat.Accept},
		{name: "parsed merge", v: testcontainer.Verdict{Verdict: "accept", Kind: "dml", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: MergeStatement [at 1:1] MERGE a USING a ON NULL WHEN MATCHED THEN DELETE ^"}, want: compat.Accept},
		{name: "parsed begin", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: BeginStatement [at 1:1] BEGIN ^"}, want: compat.Accept},
		{name: "parsed load data", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: AuxLoadDataStatement [at 1:1] LOAD DATA INTO a FROM FILES() ^"}, want: compat.Accept},
		{name: "parsed single assignment", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: SingleAssignment [at 1:1] SET a = NULL ^"}, want: compat.Accept},
		{name: "parsed struct assignment", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: AssignmentFromStruct [at 1:1] SET(a) = NULL ^"}, want: compat.Accept},

		// Parse-error fallback Statement not supported (<ResolvedKind>Stmt) -> reject
		{name: "parse fallback subpipeline", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: SubpipelineStmt [at 1:1] |> LOG(SELECT ^"}, want: compat.Reject},
		{name: "parse fallback merge", v: testcontainer.Verdict{Verdict: "accept", Kind: "dml", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: MergeStmt [at 1:1] MERGE a USING a ON NULL WHEN MATCHED THEN DELETION ^"}, want: compat.Reject},
		{name: "parse fallback begin", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: BeginStmt [at 1:1] BEGIN zzz ^"}, want: compat.Reject},
		{name: "parse fallback load data", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: AuxLoadDataStmt [at 1:1] LOAD DATA INTO FROM FILES() ^"}, want: compat.Reject},
		{name: "parse fallback undrop", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: UndropStmt [at 1:1] UNDROP CONNECTION a ^"}, want: compat.Reject},
		{name: "parse fallback assignment", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: AssignmentStmt [at 1:1] SET @@ a = ^"}, want: compat.Reject},
		{name: "bare Stmt not matched", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: "Statement not supported: Stmt"}, want: compat.Accept},

		// Lookahead tokenizer & dynamic grammar-action refusals -> reject
		{name: "except bare", v: testcontainer.Verdict{Verdict: "accept", Kind: "query", Reason: "semantic", Code: "InvalidArgument", Message: `EXCEPT must be followed by ALL, DISTINCT, or "(" [at 1:8] SELECT EXCEPT ^`}, want: compat.Reject},
		{name: "except wrapped in ddl", v: testcontainer.Verdict{Verdict: "accept", Kind: "ddl", Reason: "semantic", Code: "InvalidArgument", Message: `Error analyzing the definition of view ` + "`a`" + `: EXCEPT must be followed by ALL, DISTINCT, or "(" [at 1:48]`}, want: compat.Reject},
		{name: "duplicate dynamic label", v: testcontainer.Verdict{Verdict: "accept", Kind: "ddl", Reason: "semantic", Code: "InvalidArgument", Message: "DYNAMIC LABEL cannot be used more than once in a single element table [at 1:59] ..."}, want: compat.Reject},
		{name: "duplicate dynamic properties", v: testcontainer.Verdict{Verdict: "accept", Kind: "ddl", Reason: "semantic", Code: "InvalidArgument", Message: "DYNAMIC PROPERTIES cannot be used more than once in a single element table [at 1:64] ..."}, want: compat.Reject},

		// Non-InvalidArgument accept is never flipped by parser refusal patterns
		{name: "failed precondition stays accept", v: testcontainer.Verdict{Verdict: "accept", Kind: "ddl", Reason: "semantic", Code: "FailedPrecondition", Message: "Statement not supported: MergeStmt"}, want: compat.Accept},

		// Fail-closed infra / unknown verdicts
		{name: "unavailable", v: testcontainer.Verdict{Verdict: "error", Kind: "query", Reason: "infra", Code: "Unavailable", Message: "transport is closing"}, wantErr: true},
		{name: "deadline exceeded", v: testcontainer.Verdict{Verdict: "error", Kind: "ddl", Reason: "infra", Code: "DeadlineExceeded", Message: "context deadline exceeded"}, wantErr: true},
		{name: "generic internal", v: testcontainer.Verdict{Verdict: "error", Kind: "query", Reason: "infra", Code: "Internal", Message: "unexpected error"}, wantErr: true},
		{name: "unknown verdict", v: testcontainer.Verdict{Verdict: "maybe", Kind: "query", Code: "OK"}, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := spannerVerdict(tc.v)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("spannerVerdict(%+v) = %q, nil; want error", tc.v, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("spannerVerdict(%+v) unexpected error: %v", tc.v, err)
			}
			if got != tc.want {
				t.Fatalf("spannerVerdict(%+v) = %q, want %q", tc.v, got, tc.want)
			}
		})
	}
}

func TestCompatJudgeViewSQLSecurityRewrite(t *testing.T) {
	ask := func(sql string) (testcontainer.Verdict, error) {
		switch sql {
		case "CREATE VIEW a AS SELECT NULL",
			"CREATE VIEW a AS NULL",
			"CREATE VIEW a AS SELECT",
			"CREATE OR REPLACE VIEW s.a AS SELECT NULL",
			"CREATE VIEW `AS` AS SELECT NULL",
			"CREATE VIEW a /* AS */ AS NULL":
			return testcontainer.Verdict{
				Verdict: "accept",
				Kind:    "ddl",
				Reason:  "semantic",
				Code:    "InvalidArgument",
				Message: "View `a` is missing the SQL SECURITY clause.",
			}, nil
		case "CREATE VIEW a SQL SECURITY INVOKER AS SELECT NULL",
			"CREATE OR REPLACE VIEW s.a SQL SECURITY INVOKER AS SELECT NULL",
			"CREATE VIEW `AS` SQL SECURITY INVOKER AS SELECT NULL":
			return testcontainer.Verdict{
				Verdict: "accept",
				Kind:    "ddl",
				Reason:  "semantic",
				Code:    "InvalidArgument",
				Message: "Error analyzing the definition of view `a`: CREATE VIEW columns must be named, but column 1 has no name",
			}, nil
		case "CREATE VIEW a SQL SECURITY INVOKER AS NULL",
			"CREATE VIEW a SQL SECURITY INVOKER AS SELECT",
			"CREATE VIEW a /* AS */ SQL SECURITY INVOKER AS NULL":
			return testcontainer.Verdict{
				Verdict: "reject",
				Kind:    "ddl",
				Reason:  "syntax",
				Code:    "InvalidArgument",
				Message: "Error analyzing the definition of view `a`: Syntax error: Unexpected keyword NULL",
			}, nil
		case "CREATE VIEW loop AS SELECT 1",
			"CREATE VIEW loop SQL SECURITY INVOKER AS SELECT 1",
			"CREATE VIEW unrewritable":
			return testcontainer.Verdict{
				Verdict: "accept",
				Kind:    "ddl",
				Reason:  "semantic",
				Code:    "InvalidArgument",
				Message: "View `loop` is missing the SQL SECURITY clause.",
			}, nil
		default:
			t.Fatalf("unexpected emulator query: %q", sql)
			return testcontainer.Verdict{}, nil
		}
	}

	for _, tc := range []struct {
		sql     string
		want    string
		wantErr bool
	}{
		{"CREATE VIEW a AS SELECT NULL", compat.Accept, false},
		{"CREATE VIEW a AS NULL", compat.Reject, false},
		{"CREATE VIEW a AS SELECT", compat.Reject, false},
		{"CREATE OR REPLACE VIEW s.a AS SELECT NULL", compat.Accept, false},
		{"CREATE VIEW `AS` AS SELECT NULL", compat.Accept, false},
		{"CREATE VIEW a /* AS */ AS NULL", compat.Reject, false},
		{"CREATE VIEW loop AS SELECT 1", "", true},
		{"CREATE VIEW unrewritable", "", true},
	} {
		got, err := compatJudgeVerdict(tc.sql, ask)
		if tc.wantErr {
			if err == nil {
				t.Errorf("compatJudgeVerdict(%q) = %q, nil; want error", tc.sql, got)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Errorf("compatJudgeVerdict(%q) = (%q, %v), want (%q, nil)", tc.sql, got, err, tc.want)
		}
	}

	for _, tc := range []struct {
		sql  string
		want string
		ok   bool
	}{
		{"CREATE VIEW a AS SELECT NULL", "CREATE VIEW a SQL SECURITY INVOKER AS SELECT NULL", true},
		{"CREATE OR REPLACE VIEW s.a AS SELECT 1", "CREATE OR REPLACE VIEW s.a SQL SECURITY INVOKER AS SELECT 1", true},
		{"CREATE VIEW `a`AS SELECT 1", "CREATE VIEW `a` SQL SECURITY INVOKER AS SELECT 1", true},
		{"CREATE VIEW `AS` AS SELECT 1", "CREATE VIEW `AS` SQL SECURITY INVOKER AS SELECT 1", true},
		{"CREATE VIEW a /* AS */ AS NULL", "CREATE VIEW a /* AS */ SQL SECURITY INVOKER AS NULL", true},
		{"CREATE VIEW a -- AS\nAS NULL", "CREATE VIEW a -- AS\nSQL SECURITY INVOKER AS NULL", true},
		{"CREATE VIEW a SQL SECURITY INVOKER AS SELECT 1", "CREATE VIEW a SQL SECURITY INVOKER AS SELECT 1", false},
		{"CREATE VIEW a (c) AS SELECT 1", "CREATE VIEW a (c) AS SELECT 1", false},
		{"CREATE VIEW a", "CREATE VIEW a", false},
		{"SELECT 1", "SELECT 1", false},
	} {
		got, ok := injectViewSQLSecurity(tc.sql)
		if ok != tc.ok || got != tc.want {
			t.Errorf("injectViewSQLSecurity(%q) = (%q, %v), want (%q, %v)", tc.sql, got, ok, tc.want, tc.ok)
		}
	}
}
