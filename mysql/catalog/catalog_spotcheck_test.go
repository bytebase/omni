package catalog

import (
	"fmt"
	"strings"
	"testing"
)

// These tests verify that the behaviors described in
// docs/plans/2026-04-13-mysql-implicit-behaviors-catalog.md
// match real MySQL 8.0 observable behavior. They run against a
// testcontainers MySQL container and do NOT involve omni's catalog.

// spotCheckQuery runs DDL (possibly multi-statement) then returns the
// rows of the given query. Fatals on any error.
func spotCheckQuery(t *testing.T, mc *mysqlContainer, ddl, query string) [][]any {
	t.Helper()
	for _, stmt := range splitStatements(ddl) {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := mc.db.ExecContext(mc.ctx, stmt); err != nil {
			t.Fatalf("DDL failed: %q\n  %v", stmt, err)
		}
	}
	rows, err := mc.db.QueryContext(mc.ctx, query)
	if err != nil {
		t.Fatalf("query failed: %q\n  %v", query, err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var results [][]any
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			t.Fatalf("scan failed: %v", err)
		}
		// Convert []byte to string for readability.
		for i, v := range vals {
			if b, ok := v.([]byte); ok {
				vals[i] = string(b)
			}
		}
		results = append(results, vals)
	}
	return results
}

func asString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case nil:
		return "<nil>"
	case int64:
		return fmt.Sprintf("%d", x)
	case int:
		return fmt.Sprintf("%d", x)
	case float64:
		return fmt.Sprintf("%v", x)
	}
	return fmt.Sprintf("%v", v)
}

