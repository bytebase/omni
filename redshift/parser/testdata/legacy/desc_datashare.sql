-- Test cases for DESC DATASHARE command
-- Based on AWS Redshift documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_DESC_DATASHARE.html

-- Basic DESC DATASHARE syntax
DESC DATASHARE salesshare;

-- DESCRIBE variant (alternative keyword)
DESCRIBE DATASHARE salesshare;

-- DESC DATASHARE with NAMESPACE (for inbound datashares within same AWS account)
DESC DATASHARE salesshare OF NAMESPACE '13b8833d-17c6-4f16-8fe4-1a018f5ed00d';

-- DESC DATASHARE with ACCOUNT and NAMESPACE (for inbound datashares across AWS accounts)
DESC DATASHARE salesshare OF ACCOUNT '123456789012' NAMESPACE '13b8833d-17c6-4f16-8fe4-1a018f5ed00d';

-- DESCRIBE variant with ACCOUNT and NAMESPACE
DESCRIBE DATASHARE salesshare OF ACCOUNT '123456789012' NAMESPACE '13b8833d-17c6-4f16-8fe4-1a018f5ed00d';

-- Different datashare names
DESC DATASHARE customer_data;
DESC DATASHARE sales_analytics;
DESC DATASHARE marketing_metrics;

-- Datashare names with underscores and numbers
DESC DATASHARE datashare_001;
DESC DATASHARE test_datashare_v2;

-- Case variations
desc datashare salesshare;
DESCRIBE datashare salesshare;
Desc DataShare salesshare;

-- With different account IDs
DESC DATASHARE inventory OF ACCOUNT '999888777666' NAMESPACE 'a1b2c3d4-e5f6-7890-abcd-ef1234567890';
DESC DATASHARE finance_data OF ACCOUNT '111222333444' NAMESPACE 'ffffffff-eeee-dddd-cccc-bbbbbbbbbbbb';

-- Only NAMESPACE specified (without ACCOUNT)
DESC DATASHARE regional_sales OF NAMESPACE '12345678-1234-5678-9012-123456789012';
DESCRIBE DATASHARE product_catalog OF NAMESPACE 'abcdef12-3456-7890-abcd-ef1234567890';

-- Whitespace variations
DESC    DATASHARE    salesshare;
DESC DATASHARE salesshare    OF    ACCOUNT    '123456789012'    NAMESPACE    '13b8833d-17c6-4f16-8fe4-1a018f5ed00d';

-- Comments
-- This describes an outbound datashare
DESC DATASHARE /* production datashare */ salesshare;

-- Multi-line example
DESC DATASHARE customer_analytics
OF ACCOUNT '555666777888'
NAMESPACE 'def456ef-789a-bcde-f012-3456789abcde';