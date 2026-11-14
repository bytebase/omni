-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestCompoundStatementFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| BEGIN NOT ATOMIC END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC SELECT 1; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a INT DEFAULT 1; SET a = 2; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC BEGIN NOT ATOMIC END; END
< ok 0 0
> exec ^
: string
| BEGIN ATOMIC END
< error
| 1064 42000 …near 'ATOMIC END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC END a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC COMMIT; DECLARE a INT; END
< error
| 1064 42000 …near 'DECLARE a INT; END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC RETURN 1; END
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| BEGIN NOT ATOMIC RETURN 1; END
| ; SELECT FROM
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| EXPLAIN BEGIN NOT ATOMIC END
< error
| 1064 42000 …near 'NOT ATOMIC END' at line 1
> exec ^
: string
| SET STATEMENT max_statement_time=1 FOR BEGIN NOT ATOMIC DECLARE a INT; END
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() BEGIN NOT ATOMIC END
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() x: BEGIN NOT ATOMIC END
< error
| 1064 42000 …near 'NOT ATOMIC END' at line 1
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN NOT ATOMIC RETURN 1; END
< ok 0 0
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO BEGIN NOT ATOMIC DECLARE a INT; END
< ok 0 0
> exec ^
: string
| LOOP SELECT 1; END LOOP
< ok 0 0
> exec ^
: string
| LOOP END LOOP
< error
| 1064 42000 …near 'END LOOP' at line 1
> exec ^
: string
| a: LOOP COMMIT; END LOOP
< error
| 1064 42000 …near 'a: LOOP COMMIT; END LOOP' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC WHILE 1 DO END WHILE; END
< error
| 1064 42000 …near 'END WHILE; END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC IF 1 THEN COMMIT; ELSE END IF; END
< error
| 1064 42000 …near 'END IF; END' at line 1
> exec ^
: string
| SET STATEMENT max_statement_time=1 FOR LOOP COMMIT; END LOOP
< ok 0 0
> exec ^
: string
| FOR a IN (SELECT 1) DO COMMIT; END FOR
< ok 0 0
> exec ^
: string
| FOR a IN REVERSE (SELECT 1) DO COMMIT; END FOR
< ok 0 0
> exec ^
: string
| FOR a IN (WITH x AS (SELECT 1) SELECT * FROM x) DO COMMIT; END FOR
< ok 0 0
> exec ^
: string
| FOR a IN (SELECT 1 UNION SELECT 2) DO COMMIT; END FOR
< ok 0 0
> exec ^
: string
| FOR a IN (SELECT 1) DO END FOR
< error
| 1064 42000 …near 'END FOR' at line 1
> exec ^
: string
| FOR a IN (SELECT 1) DO COMMIT; END FOR a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| FOR a IN (SELECT 1 INTO @x) DO COMMIT; END FOR
< error
| 1323 42000 Cursor SELECT must not have INTO
> exec ^
: string
| FOR a IN (SELECT 1 INTO @x) DO COMMIT; END FOR
| ; SELECT FROM
< error
| 1323 42000 Cursor SELECT must not have INTO
> exec ^
: string
| BEGIN NOT ATOMIC x: FOR a IN (SELECT 1) DO LEAVE x; END FOR x; END
< ok 0 0
> exec ^
: string
| FOR r IN (SELECT 1 AS x) DO BEGIN DECLARE y TYPE OF r; END; END FOR
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() FOR a IN (SELECT 1) DO COMMIT; END FOR
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT FOR r IN (SELECT 1) DO RETURN 1; END FOR
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() DETERMINISTIC FOR r IN (SELECT 1) DO SELECT 1; END FOR
< ok 0 0
> exec ^
: string
| BEGIN NOT /* c */ ATOMIC IF 1 THEN SET @a=1; END IF; END
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SET STATEMENT max_statement_time=1 FOR LOOP RETURN 1; END LOOP; END
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SET STATEMENT max_statement_time=1 FOR LOOP SELECT 1; END LOOP; RETURN 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SET STATEMENT max_statement_time=1 FOR LOOP SELECT 1; END LOOP; RETURN 1; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE a INT; SET STATEMENT max_statement_time=1 FOR BEGIN NOT ATOMIC DECLARE b TYPE OF a; END; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a, b INT; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a, b, A INT; END
< error
| 1331 42000 Duplicate variable: A
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a, b, A INT; END
| ; SELECT FROM
< error
| 1331 42000 Duplicate variable: A
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a INT; DECLARE a INT; END
< error
| 1331 42000 Duplicate variable: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a INT; DECLARE a INT; END
| ; SELECT FROM
< error
| 1331 42000 Duplicate variable: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a INT; BEGIN DECLARE a INT; END; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a INT; DECLARE b TYPE OF A; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a TYPE OF a; END
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a TYPE OF a; END
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| BEGIN NOT ATOMIC BEGIN DECLARE a INT; END; BEGIN DECLARE b TYPE OF a; END; END
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| BEGIN NOT ATOMIC BEGIN DECLARE a INT; END; BEGIN DECLARE b TYPE OF a; END; END
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE b TYPE OF t.c DEFAULT 1; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE b TYPE OF d.t.c; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a ROW TYPE OF t; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a ROW TYPE OF d.t.c; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a ROW(b INT, c VARCHAR(10) CHARACTER SET latin1) DEFAULT ROW(1, 'x'); END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a ROW(b INT, B INT); END
< error
| 1060 42S21 Duplicate column name 'B'
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a ROW(b INT, B INT); END
| ; SELECT FROM
< error
| 1060 42S21 Duplicate column name 'B'
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a ROW(); END
< error
| 1064 42000 …near '); END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a ROW(b ROW(c INT)); END
< error
| 1064 42000 …near 'ROW(c INT)); END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a ROW(b INT NOT NULL); END
< error
| 1064 42000 …near 'NOT NULL); END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE `row` INT; DECLARE type INT; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CONDITION FOR 1; DECLARE a CONDITION FOR 2; END
< error
| 1332 42000 Duplicate condition: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CONDITION FOR 1; DECLARE a CONDITION FOR 2; END
| ; SELECT FROM
< error
| 1332 42000 Duplicate condition: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CONDITION FOR 1; BEGIN DECLARE a CONDITION FOR 2; END; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CONDITION FOR SQLSTATE VALUE '01000'; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CONDITION FOR NOT FOUND; END
< error
| 1064 42000 …near 'NOT FOUND; END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CONDITION FOR b; END
< error
| 1064 42000 …near 'b; END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a INT; DECLARE a CONDITION FOR 1; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR a COMMIT; END
< error
| 1319 42000 Undefined CONDITION: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR a COMMIT; END
| ; SELECT FROM
< error
| 1319 42000 Undefined CONDITION: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE c CONDITION FOR 1; DECLARE EXIT HANDLER FOR C COMMIT; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE c CONDITION FOR 1; BEGIN DECLARE EXIT HANDLER FOR c COMMIT; END; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR 1 BEGIN DECLARE c CONDITION FOR 1; END; DECLARE EXIT HANDLER FOR c COMMIT; END
< error
| 1319 42000 Undefined CONDITION: c
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR 1 BEGIN DECLARE c CONDITION FOR 1; END; DECLARE EXIT HANDLER FOR c COMMIT; END
| ; SELECT FROM
< error
| 1319 42000 Undefined CONDITION: c
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR 1, 1 COMMIT; END
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR 1, 1 COMMIT; END
| ; SELECT FROM
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR 1062 COMMIT; DECLARE CONTINUE HANDLER FOR 01062 COMMIT; END
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR 1062 COMMIT; DECLARE CONTINUE HANDLER FOR 01062 COMMIT; END
| ; SELECT FROM
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE c CONDITION FOR 1; DECLARE EXIT HANDLER FOR c, 1 COMMIT; END
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE c CONDITION FOR 1; DECLARE EXIT HANDLER FOR c, 1 COMMIT; END
| ; SELECT FROM
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE c CONDITION FOR SQLSTATE '01000'; DECLARE EXIT HANDLER FOR c, SQLSTATE '01000' COMMIT; END
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE c CONDITION FOR SQLSTATE '01000'; DECLARE EXIT HANDLER FOR c, SQLSTATE '01000' COMMIT; END
| ; SELECT FROM
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR SQLSTATE '01000', SQLWARNING COMMIT; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR NOT FOUND, NOT FOUND COMMIT; END
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE EXIT HANDLER FOR NOT FOUND, NOT FOUND COMMIT; END
| ; SELECT FROM
< error
| 1413 42000 Duplicate handler declared in the same block
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE c CONDITION FOR SQLSTATE '42000'; DECLARE EXIT HANDLER FOR c COMMIT; BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '42000' COMMIT; END; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR FOR SELECT 1; DECLARE a CURSOR FOR SELECT 1; END
< error
| 1333 42000 Duplicate cursor: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR FOR SELECT 1; DECLARE a CURSOR FOR SELECT 1; END
| ; SELECT FROM
< error
| 1333 42000 Duplicate cursor: a
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (b INT, c INT) FOR SELECT b; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (IN b INT) FOR SELECT b; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (OUT b INT) FOR SELECT b; END
< error
| 1235 42000 This version of MariaDB doesn't yet support 'OUT/INOUT cursor parameter'
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (OUT b INT) FOR SELECT b; END
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'OUT/INOUT cursor parameter'
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (b INT DEFAULT 1) FOR SELECT b; END
< error
| 1064 42000 …near 'DEFAULT 1) FOR SELECT b; END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR () FOR SELECT 1; END
< error
| 1064 42000 …near ') FOR SELECT 1; END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (x INT, X INT) FOR SELECT 1; END
< error
| 1330 42000 Duplicate parameter: X
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (x INT, X INT) FOR SELECT 1; END
| ; SELECT FROM
< error
| 1330 42000 Duplicate parameter: X
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (b ROW(x INT)) FOR SELECT 1; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (b TYPE OF t.c) FOR SELECT 1; END
< ok 0 0
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (b TYPE OF a) FOR SELECT 1; END
< error
| 1064 42000 …near ') FOR SELECT 1; END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (b ROW TYPE OF a.b.c) FOR SELECT 1; END
< error
| 1064 42000 …near '.c) FOR SELECT 1; END' at line 1
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (x INT) FOR SELECT 1; DECLARE b TYPE OF x; END
< error
| 1327 42000 Undeclared variable: x
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR (x INT) FOR SELECT 1; DECLARE b TYPE OF x; END
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: x
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR FOR SELECT 1 INTO @x; END
< error
| 1323 42000 Cursor SELECT must not have INTO
> exec ^
: string
| BEGIN NOT ATOMIC DECLARE a CURSOR FOR SELECT 1 INTO @x; END
| ; SELECT FROM
< error
| 1323 42000 Cursor SELECT must not have INTO
> exec ^
: string
| CREATE PROCEDURE p(IN a INT, OUT b ROW(x INT), INOUT c TYPE OF t.c) BEGIN END
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p(z INT) BEGIN DECLARE a TYPE OF z; END
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p(z INT) BEGIN DECLARE z INT; END
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p(z INT, Z INT) BEGIN END
< error
| 1330 42000 Duplicate parameter: Z
> exec ^
: string
| CREATE PROCEDURE p(z INT, Z INT) BEGIN END
| ; SELECT FROM
< error
| 1330 42000 Duplicate parameter: Z
> exec ^
: string
| CREATE PROCEDURE p(a TYPE OF a) BEGIN END
< error
| 1064 42000 …near ') BEGIN END' at line 1
> exec ^
: string
| CREATE PROCEDURE p(a ROW TYPE OF d.t.c) BEGIN END
< error
| 1064 42000 …near '.c) BEGIN END' at line 1
> exec ^
: string
| CREATE FUNCTION f(a ROW TYPE OF t) RETURNS INT RETURN 1
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f(z INT, z INT) RETURNS INT RETURN 1
< error
| 1330 42000 Duplicate parameter: z
> exec ^
: string
| CREATE FUNCTION f(z INT, z INT) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1330 42000 Duplicate parameter: z
