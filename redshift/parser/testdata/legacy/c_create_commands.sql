-- Test cases for CREATE commands starting with C

-- ====================================
-- CREATE EXTERNAL FUNCTION
-- ====================================

-- Basic external function
CREATE EXTERNAL FUNCTION public.f_distance (float, float, float, float)
RETURNS float
STABLE
LAMBDA 'lambda_distance'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftLambdaRole';

-- External function with no arguments
CREATE EXTERNAL FUNCTION get_current_time()
RETURNS timestamp
IMMUTABLE
LAMBDA 'lambda_get_time'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftLambdaRole';

-- External function with multiple data types
CREATE EXTERNAL FUNCTION process_json(varchar(65535), varchar(256))
RETURNS varchar(65535)
STABLE
LAMBDA 'lambda_json_processor'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftLambdaRole';

-- External function with OR REPLACE
CREATE OR REPLACE EXTERNAL FUNCTION calculate_score(integer, integer, integer)
RETURNS decimal(10,2)
STABLE
LAMBDA 'lambda_score_calculator'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftLambdaRole';

-- ====================================
-- CREATE EXTERNAL MODEL
-- ====================================

-- Basic external model from S3
CREATE EXTERNAL MODEL customer_churn_model
FROM 's3://mybucket/models/customer_churn.tar.gz'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftMLRole';

-- External model with function name
CREATE EXTERNAL MODEL sales_forecast_model
FROM 's3://mybucket/models/sales_forecast.tar.gz'
FUNCTION_NAME predict_sales
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftMLRole';

-- External model with settings
CREATE EXTERNAL MODEL fraud_detection_model
FROM 's3://mybucket/models/fraud_detection.tar.gz'
FUNCTION_NAME detect_fraud
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftMLRole'
SETTINGS (
    MAX_BATCH_SIZE '1000',
    MAX_PAYLOAD_IN_MB '5'
);

-- ====================================
-- CREATE EXTERNAL SCHEMA
-- ====================================

-- External schema from AWS Glue Data Catalog
CREATE EXTERNAL SCHEMA spectrum_schema
FROM DATA_CATALOG
DATABASE 'spectrum_db'
REGION 'us-east-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole'
CREATE EXTERNAL DATABASE IF NOT EXISTS;

-- External schema from Hive metastore
CREATE EXTERNAL SCHEMA hive_schema
FROM HIVE METASTORE
DATABASE 'hive_db'
URI 'ec2-123-456-789-012.compute-1.amazonaws.com'
PORT 9083
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole';

-- External schema from PostgreSQL
CREATE EXTERNAL SCHEMA postgres_schema
FROM POSTGRES
DATABASE 'mypostgresdb'
URI 'hostname.region.rds.amazonaws.com'
PORT 5432
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftRole'
SECRET_ARN 'arn:aws:secretsmanager:region:123456789012:secret:postgres-secret';

-- External schema from MySQL
CREATE EXTERNAL SCHEMA mysql_schema
FROM MYSQL
DATABASE 'mysqldb'
URI 'hostname.region.rds.amazonaws.com'
PORT 3306
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftRole'
SECRET_ARN 'arn:aws:secretsmanager:region:123456789012:secret:mysql-secret';

-- External schema from Redshift
CREATE EXTERNAL SCHEMA redshift_schema
FROM REDSHIFT
DATABASE 'remote_db'
REGION 'us-west-2'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftRole';

-- External schema with IF NOT EXISTS
CREATE EXTERNAL SCHEMA IF NOT EXISTS spectrum_schema2
FROM DATA_CATALOG
DATABASE 'spectrum_db2'
REGION 'us-east-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole';

-- External schema from Kinesis
CREATE EXTERNAL SCHEMA kinesis_schema
FROM KINESIS
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftKinesisRole';

-- ====================================
-- CREATE EXTERNAL VIEW
-- ====================================

-- Basic external view
CREATE EXTERNAL VIEW spectrum_view AS
SELECT col1, col2, col3
FROM spectrum_schema.external_table
WHERE col1 > 100;

