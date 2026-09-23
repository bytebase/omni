parser grammar GoogleSQLParser;

options {
	tokenVocab = GoogleSQLLexer;
}

root: stmts EOF;

stmts:
	unterminated_sql_statement (
		SEMI_SYMBOL unterminated_sql_statement
	)* SEMI_SYMBOL?;

unterminated_sql_statement:
	statement_level_hint? sql_statement_body
	| DEFINE_SYMBOL MACRO_SYMBOL {
		p.NotifyErrorListeners("Syntax error: DEFINE MACRO statements cannot be composed from other expansions", nil, nil)
	 }
	| statement_level_hint DEFINE_SYMBOL MACRO_SYMBOL {
		p.NotifyErrorListeners("Hints are not allowed on DEFINE MACRO statements", nil, nil)
	 };

sql_statement_body:
	query_statement
	| alter_statement
	| analyze_statement
	| assert_statement
	| aux_load_data_statement
	| clone_data_statement
	| dml_statement
	| merge_statement
	| truncate_statement
	| begin_statement
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
	| create_external_table_statement
	| create_external_table_function_statement
	| create_model_statement
	| create_property_graph_statement
	| create_schema_statement
	| create_external_schema_statement
	| create_snapshot_statement
	| create_table_function_statement
	| create_table_statement
	| create_view_statement
	| create_entity_statement
	// /* TODO(zp): define macro statement */ | define_macro_statement
	| define_table_statement
	| describe_statement
	| execute_immediate
	| explain_statement
	| export_data_statement
	| export_model_statement
	| export_metadata_statement
	| gql_statement
	| grant_statement
	| rename_statement
	| revoke_statement
	| rollback_statement
	| show_statement
	| drop_all_row_access_policies_statement
	| drop_statement
	| call_statement
	| import_statement
	| module_statement
	| undrop_statement;

gql_statement:
	GRAPH_SYMBOL path_expression graph_operation_block;

graph_operation_block:
	graph_composite_query_block (
		NEXT_SYMBOL graph_composite_query_block
	)*;

graph_composite_query_block:
	graph_linear_query_operation
	| graph_composite_query_prefix;

graph_composite_query_prefix:
	graph_linear_query_operation graph_set_operation_metadata graph_linear_query_operation (
		graph_set_operation_metadata graph_linear_query_operation
	)*;

graph_set_operation_metadata:
	query_set_operation_type all_or_distinct;

graph_linear_query_operation:
	graph_linear_operator_list? graph_return_operator;

graph_linear_operator_list: graph_linear_operator+;

graph_linear_operator:
	graph_match_operator
	| graph_optional_match_operator
	| graph_let_operator
	| graph_filter_operator
	| graph_order_by_operator
	| graph_page_operator
	| graph_with_operator
	| graph_for_operator
	| graph_sample_clause;

graph_sample_clause:
	TABLESAMPLE_SYMBOL identifier LR_BRACKET_SYMBOL sample_size RR_BRACKET_SYMBOL
		opt_graph_sample_clause_suffix?;

opt_graph_sample_clause_suffix:
	repeatable_clause
	| WITH_SYMBOL WEIGHT_SYMBOL repeatable_clause?
	| WITH_SYMBOL WEIGHT_SYMBOL AS_SYMBOL identifier repeatable_clause?;

graph_for_operator:
	FOR_SYMBOL identifier IN_SYMBOL expression opt_with_offset_and_alias_with_required_as?;

opt_with_offset_and_alias_with_required_as:
	WITH_SYMBOL OFFSET_SYMBOL opt_as_alias_with_required_as?;

graph_with_operator:
	WITH_SYMBOL all_or_distinct? hint? graph_return_item_list group_by_clause?;

graph_page_operator: graph_page_clause;

graph_order_by_operator: graph_order_by_clause;

graph_filter_operator:
	FILTER_SYMBOL where_clause
	| FILTER_SYMBOL expression;

graph_let_operator:
	LET_SYMBOL graph_let_variable_definition_list;

graph_let_variable_definition_list:
	graph_let_variable_definition (
		COMMA_SYMBOL graph_let_variable_definition
	)*;

graph_let_variable_definition:
	identifier EQUAL_OPERATOR expression;

graph_optional_match_operator:
	OPTIONAL_SYMBOL MATCH_SYMBOL hint? graph_pattern;

graph_match_operator: MATCH_SYMBOL hint? graph_pattern;

graph_pattern: graph_path_pattern_list where_clause?;

graph_path_pattern_list:
	graph_path_pattern (COMMA_SYMBOL hint? graph_path_pattern)*;

graph_path_pattern:
	opt_path_variable_assignment? opt_graph_search_prefix? opt_graph_path_mode_prefix?
		graph_path_pattern_expr;

graph_path_pattern_expr:
	graph_path_factor (hint? graph_path_factor)*;

graph_path_factor:
	graph_path_primary
	| graph_quantified_path_primary;

graph_quantified_path_primary:
	graph_path_primary LC_BRACKET_SYMBOL int_literal_or_parameter? COMMA_SYMBOL
		int_literal_or_parameter RC_BRACKET_SYMBOL
	| graph_path_primary LC_BRACKET_SYMBOL int_literal_or_parameter RC_BRACKET_SYMBOL;

graph_path_primary:
	graph_element_pattern
	| graph_parenthesized_path_pattern;

graph_parenthesized_path_pattern:
	LR_BRACKET_SYMBOL hint? graph_path_pattern where_clause? RR_BRACKET_SYMBOL;

graph_element_pattern: graph_node_pattern | graph_edge_pattern;

graph_edge_pattern:
	LT_OPERATOR? MINUS_OPERATOR LS_BRACKET_SYMBOL graph_element_pattern_filler RS_BRACKET_SYMBOL
		MINUS_OPERATOR
	| MINUS_OPERATOR LS_BRACKET_SYMBOL graph_element_pattern_filler RS_BRACKET_SYMBOL
		SUB_GT_BRACKET_SYMBOL
	| MINUS_OPERATOR
	| LT_OPERATOR MINUS_OPERATOR
	| SUB_GT_BRACKET_SYMBOL;

graph_node_pattern:
	LR_BRACKET_SYMBOL graph_element_pattern_filler RR_BRACKET_SYMBOL;

graph_element_pattern_filler:
	// TODO(zp): It is better to avoid using empty production because it confused listener user.
	hint? opt_graph_element_identifier? opt_is_label_expression? graph_property_specification?
	| hint? opt_graph_element_identifier opt_is_label_expression? where_clause
	| hint? opt_graph_element_identifier opt_is_label_expression? graph_property_specification
		where_clause;

graph_property_specification:
	LC_BRACKET_SYMBOL graph_property_name_and_value (
		COMMA_SYMBOL graph_property_name_and_value
	)* RC_BRACKET_SYMBOL;

graph_property_name_and_value:
	identifier COLON_SYMBOL expression;

opt_is_label_expression:
	IS_SYMBOL label_expression
	| COLON_SYMBOL label_expression;

label_expression:
	label_primary
	| label_expression BIT_AND_SYMBOL label_expression
	| label_expression STROKE_SYMBOL label_expression
	| EXCLAMATION_OPERATOR label_expression;

label_primary:
	identifier
	| MODULO_OPERATOR
	| parenthesized_label_expression;

parenthesized_label_expression:
	LR_BRACKET_SYMBOL label_expression RR_BRACKET_SYMBOL;

opt_graph_element_identifier: graph_identifier;

opt_graph_path_mode_prefix: opt_graph_path_mode path_or_paths?;

path_or_paths: PATH_SYMBOL | PATHS_SYMBOL;

opt_graph_path_mode:
	WALK_SYMBOL
	| TRAIL_SYMBOL
	| SIMPLE_SYMBOL
	| ACYCLIC_SYMBOL;

opt_graph_search_prefix:
	(ANY_SYMBOL | ALL_SYMBOL) SHORTEST_SYMBOL?;

opt_path_variable_assignment: graph_identifier EQUAL_OPERATOR;

graph_identifier:
	token_identifier
	| common_keyword_as_identifier;

graph_return_operator:
	RETURN_SYMBOL hint? all_or_distinct? graph_return_item_list group_by_clause?
		graph_order_by_clause? graph_page_clause?;

graph_page_clause:
	OFFSET_SYMBOL possibly_cast_int_literal_or_parameter LIMIT_SYMBOL
		possibly_cast_int_literal_or_parameter
	| SKIP_SYMBOL possibly_cast_int_literal_or_parameter LIMIT_SYMBOL
		possibly_cast_int_literal_or_parameter
	| OFFSET_SYMBOL possibly_cast_int_literal_or_parameter
	| SKIP_SYMBOL possibly_cast_int_literal_or_parameter
	| LIMIT_SYMBOL possibly_cast_int_literal_or_parameter;

graph_order_by_clause:
	ORDER_SYMBOL hint? BY_SYMBOL graph_ordering_expression (
		COMMA_SYMBOL graph_ordering_expression
	)*;

graph_ordering_expression:
	expression collate_clause? opt_graph_asc_or_desc? null_order?;

opt_graph_asc_or_desc:
	asc_or_desc
	| ASCENDING_SYMBOL
	| DESCENDING_SYMBOL;

graph_return_item_list:
	graph_return_item (COMMA_SYMBOL graph_return_item)*;

graph_return_item:
	expression (AS_SYMBOL identifier)?
	| MULTIPLY_OPERATOR;

undrop_statement:
	UNDROP_SYMBOL schema_object_kind opt_if_not_exists? path_expression opt_at_system_time?
		opt_options_list?;

module_statement:
	MODULE_SYMBOL path_expression opt_options_list?;

import_statement:
	IMPORT_SYMBOL import_type path_expression_or_string opt_as_or_into_alias? opt_options_list?;

opt_as_or_into_alias: (AS_SYMBOL | INTO_SYMBOL) identifier;

path_expression_or_string: path_expression | string_literal;

import_type: MODULE_SYMBOL | PROTO_SYMBOL;

call_statement:
	CALL_SYMBOL path_expression LR_BRACKET_SYMBOL (
		tvf_argument (COMMA_SYMBOL tvf_argument)*
	)? RR_BRACKET_SYMBOL;

drop_statement:
	DROP_SYMBOL PRIVILEGE_SYMBOL RESTRICTION_SYMBOL opt_if_exists? ON_SYMBOL privilege_list
		ON_SYMBOL identifier path_expression
	| DROP_SYMBOL ROW_SYMBOL ACCESS_SYMBOL POLICY_SYMBOL opt_if_exists? identifier
		on_path_expression
	| DROP_SYMBOL index_type INDEX_SYMBOL opt_if_exists? path_expression on_path_expression?
	| /* TODO(zp): Refine syntax error */ DROP_SYMBOL table_or_table_function opt_if_exists?
		maybe_dashed_path_expression opt_function_parameters?
	| DROP_SYMBOL SNAPSHOT_SYMBOL TABLE_SYMBOL opt_if_exists? maybe_dashed_path_expression
	| DROP_SYMBOL generic_entity_type opt_if_exists? path_expression
	| DROP_SYMBOL schema_object_kind opt_if_exists? path_expression opt_function_parameters?
		opt_drop_mode?;

opt_drop_mode: RESTRICT_SYMBOL | CASCADE_SYMBOL;

drop_all_row_access_policies_statement:
	DROP_SYMBOL ALL_SYMBOL ROW_SYMBOL ACCESS_SYMBOL? POLICIES_SYMBOL ON_SYMBOL path_expression;

show_statement:
	SHOW_SYMBOL show_target opt_from_path_expression? opt_like_string_literal?;

opt_like_string_literal: LIKE_SYMBOL string_literal;

show_target: MATERIALIZED_SYMBOL VIEWS_SYMBOL | identifier;

rename_statement:
	RENAME_SYMBOL identifier path_expression TO_SYMBOL path_expression;

revoke_statement:
	REVOKE_SYMBOL privileges ON_SYMBOL (identifier identifier?)? path_expression FROM_SYMBOL
		grantee_list;

grant_statement:
	GRANT_SYMBOL privileges ON_SYMBOL (identifier identifier?)? path_expression TO_SYMBOL
		grantee_list;

privileges: ALL_SYMBOL PRIVILEGES_SYMBOL? | privilege_list;

export_metadata_statement:
	EXPORT_SYMBOL table_or_table_function METADATA_SYMBOL FROM_SYMBOL maybe_dashed_path_expression
		with_connection_clause? opt_options_list?;

export_model_statement:
	EXPORT_SYMBOL MODEL_SYMBOL path_expression with_connection_clause? opt_options_list?;

export_data_statement: export_data_no_query as_query;

export_data_no_query:
	EXPORT_SYMBOL DATA_SYMBOL with_connection_clause? opt_options_list?;

explain_statement: EXPLAIN_SYMBOL unterminated_sql_statement;

execute_immediate:
	EXECUTE_SYMBOL IMMEDIATE_SYMBOL expression opt_execute_into_clause? opt_execute_using_clause?;

opt_execute_into_clause: INTO_SYMBOL identifier_list;

opt_execute_using_clause:
	USING_SYMBOL execute_using_argument_list;

execute_using_argument_list:
	execute_using_argument (COMMA_SYMBOL execute_using_argument)*;

execute_using_argument: expression (AS_SYMBOL identifier)?;

describe_statement: describe_keyword describe_info;

describe_info:
	identifier? maybe_slashed_or_dashed_path_expression opt_from_path_expression?;

opt_from_path_expression:
	FROM_SYMBOL maybe_slashed_or_dashed_path_expression;

describe_keyword: DESCRIBE_SYMBOL | DESC_SYMBOL;

define_table_statement:
	DEFINE_SYMBOL TABLE_SYMBOL path_expression options_list?;

create_entity_statement:
	CREATE_SYMBOL opt_or_replace? generic_entity_type opt_if_not_exists? path_expression
		opt_options_list? opt_generic_entity_body?;

opt_generic_entity_body: AS_SYMBOL generic_entity_body;

