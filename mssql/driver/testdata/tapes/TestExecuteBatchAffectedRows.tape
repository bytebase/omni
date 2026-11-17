-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestExecuteBatchAffectedRows default#1
> exec
| CREATE TABLE x_items (id INT PRIMARY KEY, note NVARCHAR(20));
< ok 0 -
> begin
< ok
> query
| INSERT INTO x_items VALUES (1, N'a'), (2, N'b'), (3, N'c');
| UPDATE x_items SET note = N'u' WHERE id <= 2;
| DELETE FROM x_items WHERE id = 3;
: message
< stream
< message rowsaffected 3
< message nextresultset more
< message rowsaffected 2
< message nextresultset more
< message rowsaffected 1
< message nextresultset more
< message nextresultset end
> commit
< ok
> begin
< ok
> query
| INSERT INTO x_items VALUES (10, N'x'), (11, N'y');
| UPDATE x_items SET note = N'z' WHERE id = 10;
| INSERT INTO x_items VALUES (10, N'duplicate');
: message
< stream
< message rowsaffected 2
< message nextresultset more
< message rowsaffected 1
< message nextresultset more
< message error
| 2627 1 14 3 Violation of PRIMARY KEY constraint 'PK__x_items__standin1'. Cannot insert duplicate key in object 'dbo.x_items'. The duplicate key value is (10).
< message notice
| The statement has been terminated.
< message nextresultset more
< message nextresultset end
> rollback
< ok
