package ddl

import (
	"strings"
	"testing"
)

func TestCreateAccessMethod(t *testing.T) {
	c := New()

	// Built-in AMs should already exist.
	if am := c.LookupAccessMethod("btree"); am == nil {
		t.Fatal("expected btree to exist")
	}
	if am := c.LookupAccessMethod("hash"); am == nil {
		t.Fatal("expected hash to exist")
	}
	if am := c.LookupAccessMethod("gist"); am == nil {
		t.Fatal("expected gist to exist")
	}

	// Create a custom access method.
	execSQL(t, c, `CREATE FUNCTION my_handler(internal) RETURNS index_am_handler AS 'my_handler' LANGUAGE C;`)
	execSQL(t, c, `CREATE ACCESS METHOD myam TYPE INDEX HANDLER my_handler;`)

	am := c.LookupAccessMethod("myam")
	if am == nil {
		t.Fatal("expected myam to exist")
	}
	if am.Type != 'i' {
		t.Errorf("access method type: got %c, want 'i'", am.Type)
	}
}

func TestCreateAccessMethod_Duplicate(t *testing.T) {
	c := New()
	execSQL(t, c, `CREATE FUNCTION h(internal) RETURNS index_am_handler AS 'h' LANGUAGE C;`)
	execSQL(t, c, `CREATE ACCESS METHOD dup_am TYPE INDEX HANDLER h;`)

	err := execSQLErr(c, `CREATE ACCESS METHOD dup_am TYPE INDEX HANDLER h;`)
	if err == nil {
		t.Fatal("expected error for duplicate access method")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCreateOpFamily(t *testing.T) {
	c := New()

	execSQL(t, c, `CREATE OPERATOR FAMILY my_fam USING btree;`)

	var found bool
	for _, fam := range c.opFamilies {
		if fam.Name == "my_fam" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected my_fam operator family to exist")
	}
}

func TestDefineOpClass_Default(t *testing.T) {
	c := New()

	// Create type + comparison functions + operators for the opclass.
	execSQL(t, c, `CREATE TYPE mytype;`)
	execSQL(t, c, `CREATE FUNCTION mytypein(cstring) RETURNS mytype AS 'mytypein' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE FUNCTION mytypeout(mytype) RETURNS cstring AS 'mytypeout' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE TYPE mytype (INPUT = mytypein, OUTPUT = mytypeout, LIKE = text);`)
	execSQL(t, c, `CREATE FUNCTION mytype_cmp(mytype, mytype) RETURNS integer AS 'mytype_cmp' LANGUAGE C IMMUTABLE STRICT;`)
	execSQL(t, c, `CREATE FUNCTION mytype_eq(mytype, mytype) RETURNS boolean AS 'mytype_eq' LANGUAGE C IMMUTABLE STRICT;`)
	execSQL(t, c, `CREATE OPERATOR = (LEFTARG = mytype, RIGHTARG = mytype, PROCEDURE = mytype_eq);`)

	execSQL(t, c, `CREATE OPERATOR CLASS mytype_ops DEFAULT FOR TYPE mytype USING btree AS
		OPERATOR 1 =,
		FUNCTION 1 mytype_cmp(mytype, mytype);`)

	mytypeOID := lookupTypeOID(t, c, "mytype")
	btreeAM := c.LookupAccessMethod("btree")
	if btreeAM == nil {
		t.Fatal("expected btree AM to exist")
	}

	opc := c.opClassByKey[opClassKey{amOID: btreeAM.OID, typeOID: mytypeOID}]
	if opc == nil {
		t.Fatal("expected default opclass for (btree, mytype)")
	}
	if opc.Name != "mytype_ops" {
		t.Errorf("opclass name: got %q, want %q", opc.Name, "mytype_ops")
	}
	if !opc.IsDefault {
		t.Error("expected opclass to be default")
	}
}

func TestDefineOpClass_CustomAM(t *testing.T) {
	c := New()

	execSQL(t, c, `CREATE FUNCTION my_handler(internal) RETURNS index_am_handler AS 'my_handler' LANGUAGE C;`)
	execSQL(t, c, `CREATE ACCESS METHOD custom_am TYPE INDEX HANDLER my_handler;`)

	execSQL(t, c, `CREATE TYPE mytype;`)
	execSQL(t, c, `CREATE FUNCTION mytypein(cstring) RETURNS mytype AS 'mytypein' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE FUNCTION mytypeout(mytype) RETURNS cstring AS 'mytypeout' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE TYPE mytype (INPUT = mytypein, OUTPUT = mytypeout, LIKE = text);`)
	execSQL(t, c, `CREATE FUNCTION mytype_eq(mytype, mytype) RETURNS boolean AS 'mytype_eq' LANGUAGE C IMMUTABLE STRICT;`)
	execSQL(t, c, `CREATE OPERATOR = (LEFTARG = mytype, RIGHTARG = mytype, PROCEDURE = mytype_eq);`)

	execSQL(t, c, `CREATE OPERATOR CLASS mytype_ops DEFAULT FOR TYPE mytype USING custom_am AS
		OPERATOR 1 =;`)

	// Verify it was registered.
	var found bool
	for _, opc := range c.opClasses {
		if opc.Name == "mytype_ops" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected mytype_ops opclass under custom_am")
	}
}

func TestDefineIndex_CustomAM(t *testing.T) {
	c := New()

	execSQL(t, c, `CREATE FUNCTION my_handler(internal) RETURNS index_am_handler AS 'my_handler' LANGUAGE C;`)
	execSQL(t, c, `CREATE ACCESS METHOD custom_am TYPE INDEX HANDLER my_handler;`)
	execSQL(t, c, `CREATE TABLE t (id integer);`)
	execSQL(t, c, `CREATE INDEX t_idx ON t USING custom_am (id);`)

	// Verify the index uses the custom AM.
	rel := c.GetRelation("", "t")
	if rel == nil {
		t.Fatal("expected table t")
	}
	idxs := c.IndexesOf(rel.OID)
	if len(idxs) == 0 {
		t.Fatal("expected at least one index on t")
	}
	if idxs[0].AccessMethod != "custom_am" {
		t.Errorf("index AM: got %q, want %q", idxs[0].AccessMethod, "custom_am")
	}
}

func TestDefineIndex_UnknownAM(t *testing.T) {
	c := New()

	execSQL(t, c, `CREATE TABLE t (id integer);`)
	err := execSQLErr(c, `CREATE INDEX t_idx ON t USING nonexistent_am (id);`)
	if err == nil {
		t.Fatal("expected error for unknown access method")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("unexpected error: %v", err)
	}
}
