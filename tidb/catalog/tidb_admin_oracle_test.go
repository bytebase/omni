package catalog

import "testing"

// TestAlterRangeMatchesTiDB runs one ALTER RANGE script on the catalog and on the
// server, and holds the two to the same answer for every statement. TiDB's codes and
// the catalog's differ for policies, as they do for ALTER DATABASE, so it compares
// whether each statement succeeded, which is what a change script's next statement
// depends on.
func TestAlterRangeMatchesTiDB(t *testing.T) {
	tc := startTiDBForCatalog(t)
	script := []string{
		`CREATE PLACEMENT POLICY omni_range_p PRIMARY_REGION="r" REGIONS="r"`,
		"ALTER RANGE global PLACEMENT POLICY omni_range_nosuch",
		"ALTER RANGE foo PLACEMENT POLICY DEFAULT",
		"ALTER RANGE foo PLACEMENT POLICY omni_range_nosuch",
		"ALTER RANGE global PLACEMENT POLICY omni_range_p",
		"DROP PLACEMENT POLICY omni_range_p",
		"ALTER RANGE GLOBAL PLACEMENT POLICY = `default`",
		"ALTER RANGE META PLACEMENT POLICY = 'OMNI_RANGE_P'",
		"DROP PLACEMENT POLICY omni_range_p",
		"ALTER RANGE `meta` PLACEMENT POLICY SET DEFAULT",
		"DROP PLACEMENT POLICY omni_range_p",
		"ALTER RANGE global PLACEMENT POLICY omni_range_p",
	}
	cleanup := func() {
		for _, sql := range []string{
			"ALTER RANGE global PLACEMENT POLICY DEFAULT",
			"ALTER RANGE meta PLACEMENT POLICY DEFAULT",
			"DROP PLACEMENT POLICY IF EXISTS omni_range_p",
		} {
			if _, err := tc.db.ExecContext(tc.ctx, sql); err != nil {
				t.Errorf("cleanup %q: %v", sql, err)
			}
		}
	}
	cleanup()
	t.Cleanup(cleanup)

	c := setupWithDB(t)
	for _, sql := range script {
		_, serverErr := tc.db.ExecContext(tc.ctx, sql)
		catalogErr := execOne(t, c, sql)
		if (serverErr == nil) != (catalogErr == nil) {
			t.Errorf("%s:\n  server:  %v\n  catalog: %v", sql, serverErr, catalogErr)
		}
	}
}
