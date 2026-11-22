package parser

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ScriptDOM's T-SQL grammar (../testdata/TSql160.g4, converted
// from TSql160.g by ./cmd/scriptdomg4) and its near-misses, and the token
// forms of suite:lexer, with omni's verdict and the pinned SQL Server's. The
// test reads the file and starts no server. OMNI_RECORD=1 judges the rows
// that have no verdict yet on the pinned SQL Server.
//
// testdata/README.md says what SQL Server's judge needs.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		PinRole: "parser",
		Grammar: filepath.Join(td, "TSql160.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict reads the row as one batch, as the judge sends it.
func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// SQL Server error numbers handled by judge rev 2 beyond hard rejections.
const (
	mssqlErrGotoLabel      = 133   // A GOTO statement references the label 'a' but the label has not been declared.
	mssqlErrUndeclaredVar  = 137   // Must declare the scalar variable "@a".
	mssqlErrUnusedCTE      = 422   // Common table expression defined but not used.
	mssqlErrUndeclaredTbl  = 1087  // Must declare the table variable "@a".
	mssqlErrFTSNotInstall  = 7609  // Full-Text Search is not installed, or a full-text component cannot be loaded.
	mssqlErrFTSBadLang     = 7678  // The following string is not defined as a language alias in syslanguages: a.
	mssqlErrParseType      = 10761 // Invalid data type a in function PARSE / TRY_PARSE.
	mssqlErrNotCursor      = 16948 // The variable '@a' is not a cursor variable, ...
	mssqlErrExtLangEnvVars = 39133 // CREATE/ALTER EXTERNAL LANGUAGE statement failed. The environment variables string is invalid.
)

// mssqlLookups are errors SQL Server's parser raises for a name it looks up in
// what the server holds: an object (208, 15151), a collation (448), or a
// database (911). When every error reported in a batch is a lookup, the
// statement counts as accepted.
var mssqlLookups = map[int32]bool{
	208:   true,
	448:   true,
	911:   true,
	15151: true,
}

const maxJudgeDeclared = 16

var (
	reVarName   = regexp.MustCompile(`["'](@[^"']+)["']`)
	reGotoLabel = regexp.MustCompile(`label '([^']+)'`)
	reParseType = regexp.MustCompile(`Invalid data type ([^\s]+) in function (TRY_PARSE|PARSE)\.`)

	reFTSStoplistLang = regexp.MustCompile(`(?i)^(\s*ALTER\s+FULLTEXT\s+STOPLIST\s+(?:\[[^\]]+\]|"[^"]+"|[A-Za-z_@#][A-Za-z0-9_@#$]*)\s+)(?:(?:ADD|DROP)\s+(?:N?'(?:''|[^'])*'|"(?:""|[^"])*")|DROP\s+ALL)\s+LANGUAGE\s+(?:N?'(?:''|[^'])*'|0[xX][0-9a-fA-F]*|[0-9]+|\[[^\]]+\]|"(?:""|[^"])*"|[A-Za-z_@#][A-Za-z0-9_@#$]*)`)
	reFTSIndexHead    = regexp.MustCompile(`(?i)^\s*(?:CREATE|ALTER)\s+FULLTEXT\s+INDEX\b`)
	reFTSIndexLang    = regexp.MustCompile(`(?i)\s+LANGUAGE\s+(?:N?'(?:''|[^'])*'|0[xX][0-9a-fA-F]*|[0-9]+|\[[^\]]+\]|"(?:""|[^"])*"|[A-Za-z_@#][A-Za-z0-9_@#$]*)`)

	reCTEHead        = regexp.MustCompile(`(?i)^\s*WITH\s+(?:CHANGE_TRACKING_CONTEXT\s*\([^)]*\)\s*,\s*)?(\[[^\]]+\]|"[^"]+"|[A-Za-z_@#][A-Za-z0-9_@#$]*)\s*(?:\([^)]*\))?\s*AS\s*\(`)
	reCTESelectStar  = regexp.MustCompile(`(?i)(\)\s*SELECT\s+\*)(\s*;?\s*)$`)
	reExtLangEnvVars = regexp.MustCompile(`(?i)(\bENVIRONMENT_VARIABLES\s*=\s*)(?:N?'(?:''|[^'])*')`)
)

// compatJudge asks the pinned SQL Server whether its parser accepts each
// statement, sent as a batch of its own under SET PARSEONLY ON, which runs
// nothing, and read by parseOnlyVerdictRev2. It runs only in a recording or
// verification, outside any tape.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	defPin, err := enginetest.LoadPin(testcontainer.Pins, testcontainer.Default)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if pin.Image != defPin.Image || pin.Digest != defPin.Digest || pin.Tag != defPin.Tag {
		t.Fatalf("parser pin %s differs from default server pin %s", pin, defPin)
	}
	conn := testcontainer.Judge(t)
	ctx := context.Background()
	// Create scratch databases FEDERATION and ROOT so USE FEDERATION ... and
	// USE ROOT ... resolve the database name and SQL Server parses past it
	// instead of aborting with Msg 911 before reading the rest of the batch.
	for _, db := range []string{"FEDERATION", "ROOT"} {
		stmt := fmt.Sprintf("IF DB_ID('%s') IS NULL CREATE DATABASE [%s]", db, db)
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	if _, err := conn.ExecContext(ctx, "USE ["+testcontainer.JudgeDatabase+"]"); err != nil {
		t.Fatalf("USE [%s]: %v", testcontainer.JudgeDatabase, err)
	}
	return func(stmt string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return parseOnlyVerdictRev2(stmt, testcontainer.ServerErrors, func(declarations, batch string) error {
			return testcontainer.ParseOnly(ctx, conn, testcontainer.JudgeDatabase, declarations, batch)
		})
	}
}

