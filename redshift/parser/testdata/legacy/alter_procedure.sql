-- ALTER PROCEDURE examples from AWS documentation
-- https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_PROCEDURE.html

-- Example 1: Rename a procedure with full argument specification
ALTER PROCEDURE first_quarter_revenue(volume INOUT bigint, at_price IN numeric, result OUT int)
RENAME TO quarterly_revenue;

-- Example 2: Rename a procedure with simplified argument list (only types)
ALTER PROCEDURE first_quarter_revenue(bigint, numeric)
RENAME TO quarterly_revenue;

-- Example 3: Change procedure owner
ALTER PROCEDURE quarterly_revenue(bigint, numeric)
OWNER TO etl_user;

-- Additional test cases

-- Rename procedure without arguments
ALTER PROCEDURE simple_proc RENAME TO renamed_proc;

-- Rename procedure with schema qualification
ALTER PROCEDURE myschema.my_procedure(integer, varchar) RENAME TO new_procedure;

-- Change owner to CURRENT_USER
ALTER PROCEDURE quarterly_revenue(bigint, numeric) OWNER TO CURRENT_USER;

-- Change owner to SESSION_USER
ALTER PROCEDURE quarterly_revenue(bigint, numeric) OWNER TO SESSION_USER;

-- Procedure with multiple arguments and different modes
ALTER PROCEDURE complex_proc(
    input_val IN integer,
    output_val OUT varchar,
    inout_val INOUT timestamp
) RENAME TO renamed_complex_proc;

-- Procedure with named arguments
ALTER PROCEDURE calc_stats(
    p_start_date IN date,
    p_end_date IN date,
    p_result OUT numeric
) OWNER TO analytics_user;

-- Procedure with array types
ALTER PROCEDURE array_proc(arr_param integer[]) RENAME TO new_array_proc;

-- Procedure with complex data types
ALTER PROCEDURE json_proc(data json) OWNER TO data_team;

-- Procedure with multiple IN parameters
ALTER PROCEDURE multi_param_proc(
    param1 IN varchar(100),
    param2 IN integer,
    param3 IN boolean
) RENAME TO updated_multi_param_proc;

-- Procedure with no parameters but with empty parentheses
ALTER PROCEDURE empty_param_proc() OWNER TO admin_user;

-- Mixed argument modes without names
ALTER PROCEDURE mixed_modes(IN integer, OUT varchar, INOUT numeric)
RENAME TO new_mixed_modes;

-- Schema-qualified procedure with owner change
ALTER PROCEDURE finance.calculate_revenue(date, date)
OWNER TO finance_admin;

-- Procedure with numeric precision types
ALTER PROCEDURE precise_calc(
    amount IN numeric(10,2),
    rate IN decimal(5,4),
    result OUT numeric(15,2)
) RENAME TO precise_calculation;

-- Test case with quoted identifiers
ALTER PROCEDURE "Special Proc"(integer) RENAME TO regular_proc;

-- Test case with mixed case identifiers
ALTER PROCEDURE MixedCaseProc(varchar) OWNER TO new_owner;