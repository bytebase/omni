-- SET SESSION AUTHORIZATION examples
-- https://docs.aws.amazon.com/redshift/latest/dg/r_SET_SESSION_AUTHORIZATION.html

-- Basic usage: Set session user to 'dwuser' (string literal)
SET SESSION AUTHORIZATION 'dwuser';

-- Set session user using identifier (without quotes)
SET SESSION AUTHORIZATION dwuser;

-- Set session user with mixed case identifier
SET SESSION AUTHORIZATION TestUser;

-- Set transaction-local user (only valid for current transaction)
SET LOCAL SESSION AUTHORIZATION 'dwuser';

-- Set transaction-local user using identifier
SET LOCAL SESSION AUTHORIZATION dwuser;

-- Reset session authorization to default user
SET SESSION AUTHORIZATION DEFAULT;

-- Reset transaction-local authorization to default
SET LOCAL SESSION AUTHORIZATION DEFAULT;

-- Various user name formats
SET SESSION AUTHORIZATION 'test_user';
SET SESSION AUTHORIZATION 'user123';
SET SESSION AUTHORIZATION 'admin@example.com';

-- Edge cases with special characters in user names
SET SESSION AUTHORIZATION 'user-with-dash';
SET SESSION AUTHORIZATION 'user.with.dots';
SET SESSION AUTHORIZATION 'user_with_underscore';

-- Combined with transaction blocks (typical usage pattern)
BEGIN;
SET LOCAL SESSION AUTHORIZATION 'unprivileged_user';
-- Test database operations as unprivileged_user
SELECT * FROM sensitive_table; -- This would be tested with reduced privileges
COMMIT; -- Session authorization reverts after transaction

-- Multiple authorization changes in one session
SET SESSION AUTHORIZATION 'analyst';
-- Do some work as analyst
SET SESSION AUTHORIZATION 'developer';
-- Do some work as developer
SET SESSION AUTHORIZATION DEFAULT;
-- Back to original user

-- Case sensitivity tests
SET SESSION AUTHORIZATION 'UserName';
SET SESSION AUTHORIZATION 'USERNAME';
SET SESSION AUTHORIZATION 'username';