-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseColumnGroups/CREATE_TABLE_LIKE_copies_them default#1
> exec
| CREATE TABLE s (a INT PRIMARY KEY) WITH COLUMN GROUP(each column)
< ok 0 0
> exec
| CREATE TABLE t LIKE s
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN)
< error
| 1060 42S21 Duplicate column group name 'each column group'
> exec
| ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)
< ok 0 0

== TestOracle_OceanBaseColumnGroups/a_group_that_names_its_columns default#1
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT) WITH COLUMN GROUP(g(a))
< error
| 1235 0A000 Not supported feature or function
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT) WITH COLUMN GROUP(all columns, g(a, b))
< error
| 1235 0A000 Not supported feature or function
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN GROUP(g(a, b))
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, g(a))
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t ADD COLUMN GROUP(`each`(a))
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t DROP COLUMN GROUP(g(a))
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS, g(a))
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE nosuch ADD COLUMN GROUP(g(a))
< error
| 1146 42S02 Table 'omni_t_oceanbase_catalog_column_groups_00.nosuch' doesn't exist

== TestOracle_OceanBaseColumnGroups/a_group_written_twice default#1
> exec
| CREATE TABLE t (a INT PRIMARY KEY) WITH COLUMN GROUP(all columns, all columns)
< error
| 1060 42S21 Duplicate column group name 'all column group'
> exec
| CREATE TABLE t (a INT PRIMARY KEY) WITH COLUMN GROUP(each column, each column)
< error
| 1060 42S21 Duplicate column group name 'single column group'
> exec
| CREATE TABLE t (a INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN GROUP(g(a), ALL COLUMNS, ALL COLUMNS)
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS, ALL COLUMNS, g(a))
< error
| 1060 42S21 Duplicate column group name 'all column group'
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, EACH COLUMN)
< error
| 1060 42S21 Duplicate column group name 'single column group'
> exec
| ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS, ALL COLUMNS)
< error
| 1060 42S21 Duplicate column group name 'all column group'
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN)
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN, EACH COLUMN)
< error
| 1060 42S21 Duplicate column group name 'single column group'

== TestOracle_OceanBaseColumnGroups/add_and_drop default#1
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT)
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS)
< error
| 4185 HY000 Column group 'all column group' not found
> exec
| ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)
< error
| 4185 HY000 Column group 'each column group' not found
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN)
< ok 0 0
> exec ^
< error
| 1060 42S21 Duplicate column group name 'each column group'
> exec
| ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS, EACH COLUMN)
< error
| 1060 42S21 Duplicate column group name 'each column group'
> exec
| ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN, ALL COLUMNS)
< error
| 4185 HY000 Column group 'all column group' not found
> exec
| ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, ALL COLUMNS)
< error
| 1060 42S21 Duplicate column group name 'all column group'
> exec
| ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS, EACH COLUMN)
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)
< error
| 4185 HY000 Column group 'each column group' not found
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, ALL COLUMNS)
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS)
< ok 0 0
> exec ^
< error
| 4185 HY000 Column group 'all column group' not found
> exec
| ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)
< ok 0 0

== TestOracle_OceanBaseColumnGroups/before_the_table's_name default#1
> exec
| CREATE TABLE t (a INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT) WITH COLUMN GROUP(g(a))
< error
| 1235 0A000 Not supported feature or function
> exec
| CREATE TABLE IF NOT EXISTS t (a INT) WITH COLUMN GROUP(g(a))
< error
| 1235 0A000 Not supported feature or function
> exec
| CREATE TABLE IF NOT EXISTS t (a INT) WITH COLUMN GROUP(all columns, all columns)
< error
| 1060 42S21 Duplicate column group name 'all column group'
> exec
| CREATE TABLE IF NOT EXISTS t (a INT) WITH COLUMN GROUP(all columns)
< ok 0 0
> exec
| CREATE TABLE t (a INT) WITH COLUMN GROUP(all columns)
< error
| 1050 42S01 Table 't' already exists

== TestOracle_OceanBaseColumnGroups/delayed default#1
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS) DELAYED
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS, ALL COLUMNS) DELAYED
< error
| 1060 42S21 Duplicate column group name 'all column group'
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, g(a)) DELAYED
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS) DELAYED
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN) DELAYED
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN) DELAYED
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN, EACH COLUMN) DELAYED
< error
| 1235 0A000 Not supported feature or function
> exec
| ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN) DELAYED
< error
| 1060 42S21 Duplicate column group name 'each column group'
> exec
| ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS, EACH COLUMN) DELAYED
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS) DELAYED
< error
| 1235 0A000 Not supported feature or function
