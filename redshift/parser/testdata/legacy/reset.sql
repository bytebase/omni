-- Test cases for RESET command
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_RESET.html

-- Basic RESET with specific parameter
RESET query_group;
RESET search_path;
RESET datestyle;
RESET timezone;
RESET statement_timeout;

-- RESET ALL - resets all runtime parameters and session context variables
RESET ALL;

-- Reset context variables (session context variables)
RESET app_context.user_id;
RESET app_context.tenant_id;
RESET session_context.application_name;
RESET custom_context.department;

-- Common configuration parameters that can be reset
RESET enable_result_cache_for_session;
RESET wlm_json_configuration;
RESET max_concurrency_scaling_clusters;
RESET spectrum_enable_pseudo_columns;

-- Case variations
reset query_group;
Reset ALL;
RESET all;

-- Multiple RESET statements in sequence
RESET query_group;
RESET search_path;
RESET ALL;

-- RESET after SET operations (common usage pattern)
SET query_group TO 'high_priority';
-- ... execute queries ...
RESET query_group;

SET search_path TO myschema, public;
-- ... execute queries ...
RESET search_path;

-- Session variables reset
SET app_context.user_id TO 12345;
-- ... execute queries ...
RESET app_context.user_id;

-- Error cases that should be tested (parser should handle gracefully)
-- RESET;  -- Missing parameter
-- RESET query_group search_path;  -- Multiple parameters not allowed
-- RESET (query_group);  -- Invalid syntax with parentheses