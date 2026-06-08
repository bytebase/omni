-- DROP PROCEDURE test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_PROCEDURE.html

-- Basic DROP PROCEDURE examples

-- Drop procedure with mixed parameters (from AWS docs example)
DROP PROCEDURE quarterly_revenue(volume INOUT bigint, at_price IN numeric, result OUT int);

-- Drop procedure with simple parameter list
DROP PROCEDURE calculate_total(price IN decimal, tax_rate IN decimal);

-- Drop procedure with no parameters
DROP PROCEDURE update_statistics();

-- DROP PROCEDURE with IF EXISTS

-- Drop procedure if it exists (prevents errors)
DROP PROCEDURE IF EXISTS quarterly_revenue(volume INOUT bigint, at_price IN numeric, result OUT int);

-- Drop procedure if exists with single parameter
DROP PROCEDURE IF EXISTS calculate_discount(amount IN decimal);

-- Drop procedure if exists with no parameters
DROP PROCEDURE IF EXISTS refresh_cache();

-- DROP PROCEDURE with CASCADE/RESTRICT options

-- Drop procedure with CASCADE (removes dependent objects)
DROP PROCEDURE calculate_total(price IN decimal, tax_rate IN decimal) CASCADE;

-- Drop procedure with RESTRICT (prevents drop if dependencies exist)
DROP PROCEDURE calculate_total(price IN decimal, tax_rate IN decimal) RESTRICT;

-- Drop procedure with IF EXISTS and CASCADE
DROP PROCEDURE IF EXISTS quarterly_revenue(volume INOUT bigint, at_price IN numeric, result OUT int) CASCADE;

-- Drop procedure with IF EXISTS and RESTRICT
DROP PROCEDURE IF EXISTS monthly_report(month_id IN integer) RESTRICT;

-- Schema-qualified procedure names

-- Drop procedure in specific schema
DROP PROCEDURE finance.calculate_revenue(start_date IN date, end_date IN date);

-- Drop procedure in quoted schema
DROP PROCEDURE "reporting".monthly_summary(year_val IN integer);

-- Drop procedure with quoted procedure name
DROP PROCEDURE "complex-procedure-name"(input_data IN text);

-- Drop procedure with both schema and procedure name quoted
DROP PROCEDURE "finance_schema"."quarterly_analysis"(quarter IN integer, year_val IN integer);

-- Parameter modes (IN, OUT, INOUT)

-- Drop procedure with IN parameters only
DROP PROCEDURE process_data(input_val IN integer, description IN varchar);

-- Drop procedure with OUT parameters only
DROP PROCEDURE get_results(total_count OUT integer, avg_value OUT decimal);

-- Drop procedure with INOUT parameters only
DROP PROCEDURE transform_data(data INOUT text);

-- Drop procedure with mixed parameter modes
DROP PROCEDURE complex_calc(
    input_amount IN decimal,
    multiplier IN float,
    result OUT decimal,
    status INOUT varchar
);

-- Different parameter mode orderings

-- Mode before argument name
DROP PROCEDURE process_order(IN order_id integer, OUT total_amount decimal);

-- Argument name before mode
DROP PROCEDURE calculate_tax(amount IN decimal, rate IN float, tax_amount OUT decimal);

-- Mixed ordering within same procedure
DROP PROCEDURE mixed_params(IN input_val integer, result OUT decimal, status INOUT text);

-- Data type variations

-- Drop procedure with SMALLINT
DROP PROCEDURE process_small_numbers(count IN smallint);

-- Drop procedure with INTEGER and INT (aliases)
DROP PROCEDURE process_integers(val1 IN integer, val2 IN int);

-- Drop procedure with BIGINT
DROP PROCEDURE process_large_numbers(big_num IN bigint);

-- Drop procedure with DECIMAL with precision
DROP PROCEDURE financial_calc(amount IN decimal(12,2), rate IN decimal(5,4));

-- Drop procedure with NUMERIC with precision and scale
DROP PROCEDURE precise_calculation(value IN numeric(15,6));

-- Drop procedure with REAL
DROP PROCEDURE real_number_proc(val IN real);

-- Drop procedure with DOUBLE PRECISION
DROP PROCEDURE high_precision_calc(precise_val IN double precision);

-- Drop procedure with FLOAT with precision
DROP PROCEDURE float_calculation(float_val IN float(24));

