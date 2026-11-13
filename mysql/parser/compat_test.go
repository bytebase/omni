package parser

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of MySQL's sql_yacc.yy and its near-misses, with omni's verdict
// and MySQL's. It reads the file and starts no server. The grammar is under
// the GPL, so the repo does not hold it: the test downloads it once
// (../testdata/sources.tsv). OMNI_RECORD=1 judges the rows that have no
// verdict yet on the pinned MySQL.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: enginetest.Fetched(t, filepath.Join("..", "testdata", "sources.tsv"), "sql_yacc.yy"),
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

// compatPrepareTimeout bounds one PREPARE, so a statement the server hangs on
// gets no verdict and the recording goes on.
const compatPrepareTimeout = time.Minute

// compatJudge asks the package's MySQL, through PREPARE, whether it accepts
// each statement's syntax. It fails t unless the server is the pinned image
// and release.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	ctx := context.Background()
	d, err := testcontainer.MySQL80.Database(ctx, "compat")
	if err != nil {
		t.Fatalf("MySQL: %v", err)
	}
	t.Cleanup(func() { d.Close(ctx) })
	if err := testcontainer.MySQL80.CheckPin(ctx, pin); err != nil {
		t.Fatalf("%v", err)
	}
	prepare := func(text string) error {
		ctx, cancel := context.WithTimeout(ctx, compatPrepareTimeout)
		defer cancel()
		_, err := d.DB.ExecContext(ctx, "PREPARE omni_compat FROM ?", text)
		return err
	}
	// The verdict reads 1064 as the parser's rejection, so a release that
	// answers a syntax error otherwise must fail here and not be misread.
	if code, ok := mysqlErrorCode(prepare("SELECT FROM")); !ok || code != 1064 {
		t.Fatalf("MySQL no longer answers PREPARE of a syntax error with 1064: %v", prepare("SELECT FROM"))
	}
	return func(stmt string) (string, error) {
		accepts, err := compat.PrepareVerdict(prepare(stmt), mysqlErrorCode, func() error {
			return prepare(stmt + compat.ProbeTail)
		})
		switch {
		case err != nil:
			return "", err
		case accepts:
			return compat.Accept, nil
		default:
			return compat.Reject, nil
		}
	}
}

// mysqlErrorCode returns the server's error number, and false for an error
// the server did not send.
func mysqlErrorCode(err error) (uint16, bool) {
	var me *gomysql.MySQLError
	if !errors.As(err, &me) {
		return 0, false
	}
	return me.Number, true
}
