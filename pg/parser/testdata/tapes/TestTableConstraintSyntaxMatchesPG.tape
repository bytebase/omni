-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestTableConstraintSyntaxMatchesPG default#1
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, check int)
< error
| 42601 syntax error at or near "int"
| position: 43
> exec
| SELECT 1/0;
| CREATE TABLE t (check int)
< error
| 42601 syntax error at or near "int"
| position: 35
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, primary int)
< error
| 42601 syntax error at or near "int"
| position: 45
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, unique int)
< error
| 42601 syntax error at or near "int"
| position: 44
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, foreign int)
< error
| 42601 syntax error at or near "int"
| position: 45
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, constraint int)
< error
| 42601 syntax error at or near ")"
| position: 51
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, select int)
< error
| 42601 syntax error at or near "select"
| position: 37
> exec
| SELECT 1/0;
| CREATE FOREIGN TABLE t (id int, check int) SERVER s
< error
| 42601 syntax error at or near "int"
| position: 51
> exec
| SELECT 1/0;
| ALTER TABLE t ADD check int
< error
| 42601 syntax error at or near "int"
| position: 37
> exec
| SELECT 1/0;
| ALTER TABLE t ADD primary int
< error
| 42601 syntax error at or near "int"
| position: 39
> exec
| SELECT 1/0;
| ALTER TABLE t ADD constraint int
< error
| 42601 syntax error at end of input
| position: 45
> exec
| SELECT 1/0;
| ALTER TABLE t ADD COLUMN check int
< error
| 42601 syntax error at or near "check"
| position: 38
> exec
| SELECT 1/0;
| CREATE TABLE t OF ty (check WITH OPTIONS NOT NULL)
< error
| 42601 syntax error at or near "WITH"
| position: 41
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, "check" int)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ADD COLUMN "check" int
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (between int)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, like int)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (exclude int)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, exclude int)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ADD exclude int
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t OF ty (exclude WITH OPTIONS NOT NULL)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (exclude) AS SELECT 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, exclude)
< error
| 42601 syntax error at or near ")"
| position: 43
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, exclude using gist (id with =))
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, exclude (id with =))
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ADD exclude (a WITH =)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, primary key (id), unique (id), foreign key (id) references u)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, constraint c check (id > 0))
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (id int CHECK id > 0)
< error
| 42601 syntax error at or near "id"
| position: 42
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, CHECK id > 0)
< error
| 42601 syntax error at or near "id"
| position: 43
> exec
| SELECT 1/0;
| CREATE TABLE t (id int CONSTRAINT c CHECK id > 0)
< error
| 42601 syntax error at or near "id"
| position: 55
> exec
| SELECT 1/0;
| CREATE TABLE t (id int DEFAULT 1 CHECK id > 0)
< error
| 42601 syntax error at or near "id"
| position: 52
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, CHECK)
< error
| 42601 syntax error at or near ")"
| position: 42
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, CHECK ())
< error
| 42601 syntax error at or near ")"
| position: 44
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CHECK id > 0
< error
| 42601 syntax error at or near "id"
| position: 37
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CONSTRAINT c CHECK id > 0
< error
| 42601 syntax error at or near "id"
| position: 50
> exec
| SELECT 1/0;
| CREATE DOMAIN d AS int CHECK VALUE > 0
< error
| 42601 syntax error at or near "VALUE"
| position: 42
> exec
| SELECT 1/0;
| CREATE TABLE t (id int CHECK (id > 0))
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (id int, CHECK (id > 0))
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (id int CHECK (id > 0) NO INHERIT)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CONSTRAINT c CHECK (id > 0)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE DOMAIN d AS int CHECK (VALUE > 0)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (a int GENERATED ALWAYS AS a + 1 STORED)
< error
| 42601 syntax error at or near "a"
| position: 55
> exec
| SELECT 1/0;
| CREATE TABLE t (a int GENERATED ALWAYS AS (a + 1))
< error
| 42601 syntax error at or near ")"
| position: 62
> exec
| SELECT 1/0;
| CREATE TABLE t (a int CONSTRAINT c)
< error
| 42601 syntax error at or near ")"
| position: 47
> exec
| SELECT 1/0;
| CREATE TABLE t (a int CONSTRAINT c NOT 5)
< error
| 42601 syntax error at or near "5"
| position: 52
> exec
| SELECT 1/0;
| CREATE TABLE t (a int PRIMARY)
< error
| 42601 syntax error at or near ")"
| position: 42
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, CONSTRAINT c)
< error
| 42601 syntax error at or near ")"
| position: 48
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, CONSTRAINT CHECK (a > 0))
< error
| 42601 syntax error at or near "CHECK"
| position: 47
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CONSTRAINT c UNIQUE USING INDEX
< error
| 42601 syntax error at end of input
| position: 62
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES)
< error
| 42601 syntax error at or near ")"
| position: 62
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, FOREIGN KEY (a) u)
< error
| 42601 syntax error at or near "u"
| position: 52
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, EXCLUDE (a =))
< error
| 42601 syntax error at or near "="
| position: 47
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, EXCLUDE (a WITH =) WHERE )
< error
| 42601 syntax error at or near ")"
| position: 61
> exec
| SELECT 1/0;
| CREATE INDEX i ON t ((a + ))
< error
| 42601 syntax error at or near ")"
| position: 39
> exec
| SELECT 1/0;
| CREATE TABLE t (a int GENERATED ALWAYS AS (a + 1) STORED)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (a int GENERATED BY DEFAULT AS IDENTITY)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, UNIQUE USING INDEX i)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ADD CONSTRAINT c UNIQUE USING INDEX i
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, EXCLUDE (a WITH OPERATOR(pg_catalog.=)))
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (a text, EXCLUDE USING gist (a gist_trgm_ops WITH =) WHERE (a <> ''))
< error
| 22012 division by zero