-- Drop procedure with BOOLEAN and BOOL
DROP PROCEDURE boolean_check(active IN boolean, enabled IN bool);

-- Drop procedure with CHAR with length
DROP PROCEDURE char_processing(code IN char(10));

-- Drop procedure with CHARACTER with length
DROP PROCEDURE character_handling(name IN character(50));

-- Drop procedure with VARCHAR with length
DROP PROCEDURE string_processing(description IN varchar(255));

-- Drop procedure with CHARACTER VARYING
DROP PROCEDURE text_manipulation(content IN character varying(1000));

-- Drop procedure with DATE
DROP PROCEDURE date_processing(process_date IN date);

-- Drop procedure with TIMESTAMP variations
DROP PROCEDURE timestamp_proc(ts IN timestamp);
DROP PROCEDURE timestamp_precise(ts IN timestamp(6));
DROP PROCEDURE timestamp_with_tz(ts IN timestamp with time zone);
DROP PROCEDURE timestamp_without_tz(ts IN timestamp without time zone);

-- Drop procedure with TIMESTAMPTZ variations
DROP PROCEDURE timestamptz_proc(ts IN timestamptz);
DROP PROCEDURE timestamptz_precise(ts IN timestamptz(3));

-- Drop procedure with TIME variations
DROP PROCEDURE time_proc(time_val IN time);
DROP PROCEDURE time_precise(time_val IN time(6));
DROP PROCEDURE time_with_tz(time_val IN time with time zone);
DROP PROCEDURE time_without_tz(time_val IN time without time zone);

-- Drop procedure with TIMETZ variations
DROP PROCEDURE timetz_proc(time_val IN timetz);
DROP PROCEDURE timetz_precise(time_val IN timetz(6));

-- Drop procedure with TEXT
DROP PROCEDURE text_processing(content IN text);

-- Redshift-specific data types

-- Drop procedure with SUPER
DROP PROCEDURE json_processing(data IN super);

-- Drop procedure with GEOMETRY
DROP PROCEDURE geospatial_calc(shape IN geometry);

-- Drop procedure with GEOGRAPHY
DROP PROCEDURE geographic_analysis(location IN geography);

-- Drop procedure with HLLSKETCH
DROP PROCEDURE hyperloglog_processing(sketch IN hllsketch);

-- Procedures with custom/identifier types
DROP PROCEDURE custom_type_proc(custom_data IN user_defined_type);

-- Complex parameter combinations

-- Drop procedure with many parameters (up to 32 limit)
DROP PROCEDURE complex_procedure(
    param1 IN integer,
    param2 IN varchar(50),
    param3 IN decimal(10,2),
    param4 OUT boolean,
    param5 INOUT text,
    param6 IN date,
    param7 IN timestamp,
    param8 OUT geometry,
    param9 IN super,
    param10 INOUT float
);

-- Drop procedure with maximum parameters example
DROP PROCEDURE max_params_example(
    p1 IN integer, p2 IN varchar, p3 IN decimal, p4 IN boolean, p5 IN date,
    p6 IN timestamp, p7 IN float, p8 IN text, p9 IN bigint, p10 IN smallint,
    p11 OUT integer, p12 OUT varchar, p13 OUT decimal, p14 OUT boolean, p15 OUT date,
    p16 INOUT timestamp, p17 INOUT float, p18 INOUT text, p19 INOUT bigint, p20 INOUT smallint,
    p21 IN geometry, p22 IN geography, p23 IN super, p24 IN hllsketch, p25 IN char(10),
    p26 OUT character(20), p27 OUT timestamptz, p28 OUT timetz, p29 INOUT real, p30 INOUT numeric(10,2)
);

-- Argument names are optional and ignored during drop

-- Drop procedure where argument names are omitted
DROP PROCEDURE unnamed_params(integer, varchar, OUT decimal);

-- Drop procedure with some names, some without
DROP PROCEDURE mixed_naming(integer, description varchar, OUT decimal);

-- Drop procedure with names that don't match original (names ignored)
DROP PROCEDURE name_mismatch(different_name integer, another_name varchar);

-- Real-world examples based on typical stored procedure use cases

-- Drop financial calculation procedures
DROP PROCEDURE calculate_compound_interest(
    principal IN decimal(15,2),
    rate IN decimal(5,4),
    periods IN integer,
    final_amount OUT decimal(15,2)
);

