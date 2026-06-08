-- CREATE GROUP examples for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_GROUP.html

-- Basic CREATE GROUP without any users
CREATE GROUP sales_group;

-- CREATE GROUP with a single user
CREATE GROUP admin_group WITH USER admin1;

-- CREATE GROUP with multiple users (using WITH)
CREATE GROUP admin_group WITH USER admin1, admin2;

-- CREATE GROUP with multiple users (without WITH keyword - optional)
CREATE GROUP dev_group USER developer1, developer2, developer3;

-- CREATE GROUP with mixed case names
CREATE GROUP Marketing_Team WITH USER john_doe, jane_smith;

-- CREATE GROUP with underscore in group name
CREATE GROUP data_analytics_group WITH USER analyst1, analyst2, analyst3, analyst4;

-- CREATE GROUP for different departments
CREATE GROUP hr_department WITH USER hr_manager, hr_specialist;
CREATE GROUP finance_team USER accountant1, accountant2, treasurer;
CREATE GROUP it_operations WITH USER sysadmin, netadmin, dbadmin;

-- Note: Group names starting with two underscores are reserved for internal Redshift use
-- This would fail: CREATE GROUP __internal_group;