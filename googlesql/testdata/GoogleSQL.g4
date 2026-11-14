//
// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// This file is derived from GoogleSQL's googlesql/parser/googlesql.tm,
// which is beside it as published, with the license it comes under
// (GoogleSQL.LICENSE). It was changed: googlesql/parser/cmd/tm2g4 wrote it
// from it, converting Textmapper to ANTLR 4 and keeping only what each rule
// derives. Actions, rule types, aliases and %prec are gone, and so are the
// rules unterminated_statement, the first rule, does not reach; a rule
// written name_opt is name?, and (x separator ",")+ is x (',' x)*. A token
// is written as the text it matches, and a token that matches a class of
// texts, or that the tokenizer makes in code, by its id, which tokens.tsv
// samples. Each rule is under the line of googlesql.tm it came from.
//
// Do not edit: rerun the converter.

grammar GoogleSQL;

// googlesql.tm:1151
unterminated_statement
    : pre_statement unterminated_sql_statement
    | pre_statement unterminated_script_statement
    ;

// googlesql.tm:1144
pre_statement
    : /* empty */
    ;

// googlesql.tm:1164
statement_level_hint
    : hint
    ;

// googlesql.tm:1172
unterminated_sql_statement
    : sql_statement_body
    | statement_level_hint sql_statement_body
    | 'DEFINE' 'MACRO'
    | statement_level_hint 'DEFINE' 'MACRO'
    | statement_level_hint KW_DEFINE_FOR_MACROS 'MACRO'
    ;

// googlesql.tm:1211
unterminated_unlabeled_script_statement
    : begin_end_block
    | while_statement
    | loop_statement
    | repeat_statement
    | for_in_statement
    ;

// googlesql.tm:1219
unterminated_script_statement
    : if_statement
    | case_statement
    | variable_declaration
    | break_statement
    | continue_statement
    | return_statement
    | raise_statement
    | unterminated_unlabeled_script_statement
    | label ':' unterminated_unlabeled_script_statement opt_identifier
    ;

// googlesql.tm:1256
sql_statement_body
    : sql_statement_body_no_pipe_suffix
    | sql_statement_body_maybe_pipe_suffix pipe_and_pipe_operator*
    ;

// googlesql.tm:1275
sql_statement_body_no_pipe_suffix
    : query_statement
    | subpipeline_statement
    | create_external_table_function_statement
    | create_external_table_statement
    | create_model_statement
    | create_table_function_statement
    | create_table_statement
    | create_view_statement
    | explain_statement
    | export_data_statement
    | export_model_statement
    | alter_statement
    | begin_statement
    | analyze_statement
    | assert_statement
    | aux_load_data_statement
    | clone_data_statement
    | merge_statement
    | truncate_statement
    | set_statement
    | commit_statement
    | start_batch_statement
    | run_batch_statement
    | abort_batch_statement
    | create_constant_statement
    | create_connection_statement
    | create_database_statement
    | create_function_statement
    | create_procedure_statement
    | create_index_statement
    | create_privilege_restriction_statement
    | create_row_access_policy_statement
    | create_sequence_statement
    | create_locality_group_statement
    | create_property_graph_statement
    | create_schema_statement
    | create_external_schema_statement
    | create_snapshot_statement
    | create_entity_statement
    | define_macro_statement
    | define_table_statement
    | export_metadata_statement
    | grant_statement
    | rename_statement
    | revoke_statement
    | rollback_statement
    | drop_all_row_access_policies_statement
    | drop_statement
    | import_statement
    | module_statement
    | undrop_statement
    | gql_statement
    | dml_statement
    ;

// googlesql.tm:1351
sql_statement_body_maybe_pipe_suffix
    : call_statement
    | describe_statement
    | execute_immediate
    | run_statement
    | show_statement
    ;

// googlesql.tm:1360
run_statement_arg
    : identifier ('=' | '=>') string_literal
    ;

// googlesql.tm:1366
run_statement_arg_list
    : ((run_statement_arg) (',' (run_statement_arg))*) ','?
    ;

// googlesql.tm:1372
run_statement
    : 'RUN' string_literal ('(' (run_statement_arg_list)? ')')?
    | 'RUN' path_expression '(' run_statement_arg_list? ')'
    ;

// googlesql.tm:1388
define_macro_statement
    : KW_DEFINE_FOR_MACROS 'MACRO' MACRO_BODY_TOKEN macro_body
    ;

// googlesql.tm:1427
macro_body
    : /* empty */
    | macro_token_list
    ;

// googlesql.tm:1438
macro_token_list
    : MACRO_BODY_TOKEN
    | macro_token_list MACRO_BODY_TOKEN
    ;

// googlesql.tm:1443
query_statement
    : query
    ;

// googlesql.tm:1450
subpipeline_statement
    : subpipeline_no_parens
    ;

// googlesql.tm:1457
alter_action
    : 'SET' 'OPTIONS' options_list
    | 'SET' 'AS' generic_entity_body
    | 'ADD' table_constraint_spec
    | 'ADD' primary_key_spec
    | 'ADD' 'CONSTRAINT' opt_if_not_exists identifier primary_key_or_table_constraint_spec
    | 'DROP' 'CONSTRAINT' opt_if_exists identifier
    | 'DROP' 'PRIMARY' 'KEY' opt_if_exists
    | 'ALTER' 'CONSTRAINT' opt_if_exists identifier constraint_enforcement
    | 'ALTER' 'CONSTRAINT' opt_if_exists identifier 'SET' 'OPTIONS' options_list
    | 'ADD' 'COLUMN' opt_if_not_exists table_column_definition opt_column_position opt_fill_using_expression
    | 'DROP' 'COLUMN' opt_if_exists identifier
    | 'RENAME' 'COLUMN' opt_if_exists identifier 'TO' identifier
    | 'ALTER' 'COLUMN' opt_if_exists identifier 'SET' 'DATA' 'TYPE' field_schema
    | 'ALTER' 'COLUMN' opt_if_exists identifier 'SET' 'OPTIONS' options_list
    | 'ALTER' 'COLUMN' opt_if_exists identifier 'SET' 'DEFAULT' expression
    | 'ALTER' 'COLUMN' opt_if_exists identifier 'DROP' 'DEFAULT'
    | 'ALTER' 'COLUMN' opt_if_exists identifier 'DROP' 'NOT' 'NULL'
    | 'ALTER' 'COLUMN' opt_if_exists identifier 'SET' 'GENERATED' generated_column_info_for_alter_column_action
    | 'ALTER' 'COLUMN' opt_if_exists identifier 'DROP' 'GENERATED'
    | 'RENAME' 'TO' path_expression
    | 'SET' 'DEFAULT' collate_clause
    | 'ADD' 'ROW' 'DELETION' 'POLICY' opt_if_not_exists '(' expression ')'
    | 'REPLACE' 'ROW' 'DELETION' 'POLICY' opt_if_exists '(' expression ')'
    | 'DROP' 'ROW' 'DELETION' 'POLICY' opt_if_exists
    | 'ALTER' generic_sub_entity_type opt_if_exists identifier alter_action
    | 'ADD' generic_sub_entity_type opt_if_not_exists identifier opt_options_list
    | 'DROP' generic_sub_entity_type opt_if_exists identifier
    | spanner_alter_column_action
    | spanner_set_on_delete_action
    ;

// googlesql.tm:1637
alter_action_list
    : alter_action
    | alter_action_list ',' alter_action
    ;

// googlesql.tm:1650
privilege_restriction_alter_action
    : restrict_to_clause
    | 'ADD' opt_if_not_exists possibly_empty_grantee_list
    | 'REMOVE' opt_if_exists possibly_empty_grantee_list
    ;

// googlesql.tm:1668
privilege_restriction_alter_action_list
    : privilege_restriction_alter_action
    | privilege_restriction_alter_action_list ',' privilege_restriction_alter_action
    ;

// googlesql.tm:1686
alter_index_action_list
    : alter_index_action
    | alter_index_action_list ',' alter_index_action
    ;

// googlesql.tm:1697
alter_index_action
    : 'SET' 'OPTIONS' options_list
    | 'ADD' 'COLUMN' opt_if_not_exists identifier opt_options_list
    | 'DROP' 'COLUMN' opt_if_exists identifier
    | 'ALTER' 'COLUMN' opt_if_exists identifier 'SET' 'OPTIONS' options_list
    | 'REBUILD'
    ;

// googlesql.tm:1729
row_access_policy_alter_action
    : grant_to_clause
    | 'FILTER' 'USING' '(' expression ')'
    | 'REVOKE' 'FROM' '(' grantee_list ')'
    | 'REVOKE' 'FROM' 'ALL'
    | 'RENAME' 'TO' identifier
    ;

// googlesql.tm:1756
row_access_policy_alter_action_list
    : row_access_policy_alter_action
    | row_access_policy_alter_action_list ',' row_access_policy_alter_action
    ;

// googlesql.tm:1776
schema_object_kind
    : 'AGGREGATE' 'FUNCTION'
    | 'APPROX' 'VIEW'
    | 'CONNECTION'
    | 'CONSTANT'
    | 'DATABASE'
    | 'EXTERNAL' table_or_table_function
    | 'EXTERNAL' 'SCHEMA'
    | 'FUNCTION'
    | 'INDEX'
    | 'MATERIALIZED' 'VIEW'
    | 'MODEL'
    | 'PROCEDURE'
    | 'SCHEMA'
    | 'SEQUENCE'
    | 'VIEW'
    | 'PROPERTY' 'GRAPH'
    ;

// googlesql.tm:1847
alter_statement
    : 'ALTER' table_or_table_function opt_if_exists maybe_dashed_path_expression alter_action_list
    | 'ALTER' schema_object_kind opt_if_exists path_expression alter_action_list
    | 'ALTER' generic_entity_type opt_if_exists path_expression alter_action_list
    | 'ALTER' generic_entity_type opt_if_exists alter_action_list
    | 'ALTER' 'PRIVILEGE' 'RESTRICTION' opt_if_exists 'ON' privilege_list 'ON' identifier path_expression privilege_restriction_alter_action_list
    | 'ALTER' 'ROW' 'ACCESS' 'POLICY' opt_if_exists identifier 'ON' path_expression row_access_policy_alter_action_list
    | 'ALTER' 'ALL' 'ROW' 'ACCESS' 'POLICIES' 'ON' path_expression row_access_policy_alter_action
    | 'ALTER' alter_index_type 'INDEX' opt_if_exists path_expression opt_on_path_expression alter_index_action_list
    ;

// googlesql.tm:1939
input_output_clause
    : 'INPUT' table_element_list 'OUTPUT' table_element_list
    ;

// googlesql.tm:1956
transform_clause
    : 'TRANSFORM' '(' select_list ')'
    ;

// googlesql.tm:1963
assert_statement
    : 'ASSERT' expression ('AS' string_literal)?
    ;

// googlesql.tm:1970
analyze_statement
    : 'ANALYZE' options_? ((table_and_column_info) (',' (table_and_column_info))*)?
    ;

// googlesql.tm:1978
table_and_column_info
    : maybe_dashed_path_expression opt_column_list
    ;

// googlesql.tm:1985
transaction_mode
    : 'READ' 'ONLY'
    | 'READ' 'WRITE'
    | 'ISOLATION' 'LEVEL' identifier
    | 'ISOLATION' 'LEVEL' identifier identifier
    ;

// googlesql.tm:2008
begin_statement
    : ('START' 'TRANSACTION' | 'BEGIN' 'TRANSACTION'?) ((transaction_mode) (',' (transaction_mode))*)?
    ;

// googlesql.tm:2016
set_statement
    : 'SET' 'TRANSACTION' ((transaction_mode) (',' (transaction_mode))*)
    | 'SET' identifier '=' expression
    | 'SET' named_parameter_expression '=' expression
    | 'SET' system_variable_expression '=' expression
    | 'SET' '(' identifier_list ')' '=' expression
    | 'SET' '(' ')'
    | 'SET' identifier ',' identifier_list '='
    ;

// googlesql.tm:2054
commit_statement
    : 'COMMIT' 'TRANSACTION'?
    ;

// googlesql.tm:2061
rollback_statement
    : 'ROLLBACK' 'TRANSACTION'?
    ;

// googlesql.tm:2068
start_batch_statement
    : 'START' 'BATCH' opt_identifier
    ;

// googlesql.tm:2075
run_batch_statement
    : 'RUN' 'BATCH'
    ;

// googlesql.tm:2082
abort_batch_statement
    : 'ABORT' 'BATCH'
    ;

// googlesql.tm:2089
create_constant_statement
    : 'CREATE' opt_or_replace opt_create_scope 'CONSTANT' opt_if_not_exists path_expression '=' expression
    ;

// googlesql.tm:2101
create_database_statement
    : 'CREATE' 'DATABASE' path_expression opt_options_list
    ;

// googlesql.tm:2108
unordered_options_body
    : options_ opt_as_sql_function_body_or_string
    | as_sql_function_body_or_string opt_options_list
    ;

// googlesql.tm:2126
create_function_statement
    : 'CREATE' opt_or_replace opt_create_scope 'AGGREGATE'? 'FUNCTION' opt_if_not_exists function_declaration ('RETURNS' (type | templated_parameter_type))? sql_security? determinism_level? language_or_remote_with_connection? unordered_options_body?
    ;

// googlesql.tm:2158
function_declaration
    : path_expression function_parameters
    ;

// googlesql.tm:2165
function_parameter
    : identifier? type_or_tvf_schema ('AS' identifier)? opt_default_expression ('NOT' 'AGGREGATE')?
    ;

// googlesql.tm:2181
function_parameters
    : '(' ((function_parameter) (',' (function_parameter))*)? ')'
    ;

// googlesql.tm:2188
create_procedure_statement
    : 'CREATE' opt_or_replace opt_create_scope 'PROCEDURE' opt_if_not_exists path_expression procedure_parameters ('EXTERNAL' 'SECURITY' ('INVOKER' | 'DEFINER'))? with_connection_clause? options_? (begin_end_block | 'LANGUAGE' identifier ('AS' string_literal)?)
    ;

// googlesql.tm:2220
procedure_parameters
    : '(' ((procedure_parameter) (',' (procedure_parameter))*)? ')'
    ;

// googlesql.tm:2227
procedure_parameter
    : procedure_parameter_mode? identifier (type_or_tvf_schema | ')' | ',')
    ;

// googlesql.tm:2250
procedure_parameter_mode
    : 'IN'
    | 'OUT'
    | 'INOUT'
    ;

// googlesql.tm:2265
returns_
    : 'RETURNS' type_or_tvf_schema
    ;

// googlesql.tm:2279
determinism_level
    : 'DETERMINISTIC'
    | 'NOT' 'DETERMINISTIC'
    | 'IMMUTABLE'
    | 'STABLE'
    | 'VOLATILE'
    ;

// googlesql.tm:2302
language
    : 'LANGUAGE' identifier
    ;

// googlesql.tm:2309
remote_with_connection_clause
    : 'REMOTE' with_connection_clause?
    | with_connection_clause
    ;

// googlesql.tm:2328
language_or_remote_with_connection
    : 'LANGUAGE' identifier remote_with_connection_clause?
    | remote_with_connection_clause language?
    ;

// googlesql.tm:2343
sql_security
    : 'SQL' 'SECURITY' ('INVOKER' | 'DEFINER')
    ;

// googlesql.tm:2352
as_sql_function_body_or_string
    : 'AS' sql_function_body
    | 'AS' string_literal
    ;

// googlesql.tm:2370
opt_as_sql_function_body_or_string
    : as_sql_function_body_or_string?
    ;

// googlesql.tm:2391
path_expression_or_string
    : path_expression
    | string_literal
    ;

// googlesql.tm:2402
path_expression_or_default
    : path_expression
    | 'DEFAULT'
    ;

// googlesql.tm:2413
sql_function_body
    : '(' expression ')'
    | '(' 'SELECT'
    ;

// googlesql.tm:2431
restrict_to_clause
    : 'RESTRICT' 'TO' possibly_empty_grantee_list
    ;

// googlesql.tm:2440
grant_to_clause
    : 'GRANT' 'TO' '(' grantee_list ')'
    ;

// googlesql.tm:2449
create_row_access_policy_grant_to_clause
    : grant_to_clause
    | 'TO' grantee_list
    ;

// googlesql.tm:2459
filter_using_clause
    : 'FILTER'? 'USING' '(' expression ')'
    ;

// googlesql.tm:2469
create_privilege_restriction_statement
    : 'CREATE' opt_or_replace 'PRIVILEGE' 'RESTRICTION' opt_if_not_exists 'ON' privilege_list 'ON' identifier path_expression restrict_to_clause?
    ;

