// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// This file is derived from Apache Doris's fe/fe-core/src/main/cup/sql_parser.cup
// and fe/fe-core/src/main/jflex/sql_scanner.flex, which are beside it as
// published, with the license they come under (Doris.LICENSE). It was
// changed: doris/parser/cmd/cup2g4 wrote it from them, converting CUP to
// ANTLR 4 and keeping only what each rule derives. Actions, labels and
// %prec are gone, and so are the rules the start symbol does not reach. A
// keyword is written as the text sql_scanner.flex maps to it, punctuation
// as the text its scanner rule matches, and a class of texts as the
// DorisLexer.g4 token that matches the same texts, or a token of its own
// that tokens.tsv samples. Each rule is under the line of sql_parser.cup
// it came from.
//
// Do not edit: rerun the converter.

// sql_parser.cup declares terminals sql_scanner.flex never makes: COMMENTED_PLAN_HINTS, FLOATINGPOINT_LITERAL, KW_TREE.

parser grammar DorisLegacyParser;

// sql_parser.cup:1081
stmts
    : stmt
    | stmts ';' stmt
    | import_columns_stmt
    | import_where_stmt
    | import_delete_on_stmt
    | import_sequence_stmt
    | import_preceding_filter_stmt
    ;

// sql_parser.cup:1113
import_columns_stmt
    : 'COLUMNS' '(' import_column_descs ')'
    ;

// sql_parser.cup:1120
import_column_descs
    : import_column_desc
    | import_column_descs ',' import_column_desc
    | import_column_descs ',' '(' import_column_desc ')'
    ;

// sql_parser.cup:1137
import_column_desc
    : ident
    | ident '=' expr
    ;

// sql_parser.cup:1148
import_where_stmt
    : 'WHERE' expr
    ;

// sql_parser.cup:1155
import_delete_on_stmt
    : 'DELETE' 'ON' expr
    ;

// sql_parser.cup:1162
import_sequence_stmt
    : 'ORDER' 'BY' ident
    ;

// sql_parser.cup:1169
import_preceding_filter_stmt
    : 'PRECEDING' 'FILTER' expr
    ;

// sql_parser.cup:1176
stmt
    : alter_stmt
    | create_stmt
    | switch_stmt
    | query_stmt
    | drop_stmt
    | recover_stmt
    | use_stmt
    | use_cloud_cluster_stmt
    | set_stmt
    | unset_var_stmt
    | unset_default_storage_vault_stmt
    | kill_stmt
    | kill_analysis_job_stmt
    | describe_stmt
    | show_stmt
    | grant_stmt
    | revoke_stmt
    | help_stmt
    | load_stmt
    | create_routine_load_stmt
    | pause_routine_load_stmt
    | resume_routine_load_stmt
    | stop_routine_load_stmt
    | show_routine_load_stmt
    | show_routine_load_task_stmt
    | show_create_routine_load_stmt
    | create_job_stmt
    | pause_job_stmt
    | cancel_job_task_stmt
    | stop_job_stmt
    | resume_job_stmt
    | show_create_load_stmt
    | show_create_reporitory_stmt
    | cancel_stmt
    | delete_stmt
    | sync_stmt
    | insert_stmt
    | insert_overwrite_stmt
    | update_stmt
    | backup_stmt
    | restore_stmt
    | transaction_stmt
    | unsupported_stmt
    | export_stmt
    | admin_stmt
    | truncate_stmt
    | install_plugin_stmt
    | uninstall_plugin_stmt
    | lock_tables_stmt
    | unlock_tables_stmt
    | refresh_stmt
    | clean_stmt
    | analyze_stmt
    | warm_up_stmt
    | /* empty */
    | copy_stmt
    ;

// sql_parser.cup:1305
refresh_stmt
    : 'REFRESH' 'TABLE' table_name
    | 'REFRESH' 'DATABASE' ident opt_properties
    | 'REFRESH' 'DATABASE' ident '.' ident opt_properties
    | 'REFRESH' 'CATALOG' ident opt_properties
    | 'REFRESH' 'LDAP' 'ALL'
    | 'REFRESH' 'LDAP' opt_user
    ;

// sql_parser.cup:1332
clean_stmt
    : 'CLEAN' 'LABEL' from_or_in ident
    | 'CLEAN' 'LABEL' ident from_or_in ident
    | 'CLEAN' 'ALL' 'PROFILE'
    | 'CLEAN' 'QUERY' 'STATS' 'FOR' ident
    | 'CLEAN' 'ALL' 'QUERY' 'STATS'
    | 'CLEAN' 'QUERY' 'STATS' from_or_in table_name
    ;

// sql_parser.cup:1360
install_plugin_stmt
    : 'INSTALL' 'PLUGIN' 'FROM' ident_or_text opt_properties
    ;

// sql_parser.cup:1367
uninstall_plugin_stmt
    : 'UNINSTALL' 'PLUGIN' ident_or_text
    ;

// sql_parser.cup:1375
alter_stmt
    : 'ALTER' 'TABLE' table_name alter_table_clause_list
    | 'ALTER' 'TABLE' table_name 'ADD' 'ROLLUP' add_rollup_clause_list
    | 'ALTER' 'TABLE' table_name 'DROP' 'ROLLUP' drop_rollup_clause_list
    | 'ALTER' 'VIEW' table_name opt_col_with_comment_list 'AS' query_stmt
    | 'ALTER' 'SYSTEM' alter_system_clause
    | 'ALTER' 'DATABASE' ident 'SET' 'DATA' 'QUOTA' quantity
    | 'ALTER' 'DATABASE' ident 'SET' 'REPLICA' 'QUOTA' INTEGER_VALUE
    | 'ALTER' 'DATABASE' ident 'SET' 'TRANSACTION' 'QUOTA' INTEGER_VALUE
    | 'ALTER' 'DATABASE' ident 'RENAME' ident
    | 'ALTER' 'DATABASE' ident 'SET' 'PROPERTIES' '(' key_value_map ')'
    | 'ALTER' 'CATALOG' ident 'RENAME' ident
    | 'ALTER' 'CATALOG' ident 'SET' 'PROPERTIES' '(' key_value_map ')'
    | 'ALTER' 'CATALOG' ident 'MODIFY' 'COMMENT' STRING_LITERAL
    | 'ALTER' 'RESOURCE' ident_or_text opt_properties
    | 'ALTER' 'COLOCATE' 'GROUP' colocate_group_name 'SET' '(' key_value_map ')'
    | 'ALTER' 'WORKLOAD' 'GROUP' ident_or_text opt_properties
    | 'ALTER' 'WORKLOAD' 'POLICY' ident_or_text opt_properties
    | 'ALTER' 'ROUTINE' 'LOAD' 'FOR' job_label opt_properties opt_datasource_properties
    | 'ALTER' 'SQL_BLOCK_RULE' ident opt_properties
    | 'ALTER' 'TABLE' table_name 'SET' 'STATS' '(' key_value_map ')' opt_partition_names
    | 'ALTER' 'TABLE' table_name 'MODIFY' 'COLUMN' ident 'SET' 'STATS' '(' key_value_map ')' opt_partition_names
    | 'ALTER' 'TABLE' table_name 'INDEX' ident 'MODIFY' 'COLUMN' ident 'SET' 'STATS' '(' key_value_map ')' opt_partition_names
    | 'ALTER' 'TABLE' table_name 'SET' '(' key_value_map ')'
    | 'ALTER' 'STORAGE' 'POLICY' ident_or_text opt_properties
    | 'ALTER' 'USER' opt_if_exists grant_user opt_password_option opt_comment_null
    | 'ALTER' 'REPOSITORY' ident opt_properties
    ;

// sql_parser.cup:1495
warm_up_stmt
    : 'WARM' 'UP' 'CLUSTER' ident 'WITH' 'CLUSTER' ident opt_force opt_properties
    | 'WARM' 'UP' 'CLUSTER' ident 'WITH' warm_up_list opt_force
    | 'WARM' 'UP' 'COMPUTE' 'GROUP' ident 'WITH' 'COMPUTE' 'GROUP' ident opt_force opt_properties
    | 'WARM' 'UP' 'COMPUTE' 'GROUP' ident 'WITH' warm_up_list opt_force
    ;

// sql_parser.cup:1514
warm_up_item
    : 'TABLE' table_name
    | 'TABLE' table_name 'PARTITION' ident
    ;

// sql_parser.cup:1529
warm_up_list
    : warm_up_item
    | warm_up_list 'AND' warm_up_item
    ;

// sql_parser.cup:1541
opt_datasource_properties
    : /* empty */
    | 'FROM' ident '(' key_value_map ')'
    ;

// sql_parser.cup:1552
quantity
    : INTEGER_VALUE
    | ident
    ;

// sql_parser.cup:1563
opt_profile_path
    : /* empty */
    | STRING_LITERAL
    ;

// sql_parser.cup:1571
opt_user
    : /* empty */
    | 'FOR' user
    ;

// sql_parser.cup:1581
add_rollup_clause
    : ident '(' ident_list ')' opt_dup_keys opt_from_rollup opt_properties
    ;

// sql_parser.cup:1588
add_rollup_clause_list
    : add_rollup_clause
    | add_rollup_clause_list ',' add_rollup_clause
    ;

// sql_parser.cup:1600
drop_rollup_clause
    : ident opt_properties
    ;

// sql_parser.cup:1607
drop_rollup_clause_list
    : drop_rollup_clause
    | drop_rollup_clause_list ',' drop_rollup_clause
    ;

// sql_parser.cup:1619
alter_table_clause_list
    : alter_table_clause
    | alter_table_clause_list ',' alter_table_clause
    ;

// sql_parser.cup:1631
opt_to_rollup
    : /* empty */
    | 'TO' ident
    | 'IN' ident
    ;

// sql_parser.cup:1645
opt_from_rollup
    : /* empty */
    | 'FROM' ident
    ;

// sql_parser.cup:1655
opt_col_pos
    : /* empty */
    | 'FIRST'
    | 'AFTER' ident
    ;

// sql_parser.cup:1669
opt_dup_keys
    : /* empty */
    | 'DUPLICATE' 'KEY' '(' ident_list ')'
    ;

// sql_parser.cup:1679
alter_table_clause
    : 'ADD' 'COLUMN' column_definition opt_col_pos opt_to_rollup opt_properties
    | 'ADD' 'COLUMN' '(' column_definition_list ')' opt_to_rollup opt_properties
    | 'DROP' 'COLUMN' ident opt_from_rollup opt_properties
    | 'MODIFY' 'COLUMN' column_definition opt_col_pos opt_from_rollup opt_properties
    | 'ORDER' 'BY' '(' ident_list ')' opt_from_rollup opt_properties
    | 'ADD' opt_tmp single_partition_desc opt_distribution opt_properties
    | 'DROP' opt_tmp 'PARTITION' opt_if_exists ident opt_force
    | 'DROP' opt_tmp 'PARTITION' opt_if_exists ident opt_force 'FROM' 'INDEX' ident
    | 'MODIFY' opt_tmp 'PARTITION' ident 'SET' '(' key_value_map ')'
    | 'MODIFY' opt_tmp 'PARTITION' '(' ident_list ')' 'SET' '(' key_value_map ')'
    | 'MODIFY' opt_tmp 'PARTITION' '(' '*' ')' 'SET' '(' key_value_map ')'
    | 'REPLACE' opt_partition_names 'WITH' opt_partition_names opt_force opt_properties
    | 'REPLACE' 'WITH' 'TABLE' ident opt_properties
    | 'RENAME' ident
    | 'RENAME' 'ROLLUP' ident ident
    | 'RENAME' 'PARTITION' ident ident
    | 'RENAME' 'COLUMN' ident ident
    | 'ADD' index_definition
    | 'DROP' 'INDEX' opt_if_exists ident
    | 'ENABLE' 'FEATURE' STRING_LITERAL opt_enable_feature_properties
    | 'MODIFY' 'DISTRIBUTION' opt_distribution
    | 'MODIFY' 'COMMENT' STRING_LITERAL
    | 'MODIFY' 'COLUMN' ident 'COMMENT' STRING_LITERAL
    | 'MODIFY' 'ENGINE' 'TO' ident opt_properties
    | 'ADD' opt_tmp 'PARTITIONS' 'FROM' '(' partition_key_list ')' 'TO' '(' partition_key_list ')' 'INTERVAL' INTEGER_VALUE ident opt_properties
    | 'ADD' opt_tmp 'PARTITIONS' 'FROM' '(' partition_key_list ')' 'TO' '(' partition_key_list ')' 'INTERVAL' INTEGER_VALUE opt_properties
    | 'CREATE' 'BRANCH' opt_if_not_exists ident branch_options
    | 'REPLACE' 'BRANCH' ident branch_options
    | 'CREATE' 'OR' 'REPLACE' 'BRANCH' ident branch_options
    | 'CREATE' 'TAG' opt_if_not_exists ident tag_options
    | 'REPLACE' 'TAG' ident tag_options
    | 'CREATE' 'OR' 'REPLACE' 'TAG' ident tag_options
    | 'DROP' 'BRANCH' opt_if_exists ident
    | 'DROP' 'TAG' opt_if_exists ident
    ;

// sql_parser.cup:1828
branch_options
    : opt_as_of_version opt_retain_time opt_retention_snapshots
    ;

// sql_parser.cup:1835
opt_as_of_version
    : /* empty */
    | 'AS' 'OF' 'VERSION' INTEGER_VALUE
    ;

// sql_parser.cup:1845
opt_retain_time
    : /* empty */
    | 'RETAIN' time_value_with_unit
    ;

// sql_parser.cup:1855
time_value_with_unit
    : INTEGER_VALUE 'DAYS'
    | INTEGER_VALUE 'HOURS'
    | INTEGER_VALUE 'MINUTES'
    ;

// sql_parser.cup:1870
opt_retention_snapshots
    : /* empty */
    | 'WITH' 'SNAPSHOT' 'RETENTION' INTEGER_VALUE 'SNAPSHOTS'
    | 'WITH' 'SNAPSHOT' 'RETENTION' time_value_with_unit
    | 'WITH' 'SNAPSHOT' 'RETENTION' INTEGER_VALUE 'SNAPSHOTS' time_value_with_unit
    ;

