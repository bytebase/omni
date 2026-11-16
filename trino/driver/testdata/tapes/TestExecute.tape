-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestExecute/counts_and_steps default#1
> exec
| CREATE SCHEMA memory.omni_execute
< ok 0 -
> exec
| CREATE TABLE memory.omni_execute.items (id integer, note varchar)
< ok 0 -
> exec
| INSERT INTO memory.omni_execute.items VALUES (1, 'a')
< ok 1 -
> exec
| INSERT INTO memory.omni_execute.items VALUES (2, 'b')
< ok 1 -

== TestExecute/empty_script default#1

== TestExecute/failure_keeps_the_rows_before_it default#1
> exec
| CREATE SCHEMA memory.omni_execute_failure
< ok 0 -
> exec
| CREATE TABLE memory.omni_execute_failure.items (id integer)
< ok 0 -
> exec
| INSERT INTO memory.omni_execute_failure.items VALUES (1)
< ok 1 -
> exec
| INSERT INTO memory.omni_execute_failure.missing VALUES (2)
< error
| {"Status":200,"Type":"USER_ERROR","Message":"line 1:1: Table 'memory.omni_execute_failure.missing' does not exist","SqlState":"","Name":"TABLE_NOT_FOUND","Code":46,"Line":1,"Column":1}

== TestExecute/transaction_mode_is_ignored default#1
> exec
| CREATE SCHEMA memory.omni_execute_tx
< ok 0 -
> exec
| CREATE TABLE memory.omni_execute_tx.items (id integer)
< ok 0 -
> exec
| INSERT INTO memory.omni_execute_tx.items VALUES (1)
< ok 1 -
