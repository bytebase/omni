-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12
-- pin: pg18 postgres:18.7-alpine sha256:75dcac3dfb1cd01ca3536910bbc84128759362f8b118bf8c55b38a94f68f3b73 REL_18_7

== TestQueryMergeCTEDoesNotLeakRows/default default#1
> exec
|
| 				CREATE TABLE secret_doc (id int PRIMARY KEY, secret text);
| 				INSERT INTO secret_doc VALUES (1, 'TOP-SECRET-SSN-123-45-6789');
| 			
< ok 1 -
> query
| SELECT secret FROM secret_doc LIMIT 5000;
< result
< column "secret" "TEXT"
< row
: string
| TOP-SECRET-SSN-123-45-6789
> exec
| WITH x AS (
| 				MERGE INTO secret_doc t USING secret_doc s ON t.id = s.id
| 				WHEN MATCHED THEN UPDATE SET id = t.id
| 				RETURNING t.secret AS secret
| 			) SELECT secret FROM x LIMIT 5000;
< ok 1 -

== TestQueryMergeCTEDoesNotLeakRows/pg18 pg18#1
> exec
|
| 				CREATE TABLE secret_doc (id int PRIMARY KEY, secret text);
| 				INSERT INTO secret_doc VALUES (1, 'TOP-SECRET-SSN-123-45-6789');
| 			
< ok 1 -
> query
| SELECT secret FROM secret_doc LIMIT 5000;
< result
< column "secret" "TEXT"
< row
: string
| TOP-SECRET-SSN-123-45-6789
> exec
| WITH x AS (
| 				MERGE INTO secret_doc t USING secret_doc s ON t.id = s.id
| 				WHEN MATCHED THEN UPDATE SET id = t.id
| 				RETURNING t.secret AS secret
| 			) SELECT secret FROM x LIMIT 5000;
< ok 1 -