// sql_parser.cup:1888
tag_options
    : opt_as_of_version opt_retain_time
    ;

// sql_parser.cup:1895
opt_enable_feature_properties
    : /* empty */
    | 'WITH' properties
    ;

// sql_parser.cup:1905
alter_system_clause
    : 'ADD' 'BACKEND' string_list opt_properties
    | 'DROP' 'BACKEND' string_list
    | 'DROPP' 'BACKEND' string_list
    | 'DECOMMISSION' 'BACKEND' string_list
    | 'ADD' 'OBSERVER' STRING_LITERAL
    | 'DROP' 'OBSERVER' STRING_LITERAL
    | 'ADD' 'FOLLOWER' STRING_LITERAL
    | 'DROP' 'FOLLOWER' STRING_LITERAL
    | 'ADD' 'BROKER' ident_or_text string_list
    | 'DROP' 'BROKER' ident_or_text string_list
    | 'DROP' 'ALL' 'BROKER' ident_or_text
    | 'SET' 'LOAD' 'ERRORS' 'HUB' opt_properties
    | 'MODIFY' 'BACKEND' string_list 'SET' '(' key_value_map ')'
    | 'MODIFY' 'FRONTEND' STRING_LITERAL 'HOSTNAME' STRING_LITERAL
    | 'MODIFY' 'BACKEND' STRING_LITERAL 'HOSTNAME' STRING_LITERAL
    ;

// sql_parser.cup:1971
sync_stmt
    : 'SYNC'
    ;

// sql_parser.cup:1978
opt_intermediate_type
    : /* empty */
    | 'INTERMEDIATE' type_def
    ;

// sql_parser.cup:1989
create_stmt
    : 'CREATE' 'DATABASE' opt_if_not_exists db_name opt_properties
    | 'CREATE' 'SCHEMA' opt_if_not_exists db_name
    | 'CREATE' 'CATALOG' opt_if_not_exists ident opt_comment opt_properties
    | 'CREATE' 'CATALOG' opt_if_not_exists ident 'WITH' 'RESOURCE' ident opt_comment opt_properties
    | 'CREATE' opt_var_type opt_aggregate 'FUNCTION' opt_if_not_exists function_name '(' func_args_def ')' 'RETURNS' type_def opt_intermediate_type opt_properties
    | 'CREATE' opt_var_type 'ALIAS' 'FUNCTION' opt_if_not_exists function_name '(' func_args_def ')' 'WITH' 'PARAMETER' '(' opt_ident_list ')' 'AS' expr
    | 'CREATE' opt_var_type 'TABLES' 'FUNCTION' opt_if_not_exists function_name '(' func_args_def ')' 'RETURNS' type_def opt_intermediate_type opt_properties
    | 'CREATE' opt_external 'TABLE' opt_if_not_exists table_name 'LIKE' table_name 'WITH' 'ROLLUP' '(' ident_list ')'
    | 'CREATE' opt_external 'TABLE' opt_if_not_exists table_name 'LIKE' table_name 'WITH' 'ROLLUP'
    | 'CREATE' opt_external 'TABLE' opt_if_not_exists table_name 'LIKE' table_name
    | 'CREATE' opt_external 'TABLE' opt_if_not_exists table_name '(' column_definition_list ')' opt_engine opt_keys opt_comment opt_partition opt_distribution opt_rollup opt_properties opt_ext_properties
    | 'CREATE' opt_external 'TABLE' opt_if_not_exists table_name '(' column_definition_list ',' ')' opt_engine opt_keys opt_comment opt_partition opt_distribution opt_rollup opt_properties opt_ext_properties
    | 'CREATE' opt_external 'TABLE' opt_if_not_exists table_name '(' column_definition_list ',' index_definition_list ')' opt_engine opt_keys opt_comment opt_partition opt_distribution opt_rollup opt_properties opt_ext_properties
    | 'CREATE' opt_external 'TABLE' opt_if_not_exists table_name '(' column_definition_list ',' index_definition_list ',' ')' opt_engine opt_keys opt_comment opt_partition opt_distribution opt_rollup opt_properties opt_ext_properties
    | 'CREATE' opt_external 'TABLE' opt_if_not_exists table_name opt_col_list opt_engine opt_keys opt_comment opt_partition opt_distribution opt_rollup opt_properties opt_ext_properties 'AS' query_stmt
    | 'CREATE' 'USER' opt_if_not_exists grant_user opt_user_role opt_password_option opt_comment
    | 'CREATE' opt_or_replace 'VIEW' opt_if_not_exists table_name opt_col_with_comment_list opt_comment 'AS' query_stmt
    | 'CREATE' opt_read_only 'REPOSITORY' ident 'WITH' storage_backend
    | 'CREATE' 'ROLE' opt_if_not_exists ident opt_comment
    | 'ALTER' 'ROLE' ident 'COMMENT' STRING_LITERAL
    | 'CREATE' 'FILE' STRING_LITERAL opt_db 'PROPERTIES' '(' key_value_map ')'
    | 'CREATE' 'MATERIALIZED' 'VIEW' ident 'AS' select_stmt opt_properties
    | 'CREATE' 'INDEX' opt_if_not_exists ident 'ON' table_name '(' ident_list ')' opt_index_type opt_properties opt_comment
    | 'CREATE' 'EXTERNAL' 'RESOURCE' opt_if_not_exists ident_or_text opt_properties
    | 'CREATE' 'RESOURCE' opt_if_not_exists ident_or_text opt_properties
    | 'CREATE' 'STORAGE' 'VAULT' opt_if_not_exists ident_or_text opt_properties
    | 'CREATE' 'WORKLOAD' 'GROUP' opt_if_not_exists ident_or_text opt_properties
    | 'CREATE' 'WORKLOAD' 'POLICY' opt_if_not_exists ident_or_text opt_conditions opt_actions opt_properties
    | 'CREATE' 'ENCRYPTKEY' opt_if_not_exists encryptkey_name 'AS' STRING_LITERAL
    | 'CREATE' 'SQL_BLOCK_RULE' opt_if_not_exists ident opt_properties
    | 'CREATE' 'ROW' 'POLICY' opt_if_not_exists ident 'ON' table_name 'AS' ident 'TO' user_identity 'USING' '(' expr ')'
    | 'CREATE' 'ROW' 'POLICY' opt_if_not_exists ident 'ON' table_name 'AS' ident 'TO' 'ROLE' ident 'USING' '(' expr ')'
    | 'CREATE' 'STORAGE' 'POLICY' opt_if_not_exists ident opt_properties
    | 'CREATE' 'INVERTED' 'INDEX' 'ANALYZER' opt_if_not_exists ident opt_properties
    | 'CREATE' 'INVERTED' 'INDEX' 'TOKENIZER' opt_if_not_exists ident opt_properties
    | 'CREATE' 'INVERTED' 'INDEX' 'TOKEN_FILTER' opt_if_not_exists ident opt_properties
    | 'BUILD' 'INDEX' ident 'ON' table_name opt_partition_names
    | 'CREATE' 'STAGE' opt_if_not_exists ident 'PROPERTIES' opt_key_value_map
    ;

// sql_parser.cup:2217
opt_aggregate
    : /* empty */
    | 'AGGREGATE'
    ;

// sql_parser.cup:2227
storage_backend
    : /* empty */
    | 'BROKER' ident 'ON' 'LOCATION' STRING_LITERAL opt_properties
    | 'S3' 'ON' 'LOCATION' STRING_LITERAL opt_properties
    | 'HDFS' 'ON' 'LOCATION' STRING_LITERAL opt_properties
    | 'LOCAL' 'ON' 'LOCATION' STRING_LITERAL opt_properties
    ;

// sql_parser.cup:2246
opt_read_only
    : /* empty */
    | 'READ' 'ONLY'
    ;

// sql_parser.cup:2256
grant_user
    : user_identity
    | user_identity 'IDENTIFIED' 'BY' STRING_LITERAL
    | user_identity 'IDENTIFIED' 'BY' 'PASSWORD' STRING_LITERAL
    ;

// sql_parser.cup:2274
opt_user_role
    : /* empty */
    | 'SUPERUSER'
    | 'DEFAULT' 'ROLE' STRING_LITERAL
    ;

// sql_parser.cup:2289
opt_password_option
    : opt_passwd_history_policy opt_passwd_expire_policy opt_passwd_reuse_policy opt_failed_login_attempts opt_password_lock_time opt_lock_account
    ;

// sql_parser.cup:2301
opt_passwd_expire_policy
    : /* empty */
    | 'PASSWORD_EXPIRE' passwd_expire_opt
    ;

// sql_parser.cup:2311
opt_passwd_history_policy
    : /* empty */
    | 'PASSWORD_HISTORY' passwd_history_opt
    ;

// sql_parser.cup:2321
opt_passwd_reuse_policy
    : /* empty */
    | 'PASSWORD_REUSE' 'INTERVAL' passwd_reuse_opt
    ;

// sql_parser.cup:2331
opt_lock_account
    : /* empty */
    | 'ACCOUNT_LOCK'
    | 'ACCOUNT_UNLOCK'
    ;

// sql_parser.cup:2345
passwd_expire_opt
    : 'DEFAULT'
    | 'NEVER'
    | 'INTERVAL' INTEGER_VALUE passwd_time_unit
    ;

// sql_parser.cup:2360
passwd_time_unit
    : 'DAY'
    | 'HOUR'
    | 'SECOND'
    ;

// sql_parser.cup:2375
passwd_history_opt
    : 'DEFAULT'
    | INTEGER_VALUE
    ;

// sql_parser.cup:2386
passwd_reuse_opt
    : 'DEFAULT'
    | INTEGER_VALUE 'DAY'
    ;

// sql_parser.cup:2397
opt_failed_login_attempts
    : /* empty */
    | 'FAILED_LOGIN_ATTEMPTS' INTEGER_VALUE
    ;

// sql_parser.cup:2407
opt_password_lock_time
    : /* empty */
    | 'PASSWORD_LOCK_TIME' passwd_lock_time_opt
    ;

// sql_parser.cup:2417
passwd_lock_time_opt
    : INTEGER_VALUE passwd_time_unit
    | 'UNBOUNDED'
    ;

// sql_parser.cup:2428
user
    : ident_or_text
    ;

// sql_parser.cup:2435
user_identity
    : ident_or_text
    | ident_or_text '@' ident_or_text
    | ident_or_text '@' '[' ident_or_text ']'
    ;

// sql_parser.cup:2451
help_stmt
    : 'HELP' ident_or_text
    ;

// sql_parser.cup:2459
export_stmt
    : 'EXPORT' 'TABLE' base_table_ref where_clause 'TO' STRING_LITERAL opt_properties opt_broker
    ;

// sql_parser.cup:2471
load_stmt
    : 'LOAD' 'LABEL' job_label '(' data_desc_list ')' opt_broker opt_properties opt_comment
    | 'LOAD' 'LABEL' job_label '(' data_desc_list ')' resource_desc opt_properties opt_comment
    | 'LOAD' mysql_data_desc opt_properties opt_comment
    ;

// sql_parser.cup:2494
job_label
    : ident
    | ident '.' ident
    ;

// sql_parser.cup:2506
copy_stmt
    : 'COPY' 'INTO' opt_select_hints table_name opt_col_list 'FROM' copy_from_param 'PROPERTIES' opt_key_value_map
    | 'COPY' 'INTO' opt_select_hints table_name opt_col_list 'FROM' copy_from_param
    ;

// sql_parser.cup:2517
copy_from_param
    : stage_and_pattern
    | '(' 'SELECT' copy_select_expr_list 'FROM' stage_and_pattern where_clause ')'
    ;

// sql_parser.cup:2528
stage_name
    : '@' ident
    | '@' '~'
    ;

// sql_parser.cup:2539
stage_and_pattern
    : stage_name
    | stage_name '(' STRING_LITERAL ')'
    ;

// sql_parser.cup:2550
copy_select_expr_list
    : /* empty */
    | '*'
    | expr_list
    ;

// sql_parser.cup:2561
data_desc_list
    : data_desc
    | data_desc_list ',' data_desc
    ;

// sql_parser.cup:2573
opt_merge_type
    : /* empty */
    | 'APPEND'
    | 'DELETE'
    | 'MERGE'
    | 'WITH' 'APPEND'
    | 'WITH' 'DELETE'
    | 'WITH' 'MERGE'
    ;

// sql_parser.cup:2603
opt_with_merge_type
    : /* empty */
    | 'WITH' 'APPEND'
    | 'WITH' 'DELETE'
    | 'WITH' 'MERGE'
    ;

// sql_parser.cup:2621
data_desc
    : opt_merge_type 'DATA' 'INFILE' '(' string_list ')' opt_negative 'INTO' 'TABLE' ident opt_partition_names opt_field_term opt_line_term opt_file_format opt_file_compress_type opt_col_list opt_columns_from_path opt_col_mapping_list pre_filter_clause where_clause delete_on_clause sequence_col_clause opt_properties
    | opt_merge_type 'DATA' 'FROM' 'TABLE' ident opt_negative 'INTO' 'TABLE' ident opt_partition_names opt_col_mapping_list where_clause delete_on_clause opt_properties
    ;

// sql_parser.cup:2658
mysql_data_desc
    : 'DATA' opt_local 'INFILE' STRING_LITERAL 'INTO' 'TABLE' table_name opt_partition_names opt_field_term opt_line_term opt_skip_lines opt_col_list opt_col_mapping_list opt_properties
    ;

// sql_parser.cup:2678
opt_negative
    : /* empty */
    | 'NEGATIVE'
    ;

// sql_parser.cup:2688
opt_local
    : /* empty */
    | 'LOCAL'
    ;

// sql_parser.cup:2698
opt_field_term
    : /* empty */
    | 'COLUMNS' 'TERMINATED' 'BY' STRING_LITERAL
    ;

// sql_parser.cup:2709
opt_line_term
    : /* empty */
    | 'LINES' 'TERMINATED' 'BY' STRING_LITERAL
    ;

// sql_parser.cup:2720
opt_skip_lines
    : /* empty */
    | 'IGNORE' INTEGER_VALUE 'LINES'
    | 'IGNORE' INTEGER_VALUE 'ROWS'
    ;

