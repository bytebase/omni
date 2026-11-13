-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_5_3_UserRoleManagement default#1
> exec
| CREATE USER 'testuser53'@'localhost' IDENTIFIED BY 'Password123!'
< ok 0 0
> exec
| DROP USER IF EXISTS 'testuser53'@'localhost'
< ok 0 0
> exec
| CREATE USER IF NOT EXISTS 'dropme53'@'localhost' IDENTIFIED BY 'Password123!'
< ok 0 0
> exec
| DROP USER 'dropme53'@'localhost'
< ok 0 0
> exec
| CREATE USER IF NOT EXISTS 'alterme53'@'localhost' IDENTIFIED BY 'Password123!'
< ok 0 0
> exec
| ALTER USER 'alterme53'@'localhost' IDENTIFIED BY 'NewPassword456!'
< ok 0 0
> exec
| DROP USER IF EXISTS 'alterme53'@'localhost'
< ok 0 0
> exec
| CREATE ROLE 'app_role53'
< ok 0 0
> exec
| DROP ROLE IF EXISTS 'app_role53'
< ok 0 0
> exec
| CREATE ROLE IF NOT EXISTS 'droprole53'
< ok 0 0
> exec
| DROP ROLE 'droprole53'
< ok 0 0
> exec
| CREATE USER IF NOT EXISTS 'grantee53'@'localhost' IDENTIFIED BY 'Password123!'
< ok 0 0
> exec
| GRANT SELECT, INSERT ON test.* TO 'grantee53'@'localhost'
< ok 0 0
> exec
| DROP USER IF EXISTS 'grantee53'@'localhost'
< ok 0 0
> exec
| CREATE USER IF NOT EXISTS 'revokee53'@'localhost' IDENTIFIED BY 'Password123!'
< ok 0 0
> exec
| GRANT SELECT ON test.* TO 'revokee53'@'localhost'
< ok 0 0
> exec
| REVOKE SELECT ON test.* FROM 'revokee53'@'localhost'
< ok 0 0
> exec
| DROP USER IF EXISTS 'revokee53'@'localhost'
< ok 0 0
