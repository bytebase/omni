package parser

import "strings"

// keywordsNotFunctionNames are the keywords TiDB v8.5.5 refuses as an unquoted function name,
// in every argument shape: `SELECT comment(3)` is 1064 where `SELECT foo(3)` and
// the quoted `comment`(3) parse. They are no other word's business: all of them still
// name a column or a table where their category lets them, and a qualified `d.comment(3)`
// parses.
//
// The list is measured and not derived, because this parser's keyword table is not the
// engine's: it holds words the engine lexes as plain identifiers, which do take the call.
// TestKeywordFunctionNamesOracle holds it to pingcap/tidb:v8.5.5 over every keyword the lexer has.
//
// Ref: `FunctionCallGeneric: identifier '(' ExpressionListOpt ')'` in TiDB v8.5.5 pkg/parser/parser.y: `identifier` is the plain token, where `Identifier` takes the keywords too.
var keywordsNotFunctionNames = wordSet(`
account action admin after against algorithm always any array attribute auto_id_cache
auto_increment auto_random auto_random_base avg_row_length background backup backups
batch begin bernoulli binding binding_cache bindings binlog bit block bool boolean btree
buckets builtins burstable byte cache calibrate cascaded chain checksum cipher client
close clustered column_format column_stats_usage columns comment commit committed
compact compressed compression config connection consistent constraints context cooldown
cpu current cycle data datetime deallocate declare defined definer delay_key_write digest
directory disable disabled discard disk distribute distribution distributions do dry
dryrun duplicate duration dynamic enable enabled encryption end end_time enforced engine
engines enum error errors escape event events exact exchange exec_elapsed execute
expansion expire extended failed_login_attempts faults fields file first fixed flashback
flush follower_constraints followers following found full function general global grants
handler hash help high histogram histograms_in_flight history hosts identified import
increment indexes insert_method instance invisible invoker io ipc isolation issuer job
jobs json key_block_size labels language last leader_constraints learner_constraints
learners less level list local location locked logs low master max_connections_per_hour
max_queries_per_hour max_rows max_updates_per_hour max_user_connections medium member
memory merge min_rows minvalue mode modify names national nchar never next next_row_id
no nocache nocycle nodegroup nomaxvalue nominvalue nonclustered none nowait nulls
nvarchar off offset oltp_read_only oltp_read_write oltp_write_only only open optional
pack_keys page parser partial partitioning partitions password_lock_time placement plan
plugins policy pre_split_regions preceding prepare preserve primary_region priority
privileges process processed_keys processlist profile profiles proxy purge query
query_limit quick rebuild recover redundant region regions reload remove reorganize
repair repeatable replica replicas replication reset resource respect restart restore
restores resume reuse role rollback rollup routine row_format rtree ru ru_per_sec rule
run s3 savepoint schedule secondary secondary_engine secondary_load secondary_unload
security separator sequence serial serializable session session_states shard_row_id_bits
share shutdown signed similar simple skip slave slow snapshot some source split
sql_tsi_day sql_tsi_hour sql_tsi_minute sql_tsi_month sql_tsi_quarter sql_tsi_second
sql_tsi_week sql_tsi_year start start_time stats_auto_recalc stats_buckets
stats_col_choice stats_col_list stats_extended stats_healthy stats_histograms
stats_locked stats_meta stats_persistent stats_sample_pages stats_sample_rate stats_topn
status stop storage subject subpartition subpartitions super survival_preferences swaps
switch_group switches system table_checksum tables tablesample tablespace task_types
temporary temptable text than tiflash timeout tls tpcc tpch_10 transaction triggers ttl
ttl_enable ttl_job_interval type unbounded uncommitted undefined unicode unknown
unlimited until utilization_limit validation value variables vector view visible
voter_constraints voters wait warnings watch without workload x509
`)

// tidbKeywordTokensNotFunctionNames are words TiDB's lexer turns into a keyword token
// (they are all in tidbLexerWords) that have no function-call production, yet this
// parser's forked lexer scans as plain identifiers. refusesFunctionName's keyword branch
// below keys on this parser's keyword tokens, so it never sees them, and the generic
// `identifier '(' ... ')'` path would take them: `DO queries()` and `DO last_backup(x)`
// are 1064s on TiDB v8.5.5 where `DO foo()` parses. Measured against v8.5.5; none has an
// `'(' ... ')'` arm in pkg/parser/parser.y (FunctionNameConflict, FunctionCallKeyword,
// FunctionCallNonKeyword, SumExpr, FunctionNameSequence), so refusing the call can never
// reject valid SQL. A quoted name is an identifier whatever it spells, so it is excluded.
var tidbKeywordTokensNotFunctionNames = wordSet(`
defined last_backup queries replayer tso
`)

func wordSet(words string) map[string]bool {
	set := make(map[string]bool)
	for _, word := range strings.Fields(words) {
		set[word] = true
	}
	return set
}

// refusesFunctionName reports whether the current token is a keyword the engine does not
// take as the name of the call that follows it. A quoted name lexes as an identifier
// whatever it spells, so it never comes here.
func (p *Parser) refusesFunctionName() bool {
	if p.peekNext().Type != '(' {
		return false
	}
	word := strings.ToLower(p.cur.Str)
	if p.cur.Type != tokIDENT {
		return keywordsNotFunctionNames[word]
	}
	// The word reaches here as a plain identifier in this parser's lexer, but TiDB
	// may lex it as a keyword token with no function-call production. A quoted name
	// is always an identifier, so it still takes the generic call.
	return !p.backtickQuotedAt(p.cur.Loc) && tidbKeywordTokensNotFunctionNames[word]
}
