-- CREATE EXTERNAL SCHEMA test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_SCHEMA.html

-- Basic Data Catalog external schema
CREATE EXTERNAL SCHEMA spectrum
FROM DATA CATALOG
DATABASE 'spectrumdb'
REGION 'us-west-2'
IAM_ROLE 'arn:aws:iam::123456789012:role/MySpectrumRole';

-- Data Catalog with IF NOT EXISTS
CREATE EXTERNAL SCHEMA IF NOT EXISTS spectrum_catalog
FROM DATA CATALOG
DATABASE 'my_glue_database'
REGION 'us-east-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole';

-- Data Catalog with DEFAULT IAM_ROLE
CREATE EXTERNAL SCHEMA default_role_schema
FROM DATA CATALOG
DATABASE 'analytics_db'
REGION 'eu-west-1'
IAM_ROLE DEFAULT;

-- Data Catalog with SESSION IAM_ROLE
CREATE EXTERNAL SCHEMA session_role_schema
FROM DATA CATALOG
DATABASE 'session_db'
REGION 'us-west-1'
IAM_ROLE 'SESSION';

-- Data Catalog with CATALOG_ROLE using SESSION
CREATE EXTERNAL SCHEMA catalog_session_schema
FROM DATA CATALOG
DATABASE 'shared_catalog_db'
REGION 'us-east-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole'
CATALOG_ROLE 'SESSION';

-- Data Catalog with CATALOG_ROLE using ARN
CREATE EXTERNAL SCHEMA catalog_arn_schema
FROM DATA CATALOG
DATABASE 'cross_account_db'
REGION 'us-west-2'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole'
CATALOG_ROLE 'arn:aws:iam::987654321098:role/CrossAccountCatalogRole';

-- Data Catalog with CREATE EXTERNAL DATABASE
CREATE EXTERNAL SCHEMA auto_create_db_schema
FROM DATA CATALOG
DATABASE 'new_external_db'
REGION 'us-east-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole'
CREATE EXTERNAL DATABASE IF NOT EXISTS;

-- Data Catalog with CATALOG_ID for cross-account
CREATE EXTERNAL SCHEMA cross_account_schema
FROM DATA CATALOG
DATABASE 'cross_account_database'
REGION 'us-west-2'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole'
CATALOG_ID '987654321098';

-- Data Catalog with all optional parameters
CREATE EXTERNAL SCHEMA comprehensive_catalog_schema
FROM DATA CATALOG
DATABASE 'comprehensive_db'
REGION 'eu-central-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole'
CATALOG_ROLE 'arn:aws:iam::123456789012:role/CatalogAccessRole'
CREATE EXTERNAL DATABASE IF NOT EXISTS
CATALOG_ID '123456789012';

-- Hive Metastore external schema
CREATE EXTERNAL SCHEMA hive_schema
FROM HIVE METASTORE
DATABASE 'hive_db'
URI 'thrift://myserver:9083'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole';

-- Hive Metastore with custom port
CREATE EXTERNAL SCHEMA hive_custom_port_schema
FROM HIVE METASTORE
DATABASE 'production_hive_db'
URI 'thrift://hive-metastore.company.com'
PORT 9084
IAM_ROLE 'arn:aws:iam::123456789012:role/HiveAccessRole';

-- PostgreSQL federated query schema
CREATE EXTERNAL SCHEMA postgres_federated
FROM POSTGRES
DATABASE 'postgres_production'
URI 'postgres-cluster.cluster-xyz.us-east-1.rds.amazonaws.com'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftFederatedRole'
SECRET_ARN 'arn:aws:secretsmanager:us-east-1:123456789012:secret:postgres-credentials-AbCdEf';

-- PostgreSQL with custom port and schema
CREATE EXTERNAL SCHEMA postgres_custom_schema
FROM POSTGRES
DATABASE 'analytics_db'
SCHEMA 'public'
URI 'postgres-server.company.com'
PORT 5433
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftFederatedRole'
SECRET_ARN 'arn:aws:secretsmanager:us-east-1:123456789012:secret:postgres-analytics-XyZaBc';

