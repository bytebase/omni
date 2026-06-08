-- CREATE DATABASE examples for Amazon Redshift
-- Based on: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_DATABASE.html

-- Basic database creation
CREATE DATABASE mydb;

-- Database with owner
CREATE DATABASE tickit WITH OWNER dwuser;

-- Database with owner (using = sign)
CREATE DATABASE tickit2 WITH OWNER = dwuser;

-- Database with connection limit
CREATE DATABASE salesdb WITH CONNECTION LIMIT 100;

-- Database with unlimited connections
CREATE DATABASE devdb WITH CONNECTION LIMIT UNLIMITED;

-- Database with specific isolation level
CREATE DATABASE sampledb ISOLATION LEVEL SNAPSHOT;

-- Database with serializable isolation level
CREATE DATABASE transactionaldb ISOLATION LEVEL SERIALIZABLE;

-- Case-insensitive database
CREATE DATABASE cidb COLLATE CASE_INSENSITIVE;

-- Case-insensitive database using CI shorthand
CREATE DATABASE cidb2 COLLATE CI;

-- Case-sensitive database
CREATE DATABASE csdb COLLATE CASE_SENSITIVE;

-- Case-sensitive database using CS shorthand
CREATE DATABASE csdb2 COLLATE CS;

-- Database with multiple options
CREATE DATABASE proddb
    WITH OWNER produser
    CONNECTION LIMIT 50
    COLLATE CASE_INSENSITIVE
    ISOLATION LEVEL SNAPSHOT;

-- Database with all options
CREATE DATABASE enterprisedb
    OWNER = admin
    CONNECTION LIMIT 200
    COLLATE CI
    ISOLATION LEVEL SNAPSHOT;

-- Database from integration (for data sharing) - Complex syntax commented out for now
-- CREATE DATABASE shared_db
--     FROM INTEGRATION '12345678-1234-1234-1234-123456789012'
--     DATABASE 'source_database'
--     SET
--     ACCEPTINVCHARS TRUE
--     QUERY_ALL_STATES FALSE
--     REFRESH_INTERVAL 3600
--     TRUNCATECOLUMNS TRUE
--     HISTORY_MODE FALSE;

-- Database from integration with WITH clause options - Complex syntax commented out for now
-- CREATE DATABASE hybrid_db
--     FROM INTEGRATION 'abcdef12-3456-7890-abcd-ef1234567890'
--     DATABASE 'remote_db'
--     ACCEPTINVCHARS FALSE
--     WITH OWNER dataowner
--     CONNECTION LIMIT 75
--     ISOLATION LEVEL SNAPSHOT;

-- Database with semicolon
CREATE DATABASE testdb;

-- Database without semicolon (also valid)
CREATE DATABASE testdb2