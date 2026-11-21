package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync and testdata/dump from this run")

func openTestDriver(t *testing.T, dbName string) *Driver {
	t.Helper()
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS " + quoteBacktick(dbName),
			"CREATE DATABASE " + quoteBacktick(dbName),
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.OpenWith(t, testcontainer.Options{
		Database:        dbName,
		MultiStatements: true,
	})
	t.Cleanup(func() { _ = db.Close() })
	return &Driver{db: db, databaseName: dbName}
}

func execCorpus(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	for _, chunk := range strings.Split(script, "\n\n") {
		if stmt := strings.TrimSpace(chunk); stmt != "" {
			if _, err := db.ExecContext(context.Background(), stmt); err != nil {
				t.Fatalf("%s:\n%v", stmt, err)
			}
		}
	}
}

func TestSyncDatabase(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t, "omni_sync_corpus")
	execCorpus(t, d.db, string(corpus))
	got, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	golden := filepath.Join("testdata", "sync", "default.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want := readGolden(t, golden)
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("SyncDatabase differs from %s (-want +got):\n%s", golden, diff)
	}
}

func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(b, m); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return m
}

func TestSyncInstance(t *testing.T) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		dbRows, err := db.QueryContext(ctx, "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA")
		if err != nil {
			return err
		}
		var extraDBs []string
		for dbRows.Next() {
			var name string
			if err := dbRows.Scan(&name); err != nil {
				_ = dbRows.Close()
				return err
			}
			lower := strings.ToLower(name)
			if lower == "information_schema" || lower == "mysql" || lower == "oceanbase" || lower == "sys" || lower == "test" || strings.HasPrefix(lower, "__") {
				continue
			}
			extraDBs = append(extraDBs, name)
		}
		_ = dbRows.Close()
		for _, name := range extraDBs {
			rRows, err := db.QueryContext(ctx, "SELECT ROUTINE_TYPE, ROUTINE_NAME FROM INFORMATION_SCHEMA.ROUTINES WHERE ROUTINE_SCHEMA = ?", name)
			if err != nil {
				return err
			}
			var dropRoutines []string
			for rRows.Next() {
				var rt, rn string
				if err := rRows.Scan(&rt, &rn); err != nil {
					_ = rRows.Close()
					return err
				}
				dropRoutines = append(dropRoutines, fmt.Sprintf("DROP %s IF EXISTS %s.%s", rt, quoteBacktick(name), quoteBacktick(rn)))
			}
			_ = rRows.Close()
			for _, stmt := range dropRoutines {
				if _, err := db.ExecContext(ctx, stmt); err != nil {
					return err
				}
			}
			if _, err := db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+quoteBacktick(name)); err != nil {
				return err
			}
		}

		uRows, err := db.QueryContext(ctx, "SELECT user, host FROM mysql.user WHERE user LIKE 'sync_%'")
		if err != nil {
			return err
		}
		var extraUsers [][2]string
		for uRows.Next() {
			var u, h string
			if err := uRows.Scan(&u, &h); err != nil {
				_ = uRows.Close()
				return err
			}
			extraUsers = append(extraUsers, [2]string{u, h})
		}
		_ = uRows.Close()
		for _, u := range extraUsers {
			if _, err := db.ExecContext(ctx, fmt.Sprintf("DROP USER IF EXISTS %s@%s", quoteSingle(u[0]), quoteSingle(u[1]))); err != nil {
				return err
			}
		}

		for _, stmt := range []string{
			"CREATE DATABASE sync_inst_db CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci",
			"CREATE USER 'sync_reader'@'%' IDENTIFIED BY 'reader-pass'",
			"GRANT SELECT ON sync_inst_db.* TO 'sync_reader'@'%'",
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.OpenWith(t, testcontainer.Options{MultiStatements: true})
	t.Cleanup(func() { _ = db.Close() })
	d := &Driver{db: db}

	inst, err := d.SyncInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var foundDB string
	for _, m := range inst.Databases {
		if m.Name == "sync_inst_db" {
			foundDB = fmt.Sprintf("%s %s %s", m.Name, m.CharacterSet, m.Collation)
		}
	}
	var roleList []string
	for _, r := range inst.Roles {
		if strings.HasPrefix(r.Name, "'sync_") {
			roleList = append(roleList, r.Name+": "+*r.Attribute)
		}
	}
	b, err := json.MarshalIndent(map[string]any{
		"version":             inst.Version,
		"lowerCaseTableNames": inst.LowerCaseTableNames,
		"database":            foundDB,
		"roles":               roleList,
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	out := string(b) + "\n"

	golden := filepath.Join("testdata", "sync", "instance.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(out), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), out); diff != "" {
		t.Errorf("SyncInstance differs from %s (-want +got):\n%s", golden, diff)
	}
}

func TestSyncDatabaseNotFound(t *testing.T) {
	db := testcontainer.OpenWith(t, testcontainer.Options{})
	t.Cleanup(func() { _ = db.Close() })
	d := &Driver{db: db, databaseName: "omni_no_such_database"}
	if _, err := d.SyncDatabase(context.Background()); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase missing db err = %v, want driver.ErrNotFound", err)
	}
}

func TestSyncForeignKeyReferencedDatabase(t *testing.T) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		for _, stmt := range []string{
			"DROP DATABASE IF EXISTS omni_fk_local",
			"DROP DATABASE IF EXISTS omni_fk_remote",
			"CREATE DATABASE omni_fk_local",
			"CREATE DATABASE omni_fk_remote",
		} {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
		return nil
	})
	db := testcontainer.OpenWith(t, testcontainer.Options{Database: "omni_fk_local", MultiStatements: true})
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.ExecContext(context.Background(), `
		CREATE TABLE omni_fk_remote.parent (id INT PRIMARY KEY);
		CREATE TABLE omni_fk_local.parent (id INT PRIMARY KEY);
		CREATE TABLE omni_fk_local.child (
			id INT PRIMARY KEY,
			local_id INT,
			qualified_id INT,
			remote_id INT,
			CONSTRAINT fk_local FOREIGN KEY (local_id) REFERENCES parent (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
			CONSTRAINT fk_qualified FOREIGN KEY (qualified_id) REFERENCES omni_fk_local.parent (id) ON DELETE SET NULL ON UPDATE RESTRICT,
			CONSTRAINT fk_remote FOREIGN KEY (remote_id) REFERENCES omni_fk_remote.parent (id) ON DELETE CASCADE ON UPDATE RESTRICT
		);
	`); err != nil {
		t.Fatal(err)
	}

	d := &Driver{db: db, databaseName: "omni_fk_local"}
	fks, err := d.getForeignKeyList(context.Background(), "omni_fk_local")
	if err != nil {
		t.Fatal(err)
	}
	want := []*metadatapb.ForeignKeyMetadata{
		{Name: "fk_local", Columns: []string{"local_id"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "RESTRICT", OnUpdate: "RESTRICT", MatchType: "NONE"},
		{Name: "fk_qualified", Columns: []string{"qualified_id"}, ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "SET NULL", OnUpdate: "RESTRICT", MatchType: "NONE"},
		{Name: "fk_remote", Columns: []string{"remote_id"}, ReferencedSchema: "omni_fk_remote", ReferencedTable: "parent", ReferencedColumns: []string{"id"}, OnDelete: "CASCADE", OnUpdate: "RESTRICT", MatchType: "NONE"},
	}
	if diff := cmp.Diff(want, fks["child"], protocmp.Transform()); diff != "" {
		t.Errorf("foreign keys diff (-want +got):\n%s", diff)
	}
}