// googlesql.tm:2484
create_row_access_policy_statement
    : 'CREATE' opt_or_replace 'ROW' opt_access 'POLICY' opt_if_not_exists opt_identifier 'ON' path_expression create_row_access_policy_grant_to_clause? filter_using_clause
    ;

// googlesql.tm:2504
with_partition_columns_clause
    : 'WITH' 'PARTITION' 'COLUMNS' opt_table_element_list
    ;

// googlesql.tm:2513
with_connection_clause
    : 'WITH' connection_clause
    ;

// googlesql.tm:2520
create_external_table_statement
    : 'CREATE' opt_or_replace opt_create_scope 'EXTERNAL' 'TABLE' opt_if_not_exists maybe_dashed_path_expression table_element_list? like_path_expression? default_collate_clause? with_partition_columns_clause? with_connection_clause? options_
    ;

// googlesql.tm:2544
create_external_table_function_statement
    : 'CREATE' opt_or_replace opt_create_scope 'EXTERNAL' 'TABLE' 'FUNCTION'
    ;

// googlesql.tm:2559
create_index_statement_suffix
    : partition_by_clause_prefix_no_hint opt_options_list
    | opt_options_list spanner_index_interleave_clause
    | options_
    ;

// googlesql.tm:2577
create_index_statement
    : 'CREATE' opt_or_replace opt_unique opt_spanner_null_filtered opt_index_type 'INDEX' opt_if_not_exists path_expression on_path_expression as_alias? opt_index_unnest_expression_list index_order_by_and_options opt_index_storing_list create_index_statement_suffix?
    ;

// googlesql.tm:2615
braced_graph_subquery
    : '{' graph_operation_block '}'
    | '{' 'GRAPH' path_expression graph_operation_block '}'
    ;

// googlesql.tm:2626
exists_graph_pattern_subquery
    : 'EXISTS' hint? '{' graph_pattern '}'
    | 'EXISTS' hint? '{' 'GRAPH' path_expression graph_pattern '}'
    ;

// googlesql.tm:2639
exists_linear_ops_subquery
    : 'EXISTS' hint? '{' graph_linear_operator_list '}'
    | 'EXISTS' hint? '{' 'GRAPH' path_expression graph_linear_operator_list '}'
    ;

// googlesql.tm:2652
exists_graph_subquery
    : exists_graph_pattern_subquery
    | exists_linear_ops_subquery
    | 'EXISTS' hint? braced_graph_subquery
    ;

// googlesql.tm:2662
create_property_graph_statement
    : 'CREATE' opt_or_replace opt_create_scope 'PROPERTY' 'GRAPH' opt_if_not_exists path_expression 'NODE' 'TABLES' element_table_list edge_table_clause? opt_options_list
    ;

// googlesql.tm:2681
element_table_list
    : '(' ((element_table_definition) (',' (element_table_definition))*) ','? ')'
    ;

// googlesql.tm:2688
element_table_definition
    : path_expression as_alias_with_required_as? key_clause? source_node_table_clause? dest_node_table_clause? opt_options_list opt_label_and_properties_clause dynamic_label_and_properties?
    ;

// googlesql.tm:2712
key_clause
    : 'KEY' column_list
    ;

// googlesql.tm:2719
source_node_table_clause
    : 'SOURCE' 'KEY' column_list 'REFERENCES' identifier opt_column_list
    ;

// googlesql.tm:2729
dest_node_table_clause
    : 'DESTINATION' 'KEY' column_list 'REFERENCES' identifier opt_column_list
    ;

// googlesql.tm:2739
edge_table_clause
    : 'EDGE' 'TABLES' element_table_list
    ;

// googlesql.tm:2746
opt_label_and_properties_clause
    : /* empty */
    | properties_clause
    | label_and_properties+
    ;

// googlesql.tm:2768
label_and_properties
    : ('DEFAULT' 'LABEL' opt_options_list | 'LABEL' identifier) properties_clause?
    ;

// googlesql.tm:2786
properties_clause
    : 'NO' 'PROPERTIES'
    | properties_all_columns except_column_list?
    | 'PROPERTIES' '(' derived_property_list ')'
    ;

// googlesql.tm:2809
properties_all_columns
    : 'PROPERTIES' 'ALL' 'COLUMNS'
    | 'PROPERTIES' 'ARE' 'ALL' 'COLUMNS'
    ;

// googlesql.tm:2814
except_column_list
    : 'EXCEPT' column_list
    ;

// googlesql.tm:2821
derived_property_list
    : ((derived_property) (',' (derived_property))*)
    ;

// googlesql.tm:2828
derived_property
    : expression as_alias_with_required_as? opt_options_list
    ;

// googlesql.tm:2835
dynamic_label_or_properties
    : 'DYNAMIC' 'LABEL' '(' dynamic_label ')'
    | 'DYNAMIC' 'PROPERTIES' '(' dynamic_property ')'
    ;

// googlesql.tm:2846
static_label_or_properties_keywords
    : 'NO' 'PROPERTIES'
    | 'DEFAULT' 'LABEL'
    | 'LABEL'
    | 'PROPERTIES'
    ;

// googlesql.tm:2853
dynamic_label_and_properties
    : dynamic_label_or_properties
    | dynamic_label_and_properties dynamic_label_or_properties
    | dynamic_label_or_properties static_label_or_properties_keywords
    ;

// googlesql.tm:2875
dynamic_label
    : expression
    ;

// googlesql.tm:2882
dynamic_property
    : expression
    ;

// googlesql.tm:2889
create_schema_statement
    : 'CREATE' opt_or_replace 'SCHEMA' opt_if_not_exists path_expression opt_default_collate_clause opt_options_list
    ;

// googlesql.tm:2900
create_locality_group_statement
    : 'CREATE' 'LOCALITY' 'GROUP' path_expression opt_options_list
    ;

// googlesql.tm:2911
create_external_schema_statement
    : 'CREATE' opt_or_replace opt_create_scope 'EXTERNAL' 'SCHEMA' opt_if_not_exists path_expression opt_with_connection_clause options_
    ;

// googlesql.tm:2923
create_connection_statement
    : 'CREATE' opt_or_replace 'CONNECTION' opt_if_not_exists path_expression opt_options_list
    ;

// googlesql.tm:2934
undrop_statement
    : 'UNDROP' schema_object_kind opt_if_not_exists path_expression at_system_time? opt_options_list
    ;

// googlesql.tm:2954
create_snapshot_statement
    : 'CREATE' opt_or_replace 'SNAPSHOT' 'TABLE' opt_if_not_exists maybe_dashed_path_expression 'CLONE' clone_data_source opt_options_list
    | 'CREATE' opt_or_replace 'SNAPSHOT' schema_object_kind opt_if_not_exists maybe_dashed_path_expression 'CLONE' clone_data_source opt_options_list
    ;

// googlesql.tm:2983
unordered_language_options
    : language options_?
    | options_ language?
    ;

// googlesql.tm:3000
create_table_function_statement
    : 'CREATE' opt_or_replace opt_create_scope 'TABLE' 'FUNCTION' opt_if_not_exists path_expression opt_function_parameters returns_? sql_security? unordered_language_options? opt_as_query_or_string
    ;

// googlesql.tm:3042
create_table_statement_prefix
    : 'CREATE' opt_or_replace opt_create_scope 'TABLE' opt_if_not_exists maybe_dashed_path_expression opt_table_element_list opt_spanner_table_options opt_like_path_expression opt_clone_table opt_copy_table opt_default_collate_clause opt_partition_by_clause_no_hint opt_cluster_by_clause_no_hint opt_ttl_clause opt_with_connection_clause opt_options_list
    ;

// googlesql.tm:3072
create_table_statement
    : create_table_statement_prefix opt_as_query
    ;

// googlesql.tm:3080
append_or_overwrite
    : 'INTO'
    | 'OVERWRITE'
    ;

// googlesql.tm:3092
aux_load_data_from_files_options_list
    : 'FROM' 'FILES' options_list
    ;

// googlesql.tm:3099
load_data_partitions_clause
    : 'OVERWRITE'? 'PARTITIONS' '(' expression ')'
    ;

// googlesql.tm:3115
maybe_dashed_path_expression_with_scope
    : 'TEMP' 'TABLE' maybe_dashed_path_expression
    | 'TEMPORARY' 'TABLE' maybe_dashed_path_expression
    | maybe_dashed_path_expression
    ;

// googlesql.tm:3133
aux_load_data_statement
    : 'LOAD' 'DATA' append_or_overwrite maybe_dashed_path_expression_with_scope table_element_list? load_data_partitions_clause? collate_clause? partition_by_clause_prefix_no_hint? cluster_by_clause_no_hint? options_? aux_load_data_from_files_options_list with_partition_columns_clause? with_connection_clause?
    ;

// googlesql.tm:3165
generic_entity_type_unchecked
    : IDENTIFIER
    | 'PROJECT'
    ;

// googlesql.tm:3179
generic_entity_type
    : generic_entity_type_unchecked
    ;

// googlesql.tm:3198
sub_entity_type_identifier
    : IDENTIFIER
    | 'REPLICA'
    ;

// googlesql.tm:3209
generic_sub_entity_type
    : sub_entity_type_identifier
    ;

// googlesql.tm:3221
generic_entity_body
    : json_literal
    | string_literal
    ;

// googlesql.tm:3232
opt_generic_entity_body
    : 'AS' generic_entity_body
    | /* empty */
    ;

// googlesql.tm:3243
create_entity_statement
    : 'CREATE' opt_or_replace generic_entity_type opt_if_not_exists path_expression opt_options_list opt_generic_entity_body
    ;

// googlesql.tm:3259
create_model_statement
    : 'CREATE' opt_or_replace opt_create_scope 'MODEL' opt_if_not_exists path_expression input_output_clause? transform_clause? remote_with_connection_clause? options_? as_query_or_aliased_query_list?
    ;

// googlesql.tm:3278
opt_table_element_list
    : table_element_list
    | /* empty */
    ;

// googlesql.tm:3286
table_element_list
    : table_element_list_prefix ')'
    | '(' ')'
    ;

// googlesql.tm:3301
table_element_list_prefix
    : '(' table_element
    | table_element_list_prefix ',' table_element
    | table_element_list_prefix ','
    ;

// googlesql.tm:3348
table_element
    : table_column_definition
    | table_constraint_definition
    ;

// googlesql.tm:3353
table_column_definition
    : identifier table_column_schema column_attributes? opt_options_list
    ;

// googlesql.tm:3362
table_column_schema
    : column_schema_inner opt_collate_clause opt_column_info
    | generated_column_info
    ;

// googlesql.tm:3374
simple_column_schema_inner
    : path_expression
    | 'INTERVAL'
    ;

// googlesql.tm:3388
array_column_schema_inner
    : 'ARRAY' template_type_open field_schema template_type_close
    ;

// googlesql.tm:3395
range_column_schema_inner
    : 'RANGE' template_type_open field_schema template_type_close
    ;

// googlesql.tm:3402
map_column_schema_inner
    : 'MAP' template_type_open field_schema ',' field_schema template_type_close
    ;

// googlesql.tm:3409
struct_column_field
    : column_schema_inner opt_collate_clause opt_field_attributes
    | identifier field_schema
    ;

// googlesql.tm:3438
struct_column_schema_prefix
    : 'STRUCT' template_type_open struct_column_field
    | struct_column_schema_prefix ',' struct_column_field
    ;

// googlesql.tm:3451
struct_column_schema_inner
    : 'STRUCT' template_type_open template_type_close
    | struct_column_schema_prefix template_type_close
    ;

// googlesql.tm:3459
raw_column_schema_inner
    : simple_column_schema_inner
    | array_column_schema_inner
    | struct_column_schema_inner
    | range_column_schema_inner
    | map_column_schema_inner
    ;

// googlesql.tm:3467
column_schema_inner
    : raw_column_schema_inner opt_type_parameters
    ;

// googlesql.tm:3474
generated_mode
    : 'GENERATED' 'AS'
    | 'GENERATED' 'ALWAYS' 'AS'
    | 'GENERATED' 'BY' 'DEFAULT' 'AS'
    | 'AS'
    ;

// googlesql.tm:3493
stored_mode
    : 'STORED' 'VOLATILE'
    | 'STORED'
    | /* empty */
    ;

// googlesql.tm:3508
signed_numerical_literal
    : integer_literal
    | numeric_literal
    | bignumeric_literal
    | floating_point_literal
    | '-' integer_literal
    | '-' floating_point_literal
    ;

// googlesql.tm:3530
opt_start_with
    : 'START' 'WITH' signed_numerical_literal
    | /* empty */
    ;

// googlesql.tm:3541
opt_increment_by
    : 'INCREMENT' 'BY' signed_numerical_literal
    | /* empty */
    ;

// googlesql.tm:3552
opt_maxvalue
    : 'MAXVALUE' signed_numerical_literal
    | /* empty */
    ;

// googlesql.tm:3563
opt_minvalue
    : 'MINVALUE' signed_numerical_literal
    | /* empty */
    ;

// googlesql.tm:3574
opt_cycle
    : 'CYCLE'
    | 'NO' 'CYCLE'
    | /* empty */
    ;

// googlesql.tm:3589
identity_column_info
    : 'IDENTITY' '(' opt_start_with opt_increment_by opt_maxvalue opt_minvalue opt_cycle ')'
    ;

// googlesql.tm:3600
generated_column_info
    : generated_mode '(' expression ')' stored_mode
    | generated_mode identity_column_info
    ;

// googlesql.tm:3616
generated_mode_for_alter_column_action
    : 'ALWAYS'? 'AS'
    | 'BY' 'DEFAULT' 'AS'
    ;

// googlesql.tm:3627
generated_column_info_for_alter_column_action
    : generated_mode_for_alter_column_action identity_column_info
    ;

// googlesql.tm:3636
default_column_info
    : 'DEFAULT' expression
    ;

// googlesql.tm:3648
opt_column_info
    : generated_column_info default_column_info
    | generated_column_info
    | default_column_info generated_column_info
    | default_column_info
    | /* empty */
    ;

// googlesql.tm:3687
field_schema
    : column_schema_inner opt_collate_clause opt_field_attributes opt_options_list
    ;

// googlesql.tm:3694
primary_key_column_attribute
    : 'PRIMARY' 'KEY'
    ;

// googlesql.tm:3701
foreign_key_column_attribute
    : opt_constraint_identity foreign_key_reference
    ;

// googlesql.tm:3709
hidden_column_attribute
    : 'HIDDEN'
    ;

// googlesql.tm:3716
not_null_column_attribute
    : 'NOT' 'NULL'
    ;

// googlesql.tm:3723
column_attribute
    : primary_key_column_attribute
    | foreign_key_column_attribute
    | hidden_column_attribute
    | not_null_column_attribute
    ;

// googlesql.tm:3746
column_attributes
    : column_attribute
    | column_attributes column_attribute
    | column_attributes constraint_enforcement
    ;

// googlesql.tm:3784
opt_field_attributes
    : not_null_column_attribute
    | /* empty */
    ;

// googlesql.tm:3795
column_position
    : 'PRECEDING' identifier
    | 'FOLLOWING' identifier
    ;

// googlesql.tm:3810
opt_column_position
    : column_position
    | /* empty */
    ;

// googlesql.tm:3818
fill_using_expression
    : 'FILL' 'USING' expression
    ;

// googlesql.tm:3825
opt_fill_using_expression
    : fill_using_expression
    | /* empty */
    ;

// googlesql.tm:3833
table_constraint_spec
    : 'CHECK' '(' expression ')' opt_constraint_enforcement opt_options_list
    | 'FOREIGN' 'KEY' column_list foreign_key_reference opt_constraint_enforcement opt_options_list
    ;

// googlesql.tm:3849
primary_key_element
    : identifier opt_asc_or_desc opt_null_order
    ;

// googlesql.tm:3869
primary_key_element_list_prefix
    : '(' primary_key_element
    | primary_key_element_list_prefix ',' primary_key_element
    ;

// googlesql.tm:3880
primary_key_element_list
    : primary_key_element_list_prefix ')'
    | '(' ')'
    ;

// googlesql.tm:3891
primary_key_spec
    : 'PRIMARY' 'KEY' primary_key_element_list opt_constraint_enforcement opt_options_list
    ;

// googlesql.tm:3901
primary_key_or_table_constraint_spec
    : primary_key_spec
    | table_constraint_spec
    ;

// googlesql.tm:3908
table_constraint_definition
    : primary_key_spec
    | table_constraint_spec
    | identifier identifier table_constraint_spec
    ;

