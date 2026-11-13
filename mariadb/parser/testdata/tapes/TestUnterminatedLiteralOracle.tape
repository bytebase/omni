-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestUnterminatedLiteralOracle default#1
> exec
| SELECT 'oops;
< error
| 1064 42000 …near ''oops' at line 1
> exec
| SELECT "oops;
< error
| 1064 42000 …near '"oops' at line 1
> exec
| SELECT _utf8'oops;
< error
| 1064 42000 …near ''oops' at line 1
> exec
| SELECT `oops;
< error
| 1064 42000 …near '`oops' at line 1
> exec
| SELECT @`oops;
< error
| 1064 42000 …near '`oops' at line 1
> exec
| SELECT 1 /* oops;
< error
| 1064 42000 …near '/* oops' at line 1
> exec
| /* oops
< error
| 1064 42000 …near '/* oops' at line 1
> exec
| SELECT 1 /*!40000 oops
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT 'oops';
< ok 0 0
> exec
| SELECT "oops";
< ok 0 0
> exec
| SELECT _utf8'oops';
< ok 0 0
> exec
| SELECT `oops`;
< error
| 1054 42S22 Unknown column 'oops' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `oops`;
< error
| 1054 42S22 Unknown column 'oops' in 'SELECT'
> exec ^
: string
| SELECT `oops`;
| ; SELECT FROM
< error
| 1064 42000 …near '; SELECT FROM' at line 2
> exec
| SELECT @`oops`;
< ok 0 0
> exec
| SELECT 1 /* oops */;
< ok 0 0
> exec
| /* oops */
< ok 0 0
> exec
| SELECT 1 /*!40000 + 1 */
< ok 0 0
