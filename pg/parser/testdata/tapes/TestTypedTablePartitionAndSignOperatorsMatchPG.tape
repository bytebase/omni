-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestTypedTablePartitionAndSignOperatorsMatchPG default#1
> exec
| SELECT 1/0;
| CREATE TABLE t OF typ PARTITION BY RANGE (a)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE IF NOT EXISTS t OF typ (a WITH OPTIONS NOT NULL) PARTITION BY HASH (a) USING heap
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t OF typ PARTITION BY RANGE
< error
| 42601 syntax error at end of input
| position: 53
> exec
| SELECT 1/0;
| CREATE OPERATOR === (FUNCTION = f, COMMUTATOR = +)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TYPE t SET (a = -)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TYPE t SET (a = - 1)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TYPE t SET (a = + +)
< error
| 42601 syntax error at or near "+"
| position: 37
