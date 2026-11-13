-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseViewColumnNames default#1
> exec
| CREATE VIEW v64b AS SELECT concat('ééééééééééééééééééééééééééé')
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v64b
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| concat('ééééééééééééééééééééééééééé')
> exec
| CREATE VIEW v66b AS SELECT concat('éééééééééééééééééééééééééééé')
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v66b
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| Name_exp_1
> exec
| CREATE VIEW v64c AS SELECT concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee')
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v64c
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee')
> exec
| CREATE VIEW v65c AS SELECT concat('eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee')
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| v65c
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| Name_exp_1
