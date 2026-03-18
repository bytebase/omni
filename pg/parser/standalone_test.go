package parser

import "testing"

// parseOK is a test helper that verifies SQL parses without error.
func parseOK(t *testing.T, sql string) {
	t.Helper()
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
}

func TestSetCatalog(t *testing.T) {
	tests := []string{
		"SET CATALOG 'test'",
		"SET CATALOG 'mydb'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
			CheckLocations(t, sql)
		})
	}
}

func TestDropOperatorUnary(t *testing.T) {
	tests := []string{
		"DROP OPERATOR + (integer)",
		"DROP OPERATOR + (NONE, integer)",
		"DROP OPERATOR + (integer, NONE)",
		"DROP OPERATOR IF EXISTS + (integer)",
		"DROP OPERATOR + (integer) CASCADE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
			CheckLocations(t, sql)
		})
	}
}

func TestAlterOperatorUnary(t *testing.T) {
	tests := []string{
		"ALTER EXTENSION myext ADD OPERATOR + (integer)",
		"ALTER EXTENSION myext DROP OPERATOR + (integer)",
		"ALTER EXTENSION myext ADD OPERATOR + (NONE, integer)",
		"ALTER EXTENSION myext ADD OPERATOR + (integer, NONE)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
			CheckLocations(t, sql)
		})
	}
}

func TestDdlCompatUnloggedView(t *testing.T) {
	tests := []string{
		"CREATE UNLOGGED VIEW v AS SELECT 1",
		"CREATE UNLOGGED VIEW myschema.v AS SELECT a, b FROM t",
		"CREATE UNLOGGED VIEW v (x, y) AS SELECT 1, 2",
		"CREATE UNLOGGED VIEW v AS SELECT * FROM t WITH CHECK OPTION",
		"CREATE UNLOGGED VIEW v AS SELECT * FROM t WITH LOCAL CHECK OPTION",
		"CREATE UNLOGGED VIEW v AS SELECT * FROM t WITH CASCADED CHECK OPTION",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
			CheckLocations(t, sql)
		})
	}
}

func TestDdlCompatModifyingCTE(t *testing.T) {
	tests := []string{
		"WITH moved AS (DELETE FROM t1 WHERE id > 100 RETURNING *) INSERT INTO t2 SELECT * FROM moved",
		"WITH upd AS (UPDATE t SET x = 1 RETURNING *) SELECT * FROM upd",
		"WITH ins AS (INSERT INTO t VALUES (1) RETURNING *) SELECT * FROM ins",
		"WITH del AS (DELETE FROM t RETURNING id) SELECT * FROM del",
		"WITH moved AS (DELETE FROM t1 RETURNING *) INSERT INTO t2 SELECT * FROM moved",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
			CheckLocations(t, sql)
		})
	}
}

// TestAlterGroupRename runs ALTER GROUP ... RENAME TO tests for batch 77.
func TestAlterGroupRename(t *testing.T) {
	tests := []string{
		"ALTER GROUP mygroup RENAME TO newgroup",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
			CheckLocations(t, sql)
		})
	}
}

// TestAlterColumnReset runs ALTER COLUMN ... RESET tests for batch 73.
func TestAlterColumnReset(t *testing.T) {
	tests := []string{
		"ALTER TABLE t ALTER COLUMN c RESET (n_distinct)",
		"ALTER TABLE t ALTER COLUMN c RESET (n_distinct, n_distinct_inherited)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
			CheckLocations(t, sql)
		})
	}
}

// TestReservedKeywordInFromClause tests that reserved keywords like USER,
// CURRENT_USER, etc. are accepted in FROM context via the func_table path,
// matching PostgreSQL behavior where they resolve to func_expr_common_subexpr.
func TestReservedKeywordInFromClause(t *testing.T) {
	tests := []string{
		"SELECT * FROM user",
		"SELECT * FROM current_user",
		"SELECT * FROM session_user",
		"SELECT * FROM current_date",
		"SELECT * FROM current_timestamp",
		"SELECT * FROM current_role",
		"SELECT * FROM current_catalog",
		"SELECT * FROM current_schema",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}

// TestSetRoleEquals tests that SET role = 'admin' is accepted as generic_set,
// while SET ROLE admin still works as the keyword-specific form.
func TestSetRoleEquals(t *testing.T) {
	tests := []string{
		"SET role = 'admin'",
		"SET role TO 'admin'",
		"SET ROLE admin",
		"SET ROLE 'admin'",
		"SET ROLE DEFAULT",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseOK(t, sql)
		})
	}
}
