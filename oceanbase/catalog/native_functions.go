package catalog

import "strings"

// nativeFunctions are the MySQL-compatible native functions registered in
// OceanBase CE 4.3.5.6 MySQL mode (REG_OP in ob_expr_operator_factory.cpp).
var nativeFunctions = func() map[string]bool {
	m := map[string]bool{}
	for _, name := range strings.Fields(`
abs acos addtime aes_decrypt aes_encrypt any_value area asin atan atan2 benchmark bin
bin_to_uuid bit_count bit_length can_access_column can_access_database can_access_event
can_access_resource_group can_access_routine can_access_table can_access_trigger
can_access_user can_access_view ceil ceiling centroid char_length character_length
coercibility compress concat concat_ws connection_id conv convert_tz cos cot crc32
current_role date_format datediff dayname dayofmonth dayofweek dayofyear degrees elt exp
export_set extractvalue field find_in_set floor format_bytes format_pico_time found_rows
from_base64 from_days from_unixtime get_lock greatest gtid_subset gtid_subtract hex
icu_version ifnull inet6_aton inet6_ntoa inet_aton inet_ntoa instr is_free_lock is_ipv4
is_ipv4_compat is_ipv4_mapped is_ipv6 is_used_lock is_uuid isnull json_append json_array
json_array_append json_array_insert json_contains json_contains_path json_depth
json_extract json_insert json_keys json_length json_merge json_merge_patch
json_merge_preserve json_object json_overlaps json_pretty json_quote json_remove
json_replace json_schema_valid json_schema_validation_report json_search json_set
json_storage_free json_storage_size json_type json_unquote json_valid last_day
last_insert_id lcase least length ln load_file locate log log10 log2 lower lpad ltrim
make_set makedate maketime master_pos_wait mbrcontains mbrcoveredby mbrcovers
mbrdisjoint mbrequals mbrintersects mbroverlaps mbrtouches mbrwithin md5 monthname
name_const nullif oct octet_length ord period_add period_diff pi pow power
ps_current_thread_id ps_thread_id quote radians rand random_bytes regexp_instr
regexp_like regexp_replace regexp_substr release_all_locks release_lock reverse
roles_graphml round rpad rtrim sec_to_time sha sha1 sha2 sign sin sleep soundex
source_pos_wait space sqrt st_area st_asbinary st_asgeojson st_astext st_aswkb st_aswkt
st_buffer st_buffer_strategy st_centroid st_contains st_convexhull st_crosses
st_difference st_dimension st_disjoint st_distance st_distance_sphere st_endpoint
st_envelope st_equals st_exteriorring st_frechetdistance st_geohash
st_geomcollfromtext st_geomcollfromtxt st_geomcollfromwkb st_geometrycollectionfromtext
st_geometrycollectionfromwkb st_geometryfromtext st_geometryfromwkb st_geometryn
st_geometrytype st_geomfromgeojson st_geomfromtext st_geomfromwkb st_hausdorffdistance
st_interiorringn st_intersection st_intersects st_isclosed st_isempty st_issimple
st_isvalid st_latfromgeohash st_latitude st_length st_linefromtext st_linefromwkb
st_lineinterpolatepoint st_lineinterpolatepoints st_linestringfromtext
st_linestringfromwkb st_longfromgeohash st_longitude st_makeenvelope st_mlinefromtext
st_mlinefromwkb st_mpointfromtext st_mpointfromwkb st_mpolyfromtext st_mpolyfromwkb
st_multilinestringfromtext st_multilinestringfromwkb st_multipointfromtext
st_multipointfromwkb st_multipolygonfromtext st_multipolygonfromwkb st_numgeometries
st_numinteriorring st_numinteriorrings st_numpoints st_overlaps st_pointatdistance
st_pointfromgeohash st_pointfromtext st_pointfromwkb st_pointn st_polyfromtext
st_polyfromwkb st_polygonfromtext st_polygonfromwkb st_simplify st_srid st_startpoint
st_swapxy st_symdifference st_touches st_transform st_union st_validate st_within st_x
st_y statement_digest statement_digest_text str_to_date strcmp substring_index subtime
tan time_format time_to_sec timediff to_base64 to_days to_seconds ucase uncompress
uncompressed_length unhex unix_timestamp updatexml upper uuid uuid_short uuid_to_bin
validate_password_strength version wait_for_executed_gtid_set
wait_until_sql_thread_after_gtids weekday weekofyear yearweek
`) {
		m[name] = true
	}
	return m
}()

// grammarBuiltinFunctions are built-in functions produced by dedicated grammar
// productions in OceanBase CE 4.3.5.6 MySQL mode (sql_parser_mysql_mode.y),
// plus parser-synthesized trim_* directional forms.
var grammarBuiltinFunctions = func() map[string]bool {
	m := map[string]bool{}
	for _, name := range strings.Fields(`
adddate ascii avg binary bit_and bit_or bit_xor cast char character charset coalesce
collation convert count cume_dist curdate current_date current_time current_timestamp
current_user curtime database date date_add date_sub day default dense_rank dump encrypt
extract first_value format geomcollection geometrycollection get_format group_concat
grouping grouping_id hour if insert interval json_arrayagg json_objectagg json_query
json_table json_value lag last_value lead left linestring localtime localtimestamp match
max microsecond mid min minute mod month multilinestring multipoint multipolygon normal
now nth_value ntile old_password password percent_rank point polygon position quarter
rank repeat replace right row_count row_number schema second session_user st_collect std
stddev stddev_pop stddev_samp subdate substr substring sum sysdate system_user time
timestamp timestampadd timestampdiff trim trim_both trim_from trim_leading trim_trailing
truncate user utc_date utc_time utc_timestamp values var_pop var_samp variance week
weight_string year
`) {
		m[name] = true
	}
	return m
}()

