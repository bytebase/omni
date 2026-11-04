package parser

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
)

var update = flag.Bool("update", false, "rewrite the TiDB parser corpus from the current omni and engine verdicts")

// tidbCorpusVersion is the TiDB release the corpus's engine column was
// recorded against. It must name the release tidbImage pins: the engine column
// describes that server, not another one.
const tidbCorpusVersion = "v8.5.5"

// corpusMaxReported caps the per-statement failures a run prints. Closing a
// grammar family moves dozens of rows at once, and a thousand identical lines
// is not a report.
const corpusMaxReported = 25

var tidbCorpusPath = filepath.Join("testdata", "pingcap_parser_corpus.tsv")

// corpusCase is one recorded row of the corpus: a statement and the verdict
// each side gave it when the corpus was last recorded.
type corpusCase struct {
	line   int // 1-based line in the corpus file, for a failure message
	omni   bool
	engine bool
	sql    string
}

// loadCorpus reads the corpus, returning its comment header unchanged so an
// -update run can write it back.
//
// The file is split on '\n' rather than scanned line by line because a
// statement may end in spaces or hold a CR inside a string literal, and those
// bytes are the case. The statement is the third field, so a tab inside it
// would survive too; none has one today.
func loadCorpus(t testing.TB) (header []string, cases []corpusCase) {
	t.Helper()

	b, err := os.ReadFile(tidbCorpusPath)
	if err != nil {
		t.Fatalf("read %s: %v", tidbCorpusPath, err)
	}
	for i, raw := range strings.Split(string(b), "\n") {
		switch {
		case raw == "":
			continue
		case strings.HasPrefix(raw, "#"):
			header = append(header, raw)
			continue
		}
		fields := strings.SplitN(raw, "\t", 3)
		if len(fields) != 3 {
			t.Fatalf("%s:%d: want 3 tab-separated fields, got %d", tidbCorpusPath, i+1, len(fields))
		}
		omni, err := recordedVerdict(fields[0])
		if err != nil {
			t.Fatalf("%s:%d: omni column: %v", tidbCorpusPath, i+1, err)
		}
		engine, err := recordedVerdict(fields[1])
		if err != nil {
			t.Fatalf("%s:%d: engine column: %v", tidbCorpusPath, i+1, err)
		}
		cases = append(cases, corpusCase{line: i + 1, omni: omni, engine: engine, sql: fields[2]})
	}
	if len(cases) == 0 {
		t.Fatalf("%s holds no cases", tidbCorpusPath)
	}
	return header, cases
}

func recordedVerdict(field string) (bool, error) {
	switch field {
	case "accept":
		return true, nil
	case "reject":
		return false, nil
	}
	return false, fmt.Errorf("got %q, want \"accept\" or \"reject\"", field)
}

func verdictName(accepts bool) string {
	if accepts {
		return "accept"
	}
	return "reject"
}

