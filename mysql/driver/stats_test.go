package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"io"
	"slices"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

func TestDMLStatement(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		wantErr   string
	}{
		{name: "insert", statement: "INSERT INTO t VALUES (1)"},
		{name: "replace", statement: "REPLACE INTO t VALUES (1)"},
		{name: "update", statement: "UPDATE t SET c = 1 WHERE id = 1"},
		{name: "delete", statement: "DELETE FROM t WHERE id = 1"},
		{name: "select", statement: "SELECT * FROM t", wantErr: "expected an INSERT, UPDATE or DELETE"},
		{name: "script", statement: "DELETE FROM t WHERE id = 1; SELECT 1", wantErr: "expected exactly one statement"},
		{name: "empty", statement: "", wantErr: "expected exactly one statement"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := dmlStatement(tc.statement)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("dmlStatement(%q) = (%T, %v), want error containing %q", tc.statement, stmt, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("dmlStatement(%q) = %v", tc.statement, err)
			}
			if stmt == nil {
				t.Fatalf("dmlStatement(%q) returned no statement", tc.statement)
			}
		})
	}
}

type fakeStatsConnector struct {
	queries         []string
	setVersionErr   error
	jsonPlan        string
	traditionalCols []string
	traditionalRows [][]sqldriver.Value
}

func (c *fakeStatsConnector) Connect(context.Context) (sqldriver.Conn, error) {
	return &fakeStatsConn{c: c}, nil
}

func (*fakeStatsConnector) Driver() sqldriver.Driver { return nil }

type fakeStatsConn struct {
	c *fakeStatsConnector
}

func (*fakeStatsConn) Prepare(string) (sqldriver.Stmt, error) { return nil, sqldriver.ErrSkip }
func (*fakeStatsConn) Close() error                           { return nil }
func (*fakeStatsConn) Begin() (sqldriver.Tx, error)           { return nil, sqldriver.ErrSkip }

func (f *fakeStatsConn) ExecContext(_ context.Context, query string, _ []sqldriver.NamedValue) (sqldriver.Result, error) {
	f.c.queries = append(f.c.queries, query)
	if query == "SET SESSION explain_json_format_version = 1" && f.c.setVersionErr != nil {
		return nil, f.c.setVersionErr
	}
	return sqldriver.RowsAffected(0), nil
}

func (f *fakeStatsConn) QueryContext(_ context.Context, query string, _ []sqldriver.NamedValue) (sqldriver.Rows, error) {
	f.c.queries = append(f.c.queries, query)
	if strings.HasPrefix(query, "EXPLAIN FORMAT=JSON ") {
		return &fakeStatsRows{
			cols: []string{"EXPLAIN"},
			rows: [][]sqldriver.Value{{f.c.jsonPlan}},
		}, nil
	}
	if strings.HasPrefix(query, "EXPLAIN FORMAT=TRADITIONAL ") {
		return &fakeStatsRows{
			cols: f.c.traditionalCols,
			rows: f.c.traditionalRows,
		}, nil
	}
	return nil, sqldriver.ErrSkip
}

type fakeStatsRows struct {
	cols []string
	rows [][]sqldriver.Value
	idx  int
}

func (r *fakeStatsRows) Columns() []string { return r.cols }
func (*fakeStatsRows) Close() error        { return nil }
func (r *fakeStatsRows) Next(dest []sqldriver.Value) error {
	if r.idx >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.idx])
	r.idx++
	return nil
}

func TestCountAffectedRowsFallsBackToTraditionalExplain(t *testing.T) {
	conn := &fakeStatsConnector{
		jsonPlan:        `{"query_block":{"select_id":1,"windowing":{"windows":[]}}}`,
		traditionalCols: []string{"id", "select_type", "table", "rows", "filtered"},
		traditionalRows: [][]sqldriver.Value{
			{"1", "SIMPLE", "source_rows", "10", "50.00"},
			{"1", "SIMPLE", "joined_rows", "4", "50.00"},
			{"2", "SUBQUERY", "ignored_block", "999", "100.00"},
		},
	}
	sqlDB := sql.OpenDB(conn)
	defer sqlDB.Close()
	d := &Driver{db: sqlDB}

	got, err := d.CountAffectedRows(context.Background(), "INSERT INTO target_rows SELECT id, ROW_NUMBER() OVER (ORDER BY id) FROM source_rows;")
	if err != nil {
		t.Fatalf("CountAffectedRows = %v", err)
	}
	// First query block product: 10 * 0.5 * 4 * 0.5 = 10; block 2 is ignored.
	if got != 10 {
		t.Fatalf("CountAffectedRows = %d, want 10", got)
	}
	wantQueries := []string{
		"SET SESSION explain_json_format_version = 1",
		"EXPLAIN FORMAT=JSON INSERT INTO target_rows SELECT id, ROW_NUMBER() OVER (ORDER BY id) FROM source_rows;",
		"SET SESSION explain_json_format_version = DEFAULT",
		"EXPLAIN FORMAT=TRADITIONAL INSERT INTO target_rows SELECT id, ROW_NUMBER() OVER (ORDER BY id) FROM source_rows;",
	}
	if !slices.Equal(conn.queries, wantQueries) {
		t.Fatalf("queries = %v, want %v", conn.queries, wantQueries)
	}
}

func TestExplainJSONToleratesUnknownFormatVersion(t *testing.T) {
	conn := &fakeStatsConnector{
		setVersionErr: &gomysql.MySQLError{Number: erUnknownSystemVariable, Message: "Unknown system variable 'explain_json_format_version'"},
		jsonPlan:      `{"query_block":{"select_id":1}}`,
	}
	sqlDB := sql.OpenDB(conn)
	defer sqlDB.Close()
	d := &Driver{db: sqlDB}

	plan, err := d.ExplainJSON(context.Background(), "SELECT 1")
	if err != nil {
		t.Fatalf("ExplainJSON = %v", err)
	}
	if plan != conn.jsonPlan {
		t.Fatalf("ExplainJSON = %q, want %q", plan, conn.jsonPlan)
	}
	wantQueries := []string{
		"SET SESSION explain_json_format_version = 1",
		"EXPLAIN FORMAT=JSON SELECT 1",
	}
	if !slices.Equal(conn.queries, wantQueries) {
		t.Fatalf("queries = %v, want %v", conn.queries, wantQueries)
	}
}