create_view_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? RECURSIVE_SYMBOL? VIEW_SYMBOL opt_if_not_exists?
		maybe_dashed_path_expression column_with_options_list? opt_sql_security_clause?
		opt_options_list? as_query
	| CREATE_SYMBOL opt_or_replace? MATERIALIZED_SYMBOL RECURSIVE_SYMBOL? VIEW_SYMBOL
		opt_if_not_exists? maybe_dashed_path_expression column_with_options_list?
		opt_sql_security_clause? partition_by_clause_prefix_no_hint?
		cluster_by_clause_prefix_no_hint? opt_options_list? AS_SYMBOL query_or_replica_source
	| CREATE_SYMBOL opt_or_replace? APPROX_SYMBOL RECURSIVE_SYMBOL? VIEW_SYMBOL opt_if_not_exists?
		maybe_dashed_path_expression column_with_options_list? opt_sql_security_clause?
		opt_options_list? as_query;

query_or_replica_source:
	query
	| REPLICA_SYMBOL OF_SYMBOL maybe_dashed_path_expression;

column_with_options_list:
	LR_BRACKET_SYMBOL column_with_options (
		COMMA_SYMBOL column_with_options
	)* RR_BRACKET_SYMBOL;

column_with_options: identifier opt_options_list?;

create_table_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? TABLE_SYMBOL opt_if_not_exists?
		maybe_dashed_path_expression table_element_list? opt_spanner_table_options?
		opt_like_path_expression? opt_clone_table? opt_copy_table? opt_default_collate_clause?
		partition_by_clause_prefix_no_hint? cluster_by_clause_prefix_no_hint? opt_ttl_clause?
		with_connection_clause? opt_options_list? as_query?;

opt_ttl_clause:
	ROW_SYMBOL DELETION_SYMBOL POLICY_SYMBOL LR_BRACKET_SYMBOL expression RR_BRACKET_SYMBOL;

opt_copy_table: COPY_SYMBOL copy_data_source;

copy_data_source:
	maybe_dashed_path_expression opt_at_system_time? where_clause?;

opt_clone_table: CLONE_SYMBOL clone_data_source;

opt_spanner_table_options:
	spanner_primary_key opt_spanner_interleave_in_parent_clause?;

opt_spanner_interleave_in_parent_clause:
	COMMA_SYMBOL INTERLEAVE_SYMBOL IN_SYMBOL PARENT_SYMBOL maybe_dashed_path_expression
		foreign_key_on_delete;

spanner_primary_key:
	PRIMARY_SYMBOL KEY_SYMBOL primary_key_element_list;

create_table_function_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? TABLE_SYMBOL FUNCTION_SYMBOL opt_if_not_exists?
		path_expression opt_function_parameters? opt_returns? opt_sql_security_clause?
		unordered_language_options? opt_as_query_or_string? {
			if localctx.Opt_function_parameters() == nil {
				p.NotifyErrorListeners("Syntax error: Expected (", nil, nil)
			}
		};

opt_as_query_or_string: as_query | AS_SYMBOL string_literal;

unordered_language_options:
	language opt_options_list?
	| opt_options_list language?;

opt_function_parameters:
	LR_BRACKET_SYMBOL (
		function_parameter (COMMA_SYMBOL function_parameter)*
	)? RR_BRACKET_SYMBOL;

create_snapshot_statement:
	CREATE_SYMBOL opt_or_replace? SNAPSHOT_SYMBOL (
		TABLE_SYMBOL
		| schema_object_kind
	) opt_if_not_exists? maybe_dashed_path_expression CLONE_SYMBOL clone_data_source
		opt_options_list?;

create_external_schema_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? EXTERNAL_SYMBOL SCHEMA_SYMBOL opt_if_not_exists?
		path_expression with_connection_clause? opt_options_list;

create_schema_statement:
	CREATE_SYMBOL opt_or_replace? SCHEMA_SYMBOL opt_if_not_exists? path_expression
		opt_default_collate_clause? opt_options_list?;

create_property_graph_statement:
	CREATE_SYMBOL opt_or_replace? PROPERTY_SYMBOL GRAPH_SYMBOL opt_if_not_exists path_expression
		opt_options_list? NODE_SYMBOL TABLES_SYMBOL element_table_list opt_edge_table_clause?;

opt_edge_table_clause:
	EDGE_SYMBOL TABLES_SYMBOL element_table_list;

element_table_list:
	LR_BRACKET_SYMBOL element_table_definition (
		COMMA_SYMBOL element_table_definition
	)* COMMA_SYMBOL? RR_BRACKET_SYMBOL;

element_table_definition:
	path_expression opt_as_alias_with_required_as? opt_key_clause? opt_source_node_table_clause?
		opt_dest_node_table_clause? opt_label_and_properties_clause?;

opt_label_and_properties_clause:
	properties_clause
	| label_and_properties_list;

label_and_properties_list: label_and_properties+;

label_and_properties:
	DEFAULT_SYMBOL? LABEL_SYMBOL identifier properties_clause?;

properties_clause:
	NO_SYMBOL PROPERTIES_SYMBOL
	| properties_all_columns opt_except_column_list?
	| PROPERTIES_SYMBOL LR_BRACKET_SYMBOL derived_property_list RR_BRACKET_SYMBOL;

derived_property_list:
	derived_property (COMMA_SYMBOL derived_property)*;

derived_property: expression opt_as_alias_with_required_as?;

opt_except_column_list: EXCEPT_SYMBOL column_list;

properties_all_columns:
	PROPERTIES_SYMBOL ARE_SYMBOL? ALL_SYMBOL COLUMNS_SYMBOL;

opt_dest_node_table_clause:
	DESTINATION_SYMBOL KEY_SYMBOL column_list REFERENCES_SYMBOL identifier column_list?;

opt_source_node_table_clause:
	SOURCE_SYMBOL KEY_SYMBOL column_list REFERENCES_SYMBOL identifier column_list?;

opt_key_clause: KEY_SYMBOL column_list;

create_model_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? MODEL_SYMBOL opt_if_not_exists? path_expression
		opt_input_output_clause? opt_transform_clause? remote_with_connection_clause?
		opt_options_list? opt_as_query_or_aliased_query_list?;

opt_input_output_clause:
	INPUT_SYMBOL table_element_list OUTPUT_SYMBOL table_element_list;

opt_transform_clause:
	TRANSFORM_SYMBOL LR_BRACKET_SYMBOL select_list RR_BRACKET_SYMBOL;

opt_as_query_or_aliased_query_list:
	as_query
	| AS_SYMBOL LR_BRACKET_SYMBOL aliased_query_list RR_BRACKET_SYMBOL;

aliased_query_list: aliased_query (COMMA_SYMBOL aliased_query)*;

as_query: AS_SYMBOL query;

create_external_table_function_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? EXTERNAL_SYMBOL TABLE_SYMBOL FUNCTION_SYMBOL {
		p.NotifyErrorListeners("Syntax error: CREATE EXTERNAL TABLE FUNCTION is not supported", nil, nil)
	};

create_external_table_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? EXTERNAL_SYMBOL TABLE_SYMBOL opt_if_not_exists?
		maybe_dashed_path_expression table_element_list? opt_like_path_expression?
		opt_default_collate_clause? opt_external_table_with_clauses? opt_options_list?;

opt_default_collate_clause: DEFAULT_SYMBOL collate_clause;

opt_like_path_expression:
	LIKE_SYMBOL maybe_dashed_path_expression;

create_row_access_policy_statement:
	CREATE_SYMBOL opt_or_replace? ROW_SYMBOL ACCESS_SYMBOL? POLICY_SYMBOL opt_if_not_exists?
		identifier? ON_SYMBOL path_expression create_row_access_policy_grant_to_clause?
		filter_using_clause;

filter_using_clause:
	FILTER_SYMBOL? USING_SYMBOL LR_BRACKET_SYMBOL expression RR_BRACKET_SYMBOL;

create_row_access_policy_grant_to_clause:
	grant_to_clause
	| TO_SYMBOL grantee_list;

create_privilege_restriction_statement:
	CREATE_SYMBOL opt_or_replace? PRIVILEGE_SYMBOL RESTRICTION_SYMBOL opt_if_not_exists? ON_SYMBOL
		privilege_list ON_SYMBOL identifier path_expression restrict_to_clause?;

restrict_to_clause:
	RESTRICT_SYMBOL TO_SYMBOL possibly_empty_grantee_list;

possibly_empty_grantee_list:
	LR_BRACKET_SYMBOL (
		string_literal_or_parameter (
			COMMA_SYMBOL string_literal_or_parameter
		)*
	)? RR_BRACKET_SYMBOL;

create_index_statement:
	CREATE_SYMBOL opt_or_replace? UNIQUE_SYMBOL? opt_spanner_null_filtered? index_type? INDEX_SYMBOL
		opt_if_not_exists? path_expression on_path_expression as_alias? index_unnest_expression_list
		? index_order_by_and_options index_storing_list? opt_create_index_statement_suffix?;

opt_create_index_statement_suffix:
	partition_by_clause_prefix_no_hint opt_options_list?
	| opt_options_list? spanner_index_interleave_clause
	| opt_options_list;

spanner_index_interleave_clause:
	COMMA_SYMBOL INTERLEAVE_SYMBOL IN_SYMBOL maybe_dashed_path_expression;

index_storing_list:
	STORING_SYMBOL index_storing_expression_list;

index_storing_expression_list:
	LR_BRACKET_SYMBOL expression (COMMA_SYMBOL expression)* RR_BRACKET_SYMBOL;

index_order_by_and_options:
	LR_BRACKET_SYMBOL column_ordering_and_options_expr (
		COMMA_SYMBOL column_ordering_and_options_expr
	)* RR_BRACKET_SYMBOL
	| index_all_columns;

index_all_columns:
	LR_BRACKET_SYMBOL ALL_SYMBOL COLUMNS_SYMBOL opt_with_column_options? RR_BRACKET_SYMBOL;

opt_with_column_options:
	WITH_SYMBOL COLUMN_SYMBOL OPTIONS_SYMBOL all_column_column_options;

all_column_column_options:
	LR_BRACKET_SYMBOL column_ordering_and_options_expr (
		COMMA_SYMBOL column_ordering_and_options_expr
	)* RR_BRACKET_SYMBOL;

column_ordering_and_options_expr:
	expression collate_clause? asc_or_desc? null_order? opt_options_list?;

index_unnest_expression_list:
	unnest_expression_with_opt_alias_and_offset+;

unnest_expression_with_opt_alias_and_offset:
	unnest_expression as_alias? opt_with_offset_and_alias?;

on_path_expression: ON_SYMBOL path_expression;

index_type: SEARCH_SYMBOL | VECTOR_SYMBOL;

opt_spanner_null_filtered: NULL_FILTERED_SYMBOL;

create_procedure_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? PROCEDURE_SYMBOL opt_if_not_exists?
		path_expression procedure_parameters opt_external_security_clause? with_connection_clause?
		opt_options_list? begin_end_block_or_language_as_code;

begin_end_block_or_language_as_code:
	begin_end_block
	| LANGUAGE_SYMBOL identifier opt_as_code?;

begin_end_block:
	BEGIN_SYMBOL statement_list? opt_exception_handler? END_SYMBOL;

opt_exception_handler:
	EXCEPTION_SYMBOL WHEN_SYMBOL ERROR_SYMBOL THEN_SYMBOL statement_list;

statement_list:
	unterminated_non_empty_statement_list SEMI_SYMBOL;

unterminated_non_empty_statement_list:
	unterminated_statement (SEMI_SYMBOL unterminated_statement)*;

unterminated_statement:
	unterminated_sql_statement
	| unterminated_script_statement;

unterminated_script_statement:
	if_statement
	| case_statement
	| variable_declaration
	| break_statement
	| continue_statement
	| return_statement
	| raise_statement
	| unterminated_unlabeled_script_statement
	| label COLON_SYMBOL unterminated_unlabeled_script_statement identifier?;

label: /* TODO(zp): refine label. */ identifier;

unterminated_unlabeled_script_statement:
	begin_end_block
	| while_statement
	| loop_statement
	| repeat_statement
	| for_in_statement;

for_in_statement:
	FOR_SYMBOL identifier IN_SYMBOL parenthesized_query DO_SYMBOL statement_list? END_SYMBOL
		FOR_SYMBOL;

repeat_statement:
	REPEAT_SYMBOL statement_list? until_clause END_SYMBOL REPEAT_SYMBOL;

until_clause: UNTIL_SYMBOL expression;

loop_statement:
	LOOP_SYMBOL statement_list? END_SYMBOL LOOP_SYMBOL;

while_statement:
	WHILE_SYMBOL expression DO_SYMBOL statement_list? END_SYMBOL WHILE_SYMBOL;

raise_statement:
	RAISE_SYMBOL
	| RAISE_SYMBOL USING_SYMBOL MESSAGE_SYMBOL EQUAL_OPERATOR expression;

return_statement: RETURN_SYMBOL;

continue_statement:
	CONTINUE_SYMBOL identifier?
	| ITERATE_SYMBOL identifier?;

variable_declaration:
	DECLARE_SYMBOL identifier_list type opt_default_expression?
	| DECLARE_SYMBOL identifier_list DEFAULT_SYMBOL expression;

break_statement:
	BREAK_SYMBOL identifier?
	| LEAVE_SYMBOL identifier?;

case_statement:
	CASE_SYMBOL expression? when_then_clauses opt_else? END_SYMBOL CASE_SYMBOL;

when_then_clauses:
	(WHEN_SYMBOL expression THEN_SYMBOL statement_list?)+;

if_statement:
	IF_SYMBOL expression THEN_SYMBOL statement_list? elseif_clauses? opt_else? END_SYMBOL IF_SYMBOL;

elseif_clauses:
	(ELSEIF_SYMBOL expression THEN_SYMBOL statement_list?)+;

opt_else: ELSE_SYMBOL statement_list?;