// writeCorpus rewrites the corpus from the verdicts a run measured, keeping the
// comment header and the statement order.
func writeCorpus(t *testing.T, header []string, cases []corpusCase) {
	t.Helper()
	var b strings.Builder
	for _, line := range header {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	for _, c := range cases {
		fmt.Fprintf(&b, "%s\t%s\t%s\n", verdictName(c.omni), verdictName(c.engine), c.sql)
	}
	if err := os.WriteFile(tidbCorpusPath, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write %s: %v", tidbCorpusPath, err)
	}
	t.Logf("wrote %s (%d cases)", tidbCorpusPath, len(cases))
}

// corpusTally counts the ways omni and the engine can answer one statement.
// noVerdict is its own column: a statement the server never answered for is
// not agreement, and must never be added to one of the other four.
type corpusTally struct {
	total                    int
	bothAccept               int
	engineAcceptsOmniRejects int
	engineRejectsOmniAccepts int
	bothReject               int
	noVerdict                int
	recovered                int // verdicts that needed a reconnect or a restart
}

func (ct *corpusTally) add(omni, engine bool) {
	ct.total++
	switch {
	case omni && engine:
		ct.bothAccept++
	case !omni && engine:
		ct.engineAcceptsOmniRejects++
	case omni && !engine:
		ct.engineRejectsOmniAccepts++
	default:
		ct.bothReject++
	}
}

func (ct *corpusTally) addNoVerdict() {
	ct.total++
	ct.noVerdict++
}

// report prints every number in one block, so a run answers "how far along is
// the TiDB parser" without anything else being read.
func (ct corpusTally) report(t *testing.T, engineSource string) {
	t.Helper()
	engineAccepts := ct.bothAccept + ct.engineAcceptsOmniRejects
	omniAccepts := ct.bothAccept + ct.engineRejectsOmniAccepts
	t.Logf("TiDB %s corpus (%s), engine verdicts %s:", tidbCorpusVersion, tidbCorpusPath, engineSource)
	t.Logf("  total cases                      %5d", ct.total)
	t.Logf("  both accept                      %5d", ct.bothAccept)
	t.Logf("  TiDB accepts, omni rejects (gap) %5d", ct.engineAcceptsOmniRejects)
	t.Logf("  TiDB rejects, omni accepts       %5d", ct.engineRejectsOmniAccepts)
	t.Logf("  both reject                      %5d", ct.bothReject)
	t.Logf("  no verdict                       %5d", ct.noVerdict)
	if ct.recovered > 0 {
		t.Logf("  verdicts that needed a recovery  %5d", ct.recovered)
	}
	if engineAccepts > 0 {
		t.Logf("  omni parses %d of the %d statements TiDB accepts (%.1f%%)",
			ct.bothAccept, engineAccepts, 100*float64(ct.bothAccept)/float64(engineAccepts))
	}
	if omniAccepts > 0 {
		t.Logf("  omni over-accepts %d of the %d statements it parses (%.1f%%)",
			ct.engineRejectsOmniAccepts, omniAccepts,
			100*float64(ct.engineRejectsOmniAccepts)/float64(omniAccepts))
	}
}

// TestTiDBCorpusOmniVerdicts holds tidb/parser to the omni column of the
// corpus and prints how much of the TiDB-accepted surface it covers. It needs
// no container, so it is the gate a grammar change re-baselines: close a gap,
// rerun with -update, and the diff shows which statements the change bought.
//
// The engine column is read, never written, here: a run without a server
// cannot know what the engine does.
func TestTiDBCorpusOmniVerdicts(t *testing.T) {
	header, cases := loadCorpus(t)

	var tally corpusTally
	var drift int
	for i := range cases {
		c := &cases[i]
		_, err := Parse(c.sql)
		accepts := err == nil
		if accepts != c.omni {
			drift++
			if drift <= corpusMaxReported {
				t.Errorf("%s:%d: omni now %ss %q, corpus records %s (err=%v)",
					tidbCorpusPath, c.line, verdictName(accepts), c.sql, verdictName(c.omni), err)
			}
			c.omni = accepts
		}
		tally.add(accepts, c.engine)
	}
	if drift > corpusMaxReported {
		t.Errorf("%d more rows drifted from the recorded omni verdict; rerun with -update and read the diff", drift-corpusMaxReported)
	}
	tally.report(t, "as recorded")

	if *update {
		writeCorpus(t, header, cases)
	}
}

// TestTiDBCorpusEngineVerdicts measures both parsers against each other in one
// run: omni in process, TiDB v8.5.5 in the pinned container. It reports the
// four agreement quadrants and the no-verdict count, and holds the engine
// column of the corpus to what the server actually says, so bumping the image
// surfaces as drift to read rather than as a silent change to the gap numbers.
func TestTiDBCorpusEngineVerdicts(t *testing.T) {
	if want := "pingcap/tidb:" + tidbCorpusVersion; tidbImage != want {
		t.Fatalf("the corpus engine column was recorded against %s, but this package pins %s; rerun with -update and read the drift", want, tidbImage)
	}
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)

	header, cases := loadCorpus(t)

	var tally corpusTally
	var omniDrift, engineDrift int
	var unanswered []string
	for i := range cases {
		c := &cases[i]

		_, omniErr := Parse(c.sql)
		omniAccepts := omniErr == nil
		if omniAccepts != c.omni {
			omniDrift++
			c.omni = omniAccepts
		}

		previous := "<start of corpus>"
		if i > 0 {
			previous = cases[i-1].sql
		}
		engineAccepts, recovered, err := tc.syntaxVerdict(t, c.sql, previous)
		if recovered {
			tally.recovered++
		}
		if err != nil {
			tally.addNoVerdict()
			unanswered = append(unanswered, fmt.Sprintf("%s:%d: %v", tidbCorpusPath, c.line, err))
			// A server that has stopped answering will not start again, and
			// grinding through the rest records nothing.
			if len(unanswered) > 10 {
				t.Errorf("gave up after %d statements with no verdict", len(unanswered))
				break
			}
			continue
		}

		if engineAccepts != c.engine {
			engineDrift++
			if engineDrift <= corpusMaxReported {
				t.Errorf("%s:%d: TiDB %s now %ss %q, corpus records %s",
					tidbCorpusPath, c.line, tidbCorpusVersion, verdictName(engineAccepts), c.sql, verdictName(c.engine))
			}
			c.engine = engineAccepts
		}
		tally.add(omniAccepts, engineAccepts)
	}

	if engineDrift > corpusMaxReported {
		t.Errorf("%d more rows drifted from the recorded engine verdict; the pinned image decides this column, so rerun with -update only after reading the drift",
			engineDrift-corpusMaxReported)
	}
	if omniDrift > 0 {
		t.Errorf("%d rows drifted from the recorded omni verdict; TestTiDBCorpusOmniVerdicts names them", omniDrift)
	}
	tally.report(t, "live, via PREPARE")
	for _, u := range unanswered {
		t.Errorf("no verdict: %s", u)
	}

	// A run that lost verdicts measured nothing it can record: the rows it
	// never reached would keep stale values beside rows it did.
	if *update && len(unanswered) == 0 {
		writeCorpus(t, header, cases)
	}
}

