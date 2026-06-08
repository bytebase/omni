-- CREATE MASKING POLICY test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_MASKING_POLICY.html

-- Basic CREATE MASKING POLICY statements

-- 1. Simple masking policy with constant value and explicit cast
CREATE MASKING POLICY mask_credit_card_full
WITH (credit_card VARCHAR(256))
USING ('000000XXXX0000'::TEXT);

-- 2. Masking policy with NULL value
CREATE MASKING POLICY hide_sensitive_data
WITH (sensitive_column VARCHAR(100))
USING (NULL);

-- 3. Masking policy with column reference (pass-through for some users)
CREATE MASKING POLICY passthrough_policy
WITH (data_column TEXT)
USING (data_column);

-- 4. Multiple input columns
CREATE MASKING POLICY multi_column_mask
WITH (first_name VARCHAR(50), last_name VARCHAR(100))
USING ('REDACTED');

-- 5. Numeric data masking with zero
CREATE MASKING POLICY zero_salary
WITH (salary DECIMAL(10,2))
USING (0::DECIMAL(10,2));

-- Complex masking expressions with functions

-- 6. Partial credit card masking using SUBSTRING
CREATE MASKING POLICY mask_credit_card_partial
WITH (credit_card VARCHAR(256))
USING (SUBSTRING(credit_card FROM 1 FOR 6) || 'XXXXX' || SUBSTRING(credit_card FROM 12 FOR 4));

-- 7. Email masking preserving domain
CREATE MASKING POLICY mask_email_partial
WITH (email VARCHAR(255))
USING (SUBSTRING(email FROM 1 FOR 3) || '****@' || SUBSTRING(email FROM POSITION('@' IN email) + 1));

-- 8. Phone number masking
CREATE MASKING POLICY mask_phone
WITH (phone VARCHAR(15))
USING ('(' || SUBSTRING(phone FROM 1 FOR 3) || ') XXX-XXXX');

-- 9. Date masking - year only
CREATE MASKING POLICY mask_birthdate_year
WITH (birthdate DATE)
USING (DATE_TRUNC('year', birthdate));

-- 10. Using built-in string functions
CREATE MASKING POLICY mask_with_upper
WITH (name VARCHAR(100))
USING (UPPER('REDACTED'));

-- CASE expressions for conditional masking

-- 11. Simple CASE expression based on data value
CREATE MASKING POLICY conditional_salary_mask
WITH (salary DECIMAL(10,2), role VARCHAR(50))
USING (
    CASE
        WHEN role = 'admin' THEN salary
        WHEN role = 'manager' THEN salary / 10 * 10
        ELSE 0::DECIMAL(10,2)
    END
);

-- 12. Complex CASE with multiple conditions
CREATE MASKING POLICY region_based_mask
WITH (customer_data TEXT, region VARCHAR(20), user_level INTEGER)
USING (
    CASE
        WHEN region = 'US' AND user_level > 5 THEN customer_data
        WHEN region = 'EU' THEN 'GDPR_MASKED'
        WHEN region = 'APAC' THEN SUBSTRING(customer_data FROM 1 FOR 5) || '***'
        ELSE 'INTERNATIONAL_MASKED'
    END
);

-- 13. Nested CASE expressions
CREATE MASKING POLICY nested_case_mask
WITH (data_value TEXT, department VARCHAR(30), clearance_level INTEGER)
USING (
    CASE
        WHEN department = 'FINANCE' THEN
            CASE
                WHEN clearance_level >= 8 THEN data_value
                WHEN clearance_level >= 5 THEN 'FINANCE_PARTIAL'
                ELSE 'FINANCE_RESTRICTED'
            END
        WHEN department = 'HR' THEN
            CASE
                WHEN clearance_level >= 7 THEN data_value
                ELSE 'HR_CONFIDENTIAL'
            END
        ELSE 'DEPARTMENT_MASKED'
    END
);

-- Mathematical operations

-- 14. Arithmetic operations for numeric masking
CREATE MASKING POLICY numeric_arithmetic_mask
WITH (revenue DECIMAL(15,2))
USING (revenue * 0.1);

-- 15. Division and rounding for approximate values
CREATE MASKING POLICY approximate_revenue
WITH (revenue BIGINT)
USING ((revenue / 1000) * 1000);

-- 16. Modulo operation for pattern masking
CREATE MASKING POLICY modulo_mask
WITH (id INTEGER)
USING (id % 1000);

-- Logical operations

