-- CALL command examples from AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CALL_procedure.html

-- Example 1: Simple procedure call without parameters
CALL refresh_materialized_view();

-- Example 2: Procedure with positional parameters
CALL update_customer_stats(2023, 'Q4');

-- Example 3: Procedure with named parameters using => syntax
CALL process_orders(start_date => '2023-01-01', end_date => '2023-12-31', status => 'completed');

-- Example 4: Procedure with mixed positional and named parameters
CALL generate_report(2023, report_type => 'annual', include_details => true);

-- Example 5: Schema-qualified procedure call
CALL finance.calculate_revenue('2023-01-01', '2023-12-31');

-- Example 6: Procedure with numeric parameters
CALL test_sp1(5, 'abc');

-- Example 7: Procedure with INOUT parameters (from documentation example)
CALL test_sp2(2, '2019');

-- Example 8: Procedure call with function as parameter
CALL log_query_id(pg_last_query_id());

-- Example 9: Procedure with boolean parameters
CALL update_settings(enable_logging => true, debug_mode => false);

-- Example 10: Procedure with NULL parameter
CALL reset_user_data(user_id => 12345, clear_history => NULL);

-- Example 11: Procedure with string concatenation
CALL send_notification('Alert: ' + 'System maintenance scheduled');

-- Example 12: Nested procedure call example (outer procedure)
CALL outer_proc(5);

-- Example 13: Procedure with cursor output (would be used in transaction)
CALL get_customers_by_city('MOROCCO 0', 'thecursor');

-- Example 14: Procedure that creates temp table
CALL get_customers_by_city_temp_table('MOROCCO 0', 'mytemptable');

-- Example 15: Complex procedure with multiple parameter types
CALL complex_calculation(
    input_value => 100.50,
    multiplier => 1.5,
    use_cache => true,
    description => 'Monthly calculation'
);

-- Example 16: Procedure with array-like string parameter
CALL process_batch('item1,item2,item3', batch_size => 10);

-- Example 17: Procedure with date calculations
CALL generate_monthly_report(
    start_date => '2023-01-01'::date,
    months => 12
);

-- Example 18: Procedure with quoted identifiers
CALL "Special Procedure"(param1 => 'value1');

-- Example 19: Procedure with numeric precision
CALL financial_calculation(amount => 12345.67, rate => 0.0525);

-- Example 20: Procedure with mixed case identifiers
CALL MixedCaseProc(InputParam => 'test');

-- Example 21: Procedure with parentheses in expressions
CALL calculate_total((base_amount * tax_rate) + shipping_cost);

-- Example 22: Procedure with nested function calls
CALL log_operation(
    operation_id => extract(epoch from current_timestamp),
    details => upper(trim('  operation completed  '))
);

-- Example 23: Procedure using = syntax for named parameters
CALL update_configuration(setting_name = 'timeout', setting_value = '300');

-- Example 24: Simple procedure with single numeric parameter
CALL increment_counter(1);

-- Example 25: Procedure with multiple string parameters
CALL create_user_profile(
    username => 'john_doe',
    email => 'john@example.com',
    department => 'Engineering'
);

-- Example 26: Procedure with boolean logic
CALL validate_data(
    check_duplicates => true,
    strict_mode => false,
    auto_fix => true
);

-- Example 27: Procedure with mathematical expressions
CALL calculate_metrics(
    value1 => 100 + 50,
    value2 => 200 * 1.5,
    value3 => (300 / 2) - 25
);

-- Example 28: Long procedure name with underscores
CALL very_long_procedure_name_with_many_parameters(
    param_one => 'value1',
    param_two => 'value2',
    param_three => 'value3'
);

-- Example 29: Procedure with special characters in strings
CALL process_text(input_text => 'Text with "quotes" and ''apostrophes''');

-- Example 30: Minimal procedure call with empty parameters
CALL simple_procedure();