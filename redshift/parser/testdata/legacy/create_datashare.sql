-- CREATE DATASHARE examples from AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_DATASHARE.html

-- Basic datashare creation
CREATE DATASHARE salesshare;

-- Create a datashare with public accessibility set to false (default)
CREATE DATASHARE private_share PUBLICACCESSIBLE FALSE;

-- Create a datashare with public accessibility set to true
CREATE DATASHARE public_share PUBLICACCESSIBLE TRUE;

-- Create a datashare with SET keyword before PUBLICACCESSIBLE
CREATE DATASHARE share_with_set SET PUBLICACCESSIBLE TRUE;

-- Create a datashare with equal sign in PUBLICACCESSIBLE
CREATE DATASHARE share_with_equal PUBLICACCESSIBLE = TRUE;

-- Create a datashare with SET and equal sign
CREATE DATASHARE share_set_equal SET PUBLICACCESSIBLE = FALSE;

-- Create a datashare managed by AWS Data Exchange
CREATE DATASHARE adx_share MANAGEDBY ADX;

-- Create a public datashare managed by AWS Data Exchange (from documentation example)
CREATE DATASHARE demoshare SET PUBLICACCESSIBLE TRUE, MANAGEDBY ADX;

-- Various combinations
CREATE DATASHARE combined_share1 PUBLICACCESSIBLE TRUE, MANAGEDBY ADX;
CREATE DATASHARE combined_share2 SET PUBLICACCESSIBLE = FALSE, MANAGEDBY ADX;
CREATE DATASHARE combined_share3 MANAGEDBY ADX, PUBLICACCESSIBLE TRUE;
CREATE DATASHARE combined_share4 MANAGEDBY ADX, SET PUBLICACCESSIBLE = TRUE;

-- Different identifier formats
CREATE DATASHARE "quoted_share";
CREATE DATASHARE schema_qualified.share_name;
CREATE DATASHARE _underscore_share;
CREATE DATASHARE share123;

-- Multiple options with different orderings
CREATE DATASHARE multi_option1 PUBLICACCESSIBLE FALSE, MANAGEDBY ADX;
CREATE DATASHARE multi_option2 MANAGEDBY ADX, PUBLICACCESSIBLE FALSE;

-- Edge cases
CREATE DATASHARE minimal_share ;
CREATE DATASHARE share_no_value PUBLICACCESSIBLE;
CREATE DATASHARE share_set_no_value SET PUBLICACCESSIBLE;