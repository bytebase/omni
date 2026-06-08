-- Amazon Redshift DROP USER Command Examples
-- Based on AWS documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_USER.html

-- Basic DROP USER - single user
DROP USER paulo;

-- DROP USER with multiple users
DROP USER paulo, martha;

-- DROP USER with IF EXISTS clause (single user)
DROP USER IF EXISTS paulo;

-- DROP USER with IF EXISTS clause (multiple users)
DROP USER IF EXISTS paulo, martha, john;

-- DROP USER with quoted identifiers
DROP USER "user with spaces";
DROP USER "User_With_Mixed_Case";

-- DROP USER combining regular and quoted identifiers
DROP USER normaluser, "quoted user", another_user;

-- DROP USER with IF EXISTS and quoted identifiers
DROP USER IF EXISTS "test_user", regular_user;

-- DROP USER with various valid identifier formats
DROP USER user1;
DROP USER user_123;
DROP USER user$test;
DROP USER _leading_underscore;

-- DROP USER with case variations (should be case-insensitive for keywords)
drop user test_user1;
Drop User test_user2;
DROP user test_user3;
drop USER test_user4;

-- Complex example with IF EXISTS and multiple mixed identifiers
DROP USER IF EXISTS
    regular_user1,
    "quoted user 2",
    user_with_numbers123,
    "CaseSensitiveUser",
    _underscore_user;

-- Single line with multiple users
DROP USER user1, user2, user3, user4, user5;