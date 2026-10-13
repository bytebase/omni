package parser

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestCreateTypeHeading checks the heading clauses, subtypes, incomplete
// types, and collection elements of CREATE TYPE against Oracle 23ai: every
// accepted form compiles there (or fails only on names it cannot resolve or
// on the instance's settings), and every rejected one fails with the error
// noted beside it. Rows ref_189 to ref_208, ref_216 to ref_218, ref_223 to
// ref_224, and ref_247 to ref_250 run the PLS-00103 and ORA cases against
// the engine; the others are errors Oracle reports from the text alone after
// parsing, which the reference harness counts as accepted.
func TestCreateTypeHeading(t *testing.T) {
	for _, sql := range []string{
		// FORCE and OID in either order, then AUTHID and ACCESSIBLE BY in
		// either order; SHARING first.
		"CREATE TYPE t FORCE AS OBJECT (a NUMBER)",
		"CREATE OR REPLACE TYPE t FORCE AS VARRAY(10) OF NUMBER",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDEF' AS OBJECT (a NUMBER)",
		"CREATE TYPE t FORCE OID '1234567890ABCDEF1234567890ABCDEF' AS OBJECT (a NUMBER)",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDEF' FORCE AS OBJECT (a NUMBER)",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDEF' AS VARRAY(3) OF NUMBER",
		"CREATE TYPE t AUTHID CURRENT_USER AS OBJECT (a NUMBER)",
		"CREATE TYPE t AUTHID DEFINER ACCESSIBLE BY (PROCEDURE p) AS OBJECT (a NUMBER)",
		"CREATE TYPE t ACCESSIBLE BY (PROCEDURE p) AUTHID DEFINER AS OBJECT (a NUMBER) NOT FINAL",
		"CREATE TYPE t SHARING = NONE FORCE OID '1234567890ABCDEF1234567890ABCDEF' AUTHID DEFINER ACCESSIBLE BY (PROCEDURE p) AS OBJECT (a NUMBER)",
		"CREATE TYPE t SHARING = NONE OID '1234567890ABCDEF1234567890ABCDEF' AS OBJECT (a NUMBER)",
		// METADATA and DATA are illegal outside an application root
		// (ORA-65021).
		"CREATE TYPE t SHARING = METADATA AS OBJECT (a NUMBER)",
		// DEFAULT COLLATION, checked with MAX_STRING_SIZE = EXTENDED: after
		// SHARING, FORCE, and OID, in any order with AUTHID and ACCESSIBLE BY,
		// and on a collection or a subtype too.
		"CREATE TYPE t DEFAULT COLLATION USING_NLS_COMP AS OBJECT (a NUMBER)",
		"CREATE TYPE t AUTHID DEFINER DEFAULT COLLATION USING_NLS_COMP AS OBJECT (a NUMBER)",
		"CREATE TYPE t DEFAULT COLLATION USING_NLS_COMP AUTHID DEFINER AS OBJECT (a NUMBER)",
		"CREATE TYPE t ACCESSIBLE BY (PROCEDURE p) DEFAULT COLLATION USING_NLS_COMP AS OBJECT (a NUMBER)",
		"CREATE TYPE t SHARING = NONE FORCE OID '1234567890ABCDEF1234567890ABCDEF' DEFAULT COLLATION USING_NLS_COMP AS OBJECT (a NUMBER)",
		"CREATE TYPE t DEFAULT COLLATION USING_NLS_COMP AS TABLE OF NUMBER",
		"CREATE TYPE t DEFAULT COLLATION USING_NLS_COMP UNDER base (b NUMBER)",
		"CREATE TYPE t IS OBJECT (a NUMBER)",
		"CREATE TYPE t IS TABLE OF NUMBER",
		// An incomplete type takes SHARING and OID only.
		"CREATE TYPE t",
		"CREATE TYPE t;",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDEF'",
		"CREATE TYPE t SHARING = NONE",
		// A subtype: no IS or AS, a parenthesized element list that may be
		// empty or hold methods only, FINAL and INSTANTIABLE modifiers.
		"CREATE TYPE t UNDER base (b NUMBER)",
		"CREATE TYPE t UNDER s.base (b NUMBER, MEMBER FUNCTION f RETURN NUMBER) NOT FINAL NOT INSTANTIABLE",
		"CREATE TYPE t UNDER base ()",
		"CREATE TYPE t UNDER base (MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t UNDER base (PRAGMA RESTRICT_REFERENCES (DEFAULT, WNDS))",
		"CREATE TYPE t UNDER base (b NUMBER) FINAL",
		// A subtype's MAP method depends on the supertype's (PLS-00636,
		// PLS-00646), which the text does not show.
		"CREATE TYPE t UNDER base (OVERRIDING MAP MEMBER FUNCTION m RETURN NUMBER)",
		"CREATE TYPE t FORCE UNDER base (b NUMBER)",
		"CREATE TYPE t ACCESSIBLE BY (PROCEDURE p) UNDER base (b NUMBER)",
		// The supertype's AUTHID decides (PLS-00722 when they differ).
		"CREATE TYPE t AUTHID DEFINER UNDER base (b NUMBER)",
		// Collection elements: NOT NULL bare, or PERSISTABLE after a
		// parenthesized element with NOT NULL inside.
		"CREATE TYPE t AS TABLE OF NUMBER NOT NULL",
		"CREATE TYPE t AS TABLE OF NUMBER(5) NOT NULL",
		"CREATE TYPE t AS TABLE OF (NUMBER)",
		"CREATE TYPE t AS TABLE OF (u) NOT PERSISTABLE",
		"CREATE TYPE t AS TABLE OF (NUMBER NOT NULL)",
		"CREATE TYPE t AS TABLE OF (NUMBER NOT NULL) NOT PERSISTABLE",
		"CREATE TYPE t AS VARRAY(3) OF (NUMBER) PERSISTABLE",
		"CREATE TYPE t AS VARRAY(3) OF NUMBER NOT NULL",
		"CREATE TYPE t AS VARRAY(3) OF (NUMBER NOT NULL) NOT PERSISTABLE",
		"CREATE TYPE t AS VARYING ARRAY(3) OF NUMBER",
		"CREATE TYPE t AS ARRAY(3) OF NUMBER",
		"CREATE TYPE t AS TABLE OF REF u",
		"CREATE TYPE t AS VARRAY(3) OF REF u",
		"CREATE TYPE t AS TABLE OF (REF u) NOT PERSISTABLE",
		"CREATE TYPE t AS TABLE OF VARCHAR2(10) CHARACTER SET NCHAR_CS",
		"CREATE TYPE t AS TABLE OF (VARCHAR2(10) CHARACTER SET NCHAR_CS NOT NULL) NOT PERSISTABLE",
		"CREATE TYPE t AS TABLE OF VARCHAR2(ORA_MAX_NAME_LEN)",
		"CREATE TYPE t AS VARRAY(3) OF VARCHAR2(ORA_MAX_NAME_LEN)",
		// Attributes: ANSI and REF datatypes, a named character set, and
		// PL/SQL size expressions.
		"CREATE TYPE t AS OBJECT (a CHARACTER(1 BYTE), b CHAR(2 CHAR), c CHARACTER VARYING(3), d DOUBLE PRECISION, e NATIONAL CHARACTER(2), f CHARACTER(2 CHAR))",
		"CREATE TYPE t AS OBJECT (a CHARACTER LARGE OBJECT)",
		"CREATE TYPE t AS OBJECT (a REF u, b REF s.u)",
		"CREATE TYPE t AS OBJECT (a VARCHAR2(10) CHARACTER SET NCHAR_CS)",
		"CREATE TYPE t AS OBJECT (a VARCHAR2(ORA_MAX_NAME_LEN), b VARCHAR2(10+2), c NUMBER(5+1, 2))",
		"CREATE TYPE t AS OBJECT (a NUMBER) NOT PERSISTABLE NOT FINAL",
		"CREATE TYPE t AS OBJECT (a NUMBER) NOT INSTANTIABLE NOT FINAL",
		"CREATE TYPE t AS OBJECT (a NUMBER) INSTANTIABLE FINAL",
		"CREATE TYPE BODY t IS MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END;",
	} {
		ParseAndCheck(t, sql)
	}

	for _, sql := range []string{
		// Each clause once, FORCE and OID before AUTHID and ACCESSIBLE BY,
		// SHARING first (PLS-00103), AUTHID and ACCESSIBLE BY once
		// (PLS-00371).
		"CREATE TYPE t FORCE FORCE AS OBJECT (a NUMBER)",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDE5' OID '1234567890ABCDEF1234567890ABCDE6' AS OBJECT (a NUMBER)",
		"CREATE TYPE t AUTHID DEFINER FORCE AS OBJECT (a NUMBER)",
		"CREATE TYPE t ACCESSIBLE BY (PROCEDURE p) OID '1234567890ABCDEF1234567890ABCDEF' AS OBJECT (a NUMBER)",
		"CREATE TYPE t AUTHID DEFINER AUTHID DEFINER AS OBJECT (a NUMBER)",
		"CREATE TYPE t ACCESSIBLE BY (PROCEDURE p) ACCESSIBLE BY (PROCEDURE q) AS OBJECT (a NUMBER)",
		"CREATE TYPE t FORCE SHARING = NONE AS OBJECT (a NUMBER)",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDEF' SHARING = NONE AS OBJECT (a NUMBER)",
		"CREATE TYPE t AUTHID DEFINER SHARING = NONE AS OBJECT (a NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER) NOT FINAL AUTHID DEFINER",
		"CREATE TYPE t DEFAULT COLLATION USING_NLS_COMP FORCE AS OBJECT (a NUMBER)",
		"CREATE TYPE t DEFAULT COLLATION USING_NLS_COMP OID '1234567890ABCDEF1234567890ABCDEF' AS OBJECT (a NUMBER)",
		"CREATE TYPE t DEFAULT COLLATION USING_NLS_COMP SHARING = NONE AS OBJECT (a NUMBER)",
		"CREATE TYPE t DEFAULT COLLATION USING_NLS_COMP DEFAULT COLLATION USING_NLS_COMP AS OBJECT (a NUMBER)",
		// Program unit collation is USING_NLS_COMP (PLS-00761).
		"CREATE TYPE t DEFAULT COLLATION BINARY AS OBJECT (a NUMBER)",
		// OID takes a string literal (ORA-01780) of 32 hexadecimal digits
		// (ORA-02304); SHARING its '=' (ORA-02000) and a valid value
		// (ORA-65014).
		"CREATE TYPE t OID 123 AS OBJECT (a NUMBER)",
		"CREATE TYPE t OID 'ABC' AS OBJECT (a NUMBER)",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDEG' AS OBJECT (a NUMBER)",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCD1' AS OBJECT (a NUMBER)",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCD123' AS OBJECT (a NUMBER)",
		"CREATE TYPE t SHARING NONE AS OBJECT (a NUMBER)",
		"CREATE TYPE t SHARING = FOO AS OBJECT (a NUMBER)",
		"CREATE TYPE t SHARING = EXTENDED DATA AS OBJECT (a NUMBER)",
		// IS or AS is required, except before UNDER (PLS-00103).
		"CREATE TYPE t OBJECT (a NUMBER)",
		"CREATE TYPE t TABLE OF NUMBER",
		"CREATE TYPE t AS UNDER base (b NUMBER)",
		"CREATE TYPE BODY t MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END;",
		// An incomplete type takes no FORCE, DEFAULT COLLATION, AUTHID, or
		// ACCESSIBLE BY.
		"CREATE TYPE t FORCE",
		"CREATE TYPE t DEFAULT COLLATION USING_NLS_COMP",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDEF' DEFAULT COLLATION USING_NLS_COMP",
		"CREATE TYPE t AUTHID DEFINER",
		"CREATE TYPE t ACCESSIBLE BY (PROCEDURE p)",
		"CREATE TYPE t FORCE OID '1234567890ABCDEF1234567890ABCDEF'",
		"CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDEF' FORCE",
		// AUTHID and ACCESSIBLE BY on a collection type (PLS-00262).
		"CREATE TYPE t AUTHID CURRENT_USER AS TABLE OF NUMBER",
		"CREATE TYPE t ACCESSIBLE BY (PROCEDURE p) AS TABLE OF NUMBER",
		"CREATE TYPE t AUTHID CURRENT_USER AS VARRAY(3) OF NUMBER",
		// A subtype needs its element list (PLS-00103), takes no database
		// link (PLS-00103), no PERSISTABLE (PLS-00772), and no attribute
		// after a method (PLS-00311).
		"CREATE TYPE t UNDER base",
		"CREATE TYPE t UNDER base NOT FINAL",
		"CREATE TYPE t UNDER base@lnk (b NUMBER)",
		"CREATE TYPE t UNDER base (b NUMBER,)",
		"CREATE TYPE t UNDER base (b NUMBER) NOT PERSISTABLE",
		"CREATE TYPE t UNDER base (MEMBER FUNCTION f RETURN NUMBER, b NUMBER)",
		// A subtype has no ORDER method (PLS-00646, or PLS-00154 under a
		// supertype with MAP or ORDER).
		"CREATE TYPE t UNDER base (ORDER MEMBER FUNCTION o (x t) RETURN NUMBER)",
		"CREATE TYPE t UNDER base (b NUMBER, OVERRIDING ORDER MEMBER FUNCTION o (x t) RETURN NUMBER)",
		// An object type is FINAL unless it says NOT FINAL, and a FINAL type
		// cannot be NOT INSTANTIABLE (PLS-00614).
		"CREATE TYPE t AS OBJECT (a NUMBER) NOT INSTANTIABLE",
		"CREATE TYPE t AS OBJECT (a NUMBER) FINAL NOT INSTANTIABLE",
		"CREATE TYPE t UNDER base (b NUMBER) NOT INSTANTIABLE",
		"CREATE TYPE t UNDER base (b NUMBER) NOT INSTANTIABLE FINAL",
		// Collection elements (PLS-00103): PERSISTABLE only after a
		// parenthesized element, NOT NULL inside the parentheses, once each,
		// and no object type modifiers.
		"CREATE TYPE t AS TABLE OF NUMBER NOT FINAL",
		"CREATE TYPE t AS TABLE OF u NOT PERSISTABLE",
		"CREATE TYPE t AS TABLE OF NUMBER PERSISTABLE",
		"CREATE TYPE t AS VARRAY(3) OF NUMBER NOT PERSISTABLE",
		"CREATE TYPE t AS TABLE OF (NUMBER) NOT NULL",
		"CREATE TYPE t AS TABLE OF (NUMBER) PERSISTABLE NOT NULL",
		"CREATE TYPE t AS TABLE OF ((NUMBER))",
		"CREATE TYPE t AS TABLE OF (NUMBER) NOT PERSISTABLE NOT PERSISTABLE",
		"CREATE TYPE t AS TABLE OF VARCHAR2(10) NOT NULL CHARACTER SET NCHAR_CS",
		"CREATE TYPE t AS TABLE OF REF (u)",
		"CREATE TYPE t AS VARYING (3) OF NUMBER",
		// Attribute datatypes (PLS-00103, PLS-00551, PLS-00550 for a REF
		// with a character set).
		"CREATE TYPE t AS OBJECT (a REF u CHARACTER SET NCHAR_CS)",
		"CREATE TYPE t AS OBJECT (a REF)",
		"CREATE TYPE t AS OBJECT (a DOUBLE)",
		"CREATE TYPE t AS OBJECT (a LONG VARCHAR)",
		"CREATE TYPE t AS OBJECT (a VARCHAR2(10) CHARACTER SET ANY_CS)",
		// A subtype's MAP method overrides (PLS-00636, PLS-00646).
		"CREATE TYPE c UNDER p (MAP MEMBER FUNCTION m RETURN NUMBER)",
		// A method signature once per kind (PLS-00305).
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER, MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE c UNDER p (MEMBER FUNCTION f RETURN NUMBER, MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'A', MEMBER FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'B')",
		// A MAP or ORDER method is a MEMBER one.
		"CREATE TYPE t AS OBJECT (a NUMBER, MAP MEMBER FUNCTION f RETURN NUMBER, MEMBER FUNCTION f RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, ORDER MEMBER FUNCTION f (o t) RETURN NUMBER, MEMBER FUNCTION f (o t) RETURN NUMBER)",
		// A constructor is FINAL and overrides nothing (PLS-00169).
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT FINAL CONSTRUCTOR FUNCTION t RETURN SELF AS RESULT) NOT FINAL",
		"CREATE TYPE t AS OBJECT (a NUMBER, OVERRIDING CONSTRUCTOR FUNCTION t RETURN SELF AS RESULT) NOT FINAL",
		// A type with a NOT INSTANTIABLE method is NOT INSTANTIABLE
		// (PLS-00634).
		"CREATE TYPE t AS OBJECT (a NUMBER, NOT INSTANTIABLE MEMBER FUNCTION f RETURN NUMBER) NOT FINAL",
		"CREATE TYPE c UNDER p (NOT INSTANTIABLE MEMBER FUNCTION f RETURN NUMBER) NOT FINAL",
		// SHARING = DATA is not for types (ORA-65021).
		"CREATE TYPE t SHARING = DATA AS OBJECT (a NUMBER)",
		// An ORDER method's parameter has IN mode (PLS-00528).
		"CREATE TYPE t AS OBJECT (a NUMBER, ORDER MEMBER FUNCTION o (x OUT t) RETURN NUMBER)",
		// A method bears no attribute's name (PLS-00305).
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION a RETURN NUMBER)",
		`CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION "A" RETURN NUMBER)`,
		"CREATE TYPE c UNDER p (a NUMBER, MEMBER PROCEDURE a)",
		// A default or a function property makes no overload (PLS-00305).
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER PROCEDURE p (x NUMBER), MEMBER PROCEDURE p (x NUMBER := 1))",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER, MEMBER FUNCTION f RETURN NUMBER DETERMINISTIC)",
		// A type has a name (ORA-02302).
		"CREATE TYPE",
		"CREATE TYPE;",
		// Attribute names are unique in the statement (PLS-00410).
		"CREATE TYPE t AS OBJECT (a NUMBER, a DATE)",
		`CREATE TYPE t AS OBJECT (a NUMBER, "A" DATE)`,
		"CREATE TYPE c UNDER p (a NUMBER, a DATE)",
	} {
		ParseShouldFail(t, sql)
	}
}