// sql_parser.cup:2735
separator
    : 'COLUMNS' 'TERMINATED' 'BY' STRING_LITERAL
    ;

// sql_parser.cup:2742
opt_file_format
    : /* empty */
    | 'FORMAT' 'AS' ident_or_text
    ;

// sql_parser.cup:2749
opt_file_compress_type
    : /* empty */
    | 'COMPRESS_TYPE' 'AS' ident_or_text
    ;

// sql_parser.cup:2756
opt_columns_from_path
    : /* empty */
    | 'COLUMNS' 'FROM' 'PATH' 'AS' '(' ident_list ')'
    ;

// sql_parser.cup:2763
opt_col_list
    : /* empty */
    | '(' ident_list ')'
    ;

// sql_parser.cup:2773
opt_col_with_comment_list
    : /* empty */
    | '(' col_with_comment_list ')'
    ;

// sql_parser.cup:2783
col_with_comment_list
    : col_with_comment
    | col_with_comment_list ',' col_with_comment
    ;

// sql_parser.cup:2797
col_with_comment
    : ident opt_comment
    ;

// sql_parser.cup:2804
opt_col_mapping_list
    : /* empty */
    | 'SET' '(' expr_list ')'
    ;

// sql_parser.cup:2815
opt_broker
    : /* empty */
    | 'WITH' 'S3' '(' key_value_map ')'
    | 'WITH' 'HDFS' '(' key_value_map ')'
    | 'WITH' 'LOCAL' '(' key_value_map ')'
    | 'WITH' 'BROKER' ident_or_text
    | 'WITH' 'BROKER' ident_or_text '(' key_value_map ')'
    ;

// sql_parser.cup:2841
resource_desc
    : 'WITH' 'RESOURCE' ident_or_text
    | 'WITH' 'RESOURCE' ident_or_text '(' key_value_map ')'
    ;

// sql_parser.cup:2851
create_job_stmt
    : 'CREATE' 'JOB' job_label 'ON' 'SCHEDULE' 'EVERY' INTEGER_VALUE ident opt_job_starts opt_job_ends opt_comment 'DO' stmt
    | 'CREATE' 'JOB' job_label 'ON' 'SCHEDULE' job_at_time opt_comment 'DO' stmt
    ;

// sql_parser.cup:2868
opt_job_starts
    : /* empty */
    | 'STARTS' STRING_LITERAL
    | 'STARTS' 'CURRENT_TIMESTAMP'
    ;

// sql_parser.cup:2882
job_at_time
    : /* empty */
    | 'AT' STRING_LITERAL
    | 'AT' 'CURRENT_TIMESTAMP'
    ;

// sql_parser.cup:2892
opt_job_ends
    : /* empty */
    | 'ENDS' STRING_LITERAL
    ;

// sql_parser.cup:2901
pause_job_stmt
    : 'PAUSE' 'JOB' opt_wild_where
    ;

// sql_parser.cup:2908
stop_job_stmt
    : 'DROP' 'JOB' opt_if_exists opt_wild_where
    ;

// sql_parser.cup:2914
resume_job_stmt
    : 'RESUME' 'JOB' opt_wild_where
    ;

// sql_parser.cup:2920
cancel_job_task_stmt
    : 'CANCEL' 'TASK' opt_wild_where
    ;

// sql_parser.cup:2927
create_routine_load_stmt
    : 'CREATE' 'ROUTINE' 'LOAD' job_label optional_on_ident opt_with_merge_type opt_load_property_list opt_properties 'FROM' ident '(' key_value_map ')' opt_comment
    ;

// sql_parser.cup:2939
optional_on_ident
    : 'ON' ident
    | /* empty */
    ;

// sql_parser.cup:2942
opt_load_property_list
    : /* empty */
    | load_property
    | opt_load_property_list ',' load_property
    ;

// sql_parser.cup:2957
load_property
    : separator
    | import_columns_stmt
    | import_preceding_filter_stmt
    | import_where_stmt
    | import_delete_on_stmt
    | import_sequence_stmt
    | partition_names
    ;

// sql_parser.cup:2988
pause_routine_load_stmt
    : 'PAUSE' 'ROUTINE' 'LOAD' 'FOR' job_label
    | 'PAUSE' 'ALL' 'ROUTINE' 'LOAD'
    ;

// sql_parser.cup:2999
resume_routine_load_stmt
    : 'RESUME' 'ROUTINE' 'LOAD' 'FOR' job_label
    | 'RESUME' 'ALL' 'ROUTINE' 'LOAD'
    ;

// sql_parser.cup:3010
stop_routine_load_stmt
    : 'STOP' 'ROUTINE' 'LOAD' 'FOR' job_label
    ;

// sql_parser.cup:3017
show_routine_load_stmt
    : 'SHOW' 'ROUTINE' 'LOAD' 'FOR' job_label
    | 'SHOW' 'ALL' 'ROUTINE' 'LOAD' 'FOR' job_label
    | 'SHOW' 'ROUTINE' 'LOAD' opt_wild_where
    | 'SHOW' 'ALL' 'ROUTINE' 'LOAD' opt_wild_where
    ;

// sql_parser.cup:3036
show_routine_load_task_stmt
    : 'SHOW' 'ROUTINE' 'LOAD' 'TASK' opt_db opt_wild_where
    ;

// sql_parser.cup:3043
show_create_routine_load_stmt
    : 'SHOW' 'CREATE' 'ROUTINE' 'LOAD' 'FOR' job_label
    | 'SHOW' 'ALL' 'CREATE' 'ROUTINE' 'LOAD' 'FOR' job_label
    ;

// sql_parser.cup:3054
show_create_load_stmt
    : 'SHOW' 'CREATE' 'LOAD' 'FOR' job_label
    ;

// sql_parser.cup:3061
show_create_reporitory_stmt
    : 'SHOW' 'CREATE' 'REPOSITORY' 'FOR' ident
    ;

// sql_parser.cup:3069
analyze_stmt
    : 'ANALYZE' 'TABLE' table_name opt_partition_names opt_col_list opt_with_analysis_properties opt_properties
    | 'ANALYZE' 'DATABASE' ident '.' ident opt_with_analysis_properties opt_properties
    | 'ANALYZE' 'DATABASE' ident opt_with_analysis_properties opt_properties
    ;

// sql_parser.cup:3122
grant_stmt
    : 'GRANT' privilege_list 'ON' tbl_pattern 'TO' user_identity
    | 'GRANT' privilege_list 'ON' tbl_pattern 'TO' 'ROLE' STRING_LITERAL
    | 'GRANT' privilege_list 'ON' 'RESOURCE' resource_pattern 'TO' user_identity
    | 'GRANT' privilege_list 'ON' 'RESOURCE' resource_pattern 'TO' 'ROLE' STRING_LITERAL
    | 'GRANT' privilege_list 'ON' 'CLUSTER' resource_pattern 'TO' user_identity
    | 'GRANT' privilege_list 'ON' 'CLUSTER' resource_pattern 'TO' 'ROLE' STRING_LITERAL
    | 'GRANT' privilege_list 'ON' 'COMPUTE' 'GROUP' resource_pattern 'TO' user_identity
    | 'GRANT' privilege_list 'ON' 'COMPUTE' 'GROUP' resource_pattern 'TO' 'ROLE' STRING_LITERAL
    | 'GRANT' privilege_list 'ON' 'STAGE' resource_pattern 'TO' user_identity
    | 'GRANT' privilege_list 'ON' 'STAGE' resource_pattern 'TO' 'ROLE' STRING_LITERAL
    | 'GRANT' privilege_list 'ON' 'STORAGE' 'VAULT' resource_pattern 'TO' user_identity
    | 'GRANT' privilege_list 'ON' 'STORAGE' 'VAULT' resource_pattern 'TO' 'ROLE' STRING_LITERAL
    | 'GRANT' privilege_list 'ON' 'WORKLOAD' 'GROUP' workload_group_pattern 'TO' user_identity
    | 'GRANT' privilege_list 'ON' 'WORKLOAD' 'GROUP' workload_group_pattern 'TO' 'ROLE' STRING_LITERAL
    | 'GRANT' string_list 'TO' user_identity
    ;

// sql_parser.cup:3185
tbl_pattern
    : ident_or_star
    | ident_or_star '.' ident_or_star
    | ident_or_star '.' ident_or_star '.' ident_or_star
    ;

// sql_parser.cup:3200
resource_pattern
    : ident_or_star
    | STRING_LITERAL
    ;

// sql_parser.cup:3211
workload_group_pattern
    : ident_or_star
    | STRING_LITERAL
    ;

// sql_parser.cup:3222
ident_or_star
    : '*'
    | ident
    ;

// sql_parser.cup:3234
revoke_stmt
    : 'REVOKE' privilege_list 'ON' tbl_pattern 'FROM' user_identity
    | 'REVOKE' privilege_list 'ON' tbl_pattern 'FROM' 'ROLE' STRING_LITERAL
    | 'REVOKE' privilege_list 'ON' 'RESOURCE' resource_pattern 'FROM' user_identity
    | 'REVOKE' privilege_list 'ON' 'RESOURCE' resource_pattern 'FROM' 'ROLE' STRING_LITERAL
    | 'REVOKE' privilege_list 'ON' 'CLUSTER' resource_pattern 'FROM' user_identity
    | 'REVOKE' privilege_list 'ON' 'CLUSTER' resource_pattern 'FROM' 'ROLE' STRING_LITERAL
    | 'REVOKE' privilege_list 'ON' 'COMPUTE' 'GROUP' resource_pattern 'FROM' user_identity
    | 'REVOKE' privilege_list 'ON' 'COMPUTE' 'GROUP' resource_pattern 'FROM' 'ROLE' STRING_LITERAL
    | 'REVOKE' privilege_list 'ON' 'STAGE' resource_pattern 'FROM' user_identity
    | 'REVOKE' privilege_list 'ON' 'STAGE' resource_pattern 'FROM' 'ROLE' STRING_LITERAL
    | 'REVOKE' privilege_list 'ON' 'STORAGE' 'VAULT' resource_pattern 'FROM' user_identity
    | 'REVOKE' privilege_list 'ON' 'STORAGE' 'VAULT' resource_pattern 'FROM' 'ROLE' STRING_LITERAL
    | 'REVOKE' privilege_list 'ON' 'WORKLOAD' 'GROUP' workload_group_pattern 'FROM' user_identity
    | 'REVOKE' privilege_list 'ON' 'WORKLOAD' 'GROUP' workload_group_pattern 'FROM' 'ROLE' STRING_LITERAL
    | 'REVOKE' string_list 'FROM' user_identity
    ;

// sql_parser.cup:3298
drop_stmt
    : 'DROP' 'DATABASE' opt_if_exists db_name opt_force
    | 'DROP' 'SCHEMA' opt_if_exists db_name opt_force
    | 'DROP' 'CATALOG' opt_if_exists ident
    | 'DROP' opt_var_type 'FUNCTION' opt_if_exists function_name '(' func_args_def ')'
    | 'DROP' opt_tmp 'TABLE' opt_if_exists table_name opt_force
    | 'DROP' 'USER' opt_if_exists user_identity
    | 'DROP' 'VIEW' opt_if_exists table_name
    | 'DROP' 'REPOSITORY' ident
    | 'DROP' 'ROLE' opt_if_exists ident
    | 'DROP' 'FILE' STRING_LITERAL opt_db 'PROPERTIES' '(' key_value_map ')'
    | 'DROP' 'INDEX' opt_if_exists ident 'ON' table_name
    | 'DROP' 'MATERIALIZED' 'VIEW' opt_if_exists ident 'ON' table_name
    | 'DROP' 'RESOURCE' opt_if_exists ident_or_text
    | 'DROP' 'WORKLOAD' 'GROUP' opt_if_exists ident_or_text
    | 'DROP' 'WORKLOAD' 'POLICY' opt_if_exists ident_or_text
    | 'DROP' 'ENCRYPTKEY' opt_if_exists encryptkey_name
    | 'DROP' 'SQL_BLOCK_RULE' opt_if_exists ident_list
    | 'DROP' 'ROW' 'POLICY' opt_if_exists ident 'ON' table_name
    | 'DROP' 'ROW' 'POLICY' opt_if_exists ident 'ON' table_name 'FOR' user_identity
    | 'DROP' 'ROW' 'POLICY' opt_if_exists ident 'ON' table_name 'FOR' 'ROLE' ident
    | 'DROP' 'STORAGE' 'POLICY' opt_if_exists ident
    | 'DROP' 'INVERTED' 'INDEX' 'ANALYZER' opt_if_exists ident
    | 'DROP' 'INVERTED' 'INDEX' 'TOKENIZER' opt_if_exists ident
    | 'DROP' 'INVERTED' 'INDEX' 'TOKEN_FILTER' opt_if_exists ident
    | 'DROP' 'STATS' table_name opt_col_list opt_partition_names
    | 'DROP' 'CACHED' 'STATS' table_name
    | 'DROP' 'EXPIRED' 'STATS'
    | 'DROP' 'ANALYZE' 'JOB' INTEGER_VALUE
    | 'DROP' 'STAGE' opt_if_exists ident
    ;

// sql_parser.cup:3425
recover_stmt
    : 'RECOVER' 'DATABASE' ident opt_id opt_alias
    | 'RECOVER' 'TABLE' table_name opt_id opt_alias
    | 'RECOVER' 'PARTITION' ident opt_id opt_alias 'FROM' table_name
    ;

// sql_parser.cup:3440
opt_agg_type
    : 'SUM'
    | 'MAX'
    | 'MIN'
    | 'REPLACE'
    | 'REPLACE_IF_NOT_NULL'
    | 'HLL_UNION'
    | 'BITMAP_UNION'
    | 'QUANTILE_UNION'
    | 'GENERIC'
    ;

// sql_parser.cup:3479
opt_partition
    : /* empty */
    | 'PARTITION' 'BY' 'RANGE' '(' ident_list ')' '(' opt_all_partition_desc_list ')'
    | 'PARTITION' 'BY' 'LIST' '(' ident_list ')' '(' opt_all_partition_desc_list ')'
    | 'AUTO' 'PARTITION' 'BY' 'RANGE' '(' function_call_expr ')' '(' opt_all_partition_desc_list ')'
    | 'AUTO' 'PARTITION' 'BY' 'RANGE' '(' function_name '(' expr_list ',' 'INTERVAL' expr ident ')' ')' '(' opt_all_partition_desc_list ')'
    | 'AUTO' 'PARTITION' 'BY' 'LIST' '(' expr_list ')' '(' opt_all_partition_desc_list ')'
    | 'AUTO' 'PARTITION' 'BY' 'LIST' function_call_expr '(' opt_all_partition_desc_list ')'
    ;

