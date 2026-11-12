-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
