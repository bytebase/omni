// Package testcontainer is how mssql's tests reach SQL Server and ScriptDOM:
// through a tape (enginetest/tape) that replays the pinned server's answers,
// or, in a recording (OMNI_RECORD=1), through the pinned server itself. That
// server is one testcontainer per pin role and test binary, started by the
// first test that records on it and removed by Terminate, which each
// package's TestMain calls.
package testcontainer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/testcontainers/testcontainers-go"
	tcmssql "github.com/testcontainers/testcontainers-go/modules/mssql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

// Pins is mssql/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin role of the SQL Server the tests ask.
const Default = "default"

// password is the sa password of a server this package starts.
const password = "Str0ngPa$$w0rd!"

// JudgeDatabase is the database the judge's session is in. A generated row
// names everything a, so the judge's database holds a table a with a column
// a: a name the parser looks up then resolves, and the parser reads on.
const JudgeDatabase = "a"

// Judge returns a session on the default server itself, outside any tape,
// for the judge of compat.tsv, which runs only in a recording. Its database
// is JudgeDatabase.
//
// It is a session and not a pool: database/sql resets a pooled connection
// before it hands it out again, and go-mssqldb then sends the next batch
// with RESETCONNECTION, which clears SET options and the database context.
func Judge(t testing.TB) *sql.Conn {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("testcontainer.Judge outside a recording: a replay reaches no server")
	}
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		t.Fatalf("%v", err)
	}
	in := started(t, pin)
	ctx := context.Background()
	db := sql.OpenDB(in.connector)
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("connect to SQL Server: %v", err)
	}
	t.Cleanup(func() {
		conn.Close()
		db.Close()
	})
	for _, stmt := range []string{
		"IF DB_ID('" + JudgeDatabase + "') IS NOT NULL BEGIN ALTER DATABASE " + quoteName(JudgeDatabase) + " SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE " + quoteName(JudgeDatabase) + " END",
		"CREATE DATABASE " + quoteName(JudgeDatabase),
		"USE " + quoteName(JudgeDatabase),
		"CREATE TABLE a (a INT)",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	return conn
}

// ParseOnly sends stmt to the session as a batch of its own under SET
// PARSEONLY ON, after declarations, and returns what the batch returned. The
// batch ends in a newline: go-mssqldb sends a text that looks like a
// procedure name (GO, a, BREAK) as a call of that procedure, not a batch.
// SET PARSEONLY, SET QUOTED_IDENTIFIER and SET FIPS_FLAGGER take effect
// while a batch is parsed, so a statement that sets one changes how the next
// is read, and so may USE: each check puts all four back as the session had
// them, in database.
func ParseOnly(ctx context.Context, conn *sql.Conn, database, declarations, stmt string) error {
	if _, err := conn.ExecContext(ctx, "SET PARSEONLY ON"); err != nil {
		return fmt.Errorf("SET PARSEONLY ON: %w", err)
	}
	_, parseErr := conn.ExecContext(ctx, declarations+stmt+"\n")
	restore := "SET PARSEONLY OFF\nSET QUOTED_IDENTIFIER ON\nSET FIPS_FLAGGER OFF\nUSE " + quoteName(database)
	if _, err := conn.ExecContext(ctx, restore); err != nil {
		return fmt.Errorf("restore the session after %q: %w", stmt, err)
	}
	return parseErr
}

// quoteName writes a name as a bracketed identifier.
func quoteName(name string) string { return "[" + strings.ReplaceAll(name, "]", "]]") + "]" }

// Verdict is compat.ParseOnlyVerdict on stmt, asked of the session, which
// is in database.
func Verdict(ctx context.Context, conn *sql.Conn, database, stmt string) (string, error) {
	return compat.ParseOnlyVerdict(ParseOnly(ctx, conn, database, "", stmt), ServerErrors, func(declarations string) error {
		return ParseOnly(ctx, conn, database, declarations, stmt)
	})
}

// ServerErrors reads go-mssqldb's error for compat.ParseOnlyVerdict: the
// errors SQL Server reported for a batch, in order. It is false for an error
// the server did not send.
func ServerErrors(err error) ([]compat.SQLServerError, bool) {
	var e mssql.Error
	if !errors.As(err, &e) {
		return nil, false
	}
	all := e.All
	if len(all) == 0 {
		all = []mssql.Error{e}
	}
	out := make([]compat.SQLServerError, len(all))
	for i, x := range all {
		out[i] = compat.SQLServerError{Number: x.Number, Message: x.Message}
	}
	return out, true
}

// An instance is one pinned server this run started.
type instance struct {
	once      sync.Once
	pin       enginetest.Pin
	container *tcmssql.MSSQLServerContainer
	dsn       string
	connector *mssql.Connector // to master
	err       error

	mu      sync.Mutex
	emptied map[string]bool // the top-level tests the server was emptied for
}

var shared = struct {
	sync.Mutex
	byRole map[string]*instance
}{byRole: map[string]*instance{}}

// started returns the server of pin, starting it on the first call.
func started(t testing.TB, pin enginetest.Pin) *instance {
	t.Helper()
	shared.Lock()
	in := shared.byRole[pin.Role]
	if in == nil {
		in = &instance{pin: pin, emptied: map[string]bool{}}
		shared.byRole[pin.Role] = in
	}
	shared.Unlock()
	in.once.Do(func() { in.err = in.start(context.Background()) })
	if in.err != nil {
		t.Fatalf("start SQL Server (%s): %v", pin.Role, in.err)
	}
	if in.pin != pin {
		t.Fatalf("this run started SQL Server at %s, not %s", in.pin, pin)
	}
	return in
}

func (in *instance) start(ctx context.Context) (err error) {
	// testcontainers-go panics from MustExtractDockerHost when the Docker
	// socket is reachable but access is denied, and a test must report that
	// rather than crash the binary.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	c, err := tcmssql.Run(ctx, in.pin.Image, tcmssql.WithAcceptEULA(), tcmssql.WithPassword(password))
	if c != nil {
		in.container = c
	}
	if err != nil {
		return err
	}
	if err := checkDigest(ctx, c, in.pin.Digest); err != nil {
		return err
	}
	if in.dsn, err = c.ConnectionString(ctx); err != nil {
		return err
	}
	if in.connector, err = mssql.NewConnector(in.dsn); err != nil {
		return err
	}
	db := sql.OpenDB(in.connector)
	defer db.Close()
	var version string
	if err := db.QueryRowContext(ctx, "SELECT CAST(SERVERPROPERTY('ProductVersion') AS NVARCHAR(128))").Scan(&version); err != nil {
		return fmt.Errorf("SQL Server does not answer: %w", err)
	}
	if version != in.pin.Tag {
		return fmt.Errorf("the server reports version %q, but pins.tsv pins %s", version, in.pin.Tag)
	}
	return nil
}

// Terminate removes every server this run started, and the ScriptDOM
// harness.
func Terminate() {
	stopHarness()
	shared.Lock()
	defer shared.Unlock()
	for _, in := range shared.byRole {
		if in.container != nil {
			_ = testcontainers.TerminateContainer(in.container)
		}
	}
}

// checkDigest fails unless c runs the image with the pinned repository
// digest. A tag can be pushed again; a digest cannot.
func checkDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the SQL Server container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the SQL Server image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the SQL Server container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
}
