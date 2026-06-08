-- CREATE MODEL Test Cases
-- Based on AWS Redshift documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_MODEL.html

-- Basic CREATE MODEL with table source
CREATE MODEL customer_churn_model
FROM customer_data
TARGET churn
FUNCTION predict_churn
IAM_ROLE default;

-- CREATE MODEL with SELECT statement
CREATE MODEL sales_forecast_model
FROM (
    SELECT product_id, sales_date, quantity, price
    FROM sales_history
    WHERE sales_date >= '2023-01-01'
)
TARGET quantity
FUNCTION forecast_sales
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftMLRole';

-- CREATE MODEL with job name
CREATE MODEL pretrained_model
FROM 'existing-training-job'
FUNCTION classify_sentiment
IAM_ROLE default;

-- CREATE MODEL with function parameters and return type
CREATE MODEL price_prediction_model
FROM product_features
TARGET price
FUNCTION predict_price(INTEGER, VARCHAR(50), DECIMAL(10,2))
RETURNS DECIMAL(10,2)
IAM_ROLE 'arn:aws:iam::123456789012:role/MLRole';

-- CREATE MODEL with SageMaker endpoint
CREATE MODEL external_model
FROM training_data
TARGET outcome
FUNCTION predict_outcome
SAGEMAKER 'my-endpoint'
IAM_ROLE default;

-- CREATE MODEL with SageMaker endpoint and model name
CREATE MODEL external_model_v2
FROM training_data
TARGET outcome
FUNCTION predict_outcome
SAGEMAKER 'my-endpoint':'model-v2'
IAM_ROLE default;

-- CREATE MODEL with AUTO OFF
CREATE MODEL manual_model
FROM user_behavior
TARGET conversion
FUNCTION predict_conversion
IAM_ROLE default
AUTO OFF;

-- CREATE MODEL with specific model type
CREATE MODEL xgboost_model
FROM training_dataset
TARGET label
FUNCTION classify_label
IAM_ROLE default
MODEL_TYPE XGBOOST;

-- CREATE MODEL with MLP model type
CREATE MODEL neural_network_model
FROM features_table
TARGET target_variable
FUNCTION nn_predict
IAM_ROLE default
MODEL_TYPE MLP;

-- CREATE MODEL with Linear Learner
CREATE MODEL linear_model
FROM regression_data
TARGET continuous_target
FUNCTION linear_predict
IAM_ROLE default
MODEL_TYPE LINEAR_LEARNER;

-- CREATE MODEL with K-Means clustering
CREATE MODEL clustering_model
FROM customer_segments
FUNCTION cluster_customers
IAM_ROLE default
MODEL_TYPE KMEANS;

-- CREATE MODEL with forecasting
CREATE MODEL time_series_model
FROM historical_data
TARGET sales_value
FUNCTION forecast_sales
IAM_ROLE default
MODEL_TYPE FORECAST;

-- CREATE MODEL with problem type - regression
CREATE MODEL regression_model
FROM numerical_data
TARGET continuous_value
FUNCTION predict_value
IAM_ROLE default
PROBLEM_TYPE (REGRESSION);

-- CREATE MODEL with problem type - binary classification
CREATE MODEL binary_classifier
FROM binary_features
TARGET binary_target
FUNCTION binary_predict
IAM_ROLE default
PROBLEM_TYPE (BINARY_CLASSIFICATION);

-- CREATE MODEL with problem type - multiclass classification
CREATE MODEL multiclass_classifier
FROM categorical_features
TARGET category
FUNCTION classify_category
IAM_ROLE default
PROBLEM_TYPE (MULTICLASS_CLASSIFICATION);

-- CREATE MODEL with objective
CREATE MODEL accuracy_model
FROM classification_data
TARGET class_label
FUNCTION classify
IAM_ROLE default
OBJECTIVE ('Accuracy');

-- CREATE MODEL with MSE objective
CREATE MODEL mse_model
FROM regression_features
TARGET numeric_target
FUNCTION regress
IAM_ROLE default
OBJECTIVE ('MSE');

-- CREATE MODEL with preprocessors
CREATE MODEL preprocessed_model
FROM raw_data
TARGET outcome
FUNCTION predict_processed
IAM_ROLE default
PREPROCESSORS 'normalize,encode_categorical';

-- CREATE MODEL with default hyperparameters
CREATE MODEL default_hp_model
FROM training_set
TARGET response
FUNCTION predict_response
IAM_ROLE default
HYPERPARAMETERS DEFAULT;

