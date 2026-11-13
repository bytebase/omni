package parser

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// TiDB's backup, restore, recover and flashback statements, read against a
// live pingcap/tidb:v8.5.5.
//
//	BACKUP  {DATABASE|SCHEMA} {* | db, ...} TO   'storage' [option ...]
//	BACKUP  TABLE tbl, ...                  TO   'storage' [option ...]
//	BACKUP  LOGS                            TO   'storage' [option ...]
//	RESTORE {DATABASE|SCHEMA} {* | db, ...} FROM 'storage' [option ...]
//	RESTORE TABLE tbl, ...                  FROM 'storage' [option ...]
//	RESTORE POINT                           FROM 'storage' [option ...]
//	RECOVER TABLE {BY JOB job_id | tbl [job_num]}
//	FLASHBACK TABLE tbl              [TO new_name]
//	FLASHBACK {DATABASE|SCHEMA} db   [TO new_db]
//	FLASHBACK CLUSTER                 TO {TIMESTAMP 'ts' | TSO n}
//	FLASHBACK TABLE tbl, ...          TO {TIMESTAMP 'ts' | TSO n}
//	FLASHBACK {DATABASE|SCHEMA} db    TO {TIMESTAMP 'ts' | TSO n}
//
// None of it is MySQL's. MySQL dropped BACKUP TABLE and RESTORE TABLE in 5.5,
// and mysql:8.0.46 and mariadb:11.8.9 each answer 1064 to all of the forms
// above (read on 2026-11-04); OceanBase has no equivalent either. omni's
// mysql, mariadb and oceanbase parsers refuse every one of them, so none of
// this is ported to them.
//
// The corpus writes 47 of these statements and almost nothing malformed, so
// the rejections are the point of this file: the family is a set of similar
// statements that differ in small ways, and one permissive reader would take
// every combination without the corpus noticing. Several things about the
// shape needed asking rather than reading:
//
//   - `FLASHBACK TABLE t TO timestamp` is a rename to a table called
//     timestamp, not a point-in-time flashback missing its value. Upstream
//     decides this in the lexer: TO TIMESTAMP is one token only when a
//     string literal follows it, and TO TSO only when an integer does.
//   - Every keyword in the option list has to be written bare. A backquoted
//     `CONCURRENCY`, `MB`, `minute`, `true` or `off` is a syntax error, and
//     BACKEND's and ON_DUPLICATE's value is the only place a backquoted name
//     is taken.
//   - SNAPSHOT n <unit> AGO takes only units with a fixed length, because
//     the grammar's action converts the unit to a duration and fails the
//     parse when it cannot: MONTH, QUARTER and YEAR are out, WEEK is in.
//   - A TSO must be above zero, which is also checked during the parse.
//
// Every case was measured against pingcap/tidb:v8.5.5 on 2026-11-04, and the
// expectation is what that server answered, so an image upgrade that changes
// a verdict fails here instead of passing silently.
//
// Ref: TiDB v8.5.5 parser.y:5623-5747 (BRIEStmt, BRIETables), 5759-6039 (the
// option list), 3039-3059 (RecoverTableStmt), 3067-3124
// (FlashbackToTimestampStmt), 3134-3166 (FlashbackTableStmt,
// FlashbackDatabaseStmt); pkg/parser/lexer.go:272-290 (toTimestamp, toTSO).
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-backup/
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-flashback-table/

// brieCase is a statement and whether TiDB's parser accepts it.
type brieCase struct {
	sql    string
	accept bool
}

