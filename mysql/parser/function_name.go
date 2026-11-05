package parser

import "strings"

// keywordsNotFunctionNames are the keywords MySQL 8.0 refuses as an unquoted function name,
// in every argument shape: `SELECT comment(3)` is 1064 where `SELECT foo(3)` and
// the quoted `comment`(3) parse. They are no other word's business: all of them still
// name a column or a table where their category lets them, and a qualified `d.comment(3)`
// parses.
//
// The list is measured and not derived, because this parser's keyword table is not the
// engine's: it holds words the engine lexes as plain identifiers, which do take the call.
// TestKeywordFunctionNamesOracle holds it to mysql:8.0 over every keyword the lexer has.
//
// Ref: `function_call_generic: IDENT_sys '(' opt_udf_expr_list ')'` in MySQL 8.0 sql/sql_yacc.yy: IDENT_sys is no keyword.
var keywordsNotFunctionNames = wordSet(`
account action active admin after against aggregate algorithm always any array
assign_gtids_to_anonymous_transactions at attribute authentication auto_increment
autoextend_size avg_row_length backup begin binlog bit block bool boolean btree buckets
bulk byte cache cascaded catalog_name chain challenge_response changed channel checksum
cipher class_origin client clone close code column_format column_name columns comment
commit committed compact completion component compressed compression concurrent
connection consistent constraint_catalog constraint_name constraint_schema contains
context cpu current cursor_name data datafile datetime deallocate default_auth definer
definition delay_key_write description diagnostics directory disable discard disk do
dumpfile duplicate dynamic enable encryption end ends enforced engine engine_attribute
engines enum error errors escape event events every exchange exclude execute expansion
expire export extended extent_size factor failed_login_attempts fast faults fields file
file_block_size filter finish first fixed flush following follows found full general
generate geometry get_master_public_key get_source_public_key global grants
group_replication gtid_only handler hash help histogram history host hosts identified
ignore_server_ids import inactive indexes initial initial_size initiate insert_method
install instance invisible invoker io io_thread ipc isolation issuer json
key_block_size keyring language last leaves less level list local locked locks logfile
logs master master_auto_position master_compression_algorithms master_connect_retry
master_delay master_heartbeat_period master_host master_log_file master_log_pos
master_password master_port master_public_key_path master_retry_count master_ssl
master_ssl_ca master_ssl_capath master_ssl_cert master_ssl_cipher master_ssl_crl
master_ssl_crlpath master_ssl_key master_tls_ciphersuites master_tls_version master_user
master_zstd_compression_level max_connections_per_hour max_queries_per_hour max_rows
max_size max_updates_per_hour max_user_connections medium member memory merge
message_text migrate min_rows mode modify mutex mysql_errno name names national nchar
ndb ndbcluster nested network_namespace never new next no no_wait nodegroup none nowait
nulls number nvarchar offset oj old one only open optional options ordinality
organization others owner pack_keys page parser partial partitioning partitions password
password_lock_time path persist persist_only phase plugin plugin_dir plugins port
precedes preceding prepare preserve prev privilege_checks_user privileges process
processlist profile profiles proxy query quick random read_only rebuild recover
redo_buffer_size redundant reference registration relay relay_log_file relay_log_pos
relay_thread relaylog reload remove reorganize repair repeatable replica replicas
replicate_do_db replicate_do_table replicate_ignore_db replicate_ignore_table
replicate_rewrite_db replicate_wild_do_table replicate_wild_ignore_table replication
require_row_format require_table_primary_key_check reset resource respect restart
restore resume retain returned_sqlstate returning returns reuse role rollback rollup
rotate routine row_format rtree savepoint schedule schema_name secondary
secondary_engine secondary_engine_attribute secondary_load secondary_unload security
serial serializable server session share shutdown signed simple skip slave slow snapshot
socket some soname sounds source source_auto_position source_bind
source_compression_algorithms source_connect_retry source_connection_auto_failover
source_delay source_heartbeat_period source_host source_log_file source_log_pos
source_password source_port source_public_key_path source_retry_count source_ssl
source_ssl_ca source_ssl_capath source_ssl_cert source_ssl_cipher source_ssl_crl
source_ssl_crlpath source_ssl_key source_ssl_verify_server_cert source_tls_ciphersuites
source_tls_version source_user source_zstd_compression_level sql_after_gtids
sql_after_mts_gaps sql_before_gtids sql_thread srid stacked start starts
stats_auto_recalc stats_persistent stats_sample_pages status stop storage stream string
subclass_origin subject subpartition subpartitions super suspend swaps switches
table_checksum table_name tables tablespace temporary temptable text than
thread_priority ties tls transaction triggers type types unbounded uncommitted undefined
undo_buffer_size undofile unicode uninstall unknown unregister until upgrade url use_frm
user_resources validation value variables vcpu view visible wait warnings without work
wrapper x509 xa xid xml zone
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