// googlesql.tm:3936
foreign_key_reference
    : 'REFERENCES' path_expression column_list opt_foreign_key_match opt_foreign_key_actions
    ;

// googlesql.tm:3946
opt_foreign_key_match
    : 'MATCH' foreign_key_match_mode
    | /* empty */
    ;

// googlesql.tm:3957
foreign_key_match_mode
    : 'SIMPLE'
    | 'FULL'
    | KW_NOT_SPECIAL 'DISTINCT'
    ;

// googlesql.tm:3972
opt_foreign_key_actions
    : foreign_key_on_update opt_foreign_key_on_delete
    | foreign_key_on_delete opt_foreign_key_on_update
    | /* empty */
    ;

// googlesql.tm:3993
opt_foreign_key_on_update
    : foreign_key_on_update
    | /* empty */
    ;

// googlesql.tm:4001
opt_foreign_key_on_delete
    : foreign_key_on_delete
    | /* empty */
    ;

// googlesql.tm:4009
foreign_key_on_update
    : 'ON' 'UPDATE' foreign_key_action
    ;

// googlesql.tm:4016
foreign_key_on_delete
    : 'ON' 'DELETE' foreign_key_action
    ;

// googlesql.tm:4023
foreign_key_action
    : 'NO' 'ACTION'
    | 'RESTRICT'
    | 'CASCADE'
    | 'SET' 'NULL'
    ;

// googlesql.tm:4042
opt_constraint_identity
    : 'CONSTRAINT' identifier
    | /* empty */
    ;

// googlesql.tm:4053
opt_constraint_enforcement
    : constraint_enforcement
    | /* empty */
    ;

// googlesql.tm:4061
constraint_enforcement
    : 'ENFORCED'
    | 'NOT' 'ENFORCED'
    ;

// googlesql.tm:4074
table_or_table_function
    : 'TABLE' 'FUNCTION'
    | 'TABLE'
    ;

// googlesql.tm:4085
tvf_schema_column
    : identifier type
    | type
    ;

// googlesql.tm:4096
tvf_schema_prefix
    : 'TABLE' template_type_open tvf_schema_column
    | tvf_schema_prefix ',' tvf_schema_column
    ;

// googlesql.tm:4108
tvf_schema
    : tvf_schema_prefix template_type_close
    ;

// googlesql.tm:4115
create_view_statement
    : 'CREATE' opt_or_replace opt_create_scope 'RECURSIVE'? 'VIEW' opt_if_not_exists maybe_dashed_path_expression column_with_options_list? sql_security? options_? as_query
    | 'CREATE' opt_or_replace 'MATERIALIZED' 'RECURSIVE'? 'VIEW' opt_if_not_exists maybe_dashed_path_expression column_with_options_list? sql_security? opt_partition_by_clause_no_hint cluster_by_clause_no_hint? options_? as_before_query (query | 'REPLICA' 'OF' maybe_dashed_path_expression)
    | 'CREATE' opt_or_replace 'APPROX' 'RECURSIVE'? 'VIEW' opt_if_not_exists maybe_dashed_path_expression column_with_options_list? sql_security? options_? as_query
    ;

// googlesql.tm:4169
as_before_query
    : 'AS'
    ;

// googlesql.tm:4176
as_query
    : as_before_query query_after_as
    ;

// googlesql.tm:4190
opt_as_query
    : as_query
    | /* empty */
    ;

// googlesql.tm:4201
as_query_or_string
    : as_query
    | as_before_query string_literal
    ;

// googlesql.tm:4212
opt_as_query_or_string
    : as_query_or_string?
    ;

// googlesql.tm:4219
as_query_or_aliased_query_list
    : as_query
    | as_before_query '(' aliased_query_list ')'
    ;

// googlesql.tm:4230
if_not_exists
    : 'IF' 'NOT' 'EXISTS'
    ;

// googlesql.tm:4234
opt_if_not_exists
    : if_not_exists?
    ;

// googlesql.tm:4241
describe_statement
    : describe_keyword describe_info
    ;

// googlesql.tm:4248
describe_info
    : identifier maybe_slashed_or_dashed_path_expression from_path_expression?
    | maybe_slashed_or_dashed_path_expression from_path_expression?
    ;

// googlesql.tm:4259
from_path_expression
    : 'FROM' maybe_slashed_or_dashed_path_expression
    ;

// googlesql.tm:4266
explain_statement
    : 'EXPLAIN' unterminated_sql_statement
    ;

// googlesql.tm:4273
export_data_no_query
    : 'EXPORT' 'DATA' opt_with_connection_clause opt_options_list
    ;

// googlesql.tm:4281
export_data_statement
    : export_data_no_query as_query
    ;

// googlesql.tm:4290
export_model_statement
    : 'EXPORT' 'MODEL' path_expression opt_with_connection_clause opt_options_list
    ;

// googlesql.tm:4297
export_metadata_statement
    : 'EXPORT' table_or_table_function 'METADATA' 'FROM' maybe_dashed_path_expression opt_with_connection_clause opt_options_list
    ;

// googlesql.tm:4313
grant_statement
    : 'GRANT' privileges 'ON' identifier path_expression 'TO' grantee_list
    | 'GRANT' privileges 'ON' identifier identifier path_expression 'TO' grantee_list
    | 'GRANT' privileges 'ON' path_expression 'TO' grantee_list
    ;

// googlesql.tm:4329
revoke_statement
    : 'REVOKE' privileges 'ON' identifier path_expression 'FROM' grantee_list
    | 'REVOKE' privileges 'ON' identifier identifier path_expression 'FROM' grantee_list
    | 'REVOKE' privileges 'ON' path_expression 'FROM' grantee_list
    ;

// googlesql.tm:4345
privileges
    : 'ALL' 'PRIVILEGES'?
    | privilege_list
    ;

// googlesql.tm:4356
privilege_list
    : ((privilege) (',' (privilege))*)
    ;

// googlesql.tm:4363
privilege
    : privilege_name opt_path_expression_list_with_parens
    ;

// googlesql.tm:4370
privilege_name
    : identifier
    | 'SELECT'
    ;

// googlesql.tm:4382
rename_statement
    : 'RENAME' identifier path_expression 'TO' path_expression
    ;

// googlesql.tm:4389
import_statement
    : 'IMPORT' import_type path_expression_or_string opt_as_or_into_alias opt_options_list
    ;

// googlesql.tm:4406
module_statement
    : 'MODULE' path_expression opt_options_list
    ;

// googlesql.tm:4414
column_ordering_and_options_expr
    : expression opt_collate_clause opt_asc_or_desc opt_null_order opt_options_list
    ;

// googlesql.tm:4429
index_order_by_and_options_prefix
    : '(' column_ordering_and_options_expr
    | index_order_by_and_options_prefix ',' column_ordering_and_options_expr
    ;

// googlesql.tm:4440
all_column_column_options
    : index_order_by_and_options_prefix ')'
    ;

// googlesql.tm:4444
with_column_options
    : 'WITH' 'COLUMN' 'OPTIONS' all_column_column_options
    ;

// googlesql.tm:4451
index_all_columns
    : '(' 'ALL' 'COLUMNS' with_column_options? ')'
    ;

// googlesql.tm:4465
index_order_by_and_options
    : index_order_by_and_options_prefix ')'
    | index_all_columns
    ;

// googlesql.tm:4476
index_unnest_expression_list
    : (unnest_expression_with_opt_alias_and_offset)+
    ;

// googlesql.tm:4483
opt_index_unnest_expression_list
    : index_unnest_expression_list?
    ;

// googlesql.tm:4490
index_storing_list
    : 'STORING' '(' ((expression) (',' (expression))*) ')'
    ;

// googlesql.tm:4497
opt_index_storing_list
    : index_storing_list?
    ;

// googlesql.tm:4504
column_list
    : '(' ((identifier) (',' (identifier))*) ')'
    ;

// googlesql.tm:4515
opt_column_list
    : column_list?
    ;

// googlesql.tm:4522
column_with_options
    : identifier opt_options_list
    ;

// googlesql.tm:4529
column_with_options_list
    : '(' ((column_with_options) (',' (column_with_options))*) ')'
    ;

// googlesql.tm:4536
grantee_list
    : ((string_literal_or_parameter) (',' (string_literal_or_parameter))*)
    ;

// googlesql.tm:4543
possibly_empty_grantee_list
    : '(' grantee_list ')'
    | '(' ')'
    ;

// googlesql.tm:4554
show_statement
    : 'SHOW' show_target from_path_expression? like_string_literal?
    ;

// googlesql.tm:4561
show_target
    : 'MATERIALIZED' 'VIEWS'
    | identifier
    ;

// googlesql.tm:4572
like_string_literal
    : 'LIKE' string_literal
    ;

// googlesql.tm:4579
like_path_expression
    : 'LIKE' maybe_dashed_path_expression
    ;

// googlesql.tm:4586
opt_like_path_expression
    : like_path_expression?
    ;

// googlesql.tm:4593
clone_table
    : 'CLONE' clone_data_source
    ;

// googlesql.tm:4600
opt_clone_table
    : clone_table?
    ;

// googlesql.tm:4607
copy_table
    : 'COPY' copy_data_source
    ;

// googlesql.tm:4614
opt_copy_table
    : copy_table?
    ;

// googlesql.tm:4621
all_or_distinct
    : 'ALL'
    | 'DISTINCT'
    ;

// googlesql.tm:4634
distinct
    : 'DISTINCT'
    ;

// googlesql.tm:4638
opt_distinct
    : distinct?
    ;

// googlesql.tm:4647
query_set_operation_type
    : 'UNION'
    | KW_EXCEPT_IN_SET_OP
    | 'INTERSECT'
    ;

// googlesql.tm:4665
query_primary_or_set_operation
    : query_primary
    | query_set_operation
    ;

// googlesql.tm:4670
parenthesized_query
    : '(' query ')'
    ;

// googlesql.tm:4683
select_or_from_keyword
    : 'SELECT'
    | 'FROM'
    ;

// googlesql.tm:4690
bad_keyword_after_from_query
    : 'WHERE'
    | 'SELECT'
    | 'GROUP'
    ;

// googlesql.tm:4706
bad_keyword_after_from_query_allows_parens
    : 'ORDER'
    | 'UNION'
    | 'INTERSECT'
    | KW_EXCEPT_IN_SET_OP
    | 'LIMIT'
    ;

// googlesql.tm:4729
query_without_pipe_operators
    : with_clause query_primary_or_set_operation opt_order_by_clause opt_limit_offset_clause opt_lock_mode_clause
    | with_clause_with_trailing_comma select_or_from_keyword
    | with_clause '|>'
    | query_primary_or_set_operation opt_order_by_clause opt_limit_offset_clause opt_lock_mode_clause
    | opt_with_clause from_clause opt_lock_mode_clause
    | opt_with_clause from_clause bad_keyword_after_from_query
    | opt_with_clause from_clause bad_keyword_after_from_query_allows_parens
    ;

// googlesql.tm:4818
query
    : query_without_pipe_operators
    | query pipe_and_pipe_operator
    ;

// googlesql.tm:4838
query_after_as
    : query
    | gql_query
    ;

// googlesql.tm:4844
subpipeline_bad_prefix_subquery
    : 'SELECT'
    | 'FROM'
    | 'WITH'
    ;

// googlesql.tm:4853
subpipeline_bad_prefix_not_subquery
    : identifier
    | 'WHERE'
    | 'LIMIT'
    | 'JOIN'
    | 'ORDER'
    | 'GROUP'
    | '('
    ;

// googlesql.tm:4866
subpipeline_prefix_invalid
    : '(' subpipeline_bad_prefix_subquery
    | '(' subpipeline_bad_prefix_not_subquery
    ;

// googlesql.tm:4881
subpipeline_no_parens
    : pipe_and_pipe_operator+
    ;

// googlesql.tm:4889
subpipeline_with_parens
    : '(' pipe_and_pipe_operator* ')'
    ;

// googlesql.tm:4899
subpipeline
    : subpipeline_with_parens
    | subpipeline_prefix_invalid
    ;

// googlesql.tm:4904
opt_subpipeline
    : subpipeline
    | /* empty */
    ;

// googlesql.tm:4912
pipe_and_pipe_operator
    : '|>' pipe_operator
    ;

// googlesql.tm:4920
pipe_operator
    : pipe_where
    | pipe_select
    | pipe_extend
    | pipe_rename
    | pipe_aggregate
    | pipe_group_by
    | pipe_limit_offset
    | pipe_set_operation
    | pipe_recursive_union
    | pipe_order_by
    | pipe_join
    | pipe_call
    | pipe_window
    | pipe_distinct
    | pipe_tablesample
    | pipe_as
    | pipe_describe
    | pipe_static_describe
    | pipe_assert
    | pipe_log
    | pipe_drop
    | pipe_set
    | pipe_pivot
    | pipe_unpivot
    | pipe_if
    | pipe_elseif
    | pipe_else
    | pipe_fork
    | pipe_tee
    | pipe_with
    | pipe_export_data
    | pipe_create_table
    | pipe_insert
    | pipe_match_recognize
    ;

// googlesql.tm:4957
pipe_where
    : where_clause
    ;

// googlesql.tm:4964
pipe_select
    : select_clause opt_window_clause_with_trailing_comma
    ;

// googlesql.tm:4972
pipe_limit_offset
    : limit_offset_clause
    ;

// googlesql.tm:4979
pipe_order_by
    : order_by_clause_with_opt_comma
    ;

// googlesql.tm:4986
pipe_extend
    : 'EXTEND' pipe_selection_item_list opt_window_clause_with_trailing_comma
    ;

// googlesql.tm:4998
pipe_selection_item
    : select_column_expr
    | select_column_dot_star
    ;

// googlesql.tm:5005
pipe_selection_item_with_order
    : select_column_expr opt_selection_item_order
    | select_column_dot_star
    ;

// googlesql.tm:5015
pipe_selection_item_list_no_comma
    : pipe_selection_item
    | pipe_selection_item_list_no_comma ',' pipe_selection_item
    ;

// googlesql.tm:5026
pipe_selection_item_list_no_comma_with_order
    : pipe_selection_item_with_order
    | pipe_selection_item_list_no_comma_with_order ',' pipe_selection_item_with_order
    ;

// googlesql.tm:5040
pipe_selection_item_list
    : pipe_selection_item_list_no_comma ','?
    ;

// googlesql.tm:5049
pipe_selection_item_list_with_order
    : pipe_selection_item_list_no_comma_with_order ','?
    ;

// googlesql.tm:5056
pipe_selection_item_list_with_order_or_empty
    : pipe_selection_item_list_with_order
    | /* empty */
    ;

// googlesql.tm:5067
pipe_rename_item
    : identifier 'AS'? identifier
    | identifier '.'
    ;

// googlesql.tm:5082
pipe_rename_item_list
    : pipe_rename_item
    | pipe_rename_item_list ',' pipe_rename_item
    ;

// googlesql.tm:5093
pipe_rename
    : 'RENAME' pipe_rename_item_list ','?
    ;

// googlesql.tm:5100
pipe_aggregate
    : 'AGGREGATE' opt_with_modifier pipe_selection_item_list_with_order_or_empty opt_group_by_clause_in_pipe
    ;

// googlesql.tm:5113
pipe_group_by
    : 'GROUP'
    ;

// googlesql.tm:5123
pipe_set_operation_base
    : set_operation_metadata parenthesized_query
    | set_operation_metadata table_clause
    | pipe_set_operation_base ',' parenthesized_query
    | pipe_set_operation_base ',' table_clause
    ;

// googlesql.tm:5144
pipe_set_operation
    : pipe_set_operation_base ','?
    ;

// googlesql.tm:5148
subquery_or_subpipeline
    : subpipeline_with_parens
    | parenthesized_query
    ;

// googlesql.tm:5162
pipe_recursive_union_no_alias
    : 'RECURSIVE' set_operation_metadata opt_recursion_depth_modifier subquery_or_subpipeline
    ;

// googlesql.tm:5171
pipe_recursive_union
    : pipe_recursive_union_no_alias as_alias_with_required_as?
    | pipe_recursive_union_no_alias identifier
    ;

// googlesql.tm:5185
pipe_join
    : opt_natural join_type join_hint 'JOIN' hint? table_primary opt_on_or_using_clause
    ;

// googlesql.tm:5218
pipe_call
    : 'CALL' tvf as_alias?
    ;

// googlesql.tm:5225
pipe_window
    : 'WINDOW' pipe_selection_item_list
    ;

// googlesql.tm:5235
pipe_distinct
    : 'DISTINCT'
    ;

