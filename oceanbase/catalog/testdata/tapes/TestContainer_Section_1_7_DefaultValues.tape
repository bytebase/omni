-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_1_7_DefaultValues default#1
> exec
| CREATE TABLE t_def_exjson (a JSON DEFAULT (JSON_ARRAY()))
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'a' can't have a default value
