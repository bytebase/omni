-- SHOW PROCEDURE Examples
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_PROCEDURE.html

-- Basic usage: Show procedure with argument types
SHOW PROCEDURE test_sp2(int, varchar);

-- Show procedure with no arguments
SHOW PROCEDURE simple_proc();

-- Show procedure without parentheses (when unique)
SHOW PROCEDURE unique_procedure_name;

-- Show procedure with multiple arguments
SHOW PROCEDURE calculate_totals(integer, decimal(10,2), varchar(100));

-- Show procedure with named arguments and modes
SHOW PROCEDURE process_order(IN order_id integer, INOUT status varchar, OUT result decimal);

-- Show procedure with complex argument types
SHOW PROCEDURE data_transform(
    bigint,
    timestamp,
    numeric(18,4),
    boolean,
    char(10),
    date
);

-- Show overloaded procedure with specific signature
SHOW PROCEDURE update_inventory(smallint);
SHOW PROCEDURE update_inventory(integer, varchar);
SHOW PROCEDURE update_inventory(bigint, text, timestamp);

-- Show procedure with array types
SHOW PROCEDURE array_processor(integer[]);

-- Show procedure with mixed argument modes
SHOW PROCEDURE complex_proc(
    IN input_val integer,
    OUT output_val varchar,
    INOUT bidirectional_val numeric
);

-- Show procedure with optional argument names
SHOW PROCEDURE named_args_proc(
    customer_id integer,
    order_date date,
    amount decimal(10,2)
);

-- Show procedure in specific schema (using qualified name)
SHOW PROCEDURE myschema.special_procedure(integer);

-- Show procedure with no arguments (both syntaxes)
SHOW PROCEDURE no_args_proc;
SHOW PROCEDURE no_args_proc();

-- Show procedure with text/varchar variations
SHOW PROCEDURE string_processor(text);
SHOW PROCEDURE string_processor(varchar);
SHOW PROCEDURE string_processor(varchar(255));
SHOW PROCEDURE string_processor(character varying(100));

-- Show procedure with numeric variations
SHOW PROCEDURE calc_proc(numeric);
SHOW PROCEDURE calc_proc(decimal);
SHOW PROCEDURE calc_proc(numeric(10,2));
SHOW PROCEDURE calc_proc(decimal(18,4));

-- Show procedure with real/double precision
SHOW PROCEDURE float_calc(real);
SHOW PROCEDURE float_calc(double precision);
SHOW PROCEDURE float_calc(float4);
SHOW PROCEDURE float_calc(float8);

-- Show procedure with timestamp variations
SHOW PROCEDURE time_proc(timestamp);
SHOW PROCEDURE time_proc(timestamp without time zone);
SHOW PROCEDURE time_proc(timestamp with time zone);
SHOW PROCEDURE time_proc(timestamptz);

-- Show procedure with interval type
SHOW PROCEDURE interval_calc(interval);

-- Show procedure with boolean type
SHOW PROCEDURE flag_processor(boolean);

-- Mixed case procedure names
SHOW PROCEDURE "MixedCaseProcedure"(integer);
SHOW PROCEDURE "UPPERCASE_PROC"(varchar);

-- Procedure with underscore and numbers in name
SHOW PROCEDURE proc_v2_final(integer, varchar);
SHOW PROCEDURE _private_proc(text);
SHOW PROCEDURE proc123(numeric);

-- Show procedure with all argument modes specified
SHOW PROCEDURE full_spec_proc(
    IN param1 integer,
    IN param2 varchar(100),
    OUT result1 decimal,
    OUT result2 text,
    INOUT status_flag boolean
);