-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBAdminFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| recommend index run
< ok 0 0
> exec ^
: string
| RECOMMEND INDEX RUN
< ok 0 0
> exec ^
: string
| recommend index run with a = 1
< ok 0 0
> exec ^
: string
| recommend index run with a = 1, b = 'x', c = true, d = false, e = null, f = 1.5, g = 1e3, h = 0x10, i = 0b1, j = x'10', k = b'1'
< ok 0 0
> exec ^
: string
| recommend index run with a = 'x' 'y'
< ok 0 0
> exec ^
: string
| recommend index run with a = _utf8mb4 'x'
< ok 0 0
> exec ^
: string
| recommend index run with a = _utf8mb4 0x10
< ok 0 0
> exec ^
: string
| recommend index run with a = _utf8mb4 b'1'
< ok 0 0
> exec ^
: string
| recommend index run with a = N'x'
< ok 0 0
> exec ^
: string
| recommend index run with a = -1
< error
| 1064 42000 …line 1 column 30 near "-1" 
> exec ^
: string
| recommend index run with a = (1)
< error
| 1064 42000 …line 1 column 30 near "(1)" 
> exec ^
: string
| recommend index run with a = 1 + 1
< error
| 1064 42000 …line 1 column 32 near "+ 1" 
> exec ^
: string
| recommend index run with a = b
< error
| 1064 42000 …line 1 column 30 near "b" 
> exec ^
: string
| recommend index run with a = default
< error
| 1064 42000 …line 1 column 36 near "default" 
> exec ^
: string
| recommend index run with a = @b
< error
| 1064 42000 …line 1 column 31 near "@b" 
> exec ^
: string
| recommend index run with a = now()
< error
| 1064 42000 …line 1 column 32 near "now()" 
> exec ^
: string
| recommend index run with a := 1
< error
| 1064 42000 …line 1 column 29 near ":= 1" 
> exec ^
: string
| recommend index run with a
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| recommend index run with a =
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| recommend index run with
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| recommend index run with a = 1,
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| recommend index run with a = 1 b = 2
< error
| 1064 42000 …line 1 column 32 near "b = 2" 
> exec ^
: string
| recommend index run with `a` = 1
< ok 0 0
> exec ^
: string
| recommend index run with config = 1
< ok 0 0
> exec ^
: string
| recommend index run with select = 1
< error
| 1064 42000 …line 1 column 31 near "select = 1" 
> exec ^
: string
| recommend index run with a.b = 1
< error
| 1064 42000 …line 1 column 27 near ".b = 1" 
> exec ^
: string
| recommend index run for
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| recommend index run for 'x'
< ok 0 0
> exec ^
: string
| recommend index run for "x"
< ok 0 0
> exec ^
: string
| recommend index run for 'x' 'y'
< error
| 1064 42000 …line 1 column 31 near "'y'" 
> exec ^
: string
| recommend index run for x
< error
| 1064 42000 …line 1 column 25 near "x" 
> exec ^
: string
| recommend index run for 1
< error
| 1064 42000 …line 1 column 25 near "1" 
> exec ^
: string
| recommend index run for 'x' with a = 1
< ok 0 0
> exec ^
: string
| recommend index run for 'x' with
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| recommend index run for 'x' a = 1
< error
| 1064 42000 …line 1 column 29 near "a = 1" 
> exec ^
: string
| recommend index run for (select 1)
< error
| 1064 42000 …line 1 column 25 near "(select 1)" 
> exec ^
: string
| recommend index run for select 1
< error
| 1064 42000 …line 1 column 30 near "select 1" 
> exec ^
: string
| recommend index run a = 1
< error
| 1064 42000 …line 1 column 21 near "a = 1" 
> exec ^
: string
| recommend index run with a = 1 for 'x'
< error
| 1064 42000 …line 1 column 34 near "for 'x'" 
> exec ^
: string
| recommend index show
< ok 0 0
> exec ^
: string
| recommend index show option
< error
| 1064 42000 …line 1 column 27 near "option" 
> exec ^
: string
| recommend index show 1
< error
| 1064 42000 …line 1 column 22 near "1" 
> exec ^
: string
| recommend index show with a = 1
< error
| 1064 42000 …line 1 column 25 near "with a = 1" 
> exec ^
: string
| recommend index apply
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| recommend index apply 1
< error
| 1105 HY000 unsupported action apply
> exec ^
: string
| recommend index apply 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| recommend index apply 0
< error
| 1105 HY000 unsupported action apply
> exec ^
: string
| recommend index apply 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| recommend index apply 9223372036854775807
< error
| 1105 HY000 unsupported action apply
> exec ^
: string
| recommend index apply 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| recommend index apply 18446744073709551616
< error
| 1064 42000 …line 1 column 42 near "18446744073709551616" 
> exec ^
: string
| recommend index apply -1
< error
| 1064 42000 …line 1 column 23 near "-1" 
> exec ^
: string
| recommend index apply 1.5
< error
| 1064 42000 …line 1 column 25 near "1.5" 
> exec ^
: string
| recommend index apply '1'
< error
| 1064 42000 …line 1 column 25 near "'1'" 
> exec ^
: string
| recommend index apply a
< error
| 1064 42000 …line 1 column 23 near "a" 
> exec ^
: string
| recommend index apply 1, 2
< error
| 1064 42000 …line 1 column 24 near ", 2" 
> exec ^
: string
| recommend index apply 1 2
< error
| 1064 42000 …line 1 column 25 near "2" 
> exec ^
: string
| recommend index ignore
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| recommend index ignore 1
< error
| 1105 HY000 unsupported action ignore
> exec ^
: string
| recommend index ignore 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| recommend index ignore a
< error
| 1064 42000 …line 1 column 24 near "a" 
> exec ^
: string
| recommend index ignore -1
< error
| 1064 42000 …line 1 column 24 near "-1" 
> exec ^
: string
| recommend index set
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| recommend index set a = 1
< ok 0 0
> exec ^
: string
| recommend index set a = 1, b = 2, c = 3
< ok 0 0
> exec ^
: string
| recommend index set a := 1
< error
| 1064 42000 …line 1 column 24 near ":= 1" 
> exec ^
: string
| recommend index set a
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| recommend index set a = 1,
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| recommend index set a = b
< error
| 1064 42000 …line 1 column 25 near "b" 
> exec ^
: string
| recommend index set with a = 1
< error
| 1064 42000 …line 1 column 24 near "with a = 1" 
> exec ^
: string
| recommend index set `a` = 'x' 'y'
< ok 0 0
> exec ^
: string
| recommend index
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| recommend index foo
< error
| 1064 42000 …line 1 column 19 near "foo" 
> exec ^
: string
| recommend
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| recommend indexes run
< error
| 1064 42000 …line 1 column 17 near "indexes run" 
> exec ^
: string
| recommend key run
< error
| 1064 42000 …line 1 column 13 near "key run" 
> exec ^
: string
| recommend `index` run
< error
| 1064 42000 …line 1 column 17 near "`index` run" 
> exec ^
: string
| `recommend` index run
< error
| 1064 42000 …line 1 column 11 near "`recommend` index run" 
> exec ^
: string
| recommend index `run`
< error
| 1064 42000 …line 1 column 21 near "`run`" 
> exec ^
: string
| recommend index `show`
< error
| 1064 42000 …line 1 column 22 near "`show`" 
> exec ^
: string
| recommend index `apply` 1
< error
| 1064 42000 …line 1 column 23 near "`apply` 1" 
> exec ^
: string
| recommend index `ignore` 1
< error
| 1064 42000 …line 1 column 24 near "`ignore` 1" 
> exec ^
: string
| recommend index `set` a = 1
< error
| 1064 42000 …line 1 column 21 near "`set` a = 1" 
> exec ^
: string
| recommend index run with a = _binary 'x'
< ok 0 0
> exec ^
: string
| recommend index run with a = 'x' 'y' 'z'
< ok 0 0
> exec ^
: string
| recommend index run with a = _utf8mb4 'x' 'y'
< error
| 1064 42000 …line 1 column 45 near "'y'" 
> exec ^
: string
| recommend index run with a = x'10' 'y'
< error
| 1064 42000 …line 1 column 38 near "'y'" 
> exec ^
: string
| recommend index run with a = 'y' x'10'
< error
| 1064 42000 …line 1 column 38 near "x'10'" 
> exec ^
: string
| recommend index run with a = 1 2
< error
| 1064 42000 …line 1 column 32 near "2" 
> exec ^
: string
| recommend index run with a = n'x' 'y'
< error
| 1064 42000 …line 1 column 37 near "'y'" 
> exec ^
: string
| recommend index run with a = 'x' n'y'
< error
| 1064 42000 …line 1 column 34 near "n'y'" 
> exec ^
: string
| recommend index run with a = _foo 'x'
< error
| 1064 42000 …line 1 column 33 near "_foo 'x'" 
> exec ^
: string
| recommend index run with a = b'1' b'1'
< error
| 1064 42000 …line 1 column 38 near "b'1'" 
> exec ^
: string
| recommend index run with status = 1, tidb = 2, recommend = 3, run = 4, apply = 5
< ok 0 0
> exec ^
: string
| recommend index run with a = .5
< ok 0 0
> exec ^
: string
| recommend index run with a = 1.
< ok 0 0
> exec ^
: string
| recommend index run with a = 18446744073709551616
< ok 0 0
> exec ^
: string
| recommend index run for n'x'
< error
| 1064 42000 …line 1 column 25 near "n'x'" 
> exec ^
: string
| recommend index run for _utf8mb4 'x'
< error
| 1064 42000 …line 1 column 32 near "_utf8mb4 'x'" 
> exec ^
: string
| recommend index run for x'10'
< error
| 1064 42000 …line 1 column 29 near "x'10'" 
> exec ^
: string
| recommend index apply 01
< error
| 1105 HY000 unsupported action apply
> exec ^
: string
| recommend index apply 01
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| recommend index apply 0x10
< error
| 1064 42000 …line 1 column 26 near "0x10" 
> exec ^
: string
| recommend index ignore 1e3
< error
| 1064 42000 …line 1 column 26 near "1e3" 
> exec ^
: string
| recommend index run with @a = 1
< error
| 1064 42000 …line 1 column 27 near "@a = 1" 
> exec ^
: string
| set config tikv log.level='info'
< ok 0 0
> exec ^
: string
| SET CONFIG TIKV LOG.LEVEL='info'
< ok 0 0
> exec ^
: string
| set config pd log.level='info'
< ok 0 0
> exec ^
: string
| set config tidb log.level='info'
< ok 0 0
> exec ^
: string
| set config tiflash log.level='info'
< ok 0 0
> exec ^
: string
| set config foo log.level='info'
< ok 0 0
> exec ^
: string
| set config `tikv` log.level='info'
< ok 0 0
> exec ^
: string
| set config 'tikv' log.level='info'
< ok 0 0
> exec ^
: string
| set config "tikv" log.level='info'
< ok 0 0
> exec ^
: string
| set config '127.0.0.1:3306' log.level='info'
< ok 0 0
> exec ^
: string
| set config 'a' 'b' log.level='info'
< error
| 1064 42000 …line 1 column 18 near "'b' log.level='info'" 
> exec ^
: string
| set config 1 log.level='info'
< error
| 1064 42000 …line 1 column 12 near "1 log.level='info'" 
> exec ^
: string
| set config 127.0.0.1 log.level='info'
< error
| 1064 42000 …line 1 column 16 near "127.0.0.1 log.level='info'" 
> exec ^
: string
| set config select a = 1
< error
| 1064 42000 …line 1 column 17 near "select a = 1" 
> exec ^
: string
| set config tikv a = 1
< ok 0 0
> exec ^
: string
| set config tikv a := 1
< ok 0 0
> exec ^
: string
| set config tikv a.b.c.d = 1
< ok 0 0
> exec ^
: string
| set config tikv a-b-c = 1
< ok 0 0
> exec ^
: string
| set config tikv a-b.c-d = 1
< ok 0 0
> exec ^
: string
| set config tikv `a`.`b` = 1
< ok 0 0
> exec ^
: string
| set config tikv `a`-`b` = 1
< ok 0 0
> exec ^
: string
| set config tikv `a.b` = 1
< ok 0 0
> exec ^
: string
| set config tikv 'a' = 1
< error
| 1064 42000 …line 1 column 19 near "'a' = 1" 
> exec ^
: string
| set config tikv a. = 1
< error
| 1064 42000 …line 1 column 20 near "= 1" 
> exec ^
: string
| set config tikv .a = 1
< error
| 1064 42000 …line 1 column 17 near ".a = 1" 
> exec ^
: string
| set config tikv a- = 1
< error
| 1064 42000 …line 1 column 20 near "= 1" 
> exec ^
: string
| set config tikv a..b = 1
< error
| 1064 42000 …line 1 column 19 near ".b = 1" 
> exec ^
: string
| set config tikv a--b = 1
< error
| 1064 42000 …line 1 column 19 near "-b = 1" 
> exec ^
: string
| set config tikv a-1 = 1
< error
| 1064 42000 …line 1 column 19 near "1 = 1" 
> exec ^
: string
| set config tikv a.1 = 1
< ok 0 0
> exec ^
: string
| set config tikv a-1b = 1
< ok 0 0
> exec ^
: string
| set config tikv a.select = 1
< ok 0 0
> exec ^
: string
| set config tikv select.a = 1
< ok 0 0
> exec ^
: string
| set config tikv a-select = 1
< error
| 1064 42000 …line 1 column 24 near "select = 1" 
> exec ^
: string
| set config tikv select-a = 1
< error
| 1064 42000 …line 1 column 22 near "select-a = 1" 
> exec ^
: string
| set config tikv a.b-select = 1
< error
| 1064 42000 …line 1 column 26 near "select = 1" 
> exec ^
: string
| set config tikv a . b = 1
< ok 0 0
> exec ^
: string
| set config tikv a - b = 1
< ok 0 0
> exec ^
: string
| set config tikv = 1
< error
| 1064 42000 …line 1 column 17 near "= 1" 
> exec ^
: string
| set config tikv a
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set config tikv a =
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set config tikv
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set config
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| set config tikv a = on
< ok 0 0
> exec ^
: string
| set config tikv a = off
< error
| 1054 42S22 Unknown column 'off' in ''
> exec ^
: string
| set config tikv a = off
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set config tikv a = binary
< ok 0 0
> exec ^
: string
| set config tikv a = default
< error
| 1105 HY000 Unknown DEFAULT for SET CONFIG
> exec ^
: string
| set config tikv a = default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set config tikv a = true
< ok 0 0
> exec ^
: string
| set config tikv a = 1 + 1
< ok 0 0
> exec ^
: string
| set config tikv a = (1)
< ok 0 0
> exec ^
: string
| set config tikv a = b
< error
| 1054 42S22 Unknown column 'b' in ''
> exec ^
: string
| set config tikv a = b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set config tikv a = now()
< ok 0 0
> exec ^
: string
| set config tikv a = @b
< ok 0 0
> exec ^
: string
| set config tikv a = @@sql_mode
< ok 0 0
> exec ^
: string
| set config tikv a = (select 1)
< ok 0 0
> exec ^
: string
| set config tikv a = 1, b = 2
< error
| 1064 42000 …line 1 column 22 near ", b = 2" 
> exec ^
: string
| set config tikv a = 1 b = 2
< error
| 1064 42000 …line 1 column 23 near "b = 2" 
> exec ^
: string
| set config tikv a = 1, tikv b = 2
< error
| 1064 42000 …line 1 column 22 near ", tikv b = 2" 
> exec ^
: string
| set config tikv a > 1
< error
| 1064 42000 …line 1 column 19 near "> 1" 
> exec ^
: string
| set config = 1
< ok 0 0
> exec ^
: string
| set config := 1
< ok 0 0
> exec ^
: string
| set config.a = 1
< ok 0 0
> exec ^
: string
| set config a = 1
< error
| 1064 42000 …line 1 column 14 near "= 1" 
> exec ^
: string
| set global config tikv a = 1
< error
| 1064 42000 …line 1 column 22 near "tikv a = 1" 
> exec ^
: string
| set session config tikv a = 1
< error
| 1064 42000 …line 1 column 23 near "tikv a = 1" 
> exec ^
: string
| set @@config = 1
< ok 0 0
> exec ^
: string
| set `config` tikv a = 1
< error
| 1064 42000 …line 1 column 17 near "tikv a = 1" 
> exec ^
: string
| set a = 1, config tikv b = 1
< error
| 1064 42000 …line 1 column 22 near "tikv b = 1" 
> exec ^
: string
| set config tikv a = 1 for x
< error
| 1064 42000 …line 1 column 25 near "for x" 
> exec ^
: string
| set global config = 1
< ok 0 0
> exec ^
: string
| set config config config = 1
< ok 0 0
> exec ^
: string
| set config tikv config = 1
< ok 0 0
> exec ^
: string
| set config n'a' a = 1
< error
| 1064 42000 …line 1 column 12 near "n'a' a = 1" 
> exec ^
: string
| set config _utf8mb4 'a' a = 1
< error
| 1064 42000 …line 1 column 19 near "_utf8mb4 'a' a = 1" 
> exec ^
: string
| set config tikv n'a' = 1
< error
| 1064 42000 …line 1 column 17 near "n'a' = 1" 
> exec ^
: string
| set config @a a = 1
< error
| 1064 42000 …line 1 column 13 near "@a a = 1" 
> exec ^
: string
| set config tikv a = 1;
< ok 0 0
> exec ^
: string
| set config tikv `a`.b-`c` = 1
< ok 0 0
> exec ^
: string
| set config tikv a.`select` = 1
< ok 0 0
> exec ^
: string
| set config tikv a.b. = 1
< error
| 1064 42000 …line 1 column 22 near "= 1" 
> exec ^
: string
| set config tikv a.-b = 1
< error
| 1064 42000 …line 1 column 19 near "-b = 1" 
> exec ^
: string
| set config tikv a-.b = 1
< error
| 1064 42000 …line 1 column 19 near ".b = 1" 
> exec ^
: string
| set config tikv 1a = 1
< ok 0 0
> exec ^
: string
| set config tikv 1a.b = 1
< ok 0 0
> exec ^
: string
| set config tikv a.1b = 1
< ok 0 0
> exec ^
: string
| set config 1a b = 1
< ok 0 0
> exec ^
: string
| set config tikv a-b-1 = 1
< error
| 1064 42000 …line 1 column 21 near "1 = 1" 
> exec ^
: string
| set config tikv a-1-b = 1
< error
| 1064 42000 …line 1 column 19 near "1-b = 1" 
> exec ^
: string
| set config tikv a.1.b = 1
< ok 0 0
> exec ^
: string
| set config tikv a.1.5 = 1
< ok 0 0
> exec ^
: string
| set config tikv a-1.5 = 1
< error
| 1064 42000 …line 1 column 21 near "1.5 = 1" 
> exec ^
: string
| set config tikv a.b.1 = 1
< ok 0 0
> exec ^
: string
| set config tikv a -b = 1
< ok 0 0
> exec ^
: string
| set config tikv a- b = 1
< ok 0 0
> exec ^
: string
| set config tikv a -- b
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set config tikv a = 1 = 2
< ok 0 0
> exec ^
: string
| set config tikv a = b = 2
< error
| 1054 42S22 Unknown column 'b' in ''
> exec ^
: string
| set config tikv a = b = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set config tikv a=1
< ok 0 0
> exec ^
: string
| set config tikv a:=1
< ok 0 0
> exec ^
: string
| set config tikv a.b=1
< ok 0 0
> exec ^
: string
| set config tikv a.b:=1
< ok 0 0
> exec ^
: string
| set config tikv a!=1
< error
| 1064 42000 …line 1 column 19 near "!=1" 
> exec ^
: string
| set config tikv a<=>1
< error
| 1064 42000 …line 1 column 20 near "<=>1" 
> exec ^
: string
| set config = 1, a = 2
< ok 0 0
> exec ^
: string
| set config := 1, a = 2
< ok 0 0
> exec ^
: string
| set config.a = 1, b = 2
< ok 0 0
> exec ^
: string
| set config config = 1
< error
| 1064 42000 …line 1 column 19 near "= 1" 
> exec ^
: string
| set config config.a = 1
< error
| 1064 42000 …line 1 column 18 near ".a = 1" 
> exec ^
: string
| set config config-a = 1
< error
| 1064 42000 …line 1 column 18 near "-a = 1" 
> exec ^
: string
| set config `config` a = 1
< ok 0 0
> exec ^
: string
| set config global a = 1
< ok 0 0
> exec ^
: string
| set config session a = 1
< ok 0 0
> exec ^
: string
| set config on a = 1
< error
| 1064 42000 …line 1 column 13 near "on a = 1" 
> exec ^
: string
| set config default a = 1
< error
| 1064 42000 …line 1 column 18 near "default a = 1" 
> exec ^
: string
| set config binary a = 1
< error
| 1064 42000 …line 1 column 17 near "binary a = 1" 
> exec ^
: string
| set config password a = 1
< ok 0 0
> exec ^
: string
| set config tikv `a`.1 = 1
< ok 0 0
> exec ^
: string
| set config tikv status.1 = 1
< ok 0 0
> exec ^
: string
| set config tikv select.1 = 1
< ok 0 0
> exec ^
: string
| set config tikv a.`b`.1 = 1
< ok 0 0
> exec ^
: string
| set config tikv a.	1 = 1
< ok 0 0
> exec ^
: string
| set config tikv a. /* c */ 1 = 1
< ok 0 0
> exec ^
: string
| set config tikv a./* c */1 = 1
< ok 0 0
> exec ^
: string
| set config tikv a . select = 1
< ok 0 0
> exec ^
: string
| set config tikv a. select = 1
< ok 0 0
> exec ^
: string
| set config tikv a .select = 1
< ok 0 0
> exec ^
: string
| set config tikv a	.	b = 1
< ok 0 0
> exec ^
: string
| set config tikv a-status = 1
< ok 0 0
> exec ^
: string
| set config tikv status-a = 1
< ok 0 0
> exec ^
: string
| set config tikv a.status-b = 1
< ok 0 0
> exec ^
: string
| set config tikv `a`-1b = 1
< ok 0 0
> exec ^
: string
| set config tikv a.1.select = 1
< ok 0 0
> exec ^
: string
| set config tikv a.1-select = 1
< error
| 1064 42000 …line 1 column 26 near "select = 1" 
> exec ^
: string
| set config tikv 1.a = 1
< error
| 1064 42000 …line 1 column 18 near "1.a = 1" 
> exec ^
: string
| set config tikv .1 = 1
< error
| 1064 42000 …line 1 column 18 near ".1 = 1" 
> exec ^
: string
| set config tikv a.b = 1 + 1
< ok 0 0
> exec ^
: string
| set config 'a' `b`.1 = 1
< ok 0 0
> exec ^
: string
| set config tikv @a = 1
< error
| 1064 42000 …line 1 column 18 near "@a = 1" 
> exec ^
: string
| set config tikv @@a = 1
< error
| 1064 42000 …line 1 column 19 near "@@a = 1" 
> exec ^
: string
| set config @@a b = 1
< error
| 1064 42000 …line 1 column 14 near "@@a b = 1" 
> exec ^
: string
| set config tikv a.@b = 1
< error
| 1064 42000 …line 1 column 20 near "@b = 1" 
> exec ^
: string
| pause backup logs
< ok 0 0
> exec ^
: string
| PAUSE BACKUP LOGS
< ok 0 0
> exec ^
: string
| pause backup logs gc_ttl='20220304'
< ok 0 0
> exec ^
: string
| pause backup logs gc_ttl='20220304' gc_ttl='2'
< ok 0 0
> exec ^
: string
| pause backup logs concurrency 4
< ok 0 0
> exec ^
: string
| pause backup logs foo = 1
< error
| 1064 42000 …line 1 column 21 near "foo = 1" 
> exec ^
: string
| pause backup logs gc_ttl
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| pause backup logs gc_ttl =
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| pause backup logs to 'x'
< error
| 1064 42000 …line 1 column 20 near "to 'x'" 
> exec ^
: string
| pause backup
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| pause
< error
| 1064 42000 …line 1 column 5 near "" 
> exec ^
: string
| pause logs
< error
| 1064 42000 …line 1 column 10 near "logs" 
> exec ^
: string
| pause backup log
< error
| 1064 42000 …line 1 column 16 near "log" 
> exec ^
: string
| pause `backup` logs
< error
| 1064 42000 …line 1 column 14 near "`backup` logs" 
> exec ^
: string
| pause backup `logs`
< error
| 1064 42000 …line 1 column 19 near "`logs`" 
> exec ^
: string
| `pause` backup logs
< error
| 1064 42000 …line 1 column 7 near "`pause` backup logs" 
> exec ^
: string
| pause backup logs 1
< error
| 1064 42000 …line 1 column 19 near "1" 
> exec ^
: string
| resume backup logs
< ok 0 0
> exec ^
: string
| RESUME BACKUP LOGS
< ok 0 0
> exec ^
: string
| resume backup logs gc_ttl='20220304'
< error
| 1064 42000 …line 1 column 25 near "gc_ttl='20220304'" 
> exec ^
: string
| resume backup logs concurrency 4
< error
| 1064 42000 …line 1 column 30 near "concurrency 4" 
> exec ^
: string
| resume backup logs to 'x'
< error
| 1064 42000 …line 1 column 21 near "to 'x'" 
> exec ^
: string
| resume backup
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| resume
< error
| 1064 42000 …line 1 column 6 near "" 
> exec ^
: string
| resume `backup` logs
< error
| 1064 42000 …line 1 column 15 near "`backup` logs" 
> exec ^
: string
| resume backup `logs`
< error
| 1064 42000 …line 1 column 20 near "`logs`" 
> exec ^
: string
| stop backup logs
< ok 0 0
> exec ^
: string
| STOP BACKUP LOGS
< ok 0 0
> exec ^
: string
| stop backup logs gc_ttl='2'
< error
| 1064 42000 …line 1 column 23 near "gc_ttl='2'" 
> exec ^
: string
| stop backup
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| stop backup logs to 'x'
< error
| 1064 42000 …line 1 column 19 near "to 'x'" 
> exec ^
: string
| stop `backup` logs
< error
| 1064 42000 …line 1 column 13 near "`backup` logs" 
> exec ^
: string
| stop backup `logs`
< error
| 1064 42000 …line 1 column 18 near "`logs`" 
> exec ^
: string
| purge backup logs from 'noop://'
< ok 0 0
> exec ^
: string
| PURGE BACKUP LOGS FROM 'noop://'
< ok 0 0
> exec ^
: string
| purge backup logs from 'noop://' until_ts='2012122304'
< ok 0 0
> exec ^
: string
| purge backup logs from "noop://"
< ok 0 0
> exec ^
: string
| purge backup logs from 'noop://' 'x'
< error
| 1064 42000 …line 1 column 36 near "'x'" 
> exec ^
: string
| purge backup logs from
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| purge backup logs from x
< error
| 1064 42000 …line 1 column 24 near "x" 
> exec ^
: string
| purge backup logs from 1
< error
| 1064 42000 …line 1 column 24 near "1" 
> exec ^
: string
| purge backup logs
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| purge backup logs to 'noop://'
< error
| 1064 42000 …line 1 column 20 near "to 'noop://'" 
> exec ^
: string
| purge backup logs from 'noop://' concurrency 4
< ok 0 0
> exec ^
: string
| purge backup logs from 'noop://' foo = 1
< error
| 1064 42000 …line 1 column 36 near "foo = 1" 
> exec ^
: string
| purge backup
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| purge
< error
| 1064 42000 …line 1 column 5 near "" 
> exec ^
: string
| purge `backup` logs from 'x'
< error
| 1064 42000 …line 1 column 14 near "`backup` logs from 'x'" 
> exec ^
: string
| purge backup `logs` from 'x'
< error
| 1064 42000 …line 1 column 19 near "`logs` from 'x'" 
> exec ^
: string
| purge binary logs to 'x'
< error
| 1064 42000 …line 1 column 12 near "binary logs to 'x'" 
> exec ^
: string
| purge master logs to 'x'
< error
| 1064 42000 …line 1 column 12 near "master logs to 'x'" 
> exec ^
: string
| purge binary logs before '2020-01-01'
< error
| 1064 42000 …line 1 column 12 near "binary logs before '2020-01-01'" 
> exec ^
: string
| purge a b c
< error
| 1064 42000 …line 1 column 7 near "a b c" 
> exec ^
: string
| purge logs to 'x'
< error
| 1064 42000 …line 1 column 10 near "logs to 'x'" 
> exec ^
: string
| show backup logs status
< ok 0 0
> exec ^
: string
| SHOW BACKUP LOGS STATUS
< ok 0 0
> exec ^
: string
| show backup logs status where a = 1
< error
| 1064 42000 …line 1 column 29 near "where a = 1" 
> exec ^
: string
| show backup logs status like 'x'
< error
| 1064 42000 …line 1 column 28 near "like 'x'" 
> exec ^
: string
| show backup logs status 1
< error
| 1064 42000 …line 1 column 25 near "1" 
> exec ^
: string
| show backup logs
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| show backup logs metadata from 'noop://'
< ok 0 0
> exec ^
: string
| SHOW BACKUP LOGS METADATA FROM 'noop://'
< ok 0 0
> exec ^
: string
| show backup logs metadata from "noop://"
< ok 0 0
> exec ^
: string
| show backup logs metadata from 'noop://' 'x'
< error
| 1064 42000 …line 1 column 44 near "'x'" 
> exec ^
: string
| show backup logs metadata from 'noop://' concurrency 4
< error
| 1064 42000 …line 1 column 52 near "concurrency 4" 
> exec ^
: string
| show backup logs metadata from x
< error
| 1064 42000 …line 1 column 32 near "x" 
> exec ^
: string
| show backup logs metadata from
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show backup logs metadata
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| show backup logs metadata 'noop://'
< error
| 1064 42000 …line 1 column 35 near "'noop://'" 
> exec ^
: string
| show backup metadata from 'noop://'
< ok 0 0
> exec ^
: string
| SHOW BACKUP METADATA FROM 'noop://'
< ok 0 0
> exec ^
: string
| show backup metadata from 'noop://' concurrency 4
< error
| 1064 42000 …line 1 column 47 near "concurrency 4" 
> exec ^
: string
| show backup metadata from 'x' 'y'
< error
| 1064 42000 …line 1 column 33 near "'y'" 
> exec ^
: string
| show backup metadata from x
< error
| 1064 42000 …line 1 column 27 near "x" 
> exec ^
: string
| show backup metadata from
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| show backup metadata
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| show backup
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| show `backup` logs status
< error
| 1064 42000 …line 1 column 13 near "`backup` logs status" 
> exec ^
: string
| show backup `logs` status
< error
| 1064 42000 …line 1 column 18 near "`logs` status" 
> exec ^
: string
| show backup logs `status`
< error
| 1064 42000 …line 1 column 25 near "`status`" 
> exec ^
: string
| show backup logs `metadata` from 'x'
< error
| 1064 42000 …line 1 column 27 near "`metadata` from 'x'" 
> exec ^
: string
| show backup `metadata` from 'x'
< error
| 1064 42000 …line 1 column 22 near "`metadata` from 'x'" 
> exec ^
: string
| show br job 1234
< ok 0 0
> exec ^
: string
| SHOW BR JOB 1234
< ok 0 0
> exec ^
: string
| show br job query 1234
< ok 0 0
> exec ^
: string
| SHOW BR JOB QUERY 1234
< ok 0 0
> exec ^
: string
| show br job 0
< ok 0 0
> exec ^
: string
| show br job 9223372036854775807
< ok 0 0
> exec ^
: string
| show br job 9223372036854775808
< error
| 1064 42000 …line 1 column 31 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| show br job 18446744073709551615
< error
| 1064 42000 …line 1 column 32 near ""18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| show br job 18446744073709551616
< error
| 1064 42000 …line 1 column 32 near "18446744073709551616" 
> exec ^
: string
| show br job -1
< error
| 1064 42000 …line 1 column 13 near "-1" 
> exec ^
: string
| show br job 1.5
< error
| 1064 42000 …line 1 column 15 near "1.5" 
> exec ^
: string
| show br job '1'
< error
| 1064 42000 …line 1 column 15 near "'1'" 
> exec ^
: string
| show br job a
< error
| 1064 42000 …line 1 column 13 near "a" 
> exec ^
: string
| show br job
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| show br job query
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| show br job query a
< error
| 1064 42000 …line 1 column 19 near "a" 
> exec ^
: string
| show br job query -1
< error
| 1064 42000 …line 1 column 19 near "-1" 
> exec ^
: string
| show br job query 9223372036854775808
< error
| 1064 42000 …line 1 column 37 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| show br job 1 query
< error
| 1064 42000 …line 1 column 19 near "query" 
> exec ^
: string
| show br jobs 1
< error
| 1064 42000 …line 1 column 12 near "jobs 1" 
> exec ^
: string
| show br 1
< error
| 1064 42000 …line 1 column 9 near "1" 
> exec ^
: string
| show br
< error
| 1064 42000 …line 1 column 7 near "" 
> exec ^
: string
| show `br` job 1
< error
| 1064 42000 …line 1 column 9 near "`br` job 1" 
> exec ^
: string
| show br `job` 1
< error
| 1064 42000 …line 1 column 13 near "`job` 1" 
> exec ^
: string
| show br job `query` 1
< error
| 1064 42000 …line 1 column 19 near "`query` 1" 
> exec ^
: string
| show br job 1 where a = 1
< error
| 1064 42000 …line 1 column 19 near "where a = 1" 
> exec ^
: string
| show br job 1, 2
< error
| 1064 42000 …line 1 column 14 near ", 2" 
> exec ^
: string
| cancel br job 1234
< ok 0 0
> exec ^
: string
| CANCEL BR JOB 1234
< ok 0 0
> exec ^
: string
| cancel br job
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| cancel br job 0
< ok 0 0
> exec ^
: string
| cancel br job 9223372036854775807
< ok 0 0
> exec ^
: string
| cancel br job 9223372036854775808
< error
| 1064 42000 …line 1 column 33 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| cancel br job 18446744073709551616
< error
| 1064 42000 …line 1 column 34 near "18446744073709551616" 
> exec ^
: string
| cancel br job -1
< error
| 1064 42000 …line 1 column 15 near "-1" 
> exec ^
: string
| cancel br job a
< error
| 1064 42000 …line 1 column 15 near "a" 
> exec ^
: string
| cancel br job 1, 2
< error
| 1064 42000 …line 1 column 16 near ", 2" 
> exec ^
: string
| cancel br job query 1
< error
| 1064 42000 …line 1 column 19 near "query 1" 
> exec ^
: string
| cancel br 1
< error
| 1064 42000 …line 1 column 11 near "1" 
> exec ^
: string
| cancel br
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| cancel `br` job 1
< error
| 1064 42000 …line 1 column 11 near "`br` job 1" 
> exec ^
: string
| cancel br `job` 1
< error
| 1064 42000 …line 1 column 15 near "`job` 1" 
> exec ^
: string
| `cancel` br job 1
< error
| 1064 42000 …line 1 column 8 near "`cancel` br job 1" 
> exec ^
: string
| pause backup logs gc_ttl 'x'
< ok 0 0
> exec ^
: string
| pause backup logs snapshot 5 minute ago
< ok 0 0
> exec ^
: string
| pause backup logs on duplicate replace
< ok 0 0
> exec ^
: string
| pause backup logs rate_limit 5 mb/second
< ok 0 0
> exec ^
: string
| pause backup logs , gc_ttl = 'x'
< error
| 1064 42000 …line 1 column 19 near ", gc_ttl = 'x'" 
> exec ^
: string
| purge backup logs from 'x' snapshot 5 month ago
< error
| 1064 42000 …MONTH is not a constant time interval and cannot be used here
> exec ^
: string
| purge backup logs from 'x' until_ts 'a' until_ts 'b'
< ok 0 0
> exec ^
: string
| show backup logs status;
< ok 0 0
> exec ^
: string
| show backup metadata from 'x';
< ok 0 0
> exec ^
: string
| stop backup logs;
< ok 0 0
> exec ^
: string
| show br job 1;
< ok 0 0
> exec ^
: string
| show br job  query  1;
< ok 0 0
> exec ^
: string
| show BR JOB 0x10
< error
| 1064 42000 …line 1 column 16 near "0x10" 
> exec ^
: string
| show br job 1e3
< error
| 1064 42000 …line 1 column 15 near "1e3" 
> exec ^
: string
| show br job 01
< ok 0 0
> exec ^
: string
| cancel br job 01
< ok 0 0
> exec ^
: string
| cancel br job 0x10
< error
| 1064 42000 …line 1 column 18 near "0x10" 
> exec ^
: string
| cancel
< error
| 1064 42000 …line 1 column 6 near "" 
> exec ^
: string
| cancel distribution job 1
< ok 0 0
> exec ^
: string
| CANCEL DISTRIBUTION JOB 1
< ok 0 0
> exec ^
: string
| cancel distribution job
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| cancel distribution job 0
< ok 0 0
> exec ^
: string
| cancel distribution job 9223372036854775807
< ok 0 0
> exec ^
: string
| cancel distribution job 9223372036854775808
< error
| 1064 42000 …line 1 column 43 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| cancel distribution job 18446744073709551616
< error
| 1064 42000 …line 1 column 44 near "18446744073709551616" 
> exec ^
: string
| cancel distribution job -1
< error
| 1064 42000 …line 1 column 25 near "-1" 
> exec ^
: string
| cancel distribution job a
< error
| 1064 42000 …line 1 column 25 near "a" 
> exec ^
: string
| cancel distribution job 1, 2
< error
| 1064 42000 …line 1 column 26 near ", 2" 
> exec ^
: string
| cancel distribution jobs 1
< error
| 1064 42000 …line 1 column 24 near "jobs 1" 
> exec ^
: string
| cancel distributions job 1
< error
| 1064 42000 …line 1 column 20 near "distributions job 1" 
> exec ^
: string
| cancel distribution 1
< error
| 1064 42000 …line 1 column 21 near "1" 
> exec ^
: string
| cancel `distribution` job 1
< error
| 1064 42000 …line 1 column 21 near "`distribution` job 1" 
> exec ^
: string
| cancel distribution `job` 1
< error
| 1064 42000 …line 1 column 25 near "`job` 1" 
> exec ^
: string
| cancel foo job 1
< error
| 1064 42000 …line 1 column 10 near "foo job 1" 
> exec ^
: string
| cancel distribution job 01
< ok 0 0
> exec ^
: string
| cancel import job 1
< ok 0 0
> exec ^
: string
| CANCEL IMPORT JOB 1
< ok 0 0
> exec ^
: string
| cancel import job
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| cancel import job 0
< ok 0 0
> exec ^
: string
| cancel import job 9223372036854775807
< ok 0 0
> exec ^
: string
| cancel import job 9223372036854775808
< error
| 1064 42000 …line 1 column 37 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| cancel import job -1
< error
| 1064 42000 …line 1 column 19 near "-1" 
> exec ^
: string
| cancel import job a
< error
| 1064 42000 …line 1 column 19 near "a" 
> exec ^
: string
| cancel import job 1, 2
< error
| 1064 42000 …line 1 column 20 near ", 2" 
> exec ^
: string
| cancel import jobs 1
< error
| 1064 42000 …line 1 column 18 near "jobs 1" 
> exec ^
: string
| cancel import 1
< error
| 1064 42000 …line 1 column 15 near "1" 
> exec ^
: string
| cancel `import` job 1
< error
| 1064 42000 …line 1 column 15 near "`import` job 1" 
> exec ^
: string
| kill 1
< ok 0 0
> exec ^
: string
| kill connection 1
< ok 0 0
> exec ^
: string
| kill query 1
< ok 0 0
> exec ^
: string
| kill tidb 1
< ok 0 0
> exec ^
: string
| KILL TIDB 1
< ok 0 0
> exec ^
: string
| kill tidb connection 1
< ok 0 0
> exec ^
: string
| kill tidb query 1
< ok 0 0
> exec ^
: string
| kill tidb
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| kill tidb connection
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| kill tidb query
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| kill tidb tidb 1
< error
| 1064 42000 …line 1 column 14 near "tidb 1" 
> exec ^
: string
| kill tidb connection_id()
< ok 0 0
> exec ^
: string
| kill tidb (connection_id())
< ok 0 0
> exec ^
: string
| kill tidb ((connection_id()))
< ok 0 0
> exec ^
: string
| kill tidb abs(1)
< ok 0 0
> exec ^
: string
| kill tidb abs(1, 2)
< ok 0 0
> exec ^
: string
| kill tidb replace('a','b','c')
< ok 0 0
> exec ^
: string
| kill tidb 1 + 1
< error
| 1064 42000 …line 1 column 13 near "+ 1" 
> exec ^
: string
| kill tidb -1
< error
| 1064 42000 …line 1 column 11 near "-1" 
> exec ^
: string
| kill tidb 1.5
< error
| 1064 42000 …line 1 column 13 near "1.5" 
> exec ^
: string
| kill tidb '1'
< error
| 1064 42000 …line 1 column 13 near "'1'" 
> exec ^
: string
| kill tidb a
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| kill tidb @a
< error
| 1064 42000 …line 1 column 12 near "@a" 
> exec ^
: string
| kill tidb (1)
< error
| 1064 42000 …line 1 column 12 near "1)" 
> exec ^
: string
| kill tidb connection connection_id()
< error
| 1064 42000 …line 1 column 34 near "connection_id()" 
> exec ^
: string
| kill tidb query connection_id()
< error
| 1064 42000 …line 1 column 29 near "connection_id()" 
> exec ^
: string
| kill tidb connection (1)
< error
| 1064 42000 …line 1 column 22 near "(1)" 
> exec ^
: string
| kill `tidb` 1
< error
| 1064 42000 …line 1 column 13 near "1" 
> exec ^
: string
| kill tidb 18446744073709551615
< ok 0 0
> exec ^
: string
| kill tidb 18446744073709551616
< error
| 1064 42000 …line 1 column 30 near "18446744073709551616" 
> exec ^
: string
| kill tidb 1, 2
< error
| 1064 42000 …line 1 column 12 near ", 2" 
> exec ^
: string
| kill tidb 1 2
< error
| 1064 42000 …line 1 column 13 near "2" 
> exec ^
: string
| kill connection tidb 1
< error
| 1064 42000 …line 1 column 20 near "tidb 1" 
> exec ^
: string
| kill query tidb 1
< error
| 1064 42000 …line 1 column 15 near "tidb 1" 
> exec ^
: string
| kill tidb()
< error
| 1064 42000 …line 1 column 11 near ")" 
> exec ^
: string
| kill tidb(1)
< error
| 1064 42000 …line 1 column 11 near "1)" 
> exec ^
: string
| kill `tidb`()
< ok 0 0
> exec ^
: string
| kill `tidb`(1)
< ok 0 0
> exec ^
: string
| kill 1 + 1
< error
| 1064 42000 …line 1 column 8 near "+ 1" 
> exec ^
: string
| kill -1
< error
| 1064 42000 …line 1 column 6 near "-1" 
> exec ^
: string
| kill 1.5
< error
| 1064 42000 …line 1 column 8 near "1.5" 
> exec ^
: string
| kill '1'
< error
| 1064 42000 …line 1 column 8 near "'1'" 
> exec ^
: string
| kill a
< error
| 1064 42000 …line 1 column 6 near "" 
> exec ^
: string
| kill a.b
< error
| 1064 42000 …line 1 column 7 near ".b" 
> exec ^
: string
| kill a.b()
< error
| 1064 42000 …line 1 column 7 near ".b()" 
> exec ^
: string
| kill @a
< error
| 1064 42000 …line 1 column 7 near "@a" 
> exec ^
: string
| kill (1)
< error
| 1064 42000 …line 1 column 7 near "1)" 
> exec ^
: string
| kill (connection_id())
< ok 0 0
> exec ^
: string
| kill ((connection_id()))
< ok 0 0
> exec ^
: string
| kill (connection_id()) + 1
< error
| 1064 42000 …line 1 column 24 near "+ 1" 
> exec ^
: string
| kill connection_id() + 1
< error
| 1064 42000 …line 1 column 22 near "+ 1" 
> exec ^
: string
| kill connection connection_id()
< error
| 1064 42000 …line 1 column 29 near "connection_id()" 
> exec ^
: string
| kill query connection_id()
< error
| 1064 42000 …line 1 column 24 near "connection_id()" 
> exec ^
: string
| kill connection (1)
< error
| 1064 42000 …line 1 column 17 near "(1)" 
> exec ^
: string
| kill query (1)
< error
| 1064 42000 …line 1 column 12 near "(1)" 
> exec ^
: string
| kill connection 1 + 1
< error
| 1064 42000 …line 1 column 19 near "+ 1" 
> exec ^
: string
| kill query a
< error
| 1064 42000 …line 1 column 12 near "a" 
> exec ^
: string
| kill connection a
< error
| 1064 42000 …line 1 column 17 near "a" 
> exec ^
: string
| kill connection
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| kill query
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| kill
< error
| 1064 42000 …line 1 column 4 near "" 
> exec ^
: string
| kill 18446744073709551615
< ok 0 0
> exec ^
: string
| kill 18446744073709551616
< error
| 1064 42000 …line 1 column 25 near "18446744073709551616" 
> exec ^
: string
| kill 9223372036854775808
< ok 0 0
> exec ^
: string
| kill 0x10
< error
| 1064 42000 …line 1 column 9 near "0x10" 
> exec ^
: string
| kill abs(1)
< ok 0 0
> exec ^
: string
| kill now()
< error
| 1064 42000 …line 1 column 8 near "now()" 
> exec ^
: string
| kill count(1)
< error
| 1064 42000 …line 1 column 10 near "count(1)" 
> exec ^
: string
| kill sum(1)
< error
| 1064 42000 …line 1 column 8 near "sum(1)" 
> exec ^
: string
| kill replace('a','b','c')
< ok 0 0
> exec ^
: string
| kill replace()
< error
| 1064 42000 …line 1 column 14 near ")" 
> exec ^
: string
| kill replace
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| kill `f`()
< ok 0 0
> exec ^
: string
| kill `connection_id`()
< ok 0 0
> exec ^
: string
| kill connection_id
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| kill current_user
< error
| 1064 42000 …line 1 column 17 near "current_user" 
> exec ^
: string
| kill current_user()
< error
| 1064 42000 …line 1 column 17 near "current_user()" 
> exec ^
: string
| kill database()
< error
| 1064 42000 …line 1 column 13 near "database()" 
> exec ^
: string
| kill if(1,2,3)
< error
| 1064 42000 …line 1 column 7 near "if(1,2,3)" 
> exec ^
: string
| kill cast(1 as signed)
< error
| 1064 42000 …line 1 column 9 near "cast(1 as signed)" 
> exec ^
: string
| kill recommend()
< error
| 1064 42000 …line 1 column 14 near "recommend()" 
> exec ^
: string
| kill pause()
< error
| 1064 42000 …line 1 column 10 near "pause()" 
> exec ^
: string
| kill cancel()
< error
| 1064 42000 …line 1 column 11 near "cancel()" 
> exec ^
: string
| kill br()
< error
| 1064 42000 …line 1 column 7 near "br()" 
> exec ^
: string
| kill connection()
< error
| 1064 42000 …line 1 column 16 near "()" 
> exec ^
: string
| kill query()
< error
| 1064 42000 …line 1 column 11 near "()" 
> exec ^
: string
| kill foo()
< ok 0 0
> exec ^
: string
| kill trim('a')
< error
| 1064 42000 …line 1 column 9 near "trim('a')" 
> exec ^
: string
| kill substring('a', 1)
< error
| 1064 42000 …line 1 column 14 near "substring('a', 1)" 
> exec ^
: string
| kill user()
< error
| 1064 42000 …line 1 column 9 near "user()" 
> exec ^
: string
| kill ascii('a')
< error
| 1064 42000 …line 1 column 10 near "ascii('a')" 
> exec ^
: string
| kill tidb_version()
< ok 0 0
> exec ^
: string
| kill recommend(1)
< error
| 1064 42000 …line 1 column 14 near "recommend(1)" 
> exec ^
: string
| kill run()
< error
| 1064 42000 …line 1 column 8 near "run()" 
> exec ^
: string
| kill apply()
< error
| 1064 42000 …line 1 column 10 near "apply()" 
> exec ^
: string
| kill config()
< error
| 1064 42000 …line 1 column 11 near "config()" 
> exec ^
: string
| kill plugins()
< error
| 1064 42000 …line 1 column 12 near "plugins()" 
> exec ^
: string
| kill distribution()
< error
| 1064 42000 …line 1 column 17 near "distribution()" 
> exec ^
: string
| kill job()
< error
| 1064 42000 …line 1 column 8 near "job()" 
> exec ^
: string
| kill count (1)
< ok 0 0
> exec ^
: string
| kill cast (1)
< error
| 1064 42000 …line 1 column 9 near "cast (1)" 
> exec ^
: string
| kill now ()
< error
| 1064 42000 …line 1 column 8 near "now ()" 
> exec ^
: string
| kill foo (1)
< ok 0 0
> exec ^
: string
| kill foo(1, 2)
< ok 0 0
> exec ^
: string
| kill foo(*)
< error
| 1064 42000 …line 1 column 10 near "*)" 
> exec ^
: string
| kill foo(distinct a)
< error
| 1064 42000 …line 1 column 17 near "distinct a)" 
> exec ^
: string
| kill foo(select 1)
< error
| 1064 42000 …line 1 column 15 near "select 1)" 
> exec ^
: string
| kill foo((select 1))
< ok 0 0
> exec ^
: string
| kill foo(a, b + 1)
< ok 0 0
> exec ^
: string
| kill foo(1,)
< error
| 1064 42000 …line 1 column 12 near ")" 
> exec ^
: string
| kill foo(
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| kill foo(1
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| kill (foo()
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| kill foo())
< error
| 1064 42000 …line 1 column 11 near ")" 
> exec ^
: string
| kill ((foo())
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| kill _utf8mb4()
< error
| 1064 42000 …line 1 column 13 near "_utf8mb4()" 
> exec ^
: string
| kill _foo()
< ok 0 0
> exec ^
: string
| kill _utf8mb4 ()
< error
| 1064 42000 …line 1 column 13 near "_utf8mb4 ()" 
> exec ^
: string
| kill N()
< ok 0 0
> exec ^
: string
| kill n'a'()
< error
| 1064 42000 …line 1 column 6 near "n'a'()" 
> exec ^
: string
| kill x()
< ok 0 0
> exec ^
: string
| kill b()
< ok 0 0
> exec ^
: string
| kill current_date()
< error
| 1064 42000 …line 1 column 17 near "current_date()" 
> exec ^
: string
| kill replace (1)
< ok 0 0
> exec ^
: string
| kill `replace`()
< ok 0 0
> exec ^
: string
| kill `replace`
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| kill `now`()
< ok 0 0
> exec ^
: string
| kill a.b(1)
< error
| 1064 42000 …line 1 column 7 near ".b(1)" 
> exec ^
: string
| kill .a()
< error
| 1064 42000 …line 1 column 6 near ".a()" 
> exec ^
: string
| kill now.a()
< error
| 1064 42000 …line 1 column 9 near ".a()" 
> exec ^
: string
| kill a$b()
< ok 0 0
> exec ^
: string
| kill 1a()
< ok 0 0
> exec ^
: string
| kill rank()
< error
| 1064 42000 …line 1 column 9 near "rank()" 
> exec ^
: string
| kill row_number()
< error
| 1064 42000 …line 1 column 15 near "row_number()" 
> exec ^
: string
| kill lag(1)
< error
| 1064 42000 …line 1 column 8 near "lag(1)" 
> exec ^
: string
| kill tidb foo (1)
< ok 0 0
> exec ^
: string
| kill connection foo()
< error
| 1064 42000 …line 1 column 19 near "foo()" 
> exec ^
: string
| kill foo() bar
< error
| 1064 42000 …line 1 column 14 near "bar" 
> exec ^
: string
| kill 1;
< ok 0 0
> exec ^
: string
| kill tidb 1;
< ok 0 0
> exec ^
: string
| kill connection 18446744073709551615
< ok 0 0
> exec ^
: string
| kill query 18446744073709551616
< error
| 1064 42000 …line 1 column 31 near "18446744073709551616" 
> exec ^
: string
| kill tidb query 18446744073709551615
< ok 0 0
> exec ^
: string
| kill 01
< ok 0 0
> exec ^
: string
| kill 0x1
< error
| 1064 42000 …line 1 column 8 near "0x1" 
> exec ^
: string
| kill @a()
< error
| 1064 42000 …line 1 column 7 near "@a()" 
> exec ^
: string
| kill foo(@a)
< ok 0 0
> exec ^
: string
| kill foo(a.*)
< error
| 1064 42000 …line 1 column 12 near "*)" 
> exec ^
: string
| kill foo(t.a)
< ok 0 0
> exec ^
: string
| kill foo(default)
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec ^
: string
| kill foo(default(a))
< ok 0 0
> exec ^
: string
| kill foo(interval 1 day)
< error
| 1064 42000 …line 1 column 24 near ")" 
> exec ^
: string
| kill foo(1 as a)
< error
| 1064 42000 …line 1 column 13 near "as a)" 
> exec ^
: string
| flush tidb plugins a
< ok 0 0
> exec ^
: string
| FLUSH TIDB PLUGINS A
< ok 0 0
> exec ^
: string
| flush tidb plugins a, b
< ok 0 0
> exec ^
: string
| flush tidb plugins a, b, c
< ok 0 0
> exec ^
: string
| flush tidb plugins `a`
< ok 0 0
> exec ^
: string
| flush tidb plugins status
< ok 0 0
> exec ^
: string
| flush tidb plugins tidb
< ok 0 0
> exec ^
: string
| flush tidb plugins select
< error
| 1064 42000 …line 1 column 25 near "select" 
> exec ^
: string
| flush tidb plugins 'a'
< error
| 1064 42000 …line 1 column 22 near "'a'" 
> exec ^
: string
| flush tidb plugins 1
< error
| 1064 42000 …line 1 column 20 near "1" 
> exec ^
: string
| flush tidb plugins a.b
< error
| 1064 42000 …line 1 column 21 near ".b" 
> exec ^
: string
| flush tidb plugins a,
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| flush tidb plugins a b
< error
| 1064 42000 …line 1 column 22 near "b" 
> exec ^
: string
| flush tidb plugins
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| flush tidb plugin a
< error
| 1064 42000 …line 1 column 17 near "plugin a" 
> exec ^
: string
| flush tidb
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| flush tidb, status
< error
| 1064 42000 …line 1 column 11 near ", status" 
> exec ^
: string
| flush status, tidb plugins a
< error
| 1064 42000 …line 1 column 13 near ", tidb plugins a" 
> exec ^
: string
| flush tidb plugins a, status
< ok 0 0
> exec ^
: string
| flush `tidb` plugins a
< error
| 1064 42000 …line 1 column 12 near "`tidb` plugins a" 
> exec ^
: string
| flush tidb `plugins` a
< error
| 1064 42000 …line 1 column 20 near "`plugins` a" 
> exec ^
: string
| flush local tidb plugins a
< ok 0 0
> exec ^
: string
| flush no_write_to_binlog tidb plugins a
< ok 0 0
> exec ^
: string
| flush plugins a
< error
| 1064 42000 …line 1 column 13 near "plugins a" 
> exec ^
: string
| flush foo
< error
| 1064 42000 …line 1 column 9 near "foo" 
> exec ^
: string
| flush foo, bar
< error
| 1064 42000 …line 1 column 9 near "foo, bar" 
> exec ^
: string
| flush status, privileges
< error
| 1064 42000 …line 1 column 13 near ", privileges" 
> exec ^
: string
| flush status status
< error
| 1064 42000 …line 1 column 19 near "status" 
> exec ^
: string
| flush
< error
| 1064 42000 …line 1 column 5 near "" 
> exec ^
: string
| flush user_resources
< error
| 1064 42000 …line 1 column 20 near "user_resources" 
> exec ^
: string
| flush relay logs
< error
| 1064 42000 …line 1 column 11 near "relay logs" 
> exec ^
: string
| flush optimizer_costs
< error
| 1064 42000 …line 1 column 21 near "optimizer_costs" 
> exec ^
: string
| flush des_key_file
< error
| 1064 42000 …line 1 column 18 near "des_key_file" 
> exec ^
: string
| flush query cache
< error
| 1064 42000 …line 1 column 11 near "query cache" 
> exec ^
: string
| flush binary
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| flush slow
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| flush tables t for export
< error
| 1064 42000 …line 1 column 18 near "for export" 
> exec ^
: string
| flush table with read
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| flush tables with
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| flush tables t with read lock
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| flush tables t with read lock
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| flush tables `t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| flush tables `t`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| flush tables a.b
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| flush tables a.b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| flush tables 't'
< error
| 1064 42000 …line 1 column 16 near "'t'" 
> exec ^
: string
| flush tables for export
< error
| 1064 42000 …line 1 column 16 near "for export" 
> exec ^
: string
| flush tidb tables
< error
| 1064 42000 …line 1 column 17 near "tables" 
> exec ^
: string
| flush tables, status
< error
| 1064 42000 …line 1 column 13 near ", status" 
> exec ^
: string
| flush status tables
< error
| 1064 42000 …line 1 column 19 near "tables" 
> exec ^
: string
| flush hosts, status
< error
| 1064 42000 …line 1 column 12 near ", status" 
> exec ^
: string
| flush privileges privileges
< error
| 1064 42000 …line 1 column 27 near "privileges" 
> exec ^
: string
| flush logs, status
< error
| 1064 42000 …line 1 column 11 near ", status" 
> exec ^
: string
| flush binary logs, status
< error
| 1064 42000 …line 1 column 18 near ", status" 
> exec ^
: string
| flush engine logs
< ok 0 0
> exec ^
: string
| flush error logs
< ok 0 0
> exec ^
: string
| flush general logs
< ok 0 0
> exec ^
: string
| flush slow logs
< ok 0 0
> exec ^
: string
| flush relay logs for channel 'x'
< error
| 1064 42000 …line 1 column 11 near "relay logs for channel 'x'" 
> exec ^
: string
| flush client_errors_summary
< ok 0 0
> exec ^
: string
| flush `status`
< error
| 1064 42000 …line 1 column 14 near "`status`" 
> exec ^
: string
| flush `privileges`
< error
| 1064 42000 …line 1 column 18 near "`privileges`" 
> exec ^
: string
| flush `logs`
< error
| 1064 42000 …line 1 column 12 near "`logs`" 
> exec ^
: string
| flush `hosts`
< error
| 1064 42000 …line 1 column 13 near "`hosts`" 
> exec ^
: string
| flush binary `logs`
< error
| 1064 42000 …line 1 column 19 near "`logs`" 
> exec ^
: string
| flush local status
< ok 0 0
> exec ^
: string
| flush local local status
< error
| 1064 42000 …line 1 column 17 near "local status" 
> exec ^
: string
| flush local no_write_to_binlog status
< error
| 1064 42000 …line 1 column 30 near "no_write_to_binlog status" 
> exec ^
: string
| flush no_write_to_binlog
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| flush local
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| flush local client_errors_summary
< ok 0 0
> exec ^
: string
| flush local logs
< ok 0 0
> exec ^
: string
| flush local tables
< ok 0 0
> exec ^
: string
| flush local hosts
< ok 0 0
> exec ^
: string
| flush tidb plugins @a
< error
| 1064 42000 …line 1 column 21 near "@a" 
> exec ^
: string
| alter range global placement policy mypolicy
< ok 0 0
> exec ^
: string
| ALTER RANGE GLOBAL PLACEMENT POLICY MYPOLICY
< ok 0 0
> exec ^
: string
| alter range global placement policy default
< ok 0 0
> exec ^
: string
| alter range meta placement policy mypolicy
< ok 0 0
> exec ^
: string
| alter range foo placement policy p
< ok 0 0
> exec ^
: string
| alter range `global` placement policy p
< ok 0 0
> exec ^
: string
| alter range 'global' placement policy p
< error
| 1064 42000 …line 1 column 20 near "'global' placement policy p" 
> exec ^
: string
| alter range select placement policy p
< error
| 1064 42000 …line 1 column 18 near "select placement policy p" 
> exec ^
: string
| alter range range placement policy p
< error
| 1064 42000 …line 1 column 17 near "range placement policy p" 
> exec ^
: string
| alter range a.b placement policy p
< error
| 1064 42000 …line 1 column 14 near ".b placement policy p" 
> exec ^
: string
| alter range 1 placement policy p
< error
| 1064 42000 …line 1 column 13 near "1 placement policy p" 
> exec ^
: string
| alter range global placement policy = p
< ok 0 0
> exec ^
: string
| alter range global placement policy = 'p'
< ok 0 0
> exec ^
: string
| alter range global placement policy 'p'
< ok 0 0
> exec ^
: string
| alter range global placement policy "p"
< ok 0 0
> exec ^
: string
| alter range global placement policy = default
< ok 0 0
> exec ^
: string
| alter range global placement policy set default
< ok 0 0
> exec ^
: string
| alter range global placement policy set
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| alter range global placement policy set p
< error
| 1064 42000 …line 1 column 41 near "p" 
> exec ^
: string
| alter range global placement policy = set default
< error
| 1064 42000 …line 1 column 41 near "set default" 
> exec ^
: string
| alter range global placement policy `p`
< ok 0 0
> exec ^
: string
| alter range global placement policy `default`
< ok 0 0
> exec ^
: string
| alter range global placement policy select
< error
| 1064 42000 …line 1 column 42 near "select" 
> exec ^
: string
| alter range global placement policy
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| alter range global placement policy =
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| alter range global placement
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| alter range global policy p
< error
| 1064 42000 …line 1 column 25 near "policy p" 
> exec ^
: string
| alter range global
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| alter range
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| alter range global placement policy p placement policy q
< error
| 1064 42000 …line 1 column 47 near "placement policy q" 
> exec ^
: string
| alter range global placement policy p, placement policy q
< error
| 1064 42000 …line 1 column 38 near ", placement policy q" 
> exec ^
: string
| alter range global default placement policy p
< error
| 1064 42000 …line 1 column 26 near "default placement policy p" 
> exec ^
: string
| alter range global primary_region = 'x'
< error
| 1064 42000 …line 1 column 33 near "primary_region = 'x'" 
> exec ^
: string
| alter range global placement policy 1
< error
| 1064 42000 …line 1 column 37 near "1" 
> exec ^
: string
| alter range placement policy p
< error
| 1064 42000 …line 1 column 28 near "policy p" 
> exec ^
: string
| alter range global placement `policy` p
< error
| 1064 42000 …line 1 column 37 near "`policy` p" 
> exec ^
: string
| alter range global `placement` policy p
< error
| 1064 42000 …line 1 column 30 near "`placement` policy p" 
> exec ^
: string
| alter `range` global placement policy p
< error
| 1064 42000 …line 1 column 13 near "`range` global placement policy p" 
> exec ^
: string
| alter range global placement policy default placement policy p
< error
| 1064 42000 …line 1 column 53 near "placement policy p" 
> exec ^
: string
| alter range default placement policy p
< error
| 1064 42000 …line 1 column 19 near "default placement policy p" 
> exec ^
: string
| alter range placement placement policy p
< ok 0 0
> exec ^
: string
| alter range policy placement policy p
< ok 0 0
> exec ^
: string
| alter range global placement policy placement
< ok 0 0
> exec ^
: string
| alter range global placement policy policy
< ok 0 0
> exec ^
: string
| alter range global placement policy = `default`
< ok 0 0
> exec ^
: string
| alter range global placement policy n'p'
< error
| 1064 42000 …line 1 column 37 near "n'p'" 
> exec ^
: string
| alter range global placement policy 'a' 'b'
< error
| 1064 42000 …line 1 column 43 near "'b'" 
> exec ^
: string
| alter range tikv placement policy p;
< ok 0 0
> exec ^
: string
| alter table t range placement policy p
< error
| 1064 42000 …line 1 column 19 near "range placement policy p" 
> exec ^
: string
| alter range @a placement policy p
< error
| 1064 42000 …line 1 column 14 near "@a placement policy p" 
> exec ^
: string
| alter range global placement policy @a
< error
| 1064 42000 …line 1 column 38 near "@a" 
> exec ^
: string
| select a.1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a.1b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.1b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a.1e5 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.1e5 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t.1.5 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.1.5 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `a`.1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `a`.1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select status.1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select status.1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a. 1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a. 1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a .1 from t
< error
| 1064 42000 …line 1 column 11 near ".1 from t" 
> exec ^
: string
| select a . 1 from t
< error
| 1064 42000 …line 1 column 12 near "1 from t" 
> exec ^
: string
| select select.1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select select.1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a.1 int)
< error
| 1103 42000 Incorrect table name 'a'
> exec ^
: string
| create table t (a.1 int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table a.1 (b int)
< ok 0 0
> exec ^
: string
| create table a.1b (b int)
< ok 0 0
> exec ^
: string
| select * from a.1
< error
| 1146 42S02 Table 'a.1' doesn't exist
> exec ^
: string
| select * from a.1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from a.1 b
< error
| 1146 42S02 Table 'a.1' doesn't exist
> exec ^
: string
| select * from a.1 b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table @a (b int)
< error
| 1064 42000 …line 1 column 15 near "@a (b int)" 
> exec ^
: string
| select * from @a
< error
| 1064 42000 …line 1 column 16 near "@a" 
> exec ^
: string
| alter table @a add column b int
< error
| 1064 42000 …line 1 column 14 near "@a add column b int" 
> exec ^
: string
| drop table @a
< error
| 1064 42000 …line 1 column 13 near "@a" 
> exec ^
: string
| select 1 as @a
< error
| 1064 42000 …line 1 column 14 near "@a" 
> exec ^
: string
| select 1 @a
< error
| 1064 42000 …line 1 column 11 near "@a" 
> exec ^
: string
| create table t (@a int)
< error
| 1064 42000 …line 1 column 18 near "@a int)" 
> exec ^
: string
| select foo(*)
< error
| 1064 42000 …line 1 column 12 near "*)" 
> exec ^
: string
| select abs(*)
< error
| 1064 42000 …line 1 column 12 near "*)" 
> exec ^
: string
| create table _utf8 (a int)
< error
| 1064 42000 …line 1 column 18 near "_utf8 (a int)" 
> exec ^
: string
| create table `_utf8` (a int)
< ok 0 0
> exec ^
: string
| select 1 _utf8
< error
| 1064 42000 …line 1 column 14 near "_utf8" 
> exec ^
: string
| select 1 as _utf8
< error
| 1064 42000 …line 1 column 17 near "_utf8" 
> exec ^
: string
| select a._utf8 from t
< error
| 1064 42000 …line 1 column 14 near "_utf8 from t" 
> exec ^
: string
| create table t (_utf8 int)
< error
| 1064 42000 …line 1 column 21 near "_utf8 int)" 
> exec ^
: string
| create table _foo (a int)
< ok 0 0
> exec ^
: string
| select _foo from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select _foo from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from _utf8
< error
| 1064 42000 …line 1 column 19 near "_utf8" 
> exec ^
: string
| drop table _utf8
< error
| 1064 42000 …line 1 column 16 near "_utf8" 
> exec ^
: string
| select t.*1 from t
< error
| 1064 42000 …line 1 column 11 near "1 from t" 
> exec ^
: string
| select t.*.5 from t
< error
| 1064 42000 …line 1 column 12 near ".5 from t" 
> exec ^
: string
| select t.* , a.1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.* , a.1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (1), a.b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select (1), a.b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (select 1) union (select a.1 from t)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| (select 1) union (select a.1 from t)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a.1.2.3 from t
< error
| 1064 42000 …line 1 column 13 near ".3 from t" 
> exec ^
: string
| select 1.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a.b.c.1 from t
< error
| 1064 42000 …line 1 column 13 near ".1 from t" 
> exec ^
: string
| set config tikv a = *
< error
| 1064 42000 …line 1 column 21 near "*" 
> exec ^
: string
| set config tikv a = t.*
< error
| 1064 42000 …line 1 column 23 near "*" 
> exec ^
: string
| set config tikv a = interval 1 day
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| set config tikv a = default(b)
< error
| 1105 HY000 Unknown DEFAULT for SET CONFIG
> exec ^
: string
| set config tikv a = default(b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set config tikv a = default + 1
< error
| 1064 42000 …line 1 column 29 near "+ 1" 
> exec ^
: string
| set config tikv a = on + 1
< error
| 1064 42000 …line 1 column 24 near "+ 1" 
> exec ^
: string
| set config tikv a = binary 'x'
< ok 0 0
> exec ^
: string
| set config tikv a = row(1, 2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| set config tikv a = row(1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set config tikv a = (1, 2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| set config tikv a = (1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set config tikv a = @x := 1
< ok 0 0
> exec ^
: string
| recommend index run with a = n 'x'
< error
| 1064 42000 …line 1 column 30 near "n 'x'" 
> exec ^
: string
| recommend index run with a = N  'x'
< error
| 1064 42000 …line 1 column 30 near "N  'x'" 
> exec ^
: string
| purge backup logs 'noop://'
< error
| 1064 42000 …line 1 column 27 near "'noop://'" 
> exec ^
: string
| flush tables with lock
< error
| 1064 42000 …line 1 column 22 near "lock" 
> exec ^
: string
| flush tables with read
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| flush tables t with lock
< error
| 1064 42000 …line 1 column 24 near "lock" 
> exec ^
: string
| show backup logs metadata 'x'
< error
| 1064 42000 …line 1 column 29 near "'x'" 
> exec ^
: string
| kill query foo()
< error
| 1064 42000 …line 1 column 14 near "foo()" 
> exec ^
: string
| kill tidb query foo()
< error
| 1064 42000 …line 1 column 19 near "foo()" 
> exec ^
: string
| kill connection foo
< error
| 1064 42000 …line 1 column 19 near "foo" 
> exec ^
: string
| select _utf8
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| select _utf8 from t
< error
| 1064 42000 …line 1 column 17 near "from t" 
> exec ^
: string
| select 1 _utf8
< error
| 1064 42000 …line 1 column 14 near "_utf8" 
> exec ^
: string
| select 1 as _utf8
< error
| 1064 42000 …line 1 column 17 near "_utf8" 
