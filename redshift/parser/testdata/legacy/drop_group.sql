-- DROP GROUP examples from AWS Redshift documentation
-- https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_GROUP.html

-- Basic DROP GROUP
DROP GROUP guests;

-- Another example
DROP GROUP developers;

-- Drop a group that might have privileges (need to revoke first)
-- Note: The REVOKE statement must be executed before DROP GROUP if the group has privileges
DROP GROUP sales_team;

-- Drop groups with various naming conventions
DROP GROUP admin_users;
DROP GROUP readonly_group;
DROP GROUP data_analysts;