// TestCreateTypeHeadingAST checks what the heading clauses, a subtype, and a
// collection element record.
func TestCreateTypeHeadingAST(t *testing.T) {
	ty := parseOne[*ast.CreateTypeStmt](t, "CREATE OR REPLACE TYPE s.t SHARING = NONE FORCE OID '1234567890abcdef1234567890ABCDEF' AUTHID CURRENT_USER ACCESSIBLE BY (PACKAGE k) AS OBJECT (a NUMBER, b REF s.u) NOT FINAL")
	if ty.Sharing != "NONE" || !ty.Force || ty.OID != "1234567890abcdef1234567890ABCDEF" || ty.AuthID != "CURRENT_USER" || ty.AccessibleBy.Len() != 1 {
		t.Errorf("heading = sharing %q force %v oid %q authid %q accessible %v", ty.Sharing, ty.Force, ty.OID, ty.AuthID, ty.AccessibleBy)
	}
	if b := ty.Attributes.Items[1].(*ast.ColumnDef).TypeName; !b.IsRef || b.Names.Len() != 2 {
		t.Errorf("attribute b = REF %v, %d names; want a REF to s.u", b.IsRef, b.Names.Len())
	}

	sub := parseOne[*ast.CreateTypeStmt](t, "CREATE TYPE t UNDER s.base (c DATE) NOT FINAL")
	if sub.Under == nil || sub.Under.Schema != "S" || sub.Under.Name != "BASE" || sub.Attributes.Len() != 1 ||
		!reflect.DeepEqual(sub.Modifiers, []string{"NOT FINAL"}) {
		t.Errorf("subtype = under %v, %d attributes, modifiers %q", sub.Under, sub.Attributes.Len(), sub.Modifiers)
	}

	coll := parseOne[*ast.CreateTypeStmt](t, "CREATE TYPE t AS TABLE OF (VARCHAR2(10) NOT NULL) NOT PERSISTABLE")
	if coll.AsTable == nil || !coll.ElementNotNull || !reflect.DeepEqual(coll.Modifiers, []string{"NOT PERSISTABLE"}) {
		t.Errorf("collection = %v, not null %v, modifiers %q", coll.AsTable, coll.ElementNotNull, coll.Modifiers)
	}

	if incomplete := parseOne[*ast.CreateTypeStmt](t, "CREATE TYPE t OID '1234567890ABCDEF1234567890ABCDEF'"); incomplete.OID != "1234567890ABCDEF1234567890ABCDEF" ||
		incomplete.Attributes != nil || incomplete.AsTable != nil || incomplete.AsVarray != nil || incomplete.Under != nil {
		t.Errorf("incomplete type = %+v", incomplete)
	}
}

