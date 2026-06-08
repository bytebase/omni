-- Test cases for SET command
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SET.html

-- Basic SET with TO syntax
SET datestyle TO 'SQL,DMY';
SET search_path TO myschema, public;
SET statement_timeout TO 30000;
SET timezone TO 'America/New_York';

-- Basic SET with = syntax
SET datestyle = 'ISO,MDY';
SET search_path = '$user', public;
SET statement_timeout = 60000;
SET timezone = 'UTC';

-- SET with SESSION scope (explicit)
SET SESSION datestyle TO 'German,DMY';
SET SESSION search_path TO myschema;
SET SESSION statement_timeout TO 0;
SET SESSION enable_result_cache_for_session TO off;

-- SET with LOCAL scope (transaction-specific)
SET LOCAL datestyle TO 'Postgres,MDY';
SET LOCAL search_path TO tempschema;
SET LOCAL statement_timeout TO 10000;
SET LOCAL timezone TO 'PST8PDT';

-- SET SEED for random number generation
SET SEED TO .25;
SET SEED TO 0.5;
SET SEED TO 0.99;
SET SESSION SEED TO 0.1;
SET LOCAL SEED TO 0.75;

-- SET context variables (custom application context)
SET app_context.user_id TO 123;
SET app_context.user_id TO 'sample_variable_value';
SET app_context.tenant_id TO 'company_123';
SET custom_context.department TO 'engineering';
SET session_context.application_name TO 'myapp';

-- SET with quoted values
SET datestyle TO 'SQL,DMY';
SET search_path TO 'myschema', 'public';
SET query_group TO 'priority';
SET app_context.user_name TO 'John Doe';

-- SET identity namespace
SET SESSION default_identity_namespace = 'MYCO';
SET default_identity_namespace TO 'COMPANY_ABC';
SET LOCAL default_identity_namespace = 'TEST_NAMESPACE';

-- SET with DEFAULT keyword
SET datestyle TO DEFAULT;
SET search_path TO DEFAULT;
SET timezone TO DEFAULT;
SET SESSION statement_timeout TO DEFAULT;

-- SET workload management parameters
SET query_group TO 'high_priority';
SET query_group = 'reporting';
SET wlm_json_configuration TO '{"queue_name":"etl"}';

-- SET result cache parameters
SET enable_result_cache_for_session TO on;
SET enable_result_cache_for_session = off;
SET enable_result_cache_for_session TO true;
SET enable_result_cache_for_session = false;

-- SET concurrency scaling
SET max_concurrency_scaling_clusters TO 10;
SET max_concurrency_scaling_clusters = 5;

-- SET spectrum parameters
SET spectrum_enable_pseudo_columns TO true;
SET spectrum_enable_pseudo_columns = false;

-- SET SESSION CHARACTERISTICS (transaction isolation)
SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL READ COMMITTED;
SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL SERIALIZABLE;

-- Case variations
set datestyle to 'SQL,DMY';
Set Search_Path To myschema;
SET DATESTYLE TO 'ISO,MDY';

-- Multiple SET statements in sequence (common usage pattern)
SET query_group TO 'etl_jobs';
SET statement_timeout TO 3600000;
SET enable_result_cache_for_session TO off;

-- SET before operations (typical usage)
SET query_group TO 'priority';
-- SELECT tbl, count(*) FROM stv_blocklist;
-- RESET query_group;

SET app_context.user_id TO 456;
-- SELECT * FROM user_data WHERE user_id = current_setting('app_context.user_id');
-- RESET app_context.user_id;

-- Boolean values
SET enable_result_cache_for_session TO TRUE;
SET enable_result_cache_for_session TO FALSE;
SET spectrum_enable_pseudo_columns TO true;
SET spectrum_enable_pseudo_columns TO false;

-- Numeric values
SET statement_timeout TO 0;
SET statement_timeout TO 1000;
SET statement_timeout TO 3600000;
SET max_concurrency_scaling_clusters TO 1;

-- Additional context variable examples (simplified)
SET app_context.config TO '30';
SET session_context.data TO '2024-01-01';
SET custom_context.setting TO '1000';

-- Mixed case identifiers
SET "MixedCase_Param" TO 'value';
SET app_context."UserID" TO '789';
SET "CUSTOM_CONTEXT"."DEPT_CODE" TO 'FIN';

-- Error cases that should be tested (parser should handle gracefully)
-- SET;  -- Missing parameter and value
-- SET datestyle;  -- Missing TO/= and value
-- SET datestyle TO;  -- Missing value
-- SET TO 'value';  -- Missing parameter
-- SET LOCAL;  -- Missing parameter and value
-- SET SESSION LOCAL datestyle TO 'SQL,DMY';  -- Both SESSION and LOCAL