// sql_parser.cup:3528
opt_distribution
    : /* empty */
    | 'DISTRIBUTED' 'BY' 'HASH' '(' ident_list ')' opt_distribution_number
    | 'DISTRIBUTED' 'BY' 'RANDOM' opt_distribution_number
    ;

// sql_parser.cup:3549
opt_rollup
    : /* empty */
    | 'ROLLUP' '(' add_rollup_clause_list ')'
    ;

// sql_parser.cup:3560
opt_distribution_number
    : /* empty */
    | 'BUCKETS' INTEGER_VALUE
    | 'BUCKETS' 'AUTO'
    ;

// sql_parser.cup:3576
opt_keys
    : /* empty */
    | 'DUPLICATE' 'KEY' '(' ident_list ')'
    | 'UNIQUE' 'KEY' '(' ident_list ')' opt_cluster_keys
    | 'AGGREGATE' 'KEY' '(' ident_list ')'
    ;

// sql_parser.cup:3599
opt_cluster_keys
    : /* empty */
    | 'CLUSTER' 'BY' '(' ident_list ')'
    ;

// sql_parser.cup:3610
opt_all_partition_desc_list
    : /* empty */
    | all_partition_desc_list
    ;

// sql_parser.cup:3621
all_partition_desc_list
    : all_partition_desc_list ',' single_partition_desc
    | single_partition_desc
    | all_partition_desc_list ',' multi_partition_desc
    | multi_partition_desc
    ;

// sql_parser.cup:3642
single_partition_desc
    : 'PARTITION' opt_if_not_exists ident 'VALUES' 'LESS' 'THAN' partition_key_desc opt_key_value_map
    | 'PARTITION' opt_if_not_exists ident 'VALUES' fixed_partition_key_desc opt_key_value_map
    | 'PARTITION' opt_if_not_exists ident list_partition_key_desc opt_key_value_map
    | 'PARTITION' opt_if_not_exists ident 'VALUES' 'IN' list_partition_key_desc opt_key_value_map
    ;

// sql_parser.cup:3666
multi_partition_desc
    : fixed_multi_partition_key_desc opt_key_value_map
    ;

// sql_parser.cup:3674
fixed_multi_partition_key_desc
    : 'FROM' '(' partition_key_list ')' 'TO' '(' partition_key_list ')' 'INTERVAL' INTEGER_VALUE ident
    | 'FROM' '(' partition_key_list ')' 'TO' '(' partition_key_list ')' 'INTERVAL' INTEGER_VALUE
    ;

// sql_parser.cup:3686
partition_key_desc
    : 'MAXVALUE'
    | '(' partition_key_list ')'
    ;

// sql_parser.cup:3698
list_partition_key_desc
    : /* empty */
    | '(' list_partition_values_list ')'
    ;

// sql_parser.cup:3712
list_partition_values_list
    : partition_value_list
    | list_partition_values_list ',' partition_value_list
    ;

// sql_parser.cup:3727
partition_value_list
    : STRING_LITERAL
    | INTEGER_VALUE
    | '(' partition_key_item_list ')'
    | 'NULL'
    ;

// sql_parser.cup:3749
partition_key_item_list
    : STRING_LITERAL
    | partition_key_item_list ',' STRING_LITERAL
    | INTEGER_VALUE
    | partition_key_item_list ',' INTEGER_VALUE
    | 'NULL'
    | partition_key_item_list ',' 'NULL'
    | 'MAXVALUE'
    | partition_key_item_list ',' 'MAXVALUE'
    ;

// sql_parser.cup:3788
partition_key_list
    : /* empty */
    | partition_key_list ',' STRING_LITERAL
    | partition_key_list ',' INTEGER_VALUE
    | partition_key_list ',' 'MAXVALUE'
    | STRING_LITERAL
    | INTEGER_VALUE
    | 'MAXVALUE'
    | 'NULL'
    ;

// sql_parser.cup:3827
fixed_partition_key_desc
    : '[' '(' partition_key_list ')' ',' '(' partition_key_list ')' ')'
    ;

// sql_parser.cup:3835
opt_engine
    : /* empty */
    | 'ENGINE' '=' ident
    ;

// sql_parser.cup:3841
opt_key_value_map
    : /* empty */
    | '(' key_value_map ')'
    ;

// sql_parser.cup:3851
opt_key_value_map_in_paren
    : /* empty */
    | STRING_LITERAL '=' STRING_LITERAL
    | opt_key_value_map_in_paren ',' STRING_LITERAL '=' STRING_LITERAL
    ;

// sql_parser.cup:3868
key_value_map
    : STRING_LITERAL '=' STRING_LITERAL
    | key_value_map ',' STRING_LITERAL '=' STRING_LITERAL
    ;

// sql_parser.cup:3881
type_opt_properties
    : ',' opt_properties
    | /* empty */
    ;

// sql_parser.cup:3888
opt_properties
    : /* empty */
    | properties
    ;

// sql_parser.cup:3898
policy_condition_op
    : '='
    | '>'
    | '<'
    | '>' '='
    | '<' '='
    ;

// sql_parser.cup:3920
policy_condition_value
    : INTEGER_VALUE
    | STRING_LITERAL
    | FLOATINGPOINT_LITERAL
    | DECIMAL_VALUE
    ;

// sql_parser.cup:3939
workload_policy_condition_list
    : ident policy_condition_op policy_condition_value
    | workload_policy_condition_list ',' ident policy_condition_op policy_condition_value
    ;

// sql_parser.cup:3953
conditions
    : 'CONDITIONS' '(' workload_policy_condition_list ')'
    ;

// sql_parser.cup:3960
opt_conditions
    : /* empty */
    | conditions
    ;

// sql_parser.cup:3970
workload_policy_action_list
    : 'SET_SESSION_VARIABLE' STRING_LITERAL
    | workload_policy_action_list ',' 'SET_SESSION_VARIABLE' STRING_LITERAL
    | ident
    | ident STRING_LITERAL
    | workload_policy_action_list ',' ident
    | workload_policy_action_list ',' ident STRING_LITERAL
    ;

// sql_parser.cup:4010
actions
    : 'ACTIONS' '(' workload_policy_action_list ')'
    ;

// sql_parser.cup:4017
opt_actions
    : /* empty */
    | actions
    ;

// sql_parser.cup:4027
opt_ext_properties
    : /* empty */
    | 'BROKER' properties
    ;

// sql_parser.cup:4037
properties
    : 'PROPERTIES' '(' key_value_map ')'
    ;

// sql_parser.cup:4044
column_definition_list
    : column_definition
    | column_definition_list ',' column_definition
    ;

// sql_parser.cup:4057
index_definition_list
    : index_definition
    | index_definition_list ',' index_definition
    ;

// sql_parser.cup:4070
opt_default_value
    : /* empty */
    | 'DEFAULT' STRING_LITERAL
    | 'DEFAULT' 'NULL'
    | 'DEFAULT' 'CURRENT_TIMESTAMP'
    | 'DEFAULT' 'CURRENT_TIMESTAMP' '(' INTEGER_VALUE ')'
    | 'DEFAULT' 'BITMAP_EMPTY'
    | 'DEFAULT' INTEGER_VALUE
    | 'DEFAULT' LARGE_INTEGER_LITERAL
    | 'DEFAULT' DECIMAL_VALUE
    ;

// sql_parser.cup:4109
opt_is_key
    : /* empty */
    | 'KEY'
    ;

// sql_parser.cup:4119
column_definition
    : ident type_def opt_is_key opt_nullable_type opt_auto_inc_init_value opt_default_value opt_comment
    | ident type_def opt_is_key opt_agg_type opt_nullable_type opt_auto_inc_init_value opt_default_value opt_comment
    | ident type_def opt_is_key opt_generated_always 'AS' '(' expr ')' opt_nullable_type opt_comment
    ;

// sql_parser.cup:4137
index_definition
    : 'INDEX' opt_if_not_exists ident '(' ident_list ')' opt_index_type opt_properties opt_comment
    ;

// sql_parser.cup:4144
opt_nullable_type
    : /* empty */
    | 'NULL'
    | 'NOT' 'NULL'
    ;

// sql_parser.cup:4158
opt_is_allow_null
    : /* empty */
    | 'NULL'
    | 'NOT' 'NULL'
    ;

// sql_parser.cup:4172
opt_auto_inc_init_value
    : /* empty */
    | 'AUTO_INCREMENT'
    | 'AUTO_INCREMENT' '(' INTEGER_VALUE ')'
    ;

// sql_parser.cup:4189
opt_comment
    : /* empty */
    | 'COMMENT' STRING_LITERAL
    ;

// sql_parser.cup:4200
opt_comment_null
    : /* empty */
    | 'COMMENT' STRING_LITERAL
    ;

// sql_parser.cup:4211
opt_index_type
    : /* empty */
    | 'USING' 'BITMAP'
    | 'USING' 'NGRAM_BF'
    | 'USING' 'INVERTED'
    ;

// sql_parser.cup:4233
opt_or_replace
    : /* empty */
    | 'OR' 'REPLACE'
    ;

// sql_parser.cup:4243
opt_if_exists
    : /* empty */
    | 'IF' 'EXISTS'
    ;

// sql_parser.cup:4253
opt_if_not_exists
    : /* empty */
    | 'IF' 'NOT' 'EXISTS'
    ;

// sql_parser.cup:4263
opt_external
    : /* empty */
    | 'EXTERNAL'
    ;

// sql_parser.cup:4274
opt_cached
    : /* empty */
    | 'CACHED'
    ;

// sql_parser.cup:4285
opt_force
    : /* empty */
    | 'FORCE'
    ;

// sql_parser.cup:4297
show_stmt
    : 'SHOW' show_param
    | 'SHOW' 'SQL_BLOCK_RULE' 'FOR' ident
    | 'SHOW' 'SQL_BLOCK_RULE'
    | 'SHOW' 'ROW' 'POLICY' 'FOR' user_identity
    | 'SHOW' 'ROW' 'POLICY' 'FOR' 'ROLE' ident
    | 'SHOW' 'ROW' 'POLICY'
    | 'SHOW' 'STORAGE' 'POLICY'
    | 'SHOW' 'STORAGE' 'POLICY' 'USING'
    | 'SHOW' 'STORAGE' 'POLICY' 'USING' 'FOR' ident_or_text
    | 'SHOW' 'STAGES'
    | 'SHOW' 'STORAGE' 'VAULT'
    | 'SHOW' 'STORAGE' 'VAULTS'
    | 'SHOW' 'INVERTED' 'INDEX' 'ANALYZER'
    | 'SHOW' 'INVERTED' 'INDEX' 'TOKENIZER'
    | 'SHOW' 'INVERTED' 'INDEX' 'TOKEN_FILTER'
    ;