-- 17. AND/OR operations in CASE
CREATE MASKING POLICY logical_operations_mask
WITH (data TEXT, is_public BOOLEAN, is_sensitive BOOLEAN)
USING (
    CASE
        WHEN is_public AND NOT is_sensitive THEN data
        WHEN NOT is_public OR is_sensitive THEN 'RESTRICTED'
        ELSE 'DEFAULT_MASK'
    END
);

-- 18. IS NULL/IS NOT NULL checks
CREATE MASKING POLICY null_check_mask
WITH (optional_data VARCHAR(200))
USING (
    CASE
        WHEN optional_data IS NULL THEN 'NO_DATA'
        WHEN optional_data IS NOT NULL THEN 'DATA_PRESENT'
    END
);

-- Comparison operations

-- 19. Age-based masking with comparisons
CREATE MASKING POLICY age_based_mask
WITH (age INTEGER, personal_info TEXT)
USING (
    CASE
        WHEN age >= 18 AND age <= 65 THEN 'ADULT_INFO'
        WHEN age < 18 THEN 'MINOR_INFO'
        WHEN age > 65 THEN 'SENIOR_INFO'
        ELSE 'UNKNOWN_AGE'
    END
);

-- 20. String length comparisons
CREATE MASKING POLICY length_based_mask
WITH (text_data VARCHAR(1000))
USING (
    CASE
        WHEN LENGTH(text_data) > 100 THEN 'LONG_TEXT_MASKED'
        WHEN LENGTH(text_data) > 50 THEN SUBSTRING(text_data FROM 1 FOR 50) || '...'
        ELSE text_data
    END
);

-- Type casting examples

-- 21. CAST function for type conversion
CREATE MASKING POLICY cast_function_mask
WITH (numeric_id INTEGER)
USING (CAST(numeric_id AS VARCHAR) || '-MASKED');

-- 22. :: operator for type casting
CREATE MASKING POLICY cast_operator_mask
WITH (price_text VARCHAR(20))
USING ((price_text::DECIMAL * 0.9)::VARCHAR);

-- 23. Multiple type conversions
CREATE MASKING POLICY multi_cast_mask
WITH (value_text VARCHAR(50))
USING (CAST(CAST(value_text AS INTEGER) * 100 AS VARCHAR) || '%');

-- User-defined function references (assuming they exist)

-- 24. SQL UDF reference
CREATE MASKING POLICY udf_sql_mask
WITH (credit_card VARCHAR(256))
USING (redact_credit_card_sql(credit_card));

-- 25. Python UDF reference
CREATE MASKING POLICY udf_python_mask
WITH (personal_data TEXT)
USING (redact_personal_data_python(personal_data));

-- 26. Lambda UDF reference
CREATE MASKING POLICY udf_lambda_mask
WITH (sensitive_json SUPER)
USING (lambda_json_masker(sensitive_json));

-- Advanced expressions

-- 27. Function calls with multiple parameters
CREATE MASKING POLICY multi_param_function
WITH (start_date DATE, end_date DATE)
USING (DATEDIFF('day', start_date, end_date)::VARCHAR || ' days');

-- 28. Nested function calls
CREATE MASKING POLICY nested_functions
WITH (raw_text VARCHAR(500))
USING (UPPER(TRIM(SUBSTRING(LOWER(raw_text) FROM 1 FOR 10))));

-- 29. Concatenation with multiple strings
CREATE MASKING POLICY concat_mask
WITH (category VARCHAR(50), subcategory VARCHAR(50))
USING (category || '-' || subcategory || '-MASKED');

-- 30. Complex SUBSTRING with FROM and FOR
CREATE MASKING POLICY complex_substring
WITH (long_text TEXT)
USING (SUBSTRING(long_text FROM 5 FOR 20) || '...[TRUNCATED]');

-- Different data types

-- 31. BOOLEAN data type
CREATE MASKING POLICY boolean_mask
WITH (is_active BOOLEAN)
USING (FALSE);

-- 32. DATE data type with functions
CREATE MASKING POLICY date_mask
WITH (created_date DATE)
USING (DATE_TRUNC('month', created_date));

-- 33. TIMESTAMP data type
CREATE MASKING POLICY timestamp_mask
WITH (logged_at TIMESTAMP)
USING (DATE_TRUNC('hour', logged_at));

-- 34. DECIMAL with precision and scale
CREATE MASKING POLICY decimal_mask
WITH (precise_amount DECIMAL(12,4))
USING (0.0000::DECIMAL(12,4));