opt_as_code: AS_SYMBOL string_literal;

opt_external_security_clause:
	EXTERNAL_SYMBOL SECURITY_SYMBOL external_security_clause_kind;

external_security_clause_kind: INVOKER_SYMBOL | DEFINER_SYMBOL;

procedure_parameters:
	LR_BRACKET_SYMBOL (
		procedure_parameter (COMMA_SYMBOL procedure_parameter)*
	)? RR_BRACKET_SYMBOL;

procedure_parameter:
	opt_procedure_parameter_mode? identifier type_or_tvf_schema
	| opt_procedure_parameter_mode? identifier procedure_parameter_termination {
		p.NotifyErrorListeners("Syntax error: Unexpected end of parameter. Parameters should be in the format [<parameter mode>] <parameter name> <type>. If IN/OUT/INOUT is intended to be the name of a parameter, it must be escaped with backticks", nil, nil)
	};

procedure_parameter_termination:
	RR_BRACKET_SYMBOL
	| COMMA_SYMBOL;

opt_procedure_parameter_mode:
	IN_SYMBOL
	| OUT_SYMBOL
	| INOUT_SYMBOL;

create_function_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? opt_aggregate? FUNCTION_SYMBOL opt_if_not_exists
		? function_declaration opt_function_returns? opt_sql_security_clause? opt_determinism_level?
		opt_language_or_remote_with_connection? unordered_options_body?;

opt_determinism_level:
	DETERMINISTIC_SYMBOL
	| NOT_SYMBOL DETERMINISTIC_SYMBOL
	| IMMUTABLE_SYMBOL
	| STABLE_SYMBOL
	| VOLATILE_SYMBOL;

opt_sql_security_clause:
	SQL_SYMBOL SECURITY_SYMBOL sql_security_clause_kind;

sql_security_clause_kind: INVOKER_SYMBOL | DEFINER_SYMBOL;

as_sql_function_body_or_string:
	AS_SYMBOL sql_function_body
	| AS_SYMBOL string_literal;

sql_function_body:
	LR_BRACKET_SYMBOL expression RR_BRACKET_SYMBOL
	| LR_BRACKET_SYMBOL SELECT_SYMBOL {
		p.NotifyErrorListeners("The body of each CREATE FUNCTION statement is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression", nil, nil)
	};

unordered_options_body:
	opt_options_list as_sql_function_body_or_string?
	| as_sql_function_body_or_string opt_options_list?;

opt_language_or_remote_with_connection:
	LANGUAGE_SYMBOL identifier remote_with_connection_clause?
	| remote_with_connection_clause language?;

language: LANGUAGE_SYMBOL identifier;

remote_with_connection_clause:
	REMOTE_SYMBOL with_connection_clause?
	| with_connection_clause;

with_connection_clause: WITH_SYMBOL connection_clause;

opt_function_returns: opt_returns;

opt_returns: RETURNS_SYMBOL type_or_tvf_schema;

function_declaration: path_expression function_parameters;

function_parameters:
	LR_BRACKET_SYMBOL (
		function_parameter (COMMA_SYMBOL function_parameter)*
	)? RR_BRACKET_SYMBOL;

function_parameter:
	identifier type_or_tvf_schema opt_as_alias_with_required_as? opt_default_expression?
		opt_not_aggregate?
	| type_or_tvf_schema opt_as_alias_with_required_as? opt_not_aggregate?;

opt_not_aggregate: NOT_SYMBOL AGGREGATE_SYMBOL;

opt_default_expression: DEFAULT_SYMBOL expression;

type_or_tvf_schema:
	type
	| templated_parameter_type
	| tvf_schema;

tvf_schema:
	TABLE_SYMBOL template_type_open tvf_schema_column (
		COMMA_SYMBOL tvf_schema_column
	)* template_type_close;

tvf_schema_column: identifier type | type;

templated_parameter_type: ANY_SYMBOL templated_parameter_kind;

templated_parameter_kind:
	PROTO_SYMBOL
	| ENUM_SYMBOL
	| STRUCT_SYMBOL
	| ARRAY_SYMBOL
	| identifier;

opt_aggregate: AGGREGATE_SYMBOL;

create_database_statement:
	CREATE_SYMBOL DATABASE_SYMBOL path_expression opt_options_list?;

create_connection_statement:
	CREATE_SYMBOL opt_or_replace? CONNECTION_SYMBOL opt_if_not_exists? path_expression
		opt_options_list?;

create_constant_statement:
	CREATE_SYMBOL opt_or_replace? opt_create_scope? CONSTANT_SYMBOL opt_if_not_exists?
		path_expression EQUAL_OPERATOR expression;

opt_or_replace: OR_SYMBOL REPLACE_SYMBOL;

opt_create_scope:
	TEMP_SYMBOL
	| TEMPORARY_SYMBOL
	| PUBLIC_SYMBOL
	| PRIVATE_SYMBOL;

run_batch_statement: RUN_SYMBOL BATCH_SYMBOL;

abort_batch_statement: ABORT_SYMBOL BATCH_SYMBOL;

start_batch_statement: START_SYMBOL BATCH_SYMBOL identifier?;

rollback_statement: ROLLBACK_SYMBOL TRANSACTION_SYMBOL?;

commit_statement: COMMIT_SYMBOL TRANSACTION_SYMBOL?;

set_statement:
	SET_SYMBOL TRANSACTION_SYMBOL transaction_mode_list
	| SET_SYMBOL identifier EQUAL_OPERATOR expression
	| SET_SYMBOL named_parameter_expression EQUAL_OPERATOR expression
	| SET_SYMBOL system_variable_expression EQUAL_OPERATOR expression
	| SET_SYMBOL LR_BRACKET_SYMBOL identifier_list RR_BRACKET_SYMBOL EQUAL_OPERATOR expression
	| SET_SYMBOL identifier COMMA_SYMBOL identifier EQUAL_OPERATOR {
		p.NotifyErrorListeners("Using SET with multiple variable required parentheses around the variable list", nil, nil)
	};

identifier_list: identifier (COMMA_SYMBOL identifier)*;

begin_statement:
	begin_transaction_keywords transaction_mode_list?;

begin_transaction_keywords:
	START_SYMBOL TRANSACTION_SYMBOL
	| BEGIN_SYMBOL TRANSACTION_SYMBOL?;

transaction_mode_list:
	transaction_mode (COMMA_SYMBOL transaction_mode)*;

transaction_mode:
	READ_SYMBOL ONLY_SYMBOL
	| READ_SYMBOL WRITE_SYMBOL
	| ISOLATION_SYMBOL LEVEL_SYMBOL identifier
	| ISOLATION_SYMBOL LEVEL_SYMBOL identifier identifier;

truncate_statement:
	TRUNCATE_SYMBOL TABLE_SYMBOL maybe_dashed_path_expression opt_where_expression?;

merge_statement:
	MERGE_SYMBOL INTO_SYMBOL? maybe_dashed_path_expression as_alias? USING_SYMBOL merge_source
		ON_SYMBOL expression (merge_when_clause)+;

merge_source: table_path_expression | table_subquery;

merge_when_clause:
	WHEN_SYMBOL MATCHED_SYMBOL opt_and_expression? THEN_SYMBOL merge_action
	| WHEN_SYMBOL NOT_SYMBOL MATCHED_SYMBOL by_target? opt_and_expression? THEN_SYMBOL merge_action
	| WHEN_SYMBOL NOT_SYMBOL MATCHED_SYMBOL BY_SYMBOL SOURCE_SYMBOL opt_and_expression? THEN_SYMBOL
		merge_action;

merge_action:
	INSERT_SYMBOL column_list? merge_insert_value_list_or_source_row
	| UPDATE_SYMBOL SET_SYMBOL update_item_list
	| DELETE_SYMBOL;

merge_insert_value_list_or_source_row:
	VALUES_SYMBOL insert_values_row
	| ROW_SYMBOL;

by_target: BY_SYMBOL TARGET_SYMBOL;

opt_and_expression: AND_SYMBOL expression;

statement_level_hint: hint;

// query_statement: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax
query_statement: query;

dml_statement:
	insert_statement
	| delete_statement
	| update_statement;

update_statement:
	UPDATE_SYMBOL maybe_dashed_generalized_path_expression hint? as_alias? opt_with_offset_and_alias
		? SET_SYMBOL update_item_list from_clause? opt_where_expression? opt_assert_rows_modified?
		opt_returning_clause?;

delete_statement:
	DELETE_SYMBOL FROM_SYMBOL? maybe_dashed_generalized_path_expression hint? as_alias?
		opt_with_offset_and_alias? opt_where_expression? opt_assert_rows_modified?
		opt_returning_clause?;

insert_statement:
	insert_statement_prefix column_list? insert_values_or_query opt_assert_rows_modified?
		opt_returning_clause?
	| insert_statement_prefix column_list? insert_values_list_or_table_clause on_conflict_clause
		opt_assert_rows_modified? opt_returning_clause?
	| insert_statement_prefix column_list? LR_BRACKET_SYMBOL query RR_BRACKET_SYMBOL
		on_conflict_clause opt_assert_rows_modified? opt_returning_clause?;

on_conflict_clause:
	ON_SYMBOL CONFLICT_SYMBOL opt_conflict_target? DO_SYMBOL NOTHING_SYMBOL
	| ON_SYMBOL CONFLICT_SYMBOL opt_conflict_target? DO_SYMBOL UPDATE_SYMBOL SET_SYMBOL
		update_item_list opt_where_expression?;

opt_where_expression: WHERE_SYMBOL expression;

opt_conflict_target:
	column_list
	| ON_SYMBOL UNIQUE_SYMBOL CONSTRAINT_SYMBOL identifier;

update_item_list: update_item (COMMA_SYMBOL update_item)*;

update_item: update_set_value | nested_dml_statement;

update_set_value:
	generalized_path_expression EQUAL_OPERATOR expression_or_default;

nested_dml_statement:
	LR_BRACKET_SYMBOL dml_statement RR_BRACKET_SYMBOL;

insert_values_list_or_table_clause:
	insert_values_list
	| table_clause_unreversed;

table_clause_unreversed: TABLE_SYMBOL table_clause_no_keyword;

table_clause_no_keyword:
	path_expression where_clause?
	| tvf_with_suffixes where_clause?;

opt_returning_clause:
	THEN_SYMBOL RETURN_SYMBOL select_list
	| THEN_SYMBOL RETURN_SYMBOL WITH_SYMBOL ACTION_SYMBOL select_list
	| THEN_SYMBOL RETURN_SYMBOL WITH_SYMBOL ACTION_SYMBOL AS_SYMBOL identifier select_list;

opt_assert_rows_modified:
	ASSERT_ROWS_MODIFIED_SYMBOL possibly_cast_int_literal_or_parameter;

insert_values_or_query: insert_values_list | query;

insert_values_list:
	VALUES_SYMBOL insert_values_row (
		COMMA_SYMBOL insert_values_row
	)*;

insert_values_row:
	LR_BRACKET_SYMBOL expression_or_default (
		COMMA_SYMBOL expression_or_default
	)* RR_BRACKET_SYMBOL;

expression_or_default: expression | DEFAULT_SYMBOL;

insert_statement_prefix:
	INSERT_SYMBOL opt_or_ignore_replace_update? opt_into? maybe_dashed_generalized_path_expression
		hint?;

maybe_dashed_generalized_path_expression:
	generalized_path_expression
	| dashed_path_expression;

opt_into: INTO_SYMBOL;

opt_or_ignore_replace_update:
	OR_SYMBOL IGNORE_SYMBOL
	| IGNORE_SYMBOL
	| OR_SYMBOL REPLACE_SYMBOL
	| REPLACE_SYMBOL
	| OR_SYMBOL UPDATE_SYMBOL
	| UPDATE_SYMBOL;

alter_statement:
	ALTER_SYMBOL table_or_table_function opt_if_exists? maybe_dashed_path_expression
		alter_action_list
	| ALTER_SYMBOL schema_object_kind opt_if_exists? path_expression alter_action_list
	| ALTER_SYMBOL generic_entity_type opt_if_exists? path_expression alter_action_list
	| ALTER_SYMBOL generic_entity_type opt_if_exists? alter_action_list
	| ALTER_SYMBOL PRIVILEGE_SYMBOL RESTRICTION_SYMBOL opt_if_exists? ON_SYMBOL privilege_list
		ON_SYMBOL identifier path_expression
	| ALTER_SYMBOL ROW_SYMBOL ACCESS_SYMBOL POLICY_SYMBOL opt_if_exists? identifier ON_SYMBOL
		path_expression row_access_policy_alter_action_list
	| ALTER_SYMBOL ALL_SYMBOL ROW_SYMBOL ACCESS_SYMBOL POLICIES_SYMBOL ON_SYMBOL path_expression
		row_access_policy_alter_action;

analyze_statement:
	ANALYZE_SYMBOL opt_options_list? table_and_column_info_list?;

assert_statement: ASSERT_SYMBOL expression opt_description?;

aux_load_data_statement:
	LOAD_SYMBOL DATA_SYMBOL append_or_overwrite maybe_dashed_path_expression_with_scope
		table_element_list? load_data_partitions_clause? collate_clause?
		partition_by_clause_prefix_no_hint? cluster_by_clause_prefix_no_hint? opt_options_list?
		aux_load_data_from_files_options_list opt_external_table_with_clauses?;

clone_data_statement:
	CLONE_SYMBOL DATA_SYMBOL INTO_SYMBOL maybe_dashed_path_expression FROM_SYMBOL
		clone_data_source_list;

clone_data_source_list:
	clone_data_source (UNION_SYMBOL ALL_SYMBOL clone_data_source)*;

clone_data_source:
	maybe_dashed_path_expression opt_at_system_time? where_clause?;

opt_external_table_with_clauses:
	with_partition_columns_clause with_connection_clause
	| with_partition_columns_clause
	| with_connection_clause;

