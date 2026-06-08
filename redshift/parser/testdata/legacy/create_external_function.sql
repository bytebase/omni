-- CREATE EXTERNAL FUNCTION test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_FUNCTION.html

-- Basic CREATE EXTERNAL FUNCTION examples
CREATE EXTERNAL FUNCTION exfunc_sum(INT, INT)
RETURNS INT
VOLATILE
LAMBDA 'lambda_sum'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with VARCHAR parameters and return type
CREATE EXTERNAL FUNCTION exfunc_upper(VARCHAR)
RETURNS VARCHAR
STABLE
LAMBDA 'lambda_upper'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with multiple parameters of different types
CREATE EXTERNAL FUNCTION exfunc_multiplication(INT, INT, INT)
RETURNS INT
VOLATILE
LAMBDA 'lambda_multiplication'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with RETRY_TIMEOUT option
CREATE EXTERNAL FUNCTION exfunc_with_timeout(VARCHAR(100))
RETURNS VARCHAR(200)
VOLATILE
LAMBDA 'lambda_with_timeout'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test'
RETRY_TIMEOUT 30000;

-- Function with MAX_BATCH_ROWS option
CREATE EXTERNAL FUNCTION exfunc_batch_processing(INT, VARCHAR(50))
RETURNS VARCHAR(100)
STABLE
LAMBDA 'lambda_batch_processor'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test'
MAX_BATCH_ROWS 1000;

-- Function with MAX_BATCH_SIZE option in KB
CREATE EXTERNAL FUNCTION exfunc_with_batch_size_kb(TEXT)
RETURNS TEXT
VOLATILE
LAMBDA 'lambda_text_processor'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test'
MAX_BATCH_SIZE 512 KB;

-- Function with MAX_BATCH_SIZE option in MB
CREATE EXTERNAL FUNCTION exfunc_with_batch_size_mb(SUPER)
RETURNS SUPER
STABLE
LAMBDA 'lambda_super_processor'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test'
MAX_BATCH_SIZE 10 MB;

-- Function with all optional parameters
CREATE EXTERNAL FUNCTION exfunc_complete(VARCHAR(255), INT, DECIMAL(10,2))
RETURNS JSON
VOLATILE
LAMBDA 'lambda_complete_function'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test'
RETRY_TIMEOUT 25000
MAX_BATCH_ROWS 500
MAX_BATCH_SIZE 2 MB;

-- CREATE OR REPLACE EXTERNAL FUNCTION
CREATE OR REPLACE EXTERNAL FUNCTION exfunc_replaceable(BIGINT)
RETURNS BOOLEAN
STABLE
LAMBDA 'lambda_replaceable'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with DEFAULT IAM_ROLE
CREATE EXTERNAL FUNCTION exfunc_default_role(SMALLINT, SMALLINT)
RETURNS REAL
VOLATILE
LAMBDA 'lambda_default_role'
IAM_ROLE DEFAULT;

