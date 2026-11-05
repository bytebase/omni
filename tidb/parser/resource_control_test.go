package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's resource control statements: CREATE, ALTER and DROP RESOURCE GROUP,
// SET RESOURCE GROUP, CALIBRATE RESOURCE, QUERY WATCH ADD and QUERY WATCH
// REMOVE, plus the RESOURCE GROUP clause CREATE USER and ALTER USER take.
//
// The corpus in testdata records whether a statement parses and nothing more,
// and it only holds the statements pingcap's own test writes, so it cannot see
// an option that takes a value it should refuse: of the forms below it covers
// the accepting side well and the refusing side barely. Every case here is
// checked against the parser with no container and against a live TiDB v8.5.5
// in TestResourceControlMatchTiDB, so each verdict is the engine's rather than
// a reading of its grammar file.
//
// The rules the rejections rest on, each read from v8.5.5:
//
//   - CREATE and ALTER RESOURCE GROUP require at least one option. A group
//     with none is 1064, and so is a repeat of an option the group carries.
//   - RU_PER_SEC takes an unsigned integer or UNLIMITED; PRIORITY takes LOW,
//     MEDIUM or HIGH and nothing else; BURSTABLE stands alone or takes a
//     number, TRUE or FALSE. The words MODERATED, UNLIMITED and OFF that
//     pingcap's parser on master reads for BURSTABLE are not v8.5.5's.
//   - QUERY_LIMIT and BACKGROUND take a parenthesised list, `()` or NULL. The
//     list may not be bare: `QUERY_LIMIT = EXEC_ELAPSED '10s'` is 1064.
//   - EXEC_ELAPSED and WATCH's DURATION are read by Go's time.ParseDuration,
//     so '1d' and '10' are not durations and '1h30m' and '-1s' are. DURATION
//     also takes the keyword UNLIMITED, and the string 'unlimited' stands for
//     it and so skips the check.
//   - PROCESSED_KEYS and RU take the grammar's intLit, which carries no sign, and
//     their value is Int64Num's, as a watch id is.
//   - A second ACTION or WATCH inside QUERY_LIMIT is 1064; a second
//     EXEC_ELAPSED, PROCESSED_KEYS or RU is not.
//   - CALIBRATE RESOURCE takes a workload or dynamic options, never both, and
//     WORKLOAD takes no `=`. Its DURATION is read by TiDB's own duration
//     parser, which has a 'd' unit Go's has not and refuses a bare number.
//   - QUERY WATCH REMOVE takes a watch id alone: the RESOURCE GROUP forms on
//     master are 1064 on v8.5.5, and so is a negative id.
//   - A user variable is a group name only after QUERY WATCH's RESOURCE
//     GROUP, not inside ACTION SWITCH_GROUP(). A system variable is a group
//     name nowhere: the grammar names UserVariable, which @@rg is not. A
//     backtick-quoted @rg is neither -- it is an ordinary identifier, and a
//     group name everywhere one is taken.
//   - WATCH's DURATION is lowercased before it is read, so '1H' and '1H30M'
//     are durations and 'UNLIMITED' in any case is the keyword.
//     EXEC_ELAPSED's is not lowercased, so '1H' there is a syntax error.
//   - A QUERY WATCH pattern is a SimpleExpr, which carries its own COLLATE
//     suffix but is not a whole expression: `TO 'a' COLLATE utf8mb4_bin`
//     parses and `TO 'a' AND 1` does not. It also carries SimpleExpr's JSON
//     path operators, whose operands are narrower than the rest of it: a
//     column reference on the left and a string literal on the right, and so
//     no chain of two.
//   - An INTERVAL's unit is a keyword, so a quoted one is not a unit:
//     `DURATION INTERVAL 15 MINUTE` parses and the same with MINUTE quoted
//     does not.
//   - MySQL's own resource group grammar is not TiDB's: TYPE, VCPU,
//     THREAD_PRIORITY, ENABLE, DISABLE, FORCE and SET RESOURCE GROUP's FOR
//     list are each 1064 here.
var resourceControlAccepted = []string{
	// CREATE RESOURCE GROUP: the name, IF NOT EXISTS, and DEFAULT as a name.
	`create resource group x ru_per_sec=1000`,
	`create resource group if not exists x ru_per_sec=1000`,
	"create resource group `default` ru_per_sec=100",
	// A backtick-quoted name holds whatever it likes, @ included.
	"create resource group `@rg` ru_per_sec=100",
	"alter resource group `@rg` ru_per_sec=100",
	"drop resource group `@rg`",
	"drop resource group `@`",
	"set resource group `@rg`",
	"set resource group `@@x`",
	"create resource group x ru_per_sec=1 QUERY_LIMIT=(ACTION SWITCH_GROUP(`@rg`))",
	`create resource group x RU_PER_SEC=1000`,
	`create resource group x ru_per_sec 1000`,
	`create resource group x ru_per_sec=0`,
	`create resource group x ru_per_sec=18446744073709551615`,
	`create resource group x ru_per_sec=UNLIMITED`,
	`create resource group x ru_per_sec=unlimited`,

	// PRIORITY's whole value set.
	`create resource group x priority=low`,
	`create resource group x priority=medium`,
	`create resource group x priority=high`,
	`create resource group x priority high`,

	// BURSTABLE bare and with each value v8.5.5 takes.
	`create resource group x burstable`,
	`create resource group x burstable=true`,
	`create resource group x burstable=false`,
	`create resource group x burstable=1`,
	`create resource group x burstable=0`,
	`create resource group x burstable=2`,
	`create resource group x burstable=9223372036854775807`,
	`create resource group x burstable true`,
	`create resource group x burstable 1`,

	// The separators: a comma between options, or nothing at all.
	`create resource group x ru_per_sec=20, priority=LOW, burstable`,
	`create resource group x ru_per_sec=20 priority=LOW burstable`,
	`create resource group x ru_per_sec=2000, BURSTABLE`,
	`create resource group x burstable ru_per_sec=4000`,
	`create resource group x priority=low burstable ru_per_sec=1 query_limit=() background=()`,

	// QUERY_LIMIT: the rules, the empty list, and NULL.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=()`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT ()`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=NULL`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT NULL`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED "10s")`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED='10s' PROCESSED_KEYS=100)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED='10s', PROCESSED_KEYS=100, RU=100)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU=100)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS=100)`,
	// Go's duration units, and nothing of TiDB's own.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1h30m')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1.5s')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '-1s')`,
	// A repeated rule is not a repeated option: v8.5.5 keeps the last.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1s' EXEC_ELAPSED '2s')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU 1, RU 2)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(PROCESSED_KEYS 1 PROCESSED_KEYS 2)`,

	// ACTION's whole value set, SWITCH_GROUP included.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION KILL EXEC_ELAPSED='10m')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION DRYRUN)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10m' ACTION COOLDOWN)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP(rg1))`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION=SWITCH_GROUP(rg1))`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP(default))`,
	"create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP(`default`))",

	// WATCH's types and its DURATION.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH=EXACT)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH SIMILAR DURATION '10m')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH PLAN DURATION='10m')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION UNLIMITED)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION = UNLIMITED)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION 'UNLIMITED')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '')`,
	// The action lowercases the string first, so the units may be
	// written in any case and so may the word UNLIMITED.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1H')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1M')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1MS')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1H30M')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION 'UnLiMiTeD')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION COOLDOWN WATCH EXACT DURATION '10m')`,

	// BACKGROUND: its two options, the empty list, and NULL.
	`create resource group x ru_per_sec=1000 background (task_types 'br,lightning')`,
	`create resource group x ru_per_sec=1000 background (task_types='br,lightning')`,
	`create resource group x ru_per_sec=1000 background = (task_types='')`,
	`create resource group x ru_per_sec=1000 background = (UTILIZATION_LIMIT=50)`,
	`create resource group x ru_per_sec=1000 background = (UTILIZATION_LIMIT 0)`,
	`create resource group x ru_per_sec=1000 background (task_types='br', utilization_limit=50)`,
	`create resource group x ru_per_sec=1000 background=()`,
	`create resource group x ru_per_sec=1000 background ()`,
	`create resource group x ru_per_sec=1000 background=NULL`,
	`create resource group x ru_per_sec=1000 background NULL`,

	// ALTER RESOURCE GROUP takes the same options, plus IF EXISTS.
	`alter resource group x ru_per_sec=1000`,
	`alter resource group if exists x ru_per_sec=1000`,
	`alter resource group default priority = high`,
	`alter resource group x BURSTABLE`,
	`alter resource group x BURSTABLE, ru_per_sec=3000`,
	`alter resource group x ru_per_sec=UNLIMITED, BURSTABLE`,
	`alter resource group x QUERY_LIMIT=()`,
	`alter resource group x QUERY_LIMIT=NULL`,
	`alter resource group x background NULL`,
	`alter resource group default priority=low background = ( task_types "ttl" )`,

	// DROP RESOURCE GROUP.
	`drop resource group x`,
	`drop resource group DEFAULT`,
	"drop resource group `default`",
	`drop resource group if exists x`,

	// SET RESOURCE GROUP.
	`set resource group x`,
	`set resource group default`,
	"set resource group `default`",
	"set resource group ``",

	// CALIBRATE RESOURCE: nothing, a workload, or the dynamic options.
	`calibrate resource`,
	`CALIBRATE RESOURCE`,
	`calibrate resource workload tpcc`,
	`calibrate resource workload oltp_read_write`,
	`calibrate resource workload oltp_read_only`,
	`calibrate resource workload oltp_write_only`,
	`calibrate resource workload tpch_10`,
	`calibrate resource START_TIME '2021-04-15 00:00:00'`,
	`calibrate resource START_TIME = '2021-04-15 00:00:00'`,
	`calibrate resource END_TIME now()`,
	`calibrate resource START_TIME '2023-04-01 13:00:00' END_TIME '2023-04-01 16:00:00' DURATION '20m'`,
	`calibrate resource START_TIME '2023-04-01 13:00:00',END_TIME='2023-04-01 16:00:00'`,
	`calibrate resource DURATION='20m' START_TIME '2023-04-01 13:00:00'`,
	// START_TIME and END_TIME take a whole expression.
	`calibrate resource START_TIME CURRENT_TIMESTAMP() - interval 20 minute`,
	`calibrate resource START_TIME now() - 1000 END_TIME current_timestamp()`,
	`calibrate resource START_TIME a`,
	`calibrate resource START_TIME t.a`,
	`calibrate resource START_TIME (select 1)`,
	`calibrate resource START_TIME @v`,
	`calibrate resource START_TIME -1`,
	`calibrate resource START_TIME null`,
	`calibrate resource START_TIME a is null`,
	`calibrate resource START_TIME a between 1 and 2`,
	`calibrate resource START_TIME exists (select 1)`,
	// DURATION's two forms, with TiDB's own duration units.
	`calibrate resource DURATION '20d'`,
	`calibrate resource DURATION '1d12h'`,
	`calibrate resource DURATION '1.5h'`,
	`calibrate resource DURATION ''`,
	`calibrate resource DURATION interval 15 minute`,
	`calibrate resource DURATION = interval 15 minute`,
	`calibrate resource DURATION interval '15' minute`,
	`calibrate resource DURATION interval 15 MINUTE`,
	`calibrate resource DURATION interval 15 day_hour`,
	`calibrate resource DURATION interval 15 sql_tsi_minute`,
	`calibrate resource START_TIME CURRENT_TIMESTAMP() - interval 20 minute DURATION interval 15 minute`,
	// Upstream's duplicate check reaches the first two options only, so a
	// repeat behind three of them parses.
	`calibrate resource START_TIME '2020-01-01' END_TIME '2020-01-02' DURATION '1m' DURATION '2m'`,

	// QUERY WATCH ADD: each option kind, with and without commas.
	`query watch add SQL TEXT SIMILAR to 'select 1'`,
	`query watch add SQL TEXT EXACT to 'select 1'`,
	`query watch add SQL TEXT PLAN to 'select 1'`,
	`query watch add SQL DIGEST 'b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377'`,
	`query watch add SQL DIGEST b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377`,
	"query watch add PLAN DIGEST `5e3ddd388f6012e328233dbcdda5d48f404e0536c6c54d9618233210f3d5762a`",
	`query watch add PLAN DIGEST @digest1`,
	`query watch add resource group rg1 SQL TEXT SIMILAR to 'select 1'`,
	"query watch add resource group `default` SQL TEXT SIMILAR to 'select 1'",
	`query watch add resource group default SQL DIGEST 'a'`,
	`query watch add resource group @rg SQL TEXT SIMILAR to @sql1`,
	`query watch add resource group rg1`,
	`query watch add SQL TEXT SIMILAR to 'select 1' resource group rg1`,
	`query watch add ACTION = KILL SQL TEXT SIMILAR to 'select 1'`,
	`query watch add ACTION COOLDOWN resource group rg1 SQL TEXT SIMILAR to 'select 1'`,
	`query watch add ACTION DRYRUN SQL DIGEST 'a'`,
	`query watch add ACTION SWITCH_GROUP(rg1) SQL DIGEST 'a'`,
	`query watch add ACTION KILL`,
	`query watch add SQL DIGEST 'a', resource group rg1`,
	`query watch add resource group rg1, action kill, sql digest 'a'`,
	`query watch add SQL TEXT SIMILAR to 'select 1' resource group rg1 action kill`,
	`query watch add SQL DIGEST concat('a','b')`,
	`query watch add SQL DIGEST -1`,
	// A pattern is a SimpleExpr, which takes a COLLATE suffix, and the
	// suffix repeats because the production is left-recursive.
	`query watch add SQL TEXT EXACT to 'select 1' COLLATE utf8mb4_bin`,
	`query watch add SQL DIGEST 'a' COLLATE utf8mb4_bin`,
	`query watch add PLAN DIGEST 'a' COLLATE utf8mb4_bin`,
	`query watch add SQL TEXT EXACT to 'select 1' COLLATE 'utf8mb4_bin'`,
	`query watch add SQL TEXT EXACT to 'select 1' COLLATE binary`,
	`query watch add SQL TEXT EXACT to 'select 1' COLLATE utf8mb4_bin COLLATE utf8mb4_general_ci`,
	`query watch add SQL DIGEST 'a' COLLATE utf8mb4_bin, resource group rg`,
	// SimpleExpr's JSON path operators, on a bare and a qualified column.
	`query watch add SQL DIGEST cfg->>'$.digest'`,
	`query watch add SQL DIGEST cfg->'$.digest'`,
	`query watch add SQL DIGEST d.cfg->'$.a'`,
	`query watch add SQL DIGEST a.b.cfg->'$.a'`,
	`query watch add PLAN DIGEST cfg->>'$.digest'`,
	`query watch add SQL TEXT EXACT to cfg->>'$.t'`,
	// COLLATE may follow the JSON suffix, and another option may follow both.
	`query watch add SQL DIGEST cfg->>'$.a' COLLATE utf8mb4_bin`,
	`query watch add SQL DIGEST cfg->>'$.a' resource group rg`,
	`query watch add SQL DIGEST cfg->>'$.a' action kill`,
	"query watch add resource group `@rg` SQL DIGEST 'a'",
	"query watch add resource group `@@x` SQL DIGEST 'a'",

	// QUERY WATCH REMOVE.
	`query watch remove 1`,
	`query watch remove 0`,
	`query watch remove 9223372036854775807`,
	`QUERY WATCH REMOVE 1`,

	// The RESOURCE GROUP clause on a user, which comes last.
	`CREATE USER 'root'@'127.0.0.1' IDENTIFIED BY 'hashstring' RESOURCE GROUP rg1`,
	`ALTER USER 'root'@'localhost' RESOURCE GROUP rg2`,
	`create user 'u'@'h' RESOURCE GROUP default`,
	"create user 'u'@'h' RESOURCE GROUP `rg`",
	`create user 'u'@'h' comment 'c' RESOURCE GROUP rg`,
	`create user 'u'@'h' account lock RESOURCE GROUP rg`,
	`create user 'u'@'h' require ssl RESOURCE GROUP rg`,
	`create user 'u'@'h' with max_queries_per_hour 1 RESOURCE GROUP rg`,
	`alter user if exists 'u'@'h' RESOURCE GROUP rg`,
	`alter user 'u'@'h' identified by 'p' RESOURCE GROUP rg`,

	// The words this family registers are unreserved: each stays an ordinary
	// identifier, alias, index name, savepoint name and database name.
	`create table priority (duration int, plan int, watch int, background int, ru int, exact int, similar int, digest int, unlimited int, low int, high int, start_time int, end_time int, burstable int, calibrate int, workload int, task_types int, dryrun int, cooldown int, switch_group int, tpcc int, query_limit int, exec_elapsed int, processed_keys int, utilization_limit int, oltp_read_write int, tpch_10 int)`,
	`select priority, duration, plan, watch from t`,
	`select t.plan, t.duration from t`,
	`select 1 as priority, 2 as duration, 3 as plan, 4 as watch, 5 as calibrate, 6 as unlimited`,
	`select * from t priority`,
	`select * from t as plan`,
	`select * from t duration join t2 workload on duration.a = workload.a`,
	`delete from t plan where plan.id=1`,
	`delete from t priority where priority.id=1`,
	`delete from t duration where duration.id=1`,
	`delete from t text where text.id=1`,
	`delete from t action where action.id=1`,
	`delete from t config where config.id=1`,
	`update t priority set priority.a=1`,
	`create index priority on t (a)`,
	`create index duration on t (a)`,
	`alter table t add index duration (a)`,
	`alter table t rename column a to priority`,
	`savepoint priority`,
	`savepoint duration`,
	`release savepoint plan`,
	`create database priority`,
	`create database duration`,
	`with priority as (select 1) select * from priority`,
	`create user duration`,
	`create user unlimited`,
	"create role `plan`",
	`select current_resource_group()`,
}