with_partition_columns_clause:
	WITH_SYMBOL PARTITION_SYMBOL COLUMNS_SYMBOL table_element_list?;

aux_load_data_from_files_options_list:
	FROM_SYMBOL FILES_SYMBOL options_list;

cluster_by_clause_prefix_no_hint:
	CLUSTER_SYMBOL BY_SYMBOL expression (COMMA_SYMBOL expression)*;

load_data_partitions_clause:
	OVERWRITE_SYMBOL? PARTITIONS_SYMBOL LR_BRACKET_SYMBOL expression RR_BRACKET_SYMBOL;

maybe_dashed_path_expression_with_scope:
	TEMP_SYMBOL TABLE_SYMBOL maybe_dashed_path_expression
	| TEMPORARY_SYMBOL TABLE_SYMBOL maybe_dashed_path_expression
	| maybe_dashed_path_expression;

table_element_list:
	LR_BRACKET_SYMBOL (
		table_element (COMMA_SYMBOL table_element)* COMMA_SYMBOL?
	)? RR_BRACKET_SYMBOL;

table_element:
	table_column_definition
	| table_constraint_definition;

table_constraint_definition:
	primary_key_spec
	| table_constraint_spec
	| identifier identifier table_constraint_spec;

append_or_overwrite: INTO_SYMBOL | OVERWRITE_SYMBOL;

opt_description: AS_SYMBOL string_literal;

table_and_column_info_list:
	table_and_column_info (COMMA_SYMBOL table_and_column_info)*;

table_and_column_info:
	maybe_dashed_path_expression column_list?;

row_access_policy_alter_action_list:
	row_access_policy_alter_action (
		COMMA_SYMBOL row_access_policy_alter_action
	)*;

row_access_policy_alter_action:
	grant_to_clause
	| FILTER_SYMBOL USING_SYMBOL LR_BRACKET_SYMBOL expression RR_BRACKET_SYMBOL
	| REVOKE_SYMBOL FROM_SYMBOL LR_BRACKET_SYMBOL grantee_list RR_BRACKET_SYMBOL
	| REVOKE_SYMBOL FROM_SYMBOL ALL_SYMBOL
	| RENAME_SYMBOL TO_SYMBOL identifier;

grant_to_clause:
	GRANT_SYMBOL TO_SYMBOL LR_BRACKET_SYMBOL grantee_list RR_BRACKET_SYMBOL;

grantee_list:
	string_literal_or_parameter (
		COMMA_SYMBOL string_literal_or_parameter
	)*;

privilege_list: privilege (COMMA_SYMBOL privilege)*;

privilege: privilege_name path_expression_list_with_parens?;

path_expression_list_with_parens:
	LR_BRACKET_SYMBOL path_expression_list RR_BRACKET_SYMBOL;

privilege_name: identifier | SELECT_SYMBOL;

generic_entity_type: generic_entity_type_unchecked;

generic_entity_type_unchecked: IDENTIFIER | PROJECT_SYMBOL;

schema_object_kind:
	AGGREGATE_SYMBOL FUNCTION_SYMBOL
	| APPROX_SYMBOL VIEW_SYMBOL
	| CONNECTION_SYMBOL
	| CONSTANT_SYMBOL
	| DATABASE_SYMBOL
	| EXTERNAL_SYMBOL table_or_table_function
	| EXTERNAL_SYMBOL SCHEMA_SYMBOL
	| FUNCTION_SYMBOL
	| INDEX_SYMBOL
	| MATERIALIZED_SYMBOL VIEW_SYMBOL
	| MODEL_SYMBOL
	| PROCEDURE_SYMBOL
	| SCHEMA_SYMBOL
	| VIEW_SYMBOL;

alter_action_list: alter_action (COMMA_SYMBOL alter_action)*;

alter_action:
	SET_SYMBOL OPTIONS_SYMBOL options_list
	| SET_SYMBOL AS_SYMBOL generic_entity_body
	| ADD_SYMBOL table_constraint_spec
	| ADD_SYMBOL primary_key_spec
	| ADD_SYMBOL CONSTRAINT_SYMBOL opt_if_not_exists? identifier
		primary_key_or_table_constraint_spec
	| DROP_SYMBOL CONSTRAINT_SYMBOL opt_if_exists? identifier
	| DROP_SYMBOL PRIMARY_SYMBOL KEY_SYMBOL opt_if_exists?
	| ALTER_SYMBOL CONSTRAINT_SYMBOL opt_if_exists? identifier constraint_enforcement
	| ALTER_SYMBOL CONSTRAINT_SYMBOL opt_if_exists? identifier SET_SYMBOL OPTIONS_SYMBOL
		options_list
	| ADD_SYMBOL COLUMN_SYMBOL opt_if_not_exists? table_column_definition column_position?
		fill_using_expression?
	| DROP_SYMBOL COLUMN_SYMBOL opt_if_exists? identifier
	| RENAME_SYMBOL COLUMN_SYMBOL opt_if_exists? identifier TO_SYMBOL identifier
	| ALTER_SYMBOL COLUMN_SYMBOL opt_if_exists? identifier SET_SYMBOL DATA_SYMBOL TYPE_SYMBOL
		field_schema
	| ALTER_SYMBOL COLUMN_SYMBOL opt_if_exists? identifier SET_SYMBOL OPTIONS_SYMBOL options_list
	| ALTER_SYMBOL COLUMN_SYMBOL opt_if_exists? identifier SET_SYMBOL DEFAULT_SYMBOL expression
	| ALTER_SYMBOL COLUMN_SYMBOL opt_if_exists? identifier DROP_SYMBOL DEFAULT_SYMBOL
	| ALTER_SYMBOL COLUMN_SYMBOL opt_if_exists? identifier DROP_SYMBOL NOT_SYMBOL NULL_SYMBOL
	| ALTER_SYMBOL COLUMN_SYMBOL opt_if_exists? identifier DROP_SYMBOL GENERATED_SYMBOL
	| RENAME_SYMBOL TO_SYMBOL path_expression
	| SET_SYMBOL DEFAULT_SYMBOL collate_clause
	| ADD_SYMBOL ROW_SYMBOL DELETION_SYMBOL POLICY_SYMBOL opt_if_not_exists? LR_BRACKET_SYMBOL
		expression RR_BRACKET_SYMBOL
	| REPLACE_SYMBOL ROW_SYMBOL DELETION_SYMBOL POLICY_SYMBOL opt_if_exists? LR_BRACKET_SYMBOL
		expression RR_BRACKET_SYMBOL
	| DROP_SYMBOL ROW_SYMBOL DELETION_SYMBOL POLICY_SYMBOL opt_if_exists?
	| ALTER_SYMBOL generic_sub_entity_type opt_if_exists? identifier alter_action
	| ADD_SYMBOL generic_sub_entity_type opt_if_not_exists? identifier
	| DROP_SYMBOL generic_sub_entity_type opt_if_exists? identifier
	| spanner_alter_column_action
	| spanner_set_on_delete_action;

spanner_set_on_delete_action:
	SET_SYMBOL ON_SYMBOL DELETE_SYMBOL foreign_key_action;

spanner_alter_column_action:
	ALTER_SYMBOL COLUMN_SYMBOL opt_if_exists? identifier column_schema_inner
		not_null_column_attribute? spanner_generated_or_default? opt_options_list?;

spanner_generated_or_default:
	AS_SYMBOL LR_BRACKET_SYMBOL expression RR_BRACKET_SYMBOL STORED_SYMBOL;

generic_sub_entity_type: sub_entity_type_identifier;

sub_entity_type_identifier: IDENTIFIER | REPLICA_SYMBOL;

fill_using_expression: FILL_SYMBOL USING_SYMBOL expression;

column_position:
	PRECEDING_SYMBOL identifier
	| FOLLOWING_SYMBOL identifier;

table_column_definition:
	identifier table_column_schema column_attributes? opt_options_list?;

column_attributes: column_attribute+ constraint_enforcement?;

column_attribute:
	primary_key_column_attribute
	| foreign_key_column_attribute
	| hidden_column_attribute
	| not_null_column_attribute;

primary_key_column_attribute: PRIMARY_SYMBOL KEY_SYMBOL;

foreign_key_column_attribute:
	opt_constraint_identity? foreign_key_reference;

hidden_column_attribute: HIDDEN_SYMBOL;

opt_constraint_identity: CONSTRAINT_SYMBOL identifier;

table_column_schema:
	column_schema_inner collate_clause? opt_column_info?
	| generated_column_info;

opt_column_info:
	generated_column_info invalid_default_column? {
		if localctx.Invalid_default_column() != nil {
			p.NotifyErrorListeners("Syntax error: \"DEFAULT\" and \"GENERATED ALWAYS AS\" clauses must not be both provided for the column", nil, nil)
		}
	}
	| default_column_info invalid_generated_column? {
		if localctx.Invalid_generated_column() != nil {
			p.NotifyErrorListeners("Syntax error: \"DEFAULT\" and \"GENERATED ALWAYS AS\" clauses must not be both provided for the column", nil, nil)
		}
	};

invalid_generated_column: generated_column_info;

invalid_default_column: default_column_info;

default_column_info: DEFAULT_SYMBOL expression;

generated_column_info:
	generated_mode LR_BRACKET_SYMBOL expression RR_BRACKET_SYMBOL stored_mode?
	| generated_mode identity_column_info;

identity_column_info:
	IDENTITY_SYMBOL LR_BRACKET_SYMBOL opt_start_with? opt_increment_by? opt_maxvalue? opt_minvalue?
		opt_cycle? RR_BRACKET_SYMBOL;

opt_start_with: START_SYMBOL WITH_SYMBOL signed_numeric_literal;

opt_increment_by:
	INCREMENT_SYMBOL BY_SYMBOL signed_numeric_literal;

opt_maxvalue: MAXVALUE_SYMBOL signed_numeric_literal;

opt_minvalue: MINVALUE_SYMBOL signed_numeric_literal;

opt_cycle: CYCLE_SYMBOL | NO_SYMBOL CYCLE_SYMBOL;

signed_numeric_literal:
	integer_literal
	| numeric_literal
	| bignumeric_literal
	| floating_point_literal
	| MINUS_OPERATOR integer_literal
	| MINUS_OPERATOR floating_point_literal;

// All rules reference stored_mode should make stored_mode optional.
stored_mode: STORED_SYMBOL VOLATILE_SYMBOL | STORED_SYMBOL;

generated_mode:
	GENERATED_SYMBOL AS_SYMBOL
	| GENERATED_SYMBOL ALWAYS_SYMBOL AS_SYMBOL
	| GENERATED_SYMBOL BY_SYMBOL DEFAULT_SYMBOL AS_SYMBOL
	| AS_SYMBOL;

column_schema_inner:
	raw_column_schema_inner opt_type_parameters?;

raw_column_schema_inner:
	simple_column_schema_inner
	| array_column_schema_inner
	| struct_column_schema_inner
	| range_column_schema_inner;

range_column_schema_inner:
	RANGE_SYMBOL template_type_open field_schema template_type_close;

struct_column_schema_inner:
	STRUCT_SYMBOL template_type_open (
		struct_column_field (COMMA_SYMBOL struct_column_field)*
	)? template_type_close;

struct_column_field:
	column_schema_inner collate_clause? opt_field_attributes?
	| identifier field_schema;

simple_column_schema_inner: path_expression | INTERVAL_SYMBOL;

array_column_schema_inner:
	ARRAY_SYMBOL template_type_open field_schema template_type_close;

field_schema:
	column_schema_inner collate_clause? opt_field_attributes? opt_options_list?;

opt_field_attributes: not_null_column_attribute;

not_null_column_attribute: NOT_SYMBOL NULL_SYMBOL;

primary_key_or_table_constraint_spec:
	primary_key_spec
	| table_constraint_spec;

opt_if_not_exists: IF_SYMBOL NOT_SYMBOL EXISTS_SYMBOL;

primary_key_spec:
	PRIMARY_SYMBOL KEY_SYMBOL primary_key_element_list constraint_enforcement? opt_options_list?;

primary_key_element_list:
	LR_BRACKET_SYMBOL (
		primary_key_element (COMMA_SYMBOL primary_key_element)*
	)? RR_BRACKET_SYMBOL;

primary_key_element: identifier asc_or_desc? null_order?;

table_constraint_spec:
	CHECK_SYMBOL LR_BRACKET_SYMBOL expression RR_BRACKET_SYMBOL constraint_enforcement?
		opt_options_list?
	| FOREIGN_SYMBOL KEY_SYMBOL column_list foreign_key_reference constraint_enforcement?
		opt_options_list?;

foreign_key_reference:
	REFERENCES_SYMBOL path_expression column_list opt_foreign_key_match? opt_foreign_key_action?;

opt_foreign_key_action:
	foreign_key_on_update foreign_key_on_delete?
	| foreign_key_on_delete foreign_key_on_update?;

foreign_key_on_update:
	ON_SYMBOL UPDATE_SYMBOL foreign_key_action;

foreign_key_on_delete:
	ON_SYMBOL DELETE_SYMBOL foreign_key_action;

foreign_key_action:
	NO_SYMBOL ACTION_SYMBOL
	| RESTRICT_SYMBOL
	| CASCADE_SYMBOL
	| SET_SYMBOL NULL_SYMBOL;

opt_foreign_key_match: MATCH_SYMBOL foreign_key_match_mode;

foreign_key_match_mode:
	SIMPLE_SYMBOL
	| FULL_SYMBOL
	| NOT_SYMBOL DISTINCT_SYMBOL;

column_list:
	LR_BRACKET_SYMBOL identifier (COMMA_SYMBOL identifier)* RR_BRACKET_SYMBOL;

opt_options_list: OPTIONS_SYMBOL options_list;

constraint_enforcement: NOT_SYMBOL? ENFORCED_SYMBOL;

generic_entity_body: json_literal | string_literal;

opt_if_exists: IF_SYMBOL EXISTS_SYMBOL;

table_or_table_function: TABLE_SYMBOL FUNCTION_SYMBOL?;

