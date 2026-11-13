-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseTableWithoutColumns/no_visible_column default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE a (KEY (b))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| CREATE TABLE IF NOT EXISTS a (KEY (b))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| CREATE TABLE z (KEY (zz))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| CREATE TABLE z (FOREIGN KEY (b) REFERENCES a (b))
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| CREATE TABLE z (KEY (b)) SELECT zz
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| CREATE TABLE z (b INT INVISIBLE)
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| CREATE TABLE z (b INT INVISIBLE, c INT)
< ok 0 0