-- 35. BIGINT data type
CREATE MASKING POLICY bigint_mask
WITH (large_number BIGINT)
USING (large_number / 1000000 * 1000000);

-- 36. SUPER data type (Redshift-specific)
CREATE MASKING POLICY super_mask
WITH (json_data SUPER)
USING ('{"masked": true}'::SUPER);

-- IF NOT EXISTS examples

-- 37. CREATE with IF NOT EXISTS
CREATE MASKING POLICY IF NOT EXISTS conditional_create_mask
WITH (data VARCHAR(100))
USING ('CONDITIONALLY_MASKED');

-- 38. IF NOT EXISTS with complex expression
CREATE MASKING POLICY IF NOT EXISTS complex_conditional_mask
WITH (value INTEGER, status VARCHAR(20))
USING (
    CASE
        WHEN status = 'ACTIVE' THEN CAST(value AS VARCHAR)
        ELSE 'INACTIVE_MASKED'
    END
);

-- Multiple column types in single policy

-- 39. Mixed data types
CREATE MASKING POLICY mixed_types_mask
WITH (id INTEGER, name VARCHAR(100), amount DECIMAL(10,2), created_at TIMESTAMP)
USING (
    CASE
        WHEN amount > 1000 THEN name || ' (HIGH_VALUE)'
        WHEN amount > 100 THEN 'MEDIUM_VALUE_USER'
        ELSE 'LOW_VALUE_USER'
    END
);

-- 40. Text data types variations
CREATE MASKING POLICY text_variations_mask
WITH (char_field CHAR(10), varchar_field VARCHAR(255), text_field TEXT)
USING (
    CASE
        WHEN LENGTH(text_field) > 100 THEN 'LONG_TEXT'
        WHEN LENGTH(varchar_field) > 50 THEN 'MEDIUM_TEXT'
        ELSE 'SHORT_TEXT'
    END
);

-- Parentheses and operator precedence

-- 41. Complex parentheses grouping
CREATE MASKING POLICY parentheses_mask
WITH (a INTEGER, b INTEGER, c INTEGER)
USING ((a + b) * c / 100);

-- 42. Mixed operators with precedence
CREATE MASKING POLICY precedence_mask
WITH (base_value DECIMAL(10,2), multiplier DECIMAL(5,2), offset INTEGER)
USING (base_value * multiplier + offset::DECIMAL(10,2));

-- Edge cases

-- 43. Empty string masking
CREATE MASKING POLICY empty_string_mask
WITH (data VARCHAR(100))
USING ('');

-- 44. Single character masking
CREATE MASKING POLICY single_char_mask
WITH (code CHAR(1))
USING ('X');

-- 45. Very long identifier names
CREATE MASKING POLICY very_long_policy_name_for_testing_identifier_limits
WITH (very_long_column_name_for_testing VARCHAR(50))
USING ('LONG_NAME_MASKED');

-- Comments within statements

-- 46. Line comments in expression
CREATE MASKING POLICY commented_mask
WITH (value INTEGER)
USING (
    -- This masks by setting to zero
    0
);

-- 47. Block comments in expression
CREATE MASKING POLICY block_commented_mask
WITH (text_value TEXT)
USING (
    /* This is a multi-line
       block comment explaining
       the masking strategy */
    'BLOCK_COMMENTED_MASK'
);

-- Semicolon variations

-- 48. With semicolon
CREATE MASKING POLICY with_semicolon_mask
WITH (data TEXT)
USING ('SEMICOLON_TERMINATED');

-- 49. Without semicolon
CREATE MASKING POLICY without_semicolon_mask
WITH (data TEXT)
USING ('NO_SEMICOLON')

-- 50. Multiple statements
CREATE MASKING POLICY first_batch_policy
WITH (data1 TEXT)
USING ('FIRST');

CREATE MASKING POLICY second_batch_policy
WITH (data2 TEXT)
USING ('SECOND');

-- Quoted identifiers

-- 51. Policy name with spaces
CREATE MASKING POLICY "Policy with Spaces"
WITH (data VARCHAR(100))
USING ('SPACED_POLICY_MASKED');

-- 52. Column name with spaces
CREATE MASKING POLICY quoted_column_mask
WITH ("Column With Spaces" TEXT)
USING ('QUOTED_COLUMN_MASKED');

-- 53. Mixed quoted and unquoted
CREATE MASKING POLICY "Mixed Case Policy"
WITH (regular_column VARCHAR(50), "Special Column" TEXT)
USING ('MIXED_QUOTED_MASKED');