query: query_without_pipe_operators;

query_without_pipe_operators:
	with_clause query_primary_or_set_operation order_by_clause? limit_offset_clause?
	| with_clause_with_trailing_comma select_or_from_keyword {p.NotifyErrorListeners("Syntax error: Trailing comma after the WITH clause before the main query is not allowed", nil, nil)
		}
	| with_clause PIPE_SYMBOL {p.NotifyErrorListeners("Syntax error: A pipe operator cannot follow the WITH clause before the main query; The main query usually starts with SELECT or FROM here", nil, nil)
		}
	| query_primary_or_set_operation order_by_clause? limit_offset_clause?
	| with_clause? from_clause {p.NotifyErrorListeners("Syntax error: Unexpected FROM", nil, nil)}
	// FIXME(zp): Inject the keyword from original input.
	| with_clause? from_clause bad_keyword_after_from_query {p.NotifyErrorListeners("Syntax error: <KEYWORD> not supported after FROM query; Consider using pipe operator `|>` ", nil, nil)
		}
	| with_clause? from_clause bad_keyword_after_from_query_allows_parens {p.NotifyErrorListeners("Syntax error: <KEYWORD> not supported after FROM query; Consider using pipe operator `|>` ", nil, nil)
		};

bad_keyword_after_from_query:
	WHERE_SYMBOL
	| SELECT_SYMBOL
	| GROUP_SYMBOL;

bad_keyword_after_from_query_allows_parens:
	ORDER_SYMBOL
	| UNION_SYMBOL
	| INTERSECT_SYMBOL
	| EXCEPT_SYMBOL
	| LIMIT_SYMBOL;

with_clause_with_trailing_comma: with_clause COMMA_SYMBOL;

select_or_from_keyword: SELECT_SYMBOL | FROM_SYMBOL;

query_primary_or_set_operation:
	query_primary
	| query_set_operation;

query_set_operation: query_set_operation_prefix;

query_set_operation_prefix:
	query_primary query_set_operation_item+
	| query_primary set_operation_metadata FROM_SYMBOL { p.NotifyErrorListeners("Syntax error: Unexpected FROM;FROM queries following a set operation must be parenthesized", nil, nil); 
		}
	| query_set_operation_prefix set_operation_metadata FROM_SYMBOL { p.NotifyErrorListeners("Syntax error: Unexpected FROM;FROM queries following a set operation must be parenthesized", nil, nil); 
		};

query_set_operation_item: set_operation_metadata query_primary;

query_primary:
	select
	| parenthesized_query opt_as_alias_with_required_as?;

set_operation_metadata:
	opt_corresponding_outer_mode? query_set_operation_type hint? all_or_distinct opt_strict?
		opt_column_match_suffix?;

opt_column_match_suffix:
	CORRESPONDING_SYMBOL
	| CORRESPONDING_SYMBOL BY_SYMBOL;

opt_strict: STRICT_SYMBOL;

all_or_distinct: ALL_SYMBOL | DISTINCT_SYMBOL;

query_set_operation_type:
	UNION_SYMBOL
	| EXCEPT_SYMBOL
	| INTERSECT_SYMBOL;

opt_corresponding_outer_mode:
	FULL_SYMBOL opt_outer?
	| OUTER_SYMBOL
	| LEFT_SYMBOL opt_outer?;

opt_outer: OUTER_SYMBOL;

with_clause:
	WITH_SYMBOL RECURSIVE_SYMBOL? aliased_query (
		COMMA_SYMBOL aliased_query
	)*;

aliased_query:
	identifier AS_SYMBOL parenthesized_query opt_aliased_query_modifiers?;

opt_aliased_query_modifiers: recursion_depth_modifier;

recursion_depth_modifier:
	WITH_SYMBOL DEPTH_SYMBOL opt_as_alias_with_required_as?
	| WITH_SYMBOL DEPTH_SYMBOL opt_as_alias_with_required_as? BETWEEN_SYMBOL
		possibly_unbounded_int_literal_or_parameter AND_SYMBOL
		possibly_unbounded_int_literal_or_parameter
	| WITH_SYMBOL DEPTH_SYMBOL opt_as_alias_with_required_as? MAX_SYMBOL
		possibly_unbounded_int_literal_or_parameter;

possibly_unbounded_int_literal_or_parameter:
	int_literal_or_parameter
	| UNBOUNDED_SYMBOL;

int_literal_or_parameter:
	integer_literal
	| parameter_expression
	| system_variable_expression;

order_by_clause: order_by_clause_prefix;

order_by_clause_prefix:
	ORDER_SYMBOL hint? BY_SYMBOL ordering_expression (
		COMMA_SYMBOL ordering_expression
	)*;

ordering_expression:
	expression collate_clause? asc_or_desc? null_order?;

select: select_clause from_clause? opt_clauses_following_from?;

opt_clauses_following_from:
	where_clause group_by_clause? having_clause? qualify_clause_nonreserved? window_clause?
	| opt_clauses_following_where;

opt_clauses_following_where:
	group_by_clause having_clause? qualify_clause_nonreserved? window_clause?
	| opt_clauses_following_group_by;

opt_clauses_following_group_by:
	having_clause qualify_clause_nonreserved? window_clause?
	| qualify_clause_nonreserved window_clause?
	| window_clause;

window_clause: window_clause_prefix;

window_clause_prefix:
	WINDOW_SYMBOL window_definition (
		COMMA_SYMBOL window_definition
	)*;

window_definition: identifier AS_SYMBOL window_specification;

where_clause: WHERE_SYMBOL expression;

having_clause: HAVING_SYMBOL expression;

group_by_clause: group_by_all | group_by_clause_prefix;

group_by_all: group_by_preamble ALL_SYMBOL;

select_clause:
	SELECT_SYMBOL hint? opt_select_with? all_or_distinct? opt_select_as_clause? select_list
	| SELECT_SYMBOL hint? opt_select_with? all_or_distinct? opt_select_as_clause? FROM_SYMBOL {p.NotifyErrorListeners("Syntax error: SELECT list must not be empty", nil, nil)
		};

opt_select_as_clause:
	AS_SYMBOL STRUCT_SYMBOL
	| AS_SYMBOL path_expression;

opt_select_with:
	WITH_SYMBOL identifier
	| WITH_SYMBOL identifier OPTIONS_SYMBOL options_list;

// from_clause: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#from_clause
from_clause: FROM_SYMBOL from_clause_contents;

from_clause_contents:
	table_primary from_clause_contents_suffix*
	| AT_SYMBOL {p.NotifyErrorListeners("Query parameters cannot be used in place of table names",nil,nil)
		}
	| QUESTION_SYMBOL {p.NotifyErrorListeners("Query parameters cannot be used in place of table names",nil,nil)
		}
	| ATAT_SYMBOL {p.NotifyErrorListeners("System variables cannot be used in place of table names",nil,nil)
		};

from_clause_contents_suffix:
	COMMA_SYMBOL table_primary
	| opt_natural? join_type? join_hint? JOIN_SYMBOL hint? table_primary on_or_using_clause_list?;

table_primary:
	tvf_with_suffixes
	| table_path_expression
	| LR_BRACKET_SYMBOL join RR_BRACKET_SYMBOL
	| table_subquery
	| table_primary match_recognize_clause
	| table_primary sample_clause;

tvf_with_suffixes:
	tvf_prefix_no_args RR_BRACKET_SYMBOL hint? pivot_or_unpivot_clause_and_aliases?
	| tvf_prefix RR_BRACKET_SYMBOL hint? pivot_or_unpivot_clause_and_aliases?;

pivot_or_unpivot_clause_and_aliases:
	AS_SYMBOL identifier
	| identifier
	| AS_SYMBOL identifier pivot_clause as_alias?
	| AS_SYMBOL identifier unpivot_clause as_alias?
	| AS_SYMBOL identifier qualify_clause_nonreserved {
				 p.NotifyErrorListeners("QUALIFY clause must be used in conjunction with WHERE or GROUP BY or HAVING clause", nil, nil); 
		}
	| identifier pivot_clause as_alias
	| identifier unpivot_clause as_alias
	| identifier qualify_clause_nonreserved {
				 p.NotifyErrorListeners("QUALIFY clause must be used in conjunction with WHERE or GROUP BY or HAVING clause", nil, nil); 
		}
	| pivot_clause as_alias?
	| unpivot_clause as_alias?
	| qualify_clause_nonreserved {
				 p.NotifyErrorListeners("QUALIFY clause must be used in conjunction with WHERE or GROUP BY or HAVING clause", nil, nil); 
		};

as_alias: AS_SYMBOL? identifier;

sample_clause:
	TABLESAMPLE_SYMBOL identifier LR_BRACKET_SYMBOL sample_size RR_BRACKET_SYMBOL
		opt_sample_clause_suffix;

opt_sample_clause_suffix:
	repeatable_clause
	| WITH_SYMBOL WEIGHT_SYMBOL repeatable_clause?
	| WITH_SYMBOL WEIGHT_SYMBOL identifier repeatable_clause?
	| WITH_SYMBOL WEIGHT_SYMBOL AS_SYMBOL identifier repeatable_clause?;

repeatable_clause:
	REPEATABLE_SYMBOL LR_BRACKET_SYMBOL possibly_cast_int_literal_or_parameter RR_BRACKET_SYMBOL;

possibly_cast_int_literal_or_parameter:
	cast_int_literal_or_parameter
	| int_literal_or_parameter;

cast_int_literal_or_parameter:
	CAST_SYMBOL LR_BRACKET_SYMBOL int_literal_or_parameter AS_SYMBOL type opt_format?
		RR_BRACKET_SYMBOL;

sample_size:
	sample_size_value sample_size_unit partition_by_clause_prefix_no_hint?;

sample_size_value:
	possibly_cast_int_literal_or_parameter
	| floating_point_literal;

sample_size_unit: ROWS_SYMBOL | PERCENT_SYMBOL;

partition_by_clause_prefix_no_hint:
	PARTITION_SYMBOL BY_SYMBOL expression (
		COMMA_SYMBOL expression
	)*;

match_recognize_clause:
	MATCH_RECOGNIZE_SYMBOL LR_BRACKET_SYMBOL partition_by_clause_prefix? order_by_clause
		MEASURES_SYMBOL select_list_prefix_with_as_aliases PATTERN_SYMBOL LR_BRACKET_SYMBOL
		row_pattern_expr RR_BRACKET_SYMBOL DEFINE_SYMBOL with_expression_variable_prefix
		RR_BRACKET_SYMBOL as_alias?;

row_pattern_expr:
	row_pattern_concatenation
	| row_pattern_expr STROKE_SYMBOL row_pattern_concatenation;

row_pattern_concatenation:
	row_pattern_factor
	| row_pattern_concatenation row_pattern_factor;

row_pattern_factor:
	identifier
	| LR_BRACKET_SYMBOL row_pattern_expr RR_BRACKET_SYMBOL;

select_list_prefix_with_as_aliases:
	select_column_expr_with_as_alias (
		COMMA_SYMBOL select_column_expr_with_as_alias
	)*;

select_column_expr_with_as_alias:
	expression AS_SYMBOL identifier;

table_subquery:
	parenthesized_query opt_pivot_or_unpivot_clause_and_alias?;

join: table_primary join_item*;

// join_item resolves the mutually left-recursive for [join, join_input]. join_input: join |
// table_primary;
join_item:
	opt_natural? join_type? join_hint? JOIN_SYMBOL hint? table_primary on_or_using_clause_list?;

on_or_using_clause_list: on_or_using_clause+;

on_or_using_clause: on_clause | using_clause;

using_clause:
	USING_SYMBOL LR_BRACKET_SYMBOL identifier (
		DOT_SYMBOL identifier
	)* RR_BRACKET_SYMBOL;

join_hint: HASH_SYMBOL | LOOKUP_SYMBOL;

table_path_expression:
	table_path_expression_base hint? opt_pivot_or_unpivot_clause_and_alias?
		opt_with_offset_and_alias? opt_at_system_time?;

opt_at_system_time:
	FOR_SYMBOL SYSTEM_SYMBOL TIME_SYMBOL AS_SYMBOL OF_SYMBOL expression
	| FOR_SYMBOL SYSTEM_TIME_SYMBOL AS_SYMBOL OF_SYMBOL expression;

opt_with_offset_and_alias: WITH_SYMBOL OFFSET_SYMBOL as_alias?;

opt_pivot_or_unpivot_clause_and_alias:
	AS_SYMBOL identifier
	| identifier
	| AS_SYMBOL identifier pivot_clause as_alias?
	| AS_SYMBOL identifier unpivot_clause as_alias?
	| AS_SYMBOL identifier qualify_clause_nonreserved {p.NotifyErrorListeners("QUALIFY clause must be used in conjunction with WHERE or GROUP BY or HAVING clause", nil, nil)
		}
	| identifier pivot_clause as_alias?
	| identifier unpivot_clause as_alias?
	| identifier qualify_clause_nonreserved {p.NotifyErrorListeners("QUALIFY clause must be used in conjunction with WHERE or GROUP BY or HAVING clause", nil, nil)
		}
	| pivot_clause as_alias?
	| unpivot_clause as_alias?
	| qualify_clause_nonreserved {p.NotifyErrorListeners("QUALIFY clause must be used in conjunction with WHERE or GROUP BY or HAVING clause", nil, nil)
		};

table_path_expression_base:
	unnest_expression
	| maybe_slashed_or_dashed_path_expression
	| path_expression LS_BRACKET_SYMBOL {p.NotifyErrorListeners("Syntax error: Array element access is not allowed in the FROM clause without UNNEST; Use UNNEST(<expression>)",nil,nil)
		}
	| path_expression DOT_SYMBOL LR_BRACKET_SYMBOL {p.NotifyErrorListeners("Syntax error: Generalized field access is not allowed in the FROM clause without UNNEST; Use UNNEST(<expression>)",nil,nil)
		}
	| unnest_expression LS_BRACKET_SYMBOL {p.NotifyErrorListeners("Syntax error: Array element access is not allowed in the FROM clause without UNNEST; Use UNNEST(<expression>)",nil,nil)
		}
	| unnest_expression DOT_SYMBOL LR_BRACKET_SYMBOL {p.NotifyErrorListeners("Syntax error: Generalized field access is not allowed in the FROM clause without UNNEST; Use UNNEST(<expression>)",nil,nil)
		};