// parseOnlyVerdictRev2 reads SQL Server's answer to stmt sent under SET
// PARSEONLY ON. Unlike judge rev 1 (compat.ParseOnlyVerdict), it inspects all
// errors returned in a batch (so Msg 448 + Msg 102 rejects) and retries
// context-only errors (undeclared variables, undeclared GOTO labels, unused
// CTEs, PARSE/TRY_PARSE placeholder types, EXTERNAL LANGUAGE
// ENVIRONMENT_VARIABLES strings, and FULLTEXT LANGUAGE clauses whose lookup
// aborts before the statement tail) so SQL Server parses the full statement
// honestly.
func parseOnlyVerdictRev2(stmt string, read func(error) ([]compat.SQLServerError, bool), ask func(declarations, batch string) error) (string, error) {
	var varNames []string
	varKinds := map[string]string{}
	var labels []string
	seenLabels := map[string]bool{}
	curStmt := stmt
	visited := map[string]bool{}

	buildBatch := func() (string, string) {
		var decl strings.Builder
		for _, n := range varNames {
			fmt.Fprintf(&decl, "DECLARE %s %s;\n", n, varKinds[n])
		}
		var b strings.Builder
		b.WriteString(curStmt)
		for _, lbl := range labels {
			b.WriteByte('\n')
			b.WriteString(lbl)
			b.WriteByte(':')
		}
		return decl.String(), b.String()
	}

	for {
		decls, batch := buildBatch()
		key := decls + "\x00" + batch
		if visited[key] {
			return compat.Reject, nil
		}
		visited[key] = true

		err := ask(decls, batch)
		if err == nil {
			return compat.Accept, nil
		}
		errs, ok := read(err)
		if !ok {
			return "", err
		}
		if len(errs) == 0 {
			return "", fmt.Errorf("SQL Server reported an error with no error in it: %w", err)
		}

		// Try to resolve the first retryable context error in the batch.
		retried := false
		for _, e := range errs {
			switch e.Number {
			case mssqlErrUndeclaredVar, mssqlErrUndeclaredTbl, mssqlErrNotCursor:
				kind := "INT"
				switch e.Number {
				case mssqlErrUndeclaredTbl:
					kind = "TABLE (a INT)"
				case mssqlErrNotCursor:
					kind = "CURSOR"
				}
				m := reVarName.FindStringSubmatch(e.Message)
				if m == nil {
					return "", fmt.Errorf("SQL Server's error %d names no variable: %q", e.Number, e.Message)
				}
				name := m[1]
				if prev, exists := varKinds[name]; exists && (prev == kind || e.Number != mssqlErrNotCursor) {
					return compat.Reject, nil
				}
				if len(varNames) == maxJudgeDeclared {
					return "", fmt.Errorf("SQL Server still names an undeclared variable after %d declarations: %w", maxJudgeDeclared, err)
				}
				if _, exists := varKinds[name]; !exists {
					varNames = append(varNames, name)
				}
				varKinds[name] = kind
				retried = true

			case mssqlErrGotoLabel:
				m := reGotoLabel.FindStringSubmatch(e.Message)
				if m == nil {
					return "", fmt.Errorf("SQL Server's error %d names no label: %q", e.Number, e.Message)
				}
				lbl := m[1]
				if seenLabels[lbl] {
					return compat.Reject, nil
				}
				if len(labels) == maxJudgeDeclared {
					return "", fmt.Errorf("SQL Server still names an undeclared label after %d labels: %w", maxJudgeDeclared, err)
				}
				seenLabels[lbl] = true
				labels = append(labels, lbl)
				retried = true

			case mssqlErrFTSNotInstall, mssqlErrFTSBadLang:
				if next, ok := rewriteFulltextLanguage(curStmt); ok && next != curStmt {
					curStmt = next
					retried = true
				}

			case mssqlErrUnusedCTE:
				if next, ok := rewriteUnusedCTE(curStmt); ok && next != curStmt {
					curStmt = next
					retried = true
				}

			case mssqlErrParseType:
				if next, ok := rewriteParseType(curStmt, e.Message); ok && next != curStmt {
					curStmt = next
					retried = true
				}

			case mssqlErrExtLangEnvVars:
				if next, ok := rewriteExtLangEnvVars(curStmt); ok && next != curStmt {
					curStmt = next
					retried = true
				}
			}
			if retried {
				break
			}
		}
		if retried {
			continue
		}

		// No retry applies: accept only if EVERY error in the batch is a pure
		// server lookup (never masking a syntax error such as Msg 448 + 102).
		allLookups := true
		for _, e := range errs {
			if !mssqlLookups[e.Number] {
				allLookups = false
				break
			}
		}
		if allLookups {
			return compat.Accept, nil
		}
		return compat.Reject, nil
	}
}

