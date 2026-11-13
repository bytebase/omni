-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestDML_BranchDivergence default#1
> exec
| INSERT INTO cities @ audit VALUES (1, 'San Francisco')
< error
| BRANCH_NOT_FOUND 134 1:1
| line 1:1: Branch 'audit' does not exist
> exec
| DELETE FROM orders @ audit
< error
| BRANCH_NOT_FOUND 134 1:1
| line 1:1: Branch 'audit' does not exist
> exec
| DELETE FROM orders @ audit WHERE x = 1
< error
| BRANCH_NOT_FOUND 134 1:1
| line 1:1: Branch 'audit' does not exist
> exec
| UPDATE purchases @ audit SET status = 'OVERDUE' WHERE ship_date IS NULL
< error
| BRANCH_NOT_FOUND 134 1:1
| line 1:1: Branch 'audit' does not exist
> exec
| MERGE INTO accounts @ audit t USING src s ON t.c = s.c WHEN MATCHED THEN DELETE
< error
| BRANCH_NOT_FOUND 134 1:1
| line 1:1: Branch 'audit' does not exist