maybe_slashed_or_dashed_path_expression:
	maybe_dashed_path_expression
	| slashed_path_expression;

maybe_dashed_path_expression:
	path_expression
	| dashed_path_expression;

dashed_path_expression:
	dashed_identifier
	| dashed_path_expression DOT_SYMBOL identifier;

dashed_identifier:
	identifier MINUS_OPERATOR identifier
	| dashed_identifier MINUS_OPERATOR dashed_identifier
	| identifier MINUS_OPERATOR INTEGER_LITERAL
	| dashed_identifier MINUS_OPERATOR INTEGER_LITERAL
	| identifier MINUS_OPERATOR floating_point_literal identifier
	| dashed_identifier MINUS_OPERATOR floating_point_literal identifier;

slashed_identifier:
	SLASH_SYMBOL identifier_or_integer
	| slashed_identifier slashed_identifier_separator identifier_or_integer
	| slashed_identifier slashed_identifier_separator floating_point_literal
		slashed_identifier_separator identifier_or_integer;

identifier_or_integer:
	identifier
	| INTEGER_LITERAL; // TODO(zp): SCRIPT_LABEL;

slashed_identifier_separator:
	MINUS_OPERATOR SLASH_SYMBOL COLON_SYMBOL;

slashed_path_expression:
	slashed_identifier
	| slashed_identifier slashed_identifier_separator floating_point_literal identifier;

unnest_expression:
	unnest_expression_prefix opt_array_zip_mode? RR_BRACKET_SYMBOL
	| UNNEST_SYMBOL LR_BRACKET_SYMBOL SELECT_SYMBOL {p.NotifyErrorListeners("The argument to UNNEST is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression", nil, nil)
		};

unnest_expression_prefix:
	UNNEST_SYMBOL LR_BRACKET_SYMBOL expression_with_opt_alias (
		COMMA_SYMBOL expression_with_opt_alias
	)*;

opt_array_zip_mode: COMMA_SYMBOL named_argument;

expression_with_opt_alias:
	expression opt_as_alias_with_required_as?;

tvf_prefix:
	tvf_prefix_no_args tvf_argument (COMMA_SYMBOL tvf_argument)*;

tvf_argument:
	expression
	| descriptor_argument
	| table_clause
	| model_clause
	| connection_clause
	| named_argument
	| LR_BRACKET_SYMBOL table_clause RR_BRACKET_SYMBOL {p.NotifyErrorListeners("Syntax error: Table arguments for table-valued function calls written as \"TABLE path\" must not be enclosed in parentheses. To fix this, replace (TABLE path) with TABLE path",nil,nil)
		}
	| LR_BRACKET_SYMBOL model_clause RR_BRACKET_SYMBOL {p.NotifyErrorListeners("Syntax error: Model arguments for table-valued function calls written as \"MODEL path\" must not be enclosed in parentheses. To fix this, replace (MODEL path) with MODEL path",nil,nil)
		}
	| LR_BRACKET_SYMBOL connection_clause RR_BRACKET_SYMBOL {p.NotifyErrorListeners("Syntax error: Connection arguments for table-valued function calls written as \"CONNECTION path\" must not be enclosed in parentheses. To fix this, replace (CONNECTION path) with CONNECTION path",nil,nil)
		}
	| LR_BRACKET_SYMBOL named_argument RR_BRACKET_SYMBOL {p.NotifyErrorListeners("Syntax error: Named arguments for table-valued function calls written as \"name => value\" must not be enclosed in parentheses. To fix this, replace (name => value) with name => value",nil,nil)
		}
	| SELECT_SYMBOL {p.NotifyErrorListeners("Syntax error: Each subquery argument for table-valued function calls must be enclosed in parentheses. To fix this, replace SELECT... with (SELECT...)",nil,nil)
		}
	| WITH_SYMBOL {p.NotifyErrorListeners("Syntax error: Each subquery argument for table-valued function calls must be enclosed in parentheses. To fix this, replace WITH... with (WITH...)",nil,nil)
		};

connection_clause: CONNECTION_SYMBOL path_expression_or_default;

path_expression_or_default: path_expression | DEFAULT_SYMBOL;

descriptor_argument:
	DESCRIPTOR_SYMBOL LR_BRACKET_SYMBOL descriptor_column_list RR_BRACKET_SYMBOL;

descriptor_column_list:
	descriptor_column (COMMA_SYMBOL descriptor_column)*;

descriptor_column: identifier;

table_clause:
	TABLE_SYMBOL tvf_with_suffixes
	| TABLE_SYMBOL path_expression;

model_clause: MODEL_SYMBOL path_expression;

qualify_clause_nonreserved: QUALIFY_SYMBOL expression;

unpivot_clause:
	UNPIVOT_SYMBOL unpivot_nulls_filter? LR_BRACKET_SYMBOL path_expression_list_with_opt_parens
		FOR_SYMBOL path_expression IN_SYMBOL unpivot_in_item_list RR_BRACKET_SYMBOL;

unpivot_in_item_list:
	unpivot_in_item_list_prefix RR_BRACKET_SYMBOL;

unpivot_in_item_list_prefix:
	LR_BRACKET_SYMBOL unpivot_in_item
	| unpivot_in_item_list_prefix COMMA_SYMBOL unpivot_in_item;

unpivot_in_item:
	path_expression_list_with_opt_parens opt_as_string_or_integer?;

opt_as_string_or_integer:
	AS_SYMBOL? string_literal
	| AS_SYMBOL? integer_literal;

path_expression_list_with_opt_parens:
	LR_BRACKET_SYMBOL path_expression_list RR_BRACKET_SYMBOL
	| path_expression_list;

path_expression_list:
	path_expression (COMMA_SYMBOL path_expression)*;

unpivot_nulls_filter:
	EXCLUDE_SYMBOL NULLS_SYMBOL
	| INCLUDE_SYMBOL NULLS_SYMBOL;

pivot_clause:
	PIVOT_SYMBOL LR_BRACKET_SYMBOL pivot_expression_list FOR_SYMBOL expression_higher_prec_than_and
		IN_SYMBOL LR_BRACKET_SYMBOL pivot_value_list RR_BRACKET_SYMBOL RR_BRACKET_SYMBOL;

pivot_expression_list:
	pivot_expression (COMMA_SYMBOL pivot_expression)*;

pivot_expression: expression as_alias?;

pivot_value_list: pivot_value (COMMA_SYMBOL pivot_value)*;

pivot_value: expression as_alias?;

tvf_prefix_no_args:
	path_expression
	| IF_SYMBOL LR_BRACKET_SYMBOL;

join_type:
	CROSS_SYMBOL
	| FULL_SYMBOL opt_outer?
	| INNER_SYMBOL
	| LEFT_SYMBOL opt_outer?
	| RIGHT_SYMBOL opt_outer?;

opt_natural: NATURAL_SYMBOL;

on_clause:
	ON_SYMBOL expression /* Actullay, this should be bool_expression */;

select_list:
	select_list_item (COMMA_SYMBOL select_list_item)* COMMA_SYMBOL?;

select_list_item:
	select_column_expr
	| select_column_dot_star
	| select_column_star;

select_column_star: MULTIPLY_OPERATOR star_modifiers?;

select_column_expr:
	expression
	| select_column_expr_with_as_alias
	| expression identifier;

select_column_dot_star:
	expression_higher_prec_than_and DOT_SYMBOL MULTIPLY_OPERATOR star_modifiers?;

star_modifiers:
	star_except_list
	| star_except_list? star_replace_list;

star_except_list:
	EXCEPT_SYMBOL LR_BRACKET_SYMBOL identifier (
		COMMA_SYMBOL identifier
	)* RR_BRACKET_SYMBOL;

star_replace_list:
	REPLACE_SYMBOL LR_BRACKET_SYMBOL star_replace_item (
		COMMA_SYMBOL star_replace_item
	)* RR_BRACKET_SYMBOL;

star_replace_item: expression AS_SYMBOL identifier;

// expression: https://github.com/google/zetasql/blob/194cd32b5d766d60e3ca442651d792c7fe54ea74/zetasql/parser/bison_parser.y#L7712
expression:
	expression_higher_prec_than_and
	| and_expression
	| expression OR_SYMBOL expression;

// expression_higher_prec_than_and: https://github.com/google/zetasql/blob/194cd32b5d766d60e3ca442651d792c7fe54ea74/zetasql/parser/bison_parser.y#L7747
expression_higher_prec_than_and:
	// unparenthesized_expression_higher_prec_than_and scope begin
	null_literal
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
	| expression_higher_prec_than_and LS_BRACKET_SYMBOL expression RS_BRACKET_SYMBOL
	| expression_higher_prec_than_and DOT_SYMBOL LR_BRACKET_SYMBOL path_expression RR_BRACKET_SYMBOL
	| expression_higher_prec_than_and DOT_SYMBOL identifier
	| NOT_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and like_operator any_some_all hint? unnest_expression
	| expression_higher_prec_than_and like_operator any_some_all hint?
		parenthesized_anysomeall_list_in_rhs
	| expression_higher_prec_than_and like_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and distinct_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and in_operator hint? unnest_expression {
		if localctx.Hint() != nil {
			p.NotifyErrorListeners("Syntax error: HINTs cannot be specified on IN clause with UNNEST", nil, nil)
		}
	}
	| expression_higher_prec_than_and in_operator hint? parenthesized_in_rhs
	| expression_higher_prec_than_and between_operator expression_higher_prec_than_and AND_SYMBOL
		expression_higher_prec_than_and
	| expression_higher_prec_than_and between_operator expression_higher_prec_than_and OR_SYMBOL {
		p.NotifyErrorListeners("Syntax error: Expression in BETWEEN must be parenthesized", nil, nil)
	}
	| expression_higher_prec_than_and is_operator UNKNOWN_SYMBOL
	| expression_higher_prec_than_and is_operator null_literal
	| expression_higher_prec_than_and is_operator boolean_literal
	| expression_higher_prec_than_and comparative_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and STROKE_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and CIRCUMFLEX_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and BIT_AND_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and BOOL_OR_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and shift_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and additive_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and multiplicative_operator expression_higher_prec_than_and
	| unary_operator expression_higher_prec_than_and
	// unparenthesized_expression_higher_prec_than_and scope end
	| parenthesized_expression_not_a_query
	| parenthesized_query;

expression_maybe_parenthesized_not_a_query:
	parenthesized_expression_not_a_query
	// unparenthesized_expression_higher_prec_than_and scope begin
	| null_literal
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
	| expression_higher_prec_than_and LS_BRACKET_SYMBOL expression RS_BRACKET_SYMBOL
	| expression_higher_prec_than_and DOT_SYMBOL LR_BRACKET_SYMBOL path_expression RR_BRACKET_SYMBOL
	| expression_higher_prec_than_and DOT_SYMBOL identifier
	| NOT_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and like_operator any_some_all hint? unnest_expression
	| expression_higher_prec_than_and like_operator any_some_all hint?
		parenthesized_anysomeall_list_in_rhs
	| expression_higher_prec_than_and like_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and distinct_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and in_operator hint? unnest_expression {
		if localctx.Hint() != nil {
			p.NotifyErrorListeners("Syntax error: HINTs cannot be specified on IN clause with UNNEST", nil, nil)
		}
	}
	| expression_higher_prec_than_and in_operator hint? parenthesized_in_rhs
	| expression_higher_prec_than_and between_operator expression_higher_prec_than_and AND_SYMBOL
		expression_higher_prec_than_and
	| expression_higher_prec_than_and between_operator expression_higher_prec_than_and OR_SYMBOL {
		p.NotifyErrorListeners("Syntax error: Expression in BETWEEN must be parenthesized", nil, nil)
	}
	| expression_higher_prec_than_and is_operator UNKNOWN_SYMBOL
	| expression_higher_prec_than_and is_operator null_literal
	| expression_higher_prec_than_and is_operator boolean_literal
	| expression_higher_prec_than_and comparative_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and STROKE_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and CIRCUMFLEX_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and BIT_AND_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and BOOL_OR_SYMBOL expression_higher_prec_than_and
	| expression_higher_prec_than_and shift_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and additive_operator expression_higher_prec_than_and
	| expression_higher_prec_than_and multiplicative_operator expression_higher_prec_than_and
	| unary_operator expression_higher_prec_than_and
	// unparenthesized_expression_higher_prec_than_and scope end
	| and_expression
	// Previous or_expression, replace by solving mutually left-recursive.
	| expression OR_SYMBOL expression;

parenthesized_in_rhs:
	parenthesized_query
	| LR_BRACKET_SYMBOL expression_maybe_parenthesized_not_a_query RR_BRACKET_SYMBOL
	| in_list_two_or_more_prefix RR_BRACKET_SYMBOL;

unary_operator:
	PLUS_OPERATOR
	| MINUS_OPERATOR
	| BITWISE_NOT_OPERATOR;

comparative_operator:
	EQUAL_OPERATOR
	| NOT_EQUAL_OPERATOR
	| NOT_EQUAL2_OPERATOR
	| LT_OPERATOR
	| LE_OPERATOR
	| GT_OPERATOR
	| GE_OPERATOR;

shift_operator: KL_OPERATOR | KR_OPERATOR;

additive_operator: PLUS_OPERATOR | MINUS_OPERATOR;

multiplicative_operator: MULTIPLY_OPERATOR | DIVIDE_OPERATOR;

is_operator: IS_SYMBOL NOT_SYMBOL?;

between_operator: NOT_SYMBOL? BETWEEN_SYMBOL;