// sql_parser.cup:4361
show_param
    : 'WHITELIST'
    | opt_var_type 'VARIABLES' opt_wild_where
    | 'OPEN' 'TABLES' opt_db opt_wild_where
    | 'TABLE' 'STATUS' opt_db opt_wild_where
    | 'TABLE' 'STATUS' from_or_in ident '.' ident opt_wild_where
    | opt_full 'TABLES' opt_db opt_wild_where
    | opt_full 'TABLES' from_or_in ident '.' ident opt_wild_where
    | opt_full 'VIEWS' opt_db opt_wild_where
    | opt_full 'VIEWS' from_or_in ident '.' ident opt_wild_where
    | 'TABLE' INTEGER_VALUE
    | opt_full 'PROCESSLIST'
    | procedure_or_function 'STATUS' opt_wild_where
    | opt_var_type 'STATUS' opt_wild_where
    | opt_full 'TRIGGERS' opt_db opt_wild_where
    | 'EVENTS' opt_db opt_wild_where
    | 'PLUGINS'
    | opt_storage 'ENGINES'
    | 'AUTHORS'
    | 'CREATE' 'TABLE' table_name
    | 'BRIEF' 'CREATE' 'TABLE' table_name
    | 'CREATE' 'VIEW' table_name
    | 'CREATE' 'MATERIALIZED' 'VIEW' table_name
    | 'CREATE' 'DATABASE' db_name
    | 'CREATE' 'SCHEMA' db_name
    | 'CREATE' 'CATALOG' ident
    | 'CREATE' opt_var_type 'FUNCTION' function_name '(' func_args_def ')' opt_db
    | 'DATABASES' opt_wild_where
    | 'DATABASES' 'FROM' ident opt_wild_where
    | 'DATABASE' INTEGER_VALUE
    | 'DATA' 'TYPES'
    | 'SCHEMAS' opt_wild_where
    | 'SCHEMAS' 'FROM' ident opt_wild_where
    | 'CATALOGS' opt_wild_where
    | 'CATALOG' ident
    | 'DYNAMIC' 'PARTITION' 'TABLES' opt_db
    | opt_full 'COLUMNS' from_or_in table_name opt_db opt_wild_where
    | opt_full 'FIELDS' from_or_in table_name opt_db opt_wild_where
    | 'COLLATION' opt_wild_where
    | charset opt_wild_where
    | 'PROC' STRING_LITERAL
    | 'COUNT' '(' '*' ')' 'WARNINGS'
    | 'COUNT' '(' '*' ')' 'ERRORS'
    | 'WARNINGS' limit_clause
    | 'ERRORS' limit_clause
    | 'LOAD' 'WARNINGS' opt_db opt_wild_where limit_clause
    | 'LOAD' 'WARNINGS' 'ON' STRING_LITERAL
    | 'LOAD' opt_db opt_wild_where order_by_clause limit_clause
    | 'STREAM' 'LOAD' opt_db opt_wild_where order_by_clause limit_clause
    | 'EXPORT' opt_db opt_wild_where order_by_clause limit_clause
    | 'EXPORT' from_or_in ident '.' ident opt_wild_where order_by_clause limit_clause
    | 'DELETE' opt_db
    | 'ALTER' 'TABLE' opt_alter_type opt_db opt_wild_where order_by_clause limit_clause
    | 'DATA' 'SKEW' 'FROM' base_table_ref
    | 'DATA' opt_detailed order_by_clause opt_properties
    | 'DATA' opt_detailed 'FROM' table_name order_by_clause
    | opt_tmp 'PARTITIONS' 'FROM' table_name opt_wild_where order_by_clause limit_clause
    | 'PARTITION' INTEGER_VALUE
    | 'TABLET' INTEGER_VALUE
    | 'TABLETS' 'BELONG' integer_list
    | 'TABLETS' 'FROM' table_name opt_partition_names opt_wild_where order_by_clause limit_clause
    | 'PROPERTY' opt_user opt_wild_where
    | 'ALL' 'PROPERTIES' opt_wild_where
    | 'BACKUP' opt_db opt_wild_where
    | 'RESTORE' opt_db opt_wild_where
    | 'BRIEF' 'RESTORE' opt_db opt_wild_where
    | 'BROKER'
    | 'RESOURCES' opt_wild_where order_by_clause limit_clause
    | 'WORKLOAD' 'GROUPS' opt_wild_where
    | 'BACKENDS'
    | 'TRASH' 'ON' STRING_LITERAL
    | 'TRASH'
    | 'FRONTENDS'
    | 'FRONTENDS' ident
    | 'REPOSITORIES'
    | 'SNAPSHOT' 'ON' ident opt_wild_where
    | 'ALL' 'GRANTS'
    | 'GRANTS'
    | 'GRANTS' 'FOR' user_identity
    | 'ROLES'
    | 'PRIVILEGES'
    | opt_full opt_builtin 'FUNCTIONS' opt_db opt_wild_where
    | 'GLOBAL' opt_full 'FUNCTIONS' opt_wild_where
    | 'TYPE_CAST' opt_db
    | 'FILE' opt_db
    | keys_or_index from_or_in table_name opt_db
    | 'VIEW' from_or_in table_name opt_db
    | 'TRANSACTION' opt_db opt_wild_where
    | 'QUERY' 'PROFILE' opt_profile_path limit_clause
    | 'LOAD' 'PROFILE' opt_profile_path limit_clause
    | 'CACHE' 'HOTSPOT' STRING_LITERAL
    | 'ENCRYPTKEYS' opt_db opt_wild_where
    | 'TABLE' 'STATS' table_name opt_partition_names opt_col_list
    | 'TABLE' 'STATS' INTEGER_VALUE
    | 'INDEX' 'STATS' table_name ident
    | 'COLUMN' opt_cached 'STATS' table_name opt_col_list opt_partition_names
    | 'COLUMN' 'HISTOGRAM' table_name opt_col_list
    | 'TABLE' 'CREATION' opt_db opt_wild_where
    | 'LAST' 'INSERT'
    | 'CREATE' 'MATERIALIZED' 'VIEW' ident 'ON' table_name
    | 'ANALYZE' opt_table_name opt_wild_where
    | 'ANALYZE' INTEGER_VALUE opt_wild_where
    | 'AUTO' 'ANALYZE' opt_table_name opt_wild_where
    | 'AUTO' 'JOBS' opt_table_name opt_wild_where
    | 'ANALYZE' 'TASK' 'STATUS' INTEGER_VALUE
    | 'CATALOG' 'RECYCLE' 'BIN' opt_wild_where
    | 'QUERY' 'STATS'
    | 'QUERY' 'STATS' 'FOR' ident
    | 'QUERY' 'STATS' 'FROM' table_name
    | 'QUERY' 'STATS' 'FROM' table_name 'ALL'
    | 'QUERY' 'STATS' 'FROM' table_name 'ALL' 'VERBOSE'
    | 'BUILD' 'INDEX' opt_db opt_wild_where order_by_clause limit_clause
    | 'CLUSTERS'
    | 'COMPUTE' 'GROUPS'
    | 'CONVERT_LIGHT_SCHEMA_CHANGE_PROCESS' opt_db
    | 'REPLICA' 'STATUS' 'FROM' base_table_ref opt_wild_where
    | 'REPLICA' 'DISTRIBUTION' 'FROM' base_table_ref
    | 'FRONTEND' 'CONFIG' opt_wild_where
    | 'BACKEND' 'CONFIG' opt_wild_where
    | 'BACKEND' 'CONFIG' opt_wild_where 'FROM' INTEGER_VALUE
    | 'TABLET' 'STORAGE' 'FORMAT'
    | 'TABLET' 'STORAGE' 'FORMAT' 'VERBOSE'
    | 'TABLET' 'DIAGNOSIS' INTEGER_VALUE
    | 'COPY' opt_db opt_where order_by_clause limit_clause
    | 'WARM' 'UP' 'JOB' opt_wild_where
    ;

// sql_parser.cup:4915
opt_tmp
    : /* empty */
    | 'TEMPORARY'
    ;

// sql_parser.cup:4926
keys_or_index
    : 'KEY'
    | 'KEYS'
    | 'INDEX'
    | 'INDEXES'
    ;

// sql_parser.cup:4933
opt_db
    : /* empty */
    | from_or_in ident
    ;

// sql_parser.cup:4944
charset
    : 'CHAR' 'SET'
    | 'CHARSET'
    ;

// sql_parser.cup:4949
charset_name_or_default
    : ident_or_text
    | 'DEFAULT'
    ;

// sql_parser.cup:4960
old_or_new_charset_name_or_default
    : ident_or_text
    | 'DEFAULT'
    ;

// sql_parser.cup:4971
opt_collate
    : /* empty */
    | 'COLLATE' collation_name_or_default
    ;

// sql_parser.cup:4982
collation_name_or_default
    : ident_or_text
    | 'DEFAULT'
    ;

// sql_parser.cup:4993
opt_storage
    : /* empty */
    | 'STORAGE'
    ;

// sql_parser.cup:4998
procedure_or_function
    : 'PROCEDURE'
    | 'FUNCTION'
    ;

// sql_parser.cup:5003
from_or_in
    : 'FROM'
    | 'IN'
    ;

// sql_parser.cup:5008
opt_full
    : /* empty */
    | 'FULL'
    | 'EXTENDED'
    ;

// sql_parser.cup:5023
opt_wild_where
    : /* empty */
    | 'LIKE' STRING_LITERAL
    | 'WHERE' expr
    ;

// sql_parser.cup:5035
opt_where
    : /* empty */
    | 'WHERE' expr
    ;

// sql_parser.cup:5043
opt_id
    : /* empty */
    | INTEGER_VALUE
    ;

// sql_parser.cup:5054
opt_alter_type
    : 'ROLLUP'
    | 'MATERIALIZED' 'VIEW'
    | 'COLUMN'
    ;

// sql_parser.cup:5069
opt_builtin
    : /* empty */
    | 'BUILTIN'
    ;

// sql_parser.cup:5079
opt_explain_options
    : /* empty */
    | 'VERBOSE'
    | KW_TREE
    | 'GRAPH'
    ;

// sql_parser.cup:5098
describe_stmt
    : describe_command table_name opt_partition_names
    | describe_command 'FUNCTION' table_valued_function_ref
    | describe_command table_name 'ALL'
    | describe_command opt_explain_options query_stmt
    | describe_command opt_explain_options insert_stmt
    ;

// sql_parser.cup:5123
describe_command
    : 'DESCRIBE'
    | 'DESC'
    ;

// sql_parser.cup:5129
cancel_stmt
    : 'CANCEL' cancel_param
    ;

// sql_parser.cup:5136
job_id_list
    : /* empty */
    | '(' integer_list ')'
    ;

// sql_parser.cup:5146
cancel_param
    : 'LOAD' opt_db opt_wild_where
    | 'EXPORT' opt_db opt_wild_where
    | 'ALTER' 'TABLE' opt_alter_type 'FROM' table_name job_id_list
    | 'BUILD' 'INDEX' 'ON' table_name job_id_list
    | 'DECOMMISSION' 'BACKEND' string_list
    | 'BACKUP' opt_db
    | 'RESTORE' opt_db
    | 'WARM' 'UP' 'JOB' opt_wild_where
    ;

// sql_parser.cup:5181
opt_detailed
    : /* empty */
    | 'ALL'
    ;

// sql_parser.cup:5193
delete_stmt
    : 'DELETE' 'FROM' table_name opt_partition_names opt_table_alias opt_using_clause where_clause
    ;

// sql_parser.cup:5200
opt_using_clause
    : /* empty */
    | 'USING' table_ref_list
    ;

// sql_parser.cup:5229
query_stmt
    : opt_with_clause set_operand_list opt_outfile
    | opt_with_clause set_operation_with_order_by_or_limit opt_outfile
    ;

// sql_parser.cup:5250
opt_outfile
    : /* empty */
    | 'INTO' 'OUTFILE' STRING_LITERAL opt_file_format opt_properties
    ;

// sql_parser.cup:5260
opt_with_clause
    : 'WITH' with_view_def_list
    | /* empty */
    ;

// sql_parser.cup:5267
with_view_def
    : ident 'AS' '(' query_stmt ')'
    | STRING_LITERAL 'AS' '(' query_stmt ')'
    | ident '(' ident_list ')' 'AS' '(' query_stmt ')'
    | STRING_LITERAL '(' ident_list ')' 'AS' '(' query_stmt ')'
    ;

// sql_parser.cup:5280
with_view_def_list
    : with_view_def
    | with_view_def_list ',' with_view_def
    ;

// sql_parser.cup:5305
set_operation_with_order_by_or_limit
    : set_operand_list 'LIMIT' INTEGER_VALUE
    | set_operand_list 'LIMIT' INTEGER_VALUE ',' INTEGER_VALUE
    | set_operand_list 'LIMIT' INTEGER_VALUE 'OFFSET' INTEGER_VALUE
    | set_operand_list 'ORDER' 'BY' order_by_elements
    | set_operand_list 'ORDER' 'BY' order_by_elements 'LIMIT' INTEGER_VALUE
    | set_operand_list 'ORDER' 'BY' order_by_elements 'LIMIT' INTEGER_VALUE ',' INTEGER_VALUE
    | set_operand_list 'ORDER' 'BY' order_by_elements 'LIMIT' INTEGER_VALUE 'OFFSET' INTEGER_VALUE
    ;

// sql_parser.cup:5374
set_operand
    : select_stmt
    | '(' query_stmt ')'
    ;

// sql_parser.cup:5385
set_operand_list
    : set_operand
    | set_operand_list set_op opt_set_qualifier set_operand
    ;

// sql_parser.cup:5399
set_op
    : 'UNION'
    | 'INTERSECT'
    | 'EXCEPT'
    | 'MINUS'
    ;

// sql_parser.cup:5410
opt_set_qualifier
    : /* empty */
    | 'DISTINCT'
    | 'ALL'
    ;

// sql_parser.cup:5419
switch_stmt
    : 'SWITCH' ident
    ;

// sql_parser.cup:5427
use_stmt
    : 'USE' ident
    | 'USE' ident '.' ident
    ;

// sql_parser.cup:5438
insert_overwrite_stmt
    : 'INSERT' 'OVERWRITE' 'TABLE' insert_target opt_with_label opt_col_list opt_plan_hints insert_source
    ;

// sql_parser.cup:5446
use_cloud_cluster_stmt
    : 'USE' '@' ident
    | 'USE' ident '@' ident
    | 'USE' ident '.' ident '@' ident
    ;

// sql_parser.cup:5462
insert_stmt
    : 'INSERT' 'INTO' insert_target opt_with_label opt_col_list opt_plan_hints insert_source
    | 'INSERT' 'INTO' 'DORIS_INTERNAL_TABLE_ID' '(' INTEGER_VALUE ')' opt_with_label opt_col_list opt_plan_hints insert_source
    ;

// sql_parser.cup:5476
insert_target
    : table_name opt_partition_names
    ;

// sql_parser.cup:5483
opt_with_label
    : /* empty */
    | 'WITH' 'LABEL' ident
    ;

// sql_parser.cup:5494
insert_source
    : query_stmt
    ;

// sql_parser.cup:5503
update_stmt
    : 'UPDATE' table_name opt_table_alias set_clause opt_from_clause where_clause
    ;

// sql_parser.cup:5510
opt_from_clause
    : /* empty */
    | from_clause
    ;

// sql_parser.cup:5521
set_clause
    : 'SET' assignment_list
    ;

// sql_parser.cup:5528
assignment_list
    : assignment
    | assignment_list ',' assignment
    ;

// sql_parser.cup:5542
assignment
    : column_ref '=' expr
    ;

// sql_parser.cup:5550
backup_stmt
    : 'BACKUP' 'SNAPSHOT' job_label 'TO' ident opt_backup_table_ref_list opt_properties
    ;

// sql_parser.cup:5560
unlock_tables_stmt
    : 'UNLOCK' 'TABLES'
    ;

// sql_parser.cup:5566
lock_alias
    : /* empty */
    | 'AS' ident
    | 'AS' STRING_LITERAL
    ;

// sql_parser.cup:5581
lock_table
    : table_name lock_alias 'READ'
    | table_name lock_alias 'READ' 'LOCAL'
    | table_name lock_alias 'WRITE'
    | table_name lock_alias 'LOW_PRIORITY' 'WRITE'
    ;

// sql_parser.cup:5603
opt_lock_tables_list
    : lock_table
    | opt_lock_tables_list ',' lock_table
    | /* empty */
    ;

// sql_parser.cup:5621
lock_tables_stmt
    : 'LOCK' 'TABLES' opt_lock_tables_list
    ;

// sql_parser.cup:5628
opt_backup_table_ref_list
    : backup_exclude_or_not '(' base_table_ref_list ')'
    | /* empty */
    ;

// sql_parser.cup:5639
backup_exclude_or_not
    : 'ON'
    | 'EXCLUDE'
    ;

// sql_parser.cup:5651
restore_stmt
    : 'RESTORE' 'SNAPSHOT' job_label 'FROM' ident opt_backup_table_ref_list opt_properties
    ;

// sql_parser.cup:5662
kill_stmt
    : 'KILL' INTEGER_VALUE
    | 'KILL' 'CONNECTION' INTEGER_VALUE
    | 'KILL' 'QUERY' INTEGER_VALUE
    | 'KILL' 'QUERY' STRING_LITERAL
    ;

