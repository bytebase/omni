-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestSyncInstance default#1
> exec
| -- The schema the sync tests read. Nothing is inserted: a table's row count and
| -- data size are the server's approximations, and an empty table makes both
| -- zero in every recording.
| CREATE TABLE events (
|   id UInt64,
|   name String COMMENT 'the event name',
|   at DateTime64(3, 'UTC'),
|   amount Nullable(Decimal(12, 2)),
|   tags Array(String),
|   attrs Map(String, String),
|   fixed FixedString(4),
|   flag Bool,
|   level Enum8('info' = 1, 'warn' = 2),
|   seen Date DEFAULT '2026-01-02'
| )
| ENGINE = MergeTree
| ORDER BY (id, name)
| PRIMARY KEY id
| COMMENT 'events the sync reads';
< ok 0 -
> exec
|
|
| CREATE TABLE skipped (
|   url String,
|   body String,
|   INDEX url_bloom url TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| ORDER BY url;
< ok 0 -
> exec
|
|
| CREATE VIEW event_names AS SELECT id, name FROM events;
< ok 0 -
> exec
| -- The access entities TestSyncInstance reads: a role with a setting, a user that
| -- expires and takes the role by default, a user with two authentication methods
| -- of which one never expires, grants on the database and on a column with the
| -- grant option, a quota assigned by name and one assigned to everyone but one
| -- user, a settings profile assigned to the role, and one assigned to nobody.
| --
| -- omni_testsyncinstance stands for the database the test created, which the grants name.
| CREATE ROLE omni_reader SETTINGS max_threads = 2;
< ok 0 -
> exec
|
|
| CREATE USER omni_analyst IDENTIFIED WITH sha256_password BY 'not stated' VALID UNTIL '2030-01-02 03:04:05' DEFAULT ROLE omni_reader;
< ok 0 -
> exec
|
|
| CREATE USER omni_service IDENTIFIED WITH plaintext_password BY 'a' VALID UNTIL '2031-01-01', sha256_password BY 'b';
< ok 0 -
> exec
|
|
| GRANT omni_reader TO omni_analyst;
< ok 0 -
> exec
|
|
| GRANT SELECT ON omni_testsyncinstance.* TO omni_reader;
< ok 0 -
> exec
|
|
| GRANT INSERT(id) ON omni_testsyncinstance.events TO omni_analyst WITH GRANT OPTION;
< ok 0 -
> exec
|
|
| CREATE QUOTA omni_hourly FOR INTERVAL 1 hour MAX queries = 100 TO omni_analyst;
< ok 0 -
> exec
|
|
| CREATE QUOTA omni_everyone FOR INTERVAL 1 day MAX errors = 10 TO ALL EXCEPT omni_analyst;
< ok 0 -
> exec
|
|
| CREATE SETTINGS PROFILE omni_limited SETTINGS max_memory_usage = 1000000 MIN 1 MAX 2000000 TO omni_reader;
< ok 0 -
> exec
|
|
| CREATE SETTINGS PROFILE omni_unassigned SETTINGS max_threads = 1;
< ok 0 -
> query
| SELECT VERSION()
< result
< column "VERSION()" "String" "string"
< row
: string
| 26.8.38.9
> query
|
| 		SELECT
| 			name,
| 			if(empty(valid_until) OR arrayExists(t -> toUnixTimestamp(t) = 0, valid_until), '', toString(arrayMax(valid_until)))
| 		FROM system.users
| 		ORDER BY name
< result
< column "name" "String" "string"
< column "if(or(empty(valid_until), arrayExists(lambda(tuple(t), equals(toUnixTimestamp(t), 0)), valid_until)), '', toString(arrayMax(valid_until)))" "String" "string"
< row
: string
| default
: string
|
< row
: string
| omni_analyst
: string
| 2030-01-02 03:04:05
< row
: string
| omni_service
: string
|
> query
| SELECT name FROM system.roles ORDER BY name
< result
< column "name" "String" "string"
< row
: string
| omni_reader
> query
|
| 		SELECT
| 			grantee.name,
| 			entity.name
| 		FROM system.quotas AS entity
| 		CROSS JOIN (SELECT name FROM system.users UNION DISTINCT SELECT name FROM system.roles) AS grantee
| 		WHERE has(entity.apply_to_list, grantee.name)
| 			OR (entity.apply_to_all AND NOT has(entity.apply_to_except, grantee.name))
| 		ORDER BY entity.name, grantee.name
< result
< column "grantee.name" "String" "string"
< column "name" "String" "string"
< row
: string
| default
: string
| default
< row
: string
| default
: string
| omni_everyone
< row
: string
| omni_reader
: string
| omni_everyone
< row
: string
| omni_service
: string
| omni_everyone
< row
: string
| omni_analyst
: string
| omni_hourly
> query
| SHOW CREATE QUOTA `default`
< result
< column "CREATE QUOTA default" "String" "string"
< row
: string
| CREATE QUOTA default KEYED BY user_name FOR INTERVAL 1 hour TRACKING ONLY TO default
> query
| SHOW CREATE QUOTA `omni_everyone`
< result
< column "CREATE QUOTA omni_everyone" "String" "string"
< row
: string
| CREATE QUOTA omni_everyone FOR INTERVAL 1 day MAX errors = 10 TO ALL EXCEPT omni_analyst
> query
| SHOW CREATE QUOTA `omni_hourly`
< result
< column "CREATE QUOTA omni_hourly" "String" "string"
< row
: string
| CREATE QUOTA omni_hourly FOR INTERVAL 1 hour MAX queries = 100 TO omni_analyst
> query
|
| 		SELECT
| 			grantee.name,
| 			entity.name
| 		FROM system.settings_profiles AS entity
| 		CROSS JOIN (SELECT name FROM system.users UNION DISTINCT SELECT name FROM system.roles) AS grantee
| 		WHERE has(entity.apply_to_list, grantee.name)
| 			OR (entity.apply_to_all AND NOT has(entity.apply_to_except, grantee.name))
| 		ORDER BY entity.name, grantee.name
< result
< column "grantee.name" "String" "string"
< column "name" "String" "string"
< row
: string
| omni_reader
: string
| omni_limited
> query
| SHOW CREATE SETTINGS PROFILE `omni_limited`
< result
< column "CREATE SETTINGS PROFILE omni_limited" "String" "string"
< row
: string
| CREATE SETTINGS PROFILE `omni_limited` SETTINGS max_memory_usage = 1000000 MIN 1 MAX 2000000 TO omni_reader
> query
| SHOW CREATE USER `default` SETTINGS format_display_secrets_in_show_and_select = 0
< result
< column "CREATE USER default SETTINGS format_display_secrets_in_show_and_select = 0" "String" "string"
< row
: string
| CREATE USER default IDENTIFIED WITH plaintext_password SETTINGS PROFILE `default`
> query
| SHOW GRANTS FOR `default`
< result
< column "GRANTS FOR default" "String" "string"
< row
: string
| GRANT CHECK, SHOW, SELECT, INSERT, ALTER, CREATE, DROP, UNDROP TABLE, TRUNCATE, OPTIMIZE, BACKUP, KILL QUERY, KILL TRANSACTION, MOVE PARTITION BETWEEN SHARDS, ROLE ADMIN, CREATE ROW POLICY, ALTER ROW POLICY, DROP ROW POLICY, CREATE MASKING POLICY, ALTER MASKING POLICY, DROP MASKING POLICY, CREATE QUOTA, ALTER QUOTA, DROP QUOTA, CREATE SETTINGS PROFILE, ALTER SETTINGS PROFILE, DROP SETTINGS PROFILE, ALLOW SQL SECURITY NONE, SHOW ACCESS, SHOW HANDLERS, SYSTEM, dictGet, displaySecretsInShowAndSelect, INTROSPECTION, CLUSTER, FILE, URL, REMOTE, MONGO, REDIS, MYSQL, POSTGRES, SQLITE, ODBC, JDBC, HDFS, S3, HIVE, AZURE, KAFKA, NATS, RABBITMQ, YTSAURUS, ARROW FLIGHT, BIGQUERY, DISK, SOURCES, AI FUNCTIONS ON *.* TO default WITH GRANT OPTION
< row
: string
| GRANT SOURCES ON *.* TO default WITH GRANT OPTION
< row
: string
| GRANT TABLE ENGINE ON * TO default WITH GRANT OPTION
< row
: string
| GRANT SET DEFINER ON * TO default WITH GRANT OPTION
< row
: string
| GRANT CREATE USER, ALTER USER, DROP USER, CREATE ROLE, ALTER ROLE, DROP ROLE, IMPERSONATE ON * TO default WITH GRANT OPTION
< row
: string
| GRANT ALTER NAMED COLLECTION, CREATE NAMED COLLECTION, DROP NAMED COLLECTION, SHOW NAMED COLLECTIONS, NAMED COLLECTION ON * TO default WITH GRANT OPTION
> query
| SHOW CREATE USER `omni_analyst` SETTINGS format_display_secrets_in_show_and_select = 0
< result
< column "CREATE USER omni_analyst SETTINGS format_display_secrets_in_show_and_select = 0" "String" "string"
< row
: string
| CREATE USER omni_analyst IDENTIFIED WITH sha256_password VALID UNTIL '2030-01-02 03:04:05' DEFAULT ROLE omni_reader
> query
| SHOW GRANTS FOR `omni_analyst`
< result
< column "GRANTS FOR omni_analyst" "String" "string"
< row
: string
| GRANT INSERT(id) ON omni_testsyncinstance.events TO omni_analyst WITH GRANT OPTION
< row
: string
| GRANT omni_reader TO omni_analyst
> query
| SHOW CREATE ROLE `omni_reader`
< result
< column "CREATE ROLE omni_reader" "String" "string"
< row
: string
| CREATE ROLE omni_reader SETTINGS max_threads = 2
> query
| SHOW GRANTS FOR `omni_reader`
< result
< column "GRANTS FOR omni_reader" "String" "string"
< row
: string
| GRANT SELECT ON omni_testsyncinstance.* TO omni_reader
> query
| SHOW CREATE USER `omni_service` SETTINGS format_display_secrets_in_show_and_select = 0
< result
< column "CREATE USER omni_service SETTINGS format_display_secrets_in_show_and_select = 0" "String" "string"
< row
: string
| CREATE USER omni_service IDENTIFIED WITH plaintext_password VALID UNTIL '2031-01-01 00:00:00', sha256_password
> query
| SHOW GRANTS FOR `omni_service`
< result
< column "GRANTS FOR omni_service" "String" "string"
> query
|
| 		SELECT
| 			schema_name
| 		FROM information_schema.schemata
| 		WHERE schema_name NOT IN ('INFORMATION_SCHEMA', 'information_schema', 'system')
| 		ORDER BY schema_name
< result
< column "schema_name" "String" "string"
< row
: string
| default
< row
: string
| omni_testsyncinstance
