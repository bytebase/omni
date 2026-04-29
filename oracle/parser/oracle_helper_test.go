package parser

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/sijms/go-ora/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// oracleDB wraps a real Oracle Database container connection for cross-validation testing.
type oracleDB struct {
	db      *sql.DB
	adminDB *sql.DB
	ctx     context.Context
}

var (
	oracleOnce    sync.Once
	oracleInst    *oracleDB
	oracleInitErr error
)

// startOracleDB starts a shared Oracle Database 23ai Free container. The container
// is reused across all tests via sync.Once. Returns nil and skips the test if
// Docker is unavailable or the container fails to start.
func startOracleDB(t *testing.T) *oracleDB {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping Oracle DB test in short mode")
	}

	oracleOnce.Do(func() {
		ctx := context.Background()

		req := testcontainers.ContainerRequest{
			Image:        "gvenzl/oracle-free:23-slim-faststart",
			ExposedPorts: []string{"1521/tcp"},
			Env: map[string]string{
				"ORACLE_PASSWORD":   "testpass",
				"APP_USER":          "testuser",
				"APP_USER_PASSWORD": "testpass",
			},
			WaitingFor: wait.ForLog("DATABASE IS READY TO USE!").
				WithStartupTimeout(5 * time.Minute),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			oracleInitErr = fmt.Errorf("failed to start Oracle container: %w", err)
			return
		}

		host, err := container.Host(ctx)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			oracleInitErr = fmt.Errorf("failed to get container host: %w", err)
			return
		}

		mappedPort, err := container.MappedPort(ctx, "1521")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			oracleInitErr = fmt.Errorf("failed to get mapped port: %w", err)
			return
		}

		// go-ora connection string format: oracle://user:pass@host:port/service_name
		connStr := fmt.Sprintf("oracle://testuser:testpass@%s:%s/FREEPDB1", host, mappedPort.Port())

		db, err := sql.Open("oracle", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			oracleInitErr = fmt.Errorf("failed to open database: %w", err)
			return
		}

		if err := db.PingContext(ctx); err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			oracleInitErr = fmt.Errorf("failed to ping Oracle database: %w", err)
			return
		}
		adminDB, err := openOracleAdminDB(ctx, host, mappedPort.Port())
		if err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			oracleInitErr = fmt.Errorf("failed to open Oracle admin database: %w", err)
			return
		}

		oracleInst = &oracleDB{db: db, adminDB: adminDB, ctx: ctx}
		setupOracleSchema(oracleInst)
		// Container cleanup happens when the process exits.
	})

	if oracleInitErr != nil {
		t.Fatalf("Oracle container not available: %v", oracleInitErr)
	}
	return oracleInst
}

