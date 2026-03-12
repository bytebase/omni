package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// =============================================================================
// CreateCast tests
// =============================================================================

func TestCreateCast_BinaryCoercible(t *testing.T) {
	c := New()

	// Create a user-defined type so we can cast from it.
	execSQL(t, c, `CREATE TYPE citext;`)
	execSQL(t, c, `CREATE FUNCTION citextin(cstring) RETURNS citext AS 'citextin' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE FUNCTION citextout(citext) RETURNS cstring AS 'citextout' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE TYPE citext (INPUT = citextin, OUTPUT = citextout, LIKE = text);`)

	// CREATE CAST (citext AS text) WITHOUT FUNCTION AS IMPLICIT
	execSQL(t, c, `CREATE CAST (citext AS text) WITHOUT FUNCTION AS IMPLICIT;`)

	// Verify the cast was registered.
	citextOID := lookupTypeOID(t, c, "citext")
	cast := c.LookupCast(citextOID, TEXTOID)
	if cast == nil {
		t.Fatal("expected cast citext → text to exist")
	}
	if cast.Context != 'i' {
		t.Errorf("cast context: got %c, want 'i' (implicit)", cast.Context)
	}
	if cast.Method != 'b' {
		t.Errorf("cast method: got %c, want 'b' (binary)", cast.Method)
	}
}

func TestCreateCast_WithFunction(t *testing.T) {
	c := New()

	execSQL(t, c, `CREATE TYPE citext;`)
	execSQL(t, c, `CREATE FUNCTION citextin(cstring) RETURNS citext AS 'citextin' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE FUNCTION citextout(citext) RETURNS cstring AS 'citextout' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE TYPE citext (INPUT = citextin, OUTPUT = citextout, LIKE = text);`)

	// Create cast function.
	execSQL(t, c, `CREATE FUNCTION citext(bpchar) RETURNS citext AS 'citext_cast' LANGUAGE C IMMUTABLE STRICT;`)

	citextOID := lookupTypeOID(t, c, "citext")

	// CREATE CAST (bpchar AS citext) WITH FUNCTION citext(bpchar) AS ASSIGNMENT
	execSQL(t, c, `CREATE CAST (character AS citext) WITH FUNCTION citext(character) AS ASSIGNMENT;`)

	cast := c.LookupCast(BPCHAROID, citextOID)
	if cast == nil {
		t.Fatal("expected cast bpchar → citext to exist")
	}
	if cast.Context != 'a' {
		t.Errorf("cast context: got %c, want 'a' (assignment)", cast.Context)
	}
	if cast.Method != 'f' {
		t.Errorf("cast method: got %c, want 'f' (function)", cast.Method)
	}
	if cast.Func == 0 {
		t.Error("cast function OID should not be 0")
	}
}

