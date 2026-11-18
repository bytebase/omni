package catalog

import (
	"testing"
)

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
		c := New()
		if _, err := c.Exec(sql, nil); err != nil {
			t.Errorf("Exec(%q): %v", sql, err)
		}
	}
}