var resourceControlRejected = []string{
	// BURSTABLE's number, RU, PROCESSED_KEYS and a watch id are Int64Num, whose action
	// answers 1105 for what an int64 cannot hold: 2^63 is refused and 2^63-1 is not.
	`create resource group x burstable=18446744073709551615`,
	`create resource group x burstable=9223372036854775808`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 9223372036854775808)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 9223372036854775808)`,
	`query watch remove 9223372036854775808`,
	`query watch remove 18446744073709551615`,

	// A group with no options at all.
	`create resource group x`,
	`alter resource group x`,
	`create resource group if exists x ru_per_sec=1000`,

	// A repeated option.
	`create resource group x ru_per_sec=1000 ru_per_sec=2000`,
	`create resource group x priority=LOW ru_per_sec=1000 priority=HIGH`,
	`create resource group x ru_per_sec=1000, ru_per_sec=2000`,

	// An option name the grammar does not have, and a value of the wrong kind.
	`create resource group x foo=1`,
	`create resource group x cpu ='8c'`,
	`create resource group x region ='us, 3'`,
	`create resource group x followers=0`,
	`create resource group x cpu='8c', io_read_bandwidth='2GB/s', io_write_bandwidth='200MB/s'`,
	`create resource group x ru_per_sec='check'`,
	`create resource group x ru_per_sec=`,
	`create resource group x ru_per_sec=-1`,
	`create resource group x ru_per_sec=+1000`,
	`create resource group x ru_per_sec=1.5`,
	`create resource group x ru_per_sec=18446744073709551616`,
	`create resource group x priority=mid`,
	`create resource group x priority=1`,
	`create resource group x priority='low'`,
	`create resource group x ru_per_sec=20 priority=MID BURSTABLE`,
	// A quoted option name is an identifier, not the keyword.
	"create resource group x `ru_per_sec`=1000",
	// A trailing comma ends nothing.
	`create resource group x ru_per_sec=1000,`,

	// BURSTABLE's values stop at a number, TRUE and FALSE.
	`create resource group x burstable=disable`,
	`create resource group x burstable=unlimited`,
	`create resource group x burstable=moderated`,
	`create resource group x burstable=off`,
	`create resource group x burstable=on`,
	`create resource group x burstable=null`,
	`create resource group x burstable=default`,
	`create resource group x burstable='true'`,
	`create resource group x burstable=-1`,
	`create resource group x burstable=1.5`,
	`create resource group x burstable=18446744073709551616`,
	`alter resource group x BURSTABLE = UNLIMITED ru_per_sec=4000`,
	`alter resource group x ru_per_sec=2000, BURSTABLE=yes`,

	// QUERY_LIMIT needs its parentheses, and a name it has.
	`create resource group x ru_per_sec=1000 QUERY=(EXEC_ELAPSED '10s')`,
	`create resource group x ru_per_sec=1000 LIMIT=(EXEC_ELAPSED '10s')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=EXEC_ELAPSED '10s'`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT EXEC_ELAPSED '10s'`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s'`,
	// DURATION belongs to WATCH and appears nowhere else in the list.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(DURATION '1m')`,
	// The rules' value kinds.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED 'xyz')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1d')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED 100)`,
	// EXEC_ELAPSED's string is read as written, where WATCH's is
	// lowercased first.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1H')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1M')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1MS')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(PROCESSED_KEYS='100')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(PROCESSED_KEYS -1)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU='100')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU 1.5)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU -1)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU +1)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU 18446744073709551616)`,
	// A second ACTION or WATCH.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s' ACTION DRYRUN ACTION KILL)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 100 ACTION DRYRUN ACTION KILL)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT WATCH SIMILAR)`,
	`alter resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s' ACTION DRYRUN WATCH SIMILAR DURATION '10m' ACTION COOLDOWN)`,
	// ACTION's and WATCH's value sets.
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION STOP)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP())`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP(@rg))`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH ROUGH)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION 'xyz')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1d')`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION 10)`,
	`create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1D')`,

	// BACKGROUND's option names and value kinds.
	`create resource group x ru_per_sec=1000 background=(task_types br)`,
	`create resource group x ru_per_sec=1000 background=(task_types 1)`,
	`create resource group x ru_per_sec=1000 background=(utilization_limit '50')`,
	`create resource group x ru_per_sec=1000 background = (UTILIZATION_LIMIT="NAN")`,
	`create resource group x ru_per_sec=1000 background=(utilization_limit=-1)`,
	`create resource group x ru_per_sec=1000 background=(utilization_limit 18446744073709551616)`,
	`create resource group x ru_per_sec=1000 background=(task_types 'br' task_types 'lightning')`,
	`create resource group x ru_per_sec=1000 background=(foo 1)`,
	`create resource group x ru_per_sec=1000 background=task_types 'br'`,
	// A trailing comma ends no list, inner or outer.
	`create resource group x ru_per_sec=1000 query_limit=(exec_elapsed '1s',)`,
	`create resource group x ru_per_sec=1000 background=(task_types 'br',)`,
	`create resource group x ru_per_sec=1000 query_limit=(,)`,

	// DROP takes one group and SET takes one name.
	`drop resource group x,y`,
	`drop resource group if exists x,y`,
	`drop resource group`,
	`set resource group x y`,
	`set resource group`,

	// MySQL's resource group grammar, which TiDB has none of.
	`create resource group x type=user`,
	`create resource group x type=user vcpu=0-63 thread_priority=0 enable`,
	`create resource group x vcpu=0-63`,
	`create resource group x thread_priority=0`,
	`create resource group x ru_per_sec=1000 enable`,
	`alter resource group x force`,
	`alter resource group x enable`,
	`drop resource group x force`,
	`set resource group x for 1`,

	// CALIBRATE RESOURCE.
	`calibrate resource foo`,
	`calibrate resource workload`,
	`calibrate resource workload foo`,
	`calibrate resource workload = tpcc`,
	`calibrate resource workload = oltp_read_write START_TIME '2023-04-01 13:00:00'`,
	`calibrate resource workload tpcc START_TIME '2023-04-01 13:00:00'`,
	`calibrate resource START_TIME '2023-04-01 13:00:00' workload tpcc`,
	`calibrate resource workload tpcc, workload tpcc`,
	`calibrate resource START_TIME`,
	`calibrate resource DURATION`,
	`calibrate resource DURATION 20`,
	`calibrate resource DURATION 'xyz'`,
	`calibrate resource DURATION '20'`,
	`calibrate resource DURATION '1w'`,
	`calibrate resource DURATION '-1m'`,
	`calibrate resource DURATION 'unlimited'`,
	`calibrate resource DURATION interval 15`,
	// A TimeUnit is a keyword: a quoted word is an identifier, not a unit.
	"calibrate resource DURATION interval 15 `MINUTE`",
	"calibrate resource DURATION interval 15 `day_hour`",
	`calibrate resource DURATION interval 15 'minute'`,
	`calibrate resource START_TIME '2020-01-01',`,
	// A repeat the check reaches.
	`calibrate resource START_TIME '2020-01-01' START_TIME '2020-01-02'`,
	`calibrate resource START_TIME '2020-01-01' END_TIME '2020-01-02' END_TIME '2020-01-03'`,
	`calibrate resource DURATION '1m' DURATION '2m'`,
	`calibrate resource DURATION '1m' START_TIME '2020-01-01' DURATION '2m'`,

	// QUERY WATCH.
	`query watch add`,
	`query watch add SQL SIMILAR to 'select 1'`,
	`query watch add SQL TEXT SIMILAR 'select 1'`,
	`query watch add SQL TEXT SIMILAR to`,
	`query watch add SQL TEXT FOO to 'select 1'`,
	`query watch add PLAN DIGEST`,
	`query watch add SQL DIGEST 'a' SQL DIGEST 'b'`,
	`query watch add SQL DIGEST 'a' PLAN DIGEST 'b'`,
	"query watch add resource group `default` resource group `rg1` SQL TEXT SIMILAR to 'select 1'",
	`query watch add ACTION KILL ACTION DRYRUN SQL DIGEST 'a'`,
	`query watch add ACTION SWITCH_GROUP(@rg) SQL DIGEST 'a'`,
	// RESOURCE GROUP takes a user variable, not a system variable.
	`query watch add RESOURCE GROUP @@rg SQL DIGEST 'a'`,
	`query watch add RESOURCE GROUP @@global.rg SQL DIGEST 'a'`,
	`query watch add SQL DIGEST 'a',`,
	`query watch add SQL DIGEST 'a' extra`,
	`query watch add SQL TEXT EXACT to 'select 1' COLLATE`,
	// The JSON path operators take a column reference and a string and
	// nothing else, so their result ends the pattern.
	`query watch add SQL DIGEST cfg->>1`,
	`query watch add SQL DIGEST cfg->1`,
	`query watch add SQL DIGEST 'a'->>'$.x'`,
	`query watch add SQL DIGEST (cfg)->>'$.x'`,
	`query watch add SQL DIGEST now()->>'$.a'`,
	`query watch add SQL DIGEST @v->>'$.a'`,
	`query watch add SQL DIGEST cfg->'$.a'->'$.b'`,
	`query watch add SQL DIGEST cfg->>'$.d'->>'$.e'`,
	`query watch add SQL DIGEST cfg COLLATE utf8mb4_bin ->> '$.x'`,
	`query watch add SQL DIGEST cfg->>`,
	`query watch remove`,
	`query watch remove -1`,
	`query watch remove 'a'`,
	`query watch remove @v`,
	`query watch remove 1 2`,
	`query watch remove 18446744073709551616`,
	`query watch remove resource group rg1`,
	`query watch remove resource group @rg`,
	`query watch 1`,
	`drop query watch 1`,
	`drop query watch resource group rg1`,

	// The user clause comes last and comes once.
	`create user 'u'@'h' RESOURCE GROUP rg comment 'c'`,
	`create user 'u'@'h' RESOURCE GROUP rg account lock`,
	`create user 'u'@'h' RESOURCE GROUP rg RESOURCE GROUP rg2`,
	`create user 'u'@'h' RESOURCE GROUP 'rg'`,
	`create user 'u'@'h' RESOURCE GROUP`,
	`create user 'u'@'h' RESOURCE rg`,

	// A role written on its own is a plain identifier, never a keyword: that
	// is Rolename, which is not Username.
	`CREATE ROLE RESOURCE`,
	`create role plan`,
	`create role duration`,
	`create role priority`,
	`create role action`,
	`create role text`,
	`create role r1, resource`,
	`drop role resource`,
	`set role resource`,
	`set role all except resource`,
	`set default role resource to u1`,
}

