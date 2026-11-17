package testcontainer

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestCodecRoundTrips(t *testing.T) {
	for _, in := range []*pgconn.PgError{
		{Severity: "ERROR", Code: "42601", Message: `syntax error at or near "FROM"`, Position: 8},
		{Severity: "ERROR", Code: "42501", Message: `permission denied for schema pg_auto_copy`},
		{Severity: "ERROR", Code: "42703", Message: `column "nope" does not exist`, Position: 8,
			SchemaName: "public", TableName: "lookup", ColumnName: "nope"},
		{Severity: "ERROR", Code: "23503", Message: "insert or update on table \"facts\" violates foreign key constraint",
			Detail: "Key (lookup_id)=(3) is not present in table \"lookup\".", ConstraintName: "facts_lookup_id_fkey"},
		{Severity: "FATAL", Code: "57P01", Message: "terminating connection due to administrator command"},
		{Severity: "ERROR", Code: "22023", Message: `a message with a \ and: a colon`},
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
		t.Error("a dropped connection was encoded as the cluster's answer")
	}
	if err := (Codec{}).DecodeError("no code"); err == nil {
		t.Error("a line with no SQLSTATE decoded")
	}
}

func TestReleasePatternReadsTheVersion(t *testing.T) {
	const version = "PostgreSQL 8.0.2 on i686-pc-linux-gnu, compiled by GCC gcc (GCC) 3.4.2 20041017 (Red Hat 3.4.2-6.fc3), Redshift 1.0.539507"
	matches := releasePattern.FindStringSubmatch(version)
	if len(matches) != 2 || matches[1] != "1.0.539507" {
		t.Errorf("releasePattern = %q", matches)
	}
}