// rewriteFulltextLanguage removes or normalizes a validated FULLTEXT LANGUAGE
// clause whose lookup aborted before SQL Server read the rest of the statement.
//   - ALTER FULLTEXT STOPLIST ... {ADD/DROP <str>|DROP ALL} LANGUAGE <term><tail>
//     becomes ALTER FULLTEXT STOPLIST ... DROP ALL<tail>, which keeps the
//     exact continuation grammar (required semicolon and EOF).
//   - CREATE/ALTER FULLTEXT INDEX ... LANGUAGE <term><tail> drops the optional
//     LANGUAGE <term> clause so SQL Server parses <tail>.
func rewriteFulltextLanguage(stmt string) (string, bool) {
	if loc := reFTSStoplistLang.FindStringSubmatchIndex(stmt); loc != nil {
		return stmt[loc[2]:loc[3]] + "DROP ALL" + stmt[loc[1]:], true
	}
	if reFTSIndexHead.MatchString(stmt) {
		if loc := reFTSIndexLang.FindStringIndex(stmt); loc != nil {
			prefix, suffix := stmt[:loc[0]], stmt[loc[1]:]
			sep := ""
			if len(prefix) > 0 && len(suffix) > 0 && isIdentChar(prefix[len(prefix)-1]) && isIdentChar(suffix[0]) {
				sep = " "
			}
			return prefix + sep + suffix, true
		}
	}
	return stmt, false
}

func isIdentChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_' || b == '@' || b == '#' || b == '$'
}

