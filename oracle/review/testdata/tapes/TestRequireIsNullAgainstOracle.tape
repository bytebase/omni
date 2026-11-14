-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestRequireIsNullAgainstOracle default#1
> exec
| CREATE TABLE system.pubt (id NUMBER, a NUMBER)
< ok 0 0
> exec
| GRANT SELECT ON system.pubt TO PUBLIC
< ok 0 0
> exec
| CREATE PUBLIC SYNONYM pubt FOR system.pubt
< ok 0 0
> exec
| BEGIN
|   FOR u IN (SELECT username FROM dba_users WHERE username IN ('SHOP', 'BBDATAARCHIVE')) LOOP
|     EXECUTE IMMEDIATE 'DROP USER ' || u.username || ' CASCADE';
|   END LOOP;
| END;
< ok 1 0
> exec
| CREATE USER shop IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE MATERIALIZED VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PROCEDURE, CREATE TRIGGER, CREATE ANY TABLE, COMMENT ANY TABLE TO shop
< ok 0 0
> exec
| CREATE USER bbdataarchive IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0

== TestRequireIsNullAgainstOracle default#2
> exec
| CREATE TABLE p (id NUMBER CONSTRAINT pk_p PRIMARY KEY, note VARCHAR2(20))
< ok 0 0
> exec
| CREATE TABLE t (id NUMBER CONSTRAINT pk_t PRIMARY KEY, code VARCHAR2(20) CONSTRAINT nn_code NOT NULL CONSTRAINT uk_code UNIQUE, p_id NUMBER CONSTRAINT fk_p REFERENCES p (id), n NUMBER CONSTRAINT ck_n CHECK (n > 0))
< ok 0 0
> exec
| CREATE INDEX ix_n ON t (n)
< ok 0 0
> exec
| CREATE TABLE nokey (a NUMBER)
< ok 0 0
> exec
| CREATE TABLE plain (id NUMBER CONSTRAINT pk_plain PRIMARY KEY, v NUMBER)
< ok 0 0
> exec
| CREATE UNIQUE INDEX ux_v ON plain (v)
< ok 0 0
> exec
| CREATE TABLE longt (id NUMBER CONSTRAINT pk_longt PRIMARY KEY, body LONG)
< ok 0 0
> exec
| CREATE TABLE "a'b" (id NUMBER PRIMARY KEY)
< ok 0 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a')
< ok 1 0
> exec
| INSERT INTO p VALUES (2, 'b')
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1)
< ok 1 0
> exec
| INSERT INTO t VALUES (2, 'y', 2, 2)
< ok 1 0
> exec
| INSERT INTO nokey VALUES (1)
< ok 1 0
> exec
| INSERT INTO plain VALUES (1, 1)
< ok 1 0
> exec
| INSERT INTO longt VALUES (1, 'x')
< ok 1 0
> exec
| INSERT INTO "a'b" VALUES (1)
< ok 1 0
> exec
| CREATE TABLE probe (n NUMBER, s VARCHAR2(10), c CHAR(3))
< ok 0 0
> exec
| INSERT INTO probe VALUES (NULL, NULL, NULL)
< ok 1 0
> exec
| INSERT INTO probe VALUES (1, 'a', 'a')
< ok 1 0
> exec
| INSERT INTO probe VALUES (0, ' ', ' ')
< ok 1 0
> query
| SELECT COUNT(*) FROM probe WHERE n = NULL
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (n = NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.n = NULL THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NULL = n
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (NULL = n)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF NULL = r.n THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE n <> NULL
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (n <> NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.n <> NULL THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE n != NULL
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (n != NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.n != NULL THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE n ^= NULL
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (n ^= NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.n ^= NULL THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE (n) = (NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT ((n) = (NULL))
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF (r.n) = (NULL) THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE s = ''
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (s = '')
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.s = '' THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE s <> ''
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (s <> '')
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.s <> '' THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE c = ''
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (c = '')
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.c = '' THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE c != ''
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (c != '')
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.c != '' THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE '' = ''
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT ('' = '')
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF '' = '' THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE s = N''
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (s = N'')
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.s = N'' THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE s = q'[]'
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (s = q'[]')
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.s = q'[]' THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NULL = NULL
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (NULL = NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF NULL = NULL THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE n*n = NULL
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (n*n = NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.n*r.n = NULL THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE n + 1 <> NULL
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE NOT (n + 1 <> NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 0
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.n + 1 <> NULL THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 0
> query
| SELECT COUNT(*) FROM probe WHERE n = 1
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE NOT (n = 1)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.n = 1 THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE s = ' '
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE NOT (s = ' ')
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.s = ' ' THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE c = ' '
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE NOT (c = ' ')
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.c = ' ' THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE n IS NULL
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE NOT (n IS NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 2
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.n IS NULL THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE s IS NULL
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE NOT (s IS NULL)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 2
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.s IS NULL THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE n <> 1
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| SELECT COUNT(*) FROM probe WHERE NOT (n <> 1)
< result
< column "COUNT(*)" "NUMBER"
< row
: string
| 1
> query
| WITH FUNCTION k RETURN NUMBER IS
|   cnt NUMBER := 0;
| BEGIN
|   FOR r IN (SELECT * FROM probe) LOOP
|     IF r.n <> 1 THEN cnt := cnt + 1; END IF;
|   END LOOP;
|   RETURN cnt;
| END;
| SELECT k FROM dual
< result
< column "K" "NUMBER"
< row
: string
| 1
