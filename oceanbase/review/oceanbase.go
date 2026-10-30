package review

import (
	"context"

	"github.com/bytebase/omni/review"
)

// ReviewOceanBase evaluates opts.Rules over sql, a change for OceanBase's
// MySQL mode, for every target, as Review does for MySQL. The rules are
// MySQL's, inherited with this tree, but for PRIOR_BACKUP, which Bytebase
// does not take on OceanBase.
//
// WALK_THROUGH is the MySQL walk with OceanBase's deltas, each verified
// on OceanBase CE 4.3.5: docs/specs/2026-10-27-walk-through-oceanbase.md.
//
// REQUIRE_PRIMARY_KEY does not judge a DROP COLUMN against the synced
// key, as the MySQL rule does: OceanBase refuses to drop a column of the
// primary key, so the statement never takes the key.
//
// oceanbase: https://github.com/oceanbase/oceanbase-doc/blob/V4.3.5/en-US/700.reference/500.sql-reference/100.sql-syntax/200.common-tenant-of-mysql-mode/600.sql-statement-of-mysql-mode/1600.alter-table-of-mysql-mode.md — DROP [COLUMN]
//
// ONLINE_MIGRATION is the MySQL rule too: Bytebase runs gh-ost on any
// engine whose sheet has the directive. OceanBase 4.4 reports log_bin ON
// with ROW binlogs, which gh-ost checks, but its observer rejects SHOW
// MASTER STATUS, which gh-ost's streamer reads; a deployment with
// OceanBase's Binlog Service behind obproxy serves it. The review cannot
// tell the two apart, so it does not report the request itself.
//
// bytebase: backend/plugin/advisor/builtin_rules.go;
// backend/plugin/parser/oceanbase and backend/plugin/advisor/oceanbase — the
// OceanBase registrations.
func ReviewOceanBase(ctx context.Context, sql string, opts review.Options, targets []review.Target) (*review.Result, error) {
	return reviewDialect(ctx, sql, opts, targets, dialect{oceanBase: true})
}

// checkPriorBackupOceanBase reports a change that asks for prior backup on
// OceanBase. Bytebase takes no prior backup on it and skips the backup
// without an error, so the change runs without one: the request itself is
// the finding.
//
// bytebase: backend/common/engine.go — EngineSupportPriorBackup;
// backend/runner/taskrun/database_migrate_executor.go — runStandardMigration.
func checkPriorBackupOceanBase(r *reporter) {
	r.report(review.PriorBackup, -1, review.Range{}, "the change enables prior backup, which Bytebase does not take on OceanBase: the change runs without a backup")
}
