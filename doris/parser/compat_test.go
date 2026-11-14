package parser

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/doris/internal/testcontainer"
	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of Doris's DorisParser.g4 and of its legacy CUP parser
// (converted to DorisLegacyParser.g4), and their near-misses, with omni's
// verdict and Doris's. The test reads the file and starts no server.
// OMNI_RECORD=1 judges the rows that have no verdict yet on the pinned Doris.
//
// testdata/README.md says what Doris's judge needs.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "DorisParser.g4"),
		Lexer:   filepath.Join(td, "DorisLexer.g4"),
		More:    []string{filepath.Join(td, "DorisLegacyParser.g4")},
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict reads the row as Doris does: the whole text, which may hold
// more than one statement.
func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// compatJudge asks the pinned Doris whether its parsers accept each
// statement, inside a transaction, where Doris parses a statement and runs
// none but INSERT, UPDATE, DELETE and the transaction's own, and reads the
// answer with compat.DorisVerdict. It runs only in a recording.
func compatJudge(t testing.TB, _ enginetest.Pin) func(string) (string, error) {
	t.Helper()
	conn := testcontainer.Judge(t)
	return func(sql string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		open := "BEGIN"
		if parsedWhenRun.MatchString(sql) {
			open = "ROLLBACK"
		}
		if _, err := conn.ExecContext(ctx, open); err != nil {
			return "", err
		}
		_, err := conn.ExecContext(ctx, sql)
		return compat.DorisVerdict(err, dorisMessage)
	}
}

// parsedWhenRun matches the two statements Doris's ANTLR parser builds a
// command for that hands the text to the CUP parser only when it runs:
// LoadCommand (LOAD LABEL) and CreatePolicyCommand (CREATE ROW POLICY).
// A transaction refuses them before that, so the judge sends them outside
// one, where they run: on names no table has, a load fails and a row policy
// is created.
var parsedWhenRun = regexp.MustCompile(`(?is)^\s*(LOAD\s+LABEL|CREATE\s+ROW\s+POLICY)\b`)

// dorisMessage reads a Doris answer for compat.DorisVerdict.
func dorisMessage(err error) (string, bool) {
	var me *mysql.MySQLError
	if !errors.As(err, &me) {
		return "", false
	}
	return me.Message, true
}
