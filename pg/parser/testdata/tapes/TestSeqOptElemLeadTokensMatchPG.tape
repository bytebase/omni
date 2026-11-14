-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestSeqOptElemLeadTokensMatchPG default#1
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq abort
< error
| 42601 syntax error at or near "abort"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq abort 1
< error
| 42601 syntax error at or near "abort"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq abort int
< error
| 42601 syntax error at or near "abort"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq abort BY NONE
< error
| 42601 syntax error at or near "abort"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq abort NAME omni_probe_name
< error
| 42601 syntax error at or near "abort"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq abort CYCLE
< error
| 42601 syntax error at or near "abort"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absent
< error
| 42601 syntax error at or near "absent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absent 1
< error
| 42601 syntax error at or near "absent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absent int
< error
| 42601 syntax error at or near "absent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absent BY NONE
< error
| 42601 syntax error at or near "absent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absent NAME omni_probe_name
< error
| 42601 syntax error at or near "absent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absent CYCLE
< error
| 42601 syntax error at or near "absent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absolute
< error
| 42601 syntax error at or near "absolute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absolute 1
< error
| 42601 syntax error at or near "absolute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absolute int
< error
| 42601 syntax error at or near "absolute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absolute BY NONE
< error
| 42601 syntax error at or near "absolute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absolute NAME omni_probe_name
< error
| 42601 syntax error at or near "absolute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq absolute CYCLE
< error
| 42601 syntax error at or near "absolute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq access
< error
| 42601 syntax error at or near "access"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq access 1
< error
| 42601 syntax error at or near "access"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq access int
< error
| 42601 syntax error at or near "access"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq access BY NONE
< error
| 42601 syntax error at or near "access"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq access NAME omni_probe_name
< error
| 42601 syntax error at or near "access"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq access CYCLE
< error
| 42601 syntax error at or near "access"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq action
< error
| 42601 syntax error at or near "action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq action 1
< error
| 42601 syntax error at or near "action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq action int
< error
| 42601 syntax error at or near "action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq action BY NONE
< error
| 42601 syntax error at or near "action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq action NAME omni_probe_name
< error
| 42601 syntax error at or near "action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq action CYCLE
< error
| 42601 syntax error at or near "action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq add
< error
| 42601 syntax error at or near "add"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq add 1
< error
| 42601 syntax error at or near "add"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq add int
< error
| 42601 syntax error at or near "add"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq add BY NONE
< error
| 42601 syntax error at or near "add"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq add NAME omni_probe_name
< error
| 42601 syntax error at or near "add"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq add CYCLE
< error
| 42601 syntax error at or near "add"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq admin
< error
| 42601 syntax error at or near "admin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq admin 1
< error
| 42601 syntax error at or near "admin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq admin int
< error
| 42601 syntax error at or near "admin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq admin BY NONE
< error
| 42601 syntax error at or near "admin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq admin NAME omni_probe_name
< error
| 42601 syntax error at or near "admin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq admin CYCLE
< error
| 42601 syntax error at or near "admin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq after
< error
| 42601 syntax error at or near "after"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq after 1
< error
| 42601 syntax error at or near "after"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq after int
< error
| 42601 syntax error at or near "after"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq after BY NONE
< error
| 42601 syntax error at or near "after"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq after NAME omni_probe_name
< error
| 42601 syntax error at or near "after"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq after CYCLE
< error
| 42601 syntax error at or near "after"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq aggregate
< error
| 42601 syntax error at or near "aggregate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq aggregate 1
< error
| 42601 syntax error at or near "aggregate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq aggregate int
< error
| 42601 syntax error at or near "aggregate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq aggregate BY NONE
< error
| 42601 syntax error at or near "aggregate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq aggregate NAME omni_probe_name
< error
| 42601 syntax error at or near "aggregate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq aggregate CYCLE
< error
| 42601 syntax error at or near "aggregate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq all
< error
| 42601 syntax error at or near "all"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq all 1
< error
| 42601 syntax error at or near "all"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq all int
< error
| 42601 syntax error at or near "all"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq all BY NONE
< error
| 42601 syntax error at or near "all"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq all NAME omni_probe_name
< error
| 42601 syntax error at or near "all"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq all CYCLE
< error
| 42601 syntax error at or near "all"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq also
< error
| 42601 syntax error at or near "also"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq also 1
< error
| 42601 syntax error at or near "also"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq also int
< error
| 42601 syntax error at or near "also"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq also BY NONE
< error
| 42601 syntax error at or near "also"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq also NAME omni_probe_name
< error
| 42601 syntax error at or near "also"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq also CYCLE
< error
| 42601 syntax error at or near "also"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq alter
< error
| 42601 syntax error at or near "alter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq alter 1
< error
| 42601 syntax error at or near "alter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq alter int
< error
| 42601 syntax error at or near "alter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq alter BY NONE
< error
| 42601 syntax error at or near "alter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq alter NAME omni_probe_name
< error
| 42601 syntax error at or near "alter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq alter CYCLE
< error
| 42601 syntax error at or near "alter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq always
< error
| 42601 syntax error at or near "always"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq always 1
< error
| 42601 syntax error at or near "always"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq always int
< error
| 42601 syntax error at or near "always"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq always BY NONE
< error
| 42601 syntax error at or near "always"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq always NAME omni_probe_name
< error
| 42601 syntax error at or near "always"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq always CYCLE
< error
| 42601 syntax error at or near "always"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyse
< error
| 42601 syntax error at or near "analyse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyse 1
< error
| 42601 syntax error at or near "analyse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyse int
< error
| 42601 syntax error at or near "analyse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyse BY NONE
< error
| 42601 syntax error at or near "analyse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyse NAME omni_probe_name
< error
| 42601 syntax error at or near "analyse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyse CYCLE
< error
| 42601 syntax error at or near "analyse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyze
< error
| 42601 syntax error at or near "analyze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyze 1
< error
| 42601 syntax error at or near "analyze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyze int
< error
| 42601 syntax error at or near "analyze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyze BY NONE
< error
| 42601 syntax error at or near "analyze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyze NAME omni_probe_name
< error
| 42601 syntax error at or near "analyze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq analyze CYCLE
< error
| 42601 syntax error at or near "analyze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq and
< error
| 42601 syntax error at or near "and"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq and 1
< error
| 42601 syntax error at or near "and"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq and int
< error
| 42601 syntax error at or near "and"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq and BY NONE
< error
| 42601 syntax error at or near "and"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq and NAME omni_probe_name
< error
| 42601 syntax error at or near "and"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq and CYCLE
< error
| 42601 syntax error at or near "and"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq any
< error
| 42601 syntax error at or near "any"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq any 1
< error
| 42601 syntax error at or near "any"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq any int
< error
| 42601 syntax error at or near "any"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq any BY NONE
< error
| 42601 syntax error at or near "any"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq any NAME omni_probe_name
< error
| 42601 syntax error at or near "any"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq any CYCLE
< error
| 42601 syntax error at or near "any"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq array
< error
| 42601 syntax error at or near "array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq array 1
< error
| 42601 syntax error at or near "array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq array int
< error
| 42601 syntax error at or near "array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq array BY NONE
< error
| 42601 syntax error at or near "array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq array NAME omni_probe_name
< error
| 42601 syntax error at or near "array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq array CYCLE
< error
| 42601 syntax error at or near "array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq as
< error
| 42601 syntax error at end of input
| position: 34
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq as 1
< error
| 42601 syntax error at or near "1"
| position: 35
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq as int
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asc
< error
| 42601 syntax error at or near "asc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asc 1
< error
| 42601 syntax error at or near "asc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asc int
< error
| 42601 syntax error at or near "asc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asc BY NONE
< error
| 42601 syntax error at or near "asc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asc NAME omni_probe_name
< error
| 42601 syntax error at or near "asc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asc CYCLE
< error
| 42601 syntax error at or near "asc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asensitive
< error
| 42601 syntax error at or near "asensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asensitive 1
< error
| 42601 syntax error at or near "asensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asensitive int
< error
| 42601 syntax error at or near "asensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asensitive BY NONE
< error
| 42601 syntax error at or near "asensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asensitive NAME omni_probe_name
< error
| 42601 syntax error at or near "asensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asensitive CYCLE
< error
| 42601 syntax error at or near "asensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assertion
< error
| 42601 syntax error at or near "assertion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assertion 1
< error
| 42601 syntax error at or near "assertion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assertion int
< error
| 42601 syntax error at or near "assertion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assertion BY NONE
< error
| 42601 syntax error at or near "assertion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assertion NAME omni_probe_name
< error
| 42601 syntax error at or near "assertion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assertion CYCLE
< error
| 42601 syntax error at or near "assertion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assignment
< error
| 42601 syntax error at or near "assignment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assignment 1
< error
| 42601 syntax error at or near "assignment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assignment int
< error
| 42601 syntax error at or near "assignment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assignment BY NONE
< error
| 42601 syntax error at or near "assignment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assignment NAME omni_probe_name
< error
| 42601 syntax error at or near "assignment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq assignment CYCLE
< error
| 42601 syntax error at or near "assignment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asymmetric
< error
| 42601 syntax error at or near "asymmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asymmetric 1
< error
| 42601 syntax error at or near "asymmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asymmetric int
< error
| 42601 syntax error at or near "asymmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asymmetric BY NONE
< error
| 42601 syntax error at or near "asymmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asymmetric NAME omni_probe_name
< error
| 42601 syntax error at or near "asymmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq asymmetric CYCLE
< error
| 42601 syntax error at or near "asymmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq at
< error
| 42601 syntax error at or near "at"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq at 1
< error
| 42601 syntax error at or near "at"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq at int
< error
| 42601 syntax error at or near "at"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq at BY NONE
< error
| 42601 syntax error at or near "at"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq at NAME omni_probe_name
< error
| 42601 syntax error at or near "at"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq at CYCLE
< error
| 42601 syntax error at or near "at"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq atomic
< error
| 42601 syntax error at or near "atomic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq atomic 1
< error
| 42601 syntax error at or near "atomic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq atomic int
< error
| 42601 syntax error at or near "atomic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq atomic BY NONE
< error
| 42601 syntax error at or near "atomic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq atomic NAME omni_probe_name
< error
| 42601 syntax error at or near "atomic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq atomic CYCLE
< error
| 42601 syntax error at or near "atomic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attach
< error
| 42601 syntax error at or near "attach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attach 1
< error
| 42601 syntax error at or near "attach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attach int
< error
| 42601 syntax error at or near "attach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attach BY NONE
< error
| 42601 syntax error at or near "attach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attach NAME omni_probe_name
< error
| 42601 syntax error at or near "attach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attach CYCLE
< error
| 42601 syntax error at or near "attach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attribute
< error
| 42601 syntax error at or near "attribute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attribute 1
< error
| 42601 syntax error at or near "attribute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attribute int
< error
| 42601 syntax error at or near "attribute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attribute BY NONE
< error
| 42601 syntax error at or near "attribute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attribute NAME omni_probe_name
< error
| 42601 syntax error at or near "attribute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq attribute CYCLE
< error
| 42601 syntax error at or near "attribute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq authorization
< error
| 42601 syntax error at or near "authorization"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq authorization 1
< error
| 42601 syntax error at or near "authorization"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq authorization int
< error
| 42601 syntax error at or near "authorization"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq authorization BY NONE
< error
| 42601 syntax error at or near "authorization"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq authorization NAME omni_probe_name
< error
| 42601 syntax error at or near "authorization"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq authorization CYCLE
< error
| 42601 syntax error at or near "authorization"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq backward
< error
| 42601 syntax error at or near "backward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq backward 1
< error
| 42601 syntax error at or near "backward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq backward int
< error
| 42601 syntax error at or near "backward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq backward BY NONE
< error
| 42601 syntax error at or near "backward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq backward NAME omni_probe_name
< error
| 42601 syntax error at or near "backward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq backward CYCLE
< error
| 42601 syntax error at or near "backward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq before
< error
| 42601 syntax error at or near "before"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq before 1
< error
| 42601 syntax error at or near "before"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq before int
< error
| 42601 syntax error at or near "before"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq before BY NONE
< error
| 42601 syntax error at or near "before"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq before NAME omni_probe_name
< error
| 42601 syntax error at or near "before"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq before CYCLE
< error
| 42601 syntax error at or near "before"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq begin
< error
| 42601 syntax error at or near "begin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq begin 1
< error
| 42601 syntax error at or near "begin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq begin int
< error
| 42601 syntax error at or near "begin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq begin BY NONE
< error
| 42601 syntax error at or near "begin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq begin NAME omni_probe_name
< error
| 42601 syntax error at or near "begin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq begin CYCLE
< error
| 42601 syntax error at or near "begin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq between
< error
| 42601 syntax error at or near "between"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq between 1
< error
| 42601 syntax error at or near "between"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq between int
< error
| 42601 syntax error at or near "between"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq between BY NONE
< error
| 42601 syntax error at or near "between"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq between NAME omni_probe_name
< error
| 42601 syntax error at or near "between"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq between CYCLE
< error
| 42601 syntax error at or near "between"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bigint
< error
| 42601 syntax error at or near "bigint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bigint 1
< error
| 42601 syntax error at or near "bigint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bigint int
< error
| 42601 syntax error at or near "bigint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bigint BY NONE
< error
| 42601 syntax error at or near "bigint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bigint NAME omni_probe_name
< error
| 42601 syntax error at or near "bigint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bigint CYCLE
< error
| 42601 syntax error at or near "bigint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq binary
< error
| 42601 syntax error at or near "binary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq binary 1
< error
| 42601 syntax error at or near "binary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq binary int
< error
| 42601 syntax error at or near "binary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq binary BY NONE
< error
| 42601 syntax error at or near "binary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq binary NAME omni_probe_name
< error
| 42601 syntax error at or near "binary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq binary CYCLE
< error
| 42601 syntax error at or near "binary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bit
< error
| 42601 syntax error at or near "bit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bit 1
< error
| 42601 syntax error at or near "bit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bit int
< error
| 42601 syntax error at or near "bit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bit BY NONE
< error
| 42601 syntax error at or near "bit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bit NAME omni_probe_name
< error
| 42601 syntax error at or near "bit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq bit CYCLE
< error
| 42601 syntax error at or near "bit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq boolean
< error
| 42601 syntax error at or near "boolean"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq boolean 1
< error
| 42601 syntax error at or near "boolean"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq boolean int
< error
| 42601 syntax error at or near "boolean"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq boolean BY NONE
< error
| 42601 syntax error at or near "boolean"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq boolean NAME omni_probe_name
< error
| 42601 syntax error at or near "boolean"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq boolean CYCLE
< error
| 42601 syntax error at or near "boolean"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq both
< error
| 42601 syntax error at or near "both"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq both 1
< error
| 42601 syntax error at or near "both"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq both int
< error
| 42601 syntax error at or near "both"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq both BY NONE
< error
| 42601 syntax error at or near "both"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq both NAME omni_probe_name
< error
| 42601 syntax error at or near "both"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq both CYCLE
< error
| 42601 syntax error at or near "both"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq breadth
< error
| 42601 syntax error at or near "breadth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq breadth 1
< error
| 42601 syntax error at or near "breadth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq breadth int
< error
| 42601 syntax error at or near "breadth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq breadth BY NONE
< error
| 42601 syntax error at or near "breadth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq breadth NAME omni_probe_name
< error
| 42601 syntax error at or near "breadth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq breadth CYCLE
< error
| 42601 syntax error at or near "breadth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq by
< error
| 42601 syntax error at or near "by"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq by 1
< error
| 42601 syntax error at or near "by"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq by int
< error
| 42601 syntax error at or near "by"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq by BY NONE
< error
| 42601 syntax error at or near "by"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq by NAME omni_probe_name
< error
| 42601 syntax error at or near "by"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq by CYCLE
< error
| 42601 syntax error at or near "by"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cache
< error
| 42601 syntax error at end of input
| position: 37
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cache 1
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq call
< error
| 42601 syntax error at or near "call"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq call 1
< error
| 42601 syntax error at or near "call"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq call int
< error
| 42601 syntax error at or near "call"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq call BY NONE
< error
| 42601 syntax error at or near "call"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq call NAME omni_probe_name
< error
| 42601 syntax error at or near "call"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq call CYCLE
< error
| 42601 syntax error at or near "call"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq called
< error
| 42601 syntax error at or near "called"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq called 1
< error
| 42601 syntax error at or near "called"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq called int
< error
| 42601 syntax error at or near "called"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq called BY NONE
< error
| 42601 syntax error at or near "called"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq called NAME omni_probe_name
< error
| 42601 syntax error at or near "called"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq called CYCLE
< error
| 42601 syntax error at or near "called"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascade
< error
| 42601 syntax error at or near "cascade"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascade 1
< error
| 42601 syntax error at or near "cascade"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascade int
< error
| 42601 syntax error at or near "cascade"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascade BY NONE
< error
| 42601 syntax error at or near "cascade"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascade NAME omni_probe_name
< error
| 42601 syntax error at or near "cascade"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascade CYCLE
< error
| 42601 syntax error at or near "cascade"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascaded
< error
| 42601 syntax error at or near "cascaded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascaded 1
< error
| 42601 syntax error at or near "cascaded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascaded int
< error
| 42601 syntax error at or near "cascaded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascaded BY NONE
< error
| 42601 syntax error at or near "cascaded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascaded NAME omni_probe_name
< error
| 42601 syntax error at or near "cascaded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cascaded CYCLE
< error
| 42601 syntax error at or near "cascaded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq case
< error
| 42601 syntax error at or near "case"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq case 1
< error
| 42601 syntax error at or near "case"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq case int
< error
| 42601 syntax error at or near "case"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq case BY NONE
< error
| 42601 syntax error at or near "case"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq case NAME omni_probe_name
< error
| 42601 syntax error at or near "case"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq case CYCLE
< error
| 42601 syntax error at or near "case"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cast
< error
| 42601 syntax error at or near "cast"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cast 1
< error
| 42601 syntax error at or near "cast"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cast int
< error
| 42601 syntax error at or near "cast"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cast BY NONE
< error
| 42601 syntax error at or near "cast"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cast NAME omni_probe_name
< error
| 42601 syntax error at or near "cast"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cast CYCLE
< error
| 42601 syntax error at or near "cast"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq catalog
< error
| 42601 syntax error at or near "catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq catalog 1
< error
| 42601 syntax error at or near "catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq catalog int
< error
| 42601 syntax error at or near "catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq catalog BY NONE
< error
| 42601 syntax error at or near "catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq catalog NAME omni_probe_name
< error
| 42601 syntax error at or near "catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq catalog CYCLE
< error
| 42601 syntax error at or near "catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq chain
< error
| 42601 syntax error at or near "chain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq chain 1
< error
| 42601 syntax error at or near "chain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq chain int
< error
| 42601 syntax error at or near "chain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq chain BY NONE
< error
| 42601 syntax error at or near "chain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq chain NAME omni_probe_name
< error
| 42601 syntax error at or near "chain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq chain CYCLE
< error
| 42601 syntax error at or near "chain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq char
< error
| 42601 syntax error at or near "char"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq char 1
< error
| 42601 syntax error at or near "char"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq char int
< error
| 42601 syntax error at or near "char"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq char BY NONE
< error
| 42601 syntax error at or near "char"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq char NAME omni_probe_name
< error
| 42601 syntax error at or near "char"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq char CYCLE
< error
| 42601 syntax error at or near "char"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq character
< error
| 42601 syntax error at or near "character"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq character 1
< error
| 42601 syntax error at or near "character"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq character int
< error
| 42601 syntax error at or near "character"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq character BY NONE
< error
| 42601 syntax error at or near "character"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq character NAME omni_probe_name
< error
| 42601 syntax error at or near "character"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq character CYCLE
< error
| 42601 syntax error at or near "character"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq characteristics
< error
| 42601 syntax error at or near "characteristics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq characteristics 1
< error
| 42601 syntax error at or near "characteristics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq characteristics int
< error
| 42601 syntax error at or near "characteristics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq characteristics BY NONE
< error
| 42601 syntax error at or near "characteristics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq characteristics NAME omni_probe_name
< error
| 42601 syntax error at or near "characteristics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq characteristics CYCLE
< error
| 42601 syntax error at or near "characteristics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq check
< error
| 42601 syntax error at or near "check"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq check 1
< error
| 42601 syntax error at or near "check"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq check int
< error
| 42601 syntax error at or near "check"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq check BY NONE
< error
| 42601 syntax error at or near "check"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq check NAME omni_probe_name
< error
| 42601 syntax error at or near "check"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq check CYCLE
< error
| 42601 syntax error at or near "check"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq checkpoint
< error
| 42601 syntax error at or near "checkpoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq checkpoint 1
< error
| 42601 syntax error at or near "checkpoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq checkpoint int
< error
| 42601 syntax error at or near "checkpoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq checkpoint BY NONE
< error
| 42601 syntax error at or near "checkpoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq checkpoint NAME omni_probe_name
< error
| 42601 syntax error at or near "checkpoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq checkpoint CYCLE
< error
| 42601 syntax error at or near "checkpoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq class
< error
| 42601 syntax error at or near "class"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq class 1
< error
| 42601 syntax error at or near "class"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq class int
< error
| 42601 syntax error at or near "class"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq class BY NONE
< error
| 42601 syntax error at or near "class"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq class NAME omni_probe_name
< error
| 42601 syntax error at or near "class"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq class CYCLE
< error
| 42601 syntax error at or near "class"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq close
< error
| 42601 syntax error at or near "close"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq close 1
< error
| 42601 syntax error at or near "close"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq close int
< error
| 42601 syntax error at or near "close"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq close BY NONE
< error
| 42601 syntax error at or near "close"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq close NAME omni_probe_name
< error
| 42601 syntax error at or near "close"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq close CYCLE
< error
| 42601 syntax error at or near "close"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cluster
< error
| 42601 syntax error at or near "cluster"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cluster 1
< error
| 42601 syntax error at or near "cluster"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cluster int
< error
| 42601 syntax error at or near "cluster"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cluster BY NONE
< error
| 42601 syntax error at or near "cluster"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cluster NAME omni_probe_name
< error
| 42601 syntax error at or near "cluster"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cluster CYCLE
< error
| 42601 syntax error at or near "cluster"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq coalesce
< error
| 42601 syntax error at or near "coalesce"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq coalesce 1
< error
| 42601 syntax error at or near "coalesce"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq coalesce int
< error
| 42601 syntax error at or near "coalesce"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq coalesce BY NONE
< error
| 42601 syntax error at or near "coalesce"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq coalesce NAME omni_probe_name
< error
| 42601 syntax error at or near "coalesce"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq coalesce CYCLE
< error
| 42601 syntax error at or near "coalesce"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collate
< error
| 42601 syntax error at or near "collate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collate 1
< error
| 42601 syntax error at or near "collate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collate int
< error
| 42601 syntax error at or near "collate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collate BY NONE
< error
| 42601 syntax error at or near "collate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collate NAME omni_probe_name
< error
| 42601 syntax error at or near "collate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collate CYCLE
< error
| 42601 syntax error at or near "collate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collation
< error
| 42601 syntax error at or near "collation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collation 1
< error
| 42601 syntax error at or near "collation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collation int
< error
| 42601 syntax error at or near "collation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collation BY NONE
< error
| 42601 syntax error at or near "collation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collation NAME omni_probe_name
< error
| 42601 syntax error at or near "collation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq collation CYCLE
< error
| 42601 syntax error at or near "collation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq column
< error
| 42601 syntax error at or near "column"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq column 1
< error
| 42601 syntax error at or near "column"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq column int
< error
| 42601 syntax error at or near "column"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq column BY NONE
< error
| 42601 syntax error at or near "column"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq column NAME omni_probe_name
< error
| 42601 syntax error at or near "column"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq column CYCLE
< error
| 42601 syntax error at or near "column"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq columns
< error
| 42601 syntax error at or near "columns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq columns 1
< error
| 42601 syntax error at or near "columns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq columns int
< error
| 42601 syntax error at or near "columns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq columns BY NONE
< error
| 42601 syntax error at or near "columns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq columns NAME omni_probe_name
< error
| 42601 syntax error at or near "columns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq columns CYCLE
< error
| 42601 syntax error at or near "columns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comment
< error
| 42601 syntax error at or near "comment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comment 1
< error
| 42601 syntax error at or near "comment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comment int
< error
| 42601 syntax error at or near "comment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comment BY NONE
< error
| 42601 syntax error at or near "comment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comment NAME omni_probe_name
< error
| 42601 syntax error at or near "comment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comment CYCLE
< error
| 42601 syntax error at or near "comment"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comments
< error
| 42601 syntax error at or near "comments"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comments 1
< error
| 42601 syntax error at or near "comments"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comments int
< error
| 42601 syntax error at or near "comments"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comments BY NONE
< error
| 42601 syntax error at or near "comments"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comments NAME omni_probe_name
< error
| 42601 syntax error at or near "comments"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq comments CYCLE
< error
| 42601 syntax error at or near "comments"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq commit
< error
| 42601 syntax error at or near "commit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq commit 1
< error
| 42601 syntax error at or near "commit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq commit int
< error
| 42601 syntax error at or near "commit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq commit BY NONE
< error
| 42601 syntax error at or near "commit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq commit NAME omni_probe_name
< error
| 42601 syntax error at or near "commit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq commit CYCLE
< error
| 42601 syntax error at or near "commit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq committed
< error
| 42601 syntax error at or near "committed"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq committed 1
< error
| 42601 syntax error at or near "committed"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq committed int
< error
| 42601 syntax error at or near "committed"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq committed BY NONE
< error
| 42601 syntax error at or near "committed"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq committed NAME omni_probe_name
< error
| 42601 syntax error at or near "committed"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq committed CYCLE
< error
| 42601 syntax error at or near "committed"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq compression
< error
| 42601 syntax error at or near "compression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq compression 1
< error
| 42601 syntax error at or near "compression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq compression int
< error
| 42601 syntax error at or near "compression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq compression BY NONE
< error
| 42601 syntax error at or near "compression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq compression NAME omni_probe_name
< error
| 42601 syntax error at or near "compression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq compression CYCLE
< error
| 42601 syntax error at or near "compression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq concurrently
< error
| 42601 syntax error at or near "concurrently"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq concurrently 1
< error
| 42601 syntax error at or near "concurrently"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq concurrently int
< error
| 42601 syntax error at or near "concurrently"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq concurrently BY NONE
< error
| 42601 syntax error at or near "concurrently"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq concurrently NAME omni_probe_name
< error
| 42601 syntax error at or near "concurrently"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq concurrently CYCLE
< error
| 42601 syntax error at or near "concurrently"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conditional
< error
| 42601 syntax error at or near "conditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conditional 1
< error
| 42601 syntax error at or near "conditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conditional int
< error
| 42601 syntax error at or near "conditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conditional BY NONE
< error
| 42601 syntax error at or near "conditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conditional NAME omni_probe_name
< error
| 42601 syntax error at or near "conditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conditional CYCLE
< error
| 42601 syntax error at or near "conditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq configuration
< error
| 42601 syntax error at or near "configuration"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq configuration 1
< error
| 42601 syntax error at or near "configuration"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq configuration int
< error
| 42601 syntax error at or near "configuration"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq configuration BY NONE
< error
| 42601 syntax error at or near "configuration"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq configuration NAME omni_probe_name
< error
| 42601 syntax error at or near "configuration"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq configuration CYCLE
< error
| 42601 syntax error at or near "configuration"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conflict
< error
| 42601 syntax error at or near "conflict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conflict 1
< error
| 42601 syntax error at or near "conflict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conflict int
< error
| 42601 syntax error at or near "conflict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conflict BY NONE
< error
| 42601 syntax error at or near "conflict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conflict NAME omni_probe_name
< error
| 42601 syntax error at or near "conflict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conflict CYCLE
< error
| 42601 syntax error at or near "conflict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq connection
< error
| 42601 syntax error at or near "connection"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq connection 1
< error
| 42601 syntax error at or near "connection"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq connection int
< error
| 42601 syntax error at or near "connection"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq connection BY NONE
< error
| 42601 syntax error at or near "connection"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq connection NAME omni_probe_name
< error
| 42601 syntax error at or near "connection"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq connection CYCLE
< error
| 42601 syntax error at or near "connection"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraint
< error
| 42601 syntax error at or near "constraint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraint 1
< error
| 42601 syntax error at or near "constraint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraint int
< error
| 42601 syntax error at or near "constraint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraint BY NONE
< error
| 42601 syntax error at or near "constraint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraint NAME omni_probe_name
< error
| 42601 syntax error at or near "constraint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraint CYCLE
< error
| 42601 syntax error at or near "constraint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraints
< error
| 42601 syntax error at or near "constraints"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraints 1
< error
| 42601 syntax error at or near "constraints"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraints int
< error
| 42601 syntax error at or near "constraints"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraints BY NONE
< error
| 42601 syntax error at or near "constraints"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraints NAME omni_probe_name
< error
| 42601 syntax error at or near "constraints"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq constraints CYCLE
< error
| 42601 syntax error at or near "constraints"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq content
< error
| 42601 syntax error at or near "content"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq content 1
< error
| 42601 syntax error at or near "content"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq content int
< error
| 42601 syntax error at or near "content"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq content BY NONE
< error
| 42601 syntax error at or near "content"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq content NAME omni_probe_name
< error
| 42601 syntax error at or near "content"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq content CYCLE
< error
| 42601 syntax error at or near "content"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq continue
< error
| 42601 syntax error at or near "continue"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq continue 1
< error
| 42601 syntax error at or near "continue"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq continue int
< error
| 42601 syntax error at or near "continue"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq continue BY NONE
< error
| 42601 syntax error at or near "continue"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq continue NAME omni_probe_name
< error
| 42601 syntax error at or near "continue"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq continue CYCLE
< error
| 42601 syntax error at or near "continue"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conversion
< error
| 42601 syntax error at or near "conversion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conversion 1
< error
| 42601 syntax error at or near "conversion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conversion int
< error
| 42601 syntax error at or near "conversion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conversion BY NONE
< error
| 42601 syntax error at or near "conversion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conversion NAME omni_probe_name
< error
| 42601 syntax error at or near "conversion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq conversion CYCLE
< error
| 42601 syntax error at or near "conversion"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq copy
< error
| 42601 syntax error at or near "copy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq copy 1
< error
| 42601 syntax error at or near "copy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq copy int
< error
| 42601 syntax error at or near "copy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq copy BY NONE
< error
| 42601 syntax error at or near "copy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq copy NAME omni_probe_name
< error
| 42601 syntax error at or near "copy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq copy CYCLE
< error
| 42601 syntax error at or near "copy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cost
< error
| 42601 syntax error at or near "cost"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cost 1
< error
| 42601 syntax error at or near "cost"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cost int
< error
| 42601 syntax error at or near "cost"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cost BY NONE
< error
| 42601 syntax error at or near "cost"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cost NAME omni_probe_name
< error
| 42601 syntax error at or near "cost"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cost CYCLE
< error
| 42601 syntax error at or near "cost"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq create
< error
| 42601 syntax error at or near "create"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq create 1
< error
| 42601 syntax error at or near "create"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq create int
< error
| 42601 syntax error at or near "create"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq create BY NONE
< error
| 42601 syntax error at or near "create"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq create NAME omni_probe_name
< error
| 42601 syntax error at or near "create"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq create CYCLE
< error
| 42601 syntax error at or near "create"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cross
< error
| 42601 syntax error at or near "cross"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cross 1
< error
| 42601 syntax error at or near "cross"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cross int
< error
| 42601 syntax error at or near "cross"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cross BY NONE
< error
| 42601 syntax error at or near "cross"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cross NAME omni_probe_name
< error
| 42601 syntax error at or near "cross"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cross CYCLE
< error
| 42601 syntax error at or near "cross"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq csv
< error
| 42601 syntax error at or near "csv"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq csv 1
< error
| 42601 syntax error at or near "csv"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq csv int
< error
| 42601 syntax error at or near "csv"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq csv BY NONE
< error
| 42601 syntax error at or near "csv"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq csv NAME omni_probe_name
< error
| 42601 syntax error at or near "csv"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq csv CYCLE
< error
| 42601 syntax error at or near "csv"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cube
< error
| 42601 syntax error at or near "cube"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cube 1
< error
| 42601 syntax error at or near "cube"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cube int
< error
| 42601 syntax error at or near "cube"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cube BY NONE
< error
| 42601 syntax error at or near "cube"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cube NAME omni_probe_name
< error
| 42601 syntax error at or near "cube"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cube CYCLE
< error
| 42601 syntax error at or near "cube"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current
< error
| 42601 syntax error at or near "current"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current 1
< error
| 42601 syntax error at or near "current"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current int
< error
| 42601 syntax error at or near "current"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current BY NONE
< error
| 42601 syntax error at or near "current"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current NAME omni_probe_name
< error
| 42601 syntax error at or near "current"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current CYCLE
< error
| 42601 syntax error at or near "current"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_catalog
< error
| 42601 syntax error at or near "current_catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_catalog 1
< error
| 42601 syntax error at or near "current_catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_catalog int
< error
| 42601 syntax error at or near "current_catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_catalog BY NONE
< error
| 42601 syntax error at or near "current_catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_catalog NAME omni_probe_name
< error
| 42601 syntax error at or near "current_catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_catalog CYCLE
< error
| 42601 syntax error at or near "current_catalog"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_date
< error
| 42601 syntax error at or near "current_date"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_date 1
< error
| 42601 syntax error at or near "current_date"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_date int
< error
| 42601 syntax error at or near "current_date"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_date BY NONE
< error
| 42601 syntax error at or near "current_date"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_date NAME omni_probe_name
< error
| 42601 syntax error at or near "current_date"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_date CYCLE
< error
| 42601 syntax error at or near "current_date"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_role
< error
| 42601 syntax error at or near "current_role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_role 1
< error
| 42601 syntax error at or near "current_role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_role int
< error
| 42601 syntax error at or near "current_role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_role BY NONE
< error
| 42601 syntax error at or near "current_role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_role NAME omni_probe_name
< error
| 42601 syntax error at or near "current_role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_role CYCLE
< error
| 42601 syntax error at or near "current_role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_schema
< error
| 42601 syntax error at or near "current_schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_schema 1
< error
| 42601 syntax error at or near "current_schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_schema int
< error
| 42601 syntax error at or near "current_schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_schema BY NONE
< error
| 42601 syntax error at or near "current_schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_schema NAME omni_probe_name
< error
| 42601 syntax error at or near "current_schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_schema CYCLE
< error
| 42601 syntax error at or near "current_schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_time
< error
| 42601 syntax error at or near "current_time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_time 1
< error
| 42601 syntax error at or near "current_time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_time int
< error
| 42601 syntax error at or near "current_time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_time BY NONE
< error
| 42601 syntax error at or near "current_time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_time NAME omni_probe_name
< error
| 42601 syntax error at or near "current_time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_time CYCLE
< error
| 42601 syntax error at or near "current_time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_timestamp
< error
| 42601 syntax error at or near "current_timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_timestamp 1
< error
| 42601 syntax error at or near "current_timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_timestamp int
< error
| 42601 syntax error at or near "current_timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_timestamp BY NONE
< error
| 42601 syntax error at or near "current_timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_timestamp NAME omni_probe_name
< error
| 42601 syntax error at or near "current_timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_timestamp CYCLE
< error
| 42601 syntax error at or near "current_timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_user
< error
| 42601 syntax error at or near "current_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_user 1
< error
| 42601 syntax error at or near "current_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_user int
< error
| 42601 syntax error at or near "current_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_user BY NONE
< error
| 42601 syntax error at or near "current_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_user NAME omni_probe_name
< error
| 42601 syntax error at or near "current_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq current_user CYCLE
< error
| 42601 syntax error at or near "current_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cursor
< error
| 42601 syntax error at or near "cursor"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cursor 1
< error
| 42601 syntax error at or near "cursor"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cursor int
< error
| 42601 syntax error at or near "cursor"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cursor BY NONE
< error
| 42601 syntax error at or near "cursor"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cursor NAME omni_probe_name
< error
| 42601 syntax error at or near "cursor"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cursor CYCLE
< error
| 42601 syntax error at or near "cursor"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq cycle
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq data
< error
| 42601 syntax error at or near "data"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq data 1
< error
| 42601 syntax error at or near "data"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq data int
< error
| 42601 syntax error at or near "data"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq data BY NONE
< error
| 42601 syntax error at or near "data"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq data NAME omni_probe_name
< error
| 42601 syntax error at or near "data"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq data CYCLE
< error
| 42601 syntax error at or near "data"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq database
< error
| 42601 syntax error at or near "database"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq database 1
< error
| 42601 syntax error at or near "database"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq database int
< error
| 42601 syntax error at or near "database"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq database BY NONE
< error
| 42601 syntax error at or near "database"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq database NAME omni_probe_name
< error
| 42601 syntax error at or near "database"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq database CYCLE
< error
| 42601 syntax error at or near "database"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq day
< error
| 42601 syntax error at or near "day"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq day 1
< error
| 42601 syntax error at or near "day"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq day int
< error
| 42601 syntax error at or near "day"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq day BY NONE
< error
| 42601 syntax error at or near "day"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq day NAME omni_probe_name
< error
| 42601 syntax error at or near "day"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq day CYCLE
< error
| 42601 syntax error at or near "day"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deallocate
< error
| 42601 syntax error at or near "deallocate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deallocate 1
< error
| 42601 syntax error at or near "deallocate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deallocate int
< error
| 42601 syntax error at or near "deallocate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deallocate BY NONE
< error
| 42601 syntax error at or near "deallocate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deallocate NAME omni_probe_name
< error
| 42601 syntax error at or near "deallocate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deallocate CYCLE
< error
| 42601 syntax error at or near "deallocate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dec
< error
| 42601 syntax error at or near "dec"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dec 1
< error
| 42601 syntax error at or near "dec"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dec int
< error
| 42601 syntax error at or near "dec"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dec BY NONE
< error
| 42601 syntax error at or near "dec"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dec NAME omni_probe_name
< error
| 42601 syntax error at or near "dec"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dec CYCLE
< error
| 42601 syntax error at or near "dec"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq decimal
< error
| 42601 syntax error at or near "decimal"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq decimal 1
< error
| 42601 syntax error at or near "decimal"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq decimal int
< error
| 42601 syntax error at or near "decimal"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq decimal BY NONE
< error
| 42601 syntax error at or near "decimal"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq decimal NAME omni_probe_name
< error
| 42601 syntax error at or near "decimal"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq decimal CYCLE
< error
| 42601 syntax error at or near "decimal"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq declare
< error
| 42601 syntax error at or near "declare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq declare 1
< error
| 42601 syntax error at or near "declare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq declare int
< error
| 42601 syntax error at or near "declare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq declare BY NONE
< error
| 42601 syntax error at or near "declare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq declare NAME omni_probe_name
< error
| 42601 syntax error at or near "declare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq declare CYCLE
< error
| 42601 syntax error at or near "declare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq default
< error
| 42601 syntax error at or near "default"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq default 1
< error
| 42601 syntax error at or near "default"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq default int
< error
| 42601 syntax error at or near "default"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq default BY NONE
< error
| 42601 syntax error at or near "default"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq default NAME omni_probe_name
< error
| 42601 syntax error at or near "default"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq default CYCLE
< error
| 42601 syntax error at or near "default"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq defaults
< error
| 42601 syntax error at or near "defaults"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq defaults 1
< error
| 42601 syntax error at or near "defaults"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq defaults int
< error
| 42601 syntax error at or near "defaults"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq defaults BY NONE
< error
| 42601 syntax error at or near "defaults"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq defaults NAME omni_probe_name
< error
| 42601 syntax error at or near "defaults"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq defaults CYCLE
< error
| 42601 syntax error at or near "defaults"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferrable
< error
| 42601 syntax error at or near "deferrable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferrable 1
< error
| 42601 syntax error at or near "deferrable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferrable int
< error
| 42601 syntax error at or near "deferrable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferrable BY NONE
< error
| 42601 syntax error at or near "deferrable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferrable NAME omni_probe_name
< error
| 42601 syntax error at or near "deferrable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferrable CYCLE
< error
| 42601 syntax error at or near "deferrable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferred
< error
| 42601 syntax error at or near "deferred"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferred 1
< error
| 42601 syntax error at or near "deferred"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferred int
< error
| 42601 syntax error at or near "deferred"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferred BY NONE
< error
| 42601 syntax error at or near "deferred"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferred NAME omni_probe_name
< error
| 42601 syntax error at or near "deferred"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq deferred CYCLE
< error
| 42601 syntax error at or near "deferred"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq definer
< error
| 42601 syntax error at or near "definer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq definer 1
< error
| 42601 syntax error at or near "definer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq definer int
< error
| 42601 syntax error at or near "definer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq definer BY NONE
< error
| 42601 syntax error at or near "definer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq definer NAME omni_probe_name
< error
| 42601 syntax error at or near "definer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq definer CYCLE
< error
| 42601 syntax error at or near "definer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delete
< error
| 42601 syntax error at or near "delete"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delete 1
< error
| 42601 syntax error at or near "delete"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delete int
< error
| 42601 syntax error at or near "delete"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delete BY NONE
< error
| 42601 syntax error at or near "delete"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delete NAME omni_probe_name
< error
| 42601 syntax error at or near "delete"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delete CYCLE
< error
| 42601 syntax error at or near "delete"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiter
< error
| 42601 syntax error at or near "delimiter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiter 1
< error
| 42601 syntax error at or near "delimiter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiter int
< error
| 42601 syntax error at or near "delimiter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiter BY NONE
< error
| 42601 syntax error at or near "delimiter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiter NAME omni_probe_name
< error
| 42601 syntax error at or near "delimiter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiter CYCLE
< error
| 42601 syntax error at or near "delimiter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiters
< error
| 42601 syntax error at or near "delimiters"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiters 1
< error
| 42601 syntax error at or near "delimiters"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiters int
< error
| 42601 syntax error at or near "delimiters"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiters BY NONE
< error
| 42601 syntax error at or near "delimiters"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiters NAME omni_probe_name
< error
| 42601 syntax error at or near "delimiters"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq delimiters CYCLE
< error
| 42601 syntax error at or near "delimiters"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depends
< error
| 42601 syntax error at or near "depends"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depends 1
< error
| 42601 syntax error at or near "depends"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depends int
< error
| 42601 syntax error at or near "depends"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depends BY NONE
< error
| 42601 syntax error at or near "depends"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depends NAME omni_probe_name
< error
| 42601 syntax error at or near "depends"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depends CYCLE
< error
| 42601 syntax error at or near "depends"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depth
< error
| 42601 syntax error at or near "depth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depth 1
< error
| 42601 syntax error at or near "depth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depth int
< error
| 42601 syntax error at or near "depth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depth BY NONE
< error
| 42601 syntax error at or near "depth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depth NAME omni_probe_name
< error
| 42601 syntax error at or near "depth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq depth CYCLE
< error
| 42601 syntax error at or near "depth"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq desc
< error
| 42601 syntax error at or near "desc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq desc 1
< error
| 42601 syntax error at or near "desc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq desc int
< error
| 42601 syntax error at or near "desc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq desc BY NONE
< error
| 42601 syntax error at or near "desc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq desc NAME omni_probe_name
< error
| 42601 syntax error at or near "desc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq desc CYCLE
< error
| 42601 syntax error at or near "desc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq detach
< error
| 42601 syntax error at or near "detach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq detach 1
< error
| 42601 syntax error at or near "detach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq detach int
< error
| 42601 syntax error at or near "detach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq detach BY NONE
< error
| 42601 syntax error at or near "detach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq detach NAME omni_probe_name
< error
| 42601 syntax error at or near "detach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq detach CYCLE
< error
| 42601 syntax error at or near "detach"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dictionary
< error
| 42601 syntax error at or near "dictionary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dictionary 1
< error
| 42601 syntax error at or near "dictionary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dictionary int
< error
| 42601 syntax error at or near "dictionary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dictionary BY NONE
< error
| 42601 syntax error at or near "dictionary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dictionary NAME omni_probe_name
< error
| 42601 syntax error at or near "dictionary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq dictionary CYCLE
< error
| 42601 syntax error at or near "dictionary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq disable
< error
| 42601 syntax error at or near "disable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq disable 1
< error
| 42601 syntax error at or near "disable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq disable int
< error
| 42601 syntax error at or near "disable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq disable BY NONE
< error
| 42601 syntax error at or near "disable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq disable NAME omni_probe_name
< error
| 42601 syntax error at or near "disable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq disable CYCLE
< error
| 42601 syntax error at or near "disable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq discard
< error
| 42601 syntax error at or near "discard"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq discard 1
< error
| 42601 syntax error at or near "discard"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq discard int
< error
| 42601 syntax error at or near "discard"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq discard BY NONE
< error
| 42601 syntax error at or near "discard"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq discard NAME omni_probe_name
< error
| 42601 syntax error at or near "discard"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq discard CYCLE
< error
| 42601 syntax error at or near "discard"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq distinct
< error
| 42601 syntax error at or near "distinct"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq distinct 1
< error
| 42601 syntax error at or near "distinct"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq distinct int
< error
| 42601 syntax error at or near "distinct"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq distinct BY NONE
< error
| 42601 syntax error at or near "distinct"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq distinct NAME omni_probe_name
< error
| 42601 syntax error at or near "distinct"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq distinct CYCLE
< error
| 42601 syntax error at or near "distinct"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq do
< error
| 42601 syntax error at or near "do"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq do 1
< error
| 42601 syntax error at or near "do"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq do int
< error
| 42601 syntax error at or near "do"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq do BY NONE
< error
| 42601 syntax error at or near "do"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq do NAME omni_probe_name
< error
| 42601 syntax error at or near "do"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq do CYCLE
< error
| 42601 syntax error at or near "do"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq document
< error
| 42601 syntax error at or near "document"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq document 1
< error
| 42601 syntax error at or near "document"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq document int
< error
| 42601 syntax error at or near "document"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq document BY NONE
< error
| 42601 syntax error at or near "document"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq document NAME omni_probe_name
< error
| 42601 syntax error at or near "document"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq document CYCLE
< error
| 42601 syntax error at or near "document"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq domain
< error
| 42601 syntax error at or near "domain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq domain 1
< error
| 42601 syntax error at or near "domain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq domain int
< error
| 42601 syntax error at or near "domain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq domain BY NONE
< error
| 42601 syntax error at or near "domain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq domain NAME omni_probe_name
< error
| 42601 syntax error at or near "domain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq domain CYCLE
< error
| 42601 syntax error at or near "domain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq double
< error
| 42601 syntax error at or near "double"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq double 1
< error
| 42601 syntax error at or near "double"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq double int
< error
| 42601 syntax error at or near "double"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq double BY NONE
< error
| 42601 syntax error at or near "double"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq double NAME omni_probe_name
< error
| 42601 syntax error at or near "double"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq double CYCLE
< error
| 42601 syntax error at or near "double"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq drop
< error
| 42601 syntax error at or near "drop"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq drop 1
< error
| 42601 syntax error at or near "drop"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq drop int
< error
| 42601 syntax error at or near "drop"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq drop BY NONE
< error
| 42601 syntax error at or near "drop"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq drop NAME omni_probe_name
< error
| 42601 syntax error at or near "drop"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq drop CYCLE
< error
| 42601 syntax error at or near "drop"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq each
< error
| 42601 syntax error at or near "each"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq each 1
< error
| 42601 syntax error at or near "each"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq each int
< error
| 42601 syntax error at or near "each"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq each BY NONE
< error
| 42601 syntax error at or near "each"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq each NAME omni_probe_name
< error
| 42601 syntax error at or near "each"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq each CYCLE
< error
| 42601 syntax error at or near "each"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq else
< error
| 42601 syntax error at or near "else"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq else 1
< error
| 42601 syntax error at or near "else"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq else int
< error
| 42601 syntax error at or near "else"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq else BY NONE
< error
| 42601 syntax error at or near "else"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq else NAME omni_probe_name
< error
| 42601 syntax error at or near "else"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq else CYCLE
< error
| 42601 syntax error at or near "else"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq empty
< error
| 42601 syntax error at or near "empty"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq empty 1
< error
| 42601 syntax error at or near "empty"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq empty int
< error
| 42601 syntax error at or near "empty"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq empty BY NONE
< error
| 42601 syntax error at or near "empty"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq empty NAME omni_probe_name
< error
| 42601 syntax error at or near "empty"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq empty CYCLE
< error
| 42601 syntax error at or near "empty"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enable
< error
| 42601 syntax error at or near "enable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enable 1
< error
| 42601 syntax error at or near "enable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enable int
< error
| 42601 syntax error at or near "enable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enable BY NONE
< error
| 42601 syntax error at or near "enable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enable NAME omni_probe_name
< error
| 42601 syntax error at or near "enable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enable CYCLE
< error
| 42601 syntax error at or near "enable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encoding
< error
| 42601 syntax error at or near "encoding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encoding 1
< error
| 42601 syntax error at or near "encoding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encoding int
< error
| 42601 syntax error at or near "encoding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encoding BY NONE
< error
| 42601 syntax error at or near "encoding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encoding NAME omni_probe_name
< error
| 42601 syntax error at or near "encoding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encoding CYCLE
< error
| 42601 syntax error at or near "encoding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encrypted
< error
| 42601 syntax error at or near "encrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encrypted 1
< error
| 42601 syntax error at or near "encrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encrypted int
< error
| 42601 syntax error at or near "encrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encrypted BY NONE
< error
| 42601 syntax error at or near "encrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encrypted NAME omni_probe_name
< error
| 42601 syntax error at or near "encrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq encrypted CYCLE
< error
| 42601 syntax error at or near "encrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq end
< error
| 42601 syntax error at or near "end"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq end 1
< error
| 42601 syntax error at or near "end"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq end int
< error
| 42601 syntax error at or near "end"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq end BY NONE
< error
| 42601 syntax error at or near "end"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq end NAME omni_probe_name
< error
| 42601 syntax error at or near "end"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq end CYCLE
< error
| 42601 syntax error at or near "end"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enum
< error
| 42601 syntax error at or near "enum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enum 1
< error
| 42601 syntax error at or near "enum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enum int
< error
| 42601 syntax error at or near "enum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enum BY NONE
< error
| 42601 syntax error at or near "enum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enum NAME omni_probe_name
< error
| 42601 syntax error at or near "enum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq enum CYCLE
< error
| 42601 syntax error at or near "enum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq error
< error
| 42601 syntax error at or near "error"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq error 1
< error
| 42601 syntax error at or near "error"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq error int
< error
| 42601 syntax error at or near "error"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq error BY NONE
< error
| 42601 syntax error at or near "error"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq error NAME omni_probe_name
< error
| 42601 syntax error at or near "error"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq error CYCLE
< error
| 42601 syntax error at or near "error"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq escape
< error
| 42601 syntax error at or near "escape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq escape 1
< error
| 42601 syntax error at or near "escape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq escape int
< error
| 42601 syntax error at or near "escape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq escape BY NONE
< error
| 42601 syntax error at or near "escape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq escape NAME omni_probe_name
< error
| 42601 syntax error at or near "escape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq escape CYCLE
< error
| 42601 syntax error at or near "escape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq event
< error
| 42601 syntax error at or near "event"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq event 1
< error
| 42601 syntax error at or near "event"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq event int
< error
| 42601 syntax error at or near "event"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq event BY NONE
< error
| 42601 syntax error at or near "event"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq event NAME omni_probe_name
< error
| 42601 syntax error at or near "event"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq event CYCLE
< error
| 42601 syntax error at or near "event"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq except
< error
| 42601 syntax error at or near "except"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq except 1
< error
| 42601 syntax error at or near "except"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq except int
< error
| 42601 syntax error at or near "except"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq except BY NONE
< error
| 42601 syntax error at or near "except"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq except NAME omni_probe_name
< error
| 42601 syntax error at or near "except"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq except CYCLE
< error
| 42601 syntax error at or near "except"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclude
< error
| 42601 syntax error at or near "exclude"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclude 1
< error
| 42601 syntax error at or near "exclude"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclude int
< error
| 42601 syntax error at or near "exclude"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclude BY NONE
< error
| 42601 syntax error at or near "exclude"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclude NAME omni_probe_name
< error
| 42601 syntax error at or near "exclude"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclude CYCLE
< error
| 42601 syntax error at or near "exclude"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq excluding
< error
| 42601 syntax error at or near "excluding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq excluding 1
< error
| 42601 syntax error at or near "excluding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq excluding int
< error
| 42601 syntax error at or near "excluding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq excluding BY NONE
< error
| 42601 syntax error at or near "excluding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq excluding NAME omni_probe_name
< error
| 42601 syntax error at or near "excluding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq excluding CYCLE
< error
| 42601 syntax error at or near "excluding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclusive
< error
| 42601 syntax error at or near "exclusive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclusive 1
< error
| 42601 syntax error at or near "exclusive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclusive int
< error
| 42601 syntax error at or near "exclusive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclusive BY NONE
< error
| 42601 syntax error at or near "exclusive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclusive NAME omni_probe_name
< error
| 42601 syntax error at or near "exclusive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exclusive CYCLE
< error
| 42601 syntax error at or near "exclusive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq execute
< error
| 42601 syntax error at or near "execute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq execute 1
< error
| 42601 syntax error at or near "execute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq execute int
< error
| 42601 syntax error at or near "execute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq execute BY NONE
< error
| 42601 syntax error at or near "execute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq execute NAME omni_probe_name
< error
| 42601 syntax error at or near "execute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq execute CYCLE
< error
| 42601 syntax error at or near "execute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exists
< error
| 42601 syntax error at or near "exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exists 1
< error
| 42601 syntax error at or near "exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exists int
< error
| 42601 syntax error at or near "exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exists BY NONE
< error
| 42601 syntax error at or near "exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exists NAME omni_probe_name
< error
| 42601 syntax error at or near "exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq exists CYCLE
< error
| 42601 syntax error at or near "exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq explain
< error
| 42601 syntax error at or near "explain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq explain 1
< error
| 42601 syntax error at or near "explain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq explain int
< error
| 42601 syntax error at or near "explain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq explain BY NONE
< error
| 42601 syntax error at or near "explain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq explain NAME omni_probe_name
< error
| 42601 syntax error at or near "explain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq explain CYCLE
< error
| 42601 syntax error at or near "explain"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq expression
< error
| 42601 syntax error at or near "expression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq expression 1
< error
| 42601 syntax error at or near "expression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq expression int
< error
| 42601 syntax error at or near "expression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq expression BY NONE
< error
| 42601 syntax error at or near "expression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq expression NAME omni_probe_name
< error
| 42601 syntax error at or near "expression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq expression CYCLE
< error
| 42601 syntax error at or near "expression"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extension
< error
| 42601 syntax error at or near "extension"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extension 1
< error
| 42601 syntax error at or near "extension"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extension int
< error
| 42601 syntax error at or near "extension"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extension BY NONE
< error
| 42601 syntax error at or near "extension"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extension NAME omni_probe_name
< error
| 42601 syntax error at or near "extension"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extension CYCLE
< error
| 42601 syntax error at or near "extension"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq external
< error
| 42601 syntax error at or near "external"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq external 1
< error
| 42601 syntax error at or near "external"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq external int
< error
| 42601 syntax error at or near "external"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq external BY NONE
< error
| 42601 syntax error at or near "external"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq external NAME omni_probe_name
< error
| 42601 syntax error at or near "external"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq external CYCLE
< error
| 42601 syntax error at or near "external"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extract
< error
| 42601 syntax error at or near "extract"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extract 1
< error
| 42601 syntax error at or near "extract"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extract int
< error
| 42601 syntax error at or near "extract"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extract BY NONE
< error
| 42601 syntax error at or near "extract"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extract NAME omni_probe_name
< error
| 42601 syntax error at or near "extract"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq extract CYCLE
< error
| 42601 syntax error at or near "extract"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq false
< error
| 42601 syntax error at or near "false"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq false 1
< error
| 42601 syntax error at or near "false"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq false int
< error
| 42601 syntax error at or near "false"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq false BY NONE
< error
| 42601 syntax error at or near "false"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq false NAME omni_probe_name
< error
| 42601 syntax error at or near "false"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq false CYCLE
< error
| 42601 syntax error at or near "false"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq family
< error
| 42601 syntax error at or near "family"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq family 1
< error
| 42601 syntax error at or near "family"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq family int
< error
| 42601 syntax error at or near "family"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq family BY NONE
< error
| 42601 syntax error at or near "family"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq family NAME omni_probe_name
< error
| 42601 syntax error at or near "family"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq family CYCLE
< error
| 42601 syntax error at or near "family"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq fetch
< error
| 42601 syntax error at or near "fetch"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq fetch 1
< error
| 42601 syntax error at or near "fetch"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq fetch int
< error
| 42601 syntax error at or near "fetch"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq fetch BY NONE
< error
| 42601 syntax error at or near "fetch"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq fetch NAME omni_probe_name
< error
| 42601 syntax error at or near "fetch"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq fetch CYCLE
< error
| 42601 syntax error at or near "fetch"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq filter
< error
| 42601 syntax error at or near "filter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq filter 1
< error
| 42601 syntax error at or near "filter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq filter int
< error
| 42601 syntax error at or near "filter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq filter BY NONE
< error
| 42601 syntax error at or near "filter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq filter NAME omni_probe_name
< error
| 42601 syntax error at or near "filter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq filter CYCLE
< error
| 42601 syntax error at or near "filter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq finalize
< error
| 42601 syntax error at or near "finalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq finalize 1
< error
| 42601 syntax error at or near "finalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq finalize int
< error
| 42601 syntax error at or near "finalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq finalize BY NONE
< error
| 42601 syntax error at or near "finalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq finalize NAME omni_probe_name
< error
| 42601 syntax error at or near "finalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq finalize CYCLE
< error
| 42601 syntax error at or near "finalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq first
< error
| 42601 syntax error at or near "first"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq first 1
< error
| 42601 syntax error at or near "first"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq first int
< error
| 42601 syntax error at or near "first"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq first BY NONE
< error
| 42601 syntax error at or near "first"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq first NAME omni_probe_name
< error
| 42601 syntax error at or near "first"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq first CYCLE
< error
| 42601 syntax error at or near "first"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq float
< error
| 42601 syntax error at or near "float"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq float 1
< error
| 42601 syntax error at or near "float"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq float int
< error
| 42601 syntax error at or near "float"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq float BY NONE
< error
| 42601 syntax error at or near "float"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq float NAME omni_probe_name
< error
| 42601 syntax error at or near "float"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq float CYCLE
< error
| 42601 syntax error at or near "float"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq following
< error
| 42601 syntax error at or near "following"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq following 1
< error
| 42601 syntax error at or near "following"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq following int
< error
| 42601 syntax error at or near "following"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq following BY NONE
< error
| 42601 syntax error at or near "following"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq following NAME omni_probe_name
< error
| 42601 syntax error at or near "following"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq following CYCLE
< error
| 42601 syntax error at or near "following"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq for
< error
| 42601 syntax error at or near "for"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq for 1
< error
| 42601 syntax error at or near "for"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq for int
< error
| 42601 syntax error at or near "for"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq for BY NONE
< error
| 42601 syntax error at or near "for"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq for NAME omni_probe_name
< error
| 42601 syntax error at or near "for"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq for CYCLE
< error
| 42601 syntax error at or near "for"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq force
< error
| 42601 syntax error at or near "force"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq force 1
< error
| 42601 syntax error at or near "force"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq force int
< error
| 42601 syntax error at or near "force"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq force BY NONE
< error
| 42601 syntax error at or near "force"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq force NAME omni_probe_name
< error
| 42601 syntax error at or near "force"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq force CYCLE
< error
| 42601 syntax error at or near "force"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq foreign
< error
| 42601 syntax error at or near "foreign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq foreign 1
< error
| 42601 syntax error at or near "foreign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq foreign int
< error
| 42601 syntax error at or near "foreign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq foreign BY NONE
< error
| 42601 syntax error at or near "foreign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq foreign NAME omni_probe_name
< error
| 42601 syntax error at or near "foreign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq foreign CYCLE
< error
| 42601 syntax error at or near "foreign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq format
< error
| 42601 syntax error at or near "format"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq format 1
< error
| 42601 syntax error at or near "format"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq format int
< error
| 42601 syntax error at or near "format"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq format BY NONE
< error
| 42601 syntax error at or near "format"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq format NAME omni_probe_name
< error
| 42601 syntax error at or near "format"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq format CYCLE
< error
| 42601 syntax error at or near "format"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq forward
< error
| 42601 syntax error at or near "forward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq forward 1
< error
| 42601 syntax error at or near "forward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq forward int
< error
| 42601 syntax error at or near "forward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq forward BY NONE
< error
| 42601 syntax error at or near "forward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq forward NAME omni_probe_name
< error
| 42601 syntax error at or near "forward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq forward CYCLE
< error
| 42601 syntax error at or near "forward"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq freeze
< error
| 42601 syntax error at or near "freeze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq freeze 1
< error
| 42601 syntax error at or near "freeze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq freeze int
< error
| 42601 syntax error at or near "freeze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq freeze BY NONE
< error
| 42601 syntax error at or near "freeze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq freeze NAME omni_probe_name
< error
| 42601 syntax error at or near "freeze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq freeze CYCLE
< error
| 42601 syntax error at or near "freeze"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq from
< error
| 42601 syntax error at or near "from"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq from 1
< error
| 42601 syntax error at or near "from"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq from int
< error
| 42601 syntax error at or near "from"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq from BY NONE
< error
| 42601 syntax error at or near "from"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq from NAME omni_probe_name
< error
| 42601 syntax error at or near "from"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq from CYCLE
< error
| 42601 syntax error at or near "from"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq full
< error
| 42601 syntax error at or near "full"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq full 1
< error
| 42601 syntax error at or near "full"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq full int
< error
| 42601 syntax error at or near "full"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq full BY NONE
< error
| 42601 syntax error at or near "full"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq full NAME omni_probe_name
< error
| 42601 syntax error at or near "full"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq full CYCLE
< error
| 42601 syntax error at or near "full"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq function
< error
| 42601 syntax error at or near "function"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq function 1
< error
| 42601 syntax error at or near "function"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq function int
< error
| 42601 syntax error at or near "function"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq function BY NONE
< error
| 42601 syntax error at or near "function"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq function NAME omni_probe_name
< error
| 42601 syntax error at or near "function"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq function CYCLE
< error
| 42601 syntax error at or near "function"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq functions
< error
| 42601 syntax error at or near "functions"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq functions 1
< error
| 42601 syntax error at or near "functions"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq functions int
< error
| 42601 syntax error at or near "functions"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq functions BY NONE
< error
| 42601 syntax error at or near "functions"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq functions NAME omni_probe_name
< error
| 42601 syntax error at or near "functions"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq functions CYCLE
< error
| 42601 syntax error at or near "functions"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq generated
< error
| 42601 syntax error at or near "generated"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq generated 1
< error
| 42601 syntax error at or near "generated"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq generated int
< error
| 42601 syntax error at or near "generated"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq generated BY NONE
< error
| 42601 syntax error at or near "generated"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq generated NAME omni_probe_name
< error
| 42601 syntax error at or near "generated"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq generated CYCLE
< error
| 42601 syntax error at or near "generated"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq global
< error
| 42601 syntax error at or near "global"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq global 1
< error
| 42601 syntax error at or near "global"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq global int
< error
| 42601 syntax error at or near "global"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq global BY NONE
< error
| 42601 syntax error at or near "global"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq global NAME omni_probe_name
< error
| 42601 syntax error at or near "global"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq global CYCLE
< error
| 42601 syntax error at or near "global"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grant
< error
| 42601 syntax error at or near "grant"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grant 1
< error
| 42601 syntax error at or near "grant"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grant int
< error
| 42601 syntax error at or near "grant"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grant BY NONE
< error
| 42601 syntax error at or near "grant"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grant NAME omni_probe_name
< error
| 42601 syntax error at or near "grant"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grant CYCLE
< error
| 42601 syntax error at or near "grant"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq granted
< error
| 42601 syntax error at or near "granted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq granted 1
< error
| 42601 syntax error at or near "granted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq granted int
< error
| 42601 syntax error at or near "granted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq granted BY NONE
< error
| 42601 syntax error at or near "granted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq granted NAME omni_probe_name
< error
| 42601 syntax error at or near "granted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq granted CYCLE
< error
| 42601 syntax error at or near "granted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq greatest
< error
| 42601 syntax error at or near "greatest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq greatest 1
< error
| 42601 syntax error at or near "greatest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq greatest int
< error
| 42601 syntax error at or near "greatest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq greatest BY NONE
< error
| 42601 syntax error at or near "greatest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq greatest NAME omni_probe_name
< error
| 42601 syntax error at or near "greatest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq greatest CYCLE
< error
| 42601 syntax error at or near "greatest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq group
< error
| 42601 syntax error at or near "group"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq group 1
< error
| 42601 syntax error at or near "group"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq group int
< error
| 42601 syntax error at or near "group"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq group BY NONE
< error
| 42601 syntax error at or near "group"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq group NAME omni_probe_name
< error
| 42601 syntax error at or near "group"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq group CYCLE
< error
| 42601 syntax error at or near "group"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grouping
< error
| 42601 syntax error at or near "grouping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grouping 1
< error
| 42601 syntax error at or near "grouping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grouping int
< error
| 42601 syntax error at or near "grouping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grouping BY NONE
< error
| 42601 syntax error at or near "grouping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grouping NAME omni_probe_name
< error
| 42601 syntax error at or near "grouping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq grouping CYCLE
< error
| 42601 syntax error at or near "grouping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq groups
< error
| 42601 syntax error at or near "groups"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq groups 1
< error
| 42601 syntax error at or near "groups"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq groups int
< error
| 42601 syntax error at or near "groups"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq groups BY NONE
< error
| 42601 syntax error at or near "groups"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq groups NAME omni_probe_name
< error
| 42601 syntax error at or near "groups"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq groups CYCLE
< error
| 42601 syntax error at or near "groups"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq handler
< error
| 42601 syntax error at or near "handler"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq handler 1
< error
| 42601 syntax error at or near "handler"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq handler int
< error
| 42601 syntax error at or near "handler"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq handler BY NONE
< error
| 42601 syntax error at or near "handler"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq handler NAME omni_probe_name
< error
| 42601 syntax error at or near "handler"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq handler CYCLE
< error
| 42601 syntax error at or near "handler"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq having
< error
| 42601 syntax error at or near "having"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq having 1
< error
| 42601 syntax error at or near "having"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq having int
< error
| 42601 syntax error at or near "having"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq having BY NONE
< error
| 42601 syntax error at or near "having"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq having NAME omni_probe_name
< error
| 42601 syntax error at or near "having"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq having CYCLE
< error
| 42601 syntax error at or near "having"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq header
< error
| 42601 syntax error at or near "header"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq header 1
< error
| 42601 syntax error at or near "header"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq header int
< error
| 42601 syntax error at or near "header"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq header BY NONE
< error
| 42601 syntax error at or near "header"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq header NAME omni_probe_name
< error
| 42601 syntax error at or near "header"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq header CYCLE
< error
| 42601 syntax error at or near "header"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hold
< error
| 42601 syntax error at or near "hold"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hold 1
< error
| 42601 syntax error at or near "hold"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hold int
< error
| 42601 syntax error at or near "hold"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hold BY NONE
< error
| 42601 syntax error at or near "hold"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hold NAME omni_probe_name
< error
| 42601 syntax error at or near "hold"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hold CYCLE
< error
| 42601 syntax error at or near "hold"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hour
< error
| 42601 syntax error at or near "hour"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hour 1
< error
| 42601 syntax error at or near "hour"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hour int
< error
| 42601 syntax error at or near "hour"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hour BY NONE
< error
| 42601 syntax error at or near "hour"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hour NAME omni_probe_name
< error
| 42601 syntax error at or near "hour"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq hour CYCLE
< error
| 42601 syntax error at or near "hour"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq identity
< error
| 42601 syntax error at or near "identity"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq identity 1
< error
| 42601 syntax error at or near "identity"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq identity int
< error
| 42601 syntax error at or near "identity"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq identity BY NONE
< error
| 42601 syntax error at or near "identity"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq identity NAME omni_probe_name
< error
| 42601 syntax error at or near "identity"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq identity CYCLE
< error
| 42601 syntax error at or near "identity"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq if
< error
| 42601 syntax error at or near "if"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq if 1
< error
| 42601 syntax error at or near "if"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq if int
< error
| 42601 syntax error at or near "if"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq if BY NONE
< error
| 42601 syntax error at or near "if"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq if NAME omni_probe_name
< error
| 42601 syntax error at or near "if"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq if CYCLE
< error
| 42601 syntax error at or near "if"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ilike
< error
| 42601 syntax error at or near "ilike"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ilike 1
< error
| 42601 syntax error at or near "ilike"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ilike int
< error
| 42601 syntax error at or near "ilike"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ilike BY NONE
< error
| 42601 syntax error at or near "ilike"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ilike NAME omni_probe_name
< error
| 42601 syntax error at or near "ilike"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ilike CYCLE
< error
| 42601 syntax error at or near "ilike"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immediate
< error
| 42601 syntax error at or near "immediate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immediate 1
< error
| 42601 syntax error at or near "immediate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immediate int
< error
| 42601 syntax error at or near "immediate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immediate BY NONE
< error
| 42601 syntax error at or near "immediate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immediate NAME omni_probe_name
< error
| 42601 syntax error at or near "immediate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immediate CYCLE
< error
| 42601 syntax error at or near "immediate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immutable
< error
| 42601 syntax error at or near "immutable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immutable 1
< error
| 42601 syntax error at or near "immutable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immutable int
< error
| 42601 syntax error at or near "immutable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immutable BY NONE
< error
| 42601 syntax error at or near "immutable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immutable NAME omni_probe_name
< error
| 42601 syntax error at or near "immutable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq immutable CYCLE
< error
| 42601 syntax error at or near "immutable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq implicit
< error
| 42601 syntax error at or near "implicit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq implicit 1
< error
| 42601 syntax error at or near "implicit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq implicit int
< error
| 42601 syntax error at or near "implicit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq implicit BY NONE
< error
| 42601 syntax error at or near "implicit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq implicit NAME omni_probe_name
< error
| 42601 syntax error at or near "implicit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq implicit CYCLE
< error
| 42601 syntax error at or near "implicit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq import
< error
| 42601 syntax error at or near "import"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq import 1
< error
| 42601 syntax error at or near "import"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq import int
< error
| 42601 syntax error at or near "import"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq import BY NONE
< error
| 42601 syntax error at or near "import"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq import NAME omni_probe_name
< error
| 42601 syntax error at or near "import"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq import CYCLE
< error
| 42601 syntax error at or near "import"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq in
< error
| 42601 syntax error at or near "in"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq in 1
< error
| 42601 syntax error at or near "in"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq in int
< error
| 42601 syntax error at or near "in"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq in BY NONE
< error
| 42601 syntax error at or near "in"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq in NAME omni_probe_name
< error
| 42601 syntax error at or near "in"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq in CYCLE
< error
| 42601 syntax error at or near "in"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq include
< error
| 42601 syntax error at or near "include"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq include 1
< error
| 42601 syntax error at or near "include"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq include int
< error
| 42601 syntax error at or near "include"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq include BY NONE
< error
| 42601 syntax error at or near "include"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq include NAME omni_probe_name
< error
| 42601 syntax error at or near "include"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq include CYCLE
< error
| 42601 syntax error at or near "include"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq including
< error
| 42601 syntax error at or near "including"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq including 1
< error
| 42601 syntax error at or near "including"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq including int
< error
| 42601 syntax error at or near "including"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq including BY NONE
< error
| 42601 syntax error at or near "including"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq including NAME omni_probe_name
< error
| 42601 syntax error at or near "including"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq including CYCLE
< error
| 42601 syntax error at or near "including"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq increment
< error
| 42601 syntax error at end of input
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq increment 1
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indent
< error
| 42601 syntax error at or near "indent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indent 1
< error
| 42601 syntax error at or near "indent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indent int
< error
| 42601 syntax error at or near "indent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indent BY NONE
< error
| 42601 syntax error at or near "indent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indent NAME omni_probe_name
< error
| 42601 syntax error at or near "indent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indent CYCLE
< error
| 42601 syntax error at or near "indent"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq index
< error
| 42601 syntax error at or near "index"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq index 1
< error
| 42601 syntax error at or near "index"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq index int
< error
| 42601 syntax error at or near "index"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq index BY NONE
< error
| 42601 syntax error at or near "index"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq index NAME omni_probe_name
< error
| 42601 syntax error at or near "index"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq index CYCLE
< error
| 42601 syntax error at or near "index"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indexes
< error
| 42601 syntax error at or near "indexes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indexes 1
< error
| 42601 syntax error at or near "indexes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indexes int
< error
| 42601 syntax error at or near "indexes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indexes BY NONE
< error
| 42601 syntax error at or near "indexes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indexes NAME omni_probe_name
< error
| 42601 syntax error at or near "indexes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq indexes CYCLE
< error
| 42601 syntax error at or near "indexes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherit
< error
| 42601 syntax error at or near "inherit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherit 1
< error
| 42601 syntax error at or near "inherit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherit int
< error
| 42601 syntax error at or near "inherit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherit BY NONE
< error
| 42601 syntax error at or near "inherit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherit NAME omni_probe_name
< error
| 42601 syntax error at or near "inherit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherit CYCLE
< error
| 42601 syntax error at or near "inherit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherits
< error
| 42601 syntax error at or near "inherits"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherits 1
< error
| 42601 syntax error at or near "inherits"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherits int
< error
| 42601 syntax error at or near "inherits"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherits BY NONE
< error
| 42601 syntax error at or near "inherits"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherits NAME omni_probe_name
< error
| 42601 syntax error at or near "inherits"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inherits CYCLE
< error
| 42601 syntax error at or near "inherits"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq initially
< error
| 42601 syntax error at or near "initially"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq initially 1
< error
| 42601 syntax error at or near "initially"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq initially int
< error
| 42601 syntax error at or near "initially"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq initially BY NONE
< error
| 42601 syntax error at or near "initially"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq initially NAME omni_probe_name
< error
| 42601 syntax error at or near "initially"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq initially CYCLE
< error
| 42601 syntax error at or near "initially"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inline
< error
| 42601 syntax error at or near "inline"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inline 1
< error
| 42601 syntax error at or near "inline"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inline int
< error
| 42601 syntax error at or near "inline"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inline BY NONE
< error
| 42601 syntax error at or near "inline"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inline NAME omni_probe_name
< error
| 42601 syntax error at or near "inline"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inline CYCLE
< error
| 42601 syntax error at or near "inline"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inner
< error
| 42601 syntax error at or near "inner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inner 1
< error
| 42601 syntax error at or near "inner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inner int
< error
| 42601 syntax error at or near "inner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inner BY NONE
< error
| 42601 syntax error at or near "inner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inner NAME omni_probe_name
< error
| 42601 syntax error at or near "inner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inner CYCLE
< error
| 42601 syntax error at or near "inner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inout
< error
| 42601 syntax error at or near "inout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inout 1
< error
| 42601 syntax error at or near "inout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inout int
< error
| 42601 syntax error at or near "inout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inout BY NONE
< error
| 42601 syntax error at or near "inout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inout NAME omni_probe_name
< error
| 42601 syntax error at or near "inout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq inout CYCLE
< error
| 42601 syntax error at or near "inout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq input
< error
| 42601 syntax error at or near "input"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq input 1
< error
| 42601 syntax error at or near "input"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq input int
< error
| 42601 syntax error at or near "input"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq input BY NONE
< error
| 42601 syntax error at or near "input"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq input NAME omni_probe_name
< error
| 42601 syntax error at or near "input"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq input CYCLE
< error
| 42601 syntax error at or near "input"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insensitive
< error
| 42601 syntax error at or near "insensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insensitive 1
< error
| 42601 syntax error at or near "insensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insensitive int
< error
| 42601 syntax error at or near "insensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insensitive BY NONE
< error
| 42601 syntax error at or near "insensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insensitive NAME omni_probe_name
< error
| 42601 syntax error at or near "insensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insensitive CYCLE
< error
| 42601 syntax error at or near "insensitive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insert
< error
| 42601 syntax error at or near "insert"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insert 1
< error
| 42601 syntax error at or near "insert"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insert int
< error
| 42601 syntax error at or near "insert"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insert BY NONE
< error
| 42601 syntax error at or near "insert"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insert NAME omni_probe_name
< error
| 42601 syntax error at or near "insert"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq insert CYCLE
< error
| 42601 syntax error at or near "insert"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq instead
< error
| 42601 syntax error at or near "instead"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq instead 1
< error
| 42601 syntax error at or near "instead"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq instead int
< error
| 42601 syntax error at or near "instead"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq instead BY NONE
< error
| 42601 syntax error at or near "instead"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq instead NAME omni_probe_name
< error
| 42601 syntax error at or near "instead"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq instead CYCLE
< error
| 42601 syntax error at or near "instead"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq int
< error
| 42601 syntax error at or near "int"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq int 1
< error
| 42601 syntax error at or near "int"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq int int
< error
| 42601 syntax error at or near "int"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq int BY NONE
< error
| 42601 syntax error at or near "int"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq int NAME omni_probe_name
< error
| 42601 syntax error at or near "int"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq int CYCLE
< error
| 42601 syntax error at or near "int"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq integer
< error
| 42601 syntax error at or near "integer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq integer 1
< error
| 42601 syntax error at or near "integer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq integer int
< error
| 42601 syntax error at or near "integer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq integer BY NONE
< error
| 42601 syntax error at or near "integer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq integer NAME omni_probe_name
< error
| 42601 syntax error at or near "integer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq integer CYCLE
< error
| 42601 syntax error at or near "integer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq intersect
< error
| 42601 syntax error at or near "intersect"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq intersect 1
< error
| 42601 syntax error at or near "intersect"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq intersect int
< error
| 42601 syntax error at or near "intersect"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq intersect BY NONE
< error
| 42601 syntax error at or near "intersect"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq intersect NAME omni_probe_name
< error
| 42601 syntax error at or near "intersect"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq intersect CYCLE
< error
| 42601 syntax error at or near "intersect"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq interval
< error
| 42601 syntax error at or near "interval"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq interval 1
< error
| 42601 syntax error at or near "interval"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq interval int
< error
| 42601 syntax error at or near "interval"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq interval BY NONE
< error
| 42601 syntax error at or near "interval"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq interval NAME omni_probe_name
< error
| 42601 syntax error at or near "interval"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq interval CYCLE
< error
| 42601 syntax error at or near "interval"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq into
< error
| 42601 syntax error at or near "into"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq into 1
< error
| 42601 syntax error at or near "into"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq into int
< error
| 42601 syntax error at or near "into"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq into BY NONE
< error
| 42601 syntax error at or near "into"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq into NAME omni_probe_name
< error
| 42601 syntax error at or near "into"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq into CYCLE
< error
| 42601 syntax error at or near "into"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq invoker
< error
| 42601 syntax error at or near "invoker"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq invoker 1
< error
| 42601 syntax error at or near "invoker"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq invoker int
< error
| 42601 syntax error at or near "invoker"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq invoker BY NONE
< error
| 42601 syntax error at or near "invoker"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq invoker NAME omni_probe_name
< error
| 42601 syntax error at or near "invoker"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq invoker CYCLE
< error
| 42601 syntax error at or near "invoker"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq is
< error
| 42601 syntax error at or near "is"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq is 1
< error
| 42601 syntax error at or near "is"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq is int
< error
| 42601 syntax error at or near "is"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq is BY NONE
< error
| 42601 syntax error at or near "is"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq is NAME omni_probe_name
< error
| 42601 syntax error at or near "is"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq is CYCLE
< error
| 42601 syntax error at or near "is"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isnull
< error
| 42601 syntax error at or near "isnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isnull 1
< error
| 42601 syntax error at or near "isnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isnull int
< error
| 42601 syntax error at or near "isnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isnull BY NONE
< error
| 42601 syntax error at or near "isnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isnull NAME omni_probe_name
< error
| 42601 syntax error at or near "isnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isnull CYCLE
< error
| 42601 syntax error at or near "isnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isolation
< error
| 42601 syntax error at or near "isolation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isolation 1
< error
| 42601 syntax error at or near "isolation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isolation int
< error
| 42601 syntax error at or near "isolation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isolation BY NONE
< error
| 42601 syntax error at or near "isolation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isolation NAME omni_probe_name
< error
| 42601 syntax error at or near "isolation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq isolation CYCLE
< error
| 42601 syntax error at or near "isolation"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq join
< error
| 42601 syntax error at or near "join"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq join 1
< error
| 42601 syntax error at or near "join"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq join int
< error
| 42601 syntax error at or near "join"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq join BY NONE
< error
| 42601 syntax error at or near "join"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq join NAME omni_probe_name
< error
| 42601 syntax error at or near "join"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq join CYCLE
< error
| 42601 syntax error at or near "join"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json
< error
| 42601 syntax error at or near "json"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json 1
< error
| 42601 syntax error at or near "json"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json int
< error
| 42601 syntax error at or near "json"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json BY NONE
< error
| 42601 syntax error at or near "json"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json NAME omni_probe_name
< error
| 42601 syntax error at or near "json"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json CYCLE
< error
| 42601 syntax error at or near "json"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_array
< error
| 42601 syntax error at or near "json_array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_array 1
< error
| 42601 syntax error at or near "json_array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_array int
< error
| 42601 syntax error at or near "json_array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_array BY NONE
< error
| 42601 syntax error at or near "json_array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_array NAME omni_probe_name
< error
| 42601 syntax error at or near "json_array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_array CYCLE
< error
| 42601 syntax error at or near "json_array"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_arrayagg
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_arrayagg 1
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_arrayagg int
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_arrayagg BY NONE
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_arrayagg NAME omni_probe_name
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_arrayagg CYCLE
< error
| 42601 syntax error at or near "json_arrayagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_exists
< error
| 42601 syntax error at or near "json_exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_exists 1
< error
| 42601 syntax error at or near "json_exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_exists int
< error
| 42601 syntax error at or near "json_exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_exists BY NONE
< error
| 42601 syntax error at or near "json_exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_exists NAME omni_probe_name
< error
| 42601 syntax error at or near "json_exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_exists CYCLE
< error
| 42601 syntax error at or near "json_exists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_object
< error
| 42601 syntax error at or near "json_object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_object 1
< error
| 42601 syntax error at or near "json_object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_object int
< error
| 42601 syntax error at or near "json_object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_object BY NONE
< error
| 42601 syntax error at or near "json_object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_object NAME omni_probe_name
< error
| 42601 syntax error at or near "json_object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_object CYCLE
< error
| 42601 syntax error at or near "json_object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_objectagg
< error
| 42601 syntax error at or near "json_objectagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_objectagg 1
< error
| 42601 syntax error at or near "json_objectagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_objectagg int
< error
| 42601 syntax error at or near "json_objectagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_objectagg BY NONE
< error
| 42601 syntax error at or near "json_objectagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_objectagg NAME omni_probe_name
< error
| 42601 syntax error at or near "json_objectagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_objectagg CYCLE
< error
| 42601 syntax error at or near "json_objectagg"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_query
< error
| 42601 syntax error at or near "json_query"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_query 1
< error
| 42601 syntax error at or near "json_query"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_query int
< error
| 42601 syntax error at or near "json_query"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_query BY NONE
< error
| 42601 syntax error at or near "json_query"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_query NAME omni_probe_name
< error
| 42601 syntax error at or near "json_query"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_query CYCLE
< error
| 42601 syntax error at or near "json_query"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_scalar
< error
| 42601 syntax error at or near "json_scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_scalar 1
< error
| 42601 syntax error at or near "json_scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_scalar int
< error
| 42601 syntax error at or near "json_scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_scalar BY NONE
< error
| 42601 syntax error at or near "json_scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_scalar NAME omni_probe_name
< error
| 42601 syntax error at or near "json_scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_scalar CYCLE
< error
| 42601 syntax error at or near "json_scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_serialize
< error
| 42601 syntax error at or near "json_serialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_serialize 1
< error
| 42601 syntax error at or near "json_serialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_serialize int
< error
| 42601 syntax error at or near "json_serialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_serialize BY NONE
< error
| 42601 syntax error at or near "json_serialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_serialize NAME omni_probe_name
< error
| 42601 syntax error at or near "json_serialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_serialize CYCLE
< error
| 42601 syntax error at or near "json_serialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_table
< error
| 42601 syntax error at or near "json_table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_table 1
< error
| 42601 syntax error at or near "json_table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_table int
< error
| 42601 syntax error at or near "json_table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_table BY NONE
< error
| 42601 syntax error at or near "json_table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_table NAME omni_probe_name
< error
| 42601 syntax error at or near "json_table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_table CYCLE
< error
| 42601 syntax error at or near "json_table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_value
< error
| 42601 syntax error at or near "json_value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_value 1
< error
| 42601 syntax error at or near "json_value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_value int
< error
| 42601 syntax error at or near "json_value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_value BY NONE
< error
| 42601 syntax error at or near "json_value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_value NAME omni_probe_name
< error
| 42601 syntax error at or near "json_value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq json_value CYCLE
< error
| 42601 syntax error at or near "json_value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keep
< error
| 42601 syntax error at or near "keep"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keep 1
< error
| 42601 syntax error at or near "keep"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keep int
< error
| 42601 syntax error at or near "keep"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keep BY NONE
< error
| 42601 syntax error at or near "keep"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keep NAME omni_probe_name
< error
| 42601 syntax error at or near "keep"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keep CYCLE
< error
| 42601 syntax error at or near "keep"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq key
< error
| 42601 syntax error at or near "key"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq key 1
< error
| 42601 syntax error at or near "key"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq key int
< error
| 42601 syntax error at or near "key"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq key BY NONE
< error
| 42601 syntax error at or near "key"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq key NAME omni_probe_name
< error
| 42601 syntax error at or near "key"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq key CYCLE
< error
| 42601 syntax error at or near "key"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keys
< error
| 42601 syntax error at or near "keys"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keys 1
< error
| 42601 syntax error at or near "keys"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keys int
< error
| 42601 syntax error at or near "keys"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keys BY NONE
< error
| 42601 syntax error at or near "keys"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keys NAME omni_probe_name
< error
| 42601 syntax error at or near "keys"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq keys CYCLE
< error
| 42601 syntax error at or near "keys"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq label
< error
| 42601 syntax error at or near "label"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq label 1
< error
| 42601 syntax error at or near "label"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq label int
< error
| 42601 syntax error at or near "label"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq label BY NONE
< error
| 42601 syntax error at or near "label"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq label NAME omni_probe_name
< error
| 42601 syntax error at or near "label"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq label CYCLE
< error
| 42601 syntax error at or near "label"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq language
< error
| 42601 syntax error at or near "language"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq language 1
< error
| 42601 syntax error at or near "language"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq language int
< error
| 42601 syntax error at or near "language"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq language BY NONE
< error
| 42601 syntax error at or near "language"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq language NAME omni_probe_name
< error
| 42601 syntax error at or near "language"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq language CYCLE
< error
| 42601 syntax error at or near "language"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq large
< error
| 42601 syntax error at or near "large"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq large 1
< error
| 42601 syntax error at or near "large"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq large int
< error
| 42601 syntax error at or near "large"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq large BY NONE
< error
| 42601 syntax error at or near "large"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq large NAME omni_probe_name
< error
| 42601 syntax error at or near "large"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq large CYCLE
< error
| 42601 syntax error at or near "large"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq last
< error
| 42601 syntax error at or near "last"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq last 1
< error
| 42601 syntax error at or near "last"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq last int
< error
| 42601 syntax error at or near "last"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq last BY NONE
< error
| 42601 syntax error at or near "last"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq last NAME omni_probe_name
< error
| 42601 syntax error at or near "last"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq last CYCLE
< error
| 42601 syntax error at or near "last"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lateral
< error
| 42601 syntax error at or near "lateral"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lateral 1
< error
| 42601 syntax error at or near "lateral"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lateral int
< error
| 42601 syntax error at or near "lateral"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lateral BY NONE
< error
| 42601 syntax error at or near "lateral"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lateral NAME omni_probe_name
< error
| 42601 syntax error at or near "lateral"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lateral CYCLE
< error
| 42601 syntax error at or near "lateral"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leading
< error
| 42601 syntax error at or near "leading"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leading 1
< error
| 42601 syntax error at or near "leading"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leading int
< error
| 42601 syntax error at or near "leading"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leading BY NONE
< error
| 42601 syntax error at or near "leading"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leading NAME omni_probe_name
< error
| 42601 syntax error at or near "leading"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leading CYCLE
< error
| 42601 syntax error at or near "leading"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leakproof
< error
| 42601 syntax error at or near "leakproof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leakproof 1
< error
| 42601 syntax error at or near "leakproof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leakproof int
< error
| 42601 syntax error at or near "leakproof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leakproof BY NONE
< error
| 42601 syntax error at or near "leakproof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leakproof NAME omni_probe_name
< error
| 42601 syntax error at or near "leakproof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq leakproof CYCLE
< error
| 42601 syntax error at or near "leakproof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq least
< error
| 42601 syntax error at or near "least"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq least 1
< error
| 42601 syntax error at or near "least"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq least int
< error
| 42601 syntax error at or near "least"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq least BY NONE
< error
| 42601 syntax error at or near "least"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq least NAME omni_probe_name
< error
| 42601 syntax error at or near "least"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq least CYCLE
< error
| 42601 syntax error at or near "least"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq left
< error
| 42601 syntax error at or near "left"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq left 1
< error
| 42601 syntax error at or near "left"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq left int
< error
| 42601 syntax error at or near "left"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq left BY NONE
< error
| 42601 syntax error at or near "left"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq left NAME omni_probe_name
< error
| 42601 syntax error at or near "left"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq left CYCLE
< error
| 42601 syntax error at or near "left"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq level
< error
| 42601 syntax error at or near "level"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq level 1
< error
| 42601 syntax error at or near "level"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq level int
< error
| 42601 syntax error at or near "level"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq level BY NONE
< error
| 42601 syntax error at or near "level"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq level NAME omni_probe_name
< error
| 42601 syntax error at or near "level"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq level CYCLE
< error
| 42601 syntax error at or near "level"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq like
< error
| 42601 syntax error at or near "like"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq like 1
< error
| 42601 syntax error at or near "like"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq like int
< error
| 42601 syntax error at or near "like"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq like BY NONE
< error
| 42601 syntax error at or near "like"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq like NAME omni_probe_name
< error
| 42601 syntax error at or near "like"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq like CYCLE
< error
| 42601 syntax error at or near "like"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq limit
< error
| 42601 syntax error at or near "limit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq limit 1
< error
| 42601 syntax error at or near "limit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq limit int
< error
| 42601 syntax error at or near "limit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq limit BY NONE
< error
| 42601 syntax error at or near "limit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq limit NAME omni_probe_name
< error
| 42601 syntax error at or near "limit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq limit CYCLE
< error
| 42601 syntax error at or near "limit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq listen
< error
| 42601 syntax error at or near "listen"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq listen 1
< error
| 42601 syntax error at or near "listen"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq listen int
< error
| 42601 syntax error at or near "listen"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq listen BY NONE
< error
| 42601 syntax error at or near "listen"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq listen NAME omni_probe_name
< error
| 42601 syntax error at or near "listen"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq listen CYCLE
< error
| 42601 syntax error at or near "listen"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq load
< error
| 42601 syntax error at or near "load"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq load 1
< error
| 42601 syntax error at or near "load"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq load int
< error
| 42601 syntax error at or near "load"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq load BY NONE
< error
| 42601 syntax error at or near "load"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq load NAME omni_probe_name
< error
| 42601 syntax error at or near "load"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq load CYCLE
< error
| 42601 syntax error at or near "load"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq local
< error
| 42601 syntax error at or near "local"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq local 1
< error
| 42601 syntax error at or near "local"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq local int
< error
| 42601 syntax error at or near "local"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq local BY NONE
< error
| 42601 syntax error at or near "local"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq local NAME omni_probe_name
< error
| 42601 syntax error at or near "local"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq local CYCLE
< error
| 42601 syntax error at or near "local"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtime
< error
| 42601 syntax error at or near "localtime"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtime 1
< error
| 42601 syntax error at or near "localtime"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtime int
< error
| 42601 syntax error at or near "localtime"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtime BY NONE
< error
| 42601 syntax error at or near "localtime"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtime NAME omni_probe_name
< error
| 42601 syntax error at or near "localtime"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtime CYCLE
< error
| 42601 syntax error at or near "localtime"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtimestamp
< error
| 42601 syntax error at or near "localtimestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtimestamp 1
< error
| 42601 syntax error at or near "localtimestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtimestamp int
< error
| 42601 syntax error at or near "localtimestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtimestamp BY NONE
< error
| 42601 syntax error at or near "localtimestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtimestamp NAME omni_probe_name
< error
| 42601 syntax error at or near "localtimestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq localtimestamp CYCLE
< error
| 42601 syntax error at or near "localtimestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq location
< error
| 42601 syntax error at or near "location"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq location 1
< error
| 42601 syntax error at or near "location"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq location int
< error
| 42601 syntax error at or near "location"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq location BY NONE
< error
| 42601 syntax error at or near "location"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq location NAME omni_probe_name
< error
| 42601 syntax error at or near "location"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq location CYCLE
< error
| 42601 syntax error at or near "location"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lock
< error
| 42601 syntax error at or near "lock"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lock 1
< error
| 42601 syntax error at or near "lock"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lock int
< error
| 42601 syntax error at or near "lock"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lock BY NONE
< error
| 42601 syntax error at or near "lock"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lock NAME omni_probe_name
< error
| 42601 syntax error at or near "lock"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq lock CYCLE
< error
| 42601 syntax error at or near "lock"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq locked
< error
| 42601 syntax error at or near "locked"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq locked 1
< error
| 42601 syntax error at or near "locked"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq locked int
< error
| 42601 syntax error at or near "locked"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq locked BY NONE
< error
| 42601 syntax error at or near "locked"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq locked NAME omni_probe_name
< error
| 42601 syntax error at or near "locked"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq locked CYCLE
< error
| 42601 syntax error at or near "locked"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq logged
< error
| XX000 option "logged" not recognized
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mapping
< error
| 42601 syntax error at or near "mapping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mapping 1
< error
| 42601 syntax error at or near "mapping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mapping int
< error
| 42601 syntax error at or near "mapping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mapping BY NONE
< error
| 42601 syntax error at or near "mapping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mapping NAME omni_probe_name
< error
| 42601 syntax error at or near "mapping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mapping CYCLE
< error
| 42601 syntax error at or near "mapping"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq match
< error
| 42601 syntax error at or near "match"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq match 1
< error
| 42601 syntax error at or near "match"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq match int
< error
| 42601 syntax error at or near "match"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq match BY NONE
< error
| 42601 syntax error at or near "match"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq match NAME omni_probe_name
< error
| 42601 syntax error at or near "match"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq match CYCLE
< error
| 42601 syntax error at or near "match"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq matched
< error
| 42601 syntax error at or near "matched"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq matched 1
< error
| 42601 syntax error at or near "matched"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq matched int
< error
| 42601 syntax error at or near "matched"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq matched BY NONE
< error
| 42601 syntax error at or near "matched"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq matched NAME omni_probe_name
< error
| 42601 syntax error at or near "matched"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq matched CYCLE
< error
| 42601 syntax error at or near "matched"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq materialized
< error
| 42601 syntax error at or near "materialized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq materialized 1
< error
| 42601 syntax error at or near "materialized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq materialized int
< error
| 42601 syntax error at or near "materialized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq materialized BY NONE
< error
| 42601 syntax error at or near "materialized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq materialized NAME omni_probe_name
< error
| 42601 syntax error at or near "materialized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq materialized CYCLE
< error
| 42601 syntax error at or near "materialized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq maxvalue
< error
| 42601 syntax error at end of input
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq maxvalue 1
< error
| 22023 MINVALUE (1) must be less than MAXVALUE (1)
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge
< error
| 42601 syntax error at or near "merge"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge 1
< error
| 42601 syntax error at or near "merge"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge int
< error
| 42601 syntax error at or near "merge"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge BY NONE
< error
| 42601 syntax error at or near "merge"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge NAME omni_probe_name
< error
| 42601 syntax error at or near "merge"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge CYCLE
< error
| 42601 syntax error at or near "merge"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge_action
< error
| 42601 syntax error at or near "merge_action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge_action 1
< error
| 42601 syntax error at or near "merge_action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge_action int
< error
| 42601 syntax error at or near "merge_action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge_action BY NONE
< error
| 42601 syntax error at or near "merge_action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge_action NAME omni_probe_name
< error
| 42601 syntax error at or near "merge_action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq merge_action CYCLE
< error
| 42601 syntax error at or near "merge_action"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq method
< error
| 42601 syntax error at or near "method"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq method 1
< error
| 42601 syntax error at or near "method"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq method int
< error
| 42601 syntax error at or near "method"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq method BY NONE
< error
| 42601 syntax error at or near "method"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq method NAME omni_probe_name
< error
| 42601 syntax error at or near "method"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq method CYCLE
< error
| 42601 syntax error at or near "method"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq minute
< error
| 42601 syntax error at or near "minute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq minute 1
< error
| 42601 syntax error at or near "minute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq minute int
< error
| 42601 syntax error at or near "minute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq minute BY NONE
< error
| 42601 syntax error at or near "minute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq minute NAME omni_probe_name
< error
| 42601 syntax error at or near "minute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq minute CYCLE
< error
| 42601 syntax error at or near "minute"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq minvalue
< error
| 42601 syntax error at end of input
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq minvalue 1
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mode
< error
| 42601 syntax error at or near "mode"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mode 1
< error
| 42601 syntax error at or near "mode"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mode int
< error
| 42601 syntax error at or near "mode"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mode BY NONE
< error
| 42601 syntax error at or near "mode"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mode NAME omni_probe_name
< error
| 42601 syntax error at or near "mode"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq mode CYCLE
< error
| 42601 syntax error at or near "mode"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq month
< error
| 42601 syntax error at or near "month"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq month 1
< error
| 42601 syntax error at or near "month"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq month int
< error
| 42601 syntax error at or near "month"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq month BY NONE
< error
| 42601 syntax error at or near "month"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq month NAME omni_probe_name
< error
| 42601 syntax error at or near "month"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq month CYCLE
< error
| 42601 syntax error at or near "month"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq move
< error
| 42601 syntax error at or near "move"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq move 1
< error
| 42601 syntax error at or near "move"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq move int
< error
| 42601 syntax error at or near "move"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq move BY NONE
< error
| 42601 syntax error at or near "move"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq move NAME omni_probe_name
< error
| 42601 syntax error at or near "move"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq move CYCLE
< error
| 42601 syntax error at or near "move"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq name
< error
| 42601 syntax error at or near "name"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq name 1
< error
| 42601 syntax error at or near "name"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq name int
< error
| 42601 syntax error at or near "name"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq name BY NONE
< error
| 42601 syntax error at or near "name"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq name NAME omni_probe_name
< error
| 42601 syntax error at or near "name"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq name CYCLE
< error
| 42601 syntax error at or near "name"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq names
< error
| 42601 syntax error at or near "names"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq names 1
< error
| 42601 syntax error at or near "names"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq names int
< error
| 42601 syntax error at or near "names"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq names BY NONE
< error
| 42601 syntax error at or near "names"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq names NAME omni_probe_name
< error
| 42601 syntax error at or near "names"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq names CYCLE
< error
| 42601 syntax error at or near "names"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq national
< error
| 42601 syntax error at or near "national"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq national 1
< error
| 42601 syntax error at or near "national"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq national int
< error
| 42601 syntax error at or near "national"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq national BY NONE
< error
| 42601 syntax error at or near "national"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq national NAME omni_probe_name
< error
| 42601 syntax error at or near "national"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq national CYCLE
< error
| 42601 syntax error at or near "national"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq natural
< error
| 42601 syntax error at or near "natural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq natural 1
< error
| 42601 syntax error at or near "natural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq natural int
< error
| 42601 syntax error at or near "natural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq natural BY NONE
< error
| 42601 syntax error at or near "natural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq natural NAME omni_probe_name
< error
| 42601 syntax error at or near "natural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq natural CYCLE
< error
| 42601 syntax error at or near "natural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nchar
< error
| 42601 syntax error at or near "nchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nchar 1
< error
| 42601 syntax error at or near "nchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nchar int
< error
| 42601 syntax error at or near "nchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nchar BY NONE
< error
| 42601 syntax error at or near "nchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nchar NAME omni_probe_name
< error
| 42601 syntax error at or near "nchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nchar CYCLE
< error
| 42601 syntax error at or near "nchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nested
< error
| 42601 syntax error at or near "nested"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nested 1
< error
| 42601 syntax error at or near "nested"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nested int
< error
| 42601 syntax error at or near "nested"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nested BY NONE
< error
| 42601 syntax error at or near "nested"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nested NAME omni_probe_name
< error
| 42601 syntax error at or near "nested"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nested CYCLE
< error
| 42601 syntax error at or near "nested"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq new
< error
| 42601 syntax error at or near "new"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq new 1
< error
| 42601 syntax error at or near "new"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq new int
< error
| 42601 syntax error at or near "new"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq new BY NONE
< error
| 42601 syntax error at or near "new"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq new NAME omni_probe_name
< error
| 42601 syntax error at or near "new"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq new CYCLE
< error
| 42601 syntax error at or near "new"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq next
< error
| 42601 syntax error at or near "next"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq next 1
< error
| 42601 syntax error at or near "next"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq next int
< error
| 42601 syntax error at or near "next"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq next BY NONE
< error
| 42601 syntax error at or near "next"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq next NAME omni_probe_name
< error
| 42601 syntax error at or near "next"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq next CYCLE
< error
| 42601 syntax error at or near "next"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfc
< error
| 42601 syntax error at or near "nfc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfc 1
< error
| 42601 syntax error at or near "nfc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfc int
< error
| 42601 syntax error at or near "nfc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfc BY NONE
< error
| 42601 syntax error at or near "nfc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfc NAME omni_probe_name
< error
| 42601 syntax error at or near "nfc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfc CYCLE
< error
| 42601 syntax error at or near "nfc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfd
< error
| 42601 syntax error at or near "nfd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfd 1
< error
| 42601 syntax error at or near "nfd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfd int
< error
| 42601 syntax error at or near "nfd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfd BY NONE
< error
| 42601 syntax error at or near "nfd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfd NAME omni_probe_name
< error
| 42601 syntax error at or near "nfd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfd CYCLE
< error
| 42601 syntax error at or near "nfd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkc
< error
| 42601 syntax error at or near "nfkc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkc 1
< error
| 42601 syntax error at or near "nfkc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkc int
< error
| 42601 syntax error at or near "nfkc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkc BY NONE
< error
| 42601 syntax error at or near "nfkc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkc NAME omni_probe_name
< error
| 42601 syntax error at or near "nfkc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkc CYCLE
< error
| 42601 syntax error at or near "nfkc"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkd
< error
| 42601 syntax error at or near "nfkd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkd 1
< error
| 42601 syntax error at or near "nfkd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkd int
< error
| 42601 syntax error at or near "nfkd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkd BY NONE
< error
| 42601 syntax error at or near "nfkd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkd NAME omni_probe_name
< error
| 42601 syntax error at or near "nfkd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nfkd CYCLE
< error
| 42601 syntax error at or near "nfkd"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq no
< error
| 42601 syntax error at end of input
| position: 34
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq no 1
< error
| 42601 syntax error at or near "1"
| position: 35
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq no int
< error
| 42601 syntax error at or near "int"
| position: 35
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq no BY NONE
< error
| 42601 syntax error at or near "BY"
| position: 35
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq no NAME omni_probe_name
< error
| 42601 syntax error at or near "NAME"
| position: 35
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq no CYCLE
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq none
< error
| 42601 syntax error at or near "none"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq none 1
< error
| 42601 syntax error at or near "none"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq none int
< error
| 42601 syntax error at or near "none"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq none BY NONE
< error
| 42601 syntax error at or near "none"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq none NAME omni_probe_name
< error
| 42601 syntax error at or near "none"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq none CYCLE
< error
| 42601 syntax error at or near "none"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalize
< error
| 42601 syntax error at or near "normalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalize 1
< error
| 42601 syntax error at or near "normalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalize int
< error
| 42601 syntax error at or near "normalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalize BY NONE
< error
| 42601 syntax error at or near "normalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalize NAME omni_probe_name
< error
| 42601 syntax error at or near "normalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalize CYCLE
< error
| 42601 syntax error at or near "normalize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalized
< error
| 42601 syntax error at or near "normalized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalized 1
< error
| 42601 syntax error at or near "normalized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalized int
< error
| 42601 syntax error at or near "normalized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalized BY NONE
< error
| 42601 syntax error at or near "normalized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalized NAME omni_probe_name
< error
| 42601 syntax error at or near "normalized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq normalized CYCLE
< error
| 42601 syntax error at or near "normalized"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq not
< error
| 42601 syntax error at or near "not"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq not 1
< error
| 42601 syntax error at or near "not"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq not int
< error
| 42601 syntax error at or near "not"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq not BY NONE
< error
| 42601 syntax error at or near "not"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq not NAME omni_probe_name
< error
| 42601 syntax error at or near "not"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq not CYCLE
< error
| 42601 syntax error at or near "not"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nothing
< error
| 42601 syntax error at or near "nothing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nothing 1
< error
| 42601 syntax error at or near "nothing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nothing int
< error
| 42601 syntax error at or near "nothing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nothing BY NONE
< error
| 42601 syntax error at or near "nothing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nothing NAME omni_probe_name
< error
| 42601 syntax error at or near "nothing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nothing CYCLE
< error
| 42601 syntax error at or near "nothing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notify
< error
| 42601 syntax error at or near "notify"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notify 1
< error
| 42601 syntax error at or near "notify"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notify int
< error
| 42601 syntax error at or near "notify"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notify BY NONE
< error
| 42601 syntax error at or near "notify"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notify NAME omni_probe_name
< error
| 42601 syntax error at or near "notify"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notify CYCLE
< error
| 42601 syntax error at or near "notify"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notnull
< error
| 42601 syntax error at or near "notnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notnull 1
< error
| 42601 syntax error at or near "notnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notnull int
< error
| 42601 syntax error at or near "notnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notnull BY NONE
< error
| 42601 syntax error at or near "notnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notnull NAME omni_probe_name
< error
| 42601 syntax error at or near "notnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq notnull CYCLE
< error
| 42601 syntax error at or near "notnull"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nowait
< error
| 42601 syntax error at or near "nowait"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nowait 1
< error
| 42601 syntax error at or near "nowait"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nowait int
< error
| 42601 syntax error at or near "nowait"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nowait BY NONE
< error
| 42601 syntax error at or near "nowait"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nowait NAME omni_probe_name
< error
| 42601 syntax error at or near "nowait"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nowait CYCLE
< error
| 42601 syntax error at or near "nowait"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq null
< error
| 42601 syntax error at or near "null"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq null 1
< error
| 42601 syntax error at or near "null"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq null int
< error
| 42601 syntax error at or near "null"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq null BY NONE
< error
| 42601 syntax error at or near "null"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq null NAME omni_probe_name
< error
| 42601 syntax error at or near "null"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq null CYCLE
< error
| 42601 syntax error at or near "null"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nullif
< error
| 42601 syntax error at or near "nullif"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nullif 1
< error
| 42601 syntax error at or near "nullif"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nullif int
< error
| 42601 syntax error at or near "nullif"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nullif BY NONE
< error
| 42601 syntax error at or near "nullif"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nullif NAME omni_probe_name
< error
| 42601 syntax error at or near "nullif"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nullif CYCLE
< error
| 42601 syntax error at or near "nullif"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nulls
< error
| 42601 syntax error at or near "nulls"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nulls 1
< error
| 42601 syntax error at or near "nulls"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nulls int
< error
| 42601 syntax error at or near "nulls"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nulls BY NONE
< error
| 42601 syntax error at or near "nulls"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nulls NAME omni_probe_name
< error
| 42601 syntax error at or near "nulls"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq nulls CYCLE
< error
| 42601 syntax error at or near "nulls"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq numeric
< error
| 42601 syntax error at or near "numeric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq numeric 1
< error
| 42601 syntax error at or near "numeric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq numeric int
< error
| 42601 syntax error at or near "numeric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq numeric BY NONE
< error
| 42601 syntax error at or near "numeric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq numeric NAME omni_probe_name
< error
| 42601 syntax error at or near "numeric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq numeric CYCLE
< error
| 42601 syntax error at or near "numeric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq object
< error
| 42601 syntax error at or near "object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq object 1
< error
| 42601 syntax error at or near "object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq object int
< error
| 42601 syntax error at or near "object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq object BY NONE
< error
| 42601 syntax error at or near "object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq object NAME omni_probe_name
< error
| 42601 syntax error at or near "object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq object CYCLE
< error
| 42601 syntax error at or near "object"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq of
< error
| 42601 syntax error at or near "of"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq of 1
< error
| 42601 syntax error at or near "of"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq of int
< error
| 42601 syntax error at or near "of"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq of BY NONE
< error
| 42601 syntax error at or near "of"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq of NAME omni_probe_name
< error
| 42601 syntax error at or near "of"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq of CYCLE
< error
| 42601 syntax error at or near "of"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq off
< error
| 42601 syntax error at or near "off"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq off 1
< error
| 42601 syntax error at or near "off"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq off int
< error
| 42601 syntax error at or near "off"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq off BY NONE
< error
| 42601 syntax error at or near "off"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq off NAME omni_probe_name
< error
| 42601 syntax error at or near "off"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq off CYCLE
< error
| 42601 syntax error at or near "off"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq offset
< error
| 42601 syntax error at or near "offset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq offset 1
< error
| 42601 syntax error at or near "offset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq offset int
< error
| 42601 syntax error at or near "offset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq offset BY NONE
< error
| 42601 syntax error at or near "offset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq offset NAME omni_probe_name
< error
| 42601 syntax error at or near "offset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq offset CYCLE
< error
| 42601 syntax error at or near "offset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq oids
< error
| 42601 syntax error at or near "oids"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq oids 1
< error
| 42601 syntax error at or near "oids"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq oids int
< error
| 42601 syntax error at or near "oids"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq oids BY NONE
< error
| 42601 syntax error at or near "oids"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq oids NAME omni_probe_name
< error
| 42601 syntax error at or near "oids"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq oids CYCLE
< error
| 42601 syntax error at or near "oids"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq old
< error
| 42601 syntax error at or near "old"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq old 1
< error
| 42601 syntax error at or near "old"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq old int
< error
| 42601 syntax error at or near "old"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq old BY NONE
< error
| 42601 syntax error at or near "old"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq old NAME omni_probe_name
< error
| 42601 syntax error at or near "old"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq old CYCLE
< error
| 42601 syntax error at or near "old"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq omit
< error
| 42601 syntax error at or near "omit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq omit 1
< error
| 42601 syntax error at or near "omit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq omit int
< error
| 42601 syntax error at or near "omit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq omit BY NONE
< error
| 42601 syntax error at or near "omit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq omit NAME omni_probe_name
< error
| 42601 syntax error at or near "omit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq omit CYCLE
< error
| 42601 syntax error at or near "omit"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq on
< error
| 42601 syntax error at or near "on"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq on 1
< error
| 42601 syntax error at or near "on"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq on int
< error
| 42601 syntax error at or near "on"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq on BY NONE
< error
| 42601 syntax error at or near "on"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq on NAME omni_probe_name
< error
| 42601 syntax error at or near "on"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq on CYCLE
< error
| 42601 syntax error at or near "on"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq only
< error
| 42601 syntax error at or near "only"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq only 1
< error
| 42601 syntax error at or near "only"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq only int
< error
| 42601 syntax error at or near "only"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq only BY NONE
< error
| 42601 syntax error at or near "only"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq only NAME omni_probe_name
< error
| 42601 syntax error at or near "only"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq only CYCLE
< error
| 42601 syntax error at or near "only"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq operator
< error
| 42601 syntax error at or near "operator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq operator 1
< error
| 42601 syntax error at or near "operator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq operator int
< error
| 42601 syntax error at or near "operator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq operator BY NONE
< error
| 42601 syntax error at or near "operator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq operator NAME omni_probe_name
< error
| 42601 syntax error at or near "operator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq operator CYCLE
< error
| 42601 syntax error at or near "operator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq option
< error
| 42601 syntax error at or near "option"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq option 1
< error
| 42601 syntax error at or near "option"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq option int
< error
| 42601 syntax error at or near "option"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq option BY NONE
< error
| 42601 syntax error at or near "option"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq option NAME omni_probe_name
< error
| 42601 syntax error at or near "option"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq option CYCLE
< error
| 42601 syntax error at or near "option"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq options
< error
| 42601 syntax error at or near "options"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq options 1
< error
| 42601 syntax error at or near "options"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq options int
< error
| 42601 syntax error at or near "options"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq options BY NONE
< error
| 42601 syntax error at or near "options"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq options NAME omni_probe_name
< error
| 42601 syntax error at or near "options"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq options CYCLE
< error
| 42601 syntax error at or near "options"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq or
< error
| 42601 syntax error at or near "or"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq or 1
< error
| 42601 syntax error at or near "or"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq or int
< error
| 42601 syntax error at or near "or"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq or BY NONE
< error
| 42601 syntax error at or near "or"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq or NAME omni_probe_name
< error
| 42601 syntax error at or near "or"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq or CYCLE
< error
| 42601 syntax error at or near "or"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq order
< error
| 42601 syntax error at or near "order"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq order 1
< error
| 42601 syntax error at or near "order"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq order int
< error
| 42601 syntax error at or near "order"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq order BY NONE
< error
| 42601 syntax error at or near "order"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq order NAME omni_probe_name
< error
| 42601 syntax error at or near "order"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq order CYCLE
< error
| 42601 syntax error at or near "order"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ordinality
< error
| 42601 syntax error at or near "ordinality"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ordinality 1
< error
| 42601 syntax error at or near "ordinality"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ordinality int
< error
| 42601 syntax error at or near "ordinality"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ordinality BY NONE
< error
| 42601 syntax error at or near "ordinality"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ordinality NAME omni_probe_name
< error
| 42601 syntax error at or near "ordinality"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ordinality CYCLE
< error
| 42601 syntax error at or near "ordinality"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq others
< error
| 42601 syntax error at or near "others"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq others 1
< error
| 42601 syntax error at or near "others"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq others int
< error
| 42601 syntax error at or near "others"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq others BY NONE
< error
| 42601 syntax error at or near "others"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq others NAME omni_probe_name
< error
| 42601 syntax error at or near "others"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq others CYCLE
< error
| 42601 syntax error at or near "others"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq out
< error
| 42601 syntax error at or near "out"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq out 1
< error
| 42601 syntax error at or near "out"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq out int
< error
| 42601 syntax error at or near "out"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq out BY NONE
< error
| 42601 syntax error at or near "out"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq out NAME omni_probe_name
< error
| 42601 syntax error at or near "out"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq out CYCLE
< error
| 42601 syntax error at or near "out"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq outer
< error
| 42601 syntax error at or near "outer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq outer 1
< error
| 42601 syntax error at or near "outer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq outer int
< error
| 42601 syntax error at or near "outer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq outer BY NONE
< error
| 42601 syntax error at or near "outer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq outer NAME omni_probe_name
< error
| 42601 syntax error at or near "outer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq outer CYCLE
< error
| 42601 syntax error at or near "outer"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq over
< error
| 42601 syntax error at or near "over"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq over 1
< error
| 42601 syntax error at or near "over"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq over int
< error
| 42601 syntax error at or near "over"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq over BY NONE
< error
| 42601 syntax error at or near "over"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq over NAME omni_probe_name
< error
| 42601 syntax error at or near "over"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq over CYCLE
< error
| 42601 syntax error at or near "over"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlaps
< error
| 42601 syntax error at or near "overlaps"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlaps 1
< error
| 42601 syntax error at or near "overlaps"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlaps int
< error
| 42601 syntax error at or near "overlaps"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlaps BY NONE
< error
| 42601 syntax error at or near "overlaps"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlaps NAME omni_probe_name
< error
| 42601 syntax error at or near "overlaps"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlaps CYCLE
< error
| 42601 syntax error at or near "overlaps"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlay
< error
| 42601 syntax error at or near "overlay"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlay 1
< error
| 42601 syntax error at or near "overlay"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlay int
< error
| 42601 syntax error at or near "overlay"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlay BY NONE
< error
| 42601 syntax error at or near "overlay"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlay NAME omni_probe_name
< error
| 42601 syntax error at or near "overlay"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overlay CYCLE
< error
| 42601 syntax error at or near "overlay"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overriding
< error
| 42601 syntax error at or near "overriding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overriding 1
< error
| 42601 syntax error at or near "overriding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overriding int
< error
| 42601 syntax error at or near "overriding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overriding BY NONE
< error
| 42601 syntax error at or near "overriding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overriding NAME omni_probe_name
< error
| 42601 syntax error at or near "overriding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq overriding CYCLE
< error
| 42601 syntax error at or near "overriding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owned
< error
| 42601 syntax error at end of input
| position: 37
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owned 1
< error
| 42601 syntax error at or near "1"
| position: 38
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owned int
< error
| 42601 syntax error at or near "int"
| position: 38
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owned BY NONE
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owner
< error
| 42601 syntax error at or near "owner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owner 1
< error
| 42601 syntax error at or near "owner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owner int
< error
| 42601 syntax error at or near "owner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owner BY NONE
< error
| 42601 syntax error at or near "owner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owner NAME omni_probe_name
< error
| 42601 syntax error at or near "owner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq owner CYCLE
< error
| 42601 syntax error at or near "owner"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parallel
< error
| 42601 syntax error at or near "parallel"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parallel 1
< error
| 42601 syntax error at or near "parallel"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parallel int
< error
| 42601 syntax error at or near "parallel"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parallel BY NONE
< error
| 42601 syntax error at or near "parallel"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parallel NAME omni_probe_name
< error
| 42601 syntax error at or near "parallel"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parallel CYCLE
< error
| 42601 syntax error at or near "parallel"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parameter
< error
| 42601 syntax error at or near "parameter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parameter 1
< error
| 42601 syntax error at or near "parameter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parameter int
< error
| 42601 syntax error at or near "parameter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parameter BY NONE
< error
| 42601 syntax error at or near "parameter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parameter NAME omni_probe_name
< error
| 42601 syntax error at or near "parameter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parameter CYCLE
< error
| 42601 syntax error at or near "parameter"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parser
< error
| 42601 syntax error at or near "parser"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parser 1
< error
| 42601 syntax error at or near "parser"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parser int
< error
| 42601 syntax error at or near "parser"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parser BY NONE
< error
| 42601 syntax error at or near "parser"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parser NAME omni_probe_name
< error
| 42601 syntax error at or near "parser"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq parser CYCLE
< error
| 42601 syntax error at or near "parser"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partial
< error
| 42601 syntax error at or near "partial"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partial 1
< error
| 42601 syntax error at or near "partial"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partial int
< error
| 42601 syntax error at or near "partial"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partial BY NONE
< error
| 42601 syntax error at or near "partial"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partial NAME omni_probe_name
< error
| 42601 syntax error at or near "partial"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partial CYCLE
< error
| 42601 syntax error at or near "partial"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partition
< error
| 42601 syntax error at or near "partition"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partition 1
< error
| 42601 syntax error at or near "partition"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partition int
< error
| 42601 syntax error at or near "partition"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partition BY NONE
< error
| 42601 syntax error at or near "partition"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partition NAME omni_probe_name
< error
| 42601 syntax error at or near "partition"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq partition CYCLE
< error
| 42601 syntax error at or near "partition"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq passing
< error
| 42601 syntax error at or near "passing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq passing 1
< error
| 42601 syntax error at or near "passing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq passing int
< error
| 42601 syntax error at or near "passing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq passing BY NONE
< error
| 42601 syntax error at or near "passing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq passing NAME omni_probe_name
< error
| 42601 syntax error at or near "passing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq passing CYCLE
< error
| 42601 syntax error at or near "passing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq password
< error
| 42601 syntax error at or near "password"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq password 1
< error
| 42601 syntax error at or near "password"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq password int
< error
| 42601 syntax error at or near "password"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq password BY NONE
< error
| 42601 syntax error at or near "password"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq password NAME omni_probe_name
< error
| 42601 syntax error at or near "password"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq password CYCLE
< error
| 42601 syntax error at or near "password"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq path
< error
| 42601 syntax error at or near "path"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq path 1
< error
| 42601 syntax error at or near "path"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq path int
< error
| 42601 syntax error at or near "path"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq path BY NONE
< error
| 42601 syntax error at or near "path"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq path NAME omni_probe_name
< error
| 42601 syntax error at or near "path"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq path CYCLE
< error
| 42601 syntax error at or near "path"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq placing
< error
| 42601 syntax error at or near "placing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq placing 1
< error
| 42601 syntax error at or near "placing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq placing int
< error
| 42601 syntax error at or near "placing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq placing BY NONE
< error
| 42601 syntax error at or near "placing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq placing NAME omni_probe_name
< error
| 42601 syntax error at or near "placing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq placing CYCLE
< error
| 42601 syntax error at or near "placing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plan
< error
| 42601 syntax error at or near "plan"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plan 1
< error
| 42601 syntax error at or near "plan"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plan int
< error
| 42601 syntax error at or near "plan"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plan BY NONE
< error
| 42601 syntax error at or near "plan"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plan NAME omni_probe_name
< error
| 42601 syntax error at or near "plan"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plan CYCLE
< error
| 42601 syntax error at or near "plan"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plans
< error
| 42601 syntax error at or near "plans"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plans 1
< error
| 42601 syntax error at or near "plans"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plans int
< error
| 42601 syntax error at or near "plans"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plans BY NONE
< error
| 42601 syntax error at or near "plans"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plans NAME omni_probe_name
< error
| 42601 syntax error at or near "plans"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq plans CYCLE
< error
| 42601 syntax error at or near "plans"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq policy
< error
| 42601 syntax error at or near "policy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq policy 1
< error
| 42601 syntax error at or near "policy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq policy int
< error
| 42601 syntax error at or near "policy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq policy BY NONE
< error
| 42601 syntax error at or near "policy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq policy NAME omni_probe_name
< error
| 42601 syntax error at or near "policy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq policy CYCLE
< error
| 42601 syntax error at or near "policy"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq position
< error
| 42601 syntax error at or near "position"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq position 1
< error
| 42601 syntax error at or near "position"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq position int
< error
| 42601 syntax error at or near "position"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq position BY NONE
< error
| 42601 syntax error at or near "position"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq position NAME omni_probe_name
< error
| 42601 syntax error at or near "position"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq position CYCLE
< error
| 42601 syntax error at or near "position"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preceding
< error
| 42601 syntax error at or near "preceding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preceding 1
< error
| 42601 syntax error at or near "preceding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preceding int
< error
| 42601 syntax error at or near "preceding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preceding BY NONE
< error
| 42601 syntax error at or near "preceding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preceding NAME omni_probe_name
< error
| 42601 syntax error at or near "preceding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preceding CYCLE
< error
| 42601 syntax error at or near "preceding"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq precision
< error
| 42601 syntax error at or near "precision"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq precision 1
< error
| 42601 syntax error at or near "precision"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq precision int
< error
| 42601 syntax error at or near "precision"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq precision BY NONE
< error
| 42601 syntax error at or near "precision"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq precision NAME omni_probe_name
< error
| 42601 syntax error at or near "precision"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq precision CYCLE
< error
| 42601 syntax error at or near "precision"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepare
< error
| 42601 syntax error at or near "prepare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepare 1
< error
| 42601 syntax error at or near "prepare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepare int
< error
| 42601 syntax error at or near "prepare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepare BY NONE
< error
| 42601 syntax error at or near "prepare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepare NAME omni_probe_name
< error
| 42601 syntax error at or near "prepare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepare CYCLE
< error
| 42601 syntax error at or near "prepare"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepared
< error
| 42601 syntax error at or near "prepared"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepared 1
< error
| 42601 syntax error at or near "prepared"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepared int
< error
| 42601 syntax error at or near "prepared"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepared BY NONE
< error
| 42601 syntax error at or near "prepared"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepared NAME omni_probe_name
< error
| 42601 syntax error at or near "prepared"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prepared CYCLE
< error
| 42601 syntax error at or near "prepared"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preserve
< error
| 42601 syntax error at or near "preserve"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preserve 1
< error
| 42601 syntax error at or near "preserve"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preserve int
< error
| 42601 syntax error at or near "preserve"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preserve BY NONE
< error
| 42601 syntax error at or near "preserve"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preserve NAME omni_probe_name
< error
| 42601 syntax error at or near "preserve"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq preserve CYCLE
< error
| 42601 syntax error at or near "preserve"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq primary
< error
| 42601 syntax error at or near "primary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq primary 1
< error
| 42601 syntax error at or near "primary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq primary int
< error
| 42601 syntax error at or near "primary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq primary BY NONE
< error
| 42601 syntax error at or near "primary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq primary NAME omni_probe_name
< error
| 42601 syntax error at or near "primary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq primary CYCLE
< error
| 42601 syntax error at or near "primary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prior
< error
| 42601 syntax error at or near "prior"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prior 1
< error
| 42601 syntax error at or near "prior"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prior int
< error
| 42601 syntax error at or near "prior"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prior BY NONE
< error
| 42601 syntax error at or near "prior"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prior NAME omni_probe_name
< error
| 42601 syntax error at or near "prior"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq prior CYCLE
< error
| 42601 syntax error at or near "prior"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq privileges
< error
| 42601 syntax error at or near "privileges"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq privileges 1
< error
| 42601 syntax error at or near "privileges"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq privileges int
< error
| 42601 syntax error at or near "privileges"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq privileges BY NONE
< error
| 42601 syntax error at or near "privileges"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq privileges NAME omni_probe_name
< error
| 42601 syntax error at or near "privileges"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq privileges CYCLE
< error
| 42601 syntax error at or near "privileges"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedural
< error
| 42601 syntax error at or near "procedural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedural 1
< error
| 42601 syntax error at or near "procedural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedural int
< error
| 42601 syntax error at or near "procedural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedural BY NONE
< error
| 42601 syntax error at or near "procedural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedural NAME omni_probe_name
< error
| 42601 syntax error at or near "procedural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedural CYCLE
< error
| 42601 syntax error at or near "procedural"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedure
< error
| 42601 syntax error at or near "procedure"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedure 1
< error
| 42601 syntax error at or near "procedure"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedure int
< error
| 42601 syntax error at or near "procedure"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedure BY NONE
< error
| 42601 syntax error at or near "procedure"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedure NAME omni_probe_name
< error
| 42601 syntax error at or near "procedure"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedure CYCLE
< error
| 42601 syntax error at or near "procedure"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedures
< error
| 42601 syntax error at or near "procedures"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedures 1
< error
| 42601 syntax error at or near "procedures"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedures int
< error
| 42601 syntax error at or near "procedures"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedures BY NONE
< error
| 42601 syntax error at or near "procedures"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedures NAME omni_probe_name
< error
| 42601 syntax error at or near "procedures"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq procedures CYCLE
< error
| 42601 syntax error at or near "procedures"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq program
< error
| 42601 syntax error at or near "program"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq program 1
< error
| 42601 syntax error at or near "program"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq program int
< error
| 42601 syntax error at or near "program"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq program BY NONE
< error
| 42601 syntax error at or near "program"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq program NAME omni_probe_name
< error
| 42601 syntax error at or near "program"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq program CYCLE
< error
| 42601 syntax error at or near "program"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq publication
< error
| 42601 syntax error at or near "publication"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq publication 1
< error
| 42601 syntax error at or near "publication"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq publication int
< error
| 42601 syntax error at or near "publication"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq publication BY NONE
< error
| 42601 syntax error at or near "publication"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq publication NAME omni_probe_name
< error
| 42601 syntax error at or near "publication"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq publication CYCLE
< error
| 42601 syntax error at or near "publication"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quote
< error
| 42601 syntax error at or near "quote"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quote 1
< error
| 42601 syntax error at or near "quote"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quote int
< error
| 42601 syntax error at or near "quote"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quote BY NONE
< error
| 42601 syntax error at or near "quote"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quote NAME omni_probe_name
< error
| 42601 syntax error at or near "quote"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quote CYCLE
< error
| 42601 syntax error at or near "quote"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quotes
< error
| 42601 syntax error at or near "quotes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quotes 1
< error
| 42601 syntax error at or near "quotes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quotes int
< error
| 42601 syntax error at or near "quotes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quotes BY NONE
< error
| 42601 syntax error at or near "quotes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quotes NAME omni_probe_name
< error
| 42601 syntax error at or near "quotes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq quotes CYCLE
< error
| 42601 syntax error at or near "quotes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq range
< error
| 42601 syntax error at or near "range"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq range 1
< error
| 42601 syntax error at or near "range"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq range int
< error
| 42601 syntax error at or near "range"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq range BY NONE
< error
| 42601 syntax error at or near "range"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq range NAME omni_probe_name
< error
| 42601 syntax error at or near "range"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq range CYCLE
< error
| 42601 syntax error at or near "range"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq read
< error
| 42601 syntax error at or near "read"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq read 1
< error
| 42601 syntax error at or near "read"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq read int
< error
| 42601 syntax error at or near "read"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq read BY NONE
< error
| 42601 syntax error at or near "read"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq read NAME omni_probe_name
< error
| 42601 syntax error at or near "read"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq read CYCLE
< error
| 42601 syntax error at or near "read"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq real
< error
| 42601 syntax error at or near "real"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq real 1
< error
| 42601 syntax error at or near "real"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq real int
< error
| 42601 syntax error at or near "real"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq real BY NONE
< error
| 42601 syntax error at or near "real"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq real NAME omni_probe_name
< error
| 42601 syntax error at or near "real"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq real CYCLE
< error
| 42601 syntax error at or near "real"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reassign
< error
| 42601 syntax error at or near "reassign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reassign 1
< error
| 42601 syntax error at or near "reassign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reassign int
< error
| 42601 syntax error at or near "reassign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reassign BY NONE
< error
| 42601 syntax error at or near "reassign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reassign NAME omni_probe_name
< error
| 42601 syntax error at or near "reassign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reassign CYCLE
< error
| 42601 syntax error at or near "reassign"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recheck
< error
| 42601 syntax error at or near "recheck"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recheck 1
< error
| 42601 syntax error at or near "recheck"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recheck int
< error
| 42601 syntax error at or near "recheck"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recheck BY NONE
< error
| 42601 syntax error at or near "recheck"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recheck NAME omni_probe_name
< error
| 42601 syntax error at or near "recheck"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recheck CYCLE
< error
| 42601 syntax error at or near "recheck"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recursive
< error
| 42601 syntax error at or near "recursive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recursive 1
< error
| 42601 syntax error at or near "recursive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recursive int
< error
| 42601 syntax error at or near "recursive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recursive BY NONE
< error
| 42601 syntax error at or near "recursive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recursive NAME omni_probe_name
< error
| 42601 syntax error at or near "recursive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq recursive CYCLE
< error
| 42601 syntax error at or near "recursive"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ref
< error
| 42601 syntax error at or near "ref"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ref 1
< error
| 42601 syntax error at or near "ref"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ref int
< error
| 42601 syntax error at or near "ref"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ref BY NONE
< error
| 42601 syntax error at or near "ref"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ref NAME omni_probe_name
< error
| 42601 syntax error at or near "ref"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ref CYCLE
< error
| 42601 syntax error at or near "ref"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq references
< error
| 42601 syntax error at or near "references"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq references 1
< error
| 42601 syntax error at or near "references"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq references int
< error
| 42601 syntax error at or near "references"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq references BY NONE
< error
| 42601 syntax error at or near "references"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq references NAME omni_probe_name
< error
| 42601 syntax error at or near "references"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq references CYCLE
< error
| 42601 syntax error at or near "references"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq referencing
< error
| 42601 syntax error at or near "referencing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq referencing 1
< error
| 42601 syntax error at or near "referencing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq referencing int
< error
| 42601 syntax error at or near "referencing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq referencing BY NONE
< error
| 42601 syntax error at or near "referencing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq referencing NAME omni_probe_name
< error
| 42601 syntax error at or near "referencing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq referencing CYCLE
< error
| 42601 syntax error at or near "referencing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq refresh
< error
| 42601 syntax error at or near "refresh"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq refresh 1
< error
| 42601 syntax error at or near "refresh"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq refresh int
< error
| 42601 syntax error at or near "refresh"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq refresh BY NONE
< error
| 42601 syntax error at or near "refresh"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq refresh NAME omni_probe_name
< error
| 42601 syntax error at or near "refresh"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq refresh CYCLE
< error
| 42601 syntax error at or near "refresh"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reindex
< error
| 42601 syntax error at or near "reindex"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reindex 1
< error
| 42601 syntax error at or near "reindex"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reindex int
< error
| 42601 syntax error at or near "reindex"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reindex BY NONE
< error
| 42601 syntax error at or near "reindex"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reindex NAME omni_probe_name
< error
| 42601 syntax error at or near "reindex"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reindex CYCLE
< error
| 42601 syntax error at or near "reindex"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq relative
< error
| 42601 syntax error at or near "relative"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq relative 1
< error
| 42601 syntax error at or near "relative"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq relative int
< error
| 42601 syntax error at or near "relative"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq relative BY NONE
< error
| 42601 syntax error at or near "relative"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq relative NAME omni_probe_name
< error
| 42601 syntax error at or near "relative"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq relative CYCLE
< error
| 42601 syntax error at or near "relative"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq release
< error
| 42601 syntax error at or near "release"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq release 1
< error
| 42601 syntax error at or near "release"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq release int
< error
| 42601 syntax error at or near "release"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq release BY NONE
< error
| 42601 syntax error at or near "release"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq release NAME omni_probe_name
< error
| 42601 syntax error at or near "release"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq release CYCLE
< error
| 42601 syntax error at or near "release"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rename
< error
| 42601 syntax error at or near "rename"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rename 1
< error
| 42601 syntax error at or near "rename"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rename int
< error
| 42601 syntax error at or near "rename"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rename BY NONE
< error
| 42601 syntax error at or near "rename"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rename NAME omni_probe_name
< error
| 42601 syntax error at or near "rename"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rename CYCLE
< error
| 42601 syntax error at or near "rename"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq repeatable
< error
| 42601 syntax error at or near "repeatable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq repeatable 1
< error
| 42601 syntax error at or near "repeatable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq repeatable int
< error
| 42601 syntax error at or near "repeatable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq repeatable BY NONE
< error
| 42601 syntax error at or near "repeatable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq repeatable NAME omni_probe_name
< error
| 42601 syntax error at or near "repeatable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq repeatable CYCLE
< error
| 42601 syntax error at or near "repeatable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replace
< error
| 42601 syntax error at or near "replace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replace 1
< error
| 42601 syntax error at or near "replace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replace int
< error
| 42601 syntax error at or near "replace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replace BY NONE
< error
| 42601 syntax error at or near "replace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replace NAME omni_probe_name
< error
| 42601 syntax error at or near "replace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replace CYCLE
< error
| 42601 syntax error at or near "replace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replica
< error
| 42601 syntax error at or near "replica"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replica 1
< error
| 42601 syntax error at or near "replica"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replica int
< error
| 42601 syntax error at or near "replica"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replica BY NONE
< error
| 42601 syntax error at or near "replica"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replica NAME omni_probe_name
< error
| 42601 syntax error at or near "replica"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq replica CYCLE
< error
| 42601 syntax error at or near "replica"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reset
< error
| 42601 syntax error at or near "reset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reset 1
< error
| 42601 syntax error at or near "reset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reset int
< error
| 42601 syntax error at or near "reset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reset BY NONE
< error
| 42601 syntax error at or near "reset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reset NAME omni_probe_name
< error
| 42601 syntax error at or near "reset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq reset CYCLE
< error
| 42601 syntax error at or near "reset"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq restart
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq restrict
< error
| 42601 syntax error at or near "restrict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq restrict 1
< error
| 42601 syntax error at or near "restrict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq restrict int
< error
| 42601 syntax error at or near "restrict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq restrict BY NONE
< error
| 42601 syntax error at or near "restrict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq restrict NAME omni_probe_name
< error
| 42601 syntax error at or near "restrict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq restrict CYCLE
< error
| 42601 syntax error at or near "restrict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq return
< error
| 42601 syntax error at or near "return"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq return 1
< error
| 42601 syntax error at or near "return"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq return int
< error
| 42601 syntax error at or near "return"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq return BY NONE
< error
| 42601 syntax error at or near "return"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq return NAME omni_probe_name
< error
| 42601 syntax error at or near "return"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq return CYCLE
< error
| 42601 syntax error at or near "return"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returning
< error
| 42601 syntax error at or near "returning"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returning 1
< error
| 42601 syntax error at or near "returning"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returning int
< error
| 42601 syntax error at or near "returning"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returning BY NONE
< error
| 42601 syntax error at or near "returning"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returning NAME omni_probe_name
< error
| 42601 syntax error at or near "returning"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returning CYCLE
< error
| 42601 syntax error at or near "returning"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returns
< error
| 42601 syntax error at or near "returns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returns 1
< error
| 42601 syntax error at or near "returns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returns int
< error
| 42601 syntax error at or near "returns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returns BY NONE
< error
| 42601 syntax error at or near "returns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returns NAME omni_probe_name
< error
| 42601 syntax error at or near "returns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq returns CYCLE
< error
| 42601 syntax error at or near "returns"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq revoke
< error
| 42601 syntax error at or near "revoke"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq revoke 1
< error
| 42601 syntax error at or near "revoke"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq revoke int
< error
| 42601 syntax error at or near "revoke"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq revoke BY NONE
< error
| 42601 syntax error at or near "revoke"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq revoke NAME omni_probe_name
< error
| 42601 syntax error at or near "revoke"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq revoke CYCLE
< error
| 42601 syntax error at or near "revoke"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq right
< error
| 42601 syntax error at or near "right"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq right 1
< error
| 42601 syntax error at or near "right"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq right int
< error
| 42601 syntax error at or near "right"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq right BY NONE
< error
| 42601 syntax error at or near "right"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq right NAME omni_probe_name
< error
| 42601 syntax error at or near "right"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq right CYCLE
< error
| 42601 syntax error at or near "right"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq role
< error
| 42601 syntax error at or near "role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq role 1
< error
| 42601 syntax error at or near "role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq role int
< error
| 42601 syntax error at or near "role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq role BY NONE
< error
| 42601 syntax error at or near "role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq role NAME omni_probe_name
< error
| 42601 syntax error at or near "role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq role CYCLE
< error
| 42601 syntax error at or near "role"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollback
< error
| 42601 syntax error at or near "rollback"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollback 1
< error
| 42601 syntax error at or near "rollback"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollback int
< error
| 42601 syntax error at or near "rollback"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollback BY NONE
< error
| 42601 syntax error at or near "rollback"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollback NAME omni_probe_name
< error
| 42601 syntax error at or near "rollback"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollback CYCLE
< error
| 42601 syntax error at or near "rollback"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollup
< error
| 42601 syntax error at or near "rollup"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollup 1
< error
| 42601 syntax error at or near "rollup"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollup int
< error
| 42601 syntax error at or near "rollup"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollup BY NONE
< error
| 42601 syntax error at or near "rollup"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollup NAME omni_probe_name
< error
| 42601 syntax error at or near "rollup"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rollup CYCLE
< error
| 42601 syntax error at or near "rollup"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routine
< error
| 42601 syntax error at or near "routine"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routine 1
< error
| 42601 syntax error at or near "routine"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routine int
< error
| 42601 syntax error at or near "routine"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routine BY NONE
< error
| 42601 syntax error at or near "routine"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routine NAME omni_probe_name
< error
| 42601 syntax error at or near "routine"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routine CYCLE
< error
| 42601 syntax error at or near "routine"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routines
< error
| 42601 syntax error at or near "routines"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routines 1
< error
| 42601 syntax error at or near "routines"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routines int
< error
| 42601 syntax error at or near "routines"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routines BY NONE
< error
| 42601 syntax error at or near "routines"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routines NAME omni_probe_name
< error
| 42601 syntax error at or near "routines"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq routines CYCLE
< error
| 42601 syntax error at or near "routines"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq row
< error
| 42601 syntax error at or near "row"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq row 1
< error
| 42601 syntax error at or near "row"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq row int
< error
| 42601 syntax error at or near "row"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq row BY NONE
< error
| 42601 syntax error at or near "row"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq row NAME omni_probe_name
< error
| 42601 syntax error at or near "row"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq row CYCLE
< error
| 42601 syntax error at or near "row"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rows
< error
| 42601 syntax error at or near "rows"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rows 1
< error
| 42601 syntax error at or near "rows"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rows int
< error
| 42601 syntax error at or near "rows"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rows BY NONE
< error
| 42601 syntax error at or near "rows"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rows NAME omni_probe_name
< error
| 42601 syntax error at or near "rows"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rows CYCLE
< error
| 42601 syntax error at or near "rows"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rule
< error
| 42601 syntax error at or near "rule"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rule 1
< error
| 42601 syntax error at or near "rule"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rule int
< error
| 42601 syntax error at or near "rule"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rule BY NONE
< error
| 42601 syntax error at or near "rule"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rule NAME omni_probe_name
< error
| 42601 syntax error at or near "rule"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq rule CYCLE
< error
| 42601 syntax error at or near "rule"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq savepoint
< error
| 42601 syntax error at or near "savepoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq savepoint 1
< error
| 42601 syntax error at or near "savepoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq savepoint int
< error
| 42601 syntax error at or near "savepoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq savepoint BY NONE
< error
| 42601 syntax error at or near "savepoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq savepoint NAME omni_probe_name
< error
| 42601 syntax error at or near "savepoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq savepoint CYCLE
< error
| 42601 syntax error at or near "savepoint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scalar
< error
| 42601 syntax error at or near "scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scalar 1
< error
| 42601 syntax error at or near "scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scalar int
< error
| 42601 syntax error at or near "scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scalar BY NONE
< error
| 42601 syntax error at or near "scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scalar NAME omni_probe_name
< error
| 42601 syntax error at or near "scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scalar CYCLE
< error
| 42601 syntax error at or near "scalar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schema
< error
| 42601 syntax error at or near "schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schema 1
< error
| 42601 syntax error at or near "schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schema int
< error
| 42601 syntax error at or near "schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schema BY NONE
< error
| 42601 syntax error at or near "schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schema NAME omni_probe_name
< error
| 42601 syntax error at or near "schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schema CYCLE
< error
| 42601 syntax error at or near "schema"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schemas
< error
| 42601 syntax error at or near "schemas"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schemas 1
< error
| 42601 syntax error at or near "schemas"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schemas int
< error
| 42601 syntax error at or near "schemas"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schemas BY NONE
< error
| 42601 syntax error at or near "schemas"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schemas NAME omni_probe_name
< error
| 42601 syntax error at or near "schemas"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq schemas CYCLE
< error
| 42601 syntax error at or near "schemas"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scroll
< error
| 42601 syntax error at or near "scroll"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scroll 1
< error
| 42601 syntax error at or near "scroll"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scroll int
< error
| 42601 syntax error at or near "scroll"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scroll BY NONE
< error
| 42601 syntax error at or near "scroll"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scroll NAME omni_probe_name
< error
| 42601 syntax error at or near "scroll"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq scroll CYCLE
< error
| 42601 syntax error at or near "scroll"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq search
< error
| 42601 syntax error at or near "search"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq search 1
< error
| 42601 syntax error at or near "search"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq search int
< error
| 42601 syntax error at or near "search"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq search BY NONE
< error
| 42601 syntax error at or near "search"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq search NAME omni_probe_name
< error
| 42601 syntax error at or near "search"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq search CYCLE
< error
| 42601 syntax error at or near "search"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq second
< error
| 42601 syntax error at or near "second"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq second 1
< error
| 42601 syntax error at or near "second"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq second int
< error
| 42601 syntax error at or near "second"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq second BY NONE
< error
| 42601 syntax error at or near "second"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq second NAME omni_probe_name
< error
| 42601 syntax error at or near "second"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq second CYCLE
< error
| 42601 syntax error at or near "second"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq security
< error
| 42601 syntax error at or near "security"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq security 1
< error
| 42601 syntax error at or near "security"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq security int
< error
| 42601 syntax error at or near "security"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq security BY NONE
< error
| 42601 syntax error at or near "security"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq security NAME omni_probe_name
< error
| 42601 syntax error at or near "security"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq security CYCLE
< error
| 42601 syntax error at or near "security"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq select
< error
| 42601 syntax error at or near "select"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq select 1
< error
| 42601 syntax error at or near "select"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq select int
< error
| 42601 syntax error at or near "select"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq select BY NONE
< error
| 42601 syntax error at or near "select"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq select NAME omni_probe_name
< error
| 42601 syntax error at or near "select"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq select CYCLE
< error
| 42601 syntax error at or near "select"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequence
< error
| 42601 syntax error at end of input
| position: 40
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequence 1
< error
| 42601 syntax error at or near "1"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequence int
< error
| 42601 syntax error at or near "int"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequence BY NONE
< error
| 42601 syntax error at or near "BY"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequence NAME omni_probe_name
< error
| 42601 invalid sequence option SEQUENCE NAME
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequence CYCLE
< error
| 42601 syntax error at or near "CYCLE"
| position: 41
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequences
< error
| 42601 syntax error at or near "sequences"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequences 1
< error
| 42601 syntax error at or near "sequences"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequences int
< error
| 42601 syntax error at or near "sequences"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequences BY NONE
< error
| 42601 syntax error at or near "sequences"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequences NAME omni_probe_name
< error
| 42601 syntax error at or near "sequences"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sequences CYCLE
< error
| 42601 syntax error at or near "sequences"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq serializable
< error
| 42601 syntax error at or near "serializable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq serializable 1
< error
| 42601 syntax error at or near "serializable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq serializable int
< error
| 42601 syntax error at or near "serializable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq serializable BY NONE
< error
| 42601 syntax error at or near "serializable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq serializable NAME omni_probe_name
< error
| 42601 syntax error at or near "serializable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq serializable CYCLE
< error
| 42601 syntax error at or near "serializable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq server
< error
| 42601 syntax error at or near "server"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq server 1
< error
| 42601 syntax error at or near "server"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq server int
< error
| 42601 syntax error at or near "server"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq server BY NONE
< error
| 42601 syntax error at or near "server"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq server NAME omni_probe_name
< error
| 42601 syntax error at or near "server"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq server CYCLE
< error
| 42601 syntax error at or near "server"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session
< error
| 42601 syntax error at or near "session"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session 1
< error
| 42601 syntax error at or near "session"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session int
< error
| 42601 syntax error at or near "session"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session BY NONE
< error
| 42601 syntax error at or near "session"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session NAME omni_probe_name
< error
| 42601 syntax error at or near "session"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session CYCLE
< error
| 42601 syntax error at or near "session"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session_user
< error
| 42601 syntax error at or near "session_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session_user 1
< error
| 42601 syntax error at or near "session_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session_user int
< error
| 42601 syntax error at or near "session_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session_user BY NONE
< error
| 42601 syntax error at or near "session_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session_user NAME omni_probe_name
< error
| 42601 syntax error at or near "session_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq session_user CYCLE
< error
| 42601 syntax error at or near "session_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq set
< error
| 42601 syntax error at or near "set"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq set 1
< error
| 42601 syntax error at or near "set"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq set int
< error
| 42601 syntax error at or near "set"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq set BY NONE
< error
| 42601 syntax error at or near "set"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq set NAME omni_probe_name
< error
| 42601 syntax error at or near "set"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq set CYCLE
< error
| 42601 syntax error at or near "set"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq setof
< error
| 42601 syntax error at or near "setof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq setof 1
< error
| 42601 syntax error at or near "setof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq setof int
< error
| 42601 syntax error at or near "setof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq setof BY NONE
< error
| 42601 syntax error at or near "setof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq setof NAME omni_probe_name
< error
| 42601 syntax error at or near "setof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq setof CYCLE
< error
| 42601 syntax error at or near "setof"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sets
< error
| 42601 syntax error at or near "sets"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sets 1
< error
| 42601 syntax error at or near "sets"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sets int
< error
| 42601 syntax error at or near "sets"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sets BY NONE
< error
| 42601 syntax error at or near "sets"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sets NAME omni_probe_name
< error
| 42601 syntax error at or near "sets"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sets CYCLE
< error
| 42601 syntax error at or near "sets"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq share
< error
| 42601 syntax error at or near "share"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq share 1
< error
| 42601 syntax error at or near "share"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq share int
< error
| 42601 syntax error at or near "share"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq share BY NONE
< error
| 42601 syntax error at or near "share"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq share NAME omni_probe_name
< error
| 42601 syntax error at or near "share"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq share CYCLE
< error
| 42601 syntax error at or near "share"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq show
< error
| 42601 syntax error at or near "show"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq show 1
< error
| 42601 syntax error at or near "show"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq show int
< error
| 42601 syntax error at or near "show"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq show BY NONE
< error
| 42601 syntax error at or near "show"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq show NAME omni_probe_name
< error
| 42601 syntax error at or near "show"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq show CYCLE
< error
| 42601 syntax error at or near "show"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq similar
< error
| 42601 syntax error at or near "similar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq similar 1
< error
| 42601 syntax error at or near "similar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq similar int
< error
| 42601 syntax error at or near "similar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq similar BY NONE
< error
| 42601 syntax error at or near "similar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq similar NAME omni_probe_name
< error
| 42601 syntax error at or near "similar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq similar CYCLE
< error
| 42601 syntax error at or near "similar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq simple
< error
| 42601 syntax error at or near "simple"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq simple 1
< error
| 42601 syntax error at or near "simple"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq simple int
< error
| 42601 syntax error at or near "simple"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq simple BY NONE
< error
| 42601 syntax error at or near "simple"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq simple NAME omni_probe_name
< error
| 42601 syntax error at or near "simple"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq simple CYCLE
< error
| 42601 syntax error at or near "simple"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq skip
< error
| 42601 syntax error at or near "skip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq skip 1
< error
| 42601 syntax error at or near "skip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq skip int
< error
| 42601 syntax error at or near "skip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq skip BY NONE
< error
| 42601 syntax error at or near "skip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq skip NAME omni_probe_name
< error
| 42601 syntax error at or near "skip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq skip CYCLE
< error
| 42601 syntax error at or near "skip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq smallint
< error
| 42601 syntax error at or near "smallint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq smallint 1
< error
| 42601 syntax error at or near "smallint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq smallint int
< error
| 42601 syntax error at or near "smallint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq smallint BY NONE
< error
| 42601 syntax error at or near "smallint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq smallint NAME omni_probe_name
< error
| 42601 syntax error at or near "smallint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq smallint CYCLE
< error
| 42601 syntax error at or near "smallint"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq snapshot
< error
| 42601 syntax error at or near "snapshot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq snapshot 1
< error
| 42601 syntax error at or near "snapshot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq snapshot int
< error
| 42601 syntax error at or near "snapshot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq snapshot BY NONE
< error
| 42601 syntax error at or near "snapshot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq snapshot NAME omni_probe_name
< error
| 42601 syntax error at or near "snapshot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq snapshot CYCLE
< error
| 42601 syntax error at or near "snapshot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq some
< error
| 42601 syntax error at or near "some"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq some 1
< error
| 42601 syntax error at or near "some"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq some int
< error
| 42601 syntax error at or near "some"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq some BY NONE
< error
| 42601 syntax error at or near "some"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq some NAME omni_probe_name
< error
| 42601 syntax error at or near "some"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq some CYCLE
< error
| 42601 syntax error at or near "some"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq source
< error
| 42601 syntax error at or near "source"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq source 1
< error
| 42601 syntax error at or near "source"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq source int
< error
| 42601 syntax error at or near "source"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq source BY NONE
< error
| 42601 syntax error at or near "source"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq source NAME omni_probe_name
< error
| 42601 syntax error at or near "source"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq source CYCLE
< error
| 42601 syntax error at or near "source"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sql
< error
| 42601 syntax error at or near "sql"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sql 1
< error
| 42601 syntax error at or near "sql"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sql int
< error
| 42601 syntax error at or near "sql"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sql BY NONE
< error
| 42601 syntax error at or near "sql"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sql NAME omni_probe_name
< error
| 42601 syntax error at or near "sql"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sql CYCLE
< error
| 42601 syntax error at or near "sql"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stable
< error
| 42601 syntax error at or near "stable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stable 1
< error
| 42601 syntax error at or near "stable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stable int
< error
| 42601 syntax error at or near "stable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stable BY NONE
< error
| 42601 syntax error at or near "stable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stable NAME omni_probe_name
< error
| 42601 syntax error at or near "stable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stable CYCLE
< error
| 42601 syntax error at or near "stable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq standalone
< error
| 42601 syntax error at or near "standalone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq standalone 1
< error
| 42601 syntax error at or near "standalone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq standalone int
< error
| 42601 syntax error at or near "standalone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq standalone BY NONE
< error
| 42601 syntax error at or near "standalone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq standalone NAME omni_probe_name
< error
| 42601 syntax error at or near "standalone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq standalone CYCLE
< error
| 42601 syntax error at or near "standalone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq start
< error
| 42601 syntax error at end of input
| position: 37
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq start 1
< ok 0 -
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statement
< error
| 42601 syntax error at or near "statement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statement 1
< error
| 42601 syntax error at or near "statement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statement int
< error
| 42601 syntax error at or near "statement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statement BY NONE
< error
| 42601 syntax error at or near "statement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statement NAME omni_probe_name
< error
| 42601 syntax error at or near "statement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statement CYCLE
< error
| 42601 syntax error at or near "statement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statistics
< error
| 42601 syntax error at or near "statistics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statistics 1
< error
| 42601 syntax error at or near "statistics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statistics int
< error
| 42601 syntax error at or near "statistics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statistics BY NONE
< error
| 42601 syntax error at or near "statistics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statistics NAME omni_probe_name
< error
| 42601 syntax error at or near "statistics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq statistics CYCLE
< error
| 42601 syntax error at or near "statistics"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdin
< error
| 42601 syntax error at or near "stdin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdin 1
< error
| 42601 syntax error at or near "stdin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdin int
< error
| 42601 syntax error at or near "stdin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdin BY NONE
< error
| 42601 syntax error at or near "stdin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdin NAME omni_probe_name
< error
| 42601 syntax error at or near "stdin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdin CYCLE
< error
| 42601 syntax error at or near "stdin"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdout
< error
| 42601 syntax error at or near "stdout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdout 1
< error
| 42601 syntax error at or near "stdout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdout int
< error
| 42601 syntax error at or near "stdout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdout BY NONE
< error
| 42601 syntax error at or near "stdout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdout NAME omni_probe_name
< error
| 42601 syntax error at or near "stdout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stdout CYCLE
< error
| 42601 syntax error at or near "stdout"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq storage
< error
| 42601 syntax error at or near "storage"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq storage 1
< error
| 42601 syntax error at or near "storage"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq storage int
< error
| 42601 syntax error at or near "storage"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq storage BY NONE
< error
| 42601 syntax error at or near "storage"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq storage NAME omni_probe_name
< error
| 42601 syntax error at or near "storage"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq storage CYCLE
< error
| 42601 syntax error at or near "storage"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stored
< error
| 42601 syntax error at or near "stored"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stored 1
< error
| 42601 syntax error at or near "stored"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stored int
< error
| 42601 syntax error at or near "stored"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stored BY NONE
< error
| 42601 syntax error at or near "stored"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stored NAME omni_probe_name
< error
| 42601 syntax error at or near "stored"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq stored CYCLE
< error
| 42601 syntax error at or near "stored"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strict
< error
| 42601 syntax error at or near "strict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strict 1
< error
| 42601 syntax error at or near "strict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strict int
< error
| 42601 syntax error at or near "strict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strict BY NONE
< error
| 42601 syntax error at or near "strict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strict NAME omni_probe_name
< error
| 42601 syntax error at or near "strict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strict CYCLE
< error
| 42601 syntax error at or near "strict"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq string
< error
| 42601 syntax error at or near "string"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq string 1
< error
| 42601 syntax error at or near "string"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq string int
< error
| 42601 syntax error at or near "string"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq string BY NONE
< error
| 42601 syntax error at or near "string"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq string NAME omni_probe_name
< error
| 42601 syntax error at or near "string"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq string CYCLE
< error
| 42601 syntax error at or near "string"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strip
< error
| 42601 syntax error at or near "strip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strip 1
< error
| 42601 syntax error at or near "strip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strip int
< error
| 42601 syntax error at or near "strip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strip BY NONE
< error
| 42601 syntax error at or near "strip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strip NAME omni_probe_name
< error
| 42601 syntax error at or near "strip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq strip CYCLE
< error
| 42601 syntax error at or near "strip"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq subscription
< error
| 42601 syntax error at or near "subscription"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq subscription 1
< error
| 42601 syntax error at or near "subscription"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq subscription int
< error
| 42601 syntax error at or near "subscription"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq subscription BY NONE
< error
| 42601 syntax error at or near "subscription"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq subscription NAME omni_probe_name
< error
| 42601 syntax error at or near "subscription"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq subscription CYCLE
< error
| 42601 syntax error at or near "subscription"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq substring
< error
| 42601 syntax error at or near "substring"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq substring 1
< error
| 42601 syntax error at or near "substring"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq substring int
< error
| 42601 syntax error at or near "substring"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq substring BY NONE
< error
| 42601 syntax error at or near "substring"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq substring NAME omni_probe_name
< error
| 42601 syntax error at or near "substring"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq substring CYCLE
< error
| 42601 syntax error at or near "substring"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq support
< error
| 42601 syntax error at or near "support"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq support 1
< error
| 42601 syntax error at or near "support"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq support int
< error
| 42601 syntax error at or near "support"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq support BY NONE
< error
| 42601 syntax error at or near "support"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq support NAME omni_probe_name
< error
| 42601 syntax error at or near "support"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq support CYCLE
< error
| 42601 syntax error at or near "support"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq symmetric
< error
| 42601 syntax error at or near "symmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq symmetric 1
< error
| 42601 syntax error at or near "symmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq symmetric int
< error
| 42601 syntax error at or near "symmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq symmetric BY NONE
< error
| 42601 syntax error at or near "symmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq symmetric NAME omni_probe_name
< error
| 42601 syntax error at or near "symmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq symmetric CYCLE
< error
| 42601 syntax error at or near "symmetric"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sysid
< error
| 42601 syntax error at or near "sysid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sysid 1
< error
| 42601 syntax error at or near "sysid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sysid int
< error
| 42601 syntax error at or near "sysid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sysid BY NONE
< error
| 42601 syntax error at or near "sysid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sysid NAME omni_probe_name
< error
| 42601 syntax error at or near "sysid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq sysid CYCLE
< error
| 42601 syntax error at or near "sysid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system
< error
| 42601 syntax error at or near "system"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system 1
< error
| 42601 syntax error at or near "system"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system int
< error
| 42601 syntax error at or near "system"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system BY NONE
< error
| 42601 syntax error at or near "system"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system NAME omni_probe_name
< error
| 42601 syntax error at or near "system"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system CYCLE
< error
| 42601 syntax error at or near "system"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system_user
< error
| 42601 syntax error at or near "system_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system_user 1
< error
| 42601 syntax error at or near "system_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system_user int
< error
| 42601 syntax error at or near "system_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system_user BY NONE
< error
| 42601 syntax error at or near "system_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system_user NAME omni_probe_name
< error
| 42601 syntax error at or near "system_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq system_user CYCLE
< error
| 42601 syntax error at or near "system_user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq table
< error
| 42601 syntax error at or near "table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq table 1
< error
| 42601 syntax error at or near "table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq table int
< error
| 42601 syntax error at or near "table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq table BY NONE
< error
| 42601 syntax error at or near "table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq table NAME omni_probe_name
< error
| 42601 syntax error at or near "table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq table CYCLE
< error
| 42601 syntax error at or near "table"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tables
< error
| 42601 syntax error at or near "tables"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tables 1
< error
| 42601 syntax error at or near "tables"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tables int
< error
| 42601 syntax error at or near "tables"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tables BY NONE
< error
| 42601 syntax error at or near "tables"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tables NAME omni_probe_name
< error
| 42601 syntax error at or near "tables"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tables CYCLE
< error
| 42601 syntax error at or near "tables"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablesample
< error
| 42601 syntax error at or near "tablesample"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablesample 1
< error
| 42601 syntax error at or near "tablesample"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablesample int
< error
| 42601 syntax error at or near "tablesample"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablesample BY NONE
< error
| 42601 syntax error at or near "tablesample"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablesample NAME omni_probe_name
< error
| 42601 syntax error at or near "tablesample"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablesample CYCLE
< error
| 42601 syntax error at or near "tablesample"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablespace
< error
| 42601 syntax error at or near "tablespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablespace 1
< error
| 42601 syntax error at or near "tablespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablespace int
< error
| 42601 syntax error at or near "tablespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablespace BY NONE
< error
| 42601 syntax error at or near "tablespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablespace NAME omni_probe_name
< error
| 42601 syntax error at or near "tablespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq tablespace CYCLE
< error
| 42601 syntax error at or near "tablespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq target
< error
| 42601 syntax error at or near "target"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq target 1
< error
| 42601 syntax error at or near "target"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq target int
< error
| 42601 syntax error at or near "target"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq target BY NONE
< error
| 42601 syntax error at or near "target"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq target NAME omni_probe_name
< error
| 42601 syntax error at or near "target"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq target CYCLE
< error
| 42601 syntax error at or near "target"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temp
< error
| 42601 syntax error at or near "temp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temp 1
< error
| 42601 syntax error at or near "temp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temp int
< error
| 42601 syntax error at or near "temp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temp BY NONE
< error
| 42601 syntax error at or near "temp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temp NAME omni_probe_name
< error
| 42601 syntax error at or near "temp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temp CYCLE
< error
| 42601 syntax error at or near "temp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq template
< error
| 42601 syntax error at or near "template"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq template 1
< error
| 42601 syntax error at or near "template"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq template int
< error
| 42601 syntax error at or near "template"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq template BY NONE
< error
| 42601 syntax error at or near "template"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq template NAME omni_probe_name
< error
| 42601 syntax error at or near "template"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq template CYCLE
< error
| 42601 syntax error at or near "template"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temporary
< error
| 42601 syntax error at or near "temporary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temporary 1
< error
| 42601 syntax error at or near "temporary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temporary int
< error
| 42601 syntax error at or near "temporary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temporary BY NONE
< error
| 42601 syntax error at or near "temporary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temporary NAME omni_probe_name
< error
| 42601 syntax error at or near "temporary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq temporary CYCLE
< error
| 42601 syntax error at or near "temporary"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq text
< error
| 42601 syntax error at or near "text"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq text 1
< error
| 42601 syntax error at or near "text"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq text int
< error
| 42601 syntax error at or near "text"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq text BY NONE
< error
| 42601 syntax error at or near "text"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq text NAME omni_probe_name
< error
| 42601 syntax error at or near "text"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq text CYCLE
< error
| 42601 syntax error at or near "text"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq then
< error
| 42601 syntax error at or near "then"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq then 1
< error
| 42601 syntax error at or near "then"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq then int
< error
| 42601 syntax error at or near "then"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq then BY NONE
< error
| 42601 syntax error at or near "then"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq then NAME omni_probe_name
< error
| 42601 syntax error at or near "then"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq then CYCLE
< error
| 42601 syntax error at or near "then"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ties
< error
| 42601 syntax error at or near "ties"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ties 1
< error
| 42601 syntax error at or near "ties"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ties int
< error
| 42601 syntax error at or near "ties"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ties BY NONE
< error
| 42601 syntax error at or near "ties"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ties NAME omni_probe_name
< error
| 42601 syntax error at or near "ties"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq ties CYCLE
< error
| 42601 syntax error at or near "ties"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq time
< error
| 42601 syntax error at or near "time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq time 1
< error
| 42601 syntax error at or near "time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq time int
< error
| 42601 syntax error at or near "time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq time BY NONE
< error
| 42601 syntax error at or near "time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq time NAME omni_probe_name
< error
| 42601 syntax error at or near "time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq time CYCLE
< error
| 42601 syntax error at or near "time"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq timestamp
< error
| 42601 syntax error at or near "timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq timestamp 1
< error
| 42601 syntax error at or near "timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq timestamp int
< error
| 42601 syntax error at or near "timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq timestamp BY NONE
< error
| 42601 syntax error at or near "timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq timestamp NAME omni_probe_name
< error
| 42601 syntax error at or near "timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq timestamp CYCLE
< error
| 42601 syntax error at or near "timestamp"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq to
< error
| 42601 syntax error at or near "to"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq to 1
< error
| 42601 syntax error at or near "to"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq to int
< error
| 42601 syntax error at or near "to"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq to BY NONE
< error
| 42601 syntax error at or near "to"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq to NAME omni_probe_name
< error
| 42601 syntax error at or near "to"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq to CYCLE
< error
| 42601 syntax error at or near "to"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trailing
< error
| 42601 syntax error at or near "trailing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trailing 1
< error
| 42601 syntax error at or near "trailing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trailing int
< error
| 42601 syntax error at or near "trailing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trailing BY NONE
< error
| 42601 syntax error at or near "trailing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trailing NAME omni_probe_name
< error
| 42601 syntax error at or near "trailing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trailing CYCLE
< error
| 42601 syntax error at or near "trailing"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transaction
< error
| 42601 syntax error at or near "transaction"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transaction 1
< error
| 42601 syntax error at or near "transaction"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transaction int
< error
| 42601 syntax error at or near "transaction"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transaction BY NONE
< error
| 42601 syntax error at or near "transaction"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transaction NAME omni_probe_name
< error
| 42601 syntax error at or near "transaction"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transaction CYCLE
< error
| 42601 syntax error at or near "transaction"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transform
< error
| 42601 syntax error at or near "transform"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transform 1
< error
| 42601 syntax error at or near "transform"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transform int
< error
| 42601 syntax error at or near "transform"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transform BY NONE
< error
| 42601 syntax error at or near "transform"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transform NAME omni_probe_name
< error
| 42601 syntax error at or near "transform"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq transform CYCLE
< error
| 42601 syntax error at or near "transform"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq treat
< error
| 42601 syntax error at or near "treat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq treat 1
< error
| 42601 syntax error at or near "treat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq treat int
< error
| 42601 syntax error at or near "treat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq treat BY NONE
< error
| 42601 syntax error at or near "treat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq treat NAME omni_probe_name
< error
| 42601 syntax error at or near "treat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq treat CYCLE
< error
| 42601 syntax error at or near "treat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trigger
< error
| 42601 syntax error at or near "trigger"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trigger 1
< error
| 42601 syntax error at or near "trigger"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trigger int
< error
| 42601 syntax error at or near "trigger"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trigger BY NONE
< error
| 42601 syntax error at or near "trigger"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trigger NAME omni_probe_name
< error
| 42601 syntax error at or near "trigger"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trigger CYCLE
< error
| 42601 syntax error at or near "trigger"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trim
< error
| 42601 syntax error at or near "trim"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trim 1
< error
| 42601 syntax error at or near "trim"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trim int
< error
| 42601 syntax error at or near "trim"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trim BY NONE
< error
| 42601 syntax error at or near "trim"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trim NAME omni_probe_name
< error
| 42601 syntax error at or near "trim"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trim CYCLE
< error
| 42601 syntax error at or near "trim"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq true
< error
| 42601 syntax error at or near "true"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq true 1
< error
| 42601 syntax error at or near "true"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq true int
< error
| 42601 syntax error at or near "true"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq true BY NONE
< error
| 42601 syntax error at or near "true"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq true NAME omni_probe_name
< error
| 42601 syntax error at or near "true"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq true CYCLE
< error
| 42601 syntax error at or near "true"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq truncate
< error
| 42601 syntax error at or near "truncate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq truncate 1
< error
| 42601 syntax error at or near "truncate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq truncate int
< error
| 42601 syntax error at or near "truncate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq truncate BY NONE
< error
| 42601 syntax error at or near "truncate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq truncate NAME omni_probe_name
< error
| 42601 syntax error at or near "truncate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq truncate CYCLE
< error
| 42601 syntax error at or near "truncate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trusted
< error
| 42601 syntax error at or near "trusted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trusted 1
< error
| 42601 syntax error at or near "trusted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trusted int
< error
| 42601 syntax error at or near "trusted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trusted BY NONE
< error
| 42601 syntax error at or near "trusted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trusted NAME omni_probe_name
< error
| 42601 syntax error at or near "trusted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq trusted CYCLE
< error
| 42601 syntax error at or near "trusted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq type
< error
| 42601 syntax error at or near "type"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq type 1
< error
| 42601 syntax error at or near "type"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq type int
< error
| 42601 syntax error at or near "type"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq type BY NONE
< error
| 42601 syntax error at or near "type"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq type NAME omni_probe_name
< error
| 42601 syntax error at or near "type"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq type CYCLE
< error
| 42601 syntax error at or near "type"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq types
< error
| 42601 syntax error at or near "types"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq types 1
< error
| 42601 syntax error at or near "types"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq types int
< error
| 42601 syntax error at or near "types"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq types BY NONE
< error
| 42601 syntax error at or near "types"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq types NAME omni_probe_name
< error
| 42601 syntax error at or near "types"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq types CYCLE
< error
| 42601 syntax error at or near "types"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uescape
< error
| 42601 syntax error at or near "uescape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uescape 1
< error
| 42601 syntax error at or near "uescape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uescape int
< error
| 42601 syntax error at or near "uescape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uescape BY NONE
< error
| 42601 syntax error at or near "uescape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uescape NAME omni_probe_name
< error
| 42601 syntax error at or near "uescape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uescape CYCLE
< error
| 42601 syntax error at or near "uescape"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unbounded
< error
| 42601 syntax error at or near "unbounded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unbounded 1
< error
| 42601 syntax error at or near "unbounded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unbounded int
< error
| 42601 syntax error at or near "unbounded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unbounded BY NONE
< error
| 42601 syntax error at or near "unbounded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unbounded NAME omni_probe_name
< error
| 42601 syntax error at or near "unbounded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unbounded CYCLE
< error
| 42601 syntax error at or near "unbounded"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uncommitted
< error
| 42601 syntax error at or near "uncommitted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uncommitted 1
< error
| 42601 syntax error at or near "uncommitted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uncommitted int
< error
| 42601 syntax error at or near "uncommitted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uncommitted BY NONE
< error
| 42601 syntax error at or near "uncommitted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uncommitted NAME omni_probe_name
< error
| 42601 syntax error at or near "uncommitted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq uncommitted CYCLE
< error
| 42601 syntax error at or near "uncommitted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unconditional
< error
| 42601 syntax error at or near "unconditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unconditional 1
< error
| 42601 syntax error at or near "unconditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unconditional int
< error
| 42601 syntax error at or near "unconditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unconditional BY NONE
< error
| 42601 syntax error at or near "unconditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unconditional NAME omni_probe_name
< error
| 42601 syntax error at or near "unconditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unconditional CYCLE
< error
| 42601 syntax error at or near "unconditional"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unencrypted
< error
| 42601 syntax error at or near "unencrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unencrypted 1
< error
| 42601 syntax error at or near "unencrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unencrypted int
< error
| 42601 syntax error at or near "unencrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unencrypted BY NONE
< error
| 42601 syntax error at or near "unencrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unencrypted NAME omni_probe_name
< error
| 42601 syntax error at or near "unencrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unencrypted CYCLE
< error
| 42601 syntax error at or near "unencrypted"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq union
< error
| 42601 syntax error at or near "union"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq union 1
< error
| 42601 syntax error at or near "union"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq union int
< error
| 42601 syntax error at or near "union"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq union BY NONE
< error
| 42601 syntax error at or near "union"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq union NAME omni_probe_name
< error
| 42601 syntax error at or near "union"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq union CYCLE
< error
| 42601 syntax error at or near "union"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unique
< error
| 42601 syntax error at or near "unique"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unique 1
< error
| 42601 syntax error at or near "unique"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unique int
< error
| 42601 syntax error at or near "unique"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unique BY NONE
< error
| 42601 syntax error at or near "unique"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unique NAME omni_probe_name
< error
| 42601 syntax error at or near "unique"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unique CYCLE
< error
| 42601 syntax error at or near "unique"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unknown
< error
| 42601 syntax error at or near "unknown"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unknown 1
< error
| 42601 syntax error at or near "unknown"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unknown int
< error
| 42601 syntax error at or near "unknown"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unknown BY NONE
< error
| 42601 syntax error at or near "unknown"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unknown NAME omni_probe_name
< error
| 42601 syntax error at or near "unknown"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unknown CYCLE
< error
| 42601 syntax error at or near "unknown"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unlisten
< error
| 42601 syntax error at or near "unlisten"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unlisten 1
< error
| 42601 syntax error at or near "unlisten"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unlisten int
< error
| 42601 syntax error at or near "unlisten"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unlisten BY NONE
< error
| 42601 syntax error at or near "unlisten"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unlisten NAME omni_probe_name
< error
| 42601 syntax error at or near "unlisten"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unlisten CYCLE
< error
| 42601 syntax error at or near "unlisten"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq unlogged
< error
| XX000 option "unlogged" not recognized
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq until
< error
| 42601 syntax error at or near "until"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq until 1
< error
| 42601 syntax error at or near "until"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq until int
< error
| 42601 syntax error at or near "until"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq until BY NONE
< error
| 42601 syntax error at or near "until"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq until NAME omni_probe_name
< error
| 42601 syntax error at or near "until"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq until CYCLE
< error
| 42601 syntax error at or near "until"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq update
< error
| 42601 syntax error at or near "update"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq update 1
< error
| 42601 syntax error at or near "update"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq update int
< error
| 42601 syntax error at or near "update"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq update BY NONE
< error
| 42601 syntax error at or near "update"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq update NAME omni_probe_name
< error
| 42601 syntax error at or near "update"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq update CYCLE
< error
| 42601 syntax error at or near "update"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq user
< error
| 42601 syntax error at or near "user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq user 1
< error
| 42601 syntax error at or near "user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq user int
< error
| 42601 syntax error at or near "user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq user BY NONE
< error
| 42601 syntax error at or near "user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq user NAME omni_probe_name
< error
| 42601 syntax error at or near "user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq user CYCLE
< error
| 42601 syntax error at or near "user"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq using
< error
| 42601 syntax error at or near "using"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq using 1
< error
| 42601 syntax error at or near "using"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq using int
< error
| 42601 syntax error at or near "using"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq using BY NONE
< error
| 42601 syntax error at or near "using"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq using NAME omni_probe_name
< error
| 42601 syntax error at or near "using"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq using CYCLE
< error
| 42601 syntax error at or near "using"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq vacuum
< error
| 42601 syntax error at or near "vacuum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq vacuum 1
< error
| 42601 syntax error at or near "vacuum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq vacuum int
< error
| 42601 syntax error at or near "vacuum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq vacuum BY NONE
< error
| 42601 syntax error at or near "vacuum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq vacuum NAME omni_probe_name
< error
| 42601 syntax error at or near "vacuum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq vacuum CYCLE
< error
| 42601 syntax error at or near "vacuum"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq valid
< error
| 42601 syntax error at or near "valid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq valid 1
< error
| 42601 syntax error at or near "valid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq valid int
< error
| 42601 syntax error at or near "valid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq valid BY NONE
< error
| 42601 syntax error at or near "valid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq valid NAME omni_probe_name
< error
| 42601 syntax error at or near "valid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq valid CYCLE
< error
| 42601 syntax error at or near "valid"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validate
< error
| 42601 syntax error at or near "validate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validate 1
< error
| 42601 syntax error at or near "validate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validate int
< error
| 42601 syntax error at or near "validate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validate BY NONE
< error
| 42601 syntax error at or near "validate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validate NAME omni_probe_name
< error
| 42601 syntax error at or near "validate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validate CYCLE
< error
| 42601 syntax error at or near "validate"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validator
< error
| 42601 syntax error at or near "validator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validator 1
< error
| 42601 syntax error at or near "validator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validator int
< error
| 42601 syntax error at or near "validator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validator BY NONE
< error
| 42601 syntax error at or near "validator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validator NAME omni_probe_name
< error
| 42601 syntax error at or near "validator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq validator CYCLE
< error
| 42601 syntax error at or near "validator"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq value
< error
| 42601 syntax error at or near "value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq value 1
< error
| 42601 syntax error at or near "value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq value int
< error
| 42601 syntax error at or near "value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq value BY NONE
< error
| 42601 syntax error at or near "value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq value NAME omni_probe_name
< error
| 42601 syntax error at or near "value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq value CYCLE
< error
| 42601 syntax error at or near "value"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq values
< error
| 42601 syntax error at or near "values"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq values 1
< error
| 42601 syntax error at or near "values"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq values int
< error
| 42601 syntax error at or near "values"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq values BY NONE
< error
| 42601 syntax error at or near "values"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq values NAME omni_probe_name
< error
| 42601 syntax error at or near "values"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq values CYCLE
< error
| 42601 syntax error at or near "values"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varchar
< error
| 42601 syntax error at or near "varchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varchar 1
< error
| 42601 syntax error at or near "varchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varchar int
< error
| 42601 syntax error at or near "varchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varchar BY NONE
< error
| 42601 syntax error at or near "varchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varchar NAME omni_probe_name
< error
| 42601 syntax error at or near "varchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varchar CYCLE
< error
| 42601 syntax error at or near "varchar"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq variadic
< error
| 42601 syntax error at or near "variadic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq variadic 1
< error
| 42601 syntax error at or near "variadic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq variadic int
< error
| 42601 syntax error at or near "variadic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq variadic BY NONE
< error
| 42601 syntax error at or near "variadic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq variadic NAME omni_probe_name
< error
| 42601 syntax error at or near "variadic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq variadic CYCLE
< error
| 42601 syntax error at or near "variadic"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varying
< error
| 42601 syntax error at or near "varying"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varying 1
< error
| 42601 syntax error at or near "varying"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varying int
< error
| 42601 syntax error at or near "varying"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varying BY NONE
< error
| 42601 syntax error at or near "varying"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varying NAME omni_probe_name
< error
| 42601 syntax error at or near "varying"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq varying CYCLE
< error
| 42601 syntax error at or near "varying"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq verbose
< error
| 42601 syntax error at or near "verbose"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq verbose 1
< error
| 42601 syntax error at or near "verbose"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq verbose int
< error
| 42601 syntax error at or near "verbose"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq verbose BY NONE
< error
| 42601 syntax error at or near "verbose"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq verbose NAME omni_probe_name
< error
| 42601 syntax error at or near "verbose"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq verbose CYCLE
< error
| 42601 syntax error at or near "verbose"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq version
< error
| 42601 syntax error at or near "version"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq version 1
< error
| 42601 syntax error at or near "version"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq version int
< error
| 42601 syntax error at or near "version"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq version BY NONE
< error
| 42601 syntax error at or near "version"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq version NAME omni_probe_name
< error
| 42601 syntax error at or near "version"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq version CYCLE
< error
| 42601 syntax error at or near "version"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq view
< error
| 42601 syntax error at or near "view"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq view 1
< error
| 42601 syntax error at or near "view"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq view int
< error
| 42601 syntax error at or near "view"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq view BY NONE
< error
| 42601 syntax error at or near "view"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq view NAME omni_probe_name
< error
| 42601 syntax error at or near "view"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq view CYCLE
< error
| 42601 syntax error at or near "view"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq views
< error
| 42601 syntax error at or near "views"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq views 1
< error
| 42601 syntax error at or near "views"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq views int
< error
| 42601 syntax error at or near "views"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq views BY NONE
< error
| 42601 syntax error at or near "views"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq views NAME omni_probe_name
< error
| 42601 syntax error at or near "views"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq views CYCLE
< error
| 42601 syntax error at or near "views"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq volatile
< error
| 42601 syntax error at or near "volatile"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq volatile 1
< error
| 42601 syntax error at or near "volatile"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq volatile int
< error
| 42601 syntax error at or near "volatile"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq volatile BY NONE
< error
| 42601 syntax error at or near "volatile"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq volatile NAME omni_probe_name
< error
| 42601 syntax error at or near "volatile"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq volatile CYCLE
< error
| 42601 syntax error at or near "volatile"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq when
< error
| 42601 syntax error at or near "when"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq when 1
< error
| 42601 syntax error at or near "when"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq when int
< error
| 42601 syntax error at or near "when"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq when BY NONE
< error
| 42601 syntax error at or near "when"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq when NAME omni_probe_name
< error
| 42601 syntax error at or near "when"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq when CYCLE
< error
| 42601 syntax error at or near "when"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq where
< error
| 42601 syntax error at or near "where"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq where 1
< error
| 42601 syntax error at or near "where"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq where int
< error
| 42601 syntax error at or near "where"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq where BY NONE
< error
| 42601 syntax error at or near "where"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq where NAME omni_probe_name
< error
| 42601 syntax error at or near "where"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq where CYCLE
< error
| 42601 syntax error at or near "where"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq whitespace
< error
| 42601 syntax error at or near "whitespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq whitespace 1
< error
| 42601 syntax error at or near "whitespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq whitespace int
< error
| 42601 syntax error at or near "whitespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq whitespace BY NONE
< error
| 42601 syntax error at or near "whitespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq whitespace NAME omni_probe_name
< error
| 42601 syntax error at or near "whitespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq whitespace CYCLE
< error
| 42601 syntax error at or near "whitespace"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq window
< error
| 42601 syntax error at or near "window"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq window 1
< error
| 42601 syntax error at or near "window"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq window int
< error
| 42601 syntax error at or near "window"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq window BY NONE
< error
| 42601 syntax error at or near "window"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq window NAME omni_probe_name
< error
| 42601 syntax error at or near "window"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq window CYCLE
< error
| 42601 syntax error at or near "window"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq with
< error
| 42601 syntax error at or near "with"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq with 1
< error
| 42601 syntax error at or near "with"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq with int
< error
| 42601 syntax error at or near "with"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq with BY NONE
< error
| 42601 syntax error at or near "with"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq with NAME omni_probe_name
< error
| 42601 syntax error at or near "with"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq with CYCLE
< error
| 42601 syntax error at or near "with"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq within
< error
| 42601 syntax error at or near "within"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq within 1
< error
| 42601 syntax error at or near "within"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq within int
< error
| 42601 syntax error at or near "within"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq within BY NONE
< error
| 42601 syntax error at or near "within"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq within NAME omni_probe_name
< error
| 42601 syntax error at or near "within"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq within CYCLE
< error
| 42601 syntax error at or near "within"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq without
< error
| 42601 syntax error at or near "without"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq without 1
< error
| 42601 syntax error at or near "without"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq without int
< error
| 42601 syntax error at or near "without"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq without BY NONE
< error
| 42601 syntax error at or near "without"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq without NAME omni_probe_name
< error
| 42601 syntax error at or near "without"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq without CYCLE
< error
| 42601 syntax error at or near "without"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq work
< error
| 42601 syntax error at or near "work"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq work 1
< error
| 42601 syntax error at or near "work"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq work int
< error
| 42601 syntax error at or near "work"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq work BY NONE
< error
| 42601 syntax error at or near "work"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq work NAME omni_probe_name
< error
| 42601 syntax error at or near "work"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq work CYCLE
< error
| 42601 syntax error at or near "work"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq wrapper
< error
| 42601 syntax error at or near "wrapper"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq wrapper 1
< error
| 42601 syntax error at or near "wrapper"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq wrapper int
< error
| 42601 syntax error at or near "wrapper"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq wrapper BY NONE
< error
| 42601 syntax error at or near "wrapper"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq wrapper NAME omni_probe_name
< error
| 42601 syntax error at or near "wrapper"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq wrapper CYCLE
< error
| 42601 syntax error at or near "wrapper"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq write
< error
| 42601 syntax error at or near "write"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq write 1
< error
| 42601 syntax error at or near "write"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq write int
< error
| 42601 syntax error at or near "write"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq write BY NONE
< error
| 42601 syntax error at or near "write"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq write NAME omni_probe_name
< error
| 42601 syntax error at or near "write"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq write CYCLE
< error
| 42601 syntax error at or near "write"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xml
< error
| 42601 syntax error at or near "xml"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xml 1
< error
| 42601 syntax error at or near "xml"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xml int
< error
| 42601 syntax error at or near "xml"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xml BY NONE
< error
| 42601 syntax error at or near "xml"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xml NAME omni_probe_name
< error
| 42601 syntax error at or near "xml"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xml CYCLE
< error
| 42601 syntax error at or near "xml"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlattributes
< error
| 42601 syntax error at or near "xmlattributes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlattributes 1
< error
| 42601 syntax error at or near "xmlattributes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlattributes int
< error
| 42601 syntax error at or near "xmlattributes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlattributes BY NONE
< error
| 42601 syntax error at or near "xmlattributes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlattributes NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlattributes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlattributes CYCLE
< error
| 42601 syntax error at or near "xmlattributes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlconcat
< error
| 42601 syntax error at or near "xmlconcat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlconcat 1
< error
| 42601 syntax error at or near "xmlconcat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlconcat int
< error
| 42601 syntax error at or near "xmlconcat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlconcat BY NONE
< error
| 42601 syntax error at or near "xmlconcat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlconcat NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlconcat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlconcat CYCLE
< error
| 42601 syntax error at or near "xmlconcat"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlelement
< error
| 42601 syntax error at or near "xmlelement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlelement 1
< error
| 42601 syntax error at or near "xmlelement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlelement int
< error
| 42601 syntax error at or near "xmlelement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlelement BY NONE
< error
| 42601 syntax error at or near "xmlelement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlelement NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlelement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlelement CYCLE
< error
| 42601 syntax error at or near "xmlelement"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlexists
< error
| 42601 syntax error at or near "xmlexists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlexists 1
< error
| 42601 syntax error at or near "xmlexists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlexists int
< error
| 42601 syntax error at or near "xmlexists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlexists BY NONE
< error
| 42601 syntax error at or near "xmlexists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlexists NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlexists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlexists CYCLE
< error
| 42601 syntax error at or near "xmlexists"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlforest
< error
| 42601 syntax error at or near "xmlforest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlforest 1
< error
| 42601 syntax error at or near "xmlforest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlforest int
< error
| 42601 syntax error at or near "xmlforest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlforest BY NONE
< error
| 42601 syntax error at or near "xmlforest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlforest NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlforest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlforest CYCLE
< error
| 42601 syntax error at or near "xmlforest"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlnamespaces
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlnamespaces 1
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlnamespaces int
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlnamespaces BY NONE
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlnamespaces NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlnamespaces CYCLE
< error
| 42601 syntax error at or near "xmlnamespaces"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlparse
< error
| 42601 syntax error at or near "xmlparse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlparse 1
< error
| 42601 syntax error at or near "xmlparse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlparse int
< error
| 42601 syntax error at or near "xmlparse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlparse BY NONE
< error
| 42601 syntax error at or near "xmlparse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlparse NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlparse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlparse CYCLE
< error
| 42601 syntax error at or near "xmlparse"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlpi
< error
| 42601 syntax error at or near "xmlpi"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlpi 1
< error
| 42601 syntax error at or near "xmlpi"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlpi int
< error
| 42601 syntax error at or near "xmlpi"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlpi BY NONE
< error
| 42601 syntax error at or near "xmlpi"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlpi NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlpi"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlpi CYCLE
< error
| 42601 syntax error at or near "xmlpi"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlroot
< error
| 42601 syntax error at or near "xmlroot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlroot 1
< error
| 42601 syntax error at or near "xmlroot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlroot int
< error
| 42601 syntax error at or near "xmlroot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlroot BY NONE
< error
| 42601 syntax error at or near "xmlroot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlroot NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlroot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlroot CYCLE
< error
| 42601 syntax error at or near "xmlroot"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlserialize
< error
| 42601 syntax error at or near "xmlserialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlserialize 1
< error
| 42601 syntax error at or near "xmlserialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlserialize int
< error
| 42601 syntax error at or near "xmlserialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlserialize BY NONE
< error
| 42601 syntax error at or near "xmlserialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlserialize NAME omni_probe_name
< error
| 42601 syntax error at or near "xmlserialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmlserialize CYCLE
< error
| 42601 syntax error at or near "xmlserialize"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmltable
< error
| 42601 syntax error at or near "xmltable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmltable 1
< error
| 42601 syntax error at or near "xmltable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmltable int
< error
| 42601 syntax error at or near "xmltable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmltable BY NONE
< error
| 42601 syntax error at or near "xmltable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmltable NAME omni_probe_name
< error
| 42601 syntax error at or near "xmltable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq xmltable CYCLE
< error
| 42601 syntax error at or near "xmltable"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq year
< error
| 42601 syntax error at or near "year"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq year 1
< error
| 42601 syntax error at or near "year"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq year int
< error
| 42601 syntax error at or near "year"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq year BY NONE
< error
| 42601 syntax error at or near "year"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq year NAME omni_probe_name
< error
| 42601 syntax error at or near "year"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq year CYCLE
< error
| 42601 syntax error at or near "year"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq yes
< error
| 42601 syntax error at or near "yes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq yes 1
< error
| 42601 syntax error at or near "yes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq yes int
< error
| 42601 syntax error at or near "yes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq yes BY NONE
< error
| 42601 syntax error at or near "yes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq yes NAME omni_probe_name
< error
| 42601 syntax error at or near "yes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq yes CYCLE
< error
| 42601 syntax error at or near "yes"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq zone
< error
| 42601 syntax error at or near "zone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq zone 1
< error
| 42601 syntax error at or near "zone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq zone int
< error
| 42601 syntax error at or near "zone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq zone BY NONE
< error
| 42601 syntax error at or near "zone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq zone NAME omni_probe_name
< error
| 42601 syntax error at or near "zone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| CREATE SEQUENCE omni_probe_seq zone CYCLE
< error
| 42601 syntax error at or near "zone"
| position: 32
> rollback
< ok
> begin
< ok
> exec
| ALTER SEQUENCE omni_probe_seq NO CYCLE
< error
| 42P01 relation "omni_probe_seq" does not exist
> rollback
< ok
> begin
< ok
> exec
| ALTER SEQUENCE omni_probe_seq NO MAXVALUE
< error
| 42P01 relation "omni_probe_seq" does not exist
> rollback
< ok
> begin
< ok
> exec
| ALTER SEQUENCE omni_probe_seq NO MINVALUE
< error
| 42P01 relation "omni_probe_seq" does not exist
> rollback
< ok
> begin
< ok
> exec
| ALTER SEQUENCE omni_probe_seq NO INHERIT omni_probe_table
< error
| 42P01 relation "omni_probe_seq" does not exist
> rollback
< ok
> begin
< ok
> exec
| ALTER SEQUENCE omni_probe_seq NO FORCE ROW LEVEL SECURITY
< error
| 42P01 relation "omni_probe_seq" does not exist
> rollback
< ok
