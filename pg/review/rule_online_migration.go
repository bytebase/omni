package review

import "github.com/bytebase/omni/review"

// checkOnlineMigration reports a change that asks for online migration.
// Bytebase runs online migration through gh-ost, which only exists for
// the MySQL family, so on PostgreSQL the request itself is the finding.
func checkOnlineMigration(change review.Change, r *reporter) {
	if !change.OnlineMigration {
		return
	}
	r.report(review.OnlineMigrationUnsupported, -1, review.Range{}, "PostgreSQL doesn't support online migration.\n\nTurn off **Online migration** for this change.")
}