-- CREATE MODEL with custom hyperparameters
CREATE MODEL custom_hp_model
FROM feature_matrix
TARGET target_col
FUNCTION custom_predict
IAM_ROLE default
HYPERPARAMETERS DEFAULT EXCEPT (
    learning_rate '0.1',
    max_depth '6',
    subsample '0.8'
);

-- CREATE MODEL with basic settings
CREATE MODEL settings_model
FROM data_table
TARGET prediction_target
FUNCTION predict_with_settings
IAM_ROLE default
SETTINGS (
    S3_BUCKET 'my-ml-bucket',
    TAGS 'project=ml,team=data-science'
);

-- CREATE MODEL with comprehensive settings
CREATE MODEL full_settings_model
FROM comprehensive_data
TARGET full_target
FUNCTION comprehensive_predict
IAM_ROLE 'arn:aws:iam::123456789012:role/ComprehensiveMLRole'
SETTINGS (
    S3_BUCKET 'comprehensive-ml-bucket',
    TAGS 'environment=prod,cost-center=analytics',
    KMS_KEY_ID 'arn:aws:kms:us-west-2:123456789012:key/12345678-1234-1234-1234-123456789012',
    S3_GARBAGE_COLLECT ON,
    MAX_CELLS 1000000,
    MAX_RUNTIME 7200,
    HORIZON 30,
    FREQUENCY 24,
    PERCENTILES '0.1,0.5,0.9',
    MAX_BATCH_ROWS 10000
);

-- CREATE MODEL with all optional clauses
CREATE MODEL complete_model
FROM (
    SELECT feature1, feature2, feature3, target_variable
    FROM ml_training_data
    WHERE data_quality_score > 0.8
)
TARGET target_variable
FUNCTION complete_prediction(INTEGER, DECIMAL(10,2), VARCHAR(100))
RETURNS DECIMAL(8,4)
SAGEMAKER 'production-endpoint':'model-v3'
IAM_ROLE 'arn:aws:iam::123456789012:role/ProductionMLRole'
AUTO ON
MODEL_TYPE XGBOOST
PROBLEM_TYPE (BINARY_CLASSIFICATION)
OBJECTIVE ('F1')
PREPROCESSORS 'standard_scaler,one_hot_encoder'
HYPERPARAMETERS DEFAULT EXCEPT (
    n_estimators '100',
    learning_rate '0.05',
    max_depth '8',
    min_child_weight '3'
)
SETTINGS (
    S3_BUCKET 'production-ml-models',
    TAGS 'version=3.0,model=xgboost,status=production',
    KMS_KEY_ID 'arn:aws:kms:us-east-1:123456789012:key/production-key',
    S3_GARBAGE_COLLECT OFF,
    MAX_CELLS 5000000,
    MAX_RUNTIME 10800,
    MAX_BATCH_ROWS 50000
);

-- CREATE MODEL with quoted identifiers
CREATE MODEL "Model With Spaces"
FROM "Table With Spaces"
TARGET "Target Column"
FUNCTION "Prediction Function"
IAM_ROLE default;

-- CREATE MODEL with schema-qualified names
CREATE MODEL analytics.customer_model
FROM analytics.customer_features
TARGET analytics.churn_flag
FUNCTION analytics.predict_churn
IAM_ROLE default;

-- CREATE MODEL for time series forecasting with specific settings
CREATE MODEL forecast_model
FROM time_series_data
TARGET sales_amount
FUNCTION forecast_sales_ts
IAM_ROLE default
MODEL_TYPE FORECAST
SETTINGS (
    S3_BUCKET 'forecast-models',
    HORIZON 90,
    FREQUENCY 1440,
    PERCENTILES '0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9'
);

-- CREATE MODEL with complex SELECT statement
CREATE MODEL complex_select_model
FROM (
    SELECT
        c.customer_id,
        c.age,
        c.income,
        COUNT(o.order_id) as order_count,
        AVG(o.order_total) as avg_order_value,
        MAX(o.order_date) as last_order_date,
        CASE WHEN c.subscription_status = 'active' THEN 1 ELSE 0 END as is_subscriber
    FROM customers c
    LEFT JOIN orders o ON c.customer_id = o.customer_id
    WHERE c.created_date >= '2022-01-01'
    GROUP BY c.customer_id, c.age, c.income, c.subscription_status
    HAVING COUNT(o.order_id) > 0
)
TARGET is_subscriber
FUNCTION predict_subscription_likelihood
IAM_ROLE 'arn:aws:iam::123456789012:role/AnalyticsMLRole'
AUTO ON;