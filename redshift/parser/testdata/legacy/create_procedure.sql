-- CREATE PROCEDURE test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_PROCEDURE.html

-- Basic procedure examples

-- Simple procedure with input parameters
CREATE PROCEDURE test_sp1(f1 int, f2 varchar(20)) AS $$
DECLARE
    min_val int;
BEGIN
    DROP TABLE IF EXISTS tmp_tbl;
    CREATE TEMP TABLE tmp_tbl(id int);
    INSERT INTO tmp_tbl values (f1),(10001),(10002);
    SELECT INTO min_val MIN(id) FROM tmp_tbl;
    RAISE INFO 'min_val = %, f2 = %', min_val, f2;
END;
$$ LANGUAGE plpgsql;

-- Procedure with IN, OUT, INOUT parameters
CREATE PROCEDURE test_sp2(
    f1 IN int,
    f2 INOUT varchar(256),
    out_var OUT varchar(256)
) AS $$
DECLARE
    loop_var int;
BEGIN
    IF f1 is null OR f2 is null THEN
        RAISE EXCEPTION 'input cannot be null';
    END IF;
    DROP TABLE IF EXISTS tmp_tbl;
    CREATE TEMP TABLE tmp_tbl(id int);
    FOR loop_var IN 1..f1 LOOP
        INSERT INTO tmp_tbl VALUES (loop_var);
    END LOOP;
    SELECT INTO out_var 'INFO: loop_var = ' || loop_var || ', f2 = ' || f2;
    f2 := 'updated f2';
END;
$$ LANGUAGE plpgsql;

-- Procedure with parameter mode variations
CREATE PROCEDURE test_modes(
    input_param IN integer,
    output_param OUT varchar(100),
    inout_param INOUT decimal(10,2)
) AS $$
DECLARE
    temp_value decimal(10,2);
BEGIN
    temp_value := inout_param * input_param;
    output_param := 'Calculated value: ' || temp_value::varchar;
    inout_param := temp_value;
END;
$$ LANGUAGE plpgsql;

-- OR REPLACE examples

-- Replace existing procedure
CREATE OR REPLACE PROCEDURE test_sp1(f1 int, f2 varchar(20)) AS $$
DECLARE
    min_val int;
    max_val int;
BEGIN
    DROP TABLE IF EXISTS tmp_tbl;
    CREATE TEMP TABLE tmp_tbl(id int);
    INSERT INTO tmp_tbl values (f1),(10001),(10002);
    SELECT INTO min_val, max_val MIN(id), MAX(id) FROM tmp_tbl;
    RAISE INFO 'min_val = %, max_val = %, f2 = %', min_val, max_val, f2;
END;
$$ LANGUAGE plpgsql;

-- NONATOMIC procedures

-- NONATOMIC procedure that automatically commits
CREATE PROCEDURE test_nonatomic_sp(input_id int) NONATOMIC AS $$
DECLARE
    rec record;
BEGIN
    -- Each statement auto-commits in NONATOMIC mode
    DROP TABLE IF EXISTS test_table;
    CREATE TABLE test_table(id int);
    INSERT INTO test_table VALUES (input_id);

    FOR rec IN SELECT * FROM test_table LOOP
        RAISE INFO 'Processing id: %', rec.id;
    END LOOP;

    DROP TABLE test_table;
END;
$$ LANGUAGE plpgsql;

-- NONATOMIC with transaction control
CREATE OR REPLACE PROCEDURE test_transaction_sp(batch_size int) NONATOMIC AS $$
DECLARE
    counter int := 0;
    current_id int;
BEGIN
    DROP TABLE IF EXISTS batch_table;
    CREATE TABLE batch_table(id int, processed boolean default false);

    -- Insert test data
    FOR current_id IN 1..batch_size * 3 LOOP
        INSERT INTO batch_table(id) VALUES (current_id);
        counter := counter + 1;

        -- Commit every batch_size records
        IF counter % batch_size = 0 THEN
            COMMIT;
            RAISE INFO 'Committed batch of % records', batch_size;
        END IF;
    END LOOP;

    COMMIT; -- Final commit
