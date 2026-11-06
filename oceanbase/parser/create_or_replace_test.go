package parser

import (
	"errors"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// orReplaceCase is a CREATE statement and how OceanBase parses it, as PREPARE
// reports it on 5.7.25-OceanBase_CE-v4.3.5.6 (TestCreateOrReplaceOracle,
// measured 2026-11-05; until then these were MySQL 8.0's answers):
//
//   - code 0: PREPARE succeeds.
//   - code 1064 (ER_PARSE_ERROR): rejected. near is the rest of the statement
//     from the token OceanBase rejects, so omni's ParseError must sit at
//     len(sql)-len(near). OceanBase's message quotes from that token, often the
//     token alone, and at the end of the input it quotes the last token.
//   - any other code: the statement parsed and the resolver refused it, as with
//     1368 for a CHECK OPTION on a view that cannot take one.
type orReplaceCase struct {
	sql  string
	code uint16
	near string
}

// createOrReplaceCases holds every CREATE form with OR REPLACE. OceanBase takes
// OR REPLACE only before a view, with ALGORITHM, DEFINER (only after ALGORITHM)
// and SQL SECURITY after it; anything else fails at the keyword after OR
// REPLACE, SPATIAL REFERENCE SYSTEM included.
var createOrReplaceCases = []orReplaceCase{
	{"CREATE OR REPLACE VIEW v AS SELECT 1", 0, ""},
	{"CREATE OR REPLACE ALGORITHM = MERGE VIEW v AS SELECT 1", 0, ""},
	{"CREATE OR REPLACE ALGORITHM = UNDEFINED DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v AS SELECT 1", 0, ""},
	{"CREATE OR REPLACE ALGORITHM = MERGE DEFINER = 'root'@'localhost' SQL SECURITY DEFINER VIEW v AS SELECT 1", 0, ""},
	{"CREATE OR REPLACE SQL SECURITY INVOKER VIEW v AS SELECT 1", 0, ""},
	{"CREATE OR REPLACE VIEW v (a) AS SELECT 1 WITH CHECK OPTION", 1368, ""},

	// DEFINER heads a view only after ALGORITHM; without it, CREATE DEFINER
	// heads a stored program, and OR REPLACE DEFINER is no form at all.
	{"CREATE OR REPLACE DEFINER = CURRENT_USER VIEW v AS SELECT 1", 1064, "DEFINER = CURRENT_USER VIEW v AS SELECT 1"},
	{"CREATE OR REPLACE DEFINER = 'root'@'localhost' SQL SECURITY DEFINER VIEW v AS SELECT 1", 1064, "DEFINER = 'root'@'localhost' SQL SECURITY DEFINER VIEW v AS SELECT 1"},
	{"CREATE DEFINER = CURRENT_USER VIEW v AS SELECT 1", 1064, "VIEW v AS SELECT 1"},
	{"CREATE DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v AS SELECT 1", 1064, "SQL SECURITY INVOKER VIEW v AS SELECT 1"},
	{"CREATE DEFINER = CURRENT_USER SQL SECURITY INVOKER PROCEDURE p1() SELECT 1", 1064, "SQL SECURITY INVOKER PROCEDURE p1() SELECT 1"},
	{"CREATE ALGORITHM = MERGE DEFINER = CURRENT_USER VIEW v AS SELECT 1", 0, ""},
	{"ALTER DEFINER = CURRENT_USER VIEW v AS SELECT 1", 1064, "VIEW v AS SELECT 1"},
	{"ALTER DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v AS SELECT 1", 1064, "SQL SECURITY INVOKER VIEW v AS SELECT 1"},
	{"ALTER ALGORITHM = MERGE DEFINER = CURRENT_USER VIEW v AS SELECT 1", 0, ""},
	{"ALTER SQL SECURITY INVOKER VIEW v AS SELECT 1", 0, ""},

	// OceanBase has no SPATIAL REFERENCE SYSTEM.
	{"CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME 'Greek' DEFINITION 'GEOGCS[\"Greek\",DATUM[\"Greek\",SPHEROID[\"Bessel 1841\",6377397.155,299.1528128]],PRIMEM[\"Greenwich\",0],UNIT[\"degree\",0.017453292519943278],AXIS[\"Lat\",NORTH],AXIS[\"Lon\",EAST]]'", 1064, "SPATIAL REFERENCE SYSTEM 4120 NAME 'Greek' DEFINITION 'GEOGCS[\"Greek\",DATUM[\"Greek\",SPHEROID[\"Bessel 1841\",6377397.155,299.1528128]],PRIMEM[\"Greenwich\",0],UNIT[\"degree\",0.017453292519943278],AXIS[\"Lat\",NORTH],AXIS[\"Lon\",EAST]]'"},

	// OR REPLACE heads the view clauses; it does not follow them.
	{"CREATE ALGORITHM = MERGE OR REPLACE VIEW v AS SELECT 1", 1064, "OR REPLACE VIEW v AS SELECT 1"},
	{"CREATE DEFINER = CURRENT_USER OR REPLACE VIEW v AS SELECT 1", 1064, "OR REPLACE VIEW v AS SELECT 1"},
	{"CREATE OR REPLACE OR REPLACE VIEW v AS SELECT 1", 1064, "OR REPLACE VIEW v AS SELECT 1"},
	{"CREATE OR VIEW v AS SELECT 1", 1064, "VIEW v AS SELECT 1"},
	{"CREATE OR REPLACE", 1064, ""},
	{"CREATE OR", 1064, ""},
	{"CREATE OR REPLACE TEMPORARY VIEW v AS SELECT 1", 1064, "TEMPORARY VIEW v AS SELECT 1"},
	{"CREATE OR REPLACE TEMPORARY TABLE r (id INT)", 1064, "TEMPORARY TABLE r (id INT)"},
	{"CREATE OR REPLACE SPATIAL REFERENCE SYSTEM IF NOT EXISTS 4120 NAME 'Greek'", 1064, "SPATIAL REFERENCE SYSTEM IF NOT EXISTS 4120 NAME 'Greek'"},
	{"CREATE OR REPLACE SPATIAL INDEX ix ON t (g)", 1064, "SPATIAL INDEX ix ON t (g)"},

	{"CREATE OR REPLACE TABLE r (id INT)", 1064, "TABLE r (id INT)"},
	{"CREATE OR REPLACE DATABASE d1", 1064, "DATABASE d1"},
	{"CREATE OR REPLACE SCHEMA d1", 1064, "SCHEMA d1"},
	{"CREATE OR REPLACE INDEX ix ON t (n)", 1064, "INDEX ix ON t (n)"},
	{"CREATE OR REPLACE UNIQUE INDEX ix ON t (n)", 1064, "UNIQUE INDEX ix ON t (n)"},
	{"CREATE OR REPLACE FULLTEXT INDEX ix ON t (n)", 1064, "FULLTEXT INDEX ix ON t (n)"},
	{"CREATE OR REPLACE USER u1", 1064, "USER u1"},
	{"CREATE OR REPLACE ROLE r1", 1064, "ROLE r1"},
	{"CREATE OR REPLACE FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1", 1064, "FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1"},
	{"CREATE OR REPLACE FUNCTION f1 RETURNS INTEGER SONAME 'x.so'", 1064, "FUNCTION f1 RETURNS INTEGER SONAME 'x.so'"},
	{"CREATE OR REPLACE AGGREGATE FUNCTION f1 RETURNS INTEGER SONAME 'x.so'", 1064, "AGGREGATE FUNCTION f1 RETURNS INTEGER SONAME 'x.so'"},
	{"CREATE OR REPLACE PROCEDURE p1() SELECT 1", 1064, "PROCEDURE p1() SELECT 1"},
	{"CREATE OR REPLACE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 1064, "TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1"},
	{"CREATE OR REPLACE EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", 1064, "EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1"},
	{"CREATE OR REPLACE TABLESPACE ts1 ADD DATAFILE 'x.ibd'", 1064, "TABLESPACE ts1 ADD DATAFILE 'x.ibd'"},
	{"CREATE OR REPLACE UNDO TABLESPACE ts1 ADD DATAFILE 'x.ibu'", 1064, "UNDO TABLESPACE ts1 ADD DATAFILE 'x.ibu'"},
	{"CREATE OR REPLACE SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'x')", 1064, "SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'x')"},
	{"CREATE OR REPLACE LOGFILE GROUP lg ADD UNDOFILE 'u.dat' ENGINE = NDB", 1064, "LOGFILE GROUP lg ADD UNDOFILE 'u.dat' ENGINE = NDB"},
	{"CREATE OR REPLACE RESOURCE GROUP rg TYPE = USER", 1064, "RESOURCE GROUP rg TYPE = USER"},

	// Nor does OR REPLACE take DEFINER before a stored program.
	{"CREATE OR REPLACE DEFINER = CURRENT_USER FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1", 1064, "DEFINER = CURRENT_USER FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1"},
	{"CREATE OR REPLACE DEFINER = CURRENT_USER PROCEDURE p1() SELECT 1", 1064, "DEFINER = CURRENT_USER PROCEDURE p1() SELECT 1"},
	{"CREATE OR REPLACE DEFINER = CURRENT_USER TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 1064, "DEFINER = CURRENT_USER TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1"},
	{"CREATE OR REPLACE DEFINER = CURRENT_USER EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1", 1064, "DEFINER = CURRENT_USER EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1"},
	{"CREATE OR REPLACE ALGORITHM = MERGE TABLE r (id INT)", 1064, "TABLE r (id INT)"},

	// The same objects without OR REPLACE parse.
	{"CREATE TABLE r (id INT)", 0, ""},
	{"CREATE TEMPORARY TABLE r (id INT)", 0, ""},
	{"CREATE INDEX ix ON t (n)", 0, ""},
	{"CREATE USER u1", 0, ""},
	{"CREATE FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1", 0, ""},
	{"CREATE ALGORITHM = MERGE VIEW v AS SELECT 1", 0, ""},
	{"CREATE ALGORITHM = MERGE DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v AS SELECT 1", 0, ""},
}

// parsesInEngine reports whether code means OceanBase parsed the statement.
func parsesInEngine(code uint16) bool {
	return code != 1064
}

// TestCreateOrReplace checks omni against createOrReplaceCases: it accepts
// what OceanBase parses and rejects the rest at the token OceanBase rejects.
func TestCreateOrReplace(t *testing.T) {
	for _, tc := range createOrReplaceCases {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if parsesInEngine(tc.code) {
				if err != nil {
					t.Fatalf("Parse: %v, want success", err)
				}
				return
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse error = %v, want *ParseError", err)
			}
			if want := len(tc.sql) - len(tc.near); pe.Position != want {
				t.Errorf("error at %d, want %d (near %q): %v", pe.Position, want, tc.near, pe)
			}
		})
	}
}

// TestCreateOrReplaceAST checks the nodes the accepted OR REPLACE forms build.
func TestCreateOrReplaceAST(t *testing.T) {
	views := []struct {
		sql                             string
		algorithm, definer, sqlSecurity string
	}{
		{"CREATE OR REPLACE VIEW v AS SELECT 1", "", "", ""},
		{"CREATE OR REPLACE ALGORITHM = MERGE VIEW v AS SELECT 1", "MERGE", "", ""},
		{"CREATE OR REPLACE ALGORITHM = UNDEFINED DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v AS SELECT 1", "UNDEFINED", "CURRENT_USER", "INVOKER"},
		{"CREATE OR REPLACE ALGORITHM = MERGE DEFINER = CURRENT_USER VIEW v AS SELECT 1", "MERGE", "CURRENT_USER", ""},
		{"CREATE OR REPLACE SQL SECURITY DEFINER VIEW v AS SELECT 1", "", "", "DEFINER"},
	}
	for _, tc := range views {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneOrReplace(t, tc.sql).(*nodes.CreateViewStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.CreateViewStmt", stmt)
			}
			if !stmt.OrReplace {
				t.Error("OrReplace = false, want true")
			}
			if stmt.Algorithm != tc.algorithm || stmt.Definer != tc.definer || stmt.SqlSecurity != tc.sqlSecurity {
				t.Errorf("Algorithm, Definer, SqlSecurity = %q, %q, %q; want %q, %q, %q",
					stmt.Algorithm, stmt.Definer, stmt.SqlSecurity, tc.algorithm, tc.definer, tc.sqlSecurity)
			}
			if stmt.Loc.Start != 0 || stmt.Loc.End != len(tc.sql) {
				t.Errorf("Loc = %+v, want {0 %d}", stmt.Loc, len(tc.sql))
			}
		})
	}
}

func parseOneOrReplace(t *testing.T, sql string) nodes.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) = %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}
