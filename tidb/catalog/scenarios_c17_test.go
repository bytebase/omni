package catalog

import (
	"errors"
	"strings"
	"testing"
)

// TestScenario_C17 covers section C17 (string function charset / collation
// propagation) from mysql/catalog/SCENARIOS-mysql-implicit-behavior.md, asked
// of TiDB v8.5.5 and of the omni catalog.
//
// TiDB's answers are not MySQL's (see view_collation.go): it has no
// latin1_swedish_ci or utf8mb4_0900_as_cs under the new collation framework,
// so the tables use collations it has; two IMPLICIT collations of one charset
// aggregate to the charset's _bin collation with coercibility NONE instead of
// failing, and only a comparison refuses NONE; a _bin collation wins over
// another collation at the same coercibility; and LPAD/RPAD aggregate their
// first and third arguments rather than taking the first one's.
//
// TiDB's information_schema.COLUMNS reports NULL CHARACTER_SET_NAME and
// COLLATION_NAME for every view column, so each scenario reads the column's
// charset and collation with SELECT CHARSET(c), COLLATION(c) over the view and
// holds the catalog's View.ColumnMetadata to it. A body TiDB refuses must be
// refused by the catalog with the same error.
func TestScenario_C17(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// A string literal takes the session's collation_connection. The Go
	// driver sets utf8mb4_general_ci; the catalog has no session and models
	// TiDB's server default, utf8mb4_bin, which a client that names no
	// collation gets. The pool holds one connection (scenarioContainer), so
	// the setting holds for every statement below and is put back after.
	var prevCollation string
	oracleScan(t, mc, `SELECT @@collation_connection`, &prevCollation)
	if _, err := mc.db.ExecContext(mc.ctx, `SET NAMES utf8mb4 COLLATE utf8mb4_bin`); err != nil {
		t.Fatalf("tidb SET NAMES: %v", err)
	}
	defer func() {
		if _, err := mc.db.ExecContext(mc.ctx, `SET NAMES utf8mb4 COLLATE `+prevCollation); err != nil {
			t.Errorf("tidb restore collation_connection %s: %v", prevCollation, err)
		}
	}()

	// c17ViewColumn asserts what TiDB and the catalog derive for a view
	// column: SELECT CHARSET(col), COLLATION(col) over the view, which needs a
	// row in it, against View.ColumnMetadata. information_schema.COLUMNS has
	// NULL for both on TiDB, which is asserted too.
	c17ViewColumn := func(t *testing.T, c *Catalog, view, col, wantCharset, wantCollation string) {
		t.Helper()
		var cs, co any
		if err := mc.db.QueryRowContext(mc.ctx,
			`SELECT CHARSET(`+col+`), COLLATION(`+col+`) FROM testdb.`+view+` LIMIT 1`).Scan(&cs, &co); err != nil {
			t.Errorf("tidb CHARSET/COLLATION of %s.%s: %v", view, col, err)
		} else {
			assertStringEq(t, "tidb CHARSET("+view+"."+col+")", asString(cs), wantCharset)
			assertStringEq(t, "tidb COLLATION("+view+"."+col+")", asString(co), wantCollation)
		}

		var isCS, isCo any
		if err := mc.db.QueryRowContext(mc.ctx,
			`SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
			  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?`, view, col).Scan(&isCS, &isCo); err != nil {
			t.Errorf("tidb information_schema.COLUMNS %s.%s: %v", view, col, err)
		} else if isCS != nil || isCo != nil {
			t.Errorf("tidb information_schema.COLUMNS %s.%s: charset %v, collation %v; TiDB v8.5.5 reports NULL for both", view, col, isCS, isCo)
		}

		db := c.GetDatabase("testdb")
		if db == nil {
			t.Errorf("omni: database testdb not found")
			return
		}
		v := db.Views[toLower(view)]
		if v == nil {
			t.Errorf("omni: view %s not created", view)
			return
		}
		for _, meta := range v.ColumnMetadata {
			if strings.EqualFold(meta.Name, col) {
				assertStringEq(t, "omni "+view+"."+col+" charset", meta.Charset, wantCharset)
				assertStringEq(t, "omni "+view+"."+col+" collation", meta.Collation, wantCollation)
				return
			}
		}
		t.Errorf("omni: view column %s.%s has no metadata", view, col)
	}

	// c17Coercibility asserts TiDB's COERCIBILITY(col) over a view, which the
	// catalog derives but does not store.
	c17Coercibility := func(t *testing.T, view, col string, want int) {
		t.Helper()
		var got int
		oracleScan(t, mc, `SELECT COERCIBILITY(`+col+`) FROM testdb.`+view+` LIMIT 1`, &got)
		assertIntEq(t, "tidb COERCIBILITY("+view+"."+col+")", got, want)
	}

	// c17Refused runs a statement that TiDB refuses on both and asserts the
	// same code and message from each.
	c17Refused := func(t *testing.T, c *Catalog, ddl string, wantCode int, wantMsg string) {
		t.Helper()
		_, tidbErr := mc.db.ExecContext(mc.ctx, ddl)
		if tidbErr == nil {
			t.Errorf("tidb accepted %q; want error %d", ddl, wantCode)
		} else if !strings.Contains(tidbErr.Error(), wantMsg) {
			t.Errorf("tidb %q: got %v, want %d %q", ddl, tidbErr, wantCode, wantMsg)
		}

		var omniErr error
		results, err := c.Exec(ddl, nil)
		if err != nil {
			omniErr = err
		}
		for _, r := range results {
			if r.Error != nil && omniErr == nil {
				omniErr = r.Error
			}
		}
		var catErr *Error
		switch {
		case omniErr == nil:
			t.Errorf("omni accepted %q; want error %d", ddl, wantCode)
		case !errors.As(omniErr, &catErr):
			t.Errorf("omni %q: got %v, want *Error %d", ddl, omniErr, wantCode)
		default:
			assertIntEq(t, "omni error code", catErr.Code, wantCode)
			assertStringEq(t, "omni error message", catErr.Message, wantMsg)
		}
	}

	c17Insert := func(t *testing.T, sql string) {
		t.Helper()
		if _, err := mc.db.ExecContext(mc.ctx, sql); err != nil {
			t.Errorf("tidb %q: %v", sql, err)
		}
	}

	// ---- 17.1 CONCAT identical charset/collation --------------------------
	t.Run("17_1_CONCAT_same_charset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v1 AS SELECT CONCAT(a, b) AS c FROM t`)
		c17Insert(t, `INSERT INTO testdb.t VALUES ('x', 'y')`)

		c17ViewColumn(t, c, "v1", "c", "utf8mb4", "utf8mb4_general_ci")
		c17Coercibility(t, "v1", "c", 2)

		var maxLen int
		oracleScan(t, mc, `SELECT CHARACTER_MAXIMUM_LENGTH FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v1' AND COLUMN_NAME='c'`, &maxLen)
		assertIntEq(t, "tidb v1.c CHARACTER_MAXIMUM_LENGTH", maxLen, 20)
	})

	// ---- 17.2 CONCAT mixing latin1 + utf8mb4 ------------------------------
	// At one coercibility the Unicode charset wins, in either order.
	t.Run("17_2_CONCAT_latin1_utf8mb4_superset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			b VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_bin
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v2 AS SELECT CONCAT(a, b) AS c, CONCAT(b, a) AS d FROM t`)
		c17Insert(t, `INSERT INTO testdb.t VALUES ('x', 'y')`)

		c17ViewColumn(t, c, "v2", "c", "utf8mb4", "utf8mb4_general_ci")
		c17ViewColumn(t, c, "v2", "d", "utf8mb4", "utf8mb4_general_ci")
	})

	// ---- 17.3 Two IMPLICIT collations of one charset ----------------------
	// TiDB takes CONCAT of them, giving utf8mb4_bin with coercibility NONE,
	// and refuses only a comparison: of the two columns, or of that NONE
	// result. utf8mb4_bin beats another IMPLICIT collation instead.
	t.Run("17_3_CONCAT_incompatible_collations", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
			e VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v3 AS SELECT CONCAT(a, b) AS c, CONCAT(a, e) AS d FROM t WHERE a = e`)
		c17Insert(t, `INSERT INTO testdb.t VALUES ('x', 'y', 'x')`)

		c17ViewColumn(t, c, "v3", "c", "utf8mb4", "utf8mb4_bin")
		c17Coercibility(t, "v3", "c", 1)
		c17ViewColumn(t, c, "v3", "d", "utf8mb4", "utf8mb4_bin")
		c17Coercibility(t, "v3", "d", 2)

		c17Refused(t, c, `CREATE VIEW v3a AS SELECT a FROM t WHERE a = b`, 1267,
			"Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '='")
		c17Refused(t, c, `CREATE VIEW v3b AS SELECT CONCAT(a, b) = 'x' AS c FROM t`, 1267,
			"Illegal mix of collations (utf8mb4_bin,NONE) and (utf8mb4_bin,COERCIBLE) for operation '='")
		c17Refused(t, c, `CREATE VIEW v3c AS SELECT a FROM t WHERE a LIKE b`, 1267,
			"Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation 'like'")
	})

	// ---- 17.4 CONCAT_WS NULL skipping + separator aggregation -------------
	t.Run("17_4_CONCAT_WS_nullskip", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v4 AS SELECT CONCAT_WS(',', a, b, NULL) AS c FROM t`)

		// Runtime: CONCAT(NULL,'x') is NULL; CONCAT_WS(',',NULL,'x') is 'x'.
		c17Insert(t, `INSERT INTO testdb.t VALUES (NULL, 'x')`)
		var concatRes, cwsRes any
		row := mc.db.QueryRowContext(mc.ctx,
			`SELECT CONCAT(a,b), CONCAT_WS(',',a,b) FROM testdb.t LIMIT 1`)
		if err := row.Scan(&concatRes, &cwsRes); err != nil {
			t.Errorf("tidb runtime scan: %v", err)
		} else {
			if concatRes != nil {
				t.Errorf("tidb CONCAT(NULL,'x'): got %v, want NULL", concatRes)
			}
			if s := asString(cwsRes); s != "x" {
				t.Errorf("tidb CONCAT_WS(',',NULL,'x'): got %q, want \"x\"", s)
			}
		}

		c17ViewColumn(t, c, "v4", "c", "utf8mb4", "utf8mb4_general_ci")
	})

	// ---- 17.5 A literal, introducer or not, is COERCIBLE ------------------
	// The column's charset wins over the literal's. An ascii column cannot
	// take a non-ASCII literal, which TiDB converts by value.
	t.Run("17_5_introducer_still_coercible", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET latin1 COLLATE latin1_bin,
			s VARCHAR(10) CHARACTER SET ascii
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v5a AS SELECT CONCAT(a, 'x') AS c FROM t`)
		runOnBoth(t, mc, c, `CREATE VIEW v5b AS SELECT CONCAT(a, _utf8mb4'x') AS c FROM t`)
		runOnBoth(t, mc, c, `CREATE VIEW v5c AS SELECT CONCAT('x', s) AS c FROM t`)
		c17Insert(t, `INSERT INTO testdb.t VALUES ('x', 'y')`)

		c17ViewColumn(t, c, "v5a", "c", "latin1", "latin1_bin")
		c17ViewColumn(t, c, "v5b", "c", "latin1", "latin1_bin")
		c17ViewColumn(t, c, "v5c", "c", "ascii", "ascii_bin")

		c17Refused(t, c, `CREATE VIEW v5d AS SELECT CONCAT(s, 'é') AS c FROM t`, 1267,
			"Illegal mix of collations (ascii_bin,IMPLICIT) and (utf8mb4_bin,COERCIBLE) for operation 'concat'")
	})

	// ---- 17.6 REPEAT pins to its string; LPAD / RPAD aggregate the pad -----
	// TiDB's deriveCollation aggregates LPAD and RPAD's first and third
	// arguments, so the utf8mb4 pad wins over a latin1 string.
	t.Run("17_6_REPEAT_LPAD_RPAD_charset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_bin,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			u VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
		)`)
		runOnBoth(t, mc, c, `CREATE VIEW v6a AS SELECT REPEAT(a, 3) AS c FROM t`)
		runOnBoth(t, mc, c, `CREATE VIEW v6b AS SELECT LPAD(a, 20, b) AS c FROM t`)
		runOnBoth(t, mc, c, `CREATE VIEW v6c AS SELECT RPAD(b, 20, a) AS c FROM t`)
		runOnBoth(t, mc, c, `CREATE VIEW v6d AS SELECT LPAD(b, 20, u) AS c FROM t`)
		c17Insert(t, `INSERT INTO testdb.t VALUES ('x', 'y', 'z')`)

		c17ViewColumn(t, c, "v6a", "c", "latin1", "latin1_bin")
		c17ViewColumn(t, c, "v6b", "c", "utf8mb4", "utf8mb4_general_ci")
		c17ViewColumn(t, c, "v6c", "c", "utf8mb4", "utf8mb4_general_ci")
		c17ViewColumn(t, c, "v6d", "c", "utf8mb4", "utf8mb4_bin")
		c17Coercibility(t, "v6d", "c", 1)

		c17Refused(t, c, `CREATE VIEW v6e AS SELECT a FROM t WHERE LPAD(b, 20, u) = 'x'`, 1267,
			"Illegal mix of collations (utf8mb4_bin,NONE) and (utf8mb4_bin,COERCIBLE) for operation '='")
	})

	// ---- 17.7 CONVERT(x USING cs) is IMPLICIT in cs's default collation ---
	t.Run("17_7_CONVERT_USING_pins_charset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
		)`)
		// a = b is 17.3's 1267; converting b gives utf8mb4_bin, which wins.
		runOnBoth(t, mc, c, `CREATE VIEW v7 AS
			SELECT CONCAT(a, CONVERT(b USING utf8mb4)) AS c, CONVERT(a USING latin1) AS d
			FROM t WHERE a = CONVERT(b USING utf8mb4)`)
		c17Insert(t, `INSERT INTO testdb.t VALUES ('x', 'x')`)

		c17ViewColumn(t, c, "v7", "c", "utf8mb4", "utf8mb4_bin")
		c17ViewColumn(t, c, "v7", "d", "latin1", "latin1_bin")
		c17Coercibility(t, "v7", "d", 2)
	})

	// ---- 17.8 COLLATE clause is EXPLICIT — highest precedence -------------
	t.Run("17_8_COLLATE_explicit", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t (
			a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
			l VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_bin
		)`)

		// One EXPLICIT side wins.
		runOnBoth(t, mc, c, `CREATE VIEW v8a AS SELECT CONCAT(a, b COLLATE utf8mb4_bin) AS c FROM t`)
		c17Insert(t, `INSERT INTO testdb.t VALUES ('x', 'y', 'z')`)
		c17ViewColumn(t, c, "v8a", "c", "utf8mb4", "utf8mb4_bin")
		c17Coercibility(t, "v8a", "c", 0)

		// Two EXPLICIT sides with different collations fail, and so does an
		// EXPLICIT latin1 that a utf8mb4 column cannot convert to.
		c17Refused(t, c, `CREATE VIEW v8b AS SELECT CONCAT(a COLLATE utf8mb4_general_ci, b COLLATE utf8mb4_bin) AS c FROM t`, 1267,
			"Illegal mix of collations (utf8mb4_general_ci,EXPLICIT) and (utf8mb4_bin,EXPLICIT) for operation 'concat'")
		c17Refused(t, c, `CREATE VIEW v8c AS SELECT CONCAT(l COLLATE latin1_bin, a) AS c FROM t`, 1267,
			"Illegal mix of collations (latin1_bin,EXPLICIT) and (utf8mb4_general_ci,IMPLICIT) for operation 'concat'")
		c17Refused(t, c, `CREATE VIEW v8d AS SELECT LENGTH(CONCAT(a COLLATE utf8mb4_general_ci, b COLLATE utf8mb4_bin)) AS c FROM t`, 1267,
			"Illegal mix of collations (utf8mb4_general_ci,EXPLICIT) and (utf8mb4_bin,EXPLICIT) for operation 'concat'")
	})
}
