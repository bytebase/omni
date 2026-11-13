-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestAlterTableRenameSyntaxMatchesPG default#1
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME TO u
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE IF EXISTS t RENAME TO u
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME COLUMN a TO b
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE ONLY t RENAME a TO b
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME CONSTRAINT k TO l
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t SET SCHEMA s
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER FOREIGN TABLE ft RENAME TO u
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER FOREIGN TABLE IF EXISTS ft RENAME COLUMN a TO b
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER FOREIGN TABLE ft RENAME a TO b
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER FOREIGN TABLE ft SET SCHEMA s
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME CONSTRAINT k l
< error
| 42601 syntax error at or near "l"
| position: 47
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME a b
< error
| 42601 syntax error at or near "b"
| position: 36
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME TO
< error
| 42601 syntax error at end of input
| position: 36
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME COLUMN a TO
< error
| 42601 syntax error at end of input
| position: 45
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME CONSTRAINT k TO
< error
| 42601 syntax error at end of input
| position: 49
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME a TO
< error
| 42601 syntax error at end of input
| position: 38
> exec
| SELECT 1/0;
| ALTER TABLE t RENAME
< error
| 42601 syntax error at end of input
| position: 33
> exec
| SELECT 1/0;
| ALTER TABLE t SET SCHEMA
< error
| 42601 syntax error at end of input
| position: 37
> exec
| SELECT 1/0;
| ALTER TABLE RENAME TO u
< error
| 42601 syntax error at or near "TO"
| position: 32
> exec
| SELECT 1/0;
| ALTER FOREIGN TABLE ft RENAME CONSTRAINT k TO l
< error
| 42601 syntax error at or near "CONSTRAINT"
| position: 43
> exec
| SELECT 1/0;
| ALTER FOREIGN TABLE ft RENAME TO
< error
| 42601 syntax error at end of input
| position: 45
> exec
| SELECT 1/0;
| ALTER FOREIGN TABLE ft RENAME COLUMN a TO
< error
| 42601 syntax error at end of input
| position: 54
> exec
| SELECT 1/0;
| ALTER FOREIGN TABLE ft SET SCHEMA
< error
| 42601 syntax error at end of input
| position: 46
> exec
| SELECT 1/0;
| ALTER FOREIGN TABLE RENAME TO u
< error
| 42601 syntax error at or near "TO"
| position: 40
