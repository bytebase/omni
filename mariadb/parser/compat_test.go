package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/mariadb/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of MariaDB's sql_yacc.yy and its near-misses, with omni's
// verdict and MariaDB's. It reads the file and starts no server. The grammar
// is under the GPL, so the repo does not hold it: the test downloads it once
// (../testdata/sources.tsv). The file holds the sql_mode=ORACLE grammar too,
// which omni does not parse: the rows are those of %ifdef MARIADB.
// OMNI_RECORD=1 judges the rows that have no verdict yet on the pinned
// MariaDB.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: enginetest.Fetched(t, filepath.Join("..", "testdata", "sources.tsv"), "sql_yacc.yy"),
		Defined: []string{"MARIADB"},
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

// compatJudge asks the pinned MariaDB, through PREPARE, whether it accepts
// each statement's syntax. testdata/README.md says what the answer can and
// cannot tell.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	db := testcontainer.Pinned(t, pin)
	// A prepared statement lives on the connection that prepared it.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	j := &compatJudgeMariaDB{ctx: context.Background(), db: db, answers: map[uint16]int{}}
	// The verdict rests on these three answers, so a release that gives
	// others must fail here and not be misread.
	for _, probe := range []struct {
		stmt string
		code uint16
	}{
		{"SELECT FROM", 1064},
		{"EXECUTE a", 1295},
		{"SELECT 1" + compat.ProbeTail, 1064},
	} {
		if code, ok := mariaDBErrorCode(j.prepare(probe.stmt)); !ok || code != probe.code {
			t.Fatalf("MariaDB no longer answers PREPARE of %q with %d: got %v", probe.stmt, probe.code, j.prepare(probe.stmt))
		}
	}
	t.Cleanup(func() { t.Logf("%s", j.report()) })
	return j.verdict
}

type compatJudgeMariaDB struct {
	ctx context.Context
	db  *sql.DB
	// answers counts the statements by the error PREPARE gave first, 0 for
	// none, and lookups those a name the server does not have answered.
	answers map[uint16]int
	lookups []string
}

// prepare asks the server to prepare text, without running it.
func (j *compatJudgeMariaDB) prepare(text string) error {
	ctx, cancel := context.WithTimeout(j.ctx, compatPrepareTimeout)
	defer cancel()
	_, err := j.db.ExecContext(ctx, "PREPARE omni_compat_stmt FROM ?", text)
	if err == nil {
		_, _ = j.db.ExecContext(ctx, "DEALLOCATE PREPARE omni_compat_stmt")
	}
	return err
}

// verdict reads PREPARE's answer as prepareParses does, so a row and a
// hand-written case of this package are judged alike: 1064 rejects, a name
// the server does not have counts as read, and compat.PrepareVerdict's probe
// decides every other error.
func (j *compatJudgeMariaDB) verdict(stmt string) (string, error) {
	first := j.prepare(stmt)
	code, sent := mariaDBErrorCode(first)
	if first != nil && !sent {
		return "", first
	}
	accepts := mariaDBLookupErrors[code]
	if accepts {
		j.lookups = append(j.lookups, fmt.Sprintf("%d %s", code, stmt))
	} else {
		var err error
		accepts, err = compat.PrepareVerdict(first, mariaDBErrorCode, func() error {
			return j.prepare(stmt + compat.ProbeTail)
		})
		if err != nil {
			return "", err
		}
	}
	j.answers[code]++
	if accepts {
		return compat.Accept, nil
	}
	return compat.Reject, nil
}

// mariaDBLookupErrors are the answers read as acceptance without asking the
// probe: a name MariaDB looks up while it parses and does not have. They
// are about what the server holds and not about the shape of the text, and
// omni keeps no such list. A character set (1115), a plugin or loadable
// function (1128), a system variable (1193, 1272, 4081), a collation (1253,
// 1273), a storage engine (1286), and a form the build leaves out (1289,
// SHOW PROCEDURE CODE without debugging).
var mariaDBLookupErrors = map[uint16]bool{
	1115: true, 1128: true, 1193: true, 1253: true, 1272: true,
	1273: true, 1286: true, 1289: true, 4081: true,
}

// report says how the server answered the statements this run judged.
func (j *compatJudgeMariaDB) report() string {
	var b strings.Builder
	codes := make([]uint16, 0, len(j.answers))
	for c := range j.answers {
		codes = append(codes, c)
	}
	slices.SortFunc(codes, func(a, b uint16) int { return j.answers[b] - j.answers[a] })
	b.WriteString("PREPARE's first answer, by error code (0 is none):")
	for _, c := range codes {
		fmt.Fprintf(&b, " %d×%d", c, j.answers[c])
	}
	fmt.Fprintf(&b, "\n%d statements were answered for a name the server does not have, and read as accepted:\n  %s",
		len(j.lookups), strings.Join(j.lookups, "\n  "))
	return b.String()
}

// mariaDBErrorCode returns the server's error number, and false for an error
// the server did not send.
func mariaDBErrorCode(err error) (uint16, bool) {
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) {
		return 0, false
	}
	return myErr.Number, true
}
