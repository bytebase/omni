package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	_ "github.com/sijms/go-ora/v2"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
)

func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/OracleDocs.g4, the statement pages of Oracle's
// SQL Language Reference that bnf/ quotes, converted by cmd/docsgrammar,
// and its near-misses, with omni's verdict and Oracle's. A rejection
// carries the byte Oracle reports it at. The test reads the file and starts
// no server. OMNI_RECORD=1 judges the rows that have no verdict yet on the
// pinned Oracle, by compat.OracleVerdict.
//
// testdata/README.md says what Oracle's judge needs.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "OracleDocs.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict reads the row as DBMS_SQL.PARSE does: one statement. Parse's
// error is a rejection at its byte; a text that holds no statement, or a
// second one, is rejected where the first ends.
func omniVerdict(sql string) string {
	list, err := Parse(sql)
	if err != nil {
		var pe *ParseError
		if errors.As(err, &pe) && pe.Position >= 0 && pe.Position <= len(sql) {
			return compat.RejectAt(pe.Position)
		}
		return compat.Reject
	}
	if list == nil || len(list.Items) == 0 {
		return compat.RejectAt(len(sql))
	}
	if len(list.Items) > 1 {
		return compat.Reject
	}
	return compat.Accept
}

// compatJudge asks the pinned Oracle whether it parses each statement,
// with compat.OracleParse in the judge's scratch schema, and reads the
// answer with compat.OracleVerdict. It runs only in a recording. With
// ORACLE_PARSER_REF_DSN set it asks that server instead of starting one:
// it must run the pinned release, and its user must be a scratch account,
// since every statement that runs is followed by dropping all it owns.
//
// A text Oracle gives no answer to in judgeTimeout gets compat.ErrNoAnswer.
// Some it never answers: its parse runs on in the server, and its session
// with it, holding the locks the next statements wait on. On the pinned
// server the judge ends that session's process and asks on a new session;
// on ORACLE_PARSER_REF_DSN's it cannot, and the recording fails.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	open, canEnd := judgeSessions(t, pin)
	db, end := open()
	return func(stmt string) (string, error) {
		return compat.OracleVerdict(stmt, func(text string) error {
			done := make(chan error, 1)
			session := db
			go func() {
				_, err := session.ExecContext(context.Background(), compat.OracleParse, text, testcontainer.JudgeFixture)
				done <- err
			}()
			select {
			case err := <-done:
				return err
			case <-time.After(judgeTimeout):
				if !canEnd {
					t.Fatalf("Oracle gave no answer in %s to %q, and ORACLE_PARSER_REF_DSN's session cannot be ended from here", judgeTimeout, text)
				}
				end()
				select {
				case <-done:
				case <-time.After(10 * time.Second):
				}
				db, end = open()
				return fmt.Errorf("after %s: %w", judgeTimeout, compat.ErrNoAnswer)
			}
		})
	}
}

// judgeTimeout bounds one text of the judge's. A text takes a few
// milliseconds, and DDL with the schema's rebuild after it a second on a
// loaded host.
const judgeTimeout = 15 * time.Second

// judgeSessions returns what opens a session for the judge, the pinned
// server's (testcontainer.Judge) or ORACLE_PARSER_REF_DSN's, and whether
// it can end a session the server does not answer on.
func judgeSessions(t testing.TB, pin enginetest.Pin) (open func() (*sql.DB, func()), canEnd bool) {
	t.Helper()
	dsn := os.Getenv("ORACLE_PARSER_REF_DSN")
	if dsn == "" {
		return testcontainer.Judge(t, pin), true
	}
	open = func() (*sql.DB, func()) { return testcontainer.Cuttable(dsn) }
	db, _ := open()
	defer db.Close()
	ctx := context.Background()
	var version string
	if err := db.QueryRowContext(ctx, `SELECT version_full FROM product_component_version WHERE product LIKE 'Oracle%'`).Scan(&version); err != nil {
		t.Fatalf("ORACLE_PARSER_REF_DSN: read the version: %v", err)
	}
	if version != pin.Tag {
		t.Fatalf("ORACLE_PARSER_REF_DSN runs Oracle %s, and compat.tsv's answers come from %s", version, pin.Tag)
	}
	// The scratch schema starts as OracleParse leaves it: the fixture alone.
	if err := db.QueryRowContext(ctx, "SELECT 1 FROM user_tables WHERE table_name = 'T'").Scan(new(int)); errors.Is(err, sql.ErrNoRows) {
		if _, err := db.ExecContext(ctx, testcontainer.JudgeFixture); err != nil {
			t.Fatalf("ORACLE_PARSER_REF_DSN: %s: %v", testcontainer.JudgeFixture, err)
		}
	} else if err != nil {
		t.Fatalf("ORACLE_PARSER_REF_DSN: %v", err)
	}
	return open, false
}

