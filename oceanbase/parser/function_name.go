package parser

import "strings"

// keywordsNotFunctionNames are the keywords OceanBase 4.3.5.6 refuses as an unquoted function name,
// in every argument shape: `SELECT comment(3)` is 1064 where `SELECT foo(3)` and
// the quoted `comment`(3) parse. They are no other word's business: all of them still
// name a column or a table where their category lets them, and a qualified `d.comment(3)`
// parses.
//
// The list is measured and not derived, because this parser's keyword table is not the
// engine's: it holds words the engine lexes as plain identifiers, which do take the call.
// TestKeywordFunctionNamesOceanBaseOracle holds it to oceanbase/oceanbase-ce:4.3.5-lts over every keyword the lexer has.
//
// Ref: `function_name` in sql_parser_mysql_mode.y (v4.3.5_CE_BP6), which is NAME_OB,
// STATISTICS and DUMP; any other keyword call has a rule of its own in func_expr.
var keywordsNotFunctionNames = wordSet(`
account action active admin after against aggregate algorithm always any at attribute
auto auto_increment auto_increment_mode auto_refresh autoextend_size avg_row_length
backup begin binlog bit block block_size bool boolean btree buckets byte cache cascaded
catalog_name chain changed checksum cipher class_origin client clone close code
column_format column_name columns comment commit committed compact completion compressed
compression concurrent connection consistent constraint_catalog constraint_name
constraint_schema context cpu current cursor_name data datafile datetime deallocate
default_auth definer delay_key_write description diagnostics directory disable discard
disk do dumpfile duplicate dynamic dynamic_partition_policy enable
enable_macro_block_bloom_filter encryption end ends enforced engine engines enum error
errors escape event events every exchange execute expansion expire export extended
extent_size external fast faults fields file file_block_size files filter first fixed
flush following found full general geometry global grants handler hash help histogram
host hosts identified ignore_server_ids import indexes initial_size insert_method
install instance invisible invoker io io_thread ipc isolation issuer json
key_block_size language last leaves less level list local locked locks logfile logs
manual master master_auto_position master_connect_retry master_delay
master_heartbeat_period master_host master_log_file master_log_pos master_password
master_port master_retry_count master_ssl master_ssl_ca master_ssl_capath
master_ssl_cert master_ssl_cipher master_ssl_crl master_ssl_crlpath master_ssl_key
master_user materialized max_connections_per_hour max_queries_per_hour max_rows max_size
max_updates_per_hour max_user_connections medium member memory merge message_text
migrate min_rows mode modify mutex mysql_errno name names national nchar ndb ndbcluster
nested never new next no no_wait nodegroup none nowait nulls number nvarchar offset oj
old one only open options ordinality organization overwrite owner pack_keys page
parallel parameters parser partial partitioning partitions path pctfree phase plugin
plugin_dir plugins port preceding prepare preserve prev privileges process processlist
profile profiles proxy query quick read_only rebuild recover redo_buffer_size redundant
relay relay_log_file relay_log_pos relay_thread relaylog reload remove reorganize repair
repeatable replica replica_num replication reset resource respect restart restore resume
returned_sqlstate returning returns role rollback rollup rotate routine row_format rtree
savepoint schedule schema_name security serial serializable server session sets share
shutdown signed simple skip slave slow snapshot socket some soname sounds source
sql_after_gtids sql_after_mts_gaps sql_before_gtids sql_thread sql_tsi_day sql_tsi_hour
sql_tsi_minute sql_tsi_month sql_tsi_quarter sql_tsi_second sql_tsi_week sql_tsi_year
srid stacked start starts stats_auto_recalc stats_persistent stats_sample_pages status
stop storage string subclass_origin subject subpartition subpartitions super suspend
swaps switches table_checksum table_name tables tablespace tablet_size template
temporary temptable text than transaction triggers type types unbounded uncommitted
undefined undo_buffer_size undofile unicode uninstall unknown until upgrade
use_bloom_filter use_frm user_resources validate validation value variables vector view
visible wait warnings without work wrapper x509 xa xid xml zone
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
	return p.cur.Type != tokIDENT && p.peekNext().Type == '(' && keywordsNotFunctionNames[strings.ToLower(p.cur.Str)]
}
