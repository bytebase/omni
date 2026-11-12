-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestGrantFormsOnOceanBase default#1
> exec
| CREATE USER u1@localhost
< ok 0 0
> exec
| CREATE USER `u1`@localhost
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER 'u1'@localhost
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER u1@`localhost`
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER `u1`@`localhost`
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER 'u1'@`localhost`
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER u1@'localhost'
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER `u1`@'localhost'
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER 'u1'@'localhost'
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER `user@example.com`@'localhost'
< ok 0 0
> exec
| CREATE USER u1@10.0.0.1
< ok 0 0
> exec
| CREATE USER u1@`local host`
< error
| 1064 42000 …near '@`local host`' at line 1
> exec
| CREATE USER select@localhost
< error
| 1064 42000 …near 'select@localhost' at line 1
> exec
| CREATE USER comment@localhost
< ok 0 0
> exec
| CREATE USER u1@select
< ok 0 0
> exec
| CREATE USER u1@comment
< ok 0 0
> exec
| CREATE USER u1 @localhost
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER u1 @`localhost`
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER u1 @'localhost'
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'localhost'
> exec
| CREATE USER u1@ localhost
< error
| 1064 42000 …near 'localhost' at line 1
> exec
| CREATE USER u1@ 'localhost'
< error
| 1064 42000 …near ''localhost'' at line 1
> exec
| CREATE USER u1 @ localhost
< error
| 1064 42000 …near 'localhost' at line 1
> exec
| CREATE USER u1 @ 'localhost'
< error
| 1064 42000 …near ''localhost'' at line 1
> exec
| CREATE USER u1@
< ok 0 0
> exec
| DROP USER u1@
< ok 0 0
> exec
| ALTER USER u1@ IDENTIFIED BY 'notarealpw'
< error
| 1133 42000 Can not find any matching row in the user table
> exec
| GRANT SELECT ON *.* TO u1@
< error
| 1045 42000 Access denied for user 'u1'@'' (using password: NO)
> exec
| SET DEFAULT ROLE role1 TO u1@
< error
| 3523 HY000 Unknown authorization ID `u1`@``
> exec
| DROP USER u1@localhost, `u2`@`localhost`
< error
| 1396 HY000 Operation DROP USER failed for 'u2'@'localhost'
> exec
| RENAME USER u1@localhost TO `u2`@`localhost`
< error
| 1396 HY000 Operation RENAME USER failed for 'u1'@'localhost'
> exec
| SET DEFAULT ROLE administrator, developer TO `joe`@`10.0.0.1`
< error
| 3523 HY000 Unknown authorization ID `joe`@`10.0.0.1`
> exec
| CREATE USER u1@'%' IDENTIFIED WITH mysql_native_password
< error
| 1064 42000 …near 'mysql_native_password' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password'
< error
| 1064 42000 …near ''mysql_native_password'' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED WITH `mysql_native_password`
< error
| 1064 42000 …near '`mysql_native_password`' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' BY 'notarealpw'
< error
| 1524 HY000 Plugin 'caching_sha2_password' is not loaded
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 'notarealhash'
< error
| 1524 HY000 Plugin 'caching_sha2_password' is not loaded
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS ''
< error
| 1827 42000 The password hash doesn't have the expected format. Check if the correct password algorithm is being used with the PASSWORD() function.
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 'caching_sha2_password' AS 0x2A3132333435
< error
| 1064 42000 …near '0x2A3132333435' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED WITH mysql_native_password AS 0x2A3132333435
< error
| 1064 42000 …near '0x2A3132333435' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< error
| 1827 42000 The password hash doesn't have the expected format. Check if the correct password algorithm is being used with the PASSWORD() function.
> exec
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD 0x2A3132333435
< error
| 1064 42000 …near '0x2A3132333435' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD ''
< error
| 1827 42000 The password hash doesn't have the expected format. Check if the correct password algorithm is being used with the PASSWORD() function.
> exec
| ALTER USER u1@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< error
| 1064 42000 …near 'PASSWORD 'notarealhash'' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD
< error
| 1064 42000 …near 'PASSWORD' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' AS
< error
| 1064 42000 …near 'AS' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY
< error
| 1064 42000 …near 'BY' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED BY
< error
| 1064 42000 …near 'BY' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 123
< error
| 1064 42000 …near '123' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 'mysql_native_password' BY PASSWORD 'notarealhash'
< error
| 1827 42000 The password hash doesn't have the expected format. Check if the correct password algorithm is being used with the PASSWORD() function.
> exec
| CREATE USER u1@'%' IDENTIFIED BY PASSWORD PASSWORD 'notarealhash'
< error
| 1064 42000 …near 'PASSWORD 'notarealhash'' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED BY 'notarealpw', u2@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< error
| 1827 42000 The password hash doesn't have the expected format. Check if the correct password algorithm is being used with the PASSWORD() function.
> exec
| ALTER USER u1@'%' IDENTIFIED BY 'notarealpw', u2@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< error
| 1064 42000 …near ', u2@'%' IDENTIFIED BY PASSWORD 'notarealhash'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE NONE
< ok 0 0
> exec
| CREATE USER u1@'%' REQUIRE SSL
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'%'
> exec
| CREATE USER u1@'%' REQUIRE X509
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'%'
> exec
| CREATE USER u1@'%' REQUIRE CIPHER 'AES128-SHA'
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'%'
> exec
| CREATE USER u1@'%' REQUIRE SUBJECT 's' ISSUER 'i' CIPHER 'c'
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'%'
> exec
| CREATE USER u1@'%' REQUIRE SUBJECT 's' AND ISSUER 'i'
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'%'
> exec
| CREATE USER u1@'%' REQUIRE CIPHER 'c' AND SUBJECT 's' AND ISSUER 'i'
< error
| 1396 HY000 Operation CREATE USER failed for 'u1'@'%'
> exec
| CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user, URI:spiffe://example.org/myservice'
< error
| 1064 42000 …near 'SAN 'DNS:example-user, URI:spiffe://example.org/myservice'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'
< error
| 1064 42000 …near 'TOKEN_ISSUER 'issuer-abc'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE token_issuer 'issuer-abc'
< error
| 1064 42000 …near 'token_issuer 'issuer-abc'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE CIPHER 'AES128-SHA' SAN 'DNS:example-user'
< error
| 1064 42000 …near 'SAN 'DNS:example-user'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE ISSUER 'i' AND SAN 'DNS:example-user'
< error
| 1064 42000 …near 'SAN 'DNS:example-user'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE SAN 'DNS:example-user' AND TOKEN_ISSUER 'issuer-abc'
< error
| 1064 42000 …near 'SAN 'DNS:example-user' AND TOKEN_ISSUER 'issuer-abc'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE ISSUER 'i' SUBJECT 's' CIPHER 'c' SAN 'DNS:x' TOKEN_ISSUER 'ti'
< error
| 1064 42000 …near 'SAN 'DNS:x' TOKEN_ISSUER 'ti'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE SSL AND X509
< error
| 1064 42000 …near 'AND X509' at line 1
> exec
| CREATE USER u1@'%' REQUIRE SSL AND SAN 'DNS:example-user'
< error
| 1064 42000 …near 'AND SAN 'DNS:example-user'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE X509 SAN 'DNS:example-user'
< error
| 1064 42000 …near 'SAN 'DNS:example-user'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE NONE SSL
< error
| 1064 42000 …near 'SSL' at line 1
> exec
| CREATE USER u1@'%' REQUIRE NONE AND SSL
< error
| 1064 42000 …near 'AND SSL' at line 1
> exec
| CREATE USER u1@'%' REQUIRE
< error
| 1064 42000 …near 'REQUIRE' at line 1
> exec
| CREATE USER u1@'%' REQUIRE CIPHER
< error
| 1064 42000 …near 'CIPHER' at line 1
> exec
| CREATE USER u1@'%' REQUIRE SAN
< error
| 1064 42000 …near 'SAN' at line 1
> exec
| CREATE USER u1@'%' REQUIRE TOKEN_ISSUER
< error
| 1064 42000 …near 'TOKEN_ISSUER' at line 1
> exec
| CREATE USER u1@'%' REQUIRE CIPHER 'c' AND
< error
| 1064 42000 …near 'AND' at line 1
> exec
| CREATE USER u1@'%' REQUIRE SAN 0x2A31
< error
| 1064 42000 …near 'SAN 0x2A31' at line 1
> exec
| CREATE USER u1@'%' REQUIRE `SAN` 'DNS:example-user'
< error
| 1064 42000 …near '`SAN` 'DNS:example-user'' at line 1
> exec
| CREATE USER u1@'%' REQUIRE `TOKEN_ISSUER` 'issuer-abc'
< error
| 1064 42000 …near '`TOKEN_ISSUER` 'issuer-abc'' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED WITH 'tidb_auth_token' REQUIRE token_issuer 'issuer-abc' ATTRIBUTE '{"email": "user@example.com"}'
< error
| 1064 42000 …near 'REQUIRE token_issuer 'issuer-abc' ATTRIBUTE '{"email": "user@example.com"}'' at line 1
> exec
| GRANT ALL PRIVILEGES ON zz.* TO 'u1'@'localhost' IDENTIFIED BY 'notarealpw'
< ok 0 0
> exec
| GRANT PROCESS, USAGE ON *.* TO 'u1'@'%' IDENTIFIED BY PASSWORD 'notarealhash'
< error
| 1827 42000 The password hash doesn't have the expected format. Check if the correct password algorithm is being used with the PASSWORD() function.
> exec
| REVOKE ALL PRIVILEGES ON zz.* FROM 'u1'@'localhost' IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …near 'IDENTIFIED BY 'notarealpw'' at line 1
> exec
| REVOKE SELECT ON *.* FROM u1 IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …near 'IDENTIFIED BY 'notarealpw'' at line 1
> exec
| GRANT ALL ON *.* TO u1 IDENTIFIED BY 'notarealpw' REQUIRE SSL WITH GRANT OPTION
< error
| 1064 42000 …near 'REQUIRE SSL WITH GRANT OPTION' at line 1
> exec
| GRANT ALL ON *.* TO u1 REQUIRE SSL
< error
| 1064 42000 …near 'REQUIRE SSL' at line 1
> exec
| GRANT ALL ON *.* TO u1 WITH MAX_QUERIES_PER_HOUR 10
< error
| 1064 42000 …near 'MAX_QUERIES_PER_HOUR 10' at line 1
> exec
| GRANT ALL ON *.* TO u1 REQUIRE SAN 'DNS:example-user'
< error
| 1064 42000 …near 'REQUIRE SAN 'DNS:example-user'' at line 1
> exec
| GRANT 'role1' TO u1 IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …near 'IDENTIFIED BY 'notarealpw'' at line 1
> exec
| REVOKE 'role1' FROM u1 IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …near 'IDENTIFIED BY 'notarealpw'' at line 1
> exec
| GRANT PROXY ON u1 TO u2 IDENTIFIED BY 'notarealpw'
< error
| 1064 42000 …near 'IDENTIFIED BY 'notarealpw'' at line 1
> exec
| GRANT 'role1' TO 'u1'@'localhost'
< error
| 3523 HY000 Unknown authorization ID `role1`@`%`
> exec
| GRANT 'role1', 'role2' TO 'u1'@'localhost', 'u2'@'localhost'
< error
| 3523 HY000 Unknown authorization ID `u2`@`localhost`
> exec
| GRANT 'role1'@'%', 'role2'@'%' TO 'u1'@'localhost'
< error
| 3523 HY000 Unknown authorization ID `role1`@`%`
> exec
| GRANT 'u1' TO 'u1'
< error
| 4027 HY000 User account `u1`@`%` is directly or indirectly granted to the role `u1`@`%`. The GRANT would create a loop
> exec
| REVOKE 'role1' FROM u1
< error
| 3523 HY000 Unknown authorization ID `role1`@`%`
> exec
| REVOKE 'role1', 'role2' FROM 'u1'@'localhost', 'u2'@'localhost'
< error
| 3523 HY000 Unknown authorization ID `role1`@`%`
> exec
| SET DEFAULT ROLE a, b TO u1@localhost
< error
| 3523 HY000 Unknown authorization ID `a`@`%`
> exec
| GRANT 'role1' ON *.* TO u1
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| GRANT 'role1'@'%' ON *.* TO u1
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| REVOKE 'role1' ON *.* FROM u1
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| GRANT SELECT, 'role1' ON *.* TO u1
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| GRANT role1@localhost TO u1
< error
| 3523 HY000 Unknown authorization ID `role1`@`localhost`
> exec
| GRANT role1@`localhost` TO u1
< error
| 3523 HY000 Unknown authorization ID `role1`@`localhost`
> exec
| GRANT role1@'localhost' TO u1
< error
| 3523 HY000 Unknown authorization ID `role1`@`localhost`
> exec
| GRANT comment@localhost TO u1
< ok 0 0
> exec
| REVOKE role1@localhost FROM u1
< error
| 3523 HY000 Unknown authorization ID `role1`@`localhost`
> exec
| REVOKE comment@localhost FROM u1
< ok 0 0
> exec
| GRANT role1@localhost ON *.* TO u1
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| GRANT comment@localhost ON *.* TO u1
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| GRANT role1 TO u1
< error
| 3523 HY000 Unknown authorization ID `role1`@`%`
> exec
| GRANT `select` TO u1
< error
| 3523 HY000 Unknown authorization ID `select`@`%`
> exec
| GRANT COMMENT TO u1
< error
| 3523 HY000 Unknown authorization ID `COMMENT`@`%`
> exec
| GRANT CREATE VIEW TO u1
< error
| 1064 42000 …near 'CREATE VIEW TO u1' at line 1
> exec
| GRANT REPLICATION SLAVE TO u1
< error
| 1064 42000 …near 'REPLICATION SLAVE TO u1' at line 1
> exec
| GRANT GRANT OPTION TO u1
< error
| 1064 42000 …near 'GRANT OPTION TO u1' at line 1
> exec
| GRANT GRANT OPTION ON *.* TO u1
< ok 0 0
> exec
| GRANT GRANT OPTION ON `db`.* TO u1
< ok 0 0
> exec
| GRANT GRANT OPTION, SELECT ON *.* TO u1
< ok 0 0
> exec
| GRANT SELECT, GRANT OPTION ON *.* TO u1
< ok 0 0
> exec
| GRANT GRANT ON *.* TO u1
< error
| 1064 42000 …near 'ON *.* TO u1' at line 1
> exec
| GRANT INVOKE LAMBDA ON *.* TO u1
< error
| 1064 42000 …near 'LAMBDA ON *.* TO u1' at line 1
> exec
| GRANT INVOKE SAGEMAKER ON *.* TO u1
< error
| 1064 42000 …near 'SAGEMAKER ON *.* TO u1' at line 1
> exec
| GRANT INVOKE COMPREHEND ON *.* TO u1
< error
| 1064 42000 …near 'COMPREHEND ON *.* TO u1' at line 1
> exec
| GRANT INVOKE LAMBDA, INVOKE SAGEMAKER ON *.* TO u1
< error
| 1064 42000 …near 'LAMBDA, INVOKE SAGEMAKER ON *.* TO u1' at line 1
> exec
| GRANT BACKUP_ADMIN ON *.* TO u1
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| GRANT SELECT, BACKUP_ADMIN ON *.* TO u1
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| GRANT INVOKE `LAMBDA` ON *.* TO u1
< error
| 1064 42000 …near '`LAMBDA` ON *.* TO u1' at line 1
> exec
| GRANT LOAD FROM S3 ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM S3 ON *.* TO u1' at line 1
> exec
| GRANT SELECT INTO S3 ON *.* TO u1
< error
| 1064 42000 …near 'INTO S3 ON *.* TO u1' at line 1
> exec
| GRANT LOAD FROM s3 ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM s3 ON *.* TO u1' at line 1
> exec
| GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA ON *.* TO u1' at line 1
> exec
| GRANT LOAD FROM `S3` ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM `S3` ON *.* TO u1' at line 1
> exec
| GRANT LOAD S3 ON *.* TO u1
< error
| 1064 42000 …near 'LOAD S3 ON *.* TO u1' at line 1
> exec
| GRANT SELECT S3 ON *.* TO u1
< error
| 1064 42000 …near 'S3 ON *.* TO u1' at line 1
> exec
| GRANT LOAD FROM SELECT ON *.* TO u1
< error
| 1064 42000 …near 'LOAD FROM SELECT ON *.* TO u1' at line 1
> exec
| GRANT LOAD FROM S3 TO u1
< error
| 1064 42000 …near 'LOAD FROM S3 TO u1' at line 1
> exec
| GRANT INVOKE LAMBDA TO u1
< error
| 1064 42000 …near 'LAMBDA TO u1' at line 1
> exec
| GRANT `grant` TO u1
< error
| 3523 HY000 Unknown authorization ID `grant`@`%`
> exec
| REVOKE `grant` FROM u1
< error
| 3523 HY000 Unknown authorization ID `grant`@`%`
> exec
| GRANT `option` TO u1
< error
| 3523 HY000 Unknown authorization ID `option`@`%`
> exec
| GRANT `select` TO u1
< error
| 3523 HY000 Unknown authorization ID `select`@`%`
> exec
| GRANT `grant option` TO u1
< error
| 3523 HY000 Unknown authorization ID `grant option`@`%`
> exec
| GRANT `app reader` TO u1
< error
| 3523 HY000 Unknown authorization ID `app reader`@`%`
> exec
| REVOKE `app reader` FROM u1
< error
| 3523 HY000 Unknown authorization ID `app reader`@`%`
> exec
| GRANT `app reader` ON *.* TO u1
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| GRANT `create` VIEW ON *.* TO u1
< error
| 1064 42000 …near 'VIEW ON *.* TO u1' at line 1
> exec
| GRANT `replication` SLAVE ON *.* TO u1
< error
| 1064 42000 …near 'SLAVE ON *.* TO u1' at line 1
> exec
| CREATE USER 'u1' `@localhost`
< error
| 1064 42000 …near '`@localhost`' at line 1
> exec
| DROP USER 'u1' `@localhost`
< error
| 1064 42000 …near '`@localhost`' at line 1
> exec
| CREATE USER u1@'%' IDENTIFIED
< error
| 1064 42000 …near 'IDENTIFIED' at line 1
> exec
| GRANT SELECT ON *.* TO u1 IDENTIFIED
< error
| 1064 42000 …near 'IDENTIFIED' at line 1
> exec
| REVOKE SELECT ON *.* FROM u1 IDENTIFIED
< error
| 1064 42000 …near 'IDENTIFIED' at line 1
> exec
| GRANT SELECT ON *.* TO u1 IDENTIFIED BY
< error
| 1064 42000 …near 'BY' at line 1
> exec
| RENAME USER CURRENT_USER TO CURRENT_USER
< error
| 1064 42000 …near 'CURRENT_USER TO CURRENT_USER' at line 1
> exec
| RENAME USER u1@localhost TO CURRENT_USER
< error
| 1064 42000 …near 'CURRENT_USER' at line 1
> exec
| SET PASSWORD FOR CURRENT_USER = 'notarealpw'
< error
| 1064 42000 …near 'CURRENT_USER = 'notarealpw'' at line 1
> exec
| CREATE ROLE current_user
< error
| 1064 42000 …near 'current_user' at line 1
> exec
| SET ROLE current_user
< error
| 1064 42000 …near 'current_user' at line 1
> exec
| SET DEFAULT ROLE current_user TO u1@localhost
< error
| 1064 42000 …near 'current_user TO u1@localhost' at line 1
> exec
| CREATE USER u1@'%' PASSWORD EXPIRE DEFAULT ACCOUNT UNLOCK
< error
| 1064 42000 …near 'PASSWORD EXPIRE DEFAULT ACCOUNT UNLOCK' at line 1
> exec
| CREATE USER u1@'%' ATTRIBUTE '{"email": "user@example.com"}'
< error
| 1064 42000 …near 'ATTRIBUTE '{"email": "user@example.com"}'' at line 1
> exec
| GRANT SELECT TO u1
< error
| 1064 42000 …near 'SELECT TO u1' at line 1
> exec
| GRANT usage TO u1
< error
| 1064 42000 …near 'usage TO u1' at line 1
> exec
| GRANT process TO u1
< error
| 1064 42000 …near 'process TO u1' at line 1
