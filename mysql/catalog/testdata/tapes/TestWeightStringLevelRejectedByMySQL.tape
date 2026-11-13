-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestWeightStringLevelRejectedByMySQL default#1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' as char(4)) AS `x`
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' level 1) AS `x`
< error
| 1064 42000 …near 'level 1) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' level 0) AS `x`
< error
| 1064 42000 …near 'level 0) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' level 1 asc) AS `x`
< error
| 1064 42000 …near 'level 1 asc) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' level 1 desc) AS `x`
< error
| 1064 42000 …near 'level 1 desc) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' level 1 reverse) AS `x`
< error
| 1064 42000 …near 'level 1 reverse) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' level 1 desc reverse) AS `x`
< error
| 1064 42000 …near 'level 1 desc reverse) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' level 1 reverse desc) AS `x`
< error
| 1064 42000 …near 'level 1 reverse desc) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' level 1, 2) AS `x`
< error
| 1064 42000 …near 'level 1, 2) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' level 1 - 2) AS `x`
< error
| 1064 42000 …near 'level 1 - 2) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' as char(4) level 1) AS `x`
< error
| 1064 42000 …near 'level 1) AS `x`' at line 1
> exec
| CREATE OR REPLACE VIEW `v` AS select weight_string('ab' as binary(4) level 1) AS `x`
< error
| 1064 42000 …near 'level 1) AS `x`' at line 1
