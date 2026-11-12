-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestCreateTableOptionsOceanBaseOracle default#1
> exec
| DROP DATABASE IF EXISTS obctopt
< ok 0 0
> exec
| CREATE DATABASE obctopt
< ok 1 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) checksum = 1
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) checksum = 2
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) checksum = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_AUTO_RECALC=1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) checksum = 'abc'
< error
| 1064 42000 …near ''abc'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) checksum = DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) table_checksum = 0
< error
| 1064 42000 …near 'table_checksum = 0' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) table_checksum 1
< error
| 1064 42000 …near 'table_checksum 1' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) table_checksum = 1
< error
| 1064 42000 …near 'table_checksum = 1' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) table_checksum
< error
| 1064 42000 …near 'table_checksum' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c int) table_checksum = 'abc'
< error
| 1064 42000 …near 'table_checksum = 'abc'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_AUTO_RECALC=0
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_AUTO_RECALC=1
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) stats_auto_recalc 2
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) stats_auto_recalc = 10
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_AUTO_RECALC=DEFAULT
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) stats_auto_recalc default
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_AUTO_RECALC='abc'
< error
| 1064 42000 …near ''abc'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_AUTO_RECALC='1'
< error
| 1064 42000 …near ''1'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_AUTO_RECALC=-1
< error
| 1064 42000 …near '-1' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_AUTO_RECALC=
< error
| 1064 42000 …near '=' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_BUCKETS=1
< error
| 1064 42000 …near 'STATS_BUCKETS=1' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_TOPN=1
< error
| 1064 42000 …near 'STATS_TOPN=1' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_SAMPLE_RATE=0.1
< error
| 1064 42000 …near 'STATS_SAMPLE_RATE=0.1' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) STATS_COL_CHOICE='all'
< error
| 1064 42000 …near 'STATS_COL_CHOICE='all'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int, b int) STATS_COL_LIST='a,b'
< error
| 1064 42000 …near 'STATS_COL_LIST='a,b'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a serial)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int serial default value)
< error
| 1064 42000 …near 'serial default value)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a bigint serial default value)
< error
| 1064 42000 …near 'serial default value)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int serial default value null)
< error
| 1064 42000 …near 'serial default value null)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int serial default value primary key)
< error
| 1064 42000 …near 'serial default value primary key)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int serial)
< error
| 1064 42000 …near 'serial)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int serial default)
< error
| 1064 42000 …near 'serial default)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a year)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a year(4))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a year zerofill)
< error
| 1064 42000 …near 'zerofill)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a year(4) zerofill)
< error
| 1064 42000 …near 'zerofill)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a year(4) binary)
< error
| 1064 42000 …near 'binary)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a sql_tsi_year)
< error
| 1064 42000 …near 'sql_tsi_year)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a SQL_TSI_YEAR(4))
< error
| 1064 42000 …near 'SQL_TSI_YEAR(4))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a SQL_TSI_YEAR(4) ZEROFILL)
< error
| 1064 42000 …near 'SQL_TSI_YEAR(4) ZEROFILL)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int(1))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a double(10,2))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a year())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a char())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a datetime())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a bit())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a blob())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a binary())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a decimal())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a float())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a decimal(10,))
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a decimal(,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a float(10,))
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int(,))
< error
| 1064 42000 …near ',))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0b01100001, 'b'), c2 set(0b01100001, 'b'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(x'61'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(b'01100001'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0x6162))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 set(0b0110000101100010))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0x61) charset utf8mb4)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0x61, 'b') character set utf8mb4)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum('a', 0x61))
< error
| 1291 HY000 Column 'c1' has duplicated value 'a' in ENUM 
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 set())
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(1))
< error
| 1064 42000 …near '1))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum('a',))
< error
| 1064 42000 …near '))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0x))
< error
| 1064 42000 …near '0x))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0b))
< error
| 1064 42000 …near '0b))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0xZZ))
< error
| 1064 42000 …near '0xZZ))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0x6))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0x616))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0b1))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(0b101))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(x'0b61'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(x'0b'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(x'0B61'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 set(x'0b61'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(b'1'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(b'101'))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(x''))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(b''))
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(x'6'))
< error
| 1064 42000 …near 'x'6'))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(x'616'))
< error
| 1064 42000 …near 'x'616'))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(x'0x61'))
< error
| 1064 42000 …near 'x'0x61'))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (c1 enum(b'0b01'))
< error
| 1064 42000 …near 'b'0b01'))' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE `123` (123a1 int)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (123a1 int)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (1a int)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (123abc int)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (12e int)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (1_2 int)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (123a1 int, 456b2 int)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE 1t (a int)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE 123 (123a1 int)
< error
| 1064 42000 …near '123 (123a1 int)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (1e1 int)
< error
| 1064 42000 …near '1e1 int)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (0x61 int)
< error
| 1064 42000 …near '0x61 int)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| SELECT 1e1 FROM obcto
< error
| 1146 42S02 Table 'obctopt.obcto' doesn't exist
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| SELECT 1E+10 FROM obcto
< error
| 1146 42S02 Table 'obctopt.obcto' doesn't exist
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| SELECT 1.5e3 FROM obcto
< error
| 1146 42S02 Table 'obctopt.obcto' doesn't exist
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| SELECT 1e+ FROM obcto
< error
| 1064 42000 …near 'FROM obcto' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (id INT) SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'
< error
| 1064 42000 …near 'SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (id INT) SECONDARY_ENGINE_ATTRIBUTE=
< error
| 1064 42000 …near 'SECONDARY_ENGINE_ATTRIBUTE=' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (id INT) SECONDARY_ENGINE_ATTRIBUTE
< error
| 1064 42000 …near 'SECONDARY_ENGINE_ATTRIBUTE' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (id INT SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}')
< error
| 1064 42000 …near 'SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}')' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (id INT SECONDARY_ENGINE_ATTRIBUTE=)
< error
| 1064 42000 …near 'SECONDARY_ENGINE_ATTRIBUTE=)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (id INT SECONDARY_ENGINE_ATTRIBUTE)
< error
| 1064 42000 …near 'SECONDARY_ENGINE_ATTRIBUTE)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (id INT) ENGINE_ATTRIBUTE='{"key":"value"}'
< error
| 1064 42000 …near 'ENGINE_ATTRIBUTE='{"key":"value"}'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (id INT ENGINE_ATTRIBUTE='{"key":"value"}')
< error
| 1064 42000 …near 'ENGINE_ATTRIBUTE='{"key":"value"}')' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (id INT ENGINE_ATTRIBUTE)
< error
| 1064 42000 …near 'ENGINE_ATTRIBUTE)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| ALTER TABLE obcto SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'
< error
| 1064 42000 …near 'SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| ALTER TABLE obcto add column c int SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'
< error
| 1064 42000 …near 'SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE GLOBAL TEMPORARY TABLE obcto (a int) ON COMMIT DELETE ROWS
< error
| 1064 42000 …near 'GLOBAL TEMPORARY TABLE obcto (a int) ON COMMIT DELETE ROWS' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TEMPORARY TABLE obcto (a int) ON COMMIT DELETE ROWS
< error
| 1064 42000 …near 'ON COMMIT DELETE ROWS' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TABLE obcto (a int) ON COMMIT DELETE ROWS
< error
| 1064 42000 …near 'ON COMMIT DELETE ROWS' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE GLOBAL TABLE obcto (a int)
< error
| 1064 42000 …near 'GLOBAL TABLE obcto (a int)' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| DROP GLOBAL TEMPORARY TABLE obcto
< error
| 1064 42000 …near 'GLOBAL TEMPORARY TABLE obcto' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE GLOBAL TEMPORARY VIEW obctov AS SELECT 1
< error
| 1064 42000 …near 'GLOBAL TEMPORARY VIEW obctov AS SELECT 1' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE GLOBAL TEMPORARY DATABASE obctod
< error
| 1064 42000 …near 'GLOBAL TEMPORARY DATABASE obctod' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| DROP GLOBAL TEMPORARY DATABASE obctod
< error
| 1064 42000 …near 'GLOBAL TEMPORARY DATABASE obctod' at line 1
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| CREATE TEMPORARY TABLE obcto (a int)
< ok 0 0
> exec
| USE obctopt
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto
< ok 0 0
> exec
| DROP TABLE IF EXISTS obcto_tmp
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`
< ok 0 0
> exec
| DROP TABLE IF EXISTS `123`_tmp
< error
| 1064 42000 …near '_tmp' at line 1
> exec
| DROP TABLE IF EXISTS 1t
< ok 0 0
> exec
| DROP TABLE IF EXISTS 1t_tmp
< ok 0 0
> exec
| DROP TEMPORARY TABLE obcto
< error
| 1051 42S02 Unknown table 'obctopt.obcto'
> exec
| USE obctopt
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obctopt
< ok 0 0
