package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// mysqlContainer is a test's connection to OceanBase: a replay of its tape, or
// the pinned server in a recording.
type mysqlContainer struct {
	db  *sql.DB
	ctx context.Context
}

// columnInfo holds a row from INFORMATION_SCHEMA.COLUMNS.
type columnInfo struct {
	Name, DataType, ColumnType, ColumnKey, Extra, Nullable string
	Position                                               int
	Default, Charset, Collation                            sql.NullString
	CharMaxLen, NumPrecision, NumScale                     sql.NullInt64
}

// indexInfo holds a row from INFORMATION_SCHEMA.STATISTICS.
type indexInfo struct {
	Name, ColumnName, IndexType, Nullable string
	NonUnique, SeqInIndex                 int
	Collation                             sql.NullString
}

// constraintInfo holds a row from INFORMATION_SCHEMA.TABLE_CONSTRAINTS.
type constraintInfo struct {
	Name, Type string
}

// startContainer returns t's connection to an OceanBase that holds nothing an
// earlier test left, with the database test selected
// (testcontainer.OpenEmpty). The connection closes with the test; the func
// returned is for the call sites that defer one.
func startContainer(t *testing.T) (*mysqlContainer, func()) {
	t.Helper()
	return &mysqlContainer{db: testcontainer.OpenEmpty(t), ctx: context.Background()}, func() {}
}

func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.OceanBase.Terminate()
	os.Exit(code)
}

// execSQL executes one or more SQL statements separated by semicolons.
// It respects quoted strings when splitting.
func (o *mysqlContainer) execSQL(sqlStr string) error {
	stmts := splitStatements(sqlStr)
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := o.db.ExecContext(o.ctx, stmt); err != nil {
			return fmt.Errorf("executing %q: %w", stmt, err)
		}
	}
	return nil
}

// execSQLDirect executes a single SQL statement directly without splitting on semicolons.
// This is needed for CREATE PROCEDURE/FUNCTION with BEGIN...END blocks.
func (o *mysqlContainer) execSQLDirect(sqlStr string) error {
	_, err := o.db.ExecContext(o.ctx, sqlStr)
	return err
}

// showCreateDatabase runs SHOW CREATE DATABASE and returns the CREATE DATABASE statement.
func (o *mysqlContainer) showCreateDatabase(database string) (string, error) {
	var dbName, createStmt string
	err := o.db.QueryRowContext(o.ctx, "SHOW CREATE DATABASE "+database).Scan(&dbName, &createStmt)
	if err != nil {
		return "", fmt.Errorf("SHOW CREATE DATABASE %s: %w", database, err)
	}
	return createStmt, nil
}

// showCreateTable runs SHOW CREATE TABLE and returns the CREATE TABLE statement.
func (o *mysqlContainer) showCreateTable(table string) (string, error) {
	var tableName, createStmt string
	err := o.db.QueryRowContext(o.ctx, "SHOW CREATE TABLE "+table).Scan(&tableName, &createStmt)
	if err != nil {
		return "", fmt.Errorf("SHOW CREATE TABLE %s: %w", table, err)
	}
	return createStmt, nil
}

// showCreateFunction runs SHOW CREATE FUNCTION and returns the CREATE FUNCTION statement.
func (o *mysqlContainer) showCreateFunction(name string) (string, error) {
	var funcName, sqlMode, createStmt, charSetClient, collConn, dbCollation string
	err := o.db.QueryRowContext(o.ctx, "SHOW CREATE FUNCTION "+name).Scan(
		&funcName, &sqlMode, &createStmt, &charSetClient, &collConn, &dbCollation)
	if err != nil {
		return "", fmt.Errorf("SHOW CREATE FUNCTION %s: %w", name, err)
	}
	return createStmt, nil
}

