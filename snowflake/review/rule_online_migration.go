package review

import "github.com/bytebase/omni/review"

// checkOnlineMigration reports a change that asks for online migration.
// Bytebase runs a sheet that carries the gh-ost directive through gh-ost
// whatever the engine, and gh-ost migrates MySQL tables over the MySQL
// protocol, so on Snowflake the request itself is the finding.
//
// bytebase: backend/runner/taskrun/database_migrate_executor.go — the
// ghost.IsGhostEnabled branch to runGhostMigration;
// backend/component/ghost/config.go — NewMigrationContext
func checkOnlineMigration(r *reporter) {
	r.report(review.OnlineMigration, -1, review.Range{}, "the change requests online migration, which Snowflake does not support")
}
