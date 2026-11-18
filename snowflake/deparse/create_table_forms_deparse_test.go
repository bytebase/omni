package deparse_test

import (
	"strings"
	"testing"
)

// The CREATE TABLE forms that restore, template, clone read-only, or make an
// interactive table round-trip, and CTAS keeps the clauses before its AS.
func TestDeparse_CreateTable_Forms(t *testing.T) {
	for _, sql := range []string{
		`CREATE TABLE t FROM BACKUP SET bs IDENTIFIER 'b1'`,
		`CREATE ICEBERG TABLE t FROM BACKUP SET bs IDENTIFIER 'b1'`,
		`CREATE ICEBERG TABLE t (a INT, b INT AS 1) BASE_LOCATION = 'p/'`,
		`CREATE TABLE t FROM ARCHIVE OF s AS x WHERE x.a > 1`,
		`CREATE TABLE t COPY GRANTS USING TEMPLATE (SELECT 1)`,
		`CREATE TEMPORARY READ ONLY TABLE t CLONE s COPY GRANTS`,
		`CREATE OR REPLACE INTERACTIVE TABLE t (a INT) CLUSTER BY (a) TARGET_LAG = '1 minute' WAREHOUSE = wh COMMENT = 'c' AS SELECT a FROM s`,
		`CREATE TABLE t CLUSTER BY (a) COPY GRANTS AS SELECT 1 AS a`,
		`CREATE INTERACTIVE TABLE t CLUSTER BY (a) WAREHOUSE = "my warehouse" AS SELECT 1 AS a`,
		`CREATE TABLE t (a INT) COPY TAGS COPY = (COMMENTS, GRANTS) AGGREGATION POLICY p ROW ACCESS POLICY r ON (a) WITH CONTACT (steward = u)`,
	} {
		assertRoundTrip(t, sql)
	}
}

// A view column's projection policy round-trips.
func TestDeparse_CreateView_ColumnProjectionPolicy(t *testing.T) {
	out := assertRoundTrip(t, `CREATE VIEW v (a WITH PROJECTION POLICY p) AS SELECT 1 AS a`)
	if !strings.Contains(out, "WITH PROJECTION POLICY p") {
		t.Errorf("deparsed SQL missing projection policy: %q", out)
	}
}
