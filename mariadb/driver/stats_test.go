package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"errors"
	"io"
	"slices"
	"strings"
	"sync"
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

func TestRequireOneStatement(t *testing.T) {
	if err := requireOneStatement("DELETE FROM t WHERE id = 1"); err != nil {
		t.Fatalf("one statement = %v", err)
	}
	// The probe prefixes EXPLAIN, and the connection sends every statement of a
	// request, so a second one would run.
	for _, statement := range []string{"DELETE FROM t; DROP TABLE t", "", "DELETE FROM"} {
		if err := requireOneStatement(statement); err == nil {
			t.Errorf("requireOneStatement(%q) = nil", statement)
		}
	}
}

type fakeStatsConnector struct {
	mu              sync.Mutex
	queries         []string
	queryConns      []int
	conns           int
	setVersionErr   error
	explainErr      error
	jsonPlan        string
	traditionalCols []string
	traditionalRows [][]sqldriver.Value
}

func (c *fakeStatsConnector) Connect(context.Context) (sqldriver.Conn, error) {
	c.mu.Lock()
	c.conns++
	id := c.conns
	c.mu.Unlock()
	return &fakeStatsConn{c: c, id: id}, nil
}

func (*fakeStatsConnector) Driver() sqldriver.Driver { return nil }

type fakeStatsConn struct {
	c  *fakeStatsConnector
	id int
}

func (*fakeStatsConn) Prepare(string) (sqldriver.Stmt, error) { return nil, sqldriver.ErrSkip }
func (*fakeStatsConn) Close() error                           { return nil }
func (*fakeStatsConn) Begin() (sqldriver.Tx, error)           { return nil, sqldriver.ErrSkip }

func (f *fakeStatsConn) record(query string) {
	f.c.mu.Lock()
	f.c.queries = append(f.c.queries, query)
	f.c.queryConns = append(f.c.queryConns, f.id)
	f.c.mu.Unlock()
}

func (f *fakeStatsConn) ExecContext(_ context.Context, query string, _ []sqldriver.NamedValue) (sqldriver.Result, error) {
	f.record(query)
	if query == "SET SESSION explain_json_format_version = 1" && f.c.setVersionErr != nil {
		return nil, f.c.setVersionErr
	}
	return sqldriver.RowsAffected(0), nil
}

func (f *fakeStatsConn) QueryContext(_ context.Context, query string, _ []sqldriver.NamedValue) (sqldriver.Rows, error) {
	f.record(query)
	if f.c.explainErr != nil {
		return nil, f.c.explainErr
	}
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
	sqlDB.SetMaxIdleConns(0)
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
	if conn.queryConns[0] != conn.queryConns[1] {
		t.Fatalf("SET and EXPLAIN ran on different connections: %v", conn.queryConns)
	}
}

func TestCountAffectedRowsFromJSONAndConnectionPinning(t *testing.T) {
	const (
		deleteJoinPlan = `{"query_block":{"select_id":1,"nested_loop":[
			{"table":{"table_name":"s","access_type":"ALL","loops":1,"rows":100,"filtered":10}},
			{"table":{"table_name":"big","access_type":"ref","loops":10,"rows":100,"filtered":100}}
		]}}`
		insertSelectPlan = `{"query_block":{"select_id":1,"table":{"insert":true,"table_name":"t2","access_type":"ALL"},"insert_from":{"table":{"table_name":"t","access_type":"ALL","rows":1000,"loops":1,"filtered":10}}}}`
	)
	unknownVarErr := &gomysql.MySQLError{Number: erUnknownSystemVariable, Message: "Unknown system variable 'explain_json_format_version'"}
	for _, tc := range []struct {
		name          string
		statement     string
		plan          string
		setVersionErr error
		want          int64
	}{
		{
			name:          "multi-table delete without explain_json_format_version",
			statement:     "DELETE big FROM big JOIN s ON big.s_id = s.id WHERE s.flag = 1;",
			plan:          deleteJoinPlan,
			setVersionErr: unknownVarErr,
			want:          1000,
		},
		{
			name:          "multi-table delete capped by LIMIT",
			statement:     "DELETE big FROM big JOIN s ON big.s_id = s.id WHERE s.flag = 1 LIMIT 3;",
			plan:          deleteJoinPlan,
			setVersionErr: unknownVarErr,
			want:          3,
		},
		{
			name:      "insert select pins connection for set and reset",
			statement: "INSERT INTO t2 SELECT * FROM t WHERE d = 3;",
			plan:      insertSelectPlan,
			want:      100,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			conn := &fakeStatsConnector{
				jsonPlan:      tc.plan,
				setVersionErr: tc.setVersionErr,
			}
			sqlDB := sql.OpenDB(conn)
			defer sqlDB.Close()
			sqlDB.SetMaxIdleConns(0)
			d := &Driver{db: sqlDB}

			got, err := d.CountAffectedRows(context.Background(), tc.statement)
			if err != nil {
				t.Fatalf("CountAffectedRows = %v", err)
			}
			if got != tc.want {
				t.Fatalf("CountAffectedRows = %d, want %d", got, tc.want)
			}
			for _, id := range conn.queryConns {
				if id != conn.queryConns[0] {
					t.Fatalf("SET/EXPLAIN/reset ran across multiple connections: %v", conn.queryConns)
				}
			}
		})
	}
}

