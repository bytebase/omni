package testcontainer

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/enginetest/tape"
)

func TestCodecRoundTrips(t *testing.T) {
	for _, in := range []*pgconn.PgError{
		{Severity: "ERROR", Code: "42601", Message: `syntax error at or near "FROM"`, Position: 8},
		{Severity: "ERROR", Code: "23505", Message: `duplicate key value violates unique constraint "t_pkey"`,
			Detail: "Key (a)=(1) already exists.", SchemaName: "public", TableName: "t", ConstraintName: "t_pkey"},
		{Severity: "FATAL", Code: "57P01", Message: "terminating connection due to administrator command"},
		{Severity: "ERROR", Code: "42883", Message: "function f() does not exist", Hint: "No function matches.\nYou might need to add explicit type casts.",
			InternalPosition: 8, InternalQuery: "SELECT f()", Where: `PL/pgSQL function g() line 3 at RETURN\nSQL statement "SELECT g()"`},
		{Severity: "ERROR", Code: "22023", Message: `a message with a \ and: a colon`, ColumnName: "c", DataTypeName: "d"},
	} {
		in.SeverityUnlocalized = in.Severity
		text, ok := Codec{}.EncodeError(in)
		if !ok {
			t.Fatalf("%v is not encoded", in)
		}
		var out *pgconn.PgError
		if err := (Codec{}).DecodeError(text); !errors.As(err, &out) || *out != *in {
			t.Errorf("%q decodes to %#v, want %#v", text, err, in)
		}
	}
	catalogKey := &pgconn.PgError{Severity: "ERROR", Code: "23505", Message: `duplicate key value violates unique constraint "pg_enum_typid_label_index"`,
		Detail: "Key (enumtypid, enumlabel)=(52831, a) already exists.", SchemaName: "pg_catalog", TableName: "pg_enum"}
	if text, _ := (Codec{}).EncodeError(catalogKey); !strings.Contains(text, "detail: Key (enumtypid, enumlabel)=(…, a) already exists.") || catalogKey.Detail != "Key (enumtypid, enumlabel)=(52831, a) already exists." {
		t.Errorf("a catalog key's OID is written %q, and the error left %q", text, catalogKey.Detail)
	}
	if _, ok := (Codec{}).EncodeError(errors.New("unexpected EOF")); ok {
		t.Error("a dropped connection was encoded as the engine's answer")
	}
	if err := (Codec{}).DecodeError("no code"); err == nil {
		t.Error("a line with no SQLSTATE decoded")
	}
}

func TestDatabaseNamesATestsDatabase(t *testing.T) {
	if got := Database(t); got != "omni_testdatabasenamesatestsdatabase" {
		t.Errorf("Database = %q", got)
	}
	t.Run("Sub/Test", func(t *testing.T) {
		if got := Database(t); got != "omni_testdatabasenamesatestsdatabase" {
			t.Errorf("a subtest's Database = %q, want its top-level test's", got)
		}
	})
}

func TestVersionReadsASourceTag(t *testing.T) {
	if got := Version("REL_17_12"); got != "17.12" {
		t.Errorf("Version = %q", got)
	}
}

// TestOIDNumberingIsPerServer pins that each pin role numbers the OIDs it
// reads from its own start. Two servers' counters start wherever the run
// left them, so a stand-in shared between roles would make a tape depend on
// what the recording did before instead of on the data it recorded.
func TestOIDNumberingIsPerServer(t *testing.T) {
	one, other := oidNumbering(t.Name(), Default), oidNumbering(t.Name(), PG18)
	if one == other {
		t.Fatal("two pin roles share one numbering")
	}
	if got, want := one.number(16388), int64(1<<32+1); got != want {
		t.Errorf("the first OID of a role is stand-in %d, want %d", got, want)
	}
	if got, want := one.number(16388), int64(1<<32+1); got != want {
		t.Errorf("rereading an OID changed its stand-in to %d, want %d", got, want)
	}
	if got, want := one.number(16390), int64(1<<32+2); got != want {
		t.Errorf("the next OID of a role is stand-in %d, want %d", got, want)
	}
	if got, want := other.number(16388), int64(1<<32+1); got != want {
		t.Errorf("another role's first OID is stand-in %d, want %d: its server counts on its own", got, want)
	}
	if again := oidNumbering(t.Name(), Default); again != one {
		t.Error("a role's second connection got another numbering")
	}
}

// TestCodecRenamesOIDs pins that the codec is a tape.Renamer whose stand-ins
// it matches itself, so a stand-in a test sends back reaches the server as
// the OID it stands for, and that no OID's text is one.
func TestCodecRenamesOIDs(t *testing.T) {
	var r tape.Renamer = Codec{}
	in := r.StandIn("oid:16390", 3)
	if in != "oid:03" || !r.MadeUp().MatchString(in) {
		t.Errorf("stand-in %q does not match %v", in, r.MadeUp())
	}
	if got := r.MadeUp().FindAllString("'oid:' || c.oid, oid:16390, oid:0", -1); !slices.Equal(got, []string{"oid:16390", "oid:0"}) {
		t.Errorf("MadeUp found %q", got)
	}
}