END;
$$ LANGUAGE plpgsql;

-- Security options

-- SECURITY INVOKER procedure (default)
CREATE PROCEDURE test_security_invoker(table_name varchar(100))
SECURITY INVOKER AS $$
BEGIN
    -- Runs with caller's privileges
    EXECUTE 'SELECT COUNT(*) FROM ' || table_name;
    RAISE INFO 'Accessed table: %', table_name;
END;
$$ LANGUAGE plpgsql;

-- SECURITY DEFINER procedure
CREATE PROCEDURE test_security_definer(user_id int)
SECURITY DEFINER AS $$
DECLARE
    user_count int;
BEGIN
    -- Runs with procedure owner's privileges
    SELECT COUNT(*) INTO user_count FROM pg_user WHERE usesysid = user_id;
    RAISE INFO 'User count for ID %: %', user_id, user_count;
END;
$$ LANGUAGE plpgsql;

-- SET configuration parameter examples

-- Procedure with configuration parameter
CREATE PROCEDURE test_config_sp(search_query varchar(100))
SET search_path = public, temp_schema AS $$
BEGIN
    -- search_path is temporarily set for this procedure
    DROP TABLE IF EXISTS search_results;
    CREATE TEMP TABLE search_results(result_text varchar(200));
    INSERT INTO search_results VALUES ('Found: ' || search_query);
    RAISE INFO 'Search completed for: %', search_query;
END;
$$ LANGUAGE plpgsql;

-- Multiple configuration parameters
CREATE PROCEDURE test_multi_config_sp(log_level varchar(10))
SET log_min_messages = info
SET work_mem = '256MB' AS $$
BEGIN
    RAISE INFO 'Procedure started with log level: %', log_level;
    -- Memory and logging settings are temporarily modified
    PERFORM pg_sleep(1);
    RAISE INFO 'Procedure completed';
END;
$$ LANGUAGE plpgsql;

-- Data type examples

-- Procedure with various data types
CREATE PROCEDURE test_datatypes_sp(
    small_num smallint,
    big_num bigint,
    precise_num decimal(10,2),
    float_num real,
    double_num double precision,
    bool_val boolean,
    char_val char(10),
    varchar_val varchar(255),
    date_val date,
    timestamp_val timestamp,
    text_val text
) AS $$
DECLARE
    result_text text;
BEGIN
    result_text := 'Values: ' || small_num || ', ' || big_num || ', ' || precise_num;
    result_text := result_text || ', ' || float_num || ', ' || double_num;
    result_text := result_text || ', ' || bool_val || ', ' || trim(char_val);
    result_text := result_text || ', ' || varchar_val || ', ' || date_val;
    result_text := result_text || ', ' || timestamp_val || ', ' || text_val;

    RAISE INFO 'All data types: %', result_text;
END;
$$ LANGUAGE plpgsql;

-- Redshift-specific data types
CREATE PROCEDURE test_redshift_types_sp(
    super_val super,
    geom_val geometry,
    geog_val geography,
    hll_val hllsketch
) AS $$
BEGIN
    RAISE INFO 'Super type length: %', JSON_ARRAY_LENGTH(super_val);
    RAISE INFO 'Geometry area: %', ST_AREA(geom_val);
    RAISE INFO 'Geography provided: %', CASE WHEN geog_val IS NOT NULL THEN 'Yes' ELSE 'No' END;
    RAISE INFO 'HLL cardinality: %', HLL_CARDINALITY(hll_val);
END;
$$ LANGUAGE plpgsql;

-- Control flow examples

-- Procedure with IF-ELSIF-ELSE
CREATE PROCEDURE test_control_flow_sp(grade_score int) AS $$
DECLARE
    letter_grade varchar(2);
    message varchar(100);
