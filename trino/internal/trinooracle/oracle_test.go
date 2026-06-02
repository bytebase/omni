package trinooracle

import (
	"context"
	"testing"
	"time"
)

// testOracle connects to a live Trino server for differential testing. It
// prefers $TRINO_ORACLE_URL, falling back to DefaultURL. The test is skipped
// (not failed) in -short mode or when no oracle is reachable, so the suite stays
// green without a running Trino while exercising the real classifier when one is
// up. Start one with:
//
//	docker run -d --name trino-oracle -p 18080:8080 trinodb/trino:latest
func testOracle(t *testing.T) *Oracle {
	t.Helper()
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := Connect("")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ver, err := o.Ping(ctx)
	if err != nil {
		t.Skipf("trino oracle not reachable at %s (start: docker run -d -p 18080:8080 %s): %v", o.baseURL, DefaultImage, err)
	}
	t.Logf("connected to Trino %s at %s", ver, o.baseURL)
	return o
}

// TestOracleClassification proves the oracle distinguishes a Trino parser
// rejection (SYNTAX_ERROR) from grammar-accepted statements that fail for
// semantic reasons (missing table/catalog/column, unsupported feature) or
// succeed. This is the contract every grammar node relies on.
func TestOracleClassification(t *testing.T) {
	o := testOracle(t)
	cases := []struct {
		name     string
		sql      string
		accepted bool
	}{
		// Accepted: valid syntax.
		{"valid-trivial", "SELECT 1", true},
		{"valid-tpch", "SELECT name FROM tpch.sf1.nation", true},
		{"show-catalogs", "SHOW CATALOGS", true},
		{"explain", "EXPLAIN SELECT 1", true},
		{"ddl-create-memory", "CREATE TABLE memory.default.omni_oracle_probe (x integer)", true},
		// Trino-specific syntax must be accepted.
		{"match-recognize", "SELECT * FROM tpch.sf1.orders MATCH_RECOGNIZE (PARTITION BY orderkey ORDER BY orderdate MEASURES A.totalprice AS p PATTERN (A+) DEFINE A AS true)", true},
		{"group-by-rollup", "SELECT nationkey, count(*) FROM tpch.sf1.nation GROUP BY ROLLUP (nationkey)", true},
		{"lambda", "SELECT filter(ARRAY[1,2,3], x -> x > 1)", true},
		// Accepted by the parser but rejected at a semantic stage (NOT syntax).
		{"semantic-no-table", "SELECT * FROM no_such_table_xyz", true},
		{"semantic-no-catalog", "SELECT * FROM nocat.nosch.notbl", true},
		{"semantic-no-column", "SELECT bogus_col FROM tpch.sf1.nation", true},
		// Rejected: genuine syntax errors.
		{"syntax-misspell", "SELCT 1", false},
		{"syntax-dangling-from", "SELECT * FROM", false},
		{"syntax-bad-paren", "SELECT (1 +", false},
		{"syntax-ddl-garbage", "CREATE TABLE memory.default.t (x intGAR-BAGE !!!)", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			res, err := o.CheckSyntax(ctx, tc.sql)
			if err != nil {
				t.Fatalf("CheckSyntax(%q): %v", tc.sql, err)
			}
			if res.Accepted != tc.accepted {
				t.Errorf("CheckSyntax(%q): accepted=%v (errorName=%q code=%d), want accepted=%v",
					tc.sql, res.Accepted, res.ErrorName, res.ErrorCode, tc.accepted)
			}
		})
	}
}
