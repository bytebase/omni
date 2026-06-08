-- CREATE LIBRARY test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_LIBRARY.html

-- Basic CREATE LIBRARY with S3 source and IAM role
CREATE LIBRARY f_urlparse
LANGUAGE plpythonu
FROM 's3://amzn-s3-demo-bucket/urlparse3-1.0.3.zip'
CREDENTIALS 'aws_iam_role=arn:aws:iam::123456789012:role/RedshiftRole'
REGION AS 'us-east-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftRole';

-- CREATE LIBRARY with HTTP source (no credentials needed)
CREATE LIBRARY f_urlparse
LANGUAGE plpythonu
FROM 'https://example.com/packages/urlparse3-1.0.3.zip'
IAM_ROLE default;

-- CREATE OR REPLACE LIBRARY with S3 source
CREATE OR REPLACE LIBRARY my_python_lib
LANGUAGE plpythonu
FROM 's3://my-bucket/libraries/my-library-1.2.3.zip'
REGION 'us-west-2'
IAM_ROLE 'arn:aws:iam::987654321098:role/MyRedshiftRole';

-- CREATE LIBRARY with minimal syntax (HTTP source, default IAM role)
CREATE LIBRARY simple_lib
LANGUAGE plpythonu
FROM 'https://pypi.org/packages/simple-lib-1.0.zip'
IAM_ROLE default;

-- CREATE LIBRARY with quoted library name
CREATE LIBRARY "my-library-name"
LANGUAGE plpythonu
FROM 's3://data-bucket/python-libs/special-chars-lib.zip'
IAM_ROLE 'arn:aws:iam::111222333444:role/DataRole';

-- CREATE LIBRARY with complex S3 path
CREATE LIBRARY advanced_analytics
LANGUAGE plpythonu
FROM 's3://analytics-bucket/python-libraries/2023/advanced-analytics-v2.1.0.zip'
REGION AS 'eu-west-1'
IAM_ROLE 'arn:aws:iam::555666777888:role/AnalyticsRole';

-- CREATE LIBRARY with underscores in name
CREATE LIBRARY data_processing_utils
LANGUAGE plpythonu
FROM 's3://corp-bucket/libraries/data_processing_utils_v1.5.zip'
CREDENTIALS 'aws_iam_role=arn:aws:iam::999888777666:role/DataProcessingRole'
REGION 'ap-southeast-1'
IAM_ROLE 'arn:aws:iam::999888777666:role/DataProcessingRole';

-- CREATE LIBRARY with numbers in name
CREATE LIBRARY ml_lib_v2
LANGUAGE plpythonu
FROM 'https://github.com/example/releases/download/v2.0/ml_lib_v2.zip'
IAM_ROLE default;

-- CREATE OR REPLACE with HTTPS source
CREATE OR REPLACE LIBRARY json_utils
LANGUAGE plpythonu
FROM 'https://files.example.com/libraries/json-utilities-3.4.1.zip'
IAM_ROLE 'arn:aws:iam::123123123123:role/UtilsRole';

-- CREATE LIBRARY with deep S3 path structure
CREATE LIBRARY nested_lib
LANGUAGE plpythonu
FROM 's3://enterprise-data/departments/analytics/libraries/python/2024/Q1/nested-lib-1.0.0.zip'
REGION AS 'us-east-2'
IAM_ROLE 'arn:aws:iam::456456456456:role/EnterpriseRole';

-- CREATE LIBRARY with dollar sign in IAM role (testing edge cases)
CREATE LIBRARY test_lib
LANGUAGE plpythonu
FROM 's3://test-bucket/test-lib.zip'
IAM_ROLE 'arn:aws:iam::789789789789:role/Test$Role';

-- CREATE LIBRARY with case variations (should be case-insensitive)
create library lower_case_lib
language plpythonu
from 's3://bucket/lib.zip'
iam_role default;

CREATE LIBRARY UPPER_CASE_LIB
LANGUAGE PLPYTHONU
FROM 'S3://BUCKET/LIB.ZIP'
IAM_ROLE DEFAULT;

-- CREATE LIBRARY with mixed case
Create Library Mixed_Case_Lib
Language PlPythonU
From 's3://bucket/mixed-lib.zip'
Region as 'us-west-1'
Iam_Role 'arn:aws:iam::111111111111:role/MixedRole';