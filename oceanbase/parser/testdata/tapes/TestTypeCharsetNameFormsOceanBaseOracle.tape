-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestTypeCharsetNameFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| alter table t add c char charset 'binary'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char character set 'binary'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char char set 'binary'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset 'utf8mb4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset "utf8mb4"
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset 'UTF8MB4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset utf8mb4
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset `utf8mb4`
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char collate 'utf8mb4_bin'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char collate 'binary'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset 'utf8mb4' collate 'utf8mb4_bin'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset utf8mb4 collate 'utf8mb4_bin'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset 'utf8mb4' collate utf8mb4_bin
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset 'utf8mb4'collate'utf8mb4_bin'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) charset 'utf8mb4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) character set 'utf8mb4' collate 'utf8mb4_bin'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c text charset 'utf8mb4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c text collate 'utf8mb4_bin'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c enum('a') charset 'utf8mb4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c set('a') charset 'utf8mb4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c long charset 'utf8mb4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c long varchar charset 'utf8mb4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t modify c char(10) charset 'utf8mb4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t change c c char(10) charset 'utf8mb4'
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t2 (c char charset 'utf8mb4')
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t2 (c char(1) charset 'utf8mb4' collate 'utf8mb4_bin')
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t2 (c char(1) charset 'utf8mb4' not null)
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset 'utf8mb4' default 'x'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char charset N'utf8mb4'
< error
| 1064 42000 …near 'N'utf8mb4'' at line 1
> exec ^
: string
| alter table t add c char charset _utf8mb4'utf8mb4'
< error
| 1064 42000 …near '_utf8mb4'utf8mb4'' at line 1
> exec ^
: string
| alter table t add c char charset 0x62
< error
| 1064 42000 …near '0x62' at line 1
> exec ^
: string
| alter table t add c char charset x'62'
< error
| 1064 42000 …near 'x'62'' at line 1
> exec ^
: string
| alter table t add c char charset 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| alter table t add c char charset ?
< error
| 1064 42000 …near '?' at line 1
> exec ^
: string
| alter table t add c char charset default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| alter table t add c char charset = 'utf8mb4'
< error
| 1064 42000 …near '= 'utf8mb4'' at line 1
> exec ^
: string
| alter table t add c char charset 'utf8mb4' charset 'utf8mb4'
< error
| 1064 42000 …near 'charset 'utf8mb4'' at line 1
> exec ^
: string
| alter table t add c char collate 'utf8mb4_bin' charset 'utf8mb4'
< error
| 1064 42000 …near 'charset 'utf8mb4'' at line 1
> exec ^
: string
| alter table t add c int charset 'utf8mb4'
< error
| 1064 42000 …near 'charset 'utf8mb4'' at line 1
