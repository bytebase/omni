package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// orReplaceObjectOracleCases pins, for every CREATE object omni parses other
// than TABLE (orReplaceTableOracleCases), whether MariaDB 11.8.8 takes OR
// REPLACE and what it raises for OR REPLACE with IF NOT EXISTS: ER_WRONG_USAGE
// (1221) from LEX::check_create_options while parsing. MariaDB has none of
// TABLESPACE, UNDO TABLESPACE, LOGFILE GROUP, RESOURCE GROUP, or SPATIAL
// REFERENCE SYSTEM, so OR REPLACE on them is 1064.
var orReplaceObjectOracleCases = []parseCodeCase{
	// SEQUENCE
	{"CREATE OR REPLACE SEQUENCE s1", 0},
	{"CREATE OR REPLACE TEMPORARY SEQUENCE s1", 0},
	{"CREATE OR REPLACE SEQUENCE IF NOT EXISTS cs_both", 1221},
	{"CREATE OR REPLACE TEMPORARY SEQUENCE IF NOT EXISTS s1", 1221},
	// INDEX
	{"CREATE OR REPLACE INDEX ix ON t (n)", 0},
	{"CREATE OR REPLACE INDEX ix USING BTREE ON t (n)", 0},
	{"CREATE OR REPLACE UNIQUE INDEX ux ON t (n)", 0},
	{"CREATE OR REPLACE FULLTEXT INDEX fx ON t (n)", 0},
	{"CREATE OR REPLACE SPATIAL INDEX sx ON t (n)", 0},
	{"CREATE OR REPLACE INDEX IF NOT EXISTS ix ON t (n)", 1221},
	{"CREATE OR REPLACE UNIQUE INDEX IF NOT EXISTS ux ON t (n)", 1221},
	{"CREATE OR REPLACE FULLTEXT INDEX IF NOT EXISTS fx ON t (n)", 1221},
	{"CREATE OR REPLACE SPATIAL INDEX IF NOT EXISTS sx ON t (n)", 1221},
	{"CREATE OR REPLACE INDEX IF EXISTS ix ON t (n)", 1064},
	// DATABASE / SCHEMA
	{"CREATE OR REPLACE DATABASE d1", 0},
	{"CREATE OR REPLACE SCHEMA d1", 0},
	{"CREATE OR REPLACE DATABASE IF NOT EXISTS d1", 1221},
	{"CREATE OR REPLACE SCHEMA IF NOT EXISTS d1", 1221},
	// VIEW, through each prefix path. omni does not parse IF NOT EXISTS on a
	// view yet, so it rejects the 1221 forms with a syntax error instead.
	{"CREATE OR REPLACE VIEW v1 AS SELECT 1", 0},
	{"CREATE OR REPLACE ALGORITHM=MERGE VIEW v1 AS SELECT 1", 0},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER VIEW v1 AS SELECT 1", 0},
	{"CREATE OR REPLACE SQL SECURITY INVOKER VIEW v1 AS SELECT 1", 0},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER SQL SECURITY INVOKER VIEW v1 AS SELECT 1", 0},
	{"CREATE OR REPLACE ALGORITHM=MERGE DEFINER=CURRENT_USER SQL SECURITY INVOKER VIEW v1 AS SELECT 1", 0},
	{"CREATE OR REPLACE VIEW IF NOT EXISTS v1 AS SELECT 1", 1221},
	{"CREATE OR REPLACE ALGORITHM=MERGE VIEW IF NOT EXISTS v1 AS SELECT 1", 1221},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER VIEW IF NOT EXISTS v1 AS SELECT 1", 1221},
	{"CREATE OR REPLACE SQL SECURITY INVOKER VIEW IF NOT EXISTS v1 AS SELECT 1", 1221},
	// FUNCTION and PROCEDURE, with and without DEFINER
	{"CREATE OR REPLACE FUNCTION f1() RETURNS INT RETURN 1", 0},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER FUNCTION f1() RETURNS INT RETURN 1", 0},
	{"CREATE OR REPLACE FUNCTION IF NOT EXISTS f1() RETURNS INT RETURN 1", 1221},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER FUNCTION IF NOT EXISTS f1() RETURNS INT RETURN 1", 1221},
	{"CREATE OR REPLACE PROCEDURE p1() SELECT 1", 0},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER PROCEDURE p1() SELECT 1", 0},
	{"CREATE OR REPLACE PROCEDURE IF NOT EXISTS p1() SELECT 1", 1221},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER PROCEDURE IF NOT EXISTS p1() SELECT 1", 1221},
	// Loadable functions (UDF), plain and AGGREGATE
	{"CREATE OR REPLACE FUNCTION udf1 RETURNS STRING SONAME 'x.so'", 0},
	{"CREATE OR REPLACE FUNCTION IF NOT EXISTS udf1 RETURNS STRING SONAME 'x.so'", 1221},
	{"CREATE OR REPLACE AGGREGATE FUNCTION agg1 RETURNS INTEGER SONAME 'x.so'", 0},
	{"CREATE OR REPLACE AGGREGATE FUNCTION IF NOT EXISTS agg1 RETURNS INTEGER SONAME 'x.so'", 1221},
	// TRIGGER
	{"CREATE OR REPLACE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 0},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 0},
	{"CREATE OR REPLACE TRIGGER IF NOT EXISTS tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 1221},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER TRIGGER IF NOT EXISTS tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 1221},
	// EVENT
	{"CREATE OR REPLACE EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", 0},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", 0},
	{"CREATE OR REPLACE EVENT IF NOT EXISTS e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", 1221},
	{"CREATE OR REPLACE DEFINER=CURRENT_USER EVENT IF NOT EXISTS e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", 1221},
	// USER and ROLE
	{"CREATE OR REPLACE USER u1", 0},
	{"CREATE OR REPLACE USER u1@'%' IDENTIFIED BY 'x', u2", 0},
	{"CREATE OR REPLACE USER IF NOT EXISTS u1", 1221},
	{"CREATE OR REPLACE ROLE ro1", 0},
	{"CREATE OR REPLACE ROLE IF NOT EXISTS ro1", 1221},
	// SERVER
	{"CREATE OR REPLACE SERVER sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')", 0},
	{"CREATE SERVER IF NOT EXISTS sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')", 0},
	{"CREATE OR REPLACE SERVER IF NOT EXISTS sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')", 1221},
	// MySQL objects MariaDB does not have
	{"CREATE OR REPLACE TABLESPACE ts1 ADD DATAFILE 'x.ibd'", 1064},
	{"CREATE OR REPLACE UNDO TABLESPACE uts1 ADD DATAFILE 'x.ibu'", 1064},
	{"CREATE OR REPLACE LOGFILE GROUP lg1 ADD UNDOFILE 'u.dat' ENGINE=NDB", 1064},
	{"CREATE OR REPLACE RESOURCE GROUP rg1 TYPE = USER", 1064},
	{"CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME 'x' DEFINITION 'y'", 1064},
	// OR without REPLACE
	{"CREATE OR INDEX ix ON t (n)", 1064},
}

// createModifierOracleCases pins the other CREATE prefixes the dispatcher
// consumes. TEMPORARY precedes only TABLE and SEQUENCE in MariaDB's create
// rule; IF NOT EXISTS must be spelled out in full.
var createModifierOracleCases = []parseCodeCase{
	// TEMPORARY
	{"CREATE TEMPORARY SEQUENCE IF NOT EXISTS s1", 0},
	{"CREATE TEMPORARY VIEW v1 AS SELECT 1", 1064},
	{"CREATE OR REPLACE TEMPORARY VIEW v1 AS SELECT 1", 1064},
	{"CREATE TEMPORARY ALGORITHM=MERGE VIEW v1 AS SELECT 1", 1064},
	{"CREATE TEMPORARY DEFINER=CURRENT_USER VIEW v1 AS SELECT 1", 1064},
	{"CREATE TEMPORARY SQL SECURITY INVOKER VIEW v1 AS SELECT 1", 1064},
	{"CREATE TEMPORARY INDEX ix ON t (n)", 1064},
	{"CREATE TEMPORARY UNIQUE INDEX ux ON t (n)", 1064},
	{"CREATE TEMPORARY FULLTEXT INDEX fx ON t (n)", 1064},
	{"CREATE TEMPORARY SPATIAL INDEX sx ON t (n)", 1064},
	{"CREATE TEMPORARY DATABASE d1", 1064},
	{"CREATE TEMPORARY SCHEMA d1", 1064},
	{"CREATE TEMPORARY FUNCTION f1() RETURNS INT RETURN 1", 1064},
	{"CREATE TEMPORARY AGGREGATE FUNCTION agg1 RETURNS INTEGER SONAME 'x.so'", 1064},
	{"CREATE TEMPORARY PROCEDURE p1() SELECT 1", 1064},
	{"CREATE TEMPORARY TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 1064},
	{"CREATE TEMPORARY EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", 1064},
	{"CREATE TEMPORARY USER u1", 1064},
	{"CREATE TEMPORARY ROLE ro1", 1064},
	{"CREATE TEMPORARY SERVER sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')", 1064},
	{"CREATE TEMPORARY TABLESPACE ts1 ADD DATAFILE 'x.ibd'", 1064},
	// IF NOT EXISTS, spelled in full
	{"CREATE FUNCTION IF NOT EXISTS f1() RETURNS INT RETURN 1", 0},
	{"CREATE FUNCTION IF f1() RETURNS INT RETURN 1", 1064},
	{"CREATE FUNCTION IF NOT f1() RETURNS INT RETURN 1", 1064},
	{"CREATE FUNCTION IF EXISTS f1() RETURNS INT RETURN 1", 1064},
	{"CREATE PROCEDURE IF NOT EXISTS p1() SELECT 1", 0},
	{"CREATE PROCEDURE IF EXISTS p1() SELECT 1", 1064},
	{"CREATE AGGREGATE FUNCTION IF NOT EXISTS agg1 RETURNS INTEGER SONAME 'x.so'", 0},
	{"CREATE AGGREGATE FUNCTION IF agg1 RETURNS INTEGER SONAME 'x.so'", 1064},
	{"CREATE FUNCTION IF EXISTS udf1 RETURNS STRING SONAME 'x.so'", 1064},
	{"CREATE TRIGGER IF NOT tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 1064},
	{"CREATE TRIGGER IF EXISTS tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 1064},
	{"CREATE EVENT IF e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", 1064},
	{"CREATE EVENT IF EXISTS e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", 1064},
	{"CREATE USER IF u1", 1064},
	{"CREATE USER IF NOT u1", 1064},
	{"CREATE USER IF EXISTS u1", 1064},
	{"CREATE ROLE IF ro1", 1064},
	{"CREATE ROLE IF EXISTS ro1", 1064},
	{"CREATE SERVER IF sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')", 1064},
}

// TestMariaDBOrReplaceObjectOracle asserts omni agrees with a live
// mariadb:11.8.8 on orReplaceObjectOracleCases (container test).
func TestMariaDBOrReplaceObjectOracle(t *testing.T) {
	assertParseCodes(t, orReplaceObjectOracleCases)
}

// TestMariaDBCreateModifierOracle asserts omni agrees with a live
// mariadb:11.8.8 on createModifierOracleCases (container test).
func TestMariaDBCreateModifierOracle(t *testing.T) {
	assertParseCodes(t, createModifierOracleCases)
}

// TestCreateOrReplaceRecorded checks that every CREATE statement MariaDB takes
// with OR REPLACE records it on its node and prints it as :or_replace, and that
// the plain form leaves it unset.
func TestCreateOrReplaceRecorded(t *testing.T) {
	cases := []struct {
		sql string
		tag string
	}{
		{"CREATE OR REPLACE SEQUENCE s1", "CREATE_SEQUENCE"},
		{"CREATE OR REPLACE INDEX ix ON t (n)", "CREATE_INDEX"},
		{"CREATE OR REPLACE UNIQUE INDEX ux ON t (n)", "CREATE_INDEX"},
		{"CREATE OR REPLACE FULLTEXT INDEX fx ON t (n)", "CREATE_INDEX"},
		{"CREATE OR REPLACE SPATIAL INDEX sx ON t (n)", "CREATE_INDEX"},
		{"CREATE OR REPLACE DATABASE d1", "CREATE_DATABASE"},
		{"CREATE OR REPLACE SCHEMA d1", "CREATE_DATABASE"},
		{"CREATE OR REPLACE VIEW v1 AS SELECT 1", "CREATE_VIEW"},
		{"CREATE OR REPLACE ALGORITHM=MERGE VIEW v1 AS SELECT 1", "CREATE_VIEW"},
		{"CREATE OR REPLACE DEFINER=CURRENT_USER SQL SECURITY INVOKER VIEW v1 AS SELECT 1", "CREATE_VIEW"},
		{"CREATE OR REPLACE FUNCTION f1() RETURNS INT RETURN 1", "CREATE_FUNCTION"},
		{"CREATE OR REPLACE DEFINER=CURRENT_USER FUNCTION f1() RETURNS INT RETURN 1", "CREATE_FUNCTION"},
		{"CREATE OR REPLACE FUNCTION udf1 RETURNS STRING SONAME 'x.so'", "CREATE_FUNCTION"},
		{"CREATE OR REPLACE AGGREGATE FUNCTION agg1 RETURNS INTEGER SONAME 'x.so'", "CREATE_FUNCTION"},
		{"CREATE OR REPLACE PROCEDURE p1() SELECT 1", "CREATE_FUNCTION"},
		{"CREATE OR REPLACE DEFINER=CURRENT_USER PROCEDURE p1() SELECT 1", "CREATE_FUNCTION"},
		{"CREATE OR REPLACE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", "CREATE_TRIGGER"},
		{"CREATE OR REPLACE DEFINER=CURRENT_USER TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", "CREATE_TRIGGER"},
		{"CREATE OR REPLACE EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", "CREATE_EVENT"},
		{"CREATE OR REPLACE DEFINER=CURRENT_USER EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", "CREATE_EVENT"},
		{"CREATE OR REPLACE USER u1", "CREATE_USER"},
		{"CREATE OR REPLACE ROLE ro1", "CREATE_ROLE"},
		{"CREATE OR REPLACE SERVER sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')", "CREATE_SERVER"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			node := parseSeqStmt[ast.Node](t, tc.sql)
			if got, ok := orReplaceOf(node); !ok || !got {
				t.Errorf("%T OrReplace = %v (field found %v), want true", node, got, ok)
			}
			if got, want := ast.NodeToString(node), "{"+tc.tag+" :loc 0 :or_replace true"; !strings.HasPrefix(got, want) {
				t.Errorf("NodeToString = %s, want prefix %s", got, want)
			}

			plain := strings.Replace(tc.sql, "CREATE OR REPLACE ", "CREATE ", 1)
			node = parseSeqStmt[ast.Node](t, plain)
			if got, ok := orReplaceOf(node); !ok || got {
				t.Errorf("%q: %T OrReplace = %v (field found %v), want false", plain, node, got, ok)
			}
			if got := ast.NodeToString(node); strings.Contains(got, ":or_replace") {
				t.Errorf("%q: NodeToString = %s, want no :or_replace", plain, got)
			}
		})
	}
}

// orReplaceOf returns the OrReplace field of a CREATE node and whether the node
// has one.
func orReplaceOf(node ast.Node) (bool, bool) {
	switch n := node.(type) {
	case *ast.CreateSequenceStmt:
		return n.OrReplace, true
	case *ast.CreateIndexStmt:
		return n.OrReplace, true
	case *ast.CreateDatabaseStmt:
		return n.OrReplace, true
	case *ast.CreateViewStmt:
		return n.OrReplace, true
	case *ast.CreateFunctionStmt:
		return n.OrReplace, true
	case *ast.CreateTriggerStmt:
		return n.OrReplace, true
	case *ast.CreateEventStmt:
		return n.OrReplace, true
	case *ast.CreateUserStmt:
		return n.OrReplace, true
	case *ast.CreateRoleStmt:
		return n.OrReplace, true
	case *ast.CreateServerStmt:
		return n.OrReplace, true
	}
	return false, false
}

// TestCreateOrReplaceIfNotExistsError checks the error omni gives for OR
// REPLACE with IF NOT EXISTS: MariaDB's ER_WRONG_USAGE text, at IF.
func TestCreateOrReplaceIfNotExistsError(t *testing.T) {
	for _, tc := range orReplaceObjectOracleCases {
		if tc.code != 1221 || strings.Contains(tc.sql, " VIEW ") {
			continue
		}
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %v, want *ParseError", tc.sql, err)
			}
			if pe.Message != "Incorrect usage of OR REPLACE and IF NOT EXISTS" {
				t.Errorf("Message = %q", pe.Message)
			}
			if want := strings.Index(tc.sql, "IF NOT EXISTS"); pe.Position != want {
				t.Errorf("Position = %d, want %d", pe.Position, want)
			}
		})
	}
}

// TestCreateModifierErrorPosition checks where omni reports the 1064 forms: at
// the object keyword after an OR REPLACE MariaDB does not take, and at the
// token after a TEMPORARY that precedes neither TABLE nor SEQUENCE, matching
// the "near" text MariaDB reports.
func TestCreateModifierErrorPosition(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"CREATE OR REPLACE TABLESPACE ts1 ADD DATAFILE 'x.ibd'", "TABLESPACE"},
		{"CREATE OR REPLACE UNDO TABLESPACE uts1 ADD DATAFILE 'x.ibu'", "UNDO"},
		{"CREATE OR REPLACE LOGFILE GROUP lg1 ADD UNDOFILE 'u.dat' ENGINE=NDB", "LOGFILE"},
		{"CREATE OR REPLACE RESOURCE GROUP rg1 TYPE = USER", "RESOURCE"},
		{"CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME 'x' DEFINITION 'y'", "REFERENCE"},
		{"CREATE TEMPORARY VIEW v1 AS SELECT 1", "VIEW"},
		{"CREATE OR REPLACE TEMPORARY INDEX ix ON t (n)", "INDEX"},
		{"CREATE TEMPORARY DEFINER=CURRENT_USER VIEW v1 AS SELECT 1", "DEFINER"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %v, want *ParseError", tc.sql, err)
			}
			if want := strings.Index(tc.sql, tc.near); pe.Position != want {
				t.Errorf("Position = %d, want %d (%s)", pe.Position, want, pe.Message)
			}
		})
	}
}
