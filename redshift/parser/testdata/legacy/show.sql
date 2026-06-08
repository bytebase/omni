-- SHOW command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW.html

-- Basic SHOW commands for specific parameters
SHOW query_group;
SHOW datestyle;
SHOW extra_float_digits;
SHOW search_path;
SHOW statement_timeout;
SHOW enable_case_sensitive_identifier;
SHOW enable_result_cache_for_session;
SHOW json_configuration_file;
SHOW max_concurrency_scaling_clusters;
SHOW mv_enable_aqmv_for_session;
SHOW use_fips_ssl;
SHOW wlm_json_configuration;

-- SHOW ALL command
SHOW ALL;

-- Session context variables
SHOW app_context.user_id;
SHOW app_context.client_id;
SHOW app_context.tenant_id;
SHOW application_context.user_name;
SHOW session_context.request_id;

-- Common parameter names from AWS documentation
SHOW analyze_threshold_percent;
SHOW datashare_database_mode;
SHOW default_geometry_encoding;
SHOW describe_field_name_in_uppercase;
SHOW downcase_delimited_identifier;
SHOW enable_hll;
SHOW enable_vacuum_boost;
SHOW explain_returntargetlist;
SHOW ignore_session_isolation_level_for_datashares;
SHOW json_serialization_enable;
SHOW json_serialization_parse_nested_strings;
SHOW max_cursor_result_set_size;
SHOW preserved_statements_count;
SHOW query_planning_improvement;
SHOW reader_endpoint_validated_clusters;
SHOW spectrum_enable_pseudo_columns;
SHOW timezone;

-- Additional common parameters
SHOW client_encoding;
SHOW IntervalStyle;
SHOW standard_conforming_strings;
SHOW transaction_isolation;
SHOW transaction_read_only;
SHOW transaction_deferrable;

-- Mixed case parameter names
SHOW TimeZone;
SHOW DateStyle;
SHOW IntervalStyle;

-- Parameter names with underscores
SHOW enable_interleaved_sort_key;
SHOW max_query_queue_time;
SHOW query_timeout;
SHOW idle_in_transaction_session_timeout;
SHOW lock_timeout;

-- System information parameters
SHOW server_version;
SHOW server_encoding;
SHOW application_name;

-- WLM related parameters
SHOW wlm_query_slot_count;
SHOW wlm_query_group;
SHOW wlm_memory_percent_to_use;

-- Spectrum related parameters
SHOW spectrum_query_maxerror;
SHOW spectrum_iam_role;

-- Encoding related parameters
SHOW client_encoding;
SHOW server_encoding;

-- Session parameters (using valid Redshift SHOW syntax)
SHOW current_user;
SHOW current_database;

-- Database connection parameters
SHOW current_database;
SHOW current_user;
SHOW session_user;

-- Additional context variable examples with different naming patterns
SHOW app.user_id;
SHOW application.session_id;
SHOW context.request_id;
SHOW custom_context.department_id;
SHOW my_app_context.organization_id;