BEGIN
    IF grade_score >= 90 THEN
        letter_grade := 'A';
        message := 'Excellent';
    ELSIF grade_score >= 80 THEN
        letter_grade := 'B';
        message := 'Good';
    ELSIF grade_score >= 70 THEN
        letter_grade := 'C';
        message := 'Average';
    ELSIF grade_score >= 60 THEN
        letter_grade := 'D';
        message := 'Below Average';
    ELSE
        letter_grade := 'F';
        message := 'Failing';
    END IF;

    RAISE INFO 'Score: %, Grade: %, Assessment: %', grade_score, letter_grade, message;
END;
$$ LANGUAGE plpgsql;

-- Procedure with loops
CREATE PROCEDURE test_loops_sp(max_iterations int) AS $$
DECLARE
    counter int := 1;
    factorial bigint := 1;
    current_num int;
BEGIN
    -- WHILE loop
    WHILE counter <= max_iterations LOOP
        factorial := factorial * counter;
        counter := counter + 1;
    END LOOP;

    RAISE INFO 'Factorial of %: %', max_iterations, factorial;

    -- FOR loop with range
    FOR current_num IN 1..max_iterations LOOP
        IF current_num % 2 = 0 THEN
            RAISE INFO 'Even number: %', current_num;
        END IF;
    END LOOP;

    -- FOR loop with REVERSE
    FOR current_num IN REVERSE max_iterations..1 LOOP
        IF current_num <= 3 THEN
            EXIT; -- Exit the loop early
        END IF;
        RAISE INFO 'Countdown: %', current_num;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Procedure with labeled loops and EXIT/CONTINUE
CREATE PROCEDURE test_nested_loops_sp(outer_limit int, inner_limit int) AS $$
DECLARE
    i int;
    j int;
    product int;
BEGIN
    <<outer_loop>>
    FOR i IN 1..outer_limit LOOP
        <<inner_loop>>
        FOR j IN 1..inner_limit LOOP
            product := i * j;

            -- Skip even products
            IF product % 2 = 0 THEN
                CONTINUE inner_loop;
            END IF;

            -- Exit outer loop if product > 50
            IF product > 50 THEN
                EXIT outer_loop;
            END IF;

            RAISE INFO 'Product: % x % = %', i, j, product;
        END LOOP inner_loop;
    END LOOP outer_loop;
END;
$$ LANGUAGE plpgsql;

-- Exception handling examples

-- Procedure with exception handling
CREATE PROCEDURE test_exception_sp(divisor int) AS $$
DECLARE
    result decimal(10,2);
    dividend constant int := 100;
BEGIN
    BEGIN
        result := dividend / divisor;
        RAISE INFO 'Result: % / % = %', dividend, divisor, result;
    EXCEPTION
        WHEN division_by_zero THEN
            RAISE WARNING 'Division by zero attempted';
            result := NULL;
        WHEN others THEN
            RAISE WARNING 'Unexpected error occurred: %', SQLERRM;
            result := NULL;
    END;

    IF result IS NOT NULL THEN
        INSERT INTO calculation_log VALUES (dividend, divisor, result, current_timestamp);
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Custom exception handling
CREATE PROCEDURE test_custom_exception_sp(input_value int) AS $$
DECLARE
    processed_value int;
BEGIN
    IF input_value < 0 THEN
        RAISE EXCEPTION 'Negative values not allowed'
            USING HINT = 'Please provide a positive integer',
                  ERRCODE = 'P0001';
    END IF;

    IF input_value > 1000 THEN
        RAISE EXCEPTION 'Value too large: %', input_value
            USING HINT = 'Maximum allowed value is 1000',
                  ERRCODE = 'P0002';
    END IF;

    processed_value := input_value * 2;
    RAISE INFO 'Processed value: %', processed_value;