// rewriteUnusedCTE appends `FROM <cte>` to a trailing table-less `SELECT *`
// so SQL Server evaluates the full CTE statement instead of stopping at Msg 422
// (which can also precede real rejections such as Msg 22108 for
// WITH CHANGE_TRACKING_CONTEXT + SELECT).
func rewriteUnusedCTE(stmt string) (string, bool) {
	head := reCTEHead.FindStringSubmatch(stmt)
	if head == nil {
		return stmt, false
	}
	cteName := head[1]
	if loc := reCTESelectStar.FindStringSubmatchIndex(stmt); loc != nil {
		return stmt[:loc[3]] + " FROM " + cteName + stmt[loc[4]:loc[5]], true
	}
	return stmt, false
}

// rewriteParseType replaces the rejected placeholder type in PARSE(... AS <T>
// or TRY_PARSE(... AS <T> with INT so SQL Server parses the rest of the call
// and statement.
func rewriteParseType(stmt, msg string) (string, bool) {
	m := reParseType.FindStringSubmatch(msg)
	if m == nil {
		return stmt, false
	}
	typeName, fnName := m[1], m[2]
	re := regexp.MustCompile(`(?i)(\b` + regexp.QuoteMeta(fnName) + `\s*\([^;]*?\bAS\s+)` + regexp.QuoteMeta(typeName) + `($|[^A-Za-z0-9_@#$])`)
	loc := re.FindStringSubmatchIndex(stmt)
	if loc == nil {
		return stmt, false
	}
	return stmt[:loc[3]] + "INT" + stmt[loc[4]:], true
}

// rewriteExtLangEnvVars replaces the placeholder string literal in
// ENVIRONMENT_VARIABLES = '<str>' with a valid JSON object literal so SQL
// Server parses past the value check and validates the rest of the statement.
func rewriteExtLangEnvVars(stmt string) (string, bool) {
	loc := reExtLangEnvVars.FindStringSubmatchIndex(stmt)
	if loc == nil {
		return stmt, false
	}
	return stmt[:loc[3]] + `'{"a":"b"}'` + stmt[loc[1]:], true
}

