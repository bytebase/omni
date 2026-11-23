package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// mariadbNativeFunctions are the built-in functions MariaDB 11.8.8 resolves
// through its native function table (func_array in sql/item_create.cc): they
// resolve whether written plain, backtick-quoted, or under mariadb_schema /
// oracle_schema. Measured against the pinned mariadb:11.8.8 server.
var mariadbNativeFunctions = func() map[string]bool {
	m := map[string]bool{}
	for _, name := range strings.Fields(`
abs acos add_months addtime aes_decrypt aes_encrypt area asbinary asin astext aswkb aswkt
atan atan2 benchmark bin binlog_gtid_pos bit_count bit_length boundary buffer ceil ceiling
centroid char_length character_length chr coalesce coercibility collation column_check
column_exists column_json column_list compress concat concat_operator_oracle concat_ws
connection_id contains conv convert_tz convexhull cos cot crc32 crc32c crosses database
date_format datediff dayname dayofmonth dayofweek dayofyear decode decode_histogram
decode_oracle degrees des_decrypt des_encrypt dimension disjoint elt encode encrypt
endpoint envelope equals exp export_set exteriorring extractvalue field find_in_set floor
format format_bytes format_pico_time found_rows from_base64 from_days from_unixtime
geomcollfromtext geomcollfromwkb geometrycollectionfromtext geometrycollectionfromwkb
geometryfromtext geometryfromwkb geometryn geometrytype geomfromtext geomfromwkb get_lock
glength greatest hex ifnull inet6_aton inet6_ntoa inet_aton inet_ntoa instr interiorringn
intersects is_free_lock is_ipv4 is_ipv4_compat is_ipv4_mapped is_ipv6 is_used_lock
isclosed isempty isnull isring issimple json_array json_array_append json_array_insert
json_array_intersect json_compact json_contains json_contains_path json_depth
json_detailed json_equals json_exists json_extract json_insert json_key_value json_keys
json_length json_loose json_merge json_merge_patch json_merge_preserve json_normalize
json_object json_object_filter_keys json_object_to_array json_overlaps json_pretty
json_query json_quote json_remove json_replace json_schema_valid json_search json_set
json_type json_unquote json_valid json_value kdf last_day last_insert_id lcase least
length lengthb linefromtext linefromwkb linestringfromtext linestringfromwkb ln load_file
locate log log10 log2 lower lpad lpad_oracle ltrim ltrim_oracle make_set makedate maketime
master_gtid_wait master_pos_wait mbrcontains mbrdisjoint mbrequal mbrequals mbrintersects
mbroverlaps mbrtouches mbrwithin md5 microsecond mlinefromtext mlinefromwkb mod monthname
mpointfromtext mpointfromwkb mpolyfromtext mpolyfromwkb multilinestringfromtext
multilinestringfromwkb multipointfromtext multipointfromwkb multipolygonfromtext
multipolygonfromwkb name_const natural_sort_key nullif numgeometries numinteriorrings
numpoints nvl nvl2 oct octet_length old_password ord overlaps password period_add
period_diff pi pointfromtext pointfromwkb pointn pointonsurface polyfromtext polyfromwkb
polygonfromtext polygonfromwkb pow power quarter quote radians rand random_bytes
regexp_instr regexp_replace regexp_substr release_all_locks release_lock replace_oracle
reverse round row_count rpad rpad_oracle rtrim rtrim_oracle schema schemas sec_to_time
sformat sha sha1 sha2 sign sin sleep soundex space sqrt srid st_area st_asbinary
st_asgeojson st_astext st_aswkb st_aswkt st_boundary st_buffer st_centroid st_contains
st_convexhull st_crosses st_difference st_dimension st_disjoint st_distance
st_distance_sphere st_endpoint st_envelope st_equals st_exteriorring st_geomcollfromtext
st_geomcollfromwkb st_geometrycollectionfromtext st_geometrycollectionfromwkb
st_geometryfromtext st_geometryfromwkb st_geometryn st_geometrytype st_geomfromgeojson
st_geomfromtext st_geomfromwkb st_interiorringn st_intersection st_intersects st_isclosed
st_isempty st_isring st_issimple st_length st_linefromtext st_linefromwkb
st_linestringfromtext st_linestringfromwkb st_mlinefromtext st_mlinefromwkb
st_mpointfromtext st_mpointfromwkb st_mpolyfromtext st_mpolyfromwkb
st_multilinestringfromtext st_multilinestringfromwkb st_multipointfromtext
st_multipointfromwkb st_multipolygonfromtext st_multipolygonfromwkb st_numgeometries
st_numinteriorrings st_numpoints st_overlaps st_pointfromtext st_pointfromwkb st_pointn
st_pointonsurface st_polyfromtext st_polyfromwkb st_polygonfromtext st_polygonfromwkb
st_relate st_srid st_startpoint st_symdifference st_touches st_union st_within st_x st_y
startpoint str_to_date strcmp substr_oracle substring_index subtime sys_guid tan
time_format time_to_sec timediff to_base64 to_char to_days to_seconds touches ucase
uncompress uncompressed_length unhex unix_timestamp updatexml upper uuid uuid_short
uuid_v4 uuid_v7 vec_distance vec_distance_cosine vec_distance_euclidean vec_fromtext
vec_totext version week weekday weekofyear within wsrep_last_seen_gtid
wsrep_last_written_gtid wsrep_sync_wait_upto_gtid x y yearweek
`) {
		m[name] = true
	}
	return m
}()

