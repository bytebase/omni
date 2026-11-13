-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestTypeBinaryFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| alter table t add c varchar(10) binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) binary character set ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) binary char set ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) binary character set 'ascii'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) binary charset binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) binary collate utf8mb4_bin
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) binary character set ascii collate ascii_general_ci
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) binary not null
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) binary character set ascii not null default 'x'
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c char(5) binary character set ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c character varying(5) binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c tinytext binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c text binary character set ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c text(100) binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c mediumtext binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c longtext binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c long binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c long binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c long varchar binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c enum('x','y') binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c enum('x','y') binary character set ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c set('x','y') binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c nchar(5) binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c nvarchar(5) binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c national char(5) binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c national varchar(5) binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c nchar varchar(5) binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c national char varying(5) binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t2 (c varchar(10) binary character set ascii)
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t modify c char(10) binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t change c c char(10) binary charset ascii
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t add c varchar(10) character set ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c varchar(10) charset ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c varchar(10) char set ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c varchar(10) collate utf8mb4_bin binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c varchar(10) character set ascii collate ascii_bin binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c varchar(10) not null binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c varchar(10) binary binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c varchar(10) binary character set ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c char(5) character set ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c character varying(5) charset ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c text character set ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c longtext charset ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c long charset ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c enum('x','y') character set ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| alter table t add c set('x','y') charset ascii binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t2 (c varchar(10) character set ascii binary)
< error
| 1064 42000 …near 'binary)' at line 1
> exec ^
: string
| alter table t modify c char(10) charset ascii binary
< error
| 1064 42000 …near 'binary' at line 1