-- Function with timestamp parameters
CREATE EXTERNAL FUNCTION exfunc_timestamp_processing(TIMESTAMP, DATE)
RETURNS TIMESTAMPTZ
STABLE
LAMBDA 'lambda_timestamp_processor'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with time zone types
CREATE EXTERNAL FUNCTION exfunc_timezone(TIMESTAMP WITH TIME ZONE)
RETURNS TIME WITHOUT TIME ZONE
VOLATILE
LAMBDA 'lambda_timezone_converter'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with JSONB parameters
CREATE EXTERNAL FUNCTION exfunc_json_processor(JSONB, VARCHAR(100))
RETURNS JSONB
STABLE
LAMBDA 'lambda_json_processor'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with geometric types
CREATE EXTERNAL FUNCTION exfunc_geometry(GEOMETRY)
RETURNS GEOGRAPHY
VOLATILE
LAMBDA 'lambda_geo_converter'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with HLLSKETCH type
CREATE EXTERNAL FUNCTION exfunc_hll_analysis(HLLSKETCH)
RETURNS BIGINT
STABLE
LAMBDA 'lambda_hll_estimator'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with INTERVAL type
CREATE EXTERNAL FUNCTION exfunc_interval_calc(INTERVAL, INT)
RETURNS INTERVAL
VOLATILE
LAMBDA 'lambda_interval_calculator'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with DOUBLE PRECISION
CREATE EXTERNAL FUNCTION exfunc_precision_calc(DOUBLE PRECISION, NUMERIC(15,5))
RETURNS DOUBLE PRECISION
STABLE
LAMBDA 'lambda_precision_calculator'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with no parameters
CREATE EXTERNAL FUNCTION exfunc_no_params()
RETURNS VARCHAR(50)
VOLATILE
LAMBDA 'lambda_no_params'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Complex function with many parameters
CREATE EXTERNAL FUNCTION exfunc_complex(
    VARCHAR(100),
    INT,
    DECIMAL(18,6),
    TIMESTAMP,
    BOOLEAN,
    TEXT,
    JSONB
)
RETURNS SUPER
VOLATILE
LAMBDA 'lambda_complex_processor'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test'
RETRY_TIMEOUT 60000
MAX_BATCH_ROWS 100
MAX_BATCH_SIZE 5 MB;

-- Function for data validation
CREATE EXTERNAL FUNCTION exfunc_validate_email(VARCHAR(255))
RETURNS BOOLEAN
STABLE
LAMBDA 'lambda_email_validator'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function for data transformation
CREATE EXTERNAL FUNCTION exfunc_transform_data(SUPER)
RETURNS SUPER
VOLATILE
LAMBDA 'lambda_data_transformer'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test'
MAX_BATCH_SIZE 1024 KB;

-- Function with TIMETZ type
CREATE EXTERNAL FUNCTION exfunc_time_processing(TIMETZ)
RETURNS TIME WITH TIME ZONE
STABLE
LAMBDA 'lambda_time_processor'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function for machine learning predictions
CREATE OR REPLACE EXTERNAL FUNCTION exfunc_ml_predict(
    REAL,
    REAL,
    REAL,
    REAL
)
RETURNS REAL
STABLE
LAMBDA 'lambda_ml_model'
IAM_ROLE 'arn:aws:iam::123456789012:role/ML-Redshift-Role'
RETRY_TIMEOUT 45000
MAX_BATCH_ROWS 10000;

-- Function for text analysis
CREATE EXTERNAL FUNCTION exfunc_sentiment_analysis(TEXT)
RETURNS VARCHAR(20)
STABLE
LAMBDA 'lambda_sentiment_analyzer'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-NLP-Role'
MAX_BATCH_SIZE 256 KB;

-- Function with CHAR type with length
CREATE EXTERNAL FUNCTION exfunc_char_processor(CHAR(10))
RETURNS CHAR(20)
VOLATILE
LAMBDA 'lambda_char_processor'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function with various numeric types
CREATE EXTERNAL FUNCTION exfunc_numeric_operations(
    NUMERIC(10,2),
    DECIMAL(15,4),
    INTEGER,
    BIGINT,
    SMALLINT
)
RETURNS NUMERIC(20,6)
STABLE
LAMBDA 'lambda_numeric_calculator'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test';

-- Function for API integration
CREATE EXTERNAL FUNCTION exfunc_api_call(VARCHAR(500), VARCHAR(100))
RETURNS JSONB
VOLATILE
LAMBDA 'lambda_api_gateway'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-API-Role'
RETRY_TIMEOUT 30000
MAX_BATCH_ROWS 50;

-- Function with maximum batch size without unit (defaults to bytes)
CREATE EXTERNAL FUNCTION exfunc_default_size_unit(VARCHAR(1000))
RETURNS TEXT
STABLE
LAMBDA 'lambda_text_analyzer'
IAM_ROLE 'arn:aws:iam::123456789012:role/Redshift-Exfunc-Test'
MAX_BATCH_SIZE 1048576;