// TestTriggerBody pins the header ACTION_STATEMENT repeats. The pinned server's answers
// cover a plain body and a BEGIN…END one; a quoted name, and a body that holds the header's
// own words, pin the scan.
func TestTriggerBody(t *testing.T) {
	for _, tc := range []struct {
		name, statement, want string
	}{
		{
			name:      "single statement",
			statement: "TRIGGER trg_1 BEFORE INSERT ON trg_t FOR EACH ROW SET NEW.v = 1",
			want:      "SET NEW.v = 1",
		},
		{
			name:      "compound statement",
			statement: "TRIGGER trg_1 BEFORE INSERT ON trg_t FOR EACH ROW BEGIN IF NEW.v < 0 THEN SET NEW.v = 0; END IF; END",
			want:      "BEGIN IF NEW.v < 0 THEN SET NEW.v = 0; END IF; END",
		},
		{
			name:      "quoted names",
			statement: "TRIGGER `trg odd` AFTER DELETE ON `trg_t` FOR EACH ROW INSERT INTO trg_t (id) VALUES (1)",
			want:      "INSERT INTO trg_t (id) VALUES (1)",
		},
		{
			name:      "body holds the header's words",
			statement: "TRIGGER t1 BEFORE UPDATE ON t FOR EACH ROW SET NEW.x = 'FOR EACH ROW'",
			want:      "SET NEW.x = 'FOR EACH ROW'",
		},
		{
			name:      "already a body",
			statement: "SET NEW.v = 1",
			want:      "SET NEW.v = 1",
		},
		{
			name:      "no header terminator",
			statement: "TRIGGER t1 BEFORE UPDATE ON t",
			want:      "TRIGGER t1 BEFORE UPDATE ON t",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := triggerBody(tc.statement); got != tc.want {
				t.Errorf("triggerBody(%q) = %q, want %q", tc.statement, got, tc.want)
			}
		})
	}
}