// lexerCorpus is the suite:lexer rows: the token forms the Language
// Reference gives that tokens.tsv writes one way, each in a statement that
// is otherwise plain. The pages are cited in ../testdata/tokens.tsv. A
// compat.tsv row is one line, so no form here holds a line break.
var lexerCorpus = []string{
	// Names: nonquoted, with $ and #, quoted, and what a name may not be.
	`SELECT a$b, a#c, a_1 FROM dual`,
	`SELECT "a" FROM dual`,
	`SELECT 1 "a b" FROM dual`,
	`SELECT 1 "select" FROM dual`,
	`SELECT 1 "a""b" FROM dual`,
	`SELECT 1 "" FROM dual`,
	`SELECT 1 "a FROM dual`,
	`SELECT 1 _a FROM dual`,
	`SELECT 1 $a FROM dual`,
	`SELECT 1 a.b FROM dual`,
	`SELECT 1 FROM "DUAL"`,
	`SELECT 1 FROM dual "d"`,
	`SELECT dummy FROM sys.dual`,
	`SELECT dummy FROM "SYS"."DUAL"`,
	`SELECT 1 FROM dual@a`,
	// Text literals: '' inside, empty, N, q quoting with each delimiter.
	`SELECT 'it''s' FROM dual`,
	`SELECT '' FROM dual`,
	`SELECT 'a FROM dual`,
	`SELECT N'a' FROM dual`,
	`SELECT n'a' FROM dual`,
	`SELECT q'[a'b]' FROM dual`,
	`SELECT Q'{a}' FROM dual`,
	`SELECT q'(a)' FROM dual`,
	`SELECT q'<a>' FROM dual`,
	`SELECT q'!a!' FROM dual`,
	`SELECT q'#a#' FROM dual`,
	`SELECT nq'[a]' FROM dual`,
	`SELECT NQ'[a]' FROM dual`,
	`SELECT q'[a' FROM dual`,
	`SELECT q'[a)' FROM dual`,
	`SELECT q' a ' FROM dual`,
	`SELECT 'a' 'b' FROM dual`,
	// Numeric literals.
	`SELECT 1, 1.5, .5, 1., 1e3, 1E-3, 1.5e+3 FROM dual`,
	`SELECT 1f, 1.5F, 1d, 2.5D FROM dual`,
	`SELECT -1, +1 FROM dual`,
	`SELECT 0x1F FROM dual`,
	`SELECT 1e FROM dual`,
	`SELECT 1..2 FROM dual`,
	`SELECT BINARY_FLOAT_INFINITY, BINARY_DOUBLE_NAN FROM dual`,
	// Datetime and interval literals.
	`SELECT DATE '2026-01-01' FROM dual`,
	`SELECT TIMESTAMP '2026-01-01 00:00:00' FROM dual`,
	`SELECT TIMESTAMP '2026-01-01 00:00:00 +00:00' FROM dual`,
	`SELECT INTERVAL '1' DAY FROM dual`,
	`SELECT INTERVAL '1-2' YEAR TO MONTH FROM dual`,
	`SELECT DATE 1 FROM dual`,
	// Bind variables.
	`SELECT :a FROM dual`,
	`SELECT :1 FROM dual`,
	`SELECT :"a b" FROM dual`,
	`SELECT :a :b FROM dual`,
	`SELECT : a FROM dual`,
	`SELECT 1 FROM dual WHERE dummy = :a`,
	// Comments and hints.
	`SELECT 1 FROM dual -- a`,
	`SELECT 1 /* a */ FROM dual`,
	`SELECT /* a */ 1 FROM dual`,
	`SELECT /*+ FULL(d) */ 1 FROM dual d`,
	`SELECT /*+ */ 1 FROM dual`,
	`SELECT 1 /* a FROM dual`,
	`SELECT 1 /* /* a */ */ FROM dual`,
	`SELECT 1 FROM dual REM a`,
	// Operators the lexer reads as one token.
	`SELECT 1 FROM dual WHERE 1 != 2`,
	`SELECT 1 FROM dual WHERE 1 ^= 2`,
	`SELECT 1 FROM dual WHERE 1 <> 2`,
	`SELECT 1 FROM dual WHERE 1 ~= 2`,
	`SELECT 1 FROM dual WHERE 1 < > 2`,
	`SELECT 1 FROM dual WHERE 1 ! = 2`,
	`SELECT 'a' || 'b' FROM dual`,
	`SELECT 'a' | | 'b' FROM dual`,
	`SELECT 1 FROM dual WHERE 1 >= 2 AND 1 <= 2`,
	// Case, white space and the end of the text.
	`select 1 from DUAL`,
	"SELECT\t1\tFROM\tdual",
	`SELECT 1 FROM dual;`,
	`SELECT 1 FROM dual;;`,
	`SELECT 1 FROM dual /`,
	`SELECT 1 FROM dual @`,
	`SELECT 1 FROM dual ~`,
	`SELECT 1 FROM dual \`,
	`SELECT 1 FROM dual ` + "`",
}