func TestCountAffectedRowsFallbackEmptyReturnsJSONError(t *testing.T) {
	conn := &fakeStatsConnector{
		setVersionErr:   &gomysql.MySQLError{Number: erUnknownSystemVariable, Message: "Unknown system variable 'explain_json_format_version'"},
		jsonPlan:        `{"query_block":{"select_id":1,"table":{"insert":true,"select_id":1,"table_name":"t2","access_type":"ALL"}}}`,
		traditionalCols: []string{"id", "select_type", "table", "type", "rows", "filtered"},
	}
	sqlDB := sql.OpenDB(conn)
	defer sqlDB.Close()
	d := &Driver{db: sqlDB}

	_, err := d.CountAffectedRows(context.Background(), "INSERT INTO t2 SELECT id, ROW_NUMBER() OVER () FROM t;")
	if err == nil || !strings.Contains(err.Error(), `table "t2" in the plan has no row estimate`) {
		t.Fatalf("CountAffectedRows err = %v, want missing row estimate error", err)
	}
	wantQueries := []string{
		"SET SESSION explain_json_format_version = 1",
		"EXPLAIN FORMAT=JSON INSERT INTO t2 SELECT id, ROW_NUMBER() OVER () FROM t;",
		"EXPLAIN FORMAT=TRADITIONAL INSERT INTO t2 SELECT id, ROW_NUMBER() OVER () FROM t;",
	}
	if !slices.Equal(conn.queries, wantQueries) {
		t.Fatalf("queries = %v, want %v", conn.queries, wantQueries)
	}
}

func TestExplainJSONErrors(t *testing.T) {
	t.Run("non-1193 SET error fails fast", func(t *testing.T) {
		conn := &fakeStatsConnector{
			setVersionErr: &gomysql.MySQLError{Number: 1227, Message: "Access denied"},
		}
		sqlDB := sql.OpenDB(conn)
		defer sqlDB.Close()
		d := &Driver{db: sqlDB}

		_, err := d.ExplainJSON(context.Background(), "SELECT 1")
		if err == nil || !strings.Contains(err.Error(), "failed to set explain_json_format_version") {
			t.Fatalf("ExplainJSON err = %v", err)
		}
	})

	t.Run("refused EXPLAIN returns server error", func(t *testing.T) {
		conn := &fakeStatsConnector{
			explainErr: errors.New("Table 'test.missing' doesn't exist"),
		}
		sqlDB := sql.OpenDB(conn)
		defer sqlDB.Close()
		d := &Driver{db: sqlDB}

		_, err := d.CountAffectedRows(context.Background(), "UPDATE missing SET c = 1 WHERE id = 1;")
		if err == nil || !strings.Contains(err.Error(), "Table 'test.missing' doesn't exist") {
			t.Fatalf("CountAffectedRows err = %v", err)
		}
	})
}
