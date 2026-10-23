package review

import (
	"context"

	"github.com/bytebase/omni/review"
)

// ReviewOceanBase evaluates opts.Rules over sql, a change for OceanBase's
// MySQL mode, for every target, as Review does for MySQL. Bytebase reads
// OceanBase as the MySQL dialect: it parses it with the MySQL parser and
// reviews it with the MySQL advisors, so every rule here is the MySQL one,
// but for PRIOR_BACKUP, which Bytebase does not take on OceanBase, and
// WALK_THROUGH, which waits for an oracle run against OceanBase.
//
// ONLINE_MIGRATION is the MySQL rule too: Bytebase runs gh-ost on any
// engine whose sheet has the directive. OceanBase 4.4 reports log_bin ON
// with ROW binlogs, which gh-ost checks, but its observer rejects SHOW
// MASTER STATUS, which gh-ost's streamer reads; a deployment with
// OceanBase's Binlog Service behind obproxy serves it. The review cannot
// tell the two apart, so it does not report the request itself.
//
// bytebase: backend/plugin/advisor/builtin_rules.go;
// backend/plugin/parser/mysql — the OceanBase registrations.
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