// mariadbGrammarBuiltinFunctions are MariaDB 11.8.8 built-in functions produced
// by dedicated grammar productions in sql/sql_yacc.yy (keyword, nonkeyword,
// conflict, geometry, dynamic-column, sequence, aggregate, and window functions)
// plus parser-synthesized directional TRIM forms. Unlike native functions, a
// backtick-quoted call to one of these is treated as a stored-function call.
var mariadbGrammarBuiltinFunctions = func() map[string]bool {
	m := map[string]bool{}
	for _, name := range strings.Fields(`
adddate ascii avg binary bit_and bit_or bit_xor cast char character charset column_add
column_create column_delete column_get convert count cume_dist curdate current_date
current_role current_time current_timestamp current_user curtime date date_add date_sub
day default dense_rank extract first_value geomcollection geometrycollection get_format
group_concat hour if insert interval json_arrayagg json_objectagg lag last_value lastval
lead left linestring localtime localtimestamp match max median mid min minute month
multilinestring multipoint multipolygon nextval now nth_value ntile percent_rank
percentile_cont percentile_disc point polygon position rank repeat replace right row
row_number rownum second session_user setval std stddev stddev_pop stddev_samp subdate
substr substring sum sysdate system_user time timestamp timestampadd timestampdiff trim
trim_both trim_from trim_leading trim_oracle trim_trailing truncate user utc_date utc_time
utc_timestamp value values var_pop var_samp variance weight_string year
`) {
		m[name] = true
	}
	return m
}()

