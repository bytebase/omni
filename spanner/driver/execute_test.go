package driver

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

// executeLogger records what Execute reported, as one line per event.
type executeLogger struct {
	events []string
}

func (l *executeLogger) CommandExecute(command driver.Command) {
	l.events = append(l.events, "command:"+command.Text)
}

func (l *executeLogger) CommandResponse(affectedRows int64, _ []int64, err error) {
	l.events = append(l.events, "response")
}

func (l *executeLogger) TransactionControl(op driver.TransactionOp, err error) {
	if err != nil {
		l.events = append(l.events, string(op)+":failed")
		return
	}
	l.events = append(l.events, string(op))
}

func (l *executeLogger) Retry(error, int, int) {
	l.events = append(l.events, "retry")
}

// openScratch opens a driver on a fixture database of its own in instance
// instance, so a test that writes cannot change what another test recorded.
func openScratch(t *testing.T, instance, database string) *Driver {
	t.Helper()
	testcontainer.Create(t, instance, database, fixtureDDL, fixtureRows)
	return openDriver(t, instance, database)
}

// TestExecuteTransaction runs a script's DML in one read-write transaction and
// reports its steps.
func TestExecuteTransaction(t *testing.T) {
	a := require.New(t)
	d := openScratch(t, "execute", "transaction")
	defer d.Close()

	log := &executeLogger{}
	count, err := d.Execute(context.Background(),
		`UPDATE Singers SET Rating = 5.0 WHERE Id = 1; UPDATE Singers SET Rating = 4.0 WHERE Id = 2`,
		driver.ExecuteOptions{Log: log})
	a.NoError(err)
	a.Equal(int64(2), count)
	a.Equal([]string{string(driver.Begin), string(driver.Commit)}, log.events)
}

// TestExecuteAutoCommit runs each statement as a partitioned update, whose row
// count is Spanner's lower bound, and reports no transaction it did not open.
func TestExecuteAutoCommit(t *testing.T) {
	a := require.New(t)
	d := openScratch(t, "execute", "autocommit")
	defer d.Close()

	log := &executeLogger{}
	count, err := d.Execute(context.Background(),
		`UPDATE Singers SET Active = TRUE WHERE Id = 1`,
		driver.ExecuteOptions{Transaction: driver.AutoCommit, Log: log})
	a.NoError(err)
	a.Equal(int64(1), count)
	a.Empty(log.events)
}

// TestExecuteDDL applies a script of DDL statements through the admin API.
func TestExecuteDDL(t *testing.T) {
	a := require.New(t)
	d := openScratch(t, "execute", "ddl")
	defer d.Close()

	count, err := d.Execute(context.Background(),
		`CREATE TABLE Written (Id INT64 NOT NULL, Name STRING(64)) PRIMARY KEY (Id);
		 DROP TABLE Written`, driver.ExecuteOptions{})
	a.NoError(err)
	a.Equal(int64(0), count)
}