// googlesql.tm:5242
pipe_tablesample
    : sample_clause
    ;

// googlesql.tm:5249
pipe_as
    : 'AS' identifier
    ;

// googlesql.tm:5257
pipe_describe
    : 'DESCRIBE'
    ;

// googlesql.tm:5264
pipe_static_describe
    : 'STATIC_DESCRIBE'
    ;

// googlesql.tm:5271
pipe_assert_base
    : 'ASSERT' expression
    | pipe_assert_base ',' expression
    ;

// googlesql.tm:5282
pipe_assert
    : pipe_assert_base ','?
    ;

// googlesql.tm:5289
pipe_log
    : 'LOG' hint? opt_subpipeline
    ;

// googlesql.tm:5296
identifier_in_pipe_drop
    : identifier
    | identifier '.'
    ;

// googlesql.tm:5315
identifier_list_in_pipe_drop
    : identifier_in_pipe_drop
    | identifier_list_in_pipe_drop ',' identifier_in_pipe_drop
    ;

// googlesql.tm:5326
pipe_drop
    : 'DROP' identifier_list_in_pipe_drop ','?
    ;

// googlesql.tm:5333
pipe_set_item
    : identifier '=' expression
    | identifier '.'
    ;

// googlesql.tm:5350
pipe_set_item_list
    : pipe_set_item
    | pipe_set_item_list ',' pipe_set_item
    ;

// googlesql.tm:5361
pipe_set
    : 'SET' pipe_set_item_list ','?
    ;

// googlesql.tm:5368
pipe_pivot
    : pivot_clause as_alias?
    ;

// googlesql.tm:5377
pipe_unpivot
    : unpivot_clause as_alias?
    ;

// googlesql.tm:5386
pipe_if
    : pipe_if_prefix
    | pipe_if_prefix 'ELSE' subpipeline
    ;

// googlesql.tm:5394
pipe_if_prefix
    : 'IF' hint? expression 'THEN' subpipeline
    | pipe_if_prefix pipe_if_elseif
    ;

// googlesql.tm:5406
pipe_if_elseif
    : 'ELSEIF' expression 'THEN' subpipeline
    | 'ELSE' 'IF'
    ;

// googlesql.tm:5419
pipe_elseif
    : 'ELSEIF'
    ;

// googlesql.tm:5430
pipe_else
    : 'ELSE'
    ;

// googlesql.tm:5440
pipe_fork
    : pipe_fork_impl ','?
    ;

// googlesql.tm:5447
pipe_fork_impl
    : 'FORK' hint? subpipeline
    | pipe_fork_impl ',' subpipeline
    ;

// googlesql.tm:5458
pipe_tee
    : 'TEE' hint?
    | pipe_tee_impl ','?
    ;

// googlesql.tm:5469
pipe_tee_impl
    : 'TEE' hint? subpipeline
    | pipe_tee_impl ',' subpipeline
    ;

// googlesql.tm:5480
pipe_with
    : with_clause ','?
    ;

// googlesql.tm:5487
pipe_export_data
    : export_data_no_query
    | export_data_no_query 'AS'
    ;

// googlesql.tm:5499
pipe_create_table
    : create_table_statement_prefix
    | create_table_statement_prefix 'AS'
    ;

// googlesql.tm:5511
pipe_insert
    : insert_statement_in_pipe
    ;

// googlesql.tm:5518
pipe_match_recognize
    : match_recognize_clause
    ;

// googlesql.tm:5525
opt_corresponding_outer_mode
    : (KW_FULL_IN_SET_OP | KW_FULL_IN_SET_OP 'OUTER' | 'OUTER')
    | KW_INNER_IN_SET_OP
    | KW_LEFT_IN_SET_OP 'OUTER'?
    | /* empty */
    ;

// googlesql.tm:5547
opt_strict
    : 'STRICT'
    | /* empty */
    ;

// googlesql.tm:5559
opt_column_match_suffix
    : 'CORRESPONDING'
    | 'CORRESPONDING' 'BY' column_list
    | 'BY' 'NAME'
    | 'BY' 'NAME' 'ON' column_list
    | /* empty */
    ;

// googlesql.tm:5601
query_set_operation_prefix
    : query_primary set_operation_metadata query_primary
    | query_set_operation_prefix set_operation_metadata query_primary
    | query_primary set_operation_metadata 'FROM'
    | query_set_operation_prefix set_operation_metadata 'FROM'
    ;

// googlesql.tm:5634
set_operation_metadata
    : opt_corresponding_outer_mode query_set_operation_type hint? all_or_distinct opt_strict opt_column_match_suffix
    ;

// googlesql.tm:5665
query_set_operation
    : query_set_operation_prefix
    ;

// googlesql.tm:5672
query_primary
    : select
    | table_clause_reserved
    | parenthesized_query as_alias_with_required_as?
    ;

// googlesql.tm:5696
select_clause
    : 'SELECT' hint? opt_with_modifier opt_all_or_distinct opt_select_as_clause select_list
    | 'SELECT' hint? opt_with_modifier opt_all_or_distinct opt_select_as_clause 'FROM'
    ;

// googlesql.tm:5717
select
    : select_clause opt_from_clause opt_clauses_following_from
    ;

// googlesql.tm:5726
pre_with_modifier
    : /* empty */
    ;

// googlesql.tm:5733
opt_with_modifier
    : pre_with_modifier 'WITH' identifier
    | pre_with_modifier 'WITH' identifier KW_OPTIONS_IN_WITH_OPTIONS options_list
    | pre_with_modifier
    ;

// googlesql.tm:5750
opt_select_as_clause
    : 'AS' 'STRUCT'
    | 'AS' path_expression
    | /* empty */
    ;

// googlesql.tm:5787
extra_identifier_in_hints_name
    : 'HASH'
    | 'PROTO'
    | 'PARTITION'
    ;

// googlesql.tm:5793
identifier_in_hints
    : identifier
    | extra_identifier_in_hints_name
    ;

// googlesql.tm:5801
hint_entry
    : identifier_in_hints '=' expression
    | identifier_in_hints '.' identifier_in_hints '=' expression
    ;

// googlesql.tm:5812
hint_with_body_prefix
    : OPEN_INTEGER_PREFIX_HINT integer_literal KW_OPEN_HINT '{' hint_entry
    | KW_OPEN_HINT '{' hint_entry
    | hint_with_body_prefix ',' hint_entry
    ;

// googlesql.tm:5827
hint_with_body
    : hint_with_body_prefix '}'
    ;

// googlesql.tm:5837
hint
    : KW_OPEN_INTEGER_HINT integer_literal
    | hint_with_body
    ;

// googlesql.tm:5853
opt_all_or_distinct
    : 'ALL'
    | 'DISTINCT'
    | /* empty */
    ;

// googlesql.tm:5868
select_list_prefix
    : select_column
    | select_list_prefix ',' select_column
    ;

// googlesql.tm:5879
select_list
    : select_list_prefix
    | select_list_prefix ','
    ;

// googlesql.tm:5890
star_except_list_prefix
    : 'EXCEPT' '(' identifier
    | star_except_list_prefix ',' identifier
    ;

// googlesql.tm:5901
star_except_list
    : star_except_list_prefix ')'
    ;

// googlesql.tm:5908
star_replace_item
    : expression 'AS' identifier
    ;

// googlesql.tm:5915
star_modifiers_with_replace_prefix
    : star_except_list 'REPLACE' '(' star_replace_item
    | 'REPLACE' '(' star_replace_item
    | star_modifiers_with_replace_prefix ',' star_replace_item
    ;

// googlesql.tm:5930
star_modifiers
    : star_except_list
    | star_modifiers_with_replace_prefix ')'
    ;

// googlesql.tm:5941
select_column
    : select_column_expr
    | select_column_dot_star
    | select_column_star
    ;

// googlesql.tm:5948
select_column_expr
    : expression
    | select_column_expr_with_as_alias
    | expression identifier
    | expression ATTACHED_ALIAS
    ;

// googlesql.tm:5969
select_list_prefix_with_as_aliases
    : select_column_expr_with_as_alias
    | select_list_prefix_with_as_aliases ',' select_column_expr_with_as_alias
    ;

// googlesql.tm:5980
select_column_expr_with_as_alias
    : expression 'AS' identifier
    ;

// googlesql.tm:5990
select_column_dot_star
    : expression_higher_prec_than_and '.' '*'
    | expression_higher_prec_than_and '.' '*' star_modifiers
    ;

// googlesql.tm:6012
select_column_star
    : '*'
    | '*' star_modifiers
    ;

// googlesql.tm:6026
as_alias
    : 'AS'? identifier
    ;

// googlesql.tm:6033
as_alias_with_required_as
    : 'AS' identifier
    ;

// googlesql.tm:6040
opt_as_or_into_alias
    : 'AS' identifier
    | 'INTO' identifier
    | /* empty */
    ;

// googlesql.tm:6056
opt_natural
    : 'NATURAL'
    | /* empty */
    ;

// googlesql.tm:6067
opt_int_literal_or_parameter
    : int_literal_or_parameter
    | /* empty */
    ;

// googlesql.tm:6075
int_literal_or_parameter
    : integer_literal
    | parameter_expression
    | system_variable_expression
    ;

// googlesql.tm:6081
cast_int_literal_or_parameter
    : 'CAST' '(' int_literal_or_parameter 'AS' type opt_format ')'
    ;

// googlesql.tm:6090
possibly_cast_int_literal_or_parameter
    : cast_int_literal_or_parameter
    | int_literal_or_parameter
    ;

// googlesql.tm:6095
repeatable_clause
    : 'REPEATABLE' '(' possibly_cast_int_literal_or_parameter ')'
    ;

// googlesql.tm:6102
sample_size_value
    : possibly_cast_int_literal_or_parameter
    | floating_point_literal
    ;

// googlesql.tm:6108
sample_size_unit
    : 'ROWS'
    | 'PERCENT'
    ;

// googlesql.tm:6119
sample_size
    : sample_size_value sample_size_unit opt_partition_by_clause_no_hint
    ;

// googlesql.tm:6128
opt_repeatable_clause
    : repeatable_clause
    | /* empty */
    ;

// googlesql.tm:6138
opt_sample_clause_suffix
    : repeatable_clause
    | 'WITH' 'WEIGHT' opt_repeatable_clause
    | 'WITH' 'WEIGHT' identifier opt_repeatable_clause
    | 'WITH' 'WEIGHT' 'AS' identifier opt_repeatable_clause
    | /* empty */
    ;

// googlesql.tm:6169
opt_graph_sample_clause_suffix
    : repeatable_clause
    | 'WITH' 'WEIGHT' opt_repeatable_clause
    | 'WITH' 'WEIGHT' 'AS' identifier opt_repeatable_clause
    | /* empty */
    ;

// googlesql.tm:6191
sample_clause
    : 'TABLESAMPLE' identifier '(' sample_size ')' opt_sample_clause_suffix
    ;

// googlesql.tm:6198
graph_sample_clause
    : 'TABLESAMPLE' identifier '(' sample_size ')' opt_graph_sample_clause_suffix
    ;

// googlesql.tm:6205
pivot_expression
    : expression as_alias?
    ;

// googlesql.tm:6212
pivot_expression_list
    : pivot_expression
    | pivot_expression_list ',' pivot_expression
    ;

// googlesql.tm:6223
pivot_value
    : expression as_alias?
    ;

// googlesql.tm:6230
pivot_value_list
    : pivot_value
    | pivot_value_list ',' pivot_value
    ;

// googlesql.tm:6241
pivot_clause
    : 'PIVOT' '(' pivot_expression_list 'FOR' expression_higher_prec_than_and 'IN' '(' pivot_value_list ')' ')'
    ;

// googlesql.tm:6253
opt_as_string_or_integer
    : 'AS'? string_literal
    | 'AS'? integer_literal
    | /* empty */
    ;

// googlesql.tm:6268
path_expression_list
    : path_expression
    | path_expression_list ',' path_expression
    ;

// googlesql.tm:6279
path_expression_list_with_opt_parens
    : '(' path_expression_list ')'
    | path_expression
    ;

// googlesql.tm:6290
path_expression_list_prefix
    : '(' path_expression
    | path_expression_list_prefix ',' path_expression
    ;

// googlesql.tm:6301
path_expression_list_with_parens
    : path_expression_list_prefix ')'
    ;

// googlesql.tm:6308
opt_path_expression_list_with_parens
    : path_expression_list_with_parens
    | /* empty */
    ;

// googlesql.tm:6316
unpivot_in_item
    : path_expression_list_with_opt_parens opt_as_string_or_integer
    ;

// googlesql.tm:6323
unpivot_in_item_list_prefix
    : '(' unpivot_in_item
    | unpivot_in_item_list_prefix ',' unpivot_in_item
    ;

// googlesql.tm:6334
unpivot_in_item_list
    : unpivot_in_item_list_prefix ')'
    ;

// googlesql.tm:6341
opt_unpivot_nulls_filter
    : 'EXCLUDE' 'NULLS'
    | 'INCLUDE' 'NULLS'
    | /* empty */
    ;

// googlesql.tm:6356
unpivot_clause
    : 'UNPIVOT' opt_unpivot_nulls_filter '(' path_expression_list_with_opt_parens 'FOR' path_expression 'IN' unpivot_in_item_list ')'
    ;

// googlesql.tm:6367
pivot_or_unpivot_clause
    : pivot_clause as_alias?
    | unpivot_clause as_alias?
    | qualify_clause_nonreserved
    ;

// googlesql.tm:6384
match_recognize_clause
    : KW_MATCH_RECOGNIZE_RESERVED '(' opt_partition_by_clause order_by_clause measures_clause opt_after_match_skip_clause 'PATTERN' '(' row_pattern_expr ')' 'DEFINE' with_expression_variable_prefix opt_options_list ')' as_alias?
    ;

// googlesql.tm:6402
measures_clause
    : 'MEASURES' select_list_prefix_with_as_aliases
    ;

// googlesql.tm:6409
opt_after_match_skip_clause
    : /* empty */
    | 'AFTER' 'MATCH' 'SKIP' skip_to_target
    ;

// googlesql.tm:6420
skip_to_target
    : 'PAST' 'LAST' 'ROW'
    | 'TO' 'NEXT' 'ROW'
    ;

// googlesql.tm:6435
row_pattern_expr
    : row_pattern_concatenation_or_empty
    | row_pattern_expr '|' row_pattern_concatenation_or_empty
    | row_pattern_expr '||' row_pattern_concatenation_or_empty
    ;

// googlesql.tm:6461
row_pattern_concatenation_or_empty
    : row_pattern_concatenation
    | /* empty */
    ;

// googlesql.tm:6471
row_pattern_concatenation
    : row_pattern_factor
    | row_pattern_concatenation row_pattern_factor
    ;

// googlesql.tm:6480
row_pattern_factor
    : row_pattern_primary
    | row_pattern_anchor
    | quantified_row_pattern
    ;

// googlesql.tm:6486
row_pattern_anchor
    : '^'
    | '$'
    ;

// googlesql.tm:6501
quantified_row_pattern
    : row_pattern_primary row_pattern_quantifier
    ;

// googlesql.tm:6508
row_pattern_primary
    : identifier
    | '(' row_pattern_expr ')'
    ;

// googlesql.tm:6523
row_pattern_quantifier
    : potentially_reluctant_quantifier
    | potentially_reluctant_quantifier '?'
    | '{' int_literal_or_parameter '}'
    ;

// googlesql.tm:6540
plus_or_star_quantifier
    : '*'
    | '+'
    ;

// googlesql.tm:6554
potentially_reluctant_quantifier
    : '?'
    | plus_or_star_quantifier
    | '{' opt_int_literal_or_parameter ',' opt_int_literal_or_parameter '}'
    ;

// googlesql.tm:6574
table_subquery
    : parenthesized_query as_alias? pivot_or_unpivot_clause?
    ;

// googlesql.tm:6588
table_clause_no_keyword
    : path_expression opt_where_clause
    | tvf as_alias? pivot_or_unpivot_clause? opt_where_clause
    ;

// googlesql.tm:6599
table_clause_reserved
    : KW_TABLE_FOR_TABLE_CLAUSE table_clause_no_keyword
    ;

// googlesql.tm:6606
table_clause_unreserved
    : 'TABLE' table_clause_no_keyword
    ;

// googlesql.tm:6613
table_clause
    : table_clause_reserved
    | table_clause_unreserved
    ;

// googlesql.tm:6618
model_clause
    : 'MODEL' path_expression
    ;

// googlesql.tm:6625
connection_clause
    : 'CONNECTION' path_expression_or_default
    ;

// googlesql.tm:6632
descriptor_column
    : identifier
    ;

