package parser

import "testing"

// ============================================================================
// MariaDB 11.8 RETURNING conformance oracle (BYT-9135 P2).
//
// Full parity (0 OVER / 0 GAP): RETURNING is accepted on INSERT/REPLACE and
// single-table DELETE, rejected on UPDATE and multi-table DELETE, and is a
// RESERVED word (cannot be a column name or alias — a MySQL-8.0 divergence).
// The reserved-word reject edges are the durable guard: if a future change
// un-reserves RETURNING (or a keyword refactor regresses it), `INSERT … SELECT …
// RETURNING` silently mis-parses (the SELECT eats RETURNING as an alias) and a
// pile of column/alias over-accepts return — all caught here.
// ============================================================================

var returningOracleCorpus = []string{
	// --- accept: INSERT (every form) ---
	"INSERT INTO rt (id, name) VALUES (1, 'a') RETURNING id",
	"INSERT INTO rt (id, name) VALUES (2, 'b') RETURNING id, name, id + 100 AS big",
	"INSERT INTO rt (id, name) VALUES (3, 'c') RETURNING *",
	"INSERT INTO rt SELECT 9, 'q' RETURNING id, name",
	"INSERT INTO rt (SELECT 1, 'x') RETURNING id",
	"INSERT INTO rt SET id = 1, name = 'x' RETURNING id",
	"INSERT INTO rt (id, name) VALUES (4, 'd') ON DUPLICATE KEY UPDATE name = 'x' RETURNING id",
	"INSERT INTO rt (id) VALUES (1) RETURNING id, (SELECT COUNT(*) FROM rt2)",
	"INSERT INTO rt (id) VALUES (1) RETURNING COUNT(*)",
	// --- accept: REPLACE (donorless arm) ---
	"REPLACE INTO rt (id, name) VALUES (1, 'z') RETURNING id, name",
	"REPLACE INTO rt (id, name) VALUES (1, 'z') RETURNING *",
	// --- accept: single-table DELETE ---
	"DELETE FROM rt WHERE id = 3 RETURNING id, name",
	"DELETE FROM rt WHERE id = 2 RETURNING *",
	"DELETE FROM rt RETURNING COUNT(*)",

	// --- reject (1064): no UPDATE RETURNING / multi-table DELETE / bare list ---
	"UPDATE rt SET name = 'x' WHERE id = 1 RETURNING id",
	"DELETE rt FROM rt JOIN rt2 ON rt.id = rt2.id RETURNING rt.id",       // multi-table delete, syntax 1 (JOIN)
	"DELETE FROM rt USING rt JOIN rt2 ON rt.id = rt2.id RETURNING rt.id", // multi-table delete, syntax 2 (USING → the Using==nil guard half)
	"INSERT INTO rt (id) VALUES (1) RETURNING",

	// --- reject (1064): RETURNING is reserved — not a column name or alias ---
	"CREATE TABLE rret (returning INT)",
	"SELECT returning FROM rt",
	"SELECT 1 RETURNING",
	"SELECT 1 AS RETURNING",
	"SELECT 1 RETURNING x",
}

// TestMariaDBReturningOracle asserts omni agrees with a live MariaDB 11.8.8 on
// every RETURNING statement (Short-gated container test).
func TestMariaDBReturningOracle(t *testing.T) {
	o := startMariaDB(t)

	for _, ddl := range []string{
		"CREATE TABLE IF NOT EXISTS rt (id INT PRIMARY KEY, name VARCHAR(20))",
		"CREATE TABLE IF NOT EXISTS rt2 (id INT, name VARCHAR(20))",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}

	for _, sqlStr := range returningOracleCorpus {
		_, perr := Parse(sqlStr)
		omniAccepts := perr == nil
		v, code, msg := o.classify(sqlStr)
		if v == vErrored {
			t.Errorf("ERRORED (infra, not a parse verdict) for %q: %s", sqlStr, msg)
			continue
		}
		mdbAccepts := v == vAccepted
		if omniAccepts != mdbAccepts {
			t.Errorf("DIVERGENCE omni=%v mariadb=%v (code %d) for %q: %s",
				omniAccepts, mdbAccepts, code, sqlStr, msg)
		}
	}
}
