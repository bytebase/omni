-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestGrantFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE USER u1@localhost
< ok 0 0
> exec ^
: string
| CREATE USER `u1`@localhost
< ok 0 0
> exec ^
: string
| CREATE USER 'u1'@localhost
< ok 0 0
> exec ^
: string
| CREATE USER u1@`localhost`
< ok 0 0
> exec ^
: string
| CREATE USER `u1`@`localhost`
< ok 0 0
> exec ^
: string
| CREATE USER 'u1'@`localhost`
< ok 0 0
> exec ^
: string
| CREATE USER u1@'localhost'
< ok 0 0
> exec ^
: string
| CREATE USER `u1`@'localhost'
< ok 0 0
> exec ^
: string
| CREATE USER 'u1'@'localhost'
< ok 0 0
> exec ^
: string
| CREATE USER `user@example.com`@'localhost'
< ok 0 0
> exec ^
: string
| CREATE USER u1@10.0.0.1
< ok 0 0
> exec ^
: string
| CREATE USER u1@`local host`
< ok 0 0
> exec ^
: string
| CREATE USER select@localhost
< error
| 1064 42000 …line 1 column 18 near "select@localhost" 
> exec ^
: string
| CREATE USER comment@localhost
< ok 0 0
> exec ^
: string
| CREATE USER u1@select
< ok 0 0
> exec ^
: string
| CREATE USER u1@comment
< ok 0 0
> exec ^
: string
| CREATE USER u1 @localhost
< ok 0 0
> exec ^
: string
| CREATE USER u1 @`localhost`
< ok 0 0
> exec ^
: string
| CREATE USER u1 @'localhost'
< ok 0 0
> exec ^
: string
| CREATE USER u1@ localhost
< error
| 1064 42000 …line 1 column 25 near "localhost" 
> exec ^
: string
| CREATE USER u1@ 'localhost'
< error
| 1064 42000 …line 1 column 27 near "'localhost'" 
> exec ^
: string
| CREATE USER u1 @ localhost
< error
| 1064 42000 …line 1 column 26 near "localhost" 
> exec ^
: string
| CREATE USER u1 @ 'localhost'
< error
| 1064 42000 …line 1 column 28 near "'localhost'" 
> exec ^
: string
| CREATE USER u1@
< ok 0 0
> exec ^
: string
| DROP USER u1@
< ok 0 0
> exec ^
: string
| ALTER USER u1@ IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec ^
: string
| GRANT SELECT ON *.* TO u1@
< ok 0 0
> exec ^
: string
| SET DEFAULT ROLE role1 TO u1@
< ok 0 0
> exec ^
: string
| DROP USER u1@localhost, `u2`@`localhost`
< ok 0 0
> exec ^
: string
| RENAME USER u1@localhost TO `u2`@`localhost`
< ok 0 0
> exec ^
: string
| SET DEFAULT ROLE administrator, developer TO `joe`@`10.0.0.1`
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH mysql_native_password
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH `mysql_native_password`
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' BY 'notarealpw'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 'notarealhash'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS ''
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 0x2A3132333435
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH mysql_native_password AS 0x2A3132333435
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD 0x2A3132333435
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD ''
< ok 0 0
> exec ^
: string
| ALTER USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS
< error
| 1064 42000 …line 1 column 61 near "" 
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY
< error
| 1064 42000 …line 1 column 61 near "" 
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 123
< error
| 1064 42000 …line 1 column 38 near "123" 
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY PASSWORD 'notarealhash'
< error
| 1064 42000 …line 1 column 70 near "PASSWORD 'notarealhash'" 
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD PASSWORD 'notarealhash'
< error
| 1064 42000 …line 1 column 50 near "PASSWORD 'notarealhash'" 
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY 'notarealpw', u2@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< ok 0 0
> exec ^
: string
| ALTER USER u1@'%' IDENTIFIED BY 'notarealpw', u2@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE NONE
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SSL
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE X509
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE CIPHER 'AES128-SHA'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SUBJECT 's' ISSUER 'i' CIPHER 'c'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SUBJECT 's' AND ISSUER 'i'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE CIPHER 'c' AND SUBJECT 's' AND ISSUER 'i'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user, URI:spiffe://example.org/myservice'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE token_issuer 'issuer-abc'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE CIPHER 'AES128-SHA' SAN 'DNS:example-user'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE ISSUER 'i' AND SAN 'DNS:example-user'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user' AND TOKEN_ISSUER 'issuer-abc'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE ISSUER 'i' SUBJECT 's' CIPHER 'c' SAN 'DNS:x' TOKEN_ISSUER 'ti'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SSL AND X509
< error
| 1064 42000 …line 1 column 34 near "AND X509" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SSL AND SAN 'DNS:example-user'
< error
| 1064 42000 …line 1 column 34 near "AND SAN 'DNS:example-user'" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE X509 SAN 'DNS:example-user'
< error
| 1064 42000 …line 1 column 35 near "SAN 'DNS:example-user'" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE NONE SSL
< error
| 1064 42000 …line 1 column 35 near "SSL" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE NONE AND SSL
< error
| 1064 42000 …line 1 column 35 near "AND SSL" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE CIPHER
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SAN
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE TOKEN_ISSUER
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE CIPHER 'c' AND
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SAN 0x2A31
< error
| 1064 42000 …line 1 column 37 near "0x2A31" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE `SAN` 'DNS:example-user'
< error
| 1064 42000 …line 1 column 32 near "`SAN` 'DNS:example-user'" 
> exec ^
: string
| CREATE USER u1@'%' REQUIRE `TOKEN_ISSUER` 'issuer-abc'
< error
| 1064 42000 …line 1 column 41 near "`TOKEN_ISSUER` 'issuer-abc'" 
> exec ^
: string
| CREATE USER u1@'%' PASSWORD EXPIRE DEFAULT ACCOUNT UNLOCK
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' ATTRIBUTE '{"email": "user@example.com"}'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'tidb_auth_token' REQUIRE token_issuer 'issuer-abc' ATTRIBUTE '{"email": "user@example.com"}'
< ok 0 0
> exec ^
: string
| GRANT ALL PRIVILEGES ON zz.* TO 'u1'@'localhost' IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec ^
: string
| GRANT PROCESS, USAGE ON *.* TO 'u1'@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< ok 0 0
> exec ^
: string
| REVOKE ALL PRIVILEGES ON zz.* FROM 'u1'@'localhost' IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec ^
: string
| REVOKE SELECT ON *.* FROM u1 IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 IDENTIFIED BY 'notarealpw' REQUIRE SSL WITH GRANT OPTION
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 REQUIRE SSL
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR 10
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 REQUIRE SAN 'DNS:example-user'
< ok 0 0
> exec ^
: string
| GRANT 'role1' TO u1 IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …line 1 column 30 near "IDENTIFIED BY 'notarealpw'" 
> exec ^
: string
| REVOKE 'role1' FROM u1 IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …line 1 column 33 near "IDENTIFIED BY 'notarealpw'" 
> exec ^
: string
| GRANT PROXY ON u1 TO u2 IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …line 1 column 34 near "IDENTIFIED BY 'notarealpw'" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH OPTION
< error
| 1064 42000 …line 1 column 34 near "OPTION" 
> exec ^
: string
| GRANT SELECT ON db1.t TO u1 WITH GRANT
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 REQUIRE SSL WITH GRANT
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_UPDATES_PER_HOUR 1
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_CONNECTIONS_PER_HOUR 1
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_USER_CONNECTIONS 1
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR
< error
| 1064 42000 …line 1 column 48 near "" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION MAX_QUERIES_PER_HOUR 1
< error
| 1064 42000 …line 1 column 61 near "MAX_QUERIES_PER_HOUR 1" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR 1 GRANT OPTION
< error
| 1064 42000 …line 1 column 56 near "GRANT OPTION" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR 1 MAX_UPDATES_PER_HOUR 2
< error
| 1064 42000 …line 1 column 71 near "MAX_UPDATES_PER_HOUR 2" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT MAX_QUERIES_PER_HOUR 1
< error
| 1064 42000 …line 1 column 54 near "MAX_QUERIES_PER_HOUR 1" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION WITH GRANT OPTION
< error
| 1064 42000 …line 1 column 45 near "WITH GRANT OPTION" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION, GRANT OPTION
< error
| 1064 42000 …line 1 column 41 near ", GRANT OPTION" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH `GRANT` OPTION
< error
| 1064 42000 …line 1 column 35 near "`GRANT` OPTION" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT `OPTION`
< error
| 1064 42000 …line 1 column 42 near "`OPTION`" 
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH GRANT OPTION
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH GRANT OPTION
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH MAX_USER_CONNECTIONS 1
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH MAX_USER_CONNECTIONS 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH GRANT
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH GRANT OPTION MAX_QUERIES_PER_HOUR 1
< error
| 1064 42000 …line 1 column 62 near "MAX_QUERIES_PER_HOUR 1" 
> exec ^
: string
| GRANT role1 TO u1 WITH
< error
| 1064 42000 …line 1 column 22 near "WITH" 
> exec ^
: string
| GRANT role1 TO u1 WITH ADMIN
< error
| 1064 42000 …line 1 column 22 near "WITH ADMIN" 
> exec ^
: string
| GRANT role1 TO u1 WITH ADMIN OPTION
< error
| 1064 42000 …line 1 column 22 near "WITH ADMIN OPTION" 
> exec ^
: string
| GRANT role1, role2 TO u1, u2 WITH ADMIN OPTION
< error
| 1064 42000 …line 1 column 33 near "WITH ADMIN OPTION" 
> exec ^
: string
| GRANT 'role1' TO 'u1'@'%' WITH ADMIN OPTION
< error
| 1064 42000 …line 1 column 30 near "WITH ADMIN OPTION" 
> exec ^
: string
| GRANT role1 TO u1 WITH GRANT OPTION
< error
| 1064 42000 …line 1 column 22 near "WITH GRANT OPTION" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2
< error
| 1064 42000 …line 1 column 25 near "AS u2" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS 'u2'@'%'
< error
| 1064 42000 …line 1 column 25 near "AS 'u2'@'%'" 
> exec ^
: string
| GRANT SELECT ON db1.t TO u1 AS u2
< error
| 1064 42000 …line 1 column 30 near "AS u2" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2 WITH ROLE DEFAULT
< error
| 1064 42000 …line 1 column 25 near "AS u2 WITH ROLE DEFAULT" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2 WITH ROLE NONE
< error
| 1064 42000 …line 1 column 25 near "AS u2 WITH ROLE NONE" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2 WITH ROLE ALL EXCEPT role1
< error
| 1064 42000 …line 1 column 25 near "AS u2 WITH ROLE ALL EXCEPT role1" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2 WITH ROLE role1, role2
< error
| 1064 42000 …line 1 column 25 near "AS u2 WITH ROLE role1, role2" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION AS u2
< error
| 1064 42000 …line 1 column 43 near "AS u2" 
> exec ^
: string
| GRANT ALL ON *.* TO u1 REQUIRE SSL AS u2
< error
| 1064 42000 …line 1 column 37 near "AS u2" 
> exec ^
: string
| GRANT 'role1' TO 'u1'@'localhost'
< ok 0 0
> exec ^
: string
| GRANT 'role1', 'role2' TO 'u1'@'localhost', 'u2'@'localhost'
< ok 0 0
> exec ^
: string
| GRANT 'role1'@'%', 'role2'@'%' TO 'u1'@'localhost'
< ok 0 0
> exec ^
: string
| GRANT 'u1' TO 'u1'
< ok 0 0
> exec ^
: string
| REVOKE 'role1' FROM u1
< ok 0 0
> exec ^
: string
| REVOKE 'role1', 'role2' FROM 'u1'@'localhost', 'u2'@'localhost'
< ok 0 0
> exec ^
: string
| SET DEFAULT ROLE a, b TO u1@localhost
< ok 0 0
> exec ^
: string
| GRANT 'role1' ON *.* TO u1
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| GRANT 'role1'@'%' ON *.* TO u1
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| REVOKE 'role1' ON *.* FROM u1
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| GRANT SELECT, 'role1' ON *.* TO u1
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| GRANT role1@localhost TO u1
< ok 0 0
> exec ^
: string
| GRANT role1@`localhost` TO u1
< ok 0 0
> exec ^
: string
| GRANT role1@'localhost' TO u1
< ok 0 0
> exec ^
: string
| GRANT comment@localhost TO u1
< ok 0 0
> exec ^
: string
| REVOKE role1@localhost FROM u1
< ok 0 0
> exec ^
: string
| REVOKE comment@localhost FROM u1
< ok 0 0
> exec ^
: string
| GRANT role1@localhost ON *.* TO u1
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| GRANT comment@localhost ON *.* TO u1
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| GRANT role1 TO u1
< ok 0 0
> exec ^
: string
| GRANT `select` TO u1
< ok 0 0
> exec ^
: string
| GRANT SELECT TO u1
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| GRANT COMMENT TO u1
< error
| 1064 42000 …line 1 column 16 near "TO u1" 
> exec ^
: string
| GRANT S3 ON *.* TO u1
< error
| 1064 42000 …line 1 column 11 near "ON *.* TO u1" 
> exec ^
: string
| GRANT COMMENT ON *.* TO u1
< error
| 1064 42000 …line 1 column 16 near "ON *.* TO u1" 
> exec ^
: string
| GRANT COMMENT, SELECT ON *.* TO u1
< error
| 1064 42000 …line 1 column 14 near ", SELECT ON *.* TO u1" 
> exec ^
: string
| GRANT usage TO u1
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| GRANT process TO u1
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| GRANT CREATE VIEW TO u1
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| GRANT REPLICATION SLAVE TO u1
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| GRANT GRANT OPTION TO u1
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| GRANT SEND_CREDENTIALS_TO_TIKV TO CURRENT_USER
< error
| 1064 42000 …line 1 column 33 near "TO CURRENT_USER" 
> exec ^
: string
| GRANT SEND_CREDENTIALS_TO_TIKV TO u1
< error
| 1064 42000 …line 1 column 33 near "TO u1" 
> exec ^
: string
| GRANT SEND_CREDENTIALS_TO_TIKV ON *.* TO u1
< error
| 1064 42000 …line 1 column 33 near "ON *.* TO u1" 
> exec ^
: string
| GRANT `send_credentials_to_tikv` TO u1
< ok 0 0
> exec ^
: string
| GRANT GRANT OPTION ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT GRANT OPTION ON `db`.* TO u1
< ok 0 0
> exec ^
: string
| GRANT GRANT OPTION, SELECT ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT SELECT, GRANT OPTION ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT GRANT ON *.* TO u1
< error
| 1064 42000 …line 1 column 14 near "ON *.* TO u1" 
> exec ^
: string
| GRANT INVOKE LAMBDA ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT INVOKE SAGEMAKER ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT INVOKE COMPREHEND ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT INVOKE LAMBDA, INVOKE SAGEMAKER ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT BACKUP_ADMIN ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT SELECT, BACKUP_ADMIN ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT INVOKE `LAMBDA` ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT LOAD FROM S3 ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT SELECT INTO S3 ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT LOAD FROM s3 ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT LOAD FROM `S3` ON *.* TO u1
< error
| 1064 42000 …line 1 column 20 near "`S3` ON *.* TO u1" 
> exec ^
: string
| GRANT LOAD S3 ON *.* TO u1
< error
| 1064 42000 …line 1 column 13 near "S3 ON *.* TO u1" 
> exec ^
: string
| GRANT SELECT S3 ON *.* TO u1
< error
| 1064 42000 …line 1 column 15 near "S3 ON *.* TO u1" 
> exec ^
: string
| GRANT LOAD FROM SELECT ON *.* TO u1
< error
| 1064 42000 …line 1 column 22 near "SELECT ON *.* TO u1" 
> exec ^
: string
| GRANT LOAD FROM S3 TO u1
< ok 0 0
> exec ^
: string
| GRANT INVOKE LAMBDA TO u1
< ok 0 0
> exec ^
: string
| GRANT BACKUP_ADMIN TO u1
< ok 0 0
> exec ^
: string
| GRANT `grant` TO u1
< ok 0 0
> exec ^
: string
| REVOKE `grant` FROM u1
< ok 0 0
> exec ^
: string
| GRANT `option` TO u1
< ok 0 0
> exec ^
: string
| GRANT `select` TO u1
< ok 0 0
> exec ^
: string
| GRANT `grant option` TO u1
< ok 0 0
> exec ^
: string
| GRANT `app reader` TO u1
< ok 0 0
> exec ^
: string
| REVOKE `app reader` FROM u1
< ok 0 0
> exec ^
: string
| GRANT `app reader` ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT `create` VIEW ON *.* TO u1
< error
| 1064 42000 …line 1 column 19 near "VIEW ON *.* TO u1" 
> exec ^
: string
| GRANT `replication` SLAVE ON *.* TO u1
< error
| 1064 42000 …line 1 column 25 near "SLAVE ON *.* TO u1" 
> exec ^
: string
| CREATE USER 'u1' `@localhost`
< error
| 1064 42000 …line 1 column 29 near "`@localhost`" 
> exec ^
: string
| DROP USER 'u1' `@localhost`
< error
| 1064 42000 …line 1 column 27 near "`@localhost`" 
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| GRANT SELECT ON *.* TO u1 IDENTIFIED
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| REVOKE SELECT ON *.* FROM u1 IDENTIFIED
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| GRANT SELECT ON *.* TO u1 IDENTIFIED BY
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| CREATE USER CURRENT_USER
< ok 0 0
> exec ^
: string
| CREATE USER CURRENT_USER()
< ok 0 0
> exec ^
: string
| ALTER USER CURRENT_USER IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec ^
: string
| DROP USER CURRENT_USER
< ok 0 0
> exec ^
: string
| RENAME USER CURRENT_USER TO CURRENT_USER
< ok 0 0
> exec ^
: string
| GRANT 'role1' TO CURRENT_USER
< ok 0 0
> exec ^
: string
| REVOKE 'role1' FROM CURRENT_USER
< ok 0 0
> exec ^
: string
| GRANT PROXY ON CURRENT_USER TO CURRENT_USER
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| GRANT PROXY ON CURRENT_USER TO CURRENT_USER
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| GRANT SELECT ON *.* TO CURRENT_USER
< ok 0 0
> exec ^
: string
| SET PASSWORD FOR CURRENT_USER = 'notarealpw'
< ok 0 0
> exec ^
: string
| SET DEFAULT ROLE NONE TO CURRENT_USER
< ok 0 0
> exec ^
: string
| CREATE USER u1@localhost DEFAULT ROLE current_user
< error
| 1064 42000 …line 1 column 32 near "DEFAULT ROLE current_user" 
> exec ^
: string
| ALTER USER u1@localhost DEFAULT ROLE current_user
< error
| 1064 42000 …line 1 column 31 near "DEFAULT ROLE current_user" 
> exec ^
: string
| ALTER USER USER() IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec ^
: string
| ALTER USER USER(IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …line 1 column 26 near "IDENTIFIED BY 'notarealpw'" 
> exec ^
: string
| CREATE USER IF NOT EXISTS u
< ok 0 0
> exec ^
: string
| ALTER USER IF EXISTS u IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec ^
: string
| DROP USER IF EXISTS u
< ok 0 0
> exec ^
: string
| RENAME USER a TO b, c TO d
< ok 0 0
> exec ^
: string
| CREATE USER IF EXISTS u
< error
| 1064 42000 …line 1 column 21 near "EXISTS u" 
> exec ^
: string
| CREATE USER IF NOT u
< error
| 1064 42000 …line 1 column 20 near "u" 
> exec ^
: string
| ALTER USER IF u
< error
| 1064 42000 …line 1 column 15 near "u" 
> exec ^
: string
| DROP USER IF u
< error
| 1064 42000 …line 1 column 14 near "u" 
> exec ^
: string
| RENAME USER a b
< error
| 1064 42000 …line 1 column 15 near "b" 
> exec ^
: string
| RENAME USER a b, c TO d
< error
| 1064 42000 …line 1 column 15 near "b, c TO d" 
> exec ^
: string
| RENAME USER a TO b, c d
< error
| 1064 42000 …line 1 column 23 near "d" 