-- External view with column list
CREATE EXTERNAL VIEW sales_summary_view (region, total_sales, avg_price) AS
SELECT region, SUM(sales_amount), AVG(unit_price)
FROM spectrum_schema.sales_data
GROUP BY region;

-- External view with OR REPLACE
CREATE OR REPLACE EXTERNAL VIEW customer_360_view AS
SELECT c.customer_id, c.name, o.order_count, p.total_spent
FROM spectrum_schema.customers c
JOIN spectrum_schema.order_summary o ON c.customer_id = o.customer_id
JOIN spectrum_schema.payment_summary p ON c.customer_id = p.customer_id;

-- ====================================
-- CREATE IDENTITY PROVIDER
-- ====================================

-- Azure identity provider
CREATE IDENTITY PROVIDER azure_idp
TYPE 'Azure'
NAMESPACE 'aad:myazuretenant.onmicrosoft.com';

-- Azure identity provider with provider URL
CREATE IDENTITY PROVIDER azure_idp_with_params
TYPE 'Azure'
PROVIDER_URL 'https://login.microsoftonline.com/tenant-id/v2.0';

-- Custom identity provider
CREATE IDENTITY PROVIDER custom_saml_idp
TYPE 'Custom'
PROVIDER_URL 'https://idp.mycompany.com';

-- ====================================
-- CREATE LIBRARY
-- ====================================

-- Basic Python library
CREATE LIBRARY pandas_lib
LANGUAGE PLPYTHONU
FROM 's3://mybucket/python-libs/pandas-layer.zip';

-- Library with credentials
CREATE LIBRARY numpy_lib
LANGUAGE PLPYTHONU
FROM 's3://mybucket/python-libs/numpy-layer.zip'
CREDENTIALS 'aws_access_key_id=AKIAIOSFODNN7EXAMPLE;aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY';

-- Library with region
CREATE LIBRARY scipy_lib
LANGUAGE PLPYTHONU
FROM 's3://mybucket/python-libs/scipy-layer.zip'
REGION AS 'us-west-2';

-- Library with OR REPLACE
CREATE OR REPLACE LIBRARY ml_utils_lib
LANGUAGE PLPYTHONU
FROM 's3://mybucket/python-libs/ml-utils.zip'
CREDENTIALS 'aws_access_key_id=AKIAIOSFODNN7EXAMPLE;aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'
REGION AS 'us-east-1';

-- ====================================
-- CREATE MASKING POLICY
-- ====================================

-- Basic masking policy for SSN
CREATE MASKING POLICY mask_ssn
WITH (ssn varchar(11))
USING (
    CASE
        WHEN current_user = 'admin' THEN ssn
        ELSE 'XXX-XX-' || SUBSTRING(ssn, 8, 4)
    END
);

-- Masking policy for credit card with role-based access
CREATE MASKING POLICY mask_credit_card
WITH (cc_number varchar(19))
USING (
    CASE
        WHEN current_user IN ('finance_admin', 'compliance_officer') THEN cc_number
        WHEN current_user_has_role('finance_read') THEN 'XXXX-XXXX-XXXX-' || SUBSTRING(cc_number, 16, 4)
        ELSE 'XXXX-XXXX-XXXX-XXXX'
    END
);

-- Masking policy with USING clause
CREATE MASKING POLICY mask_salary
WITH (salary decimal(10,2))
USING (
    CASE
        WHEN is_manager = true AND department = user_department THEN salary
        WHEN is_hr = true THEN salary
        ELSE 0
    END
);

-- Complex masking policy for email
CREATE MASKING POLICY mask_email
WITH (email varchar(255))
USING (
    CASE
        WHEN current_user = 'data_analyst' THEN
            SUBSTRING(email, 1, 3) || '***@' || SPLIT_PART(email, '@', 2)
        WHEN pg_has_role(current_user, 'marketing_team', 'MEMBER') THEN
            email
        ELSE
            'hidden@example.com'
    END
);