-- PostgreSQL with DEFAULT IAM_ROLE
CREATE EXTERNAL SCHEMA postgres_default_role
FROM POSTGRES
DATABASE 'reporting_db'
URI 'postgres.internal.company.com'
IAM_ROLE DEFAULT
SECRET_ARN 'arn:aws:secretsmanager:us-east-1:123456789012:secret:postgres-reporting-123456';

-- MySQL federated query schema
CREATE EXTERNAL SCHEMA mysql_federated
FROM MYSQL
DATABASE 'mysql_production'
URI 'mysql-cluster.cluster-abc.us-west-2.rds.amazonaws.com'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftFederatedRole'
SECRET_ARN 'arn:aws:secretsmanager:us-west-2:123456789012:secret:mysql-credentials-DefGhi';

-- MySQL with custom port
CREATE EXTERNAL SCHEMA mysql_custom_port
FROM MYSQL
DATABASE 'legacy_mysql_db'
URI 'mysql-legacy.company.com'
PORT 3307
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftFederatedRole'
SECRET_ARN 'arn:aws:secretsmanager:us-west-2:123456789012:secret:mysql-legacy-789012';

-- Kinesis streaming schema
CREATE EXTERNAL SCHEMA kinesis_stream_schema
FROM KINESIS
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftKinesisRole';

-- Kinesis with DEFAULT IAM_ROLE
CREATE EXTERNAL SCHEMA kinesis_default_role
FROM KINESIS
IAM_ROLE DEFAULT;

-- Kafka streaming schema with URI
CREATE EXTERNAL SCHEMA kafka_stream_schema
FROM KAFKA
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftKafkaRole'
URI 'kafka-broker1:9092,kafka-broker2:9092,kafka-broker3:9092';

-- Kafka with authentication none
CREATE EXTERNAL SCHEMA kafka_no_auth
FROM KAFKA
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftKafkaRole'
URI 'kafka-cluster.company.com:9092'
AUTHENTICATION NONE;

-- Kafka with IAM authentication
CREATE EXTERNAL SCHEMA kafka_iam_auth
FROM KAFKA
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftKafkaRole'
URI 'kafka-secure.company.com:9093'
AUTHENTICATION IAM;

-- Kafka with mTLS authentication
CREATE EXTERNAL SCHEMA kafka_mtls_auth
FROM KAFKA
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftKafkaRole'
URI 'kafka-mtls.company.com:9094'
AUTHENTICATION MTLS
AUTHENTICATION_ARN 'arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012';

-- MSK (Managed Streaming for Kafka) schema
CREATE EXTERNAL SCHEMA msk_stream_schema
FROM MSK
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftMSKRole'
URI 'b-1.msk-cluster.abc123.c2.kafka.us-east-1.amazonaws.com:9092';

-- MSK with IAM authentication
CREATE EXTERNAL SCHEMA msk_iam_auth
FROM MSK
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftMSKRole'
URI 'b-1.secure-msk.def456.c2.kafka.us-west-2.amazonaws.com:9098'
AUTHENTICATION IAM;

-- MSK with mTLS authentication
CREATE EXTERNAL SCHEMA msk_mtls_auth
FROM MSK
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftMSKRole'
URI 'b-1.mtls-msk.ghi789.c2.kafka.eu-west-1.amazonaws.com:9094'
AUTHENTICATION MTLS
AUTHENTICATION_ARN 'arn:aws:acm:eu-west-1:123456789012:certificate/87654321-4321-4321-4321-210987654321';

-- Cross-database query schema (Redshift to Redshift)
CREATE EXTERNAL SCHEMA cross_db_schema
FROM REDSHIFT
DATABASE 'prod_analytics'
SCHEMA 'reporting';

-- Cross-database without explicit schema
CREATE EXTERNAL SCHEMA cross_db_no_schema
FROM REDSHIFT
DATABASE 'data_warehouse';

-- Data Catalog without explicit "DATA CATALOG" keywords
CREATE EXTERNAL SCHEMA implicit_catalog
FROM
DATABASE 'implicit_db'
REGION 'ap-southeast-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftSpectrumRole';

-- Complex schema names with underscores
CREATE EXTERNAL SCHEMA data_lake_analytics_schema
FROM DATA CATALOG
DATABASE 'data_lake_db'
REGION 'us-east-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/DataLakeAccessRole';