in_operator: NOT_SYMBOL? IN_SYMBOL;

distinct_operator:
	IS_SYMBOL NOT_SYMBOL? DISTINCT_SYMBOL FROM_SYMBOL;

parenthesized_query: LR_BRACKET_SYMBOL query RR_BRACKET_SYMBOL;

parenthesized_expression_not_a_query:
	LR_BRACKET_SYMBOL (
		expression_maybe_parenthesized_not_a_query
	) RR_BRACKET_SYMBOL;

parenthesized_anysomeall_list_in_rhs:
	parenthesized_query
	| parenthesized_expression_not_a_query
	| in_list_two_or_more_prefix RR_BRACKET_SYMBOL;

and_expression:
	expression_higher_prec_than_and AND_SYMBOL expression_higher_prec_than_and (
		AND_SYMBOL expression_higher_prec_than_and
	)*;

in_list_two_or_more_prefix:
	LR_BRACKET_SYMBOL expression COMMA_SYMBOL expression (
		COMMA_SYMBOL expression
	)*;

any_some_all: ANY_SYMBOL | SOME_SYMBOL | ALL_SYMBOL;

like_operator: LIKE_SYMBOL | NOT_SYMBOL LIKE_SYMBOL;

expression_subquery_with_keyword:
	ARRAY_SYMBOL parenthesized_query
	| EXISTS_SYMBOL hint? parenthesized_query;

struct_constructor:
	struct_constructor_prefix_with_keyword RR_BRACKET_SYMBOL
	| struct_constructor_prefix_with_keyword_no_arg RR_BRACKET_SYMBOL
	| struct_constructor_prefix_without_keyword RR_BRACKET_SYMBOL;

struct_constructor_prefix_with_keyword:
	struct_constructor_prefix_with_keyword_no_arg struct_constructor_arg (
		COMMA_SYMBOL struct_constructor_arg
	)*;

struct_constructor_arg:
	expression opt_as_alias_with_required_as?;

struct_constructor_prefix_without_keyword:
	LR_BRACKET_SYMBOL expression COMMA_SYMBOL expression (
		COMMA_SYMBOL expression
	)*;

struct_constructor_prefix_with_keyword_no_arg:
	struct_type LR_BRACKET_SYMBOL
	| STRUCT_SYMBOL LR_BRACKET_SYMBOL;

interval_expression:
	INTERVAL_SYMBOL expression identifier (TO_SYMBOL identifier)?;

function_call_expression_with_clauses:
	// NOTE: zetasql bison.y is LALR(1) parser, it checks the first rule should be path_expression
	// in action code instead of use expression directly to avoid parser ambiguous.
	path_expression LR_BRACKET_SYMBOL DISTINCT_SYMBOL? function_call_expression_with_clauses_suffix
	| function_name_from_keyword LR_BRACKET_SYMBOL function_call_expression_with_clauses_suffix;

function_call_expression_with_clauses_suffix:
	(
		// Empty argument list.
		opt_having_or_group_by_modifier? order_by_clause? limit_offset_clause? RR_BRACKET_SYMBOL
		// Non empty argument list.
		| (
			(function_call_argument | MULTIPLY_OPERATOR) (
				COMMA_SYMBOL function_call_argument
			)*
		) opt_null_handling_modifier? opt_having_or_group_by_modifier? clamped_between_modifier?
			with_report_modifier? order_by_clause? limit_offset_clause? RR_BRACKET_SYMBOL
	) hint? with_group_rows? over_clause?;

over_clause: OVER_SYMBOL window_specification;

window_specification:
	identifier
	| LR_BRACKET_SYMBOL identifier? partition_by_clause? order_by_clause? opt_window_frame_clause?
		RR_BRACKET_SYMBOL;

opt_window_frame_clause:
	frame_unit BETWEEN_SYMBOL window_frame_bound AND_SYMBOL window_frame_bound
	| frame_unit window_frame_bound;

window_frame_bound:
	UNBOUNDED_SYMBOL preceding_or_following
	| CURRENT_SYMBOL ROW_SYMBOL
	| expression preceding_or_following;

preceding_or_following: PRECEDING_SYMBOL | FOLLOWING_SYMBOL;

frame_unit: ROWS_SYMBOL | RANGE_SYMBOL;

partition_by_clause: partition_by_clause_prefix;

partition_by_clause_prefix:
	PARTITION_SYMBOL hint? BY_SYMBOL expression (
		COMMA_SYMBOL expression
	)*;

with_group_rows:
	WITH_SYMBOL GROUP_SYMBOL ROWS_SYMBOL /* XXX(zp): query = parenthesized_query*/;

with_report_modifier:
	WITH_SYMBOL REPORT_SYMBOL with_report_format;

clamped_between_modifier:
	CLAMPED_SYMBOL expression_higher_prec_than_and AND_SYMBOL expression;

with_report_format: options_list;

options_list:
	options_list_prefix RR_BRACKET_SYMBOL
	| LR_BRACKET_SYMBOL RR_BRACKET_SYMBOL;

options_list_prefix:
	LR_BRACKET_SYMBOL options_entry (COMMA_SYMBOL options_entry)*;

options_entry:
	identifier_in_hints options_assignment_operator expression_or_proto;

expression_or_proto: PROTO_SYMBOL | expression;

options_assignment_operator:
	EQUAL_OPERATOR
	| PLUS_EQUAL_SYMBOL
	| SUB_EQUAL_SYMBOL;

opt_null_handling_modifier:
	IGNORE_SYMBOL NULLS_SYMBOL
	| RESPECT_SYMBOL NULLS_SYMBOL;

function_call_argument:
	expression opt_as_alias_with_required_as?
	| named_argument
	| lambda_argument
	| sequence_arg
	| SELECT_SYMBOL { p.NotifyErrorListeners("Each function argument is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression", nil, nil); 
		};

sequence_arg: SEQUENCE_SYMBOL path_expression;

named_argument:
	identifier EQUAL_GT_BRACKET_SYMBOL expression
	| identifier EQUAL_GT_BRACKET_SYMBOL lambda_argument;

lambda_argument:
	lambda_argument_list SUB_GT_BRACKET_SYMBOL expression;

lambda_argument_list:
	/* XXX(zp): expr kind check expression*/ expression
	| LR_BRACKET_SYMBOL RR_BRACKET_SYMBOL;

limit_offset_clause:
	LIMIT_SYMBOL expression OFFSET_SYMBOL expression
	| LIMIT_SYMBOL expression;

opt_having_or_group_by_modifier:
	HAVING_SYMBOL MAX_SYMBOL expression
	| HAVING_SYMBOL MIN_SYMBOL expression group_by_clause_prefix;

group_by_clause_prefix:
	group_by_preamble grouping_item (COMMA_SYMBOL grouping_item)*;

group_by_preamble: GROUP_SYMBOL hint? opt_and_order? BY_SYMBOL;

opt_and_order: AND_SYMBOL ORDER_SYMBOL;

hint:
	/*XXX(zp): ABORT_CHECK*/ AT_SYMBOL integer_literal
	| hint_with_body;
hint_with_body: hint_with_body_prefix RC_BRACKET_SYMBOL;

hint_with_body_prefix:
	AT_SYMBOL (integer_literal AT_SYMBOL)? LC_BRACKET_SYMBOL hint_entry (
		COMMA_SYMBOL hint_entry
	)*;

hint_entry:
	identifier_in_hints EQUAL_OPERATOR expression
	| identifier_in_hints DOT_SYMBOL identifier_in_hints EQUAL_OPERATOR expression;

identifier_in_hints:
	identifier
	| extra_identifier_in_hints_name;

extra_identifier_in_hints_name:
	HASH_SYMBOL
	| PROTO_SYMBOL
	| PARTITION_SYMBOL;

grouping_item:
	LR_BRACKET_SYMBOL RR_BRACKET_SYMBOL
	| expression opt_as_alias_with_required_as? opt_grouping_item_order?
	| rollup_list RR_BRACKET_SYMBOL
	| cube_list RR_BRACKET_SYMBOL
	| grouping_set_list RR_BRACKET_SYMBOL;

grouping_set_list:
	GROUPING_SYMBOL SETS_SYMBOL LR_BRACKET_SYMBOL grouping_set (
		COMMA_SYMBOL grouping_set
	)*;

grouping_set:
	LR_BRACKET_SYMBOL RR_BRACKET_SYMBOL
	| expression
	| rollup_list RR_BRACKET_SYMBOL
	| cube_list RR_BRACKET_SYMBOL;

cube_list:
	CUBE_SYMBOL LR_BRACKET_SYMBOL (COMMA_SYMBOL expression)*;

rollup_list:
	ROLLUP_SYMBOL LR_BRACKET_SYMBOL expression (
		COMMA_SYMBOL expression
	)*;

opt_as_alias_with_required_as: AS_SYMBOL identifier;

opt_grouping_item_order: opt_selection_item_order | null_order;

opt_selection_item_order: asc_or_desc null_order?;

asc_or_desc: ASC_SYMBOL | DESC_SYMBOL;

null_order:
	NULLS_SYMBOL FIRST_SYMBOL
	| NULLS_SYMBOL LAST_SYMBOL;

function_name_from_keyword:
	IF_SYMBOL
	| GROUPING_SYMBOL
	| LEFT_SYMBOL
	| RIGHT_SYMBOL
	| COLLATE_SYMBOL
	| RANGE_SYMBOL;

replace_fields_expression:
	replace_fields_prefix RR_BRACKET_SYMBOL;

replace_fields_prefix:
	REPLACE_FIELDS_SYMBOL LR_BRACKET_SYMBOL expression COMMA_SYMBOL replace_fields_arg (
		COMMA_SYMBOL replace_fields_arg
	)*;

replace_fields_arg:
	expression AS_SYMBOL generalized_path_expression
	| expression AS_SYMBOL generalized_extension_path;

generalized_path_expression:
	identifier
	| generalized_path_expression DOT_SYMBOL generalized_extension_path
	| generalized_path_expression DOT_SYMBOL identifier
	| generalized_path_expression LS_BRACKET_SYMBOL expression RS_BRACKET_SYMBOL;

generalized_extension_path:
	LR_BRACKET_SYMBOL path_expression RR_BRACKET_SYMBOL
	| generalized_extension_path DOT_SYMBOL LR_BRACKET_SYMBOL path_expression RR_BRACKET_SYMBOL
	| generalized_extension_path DOT_SYMBOL identifier;

with_expression:
	/* XXX(zp): zetasql Yacc implement this in lookahead_transformer. */ WITH_SYMBOL
		LR_BRACKET_SYMBOL with_expression_variable_prefix COMMA_SYMBOL expression RR_BRACKET_SYMBOL;

with_expression_variable_prefix:
	with_expression_variable (
		COMMA_SYMBOL with_expression_variable
	)*;

with_expression_variable: identifier AS_SYMBOL expression;

extract_expression:
	extract_expression_base RR_BRACKET_SYMBOL
	| extract_expression_base AT_SYMBOL TIME_SYMBOL ZONE_SYMBOL expression RR_BRACKET_SYMBOL;

extract_expression_base:
	EXTRACT_SYMBOL LR_BRACKET_SYMBOL expression FROM_SYMBOL expression;

opt_format: FORMAT_SYMBOL expression opt_at_time_zone?;

opt_at_time_zone: AT_SYMBOL TIME_SYMBOL ZONE_SYMBOL expression;

cast_expression:
	CAST_SYMBOL LR_BRACKET_SYMBOL expression AS_SYMBOL type opt_format? RR_BRACKET_SYMBOL
	| CAST_SYMBOL LR_BRACKET_SYMBOL CAST_SYMBOL { p.NotifyErrorListeners("The argument to CAST is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression", nil, nil); 
		}
	| SAFE_CAST_SYMBOL LR_BRACKET_SYMBOL expression AS_SYMBOL type opt_format? RR_BRACKET_SYMBOL
	| SAFE_CAST_SYMBOL LR_BRACKET_SYMBOL SAFE_CAST_SYMBOL { p.NotifyErrorListeners("The argument to CAST is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression", nil, nil); 
		};

case_expression:
	case_expression_prefix END_SYMBOL
	| case_expression_prefix ELSE_SYMBOL expression END_SYMBOL;

case_expression_prefix:
	case_no_value_expression_prefix
	| case_value_expression_prefix;

case_value_expression_prefix:
	CASE_SYMBOL expression (
		WHEN_SYMBOL expression THEN_SYMBOL expression
	)+;

case_no_value_expression_prefix:
	CASE_SYMBOL (WHEN_SYMBOL expression THEN_SYMBOL expression)+;

struct_braced_constructor:
	stype = struct_type ctor = braced_constructor
	| STRUCT_SYMBOL ctor = braced_constructor;

braced_new_constructor: NEW_SYMBOL type_name new_constructor;

braced_constructor:
	braced_constructor_start RC_BRACKET_SYMBOL
	| braced_constructor_prefix RC_BRACKET_SYMBOL
	// Allow trailing comma in braced constructor. | braced_constructor_prefix
	COMMA_SYMBOL RC_BRACKET_SYMBOL;

braced_constructor_start: RC_BRACKET_SYMBOL;

braced_constructor_prefix:
	braced_constructor_start braced_constructor_field
	| braced_constructor_start braced_constructor_extension
	| braced_constructor_prefix COMMA_SYMBOL braced_constructor_field
	| braced_constructor_prefix braced_constructor_field
	| braced_constructor_prefix COMMA_SYMBOL braced_constructor_extension;

braced_constructor_field:
	braced_constructor_lhs braced_constructor_field_value;

braced_constructor_lhs: generalized_path_expression;

braced_constructor_field_value:
	COLON_SYMBOL expression
	| braced_constructor;

braced_constructor_extension:
	LR_BRACKET_SYMBOL path_expression RR_BRACKET_SYMBOL;

new_constructor:
	new_constructor_prefix RR_BRACKET_SYMBOL
	| new_constructor_prefix_no_arg RR_BRACKET_SYMBOL;

