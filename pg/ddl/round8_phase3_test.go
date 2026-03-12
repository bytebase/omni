package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Round 8 Phase 3: Sequence & Schema Improvements
// =============================================================================

// -----------------------------------------------------------------------------
// sequence.go: Strict sequence type validation (resolveSeqType rejects unknowns)
// -----------------------------------------------------------------------------

func TestSequenceInvalidTypeError(t *testing.T) {
	c := New()
	stmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "as",
				Arg:     &nodes.String{Str: "text"},
			},
		}},
	}
	err := c.DefineSequence(stmt)
	assertCode(t, err, CodeUndefinedObject)
	if !strings.Contains(err.Error(), "text") {
		t.Errorf("expected error to mention type name 'text', got: %s", err)
	}
}

func TestSequenceInvalidTypeFloat(t *testing.T) {
	c := New()
	stmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "as",
				Arg:     &nodes.String{Str: "float8"},
			},
		}},
	}
	err := c.DefineSequence(stmt)
	assertCode(t, err, CodeUndefinedObject)
	if !strings.Contains(err.Error(), "float8") {
		t.Errorf("expected error to mention type name 'float8', got: %s", err)
	}
}

func TestSequenceValidTypeSmallint(t *testing.T) {
	c := New()
	stmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "as",
				Arg:     &nodes.String{Str: "smallint"},
			},
		}},
	}
	err := c.DefineSequence(stmt)
	if err != nil {
		t.Fatal(err)
	}
	seq, _ := c.findSequence("", "seq1")
	if seq.TypeOID != INT2OID {
		t.Errorf("expected INT2OID, got %d", seq.TypeOID)
	}
}

func TestSequenceValidTypeInteger(t *testing.T) {
	c := New()
	stmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "as",
				Arg:     &nodes.String{Str: "integer"},
			},
		}},
	}
	err := c.DefineSequence(stmt)
	if err != nil {
		t.Fatal(err)
	}
	seq, _ := c.findSequence("", "seq1")
	if seq.TypeOID != INT4OID {
		t.Errorf("expected INT4OID, got %d", seq.TypeOID)
	}
}

func TestSequenceValidTypeBigint(t *testing.T) {
	c := New()
	stmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "as",
				Arg:     &nodes.String{Str: "bigint"},
			},
		}},
	}
	err := c.DefineSequence(stmt)
	if err != nil {
		t.Fatal(err)
	}
	seq, _ := c.findSequence("", "seq1")
	if seq.TypeOID != INT8OID {
		t.Errorf("expected INT8OID, got %d", seq.TypeOID)
	}
}

// -----------------------------------------------------------------------------
// sequence.go: SEQUENCE_NAME rejection in CREATE SEQUENCE
// -----------------------------------------------------------------------------

func TestSequenceNameOptionError(t *testing.T) {
	c := New()
	stmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "sequence_name",
				Arg:     &nodes.String{Str: "myseq"},
			},
		}},
	}
	err := c.DefineSequence(stmt)
	assertCode(t, err, CodeSyntaxError)
	if !strings.Contains(err.Error(), "conflicting or redundant options") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// schemacmds.go: CREATE SCHEMA AUTHORIZATION
// -----------------------------------------------------------------------------

func TestCreateSchemaAuthorizationNoName(t *testing.T) {
	c := New()
	// CREATE SCHEMA AUTHORIZATION alice -> schema name defaults to "alice"
	stmt := &nodes.CreateSchemaStmt{
		Authrole: &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CSTRING),
			Rolename: "alice",
		},
	}
	err := c.CreateSchemaCommand(stmt)
	if err != nil {
		t.Fatal(err)
	}
	s := c.GetSchema("alice")
	if s == nil {
		t.Fatal("expected schema 'alice' to exist")
	}
	if s.Owner != "alice" {
		t.Errorf("expected owner 'alice', got %q", s.Owner)
	}
}

func TestCreateSchemaAuthorizationWithName(t *testing.T) {
	c := New()
	// CREATE SCHEMA myschema AUTHORIZATION alice -> schema name is "myschema", owner is "alice"
	stmt := &nodes.CreateSchemaStmt{
		Schemaname: "myschema",
		Authrole: &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CSTRING),
			Rolename: "alice",
		},
	}
	err := c.CreateSchemaCommand(stmt)
	if err != nil {
		t.Fatal(err)
	}
	s := c.GetSchema("myschema")
	if s == nil {
		t.Fatal("expected schema 'myschema' to exist")
	}
	if s.Owner != "alice" {
		t.Errorf("expected owner 'alice', got %q", s.Owner)
	}
}

func TestCreateSchemaAuthorizationPgPrefixError(t *testing.T) {
	c := New()
	// CREATE SCHEMA AUTHORIZATION pg_admin -> schema name would be "pg_admin" which is reserved
	stmt := &nodes.CreateSchemaStmt{
		Authrole: &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CSTRING),
			Rolename: "pg_admin",
		},
	}
	err := c.CreateSchemaCommand(stmt)
	assertCode(t, err, CodeReservedName)
	if !strings.Contains(err.Error(), "pg_") {
		t.Errorf("expected error to mention pg_ prefix, got: %s", err)
	}
}

