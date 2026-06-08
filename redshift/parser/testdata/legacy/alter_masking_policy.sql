-- ALTER MASKING POLICY test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_MASKING_POLICY.html

-- Basic ALTER MASKING POLICY statements

-- 1. Simple masking expression with a constant
ALTER MASKING POLICY mask_credit_card
USING ('XXXX-XXXX-XXXX-XXXX');

-- 2. Masking expression with explicit cast (required for constants)
ALTER MASKING POLICY mask_ssn
USING ('XXX-XX-XXXX'::VARCHAR(11));

-- 3. Masking expression using NULL
ALTER MASKING POLICY hide_sensitive_data
USING (NULL);

-- 4. Masking expression with column reference
ALTER MASKING POLICY partial_mask
USING (column_name);

-- Complex masking expressions

-- 5. Using built-in functions
ALTER MASKING POLICY mask_email
USING (SUBSTRING(email FROM 1 FOR 3) || '****@****');

-- 6. Using mathematical operations
ALTER MASKING POLICY mask_salary
USING (salary * 0);

-- 7. Using CASE expression for conditional masking
ALTER MASKING POLICY conditional_mask
USING (
    CASE
        WHEN user_level > 5 THEN actual_value
        ELSE 'REDACTED'
    END
);

-- 8. Multiple column masking with complex logic
ALTER MASKING POLICY complex_mask
USING (
    CASE
        WHEN role = 'admin' THEN original_value
        WHEN role = 'user' THEN SUBSTRING(original_value FROM 1 FOR 3) || '***'
        ELSE 'HIDDEN'
    END
);

-- 9. Using concatenation and string functions
ALTER MASKING POLICY mask_phone
USING (
    '(' || SUBSTRING(phone FROM 1 FOR 3) || ') XXX-XXXX'
);

-- 10. Numeric masking with arithmetic
ALTER MASKING POLICY mask_revenue
USING (
    CASE
        WHEN amount > 1000000 THEN 1000000
        WHEN amount > 100000 THEN amount / 10 * 10
        ELSE 0
    END
);

-- 11. Date masking expressions
ALTER MASKING POLICY mask_birthdate
USING (
    DATE_TRUNC('year', birthdate)
);

-- 12. Boolean masking
ALTER MASKING POLICY mask_flag
USING (
    CASE
        WHEN is_sensitive THEN FALSE
        ELSE is_sensitive
    END
);

-- User-defined function examples

-- 13. Using a SQL UDF
ALTER MASKING POLICY mask_with_udf
USING (
    my_masking_function(input_column)
);

-- 14. Using a Python UDF
ALTER MASKING POLICY mask_with_python
USING (
    python_mask_function(data_column, user_role)
);

-- 15. Using an AWS Lambda UDF
ALTER MASKING POLICY mask_with_lambda
USING (
    lambda_mask_handler(sensitive_data)
);

-- Edge cases and special scenarios

-- 16. Empty parentheses (would be invalid but testing parser)
-- ALTER MASKING POLICY invalid_mask USING ();

-- 17. Multiple nested functions
ALTER MASKING POLICY nested_functions
USING (
    UPPER(TRIM(SUBSTRING(LOWER(input_text) FROM 1 FOR 5)))
);

-- 18. Complex nested CASE expressions
ALTER MASKING POLICY nested_case
USING (
    CASE
        WHEN region = 'US' THEN
            CASE
                WHEN state = 'CA' THEN 'CALIFORNIA-MASKED'
                ELSE 'US-MASKED'
            END
        WHEN region = 'EU' THEN 'EU-GDPR-MASKED'
        ELSE 'INTERNATIONAL-MASKED'
    END
);

-- 19. Using IS NULL/IS NOT NULL
ALTER MASKING POLICY null_check_mask
USING (
    CASE
        WHEN input_value IS NULL THEN 'NO DATA'
        WHEN input_value IS NOT NULL THEN 'MASKED'
    END
);

-- 20. Logical operators
ALTER MASKING POLICY logic_mask
USING (
    CASE
        WHEN is_public AND NOT is_sensitive THEN original_value
        WHEN is_internal OR is_sensitive THEN 'RESTRICTED'
        ELSE 'DEFAULT-MASK'
    END
);

-- 21. Comparison operators
ALTER MASKING POLICY comparison_mask
USING (
    CASE
        WHEN age >= 18 AND age <= 65 THEN 'ADULT'
        WHEN age < 18 THEN 'MINOR'
        ELSE 'SENIOR'
    END
);

-- 22. Mixed data type casting
ALTER MASKING POLICY type_cast_mask
USING (
    CAST(numeric_value AS VARCHAR) || '-MASKED'
);

-- 23. Alternative cast syntax using ::
ALTER MASKING POLICY alt_cast_mask
USING (
    text_value::INTEGER * 0
);

-- Test with different identifier formats

-- 24. Quoted identifiers
ALTER MASKING POLICY "Mask Policy with Spaces"
USING ('MASKED');

-- 25. Mixed case identifiers
ALTER MASKING POLICY MixedCasePolicy
USING (DefaultMaskValue);

-- 26. Underscore in policy names
ALTER MASKING POLICY mask_policy_with_underscores
USING ('_MASKED_');

-- 27. Numbers in policy names
ALTER MASKING POLICY mask_policy_2024
USING ('MASKED_2024');

-- With and without semicolons

-- 28. With semicolon
ALTER MASKING POLICY with_semicolon
USING ('MASKED');

-- 29. Without semicolon
ALTER MASKING POLICY without_semicolon
USING ('MASKED')

-- 30. Multiple statements
ALTER MASKING POLICY first_policy USING ('FIRST');
ALTER MASKING POLICY second_policy USING ('SECOND');

-- Comments in statements

-- 31. With inline comments
ALTER MASKING POLICY commented_policy
USING (
    -- This masks sensitive data
    'COMMENTED_MASK'
);

-- 32. With block comments
ALTER MASKING POLICY block_comment_policy
USING (
    /* Multi-line
       comment block */
    'BLOCK_MASKED'
);