// TestANSIDatatypes checks the ANSI spellings of the character and floating
// types, and REF types, in SQL and PL/SQL against Oracle 23ai.
func TestANSIDatatypes(t *testing.T) {
	both := []string{
		"CHARACTER", "CHARACTER(1 BYTE)", "CHARACTER(2 CHAR)", "CHARACTER VARYING(3)", "CHARACTER VARYING(3 CHAR)",
		"CHAR VARYING(3)", "NATIONAL CHARACTER(2)", "NATIONAL CHAR(2)", "NATIONAL CHARACTER VARYING(3)",
		"NATIONAL CHAR VARYING(3)", "NCHAR VARYING(3)", "NATIONAL CHARACTER", "DOUBLE PRECISION", "REF u", "REF s.u",
	}
	for _, datatype := range both {
		ParseAndCheck(t, "CREATE TABLE t (x "+datatype+")")
		ParseAndCheck(t, "DECLARE x "+datatype+"; BEGIN NULL; END;")
	}
	// Each context has spellings of its own.
	ParseAndCheck(t, "CREATE TABLE t (x LONG VARCHAR)")
	ParseShouldFail(t, "DECLARE x LONG VARCHAR; BEGIN NULL; END;")
	ParseAndCheck(t, "DECLARE x CHARACTER LARGE OBJECT; BEGIN NULL; END;")
	ParseShouldFail(t, "CREATE TABLE t (x CHARACTER LARGE OBJECT)")
	ParseAndCheck(t, "DECLARE x LONG RAW(32767); y LONG(100); BEGIN NULL; END;")
	ParseShouldFail(t, "CREATE TABLE t (x LONG RAW(10))")
	ParseShouldFail(t, "CREATE TABLE t (x LONG(10))")
	// DOUBLE takes PRECISION (ORA-02000, PLS-00103).
	ParseShouldFail(t, "CREATE TABLE t (x DOUBLE)")
	ParseShouldFail(t, "DECLARE x DOUBLE; BEGIN NULL; END;")
	// A VARRAY limit is any expression; a literal one is from 1 through
	// 2147483647 (PLS-00537, PLS-00325).
	ParseAndCheck(t, "CREATE TYPE t AS ARRAY(1+1) OF NUMBER")
	ParseAndCheck(t, "CREATE TYPE t AS VARRAY(2147483647) OF NUMBER")
	ParseAndCheck(t, "DECLARE c CONSTANT PLS_INTEGER := 3; TYPE a IS VARRAY(c) OF NUMBER; TYPE b IS VARRAY(1+1) OF NUMBER; BEGIN NULL; END;")
	// A numeric literal rounds half away from zero first.
	for _, limit := range []string{"ARRAY(0.5)", "ARRAY(0.9)", "ARRAY(2147483647.4)", "ARRAY(1e1)"} {
		ParseAndCheck(t, "CREATE TYPE t AS "+limit+" OF NUMBER")
	}
	ParseShouldFail(t, "DECLARE TYPE a IS VARRAY(0.4) OF NUMBER; BEGIN NULL; END;")
	// AUTHID takes its values quoted too.
	ParseAndCheck(t, `CREATE TYPE t AUTHID "DEFINER" AS OBJECT (a NUMBER)`)
	ParseAndCheck(t, `CREATE TYPE t AUTHID "CURRENT_USER" AS OBJECT (a NUMBER)`)
	for _, limit := range []string{
		"ARRAY(0)", "ARRAY(-1)", "ARRAY(2147483648)", "VARRAY(0)", "VARRAY(2147483648)", "VARYING ARRAY(0)",
		"ARRAY(0.0)", "ARRAY(-1.0)", "ARRAY(2147483648.0)", "ARRAY(0.4)", "ARRAY(-0.4)", "ARRAY(2147483647.5)", "ARRAY(1e10)",
	} {
		ParseShouldFail(t, "CREATE TYPE t AS "+limit+" OF NUMBER")
	}
	ParseShouldFail(t, "DECLARE TYPE a IS VARRAY(0) OF NUMBER; BEGIN NULL; END;")
	ParseShouldFail(t, "DECLARE TYPE a IS VARRAY(-1) OF NUMBER; BEGIN NULL; END;")
	// Overloads that differ in a quoted spelling of an ANSI type alone are
	// one (PLS-00305), even with a type of that name.
	ParseShouldFail(t, `CREATE TYPE t AS OBJECT (a NUMBER, MEMBER PROCEDURE p (x CHARACTER), MEMBER PROCEDURE p (x "CHARACTER"))`)
	// A method's SELF written in the mode of its implicit one is that one
	// (PLS-00305): IN for a member function, IN OUT for a member procedure
	// or a constructor, whose SELF is IN OUT (PLS-00660); a static method
	// has none (PLS-00586).
	for _, methods := range []string{
		"MEMBER FUNCTION f RETURN NUMBER, MEMBER FUNCTION f (SELF IN t) RETURN NUMBER",
		"MEMBER FUNCTION f RETURN NUMBER, MEMBER FUNCTION f (SELF t) RETURN NUMBER",
		"MEMBER FUNCTION f (x NUMBER) RETURN NUMBER, MEMBER FUNCTION f (SELF IN t, x NUMBER) RETURN NUMBER",
		"MEMBER PROCEDURE p, MEMBER PROCEDURE p (SELF IN OUT t)",
		"MAP MEMBER FUNCTION m RETURN NUMBER, MEMBER FUNCTION m (SELF IN t) RETURN NUMBER",
		"CONSTRUCTOR FUNCTION t RETURN SELF AS RESULT, CONSTRUCTOR FUNCTION t (SELF IN OUT t) RETURN SELF AS RESULT",
		"CONSTRUCTOR FUNCTION t (SELF IN t) RETURN SELF AS RESULT",
		"STATIC FUNCTION f (SELF IN t) RETURN NUMBER",
	} {
		ParseShouldFail(t, "CREATE TYPE t AS OBJECT (a NUMBER, "+methods+")")
	}
	ParseAndCheck(t, "CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER, MEMBER FUNCTION f (SELF IN OUT t) RETURN NUMBER, MEMBER PROCEDURE p, MEMBER PROCEDURE p (SELF IN t))")
	ParseShouldFail(t, "ALTER TYPE t ADD MEMBER FUNCTION f RETURN NUMBER, ADD MEMBER FUNCTION f (SELF IN t) RETURN NUMBER CASCADE")
	ParseShouldFail(t, "ALTER TYPE t ADD STATIC FUNCTION f (SELF IN t) RETURN NUMBER CASCADE")
	// A SELF is of its type, no predefined datatype (PLS-00382).
	for _, methods := range []string{
		"CONSTRUCTOR FUNCTION t (SELF IN OUT NUMBER) RETURN SELF AS RESULT",
		"MEMBER FUNCTION f (SELF IN NUMBER) RETURN NUMBER",
		"MEMBER PROCEDURE p (SELF IN OUT VARCHAR2)",
	} {
		ParseShouldFail(t, "CREATE TYPE t AS OBJECT (a NUMBER, "+methods+")")
	}
	ParseAndCheck(t, "CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f (SELF IN s.t) RETURN NUMBER)")
	ParseShouldFail(t, "ALTER TYPE t ADD MEMBER FUNCTION f (SELF IN DATE) RETURN NUMBER CASCADE")
	// A VARRAY limit past any exponent bound fails before it is built
	// (PLS-00569, PLS-00537).
	ParseShouldFail(t, "CREATE TYPE t AS VARRAY(1e1000000000) OF NUMBER")
	ParseShouldFail(t, "CREATE TYPE t AS VARRAY(1e-1000000000) OF NUMBER")
	// ALTER TABLE ADD lists check their REF constraints against their own
	// columns (ORA-22893, ORA-22888).
	ParseShouldFail(t, "ALTER TABLE r ADD (n NUMBER, SCOPE FOR (n) IS o)")
	ParseShouldFail(t, "ALTER TABLE r ADD (n NUMBER, REF (n) WITH ROWID)")
	ParseShouldFail(t, "ALTER TABLE r ADD (c REF t SCOPE IS o, SCOPE FOR (c) IS o)")
	ParseAndCheck(t, "ALTER TABLE r ADD (c REF t, SCOPE FOR (c) IS o)")
	// A REF column takes a CHECK constraint.
	ParseAndCheck(t, "CREATE TABLE r (c REF t CHECK (c IS NOT NULL), d REF t CONSTRAINT ck CHECK (1 = 1))")
	// A column named CONSTRAINT takes a reserved datatype word alone.
	for _, datatype := range []string{"NUMBER", "CHAR(1)", "VARCHAR2(10)", "VARCHAR(10)", "DATE", "INTEGER", "SMALLINT", "DECIMAL", "FLOAT", "LONG", "LONG RAW", "RAW(1)", "ROWID"} {
		ParseAndCheck(t, "CREATE TABLE t (constraint "+datatype+")")
	}
	for _, datatype := range []string{"CLOB", "BLOB", "NCLOB", "NCHAR(1)", "NVARCHAR2(1)", "TIMESTAMP", "INTERVAL DAY TO SECOND", "JSON", "REF u", "CHARACTER(1)", "DEC(5)", "DOUBLE PRECISION", "NATIONAL CHAR(1)", "NCHAR VARYING(3)"} {
		ParseShouldFail(t, "CREATE TABLE t (constraint "+datatype+")")
	}
	// A REF scope or referenced table names its table after a schema
	// (ORA-00903), and a literal past int64 is no PL/SQL length (PLS-00325).
	ParseShouldFail(t, "CREATE TABLE r (c REF t SCOPE IS s.)")
	ParseShouldFail(t, "CREATE TABLE r (c REF t, SCOPE FOR (c) IS s.)")
	ParseShouldFail(t, "CREATE TABLE r (c REF t REFERENCES s.)")
	ParseShouldFail(t, "CREATE TYPE t AS OBJECT (a VARCHAR2(9223372036854775808))")
	ParseShouldFail(t, "DECLARE x NUMBER(99999999999999999999); BEGIN NULL; END;")
	// ALTER TABLE ADD takes out-of-line REF constraints, alone or in a list
	// with columns and other constraints.
	for _, sql := range []string{
		"ALTER TABLE r ADD SCOPE FOR (c) IS o", "ALTER TABLE r ADD REF (d) WITH ROWID",
		"ALTER TABLE r ADD (SCOPE FOR (d) IS o)", "ALTER TABLE r ADD (REF (c) WITH ROWID)",
		"ALTER TABLE r ADD (SCOPE FOR (c) IS o, REF (d) WITH ROWID)", "ALTER TABLE r ADD (m NUMBER, SCOPE FOR (d) IS o)",
	} {
		ParseAndCheck(t, sql)
	}
	cmds := parseOne[*ast.AlterTableStmt](t, "ALTER TABLE r ADD (m NUMBER, CONSTRAINT pk PRIMARY KEY (n), SCOPE FOR (d) IS o)").Actions
	if cmds.Len() != 3 || cmds.Items[0].(*ast.AlterTableCmd).ColumnDef == nil ||
		cmds.Items[1].(*ast.AlterTableCmd).Constraint.Type != ast.CONSTRAINT_PRIMARY ||
		cmds.Items[2].(*ast.AlterTableCmd).Constraint.Type != ast.CONSTRAINT_SCOPE {
		t.Errorf("ALTER TABLE ADD list = %s", ast.NodeToString(cmds))
	}
	ParseShouldFail(t, "ALTER TABLE r ADD SCOPE FOR (c) IS o@lnk")
	// A MAP method and a STATIC one may share a signature.
	ParseAndCheck(t, "CREATE TYPE t AS OBJECT (a NUMBER, MAP MEMBER FUNCTION f RETURN NUMBER, STATIC FUNCTION f RETURN NUMBER)")
	// A REF targets no predefined datatype (ORA-02302, PLS-00532) but one a
	// declaration in scope redeclares.
	for _, datatype := range []string{"REF NUMBER", "REF DATE", "REF VARCHAR2", "REF BLOB"} {
		ParseShouldFail(t, "CREATE TABLE t (x "+datatype+")")
	}
	for _, datatype := range []string{"REF DATE", "REF PLS_INTEGER", "REF SYS_REFCURSOR"} {
		ParseShouldFail(t, "DECLARE x "+datatype+"; BEGIN NULL; END;")
	}
	ParseShouldFail(t, "CREATE TYPE t AS OBJECT (x REF NUMBER)")
	ParseAndCheck(t, "CREATE TABLE t (x REF s.number_t, y REF date_t)")
	// A SQL size is 1 or more (ORA-01723) and at most 2000 for a fixed-length
	// type, 32767 for a varying one or RAW (ORA-00910).
	ParseAndCheck(t, "CREATE TABLE t (a CHAR(2000), b NCHAR(1000), c CHAR(2000 CHAR), d VARCHAR2(1), e RAW(1))")
	for _, datatype := range []string{
		"CHAR(0)", "CHARACTER(0)", "CHARACTER VARYING(0)", "NATIONAL CHARACTER VARYING(0)", "VARCHAR2(0)", "NCHAR(0)",
		"NVARCHAR2(0)", "RAW(0)", "CHAR(2001)", "CHAR(2001 CHAR)", "NCHAR(2001)", "VARCHAR2(32768)", "VARCHAR2(32768 CHAR)",
		"NVARCHAR2(32768)", "RAW(32768)",
	} {
		ParseShouldFail(t, "CREATE TABLE t (x "+datatype+")")
		ParseShouldFail(t, "SELECT CAST('a' AS "+datatype+") FROM dual")
	}
	// DEC is DECIMAL's ANSI spelling, in SQL and PL/SQL.
	ParseAndCheck(t, "CREATE TABLE t (x DEC(10,2), y DEC, z DEC(5), w DEC(*,2))")
	ParseAndCheck(t, "DECLARE x DEC(10,2); BEGIN NULL; END;")
	// SQL bounds a NUMBER precision to 1 through 38 (ORA-01727) and its
	// scale to -84 through 127 (ORA-01728), a FLOAT precision to 1 through
	// 126 (ORA-01724), and a datetime or interval precision to 0 through 9
	// (ORA-30088).
	ParseAndCheck(t, "CREATE TABLE t (a NUMBER(38), b DEC(1), c NUMBER(5,-84), d NUMERIC(5,127), e NUMBER(*,127), f FLOAT(126), g FLOAT(1))")
	ParseAndCheck(t, "CREATE TABLE t (a TIMESTAMP(9), b TIMESTAMP(0) WITH TIME ZONE, c INTERVAL DAY(9) TO SECOND(9), d INTERVAL YEAR(0) TO MONTH)")
	// PL/SQL, an object type's attributes, and a collection's elements bound
	// a literal length to 1 through 32767 (PLS-00215, LONG and LONG RAW
	// among them), a NUMBER precision and scale (PLS-00216, PLS-00217), a
	// FLOAT precision (PLS-00623), and a datetime or interval precision
	// (PLS-00260); a type a declaration shadows is the engine's.
	ParseAndCheck(t, "DECLARE a LONG(32767); b LONG RAW(32767); c VARCHAR2(32767 CHAR); d NUMBER(38, -84); e FLOAT(126); f TIMESTAMP(9); g INTERVAL DAY(9) TO SECOND(9); BEGIN NULL; END;")
	ParseAndCheck(t, "DECLARE SUBTYPE number IS VARCHAR2(10); x number(100); BEGIN NULL; END;")
	ParseAndCheck(t, "DECLARE x NUMBER(5.5); y VARCHAR2(1e1); BEGIN NULL; END;")
	for _, decl := range []string{
		"LONG(0)", "LONG RAW(0)", "VARCHAR2(32768)", "VARCHAR2(0)", "NVARCHAR2(0)", "CHAR(32768)", "STRING(32768)",
		"RAW(32768)", "VARCHAR2(32768 CHAR)", "NUMBER(0)", "NUMBER(5, -85)", "DECIMAL(39)", "FLOAT(0)",
		"TIMESTAMP(10)", "INTERVAL DAY(10) TO SECOND", "INTERVAL DAY TO SECOND(10)", "INTERVAL YEAR(10) TO MONTH",
	} {
		ParseShouldFail(t, "DECLARE x "+decl+"; BEGIN NULL; END;")
	}
	ParseShouldFail(t, "CREATE TYPE t AS OBJECT (a NUMBER(39))")
	ParseShouldFail(t, "CREATE TYPE t AS VARRAY(2) OF TIMESTAMP(10)")
	for _, datatype := range []string{
		"NUMBER(0)", "NUMBER(39)", "DEC(0)", "DEC(39)", "DECIMAL(0)", "NUMERIC(39)", "NUMBER(5,-85)", "NUMBER(5,128)",
		"NUMBER(*,128)", "DECIMAL(5,128)", "FLOAT(0)", "FLOAT(127)", "TIMESTAMP(10)", "TIMESTAMP(10) WITH TIME ZONE",
		"INTERVAL YEAR(10) TO MONTH", "INTERVAL DAY(9) TO SECOND(10)",
	} {
		ParseShouldFail(t, "CREATE TABLE t (x "+datatype+")")
		ParseShouldFail(t, "SELECT CAST(NULL AS "+datatype+") FROM dual")
	}
	// [ NATIONAL ] CHARACTER LARGE OBJECT and NCHAR LARGE OBJECT are PL/SQL's
	// (ORA-03062 in SQL); NATIONAL CHAR LARGE OBJECT is no type (PLS-00103).
	ParseAndCheck(t, "DECLARE x NATIONAL CHARACTER LARGE OBJECT; y NCHAR LARGE OBJECT; BEGIN NULL; END;")
	ParseShouldFail(t, "DECLARE x NATIONAL CHAR LARGE OBJECT; BEGIN NULL; END;")
	ParseShouldFail(t, "CREATE TABLE t (x NATIONAL CHARACTER LARGE OBJECT)")
	// NATIONAL names no item for %TYPE (PLS-00103), and DOUBLE PRECISION
	// is numeric in a package body too (PLS-00550).
	ParseShouldFail(t, "DECLARE national NUMBER; x national%TYPE; BEGIN NULL; END;")
	ParseShouldFail(t, "CREATE OR REPLACE PACKAGE BODY k AS x DOUBLE PRECISION CHARACTER SET NCHAR_CS; END;")
	// A VARYING type has a size in SQL (ORA-00906, ORA-00910 for empty
	// parentheses), not in PL/SQL, an attribute, or a collection element.
	for _, datatype := range []string{"CHARACTER VARYING", "CHAR VARYING", "NCHAR VARYING", "NATIONAL CHARACTER VARYING"} {
		ParseShouldFail(t, "CREATE TABLE t (x "+datatype+")")
		ParseShouldFail(t, "CREATE TABLE t (x "+datatype+"())")
		ParseAndCheck(t, "DECLARE x "+datatype+"; BEGIN NULL; END;")
	}
	// A REF column takes SCOPE IS and WITH ROWID, inline or out of line,
	// among its other constraints; another column takes neither
	// (ORA-22893), and one SCOPE (ORA-22888).
	for _, cols := range []string{
		"c REF t SCOPE IS o", "c REF t WITH ROWID", "c REF t SCOPE IS s.o NOT NULL",
		"c REF t NOT NULL SCOPE IS o", "c REF t DEFAULT NULL SCOPE IS o", "c REF t WITH ROWID SCOPE IS o",
		"c REF t SCOPE IS o WITH ROWID", "c REF t CONSTRAINT fk REFERENCES o ON DELETE CASCADE",
		"c REF t, SCOPE FOR (c) IS o, d NUMBER", "c REF t, REF (c) WITH ROWID",
		"SCOPE FOR (c) IS o, c REF t", "o h, SCOPE FOR (o.r) IS x", "o h, REF (o.r) WITH ROWID",
		"c REF t WITH ROWID, REF (c) WITH ROWID", "c REF t REFERENCES o",
		"o h, FOREIGN KEY (o.r) REFERENCES x", "o h, CONSTRAINT k FOREIGN KEY (o.r) REFERENCES x",
		"c REF t WITH ROWID REFERENCES o", "c REF t, SCOPE FOR (c) IS o, REF (c) WITH ROWID",
		// An attribute path is not the quoted column that spells it.
		`o h, "O.R" NUMBER, SCOPE FOR (o.r) IS x`, `o h, "O.R" NUMBER, REF (o.r) WITH ROWID`,
		`o h, "O.R" REF t, SCOPE FOR (o.r) IS x, SCOPE FOR ("O.R") IS x`,
		`o h, "O.R" h, SCOPE FOR (o."R.X") IS x, SCOPE FOR ("O.R".x) IS x`,
		// An attribute is a REF where a SCOPE or WITH ROWID says so; a
		// scalar one takes a referenced column list and keys.
		"o h, FOREIGN KEY (o.n) REFERENCES u (a)", "o h, REF (o.r) WITH ROWID, FOREIGN KEY (o.r) REFERENCES x",
		"c REF t WITH ROWID, d NUMBER, PRIMARY KEY (d)", "o h, d NUMBER, FOREIGN KEY (o.n, d) REFERENCES p (a, b)",
	} {
		ParseAndCheck(t, "CREATE TABLE r ("+cols+")")
	}
	for _, cols := range []string{
		"c NUMBER SCOPE IS o", "c NUMBER WITH ROWID", "c REF t SCOPE o",
		"c REF t SCOPE IS o SCOPE IS o", "c REF t, SCOPE FOR (c) o", "c REF t, REF (c) ROWID",
		"c NUMBER, SCOPE FOR (c) IS o", "c NUMBER, REF (c) WITH ROWID",
		// One SCOPE, inline or out of line (ORA-22888), and a REF column's
		// referential constraint names no columns (ORA-30726).
		"c REF t SCOPE IS o, SCOPE FOR (c) IS o", "c REF t, SCOPE FOR (c) IS o, SCOPE FOR (c) IS o",
		"c REF t REFERENCES o(n)", "c REF t, FOREIGN KEY (c) REFERENCES o(n)",
		// The scope table is local (ORA-25124).
		"c REF t SCOPE IS o@lnk", "c REF t, SCOPE FOR (c) IS o@lnk",
		`o h, "O.R" NUMBER, SCOPE FOR ("O.R") IS x`, "o h, SCOPE FOR (o.r) IS x, SCOPE FOR (o.r) IS x",
		// SCOPE is a keyword there (ORA-03062, ORA-00902).
		`c REF t "SCOPE" IS o`, `c REF t, "SCOPE" FOR (c) IS o`, `c REF t WITH "ROWID"`,
		// Not both a SCOPE and a referential constraint, whatever the tables
		// (ORA-22896), and one referential constraint (ORA-30771).
		"c REF t SCOPE IS o REFERENCES o", "c REF t SCOPE IS o REFERENCES p", "c REF t REFERENCES p SCOPE IS o",
		"c REF t, SCOPE FOR (c) IS o, FOREIGN KEY (c) REFERENCES p", "c REF t REFERENCES p, SCOPE FOR (c) IS o",
		"SCOPE FOR (c) IS o, c REF t REFERENCES p", "c REF t SCOPE IS o, FOREIGN KEY (c) REFERENCES o",
		"o h, SCOPE FOR (o.r) IS x, FOREIGN KEY (o.r) REFERENCES y",
		"c REF t REFERENCES o REFERENCES p", "c REF t REFERENCES o, FOREIGN KEY (c) REFERENCES p",
		// An attribute's column is of an object type (ORA-02337).
		"n NUMBER, SCOPE FOR (n.x) IS o", "n NUMBER, FOREIGN KEY (n.x) REFERENCES o", "n NUMBER, REF (n.x) WITH ROWID",
		"c REF t, SCOPE FOR (c.n) IS o", "v VARCHAR2(10), SCOPE FOR (v.x) IS o", "d DATE, REF (d.x) WITH ROWID",
		"n NUMBER, d NUMBER, FOREIGN KEY (n.x, d) REFERENCES p (a, b)",
		// A REF is in no key (ORA-02329), a composite foreign key among them
		// (ORA-02267).
		"c REF t PRIMARY KEY", "c REF t UNIQUE", "c REF t, d NUMBER, PRIMARY KEY (c, d)", "c REF t, d NUMBER, UNIQUE (c, d)",
		"c REF t, d NUMBER, FOREIGN KEY (c, d) REFERENCES p (a, b)", "c REF t, d NUMBER, FOREIGN KEY (d, c) REFERENCES p",
		"o h, REF (o.r) WITH ROWID, d NUMBER, FOREIGN KEY (o.r, d) REFERENCES p (a, b)",
		// The REF rules hold for an attribute a WITH ROWID shows to be one.
		"o h, REF (o.r) WITH ROWID, FOREIGN KEY (o.r) REFERENCES x (n)",
		"o h, REF (o.r) WITH ROWID, FOREIGN KEY (o.r) REFERENCES x, FOREIGN KEY (o.r) REFERENCES y",
		// A referenced table is local (ORA-02021).
		"c REF t REFERENCES o@lnk", "c REF t, FOREIGN KEY (c) REFERENCES o@lnk",
		"n NUMBER REFERENCES p@lnk (a)", "n NUMBER, FOREIGN KEY (n) REFERENCES p@lnk (a)",
	} {
		ParseShouldFail(t, "CREATE TABLE r ("+cols+")")
	}
	ParseShouldFail(t, "ALTER TABLE r ADD CONSTRAINT k FOREIGN KEY (n) REFERENCES p@lnk (a)")
	ParseShouldFail(t, "ALTER TABLE r MODIFY (n REFERENCES p@lnk (a))")
	ParseShouldFail(t, "ALTER TABLE r ADD (m NUMBER REFERENCES p@lnk (a))")
	ParseAndCheck(t, "CREATE TYPE t AS OBJECT (x CHARACTER VARYING)")
	ParseAndCheck(t, "CREATE TYPE t AS TABLE OF CHAR VARYING")
	// A REF names an object type, optionally schema-qualified: no %TYPE or
	// %ROWTYPE (ORA-00911, PLS-00103) and no third component (ORA-03099).
	ParseShouldFail(t, "CREATE TABLE t (x REF u%TYPE)")
	ParseShouldFail(t, "CREATE TABLE t (x REF s.u.a)")
	ParseShouldFail(t, "DECLARE x REF u%TYPE; BEGIN NULL; END;")
	ParseShouldFail(t, "CREATE TYPE t AS OBJECT (x REF u%ROWTYPE)")

	tn := parseOne[*ast.CreateTableStmt](t, "CREATE TABLE t (x NATIONAL CHARACTER VARYING(3))").Columns.Items[0].(*ast.ColumnDef).TypeName
	if got := typeNameWords(tn); !reflect.DeepEqual(got, []string{"NATIONAL", "CHARACTER", "VARYING"}) || tn.TypeMods.Len() != 1 {
		t.Errorf("NATIONAL CHARACTER VARYING(3) = %q with %d modifiers", got, tn.TypeMods.Len())
	}
}

