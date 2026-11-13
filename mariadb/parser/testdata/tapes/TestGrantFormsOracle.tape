-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestGrantFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| GRANT SELECT TO u1
< error
| 1064 42000 …near 'TO u1' at line 1
> exec ^
: string
| GRANT usage TO u1
< error
| 1064 42000 …near 'TO u1' at line 1
> exec ^
: string
| GRANT process TO u1
< ok 0 0
> exec ^
: string
| GRANT `select` TO u1
< ok 0 0
> exec ^
: string
| GRANT ALL TO u1
< error
| 1064 42000 …near 'TO u1' at line 1
> exec ^
: string
| GRANT r1 ON *.* TO u1
< error
| 1064 42000 …near 'ON *.* TO u1' at line 1
> exec ^
: string
| GRANT `select` ON *.* TO u1
< error
| 1064 42000 …near 'ON *.* TO u1' at line 1
> exec ^
: string
| GRANT BACKUP_ADMIN ON *.* TO u1
< error
| 1064 42000 …near 'ON *.* TO u1' at line 1
> exec ^
: string
| GRANT SELECT, BACKUP_ADMIN ON *.* TO u1
< error
| 1064 42000 …near 'BACKUP_ADMIN ON *.* TO u1' at line 1
> exec ^
: string
| GRANT BINLOG MONITOR ON *.* TO u1
< ok 0 0
> exec ^
: string
| GRANT REPLICATION SLAVE ADMIN ON *.* TO u1
< ok 0 0
> exec ^
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
| 1064 42000 …near 'select@localhost' at line 1
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
| 1064 42000 …near 'localhost' at line 1
> exec ^
: string
| CREATE USER u1@ 'localhost'
< error
| 1064 42000 …near ''localhost'' at line 1
> exec ^
: string
| CREATE USER u1 @ localhost
< error
| 1064 42000 …near 'localhost' at line 1
> exec ^
: string
| CREATE USER u1 @ 'localhost'
< error
| 1064 42000 …near ''localhost'' at line 1
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
< error
| 1064 42000 …near 'TO u1@' at line 1
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
< error
| 1064 42000 …near 'TO `joe`@`10.0.0.1`' at line 1
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
< error
| 1064 42000 …near 'BY 'notarealpw'' at line 1
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
< error
| 1064 42000 …near '0x2A3132333435' at line 1
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH mysql_native_password AS 0x2A3132333435
< error
| 1064 42000 …near '0x2A3132333435' at line 1
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD 0x2A3132333435
< error
| 1064 42000 …near '0x2A3132333435' at line 1
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
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY
< error
| 1064 42000 …near 'BY' at line 1
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 123
< error
| 1064 42000 …near '123' at line 1
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY PASSWORD 'notarealhash'
< error
| 1064 42000 …near 'BY PASSWORD 'notarealhash'' at line 1
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD PASSWORD 'notarealhash'
< error
| 1064 42000 …near 'PASSWORD 'notarealhash'' at line 1
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
< error
| 1064 42000 …near 'SAN 'DNS:example-user, URI:spiffe://example.org/myservice'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'
< error
| 1064 42000 …near 'TOKEN_ISSUER 'issuer-abc'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE token_issuer 'issuer-abc'
< error
| 1064 42000 …near 'token_issuer 'issuer-abc'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE CIPHER 'AES128-SHA' SAN 'DNS:example-user'
< error
| 1064 42000 …near 'SAN 'DNS:example-user'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE ISSUER 'i' AND SAN 'DNS:example-user'
< error
| 1064 42000 …near 'SAN 'DNS:example-user'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user' AND TOKEN_ISSUER 'issuer-abc'
< error
| 1064 42000 …near 'SAN 'DNS:example-user' AND TOKEN_ISSUER 'issuer-abc'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE ISSUER 'i' SUBJECT 's' CIPHER 'c' SAN 'DNS:x' TOKEN_ISSUER 'ti'
< error
| 1064 42000 …near 'SAN 'DNS:x' TOKEN_ISSUER 'ti'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SSL AND X509
< error
| 1064 42000 …near 'AND X509' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SSL AND SAN 'DNS:example-user'
< error
| 1064 42000 …near 'AND SAN 'DNS:example-user'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE X509 SAN 'DNS:example-user'
< error
| 1064 42000 …near 'SAN 'DNS:example-user'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE NONE SSL
< error
| 1064 42000 …near 'SSL' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE NONE AND SSL
< error
| 1064 42000 …near 'AND SSL' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE CIPHER
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SAN
< error
| 1064 42000 …near 'SAN' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE TOKEN_ISSUER
< error
| 1064 42000 …near 'TOKEN_ISSUER' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE CIPHER 'c' AND
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE SAN 0x2A31
< error
| 1064 42000 …near 'SAN 0x2A31' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE `SAN` 'DNS:example-user'
< error
| 1064 42000 …near '`SAN` 'DNS:example-user'' at line 1
> exec ^
: string
| CREATE USER u1@'%' REQUIRE `TOKEN_ISSUER` 'issuer-abc'
< error
| 1064 42000 …near '`TOKEN_ISSUER` 'issuer-abc'' at line 1
> exec ^
: string
| CREATE USER u1@'%' PASSWORD EXPIRE DEFAULT ACCOUNT UNLOCK
< ok 0 0
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED WITH 'tidb_auth_token' REQUIRE token_issuer 'issuer-abc' ATTRIBUTE '{"email": "user@example.com"}'
< error
| 1064 42000 …near 'token_issuer 'issuer-abc' ATTRIBUTE '{"email": "user@example.com"}'' at line 1
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
< error
| 1064 42000 …near 'IDENTIFIED BY 'notarealpw'' at line 1
> exec ^
: string
| REVOKE SELECT ON *.* FROM u1 IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …near 'IDENTIFIED BY 'notarealpw'' at line 1
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
< error
| 1064 42000 …near 'SAN 'DNS:example-user'' at line 1
> exec ^
: string
| GRANT 'role1' TO u1 IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec ^
: string
| REVOKE 'role1' FROM u1 IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …near 'IDENTIFIED BY 'notarealpw'' at line 1
> exec ^
: string
| GRANT PROXY ON u1 TO u2 IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH OPTION
< error
| 1064 42000 …near 'OPTION' at line 1
> exec ^
: string
| GRANT SELECT ON db1.t TO u1 WITH GRANT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 REQUIRE SSL WITH GRANT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION MAX_QUERIES_PER_HOUR 1
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT MAX_QUERIES_PER_HOUR 1
< error
| 1064 42000 …near 'MAX_QUERIES_PER_HOUR 1' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH `GRANT` OPTION
< error
| 1064 42000 …near '`GRANT` OPTION' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT `OPTION`
< error
| 1064 42000 …near '`OPTION`' at line 1
> exec ^
: string
| GRANT role1 TO u1 WITH `ADMIN` OPTION
< error
| 1064 42000 …near '`ADMIN` OPTION' at line 1
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH GRANT OPTION
< ok 0 0
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH GRANT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH MAX_QUERIES_PER_HOUR 1
< error
| 1064 42000 …near 'MAX_QUERIES_PER_HOUR 1' at line 1
> exec ^
: string
| GRANT role1 TO u1 WITH ADMIN OPTION
< ok 0 0
> exec ^
: string
| GRANT role1 TO u1 WITH ADMIN
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT role1 TO u1 WITH
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT role1 TO u1 WITH GRANT OPTION
< error
| 1064 42000 …near 'GRANT OPTION' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION GRANT OPTION
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR 1 GRANT OPTION
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION MAX_QUERIES_PER_HOUR 1 GRANT OPTION
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR 1 MAX_QUERIES_PER_HOUR 2
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_USER_CONNECTIONS 1 GRANT OPTION MAX_UPDATES_PER_HOUR 1
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 REQUIRE SSL WITH MAX_USER_CONNECTIONS 1 GRANT OPTION
< ok 0 0
> exec ^
: string
| GRANT SELECT ON db1.t TO u1 WITH GRANT OPTION GRANT OPTION
< ok 0 0
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION, GRANT OPTION
< error
| 1064 42000 …near ' GRANT OPTION' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION WITH GRANT OPTION
< error
| 1064 42000 …near 'WITH GRANT OPTION' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION GRANT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR 1 GRANT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION MAX_QUERIES_PER_HOUR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH GRANT OPTION GRANT OPTION
< error
| 1064 42000 …near 'GRANT OPTION' at line 1
> exec ^
: string
| GRANT PROXY ON u1 TO u2 WITH GRANT OPTION MAX_QUERIES_PER_HOUR 1
< error
| 1064 42000 …near 'MAX_QUERIES_PER_HOUR 1' at line 1
> exec ^
: string
| GRANT role1 TO u1 WITH ADMIN OPTION GRANT OPTION
< error
| 1064 42000 …near 'GRANT OPTION' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2
< error
| 1064 42000 …near 'AS u2' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS 'u2'@'%'
< error
| 1064 42000 …near 'AS 'u2'@'%'' at line 1
> exec ^
: string
| GRANT SELECT ON db1.t TO u1 AS u2
< error
| 1064 42000 …near 'AS u2' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2 WITH ROLE DEFAULT
< error
| 1064 42000 …near 'AS u2 WITH ROLE DEFAULT' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2 WITH ROLE NONE
< error
| 1064 42000 …near 'AS u2 WITH ROLE NONE' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2 WITH ROLE ALL EXCEPT role1
< error
| 1064 42000 …near 'AS u2 WITH ROLE ALL EXCEPT role1' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 AS u2 WITH ROLE role1, role2
< error
| 1064 42000 …near 'AS u2 WITH ROLE role1, role2' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u1 WITH GRANT OPTION AS u2
< error
| 1064 42000 …near 'AS u2' at line 1
> exec ^
: string
| GRANT 'role1' TO 'u1'@'localhost'
< ok 0 0
> exec ^
: string
| GRANT 'role1', 'role2' TO 'u1'@'localhost', 'u2'@'localhost'
< error
| 1064 42000 …near ' 'role2' TO 'u1'@'localhost', 'u2'@'localhost'' at line 1
> exec ^
: string
| GRANT 'role1'@'%', 'role2'@'%' TO 'u1'@'localhost'
< error
| 1064 42000 …near '@'%', 'role2'@'%' TO 'u1'@'localhost'' at line 1
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
< error
| 1064 42000 …near ' 'role2' FROM 'u1'@'localhost', 'u2'@'localhost'' at line 1
> exec ^
: string
| SET DEFAULT ROLE a, b TO u1@localhost
< error
| 1064 42000 …near 'TO u1@localhost' at line 1
> exec ^
: string
| GRANT 'role1' ON *.* TO u1
< error
| 1064 42000 …near 'ON *.* TO u1' at line 1
> exec ^
: string
| GRANT 'role1'@'%' ON *.* TO u1
< error
| 1064 42000 …near '@'%' ON *.* TO u1' at line 1
> exec ^
: string
| REVOKE 'role1' ON *.* FROM u1
< error
| 1064 42000 …near 'ON *.* FROM u1' at line 1
> exec ^
: string
| GRANT SELECT, 'role1' ON *.* TO u1
< error
| 1064 42000 …near ''role1' ON *.* TO u1' at line 1
> exec ^
: string
| GRANT role1@localhost TO u1
< error
| 1064 42000 …near '@localhost TO u1' at line 1
> exec ^
: string
| GRANT role1@`localhost` TO u1
< error
| 1064 42000 …near '@`localhost` TO u1' at line 1
> exec ^
: string
| GRANT role1@'localhost' TO u1
< error
| 1064 42000 …near '@'localhost' TO u1' at line 1
> exec ^
: string
| GRANT comment@localhost TO u1
< error
| 1064 42000 …near '@localhost TO u1' at line 1
> exec ^
: string
| REVOKE role1@localhost FROM u1
< error
| 1064 42000 …near '@localhost FROM u1' at line 1
> exec ^
: string
| REVOKE comment@localhost FROM u1
< error
| 1064 42000 …near '@localhost FROM u1' at line 1
> exec ^
: string
| GRANT role1@localhost ON *.* TO u1
< error
| 1064 42000 …near '@localhost ON *.* TO u1' at line 1
> exec ^
: string
| GRANT comment@localhost ON *.* TO u1
< error
| 1064 42000 …near '@localhost ON *.* TO u1' at line 1
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
| GRANT COMMENT TO u1
< ok 0 0
> exec ^
: string
| GRANT process TO u1
< ok 0 0
> exec ^
: string
| GRANT CREATE VIEW TO u1
< error
| 1064 42000 …near 'TO u1' at line 1
> exec ^
: string
| GRANT REPLICATION SLAVE TO u1
< error
| 1064 42000 …near 'TO u1' at line 1
> exec ^
: string
| GRANT GRANT OPTION TO u1
< error
| 1064 42000 …near 'TO u1' at line 1
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
| 1064 42000 …near 'ON *.* TO u1' at line 1
> exec ^
: string
| GRANT INVOKE LAMBDA ON *.* TO u1
< error
| 1064 42000 …near 'LAMBDA ON *.* TO u1' at line 1
> exec ^
: string
| GRANT INVOKE SAGEMAKER ON *.* TO u1
< error
| 1064 42000 …near 'SAGEMAKER ON *.* TO u1' at line 1
> exec ^
: string
| GRANT INVOKE COMPREHEND ON *.* TO u1
< error
| 1064 42000 …near 'COMPREHEND ON *.* TO u1' at line 1
> exec ^
: string
| GRANT INVOKE LAMBDA, INVOKE SAGEMAKER ON *.* TO u1
< error
| 1064 42000 …near 'LAMBDA, INVOKE SAGEMAKER ON *.* TO u1' at line 1
> exec ^
: string
| GRANT INVOKE `LAMBDA` ON *.* TO u1
< error
| 1064 42000 …near '`LAMBDA` ON *.* TO u1' at line 1
> exec ^
: string
| GRANT LOAD FROM S3 ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM S3 ON *.* TO u1' at line 1
> exec ^
: string
| GRANT SELECT INTO S3 ON *.* TO u1
< error
| 1064 42000 …near 'INTO S3 ON *.* TO u1' at line 1
> exec ^
: string
| GRANT LOAD FROM s3 ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM s3 ON *.* TO u1' at line 1
> exec ^
: string
| GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA ON *.* TO u1' at line 1
> exec ^
: string
| GRANT LOAD FROM `S3` ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM `S3` ON *.* TO u1' at line 1
> exec ^
: string
| GRANT LOAD S3 ON *.* TO u1
< error
| 1064 42000 …near 'LOAD S3 ON *.* TO u1' at line 1
> exec ^
: string
| GRANT SELECT S3 ON *.* TO u1
< error
| 1064 42000 …near 'S3 ON *.* TO u1' at line 1
> exec ^
: string
| GRANT LOAD FROM SELECT ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM SELECT ON *.* TO u1' at line 1
> exec ^
: string
| GRANT LOAD FROM S3 TO u1
< error
| 1064 42000 …near 'LOAD FROM S3 TO u1' at line 1
> exec ^
: string
| GRANT INVOKE LAMBDA TO u1
< error
| 1064 42000 …near 'LAMBDA TO u1' at line 1
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
| GRANT `create` VIEW ON *.* TO u1
< error
| 1064 42000 …near 'VIEW ON *.* TO u1' at line 1
> exec ^
: string
| GRANT `replication` SLAVE ON *.* TO u1
< error
| 1064 42000 …near 'SLAVE ON *.* TO u1' at line 1
> exec ^
: string
| CREATE USER 'u1' `@localhost`
< error
| 1064 42000 …near '`@localhost`' at line 1
> exec ^
: string
| DROP USER 'u1' `@localhost`
< error
| 1064 42000 …near '`@localhost`' at line 1
> exec ^
: string
| CREATE USER u1@'%' IDENTIFIED
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| GRANT SELECT ON *.* TO u1 IDENTIFIED
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| REVOKE SELECT ON *.* FROM u1 IDENTIFIED
< error
| 1064 42000 …near 'IDENTIFIED' at line 1
> exec ^
: string
| GRANT SELECT ON *.* TO u1 IDENTIFIED BY
< error
| 1064 42000 …near '' at line 1
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
< ok 0 0
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
| CREATE ROLE current_user
< error
| 1064 42000 …near 'current_user' at line 1
> exec ^
: string
| DROP ROLE current_user
< error
| 1064 42000 …near 'current_user' at line 1
> exec ^
: string
| SET ROLE current_user
< error
| 1064 42000 …near 'current_user' at line 1
> exec ^
: string
| CREATE USER u1@localhost DEFAULT ROLE current_user
< error
| 1064 42000 …near 'DEFAULT ROLE current_user' at line 1
> exec ^
: string
| ALTER USER u1@localhost DEFAULT ROLE current_user
< error
| 1064 42000 …near 'DEFAULT ROLE current_user' at line 1
> exec ^
: string
| SHOW GRANTS FOR u1@localhost USING current_user
< error
| 1064 42000 …near 'USING current_user' at line 1
> exec ^
: string
| SET DEFAULT ROLE CURRENT_USER
< error
| 1064 42000 …near 'CURRENT_USER' at line 1
> exec ^
: string
| SET DEFAULT ROLE CURRENT_USER() FOR 'u'@'localhost'
< error
| 1064 42000 …near 'CURRENT_USER() FOR 'u'@'localhost'' at line 1
> exec ^
: string
| SET DEFAULT ROLE 'CURRENT_USER'
< error
| 1959 OP000 Invalid role specification `CURRENT_USER`
> exec ^
: string
| SET DEFAULT ROLE 'CURRENT_USER'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SET DEFAULT ROLE role1 FOR CURRENT_USER
< error
| 1959 OP000 Invalid role specification `role1`
> exec ^
: string
| SET DEFAULT ROLE role1 FOR CURRENT_USER
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