// googlesql.tm:6639
descriptor_column_list
    : descriptor_column
    | descriptor_column_list ',' descriptor_column
    ;

// googlesql.tm:6650
descriptor_argument
    : 'DESCRIPTOR' '(' descriptor_column_list ')'
    ;

// googlesql.tm:6657
input_table_argument
    : 'INPUT' 'TABLE'
    ;

// googlesql.tm:6664
named_non_expr_argument
    : identifier '=>' table_clause
    ;

// googlesql.tm:6673
tvf_argument
    : expression
    | descriptor_argument
    | table_clause
    | model_clause
    | connection_clause
    | named_argument
    | named_non_expr_argument
    | input_table_argument
    | '(' model_clause ')'
    | '(' connection_clause ')'
    | '(' named_argument ')'
    | 'SELECT'
    | 'WITH'
    ;

// googlesql.tm:6750
tvf
    : (path_expression | 'IF') '(' ((tvf_argument) (',' (tvf_argument))*)? ')' hint?
    ;

// googlesql.tm:6762
table_path_expression_base
    : unnest_expression
    | maybe_slashed_or_dashed_path_expression
    | path_expression '['
    | path_expression '.' '('
    | unnest_expression '['
    | unnest_expression '.' '('
    ;

// googlesql.tm:6798
table_path_expression
    : table_path_expression_base hint? as_alias? (with_offset as_alias?)? pivot_or_unpivot_clause? at_system_time?
    ;

// googlesql.tm:6818
table_primary
    : 'LATERAL'? tvf as_alias? pivot_or_unpivot_clause?
    | table_path_expression
    | '(' join ')'
    | table_subquery
    | 'LATERAL' table_subquery
    | graph_table_query
    | table_primary match_recognize_clause
    | table_primary sample_clause
    ;

// googlesql.tm:6860
gql_query
    : 'GRAPH' path_expression graph_operation_block
    ;

// googlesql.tm:6870
gql_statement
    : gql_query
    ;

// googlesql.tm:6877
graph_table_query
    : KW_GRAPH_TABLE_RESERVED '(' path_expression graph_match_operator graph_shape_clause ')' as_alias?
    | KW_GRAPH_TABLE_RESERVED '(' path_expression graph_operation_block ')' as_alias?
    ;

// googlesql.tm:6897
graph_shape_clause
    : 'COLUMNS' '(' select_list ')'
    | /* empty */
    ;

// googlesql.tm:6908
graph_return_item
    : expression
    | expression 'AS' identifier
    | '*'
    ;

// googlesql.tm:6926
graph_return_item_list
    : graph_return_item
    | graph_return_item_list ',' graph_return_item
    ;

// googlesql.tm:6937
graph_return_operator
    : 'RETURN' hint? opt_all_or_distinct graph_return_item_list opt_group_by_clause opt_graph_order_by_clause graph_offset_clause? graph_limit_clause?
    ;

// googlesql.tm:6971
opt_graph_asc_or_desc
    : opt_asc_or_desc
    | 'ASCENDING'
    | 'DESCENDING'
    ;

// googlesql.tm:6983
graph_ordering_expression
    : expression opt_collate_clause opt_graph_asc_or_desc opt_null_order
    ;

// googlesql.tm:6994
graph_order_by_clause_prefix
    : 'ORDER' hint? 'BY' graph_ordering_expression
    | graph_order_by_clause_prefix ',' graph_ordering_expression
    ;

// googlesql.tm:7006
graph_order_by_clause
    : graph_order_by_clause_prefix
    ;

// googlesql.tm:7013
opt_graph_order_by_clause
    : graph_order_by_clause
    | /* empty */
    ;

// googlesql.tm:7024
graph_order_by_operator
    : graph_order_by_clause
    ;

// googlesql.tm:7031
graph_limit_clause
    : 'LIMIT' possibly_cast_int_literal_or_parameter
    ;

// googlesql.tm:7038
graph_offset_clause
    : ('OFFSET' | 'SKIP') possibly_cast_int_literal_or_parameter
    ;

// googlesql.tm:7045
graph_match_operator
    : 'MATCH' hint? graph_pattern
    ;

// googlesql.tm:7052
graph_optional_match_operator
    : 'OPTIONAL' 'MATCH' hint? graph_pattern
    ;

// googlesql.tm:7061
graph_let_operator
    : 'LET' graph_let_variable_definition_list
    ;

// googlesql.tm:7068
graph_let_variable_definition_list
    : graph_let_variable_definition
    | graph_let_variable_definition_list ',' graph_let_variable_definition
    ;

// googlesql.tm:7079
graph_let_variable_definition
    : identifier '=' expression
    ;

// googlesql.tm:7086
graph_filter_operator
    : 'FILTER' where_clause
    | 'FILTER' expression
    ;

// googlesql.tm:7098
graph_with_operator
    : 'WITH' opt_all_or_distinct hint? graph_return_item_list opt_group_by_clause
    ;

// googlesql.tm:7113
graph_for_operator
    : 'FOR' identifier 'IN' expression (with_offset as_alias_with_required_as?)?
    ;

// googlesql.tm:7127
graph_sample_operator
    : graph_sample_clause
    ;

// googlesql.tm:7134
graph_call_operator
    : graph_call_operator_core
    | 'OPTIONAL' graph_call_operator_core
    ;

// googlesql.tm:7144
graph_call_operator_core
    : 'CALL' opt_per_clause tvf opt_yield_clause
    | 'CALL' opt_per_clause braced_graph_subquery
    | 'CALL' parenthesized_identifier_list braced_graph_subquery
    ;

// googlesql.tm:7162
parenthesized_identifier_list
    : '(' identifier_list ')'
    | '(' ')'
    ;

// googlesql.tm:7173
opt_per_clause
    : 'PER' parenthesized_identifier_list
    | /* empty */
    ;

// googlesql.tm:7184
opt_yield_clause
    : 'YIELD' yield_item_list
    | /* empty */
    ;

// googlesql.tm:7195
yield_item_list
    : yield_item
    | yield_item_list ',' yield_item
    ;

// googlesql.tm:7206
yield_item
    : identifier as_alias_with_required_as?
    ;

// googlesql.tm:7219
graph_linear_op
    : graph_match_operator
    | graph_optional_match_operator
    | graph_let_operator
    | graph_filter_operator
    | graph_order_by_operator
    | graph_offset_clause
    | graph_limit_clause
    | graph_with_operator
    | graph_for_operator
    | graph_sample_operator
    | graph_call_operator
    ;

// googlesql.tm:7242
graph_linear_operator_list
    : graph_linear_op+
    ;

// googlesql.tm:7272
graph_linear_query_operation
    : graph_return_operator
    | graph_linear_operator_list graph_return_operator
    ;

// googlesql.tm:7283
graph_set_operation_metadata
    : opt_corresponding_outer_mode query_set_operation_type all_or_distinct
    ;

// googlesql.tm:7293
graph_composite_query_prefix
    : graph_linear_query_operation graph_set_operation_metadata graph_linear_query_operation
    | graph_composite_query_prefix graph_set_operation_metadata graph_linear_query_operation
    ;

// googlesql.tm:7310
graph_composite_query_block
    : graph_linear_query_operation
    | graph_composite_query_prefix
    ;

// googlesql.tm:7315
graph_operation_block
    : graph_composite_query_block
    | graph_operation_block 'NEXT' graph_composite_query_block
    ;

// googlesql.tm:7327
graph_pattern
    : graph_path_pattern_list opt_where_clause
    ;

// googlesql.tm:7336
graph_path_pattern_list
    : graph_path_pattern
    | graph_path_pattern_list ',' hint? graph_path_pattern
    ;

// googlesql.tm:7350
graph_search_prefix
    : 'ANY'
    | 'ANY' 'SHORTEST'
    | 'ANY' 'CHEAPEST'
    | 'ANY' int_literal_or_parameter
    | 'SHORTEST' int_literal_or_parameter
    | 'CHEAPEST' int_literal_or_parameter
    | 'ALL'
    | 'ALL' 'SHORTEST'
    | 'ALL' 'CHEAPEST'
    ;

// googlesql.tm:7401
graph_path_mode
    : 'WALK'
    | 'TRAIL'
    | 'SIMPLE'
    | 'ACYCLIC'
    ;

// googlesql.tm:7424
graph_path_pattern
    : (graph_identifier '=')? graph_search_prefix? graph_path_mode? ('PATH' | 'PATHS')? graph_path_pattern_expr
    ;

// googlesql.tm:7443
graph_path_pattern_expr
    : graph_path_factor
    | graph_path_pattern_expr hint? graph_path_factor
    ;

// googlesql.tm:7512
graph_path_factor
    : graph_path_primary
    | graph_quantified_path_primary
    ;

// googlesql.tm:7517
graph_path_primary
    : graph_element_pattern
    | graph_parenthesized_path_pattern
    ;

// googlesql.tm:7528
graph_parenthesized_path_pattern
    : '(' hint? graph_path_pattern opt_where_clause ')'
    ;

// googlesql.tm:7554
graph_quantifier
    : '{' opt_int_literal_or_parameter ',' opt_int_literal_or_parameter '}'
    | '{' int_literal_or_parameter '}'
    | plus_or_star_quantifier
    ;

// googlesql.tm:7575
graph_quantified_path_primary
    : graph_path_primary graph_quantifier
    ;

// googlesql.tm:7597
graph_element_pattern
    : graph_node_pattern
    | graph_edge_pattern
    ;

// googlesql.tm:7602
graph_identifier
    : token_identifier
    | common_keyword_as_identifier
    ;

// googlesql.tm:7610
opt_graph_element_identifier
    : graph_identifier
    | /* empty */
    ;

// googlesql.tm:7618
opt_graph_cost
    : 'COST' expression
    | /* empty */
    ;

// googlesql.tm:7629
graph_element_pattern_filler
    : hint? opt_graph_element_identifier is_label_expression? opt_graph_property_specification opt_where_clause opt_graph_cost
    ;

// googlesql.tm:7642
graph_property_specification
    : graph_property_specification_prefix '}'
    ;

// googlesql.tm:7649
opt_graph_property_specification
    : graph_property_specification
    | /* empty */
    ;

// googlesql.tm:7657
graph_property_specification_prefix
    : '{' graph_property_name_and_value
    | graph_property_specification_prefix ',' graph_property_name_and_value
    ;

// googlesql.tm:7668
graph_property_name_and_value
    : identifier ':' expression
    ;

// googlesql.tm:7675
graph_node_pattern
    : '(' graph_element_pattern_filler ')'
    ;

// googlesql.tm:7684
graph_edge_pattern
    : '-' '[' graph_element_pattern_filler ']' '-'
    | '<' '-' '[' graph_element_pattern_filler ']' '-'
    | '-' '[' graph_element_pattern_filler ']' '->'
    | '-'
    | '<' '-'
    | '->'
    ;

// googlesql.tm:7727
is_label_expression
    : ('IS' | ':') label_expression
    ;

// googlesql.tm:7734
label_expression
    : label_primary
    | label_expression '&' label_expression
    | label_expression '|' label_expression
    | '!' label_expression
    ;

// googlesql.tm:7754
label_primary
    : identifier
    | '%'
    | parenthesized_label_expression
    ;

// googlesql.tm:7766
parenthesized_label_expression
    : '(' label_expression ')'
    ;

// googlesql.tm:7777
graph_expression
    : expression_higher_prec_than_and edge_source_endpoint_operator expression_higher_prec_than_and
    | expression_higher_prec_than_and edge_dest_endpoint_operator expression_higher_prec_than_and
    | expression_higher_prec_than_and is_labeled_operator label_expression
    | expression_higher_prec_than_and in_operator braced_graph_subquery
    ;

// googlesql.tm:7816
edge_source_endpoint_operator
    : 'IS' 'SOURCE'
    | 'IS' 'SOURCE' 'OF'
    | 'IS' 'NOT' 'SOURCE'
    | 'IS' 'NOT' 'SOURCE' 'OF'
    ;

// googlesql.tm:7841
edge_dest_endpoint_operator
    : 'IS' 'DESTINATION'
    | 'IS' 'DESTINATION' 'OF'
    | 'IS' 'NOT' 'DESTINATION'
    | 'IS' 'NOT' 'DESTINATION' 'OF'
    ;

// googlesql.tm:7866
is_labeled_operator
    : 'IS' 'LABELED'
    | 'IS' 'NOT' 'LABELED'
    ;

// googlesql.tm:7880
at_system_time
    : 'FOR' ('SYSTEM' 'TIME' | 'SYSTEM_TIME') 'AS' 'OF' expression
    ;

// googlesql.tm:7887
on_clause
    : 'ON' expression
    ;

// googlesql.tm:7894
using_clause_prefix
    : 'USING' '(' identifier
    | using_clause_prefix ',' identifier
    ;

// googlesql.tm:7905
using_clause
    : using_clause_prefix ')'
    ;

// googlesql.tm:7912
on_or_using_clause_lists
    : on_or_using_clause+
    ;

// googlesql.tm:7930
on_or_using_clause
    : on_clause
    | using_clause
    ;

// googlesql.tm:7935
opt_on_or_using_clause
    : on_or_using_clause?
    ;

// googlesql.tm:7943
join_type_base
    : 'CROSS'
    | 'FULL' 'OUTER'?
    | 'INNER'
    | 'LEFT' 'OUTER'?
    | 'RIGHT' 'OUTER'?
    ;

// googlesql.tm:7966
join_type
    : join_type_base?
    ;

// googlesql.tm:7974
join_hint_base
    : 'HASH'
    | 'LOOKUP'
    ;

// googlesql.tm:7985
join_hint
    : join_hint_base?
    ;

// googlesql.tm:7992
join_input
    : join
    | table_primary
    ;

// googlesql.tm:8004
join
    : join_input opt_natural join_type join_hint 'JOIN' hint? table_primary on_or_using_clause_lists?
    ;

// googlesql.tm:8026
from_clause_contents
    : table_primary
    | from_clause_contents ',' table_primary
    | from_clause_contents opt_natural join_type join_hint 'JOIN' hint? table_primary on_or_using_clause_lists?
    | '@'
    | '?'
    | '@@'
    ;

// googlesql.tm:8107
opt_from_clause
    : from_clause?
    ;

// googlesql.tm:8114
from_clause
    : 'FROM' from_clause_contents
    ;

// googlesql.tm:8138
opt_clauses_following_from
    : where_clause opt_group_by_clause opt_having_clause qualify_clause? window_clause?
    | opt_clauses_following_where
    ;

// googlesql.tm:8151
opt_clauses_following_where
    : group_by_clause opt_having_clause qualify_clause? window_clause?
    | opt_clauses_following_group_by
    ;

// googlesql.tm:8165
opt_clauses_following_group_by
    : having_clause qualify_clause? window_clause?
    | qualify_clause_reserved? window_clause?
    ;

// googlesql.tm:8178
where_clause
    : 'WHERE' expression
    ;

// googlesql.tm:8185
opt_where_clause
    : where_clause?
    ;

// googlesql.tm:8192
rollup_list
    : 'ROLLUP' '(' ((expression) (',' (expression))*) ')'
    ;

// googlesql.tm:8199
cube_list
    : 'CUBE' '(' ((expression) (',' (expression))*) ')'
    ;

// googlesql.tm:8206
grouping_set
    : '(' ')'
    | expression
    | rollup_list
    | cube_list
    ;

// googlesql.tm:8238
selection_item_order
    : asc_or_desc opt_null_order
    ;

// googlesql.tm:8247
opt_selection_item_order
    : selection_item_order?
    ;

// googlesql.tm:8255
opt_grouping_item_order
    : opt_selection_item_order
    | null_order
    ;

// googlesql.tm:8264
grouping_item_base
    : '(' ')'
    | rollup_list
    | cube_list
    | 'GROUPING' 'SETS' '(' ((grouping_set) (',' (grouping_set))*) ')'
    ;

// googlesql.tm:8288
grouping_item
    : grouping_item_base
    | expression
    ;

// googlesql.tm:8298
grouping_item_in_pipe
    : grouping_item_base
    | expression as_alias? opt_grouping_item_order
    ;

// googlesql.tm:8306
group_by_preamble
    : 'GROUP' hint? 'BY'
    ;

// googlesql.tm:8314
group_by_preamble_in_pipe
    : 'GROUP' hint? ('AND' 'ORDER')? 'BY'
    ;

// googlesql.tm:8327
group_by_clause_prefix
    : group_by_preamble grouping_item
    | group_by_clause_prefix ',' grouping_item
    ;