func TestResourceControlAccepted(t *testing.T) {
	for _, sql := range resourceControlAccepted {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

func TestResourceControlRejected(t *testing.T) {
	for _, sql := range resourceControlRejected {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but TiDB v8.5.5 rejects it", sql)
			}
		})
	}
}

// TestResourceControlMatchTiDB holds the two tables above to the engine
// itself, so the verdicts above stay a measurement of v8.5.5 and not a
// reading of its grammar file.
func TestResourceControlMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for want, cases := range map[bool][]string{
		true:  resourceControlAccepted,
		false: resourceControlRejected,
	} {
		for _, sql := range cases {
			t.Run(sql, func(t *testing.T) {
				accepts, err := tc.prepareSyntax(sql)
				if err != nil {
					t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
				}
				if accepts != want {
					t.Errorf("TiDB %ss %q, this file records %s",
						verdictName(accepts), sql, verdictName(want))
				}
			})
		}
	}
}

// one parses sql and returns its single statement. The corpus records only
// whether a statement parses, so these cases assert the tree: an option list
// that parses into the wrong names and values would pass every test above.
func one(t *testing.T, sql string) nodes.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	if list == nil || len(list.Items) != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}

func TestResourceGroupOptionsAST(t *testing.T) {
	stmt, ok := one(t, `CREATE RESOURCE GROUP IF NOT EXISTS rg RU_PER_SEC = 2000, PRIORITY = HIGH, BURSTABLE = FALSE`).(*nodes.CreateResourceGroupStmt)
	if !ok {
		t.Fatalf("not a CreateResourceGroupStmt")
	}
	if !stmt.IfNotExists || stmt.Name != "rg" {
		t.Errorf("if_not_exists=%v name=%q, want true and rg", stmt.IfNotExists, stmt.Name)
	}
	if len(stmt.Options) != 3 {
		t.Fatalf("got %d options, want 3", len(stmt.Options))
	}
	if o := stmt.Options[0]; o.Name != "RU_PER_SEC" || o.IntValue != 2000 || o.Unlimited {
		t.Errorf("option 0 = %+v, want RU_PER_SEC 2000", o)
	}
	if o := stmt.Options[1]; o.Name != "PRIORITY" || o.Value != "HIGH" {
		t.Errorf("option 1 = %+v, want PRIORITY HIGH", o)
	}
	if o := stmt.Options[2]; o.Name != "BURSTABLE" || o.Burstable {
		t.Errorf("option 2 = %+v, want BURSTABLE false", o)
	}

	// RU_PER_SEC = UNLIMITED is a flag, not a rate, and bare BURSTABLE is on.
	stmt = one(t, `CREATE RESOURCE GROUP rg RU_PER_SEC = UNLIMITED BURSTABLE`).(*nodes.CreateResourceGroupStmt)
	if o := stmt.Options[0]; !o.Unlimited || o.IntValue != 0 {
		t.Errorf("RU_PER_SEC = UNLIMITED gave %+v, want unlimited with no rate", o)
	}
	if o := stmt.Options[1]; !o.Burstable {
		t.Errorf("bare BURSTABLE gave %+v, want true", o)
	}
	// A zero value turns BURSTABLE off.
	stmt = one(t, `CREATE RESOURCE GROUP rg BURSTABLE = 0`).(*nodes.CreateResourceGroupStmt)
	if o := stmt.Options[0]; o.Burstable {
		t.Errorf("BURSTABLE = 0 gave %+v, want false", o)
	}
}