// assertTiDBVersion fails unless the container runs the release the corpus was
// recorded on. VERSION() reports the MySQL version TiDB claims with the TiDB
// release appended, as in "8.0.11-TiDB-v8.5.5".
func assertTiDBVersion(t *testing.T, tc *tidbContainer) {
	t.Helper()
	var version string
	if err := tc.db.QueryRowContext(tc.ctx, "SELECT VERSION()").Scan(&version); err != nil {
		t.Fatalf("read the TiDB version: %v", err)
	}
	if !strings.HasSuffix(version, "-TiDB-"+tidbCorpusVersion) {
		t.Fatalf("the server reports %q, but the corpus was recorded against TiDB %s", version, tidbCorpusVersion)
	}
}

// syntaxVerdict asks TiDB whether it accepts a statement's syntax, recovering
// from a lost connection or a lost server before giving up. recovered reports
// whether it had to. A non-nil error means the server gave no verdict at all,
// which is not a rejection and not agreement; it names the statement that went
// unanswered and the one before it, because the statement that destroys a run
// is the one that ran last, not the one that reported the loss.
func (tc *tidbContainer) syntaxVerdict(t *testing.T, stmt, previous string) (accepts, recovered bool, err error) {
	t.Helper()
	accepts, err = tc.prepareSyntax(stmt)
	if err == nil {
		return accepts, false, nil
	}
	lost := err
	t.Logf("no verdict for %q (preceded by %q): %v; recovering", stmt, previous, lost)

	// A statement can drop the connection it ran on without taking the server
	// with it. Rebuilding the pool costs one round trip and buys back every
	// later verdict.
	if err = tc.reconnect(); err == nil {
		if accepts, err = tc.prepareSyntax(stmt); err == nil {
			return accepts, true, nil
		}
	}

	// The server itself is gone, and no pool reaches a process that has exited.
	if err = tc.restart(); err == nil {
		if accepts, err = tc.prepareSyntax(stmt); err == nil {
			return accepts, true, nil
		}
	}
	return false, true, fmt.Errorf("no verdict for %q (preceded by %q) after a reconnect and a restart: %w (first failure: %v)",
		stmt, previous, err, lost)
}

// TestTiDBCorpusHarnessSurvivesALostServer pins the resilience the corpus run
// depends on. An earlier measuring pass ran the corpus by executing it and
// stopped getting verdicts partway through; `shutdown` is in the corpus, it
// succeeds, and it takes the server with it. The corpus harness prepares
// rather than executes, so it never reaches that state, but a container can
// still die under a test for its own reasons, and a run that then recorded
// silence as agreement would report a gap count about Docker.
//
// The test puts the harness in both states deliberately: a connection that is
// gone, and a server that is gone. Each must come back with a verdict and say
// it had to recover. The restart leaves a fresh server behind; every other
// container test in this package creates the database it needs, so none of
// them depends on what this one discards.
func TestTiDBCorpusHarnessSurvivesALostServer(t *testing.T) {
	tc := startTiDB(t)

	// A lost connection. database/sql may well heal this itself; what matters
	// is that a verdict comes back rather than an error that reads as one.
	_, _ = tc.db.ExecContext(tc.ctx, "KILL CONNECTION_ID()")
	accepts, _, err := tc.syntaxVerdict(t, "SELECT 1", "KILL CONNECTION_ID()")
	if err != nil {
		t.Fatalf("after a lost connection: %v", err)
	}
	if !accepts {
		t.Errorf("after a lost connection: TiDB rejected SELECT 1")
	}

	// A lost server. TiDB stops gracefully, so the process keeps answering for
	// a moment after the statement returns; wait for it to go, or the verdict
	// that follows would come from the server this test meant to destroy.
	_, _ = tc.db.ExecContext(tc.ctx, "SHUTDOWN")
	waitForTiDBGone(t, tc)
	accepts, recovered, err := tc.syntaxVerdict(t, "SELECT 1", "SHUTDOWN")
	if err != nil {
		t.Fatalf("after a lost server: %v", err)
	}
	if !recovered {
		t.Errorf("after a lost server: got a verdict without recovering, so the server never went away")
	}
	if !accepts {
		t.Errorf("after a lost server: TiDB rejected SELECT 1")
	}

	// The replacement server must be the pinned release, or every later
	// verdict in the run describes a different engine.
	assertTiDBVersion(t, tc)
}