// googlesql.tm:8338
group_by_clause_prefix_in_pipe
    : group_by_preamble_in_pipe grouping_item_in_pipe
    | group_by_clause_prefix_in_pipe ',' grouping_item_in_pipe
    ;

// googlesql.tm:8351
group_by_all
    : group_by_preamble 'ALL'
    ;

// googlesql.tm:8361
group_by_clause
    : group_by_all
    | group_by_clause_prefix
    ;

// googlesql.tm:8367
opt_group_by_clause
    : group_by_clause?
    ;

// googlesql.tm:8386
group_by_clause_in_pipe
    : group_by_clause_prefix_in_pipe ','?
    ;

// googlesql.tm:8393
opt_group_by_clause_in_pipe
    : group_by_clause_in_pipe?
    ;

// googlesql.tm:8400
having_clause
    : 'HAVING' expression
    ;

// googlesql.tm:8407
opt_having_clause
    : having_clause?
    ;

// googlesql.tm:8414
window_definition
    : identifier 'AS' window_specification
    ;

// googlesql.tm:8421
window_clause
    : 'WINDOW' window_definition
    | window_clause ',' window_definition
    ;

// googlesql.tm:8432
window_clause_with_trailing_comma
    : window_clause ','?
    ;

// googlesql.tm:8439
opt_window_clause_with_trailing_comma
    : window_clause_with_trailing_comma?
    ;

// googlesql.tm:8446
qualify_clause
    : qualify_clause_reserved
    | qualify_clause_nonreserved
    ;

// googlesql.tm:8457
qualify_clause_reserved
    : KW_QUALIFY_RESERVED expression
    ;

// googlesql.tm:8467
qualify_clause_nonreserved
    : 'QUALIFY' expression
    ;

// googlesql.tm:8477
limit_offset_clause
    : limit_expression 'OFFSET' expression
    | limit_expression
    | limit_all 'OFFSET' expression
    | limit_all
    ;

// googlesql.tm:8497
limit_expression
    : 'LIMIT' expression
    ;

// googlesql.tm:8504
limit_all
    : 'LIMIT' 'ALL'
    ;

// googlesql.tm:8512
opt_limit_offset_clause
    : limit_offset_clause?
    ;

// googlesql.tm:8519
lock_mode_clause
    : KW_FOR_BEFORE_LOCK_MODE 'UPDATE'
    ;

// googlesql.tm:8531
opt_lock_mode_clause
    : lock_mode_clause?
    ;

// googlesql.tm:8538
having_modifier
    : 'HAVING' ('MAX' | 'MIN') expression
    ;

// googlesql.tm:8549
clamped_between_modifier
    : 'CLAMPED' 'BETWEEN' expression_higher_prec_than_and 'AND' expression
    ;

// googlesql.tm:8558
with_report_modifier
    : 'WITH' 'REPORT' options_list?
    ;

// googlesql.tm:8565
null_handling_modifier
    : (('IGNORE' | 'RESPECT') 'NULLS')?
    ;

// googlesql.tm:8578
possibly_unbounded_int_literal_or_parameter
    : int_literal_or_parameter
    | 'UNBOUNDED'
    ;

// googlesql.tm:8589
recursion_depth_modifier
    : 'WITH' 'DEPTH' as_alias_with_required_as?
    | 'WITH' 'DEPTH' as_alias_with_required_as? 'BETWEEN' possibly_unbounded_int_literal_or_parameter 'AND' possibly_unbounded_int_literal_or_parameter
    | 'WITH' 'DEPTH' as_alias_with_required_as? 'MAX' possibly_unbounded_int_literal_or_parameter
    ;

// googlesql.tm:8626
opt_recursion_depth_modifier
    : recursion_depth_modifier?
    ;

// googlesql.tm:8642
aliased_query_with_overridden_next_token_lookback
    : '(' query ')'
    ;

// googlesql.tm:8654
aliased_query
    : identifier 'AS' aliased_query_with_overridden_next_token_lookback recursion_depth_modifier?
    ;

// googlesql.tm:8666
aliased_query_list
    : aliased_query
    | aliased_query_list ',' aliased_query
    ;

// googlesql.tm:8677
with_clause_entry
    : aliased_query
    | identifier '(' ')' 'AS' 'GROUP' 'ROWS'
    ;

// googlesql.tm:8693
with_clause
    : 'WITH' 'RECURSIVE'? with_clause_entry
    | with_clause ',' with_clause_entry
    ;

// googlesql.tm:8705
opt_with_clause
    : with_clause?
    ;

// googlesql.tm:8712
opt_with_connection_clause
    : with_connection_clause?
    ;

// googlesql.tm:8719
with_clause_with_trailing_comma
    : with_clause ','
    ;

// googlesql.tm:8726
asc_or_desc
    : 'ASC'
    | 'DESC'
    ;

// googlesql.tm:8737
opt_asc_or_desc
    : asc_or_desc?
    ;

// googlesql.tm:8744
null_order
    : 'NULLS' 'FIRST'
    | 'NULLS' 'LAST'
    ;

// googlesql.tm:8759
opt_null_order
    : null_order?
    ;

// googlesql.tm:8766
string_literal_or_parameter
    : string_literal
    | parameter_expression
    | system_variable_expression
    ;

// googlesql.tm:8775
collate_clause
    : 'COLLATE' string_literal_or_parameter
    ;

// googlesql.tm:8782
opt_collate_clause
    : collate_clause?
    ;

// googlesql.tm:8789
default_collate_clause
    : 'DEFAULT' collate_clause
    ;

// googlesql.tm:8796
opt_default_collate_clause
    : default_collate_clause?
    ;

// googlesql.tm:8803
ordering_expression
    : expression collate_clause? opt_asc_or_desc null_order?
    ;

// googlesql.tm:8812
order_by_clause_prefix
    : 'ORDER' hint? 'BY' ordering_expression
    | order_by_clause_prefix ',' ordering_expression
    ;

// googlesql.tm:8823
order_by_clause
    : order_by_clause_prefix
    ;

// googlesql.tm:8830
order_by_clause_with_opt_comma
    : order_by_clause_prefix ','?
    ;

// googlesql.tm:8837
opt_order_by_clause
    : order_by_clause?
    ;

// googlesql.tm:8844
parenthesized_in_rhs
    : parenthesized_query
    | '(' expression_maybe_parenthesized_not_a_query ')'
    | in_list_two_or_more_prefix ')'
    ;

// googlesql.tm:8862
parenthesized_anysomeall_list_rhs
    : parenthesized_query
    | '(' expression_maybe_parenthesized_not_a_query ')'
    | in_list_two_or_more_prefix ')'
    ;

// googlesql.tm:8903
in_list_two_or_more_prefix
    : '(' expression ',' expression
    | in_list_two_or_more_prefix ',' expression
    ;

// googlesql.tm:8916
expression_with_opt_alias
    : expression as_alias_with_required_as?
    ;

// googlesql.tm:8923
unnest_expression_prefix
    : 'UNNEST' '(' expression_with_opt_alias
    | unnest_expression_prefix ',' expression_with_opt_alias
    ;

// googlesql.tm:8934
unnest_expression
    : unnest_expression_prefix (',' named_argument)? ')'
    | 'UNNEST' '(' 'SELECT'
    ;

// googlesql.tm:8949
unnest_expression_with_opt_alias_and_offset
    : unnest_expression as_alias? with_offset_and_alias?
    ;

// googlesql.tm:8957
comparative_operator
    : '='
    | '!='
    | '<>'
    | '<'
    | '<='
    | '>'
    | '>='
    ;

// googlesql.tm:8988
additive_operator
    : '+'
    | '-'
    ;

// googlesql.tm:8999
multiplicative_operator
    : '*'
    | '/'
    ;

// googlesql.tm:9011
shift_operator
    : '<<'
    | '>>'
    ;

// googlesql.tm:9023
import_type
    : 'MODULE'
    | 'PROTO'
    ;

// googlesql.tm:9035
any_some_all
    : 'ANY'
    | 'SOME'
    | 'ALL'
    ;

// googlesql.tm:9057
like_operator
    : 'LIKE'
    | KW_NOT_SPECIAL 'LIKE'
    ;

// googlesql.tm:9072
between_operator
    : 'BETWEEN'
    | KW_NOT_SPECIAL 'BETWEEN'
    ;

// googlesql.tm:9086
distinct_operator
    : 'IS' 'DISTINCT' 'FROM'
    | 'IS' KW_NOT_SPECIAL 'DISTINCT' 'FROM'
    ;

// googlesql.tm:9101
in_operator
    : 'IN'
    | KW_NOT_SPECIAL 'IN'
    ;

// googlesql.tm:9116
is_operator
    : 'IS'
    | 'IS' 'NOT'
    ;

// googlesql.tm:9129
unary_operator
    : '+'
    | '-'
    | '~'
    ;

// googlesql.tm:9147
with_expression_variable
    : identifier 'AS' expression
    ;

// googlesql.tm:9155
with_expression_variable_prefix
    : with_expression_variable
    | with_expression_variable_prefix ',' with_expression_variable
    ;

// googlesql.tm:9166
with_expression
    : KW_WITH_STARTING_WITH_EXPRESSION '(' with_expression_variable_prefix ',' expression ')'
    ;

// googlesql.tm:9179
expression
    : expression_higher_prec_than_and
    | and_expression
    | or_expression
    ;

// googlesql.tm:9185
or_expression
    : expression 'OR' expression
    ;

// googlesql.tm:9198
and_expression
    : and_expression 'AND' expression_higher_prec_than_and
    | expression_higher_prec_than_and 'AND' expression_higher_prec_than_and
    ;

// googlesql.tm:9215
expression_higher_prec_than_and
    : unparenthesized_expression_higher_prec_than_and
    | parenthesized_expression_not_a_query
    | parenthesized_query
    ;

// googlesql.tm:9227
expression_maybe_parenthesized_not_a_query
    : parenthesized_expression_not_a_query
    | unparenthesized_expression_higher_prec_than_and
    | and_expression
    | or_expression
    ;

// googlesql.tm:9237
parenthesized_expression_not_a_query
    : '(' expression_maybe_parenthesized_not_a_query ')'
    ;

// googlesql.tm:9249
unparenthesized_expression_higher_prec_than_and
    : null_literal
    | boolean_literal
    | string_literal
    | bytes_literal
    | integer_literal
    | numeric_literal
    | bignumeric_literal
    | json_literal
    | floating_point_literal
    | date_or_time_literal
    | range_literal
    | parameter_expression
    | system_variable_expression
    | array_constructor
    | new_constructor
    | braced_constructor
    | braced_new_constructor
    | struct_braced_constructor
    | case_expression
    | cast_expression
    | extract_expression
    | with_expression
    | replace_fields_expression
    | function_call_expression_with_clauses
    | interval_expression
    | identifier
    | struct_constructor
    | expression_subquery_with_keyword
    | expression_higher_prec_than_and '[' expression ']'
    | expression_higher_prec_than_and '.' '(' path_expression ')'
    | expression_higher_prec_than_and '.' identifier_or_function_name_from_keyword
    | 'NOT' expression_higher_prec_than_and
    | expression_higher_prec_than_and like_operator any_some_all hint? unnest_expression
    | expression_higher_prec_than_and comparative_operator any_some_all hint? unnest_expression
    | expression_higher_prec_than_and like_operator any_some_all hint? parenthesized_anysomeall_list_rhs
    | expression_higher_prec_than_and comparative_operator any_some_all hint? parenthesized_anysomeall_list_rhs
    | expression_higher_prec_than_and like_operator expression_higher_prec_than_and
    | expression_higher_prec_than_and distinct_operator expression_higher_prec_than_and
    | expression_higher_prec_than_and in_operator hint? unnest_expression
    | expression_higher_prec_than_and in_operator hint? parenthesized_in_rhs
    | expression_higher_prec_than_and between_operator expression_higher_prec_than_and 'AND' expression_higher_prec_than_and
    | expression_higher_prec_than_and between_operator expression_higher_prec_than_and 'OR'
    | expression_higher_prec_than_and is_operator 'UNKNOWN'
    | expression_higher_prec_than_and is_operator null_literal
    | expression_higher_prec_than_and is_operator boolean_literal
    | expression_higher_prec_than_and comparative_operator expression_higher_prec_than_and
    | expression_higher_prec_than_and '|' expression_higher_prec_than_and
    | expression_higher_prec_than_and '^' expression_higher_prec_than_and
    | expression_higher_prec_than_and '&' expression_higher_prec_than_and
    | expression_higher_prec_than_and '||' expression_higher_prec_than_and
    | expression_higher_prec_than_and shift_operator expression_higher_prec_than_and
    | expression_higher_prec_than_and additive_operator expression_higher_prec_than_and
    | expression_higher_prec_than_and multiplicative_operator expression_higher_prec_than_and
    | unary_operator expression_higher_prec_than_and
    | graph_expression
    ;

// googlesql.tm:9699
path_expression
    : identifier
    | path_expression '.' identifier
    ;

// googlesql.tm:9710
dashed_identifier
    : identifier '-' identifier
    | dashed_identifier '-' identifier
    | identifier '-' INTEGER_LITERAL
    | dashed_identifier '-' INTEGER_LITERAL
    | identifier '-' FLOATING_POINT_LITERAL identifier
    | dashed_identifier '-' FLOATING_POINT_LITERAL identifier
    ;

// googlesql.tm:9806
dashed_path_expression
    : dashed_identifier
    | dashed_path_expression '.' identifier
    ;

// googlesql.tm:9823
maybe_dashed_path_expression
    : path_expression
    | dashed_path_expression
    ;

// googlesql.tm:9851
maybe_slashed_or_dashed_path_expression
    : maybe_dashed_path_expression
    | slashed_path_expression
    ;

// googlesql.tm:9878
slashed_identifier_separator
    : '-'
    | '/'
    | ':'
    ;

// googlesql.tm:9887
identifier_or_integer
    : identifier
    | INTEGER_LITERAL
    | SCRIPT_LABEL
    ;

// googlesql.tm:9905
slashed_identifier
    : '/' identifier_or_integer
    | slashed_identifier slashed_identifier_separator identifier_or_integer
    | slashed_identifier slashed_identifier_separator FLOATING_POINT_LITERAL slashed_identifier_separator identifier_or_integer
    ;

// googlesql.tm:9994
slashed_path_expression
    : slashed_identifier
    | slashed_identifier slashed_identifier_separator FLOATING_POINT_LITERAL identifier
    | slashed_identifier slashed_identifier_separator FLOATING_POINT_LITERAL '.' identifier
    | slashed_path_expression '.' identifier
    ;

// googlesql.tm:10109
array_constructor_prefix_no_expressions
    : 'ARRAY' '['
    | '['
    | array_type '['
    ;

// googlesql.tm:10124
array_constructor_prefix
    : array_constructor_prefix_no_expressions expression
    | array_constructor_prefix ',' expression
    ;

// googlesql.tm:10135
array_constructor
    : array_constructor_prefix_no_expressions ']'
    | array_constructor_prefix ']'
    ;

// googlesql.tm:10146
range_literal
    : range_type string_literal
    ;

// googlesql.tm:10153
date_or_time_literal_kind
    : 'DATE'
    | 'DATETIME'
    | 'TIME'
    | 'TIMESTAMP'
    ;

// googlesql.tm:10172
date_or_time_literal
    : date_or_time_literal_kind string_literal
    ;

// googlesql.tm:10181
interval_expression
    : 'INTERVAL' expression identifier
    | 'INTERVAL' expression identifier 'TO' identifier
    ;

// googlesql.tm:10192
parameter_expression
    : named_parameter_expression
    | '?'
    ;

// googlesql.tm:10204
named_parameter_expression
    : '@' identifier
    ;

// googlesql.tm:10214
type_name
    : path_expression
    | 'INTERVAL'
    ;

// googlesql.tm:10228
template_type_open
    : '<'
    ;

// googlesql.tm:10238
template_type_close
    : '>'
    ;

// googlesql.tm:10248
array_type
    : 'ARRAY' template_type_open type template_type_close
    ;

// googlesql.tm:10255
struct_field
    : identifier type
    | type
    ;

// googlesql.tm:10266
struct_type_prefix
    : 'STRUCT' template_type_open struct_field
    | struct_type_prefix ',' struct_field
    ;

// googlesql.tm:10277
struct_type
    : 'STRUCT' template_type_open template_type_close
    | struct_type_prefix template_type_close
    ;

// googlesql.tm:10288
range_type
    : 'RANGE' template_type_open type template_type_close
    ;

// googlesql.tm:10295
function_type_prefix
    : 'FUNCTION' template_type_open '(' type
    | function_type_prefix ',' type
    ;