func TestCreateSchemaNoAuthNoOwner(t *testing.T) {
	c := New()
	stmt := &nodes.CreateSchemaStmt{
		Schemaname: "myschema",
	}
	err := c.CreateSchemaCommand(stmt)
	if err != nil {
		t.Fatal(err)
	}
	s := c.GetSchema("myschema")
	if s == nil {
		t.Fatal("expected schema 'myschema' to exist")
	}
	if s.Owner != "" {
		t.Errorf("expected empty owner, got %q", s.Owner)
	}
}

// -----------------------------------------------------------------------------
// alter.go: RENAME SCHEMA to pg_ prefix rejection
// -----------------------------------------------------------------------------

func TestRenameSchemaReservedPrefixError(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("myschema", false))

	err := c.ExecRenameStmt(makeRenameSchemaStmt("myschema", "pg_badname"))
	assertCode(t, err, CodeReservedName)
	if !strings.Contains(err.Error(), "unacceptable schema name") {
		t.Errorf("expected 'unacceptable schema name', got: %s", err)
	}
	if !strings.Contains(err.Error(), "pg_") {
		t.Errorf("expected error to mention pg_ prefix, got: %s", err)
	}
}

func TestRenameSchemaReservedPrefixPgTemp(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("myschema", false))

	err := c.ExecRenameStmt(makeRenameSchemaStmt("myschema", "pg_temp_999"))
	assertCode(t, err, CodeReservedName)
}

func TestRenameSchemaValidNameOK(t *testing.T) {
	c := New()
	c.CreateSchemaCommand(makeCreateSchemaStmt("old_schema", false))

	err := c.ExecRenameStmt(makeRenameSchemaStmt("old_schema", "new_schema"))
	if err != nil {
		t.Fatal(err)
	}
	if c.GetSchema("old_schema") != nil {
		t.Error("old schema name still resolves")
	}
	if c.GetSchema("new_schema") == nil {
		t.Error("new schema name not found")
	}
}

// -----------------------------------------------------------------------------
// sequence.go: OWNED BY view support
// -----------------------------------------------------------------------------

func TestSequenceOwnedByViewOK(t *testing.T) {
	c := New()
	// Create a table first (needed for the view).
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// Create a view.
	viewStmt := &nodes.ViewStmt{
		View: &nodes.RangeVar{Relname: "v"},
		Query: &nodes.SelectStmt{
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Name: "id",
					Val:  &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "id"}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "t"},
			}},
		},
	}
	if err := c.DefineView(viewStmt); err != nil {
		t.Fatal(err)
	}
	// Create a sequence and set OWNED BY on the view's column.
	seqStmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "owned_by",
				Arg: &nodes.List{Items: []nodes.Node{
					&nodes.String{Str: "v"},
					&nodes.String{Str: "id"},
				}},
			},
		}},
	}
	err := c.DefineSequence(seqStmt)
	if err != nil {
		t.Fatal(err)
	}
	seq, _ := c.findSequence("", "seq1")
	if seq.OwnerRelOID == 0 {
		t.Error("expected sequence to have owner relation OID set")
	}
}

func TestSequenceOwnedByMatviewError(t *testing.T) {
	// Materialized views (relkind 'm') should still be rejected by setSequenceOwner.
	c := New()
	// Create a table.
	if err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r'); err != nil {
		t.Fatal(err)
	}
	// Create a materialized view.
	matviewStmt := &nodes.CreateTableAsStmt{
		Into: &nodes.IntoClause{
			Rel: &nodes.RangeVar{Relname: "mv"},
		},
		Query: &nodes.SelectStmt{
			TargetList: &nodes.List{Items: []nodes.Node{
				&nodes.ResTarget{
					Name: "id",
					Val:  &nodes.ColumnRef{Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "id"}}}},
				},
			}},
			FromClause: &nodes.List{Items: []nodes.Node{
				&nodes.RangeVar{Relname: "t"},
			}},
		},
		Objtype: nodes.OBJECT_MATVIEW,
	}
	if err := c.ExecCreateTableAs(matviewStmt); err != nil {
		t.Fatal(err)
	}
	// Create sequence with OWNED BY on the matview's column — should fail.
	seqStmt := &nodes.CreateSeqStmt{
		Sequence: &nodes.RangeVar{Relname: "seq1"},
		Options: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "owned_by",
				Arg: &nodes.List{Items: []nodes.Node{
					&nodes.String{Str: "mv"},
					&nodes.String{Str: "id"},
				}},
			},
		}},
	}
	err := c.DefineSequence(seqStmt)
	assertCode(t, err, CodeInvalidParameterValue)
	if !strings.Contains(err.Error(), "materialized view") {
		t.Errorf("expected error to mention 'materialized view', got: %s", err)
	}
}