// TestAlterTypeElements checks the attribute and method lists of ALTER TYPE
// against Oracle 23ai: each item has its own ADD, DROP, or MODIFY, a list
// holds methods or attributes but not both, and a method may carry
// inheritance clauses. The rejected forms fail with ORA-06545.
func TestAlterTypeElements(t *testing.T) {
	for _, sql := range []string{
		"ALTER TYPE t MODIFY ATTRIBUTE b VARCHAR2(20) CASCADE",
		"ALTER TYPE t MODIFY ATTRIBUTE b VARCHAR2(ORA_MAX_NAME_LEN) CASCADE",
		"ALTER TYPE t ADD ATTRIBUTE c NUMBER CASCADE",
		"ALTER TYPE t ADD ATTRIBUTE (d NUMBER, e DATE) CASCADE",
		"ALTER TYPE t ADD ATTRIBUTE g REF u CASCADE",
		"ALTER TYPE t DROP ATTRIBUTE e CASCADE",
		"ALTER TYPE t DROP ATTRIBUTE (d, e) CASCADE",
		"ALTER TYPE t ADD ATTRIBUTE c NUMBER, MODIFY ATTRIBUTE b VARCHAR2(30) CASCADE",
		"ALTER TYPE t DROP ATTRIBUTE e, DROP ATTRIBUTE d CASCADE",
		"ALTER TYPE t ADD MEMBER FUNCTION g RETURN NUMBER, ADD MEMBER FUNCTION h RETURN NUMBER CASCADE",
		"ALTER TYPE t DROP MEMBER FUNCTION g RETURN NUMBER, DROP MEMBER FUNCTION h RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MEMBER FUNCTION g RETURN NUMBER, DROP MEMBER FUNCTION f RETURN NUMBER CASCADE",
		"ALTER TYPE s ADD OVERRIDING MEMBER FUNCTION f RETURN NUMBER CASCADE",
		"ALTER TYPE s DROP OVERRIDING MEMBER PROCEDURE p (a IN NUMBER, b IN OUT NOCOPY CLOB), DROP OVERRIDING MEMBER PROCEDURE q CASCADE",
		"ALTER TYPE t ADD NOT FINAL MEMBER FUNCTION nf RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD STATIC FUNCTION s RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MAP MEMBER FUNCTION m RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD CONSTRUCTOR FUNCTION t (a NUMBER) RETURN SELF AS RESULT CASCADE",
		"ALTER TYPE t ADD ORDER MEMBER FUNCTION o (x t) RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MAP MEMBER FUNCTION m RETURN NUMBER, ADD MEMBER FUNCTION f (x NUMBER) RETURN NUMBER CASCADE",
		"ALTER TYPE t DROP ATTRIBUTE (a, b) CASCADE",
		`ALTER TYPE t ADD ATTRIBUTE c NUMBER, ADD ATTRIBUTE "c" DATE CASCADE`,
		// Overloads differ in a parameter name or type, or the return type.
		"ALTER TYPE t ADD MEMBER FUNCTION f RETURN NUMBER, ADD MEMBER FUNCTION f (x NUMBER) RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MEMBER FUNCTION f (x NUMBER) RETURN NUMBER, ADD MEMBER FUNCTION f (y NUMBER) RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MEMBER FUNCTION f RETURN NUMBER, ADD MEMBER FUNCTION f RETURN VARCHAR2 CASCADE",
		`ALTER TYPE t ADD MEMBER PROCEDURE p, ADD MEMBER PROCEDURE "p" CASCADE`,
		`CREATE TYPE t AS OBJECT (a NUMBER, "a" DATE)`,
		// Overloads may differ in a parameter's mode.
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER PROCEDURE p (x IN NUMBER), MEMBER PROCEDURE p (x OUT NUMBER))",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER PROCEDURE p (x NUMBER), MEMBER PROCEDURE p (x IN OUT NOCOPY NUMBER))",
		"ALTER TYPE t ADD MEMBER PROCEDURE p (x IN NUMBER), ADD MEMBER PROCEDURE p (x OUT NUMBER) CASCADE",
		"ALTER TYPE t ADD ORDER MEMBER FUNCTION o (x IN t) RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD FINAL CONSTRUCTOR FUNCTION t RETURN SELF AS RESULT CASCADE",
		"ALTER TYPE t DROP CONSTRUCTOR FUNCTION t (x NUMBER) RETURN SELF AS RESULT CASCADE",
		// A call specification, as in CREATE TYPE.
		"ALTER TYPE t ADD MEMBER FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int' CASCADE",
		"ALTER TYPE t ADD MEMBER FUNCTION f RETURN NUMBER IS LANGUAGE JAVA NAME 'X.f() return int', ADD MEMBER FUNCTION g RETURN NUMBER CASCADE",
		// A collection's new element type, as in CREATE TYPE AS TABLE OF.
		"ALTER TYPE v MODIFY ELEMENT TYPE CHARACTER VARYING CASCADE",
		"ALTER TYPE v MODIFY ELEMENT TYPE VARCHAR2 CASCADE",
		"CREATE TYPE c UNDER p (OVERRIDING MAP MEMBER FUNCTION m RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER, MEMBER FUNCTION f RETURN VARCHAR2)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f (x NUMBER) RETURN NUMBER, MEMBER FUNCTION f (x VARCHAR2) RETURN NUMBER)",
		"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER PROCEDURE p, STATIC PROCEDURE p)",
		"CREATE TYPE c UNDER p (NOT INSTANTIABLE MEMBER FUNCTION f RETURN NUMBER) NOT INSTANTIABLE NOT FINAL",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		"ALTER TYPE t ADD MEMBER FUNCTION g RETURN NUMBER, ADD ATTRIBUTE d NUMBER CASCADE",
		"ALTER TYPE t ADD ATTRIBUTE e NUMBER, ADD MEMBER FUNCTION k RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MEMBER FUNCTION m RETURN NUMBER, MEMBER FUNCTION n RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD ATTRIBUTE",
		"ALTER TYPE t MODIFY ATTRIBUTE (b VARCHAR2(30) CASCADE",
		"ALTER TYPE t ADD MAP FUNCTION m RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MEMBER g RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MEMBER PROCEDURE p RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MEMBER FUNCTION g CASCADE",
		"ALTER TYPE t ADD CONSTRUCTOR FUNCTION t (a NUMBER) RETURN SELF CASCADE",
		"ALTER TYPE t ADD CONSTRUCTOR PROCEDURE p CASCADE",
		"ALTER TYPE t ADD NOT OVERRIDING MEMBER FUNCTION no RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD FINAL OVERRIDING MEMBER FUNCTION fo RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD OVERRIDING ATTRIBUTE z NUMBER CASCADE",
		// Inheritance clauses that conflict with each other or the method
		// kind (ORA-22324 or ORA-06545 on the PLS-00169 rule).
		"ALTER TYPE t ADD OVERRIDING STATIC FUNCTION g RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD NOT INSTANTIABLE STATIC FUNCTION h RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD NOT INSTANTIABLE FINAL MEMBER FUNCTION k RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD NOT INSTANTIABLE CONSTRUCTOR FUNCTION t (x NUMBER) RETURN SELF AS RESULT CASCADE",
		"ALTER TYPE t ADD NOT FINAL CONSTRUCTOR FUNCTION t RETURN SELF AS RESULT CASCADE",
		"ALTER TYPE t ADD OVERRIDING CONSTRUCTOR FUNCTION t RETURN SELF AS RESULT CASCADE",
		"ALTER TYPE t ADD ORDER MEMBER FUNCTION o (x OUT t) RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD ORDER MEMBER FUNCTION o (x IN OUT t) RETURN NUMBER CASCADE",
		// MODIFY ELEMENT TYPE widens a scalar element (ORA-22324).
		"ALTER TYPE v MODIFY ELEMENT TYPE REF t CASCADE",
		"ALTER TYPE v MODIFY ELEMENT TYPE SYS_REFCURSOR CASCADE",
		// No EXTERNAL call spec and no PL/SQL body (ORA-06545).
		`ALTER TYPE t ADD MEMBER FUNCTION f RETURN NUMBER AS EXTERNAL NAME "f" LIBRARY lib CASCADE`,
		"ALTER TYPE t ADD MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; CASCADE",
		// A constructor bears its type's name (ORA-22324 on PLS-00658).
		"ALTER TYPE t ADD CONSTRUCTOR FUNCTION wrong (x NUMBER) RETURN SELF AS RESULT CASCADE",
		"ALTER TYPE t ADD MEMBER FUNCTION g RETURN NUMBER, ADD CONSTRUCTOR FUNCTION wrong RETURN SELF AS RESULT CASCADE",
		// MAP takes no parameter besides SELF and ORDER exactly one, and a
		// list adds at most one of them (ORA-22324).
		"ALTER TYPE t ADD MAP MEMBER FUNCTION m (x NUMBER) RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD ORDER MEMBER FUNCTION o RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD ORDER MEMBER FUNCTION o (x t, y t) RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MAP MEMBER FUNCTION m RETURN NUMBER, ADD ORDER MEMBER FUNCTION o (x t) RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD ORDER MEMBER FUNCTION o (x t) RETURN NUMBER, ADD ORDER MEMBER FUNCTION o2 (x t) RETURN NUMBER CASCADE",
		// A function returns a datatype (ORA-06545).
		"ALTER TYPE t ADD MEMBER FUNCTION good RETURN NUMBER, ADD MEMBER FUNCTION bad RETURN, ADD MEMBER FUNCTION nx RETURN NUMBER CASCADE",
		// An attribute appears once in the statement (ORA-22324).
		"ALTER TYPE t MODIFY ATTRIBUTE a VARCHAR2(20), MODIFY ATTRIBUTE a VARCHAR2(30) CASCADE",
		"ALTER TYPE t MODIFY ATTRIBUTE (a VARCHAR2(20), a VARCHAR2(30)) CASCADE",
		"ALTER TYPE t DROP ATTRIBUTE a, DROP ATTRIBUTE a CASCADE",
		"ALTER TYPE t DROP ATTRIBUTE a, MODIFY ATTRIBUTE a VARCHAR2(30) CASCADE",
		`ALTER TYPE t ADD ATTRIBUTE c NUMBER, ADD ATTRIBUTE "C" DATE CASCADE`,
		// And a method signature once, whatever the verb or kind (ORA-22324).
		"ALTER TYPE t ADD MEMBER FUNCTION f RETURN NUMBER, ADD MEMBER FUNCTION f RETURN NUMBER CASCADE",
		"ALTER TYPE t DROP MEMBER FUNCTION g RETURN NUMBER, DROP MEMBER FUNCTION g RETURN NUMBER CASCADE",
		"ALTER TYPE t DROP MEMBER FUNCTION g RETURN NUMBER, ADD MEMBER FUNCTION g RETURN NUMBER CASCADE",
		"ALTER TYPE t ADD MEMBER PROCEDURE p, ADD STATIC PROCEDURE p CASCADE",
		`ALTER TYPE t ADD MEMBER PROCEDURE p, ADD MEMBER PROCEDURE "P" CASCADE`,
		"ALTER TYPE t ADD MEMBER FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'A', ADD MEMBER FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'B' CASCADE",
	} {
		ParseShouldFail(t, sql)
	}

	alter := parseOne[*ast.AlterTypeStmt](t, "ALTER TYPE s DROP OVERRIDING MEMBER PROCEDURE p (a NUMBER), ADD MAP MEMBER FUNCTION m RETURN NUMBER CASCADE")
	if alter.Action != "DROP_METHOD" || len(alter.Methods) != 2 || !alter.Cascade {
		t.Fatalf("ALTER TYPE = action %q, %d methods, cascade %v", alter.Action, len(alter.Methods), alter.Cascade)
	}
	if m := alter.Methods[0]; m.Action != "DROP" || !reflect.DeepEqual(m.Modifiers, []string{"OVERRIDING"}) || m.Kind != "MEMBER" ||
		m.Type != "PROCEDURE" || m.Name != "P" || len(m.Params) != 1 {
		t.Errorf("first method = %+v", m)
	}
	if m := alter.Methods[1]; m.Action != "ADD" || m.Kind != "MAP MEMBER" || m.Type != "FUNCTION" || m.Return == nil {
		t.Errorf("second method = %+v", m)
	}

	// A method's location starts at its verb.
	const items = "ALTER TYPE s DROP OVERRIDING MEMBER PROCEDURE p (a NUMBER), ADD MAP MEMBER FUNCTION m RETURN NUMBER CASCADE"
	for i, want := range []string{"DROP OVERRIDING MEMBER PROCEDURE p (a NUMBER)", "ADD MAP MEMBER FUNCTION m RETURN NUMBER"} {
		if m := parseOne[*ast.AlterTypeStmt](t, items).Methods[i]; items[m.Loc.Start:m.Loc.End] != want {
			t.Errorf("method %d spans %q, want %q", i, items[m.Loc.Start:m.Loc.End], want)
		}
	}

	const ctor = "ALTER TYPE t ADD CONSTRUCTOR FUNCTION t (a NUMBER) RETURN SELF AS RESULT"
	ret := parseOne[*ast.AlterTypeStmt](t, ctor).Methods[0].Return
	if got := ctor[ret.Loc.Start:ret.Loc.End]; got != "SELF AS RESULT" {
		t.Errorf("constructor return spans %q, want SELF AS RESULT", got)
	}

	attrs := parseOne[*ast.AlterTypeStmt](t, "ALTER TYPE t ADD ATTRIBUTE (c NUMBER, d DATE), MODIFY ATTRIBUTE b VARCHAR2(30), DROP ATTRIBUTE e")
	var got []string
	for _, attr := range attrs.Attributes {
		got = append(got, attr.Action+" "+attr.Name)
	}
	if want := []string{"ADD C", "ADD D", "MODIFY B", "DROP E"}; attrs.Action != "ADD_ATTRIBUTE" || !reflect.DeepEqual(got, want) {
		t.Errorf("ALTER TYPE attributes = %q %q, want ADD_ATTRIBUTE %q", attrs.Action, got, want)
	}
}

