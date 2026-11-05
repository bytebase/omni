package parser

import (
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TiDB's RECOMMEND INDEX, SET CONFIG, FLUSH TIDB PLUGINS, KILL TIDB, ALTER RANGE and
// CANCEL statements, and the log backup task and BR job arms of its BRIEStmt, are TiDB's
// alone: tidb/parser owns that grammar, and this file records that MariaDB does not have it.
//
// Every statement below was measured against mariadb:11.8.8 through PREPARE on 2026-11-05.

// tidbAdminStatements are statements only TiDB has, refused with a 1064.
var tidbAdminStatements = []string{
	// RECOMMEND INDEX, the advisor.
	"RECOMMEND INDEX RUN",
	"RECOMMEND INDEX RUN WITH a = 1",
	"RECOMMEND INDEX RUN FOR 'x'",
	"RECOMMEND INDEX SHOW",
	"RECOMMEND INDEX APPLY 1",
	"RECOMMEND INDEX IGNORE 1",
	"RECOMMEND INDEX SET a = 1",
	// SET CONFIG, which is not an assignment here.
	"SET CONFIG tikv log.level = 'info'",
	"SET CONFIG '127.0.0.1:3306' a = 1",
	"SET CONFIG tikv a-b.c = ON",
	// The log backup tasks and BR jobs.
	"PAUSE BACKUP LOGS",
	"RESUME BACKUP LOGS",
	"STOP BACKUP LOGS",
	"PURGE BACKUP LOGS FROM 'noop://'",
	"SHOW BACKUP LOGS STATUS",
	"SHOW BACKUP LOGS METADATA FROM 'noop://'",
	"SHOW BACKUP METADATA FROM 'noop://'",
	"SHOW BR JOB 1",
	"SHOW BR JOB QUERY 1",
	"CANCEL BR JOB 1",
	// CANCEL DISTRIBUTION JOB.
	"CANCEL DISTRIBUTION JOB 1",
	// FLUSH TIDB PLUGINS.
	"FLUSH TIDB PLUGINS a",
	"FLUSH TIDB PLUGINS a, b",
	"FLUSH LOCAL TIDB PLUGINS a",
	// ALTER RANGE.
	"ALTER RANGE global PLACEMENT POLICY p",
	"ALTER RANGE global PLACEMENT POLICY DEFAULT",
	"ALTER RANGE global PLACEMENT POLICY SET DEFAULT",
}

// tidbKillStatements are refused too, but not with a 1064. MariaDB's KILL takes an
// expression and looks a name up while it parses, so tidb answers 1054 "Unknown column"
// before the server reads what follows it. The same statement with a syntax error after
// it answers 1054 again, which is how TestTiDBAdminStatementsRefusedOracle tells that
// from a statement that parsed.
var tidbKillStatements = []string{
	"KILL TIDB 1",
	"KILL TIDB CONNECTION 1",
	"KILL TIDB QUERY 1",
	"KILL TIDB CONNECTION_ID()",
}

// TestTiDBAdminStatementsRefused holds this parser to tidbAdminStatements. It needs no
// server, so it guards the grammar wherever it runs.
func TestTiDBAdminStatementsRefused(t *testing.T) {
	for _, sql := range tidbAdminStatements {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("omni accepts %q, which MariaDB 11.8.8 refuses", sql)
			}
		})
	}
}

// TestTiDBKillStatementsRefused is the same for the KILL TIDB forms.
func TestTiDBKillStatementsRefused(t *testing.T) {
	for _, sql := range tidbKillStatements {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("omni accepts %q, which MariaDB 11.8.8 refuses", sql)
			}
		})
	}
}

// TestTiDBAdminStatementsRefusedOracle holds both lists to the live server
// (container test).
func TestTiDBAdminStatementsRefusedOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, sql := range tidbAdminStatements {
		t.Run(sql, func(t *testing.T) {
			if o.prepareAccepts(t, sql) {
				t.Errorf("MariaDB 11.8.8 accepts %q", sql)
			}
		})
	}
	prepareCode := func(t *testing.T, sql string) uint16 {
		t.Helper()
		conn, err := o.db.Conn(o.ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		if _, err := conn.ExecContext(o.ctx, "SET @omni_kill = ?", sql); err != nil {
			t.Fatal(err)
		}
		_, err = conn.ExecContext(o.ctx, "PREPARE omni_kill_stmt FROM @omni_kill")
		if err == nil {
			_, _ = conn.ExecContext(o.ctx, "DEALLOCATE PREPARE omni_kill_stmt")
			return 0
		}
		var myErr *gomysql.MySQLError
		if !errors.As(err, &myErr) {
			t.Fatalf("no verdict for %q: %v", sql, err)
		}
		return myErr.Number
	}
	for _, sql := range tidbKillStatements {
		t.Run(sql, func(t *testing.T) {
			if got := prepareCode(t, sql); got != 1054 {
				t.Errorf("MariaDB 11.8.8 answers %d for %q, want 1054", got, sql)
			}
			if got := prepareCode(t, sql+"\n; SELECT FROM"); got != 1054 {
				t.Errorf("with a syntax error after it MariaDB 11.8.8 answers %d for %q: it read past the column", got, sql)
			}
		})
	}
}
