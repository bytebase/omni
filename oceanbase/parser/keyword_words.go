package parser

import "strings"

// oceanBaseKeywords are the words OceanBase's lexer makes a keyword token of: the 1,143
// of Mysql_none_reserved_keywords (src/sql/parser/non_reserved_keywords_mysql_mode.c,
// oceanbase v4.3.5_CE_BP6), which {identifier} is looked up in before it becomes a
// NAME_OB. The table holds reserved and unreserved words alike, and far more than this
// parser's own keyword table: PRIMARY_ZONE is a keyword there and a plain identifier
// here.
//
// It matters where the grammar writes NAME_OB with no unreserved_keyword beside it. A
// keyword is not a name there, reserved or not, until it is quoted.
var oceanBaseKeywords = wordSet(`
_st_asmvt access access_info accessible accessid accesskey accesstype account action activate
active add adddate admin after against aggregate algorithm all all_meta all_user allow alter
always analyse analyze and any api_mode approx approx_count_distinct
approx_count_distinct_synopsis approx_count_distinct_synopsis_merge approximate arbitration
archivelog array array_agg array_filter array_first array_map array_sortby as asc ascii asensitive
asis asynchronous at attribute audit authors auto auto_increment auto_increment_cache_size
auto_increment_mode auto_refresh autoextend_size availability avg avg_row_length backed backup
backup_backup_dest backup_copies backuppiece backupround backupset badfile balance bandwidth base
baseline baseline_id basic before begin between bigint bigint_precision binary binary_format
binding binlog bit bit_and bit_or bit_xor blob block block_size bool boolean bootstrap both
boundary_column boundary_column_unit breadth btree buckets buffer_size by byte cache
calc_partition_id calibration calibration_info call cancel cascade cascaded case cast catalog
catalog_name catalogs chain change changed char character charset check checkpoint checksum chunk
cipher class_origin clean clear client clone close cluster cluster_id cluster_name coalesce code
collate collation collect_statistics_on_create column column_bloom_filter column_format
column_index_type column_name columns comment commit committed compact complete completion
compressed compression compression_block_size compression_code computation compute concurrent
concurrent_limiting_rule condensed condition conditional connect connection consistent constraint
constraint_catalog constraint_name constraint_schema contains context continue contributors
convert copy cosine count cpu create create_timestamp cross ctx_id ctxcat cube cume_dist curdate
current current_date current_time current_timestamp current_user cursor cursor_name curtime cycle
dag data data_disk_size data_source data_table_id database database_id databases datafile date
date_add date_sub datetime day day_hour day_microsecond day_minute day_second deallocate dec
decimal declare decrypt decryption default default_auth default_lob_inrow_threshold
default_tablegroup definer delay delay_key_write delayed delete delete_insert demand dense_rank
depth des_key_file desc describe description destination deterministic diagnostics dict_table
directory disable disallow discard disconnect disk diskgroup distinct distinctrow div do dot
double drop dual dump dumpfile duplicate duplicate_read_consistency duplicate_scope dynamic
dynamic_partition_policy each effective else elseif empty empty_field_as_null enable
enable_arbitration_service enable_extended_rowid enable_macro_block_bloom_filter enclosed encoding
encrypt encrypted encryption end endpoint ends enforced engine engines entity enum error
error_code errors escape escaped estimate euclidean event events every except exchange excluding
exclusive execute exists exit expansion expire expire_info expire_time expired explain export
extended extended_noaddr extent_size external external_storage_dest extract failover fast faults
fetch field_delimiter field_enclosed_by field_optionally_enclosed_by fields file file_block_size
file_extension file_id files filter final_count first first_value fixed flashback float float4
float8 flush follower following for force foreign format found fragmentation freeze frequency from
frozen full fulltext function general generated geomcollection geometry geometrycollection get
get_format global global_name grant grants group group_concat grouping grouping_id gts handler
hash having heap help hidden high_priority histogram host hosts hot_retention hour
hour_microsecond hour_minute hour_second hybrid hybrid_hist id idc identified if ignore
ignore_last_empty_column ignore_server_ids ik_mode ilogcache immediate import in including
inconsistent incr increment incremental index index_table_id indexed indexes individual infile
info initial_size inner innodb inout inplace insensitive insert insert_method install instance
instant int int1 int2 int3 int4 int8 integer intersect into invisible invoker io io_after_gtids
io_before_gtids io_thread iops_weight ipc is isnull isolate isolation issuer iterate job join json
json_arrayagg json_objectagg json_query json_table json_value key key_block_size key_version keys
keyword kill kv_attributes kvcache lag language last last_refresh_scn last_value lateral lead
leader leading leak leak_mod leak_rate leave leaves left less level lib license like limit
line_delimiter linear lines linestring link list listagg ln load lob_inrow_threshold local
locality localtime localtimestamp location lock locked locks log log_disk_size log_restore_source
logfile logical_reads logonly_replica_num logs long longblob longtext loop low_priority ls major
manhattan manual map master master_auto_position master_bind master_connect_retry master_delay
master_heartbeat_period master_host master_log_file master_log_pos master_password master_port
master_retry_count master_server_id master_ssl master_ssl_ca master_ssl_capath master_ssl_cert
master_ssl_cipher master_ssl_crl master_ssl_crlpath master_ssl_key master_ssl_verify_server_cert
master_user match materialized max max_concurrency max_connections_per_hour max_cpu max_file_size
max_iops max_net_bandwidth max_ngram_size max_queries_per_hour max_rows max_size max_token_size
max_updates_per_hour max_used_part_id max_user_connections maximize maxvalue medium mediumblob
mediumint mediumtext member memory memory_size memstore_percent memtable merge merge_engine
message_text meta micro_index_clustered microsecond middleint migrate migration min min_cpu
min_iops min_max min_ngram_size min_rows min_token_size minor minus minute minute_microsecond
minute_second minvalue mismatch mod mode modifies modify module month move multilinestring
multipoint multipolygon multivalue mutex my_name mysql_errno name names namespace national natural
nchar ndb ndbcluster nested net_bandwidth_weight network never new next ngram_token_size no
no_parallel no_rewrite no_wait no_write_to_binlog noarchivelog noaudit nocache nocycle nodegroup
nomaxvalue nominvalue none noorder noparallel normal now nowait nth_value ntile null_if nulls
number numeric nvarchar obconfig_url object object_id object_name_display obsolete occur of off
offset oj old old_key old_password on on_error one one_shot only open optimize optimizer_costs
option optionally options or order ordinality organization orig_default out outer outfile outline
over overwrite owner p_chunk p_entity pack_keys page parallel parameters parse_header parser
parser_properties partial partial_update partition partition_id partition_type partitioning
partitions password path pattern pause paxos_replica_num pctfree per percent_rank percentage
percentile_cont performance phase physical pl plainaccess plan planregress plugin plugin_dir
plugins plus point policy polygon pool port position preceding precision precreate_time prepare
preserve pretty pretty_color prev preview primary primary_zone priority privileges procedure
process processlist proctime profile profiles progressive_merge_num project_name properties
protection proxy ps public purge quantifier_table quarter query query_response_time queue_time
quick quota_name range rank rb_and_agg rb_and_cardinality_agg rb_build_agg rb_iterate rb_or_agg
rb_or_cardinality_agg read read_consistency read_only read_write reads real rebuild recover
recovery recovery_window recursive recycle recyclebin redo_buffer_size redo_transport_options
redofile redundancy redundant references refresh regexp region reject relay relay_log_file
relay_log_pos relay_thread relaylog release reload remap remote_oss remove rename reorganize
repair repeat repeatable replace replica replica_num replica_type replication report require reset
resignal resource resource_pool resource_pool_list respect restart restore restrict resume return
returned_sqlstate returning returns reverse revoke rewrite right rlike roaringbitmap role rollback
rolling rollup root rootservice rootservice_list roottable rotate routine row row_count row_format
row_group_size row_index_stride row_number rowid rows rt rtree rule rules run sample savepoint
scalars schedule schema schema_id schema_name schemas scn scope second second_microsecond security
seed select semistruct_encoding_type sensitive separator sequence sequences serial serializable
server server_ip server_port server_type service session session_user set set_master_cluster
set_slave_cluster set_tp sets sharding share shared shared_storage_dest shared_storage_info show
shutdown signal signed simple single size skewonly skip skip_blank_lines skip_header skip_index
slave slog slot_idx slow smallint snapshot socket some soname sounds source sparsevector spatial
specific spfile split sql sql_after_gtids sql_after_mts_gaps sql_before_gtids sql_big_result
sql_buffer_result sql_cache sql_calc_found_rows sql_id sql_no_cache sql_small_result sql_thread
sql_tsi_day sql_tsi_hour sql_tsi_minute sql_tsi_month sql_tsi_quarter sql_tsi_second sql_tsi_week
sql_tsi_year sqlexception sqlstate sqlwarning srid ss_micro_cache ssl stacked standby start
starting starts stat statement statement_id statements statistics stats_auto_recalc
stats_persistent stats_sample_pages status std stddev stddev_pop stddev_samp stop stopword_table
storage storage_cache_policy storage_cache_policy_executor storage_format_version store stored
storing straight_join string stripe_size strong ststoken subclass_origin subdate subject
subpartition subpartitions substr substring successful sum sum_opnsize super suspend swaps switch
switches switchover synchronization synchronous sysdate system system_user table table_checksum
table_id table_mode table_name tablegroup tablegroup_id tablegroups tables tablespace tablet
tablet_id tablet_max_size tablet_size tablets task task_id template temporary temptable tenant
tenant_id tenant_sts_credential terminated text than then throttle time time_unit time_zone
time_zone_info times timestamp timestampadd timestampdiff tinyblob tinyint tinytext to
top_k_fre_hist tp_name tp_no trace traditional trailing transaction transfer trigger triggers trim
trim_space truncate ttl tunnel_endpoint type types unbounded uncommitted unconditional undefined
undo undo_buffer_size undofile unicode uninstall union unique unit unit_group unit_num unknown
unlimited unlock unlocked unnest unsigned until unusual up update upgrade usage use
use_bloom_filter use_frm user user_resources user_specified using utc_date utc_time utc_timestamp
valid validate validation value values var_pop var_samp varbinary varchar varcharacter variables
variance varying vector vector_distance vector_similarity verbose verify view virtual
virtual_column_id visible wait warnings wash weak week weight_string when whenever where while
window with without work wrapper write x509 xa xid xml xor year year_month zerofill zone zone_list
zone_type
`)

// parseNameOB reads NAME_OB, an identifier that is no keyword of OceanBase's: the name
// of an index in an index hint is one (index_element), where `USE KEY (primary_zone)` is
// a 1064 and USE KEY (`primary_zone`) parses.
func (p *Parser) parseNameOB() (string, int, error) {
	if (p.cur.Type == tokIDENT || p.cur.Type >= 700) && !p.backtickQuotedAt(p.cur.Loc) &&
		oceanBaseKeywords[strings.ToLower(p.cur.Str)] {
		return "", 0, p.syntaxErrorAtCur()
	}
	return p.parseIdent()
}