func TestResourceGroupQueryLimitAST(t *testing.T) {
	stmt := one(t, `ALTER RESOURCE GROUP IF EXISTS rg QUERY_LIMIT = (EXEC_ELAPSED '10s' PROCESSED_KEYS = 100 RU 200 ACTION SWITCH_GROUP(rg2) WATCH SIMILAR DURATION '1m')`).(*nodes.AlterResourceGroupStmt)
	if !stmt.IfExists {
		t.Errorf("if_exists not set")
	}
	if len(stmt.Options) != 1 || stmt.Options[0].Name != "QUERY_LIMIT" {
		t.Fatalf("got options %+v, want one QUERY_LIMIT", stmt.Options)
	}
	rules := stmt.Options[0].Runaway
	if len(rules) != 5 {
		t.Fatalf("got %d rules, want 5", len(rules))
	}
	if r := rules[0]; r.Name != "EXEC_ELAPSED" || r.Value != "10s" {
		t.Errorf("rule 0 = %+v, want EXEC_ELAPSED 10s", r)
	}
	if r := rules[1]; r.Name != "PROCESSED_KEYS" || r.IntValue != 100 {
		t.Errorf("rule 1 = %+v, want PROCESSED_KEYS 100", r)
	}
	if r := rules[2]; r.Name != "RU" || r.IntValue != 200 {
		t.Errorf("rule 2 = %+v, want RU 200", r)
	}
	if r := rules[3]; r.Name != "ACTION" || r.Action != "SWITCH_GROUP" || r.SwitchGroup != "rg2" {
		t.Errorf("rule 3 = %+v, want ACTION SWITCH_GROUP(rg2)", r)
	}
	if r := rules[4]; r.Name != "WATCH" || r.Watch != "SIMILAR" || r.Duration != "1m" || r.DurationUnlimited {
		t.Errorf("rule 4 = %+v, want WATCH SIMILAR DURATION 1m", r)
	}

	// DURATION = UNLIMITED is the keyword, not a duration string, and the
	// string 'unlimited' in any case means the same thing: the action
	// lowercases it and keeps no duration for either, so both land on the
	// flag rather than one of them looking like a duration called
	// "unlimited".
	for _, sql := range []string{
		`ALTER RESOURCE GROUP rg QUERY_LIMIT = (WATCH EXACT DURATION = UNLIMITED)`,
		`ALTER RESOURCE GROUP rg QUERY_LIMIT = (WATCH EXACT DURATION 'unlimited')`,
		`ALTER RESOURCE GROUP rg QUERY_LIMIT = (WATCH EXACT DURATION 'UnLiMiTeD')`,
	} {
		got := one(t, sql).(*nodes.AlterResourceGroupStmt).Options[0].Runaway[0]
		if !got.DurationUnlimited || got.Duration != "" {
			t.Errorf("%s gave %+v, want the unlimited flag and no string", sql, got)
		}
	}
	// A duration written in upper case keeps its text and is still read.
	stmt = one(t, `ALTER RESOURCE GROUP rg QUERY_LIMIT = (WATCH EXACT DURATION '1H30M')`).(*nodes.AlterResourceGroupStmt)
	if r := stmt.Options[0].Runaway[0]; r.Duration != "1H30M" || r.DurationUnlimited {
		t.Errorf("DURATION '1H30M' gave %+v, want the duration as written", r)
	}

	// QUERY_LIMIT = NULL and QUERY_LIMIT = () are the option with no rules,
	// which is not the same as the option being absent.
	for _, sql := range []string{
		`ALTER RESOURCE GROUP rg QUERY_LIMIT = NULL`,
		`ALTER RESOURCE GROUP rg QUERY_LIMIT = ()`,
	} {
		stmt = one(t, sql).(*nodes.AlterResourceGroupStmt)
		if len(stmt.Options) != 1 || stmt.Options[0].Name != "QUERY_LIMIT" || stmt.Options[0].Runaway != nil {
			t.Errorf("%s gave %+v, want one QUERY_LIMIT with no rules", sql, stmt.Options)
		}
	}
}