new_constructor_prefix:
	new_constructor_prefix_no_arg new_constructor_arg (
		COMMA_SYMBOL new_constructor_arg
	)*;

new_constructor_prefix_no_arg:
	NEW_SYMBOL type_name LR_BRACKET_SYMBOL;

new_constructor_arg:
	expression
	| expression AS_SYMBOL identifier
	| expression AS_SYMBOL LR_BRACKET_SYMBOL path_expression RR_BRACKET_SYMBOL;

array_constructor:
	array_constructor_prefix_no_expressions RS_BRACKET_SYMBOL
	| array_constructor_prefix RS_BRACKET_SYMBOL;

array_constructor_prefix:
	array_constructor_prefix_no_expressions expression (
		COMMA_SYMBOL expression
	)*;

array_constructor_prefix_no_expressions:
	ARRAY_SYMBOL LS_BRACKET_SYMBOL
	| LS_BRACKET_SYMBOL
	| array_type LS_BRACKET_SYMBOL;

range_literal: range_type string_literal;

range_type:
	RANGE_SYMBOL template_type_open type template_type_close;

type: raw_type opt_type_parameters? collate_clause?;

collate_clause: COLLATE_SYMBOL string_literal_or_parameter;

string_literal_or_parameter:
	string_literal
	| parameter_expression
	| system_variable_expression;

system_variable_expression: ATAT_SYMBOL path_expression;

parameter_expression:
	named_parameter_expression
	| QUESTION_SYMBOL;

named_parameter_expression: AT_SYMBOL identifier;

// This is opt_type_parameters in zetasql yacc, but here prefer to use ? in ANTLR.
opt_type_parameters:
	type_parameters_prefix RR_BRACKET_SYMBOL
	| type_parameters_prefix COMMA_SYMBOL RR_BRACKET_SYMBOL { p.NotifyErrorListeners("Syntax error: Trailing comma in type parameters list is not allowed.", nil, nil); 
		};

type_parameters_prefix:
	LR_BRACKET_SYMBOL type_parameter (
		COMMA_SYMBOL type_parameter
	)*;

type_parameter:
	integer_literal
	| boolean_literal
	| string_literal
	| bytes_literal
	| floating_point_literal
	| MAX_SYMBOL;

raw_type:
	array_type
	| struct_type
	| type_name
	| range_type
	| function_type
	| map_type;

map_type:
	MAP_SYMBOL template_type_open key_type = type COMMA_SYMBOL value_type = type template_type_close
		;

function_type:
	FUNCTION_SYMBOL template_type_open LR_BRACKET_SYMBOL RR_BRACKET_SYMBOL SUB_GT_BRACKET_SYMBOL
		return_type = type template_type_close
	| FUNCTION_SYMBOL template_type_open arg_type = type SUB_GT_BRACKET_SYMBOL return_type = type
		template_type_close
	| arg_list = function_type_prefix RR_BRACKET_SYMBOL SUB_GT_BRACKET_SYMBOL return_type = type
		template_type_close;

function_type_prefix:
	FUNCTION_SYMBOL template_type_open LR_BRACKET_SYMBOL type (
		COMMA_SYMBOL type
	)*;

type_name: path_expression | INTERVAL_SYMBOL;

path_expression: identifier (DOT_SYMBOL identifier)*;

identifier: token_identifier | keyword_as_identifier;

keyword_as_identifier:
	common_keyword_as_identifier
	| SIMPLE_SYMBOL;

common_keyword_as_identifier:
	ABORT_SYMBOL
	| ACCESS_SYMBOL
	| ACTION_SYMBOL
	| AGGREGATE_SYMBOL
	| ADD_SYMBOL
	| ALTER_SYMBOL
	| ALWAYS_SYMBOL
	| ANALYZE_SYMBOL
	| APPROX_SYMBOL
	| ARE_SYMBOL
	| ASSERT_SYMBOL
	| BATCH_SYMBOL
	| BEGIN_SYMBOL
	| BIGDECIMAL_SYMBOL
	| BIGNUMERIC_SYMBOL
	| BREAK_SYMBOL
	| CALL_SYMBOL
	| CASCADE_SYMBOL
	| CHECK_SYMBOL
	| CLAMPED_SYMBOL
	| CONFLICT_SYMBOL
	| CLONE_SYMBOL
	| COPY_SYMBOL
	| CLUSTER_SYMBOL
	| COLUMN_SYMBOL
	| COLUMNS_SYMBOL
	| COMMIT_SYMBOL
	| CONNECTION_SYMBOL
	| CONSTANT_SYMBOL
	| CONSTRAINT_SYMBOL
	| CONTINUE_SYMBOL
	| CORRESPONDING_SYMBOL
	| CYCLE_SYMBOL
	| DATA_SYMBOL
	| DATABASE_SYMBOL
	| DATE_SYMBOL
	| DATETIME_SYMBOL
	| DECIMAL_SYMBOL
	| DECLARE_SYMBOL
	| DEFINER_SYMBOL
	| DELETE_SYMBOL
	| DELETION_SYMBOL
	| DEPTH_SYMBOL
	| DESCRIBE_SYMBOL
	| DETERMINISTIC_SYMBOL
	| DO_SYMBOL
	| DROP_SYMBOL
	| ELSEIF_SYMBOL
	| ENFORCED_SYMBOL
	| ERROR_SYMBOL
	| EXCEPTION_SYMBOL
	| EXECUTE_SYMBOL
	| EXPLAIN_SYMBOL
	| EXPORT_SYMBOL
	| EXTEND_SYMBOL
	| EXTERNAL_SYMBOL
	| FILES_SYMBOL
	| FILTER_SYMBOL
	| FILL_SYMBOL
	| FIRST_SYMBOL
	| FOREIGN_SYMBOL
	| FORMAT_SYMBOL
	| FUNCTION_SYMBOL
	| GENERATED_SYMBOL
	| GRANT_SYMBOL
	| GROUP_ROWS_SYMBOL
	| HIDDEN_SYMBOL
	| IDENTITY_SYMBOL
	| IMMEDIATE_SYMBOL
	| IMMUTABLE_SYMBOL
	| IMPORT_SYMBOL
	| INCLUDE_SYMBOL
	| INCREMENT_SYMBOL
	| INDEX_SYMBOL
	| INOUT_SYMBOL
	| INPUT_SYMBOL
	| INSERT_SYMBOL
	| INVOKER_SYMBOL
	| ISOLATION_SYMBOL
	| ITERATE_SYMBOL
	| JSON_SYMBOL
	| KEY_SYMBOL
	| LANGUAGE_SYMBOL
	| LAST_SYMBOL
	| LEAVE_SYMBOL
	| LEVEL_SYMBOL
	| LOAD_SYMBOL
	| LOOP_SYMBOL
	| MACRO_SYMBOL
	| MAP_SYMBOL
	| MATCH_SYMBOL
	| KW_MATCH_RECOGNIZE_NONRESERVED_SYMBOL
	| MATCHED_SYMBOL
	| MATERIALIZED_SYMBOL
	| MAX_SYMBOL
	| MAXVALUE_SYMBOL
	| MEASURES_SYMBOL
	| MESSAGE_SYMBOL
	| METADATA_SYMBOL
	| MIN_SYMBOL
	| MINVALUE_SYMBOL
	| MODEL_SYMBOL
	| MODULE_SYMBOL
	| NUMERIC_SYMBOL
	| OFFSET_SYMBOL
	| ONLY_SYMBOL
	| OPTIONS_SYMBOL
	| OUT_SYMBOL
	| OUTPUT_SYMBOL
	| OVERWRITE_SYMBOL
	| PARTITIONS_SYMBOL
	| PATTERN_SYMBOL
	| PERCENT_SYMBOL
	| PIVOT_SYMBOL
	| POLICIES_SYMBOL
	| POLICY_SYMBOL
	| PRIMARY_SYMBOL
	| PRIVATE_SYMBOL
	| PRIVILEGE_SYMBOL
	| PRIVILEGES_SYMBOL
	| PROCEDURE_SYMBOL
	| PROJECT_SYMBOL
	| PUBLIC_SYMBOL
	| RAISE_SYMBOL
	| READ_SYMBOL
	| REFERENCES_SYMBOL
	| REMOTE_SYMBOL
	| REMOVE_SYMBOL
	| RENAME_SYMBOL
	| REPEAT_SYMBOL
	| REPEATABLE_SYMBOL
	| REPLACE_SYMBOL
	| REPLACE_FIELDS_SYMBOL
	| REPLICA_SYMBOL
	| REPORT_SYMBOL
	| RESTRICT_SYMBOL
	| RESTRICTION_SYMBOL
	| RETURNS_SYMBOL
	| RETURN_SYMBOL
	| REVOKE_SYMBOL
	| ROLLBACK_SYMBOL
	| ROW_SYMBOL
	| RUN_SYMBOL
	| SAFE_CAST_SYMBOL
	| SCHEMA_SYMBOL
	| SEARCH_SYMBOL
	| SECURITY_SYMBOL
	| SEQUENCE_SYMBOL
	| SETS_SYMBOL
	| SHOW_SYMBOL
	| SNAPSHOT_SYMBOL
	| SOURCE_SYMBOL
	| SQL_SYMBOL
	| STABLE_SYMBOL
	| START_SYMBOL
	| STATIC_DESCRIBE_SYMBOL
	| STORED_SYMBOL
	| STORING_SYMBOL
	| STRICT_SYMBOL
	| SYSTEM_SYMBOL
	| SYSTEM_TIME_SYMBOL
	| TABLE_SYMBOL
	| TABLES_SYMBOL
	| TARGET_SYMBOL
	| TEMP_SYMBOL
	| TEMPORARY_SYMBOL
	| TIME_SYMBOL
	| TIMESTAMP_SYMBOL
	| TRANSACTION_SYMBOL
	| TRANSFORM_SYMBOL
	| TRUNCATE_SYMBOL
	| TYPE_SYMBOL
	| UNDROP_SYMBOL
	| UNIQUE_SYMBOL
	| UNKNOWN_SYMBOL
	| UNPIVOT_SYMBOL
	| UNTIL_SYMBOL
	| UPDATE_SYMBOL
	| VALUE_SYMBOL
	| VALUES_SYMBOL
	| VECTOR_SYMBOL
	| VIEW_SYMBOL
	| VIEWS_SYMBOL
	| VOLATILE_SYMBOL
	| WEIGHT_SYMBOL
	| WHILE_SYMBOL
	| WRITE_SYMBOL
	| ZONE_SYMBOL
	| DESCRIPTOR_SYMBOL
	| INTERLEAVE_SYMBOL
	| NULL_FILTERED_SYMBOL
	| PARENT_SYMBOL
	| DESTINATION_SYMBOL
	| PROPERTY_SYMBOL
	| GRAPH_SYMBOL
	| NODE_SYMBOL
	| PROPERTIES_SYMBOL
	| LABEL_SYMBOL
	| EDGE_SYMBOL
	| NEXT_SYMBOL
	| ASCENDING_SYMBOL
	| DESCENDING_SYMBOL
	| SKIP_SYMBOL
	| PATH_SYMBOL
	| PATHS_SYMBOL
	| WALK_SYMBOL
	| TRAIL_SYMBOL
	| ACYCLIC_SYMBOL
	| OPTIONAL_SYMBOL
	| LET_SYMBOL;

token_identifier: IDENTIFIER;

struct_type:
	STRUCT_SYMBOL template_type_open template_type_close
	| struct_type_prefix template_type_close;

struct_type_prefix:
	STRUCT_SYMBOL template_type_open struct_field (
		COMMA_SYMBOL struct_field
	)*;

struct_field: identifier type | type;

array_type:
	ARRAY_SYMBOL template_type_open type template_type_close;

template_type_open: LT_OPERATOR;

template_type_close: GT_OPERATOR;

date_or_time_literal: date_or_time_literal_kind string_literal;

date_or_time_literal_kind:
	DATE_SYMBOL
	| TIME_SYMBOL
	| DATETIME_SYMBOL
	| TIMESTAMP_SYMBOL;

floating_point_literal: FLOATING_POINT_LITERAL;

json_literal: JSON_SYMBOL string_literal;

bignumeric_literal: bignumeric_literal_prefix string_literal;

bignumeric_literal_prefix:
	BIGNUMERIC_SYMBOL
	| BIGDECIMAL_SYMBOL;

numeric_literal: numeric_literal_prefix string_literal;

numeric_literal_prefix: NUMERIC_SYMBOL | DECIMAL_SYMBOL;

integer_literal: INTEGER_LITERAL;

bytes_literal:
	bytes_literal_component
	| bytes_literal bytes_literal_component {
	 literalStopIndex, componentStartIndex := localctx.Bytes_literal().GetStop().GetStop(), localctx.Bytes_literal_component().GetStart().GetStart() 
	 if literalStopIndex + 1 == componentStartIndex { p.NotifyErrorListeners("Syntax error: concatenated bytes literals must be separated by whitespace or comments.", nil, nil) } 
		}
	| bytes_literal string_literal_component {p.NotifyErrorListeners("Syntax error: string and bytes literals cannot be concatenated.", nil,
	 nil); };

null_literal: NULL_SYMBOL;

boolean_literal: TRUE_SYMBOL | FALSE_SYMBOL;

string_literal:
	string_literal_component
	| string_literal string_literal_component {
	 literalStopIndex, componentStartIndex := localctx.String_literal().GetStop().GetStop(), localctx.String_literal_component().GetStart().GetStart() 
	 if literalStopIndex + 1 == componentStartIndex { p.NotifyErrorListeners("Syntax error: concatenated string literals must be separated by whitespace or comments.", nil, nil) } 
		}
	| string_literal bytes_literal_component {p.NotifyErrorListeners("Syntax error: string and bytes literals cannot be concatenated.", nil, nil); 
		};

string_literal_component: STRING_LITERAL;

bytes_literal_component: BYTES_LITERAL;