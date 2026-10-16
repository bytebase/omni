package review

import "github.com/bytebase/omni/review"

// checkOnlineMigration reports a change that asks for online migration.
// Bytebase runs online migration through gh-ost, which only exists for
// the MySQL family, so on SQL Server the request itself is the finding.
func checkOnlineMigration(r *reporter) {
	r.report(review.OnlineMigration, -1, review.Range{}, "the change requests online migration, which SQL Server does not support")
}