func TestResourceGroupBackgroundAST(t *testing.T) {
	stmt := one(t, `CREATE RESOURCE GROUP rg RU_PER_SEC = 1 BACKGROUND = (TASK_TYPES 'br,lightning', UTILIZATION_LIMIT = 30)`).(*nodes.CreateResourceGroupStmt)
	opts := stmt.Options[1].Background
	if len(opts) != 2 {
		t.Fatalf("got %d background options, want 2", len(opts))
	}
	if o := opts[0]; o.Name != "TASK_TYPES" || o.Value != "br,lightning" {
		t.Errorf("background 0 = %+v, want TASK_TYPES br,lightning", o)
	}
	if o := opts[1]; o.Name != "UTILIZATION_LIMIT" || o.IntValue != 30 {
		t.Errorf("background 1 = %+v, want UTILIZATION_LIMIT 30", o)
	}
}

func TestDropAndSetResourceGroupAST(t *testing.T) {
	drop := one(t, `DROP RESOURCE GROUP IF EXISTS rg`).(*nodes.DropResourceGroupStmt)
	if !drop.IfExists || drop.Name != "rg" {
		t.Errorf("drop = %+v, want if_exists rg", drop)
	}
	// DEFAULT is a group name here, not the keyword's usual meaning.
	drop = one(t, `DROP RESOURCE GROUP DEFAULT`).(*nodes.DropResourceGroupStmt)
	if drop.Name != "DEFAULT" {
		t.Errorf("drop name = %q, want DEFAULT", drop.Name)
	}
	set := one(t, `SET RESOURCE GROUP rg`).(*nodes.SetResourceGroupStmt)
	if set.Name != "rg" {
		t.Errorf("set name = %q, want rg", set.Name)
	}
}

