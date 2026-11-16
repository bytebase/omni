package driver

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	omnidriver "github.com/bytebase/omni/driver"
	oracleparser "github.com/bytebase/omni/oracle/parser"
)

type stepLogger struct {
	steps []string
}

func (l *stepLogger) CommandExecute(command omnidriver.Command) {
	l.steps = append(l.steps, fmt.Sprintf("exec[%d:%d]", command.ByteStart, command.ByteEnd))
}

func (l *stepLogger) CommandResponse(affected int64, perResult []int64, err error) {
	if err != nil {
		l.steps = append(l.steps, "resp:err")
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("resp:%d:%v", affected, perResult))
}

func (l *stepLogger) TransactionControl(op omnidriver.TransactionOp, err error) {
	if err != nil {
		l.steps = append(l.steps, fmt.Sprintf("tx:%s:err", op))
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("tx:%s", op))
}

func (l *stepLogger) Retry(_ error, attempt, max int) {
	l.steps = append(l.steps, fmt.Sprintf("retry:%d/%d", attempt, max))
}

// querySetup fills the schema with the rows the query tests read.
const querySetup = `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_APP;
CREATE TABLE OMNI_APP.Q (
	ID NUMBER(4) NOT NULL,
	LABEL VARCHAR2(20),
	AMOUNT NUMBER(8, 2),
	FLAG NUMBER(1),
	CREATED DATE,
	AT TIMESTAMP(3),
	AT_ZONE TIMESTAMP WITH TIME ZONE,
	DOC CLOB,
	BIN BLOB,
	CONSTRAINT Q_PK PRIMARY KEY (ID)
);
INSERT INTO OMNI_APP.Q (ID, LABEL, AMOUNT, FLAG, CREATED, AT, AT_ZONE, DOC, BIN)
	VALUES (1, 'a', 1.5, 1, DATE '2026-03-04', TIMESTAMP '2026-03-04 05:06:07.891', TIMESTAMP '2026-03-04 05:06:07.891 +02:00', 'doc', HEXTORAW('01'));
INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (2, 'b');
INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (3, NULL);
`

