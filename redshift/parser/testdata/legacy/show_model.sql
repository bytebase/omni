-- SHOW MODEL test cases for Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_MODEL.html

-- Test Case 1: Show all models
-- This command shows all models the user has access to
SHOW MODEL ALL;

-- Test Case 2: Show specific model details
-- This displays detailed information about the customer_churn model
SHOW MODEL customer_churn;

-- Test Case 3: Show model with schema-qualified name
-- Testing with a fully qualified model name
SHOW MODEL public.customer_churn;

-- Test Case 4: Show model with quoted identifier
-- Testing with a quoted model name that may contain special characters
SHOW MODEL "my-model-2023";

-- Test Case 5: Show model with mixed case identifier
-- Testing with a mixed case model name (case-insensitive in Redshift)
SHOW MODEL CustomerSegmentation;

-- Test Case 6: Show model with underscore in name
-- Common naming convention for ML models
SHOW MODEL sales_forecast_model;

-- Test Case 7: Show model with numeric suffix
-- Testing model names with version numbers
SHOW MODEL churn_model_v2;

-- Test Case 8: Show all models (lowercase)
-- Testing case insensitivity of the ALL keyword
show model all;

-- Test Case 9: Show specific model (mixed case command)
-- Testing case insensitivity of the SHOW MODEL command
Show Model product_recommendation;

-- Test Case 10: Show model with long name
-- Testing longer model names
SHOW MODEL customer_lifetime_value_prediction_model_2023;