func TestCalibrateResourceAST(t *testing.T) {
	// The workload form.
	stmt := one(t, `CALIBRATE RESOURCE WORKLOAD TPCH_10`).(*nodes.CalibrateResourceStmt)
	if stmt.Workload != "TPCH_10" || len(stmt.Options) != 0 {
		t.Errorf("got workload %q with %d options, want TPCH_10 and none", stmt.Workload, len(stmt.Options))
	}
	// CALIBRATE RESOURCE alone carries neither.
	stmt = one(t, `CALIBRATE RESOURCE`).(*nodes.CalibrateResourceStmt)
	if stmt.Workload != "" || len(stmt.Options) != 0 {
		t.Errorf("got workload %q with %d options, want neither", stmt.Workload, len(stmt.Options))
	}

	// The dynamic form: START_TIME and END_TIME carry an expression,
	// DURATION a duration string.
	stmt = one(t, `CALIBRATE RESOURCE START_TIME now() - interval 20 minute, END_TIME = now(), DURATION '20d'`).(*nodes.CalibrateResourceStmt)
	if stmt.Workload != "" {
		t.Errorf("workload = %q, want empty", stmt.Workload)
	}
	if len(stmt.Options) != 3 {
		t.Fatalf("got %d options, want 3", len(stmt.Options))
	}
	if o := stmt.Options[0]; o.Name != "START_TIME" || o.Ts == nil {
		t.Errorf("option 0 = %+v, want START_TIME with an expression", o)
	}
	if _, ok := stmt.Options[0].Ts.(*nodes.BinaryExpr); !ok {
		t.Errorf("START_TIME now() - interval 20 minute gave %T, want a BinaryExpr", stmt.Options[0].Ts)
	}
	if o := stmt.Options[1]; o.Name != "END_TIME" || o.Ts == nil {
		t.Errorf("option 1 = %+v, want END_TIME with an expression", o)
	}
	if o := stmt.Options[2]; o.Name != "DURATION" || o.Value != "20d" || o.Ts != nil {
		t.Errorf("option 2 = %+v, want DURATION 20d", o)
	}

	// DURATION INTERVAL expr unit carries the expression and the unit
	// instead of a string, and the unit comes out normalized.
	stmt = one(t, `CALIBRATE RESOURCE DURATION INTERVAL 15 MINUTE`).(*nodes.CalibrateResourceStmt)
	if o := stmt.Options[0]; o.Name != "DURATION" || o.Unit != "MINUTE" || o.Ts == nil || o.Value != "" {
		t.Errorf("option = %+v, want DURATION INTERVAL with unit MINUTE", o)
	}
	stmt = one(t, `CALIBRATE RESOURCE DURATION INTERVAL 15 sql_tsi_minute`).(*nodes.CalibrateResourceStmt)
	if o := stmt.Options[0]; o.Unit != "MINUTE" {
		t.Errorf("unit = %q, want MINUTE, the unit SQL_TSI_MINUTE names", o.Unit)
	}
}