func TestQuery(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	results, err := s.Query(ctx, "SELECT ID, LABEL, AMOUNT, FLAG, CREATED, AT, DOC FROM OMNI_APP.Q ORDER BY ID", omnidriver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	got := results[0]
	if got.Err != nil {
		t.Fatalf("query failed: %v", got.Err)
	}
	wantColumns := []omnidriver.Column{
		{Name: "ID", TypeName: "NUMBER"},
		{Name: "LABEL", TypeName: "NCHAR"},
		{Name: "AMOUNT", TypeName: "NUMBER"},
		{Name: "FLAG", TypeName: "NUMBER"},
		{Name: "CREATED", TypeName: "DATE"},
		{Name: "AT", TypeName: "TIMESTAMPDTY"},
		{Name: "DOC", TypeName: "LONGVARCHAR"},
	}
	if diff := cmp.Diff(wantColumns, got.Columns); diff != "" {
		t.Errorf("columns (-want +got):\n%s", diff)
	}
	if len(got.Rows) != 3 {
		t.Fatalf("rows = %d, want 3", len(got.Rows))
	}
	// A NUMBER arrives as text, which the driver keeps.
	if got.Rows[0][0].Kind != omnidriver.String || got.Rows[0][0].String != "1" {
		t.Errorf("ID = %+v, want the text 1", got.Rows[0][0])
	}
	if got.Rows[2][1].Kind != omnidriver.Null {
		t.Errorf("a NULL label = %+v", got.Rows[2][1])
	}
	if got.Rows[1][2].Kind != omnidriver.Null {
		t.Errorf("an unset amount = %+v", got.Rows[1][2])
	}
	// A DATE and a TIMESTAMP keep their wall clock; only TIMESTAMP WITH TIME
	// ZONE carries a zone.
	if got.Rows[0][4].Kind != omnidriver.Timestamp || got.Rows[0][4].Time.Format("2006-01-02 15:04:05") != "2026-03-04 00:00:00" {
		t.Errorf("DATE = %+v", got.Rows[0][4])
	}
	if got.Rows[0][5].Kind != omnidriver.Timestamp || got.Rows[0][5].Time.Format("2006-01-02 15:04:05.000") != "2026-03-04 05:06:07.891" {
		t.Errorf("TIMESTAMP = %+v", got.Rows[0][5])
	}
	// go-ora reports no fractional-second size for a TIMESTAMP.
	if got.Rows[0][5].Precision != 0 {
		t.Errorf("TIMESTAMP precision = %d, want the client's 0", got.Rows[0][5].Precision)
	}
}

func TestQueryTimeZone(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	results, err := s.Query(ctx, "SELECT AT_ZONE, BIN FROM OMNI_APP.Q WHERE ID = 1", omnidriver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	row := results[0].Rows[0]
	if row[0].Kind != omnidriver.TimestampTZ {
		t.Fatalf("TIMESTAMP WITH TIME ZONE = %+v", row[0])
	}
	// The value is the instant it names, rendered in the session's zone.
	if got := row[0].Time.UTC().Format("2006-01-02 15:04:05"); got != "2026-03-04 03:06:07" {
		t.Errorf("TIMESTAMP WITH TIME ZONE = %s, want the 03:06:07 UTC instant", got)
	}
	// go-ora names BLOB LONGVARBINARY, which the driver reads as text, as
	// Bytebase does.
	if row[1].Kind != omnidriver.String || len(row[1].String) == 0 {
		t.Errorf("BLOB = %+v, want the client's text", row[1])
	}
}

// TestQuerySchema pins that QueryOptions.Schema selects the schema an
// unqualified statement resolves in.
func TestQuerySchema(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_APP;
CREATE TABLE OMNI_APP.S (ID NUMBER);
INSERT INTO OMNI_APP.S (ID) VALUES (1);
`)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	// The session runs as SYSTEM, so an unqualified name reaches no table of
	// OMNI_APP's without the option.
	results, err := s.Query(ctx, "SELECT ID FROM S", omnidriver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Err == nil {
		t.Fatalf("an unqualified name resolved without QueryOptions.Schema")
	}

	results, err = s.Query(ctx, "SELECT ID FROM S", omnidriver.QueryOptions{Schema: testSchema})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Err != nil {
		t.Fatalf("QueryOptions.Schema did not select the schema: %v", results[0].Err)
	}
	if len(results[0].Rows) != 1 || results[0].Rows[0][0].String != "1" {
		t.Errorf("rows = %+v, want the one row", results[0].Rows)
	}
}

// TestVersionFromServer pins that the driver asks the server for its release and
// reads the configured version only when the server does not answer.
func TestVersionFromServer(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	d.engineVersion = "11.1.0.0.0"

	version, err := d.version(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if version.First != 23 || version.Second != 26 {
		t.Errorf("version = %+v, want the server's 23.26", version)
	}
}

// TestQuerySelectDeclaringPLSQL pins that a SELECT whose WITH clause declares
// PL/SQL runs as a query and returns its rows, though the read-only validator
// refuses that form.
func TestQuerySelectDeclaringPLSQL(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	results, err := s.Query(ctx, "WITH FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; SELECT f FROM dual", omnidriver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Err != nil {
		t.Fatalf("the query failed: %v", results[0].Err)
	}
	if len(results[0].Rows) != 1 || len(results[0].Rows[0]) != 1 || results[0].Rows[0][0].String != "1" {
		t.Errorf("rows = %+v, want the one row the function returns", results[0].Rows)
	}
}

// TestQueryExplainStopsAtTheFailingStatement pins that a failing EXPLAIN keeps
// the plans of the statements before it, and reports its error on its own
// result.
func TestQueryExplainStopsAtTheFailingStatement(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	results, err := s.Query(ctx, "SELECT ID FROM OMNI_APP.Q; SELECT NOSUCH FROM OMNI_APP.Q; SELECT ID FROM OMNI_APP.Q", omnidriver.QueryOptions{Explain: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2: the statement after the failure does not run", len(results))
	}
	if results[0].Err != nil || results[0].Plan == nil {
		t.Errorf("first statement = %+v, want its plan", results[0])
	}
	if results[1].Err == nil || !strings.Contains(results[1].Err.Error(), "NOSUCH") {
		t.Errorf("second statement error = %v, want the plan failure", results[1].Err)
	}
}

// TestQueryReportsAffectedRows pins that a write in a script reports the rows it
// changed rather than returning data.
func TestQueryReportsAffectedRows(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	results, err := s.Query(ctx, "UPDATE OMNI_APP.Q SET LABEL = 'x' WHERE ID <= 2", omnidriver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	if results[0].AffectedRows == nil || *results[0].AffectedRows != 2 {
		t.Errorf("affected rows = %v, want 2", results[0].AffectedRows)
	}
	if len(results[0].Rows) != 0 {
		t.Errorf("rows = %d, want none", len(results[0].Rows))
	}
}

// TestQueryStopsAtTheFailingStatement pins that a script stops at its first
// failure and the failure is reported on the result of that statement.
func TestQueryStopsAtTheFailingStatement(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	results, err := s.Query(ctx, "SELECT ID FROM OMNI_APP.Q WHERE ID = 1; SELECT NOSUCH FROM OMNI_APP.Q; SELECT ID FROM OMNI_APP.Q WHERE ID = 2", omnidriver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2: the statement after the failure does not run", len(results))
	}
	if results[0].Err != nil {
		t.Errorf("first statement failed: %v", results[0].Err)
	}
	if results[1].Err == nil || !strings.Contains(results[1].Err.Error(), "NOSUCH") {
		t.Errorf("second statement error = %v, want the invalid identifier", results[1].Err)
	}
}

// TestQueryLimit pins that QueryOptions.Limit bounds a statement it can rewrite
// and sets a result Truncated past MaxResultBytes.
func TestQueryLimit(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	results, err := s.Query(ctx, "SELECT ID FROM OMNI_APP.Q ORDER BY ID", omnidriver.QueryOptions{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(results[0].Rows); got != 2 {
		t.Errorf("rows = %d, want 2", got)
	}
	if !strings.Contains(strings.ToUpper(results[0].Statement), "ROWNUM") && !strings.Contains(strings.ToUpper(results[0].Statement), "FETCH NEXT") {
		t.Errorf("statement = %q, want a limit rewrite", results[0].Statement)
	}

	// One row exceeds a byte, so reading stops after it and says so.
	results, err = s.Query(ctx, "SELECT ID, LABEL FROM OMNI_APP.Q ORDER BY ID", omnidriver.QueryOptions{MaxResultBytes: 1})
	if err != nil {
		t.Fatal(err)
	}
	if !results[0].Truncated {
		t.Errorf("reading stopped at MaxResultBytes but the result does not say so")
	}
	if got := len(results[0].Rows); got != 1 {
		t.Errorf("rows = %d, want the row that passed the cap", got)
	}
}

// TestQueryExplain pins that Explain returns the plan instead of the rows, in
// the requested format.
func TestQueryExplain(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	// The default is the plan as Oracle states it, its printed rows.
	results, err := s.Query(ctx, "SELECT ID FROM OMNI_APP.Q", omnidriver.QueryOptions{Explain: true})
	if err != nil {
		t.Fatal(err)
	}
	res := results[0]
	if res.Err != nil {
		t.Fatalf("explain failed: %v", res.Err)
	}
	if res.Plan == nil || res.Plan.Format != "TEXT" {
		t.Fatalf("plan = %+v, want TEXT", res.Plan)
	}
	if len(res.Rows) == 0 || res.Columns[0].Name != "Operation" {
		t.Errorf("text plan = %+v, want the printed rows", res.Columns)
	}

	results, err = s.Query(ctx, "SELECT ID FROM OMNI_APP.Q", omnidriver.QueryOptions{Explain: true, ExplainFormat: "JSON"})
	if err != nil {
		t.Fatal(err)
	}
	res = results[0]
	if res.Err != nil {
		t.Fatalf("explain failed: %v", res.Err)
	}
	if res.Plan == nil || res.Plan.Format != "JSON" {
		t.Fatalf("plan = %+v, want JSON", res.Plan)
	}
	if len(res.Rows) != 1 || len(res.Rows[0]) != 1 || res.Rows[0][0].Kind != omnidriver.String {
		t.Fatalf("explain rows = %+v, want one plan document", res.Rows)
	}
	if !strings.Contains(res.Rows[0][0].String, `"operation"`) {
		t.Errorf("plan = %s, want the operators", res.Rows[0][0].String)
	}
}

// commandSteps returns the steps Execute logs for script's commands, paired with
// the responses the test expects for them.
func commandSteps(t *testing.T, script string, responses ...string) []string {
	t.Helper()
	var steps []string
	i := 0
	for _, statement := range oracleparser.Split(script) {
		if statement.Empty() {
			continue
		}
		if i >= len(responses) {
			t.Fatalf("script has more commands than the test's %d responses", len(responses))
		}
		steps = append(steps, fmt.Sprintf("exec[%d:%d]", statement.ByteStart, statement.ByteEnd), responses[i])
		i++
		// The responses end where the script stops: a failing statement's
		// successors run not at all.
		if i == len(responses) {
			break
		}
	}
	if i != len(responses) {
		t.Fatalf("script has %d commands, the test expected %d", i, len(responses))
	}
	return steps
}

func TestExecute(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)

	script := `
INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (4, 'd');
UPDATE OMNI_APP.Q SET LABEL = 'e' WHERE ID = 4;
`
	log := &stepLogger{}
	affected, err := d.Execute(ctx, script, omnidriver.ExecuteOptions{Transaction: omnidriver.Transactional, Log: log})
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if affected != 2 {
		t.Errorf("affected = %d, want 2", affected)
	}
	want := append([]string{"tx:BEGIN"}, commandSteps(t, script, "resp:1:[1]", "resp:1:[1]")...)
	want = append(want, "tx:COMMIT")
	if diff := cmp.Diff(want, log.steps); diff != "" {
		t.Errorf("steps (-want +got):\n%s", diff)
	}

	// The transaction rolled nothing back: both statements are readable.
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()
	results, err := s.Query(ctx, "SELECT LABEL FROM OMNI_APP.Q WHERE ID = 4", omnidriver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got := results[0].Rows[0][0]; got.Kind != omnidriver.String || got.String != "e" {
		t.Errorf("label = %+v, want e", got)
	}
}

func TestExecuteFailureRollsBack(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)

	script := `
INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (5, 'f');
INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (5, 'g');
`
	log := &stepLogger{}
	_, err := d.Execute(ctx, script, omnidriver.ExecuteOptions{Transaction: omnidriver.Transactional, Log: log})
	if err == nil {
		t.Fatal("the primary key violation succeeded")
	}
	want := append([]string{"tx:BEGIN"}, commandSteps(t, script, "resp:1:[1]", "resp:err")...)
	want = append(want, "tx:ROLLBACK")
	if diff := cmp.Diff(want, log.steps); diff != "" {
		t.Errorf("steps (-want +got):\n%s", diff)
	}

	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()
	results, err := s.Query(ctx, "SELECT COUNT(*) FROM OMNI_APP.Q WHERE ID = 5", omnidriver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got := results[0].Rows[0][0]; got.Kind != omnidriver.String || got.String != "0" {
		t.Errorf("committed rows = %+v, want the rollback to leave none", got)
	}
}

// TestExecuteAutoCommitStopsAtTheFailingStatement pins that a script outside a
// transaction keeps what it ran before the failure.
func TestExecuteAutoCommitStopsAtTheFailingStatement(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, querySetup)

	script := `
INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (6, 'h');
INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (6, 'i');
INSERT INTO OMNI_APP.Q (ID, LABEL) VALUES (7, 'j');
`
	log := &stepLogger{}
	_, err := d.Execute(ctx, script, omnidriver.ExecuteOptions{Transaction: omnidriver.AutoCommit, Log: log})
	if err == nil {
		t.Fatal("the primary key violation succeeded")
	}
	want := commandSteps(t, script, "resp:1:[1]", "resp:err")
	if diff := cmp.Diff(want, log.steps); diff != "" {
		t.Errorf("steps (-want +got):\n%s", diff)
	}

	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()
	results, err := s.Query(ctx, "SELECT ID FROM OMNI_APP.Q WHERE ID >= 6 ORDER BY ID", omnidriver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(results[0].Rows); got != 1 {
		t.Errorf("rows = %d, want the one the auto-commit kept", got)
	}
}

// TestOwnsPrivateDatabaseLink pins the probe a caller asks before it runs a
// statement that could reach another database over a private link.
func TestOwnsPrivateDatabaseLink(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
`)

	owns, err := d.OwnsPrivateDatabaseLink(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if owns {
		t.Fatal("the account owns a private database link it did not make")
	}

	execCorpus(t, d.db, "CREATE DATABASE LINK omni_link CONNECT TO omni IDENTIFIED BY omni USING 'nosuch.example.com:1521/XE'")
	owns, err = d.OwnsPrivateDatabaseLink(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !owns {
		t.Error("the account's own private database link went unreported")
	}
	// The link is the instance's, not the schema's, so a later test would read it.
	execCorpus(t, d.db, "DROP DATABASE LINK omni_link")
}

// TestQueryLocalTimeZone pins that a TIMESTAMP WITH LOCAL TIME ZONE is printed in
// the session's zone, as SQL*Plus prints it, including a zone the script changes
// in the middle.
func TestQueryLocalTimeZone(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_APP;
CREATE TABLE OMNI_APP.LTZ (ID NUMBER, AT TIMESTAMP WITH LOCAL TIME ZONE);
ALTER SESSION SET TIME_ZONE = '+02:00';
INSERT INTO OMNI_APP.LTZ (ID, AT) VALUES (1, TIMESTAMP '2026-03-04 05:06:07');
`)
	s, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	for _, test := range []struct {
		name string
		zone string
		want string
	}{
		// The row was stored as 03:06:07 UTC.
		{name: "offset", zone: "+02:00", want: "2026-03-04 05:06:07"},
		{name: "utc", zone: "+00:00", want: "2026-03-04 03:06:07"},
		{name: "region", zone: "Europe/Berlin", want: "2026-03-04 04:06:07"},
	} {
		t.Run(test.name, func(t *testing.T) {
			results, err := s.Query(ctx, fmt.Sprintf("ALTER SESSION SET TIME_ZONE = '%s'; SELECT AT FROM OMNI_APP.LTZ", test.zone), omnidriver.QueryOptions{})
			if err != nil {
				t.Fatal(err)
			}
			row := results[len(results)-1].Rows[0][0]
			if row.Kind != omnidriver.Timestamp {
				t.Fatalf("value = %+v, want a wall clock", row)
			}
			if got := row.Time.Format(time.DateTime); got != test.want {
				t.Errorf("value = %s, want %s", got, test.want)
			}
		})
	}
}

// TestExecuteCreateDatabase pins that Oracle refuses a database-creation step:
// a schema is a user, which is not what a change script creates.
func TestExecuteCreateDatabase(t *testing.T) {
	d := openTestDriver(t)
	if _, err := d.Execute(context.Background(), "SELECT 1 FROM DUAL", omnidriver.ExecuteOptions{CreateDatabase: true}); err == nil {
		t.Fatal("a create-database step succeeded")
	}
}

// TestCountAffectedRows pins that the estimate comes from the plan of the
// statement, planned with dynamic sampling.
func TestCountAffectedRows(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE USER OMNI_APP NO AUTHENTICATION QUOTA UNLIMITED ON USERS;
GRANT CREATE SESSION, CREATE TABLE TO OMNI_APP;
CREATE TABLE OMNI_APP.BIG AS SELECT LEVEL AS ID, LEVEL AS V FROM DUAL CONNECT BY LEVEL <= 1000;
CREATE TABLE OMNI_APP.MERGE_TARGET (ID NUMBER, V NUMBER);
`)
	got, err := d.CountAffectedRows(ctx, "UPDATE OMNI_APP.BIG SET V = 1 WHERE ID < 10;\n")
	if err != nil {
		t.Fatal(err)
	}
	if got <= 0 {
		t.Errorf("estimate = %d, want the plan's cardinality", got)
	}

	// A MERGE estimates the rows it reads, not the statement row dynamic
	// sampling understates.
	got, err = d.CountAffectedRows(ctx, "MERGE INTO OMNI_APP.MERGE_TARGET t USING (SELECT ID, V FROM OMNI_APP.BIG WHERE ID < 10) s ON (t.ID = s.ID) WHEN NOT MATCHED THEN INSERT (ID, V) VALUES (s.ID, s.V)")
	if err != nil {
		t.Fatal(err)
	}
	if got < 5 {
		t.Errorf("MERGE estimate = %d, want the rows below the MERGE operation", got)
	}

	// An INSERT ALL writes each row it reads to every INTO target.
	got, err = d.CountAffectedRows(ctx, "INSERT ALL INTO OMNI_APP.BIG (ID) VALUES (1) INTO OMNI_APP.BIG (ID) VALUES (2) SELECT ID FROM OMNI_APP.BIG")
	if err != nil {
		t.Fatal(err)
	}
	if got < 1000 {
		t.Errorf("INSERT ALL estimate = %d, want the rows it reads times its targets", got)
	}

	// An INSERT FIRST writes each row at most once, so its targets do not
	// multiply the estimate.
	got, err = d.CountAffectedRows(ctx, "INSERT FIRST WHEN ID < 5 THEN INTO OMNI_APP.BIG (ID) VALUES (ID) ELSE INTO OMNI_APP.BIG (ID) VALUES (ID) SELECT ID FROM OMNI_APP.BIG")
	if err != nil {
		t.Fatal(err)
	}
	if got >= 1500 {
		t.Errorf("INSERT FIRST estimate = %d, want the rows it reads, not their product with its targets", got)
	}
}