// googlesql.tm:10306
function_type
    : 'FUNCTION' template_type_open '(' ')' '->' type template_type_close
    | 'FUNCTION' template_type_open type '->' type template_type_close
    | function_type_prefix ')' '->' type template_type_close
    ;

// googlesql.tm:10324
map_type
    : 'MAP' template_type_open type ',' type template_type_close
    ;

// googlesql.tm:10332
raw_type
    : array_type
    | struct_type
    | type_name
    | range_type
    | function_type
    | map_type
    ;

// googlesql.tm:10341
type_parameter
    : integer_literal
    | boolean_literal
    | string_literal
    | bytes_literal
    | floating_point_literal
    | 'MAX'
    ;

// googlesql.tm:10362
type_parameters_prefix
    : '(' type_parameter
    | type_parameters_prefix ',' type_parameter
    ;

// googlesql.tm:10373
type_parameters
    : type_parameters_prefix ')'
    | type_parameters_prefix ',' ')'
    ;

// googlesql.tm:10386
opt_type_parameters
    : type_parameters?
    ;

// googlesql.tm:10393
type
    : raw_type type_parameters? collate_clause?
    ;

// googlesql.tm:10400
templated_parameter_kind
    : 'PROTO'
    | 'ENUM'
    | 'STRUCT'
    | 'ARRAY'
    | identifier
    ;

// googlesql.tm:10431
templated_parameter_type
    : 'ANY' templated_parameter_kind
    ;

// googlesql.tm:10440
type_or_tvf_schema
    : type
    | templated_parameter_type
    | tvf_schema
    ;

// googlesql.tm:10446
new_constructor_prefix_no_arg
    : 'NEW' type_name '('
    ;

// googlesql.tm:10453
new_constructor_arg
    : expression
    | expression 'AS' identifier
    | expression 'AS' '(' path_expression ')'
    ;

// googlesql.tm:10470
new_constructor_prefix
    : new_constructor_prefix_no_arg new_constructor_arg
    | new_constructor_prefix ',' new_constructor_arg
    ;

// googlesql.tm:10481
new_constructor
    : new_constructor_prefix ')'
    | new_constructor_prefix_no_arg ')'
    ;

// googlesql.tm:10492
braced_constructor_extension_expression
    : '(' path_expression ')' ('[' expression ']')? ('.' generalized_path_expression)?
    ;

// googlesql.tm:10509
braced_constructor_field
    : (braced_constructor_extension_expression | generalized_path_expression) ('*' | '?')? (':' expression | braced_constructor)
    ;

// googlesql.tm:10533
braced_constructor_field_following_omitted_comma
    : generalized_path_expression ('*' | '?')? (':' expression | braced_constructor)
    ;

// googlesql.tm:10555
braced_constructor_field_list_continuation
    : ',' braced_constructor_field
    | braced_constructor_field_following_omitted_comma
    ;

// googlesql.tm:10571
braced_constructor
    : '{' (braced_constructor_field (braced_constructor_field_list_continuation)* ','?)? '}'
    ;

// googlesql.tm:10582
braced_new_constructor
    : 'NEW' type_name braced_constructor
    ;

// googlesql.tm:10589
struct_braced_constructor
    : struct_type braced_constructor
    | 'STRUCT' braced_constructor
    ;

// googlesql.tm:10600
case_no_value_expression_prefix
    : 'CASE' 'WHEN' expression 'THEN' expression
    | case_no_value_expression_prefix 'WHEN' expression 'THEN' expression
    ;

// googlesql.tm:10611
case_value_expression_prefix
    : 'CASE' expression 'WHEN' expression 'THEN' expression
    | case_value_expression_prefix 'WHEN' expression 'THEN' expression
    ;

// googlesql.tm:10622
case_expression_prefix
    : case_no_value_expression_prefix
    | case_value_expression_prefix
    ;

// googlesql.tm:10627
case_expression
    : case_expression_prefix 'END'
    | case_expression_prefix 'ELSE' expression 'END'
    ;

// googlesql.tm:10638
at_time_zone
    : 'AT' 'TIME' 'ZONE' expression
    ;

// googlesql.tm:10645
opt_format
    : format?
    ;

// googlesql.tm:10652
format
    : 'FORMAT' expression at_time_zone?
    ;

// googlesql.tm:10659
cast_expression
    : ('CAST' | 'SAFE_CAST') '(' expression 'AS' type format? ')'
    | 'CAST' '(' 'SELECT'
    | 'SAFE_CAST' '(' 'SELECT'
    ;

// googlesql.tm:10684
extract_expression_base
    : 'EXTRACT' '(' expression 'FROM' expression
    ;

// googlesql.tm:10691
extract_expression
    : extract_expression_base ')'
    | extract_expression_base 'AT' 'TIME' 'ZONE' expression ')'
    ;

// googlesql.tm:10702
replace_fields_arg
    : expression 'AS' generalized_path_expression
    | expression 'AS' generalized_extension_path
    ;

// googlesql.tm:10713
replace_fields_prefix
    : 'REPLACE_FIELDS' '(' expression ',' replace_fields_arg
    | replace_fields_prefix ',' replace_fields_arg
    ;

// googlesql.tm:10724
replace_fields_expression
    : replace_fields_prefix ')'
    ;

// googlesql.tm:10731
function_name_from_keyword
    : 'IF'
    | 'GROUPING'
    | 'LEFT'
    | 'RIGHT'
    | 'COLLATE'
    | 'RANGE'
    ;

// googlesql.tm:10762
identifier_or_function_name_from_keyword
    : identifier
    | function_name_from_keyword
    ;

// googlesql.tm:10787
function_call_expression_base
    : expression_higher_prec_than_and '(' opt_distinct
    | function_name_from_keyword '(' opt_distinct
    ;

// googlesql.tm:10873
function_call_argument
    : expression as_alias_with_required_as?
    | named_argument
    | lambda_argument
    | sequence_arg
    | 'SELECT'
    ;

// googlesql.tm:10900
sequence_arg
    : 'SEQUENCE' path_expression
    ;

// googlesql.tm:10907
named_argument
    : identifier '=>' expression
    | identifier '=>' lambda_argument
    | identifier '=>' input_table_argument
    ;

// googlesql.tm:10922
lambda_argument
    : lambda_argument_list '->' expression
    ;

// googlesql.tm:10941
lambda_argument_list
    : expression
    | '(' ')'
    ;

// googlesql.tm:10959
function_call_expression_with_args_prefix
    : function_call_expression_base function_call_argument
    | function_call_expression_base '*'
    | function_call_expression_with_args_prefix ',' function_call_argument
    ;

// googlesql.tm:10983
function_call_expression
    : function_call_expression_base null_handling_modifier opt_where_clause having_modifier? (group_by_clause_prefix having_clause?)? with_report_modifier? order_by_clause? opt_limit_offset_clause ')'
    | function_call_expression_with_args_prefix null_handling_modifier opt_where_clause having_modifier? (group_by_clause_prefix having_clause?)? clamped_between_modifier? with_report_modifier? order_by_clause? opt_limit_offset_clause ')'
    ;

// googlesql.tm:11031
opt_identifier
    : identifier?
    ;

// googlesql.tm:11038
partition_by_clause_prefix
    : 'PARTITION' hint? 'BY' expression
    | partition_by_clause_prefix ',' expression
    ;

// googlesql.tm:11049
opt_partition_by_clause
    : partition_by_clause_prefix?
    ;

// googlesql.tm:11056
partition_by_clause_prefix_no_hint
    : 'PARTITION' 'BY' expression
    | partition_by_clause_prefix_no_hint ',' expression
    ;

// googlesql.tm:11067
opt_partition_by_clause_no_hint
    : partition_by_clause_prefix_no_hint?
    ;

// googlesql.tm:11074
cluster_by_clause_prefix_no_hint
    : 'CLUSTER' 'BY' expression
    | cluster_by_clause_prefix_no_hint ',' expression
    ;

// googlesql.tm:11085
opt_cluster_by_clause_no_hint
    : cluster_by_clause_no_hint?
    ;

// googlesql.tm:11092
cluster_by_clause_no_hint
    : cluster_by_clause_prefix_no_hint
    ;

// googlesql.tm:11099
ttl_clause
    : 'ROW' 'DELETION' 'POLICY' '(' expression ')'
    ;

// googlesql.tm:11109
opt_ttl_clause
    : ttl_clause?
    ;

// googlesql.tm:11117
preceding_or_following
    : 'PRECEDING'
    | 'FOLLOWING'
    ;

// googlesql.tm:11128
window_frame_bound
    : 'UNBOUNDED' preceding_or_following
    | 'CURRENT' 'ROW'
    | expression preceding_or_following
    ;

// googlesql.tm:11156
frame_unit
    : 'ROWS'
    | 'RANGE'
    ;

// googlesql.tm:11167
window_frame_clause
    : frame_unit 'BETWEEN' window_frame_bound 'AND' window_frame_bound
    | frame_unit window_frame_bound
    ;

// googlesql.tm:11180
window_specification
    : identifier
    | '(' identifier? opt_partition_by_clause order_by_clause? window_frame_clause? ')'
    ;

// googlesql.tm:11195
function_call_expression_with_clauses
    : function_call_expression hint? over_clause?
    | function_call_expression braced_constructor
    ;

// googlesql.tm:11211
over_clause
    : 'OVER' window_specification
    ;

// googlesql.tm:11218
struct_constructor_prefix_with_keyword_no_arg
    : struct_type '('
    | 'STRUCT' '('
    ;

// googlesql.tm:11229
struct_constructor_prefix_with_keyword
    : struct_constructor_prefix_with_keyword_no_arg struct_constructor_arg
    | struct_constructor_prefix_with_keyword ',' struct_constructor_arg
    ;

// googlesql.tm:11240
struct_constructor_arg
    : expression as_alias_with_required_as?
    ;

// googlesql.tm:11247
struct_constructor_prefix_without_keyword
    : '(' expression ',' expression
    | struct_constructor_prefix_without_keyword ',' expression
    ;

// googlesql.tm:11260
struct_constructor
    : struct_constructor_prefix_with_keyword ')'
    | struct_constructor_prefix_with_keyword_no_arg ')'
    | struct_constructor_prefix_without_keyword ')'
    ;

// googlesql.tm:11275
expression_subquery_with_keyword
    : 'ARRAY' parenthesized_query
    | 'EXISTS' hint? parenthesized_query
    | 'ARRAY' braced_graph_subquery
    | exists_graph_subquery
    | 'VALUE' hint? braced_graph_subquery
    ;

// googlesql.tm:11306
null_literal
    : 'NULL'
    ;

// googlesql.tm:11317
boolean_literal
    : 'TRUE'
    | 'FALSE'
    ;

// googlesql.tm:11338
string_literal_component
    : STRING_LITERAL
    ;

// googlesql.tm:11369
string_literal
    : string_literal_component
    | string_literal string_literal_component
    | string_literal bytes_literal_component
    ;

// googlesql.tm:11396
bytes_literal_component
    : BYTES_LITERAL
    ;

// googlesql.tm:11431
bytes_literal
    : bytes_literal_component
    | bytes_literal bytes_literal_component
    | bytes_literal string_literal_component
    ;

// googlesql.tm:11457
integer_literal
    : INTEGER_LITERAL
    ;

// googlesql.tm:11466
numeric_literal_prefix
    : 'NUMERIC'
    | 'DECIMAL'
    ;

// googlesql.tm:11471
numeric_literal
    : numeric_literal_prefix string_literal
    ;

// googlesql.tm:11478
bignumeric_literal_prefix
    : 'BIGNUMERIC'
    | 'BIGDECIMAL'
    ;

// googlesql.tm:11483
bignumeric_literal
    : bignumeric_literal_prefix string_literal
    ;

// googlesql.tm:11490
json_literal
    : 'JSON' string_literal
    ;

// googlesql.tm:11497
floating_point_literal
    : FLOATING_POINT_LITERAL
    ;

// googlesql.tm:11506
token_identifier
    : IDENTIFIER
    ;

// googlesql.tm:11539
identifier
    : token_identifier
    | keyword_as_identifier
    ;

// googlesql.tm:11547
label
    : SCRIPT_LABEL
    ;

// googlesql.tm:11580
system_variable_expression
    : '@@' path_expression
    ;

// googlesql.tm:11593
common_keyword_as_identifier
    : 'ABORT'
    | 'ACCESS'
    | 'ACTION'
    | 'AGGREGATE'
    | 'ADD'
    | 'AFTER'
    | 'ALTER'
    | 'ALWAYS'
    | 'ANALYZE'
    | 'APPROX'
    | 'ARE'
    | 'ASCENDING'
    | 'ASSERT'
    | 'BATCH'
    | 'BEGIN'
    | 'BIGDECIMAL'
    | 'BIGNUMERIC'
    | 'BREAK'
    | 'CALL'
    | 'CASCADE'
    | 'CHEAPEST'
    | 'CHECK'
    | 'CLAMPED'
    | 'CLONE'
    | 'COPY'
    | 'COST'
    | 'CLUSTER'
    | 'COLUMN'
    | 'COLUMNS'
    | 'COMMIT'
    | 'CONFLICT'
    | 'CONNECTION'
    | 'CONSTANT'
    | 'CONSTRAINT'
    | 'CONTINUE'
    | 'CORRESPONDING'
    | 'CYCLE'
    | 'DATA'
    | 'DATABASE'
    | 'DATE'
    | 'DATETIME'
    | 'DECIMAL'
    | 'DECLARE'
    | 'DEFINER'
    | 'DELETE'
    | 'DELETION'
    | 'DEPTH'
    | 'DESCENDING'
    | 'DESCRIBE'
    | 'DESTINATION'
    | 'DETERMINISTIC'
    | 'DO'
    | 'DROP'
    | 'DYNAMIC'
    | 'EDGE'
    | 'ELSEIF'
    | 'ENFORCED'
    | 'ERROR'
    | 'EXCEPTION'
    | 'EXECUTE'
    | 'EXPLAIN'
    | 'EXPORT'
    | 'EXTEND'
    | 'EXTERNAL'
    | 'FILES'
    | 'FILTER'
    | 'FILL'
    | 'FIRST'
    | 'FOREIGN'
    | 'FORK'
    | 'FORMAT'
    | 'FUNCTION'
    | 'GENERATED'
    | 'GRAPH_TABLE'
    | 'GRAPH'
    | 'GRANT'
    | 'GROUP_ROWS'
    | 'HIDDEN'
    | 'IDENTITY'
    | 'IMMEDIATE'
    | 'IMMUTABLE'
    | 'IMPORT'
    | 'INCLUDE'
    | 'INCREMENT'
    | 'INDEX'
    | 'INOUT'
    | 'INPUT'
    | 'INSERT'
    | 'INVOKER'
    | 'ISOLATION'
    | 'ITERATE'
    | 'JSON'
    | 'KEY'
    | 'LABEL'
    | 'LABELED'
    | 'LANGUAGE'
    | 'LAST'
    | 'LEAVE'
    | 'LET'
    | 'LEVEL'
    | 'LOAD'
    | 'LOCALITY'
    | 'LOG'
    | 'LOOP'
    | 'MACRO'
    | 'MAP'
    | 'MATCH'
    | 'MATCH_RECOGNIZE'
    | 'MATCHED'
    | 'MATERIALIZED'
    | 'MAX'
    | 'MAXVALUE'
    | 'MEASURES'
    | 'MESSAGE'
    | 'METADATA'
    | 'MIN'
    | 'MINVALUE'
    | 'MODEL'
    | 'MODULE'
    | 'NAME'
    | 'NEXT'
    | 'NODE'
    | 'NOTHING'
    | 'NUMERIC'
    | 'OFFSET'
    | 'ONLY'
    | 'OPTIONAL'
    | 'OPTIONS'
    | 'OUT'
    | 'OUTPUT'
    | 'OVERWRITE'
    | 'PARTITIONS'
    | 'PAST'
    | 'PATTERN'
    | 'PER'
    | 'PERCENT'
    | 'PIVOT'
    | 'POLICIES'
    | 'POLICY'
    | 'PRIMARY'
    | 'PRIVATE'
    | 'PRIVILEGE'
    | 'PRIVILEGES'
    | 'PROCEDURE'
    | 'PROJECT'
    | 'PROPERTIES'
    | 'PROPERTY'
    | 'PUBLIC'
    | 'QUALIFY'
    | 'RAISE'
    | 'READ'
    | 'REBUILD'
    | 'REFERENCES'
    | 'REMOTE'
    | 'REMOVE'
    | 'RENAME'
    | 'REPEAT'
    | 'REPEATABLE'
    | 'REPLACE'
    | 'REPLACE_FIELDS'
    | 'REPLICA'
    | 'REPORT'
    | 'RESTRICT'
    | 'RESTRICTION'
    | 'RETURNS'
    | 'RETURN'
    | 'REVOKE'
    | 'ROLLBACK'
    | 'ROW'
    | 'RUN'
    | 'SAFE_CAST'
    | 'SCHEMA'
    | 'SEARCH'
    | 'SECURITY'
    | 'SEQUENCE'
    | 'SETS'
    | 'SHORTEST'
    | 'SHOW'
    | 'SKIP'
    | 'SNAPSHOT'
    | 'SOURCE'
    | 'SQL'
    | 'STABLE'
    | 'START'
    | 'STATIC_DESCRIBE'
    | 'STORED'
    | 'STORING'
    | 'STRICT'
    | 'SYSTEM'
    | 'SYSTEM_TIME'
    | 'TABLE'
    | 'TABLES'
    | 'TARGET'
    | 'TEMP'
    | 'TEMPORARY'
    | 'TEE'
    | 'TIME'
    | 'TIMESTAMP'
    | 'TRANSACTION'
    | 'TRANSFORM'
    | 'TRUNCATE'
    | 'TYPE'
    | 'UNDROP'
    | 'UNIQUE'
    | 'UNKNOWN'
    | 'UNPIVOT'
    | 'UNTIL'
    | 'UPDATE'
    | 'VALUE'
    | 'VALUES'
    | 'VECTOR'
    | 'VIEW'
    | 'VIEWS'
    | 'VOLATILE'
    | 'WEIGHT'
    | 'WHILE'
    | 'WRITE'
    | 'YIELD'
    | 'ZONE'
    | 'DESCRIPTOR'
    | 'INTERLEAVE'
    | 'NULL_FILTERED'
    | 'PARENT'
    ;

