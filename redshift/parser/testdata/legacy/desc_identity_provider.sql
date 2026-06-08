-- DESC IDENTITY PROVIDER command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DESC_IDENTITY_PROVIDER.html

-- Basic DESC IDENTITY PROVIDER
DESC IDENTITY PROVIDER azure_idp;

-- Alternative syntax with DESCRIBE
DESCRIBE IDENTITY PROVIDER azure_idp;

-- Identity provider with different naming patterns
DESC IDENTITY PROVIDER my_identity_provider;
DESC IDENTITY PROVIDER okta_provider;
DESC IDENTITY PROVIDER saml_idp_001;

-- Quoted identifier names
DESC IDENTITY PROVIDER "Azure-IDP";
DESC IDENTITY PROVIDER "My Identity Provider";
DESC IDENTITY PROVIDER "okta_provider_2023";

-- Case variations (should work due to case-insensitive keywords)
desc identity provider test_idp;
Desc Identity Provider test_idp;
DESCRIBE identity provider test_idp;

-- Names with numbers and underscores
DESC IDENTITY PROVIDER idp_123;
DESC IDENTITY PROVIDER provider_v2;
DESC IDENTITY PROVIDER auth0_provider_prod;

-- Long identifier names
DESC IDENTITY PROVIDER very_long_identity_provider_name_for_testing;
DESCRIBE IDENTITY PROVIDER corporate_sso_identity_provider_production;

-- Names that might contain SQL keywords (quoted)
DESC IDENTITY PROVIDER "select";
DESC IDENTITY PROVIDER "table";
DESC IDENTITY PROVIDER "identity";

-- Realistic identity provider names from different vendors
DESC IDENTITY PROVIDER azure_ad_corp;
DESC IDENTITY PROVIDER okta_enterprise;
DESC IDENTITY PROVIDER auth0_production;
DESC IDENTITY PROVIDER ping_identity;
DESC IDENTITY PROVIDER onelogin_sso;
DESC IDENTITY PROVIDER google_workspace;
DESC IDENTITY PROVIDER microsoft_entra;