// waitForTiDBGone blocks until the pool can no longer reach the server.
func waitForTiDBGone(t *testing.T, tc *tidbContainer) {
	t.Helper()
	deadline := time.Now().Add(60 * time.Second)
	for {
		if err := tc.db.PingContext(tc.ctx); err != nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("the server still answers 60s after SHUTDOWN")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// TestPrepareSyntaxHoldsOneConnection pins the connection the bind and the
// PREPARE share. The statement reaches the server as a session variable, and a
// session variable lives on one connection: a PREPARE that reaches a different
// one reads a variable nothing ever set, and the server answers that with an
// error prepareVerdict cannot tell from a verdict about the statement. A
// corpus row recorded from one of those describes a statement the server never
// saw, and nothing later in the run can detect it.
//
// A pool told to keep no idle connection is the cheapest way to force the two
// apart. In a run it is a connection dropped between them that does it:
// database/sql retries a statement whose connection went bad on a fresh one,
// where the variable is unset, so the loss would be recorded as a verdict
// instead of reaching the recovery in syntaxVerdict.
func TestPrepareSyntaxHoldsOneConnection(t *testing.T) {
	tc := startTiDB(t)

	// The pool is shared with every other test in this package; two is the
	// database/sql default this package relies on everywhere else.
	defer tc.db.SetMaxIdleConns(2)
	tc.db.SetMaxIdleConns(0)

	for _, c := range []struct {
		stmt string
		want bool
	}{
		{"SELECT 1", true},
		{"SELEKT 1", false},
	} {
		accepts, err := tc.prepareSyntax(c.stmt)
		if err != nil {
			t.Fatalf("no verdict for %q: %v", c.stmt, err)
		}
		if accepts != c.want {
			t.Errorf("%q: accepts=%v, want %v; the bind and the PREPARE reached different connections", c.stmt, accepts, c.want)
		}
	}
}

// prepareSyntax reports whether TiDB accepts the statement's syntax, without
// running it: PREPARE parses the text and stops there. A non-nil error means
// the server gave no verdict at all.
//
// Executing the corpus cannot measure this, because statements in it destroy
// the run. `shutdown` succeeds and stops the server, so nothing after it gets
// a verdict and no reconnect reaches a process that has exited.
// `SET PASSWORD = 'password'` succeeds and changes the connecting user's
// password, so the next new connection is refused and the run dies as soon as
// anything replaces the one it authenticated. `load stats '...'` and the
// PLAN REPLAYER file forms make the server ask the client for a local file,
// which the driver refuses with an error that is not a server verdict at all.
// PREPARE runs none of them.
//
// The statement travels as a bound parameter, so nothing in it has to be
// escaped and no byte in it can end the literal that carries it. The corpus
// holds rows that are not valid UTF-8, rows ending in spaces, and rows with a
// CR inside a string literal; all of them reach the server unchanged this way.
//
// The three statements run on one connection the caller holds, because the
// variable that carries the statement lives on the connection that set it.
// Taken from the pool separately they can land on different connections, and
// then PREPARE reads a variable nothing set and answers 1064 -- a rejection
// recorded against a statement the server never saw. Holding the connection
// turns a loss between them into an error instead, which is what syntaxVerdict
// recovers from.
func (tc *tidbContainer) prepareSyntax(stmt string) (bool, error) {
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		return false, fmt.Errorf("taking a connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(tc.ctx, "SET @omni_corpus = ?", stmt); err != nil {
		return false, fmt.Errorf("binding the statement: %w", err)
	}
	_, err = conn.ExecContext(tc.ctx, "PREPARE omni_corpus_stmt FROM @omni_corpus")
	if err == nil {
		_, _ = conn.ExecContext(tc.ctx, "DEALLOCATE PREPARE omni_corpus_stmt")
	}
	return prepareVerdict(err)
}

// prepareVerdict reads the error PREPARE returned. ER_PARSE_ERROR (1064) is
// the parser's rejection. Any other server error means the text parsed and
// failed later: 1295 for a statement the prepared-statement protocol does not
// carry, 1146 for a table that does not exist.
//
// Anything that is not a server error — a dropped connection, a timeout, a
// driver refusing a local-file request — is no verdict, and comes back as an
// error so the caller recovers or fails rather than records agreement.
func prepareVerdict(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	var myErr *gomysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number != 1064, nil
	}
	return false, err
}
