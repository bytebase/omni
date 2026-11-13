package parser

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/tidb/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of TiDB's parser.y and its near-misses, with omni's verdict and
// TiDB's. It reads the file and starts no server. OMNI_RECORD=1 judges the
// rows that have no verdict yet on the pinned TiDB.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: filepath.Join("..", "testdata", "parser.y"),
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

// compatJudge asks the pinned TiDB, through PREPARE, whether it accepts each
// statement's syntax. It runs only in a recording, on a connection outside any
// tape. The statement is bound to a user variable and prepared from it, so it
// reaches the server byte for byte: rows hold invalid UTF-8, trailing spaces
// and a CR inside a literal. The variable lives on one connection, so the
// judge holds one.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	ctx := context.Background()
	conn, err := testcontainer.Judge(t).Conn(ctx)
	if err != nil {
		t.Fatalf("connect to TiDB: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	prepare := func(text string) error {
		if _, err := conn.ExecContext(ctx, "SET @omni_compat = ?", text); err != nil {
			// %v, not %w: a failed bind is no verdict, whatever its code.
			return fmt.Errorf("binding the statement: %v", err)
		}
		_, err := conn.ExecContext(ctx, "PREPARE omni_compat_stmt FROM @omni_compat")
		return err
	}
	return func(sql string) (string, error) {
		accepts, err := prepareVerdict(prepare(sql), func() error { return prepare(sql + compat.ProbeTail) })
		if err != nil {
			return "", err
		}
		if accepts {
			return compat.Accept, nil
		}
		return compat.Reject, nil
	}
}
