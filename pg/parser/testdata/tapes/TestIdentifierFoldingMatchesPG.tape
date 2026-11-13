-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestIdentifierFoldingMatchesPG default#1
> query
| SELECT 1 AS MixedCase
< result
< column "mixedcase" "INT4"
< row
: int64 1
> query
| SELECT 1 AS Äb
< result
< column "Äb" "INT4"
< row
: int64 1
> query
| SELECT 1 AS ÉCOLE
< result
< column "École" "INT4"
< row
: int64 1
> query
| SELECT 1 AS ǅx
< result
< column "ǅx" "INT4"
< row
: int64 1
> query
| SELECT 1 AS İSTANBUL
< result
< column "İstanbul" "INT4"
< row
: int64 1
> query
| SELECT 1 AS KEY
< result
< column "Key" "INT4"
< row
: int64 1
> query
| SELECT 1 AS ſELECT
< result
< column "ſelect" "INT4"
< row
: int64 1
> query
| SELECT 1 AS "MixedCase"
< result
< column "MixedCase" "INT4"
< row
: int64 1
> query
| SELECT 1 AS AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
< result
< column "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" "INT4"
< row
: int64 1
> query
| SELECT 1 AS aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaÄxyz
< result
< column "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" "INT4"
< row
: int64 1
> query
| SELECT 1 AS aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaÄxyz
< result
< column "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaÄ" "INT4"
< row
: int64 1
> query
| SELECT 1 AS aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa€x
< result
< column "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa€" "INT4"
< row
: int64 1
> query
| SELECT 1 AS aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa€x
< result
< column "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" "INT4"
< row
: int64 1
> query
| SELECT 1 AS aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa😀x
< result
< column "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" "INT4"
< row
: int64 1
> query
| SELECT 1 AS aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
< result
< column "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" "INT4"
< row
: int64 1
> query
| SELECT 1 AS "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBÉxyz"
< result
< column "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" "INT4"
< row
: int64 1
> query
| SELECT 1 AS "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
< result
< column "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" "INT4"
< row
: int64 1
> query
| SELECT 1 AS U&"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\00C9xyz"
< result
< column "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc" "INT4"
< row
: int64 1
> query
| SELECT 1 AS U&"ccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\00C9xyz"
< result
< column "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccÉ" "INT4"
< row
: int64 1
> query
| SELECT 1 AS U&"d!00C9" UESCAPE '!'
< result
< column "dÉ" "INT4"
< row
: int64 1
