-- ALTER GROUP test cases for Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_GROUP.html

-- Test 1: Basic ADD USER - add single user to a group
ALTER GROUP admin_group ADD USER dwuser;

-- Test 2: ADD USER with multiple users
ALTER GROUP admin_group ADD USER u1, u2;

-- Test 3: ADD USER with more users
ALTER GROUP sales_group ADD USER alice, bob, charlie;

-- Test 4: Basic DROP USER - remove single user from a group
ALTER GROUP admin_group DROP USER dwuser;

-- Test 5: DROP USER with multiple users
ALTER GROUP admin_group DROP USER u1, u2;

-- Test 6: DROP USER with more users
ALTER GROUP sales_group DROP USER alice, bob, charlie;

-- Test 7: Basic RENAME TO
ALTER GROUP admin_group RENAME TO administrators;

-- Test 8: RENAME TO with different name
ALTER GROUP sales_group RENAME TO sales_team;

-- Test 9: Group names with underscores
ALTER GROUP test_group ADD USER test_user;
ALTER GROUP test_group RENAME TO new_test_group;

-- Test 10: Mixed case group and user names
ALTER GROUP TestGroup ADD USER TestUser;
ALTER GROUP TestGroup DROP USER TestUser;
ALTER GROUP TestGroup RENAME TO NewTestGroup;

-- Test 11: Quoted identifiers
ALTER GROUP "admin group" ADD USER "dw user";
ALTER GROUP "admin group" DROP USER "dw user";
ALTER GROUP "admin group" RENAME TO "administrator group";

-- Test 12: Long user lists
ALTER GROUP developers ADD USER dev1, dev2, dev3, dev4, dev5, dev6, dev7, dev8, dev9, dev10;
ALTER GROUP developers DROP USER dev1, dev2, dev3, dev4, dev5, dev6, dev7, dev8, dev9, dev10;

-- Test 13: Various group name patterns
ALTER GROUP grp123 ADD USER usr456;
ALTER GROUP group_123_test ADD USER user_456_test;
ALTER GROUP grp$test ADD USER usr$test;

-- Test 14: User names with numbers and special characters
ALTER GROUP data_analysts ADD USER analyst1, analyst2, analyst_lead;
ALTER GROUP data_analysts ADD USER data$user, test_user123, user_456;

-- Test 15: Comments in statements
ALTER GROUP admin_group ADD USER dwuser; -- Adding user to admin group
ALTER GROUP admin_group DROP USER dwuser; -- Removing user from admin group
ALTER GROUP admin_group RENAME TO administrators; -- Renaming the group

-- Test 16: Multi-line formatting
ALTER GROUP admin_group
    ADD USER
        user1,
        user2,
        user3;

ALTER GROUP admin_group
    DROP USER
        user1,
        user2,
        user3;

ALTER GROUP admin_group
    RENAME TO
        new_admin_group;

-- Test 17: Reserved keywords as identifiers (quoted)
ALTER GROUP "group" ADD USER "user";
ALTER GROUP "alter" ADD USER "add";
ALTER GROUP "drop" RENAME TO "rename";

-- Test 18: Empty group operations (edge cases for testing)
ALTER GROUP empty_group ADD USER single_user;
ALTER GROUP empty_group DROP USER single_user;

-- Test 19: Sequential operations on same group
ALTER GROUP project_group ADD USER member1;
ALTER GROUP project_group ADD USER member2, member3;
ALTER GROUP project_group DROP USER member1;
ALTER GROUP project_group RENAME TO project_team;

-- Test 20: Groups with similar names
ALTER GROUP group1 ADD USER user1;
ALTER GROUP group2 ADD USER user2;
ALTER GROUP group11 ADD USER user11;
ALTER GROUP group_1 ADD USER user_1;