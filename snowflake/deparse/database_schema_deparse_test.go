package deparse_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/bytebase/omni/snowflake/parser"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// DATABASE / SCHEMA DDL round-trips for the documented forms: open-ended
// properties, catalog-linked databases, contacts, clone/backup/replica sources,
// versioned schemas and the newer ALTER DATABASE actions.
// ---------------------------------------------------------------------------

func TestDeparse_CreateDatabaseDocumented(t *testing.T) {
	assertRoundTrip(t, "CREATE OR ALTER DATABASE d DATA_RETENTION_TIME_IN_DAYS = 1 EXTERNAL_VOLUME = v LOG_LEVEL = 'INFO' OBJECT_VISIBILITY = PRIVILEGED")
	assertRoundTrip(t, "CREATE DATABASE d LINKED_CATALOG = (CATALOG = 'cat', ALLOWED_NAMESPACES = ('a', 'b'), SYNC_INTERVAL_SECONDS = 60) CATALOG_CASE_SENSITIVITY = CASE_INSENSITIVE")
	assertRoundTrip(t, `CREATE DATABASE d WITH TAG (t = 'x') WITH CONTACT (SUPPORT = c1, ACCESS_APPROVAL = db.sc."Mixed") OBJECT_VISIBILITY = PRIVILEGED`)
	assertRoundTrip(t, "CREATE DATABASE d CLONE s IGNORE TABLES WITH INSUFFICIENT DATA RETENTION IGNORE HYBRID TABLES")
	assertRoundTrip(t, "CREATE DATABASE d FROM LISTING 'GZ1Z2Z3'")
	assertRoundTrip(t, "CREATE DATABASE d FROM BACKUP SET bs IDENTIFIER 'id1'")
	assertRoundTrip(t, "CREATE DATABASE d AS REPLICA OF acct.primary_db DATA_RETENTION_TIME_IN_DAYS = 1")
}

func TestDeparse_CreateSchemaDocumented(t *testing.T) {
	assertRoundTrip(t, "CREATE SCHEMA s CLONE t IGNORE HYBRID TABLES WITH MANAGED ACCESS CATALOG = c WITH CONTACT (SUPPORT = c1)")
	assertRoundTrip(t, "CREATE SCHEMA s FROM BACKUP SET bs IDENTIFIER 'id1'")
	assertRoundTrip(t, "CREATE OR ALTER VERSIONED SCHEMA s WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 1 COMMENT = 'x'")
}

func TestDeparse_AlterDatabaseDocumented(t *testing.T) {
	assertRoundTrip(t, "ALTER DATABASE d SET DATA_RETENTION_TIME_IN_DAYS = 1 CONTACT SUPPORT = c1, ACCESS_APPROVAL = c2 ENABLE_DATA_COMPACTION = TRUE")
	assertRoundTrip(t, "ALTER DATABASE d UNSET CONTACT SUPPORT, COMMENT")
	assertRoundTrip(t, "ALTER DATABASE d SET FEATURE POLICY p FORCE")
	assertRoundTrip(t, "ALTER DATABASE d UNSET FEATURE POLICY")
	assertRoundTrip(t, "ALTER DATABASE IF EXISTS d SUSPEND DISCOVERY")
	assertRoundTrip(t, "ALTER DATABASE d RESUME DISCOVERY")
	assertRoundTrip(t, "ALTER DATABASE d UPDATE LINKED_CATALOG ADD ('a', 'b') TO ALLOWED_NAMESPACES")
	assertRoundTrip(t, "ALTER DATABASE d UPDATE LINKED_CATALOG REMOVE ('a') FROM BLOCKED_NAMESPACES")
	assertRoundTrip(t, "ALTER DATABASE d UPDATE LINKED_CATALOG UNSET ALLOWED_NAMESPACES")
	assertRoundTrip(t, "ALTER DATABASE d UPDATE LINKED_CATALOG SET ALLOWED_WRITE_OPERATIONS = NONE")
	assertRoundTrip(t, "ALTER DATABASE d SET TAG t1 = 'a', t2 = 'b'")
	assertRoundTrip(t, "ALTER DATABASE d UNSET TAG t1, t2")
}

func TestDeparse_AlterSchemaDocumented(t *testing.T) {
	assertRoundTrip(t, "ALTER SCHEMA s SET LOG_LEVEL = 'INFO' CONTACT SUPPORT = c STORAGE_SERIALIZATION_POLICY = OPTIMIZED")
	assertRoundTrip(t, "ALTER SCHEMA s UNSET CONTACT SUPPORT, REPLACE_INVALID_CHARACTERS")
	assertRoundTrip(t, "ALTER SCHEMA s SET TAG t = 'a'")
	assertRoundTrip(t, "ALTER SCHEMA s UNSET TAG t")
}

// TestDeparse_DBSchemaTagsUnparenthesized pins that ALTER ... SET/UNSET TAG is
// written in the documented, unparenthesized form even when the input used the
// parenthesized one omni also accepts.
func TestDeparse_DBSchemaTagsUnparenthesized(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER DATABASE d SET TAG (t = 'a')": "ALTER DATABASE d SET TAG t = 'a'",
		"ALTER SCHEMA s UNSET TAG (t1, t2)":  "ALTER SCHEMA s UNSET TAG t1, t2",
	} {
		f, err := parser.Parse(in)
		require.NoError(t, err)
		out, err := deparse.DeparseFile(f)
		require.NoError(t, err)
		require.Equal(t, want, out)
	}
}
