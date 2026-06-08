-- CREATE FUNCTION test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_FUNCTION.html

-- Basic Python UDF examples

-- Python UDF with STABLE volatility (recommended prefix 'f_')
CREATE FUNCTION f_py_greater (a float, b float)
RETURNS float STABLE AS $$
    if a > b:
        return a
    return b
$$ LANGUAGE plpythonu;

-- Python UDF with multiple parameter types
CREATE FUNCTION f_py_calculator (x integer, y integer, operation varchar)
RETURNS integer VOLATILE AS $$
    if operation == 'add':
        return x + y
    elif operation == 'subtract':
        return x - y
    elif operation == 'multiply':
        return x * y
    elif operation == 'divide':
        if y != 0:
            return x // y
        else:
            return None
    return 0
$$ LANGUAGE plpythonu;

-- Python UDF with IMMUTABLE volatility
CREATE FUNCTION f_py_square (num float)
RETURNS float IMMUTABLE AS $$
    return num * num
$$ LANGUAGE plpythonu;

-- Python UDF with boolean return
CREATE FUNCTION f_py_is_even (num integer)
RETURNS boolean STABLE AS $$
    return num % 2 == 0
$$ LANGUAGE plpythonu;

-- Python UDF with string manipulation
CREATE FUNCTION f_py_clean_string (input_str varchar)
RETURNS varchar STABLE AS $$
    if input_str is None:
        return ''
    return input_str.strip().lower()
$$ LANGUAGE plpythonu;

-- Python UDF with date handling
CREATE FUNCTION f_py_days_between (start_date date, end_date date)
RETURNS integer STABLE AS $$
    if start_date is None or end_date is None:
        return None
    diff = end_date - start_date
    return diff.days
$$ LANGUAGE plpythonu;

-- Basic SQL UDF examples

-- SQL UDF with STABLE volatility
CREATE FUNCTION f_sql_greater (float, float)
RETURNS float STABLE AS $$
    SELECT CASE WHEN $1 > $2 THEN $1 ELSE $2 END
$$ LANGUAGE sql;

-- SQL UDF with multiple parameters
CREATE FUNCTION f_sql_discount_price (original_price decimal, discount_percent decimal)
RETURNS decimal STABLE AS $$
    SELECT $1 * (1 - $2 / 100)
$$ LANGUAGE sql;

-- SQL UDF with IMMUTABLE volatility
CREATE FUNCTION f_sql_circle_area (radius float)
RETURNS float IMMUTABLE AS $$
    SELECT 3.14159 * $1 * $1
$$ LANGUAGE sql;

-- SQL UDF with string operations
CREATE FUNCTION f_sql_full_name (first_name varchar, last_name varchar)
RETURNS varchar STABLE AS $$
    SELECT $1 || ' ' || $2
$$ LANGUAGE sql;

-- SQL UDF with conditional logic
CREATE FUNCTION f_sql_grade_letter (score integer)
RETURNS varchar STABLE AS $$
    SELECT CASE
        WHEN $1 >= 90 THEN 'A'
        WHEN $1 >= 80 THEN 'B'
        WHEN $1 >= 70 THEN 'C'
        WHEN $1 >= 60 THEN 'D'
        ELSE 'F'
    END
$$ LANGUAGE sql;

-- OR REPLACE examples

-- Replace existing function
CREATE OR REPLACE FUNCTION f_py_greater (a float, b float)
RETURNS float STABLE AS $$
    # Updated version with null handling
    if a is None and b is None:
        return None
    elif a is None:
        return b
    elif b is None:
        return a
    else:
        return a if a > b else b
$$ LANGUAGE plpythonu;

-- Replace SQL function
CREATE OR REPLACE FUNCTION f_sql_greater (float, float)
RETURNS float STABLE AS $$
    SELECT CASE
        WHEN $1 IS NULL AND $2 IS NULL THEN NULL
        WHEN $1 IS NULL THEN $2
        WHEN $2 IS NULL THEN $1
        WHEN $1 > $2 THEN $1
        ELSE $2
    END
$$ LANGUAGE sql;

-- Advanced data type examples

-- Function with DECIMAL parameters
CREATE FUNCTION f_sql_compound_interest (principal decimal(10,2), rate decimal(5,4), years integer)
RETURNS decimal(10,2) STABLE AS $$
    SELECT $1 * POWER(1 + $2, $3)
$$ LANGUAGE sql;

-- Function with TIMESTAMP parameters
CREATE FUNCTION f_sql_business_days (start_ts timestamp, end_ts timestamp)
RETURNS integer STABLE AS $$
    SELECT CASE
        WHEN $1 IS NULL OR $2 IS NULL THEN NULL
        ELSE DATEDIFF(day, $1, $2)
    END
