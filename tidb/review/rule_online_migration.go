package review

import "github.com/bytebase/omni/review"

// checkOnlineMigration reports a change that asks for online migration.
// Bytebase runs online migration through gh-ost on any engine whose sheet
// has the directive, and gh-ost reads the server's binary log, which TiDB
// does not have: gh-ost's check of it, select @@global.log_bin, fails on
// TiDB with an unknown system variable. So on TiDB the request itself is
// the finding.
//
// gh-ost (bytebase/gh-ost2): go/logic/inspect.go — validateBinlogs.
func checkOnlineMigration(r *reporter) {
	r.report(review.OnlineMigrationUnsupported, -1, review.Range{}, "TiDB doesn't support online migration.\n\ngh-ost needs a binary log, which TiDB doesn't have. Turn off **Online migration** for this change.")
}