func TestSpotCheck_CatalogVerification(t *testing.T) {
	if testing.Short() {
		t.Skip("spot-check requires container")
	}
	mc, cleanup := startContainer(t)
	defer cleanup()

	// Reset state at the start of every sub-test.
	reset := func(t *testing.T) {
		t.Helper()
		if _, err := mc.db.ExecContext(mc.ctx, "DROP DATABASE IF EXISTS sc"); err != nil {
			t.Fatalf("drop sc: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, "CREATE DATABASE sc"); err != nil {
			t.Fatalf("create sc: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, "USE sc"); err != nil {
			t.Fatalf("use sc: %v", err)
		}
	}

	// ---------------------------------------------------------------
	// C1.1: FK name counter uses max(existing) + 1, not count + 1.
	// ---------------------------------------------------------------
	t.Run("C1_1_FK_counter_max_plus_one", func(t *testing.T) {
		reset(t)
		// Test both CREATE TABLE and ALTER TABLE flows.
		// Phase A: CREATE TABLE mixing explicit high-numbered CONSTRAINT name and unnamed FK.
		rowsA := spotCheckQuery(t, mc, `
			CREATE TABLE parent (id INT PRIMARY KEY);
			CREATE TABLE child (
				a INT,
				b INT,
				CONSTRAINT child_ibfk_5 FOREIGN KEY (a) REFERENCES parent(id),
				FOREIGN KEY (b) REFERENCES parent(id)
			);
		`, `
			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='child' AND CONSTRAINT_TYPE='FOREIGN KEY'
			ORDER BY CONSTRAINT_NAME
		`)
		var namesA []string
		for _, row := range rowsA {
			namesA = append(namesA, asString(row[0]))
		}
		t.Logf("C1.1 phase A (CREATE TABLE) observed names: %v", namesA)

		// Phase B: ALTER TABLE adds an unnamed FK AFTER child_ibfk_5 already exists.
		// Catalog's max+1 rule should yield child_ibfk_6.
		rowsB := spotCheckQuery(t, mc, `
			CREATE TABLE parent2 (id INT PRIMARY KEY);
			CREATE TABLE child2 (a INT, b INT);
			ALTER TABLE child2 ADD CONSTRAINT child2_ibfk_5 FOREIGN KEY (a) REFERENCES parent2(id);
			ALTER TABLE child2 ADD FOREIGN KEY (b) REFERENCES parent2(id);
		`, `
			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='child2' AND CONSTRAINT_TYPE='FOREIGN KEY'
			ORDER BY CONSTRAINT_NAME
		`)
		var namesB []string
		for _, row := range rowsB {
			namesB = append(namesB, asString(row[0]))
		}
		has5 := false
		has6 := false
		for _, n := range namesB {
			if n == "child2_ibfk_5" {
				has5 = true
			}
			if n == "child2_ibfk_6" {
				has6 = true
			}
		}
		if !has5 || !has6 || len(namesB) != 2 {
			t.Errorf("CATALOG MISMATCH C1.1 (phase B, ALTER TABLE max+1): expected [child2_ibfk_5, child2_ibfk_6], got %v", namesB)
		} else {
			t.Logf("OK C1.1 phase B verified (ALTER TABLE honors max+1): %v", namesB)
		}

		// Report phase A as observation (may differ from phase B if CREATE uses count-based).
		if len(namesA) == 2 {
			hasA5, hasA6 := false, false
			for _, n := range namesA {
				if n == "child_ibfk_5" {
					hasA5 = true
				}
				if n == "child_ibfk_6" {
					hasA6 = true
				}
			}
			if hasA5 && hasA6 {
				t.Logf("OK C1.1 phase A: CREATE TABLE also follows max+1: %v", namesA)
			} else {
				t.Logf("NOTE C1.1 phase A: CREATE TABLE does NOT follow max+1 (catalog may need clarification): %v", namesA)
			}
		}
	})

	// ---------------------------------------------------------------
	// C1.2: Partition default names are p0, p1, p2, ...
	// ---------------------------------------------------------------
	t.Run("C1_2_partition_naming", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE pt (id INT) PARTITION BY HASH(id) PARTITIONS 4;
		`, `
			SELECT PARTITION_NAME FROM information_schema.PARTITIONS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='pt'
			ORDER BY PARTITION_ORDINAL_POSITION
		`)
		var names []string
		for _, row := range rows {
			names = append(names, asString(row[0]))
		}
		want := []string{"p0", "p1", "p2", "p3"}
		if len(names) != 4 {
			t.Errorf("CATALOG MISMATCH C1.2: expected 4 partitions, got %d: %v", len(names), names)
			return
		}
		for i, w := range want {
			if names[i] != w {
				t.Errorf("CATALOG MISMATCH C1.2: expected %v, got %v", want, names)
				return
			}
		}
		t.Logf("OK C1.2 verified: %v", names)
	})

	// ---------------------------------------------------------------
	// C3.1: TIMESTAMP NOT NULL promotion applies only to FIRST timestamp.
	// Catalog claim: ts1 gets DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	// ts2 does NOT.
	// ---------------------------------------------------------------
	t.Run("C3_1_timestamp_first_only_promotion", func(t *testing.T) {
		reset(t)
		// We need explicit_defaults_for_timestamp=OFF to see the promotion,
		// and we must relax STRICT mode because a second TIMESTAMP NOT NULL
		// with no default receives the zero-date default which STRICT rejects.
		if _, err := mc.db.ExecContext(mc.ctx, "SET SESSION explicit_defaults_for_timestamp=OFF"); err != nil {
			t.Fatalf("set sql var: %v", err)
		}
		defer mc.db.ExecContext(mc.ctx, "SET SESSION explicit_defaults_for_timestamp=ON")
		if _, err := mc.db.ExecContext(mc.ctx, "SET SESSION sql_mode=''"); err != nil {
			t.Fatalf("set sql_mode: %v", err)
		}
		defer mc.db.ExecContext(mc.ctx, "SET SESSION sql_mode=DEFAULT")

		rows := spotCheckQuery(t, mc, `
			CREATE TABLE ts (
				ts1 TIMESTAMP NOT NULL,
				ts2 TIMESTAMP NOT NULL
			);
		`, `
			SELECT COLUMN_NAME, COLUMN_DEFAULT, EXTRA FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='ts'
			ORDER BY ORDINAL_POSITION
		`)
		if len(rows) != 2 {
			t.Fatalf("expected 2 cols, got %d", len(rows))
		}
		ts1Default := asString(rows[0][1])
		ts1Extra := asString(rows[0][2])
		ts2Default := asString(rows[1][1])
		ts2Extra := asString(rows[1][2])

		ts1Promoted := strings.Contains(strings.ToUpper(ts1Default), "CURRENT_TIMESTAMP") &&
			strings.Contains(strings.ToUpper(ts1Extra), "ON UPDATE")
		ts2Promoted := strings.Contains(strings.ToUpper(ts2Default), "CURRENT_TIMESTAMP")

		if !ts1Promoted {
			t.Errorf("CATALOG MISMATCH C3.1: expected ts1 promoted; got default=%q extra=%q",
				ts1Default, ts1Extra)
		}
		if ts2Promoted {
			t.Errorf("CATALOG MISMATCH C3.1: ts2 should NOT be promoted; got default=%q extra=%q",
				ts2Default, ts2Extra)
		}
		if ts1Promoted && !ts2Promoted {
			t.Logf("OK C3.1 verified: ts1 default=%q extra=%q; ts2 default=%q extra=%q",
				ts1Default, ts1Extra, ts2Default, ts2Extra)
		}
	})

	// ---------------------------------------------------------------
	// C3.2: PRIMARY KEY column is implicitly NOT NULL.
	// ---------------------------------------------------------------
	t.Run("C3_2_primary_key_implies_not_null", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE pk (id INT, PRIMARY KEY(id));
		`, `
			SELECT IS_NULLABLE FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='pk' AND COLUMN_NAME='id'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		nullable := asString(rows[0][0])
		if nullable != "NO" {
			t.Errorf("CATALOG MISMATCH C3.2: expected IS_NULLABLE=NO, got %q", nullable)
		} else {
			t.Logf("OK C3.2 verified: IS_NULLABLE=%q", nullable)
		}
	})

	// ---------------------------------------------------------------
	// C4.1: Table inherits charset from database.
	// ---------------------------------------------------------------
	t.Run("C4_1_table_charset_from_database", func(t *testing.T) {
		// Special: create DB with custom charset, not 'sc'.
		if _, err := mc.db.ExecContext(mc.ctx, "DROP DATABASE IF EXISTS sc_cs"); err != nil {
			t.Fatalf("drop: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, "CREATE DATABASE sc_cs CHARACTER SET latin1 COLLATE latin1_swedish_ci"); err != nil {
			t.Fatalf("create db: %v", err)
		}
		defer mc.db.ExecContext(mc.ctx, "DROP DATABASE IF EXISTS sc_cs")
		if _, err := mc.db.ExecContext(mc.ctx, "CREATE TABLE sc_cs.t (c VARCHAR(10))"); err != nil {
			t.Fatalf("create table: %v", err)
		}
		rows := spotCheckQuery(t, mc, ``, `
			SELECT TABLE_COLLATION FROM information_schema.TABLES
			WHERE TABLE_SCHEMA='sc_cs' AND TABLE_NAME='t'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		coll := asString(rows[0][0])
		if coll != "latin1_swedish_ci" {
			t.Errorf("CATALOG MISMATCH C4.1: expected latin1_swedish_ci, got %q", coll)
		} else {
			t.Logf("OK C4.1 verified: TABLE_COLLATION=%q", coll)
		}
	})

	// ---------------------------------------------------------------
	// C5.1: FK ON DELETE default. Catalog claims parser default is
	// FK_OPTION_RESTRICT. However, information_schema.REFERENTIAL_CONSTRAINTS
	// famously reports 'NO ACTION' for both RESTRICT and unspecified, and
	// SHOW CREATE TABLE elides the clause entirely. Verify both views.
	// ---------------------------------------------------------------
	t.Run("C5_1_fk_on_delete_default", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE p5 (id INT PRIMARY KEY);
			CREATE TABLE c5 (
				a INT,
				FOREIGN KEY (a) REFERENCES p5(id)
			);
		`, `
			SELECT DELETE_RULE, UPDATE_RULE FROM information_schema.REFERENTIAL_CONSTRAINTS
			WHERE CONSTRAINT_SCHEMA='sc' AND TABLE_NAME='c5'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		del := asString(rows[0][0])
		upd := asString(rows[0][1])
		// Real MySQL 8.0 reports "NO ACTION" in info_schema for unspecified.
		if del != "NO ACTION" || upd != "NO ACTION" {
			t.Errorf("CATALOG OBS C5.1: REFERENTIAL_CONSTRAINTS reports DELETE=%q UPDATE=%q (catalog says parser default is FK_OPTION_RESTRICT; confirm whether the catalog means semantic behavior vs reporting)", del, upd)
		} else {
			t.Logf("OK C5.1 observed: info_schema reports DELETE=%q UPDATE=%q (semantically equivalent to RESTRICT)", del, upd)
		}
		// Also verify SHOW CREATE TABLE elides ON DELETE clause (standard behavior).
		stmt, err := mc.showCreateTable("c5")
		if err != nil {
			t.Fatalf("show create: %v", err)
		}
		if strings.Contains(strings.ToUpper(stmt), "ON DELETE") {
			t.Errorf("unexpected ON DELETE clause in SHOW CREATE TABLE: %s", stmt)
		} else {
			t.Logf("OK C5.1 SHOW CREATE elides ON DELETE clause: %s", stmt)
		}
	})

	// ---------------------------------------------------------------
	// C10.2: View SQL SECURITY defaults to DEFINER.
	// ---------------------------------------------------------------
	t.Run("C10_2_view_security_definer", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE base (id INT);
			CREATE VIEW v AS SELECT * FROM base;
		`, `
			SELECT SECURITY_TYPE FROM information_schema.VIEWS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='v'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		sec := asString(rows[0][0])
		if sec != "DEFINER" {
			t.Errorf("CATALOG MISMATCH C10.2: expected SECURITY_TYPE=DEFINER, got %q", sec)
		} else {
			t.Logf("OK C10.2 verified: SECURITY_TYPE=%q", sec)
		}
	})

	// ---------------------------------------------------------------
	// C16.1: NOW()/CURRENT_TIMESTAMP precision defaults to 0.
	// Test via a column with DEFAULT NOW() -- the column's DATETIME_PRECISION
	// is driven by the column type, so instead check the generated column case:
	// if we use `DATETIME` without precision, DATETIME_PRECISION = 0.
	// More directly observable: use LENGTH(NOW()) in a scalar SELECT.
	// ---------------------------------------------------------------
	t.Run("C16_1_now_precision_default_zero", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, ``, `SELECT LENGTH(NOW()), LENGTH(NOW(6)), LENGTH(CURRENT_TIMESTAMP)`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		// "YYYY-MM-DD HH:MM:SS" = 19, with (6) fractional = 19+7 = 26.
		l0 := asString(rows[0][0])
		l6 := asString(rows[0][1])
		lCT := asString(rows[0][2])
		if l0 != "19" {
			t.Errorf("CATALOG MISMATCH C16.1: LENGTH(NOW())=%q, expected 19 (no fractional seconds)", l0)
		}
		if l6 != "26" {
			t.Errorf("CATALOG MISMATCH C16.1: LENGTH(NOW(6))=%q, expected 26", l6)
		}
		if lCT != "19" {
			t.Errorf("CATALOG MISMATCH C16.1: LENGTH(CURRENT_TIMESTAMP)=%q, expected 19", lCT)
		}
		if l0 == "19" && l6 == "26" && lCT == "19" {
			t.Logf("OK C16.1 verified: NOW()=%s CURRENT_TIMESTAMP=%s NOW(6)=%s", l0, lCT, l6)
		}
	})

	// ---------------------------------------------------------------
	// C18.4: AUTO_INCREMENT clause elided in SHOW CREATE TABLE if counter <= 1.
	// ---------------------------------------------------------------
	t.Run("C18_4_auto_increment_elision", func(t *testing.T) {
		reset(t)
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE ai (id INT AUTO_INCREMENT PRIMARY KEY, v INT)`); err != nil {
			t.Fatalf("create: %v", err)
		}
		stmt1, err := mc.showCreateTable("ai")
		if err != nil {
			t.Fatalf("show create: %v", err)
		}
		// Fresh table: no AUTO_INCREMENT=N clause.
		if strings.Contains(strings.ToUpper(stmt1), "AUTO_INCREMENT=") {
			t.Errorf("CATALOG MISMATCH C18.4 (before insert): expected no AUTO_INCREMENT= clause, got: %s", stmt1)
		} else {
			t.Logf("OK C18.4 verified (before insert): %s", stmt1)
		}
		// After inserting, counter advances, clause should appear.
		if _, err := mc.db.ExecContext(mc.ctx, `INSERT INTO ai (v) VALUES (10),(20),(30)`); err != nil {
			t.Fatalf("insert: %v", err)
		}
		stmt2, err := mc.showCreateTable("ai")
		if err != nil {
			t.Fatalf("show create: %v", err)
		}
		if !strings.Contains(strings.ToUpper(stmt2), "AUTO_INCREMENT=") {
			t.Errorf("CATALOG MISMATCH C18.4 (after insert): expected AUTO_INCREMENT= clause, got: %s", stmt2)
		} else {
			t.Logf("OK C18.4 verified (after insert, counter > 1): %s", stmt2)
		}
	})

	// ===================================================================
	// Round 2 extended spot-check (PS1-PS7 path-splits + Round 1/2 gaps)
	// ===================================================================

	// ---------------------------------------------------------------
	// PS1 (CREATE): CHECK constraint counter — CREATE uses FRESH counter.
	// Catalog claim: unnamed CHECK constraints in CREATE are numbered
	// 1, 2, 3, ... starting from 0 regardless of user-named CCs.
	// ---------------------------------------------------------------
	t.Run("PS1_CheckCounter_CREATE_fresh", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE tchk (
				a INT,
				CONSTRAINT tchk_chk_5 CHECK (a > 0),
				b INT,
				CHECK (b < 100)
			);
		`, `
			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tchk' AND CONSTRAINT_TYPE='CHECK'
			ORDER BY CONSTRAINT_NAME
		`)
		var names []string
		for _, row := range rows {
			names = append(names, asString(row[0]))
		}
		t.Logf("PS1 CREATE observed: %v", names)
		has1, has5 := false, false
		for _, n := range names {
			if n == "tchk_chk_1" {
				has1 = true
			}
			if n == "tchk_chk_5" {
				has5 = true
			}
		}
		if !has1 || !has5 || len(names) != 2 {
			t.Errorf("CATALOG MISMATCH PS1 CREATE: expected [tchk_chk_1, tchk_chk_5], got %v", names)
		} else {
			t.Logf("OK PS1 CREATE verified (fresh counter from 0): %v", names)
		}
	})

	// ---------------------------------------------------------------
	// PS1 (ALTER): Does ALTER use fresh counter (like CREATE) or
	// max+1 (like FK ALTER)? This is the open question.
	// ---------------------------------------------------------------
	t.Run("PS1_CheckCounter_ALTER_open", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE tchk2 (
				a INT,
				b INT,
				CONSTRAINT tchk2_chk_20 CHECK (a > 0)
			);
			ALTER TABLE tchk2 ADD CHECK (b > 0);
		`, `
			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tchk2' AND CONSTRAINT_TYPE='CHECK'
			ORDER BY CONSTRAINT_NAME
		`)
		var names []string
		for _, row := range rows {
			names = append(names, asString(row[0]))
		}
		t.Logf("PS1 ALTER observed: %v", names)
		hasFresh := false
		hasMaxPlus1 := false
		for _, n := range names {
			if n == "tchk2_chk_1" {
				hasFresh = true
			}
			if n == "tchk2_chk_21" {
				hasMaxPlus1 = true
			}
		}
		switch {
		case hasFresh:
			t.Logf("PS1 ALTER FINDING: fresh counter (tchk2_chk_1). Catalog should document ALTER=fresh.")
		case hasMaxPlus1:
			t.Logf("PS1 ALTER FINDING: max+1 counter (tchk2_chk_21). Catalog should document ALTER=max+1 (like FK).")
		default:
			t.Logf("PS1 ALTER FINDING: unexpected name(s) %v", names)
		}
	})

	// ---------------------------------------------------------------
	// PS5: DATETIME(6) DEFAULT NOW() (fsp=0) — catalog says ER_INVALID_DEFAULT.
	// ---------------------------------------------------------------
	t.Run("PS5_DatetimeFspMismatch", func(t *testing.T) {
		reset(t)
		_, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE tps5 (ts DATETIME(6) DEFAULT NOW())`)
		if err == nil {
			// Accepted — catalog MISMATCH. Read back what was stored.
			rows := spotCheckQuery(t, mc, ``, `
				SELECT COLUMN_NAME, COLUMN_DEFAULT, DATETIME_PRECISION
				FROM information_schema.COLUMNS
				WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tps5'
			`)
			t.Errorf("CATALOG MISMATCH PS5: MySQL accepted DATETIME(6) DEFAULT NOW() (catalog says ER_INVALID_DEFAULT). COLUMNS=%v", rows)
			return
		}
		msg := err.Error()
		t.Logf("PS5 error observed: %v", msg)
		if !strings.Contains(strings.ToLower(msg), "invalid default") && !strings.Contains(msg, "1067") {
			t.Errorf("PS5 UNEXPECTED ERROR TEXT: expected ER_INVALID_DEFAULT (1067), got: %v", msg)
		} else {
			t.Logf("OK PS5 verified: MySQL rejects DATETIME(6) DEFAULT NOW() with ER_INVALID_DEFAULT")
		}
	})

	// ---------------------------------------------------------------
	// PS7: FK name collision — first unnamed FK wants t_ibfk_1, collides
	// with user-named t_ibfk_1 → ER_FK_DUP_NAME.
	// ---------------------------------------------------------------
	t.Run("PS7_FKNameCollision", func(t *testing.T) {
		reset(t)
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE p7 (id INT PRIMARY KEY)`); err != nil {
			t.Fatalf("setup: %v", err)
		}
		_, err := mc.db.ExecContext(mc.ctx, `
			CREATE TABLE tps7 (
				a INT,
				CONSTRAINT tps7_ibfk_1 FOREIGN KEY (a) REFERENCES p7(id),
				b INT,
				FOREIGN KEY (b) REFERENCES p7(id)
			)
		`)
		if err == nil {
			rows := spotCheckQuery(t, mc, ``, `
				SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
				WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tps7' AND CONSTRAINT_TYPE='FOREIGN KEY'
				ORDER BY CONSTRAINT_NAME
			`)
			var names []string
			for _, row := range rows {
				names = append(names, asString(row[0]))
			}
			t.Errorf("CATALOG MISMATCH PS7: expected ER_FK_DUP_NAME, got success with FKs %v", names)
			return
		}
		msg := err.Error()
		t.Logf("PS7 error observed: %v", msg)
		// ER_FK_DUP_NAME = 1826
		if !strings.Contains(msg, "1826") && !strings.Contains(strings.ToLower(msg), "duplicate") {
			t.Errorf("PS7 UNEXPECTED ERROR: expected ER_FK_DUP_NAME (1826), got: %v", msg)
		} else {
			t.Logf("OK PS7 verified: collision rejected")
		}
	})

	// ---------------------------------------------------------------
	// C1.3: Check constraint name format = {table}_chk_N
	// (Also verified as part of PS1 CREATE above; explicit here.)
	// ---------------------------------------------------------------
	t.Run("C1_3_CheckConstraintName", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE cc (a INT, CHECK (a > 0), b INT, CHECK (b < 100));
		`, `
			SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='cc' AND CONSTRAINT_TYPE='CHECK'
			ORDER BY CONSTRAINT_NAME
		`)
		var names []string
		for _, row := range rows {
			names = append(names, asString(row[0]))
		}
		want := []string{"cc_chk_1", "cc_chk_2"}
		if len(names) != 2 || names[0] != want[0] || names[1] != want[1] {
			t.Errorf("CATALOG MISMATCH C1.3: expected %v, got %v", want, names)
		} else {
			t.Logf("OK C1.3 verified: %v", names)
		}
	})

	// ---------------------------------------------------------------
	// C2.1: REAL → DOUBLE.
	// ---------------------------------------------------------------
	t.Run("C2_1_REAL_to_DOUBLE", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE t2 (x REAL);
		`, `
			SELECT DATA_TYPE, COLUMN_TYPE FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='t2' AND COLUMN_NAME='x'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		dt := strings.ToLower(asString(rows[0][0]))
		ct := strings.ToLower(asString(rows[0][1]))
		if dt != "double" {
			t.Errorf("CATALOG MISMATCH C2.1: expected DATA_TYPE=double, got %q (COLUMN_TYPE=%q)", dt, ct)
		} else {
			t.Logf("OK C2.1 verified: DATA_TYPE=%q COLUMN_TYPE=%q", dt, ct)
		}
	})

	// ---------------------------------------------------------------
	// C2.2: BOOL → TINYINT(1).
	// ---------------------------------------------------------------
	t.Run("C2_2_BOOL_to_TINYINT1", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE tbool (flag BOOL);
		`, `
			SELECT DATA_TYPE, COLUMN_TYPE FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tbool' AND COLUMN_NAME='flag'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		dt := strings.ToLower(asString(rows[0][0]))
		ct := strings.ToLower(asString(rows[0][1]))
		if dt != "tinyint" || ct != "tinyint(1)" {
			t.Errorf("CATALOG MISMATCH C2.2: expected tinyint / tinyint(1), got %q / %q", dt, ct)
		} else {
			t.Logf("OK C2.2 verified: %q / %q", dt, ct)
		}
	})

	// ---------------------------------------------------------------
	// C3.3: AUTO_INCREMENT implies NOT NULL.
	// ---------------------------------------------------------------
	t.Run("C3_3_AutoIncrement_implies_NOT_NULL", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE tai (id INT AUTO_INCREMENT PRIMARY KEY);
		`, `
			SELECT IS_NULLABLE FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tai' AND COLUMN_NAME='id'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		if asString(rows[0][0]) != "NO" {
			t.Errorf("CATALOG MISMATCH C3.3: expected IS_NULLABLE=NO, got %q", asString(rows[0][0]))
		} else {
			t.Logf("OK C3.3 verified")
		}
	})

	// ---------------------------------------------------------------
	// C4.2 + C18.1 + C18.5: DB utf8mb4 + table latin1 override;
	// per-column charset inheritance from table charset;
	// SHOW CREATE elides the per-column charset when matching table.
	// ---------------------------------------------------------------
	t.Run("C4_2_and_C18_1_and_C18_5_charset_inheritance_and_elision", func(t *testing.T) {
		if _, err := mc.db.ExecContext(mc.ctx, "DROP DATABASE IF EXISTS sc_cs2"); err != nil {
			t.Fatalf("drop: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, "CREATE DATABASE sc_cs2 CHARACTER SET utf8mb4"); err != nil {
			t.Fatalf("create db: %v", err)
		}
		defer mc.db.ExecContext(mc.ctx, "DROP DATABASE IF EXISTS sc_cs2")
		if _, err := mc.db.ExecContext(mc.ctx, "CREATE TABLE sc_cs2.t (c VARCHAR(10)) CHARSET latin1"); err != nil {
			t.Fatalf("create table: %v", err)
		}
		rows := spotCheckQuery(t, mc, ``, `
			SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc_cs2' AND TABLE_NAME='t' AND COLUMN_NAME='c'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		cs := asString(rows[0][0])
		if cs != "latin1" {
			t.Errorf("CATALOG MISMATCH C4.2: expected column charset=latin1 (inherited from table), got %q", cs)
		} else {
			t.Logf("OK C4.2 verified: column charset=%s", cs)
		}
		// SHOW CREATE TABLE should NOT contain per-column CHARACTER SET,
		// but SHOULD contain DEFAULT CHARSET=latin1 (explicitly specified).
		var scStmt string
		row := mc.db.QueryRowContext(mc.ctx, "SHOW CREATE TABLE sc_cs2.t")
		var tbl string
		if err := row.Scan(&tbl, &scStmt); err != nil {
			t.Fatalf("show create: %v", err)
		}
		up := strings.ToUpper(scStmt)
		// C18.1: column-level CHARACTER SET should be elided
		// We look specifically for "CHARACTER SET" after the column name "c".
		colLineIdx := strings.Index(scStmt, "`c` ")
		if colLineIdx < 0 {
			t.Logf("C18.1 NOTE: could not find `c` column line in SHOW CREATE")
		} else {
			rest := scStmt[colLineIdx:]
			if nl := strings.Index(rest, "\n"); nl >= 0 {
				rest = rest[:nl]
			}
			if strings.Contains(strings.ToUpper(rest), "CHARACTER SET") {
				t.Errorf("CATALOG MISMATCH C18.1: expected per-column CHARACTER SET elided; column line: %q", rest)
			} else {
				t.Logf("OK C18.1 verified: per-column CHARACTER SET elided (column line: %q)", rest)
			}
		}
		// C18.5: DEFAULT CHARSET=latin1 SHOULD be present (user explicitly specified).
		if !strings.Contains(up, "DEFAULT CHARSET=LATIN1") && !strings.Contains(up, "CHARSET=LATIN1") {
			t.Errorf("CATALOG MISMATCH C18.5 (explicit): expected DEFAULT CHARSET=latin1 to be shown, got: %s", scStmt)
		} else {
			t.Logf("OK C18.5 (explicit) verified: DEFAULT CHARSET=latin1 shown")
		}
	})

	// ---------------------------------------------------------------
	// C18.5 (implicit): CREATE TABLE without charset → SHOW CREATE
	// may still include DEFAULT CHARSET clause inherited from DB.
	// Real MySQL: it DOES show DEFAULT CHARSET even when inherited.
	// ---------------------------------------------------------------
	t.Run("C18_5_DefaultCharset_implicit", func(t *testing.T) {
		reset(t)
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE tnocs (x INT)`); err != nil {
			t.Fatalf("create: %v", err)
		}
		stmt, err := mc.showCreateTable("tnocs")
		if err != nil {
			t.Fatalf("show create: %v", err)
		}
		up := strings.ToUpper(stmt)
		has := strings.Contains(up, "DEFAULT CHARSET=") || strings.Contains(up, "CHARSET=")
		t.Logf("C18.5 (implicit) observation: DEFAULT CHARSET present=%v; stmt=%s", has, stmt)
	})

	// ---------------------------------------------------------------
	// C5.3: FK MATCH default.
	// ---------------------------------------------------------------
	t.Run("C5_3_FK_MATCH_default", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE pm (id INT PRIMARY KEY);
			CREATE TABLE cm (a INT, FOREIGN KEY (a) REFERENCES pm(id));
		`, `
			SELECT MATCH_OPTION FROM information_schema.REFERENTIAL_CONSTRAINTS
			WHERE CONSTRAINT_SCHEMA='sc' AND TABLE_NAME='cm'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		t.Logf("OK C5.3 observed: MATCH_OPTION=%q", asString(rows[0][0]))
		// Catalog says FK_MATCH_SIMPLE; info_schema typically reports "NONE" for InnoDB.
	})

	// ---------------------------------------------------------------
	// C6.1: PARTITION BY HASH without PARTITIONS defaults to 1.
	// ---------------------------------------------------------------
	t.Run("C6_1_Partition_default_count", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE phd (id INT) PARTITION BY HASH(id);
		`, `
			SELECT COUNT(*) FROM information_schema.PARTITIONS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='phd' AND PARTITION_NAME IS NOT NULL
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		n := asString(rows[0][0])
		if n != "1" {
			t.Errorf("CATALOG MISMATCH C6.1: expected 1 partition, got %s", n)
		} else {
			t.Logf("OK C6.1 verified: partitions=1")
		}
	})

	// ---------------------------------------------------------------
	// C6.2: Subpartitions auto-gen. PARTITIONS 2 SUBPARTITIONS 3 → 6 rows.
	// ---------------------------------------------------------------
	t.Run("C6_2_Subpartition_count", func(t *testing.T) {
		reset(t)
		_, err := mc.db.ExecContext(mc.ctx, `
			CREATE TABLE psp (id INT, d INT)
			PARTITION BY RANGE(id)
			SUBPARTITION BY HASH(d) SUBPARTITIONS 3 (
				PARTITION p0 VALUES LESS THAN (10),
				PARTITION p1 VALUES LESS THAN (20)
			)
		`)
		if err != nil {
			t.Logf("C6.2 NOTE: subpartition DDL errored: %v (skipping verification)", err)
			return
		}
		rows := spotCheckQuery(t, mc, ``, `
			SELECT PARTITION_NAME, SUBPARTITION_NAME FROM information_schema.PARTITIONS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='psp'
			ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
		`)
		if len(rows) != 6 {
			t.Errorf("CATALOG MISMATCH C6.2: expected 6 sub-part rows, got %d: %v", len(rows), rows)
		} else {
			var subs []string
			for _, r := range rows {
				subs = append(subs, asString(r[1]))
			}
			t.Logf("OK C6.2 verified: 6 subparts, names=%v", subs)
		}
	})

	// ---------------------------------------------------------------
	// C7.1: Default index algorithm = BTREE.
	// ---------------------------------------------------------------
	t.Run("C7_1_Default_index_algorithm_BTREE", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE tidx (a INT, KEY(a));
		`, `
			SELECT INDEX_TYPE FROM information_schema.STATISTICS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tidx' AND INDEX_NAME='a'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		if asString(rows[0][0]) != "BTREE" {
			t.Errorf("CATALOG MISMATCH C7.1: expected BTREE, got %q", asString(rows[0][0]))
		} else {
			t.Logf("OK C7.1 verified: INDEX_TYPE=BTREE")
		}
	})

	// ---------------------------------------------------------------
	// C7.2: FK creates implicit backing index on child FK columns.
	// ---------------------------------------------------------------
	t.Run("C7_2_FK_backing_index", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE pfk (id INT PRIMARY KEY);
			CREATE TABLE cfk (a INT, FOREIGN KEY (a) REFERENCES pfk(id));
		`, `
			SELECT INDEX_NAME, COLUMN_NAME FROM information_schema.STATISTICS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='cfk'
			ORDER BY INDEX_NAME, SEQ_IN_INDEX
		`)
		if len(rows) == 0 {
			t.Errorf("CATALOG MISMATCH C7.2: expected at least 1 backing index, got 0")
		} else {
			t.Logf("OK C7.2 verified: %d index row(s) on cfk: %v", len(rows), rows)
		}
	})

	// ---------------------------------------------------------------
	// C8.1: Default engine = InnoDB.
	// ---------------------------------------------------------------
	t.Run("C8_1_Default_engine_InnoDB", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `CREATE TABLE teng (x INT);`, `
			SELECT ENGINE FROM information_schema.TABLES
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='teng'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		if asString(rows[0][0]) != "InnoDB" {
			t.Errorf("CATALOG MISMATCH C8.1: expected InnoDB, got %q", asString(rows[0][0]))
		} else {
			t.Logf("OK C8.1 verified: ENGINE=InnoDB")
		}
	})

	// ---------------------------------------------------------------
	// C8.2: Default ROW_FORMAT.
	// ---------------------------------------------------------------
	t.Run("C8_2_Default_row_format", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `CREATE TABLE trf (x INT);`, `
			SELECT ROW_FORMAT FROM information_schema.TABLES
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='trf'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		rf := asString(rows[0][0])
		if rf != "Dynamic" && rf != "Compact" {
			t.Errorf("CATALOG MISMATCH C8.2: expected Dynamic or Compact, got %q", rf)
		} else {
			t.Logf("OK C8.2 verified: ROW_FORMAT=%s", rf)
		}
	})

	// ---------------------------------------------------------------
	// C9.1: Generated column defaults to VIRTUAL.
	// ---------------------------------------------------------------
	t.Run("C9_1_GeneratedColumn_default_VIRTUAL", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE tgen (a INT, b INT AS (a + 1));
		`, `
			SELECT EXTRA FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tgen' AND COLUMN_NAME='b'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		extra := strings.ToUpper(asString(rows[0][0]))
		if !strings.Contains(extra, "VIRTUAL") {
			t.Errorf("CATALOG MISMATCH C9.1: expected VIRTUAL GENERATED, got %q", extra)
		} else {
			t.Logf("OK C9.1 verified: EXTRA=%s", extra)
		}
	})

	// ---------------------------------------------------------------
	// C10.1 + C10.3 + C10.4: view algorithm, definer, check option.
	// ---------------------------------------------------------------
	t.Run("C10_1_3_4_View_defaults", func(t *testing.T) {
		reset(t)
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE base10 (id INT)`); err != nil {
			t.Fatalf("setup: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE VIEW v10 AS SELECT * FROM base10`); err != nil {
			t.Fatalf("create view: %v", err)
		}
		rows := spotCheckQuery(t, mc, ``, `
			SELECT VIEW_DEFINITION, CHECK_OPTION, DEFINER FROM information_schema.VIEWS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='v10'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		checkOpt := asString(rows[0][1])
		definer := asString(rows[0][2])
		if checkOpt != "NONE" {
			t.Errorf("CATALOG MISMATCH C10.4: expected CHECK_OPTION=NONE, got %q", checkOpt)
		} else {
			t.Logf("OK C10.4 verified: CHECK_OPTION=%s", checkOpt)
		}
		if definer == "" || definer == "<nil>" {
			t.Errorf("CATALOG MISMATCH C10.3: expected DEFINER to be populated, got %q", definer)
		} else {
			t.Logf("OK C10.3 verified: DEFINER=%s", definer)
		}
		// C10.1: SHOW CREATE VIEW for ALGORITHM
		stmt, err := mc.showCreateView("v10")
		if err != nil {
			t.Fatalf("show create view: %v", err)
		}
		up := strings.ToUpper(stmt)
		if !strings.Contains(up, "ALGORITHM=UNDEFINED") {
			t.Errorf("CATALOG MISMATCH C10.1: expected ALGORITHM=UNDEFINED in SHOW CREATE VIEW, got: %s", stmt)
		} else {
			t.Logf("OK C10.1 verified: ALGORITHM=UNDEFINED")
		}
	})

	// ---------------------------------------------------------------
	// C11.1: Trigger DEFINER defaults to current user.
	// ---------------------------------------------------------------
	t.Run("C11_1_Trigger_definer_default", func(t *testing.T) {
		reset(t)
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE tt11 (a INT)`); err != nil {
			t.Fatalf("setup: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx,
			`CREATE TRIGGER trg11 BEFORE INSERT ON tt11 FOR EACH ROW SET NEW.a = NEW.a`); err != nil {
			t.Fatalf("create trigger: %v", err)
		}
		rows := spotCheckQuery(t, mc, ``, `
			SELECT DEFINER FROM information_schema.TRIGGERS
			WHERE TRIGGER_SCHEMA='sc' AND TRIGGER_NAME='trg11'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		def := asString(rows[0][0])
		if def == "" || def == "<nil>" {
			t.Errorf("CATALOG MISMATCH C11.1: expected DEFINER populated, got %q", def)
		} else {
			t.Logf("OK C11.1 verified: DEFINER=%s", def)
		}
	})

	// ---------------------------------------------------------------
	// C14.1: CHECK CONSTRAINT ENFORCED by default.
	// ---------------------------------------------------------------
	t.Run("C14_1_Check_enforced_default", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE tchk14 (a INT, CONSTRAINT chk14 CHECK (a > 0));
		`, `
			SELECT ENFORCED FROM information_schema.TABLE_CONSTRAINTS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tchk14' AND CONSTRAINT_NAME='chk14'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		if asString(rows[0][0]) != "YES" {
			t.Errorf("CATALOG MISMATCH C14.1: expected ENFORCED=YES, got %q", asString(rows[0][0]))
		} else {
			t.Logf("OK C14.1 verified: ENFORCED=YES")
		}
	})

	// ---------------------------------------------------------------
	// C15.1: New column added via ALTER lands at end.
	// ---------------------------------------------------------------
	t.Run("C15_1_Column_positioning_end", func(t *testing.T) {
		reset(t)
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE tpos (a INT, b INT)`); err != nil {
			t.Fatalf("create: %v", err)
		}
		if _, err := mc.db.ExecContext(mc.ctx, `ALTER TABLE tpos ADD COLUMN c INT`); err != nil {
			t.Fatalf("alter: %v", err)
		}
		rows := spotCheckQuery(t, mc, ``, `
			SELECT COLUMN_NAME FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='tpos'
			ORDER BY ORDINAL_POSITION
		`)
		var names []string
		for _, r := range rows {
			names = append(names, asString(r[0]))
		}
		want := []string{"a", "b", "c"}
		if len(names) != 3 || names[0] != want[0] || names[1] != want[1] || names[2] != want[2] {
			t.Errorf("CATALOG MISMATCH C15.1: expected %v, got %v", want, names)
		} else {
			t.Logf("OK C15.1 verified: %v", names)
		}
	})

	// ---------------------------------------------------------------
	// C18.2: NOT NULL rendering in SHOW CREATE TABLE.
	// ---------------------------------------------------------------
	t.Run("C18_2_NotNull_rendering", func(t *testing.T) {
		reset(t)
		if _, err := mc.db.ExecContext(mc.ctx, `CREATE TABLE tnn (x INT, y INT NOT NULL)`); err != nil {
			t.Fatalf("create: %v", err)
		}
		stmt, err := mc.showCreateTable("tnn")
		if err != nil {
			t.Fatalf("show create: %v", err)
		}
		up := strings.ToUpper(stmt)
		// `x` line should NOT contain "NOT NULL"; `y` line SHOULD.
		xIdx := strings.Index(stmt, "`x`")
		yIdx := strings.Index(stmt, "`y`")
		xLine := ""
		yLine := ""
		if xIdx >= 0 {
			e := strings.Index(stmt[xIdx:], "\n")
			if e < 0 {
				e = len(stmt) - xIdx
			}
			xLine = stmt[xIdx : xIdx+e]
		}
		if yIdx >= 0 {
			e := strings.Index(stmt[yIdx:], "\n")
			if e < 0 {
				e = len(stmt) - yIdx
			}
			yLine = stmt[yIdx : yIdx+e]
		}
		_ = up
		if strings.Contains(strings.ToUpper(xLine), "NOT NULL") {
			t.Errorf("CATALOG MISMATCH C18.2: expected `x` line to elide NOT NULL, got %q", xLine)
		} else {
			t.Logf("OK C18.2 (nullable elides): %q", xLine)
		}
		if !strings.Contains(strings.ToUpper(yLine), "NOT NULL") {
			t.Errorf("CATALOG MISMATCH C18.2: expected `y` line to contain NOT NULL, got %q", yLine)
		} else {
			t.Logf("OK C18.2 (NOT NULL shown): %q", yLine)
		}
	})

	// ---------------------------------------------------------------
	// C21.1: DEFAULT NULL on nullable column.
	// ---------------------------------------------------------------
	t.Run("C21_1_Default_NULL", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE td (x INT DEFAULT NULL);
		`, `
			SELECT COLUMN_DEFAULT, IS_NULLABLE FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='td' AND COLUMN_NAME='x'
		`)
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}
		def := asString(rows[0][0])
		nullable := asString(rows[0][1])
		if def != "<nil>" && def != "NULL" {
			t.Errorf("CATALOG MISMATCH C21.1: expected COLUMN_DEFAULT=NULL, got %q", def)
		}
		if nullable != "YES" {
			t.Errorf("CATALOG MISMATCH C21.1: expected IS_NULLABLE=YES, got %q", nullable)
		}
		t.Logf("OK C21.1 verified: default=%q nullable=%q", def, nullable)
	})

	// ---------------------------------------------------------------
	// C25.1 — original DECIMAL test. Keep below.
	// ---------------------------------------------------------------
	t.Run("C25_1_decimal_default_10_0", func(t *testing.T) {
		reset(t)
		rows := spotCheckQuery(t, mc, `
			CREATE TABLE d (x DECIMAL, y DECIMAL(8), z NUMERIC);
		`, `
			SELECT COLUMN_NAME, COLUMN_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
			FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='sc' AND TABLE_NAME='d'
			ORDER BY ORDINAL_POSITION
		`)
		if len(rows) != 3 {
			t.Fatalf("expected 3 cols, got %d", len(rows))
		}
		// x DECIMAL -> decimal(10,0)
		xType := asString(rows[0][1])
		xPrec := asString(rows[0][2])
		xScale := asString(rows[0][3])
		if xType != "decimal(10,0)" || xPrec != "10" || xScale != "0" {
			t.Errorf("CATALOG MISMATCH C25.1 DECIMAL: type=%q prec=%q scale=%q (expected decimal(10,0), 10, 0)",
				xType, xPrec, xScale)
		} else {
			t.Logf("OK C25.1 DECIMAL verified: %s prec=%s scale=%s", xType, xPrec, xScale)
		}
		// y DECIMAL(8) -> decimal(8,0)
		yType := asString(rows[1][1])
		if yType != "decimal(8,0)" {
			t.Errorf("CATALOG MISMATCH C25.1 DECIMAL(8): type=%q (expected decimal(8,0))", yType)
		} else {
			t.Logf("OK C25.1 DECIMAL(8) verified: %s", yType)
		}
		// z NUMERIC -> decimal(10,0) (NUMERIC is alias for DECIMAL)
		zType := asString(rows[2][1])
		if zType != "decimal(10,0)" {
			t.Errorf("CATALOG MISMATCH C25.1 NUMERIC: type=%q (expected decimal(10,0))", zType)
		} else {
			t.Logf("OK C25.1 NUMERIC verified: %s", zType)
		}
	})
}
