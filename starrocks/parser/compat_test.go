package parser

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/starrocks/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of StarRocks's StarRocks.g4 and its near-misses, with omni's
// verdict and StarRocks's. The test reads the file and starts no server.
// OMNI_RECORD=1 judges the rows that have no verdict yet on the pinned
// StarRocks.
//
// testdata/README.md says what StarRocks's judge needs.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "StarRocks.g4"),
		Lexer:   filepath.Join(td, "StarRocksLex.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict reads the row as StarRocks does: the whole text, which may
// hold more than one statement.
func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// compatJudge asks the pinned StarRocks whether its parser accepts each
// statement, sent behind compat.StarRocksProbe in one request, so that
// StarRocks parses the statement and runs none of it, and reads the answer
// with compat.StarRocksVerdict. It runs only in a recording.
func compatJudge(t testing.TB, _ enginetest.Pin) func(string) (string, error) {
	t.Helper()
	conn := testcontainer.Judge(t)
	return func(sql string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_, err := conn.ExecContext(ctx, compat.StarRocksProbe+";\n"+sql)
		return compat.StarRocksVerdict(err, starRocksMessage)
	}
}

// starRocksMessage reads a StarRocks answer for compat.StarRocksVerdict.
func starRocksMessage(err error) (string, bool) {
	var me *mysql.MySQLError
	if !errors.As(err, &me) {
		return "", false
	}
	return me.Message, true
}
