package validate

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/bytebase/omni/mysql/parser"
)

type routineValidationOracle struct {
	db  *sql.DB
	ctx context.Context
}

func startRoutineValidationOracle(t *testing.T) *routineValidationOracle {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}

	ctx := context.Background()
	container, err := tcmysql.Run(ctx, "mysql:8.0",
		tcmysql.WithDatabase("test"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("failed to start MySQL container: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })

	connStr, err := container.ConnectionString(ctx, "parseTime=true", "multiStatements=true")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	return &routineValidationOracle{db: db, ctx: ctx}
}

func TestRoutineValidationOracleAlignment(t *testing.T) {
	o := startRoutineValidationOracle(t)
	mustExecRoutineOracle(t, o, "DROP TABLE IF EXISTS t")
	mustExecRoutineOracle(t, o, "CREATE TABLE t (id INT PRIMARY KEY, v INT, c VARCHAR(50))")
	mustExecRoutineOracle(t, o, "DROP TABLE IF EXISTS audit_log")
	mustExecRoutineOracle(t, o, "CREATE TABLE audit_log (id INT PRIMARY KEY AUTO_INCREMENT, action VARCHAR(50))")

	type probe struct {
		name    string
		dropSQL string
		sql     string
	}

	dropProc := func(name string) string { return "DROP PROCEDURE IF EXISTS " + name }
	dropFunc := func(name string) string { return "DROP FUNCTION IF EXISTS " + name }
	dropTrig := func(name string) string { return "DROP TRIGGER IF EXISTS " + name }
	dropEvent := func(name string) string { return "DROP EVENT IF EXISTS " + name }

	probes := []probe{
		{
			name:    "procedure with simple BEGIN body",
			dropSQL: dropProc("p_ok"),
			sql:     "CREATE PROCEDURE p_ok() BEGIN SELECT 1; END",
		},
		{
			name:    "function with RETURN",
			dropSQL: dropFunc("f_ok"),
			sql:     "CREATE FUNCTION f_ok(a INT) RETURNS INT DETERMINISTIC RETURN a + 1",
		},
		{
			name:    "trigger with required FOR EACH ROW",
			dropSQL: dropTrig("trg_ok"),
			sql:     "CREATE TRIGGER trg_ok BEFORE INSERT ON t FOR EACH ROW SET NEW.v = COALESCE(NEW.v, 0)",
		},
		{
			name:    "event with DO body",
			dropSQL: dropEvent("ev_ok"),
			sql:     "CREATE EVENT ev_ok ON SCHEDULE EVERY 1 HOUR DO BEGIN INSERT INTO audit_log(action) VALUES ('tick'); END",
		},
		{
			name:    "LEAVE undeclared label",
			dropSQL: dropProc("p_leave_undeclared"),
			sql:     "CREATE PROCEDURE p_leave_undeclared() BEGIN LEAVE nowhere; END",
		},
		{
			name:    "ITERATE undeclared label",
			dropSQL: dropProc("p_iterate_undeclared"),
			sql:     "CREATE PROCEDURE p_iterate_undeclared() BEGIN ITERATE nowhere; END",
		},
		{
			name:    "ITERATE BEGIN label",
			dropSQL: dropProc("p_iterate_begin_label"),
			sql:     "CREATE PROCEDURE p_iterate_begin_label() my_block: BEGIN ITERATE my_block; END my_block",
		},
		{
			name:    "OPEN undeclared cursor",
			dropSQL: dropProc("p_open_undeclared"),
			sql:     "CREATE PROCEDURE p_open_undeclared() BEGIN OPEN nope; END",
		},
		{
			name:    "FETCH undeclared cursor",
			dropSQL: dropProc("p_fetch_undeclared"),
			sql:     "CREATE PROCEDURE p_fetch_undeclared() BEGIN FETCH nope INTO @x; END",
		},
		{
			name:    "CLOSE undeclared cursor",
			dropSQL: dropProc("p_close_undeclared"),
			sql:     "CREATE PROCEDURE p_close_undeclared() BEGIN CLOSE nope; END",
		},
		{
			name:    "HANDLER references undeclared condition",
			dropSQL: dropProc("p_handler_condition_undeclared"),
			sql:     "CREATE PROCEDURE p_handler_condition_undeclared() BEGIN DECLARE EXIT HANDLER FOR no_such SET @e = 1; SELECT 1; END",
		},
		{
			name:    "SET undeclared local variable",
			dropSQL: dropProc("p_set_undeclared"),
			sql:     "CREATE PROCEDURE p_set_undeclared() BEGIN SET nope = 1; END",
		},
		{
			name:    "duplicate DECLARE variable",
			dropSQL: dropProc("p_duplicate_var"),
			sql:     "CREATE PROCEDURE p_duplicate_var() BEGIN DECLARE x INT; DECLARE x INT; SELECT 1; END",
		},
		{
			name:    "duplicate DECLARE cursor",
			dropSQL: dropProc("p_duplicate_cursor"),
			sql:     "CREATE PROCEDURE p_duplicate_cursor() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE c CURSOR FOR SELECT 2; SELECT 1; END",
		},
		{
			name:    "duplicate DECLARE condition",
			dropSQL: dropProc("p_duplicate_condition"),
			sql:     "CREATE PROCEDURE p_duplicate_condition() BEGIN DECLARE dk CONDITION FOR SQLSTATE '23000'; DECLARE dk CONDITION FOR SQLSTATE '23001'; SELECT 1; END",
		},
		{
			name:    "duplicate nested label",
			dropSQL: dropProc("p_duplicate_label"),
			sql:     "CREATE PROCEDURE p_duplicate_label() lbl: BEGIN lbl: BEGIN SELECT 1; END lbl; END lbl",
		},
		{
			name:    "duplicate HANDLER condition value",
			dropSQL: dropProc("p_duplicate_handler_condition"),
			sql:     "CREATE PROCEDURE p_duplicate_handler_condition() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '23000', SQLSTATE '23000' SET @e = 1; SELECT 1; END",
		},
		{
			name:    "function with no RETURN",
			dropSQL: dropFunc("f_no_return"),
			sql:     "CREATE FUNCTION f_no_return() RETURNS INT DETERMINISTIC BEGIN SELECT 1; END",
		},
		{
			name:    "function with SIGNAL but no RETURN",
			dropSQL: dropFunc("f_signal_only"),
			sql:     "CREATE FUNCTION f_signal_only() RETURNS INT DETERMINISTIC BEGIN SIGNAL SQLSTATE '45000'; END",
		},
		{
			name:    "RETURN inside procedure",
			dropSQL: dropProc("p_return"),
			sql:     "CREATE PROCEDURE p_return() BEGIN RETURN 1; END",
		},
		{
			name:    "RETURN inside trigger",
			dropSQL: dropTrig("trg_return"),
			sql:     "CREATE TRIGGER trg_return BEFORE INSERT ON t FOR EACH ROW BEGIN RETURN 1; END",
		},
		{
			name:    "RETURN inside event",
			dropSQL: dropEvent("ev_return"),
			sql:     "CREATE EVENT ev_return ON SCHEDULE EVERY 1 HOUR DO BEGIN RETURN 1; END",
		},
		{
			name:    "LEAVE outer label from handler body",
			dropSQL: dropProc("p_handler_leave_outer"),
			sql: `CREATE PROCEDURE p_handler_leave_outer() my_block: BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN LEAVE my_block; END;
    SELECT 1;
END my_block`,
		},
		{
			name:    "function with RETURN on one branch",
			dropSQL: dropFunc("f_return_one_branch"),
			sql:     "CREATE FUNCTION f_return_one_branch(x INT) RETURNS INT DETERMINISTIC BEGIN IF x > 0 THEN RETURN 1; END IF; END",
		},
	}

	for _, p := range probes {
		t.Run(p.name, func(t *testing.T) {
			_, _ = o.db.ExecContext(o.ctx, p.dropSQL)

			omniAccept, omniDetail := parseAndValidateAccepts(p.sql)
			mysqlAccept, mysqlDetail := mysqlExecAccepts(o, p.sql)
			if omniAccept != mysqlAccept {
				t.Fatalf("acceptance mismatch\nomni accept: %v (%s)\nmysql accept: %v (%s)\nsql: %s",
					omniAccept, omniDetail, mysqlAccept, mysqlDetail, p.sql)
			}
		})
	}
}

func parseAndValidateAccepts(sql string) (bool, string) {
	list, err := parser.Parse(sql)
	if err != nil {
		return false, trimFirstLine(err.Error())
	}
	diags := Validate(list, Options{})
	if len(diags) == 0 {
		return true, "ok"
	}
	codes := make([]string, 0, len(diags))
	for _, d := range diags {
		codes = append(codes, d.Code)
	}
	return false, strings.Join(codes, ",")
}

func mysqlExecAccepts(o *routineValidationOracle, sql string) (bool, string) {
	_, err := o.db.ExecContext(o.ctx, sql)
	if err == nil {
		return true, "ok"
	}
	if myErr, ok := err.(*gomysql.MySQLError); ok {
		return false, "ERR " + mysqlOracleItoa(int(myErr.Number)) + ": " + trimFirstLine(myErr.Message)
	}
	return false, trimFirstLine(err.Error())
}

func mustExecRoutineOracle(t *testing.T, o *routineValidationOracle, sql string) {
	t.Helper()
	if _, err := o.db.ExecContext(o.ctx, sql); err != nil {
		t.Fatalf("seed failed: %v\nsql: %s", err, sql)
	}
}

func trimFirstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func mysqlOracleItoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
