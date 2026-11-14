-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestPackageFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE PACKAGE pk END
< ok 0 0
> exec ^
: string
| CREATE OR REPLACE PACKAGE pk END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE IF NOT EXISTS pk END
< ok 0 0
> exec ^
: string
| CREATE OR REPLACE PACKAGE IF NOT EXISTS pk END
< error
| 1221 HY000 Incorrect usage of OR REPLACE and IF NOT EXISTS
> exec ^
: string
| CREATE OR REPLACE PACKAGE IF NOT EXISTS pk END
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of OR REPLACE and IF NOT EXISTS
> exec ^
: string
| CREATE DEFINER = root@localhost PACKAGE pk END
< ok 0 0
> exec ^
: string
| CREATE DEFINER = CURRENT_USER() PACKAGE pk END
< ok 0 0
> exec ^
: string
| CREATE DEFINER = root SQL SECURITY INVOKER PACKAGE pk END
< error
| 1064 42000 …near 'PACKAGE pk END' at line 1
> exec ^
: string
| CREATE PACKAGE d.pk END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE pk END pk
< error
| 1064 42000 …near 'pk' at line 1
> exec ^
: string
| CREATE PACKAGE `BODY` END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE body END
< error
| 1305 42000 PACKAGE test.END does not exist
> exec ^
: string
| CREATE PACKAGE body END
| ; SELECT FROM
< error
| 1305 42000 PACKAGE test.END does not exist
> exec ^
: string
| CREATE TEMPORARY PACKAGE pk END
< error
| 1064 42000 …near 'PACKAGE pk END' at line 1
> exec ^
: string
| CREATE AGGREGATE PACKAGE pk END
< error
| 1064 42000 …near 'PACKAGE pk END' at line 1
> exec ^
: string
| CREATE PACKAGE pk SQL SECURITY INVOKER COMMENT 'x' SQL SECURITY DEFINER END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE pk SQL SECURITY pk END
< error
| 1064 42000 …near 'pk END' at line 1
> exec ^
: string
| CREATE PACKAGE pk LANGUAGE SQL END
< error
| 1064 42000 …near 'LANGUAGE SQL END' at line 1
> exec ^
: string
| CREATE PACKAGE pk DETERMINISTIC END
< error
| 1064 42000 …near 'DETERMINISTIC END' at line 1
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p(IN x INT, OUT y INT) COMMENT 'c' SQL SECURITY INVOKER; FUNCTION f(x INT) RETURNS INT DETERMINISTIC NO SQL; END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p() NOT DETERMINISTIC MODIFIES SQL DATA; END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p(); FUNCTION p() RETURNS INT; END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p(); PROCEDURE P(); END
< error
| 1304 42000 PROCEDURE pk.P already exists
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p(); PROCEDURE P(); END
| ; SELECT FROM
< error
| 1304 42000 PROCEDURE pk.P already exists
> exec ^
: string
| CREATE PACKAGE pk FUNCTION f() RETURNS INT; FUNCTION f() RETURNS INT; END
< error
| 1304 42000 FUNCTION pk.f already exists
> exec ^
: string
| CREATE PACKAGE pk FUNCTION f() RETURNS INT; FUNCTION f() RETURNS INT; END
| ; SELECT FROM
< error
| 1304 42000 FUNCTION pk.f already exists
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p() BEGIN END; END
< error
| 1064 42000 …near 'BEGIN END; END' at line 1
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p; END
< error
| 1064 42000 …near '; END' at line 1
> exec ^
: string
| CREATE PACKAGE pk FUNCTION f RETURNS INT; END
< error
| 1064 42000 …near 'RETURNS INT; END' at line 1
> exec ^
: string
| CREATE PACKAGE pk FUNCTION f() RETURNS INT AGGREGATE; END
< error
| 1064 42000 …near 'AGGREGATE; END' at line 1
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE d.p(); END
< error
| 1064 42000 …near '.p(); END' at line 1
> exec ^
: string
| CREATE PACKAGE pk DECLARE x INT; END
< error
| 1064 42000 …near 'DECLARE x INT; END' at line 1
> exec ^
: string
| CREATE PACKAGE `package` PROCEDURE `end`(); END
< ok 0 0
> exec ^
: string
| CREATE /* c */ PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE DEFINER = 'a b'@'localhost' PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE OR /* x */ REPLACE DEFINER = CURRENT_USER ( ) PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE DEFINER = 用户 PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE DEFINER = a$b@localhost PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE /*! DEFINER = root */ PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE /*M!100000 OR REPLACE */ PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p() LANGUAGE SQL; END
< ok 0 0
> exec ^
: string
| CREATE /*! /*! DEFINER = root */ PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p(x INT, x INT); END
< error
| 1330 42000 Duplicate parameter: x
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p(x INT, x INT); END
| ; SELECT FROM
< error
| 1330 42000 Duplicate parameter: x
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p(a ROW(x INT), b TYPE OF t.c, c ROW TYPE OF t); END
< ok 0 0
> exec ^
: string
| CREATE PACKAGE pk FUNCTION f(a TYPE OF a) RETURNS INT; END
< error
| 1064 42000 …near ') RETURNS INT; END' at line 1
> exec ^
: string
| CREATE DEFINER=root@127.0.0.1 PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE DEFINER = CURRENT_USER( PACKAGE pk END
< error
| 1064 42000 …near 'PACKAGE pk END' at line 1
> exec ^
: string
| CREATE DEFINER PACKAGE pk END
< error
| 1064 42000 …near 'PACKAGE pk END' at line 1
> exec ^
: string
| CREATE DEFINER = PACKAGE pk END
< error
| 1064 42000 …near 'pk END' at line 1
> exec ^
: string
| CREATE DEFINER = root@ PACKAGE pk END
< ok 0 0
> exec ^
: string
| CREATE DEFINER=root@ PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE DEFINER=CURRENT_ROLE() PACKAGE pk PROCEDURE p(); END
< ok 0 0
> exec ^
: string
| CREATE DEFINER=root@package VIEW v AS SELECT 1
< ok 0 0
> exec ^
: string
| CREATE DEFINER = root@ 'h' VIEW v AS SELECT 1
< error
| 1064 42000 …near ''h' VIEW v AS SELECT 1' at line 1
> exec ^
: string
| CREATE DEFINER = 'root'@ VIEW v AS SELECT 1
< ok 0 0
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE `p.q`(); END
< error
| 1458 42000 Incorrect routine name 'p.q'
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE `p.q`(); END
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'p.q'
> exec ^
: string
| CREATE PACKAGE pk FUNCTION `f.g`() RETURNS INT; END
< error
| 1458 42000 Incorrect routine name 'f.g'
> exec ^
: string
| CREATE PACKAGE pk FUNCTION `f.g`() RETURNS INT; END
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'f.g'
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE ``(); END
< error
| 1458 42000 Incorrect routine name ''
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE ``(); END
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name ''
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE `p `(); END
< error
| 1458 42000 Incorrect routine name 'p '
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE `p `(); END
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'p '
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`(); END
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`(); END
| ; SELECT FROM
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec ^
: string
| CREATE PACKAGE pk FUNCTION f(x INT, X INT) RETURNS INT; END
< error
| 1330 42000 Duplicate parameter: X
> exec ^
: string
| CREATE PACKAGE pk FUNCTION f(x INT, X INT) RETURNS INT; END
| ; SELECT FROM
< error
| 1330 42000 Duplicate parameter: X
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p() CONTAINS; END
< error
| 1064 42000 …near '; END' at line 1
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p() READS DATA; END
< error
| 1064 42000 …near 'DATA; END' at line 1
> exec ^
: string
| CREATE PACKAGE pk PROCEDURE p() SQL SECURITY garbage; END
< error
| 1064 42000 …near 'garbage; END' at line 1
> exec ^
: string
| DROP PACKAGE pk
< ok 0 0
> exec ^
: string
| DROP PACKAGE IF EXISTS d.pk
< ok 0 0
> exec ^
: string
| DROP PACKAGE BODY IF EXISTS pk
< ok 0 0
> exec ^
: string
| DROP PACKAGE BODY BODY
< ok 0 0
> exec ^
: string
| DROP PACKAGE BODY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| DROP PACKAGE IF EXISTS
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| DROP TEMPORARY PACKAGE pk
< error
| 1064 42000 …near 'PACKAGE pk' at line 1
> exec ^
: string
| SHOW CREATE PACKAGE pk
< ok 0 0
> exec ^
: string
| SHOW CREATE PACKAGE BODY pk
< ok 0 0
> exec ^
: string
| SHOW CREATE PACKAGE BODY
< ok 0 0
> exec ^
: string
| SHOW CREATE PACKAGE BODY.pk
< ok 0 0
> exec ^
: string
| SHOW CREATE PACKAGE BODY d.pk
< ok 0 0
> exec ^
: string
| SHOW PACKAGE STATUS
< ok 0 0
> exec ^
: string
| SHOW PACKAGE STATUS LIKE 'p%'
< ok 0 0
> exec ^
: string
| SHOW PACKAGE BODY STATUS WHERE Db = 'd'
< ok 0 0
> exec ^
: string
| SHOW PACKAGE BODY CODE pk
< error
| 1289 HY000 The 'SHOW PROCEDURE|FUNCTION CODE' feature is disabled; you need MariaDB built with '--with-debug' to have it working
> exec ^
: string
| SHOW PACKAGE CODE pk
< error
| 1064 42000 …near 'CODE pk' at line 1
> exec ^
: string
| SHOW PACKAGE BODY CODE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW PACKAGE
< error
| 1064 42000 …near '' at line 1