DROP PROCEDURE process_monthly_billing(
    customer_id IN integer,
    billing_month IN integer,
    billing_year IN integer,
    total_charges OUT decimal(12,2),
    status OUT varchar(20)
);

-- Drop data processing procedures
DROP PROCEDURE etl_transform_customer_data(
    source_table IN varchar(100),
    target_table IN varchar(100),
    transformation_rules IN super,
    rows_processed OUT bigint
);

DROP PROCEDURE validate_data_quality(
    table_name IN varchar(100),
    validation_rules IN super,
    error_count OUT integer,
    validation_report OUT text
);

-- Drop reporting procedures
DROP PROCEDURE generate_sales_report(
    start_date IN date,
    end_date IN date,
    region IN varchar(50),
    report_data OUT super
);

DROP PROCEDURE update_dashboard_metrics(
    metric_date IN date,
    refresh_cache IN boolean
);

-- Drop maintenance procedures
DROP PROCEDURE cleanup_temporary_tables(
    retention_days IN integer,
    tables_dropped OUT integer
);

DROP PROCEDURE analyze_table_statistics(
    schema_name IN varchar(100),
    table_pattern IN varchar(100)
);

-- Edge cases and special scenarios

-- Drop procedure with very long name
DROP PROCEDURE very_long_procedure_name_that_demonstrates_maximum_identifier_length_support(param IN integer);

-- Drop procedure with quoted identifiers containing special characters
DROP PROCEDURE "procedure-with-hyphens"(input_val IN integer);
DROP PROCEDURE "UPPERCASE_PROCEDURE"(data IN varchar);
DROP PROCEDURE "procedure with spaces"(value IN decimal);

-- Drop procedure in system-like schemas
DROP PROCEDURE pg_catalog.system_procedure(config IN text);
DROP PROCEDURE information_schema.info_procedure(metadata IN super);

-- Multiple procedures with same name but different signatures
DROP PROCEDURE overloaded_proc(param IN integer);
DROP PROCEDURE overloaded_proc(param IN varchar);
DROP PROCEDURE overloaded_proc(param1 IN integer, param2 IN varchar);

-- Procedures with complex data type specifications
DROP PROCEDURE complex_types_proc(
    precise_decimal IN decimal(38,18),
    big_varchar IN varchar(65535),
    precise_timestamp IN timestamp(6) with time zone,
    high_precision_float IN float(53),
    long_char IN character(4096)
);

-- Default behavior demonstrations

-- Drop procedure with default behavior (RESTRICT implied)
DROP PROCEDURE default_behavior_proc(data IN text);

-- Drop procedure with explicit CASCADE
DROP PROCEDURE cascade_example_proc(input_val IN integer) CASCADE;

-- Drop procedure with explicit RESTRICT
DROP PROCEDURE restrict_example_proc(config IN super) RESTRICT;

-- Combined options testing

-- All combinations of IF EXISTS with CASCADE/RESTRICT
DROP PROCEDURE IF EXISTS combo_test_cascade(val IN integer) CASCADE;
DROP PROCEDURE IF EXISTS combo_test_restrict(val IN varchar) RESTRICT;

-- Schema-qualified with IF EXISTS and CASCADE/RESTRICT
DROP PROCEDURE IF EXISTS myschema.qualified_cascade_proc(data IN text) CASCADE;
DROP PROCEDURE IF EXISTS "quoted_schema".qualified_restrict_proc(value IN decimal) RESTRICT;

-- Complex real-world scenarios

-- Drop machine learning model procedures
DROP PROCEDURE train_customer_segmentation_model(
    training_data_table IN varchar(100),
    model_name IN varchar(100),
    hyperparameters IN super,
    model_accuracy OUT decimal(5,4)
);

-- Drop data science procedures
DROP PROCEDURE calculate_statistical_summary(
    dataset_table IN varchar(100),
    columns_list IN super,
    summary_stats OUT super
);

-- Drop integration procedures
DROP PROCEDURE sync_external_data(
    source_system IN varchar(50),
    sync_timestamp IN timestamp,
    records_updated OUT bigint,
    sync_status OUT varchar(20)
);

-- Drop monitoring procedures
DROP PROCEDURE check_system_health(
    component_name IN varchar(100),
    health_status OUT varchar(20),
    metrics OUT super
);