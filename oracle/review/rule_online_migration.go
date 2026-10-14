package review

import "github.com/bytebase/omni/review"

// checkOnlineMigration reports a change that asks for online migration.
// Bytebase runs every change whose sheet carries the gh-ost directive
// through gh-ost, which only exists for the MySQL family, so on Oracle the
// request itself is the finding.
//
// bytebase: backend/runner/taskrun/database_migrate_executor.go — RunOnce (ghost.IsGhostEnabled)
func checkOnlineMigration(r *reporter) {
	r.report(review.OnlineMigration, -1, review.Range{}, "the change requests online migration, which Oracle does not support")
}