func openOracleAdminDB(ctx context.Context, host, port string) (*sql.DB, error) {
	connStr := fmt.Sprintf("oracle://system:testpass@%s:%s/FREEPDB1", host, port)
	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// canExecute checks whether Oracle DB accepts the given SQL statement.
// It uses EXPLAIN PLAN to validate the syntax without executing the statement.
// For DDL statements that cannot be explained, it uses DBMS_SQL.PARSE()
// to validate syntax without execution, avoiding side effects.
//
// NOTE: For DDL that must actually be executed (e.g., Stage 6 catalog tests),
// use execSQL() with explicit cleanup instead.
func (o *oracleDB) canExecute(sqlStr string) error {
	// Try EXPLAIN PLAN first — works for DML (SELECT, INSERT, UPDATE, DELETE, MERGE).
	explainSQL := fmt.Sprintf("EXPLAIN PLAN FOR %s", sqlStr)
	_, err := o.db.ExecContext(o.ctx, explainSQL)
	if err == nil {
		return nil
	}

	// For DDL and other statements that EXPLAIN PLAN doesn't support,
	// use DBMS_SQL.PARSE to validate syntax without executing.
	// DBMS_SQL.PARSE parses the statement and raises an exception on syntax errors,
	// but does not execute it — no side effects on the database state.
	parseCheck := fmt.Sprintf(`DECLARE
  c INTEGER := DBMS_SQL.OPEN_CURSOR;
BEGIN
  DBMS_SQL.PARSE(c, '%s', DBMS_SQL.NATIVE);
  DBMS_SQL.CLOSE_CURSOR(c);
EXCEPTION
  WHEN OTHERS THEN
    DBMS_SQL.CLOSE_CURSOR(c);
    RAISE;
END;`, strings.ReplaceAll(sqlStr, "'", "''"))
	_, parseErr := o.db.ExecContext(o.ctx, parseCheck)
	return parseErr
}

// execSQL executes a SQL statement on the Oracle container, failing the test on error.
func (o *oracleDB) execSQL(t *testing.T, sqlStr string) {
	t.Helper()
	_, err := o.db.ExecContext(o.ctx, sqlStr)
	if err != nil {
		t.Fatalf("Oracle exec failed: %v\nSQL: %s", err, truncateOracleSQL(sqlStr, 500))
	}
}

// setupOracleSchema creates tables and objects needed by the SQL corpus
// so that EXPLAIN PLAN / DBMS_SQL.PARSE can validate statements that
// reference these objects. Without this, most DML/queries get rejected
// with ORA-00942 (table or view does not exist).
func setupOracleSchema(o *oracleDB) {
	ddls := []string{
		// HR-style schema used by most corpus statements
		`CREATE TABLE departments (
			department_id NUMBER PRIMARY KEY,
			department_name VARCHAR2(100),
			location_id NUMBER
		)`,
		`CREATE TABLE employees (
			employee_id NUMBER PRIMARY KEY,
			first_name VARCHAR2(50),
			last_name VARCHAR2(100) NOT NULL,
			email VARCHAR2(100),
			salary NUMBER(10,2),
			commission_pct NUMBER(4,2),
			department_id NUMBER REFERENCES departments(department_id),
			manager_id NUMBER,
			job_id VARCHAR2(20),
			created_date DATE DEFAULT SYSDATE
		)`,
		`CREATE TABLE customers (
			customer_id NUMBER PRIMARY KEY,
			name VARCHAR2(200)
		)`,
		`CREATE TABLE orders (
			order_id NUMBER PRIMARY KEY,
			customer_id NUMBER REFERENCES customers(customer_id),
			order_date DATE,
			amount NUMBER(10,2)
		)`,
		`CREATE TABLE t (
			id NUMBER,
			a NUMBER,
			b NUMBER,
			c NUMBER,
			x NUMBER,
			doc CLOB
		)`,
		`CREATE TABLE s (
			id NUMBER,
			x NUMBER
		)`,
		// Tables referenced in DML corpus
		`CREATE TABLE emp_backup AS SELECT * FROM employees WHERE 1=0`,
		`CREATE TABLE high_earners AS SELECT * FROM employees WHERE 1=0`,
		`CREATE TABLE mid_earners AS SELECT * FROM employees WHERE 1=0`,
		`CREATE TABLE low_earners AS SELECT * FROM employees WHERE 1=0`,
		`CREATE TABLE sales (prod_id NUMBER, amount NUMBER)`,
		`CREATE TABLE employee_bonuses (employee_id NUMBER, bonus NUMBER)`,
		// Tables referenced in admin corpus
		`CREATE TABLE temp_results (id NUMBER, result VARCHAR2(200))`,
		`CREATE TABLE temp_data (id NUMBER)`,
		`CREATE TABLE old_table (id NUMBER)`,
		`CREATE TABLE test_table (id NUMBER)`,
		// Views/objects for advanced queries
		`CREATE VIEW sales_view AS SELECT 'US' AS country, 'Bounce' AS product, 2023 AS year, 1000 AS sales FROM dual`,
		`CREATE TABLE j_purchaseorder (po_document CLOB)`,
	}

	for _, ddl := range ddls {
		_, _ = o.db.ExecContext(o.ctx, ddl)
		// Ignore errors (table may already exist from a previous run)
	}
}

// truncateOracleSQL truncates a SQL string for error messages.
func truncateOracleSQL(sqlStr string, maxLen int) string {
	if len(sqlStr) <= maxLen {
		return sqlStr
	}
	return sqlStr[:maxLen] + "..."
}
