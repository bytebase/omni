-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestExecuteAutoCommitKeepsSessionState default#1
> exec
| CREATE TABLE s_items (id INT PRIMARY KEY);
< ok 0 -
> query
| CREATE TABLE #scratch (id INT);
| INSERT INTO #scratch VALUES (1), (2);
: message
< stream
< message nextresultset more
< message rowsaffected 2
< message nextresultset more
< message nextresultset end
> exec
| INSERT INTO s_items SELECT id FROM #scratch;
< ok 2 -
> query
| SELECT COUNT(*) FROM s_items;
< result
< column "" "INT"
< row
: int64 2
