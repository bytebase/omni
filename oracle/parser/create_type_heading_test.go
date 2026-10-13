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
		"CREATE TYPE t SHARING = DATA AS OBJECT (a NUMBER)",
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
	// A VARYING type has a size in SQL (ORA-00906), not in PL/SQL, an
	// attribute, or a collection element.
	for _, datatype := range []string{"CHARACTER VARYING", "CHAR VARYING", "NCHAR VARYING", "NATIONAL CHARACTER VARYING"} {
		ParseShouldFail(t, "CREATE TABLE t (x "+datatype+")")
		ParseAndCheck(t, "DECLARE x "+datatype+"; BEGIN NULL; END;")
	}
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
