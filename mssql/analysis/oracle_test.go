package analysis

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"os"
	"testing"

	gomssqldb "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
)

// TestMain runs the tests through the tape harness and removes the server a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// oracleConn is t's session on the pinned server, in the test's own database.
func oracleConn(t *testing.T) *sql.Conn {
	t.Helper()
	conn, err := testcontainer.Open(t).Conn(context.Background())
	if err != nil {
		t.Fatalf("connect to SQL Server: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

// oracleExec runs a statement that returns no rows.
func oracleExec(t *testing.T, conn *sql.Conn, stmt string) {
	t.Helper()
	if _, err := conn.ExecContext(context.Background(), stmt); err != nil {
		t.Fatalf("%s: %v", stmt, err)
	}
}

// oracleQuery runs a statement and returns its rows as strings, one slice per
// row, so a test can compare values without caring for their types.
func oracleQuery(t *testing.T, conn *sql.Conn, stmt string) [][]string {
	t.Helper()
	rows, err := conn.QueryContext(context.Background(), stmt)
	if err != nil {
		t.Fatalf("%s: %v", stmt, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("%s: %v", stmt, err)
	}
	var out [][]string
	for rows.Next() {
		cells := make([]any, len(columns))
		for i := range cells {
			cells[i] = new(sql.RawBytes)
		}
		if err := rows.Scan(cells...); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
		row := make([]string, len(cells))
		for i, cell := range cells {
			value := *cell.(*sql.RawBytes)
			if value == nil {
				row[i] = "NULL"
				continue
			}
			row[i] = string(value)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("%s: %v", stmt, err)
	}
	return out
}

// oracleError runs a statement that must fail and returns the server's error
// number.
func oracleError(t *testing.T, conn *sql.Conn, stmt string) int32 {
	t.Helper()
	_, err := conn.ExecContext(context.Background(), stmt)
	var serverError gomssqldb.Error
	if !errors.As(err, &serverError) {
		t.Fatalf("%s: got %v, want SQL Server's error", stmt, err)
	}
	return serverError.Number
}

// oracleBatch runs a script the way the driver does, one request per GO
// delimited batch, and returns each batch's error, nil for a batch that ran.
// It is how a test asks the server what a batch boundary does to state.
func oracleBatch(t *testing.T, conn *sql.Conn, script string) []error {
	t.Helper()
	batcher := NewBatcher(script)
	var errs []error
	for {
		command, err := batcher.Next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				t.Fatalf("split %q: %v", script, err)
			}
			if batch := batcher.Batch(); batch != nil && batch.Text != "" {
				errs = append(errs, oracleRunBatch(t, conn, batch.Text))
			}
			break
		}
		if command == nil {
			continue
		}
		errs = append(errs, oracleRunBatch(t, conn, batcher.Batch().Text))
		batcher.Reset(nil)
	}
	return errs
}

func oracleRunBatch(t *testing.T, conn *sql.Conn, batch string) error {
	t.Helper()
	rows, err := conn.QueryContext(context.Background(), batch)
	if err != nil {
		return err
	}
	defer rows.Close()
	for {
		for rows.Next() {
		}
		if !rows.NextResultSet() {
			break
		}
	}
	return rows.Err()
}

// TestOracleLinkedServerRead pins what SQL Server answers when a statement
// reads another server, which is what the analysis refuses: the connected
// instance's metadata names no such table, and a loopback linked server hands
// back the same rows a local query does.
func TestOracleLinkedServerRead(t *testing.T) {
	conn := oracleConn(t)
	database := testcontainer.Database(t)
	oracleExec(t, conn, "CREATE TABLE dbo.base (a INT, secret NVARCHAR(20))")
	oracleExec(t, conn, "INSERT dbo.base VALUES (1, N'alice'), (2, N'bob')")

	// The server resolves a four-part name, OPENQUERY, and EXEC AT through
	// sys.servers; a server that is not registered is error 7202 for each.
	for _, stmt := range []string{
		"SELECT a FROM srv." + database + ".dbo.base",
		"SELECT * FROM OPENQUERY(srv, 'SELECT a FROM dbo.base')",
		"EXEC ('SELECT 1') AT srv",
	} {
		if got, want := oracleError(t, conn, stmt), int32(7202); got != want {
			t.Errorf("%s: error %d, want %d", stmt, got, want)
		}
	}

	// A linked server to the server itself, so a test needs no second server.
	oracleExec(t, conn, "EXEC sp_addlinkedserver @server = 'loop', @srvproduct = '', @provider = 'MSOLEDBSQL', @datasrc = '127.0.0.1,1433'")
	oracleExec(t, conn, "EXEC sp_addlinkedsrvlogin @rmtsrvname = 'loop', @useself = 'TRUE'")

	// The four-part name reads the table's rows, and no object of the
	// connected instance describes it.
	local := oracleQuery(t, conn, "SELECT a, secret FROM dbo.base ORDER BY a")
	remote := oracleQuery(t, conn, "SELECT a, secret FROM loop."+database+".dbo.base ORDER BY a")
	if len(local) != 2 || len(remote) != len(local) {
		t.Fatalf("four-part name returned %v, local query %v", remote, local)
	}
	for i := range local {
		if remote[i][0] != local[i][0] || remote[i][1] != local[i][1] {
			t.Errorf("four-part name row %d is %v, local query %v", i, remote[i], local[i])
		}
	}
	if remote[1][1] != "bob" {
		t.Errorf("four-part name returned %v, want the secret bob in row 2", remote)
	}
	if got := oracleQuery(t, conn, "SELECT OBJECT_ID('loop."+database+".dbo.base'), OBJECT_ID('dbo.base')"); got[0][0] != "NULL" || got[0][1] == "NULL" {
		t.Errorf("object ids: got %v, want NULL for the remote table and an id for the local one", got)
	}

	// OPENQUERY's column names are the pass-through query's, which the local
	// statement does not spell.
	if got := oracleQuery(t, conn, "SELECT o.bp FROM OPENQUERY(loop, 'SELECT a, a + 1 AS bp FROM "+database+".dbo.base') AS o ORDER BY o.bp"); len(got) != 2 || got[0][0] != "2" || got[1][0] != "3" {
		t.Errorf("OPENQUERY: got %v, want 2 and 3", got)
	}

	// EXEC AT needs rpc out, which the linked server does not have.
	if got, want := oracleError(t, conn, "EXEC ('SELECT 1') AT loop"), int32(7411); got != want {
		t.Errorf("EXEC AT: error %d, want %d", got, want)
	}

	// A local temp table can be filled from the remote table, so what masking
	// would withhold from a direct read can come back through a #temp.
	oracleExec(t, conn, "SELECT a, secret INTO #lt FROM loop."+database+".dbo.base")
	if got := oracleQuery(t, conn, "SELECT a, secret FROM #lt ORDER BY a"); len(got) != 2 || got[1][1] != "bob" {
		t.Errorf("temp table from the four-part name: got %v", got)
	}
}

// TestOracleTempTableRows pins where a #temp table's rows come from: the
// batch's own statements, in every shape SQL Server accepts. A temp table
// outlives GO, so no batch boundary bounds what a later read returns.
func TestOracleTempTableRows(t *testing.T) {
	conn := oracleConn(t)
	oracleExec(t, conn, "CREATE TABLE dbo.base (a INT, secret NVARCHAR(20))")
	oracleExec(t, conn, "INSERT dbo.base VALUES (1, N'alice'), (2, N'bob')")
	oracleExec(t, conn, "CREATE PROCEDURE dbo.p AS SELECT a, secret FROM dbo.base")

	// SELECT ... INTO copies the base table's rows.
	oracleExec(t, conn, "SELECT a, secret INTO #into FROM dbo.base")
	if got := oracleQuery(t, conn, "SELECT a, secret FROM #into ORDER BY a"); len(got) != 2 || got[0][1] != "alice" || got[1][1] != "bob" {
		t.Errorf("#into: got %v, want alice and bob", got)
	}
	// INSERT ... VALUES takes the value an expression reads.
	oracleExec(t, conn, "CREATE TABLE #values (a INT, secret NVARCHAR(20))")
	oracleExec(t, conn, "INSERT #values VALUES ((SELECT MIN(a) FROM dbo.base), (SELECT secret FROM dbo.base WHERE a = 1))")
	if got := oracleQuery(t, conn, "SELECT a, secret FROM #values"); len(got) != 1 || got[0][1] != "alice" {
		t.Errorf("#values: got %v, want alice", got)
	}
	// INSERT ... SELECT copies a projection.
	oracleExec(t, conn, "CREATE TABLE #select (a INT, secret NVARCHAR(20))")
	oracleExec(t, conn, "INSERT #select SELECT a, secret FROM dbo.base")
	if got := oracleQuery(t, conn, "SELECT a, secret FROM #select ORDER BY a"); len(got) != 2 || got[1][1] != "bob" {
		t.Errorf("#select: got %v, want bob in row 2", got)
	}
	// UPDATE ... SET ... FROM replaces a value with another table's column.
	oracleExec(t, conn, "CREATE TABLE #update (a INT, secret NVARCHAR(20))")
	oracleExec(t, conn, "INSERT #update (a) VALUES (1)")
	oracleExec(t, conn, "UPDATE #update SET secret = b.secret FROM dbo.base AS b WHERE b.a = #update.a")
	if got := oracleQuery(t, conn, "SELECT a, secret FROM #update"); len(got) != 1 || got[0][1] != "alice" {
		t.Errorf("#update: got %v, want alice", got)
	}
	// MERGE ... USING reads a source query.
	oracleExec(t, conn, "CREATE TABLE #merge (a INT, secret NVARCHAR(20))")
	oracleExec(t, conn, "MERGE #merge AS tgt USING (SELECT a, secret FROM dbo.base) AS src ON tgt.a = src.a WHEN NOT MATCHED THEN INSERT (a, secret) VALUES (src.a, src.secret);")
	if got := oracleQuery(t, conn, "SELECT a, secret FROM #merge ORDER BY a"); len(got) != 2 || got[1][1] != "bob" {
		t.Errorf("#merge: got %v, want bob in row 2", got)
	}
	// INSERT ... EXEC takes the rows a procedure returns, which no column of
	// the statement names.
	oracleExec(t, conn, "CREATE TABLE #exec (a INT, secret NVARCHAR(20))")
	oracleExec(t, conn, "INSERT #exec EXEC dbo.p")
	if got := oracleQuery(t, conn, "SELECT a, secret FROM #exec ORDER BY a"); len(got) != 2 || got[0][1] != "alice" {
		t.Errorf("#exec: got %v, want alice and bob", got)
	}

	// A temp table created in one batch is read in the next.
	for i, err := range oracleBatch(t, conn, "SELECT a, secret INTO #across FROM dbo.base\nGO\nSELECT COUNT(*) AS n FROM #across") {
		if err != nil {
			t.Fatalf("batch %d: %v", i, err)
		}
	}
	if got := oracleQuery(t, conn, "SELECT a, secret FROM #across ORDER BY a"); len(got) != 2 || got[0][1] != "alice" {
		t.Errorf("#across after GO: got %v, want alice and bob", got)
	}
}

// TestOracleVariableRows pins where a variable's value comes from: a column of
// the batch, when a statement of the batch assigned it. A batch boundary ends
// the variable's scope, unlike a temp table's.
func TestOracleVariableRows(t *testing.T) {
	conn := oracleConn(t)
	oracleExec(t, conn, "CREATE TABLE dbo.base (a INT, secret NVARCHAR(20))")
	oracleExec(t, conn, "INSERT dbo.base VALUES (1, N'alice'), (2, N'bob')")

	// DECLARE, SET, and SELECT @v = ... each take a base column's value into a
	// variable that a later statement of the batch reads.
	for _, script := range []string{
		"DECLARE @v NVARCHAR(20) = (SELECT secret FROM dbo.base WHERE a = 1); SELECT @v AS v",
		"DECLARE @v NVARCHAR(20); SET @v = (SELECT secret FROM dbo.base WHERE a = 2); SELECT @v AS v",
		"DECLARE @v NVARCHAR(20); SELECT @v = secret FROM dbo.base WHERE a = 1; SELECT @v AS v",
	} {
		if got := oracleQuery(t, conn, script); len(got) != 1 || (got[0][0] != "alice" && got[0][0] != "bob") {
			t.Errorf("%s: got %v, want the base table's secret", script, got)
		}
	}

	// A SELECT that assigns a variable returns no result set.
	rows, err := conn.QueryContext(context.Background(), "DECLARE @v INT; SELECT @v = a FROM dbo.base")
	if err != nil {
		t.Fatalf("SELECT @v = a: %v", err)
	}
	if columns, err := rows.Columns(); err != nil || len(columns) != 0 {
		t.Errorf("SELECT @v = a: columns %v, %v, want none", columns, err)
	}
	rows.Close()

	// A SELECT that assigns a variable and retrieves a column is error 141.
	if got, want := oracleError(t, conn, "DECLARE @v INT; SELECT @v = a, secret FROM dbo.base"), int32(141); got != want {
		t.Errorf("mixed assignment and retrieval: error %d, want %d", got, want)
	}

	// The next batch no longer has the variable, even on the same session.
	errs := oracleBatch(t, conn, "DECLARE @v INT = 1\nGO\nSELECT @v AS v")
	if len(errs) != 2 || errs[0] != nil {
		t.Fatalf("batches: got %v", errs)
	}
	var serverError gomssqldb.Error
	if !errors.As(errs[1], &serverError) || serverError.Number != 137 {
		t.Errorf("SELECT @v in the next batch: got %v, want error 137", errs[1])
	}
}

// TestOracleBackupTop verifies that the SELECT ... INTO generated for an
// UPDATE TOP statement executes on SQL Server and copies at most the TOP rows.
func TestOracleBackupTop(t *testing.T) {
	conn := oracleConn(t)
	database := testcontainer.Database(t)
	oracleExec(t, conn, "CREATE TABLE dbo.test (c1 INT, c2 INT)")
	oracleExec(t, conn, "INSERT dbo.test VALUES (1, 2), (2, 2), (3, 2), (4, 2)")

	backups, err := TransformDMLToSelect(
		parseScript(t, "UPDATE TOP (3) test SET c1 = 9 WHERE c2 = 2;"),
		database, database, "rollback", false,
	)
	if err != nil {
		t.Fatalf("TransformDMLToSelect: %v", err)
	}
	if len(backups) != 1 {
		t.Fatalf("got %d backups, want 1", len(backups))
	}
	oracleExec(t, conn, backups[0].Statement)
	if got := oracleQuery(t, conn, "SELECT COUNT(*) FROM dbo.["+backups[0].TargetTableName+"]"); len(got) != 1 || got[0][0] != "3" {
		t.Errorf("backup row count: got %v, want 3", got)
	}
}