END;
$$ LANGUAGE plpgsql;

-- Advanced SQL operations

-- Procedure with dynamic SQL
CREATE PROCEDURE test_dynamic_sql_sp(table_name varchar(100), column_name varchar(100)) AS $$
DECLARE
    sql_query text;
    record_count int;
BEGIN
    -- Build dynamic query
    sql_query := 'SELECT COUNT(*) FROM ' || quote_ident(table_name) ||
                 ' WHERE ' || quote_ident(column_name) || ' IS NOT NULL';

    -- Execute dynamic SQL
    EXECUTE sql_query INTO record_count;

    RAISE INFO 'Records in %.% with non-null values: %',
               table_name, column_name, record_count;
END;
$$ LANGUAGE plpgsql;

-- Procedure with cursor operations
CREATE PROCEDURE test_cursor_sp(department_id int) AS $$
DECLARE
    emp_cursor CURSOR FOR
        SELECT employee_id, first_name, last_name, salary
        FROM employees
        WHERE dept_id = department_id
        ORDER BY salary DESC;
    emp_rec RECORD;
    total_salary decimal(12,2) := 0;
    emp_count int := 0;
BEGIN
    FOR emp_rec IN emp_cursor LOOP
        total_salary := total_salary + emp_rec.salary;
        emp_count := emp_count + 1;

        RAISE INFO 'Employee: % % - Salary: %',
                   emp_rec.first_name, emp_rec.last_name, emp_rec.salary;
    END LOOP;

    IF emp_count > 0 THEN
        RAISE INFO 'Department % - Total employees: %, Total salary: %, Average: %',
                   department_id, emp_count, total_salary, (total_salary / emp_count);
    ELSE
        RAISE INFO 'No employees found in department %', department_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Procedure with RETURN statements
CREATE PROCEDURE test_return_sp(operation varchar(10), value1 int, value2 int) AS $$
DECLARE
    result int;
BEGIN
    CASE operation
        WHEN 'add' THEN
            result := value1 + value2;
        WHEN 'subtract' THEN
            result := value1 - value2;
        WHEN 'multiply' THEN
            result := value1 * value2;
        WHEN 'divide' THEN
            IF value2 = 0 THEN
                RAISE INFO 'Cannot divide by zero';
                RETURN; -- Early return
            END IF;
            result := value1 / value2;
        ELSE
            RAISE INFO 'Unknown operation: %', operation;
            RETURN; -- Early return
    END CASE;

    RAISE INFO 'Operation: % % % = %', value1, operation, value2, result;
END;
$$ LANGUAGE plpgsql;

-- Complex real-world examples

-- Data processing procedure
CREATE PROCEDURE process_sales_data_sp(process_date date) AS $$
DECLARE
    processed_count int := 0;
    error_count int := 0;
    sales_rec RECORD;
BEGIN
    -- Process sales records for given date
    FOR sales_rec IN
        SELECT order_id, customer_id, amount, status
        FROM sales_orders
        WHERE order_date = process_date AND status = 'pending'
    LOOP
        BEGIN
            -- Validate and process each record
            IF sales_rec.amount > 0 AND sales_rec.customer_id IS NOT NULL THEN
                UPDATE sales_orders
                SET status = 'processed', processed_timestamp = current_timestamp
                WHERE order_id = sales_rec.order_id;

                processed_count := processed_count + 1;
            ELSE
                UPDATE sales_orders
                SET status = 'error', error_message = 'Invalid data'
                WHERE order_id = sales_rec.order_id;

                error_count := error_count + 1;
            END IF;

        EXCEPTION
            WHEN OTHERS THEN
                error_count := error_count + 1;
                RAISE WARNING 'Error processing order %: %',
                             sales_rec.order_id, SQLERRM;
        END;
    END LOOP;

    RAISE INFO 'Processing complete for %: % processed, % errors',
               process_date, processed_count, error_count;
