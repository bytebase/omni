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
	r.report(review.PriorBackupUnsupported, -1, review.Range{}, "Bytebase doesn't take a prior backup on Snowflake.\n\nThe change runs without one. Turn off **Prior backup** for this change.")
}