// sql_parser.cup:5681
kill_analysis_job_stmt
    : 'KILL' 'ANALYZE' INTEGER_VALUE
    ;

// sql_parser.cup:5690
set_stmt
    : 'SET' start_option_value_list
    | 'SET' 'PROPERTY' opt_user user_property_list
    | 'SET' ident 'AS' 'DEFAULT' 'STORAGE' 'VAULT'
    ;

// sql_parser.cup:5706
unset_var_stmt
    : 'UNSET' opt_var_type 'VARIABLE' variable_name
    | 'UNSET' opt_var_type 'VARIABLE' 'ALL'
    ;

// sql_parser.cup:5717
unset_default_storage_vault_stmt
    : 'UNSET' 'DEFAULT' 'STORAGE' 'VAULT'
    ;

// sql_parser.cup:5723
user_property_list
    : user_property
    | user_property_list ',' user_property
    ;

// sql_parser.cup:5735
user_property
    : STRING_LITERAL equal STRING_LITERAL
    | STRING_LITERAL equal 'NULL'
    ;

// sql_parser.cup:5747
start_option_value_list
    : option_value_no_option_type option_value_list_continued
    | 'TRANSACTION' transaction_characteristics
    | option_type start_option_value_list_following_option_type
    ;

// sql_parser.cup:5774
start_option_value_list_following_option_type
    : option_value_follow_option_type option_value_list_continued
    | 'TRANSACTION' transaction_characteristics
    ;

// sql_parser.cup:5787
option_value_list_continued
    : /* empty */
    | ',' option_value_list
    ;

// sql_parser.cup:5798
option_value_list
    : option_value
    | option_value_list ',' option_value
    ;

// sql_parser.cup:5810
option_value
    : option_type option_value_follow_option_type
    | option_value_no_option_type
    ;

// sql_parser.cup:5822
option_value_follow_option_type
    : variable_name equal set_expr_or_default
    ;

// sql_parser.cup:5829
option_value_no_option_type
    : variable_name equal set_expr_or_default
    | '@' ident_or_text equal expr
    | '@' '@' variable_name equal set_expr_or_default
    | '@' '@' var_ident_type variable_name equal set_expr_or_default
    | charset old_or_new_charset_name_or_default
    | 'NAMES' equal expr
    | 'NAMES' charset_name_or_default opt_collate
    | 'PASSWORD' equal text_or_password
    | 'PASSWORD' 'FOR' user_identity equal text_or_password
    | 'LDAP_ADMIN_PASSWORD' equal text_or_password
    ;

// sql_parser.cup:5876
variable_name
    : ident
    ;

// sql_parser.cup:5883
text_or_password
    : STRING_LITERAL
    | 'PASSWORD' '(' STRING_LITERAL ')'
    ;

// sql_parser.cup:5896
option_type
    : 'GLOBAL'
    | 'LOCAL'
    | 'SESSION'
    ;

// sql_parser.cup:5911
opt_var_type
    : /* empty */
    | 'GLOBAL'
    | 'LOCAL'
    | 'SESSION'
    ;

// sql_parser.cup:5922
var_ident_type
    : 'GLOBAL' '.'
    | 'LOCAL' '.'
    | 'SESSION' '.'
    ;

// sql_parser.cup:5937
equal
    : '='
    | ':='
    ;

// sql_parser.cup:5942
transaction_characteristics
    : transaction_access_mode
    | isolation_level
    | transaction_access_mode ',' isolation_level
    | isolation_level ',' transaction_access_mode
    ;

// sql_parser.cup:5949
transaction_access_mode
    : 'READ' 'ONLY'
    | 'READ' 'WRITE'
    ;

// sql_parser.cup:5954
isolation_level
    : 'ISOLATION' 'LEVEL' isolation_types
    ;

// sql_parser.cup:5958
isolation_types
    : 'READ' 'UNCOMMITTED'
    | 'READ' 'COMMITTED'
    | 'REPEATABLE' 'READ'
    | 'SERIALIZABLE'
    ;

// sql_parser.cup:5965
set_expr_or_default
    : 'DEFAULT'
    | 'ON'
    | 'ALL'
    | expr
    ;

// sql_parser.cup:5984
select_stmt
    : select_clause limit_clause
    | select_clause from_clause where_clause group_by_clause having_clause order_by_clause limit_clause
    | value_clause order_by_clause limit_clause
    ;

// sql_parser.cup:6006
value_clause
    : 'VALUES' row_value
    | value_clause ',' row_value
    ;

// sql_parser.cup:6018
row_value
    : '(' opt_values ')'
    ;

// sql_parser.cup:6025
opt_values
    : values
    | /* empty */
    ;

// sql_parser.cup:6036
values
    : expr_or_default
    | values ',' expr_or_default
    ;

// sql_parser.cup:6048
expr_or_default
    : expr
    | 'DEFAULT'
    ;

// sql_parser.cup:6081
select_clause
    : 'SELECT' opt_select_hints select_list
    | 'SELECT' opt_select_hints 'ALL' select_list
    | 'SELECT' opt_select_hints 'DISTINCT' select_list
    ;

// sql_parser.cup:6100
query_hint_parameter_key
    : literal
    | variable_name
    ;

// sql_parser.cup:6111
query_hint_parameter
    : query_hint_parameter_key
    | query_hint_parameter_key equal literal_or_ident
    ;

// sql_parser.cup:6123
query_hint_parameters
    : query_hint_parameters ',' query_hint_parameter
    | query_hint_parameter
    | /* empty */
    ;

// sql_parser.cup:6141
query_hint
    : ident '(' query_hint_parameters ')'
    ;

// sql_parser.cup:6148
query_hints
    : query_hints query_hint
    | query_hint
    ;

// sql_parser.cup:6163
literal_or_ident
    : literal
    | ident
    ;

// sql_parser.cup:6174
opt_select_hints
    : COMMENTED_PLAN_HINT_START query_hints COMMENTED_PLAN_HINT_END
    | /* empty */
    ;

// sql_parser.cup:6185
select_list
    : select_sublist
    | '*'
    | '*' 'EXCEPT' '(' select_sublist ')'
    ;

// sql_parser.cup:6204
select_sublist
    : select_sublist ',' select_list_item
    | select_sublist ',' '*'
    | '*' ',' select_list_item
    | select_list_item
    ;

// sql_parser.cup:6231
select_list_item
    : expr select_alias
    | star_expr
    ;

// sql_parser.cup:6242
select_alias
    : /* empty */
    | 'AS' ident
    | ident
    | 'AS' STRING_LITERAL
    | STRING_LITERAL
    ;

// sql_parser.cup:6265
star_expr
    : ident '.' '*'
    | ident '.' ident '.' '*'
    | ident '.' ident '.' ident '.' '*'
    ;

// sql_parser.cup:6282
opt_table_name
    : /* empty */
    | table_name
    ;

// sql_parser.cup:6292
table_name
    : ident
    | ident '.' ident
    | ident '.' ident '.' ident
    ;

// sql_parser.cup:6301
db_name
    : ident
    | ident '.' ident
    ;

// sql_parser.cup:6308
colocate_group_name
    : ident
    | ident '.' ident
    ;

// sql_parser.cup:6319
opt_scan_params_ref
    : /* empty */
    | '@' ident '(' opt_key_value_map_in_paren ')'
    ;

// sql_parser.cup:6330
encryptkey_name
    : ident
    | ident '.' ident
    ;

// sql_parser.cup:6341
function_name
    : type_function_name
    | ident '.' type_function_name
    ;

// sql_parser.cup:6348
type_function_name
    : ident
    | type_func_name_keyword
    ;

// sql_parser.cup:6355
from_clause
    : 'FROM' table_ref_list
    | 'FROM' 'DUAL'
    ;

// sql_parser.cup:6362
table_ref_list
    : table_ref opt_sort_hints
    | table_ref_list ',' table_ref opt_sort_hints
    | table_ref_list join_operator opt_plan_hints table_ref opt_sort_hints
    | table_ref_list join_operator opt_plan_hints table_ref opt_sort_hints 'ON' expr
    | table_ref_list join_operator opt_plan_hints table_ref opt_sort_hints 'USING' '(' ident_list ')'
    ;

// sql_parser.cup:6406
table_ref
    : base_table_ref opt_lateral_view_ref_list
    | inline_view_ref opt_lateral_view_ref_list
    | table_valued_function_ref
    ;

// sql_parser.cup:6423
table_valued_function_ref
    : ident '(' opt_key_value_map_in_paren ')' opt_table_alias
    ;

// sql_parser.cup:6430
inline_view_ref
    : '(' query_stmt ')' opt_table_alias
    ;

// sql_parser.cup:6437
base_table_ref_list
    : base_table_ref
    | base_table_ref_list ',' base_table_ref
    ;

// sql_parser.cup:6451
base_table_ref
    : table_name opt_scan_params_ref opt_table_snapshot opt_partition_names opt_tablet_list opt_table_alias opt_table_sample opt_common_hints
    ;

// sql_parser.cup:6458
opt_table_snapshot
    : /* empty */
    | table_snapshot
    ;

// sql_parser.cup:6469
table_snapshot
    : 'FOR' 'VERSION' 'AS' 'OF' INTEGER_VALUE
    | 'FOR' 'TIME' 'AS' 'OF' STRING_LITERAL
    ;

// sql_parser.cup:6480
opt_common_hints
    : COMMENTED_PLAN_HINT_START ident_list COMMENTED_PLAN_HINT_END
    | '[' ident_list ']'
    | /* empty */
    ;

// sql_parser.cup:6495
opt_alias
    : /* empty */
    | 'AS' ident
    ;

// sql_parser.cup:6506
opt_table_alias
    : /* empty */
    | ident
    | 'AS' ident
    | '=' ident
    ;

// sql_parser.cup:6525
opt_partition_names
    : /* empty */
    | partition_names
    ;

// sql_parser.cup:6536
opt_tablet_list
    : /* empty */
    | tablet_list
    ;

// sql_parser.cup:6547
tablet_list
    : 'TABLET' '(' integer_list ')'
    ;

// sql_parser.cup:6554
partition_names
    : 'PARTITION' '(' ident_list ')'
    | 'TEMPORARY' 'PARTITION' '(' ident_list ')'
    | 'PARTITIONS' '(' ident_list ')'
    | 'TEMPORARY' 'PARTITIONS' '(' ident_list ')'
    | 'PARTITION' ident
    | 'TEMPORARY' 'PARTITION' ident
    | 'PARTITIONS' '(' '*' ')'
    | 'PARTITION' '(' '*' ')'
    | 'PARTITIONS' 'WITH' 'RECENT' INTEGER_VALUE
    ;

// sql_parser.cup:6593
opt_table_sample
    : /* empty */
    | table_sample
    ;

// sql_parser.cup:6604
table_sample
    : 'TABLESAMPLE' '(' INTEGER_VALUE 'PERCENT' ')'
    | 'TABLESAMPLE' '(' INTEGER_VALUE 'ROWS' ')'
    | 'TABLESAMPLE' '(' INTEGER_VALUE 'PERCENT' ')' 'REPEATABLE' INTEGER_VALUE
    | 'TABLESAMPLE' '(' INTEGER_VALUE 'ROWS' ')' 'REPEATABLE' INTEGER_VALUE
    ;

// sql_parser.cup:6623
opt_lateral_view_ref_list
    : /* empty */
    | lateral_view_ref_list
    ;

// sql_parser.cup:6634
lateral_view_ref_list
    : lateral_view_ref
    | lateral_view_ref lateral_view_ref_list
    ;

// sql_parser.cup:6648
lateral_view_ref
    : 'LATERAL' 'VIEW' function_call_expr ident 'AS' ident
    ;

// sql_parser.cup:6655
join_operator
    : opt_inner 'JOIN'
    | 'LEFT' opt_outer 'JOIN'
    | 'RIGHT' opt_outer 'JOIN'
    | 'FULL' opt_outer 'JOIN'
    | 'LEFT' 'SEMI' 'JOIN'
    | 'RIGHT' 'SEMI' 'JOIN'
    | 'LEFT' 'ANTI' 'JOIN'
    | 'RIGHT' 'ANTI' 'JOIN'
    | 'CROSS' 'JOIN'
    | 'NATURAL' 'JOIN'
    ;

// sql_parser.cup:6682
opt_inner
    : 'INNER'
    | /* empty */
    ;

// sql_parser.cup:6687
opt_outer
    : 'OUTER'
    | /* empty */
    ;

// sql_parser.cup:6692
opt_plan_hints
    : COMMENTED_PLAN_HINTS
    | COMMENTED_PLAN_HINT_START ident_list COMMENTED_PLAN_HINT_END
    | '[' ident_list ']'
    | /* empty */
    ;

// sql_parser.cup:6719
opt_sort_hints
    : '[' ident_list ']'
    | /* empty */
    ;

// sql_parser.cup:6726
ident_list
    : ident
    | ident_list ',' ident
    ;

// sql_parser.cup:6740
opt_ident_list
    : /* empty */
    | ident_list
    ;

// sql_parser.cup:6750
with_analysis_properties
    : 'SYNC'
    | 'INCREMENTAL'
    | 'SAMPLE' 'PERCENT' INTEGER_VALUE
    | 'SAMPLE' 'ROWS' INTEGER_VALUE
    | 'BUCKETS' INTEGER_VALUE
    | 'PERIOD' INTEGER_VALUE
    | 'HISTOGRAM'
    | 'CRON' STRING_LITERAL
    | 'FULL'
    | 'SQL'
    ;

// sql_parser.cup:6814
opt_with_analysis_properties
    : /* empty */
    | opt_with_analysis_properties 'WITH' with_analysis_properties
    ;

// sql_parser.cup:6828
expr_list
    : expr
    | expr_list ',' expr
    ;

// sql_parser.cup:6842
where_clause
    : /* empty */
    | 'WHERE' expr
    ;

// sql_parser.cup:6849
delete_on_clause
    : /* empty */
    | 'DELETE' 'ON' expr
    ;

// sql_parser.cup:6856
sequence_col_clause
    : /* empty */
    | 'ORDER' 'BY' ident
    ;