var brieCases = []brieCase{
	// BACKUP and RESTORE targets. SCHEMA is DATABASE, the asterisk is the bare
	// token and stands alone, and the TABLE form takes a list of names.
	{"BACKUP DATABASE * TO 'noop://'", true},
	{"BACKUP SCHEMA * TO 'noop://'", true},
	{"BACKUP DATABASE a TO 'noop://'", true},
	{"BACKUP DATABASE a, b, c TO 'noop://'", true},
	{"BACKUP SCHEMA a TO 'noop://'", true},
	{"BACKUP TABLE a TO 'noop://'", true},
	{"BACKUP TABLE a.b TO 'noop://'", true},
	{"BACKUP TABLE a.b, c.d, e TO 'noop://'", true},
	{"BACKUP TABLE `a`.`b` TO 'noop://'", true},
	{"BACKUP TABLE *.a TO 'noop://'", true},
	{"BACKUP DATABASE `*` TO 'noop://'", true},
	{"RESTORE DATABASE * FROM 'noop://'", true},
	{"RESTORE DATABASE a, b FROM 'noop://'", true},
	{"RESTORE SCHEMA a FROM 'noop://'", true},
	{"RESTORE TABLE a.b FROM 'noop://'", true},
	{"BACKUP DATABASE * TO \"noop://\"", true},
	// Neither target is optional, the asterisk mixes with nothing, a database
	// name is never qualified, and a table name is never a wildcard.
	{"BACKUP TO 'noop://'", false},
	{"BACKUP DATABASE TO 'noop://'", false},
	{"BACKUP TABLE TO 'noop://'", false},
	{"BACKUP DATABASE *, a TO 'noop://'", false},
	{"BACKUP DATABASE a, * TO 'noop://'", false},
	{"BACKUP DATABASE a.b TO 'noop://'", false},
	{"BACKUP TABLE * TO 'noop://'", false},
	{"BACKUP TABLE a.* TO 'noop://'", false},
	{"BACKUP DATABASE a, TO 'noop://'", false},
	{"BACKUP DATABASES a TO 'noop://'", false},
	// The destination is mandatory, is one string literal, and is reached by TO
	// for BACKUP and by FROM for RESTORE and never the other way round.
	{"BACKUP DATABASE a", false},
	{"BACKUP DATABASE a TO", false},
	{"BACKUP DATABASE a TO 'noop' '://'", false},
	{"BACKUP DATABASE a TO noop", false},
	{"BACKUP DATABASE a TO 1", false},
	{"BACKUP DATABASE a FROM 'noop://'", false},
	{"RESTORE DATABASE a TO 'noop://'", false},
	{"BACKUP DATABASE a TO 'noop://', b TO 'noop://'", false},
	// BACKUP LOGS and RESTORE POINT name nothing and take the same option list.
	{"BACKUP LOGS TO 'noop://'", true},
	{"BACKUP LOGS TO 'noop://' start_ts = '20220304'", true},
	{"BACKUP LOGS", false},
	{"BACKUP LOGS TO", false},
	{"BACKUP LOGS a TO 'noop://'", false},
	{"RESTORE POINT FROM 'noop://'", true},
	{"RESTORE POINT FROM 'noop://' full_backup_storage = 'noop://full'", true},
	{"RESTORE POINT FROM 'noop://' full_backup_storage='noop://full' start_ts='20230101' restored_ts='20230123'", true},
	{"RESTORE POINT", false},
	{"RESTORE POINT TO 'noop://'", false},
	{"RESTORE POINTS FROM 'noop://'", false},
	// The option list may be empty, repeats an option without complaint, and
	// accepts an equals sign in front of every value or none.
	{"BACKUP TABLE a TO 'noop://' concurrency 4", true},
	{"BACKUP TABLE a TO 'noop://' concurrency = 4", true},
	{"BACKUP TABLE a TO 'noop://' CONCURRENCY 4", true},
	{"BACKUP TABLE a TO 'noop://' concurrency 4 concurrency 8", true},
	{"BACKUP TABLE a TO 'noop://' checksum_concurrency 4 compression_level 4 ignore_stats 1 compression_type 'lz4'", true},
	{"RESTORE TABLE a FROM 'noop://' concurrency 40 checksum 0 online 1", true},
	{"RESTORE TABLE a FROM 'noop://' checksum_concurrency 4 wait_tiflash_ready 1 with_sys_table 1", true},
	{"BACKUP TABLE a TO 'noop://' resume 1", true},
	// An option name is one word the grammar knows, written bare.
	{"BACKUP TABLE a TO 'noop://' nosuchopt 1", false},
	{"BACKUP TABLE a TO 'noop://' `concurrency` 4", false},
	{"BACKUP TABLE a TO 'noop://' concurrency", false},
	{"BACKUP TABLE a TO 'noop://' 4", false},
	{"BACKUP TABLE a TO 'noop://' select 1", false},
	// An integer option takes LengthNum: unsigned, no sign, inside uint64.
	{"BACKUP TABLE a TO 'noop://' concurrency 0", true},
	{"BACKUP TABLE a TO 'noop://' concurrency 18446744073709551615", true},
	{"BACKUP TABLE a TO 'noop://' concurrency 18446744073709551616", false},
	{"BACKUP TABLE a TO 'noop://' concurrency -1", false},
	{"BACKUP TABLE a TO 'noop://' concurrency 'x'", false},
	{"BACKUP TABLE a TO 'noop://' concurrency 1.5", false},
	// A boolean option takes a number or the word TRUE or FALSE.
	{"BACKUP TABLE a TO 'noop://' online 1", true},
	{"BACKUP TABLE a TO 'noop://' online true", true},
	{"BACKUP TABLE a TO 'noop://' online FALSE", true},
	{"BACKUP TABLE a TO 'noop://' online = TRUE", true},
	{"BACKUP TABLE a TO 'noop://' online 2", true},
	{"BACKUP TABLE a TO 'noop://' online `true`", false},
	{"BACKUP TABLE a TO 'noop://' online 'x'", false},
	{"BACKUP TABLE a TO 'noop://' online -0", false},
	{"BACKUP TABLE a TO 'noop://' send_credentials_to_tikv = 1", true},
	{"BACKUP TABLE a TO 'noop://' checkpoint 1 skip_schema_files 1 strict_format 1", true},
	{"BACKUP TABLE a TO 'noop://' csv_not_null 1 csv_backslash_escape 1 csv_trim_last_separators 1", true},
	{"BACKUP TABLE a TO 'noop://' load_stats 0", true},
	// A string option takes one string literal and nothing else.
	{"BACKUP TABLE a TO 'noop://' compression_type 'lz4'", true},
	{"BACKUP TABLE a TO 'noop://' compression_type lz4", false},
	{"BACKUP TABLE a TO 'noop://' tikv_importer = 'x'", true},
	{"BACKUP TABLE a TO 'noop://' encryption_method 'aes128-ctr' encryption_keyfile 'f'", true},
	{"BACKUP TABLE a TO 'noop://' csv_separator ',' csv_delimiter '\"' csv_null 'N'", true},
	{"BACKUP TABLE a TO 'noop://' gc_ttl 'x' until_ts 'x' start_ts 'x' restored_ts 'x'", true},
	{"BACKUP TABLE a TO 'noop://' full_backup_storage 'x'", true},
	{"BACKUP TABLE a TO 'noop://' compression_type 1", false},
	// BACKEND and ON_DUPLICATE take a string, an identifier, or IGNORE or
	// REPLACE. They are the one place a backquoted name is a value here.
	{"BACKUP TABLE a TO 'noop://' backend = 'local'", true},
	{"BACKUP TABLE a TO 'noop://' backend = local", true},
	{"BACKUP TABLE a TO 'noop://' backend local", true},
	{"BACKUP TABLE a TO 'noop://' backend = `ignore`", true},
	{"BACKUP TABLE a TO 'noop://' backend = ignore", true},
	{"BACKUP TABLE a TO 'noop://' backend = replace", true},
	{"BACKUP TABLE a TO 'noop://' `backend` = 'local'", false},
	{"BACKUP TABLE a TO 'noop://' backend = 1", false},
	{"BACKUP TABLE a TO 'noop://' on_duplicate = ignore", true},
	{"BACKUP TABLE a TO 'noop://' ON DUPLICATE = ignore", true},
	{"BACKUP TABLE a TO 'noop://' ON `DUPLICATE` = ignore", false},
	{"BACKUP TABLE a TO 'noop://' ON = ignore", false},
	// SNAPSHOT takes a timestamp, a TSO, or a duration back from now. The unit
	// must be one with a fixed length, so MONTH, QUARTER and YEAR are out.
	{"BACKUP DATABASE * TO 'noop://' snapshot = '2020-03-18 18:13:54'", true},
	{"BACKUP DATABASE * TO 'noop://' snapshot = 1234567890", true},
	{"BACKUP DATABASE * TO 'noop://' snapshot 5 minute ago", true},
	{"BACKUP DATABASE * TO 'noop://' snapshot = 5 sql_tsi_hour ago", true},
	{"BACKUP DATABASE * TO 'noop://' snapshot 5 microsecond ago", true},
	{"BACKUP DATABASE * TO 'noop://' snapshot 5 week ago", true},
	{"BACKUP DATABASE * TO 'noop://' snapshot 5 month ago", false},
	{"BACKUP DATABASE * TO 'noop://' snapshot 5 year ago", false},
	{"BACKUP DATABASE * TO 'noop://' snapshot 5 sql_tsi_month ago", false},
	{"BACKUP DATABASE * TO 'noop://' snapshot 5 minute", false},
	{"BACKUP DATABASE * TO 'noop://' snapshot 5 `minute` ago", false},
	{"BACKUP DATABASE * TO 'noop://' snapshot 5 minute `ago`", false},
	{"BACKUP DATABASE * TO 'noop://' snapshot ago", false},
	{"BACKUP DATABASE * TO 'noop://' snapshot 18446744073709551616", false},
	// An option list is not order-dependent. A numeric SNAPSHOT or
	// LAST_BACKUP ends at the number, so the word after it starts the
	// next option rather than a duration; a unit, once written, still
	// needs its AGO.
	{"BACKUP DATABASE * TO 'noop://' SNAPSHOT 1234567890 CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' CONCURRENCY 4 SNAPSHOT 1234567890", true},
	{"BACKUP DATABASE * TO 'noop://' SNAPSHOT 1234567890 ONLINE 1", true},
	{"BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 SNAPSHOT 6", true},
	{"BACKUP DATABASE * TO 'noop://' SNAPSHOT = 5 minute ago CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 day ago ONLINE 1 CONCURRENCY 4 checksum off", true},
	{"BACKUP DATABASE * TO 'noop://' SNAPSHOT 1234567890 week ago", true},
	{"BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 nosuchopt 4", false},
	{"BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 microsecond CONCURRENCY 4", false},
	{"BACKUP DATABASE * TO 'noop://' LAST_BACKUP 5 CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' csv_header 1 CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' csv_header fields CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' checksum 1 CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' checksum off CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' online 1 CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' rate_limit 500 MB/second CONCURRENCY 4", true},
	{"BACKUP DATABASE * TO 'noop://' concurrency 4 backend local", true},
	{"BACKUP DATABASE * TO 'noop://' backend local CONCURRENCY 4", true},
	// BACKEND's value is an Identifier, and CONCURRENCY is an
	// unreserved keyword upstream (parser.y:7097), so this reads
	// CONCURRENCY as the value and is left with a bare 4 that starts
	// no option: 1064.
	{"BACKUP DATABASE * TO 'noop://' backend CONCURRENCY 4", false},
	// LAST_BACKUP takes a timestamp or a TSO, and no duration.
	{"BACKUP DATABASE * TO 'noop://' LAST_BACKUP = '2020-02-02 14:14:14'", true},
	{"BACKUP DATABASE * TO 'noop://' LAST_BACKUP = 1234567890", true},
	{"BACKUP DATABASE * TO 'noop://' LAST_BACKUP 1234567890", true},
	{"BACKUP DATABASE * TO 'noop://' LAST_BACKUP 5 minute ago", false},
	{"BACKUP DATABASE * TO 'noop://' LAST_BACKUP", false},
	// RATE_LIMIT is spelled in MB per second and in nothing else.
	{"BACKUP DATABASE * TO 'noop://' rate_limit 500 MB/second", true},
	{"BACKUP DATABASE * TO 'noop://' rate_limit = 500 MB / second", true},
	{"BACKUP DATABASE * TO 'noop://' rate_limit 500 mb/second", true},
	{"BACKUP DATABASE * TO 'noop://' rate_limit 500 MB", false},
	{"BACKUP DATABASE * TO 'noop://' rate_limit 500 GB/second", false},
	{"BACKUP DATABASE * TO 'noop://' rate_limit 500 MB/minute", false},
	{"BACKUP DATABASE * TO 'noop://' rate_limit 500 `MB`/second", false},
	{"BACKUP DATABASE * TO 'noop://' rate_limit 500", false},
	// CSV_HEADER takes a number or the word FIELDS or COLUMNS.
	{"BACKUP TABLE a TO 'noop://' csv_header = 1", true},
	{"BACKUP TABLE a TO 'noop://' csv_header = fields", true},
	{"BACKUP TABLE a TO 'noop://' csv_header COLUMNS", true},
	{"BACKUP TABLE a TO 'noop://' csv_header = `fields`", false},
	{"BACKUP TABLE a TO 'noop://' csv_header = 'fields'", false},
	// CHECKSUM and ANALYZE each take a boolean or a level.
	{"RESTORE TABLE g FROM 'noop://' checksum off", true},
	{"RESTORE TABLE g FROM 'noop://' checksum optional", true},
	{"RESTORE TABLE g FROM 'noop://' checksum required", true},
	{"RESTORE TABLE g FROM 'noop://' checksum 0", true},
	{"RESTORE TABLE g FROM 'noop://' checksum true", true},
	{"RESTORE TABLE g FROM 'noop://' analyze off", true},
	{"RESTORE TABLE g FROM 'noop://' analyze 1", true},
	{"RESTORE TABLE g FROM 'noop://' checksum `off`", false},
	{"RESTORE TABLE g FROM 'noop://' checksum 'off'", false},
	{"RESTORE TABLE g FROM 'noop://' checksum maybe", false},
	// The whole URL, however long, is still one string literal.
	{"backup database * to 's3://bucket/path/?provider=alibaba&region=us-west-9&storage-class=glacier&sse=AES256&acl=authenticated-read&use-accelerate-endpoint=1' send_credentials_to_tikv = 1", true},
	{"backup table x to 's3://bucket/path/?endpoint=https://test-cluster-s3.local&access-key=aaaaaaaaa&secret-access-key=bbbbbbbb&force-path-style=1'", true},
	{"restore database * from 'gcs://bucket/path/?endpoint=https://test-cluster.gcs.local&storage-class=coldline&predefined-acl=OWNER&credentials-file=/data/private/creds.json'", true},
	// RECOVER TABLE BY JOB names one DDL job and no table.
	{"RECOVER TABLE BY JOB 11", true},
	{"recover table by job 0", true},
	{"RECOVER TABLE BY JOB 9223372036854775807", true},
	{"RECOVER TABLE BY JOB", false},
	{"RECOVER TABLE BY JOB -1", false},
	{"RECOVER TABLE BY JOB 9223372036854775808", false},
	{"RECOVER TABLE BY JOB 11,12,13", false},
	{"RECOVER TABLE BY JOB 'x'", false},
	{"RECOVER TABLE BY JOB t1", false},
	{"RECOVER TABLE BY 11", false},
	// RECOVER TABLE names one table and at most one number after it.
	{"RECOVER TABLE t1", true},
	{"RECOVER TABLE t1 100", true},
	{"RECOVER TABLE t1 0", true},
	{"RECOVER TABLE a.b", true},
	{"RECOVER TABLE a.b 100", true},
	{"RECOVER TABLE *.a", true},
	{"RECOVER TABLE t1,t2", false},
	{"RECOVER TABLE t1 100 200", false},
	{"RECOVER TABLE t1 abc", false},
	{"RECOVER TABLE t1 -1", false},
	{"RECOVER TABLE", false},
	{"RECOVER", false},
	{"RECOVER DATABASE d", false},
	{"RECOVER TABLE t1 'x'", false},
	// FLASHBACK TABLE and FLASHBACK DATABASE bring back a dropped object,
	// optionally under another name. The new name is a bare identifier.
	{"flashback table t", true},
	{"flashback table t TO t1", true},
	{"flashback table a.b to c", true},
	{"flashback table `t` to `t1`", true},
	{"flashback database db1", true},
	{"flashback schema db1", true},
	{"flashback database db1 to db2", true},
	{"flashback schema db1 to db2", true},
	{"flashback table a.b to c.d", false},
	{"flashback database a.b", false},
	{"flashback schema d to e.f", false},
	{"flashback table t,t1", false},
	{"flashback table t,t1 to t2", false},
	{"flashback table t to", false},
	{"flashback table", false},
	{"flashback database", false},
	{"flashback", false},
	{"flashback table t to 'x'", false},
	{"flashback table t to 1", false},
	// TO TIMESTAMP and TO TSO are the point-in-time forms, and only when the
	// literal that upstream's lexer looks for is there. Without it, the word
	// after TO is the new name, so `TO timestamp` and `TO tso` are renames.
	{"flashback cluster to timestamp '2021-05-26 16:45:26'", true},
	{"flashback cluster to tso 445494955052105721", true},
	{"flashback table t to timestamp '2021-05-26 16:45:26'", true},
	{"flashback table t,t1 to timestamp '2021-05-26 16:45:26'", true},
	{"flashback table t to tso 445494955052105722", true},
	{"flashback table t,t1 to tso 445494955052105723", true},
	{"flashback database test to timestamp '2021-05-26 16:45:26'", true},
	{"flashback schema test to timestamp '2021-05-26 16:45:26'", true},
	{"flashback database test to tso 445494955052105724", true},
	{"flashback schema test to tso 445494955052105725", true},
	{"flashback cluster to timestamp \"2021-05-26 16:45:26\"", true},
	{"flashback table t TO timestamp", true},
	{"flashback table t TO tso", true},
	{"flashback table t TO `timestamp`", true},
	{"flashback database d TO timestamp", true},
	{"flashback database d TO tso", true},
	{"flashback table t to timestamp 'a' 'b'", false},
	{"flashback table t to TIMESTAMP", true},
	{"flashback table to timestamp '2021-05-26 16:45:26'", false},
	{"flashback database to timestamp '2021-05-26 16:45:26'", false},
	{"flashback table to tso 445494955052105726", false},
	{"flashback database to tso 445494955052105727", false},
	// The cluster form has no new-name arm, so every TO it takes is a point.
	{"flashback cluster", false},
	{"flashback cluster to timestamp", false},
	{"flashback cluster to tso", false},
	{"flashback cluster to c", false},
	{"flashback cluster to timestamp 'x' to tso 1", false},
	{"flashback cluster to timestamp DATE_SUB(NOW(), INTERVAL 3 SECOND)", false},
	{"flashback cluster to timestamp TIDB_BOUNDED_STALENESS(DATE_SUB(NOW(), INTERVAL 3 SECOND), NOW())", false},
	{"flashback `cluster` to timestamp '2021-05-26 16:45:26'", false},
	// A TSO is an unsigned literal above zero that fits in a uint64.
	{"flashback schema test to tso 18446744073709551615", true},
	{"flashback schema test to tso 0", false},
	{"flashback schema test to tso -100", false},
	{"flashback schema test to tso 18446744073709551616", false},
	{"flashback table t to tso 1.5", false},
	{"flashback table t to tso +1", false},
	// FLASHBACK is an unreserved keyword and stays a name everywhere else.
	{"CREATE TABLE flashback (flashback int)", true},
	{"SELECT flashback FROM flashback", true},
	{"SELECT * FROM t flashback", true},
	{"SELECT * FROM t AS flashback", true},
	{"DELETE FROM t flashback WHERE flashback.a = 1", true},
	{"UPDATE t SET flashback = 1", true},
	{"CREATE INDEX flashback ON t (a)", true},
	{"CREATE DATABASE flashback", true},
	{"SELECT 1 AS flashback", true},
	{"SELECT * FROM flashback.t", true},
	// So are BACKUP, RESTORE and RECOVER, which this family did not change.
	{"CREATE TABLE t (backup int, restore int, recover int)", true},
	{"SELECT * FROM t backup", true},
	{"SELECT * FROM t restore", true},
	{"SELECT * FROM t recover", true},
	{"DELETE FROM t backup WHERE backup.a = 1", true},
}

// TestBRIEFormsOracle asserts omni agrees with the pinned TiDB on brieCases
// (container test). The verdict comes from PREPARE, which parses the
// statement without running it: 1064 is the parser's rejection, and any other
// server error means the text parsed and failed later.
func TestBRIEFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range brieCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(c.sql)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbRelease, accepts, c.accept, c.sql)
			}
			_, perr := Parse(c.sql)
			if omniAccepts := perr == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbRelease, c.accept, c.sql, perr)
			}
		})
	}
}

// TestBRIEForms holds omni to the same verdicts without a server, so a run
// with no Docker still reports a widened or narrowed grammar.
func TestBRIEForms(t *testing.T) {
	for _, c := range brieCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", accepts, c.accept, c.sql, err)
			}
		})
	}
}

// qualified writes a table reference the way the tree tests name one.
func qualified(ref *ast.TableRef) string {
	if ref == nil {
		return "<nil>"
	}
	if ref.Schema != "" {
		return ref.Schema + "." + ref.Name
	}
	return ref.Name
}

// parseOne parses one statement and fails unless that is all there was.
func parseOne(t *testing.T, sql string) ast.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) read %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}
