-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestKillFunctionNameSweepMatchesTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| KILL _armscii8()
< error
| 1064 42000 …line 1 column 14 near "_armscii8()" 
> exec ^
: string
| KILL _armscii8 (1)
< error
| 1064 42000 …line 1 column 14 near "_armscii8 (1)" 
> exec ^
: string
| SET CONFIG tikv _armscii8 = 1
< error
| 1064 42000 …line 1 column 25 near "_armscii8 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_armscii8 = 1
< error
| 1064 42000 …line 1 column 27 near "_armscii8 = 1" 
> exec ^
: string
| KILL _ascii()
< error
| 1064 42000 …line 1 column 11 near "_ascii()" 
> exec ^
: string
| KILL _ascii (1)
< error
| 1064 42000 …line 1 column 11 near "_ascii (1)" 
> exec ^
: string
| SET CONFIG tikv _ascii = 1
< error
| 1064 42000 …line 1 column 22 near "_ascii = 1" 
> exec ^
: string
| SET CONFIG tikv a-_ascii = 1
< error
| 1064 42000 …line 1 column 24 near "_ascii = 1" 
> exec ^
: string
| KILL _big5()
< error
| 1064 42000 …line 1 column 10 near "_big5()" 
> exec ^
: string
| KILL _big5 (1)
< error
| 1064 42000 …line 1 column 10 near "_big5 (1)" 
> exec ^
: string
| SET CONFIG tikv _big5 = 1
< error
| 1064 42000 …line 1 column 21 near "_big5 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_big5 = 1
< error
| 1064 42000 …line 1 column 23 near "_big5 = 1" 
> exec ^
: string
| KILL _binary()
< error
| 1064 42000 …line 1 column 12 near "_binary()" 
> exec ^
: string
| KILL _binary (1)
< error
| 1064 42000 …line 1 column 12 near "_binary (1)" 
> exec ^
: string
| SET CONFIG tikv _binary = 1
< error
| 1064 42000 …line 1 column 23 near "_binary = 1" 
> exec ^
: string
| SET CONFIG tikv a-_binary = 1
< error
| 1064 42000 …line 1 column 25 near "_binary = 1" 
> exec ^
: string
| KILL _cp1250()
< error
| 1064 42000 …line 1 column 12 near "_cp1250()" 
> exec ^
: string
| KILL _cp1250 (1)
< error
| 1064 42000 …line 1 column 12 near "_cp1250 (1)" 
> exec ^
: string
| SET CONFIG tikv _cp1250 = 1
< error
| 1064 42000 …line 1 column 23 near "_cp1250 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_cp1250 = 1
< error
| 1064 42000 …line 1 column 25 near "_cp1250 = 1" 
> exec ^
: string
| KILL _cp1251()
< error
| 1064 42000 …line 1 column 12 near "_cp1251()" 
> exec ^
: string
| KILL _cp1251 (1)
< error
| 1064 42000 …line 1 column 12 near "_cp1251 (1)" 
> exec ^
: string
| SET CONFIG tikv _cp1251 = 1
< error
| 1064 42000 …line 1 column 23 near "_cp1251 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_cp1251 = 1
< error
| 1064 42000 …line 1 column 25 near "_cp1251 = 1" 
> exec ^
: string
| KILL _cp1256()
< error
| 1064 42000 …line 1 column 12 near "_cp1256()" 
> exec ^
: string
| KILL _cp1256 (1)
< error
| 1064 42000 …line 1 column 12 near "_cp1256 (1)" 
> exec ^
: string
| SET CONFIG tikv _cp1256 = 1
< error
| 1064 42000 …line 1 column 23 near "_cp1256 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_cp1256 = 1
< error
| 1064 42000 …line 1 column 25 near "_cp1256 = 1" 
> exec ^
: string
| KILL _cp1257()
< error
| 1064 42000 …line 1 column 12 near "_cp1257()" 
> exec ^
: string
| KILL _cp1257 (1)
< error
| 1064 42000 …line 1 column 12 near "_cp1257 (1)" 
> exec ^
: string
| SET CONFIG tikv _cp1257 = 1
< error
| 1064 42000 …line 1 column 23 near "_cp1257 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_cp1257 = 1
< error
| 1064 42000 …line 1 column 25 near "_cp1257 = 1" 
> exec ^
: string
| KILL _cp850()
< error
| 1064 42000 …line 1 column 11 near "_cp850()" 
> exec ^
: string
| KILL _cp850 (1)
< error
| 1064 42000 …line 1 column 11 near "_cp850 (1)" 
> exec ^
: string
| SET CONFIG tikv _cp850 = 1
< error
| 1064 42000 …line 1 column 22 near "_cp850 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_cp850 = 1
< error
| 1064 42000 …line 1 column 24 near "_cp850 = 1" 
> exec ^
: string
| KILL _cp852()
< error
| 1064 42000 …line 1 column 11 near "_cp852()" 
> exec ^
: string
| KILL _cp852 (1)
< error
| 1064 42000 …line 1 column 11 near "_cp852 (1)" 
> exec ^
: string
| SET CONFIG tikv _cp852 = 1
< error
| 1064 42000 …line 1 column 22 near "_cp852 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_cp852 = 1
< error
| 1064 42000 …line 1 column 24 near "_cp852 = 1" 
> exec ^
: string
| KILL _cp866()
< error
| 1064 42000 …line 1 column 11 near "_cp866()" 
> exec ^
: string
| KILL _cp866 (1)
< error
| 1064 42000 …line 1 column 11 near "_cp866 (1)" 
> exec ^
: string
| SET CONFIG tikv _cp866 = 1
< error
| 1064 42000 …line 1 column 22 near "_cp866 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_cp866 = 1
< error
| 1064 42000 …line 1 column 24 near "_cp866 = 1" 
> exec ^
: string
| KILL _cp932()
< error
| 1064 42000 …line 1 column 11 near "_cp932()" 
> exec ^
: string
| KILL _cp932 (1)
< error
| 1064 42000 …line 1 column 11 near "_cp932 (1)" 
> exec ^
: string
| SET CONFIG tikv _cp932 = 1
< error
| 1064 42000 …line 1 column 22 near "_cp932 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_cp932 = 1
< error
| 1064 42000 …line 1 column 24 near "_cp932 = 1" 
> exec ^
: string
| KILL _dec8()
< error
| 1064 42000 …line 1 column 10 near "_dec8()" 
> exec ^
: string
| KILL _dec8 (1)
< error
| 1064 42000 …line 1 column 10 near "_dec8 (1)" 
> exec ^
: string
| SET CONFIG tikv _dec8 = 1
< error
| 1064 42000 …line 1 column 21 near "_dec8 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_dec8 = 1
< error
| 1064 42000 …line 1 column 23 near "_dec8 = 1" 
> exec ^
: string
| KILL _eucjpms()
< error
| 1064 42000 …line 1 column 13 near "_eucjpms()" 
> exec ^
: string
| KILL _eucjpms (1)
< error
| 1064 42000 …line 1 column 13 near "_eucjpms (1)" 
> exec ^
: string
| SET CONFIG tikv _eucjpms = 1
< error
| 1064 42000 …line 1 column 24 near "_eucjpms = 1" 
> exec ^
: string
| SET CONFIG tikv a-_eucjpms = 1
< error
| 1064 42000 …line 1 column 26 near "_eucjpms = 1" 
> exec ^
: string
| KILL _euckr()
< error
| 1064 42000 …line 1 column 11 near "_euckr()" 
> exec ^
: string
| KILL _euckr (1)
< error
| 1064 42000 …line 1 column 11 near "_euckr (1)" 
> exec ^
: string
| SET CONFIG tikv _euckr = 1
< error
| 1064 42000 …line 1 column 22 near "_euckr = 1" 
> exec ^
: string
| SET CONFIG tikv a-_euckr = 1
< error
| 1064 42000 …line 1 column 24 near "_euckr = 1" 
> exec ^
: string
| KILL _foo()
< ok 0 0
> exec ^
: string
| KILL _foo (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv _foo = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-_foo = 1
< ok 0 0
> exec ^
: string
| KILL _gb18030()
< error
| 1064 42000 …line 1 column 13 near "_gb18030()" 
> exec ^
: string
| KILL _gb18030 (1)
< error
| 1064 42000 …line 1 column 13 near "_gb18030 (1)" 
> exec ^
: string
| SET CONFIG tikv _gb18030 = 1
< error
| 1064 42000 …line 1 column 24 near "_gb18030 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_gb18030 = 1
< error
| 1064 42000 …line 1 column 26 near "_gb18030 = 1" 
> exec ^
: string
| KILL _gb2312()
< error
| 1064 42000 …line 1 column 12 near "_gb2312()" 
> exec ^
: string
| KILL _gb2312 (1)
< error
| 1064 42000 …line 1 column 12 near "_gb2312 (1)" 
> exec ^
: string
| SET CONFIG tikv _gb2312 = 1
< error
| 1064 42000 …line 1 column 23 near "_gb2312 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_gb2312 = 1
< error
| 1064 42000 …line 1 column 25 near "_gb2312 = 1" 
> exec ^
: string
| KILL _gbk()
< error
| 1064 42000 …line 1 column 9 near "_gbk()" 
> exec ^
: string
| KILL _gbk (1)
< error
| 1064 42000 …line 1 column 9 near "_gbk (1)" 
> exec ^
: string
| SET CONFIG tikv _gbk = 1
< error
| 1064 42000 …line 1 column 20 near "_gbk = 1" 
> exec ^
: string
| SET CONFIG tikv a-_gbk = 1
< error
| 1064 42000 …line 1 column 22 near "_gbk = 1" 
> exec ^
: string
| KILL _geostd8()
< error
| 1064 42000 …line 1 column 13 near "_geostd8()" 
> exec ^
: string
| KILL _geostd8 (1)
< error
| 1064 42000 …line 1 column 13 near "_geostd8 (1)" 
> exec ^
: string
| SET CONFIG tikv _geostd8 = 1
< error
| 1064 42000 …line 1 column 24 near "_geostd8 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_geostd8 = 1
< error
| 1064 42000 …line 1 column 26 near "_geostd8 = 1" 
> exec ^
: string
| KILL _greek()
< error
| 1064 42000 …line 1 column 11 near "_greek()" 
> exec ^
: string
| KILL _greek (1)
< error
| 1064 42000 …line 1 column 11 near "_greek (1)" 
> exec ^
: string
| SET CONFIG tikv _greek = 1
< error
| 1064 42000 …line 1 column 22 near "_greek = 1" 
> exec ^
: string
| SET CONFIG tikv a-_greek = 1
< error
| 1064 42000 …line 1 column 24 near "_greek = 1" 
> exec ^
: string
| KILL _hebrew()
< error
| 1064 42000 …line 1 column 12 near "_hebrew()" 
> exec ^
: string
| KILL _hebrew (1)
< error
| 1064 42000 …line 1 column 12 near "_hebrew (1)" 
> exec ^
: string
| SET CONFIG tikv _hebrew = 1
< error
| 1064 42000 …line 1 column 23 near "_hebrew = 1" 
> exec ^
: string
| SET CONFIG tikv a-_hebrew = 1
< error
| 1064 42000 …line 1 column 25 near "_hebrew = 1" 
> exec ^
: string
| KILL _hp8()
< error
| 1064 42000 …line 1 column 9 near "_hp8()" 
> exec ^
: string
| KILL _hp8 (1)
< error
| 1064 42000 …line 1 column 9 near "_hp8 (1)" 
> exec ^
: string
| SET CONFIG tikv _hp8 = 1
< error
| 1064 42000 …line 1 column 20 near "_hp8 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_hp8 = 1
< error
| 1064 42000 …line 1 column 22 near "_hp8 = 1" 
> exec ^
: string
| KILL _keybcs2()
< error
| 1064 42000 …line 1 column 13 near "_keybcs2()" 
> exec ^
: string
| KILL _keybcs2 (1)
< error
| 1064 42000 …line 1 column 13 near "_keybcs2 (1)" 
> exec ^
: string
| SET CONFIG tikv _keybcs2 = 1
< error
| 1064 42000 …line 1 column 24 near "_keybcs2 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_keybcs2 = 1
< error
| 1064 42000 …line 1 column 26 near "_keybcs2 = 1" 
> exec ^
: string
| KILL _koi8r()
< error
| 1064 42000 …line 1 column 11 near "_koi8r()" 
> exec ^
: string
| KILL _koi8r (1)
< error
| 1064 42000 …line 1 column 11 near "_koi8r (1)" 
> exec ^
: string
| SET CONFIG tikv _koi8r = 1
< error
| 1064 42000 …line 1 column 22 near "_koi8r = 1" 
> exec ^
: string
| SET CONFIG tikv a-_koi8r = 1
< error
| 1064 42000 …line 1 column 24 near "_koi8r = 1" 
> exec ^
: string
| KILL _koi8u()
< error
| 1064 42000 …line 1 column 11 near "_koi8u()" 
> exec ^
: string
| KILL _koi8u (1)
< error
| 1064 42000 …line 1 column 11 near "_koi8u (1)" 
> exec ^
: string
| SET CONFIG tikv _koi8u = 1
< error
| 1064 42000 …line 1 column 22 near "_koi8u = 1" 
> exec ^
: string
| SET CONFIG tikv a-_koi8u = 1
< error
| 1064 42000 …line 1 column 24 near "_koi8u = 1" 
> exec ^
: string
| KILL _latin1()
< error
| 1064 42000 …line 1 column 12 near "_latin1()" 
> exec ^
: string
| KILL _latin1 (1)
< error
| 1064 42000 …line 1 column 12 near "_latin1 (1)" 
> exec ^
: string
| SET CONFIG tikv _latin1 = 1
< error
| 1064 42000 …line 1 column 23 near "_latin1 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_latin1 = 1
< error
| 1064 42000 …line 1 column 25 near "_latin1 = 1" 
> exec ^
: string
| KILL _latin2()
< error
| 1064 42000 …line 1 column 12 near "_latin2()" 
> exec ^
: string
| KILL _latin2 (1)
< error
| 1064 42000 …line 1 column 12 near "_latin2 (1)" 
> exec ^
: string
| SET CONFIG tikv _latin2 = 1
< error
| 1064 42000 …line 1 column 23 near "_latin2 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_latin2 = 1
< error
| 1064 42000 …line 1 column 25 near "_latin2 = 1" 
> exec ^
: string
| KILL _latin5()
< error
| 1064 42000 …line 1 column 12 near "_latin5()" 
> exec ^
: string
| KILL _latin5 (1)
< error
| 1064 42000 …line 1 column 12 near "_latin5 (1)" 
> exec ^
: string
| SET CONFIG tikv _latin5 = 1
< error
| 1064 42000 …line 1 column 23 near "_latin5 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_latin5 = 1
< error
| 1064 42000 …line 1 column 25 near "_latin5 = 1" 
> exec ^
: string
| KILL _latin7()
< error
| 1064 42000 …line 1 column 12 near "_latin7()" 
> exec ^
: string
| KILL _latin7 (1)
< error
| 1064 42000 …line 1 column 12 near "_latin7 (1)" 
> exec ^
: string
| SET CONFIG tikv _latin7 = 1
< error
| 1064 42000 …line 1 column 23 near "_latin7 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_latin7 = 1
< error
| 1064 42000 …line 1 column 25 near "_latin7 = 1" 
> exec ^
: string
| KILL _macce()
< error
| 1064 42000 …line 1 column 11 near "_macce()" 
> exec ^
: string
| KILL _macce (1)
< error
| 1064 42000 …line 1 column 11 near "_macce (1)" 
> exec ^
: string
| SET CONFIG tikv _macce = 1
< error
| 1064 42000 …line 1 column 22 near "_macce = 1" 
> exec ^
: string
| SET CONFIG tikv a-_macce = 1
< error
| 1064 42000 …line 1 column 24 near "_macce = 1" 
> exec ^
: string
| KILL _macroman()
< error
| 1064 42000 …line 1 column 14 near "_macroman()" 
> exec ^
: string
| KILL _macroman (1)
< error
| 1064 42000 …line 1 column 14 near "_macroman (1)" 
> exec ^
: string
| SET CONFIG tikv _macroman = 1
< error
| 1064 42000 …line 1 column 25 near "_macroman = 1" 
> exec ^
: string
| SET CONFIG tikv a-_macroman = 1
< error
| 1064 42000 …line 1 column 27 near "_macroman = 1" 
> exec ^
: string
| KILL _sjis()
< error
| 1064 42000 …line 1 column 10 near "_sjis()" 
> exec ^
: string
| KILL _sjis (1)
< error
| 1064 42000 …line 1 column 10 near "_sjis (1)" 
> exec ^
: string
| SET CONFIG tikv _sjis = 1
< error
| 1064 42000 …line 1 column 21 near "_sjis = 1" 
> exec ^
: string
| SET CONFIG tikv a-_sjis = 1
< error
| 1064 42000 …line 1 column 23 near "_sjis = 1" 
> exec ^
: string
| KILL _swe7()
< error
| 1064 42000 …line 1 column 10 near "_swe7()" 
> exec ^
: string
| KILL _swe7 (1)
< error
| 1064 42000 …line 1 column 10 near "_swe7 (1)" 
> exec ^
: string
| SET CONFIG tikv _swe7 = 1
< error
| 1064 42000 …line 1 column 21 near "_swe7 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_swe7 = 1
< error
| 1064 42000 …line 1 column 23 near "_swe7 = 1" 
> exec ^
: string
| KILL _tis620()
< error
| 1064 42000 …line 1 column 12 near "_tis620()" 
> exec ^
: string
| KILL _tis620 (1)
< error
| 1064 42000 …line 1 column 12 near "_tis620 (1)" 
> exec ^
: string
| SET CONFIG tikv _tis620 = 1
< error
| 1064 42000 …line 1 column 23 near "_tis620 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_tis620 = 1
< error
| 1064 42000 …line 1 column 25 near "_tis620 = 1" 
> exec ^
: string
| KILL _ucs2()
< error
| 1064 42000 …line 1 column 10 near "_ucs2()" 
> exec ^
: string
| KILL _ucs2 (1)
< error
| 1064 42000 …line 1 column 10 near "_ucs2 (1)" 
> exec ^
: string
| SET CONFIG tikv _ucs2 = 1
< error
| 1064 42000 …line 1 column 21 near "_ucs2 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_ucs2 = 1
< error
| 1064 42000 …line 1 column 23 near "_ucs2 = 1" 
> exec ^
: string
| KILL _ujis()
< error
| 1064 42000 …line 1 column 10 near "_ujis()" 
> exec ^
: string
| KILL _ujis (1)
< error
| 1064 42000 …line 1 column 10 near "_ujis (1)" 
> exec ^
: string
| SET CONFIG tikv _ujis = 1
< error
| 1064 42000 …line 1 column 21 near "_ujis = 1" 
> exec ^
: string
| SET CONFIG tikv a-_ujis = 1
< error
| 1064 42000 …line 1 column 23 near "_ujis = 1" 
> exec ^
: string
| KILL _utf16()
< error
| 1064 42000 …line 1 column 11 near "_utf16()" 
> exec ^
: string
| KILL _utf16 (1)
< error
| 1064 42000 …line 1 column 11 near "_utf16 (1)" 
> exec ^
: string
| SET CONFIG tikv _utf16 = 1
< error
| 1064 42000 …line 1 column 22 near "_utf16 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_utf16 = 1
< error
| 1064 42000 …line 1 column 24 near "_utf16 = 1" 
> exec ^
: string
| KILL _utf16le()
< error
| 1064 42000 …line 1 column 13 near "_utf16le()" 
> exec ^
: string
| KILL _utf16le (1)
< error
| 1064 42000 …line 1 column 13 near "_utf16le (1)" 
> exec ^
: string
| SET CONFIG tikv _utf16le = 1
< error
| 1064 42000 …line 1 column 24 near "_utf16le = 1" 
> exec ^
: string
| SET CONFIG tikv a-_utf16le = 1
< error
| 1064 42000 …line 1 column 26 near "_utf16le = 1" 
> exec ^
: string
| KILL _utf32()
< error
| 1064 42000 …line 1 column 11 near "_utf32()" 
> exec ^
: string
| KILL _utf32 (1)
< error
| 1064 42000 …line 1 column 11 near "_utf32 (1)" 
> exec ^
: string
| SET CONFIG tikv _utf32 = 1
< error
| 1064 42000 …line 1 column 22 near "_utf32 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_utf32 = 1
< error
| 1064 42000 …line 1 column 24 near "_utf32 = 1" 
> exec ^
: string
| KILL _utf8()
< error
| 1064 42000 …line 1 column 10 near "_utf8()" 
> exec ^
: string
| KILL _utf8 (1)
< error
| 1064 42000 …line 1 column 10 near "_utf8 (1)" 
> exec ^
: string
| SET CONFIG tikv _utf8 = 1
< error
| 1064 42000 …line 1 column 21 near "_utf8 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_utf8 = 1
< error
| 1064 42000 …line 1 column 23 near "_utf8 = 1" 
> exec ^
: string
| KILL _utf8mb3()
< error
| 1064 42000 …line 1 column 13 near "_utf8mb3()" 
> exec ^
: string
| KILL _utf8mb3 (1)
< error
| 1064 42000 …line 1 column 13 near "_utf8mb3 (1)" 
> exec ^
: string
| SET CONFIG tikv _utf8mb3 = 1
< error
| 1064 42000 …line 1 column 24 near "_utf8mb3 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_utf8mb3 = 1
< error
| 1064 42000 …line 1 column 26 near "_utf8mb3 = 1" 
> exec ^
: string
| KILL _utf8mb4()
< error
| 1064 42000 …line 1 column 13 near "_utf8mb4()" 
> exec ^
: string
| KILL _utf8mb4 (1)
< error
| 1064 42000 …line 1 column 13 near "_utf8mb4 (1)" 
> exec ^
: string
| SET CONFIG tikv _utf8mb4 = 1
< error
| 1064 42000 …line 1 column 24 near "_utf8mb4 = 1" 
> exec ^
: string
| SET CONFIG tikv a-_utf8mb4 = 1
< error
| 1064 42000 …line 1 column 26 near "_utf8mb4 = 1" 
> exec ^
: string
| KILL absent()
< ok 0 0
> exec ^
: string
| KILL absent (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv absent = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-absent = 1
< ok 0 0
> exec ^
: string
| KILL accessible()
< ok 0 0
> exec ^
: string
| KILL accessible (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv accessible = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-accessible = 1
< ok 0 0
> exec ^
: string
| KILL account()
< error
| 1064 42000 …line 1 column 12 near "account()" 
> exec ^
: string
| KILL account (1)
< error
| 1064 42000 …line 1 column 12 near "account (1)" 
> exec ^
: string
| SET CONFIG tikv account = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-account = 1
< ok 0 0
> exec ^
: string
| KILL action()
< error
| 1064 42000 …line 1 column 11 near "action()" 
> exec ^
: string
| KILL action (1)
< error
| 1064 42000 …line 1 column 11 near "action (1)" 
> exec ^
: string
| SET CONFIG tikv action = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-action = 1
< ok 0 0
> exec ^
: string
| KILL active()
< ok 0 0
> exec ^
: string
| KILL active (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv active = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-active = 1
< ok 0 0
> exec ^
: string
| KILL add()
< error
| 1064 42000 …line 1 column 8 near "add()" 
> exec ^
: string
| KILL add (1)
< error
| 1064 42000 …line 1 column 8 near "add (1)" 
> exec ^
: string
| SET CONFIG tikv add = 1
< error
| 1064 42000 …line 1 column 19 near "add = 1" 
> exec ^
: string
| SET CONFIG tikv a-add = 1
< error
| 1064 42000 …line 1 column 21 near "add = 1" 
> exec ^
: string
| KILL adddate()
< error
| 1064 42000 …line 1 column 12 near "adddate()" 
> exec ^
: string
| KILL adddate (1)
< error
| 1064 42000 …line 1 column 12 near "adddate (1)" 
> exec ^
: string
| SET CONFIG tikv adddate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-adddate = 1
< ok 0 0
> exec ^
: string
| KILL admin()
< error
| 1064 42000 …line 1 column 10 near "admin()" 
> exec ^
: string
| KILL admin (1)
< error
| 1064 42000 …line 1 column 10 near "admin (1)" 
> exec ^
: string
| SET CONFIG tikv admin = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-admin = 1
< ok 0 0
> exec ^
: string
| KILL advise()
< error
| 1064 42000 …line 1 column 11 near "advise()" 
> exec ^
: string
| KILL advise (1)
< error
| 1064 42000 …line 1 column 11 near "advise (1)" 
> exec ^
: string
| SET CONFIG tikv advise = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-advise = 1
< ok 0 0
> exec ^
: string
| KILL affinity()
< error
| 1064 42000 …line 1 column 13 near "affinity()" 
> exec ^
: string
| KILL affinity (1)
< error
| 1064 42000 …line 1 column 13 near "affinity (1)" 
> exec ^
: string
| SET CONFIG tikv affinity = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-affinity = 1
< ok 0 0
> exec ^
: string
| KILL after()
< error
| 1064 42000 …line 1 column 10 near "after()" 
> exec ^
: string
| KILL after (1)
< error
| 1064 42000 …line 1 column 10 near "after (1)" 
> exec ^
: string
| SET CONFIG tikv after = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-after = 1
< ok 0 0
> exec ^
: string
| KILL against()
< error
| 1064 42000 …line 1 column 12 near "against()" 
> exec ^
: string
| KILL against (1)
< error
| 1064 42000 …line 1 column 12 near "against (1)" 
> exec ^
: string
| SET CONFIG tikv against = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-against = 1
< ok 0 0
> exec ^
: string
| KILL aggregate()
< ok 0 0
> exec ^
: string
| KILL aggregate (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv aggregate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-aggregate = 1
< ok 0 0
> exec ^
: string
| KILL ago()
< error
| 1064 42000 …line 1 column 8 near "ago()" 
> exec ^
: string
| KILL ago (1)
< error
| 1064 42000 …line 1 column 8 near "ago (1)" 
> exec ^
: string
| SET CONFIG tikv ago = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ago = 1
< ok 0 0
> exec ^
: string
| KILL algorithm()
< error
| 1064 42000 …line 1 column 14 near "algorithm()" 
> exec ^
: string
| KILL algorithm (1)
< error
| 1064 42000 …line 1 column 14 near "algorithm (1)" 
> exec ^
: string
| SET CONFIG tikv algorithm = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-algorithm = 1
< ok 0 0
> exec ^
: string
| KILL all()
< error
| 1064 42000 …line 1 column 8 near "all()" 
> exec ^
: string
| KILL all (1)
< error
| 1064 42000 …line 1 column 8 near "all (1)" 
> exec ^
: string
| SET CONFIG tikv all = 1
< error
| 1064 42000 …line 1 column 19 near "all = 1" 
> exec ^
: string
| SET CONFIG tikv a-all = 1
< error
| 1064 42000 …line 1 column 21 near "all = 1" 
> exec ^
: string
| KILL allow_missing_files()
< ok 0 0
> exec ^
: string
| KILL allow_missing_files (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv allow_missing_files = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-allow_missing_files = 1
< ok 0 0
> exec ^
: string
| KILL alter()
< error
| 1064 42000 …line 1 column 10 near "alter()" 
> exec ^
: string
| KILL alter (1)
< error
| 1064 42000 …line 1 column 10 near "alter (1)" 
> exec ^
: string
| SET CONFIG tikv alter = 1
< error
| 1064 42000 …line 1 column 21 near "alter = 1" 
> exec ^
: string
| SET CONFIG tikv a-alter = 1
< error
| 1064 42000 …line 1 column 23 near "alter = 1" 
> exec ^
: string
| KILL always()
< error
| 1064 42000 …line 1 column 11 near "always()" 
> exec ^
: string
| KILL always (1)
< error
| 1064 42000 …line 1 column 11 near "always (1)" 
> exec ^
: string
| SET CONFIG tikv always = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-always = 1
< ok 0 0
> exec ^
: string
| KILL analyze()
< error
| 1064 42000 …line 1 column 12 near "analyze()" 
> exec ^
: string
| KILL analyze (1)
< error
| 1064 42000 …line 1 column 12 near "analyze (1)" 
> exec ^
: string
| SET CONFIG tikv analyze = 1
< error
| 1064 42000 …line 1 column 23 near "analyze = 1" 
> exec ^
: string
| SET CONFIG tikv a-analyze = 1
< error
| 1064 42000 …line 1 column 25 near "analyze = 1" 
> exec ^
: string
| KILL and()
< error
| 1064 42000 …line 1 column 8 near "and()" 
> exec ^
: string
| KILL and (1)
< error
| 1064 42000 …line 1 column 8 near "and (1)" 
> exec ^
: string
| SET CONFIG tikv and = 1
< error
| 1064 42000 …line 1 column 19 near "and = 1" 
> exec ^
: string
| SET CONFIG tikv a-and = 1
< error
| 1064 42000 …line 1 column 21 near "and = 1" 
> exec ^
: string
| KILL any()
< error
| 1064 42000 …line 1 column 8 near "any()" 
> exec ^
: string
| KILL any (1)
< error
| 1064 42000 …line 1 column 8 near "any (1)" 
> exec ^
: string
| SET CONFIG tikv any = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-any = 1
< ok 0 0
> exec ^
: string
| KILL apply()
< error
| 1064 42000 …line 1 column 10 near "apply()" 
> exec ^
: string
| KILL apply (1)
< error
| 1064 42000 …line 1 column 10 near "apply (1)" 
> exec ^
: string
| SET CONFIG tikv apply = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-apply = 1
< ok 0 0
> exec ^
: string
| KILL approx_count_distinct()
< error
| 1064 42000 …line 1 column 26 near "approx_count_distinct()" 
> exec ^
: string
| KILL approx_count_distinct (1)
< error
| 1064 42000 …line 1 column 26 near "approx_count_distinct (1)" 
> exec ^
: string
| SET CONFIG tikv approx_count_distinct = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-approx_count_distinct = 1
< ok 0 0
> exec ^
: string
| KILL approx_percentile()
< error
| 1064 42000 …line 1 column 22 near "approx_percentile()" 
> exec ^
: string
| KILL approx_percentile (1)
< error
| 1064 42000 …line 1 column 22 near "approx_percentile (1)" 
> exec ^
: string
| SET CONFIG tikv approx_percentile = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-approx_percentile = 1
< ok 0 0
> exec ^
: string
| KILL array()
< error
| 1064 42000 …line 1 column 10 near "array()" 
> exec ^
: string
| KILL array (1)
< error
| 1064 42000 …line 1 column 10 near "array (1)" 
> exec ^
: string
| SET CONFIG tikv array = 1
< error
| 1064 42000 …line 1 column 21 near "array = 1" 
> exec ^
: string
| SET CONFIG tikv a-array = 1
< error
| 1064 42000 …line 1 column 23 near "array = 1" 
> exec ^
: string
| KILL as()
< error
| 1064 42000 …line 1 column 7 near "as()" 
> exec ^
: string
| KILL as (1)
< error
| 1064 42000 …line 1 column 7 near "as (1)" 
> exec ^
: string
| SET CONFIG tikv as = 1
< error
| 1064 42000 …line 1 column 18 near "as = 1" 
> exec ^
: string
| SET CONFIG tikv a-as = 1
< error
| 1064 42000 …line 1 column 20 near "as = 1" 
> exec ^
: string
| KILL asc()
< error
| 1064 42000 …line 1 column 8 near "asc()" 
> exec ^
: string
| KILL asc (1)
< error
| 1064 42000 …line 1 column 8 near "asc (1)" 
> exec ^
: string
| SET CONFIG tikv asc = 1
< error
| 1064 42000 …line 1 column 19 near "asc = 1" 
> exec ^
: string
| SET CONFIG tikv a-asc = 1
< error
| 1064 42000 …line 1 column 21 near "asc = 1" 
> exec ^
: string
| KILL ascii()
< error
| 1064 42000 …line 1 column 10 near "ascii()" 
> exec ^
: string
| KILL ascii (1)
< error
| 1064 42000 …line 1 column 10 near "ascii (1)" 
> exec ^
: string
| SET CONFIG tikv ascii = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ascii = 1
< ok 0 0
> exec ^
: string
| KILL asensitive()
< ok 0 0
> exec ^
: string
| KILL asensitive (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv asensitive = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-asensitive = 1
< ok 0 0
> exec ^
: string
| KILL assign_gtids_to_anonymous_transactions()
< ok 0 0
> exec ^
: string
| KILL assign_gtids_to_anonymous_transactions (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv assign_gtids_to_anonymous_transactions = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-assign_gtids_to_anonymous_transactions = 1
< ok 0 0
> exec ^
: string
| KILL at()
< ok 0 0
> exec ^
: string
| KILL at (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv at = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-at = 1
< ok 0 0
> exec ^
: string
| KILL attribute()
< error
| 1064 42000 …line 1 column 14 near "attribute()" 
> exec ^
: string
| KILL attribute (1)
< error
| 1064 42000 …line 1 column 14 near "attribute (1)" 
> exec ^
: string
| SET CONFIG tikv attribute = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-attribute = 1
< ok 0 0
> exec ^
: string
| KILL attributes()
< error
| 1064 42000 …line 1 column 15 near "attributes()" 
> exec ^
: string
| KILL attributes (1)
< error
| 1064 42000 …line 1 column 15 near "attributes (1)" 
> exec ^
: string
| SET CONFIG tikv attributes = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-attributes = 1
< ok 0 0
> exec ^
: string
| KILL authentication()
< ok 0 0
> exec ^
: string
| KILL authentication (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv authentication = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-authentication = 1
< ok 0 0
> exec ^
: string
| KILL auto()
< ok 0 0
> exec ^
: string
| KILL auto (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv auto = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-auto = 1
< ok 0 0
> exec ^
: string
| KILL auto_id_cache()
< error
| 1064 42000 …line 1 column 18 near "auto_id_cache()" 
> exec ^
: string
| KILL auto_id_cache (1)
< error
| 1064 42000 …line 1 column 18 near "auto_id_cache (1)" 
> exec ^
: string
| SET CONFIG tikv auto_id_cache = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-auto_id_cache = 1
< ok 0 0
> exec ^
: string
| KILL auto_increment()
< error
| 1064 42000 …line 1 column 19 near "auto_increment()" 
> exec ^
: string
| KILL auto_increment (1)
< error
| 1064 42000 …line 1 column 19 near "auto_increment (1)" 
> exec ^
: string
| SET CONFIG tikv auto_increment = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-auto_increment = 1
< ok 0 0
> exec ^
: string
| KILL auto_random()
< error
| 1064 42000 …line 1 column 16 near "auto_random()" 
> exec ^
: string
| KILL auto_random (1)
< error
| 1064 42000 …line 1 column 16 near "auto_random (1)" 
> exec ^
: string
| SET CONFIG tikv auto_random = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-auto_random = 1
< ok 0 0
> exec ^
: string
| KILL auto_random_base()
< error
| 1064 42000 …line 1 column 21 near "auto_random_base()" 
> exec ^
: string
| KILL auto_random_base (1)
< error
| 1064 42000 …line 1 column 21 near "auto_random_base (1)" 
> exec ^
: string
| SET CONFIG tikv auto_random_base = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-auto_random_base = 1
< ok 0 0
> exec ^
: string
| KILL auto_refresh()
< ok 0 0
> exec ^
: string
| KILL auto_refresh (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv auto_refresh = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-auto_refresh = 1
< ok 0 0
> exec ^
: string
| KILL auto_refresh_source()
< ok 0 0
> exec ^
: string
| KILL auto_refresh_source (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv auto_refresh_source = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-auto_refresh_source = 1
< ok 0 0
> exec ^
: string
| KILL autoextend_size()
< ok 0 0
> exec ^
: string
| KILL autoextend_size (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv autoextend_size = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-autoextend_size = 1
< ok 0 0
> exec ^
: string
| KILL avg()
< error
| 1064 42000 …line 1 column 8 near "avg()" 
> exec ^
: string
| KILL avg (1)
< error
| 1064 42000 …line 1 column 8 near "avg (1)" 
> exec ^
: string
| SET CONFIG tikv avg = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-avg = 1
< ok 0 0
> exec ^
: string
| KILL avg_row_length()
< error
| 1064 42000 …line 1 column 19 near "avg_row_length()" 
> exec ^
: string
| KILL avg_row_length (1)
< error
| 1064 42000 …line 1 column 19 near "avg_row_length (1)" 
> exec ^
: string
| SET CONFIG tikv avg_row_length = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-avg_row_length = 1
< ok 0 0
> exec ^
: string
| KILL backend()
< error
| 1064 42000 …line 1 column 12 near "backend()" 
> exec ^
: string
| KILL backend (1)
< error
| 1064 42000 …line 1 column 12 near "backend (1)" 
> exec ^
: string
| SET CONFIG tikv backend = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-backend = 1
< ok 0 0
> exec ^
: string
| KILL background()
< error
| 1064 42000 …line 1 column 15 near "background()" 
> exec ^
: string
| KILL background (1)
< error
| 1064 42000 …line 1 column 15 near "background (1)" 
> exec ^
: string
| SET CONFIG tikv background = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-background = 1
< ok 0 0
> exec ^
: string
| KILL backup()
< error
| 1064 42000 …line 1 column 11 near "backup()" 
> exec ^
: string
| KILL backup (1)
< error
| 1064 42000 …line 1 column 11 near "backup (1)" 
> exec ^
: string
| SET CONFIG tikv backup = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-backup = 1
< ok 0 0
> exec ^
: string
| KILL backups()
< error
| 1064 42000 …line 1 column 12 near "backups()" 
> exec ^
: string
| KILL backups (1)
< error
| 1064 42000 …line 1 column 12 near "backups (1)" 
> exec ^
: string
| SET CONFIG tikv backups = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-backups = 1
< ok 0 0
> exec ^
: string
| KILL batch()
< error
| 1064 42000 …line 1 column 10 near "batch()" 
> exec ^
: string
| KILL batch (1)
< error
| 1064 42000 …line 1 column 10 near "batch (1)" 
> exec ^
: string
| SET CONFIG tikv batch = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-batch = 1
< ok 0 0
> exec ^
: string
| KILL bdr()
< error
| 1064 42000 …line 1 column 8 near "bdr()" 
> exec ^
: string
| KILL bdr (1)
< error
| 1064 42000 …line 1 column 8 near "bdr (1)" 
> exec ^
: string
| SET CONFIG tikv bdr = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bdr = 1
< ok 0 0
> exec ^
: string
| KILL before()
< ok 0 0
> exec ^
: string
| KILL before (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv before = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-before = 1
< ok 0 0
> exec ^
: string
| KILL begin()
< error
| 1064 42000 …line 1 column 10 near "begin()" 
> exec ^
: string
| KILL begin (1)
< error
| 1064 42000 …line 1 column 10 near "begin (1)" 
> exec ^
: string
| SET CONFIG tikv begin = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-begin = 1
< ok 0 0
> exec ^
: string
| KILL bernoulli()
< error
| 1064 42000 …line 1 column 14 near "bernoulli()" 
> exec ^
: string
| KILL bernoulli (1)
< error
| 1064 42000 …line 1 column 14 near "bernoulli (1)" 
> exec ^
: string
| SET CONFIG tikv bernoulli = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bernoulli = 1
< ok 0 0
> exec ^
: string
| KILL between()
< error
| 1064 42000 …line 1 column 12 near "between()" 
> exec ^
: string
| KILL between (1)
< error
| 1064 42000 …line 1 column 12 near "between (1)" 
> exec ^
: string
| SET CONFIG tikv between = 1
< error
| 1064 42000 …line 1 column 23 near "between = 1" 
> exec ^
: string
| SET CONFIG tikv a-between = 1
< error
| 1064 42000 …line 1 column 25 near "between = 1" 
> exec ^
: string
| KILL bigint()
< error
| 1064 42000 …line 1 column 11 near "bigint()" 
> exec ^
: string
| KILL bigint (1)
< error
| 1064 42000 …line 1 column 11 near "bigint (1)" 
> exec ^
: string
| SET CONFIG tikv bigint = 1
< error
| 1064 42000 …line 1 column 22 near "bigint = 1" 
> exec ^
: string
| SET CONFIG tikv a-bigint = 1
< error
| 1064 42000 …line 1 column 24 near "bigint = 1" 
> exec ^
: string
| KILL binary()
< error
| 1064 42000 …line 1 column 11 near "binary()" 
> exec ^
: string
| KILL binary (1)
< error
| 1064 42000 …line 1 column 11 near "binary (1)" 
> exec ^
: string
| SET CONFIG tikv binary = 1
< error
| 1064 42000 …line 1 column 22 near "binary = 1" 
> exec ^
: string
| SET CONFIG tikv a-binary = 1
< error
| 1064 42000 …line 1 column 24 near "binary = 1" 
> exec ^
: string
| KILL binding()
< error
| 1064 42000 …line 1 column 12 near "binding()" 
> exec ^
: string
| KILL binding (1)
< error
| 1064 42000 …line 1 column 12 near "binding (1)" 
> exec ^
: string
| SET CONFIG tikv binding = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-binding = 1
< ok 0 0
> exec ^
: string
| KILL binding_cache()
< error
| 1064 42000 …line 1 column 18 near "binding_cache()" 
> exec ^
: string
| KILL binding_cache (1)
< error
| 1064 42000 …line 1 column 18 near "binding_cache (1)" 
> exec ^
: string
| SET CONFIG tikv binding_cache = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-binding_cache = 1
< ok 0 0
> exec ^
: string
| KILL bindings()
< error
| 1064 42000 …line 1 column 13 near "bindings()" 
> exec ^
: string
| KILL bindings (1)
< error
| 1064 42000 …line 1 column 13 near "bindings (1)" 
> exec ^
: string
| SET CONFIG tikv bindings = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bindings = 1
< ok 0 0
> exec ^
: string
| KILL binlog()
< error
| 1064 42000 …line 1 column 11 near "binlog()" 
> exec ^
: string
| KILL binlog (1)
< error
| 1064 42000 …line 1 column 11 near "binlog (1)" 
> exec ^
: string
| SET CONFIG tikv binlog = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-binlog = 1
< ok 0 0
> exec ^
: string
| KILL bit()
< error
| 1064 42000 …line 1 column 8 near "bit()" 
> exec ^
: string
| KILL bit (1)
< error
| 1064 42000 …line 1 column 8 near "bit (1)" 
> exec ^
: string
| SET CONFIG tikv bit = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bit = 1
< ok 0 0
> exec ^
: string
| KILL bit_and()
< error
| 1064 42000 …line 1 column 12 near "bit_and()" 
> exec ^
: string
| KILL bit_and (1)
< error
| 1064 42000 …line 1 column 12 near "bit_and (1)" 
> exec ^
: string
| SET CONFIG tikv bit_and = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bit_and = 1
< ok 0 0
> exec ^
: string
| KILL bit_or()
< error
| 1064 42000 …line 1 column 11 near "bit_or()" 
> exec ^
: string
| KILL bit_or (1)
< error
| 1064 42000 …line 1 column 11 near "bit_or (1)" 
> exec ^
: string
| SET CONFIG tikv bit_or = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bit_or = 1
< ok 0 0
> exec ^
: string
| KILL bit_xor()
< error
| 1064 42000 …line 1 column 12 near "bit_xor()" 
> exec ^
: string
| KILL bit_xor (1)
< error
| 1064 42000 …line 1 column 12 near "bit_xor (1)" 
> exec ^
: string
| SET CONFIG tikv bit_xor = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bit_xor = 1
< ok 0 0
> exec ^
: string
| KILL blob()
< error
| 1064 42000 …line 1 column 9 near "blob()" 
> exec ^
: string
| KILL blob (1)
< error
| 1064 42000 …line 1 column 9 near "blob (1)" 
> exec ^
: string
| SET CONFIG tikv blob = 1
< error
| 1064 42000 …line 1 column 20 near "blob = 1" 
> exec ^
: string
| SET CONFIG tikv a-blob = 1
< error
| 1064 42000 …line 1 column 22 near "blob = 1" 
> exec ^
: string
| KILL block()
< error
| 1064 42000 …line 1 column 10 near "block()" 
> exec ^
: string
| KILL block (1)
< error
| 1064 42000 …line 1 column 10 near "block (1)" 
> exec ^
: string
| SET CONFIG tikv block = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-block = 1
< ok 0 0
> exec ^
: string
| KILL bool()
< error
| 1064 42000 …line 1 column 9 near "bool()" 
> exec ^
: string
| KILL bool (1)
< error
| 1064 42000 …line 1 column 9 near "bool (1)" 
> exec ^
: string
| SET CONFIG tikv bool = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bool = 1
< ok 0 0
> exec ^
: string
| KILL boolean()
< error
| 1064 42000 …line 1 column 12 near "boolean()" 
> exec ^
: string
| KILL boolean (1)
< error
| 1064 42000 …line 1 column 12 near "boolean (1)" 
> exec ^
: string
| SET CONFIG tikv boolean = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-boolean = 1
< ok 0 0
> exec ^
: string
| KILL both()
< error
| 1064 42000 …line 1 column 9 near "both()" 
> exec ^
: string
| KILL both (1)
< error
| 1064 42000 …line 1 column 9 near "both (1)" 
> exec ^
: string
| SET CONFIG tikv both = 1
< error
| 1064 42000 …line 1 column 20 near "both = 1" 
> exec ^
: string
| SET CONFIG tikv a-both = 1
< error
| 1064 42000 …line 1 column 22 near "both = 1" 
> exec ^
: string
| KILL bound()
< error
| 1064 42000 …line 1 column 10 near "bound()" 
> exec ^
: string
| KILL bound (1)
< error
| 1064 42000 …line 1 column 10 near "bound (1)" 
> exec ^
: string
| SET CONFIG tikv bound = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bound = 1
< ok 0 0
> exec ^
: string
| KILL br()
< error
| 1064 42000 …line 1 column 7 near "br()" 
> exec ^
: string
| KILL br (1)
< error
| 1064 42000 …line 1 column 7 near "br (1)" 
> exec ^
: string
| SET CONFIG tikv br = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-br = 1
< ok 0 0
> exec ^
: string
| KILL brief()
< error
| 1064 42000 …line 1 column 10 near "brief()" 
> exec ^
: string
| KILL brief (1)
< error
| 1064 42000 …line 1 column 10 near "brief (1)" 
> exec ^
: string
| SET CONFIG tikv brief = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-brief = 1
< ok 0 0
> exec ^
: string
| KILL btree()
< error
| 1064 42000 …line 1 column 10 near "btree()" 
> exec ^
: string
| KILL btree (1)
< error
| 1064 42000 …line 1 column 10 near "btree (1)" 
> exec ^
: string
| SET CONFIG tikv btree = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-btree = 1
< ok 0 0
> exec ^
: string
| KILL buckets()
< error
| 1064 42000 …line 1 column 12 near "buckets()" 
> exec ^
: string
| KILL buckets (1)
< error
| 1064 42000 …line 1 column 12 near "buckets (1)" 
> exec ^
: string
| SET CONFIG tikv buckets = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-buckets = 1
< ok 0 0
> exec ^
: string
| KILL builtins()
< error
| 1064 42000 …line 1 column 13 near "builtins()" 
> exec ^
: string
| KILL builtins (1)
< error
| 1064 42000 …line 1 column 13 near "builtins (1)" 
> exec ^
: string
| SET CONFIG tikv builtins = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-builtins = 1
< ok 0 0
> exec ^
: string
| KILL bulk()
< ok 0 0
> exec ^
: string
| KILL bulk (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv bulk = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-bulk = 1
< ok 0 0
> exec ^
: string
| KILL burstable()
< error
| 1064 42000 …line 1 column 14 near "burstable()" 
> exec ^
: string
| KILL burstable (1)
< error
| 1064 42000 …line 1 column 14 near "burstable (1)" 
> exec ^
: string
| SET CONFIG tikv burstable = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-burstable = 1
< ok 0 0
> exec ^
: string
| KILL by()
< error
| 1064 42000 …line 1 column 7 near "by()" 
> exec ^
: string
| KILL by (1)
< error
| 1064 42000 …line 1 column 7 near "by (1)" 
> exec ^
: string
| SET CONFIG tikv by = 1
< error
| 1064 42000 …line 1 column 18 near "by = 1" 
> exec ^
: string
| SET CONFIG tikv a-by = 1
< error
| 1064 42000 …line 1 column 20 near "by = 1" 
> exec ^
: string
| KILL byte()
< error
| 1064 42000 …line 1 column 9 near "byte()" 
> exec ^
: string
| KILL byte (1)
< error
| 1064 42000 …line 1 column 9 near "byte (1)" 
> exec ^
: string
| SET CONFIG tikv byte = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-byte = 1
< ok 0 0
> exec ^
: string
| KILL cache()
< error
| 1064 42000 …line 1 column 10 near "cache()" 
> exec ^
: string
| KILL cache (1)
< error
| 1064 42000 …line 1 column 10 near "cache (1)" 
> exec ^
: string
| SET CONFIG tikv cache = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cache = 1
< ok 0 0
> exec ^
: string
| KILL calibrate()
< error
| 1064 42000 …line 1 column 14 near "calibrate()" 
> exec ^
: string
| KILL calibrate (1)
< error
| 1064 42000 …line 1 column 14 near "calibrate (1)" 
> exec ^
: string
| SET CONFIG tikv calibrate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-calibrate = 1
< ok 0 0
> exec ^
: string
| KILL call()
< error
| 1064 42000 …line 1 column 9 near "call()" 
> exec ^
: string
| KILL call (1)
< error
| 1064 42000 …line 1 column 9 near "call (1)" 
> exec ^
: string
| SET CONFIG tikv call = 1
< error
| 1064 42000 …line 1 column 20 near "call = 1" 
> exec ^
: string
| SET CONFIG tikv a-call = 1
< error
| 1064 42000 …line 1 column 22 near "call = 1" 
> exec ^
: string
| KILL cancel()
< error
| 1064 42000 …line 1 column 11 near "cancel()" 
> exec ^
: string
| KILL cancel (1)
< error
| 1064 42000 …line 1 column 11 near "cancel (1)" 
> exec ^
: string
| SET CONFIG tikv cancel = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cancel = 1
< ok 0 0
> exec ^
: string
| KILL capture()
< error
| 1064 42000 …line 1 column 12 near "capture()" 
> exec ^
: string
| KILL capture (1)
< error
| 1064 42000 …line 1 column 12 near "capture (1)" 
> exec ^
: string
| SET CONFIG tikv capture = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-capture = 1
< ok 0 0
> exec ^
: string
| KILL cardinality()
< error
| 1064 42000 …line 1 column 16 near "cardinality()" 
> exec ^
: string
| KILL cardinality (1)
< error
| 1064 42000 …line 1 column 16 near "cardinality (1)" 
> exec ^
: string
| SET CONFIG tikv cardinality = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cardinality = 1
< ok 0 0
> exec ^
: string
| KILL cascade()
< error
| 1064 42000 …line 1 column 12 near "cascade()" 
> exec ^
: string
| KILL cascade (1)
< error
| 1064 42000 …line 1 column 12 near "cascade (1)" 
> exec ^
: string
| SET CONFIG tikv cascade = 1
< error
| 1064 42000 …line 1 column 23 near "cascade = 1" 
> exec ^
: string
| SET CONFIG tikv a-cascade = 1
< error
| 1064 42000 …line 1 column 25 near "cascade = 1" 
> exec ^
: string
| KILL cascaded()
< error
| 1064 42000 …line 1 column 13 near "cascaded()" 
> exec ^
: string
| KILL cascaded (1)
< error
| 1064 42000 …line 1 column 13 near "cascaded (1)" 
> exec ^
: string
| SET CONFIG tikv cascaded = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cascaded = 1
< ok 0 0
> exec ^
: string
| KILL case()
< error
| 1064 42000 …line 1 column 9 near "case()" 
> exec ^
: string
| KILL case (1)
< error
| 1064 42000 …line 1 column 9 near "case (1)" 
> exec ^
: string
| SET CONFIG tikv case = 1
< error
| 1064 42000 …line 1 column 20 near "case = 1" 
> exec ^
: string
| SET CONFIG tikv a-case = 1
< error
| 1064 42000 …line 1 column 22 near "case = 1" 
> exec ^
: string
| KILL cast()
< error
| 1064 42000 …line 1 column 9 near "cast()" 
> exec ^
: string
| KILL cast (1)
< error
| 1064 42000 …line 1 column 9 near "cast (1)" 
> exec ^
: string
| SET CONFIG tikv cast = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cast = 1
< ok 0 0
> exec ^
: string
| KILL catalog_name()
< ok 0 0
> exec ^
: string
| KILL catalog_name (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv catalog_name = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-catalog_name = 1
< ok 0 0
> exec ^
: string
| KILL causal()
< error
| 1064 42000 …line 1 column 11 near "causal()" 
> exec ^
: string
| KILL causal (1)
< error
| 1064 42000 …line 1 column 11 near "causal (1)" 
> exec ^
: string
| SET CONFIG tikv causal = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-causal = 1
< ok 0 0
> exec ^
: string
| KILL chain()
< error
| 1064 42000 …line 1 column 10 near "chain()" 
> exec ^
: string
| KILL chain (1)
< error
| 1064 42000 …line 1 column 10 near "chain (1)" 
> exec ^
: string
| SET CONFIG tikv chain = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-chain = 1
< ok 0 0
> exec ^
: string
| KILL challenge_response()
< ok 0 0
> exec ^
: string
| KILL challenge_response (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv challenge_response = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-challenge_response = 1
< ok 0 0
> exec ^
: string
| KILL change()
< error
| 1064 42000 …line 1 column 11 near "change()" 
> exec ^
: string
| KILL change (1)
< error
| 1064 42000 …line 1 column 11 near "change (1)" 
> exec ^
: string
| SET CONFIG tikv change = 1
< error
| 1064 42000 …line 1 column 22 near "change = 1" 
> exec ^
: string
| SET CONFIG tikv a-change = 1
< error
| 1064 42000 …line 1 column 24 near "change = 1" 
> exec ^
: string
| KILL changed()
< ok 0 0
> exec ^
: string
| KILL changed (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv changed = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-changed = 1
< ok 0 0
> exec ^
: string
| KILL channel()
< ok 0 0
> exec ^
: string
| KILL channel (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv channel = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-channel = 1
< ok 0 0
> exec ^
: string
| KILL char()
< error
| 1064 42000 …line 1 column 9 near "char()" 
> exec ^
: string
| KILL char (1)
< error
| 1064 42000 …line 1 column 9 near "char (1)" 
> exec ^
: string
| SET CONFIG tikv char = 1
< error
| 1064 42000 …line 1 column 20 near "char = 1" 
> exec ^
: string
| SET CONFIG tikv a-char = 1
< error
| 1064 42000 …line 1 column 22 near "char = 1" 
> exec ^
: string
| KILL character()
< error
| 1064 42000 …line 1 column 14 near "character()" 
> exec ^
: string
| KILL character (1)
< error
| 1064 42000 …line 1 column 14 near "character (1)" 
> exec ^
: string
| SET CONFIG tikv character = 1
< error
| 1064 42000 …line 1 column 25 near "character = 1" 
> exec ^
: string
| SET CONFIG tikv a-character = 1
< error
| 1064 42000 …line 1 column 27 near "character = 1" 
> exec ^
: string
| KILL charset()
< error
| 1064 42000 …line 1 column 12 near "charset()" 
> exec ^
: string
| KILL charset (1)
< error
| 1064 42000 …line 1 column 12 near "charset (1)" 
> exec ^
: string
| SET CONFIG tikv charset = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-charset = 1
< ok 0 0
> exec ^
: string
| KILL check()
< error
| 1064 42000 …line 1 column 10 near "check()" 
> exec ^
: string
| KILL check (1)
< error
| 1064 42000 …line 1 column 10 near "check (1)" 
> exec ^
: string
| SET CONFIG tikv check = 1
< error
| 1064 42000 …line 1 column 21 near "check = 1" 
> exec ^
: string
| SET CONFIG tikv a-check = 1
< error
| 1064 42000 …line 1 column 23 near "check = 1" 
> exec ^
: string
| KILL checkpoint()
< error
| 1064 42000 …line 1 column 15 near "checkpoint()" 
> exec ^
: string
| KILL checkpoint (1)
< error
| 1064 42000 …line 1 column 15 near "checkpoint (1)" 
> exec ^
: string
| SET CONFIG tikv checkpoint = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-checkpoint = 1
< ok 0 0
> exec ^
: string
| KILL checksum()
< error
| 1064 42000 …line 1 column 13 near "checksum()" 
> exec ^
: string
| KILL checksum (1)
< error
| 1064 42000 …line 1 column 13 near "checksum (1)" 
> exec ^
: string
| SET CONFIG tikv checksum = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-checksum = 1
< ok 0 0
> exec ^
: string
| KILL checksum_concurrency()
< error
| 1064 42000 …line 1 column 25 near "checksum_concurrency()" 
> exec ^
: string
| KILL checksum_concurrency (1)
< error
| 1064 42000 …line 1 column 25 near "checksum_concurrency (1)" 
> exec ^
: string
| SET CONFIG tikv checksum_concurrency = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-checksum_concurrency = 1
< ok 0 0
> exec ^
: string
| KILL cipher()
< error
| 1064 42000 …line 1 column 11 near "cipher()" 
> exec ^
: string
| KILL cipher (1)
< error
| 1064 42000 …line 1 column 11 near "cipher (1)" 
> exec ^
: string
| SET CONFIG tikv cipher = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cipher = 1
< ok 0 0
> exec ^
: string
| KILL class_origin()
< ok 0 0
> exec ^
: string
| KILL class_origin (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv class_origin = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-class_origin = 1
< ok 0 0
> exec ^
: string
| KILL cleanup()
< error
| 1064 42000 …line 1 column 12 near "cleanup()" 
> exec ^
: string
| KILL cleanup (1)
< error
| 1064 42000 …line 1 column 12 near "cleanup (1)" 
> exec ^
: string
| SET CONFIG tikv cleanup = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cleanup = 1
< ok 0 0
> exec ^
: string
| KILL client()
< error
| 1064 42000 …line 1 column 11 near "client()" 
> exec ^
: string
| KILL client (1)
< error
| 1064 42000 …line 1 column 11 near "client (1)" 
> exec ^
: string
| SET CONFIG tikv client = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-client = 1
< ok 0 0
> exec ^
: string
| KILL client_errors_summary()
< error
| 1064 42000 …line 1 column 26 near "client_errors_summary()" 
> exec ^
: string
| KILL client_errors_summary (1)
< error
| 1064 42000 …line 1 column 26 near "client_errors_summary (1)" 
> exec ^
: string
| SET CONFIG tikv client_errors_summary = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-client_errors_summary = 1
< ok 0 0
> exec ^
: string
| KILL clone()
< ok 0 0
> exec ^
: string
| KILL clone (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv clone = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-clone = 1
< ok 0 0
> exec ^
: string
| KILL close()
< error
| 1064 42000 …line 1 column 10 near "close()" 
> exec ^
: string
| KILL close (1)
< error
| 1064 42000 …line 1 column 10 near "close (1)" 
> exec ^
: string
| SET CONFIG tikv close = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-close = 1
< ok 0 0
> exec ^
: string
| KILL cluster()
< error
| 1064 42000 …line 1 column 12 near "cluster()" 
> exec ^
: string
| KILL cluster (1)
< error
| 1064 42000 …line 1 column 12 near "cluster (1)" 
> exec ^
: string
| SET CONFIG tikv cluster = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cluster = 1
< ok 0 0
> exec ^
: string
| KILL clustered()
< error
| 1064 42000 …line 1 column 14 near "clustered()" 
> exec ^
: string
| KILL clustered (1)
< error
| 1064 42000 …line 1 column 14 near "clustered (1)" 
> exec ^
: string
| SET CONFIG tikv clustered = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-clustered = 1
< ok 0 0
> exec ^
: string
| KILL cmsketch()
< error
| 1064 42000 …line 1 column 13 near "cmsketch()" 
> exec ^
: string
| KILL cmsketch (1)
< error
| 1064 42000 …line 1 column 13 near "cmsketch (1)" 
> exec ^
: string
| SET CONFIG tikv cmsketch = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cmsketch = 1
< ok 0 0
> exec ^
: string
| KILL coalesce()
< error
| 1064 42000 …line 1 column 13 near "coalesce()" 
> exec ^
: string
| KILL coalesce (1)
< error
| 1064 42000 …line 1 column 13 near "coalesce (1)" 
> exec ^
: string
| SET CONFIG tikv coalesce = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-coalesce = 1
< ok 0 0
> exec ^
: string
| KILL code()
< ok 0 0
> exec ^
: string
| KILL code (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv code = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-code = 1
< ok 0 0
> exec ^
: string
| KILL collate()
< error
| 1064 42000 …line 1 column 12 near "collate()" 
> exec ^
: string
| KILL collate (1)
< error
| 1064 42000 …line 1 column 12 near "collate (1)" 
> exec ^
: string
| SET CONFIG tikv collate = 1
< error
| 1064 42000 …line 1 column 23 near "collate = 1" 
> exec ^
: string
| SET CONFIG tikv a-collate = 1
< error
| 1064 42000 …line 1 column 25 near "collate = 1" 
> exec ^
: string
| KILL collation()
< error
| 1064 42000 …line 1 column 14 near "collation()" 
> exec ^
: string
| KILL collation (1)
< error
| 1064 42000 …line 1 column 14 near "collation (1)" 
> exec ^
: string
| SET CONFIG tikv collation = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-collation = 1
< ok 0 0
> exec ^
: string
| KILL column()
< error
| 1064 42000 …line 1 column 11 near "column()" 
> exec ^
: string
| KILL column (1)
< error
| 1064 42000 …line 1 column 11 near "column (1)" 
> exec ^
: string
| SET CONFIG tikv column = 1
< error
| 1064 42000 …line 1 column 22 near "column = 1" 
> exec ^
: string
| SET CONFIG tikv a-column = 1
< error
| 1064 42000 …line 1 column 24 near "column = 1" 
> exec ^
: string
| KILL column_format()
< error
| 1064 42000 …line 1 column 18 near "column_format()" 
> exec ^
: string
| KILL column_format (1)
< error
| 1064 42000 …line 1 column 18 near "column_format (1)" 
> exec ^
: string
| SET CONFIG tikv column_format = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-column_format = 1
< ok 0 0
> exec ^
: string
| KILL column_name()
< ok 0 0
> exec ^
: string
| KILL column_name (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv column_name = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-column_name = 1
< ok 0 0
> exec ^
: string
| KILL column_stats_usage()
< error
| 1064 42000 …line 1 column 23 near "column_stats_usage()" 
> exec ^
: string
| KILL column_stats_usage (1)
< error
| 1064 42000 …line 1 column 23 near "column_stats_usage (1)" 
> exec ^
: string
| SET CONFIG tikv column_stats_usage = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-column_stats_usage = 1
< ok 0 0
> exec ^
: string
| KILL columns()
< error
| 1064 42000 …line 1 column 12 near "columns()" 
> exec ^
: string
| KILL columns (1)
< error
| 1064 42000 …line 1 column 12 near "columns (1)" 
> exec ^
: string
| SET CONFIG tikv columns = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-columns = 1
< ok 0 0
> exec ^
: string
| KILL comment()
< error
| 1064 42000 …line 1 column 12 near "comment()" 
> exec ^
: string
| KILL comment (1)
< error
| 1064 42000 …line 1 column 12 near "comment (1)" 
> exec ^
: string
| SET CONFIG tikv comment = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-comment = 1
< ok 0 0
> exec ^
: string
| KILL commit()
< error
| 1064 42000 …line 1 column 11 near "commit()" 
> exec ^
: string
| KILL commit (1)
< error
| 1064 42000 …line 1 column 11 near "commit (1)" 
> exec ^
: string
| SET CONFIG tikv commit = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-commit = 1
< ok 0 0
> exec ^
: string
| KILL committed()
< error
| 1064 42000 …line 1 column 14 near "committed()" 
> exec ^
: string
| KILL committed (1)
< error
| 1064 42000 …line 1 column 14 near "committed (1)" 
> exec ^
: string
| SET CONFIG tikv committed = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-committed = 1
< ok 0 0
> exec ^
: string
| KILL compact()
< error
| 1064 42000 …line 1 column 12 near "compact()" 
> exec ^
: string
| KILL compact (1)
< error
| 1064 42000 …line 1 column 12 near "compact (1)" 
> exec ^
: string
| SET CONFIG tikv compact = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-compact = 1
< ok 0 0
> exec ^
: string
| KILL completion()
< ok 0 0
> exec ^
: string
| KILL completion (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv completion = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-completion = 1
< ok 0 0
> exec ^
: string
| KILL component()
< ok 0 0
> exec ^
: string
| KILL component (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv component = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-component = 1
< ok 0 0
> exec ^
: string
| KILL compressed()
< error
| 1064 42000 …line 1 column 15 near "compressed()" 
> exec ^
: string
| KILL compressed (1)
< error
| 1064 42000 …line 1 column 15 near "compressed (1)" 
> exec ^
: string
| SET CONFIG tikv compressed = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-compressed = 1
< ok 0 0
> exec ^
: string
| KILL compression()
< error
| 1064 42000 …line 1 column 16 near "compression()" 
> exec ^
: string
| KILL compression (1)
< error
| 1064 42000 …line 1 column 16 near "compression (1)" 
> exec ^
: string
| SET CONFIG tikv compression = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-compression = 1
< ok 0 0
> exec ^
: string
| KILL compression_level()
< error
| 1064 42000 …line 1 column 22 near "compression_level()" 
> exec ^
: string
| KILL compression_level (1)
< error
| 1064 42000 …line 1 column 22 near "compression_level (1)" 
> exec ^
: string
| SET CONFIG tikv compression_level = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-compression_level = 1
< ok 0 0
> exec ^
: string
| KILL compression_type()
< error
| 1064 42000 …line 1 column 21 near "compression_type()" 
> exec ^
: string
| KILL compression_type (1)
< error
| 1064 42000 …line 1 column 21 near "compression_type (1)" 
> exec ^
: string
| SET CONFIG tikv compression_type = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-compression_type = 1
< ok 0 0
> exec ^
: string
| KILL concat()
< ok 0 0
> exec ^
: string
| KILL concat (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv concat = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-concat = 1
< ok 0 0
> exec ^
: string
| KILL concurrency()
< error
| 1064 42000 …line 1 column 16 near "concurrency()" 
> exec ^
: string
| KILL concurrency (1)
< error
| 1064 42000 …line 1 column 16 near "concurrency (1)" 
> exec ^
: string
| SET CONFIG tikv concurrency = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-concurrency = 1
< ok 0 0
> exec ^
: string
| KILL concurrent()
< ok 0 0
> exec ^
: string
| KILL concurrent (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv concurrent = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-concurrent = 1
< ok 0 0
> exec ^
: string
| KILL condition()
< ok 0 0
> exec ^
: string
| KILL condition (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv condition = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-condition = 1
< ok 0 0
> exec ^
: string
| KILL config()
< error
| 1064 42000 …line 1 column 11 near "config()" 
> exec ^
: string
| KILL config (1)
< error
| 1064 42000 …line 1 column 11 near "config (1)" 
> exec ^
: string
| SET CONFIG tikv config = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-config = 1
< ok 0 0
> exec ^
: string
| KILL connection()
< error
| 1064 42000 …line 1 column 16 near "()" 
> exec ^
: string
| KILL connection (1)
< error
| 1064 42000 …line 1 column 17 near "(1)" 
> exec ^
: string
| SET CONFIG tikv connection = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-connection = 1
< ok 0 0
> exec ^
: string
| KILL consistency()
< error
| 1064 42000 …line 1 column 16 near "consistency()" 
> exec ^
: string
| KILL consistency (1)
< error
| 1064 42000 …line 1 column 16 near "consistency (1)" 
> exec ^
: string
| SET CONFIG tikv consistency = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-consistency = 1
< ok 0 0
> exec ^
: string
| KILL consistent()
< error
| 1064 42000 …line 1 column 15 near "consistent()" 
> exec ^
: string
| KILL consistent (1)
< error
| 1064 42000 …line 1 column 15 near "consistent (1)" 
> exec ^
: string
| SET CONFIG tikv consistent = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-consistent = 1
< ok 0 0
> exec ^
: string
| KILL constraint()
< error
| 1064 42000 …line 1 column 15 near "constraint()" 
> exec ^
: string
| KILL constraint (1)
< error
| 1064 42000 …line 1 column 15 near "constraint (1)" 
> exec ^
: string
| SET CONFIG tikv constraint = 1
< error
| 1064 42000 …line 1 column 26 near "constraint = 1" 
> exec ^
: string
| SET CONFIG tikv a-constraint = 1
< error
| 1064 42000 …line 1 column 28 near "constraint = 1" 
> exec ^
: string
| KILL constraint_catalog()
< ok 0 0
> exec ^
: string
| KILL constraint_catalog (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv constraint_catalog = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-constraint_catalog = 1
< ok 0 0
> exec ^
: string
| KILL constraint_name()
< ok 0 0
> exec ^
: string
| KILL constraint_name (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv constraint_name = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-constraint_name = 1
< ok 0 0
> exec ^
: string
| KILL constraint_schema()
< ok 0 0
> exec ^
: string
| KILL constraint_schema (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv constraint_schema = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-constraint_schema = 1
< ok 0 0
> exec ^
: string
| KILL constraints()
< error
| 1064 42000 …line 1 column 16 near "constraints()" 
> exec ^
: string
| KILL constraints (1)
< error
| 1064 42000 …line 1 column 16 near "constraints (1)" 
> exec ^
: string
| SET CONFIG tikv constraints = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-constraints = 1
< ok 0 0
> exec ^
: string
| KILL contains()
< ok 0 0
> exec ^
: string
| KILL contains (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv contains = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-contains = 1
< ok 0 0
> exec ^
: string
| KILL context()
< error
| 1064 42000 …line 1 column 12 near "context()" 
> exec ^
: string
| KILL context (1)
< error
| 1064 42000 …line 1 column 12 near "context (1)" 
> exec ^
: string
| SET CONFIG tikv context = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-context = 1
< ok 0 0
> exec ^
: string
| KILL continue()
< error
| 1064 42000 …line 1 column 13 near "continue()" 
> exec ^
: string
| KILL continue (1)
< error
| 1064 42000 …line 1 column 13 near "continue (1)" 
> exec ^
: string
| SET CONFIG tikv continue = 1
< error
| 1064 42000 …line 1 column 24 near "continue = 1" 
> exec ^
: string
| SET CONFIG tikv a-continue = 1
< error
| 1064 42000 …line 1 column 26 near "continue = 1" 
> exec ^
: string
| KILL convert()
< error
| 1064 42000 …line 1 column 12 near "convert()" 
> exec ^
: string
| KILL convert (1)
< error
| 1064 42000 …line 1 column 12 near "convert (1)" 
> exec ^
: string
| SET CONFIG tikv convert = 1
< error
| 1064 42000 …line 1 column 23 near "convert = 1" 
> exec ^
: string
| SET CONFIG tikv a-convert = 1
< error
| 1064 42000 …line 1 column 25 near "convert = 1" 
> exec ^
: string
| KILL cooldown()
< error
| 1064 42000 …line 1 column 13 near "cooldown()" 
> exec ^
: string
| KILL cooldown (1)
< error
| 1064 42000 …line 1 column 13 near "cooldown (1)" 
> exec ^
: string
| SET CONFIG tikv cooldown = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cooldown = 1
< ok 0 0
> exec ^
: string
| KILL copy()
< error
| 1064 42000 …line 1 column 9 near "copy()" 
> exec ^
: string
| KILL copy (1)
< error
| 1064 42000 …line 1 column 9 near "copy (1)" 
> exec ^
: string
| SET CONFIG tikv copy = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-copy = 1
< ok 0 0
> exec ^
: string
| KILL correlation()
< error
| 1064 42000 …line 1 column 16 near "correlation()" 
> exec ^
: string
| KILL correlation (1)
< error
| 1064 42000 …line 1 column 16 near "correlation (1)" 
> exec ^
: string
| SET CONFIG tikv correlation = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-correlation = 1
< ok 0 0
> exec ^
: string
| KILL count()
< error
| 1064 42000 …line 1 column 10 near "count()" 
> exec ^
: string
| KILL count (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv count = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-count = 1
< ok 0 0
> exec ^
: string
| KILL cpu()
< error
| 1064 42000 …line 1 column 8 near "cpu()" 
> exec ^
: string
| KILL cpu (1)
< error
| 1064 42000 …line 1 column 8 near "cpu (1)" 
> exec ^
: string
| SET CONFIG tikv cpu = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cpu = 1
< ok 0 0
> exec ^
: string
| KILL create()
< error
| 1064 42000 …line 1 column 11 near "create()" 
> exec ^
: string
| KILL create (1)
< error
| 1064 42000 …line 1 column 11 near "create (1)" 
> exec ^
: string
| SET CONFIG tikv create = 1
< error
| 1064 42000 …line 1 column 22 near "create = 1" 
> exec ^
: string
| SET CONFIG tikv a-create = 1
< error
| 1064 42000 …line 1 column 24 near "create = 1" 
> exec ^
: string
| KILL cross()
< error
| 1064 42000 …line 1 column 10 near "cross()" 
> exec ^
: string
| KILL cross (1)
< error
| 1064 42000 …line 1 column 10 near "cross (1)" 
> exec ^
: string
| SET CONFIG tikv cross = 1
< error
| 1064 42000 …line 1 column 21 near "cross = 1" 
> exec ^
: string
| SET CONFIG tikv a-cross = 1
< error
| 1064 42000 …line 1 column 23 near "cross = 1" 
> exec ^
: string
| KILL csv_backslash_escape()
< error
| 1064 42000 …line 1 column 25 near "csv_backslash_escape()" 
> exec ^
: string
| KILL csv_backslash_escape (1)
< error
| 1064 42000 …line 1 column 25 near "csv_backslash_escape (1)" 
> exec ^
: string
| SET CONFIG tikv csv_backslash_escape = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-csv_backslash_escape = 1
< ok 0 0
> exec ^
: string
| KILL csv_delimiter()
< error
| 1064 42000 …line 1 column 18 near "csv_delimiter()" 
> exec ^
: string
| KILL csv_delimiter (1)
< error
| 1064 42000 …line 1 column 18 near "csv_delimiter (1)" 
> exec ^
: string
| SET CONFIG tikv csv_delimiter = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-csv_delimiter = 1
< ok 0 0
> exec ^
: string
| KILL csv_header()
< error
| 1064 42000 …line 1 column 15 near "csv_header()" 
> exec ^
: string
| KILL csv_header (1)
< error
| 1064 42000 …line 1 column 15 near "csv_header (1)" 
> exec ^
: string
| SET CONFIG tikv csv_header = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-csv_header = 1
< ok 0 0
> exec ^
: string
| KILL csv_not_null()
< error
| 1064 42000 …line 1 column 17 near "csv_not_null()" 
> exec ^
: string
| KILL csv_not_null (1)
< error
| 1064 42000 …line 1 column 17 near "csv_not_null (1)" 
> exec ^
: string
| SET CONFIG tikv csv_not_null = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-csv_not_null = 1
< ok 0 0
> exec ^
: string
| KILL csv_null()
< error
| 1064 42000 …line 1 column 13 near "csv_null()" 
> exec ^
: string
| KILL csv_null (1)
< error
| 1064 42000 …line 1 column 13 near "csv_null (1)" 
> exec ^
: string
| SET CONFIG tikv csv_null = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-csv_null = 1
< ok 0 0
> exec ^
: string
| KILL csv_separator()
< error
| 1064 42000 …line 1 column 18 near "csv_separator()" 
> exec ^
: string
| KILL csv_separator (1)
< error
| 1064 42000 …line 1 column 18 near "csv_separator (1)" 
> exec ^
: string
| SET CONFIG tikv csv_separator = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-csv_separator = 1
< ok 0 0
> exec ^
: string
| KILL csv_trim_last_separators()
< error
| 1064 42000 …line 1 column 29 near "csv_trim_last_separators()" 
> exec ^
: string
| KILL csv_trim_last_separators (1)
< error
| 1064 42000 …line 1 column 29 near "csv_trim_last_separators (1)" 
> exec ^
: string
| SET CONFIG tikv csv_trim_last_separators = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-csv_trim_last_separators = 1
< ok 0 0
> exec ^
: string
| KILL cube()
< ok 0 0
> exec ^
: string
| KILL cube (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv cube = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cube = 1
< ok 0 0
> exec ^
: string
| KILL cume_dist()
< error
| 1064 42000 …line 1 column 14 near "cume_dist()" 
> exec ^
: string
| KILL cume_dist (1)
< error
| 1064 42000 …line 1 column 14 near "cume_dist (1)" 
> exec ^
: string
| SET CONFIG tikv cume_dist = 1
< error
| 1064 42000 …line 1 column 25 near "cume_dist = 1" 
> exec ^
: string
| SET CONFIG tikv a-cume_dist = 1
< error
| 1064 42000 …line 1 column 27 near "cume_dist = 1" 
> exec ^
: string
| KILL curdate()
< error
| 1064 42000 …line 1 column 12 near "curdate()" 
> exec ^
: string
| KILL curdate (1)
< error
| 1064 42000 …line 1 column 12 near "curdate (1)" 
> exec ^
: string
| SET CONFIG tikv curdate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-curdate = 1
< ok 0 0
> exec ^
: string
| KILL current()
< error
| 1064 42000 …line 1 column 12 near "current()" 
> exec ^
: string
| KILL current (1)
< error
| 1064 42000 …line 1 column 12 near "current (1)" 
> exec ^
: string
| SET CONFIG tikv current = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-current = 1
< ok 0 0
> exec ^
: string
| KILL current_date()
< error
| 1064 42000 …line 1 column 17 near "current_date()" 
> exec ^
: string
| KILL current_date (1)
< error
| 1064 42000 …line 1 column 17 near "current_date (1)" 
> exec ^
: string
| SET CONFIG tikv current_date = 1
< error
| 1064 42000 …line 1 column 28 near "current_date = 1" 
> exec ^
: string
| SET CONFIG tikv a-current_date = 1
< error
| 1064 42000 …line 1 column 30 near "current_date = 1" 
> exec ^
: string
| KILL current_role()
< error
| 1064 42000 …line 1 column 17 near "current_role()" 
> exec ^
: string
| KILL current_role (1)
< error
| 1064 42000 …line 1 column 17 near "current_role (1)" 
> exec ^
: string
| KILL current_time()
< error
| 1064 42000 …line 1 column 17 near "current_time()" 
> exec ^
: string
| KILL current_time (1)
< error
| 1064 42000 …line 1 column 17 near "current_time (1)" 
> exec ^
: string
| SET CONFIG tikv current_time = 1
< error
| 1064 42000 …line 1 column 28 near "current_time = 1" 
> exec ^
: string
| SET CONFIG tikv a-current_time = 1
< error
| 1064 42000 …line 1 column 30 near "current_time = 1" 
> exec ^
: string
| KILL current_timestamp()
< error
| 1064 42000 …line 1 column 22 near "current_timestamp()" 
> exec ^
: string
| KILL current_timestamp (1)
< error
| 1064 42000 …line 1 column 22 near "current_timestamp (1)" 
> exec ^
: string
| SET CONFIG tikv current_timestamp = 1
< error
| 1064 42000 …line 1 column 33 near "current_timestamp = 1" 
> exec ^
: string
| SET CONFIG tikv a-current_timestamp = 1
< error
| 1064 42000 …line 1 column 35 near "current_timestamp = 1" 
> exec ^
: string
| KILL current_user()
< error
| 1064 42000 …line 1 column 17 near "current_user()" 
> exec ^
: string
| KILL current_user (1)
< error
| 1064 42000 …line 1 column 17 near "current_user (1)" 
> exec ^
: string
| SET CONFIG tikv current_user = 1
< error
| 1064 42000 …line 1 column 28 near "current_user = 1" 
> exec ^
: string
| SET CONFIG tikv a-current_user = 1
< error
| 1064 42000 …line 1 column 30 near "current_user = 1" 
> exec ^
: string
| KILL cursor()
< error
| 1064 42000 …line 1 column 11 near "cursor()" 
> exec ^
: string
| KILL cursor (1)
< error
| 1064 42000 …line 1 column 11 near "cursor (1)" 
> exec ^
: string
| SET CONFIG tikv cursor = 1
< error
| 1064 42000 …line 1 column 22 near "cursor = 1" 
> exec ^
: string
| SET CONFIG tikv a-cursor = 1
< error
| 1064 42000 …line 1 column 24 near "cursor = 1" 
> exec ^
: string
| KILL cursor_name()
< ok 0 0
> exec ^
: string
| KILL cursor_name (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv cursor_name = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cursor_name = 1
< ok 0 0
> exec ^
: string
| KILL curtime()
< error
| 1064 42000 …line 1 column 12 near "curtime()" 
> exec ^
: string
| KILL curtime (1)
< error
| 1064 42000 …line 1 column 12 near "curtime (1)" 
> exec ^
: string
| SET CONFIG tikv curtime = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-curtime = 1
< ok 0 0
> exec ^
: string
| KILL cycle()
< error
| 1064 42000 …line 1 column 10 near "cycle()" 
> exec ^
: string
| KILL cycle (1)
< error
| 1064 42000 …line 1 column 10 near "cycle (1)" 
> exec ^
: string
| SET CONFIG tikv cycle = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-cycle = 1
< ok 0 0
> exec ^
: string
| KILL data()
< error
| 1064 42000 …line 1 column 9 near "data()" 
> exec ^
: string
| KILL data (1)
< error
| 1064 42000 …line 1 column 9 near "data (1)" 
> exec ^
: string
| SET CONFIG tikv data = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-data = 1
< ok 0 0
> exec ^
: string
| KILL database()
< error
| 1064 42000 …line 1 column 13 near "database()" 
> exec ^
: string
| KILL database (1)
< error
| 1064 42000 …line 1 column 13 near "database (1)" 
> exec ^
: string
| SET CONFIG tikv database = 1
< error
| 1064 42000 …line 1 column 24 near "database = 1" 
> exec ^
: string
| SET CONFIG tikv a-database = 1
< error
| 1064 42000 …line 1 column 26 near "database = 1" 
> exec ^
: string
| KILL databases()
< error
| 1064 42000 …line 1 column 14 near "databases()" 
> exec ^
: string
| KILL databases (1)
< error
| 1064 42000 …line 1 column 14 near "databases (1)" 
> exec ^
: string
| SET CONFIG tikv databases = 1
< error
| 1064 42000 …line 1 column 25 near "databases = 1" 
> exec ^
: string
| SET CONFIG tikv a-databases = 1
< error
| 1064 42000 …line 1 column 27 near "databases = 1" 
> exec ^
: string
| KILL datafile()
< ok 0 0
> exec ^
: string
| KILL datafile (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv datafile = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-datafile = 1
< ok 0 0
> exec ^
: string
| KILL date()
< error
| 1064 42000 …line 1 column 9 near "date()" 
> exec ^
: string
| KILL date (1)
< error
| 1064 42000 …line 1 column 9 near "date (1)" 
> exec ^
: string
| SET CONFIG tikv date = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-date = 1
< ok 0 0
> exec ^
: string
| KILL date_add()
< error
| 1064 42000 …line 1 column 13 near "date_add()" 
> exec ^
: string
| KILL date_add (1)
< error
| 1064 42000 …line 1 column 13 near "date_add (1)" 
> exec ^
: string
| SET CONFIG tikv date_add = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-date_add = 1
< ok 0 0
> exec ^
: string
| KILL date_sub()
< error
| 1064 42000 …line 1 column 13 near "date_sub()" 
> exec ^
: string
| KILL date_sub (1)
< error
| 1064 42000 …line 1 column 13 near "date_sub (1)" 
> exec ^
: string
| SET CONFIG tikv date_sub = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-date_sub = 1
< ok 0 0
> exec ^
: string
| KILL datetime()
< error
| 1064 42000 …line 1 column 13 near "datetime()" 
> exec ^
: string
| KILL datetime (1)
< error
| 1064 42000 …line 1 column 13 near "datetime (1)" 
> exec ^
: string
| SET CONFIG tikv datetime = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-datetime = 1
< ok 0 0
> exec ^
: string
| KILL day()
< error
| 1064 42000 …line 1 column 8 near "day()" 
> exec ^
: string
| KILL day (1)
< error
| 1064 42000 …line 1 column 8 near "day (1)" 
> exec ^
: string
| SET CONFIG tikv day = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-day = 1
< ok 0 0
> exec ^
: string
| KILL day_hour()
< error
| 1064 42000 …line 1 column 13 near "day_hour()" 
> exec ^
: string
| KILL day_hour (1)
< error
| 1064 42000 …line 1 column 13 near "day_hour (1)" 
> exec ^
: string
| SET CONFIG tikv day_hour = 1
< error
| 1064 42000 …line 1 column 24 near "day_hour = 1" 
> exec ^
: string
| SET CONFIG tikv a-day_hour = 1
< error
| 1064 42000 …line 1 column 26 near "day_hour = 1" 
> exec ^
: string
| KILL day_microsecond()
< error
| 1064 42000 …line 1 column 20 near "day_microsecond()" 
> exec ^
: string
| KILL day_microsecond (1)
< error
| 1064 42000 …line 1 column 20 near "day_microsecond (1)" 
> exec ^
: string
| SET CONFIG tikv day_microsecond = 1
< error
| 1064 42000 …line 1 column 31 near "day_microsecond = 1" 
> exec ^
: string
| SET CONFIG tikv a-day_microsecond = 1
< error
| 1064 42000 …line 1 column 33 near "day_microsecond = 1" 
> exec ^
: string
| KILL day_minute()
< error
| 1064 42000 …line 1 column 15 near "day_minute()" 
> exec ^
: string
| KILL day_minute (1)
< error
| 1064 42000 …line 1 column 15 near "day_minute (1)" 
> exec ^
: string
| SET CONFIG tikv day_minute = 1
< error
| 1064 42000 …line 1 column 26 near "day_minute = 1" 
> exec ^
: string
| SET CONFIG tikv a-day_minute = 1
< error
| 1064 42000 …line 1 column 28 near "day_minute = 1" 
> exec ^
: string
| KILL day_second()
< error
| 1064 42000 …line 1 column 15 near "day_second()" 
> exec ^
: string
| KILL day_second (1)
< error
| 1064 42000 …line 1 column 15 near "day_second (1)" 
> exec ^
: string
| SET CONFIG tikv day_second = 1
< error
| 1064 42000 …line 1 column 26 near "day_second = 1" 
> exec ^
: string
| SET CONFIG tikv a-day_second = 1
< error
| 1064 42000 …line 1 column 28 near "day_second = 1" 
> exec ^
: string
| KILL ddl()
< error
| 1064 42000 …line 1 column 8 near "ddl()" 
> exec ^
: string
| KILL ddl (1)
< error
| 1064 42000 …line 1 column 8 near "ddl (1)" 
> exec ^
: string
| SET CONFIG tikv ddl = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ddl = 1
< ok 0 0
> exec ^
: string
| KILL deallocate()
< error
| 1064 42000 …line 1 column 15 near "deallocate()" 
> exec ^
: string
| KILL deallocate (1)
< error
| 1064 42000 …line 1 column 15 near "deallocate (1)" 
> exec ^
: string
| SET CONFIG tikv deallocate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-deallocate = 1
< ok 0 0
> exec ^
: string
| KILL dec()
< error
| 1064 42000 …line 1 column 8 near "dec()" 
> exec ^
: string
| KILL dec (1)
< error
| 1064 42000 …line 1 column 8 near "dec (1)" 
> exec ^
: string
| SET CONFIG tikv dec = 1
< error
| 1064 42000 …line 1 column 19 near "dec = 1" 
> exec ^
: string
| SET CONFIG tikv a-dec = 1
< error
| 1064 42000 …line 1 column 21 near "dec = 1" 
> exec ^
: string
| KILL decimal()
< error
| 1064 42000 …line 1 column 12 near "decimal()" 
> exec ^
: string
| KILL decimal (1)
< error
| 1064 42000 …line 1 column 12 near "decimal (1)" 
> exec ^
: string
| SET CONFIG tikv decimal = 1
< error
| 1064 42000 …line 1 column 23 near "decimal = 1" 
> exec ^
: string
| SET CONFIG tikv a-decimal = 1
< error
| 1064 42000 …line 1 column 25 near "decimal = 1" 
> exec ^
: string
| KILL declare()
< error
| 1064 42000 …line 1 column 12 near "declare()" 
> exec ^
: string
| KILL declare (1)
< error
| 1064 42000 …line 1 column 12 near "declare (1)" 
> exec ^
: string
| SET CONFIG tikv declare = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-declare = 1
< ok 0 0
> exec ^
: string
| KILL default()
< error
| 1064 42000 …line 1 column 12 near "default()" 
> exec ^
: string
| KILL default (1)
< error
| 1064 42000 …line 1 column 12 near "default (1)" 
> exec ^
: string
| SET CONFIG tikv default = 1
< error
| 1064 42000 …line 1 column 23 near "default = 1" 
> exec ^
: string
| SET CONFIG tikv a-default = 1
< error
| 1064 42000 …line 1 column 25 near "default = 1" 
> exec ^
: string
| KILL default_auth()
< ok 0 0
> exec ^
: string
| KILL default_auth (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv default_auth = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-default_auth = 1
< ok 0 0
> exec ^
: string
| KILL defined()
< error
| 1064 42000 …line 1 column 12 near "defined()" 
> exec ^
: string
| KILL defined (1)
< error
| 1064 42000 …line 1 column 12 near "defined (1)" 
> exec ^
: string
| SET CONFIG tikv defined = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-defined = 1
< ok 0 0
> exec ^
: string
| KILL definer()
< error
| 1064 42000 …line 1 column 12 near "definer()" 
> exec ^
: string
| KILL definer (1)
< error
| 1064 42000 …line 1 column 12 near "definer (1)" 
> exec ^
: string
| SET CONFIG tikv definer = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-definer = 1
< ok 0 0
> exec ^
: string
| KILL definition()
< ok 0 0
> exec ^
: string
| KILL definition (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv definition = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-definition = 1
< ok 0 0
> exec ^
: string
| KILL delay_key_write()
< error
| 1064 42000 …line 1 column 20 near "delay_key_write()" 
> exec ^
: string
| KILL delay_key_write (1)
< error
| 1064 42000 …line 1 column 20 near "delay_key_write (1)" 
> exec ^
: string
| SET CONFIG tikv delay_key_write = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-delay_key_write = 1
< ok 0 0
> exec ^
: string
| KILL delayed()
< error
| 1064 42000 …line 1 column 12 near "delayed()" 
> exec ^
: string
| KILL delayed (1)
< error
| 1064 42000 …line 1 column 12 near "delayed (1)" 
> exec ^
: string
| SET CONFIG tikv delayed = 1
< error
| 1064 42000 …line 1 column 23 near "delayed = 1" 
> exec ^
: string
| SET CONFIG tikv a-delayed = 1
< error
| 1064 42000 …line 1 column 25 near "delayed = 1" 
> exec ^
: string
| KILL delete()
< error
| 1064 42000 …line 1 column 11 near "delete()" 
> exec ^
: string
| KILL delete (1)
< error
| 1064 42000 …line 1 column 11 near "delete (1)" 
> exec ^
: string
| SET CONFIG tikv delete = 1
< error
| 1064 42000 …line 1 column 22 near "delete = 1" 
> exec ^
: string
| SET CONFIG tikv a-delete = 1
< error
| 1064 42000 …line 1 column 24 near "delete = 1" 
> exec ^
: string
| KILL dense_rank()
< error
| 1064 42000 …line 1 column 15 near "dense_rank()" 
> exec ^
: string
| KILL dense_rank (1)
< error
| 1064 42000 …line 1 column 15 near "dense_rank (1)" 
> exec ^
: string
| SET CONFIG tikv dense_rank = 1
< error
| 1064 42000 …line 1 column 26 near "dense_rank = 1" 
> exec ^
: string
| SET CONFIG tikv a-dense_rank = 1
< error
| 1064 42000 …line 1 column 28 near "dense_rank = 1" 
> exec ^
: string
| KILL dependency()
< error
| 1064 42000 …line 1 column 15 near "dependency()" 
> exec ^
: string
| KILL dependency (1)
< error
| 1064 42000 …line 1 column 15 near "dependency (1)" 
> exec ^
: string
| SET CONFIG tikv dependency = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-dependency = 1
< ok 0 0
> exec ^
: string
| KILL depth()
< error
| 1064 42000 …line 1 column 10 near "depth()" 
> exec ^
: string
| KILL depth (1)
< error
| 1064 42000 …line 1 column 10 near "depth (1)" 
> exec ^
: string
| SET CONFIG tikv depth = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-depth = 1
< ok 0 0
> exec ^
: string
| KILL desc()
< error
| 1064 42000 …line 1 column 9 near "desc()" 
> exec ^
: string
| KILL desc (1)
< error
| 1064 42000 …line 1 column 9 near "desc (1)" 
> exec ^
: string
| SET CONFIG tikv desc = 1
< error
| 1064 42000 …line 1 column 20 near "desc = 1" 
> exec ^
: string
| SET CONFIG tikv a-desc = 1
< error
| 1064 42000 …line 1 column 22 near "desc = 1" 
> exec ^
: string
| KILL describe()
< error
| 1064 42000 …line 1 column 13 near "describe()" 
> exec ^
: string
| KILL describe (1)
< error
| 1064 42000 …line 1 column 13 near "describe (1)" 
> exec ^
: string
| SET CONFIG tikv describe = 1
< error
| 1064 42000 …line 1 column 24 near "describe = 1" 
> exec ^
: string
| SET CONFIG tikv a-describe = 1
< error
| 1064 42000 …line 1 column 26 near "describe = 1" 
> exec ^
: string
| KILL description()
< ok 0 0
> exec ^
: string
| KILL description (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv description = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-description = 1
< ok 0 0
> exec ^
: string
| KILL deterministic()
< ok 0 0
> exec ^
: string
| KILL deterministic (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv deterministic = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-deterministic = 1
< ok 0 0
> exec ^
: string
| KILL diagnostics()
< ok 0 0
> exec ^
: string
| KILL diagnostics (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv diagnostics = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-diagnostics = 1
< ok 0 0
> exec ^
: string
| KILL digest()
< error
| 1064 42000 …line 1 column 11 near "digest()" 
> exec ^
: string
| KILL digest (1)
< error
| 1064 42000 …line 1 column 11 near "digest (1)" 
> exec ^
: string
| SET CONFIG tikv digest = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-digest = 1
< ok 0 0
> exec ^
: string
| KILL directory()
< error
| 1064 42000 …line 1 column 14 near "directory()" 
> exec ^
: string
| KILL directory (1)
< error
| 1064 42000 …line 1 column 14 near "directory (1)" 
> exec ^
: string
| SET CONFIG tikv directory = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-directory = 1
< ok 0 0
> exec ^
: string
| KILL disable()
< error
| 1064 42000 …line 1 column 12 near "disable()" 
> exec ^
: string
| KILL disable (1)
< error
| 1064 42000 …line 1 column 12 near "disable (1)" 
> exec ^
: string
| SET CONFIG tikv disable = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-disable = 1
< ok 0 0
> exec ^
: string
| KILL disabled()
< error
| 1064 42000 …line 1 column 13 near "disabled()" 
> exec ^
: string
| KILL disabled (1)
< error
| 1064 42000 …line 1 column 13 near "disabled (1)" 
> exec ^
: string
| SET CONFIG tikv disabled = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-disabled = 1
< ok 0 0
> exec ^
: string
| KILL discard()
< error
| 1064 42000 …line 1 column 12 near "discard()" 
> exec ^
: string
| KILL discard (1)
< error
| 1064 42000 …line 1 column 12 near "discard (1)" 
> exec ^
: string
| SET CONFIG tikv discard = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-discard = 1
< ok 0 0
> exec ^
: string
| KILL disk()
< error
| 1064 42000 …line 1 column 9 near "disk()" 
> exec ^
: string
| KILL disk (1)
< error
| 1064 42000 …line 1 column 9 near "disk (1)" 
> exec ^
: string
| SET CONFIG tikv disk = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-disk = 1
< ok 0 0
> exec ^
: string
| KILL distinct()
< error
| 1064 42000 …line 1 column 13 near "distinct()" 
> exec ^
: string
| KILL distinct (1)
< error
| 1064 42000 …line 1 column 13 near "distinct (1)" 
> exec ^
: string
| SET CONFIG tikv distinct = 1
< error
| 1064 42000 …line 1 column 24 near "distinct = 1" 
> exec ^
: string
| SET CONFIG tikv a-distinct = 1
< error
| 1064 42000 …line 1 column 26 near "distinct = 1" 
> exec ^
: string
| KILL distinctrow()
< error
| 1064 42000 …line 1 column 16 near "distinctrow()" 
> exec ^
: string
| KILL distinctrow (1)
< error
| 1064 42000 …line 1 column 16 near "distinctrow (1)" 
> exec ^
: string
| SET CONFIG tikv distinctrow = 1
< error
| 1064 42000 …line 1 column 27 near "distinctrow = 1" 
> exec ^
: string
| SET CONFIG tikv a-distinctrow = 1
< error
| 1064 42000 …line 1 column 29 near "distinctrow = 1" 
> exec ^
: string
| KILL distribute()
< error
| 1064 42000 …line 1 column 15 near "distribute()" 
> exec ^
: string
| KILL distribute (1)
< error
| 1064 42000 …line 1 column 15 near "distribute (1)" 
> exec ^
: string
| SET CONFIG tikv distribute = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-distribute = 1
< ok 0 0
> exec ^
: string
| KILL distribution()
< error
| 1064 42000 …line 1 column 17 near "distribution()" 
> exec ^
: string
| KILL distribution (1)
< error
| 1064 42000 …line 1 column 17 near "distribution (1)" 
> exec ^
: string
| SET CONFIG tikv distribution = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-distribution = 1
< ok 0 0
> exec ^
: string
| KILL distributions()
< error
| 1064 42000 …line 1 column 18 near "distributions()" 
> exec ^
: string
| KILL distributions (1)
< error
| 1064 42000 …line 1 column 18 near "distributions (1)" 
> exec ^
: string
| SET CONFIG tikv distributions = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-distributions = 1
< ok 0 0
> exec ^
: string
| KILL div()
< error
| 1064 42000 …line 1 column 8 near "div()" 
> exec ^
: string
| KILL div (1)
< error
| 1064 42000 …line 1 column 8 near "div (1)" 
> exec ^
: string
| SET CONFIG tikv div = 1
< error
| 1064 42000 …line 1 column 19 near "div = 1" 
> exec ^
: string
| SET CONFIG tikv a-div = 1
< error
| 1064 42000 …line 1 column 21 near "div = 1" 
> exec ^
: string
| KILL do()
< error
| 1064 42000 …line 1 column 7 near "do()" 
> exec ^
: string
| KILL do (1)
< error
| 1064 42000 …line 1 column 7 near "do (1)" 
> exec ^
: string
| SET CONFIG tikv do = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-do = 1
< ok 0 0
> exec ^
: string
| KILL dot()
< error
| 1064 42000 …line 1 column 8 near "dot()" 
> exec ^
: string
| KILL dot (1)
< error
| 1064 42000 …line 1 column 8 near "dot (1)" 
> exec ^
: string
| SET CONFIG tikv dot = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-dot = 1
< ok 0 0
> exec ^
: string
| KILL double()
< error
| 1064 42000 …line 1 column 11 near "double()" 
> exec ^
: string
| KILL double (1)
< error
| 1064 42000 …line 1 column 11 near "double (1)" 
> exec ^
: string
| SET CONFIG tikv double = 1
< error
| 1064 42000 …line 1 column 22 near "double = 1" 
> exec ^
: string
| SET CONFIG tikv a-double = 1
< error
| 1064 42000 …line 1 column 24 near "double = 1" 
> exec ^
: string
| KILL drop()
< error
| 1064 42000 …line 1 column 9 near "drop()" 
> exec ^
: string
| KILL drop (1)
< error
| 1064 42000 …line 1 column 9 near "drop (1)" 
> exec ^
: string
| SET CONFIG tikv drop = 1
< error
| 1064 42000 …line 1 column 20 near "drop = 1" 
> exec ^
: string
| SET CONFIG tikv a-drop = 1
< error
| 1064 42000 …line 1 column 22 near "drop = 1" 
> exec ^
: string
| KILL dry()
< error
| 1064 42000 …line 1 column 8 near "dry()" 
> exec ^
: string
| KILL dry (1)
< error
| 1064 42000 …line 1 column 8 near "dry (1)" 
> exec ^
: string
| SET CONFIG tikv dry = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-dry = 1
< ok 0 0
> exec ^
: string
| KILL dryrun()
< error
| 1064 42000 …line 1 column 11 near "dryrun()" 
> exec ^
: string
| KILL dryrun (1)
< error
| 1064 42000 …line 1 column 11 near "dryrun (1)" 
> exec ^
: string
| SET CONFIG tikv dryrun = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-dryrun = 1
< ok 0 0
> exec ^
: string
| KILL dual()
< error
| 1064 42000 …line 1 column 9 near "dual()" 
> exec ^
: string
| KILL dual (1)
< error
| 1064 42000 …line 1 column 9 near "dual (1)" 
> exec ^
: string
| SET CONFIG tikv dual = 1
< error
| 1064 42000 …line 1 column 20 near "dual = 1" 
> exec ^
: string
| SET CONFIG tikv a-dual = 1
< error
| 1064 42000 …line 1 column 22 near "dual = 1" 
> exec ^
: string
| KILL duality()
< ok 0 0
> exec ^
: string
| KILL duality (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv duality = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-duality = 1
< ok 0 0
> exec ^
: string
| KILL dump()
< error
| 1064 42000 …line 1 column 9 near "dump()" 
> exec ^
: string
| KILL dump (1)
< error
| 1064 42000 …line 1 column 9 near "dump (1)" 
> exec ^
: string
| SET CONFIG tikv dump = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-dump = 1
< ok 0 0
> exec ^
: string
| KILL dumpfile()
< ok 0 0
> exec ^
: string
| KILL dumpfile (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv dumpfile = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-dumpfile = 1
< ok 0 0
> exec ^
: string
| KILL duplicate()
< error
| 1064 42000 …line 1 column 14 near "duplicate()" 
> exec ^
: string
| KILL duplicate (1)
< error
| 1064 42000 …line 1 column 14 near "duplicate (1)" 
> exec ^
: string
| SET CONFIG tikv duplicate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-duplicate = 1
< ok 0 0
> exec ^
: string
| KILL duration()
< error
| 1064 42000 …line 1 column 13 near "duration()" 
> exec ^
: string
| KILL duration (1)
< error
| 1064 42000 …line 1 column 13 near "duration (1)" 
> exec ^
: string
| SET CONFIG tikv duration = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-duration = 1
< ok 0 0
> exec ^
: string
| KILL dynamic()
< error
| 1064 42000 …line 1 column 12 near "dynamic()" 
> exec ^
: string
| KILL dynamic (1)
< error
| 1064 42000 …line 1 column 12 near "dynamic (1)" 
> exec ^
: string
| SET CONFIG tikv dynamic = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-dynamic = 1
< ok 0 0
> exec ^
: string
| KILL each()
< ok 0 0
> exec ^
: string
| KILL each (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv each = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-each = 1
< ok 0 0
> exec ^
: string
| KILL else()
< error
| 1064 42000 …line 1 column 9 near "else()" 
> exec ^
: string
| KILL else (1)
< error
| 1064 42000 …line 1 column 9 near "else (1)" 
> exec ^
: string
| SET CONFIG tikv else = 1
< error
| 1064 42000 …line 1 column 20 near "else = 1" 
> exec ^
: string
| SET CONFIG tikv a-else = 1
< error
| 1064 42000 …line 1 column 22 near "else = 1" 
> exec ^
: string
| KILL elseif()
< error
| 1064 42000 …line 1 column 11 near "elseif()" 
> exec ^
: string
| KILL elseif (1)
< error
| 1064 42000 …line 1 column 11 near "elseif (1)" 
> exec ^
: string
| SET CONFIG tikv elseif = 1
< error
| 1064 42000 …line 1 column 22 near "elseif = 1" 
> exec ^
: string
| SET CONFIG tikv a-elseif = 1
< error
| 1064 42000 …line 1 column 24 near "elseif = 1" 
> exec ^
: string
| KILL empty()
< ok 0 0
> exec ^
: string
| KILL empty (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv empty = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-empty = 1
< ok 0 0
> exec ^
: string
| KILL enable()
< error
| 1064 42000 …line 1 column 11 near "enable()" 
> exec ^
: string
| KILL enable (1)
< error
| 1064 42000 …line 1 column 11 near "enable (1)" 
> exec ^
: string
| SET CONFIG tikv enable = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-enable = 1
< ok 0 0
> exec ^
: string
| KILL enabled()
< error
| 1064 42000 …line 1 column 12 near "enabled()" 
> exec ^
: string
| KILL enabled (1)
< error
| 1064 42000 …line 1 column 12 near "enabled (1)" 
> exec ^
: string
| SET CONFIG tikv enabled = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-enabled = 1
< ok 0 0
> exec ^
: string
| KILL enclosed()
< error
| 1064 42000 …line 1 column 13 near "enclosed()" 
> exec ^
: string
| KILL enclosed (1)
< error
| 1064 42000 …line 1 column 13 near "enclosed (1)" 
> exec ^
: string
| SET CONFIG tikv enclosed = 1
< error
| 1064 42000 …line 1 column 24 near "enclosed = 1" 
> exec ^
: string
| SET CONFIG tikv a-enclosed = 1
< error
| 1064 42000 …line 1 column 26 near "enclosed = 1" 
> exec ^
: string
| KILL encryption()
< error
| 1064 42000 …line 1 column 15 near "encryption()" 
> exec ^
: string
| KILL encryption (1)
< error
| 1064 42000 …line 1 column 15 near "encryption (1)" 
> exec ^
: string
| SET CONFIG tikv encryption = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-encryption = 1
< ok 0 0
> exec ^
: string
| KILL encryption_keyfile()
< error
| 1064 42000 …line 1 column 23 near "encryption_keyfile()" 
> exec ^
: string
| KILL encryption_keyfile (1)
< error
| 1064 42000 …line 1 column 23 near "encryption_keyfile (1)" 
> exec ^
: string
| SET CONFIG tikv encryption_keyfile = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-encryption_keyfile = 1
< ok 0 0
> exec ^
: string
| KILL encryption_method()
< error
| 1064 42000 …line 1 column 22 near "encryption_method()" 
> exec ^
: string
| KILL encryption_method (1)
< error
| 1064 42000 …line 1 column 22 near "encryption_method (1)" 
> exec ^
: string
| SET CONFIG tikv encryption_method = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-encryption_method = 1
< ok 0 0
> exec ^
: string
| KILL end()
< error
| 1064 42000 …line 1 column 8 near "end()" 
> exec ^
: string
| KILL end (1)
< error
| 1064 42000 …line 1 column 8 near "end (1)" 
> exec ^
: string
| SET CONFIG tikv end = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-end = 1
< ok 0 0
> exec ^
: string
| KILL end_time()
< error
| 1064 42000 …line 1 column 13 near "end_time()" 
> exec ^
: string
| KILL end_time (1)
< error
| 1064 42000 …line 1 column 13 near "end_time (1)" 
> exec ^
: string
| SET CONFIG tikv end_time = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-end_time = 1
< ok 0 0
> exec ^
: string
| KILL ends()
< ok 0 0
> exec ^
: string
| KILL ends (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv ends = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ends = 1
< ok 0 0
> exec ^
: string
| KILL enforced()
< error
| 1064 42000 …line 1 column 13 near "enforced()" 
> exec ^
: string
| KILL enforced (1)
< error
| 1064 42000 …line 1 column 13 near "enforced (1)" 
> exec ^
: string
| SET CONFIG tikv enforced = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-enforced = 1
< ok 0 0
> exec ^
: string
| KILL engine()
< error
| 1064 42000 …line 1 column 11 near "engine()" 
> exec ^
: string
| KILL engine (1)
< error
| 1064 42000 …line 1 column 11 near "engine (1)" 
> exec ^
: string
| SET CONFIG tikv engine = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-engine = 1
< ok 0 0
> exec ^
: string
| KILL engine_attribute()
< ok 0 0
> exec ^
: string
| KILL engine_attribute (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv engine_attribute = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-engine_attribute = 1
< ok 0 0
> exec ^
: string
| KILL engines()
< error
| 1064 42000 …line 1 column 12 near "engines()" 
> exec ^
: string
| KILL engines (1)
< error
| 1064 42000 …line 1 column 12 near "engines (1)" 
> exec ^
: string
| SET CONFIG tikv engines = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-engines = 1
< ok 0 0
> exec ^
: string
| KILL enum()
< error
| 1064 42000 …line 1 column 9 near "enum()" 
> exec ^
: string
| KILL enum (1)
< error
| 1064 42000 …line 1 column 9 near "enum (1)" 
> exec ^
: string
| SET CONFIG tikv enum = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-enum = 1
< ok 0 0
> exec ^
: string
| KILL error()
< error
| 1064 42000 …line 1 column 10 near "error()" 
> exec ^
: string
| KILL error (1)
< error
| 1064 42000 …line 1 column 10 near "error (1)" 
> exec ^
: string
| SET CONFIG tikv error = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-error = 1
< ok 0 0
> exec ^
: string
| KILL errors()
< error
| 1064 42000 …line 1 column 11 near "errors()" 
> exec ^
: string
| KILL errors (1)
< error
| 1064 42000 …line 1 column 11 near "errors (1)" 
> exec ^
: string
| SET CONFIG tikv errors = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-errors = 1
< ok 0 0
> exec ^
: string
| KILL escape()
< error
| 1064 42000 …line 1 column 11 near "escape()" 
> exec ^
: string
| KILL escape (1)
< error
| 1064 42000 …line 1 column 11 near "escape (1)" 
> exec ^
: string
| SET CONFIG tikv escape = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-escape = 1
< ok 0 0
> exec ^
: string
| KILL escaped()
< error
| 1064 42000 …line 1 column 12 near "escaped()" 
> exec ^
: string
| KILL escaped (1)
< error
| 1064 42000 …line 1 column 12 near "escaped (1)" 
> exec ^
: string
| SET CONFIG tikv escaped = 1
< error
| 1064 42000 …line 1 column 23 near "escaped = 1" 
> exec ^
: string
| SET CONFIG tikv a-escaped = 1
< error
| 1064 42000 …line 1 column 25 near "escaped = 1" 
> exec ^
: string
| KILL event()
< error
| 1064 42000 …line 1 column 10 near "event()" 
> exec ^
: string
| KILL event (1)
< error
| 1064 42000 …line 1 column 10 near "event (1)" 
> exec ^
: string
| SET CONFIG tikv event = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-event = 1
< ok 0 0
> exec ^
: string
| KILL events()
< error
| 1064 42000 …line 1 column 11 near "events()" 
> exec ^
: string
| KILL events (1)
< error
| 1064 42000 …line 1 column 11 near "events (1)" 
> exec ^
: string
| SET CONFIG tikv events = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-events = 1
< ok 0 0
> exec ^
: string
| KILL every()
< ok 0 0
> exec ^
: string
| KILL every (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv every = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-every = 1
< ok 0 0
> exec ^
: string
| KILL evolve()
< error
| 1064 42000 …line 1 column 11 near "evolve()" 
> exec ^
: string
| KILL evolve (1)
< error
| 1064 42000 …line 1 column 11 near "evolve (1)" 
> exec ^
: string
| SET CONFIG tikv evolve = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-evolve = 1
< ok 0 0
> exec ^
: string
| KILL exact()
< error
| 1064 42000 …line 1 column 10 near "exact()" 
> exec ^
: string
| KILL exact (1)
< error
| 1064 42000 …line 1 column 10 near "exact (1)" 
> exec ^
: string
| SET CONFIG tikv exact = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-exact = 1
< ok 0 0
> exec ^
: string
| KILL except()
< error
| 1064 42000 …line 1 column 11 near "except()" 
> exec ^
: string
| KILL except (1)
< error
| 1064 42000 …line 1 column 11 near "except (1)" 
> exec ^
: string
| SET CONFIG tikv except = 1
< error
| 1064 42000 …line 1 column 22 near "except = 1" 
> exec ^
: string
| SET CONFIG tikv a-except = 1
< error
| 1064 42000 …line 1 column 24 near "except = 1" 
> exec ^
: string
| KILL exchange()
< error
| 1064 42000 …line 1 column 13 near "exchange()" 
> exec ^
: string
| KILL exchange (1)
< error
| 1064 42000 …line 1 column 13 near "exchange (1)" 
> exec ^
: string
| SET CONFIG tikv exchange = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-exchange = 1
< ok 0 0
> exec ^
: string
| KILL exclude()
< ok 0 0
> exec ^
: string
| KILL exclude (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv exclude = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-exclude = 1
< ok 0 0
> exec ^
: string
| KILL exclusive()
< error
| 1064 42000 …line 1 column 14 near "exclusive()" 
> exec ^
: string
| KILL exclusive (1)
< error
| 1064 42000 …line 1 column 14 near "exclusive (1)" 
> exec ^
: string
| SET CONFIG tikv exclusive = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-exclusive = 1
< ok 0 0
> exec ^
: string
| KILL exec_elapsed()
< error
| 1064 42000 …line 1 column 17 near "exec_elapsed()" 
> exec ^
: string
| KILL exec_elapsed (1)
< error
| 1064 42000 …line 1 column 17 near "exec_elapsed (1)" 
> exec ^
: string
| SET CONFIG tikv exec_elapsed = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-exec_elapsed = 1
< ok 0 0
> exec ^
: string
| KILL execute()
< error
| 1064 42000 …line 1 column 12 near "execute()" 
> exec ^
: string
| KILL execute (1)
< error
| 1064 42000 …line 1 column 12 near "execute (1)" 
> exec ^
: string
| SET CONFIG tikv execute = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-execute = 1
< ok 0 0
> exec ^
: string
| KILL exists()
< error
| 1064 42000 …line 1 column 11 near "exists()" 
> exec ^
: string
| KILL exists (1)
< error
| 1064 42000 …line 1 column 11 near "exists (1)" 
> exec ^
: string
| SET CONFIG tikv exists = 1
< error
| 1064 42000 …line 1 column 22 near "exists = 1" 
> exec ^
: string
| SET CONFIG tikv a-exists = 1
< error
| 1064 42000 …line 1 column 24 near "exists = 1" 
> exec ^
: string
| KILL exit()
< error
| 1064 42000 …line 1 column 9 near "exit()" 
> exec ^
: string
| KILL exit (1)
< error
| 1064 42000 …line 1 column 9 near "exit (1)" 
> exec ^
: string
| SET CONFIG tikv exit = 1
< error
| 1064 42000 …line 1 column 20 near "exit = 1" 
> exec ^
: string
| SET CONFIG tikv a-exit = 1
< error
| 1064 42000 …line 1 column 22 near "exit = 1" 
> exec ^
: string
| KILL expansion()
< error
| 1064 42000 …line 1 column 14 near "expansion()" 
> exec ^
: string
| KILL expansion (1)
< error
| 1064 42000 …line 1 column 14 near "expansion (1)" 
> exec ^
: string
| SET CONFIG tikv expansion = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-expansion = 1
< ok 0 0
> exec ^
: string
| KILL expire()
< error
| 1064 42000 …line 1 column 11 near "expire()" 
> exec ^
: string
| KILL expire (1)
< error
| 1064 42000 …line 1 column 11 near "expire (1)" 
> exec ^
: string
| SET CONFIG tikv expire = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-expire = 1
< ok 0 0
> exec ^
: string
| KILL explain()
< error
| 1064 42000 …line 1 column 12 near "explain()" 
> exec ^
: string
| KILL explain (1)
< error
| 1064 42000 …line 1 column 12 near "explain (1)" 
> exec ^
: string
| SET CONFIG tikv explain = 1
< error
| 1064 42000 …line 1 column 23 near "explain = 1" 
> exec ^
: string
| SET CONFIG tikv a-explain = 1
< error
| 1064 42000 …line 1 column 25 near "explain = 1" 
> exec ^
: string
| KILL export()
< ok 0 0
> exec ^
: string
| KILL export (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv export = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-export = 1
< ok 0 0
> exec ^
: string
| KILL expr_pushdown_blacklist()
< error
| 1064 42000 …line 1 column 28 near "expr_pushdown_blacklist()" 
> exec ^
: string
| KILL expr_pushdown_blacklist (1)
< error
| 1064 42000 …line 1 column 28 near "expr_pushdown_blacklist (1)" 
> exec ^
: string
| SET CONFIG tikv expr_pushdown_blacklist = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-expr_pushdown_blacklist = 1
< ok 0 0
> exec ^
: string
| KILL extended()
< error
| 1064 42000 …line 1 column 13 near "extended()" 
> exec ^
: string
| KILL extended (1)
< error
| 1064 42000 …line 1 column 13 near "extended (1)" 
> exec ^
: string
| SET CONFIG tikv extended = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-extended = 1
< ok 0 0
> exec ^
: string
| KILL extent_size()
< ok 0 0
> exec ^
: string
| KILL extent_size (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv extent_size = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-extent_size = 1
< ok 0 0
> exec ^
: string
| KILL external()
< ok 0 0
> exec ^
: string
| KILL external (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv external = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-external = 1
< ok 0 0
> exec ^
: string
| KILL external_format()
< ok 0 0
> exec ^
: string
| KILL external_format (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv external_format = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-external_format = 1
< ok 0 0
> exec ^
: string
| KILL extract()
< error
| 1064 42000 …line 1 column 12 near "extract()" 
> exec ^
: string
| KILL extract (1)
< error
| 1064 42000 …line 1 column 12 near "extract (1)" 
> exec ^
: string
| SET CONFIG tikv extract = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-extract = 1
< ok 0 0
> exec ^
: string
| KILL factor()
< ok 0 0
> exec ^
: string
| KILL factor (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv factor = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-factor = 1
< ok 0 0
> exec ^
: string
| KILL failed_login_attempts()
< error
| 1064 42000 …line 1 column 26 near "failed_login_attempts()" 
> exec ^
: string
| KILL failed_login_attempts (1)
< error
| 1064 42000 …line 1 column 26 near "failed_login_attempts (1)" 
> exec ^
: string
| SET CONFIG tikv failed_login_attempts = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-failed_login_attempts = 1
< ok 0 0
> exec ^
: string
| KILL false()
< error
| 1064 42000 …line 1 column 10 near "false()" 
> exec ^
: string
| KILL false (1)
< error
| 1064 42000 …line 1 column 10 near "false (1)" 
> exec ^
: string
| SET CONFIG tikv false = 1
< error
| 1064 42000 …line 1 column 21 near "false = 1" 
> exec ^
: string
| SET CONFIG tikv a-false = 1
< error
| 1064 42000 …line 1 column 23 near "false = 1" 
> exec ^
: string
| KILL fast()
< ok 0 0
> exec ^
: string
| KILL fast (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv fast = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-fast = 1
< ok 0 0
> exec ^
: string
| KILL faults()
< error
| 1064 42000 …line 1 column 11 near "faults()" 
> exec ^
: string
| KILL faults (1)
< error
| 1064 42000 …line 1 column 11 near "faults (1)" 
> exec ^
: string
| SET CONFIG tikv faults = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-faults = 1
< ok 0 0
> exec ^
: string
| KILL fetch()
< error
| 1064 42000 …line 1 column 10 near "fetch()" 
> exec ^
: string
| KILL fetch (1)
< error
| 1064 42000 …line 1 column 10 near "fetch (1)" 
> exec ^
: string
| SET CONFIG tikv fetch = 1
< error
| 1064 42000 …line 1 column 21 near "fetch = 1" 
> exec ^
: string
| SET CONFIG tikv a-fetch = 1
< error
| 1064 42000 …line 1 column 23 near "fetch = 1" 
> exec ^
: string
| KILL fields()
< error
| 1064 42000 …line 1 column 11 near "fields()" 
> exec ^
: string
| KILL fields (1)
< error
| 1064 42000 …line 1 column 11 near "fields (1)" 
> exec ^
: string
| SET CONFIG tikv fields = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-fields = 1
< ok 0 0
> exec ^
: string
| KILL file()
< error
| 1064 42000 …line 1 column 9 near "file()" 
> exec ^
: string
| KILL file (1)
< error
| 1064 42000 …line 1 column 9 near "file (1)" 
> exec ^
: string
| SET CONFIG tikv file = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-file = 1
< ok 0 0
> exec ^
: string
| KILL file_block_size()
< ok 0 0
> exec ^
: string
| KILL file_block_size (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv file_block_size = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-file_block_size = 1
< ok 0 0
> exec ^
: string
| KILL file_format()
< ok 0 0
> exec ^
: string
| KILL file_format (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv file_format = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-file_format = 1
< ok 0 0
> exec ^
: string
| KILL file_name()
< ok 0 0
> exec ^
: string
| KILL file_name (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv file_name = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-file_name = 1
< ok 0 0
> exec ^
: string
| KILL file_pattern()
< ok 0 0
> exec ^
: string
| KILL file_pattern (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv file_pattern = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-file_pattern = 1
< ok 0 0
> exec ^
: string
| KILL file_prefix()
< ok 0 0
> exec ^
: string
| KILL file_prefix (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv file_prefix = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-file_prefix = 1
< ok 0 0
> exec ^
: string
| KILL files()
< ok 0 0
> exec ^
: string
| KILL files (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv files = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-files = 1
< ok 0 0
> exec ^
: string
| KILL filter()
< ok 0 0
> exec ^
: string
| KILL filter (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv filter = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-filter = 1
< ok 0 0
> exec ^
: string
| KILL finish()
< ok 0 0
> exec ^
: string
| KILL finish (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv finish = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-finish = 1
< ok 0 0
> exec ^
: string
| KILL first()
< error
| 1064 42000 …line 1 column 10 near "first()" 
> exec ^
: string
| KILL first (1)
< error
| 1064 42000 …line 1 column 10 near "first (1)" 
> exec ^
: string
| SET CONFIG tikv first = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-first = 1
< ok 0 0
> exec ^
: string
| KILL first_value()
< error
| 1064 42000 …line 1 column 16 near "first_value()" 
> exec ^
: string
| KILL first_value (1)
< error
| 1064 42000 …line 1 column 16 near "first_value (1)" 
> exec ^
: string
| SET CONFIG tikv first_value = 1
< error
| 1064 42000 …line 1 column 27 near "first_value = 1" 
> exec ^
: string
| SET CONFIG tikv a-first_value = 1
< error
| 1064 42000 …line 1 column 29 near "first_value = 1" 
> exec ^
: string
| KILL fixed()
< error
| 1064 42000 …line 1 column 10 near "fixed()" 
> exec ^
: string
| KILL fixed (1)
< error
| 1064 42000 …line 1 column 10 near "fixed (1)" 
> exec ^
: string
| SET CONFIG tikv fixed = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-fixed = 1
< ok 0 0
> exec ^
: string
| KILL flashback()
< error
| 1064 42000 …line 1 column 14 near "flashback()" 
> exec ^
: string
| KILL flashback (1)
< error
| 1064 42000 …line 1 column 14 near "flashback (1)" 
> exec ^
: string
| SET CONFIG tikv flashback = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-flashback = 1
< ok 0 0
> exec ^
: string
| KILL float()
< error
| 1064 42000 …line 1 column 10 near "float()" 
> exec ^
: string
| KILL float (1)
< error
| 1064 42000 …line 1 column 10 near "float (1)" 
> exec ^
: string
| SET CONFIG tikv float = 1
< error
| 1064 42000 …line 1 column 21 near "float = 1" 
> exec ^
: string
| SET CONFIG tikv a-float = 1
< error
| 1064 42000 …line 1 column 23 near "float = 1" 
> exec ^
: string
| KILL float4()
< error
| 1064 42000 …line 1 column 11 near "float4()" 
> exec ^
: string
| KILL float4 (1)
< error
| 1064 42000 …line 1 column 11 near "float4 (1)" 
> exec ^
: string
| SET CONFIG tikv float4 = 1
< error
| 1064 42000 …line 1 column 22 near "float4 = 1" 
> exec ^
: string
| SET CONFIG tikv a-float4 = 1
< error
| 1064 42000 …line 1 column 24 near "float4 = 1" 
> exec ^
: string
| KILL float8()
< error
| 1064 42000 …line 1 column 11 near "float8()" 
> exec ^
: string
| KILL float8 (1)
< error
| 1064 42000 …line 1 column 11 near "float8 (1)" 
> exec ^
: string
| SET CONFIG tikv float8 = 1
< error
| 1064 42000 …line 1 column 22 near "float8 = 1" 
> exec ^
: string
| SET CONFIG tikv a-float8 = 1
< error
| 1064 42000 …line 1 column 24 near "float8 = 1" 
> exec ^
: string
| KILL flush()
< error
| 1064 42000 …line 1 column 10 near "flush()" 
> exec ^
: string
| KILL flush (1)
< error
| 1064 42000 …line 1 column 10 near "flush (1)" 
> exec ^
: string
| SET CONFIG tikv flush = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-flush = 1
< ok 0 0
> exec ^
: string
| KILL follower()
< error
| 1064 42000 …line 1 column 13 near "follower()" 
> exec ^
: string
| KILL follower (1)
< error
| 1064 42000 …line 1 column 13 near "follower (1)" 
> exec ^
: string
| SET CONFIG tikv follower = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-follower = 1
< ok 0 0
> exec ^
: string
| KILL follower_constraints()
< error
| 1064 42000 …line 1 column 25 near "follower_constraints()" 
> exec ^
: string
| KILL follower_constraints (1)
< error
| 1064 42000 …line 1 column 25 near "follower_constraints (1)" 
> exec ^
: string
| SET CONFIG tikv follower_constraints = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-follower_constraints = 1
< ok 0 0
> exec ^
: string
| KILL followers()
< error
| 1064 42000 …line 1 column 14 near "followers()" 
> exec ^
: string
| KILL followers (1)
< error
| 1064 42000 …line 1 column 14 near "followers (1)" 
> exec ^
: string
| SET CONFIG tikv followers = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-followers = 1
< ok 0 0
> exec ^
: string
| KILL following()
< error
| 1064 42000 …line 1 column 14 near "following()" 
> exec ^
: string
| KILL following (1)
< error
| 1064 42000 …line 1 column 14 near "following (1)" 
> exec ^
: string
| SET CONFIG tikv following = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-following = 1
< ok 0 0
> exec ^
: string
| KILL follows()
< ok 0 0
> exec ^
: string
| KILL follows (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv follows = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-follows = 1
< ok 0 0
> exec ^
: string
| KILL foo()
< ok 0 0
> exec ^
: string
| KILL foo (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv foo = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-foo = 1
< ok 0 0
> exec ^
: string
| KILL for()
< error
| 1064 42000 …line 1 column 8 near "for()" 
> exec ^
: string
| KILL for (1)
< error
| 1064 42000 …line 1 column 8 near "for (1)" 
> exec ^
: string
| SET CONFIG tikv for = 1
< error
| 1064 42000 …line 1 column 19 near "for = 1" 
> exec ^
: string
| SET CONFIG tikv a-for = 1
< error
| 1064 42000 …line 1 column 21 near "for = 1" 
> exec ^
: string
| KILL force()
< error
| 1064 42000 …line 1 column 10 near "force()" 
> exec ^
: string
| KILL force (1)
< error
| 1064 42000 …line 1 column 10 near "force (1)" 
> exec ^
: string
| SET CONFIG tikv force = 1
< error
| 1064 42000 …line 1 column 21 near "force = 1" 
> exec ^
: string
| SET CONFIG tikv a-force = 1
< error
| 1064 42000 …line 1 column 23 near "force = 1" 
> exec ^
: string
| KILL foreign()
< error
| 1064 42000 …line 1 column 12 near "foreign()" 
> exec ^
: string
| KILL foreign (1)
< error
| 1064 42000 …line 1 column 12 near "foreign (1)" 
> exec ^
: string
| SET CONFIG tikv foreign = 1
< error
| 1064 42000 …line 1 column 23 near "foreign = 1" 
> exec ^
: string
| SET CONFIG tikv a-foreign = 1
< error
| 1064 42000 …line 1 column 25 near "foreign = 1" 
> exec ^
: string
| KILL format()
< error
| 1064 42000 …line 1 column 11 near "format()" 
> exec ^
: string
| KILL format (1)
< error
| 1064 42000 …line 1 column 11 near "format (1)" 
> exec ^
: string
| SET CONFIG tikv format = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-format = 1
< ok 0 0
> exec ^
: string
| KILL found()
< error
| 1064 42000 …line 1 column 10 near "found()" 
> exec ^
: string
| KILL found (1)
< error
| 1064 42000 …line 1 column 10 near "found (1)" 
> exec ^
: string
| SET CONFIG tikv found = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-found = 1
< ok 0 0
> exec ^
: string
| KILL from()
< error
| 1064 42000 …line 1 column 9 near "from()" 
> exec ^
: string
| KILL from (1)
< error
| 1064 42000 …line 1 column 9 near "from (1)" 
> exec ^
: string
| SET CONFIG tikv from = 1
< error
| 1064 42000 …line 1 column 20 near "from = 1" 
> exec ^
: string
| SET CONFIG tikv a-from = 1
< error
| 1064 42000 …line 1 column 22 near "from = 1" 
> exec ^
: string
| KILL full()
< error
| 1064 42000 …line 1 column 9 near "full()" 
> exec ^
: string
| KILL full (1)
< error
| 1064 42000 …line 1 column 9 near "full (1)" 
> exec ^
: string
| SET CONFIG tikv full = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-full = 1
< ok 0 0
> exec ^
: string
| KILL full_backup_storage()
< error
| 1064 42000 …line 1 column 24 near "full_backup_storage()" 
> exec ^
: string
| KILL full_backup_storage (1)
< error
| 1064 42000 …line 1 column 24 near "full_backup_storage (1)" 
> exec ^
: string
| SET CONFIG tikv full_backup_storage = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-full_backup_storage = 1
< ok 0 0
> exec ^
: string
| KILL fulltext()
< error
| 1064 42000 …line 1 column 13 near "fulltext()" 
> exec ^
: string
| KILL fulltext (1)
< error
| 1064 42000 …line 1 column 13 near "fulltext (1)" 
> exec ^
: string
| SET CONFIG tikv fulltext = 1
< error
| 1064 42000 …line 1 column 24 near "fulltext = 1" 
> exec ^
: string
| SET CONFIG tikv a-fulltext = 1
< error
| 1064 42000 …line 1 column 26 near "fulltext = 1" 
> exec ^
: string
| KILL function()
< error
| 1064 42000 …line 1 column 13 near "function()" 
> exec ^
: string
| KILL function (1)
< error
| 1064 42000 …line 1 column 13 near "function (1)" 
> exec ^
: string
| SET CONFIG tikv function = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-function = 1
< ok 0 0
> exec ^
: string
| KILL gc_ttl()
< error
| 1064 42000 …line 1 column 11 near "gc_ttl()" 
> exec ^
: string
| KILL gc_ttl (1)
< error
| 1064 42000 …line 1 column 11 near "gc_ttl (1)" 
> exec ^
: string
| SET CONFIG tikv gc_ttl = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-gc_ttl = 1
< ok 0 0
> exec ^
: string
| KILL general()
< error
| 1064 42000 …line 1 column 12 near "general()" 
> exec ^
: string
| KILL general (1)
< error
| 1064 42000 …line 1 column 12 near "general (1)" 
> exec ^
: string
| SET CONFIG tikv general = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-general = 1
< ok 0 0
> exec ^
: string
| KILL generate()
< ok 0 0
> exec ^
: string
| KILL generate (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv generate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-generate = 1
< ok 0 0
> exec ^
: string
| KILL generated()
< error
| 1064 42000 …line 1 column 14 near "generated()" 
> exec ^
: string
| KILL generated (1)
< error
| 1064 42000 …line 1 column 14 near "generated (1)" 
> exec ^
: string
| SET CONFIG tikv generated = 1
< error
| 1064 42000 …line 1 column 25 near "generated = 1" 
> exec ^
: string
| SET CONFIG tikv a-generated = 1
< error
| 1064 42000 …line 1 column 27 near "generated = 1" 
> exec ^
: string
| KILL geomcollection()
< ok 0 0
> exec ^
: string
| KILL geomcollection (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv geomcollection = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-geomcollection = 1
< ok 0 0
> exec ^
: string
| KILL geometry()
< ok 0 0
> exec ^
: string
| KILL geometry (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv geometry = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-geometry = 1
< ok 0 0
> exec ^
: string
| KILL geometrycollection()
< ok 0 0
> exec ^
: string
| KILL geometrycollection (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv geometrycollection = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-geometrycollection = 1
< ok 0 0
> exec ^
: string
| KILL get()
< ok 0 0
> exec ^
: string
| KILL get (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv get = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-get = 1
< ok 0 0
> exec ^
: string
| KILL get_format()
< error
| 1064 42000 …line 1 column 15 near "get_format()" 
> exec ^
: string
| KILL get_format (1)
< error
| 1064 42000 …line 1 column 15 near "get_format (1)" 
> exec ^
: string
| SET CONFIG tikv get_format = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-get_format = 1
< ok 0 0
> exec ^
: string
| KILL get_source_public_key()
< ok 0 0
> exec ^
: string
| KILL get_source_public_key (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv get_source_public_key = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-get_source_public_key = 1
< ok 0 0
> exec ^
: string
| KILL global()
< error
| 1064 42000 …line 1 column 11 near "global()" 
> exec ^
: string
| KILL global (1)
< error
| 1064 42000 …line 1 column 11 near "global (1)" 
> exec ^
: string
| SET CONFIG tikv global = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-global = 1
< ok 0 0
> exec ^
: string
| KILL grant()
< error
| 1064 42000 …line 1 column 10 near "grant()" 
> exec ^
: string
| KILL grant (1)
< error
| 1064 42000 …line 1 column 10 near "grant (1)" 
> exec ^
: string
| SET CONFIG tikv grant = 1
< error
| 1064 42000 …line 1 column 21 near "grant = 1" 
> exec ^
: string
| SET CONFIG tikv a-grant = 1
< error
| 1064 42000 …line 1 column 23 near "grant = 1" 
> exec ^
: string
| KILL grants()
< error
| 1064 42000 …line 1 column 11 near "grants()" 
> exec ^
: string
| KILL grants (1)
< error
| 1064 42000 …line 1 column 11 near "grants (1)" 
> exec ^
: string
| SET CONFIG tikv grants = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-grants = 1
< ok 0 0
> exec ^
: string
| KILL greatest()
< ok 0 0
> exec ^
: string
| KILL greatest (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv greatest = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-greatest = 1
< ok 0 0
> exec ^
: string
| KILL group()
< error
| 1064 42000 …line 1 column 10 near "group()" 
> exec ^
: string
| KILL group (1)
< error
| 1064 42000 …line 1 column 10 near "group (1)" 
> exec ^
: string
| SET CONFIG tikv group = 1
< error
| 1064 42000 …line 1 column 21 near "group = 1" 
> exec ^
: string
| SET CONFIG tikv a-group = 1
< error
| 1064 42000 …line 1 column 23 near "group = 1" 
> exec ^
: string
| KILL group_concat()
< error
| 1064 42000 …line 1 column 17 near "group_concat()" 
> exec ^
: string
| KILL group_concat (1)
< error
| 1064 42000 …line 1 column 17 near "group_concat (1)" 
> exec ^
: string
| SET CONFIG tikv group_concat = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-group_concat = 1
< ok 0 0
> exec ^
: string
| KILL group_replication()
< ok 0 0
> exec ^
: string
| KILL group_replication (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv group_replication = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-group_replication = 1
< ok 0 0
> exec ^
: string
| KILL grouping()
< ok 0 0
> exec ^
: string
| KILL grouping (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv grouping = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-grouping = 1
< ok 0 0
> exec ^
: string
| KILL groups()
< error
| 1064 42000 …line 1 column 11 near "groups()" 
> exec ^
: string
| KILL groups (1)
< error
| 1064 42000 …line 1 column 11 near "groups (1)" 
> exec ^
: string
| SET CONFIG tikv groups = 1
< error
| 1064 42000 …line 1 column 22 near "groups = 1" 
> exec ^
: string
| SET CONFIG tikv a-groups = 1
< error
| 1064 42000 …line 1 column 24 near "groups = 1" 
> exec ^
: string
| KILL gtid_only()
< ok 0 0
> exec ^
: string
| KILL gtid_only (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv gtid_only = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-gtid_only = 1
< ok 0 0
> exec ^
: string
| KILL gtids()
< ok 0 0
> exec ^
: string
| KILL gtids (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv gtids = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-gtids = 1
< ok 0 0
> exec ^
: string
| KILL guided()
< ok 0 0
> exec ^
: string
| KILL guided (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv guided = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-guided = 1
< ok 0 0
> exec ^
: string
| KILL handler()
< error
| 1064 42000 …line 1 column 12 near "handler()" 
> exec ^
: string
| KILL handler (1)
< error
| 1064 42000 …line 1 column 12 near "handler (1)" 
> exec ^
: string
| SET CONFIG tikv handler = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-handler = 1
< ok 0 0
> exec ^
: string
| KILL hash()
< error
| 1064 42000 …line 1 column 9 near "hash()" 
> exec ^
: string
| KILL hash (1)
< error
| 1064 42000 …line 1 column 9 near "hash (1)" 
> exec ^
: string
| SET CONFIG tikv hash = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-hash = 1
< ok 0 0
> exec ^
: string
| KILL having()
< error
| 1064 42000 …line 1 column 11 near "having()" 
> exec ^
: string
| KILL having (1)
< error
| 1064 42000 …line 1 column 11 near "having (1)" 
> exec ^
: string
| SET CONFIG tikv having = 1
< error
| 1064 42000 …line 1 column 22 near "having = 1" 
> exec ^
: string
| SET CONFIG tikv a-having = 1
< error
| 1064 42000 …line 1 column 24 near "having = 1" 
> exec ^
: string
| KILL header()
< ok 0 0
> exec ^
: string
| KILL header (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv header = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-header = 1
< ok 0 0
> exec ^
: string
| KILL help()
< error
| 1064 42000 …line 1 column 9 near "help()" 
> exec ^
: string
| KILL help (1)
< error
| 1064 42000 …line 1 column 9 near "help (1)" 
> exec ^
: string
| SET CONFIG tikv help = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-help = 1
< ok 0 0
> exec ^
: string
| KILL high()
< error
| 1064 42000 …line 1 column 9 near "high()" 
> exec ^
: string
| KILL high (1)
< error
| 1064 42000 …line 1 column 9 near "high (1)" 
> exec ^
: string
| SET CONFIG tikv high = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-high = 1
< ok 0 0
> exec ^
: string
| KILL high_priority()
< error
| 1064 42000 …line 1 column 18 near "high_priority()" 
> exec ^
: string
| KILL high_priority (1)
< error
| 1064 42000 …line 1 column 18 near "high_priority (1)" 
> exec ^
: string
| SET CONFIG tikv high_priority = 1
< error
| 1064 42000 …line 1 column 29 near "high_priority = 1" 
> exec ^
: string
| SET CONFIG tikv a-high_priority = 1
< error
| 1064 42000 …line 1 column 31 near "high_priority = 1" 
> exec ^
: string
| KILL histogram()
< error
| 1064 42000 …line 1 column 14 near "histogram()" 
> exec ^
: string
| KILL histogram (1)
< error
| 1064 42000 …line 1 column 14 near "histogram (1)" 
> exec ^
: string
| SET CONFIG tikv histogram = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-histogram = 1
< ok 0 0
> exec ^
: string
| KILL histograms_in_flight()
< error
| 1064 42000 …line 1 column 25 near "histograms_in_flight()" 
> exec ^
: string
| KILL histograms_in_flight (1)
< error
| 1064 42000 …line 1 column 25 near "histograms_in_flight (1)" 
> exec ^
: string
| SET CONFIG tikv histograms_in_flight = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-histograms_in_flight = 1
< ok 0 0
> exec ^
: string
| KILL history()
< error
| 1064 42000 …line 1 column 12 near "history()" 
> exec ^
: string
| KILL history (1)
< error
| 1064 42000 …line 1 column 12 near "history (1)" 
> exec ^
: string
| SET CONFIG tikv history = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-history = 1
< ok 0 0
> exec ^
: string
| KILL hnsw()
< error
| 1064 42000 …line 1 column 9 near "hnsw()" 
> exec ^
: string
| KILL hnsw (1)
< error
| 1064 42000 …line 1 column 9 near "hnsw (1)" 
> exec ^
: string
| SET CONFIG tikv hnsw = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-hnsw = 1
< ok 0 0
> exec ^
: string
| KILL host()
< ok 0 0
> exec ^
: string
| KILL host (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv host = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-host = 1
< ok 0 0
> exec ^
: string
| KILL hosts()
< error
| 1064 42000 …line 1 column 10 near "hosts()" 
> exec ^
: string
| KILL hosts (1)
< error
| 1064 42000 …line 1 column 10 near "hosts (1)" 
> exec ^
: string
| SET CONFIG tikv hosts = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-hosts = 1
< ok 0 0
> exec ^
: string
| KILL hour()
< error
| 1064 42000 …line 1 column 9 near "hour()" 
> exec ^
: string
| KILL hour (1)
< error
| 1064 42000 …line 1 column 9 near "hour (1)" 
> exec ^
: string
| SET CONFIG tikv hour = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-hour = 1
< ok 0 0
> exec ^
: string
| KILL hour_microsecond()
< error
| 1064 42000 …line 1 column 21 near "hour_microsecond()" 
> exec ^
: string
| KILL hour_microsecond (1)
< error
| 1064 42000 …line 1 column 21 near "hour_microsecond (1)" 
> exec ^
: string
| SET CONFIG tikv hour_microsecond = 1
< error
| 1064 42000 …line 1 column 32 near "hour_microsecond = 1" 
> exec ^
: string
| SET CONFIG tikv a-hour_microsecond = 1
< error
| 1064 42000 …line 1 column 34 near "hour_microsecond = 1" 
> exec ^
: string
| KILL hour_minute()
< error
| 1064 42000 …line 1 column 16 near "hour_minute()" 
> exec ^
: string
| KILL hour_minute (1)
< error
| 1064 42000 …line 1 column 16 near "hour_minute (1)" 
> exec ^
: string
| SET CONFIG tikv hour_minute = 1
< error
| 1064 42000 …line 1 column 27 near "hour_minute = 1" 
> exec ^
: string
| SET CONFIG tikv a-hour_minute = 1
< error
| 1064 42000 …line 1 column 29 near "hour_minute = 1" 
> exec ^
: string
| KILL hour_second()
< error
| 1064 42000 …line 1 column 16 near "hour_second()" 
> exec ^
: string
| KILL hour_second (1)
< error
| 1064 42000 …line 1 column 16 near "hour_second (1)" 
> exec ^
: string
| SET CONFIG tikv hour_second = 1
< error
| 1064 42000 …line 1 column 27 near "hour_second = 1" 
> exec ^
: string
| SET CONFIG tikv a-hour_second = 1
< error
| 1064 42000 …line 1 column 29 near "hour_second = 1" 
> exec ^
: string
| KILL hypo()
< error
| 1064 42000 …line 1 column 9 near "hypo()" 
> exec ^
: string
| KILL hypo (1)
< error
| 1064 42000 …line 1 column 9 near "hypo (1)" 
> exec ^
: string
| SET CONFIG tikv hypo = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-hypo = 1
< ok 0 0
> exec ^
: string
| KILL identified()
< error
| 1064 42000 …line 1 column 15 near "identified()" 
> exec ^
: string
| KILL identified (1)
< error
| 1064 42000 …line 1 column 15 near "identified (1)" 
> exec ^
: string
| SET CONFIG tikv identified = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-identified = 1
< ok 0 0
> exec ^
: string
| KILL if()
< error
| 1064 42000 …line 1 column 7 near "if()" 
> exec ^
: string
| KILL if (1)
< error
| 1064 42000 …line 1 column 7 near "if (1)" 
> exec ^
: string
| SET CONFIG tikv if = 1
< error
| 1064 42000 …line 1 column 18 near "if = 1" 
> exec ^
: string
| SET CONFIG tikv a-if = 1
< error
| 1064 42000 …line 1 column 20 near "if = 1" 
> exec ^
: string
| KILL ignore()
< error
| 1064 42000 …line 1 column 11 near "ignore()" 
> exec ^
: string
| KILL ignore (1)
< error
| 1064 42000 …line 1 column 11 near "ignore (1)" 
> exec ^
: string
| SET CONFIG tikv ignore = 1
< error
| 1064 42000 …line 1 column 22 near "ignore = 1" 
> exec ^
: string
| SET CONFIG tikv a-ignore = 1
< error
| 1064 42000 …line 1 column 24 near "ignore = 1" 
> exec ^
: string
| KILL ignore_server_ids()
< ok 0 0
> exec ^
: string
| KILL ignore_server_ids (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv ignore_server_ids = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ignore_server_ids = 1
< ok 0 0
> exec ^
: string
| KILL ignore_stats()
< error
| 1064 42000 …line 1 column 17 near "ignore_stats()" 
> exec ^
: string
| KILL ignore_stats (1)
< error
| 1064 42000 …line 1 column 17 near "ignore_stats (1)" 
> exec ^
: string
| SET CONFIG tikv ignore_stats = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ignore_stats = 1
< ok 0 0
> exec ^
: string
| KILL ilike()
< error
| 1064 42000 …line 1 column 10 near "ilike()" 
> exec ^
: string
| KILL ilike (1)
< error
| 1064 42000 …line 1 column 10 near "ilike (1)" 
> exec ^
: string
| KILL import()
< error
| 1064 42000 …line 1 column 11 near "import()" 
> exec ^
: string
| KILL import (1)
< error
| 1064 42000 …line 1 column 11 near "import (1)" 
> exec ^
: string
| SET CONFIG tikv import = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-import = 1
< ok 0 0
> exec ^
: string
| KILL imports()
< error
| 1064 42000 …line 1 column 12 near "imports()" 
> exec ^
: string
| KILL imports (1)
< error
| 1064 42000 …line 1 column 12 near "imports (1)" 
> exec ^
: string
| SET CONFIG tikv imports = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-imports = 1
< ok 0 0
> exec ^
: string
| KILL in()
< error
| 1064 42000 …line 1 column 7 near "in()" 
> exec ^
: string
| KILL in (1)
< error
| 1064 42000 …line 1 column 7 near "in (1)" 
> exec ^
: string
| SET CONFIG tikv in = 1
< error
| 1064 42000 …line 1 column 18 near "in = 1" 
> exec ^
: string
| SET CONFIG tikv a-in = 1
< error
| 1064 42000 …line 1 column 20 near "in = 1" 
> exec ^
: string
| KILL inactive()
< ok 0 0
> exec ^
: string
| KILL inactive (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv inactive = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-inactive = 1
< ok 0 0
> exec ^
: string
| KILL increment()
< error
| 1064 42000 …line 1 column 14 near "increment()" 
> exec ^
: string
| KILL increment (1)
< error
| 1064 42000 …line 1 column 14 near "increment (1)" 
> exec ^
: string
| SET CONFIG tikv increment = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-increment = 1
< ok 0 0
> exec ^
: string
| KILL incremental()
< error
| 1064 42000 …line 1 column 16 near "incremental()" 
> exec ^
: string
| KILL incremental (1)
< error
| 1064 42000 …line 1 column 16 near "incremental (1)" 
> exec ^
: string
| SET CONFIG tikv incremental = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-incremental = 1
< ok 0 0
> exec ^
: string
| KILL index()
< error
| 1064 42000 …line 1 column 10 near "index()" 
> exec ^
: string
| KILL index (1)
< error
| 1064 42000 …line 1 column 10 near "index (1)" 
> exec ^
: string
| SET CONFIG tikv index = 1
< error
| 1064 42000 …line 1 column 21 near "index = 1" 
> exec ^
: string
| SET CONFIG tikv a-index = 1
< error
| 1064 42000 …line 1 column 23 near "index = 1" 
> exec ^
: string
| KILL indexes()
< error
| 1064 42000 …line 1 column 12 near "indexes()" 
> exec ^
: string
| KILL indexes (1)
< error
| 1064 42000 …line 1 column 12 near "indexes (1)" 
> exec ^
: string
| SET CONFIG tikv indexes = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-indexes = 1
< ok 0 0
> exec ^
: string
| KILL infile()
< error
| 1064 42000 …line 1 column 11 near "infile()" 
> exec ^
: string
| KILL infile (1)
< error
| 1064 42000 …line 1 column 11 near "infile (1)" 
> exec ^
: string
| SET CONFIG tikv infile = 1
< error
| 1064 42000 …line 1 column 22 near "infile = 1" 
> exec ^
: string
| SET CONFIG tikv a-infile = 1
< error
| 1064 42000 …line 1 column 24 near "infile = 1" 
> exec ^
: string
| KILL initial()
< ok 0 0
> exec ^
: string
| KILL initial (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv initial = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-initial = 1
< ok 0 0
> exec ^
: string
| KILL initial_size()
< ok 0 0
> exec ^
: string
| KILL initial_size (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv initial_size = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-initial_size = 1
< ok 0 0
> exec ^
: string
| KILL initiate()
< ok 0 0
> exec ^
: string
| KILL initiate (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv initiate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-initiate = 1
< ok 0 0
> exec ^
: string
| KILL inner()
< error
| 1064 42000 …line 1 column 10 near "inner()" 
> exec ^
: string
| KILL inner (1)
< error
| 1064 42000 …line 1 column 10 near "inner (1)" 
> exec ^
: string
| SET CONFIG tikv inner = 1
< error
| 1064 42000 …line 1 column 21 near "inner = 1" 
> exec ^
: string
| SET CONFIG tikv a-inner = 1
< error
| 1064 42000 …line 1 column 23 near "inner = 1" 
> exec ^
: string
| KILL inout()
< error
| 1064 42000 …line 1 column 10 near "inout()" 
> exec ^
: string
| KILL inout (1)
< error
| 1064 42000 …line 1 column 10 near "inout (1)" 
> exec ^
: string
| SET CONFIG tikv inout = 1
< error
| 1064 42000 …line 1 column 21 near "inout = 1" 
> exec ^
: string
| SET CONFIG tikv a-inout = 1
< error
| 1064 42000 …line 1 column 23 near "inout = 1" 
> exec ^
: string
| KILL inplace()
< error
| 1064 42000 …line 1 column 12 near "inplace()" 
> exec ^
: string
| KILL inplace (1)
< error
| 1064 42000 …line 1 column 12 near "inplace (1)" 
> exec ^
: string
| SET CONFIG tikv inplace = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-inplace = 1
< ok 0 0
> exec ^
: string
| KILL insensitive()
< ok 0 0
> exec ^
: string
| KILL insensitive (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv insensitive = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-insensitive = 1
< ok 0 0
> exec ^
: string
| KILL insert()
< error
| 1064 42000 …line 1 column 11 near "insert()" 
> exec ^
: string
| KILL insert (1)
< error
| 1064 42000 …line 1 column 11 near "insert (1)" 
> exec ^
: string
| SET CONFIG tikv insert = 1
< error
| 1064 42000 …line 1 column 22 near "insert = 1" 
> exec ^
: string
| SET CONFIG tikv a-insert = 1
< error
| 1064 42000 …line 1 column 24 near "insert = 1" 
> exec ^
: string
| KILL insert_method()
< error
| 1064 42000 …line 1 column 18 near "insert_method()" 
> exec ^
: string
| KILL insert_method (1)
< error
| 1064 42000 …line 1 column 18 near "insert_method (1)" 
> exec ^
: string
| SET CONFIG tikv insert_method = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-insert_method = 1
< ok 0 0
> exec ^
: string
| KILL install()
< ok 0 0
> exec ^
: string
| KILL install (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv install = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-install = 1
< ok 0 0
> exec ^
: string
| KILL instance()
< error
| 1064 42000 …line 1 column 13 near "instance()" 
> exec ^
: string
| KILL instance (1)
< error
| 1064 42000 …line 1 column 13 near "instance (1)" 
> exec ^
: string
| SET CONFIG tikv instance = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-instance = 1
< ok 0 0
> exec ^
: string
| KILL instant()
< error
| 1064 42000 …line 1 column 12 near "instant()" 
> exec ^
: string
| KILL instant (1)
< error
| 1064 42000 …line 1 column 12 near "instant (1)" 
> exec ^
: string
| SET CONFIG tikv instant = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-instant = 1
< ok 0 0
> exec ^
: string
| KILL int()
< error
| 1064 42000 …line 1 column 8 near "int()" 
> exec ^
: string
| KILL int (1)
< error
| 1064 42000 …line 1 column 8 near "int (1)" 
> exec ^
: string
| SET CONFIG tikv int = 1
< error
| 1064 42000 …line 1 column 19 near "int = 1" 
> exec ^
: string
| SET CONFIG tikv a-int = 1
< error
| 1064 42000 …line 1 column 21 near "int = 1" 
> exec ^
: string
| KILL int1()
< error
| 1064 42000 …line 1 column 9 near "int1()" 
> exec ^
: string
| KILL int1 (1)
< error
| 1064 42000 …line 1 column 9 near "int1 (1)" 
> exec ^
: string
| SET CONFIG tikv int1 = 1
< error
| 1064 42000 …line 1 column 20 near "int1 = 1" 
> exec ^
: string
| SET CONFIG tikv a-int1 = 1
< error
| 1064 42000 …line 1 column 22 near "int1 = 1" 
> exec ^
: string
| KILL int2()
< error
| 1064 42000 …line 1 column 9 near "int2()" 
> exec ^
: string
| KILL int2 (1)
< error
| 1064 42000 …line 1 column 9 near "int2 (1)" 
> exec ^
: string
| SET CONFIG tikv int2 = 1
< error
| 1064 42000 …line 1 column 20 near "int2 = 1" 
> exec ^
: string
| SET CONFIG tikv a-int2 = 1
< error
| 1064 42000 …line 1 column 22 near "int2 = 1" 
> exec ^
: string
| KILL int3()
< error
| 1064 42000 …line 1 column 9 near "int3()" 
> exec ^
: string
| KILL int3 (1)
< error
| 1064 42000 …line 1 column 9 near "int3 (1)" 
> exec ^
: string
| SET CONFIG tikv int3 = 1
< error
| 1064 42000 …line 1 column 20 near "int3 = 1" 
> exec ^
: string
| SET CONFIG tikv a-int3 = 1
< error
| 1064 42000 …line 1 column 22 near "int3 = 1" 
> exec ^
: string
| KILL int4()
< error
| 1064 42000 …line 1 column 9 near "int4()" 
> exec ^
: string
| KILL int4 (1)
< error
| 1064 42000 …line 1 column 9 near "int4 (1)" 
> exec ^
: string
| SET CONFIG tikv int4 = 1
< error
| 1064 42000 …line 1 column 20 near "int4 = 1" 
> exec ^
: string
| SET CONFIG tikv a-int4 = 1
< error
| 1064 42000 …line 1 column 22 near "int4 = 1" 
> exec ^
: string
| KILL int8()
< error
| 1064 42000 …line 1 column 9 near "int8()" 
> exec ^
: string
| KILL int8 (1)
< error
| 1064 42000 …line 1 column 9 near "int8 (1)" 
> exec ^
: string
| SET CONFIG tikv int8 = 1
< error
| 1064 42000 …line 1 column 20 near "int8 = 1" 
> exec ^
: string
| SET CONFIG tikv a-int8 = 1
< error
| 1064 42000 …line 1 column 22 near "int8 = 1" 
> exec ^
: string
| KILL integer()
< error
| 1064 42000 …line 1 column 12 near "integer()" 
> exec ^
: string
| KILL integer (1)
< error
| 1064 42000 …line 1 column 12 near "integer (1)" 
> exec ^
: string
| SET CONFIG tikv integer = 1
< error
| 1064 42000 …line 1 column 23 near "integer = 1" 
> exec ^
: string
| SET CONFIG tikv a-integer = 1
< error
| 1064 42000 …line 1 column 25 near "integer = 1" 
> exec ^
: string
| KILL internal()
< error
| 1064 42000 …line 1 column 13 near "internal()" 
> exec ^
: string
| KILL internal (1)
< error
| 1064 42000 …line 1 column 13 near "internal (1)" 
> exec ^
: string
| SET CONFIG tikv internal = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-internal = 1
< ok 0 0
> exec ^
: string
| KILL intersect()
< error
| 1064 42000 …line 1 column 14 near "intersect()" 
> exec ^
: string
| KILL intersect (1)
< error
| 1064 42000 …line 1 column 14 near "intersect (1)" 
> exec ^
: string
| SET CONFIG tikv intersect = 1
< error
| 1064 42000 …line 1 column 25 near "intersect = 1" 
> exec ^
: string
| SET CONFIG tikv a-intersect = 1
< error
| 1064 42000 …line 1 column 27 near "intersect = 1" 
> exec ^
: string
| KILL interval()
< error
| 1064 42000 …line 1 column 13 near "interval()" 
> exec ^
: string
| KILL interval (1)
< error
| 1064 42000 …line 1 column 13 near "interval (1)" 
> exec ^
: string
| SET CONFIG tikv interval = 1
< error
| 1064 42000 …line 1 column 24 near "interval = 1" 
> exec ^
: string
| SET CONFIG tikv a-interval = 1
< error
| 1064 42000 …line 1 column 26 near "interval = 1" 
> exec ^
: string
| KILL into()
< error
| 1064 42000 …line 1 column 9 near "into()" 
> exec ^
: string
| KILL into (1)
< error
| 1064 42000 …line 1 column 9 near "into (1)" 
> exec ^
: string
| SET CONFIG tikv into = 1
< error
| 1064 42000 …line 1 column 20 near "into = 1" 
> exec ^
: string
| SET CONFIG tikv a-into = 1
< error
| 1064 42000 …line 1 column 22 near "into = 1" 
> exec ^
: string
| KILL invisible()
< error
| 1064 42000 …line 1 column 14 near "invisible()" 
> exec ^
: string
| KILL invisible (1)
< error
| 1064 42000 …line 1 column 14 near "invisible (1)" 
> exec ^
: string
| SET CONFIG tikv invisible = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-invisible = 1
< ok 0 0
> exec ^
: string
| KILL invoker()
< error
| 1064 42000 …line 1 column 12 near "invoker()" 
> exec ^
: string
| KILL invoker (1)
< error
| 1064 42000 …line 1 column 12 near "invoker (1)" 
> exec ^
: string
| SET CONFIG tikv invoker = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-invoker = 1
< ok 0 0
> exec ^
: string
| KILL io()
< error
| 1064 42000 …line 1 column 7 near "io()" 
> exec ^
: string
| KILL io (1)
< error
| 1064 42000 …line 1 column 7 near "io (1)" 
> exec ^
: string
| SET CONFIG tikv io = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-io = 1
< ok 0 0
> exec ^
: string
| KILL io_after_gtids()
< ok 0 0
> exec ^
: string
| KILL io_after_gtids (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv io_after_gtids = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-io_after_gtids = 1
< ok 0 0
> exec ^
: string
| KILL io_before_gtids()
< ok 0 0
> exec ^
: string
| KILL io_before_gtids (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv io_before_gtids = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-io_before_gtids = 1
< ok 0 0
> exec ^
: string
| KILL io_read_bandwidth()
< error
| 1064 42000 …line 1 column 22 near "io_read_bandwidth()" 
> exec ^
: string
| KILL io_read_bandwidth (1)
< error
| 1064 42000 …line 1 column 22 near "io_read_bandwidth (1)" 
> exec ^
: string
| SET CONFIG tikv io_read_bandwidth = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-io_read_bandwidth = 1
< ok 0 0
> exec ^
: string
| KILL io_thread()
< ok 0 0
> exec ^
: string
| KILL io_thread (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv io_thread = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-io_thread = 1
< ok 0 0
> exec ^
: string
| KILL io_write_bandwidth()
< error
| 1064 42000 …line 1 column 23 near "io_write_bandwidth()" 
> exec ^
: string
| KILL io_write_bandwidth (1)
< error
| 1064 42000 …line 1 column 23 near "io_write_bandwidth (1)" 
> exec ^
: string
| SET CONFIG tikv io_write_bandwidth = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-io_write_bandwidth = 1
< ok 0 0
> exec ^
: string
| KILL ipc()
< error
| 1064 42000 …line 1 column 8 near "ipc()" 
> exec ^
: string
| KILL ipc (1)
< error
| 1064 42000 …line 1 column 8 near "ipc (1)" 
> exec ^
: string
| SET CONFIG tikv ipc = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ipc = 1
< ok 0 0
> exec ^
: string
| KILL is()
< error
| 1064 42000 …line 1 column 7 near "is()" 
> exec ^
: string
| KILL is (1)
< error
| 1064 42000 …line 1 column 7 near "is (1)" 
> exec ^
: string
| SET CONFIG tikv is = 1
< error
| 1064 42000 …line 1 column 18 near "is = 1" 
> exec ^
: string
| SET CONFIG tikv a-is = 1
< error
| 1064 42000 …line 1 column 20 near "is = 1" 
> exec ^
: string
| KILL isolation()
< error
| 1064 42000 …line 1 column 14 near "isolation()" 
> exec ^
: string
| KILL isolation (1)
< error
| 1064 42000 …line 1 column 14 near "isolation (1)" 
> exec ^
: string
| SET CONFIG tikv isolation = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-isolation = 1
< ok 0 0
> exec ^
: string
| KILL issuer()
< error
| 1064 42000 …line 1 column 11 near "issuer()" 
> exec ^
: string
| KILL issuer (1)
< error
| 1064 42000 …line 1 column 11 near "issuer (1)" 
> exec ^
: string
| SET CONFIG tikv issuer = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-issuer = 1
< ok 0 0
> exec ^
: string
| KILL iterate()
< error
| 1064 42000 …line 1 column 12 near "iterate()" 
> exec ^
: string
| KILL iterate (1)
< error
| 1064 42000 …line 1 column 12 near "iterate (1)" 
> exec ^
: string
| SET CONFIG tikv iterate = 1
< error
| 1064 42000 …line 1 column 23 near "iterate = 1" 
> exec ^
: string
| SET CONFIG tikv a-iterate = 1
< error
| 1064 42000 …line 1 column 25 near "iterate = 1" 
> exec ^
: string
| KILL job()
< error
| 1064 42000 …line 1 column 8 near "job()" 
> exec ^
: string
| KILL job (1)
< error
| 1064 42000 …line 1 column 8 near "job (1)" 
> exec ^
: string
| SET CONFIG tikv job = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-job = 1
< ok 0 0
> exec ^
: string
| KILL jobs()
< error
| 1064 42000 …line 1 column 9 near "jobs()" 
> exec ^
: string
| KILL jobs (1)
< error
| 1064 42000 …line 1 column 9 near "jobs (1)" 
> exec ^
: string
| SET CONFIG tikv jobs = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-jobs = 1
< ok 0 0
> exec ^
: string
| KILL join()
< error
| 1064 42000 …line 1 column 9 near "join()" 
> exec ^
: string
| KILL join (1)
< error
| 1064 42000 …line 1 column 9 near "join (1)" 
> exec ^
: string
| SET CONFIG tikv join = 1
< error
| 1064 42000 …line 1 column 20 near "join = 1" 
> exec ^
: string
| SET CONFIG tikv a-join = 1
< error
| 1064 42000 …line 1 column 22 near "join = 1" 
> exec ^
: string
| KILL json()
< error
| 1064 42000 …line 1 column 9 near "json()" 
> exec ^
: string
| KILL json (1)
< error
| 1064 42000 …line 1 column 9 near "json (1)" 
> exec ^
: string
| SET CONFIG tikv json = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-json = 1
< ok 0 0
> exec ^
: string
| KILL json_arrayagg()
< error
| 1064 42000 …line 1 column 18 near "json_arrayagg()" 
> exec ^
: string
| KILL json_arrayagg (1)
< error
| 1064 42000 …line 1 column 18 near "json_arrayagg (1)" 
> exec ^
: string
| SET CONFIG tikv json_arrayagg = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-json_arrayagg = 1
< ok 0 0
> exec ^
: string
| KILL json_duality_object()
< ok 0 0
> exec ^
: string
| KILL json_duality_object (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv json_duality_object = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-json_duality_object = 1
< ok 0 0
> exec ^
: string
| KILL json_objectagg()
< error
| 1064 42000 …line 1 column 19 near "json_objectagg()" 
> exec ^
: string
| KILL json_objectagg (1)
< error
| 1064 42000 …line 1 column 19 near "json_objectagg (1)" 
> exec ^
: string
| SET CONFIG tikv json_objectagg = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-json_objectagg = 1
< ok 0 0
> exec ^
: string
| KILL json_table()
< ok 0 0
> exec ^
: string
| KILL json_table (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv json_table = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-json_table = 1
< ok 0 0
> exec ^
: string
| KILL json_value()
< ok 0 0
> exec ^
: string
| KILL json_value (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv json_value = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-json_value = 1
< ok 0 0
> exec ^
: string
| KILL key()
< error
| 1064 42000 …line 1 column 8 near "key()" 
> exec ^
: string
| KILL key (1)
< error
| 1064 42000 …line 1 column 8 near "key (1)" 
> exec ^
: string
| SET CONFIG tikv key = 1
< error
| 1064 42000 …line 1 column 19 near "key = 1" 
> exec ^
: string
| SET CONFIG tikv a-key = 1
< error
| 1064 42000 …line 1 column 21 near "key = 1" 
> exec ^
: string
| KILL key_block_size()
< error
| 1064 42000 …line 1 column 19 near "key_block_size()" 
> exec ^
: string
| KILL key_block_size (1)
< error
| 1064 42000 …line 1 column 19 near "key_block_size (1)" 
> exec ^
: string
| SET CONFIG tikv key_block_size = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-key_block_size = 1
< ok 0 0
> exec ^
: string
| KILL keyring()
< ok 0 0
> exec ^
: string
| KILL keyring (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv keyring = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-keyring = 1
< ok 0 0
> exec ^
: string
| KILL keys()
< error
| 1064 42000 …line 1 column 9 near "keys()" 
> exec ^
: string
| KILL keys (1)
< error
| 1064 42000 …line 1 column 9 near "keys (1)" 
> exec ^
: string
| SET CONFIG tikv keys = 1
< error
| 1064 42000 …line 1 column 20 near "keys = 1" 
> exec ^
: string
| SET CONFIG tikv a-keys = 1
< error
| 1064 42000 …line 1 column 22 near "keys = 1" 
> exec ^
: string
| KILL kill()
< error
| 1064 42000 …line 1 column 9 near "kill()" 
> exec ^
: string
| KILL kill (1)
< error
| 1064 42000 …line 1 column 9 near "kill (1)" 
> exec ^
: string
| SET CONFIG tikv kill = 1
< error
| 1064 42000 …line 1 column 20 near "kill = 1" 
> exec ^
: string
| SET CONFIG tikv a-kill = 1
< error
| 1064 42000 …line 1 column 22 near "kill = 1" 
> exec ^
: string
| KILL labels()
< error
| 1064 42000 …line 1 column 11 near "labels()" 
> exec ^
: string
| KILL labels (1)
< error
| 1064 42000 …line 1 column 11 near "labels (1)" 
> exec ^
: string
| SET CONFIG tikv labels = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-labels = 1
< ok 0 0
> exec ^
: string
| KILL lag()
< error
| 1064 42000 …line 1 column 8 near "lag()" 
> exec ^
: string
| KILL lag (1)
< error
| 1064 42000 …line 1 column 8 near "lag (1)" 
> exec ^
: string
| SET CONFIG tikv lag = 1
< error
| 1064 42000 …line 1 column 19 near "lag = 1" 
> exec ^
: string
| SET CONFIG tikv a-lag = 1
< error
| 1064 42000 …line 1 column 21 near "lag = 1" 
> exec ^
: string
| KILL language()
< error
| 1064 42000 …line 1 column 13 near "language()" 
> exec ^
: string
| KILL language (1)
< error
| 1064 42000 …line 1 column 13 near "language (1)" 
> exec ^
: string
| SET CONFIG tikv language = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-language = 1
< ok 0 0
> exec ^
: string
| KILL last()
< error
| 1064 42000 …line 1 column 9 near "last()" 
> exec ^
: string
| KILL last (1)
< error
| 1064 42000 …line 1 column 9 near "last (1)" 
> exec ^
: string
| SET CONFIG tikv last = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-last = 1
< ok 0 0
> exec ^
: string
| KILL last_backup()
< error
| 1064 42000 …line 1 column 16 near "last_backup()" 
> exec ^
: string
| KILL last_backup (1)
< error
| 1064 42000 …line 1 column 16 near "last_backup (1)" 
> exec ^
: string
| SET CONFIG tikv last_backup = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-last_backup = 1
< ok 0 0
> exec ^
: string
| KILL last_value()
< error
| 1064 42000 …line 1 column 15 near "last_value()" 
> exec ^
: string
| KILL last_value (1)
< error
| 1064 42000 …line 1 column 15 near "last_value (1)" 
> exec ^
: string
| SET CONFIG tikv last_value = 1
< error
| 1064 42000 …line 1 column 26 near "last_value = 1" 
> exec ^
: string
| SET CONFIG tikv a-last_value = 1
< error
| 1064 42000 …line 1 column 28 near "last_value = 1" 
> exec ^
: string
| KILL lastval()
< error
| 1064 42000 …line 1 column 12 near "lastval()" 
> exec ^
: string
| KILL lastval (1)
< error
| 1064 42000 …line 1 column 12 near "lastval (1)" 
> exec ^
: string
| SET CONFIG tikv lastval = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-lastval = 1
< ok 0 0
> exec ^
: string
| KILL lateral()
< ok 0 0
> exec ^
: string
| KILL lateral (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv lateral = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-lateral = 1
< ok 0 0
> exec ^
: string
| KILL lead()
< error
| 1064 42000 …line 1 column 9 near "lead()" 
> exec ^
: string
| KILL lead (1)
< error
| 1064 42000 …line 1 column 9 near "lead (1)" 
> exec ^
: string
| SET CONFIG tikv lead = 1
< error
| 1064 42000 …line 1 column 20 near "lead = 1" 
> exec ^
: string
| SET CONFIG tikv a-lead = 1
< error
| 1064 42000 …line 1 column 22 near "lead = 1" 
> exec ^
: string
| KILL leader()
< error
| 1064 42000 …line 1 column 11 near "leader()" 
> exec ^
: string
| KILL leader (1)
< error
| 1064 42000 …line 1 column 11 near "leader (1)" 
> exec ^
: string
| SET CONFIG tikv leader = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-leader = 1
< ok 0 0
> exec ^
: string
| KILL leader_constraints()
< error
| 1064 42000 …line 1 column 23 near "leader_constraints()" 
> exec ^
: string
| KILL leader_constraints (1)
< error
| 1064 42000 …line 1 column 23 near "leader_constraints (1)" 
> exec ^
: string
| SET CONFIG tikv leader_constraints = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-leader_constraints = 1
< ok 0 0
> exec ^
: string
| KILL leading()
< error
| 1064 42000 …line 1 column 12 near "leading()" 
> exec ^
: string
| KILL leading (1)
< error
| 1064 42000 …line 1 column 12 near "leading (1)" 
> exec ^
: string
| SET CONFIG tikv leading = 1
< error
| 1064 42000 …line 1 column 23 near "leading = 1" 
> exec ^
: string
| SET CONFIG tikv a-leading = 1
< error
| 1064 42000 …line 1 column 25 near "leading = 1" 
> exec ^
: string
| KILL learner()
< error
| 1064 42000 …line 1 column 12 near "learner()" 
> exec ^
: string
| KILL learner (1)
< error
| 1064 42000 …line 1 column 12 near "learner (1)" 
> exec ^
: string
| SET CONFIG tikv learner = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-learner = 1
< ok 0 0
> exec ^
: string
| KILL learner_constraints()
< error
| 1064 42000 …line 1 column 24 near "learner_constraints()" 
> exec ^
: string
| KILL learner_constraints (1)
< error
| 1064 42000 …line 1 column 24 near "learner_constraints (1)" 
> exec ^
: string
| SET CONFIG tikv learner_constraints = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-learner_constraints = 1
< ok 0 0
> exec ^
: string
| KILL learners()
< error
| 1064 42000 …line 1 column 13 near "learners()" 
> exec ^
: string
| KILL learners (1)
< error
| 1064 42000 …line 1 column 13 near "learners (1)" 
> exec ^
: string
| SET CONFIG tikv learners = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-learners = 1
< ok 0 0
> exec ^
: string
| KILL least()
< ok 0 0
> exec ^
: string
| KILL least (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv least = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-least = 1
< ok 0 0
> exec ^
: string
| KILL leave()
< error
| 1064 42000 …line 1 column 10 near "leave()" 
> exec ^
: string
| KILL leave (1)
< error
| 1064 42000 …line 1 column 10 near "leave (1)" 
> exec ^
: string
| SET CONFIG tikv leave = 1
< error
| 1064 42000 …line 1 column 21 near "leave = 1" 
> exec ^
: string
| SET CONFIG tikv a-leave = 1
< error
| 1064 42000 …line 1 column 23 near "leave = 1" 
> exec ^
: string
| KILL leaves()
< ok 0 0
> exec ^
: string
| KILL leaves (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv leaves = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-leaves = 1
< ok 0 0
> exec ^
: string
| KILL left()
< error
| 1064 42000 …line 1 column 9 near "left()" 
> exec ^
: string
| KILL left (1)
< error
| 1064 42000 …line 1 column 9 near "left (1)" 
> exec ^
: string
| SET CONFIG tikv left = 1
< error
| 1064 42000 …line 1 column 20 near "left = 1" 
> exec ^
: string
| SET CONFIG tikv a-left = 1
< error
| 1064 42000 …line 1 column 22 near "left = 1" 
> exec ^
: string
| KILL less()
< error
| 1064 42000 …line 1 column 9 near "less()" 
> exec ^
: string
| KILL less (1)
< error
| 1064 42000 …line 1 column 9 near "less (1)" 
> exec ^
: string
| SET CONFIG tikv less = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-less = 1
< ok 0 0
> exec ^
: string
| KILL level()
< error
| 1064 42000 …line 1 column 10 near "level()" 
> exec ^
: string
| KILL level (1)
< error
| 1064 42000 …line 1 column 10 near "level (1)" 
> exec ^
: string
| SET CONFIG tikv level = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-level = 1
< ok 0 0
> exec ^
: string
| KILL library()
< ok 0 0
> exec ^
: string
| KILL library (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv library = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-library = 1
< ok 0 0
> exec ^
: string
| KILL like()
< error
| 1064 42000 …line 1 column 9 near "like()" 
> exec ^
: string
| KILL like (1)
< error
| 1064 42000 …line 1 column 9 near "like (1)" 
> exec ^
: string
| SET CONFIG tikv like = 1
< error
| 1064 42000 …line 1 column 20 near "like = 1" 
> exec ^
: string
| SET CONFIG tikv a-like = 1
< error
| 1064 42000 …line 1 column 22 near "like = 1" 
> exec ^
: string
| KILL limit()
< error
| 1064 42000 …line 1 column 10 near "limit()" 
> exec ^
: string
| KILL limit (1)
< error
| 1064 42000 …line 1 column 10 near "limit (1)" 
> exec ^
: string
| SET CONFIG tikv limit = 1
< error
| 1064 42000 …line 1 column 21 near "limit = 1" 
> exec ^
: string
| SET CONFIG tikv a-limit = 1
< error
| 1064 42000 …line 1 column 23 near "limit = 1" 
> exec ^
: string
| KILL linear()
< error
| 1064 42000 …line 1 column 11 near "linear()" 
> exec ^
: string
| KILL linear (1)
< error
| 1064 42000 …line 1 column 11 near "linear (1)" 
> exec ^
: string
| SET CONFIG tikv linear = 1
< error
| 1064 42000 …line 1 column 22 near "linear = 1" 
> exec ^
: string
| SET CONFIG tikv a-linear = 1
< error
| 1064 42000 …line 1 column 24 near "linear = 1" 
> exec ^
: string
| KILL lines()
< error
| 1064 42000 …line 1 column 10 near "lines()" 
> exec ^
: string
| KILL lines (1)
< error
| 1064 42000 …line 1 column 10 near "lines (1)" 
> exec ^
: string
| SET CONFIG tikv lines = 1
< error
| 1064 42000 …line 1 column 21 near "lines = 1" 
> exec ^
: string
| SET CONFIG tikv a-lines = 1
< error
| 1064 42000 …line 1 column 23 near "lines = 1" 
> exec ^
: string
| KILL linestring()
< ok 0 0
> exec ^
: string
| KILL linestring (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv linestring = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-linestring = 1
< ok 0 0
> exec ^
: string
| KILL list()
< error
| 1064 42000 …line 1 column 9 near "list()" 
> exec ^
: string
| KILL list (1)
< error
| 1064 42000 …line 1 column 9 near "list (1)" 
> exec ^
: string
| SET CONFIG tikv list = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-list = 1
< ok 0 0
> exec ^
: string
| KILL load()
< error
| 1064 42000 …line 1 column 9 near "load()" 
> exec ^
: string
| KILL load (1)
< error
| 1064 42000 …line 1 column 9 near "load (1)" 
> exec ^
: string
| SET CONFIG tikv load = 1
< error
| 1064 42000 …line 1 column 20 near "load = 1" 
> exec ^
: string
| SET CONFIG tikv a-load = 1
< error
| 1064 42000 …line 1 column 22 near "load = 1" 
> exec ^
: string
| KILL load_stats()
< error
| 1064 42000 …line 1 column 15 near "load_stats()" 
> exec ^
: string
| KILL load_stats (1)
< error
| 1064 42000 …line 1 column 15 near "load_stats (1)" 
> exec ^
: string
| SET CONFIG tikv load_stats = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-load_stats = 1
< ok 0 0
> exec ^
: string
| KILL local()
< error
| 1064 42000 …line 1 column 10 near "local()" 
> exec ^
: string
| KILL local (1)
< error
| 1064 42000 …line 1 column 10 near "local (1)" 
> exec ^
: string
| SET CONFIG tikv local = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-local = 1
< ok 0 0
> exec ^
: string
| KILL localtime()
< error
| 1064 42000 …line 1 column 14 near "localtime()" 
> exec ^
: string
| KILL localtime (1)
< error
| 1064 42000 …line 1 column 14 near "localtime (1)" 
> exec ^
: string
| SET CONFIG tikv localtime = 1
< error
| 1064 42000 …line 1 column 25 near "localtime = 1" 
> exec ^
: string
| SET CONFIG tikv a-localtime = 1
< error
| 1064 42000 …line 1 column 27 near "localtime = 1" 
> exec ^
: string
| KILL localtimestamp()
< error
| 1064 42000 …line 1 column 19 near "localtimestamp()" 
> exec ^
: string
| KILL localtimestamp (1)
< error
| 1064 42000 …line 1 column 19 near "localtimestamp (1)" 
> exec ^
: string
| SET CONFIG tikv localtimestamp = 1
< error
| 1064 42000 …line 1 column 30 near "localtimestamp = 1" 
> exec ^
: string
| SET CONFIG tikv a-localtimestamp = 1
< error
| 1064 42000 …line 1 column 32 near "localtimestamp = 1" 
> exec ^
: string
| KILL location()
< error
| 1064 42000 …line 1 column 13 near "location()" 
> exec ^
: string
| KILL location (1)
< error
| 1064 42000 …line 1 column 13 near "location (1)" 
> exec ^
: string
| SET CONFIG tikv location = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-location = 1
< ok 0 0
> exec ^
: string
| KILL lock()
< error
| 1064 42000 …line 1 column 9 near "lock()" 
> exec ^
: string
| KILL lock (1)
< error
| 1064 42000 …line 1 column 9 near "lock (1)" 
> exec ^
: string
| SET CONFIG tikv lock = 1
< error
| 1064 42000 …line 1 column 20 near "lock = 1" 
> exec ^
: string
| SET CONFIG tikv a-lock = 1
< error
| 1064 42000 …line 1 column 22 near "lock = 1" 
> exec ^
: string
| KILL locked()
< error
| 1064 42000 …line 1 column 11 near "locked()" 
> exec ^
: string
| KILL locked (1)
< error
| 1064 42000 …line 1 column 11 near "locked (1)" 
> exec ^
: string
| SET CONFIG tikv locked = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-locked = 1
< ok 0 0
> exec ^
: string
| KILL locks()
< ok 0 0
> exec ^
: string
| KILL locks (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv locks = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-locks = 1
< ok 0 0
> exec ^
: string
| KILL log()
< error
| 1064 42000 …line 1 column 8 near "log()" 
> exec ^
: string
| KILL log (1)
< error
| 1064 42000 …line 1 column 8 near "log (1)" 
> exec ^
: string
| SET CONFIG tikv log = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-log = 1
< ok 0 0
> exec ^
: string
| KILL logfile()
< ok 0 0
> exec ^
: string
| KILL logfile (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv logfile = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-logfile = 1
< ok 0 0
> exec ^
: string
| KILL logs()
< error
| 1064 42000 …line 1 column 9 near "logs()" 
> exec ^
: string
| KILL logs (1)
< error
| 1064 42000 …line 1 column 9 near "logs (1)" 
> exec ^
: string
| SET CONFIG tikv logs = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-logs = 1
< ok 0 0
> exec ^
: string
| KILL long()
< error
| 1064 42000 …line 1 column 9 near "long()" 
> exec ^
: string
| KILL long (1)
< error
| 1064 42000 …line 1 column 9 near "long (1)" 
> exec ^
: string
| SET CONFIG tikv long = 1
< error
| 1064 42000 …line 1 column 20 near "long = 1" 
> exec ^
: string
| SET CONFIG tikv a-long = 1
< error
| 1064 42000 …line 1 column 22 near "long = 1" 
> exec ^
: string
| KILL longblob()
< error
| 1064 42000 …line 1 column 13 near "longblob()" 
> exec ^
: string
| KILL longblob (1)
< error
| 1064 42000 …line 1 column 13 near "longblob (1)" 
> exec ^
: string
| SET CONFIG tikv longblob = 1
< error
| 1064 42000 …line 1 column 24 near "longblob = 1" 
> exec ^
: string
| SET CONFIG tikv a-longblob = 1
< error
| 1064 42000 …line 1 column 26 near "longblob = 1" 
> exec ^
: string
| KILL longtext()
< error
| 1064 42000 …line 1 column 13 near "longtext()" 
> exec ^
: string
| KILL longtext (1)
< error
| 1064 42000 …line 1 column 13 near "longtext (1)" 
> exec ^
: string
| SET CONFIG tikv longtext = 1
< error
| 1064 42000 …line 1 column 24 near "longtext = 1" 
> exec ^
: string
| SET CONFIG tikv a-longtext = 1
< error
| 1064 42000 …line 1 column 26 near "longtext = 1" 
> exec ^
: string
| KILL loop()
< ok 0 0
> exec ^
: string
| KILL loop (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv loop = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-loop = 1
< ok 0 0
> exec ^
: string
| KILL low()
< error
| 1064 42000 …line 1 column 8 near "low()" 
> exec ^
: string
| KILL low (1)
< error
| 1064 42000 …line 1 column 8 near "low (1)" 
> exec ^
: string
| SET CONFIG tikv low = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-low = 1
< ok 0 0
> exec ^
: string
| KILL low_priority()
< error
| 1064 42000 …line 1 column 17 near "low_priority()" 
> exec ^
: string
| KILL low_priority (1)
< error
| 1064 42000 …line 1 column 17 near "low_priority (1)" 
> exec ^
: string
| SET CONFIG tikv low_priority = 1
< error
| 1064 42000 …line 1 column 28 near "low_priority = 1" 
> exec ^
: string
| SET CONFIG tikv a-low_priority = 1
< error
| 1064 42000 …line 1 column 30 near "low_priority = 1" 
> exec ^
: string
| KILL manual()
< ok 0 0
> exec ^
: string
| KILL manual (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv manual = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-manual = 1
< ok 0 0
> exec ^
: string
| KILL master()
< error
| 1064 42000 …line 1 column 11 near "master()" 
> exec ^
: string
| KILL master (1)
< error
| 1064 42000 …line 1 column 11 near "master (1)" 
> exec ^
: string
| SET CONFIG tikv master = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-master = 1
< ok 0 0
> exec ^
: string
| KILL master_log_file()
< ok 0 0
> exec ^
: string
| KILL master_log_file (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv master_log_file = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-master_log_file = 1
< ok 0 0
> exec ^
: string
| KILL match()
< error
| 1064 42000 …line 1 column 10 near "match()" 
> exec ^
: string
| KILL match (1)
< error
| 1064 42000 …line 1 column 10 near "match (1)" 
> exec ^
: string
| SET CONFIG tikv match = 1
< error
| 1064 42000 …line 1 column 21 near "match = 1" 
> exec ^
: string
| SET CONFIG tikv a-match = 1
< error
| 1064 42000 …line 1 column 23 near "match = 1" 
> exec ^
: string
| KILL materialized()
< ok 0 0
> exec ^
: string
| KILL materialized (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv materialized = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-materialized = 1
< ok 0 0
> exec ^
: string
| KILL max()
< error
| 1064 42000 …line 1 column 8 near "max()" 
> exec ^
: string
| KILL max (1)
< error
| 1064 42000 …line 1 column 8 near "max (1)" 
> exec ^
: string
| SET CONFIG tikv max = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-max = 1
< ok 0 0
> exec ^
: string
| KILL max_connections_per_hour()
< error
| 1064 42000 …line 1 column 29 near "max_connections_per_hour()" 
> exec ^
: string
| KILL max_connections_per_hour (1)
< error
| 1064 42000 …line 1 column 29 near "max_connections_per_hour (1)" 
> exec ^
: string
| SET CONFIG tikv max_connections_per_hour = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-max_connections_per_hour = 1
< ok 0 0
> exec ^
: string
| KILL max_idxnum()
< error
| 1064 42000 …line 1 column 15 near "max_idxnum()" 
> exec ^
: string
| KILL max_idxnum (1)
< error
| 1064 42000 …line 1 column 15 near "max_idxnum (1)" 
> exec ^
: string
| SET CONFIG tikv max_idxnum = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-max_idxnum = 1
< ok 0 0
> exec ^
: string
| KILL max_minutes()
< error
| 1064 42000 …line 1 column 16 near "max_minutes()" 
> exec ^
: string
| KILL max_minutes (1)
< error
| 1064 42000 …line 1 column 16 near "max_minutes (1)" 
> exec ^
: string
| SET CONFIG tikv max_minutes = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-max_minutes = 1
< ok 0 0
> exec ^
: string
| KILL max_queries_per_hour()
< error
| 1064 42000 …line 1 column 25 near "max_queries_per_hour()" 
> exec ^
: string
| KILL max_queries_per_hour (1)
< error
| 1064 42000 …line 1 column 25 near "max_queries_per_hour (1)" 
> exec ^
: string
| SET CONFIG tikv max_queries_per_hour = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-max_queries_per_hour = 1
< ok 0 0
> exec ^
: string
| KILL max_rows()
< error
| 1064 42000 …line 1 column 13 near "max_rows()" 
> exec ^
: string
| KILL max_rows (1)
< error
| 1064 42000 …line 1 column 13 near "max_rows (1)" 
> exec ^
: string
| SET CONFIG tikv max_rows = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-max_rows = 1
< ok 0 0
> exec ^
: string
| KILL max_size()
< ok 0 0
> exec ^
: string
| KILL max_size (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv max_size = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-max_size = 1
< ok 0 0
> exec ^
: string
| KILL max_updates_per_hour()
< error
| 1064 42000 …line 1 column 25 near "max_updates_per_hour()" 
> exec ^
: string
| KILL max_updates_per_hour (1)
< error
| 1064 42000 …line 1 column 25 near "max_updates_per_hour (1)" 
> exec ^
: string
| SET CONFIG tikv max_updates_per_hour = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-max_updates_per_hour = 1
< ok 0 0
> exec ^
: string
| KILL max_user_connections()
< error
| 1064 42000 …line 1 column 25 near "max_user_connections()" 
> exec ^
: string
| KILL max_user_connections (1)
< error
| 1064 42000 …line 1 column 25 near "max_user_connections (1)" 
> exec ^
: string
| SET CONFIG tikv max_user_connections = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-max_user_connections = 1
< ok 0 0
> exec ^
: string
| KILL maxvalue()
< error
| 1064 42000 …line 1 column 13 near "maxvalue()" 
> exec ^
: string
| KILL maxvalue (1)
< error
| 1064 42000 …line 1 column 13 near "maxvalue (1)" 
> exec ^
: string
| SET CONFIG tikv maxvalue = 1
< error
| 1064 42000 …line 1 column 24 near "maxvalue = 1" 
> exec ^
: string
| SET CONFIG tikv a-maxvalue = 1
< error
| 1064 42000 …line 1 column 26 near "maxvalue = 1" 
> exec ^
: string
| KILL mb()
< error
| 1064 42000 …line 1 column 7 near "mb()" 
> exec ^
: string
| KILL mb (1)
< error
| 1064 42000 …line 1 column 7 near "mb (1)" 
> exec ^
: string
| SET CONFIG tikv mb = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-mb = 1
< ok 0 0
> exec ^
: string
| KILL medium()
< error
| 1064 42000 …line 1 column 11 near "medium()" 
> exec ^
: string
| KILL medium (1)
< error
| 1064 42000 …line 1 column 11 near "medium (1)" 
> exec ^
: string
| SET CONFIG tikv medium = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-medium = 1
< ok 0 0
> exec ^
: string
| KILL mediumblob()
< error
| 1064 42000 …line 1 column 15 near "mediumblob()" 
> exec ^
: string
| KILL mediumblob (1)
< error
| 1064 42000 …line 1 column 15 near "mediumblob (1)" 
> exec ^
: string
| SET CONFIG tikv mediumblob = 1
< error
| 1064 42000 …line 1 column 26 near "mediumblob = 1" 
> exec ^
: string
| SET CONFIG tikv a-mediumblob = 1
< error
| 1064 42000 …line 1 column 28 near "mediumblob = 1" 
> exec ^
: string
| KILL mediumint()
< error
| 1064 42000 …line 1 column 14 near "mediumint()" 
> exec ^
: string
| KILL mediumint (1)
< error
| 1064 42000 …line 1 column 14 near "mediumint (1)" 
> exec ^
: string
| SET CONFIG tikv mediumint = 1
< error
| 1064 42000 …line 1 column 25 near "mediumint = 1" 
> exec ^
: string
| SET CONFIG tikv a-mediumint = 1
< error
| 1064 42000 …line 1 column 27 near "mediumint = 1" 
> exec ^
: string
| KILL mediumtext()
< error
| 1064 42000 …line 1 column 15 near "mediumtext()" 
> exec ^
: string
| KILL mediumtext (1)
< error
| 1064 42000 …line 1 column 15 near "mediumtext (1)" 
> exec ^
: string
| SET CONFIG tikv mediumtext = 1
< error
| 1064 42000 …line 1 column 26 near "mediumtext = 1" 
> exec ^
: string
| SET CONFIG tikv a-mediumtext = 1
< error
| 1064 42000 …line 1 column 28 near "mediumtext = 1" 
> exec ^
: string
| KILL member()
< error
| 1064 42000 …line 1 column 11 near "member()" 
> exec ^
: string
| KILL member (1)
< error
| 1064 42000 …line 1 column 11 near "member (1)" 
> exec ^
: string
| SET CONFIG tikv member = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-member = 1
< ok 0 0
> exec ^
: string
| KILL memory()
< error
| 1064 42000 …line 1 column 11 near "memory()" 
> exec ^
: string
| KILL memory (1)
< error
| 1064 42000 …line 1 column 11 near "memory (1)" 
> exec ^
: string
| SET CONFIG tikv memory = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-memory = 1
< ok 0 0
> exec ^
: string
| KILL merge()
< error
| 1064 42000 …line 1 column 10 near "merge()" 
> exec ^
: string
| KILL merge (1)
< error
| 1064 42000 …line 1 column 10 near "merge (1)" 
> exec ^
: string
| SET CONFIG tikv merge = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-merge = 1
< ok 0 0
> exec ^
: string
| KILL message_text()
< ok 0 0
> exec ^
: string
| KILL message_text (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv message_text = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-message_text = 1
< ok 0 0
> exec ^
: string
| KILL metadata()
< error
| 1064 42000 …line 1 column 13 near "metadata()" 
> exec ^
: string
| KILL metadata (1)
< error
| 1064 42000 …line 1 column 13 near "metadata (1)" 
> exec ^
: string
| SET CONFIG tikv metadata = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-metadata = 1
< ok 0 0
> exec ^
: string
| KILL microsecond()
< error
| 1064 42000 …line 1 column 16 near "microsecond()" 
> exec ^
: string
| KILL microsecond (1)
< error
| 1064 42000 …line 1 column 16 near "microsecond (1)" 
> exec ^
: string
| SET CONFIG tikv microsecond = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-microsecond = 1
< ok 0 0
> exec ^
: string
| KILL mid()
< error
| 1064 42000 …line 1 column 8 near "mid()" 
> exec ^
: string
| KILL mid (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv mid = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-mid = 1
< ok 0 0
> exec ^
: string
| KILL middleint()
< error
| 1064 42000 …line 1 column 14 near "middleint()" 
> exec ^
: string
| KILL middleint (1)
< error
| 1064 42000 …line 1 column 14 near "middleint (1)" 
> exec ^
: string
| SET CONFIG tikv middleint = 1
< error
| 1064 42000 …line 1 column 25 near "middleint = 1" 
> exec ^
: string
| SET CONFIG tikv a-middleint = 1
< error
| 1064 42000 …line 1 column 27 near "middleint = 1" 
> exec ^
: string
| KILL migrate()
< ok 0 0
> exec ^
: string
| KILL migrate (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv migrate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-migrate = 1
< ok 0 0
> exec ^
: string
| KILL min()
< error
| 1064 42000 …line 1 column 8 near "min()" 
> exec ^
: string
| KILL min (1)
< error
| 1064 42000 …line 1 column 8 near "min (1)" 
> exec ^
: string
| SET CONFIG tikv min = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-min = 1
< ok 0 0
> exec ^
: string
| KILL min_rows()
< error
| 1064 42000 …line 1 column 13 near "min_rows()" 
> exec ^
: string
| KILL min_rows (1)
< error
| 1064 42000 …line 1 column 13 near "min_rows (1)" 
> exec ^
: string
| SET CONFIG tikv min_rows = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-min_rows = 1
< ok 0 0
> exec ^
: string
| KILL minute()
< error
| 1064 42000 …line 1 column 11 near "minute()" 
> exec ^
: string
| KILL minute (1)
< error
| 1064 42000 …line 1 column 11 near "minute (1)" 
> exec ^
: string
| SET CONFIG tikv minute = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-minute = 1
< ok 0 0
> exec ^
: string
| KILL minute_microsecond()
< error
| 1064 42000 …line 1 column 23 near "minute_microsecond()" 
> exec ^
: string
| KILL minute_microsecond (1)
< error
| 1064 42000 …line 1 column 23 near "minute_microsecond (1)" 
> exec ^
: string
| SET CONFIG tikv minute_microsecond = 1
< error
| 1064 42000 …line 1 column 34 near "minute_microsecond = 1" 
> exec ^
: string
| SET CONFIG tikv a-minute_microsecond = 1
< error
| 1064 42000 …line 1 column 36 near "minute_microsecond = 1" 
> exec ^
: string
| KILL minute_second()
< error
| 1064 42000 …line 1 column 18 near "minute_second()" 
> exec ^
: string
| KILL minute_second (1)
< error
| 1064 42000 …line 1 column 18 near "minute_second (1)" 
> exec ^
: string
| SET CONFIG tikv minute_second = 1
< error
| 1064 42000 …line 1 column 29 near "minute_second = 1" 
> exec ^
: string
| SET CONFIG tikv a-minute_second = 1
< error
| 1064 42000 …line 1 column 31 near "minute_second = 1" 
> exec ^
: string
| KILL minvalue()
< error
| 1064 42000 …line 1 column 13 near "minvalue()" 
> exec ^
: string
| KILL minvalue (1)
< error
| 1064 42000 …line 1 column 13 near "minvalue (1)" 
> exec ^
: string
| SET CONFIG tikv minvalue = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-minvalue = 1
< ok 0 0
> exec ^
: string
| KILL mod()
< error
| 1064 42000 …line 1 column 8 near "mod()" 
> exec ^
: string
| KILL mod (1)
< error
| 1064 42000 …line 1 column 8 near "mod (1)" 
> exec ^
: string
| SET CONFIG tikv mod = 1
< error
| 1064 42000 …line 1 column 19 near "mod = 1" 
> exec ^
: string
| SET CONFIG tikv a-mod = 1
< error
| 1064 42000 …line 1 column 21 near "mod = 1" 
> exec ^
: string
| KILL mode()
< error
| 1064 42000 …line 1 column 9 near "mode()" 
> exec ^
: string
| KILL mode (1)
< error
| 1064 42000 …line 1 column 9 near "mode (1)" 
> exec ^
: string
| SET CONFIG tikv mode = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-mode = 1
< ok 0 0
> exec ^
: string
| KILL modifies()
< ok 0 0
> exec ^
: string
| KILL modifies (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv modifies = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-modifies = 1
< ok 0 0
> exec ^
: string
| KILL modify()
< error
| 1064 42000 …line 1 column 11 near "modify()" 
> exec ^
: string
| KILL modify (1)
< error
| 1064 42000 …line 1 column 11 near "modify (1)" 
> exec ^
: string
| SET CONFIG tikv modify = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-modify = 1
< ok 0 0
> exec ^
: string
| KILL month()
< error
| 1064 42000 …line 1 column 10 near "month()" 
> exec ^
: string
| KILL month (1)
< error
| 1064 42000 …line 1 column 10 near "month (1)" 
> exec ^
: string
| SET CONFIG tikv month = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-month = 1
< ok 0 0
> exec ^
: string
| KILL multilinestring()
< ok 0 0
> exec ^
: string
| KILL multilinestring (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv multilinestring = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-multilinestring = 1
< ok 0 0
> exec ^
: string
| KILL multipoint()
< ok 0 0
> exec ^
: string
| KILL multipoint (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv multipoint = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-multipoint = 1
< ok 0 0
> exec ^
: string
| KILL multipolygon()
< ok 0 0
> exec ^
: string
| KILL multipolygon (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv multipolygon = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-multipolygon = 1
< ok 0 0
> exec ^
: string
| KILL mutex()
< ok 0 0
> exec ^
: string
| KILL mutex (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv mutex = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-mutex = 1
< ok 0 0
> exec ^
: string
| KILL mysql_errno()
< ok 0 0
> exec ^
: string
| KILL mysql_errno (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv mysql_errno = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-mysql_errno = 1
< ok 0 0
> exec ^
: string
| KILL name()
< ok 0 0
> exec ^
: string
| KILL name (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv name = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-name = 1
< ok 0 0
> exec ^
: string
| KILL names()
< error
| 1064 42000 …line 1 column 10 near "names()" 
> exec ^
: string
| KILL names (1)
< error
| 1064 42000 …line 1 column 10 near "names (1)" 
> exec ^
: string
| SET CONFIG tikv names = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-names = 1
< ok 0 0
> exec ^
: string
| KILL national()
< error
| 1064 42000 …line 1 column 13 near "national()" 
> exec ^
: string
| KILL national (1)
< error
| 1064 42000 …line 1 column 13 near "national (1)" 
> exec ^
: string
| SET CONFIG tikv national = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-national = 1
< ok 0 0
> exec ^
: string
| KILL natural()
< error
| 1064 42000 …line 1 column 12 near "natural()" 
> exec ^
: string
| KILL natural (1)
< error
| 1064 42000 …line 1 column 12 near "natural (1)" 
> exec ^
: string
| SET CONFIG tikv natural = 1
< error
| 1064 42000 …line 1 column 23 near "natural = 1" 
> exec ^
: string
| SET CONFIG tikv a-natural = 1
< error
| 1064 42000 …line 1 column 25 near "natural = 1" 
> exec ^
: string
| KILL nchar()
< error
| 1064 42000 …line 1 column 10 near "nchar()" 
> exec ^
: string
| KILL nchar (1)
< error
| 1064 42000 …line 1 column 10 near "nchar (1)" 
> exec ^
: string
| SET CONFIG tikv nchar = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nchar = 1
< ok 0 0
> exec ^
: string
| KILL ndb()
< ok 0 0
> exec ^
: string
| KILL ndb (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv ndb = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ndb = 1
< ok 0 0
> exec ^
: string
| KILL ndbcluster()
< ok 0 0
> exec ^
: string
| KILL ndbcluster (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv ndbcluster = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ndbcluster = 1
< ok 0 0
> exec ^
: string
| KILL nested()
< ok 0 0
> exec ^
: string
| KILL nested (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv nested = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nested = 1
< ok 0 0
> exec ^
: string
| KILL network_namespace()
< ok 0 0
> exec ^
: string
| KILL network_namespace (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv network_namespace = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-network_namespace = 1
< ok 0 0
> exec ^
: string
| KILL never()
< error
| 1064 42000 …line 1 column 10 near "never()" 
> exec ^
: string
| KILL never (1)
< error
| 1064 42000 …line 1 column 10 near "never (1)" 
> exec ^
: string
| SET CONFIG tikv never = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-never = 1
< ok 0 0
> exec ^
: string
| KILL new()
< ok 0 0
> exec ^
: string
| KILL new (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv new = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-new = 1
< ok 0 0
> exec ^
: string
| KILL next()
< error
| 1064 42000 …line 1 column 9 near "next()" 
> exec ^
: string
| KILL next (1)
< error
| 1064 42000 …line 1 column 9 near "next (1)" 
> exec ^
: string
| SET CONFIG tikv next = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-next = 1
< ok 0 0
> exec ^
: string
| KILL next_row_id()
< error
| 1064 42000 …line 1 column 16 near "next_row_id()" 
> exec ^
: string
| KILL next_row_id (1)
< error
| 1064 42000 …line 1 column 16 near "next_row_id (1)" 
> exec ^
: string
| SET CONFIG tikv next_row_id = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-next_row_id = 1
< ok 0 0
> exec ^
: string
| KILL nextval()
< error
| 1064 42000 …line 1 column 12 near "nextval()" 
> exec ^
: string
| KILL nextval (1)
< error
| 1064 42000 …line 1 column 12 near "nextval (1)" 
> exec ^
: string
| SET CONFIG tikv nextval = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nextval = 1
< ok 0 0
> exec ^
: string
| KILL no()
< error
| 1064 42000 …line 1 column 7 near "no()" 
> exec ^
: string
| KILL no (1)
< error
| 1064 42000 …line 1 column 7 near "no (1)" 
> exec ^
: string
| SET CONFIG tikv no = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-no = 1
< ok 0 0
> exec ^
: string
| KILL no_wait()
< ok 0 0
> exec ^
: string
| KILL no_wait (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv no_wait = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-no_wait = 1
< ok 0 0
> exec ^
: string
| KILL no_write_to_binlog()
< error
| 1064 42000 …line 1 column 23 near "no_write_to_binlog()" 
> exec ^
: string
| KILL no_write_to_binlog (1)
< error
| 1064 42000 …line 1 column 23 near "no_write_to_binlog (1)" 
> exec ^
: string
| SET CONFIG tikv no_write_to_binlog = 1
< error
| 1064 42000 …line 1 column 34 near "no_write_to_binlog = 1" 
> exec ^
: string
| SET CONFIG tikv a-no_write_to_binlog = 1
< error
| 1064 42000 …line 1 column 36 near "no_write_to_binlog = 1" 
> exec ^
: string
| KILL nocache()
< error
| 1064 42000 …line 1 column 12 near "nocache()" 
> exec ^
: string
| KILL nocache (1)
< error
| 1064 42000 …line 1 column 12 near "nocache (1)" 
> exec ^
: string
| SET CONFIG tikv nocache = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nocache = 1
< ok 0 0
> exec ^
: string
| KILL nocycle()
< error
| 1064 42000 …line 1 column 12 near "nocycle()" 
> exec ^
: string
| KILL nocycle (1)
< error
| 1064 42000 …line 1 column 12 near "nocycle (1)" 
> exec ^
: string
| SET CONFIG tikv nocycle = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nocycle = 1
< ok 0 0
> exec ^
: string
| KILL node_id()
< error
| 1064 42000 …line 1 column 12 near "node_id()" 
> exec ^
: string
| KILL node_id (1)
< error
| 1064 42000 …line 1 column 12 near "node_id (1)" 
> exec ^
: string
| SET CONFIG tikv node_id = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-node_id = 1
< ok 0 0
> exec ^
: string
| KILL node_state()
< error
| 1064 42000 …line 1 column 15 near "node_state()" 
> exec ^
: string
| KILL node_state (1)
< error
| 1064 42000 …line 1 column 15 near "node_state (1)" 
> exec ^
: string
| SET CONFIG tikv node_state = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-node_state = 1
< ok 0 0
> exec ^
: string
| KILL nodegroup()
< error
| 1064 42000 …line 1 column 14 near "nodegroup()" 
> exec ^
: string
| KILL nodegroup (1)
< error
| 1064 42000 …line 1 column 14 near "nodegroup (1)" 
> exec ^
: string
| SET CONFIG tikv nodegroup = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nodegroup = 1
< ok 0 0
> exec ^
: string
| KILL nomaxvalue()
< error
| 1064 42000 …line 1 column 15 near "nomaxvalue()" 
> exec ^
: string
| KILL nomaxvalue (1)
< error
| 1064 42000 …line 1 column 15 near "nomaxvalue (1)" 
> exec ^
: string
| SET CONFIG tikv nomaxvalue = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nomaxvalue = 1
< ok 0 0
> exec ^
: string
| KILL nominvalue()
< error
| 1064 42000 …line 1 column 15 near "nominvalue()" 
> exec ^
: string
| KILL nominvalue (1)
< error
| 1064 42000 …line 1 column 15 near "nominvalue (1)" 
> exec ^
: string
| SET CONFIG tikv nominvalue = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nominvalue = 1
< ok 0 0
> exec ^
: string
| KILL nonclustered()
< error
| 1064 42000 …line 1 column 17 near "nonclustered()" 
> exec ^
: string
| KILL nonclustered (1)
< error
| 1064 42000 …line 1 column 17 near "nonclustered (1)" 
> exec ^
: string
| SET CONFIG tikv nonclustered = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nonclustered = 1
< ok 0 0
> exec ^
: string
| KILL none()
< error
| 1064 42000 …line 1 column 9 near "none()" 
> exec ^
: string
| KILL none (1)
< error
| 1064 42000 …line 1 column 9 near "none (1)" 
> exec ^
: string
| SET CONFIG tikv none = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-none = 1
< ok 0 0
> exec ^
: string
| KILL not()
< error
| 1064 42000 …line 1 column 8 near "not()" 
> exec ^
: string
| KILL not (1)
< error
| 1064 42000 …line 1 column 8 near "not (1)" 
> exec ^
: string
| SET CONFIG tikv not = 1
< error
| 1064 42000 …line 1 column 19 near "not = 1" 
> exec ^
: string
| SET CONFIG tikv a-not = 1
< error
| 1064 42000 …line 1 column 21 near "not = 1" 
> exec ^
: string
| KILL now()
< error
| 1064 42000 …line 1 column 8 near "now()" 
> exec ^
: string
| KILL now (1)
< error
| 1064 42000 …line 1 column 8 near "now (1)" 
> exec ^
: string
| SET CONFIG tikv now = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-now = 1
< ok 0 0
> exec ^
: string
| KILL nowait()
< error
| 1064 42000 …line 1 column 11 near "nowait()" 
> exec ^
: string
| KILL nowait (1)
< error
| 1064 42000 …line 1 column 11 near "nowait (1)" 
> exec ^
: string
| SET CONFIG tikv nowait = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nowait = 1
< ok 0 0
> exec ^
: string
| KILL nth_value()
< error
| 1064 42000 …line 1 column 14 near "nth_value()" 
> exec ^
: string
| KILL nth_value (1)
< error
| 1064 42000 …line 1 column 14 near "nth_value (1)" 
> exec ^
: string
| SET CONFIG tikv nth_value = 1
< error
| 1064 42000 …line 1 column 25 near "nth_value = 1" 
> exec ^
: string
| SET CONFIG tikv a-nth_value = 1
< error
| 1064 42000 …line 1 column 27 near "nth_value = 1" 
> exec ^
: string
| KILL ntile()
< error
| 1064 42000 …line 1 column 10 near "ntile()" 
> exec ^
: string
| KILL ntile (1)
< error
| 1064 42000 …line 1 column 10 near "ntile (1)" 
> exec ^
: string
| SET CONFIG tikv ntile = 1
< error
| 1064 42000 …line 1 column 21 near "ntile = 1" 
> exec ^
: string
| SET CONFIG tikv a-ntile = 1
< error
| 1064 42000 …line 1 column 23 near "ntile = 1" 
> exec ^
: string
| KILL null()
< error
| 1064 42000 …line 1 column 9 near "null()" 
> exec ^
: string
| KILL null (1)
< error
| 1064 42000 …line 1 column 9 near "null (1)" 
> exec ^
: string
| SET CONFIG tikv null = 1
< error
| 1064 42000 …line 1 column 20 near "null = 1" 
> exec ^
: string
| SET CONFIG tikv a-null = 1
< error
| 1064 42000 …line 1 column 22 near "null = 1" 
> exec ^
: string
| KILL nullif()
< ok 0 0
> exec ^
: string
| KILL nullif (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv nullif = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nullif = 1
< ok 0 0
> exec ^
: string
| KILL nulls()
< error
| 1064 42000 …line 1 column 10 near "nulls()" 
> exec ^
: string
| KILL nulls (1)
< error
| 1064 42000 …line 1 column 10 near "nulls (1)" 
> exec ^
: string
| SET CONFIG tikv nulls = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nulls = 1
< ok 0 0
> exec ^
: string
| KILL number()
< ok 0 0
> exec ^
: string
| KILL number (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv number = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-number = 1
< ok 0 0
> exec ^
: string
| KILL numeric()
< error
| 1064 42000 …line 1 column 12 near "numeric()" 
> exec ^
: string
| KILL numeric (1)
< error
| 1064 42000 …line 1 column 12 near "numeric (1)" 
> exec ^
: string
| SET CONFIG tikv numeric = 1
< error
| 1064 42000 …line 1 column 23 near "numeric = 1" 
> exec ^
: string
| SET CONFIG tikv a-numeric = 1
< error
| 1064 42000 …line 1 column 25 near "numeric = 1" 
> exec ^
: string
| KILL nvarchar()
< error
| 1064 42000 …line 1 column 13 near "nvarchar()" 
> exec ^
: string
| KILL nvarchar (1)
< error
| 1064 42000 …line 1 column 13 near "nvarchar (1)" 
> exec ^
: string
| SET CONFIG tikv nvarchar = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-nvarchar = 1
< ok 0 0
> exec ^
: string
| KILL of()
< error
| 1064 42000 …line 1 column 7 near "of()" 
> exec ^
: string
| KILL of (1)
< error
| 1064 42000 …line 1 column 7 near "of (1)" 
> exec ^
: string
| SET CONFIG tikv of = 1
< error
| 1064 42000 …line 1 column 18 near "of = 1" 
> exec ^
: string
| SET CONFIG tikv a-of = 1
< error
| 1064 42000 …line 1 column 20 near "of = 1" 
> exec ^
: string
| KILL off()
< error
| 1064 42000 …line 1 column 8 near "off()" 
> exec ^
: string
| KILL off (1)
< error
| 1064 42000 …line 1 column 8 near "off (1)" 
> exec ^
: string
| SET CONFIG tikv off = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-off = 1
< ok 0 0
> exec ^
: string
| KILL offset()
< error
| 1064 42000 …line 1 column 11 near "offset()" 
> exec ^
: string
| KILL offset (1)
< error
| 1064 42000 …line 1 column 11 near "offset (1)" 
> exec ^
: string
| SET CONFIG tikv offset = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-offset = 1
< ok 0 0
> exec ^
: string
| KILL oj()
< ok 0 0
> exec ^
: string
| KILL oj (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv oj = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-oj = 1
< ok 0 0
> exec ^
: string
| KILL old()
< ok 0 0
> exec ^
: string
| KILL old (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv old = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-old = 1
< ok 0 0
> exec ^
: string
| KILL oltp_read_only()
< error
| 1064 42000 …line 1 column 19 near "oltp_read_only()" 
> exec ^
: string
| KILL oltp_read_only (1)
< error
| 1064 42000 …line 1 column 19 near "oltp_read_only (1)" 
> exec ^
: string
| SET CONFIG tikv oltp_read_only = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-oltp_read_only = 1
< ok 0 0
> exec ^
: string
| KILL oltp_read_write()
< error
| 1064 42000 …line 1 column 20 near "oltp_read_write()" 
> exec ^
: string
| KILL oltp_read_write (1)
< error
| 1064 42000 …line 1 column 20 near "oltp_read_write (1)" 
> exec ^
: string
| SET CONFIG tikv oltp_read_write = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-oltp_read_write = 1
< ok 0 0
> exec ^
: string
| KILL oltp_write_only()
< error
| 1064 42000 …line 1 column 20 near "oltp_write_only()" 
> exec ^
: string
| KILL oltp_write_only (1)
< error
| 1064 42000 …line 1 column 20 near "oltp_write_only (1)" 
> exec ^
: string
| SET CONFIG tikv oltp_write_only = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-oltp_write_only = 1
< ok 0 0
> exec ^
: string
| KILL on()
< error
| 1064 42000 …line 1 column 7 near "on()" 
> exec ^
: string
| KILL on (1)
< error
| 1064 42000 …line 1 column 7 near "on (1)" 
> exec ^
: string
| SET CONFIG tikv on = 1
< error
| 1064 42000 …line 1 column 18 near "on = 1" 
> exec ^
: string
| SET CONFIG tikv a-on = 1
< error
| 1064 42000 …line 1 column 20 near "on = 1" 
> exec ^
: string
| KILL on_duplicate()
< error
| 1064 42000 …line 1 column 17 near "on_duplicate()" 
> exec ^
: string
| KILL on_duplicate (1)
< error
| 1064 42000 …line 1 column 17 near "on_duplicate (1)" 
> exec ^
: string
| SET CONFIG tikv on_duplicate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-on_duplicate = 1
< ok 0 0
> exec ^
: string
| KILL one()
< ok 0 0
> exec ^
: string
| KILL one (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv one = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-one = 1
< ok 0 0
> exec ^
: string
| KILL online()
< error
| 1064 42000 …line 1 column 11 near "online()" 
> exec ^
: string
| KILL online (1)
< error
| 1064 42000 …line 1 column 11 near "online (1)" 
> exec ^
: string
| SET CONFIG tikv online = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-online = 1
< ok 0 0
> exec ^
: string
| KILL only()
< error
| 1064 42000 …line 1 column 9 near "only()" 
> exec ^
: string
| KILL only (1)
< error
| 1064 42000 …line 1 column 9 near "only (1)" 
> exec ^
: string
| SET CONFIG tikv only = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-only = 1
< ok 0 0
> exec ^
: string
| KILL open()
< error
| 1064 42000 …line 1 column 9 near "open()" 
> exec ^
: string
| KILL open (1)
< error
| 1064 42000 …line 1 column 9 near "open (1)" 
> exec ^
: string
| SET CONFIG tikv open = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-open = 1
< ok 0 0
> exec ^
: string
| KILL opt_rule_blacklist()
< error
| 1064 42000 …line 1 column 23 near "opt_rule_blacklist()" 
> exec ^
: string
| KILL opt_rule_blacklist (1)
< error
| 1064 42000 …line 1 column 23 near "opt_rule_blacklist (1)" 
> exec ^
: string
| SET CONFIG tikv opt_rule_blacklist = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-opt_rule_blacklist = 1
< ok 0 0
> exec ^
: string
| KILL optimistic()
< error
| 1064 42000 …line 1 column 15 near "optimistic()" 
> exec ^
: string
| KILL optimistic (1)
< error
| 1064 42000 …line 1 column 15 near "optimistic (1)" 
> exec ^
: string
| SET CONFIG tikv optimistic = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-optimistic = 1
< ok 0 0
> exec ^
: string
| KILL optimize()
< error
| 1064 42000 …line 1 column 13 near "optimize()" 
> exec ^
: string
| KILL optimize (1)
< error
| 1064 42000 …line 1 column 13 near "optimize (1)" 
> exec ^
: string
| SET CONFIG tikv optimize = 1
< error
| 1064 42000 …line 1 column 24 near "optimize = 1" 
> exec ^
: string
| SET CONFIG tikv a-optimize = 1
< error
| 1064 42000 …line 1 column 26 near "optimize = 1" 
> exec ^
: string
| KILL optimizer_costs()
< ok 0 0
> exec ^
: string
| KILL optimizer_costs (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv optimizer_costs = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-optimizer_costs = 1
< ok 0 0
> exec ^
: string
| KILL option()
< error
| 1064 42000 …line 1 column 11 near "option()" 
> exec ^
: string
| KILL option (1)
< error
| 1064 42000 …line 1 column 11 near "option (1)" 
> exec ^
: string
| SET CONFIG tikv option = 1
< error
| 1064 42000 …line 1 column 22 near "option = 1" 
> exec ^
: string
| SET CONFIG tikv a-option = 1
< error
| 1064 42000 …line 1 column 24 near "option = 1" 
> exec ^
: string
| KILL optional()
< error
| 1064 42000 …line 1 column 13 near "optional()" 
> exec ^
: string
| KILL optional (1)
< error
| 1064 42000 …line 1 column 13 near "optional (1)" 
> exec ^
: string
| SET CONFIG tikv optional = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-optional = 1
< ok 0 0
> exec ^
: string
| KILL optionally()
< error
| 1064 42000 …line 1 column 15 near "optionally()" 
> exec ^
: string
| KILL optionally (1)
< error
| 1064 42000 …line 1 column 15 near "optionally (1)" 
> exec ^
: string
| SET CONFIG tikv optionally = 1
< error
| 1064 42000 …line 1 column 26 near "optionally = 1" 
> exec ^
: string
| SET CONFIG tikv a-optionally = 1
< error
| 1064 42000 …line 1 column 28 near "optionally = 1" 
> exec ^
: string
| KILL options()
< ok 0 0
> exec ^
: string
| KILL options (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv options = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-options = 1
< ok 0 0
> exec ^
: string
| KILL or()
< error
| 1064 42000 …line 1 column 7 near "or()" 
> exec ^
: string
| KILL or (1)
< error
| 1064 42000 …line 1 column 7 near "or (1)" 
> exec ^
: string
| SET CONFIG tikv or = 1
< error
| 1064 42000 …line 1 column 18 near "or = 1" 
> exec ^
: string
| SET CONFIG tikv a-or = 1
< error
| 1064 42000 …line 1 column 20 near "or = 1" 
> exec ^
: string
| KILL order()
< error
| 1064 42000 …line 1 column 10 near "order()" 
> exec ^
: string
| KILL order (1)
< error
| 1064 42000 …line 1 column 10 near "order (1)" 
> exec ^
: string
| SET CONFIG tikv order = 1
< error
| 1064 42000 …line 1 column 21 near "order = 1" 
> exec ^
: string
| SET CONFIG tikv a-order = 1
< error
| 1064 42000 …line 1 column 23 near "order = 1" 
> exec ^
: string
| KILL ordinality()
< ok 0 0
> exec ^
: string
| KILL ordinality (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv ordinality = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ordinality = 1
< ok 0 0
> exec ^
: string
| KILL organization()
< ok 0 0
> exec ^
: string
| KILL organization (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv organization = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-organization = 1
< ok 0 0
> exec ^
: string
| KILL others()
< ok 0 0
> exec ^
: string
| KILL others (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv others = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-others = 1
< ok 0 0
> exec ^
: string
| KILL out()
< error
| 1064 42000 …line 1 column 8 near "out()" 
> exec ^
: string
| KILL out (1)
< error
| 1064 42000 …line 1 column 8 near "out (1)" 
> exec ^
: string
| SET CONFIG tikv out = 1
< error
| 1064 42000 …line 1 column 19 near "out = 1" 
> exec ^
: string
| SET CONFIG tikv a-out = 1
< error
| 1064 42000 …line 1 column 21 near "out = 1" 
> exec ^
: string
| KILL outer()
< error
| 1064 42000 …line 1 column 10 near "outer()" 
> exec ^
: string
| KILL outer (1)
< error
| 1064 42000 …line 1 column 10 near "outer (1)" 
> exec ^
: string
| SET CONFIG tikv outer = 1
< error
| 1064 42000 …line 1 column 21 near "outer = 1" 
> exec ^
: string
| SET CONFIG tikv a-outer = 1
< error
| 1064 42000 …line 1 column 23 near "outer = 1" 
> exec ^
: string
| KILL outfile()
< error
| 1064 42000 …line 1 column 12 near "outfile()" 
> exec ^
: string
| KILL outfile (1)
< error
| 1064 42000 …line 1 column 12 near "outfile (1)" 
> exec ^
: string
| SET CONFIG tikv outfile = 1
< error
| 1064 42000 …line 1 column 23 near "outfile = 1" 
> exec ^
: string
| SET CONFIG tikv a-outfile = 1
< error
| 1064 42000 …line 1 column 25 near "outfile = 1" 
> exec ^
: string
| KILL over()
< error
| 1064 42000 …line 1 column 9 near "over()" 
> exec ^
: string
| KILL over (1)
< error
| 1064 42000 …line 1 column 9 near "over (1)" 
> exec ^
: string
| SET CONFIG tikv over = 1
< error
| 1064 42000 …line 1 column 20 near "over = 1" 
> exec ^
: string
| SET CONFIG tikv a-over = 1
< error
| 1064 42000 …line 1 column 22 near "over = 1" 
> exec ^
: string
| KILL overlay()
< ok 0 0
> exec ^
: string
| KILL overlay (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv overlay = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-overlay = 1
< ok 0 0
> exec ^
: string
| KILL owner()
< ok 0 0
> exec ^
: string
| KILL owner (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv owner = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-owner = 1
< ok 0 0
> exec ^
: string
| KILL pack_keys()
< error
| 1064 42000 …line 1 column 14 near "pack_keys()" 
> exec ^
: string
| KILL pack_keys (1)
< error
| 1064 42000 …line 1 column 14 near "pack_keys (1)" 
> exec ^
: string
| SET CONFIG tikv pack_keys = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-pack_keys = 1
< ok 0 0
> exec ^
: string
| KILL page()
< error
| 1064 42000 …line 1 column 9 near "page()" 
> exec ^
: string
| KILL page (1)
< error
| 1064 42000 …line 1 column 9 near "page (1)" 
> exec ^
: string
| SET CONFIG tikv page = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-page = 1
< ok 0 0
> exec ^
: string
| KILL parallel()
< ok 0 0
> exec ^
: string
| KILL parallel (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv parallel = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-parallel = 1
< ok 0 0
> exec ^
: string
| KILL parameters()
< ok 0 0
> exec ^
: string
| KILL parameters (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv parameters = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-parameters = 1
< ok 0 0
> exec ^
: string
| KILL parse_tree()
< ok 0 0
> exec ^
: string
| KILL parse_tree (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv parse_tree = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-parse_tree = 1
< ok 0 0
> exec ^
: string
| KILL parser()
< error
| 1064 42000 …line 1 column 11 near "parser()" 
> exec ^
: string
| KILL parser (1)
< error
| 1064 42000 …line 1 column 11 near "parser (1)" 
> exec ^
: string
| SET CONFIG tikv parser = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-parser = 1
< ok 0 0
> exec ^
: string
| KILL partial()
< error
| 1064 42000 …line 1 column 12 near "partial()" 
> exec ^
: string
| KILL partial (1)
< error
| 1064 42000 …line 1 column 12 near "partial (1)" 
> exec ^
: string
| SET CONFIG tikv partial = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-partial = 1
< ok 0 0
> exec ^
: string
| KILL partition()
< error
| 1064 42000 …line 1 column 14 near "partition()" 
> exec ^
: string
| KILL partition (1)
< error
| 1064 42000 …line 1 column 14 near "partition (1)" 
> exec ^
: string
| SET CONFIG tikv partition = 1
< error
| 1064 42000 …line 1 column 25 near "partition = 1" 
> exec ^
: string
| SET CONFIG tikv a-partition = 1
< error
| 1064 42000 …line 1 column 27 near "partition = 1" 
> exec ^
: string
| KILL partitioning()
< error
| 1064 42000 …line 1 column 17 near "partitioning()" 
> exec ^
: string
| KILL partitioning (1)
< error
| 1064 42000 …line 1 column 17 near "partitioning (1)" 
> exec ^
: string
| SET CONFIG tikv partitioning = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-partitioning = 1
< ok 0 0
> exec ^
: string
| KILL partitions()
< error
| 1064 42000 …line 1 column 15 near "partitions()" 
> exec ^
: string
| KILL partitions (1)
< error
| 1064 42000 …line 1 column 15 near "partitions (1)" 
> exec ^
: string
| SET CONFIG tikv partitions = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-partitions = 1
< ok 0 0
> exec ^
: string
| KILL password()
< error
| 1064 42000 …line 1 column 13 near "password()" 
> exec ^
: string
| KILL password (1)
< error
| 1064 42000 …line 1 column 13 near "password (1)" 
> exec ^
: string
| SET CONFIG tikv password = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-password = 1
< ok 0 0
> exec ^
: string
| KILL password_lock_time()
< error
| 1064 42000 …line 1 column 23 near "password_lock_time()" 
> exec ^
: string
| KILL password_lock_time (1)
< error
| 1064 42000 …line 1 column 23 near "password_lock_time (1)" 
> exec ^
: string
| SET CONFIG tikv password_lock_time = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-password_lock_time = 1
< ok 0 0
> exec ^
: string
| KILL path()
< ok 0 0
> exec ^
: string
| KILL path (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv path = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-path = 1
< ok 0 0
> exec ^
: string
| KILL pause()
< error
| 1064 42000 …line 1 column 10 near "pause()" 
> exec ^
: string
| KILL pause (1)
< error
| 1064 42000 …line 1 column 10 near "pause (1)" 
> exec ^
: string
| SET CONFIG tikv pause = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-pause = 1
< ok 0 0
> exec ^
: string
| KILL per_db()
< error
| 1064 42000 …line 1 column 11 near "per_db()" 
> exec ^
: string
| KILL per_db (1)
< error
| 1064 42000 …line 1 column 11 near "per_db (1)" 
> exec ^
: string
| SET CONFIG tikv per_db = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-per_db = 1
< ok 0 0
> exec ^
: string
| KILL per_table()
< error
| 1064 42000 …line 1 column 14 near "per_table()" 
> exec ^
: string
| KILL per_table (1)
< error
| 1064 42000 …line 1 column 14 near "per_table (1)" 
> exec ^
: string
| SET CONFIG tikv per_table = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-per_table = 1
< ok 0 0
> exec ^
: string
| KILL percent()
< error
| 1064 42000 …line 1 column 12 near "percent()" 
> exec ^
: string
| KILL percent (1)
< error
| 1064 42000 …line 1 column 12 near "percent (1)" 
> exec ^
: string
| SET CONFIG tikv percent = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-percent = 1
< ok 0 0
> exec ^
: string
| KILL percent_rank()
< error
| 1064 42000 …line 1 column 17 near "percent_rank()" 
> exec ^
: string
| KILL percent_rank (1)
< error
| 1064 42000 …line 1 column 17 near "percent_rank (1)" 
> exec ^
: string
| SET CONFIG tikv percent_rank = 1
< error
| 1064 42000 …line 1 column 28 near "percent_rank = 1" 
> exec ^
: string
| SET CONFIG tikv a-percent_rank = 1
< error
| 1064 42000 …line 1 column 30 near "percent_rank = 1" 
> exec ^
: string
| KILL persist()
< ok 0 0
> exec ^
: string
| KILL persist (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv persist = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-persist = 1
< ok 0 0
> exec ^
: string
| KILL persist_only()
< ok 0 0
> exec ^
: string
| KILL persist_only (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv persist_only = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-persist_only = 1
< ok 0 0
> exec ^
: string
| KILL pessimistic()
< error
| 1064 42000 …line 1 column 16 near "pessimistic()" 
> exec ^
: string
| KILL pessimistic (1)
< error
| 1064 42000 …line 1 column 16 near "pessimistic (1)" 
> exec ^
: string
| SET CONFIG tikv pessimistic = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-pessimistic = 1
< ok 0 0
> exec ^
: string
| KILL phase()
< ok 0 0
> exec ^
: string
| KILL phase (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv phase = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-phase = 1
< ok 0 0
> exec ^
: string
| KILL placement()
< error
| 1064 42000 …line 1 column 14 near "placement()" 
> exec ^
: string
| KILL placement (1)
< error
| 1064 42000 …line 1 column 14 near "placement (1)" 
> exec ^
: string
| SET CONFIG tikv placement = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-placement = 1
< ok 0 0
> exec ^
: string
| KILL placing()
< ok 0 0
> exec ^
: string
| KILL placing (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv placing = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-placing = 1
< ok 0 0
> exec ^
: string
| KILL plan()
< error
| 1064 42000 …line 1 column 9 near "plan()" 
> exec ^
: string
| KILL plan (1)
< error
| 1064 42000 …line 1 column 9 near "plan (1)" 
> exec ^
: string
| SET CONFIG tikv plan = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-plan = 1
< ok 0 0
> exec ^
: string
| KILL plan_cache()
< error
| 1064 42000 …line 1 column 15 near "plan_cache()" 
> exec ^
: string
| KILL plan_cache (1)
< error
| 1064 42000 …line 1 column 15 near "plan_cache (1)" 
> exec ^
: string
| SET CONFIG tikv plan_cache = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-plan_cache = 1
< ok 0 0
> exec ^
: string
| KILL plugin()
< ok 0 0
> exec ^
: string
| KILL plugin (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv plugin = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-plugin = 1
< ok 0 0
> exec ^
: string
| KILL plugin_dir()
< ok 0 0
> exec ^
: string
| KILL plugin_dir (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv plugin_dir = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-plugin_dir = 1
< ok 0 0
> exec ^
: string
| KILL plugins()
< error
| 1064 42000 …line 1 column 12 near "plugins()" 
> exec ^
: string
| KILL plugins (1)
< error
| 1064 42000 …line 1 column 12 near "plugins (1)" 
> exec ^
: string
| SET CONFIG tikv plugins = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-plugins = 1
< ok 0 0
> exec ^
: string
| KILL point()
< error
| 1064 42000 …line 1 column 10 near "point()" 
> exec ^
: string
| KILL point (1)
< error
| 1064 42000 …line 1 column 10 near "point (1)" 
> exec ^
: string
| SET CONFIG tikv point = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-point = 1
< ok 0 0
> exec ^
: string
| KILL policy()
< error
| 1064 42000 …line 1 column 11 near "policy()" 
> exec ^
: string
| KILL policy (1)
< error
| 1064 42000 …line 1 column 11 near "policy (1)" 
> exec ^
: string
| SET CONFIG tikv policy = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-policy = 1
< ok 0 0
> exec ^
: string
| KILL polygon()
< ok 0 0
> exec ^
: string
| KILL polygon (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv polygon = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-polygon = 1
< ok 0 0
> exec ^
: string
| KILL port()
< ok 0 0
> exec ^
: string
| KILL port (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv port = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-port = 1
< ok 0 0
> exec ^
: string
| KILL position()
< error
| 1064 42000 …line 1 column 13 near "position()" 
> exec ^
: string
| KILL position (1)
< error
| 1064 42000 …line 1 column 13 near "position (1)" 
> exec ^
: string
| SET CONFIG tikv position = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-position = 1
< ok 0 0
> exec ^
: string
| KILL pre_split_regions()
< error
| 1064 42000 …line 1 column 22 near "pre_split_regions()" 
> exec ^
: string
| KILL pre_split_regions (1)
< error
| 1064 42000 …line 1 column 22 near "pre_split_regions (1)" 
> exec ^
: string
| SET CONFIG tikv pre_split_regions = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-pre_split_regions = 1
< ok 0 0
> exec ^
: string
| KILL precedes()
< ok 0 0
> exec ^
: string
| KILL precedes (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv precedes = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-precedes = 1
< ok 0 0
> exec ^
: string
| KILL preceding()
< error
| 1064 42000 …line 1 column 14 near "preceding()" 
> exec ^
: string
| KILL preceding (1)
< error
| 1064 42000 …line 1 column 14 near "preceding (1)" 
> exec ^
: string
| SET CONFIG tikv preceding = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-preceding = 1
< ok 0 0
> exec ^
: string
| KILL precision()
< error
| 1064 42000 …line 1 column 14 near "precision()" 
> exec ^
: string
| KILL precision (1)
< error
| 1064 42000 …line 1 column 14 near "precision (1)" 
> exec ^
: string
| SET CONFIG tikv precision = 1
< error
| 1064 42000 …line 1 column 25 near "precision = 1" 
> exec ^
: string
| SET CONFIG tikv a-precision = 1
< error
| 1064 42000 …line 1 column 27 near "precision = 1" 
> exec ^
: string
| KILL predicate()
< error
| 1064 42000 …line 1 column 14 near "predicate()" 
> exec ^
: string
| KILL predicate (1)
< error
| 1064 42000 …line 1 column 14 near "predicate (1)" 
> exec ^
: string
| SET CONFIG tikv predicate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-predicate = 1
< ok 0 0
> exec ^
: string
| KILL prepare()
< error
| 1064 42000 …line 1 column 12 near "prepare()" 
> exec ^
: string
| KILL prepare (1)
< error
| 1064 42000 …line 1 column 12 near "prepare (1)" 
> exec ^
: string
| SET CONFIG tikv prepare = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-prepare = 1
< ok 0 0
> exec ^
: string
| KILL preserve()
< error
| 1064 42000 …line 1 column 13 near "preserve()" 
> exec ^
: string
| KILL preserve (1)
< error
| 1064 42000 …line 1 column 13 near "preserve (1)" 
> exec ^
: string
| SET CONFIG tikv preserve = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-preserve = 1
< ok 0 0
> exec ^
: string
| KILL prev()
< ok 0 0
> exec ^
: string
| KILL prev (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv prev = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-prev = 1
< ok 0 0
> exec ^
: string
| KILL primary()
< error
| 1064 42000 …line 1 column 12 near "primary()" 
> exec ^
: string
| KILL primary (1)
< error
| 1064 42000 …line 1 column 12 near "primary (1)" 
> exec ^
: string
| SET CONFIG tikv primary = 1
< error
| 1064 42000 …line 1 column 23 near "primary = 1" 
> exec ^
: string
| SET CONFIG tikv a-primary = 1
< error
| 1064 42000 …line 1 column 25 near "primary = 1" 
> exec ^
: string
| KILL primary_region()
< error
| 1064 42000 …line 1 column 19 near "primary_region()" 
> exec ^
: string
| KILL primary_region (1)
< error
| 1064 42000 …line 1 column 19 near "primary_region (1)" 
> exec ^
: string
| SET CONFIG tikv primary_region = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-primary_region = 1
< ok 0 0
> exec ^
: string
| KILL priority()
< error
| 1064 42000 …line 1 column 13 near "priority()" 
> exec ^
: string
| KILL priority (1)
< error
| 1064 42000 …line 1 column 13 near "priority (1)" 
> exec ^
: string
| SET CONFIG tikv priority = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-priority = 1
< ok 0 0
> exec ^
: string
| KILL privilege_checks_user()
< ok 0 0
> exec ^
: string
| KILL privilege_checks_user (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv privilege_checks_user = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-privilege_checks_user = 1
< ok 0 0
> exec ^
: string
| KILL privileges()
< error
| 1064 42000 …line 1 column 15 near "privileges()" 
> exec ^
: string
| KILL privileges (1)
< error
| 1064 42000 …line 1 column 15 near "privileges (1)" 
> exec ^
: string
| SET CONFIG tikv privileges = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-privileges = 1
< ok 0 0
> exec ^
: string
| KILL procedure()
< error
| 1064 42000 …line 1 column 14 near "procedure()" 
> exec ^
: string
| KILL procedure (1)
< error
| 1064 42000 …line 1 column 14 near "procedure (1)" 
> exec ^
: string
| SET CONFIG tikv procedure = 1
< error
| 1064 42000 …line 1 column 25 near "procedure = 1" 
> exec ^
: string
| SET CONFIG tikv a-procedure = 1
< error
| 1064 42000 …line 1 column 27 near "procedure = 1" 
> exec ^
: string
| KILL process()
< error
| 1064 42000 …line 1 column 12 near "process()" 
> exec ^
: string
| KILL process (1)
< error
| 1064 42000 …line 1 column 12 near "process (1)" 
> exec ^
: string
| SET CONFIG tikv process = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-process = 1
< ok 0 0
> exec ^
: string
| KILL processed_keys()
< error
| 1064 42000 …line 1 column 19 near "processed_keys()" 
> exec ^
: string
| KILL processed_keys (1)
< error
| 1064 42000 …line 1 column 19 near "processed_keys (1)" 
> exec ^
: string
| SET CONFIG tikv processed_keys = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-processed_keys = 1
< ok 0 0
> exec ^
: string
| KILL processlist()
< error
| 1064 42000 …line 1 column 16 near "processlist()" 
> exec ^
: string
| KILL processlist (1)
< error
| 1064 42000 …line 1 column 16 near "processlist (1)" 
> exec ^
: string
| SET CONFIG tikv processlist = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-processlist = 1
< ok 0 0
> exec ^
: string
| KILL profile()
< error
| 1064 42000 …line 1 column 12 near "profile()" 
> exec ^
: string
| KILL profile (1)
< error
| 1064 42000 …line 1 column 12 near "profile (1)" 
> exec ^
: string
| SET CONFIG tikv profile = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-profile = 1
< ok 0 0
> exec ^
: string
| KILL profiles()
< error
| 1064 42000 …line 1 column 13 near "profiles()" 
> exec ^
: string
| KILL profiles (1)
< error
| 1064 42000 …line 1 column 13 near "profiles (1)" 
> exec ^
: string
| SET CONFIG tikv profiles = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-profiles = 1
< ok 0 0
> exec ^
: string
| KILL proxy()
< error
| 1064 42000 …line 1 column 10 near "proxy()" 
> exec ^
: string
| KILL proxy (1)
< error
| 1064 42000 …line 1 column 10 near "proxy (1)" 
> exec ^
: string
| SET CONFIG tikv proxy = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-proxy = 1
< ok 0 0
> exec ^
: string
| KILL purge()
< error
| 1064 42000 …line 1 column 10 near "purge()" 
> exec ^
: string
| KILL purge (1)
< error
| 1064 42000 …line 1 column 10 near "purge (1)" 
> exec ^
: string
| SET CONFIG tikv purge = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-purge = 1
< ok 0 0
> exec ^
: string
| KILL qualify()
< ok 0 0
> exec ^
: string
| KILL qualify (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv qualify = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-qualify = 1
< ok 0 0
> exec ^
: string
| KILL quarter()
< error
| 1064 42000 …line 1 column 12 near "quarter()" 
> exec ^
: string
| KILL quarter (1)
< error
| 1064 42000 …line 1 column 12 near "quarter (1)" 
> exec ^
: string
| SET CONFIG tikv quarter = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-quarter = 1
< ok 0 0
> exec ^
: string
| KILL queries()
< error
| 1064 42000 …line 1 column 12 near "queries()" 
> exec ^
: string
| KILL queries (1)
< error
| 1064 42000 …line 1 column 12 near "queries (1)" 
> exec ^
: string
| SET CONFIG tikv queries = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-queries = 1
< ok 0 0
> exec ^
: string
| KILL query()
< error
| 1064 42000 …line 1 column 11 near "()" 
> exec ^
: string
| KILL query (1)
< error
| 1064 42000 …line 1 column 12 near "(1)" 
> exec ^
: string
| SET CONFIG tikv query = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-query = 1
< ok 0 0
> exec ^
: string
| KILL query_limit()
< error
| 1064 42000 …line 1 column 16 near "query_limit()" 
> exec ^
: string
| KILL query_limit (1)
< error
| 1064 42000 …line 1 column 16 near "query_limit (1)" 
> exec ^
: string
| SET CONFIG tikv query_limit = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-query_limit = 1
< ok 0 0
> exec ^
: string
| KILL quick()
< error
| 1064 42000 …line 1 column 10 near "quick()" 
> exec ^
: string
| KILL quick (1)
< error
| 1064 42000 …line 1 column 10 near "quick (1)" 
> exec ^
: string
| SET CONFIG tikv quick = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-quick = 1
< ok 0 0
> exec ^
: string
| KILL random()
< ok 0 0
> exec ^
: string
| KILL random (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv random = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-random = 1
< ok 0 0
> exec ^
: string
| KILL range()
< error
| 1064 42000 …line 1 column 10 near "range()" 
> exec ^
: string
| KILL range (1)
< error
| 1064 42000 …line 1 column 10 near "range (1)" 
> exec ^
: string
| SET CONFIG tikv range = 1
< error
| 1064 42000 …line 1 column 21 near "range = 1" 
> exec ^
: string
| SET CONFIG tikv a-range = 1
< error
| 1064 42000 …line 1 column 23 near "range = 1" 
> exec ^
: string
| KILL rank()
< error
| 1064 42000 …line 1 column 9 near "rank()" 
> exec ^
: string
| KILL rank (1)
< error
| 1064 42000 …line 1 column 9 near "rank (1)" 
> exec ^
: string
| SET CONFIG tikv rank = 1
< error
| 1064 42000 …line 1 column 20 near "rank = 1" 
> exec ^
: string
| SET CONFIG tikv a-rank = 1
< error
| 1064 42000 …line 1 column 22 near "rank = 1" 
> exec ^
: string
| KILL rate_limit()
< error
| 1064 42000 …line 1 column 15 near "rate_limit()" 
> exec ^
: string
| KILL rate_limit (1)
< error
| 1064 42000 …line 1 column 15 near "rate_limit (1)" 
> exec ^
: string
| SET CONFIG tikv rate_limit = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-rate_limit = 1
< ok 0 0
> exec ^
: string
| KILL read()
< error
| 1064 42000 …line 1 column 9 near "read()" 
> exec ^
: string
| KILL read (1)
< error
| 1064 42000 …line 1 column 9 near "read (1)" 
> exec ^
: string
| SET CONFIG tikv read = 1
< error
| 1064 42000 …line 1 column 20 near "read = 1" 
> exec ^
: string
| SET CONFIG tikv a-read = 1
< error
| 1064 42000 …line 1 column 22 near "read = 1" 
> exec ^
: string
| KILL read_only()
< ok 0 0
> exec ^
: string
| KILL read_only (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv read_only = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-read_only = 1
< ok 0 0
> exec ^
: string
| KILL read_write()
< ok 0 0
> exec ^
: string
| KILL read_write (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv read_write = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-read_write = 1
< ok 0 0
> exec ^
: string
| KILL reads()
< ok 0 0
> exec ^
: string
| KILL reads (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv reads = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-reads = 1
< ok 0 0
> exec ^
: string
| KILL real()
< error
| 1064 42000 …line 1 column 9 near "real()" 
> exec ^
: string
| KILL real (1)
< error
| 1064 42000 …line 1 column 9 near "real (1)" 
> exec ^
: string
| SET CONFIG tikv real = 1
< error
| 1064 42000 …line 1 column 20 near "real = 1" 
> exec ^
: string
| SET CONFIG tikv a-real = 1
< error
| 1064 42000 …line 1 column 22 near "real = 1" 
> exec ^
: string
| KILL rebuild()
< error
| 1064 42000 …line 1 column 12 near "rebuild()" 
> exec ^
: string
| KILL rebuild (1)
< error
| 1064 42000 …line 1 column 12 near "rebuild (1)" 
> exec ^
: string
| SET CONFIG tikv rebuild = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-rebuild = 1
< ok 0 0
> exec ^
: string
| KILL recent()
< error
| 1064 42000 …line 1 column 11 near "recent()" 
> exec ^
: string
| KILL recent (1)
< error
| 1064 42000 …line 1 column 11 near "recent (1)" 
> exec ^
: string
| SET CONFIG tikv recent = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-recent = 1
< ok 0 0
> exec ^
: string
| KILL recommend()
< error
| 1064 42000 …line 1 column 14 near "recommend()" 
> exec ^
: string
| KILL recommend (1)
< error
| 1064 42000 …line 1 column 14 near "recommend (1)" 
> exec ^
: string
| SET CONFIG tikv recommend = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-recommend = 1
< ok 0 0
> exec ^
: string
| KILL recover()
< error
| 1064 42000 …line 1 column 12 near "recover()" 
> exec ^
: string
| KILL recover (1)
< error
| 1064 42000 …line 1 column 12 near "recover (1)" 
> exec ^
: string
| SET CONFIG tikv recover = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-recover = 1
< ok 0 0
> exec ^
: string
| KILL recursive()
< error
| 1064 42000 …line 1 column 14 near "recursive()" 
> exec ^
: string
| KILL recursive (1)
< error
| 1064 42000 …line 1 column 14 near "recursive (1)" 
> exec ^
: string
| SET CONFIG tikv recursive = 1
< error
| 1064 42000 …line 1 column 25 near "recursive = 1" 
> exec ^
: string
| SET CONFIG tikv a-recursive = 1
< error
| 1064 42000 …line 1 column 27 near "recursive = 1" 
> exec ^
: string
| KILL redo_buffer_size()
< ok 0 0
> exec ^
: string
| KILL redo_buffer_size (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv redo_buffer_size = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-redo_buffer_size = 1
< ok 0 0
> exec ^
: string
| KILL redundant()
< error
| 1064 42000 …line 1 column 14 near "redundant()" 
> exec ^
: string
| KILL redundant (1)
< error
| 1064 42000 …line 1 column 14 near "redundant (1)" 
> exec ^
: string
| SET CONFIG tikv redundant = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-redundant = 1
< ok 0 0
> exec ^
: string
| KILL reference()
< ok 0 0
> exec ^
: string
| KILL reference (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv reference = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-reference = 1
< ok 0 0
> exec ^
: string
| KILL references()
< error
| 1064 42000 …line 1 column 15 near "references()" 
> exec ^
: string
| KILL references (1)
< error
| 1064 42000 …line 1 column 15 near "references (1)" 
> exec ^
: string
| SET CONFIG tikv references = 1
< error
| 1064 42000 …line 1 column 26 near "references = 1" 
> exec ^
: string
| SET CONFIG tikv a-references = 1
< error
| 1064 42000 …line 1 column 28 near "references = 1" 
> exec ^
: string
| KILL regexp()
< error
| 1064 42000 …line 1 column 11 near "regexp()" 
> exec ^
: string
| KILL regexp (1)
< error
| 1064 42000 …line 1 column 11 near "regexp (1)" 
> exec ^
: string
| SET CONFIG tikv regexp = 1
< error
| 1064 42000 …line 1 column 22 near "regexp = 1" 
> exec ^
: string
| SET CONFIG tikv a-regexp = 1
< error
| 1064 42000 …line 1 column 24 near "regexp = 1" 
> exec ^
: string
| KILL region()
< error
| 1064 42000 …line 1 column 11 near "region()" 
> exec ^
: string
| KILL region (1)
< error
| 1064 42000 …line 1 column 11 near "region (1)" 
> exec ^
: string
| SET CONFIG tikv region = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-region = 1
< ok 0 0
> exec ^
: string
| KILL regions()
< error
| 1064 42000 …line 1 column 12 near "regions()" 
> exec ^
: string
| KILL regions (1)
< error
| 1064 42000 …line 1 column 12 near "regions (1)" 
> exec ^
: string
| SET CONFIG tikv regions = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-regions = 1
< ok 0 0
> exec ^
: string
| KILL registration()
< ok 0 0
> exec ^
: string
| KILL registration (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv registration = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-registration = 1
< ok 0 0
> exec ^
: string
| KILL relational()
< ok 0 0
> exec ^
: string
| KILL relational (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv relational = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-relational = 1
< ok 0 0
> exec ^
: string
| KILL relay()
< ok 0 0
> exec ^
: string
| KILL relay (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv relay = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-relay = 1
< ok 0 0
> exec ^
: string
| KILL relay_log_file()
< ok 0 0
> exec ^
: string
| KILL relay_log_file (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv relay_log_file = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-relay_log_file = 1
< ok 0 0
> exec ^
: string
| KILL relay_log_pos()
< ok 0 0
> exec ^
: string
| KILL relay_log_pos (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv relay_log_pos = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-relay_log_pos = 1
< ok 0 0
> exec ^
: string
| KILL relay_thread()
< ok 0 0
> exec ^
: string
| KILL relay_thread (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv relay_thread = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-relay_thread = 1
< ok 0 0
> exec ^
: string
| KILL relaylog()
< ok 0 0
> exec ^
: string
| KILL relaylog (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv relaylog = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-relaylog = 1
< ok 0 0
> exec ^
: string
| KILL release()
< error
| 1064 42000 …line 1 column 12 near "release()" 
> exec ^
: string
| KILL release (1)
< error
| 1064 42000 …line 1 column 12 near "release (1)" 
> exec ^
: string
| SET CONFIG tikv release = 1
< error
| 1064 42000 …line 1 column 23 near "release = 1" 
> exec ^
: string
| SET CONFIG tikv a-release = 1
< error
| 1064 42000 …line 1 column 25 near "release = 1" 
> exec ^
: string
| KILL reload()
< error
| 1064 42000 …line 1 column 11 near "reload()" 
> exec ^
: string
| KILL reload (1)
< error
| 1064 42000 …line 1 column 11 near "reload (1)" 
> exec ^
: string
| SET CONFIG tikv reload = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-reload = 1
< ok 0 0
> exec ^
: string
| KILL remove()
< error
| 1064 42000 …line 1 column 11 near "remove()" 
> exec ^
: string
| KILL remove (1)
< error
| 1064 42000 …line 1 column 11 near "remove (1)" 
> exec ^
: string
| SET CONFIG tikv remove = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-remove = 1
< ok 0 0
> exec ^
: string
| KILL rename()
< error
| 1064 42000 …line 1 column 11 near "rename()" 
> exec ^
: string
| KILL rename (1)
< error
| 1064 42000 …line 1 column 11 near "rename (1)" 
> exec ^
: string
| SET CONFIG tikv rename = 1
< error
| 1064 42000 …line 1 column 22 near "rename = 1" 
> exec ^
: string
| SET CONFIG tikv a-rename = 1
< error
| 1064 42000 …line 1 column 24 near "rename = 1" 
> exec ^
: string
| KILL reorganize()
< error
| 1064 42000 …line 1 column 15 near "reorganize()" 
> exec ^
: string
| KILL reorganize (1)
< error
| 1064 42000 …line 1 column 15 near "reorganize (1)" 
> exec ^
: string
| SET CONFIG tikv reorganize = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-reorganize = 1
< ok 0 0
> exec ^
: string
| KILL repair()
< error
| 1064 42000 …line 1 column 11 near "repair()" 
> exec ^
: string
| KILL repair (1)
< error
| 1064 42000 …line 1 column 11 near "repair (1)" 
> exec ^
: string
| SET CONFIG tikv repair = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-repair = 1
< ok 0 0
> exec ^
: string
| KILL repeat()
< error
| 1064 42000 …line 1 column 11 near "repeat()" 
> exec ^
: string
| KILL repeat (1)
< error
| 1064 42000 …line 1 column 11 near "repeat (1)" 
> exec ^
: string
| SET CONFIG tikv repeat = 1
< error
| 1064 42000 …line 1 column 22 near "repeat = 1" 
> exec ^
: string
| SET CONFIG tikv a-repeat = 1
< error
| 1064 42000 …line 1 column 24 near "repeat = 1" 
> exec ^
: string
| KILL repeatable()
< error
| 1064 42000 …line 1 column 15 near "repeatable()" 
> exec ^
: string
| KILL repeatable (1)
< error
| 1064 42000 …line 1 column 15 near "repeatable (1)" 
> exec ^
: string
| SET CONFIG tikv repeatable = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-repeatable = 1
< ok 0 0
> exec ^
: string
| KILL replace()
< error
| 1064 42000 …line 1 column 14 near ")" 
> exec ^
: string
| KILL replace (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv replace = 1
< error
| 1064 42000 …line 1 column 23 near "replace = 1" 
> exec ^
: string
| SET CONFIG tikv a-replace = 1
< error
| 1064 42000 …line 1 column 25 near "replace = 1" 
> exec ^
: string
| KILL replayer()
< error
| 1064 42000 …line 1 column 13 near "replayer()" 
> exec ^
: string
| KILL replayer (1)
< error
| 1064 42000 …line 1 column 13 near "replayer (1)" 
> exec ^
: string
| SET CONFIG tikv replayer = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replayer = 1
< ok 0 0
> exec ^
: string
| KILL replica()
< error
| 1064 42000 …line 1 column 12 near "replica()" 
> exec ^
: string
| KILL replica (1)
< error
| 1064 42000 …line 1 column 12 near "replica (1)" 
> exec ^
: string
| SET CONFIG tikv replica = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replica = 1
< ok 0 0
> exec ^
: string
| KILL replicas()
< error
| 1064 42000 …line 1 column 13 near "replicas()" 
> exec ^
: string
| KILL replicas (1)
< error
| 1064 42000 …line 1 column 13 near "replicas (1)" 
> exec ^
: string
| SET CONFIG tikv replicas = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replicas = 1
< ok 0 0
> exec ^
: string
| KILL replicate_do_db()
< ok 0 0
> exec ^
: string
| KILL replicate_do_db (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv replicate_do_db = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replicate_do_db = 1
< ok 0 0
> exec ^
: string
| KILL replicate_do_table()
< ok 0 0
> exec ^
: string
| KILL replicate_do_table (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv replicate_do_table = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replicate_do_table = 1
< ok 0 0
> exec ^
: string
| KILL replicate_ignore_db()
< ok 0 0
> exec ^
: string
| KILL replicate_ignore_db (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv replicate_ignore_db = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replicate_ignore_db = 1
< ok 0 0
> exec ^
: string
| KILL replicate_ignore_table()
< ok 0 0
> exec ^
: string
| KILL replicate_ignore_table (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv replicate_ignore_table = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replicate_ignore_table = 1
< ok 0 0
> exec ^
: string
| KILL replicate_rewrite_db()
< ok 0 0
> exec ^
: string
| KILL replicate_rewrite_db (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv replicate_rewrite_db = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replicate_rewrite_db = 1
< ok 0 0
> exec ^
: string
| KILL replicate_wild_do_table()
< ok 0 0
> exec ^
: string
| KILL replicate_wild_do_table (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv replicate_wild_do_table = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replicate_wild_do_table = 1
< ok 0 0
> exec ^
: string
| KILL replicate_wild_ignore_table()
< ok 0 0
> exec ^
: string
| KILL replicate_wild_ignore_table (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv replicate_wild_ignore_table = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replicate_wild_ignore_table = 1
< ok 0 0
> exec ^
: string
| KILL replication()
< error
| 1064 42000 …line 1 column 16 near "replication()" 
> exec ^
: string
| KILL replication (1)
< error
| 1064 42000 …line 1 column 16 near "replication (1)" 
> exec ^
: string
| SET CONFIG tikv replication = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-replication = 1
< ok 0 0
> exec ^
: string
| KILL require()
< error
| 1064 42000 …line 1 column 12 near "require()" 
> exec ^
: string
| KILL require (1)
< error
| 1064 42000 …line 1 column 12 near "require (1)" 
> exec ^
: string
| SET CONFIG tikv require = 1
< error
| 1064 42000 …line 1 column 23 near "require = 1" 
> exec ^
: string
| SET CONFIG tikv a-require = 1
< error
| 1064 42000 …line 1 column 25 near "require = 1" 
> exec ^
: string
| KILL require_row_format()
< ok 0 0
> exec ^
: string
| KILL require_row_format (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv require_row_format = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-require_row_format = 1
< ok 0 0
> exec ^
: string
| KILL require_table_primary_key_check()
< ok 0 0
> exec ^
: string
| KILL require_table_primary_key_check (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv require_table_primary_key_check = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-require_table_primary_key_check = 1
< ok 0 0
> exec ^
: string
| KILL required()
< error
| 1064 42000 …line 1 column 13 near "required()" 
> exec ^
: string
| KILL required (1)
< error
| 1064 42000 …line 1 column 13 near "required (1)" 
> exec ^
: string
| SET CONFIG tikv required = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-required = 1
< ok 0 0
> exec ^
: string
| KILL reset()
< error
| 1064 42000 …line 1 column 10 near "reset()" 
> exec ^
: string
| KILL reset (1)
< error
| 1064 42000 …line 1 column 10 near "reset (1)" 
> exec ^
: string
| SET CONFIG tikv reset = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-reset = 1
< ok 0 0
> exec ^
: string
| KILL resignal()
< ok 0 0
> exec ^
: string
| KILL resignal (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv resignal = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-resignal = 1
< ok 0 0
> exec ^
: string
| KILL resource()
< error
| 1064 42000 …line 1 column 13 near "resource()" 
> exec ^
: string
| KILL resource (1)
< error
| 1064 42000 …line 1 column 13 near "resource (1)" 
> exec ^
: string
| SET CONFIG tikv resource = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-resource = 1
< ok 0 0
> exec ^
: string
| KILL respect()
< error
| 1064 42000 …line 1 column 12 near "respect()" 
> exec ^
: string
| KILL respect (1)
< error
| 1064 42000 …line 1 column 12 near "respect (1)" 
> exec ^
: string
| SET CONFIG tikv respect = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-respect = 1
< ok 0 0
> exec ^
: string
| KILL restart()
< error
| 1064 42000 …line 1 column 12 near "restart()" 
> exec ^
: string
| KILL restart (1)
< error
| 1064 42000 …line 1 column 12 near "restart (1)" 
> exec ^
: string
| SET CONFIG tikv restart = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-restart = 1
< ok 0 0
> exec ^
: string
| KILL restore()
< error
| 1064 42000 …line 1 column 12 near "restore()" 
> exec ^
: string
| KILL restore (1)
< error
| 1064 42000 …line 1 column 12 near "restore (1)" 
> exec ^
: string
| SET CONFIG tikv restore = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-restore = 1
< ok 0 0
> exec ^
: string
| KILL restored_ts()
< error
| 1064 42000 …line 1 column 16 near "restored_ts()" 
> exec ^
: string
| KILL restored_ts (1)
< error
| 1064 42000 …line 1 column 16 near "restored_ts (1)" 
> exec ^
: string
| SET CONFIG tikv restored_ts = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-restored_ts = 1
< ok 0 0
> exec ^
: string
| KILL restores()
< error
| 1064 42000 …line 1 column 13 near "restores()" 
> exec ^
: string
| KILL restores (1)
< error
| 1064 42000 …line 1 column 13 near "restores (1)" 
> exec ^
: string
| SET CONFIG tikv restores = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-restores = 1
< ok 0 0
> exec ^
: string
| KILL restrict()
< error
| 1064 42000 …line 1 column 13 near "restrict()" 
> exec ^
: string
| KILL restrict (1)
< error
| 1064 42000 …line 1 column 13 near "restrict (1)" 
> exec ^
: string
| SET CONFIG tikv restrict = 1
< error
| 1064 42000 …line 1 column 24 near "restrict = 1" 
> exec ^
: string
| SET CONFIG tikv a-restrict = 1
< error
| 1064 42000 …line 1 column 26 near "restrict = 1" 
> exec ^
: string
| KILL resume()
< error
| 1064 42000 …line 1 column 11 near "resume()" 
> exec ^
: string
| KILL resume (1)
< error
| 1064 42000 …line 1 column 11 near "resume (1)" 
> exec ^
: string
| SET CONFIG tikv resume = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-resume = 1
< ok 0 0
> exec ^
: string
| KILL retain()
< ok 0 0
> exec ^
: string
| KILL retain (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv retain = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-retain = 1
< ok 0 0
> exec ^
: string
| KILL return()
< ok 0 0
> exec ^
: string
| KILL return (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv return = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-return = 1
< ok 0 0
> exec ^
: string
| KILL returned_sqlstate()
< ok 0 0
> exec ^
: string
| KILL returned_sqlstate (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv returned_sqlstate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-returned_sqlstate = 1
< ok 0 0
> exec ^
: string
| KILL returning()
< ok 0 0
> exec ^
: string
| KILL returning (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv returning = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-returning = 1
< ok 0 0
> exec ^
: string
| KILL returns()
< ok 0 0
> exec ^
: string
| KILL returns (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv returns = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-returns = 1
< ok 0 0
> exec ^
: string
| KILL reuse()
< error
| 1064 42000 …line 1 column 10 near "reuse()" 
> exec ^
: string
| KILL reuse (1)
< error
| 1064 42000 …line 1 column 10 near "reuse (1)" 
> exec ^
: string
| SET CONFIG tikv reuse = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-reuse = 1
< ok 0 0
> exec ^
: string
| KILL reverse()
< error
| 1064 42000 …line 1 column 12 near "reverse()" 
> exec ^
: string
| KILL reverse (1)
< error
| 1064 42000 …line 1 column 12 near "reverse (1)" 
> exec ^
: string
| SET CONFIG tikv reverse = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-reverse = 1
< ok 0 0
> exec ^
: string
| KILL revoke()
< error
| 1064 42000 …line 1 column 11 near "revoke()" 
> exec ^
: string
| KILL revoke (1)
< error
| 1064 42000 …line 1 column 11 near "revoke (1)" 
> exec ^
: string
| SET CONFIG tikv revoke = 1
< error
| 1064 42000 …line 1 column 22 near "revoke = 1" 
> exec ^
: string
| SET CONFIG tikv a-revoke = 1
< error
| 1064 42000 …line 1 column 24 near "revoke = 1" 
> exec ^
: string
| KILL right()
< error
| 1064 42000 …line 1 column 10 near "right()" 
> exec ^
: string
| KILL right (1)
< error
| 1064 42000 …line 1 column 10 near "right (1)" 
> exec ^
: string
| SET CONFIG tikv right = 1
< error
| 1064 42000 …line 1 column 21 near "right = 1" 
> exec ^
: string
| SET CONFIG tikv a-right = 1
< error
| 1064 42000 …line 1 column 23 near "right = 1" 
> exec ^
: string
| KILL rlike()
< error
| 1064 42000 …line 1 column 10 near "rlike()" 
> exec ^
: string
| KILL rlike (1)
< error
| 1064 42000 …line 1 column 10 near "rlike (1)" 
> exec ^
: string
| SET CONFIG tikv rlike = 1
< error
| 1064 42000 …line 1 column 21 near "rlike = 1" 
> exec ^
: string
| SET CONFIG tikv a-rlike = 1
< error
| 1064 42000 …line 1 column 23 near "rlike = 1" 
> exec ^
: string
| KILL role()
< error
| 1064 42000 …line 1 column 9 near "role()" 
> exec ^
: string
| KILL role (1)
< error
| 1064 42000 …line 1 column 9 near "role (1)" 
> exec ^
: string
| SET CONFIG tikv role = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-role = 1
< ok 0 0
> exec ^
: string
| KILL rollback()
< error
| 1064 42000 …line 1 column 13 near "rollback()" 
> exec ^
: string
| KILL rollback (1)
< error
| 1064 42000 …line 1 column 13 near "rollback (1)" 
> exec ^
: string
| SET CONFIG tikv rollback = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-rollback = 1
< ok 0 0
> exec ^
: string
| KILL rollup()
< error
| 1064 42000 …line 1 column 11 near "rollup()" 
> exec ^
: string
| KILL rollup (1)
< error
| 1064 42000 …line 1 column 11 near "rollup (1)" 
> exec ^
: string
| SET CONFIG tikv rollup = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-rollup = 1
< ok 0 0
> exec ^
: string
| KILL rotate()
< ok 0 0
> exec ^
: string
| KILL rotate (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv rotate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-rotate = 1
< ok 0 0
> exec ^
: string
| KILL routine()
< error
| 1064 42000 …line 1 column 12 near "routine()" 
> exec ^
: string
| KILL routine (1)
< error
| 1064 42000 …line 1 column 12 near "routine (1)" 
> exec ^
: string
| SET CONFIG tikv routine = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-routine = 1
< ok 0 0
> exec ^
: string
| KILL row()
< error
| 1064 42000 …line 1 column 8 near "row()" 
> exec ^
: string
| KILL row (1)
< error
| 1064 42000 …line 1 column 8 near "row (1)" 
> exec ^
: string
| SET CONFIG tikv row = 1
< error
| 1064 42000 …line 1 column 19 near "row = 1" 
> exec ^
: string
| SET CONFIG tikv a-row = 1
< error
| 1064 42000 …line 1 column 21 near "row = 1" 
> exec ^
: string
| KILL row_count()
< error
| 1064 42000 …line 1 column 14 near "row_count()" 
> exec ^
: string
| KILL row_count (1)
< error
| 1064 42000 …line 1 column 14 near "row_count (1)" 
> exec ^
: string
| SET CONFIG tikv row_count = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-row_count = 1
< ok 0 0
> exec ^
: string
| KILL row_format()
< error
| 1064 42000 …line 1 column 15 near "row_format()" 
> exec ^
: string
| KILL row_format (1)
< error
| 1064 42000 …line 1 column 15 near "row_format (1)" 
> exec ^
: string
| SET CONFIG tikv row_format = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-row_format = 1
< ok 0 0
> exec ^
: string
| KILL row_number()
< error
| 1064 42000 …line 1 column 15 near "row_number()" 
> exec ^
: string
| KILL row_number (1)
< error
| 1064 42000 …line 1 column 15 near "row_number (1)" 
> exec ^
: string
| SET CONFIG tikv row_number = 1
< error
| 1064 42000 …line 1 column 26 near "row_number = 1" 
> exec ^
: string
| SET CONFIG tikv a-row_number = 1
< error
| 1064 42000 …line 1 column 28 near "row_number = 1" 
> exec ^
: string
| KILL rows()
< error
| 1064 42000 …line 1 column 9 near "rows()" 
> exec ^
: string
| KILL rows (1)
< error
| 1064 42000 …line 1 column 9 near "rows (1)" 
> exec ^
: string
| SET CONFIG tikv rows = 1
< error
| 1064 42000 …line 1 column 20 near "rows = 1" 
> exec ^
: string
| SET CONFIG tikv a-rows = 1
< error
| 1064 42000 …line 1 column 22 near "rows = 1" 
> exec ^
: string
| KILL rtree()
< error
| 1064 42000 …line 1 column 10 near "rtree()" 
> exec ^
: string
| KILL rtree (1)
< error
| 1064 42000 …line 1 column 10 near "rtree (1)" 
> exec ^
: string
| SET CONFIG tikv rtree = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-rtree = 1
< ok 0 0
> exec ^
: string
| KILL ru()
< error
| 1064 42000 …line 1 column 7 near "ru()" 
> exec ^
: string
| KILL ru (1)
< error
| 1064 42000 …line 1 column 7 near "ru (1)" 
> exec ^
: string
| SET CONFIG tikv ru = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ru = 1
< ok 0 0
> exec ^
: string
| KILL ru_per_sec()
< error
| 1064 42000 …line 1 column 15 near "ru_per_sec()" 
> exec ^
: string
| KILL ru_per_sec (1)
< error
| 1064 42000 …line 1 column 15 near "ru_per_sec (1)" 
> exec ^
: string
| SET CONFIG tikv ru_per_sec = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ru_per_sec = 1
< ok 0 0
> exec ^
: string
| KILL rule()
< error
| 1064 42000 …line 1 column 9 near "rule()" 
> exec ^
: string
| KILL rule (1)
< error
| 1064 42000 …line 1 column 9 near "rule (1)" 
> exec ^
: string
| SET CONFIG tikv rule = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-rule = 1
< ok 0 0
> exec ^
: string
| KILL run()
< error
| 1064 42000 …line 1 column 8 near "run()" 
> exec ^
: string
| KILL run (1)
< error
| 1064 42000 …line 1 column 8 near "run (1)" 
> exec ^
: string
| SET CONFIG tikv run = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-run = 1
< ok 0 0
> exec ^
: string
| KILL running()
< error
| 1064 42000 …line 1 column 12 near "running()" 
> exec ^
: string
| KILL running (1)
< error
| 1064 42000 …line 1 column 12 near "running (1)" 
> exec ^
: string
| SET CONFIG tikv running = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-running = 1
< ok 0 0
> exec ^
: string
| KILL s3()
< error
| 1064 42000 …line 1 column 7 near "s3()" 
> exec ^
: string
| KILL s3 (1)
< error
| 1064 42000 …line 1 column 7 near "s3 (1)" 
> exec ^
: string
| SET CONFIG tikv s3 = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-s3 = 1
< ok 0 0
> exec ^
: string
| KILL samplerate()
< error
| 1064 42000 …line 1 column 15 near "samplerate()" 
> exec ^
: string
| KILL samplerate (1)
< error
| 1064 42000 …line 1 column 15 near "samplerate (1)" 
> exec ^
: string
| SET CONFIG tikv samplerate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-samplerate = 1
< ok 0 0
> exec ^
: string
| KILL samples()
< error
| 1064 42000 …line 1 column 12 near "samples()" 
> exec ^
: string
| KILL samples (1)
< error
| 1064 42000 …line 1 column 12 near "samples (1)" 
> exec ^
: string
| SET CONFIG tikv samples = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-samples = 1
< ok 0 0
> exec ^
: string
| KILL san()
< error
| 1064 42000 …line 1 column 8 near "san()" 
> exec ^
: string
| KILL san (1)
< error
| 1064 42000 …line 1 column 8 near "san (1)" 
> exec ^
: string
| SET CONFIG tikv san = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-san = 1
< ok 0 0
> exec ^
: string
| KILL savepoint()
< error
| 1064 42000 …line 1 column 14 near "savepoint()" 
> exec ^
: string
| KILL savepoint (1)
< error
| 1064 42000 …line 1 column 14 near "savepoint (1)" 
> exec ^
: string
| SET CONFIG tikv savepoint = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-savepoint = 1
< ok 0 0
> exec ^
: string
| KILL schedule()
< error
| 1064 42000 …line 1 column 13 near "schedule()" 
> exec ^
: string
| KILL schedule (1)
< error
| 1064 42000 …line 1 column 13 near "schedule (1)" 
> exec ^
: string
| SET CONFIG tikv schedule = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-schedule = 1
< ok 0 0
> exec ^
: string
| KILL schema()
< error
| 1064 42000 …line 1 column 11 near "schema()" 
> exec ^
: string
| KILL schema (1)
< error
| 1064 42000 …line 1 column 11 near "schema (1)" 
> exec ^
: string
| SET CONFIG tikv schema = 1
< error
| 1064 42000 …line 1 column 22 near "schema = 1" 
> exec ^
: string
| SET CONFIG tikv a-schema = 1
< error
| 1064 42000 …line 1 column 24 near "schema = 1" 
> exec ^
: string
| KILL schema_name()
< ok 0 0
> exec ^
: string
| KILL schema_name (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv schema_name = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-schema_name = 1
< ok 0 0
> exec ^
: string
| KILL schemas()
< error
| 1064 42000 …line 1 column 12 near "schemas()" 
> exec ^
: string
| KILL schemas (1)
< error
| 1064 42000 …line 1 column 12 near "schemas (1)" 
> exec ^
: string
| SET CONFIG tikv schemas = 1
< error
| 1064 42000 …line 1 column 23 near "schemas = 1" 
> exec ^
: string
| SET CONFIG tikv a-schemas = 1
< error
| 1064 42000 …line 1 column 25 near "schemas = 1" 
> exec ^
: string
| KILL second()
< error
| 1064 42000 …line 1 column 11 near "second()" 
> exec ^
: string
| KILL second (1)
< error
| 1064 42000 …line 1 column 11 near "second (1)" 
> exec ^
: string
| SET CONFIG tikv second = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-second = 1
< ok 0 0
> exec ^
: string
| KILL second_microsecond()
< error
| 1064 42000 …line 1 column 23 near "second_microsecond()" 
> exec ^
: string
| KILL second_microsecond (1)
< error
| 1064 42000 …line 1 column 23 near "second_microsecond (1)" 
> exec ^
: string
| SET CONFIG tikv second_microsecond = 1
< error
| 1064 42000 …line 1 column 34 near "second_microsecond = 1" 
> exec ^
: string
| SET CONFIG tikv a-second_microsecond = 1
< error
| 1064 42000 …line 1 column 36 near "second_microsecond = 1" 
> exec ^
: string
| KILL secondary()
< error
| 1064 42000 …line 1 column 14 near "secondary()" 
> exec ^
: string
| KILL secondary (1)
< error
| 1064 42000 …line 1 column 14 near "secondary (1)" 
> exec ^
: string
| SET CONFIG tikv secondary = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-secondary = 1
< ok 0 0
> exec ^
: string
| KILL secondary_engine()
< error
| 1064 42000 …line 1 column 21 near "secondary_engine()" 
> exec ^
: string
| KILL secondary_engine (1)
< error
| 1064 42000 …line 1 column 21 near "secondary_engine (1)" 
> exec ^
: string
| SET CONFIG tikv secondary_engine = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-secondary_engine = 1
< ok 0 0
> exec ^
: string
| KILL secondary_engine_attribute()
< ok 0 0
> exec ^
: string
| KILL secondary_engine_attribute (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv secondary_engine_attribute = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-secondary_engine_attribute = 1
< ok 0 0
> exec ^
: string
| KILL secondary_load()
< error
| 1064 42000 …line 1 column 19 near "secondary_load()" 
> exec ^
: string
| KILL secondary_load (1)
< error
| 1064 42000 …line 1 column 19 near "secondary_load (1)" 
> exec ^
: string
| SET CONFIG tikv secondary_load = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-secondary_load = 1
< ok 0 0
> exec ^
: string
| KILL secondary_unload()
< error
| 1064 42000 …line 1 column 21 near "secondary_unload()" 
> exec ^
: string
| KILL secondary_unload (1)
< error
| 1064 42000 …line 1 column 21 near "secondary_unload (1)" 
> exec ^
: string
| SET CONFIG tikv secondary_unload = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-secondary_unload = 1
< ok 0 0
> exec ^
: string
| KILL security()
< error
| 1064 42000 …line 1 column 13 near "security()" 
> exec ^
: string
| KILL security (1)
< error
| 1064 42000 …line 1 column 13 near "security (1)" 
> exec ^
: string
| SET CONFIG tikv security = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-security = 1
< ok 0 0
> exec ^
: string
| KILL select()
< error
| 1064 42000 …line 1 column 11 near "select()" 
> exec ^
: string
| KILL select (1)
< error
| 1064 42000 …line 1 column 11 near "select (1)" 
> exec ^
: string
| SET CONFIG tikv select = 1
< error
| 1064 42000 …line 1 column 22 near "select = 1" 
> exec ^
: string
| SET CONFIG tikv a-select = 1
< error
| 1064 42000 …line 1 column 24 near "select = 1" 
> exec ^
: string
| KILL send_credentials_to_tikv()
< error
| 1064 42000 …line 1 column 29 near "send_credentials_to_tikv()" 
> exec ^
: string
| KILL send_credentials_to_tikv (1)
< error
| 1064 42000 …line 1 column 29 near "send_credentials_to_tikv (1)" 
> exec ^
: string
| SET CONFIG tikv send_credentials_to_tikv = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-send_credentials_to_tikv = 1
< ok 0 0
> exec ^
: string
| KILL sensitive()
< ok 0 0
> exec ^
: string
| KILL sensitive (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv sensitive = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sensitive = 1
< ok 0 0
> exec ^
: string
| KILL separator()
< error
| 1064 42000 …line 1 column 14 near "separator()" 
> exec ^
: string
| KILL separator (1)
< error
| 1064 42000 …line 1 column 14 near "separator (1)" 
> exec ^
: string
| SET CONFIG tikv separator = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-separator = 1
< ok 0 0
> exec ^
: string
| KILL sequence()
< error
| 1064 42000 …line 1 column 13 near "sequence()" 
> exec ^
: string
| KILL sequence (1)
< error
| 1064 42000 …line 1 column 13 near "sequence (1)" 
> exec ^
: string
| SET CONFIG tikv sequence = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sequence = 1
< ok 0 0
> exec ^
: string
| KILL serial()
< error
| 1064 42000 …line 1 column 11 near "serial()" 
> exec ^
: string
| KILL serial (1)
< error
| 1064 42000 …line 1 column 11 near "serial (1)" 
> exec ^
: string
| SET CONFIG tikv serial = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-serial = 1
< ok 0 0
> exec ^
: string
| KILL serializable()
< error
| 1064 42000 …line 1 column 17 near "serializable()" 
> exec ^
: string
| KILL serializable (1)
< error
| 1064 42000 …line 1 column 17 near "serializable (1)" 
> exec ^
: string
| SET CONFIG tikv serializable = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-serializable = 1
< ok 0 0
> exec ^
: string
| KILL server()
< ok 0 0
> exec ^
: string
| KILL server (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv server = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-server = 1
< ok 0 0
> exec ^
: string
| KILL session()
< error
| 1064 42000 …line 1 column 12 near "session()" 
> exec ^
: string
| KILL session (1)
< error
| 1064 42000 …line 1 column 12 near "session (1)" 
> exec ^
: string
| SET CONFIG tikv session = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-session = 1
< ok 0 0
> exec ^
: string
| KILL session_states()
< error
| 1064 42000 …line 1 column 19 near "session_states()" 
> exec ^
: string
| KILL session_states (1)
< error
| 1064 42000 …line 1 column 19 near "session_states (1)" 
> exec ^
: string
| SET CONFIG tikv session_states = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-session_states = 1
< ok 0 0
> exec ^
: string
| KILL session_user()
< ok 0 0
> exec ^
: string
| KILL session_user (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv session_user = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-session_user = 1
< ok 0 0
> exec ^
: string
| KILL set()
< error
| 1064 42000 …line 1 column 8 near "set()" 
> exec ^
: string
| KILL set (1)
< error
| 1064 42000 …line 1 column 8 near "set (1)" 
> exec ^
: string
| SET CONFIG tikv set = 1
< error
| 1064 42000 …line 1 column 19 near "set = 1" 
> exec ^
: string
| SET CONFIG tikv a-set = 1
< error
| 1064 42000 …line 1 column 21 near "set = 1" 
> exec ^
: string
| KILL sets()
< ok 0 0
> exec ^
: string
| KILL sets (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv sets = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sets = 1
< ok 0 0
> exec ^
: string
| KILL setval()
< error
| 1064 42000 …line 1 column 11 near "setval()" 
> exec ^
: string
| KILL setval (1)
< error
| 1064 42000 …line 1 column 11 near "setval (1)" 
> exec ^
: string
| SET CONFIG tikv setval = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-setval = 1
< ok 0 0
> exec ^
: string
| KILL shard_row_id_bits()
< error
| 1064 42000 …line 1 column 22 near "shard_row_id_bits()" 
> exec ^
: string
| KILL shard_row_id_bits (1)
< error
| 1064 42000 …line 1 column 22 near "shard_row_id_bits (1)" 
> exec ^
: string
| SET CONFIG tikv shard_row_id_bits = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-shard_row_id_bits = 1
< ok 0 0
> exec ^
: string
| KILL share()
< error
| 1064 42000 …line 1 column 10 near "share()" 
> exec ^
: string
| KILL share (1)
< error
| 1064 42000 …line 1 column 10 near "share (1)" 
> exec ^
: string
| SET CONFIG tikv share = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-share = 1
< ok 0 0
> exec ^
: string
| KILL shared()
< error
| 1064 42000 …line 1 column 11 near "shared()" 
> exec ^
: string
| KILL shared (1)
< error
| 1064 42000 …line 1 column 11 near "shared (1)" 
> exec ^
: string
| SET CONFIG tikv shared = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-shared = 1
< ok 0 0
> exec ^
: string
| KILL show()
< error
| 1064 42000 …line 1 column 9 near "show()" 
> exec ^
: string
| KILL show (1)
< error
| 1064 42000 …line 1 column 9 near "show (1)" 
> exec ^
: string
| SET CONFIG tikv show = 1
< error
| 1064 42000 …line 1 column 20 near "show = 1" 
> exec ^
: string
| SET CONFIG tikv a-show = 1
< error
| 1064 42000 …line 1 column 22 near "show = 1" 
> exec ^
: string
| KILL shutdown()
< error
| 1064 42000 …line 1 column 13 near "shutdown()" 
> exec ^
: string
| KILL shutdown (1)
< error
| 1064 42000 …line 1 column 13 near "shutdown (1)" 
> exec ^
: string
| SET CONFIG tikv shutdown = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-shutdown = 1
< ok 0 0
> exec ^
: string
| KILL signal()
< ok 0 0
> exec ^
: string
| KILL signal (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv signal = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-signal = 1
< ok 0 0
> exec ^
: string
| KILL signed()
< error
| 1064 42000 …line 1 column 11 near "signed()" 
> exec ^
: string
| KILL signed (1)
< error
| 1064 42000 …line 1 column 11 near "signed (1)" 
> exec ^
: string
| SET CONFIG tikv signed = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-signed = 1
< ok 0 0
> exec ^
: string
| KILL similar()
< error
| 1064 42000 …line 1 column 12 near "similar()" 
> exec ^
: string
| KILL similar (1)
< error
| 1064 42000 …line 1 column 12 near "similar (1)" 
> exec ^
: string
| SET CONFIG tikv similar = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-similar = 1
< ok 0 0
> exec ^
: string
| KILL simple()
< error
| 1064 42000 …line 1 column 11 near "simple()" 
> exec ^
: string
| KILL simple (1)
< error
| 1064 42000 …line 1 column 11 near "simple (1)" 
> exec ^
: string
| SET CONFIG tikv simple = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-simple = 1
< ok 0 0
> exec ^
: string
| KILL skip()
< error
| 1064 42000 …line 1 column 9 near "skip()" 
> exec ^
: string
| KILL skip (1)
< error
| 1064 42000 …line 1 column 9 near "skip (1)" 
> exec ^
: string
| SET CONFIG tikv skip = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-skip = 1
< ok 0 0
> exec ^
: string
| KILL skip_schema_files()
< error
| 1064 42000 …line 1 column 22 near "skip_schema_files()" 
> exec ^
: string
| KILL skip_schema_files (1)
< error
| 1064 42000 …line 1 column 22 near "skip_schema_files (1)" 
> exec ^
: string
| SET CONFIG tikv skip_schema_files = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-skip_schema_files = 1
< ok 0 0
> exec ^
: string
| KILL slave()
< error
| 1064 42000 …line 1 column 10 near "slave()" 
> exec ^
: string
| KILL slave (1)
< error
| 1064 42000 …line 1 column 10 near "slave (1)" 
> exec ^
: string
| SET CONFIG tikv slave = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-slave = 1
< ok 0 0
> exec ^
: string
| KILL slow()
< error
| 1064 42000 …line 1 column 9 near "slow()" 
> exec ^
: string
| KILL slow (1)
< error
| 1064 42000 …line 1 column 9 near "slow (1)" 
> exec ^
: string
| SET CONFIG tikv slow = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-slow = 1
< ok 0 0
> exec ^
: string
| KILL smallint()
< error
| 1064 42000 …line 1 column 13 near "smallint()" 
> exec ^
: string
| KILL smallint (1)
< error
| 1064 42000 …line 1 column 13 near "smallint (1)" 
> exec ^
: string
| SET CONFIG tikv smallint = 1
< error
| 1064 42000 …line 1 column 24 near "smallint = 1" 
> exec ^
: string
| SET CONFIG tikv a-smallint = 1
< error
| 1064 42000 …line 1 column 26 near "smallint = 1" 
> exec ^
: string
| KILL snapshot()
< error
| 1064 42000 …line 1 column 13 near "snapshot()" 
> exec ^
: string
| KILL snapshot (1)
< error
| 1064 42000 …line 1 column 13 near "snapshot (1)" 
> exec ^
: string
| SET CONFIG tikv snapshot = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-snapshot = 1
< ok 0 0
> exec ^
: string
| KILL socket()
< ok 0 0
> exec ^
: string
| KILL socket (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv socket = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-socket = 1
< ok 0 0
> exec ^
: string
| KILL some()
< error
| 1064 42000 …line 1 column 9 near "some()" 
> exec ^
: string
| KILL some (1)
< error
| 1064 42000 …line 1 column 9 near "some (1)" 
> exec ^
: string
| SET CONFIG tikv some = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-some = 1
< ok 0 0
> exec ^
: string
| KILL soname()
< ok 0 0
> exec ^
: string
| KILL soname (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv soname = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-soname = 1
< ok 0 0
> exec ^
: string
| KILL sounds()
< ok 0 0
> exec ^
: string
| KILL sounds (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv sounds = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sounds = 1
< ok 0 0
> exec ^
: string
| KILL source()
< error
| 1064 42000 …line 1 column 11 near "source()" 
> exec ^
: string
| KILL source (1)
< error
| 1064 42000 …line 1 column 11 near "source (1)" 
> exec ^
: string
| SET CONFIG tikv source = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source = 1
< ok 0 0
> exec ^
: string
| KILL source_auto_position()
< ok 0 0
> exec ^
: string
| KILL source_auto_position (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_auto_position = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_auto_position = 1
< ok 0 0
> exec ^
: string
| KILL source_bind()
< ok 0 0
> exec ^
: string
| KILL source_bind (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_bind = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_bind = 1
< ok 0 0
> exec ^
: string
| KILL source_compression_algorithms()
< ok 0 0
> exec ^
: string
| KILL source_compression_algorithms (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_compression_algorithms = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_compression_algorithms = 1
< ok 0 0
> exec ^
: string
| KILL source_connect_retry()
< ok 0 0
> exec ^
: string
| KILL source_connect_retry (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_connect_retry = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_connect_retry = 1
< ok 0 0
> exec ^
: string
| KILL source_connection_auto_failover()
< ok 0 0
> exec ^
: string
| KILL source_connection_auto_failover (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_connection_auto_failover = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_connection_auto_failover = 1
< ok 0 0
> exec ^
: string
| KILL source_delay()
< ok 0 0
> exec ^
: string
| KILL source_delay (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_delay = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_delay = 1
< ok 0 0
> exec ^
: string
| KILL source_heartbeat_period()
< ok 0 0
> exec ^
: string
| KILL source_heartbeat_period (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_heartbeat_period = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_heartbeat_period = 1
< ok 0 0
> exec ^
: string
| KILL source_host()
< ok 0 0
> exec ^
: string
| KILL source_host (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_host = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_host = 1
< ok 0 0
> exec ^
: string
| KILL source_log_file()
< ok 0 0
> exec ^
: string
| KILL source_log_file (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_log_file = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_log_file = 1
< ok 0 0
> exec ^
: string
| KILL source_log_pos()
< ok 0 0
> exec ^
: string
| KILL source_log_pos (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_log_pos = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_log_pos = 1
< ok 0 0
> exec ^
: string
| KILL source_password()
< ok 0 0
> exec ^
: string
| KILL source_password (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_password = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_password = 1
< ok 0 0
> exec ^
: string
| KILL source_port()
< ok 0 0
> exec ^
: string
| KILL source_port (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_port = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_port = 1
< ok 0 0
> exec ^
: string
| KILL source_public_key_path()
< ok 0 0
> exec ^
: string
| KILL source_public_key_path (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_public_key_path = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_public_key_path = 1
< ok 0 0
> exec ^
: string
| KILL source_retry_count()
< ok 0 0
> exec ^
: string
| KILL source_retry_count (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_retry_count = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_retry_count = 1
< ok 0 0
> exec ^
: string
| KILL source_ssl()
< ok 0 0
> exec ^
: string
| KILL source_ssl (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_ssl = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_ssl = 1
< ok 0 0
> exec ^
: string
| KILL source_ssl_ca()
< ok 0 0
> exec ^
: string
| KILL source_ssl_ca (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_ssl_ca = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_ssl_ca = 1
< ok 0 0
> exec ^
: string
| KILL source_ssl_capath()
< ok 0 0
> exec ^
: string
| KILL source_ssl_capath (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_ssl_capath = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_ssl_capath = 1
< ok 0 0
> exec ^
: string
| KILL source_ssl_cert()
< ok 0 0
> exec ^
: string
| KILL source_ssl_cert (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_ssl_cert = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_ssl_cert = 1
< ok 0 0
> exec ^
: string
| KILL source_ssl_cipher()
< ok 0 0
> exec ^
: string
| KILL source_ssl_cipher (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_ssl_cipher = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_ssl_cipher = 1
< ok 0 0
> exec ^
: string
| KILL source_ssl_crl()
< ok 0 0
> exec ^
: string
| KILL source_ssl_crl (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_ssl_crl = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_ssl_crl = 1
< ok 0 0
> exec ^
: string
| KILL source_ssl_crlpath()
< ok 0 0
> exec ^
: string
| KILL source_ssl_crlpath (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_ssl_crlpath = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_ssl_crlpath = 1
< ok 0 0
> exec ^
: string
| KILL source_ssl_key()
< ok 0 0
> exec ^
: string
| KILL source_ssl_key (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_ssl_key = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_ssl_key = 1
< ok 0 0
> exec ^
: string
| KILL source_ssl_verify_server_cert()
< ok 0 0
> exec ^
: string
| KILL source_ssl_verify_server_cert (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_ssl_verify_server_cert = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_ssl_verify_server_cert = 1
< ok 0 0
> exec ^
: string
| KILL source_tls_ciphersuites()
< ok 0 0
> exec ^
: string
| KILL source_tls_ciphersuites (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_tls_ciphersuites = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_tls_ciphersuites = 1
< ok 0 0
> exec ^
: string
| KILL source_tls_version()
< ok 0 0
> exec ^
: string
| KILL source_tls_version (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_tls_version = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_tls_version = 1
< ok 0 0
> exec ^
: string
| KILL source_user()
< ok 0 0
> exec ^
: string
| KILL source_user (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_user = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_user = 1
< ok 0 0
> exec ^
: string
| KILL source_zstd_compression_level()
< ok 0 0
> exec ^
: string
| KILL source_zstd_compression_level (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv source_zstd_compression_level = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-source_zstd_compression_level = 1
< ok 0 0
> exec ^
: string
| KILL spatial()
< error
| 1064 42000 …line 1 column 12 near "spatial()" 
> exec ^
: string
| KILL spatial (1)
< error
| 1064 42000 …line 1 column 12 near "spatial (1)" 
> exec ^
: string
| SET CONFIG tikv spatial = 1
< error
| 1064 42000 …line 1 column 23 near "spatial = 1" 
> exec ^
: string
| SET CONFIG tikv a-spatial = 1
< error
| 1064 42000 …line 1 column 25 near "spatial = 1" 
> exec ^
: string
| KILL specific()
< ok 0 0
> exec ^
: string
| KILL specific (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv specific = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-specific = 1
< ok 0 0
> exec ^
: string
| KILL split()
< error
| 1064 42000 …line 1 column 10 near "split()" 
> exec ^
: string
| KILL split (1)
< error
| 1064 42000 …line 1 column 10 near "split (1)" 
> exec ^
: string
| SET CONFIG tikv split = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-split = 1
< ok 0 0
> exec ^
: string
| KILL sql()
< error
| 1064 42000 …line 1 column 8 near "sql()" 
> exec ^
: string
| KILL sql (1)
< error
| 1064 42000 …line 1 column 8 near "sql (1)" 
> exec ^
: string
| SET CONFIG tikv sql = 1
< error
| 1064 42000 …line 1 column 19 near "sql = 1" 
> exec ^
: string
| SET CONFIG tikv a-sql = 1
< error
| 1064 42000 …line 1 column 21 near "sql = 1" 
> exec ^
: string
| KILL sql_after_gtids()
< ok 0 0
> exec ^
: string
| KILL sql_after_gtids (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv sql_after_gtids = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_after_gtids = 1
< ok 0 0
> exec ^
: string
| KILL sql_after_mts_gaps()
< ok 0 0
> exec ^
: string
| KILL sql_after_mts_gaps (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv sql_after_mts_gaps = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_after_mts_gaps = 1
< ok 0 0
> exec ^
: string
| KILL sql_before_gtids()
< ok 0 0
> exec ^
: string
| KILL sql_before_gtids (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv sql_before_gtids = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_before_gtids = 1
< ok 0 0
> exec ^
: string
| KILL sql_big_result()
< error
| 1064 42000 …line 1 column 19 near "sql_big_result()" 
> exec ^
: string
| KILL sql_big_result (1)
< error
| 1064 42000 …line 1 column 19 near "sql_big_result (1)" 
> exec ^
: string
| SET CONFIG tikv sql_big_result = 1
< error
| 1064 42000 …line 1 column 30 near "sql_big_result = 1" 
> exec ^
: string
| SET CONFIG tikv a-sql_big_result = 1
< error
| 1064 42000 …line 1 column 32 near "sql_big_result = 1" 
> exec ^
: string
| KILL sql_buffer_result()
< error
| 1064 42000 …line 1 column 22 near "sql_buffer_result()" 
> exec ^
: string
| KILL sql_buffer_result (1)
< error
| 1064 42000 …line 1 column 22 near "sql_buffer_result (1)" 
> exec ^
: string
| SET CONFIG tikv sql_buffer_result = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_buffer_result = 1
< ok 0 0
> exec ^
: string
| KILL sql_cache()
< error
| 1064 42000 …line 1 column 14 near "sql_cache()" 
> exec ^
: string
| KILL sql_cache (1)
< error
| 1064 42000 …line 1 column 14 near "sql_cache (1)" 
> exec ^
: string
| SET CONFIG tikv sql_cache = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_cache = 1
< ok 0 0
> exec ^
: string
| KILL sql_calc_found_rows()
< error
| 1064 42000 …line 1 column 24 near "sql_calc_found_rows()" 
> exec ^
: string
| KILL sql_calc_found_rows (1)
< error
| 1064 42000 …line 1 column 24 near "sql_calc_found_rows (1)" 
> exec ^
: string
| SET CONFIG tikv sql_calc_found_rows = 1
< error
| 1064 42000 …line 1 column 35 near "sql_calc_found_rows = 1" 
> exec ^
: string
| SET CONFIG tikv a-sql_calc_found_rows = 1
< error
| 1064 42000 …line 1 column 37 near "sql_calc_found_rows = 1" 
> exec ^
: string
| KILL sql_no_cache()
< error
| 1064 42000 …line 1 column 17 near "sql_no_cache()" 
> exec ^
: string
| KILL sql_no_cache (1)
< error
| 1064 42000 …line 1 column 17 near "sql_no_cache (1)" 
> exec ^
: string
| SET CONFIG tikv sql_no_cache = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_no_cache = 1
< ok 0 0
> exec ^
: string
| KILL sql_small_result()
< error
| 1064 42000 …line 1 column 21 near "sql_small_result()" 
> exec ^
: string
| KILL sql_small_result (1)
< error
| 1064 42000 …line 1 column 21 near "sql_small_result (1)" 
> exec ^
: string
| SET CONFIG tikv sql_small_result = 1
< error
| 1064 42000 …line 1 column 32 near "sql_small_result = 1" 
> exec ^
: string
| SET CONFIG tikv a-sql_small_result = 1
< error
| 1064 42000 …line 1 column 34 near "sql_small_result = 1" 
> exec ^
: string
| KILL sql_thread()
< ok 0 0
> exec ^
: string
| KILL sql_thread (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv sql_thread = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_thread = 1
< ok 0 0
> exec ^
: string
| KILL sql_tsi_day()
< error
| 1064 42000 …line 1 column 16 near "sql_tsi_day()" 
> exec ^
: string
| KILL sql_tsi_day (1)
< error
| 1064 42000 …line 1 column 16 near "sql_tsi_day (1)" 
> exec ^
: string
| SET CONFIG tikv sql_tsi_day = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_tsi_day = 1
< ok 0 0
> exec ^
: string
| KILL sql_tsi_hour()
< error
| 1064 42000 …line 1 column 17 near "sql_tsi_hour()" 
> exec ^
: string
| KILL sql_tsi_hour (1)
< error
| 1064 42000 …line 1 column 17 near "sql_tsi_hour (1)" 
> exec ^
: string
| SET CONFIG tikv sql_tsi_hour = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_tsi_hour = 1
< ok 0 0
> exec ^
: string
| KILL sql_tsi_minute()
< error
| 1064 42000 …line 1 column 19 near "sql_tsi_minute()" 
> exec ^
: string
| KILL sql_tsi_minute (1)
< error
| 1064 42000 …line 1 column 19 near "sql_tsi_minute (1)" 
> exec ^
: string
| SET CONFIG tikv sql_tsi_minute = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_tsi_minute = 1
< ok 0 0
> exec ^
: string
| KILL sql_tsi_month()
< error
| 1064 42000 …line 1 column 18 near "sql_tsi_month()" 
> exec ^
: string
| KILL sql_tsi_month (1)
< error
| 1064 42000 …line 1 column 18 near "sql_tsi_month (1)" 
> exec ^
: string
| SET CONFIG tikv sql_tsi_month = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_tsi_month = 1
< ok 0 0
> exec ^
: string
| KILL sql_tsi_quarter()
< error
| 1064 42000 …line 1 column 20 near "sql_tsi_quarter()" 
> exec ^
: string
| KILL sql_tsi_quarter (1)
< error
| 1064 42000 …line 1 column 20 near "sql_tsi_quarter (1)" 
> exec ^
: string
| SET CONFIG tikv sql_tsi_quarter = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_tsi_quarter = 1
< ok 0 0
> exec ^
: string
| KILL sql_tsi_second()
< error
| 1064 42000 …line 1 column 19 near "sql_tsi_second()" 
> exec ^
: string
| KILL sql_tsi_second (1)
< error
| 1064 42000 …line 1 column 19 near "sql_tsi_second (1)" 
> exec ^
: string
| SET CONFIG tikv sql_tsi_second = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_tsi_second = 1
< ok 0 0
> exec ^
: string
| KILL sql_tsi_week()
< error
| 1064 42000 …line 1 column 17 near "sql_tsi_week()" 
> exec ^
: string
| KILL sql_tsi_week (1)
< error
| 1064 42000 …line 1 column 17 near "sql_tsi_week (1)" 
> exec ^
: string
| SET CONFIG tikv sql_tsi_week = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_tsi_week = 1
< ok 0 0
> exec ^
: string
| KILL sql_tsi_year()
< error
| 1064 42000 …line 1 column 17 near "sql_tsi_year()" 
> exec ^
: string
| KILL sql_tsi_year (1)
< error
| 1064 42000 …line 1 column 17 near "sql_tsi_year (1)" 
> exec ^
: string
| SET CONFIG tikv sql_tsi_year = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sql_tsi_year = 1
< ok 0 0
> exec ^
: string
| KILL sqlexception()
< error
| 1064 42000 …line 1 column 17 near "sqlexception()" 
> exec ^
: string
| KILL sqlexception (1)
< error
| 1064 42000 …line 1 column 17 near "sqlexception (1)" 
> exec ^
: string
| SET CONFIG tikv sqlexception = 1
< error
| 1064 42000 …line 1 column 28 near "sqlexception = 1" 
> exec ^
: string
| SET CONFIG tikv a-sqlexception = 1
< error
| 1064 42000 …line 1 column 30 near "sqlexception = 1" 
> exec ^
: string
| KILL sqlstate()
< error
| 1064 42000 …line 1 column 13 near "sqlstate()" 
> exec ^
: string
| KILL sqlstate (1)
< error
| 1064 42000 …line 1 column 13 near "sqlstate (1)" 
> exec ^
: string
| SET CONFIG tikv sqlstate = 1
< error
| 1064 42000 …line 1 column 24 near "sqlstate = 1" 
> exec ^
: string
| SET CONFIG tikv a-sqlstate = 1
< error
| 1064 42000 …line 1 column 26 near "sqlstate = 1" 
> exec ^
: string
| KILL sqlwarning()
< error
| 1064 42000 …line 1 column 15 near "sqlwarning()" 
> exec ^
: string
| KILL sqlwarning (1)
< error
| 1064 42000 …line 1 column 15 near "sqlwarning (1)" 
> exec ^
: string
| SET CONFIG tikv sqlwarning = 1
< error
| 1064 42000 …line 1 column 26 near "sqlwarning = 1" 
> exec ^
: string
| SET CONFIG tikv a-sqlwarning = 1
< error
| 1064 42000 …line 1 column 28 near "sqlwarning = 1" 
> exec ^
: string
| KILL srid()
< ok 0 0
> exec ^
: string
| KILL srid (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv srid = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-srid = 1
< ok 0 0
> exec ^
: string
| KILL ssl()
< error
| 1064 42000 …line 1 column 8 near "ssl()" 
> exec ^
: string
| KILL ssl (1)
< error
| 1064 42000 …line 1 column 8 near "ssl (1)" 
> exec ^
: string
| SET CONFIG tikv ssl = 1
< error
| 1064 42000 …line 1 column 19 near "ssl = 1" 
> exec ^
: string
| SET CONFIG tikv a-ssl = 1
< error
| 1064 42000 …line 1 column 21 near "ssl = 1" 
> exec ^
: string
| KILL st_collect()
< ok 0 0
> exec ^
: string
| KILL st_collect (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv st_collect = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-st_collect = 1
< ok 0 0
> exec ^
: string
| KILL stacked()
< ok 0 0
> exec ^
: string
| KILL stacked (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv stacked = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stacked = 1
< ok 0 0
> exec ^
: string
| KILL staleness()
< error
| 1064 42000 …line 1 column 14 near "staleness()" 
> exec ^
: string
| KILL staleness (1)
< error
| 1064 42000 …line 1 column 14 near "staleness (1)" 
> exec ^
: string
| SET CONFIG tikv staleness = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-staleness = 1
< ok 0 0
> exec ^
: string
| KILL start()
< error
| 1064 42000 …line 1 column 10 near "start()" 
> exec ^
: string
| KILL start (1)
< error
| 1064 42000 …line 1 column 10 near "start (1)" 
> exec ^
: string
| SET CONFIG tikv start = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-start = 1
< ok 0 0
> exec ^
: string
| KILL start_time()
< error
| 1064 42000 …line 1 column 15 near "start_time()" 
> exec ^
: string
| KILL start_time (1)
< error
| 1064 42000 …line 1 column 15 near "start_time (1)" 
> exec ^
: string
| SET CONFIG tikv start_time = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-start_time = 1
< ok 0 0
> exec ^
: string
| KILL start_ts()
< error
| 1064 42000 …line 1 column 13 near "start_ts()" 
> exec ^
: string
| KILL start_ts (1)
< error
| 1064 42000 …line 1 column 13 near "start_ts (1)" 
> exec ^
: string
| SET CONFIG tikv start_ts = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-start_ts = 1
< ok 0 0
> exec ^
: string
| KILL starting()
< error
| 1064 42000 …line 1 column 13 near "starting()" 
> exec ^
: string
| KILL starting (1)
< error
| 1064 42000 …line 1 column 13 near "starting (1)" 
> exec ^
: string
| SET CONFIG tikv starting = 1
< error
| 1064 42000 …line 1 column 24 near "starting = 1" 
> exec ^
: string
| SET CONFIG tikv a-starting = 1
< error
| 1064 42000 …line 1 column 26 near "starting = 1" 
> exec ^
: string
| KILL starts()
< ok 0 0
> exec ^
: string
| KILL starts (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv starts = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-starts = 1
< ok 0 0
> exec ^
: string
| KILL statistics()
< error
| 1064 42000 …line 1 column 15 near "statistics()" 
> exec ^
: string
| KILL statistics (1)
< error
| 1064 42000 …line 1 column 15 near "statistics (1)" 
> exec ^
: string
| SET CONFIG tikv statistics = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-statistics = 1
< ok 0 0
> exec ^
: string
| KILL stats()
< error
| 1064 42000 …line 1 column 10 near "stats()" 
> exec ^
: string
| KILL stats (1)
< error
| 1064 42000 …line 1 column 10 near "stats (1)" 
> exec ^
: string
| SET CONFIG tikv stats = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats = 1
< ok 0 0
> exec ^
: string
| KILL stats_auto_recalc()
< error
| 1064 42000 …line 1 column 22 near "stats_auto_recalc()" 
> exec ^
: string
| KILL stats_auto_recalc (1)
< error
| 1064 42000 …line 1 column 22 near "stats_auto_recalc (1)" 
> exec ^
: string
| SET CONFIG tikv stats_auto_recalc = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_auto_recalc = 1
< ok 0 0
> exec ^
: string
| KILL stats_buckets()
< error
| 1064 42000 …line 1 column 18 near "stats_buckets()" 
> exec ^
: string
| KILL stats_buckets (1)
< error
| 1064 42000 …line 1 column 18 near "stats_buckets (1)" 
> exec ^
: string
| SET CONFIG tikv stats_buckets = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_buckets = 1
< ok 0 0
> exec ^
: string
| KILL stats_col_choice()
< error
| 1064 42000 …line 1 column 21 near "stats_col_choice()" 
> exec ^
: string
| KILL stats_col_choice (1)
< error
| 1064 42000 …line 1 column 21 near "stats_col_choice (1)" 
> exec ^
: string
| SET CONFIG tikv stats_col_choice = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_col_choice = 1
< ok 0 0
> exec ^
: string
| KILL stats_col_list()
< error
| 1064 42000 …line 1 column 19 near "stats_col_list()" 
> exec ^
: string
| KILL stats_col_list (1)
< error
| 1064 42000 …line 1 column 19 near "stats_col_list (1)" 
> exec ^
: string
| SET CONFIG tikv stats_col_list = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_col_list = 1
< ok 0 0
> exec ^
: string
| KILL stats_extended()
< error
| 1064 42000 …line 1 column 19 near "stats_extended()" 
> exec ^
: string
| KILL stats_extended (1)
< error
| 1064 42000 …line 1 column 19 near "stats_extended (1)" 
> exec ^
: string
| SET CONFIG tikv stats_extended = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_extended = 1
< ok 0 0
> exec ^
: string
| KILL stats_healthy()
< error
| 1064 42000 …line 1 column 18 near "stats_healthy()" 
> exec ^
: string
| KILL stats_healthy (1)
< error
| 1064 42000 …line 1 column 18 near "stats_healthy (1)" 
> exec ^
: string
| SET CONFIG tikv stats_healthy = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_healthy = 1
< ok 0 0
> exec ^
: string
| KILL stats_histograms()
< error
| 1064 42000 …line 1 column 21 near "stats_histograms()" 
> exec ^
: string
| KILL stats_histograms (1)
< error
| 1064 42000 …line 1 column 21 near "stats_histograms (1)" 
> exec ^
: string
| SET CONFIG tikv stats_histograms = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_histograms = 1
< ok 0 0
> exec ^
: string
| KILL stats_locked()
< error
| 1064 42000 …line 1 column 17 near "stats_locked()" 
> exec ^
: string
| KILL stats_locked (1)
< error
| 1064 42000 …line 1 column 17 near "stats_locked (1)" 
> exec ^
: string
| SET CONFIG tikv stats_locked = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_locked = 1
< ok 0 0
> exec ^
: string
| KILL stats_meta()
< error
| 1064 42000 …line 1 column 15 near "stats_meta()" 
> exec ^
: string
| KILL stats_meta (1)
< error
| 1064 42000 …line 1 column 15 near "stats_meta (1)" 
> exec ^
: string
| SET CONFIG tikv stats_meta = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_meta = 1
< ok 0 0
> exec ^
: string
| KILL stats_options()
< error
| 1064 42000 …line 1 column 18 near "stats_options()" 
> exec ^
: string
| KILL stats_options (1)
< error
| 1064 42000 …line 1 column 18 near "stats_options (1)" 
> exec ^
: string
| SET CONFIG tikv stats_options = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_options = 1
< ok 0 0
> exec ^
: string
| KILL stats_persistent()
< error
| 1064 42000 …line 1 column 21 near "stats_persistent()" 
> exec ^
: string
| KILL stats_persistent (1)
< error
| 1064 42000 …line 1 column 21 near "stats_persistent (1)" 
> exec ^
: string
| SET CONFIG tikv stats_persistent = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_persistent = 1
< ok 0 0
> exec ^
: string
| KILL stats_sample_pages()
< error
| 1064 42000 …line 1 column 23 near "stats_sample_pages()" 
> exec ^
: string
| KILL stats_sample_pages (1)
< error
| 1064 42000 …line 1 column 23 near "stats_sample_pages (1)" 
> exec ^
: string
| SET CONFIG tikv stats_sample_pages = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_sample_pages = 1
< ok 0 0
> exec ^
: string
| KILL stats_sample_rate()
< error
| 1064 42000 …line 1 column 22 near "stats_sample_rate()" 
> exec ^
: string
| KILL stats_sample_rate (1)
< error
| 1064 42000 …line 1 column 22 near "stats_sample_rate (1)" 
> exec ^
: string
| SET CONFIG tikv stats_sample_rate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_sample_rate = 1
< ok 0 0
> exec ^
: string
| KILL stats_topn()
< error
| 1064 42000 …line 1 column 15 near "stats_topn()" 
> exec ^
: string
| KILL stats_topn (1)
< error
| 1064 42000 …line 1 column 15 near "stats_topn (1)" 
> exec ^
: string
| SET CONFIG tikv stats_topn = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stats_topn = 1
< ok 0 0
> exec ^
: string
| KILL status()
< error
| 1064 42000 …line 1 column 11 near "status()" 
> exec ^
: string
| KILL status (1)
< error
| 1064 42000 …line 1 column 11 near "status (1)" 
> exec ^
: string
| SET CONFIG tikv status = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-status = 1
< ok 0 0
> exec ^
: string
| KILL std()
< error
| 1064 42000 …line 1 column 8 near "std()" 
> exec ^
: string
| KILL std (1)
< error
| 1064 42000 …line 1 column 8 near "std (1)" 
> exec ^
: string
| SET CONFIG tikv std = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-std = 1
< ok 0 0
> exec ^
: string
| KILL stddev()
< error
| 1064 42000 …line 1 column 11 near "stddev()" 
> exec ^
: string
| KILL stddev (1)
< error
| 1064 42000 …line 1 column 11 near "stddev (1)" 
> exec ^
: string
| SET CONFIG tikv stddev = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stddev = 1
< ok 0 0
> exec ^
: string
| KILL stddev_pop()
< error
| 1064 42000 …line 1 column 15 near "stddev_pop()" 
> exec ^
: string
| KILL stddev_pop (1)
< error
| 1064 42000 …line 1 column 15 near "stddev_pop (1)" 
> exec ^
: string
| SET CONFIG tikv stddev_pop = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stddev_pop = 1
< ok 0 0
> exec ^
: string
| KILL stddev_samp()
< error
| 1064 42000 …line 1 column 16 near "stddev_samp()" 
> exec ^
: string
| KILL stddev_samp (1)
< error
| 1064 42000 …line 1 column 16 near "stddev_samp (1)" 
> exec ^
: string
| SET CONFIG tikv stddev_samp = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stddev_samp = 1
< ok 0 0
> exec ^
: string
| KILL stop()
< error
| 1064 42000 …line 1 column 9 near "stop()" 
> exec ^
: string
| KILL stop (1)
< error
| 1064 42000 …line 1 column 9 near "stop (1)" 
> exec ^
: string
| SET CONFIG tikv stop = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stop = 1
< ok 0 0
> exec ^
: string
| KILL storage()
< error
| 1064 42000 …line 1 column 12 near "storage()" 
> exec ^
: string
| KILL storage (1)
< error
| 1064 42000 …line 1 column 12 near "storage (1)" 
> exec ^
: string
| SET CONFIG tikv storage = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-storage = 1
< ok 0 0
> exec ^
: string
| KILL stored()
< error
| 1064 42000 …line 1 column 11 near "stored()" 
> exec ^
: string
| KILL stored (1)
< error
| 1064 42000 …line 1 column 11 near "stored (1)" 
> exec ^
: string
| SET CONFIG tikv stored = 1
< error
| 1064 42000 …line 1 column 22 near "stored = 1" 
> exec ^
: string
| SET CONFIG tikv a-stored = 1
< error
| 1064 42000 …line 1 column 24 near "stored = 1" 
> exec ^
: string
| KILL straight_join()
< error
| 1064 42000 …line 1 column 18 near "straight_join()" 
> exec ^
: string
| KILL straight_join (1)
< error
| 1064 42000 …line 1 column 18 near "straight_join (1)" 
> exec ^
: string
| SET CONFIG tikv straight_join = 1
< error
| 1064 42000 …line 1 column 29 near "straight_join = 1" 
> exec ^
: string
| SET CONFIG tikv a-straight_join = 1
< error
| 1064 42000 …line 1 column 31 near "straight_join = 1" 
> exec ^
: string
| KILL stream()
< ok 0 0
> exec ^
: string
| KILL stream (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv stream = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-stream = 1
< ok 0 0
> exec ^
: string
| KILL strict()
< error
| 1064 42000 …line 1 column 11 near "strict()" 
> exec ^
: string
| KILL strict (1)
< error
| 1064 42000 …line 1 column 11 near "strict (1)" 
> exec ^
: string
| SET CONFIG tikv strict = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-strict = 1
< ok 0 0
> exec ^
: string
| KILL strict_format()
< error
| 1064 42000 …line 1 column 18 near "strict_format()" 
> exec ^
: string
| KILL strict_format (1)
< error
| 1064 42000 …line 1 column 18 near "strict_format (1)" 
> exec ^
: string
| SET CONFIG tikv strict_format = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-strict_format = 1
< ok 0 0
> exec ^
: string
| KILL strict_load()
< ok 0 0
> exec ^
: string
| KILL strict_load (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv strict_load = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-strict_load = 1
< ok 0 0
> exec ^
: string
| KILL string()
< ok 0 0
> exec ^
: string
| KILL string (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv string = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-string = 1
< ok 0 0
> exec ^
: string
| KILL strong()
< error
| 1064 42000 …line 1 column 11 near "strong()" 
> exec ^
: string
| KILL strong (1)
< error
| 1064 42000 …line 1 column 11 near "strong (1)" 
> exec ^
: string
| SET CONFIG tikv strong = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-strong = 1
< ok 0 0
> exec ^
: string
| KILL subclass_origin()
< ok 0 0
> exec ^
: string
| KILL subclass_origin (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv subclass_origin = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-subclass_origin = 1
< ok 0 0
> exec ^
: string
| KILL subdate()
< error
| 1064 42000 …line 1 column 12 near "subdate()" 
> exec ^
: string
| KILL subdate (1)
< error
| 1064 42000 …line 1 column 12 near "subdate (1)" 
> exec ^
: string
| SET CONFIG tikv subdate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-subdate = 1
< ok 0 0
> exec ^
: string
| KILL subject()
< error
| 1064 42000 …line 1 column 12 near "subject()" 
> exec ^
: string
| KILL subject (1)
< error
| 1064 42000 …line 1 column 12 near "subject (1)" 
> exec ^
: string
| SET CONFIG tikv subject = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-subject = 1
< ok 0 0
> exec ^
: string
| KILL subpartition()
< error
| 1064 42000 …line 1 column 17 near "subpartition()" 
> exec ^
: string
| KILL subpartition (1)
< error
| 1064 42000 …line 1 column 17 near "subpartition (1)" 
> exec ^
: string
| SET CONFIG tikv subpartition = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-subpartition = 1
< ok 0 0
> exec ^
: string
| KILL subpartitions()
< error
| 1064 42000 …line 1 column 18 near "subpartitions()" 
> exec ^
: string
| KILL subpartitions (1)
< error
| 1064 42000 …line 1 column 18 near "subpartitions (1)" 
> exec ^
: string
| SET CONFIG tikv subpartitions = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-subpartitions = 1
< ok 0 0
> exec ^
: string
| KILL substr()
< error
| 1064 42000 …line 1 column 11 near "substr()" 
> exec ^
: string
| KILL substr (1)
< error
| 1064 42000 …line 1 column 11 near "substr (1)" 
> exec ^
: string
| SET CONFIG tikv substr = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-substr = 1
< ok 0 0
> exec ^
: string
| KILL substring()
< error
| 1064 42000 …line 1 column 14 near "substring()" 
> exec ^
: string
| KILL substring (1)
< error
| 1064 42000 …line 1 column 14 near "substring (1)" 
> exec ^
: string
| SET CONFIG tikv substring = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-substring = 1
< ok 0 0
> exec ^
: string
| KILL sum()
< error
| 1064 42000 …line 1 column 8 near "sum()" 
> exec ^
: string
| KILL sum (1)
< error
| 1064 42000 …line 1 column 8 near "sum (1)" 
> exec ^
: string
| SET CONFIG tikv sum = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sum = 1
< ok 0 0
> exec ^
: string
| KILL super()
< error
| 1064 42000 …line 1 column 10 near "super()" 
> exec ^
: string
| KILL super (1)
< error
| 1064 42000 …line 1 column 10 near "super (1)" 
> exec ^
: string
| SET CONFIG tikv super = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-super = 1
< ok 0 0
> exec ^
: string
| KILL survival_preferences()
< error
| 1064 42000 …line 1 column 25 near "survival_preferences()" 
> exec ^
: string
| KILL survival_preferences (1)
< error
| 1064 42000 …line 1 column 25 near "survival_preferences (1)" 
> exec ^
: string
| SET CONFIG tikv survival_preferences = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-survival_preferences = 1
< ok 0 0
> exec ^
: string
| KILL suspend()
< ok 0 0
> exec ^
: string
| KILL suspend (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv suspend = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-suspend = 1
< ok 0 0
> exec ^
: string
| KILL swaps()
< error
| 1064 42000 …line 1 column 10 near "swaps()" 
> exec ^
: string
| KILL swaps (1)
< error
| 1064 42000 …line 1 column 10 near "swaps (1)" 
> exec ^
: string
| SET CONFIG tikv swaps = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-swaps = 1
< ok 0 0
> exec ^
: string
| KILL switch_group()
< error
| 1064 42000 …line 1 column 17 near "switch_group()" 
> exec ^
: string
| KILL switch_group (1)
< error
| 1064 42000 …line 1 column 17 near "switch_group (1)" 
> exec ^
: string
| SET CONFIG tikv switch_group = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-switch_group = 1
< ok 0 0
> exec ^
: string
| KILL switches()
< error
| 1064 42000 …line 1 column 13 near "switches()" 
> exec ^
: string
| KILL switches (1)
< error
| 1064 42000 …line 1 column 13 near "switches (1)" 
> exec ^
: string
| SET CONFIG tikv switches = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-switches = 1
< ok 0 0
> exec ^
: string
| KILL sysdate()
< error
| 1064 42000 …line 1 column 12 near "sysdate()" 
> exec ^
: string
| KILL sysdate (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv sysdate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-sysdate = 1
< ok 0 0
> exec ^
: string
| KILL system()
< error
| 1064 42000 …line 1 column 11 near "system()" 
> exec ^
: string
| KILL system (1)
< error
| 1064 42000 …line 1 column 11 near "system (1)" 
> exec ^
: string
| SET CONFIG tikv system = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-system = 1
< ok 0 0
> exec ^
: string
| KILL system_time()
< error
| 1064 42000 …line 1 column 16 near "system_time()" 
> exec ^
: string
| KILL system_time (1)
< error
| 1064 42000 …line 1 column 16 near "system_time (1)" 
> exec ^
: string
| SET CONFIG tikv system_time = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-system_time = 1
< ok 0 0
> exec ^
: string
| KILL system_user()
< ok 0 0
> exec ^
: string
| KILL system_user (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv system_user = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-system_user = 1
< ok 0 0
> exec ^
: string
| KILL table()
< error
| 1064 42000 …line 1 column 10 near "table()" 
> exec ^
: string
| KILL table (1)
< error
| 1064 42000 …line 1 column 10 near "table (1)" 
> exec ^
: string
| SET CONFIG tikv table = 1
< error
| 1064 42000 …line 1 column 21 near "table = 1" 
> exec ^
: string
| SET CONFIG tikv a-table = 1
< error
| 1064 42000 …line 1 column 23 near "table = 1" 
> exec ^
: string
| KILL table_checksum()
< error
| 1064 42000 …line 1 column 19 near "table_checksum()" 
> exec ^
: string
| KILL table_checksum (1)
< error
| 1064 42000 …line 1 column 19 near "table_checksum (1)" 
> exec ^
: string
| SET CONFIG tikv table_checksum = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-table_checksum = 1
< ok 0 0
> exec ^
: string
| KILL table_name()
< ok 0 0
> exec ^
: string
| KILL table_name (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv table_name = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-table_name = 1
< ok 0 0
> exec ^
: string
| KILL tables()
< error
| 1064 42000 …line 1 column 11 near "tables()" 
> exec ^
: string
| KILL tables (1)
< error
| 1064 42000 …line 1 column 11 near "tables (1)" 
> exec ^
: string
| SET CONFIG tikv tables = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tables = 1
< ok 0 0
> exec ^
: string
| KILL tablesample()
< error
| 1064 42000 …line 1 column 16 near "tablesample()" 
> exec ^
: string
| KILL tablesample (1)
< error
| 1064 42000 …line 1 column 16 near "tablesample (1)" 
> exec ^
: string
| SET CONFIG tikv tablesample = 1
< error
| 1064 42000 …line 1 column 27 near "tablesample = 1" 
> exec ^
: string
| SET CONFIG tikv a-tablesample = 1
< error
| 1064 42000 …line 1 column 29 near "tablesample = 1" 
> exec ^
: string
| KILL tablespace()
< error
| 1064 42000 …line 1 column 15 near "tablespace()" 
> exec ^
: string
| KILL tablespace (1)
< error
| 1064 42000 …line 1 column 15 near "tablespace (1)" 
> exec ^
: string
| SET CONFIG tikv tablespace = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tablespace = 1
< ok 0 0
> exec ^
: string
| KILL target()
< error
| 1064 42000 …line 1 column 11 near "target()" 
> exec ^
: string
| KILL target (1)
< error
| 1064 42000 …line 1 column 11 near "target (1)" 
> exec ^
: string
| SET CONFIG tikv target = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-target = 1
< ok 0 0
> exec ^
: string
| KILL task_types()
< error
| 1064 42000 …line 1 column 15 near "task_types()" 
> exec ^
: string
| KILL task_types (1)
< error
| 1064 42000 …line 1 column 15 near "task_types (1)" 
> exec ^
: string
| SET CONFIG tikv task_types = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-task_types = 1
< ok 0 0
> exec ^
: string
| KILL temporary()
< error
| 1064 42000 …line 1 column 14 near "temporary()" 
> exec ^
: string
| KILL temporary (1)
< error
| 1064 42000 …line 1 column 14 near "temporary (1)" 
> exec ^
: string
| SET CONFIG tikv temporary = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-temporary = 1
< ok 0 0
> exec ^
: string
| KILL temptable()
< error
| 1064 42000 …line 1 column 14 near "temptable()" 
> exec ^
: string
| KILL temptable (1)
< error
| 1064 42000 …line 1 column 14 near "temptable (1)" 
> exec ^
: string
| SET CONFIG tikv temptable = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-temptable = 1
< ok 0 0
> exec ^
: string
| KILL terminated()
< error
| 1064 42000 …line 1 column 15 near "terminated()" 
> exec ^
: string
| KILL terminated (1)
< error
| 1064 42000 …line 1 column 15 near "terminated (1)" 
> exec ^
: string
| SET CONFIG tikv terminated = 1
< error
| 1064 42000 …line 1 column 26 near "terminated = 1" 
> exec ^
: string
| SET CONFIG tikv a-terminated = 1
< error
| 1064 42000 …line 1 column 28 near "terminated = 1" 
> exec ^
: string
| KILL text()
< error
| 1064 42000 …line 1 column 9 near "text()" 
> exec ^
: string
| KILL text (1)
< error
| 1064 42000 …line 1 column 9 near "text (1)" 
> exec ^
: string
| SET CONFIG tikv text = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-text = 1
< ok 0 0
> exec ^
: string
| KILL than()
< error
| 1064 42000 …line 1 column 9 near "than()" 
> exec ^
: string
| KILL than (1)
< error
| 1064 42000 …line 1 column 9 near "than (1)" 
> exec ^
: string
| SET CONFIG tikv than = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-than = 1
< ok 0 0
> exec ^
: string
| KILL then()
< error
| 1064 42000 …line 1 column 9 near "then()" 
> exec ^
: string
| KILL then (1)
< error
| 1064 42000 …line 1 column 9 near "then (1)" 
> exec ^
: string
| SET CONFIG tikv then = 1
< error
| 1064 42000 …line 1 column 20 near "then = 1" 
> exec ^
: string
| SET CONFIG tikv a-then = 1
< error
| 1064 42000 …line 1 column 22 near "then = 1" 
> exec ^
: string
| KILL thread_priority()
< ok 0 0
> exec ^
: string
| KILL thread_priority (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv thread_priority = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-thread_priority = 1
< ok 0 0
> exec ^
: string
| KILL tidb()
< error
| 1064 42000 …line 1 column 11 near ")" 
> exec ^
: string
| KILL tidb (1)
< error
| 1064 42000 …line 1 column 12 near "1)" 
> exec ^
: string
| SET CONFIG tikv tidb = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tidb = 1
< ok 0 0
> exec ^
: string
| KILL tidb_current_tso()
< error
| 1064 42000 …line 1 column 21 near "tidb_current_tso()" 
> exec ^
: string
| KILL tidb_current_tso (1)
< error
| 1064 42000 …line 1 column 21 near "tidb_current_tso (1)" 
> exec ^
: string
| KILL tidb_json()
< error
| 1064 42000 …line 1 column 14 near "tidb_json()" 
> exec ^
: string
| KILL tidb_json (1)
< error
| 1064 42000 …line 1 column 14 near "tidb_json (1)" 
> exec ^
: string
| SET CONFIG tikv tidb_json = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tidb_json = 1
< ok 0 0
> exec ^
: string
| KILL ties()
< ok 0 0
> exec ^
: string
| KILL ties (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv ties = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ties = 1
< ok 0 0
> exec ^
: string
| KILL tiflash()
< error
| 1064 42000 …line 1 column 12 near "tiflash()" 
> exec ^
: string
| KILL tiflash (1)
< error
| 1064 42000 …line 1 column 12 near "tiflash (1)" 
> exec ^
: string
| SET CONFIG tikv tiflash = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tiflash = 1
< ok 0 0
> exec ^
: string
| KILL tikv_importer()
< error
| 1064 42000 …line 1 column 18 near "tikv_importer()" 
> exec ^
: string
| KILL tikv_importer (1)
< error
| 1064 42000 …line 1 column 18 near "tikv_importer (1)" 
> exec ^
: string
| SET CONFIG tikv tikv_importer = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tikv_importer = 1
< ok 0 0
> exec ^
: string
| KILL time()
< error
| 1064 42000 …line 1 column 9 near "time()" 
> exec ^
: string
| KILL time (1)
< error
| 1064 42000 …line 1 column 9 near "time (1)" 
> exec ^
: string
| SET CONFIG tikv time = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-time = 1
< ok 0 0
> exec ^
: string
| KILL timeout()
< error
| 1064 42000 …line 1 column 12 near "timeout()" 
> exec ^
: string
| KILL timeout (1)
< error
| 1064 42000 …line 1 column 12 near "timeout (1)" 
> exec ^
: string
| SET CONFIG tikv timeout = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-timeout = 1
< ok 0 0
> exec ^
: string
| KILL timestamp()
< error
| 1064 42000 …line 1 column 14 near "timestamp()" 
> exec ^
: string
| KILL timestamp (1)
< error
| 1064 42000 …line 1 column 14 near "timestamp (1)" 
> exec ^
: string
| SET CONFIG tikv timestamp = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-timestamp = 1
< ok 0 0
> exec ^
: string
| KILL timestampadd()
< error
| 1064 42000 …line 1 column 17 near "timestampadd()" 
> exec ^
: string
| KILL timestampadd (1)
< error
| 1064 42000 …line 1 column 17 near "timestampadd (1)" 
> exec ^
: string
| SET CONFIG tikv timestampadd = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-timestampadd = 1
< ok 0 0
> exec ^
: string
| KILL timestampdiff()
< error
| 1064 42000 …line 1 column 18 near "timestampdiff()" 
> exec ^
: string
| KILL timestampdiff (1)
< error
| 1064 42000 …line 1 column 18 near "timestampdiff (1)" 
> exec ^
: string
| SET CONFIG tikv timestampdiff = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-timestampdiff = 1
< ok 0 0
> exec ^
: string
| KILL tinyblob()
< error
| 1064 42000 …line 1 column 13 near "tinyblob()" 
> exec ^
: string
| KILL tinyblob (1)
< error
| 1064 42000 …line 1 column 13 near "tinyblob (1)" 
> exec ^
: string
| SET CONFIG tikv tinyblob = 1
< error
| 1064 42000 …line 1 column 24 near "tinyblob = 1" 
> exec ^
: string
| SET CONFIG tikv a-tinyblob = 1
< error
| 1064 42000 …line 1 column 26 near "tinyblob = 1" 
> exec ^
: string
| KILL tinyint()
< error
| 1064 42000 …line 1 column 12 near "tinyint()" 
> exec ^
: string
| KILL tinyint (1)
< error
| 1064 42000 …line 1 column 12 near "tinyint (1)" 
> exec ^
: string
| SET CONFIG tikv tinyint = 1
< error
| 1064 42000 …line 1 column 23 near "tinyint = 1" 
> exec ^
: string
| SET CONFIG tikv a-tinyint = 1
< error
| 1064 42000 …line 1 column 25 near "tinyint = 1" 
> exec ^
: string
| KILL tinytext()
< error
| 1064 42000 …line 1 column 13 near "tinytext()" 
> exec ^
: string
| KILL tinytext (1)
< error
| 1064 42000 …line 1 column 13 near "tinytext (1)" 
> exec ^
: string
| SET CONFIG tikv tinytext = 1
< error
| 1064 42000 …line 1 column 24 near "tinytext = 1" 
> exec ^
: string
| SET CONFIG tikv a-tinytext = 1
< error
| 1064 42000 …line 1 column 26 near "tinytext = 1" 
> exec ^
: string
| KILL tls()
< error
| 1064 42000 …line 1 column 8 near "tls()" 
> exec ^
: string
| KILL tls (1)
< error
| 1064 42000 …line 1 column 8 near "tls (1)" 
> exec ^
: string
| SET CONFIG tikv tls = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tls = 1
< ok 0 0
> exec ^
: string
| KILL to()
< error
| 1064 42000 …line 1 column 7 near "to()" 
> exec ^
: string
| KILL to (1)
< error
| 1064 42000 …line 1 column 7 near "to (1)" 
> exec ^
: string
| SET CONFIG tikv to = 1
< error
| 1064 42000 …line 1 column 18 near "to = 1" 
> exec ^
: string
| SET CONFIG tikv a-to = 1
< error
| 1064 42000 …line 1 column 20 near "to = 1" 
> exec ^
: string
| KILL token_issuer()
< error
| 1064 42000 …line 1 column 17 near "token_issuer()" 
> exec ^
: string
| KILL token_issuer (1)
< error
| 1064 42000 …line 1 column 17 near "token_issuer (1)" 
> exec ^
: string
| SET CONFIG tikv token_issuer = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-token_issuer = 1
< ok 0 0
> exec ^
: string
| KILL tokudb_default()
< error
| 1064 42000 …line 1 column 19 near "tokudb_default()" 
> exec ^
: string
| KILL tokudb_default (1)
< error
| 1064 42000 …line 1 column 19 near "tokudb_default (1)" 
> exec ^
: string
| SET CONFIG tikv tokudb_default = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tokudb_default = 1
< ok 0 0
> exec ^
: string
| KILL tokudb_fast()
< error
| 1064 42000 …line 1 column 16 near "tokudb_fast()" 
> exec ^
: string
| KILL tokudb_fast (1)
< error
| 1064 42000 …line 1 column 16 near "tokudb_fast (1)" 
> exec ^
: string
| SET CONFIG tikv tokudb_fast = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tokudb_fast = 1
< ok 0 0
> exec ^
: string
| KILL tokudb_lzma()
< error
| 1064 42000 …line 1 column 16 near "tokudb_lzma()" 
> exec ^
: string
| KILL tokudb_lzma (1)
< error
| 1064 42000 …line 1 column 16 near "tokudb_lzma (1)" 
> exec ^
: string
| SET CONFIG tikv tokudb_lzma = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tokudb_lzma = 1
< ok 0 0
> exec ^
: string
| KILL tokudb_quicklz()
< error
| 1064 42000 …line 1 column 19 near "tokudb_quicklz()" 
> exec ^
: string
| KILL tokudb_quicklz (1)
< error
| 1064 42000 …line 1 column 19 near "tokudb_quicklz (1)" 
> exec ^
: string
| SET CONFIG tikv tokudb_quicklz = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tokudb_quicklz = 1
< ok 0 0
> exec ^
: string
| KILL tokudb_small()
< error
| 1064 42000 …line 1 column 17 near "tokudb_small()" 
> exec ^
: string
| KILL tokudb_small (1)
< error
| 1064 42000 …line 1 column 17 near "tokudb_small (1)" 
> exec ^
: string
| SET CONFIG tikv tokudb_small = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tokudb_small = 1
< ok 0 0
> exec ^
: string
| KILL tokudb_snappy()
< error
| 1064 42000 …line 1 column 18 near "tokudb_snappy()" 
> exec ^
: string
| KILL tokudb_snappy (1)
< error
| 1064 42000 …line 1 column 18 near "tokudb_snappy (1)" 
> exec ^
: string
| SET CONFIG tikv tokudb_snappy = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tokudb_snappy = 1
< ok 0 0
> exec ^
: string
| KILL tokudb_uncompressed()
< error
| 1064 42000 …line 1 column 24 near "tokudb_uncompressed()" 
> exec ^
: string
| KILL tokudb_uncompressed (1)
< error
| 1064 42000 …line 1 column 24 near "tokudb_uncompressed (1)" 
> exec ^
: string
| SET CONFIG tikv tokudb_uncompressed = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tokudb_uncompressed = 1
< ok 0 0
> exec ^
: string
| KILL tokudb_zlib()
< error
| 1064 42000 …line 1 column 16 near "tokudb_zlib()" 
> exec ^
: string
| KILL tokudb_zlib (1)
< error
| 1064 42000 …line 1 column 16 near "tokudb_zlib (1)" 
> exec ^
: string
| SET CONFIG tikv tokudb_zlib = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tokudb_zlib = 1
< ok 0 0
> exec ^
: string
| KILL tokudb_zstd()
< error
| 1064 42000 …line 1 column 16 near "tokudb_zstd()" 
> exec ^
: string
| KILL tokudb_zstd (1)
< error
| 1064 42000 …line 1 column 16 near "tokudb_zstd (1)" 
> exec ^
: string
| SET CONFIG tikv tokudb_zstd = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tokudb_zstd = 1
< ok 0 0
> exec ^
: string
| KILL top()
< error
| 1064 42000 …line 1 column 8 near "top()" 
> exec ^
: string
| KILL top (1)
< error
| 1064 42000 …line 1 column 8 near "top (1)" 
> exec ^
: string
| SET CONFIG tikv top = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-top = 1
< ok 0 0
> exec ^
: string
| KILL topn()
< error
| 1064 42000 …line 1 column 9 near "topn()" 
> exec ^
: string
| KILL topn (1)
< error
| 1064 42000 …line 1 column 9 near "topn (1)" 
> exec ^
: string
| SET CONFIG tikv topn = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-topn = 1
< ok 0 0
> exec ^
: string
| KILL tpcc()
< error
| 1064 42000 …line 1 column 9 near "tpcc()" 
> exec ^
: string
| KILL tpcc (1)
< error
| 1064 42000 …line 1 column 9 near "tpcc (1)" 
> exec ^
: string
| SET CONFIG tikv tpcc = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tpcc = 1
< ok 0 0
> exec ^
: string
| KILL tpch_10()
< error
| 1064 42000 …line 1 column 12 near "tpch_10()" 
> exec ^
: string
| KILL tpch_10 (1)
< error
| 1064 42000 …line 1 column 12 near "tpch_10 (1)" 
> exec ^
: string
| SET CONFIG tikv tpch_10 = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tpch_10 = 1
< ok 0 0
> exec ^
: string
| KILL trace()
< error
| 1064 42000 …line 1 column 10 near "trace()" 
> exec ^
: string
| KILL trace (1)
< error
| 1064 42000 …line 1 column 10 near "trace (1)" 
> exec ^
: string
| SET CONFIG tikv trace = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-trace = 1
< ok 0 0
> exec ^
: string
| KILL traditional()
< error
| 1064 42000 …line 1 column 16 near "traditional()" 
> exec ^
: string
| KILL traditional (1)
< error
| 1064 42000 …line 1 column 16 near "traditional (1)" 
> exec ^
: string
| SET CONFIG tikv traditional = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-traditional = 1
< ok 0 0
> exec ^
: string
| KILL trailing()
< error
| 1064 42000 …line 1 column 13 near "trailing()" 
> exec ^
: string
| KILL trailing (1)
< error
| 1064 42000 …line 1 column 13 near "trailing (1)" 
> exec ^
: string
| SET CONFIG tikv trailing = 1
< error
| 1064 42000 …line 1 column 24 near "trailing = 1" 
> exec ^
: string
| SET CONFIG tikv a-trailing = 1
< error
| 1064 42000 …line 1 column 26 near "trailing = 1" 
> exec ^
: string
| KILL transaction()
< error
| 1064 42000 …line 1 column 16 near "transaction()" 
> exec ^
: string
| KILL transaction (1)
< error
| 1064 42000 …line 1 column 16 near "transaction (1)" 
> exec ^
: string
| SET CONFIG tikv transaction = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-transaction = 1
< ok 0 0
> exec ^
: string
| KILL translate()
< error
| 1064 42000 …line 1 column 14 near "translate()" 
> exec ^
: string
| KILL translate (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv translate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-translate = 1
< ok 0 0
> exec ^
: string
| KILL trigger()
< error
| 1064 42000 …line 1 column 12 near "trigger()" 
> exec ^
: string
| KILL trigger (1)
< error
| 1064 42000 …line 1 column 12 near "trigger (1)" 
> exec ^
: string
| SET CONFIG tikv trigger = 1
< error
| 1064 42000 …line 1 column 23 near "trigger = 1" 
> exec ^
: string
| SET CONFIG tikv a-trigger = 1
< error
| 1064 42000 …line 1 column 25 near "trigger = 1" 
> exec ^
: string
| KILL triggers()
< error
| 1064 42000 …line 1 column 13 near "triggers()" 
> exec ^
: string
| KILL triggers (1)
< error
| 1064 42000 …line 1 column 13 near "triggers (1)" 
> exec ^
: string
| SET CONFIG tikv triggers = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-triggers = 1
< ok 0 0
> exec ^
: string
| KILL trim()
< error
| 1064 42000 …line 1 column 9 near "trim()" 
> exec ^
: string
| KILL trim (1)
< error
| 1064 42000 …line 1 column 9 near "trim (1)" 
> exec ^
: string
| SET CONFIG tikv trim = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-trim = 1
< ok 0 0
> exec ^
: string
| KILL true()
< error
| 1064 42000 …line 1 column 9 near "true()" 
> exec ^
: string
| KILL true (1)
< error
| 1064 42000 …line 1 column 9 near "true (1)" 
> exec ^
: string
| SET CONFIG tikv true = 1
< error
| 1064 42000 …line 1 column 20 near "true = 1" 
> exec ^
: string
| SET CONFIG tikv a-true = 1
< error
| 1064 42000 …line 1 column 22 near "true = 1" 
> exec ^
: string
| KILL true_card_cost()
< error
| 1064 42000 …line 1 column 19 near "true_card_cost()" 
> exec ^
: string
| KILL true_card_cost (1)
< error
| 1064 42000 …line 1 column 19 near "true_card_cost (1)" 
> exec ^
: string
| SET CONFIG tikv true_card_cost = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-true_card_cost = 1
< ok 0 0
> exec ^
: string
| KILL truncate()
< error
| 1064 42000 …line 1 column 13 near "truncate()" 
> exec ^
: string
| KILL truncate (1)
< error
| 1064 42000 …line 1 column 13 near "truncate (1)" 
> exec ^
: string
| SET CONFIG tikv truncate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-truncate = 1
< ok 0 0
> exec ^
: string
| KILL tso()
< error
| 1064 42000 …line 1 column 8 near "tso()" 
> exec ^
: string
| KILL tso (1)
< error
| 1064 42000 …line 1 column 8 near "tso (1)" 
> exec ^
: string
| SET CONFIG tikv tso = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-tso = 1
< ok 0 0
> exec ^
: string
| KILL ttl()
< error
| 1064 42000 …line 1 column 8 near "ttl()" 
> exec ^
: string
| KILL ttl (1)
< error
| 1064 42000 …line 1 column 8 near "ttl (1)" 
> exec ^
: string
| SET CONFIG tikv ttl = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ttl = 1
< ok 0 0
> exec ^
: string
| KILL ttl_enable()
< error
| 1064 42000 …line 1 column 15 near "ttl_enable()" 
> exec ^
: string
| KILL ttl_enable (1)
< error
| 1064 42000 …line 1 column 15 near "ttl_enable (1)" 
> exec ^
: string
| SET CONFIG tikv ttl_enable = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ttl_enable = 1
< ok 0 0
> exec ^
: string
| KILL ttl_job_interval()
< error
| 1064 42000 …line 1 column 21 near "ttl_job_interval()" 
> exec ^
: string
| KILL ttl_job_interval (1)
< error
| 1064 42000 …line 1 column 21 near "ttl_job_interval (1)" 
> exec ^
: string
| SET CONFIG tikv ttl_job_interval = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-ttl_job_interval = 1
< ok 0 0
> exec ^
: string
| KILL type()
< error
| 1064 42000 …line 1 column 9 near "type()" 
> exec ^
: string
| KILL type (1)
< error
| 1064 42000 …line 1 column 9 near "type (1)" 
> exec ^
: string
| SET CONFIG tikv type = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-type = 1
< ok 0 0
> exec ^
: string
| KILL types()
< ok 0 0
> exec ^
: string
| KILL types (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv types = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-types = 1
< ok 0 0
> exec ^
: string
| KILL unbounded()
< error
| 1064 42000 …line 1 column 14 near "unbounded()" 
> exec ^
: string
| KILL unbounded (1)
< error
| 1064 42000 …line 1 column 14 near "unbounded (1)" 
> exec ^
: string
| SET CONFIG tikv unbounded = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-unbounded = 1
< ok 0 0
> exec ^
: string
| KILL uncommitted()
< error
| 1064 42000 …line 1 column 16 near "uncommitted()" 
> exec ^
: string
| KILL uncommitted (1)
< error
| 1064 42000 …line 1 column 16 near "uncommitted (1)" 
> exec ^
: string
| SET CONFIG tikv uncommitted = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-uncommitted = 1
< ok 0 0
> exec ^
: string
| KILL undefined()
< error
| 1064 42000 …line 1 column 14 near "undefined()" 
> exec ^
: string
| KILL undefined (1)
< error
| 1064 42000 …line 1 column 14 near "undefined (1)" 
> exec ^
: string
| SET CONFIG tikv undefined = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-undefined = 1
< ok 0 0
> exec ^
: string
| KILL undo()
< ok 0 0
> exec ^
: string
| KILL undo (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv undo = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-undo = 1
< ok 0 0
> exec ^
: string
| KILL undo_buffer_size()
< ok 0 0
> exec ^
: string
| KILL undo_buffer_size (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv undo_buffer_size = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-undo_buffer_size = 1
< ok 0 0
> exec ^
: string
| KILL undofile()
< ok 0 0
> exec ^
: string
| KILL undofile (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv undofile = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-undofile = 1
< ok 0 0
> exec ^
: string
| KILL unicode()
< error
| 1064 42000 …line 1 column 12 near "unicode()" 
> exec ^
: string
| KILL unicode (1)
< error
| 1064 42000 …line 1 column 12 near "unicode (1)" 
> exec ^
: string
| SET CONFIG tikv unicode = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-unicode = 1
< ok 0 0
> exec ^
: string
| KILL uninstall()
< ok 0 0
> exec ^
: string
| KILL uninstall (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv uninstall = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-uninstall = 1
< ok 0 0
> exec ^
: string
| KILL union()
< error
| 1064 42000 …line 1 column 10 near "union()" 
> exec ^
: string
| KILL union (1)
< error
| 1064 42000 …line 1 column 10 near "union (1)" 
> exec ^
: string
| SET CONFIG tikv union = 1
< error
| 1064 42000 …line 1 column 21 near "union = 1" 
> exec ^
: string
| SET CONFIG tikv a-union = 1
< error
| 1064 42000 …line 1 column 23 near "union = 1" 
> exec ^
: string
| KILL unique()
< error
| 1064 42000 …line 1 column 11 near "unique()" 
> exec ^
: string
| KILL unique (1)
< error
| 1064 42000 …line 1 column 11 near "unique (1)" 
> exec ^
: string
| SET CONFIG tikv unique = 1
< error
| 1064 42000 …line 1 column 22 near "unique = 1" 
> exec ^
: string
| SET CONFIG tikv a-unique = 1
< error
| 1064 42000 …line 1 column 24 near "unique = 1" 
> exec ^
: string
| KILL unknown()
< error
| 1064 42000 …line 1 column 12 near "unknown()" 
> exec ^
: string
| KILL unknown (1)
< error
| 1064 42000 …line 1 column 12 near "unknown (1)" 
> exec ^
: string
| SET CONFIG tikv unknown = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-unknown = 1
< ok 0 0
> exec ^
: string
| KILL unlimited()
< error
| 1064 42000 …line 1 column 14 near "unlimited()" 
> exec ^
: string
| KILL unlimited (1)
< error
| 1064 42000 …line 1 column 14 near "unlimited (1)" 
> exec ^
: string
| SET CONFIG tikv unlimited = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-unlimited = 1
< ok 0 0
> exec ^
: string
| KILL unlock()
< error
| 1064 42000 …line 1 column 11 near "unlock()" 
> exec ^
: string
| KILL unlock (1)
< error
| 1064 42000 …line 1 column 11 near "unlock (1)" 
> exec ^
: string
| SET CONFIG tikv unlock = 1
< error
| 1064 42000 …line 1 column 22 near "unlock = 1" 
> exec ^
: string
| SET CONFIG tikv a-unlock = 1
< error
| 1064 42000 …line 1 column 24 near "unlock = 1" 
> exec ^
: string
| KILL unregister()
< ok 0 0
> exec ^
: string
| KILL unregister (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv unregister = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-unregister = 1
< ok 0 0
> exec ^
: string
| KILL unset()
< error
| 1064 42000 …line 1 column 10 near "unset()" 
> exec ^
: string
| KILL unset (1)
< error
| 1064 42000 …line 1 column 10 near "unset (1)" 
> exec ^
: string
| SET CONFIG tikv unset = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-unset = 1
< ok 0 0
> exec ^
: string
| KILL unsigned()
< error
| 1064 42000 …line 1 column 13 near "unsigned()" 
> exec ^
: string
| KILL unsigned (1)
< error
| 1064 42000 …line 1 column 13 near "unsigned (1)" 
> exec ^
: string
| SET CONFIG tikv unsigned = 1
< error
| 1064 42000 …line 1 column 24 near "unsigned = 1" 
> exec ^
: string
| SET CONFIG tikv a-unsigned = 1
< error
| 1064 42000 …line 1 column 26 near "unsigned = 1" 
> exec ^
: string
| KILL until()
< error
| 1064 42000 …line 1 column 10 near "until()" 
> exec ^
: string
| KILL until (1)
< error
| 1064 42000 …line 1 column 10 near "until (1)" 
> exec ^
: string
| SET CONFIG tikv until = 1
< error
| 1064 42000 …line 1 column 21 near "until = 1" 
> exec ^
: string
| SET CONFIG tikv a-until = 1
< error
| 1064 42000 …line 1 column 23 near "until = 1" 
> exec ^
: string
| KILL until_ts()
< error
| 1064 42000 …line 1 column 13 near "until_ts()" 
> exec ^
: string
| KILL until_ts (1)
< error
| 1064 42000 …line 1 column 13 near "until_ts (1)" 
> exec ^
: string
| SET CONFIG tikv until_ts = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-until_ts = 1
< ok 0 0
> exec ^
: string
| KILL update()
< error
| 1064 42000 …line 1 column 11 near "update()" 
> exec ^
: string
| KILL update (1)
< error
| 1064 42000 …line 1 column 11 near "update (1)" 
> exec ^
: string
| SET CONFIG tikv update = 1
< error
| 1064 42000 …line 1 column 22 near "update = 1" 
> exec ^
: string
| SET CONFIG tikv a-update = 1
< error
| 1064 42000 …line 1 column 24 near "update = 1" 
> exec ^
: string
| KILL upgrade()
< ok 0 0
> exec ^
: string
| KILL upgrade (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv upgrade = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-upgrade = 1
< ok 0 0
> exec ^
: string
| KILL uri()
< ok 0 0
> exec ^
: string
| KILL uri (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv uri = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-uri = 1
< ok 0 0
> exec ^
: string
| KILL url()
< ok 0 0
> exec ^
: string
| KILL url (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv url = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-url = 1
< ok 0 0
> exec ^
: string
| KILL usage()
< error
| 1064 42000 …line 1 column 10 near "usage()" 
> exec ^
: string
| KILL usage (1)
< error
| 1064 42000 …line 1 column 10 near "usage (1)" 
> exec ^
: string
| SET CONFIG tikv usage = 1
< error
| 1064 42000 …line 1 column 21 near "usage = 1" 
> exec ^
: string
| SET CONFIG tikv a-usage = 1
< error
| 1064 42000 …line 1 column 23 near "usage = 1" 
> exec ^
: string
| KILL use()
< error
| 1064 42000 …line 1 column 8 near "use()" 
> exec ^
: string
| KILL use (1)
< error
| 1064 42000 …line 1 column 8 near "use (1)" 
> exec ^
: string
| SET CONFIG tikv use = 1
< error
| 1064 42000 …line 1 column 19 near "use = 1" 
> exec ^
: string
| SET CONFIG tikv a-use = 1
< error
| 1064 42000 …line 1 column 21 near "use = 1" 
> exec ^
: string
| KILL use_frm()
< ok 0 0
> exec ^
: string
| KILL use_frm (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv use_frm = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-use_frm = 1
< ok 0 0
> exec ^
: string
| KILL user()
< error
| 1064 42000 …line 1 column 9 near "user()" 
> exec ^
: string
| KILL user (1)
< error
| 1064 42000 …line 1 column 9 near "user (1)" 
> exec ^
: string
| SET CONFIG tikv user = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-user = 1
< ok 0 0
> exec ^
: string
| KILL user_resources()
< ok 0 0
> exec ^
: string
| KILL user_resources (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv user_resources = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-user_resources = 1
< ok 0 0
> exec ^
: string
| KILL using()
< error
| 1064 42000 …line 1 column 10 near "using()" 
> exec ^
: string
| KILL using (1)
< error
| 1064 42000 …line 1 column 10 near "using (1)" 
> exec ^
: string
| SET CONFIG tikv using = 1
< error
| 1064 42000 …line 1 column 21 near "using = 1" 
> exec ^
: string
| SET CONFIG tikv a-using = 1
< error
| 1064 42000 …line 1 column 23 near "using = 1" 
> exec ^
: string
| KILL utc_date()
< error
| 1064 42000 …line 1 column 13 near "utc_date()" 
> exec ^
: string
| KILL utc_date (1)
< error
| 1064 42000 …line 1 column 13 near "utc_date (1)" 
> exec ^
: string
| SET CONFIG tikv utc_date = 1
< error
| 1064 42000 …line 1 column 24 near "utc_date = 1" 
> exec ^
: string
| SET CONFIG tikv a-utc_date = 1
< error
| 1064 42000 …line 1 column 26 near "utc_date = 1" 
> exec ^
: string
| KILL utc_time()
< error
| 1064 42000 …line 1 column 13 near "utc_time()" 
> exec ^
: string
| KILL utc_time (1)
< error
| 1064 42000 …line 1 column 13 near "utc_time (1)" 
> exec ^
: string
| SET CONFIG tikv utc_time = 1
< error
| 1064 42000 …line 1 column 24 near "utc_time = 1" 
> exec ^
: string
| SET CONFIG tikv a-utc_time = 1
< error
| 1064 42000 …line 1 column 26 near "utc_time = 1" 
> exec ^
: string
| KILL utc_timestamp()
< error
| 1064 42000 …line 1 column 18 near "utc_timestamp()" 
> exec ^
: string
| KILL utc_timestamp (1)
< error
| 1064 42000 …line 1 column 18 near "utc_timestamp (1)" 
> exec ^
: string
| SET CONFIG tikv utc_timestamp = 1
< error
| 1064 42000 …line 1 column 29 near "utc_timestamp = 1" 
> exec ^
: string
| SET CONFIG tikv a-utc_timestamp = 1
< error
| 1064 42000 …line 1 column 31 near "utc_timestamp = 1" 
> exec ^
: string
| KILL utilization_limit()
< error
| 1064 42000 …line 1 column 22 near "utilization_limit()" 
> exec ^
: string
| KILL utilization_limit (1)
< error
| 1064 42000 …line 1 column 22 near "utilization_limit (1)" 
> exec ^
: string
| SET CONFIG tikv utilization_limit = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-utilization_limit = 1
< ok 0 0
> exec ^
: string
| KILL validate()
< ok 0 0
> exec ^
: string
| KILL validate (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv validate = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-validate = 1
< ok 0 0
> exec ^
: string
| KILL validation()
< error
| 1064 42000 …line 1 column 15 near "validation()" 
> exec ^
: string
| KILL validation (1)
< error
| 1064 42000 …line 1 column 15 near "validation (1)" 
> exec ^
: string
| SET CONFIG tikv validation = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-validation = 1
< ok 0 0
> exec ^
: string
| KILL value()
< error
| 1064 42000 …line 1 column 10 near "value()" 
> exec ^
: string
| KILL value (1)
< error
| 1064 42000 …line 1 column 10 near "value (1)" 
> exec ^
: string
| SET CONFIG tikv value = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-value = 1
< ok 0 0
> exec ^
: string
| KILL values()
< error
| 1064 42000 …line 1 column 11 near "values()" 
> exec ^
: string
| KILL values (1)
< error
| 1064 42000 …line 1 column 11 near "values (1)" 
> exec ^
: string
| SET CONFIG tikv values = 1
< error
| 1064 42000 …line 1 column 22 near "values = 1" 
> exec ^
: string
| SET CONFIG tikv a-values = 1
< error
| 1064 42000 …line 1 column 24 near "values = 1" 
> exec ^
: string
| KILL var_pop()
< error
| 1064 42000 …line 1 column 12 near "var_pop()" 
> exec ^
: string
| KILL var_pop (1)
< error
| 1064 42000 …line 1 column 12 near "var_pop (1)" 
> exec ^
: string
| SET CONFIG tikv var_pop = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-var_pop = 1
< ok 0 0
> exec ^
: string
| KILL var_samp()
< error
| 1064 42000 …line 1 column 13 near "var_samp()" 
> exec ^
: string
| KILL var_samp (1)
< error
| 1064 42000 …line 1 column 13 near "var_samp (1)" 
> exec ^
: string
| SET CONFIG tikv var_samp = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-var_samp = 1
< ok 0 0
> exec ^
: string
| KILL varbinary()
< error
| 1064 42000 …line 1 column 14 near "varbinary()" 
> exec ^
: string
| KILL varbinary (1)
< error
| 1064 42000 …line 1 column 14 near "varbinary (1)" 
> exec ^
: string
| SET CONFIG tikv varbinary = 1
< error
| 1064 42000 …line 1 column 25 near "varbinary = 1" 
> exec ^
: string
| SET CONFIG tikv a-varbinary = 1
< error
| 1064 42000 …line 1 column 27 near "varbinary = 1" 
> exec ^
: string
| KILL varchar()
< error
| 1064 42000 …line 1 column 12 near "varchar()" 
> exec ^
: string
| KILL varchar (1)
< error
| 1064 42000 …line 1 column 12 near "varchar (1)" 
> exec ^
: string
| SET CONFIG tikv varchar = 1
< error
| 1064 42000 …line 1 column 23 near "varchar = 1" 
> exec ^
: string
| SET CONFIG tikv a-varchar = 1
< error
| 1064 42000 …line 1 column 25 near "varchar = 1" 
> exec ^
: string
| KILL varcharacter()
< error
| 1064 42000 …line 1 column 17 near "varcharacter()" 
> exec ^
: string
| KILL varcharacter (1)
< error
| 1064 42000 …line 1 column 17 near "varcharacter (1)" 
> exec ^
: string
| SET CONFIG tikv varcharacter = 1
< error
| 1064 42000 …line 1 column 28 near "varcharacter = 1" 
> exec ^
: string
| SET CONFIG tikv a-varcharacter = 1
< error
| 1064 42000 …line 1 column 30 near "varcharacter = 1" 
> exec ^
: string
| KILL variables()
< error
| 1064 42000 …line 1 column 14 near "variables()" 
> exec ^
: string
| KILL variables (1)
< error
| 1064 42000 …line 1 column 14 near "variables (1)" 
> exec ^
: string
| SET CONFIG tikv variables = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-variables = 1
< ok 0 0
> exec ^
: string
| KILL variance()
< error
| 1064 42000 …line 1 column 13 near "variance()" 
> exec ^
: string
| KILL variance (1)
< error
| 1064 42000 …line 1 column 13 near "variance (1)" 
> exec ^
: string
| SET CONFIG tikv variance = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-variance = 1
< ok 0 0
> exec ^
: string
| KILL varying()
< error
| 1064 42000 …line 1 column 12 near "varying()" 
> exec ^
: string
| KILL varying (1)
< error
| 1064 42000 …line 1 column 12 near "varying (1)" 
> exec ^
: string
| SET CONFIG tikv varying = 1
< error
| 1064 42000 …line 1 column 23 near "varying = 1" 
> exec ^
: string
| SET CONFIG tikv a-varying = 1
< error
| 1064 42000 …line 1 column 25 near "varying = 1" 
> exec ^
: string
| KILL vcpu()
< ok 0 0
> exec ^
: string
| KILL vcpu (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv vcpu = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-vcpu = 1
< ok 0 0
> exec ^
: string
| KILL vector()
< error
| 1064 42000 …line 1 column 11 near "vector()" 
> exec ^
: string
| KILL vector (1)
< error
| 1064 42000 …line 1 column 11 near "vector (1)" 
> exec ^
: string
| SET CONFIG tikv vector = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-vector = 1
< ok 0 0
> exec ^
: string
| KILL verbose()
< error
| 1064 42000 …line 1 column 12 near "verbose()" 
> exec ^
: string
| KILL verbose (1)
< error
| 1064 42000 …line 1 column 12 near "verbose (1)" 
> exec ^
: string
| SET CONFIG tikv verbose = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-verbose = 1
< ok 0 0
> exec ^
: string
| KILL verify_key_constraints()
< ok 0 0
> exec ^
: string
| KILL verify_key_constraints (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv verify_key_constraints = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-verify_key_constraints = 1
< ok 0 0
> exec ^
: string
| KILL view()
< error
| 1064 42000 …line 1 column 9 near "view()" 
> exec ^
: string
| KILL view (1)
< error
| 1064 42000 …line 1 column 9 near "view (1)" 
> exec ^
: string
| SET CONFIG tikv view = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-view = 1
< ok 0 0
> exec ^
: string
| KILL virtual()
< error
| 1064 42000 …line 1 column 12 near "virtual()" 
> exec ^
: string
| KILL virtual (1)
< error
| 1064 42000 …line 1 column 12 near "virtual (1)" 
> exec ^
: string
| SET CONFIG tikv virtual = 1
< error
| 1064 42000 …line 1 column 23 near "virtual = 1" 
> exec ^
: string
| SET CONFIG tikv a-virtual = 1
< error
| 1064 42000 …line 1 column 25 near "virtual = 1" 
> exec ^
: string
| KILL visible()
< error
| 1064 42000 …line 1 column 12 near "visible()" 
> exec ^
: string
| KILL visible (1)
< error
| 1064 42000 …line 1 column 12 near "visible (1)" 
> exec ^
: string
| SET CONFIG tikv visible = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-visible = 1
< ok 0 0
> exec ^
: string
| KILL voter()
< error
| 1064 42000 …line 1 column 10 near "voter()" 
> exec ^
: string
| KILL voter (1)
< error
| 1064 42000 …line 1 column 10 near "voter (1)" 
> exec ^
: string
| SET CONFIG tikv voter = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-voter = 1
< ok 0 0
> exec ^
: string
| KILL voter_constraints()
< error
| 1064 42000 …line 1 column 22 near "voter_constraints()" 
> exec ^
: string
| KILL voter_constraints (1)
< error
| 1064 42000 …line 1 column 22 near "voter_constraints (1)" 
> exec ^
: string
| SET CONFIG tikv voter_constraints = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-voter_constraints = 1
< ok 0 0
> exec ^
: string
| KILL voters()
< error
| 1064 42000 …line 1 column 11 near "voters()" 
> exec ^
: string
| KILL voters (1)
< error
| 1064 42000 …line 1 column 11 near "voters (1)" 
> exec ^
: string
| SET CONFIG tikv voters = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-voters = 1
< ok 0 0
> exec ^
: string
| KILL wait()
< error
| 1064 42000 …line 1 column 9 near "wait()" 
> exec ^
: string
| KILL wait (1)
< error
| 1064 42000 …line 1 column 9 near "wait (1)" 
> exec ^
: string
| SET CONFIG tikv wait = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-wait = 1
< ok 0 0
> exec ^
: string
| KILL wait_tiflash_ready()
< error
| 1064 42000 …line 1 column 23 near "wait_tiflash_ready()" 
> exec ^
: string
| KILL wait_tiflash_ready (1)
< error
| 1064 42000 …line 1 column 23 near "wait_tiflash_ready (1)" 
> exec ^
: string
| SET CONFIG tikv wait_tiflash_ready = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-wait_tiflash_ready = 1
< ok 0 0
> exec ^
: string
| KILL warnings()
< error
| 1064 42000 …line 1 column 13 near "warnings()" 
> exec ^
: string
| KILL warnings (1)
< error
| 1064 42000 …line 1 column 13 near "warnings (1)" 
> exec ^
: string
| SET CONFIG tikv warnings = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-warnings = 1
< ok 0 0
> exec ^
: string
| KILL watch()
< error
| 1064 42000 …line 1 column 10 near "watch()" 
> exec ^
: string
| KILL watch (1)
< error
| 1064 42000 …line 1 column 10 near "watch (1)" 
> exec ^
: string
| SET CONFIG tikv watch = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-watch = 1
< ok 0 0
> exec ^
: string
| KILL week()
< error
| 1064 42000 …line 1 column 9 near "week()" 
> exec ^
: string
| KILL week (1)
< error
| 1064 42000 …line 1 column 9 near "week (1)" 
> exec ^
: string
| SET CONFIG tikv week = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-week = 1
< ok 0 0
> exec ^
: string
| KILL weight_string()
< error
| 1064 42000 …line 1 column 18 near "weight_string()" 
> exec ^
: string
| KILL weight_string (1)
< error
| 1064 42000 …line 1 column 18 near "weight_string (1)" 
> exec ^
: string
| SET CONFIG tikv weight_string = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-weight_string = 1
< ok 0 0
> exec ^
: string
| KILL when()
< error
| 1064 42000 …line 1 column 9 near "when()" 
> exec ^
: string
| KILL when (1)
< error
| 1064 42000 …line 1 column 9 near "when (1)" 
> exec ^
: string
| SET CONFIG tikv when = 1
< error
| 1064 42000 …line 1 column 20 near "when = 1" 
> exec ^
: string
| SET CONFIG tikv a-when = 1
< error
| 1064 42000 …line 1 column 22 near "when = 1" 
> exec ^
: string
| KILL where()
< error
| 1064 42000 …line 1 column 10 near "where()" 
> exec ^
: string
| KILL where (1)
< error
| 1064 42000 …line 1 column 10 near "where (1)" 
> exec ^
: string
| SET CONFIG tikv where = 1
< error
| 1064 42000 …line 1 column 21 near "where = 1" 
> exec ^
: string
| SET CONFIG tikv a-where = 1
< error
| 1064 42000 …line 1 column 23 near "where = 1" 
> exec ^
: string
| KILL while()
< error
| 1064 42000 …line 1 column 10 near "while()" 
> exec ^
: string
| KILL while (1)
< error
| 1064 42000 …line 1 column 10 near "while (1)" 
> exec ^
: string
| SET CONFIG tikv while = 1
< error
| 1064 42000 …line 1 column 21 near "while = 1" 
> exec ^
: string
| SET CONFIG tikv a-while = 1
< error
| 1064 42000 …line 1 column 23 near "while = 1" 
> exec ^
: string
| KILL width()
< error
| 1064 42000 …line 1 column 10 near "width()" 
> exec ^
: string
| KILL width (1)
< error
| 1064 42000 …line 1 column 10 near "width (1)" 
> exec ^
: string
| SET CONFIG tikv width = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-width = 1
< ok 0 0
> exec ^
: string
| KILL window()
< error
| 1064 42000 …line 1 column 11 near "window()" 
> exec ^
: string
| KILL window (1)
< error
| 1064 42000 …line 1 column 11 near "window (1)" 
> exec ^
: string
| SET CONFIG tikv window = 1
< error
| 1064 42000 …line 1 column 22 near "window = 1" 
> exec ^
: string
| SET CONFIG tikv a-window = 1
< error
| 1064 42000 …line 1 column 24 near "window = 1" 
> exec ^
: string
| KILL with()
< error
| 1064 42000 …line 1 column 9 near "with()" 
> exec ^
: string
| KILL with (1)
< error
| 1064 42000 …line 1 column 9 near "with (1)" 
> exec ^
: string
| SET CONFIG tikv with = 1
< error
| 1064 42000 …line 1 column 20 near "with = 1" 
> exec ^
: string
| SET CONFIG tikv a-with = 1
< error
| 1064 42000 …line 1 column 22 near "with = 1" 
> exec ^
: string
| KILL with_sys_table()
< error
| 1064 42000 …line 1 column 19 near "with_sys_table()" 
> exec ^
: string
| KILL with_sys_table (1)
< error
| 1064 42000 …line 1 column 19 near "with_sys_table (1)" 
> exec ^
: string
| SET CONFIG tikv with_sys_table = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-with_sys_table = 1
< ok 0 0
> exec ^
: string
| KILL without()
< error
| 1064 42000 …line 1 column 12 near "without()" 
> exec ^
: string
| KILL without (1)
< error
| 1064 42000 …line 1 column 12 near "without (1)" 
> exec ^
: string
| SET CONFIG tikv without = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-without = 1
< ok 0 0
> exec ^
: string
| KILL work()
< ok 0 0
> exec ^
: string
| KILL work (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv work = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-work = 1
< ok 0 0
> exec ^
: string
| KILL workload()
< error
| 1064 42000 …line 1 column 13 near "workload()" 
> exec ^
: string
| KILL workload (1)
< error
| 1064 42000 …line 1 column 13 near "workload (1)" 
> exec ^
: string
| SET CONFIG tikv workload = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-workload = 1
< ok 0 0
> exec ^
: string
| KILL wrapper()
< ok 0 0
> exec ^
: string
| KILL wrapper (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv wrapper = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-wrapper = 1
< ok 0 0
> exec ^
: string
| KILL write()
< error
| 1064 42000 …line 1 column 10 near "write()" 
> exec ^
: string
| KILL write (1)
< error
| 1064 42000 …line 1 column 10 near "write (1)" 
> exec ^
: string
| SET CONFIG tikv write = 1
< error
| 1064 42000 …line 1 column 21 near "write = 1" 
> exec ^
: string
| SET CONFIG tikv a-write = 1
< error
| 1064 42000 …line 1 column 23 near "write = 1" 
> exec ^
: string
| KILL x509()
< error
| 1064 42000 …line 1 column 9 near "x509()" 
> exec ^
: string
| KILL x509 (1)
< error
| 1064 42000 …line 1 column 9 near "x509 (1)" 
> exec ^
: string
| SET CONFIG tikv x509 = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-x509 = 1
< ok 0 0
> exec ^
: string
| KILL xa()
< ok 0 0
> exec ^
: string
| KILL xa (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv xa = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-xa = 1
< ok 0 0
> exec ^
: string
| KILL xid()
< ok 0 0
> exec ^
: string
| KILL xid (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv xid = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-xid = 1
< ok 0 0
> exec ^
: string
| KILL xml()
< ok 0 0
> exec ^
: string
| KILL xml (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv xml = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-xml = 1
< ok 0 0
> exec ^
: string
| KILL xor()
< error
| 1064 42000 …line 1 column 8 near "xor()" 
> exec ^
: string
| KILL xor (1)
< error
| 1064 42000 …line 1 column 8 near "xor (1)" 
> exec ^
: string
| SET CONFIG tikv xor = 1
< error
| 1064 42000 …line 1 column 19 near "xor = 1" 
> exec ^
: string
| SET CONFIG tikv a-xor = 1
< error
| 1064 42000 …line 1 column 21 near "xor = 1" 
> exec ^
: string
| KILL year()
< error
| 1064 42000 …line 1 column 9 near "year()" 
> exec ^
: string
| KILL year (1)
< error
| 1064 42000 …line 1 column 9 near "year (1)" 
> exec ^
: string
| SET CONFIG tikv year = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-year = 1
< ok 0 0
> exec ^
: string
| KILL year_month()
< error
| 1064 42000 …line 1 column 15 near "year_month()" 
> exec ^
: string
| KILL year_month (1)
< error
| 1064 42000 …line 1 column 15 near "year_month (1)" 
> exec ^
: string
| SET CONFIG tikv year_month = 1
< error
| 1064 42000 …line 1 column 26 near "year_month = 1" 
> exec ^
: string
| SET CONFIG tikv a-year_month = 1
< error
| 1064 42000 …line 1 column 28 near "year_month = 1" 
> exec ^
: string
| KILL zerofill()
< error
| 1064 42000 …line 1 column 13 near "zerofill()" 
> exec ^
: string
| KILL zerofill (1)
< error
| 1064 42000 …line 1 column 13 near "zerofill (1)" 
> exec ^
: string
| SET CONFIG tikv zerofill = 1
< error
| 1064 42000 …line 1 column 24 near "zerofill = 1" 
> exec ^
: string
| SET CONFIG tikv a-zerofill = 1
< error
| 1064 42000 …line 1 column 26 near "zerofill = 1" 
> exec ^
: string
| KILL zone()
< ok 0 0
> exec ^
: string
| KILL zone (1)
< ok 0 0
> exec ^
: string
| SET CONFIG tikv zone = 1
< ok 0 0
> exec ^
: string
| SET CONFIG tikv a-zone = 1
< ok 0 0