// mariadbFuncCollisionNames are the 748 words in MariaDB 11.8.8's lexer symbol
// tables (information_schema.KEYWORDS / symbols[] + sql_functions[]): when a
// 1-part or 2-part function call fails to resolve as a built-in or stored
// function and its name is in this set, MariaDB reports 1630
// (ER_FUNC_INEXISTENT_NAME_COLLISION) rather than 1305 (ER_SP_DOES_NOT_EXIST).
// Measured against the pinned mariadb:11.8.8 server.
var mariadbFuncCollisionNames = func() map[string]bool {
	m := map[string]bool{}
	for _, name := range strings.Fields(`
!= && << <= <=> <> >= >> accessible account action add adddate admin after against
aggregate algorithm all alter always analyze and any as asc ascii asensitive at atomic
authors auto auto_increment autoextend_size avg avg_row_length backup before begin between
bigint binary binlog bit bit_and bit_or bit_xor blob block body bool boolean both btree by
byte cache call cascade cascaded case cast catalog_name chain change changed channel char
character charset check checkpoint checksum cipher class_origin client clob close coalesce
code collate collation column column_add column_check column_create column_delete
column_get column_name columns comment commit committed compact completion compressed
concurrent condition connection consistent constraint constraint_catalog constraint_name
constraint_schema contains context continue contributors convert count cpu create cross
cube cume_dist curdate current current_date current_pos current_role current_time
current_timestamp current_user cursor cursor_name curtime cycle data database databases
datafile date date_add date_sub datetime day day_hour day_microsecond day_minute
day_second deallocate dec decimal declare default definer delay_key_write delayed delete
delete_domain_id dense_rank des_key_file desc describe deterministic diagnostics directory
disable discard disk distinct distinctrow div do do_domain_ids double drop dual dumpfile
duplicate dynamic each else elseif elsif empty enable enclosed end ends engine engines
enum error errors escape escaped event events every examined except exception exchange
exclude execute exists exit expansion expire explain export extended extent_size extract
false fast faults federated fetch fields file first first_value fixed float float4 float8
flush following follows for force foreign format found from full fulltext function general
generated get get_format global goto grant grants group group_concat handler hard hash
having help high_priority history host hosts hour hour_microsecond hour_minute hour_second
id identified if ignore ignore_domain_ids ignore_server_ids ignored immediate import in
increment index indexes infile initial_size inner inout insensitive insert insert_method
install int int1 int2 int3 int4 int8 integer intersect interval into invisible invoker io
io_thread ipc is isolation isopen issuer iterate join json json_arrayagg json_objectagg
json_table key key_block_size keys kill lag language last last_value lastval lead leading
leave leaves left less level like limit linear lines list load local localtime
localtimestamp lock locked locks logfile logs long longblob longtext loop low_priority
master master_connect_retry master_delay master_demote_to_replica master_demote_to_slave
master_gtid_pos master_heartbeat_period master_host master_log_file master_log_pos
master_password master_port master_server_id master_ssl master_ssl_ca master_ssl_capath
master_ssl_cert master_ssl_cipher master_ssl_crl master_ssl_crlpath master_ssl_key
master_ssl_verify_server_cert master_use_gtid master_user match max
max_connections_per_hour max_queries_per_hour max_rows max_size max_statement_time
max_updates_per_hour max_user_connections maxvalue median medium mediumblob mediumint
mediumtext memory merge message_text microsecond mid middleint migrate min min_rows minus
minute minute_microsecond minute_second minvalue mod mode modifies modify monitor month
mutex mysql mysql_errno name names national natural nchar nested never next nextval no
no_wait no_write_to_binlog nocache nocopy nocycle nodegroup nomaxvalue nominvalue none not
notfound now nowait nth_value ntile null number numeric nvarchar of offset old_password on
one online only open optimize option optionally options or order ordinality others out
outer outfile over overlaps owner pack_keys package page page_checksum parse_vcol_expr
parser partial partition partitioning partitions password path percent_rank
percentile_cont percentile_disc period persistent phase plugin plugins port portion
position precedes preceding precision prepare preserve prev previous primary privileges
procedure process processlist profile profiles proxy purge quarter query quick raise range
rank raw read read_only read_write reads real rebuild record recover recursive
redo_buffer_size redofile redundant ref_system_id references regexp relay relay_log_file
relay_log_pos relay_thread relaylog release reload remove rename reorganize repair repeat
repeatable replace replay replica replica_pos replicas replication require reset resignal
restart restore restrict resume return returned_sqlstate returning returns reuse reverse
revoke right rlike role rollback rollup routine row row_count row_format row_number
rowcount rownum rows rowtype rtree savepoint schedule schema schema_name schemas second
second_microsecond security select sensitive separator sequence serial serializable server
session session_user set setval share show shutdown signal signed simple skip slave
slave_pos slaves slow smallint snapshot socket soft some soname sounds source spatial
specific sql sql_after_gtids sql_before_gtids sql_big_result sql_buffer_result sql_cache
sql_calc_found_rows sql_no_cache sql_small_result sql_thread sql_tsi_day sql_tsi_hour
sql_tsi_minute sql_tsi_month sql_tsi_quarter sql_tsi_second sql_tsi_week sql_tsi_year
sqlexception sqlstate sqlwarning ssl stage start starting starts statement
stats_auto_recalc stats_persistent stats_sample_pages status std stddev stddev_pop
stddev_samp stop storage stored straight_join string subclass_origin subdate subject
subpartition subpartitions substr substring sum super suspend swaps switches sysdate
system system_time system_user table table_checksum table_name tables tablespace temporary
temptable terminated text than then threads ties time timestamp timestampadd timestampdiff
tinyblob tinyint tinytext to trailing transaction transactional trigger triggers trim
trim_oracle true truncate type unbounded uncommitted undefined undo undo_buffer_size
undofile unicode uninstall union unique unknown unlock unsigned until update upgrade usage
use use_frm user user_resources using utc_date utc_time utc_timestamp validation value
values var_pop var_samp varbinary varchar varchar2 varcharacter variables variance varying
vector versioning via view virtual visible wait warnings week weight_string when where
while window with within without work wrapper write x509 xa xml xor year year_month
zerofill ||
`) {
		m[name] = true
	}
	return m
}()

// isMariaDBBuiltinFunction reports whether lowerName resolves as a MariaDB
// 11.8.8 built-in function. When quoted is true (backtick-quoted function
// name), only native func_array functions resolve as built-ins; grammar
// productions fall through to stored-function lookup.
func isMariaDBBuiltinFunction(lowerName string, quoted bool) bool {
	if mariadbNativeFunctions[lowerName] {
		return true
	}
	if quoted {
		return false
	}
	return mariadbGrammarBuiltinFunctions[lowerName]
}

// isMariaDBSchemaBuiltinFunction reports whether a schema-qualified call under
// mariadb_schema or oracle_schema resolves to a built-in function in MariaDB
// 11.8.8. Native functions resolve with or without backticks; the schema-aware
// grammar forms (REPLACE, SUBSTR/SUBSTRING, TRIM) resolve when unquoted and
// with a valid built-in argument shape.
func isMariaDBSchemaBuiltinFunction(fc *nodes.FuncCallExpr, lowerName string) bool {
	if mariadbNativeFunctions[lowerName] {
		return true
	}
	if fc == nil || fc.NameQuoted {
		return false
	}
	switch lowerName {
	case "trim_leading", "trim_trailing", "trim_both":
		return true
	case "trim":
		return fc.SpecialSyntax || len(fc.Args) == 1 || len(fc.Args) == 2
	case "substring", "substr":
		return fc.SpecialSyntax || len(fc.Args) == 2 || len(fc.Args) == 3
	case "replace":
		return len(fc.Args) == 3
	default:
		return false
	}
}

// isMariaDBFuncCollisionName reports whether lowerName collides with a MariaDB
// 11.8.8 lexer keyword or function symbol (causing 1630 instead of 1305 on a
// missing 1-part or 2-part function call).
func isMariaDBFuncCollisionName(lowerName string) bool {
	return mariadbFuncCollisionNames[lowerName]
}
