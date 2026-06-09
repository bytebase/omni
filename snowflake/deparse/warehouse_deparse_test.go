package deparse_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/bytebase/omni/snowflake/parser"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// WAREHOUSE DDL round-trips (gap-warehouse)
//
// assertRoundTrip compares ASTs (Loc-stripped), so cosmetic normalizations the
// deparser applies — dropping the optional leading WITH, uppercasing option
// names/keys, emitting DROP TABLES as REMOVE TABLES (the same AST action) — are
// expected and verified to be AST-equivalent.
// ---------------------------------------------------------------------------

func TestDeparse_CreateWarehouse(t *testing.T) {
	assertRoundTrip(t, "CREATE WAREHOUSE my_wh")
	assertRoundTrip(t, "CREATE OR REPLACE WAREHOUSE my_wh WITH WAREHOUSE_SIZE = 'X-LARGE'")
	assertRoundTrip(t, "CREATE OR REPLACE WAREHOUSE my_wh WAREHOUSE_SIZE = LARGE INITIALLY_SUSPENDED = TRUE")
	assertRoundTrip(t, "CREATE WAREHOUSE so_warehouse WITH WAREHOUSE_TYPE = 'SNOWPARK-OPTIMIZED' WAREHOUSE_SIZE = XLARGE RESOURCE_CONSTRAINT = 'MEMORY_16X_x86'")
	assertRoundTrip(t, "CREATE OR ALTER WAREHOUSE so_warehouse WAREHOUSE_TYPE = 'SNOWPARK-OPTIMIZED' AUTO_RESUME = TRUE COMMENT = 'Snowpark warehouse for ingestion'")
	assertRoundTrip(t, "CREATE WAREHOUSE w AUTO_SUSPEND = 60 GENERATION = '1'")
	assertRoundTrip(t, "CREATE WAREHOUSE IF NOT EXISTS w WAREHOUSE_SIZE = SMALL WITH TAG (cost_center = 'eng', env = 'prod')")
}

func TestDeparse_AlterWarehouse(t *testing.T) {
	assertRoundTrip(t, "ALTER WAREHOUSE IF EXISTS wh1 RENAME TO wh2")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh SUSPEND")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh RESUME")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh RESUME IF SUSPENDED")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh ABORT ALL QUERIES")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh SET WAREHOUSE_SIZE = MEDIUM")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh SET GENERATION = '2' AUTO_SUSPEND = 120")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh UNSET STATEMENT_TIMEOUT_IN_SECONDS")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh UNSET AUTO_SUSPEND, COMMENT")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh SET TAG cost_center = 'eng', env = 'prod'")
	assertRoundTrip(t, "ALTER WAREHOUSE my_wh UNSET TAG cost_center, env")
	assertRoundTrip(t, "ALTER WAREHOUSE interactive_demo ADD TABLES (orders, customers)")
	assertRoundTrip(t, "ALTER WAREHOUSE interactive_demo REMOVE TABLES (orders, customers)")
}

// TestDeparse_AlterWarehouse_DropTablesNormalizes documents that DROP TABLES and
// REMOVE TABLES parse to the same action and the deparser emits the REMOVE
// spelling — exercising the exact deparsed string, not just AST equivalence.
func TestDeparse_AlterWarehouse_DropTablesNormalizes(t *testing.T) {
	f, err := parser.Parse("ALTER WAREHOUSE w DROP TABLES (a, b)")
	require.NoError(t, err)
	out, err := deparse.DeparseFile(f)
	require.NoError(t, err)
	require.Equal(t, "ALTER WAREHOUSE w REMOVE TABLES (a, b)", out)
}
