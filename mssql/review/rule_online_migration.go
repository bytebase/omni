package review

import "github.com/bytebase/omni/review"

// checkOnlineMigration reports a change that asks for online migration.
// Bytebase runs online migration through gh-ost, which only exists for
// the MySQL family, so on SQL Server the request itself is the finding.
func checkOnlineMigration(r *reporter) {
	r.report(review.OnlineMigrationUnsupported, -1, review.Range{}, "SQL Server doesn't support online migration.\n\nTurn off **Online migration** for this change.")
}