// oceanbaseBuiltinFunctions are OceanBase CE 4.3.5.6 MySQL-mode specific
// built-in functions registered via REG_OP and get_function_alias_name in
// src/sql/engine/expr/ob_expr_operator_factory.cpp (conversion, vector,
// full-text, roaring bitmap, array/map, GIS _st_*, partition, SCN, crypto,
// audit, and diagnostic functions).
var oceanbaseBuiltinFunctions = func() map[string]bool {
	m := map[string]bool{}
	for _, name := range strings.Fields(`
_st_asewkb _st_asewkt _st_asmvtgeom _st_bestsrid _st_buffer _st_clipbybox2d _st_covers
_st_dwithin _st_geogfromtext _st_geographyfromtext _st_geohash _st_geometrytype
_st_geomfromewkb _st_geomfromewkt _st_iscollection _st_makeenvelope _st_makepoint
_st_makevalid _st_numinteriorrings _st_point _st_pointonsurface _st_setsrid _st_touches
_st_transform approx_cosine_distance approx_count_distinct
approx_count_distinct_synopsis approx_count_distinct_synopsis_merge approx_inner_product
approx_l1_distance approx_l2_distance approx_l2_squared array array_agg array_append
array_avg array_compact array_concat array_contains array_contains_all array_difference
array_distinct array_except array_filter array_first array_intersect array_length
array_map array_max array_min array_overlaps array_position array_prepend array_range
array_remove array_slice array_sort array_sortby array_sum array_to_string array_union
audit_log_filter_remove_filter audit_log_filter_remove_user audit_log_filter_set_filter
audit_log_filter_set_user bm25 calc_odps_size calc_partition_id calc_partition_tablet_id
calc_tablet_id cardinality check_catalog_access contains cosine_distance
cosine_similarity current_catalog current_scn current_user_priv decode decode_trace_id
des_decrypt des_encrypt des_hex_str doc_id doc_length effective_tenant
effective_tenant_id element_at encode enhanced_aes_decrypt enhanced_aes_encrypt
estimate_ndv extract_cert_expired_time get_path host_ip hybrid_hist inner_product
inner_product_similarity int2ip ip2int is_enabled_role is_serving_tenant key_version
keyvalue l1_distance l2_distance l2_norm l2_similarity l2_squared last_execution_id
last_refresh_scn last_trace_id listagg lnnvl map map_concat map_entries map_filter
map_keys map_values md5_concat_ws merging_frozen_time mysql_port mysql_proc_info
negative_inner_product nvl nvl2 ob_transaction_id ob_version ora_decode ora_trunc
part_id partition_hash percentile_cont percentile_disc random randstr range_placement
rb_and rb_and_agg rb_and_cardinality rb_and_cardinality_agg rb_and_null2empty
rb_and_null2empty_cardinality rb_andnot rb_andnot_cardinality rb_andnot_null2empty
rb_andnot_null2empty_cardinality rb_build rb_build_agg rb_build_empty rb_build_varbinary
rb_calc_cardinality rb_cardinality rb_contains rb_equals rb_from_string rb_insert
rb_is_empty rb_iterate rb_max rb_min rb_or rb_or_agg rb_or_cardinality
rb_or_cardinality_agg rb_or_null2empty rb_or_null2empty_cardinality rb_remove rb_select
rb_to_array rb_to_string rb_to_varbinary rb_xor rb_xor_cardinality rpc_port
scn_to_timestamp set_collation sm3 sm4_decrypt sm4_encrypt spiv_dim spiv_value split
split_part stmt_id string_to_array sum_opnsize svr_ip svr_port sys_op_opnsize
sys_privilege_check sys_view_bigint_param time_to_usec timestamp_nvl timestamp_to_scn
to_char to_outfile_row to_pinyin tokenize top_k_fre_hist uniform unnest url_decode
url_encode usec_to_time user_can_access_obj vec_cosine_distance vec_data
vec_inner_product vec_ivf_center_id vec_ivf_center_vector vec_ivf_flat_data_vector
vec_ivf_meta_id vec_ivf_meta_vector vec_ivf_pq_center_id vec_ivf_pq_center_ids
vec_ivf_pq_center_vector vec_ivf_sq8_data_vector vec_key vec_l1_distance
vec_l2_distance vec_l2_squared vec_negative_inner_product vec_scn vec_type vec_vector
vec_vid vector_dims vector_distance vector_norm vector_similarity word_count
word_segment ws zipf
`) {
		m[name] = true
	}
	return m
}()

// isBuiltinOceanBaseFunction reports whether lowerName is an OceanBase 4.3.5.6
// MySQL-mode built-in function.
func isBuiltinOceanBaseFunction(lowerName string) bool {
	return nativeFunctions[lowerName] ||
		grammarBuiltinFunctions[lowerName] ||
		oceanbaseBuiltinFunctions[lowerName]
}

func isSystemDatabaseName(name string) bool {
	switch toLower(name) {
	case "information_schema", "performance_schema", "mysql", "sys", "oceanbase":
		return true
	}
	return false
}