-- Schema for IoT data from Kinesis
CREATE EXTERNAL SCHEMA iot_telemetry_schema
FROM KINESIS
IAM_ROLE 'arn:aws:iam::123456789012:role/IoTDataProcessingRole';

-- Schema for real-time analytics with Kafka
CREATE EXTERNAL SCHEMA realtime_analytics
FROM KAFKA
IAM_ROLE 'arn:aws:iam::123456789012:role/RealtimeAnalyticsRole'
URI 'kafka1.analytics.com:9092,kafka2.analytics.com:9092'
AUTHENTICATION IAM;

-- Multi-tenant schema with catalog role
CREATE EXTERNAL SCHEMA tenant_a_schema
FROM DATA CATALOG
DATABASE 'tenant_a_data'
REGION 'us-west-2'
IAM_ROLE 'arn:aws:iam::123456789012:role/MultiTenantRole'
CATALOG_ROLE 'arn:aws:iam::111111111111:role/TenantACatalogRole';

-- Development environment schema
CREATE EXTERNAL SCHEMA IF NOT EXISTS dev_spectrum
FROM DATA CATALOG
DATABASE 'development_db'
REGION 'us-east-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/DevEnvironmentRole'
CREATE EXTERNAL DATABASE IF NOT EXISTS;

-- Production PostgreSQL mirror
CREATE EXTERNAL SCHEMA prod_postgres_mirror
FROM POSTGRES
DATABASE 'production_mirror'
SCHEMA 'analytics'
URI 'readonly-postgres.prod.company.com'
PORT 5432
IAM_ROLE 'arn:aws:iam::123456789012:role/ProductionReadOnlyRole'
SECRET_ARN 'arn:aws:secretsmanager:us-east-1:123456789012:secret:prod-postgres-readonly-AbCdEf';

-- Legacy MySQL integration
CREATE EXTERNAL SCHEMA legacy_mysql_integration
FROM MYSQL
DATABASE 'legacy_erp_system'
URI 'mysql-legacy.internal.company.com'
PORT 3306
IAM_ROLE 'arn:aws:iam::123456789012:role/LegacySystemIntegrationRole'
SECRET_ARN 'arn:aws:secretsmanager:us-east-1:123456789012:secret:legacy-mysql-credentials-XyZ123';

-- Hive data lake integration
CREATE EXTERNAL SCHEMA hive_data_lake
FROM HIVE METASTORE
DATABASE 'enterprise_data_lake'
URI 'thrift://hive-metastore.datalake.company.com:9083'
IAM_ROLE 'arn:aws:iam::123456789012:role/DataLakeHiveRole';

-- Cross-region data access
CREATE EXTERNAL SCHEMA cross_region_data
FROM DATA CATALOG
DATABASE 'cross_region_analytics'
REGION 'eu-central-1'
IAM_ROLE 'arn:aws:iam::123456789012:role/CrossRegionAccessRole'
CATALOG_ID '987654321098';

-- Streaming schema for log analysis
CREATE EXTERNAL SCHEMA log_analysis_stream
FROM KINESIS
IAM_ROLE 'arn:aws:iam::123456789012:role/LogAnalysisRole';

-- High-security Kafka with mTLS
CREATE EXTERNAL SCHEMA secure_kafka_feed
FROM KAFKA
IAM_ROLE 'arn:aws:iam::123456789012:role/SecureDataFeedRole'
URI 'secure-kafka.compliance.company.com:9094'
AUTHENTICATION MTLS
AUTHENTICATION_ARN 'arn:aws:acm:us-east-1:123456789012:certificate/security-cert-123456';

-- MSK for financial data
CREATE EXTERNAL SCHEMA financial_data_msk
FROM MSK
IAM_ROLE 'arn:aws:iam::123456789012:role/FinancialDataRole'
URI 'b-1.financial-msk.xyz789.c2.kafka.us-east-1.amazonaws.com:9098'
AUTHENTICATION IAM;

-- Simple cross-database query for reporting
CREATE EXTERNAL SCHEMA reporting_cross_db
FROM REDSHIFT
DATABASE 'reporting_warehouse'
SCHEMA 'public';