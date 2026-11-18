package deparse_test

import "testing"

// Objects read by their documented shape (ast.ObjectStmt) round-trip.
func TestDeparse_ObjectStmt(t *testing.T) {
	for _, sql := range []string{
		"CREATE COMPUTE POOL IF NOT EXISTS p FOR APPLICATION app MIN_NODES = 1 MAX_NODES = 2 INSTANCE_FAMILY = CPU_X64_XS AUTO_RESUME = TRUE WITH TAG (db.s.t = 'v') COMMENT = 'c''s' BACKUP_INSTANCE_FAMILIES = ('a', 'b')",
		"CREATE COMPUTE POOL \"My Pool\" MIN_NODES = 1 MAX_NODES = 1 INSTANCE_FAMILY = \"fam\"",
		"ALTER COMPUTE POOL IF EXISTS p SUSPEND",
		"ALTER COMPUTE POOL p RESUME",
		"ALTER COMPUTE POOL p STOP ALL OF TYPE a, b",
		"ALTER COMPUTE POOL p SET MIN_NODES = 1 TAG t = 's', u = 's'",
		"ALTER COMPUTE POOL p UNSET AUTO_SUSPEND_SECS, COMMENT",
		"DROP COMPUTE POOL IF EXISTS p",
		"ALTER CORTEX SEARCH SERVICE s SET AUTO_SUSPEND = NULL COMMENT = 'c'",
		"CREATE OR REPLACE TEMPORARY SECURE AGENT db.s.a COPY GRANTS COMMENT = 'c' FROM SPECIFICATION 's'",
		"CREATE SECURE AGENT a COMMENT = 'c' FROM @db.s.st/dir",
	} {
		assertRoundTrip(t, sql)
	}
}