func TestQueryWatchAST(t *testing.T) {
	stmt := one(t, `QUERY WATCH ADD RESOURCE GROUP rg1 ACTION KILL SQL TEXT SIMILAR TO 'select 1'`).(*nodes.AddQueryWatchStmt)
	if len(stmt.Options) != 3 {
		t.Fatalf("got %d options, want 3", len(stmt.Options))
	}
	if o := stmt.Options[0]; o.Name != "RESOURCE_GROUP" || o.ResourceGroup != "rg1" || o.ResourceGroupExpr != nil {
		t.Errorf("option 0 = %+v, want RESOURCE_GROUP rg1", o)
	}
	if o := stmt.Options[1]; o.Name != "ACTION" || o.Action != "KILL" {
		t.Errorf("option 1 = %+v, want ACTION KILL", o)
	}
	o := stmt.Options[2]
	if o.Name != "SQL_TEXT" || o.Watch != "SIMILAR" || o.Pattern == nil {
		t.Errorf("option 2 = %+v, want SQL_TEXT SIMILAR with a pattern", o)
	}
	if lit, ok := o.Pattern.(*nodes.StringLit); !ok || lit.Value != "select 1" {
		t.Errorf("pattern = %#v, want the string literal select 1", o.Pattern)
	}

	// A user variable group name lands in the expression field, not the name.
	stmt = one(t, `QUERY WATCH ADD RESOURCE GROUP @rg SQL DIGEST 'abc'`).(*nodes.AddQueryWatchStmt)
	if o := stmt.Options[0]; o.ResourceGroup != "" || o.ResourceGroupExpr == nil {
		t.Errorf("option 0 = %+v, want the group as an expression", o)
	}
	if o := stmt.Options[1]; o.Name != "SQL_DIGEST" || o.Watch != "" {
		t.Errorf("option 1 = %+v, want SQL_DIGEST with no watch type", o)
	}
	// A quoted @rg is an identifier, so it lands the other way round. The
	// lexer gives both the same kind of token with the same Str, so a check
	// on the text alone would put this one in the expression field; the
	// verdict would still be accept and the tree would be wrong.
	stmt = one(t, "QUERY WATCH ADD RESOURCE GROUP `@rg` SQL DIGEST 'abc'").(*nodes.AddQueryWatchStmt)
	if o := stmt.Options[0]; o.ResourceGroup != "@rg" || o.ResourceGroupExpr != nil {
		t.Errorf("option 0 = %+v, want the group as the name @rg", o)
	}
	// The pattern's COLLATE suffix is part of the pattern.
	stmt = one(t, `QUERY WATCH ADD SQL TEXT EXACT TO 'select 1' COLLATE utf8mb4_bin`).(*nodes.AddQueryWatchStmt)
	coll, ok := stmt.Options[0].Pattern.(*nodes.CollateExpr)
	if !ok {
		t.Fatalf("pattern = %T, want a CollateExpr", stmt.Options[0].Pattern)
	}
	if coll.Collation != "utf8mb4_bin" {
		t.Errorf("collation = %q, want utf8mb4_bin", coll.Collation)
	}
	if lit, ok := coll.Expr.(*nodes.StringLit); !ok || lit.Value != "select 1" {
		t.Errorf("collated expression = %#v, want the string literal select 1", coll.Expr)
	}
	// A JSON path suffix comes out as the operator it is, with the column on
	// the left and the path string on the right, and COLLATE wraps all of it.
	stmt = one(t, `QUERY WATCH ADD SQL DIGEST t.cfg->>'$.digest' COLLATE utf8mb4_bin`).(*nodes.AddQueryWatchStmt)
	coll, ok = stmt.Options[0].Pattern.(*nodes.CollateExpr)
	if !ok {
		t.Fatalf("pattern = %T, want a CollateExpr", stmt.Options[0].Pattern)
	}
	bin, ok := coll.Expr.(*nodes.BinaryExpr)
	if !ok {
		t.Fatalf("collated expression = %T, want a BinaryExpr", coll.Expr)
	}
	if bin.Op != nodes.BinOpJsonUnquote {
		t.Errorf("op = %v, want the ->> operator", bin.Op)
	}
	if ref, ok := bin.Left.(*nodes.ColumnRef); !ok || ref.Table != "t" || ref.Column != "cfg" {
		t.Errorf("left = %#v, want the column t.cfg", bin.Left)
	}
	if lit, ok := bin.Right.(*nodes.StringLit); !ok || lit.Value != "$.digest" {
		t.Errorf("right = %#v, want the path string $.digest", bin.Right)
	}
	// `->` is the other operator, and carries no unquote.
	stmt = one(t, `QUERY WATCH ADD SQL DIGEST cfg->'$.a'`).(*nodes.AddQueryWatchStmt)
	if bin, ok := stmt.Options[0].Pattern.(*nodes.BinaryExpr); !ok || bin.Op != nodes.BinOpJsonExtract {
		t.Errorf("pattern = %#v, want the -> operator", stmt.Options[0].Pattern)
	}
	// PLAN DIGEST is its own name.
	stmt = one(t, `QUERY WATCH ADD PLAN DIGEST 'abc'`).(*nodes.AddQueryWatchStmt)
	if o := stmt.Options[0]; o.Name != "PLAN_DIGEST" {
		t.Errorf("option 0 = %+v, want PLAN_DIGEST", o)
	}

	remove := one(t, `QUERY WATCH REMOVE 42`).(*nodes.DropQueryWatchStmt)
	if remove.ID != 42 {
		t.Errorf("id = %d, want 42", remove.ID)
	}
}

