package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TiDB's RECOMMEND INDEX, SET CONFIG, ALTER RANGE, CANCEL, KILL and FLUSH, and the log
// task and BR job statements, in the WALK_THROUGH walk.
//
// None creates, drops, or renames a name the walk keeps, so each is lowered the way
// ANALYZE TABLE is: followed, recording nothing. What that has to earn is the statement
// after it. A node the walk does not lower returns nil from lower(), which ends the walk
// there, so one KILL in the middle of a change script would hide every error below it.
//
// Each can still fail on the server, on something the walk does not hold: an unknown
// advisor option (1105), an unknown type or instance for SET CONFIG (1105), a range or
// placement policy (1105, 8239), a plugin (1105), a job, or a cluster with no TiKV (1105
// for the log tasks on the harness's single tidb-server). TestWalkThroughAgainstTiDB
// holds the forms the single server does run to the engine; the rest are pinned here.
//
// These tests need no container and run unconditionally.

// tidbAdminWalkStatements are the forms the walk must follow.
var tidbAdminWalkStatements = []string{
	`RECOMMEND INDEX RUN`,
	`RECOMMEND INDEX RUN FOR 'select * from t' WITH max_num_index = 1`,
	`RECOMMEND INDEX SHOW`,
	`RECOMMEND INDEX APPLY 1`,
	`RECOMMEND INDEX IGNORE 1`,
	`RECOMMEND INDEX SET max_num_index = 1`,
	`SET CONFIG tikv log.level = 'info'`,
	`SET CONFIG '127.0.0.1:3306' auto-compaction-mode = TRUE`,
	`ALTER RANGE global PLACEMENT POLICY mypolicy`,
	`ALTER RANGE meta PLACEMENT POLICY DEFAULT`,
	`CANCEL DISTRIBUTION JOB 1`,
	`KILL 23123`,
	`KILL TIDB 23123`,
	`KILL TIDB CONNECTION 23123`,
	`KILL TIDB QUERY 23123`,
	`KILL CONNECTION_ID()`,
	`FLUSH TIDB PLUGINS a, b`,
	`STOP BACKUP LOGS`,
	`PAUSE BACKUP LOGS gc_ttl = '20220304'`,
	`RESUME BACKUP LOGS`,
	`PURGE BACKUP LOGS FROM 'noop://' until_ts = '2012122304'`,
	`SHOW BACKUP LOGS STATUS`,
	`SHOW BACKUP LOGS METADATA FROM 'noop://'`,
	`SHOW BACKUP METADATA FROM 'noop://'`,
	`SHOW BR JOB 1234`,
	`SHOW BR JOB QUERY 1234`,
	`CANCEL BR JOB 1234`,
}

func TestWalkThroughFollowsTiDBAdminStatements(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range tidbAdminWalkStatements {
		t.Run(sql, func(t *testing.T) {
			// On its own it reports nothing: it changes no name.
			noFinding(t, sql+";", shop)

			// The walk carries on past it, so the error below is still found,
			// and it is reported against statement 2.
			oneFinding(t, sql+";\nALTER TABLE t ADD COLUMN c int;\nALTER TABLE t ADD COLUMN c int;",
				2, "Column `c` already exists in `t`.", shop)

			// And the names a statement above it set are still there below it,
			// so the walk did not lose its state either.
			noFinding(t, "ALTER TABLE t ADD COLUMN c int;\n"+sql+";\nALTER TABLE t DROP COLUMN c;", shop)
			oneFinding(t, "DROP TABLE t;\n"+sql+";\nALTER TABLE t ADD COLUMN c int;",
				2, "Table `t` doesn't exist.", shop)
		})
	}
}

// TestWalkThroughClaimsNothingAboutWhatTheStatementNames pins that the walk does not
// read a name out of these statements as if it were one it holds. RECOMMEND INDEX RUN
// FOR answers 1146 on v8.5.5 for a table its text does not find, but that text is a
// string the parser never opens; `t` in ALTER RANGE t is a range, not the table.
func TestWalkThroughClaimsNothingAboutWhatTheStatementNames(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "RECOMMEND INDEX RUN FOR 'select * from nosuch';", shop)
	noFinding(t, "ALTER RANGE t PLACEMENT POLICY nosuch;", shop)
	noFinding(t, "SET CONFIG t nosuch.item = 1;", shop)
	noFinding(t, "FLUSH TIDB PLUGINS t;", shop)
	// Nothing here makes the name either: the table is still the one that is missing.
	oneFinding(t, "ALTER RANGE nosuch PLACEMENT POLICY nosuch;\nDROP TABLE nosuch;",
		1, "Table `nosuch` doesn't exist.", shop)
}

// TestWalkThroughStillEndsAtTheOtherBRIEKinds pins the line the log task kinds draw:
// BACKUP fails on a table that is not there and RESTORE brings names back, so the walk
// still ends at them, where it follows SHOW BR JOB beside them.
func TestWalkThroughStillEndsAtTheOtherBRIEKinds(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	for _, sql := range []string{
		`BACKUP DATABASE * TO 'noop://'`,
		`BACKUP LOGS TO 'noop://'`,
		`RESTORE DATABASE * FROM 'noop://'`,
		`RESTORE POINT FROM 'noop://'`,
	} {
		t.Run(sql, func(t *testing.T) {
			noFinding(t, sql+";\nALTER TABLE t DROP COLUMN zz;", shop)
		})
	}
}
