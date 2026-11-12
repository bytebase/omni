-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestIntroducerWordsMatchOceanBase default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_name_token_forms
< ok 1 0
> exec
| USE omni_name_token_forms
< ok 0 0
> query
| SELECT CHARACTER_SET_NAME FROM information_schema.CHARACTER_SETS UNION SELECT COLLATION_NAME FROM information_schema.COLLATIONS
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| binary
< row
: bytes
| utf8mb4
< row
: bytes
| gbk
< row
: bytes
| utf16
< row
: bytes
| gb18030
< row
: bytes
| latin1
< row
: bytes
| gb2312
< row
: bytes
| gb18030_2022
< row
: bytes
| ascii
< row
: bytes
| tis620
< row
: bytes
| ujis
< row
: bytes
| euckr
< row
: bytes
| eucjpms
< row
: bytes
| cp932
< row
: bytes
| utf16le
< row
: bytes
| sjis
< row
: bytes
| big5
< row
: bytes
| hkscs
< row
: bytes
| hkscs31
< row
: bytes
| dec8
< row
: bytes
| cp850
< row
: bytes
| hp8
< row
: bytes
| macroman
< row
: bytes
| swe7
< row
: bytes
| utf8mb4_general_ci
< row
: bytes
| utf8mb4_bin
< row
: bytes
| gbk_chinese_ci
< row
: bytes
| gbk_bin
< row
: bytes
| utf16_general_ci
< row
: bytes
| utf16_bin
< row
: bytes
| gb18030_chinese_ci
< row
: bytes
| gb18030_bin
< row
: bytes
| latin1_swedish_ci
< row
: bytes
| latin1_german1_ci
< row
: bytes
| latin1_danish_ci
< row
: bytes
| latin1_german2_ci
< row
: bytes
| latin1_general_ci
< row
: bytes
| latin1_general_cs
< row
: bytes
| latin1_spanish_ci
< row
: bytes
| latin1_bin
< row
: bytes
| gb18030_2022_bin
< row
: bytes
| gb18030_2022_chinese_ci
< row
: bytes
| gb18030_2022_chinese_cs
< row
: bytes
| gb18030_2022_radical_ci
< row
: bytes
| gb18030_2022_radical_cs
< row
: bytes
| gb18030_2022_stroke_ci
< row
: bytes
| gb18030_2022_stroke_cs
< row
: bytes
| gb2312_chinese_ci
< row
: bytes
| gb2312_bin
< row
: bytes
| ascii_general_ci
< row
: bytes
| ascii_bin
< row
: bytes
| tis620_thai_ci
< row
: bytes
| tis620_bin
< row
: bytes
| ujis_japanese_ci
< row
: bytes
| ujis_bin
< row
: bytes
| euckr_korean_ci
< row
: bytes
| euckr_bin
< row
: bytes
| eucjpms_japanese_ci
< row
: bytes
| eucjpms_bin
< row
: bytes
| cp932_japanese_ci
< row
: bytes
| cp932_bin
< row
: bytes
| utf16le_general_ci
< row
: bytes
| utf16le_bin
< row
: bytes
| sjis_japanese_ci
< row
: bytes
| sjis_bin
< row
: bytes
| big5_chinese_ci
< row
: bytes
| big5_bin
< row
: bytes
| hkscs_bin
< row
: bytes
| hkscs31_bin
< row
: bytes
| utf16_unicode_ci
< row
: bytes
| utf16_icelandic_ci
< row
: bytes
| utf16_latvian_ci
< row
: bytes
| utf16_romanian_ci
< row
: bytes
| utf16_slovenian_ci
< row
: bytes
| utf16_polish_ci
< row
: bytes
| utf16_estonian_ci
< row
: bytes
| utf16_spanish_ci
< row
: bytes
| utf16_swedish_ci
< row
: bytes
| utf16_turkish_ci
< row
: bytes
| utf16_czech_ci
< row
: bytes
| utf16_danish_ci
< row
: bytes
| utf16_lithuanian_ci
< row
: bytes
| utf16_slovak_ci
< row
: bytes
| utf16_spanish2_ci
< row
: bytes
| utf16_roman_ci
< row
: bytes
| utf16_persian_ci
< row
: bytes
| utf16_esperanto_ci
< row
: bytes
| utf16_hungarian_ci
< row
: bytes
| utf16_sinhala_ci
< row
: bytes
| utf16_german2_ci
< row
: bytes
| utf16_croatian_ci
< row
: bytes
| utf16_unicode_520_ci
< row
: bytes
| utf16_vietnamese_ci
< row
: bytes
| utf8mb4_unicode_ci
< row
: bytes
| utf8mb4_icelandic_ci
< row
: bytes
| utf8mb4_latvian_ci
< row
: bytes
| utf8mb4_romanian_ci
< row
: bytes
| utf8mb4_slovenian_ci
< row
: bytes
| utf8mb4_polish_ci
< row
: bytes
| utf8mb4_estonian_ci
< row
: bytes
| utf8mb4_spanish_ci
< row
: bytes
| utf8mb4_swedish_ci
< row
: bytes
| utf8mb4_turkish_ci
< row
: bytes
| utf8mb4_czech_ci
< row
: bytes
| utf8mb4_danish_ci
< row
: bytes
| utf8mb4_lithuanian_ci
< row
: bytes
| utf8mb4_slovak_ci
< row
: bytes
| utf8mb4_spanish2_ci
< row
: bytes
| utf8mb4_roman_ci
< row
: bytes
| utf8mb4_persian_ci
< row
: bytes
| utf8mb4_esperanto_ci
< row
: bytes
| utf8mb4_hungarian_ci
< row
: bytes
| utf8mb4_sinhala_ci
< row
: bytes
| utf8mb4_german2_ci
< row
: bytes
| utf8mb4_croatian_ci
< row
: bytes
| utf8mb4_unicode_520_ci
< row
: bytes
| utf8mb4_vietnamese_ci
< row
: bytes
| dec8_swedish_ci
< row
: bytes
| dec8_bin
< row
: bytes
| cp850_general_ci
< row
: bytes
| cp850_bin
< row
: bytes
| hp8_english_ci
< row
: bytes
| hp8_bin
< row
: bytes
| macroman_general_ci
< row
: bytes
| macroman_bin
< row
: bytes
| swe7_swedish_ci
< row
: bytes
| swe7_bin
< row
: bytes
| utf8mb4_0900_ai_ci
< row
: bytes
| utf8mb4_de_pb_0900_ai_ci
< row
: bytes
| utf8mb4_is_0900_ai_ci
< row
: bytes
| utf8mb4_lv_0900_ai_ci
< row
: bytes
| utf8mb4_ro_0900_ai_ci
< row
: bytes
| utf8mb4_sl_0900_ai_ci
< row
: bytes
| utf8mb4_pl_0900_ai_ci
< row
: bytes
| utf8mb4_et_0900_ai_ci
< row
: bytes
| utf8mb4_es_0900_ai_ci
< row
: bytes
| utf8mb4_sv_0900_ai_ci
< row
: bytes
| utf8mb4_tr_0900_ai_ci
< row
: bytes
| utf8mb4_cs_0900_ai_ci
< row
: bytes
| utf8mb4_da_0900_ai_ci
< row
: bytes
| utf8mb4_lt_0900_ai_ci
< row
: bytes
| utf8mb4_sk_0900_ai_ci
< row
: bytes
| utf8mb4_es_trad_0900_ai_ci
< row
: bytes
| utf8mb4_la_0900_ai_ci
< row
: bytes
| utf8mb4_eo_0900_ai_ci
< row
: bytes
| utf8mb4_hu_0900_ai_ci
< row
: bytes
| utf8mb4_hr_0900_ai_ci
< row
: bytes
| utf8mb4_vi_0900_ai_ci
< row
: bytes
| utf8mb4_0900_as_cs
< row
: bytes
| utf8mb4_de_pb_0900_as_cs
< row
: bytes
| utf8mb4_is_0900_as_cs
< row
: bytes
| utf8mb4_lv_0900_as_cs
< row
: bytes
| utf8mb4_ro_0900_as_cs
< row
: bytes
| utf8mb4_sl_0900_as_cs
< row
: bytes
| utf8mb4_pl_0900_as_cs
< row
: bytes
| utf8mb4_et_0900_as_cs
< row
: bytes
| utf8mb4_es_0900_as_cs
< row
: bytes
| utf8mb4_sv_0900_as_cs
< row
: bytes
| utf8mb4_tr_0900_as_cs
< row
: bytes
| utf8mb4_cs_0900_as_cs
< row
: bytes
| utf8mb4_da_0900_as_cs
< row
: bytes
| utf8mb4_lt_0900_as_cs
< row
: bytes
| utf8mb4_sk_0900_as_cs
< row
: bytes
| utf8mb4_es_trad_0900_as_cs
< row
: bytes
| utf8mb4_la_0900_as_cs
< row
: bytes
| utf8mb4_eo_0900_as_cs
< row
: bytes
| utf8mb4_hu_0900_as_cs
< row
: bytes
| utf8mb4_hr_0900_as_cs
< row
: bytes
| utf8mb4_vi_0900_as_cs
< row
: bytes
| utf8mb4_ja_0900_as_cs
< row
: bytes
| utf8mb4_ja_0900_as_cs_ks
< row
: bytes
| utf8mb4_0900_as_ci
< row
: bytes
| utf8mb4_ru_0900_ai_ci
< row
: bytes
| utf8mb4_ru_0900_as_cs
< row
: bytes
| utf8mb4_zh_0900_as_cs
< row
: bytes
| utf8mb4_0900_bin
< row
: bytes
| utf8mb4_nb_0900_ai_ci
< row
: bytes
| utf8mb4_nb_0900_as_cs
< row
: bytes
| utf8mb4_nn_0900_ai_ci
< row
: bytes
| utf8mb4_nn_0900_as_cs
< row
: bytes
| utf8mb4_sr_latn_0900_ai_ci
< row
: bytes
| utf8mb4_sr_latn_0900_as_cs
< row
: bytes
| utf8mb4_bs_0900_ai_ci
< row
: bytes
| utf8mb4_bs_0900_as_cs
< row
: bytes
| utf8mb4_bg_0900_ai_ci
< row
: bytes
| utf8mb4_bg_0900_as_cs
< row
: bytes
| utf8mb4_gl_0900_ai_ci
< row
: bytes
| utf8mb4_gl_0900_as_cs
< row
: bytes
| utf8mb4_mn_cyrl_0900_ai_ci
< row
: bytes
| utf8mb4_mn_cyrl_0900_as_cs
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| CREATE TABLE _armscii8 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ascii (b int)
< error
| 1064 42000 …near '_ascii (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._ascii
< error
| 1146 42S02 Table 't._ascii' doesn't exist
> exec ^
: string
| SELECT * FROM _ascii.t
< error
| 1146 42S02 Table '_ascii.t' doesn't exist
> exec ^
: string
| CREATE TABLE _ascii_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ascii_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _big5 (b int)
< error
| 1064 42000 …near '_big5 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._big5
< error
| 1146 42S02 Table 't._big5' doesn't exist
> exec ^
: string
| SELECT * FROM _big5.t
< error
| 1146 42S02 Table '_big5.t' doesn't exist
> exec ^
: string
| CREATE TABLE _big5_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _big5_chinese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _binary (b int)
< error
| 1064 42000 …near '_binary (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._binary
< error
| 1146 42S02 Table 't._binary' doesn't exist
> exec ^
: string
| SELECT * FROM _binary.t
< error
| 1146 42S02 Table '_binary.t' doesn't exist
> exec ^
: string
| CREATE TABLE _cp1250 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1251 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1256 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1257 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp850 (b int)
< error
| 1064 42000 …near '_cp850 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._cp850
< error
| 1146 42S02 Table 't._cp850' doesn't exist
> exec ^
: string
| SELECT * FROM _cp850.t
< error
| 1146 42S02 Table '_cp850.t' doesn't exist
> exec ^
: string
| CREATE TABLE _cp850_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp850_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp852 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp866 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp932 (b int)
< error
| 1064 42000 …near '_cp932 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._cp932
< error
| 1146 42S02 Table 't._cp932' doesn't exist
> exec ^
: string
| SELECT * FROM _cp932.t
< error
| 1146 42S02 Table '_cp932.t' doesn't exist
> exec ^
: string
| CREATE TABLE _cp932_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp932_japanese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _dec8 (b int)
< error
| 1064 42000 …near '_dec8 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._dec8
< error
| 1146 42S02 Table 't._dec8' doesn't exist
> exec ^
: string
| SELECT * FROM _dec8.t
< error
| 1146 42S02 Table '_dec8.t' doesn't exist
> exec ^
: string
| CREATE TABLE _dec8_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _dec8_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _default (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _eucjpms (b int)
< error
| 1064 42000 …near '_eucjpms (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._eucjpms
< error
| 1146 42S02 Table 't._eucjpms' doesn't exist
> exec ^
: string
| SELECT * FROM _eucjpms.t
< error
| 1146 42S02 Table '_eucjpms.t' doesn't exist
> exec ^
: string
| CREATE TABLE _eucjpms_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _eucjpms_japanese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _euckr (b int)
< error
| 1064 42000 …near '_euckr (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._euckr
< error
| 1146 42S02 Table 't._euckr' doesn't exist
> exec ^
: string
| SELECT * FROM _euckr.t
< error
| 1146 42S02 Table '_euckr.t' doesn't exist
> exec ^
: string
| CREATE TABLE _euckr_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _euckr_korean_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _filename (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _foo (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030 (b int)
< error
| 1064 42000 …near '_gb18030 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._gb18030
< error
| 1146 42S02 Table 't._gb18030' doesn't exist
> exec ^
: string
| SELECT * FROM _gb18030.t
< error
| 1146 42S02 Table '_gb18030.t' doesn't exist
> exec ^
: string
| CREATE TABLE _gb18030_2022 (b int)
< error
| 1064 42000 …near '_gb18030_2022 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._gb18030_2022
< error
| 1146 42S02 Table 't._gb18030_2022' doesn't exist
> exec ^
: string
| SELECT * FROM _gb18030_2022.t
< error
| 1146 42S02 Table '_gb18030_2022.t' doesn't exist
> exec ^
: string
| CREATE TABLE _gb18030_2022_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030_2022_chinese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030_2022_chinese_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030_2022_radical_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030_2022_radical_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030_2022_stroke_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030_2022_stroke_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030_chinese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb2312 (b int)
< error
| 1064 42000 …near '_gb2312 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._gb2312
< error
| 1146 42S02 Table 't._gb2312' doesn't exist
> exec ^
: string
| SELECT * FROM _gb2312.t
< error
| 1146 42S02 Table '_gb2312.t' doesn't exist
> exec ^
: string
| CREATE TABLE _gb2312_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb2312_chinese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gbk (b int)
< error
| 1064 42000 …near '_gbk (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._gbk
< error
| 1146 42S02 Table 't._gbk' doesn't exist
> exec ^
: string
| SELECT * FROM _gbk.t
< error
| 1146 42S02 Table '_gbk.t' doesn't exist
> exec ^
: string
| CREATE TABLE _gbk_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gbk_chinese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _geostd8 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _greek (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hebrew (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hkscs (b int)
< error
| 1064 42000 …near '_hkscs (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._hkscs
< error
| 1146 42S02 Table 't._hkscs' doesn't exist
> exec ^
: string
| SELECT * FROM _hkscs.t
< error
| 1146 42S02 Table '_hkscs.t' doesn't exist
> exec ^
: string
| CREATE TABLE _hkscs31 (b int)
< error
| 1064 42000 …near '_hkscs31 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._hkscs31
< error
| 1146 42S02 Table 't._hkscs31' doesn't exist
> exec ^
: string
| SELECT * FROM _hkscs31.t
< error
| 1146 42S02 Table '_hkscs31.t' doesn't exist
> exec ^
: string
| CREATE TABLE _hkscs31_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hkscs_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hp8 (b int)
< error
| 1064 42000 …near '_hp8 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._hp8
< error
| 1146 42S02 Table 't._hp8' doesn't exist
> exec ^
: string
| SELECT * FROM _hp8.t
< error
| 1146 42S02 Table '_hp8.t' doesn't exist
> exec ^
: string
| CREATE TABLE _hp8_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hp8_english_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _keybcs2 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8r (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8u (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin1 (b int)
< error
| 1064 42000 …near '_latin1 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._latin1
< error
| 1146 42S02 Table 't._latin1' doesn't exist
> exec ^
: string
| SELECT * FROM _latin1.t
< error
| 1146 42S02 Table '_latin1.t' doesn't exist
> exec ^
: string
| CREATE TABLE _latin1_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin1_danish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin1_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin1_general_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin1_german1_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin1_german2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin1_spanish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin1_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin2 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin5 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin7 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _macce (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _macroman (b int)
< error
| 1064 42000 …near '_macroman (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._macroman
< error
| 1146 42S02 Table 't._macroman' doesn't exist
> exec ^
: string
| SELECT * FROM _macroman.t
< error
| 1146 42S02 Table '_macroman.t' doesn't exist
> exec ^
: string
| CREATE TABLE _macroman_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _macroman_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _sjis (b int)
< error
| 1064 42000 …near '_sjis (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._sjis
< error
| 1146 42S02 Table 't._sjis' doesn't exist
> exec ^
: string
| SELECT * FROM _sjis.t
< error
| 1146 42S02 Table '_sjis.t' doesn't exist
> exec ^
: string
| CREATE TABLE _sjis_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _sjis_japanese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _swe7 (b int)
< error
| 1064 42000 …near '_swe7 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._swe7
< error
| 1146 42S02 Table 't._swe7' doesn't exist
> exec ^
: string
| SELECT * FROM _swe7.t
< error
| 1146 42S02 Table '_swe7.t' doesn't exist
> exec ^
: string
| CREATE TABLE _swe7_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _swe7_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _tis620 (b int)
< error
| 1064 42000 …near '_tis620 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._tis620
< error
| 1146 42S02 Table 't._tis620' doesn't exist
> exec ^
: string
| SELECT * FROM _tis620.t
< error
| 1146 42S02 Table '_tis620.t' doesn't exist
> exec ^
: string
| CREATE TABLE _tis620_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _tis620_thai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ujis (b int)
< error
| 1064 42000 …near '_ujis (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._ujis
< error
| 1146 42S02 Table 't._ujis' doesn't exist
> exec ^
: string
| SELECT * FROM _ujis.t
< error
| 1146 42S02 Table '_ujis.t' doesn't exist
> exec ^
: string
| CREATE TABLE _ujis_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ujis_japanese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16 (b int)
< error
| 1064 42000 …near '_utf16 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._utf16
< error
| 1146 42S02 Table 't._utf16' doesn't exist
> exec ^
: string
| SELECT * FROM _utf16.t
< error
| 1146 42S02 Table '_utf16.t' doesn't exist
> exec ^
: string
| CREATE TABLE _utf16_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_croatian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_czech_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_danish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_esperanto_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_estonian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_german2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_hungarian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_icelandic_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_latvian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_lithuanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_persian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_polish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_roman_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_romanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_sinhala_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_slovak_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_slovenian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_spanish2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_spanish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_turkish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_unicode_520_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_unicode_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_vietnamese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16le (b int)
< error
| 1064 42000 …near '_utf16le (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._utf16le
< error
| 1146 42S02 Table 't._utf16le' doesn't exist
> exec ^
: string
| SELECT * FROM _utf16le.t
< error
| 1146 42S02 Table '_utf16le.t' doesn't exist
> exec ^
: string
| CREATE TABLE _utf16le_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16le_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8 (b int)
< error
| 1064 42000 …near '_utf8 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._utf8
< error
| 1146 42S02 Table 't._utf8' doesn't exist
> exec ^
: string
| SELECT * FROM _utf8.t
< error
| 1146 42S02 Table '_utf8.t' doesn't exist
> exec ^
: string
| CREATE TABLE _utf8mb3 (b int)
< error
| 1064 42000 …near '_utf8mb3 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._utf8mb3
< error
| 1146 42S02 Table 't._utf8mb3' doesn't exist
> exec ^
: string
| SELECT * FROM _utf8mb3.t
< error
| 1146 42S02 Table '_utf8mb3.t' doesn't exist
> exec ^
: string
| CREATE TABLE _utf8mb4 (b int)
< error
| 1064 42000 …near '_utf8mb4 (b int)' at line 1
> exec ^
: string
| SELECT * FROM t._utf8mb4
< error
| 1146 42S02 Table 't._utf8mb4' doesn't exist
> exec ^
: string
| SELECT * FROM _utf8mb4.t
< error
| 1146 42S02 Table '_utf8mb4.t' doesn't exist
> exec ^
: string
| CREATE TABLE _utf8mb4_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_0900_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_0900_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_bg_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_bg_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_bs_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_bs_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_croatian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_cs_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_cs_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_czech_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_da_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_da_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_danish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_de_pb_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_de_pb_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_eo_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_eo_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_es_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_es_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_es_trad_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_es_trad_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_esperanto_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_estonian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_et_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_et_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_german2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_gl_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_gl_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_hr_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_hr_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_hu_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_hu_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_hungarian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_icelandic_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_is_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_is_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_ja_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_ja_0900_as_cs_ks (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_la_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_la_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_latvian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_lithuanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_lt_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_lt_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_lv_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_lv_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_mn_cyrl_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_mn_cyrl_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_nb_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_nb_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_nn_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_nn_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_persian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_pl_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_pl_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_polish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_ro_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_ro_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_roman_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_romanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_ru_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_ru_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_sinhala_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_sk_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_sk_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_sl_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_sl_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_slovak_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_slovenian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_spanish2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_spanish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_sr_latn_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_sr_latn_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_sv_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_sv_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_tr_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_tr_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_turkish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_unicode_520_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_unicode_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_vi_0900_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_vi_0900_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_vietnamese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_zh_0900_as_cs (b int)
< ok 0 0
