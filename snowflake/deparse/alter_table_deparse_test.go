package deparse_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/bytebase/omni/snowflake/parser"
	"github.com/stretchr/testify/require"
)

// ALTER TABLE data metric functions, governance actions, newer properties, and
// ALTER ICEBERG TABLE round-trip through the deparser.

func TestDeparse_AlterTableDataMetricFunction(t *testing.T) {
	assertRoundTrip(t, "ALTER TABLE t SET DATA_METRIC_SCHEDULE = 'TRIGGER_ON_CHANGES'")
	assertRoundTrip(t, "ALTER TABLE t UNSET DATA_METRIC_SCHEDULE")
	assertRoundTrip(t, "ALTER TABLE t ADD DATA METRIC FUNCTION s.core.null_count ON (c1, TABLE r (k)) "+
		"WITHIN GROUP (g) GROUP LIMIT 10 FILTER (c1 > 0) EXPECTATION e1 (value = 0), e2 (value < 5) "+
		"EXECUTE AS ROLE r1 ANOMALY_DETECTION = TRUE SENSITIVITY = 'LOW', m2 ON ()")
	assertRoundTrip(t, "ALTER TABLE t DROP DATA METRIC FUNCTION m ON (a), m2 ON (b)")
	assertRoundTrip(t, "ALTER TABLE t MODIFY DATA METRIC FUNCTION m ON (a) SUSPEND, m ON (b) RESUME")
	assertRoundTrip(t, "ALTER TABLE t MODIFY DATA METRIC FUNCTION m ON (a) MODIFY EXPECTATION e (value > 1)")
	assertRoundTrip(t, "ALTER TABLE t MODIFY DATA METRIC FUNCTION m ON (a) DROP EXPECTATION e, e2")
	assertRoundTrip(t, "ALTER TABLE t MODIFY DATA METRIC FUNCTION m ON (a) FILTER (), m ON (b) FILTER (b > 1)")
	assertRoundTrip(t, "ALTER TABLE t MODIFY DATA METRIC FUNCTION m SET ANOMALY_DETECTION = TRUE SENSITIVITY = 'HIGH'")
}

func TestDeparse_AlterTableGovernance(t *testing.T) {
	assertRoundTrip(t, "ALTER TABLE t SET TAG t1 = 'x', t2 = 'y'")
	assertRoundTrip(t, "ALTER TABLE t ALTER c1 SET TAG t1 = 'x', c2 UNSET TAG t2")
	assertRoundTrip(t, "ALTER TABLE t SET AGGREGATION POLICY p ENTITY KEY (a) FORCE")
	assertRoundTrip(t, "ALTER TABLE t UNSET AGGREGATION POLICY")
	assertRoundTrip(t, "ALTER TABLE t SET JOIN POLICY p FORCE")
	assertRoundTrip(t, "ALTER TABLE t UNSET JOIN POLICY")
	assertRoundTrip(t, "ALTER TABLE t ALTER c SET PROJECTION POLICY p FORCE")
	assertRoundTrip(t, "ALTER TABLE t MODIFY COLUMN c UNSET PROJECTION POLICY")
	assertRoundTrip(t, "ALTER TABLE t ADD STORAGE LIFECYCLE POLICY p ON (a, b)")
	assertRoundTrip(t, "ALTER TABLE t DROP STORAGE LIFECYCLE POLICY")
	assertRoundTrip(t, "ALTER TABLE t ALTER CONSTRAINT c1 (a, b) NOT ENFORCED NOVALIDATE RELY")
	assertRoundTrip(t, "ALTER TABLE t MODIFY PRIMARY KEY (a) ENFORCED")
	assertRoundTrip(t, "ALTER TABLE t SET ERROR_LOGGING = TRUE CONTACT steward = alice, support = bob ROW_TIMESTAMP = FALSE")
	assertRoundTrip(t, "ALTER TABLE t UNSET CONTACT steward, ENABLE_SCHEMA_EVOLUTION")
	assertRoundTrip(t, "ALTER TABLE t ADD SEARCH OPTIMIZATION ON FULL_TEXT(a, ANALYZER => 'UNICODE_ANALYZER')")
	assertRoundTrip(t, "ALTER TABLE t DROP SEARCH OPTIMIZATION ON c1, 1")
}

func TestDeparse_AlterIcebergTable(t *testing.T) {
	assertRoundTrip(t, "ALTER ICEBERG TABLE IF EXISTS t REFRESH 'metadata/v2.json'")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t CONVERT TO MANAGED BASE_LOCATION = 'p' DROP_FROM_SOURCE_CATALOG = TRUE")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t OPTIMIZE (REMOVE_POSITION_DELETES)")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t ADD PARTITION BY (BUCKET(10, c))")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t REPLACE PARTITION BY (c) WITH (DAY(c))")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t CREATE OR REPLACE VERSION_TAG v AS OF VERSION 12 RETAIN 3 DAYS")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t CREATE VERSION_TAG IF NOT EXISTS v AS OF TIMESTAMP 1700000000")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t ALTER VERSION_TAG v AS OF VERSION 3")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t DROP VERSION_TAG IF EXISTS v")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t ALTER COLUMN c SET WRITE DEFAULT 1, ALTER COLUMN d DROP WRITE DEFAULT")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t ALTER COLUMN c SET DATA TYPE INT RENAME FIELDS")
	assertRoundTrip(t, "ALTER ICEBERG TABLE t SET TARGET_FILE_SIZE = '64MB' CONTACT (steward = alice)")

	file, err := parser.Parse("ALTER ICEBERG TABLE t UPGRADE TO ICEBERG VERSION 3")
	require.NoError(t, err)
	out, err := deparse.DeparseFile(file)
	require.NoError(t, err)
	require.Equal(t, "ALTER ICEBERG TABLE t UPGRADE TO ICEBERG VERSION 3", out)
}
