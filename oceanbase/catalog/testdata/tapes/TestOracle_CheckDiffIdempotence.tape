-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_CheckDiffIdempotence default#1
> exec
| DROP DATABASE IF EXISTS chk_unnamed_single
< ok 0 0
> exec
| CREATE DATABASE chk_unnamed_single
< ok 1 0
> exec
| USE chk_unnamed_single
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_unnamed_single.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000001` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_named_single
< ok 0 0
> exec
| CREATE DATABASE chk_named_single
< ok 1 0
> exec
| USE chk_named_single
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT chk_a CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_named_single.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `chk_a` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_multiple_unnamed
< ok 0 0
> exec
| CREATE DATABASE chk_multiple_unnamed
< ok 1 0
> exec
| USE chk_multiple_unnamed
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CHECK (b < 100), CHECK (a < b))
< ok 0 0
> query
| SHOW CREATE TABLE chk_multiple_unnamed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000002` CHECK ((`a` > 0)),
|   CONSTRAINT `t_OBCHECK_0000000000000003` CHECK ((`b` < 100)),
|   CONSTRAINT `t_OBCHECK_0000000000000004` CHECK ((`a` < `b`))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_mixed_named_unnamed
< ok 0 0
> exec
| CREATE DATABASE chk_mixed_named_unnamed
< ok 1 0
> exec
| USE chk_mixed_named_unnamed
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck1 CHECK (a > 0), CHECK (b < 100))
< ok 0 0
> query
| SHOW CREATE TABLE chk_mixed_named_unnamed.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ck1` CHECK ((`a` > 0)),
|   CONSTRAINT `t_OBCHECK_0000000000000005` CHECK ((`b` < 100))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_not_enforced
< ok 0 0
> exec
| CREATE DATABASE chk_not_enforced
< ok 1 0
> exec
| USE chk_not_enforced
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE chk_not_enforced.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `ne` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_enforced_explicit
< ok 0 0
> exec
| CREATE DATABASE chk_enforced_explicit
< ok 1 0
> exec
| USE chk_enforced_explicit
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT en CHECK (a > 0) ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE chk_enforced_explicit.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `en` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_compound_and
< ok 0 0
> exec
| CREATE DATABASE chk_compound_and
< ok 1 0
> exec
| USE chk_compound_and
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))
< ok 0 0
> query
| SHOW CREATE TABLE chk_compound_and.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000006` CHECK (((`a` >= 0) and (`b` <= 100)))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_compound_or
< ok 0 0
> exec
| CREATE DATABASE chk_compound_or
< ok 1 0
> exec
| USE chk_compound_or
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a = 1 OR b = 2))
< ok 0 0
> query
| SHOW CREATE TABLE chk_compound_or.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000007` CHECK (((`a` = 1) or (`b` = 2)))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_between
< ok 0 0
> exec
| CREATE DATABASE chk_between
< ok 1 0
> exec
| USE chk_between
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a BETWEEN 1 AND 10))
< ok 0 0
> query
| SHOW CREATE TABLE chk_between.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000008` CHECK (((`a` >= 1) and (`a` <= 10)))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_between_forms
< ok 0 0
> exec
| CREATE DATABASE chk_between_forms
< ok 1 0
> exec
| USE chk_between_forms
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, c INT, d DATE, CHECK (a NOT BETWEEN 20 AND 30), CHECK (a BETWEEN b AND c), CHECK (a BETWEEN b AND 10), CHECK (a + 1 BETWEEN 1 AND 2 * 5), CHECK (d BETWEEN '2020-01-01' AND '2021-01-01'), g INT GENERATED ALWAYS AS (a BETWEEN 1 AND 10) VIRTUAL)
< ok 0 0
> query
| SHOW CREATE TABLE chk_between_forms.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` date DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS (((`a` >= 1) and (`a` <= 10))) VIRTUAL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000009` CHECK ((not(((`a` >= 20) and (`a` <= 30))))),
|   CONSTRAINT `t_OBCHECK_0000000000000010` CHECK ((`a` between `b` and `c`)),
|   CONSTRAINT `t_OBCHECK_0000000000000011` CHECK ((`a` between `b` and 10)),
|   CONSTRAINT `t_OBCHECK_0000000000000012` CHECK ((((`a` + 1) >= 1) and ((`a` + 1) <= (2 * 5)))),
|   CONSTRAINT `t_OBCHECK_0000000000000013` CHECK (((`d` >= '2020-01-01') and (`d` <= '2021-01-01')))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_in_list
< ok 0 0
> exec
| CREATE DATABASE chk_in_list
< ok 1 0
> exec
| USE chk_in_list
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a IN (1, 2, 3)))
< ok 0 0
> query
| SHOW CREATE TABLE chk_in_list.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000014` CHECK ((`a` in (1,2,3)))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_string_literal
< ok 0 0
> exec
| CREATE DATABASE chk_string_literal
< ok 1 0
> exec
| USE chk_string_literal
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (s = 'hello')) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE chk_string_literal.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000015` CHECK ((`s` = 'hello'))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_function_call
< ok 0 0
> exec
| CREATE DATABASE chk_function_call
< ok 1 0
> exec
| USE chk_function_call
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (CHAR_LENGTH(s) > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_function_call.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000016` CHECK ((CHAR_LENGTH(`s`) > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_backtick_idents
< ok 0 0
> exec
| CREATE DATABASE chk_backtick_idents
< ok 1 0
> exec
| USE chk_backtick_idents
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (`a` < `b`))
< ok 0 0
> query
| SHOW CREATE TABLE chk_backtick_idents.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000017` CHECK ((`a` < `b`))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_odd_spacing
< ok 0 0
> exec
| CREATE DATABASE chk_odd_spacing
< ok 1 0
> exec
| USE chk_odd_spacing
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (   a    >0   ))
< ok 0 0
> query
| SHOW CREATE TABLE chk_odd_spacing.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000018` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_rich_compound
< ok 0 0
> exec
| CREATE DATABASE chk_rich_compound
< ok 1 0
> exec
| USE chk_rich_compound
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, s VARCHAR(20), CONSTRAINT c1 CHECK (UPPER(s) = 'HI' AND a IN (1,2,3)), CONSTRAINT c2 CHECK (a < b) NOT ENFORCED) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE chk_rich_compound.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `c1` CHECK (((UPPER(`s`) = 'HI') and (`a` in (1,2,3)))),
|   CONSTRAINT `c2` CHECK ((`a` < `b`))  NOT ENFORCED
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS chk_column_level
< ok 0 0
> exec
| CREATE DATABASE chk_column_level
< ok 1 0
> exec
| USE chk_column_level
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE chk_column_level.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `t_OBCHECK_0000000000000019` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
