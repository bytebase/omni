package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// execOne runs one statement and returns the error Exec recorded for it.
func execOne(t *testing.T, c *Catalog, sql string) error {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	if len(results) != 1 {
		t.Fatalf("Exec(%q): %d results, want one", sql, len(results))
	}
	return results[0].Error
}

// TestGrantRedshiftOnlyTargetsRecordNothing holds the grants on the Redshift
// objects the catalog does not model — a datashare, a COPY job, a model, a
// template, a warehouse permission — out of PostgreSQL's ACL validation: Exec
// records nothing and names no object to resolve.
func TestGrantRedshiftOnlyTargetsRecordNothing(t *testing.T) {
	for _, sql := range []string{
		"GRANT ALTER ON DATASHARE salesshare TO analyst",
		"GRANT USAGE ON DATASHARE salesshare TO ACCOUNT '123456789012' VIA DATA CATALOG",
		"REVOKE USAGE ON DATASHARE salesshare FROM NAMESPACE 'a3f3ae8c'",
		"GRANT ALTER ON COPY JOB nightly TO loader",
		"GRANT EXECUTE ON MODEL churn TO ml_user",
		"GRANT CONNECT ON WORKGROUP TO AWSIDC:alice",
		"GRANT DEBUG ON DATABASE federated TO ADMINS OF ACCOUNT '1234'",
		"GRANT CREATE MODEL TO ml_admin WITH GRANT OPTION",
		"GRANT ALTER ON TEMPLATE nightly TO loader",
	} {
		if err := execOne(t, New(), sql); err != nil {
			t.Errorf("Exec(%q): %v", sql, err)
		}
	}
}

// TestGrantExternalObjectsResolveAsRelations holds an external table on the
// relation path and an external schema on the schema path, where the loaded
// metadata installs them, so a schema diff sees their grants; their privileges
// are the ones r_GRANT and r_REVOKE list, not PostgreSQL's.
func TestGrantExternalObjectsResolveAsRelations(t *testing.T) {
	c := New()
	for _, sql := range []string{
		"CREATE SCHEMA spectrum",
		"CREATE TABLE spectrum.sales (id INT)",
	} {
		if err := execOne(t, c, sql); err != nil {
			t.Fatalf("Exec(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"GRANT SELECT ON EXTERNAL TABLE spectrum.sales TO IAM_ROLE 'arn:aws:iam::1:role/r'",
		"GRANT ALTER, DROP ON EXTERNAL TABLE spectrum.sales TO PUBLIC",
		"GRANT CREATE, ALTER ON EXTERNAL SCHEMA spectrum TO IAM_ROLE 'arn:aws:iam::1:role/r'",
		"REVOKE ALTER ON EXTERNAL SCHEMA spectrum FROM IAM_ROLE 'arn:aws:iam::1:role/r'",
	} {
		if err := execOne(t, c, sql); err != nil {
			t.Errorf("Exec(%q): %v", sql, err)
		}
	}
	if len(c.grants) == 0 {
		t.Fatal("no grant recorded for the external table or schema")
	}
	// The catalog holds the object type to the privileges r_GRANT and r_REVOKE
	// list for it, not PostgreSQL's relation or schema set.
	for _, tc := range []struct {
		objtype nodes.ObjectType
		priv    string
		code    byte
		valid   bool
	}{
		{nodes.OBJECT_FOREIGN_TABLE, "alter", 'r', true},
		{nodes.OBJECT_FOREIGN_TABLE, "truncate", 'r', false},
		{nodes.OBJECT_EXTERNAL_SCHEMA, "alter", 'n', true},
		{nodes.OBJECT_EXTERNAL_SCHEMA, "usage", 'n', false},
	} {
		if got := validGrantPrivilege(tc.objtype, tc.priv, tc.code); got != tc.valid {
			t.Errorf("validGrantPrivilege(%v, %q, %q) = %v, want %v", tc.objtype, tc.priv, string(tc.code), got, tc.valid)
		}
	}
}

// TestGrantRedshiftOnlyTargetsRefusePublicGrantOption holds the one check the
// catalog keeps for every Redshift-only target: a grant option cannot go to
// PUBLIC.
func TestGrantRedshiftOnlyTargetsRefusePublicGrantOption(t *testing.T) {
	for _, sql := range []string{
		"GRANT EXECUTE ON MODEL m TO PUBLIC WITH GRANT OPTION",
		"GRANT SELECT ON EXTERNAL TABLE s.t TO PUBLIC WITH GRANT OPTION",
		"GRANT ALTER ON TEMPLATE t TO PUBLIC WITH GRANT OPTION",
	} {
		err := execOne(t, New(), sql)
		assertCode(t, err, CodeInvalidGrantOperation)
	}
}
