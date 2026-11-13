package testcontainer

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
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
