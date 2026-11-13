-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestIntroducerWordsMatchMariaDB default#1
> query
| SELECT CHARACTER_SET_NAME FROM information_schema.CHARACTER_SETS UNION SELECT COLLATION_NAME FROM information_schema.COLLATIONS
< result
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| big5
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
| koi8r
< row
: bytes
| latin1
< row
: bytes
| latin2
< row
: bytes
| swe7
< row
: bytes
| ascii
< row
: bytes
| ujis
< row
: bytes
| sjis
< row
: bytes
| hebrew
< row
: bytes
| tis620
< row
: bytes
| euckr
< row
: bytes
| koi8u
< row
: bytes
| gb2312
< row
: bytes
| greek
< row
: bytes
| cp1250
< row
: bytes
| gbk
< row
: bytes
| latin5
< row
: bytes
| armscii8
< row
: bytes
| utf8mb3
< row
: bytes
| ucs2
< row
: bytes
| cp866
< row
: bytes
| keybcs2
< row
: bytes
| macce
< row
: bytes
| macroman
< row
: bytes
| cp852
< row
: bytes
| latin7
< row
: bytes
| utf8mb4
< row
: bytes
| cp1251
< row
: bytes
| utf16
< row
: bytes
| utf16le
< row
: bytes
| cp1256
< row
: bytes
| cp1257
< row
: bytes
| utf32
< row
: bytes
| binary
< row
: bytes
| geostd8
< row
: bytes
| cp932
< row
: bytes
| eucjpms
< row
: bytes
| big5_chinese_ci
< row
: bytes
| big5_bin
< row
: bytes
| big5_chinese_nopad_ci
< row
: bytes
| big5_nopad_bin
< row
: bytes
| dec8_swedish_ci
< row
: bytes
| dec8_bin
< row
: bytes
| dec8_swedish_nopad_ci
< row
: bytes
| dec8_nopad_bin
< row
: bytes
| cp850_general_ci
< row
: bytes
| cp850_bin
< row
: bytes
| cp850_general_nopad_ci
< row
: bytes
| cp850_nopad_bin
< row
: bytes
| hp8_english_ci
< row
: bytes
| hp8_bin
< row
: bytes
| hp8_english_nopad_ci
< row
: bytes
| hp8_nopad_bin
< row
: bytes
| koi8r_general_ci
< row
: bytes
| koi8r_bin
< row
: bytes
| koi8r_general_nopad_ci
< row
: bytes
| koi8r_nopad_bin
< row
: bytes
| latin1_german1_ci
< row
: bytes
| latin1_swedish_ci
< row
: bytes
| latin1_danish_ci
< row
: bytes
| latin1_german2_ci
< row
: bytes
| latin1_bin
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
| latin1_swedish_nopad_ci
< row
: bytes
| latin1_nopad_bin
< row
: bytes
| latin2_czech_cs
< row
: bytes
| latin2_general_ci
< row
: bytes
| latin2_hungarian_ci
< row
: bytes
| latin2_croatian_ci
< row
: bytes
| latin2_bin
< row
: bytes
| latin2_general_nopad_ci
< row
: bytes
| latin2_nopad_bin
< row
: bytes
| swe7_swedish_ci
< row
: bytes
| swe7_bin
< row
: bytes
| swe7_swedish_nopad_ci
< row
: bytes
| swe7_nopad_bin
< row
: bytes
| ascii_general_ci
< row
: bytes
| ascii_bin
< row
: bytes
| ascii_general_nopad_ci
< row
: bytes
| ascii_nopad_bin
< row
: bytes
| ujis_japanese_ci
< row
: bytes
| ujis_bin
< row
: bytes
| ujis_japanese_nopad_ci
< row
: bytes
| ujis_nopad_bin
< row
: bytes
| sjis_japanese_ci
< row
: bytes
| sjis_bin
< row
: bytes
| sjis_japanese_nopad_ci
< row
: bytes
| sjis_nopad_bin
< row
: bytes
| hebrew_general_ci
< row
: bytes
| hebrew_bin
< row
: bytes
| hebrew_general_nopad_ci
< row
: bytes
| hebrew_nopad_bin
< row
: bytes
| tis620_thai_ci
< row
: bytes
| tis620_bin
< row
: bytes
| tis620_thai_nopad_ci
< row
: bytes
| tis620_nopad_bin
< row
: bytes
| euckr_korean_ci
< row
: bytes
| euckr_bin
< row
: bytes
| euckr_korean_nopad_ci
< row
: bytes
| euckr_nopad_bin
< row
: bytes
| koi8u_general_ci
< row
: bytes
| koi8u_bin
< row
: bytes
| koi8u_general_nopad_ci
< row
: bytes
| koi8u_nopad_bin
< row
: bytes
| gb2312_chinese_ci
< row
: bytes
| gb2312_bin
< row
: bytes
| gb2312_chinese_nopad_ci
< row
: bytes
| gb2312_nopad_bin
< row
: bytes
| greek_general_ci
< row
: bytes
| greek_bin
< row
: bytes
| greek_general_nopad_ci
< row
: bytes
| greek_nopad_bin
< row
: bytes
| cp1250_general_ci
< row
: bytes
| cp1250_czech_cs
< row
: bytes
| cp1250_croatian_ci
< row
: bytes
| cp1250_bin
< row
: bytes
| cp1250_polish_ci
< row
: bytes
| cp1250_general_nopad_ci
< row
: bytes
| cp1250_nopad_bin
< row
: bytes
| gbk_chinese_ci
< row
: bytes
| gbk_bin
< row
: bytes
| gbk_chinese_nopad_ci
< row
: bytes
| gbk_nopad_bin
< row
: bytes
| latin5_turkish_ci
< row
: bytes
| latin5_bin
< row
: bytes
| latin5_turkish_nopad_ci
< row
: bytes
| latin5_nopad_bin
< row
: bytes
| armscii8_general_ci
< row
: bytes
| armscii8_bin
< row
: bytes
| armscii8_general_nopad_ci
< row
: bytes
| armscii8_nopad_bin
< row
: bytes
| utf8mb3_general_ci
< row
: bytes
| utf8mb3_bin
< row
: bytes
| utf8mb3_unicode_ci
< row
: bytes
| utf8mb3_icelandic_ci
< row
: bytes
| utf8mb3_latvian_ci
< row
: bytes
| utf8mb3_romanian_ci
< row
: bytes
| utf8mb3_slovenian_ci
< row
: bytes
| utf8mb3_polish_ci
< row
: bytes
| utf8mb3_estonian_ci
< row
: bytes
| utf8mb3_spanish_ci
< row
: bytes
| utf8mb3_swedish_ci
< row
: bytes
| utf8mb3_turkish_ci
< row
: bytes
| utf8mb3_czech_ci
< row
: bytes
| utf8mb3_danish_ci
< row
: bytes
| utf8mb3_lithuanian_ci
< row
: bytes
| utf8mb3_slovak_ci
< row
: bytes
| utf8mb3_spanish2_ci
< row
: bytes
| utf8mb3_roman_ci
< row
: bytes
| utf8mb3_persian_ci
< row
: bytes
| utf8mb3_esperanto_ci
< row
: bytes
| utf8mb3_hungarian_ci
< row
: bytes
| utf8mb3_sinhala_ci
< row
: bytes
| utf8mb3_german2_ci
< row
: bytes
| utf8mb3_croatian_mysql561_ci
< row
: bytes
| utf8mb3_unicode_520_ci
< row
: bytes
| utf8mb3_vietnamese_ci
< row
: bytes
| utf8mb3_general_mysql500_ci
< row
: bytes
| utf8mb3_croatian_ci
< row
: bytes
| utf8mb3_myanmar_ci
< row
: bytes
| utf8mb3_thai_520_w2
< row
: bytes
| utf8mb3_general1400_as_ci
< row
: bytes
| utf8mb3_general_nopad_ci
< row
: bytes
| utf8mb3_nopad_bin
< row
: bytes
| utf8mb3_unicode_nopad_ci
< row
: bytes
| utf8mb3_unicode_520_nopad_ci
< row
: bytes
| ucs2_general_ci
< row
: bytes
| ucs2_bin
< row
: bytes
| ucs2_unicode_ci
< row
: bytes
| ucs2_icelandic_ci
< row
: bytes
| ucs2_latvian_ci
< row
: bytes
| ucs2_romanian_ci
< row
: bytes
| ucs2_slovenian_ci
< row
: bytes
| ucs2_polish_ci
< row
: bytes
| ucs2_estonian_ci
< row
: bytes
| ucs2_spanish_ci
< row
: bytes
| ucs2_swedish_ci
< row
: bytes
| ucs2_turkish_ci
< row
: bytes
| ucs2_czech_ci
< row
: bytes
| ucs2_danish_ci
< row
: bytes
| ucs2_lithuanian_ci
< row
: bytes
| ucs2_slovak_ci
< row
: bytes
| ucs2_spanish2_ci
< row
: bytes
| ucs2_roman_ci
< row
: bytes
| ucs2_persian_ci
< row
: bytes
| ucs2_esperanto_ci
< row
: bytes
| ucs2_hungarian_ci
< row
: bytes
| ucs2_sinhala_ci
< row
: bytes
| ucs2_german2_ci
< row
: bytes
| ucs2_croatian_mysql561_ci
< row
: bytes
| ucs2_unicode_520_ci
< row
: bytes
| ucs2_vietnamese_ci
< row
: bytes
| ucs2_general_mysql500_ci
< row
: bytes
| ucs2_croatian_ci
< row
: bytes
| ucs2_myanmar_ci
< row
: bytes
| ucs2_thai_520_w2
< row
: bytes
| ucs2_general_nopad_ci
< row
: bytes
| ucs2_nopad_bin
< row
: bytes
| ucs2_unicode_nopad_ci
< row
: bytes
| ucs2_unicode_520_nopad_ci
< row
: bytes
| cp866_general_ci
< row
: bytes
| cp866_bin
< row
: bytes
| cp866_general_nopad_ci
< row
: bytes
| cp866_nopad_bin
< row
: bytes
| keybcs2_general_ci
< row
: bytes
| keybcs2_bin
< row
: bytes
| keybcs2_general_nopad_ci
< row
: bytes
| keybcs2_nopad_bin
< row
: bytes
| macce_general_ci
< row
: bytes
| macce_bin
< row
: bytes
| macce_general_nopad_ci
< row
: bytes
| macce_nopad_bin
< row
: bytes
| macroman_general_ci
< row
: bytes
| macroman_bin
< row
: bytes
| macroman_general_nopad_ci
< row
: bytes
| macroman_nopad_bin
< row
: bytes
| cp852_general_ci
< row
: bytes
| cp852_bin
< row
: bytes
| cp852_general_nopad_ci
< row
: bytes
| cp852_nopad_bin
< row
: bytes
| latin7_estonian_cs
< row
: bytes
| latin7_general_ci
< row
: bytes
| latin7_general_cs
< row
: bytes
| latin7_bin
< row
: bytes
| latin7_general_nopad_ci
< row
: bytes
| latin7_nopad_bin
< row
: bytes
| utf8mb4_general_ci
< row
: bytes
| utf8mb4_bin
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
| utf8mb4_croatian_mysql561_ci
< row
: bytes
| utf8mb4_unicode_520_ci
< row
: bytes
| utf8mb4_vietnamese_ci
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
| utf8mb4_0900_as_ci
< row
: bytes
| utf8mb4_0900_bin
< row
: bytes
| utf8mb4_croatian_ci
< row
: bytes
| utf8mb4_myanmar_ci
< row
: bytes
| utf8mb4_thai_520_w2
< row
: bytes
| utf8mb4_general1400_as_ci
< row
: bytes
| utf8mb4_general_nopad_ci
< row
: bytes
| utf8mb4_nopad_bin
< row
: bytes
| utf8mb4_unicode_nopad_ci
< row
: bytes
| utf8mb4_unicode_520_nopad_ci
< row
: bytes
| uca1400_ai_ci
< row
: bytes
| uca1400_ai_cs
< row
: bytes
| uca1400_as_ci
< row
: bytes
| uca1400_as_cs
< row
: bytes
| uca1400_nopad_ai_ci
< row
: bytes
| uca1400_nopad_ai_cs
< row
: bytes
| uca1400_nopad_as_ci
< row
: bytes
| uca1400_nopad_as_cs
< row
: bytes
| uca1400_icelandic_ai_ci
< row
: bytes
| uca1400_icelandic_ai_cs
< row
: bytes
| uca1400_icelandic_as_ci
< row
: bytes
| uca1400_icelandic_as_cs
< row
: bytes
| uca1400_icelandic_nopad_ai_ci
< row
: bytes
| uca1400_icelandic_nopad_ai_cs
< row
: bytes
| uca1400_icelandic_nopad_as_ci
< row
: bytes
| uca1400_icelandic_nopad_as_cs
< row
: bytes
| uca1400_latvian_ai_ci
< row
: bytes
| uca1400_latvian_ai_cs
< row
: bytes
| uca1400_latvian_as_ci
< row
: bytes
| uca1400_latvian_as_cs
< row
: bytes
| uca1400_latvian_nopad_ai_ci
< row
: bytes
| uca1400_latvian_nopad_ai_cs
< row
: bytes
| uca1400_latvian_nopad_as_ci
< row
: bytes
| uca1400_latvian_nopad_as_cs
< row
: bytes
| uca1400_romanian_ai_ci
< row
: bytes
| uca1400_romanian_ai_cs
< row
: bytes
| uca1400_romanian_as_ci
< row
: bytes
| uca1400_romanian_as_cs
< row
: bytes
| uca1400_romanian_nopad_ai_ci
< row
: bytes
| uca1400_romanian_nopad_ai_cs
< row
: bytes
| uca1400_romanian_nopad_as_ci
< row
: bytes
| uca1400_romanian_nopad_as_cs
< row
: bytes
| uca1400_slovenian_ai_ci
< row
: bytes
| uca1400_slovenian_ai_cs
< row
: bytes
| uca1400_slovenian_as_ci
< row
: bytes
| uca1400_slovenian_as_cs
< row
: bytes
| uca1400_slovenian_nopad_ai_ci
< row
: bytes
| uca1400_slovenian_nopad_ai_cs
< row
: bytes
| uca1400_slovenian_nopad_as_ci
< row
: bytes
| uca1400_slovenian_nopad_as_cs
< row
: bytes
| uca1400_polish_ai_ci
< row
: bytes
| uca1400_polish_ai_cs
< row
: bytes
| uca1400_polish_as_ci
< row
: bytes
| uca1400_polish_as_cs
< row
: bytes
| uca1400_polish_nopad_ai_ci
< row
: bytes
| uca1400_polish_nopad_ai_cs
< row
: bytes
| uca1400_polish_nopad_as_ci
< row
: bytes
| uca1400_polish_nopad_as_cs
< row
: bytes
| uca1400_estonian_ai_ci
< row
: bytes
| uca1400_estonian_ai_cs
< row
: bytes
| uca1400_estonian_as_ci
< row
: bytes
| uca1400_estonian_as_cs
< row
: bytes
| uca1400_estonian_nopad_ai_ci
< row
: bytes
| uca1400_estonian_nopad_ai_cs
< row
: bytes
| uca1400_estonian_nopad_as_ci
< row
: bytes
| uca1400_estonian_nopad_as_cs
< row
: bytes
| uca1400_spanish_ai_ci
< row
: bytes
| uca1400_spanish_ai_cs
< row
: bytes
| uca1400_spanish_as_ci
< row
: bytes
| uca1400_spanish_as_cs
< row
: bytes
| uca1400_spanish_nopad_ai_ci
< row
: bytes
| uca1400_spanish_nopad_ai_cs
< row
: bytes
| uca1400_spanish_nopad_as_ci
< row
: bytes
| uca1400_spanish_nopad_as_cs
< row
: bytes
| uca1400_swedish_ai_ci
< row
: bytes
| uca1400_swedish_ai_cs
< row
: bytes
| uca1400_swedish_as_ci
< row
: bytes
| uca1400_swedish_as_cs
< row
: bytes
| uca1400_swedish_nopad_ai_ci
< row
: bytes
| uca1400_swedish_nopad_ai_cs
< row
: bytes
| uca1400_swedish_nopad_as_ci
< row
: bytes
| uca1400_swedish_nopad_as_cs
< row
: bytes
| uca1400_turkish_ai_ci
< row
: bytes
| uca1400_turkish_ai_cs
< row
: bytes
| uca1400_turkish_as_ci
< row
: bytes
| uca1400_turkish_as_cs
< row
: bytes
| uca1400_turkish_nopad_ai_ci
< row
: bytes
| uca1400_turkish_nopad_ai_cs
< row
: bytes
| uca1400_turkish_nopad_as_ci
< row
: bytes
| uca1400_turkish_nopad_as_cs
< row
: bytes
| uca1400_czech_ai_ci
< row
: bytes
| uca1400_czech_ai_cs
< row
: bytes
| uca1400_czech_as_ci
< row
: bytes
| uca1400_czech_as_cs
< row
: bytes
| uca1400_czech_nopad_ai_ci
< row
: bytes
| uca1400_czech_nopad_ai_cs
< row
: bytes
| uca1400_czech_nopad_as_ci
< row
: bytes
| uca1400_czech_nopad_as_cs
< row
: bytes
| uca1400_danish_ai_ci
< row
: bytes
| uca1400_danish_ai_cs
< row
: bytes
| uca1400_danish_as_ci
< row
: bytes
| uca1400_danish_as_cs
< row
: bytes
| uca1400_danish_nopad_ai_ci
< row
: bytes
| uca1400_danish_nopad_ai_cs
< row
: bytes
| uca1400_danish_nopad_as_ci
< row
: bytes
| uca1400_danish_nopad_as_cs
< row
: bytes
| uca1400_lithuanian_ai_ci
< row
: bytes
| uca1400_lithuanian_ai_cs
< row
: bytes
| uca1400_lithuanian_as_ci
< row
: bytes
| uca1400_lithuanian_as_cs
< row
: bytes
| uca1400_lithuanian_nopad_ai_ci
< row
: bytes
| uca1400_lithuanian_nopad_ai_cs
< row
: bytes
| uca1400_lithuanian_nopad_as_ci
< row
: bytes
| uca1400_lithuanian_nopad_as_cs
< row
: bytes
| uca1400_slovak_ai_ci
< row
: bytes
| uca1400_slovak_ai_cs
< row
: bytes
| uca1400_slovak_as_ci
< row
: bytes
| uca1400_slovak_as_cs
< row
: bytes
| uca1400_slovak_nopad_ai_ci
< row
: bytes
| uca1400_slovak_nopad_ai_cs
< row
: bytes
| uca1400_slovak_nopad_as_ci
< row
: bytes
| uca1400_slovak_nopad_as_cs
< row
: bytes
| uca1400_spanish2_ai_ci
< row
: bytes
| uca1400_spanish2_ai_cs
< row
: bytes
| uca1400_spanish2_as_ci
< row
: bytes
| uca1400_spanish2_as_cs
< row
: bytes
| uca1400_spanish2_nopad_ai_ci
< row
: bytes
| uca1400_spanish2_nopad_ai_cs
< row
: bytes
| uca1400_spanish2_nopad_as_ci
< row
: bytes
| uca1400_spanish2_nopad_as_cs
< row
: bytes
| uca1400_roman_ai_ci
< row
: bytes
| uca1400_roman_ai_cs
< row
: bytes
| uca1400_roman_as_ci
< row
: bytes
| uca1400_roman_as_cs
< row
: bytes
| uca1400_roman_nopad_ai_ci
< row
: bytes
| uca1400_roman_nopad_ai_cs
< row
: bytes
| uca1400_roman_nopad_as_ci
< row
: bytes
| uca1400_roman_nopad_as_cs
< row
: bytes
| uca1400_persian_ai_ci
< row
: bytes
| uca1400_persian_ai_cs
< row
: bytes
| uca1400_persian_as_ci
< row
: bytes
| uca1400_persian_as_cs
< row
: bytes
| uca1400_persian_nopad_ai_ci
< row
: bytes
| uca1400_persian_nopad_ai_cs
< row
: bytes
| uca1400_persian_nopad_as_ci
< row
: bytes
| uca1400_persian_nopad_as_cs
< row
: bytes
| uca1400_esperanto_ai_ci
< row
: bytes
| uca1400_esperanto_ai_cs
< row
: bytes
| uca1400_esperanto_as_ci
< row
: bytes
| uca1400_esperanto_as_cs
< row
: bytes
| uca1400_esperanto_nopad_ai_ci
< row
: bytes
| uca1400_esperanto_nopad_ai_cs
< row
: bytes
| uca1400_esperanto_nopad_as_ci
< row
: bytes
| uca1400_esperanto_nopad_as_cs
< row
: bytes
| uca1400_hungarian_ai_ci
< row
: bytes
| uca1400_hungarian_ai_cs
< row
: bytes
| uca1400_hungarian_as_ci
< row
: bytes
| uca1400_hungarian_as_cs
< row
: bytes
| uca1400_hungarian_nopad_ai_ci
< row
: bytes
| uca1400_hungarian_nopad_ai_cs
< row
: bytes
| uca1400_hungarian_nopad_as_ci
< row
: bytes
| uca1400_hungarian_nopad_as_cs
< row
: bytes
| uca1400_sinhala_ai_ci
< row
: bytes
| uca1400_sinhala_ai_cs
< row
: bytes
| uca1400_sinhala_as_ci
< row
: bytes
| uca1400_sinhala_as_cs
< row
: bytes
| uca1400_sinhala_nopad_ai_ci
< row
: bytes
| uca1400_sinhala_nopad_ai_cs
< row
: bytes
| uca1400_sinhala_nopad_as_ci
< row
: bytes
| uca1400_sinhala_nopad_as_cs
< row
: bytes
| uca1400_german2_ai_ci
< row
: bytes
| uca1400_german2_ai_cs
< row
: bytes
| uca1400_german2_as_ci
< row
: bytes
| uca1400_german2_as_cs
< row
: bytes
| uca1400_german2_nopad_ai_ci
< row
: bytes
| uca1400_german2_nopad_ai_cs
< row
: bytes
| uca1400_german2_nopad_as_ci
< row
: bytes
| uca1400_german2_nopad_as_cs
< row
: bytes
| uca1400_vietnamese_ai_ci
< row
: bytes
| uca1400_vietnamese_ai_cs
< row
: bytes
| uca1400_vietnamese_as_ci
< row
: bytes
| uca1400_vietnamese_as_cs
< row
: bytes
| uca1400_vietnamese_nopad_ai_ci
< row
: bytes
| uca1400_vietnamese_nopad_ai_cs
< row
: bytes
| uca1400_vietnamese_nopad_as_ci
< row
: bytes
| uca1400_vietnamese_nopad_as_cs
< row
: bytes
| uca1400_croatian_ai_ci
< row
: bytes
| uca1400_croatian_ai_cs
< row
: bytes
| uca1400_croatian_as_ci
< row
: bytes
| uca1400_croatian_as_cs
< row
: bytes
| uca1400_croatian_nopad_ai_ci
< row
: bytes
| uca1400_croatian_nopad_ai_cs
< row
: bytes
| uca1400_croatian_nopad_as_ci
< row
: bytes
| uca1400_croatian_nopad_as_cs
< row
: bytes
| cp1251_bulgarian_ci
< row
: bytes
| cp1251_ukrainian_ci
< row
: bytes
| cp1251_bin
< row
: bytes
| cp1251_general_ci
< row
: bytes
| cp1251_general_cs
< row
: bytes
| cp1251_nopad_bin
< row
: bytes
| cp1251_general_nopad_ci
< row
: bytes
| utf16_general_ci
< row
: bytes
| utf16_bin
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
| utf16_croatian_mysql561_ci
< row
: bytes
| utf16_unicode_520_ci
< row
: bytes
| utf16_vietnamese_ci
< row
: bytes
| utf16_croatian_ci
< row
: bytes
| utf16_myanmar_ci
< row
: bytes
| utf16_thai_520_w2
< row
: bytes
| utf16_general_nopad_ci
< row
: bytes
| utf16_nopad_bin
< row
: bytes
| utf16_unicode_nopad_ci
< row
: bytes
| utf16_unicode_520_nopad_ci
< row
: bytes
| utf16le_general_ci
< row
: bytes
| utf16le_bin
< row
: bytes
| utf16le_general_nopad_ci
< row
: bytes
| utf16le_nopad_bin
< row
: bytes
| cp1256_general_ci
< row
: bytes
| cp1256_bin
< row
: bytes
| cp1256_general_nopad_ci
< row
: bytes
| cp1256_nopad_bin
< row
: bytes
| cp1257_lithuanian_ci
< row
: bytes
| cp1257_bin
< row
: bytes
| cp1257_general_ci
< row
: bytes
| cp1257_nopad_bin
< row
: bytes
| cp1257_general_nopad_ci
< row
: bytes
| utf32_general_ci
< row
: bytes
| utf32_bin
< row
: bytes
| utf32_unicode_ci
< row
: bytes
| utf32_icelandic_ci
< row
: bytes
| utf32_latvian_ci
< row
: bytes
| utf32_romanian_ci
< row
: bytes
| utf32_slovenian_ci
< row
: bytes
| utf32_polish_ci
< row
: bytes
| utf32_estonian_ci
< row
: bytes
| utf32_spanish_ci
< row
: bytes
| utf32_swedish_ci
< row
: bytes
| utf32_turkish_ci
< row
: bytes
| utf32_czech_ci
< row
: bytes
| utf32_danish_ci
< row
: bytes
| utf32_lithuanian_ci
< row
: bytes
| utf32_slovak_ci
< row
: bytes
| utf32_spanish2_ci
< row
: bytes
| utf32_roman_ci
< row
: bytes
| utf32_persian_ci
< row
: bytes
| utf32_esperanto_ci
< row
: bytes
| utf32_hungarian_ci
< row
: bytes
| utf32_sinhala_ci
< row
: bytes
| utf32_german2_ci
< row
: bytes
| utf32_croatian_mysql561_ci
< row
: bytes
| utf32_unicode_520_ci
< row
: bytes
| utf32_vietnamese_ci
< row
: bytes
| utf32_croatian_ci
< row
: bytes
| utf32_myanmar_ci
< row
: bytes
| utf32_thai_520_w2
< row
: bytes
| utf32_general_nopad_ci
< row
: bytes
| utf32_nopad_bin
< row
: bytes
| utf32_unicode_nopad_ci
< row
: bytes
| utf32_unicode_520_nopad_ci
< row
: bytes
| geostd8_general_ci
< row
: bytes
| geostd8_bin
< row
: bytes
| geostd8_general_nopad_ci
< row
: bytes
| geostd8_nopad_bin
< row
: bytes
| cp932_japanese_ci
< row
: bytes
| cp932_bin
< row
: bytes
| cp932_japanese_nopad_ci
< row
: bytes
| cp932_nopad_bin
< row
: bytes
| eucjpms_japanese_ci
< row
: bytes
| eucjpms_bin
< row
: bytes
| eucjpms_japanese_nopad_ci
< row
: bytes
| eucjpms_nopad_bin
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE _armscii8 (b int)
< error
| 1064 42000 …near '_armscii8 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._armscii8
< error
| 1146 42S02 Table 'd._armscii8' doesn't exist
> exec ^
: string
| SELECT * FROM d._armscii8
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _armscii8_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _armscii8_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _armscii8_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _armscii8_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ascii (b int)
< error
| 1064 42000 …near '_ascii (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._ascii
< error
| 1146 42S02 Table 'd._ascii' doesn't exist
> exec ^
: string
| SELECT * FROM d._ascii
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _ascii_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ascii_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _big5 (b int)
< error
| 1064 42000 …near '_big5 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._big5
< error
| 1146 42S02 Table 'd._big5' doesn't exist
> exec ^
: string
| SELECT * FROM d._big5
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _big5_chinese_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _big5_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _binary (b int)
< error
| 1064 42000 …near '_binary (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._binary
< error
| 1146 42S02 Table 'd._binary' doesn't exist
> exec ^
: string
| SELECT * FROM d._binary
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _cp1250 (b int)
< error
| 1064 42000 …near '_cp1250 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._cp1250
< error
| 1146 42S02 Table 'd._cp1250' doesn't exist
> exec ^
: string
| SELECT * FROM d._cp1250
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _cp1250_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1250_croatian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1250_czech_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1250_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1250_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1250_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1250_polish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1251 (b int)
< error
| 1064 42000 …near '_cp1251 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._cp1251
< error
| 1146 42S02 Table 'd._cp1251' doesn't exist
> exec ^
: string
| SELECT * FROM d._cp1251
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _cp1251_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1251_bulgarian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1251_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1251_general_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1251_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1251_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1251_ukrainian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1256 (b int)
< error
| 1064 42000 …near '_cp1256 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._cp1256
< error
| 1146 42S02 Table 'd._cp1256' doesn't exist
> exec ^
: string
| SELECT * FROM d._cp1256
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _cp1256_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1256_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1256_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1256_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1257 (b int)
< error
| 1064 42000 …near '_cp1257 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._cp1257
< error
| 1146 42S02 Table 'd._cp1257' doesn't exist
> exec ^
: string
| SELECT * FROM d._cp1257
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _cp1257_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1257_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1257_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1257_lithuanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp1257_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp850 (b int)
< error
| 1064 42000 …near '_cp850 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._cp850
< error
| 1146 42S02 Table 'd._cp850' doesn't exist
> exec ^
: string
| SELECT * FROM d._cp850
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _cp850_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp850_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp852 (b int)
< error
| 1064 42000 …near '_cp852 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._cp852
< error
| 1146 42S02 Table 'd._cp852' doesn't exist
> exec ^
: string
| SELECT * FROM d._cp852
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _cp852_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp852_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp852_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp852_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp866 (b int)
< error
| 1064 42000 …near '_cp866 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._cp866
< error
| 1146 42S02 Table 'd._cp866' doesn't exist
> exec ^
: string
| SELECT * FROM d._cp866
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _cp866_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp866_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp866_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp866_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp932 (b int)
< error
| 1064 42000 …near '_cp932 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._cp932
< error
| 1146 42S02 Table 'd._cp932' doesn't exist
> exec ^
: string
| SELECT * FROM d._cp932
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _cp932_japanese_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _cp932_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _dec8 (b int)
< error
| 1064 42000 …near '_dec8 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._dec8
< error
| 1146 42S02 Table 'd._dec8' doesn't exist
> exec ^
: string
| SELECT * FROM d._dec8
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _dec8_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _dec8_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _dec8_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _dec8_swedish_nopad_ci (b int)
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
| SELECT * FROM d._eucjpms
< error
| 1146 42S02 Table 'd._eucjpms' doesn't exist
> exec ^
: string
| SELECT * FROM d._eucjpms
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _eucjpms_japanese_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _eucjpms_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _euckr (b int)
< error
| 1064 42000 …near '_euckr (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._euckr
< error
| 1146 42S02 Table 'd._euckr' doesn't exist
> exec ^
: string
| SELECT * FROM d._euckr
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _euckr_korean_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _euckr_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _filename (b int)
< error
| 1064 42000 …near '_filename (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._filename
< error
| 1146 42S02 Table 'd._filename' doesn't exist
> exec ^
: string
| SELECT * FROM d._filename
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _foo (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb18030 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb2312 (b int)
< error
| 1064 42000 …near '_gb2312 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._gb2312
< error
| 1146 42S02 Table 'd._gb2312' doesn't exist
> exec ^
: string
| SELECT * FROM d._gb2312
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _gb2312_chinese_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gb2312_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gbk (b int)
< error
| 1064 42000 …near '_gbk (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._gbk
< error
| 1146 42S02 Table 'd._gbk' doesn't exist
> exec ^
: string
| SELECT * FROM d._gbk
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _gbk_chinese_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _gbk_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _geostd8 (b int)
< error
| 1064 42000 …near '_geostd8 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._geostd8
< error
| 1146 42S02 Table 'd._geostd8' doesn't exist
> exec ^
: string
| SELECT * FROM d._geostd8
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _geostd8_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _geostd8_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _geostd8_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _geostd8_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _greek (b int)
< error
| 1064 42000 …near '_greek (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._greek
< error
| 1146 42S02 Table 'd._greek' doesn't exist
> exec ^
: string
| SELECT * FROM d._greek
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _greek_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _greek_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _greek_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _greek_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hebrew (b int)
< error
| 1064 42000 …near '_hebrew (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._hebrew
< error
| 1146 42S02 Table 'd._hebrew' doesn't exist
> exec ^
: string
| SELECT * FROM d._hebrew
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _hebrew_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hebrew_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hebrew_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hebrew_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hp8 (b int)
< error
| 1064 42000 …near '_hp8 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._hp8
< error
| 1146 42S02 Table 'd._hp8' doesn't exist
> exec ^
: string
| SELECT * FROM d._hp8
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _hp8_english_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _hp8_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _keybcs2 (b int)
< error
| 1064 42000 …near '_keybcs2 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._keybcs2
< error
| 1146 42S02 Table 'd._keybcs2' doesn't exist
> exec ^
: string
| SELECT * FROM d._keybcs2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _keybcs2_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _keybcs2_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _keybcs2_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _keybcs2_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8r (b int)
< error
| 1064 42000 …near '_koi8r (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._koi8r
< error
| 1146 42S02 Table 'd._koi8r' doesn't exist
> exec ^
: string
| SELECT * FROM d._koi8r
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _koi8r_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8r_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8r_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8r_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8u (b int)
< error
| 1064 42000 …near '_koi8u (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._koi8u
< error
| 1146 42S02 Table 'd._koi8u' doesn't exist
> exec ^
: string
| SELECT * FROM d._koi8u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _koi8u_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8u_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8u_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _koi8u_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin1 (b int)
< error
| 1064 42000 …near '_latin1 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._latin1
< error
| 1146 42S02 Table 'd._latin1' doesn't exist
> exec ^
: string
| SELECT * FROM d._latin1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _latin1_nopad_bin (b int)
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
| CREATE TABLE _latin1_swedish_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin2 (b int)
< error
| 1064 42000 …near '_latin2 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._latin2
< error
| 1146 42S02 Table 'd._latin2' doesn't exist
> exec ^
: string
| SELECT * FROM d._latin2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _latin2_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin2_croatian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin2_czech_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin2_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin2_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin2_hungarian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin2_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin5 (b int)
< error
| 1064 42000 …near '_latin5 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._latin5
< error
| 1146 42S02 Table 'd._latin5' doesn't exist
> exec ^
: string
| SELECT * FROM d._latin5
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _latin5_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin5_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin5_turkish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin5_turkish_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin7 (b int)
< error
| 1064 42000 …near '_latin7 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._latin7
< error
| 1146 42S02 Table 'd._latin7' doesn't exist
> exec ^
: string
| SELECT * FROM d._latin7
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _latin7_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin7_estonian_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin7_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin7_general_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin7_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _latin7_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _macce (b int)
< error
| 1064 42000 …near '_macce (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._macce
< error
| 1146 42S02 Table 'd._macce' doesn't exist
> exec ^
: string
| SELECT * FROM d._macce
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _macce_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _macce_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _macce_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _macce_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _macroman (b int)
< error
| 1064 42000 …near '_macroman (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._macroman
< error
| 1146 42S02 Table 'd._macroman' doesn't exist
> exec ^
: string
| SELECT * FROM d._macroman
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _macroman_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _macroman_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _sjis (b int)
< error
| 1064 42000 …near '_sjis (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._sjis
< error
| 1146 42S02 Table 'd._sjis' doesn't exist
> exec ^
: string
| SELECT * FROM d._sjis
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _sjis_japanese_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _sjis_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _swe7 (b int)
< error
| 1064 42000 …near '_swe7 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._swe7
< error
| 1146 42S02 Table 'd._swe7' doesn't exist
> exec ^
: string
| SELECT * FROM d._swe7
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _swe7_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _swe7_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _swe7_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _swe7_swedish_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _tis620 (b int)
< error
| 1064 42000 …near '_tis620 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._tis620
< error
| 1146 42S02 Table 'd._tis620' doesn't exist
> exec ^
: string
| SELECT * FROM d._tis620
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _tis620_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _tis620_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _tis620_thai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _tis620_thai_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_croatian_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_croatian_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_croatian_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_croatian_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_croatian_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_croatian_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_croatian_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_croatian_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_czech_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_czech_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_czech_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_czech_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_czech_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_czech_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_czech_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_czech_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_danish_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_danish_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_danish_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_danish_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_danish_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_danish_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_danish_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_danish_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_esperanto_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_esperanto_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_esperanto_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_esperanto_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_esperanto_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_esperanto_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_esperanto_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_esperanto_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_estonian_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_estonian_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_estonian_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_estonian_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_estonian_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_estonian_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_estonian_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_estonian_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_german2_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_german2_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_german2_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_german2_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_german2_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_german2_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_german2_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_german2_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_hungarian_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_hungarian_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_hungarian_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_hungarian_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_hungarian_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_hungarian_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_hungarian_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_hungarian_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_icelandic_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_icelandic_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_icelandic_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_icelandic_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_icelandic_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_icelandic_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_icelandic_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_icelandic_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_latvian_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_latvian_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_latvian_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_latvian_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_latvian_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_latvian_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_latvian_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_latvian_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_lithuanian_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_lithuanian_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_lithuanian_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_lithuanian_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_lithuanian_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_lithuanian_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_lithuanian_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_lithuanian_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_persian_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_persian_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_persian_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_persian_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_persian_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_persian_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_persian_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_persian_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_polish_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_polish_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_polish_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_polish_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_polish_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_polish_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_polish_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_polish_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_roman_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_roman_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_roman_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_roman_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_roman_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_roman_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_roman_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_roman_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_romanian_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_romanian_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_romanian_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_romanian_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_romanian_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_romanian_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_romanian_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_romanian_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_sinhala_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_sinhala_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_sinhala_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_sinhala_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_sinhala_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_sinhala_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_sinhala_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_sinhala_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovak_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovak_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovak_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovak_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovak_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovak_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovak_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovak_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovenian_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovenian_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovenian_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovenian_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovenian_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovenian_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovenian_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_slovenian_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish2_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish2_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish2_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish2_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish2_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish2_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish2_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish2_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_spanish_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_swedish_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_swedish_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_swedish_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_swedish_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_swedish_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_swedish_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_swedish_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_swedish_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_turkish_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_turkish_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_turkish_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_turkish_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_turkish_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_turkish_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_turkish_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_turkish_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_vietnamese_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_vietnamese_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_vietnamese_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_vietnamese_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_vietnamese_nopad_ai_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_vietnamese_nopad_ai_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_vietnamese_nopad_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _uca1400_vietnamese_nopad_as_cs (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2 (b int)
< error
| 1064 42000 …near '_ucs2 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._ucs2
< error
| 1146 42S02 Table 'd._ucs2' doesn't exist
> exec ^
: string
| SELECT * FROM d._ucs2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _ucs2_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_croatian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_croatian_mysql561_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_czech_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_danish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_esperanto_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_estonian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_general_mysql500_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_german2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_hungarian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_icelandic_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_latvian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_lithuanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_myanmar_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_persian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_polish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_roman_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_romanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_sinhala_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_slovak_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_slovenian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_spanish2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_spanish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_thai_520_w2 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_turkish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_unicode_520_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_unicode_520_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_unicode_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_unicode_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ucs2_vietnamese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ujis (b int)
< error
| 1064 42000 …near '_ujis (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._ujis
< error
| 1146 42S02 Table 'd._ujis' doesn't exist
> exec ^
: string
| SELECT * FROM d._ujis
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _ujis_japanese_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _ujis_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16 (b int)
< error
| 1064 42000 …near '_utf16 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._utf16
< error
| 1146 42S02 Table 'd._utf16' doesn't exist
> exec ^
: string
| SELECT * FROM d._utf16
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _utf16_croatian_mysql561_ci (b int)
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
| CREATE TABLE _utf16_general_nopad_ci (b int)
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
| CREATE TABLE _utf16_myanmar_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_nopad_bin (b int)
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
| CREATE TABLE _utf16_thai_520_w2 (b int)
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
| CREATE TABLE _utf16_unicode_520_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_unicode_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16_unicode_nopad_ci (b int)
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
| SELECT * FROM d._utf16le
< error
| 1146 42S02 Table 'd._utf16le' doesn't exist
> exec ^
: string
| SELECT * FROM d._utf16le
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _utf16le_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf16le_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32 (b int)
< error
| 1064 42000 …near '_utf32 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._utf32
< error
| 1146 42S02 Table 'd._utf32' doesn't exist
> exec ^
: string
| SELECT * FROM d._utf32
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _utf32_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_croatian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_croatian_mysql561_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_czech_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_danish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_esperanto_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_estonian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_german2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_hungarian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_icelandic_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_latvian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_lithuanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_myanmar_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_persian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_polish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_roman_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_romanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_sinhala_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_slovak_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_slovenian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_spanish2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_spanish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_thai_520_w2 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_turkish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_unicode_520_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_unicode_520_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_unicode_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_unicode_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf32_vietnamese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8 (b int)
< error
| 1064 42000 …near '_utf8 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._utf8
< error
| 1146 42S02 Table 'd._utf8' doesn't exist
> exec ^
: string
| SELECT * FROM d._utf8
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _utf8mb3 (b int)
< error
| 1064 42000 …near '_utf8mb3 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._utf8mb3
< error
| 1146 42S02 Table 'd._utf8mb3' doesn't exist
> exec ^
: string
| SELECT * FROM d._utf8mb3
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE _utf8mb3_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_croatian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_croatian_mysql561_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_czech_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_danish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_esperanto_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_estonian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_general1400_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_general_mysql500_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_german2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_hungarian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_icelandic_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_latvian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_lithuanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_myanmar_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_nopad_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_persian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_polish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_roman_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_romanian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_sinhala_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_slovak_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_slovenian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_spanish2_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_spanish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_swedish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_thai_520_w2 (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_turkish_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_unicode_520_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_unicode_520_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_unicode_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_unicode_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb3_vietnamese_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4 (b int)
< error
| 1064 42000 …near '_utf8mb4 (b int)' at line 1
> exec ^
: string
| SELECT * FROM d._utf8mb4
< error
| 1146 42S02 Table 'd._utf8mb4' doesn't exist
> exec ^
: string
| SELECT * FROM d._utf8mb4
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| CREATE TABLE _utf8mb4_bin (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_croatian_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_croatian_mysql561_ci (b int)
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
| CREATE TABLE _utf8mb4_general1400_as_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_general_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_general_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_german2_ci (b int)
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
| CREATE TABLE _utf8mb4_myanmar_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_nopad_bin (b int)
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
| CREATE TABLE _utf8mb4_thai_520_w2 (b int)
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
| CREATE TABLE _utf8mb4_unicode_520_nopad_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_unicode_ci (b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8mb4_unicode_nopad_ci (b int)
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