// showCreateProcedure runs SHOW CREATE PROCEDURE and returns the CREATE PROCEDURE statement.
func (o *mysqlContainer) showCreateProcedure(name string) (string, error) {
	var procName, sqlMode, createStmt, charSetClient, collConn, dbCollation string
	err := o.db.QueryRowContext(o.ctx, "SHOW CREATE PROCEDURE "+name).Scan(
		&procName, &sqlMode, &createStmt, &charSetClient, &collConn, &dbCollation)
	if err != nil {
		return "", fmt.Errorf("SHOW CREATE PROCEDURE %s: %w", name, err)
	}
	return createStmt, nil
}

// showCreateView runs SHOW CREATE VIEW and returns the CREATE VIEW statement.
func (o *mysqlContainer) showCreateView(name string) (string, error) {
	var viewName, createStmt, charSetClient, collConn string
	err := o.db.QueryRowContext(o.ctx, "SHOW CREATE VIEW "+name).Scan(
		&viewName, &createStmt, &charSetClient, &collConn)
	if err != nil {
		return "", fmt.Errorf("SHOW CREATE VIEW %s: %w", name, err)
	}
	return createStmt, nil
}

// queryColumns queries INFORMATION_SCHEMA.COLUMNS for the given table.
func (o *mysqlContainer) queryColumns(database, table string) ([]columnInfo, error) {
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT COLUMN_NAME, ORDINAL_POSITION, DATA_TYPE, COLUMN_TYPE,
		       IS_NULLABLE, COLUMN_DEFAULT, COLUMN_KEY, EXTRA,
		       CHARACTER_SET_NAME, COLLATION_NAME,
		       CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION`,
		database, table)
	if err != nil {
		return nil, fmt.Errorf("querying columns: %w", err)
	}
	defer rows.Close()

	var cols []columnInfo
	for rows.Next() {
		var c columnInfo
		if err := rows.Scan(
			&c.Name, &c.Position, &c.DataType, &c.ColumnType,
			&c.Nullable, &c.Default, &c.ColumnKey, &c.Extra,
			&c.Charset, &c.Collation,
			&c.CharMaxLen, &c.NumPrecision, &c.NumScale,
		); err != nil {
			return nil, fmt.Errorf("scanning column row: %w", err)
		}
		cols = append(cols, c)
	}
	return cols, rows.Err()
}

// queryIndexes queries INFORMATION_SCHEMA.STATISTICS for the given table.
func (o *mysqlContainer) queryIndexes(database, table string) ([]indexInfo, error) {
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT INDEX_NAME, SEQ_IN_INDEX, COLUMN_NAME, COLLATION,
		       NON_UNIQUE, INDEX_TYPE, NULLABLE
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY INDEX_NAME, SEQ_IN_INDEX`,
		database, table)
	if err != nil {
		return nil, fmt.Errorf("querying indexes: %w", err)
	}
	defer rows.Close()

	var idxs []indexInfo
	for rows.Next() {
		var idx indexInfo
		if err := rows.Scan(
			&idx.Name, &idx.SeqInIndex, &idx.ColumnName, &idx.Collation,
			&idx.NonUnique, &idx.IndexType, &idx.Nullable,
		); err != nil {
			return nil, fmt.Errorf("scanning index row: %w", err)
		}
		idxs = append(idxs, idx)
	}
	return idxs, rows.Err()
}

// queryConstraints queries INFORMATION_SCHEMA.TABLE_CONSTRAINTS for the given table.
func (o *mysqlContainer) queryConstraints(database, table string) ([]constraintInfo, error) {
	rows, err := o.db.QueryContext(o.ctx, `
		SELECT CONSTRAINT_NAME, CONSTRAINT_TYPE
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY CONSTRAINT_NAME`,
		database, table)
	if err != nil {
		return nil, fmt.Errorf("querying constraints: %w", err)
	}
	defer rows.Close()

	var cs []constraintInfo
	for rows.Next() {
		var c constraintInfo
		if err := rows.Scan(&c.Name, &c.Type); err != nil {
			return nil, fmt.Errorf("scanning constraint row: %w", err)
		}
		cs = append(cs, c)
	}
	return cs, rows.Err()
}

