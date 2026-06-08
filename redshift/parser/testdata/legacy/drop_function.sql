-- DROP FUNCTION test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_FUNCTION.html

-- Basic DROP FUNCTION examples

-- Drop function with integer parameter (from AWS docs example)
DROP FUNCTION f_sqrt(int);

-- Drop function with CASCADE option (from AWS docs example)
DROP FUNCTION f_sqrt(int) CASCADE;

-- Drop function with RESTRICT option (explicit default behavior)
DROP FUNCTION f_sqrt(int) RESTRICT;

-- DROP FUNCTION with IF EXISTS

-- Drop function if it exists (prevents errors)
DROP FUNCTION IF EXISTS f_sqrt(int);

-- Drop function if exists with CASCADE
DROP FUNCTION IF EXISTS f_sqrt(int) CASCADE;

-- Drop function if exists with RESTRICT
DROP FUNCTION IF EXISTS f_sqrt(int) RESTRICT;

-- Schema-qualified function names

-- Drop function in specific schema
DROP FUNCTION myschema.f_calculate(integer, float);

-- Drop function in quoted schema
DROP FUNCTION "my_schema".f_process(varchar);

-- Drop function with quoted function name
DROP FUNCTION "f_special_calc"(decimal);

-- Drop function with both schema and function name quoted
DROP FUNCTION "my_schema"."f_special_function"(text);

-- Functions with various data types (matching CREATE FUNCTION examples)

-- Python UDF examples to drop
DROP FUNCTION f_py_greater(float, float);
DROP FUNCTION f_py_calculator(integer, integer, varchar);
DROP FUNCTION f_py_square(float);
DROP FUNCTION f_py_is_even(integer);
DROP FUNCTION f_py_clean_string(varchar);
DROP FUNCTION f_py_days_between(date, date);

-- SQL UDF examples to drop
DROP FUNCTION f_sql_greater(float, float);
DROP FUNCTION f_sql_discount_price(decimal, decimal);
DROP FUNCTION f_sql_circle_area(float);
DROP FUNCTION f_sql_full_name(varchar, varchar);
DROP FUNCTION f_sql_grade_letter(integer);

-- Functions with specific data type variations

-- Drop function with SMALLINT
DROP FUNCTION f_process_small(smallint);

-- Drop function with BIGINT
DROP FUNCTION f_process_big(bigint);

-- Drop function with DECIMAL with precision
DROP FUNCTION f_compound_interest(decimal(10,2), decimal(5,4), integer);

-- Drop function with NUMERIC with precision and scale
DROP FUNCTION f_financial_calc(numeric(15,2));

-- Drop function with REAL
DROP FUNCTION f_real_calc(real);

-- Drop function with DOUBLE PRECISION
DROP FUNCTION f_precise_calc(double precision);

-- Drop function with FLOAT with precision
DROP FUNCTION f_float_calc(float(24));

-- Drop function with BOOLEAN
DROP FUNCTION f_bool_check(boolean);

-- Drop function with BOOL (alias)
DROP FUNCTION f_bool_alias(bool);

-- Drop function with CHAR with length
DROP FUNCTION f_initial(char(50));

-- Drop function with CHARACTER with length
DROP FUNCTION f_character_proc(character(10));

-- Drop function with VARCHAR with length
DROP FUNCTION f_format_phone(varchar(15));

-- Drop function with CHARACTER VARYING
DROP FUNCTION f_varying_text(character varying(100));

-- Drop function with DATE
DROP FUNCTION f_date_calc(date);

-- Drop function with TIMESTAMP
DROP FUNCTION f_timestamp_proc(timestamp);

-- Drop function with TIMESTAMP with precision
DROP FUNCTION f_timestamp_precise(timestamp(6));

-- Drop function with TIMESTAMP WITH TIME ZONE
DROP FUNCTION f_timestamp_tz(timestamp with time zone);

-- Drop function with TIMESTAMP WITHOUT TIME ZONE
DROP FUNCTION f_timestamp_no_tz(timestamp without time zone);

-- Drop function with TIMESTAMPTZ
DROP FUNCTION f_timestamptz_proc(timestamptz);

-- Drop function with TIMESTAMPTZ with precision
DROP FUNCTION f_timestamptz_precise(timestamptz(3));

-- Drop function with TIME
DROP FUNCTION f_time_proc(time);

-- Drop function with TIME with precision
DROP FUNCTION f_time_precise(time(6));

-- Drop function with TIME WITH TIME ZONE
DROP FUNCTION f_time_tz(time with time zone);

-- Drop function with TIME WITHOUT TIME ZONE
DROP FUNCTION f_time_no_tz(time without time zone);