// TestTemporalDefaultWithPrecision pins the precision a temporal default takes from its
// column type, which OceanBase leaves out of COLUMN_DEFAULT.
func TestTemporalDefaultWithPrecision(t *testing.T) {
	for _, tc := range []struct {
		value, columnType, want string
	}{
		{"CURRENT_TIMESTAMP", "datetime(6)", "CURRENT_TIMESTAMP(6)"},
		{"CURRENT_TIMESTAMP", "timestamp(3)", "CURRENT_TIMESTAMP(3)"},
		{"CURRENT_TIMESTAMP", "datetime", "CURRENT_TIMESTAMP"},
		{"CURRENT_TIMESTAMP(6)", "datetime(6)", "CURRENT_TIMESTAMP(6)"},
		{"CURRENT_DATE", "date", "CURRENT_DATE"},
		{"CURRENT_TIME", "time(2)", "CURRENT_TIME(2)"},
	} {
		if got := temporalDefaultWithPrecision(tc.value, tc.columnType); got != tc.want {
			t.Errorf("temporalDefaultWithPrecision(%q, %q) = %q, want %q", tc.value, tc.columnType, got, tc.want)
		}
	}
}

// TestEscapeEventComment pins how a comment information_schema.EVENTS states becomes a
// string literal under the event's own sql_mode.
func TestEscapeEventComment(t *testing.T) {
	const strictMode = "STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_AUTO_CREATE_USER"
	for _, tc := range []struct {
		comment, sqlMode, want string
	}{
		{"owner's once", strictMode, "owner''s once"},
		{`C:\new path`, strictMode, `C:\\new path`},
		{`trailing\`, strictMode, `trailing\\`},
		{`C:\new path`, "NO_BACKSLASH_ESCAPES," + strictMode, `C:\new path`},
		{`owner's \new`, "NO_BACKSLASH_ESCAPES", `owner''s \new`},
	} {
		if got := escapeEventComment(tc.comment, tc.sqlMode); got != tc.want {
			t.Errorf("escapeEventComment(%q, %q) = %q, want %q", tc.comment, tc.sqlMode, got, tc.want)
		}
	}
}

// eventRowsDriver answers the event query with the rows a test sets, so the
// recreated-database case needs no server.
type eventRowsDriver struct{ rows [][]sqldriver.Value }

func (eventRowsDriver) Open(string) (sqldriver.Conn, error) { return eventRowsConn{}, nil }

type eventRowsConn struct{}

func (eventRowsConn) Prepare(string) (sqldriver.Stmt, error) { return nil, sqldriver.ErrSkip }
func (eventRowsConn) Close() error                           { return nil }
func (eventRowsConn) Begin() (sqldriver.Tx, error)           { return nil, sqldriver.ErrSkip }

func (eventRowsConn) QueryContext(context.Context, string, []sqldriver.NamedValue) (sqldriver.Rows, error) {
	return &eventRows{d: testEventRows}, nil
}

// eventRows returns the duplicate rows a dropped database's event and the
// database created again under its name both list.
type eventRows struct {
	d    [][]sqldriver.Value
	idx  int
	done bool
}

func (r *eventRows) Columns() []string {
	return []string{"EVENT_NAME", "EVENT_TYPE", "EXECUTE_AT", "INTERVAL_VALUE", "INTERVAL_FIELD", "STARTS", "ENDS", "STATUS", "ON_COMPLETION", "EVENT_COMMENT", "TIME_ZONE", "SQL_MODE", "EVENT_DEFINITION"}
}
func (r *eventRows) Close() error { return nil }
func (r *eventRows) Next(dest []sqldriver.Value) error {
	if r.idx >= len(r.d) {
		return io.EOF
	}
	copy(dest, r.d[r.idx])
	r.idx++
	return nil
}

var testEventRows = [][]sqldriver.Value{
	{"e", "RECURRING", "", "1", "DAY", "2030-01-01 00:00:00", "", "DISABLED", "PRESERVE", "", "SYSTEM", "ONLY_FULL_GROUP_BY", "INSERT INTO t VALUES (1)"},
	{"e", "RECURRING", "", "1", "DAY", "2030-01-01 00:00:00", "", "DISABLED", "PRESERVE", "", "SYSTEM", "ONLY_FULL_GROUP_BY", "INSERT INTO t VALUES (1)"},
}

// TestGetEventListDeduplicatesRecreatedDatabaseEvents pins the case
// information_schema.EVENTS cannot state: a database created again under a
// dropped one's name lists that database's event until OceanBase forgets it,
// and the sync must list the name once.
func TestGetEventListDeduplicatesRecreatedDatabaseEvents(t *testing.T) {
	sql.Register("test_oceanbase_event_rows", eventRowsDriver{})
	db, err := sql.Open("test_oceanbase_event_rows", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	d := &Driver{db: db, databaseName: "db"}
	events, err := d.getEventList(context.Background(), "db")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("getEventList returned %d events, want 1: %+v", len(events), events)
	}
}

// TestExternalTableFacts reads the LOCATION, PATTERN and FORMAT clauses SHOW
// CREATE TABLE states, and leaves them empty when the statement has none.
func TestExternalTableFacts(t *testing.T) {
	ddl := "CREATE EXTERNAL TABLE `ext_t` (\n" +
		"  `c1` int(11) GENERATED ALWAYS AS (metadata$filecol1),\n" +
		"  `c2` varchar(10) GENERATED ALWAYS AS (metadata$filecol2)\n" +
		") \n" +
		"LOCATION='file:///tmp'\n" +
		"AUTO_REFRESH = OFF\n" +
		"FORMAT (\n" +
		"  TYPE = 'CSV',\n" +
		"  ENCODING = 'utf8mb4'\n" +
		") DEFAULT CHARSET = utf8mb4"
	location, pattern, autoRefresh, format := externalTableFacts(ddl)
	if location != "file:///tmp" {
		t.Errorf("location = %q, want %q", location, "file:///tmp")
	}
	if format != "TYPE = 'CSV', ENCODING = 'utf8mb4'" {
		t.Errorf("format = %q, want %q", format, "TYPE = 'CSV', ENCODING = 'utf8mb4'")
	}
	if pattern != "" {
		t.Errorf("pattern = %q, want empty", pattern)
	}
	if autoRefresh != "OFF" {
		t.Errorf("auto refresh = %q, want %q", autoRefresh, "OFF")
	}

	location, pattern, autoRefresh, format = externalTableFacts("CREATE EXTERNAL TABLE `ext_t` (`c1` int) LOCATION='file:///tmp' PATTERN='.*[.]csv' AUTO_REFRESH = IMMEDIATE")
	if location != "file:///tmp" || pattern != ".*[.]csv" || autoRefresh != "IMMEDIATE" || format != "" {
		t.Errorf("without FORMAT = %q, %q, %q, %q; want %q, %q, %q, empty", location, pattern, autoRefresh, format, "file:///tmp", ".*[.]csv", "IMMEDIATE")
	}

	// A MySQL backslash escape in a column comment must not end the quote
	// before LOCATION and FORMAT.
	ddl = "CREATE EXTERNAL TABLE `ext_t` (`c1` int COMMENT 'owner\\'s note') " +
		"LOCATION='file:///tmp/it\\'s' FORMAT (TYPE = 'CSV')"
	location, pattern, autoRefresh, format = externalTableFacts(ddl)
	if location != "file:///tmp/it's" {
		t.Errorf("escaped location = %q, want %q", location, "file:///tmp/it's")
	}
	if pattern != "" {
		t.Errorf("escaped pattern = %q, want empty", pattern)
	}
	if autoRefresh != "" {
		t.Errorf("escaped auto refresh = %q, want empty", autoRefresh)
	}
	if format != "TYPE = 'CSV'" {
		t.Errorf("escaped format = %q, want %q", format, "TYPE = 'CSV'")
	}
}

type stubQueryDriver struct {
	err error
}

func (s stubQueryDriver) Open(string) (sqldriver.Conn, error) { return stubQueryConn(s), nil }

type stubQueryConn struct {
	err error
}

func (stubQueryConn) Prepare(string) (sqldriver.Stmt, error) { return nil, sqldriver.ErrSkip }
func (stubQueryConn) Close() error                           { return nil }
func (stubQueryConn) Begin() (sqldriver.Tx, error)           { return nil, sqldriver.ErrSkip }
func (c stubQueryConn) QueryContext(context.Context, string, []sqldriver.NamedValue) (sqldriver.Rows, error) {
	return nil, c.err
}

var (
	stubMViewErrMu sync.Mutex
	stubMViewErr   error
)

type stubMViewSwitchDriver struct{}

func (stubMViewSwitchDriver) Open(string) (sqldriver.Conn, error) { return stubMViewSwitchConn{}, nil }

type stubMViewSwitchConn struct{}

func (stubMViewSwitchConn) Prepare(string) (sqldriver.Stmt, error) { return nil, sqldriver.ErrSkip }
func (stubMViewSwitchConn) Close() error                           { return nil }
func (stubMViewSwitchConn) Begin() (sqldriver.Tx, error)           { return nil, sqldriver.ErrSkip }
func (stubMViewSwitchConn) QueryContext(context.Context, string, []sqldriver.NamedValue) (sqldriver.Rows, error) {
	stubMViewErrMu.Lock()
	err := stubMViewErr
	stubMViewErrMu.Unlock()
	return nil, err
}

func init() {
	sql.Register("test_oceanbase_mview_err", stubMViewSwitchDriver{})
}

func TestListMaterializedViewLogsErrorHandling(t *testing.T) {
	db, err := sql.Open("test_oceanbase_mview_err", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	d := &Driver{db: db, databaseName: "db"}

	for _, tc := range []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "privilege denied suppressed",
			err:     &gomysql.MySQLError{Number: 1142, Message: "SELECT command denied"},
			wantErr: false,
		},
		{
			name:    "catalog table missing suppressed",
			err:     &gomysql.MySQLError{Number: 1146, Message: "Table doesn't exist"},
			wantErr: false,
		},
		{
			name:    "transient connection error propagated",
			err:     &gomysql.MySQLError{Number: 1105, Message: "unknown internal failure"},
			wantErr: true,
		},
		{
			name:    "generic driver error propagated",
			err:     errors.New("connection reset by peer"),
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stubMViewErrMu.Lock()
			stubMViewErr = tc.err
			stubMViewErrMu.Unlock()

			hasLog, logTables, gotErr := d.listMaterializedViewLogs(context.Background(), "db")
			if tc.wantErr {
				if gotErr == nil {
					t.Fatalf("listMaterializedViewLogs() err = nil, want %v", tc.err)
				}
				return
			}
			if gotErr != nil || hasLog != nil || logTables != nil {
				t.Fatalf("listMaterializedViewLogs() = (%v, %v, %v), want (nil, nil, nil)", hasLog, logTables, gotErr)
			}
		})
	}
}
