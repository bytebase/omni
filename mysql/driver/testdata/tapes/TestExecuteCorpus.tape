-- pin: lctn1 mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46
-- pin: mysql84 mysql:8.4.13 sha256:a81d043930e51745522bd5c430961cf8286f0a23f0a9a2a6e0108f7ea71b0902 mysql-8.4.13

== TestExecuteCorpus/lctn1 lctn1#1
> begin
< ok
> exec
| /*app=bytebase*/ CREATE TABLE eitems (id INT PRIMARY KEY, note VARCHAR(64));
< ok 0 0
> exec
| /*app=bytebase*/ 
| INSERT INTO eitems VALUES (1, 'a'), (2, 'b'), (3, 'c');
< ok 3 0
> exec
| /*app=bytebase*/ 
| UPDATE eitems SET note = 'z' WHERE id <= 2;
< ok 2 0
> commit
< ok
> begin
< ok
> exec
| /*app=bytebase*/ INSERT INTO eitems VALUES (10, 'rb');
< ok 1 0
> exec
| /*app=bytebase*/ 
| INSERT INTO eitems VALUES (1, 'dup');
< error
| 1062 23000 Duplicate entry '1' for key 'eitems.PRIMARY'
> rollback
< ok
> exec
| /*app=bytebase*/ INSERT INTO eitems VALUES (20, 'kept');
< ok 1 0
> exec
| /*app=bytebase*/ 
| INSERT INTO eitems VALUES (1, 'dup');
< error
| 1062 23000 Duplicate entry '1' for key 'eitems.PRIMARY'
> begin
< ok
> exec
| /*app=bytebase*/ CREATE TABLE emixed (id INT PRIMARY KEY);
< ok 0 0
> exec
| /*app=bytebase*/ 
| INSERT INTO emixed VALUES (1);
< ok 1 0
> exec ^
< error
| 1062 23000 Duplicate entry '1' for key 'emixed.PRIMARY'
> rollback
< ok
> exec
| /*app=bytebase*/ CREATE TABLE ecreated (id INT PRIMARY KEY);
< ok 0 0
> exec
| /*app=bytebase*/ 
| INSERT INTO ecreated VALUES (1);
< ok 1 0
> exec
| /*app=bytebase*/ CREATE PROCEDURE ecounter(OUT n INT)
| BEGIN
|   SELECT COUNT(*) INTO n FROM eitems;
| END
< ok 0 0
> query
| /*app=bytebase*/ SELECT
| 				(SELECT COUNT(*) FROM eitems),
| 				(SELECT COUNT(*) FROM eitems WHERE id = 10),
| 				(SELECT COUNT(*) FROM eitems WHERE id = 20),
| 				(SELECT COUNT(*) FROM emixed),
| 				(SELECT COUNT(*) FROM ecreated),
| 				(SELECT COUNT(*) FROM information_schema.routines WHERE routine_schema = DATABASE() AND routine_name = 'ecounter');
< result
< column "(SELECT COUNT(*) FROM eitems)" "BIGINT"
< column "(SELECT COUNT(*) FROM eitems WHERE id = 10)" "BIGINT"
< column "(SELECT COUNT(*) FROM eitems WHERE id = 20)" "BIGINT"
< column "(SELECT COUNT(*) FROM emixed)" "BIGINT"
< column "(SELECT COUNT(*) FROM ecreated)" "BIGINT"
< column "(SELECT COUNT(*) FROM information_schema.routines WHERE routine_schema = DATABASE() AND routine_name = 'ecounter')" "BIGINT"
< row
: int64 4
: int64 0
: int64 1
: int64 1
: int64 1
: int64 1

== TestExecuteCorpus/mysql84 mysql84#1
> begin
< ok
> exec
| /*app=bytebase*/ CREATE TABLE eitems (id INT PRIMARY KEY, note VARCHAR(64));
< ok 0 0
> exec
| /*app=bytebase*/ 
| INSERT INTO eitems VALUES (1, 'a'), (2, 'b'), (3, 'c');
< ok 3 0
> exec
| /*app=bytebase*/ 
| UPDATE eitems SET note = 'z' WHERE id <= 2;
< ok 2 0
> commit
< ok
> begin
< ok
> exec
| /*app=bytebase*/ INSERT INTO eitems VALUES (10, 'rb');
< ok 1 0
> exec
| /*app=bytebase*/ 
| INSERT INTO eitems VALUES (1, 'dup');
< error
| 1062 23000 Duplicate entry '1' for key 'eitems.PRIMARY'
> rollback
< ok
> exec
| /*app=bytebase*/ INSERT INTO eitems VALUES (20, 'kept');
< ok 1 0
> exec
| /*app=bytebase*/ 
| INSERT INTO eitems VALUES (1, 'dup');
< error
| 1062 23000 Duplicate entry '1' for key 'eitems.PRIMARY'
> begin
< ok
> exec
| /*app=bytebase*/ CREATE TABLE emixed (id INT PRIMARY KEY);
< ok 0 0
> exec
| /*app=bytebase*/ 
| INSERT INTO emixed VALUES (1);
< ok 1 0
> exec ^
< error
| 1062 23000 Duplicate entry '1' for key 'emixed.PRIMARY'
> rollback
< ok
> exec
| /*app=bytebase*/ CREATE TABLE ecreated (id INT PRIMARY KEY);
< ok 0 0
> exec
| /*app=bytebase*/ 
| INSERT INTO ecreated VALUES (1);
< ok 1 0
> exec
| /*app=bytebase*/ CREATE PROCEDURE ecounter(OUT n INT)
| BEGIN
|   SELECT COUNT(*) INTO n FROM eitems;
| END
< ok 0 0
> query
| /*app=bytebase*/ SELECT
| 				(SELECT COUNT(*) FROM eitems),
| 				(SELECT COUNT(*) FROM eitems WHERE id = 10),
| 				(SELECT COUNT(*) FROM eitems WHERE id = 20),
| 				(SELECT COUNT(*) FROM emixed),
| 				(SELECT COUNT(*) FROM ecreated),
| 				(SELECT COUNT(*) FROM information_schema.routines WHERE routine_schema = DATABASE() AND routine_name = 'ecounter');
< result
< column "(SELECT COUNT(*) FROM eitems)" "BIGINT"
< column "(SELECT COUNT(*) FROM eitems WHERE id = 10)" "BIGINT"
< column "(SELECT COUNT(*) FROM eitems WHERE id = 20)" "BIGINT"
< column "(SELECT COUNT(*) FROM emixed)" "BIGINT"
< column "(SELECT COUNT(*) FROM ecreated)" "BIGINT"
< column "(SELECT COUNT(*) FROM information_schema.routines WHERE routine_schema = DATABASE() AND routine_name = 'ecounter')" "BIGINT"
< row
: int64 4
: int64 0
: int64 1
: int64 1
: int64 1
: int64 1