// googlesql.tm:11834
keyword_as_identifier
    : common_keyword_as_identifier
    | 'SIMPLE'
    | 'ACYCLIC'
    | 'PATH'
    | 'PATHS'
    | 'TRAIL'
    | 'WALK'
    ;

// googlesql.tm:11845
or_replace
    : 'OR' 'REPLACE'
    ;

// googlesql.tm:11849
opt_or_replace
    : or_replace?
    ;

// googlesql.tm:11856
create_scope
    : 'TEMP'
    | 'TEMPORARY'
    | 'PUBLIC'
    | 'PRIVATE'
    ;

// googlesql.tm:11875
opt_create_scope
    : create_scope?
    ;

// googlesql.tm:11882
unique
    : 'UNIQUE'
    ;

// googlesql.tm:11886
opt_unique
    : unique?
    ;

// googlesql.tm:11893
describe_keyword
    : 'DESCRIBE'
    | 'DESC'
    ;

// googlesql.tm:11898
options_entry
    : identifier_in_hints options_assignment_operator expression_or_proto
    ;

// googlesql.tm:11907
options_assignment_operator
    : '='
    | '+='
    | '-='
    ;

// googlesql.tm:11922
expression_or_proto
    : 'PROTO'
    | expression
    ;

// googlesql.tm:11932
options_list
    : '(' ((options_entry) (',' (options_entry))*)? ')'
    ;

// googlesql.tm:11939
options_
    : 'OPTIONS' options_list
    ;

// googlesql.tm:11946
opt_options_list
    : options_?
    ;

// googlesql.tm:11953
define_table_statement
    : 'DEFINE' 'TABLE' path_expression options_list
    ;

// googlesql.tm:11960
dml_statement
    : insert_statement
    | delete_statement
    | update_statement
    ;

// googlesql.tm:11969
where_expression
    : 'WHERE' expression
    ;

// googlesql.tm:11976
assert_rows_modified
    : 'ASSERT_ROWS_MODIFIED' possibly_cast_int_literal_or_parameter
    ;

// googlesql.tm:11983
returning_clause
    : 'THEN' 'RETURN' select_list
    | 'THEN' 'RETURN' 'WITH' 'ACTION' select_list
    | 'THEN' 'RETURN' 'WITH' 'ACTION' 'AS' identifier select_list
    ;

// googlesql.tm:12002
conflict_target
    : column_list
    | 'ON' 'UNIQUE' 'CONSTRAINT' identifier
    ;

// googlesql.tm:12013
on_conflict_clause
    : 'ON' 'CONFLICT' conflict_target? 'DO' ('NOTHING' | 'UPDATE' 'SET' ((update_item) (',' (update_item))*) where_expression?)
    ;

// googlesql.tm:12026
insert_mode
    : 'OR'? 'IGNORE'
    | (KW_REPLACE_AFTER_INSERT | 'OR' 'REPLACE')
    | (KW_UPDATE_AFTER_INSERT | 'OR' 'UPDATE')
    ;

// googlesql.tm:12041
insert_source_and_opt_on_conflict
    : query
    | '(' query ')' on_conflict_clause
    | 'VALUES' ((insert_values_row) (',' (insert_values_row))*) on_conflict_clause?
    | table_clause_unreserved on_conflict_clause?
    ;

// googlesql.tm:12073
identifierPreceedsParenthsizedQuery
    : /* empty */
    ;

// googlesql.tm:12080
insert_statement
    : 'INSERT' insert_mode? 'INTO'? maybe_dashed_generalized_path_expression identifierPreceedsParenthsizedQuery hint? column_list? insert_source_and_opt_on_conflict assert_rows_modified? returning_clause?
    ;

// googlesql.tm:12107
insert_statement_in_pipe_prefix
    : 'INSERT' insert_mode? 'INTO'? maybe_dashed_generalized_path_expression hint? column_list? on_conflict_clause?
    ;

// googlesql.tm:12120
insert_statement_in_pipe
    : insert_statement_in_pipe_prefix assert_rows_modified? returning_clause?
    ;

// googlesql.tm:12129
copy_data_source
    : maybe_dashed_path_expression at_system_time? opt_where_clause
    ;

// googlesql.tm:12136
clone_data_source
    : maybe_dashed_path_expression at_system_time? opt_where_clause
    ;

// googlesql.tm:12143
clone_data_statement
    : 'CLONE' 'DATA' 'INTO' maybe_dashed_path_expression 'FROM' ((clone_data_source) ('UNION' 'ALL' (clone_data_source))*)
    ;

// googlesql.tm:12152
expression_or_default
    : expression
    | 'DEFAULT'
    ;

// googlesql.tm:12160
insert_values_row
    : '(' ((expression_or_default) (',' (expression_or_default))*) ')'
    ;

// googlesql.tm:12167
delete_statement
    : 'DELETE' 'FROM'? maybe_dashed_generalized_path_expression hint? as_alias? with_offset_and_alias? where_expression? assert_rows_modified? returning_clause?
    ;

// googlesql.tm:12177
with_offset
    : 'WITH' 'OFFSET'
    ;

// googlesql.tm:12184
with_offset_and_alias
    : with_offset as_alias?
    ;

// googlesql.tm:12191
update_statement
    : 'UPDATE' maybe_dashed_generalized_path_expression hint? as_alias? with_offset_and_alias? 'SET' ((update_item) (',' (update_item))*) opt_from_clause where_expression? assert_rows_modified? returning_clause?
    ;

// googlesql.tm:12202
truncate_statement
    : 'TRUNCATE' 'TABLE' maybe_dashed_path_expression where_expression?
    ;

// googlesql.tm:12209
nested_dml_statement
    : '(' dml_statement ')'
    ;

// googlesql.tm:12224
generalized_path_expression
    : identifier
    | generalized_path_expression '.' generalized_extension_path
    | generalized_path_expression '.' identifier
    | generalized_path_expression '[' expression ']'
    ;

// googlesql.tm:12254
maybe_dashed_generalized_path_expression
    : generalized_path_expression
    | dashed_path_expression
    ;

// googlesql.tm:12292
generalized_extension_path
    : '(' path_expression ')'
    | generalized_extension_path '.' '(' path_expression ')'
    | generalized_extension_path '.' identifier
    ;

// googlesql.tm:12308
update_set_value
    : generalized_path_expression '=' expression_or_default
    ;

// googlesql.tm:12315
update_item
    : update_set_value
    | nested_dml_statement
    ;

// googlesql.tm:12326
merge_insert_value_list_or_source_row
    : 'VALUES' insert_values_row
    | 'ROW'
    ;

// googlesql.tm:12337
merge_action
    : 'INSERT' opt_column_list merge_insert_value_list_or_source_row
    | 'UPDATE' 'SET' ((update_item) (',' (update_item))*)
    | 'DELETE'
    ;

// googlesql.tm:12356
merge_when_clause
    : 'WHEN' ('MATCHED' | 'NOT' 'MATCHED' ('BY' 'TARGET')? | 'NOT' 'MATCHED' 'BY' 'SOURCE') ('AND' expression)? 'THEN' merge_action
    ;

// googlesql.tm:12374
merge_source
    : table_path_expression
    | table_subquery
    ;

// googlesql.tm:12382
merge_statement_prefix
    : 'MERGE' 'INTO'? maybe_dashed_path_expression as_alias? 'USING' merge_source 'ON' expression
    ;

// googlesql.tm:12391
merge_statement
    : merge_statement_prefix merge_when_clause+
    ;

// googlesql.tm:12398
call_statement_with_args_prefix
    : 'CALL' path_expression '(' tvf_argument
    | call_statement_with_args_prefix ',' tvf_argument
    ;

// googlesql.tm:12409
call_statement
    : call_statement_with_args_prefix ')'
    | 'CALL' path_expression '(' ')'
    ;

// googlesql.tm:12420
opt_function_parameters
    : function_parameters?
    ;

// googlesql.tm:12428
opt_if_exists
    : 'IF' 'EXISTS'
    | /* empty */
    ;

// googlesql.tm:12440
access
    : 'ACCESS'
    ;

// googlesql.tm:12448
opt_access
    : access?
    ;

// googlesql.tm:12456
drop_all_row_access_policies_statement
    : 'DROP' 'ALL' 'ROW' access? 'POLICIES' 'ON' path_expression
    ;

// googlesql.tm:12465
on_path_expression
    : 'ON' path_expression
    ;

// googlesql.tm:12472
opt_on_path_expression
    : on_path_expression?
    ;

// googlesql.tm:12479
drop_mode
    : 'RESTRICT'
    | 'CASCADE'
    ;

// googlesql.tm:12490
drop_statement
    : 'DROP' 'PRIVILEGE' 'RESTRICTION' opt_if_exists 'ON' privilege_list 'ON' identifier path_expression
    | 'DROP' 'ROW' 'ACCESS' 'POLICY' opt_if_exists identifier on_path_expression
    | 'DROP' index_type 'INDEX' opt_if_exists path_expression on_path_expression?
    | 'DROP' table_or_table_function opt_if_exists maybe_dashed_path_expression function_parameters?
    | 'DROP' 'SNAPSHOT' 'TABLE' opt_if_exists maybe_dashed_path_expression
    | 'DROP' generic_entity_type opt_if_exists path_expression
    | 'DROP' schema_object_kind opt_if_exists path_expression function_parameters? drop_mode?
    ;

// googlesql.tm:12619
alter_index_type
    : 'SEARCH'
    | 'VECTOR'
    ;

// googlesql.tm:12630
index_type
    : 'SEARCH'
    | 'VECTOR'
    ;

// googlesql.tm:12641
opt_index_type
    : index_type?
    ;

// googlesql.tm:12648
unterminated_non_empty_statement_list
    : unterminated_statement
    | unterminated_non_empty_statement_list ';' unterminated_statement
    ;

// googlesql.tm:12684
execute_into_clause
    : 'INTO' identifier_list
    ;

// googlesql.tm:12691
execute_using_argument
    : expression 'AS' identifier
    | expression
    ;

// googlesql.tm:12705
execute_using_argument_list
    : execute_using_argument
    | execute_using_argument_list ',' execute_using_argument
    ;

// googlesql.tm:12716
execute_using_clause
    : 'USING' execute_using_argument_list
    ;

// googlesql.tm:12723
execute_immediate
    : 'EXECUTE' 'IMMEDIATE' expression execute_into_clause? execute_using_clause?
    ;

// googlesql.tm:12749
statement_list
    : unterminated_non_empty_statement_list ';'
    | /* empty */
    ;

// googlesql.tm:12761
else_clause
    : 'ELSE' statement_list
    ;

// googlesql.tm:12768
elseif_clauses
    : 'ELSEIF' expression 'THEN' statement_list
    | elseif_clauses 'ELSEIF' expression 'THEN' statement_list
    ;

// googlesql.tm:12781
if_statement_unclosed
    : 'IF' expression 'THEN' statement_list elseif_clauses? else_clause?
    ;

// googlesql.tm:12788
if_statement
    : if_statement_unclosed 'END' 'IF'
    ;

// googlesql.tm:12796
when_then_clauses
    : 'WHEN' expression 'THEN' statement_list
    | when_then_clauses 'WHEN' expression 'THEN' statement_list
    ;

// googlesql.tm:12809
case_statement
    : 'CASE' expression? when_then_clauses else_clause? 'END' 'CASE'
    ;

// googlesql.tm:12819
begin_end_block
    : 'BEGIN' statement_list exception_handler? 'END'
    ;

// googlesql.tm:12827
exception_handler
    : 'EXCEPTION' 'WHEN' 'ERROR' 'THEN' statement_list
    ;

// googlesql.tm:12836
default_expression
    : 'DEFAULT' expression
    ;

// googlesql.tm:12843
opt_default_expression
    : default_expression?
    ;

// googlesql.tm:12850
identifier_list
    : ((identifier) (',' (identifier))*)
    ;

// googlesql.tm:12857
variable_declaration
    : 'DECLARE' identifier_list type default_expression?
    | 'DECLARE' identifier_list 'DEFAULT' expression
    ;

// googlesql.tm:12868
loop_statement
    : 'LOOP' statement_list 'END' 'LOOP'
    ;

// googlesql.tm:12878
while_statement
    : 'WHILE' expression 'DO' statement_list 'END' 'WHILE'
    ;

// googlesql.tm:12889
until_clause
    : 'UNTIL' expression
    ;

// googlesql.tm:12896
repeat_statement
    : 'REPEAT' statement_list until_clause 'END' 'REPEAT'
    ;

// googlesql.tm:12909
for_in_statement
    : 'FOR' identifier 'IN' parenthesized_query 'DO' statement_list 'END' 'FOR'
    ;

// googlesql.tm:12924
break_statement
    : 'BREAK' opt_identifier
    | 'LEAVE' opt_identifier
    ;

// googlesql.tm:12953
continue_statement
    : 'CONTINUE' opt_identifier
    | 'ITERATE' opt_identifier
    ;

// googlesql.tm:12984
return_statement
    : 'RETURN'
    ;

// googlesql.tm:12991
raise_statement
    : 'RAISE'
    | 'RAISE' 'USING' 'MESSAGE' '=' expression
    ;

// googlesql.tm:13534
spanner_primary_key
    : 'PRIMARY' 'KEY' primary_key_element_list
    ;

// googlesql.tm:13546
spanner_index_interleave_clause
    : ',' 'INTERLEAVE' 'IN' maybe_dashed_path_expression
    ;

// googlesql.tm:13560
opt_spanner_interleave_in_parent_clause
    : ',' 'INTERLEAVE' 'IN' 'PARENT' maybe_dashed_path_expression opt_foreign_key_on_delete
    | /* empty */
    ;

// googlesql.tm:13581
opt_spanner_table_options
    : spanner_primary_key opt_spanner_interleave_in_parent_clause
    | /* empty */
    ;

// googlesql.tm:13598
spanner_null_filtered
    : 'NULL_FILTERED'
    ;

// googlesql.tm:13610
opt_spanner_null_filtered
    : spanner_null_filtered?
    ;

// googlesql.tm:13619
spanner_generated_or_default
    : 'AS' '(' expression ')' 'STORED'
    | default_column_info
    ;

// googlesql.tm:13629
spanner_not_null_attribute
    : not_null_column_attribute
    ;

// googlesql.tm:13638
spanner_alter_column_action
    : 'ALTER' 'COLUMN' opt_if_exists identifier column_schema_inner spanner_not_null_attribute? spanner_generated_or_default? opt_options_list
    ;

// googlesql.tm:13665
spanner_set_on_delete_action
    : 'SET' 'ON' 'DELETE' foreign_key_action
    ;

// googlesql.tm:13678
create_sequence_statement
    : 'CREATE' opt_or_replace 'SEQUENCE' opt_if_not_exists path_expression options_?
    ;