// TestParseOnlyVerdictRev2 holds judge rev 2's verdict logic to both valid
// context-only statements and broken near-misses without requiring a live
// server.
func TestParseOnlyVerdictRev2(t *testing.T) {
	makeErr := func(errs ...compat.SQLServerError) error {
		return fmt.Errorf("mssql: %w", errStub{errs})
	}
	readErr := func(err error) ([]compat.SQLServerError, bool) {
		var s errStub
		if !errors.As(err, &s) {
			return nil, false
		}
		return s.errs, true
	}

	t.Run("multi_error_lookup_and_syntax_rejects", func(t *testing.T) {
		v, err := parseOnlyVerdictRev2("ALTER DATABASE a COLLATE a zzz", readErr, func(_, _ string) error {
			return makeErr(
				compat.SQLServerError{Number: 448, Message: "Invalid collation 'a'."},
				compat.SQLServerError{Number: 102, Message: "Incorrect syntax near 'zzz'."},
			)
		})
		if err != nil || v != compat.Reject {
			t.Fatalf("got (%q, %v), want (%q, nil)", v, err, compat.Reject)
		}
	})

	t.Run("pure_lookup_accepts", func(t *testing.T) {
		v, err := parseOnlyVerdictRev2("ALTER DATABASE a COLLATE a", readErr, func(_, _ string) error {
			return makeErr(compat.SQLServerError{Number: 448, Message: "Invalid collation 'a'."})
		})
		if err != nil || v != compat.Accept {
			t.Fatalf("got (%q, %v), want (%q, nil)", v, err, compat.Accept)
		}
	})

	t.Run("goto_label_retry_accepts_valid_and_rejects_broken", func(t *testing.T) {
		v, err := parseOnlyVerdictRev2("GOTO a", readErr, func(_, batch string) error {
			if batch == "GOTO a" {
				return makeErr(compat.SQLServerError{Number: 133, Message: "A GOTO statement references the label 'a' but the label has not been declared."})
			}
			if batch == "GOTO a\na:" {
				return nil
			}
			return makeErr(compat.SQLServerError{Number: 102, Message: "unexpected"})
		})
		if err != nil || v != compat.Accept {
			t.Fatalf("valid GOTO: got (%q, %v), want (%q, nil)", v, err, compat.Accept)
		}

		v, err = parseOnlyVerdictRev2("GOTO a zzz", readErr, func(_, _ string) error {
			return makeErr(compat.SQLServerError{Number: 102, Message: "Incorrect syntax near 'zzz'."})
		})
		if err != nil || v != compat.Reject {
			t.Fatalf("broken GOTO: got (%q, %v), want (%q, nil)", v, err, compat.Reject)
		}
	})

	t.Run("unused_cte_retry_distinguishes_valid_from_change_tracking", func(t *testing.T) {
		v, err := parseOnlyVerdictRev2("WITH a AS(SELECT *) SELECT *", readErr, func(_, batch string) error {
			switch batch {
			case "WITH a AS(SELECT *) SELECT *":
				return makeErr(compat.SQLServerError{Number: 422, Message: "Common table expression defined but not used."})
			case "WITH a AS(SELECT *) SELECT * FROM a":
				return nil
			default:
				return makeErr(compat.SQLServerError{Number: 102, Message: "unexpected"})
			}
		})
		if err != nil || v != compat.Accept {
			t.Fatalf("valid CTE: got (%q, %v), want (%q, nil)", v, err, compat.Accept)
		}

		v, err = parseOnlyVerdictRev2("WITH CHANGE_TRACKING_CONTEXT(0x00), a AS(SELECT *) SELECT *", readErr, func(_, batch string) error {
			switch batch {
			case "WITH CHANGE_TRACKING_CONTEXT(0x00), a AS(SELECT *) SELECT *":
				return makeErr(compat.SQLServerError{Number: 422, Message: "Common table expression defined but not used."})
			case "WITH CHANGE_TRACKING_CONTEXT(0x00), a AS(SELECT *) SELECT * FROM a":
				return makeErr(compat.SQLServerError{Number: 22108, Message: "The CHANGE_TRACKING_CONTEXT WITH clause cannot be used with a SELECT statement."})
			default:
				return makeErr(compat.SQLServerError{Number: 102, Message: "unexpected"})
			}
		})
		if err != nil || v != compat.Reject {
			t.Fatalf("CHANGE_TRACKING CTE: got (%q, %v), want (%q, nil)", v, err, compat.Reject)
		}
	})

	t.Run("parse_type_retry_validates_tail", func(t *testing.T) {
		v, err := parseOnlyVerdictRev2("SELECT PARSE(1 AS a) zzz", readErr, func(_, batch string) error {
			switch batch {
			case "SELECT PARSE(1 AS a) zzz":
				return makeErr(compat.SQLServerError{Number: 10761, Message: "Invalid data type a in function PARSE."})
			case "SELECT PARSE(1 AS INT) zzz":
				return nil
			default:
				return makeErr(compat.SQLServerError{Number: 102, Message: "unexpected"})
			}
		})
		if err != nil || v != compat.Accept {
			t.Fatalf("valid PARSE: got (%q, %v), want (%q, nil)", v, err, compat.Accept)
		}

		v, err = parseOnlyVerdictRev2("SELECT PARSE(1 AS a) FROM", readErr, func(_, batch string) error {
			switch batch {
			case "SELECT PARSE(1 AS a) FROM":
				return makeErr(compat.SQLServerError{Number: 10761, Message: "Invalid data type a in function PARSE."})
			case "SELECT PARSE(1 AS INT) FROM":
				return makeErr(compat.SQLServerError{Number: 102, Message: "Incorrect syntax near 'FROM'."})
			default:
				return nil
			}
		})
		if err != nil || v != compat.Reject {
			t.Fatalf("broken PARSE tail: got (%q, %v), want (%q, nil)", v, err, compat.Reject)
		}
	})

	t.Run("fulltext_language_tail_probe", func(t *testing.T) {
		cases := []struct {
			in, rewritten string
			retryErr      error
			want          string
		}{
			{
				in:        "ALTER FULLTEXT INDEX ON a ADD(a LANGUAGE 'a')",
				rewritten: "ALTER FULLTEXT INDEX ON a ADD(a)",
				want:      compat.Accept,
			},
			{
				in:        "ALTER FULLTEXT INDEX ON a ADD(a LANGUAGE 'a' zzz)",
				rewritten: "ALTER FULLTEXT INDEX ON a ADD(a zzz)",
				retryErr:  makeErr(compat.SQLServerError{Number: 102, Message: "Incorrect syntax near 'zzz'."}),
				want:      compat.Reject,
			},
			{
				in:        "ALTER FULLTEXT STOPLIST a ADD 'a' LANGUAGE 'a';",
				rewritten: "ALTER FULLTEXT STOPLIST a DROP ALL;",
				want:      compat.Accept,
			},
			{
				in:        "ALTER FULLTEXT STOPLIST a ADD 'a' LANGUAGE 'a'",
				rewritten: "ALTER FULLTEXT STOPLIST a DROP ALL",
				retryErr:  makeErr(compat.SQLServerError{Number: 10736, Message: "A full-text stoplist statement must be terminated by a semi-colon (;)."}),
				want:      compat.Reject,
			},
			{
				in:        "ALTER FULLTEXT STOPLIST a DROP ALL LANGUAGE 'a' zzz",
				rewritten: "ALTER FULLTEXT STOPLIST a DROP ALL zzz",
				retryErr:  makeErr(compat.SQLServerError{Number: 102, Message: "Incorrect syntax near 'zzz'."}),
				want:      compat.Reject,
			},
		}
		for _, tc := range cases {
			v, err := parseOnlyVerdictRev2(tc.in, readErr, func(_, batch string) error {
				if batch == tc.in {
					return makeErr(compat.SQLServerError{Number: 7678, Message: "not defined"})
				}
				if batch != tc.rewritten {
					t.Errorf("rewrite(%q) = %q, want %q", tc.in, batch, tc.rewritten)
				}
				return tc.retryErr
			})
			if err != nil || v != tc.want {
				t.Errorf("%q: got (%q, %v), want (%q, nil)", tc.in, v, err, tc.want)
			}
		}
	})

	t.Run("ext_lang_env_vars_retry", func(t *testing.T) {
		validIn := "ALTER EXTERNAL LANGUAGE a SET(CONTENT = 0x00, FILE_NAME = 'a', ENVIRONMENT_VARIABLES = 'a')"
		validOut := `ALTER EXTERNAL LANGUAGE a SET(CONTENT = 0x00, FILE_NAME = 'a', ENVIRONMENT_VARIABLES = '{"a":"b"}')`
		v, err := parseOnlyVerdictRev2(validIn, readErr, func(_, batch string) error {
			if batch == validIn {
				return makeErr(compat.SQLServerError{Number: 39133, Message: "invalid env vars"})
			}
			if batch == validOut {
				return nil
			}
			return makeErr(compat.SQLServerError{Number: 102, Message: "unexpected"})
		})
		if err != nil || v != compat.Accept {
			t.Fatalf("valid env vars: got (%q, %v), want (%q, nil)", v, err, compat.Accept)
		}

		badIn := "ALTER EXTERNAL LANGUAGE a SET(CONTENT = 0x00, FILE_NAME = 'a', ENVIRONMENT_VARIABLES = 'a' zzz)"
		badOut := `ALTER EXTERNAL LANGUAGE a SET(CONTENT = 0x00, FILE_NAME = 'a', ENVIRONMENT_VARIABLES = '{"a":"b"}' zzz)`
		v, err = parseOnlyVerdictRev2(badIn, readErr, func(_, batch string) error {
			if batch == badIn {
				return makeErr(compat.SQLServerError{Number: 39133, Message: "invalid env vars"})
			}
			if batch == badOut {
				return makeErr(compat.SQLServerError{Number: 102, Message: "Incorrect syntax near 'zzz'."})
			}
			return nil
		})
		if err != nil || v != compat.Reject {
			t.Fatalf("broken env vars: got (%q, %v), want (%q, nil)", v, err, compat.Reject)
		}
	})
}

type errStub struct {
	errs []compat.SQLServerError
}

func (errStub) Error() string { return "stub" }
