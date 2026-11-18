package catalog

import (
	"testing"
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

// TestGrantRedshiftOnlyTargetsRecordNothing holds the grants on Redshift's own
// object types out of PostgreSQL's ACL validation: the catalog models none of
// them, so Exec records nothing and names no object to resolve.
func TestGrantRedshiftOnlyTargetsRecordNothing(t *testing.T) {
	for _, sql := range []string{
		"GRANT SELECT ON EXTERNAL TABLE spectrum.sales TO IAM_ROLE 'arn:aws:iam::1:role/r'",
		"GRANT CREATE ON EXTERNAL SCHEMA spectrum TO IAM_ROLE 'arn:aws:iam::1:role/r'",
		"REVOKE ALTER ON EXTERNAL SCHEMA spectrum FROM IAM_ROLE 'arn:aws:iam::1:role/r'",
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

// TestGrantRedshiftOnlyTargetsRefusePublicGrantOption holds the one check the
// catalog keeps for a Redshift-only target: a grant option cannot go to PUBLIC.
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
