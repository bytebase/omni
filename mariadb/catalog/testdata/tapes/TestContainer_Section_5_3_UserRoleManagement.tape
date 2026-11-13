-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| GRANT SELECT, INSERT ON catalog_shared.* TO 'grantee53'@'localhost'
< ok 0 0
> exec
| DROP USER IF EXISTS 'grantee53'@'localhost'
< ok 0 0
> exec
| CREATE USER IF NOT EXISTS 'revokee53'@'localhost' IDENTIFIED BY 'Password123!'
< ok 0 0
> exec
| GRANT SELECT ON catalog_shared.* TO 'revokee53'@'localhost'
< ok 0 0
> exec
| REVOKE SELECT ON catalog_shared.* FROM 'revokee53'@'localhost'
< ok 0 0
> exec
| DROP USER IF EXISTS 'revokee53'@'localhost'
< ok 0 0
