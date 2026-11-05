package parser

import (
	"testing"
)

// TiDB's RECOMMEND INDEX, SET CONFIG, FLUSH TIDB PLUGINS, KILL TIDB, ALTER RANGE and
// CANCEL statements, and the log backup task and BR job arms of its BRIEStmt, are TiDB's
// alone: tidb/parser owns that grammar, and this file records that MySQL does not have it.
//
// Every statement below was measured against mysql:8.0 through PREPARE on 2026-11-05.

// tidbAdminStatements are statements only TiDB has.
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
	// KILL TIDB. KILL takes an expression here, so tidb is a column, and what follows it is the error.
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
				t.Errorf("omni accepts %q, which MySQL 8.0 refuses", sql)
			}
		})
	}
}

// TestTiDBAdminStatementsRefusedOracle holds the same list to the live server
// (container test).
func TestTiDBAdminStatementsRefusedOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, sql := range tidbAdminStatements {
		t.Run(sql, func(t *testing.T) {
			if o.prepareAccepts(t, sql) {
				t.Errorf("MySQL 8.0 accepts %q", sql)
			}
		})
	}
}