$$ LANGUAGE sql;

-- Function with VARCHAR with length
CREATE FUNCTION f_sql_format_phone (phone varchar(15))
RETURNS varchar(20) STABLE AS $$
    SELECT '(' || SUBSTRING($1, 1, 3) || ') ' ||
           SUBSTRING($1, 4, 3) || '-' ||
           SUBSTRING($1, 7, 4)
$$ LANGUAGE sql;

-- Function with CHAR parameters
CREATE FUNCTION f_sql_initial (first_name char(50))
RETURNS char(1) STABLE AS $$
    SELECT LEFT(TRIM($1), 1)
$$ LANGUAGE sql;

-- Complex Python UDF examples

-- Python UDF with JSON-like string manipulation
CREATE FUNCTION f_py_extract_json_field (json_str varchar, field_name varchar)
RETURNS varchar VOLATILE AS $$
    import json
    try:
        data = json.loads(json_str)
        return str(data.get(field_name, ''))
    except:
        return ''
$$ LANGUAGE plpythonu;

-- Python UDF with mathematical operations
CREATE FUNCTION f_py_fibonacci (n integer)
RETURNS bigint STABLE AS $$
    if n <= 0:
        return 0
    elif n == 1:
        return 1
    else:
        a, b = 0, 1
        for i in range(2, n + 1):
            a, b = b, a + b
        return b
$$ LANGUAGE plpythonu;

-- Python UDF with array-like operations
CREATE FUNCTION f_py_word_count (text_input varchar)
RETURNS integer STABLE AS $$
    if text_input is None:
        return 0
    words = text_input.split()
    return len(words)
$$ LANGUAGE plpythonu;

-- Edge cases and special scenarios

-- Function with no parameters
CREATE FUNCTION f_sql_current_year ()
RETURNS integer STABLE AS $$
    SELECT EXTRACT(year FROM GETDATE())
$$ LANGUAGE sql;

-- Function with ANYELEMENT type (generic)
CREATE FUNCTION f_sql_first_non_null (anyelement, anyelement)
RETURNS anyelement STABLE AS $$
    SELECT COALESCE($1, $2)
$$ LANGUAGE sql;

-- Function with maximum parameters (approaching the 32 limit)
CREATE FUNCTION f_py_sum_many (
    a integer, b integer, c integer, d integer, e integer,
    f integer, g integer, h integer, i integer, j integer
)
RETURNS integer STABLE AS $$
    values = [a, b, c, d, e, f, g, h, i, j]
    return sum(v for v in values if v is not None)
$$ LANGUAGE plpythonu;

-- VOLATILE function example (different results possible)
CREATE FUNCTION f_py_random_int (min_val integer, max_val integer)
RETURNS integer VOLATILE AS $$
    import random
    return random.randint(min_val, max_val)
$$ LANGUAGE plpythonu;

-- Function with SUPER data type (Redshift-specific)
CREATE FUNCTION f_sql_super_length (input_super super)
RETURNS integer STABLE AS $$
    SELECT JSON_ARRAY_LENGTH($1)
$$ LANGUAGE sql;

-- Function with GEOMETRY type
CREATE FUNCTION f_sql_geometry_area (geom geometry)
RETURNS float STABLE AS $$
    SELECT ST_AREA($1)
$$ LANGUAGE sql;

-- Function with GEOGRAPHY type
CREATE FUNCTION f_sql_geography_distance (geog1 geography, geog2 geography)
RETURNS float STABLE AS $$
    SELECT ST_DISTANCE($1, $2)
$$ LANGUAGE sql;

-- Function with HLLSKETCH type
CREATE FUNCTION f_sql_hll_cardinality (sketch hllsketch)
RETURNS bigint STABLE AS $$
    SELECT HLL_CARDINALITY($1)
$$ LANGUAGE sql;

-- Complex SQL UDF with subquery-like logic
CREATE FUNCTION f_sql_tax_calculation (income decimal(12,2), tax_rate decimal(5,4))
RETURNS decimal(12,2) STABLE AS $$
    SELECT CASE
        WHEN $1 <= 0 THEN 0
        WHEN $2 < 0 THEN 0
        WHEN $2 > 1 THEN $1
        ELSE $1 * $2
    END
$$ LANGUAGE sql;

-- Functions demonstrating nested calls capability
CREATE FUNCTION f_sql_nested_discount (price decimal(10,2), discount1 decimal(5,4), discount2 decimal(5,4))
RETURNS decimal(10,2) STABLE AS $$
    SELECT f_sql_discount_price(f_sql_discount_price($1, $2), $3)
$$ LANGUAGE sql;