// queryStrings reads a result of one text column.
func queryStrings(ctr *mysqlContainer, query string) ([]string, error) {
	rows, err := ctr.db.QueryContext(ctr.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, rows.Err()
}

// splitStatements splits SQL text on semicolons, respecting single quotes,
// double quotes, and backtick-quoted identifiers.
func splitStatements(sqlStr string) []string {
	var stmts []string
	var current strings.Builder
	var inQuote rune // 0 means not in a quote
	var prevChar rune

	for _, ch := range sqlStr {
		switch {
		case inQuote != 0:
			current.WriteRune(ch)
			// End quote only if matching quote and not escaped by backslash.
			if ch == inQuote && prevChar != '\\' {
				inQuote = 0
			}
		case ch == '\'' || ch == '"' || ch == '`':
			inQuote = ch
			current.WriteRune(ch)
		case ch == ';':
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				stmts = append(stmts, stmt)
			}
			current.Reset()
		default:
			current.WriteRune(ch)
		}
		prevChar = ch
	}

	// Remaining text after the last semicolon (or if no semicolons).
	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		stmts = append(stmts, stmt)
	}
	return stmts
}

// normalizeWhitespace collapses runs of whitespace to a single space and trims.
func normalizeWhitespace(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func TestContainerSmoke(t *testing.T) {

	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Create a simple table.
	err := ctr.execSQL("CREATE TABLE t1 (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL)")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Verify SHOW CREATE TABLE works.
	createStmt, err := ctr.showCreateTable("t1")
	if err != nil {
		t.Fatalf("SHOW CREATE TABLE failed: %v", err)
	}
	if !strings.Contains(createStmt, "CREATE TABLE") {
		t.Errorf("expected CREATE TABLE in output, got: %s", createStmt)
	}
	t.Logf("SHOW CREATE TABLE t1:\n%s", createStmt)

	// Verify queryColumns works.
	cols, err := ctr.queryColumns("test", "t1")
	if err != nil {
		t.Fatalf("queryColumns failed: %v", err)
	}
	if len(cols) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(cols))
	}
	if cols[0].Name != "id" {
		t.Errorf("expected first column 'id', got %q", cols[0].Name)
	}
	if cols[1].Name != "name" {
		t.Errorf("expected second column 'name', got %q", cols[1].Name)
	}
	if cols[1].Nullable != "NO" {
		t.Errorf("expected 'name' to be NOT NULL, got Nullable=%q", cols[1].Nullable)
	}
}

func TestSplitStatements(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "CREATE TABLE t1 (id INT); CREATE TABLE t2 (id INT)",
			want:  []string{"CREATE TABLE t1 (id INT)", "CREATE TABLE t2 (id INT)"},
		},
		{
			input: "INSERT INTO t1 VALUES ('a;b'); SELECT 1",
			want:  []string{"INSERT INTO t1 VALUES ('a;b')", "SELECT 1"},
		},
		{
			input: `SELECT "col;name" FROM t1`,
			want:  []string{`SELECT "col;name" FROM t1`},
		},
		{
			input: "SELECT `col;name` FROM t1",
			want:  []string{"SELECT `col;name` FROM t1"},
		},
		{
			input: "",
			want:  nil,
		},
		{
			input: "  ;  ;  ",
			want:  nil,
		},
	}
	for _, tt := range tests {
		got := splitStatements(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("splitStatements(%q): got %d stmts, want %d", tt.input, len(got), len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("splitStatements(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestNormalizeWhitespace(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"  hello   world  ", "hello world"},
		{"a\n\tb", "a b"},
		{"already clean", "already clean"},
		{"", ""},
	}
	for _, tt := range tests {
		got := normalizeWhitespace(tt.input)
		if got != tt.want {
			t.Errorf("normalizeWhitespace(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
