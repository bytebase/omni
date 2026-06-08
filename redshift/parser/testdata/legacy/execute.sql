-- EXECUTE command test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_EXECUTE.html

-- Basic EXECUTE without parameters
EXECUTE simple_plan;

-- EXECUTE with single parameter
EXECUTE prep_select_plan(2);

-- EXECUTE with multiple parameters (integers)
EXECUTE prep_insert_plan(1, 2);

-- EXECUTE with mixed parameter types (int, string)
EXECUTE prep_insert_plan(1, 'one');

-- EXECUTE with string parameters
EXECUTE prep_update_plan('John', 'Doe');

-- EXECUTE with numeric parameters (decimals)
EXECUTE price_query(19.99, 100.50);

-- EXECUTE with boolean parameters
EXECUTE status_check(true, false);

-- EXECUTE with NULL parameter
EXECUTE null_test(null);

-- Examples from AWS documentation:

-- Example 1: Inserting multiple rows into a temporary table
-- PREPARE prep_insert_plan (int, char) AS insert into prep1 values ($1, $2);
EXECUTE prep_insert_plan(1, 'one');
EXECUTE prep_insert_plan(2, 'two');
EXECUTE prep_insert_plan(3, 'three');

-- Example 2: Selecting rows with a parameterized condition
-- PREPARE prep_select_plan (int) AS select * from prep1 where c1 = $1;
EXECUTE prep_select_plan(2);
EXECUTE prep_select_plan(3);

-- Additional test cases:

-- EXECUTE with complex identifiers
EXECUTE my_complex_plan_name;
EXECUTE plan_with_underscore_123;

-- EXECUTE with quoted strings containing special characters
EXECUTE text_plan('It''s a test');
EXECUTE special_chars('Hello, World!');

-- EXECUTE with negative numbers
EXECUTE math_plan(-1, -99.99);

-- EXECUTE with scientific notation
EXECUTE scientific_plan(1.23e-4, 5.67E+10);

-- EXECUTE with large numbers
EXECUTE big_numbers(999999999, 123456789.987654321);

-- EXECUTE with multiple mixed parameters
EXECUTE complex_plan(1, 'text', true, null, 3.14, false);

-- EXECUTE with single quotes in string literals
EXECUTE quote_test('Don''t worry');
EXECUTE apostrophe_test('That''s correct');

-- EXECUTE with empty string
EXECUTE empty_string_test('');

-- EXECUTE with zero values
EXECUTE zero_test(0, 0.0);

-- Case variations (should be case-insensitive)
execute lowercase_plan;
Execute mixed_case_plan;
EXECUTE uppercase_plan;