END;
$$ LANGUAGE plpgsql;

-- Procedure with transaction management (NONATOMIC)
CREATE PROCEDURE batch_update_sp(batch_size int) NONATOMIC AS $$
DECLARE
    update_count int := 0;
    total_updated int := 0;
    batch_start_time timestamp;
BEGIN
    batch_start_time := current_timestamp;

    LOOP
        -- Update a batch of records
        UPDATE customer_status
        SET last_updated = current_timestamp
        WHERE last_updated < current_date - interval '30 days'
        AND status = 'active'
        LIMIT batch_size;

        GET DIAGNOSTICS update_count = ROW_COUNT;

        IF update_count = 0 THEN
            EXIT; -- No more records to update
        END IF;

        total_updated := total_updated + update_count;

        -- Commit the batch
        COMMIT;

        RAISE INFO 'Updated batch of % records, total: %',
                   update_count, total_updated;

        -- Small delay between batches
        PERFORM pg_sleep(0.1);
    END LOOP;

    RAISE INFO 'Batch update completed: % total records updated in %',
               total_updated, (current_timestamp - batch_start_time);
END;
$$ LANGUAGE plpgsql;

-- Procedure with no parameters
CREATE PROCEDURE maintenance_cleanup_sp() AS $$
DECLARE
    cleanup_count int;
    start_time timestamp;
BEGIN
    start_time := current_timestamp;

    -- Clean up temporary tables
    DROP TABLE IF EXISTS temp_processing_table;
    DROP TABLE IF EXISTS temp_staging_table;

    -- Clean up old log entries
    DELETE FROM procedure_logs
    WHERE log_date < current_date - interval '90 days';

    GET DIAGNOSTICS cleanup_count = ROW_COUNT;

    RAISE INFO 'Maintenance cleanup completed: % old logs removed in %',
               cleanup_count, (current_timestamp - start_time);
END;
$$ LANGUAGE plpgsql;

-- Edge cases and error scenarios

-- Procedure with maximum input parameters (approaching 32 limit)
CREATE PROCEDURE test_max_params_sp(
    p1 int, p2 int, p3 int, p4 int, p5 int, p6 int, p7 int, p8 int,
    p9 int, p10 int, p11 int, p12 int, p13 int, p14 int, p15 int, p16 int,
    p17 int, p18 int, p19 int, p20 int, p21 int, p22 int, p23 int, p24 int,
    p25 int, p26 int, p27 int, p28 int, p29 int, p30 int
) AS $$
DECLARE
    param_sum bigint := 0;
    param_array int[] := ARRAY[p1,p2,p3,p4,p5,p6,p7,p8,p9,p10,
                              p11,p12,p13,p14,p15,p16,p17,p18,p19,p20,
                              p21,p22,p23,p24,p25,p26,p27,p28,p29,p30];
    param_value int;
BEGIN
    FOREACH param_value IN ARRAY param_array LOOP
        param_sum := param_sum + COALESCE(param_value, 0);
    END LOOP;

    RAISE INFO 'Sum of 30 parameters: %', param_sum;
END;
$$ LANGUAGE plpgsql;

-- Procedure with mixed IN/OUT/INOUT parameters (max output limit test)
CREATE PROCEDURE test_mixed_params_sp(
    in1 IN int, in2 IN varchar(50), in3 IN date,
    out1 OUT int, out2 OUT varchar(100), out3 OUT timestamp,
    inout1 INOUT decimal(10,2), inout2 INOUT boolean
) AS $$
BEGIN
    -- Process inputs and set outputs
    out1 := in1 * 2;
    out2 := 'Processed: ' || in2;
    out3 := in3::timestamp + interval '1 day';

    -- Modify INOUT parameters
    inout1 := inout1 * 1.1;
    inout2 := NOT inout2;

    RAISE INFO 'Mixed parameter processing completed';
END;
$$ LANGUAGE plpgsql;