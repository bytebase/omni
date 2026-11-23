-- The access entities TestSyncInstance reads: a role with a setting, a user that
-- expires and takes the role by default, a user with two authentication methods
-- of which one never expires, grants on the database and on a column with the
-- grant option, a quota assigned by name and one assigned to everyone but one
-- user, a settings profile assigned to the role, and one assigned to nobody.
--
-- $database stands for the database the test created, which the grants name.
CREATE ROLE omni_reader SETTINGS max_threads = 2;

CREATE USER omni_analyst IDENTIFIED WITH sha256_password BY 'not stated' VALID UNTIL '2030-01-02 03:04:05' DEFAULT ROLE omni_reader;

CREATE USER omni_service IDENTIFIED WITH plaintext_password BY 'a' VALID UNTIL '2031-01-01', sha256_password BY 'b';

GRANT omni_reader TO omni_analyst;

GRANT SELECT ON $database.* TO omni_reader;

GRANT INSERT(id) ON $database.events TO omni_analyst WITH GRANT OPTION;

CREATE QUOTA omni_hourly FOR INTERVAL 1 hour MAX queries = 100 TO omni_analyst;

CREATE QUOTA omni_everyone FOR INTERVAL 1 day MAX errors = 10 TO ALL EXCEPT omni_analyst;

CREATE SETTINGS PROFILE omni_limited SETTINGS max_memory_usage = 1000000 MIN 1 MAX 2000000 TO omni_reader;

CREATE SETTINGS PROFILE omni_unassigned SETTINGS max_threads = 1;