func TestCreateCast_Duplicate(t *testing.T) {
	c := New()

	execSQL(t, c, `CREATE TYPE mytype;`)
	execSQL(t, c, `CREATE FUNCTION mytypein(cstring) RETURNS mytype AS 'mytypein' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE FUNCTION mytypeout(mytype) RETURNS cstring AS 'mytypeout' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE TYPE mytype (INPUT = mytypein, OUTPUT = mytypeout, LIKE = text);`)

	execSQL(t, c, `CREATE CAST (mytype AS text) WITHOUT FUNCTION AS IMPLICIT;`)

	// Duplicate cast should error.
	err := execSQLErr(c, `CREATE CAST (mytype AS text) WITHOUT FUNCTION AS IMPLICIT;`)
	if err == nil {
		t.Fatal("expected error for duplicate cast")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestCreateCast_UsedByCoercion(t *testing.T) {
	c := New()

	execSQL(t, c, `CREATE TYPE mytype;`)
	execSQL(t, c, `CREATE FUNCTION mytypein(cstring) RETURNS mytype AS 'mytypein' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE FUNCTION mytypeout(mytype) RETURNS cstring AS 'mytypeout' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE TYPE mytype (INPUT = mytypein, OUTPUT = mytypeout, LIKE = text);`)

	// Register an implicit cast mytype → text.
	execSQL(t, c, `CREATE CAST (mytype AS text) WITHOUT FUNCTION AS IMPLICIT;`)

	// Verify FindCoercionPathway can find it.
	mytypeOID := lookupTypeOID(t, c, "mytype")
	pathway, _ := c.FindCoercionPathway(mytypeOID, TEXTOID, 'i')
	if pathway == CoercionNone {
		t.Error("expected FindCoercionPathway to find the user-defined cast")
	}
}

// =============================================================================
// DefineOperator tests
// =============================================================================

func TestDefineOperator_Binary(t *testing.T) {
	c := New()

	execSQL(t, c, `CREATE TYPE mytype;`)
	execSQL(t, c, `CREATE FUNCTION mytypein(cstring) RETURNS mytype AS 'mytypein' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE FUNCTION mytypeout(mytype) RETURNS cstring AS 'mytypeout' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE TYPE mytype (INPUT = mytypein, OUTPUT = mytypeout, LIKE = text);`)
	execSQL(t, c, `CREATE FUNCTION mytype_eq(mytype, mytype) RETURNS boolean AS 'mytype_eq' LANGUAGE C IMMUTABLE STRICT;`)

	// CREATE OPERATOR = (LEFTARG = mytype, RIGHTARG = mytype, PROCEDURE = mytype_eq)
	execSQL(t, c, `CREATE OPERATOR = (LEFTARG = mytype, RIGHTARG = mytype, PROCEDURE = mytype_eq, COMMUTATOR = =, HASHES, MERGES);`)

	mytypeOID := lookupTypeOID(t, c, "mytype")
	ops := c.LookupOperatorExact("=", mytypeOID, mytypeOID)
	if len(ops) == 0 {
		t.Fatal("expected = operator for mytype to exist")
	}
	op := ops[0]
	if op.Result != BOOLOID {
		t.Errorf("operator result: got %d, want %d (BOOLOID)", op.Result, BOOLOID)
	}
	if !op.CanHash {
		t.Error("expected operator to support hashing")
	}
	if !op.CanMerge {
		t.Error("expected operator to support merge joins")
	}
	// Self-commutator for = with same types.
	if op.Com != op.OID {
		t.Errorf("expected self-commutator, got com=%d, oid=%d", op.Com, op.OID)
	}
}

func TestDefineOperator_MissingFunction(t *testing.T) {
	c := New()

	err := execSQLErr(c, `CREATE OPERATOR @@ (LEFTARG = integer, RIGHTARG = integer);`)
	if err == nil {
		t.Fatal("expected error for missing PROCEDURE")
	}
	if !strings.Contains(err.Error(), "operator function must be specified") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDefineOperator_MissingArgs(t *testing.T) {
	c := New()

	err := execSQLErr(c, `CREATE OPERATOR @@ (PROCEDURE = int4eq);`)
	if err == nil {
		t.Fatal("expected error for missing arg types")
	}
	if !strings.Contains(err.Error(), "operator argument types must be specified") {
		t.Errorf("unexpected error: %v", err)
	}
}

// =============================================================================
// DefineAggregate tests
// =============================================================================

func TestDefineAggregate_Simple(t *testing.T) {
	c := New()

	execSQL(t, c, `CREATE TYPE citext;`)
	execSQL(t, c, `CREATE FUNCTION citextin(cstring) RETURNS citext AS 'citextin' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE FUNCTION citextout(citext) RETURNS cstring AS 'citextout' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE TYPE citext (INPUT = citextin, OUTPUT = citextout, LIKE = text);`)

	citextOID := lookupTypeOID(t, c, "citext")

	// Create a transition function for the aggregate.
	execSQL(t, c, `CREATE FUNCTION citext_smaller(citext, citext) RETURNS citext AS 'citext_smaller' LANGUAGE C IMMUTABLE STRICT;`)

	// CREATE AGGREGATE min(citext) (SFUNC = citext_smaller, STYPE = citext)
	execSQL(t, c, `CREATE AGGREGATE min(citext) (SFUNC = citext_smaller, STYPE = citext);`)

	// Verify the aggregate was registered.
	procs := c.LookupProcByName("min")
	var found bool
	for _, p := range procs {
		if p.Kind == PROKIND_AGGREGATE && len(p.ArgTypes) == 1 && p.ArgTypes[0] == citextOID {
			found = true
			if p.RetType != citextOID {
				t.Errorf("aggregate return type: got %d, want %d", p.RetType, citextOID)
			}
			break
		}
	}
	if !found {
		t.Fatal("expected min(citext) aggregate to be registered")
	}
}

func TestDefineAggregate_MissingSfunc(t *testing.T) {
	c := New()

	err := execSQLErr(c, `CREATE AGGREGATE bad_agg(integer) (STYPE = integer);`)
	if err == nil {
		t.Fatal("expected error for missing sfunc")
	}
	if !strings.Contains(err.Error(), "aggregate sfunc must be specified") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDefineAggregate_MissingStype(t *testing.T) {
	c := New()

	err := execSQLErr(c, `CREATE AGGREGATE bad_agg(integer) (SFUNC = int4pl);`)
	if err == nil {
		t.Fatal("expected error for missing stype")
	}
	if !strings.Contains(err.Error(), "aggregate stype must be specified") {
		t.Errorf("unexpected error: %v", err)
	}
}

// =============================================================================
// Implicit shell type tests
// =============================================================================

func TestImplicitShellType_CLanguage(t *testing.T) {
	c := New()

	// C function referencing an unknown return type should auto-create a shell.
	execSQL(t, c, `CREATE FUNCTION foo_in(cstring) RETURNS foo_type AS 'foo_in' LANGUAGE C STRICT IMMUTABLE;`)

	// Verify the shell type exists.
	fooOID := lookupTypeOID(t, c, "foo_type")
	bt := c.TypeByOID(fooOID)
	if bt == nil {
		t.Fatal("expected shell type foo_type to exist")
	}
	if bt.IsDefined {
		t.Error("expected shell type to be not yet defined")
	}

	// Verify a warning was emitted.
	warnings := c.DrainWarnings()
	var foundWarning bool
	for _, w := range warnings {
		if strings.Contains(w.Message, "is not yet defined") {
			foundWarning = true
			break
		}
	}
	if !foundWarning {
		t.Error("expected warning about type not yet defined")
	}
}

func TestImplicitShellType_SQLRejects(t *testing.T) {
	c := New()

	// SQL function should NOT auto-create shell types — should get type error.
	err := execSQLErr(c, `CREATE FUNCTION bar(text) RETURNS nonexistent_type AS 'body' LANGUAGE SQL;`)
	if err == nil {
		t.Fatal("expected error for SQL function with nonexistent return type")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestImplicitShellType_ThenDefineType(t *testing.T) {
	c := New()

	// Pattern used by ltree: function creates shell, then DefineType upgrades it.
	execSQL(t, c, `CREATE FUNCTION ltree_in(cstring) RETURNS ltree AS 'ltree_in' LANGUAGE C STRICT IMMUTABLE;`)
	execSQL(t, c, `CREATE FUNCTION ltree_out(ltree) RETURNS cstring AS 'ltree_out' LANGUAGE C STRICT IMMUTABLE;`)
	c.DrainWarnings() // clear shell warnings

	// Now define the full type — should upgrade the shell.
	execSQL(t, c, `CREATE TYPE ltree (INPUT = ltree_in, OUTPUT = ltree_out, INTERNALLENGTH = VARIABLE, STORAGE = extended);`)

	bt := c.TypeByOID(lookupTypeOID(t, c, "ltree"))
	if bt == nil {
		t.Fatal("expected ltree type to exist")
	}
	if !bt.IsDefined {
		t.Error("expected ltree type to be fully defined after CREATE TYPE")
	}
	if bt.Type != 'b' {
		t.Errorf("expected base type, got %c", bt.Type)
	}
}

// =============================================================================
// Integration test: citext extension subset replay
// =============================================================================

func TestExtensionReplay_Citext(t *testing.T) {
	c := New()

	// Replay a simplified subset of citext--1.6.sql.
	sql := `
-- Shell type
CREATE TYPE citext;

-- I/O functions
CREATE FUNCTION citextin(cstring) RETURNS citext AS 'citextin' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION citextout(citext) RETURNS cstring AS 'citextout' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION citextrecv(internal) RETURNS citext AS 'citextrecv' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;
CREATE FUNCTION citextsend(citext) RETURNS bytea AS 'citextsend' LANGUAGE C STRICT IMMUTABLE PARALLEL SAFE;

-- Full type definition
CREATE TYPE citext (
    INPUT          = citextin,
    OUTPUT         = citextout,
    RECEIVE        = citextrecv,
    SEND           = citextsend,
    INTERNALLENGTH = VARIABLE,
    STORAGE        = extended,
    CATEGORY       = 'S',
    PREFERRED      = false
);

-- Casts
CREATE CAST (citext AS text)    WITHOUT FUNCTION AS IMPLICIT;
CREATE CAST (citext AS character varying) WITHOUT FUNCTION AS IMPLICIT;
CREATE CAST (citext AS character)         WITHOUT FUNCTION AS ASSIGNMENT;
CREATE CAST (text AS citext)    WITHOUT FUNCTION AS ASSIGNMENT;
CREATE CAST (character varying AS citext) WITHOUT FUNCTION AS ASSIGNMENT;
CREATE CAST (character AS citext)         WITHOUT FUNCTION AS ASSIGNMENT;

-- Comparison functions
CREATE FUNCTION citext_eq(citext, citext) RETURNS boolean AS 'citext_eq' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_ne(citext, citext) RETURNS boolean AS 'citext_ne' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_lt(citext, citext) RETURNS boolean AS 'citext_lt' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_le(citext, citext) RETURNS boolean AS 'citext_le' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_gt(citext, citext) RETURNS boolean AS 'citext_gt' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE FUNCTION citext_ge(citext, citext) RETURNS boolean AS 'citext_ge' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;

-- Operators
CREATE OPERATOR = (
    LEFTARG    = citext,
    RIGHTARG   = citext,
    COMMUTATOR = =,
    NEGATOR    = <>,
    PROCEDURE  = citext_eq,
    RESTRICT   = eqsel,
    JOIN       = eqjoinsel,
    HASHES,
    MERGES
);
CREATE OPERATOR <> (
    LEFTARG    = citext,
    RIGHTARG   = citext,
    NEGATOR    = =,
    COMMUTATOR = <>,
    PROCEDURE  = citext_ne,
    RESTRICT   = neqsel,
    JOIN       = neqjoinsel
);
CREATE OPERATOR < (
    LEFTARG    = citext,
    RIGHTARG   = citext,
    NEGATOR    = >=,
    PROCEDURE  = citext_lt,
    RESTRICT   = scalarltsel,
    JOIN       = scalarltjoinsel
);

-- Aggregate: min
CREATE FUNCTION citext_smaller(citext, citext) RETURNS citext AS 'citext_smaller' LANGUAGE C IMMUTABLE STRICT PARALLEL SAFE;
CREATE AGGREGATE min(citext) (
    SFUNC = citext_smaller,
    STYPE = citext,
    SORTOP = <,
    PARALLEL = safe
);
`
	stmts := mustParse(t, sql)
	for i, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("statement %d failed: %v", i+1, err)
		}
	}

	// Verify type.
	citextOID := lookupTypeOID(t, c, "citext")
	bt := c.TypeByOID(citextOID)
	if bt == nil || !bt.IsDefined {
		t.Fatal("expected citext type to be fully defined")
	}
	if bt.Category != 'S' {
		t.Errorf("citext category: got %c, want 'S'", bt.Category)
	}

	// Verify casts.
	if cast := c.LookupCast(citextOID, TEXTOID); cast == nil {
		t.Error("expected cast citext → text")
	}
	if cast := c.LookupCast(TEXTOID, citextOID); cast == nil {
		t.Error("expected cast text → citext")
	}

	// Verify operators.
	if ops := c.LookupOperatorExact("=", citextOID, citextOID); len(ops) == 0 {
		t.Error("expected = operator for citext")
	}
	if ops := c.LookupOperatorExact("<>", citextOID, citextOID); len(ops) == 0 {
		t.Error("expected <> operator for citext")
	}
	if ops := c.LookupOperatorExact("<", citextOID, citextOID); len(ops) == 0 {
		t.Error("expected < operator for citext")
	}

	// Verify aggregate.
	procs := c.LookupProcByName("min")
	var foundAgg bool
	for _, p := range procs {
		if p.Kind == PROKIND_AGGREGATE && len(p.ArgTypes) == 1 && p.ArgTypes[0] == citextOID {
			foundAgg = true
			break
		}
	}
	if !foundAgg {
		t.Error("expected min(citext) aggregate to be registered")
	}
}

// =============================================================================
// CreateExtension tests
// =============================================================================

func TestCreateExtension_Citext(t *testing.T) {
	c := New()
	execSQL(t, c, `CREATE EXTENSION citext;`)

	// Verify type exists and is fully defined.
	citextOID := lookupTypeOID(t, c, "citext")
	bt := c.TypeByOID(citextOID)
	if bt == nil || !bt.IsDefined {
		t.Fatal("expected citext type to be fully defined")
	}
	if bt.Category != 'S' {
		t.Errorf("citext category: got %c, want 'S'", bt.Category)
	}

	// Verify casts.
	if cast := c.LookupCast(citextOID, TEXTOID); cast == nil {
		t.Error("expected cast citext -> text")
	}
	if cast := c.LookupCast(TEXTOID, citextOID); cast == nil {
		t.Error("expected cast text -> citext")
	}
	if cast := c.LookupCast(citextOID, VARCHAROID); cast == nil {
		t.Error("expected cast citext -> varchar")
	}

	// Verify operators.
	if ops := c.LookupOperatorExact("=", citextOID, citextOID); len(ops) == 0 {
		t.Error("expected = operator for citext")
	}
	if ops := c.LookupOperatorExact("<>", citextOID, citextOID); len(ops) == 0 {
		t.Error("expected <> operator for citext")
	}
	if ops := c.LookupOperatorExact("<", citextOID, citextOID); len(ops) == 0 {
		t.Error("expected < operator for citext")
	}

	// Verify aggregates.
	procs := c.LookupProcByName("min")
	var foundMin bool
	for _, p := range procs {
		if p.Kind == PROKIND_AGGREGATE && len(p.ArgTypes) == 1 && p.ArgTypes[0] == citextOID {
			foundMin = true
			break
		}
	}
	if !foundMin {
		t.Error("expected min(citext) aggregate")
	}
}

func TestCreateExtension_Hstore(t *testing.T) {
	c := New()
	execSQL(t, c, `CREATE EXTENSION hstore;`)

	// Verify type exists.
	hstoreOID := lookupTypeOID(t, c, "hstore")
	bt := c.TypeByOID(hstoreOID)
	if bt == nil || !bt.IsDefined {
		t.Fatal("expected hstore type to be fully defined")
	}

	// Verify key operators.
	if ops := c.LookupOperatorExact("->", hstoreOID, TEXTOID); len(ops) == 0 {
		t.Error("expected -> operator for hstore")
	}
	if ops := c.LookupOperatorExact("=", hstoreOID, hstoreOID); len(ops) == 0 {
		t.Error("expected = operator for hstore")
	}

	// Verify JSON cast.
	if cast := c.LookupCast(hstoreOID, JSONOID); cast == nil {
		t.Error("expected cast hstore -> json")
	}
}

func TestCreateExtension_Schema(t *testing.T) {
	c := New()
	execSQL(t, c, `CREATE SCHEMA utils;`)
	execSQL(t, c, `CREATE EXTENSION hstore SCHEMA utils;`)

	// Verify hstore type is in utils schema.
	s := c.GetSchema("utils")
	if s == nil {
		t.Fatal("expected utils schema")
	}
	// Look up type via schema-qualified search.
	key := typeKey{ns: s.OID, name: "hstore"}
	bt := c.typeByName[key]
	if bt == nil {
		t.Fatal("expected hstore type in utils schema")
	}
	if !bt.IsDefined {
		t.Error("expected hstore type to be fully defined")
	}
}

func TestCreateExtension_IfNotExists(t *testing.T) {
	c := New()
	execSQL(t, c, `CREATE EXTENSION citext;`)

	// Second CREATE EXTENSION IF NOT EXISTS should not error.
	err := execSQLErr(c, `CREATE EXTENSION IF NOT EXISTS citext;`)
	if err != nil {
		t.Fatalf("IF NOT EXISTS should not error: %v", err)
	}
	warnings := c.DrainWarnings()
	var foundSkip bool
	for _, w := range warnings {
		if strings.Contains(w.Message, "already exists") {
			foundSkip = true
			break
		}
	}
	if !foundSkip {
		t.Error("expected 'already exists' warning")
	}
}

func TestCreateExtension_PGVector(t *testing.T) {
	c := New()
	execSQL(t, c, `CREATE EXTENSION vector;`)

	// Verify vector type with typmod support.
	vectorOID := lookupTypeOID(t, c, "vector")
	bt := c.TypeByOID(vectorOID)
	if bt == nil || !bt.IsDefined {
		t.Fatal("expected vector type to be fully defined")
	}
	if bt.ModIn == 0 {
		t.Error("expected vector type to have TYPMOD_IN set")
	}

	// Verify halfvec type with typmod support.
	halfvecOID := lookupTypeOID(t, c, "halfvec")
	hvbt := c.TypeByOID(halfvecOID)
	if hvbt == nil || !hvbt.IsDefined {
		t.Fatal("expected halfvec type to be fully defined")
	}
	if hvbt.ModIn == 0 {
		t.Error("expected halfvec type to have TYPMOD_IN set")
	}

	// Verify halfvec(384) resolves correctly.
	oid, typmod, err := c.ResolveType(TypeName{Name: "halfvec", TypeMod: 384})
	if err != nil {
		t.Fatalf("halfvec(384) resolve failed: %v", err)
	}
	if oid != halfvecOID {
		t.Errorf("halfvec(384) OID: got %d, want %d", oid, halfvecOID)
	}
	if typmod != 384 {
		t.Errorf("halfvec(384) typmod: got %d, want 384", typmod)
	}

	// Verify hnsw access method exists.
	am := c.LookupAccessMethod("hnsw")
	if am == nil {
		t.Fatal("expected hnsw access method to exist")
	}

	// Verify distance operators.
	if ops := c.LookupOperatorExact("<->", vectorOID, vectorOID); len(ops) == 0 {
		t.Error("expected <-> operator for vector")
	}
	if ops := c.LookupOperatorExact("<=>", halfvecOID, halfvecOID); len(ops) == 0 {
		t.Error("expected <=> operator for halfvec")
	}
}

func TestCreateExtension_PGVector_Index(t *testing.T) {
	c := New()
	execSQL(t, c, `CREATE EXTENSION vector;`)

	// Reproduce issue_295_pgvector_typmod: CREATE TABLE + CREATE INDEX USING hnsw.
	execSQL(t, c, `CREATE TABLE activity (
		id bigserial PRIMARY KEY,
		embedding halfvec(384)
	);`)

	// This is the key test: CREATE INDEX USING hnsw with a named opclass.
	// The hnsw AM must be registered for this to succeed.
	execSQL(t, c, `CREATE INDEX activity_embedding_idx
		ON activity USING hnsw (embedding halfvec_cosine_ops);`)

	rel := c.GetRelation("", "activity")
	if rel == nil {
		t.Fatal("expected table activity")
	}

	// Verify the embedding column type has correct typmod.
	for _, col := range rel.Columns {
		if col.Name == "embedding" {
			if col.TypeMod != 384 {
				t.Errorf("embedding typmod: got %d, want 384", col.TypeMod)
			}
			break
		}
	}

	// Verify the index uses hnsw.
	idxs := c.IndexesOf(rel.OID)
	var found bool
	for _, idx := range idxs {
		if idx.Name == "activity_embedding_idx" {
			if idx.AccessMethod != "hnsw" {
				t.Errorf("index AM: got %q, want %q", idx.AccessMethod, "hnsw")
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("expected activity_embedding_idx index")
	}
}

func TestRegisterExtensionSQL(t *testing.T) {
	// Register a custom extension.
	RegisterExtensionSQL("my_ext", `
		CREATE TYPE my_ext_type;
		CREATE FUNCTION my_ext_in(cstring) RETURNS my_ext_type AS 'my_ext_in' LANGUAGE C STRICT IMMUTABLE;
		CREATE FUNCTION my_ext_out(my_ext_type) RETURNS cstring AS 'my_ext_out' LANGUAGE C STRICT IMMUTABLE;
		CREATE TYPE my_ext_type (INPUT = my_ext_in, OUTPUT = my_ext_out, INTERNALLENGTH = VARIABLE);
	`)

	c := New()
	execSQL(t, c, `CREATE EXTENSION my_ext;`)

	// Verify the type was created.
	oid := lookupTypeOID(t, c, "my_ext_type")
	bt := c.TypeByOID(oid)
	if bt == nil || !bt.IsDefined {
		t.Fatal("expected my_ext_type to be fully defined")
	}

	// Clean up global state.
	delete(extensionScripts, "my_ext")
}

func TestCreateExtension_Unknown(t *testing.T) {
	c := New()

	// Unknown extension should emit warning, not error.
	err := execSQLErr(c, `CREATE EXTENSION nonexistent_ext;`)
	if err != nil {
		t.Fatalf("unknown extension should not error: %v", err)
	}
	warnings := c.DrainWarnings()
	var foundWarning bool
	for _, w := range warnings {
		if strings.Contains(w.Message, "not bundled") {
			foundWarning = true
			break
		}
	}
	if !foundWarning {
		t.Error("expected 'not bundled' warning for unknown extension")
	}
}

// =============================================================================
// Helpers
// =============================================================================

// execSQL parses and executes SQL, fataling on error.
func execSQL(t *testing.T, c *Catalog, sql string) {
	t.Helper()
	stmts := mustParse(t, sql)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("execSQL(%q): %v", sql, err)
		}
	}
}

// execSQLErr parses and executes SQL, returning the first error.
func execSQLErr(c *Catalog, sql string) error {
	list, err := pgparser.Parse(sql)
	if err != nil {
		return err
	}
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		var node nodes.Node
		if raw, ok := item.(*nodes.RawStmt); ok {
			node = raw.Stmt.(nodes.Node)
		} else {
			node = item.(nodes.Node)
		}
		if err := c.ProcessUtility(node); err != nil {
			return err
		}
	}
	return nil
}

// mustParse parses SQL and returns statement nodes.
func mustParse(t *testing.T, sql string) []nodes.Node {
	t.Helper()
	list, err := pgparser.Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v\nSQL: %s", err, sql)
	}
	if list == nil || len(list.Items) == 0 {
		t.Fatalf("parse returned no statements\nSQL: %s", sql)
	}
	out := make([]nodes.Node, len(list.Items))
	for i, item := range list.Items {
		if raw, ok := item.(*nodes.RawStmt); ok {
			out[i] = raw.Stmt
		} else {
			out[i] = item
		}
	}
	return out
}

// lookupTypeOID resolves a type name to its OID via the catalog.
func lookupTypeOID(t *testing.T, c *Catalog, name string) uint32 {
	t.Helper()
	oid, _, err := c.ResolveType(TypeName{Name: name, TypeMod: -1})
	if err != nil {
		t.Fatalf("type %q not found: %v", name, err)
	}
	return oid
}
