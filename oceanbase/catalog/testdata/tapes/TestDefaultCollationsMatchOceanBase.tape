-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDefaultCollationsMatchOceanBase default#1
> query
| SHOW CHARACTER SET
< result
< column "Charset" "VARCHAR"
< column "Description" "VARCHAR"
< column "Default collation" "VARCHAR"
< column "Maxlen" "BIGINT"
< row
: bytes
| binary
: bytes
| Binary pseudo charset
: bytes
| binary
: int64 1
< row
: bytes
| utf8mb4
: bytes
| UTF-8 Unicode
: bytes
| utf8mb4_general_ci
: int64 4
< row
: bytes
| gbk
: bytes
| GBK charset
: bytes
| gbk_chinese_ci
: int64 2
< row
: bytes
| utf16
: bytes
| UTF-16 Unicode
: bytes
| utf16_general_ci
: int64 4
< row
: bytes
| gb18030
: bytes
| GB18030 charset
: bytes
| gb18030_chinese_ci
: int64 4
< row
: bytes
| latin1
: bytes
| cp1252 West European
: bytes
| latin1_swedish_ci
: int64 1
< row
: bytes
| gb2312
: bytes
| GB2312 Simplified Chinese
: bytes
| gb2312_chinese_ci
: int64 2
< row
: bytes
| gb18030_2022
: bytes
| GB18030-2022 charset
: bytes
| gb18030_2022_chinese_ci
: int64 4
< row
: bytes
| ascii
: bytes
| US ASCII
: bytes
| ascii_general_ci
: int64 1
< row
: bytes
| tis620
: bytes
| TIS620 Thai
: bytes
| tis620_thai_ci
: int64 1
< row
: bytes
| ujis
: bytes
| EUC-JP Japanese
: bytes
| ujis_japanese_ci
: int64 3
< row
: bytes
| euckr
: bytes
| EUC-KR Korean
: bytes
| euckr_korean_ci
: int64 2
< row
: bytes
| eucjpms
: bytes
| UJIS for Windows Japanese
: bytes
| eucjpms_japanese_ci
: int64 3
< row
: bytes
| cp932
: bytes
| SJIS for Windows Japanese
: bytes
| cp932_japanese_ci
: int64 2
< row
: bytes
| utf16le
: bytes
| UTF-16LE Unicode
: bytes
| utf16le_general_ci
: int64 4
< row
: bytes
| sjis
: bytes
| SJIS
: bytes
| sjis_japanese_ci
: int64 2
< row
: bytes
| big5
: bytes
| BIG5
: bytes
| big5_chinese_ci
: int64 2
< row
: bytes
| hkscs
: bytes
| HKSCS
: bytes
| hkscs_bin
: int64 2
< row
: bytes
| hkscs31
: bytes
| HKSCS-ISO UNICODE 31
: bytes
| hkscs31_bin
: int64 2
< row
: bytes
| dec8
: bytes
| DEC West European
: bytes
| dec8_swedish_ci
: int64 1
< row
: bytes
| cp850
: bytes
| DOS West European
: bytes
| cp850_general_ci
: int64 1
< row
: bytes
| hp8
: bytes
| HP West European
: bytes
| hp8_english_ci
: int64 1
< row
: bytes
| macroman
: bytes
| Mac West European
: bytes
| macroman_general_ci
: int64 1
< row
: bytes
| swe7
: bytes
| 7bit West European
: bytes
| swe7_swedish_ci
: int64 1