// typeNameWords returns the words of a datatype's name.
func typeNameWords(tn *ast.TypeName) []string {
	var words []string
	for _, item := range tn.Names.Items {
		words = append(words, item.(*ast.String).Str)
	}
	return words
}

// TestTypeAttributeAndElementTypes checks the predefined types Oracle 23ai
// refuses as an object type attribute (PLS-00530, PLS-00218, PLS-00989) and
// as a collection element (PLS-00531, PLS-00990), for every user, SYS
// included. The PL/SQL-only types are legal in a NOT PERSISTABLE type.
func TestTypeAttributeAndElementTypes(t *testing.T) {
	// Refused in a persistable type, legal in a NOT PERSISTABLE one. ALTER
	// TYPE does not show the type's persistability, so it takes them.
	for _, datatype := range []string{
		"LONG", "LONG RAW", "LONG(100)", "LONG RAW(32767)", "PLS_INTEGER", "BINARY_INTEGER",
		"NATURAL", "POSITIVE", "SIGNTYPE", "ROWID", "UROWID", "MLSLABEL",
	} {
		ParseShouldFail(t, "CREATE TYPE t AS OBJECT (x "+datatype+")")
		ParseShouldFail(t, "CREATE TYPE t AS OBJECT (x "+datatype+") PERSISTABLE")
		ParseShouldFail(t, "CREATE TYPE t AS TABLE OF "+datatype)
		ParseShouldFail(t, "CREATE TYPE t AS TABLE OF ("+datatype+") PERSISTABLE")
		ParseAndCheck(t, "CREATE TYPE t AS OBJECT (x "+datatype+") NOT PERSISTABLE")
		ParseAndCheck(t, "CREATE TYPE t AS OBJECT (x "+datatype+") NOT FINAL NOT PERSISTABLE")
		ParseAndCheck(t, "CREATE TYPE t AS TABLE OF ("+datatype+") NOT PERSISTABLE")
		ParseAndCheck(t, "CREATE TYPE t AS VARRAY(3) OF ("+datatype+") NOT PERSISTABLE")
		ParseAndCheck(t, "ALTER TYPE t ADD ATTRIBUTE x "+datatype+" CASCADE")
		// A subtype's persistability is its supertype's.
		ParseAndCheck(t, "CREATE TYPE t UNDER base (x "+datatype+")")
	}
	// Refused everywhere.
	for _, datatype := range []string{"JSON", "VECTOR", "SYS_REFCURSOR"} {
		ParseShouldFail(t, "CREATE TYPE t AS OBJECT (x "+datatype+") NOT PERSISTABLE")
		ParseShouldFail(t, "CREATE TYPE t AS TABLE OF ("+datatype+") NOT PERSISTABLE")
		ParseShouldFail(t, "CREATE TYPE t UNDER base (x "+datatype+")")
		ParseShouldFail(t, "ALTER TYPE t ADD ATTRIBUTE x "+datatype+" CASCADE")
	}
	// The NOT NULL subtypes are refused as attributes (PLS-00218) but not as
	// elements of a NOT PERSISTABLE collection; SIMPLE_FLOAT and
	// SIMPLE_DOUBLE are elements of any collection.
	for _, datatype := range []string{"SIMPLE_INTEGER", "NATURALN", "POSITIVEN", "SIMPLE_FLOAT", "SIMPLE_DOUBLE"} {
		ParseShouldFail(t, "CREATE TYPE t AS OBJECT (x "+datatype+") NOT PERSISTABLE")
		ParseShouldFail(t, "ALTER TYPE t ADD ATTRIBUTE x "+datatype+" CASCADE")
		ParseAndCheck(t, "CREATE TYPE t AS TABLE OF ("+datatype+") NOT PERSISTABLE")
	}
	for _, datatype := range []string{"SIMPLE_INTEGER", "NATURALN", "POSITIVEN"} {
		ParseShouldFail(t, "CREATE TYPE t AS TABLE OF "+datatype)
	}
	for _, datatype := range []string{"SIMPLE_FLOAT", "SIMPLE_DOUBLE"} {
		ParseAndCheck(t, "CREATE TYPE t AS TABLE OF "+datatype)
	}
	for _, datatype := range []string{
		"BOOLEAN", "BOOL", "CHARACTER LARGE OBJECT", "BINARY_FLOAT", "REAL", "INT", "DEC", "NUMERIC",
		"TIMESTAMP WITH LOCAL TIME ZONE", "XMLTYPE", "CLOB", "NCLOB", "BFILE",
	} {
		ParseAndCheck(t, "CREATE TYPE t AS OBJECT (x "+datatype+")")
		ParseAndCheck(t, "CREATE TYPE t AS TABLE OF "+datatype)
	}
	// ADD and MODIFY ATTRIBUTE name the datatype (ORA-06545).
	ParseShouldFail(t, "ALTER TYPE t ADD ATTRIBUTE c")
	ParseShouldFail(t, "ALTER TYPE t MODIFY ATTRIBUTE b")
	ParseShouldFail(t, "ALTER TYPE t ADD ATTRIBUTE (c)")
	// LONG takes one length (PLS-00572); a REF no character set (PLS-00550).
	ParseShouldFail(t, "DECLARE x LONG(10, 2); BEGIN NULL; END;")
	ParseShouldFail(t, "DECLARE x LONG RAW(10 BYTE); BEGIN NULL; END;")
	ParseShouldFail(t, "DECLARE x REF u CHARACTER SET NCHAR_CS; BEGIN NULL; END;")
	ParseAndCheck(t, "DECLARE x LONG(100); y LONG RAW(32767); BEGIN NULL; END;")
	// A national type takes no BYTE (ORA-00907, PLS-00639), and CHAR only
	// outside SQL.
	for _, datatype := range []string{"NCHAR(2 CHAR)", "NVARCHAR2(2 CHAR)", "NATIONAL CHARACTER(2 CHAR)", "NCHAR VARYING(2 CHAR)"} {
		ParseShouldFail(t, "CREATE TABLE t (x "+datatype+")")
		ParseAndCheck(t, "DECLARE x "+datatype+"; BEGIN NULL; END;")
		ParseAndCheck(t, "CREATE TYPE t AS OBJECT (x "+datatype+")")
	}
	for _, datatype := range []string{"NCHAR(2 BYTE)", "NVARCHAR2(2 BYTE)", "NATIONAL CHARACTER(2 BYTE)", "NATIONAL CHAR VARYING(2 BYTE)"} {
		ParseShouldFail(t, "CREATE TABLE t (x "+datatype+")")
		ParseShouldFail(t, "DECLARE x "+datatype+"; BEGIN NULL; END;")
		ParseShouldFail(t, "CREATE TYPE t AS OBJECT (x "+datatype+")")
	}
	// DOUBLE takes PRECISION even where a subtype or type named DOUBLE is in
	// scope (PLS-00103); only a quoted "DOUBLE" names it.
	ParseShouldFail(t, "DECLARE SUBTYPE double IS NUMBER; x double; BEGIN NULL; END;")
	ParseShouldFail(t, "CREATE PACKAGE BODY k AS x double; END;")
	ParseAndCheck(t, "DECLARE SUBTYPE double IS NUMBER; x \"DOUBLE\"; BEGIN NULL; END;")
}
