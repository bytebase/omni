package review

import "github.com/bytebase/omni/review"

// checkPriorBackup reports a change that asks for prior backup. Bytebase
// has no prior backup for Snowflake: EngineSupportPriorBackup excludes it,
// so the task skips the backup and runs the change without one, and no
// parser registers a TransformDMLToSelect for it. The request itself is the
// finding.
//
// bytebase: backend/common/engine.go — EngineSupportPriorBackup;
// backend/runner/taskrun/database_migrate_executor.go — the
// EngineSupportPriorBackup check before backupData;
// backend/plugin/parser/{mysql,tidb,tsql,plsql,pg}/backup.go — the
// RegisterTransformDMLToSelect calls, none for Snowflake
func checkPriorBackup(r *reporter) {
	r.report(review.PriorBackup, -1, review.Range{}, "the change enables prior backup, which Bytebase does not take on Snowflake: the change runs without a backup")
}