// sql_parser.cup:6863
pre_filter_clause
    : /* empty */
    | 'PRECEDING' 'FILTER' expr
    ;

// sql_parser.cup:6870
grouping_set
    : '(' ')'
    | '(' expr_list ')'
    ;

// sql_parser.cup:6880
grouping_set_list
    : grouping_set
    | grouping_set_list ',' grouping_set
    ;

// sql_parser.cup:6894
grouping_elements
    : expr_list
    | 'GROUPING' 'SETS' '(' grouping_set_list ')'
    | 'CUBE' '(' expr_list ')'
    | 'ROLLUP' '(' expr_list ')'
    ;

// sql_parser.cup:6913
group_by_clause
    : 'GROUP' 'BY' grouping_elements
    | /* empty */
    ;

// sql_parser.cup:6920
having_clause
    : 'HAVING' expr
    | /* empty */
    ;

// sql_parser.cup:6927
order_by_clause
    : 'ORDER' 'BY' order_by_elements
    | /* empty */
    ;

// sql_parser.cup:6934
order_by_elements
    : order_by_element
    | order_by_elements ',' order_by_element
    ;

// sql_parser.cup:6952
order_by_element
    : expr opt_order_param opt_nulls_order_param
    ;

// sql_parser.cup:6963
opt_order_param
    : 'ASC'
    | 'DESC'
    | /* empty */
    ;

// sql_parser.cup:6972
opt_nulls_order_param
    : 'NULLS' 'FIRST'
    | 'NULLS' 'LAST'
    | /* empty */
    ;

// sql_parser.cup:6981
limit_clause
    : 'LIMIT' INTEGER_VALUE
    | /* empty */
    | 'LIMIT' INTEGER_VALUE ',' INTEGER_VALUE
    | 'LIMIT' INTEGER_VALUE 'OFFSET' INTEGER_VALUE
    ;

// sql_parser.cup:6992
type
    : 'TINYINT' opt_field_length
    | 'SMALLINT' opt_field_length
    | opt_signed_unsigned 'INT' opt_field_length
    | 'UNSIGNED' 'INT' opt_field_length
    | 'BIGINT' opt_field_length
    | 'LARGEINT' opt_field_length
    | 'BOOLEAN'
    | 'FLOAT'
    | 'DOUBLE'
    | 'DATE'
    | 'DATETIME' '(' INTEGER_VALUE ')'
    | 'DATETIME'
    | 'DATEV1'
    | 'DATETIMEV1'
    | 'TIME' '(' INTEGER_VALUE ')'
    | 'TIME'
    | 'DATEV2'
    | 'DATETIMEV2' '(' INTEGER_VALUE ')'
    | 'DATETIMEV2'
    | 'IPV4'
    | 'IPV6'
    | 'BITMAP'
    | 'QUANTILE_STATE'
    | 'STRING'
    | 'JSON'
    | 'JSONB'
    | 'TEXT'
    | 'VARIANT'
    | 'VARCHAR' '(' INTEGER_VALUE ')'
    | 'VARCHAR' '(' ident_or_text ')'
    | 'VARCHAR' '(' '*' ')'
    | 'VARCHAR'
    | 'ARRAY' '<' type '>'
    | 'MAP' '<' type ',' type '>'
    | 'STRUCT' '<' struct_field_list '>'
    | 'VARIANT' '<' variant_field_list type_opt_properties '>'
    | 'VARIANT' '<' properties '>'
    | 'CHAR' '(' INTEGER_VALUE ')'
    | 'CHAR' '(' ident_or_text ')'
    | 'CHAR'
    | 'DECIMAL' '(' INTEGER_VALUE ')'
    | 'DECIMAL' '(' INTEGER_VALUE ',' INTEGER_VALUE ')'
    | 'DECIMAL'
    | 'DECIMAL' '(' ident_or_text ')'
    | 'DECIMAL' '(' ident_or_text ',' ident_or_text ')'
    | 'DECIMALV3' '(' INTEGER_VALUE ')'
    | 'DECIMALV3' '(' INTEGER_VALUE ',' INTEGER_VALUE ')'
    | 'DECIMALV3'
    | 'DECIMALV3' '(' ident_or_text ')'
    | 'DECIMALV3' '(' ident_or_text ',' ident_or_text ')'
    | 'DECIMALV2' '(' INTEGER_VALUE ')'
    | 'DECIMALV2' '(' INTEGER_VALUE ',' INTEGER_VALUE ')'
    | 'DECIMALV2'
    | 'DECIMALV2' '(' ident_or_text ')'
    | 'DECIMALV2' '(' ident_or_text ',' ident_or_text ')'
    | 'HLL'
    | 'AGG_STATE' '<' IDENTIFIER '(' type_def_nullable_list ')' '>'
    | 'ALL'
    ;

// sql_parser.cup:7126
opt_field_length
    : '(' INTEGER_VALUE ')'
    | /* empty */
    ;

// sql_parser.cup:7135
opt_signed_unsigned
    : /* empty */
    | 'SIGNED'
    ;

// sql_parser.cup:7142
type_def
    : type
    ;

// sql_parser.cup:7147
type_def_list
    : type_def
    | type_def_list ',' type_def
    ;

// sql_parser.cup:7159
type_def_nullable
    : type opt_is_allow_null
    ;

// sql_parser.cup:7164
type_def_nullable_list
    : type_def_nullable
    | type_def_nullable_list ',' type_def_nullable
    ;

// sql_parser.cup:7177
func_args_def
    : /* empty */
    | type_def_list
    | '...'
    | type_def_list ',' '...'
    ;

// sql_parser.cup:7193
cast_expr
    : 'CAST' '(' expr 'AS' type_def ')'
    ;

// sql_parser.cup:7214
case_expr
    : 'CASE' expr case_when_clause_list case_else_clause 'END'
    | 'CASE' case_when_clause_list case_else_clause 'END'
    ;

// sql_parser.cup:7227
case_when_clause_list
    : 'WHEN' expr 'THEN' expr
    | case_when_clause_list 'WHEN' expr 'THEN' expr
    ;

// sql_parser.cup:7242
case_else_clause
    : 'ELSE' expr
    | /* empty */
    ;

// sql_parser.cup:7249
sign_chain_expr
    : '-' expr
    | '+' expr
    ;

// sql_parser.cup:7264
expr
    : non_pred_expr opt_collate
    | predicate
    ;

// sql_parser.cup:7271
function_call_expr
    : function_name '(' ')'
    | 'ADD' '(' function_params ')'
    | function_name '(' function_params ')'
    | 'LIKE' '(' function_params ')'
    | 'REGEXP' '(' function_params ')'
    | function_name '(' ident '->' expr ',' function_params ')'
    | function_name '(' '(' ident ',' ident_list ')' '->' expr ',' function_params ')'
    ;

// sql_parser.cup:7310
array_literal
    : '[' ']'
    | '[' expr_list ']'
    ;

// sql_parser.cup:7321
array_expr
    : 'ARRAY' '(' function_params ')'
    | 'ARRAY' '(' ')'
    ;

// sql_parser.cup:7332
kv_list
    : expr ':' expr
    | kv_list ',' expr ':' expr
    ;

// sql_parser.cup:7348
map_literal
    : '{' '}'
    | '{' kv_list '}'
    ;

// sql_parser.cup:7359
map_expr
    : 'MAP' '(' function_params ')'
    | 'MAP' '(' ')'
    ;

// sql_parser.cup:7370
struct_field
    : ident ':' type opt_comment
    ;

// sql_parser.cup:7375
struct_field_list
    : struct_field
    | struct_field_list ',' struct_field
    ;

// sql_parser.cup:7387
variant_field
    : STRING_LITERAL ':' type opt_comment
    ;

// sql_parser.cup:7392
variant_field_list
    : variant_field
    | variant_field_list ',' variant_field
    ;

// sql_parser.cup:7404
struct_literal
    : '{' expr_list '}'
    ;

// sql_parser.cup:7411
exists_predicate
    : 'EXISTS' subquery
    ;

// sql_parser.cup:7416
non_pred_expr
    : sign_chain_expr
    | '@' ident
    | '@' '@' ident
    | '@' '@' var_ident_type ident
    | literal
    | array_expr
    | array_literal
    | map_expr
    | map_literal
    | struct_literal
    | function_call_expr
    | 'DATE' STRING_LITERAL
    | 'DATEV1' STRING_LITERAL
    | 'DATEV2' STRING_LITERAL
    | 'TIMESTAMP' STRING_LITERAL
    | 'EXTRACT' '(' function_name 'FROM' func_arg_list ')'
    | function_name '(' function_params order_by_clause ')'
    | analytic_expr
    | 'IF' '(' expr_list ')'
    | expr_pipe_list
    | cast_expr
    | case_expr
    | column_ref
    | column_subscript
    | column_slice
    | timestamp_arithmetic_expr
    | arithmetic_expr
    | '(' non_pred_expr ')'
    | 'TRIM' '(' function_params ')'
    | 'DATABASE' '(' ')'
    | 'SCHEMA' '(' ')'
    | 'USER' '(' ')'
    | 'CURRENT_USER' '(' ')'
    | 'CURRENT_CATALOG' '(' ')'
    | 'CONNECTION_ID' '(' ')'
    | 'PASSWORD' '(' STRING_LITERAL ')'
    | subquery
    | 'NULL' 'IS' 'NULL'
    | 'NULL' 'IS' 'NOT' 'NULL'
    | 'CONVERT' '(' expr ',' type_def ')'
    | 'KEY' encryptkey_name
    | expr 'IS' 'TRUE'
    | expr 'IS' 'NOT' 'TRUE'
    | expr 'IS' 'FALSE'
    | expr 'IS' 'NOT' 'FALSE'
    | 'CONVERT' '(' expr 'USING' ident ')'
    | 'CHAR' '(' expr_list opt_using_charset ')'
    ;

// sql_parser.cup:7545
opt_using_charset
    : /* empty */
    | 'USING' ident
    ;

// sql_parser.cup:7556
expr_pipe_list
    : expr '||' expr
    | expr_pipe_list '||' expr
    ;

// sql_parser.cup:7571
func_arg_list
    : expr
    | func_arg_list ',' expr
    ;

// sql_parser.cup:7585
analytic_expr
    : function_call_expr 'OVER' '(' opt_partition_by_clause order_by_clause opt_window_clause ')'
    ;

// sql_parser.cup:7599
opt_partition_by_clause
    : 'PARTITION' 'BY' expr_list
    | /* empty */
    ;

// sql_parser.cup:7606
opt_window_clause
    : window_type window_boundary
    | window_type 'BETWEEN' window_boundary 'AND' window_boundary
    | /* empty */
    ;

// sql_parser.cup:7615
window_type
    : 'ROWS'
    | 'RANGE'
    ;

// sql_parser.cup:7622
window_boundary
    : 'UNBOUNDED' 'PRECEDING'
    | 'UNBOUNDED' 'FOLLOWING'
    | 'CURRENT' 'ROW'
    | expr 'PRECEDING'
    | expr 'FOLLOWING'
    ;

// sql_parser.cup:7643
arithmetic_expr
    : expr '*' expr
    | expr '/' expr
    | expr '%' expr
    | expr 'DIV' expr
    | expr '+' expr
    | expr '-' expr
    | expr '&' expr
    | expr '|' expr
    | expr '^' expr
    | '~' expr
    ;

// sql_parser.cup:7669
timestamp_arithmetic_expr
    : 'INTERVAL' expr ident '+' expr
    | expr '+' 'INTERVAL' expr ident
    | expr '-' 'INTERVAL' expr ident
    | function_name '(' expr_list ',' 'INTERVAL' expr ident ')'
    | function_name '(' time_unit ',' expr ',' expr ')'
    ;

// sql_parser.cup:7716
literal
    : INTEGER_VALUE
    | LARGE_INTEGER_LITERAL
    | FLOATINGPOINT_LITERAL
    | DECIMAL_VALUE
    | STRING_LITERAL
    | 'TRUE'
    | 'FALSE'
    | 'NULL'
    | '?'
    | '%'
    | UNMATCHED_STRING_LITERAL expr
    | NUMERIC_OVERFLOW
    ;

// sql_parser.cup:7756
function_params
    : '*'
    | 'ALL' '*'
    | expr_list
    | 'ALL' expr_list
    | 'DISTINCT' expr_list
    ;

// sql_parser.cup:7769
predicate
    : expr 'IS' 'NULL'
    | 'ISNULL' '(' expr ')'
    | expr 'IS' 'NOT' 'NULL'
    | between_predicate
    | comparison_predicate
    | compound_predicate
    | in_predicate
    | exists_predicate
    | like_predicate
    | match_predicate
    | '(' predicate ')'
    ;

// sql_parser.cup:7797
comparison_predicate
    : expr '=' expr
    | expr '!' '=' expr
    | expr '<' '>' expr
    | expr '<' '=' expr
    | expr '>' '=' expr
    | expr '<' expr
    | expr '>' expr
    | expr '<' '=' '>' expr
    ;

// sql_parser.cup:7816
like_predicate
    : expr 'LIKE' expr
    | expr 'REGEXP' expr
    | expr 'NOT' 'LIKE' expr
    | expr 'NOT' 'REGEXP' expr
    ;

// sql_parser.cup:7829
match_predicate
    : expr 'MATCH_ANY' expr
    | expr 'MATCH' expr
    | expr 'MATCH_ALL' expr
    | expr 'MATCH_PHRASE' expr
    | expr 'MATCH_PHRASE_PREFIX' expr
    | expr 'MATCH_REGEXP' expr
    | expr 'MATCH_PHRASE_EDGE' expr
    ;

// sql_parser.cup:7848
between_predicate
    : expr 'BETWEEN' non_pred_expr 'AND' expr
    | expr 'BETWEEN' predicate 'AND' expr
    | expr 'NOT' 'BETWEEN' non_pred_expr 'AND' expr
    | expr 'NOT' 'BETWEEN' predicate 'AND' expr
    ;

// sql_parser.cup:7859
in_predicate
    : expr 'IN' '(' expr_list ')'
    | expr 'NOT' 'IN' '(' expr_list ')'
    | expr 'IN' subquery
    | expr 'NOT' 'IN' subquery
    ;

