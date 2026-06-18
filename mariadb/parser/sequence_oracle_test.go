package parser

import "testing"

// ============================================================================
// MariaDB 11.8 SEQUENCE conformance oracle (BYT-9135).
//
// Unlike TestMariaDBDivergenceInventory (which pins a subtractive allowlist),
// the sequence surface is FULL PARITY: omni must agree with the container on
// every statement (0 OVER, 0 GAP). sequenceOracleCorpus is the consolidated
// container-truth set — the option matrix, the NEXT/PREVIOUS VALUE FOR position
// matrix, and the non-reserved-keyword identifier edges — each verified parse
// accept/reject against a live mariadb:11.8.8. Any disagreement fails.
//
// The identifier edges are the durable regression guard: a future keyword
// addition (P3 system-versioning adds PERIOD/SYSTEM_TIME on the same
// keywordCategories surface) that accidentally reserves one of these words would
// break its column/alias use and surface here as a live divergence.
// ============================================================================

var sequenceOracleCorpus = []string{
	// --- CREATE SEQUENCE: flags + name forms (accept) ---
	"CREATE SEQUENCE cs_plain",
	"CREATE SEQUENCE IF NOT EXISTS cs_ine",
	"CREATE OR REPLACE SEQUENCE cs_cor",
	"CREATE OR REPLACE SEQUENCE IF NOT EXISTS cs_both",
	"CREATE TEMPORARY SEQUENCE cs_tmp",
	// --- START / INCREMENT: WITH | BY | = | bare | signed ---
	"CREATE OR REPLACE SEQUENCE cs1 START WITH 1000",
	"CREATE OR REPLACE SEQUENCE cs2 START = 5",
	"CREATE OR REPLACE SEQUENCE cs3 START 5",
	"CREATE OR REPLACE SEQUENCE cs4 INCREMENT BY 5",
	"CREATE OR REPLACE SEQUENCE cs5 INCREMENT = 5",
	"CREATE OR REPLACE SEQUENCE cs6 INCREMENT 5",
	"CREATE OR REPLACE SEQUENCE cs7 INCREMENT BY -1 MAXVALUE -1 MINVALUE -100",
	"CREATE OR REPLACE SEQUENCE cs8 INCREMENT BY +5 START WITH +5",
	// --- MIN/MAX/CACHE/CYCLE + NO-forms ---
	"CREATE OR REPLACE SEQUENCE cs9 MINVALUE 10 MAXVALUE 9999",
	"CREATE OR REPLACE SEQUENCE cs10 MINVALUE = 1 MAXVALUE = 99",
	"CREATE OR REPLACE SEQUENCE cs11 NOMINVALUE NOMAXVALUE",
	"CREATE OR REPLACE SEQUENCE cs12 NO MINVALUE NO MAXVALUE",
	"CREATE OR REPLACE SEQUENCE cs13 CACHE 50",
	"CREATE OR REPLACE SEQUENCE cs14 CACHE = 50",
	"CREATE OR REPLACE SEQUENCE cs15 NOCACHE",
	"CREATE OR REPLACE SEQUENCE cs16 MINVALUE 1 MAXVALUE 5 CYCLE",
	"CREATE OR REPLACE SEQUENCE cs17 MINVALUE 1 MAXVALUE 5 NOCYCLE",
	"CREATE OR REPLACE SEQUENCE cs18 START WITH 100 MINVALUE 10 MAXVALUE 100000 INCREMENT BY 5 CACHE 20 CYCLE",
	// --- AS int_type (any-order) ---
	"CREATE OR REPLACE SEQUENCE cs19 AS BIGINT",
	"CREATE OR REPLACE SEQUENCE cs20 AS INT UNSIGNED",
	"CREATE OR REPLACE SEQUENCE cs21 START WITH 1 AS BIGINT",
	"CREATE OR REPLACE SEQUENCE cs22 AS TINYINT START WITH 1",

	// --- CREATE rejects (1064) ---
	"CREATE SEQUENCE cr1 RESTART",
	"CREATE SEQUENCE cr2 RESTART WITH 5",
	"CREATE SEQUENCE cr3 CACHE",
	"CREATE SEQUENCE cr4 INCREMENT",
	"CREATE SEQUENCE cr5 START",
	"CREATE SEQUENCE cr6 INCREMENT BY",
	"CREATE SEQUENCE cr7 NO CACHE",
	"CREATE SEQUENCE cr8 NO CYCLE",
	"CREATE SEQUENCE cr9 AS VARCHAR(10)",
	"CREATE SEQUENCE cr10 AS DECIMAL(10)",
	"CREATE SEQUENCE cr11 AS INT(11)",
	"CREATE SEQUENCE cr12 AS DATE",
	"CREATE SEQUENCE cr13 START WITH 1+1",
	"CREATE SEQUENCE cr14 START WITH (5)",
	"CREATE SEQUENCE cr15 CACHE CYCLE",
	"CREATE SEQUENCE",

	// --- ALTER SEQUENCE (accept) ---
	"ALTER SEQUENCE sq RESTART",
	"ALTER SEQUENCE sq RESTART WITH 500",
	"ALTER SEQUENCE sq RESTART = 5",
	"ALTER SEQUENCE sq RESTART 5",
	"ALTER SEQUENCE sq INCREMENT BY 10",
	"ALTER SEQUENCE sq MINVALUE 1 NOCACHE",
	"ALTER SEQUENCE IF EXISTS sq RESTART WITH 1",
	"ALTER SEQUENCE sq START WITH 5",
	"ALTER SEQUENCE sq AS BIGINT",
	"ALTER SEQUENCE sq RESTART INCREMENT BY 5",
	"ALTER SEQUENCE sq INCREMENT BY 5 RESTART",
	// --- ALTER rejects (1064) ---
	"ALTER SEQUENCE sq", // requires at least one option

	// --- DROP SEQUENCE ---
	"DROP SEQUENCE IF EXISTS d_a",
	"DROP SEQUENCE IF EXISTS d_a, d_b",
	"DROP TEMPORARY SEQUENCE IF EXISTS d_c",
	"DROP SEQUENCE d_a,", // trailing comma (1064)

	// --- NEXT / PREVIOUS VALUE FOR positions (all parse-accept; OVER rejects) ---
	"SELECT NEXT VALUE FOR sq",
	"SELECT PREVIOUS VALUE FOR sq",
	"SELECT NEXT VALUE FOR sq + 1",
	"INSERT INTO pt (id, v) VALUES (NEXT VALUE FOR sq, 1)",
	"INSERT INTO pt (id) VALUES (PREVIOUS VALUE FOR sq)",
	"UPDATE pt SET id = NEXT VALUE FOR sq WHERE v = 1",
	"SELECT * FROM pt WHERE id = NEXT VALUE FOR sq",
	"SELECT id FROM pt HAVING id = NEXT VALUE FOR sq",
	"SELECT id FROM pt ORDER BY NEXT VALUE FOR sq",
	"CREATE OR REPLACE TABLE dseq1 (id INT DEFAULT NEXT VALUE FOR sq, v INT)",
	"CREATE OR REPLACE TABLE dseq2 (id INT DEFAULT (NEXT VALUE FOR sq), v INT)",
	"CREATE TABLE cseq1 (id INT CHECK (id < NEXT VALUE FOR sq))",
	"CREATE TABLE gseq1 (id INT, g INT AS (NEXT VALUE FOR sq))",
	"SELECT ABS(NEXT VALUE FOR sq)",
	// --- value-function rejects (1064) ---
	"SELECT NEXT VALUE FOR sq OVER (ORDER BY 1)",
	"SELECT next value",     // bare NEXT VALUE without FOR
	"SELECT previous value", // bare PREVIOUS VALUE without FOR

	// --- identifier edges: non-reserved keywords as columns/aliases (durable guard) ---
	"SELECT next FROM pt",
	"SELECT previous FROM pt",
	"SELECT id AS next FROM pt",
	"SELECT id AS previous FROM pt",
	"SELECT next, value FROM pt",
	"CREATE OR REPLACE TABLE icols (sequence INT, previous INT, increment INT, cycle INT, minvalue INT, nocache INT, nocycle INT, nominvalue INT, nomaxvalue INT)",

	// --- adversarial accepts: case-insensitivity, quoting, comments ---
	"select nExT vAlUe fOr sq",
	"CREATE OR REPLACE SEQUENCE `seq two` START WITH 1",
	"SELECT NEXT /* c */ VALUE FOR sq",
}

// TestMariaDBSequenceOracle asserts omni's parser agrees with a live MariaDB
// 11.8.8 on every statement in sequenceOracleCorpus (Short-gated container test).
func TestMariaDBSequenceOracle(t *testing.T) {
	o := startMariaDB(t)

	// Seed the objects the position/identifier cases reference so accepts execute
	// cleanly; non-existence would still classify as accepted (4091 != 1064), but
	// seeding keeps the signal a clean code-0 accept.
	for _, ddl := range []string{
		"CREATE TABLE IF NOT EXISTS pt (id INT, v INT)",
		"CREATE SEQUENCE IF NOT EXISTS sq",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}

	for _, sqlStr := range sequenceOracleCorpus {
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