-- Drop function with TIMETZ
DROP FUNCTION f_timetz_proc(timetz);

-- Drop function with TIMETZ with precision
DROP FUNCTION f_timetz_precise(timetz(6));

-- Drop function with ANYELEMENT
DROP FUNCTION f_first_non_null(anyelement, anyelement);

-- Drop function with TEXT
DROP FUNCTION f_text_process(text);

-- Redshift-specific data types

-- Drop function with SUPER
DROP FUNCTION f_super_length(super);

-- Drop function with GEOMETRY
DROP FUNCTION f_geometry_area(geometry);

-- Drop function with GEOGRAPHY
DROP FUNCTION f_geography_distance(geography, geography);

-- Drop function with HLLSKETCH
DROP FUNCTION f_hll_cardinality(hllsketch);

-- Functions with multiple parameters

-- Drop function with two parameters
DROP FUNCTION f_two_params(integer, varchar);

-- Drop function with three parameters
DROP FUNCTION f_three_params(decimal, float, boolean);

-- Drop function with many parameters (demonstrating up to 32 limit)
DROP FUNCTION f_sum_many(
    integer, integer, integer, integer, integer,
    integer, integer, integer, integer, integer
);

-- Drop function with mixed data types
DROP FUNCTION f_mixed_types(
    smallint, bigint, decimal(10,2), varchar(50),
    boolean, date, timestamp, geometry
);

-- Functions with no parameters

-- Drop function with empty parameter list
DROP FUNCTION f_current_year();

-- Drop function with no parameters using IF EXISTS
DROP FUNCTION IF EXISTS f_no_params();

-- Functions with argument names (names are ignored during drop)

-- Drop function where original had argument names (Python UDF style)
DROP FUNCTION f_py_with_names(integer, varchar);

-- Drop function specifying argument names (optional, ignored)
DROP FUNCTION f_with_arg_names(input_val integer, input_text varchar);

-- Mixed cases with argument names and without
DROP FUNCTION f_mixed_args(integer, description varchar, active boolean);

-- Edge cases and special scenarios

-- Drop function with CASCADE and IF EXISTS combined
DROP FUNCTION IF EXISTS f_cascade_example(integer) CASCADE;

-- Drop function with RESTRICT and IF EXISTS combined
DROP FUNCTION IF EXISTS f_restrict_example(varchar) RESTRICT;

-- Drop function with very long name
DROP FUNCTION f_very_long_function_name_that_demonstrates_identifier_length(integer);

-- Drop function with quoted identifiers containing special characters
DROP FUNCTION "f_special-function"(integer);
DROP FUNCTION "F_UPPER_CASE_FUNCTION"(varchar);

-- Drop function in system-like schema names
DROP FUNCTION pg_catalog.f_system_function(integer);
DROP FUNCTION information_schema.f_info_function(text);

-- Drop functions that might have been created with OR REPLACE
DROP FUNCTION f_replaceable_function(float, float);

-- Complex data type combinations

-- Drop function with nested precision specifications
DROP FUNCTION f_complex_numeric(
    decimal(15,4),
    numeric(20,6),
    float(53),
    timestamp(6) with time zone,
    varchar(255)
);

-- Functions demonstrating all three drop behaviors

-- Default behavior (RESTRICT implied)
DROP FUNCTION f_default_behavior(integer);

-- Explicit CASCADE (drops dependent objects)
DROP FUNCTION f_with_cascade(varchar) CASCADE;

-- Explicit RESTRICT (prevents drop if dependencies exist)
DROP FUNCTION f_with_restrict(boolean) RESTRICT;

-- Real-world scenario examples

-- Drop mathematical functions
DROP FUNCTION f_calculate_distance(float, float, float, float);
DROP FUNCTION f_compound_interest_calculator(decimal(12,2), decimal(5,4), integer);

-- Drop string manipulation functions
DROP FUNCTION f_clean_phone_number(varchar);
DROP FUNCTION f_extract_domain(varchar);

-- Drop date/time functions
DROP FUNCTION f_business_days_between(date, date);
DROP FUNCTION f_format_timestamp(timestamp);

-- Drop validation functions
DROP FUNCTION f_is_valid_email(varchar);
DROP FUNCTION f_check_data_quality(super);

-- Drop aggregation helper functions
DROP FUNCTION f_weighted_average(decimal, decimal);
DROP FUNCTION f_percentile_rank(integer, integer);

-- Functions with custom type names (identifiers)
DROP FUNCTION f_custom_type_function(my_custom_type);
DROP FUNCTION f_enum_function(status_enum, priority_level);