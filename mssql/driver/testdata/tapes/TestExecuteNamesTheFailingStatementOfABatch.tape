-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestExecuteNamesTheFailingStatementOfABatch default#1
> exec
| CREATE TABLE b_items (id INT PRIMARY KEY, note NVARCHAR(64));
< ok 0 -
> begin
< ok
> exec
| INSERT INTO b_items VALUES (1, N'a');
< ok 1 -
> commit
< ok
> begin
< ok
> exec
| INSERT INTO b_items VALUES (2, N'b');
| INSERT INTO b_items VALUES (1, N'duplicate');
< error
| 2627 1 14 2 Violation of PRIMARY KEY constraint 'PK__b_items__standin1'. Cannot insert duplicate key in object 'dbo.b_items'. The duplicate key value is (1).
> rollback
< ok
