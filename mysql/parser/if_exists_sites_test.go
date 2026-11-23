package parser

import "testing"

// TestIfExistsNeedsExists pins that IF is followed by EXISTS or NOT EXISTS at
// every if_exists and opt_if_not_exists site of sql_yacc.yy omni read
// leniently. mysql:8.0.46's verdicts are the <alternative>+if_exists and
// <alternative>+opt_if_not_exists rows of testdata/compat.tsv.
func TestIfExistsNeedsExists(t *testing.T) {
	for _, sql := range []string{
		"CREATE ROLE IF NOT EXISTS r",
		"CREATE USER IF NOT EXISTS u",
		"CREATE EVENT IF NOT EXISTS e ON SCHEDULE AT NOW() DO COMMIT",
		"CREATE FUNCTION IF NOT EXISTS f RETURNS STRING SONAME 's'",
		"CREATE AGGREGATE FUNCTION IF NOT EXISTS f RETURNS STRING SONAME 's'",
		"CREATE PROCEDURE IF NOT EXISTS p() COMMIT",
		"DROP FUNCTION IF EXISTS d.f",
		"DROP PROCEDURE IF EXISTS p",
		"DROP USER IF EXISTS u",
		"DROP ROLE IF EXISTS r",
		"DROP EVENT IF EXISTS e",
		"DROP TRIGGER IF EXISTS t",
		"DROP SERVER IF EXISTS s",
		"DROP SPATIAL REFERENCE SYSTEM IF EXISTS 1",
		"REVOKE IF EXISTS SELECT ON * FROM u",
		"REVOKE IF EXISTS ALL, GRANT OPTION FROM u",
		"REVOKE IF EXISTS PROXY ON a FROM u",
		"REVOKE IF EXISTS r FROM u",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("%s: %v", sql, err)
		}
	}
	for _, sql := range []string{
		"CREATE ROLE IF EXISTS r",
		"CREATE ROLE IF NOT r",
		"CREATE USER IF EXISTS u",
		"CREATE USER IF NOT u",
		"CREATE EVENT IF EXISTS e ON SCHEDULE AT NOW() DO COMMIT",
		"CREATE EVENT IF NOT e ON SCHEDULE AT NOW() DO COMMIT",
		"CREATE FUNCTION IF EXISTS f RETURNS STRING SONAME 's'",
		"CREATE AGGREGATE FUNCTION IF NOT f RETURNS STRING SONAME 's'",
		"CREATE PROCEDURE IF EXISTS p() COMMIT",
		"CREATE PROCEDURE IF NOT p() COMMIT",
		"CREATE TRIGGER IF EXISTS t BEFORE INSERT ON a FOR EACH ROW SET @x = 1",
		"DROP FUNCTION IF d.f",
		"DROP FUNCTION IF f",
		"DROP PROCEDURE IF p",
		"DROP USER IF u",
		"DROP ROLE IF r",
		"DROP EVENT IF e",
		"DROP TRIGGER IF t",
		"DROP SERVER IF s",
		"DROP SPATIAL REFERENCE SYSTEM IF 1",
		"REVOKE IF SELECT ON * FROM u",
		"REVOKE IF ALL, GRANT OPTION FROM u",
		"REVOKE IF PROXY ON a FROM u",
		"REVOKE IF r FROM u",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("%s: no error", sql)
		}
	}
}