func TestUserResourceGroupAST(t *testing.T) {
	create := one(t, `CREATE USER 'u'@'h' COMMENT 'c' RESOURCE GROUP rg1`).(*nodes.CreateUserStmt)
	if create.ResourceGroup != "rg1" || create.Comment != "c" {
		t.Errorf("create = resource_group %q comment %q, want rg1 and c", create.ResourceGroup, create.Comment)
	}
	alter := one(t, `ALTER USER 'u'@'h' RESOURCE GROUP DEFAULT`).(*nodes.AlterUserStmt)
	if alter.ResourceGroup != "DEFAULT" {
		t.Errorf("alter resource_group = %q, want DEFAULT", alter.ResourceGroup)
	}
	// Without the clause the field stays empty rather than picking up a
	// neighbouring option.
	create = one(t, `CREATE USER 'u'@'h' COMMENT 'c'`).(*nodes.CreateUserStmt)
	if create.ResourceGroup != "" {
		t.Errorf("resource_group = %q, want empty", create.ResourceGroup)
	}
}

// TestResourceControlLoc asserts the spans, which no corpus row can see. End
// is exclusive and stops at the statement's last token.
func TestResourceControlLoc(t *testing.T) {
	sql := `  CREATE RESOURCE GROUP rg RU_PER_SEC = 1000  `
	stmt := one(t, sql).(*nodes.CreateResourceGroupStmt)
	if stmt.Loc.Start != 2 || stmt.Loc.End != 44 {
		t.Errorf("stmt loc = %+v, want {2 44}", stmt.Loc)
	}
	if o := stmt.Options[0].Loc; o.Start != 27 || o.End != 44 {
		t.Errorf("option loc = %+v, want {27 44}", o)
	}

	cal := one(t, `CALIBRATE RESOURCE DURATION '20m'`).(*nodes.CalibrateResourceStmt)
	if cal.Loc.Start != 0 || cal.Loc.End != 33 {
		t.Errorf("calibrate loc = %+v, want {0 33}", cal.Loc)
	}
	watch := one(t, `QUERY WATCH REMOVE 7`).(*nodes.DropQueryWatchStmt)
	if watch.Loc.Start != 0 || watch.Loc.End != 20 {
		t.Errorf("watch loc = %+v, want {0 20}", watch.Loc)
	}
}
