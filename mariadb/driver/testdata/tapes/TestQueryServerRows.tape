-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestQueryServerRows default#1
> exec
| CREATE TABLE ritems (id INT PRIMARY KEY, d INT);
| INSERT INTO ritems VALUES (1, 0), (2, 0), (3, 1);
| CREATE PROCEDURE rcount() BEGIN SELECT COUNT(*) AS n FROM ritems; END
< ok 0 0
> query
| INSERT INTO ritems VALUES (4, 1) RETURNING id;
< result
< column "id" "INT" "int32"
< row
: int64 4
> query
| REPLACE INTO ritems VALUES (4, 2) RETURNING d;
< result
< column "d" "INT"
< row
: int64 2
> query
| DELETE FROM ritems WHERE d = 0 RETURNING id;
< result
< column "id" "INT" "int32"
< row
: int64 1
< row
: int64 2
> exec
| DELETE FROM ritems WHERE d = 1 RETURNING id;
< ok 0 0
> exec
| INSERT INTO ritems VALUES (5, 1) RETURNING id;
< ok 0 0
> query
| CALL rcount();
< result
< column "n" "BIGINT" "int64"
< row
: int64 2
> exec
| CALL rcount();
< ok 0 0
> query
| SELECT id, d FROM ritems ORDER BY id;
< result
< column "id" "INT" "int32"
< column "d" "INT"
< row
: int64 4
: int64 2
< row
: int64 5
: int64 1
> query ^
< result
< column "id" "INT" "int32"
< column "d" "INT"
< row
: int64 4
: int64 2
< row
: int64 5
: int64 1
> query
| CHECK TABLE ritems;
< result
< column "Table" "VARCHAR"
< column "Op" "VARCHAR"
< column "Msg_type" "VARCHAR"
< column "Msg_text" "MEDIUMTEXT"
< row
: bytes
| omni_rows_test.ritems
: bytes
| check
: bytes
| status
: bytes
| OK
> exec
| CREATE TEMPORARY TABLE leak (id INT);
< ok 0 0
> query
| SELECT * FROM leak;
< error
| 1146 42S02 Table 'omni_rows_test.leak' doesn't exist
