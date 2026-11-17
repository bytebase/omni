-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestQueryBatchAffectedRows default#1
> exec
| CREATE TABLE r_items (id INT PRIMARY KEY, note NVARCHAR(20));
< ok 0 -
> query
| INSERT INTO r_items VALUES (1, N'a'), (2, N'b'), (3, N'c');
|
| UPDATE r_items SET note = N'u' WHERE id <= 2;
|
| DELETE FROM r_items WHERE id = 3;
|
: message
< stream
< message rowsaffected 3
< message nextresultset more
< message rowsaffected 2
< message nextresultset more
< message rowsaffected 1
< message nextresultset more
< message nextresultset end
> query
| INSERT INTO r_items VALUES (10, N'x'), (11, N'y');
|
| UPDATE r_items SET note = N'z' WHERE id = 10;
|
| INSERT INTO r_items VALUES (10, N'duplicate');
|
| DELETE FROM r_items WHERE id = 11;
|
: message
< stream
< message rowsaffected 2
< message nextresultset more
< message rowsaffected 1
< message nextresultset more
< message error
| 2627 1 14 5 Violation of PRIMARY KEY constraint 'PK__r_items__standin1'. Cannot insert duplicate key in object 'dbo.r_items'. The duplicate key value is (10).
< message notice
| The statement has been terminated.
< message nextresultset more
< message rowsaffected 1
< message nextresultset more
< message nextresultset end
> query
| SET NOCOUNT ON;
|
| INSERT INTO r_items VALUES (20, N'n'), (21, N'm');
|
| UPDATE r_items SET note = N'q' WHERE id >= 20;
|
: message
< stream
< message nextresultset more
< message nextresultset more
< message nextresultset more
< message nextresultset end