// TestExecuteRollback rolls back every statement of a script that failed
// partway, and reports the transaction it abandoned.
func TestExecuteRollback(t *testing.T) {
	a := require.New(t)
	d := openScratch(t, "execute", "rollback")
	defer d.Close()

	log := &executeLogger{}
	count, err := d.Execute(context.Background(),
		`UPDATE Singers SET Rating = 1.0 WHERE Id = 1; UPDATE Missing SET a = 1`,
		driver.ExecuteOptions{Log: log})
	a.Error(err)
	a.Equal(int64(0), count)
	a.Equal([]string{string(driver.Begin), string(driver.Rollback) + ":failed"}, log.events)

	// The first statement's write is not there: the transaction rolled back.
	results, err := sessionOf(t, d).Query(context.Background(), `SELECT Rating FROM Singers WHERE Id = 1`, driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Equal(driver.Value{Kind: driver.Float64, Float64: 4.5}, results[0].Rows[0][0])
}

// TestExecuteMixedDDLAndDML pins what a script that holds both kinds does: the
// whole script goes to the admin API, which refuses the DML statement. Bytebase's
// driver routed a script with any DDL the same way, and the move keeps it; the
// fix is a change of its own.
func TestExecuteMixedDDLAndDML(t *testing.T) {
	a := require.New(t)
	d := openScratch(t, "execute", "mixed")
	defer d.Close()

	_, err := d.Execute(context.Background(),
		`CREATE TABLE Mixed (Id INT64 NOT NULL) PRIMARY KEY (Id); UPDATE Singers SET Rating = 3.0 WHERE Id = 1`,
		driver.ExecuteOptions{})
	a.Error(err)
	a.Contains(err.Error(), "InvalidArgument")
}

// TestExecuteCreateDatabase creates a database from the script's first statement
// and applies the rest as its extra DDL.
func TestExecuteCreateDatabase(t *testing.T) {
	a := require.New(t)
	testcontainer.Create(t, "createdb", "createdb", nil, nil)
	d := openDriver(t, "createdb", "createdb")
	defer d.Close()

	count, err := d.Execute(context.Background(),
		`CREATE DATABASE `+"`created`"+`; CREATE TABLE Written (Id INT64 NOT NULL) PRIMARY KEY (Id)`,
		driver.ExecuteOptions{CreateDatabase: true})
	a.NoError(err)
	a.Equal(int64(0), count)

	created := openDriver(t, "createdb", "created")
	defer created.Close()
	var dumped strings.Builder
	a.NoError(created.Dump(context.Background(), &dumped, nil))
	a.Contains(dumped.String(), "CREATE TABLE Written")
}

// TestExecuteCreateDatabaseWithAHeaderComment creates a database from a script
// that begins with a comment: the comment is part of the statement's text now,
// and asking the lexer for the first two keywords skips it, as Bytebase's
// comment-stripping splitter did.
func TestExecuteCreateDatabaseWithAHeaderComment(t *testing.T) {
	a := require.New(t)
	testcontainer.Create(t, "createdb", "createdb", nil, nil)
	d := openDriver(t, "createdb", "createdb")
	defer d.Close()

	count, err := d.Execute(context.Background(),
		"-- provision the database\n/* and its schema */ CREATE DATABASE `createdheader`; CREATE TABLE Written (Id INT64 NOT NULL) PRIMARY KEY (Id)",
		driver.ExecuteOptions{CreateDatabase: true})
	a.NoError(err)
	a.Equal(int64(0), count)

	created := openDriver(t, "createdb", "createdheader")
	defer created.Close()
	var dumped strings.Builder
	a.NoError(created.Dump(context.Background(), &dumped, nil))
	a.Contains(dumped.String(), "CREATE TABLE Written")
}

// TestExecuteCreateDatabaseRefusesOtherStatements refuses a creation script whose
// first statement is not CREATE DATABASE, which is what Bytebase's check did.
func TestExecuteCreateDatabaseRefusesOtherStatements(t *testing.T) {
	a := require.New(t)
	d := openScratch(t, "execute", "createcheck")
	defer d.Close()

	_, err := d.Execute(context.Background(), `DROP TABLE Singers`, driver.ExecuteOptions{CreateDatabase: true})
	a.EqualError(err, `expect the first entry of the sanitized SQLs to start with 'CREATE DATABASE', sql DROP TABLE Singers`)

	_, err = d.Execute(context.Background(), ``, driver.ExecuteOptions{CreateDatabase: true})
	a.EqualError(err, "expect sanitized SQLs to have at least one entry, original statement: ")
}

// TestExecutePostgres runs a PostgreSQL script's DML in one read-write
// transaction and its DDL through the admin API, in PostgreSQL's SQL.
func TestExecutePostgres(t *testing.T) {
	a := require.New(t)
	d := openPostgresFixture(t, "pgexecute")
	defer d.Close()

	log := &executeLogger{}
	count, err := d.Execute(context.Background(),
		`UPDATE Singers SET rating = 5.0 WHERE id = 1; UPDATE Singers SET rating = 4.0 WHERE id = 2`,
		driver.ExecuteOptions{Log: log})
	a.NoError(err)
	a.Equal(int64(2), count)
	a.Equal([]string{string(driver.Begin), string(driver.Commit)}, log.events)

	count, err = d.Execute(context.Background(),
		`CREATE TABLE written (id bigint NOT NULL, name text, PRIMARY KEY (id)); DROP TABLE written`,
		driver.ExecuteOptions{})
	a.NoError(err)
	a.Equal(int64(0), count)
}

// TestExecutePostgresAutoCommit runs each PostgreSQL statement as a partitioned
// update, whose row count is Spanner's lower bound.
func TestExecutePostgresAutoCommit(t *testing.T) {
	a := require.New(t)
	d := openPostgresFixture(t, "pgexecute")
	defer d.Close()

	count, err := d.Execute(context.Background(),
		`UPDATE Singers SET active = TRUE WHERE id = 1`,
		driver.ExecuteOptions{Transaction: driver.AutoCommit})
	a.NoError(err)
	a.Equal(int64(1), count)
}

// TestExecuteCreateDatabasePostgres creates a database from a PostgreSQL
// creation script: the double-quoted name is the dialect signal the admin API
// requires, and the rest of the script is applied in that dialect too.
func TestExecuteCreateDatabasePostgres(t *testing.T) {
	a := require.New(t)
	testcontainer.CreatePostgres(t, "pgcreatedb", "pgcreatedb", nil, nil)
	d := openDriverWithDialect(t, "pgcreatedb", "pgcreatedb", postgresSQL)
	defer d.Close()

	count, err := d.Execute(context.Background(),
		`CREATE DATABASE "pgcreated"; CREATE TABLE written (id bigint NOT NULL, PRIMARY KEY (id))`,
		driver.ExecuteOptions{CreateDatabase: true})
	a.NoError(err)
	a.Equal(int64(0), count)

	created := openDriverWithDialect(t, "pgcreatedb", "pgcreated", postgresSQL)
	defer created.Close()
	// The created database answers the PostgreSQL catalog, which is how the test
	// reads the dialect the server created it with.
	metadata, err := created.SyncDatabase(context.Background())
	a.NoError(err)
	a.Len(metadata.Schemas, 1)
	a.Equal("public", metadata.Schemas[0].Name)
	a.Len(metadata.Schemas[0].Tables, 1)
	a.Equal("written", metadata.Schemas[0].Tables[0].Name)
}

// TestExecuteCreateDatabaseWithAHashHeader creates a database from a GoogleSQL
// script that begins with a `#` line comment: the dialect is read with
// GoogleSQL's lexer first, which skips that comment rather than taking a quoted
// word in it for the database name.
func TestExecuteCreateDatabaseWithAHashHeader(t *testing.T) {
	a := require.New(t)
	testcontainer.Create(t, "createdb", "createdb", nil, nil)
	d := openDriver(t, "createdb", "createdb")
	defer d.Close()

	count, err := d.Execute(context.Background(),
		"# provision the \"warehouse\"\nCREATE DATABASE `createdhash`; CREATE TABLE Written (Id INT64 NOT NULL) PRIMARY KEY (Id)",
		driver.ExecuteOptions{CreateDatabase: true})
	a.NoError(err)
	a.Equal(int64(0), count)

	created := openDriver(t, "createdb", "createdhash")
	defer created.Close()
	var dumped strings.Builder
	a.NoError(created.Dump(context.Background(), &dumped, nil))
	a.Contains(dumped.String(), "CREATE TABLE Written")
}

// TestExecuteCreateDatabaseRefreshesTheDialect creates the database the driver is
// open on with a PostgreSQL creation script. Open found no database to read a
// dialect from and left GoogleSQL, so the creation is what the query after it
// splits by: a dollar-quoted string with a semicolon in it stays one statement.
func TestExecuteCreateDatabaseRefreshesTheDialect(t *testing.T) {
	a := require.New(t)
	// The instance exists; the database the driver is opened on does not.
	testcontainer.CreatePostgres(t, "pgfresh", "seed", nil, nil)
	d := openDriverWithDialect(t, "pgfresh", "pgfresh", googleSQL)
	defer d.Close()

	count, err := d.Execute(context.Background(),
		`CREATE DATABASE "pgfresh"`, driver.ExecuteOptions{CreateDatabase: true})
	a.NoError(err)
	a.Equal(int64(0), count)

	results, err := sessionOf(t, d).Query(context.Background(),
		`SELECT $$a;b$$ AS x`, driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Len(results, 1)
	a.NoError(results[0].Err)
	a.Equal([]driver.Value{{Kind: driver.String, String: "a;b"}}, results[0].Rows[0])
}
