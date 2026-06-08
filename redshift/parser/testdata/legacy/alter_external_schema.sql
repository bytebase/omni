-- ALTER EXTERNAL SCHEMA examples for Amazon Redshift
-- Based on: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_EXTERNAL_SCHEMA.html

-- 1. Modify IAM role for external schema
-- Using default IAM role
ALTER EXTERNAL SCHEMA spectrum_schema
IAM_ROLE default;

-- Using SESSION IAM role
ALTER EXTERNAL SCHEMA spectrum_schema
IAM_ROLE 'SESSION';

-- Using specific IAM role ARN
ALTER EXTERNAL SCHEMA spectrum_schema
IAM_ROLE 'arn:aws:iam::012345678901:role/testrole';

-- 2. Implementing mTLS authentication with ACM certificate
ALTER EXTERNAL SCHEMA kafka_schema
AUTHENTICATION mtls
AUTHENTICATION_ARN 'arn:aws:acm:us-east-1:444455556666:certificate/certificate_ID';

-- 3. Modifying mTLS authentication with Secrets Manager
ALTER EXTERNAL SCHEMA kafka_schema
AUTHENTICATION mtls
SECRET_ARN 'arn:aws:secretsmanager:us-east-1:012345678910:secret:myMTLSSecret';

-- 4. Modifying Kafka cluster URI
ALTER EXTERNAL SCHEMA kafka_schema
URI 'lkc-ghidef-67890.centralus.azure.glb.confluent.cloud:9092';

-- 5. Setting authentication to none
ALTER EXTERNAL SCHEMA spectrum_schema
AUTHENTICATION none;

-- 6. Setting authentication to IAM
ALTER EXTERNAL SCHEMA spectrum_schema
AUTHENTICATION iam;

-- 7. Multiple clauses in single statement
ALTER EXTERNAL SCHEMA kafka_schema
AUTHENTICATION mtls
AUTHENTICATION_ARN 'arn:aws:acm:us-east-1:444455556666:certificate/new_certificate_ID'
URI 'new-kafka-cluster.region.provider.com:9092';

-- 8. Another example with multiple clauses
ALTER EXTERNAL SCHEMA spectrum_schema
IAM_ROLE 'arn:aws:iam::123456789012:role/MySpectrumRole'
AUTHENTICATION iam;

-- Note: To change the owner of an external schema, use ALTER SCHEMA command:
-- This works for both regular and external schemas in Redshift

-- Change owner to specific user
ALTER SCHEMA spectrum_schema OWNER TO dwuser;

-- Change owner to another user
ALTER SCHEMA kafka_schema OWNER TO analytics_user;

-- Note: Only the schema owner, a superuser, or a user with ALTER privilege
-- on the schema can execute ALTER SCHEMA commands