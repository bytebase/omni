package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of OceanBase's sql_parser_mysql_mode.y and its near-misses, with
// omni's verdict and OceanBase's. It reads the file and starts no server.
// The grammar is under the Mulan PubL, so the repo does not hold it: the test
// downloads it once (../testdata/sources.tsv). OMNI_RECORD=1 judges the rows
// that have no verdict yet on the pinned OceanBase.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		PinRole: "parser",
		Grammar: enginetest.Fetched(t, filepath.Join("..", "testdata", "sources.tsv"), "sql_parser_mysql_mode.y"),
		Tokens:  filepath.Join("..", "testdata", "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// compatJudge asks the package's OceanBase, through PREPARE, whether it
// accepts each statement's syntax. It fails t unless the server is the
// pinned image and release.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	ctx := context.Background()
	dsn, err := testcontainer.OceanBase.DSN(ctx)
	if err != nil {
		t.Fatalf("OceanBase: %v", err)
	}
	if err := testcontainer.OceanBase.CheckPin(ctx, pin); err != nil {
		t.Fatalf("%v", err)
	}
	// The judge asks the server itself: its answers go to compat.tsv, not to
	// a tape.
	db, err := testcontainer.OpenOceanBase(dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	// PREPARE reads a user variable, which lives on the connection that set it.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	j := &compatJudgeOceanBase{t: t, ctx: ctx, db: db}
	// The verdict tells the two kinds of 1064 apart by their message (see
	// oceanBaseErrorCode), so a release that words them otherwise must fail
	// here and not be misread.
	if code, ok := oceanBaseErrorCode(j.prepare("SELECT FROM")); !ok || code != 1064 {
		t.Fatalf("OceanBase no longer answers a syntax error as %d %q: got %v", 1064, oceanBaseSyntaxError, j.prepare("SELECT FROM"))
	}
	if code, ok := oceanBaseErrorCode(j.prepare("EXECUTE a")); !ok || code != oceanBaseParseErrorAfterParsing {
		t.Fatalf("OceanBase no longer answers PREPARE of an EXECUTE with a bare 1064: got %v", j.prepare("EXECUTE a"))
	}
	if accepts, err := j.ask("DO a zzz"); err != nil || accepts {
		t.Fatalf("OceanBase judge rev 2 must reject DO a zzz (1149 with probe 1064 at line 1): accepts=%v err=%v", accepts, err)
	}
	if accepts, err := j.ask("SELECT RANK() OVER(ROWS 1 FOLLOWING)"); err != nil || !accepts {
		t.Fatalf("OceanBase judge rev 2 must accept post-parse 1149 with probe 1064 at line 2: accepts=%v err=%v", accepts, err)
	}
	if accepts, err := j.ask("SELECT /*+ TRANS_PARAM('s') */ a"); err != nil || !accepts {
		t.Fatalf("OceanBase judge rev 2 must accept ignored-hint 1054: accepts=%v err=%v", accepts, err)
	}
	return j.verdict
}

type compatJudgeOceanBase struct {
	t   testing.TB
	ctx context.Context
	db  *sql.DB
}

// compatPrepareTimeout bounds one PREPARE, so a statement the server hangs on
// gets no verdict and the recording goes on.
const compatPrepareTimeout = time.Minute

// prepare asks the server to prepare text, without running it. The text is
// bound to a user variable, so its quoting stays the server's business. An
// error from the binding is no answer about the text, and is returned as
// one the server did not send.
func (j *compatJudgeOceanBase) prepare(text string) error {
	ctx, cancel := context.WithTimeout(j.ctx, compatPrepareTimeout)
	defer cancel()
	if _, err := j.db.ExecContext(ctx, "SET @omni_compat = ?", text); err != nil {
		return fmt.Errorf("binding the statement: %v", err)
	}
	_, err := j.db.ExecContext(ctx, "PREPARE omni_compat_stmt FROM @omni_compat")
	if err == nil {
		_, _ = j.db.ExecContext(ctx, "DEALLOCATE PREPARE omni_compat_stmt")
	}
	return err
}

// ask reads PREPARE's answer through compat.PrepareVerdict. OceanBase
// refuses to prepare a text of two statements before it parses either
// (1235, with the probe's tail too, which reads as a rejection), so the
// tail stays inside the statement, and one statement is all a row can hold.
//
// A 1149 ("You have an error in your SQL syntax", with no near clause) can
// come either from a post-parse resolver check (SELECT RANK() OVER(ROWS 1
// FOLLOWING), whose probe fails at line 2 near '@@@ (((') or from a parse-time
// syntax error (DO a zzz, whose probe fails at line 1 near 'a zzz\n @@@ (((').
// When the first error is 1149 and the probe's 1064 is not reported on the
// probe tail's own line, the parser failed earlier in the statement and the
// row is rejected.
func (j *compatJudgeOceanBase) ask(stmt string) (bool, error) {
	return oceanBaseAskVerdict(stmt, func() error {
		return j.prepare(stmt)
	}, func() error {
		return j.prepare(stmt + compat.ProbeTailInStatement)
	})
}

func oceanBaseAskVerdict(stmt string, first, probe func() error) (bool, error) {
	firstErr := first()
	firstCode, firstOK := oceanBaseErrorCode(firstErr)
	probeLine := fmt.Sprintf(" at line %d", strings.Count(stmt, "\n")+2)
	return compat.PrepareVerdict(firstErr, oceanBaseErrorCode, func() error {
		tailErr := probe()
		var myErr *gomysql.MySQLError
		if firstOK && firstCode == 1149 &&
			errors.As(tailErr, &myErr) && myErr.Number == 1064 &&
			strings.HasPrefix(myErr.Message, oceanBaseSyntaxError) &&
			!strings.HasSuffix(myErr.Message, probeLine) {
			return firstErr
		}
		return tailErr
	})
}

// verdict judges one statement, and survives the statements that crash the
// server: SET ROLE ALL EXCEPT ALL is a null dereference in the observer,
// which the container outlives.
//
// When the server dies, the judge starts it again and asks the same
// statement again, because a server can die of something else, and a
// statement asked while it was down proves nothing. A statement that kills
// the server again is judged by the probe's question alone. Had it a syntax
// error, the parser would have answered 1064 and nothing would have
// crashed; so when the parser reads all of it and reports its syntax error
// at the tail, the statement parses, the crash comes after parsing, and
// that is acceptance. A statement that kills the server with the tail too
// gets no verdict.
func (j *compatJudgeOceanBase) verdict(stmt string) (string, error) {
	var noAnswer error
	for attempt := 0; attempt < 2; attempt++ {
		accepts, err := j.ask(stmt)
		if err == nil {
			return compatVerdict(accepts), nil
		}
		noAnswer = err
		if !j.recover(stmt) {
			return "", noAnswer // the server is up: a dropped connection, a timeout
		}
	}
	tailErr := j.prepare(stmt + compat.ProbeTailInStatement)
	var myErr *gomysql.MySQLError
	if errors.As(tailErr, &myErr) && myErr.Number == 1064 && strings.HasPrefix(myErr.Message, oceanBaseSyntaxError) && strings.Contains(myErr.Message, "(((") {
		j.t.Logf("crashes OceanBase after it parses, judged by the probe alone: %q", stmt)
		return compat.Accept, nil
	}
	j.recover(stmt)
	return "", fmt.Errorf("crashes the server, and the probe alone answers %v", tailErr)
}

// recover reports whether the server died, and starts it again if so. A
// lost connection alone does not say: the observer takes a moment to die,
// so recover waits until the server either answers or its process is gone.
func (j *compatJudgeOceanBase) recover(stmt string) (crashed bool) {
	time.Sleep(500 * time.Millisecond)
	deadline := time.Now().Add(2 * time.Minute)
	for {
		alive, err := testcontainer.OceanBase.Alive(j.ctx)
		if err != nil {
			j.t.Fatalf("after %q: %v", stmt, err)
		}
		if !alive {
			if err := testcontainer.OceanBase.Restart(j.ctx); err != nil {
				j.t.Fatalf("OceanBase died on %q and does not start again: %v", stmt, err)
			}
			return true
		}
		ctx, cancel := context.WithTimeout(j.ctx, 5*time.Second)
		err = j.db.PingContext(ctx)
		cancel()
		if err == nil {
			return false
		}
		if time.Now().After(deadline) {
			j.t.Fatalf("after %q OceanBase neither answers nor exits: %v", stmt, err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func compatVerdict(accepts bool) string {
	if accepts {
		return compat.Accept
	}
	return compat.Reject
}

// oceanBaseParseErrorAfterParsing stands for the 1064 OceanBase answers with
// the bare message "Parse error". It is OceanBase's own number for the
// error (OB_ERR_PARSE_SQL, 5001), which never travels: the server sends it
// as 1064.
const oceanBaseParseErrorAfterParsing = 5001

// oceanBaseSyntaxError starts the message of the 1064 OceanBase's parser
// raises, which goes on to quote the text near the error.
const oceanBaseSyntaxError = "You have an error in your SQL syntax"

// oceanBaseErrorCode returns the server's error number, and false for an
// error the server did not send.
//
// OceanBase answers 1064 in two ways. Its parser's syntax error says where:
// "You have an error in your SQL syntax; ... near 'x' at line 1". The same
// number with only "Parse error" comes from the layers after the parser,
// for an option value they cannot read (CREATE RESOURCE UNIT a MIN_CPU 's'),
// and from PREPARE itself, which does not carry PREPARE, EXECUTE or
// DEALLOCATE, though the server takes each when it is sent directly. Such
// a statement has parsed: with the probe's tail after it, the answer is the
// parser's syntax error at the tail. So the bare form gets a number of its
// own and the probe decides it, like any other error raised after parsing.
func oceanBaseErrorCode(err error) (uint16, bool) {
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) {
		return 0, false
	}
	if myErr.Number == 1064 && !strings.HasPrefix(myErr.Message, oceanBaseSyntaxError) {
		return oceanBaseParseErrorAfterParsing, true
	}
	return myErr.Number, true
}

func TestOceanBaseAskVerdict(t *testing.T) {
	syntax1149 := &gomysql.MySQLError{
		Number:  1149,
		Message: oceanBaseSyntaxError + "; check the manual that corresponds to your OceanBase version for the right syntax to use",
	}
	unknownCol1054 := &gomysql.MySQLError{
		Number:  1054,
		Message: "Unknown column 'a' in 'field list'",
	}
	probeLine1 := &gomysql.MySQLError{
		Number:  1064,
		Message: oceanBaseSyntaxError + "; check the manual that corresponds to your OceanBase version for the right syntax to use near 'a zzz\n @@@ (((' at line 1",
	}
	probeLine2 := &gomysql.MySQLError{
		Number:  1064,
		Message: oceanBaseSyntaxError + "; check the manual that corresponds to your OceanBase version for the right syntax to use near '@@@ (((' at line 2",
	}

	cases := []struct {
		name  string
		stmt  string
		first error
		probe error
		want  bool
	}{
		{
			name:  "1149 with early probe 1064 at line 1 rejects",
			stmt:  "DO a zzz",
			first: syntax1149,
			probe: probeLine1,
			want:  false,
		},
		{
			name:  "1149 with probe 1064 at line 2 accepts",
			stmt:  "SELECT RANK() OVER(ROWS 1 FOLLOWING)",
			first: syntax1149,
			probe: probeLine2,
			want:  true,
		},
		{
			name:  "1054 ignored-hint with line 1 hint syntax error still accepts",
			stmt:  "SELECT /*+ TRANS_PARAM('s') */ a",
			first: unknownCol1054,
			probe: probeLine1,
			want:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := oceanBaseAskVerdict(tc.stmt, func() error { return tc.first }, func() error { return tc.probe })
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got accepts=%v, want %v", got, tc.want)
			}
		})
	}
}
