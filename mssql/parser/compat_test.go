package parser

import (
	"context"
	"path/filepath"
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

// compatJudge asks the pinned SQL Server whether its parser accepts each
// statement, sent as a batch of its own under SET PARSEONLY ON, which runs
// nothing, and read by compat.ParseOnlyVerdict. It runs only in a
// recording, outside any tape.
func compatJudge(t testing.TB, _ enginetest.Pin) func(string) (string, error) {
	t.Helper()
	conn := testcontainer.Judge(t)
	return func(sql string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return testcontainer.Verdict(ctx, conn, testcontainer.JudgeDatabase, sql)
	}
}