// sql_parser.cup:7870
subquery
    : '(' subquery ')'
    | '(' query_stmt ')'
    ;

// sql_parser.cup:7877
compound_predicate
    : expr 'AND' expr
    | expr 'OR' expr
    | 'NOT' expr
    | '!' expr
    ;

// sql_parser.cup:7888
sub_column_path
    : ident
    | sub_column_path '.' ident
    ;

// sql_parser.cup:7900
column_ref
    : ident
    | ident '.' ident
    | ident '.' ident '.' ident
    | ident '.' ident '.' ident '.' ident
    | ident ':' sub_column_path
    | ident '.' ident ':' sub_column_path
    | ident '.' ident '.' ident ':' sub_column_path
    | ident '.' ident '.' ident '.' ident ':' sub_column_path
    ;

// sql_parser.cup:7936
column_subscript
    : expr '[' expr ']'
    ;

// sql_parser.cup:7945
column_slice
    : expr '[' expr ':' ']'
    | expr '[' expr ':' expr ']'
    ;

// sql_parser.cup:7961
privilege_type
    : ident opt_col_list
    | 'ALL'
    ;

// sql_parser.cup:7976
privilege_list
    : privilege_list ',' privilege_type
    | privilege_type
    ;

// sql_parser.cup:7988
string_list
    : string_list ',' STRING_LITERAL
    | STRING_LITERAL
    ;

// sql_parser.cup:8000
integer_list
    : integer_list ',' INTEGER_VALUE
    | INTEGER_VALUE
    ;

// sql_parser.cup:8012
admin_stmt
    : 'ADMIN' 'SHOW' 'REPLICA' 'STATUS' 'FROM' base_table_ref opt_wild_where
    | 'ADMIN' 'SHOW' 'REPLICA' 'DISTRIBUTION' 'FROM' base_table_ref
    | 'ADMIN' 'SET' 'REPLICA' 'STATUS' 'PROPERTIES' '(' key_value_map ')'
    | 'ADMIN' 'SET' 'REPLICA' 'VERSION' 'PROPERTIES' '(' key_value_map ')'
    | 'ADMIN' 'REPAIR' 'TABLE' base_table_ref
    | 'ADMIN' 'CANCEL' 'REPAIR' 'TABLE' base_table_ref
    | 'ADMIN' 'COMPACT' 'TABLE' base_table_ref opt_wild_where
    | 'ADMIN' 'SET' 'FRONTEND' 'CONFIG' opt_key_value_map
    | 'ADMIN' 'SET' 'ALL' 'FRONTENDS' 'CONFIG' opt_key_value_map
    | 'ADMIN' 'SET' 'FRONTEND' 'CONFIG' opt_key_value_map 'ALL'
    | 'ADMIN' 'SHOW' 'FRONTEND' 'CONFIG' opt_wild_where
    | 'ADMIN' 'SHOW' 'BACKEND' 'CONFIG' opt_wild_where
    | 'ADMIN' 'SHOW' 'BACKEND' 'CONFIG' opt_wild_where 'FROM' INTEGER_VALUE
    | 'ADMIN' 'CHECK' 'TABLET' '(' integer_list ')' opt_properties
    | 'ADMIN' 'REBALANCE' 'DISK' 'ON' '(' string_list ')'
    | 'ADMIN' 'REBALANCE' 'DISK'
    | 'ADMIN' 'CANCEL' 'REBALANCE' 'DISK' 'ON' '(' string_list ')'
    | 'ADMIN' 'CANCEL' 'REBALANCE' 'DISK'
    | 'ADMIN' 'CLEAN' 'TRASH' 'ON' '(' string_list ')'
    | 'ADMIN' 'CLEAN' 'TRASH'
    | 'ADMIN' 'SET' 'TABLE' table_name 'PARTITION' 'VERSION' opt_properties
    | 'ADMIN' 'DIAGNOSE' 'TABLET' INTEGER_VALUE
    | 'ADMIN' 'SHOW' 'TABLET' 'STORAGE' 'FORMAT'
    | 'ADMIN' 'SHOW' 'TABLET' 'STORAGE' 'FORMAT' 'VERBOSE'
    | 'ADMIN' 'COPY' 'TABLET' INTEGER_VALUE opt_properties
    | 'ADMIN' 'SET' 'TABLE' table_name 'STATUS' opt_properties
    ;

// sql_parser.cup:8125
truncate_stmt
    : 'TRUNCATE' 'TABLE' base_table_ref
    ;

// sql_parser.cup:8132
transaction_stmt
    : 'BEGIN'
    | 'BEGIN' 'WITH' 'LABEL' transaction_label
    | 'COMMIT' opt_work opt_chain opt_release
    | 'ROLLBACK' opt_work opt_chain opt_release
    ;

// sql_parser.cup:8151
transaction_label
    : /* empty */
    | ident
    ;

// sql_parser.cup:8162
unsupported_stmt
    : 'START' 'TRANSACTION' opt_with_consistent_snapshot
    ;

// sql_parser.cup:8169
opt_with_consistent_snapshot
    : /* empty */
    | 'WITH' 'CONSISTENT' 'SNAPSHOT'
    ;

// sql_parser.cup:8179
opt_work
    : /* empty */
    | 'WORK'
    ;

// sql_parser.cup:8188
opt_generated_always
    : /* empty */
    | 'GENERATED' 'ALWAYS'
    ;

// sql_parser.cup:8198
opt_chain
    : /* empty */
    | 'AND' 'NO' 'CHAIN'
    | 'AND' 'CHAIN'
    ;

// sql_parser.cup:8212
opt_release
    : /* empty */
    | 'RELEASE'
    | 'NO' 'RELEASE'
    ;

// sql_parser.cup:8226
type_func_name_keyword
    : 'LEFT'
    | 'RIGHT'
    ;

// sql_parser.cup:8234
keyword
    : 'AFTER'
    | 'ALWAYS'
    | 'AGGREGATE'
    | 'ALIAS'
    | 'AUTHORS'
    | 'ARRAY'
    | 'AUTO_INCREMENT'
    | 'BACKUP'
    | 'BEGIN'
    | 'BIN'
    | 'BITMAP'
    | 'BITMAP_EMPTY'
    | 'QUANTILE_STATE'
    | 'AGG_STATE'
    | 'BITMAP_UNION'
    | 'NGRAM_BF'
    | 'QUANTILE_UNION'
    | 'BLOB'
    | 'BOOLEAN'
    | 'BRANCH'
    | 'BRIEF'
    | 'BROKER'
    | 'S3'
    | 'HDFS'
    | 'BACKENDS'
    | 'BUILTIN'
    | 'BUILD'
    | 'CACHED'
    | 'CHAIN'
    | 'CHAR'
    | 'CHARSET'
    | 'CHECK'
    | 'COLUMNS'
    | 'COMMENT'
    | 'COMMITTED'
    | 'COMPLETE'
    | 'CONSISTENT'
    | 'COLLATION'
    | 'COMMIT'
    | 'COMPACT'
    | 'CONFIG'
    | 'CONNECTION'
    | 'CONNECTION_ID'
    | 'CONVERT'
    | 'COPY'
    | 'CREATION'
    | 'DATA'
    | 'DATE'
    | 'DATETIME'
    | 'DATEV2'
    | 'DATETIMEV2'
    | 'DATEV1'
    | 'DATETIMEV1'
    | 'DAYS'
    | 'DECIMAL'
    | 'DEFERRED'
    | 'DEMAND'
    | 'DECIMALV2'
    | 'DECIMALV3'
    | 'DIAGNOSE'
    | 'DIAGNOSIS'
    | 'DISTINCTPC'
    | 'DISTINCTPCSA'
    | 'DO'
    | 'BUCKETS'
    | 'FAST'
    | 'FAILED_LOGIN_ATTEMPTS'
    | 'FILE'
    | 'FIELDS'
    | 'FILTER'
    | 'FIRST'
    | 'FORMAT'
    | 'HLL_UNION'
    | 'IMMEDIATE'
    | 'PATH'
    | 'PROFILE'
    | 'FUNCTION'
    | 'END'
    | 'ENGINE'
    | 'ENGINES'
    | 'ERRORS'
    | 'EXCLUDE'
    | 'EVENTS'
    | 'EXTERNAL'
    | 'GLOBAL'
    | 'GENERATED'
    | 'GENERIC'
    | 'GRAPH'
    | 'HASH'
    | 'HELP'
    | 'HOSTNAME'
    | 'HOURS'
    | 'HUB'
    | 'IDENTIFIED'
    | 'INDEXES'
    | 'INVERTED'
    | 'ISNULL'
    | 'ISOLATION'
    | 'JOB'
    | 'JOBS'
    | 'JSON'
    | 'JSONB'
    | 'ENCRYPTKEY'
    | 'ENCRYPTKEYS'
    | 'LABEL'
    | 'LAST'
    | 'LESS'
    | 'LEVEL'
    | 'LOCAL'
    | 'LOCATION'
    | 'LOCK'
    | 'UNLOCK'
    | 'MATERIALIZED'
    | 'MERGE'
    | 'MINUTES'
    | 'MODIFY'
    | 'NAME'
    | 'NAMES'
    | 'NEGATIVE'
    | 'NEVER'
    | 'NEXT'
    | 'NO'
    | 'NULLS'
    | 'OF'
    | 'OFFSET'
    | 'ONLY'
    | 'OPEN'
    | 'PARAMETER'
    | 'PARTITIONS'
    | 'PASSWORD'
    | 'PASSWORD_EXPIRE'
    | 'PASSWORD_REUSE'
    | 'PASSWORD_HISTORY'
    | 'PASSWORD_LOCK_TIME'
    | 'LDAP'
    | 'LDAP_ADMIN_PASSWORD'
    | 'PLUGIN'
    | 'PLUGINS'
    | 'PROC'
    | 'PROCESSLIST'
    | 'PROPERTIES'
    | 'CONDITIONS'
    | 'ACTIONS'
    | 'SET_SESSION_VARIABLE'
    | 'PROPERTY'
    | 'QUERY'
    | 'QUOTA'
    | 'RANDOM'
    | 'RECOVER'
    | 'REFRESH'
    | 'REPEATABLE'
    | 'REPLACE'
    | 'REPLACE_IF_NOT_NULL'
    | 'REPOSITORY'
    | 'REPOSITORIES'
    | 'RESOURCE'
    | 'RESOURCES'
    | 'RESTORE'
    | 'RETAIN'
    | 'RETENTION'
    | 'RETURNS'
    | 'ROLLBACK'
    | 'SCHEDULE'
    | 'SCHEMA'
    | 'SERIALIZABLE'
    | 'SESSION'
    | 'SKEW'
    | 'SNAPSHOT'
    | 'SNAPSHOTS'
    | 'SONAME'
    | 'SPLIT'
    | 'START'
    | 'STATUS'
    | 'STATS'
    | 'STORAGE'
    | 'STREAM'
    | 'STREAMING'
    | 'STRUCT'
    | 'STRING'
    | 'TABLES'
    | 'TEMPORARY'
    | 'THAN'
    | 'TIME'
    | 'TIMESTAMP'
    | 'TRANSACTION'
    | KW_TREE
    | 'TEXT'
    | 'TRIGGERS'
    | 'TRUNCATE'
    | 'TYPE'
    | 'TYPES'
    | 'UNCOMMITTED'
    | 'USER'
    | 'VARIABLE'
    | 'VARIABLES'
    | 'VALUE'
    | 'VERBOSE'
    | 'VERSION'
    | 'VIEW'
    | 'VIEWS'
    | 'WARNINGS'
    | 'WORK'
    | 'CLUSTER'
    | 'CLUSTERS'
    | 'COMPUTE'
    | 'LINK'
    | 'MIGRATE'
    | 'MIGRATIONS'
    | 'COUNT'
    | 'SUM'
    | 'MIN'
    | 'MAX'
    | 'FREE'
    | 'TAG'
    | 'TASK'
    | 'TASKS'
    | 'ROUTINE'
    | 'EVERY'
    | 'AT'
    | 'ROLLUP'
    | 'STARTS'
    | 'ENDS'
    | 'PAUSE'
    | 'RESUME'
    | 'STOP'
    | 'GROUPING'
    | 'GROUPS'
    | 'DYNAMIC'
    | time_unit
    | 'ENABLE'
    | 'FEATURE'
    | 'FRONTENDS'
    | 'MAP'
    | 'VARCHAR'
    | 'POLICY'
    | 'CURRENT_TIMESTAMP'
    | 'CATALOG'
    | 'CATALOGS'
    | 'RECYCLE'
    | 'MTMV'
    | 'LINES'
    | 'IGNORE'
    | 'HISTOGRAM'
    | 'CURRENT_CATALOG'
    | 'EXPIRED'
    | 'SAMPLE'
    | 'INCREMENTAL'
    | 'PERIOD'
    | 'PERCENT'
    | 'CRON'
    | 'UNSET'
    | 'BELONG'
    | 'STAGE'
    | 'VAULT'
    | 'VAULTS'
    | 'VARIANT'
    | 'IPV4'
    | 'IPV6'
    | 'MATCH_ANY'
    | 'MATCH_ALL'
    | 'MATCH_PHRASE'
    | 'MATCH_PHRASE_PREFIX'
    | 'MATCH_REGEXP'
    | 'MATCH_PHRASE_EDGE'
    | 'SQL'
    | 'CACHE'
    | 'COLOCATE'
    | 'COMPRESS_TYPE'
    | 'DORIS_INTERNAL_TABLE_ID'
    | 'HOTSPOT'
    | 'PRIVILEGES'
    | 'RECENT'
    | 'STAGES'
    | 'WARM'
    | 'UP'
    | 'CONVERT_LIGHT_SCHEMA_CHANGE_PROCESS'
    | 'ANALYZER'
    | 'TOKENIZER'
    | 'TOKEN_FILTER'
    ;

// sql_parser.cup:8796
ident
    : IDENTIFIER
    | keyword
    ;

// sql_parser.cup:8808
ident_or_text
    : ident
    | STRING_LITERAL
    ;

// sql_parser.cup:8819
time_unit
    : 'YEAR'
    | 'MONTH'
    | 'WEEK'
    | 'DAY'
    | 'HOUR'
    | 'MINUTE'
    | 'SECOND'
    ;