// corpusSuites are the corpora that are suite: rows of compat.tsv, in the
// order the rows were added: the lexical forms above, and the statements
// of testdata/coverage/reference_oracle.tsv and compat_oracle.tsv, which a
// test once asked a live Oracle about, with {run} written R1.
func corpusSuites(t *testing.T) []struct {
	name   string
	corpus []string
} {
	t.Helper()
	coverage := func(file string) []string {
		rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", file))
		var out []string
		for _, r := range rows {
			out = append(out, strings.ReplaceAll(r.Fields["sql"], "{run}", "R1"))
		}
		return out
	}
	return []struct {
		name   string
		corpus []string
	}{
		{"lexer", lexerCorpus},
		{"reference", coverage("reference_oracle.tsv")},
		{"manifest", coverage("compat_oracle.tsv")},
	}
}

// TestCompatibilitySuitesAreTheCorpora fails while the suite: rows of
// compat.tsv are not the corpora of corpusSuites: each statement once,
// under the first suite that holds it. A
// statement added to a corpus is added as a row with ? in both verdict
// cells, and judged by a recording.
func TestCompatibilitySuitesAreTheCorpora(t *testing.T) {
	f, err := compat.ReadFile(filepath.Join("testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	var want, got []string
	for _, s := range corpusSuites(t) {
		for _, stmt := range s.corpus {
			if !seen[stmt] {
				seen[stmt] = true
				want = append(want, "suite:"+s.name+"\t"+stmt)
			}
		}
	}
	for _, r := range f.Rows {
		if strings.HasPrefix(r.Rule, "suite:") {
			got = append(got, r.Rule+"\t"+r.Statement)
		}
	}
	if slices.Equal(got, want) {
		return
	}
	for i := 0; i < max(len(got), len(want)); i++ {
		g, w := "(none)", "(none)"
		if i < len(got) {
			g = got[i]
		}
		if i < len(want) {
			w = want[i]
		}
		if g != w {
			t.Fatalf("suite row %d of testdata/compat.tsv is %q, and the corpora want %q. A new statement is a row with ? in both verdict cells,\n  ?\t?\t-\tsuite:<name>\t<statement>\njudged by a recording:\n  OMNI_RECORD=1 go test ./oracle/parser/ -run '^TestCompatibility$'", i+1, g, w)
		}
	}
}
