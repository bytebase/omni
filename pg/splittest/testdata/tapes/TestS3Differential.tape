-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestS3Differential default#1
> exec
| COPY dix FROM stdin;
| a;b
| c;d;e
| \.
|  COPY dwk FROM stdin;
| a;b
| 1	\N
| \.
|
| COPY ujc FROM stdin;
| a;b
| c;d;e
| \.
|
| + /*/*/*select 1;;;a;b*/*/*/ ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dix FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dix FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = "select 1;" + /*$x$;;*/ , e'\';\';' AS x;
| SELECT $pin$end;; );($pin$ SELECT $$END ;end;end;$$ 1 $cur$/*;$cux$$cur$ +; -- trailing note;
| SELECT "$;$--;" AS x ";;;--;" + 'end;end;' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "select 1;" + /*$x$;;*/ , e'\';\';' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "select 1;" + /*$x$;;*/ , e'\';\';' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + b'01;10' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + b'01;10' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + b'01;10' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x --END ;
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --END ;
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --END ;
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x CREATE FUNCTION lx_g() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(--;'; SELECT 1; END col;
| COPY hes FROM stdin;
| 1	\N
| a;b
| \.
|
| 1 we$f$a , CREATE FUNCTION xtll() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END SELECT; -- trailing note;
| COPY uuq FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION lx_g() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(--;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION lx_g() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(--;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT '$;$/*;;' FROM t WHERE a = U&'$;$;;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 45
> exec
| BEGIN
< ok 0 -
> exec
| SELECT '$;$/*;;' FROM t WHERE a = U&'$;$;;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 45
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT '$;$/*;;' FROM t WHERE a = U&'$;$;;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 45
> exec
| COPY _kx FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _kx FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _kx FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY pgv FROM stdin;
| a;b
| a;b
| \.
|  , (SELECT 1; SELECT 2) SELECT tq$e$t AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pgv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pgv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + CREATE FUNCTION wyoz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END FROM t WHERE a = /*; ;end;*/ SELECT --/*;);(
|  FROM t WHERE a =;
| AS x U&'; ' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION wyoz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END FROM t WHERE a = /*; ;end;*/ SELECT --/*;);(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION wyoz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END FROM t WHERE a = /*; ;end;*/ SELECT --/*;);(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY a_b FROM stdin;
| 1	\N
| c;d;e
| \.
| FROM t WHERE a = U&'select 1;a;b' , (SELECT 1; SELECT 2) AS x ik$q$b FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY a_b FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY a_b FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 --; ;;$;$
|  FROM t WHERE a =;
| COPY qfn FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --; ;;$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --; ;;$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 '; ' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '; ' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '; ' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 /*/*/*);(end;$x$;*/*/*/ FROM t WHERE a = (SELECT 1; SELECT 2) ,; SELECT --$;$
|  AS x (SELECT 1; SELECT 2) ,;
| , E'\\\\' FROM t WHERE a = --);(
|  1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*);(end;$x$;*/*/*/ FROM t WHERE a = (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*);(end;$x$;*/*/*/ FROM t WHERE a = (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 gr$q$x FROM t WHERE a = '; \' col;
| COPY ymo FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 gr$q$x FROM t WHERE a = '; \' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 gr$q$x FROM t WHERE a = '; \' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT /*/*/*);(*/*/*/ 1; -- trailing note;
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;COPY b_a FROM stdin;
| a;b
| a;b
| \.
|  , ';' 1 --;;
|  SELECT e''';\';\';' +;
< error
| 42601 syntax error at or near "AS"
| position: 45
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*/*);(*/*/*/ 1;
< ok 1 -
> exec
|  -- trailing note;
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 20
> exec
| ROLLBACK
< ok 0 -
> exec
|  -- trailing note;
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 20
> exec
| AS x kz$h$d ,;
| , ";;$;$" FROM t WHERE a = U&'/*;END ;' AS x $mrcl$; );(;$mrcl$ AS x;
| AS x x'01;10' ,; -- trailing note;
| SELECT U&');(/*;select 1;' , $$end;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x kz$h$d ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x kz$h$d ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY sgf FROM stdin;
| c;d;e
| a;b
| \.
|
| 1 X'01;10' 1 E'\n;\n;'';' ,; -- trailing note;
| COPY cwr FROM stdin;
| a;b
| 1	\N
| \.
|
| 1 (SELECT 1; SELECT 2) AS x $$$x$;$$ SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sgf FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sgf FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT '/*;;\' SELECT "a;b" , U&'end;' +;
< error
| 42601 syntax error at or near ""a;b""
| position: 23
> exec
| BEGIN
< ok 0 -
> exec
| SELECT '/*;;\' SELECT "a;b" , U&'end;' +;
< error
| 42601 syntax error at or near ""a;b""
| position: 23
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT '/*;;\' SELECT "a;b" , U&'end;' +;
< error
| 42601 syntax error at or near ""a;b""
| position: 23
> exec
| COPY mxm FROM stdin;
| a;b
| a;b
| \.
|
| SELECT e''';' col;
| COPY lnt FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mxm FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mxm FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY qlp FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qlp FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qlp FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x 'end;);(' ,;
| , ");(;;end;" , '--;; \' , ";; " 1;
| col --);(
|  col --end;select 1;
|  col ";;/*;$;$" 1;
| FROM t WHERE a = ';' , /*--;;/ *;*/ FROM t WHERE a = --; select 1;END ;
|  ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x 'end;);(' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x 'end;);(' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = --; --;;
|  FROM t WHERE a = aa$i$j col $$a;b$x$;$$ FROM t WHERE a =;
| 1 ve$k$f FROM t WHERE a = U&'a;b--;' + --select 1;end;--;
|  ,;
| SELECT e'\';\';\n;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --; --;;
|  FROM t WHERE a = aa$i$j col $$a;b$x$;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --; --;;
|  FROM t WHERE a = aa$i$j col $$a;b$x$;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = lb$u$y FROM t WHERE a = X'01;10' + e'\n;' SELECT;FROM t WHERE a = x'01;10' FROM t WHERE a = CREATE FUNCTION mlfg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$/*;'; END col /*/*--;*/*/ SELECT;
| COPY ckj FROM stdin;
| c;d;e
| c;d;e
| \.
|
| SELECT --);(end;$;$
|  , x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = lb$u$y FROM t WHERE a = X'01;10' + e'\n;' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = lb$u$y FROM t WHERE a = X'01;10' + e'\n;' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 "end;" SELECT $qbml$;;;$qbmx$$qbml$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "end;" SELECT $qbml$;;;$qbmx$$qbml$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "end;" SELECT $qbml$;;;$qbmx$$qbml$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , CREATE FUNCTION xawv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END + b'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION xawv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END + b'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION xawv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END + b'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 e'\\'';\n;' + U&'select 1;' col iz$k$i +;, $$--;end;;;$$ , --);(; 
|  FROM t WHERE a = ");(" col; -- trailing note;
| AS x /*;$x$;select 1;*/ + (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\\'';\n;' + U&'select 1;' col iz$k$i +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\\'';\n;' + U&'select 1;' col iz$k$i +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col U&'--;' AS x e''';\';' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'--;' AS x e''';\';' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'--;' AS x e''';\';' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY qvw FROM stdin;
| c;d;e
| 1	\N
| \.
|  SELECT X'01;10' , /*$x$;*/ FROM t WHERE a =;
| SELECT ';select 1;--;\' ,;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qvw FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qvw FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + CREATE FUNCTION _zvx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END 1 e'\\\n;\';' FROM t WHERE a = --/*;; );(
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION _zvx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END 1 e'\\\n;\';' FROM t WHERE a = --/*;; );(
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION _zvx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END 1 e'\\\n;\';' FROM t WHERE a = --/*;; );(
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , '$;$a;b' SELECT;
| FROM t WHERE a = --;;select 1;
|  FROM t WHERE a =;
| SELECT $$;$$ SELECT /*$x$;--;;;*/ SELECT --; 
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , '$;$a;b' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , '$;$a;b' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY nd_ FROM stdin;
| c;d;e
| c;d;e
| \.
|
| col $$select 1;$$ FROM t WHERE a =;
| FROM t WHERE a = 'a;b\' 1 X'01;10' FROM t WHERE a = /*/*/ *;;*/*/ +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY nd_ FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY nd_ FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT --;);(
|  FROM t WHERE a = ");();(" + ho$q$p SELECT;
| AS x x'01;10' AS x;
| , /*--;; */ ,;
< error
| 42601 syntax error at or near "SELECT"
| position: 51
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --;);(
|  FROM t WHERE a = ");();(" + ho$q$p SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 51
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --;);(
|  FROM t WHERE a = ");();(" + ho$q$p SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 51
> exec
| COPY zxx FROM stdin;
| 1	\N
| 1	\N
| \.
|
| COPY hpd FROM stdin;
| a;b
| a;b
| \.
|
| col ";;" col U&'a;bEND ;--;' 1 "a;b" 1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zxx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zxx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY _xy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _xy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _xy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY oux FROM stdin;
| a;b
| 1	\N
| \.
|  FROM t WHERE a = (SELECT 1; SELECT 2) AS x --select 1;END ;
|  + (SELECT 1; SELECT 2) AS x;
| col --end;a;b$x$;
|  AS x --;$;$
|  1 (SELECT 1; SELECT 2) col; -- trailing note;
| COPY xlr FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY oux FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY oux FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY fen FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fen FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fen FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x ';$;$' 1 CREATE FUNCTION zztv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; '; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ';$;$' 1 CREATE FUNCTION zztv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; '; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ';$;$' 1 CREATE FUNCTION zztv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; '; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = '$;$end;' + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = '$;$end;' + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = '$;$end;' + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 E'\';\\'';' , sj$f$e col;
| + dy$m$e +;
| 1 ");(; select 1;" col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 E'\';\\'';' , sj$f$e col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 E'\';\\'';' , sj$f$e col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY xjr FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| 1 ');(' AS x;
| COPY fva FROM stdin;
| a;b
| a;b
| \.
|
| COPY g_k FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xjr FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xjr FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 $baj$;;;;$bax$$baj$ + U&'/*;;;--;' col;
| , e'\\'';\';' col ';;END ;/*;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $baj$;;;;$bax$$baj$ + U&'/*;;;--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $baj$;;;;$bax$$baj$ + U&'/*;;;--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col --END ;; select 1;
|  FROM t WHERE a = E'\';'';' + cc$v$q col; col x'01;10' col --/*;$;$;;
|  ,;AS x --a;b$x$;
|  SELECT "; );(;;" + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --END ;; select 1;
|  FROM t WHERE a = E'\';'';' + cc$v$q col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --END ;; select 1;
|  FROM t WHERE a = E'\';'';' + cc$v$q col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col "$x$;$x$;;;" AS x $$$;$$$ 1 (SELECT 1; SELECT 2) col;
| FROM t WHERE a = $bo$/*;$bx$$bo$ FROM t WHERE a = X'01;10' col; col --/*;
|  + U&'a;b' 1 /*/*/*END ;);($;$*/*/*/ AS x;
| , CREATE FUNCTION jrhh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "$x$;$x$;;;" AS x $$$;$$$ 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "$x$;$x$;;;" AS x $$$;$$$ 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , "; /*;" SELECT CREATE FUNCTION eaox() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;; );('; SELECT 1; END FROM t WHERE a = $$--;$$ +; + --);(--;
|  AS x; COPY __s FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "; /*;" SELECT CREATE FUNCTION eaox() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;; );('; SELECT 1; END FROM t WHERE a = $$--;$$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "; /*;" SELECT CREATE FUNCTION eaox() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;; );('; SELECT 1; END FROM t WHERE a = $$--;$$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x U&'$x$;a;bselect 1;' +;COPY otx FROM stdin;
| a;b
| c;d;e
| \.
|
| + ';; select 1;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'$x$;a;bselect 1;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'$x$;a;bselect 1;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) +; SELECT --end;; select 1;
|  , X'01;10' AS x;COPY jgu FROM stdin;
| a;b
| c;d;e
| \.
|
| SELECT /*/*/*; $x$;*/*/*/ FROM t WHERE a = --END ;
|  FROM t WHERE a = CREATE FUNCTION ltei() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = "/*;/*;" + (SELECT 1; SELECT 2) FROM t WHERE a = "end;select 1;a;b" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "/*;/*;" + (SELECT 1; SELECT 2) FROM t WHERE a = "end;select 1;a;b" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "/*;/*;" + (SELECT 1; SELECT 2) FROM t WHERE a = "end;select 1;a;b" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY xem FROM stdin;
| 1	\N
| 1	\N
| \.
|
| SELECT --; --;
|  col U&'--;end;' 1 X'01;10' AS x; SELECT '--;\' col gl$z$y +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xem FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xem FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY bpl FROM stdin;
| a;b
| a;b
| \.
|
| COPY oje FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bpl FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bpl FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + fa$q$c col CREATE FUNCTION alkz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(select 1;'; SELECT 1; END 1;col X'01;10' SELECT ';select 1;END ;' FROM t WHERE a = '$;$;' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + fa$q$c col CREATE FUNCTION alkz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(select 1;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + fa$q$c col CREATE FUNCTION alkz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(select 1;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = _b$e$s AS x CREATE FUNCTION flfw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;--;;'; END + 'a;b' SELECT;
| COPY qid FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = _b$e$s AS x CREATE FUNCTION flfw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;--;;'; END + 'a;b' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = _b$e$s AS x CREATE FUNCTION flfw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;--;;'; END + 'a;b' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , (SELECT 1; SELECT 2) col U&'/*;a;b' SELECT;
| , B'01;10' , ");(" 1 (SELECT 1; SELECT 2) AS x;
| , U&'END ;' FROM t WHERE a =;
| COPY hlk FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) col U&'/*;a;b' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) col U&'/*;a;b' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT U&';;' 1 X'01;10' AS x /*/*/*END ;/ *;*/*/*/ SELECT;+ $$$;@;$$$ SELECT po$f$c +;
| col U&'$x$;;END ;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&';;' 1 X'01;10' AS x /*/*/*END ;/ *;*/*/*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&';;' 1 X'01;10' AS x /*/*/*END ;/ *;*/*/*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| col /*$;$*/ , E'\n;' AS x;
| + --;;select 1;
|  ,;AS x e'\n;' SELECT CREATE FUNCTION gouj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;;;'; END col (SELECT 1; SELECT 2) +; AS x ';;\' + /*/*$x$;*/*/ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*$;$*/ , E'\n;' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*$;$*/ , E'\n;' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = E'\';' AS x "$;$select 1;" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E'\';' AS x "$;$select 1;" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\';' AS x "$;$select 1;" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT --$x$;
|  SELECT;
| SELECT CREATE FUNCTION hden() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;end;;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "SELECT"
| position: 16
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --$x$;
|  SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 16
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --$x$;
|  SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 16
> exec
| FROM t WHERE a = E'\\'';\';' + CREATE FUNCTION acpt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END , "$;$;;" SELECT;
| COPY usd FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E'\\'';\';' + CREATE FUNCTION acpt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END , "$;$;;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\\'';\';' + CREATE FUNCTION acpt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END , "$;$;;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY yva FROM stdin;
| 1	\N
| a;b
| \.
|
| AS x ');(\' AS x E'\n;'';' +; -- trailing note;
| COPY equ FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yva FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yva FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 $$/*;$$ 1 $bfnr$a;b$bfnr$ +; -- trailing note;
| SELECT oi$u$h 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
| COPY uhs FROM stdin;
| a;b
| 1	\N
| \.
|
| FROM t WHERE a = e'\n;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $$/*;$$ 1 $bfnr$a;b$bfnr$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $$/*;$$ 1 $bfnr$a;b$bfnr$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = U&';' + b'01;10' + X'01;10' AS x;
| AS x /*/*a;b;;*/*/ AS x "end;; --;" 1 e'\\\\' +; COPY oj_ FROM stdin;
| c;d;e
| c;d;e
| \.
|  AS x 'END ;select 1;; \' 1 hz$x$f 1 yr$n$r 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&';' + b'01;10' + X'01;10' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&';' + b'01;10' + X'01;10' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + e'\\'';' FROM t WHERE a = B'01;10' 1; COPY yon FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + e'\\'';' FROM t WHERE a = B'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + e'\\'';' FROM t WHERE a = B'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + 'select 1;; ' col U&');(;;--;' FROM t WHERE a = --$;$; 
|  ,;
| FROM t WHERE a = U&'a;b;;;' +; COPY yyt FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + 'select 1;; ' col U&');(;;--;' FROM t WHERE a = --$;$; 
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + 'select 1;; ' col U&');(;;--;' FROM t WHERE a = --$;$; 
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT CREATE FUNCTION ciq_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(a;b;'; SELECT 1; END AS x /*/*/*a;bEND ;*/*/*/ AS x ft$h$x AS x; COPY hzf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION ciq_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(a;b;'; SELECT 1; END AS x /*/*/*a;bEND ;*/*/*/ AS x ft$h$x AS x;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION ciq_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(a;b;'; SELECT 1; END AS x /*/*/*a;bEND ;*/*/*/ AS x ft$h$x AS x;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| COPY nbn FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY nbn FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY nbn FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 "select 1;END ;" ,; COPY wux FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "select 1;END ;" ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "select 1;END ;" ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + $$/*;);($$ , zd$l$d col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $$/*;);($$ , zd$l$d col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $$/*;);($$ , zd$l$d col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = ");($x$;" col; COPY hcc FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ");($x$;" col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ");($x$;" col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY dsm FROM stdin;
| 1	\N
| a;b
| \.
|
| 1 $_$$x$;$x$$_$ AS x (SELECT 1; SELECT 2) col CREATE FUNCTION emjt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END AS x;
| FROM t WHERE a = $$$x$;end;select 1;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dsm FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dsm FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY pvr FROM stdin;
| 1	\N
| c;d;e
| \.
|  , $rnl$);($rnl$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pvr FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pvr FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x "END ;END ;" ,;
| COPY m_o FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "END ;END ;" ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "END ;END ;" ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col X'01;10' SELECT;
| FROM t WHERE a = ';;\' col /*);(--;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col X'01;10' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col X'01;10' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + '$;$$x$;;;' 1 "--;$x$;" SELECT;
| col cx$p$e , "$x$;;;$x$;" +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + '$;$$x$;;;' 1 "--;$x$;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + '$;$$x$;;;' 1 "--;$x$;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY ngx FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ngx FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ngx FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = B'01;10' col;
| FROM t WHERE a = --;;
|  col; -- trailing note;
| AS x /*/*/*END ;a;b*/*/*/ col E'\\' AS x '$x$;' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = B'01;10' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = B'01;10' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY i_m FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| col U&'$;$);($x$;' FROM t WHERE a = E'\';\';\';' AS x;
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =; -- trailing note;
| COPY tyf FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY i_m FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY i_m FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY hsp FROM stdin;
| a;b
| c;d;e
| \.
|
| + b'01;10' 1 ");(" SELECT e'\n;\n;\n;' col;
| COPY fkb FROM stdin;
| 1	\N
| c;d;e
| \.
|
| , CREATE FUNCTION eiaq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; END col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hsp FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hsp FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col --;;end;end;
|  SELECT '--;' 1 U&'$;$' FROM t WHERE a =;COPY yc_ FROM stdin;
| a;b
| a;b
| \.
|
| + /*end;$x$;*/ SELECT 'select 1;$x$;);(' 1 /*/*/*select 1;;/ *;*/*/*/ SELECT;COPY nbj FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --;;end;end;
|  SELECT '--;' 1 U&'$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --;;end;end;
|  SELECT '--;' 1 U&'$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) AS x b'01;10' + "; );(" col;
| COPY tvl FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) AS x b'01;10' + "; );(" col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) AS x b'01;10' + "; );(" col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| col /*a;b; ;;*/ SELECT ut$u$q FROM t WHERE a = X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*a;b; ;;*/ SELECT ut$u$q FROM t WHERE a = X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*a;b; ;;*/ SELECT ut$u$q FROM t WHERE a = X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 $u$END ;;;select 1;$u$ AS x tv$t$x AS x "a;b" 1;
| FROM t WHERE a = ');(;\' FROM t WHERE a = CREATE FUNCTION vkgi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;;'; SELECT 1; END ,;
| + $qq$--;$x$;);($qx$$qq$ col (SELECT 1; SELECT 2) + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $u$END ;;;select 1;$u$ AS x tv$t$x AS x "a;b" 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $u$END ;;;select 1;$u$ AS x tv$t$x AS x "a;b" 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY pgt FROM stdin;
| a;b
| a;b
| \.
|
| + /*$x$;*/ +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pgt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pgt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col $u$$;$end;;;$x$$u$ 1 cy$_$j + U&'/*;END ;' +; -- trailing note;
| COPY lva FROM stdin;
| a;b
| c;d;e
| \.
| , U&';;);($;$' AS x E'\n;\n;' , CREATE FUNCTION uflq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;'; SELECT 1; END col;
| , --end;;
|  , /*/*; ; */*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $u$$;$end;;;$x$$u$ 1 cy$_$j + U&'/*;END ;' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $u$$;$end;;;$x$$u$ 1 cy$_$j + U&'/*;END ;' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY mwv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mwv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mwv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY ojy FROM stdin;
| 1	\N
| 1	\N
| \.
|
| + CREATE FUNCTION bnnj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;; '; END FROM t WHERE a = ";" ,;
| 1 ");(" col yp$n$o SELECT (SELECT 1; SELECT 2) AS x;COPY ysh FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ojy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ojy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + U&'END ;);(' FROM t WHERE a =; -- trailing note;
| 1 e''';'';' 1 qs$g$o + U&'--;select 1;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'END ;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'END ;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + (SELECT 1; SELECT 2) 1; AS x --END ;
|  FROM t WHERE a = --$x$;end;
|  AS x; AS x '/*;);(\' SELECT; -- trailing note;
| FROM t WHERE a = --;;/*;
|  col /*$x$;select 1;$x$;*/ , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY qxu FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qxu FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qxu FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 /*/*end;;;;;*/*/ SELECT $$$;$);($$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*end;;;;;*/*/ SELECT $$$;$);($$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*end;;;;;*/*/ SELECT $$$;$);($$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col U&'; ' AS x "/*;;);(" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'; ' AS x "/*;;);(" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'; ' AS x "/*;;);(" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x /*; a;b);(*/ ,; AS x $w$$;$$x$$w$ col X'01;10' ,; AS x fn$r$l + (SELECT 1; SELECT 2) SELECT;FROM t WHERE a = 'END ;);(end;' , --$;$);();(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*; a;b);(*/ ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*; a;b);(*/ ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , $h$--;$h$ FROM t WHERE a =;
| col (SELECT 1; SELECT 2) + --; ; 
|  SELECT CREATE FUNCTION gbrd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$'; SELECT 1; END 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $h$--;$h$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $h$--;$h$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col CREATE FUNCTION drcf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';--;select 1;'; END AS x 'a;bend;\' ,;+ /*/*/*;;*/*/*/ 1;
| AS x $$/*;; $$ AS x CREATE FUNCTION _oza() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END +;
| 1 x'01;10' SELECT --a;b
|  col ');(' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION drcf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';--;select 1;'; END AS x 'a;bend;\' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION drcf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';--;select 1;'; END AS x 'a;bend;\' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = qs$s$k + --;
|  + CREATE FUNCTION axci() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = qs$s$k + --;
|  + CREATE FUNCTION axci() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = qs$s$k + --;
|  + CREATE FUNCTION axci() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , ");(" FROM t WHERE a = /*/*END ;*/*/ 1 $m$end;$m$ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ");(" FROM t WHERE a = /*/*END ;*/*/ 1 $m$end;$m$ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ");(" FROM t WHERE a = /*/*END ;*/*/ 1 $m$end;$m$ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col ');(select 1;\' col;
| 1 ");($;$" col rl$b$v ,;
| SELECT e'\\\';\n;' AS x U&';;' , $$;;);(select 1;$$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ');(select 1;\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ');(select 1;\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 /*/*/ *;;*/*/ AS x;
| FROM t WHERE a = $$; END ;$$ 1 "$x$;end;" +;
| SELECT e'\\'';' +; -- trailing note;
| + U&'a;bselect 1;; ' , E'\';\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/ *;;*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/ *;;*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col U&'END ;END ;' 1 /*;;*/ +;
| col --a;b
|  + U&'/*;;;' AS x; + ws$e$q SELECT CREATE FUNCTION wnus() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$--;--;'; SELECT 1; END + --;
|  SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'END ;END ;' 1 /*;;*/ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'END ;END ;' 1 /*;;*/ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 ';' + B'01;10' FROM t WHERE a =;
| AS x (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ';' + B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ';' + B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x /*/*end;*/*/ SELECT;
| + CREATE FUNCTION tnxu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;END ;'; SELECT 1; END 1 ';;;);(' + 'end;select 1;' ,;
| col "--;$;$" + $$$x$;/*;$$ ,;
| AS x uy$o$n + U&'end;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*end;*/*/ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*end;*/*/ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + E''';\n;\n;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E''';\n;\n;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E''';\n;\n;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT U&'--;' AS x ';$x$;END ;\' AS x --;;$;$;
|  1;
| SELECT "$x$;" FROM t WHERE a = (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "';$x$;END ;\\'"
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'--;' AS x ';$x$;END ;\' AS x --;;$;$;
|  1;
< error
| 42601 syntax error at or near "';$x$;END ;\\'"
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'--;' AS x ';$x$;END ;\' AS x --;;$;$;
|  1;
< error
| 42601 syntax error at or near "';$x$;END ;\\'"
| position: 21
> exec
| 1 "; ; select 1;" + U&'--;;;$x$;' + (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "; ; select 1;" + U&'--;;;$x$;' + (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "; ; select 1;" + U&'--;;;$x$;' + (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + x'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + x'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + x'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY tjf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY tjf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY tjf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col CREATE FUNCTION _nza() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;select 1;'; END 1 CREATE FUNCTION vmqp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END FROM t WHERE a = --end;select 1;--;
|  ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION _nza() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;select 1;'; END 1 CREATE FUNCTION vmqp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END FROM t WHERE a = --end;select 1;--;
|  ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION _nza() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;select 1;'; END 1 CREATE FUNCTION vmqp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END FROM t WHERE a = --end;select 1;--;
|  ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY ovb FROM stdin;
| a;b
| 1	\N
| \.
|
| COPY tmy FROM stdin;
| a;b
| c;d;e
| \.
|
| SELECT /*/*/*END ;$x$;$;$*/*/*/ , E'\n;'';' FROM t WHERE a = E'\';\\\\' col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ovb FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ovb FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x CREATE FUNCTION vedl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END 1;
| SELECT cq$f$g FROM t WHERE a = --;;);(
|  ,; + $sau$);(;$sau$ FROM t WHERE a = ff$f$w ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION vedl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION vedl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col "/*;" AS x;
| 1 E'\n;\n;\\' AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "/*;" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "/*;" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) SELECT '; ' + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) SELECT '; ' + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) SELECT '; ' + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY hli FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hli FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hli FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x U&'$;$END ;' 1 (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'$;$END ;' 1 (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'$;$END ;' 1 (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col U&'; ' 1 '; \' 1;
| COPY ckt FROM stdin;
| a;b
| c;d;e
| \.
|  , E'\\'';'';' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'; ' 1 '; \' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'; ' 1 '; \' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + --/*;
|  + e''';\n;\';' SELECT;
| COPY _wk FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --/*;
|  + e''';\n;\';' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --/*;
|  + e''';\n;\';' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = (SELECT 1; SELECT 2) col;
| col e'\';\\' +;
| COPY pzi FROM stdin;
| c;d;e
| c;d;e
| \.
|
| col --end;
|  AS x;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| COPY qia FROM stdin;
| 1	\N
| a;b
| \.
|  -- trailing note;
| AS x U&'a;b/*;;;' col CREATE FUNCTION ryiq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;end;'; SELECT 1; END 1 U&'a;b' FROM t WHERE a =;
| AS x /*; $x$;;;*/ 1 "$;$" ,;
| COPY yva FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qia FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qia FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = nt$b$b col sw$k$f col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = nt$b$b col sw$k$f col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = nt$b$b col sw$k$f col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY org FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY org FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY org FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , b'01;10' + CREATE FUNCTION qrgi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;$x$;select 1;'; SELECT 1; END AS x b'01;10' AS x; 1 $$END ;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , b'01;10' + CREATE FUNCTION qrgi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;$x$;select 1;'; SELECT 1; END AS x b'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , b'01;10' + CREATE FUNCTION qrgi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;$x$;select 1;'; SELECT 1; END AS x b'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY opm FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY opm FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY opm FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col rl$l$c 1;
| COPY pcs FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| + CREATE FUNCTION v_qz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END + vk$a$z +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col rl$l$c 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col rl$l$c 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY yly FROM stdin;
| c;d;e
| a;b
| \.
|  -- trailing note;
| col E'\n;\';'';' AS x B'01;10' FROM t WHERE a =;
| AS x /*/*/*;;END ;*/*/*/ ,;+ (SELECT 1; SELECT 2) AS x $$$x$;; $$ col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yly FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yly FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , U&'; --;select 1;' , $$;$x$;$$ AS x e'\n;\';' +;
| FROM t WHERE a = e'\\\';' , X'01;10' col ");(--;--;" 1;
| FROM t WHERE a = $f$$x$;;;$f$ FROM t WHERE a = 'select 1;--;);(' col;
| col CREATE FUNCTION vsni() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END col --/*;select 1;/*;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'; --;select 1;' , $$;$x$;$$ AS x e'\n;\';' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'; --;select 1;' , $$;$x$;$$ AS x e'\n;\';' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , e'\';\n;' col U&'; ;;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\';\n;' col U&'; ;;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\';\n;' col U&'; ;;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + $$select 1;--;$$ 1 CREATE FUNCTION zqow() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; '; END , $vm$$x$;$;$/*;$vx$$vm$ AS x;
| col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $$select 1;--;$$ 1 CREATE FUNCTION zqow() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; '; END , $vm$$x$;$;$/*;$vx$$vm$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $$select 1;--;$$ 1 CREATE FUNCTION zqow() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; '; END , $vm$$x$;$;$/*;$vx$$vm$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x CREATE FUNCTION hsgy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END FROM t WHERE a = U&';select 1;' 1 /*/*/*/ *;$;$*/*/*/ FROM t WHERE a =; SELECT "$x$;$x$;" ,; -- trailing note;
| AS x U&'; --;);(' FROM t WHERE a = "end;select 1;" + e'\\\\\';' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION hsgy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END FROM t WHERE a = U&';select 1;' 1 /*/*/*/ *;$;$*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION hsgy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END FROM t WHERE a = U&';select 1;' 1 /*/*/*/ *;$;$*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = $$; end;--;$$ , /*a;b;*/ +;
| COPY c_r FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$; end;--;$$ , /*a;b;*/ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$; end;--;$$ , /*a;b;*/ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = "--;$;$" AS x;
| COPY fwc FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "--;$;$" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "--;$;$" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY fpj FROM stdin;
| a;b
| c;d;e
| \.
| COPY xge FROM stdin;
| a;b
| c;d;e
| \.
|
| col U&'$;$;;' 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fpj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fpj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , /*/*$x$;a;b;*/*/ col;
| 1 ');($x$;$x$;' +;
| col E'\';\';' AS x 'END ;);(' +;
| 1 CREATE FUNCTION xhr_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$x$;/*;'; SELECT 1; END col --; 
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*$x$;a;b;*/*/ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*$x$;a;b;*/*/ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = CREATE FUNCTION mnxl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;);(/*;'; SELECT 1; END +;, e'\n;' col '/*;$x$;' 1 sk$q$x +;
| col --END ;
|  FROM t WHERE a = df$f$t AS x vt$b$t ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION mnxl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;);(/*;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION mnxl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;);(/*;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col ----;; 
|  col; , b'01;10' AS x CREATE FUNCTION inm_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;a;b'; END AS x e''';'';' col;
| 1 E'\n;\n;' + U&'--;$x$;' + "$x$;end;" ,;FROM t WHERE a = 'select 1;a;bend;' 1 b'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ----;; 
|  col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ----;; 
|  col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + E'\n;\\'';' SELECT /*; select 1;);(*/ ,;
| COPY jyh FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\n;\\'';' SELECT /*; select 1;);(*/ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\n;\\'';' SELECT /*; select 1;);(*/ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = e'\';' +;1 /*select 1;--;*/ , fx$l$b SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = e'\';' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = e'\';' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , om$_$h ,; -- trailing note;
| , "$x$;" SELECT;
| COPY guf FROM stdin;
| c;d;e
| 1	\N
| \.
|
| FROM t WHERE a = CREATE FUNCTION xfqg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);('; SELECT 1; END , "END ;;" col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , om$_$h ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , om$_$h ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) col; COPY fpz FROM stdin;
| a;b
| c;d;e
| \.
|  , U&'END ;' FROM t WHERE a = e'\\'';\n;' + /*/*;;*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col xy$u$k AS x "$;$" AS x;
| , b'01;10' col no$j$l ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col xy$u$k AS x "$;$" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col xy$u$k AS x "$;$" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col "a;b/*;;;" FROM t WHERE a = '$;$;;select 1;' ,;
| AS x (SELECT 1; SELECT 2) + CREATE FUNCTION ienz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;; ;;'; END FROM t WHERE a = '$;$\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "a;b/*;;;" FROM t WHERE a = '$;$;;select 1;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "a;b/*;;;" FROM t WHERE a = '$;$;;select 1;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY sop FROM stdin;
| a;b
| c;d;e
| \.
|
| SELECT --a;b
|  ,; -- trailing note;
| col X'01;10' ,;
| , /*/*/*end;*/*/*/ + B'01;10' AS x qc$o$x col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sop FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sop FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 /*/*;*/*/ , 'end;--;END ;\' +;col CREATE FUNCTION iveh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;; ;;'; END 1 U&';;' col;
| 1 (SELECT 1; SELECT 2) FROM t WHERE a = (SELECT 1; SELECT 2) + --END ;
|  1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*;*/*/ , 'end;--;END ;\' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*;*/*/ , 'end;--;END ;\' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY qto FROM stdin;
| a;b
| 1	\N
| \.
|
| AS x x'01;10' SELECT CREATE FUNCTION _ypt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END ,;col CREATE FUNCTION z_my() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;$x$;'; SELECT 1; END AS x $zm$/*;select 1;select 1;$zm$ FROM t WHERE a =;
| 1 (SELECT 1; SELECT 2) + (SELECT 1; SELECT 2) SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qto FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qto FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 wb$p$q SELECT /*/*;; ; */*/ , $$);();(--;$$ ,;
| 1 $$end;;;$$ 1;
| FROM t WHERE a = /*/*/*a;b*/*/*/ FROM t WHERE a = CREATE FUNCTION grix() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; END col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 wb$p$q SELECT /*/*;; ; */*/ , $$);();(--;$$ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 wb$p$q SELECT /*/*;; ; */*/ , $$);();(--;$$ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY upq FROM stdin;
| 1	\N
| a;b
| \.
|
| + /*--;*/ SELECT;
| FROM t WHERE a = $djfc$select 1;/*;);($djfx$$djfc$ col /*a;b*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY upq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY upq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY gez FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gez FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gez FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + CREATE FUNCTION plwx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;END ;'; END SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION plwx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;END ;'; END SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION plwx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;END ;'; END SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT X'01;10' 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT X'01;10' 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT X'01;10' 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| FROM t WHERE a = 'end;' 1 x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = 'end;' 1 x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = 'end;' 1 x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x /*/*/*a;bselect 1;*/*/*/ AS x;
| COPY gvw FROM stdin;
| 1	\N
| 1	\N
| \.
| 1 b'01;10' +;
| col ');($x$;' SELECT ";;$x$;; " ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*a;bselect 1;*/*/*/ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*a;bselect 1;*/*/*/ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT U&'; ' SELECT e'\';' SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "e'\\';'"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'; ' SELECT e'\';' SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "e'\\';'"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'; ' SELECT e'\';' SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "e'\\';'"
| position: 22
> exec
| 1 X'01;10' col;
| COPY jff FROM stdin;
| c;d;e
| a;b
| \.
| SELECT B'01;10' col;
| SELECT X'01;10' SELECT 'a;b/*;/*;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 X'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 X'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) , /*/*/*$;$*/*/*/ 1 $$;;;$$ FROM t WHERE a =; SELECT x'01;10' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) , /*/*/*$;$*/*/*/ 1 $$;;;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) , /*/*/*$;$*/*/*/ 1 $$;;;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY pyl FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 --a;b; ;;
|  1 $$;;a;b; $$ + e'\\\\'';' ,; COPY def FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pyl FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pyl FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = ";;; ; " SELECT /*/*; / *;*/*/ col "END ;;" 1; COPY yzv FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ";;; ; " SELECT /*/*; / *;*/*/ col "END ;;" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ";;; ; " SELECT /*/*; / *;*/*/ col "END ;;" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x vb$q$q 1 (SELECT 1; SELECT 2) 1 CREATE FUNCTION cqrj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bselect 1;'; SELECT 1; END FROM t WHERE a =; SELECT (SELECT 1; SELECT 2) +;
| COPY qew FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x vb$q$q 1 (SELECT 1; SELECT 2) 1 CREATE FUNCTION cqrj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bselect 1;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x vb$q$q 1 (SELECT 1; SELECT 2) 1 CREATE FUNCTION cqrj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bselect 1;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x '$x$;' col --a;b/*;
|  ,; -- trailing note;
| COPY x__ FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x '$x$;' col --a;b/*;
|  ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x '$x$;' col --a;b/*;
|  ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 e'\n;\\\';' FROM t WHERE a = /*END ;*/ + '/*;\' FROM t WHERE a =;
| FROM t WHERE a = E'\\\';\n;' SELECT 'select 1;end;);(' + /*/*/*a;b);(END ;*/*/*/ col;FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a = $yj$--;$yx$$yj$ SELECT;
| 1 U&'; select 1;' FROM t WHERE a = 'a;b;;a;b' AS x U&'$x$;select 1;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\n;\\\';' FROM t WHERE a = /*END ;*/ + '/*;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\n;\\\';' FROM t WHERE a = /*END ;*/ + '/*;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT $e$;;;$x$$e$ AS x;
| + CREATE FUNCTION qzmx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$$x$;'; END FROM t WHERE a = b'01;10' col 'END ;\' ,;
| col sb$t$f AS x $iav_$/*;;select 1;$iavx$$iav_$ AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 27
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $e$;;;$x$$e$ AS x;
< ok 1 -
> exec
|
| + CREATE FUNCTION qzmx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$$x$;'; END FROM t WHERE a = b'01;10' col 'END ;\' ,;
< error
| 42601 syntax error at or near "+"
| position: 2
> exec
| ROLLBACK
< ok 0 -
> exec
|
| + CREATE FUNCTION qzmx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$$x$;'; END FROM t WHERE a = b'01;10' col 'END ;\' ,;
< error
| 42601 syntax error at or near "+"
| position: 2
> exec
| FROM t WHERE a = ve$v$l col U&'/*;' 1 "; " col; 1 pn$_$n AS x U&'end;a;b' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ve$v$l col U&'/*;' 1 "; " col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ve$v$l col U&'/*;' 1 "; " col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY wkz FROM stdin;
| c;d;e
| c;d;e
| \.
|
| SELECT "select 1;" FROM t WHERE a =;, CREATE FUNCTION mvqn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END col $q$$x$;;$q$ AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wkz FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wkz FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x CREATE FUNCTION wobr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END , ve$b$c + "$x$;" +;
| FROM t WHERE a = --end;--;END ;
|  col $g$;$g$ ,;
| col X'01;10' FROM t WHERE a =;
| , $$END ;$$ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION wobr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END , ve$b$c + "$x$;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION wobr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END , ve$b$c + "$x$;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 b'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 b'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 b'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , B'01;10' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , B'01;10' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , B'01;10' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 ");(; " + $$; $$ AS x /*/*$x$;);(*/*/ AS x;1 (SELECT 1; SELECT 2) AS x $opgc$$x$;/*;$;$$opgx$$opgc$ 1 '$;$; end;' col;
| + X'01;10' + /*/*end;/ *;*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ");(; " + $$; $$ AS x /*/*$x$;);(*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ");(; " + $$; $$ AS x /*/*$x$;);(*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT gm$t$z +;
| + B'01;10' FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| BEGIN
< ok 0 -
> exec
| SELECT gm$t$z +;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT gm$t$z +;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| 1 U&'select 1;;$;$' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'select 1;;$;$' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'select 1;;$;$' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY ork FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| + --$;$
|  AS x;
| SELECT '; ' AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ork FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ork FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT ";--;" AS x;
| COPY stw FROM stdin;
| c;d;e
| c;d;e
| \.
|
| , /*/*$x$;$x$;;*/*/ , tn$w$g , 'select 1;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 42
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ";--;" AS x;
< error
| 42703 column ";--;" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| COPY stw FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 23
> exec
| COPY blo FROM stdin;
| 1	\N
| 1	\N
| \.
|
| + (SELECT 1; SELECT 2) + (SELECT 1; SELECT 2) 1 --$x$;
|  1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY blo FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY blo FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT (SELECT 1; SELECT 2) col U&');(select 1;' FROM t WHERE a =;
| SELECT --$x$;end;;
|  AS x;
| COPY wie FROM stdin;
| c;d;e
| c;d;e
| \.
|
| SELECT ki$b$j + --END ;;;; 
|  1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col U&');(select 1;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col U&');(select 1;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| , /*/*end;*/*/ , "; ; " 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*end;*/*/ , "; ; " 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*end;*/*/ , "; ; " 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 CREATE FUNCTION tcjy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;--;);('; SELECT 1; END 1 "; $x$;a;b" , /*/*/*$x$;END ;;;*/*/*/ col; -- trailing note;
| col U&');(END ;;;' AS x;
| FROM t WHERE a = E'\n;\';\';' , 'end;end;' AS x --end;
|  ,;
| 1 'end;$;$;;\' + U&'/*;$;$' AS x "$;$" col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION tcjy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;--;);('; SELECT 1; END 1 "; $x$;a;b" , /*/*/*$x$;END ;;;*/*/*/ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION tcjy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;--;);('; SELECT 1; END 1 "; $x$;a;b" , /*/*/*$x$;END ;;;*/*/*/ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT --END ;
|  FROM t WHERE a = '--;a;b$;$' FROM t WHERE a = E'\n;\\\\' ,;
| COPY hyw FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 46
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --END ;
|  FROM t WHERE a = '--;a;b$;$' FROM t WHERE a = E'\n;\\\\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 46
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --END ;
|  FROM t WHERE a = '--;a;b$;$' FROM t WHERE a = E'\n;\\\\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 46
> exec
| col '$x$;--;' 1;
| COPY put FROM stdin;
| c;d;e
| 1	\N
| \.
|
| + B'01;10' FROM t WHERE a = 'a;b);(' FROM t WHERE a = X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '$x$;--;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '$x$;--;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY m_b FROM stdin;
| 1	\N
| c;d;e
| \.
|  COPY dlk FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| COPY gsy FROM stdin;
| c;d;e
| c;d;e
| \.
|
| FROM t WHERE a = $fj_p$$;$/*;$fj_x$$fj_p$ ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY m_b FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY m_b FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY puc FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY puc FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY puc FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY kcn FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| COPY ezc FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kcn FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kcn FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY jxy FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY xwn FROM stdin;
| 1	\N
| a;b
| \.
|
| AS x b'01;10' AS x; col ';;END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jxy FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jxy FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 --a;b;$;$
|  , U&'a;b' SELECT da$_$m AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --a;b;$;$
|  , U&'a;b' SELECT da$_$m AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --a;b;$;$
|  , U&'a;b' SELECT da$_$m AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + CREATE FUNCTION ilij() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;; '; SELECT 1; END col;
| COPY jtr FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION ilij() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;; '; SELECT 1; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION ilij() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;; '; SELECT 1; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 CREATE FUNCTION czvw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;--;'; END col (SELECT 1; SELECT 2) ,;
| + $nru$select 1;$nru$ 1; -- trailing note;
| SELECT (SELECT 1; SELECT 2) AS x E'\\' FROM t WHERE a = "end;" 1;
| AS x E'\n;' col '; select 1;select 1;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION czvw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;--;'; END col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION czvw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;--;'; END col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , E'\n;'';\n;' , E'\';' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , E'\n;'';\n;' , E'\';' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , E'\n;'';\n;' , E'\';' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY yxp FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yxp FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yxp FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col "END ;select 1;select 1;" AS x; -- trailing note;
| SELECT "select 1;--;a;b" FROM t WHERE a = '$;$a;ba;b' 1 CREATE FUNCTION idpj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END col;
| 1 '$x$;\' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "END ;select 1;select 1;" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "END ;select 1;select 1;" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , U&'select 1;a;b$x$;' 1;
| COPY w_e FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'select 1;a;b$x$;' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'select 1;a;b$x$;' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col ----;end;
|  +;FROM t WHERE a = '$;$\' 1 e'\';\\\';' ,;
| AS x hn$u$f col s_$_$r , "END ;select 1;select 1;" +; , (SELECT 1; SELECT 2) col --);(;;; 
|  +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ----;end;
|  +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ----;end;
|  +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , e'\';\n;' col;
| , b'01;10' +;
| , ');(/*;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\';\n;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\';\n;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col E'\\'';'';' AS x /*/*/*;select 1;$x$;*/*/*/ SELECT;
| SELECT $$; select 1;$$ col ra$t$t , U&'/*;--;END ;' ,; -- trailing note;
| AS x --END ;$x$;a;b
|  SELECT $$$x$;/*;select 1;$$ AS x /*/*;;;;; */*/ col;
| 1 /*/*/*END ;select 1;--;*/*/*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col E'\\'';'';' AS x /*/*/*;select 1;$x$;*/*/*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col E'\\'';'';' AS x /*/*/*;select 1;$x$;*/*/*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col $ic$end;a;b$ic$ + '$;$/*;' AS x; -- trailing note;
| SELECT --/*;a;b
|  AS x --END ;;;
|  col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $ic$end;a;b$ic$ + '$;$/*;' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $ic$end;a;b$ic$ + '$;$/*;' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = --;/*;
|  FROM t WHERE a = "a;b; a;b" +; 1 U&'$;$' 1 e'\n;' SELECT /*/*/ *;a;b);(*/*/ SELECT; + $$end;$$ col --$x$;;;
|  , B'01;10' ,; FROM t WHERE a = $$; ;$$ 1 (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --;/*;
|  FROM t WHERE a = "a;b; a;b" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --;/*;
|  FROM t WHERE a = "a;b; a;b" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 CREATE FUNCTION xzsz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;; ;;'; END 1 b'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION xzsz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;; ;;'; END 1 b'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION xzsz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;; ;;'; END 1 b'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT $p$END ;$p$ AS x; , (SELECT 1; SELECT 2) col U&'end;a;b' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 26
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $p$END ;$p$ AS x;
< ok 1 -
> exec
|  , (SELECT 1; SELECT 2) col U&'end;a;b' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 2
> exec
| ROLLBACK
< ok 0 -
> exec
|  , (SELECT 1; SELECT 2) col U&'end;a;b' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 2
> exec
| col '--;' SELECT;1 on$b$m AS x E'\\\\\n;' 1; SELECT $d$--;$x$$d$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '--;' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '--;' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x B'01;10' col U&'a;b' +;
| + E''';\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x B'01;10' col U&'a;b' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x B'01;10' col U&'a;b' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT ';$x$;' , $k$;a;b; $x$$k$ SELECT;
| COPY kzy FROM stdin;
| 1	\N
| 1	\N
| \.
|
| SELECT CREATE FUNCTION cigj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END + /*/*/*$x$;*/*/*/ + CREATE FUNCTION ijuv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; '; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 63
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ';$x$;' , $k$;a;b; $x$$k$ SELECT;
< ok 1 -
> exec
|
| COPY kzy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 23
> exec
| ROLLBACK
< ok 0 -
> exec
|
| COPY kzy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 23
> exec
| FROM t WHERE a = $$$x$;$$ + /*/*$x$;*/*/ + U&'--;/*;end;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$$x$;$$ + /*/*$x$;*/*/ + U&'--;/*;end;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$$x$;$$ + /*/*$x$;*/*/ + U&'--;/*;end;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col /*/*/*;;; */*/*/ + "$x$;a;b$x$;" AS x x'01;10' 1; , 'end;--;;\' 1 fd$g$q col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*/*;;; */*/*/ + "$x$;a;b$x$;" AS x x'01;10' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*/*;;; */*/*/ + "$x$;a;b$x$;" AS x x'01;10' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x $$;;/*;$$ SELECT E'\\\\' FROM t WHERE a = U&'END ;' FROM t WHERE a =;
| 1 e''';\';\\' FROM t WHERE a = (SELECT 1; SELECT 2) col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $$;;/*;$$ SELECT E'\\\\' FROM t WHERE a = U&'END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $$;;/*;$$ SELECT E'\\\\' FROM t WHERE a = U&'END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 $$a;b$$ , 'end;' FROM t WHERE a =;
| FROM t WHERE a = $$$;$--;$$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $$a;b$$ , 'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $$a;b$$ , 'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY cvk FROM stdin;
| 1	\N
| c;d;e
| \.
|
| SELECT /*select 1;end;*/ col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cvk FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cvk FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 $$;; ;$$ col;
| + "end;" + E'\';' , --$;$
|  AS x;AS x b'01;10' 1 e''';\';' 1;
| COPY tbf FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $$;; ;$$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $$;; ;$$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT b'01;10' , (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) FROM t WHERE a =; COPY mg_ FROM stdin;
| 1	\N
| a;b
| \.
|  -- trailing note;
| COPY rbi FROM stdin;
| c;d;e
| 1	\N
| \.
|  + xg$q$d 1 /*a;b*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 28
> exec
| BEGIN
< ok 0 -
> exec
| SELECT b'01;10' , (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 28
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT b'01;10' , (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 28
> exec
| SELECT (SELECT 1; SELECT 2) 1;
| 1 $$end;$x$;$$ col --);(
|  1 e'\';' FROM t WHERE a =; SELECT U&'$;$select 1;' AS x U&'$;$;;; ' ,;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| col 'a;b);(/*;\' + x'01;10' col;
| AS x (SELECT 1; SELECT 2) + "/*;--;end;" 1;
| SELECT b'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'a;b);(/*;\' + x'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'a;b);(/*;\' + x'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY wee FROM stdin;
| a;b
| 1	\N
| \.
|
| FROM t WHERE a = CREATE FUNCTION owak() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;END ;'; SELECT 1; END col ';;/*;\' 1 U&'a;b;;' AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wee FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wee FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT E'\';' SELECT;
< ok 1 -
> exec
| BEGIN
< ok 0 -
> exec
| SELECT E'\';' SELECT;
< ok 1 -
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$;; ;;$$ 1 E'\n;\';' SELECT $$;;$x$;$$ 1;
| + 'select 1;END ;end;\' AS x;
< error
| 42601 syntax error at or near "1"
| position: 18
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$;; ;;$$ 1 E'\n;\';' SELECT $$;;$x$;$$ 1;
< error
| 42601 syntax error at or near "1"
| position: 18
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$;; ;;$$ 1 E'\n;\';' SELECT $$;;$x$;$$ 1;
< error
| 42601 syntax error at or near "1"
| position: 18
> exec
| + U&';' AS x fa$o$f col CREATE FUNCTION paem() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END +; -- trailing note;
| AS x "/*;/*;" FROM t WHERE a =;
| COPY xbc FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&';' AS x fa$o$f col CREATE FUNCTION paem() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&';' AS x fa$o$f col CREATE FUNCTION paem() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY zgz FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zgz FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zgz FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col "$x$;" col CREATE FUNCTION apvs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; END col 'a;b\' AS x; 1 "END ;--;" ,;
| AS x x'01;10' AS x yj$p$u SELECT "--;;" ,; col (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "$x$;" col CREATE FUNCTION apvs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; END col 'a;b\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "$x$;" col CREATE FUNCTION apvs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; END col 'a;b\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + --;a;bEND ;
|  + (SELECT 1; SELECT 2) , xi$b$o AS x;
| 1 '--;/*;\' col CREATE FUNCTION pran() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;end;--;'; SELECT 1; END +;
| FROM t WHERE a = '; ;\' SELECT (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --;a;bEND ;
|  + (SELECT 1; SELECT 2) , xi$b$o AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --;a;bEND ;
|  + (SELECT 1; SELECT 2) , xi$b$o AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col $$END ;$$ FROM t WHERE a = "select 1;a;b" ,;
| FROM t WHERE a = "$;$--;" SELECT (SELECT 1; SELECT 2) 1 $xky$; ;;$xkx$$xky$ FROM t WHERE a =;
| col E'\\' 1 X'01;10' AS x U&'/*;; ' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$END ;$$ FROM t WHERE a = "select 1;a;b" ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$END ;$$ FROM t WHERE a = "select 1;a;b" ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY _e_ FROM stdin;
| a;b
| a;b
| \.
|
| COPY eci FROM stdin;
| a;b
| a;b
| \.
| FROM t WHERE a = '; ;' + --; ;END ;
|  col e'\n;\';' ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _e_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _e_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY bka FROM stdin;
| 1	\N
| 1	\N
| \.
|  AS x --;; 
|  +;
| , CREATE FUNCTION syt_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$x$;'; END 1;
| 1 $$; end;$$ AS x "end;end;;;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bka FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bka FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + yw$i$s , --end;$x$;end;
|  col ";;" SELECT; COPY mcv FROM stdin;
| c;d;e
| 1	\N
| \.
|  SELECT '$;$);($x$;' , 'end;select 1;/*;\' , $ry$--;$ry$ FROM t WHERE a =;
| 1 ');(' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + yw$i$s , --end;$x$;end;
|  col ";;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + yw$i$s , --end;$x$;end;
|  col ";;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY xiz FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY iih FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xiz FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xiz FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = $jjq$END ;$jjq$ 1 --);(select 1;);(
|  1;
| SELECT --; 
|  col b'01;10' , /*/*/*;*/*/*/ SELECT;
| COPY inb FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $jjq$END ;$jjq$ 1 --);(select 1;);(
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $jjq$END ;$jjq$ 1 --);(select 1;);(
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT e''';' + "select 1;--;" 1; -- trailing note;
| FROM t WHERE a = U&'end;end;' col; col /*/*/*$x$;/ *;--;*/*/*/ SELECT b'01;10' +;
| COPY tci FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 32
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e''';' + "select 1;--;" 1;
< error
| 42601 syntax error at or near "1"
| position: 32
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT e''';' + "select 1;--;" 1;
< error
| 42601 syntax error at or near "1"
| position: 32
> exec
| 1 "END ;" ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "END ;" ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "END ;" ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col (SELECT 1; SELECT 2) AS x CREATE FUNCTION tfdw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$;$END ;'; END AS x /*/*select 1;--;$;$*/*/ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) AS x CREATE FUNCTION tfdw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$;$END ;'; END AS x /*/*select 1;--;$;$*/*/ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) AS x CREATE FUNCTION tfdw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$;$END ;'; END AS x /*/*select 1;--;$;$*/*/ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY wrc FROM stdin;
| a;b
| a;b
| \.
|
| SELECT x'01;10' SELECT $m$--;;;end;$m$ SELECT x'01;10' SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wrc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wrc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , '$;$--;' FROM t WHERE a =; -- trailing note;
| , p_$e$d col (SELECT 1; SELECT 2) 1;
| COPY cdy FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , '$;$--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , '$;$--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x ln$t$a SELECT --;
|  , 'select 1;' AS x; -- trailing note;
| 1 --select 1;$x$;
|  , $$--;end;END ;$$ FROM t WHERE a =;
| COPY fne FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ln$t$a SELECT --;
|  , 'select 1;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ln$t$a SELECT --;
|  , 'select 1;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col '/*;/*;' col $$--;--;a;b$$ 1;
| col (SELECT 1; SELECT 2) AS x;
| FROM t WHERE a = --/*;END ;
|  FROM t WHERE a =;AS x CREATE FUNCTION yksi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END col U&';;);(end;' + CREATE FUNCTION _cyf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;; '; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '/*;/*;' col $$--;--;a;b$$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '/*;/*;' col $$--;--;a;b$$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT b'01;10' col;
| COPY aqa FROM stdin;
| a;b
| a;b
| \.
| + e''';\n;' SELECT $qawz$);($qawx$$qawz$ col CREATE FUNCTION qpyr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "a"
| position: 43
> exec
| BEGIN
< ok 0 -
> exec
| SELECT b'01;10' col;
< error
| 22P02 ";" is not a valid binary digit
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| COPY aqa FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 23
> exec
| 1 E'\';\\\';' , CREATE FUNCTION tzyp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);('; END AS x;
| FROM t WHERE a = b'01;10' SELECT; COPY lnk FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 E'\';\\\';' , CREATE FUNCTION tzyp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);('; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 E'\';\\\';' , CREATE FUNCTION tzyp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);('; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x 'end;' AS x;
| 1 pm$_$y SELECT '/*;\' AS x --END ;$x$;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x 'end;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x 'end;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = za$y$p FROM t WHERE a = CREATE FUNCTION __lb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $;$'; SELECT 1; END , $lqac$; $lqac$ AS x;
| + --$x$;
|  SELECT ';;;;select 1;' FROM t WHERE a = U&'/*;$;$/*;' SELECT;
| col (SELECT 1; SELECT 2) AS x CREATE FUNCTION rvtn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END AS x;
| COPY pke FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = za$y$p FROM t WHERE a = CREATE FUNCTION __lb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $;$'; SELECT 1; END , $lqac$; $lqac$ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = za$y$p FROM t WHERE a = CREATE FUNCTION __lb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $;$'; SELECT 1; END , $lqac$; $lqac$ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY ink FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ink FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ink FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT X'01;10' AS x /*/*a;b*/*/ col (SELECT 1; SELECT 2) 1;SELECT ";;/*;;;" FROM t WHERE a =;
| FROM t WHERE a = 'a;b' SELECT yc$v$q col $y$; $y$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 34
> exec
| BEGIN
< ok 0 -
> exec
| SELECT X'01;10' AS x /*/*a;b*/*/ col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 34
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT X'01;10' AS x /*/*a;b*/*/ col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 34
> exec
| SELECT $g$END ;$g$ SELECT e'\';\\\';' +; -- trailing note;
| COPY hea FROM stdin;
| 1	\N
| 1	\N
| \.
|  , U&'--;' AS x --select 1;--;$;$
|  AS x;
| , "--;a;b--;" , CREATE FUNCTION sdqh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;end;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "e'\\';\\\\\\';'"
| position: 27
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $g$END ;$g$ SELECT e'\';\\\';' +;
< error
| 42601 syntax error at or near "e'\\';\\\\\\';'"
| position: 27
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $g$END ;$g$ SELECT e'\';\\\';' +;
< error
| 42601 syntax error at or near "e'\\';\\\\\\';'"
| position: 27
> exec
| col CREATE FUNCTION xeoo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;b'; END col CREATE FUNCTION ykbz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;);('; END +;
| SELECT CREATE FUNCTION fmqo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;a;b'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION xeoo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;b'; END col CREATE FUNCTION ykbz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;);('; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION xeoo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;b'; END col CREATE FUNCTION ykbz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;);('; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x ye$z$i AS x $felj$END ;;--;$felx$$felj$ FROM t WHERE a = "END ;" ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ye$z$i AS x $felj$END ;;--;$felx$$felj$ FROM t WHERE a = "END ;" ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ye$z$i AS x $felj$END ;;--;$felx$$felj$ FROM t WHERE a = "END ;" ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT 'a;b' AS x;
< ok 1 -
> exec
| BEGIN
< ok 0 -
> exec
| SELECT 'a;b' AS x;
< ok 1 -
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) 1 "select 1;a;bend;" 1 'END ;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) 1 "select 1;a;bend;" 1 'END ;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) 1 "select 1;a;bend;" 1 'END ;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY exl FROM stdin;
| c;d;e
| 1	\N
| \.
| AS x $tb$/*;$tx$$tb$ 1; -- trailing note;
| AS x $s$END ;$x$$s$ +;
| COPY ciw FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY exl FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY exl FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY bsx FROM stdin;
| 1	\N
| 1	\N
| \.
|
| SELECT $nmpy$select 1;$x$;$x$;$nmpx$$nmpy$ + (SELECT 1; SELECT 2) FROM t WHERE a =;
| FROM t WHERE a = E'\n;\\' SELECT fe$x$z col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bsx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bsx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 ----;
|  col $o$END ;$x$$o$ , CREATE FUNCTION ddsi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ----;
|  col $o$END ;$x$$o$ , CREATE FUNCTION ddsi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ----;
|  col $o$END ;$x$$o$ , CREATE FUNCTION ddsi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x --$x$;select 1;--;
|  FROM t WHERE a =;
| COPY ekl FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --$x$;select 1;--;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --$x$;select 1;--;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY cdf FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY eyy FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cdf FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cdf FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY vqm FROM stdin;
| a;b
| a;b
| \.
|
| SELECT /*/*;;*/*/ + /*/*/*$x$;*/*/*/ col /*/*$x$;$;$*/*/ SELECT;
| col (SELECT 1; SELECT 2) SELECT "end;;;" FROM t WHERE a = U&'--;end;' SELECT;
| AS x $$);(;;;$$ col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vqm FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vqm FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col /*/*/ *;; */*/ col;
| FROM t WHERE a = CREATE FUNCTION kzod() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;'; SELECT 1; END FROM t WHERE a =;+ e''';\\' 1 /*/*/*;;;*/*/*/ + ");(/*;--;" FROM t WHERE a =;
| AS x U&'end;' SELECT "$;$" col $nivf$);(END ;$nivx$$nivf$ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*/ *;; */*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*/ *;; */*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = '/*;end;END ;' 1 (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = '/*;end;END ;' 1 (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = '/*;end;END ;' 1 (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col py$r$p +;
| 1 --a;b
|  , U&'a;b' AS x zf$o$x FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col py$r$p +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col py$r$p +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) 1 CREATE FUNCTION dbjj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;--;'; SELECT 1; END +;
| FROM t WHERE a = vv$l$u FROM t WHERE a = /*/*;;$;$select 1;*/*/ + CREATE FUNCTION ghph() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;;;'; SELECT 1; END FROM t WHERE a =;SELECT b'01;10' AS x b'01;10' 1 "--;;END ;" col;
| SELECT /*/*/*; a;b*/*/*/ SELECT /*/*/ *;select 1;*/*/ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) 1 CREATE FUNCTION dbjj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;--;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) 1 CREATE FUNCTION dbjj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;--;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col X'01;10' 1 '--;end;' SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col X'01;10' 1 '--;end;' SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col X'01;10' 1 '--;end;' SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = ----;$;$select 1;
|  FROM t WHERE a = "end;end;END ;" +;
| col (SELECT 1; SELECT 2) SELECT ok$p$m +;
| col ----;
|  SELECT CREATE FUNCTION zhrd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x x'01;10' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ----;$;$select 1;
|  FROM t WHERE a = "end;end;END ;" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ----;$;$select 1;
|  FROM t WHERE a = "end;end;END ;" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x $jqhx$/*;$x$;$x$;@$jqhx$ FROM t WHERE a =;
| , E'\n;\\\\' FROM t WHERE a = e''';' AS x $awzu$$;$$awzu$ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $jqhx$/*;$x$;$x$;@$jqhx$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $jqhx$/*;$x$;$x$;@$jqhx$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , "/*;end;" SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "/*;end;" SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "/*;end;" SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY amn FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY amn FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY amn FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x (SELECT 1; SELECT 2) col 'END ;\' 1; -- trailing note;
| COPY sww FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) col 'END ;\' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) col 'END ;\' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY roc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY roc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY roc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ath FROM stdin;
| 1	\N
| c;d;e
| \.
|
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x;SELECT "end;" ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ath FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ath FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY bgy FROM stdin;
| 1	\N
| c;d;e
| \.
|
| , $$; ;;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bgy FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bgy FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col U&'; ;' col 'end;END ;$x$;\' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'; ;' col 'end;END ;$x$;\' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'; ;' col 'end;END ;$x$;\' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY ved FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY zyc FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ved FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ved FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 "$x$;" SELECT E'\';\\' , '$;$end;' ,;
| COPY aoz FROM stdin;
| a;b
| a;b
| \.
|
| AS x $yv$a;b; $yx$$yv$ FROM t WHERE a = /*/*/*;*/*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "$x$;" SELECT E'\';\\' , '$;$end;' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "$x$;" SELECT E'\';\\' , '$;$end;' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = --a;b;a;b
|  col /*/*/*/ *;;/ *;*/*/*/ AS x --;;END ;$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --a;b;a;b
|  col /*/*/*/ *;;/ *;*/*/*/ AS x --;;END ;$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --a;b;a;b
|  col /*/*/*/ *;;/ *;*/*/*/ AS x --;;END ;$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , "; );($;$" 1 U&';select 1;' AS x;
| AS x $$select 1;a;b$$ 1 CREATE FUNCTION sxpn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "; );($;$" 1 U&';select 1;' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "; );($;$" 1 U&';select 1;' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 B'01;10' 1 e'\n;\n;\n;' SELECT (SELECT 1; SELECT 2) SELECT;+ (SELECT 1; SELECT 2) FROM t WHERE a = /*END ;;*/ , x'01;10' 1;
| COPY htn FROM stdin;
| 1	\N
| c;d;e
| \.
|  col --$x$;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 B'01;10' 1 e'\n;\n;\n;' SELECT (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 B'01;10' 1 e'\n;\n;\n;' SELECT (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY fgp FROM stdin;
| a;b
| a;b
| \.
|
| 1 $ggaj$$x$;);(; $ggax$$ggaj$ ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fgp FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fgp FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT e'\n;\';\n;' , "/*;" SELECT "/*;" +; -- trailing note;
| AS x ');(;;' AS x al$c$_ SELECT rv$g$r ,;
| SELECT b'01;10' +;
< error
| 42601 syntax error at or near ""/*;""
| position: 36
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\n;\';\n;' , "/*;" SELECT "/*;" +;
< error
| 42601 syntax error at or near ""/*;""
| position: 36
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT e'\n;\';\n;' , "/*;" SELECT "/*;" +;
< error
| 42601 syntax error at or near ""/*;""
| position: 36
> exec
| COPY fyg FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fyg FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fyg FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x U&'$x$;' col (SELECT 1; SELECT 2) 1;
| COPY shf FROM stdin;
| 1	\N
| c;d;e
| \.
|
| SELECT e''';' FROM t WHERE a = x'01;10' 1 CREATE FUNCTION wt_v() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'$x$;' col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'$x$;' col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) 1 $_qa$END ;$_qx$$_qa$ , "$x$;" FROM t WHERE a =; -- trailing note;
| COPY wni FROM stdin;
| a;b
| 1	\N
| \.
|
| AS x $gkoy$; $x$;$gkox$$gkoy$ , /*$x$;*/ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) 1 $_qa$END ;$_qx$$_qa$ , "$x$;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) 1 $_qa$END ;$_qx$$_qa$ , "$x$;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY rav FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| COPY _pl FROM stdin;
| a;b
| a;b
| \.
| COPY gma FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY rav FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY rav FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 (SELECT 1; SELECT 2) 1 U&'a;ba;b' ,;
| COPY rd_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) 1 U&'a;ba;b' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) 1 U&'a;ba;b' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , E'\n;\n;\\' + je$a$z AS x;
| COPY dfw FROM stdin;
| c;d;e
| c;d;e
| \.
|
| SELECT "$;$END ;" 1 U&');(/*;--;' ,;SELECT (SELECT 1; SELECT 2) col (SELECT 1; SELECT 2) FROM t WHERE a = "/*;);(" FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , E'\n;\n;\\' + je$a$z AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , E'\n;\n;\\' + je$a$z AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col 'END ;\' col /*;;select 1;/ *;*/ col;
| SELECT U&';;' , (SELECT 1; SELECT 2) col; COPY sx_ FROM stdin;
| a;b
| c;d;e
| \.
| COPY jyx FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'END ;\' col /*;;select 1;/ *;*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'END ;\' col /*;;select 1;/ *;*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + /*--;select 1;END ;*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*--;select 1;END ;*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*--;select 1;END ;*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) , e''';' FROM t WHERE a =; -- trailing note;
| col (SELECT 1; SELECT 2) + E'\';\n;' FROM t WHERE a = /*/*/*/ *;*/*/*/ SELECT;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) , e''';' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) , e''';' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| col E'\\\n;'';' + CREATE FUNCTION smwn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$x$;END ;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col E'\\\n;'';' + CREATE FUNCTION smwn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$x$;END ;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col E'\\\n;'';' + CREATE FUNCTION smwn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$x$;END ;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + x'01;10' + CREATE FUNCTION syjf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;end;; '; SELECT 1; END AS x;
| COPY jvi FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + x'01;10' + CREATE FUNCTION syjf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;end;; '; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + x'01;10' + CREATE FUNCTION syjf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;end;; '; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY r_n FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY r_n FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY r_n FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , /*/*$x$;$;$;*/*/ + $lyz$END ;;$lyz$ FROM t WHERE a = CREATE FUNCTION adzm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;$;$'; SELECT 1; END AS x; AS x b'01;10' 1; -- trailing note;
| SELECT x'01;10' + /*select 1;);(; */ 1;
| FROM t WHERE a = "END ;$x$;" SELECT U&'a;b--;;' AS x /*/*;select 1;$;$*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*$x$;$;$;*/*/ + $lyz$END ;;$lyz$ FROM t WHERE a = CREATE FUNCTION adzm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;$;$'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*$x$;$;$;*/*/ + $lyz$END ;;$lyz$ FROM t WHERE a = CREATE FUNCTION adzm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;$;$'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 U&'--;' FROM t WHERE a = E'\';' , CREATE FUNCTION jq_t() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;END ;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'--;' FROM t WHERE a = E'\';' , CREATE FUNCTION jq_t() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;END ;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'--;' FROM t WHERE a = E'\';' , CREATE FUNCTION jq_t() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;END ;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY _hw FROM stdin;
| c;d;e
| c;d;e
| \.
|
| + CREATE FUNCTION dacj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;'; SELECT 1; END , E'\n;\';\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _hw FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _hw FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + $$/*;--;select 1;$$ , (SELECT 1; SELECT 2) , /*;$;$/ *;*/ AS x; -- trailing note;
| COPY ywg FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $$/*;--;select 1;$$ , (SELECT 1; SELECT 2) , /*;$;$/ *;*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $$/*;--;select 1;$$ , (SELECT 1; SELECT 2) , /*;$;$/ *;*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , ph$o$w , $vh$$x$;$x$;$vx$$vh$ + b'01;10' FROM t WHERE a =; FROM t WHERE a = (SELECT 1; SELECT 2) AS x "$;$a;b" +;
| COPY qvh FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ph$o$w , $vh$$x$;$x$;$vx$$vh$ + b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ph$o$w , $vh$$x$;$x$;$vx$$vh$ + b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 e'\\\';' + "select 1;end;a;b" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\\\';' + "select 1;end;a;b" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\\\';' + "select 1;end;a;b" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x /*/*--;end;*/*/ col b'01;10' FROM t WHERE a = "$x$;--;" 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*--;end;*/*/ col b'01;10' FROM t WHERE a = "$x$;--;" 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*--;end;*/*/ col b'01;10' FROM t WHERE a = "$x$;--;" 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) col;SELECT --a;b--;select 1;
|  + pg$p$n + e'\n;'';\n;' AS x; -- trailing note;
| , CREATE FUNCTION qaeb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';$;$--;'; SELECT 1; END + /*;;*/ AS x '$x$;/*;\' SELECT;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| SELECT x'01;10' 1 $avt$end;;; $avx$$avt$ col;
| COPY _yb FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY lqe FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' 1 $avt$end;;; $avx$$avt$ col;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT x'01;10' 1 $avt$end;;; $avx$$avt$ col;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| 1 'select 1;$x$;' col; -- trailing note;
| 1 $$$x$;$x$;select 1;$$ AS x;
| , CREATE FUNCTION aglk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END FROM t WHERE a =; -- trailing note;
| FROM t WHERE a = X'01;10' + b'01;10' FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 'select 1;$x$;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 'select 1;$x$;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY qbq FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| FROM t WHERE a = 'a;b);(\' 1 $yes$);($yes$ , --;;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qbq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qbq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT /*/*);(END ;a;b*/*/ SELECT b'01;10' 1;
| 1 e''';'';\';' 1 (SELECT 1; SELECT 2) + (SELECT 1; SELECT 2) col;
| COPY fhp FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "SELECT"
| position: 28
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*);(END ;a;b*/*/ SELECT b'01;10' 1;
< error
| 42601 syntax error at or near "SELECT"
| position: 28
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*/*);(END ;a;b*/*/ SELECT b'01;10' 1;
< error
| 42601 syntax error at or near "SELECT"
| position: 28
> exec
| COPY _xd FROM stdin;
| c;d;e
| 1	\N
| \.
|  , U&'select 1;$x$;; ' ,;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _xd FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _xd FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + "select 1;$x$;END ;" FROM t WHERE a = CREATE FUNCTION uold() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; ;);('; END 1;
| col (SELECT 1; SELECT 2) FROM t WHERE a = $$end;$x$;$$ AS x U&'--;' AS x;
| , $ibd$a;b;;$ibx$$ibd$ + $ibml$select 1;a;b$;$$ibml$ + $$);(/*;$$ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "select 1;$x$;END ;" FROM t WHERE a = CREATE FUNCTION uold() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; ;);('; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "select 1;$x$;END ;" FROM t WHERE a = CREATE FUNCTION uold() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; ;);('; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + --; 
|  , b'01;10' AS x X'01;10' col;
| 1 U&'$;$end;);(' FROM t WHERE a = $$a;b/*;$$ +; -- trailing note;
| , --; ;;
|  AS x;
| , hr$s$i col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --; 
|  , b'01;10' AS x X'01;10' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --; 
|  , b'01;10' AS x X'01;10' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) + --/*;END ;
|  SELECT;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) + --/*;END ;
|  SELECT;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) + --/*;END ;
|  SELECT;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| SELECT CREATE FUNCTION ukgv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END SELECT x'01;10' FROM t WHERE a =; -- trailing note;
| + E'\\'';\\' + (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION ukgv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END SELECT x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION ukgv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END SELECT x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| 1 X'01;10' 1;
| AS x CREATE FUNCTION dwkb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;$;$'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 X'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 X'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT CREATE FUNCTION auoc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END FROM t WHERE a = bs$n$n 1; -- trailing note;
| FROM t WHERE a = '; ' col /*/*/*end;a;b; */*/*/ SELECT E'\\' AS x; AS x /*/*--;*/*/ 1; -- trailing note;
| , /*/*/*);(*/*/*/ col --);($;$
|  1 "/*;end;;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION auoc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END FROM t WHERE a = bs$n$n 1;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION auoc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END FROM t WHERE a = bs$n$n 1;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| col (SELECT 1; SELECT 2) col "/*;" col (SELECT 1; SELECT 2) ,;
| COPY g_r FROM stdin;
| 1	\N
| a;b
| \.
|
| FROM t WHERE a = E'\';'';' AS x (SELECT 1; SELECT 2) 1 U&'END ;' col;
| , ';;a;b/*;\' col /*/*/*; ;;; */*/*/ + /*END ;;*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) col "/*;" col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) col "/*;" col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col pj$o$c ,; -- trailing note;
| SELECT ----;;
|  AS x CREATE FUNCTION wzbr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;a;b'; END ,;
| col e'\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col pj$o$c ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col pj$o$c ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , U&';a;b' AS x ee$t$h 1; -- trailing note;
| , '$x$;/*;\' + U&';;);(END ;' FROM t WHERE a = ';$;$' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&';a;b' AS x ee$t$h 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&';a;b' AS x ee$t$h 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col CREATE FUNCTION kysr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END col B'01;10' FROM t WHERE a = $$--;$$ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION kysr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END col B'01;10' FROM t WHERE a = $$--;$$ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION kysr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END col B'01;10' FROM t WHERE a = $$--;$$ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY dmu FROM stdin;
| 1	\N
| c;d;e
| \.
| , U&'$;$end;END ;' SELECT 'select 1;;$;$' col; COPY _ct FROM stdin;
| a;b
| 1	\N
| \.
|
| + ht$d$l 1 B'01;10' SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dmu FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dmu FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 "END ;);(;;" col B'01;10' 1 U&';;' +; SELECT '/*;$;$select 1;\' AS x "a;b);(;;" FROM t WHERE a =;
| + $lq_$end;);(end;$lq_$ + /*/*/*END ;);(;*/*/*/ FROM t WHERE a = x'01;10' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "END ;);(;;" col B'01;10' 1 U&';;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "END ;);(;;" col B'01;10' 1 U&';;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY qpx FROM stdin;
| 1	\N
| a;b
| \.
|
| col $$select 1;$$ col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qpx FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qpx FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = X'01;10' SELECT ");(" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = X'01;10' SELECT ");(" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = X'01;10' SELECT ");(" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , dc$g$u SELECT hk$e$j ,; -- trailing note;
| COPY rfa FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , dc$g$u SELECT hk$e$j ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , dc$g$u SELECT hk$e$j ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY nt_ FROM stdin;
| a;b
| 1	\N
| \.
| col $$END ;$$ col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY nt_ FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY nt_ FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ind FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY crz FROM stdin;
| c;d;e
| 1	\N
| \.
|
| SELECT B'01;10' + b'01;10' +;, (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ind FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ind FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + --$;$;;a;b
|  1; AS x ";;" ,; -- trailing note;
| AS x (SELECT 1; SELECT 2) AS x ');(\' AS x;
| COPY gbl FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --$;$;;a;b
|  1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --$;$;;a;b
|  1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 b'01;10' , ';\' AS x; -- trailing note;
| COPY cnk FROM stdin;
| a;b
| a;b
| \.
|
| COPY qtz FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 b'01;10' , ';\' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 b'01;10' , ';\' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + U&';;);(' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&';;);(' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&';;);(' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT --END ;
|  1;
| FROM t WHERE a = CREATE FUNCTION orgx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;'; SELECT 1; END 1;
| SELECT CREATE FUNCTION lwti() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;end;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "FROM"
| position: 20
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --END ;
|  1;
< ok 1 -
> exec
|
| FROM t WHERE a = CREATE FUNCTION orgx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 2
> exec
| ROLLBACK
< ok 0 -
> exec
|
| FROM t WHERE a = CREATE FUNCTION orgx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 2
> exec
| col U&');(;;end;' FROM t WHERE a = E'\\\';' 1;
| COPY oho FROM stdin;
| a;b
| 1	\N
| \.
|
| FROM t WHERE a = --select 1;
|  , 'end;$;$a;b' SELECT;
| COPY bsn FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&');(;;end;' FROM t WHERE a = E'\\\';' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&');(;;end;' FROM t WHERE a = E'\\\';' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a = U&';END ;/*;' AS x "$;$" col;
| COPY loe FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a = U&';END ;/*;' AS x "$;$" col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a = U&';END ;/*;' AS x "$;$" col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = U&'end;END ;$;$' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'end;END ;$;$' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'end;END ;$;$' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY tmy FROM stdin;
| 1	\N
| a;b
| \.
|  SELECT --select 1;$x$;/*;
|  1 fh$j$n +;
| FROM t WHERE a = /*/*/*);(select 1;*/*/*/ ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY tmy FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY tmy FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + $cbu$; ; $cbu$ col U&'--;--;' SELECT /*$x$;*/ 1; -- trailing note;
| col f_$k$j 1 e''';\n;\';' FROM t WHERE a =;
| col U&';;end;' FROM t WHERE a =; COPY cys FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $cbu$; ; $cbu$ col U&'--;--;' SELECT /*$x$;*/ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $cbu$; ; $cbu$ col U&'--;--;' SELECT /*$x$;*/ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY bhd FROM stdin;
| c;d;e
| c;d;e
| \.
|
| 1 $$$x$;a;b--;$$ 1 'select 1;' 1;col (SELECT 1; SELECT 2) AS x CREATE FUNCTION fdsn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;END ;'; SELECT 1; END SELECT (SELECT 1; SELECT 2) FROM t WHERE a =; col X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bhd FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bhd FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) 1;
| col CREATE FUNCTION yaoi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END SELECT; col $$);($$ , E'\';' AS x; -- trailing note;
| FROM t WHERE a = E'\\'';'';' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , --$x$;
|  SELECT /*/*/*/ *;;*/*/*/ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --$x$;
|  SELECT /*/*/*/ *;;*/*/*/ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --$x$;
|  SELECT /*/*/*/ *;;*/*/*/ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = $nkhf$;;$nkhf$ AS x $on$--;;;;$ox$$on$ SELECT X'01;10' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $nkhf$;;$nkhf$ AS x $on$--;;;;$ox$$on$ SELECT X'01;10' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $nkhf$;;$nkhf$ AS x $on$--;;;;$ox$$on$ SELECT X'01;10' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + /*/*/*;;--;);(*/*/*/ 1 /*/*/*select 1;;*/*/*/ AS x; -- trailing note;
| AS x b'01;10' FROM t WHERE a = /*);(/ *;END ;*/ + b'01;10' SELECT;
| COPY wlc FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*/*;;--;);(*/*/*/ 1 /*/*/*select 1;;*/*/*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*/*;;--;);(*/*/*/ 1 /*/*/*select 1;;*/*/*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + $q$/*;$q$ FROM t WHERE a =;SELECT '$x$;end;\' AS x CREATE FUNCTION whke() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END FROM t WHERE a =;col ");(/*;" 1;FROM t WHERE a = $rfb$;);(--;$rfb$ , ');(END ;$;$\' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $q$/*;$q$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $q$/*;$q$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col 'end;' FROM t WHERE a =;1 /*/*/*; */*/*/ col $ki$;$ki$ 1 /*$x$;);(*/ AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 X'01;10' FROM t WHERE a =;
| COPY pzk FROM stdin;
| a;b
| a;b
| \.
| SELECT e'\n;' col "$;$a;ba;b" SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + /*/*/*--;*/*/*/ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*/*--;*/*/*/ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*/*--;*/*/*/ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 --; select 1;
|  , e'\\' , --);(
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --; select 1;
|  , e'\\' , --);(
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --; select 1;
|  , e'\\' , --);(
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 /*END ;$x$;; */ col ';;\' col "$x$;);();(" col;AS x '--;$;$\' ,;
| col B'01;10' ,;
| COPY cqv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*END ;$x$;; */ col ';;\' col "$x$;);();(" col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*END ;$x$;; */ col ';;\' col "$x$;);();(" col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 --select 1;$x$;/*;
|  AS x /*/*/*$;$*/*/*/ AS x (SELECT 1; SELECT 2) 1;
| 1 hu$v$t AS x e'\n;' 1 ga$e$g 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --select 1;$x$;/*;
|  AS x /*/*/*$;$*/*/*/ AS x (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --select 1;$x$;/*;
|  AS x /*/*/*$;$*/*/*/ AS x (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY aex FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY aex FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY aex FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + '--;a;b' FROM t WHERE a = x'01;10' + en$x$f FROM t WHERE a =;
| , U&'$;$END ;' 1 x'01;10' , '; ;;);(' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + '--;a;b' FROM t WHERE a = x'01;10' + en$x$f FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + '--;a;b' FROM t WHERE a = x'01;10' + en$x$f FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT CREATE FUNCTION vhgs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;select 1;'; SELECT 1; END + $$;--;$$ AS x --a;b; a;b
|  SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION vhgs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;select 1;'; SELECT 1; END + $$;--;$$ AS x --a;b; a;b
|  SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION vhgs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;select 1;'; SELECT 1; END + $$;--;$$ AS x --a;b; a;b
|  SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| col "end;select 1;" col E'\n;\\\n;' col e'\';\';' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "end;select 1;" col E'\n;\\\n;' col e'\';\';' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "end;select 1;" col E'\n;\\\n;' col e'\';\';' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT e'\\\n;' AS x --$;$$;$
|  FROM t WHERE a = b'01;10' SELECT;
| COPY dgm FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "SELECT"
| position: 58
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\\\n;' AS x --$;$$;$
|  FROM t WHERE a = b'01;10' SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 58
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT e'\\\n;' AS x --$;$$;$
|  FROM t WHERE a = b'01;10' SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 58
> exec
| COPY _p_ FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| , B'01;10' AS x qb$e$_ AS x;
| 1 "end;;;" SELECT;
| AS x fd$a$g col CREATE FUNCTION xejv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;a;b'; SELECT 1; END +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _p_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _p_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 $$/*;END ;select 1;$$ col (SELECT 1; SELECT 2) SELECT e''';' +;
| 1 b'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $$/*;END ;select 1;$$ col (SELECT 1; SELECT 2) SELECT e''';' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $$/*;END ;select 1;$$ col (SELECT 1; SELECT 2) SELECT e''';' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY yfk FROM stdin;
| a;b
| 1	\N
| \.
|  COPY ung FROM stdin;
| 1	\N
| 1	\N
| \.
|
| col $$--;/*;$$ col --end;--;
|  col "select 1;END ;" ,;
| SELECT E'\n;' col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yfk FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yfk FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY atf FROM stdin;
| a;b
| a;b
| \.
|
| col ";" , U&'$x$;; ' SELECT;
| col '$;$a;b' AS x (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY atf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY atf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x --END ;--;
|  SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --END ;--;
|  SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --END ;--;
|  SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col e'\';\n;\n;' SELECT 'end;--;' FROM t WHERE a = wy$_$p +; FROM t WHERE a = dp$g$b SELECT --$x$;end;
|  AS x qb$s$p ,; -- trailing note;
| col E''';\\' SELECT E'\\\';\\' + '; $x$;' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e'\';\n;\n;' SELECT 'end;--;' FROM t WHERE a = wy$_$p +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e'\';\n;\n;' SELECT 'end;--;' FROM t WHERE a = wy$_$p +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + '$x$;END ;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + '$x$;END ;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + '$x$;END ;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY owq FROM stdin;
| c;d;e
| c;d;e
| \.
|
| FROM t WHERE a = U&'select 1;' SELECT --);($x$;
|  FROM t WHERE a =; -- trailing note;
| COPY kps FROM stdin;
| a;b
| a;b
| \.
|
| col $$end;);(/*;$$ +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY owq FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY owq FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x "select 1;;;--;" 1 e'\';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "select 1;;;--;" 1 e'\';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "select 1;;;--;" 1 e'\';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , (SELECT 1; SELECT 2) 1 e'\n;\\' AS x $$; $$ FROM t WHERE a =;
| 1 U&';;--;/*;' AS x U&'/*;--;END ;' ,;
| COPY cjf FROM stdin;
| c;d;e
| 1	\N
| \.
|  AS x '$x$;a;b;;\' 1 px$i$t 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) 1 e'\n;\\' AS x $$; $$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) 1 e'\n;\\' AS x $$; $$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY gxv FROM stdin;
| c;d;e
| a;b
| \.
|  col '$x$;;' SELECT $$end;$$ 1; 1 $$select 1;select 1;$;$$$ SELECT /*/*end;$;$*/*/ 1;
| 1 /*/*; / *;;;*/*/ , x'01;10' SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gxv FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gxv FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY owt FROM stdin;
| a;b
| a;b
| \.
|
| col U&'; ; ' 1 ';;);(end;' SELECT 'a;b$;$--;\' +;
| FROM t WHERE a = --/*;);(a;b
|  + E'\';' , E'\\' AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY owt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY owt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 U&';$;$/*;' +;
| col $zosj$/*;--;$x$;$zosx$$zosj$ col B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&';$;$/*;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&';$;$/*;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 --END ;END ;a;b
|  SELECT; -- trailing note;
| col CREATE FUNCTION tfbv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;'; SELECT 1; END AS x /*/ *;*/ col;
| , $$END ;; ;;$$ AS x E''';'';' AS x;
| , (SELECT 1; SELECT 2) AS x E'\\' AS x /*END ;*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --END ;END ;a;b
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --END ;END ;a;b
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col U&'$x$;' 1 $$/*;;;; $$ FROM t WHERE a =;
| col CREATE FUNCTION ddif() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;select 1;/*;'; END 1 x'01;10' , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'$x$;' 1 $$/*;;;; $$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'$x$;' 1 $$/*;;;; $$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = --;
|  col bq$x$w 1 "/*;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --;
|  col bq$x$w 1 "/*;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --;
|  col bq$x$w 1 "/*;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = $$/*;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$/*;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$/*;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 mu$s$h +; -- trailing note;
| col "$;$" AS x /*/*select 1;; */*/ 1 to$m$_ FROM t WHERE a =;
| FROM t WHERE a = /*/*END ;select 1;*/*/ FROM t WHERE a = U&';end;$x$;' 1 --/*;END ;
|  +;COPY pbd FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 mu$s$h +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 mu$s$h +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x $vco$$x$;select 1;$vcx$$vco$ FROM t WHERE a =;
| , /*/*/*end;;*/*/*/ , (SELECT 1; SELECT 2) AS x; FROM t WHERE a = E''';'';' + /*--;--;*/ AS x b'01;10' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $vco$$x$;select 1;$vcx$$vco$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $vco$$x$;select 1;$vcx$$vco$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY vcb FROM stdin;
| a;b
| c;d;e
| \.
|
| AS x 'select 1;;' AS x U&'END ;);(' , '$x$;$x$;$x$;' 1; AS x ';;);(' +; -- trailing note;
| col /*/*;;END ;; */*/ , U&'/*;' SELECT '$x$;select 1;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vcb FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vcb FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT $$;;a;b);($$ 1;
< error
| 42601 syntax error at or near "1"
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$;;a;b);($$ 1;
< error
| 42601 syntax error at or near "1"
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$;;a;b);($$ 1;
< error
| 42601 syntax error at or near "1"
| position: 21
> exec
| , $$);(a;b;;$$ AS x;
| FROM t WHERE a = /*/*/*;;; */*/*/ col $$select 1;; select 1;$$ SELECT;
| FROM t WHERE a = --a;bend;);(
|  , B'01;10' SELECT $ly$--;);($lx$$ly$ AS x; col U&';;$x$;' 1 /*/*;;end;$;$*/*/ + e''';\n;' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $$);(a;b;;$$ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $$);(a;b;;$$ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY hje FROM stdin;
| 1	\N
| a;b
| \.
|  AS x E''';'';' AS x $f$;;/*;$x$$f$ col gc$h$s FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hje FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hje FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x U&'/*;' 1;
| FROM t WHERE a = 'END ;/*;' col e''';' AS x;
| col e''';' 1 U&';;' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'/*;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'/*;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT /*/*/*end;*/*/*/ 1 E'\';'';\';' + (SELECT 1; SELECT 2) FROM t WHERE a =;
| SELECT "$x$;" col B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "E'\\';'';\\';'"
| position: 27
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*/*end;*/*/*/ 1 E'\';'';\';' + (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "E'\\';'';\\';'"
| position: 27
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*/*/*end;*/*/*/ 1 E'\';'';\';' + (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "E'\\';'';\\';'"
| position: 27
> exec
| COPY r_g FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY r_g FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY r_g FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col U&'END ;' FROM t WHERE a = e'\';' col; COPY hea FROM stdin;
| 1	\N
| c;d;e
| \.
|
| col U&');(' AS x CREATE FUNCTION tqqm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'END ;' FROM t WHERE a = e'\';' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'END ;' FROM t WHERE a = e'\';' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY uhl FROM stdin;
| a;b
| a;b
| \.
| 1 (SELECT 1; SELECT 2) col;
| AS x gd$m$f AS x --end;;;
|  , U&'--;; ' col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY uhl FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY uhl FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY czo FROM stdin;
| 1	\N
| 1	\N
| \.
|
| + U&'END ;' ,;
| SELECT ");($x$;" 1 U&';;/*;' , (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY czo FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY czo FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY hvb FROM stdin;
| 1	\N
| 1	\N
| \.
| AS x $$--;$x$;$$ 1 '; $;$;' , --;
|  1;
| 1 e'\';\';' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hvb FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hvb FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 CREATE FUNCTION kjxv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;END ;'; SELECT 1; END SELECT X'01;10' , $$select 1;; $$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION kjxv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;END ;'; SELECT 1; END SELECT X'01;10' , $$select 1;; $$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION kjxv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;END ;'; SELECT 1; END SELECT X'01;10' , $$select 1;; $$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT --$x$;--;; 
|  +; -- trailing note;
| AS x rx$v$u SELECT;
| COPY kfs FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --$x$;--;; 
|  +;
< error
| 42601 syntax error at or near ";"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --$x$;--;; 
|  +;
< error
| 42601 syntax error at or near ";"
| position: 22
> exec
| 1 /*/*$;$--;*/*/ AS x;
| FROM t WHERE a = CREATE FUNCTION cqhj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END FROM t WHERE a =; + "/*;;" col --END ;a;b
|  FROM t WHERE a =;
| col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*$;$--;*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*$;$--;*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + b'01;10' ,;COPY nac FROM stdin;
| 1	\N
| a;b
| \.
|
| col '$;$/*;' , X'01;10' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + b'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + b'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT U&';select 1;);(' , --$;$a;b
|  AS x; -- trailing note;
| 1 CREATE FUNCTION _vlu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$--;select 1;'; SELECT 1; END AS x; SELECT e''';' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 38
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&';select 1;);(' , --$;$a;b
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 38
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&';select 1;);(' , --$;$a;b
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 38
> exec
| COPY ico FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ico FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ico FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , ');(END ;' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ');(END ;' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ');(END ;' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY lda FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lda FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lda FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x CREATE FUNCTION tdya() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END SELECT CREATE FUNCTION pbzl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;/*;'; SELECT 1; END + U&'$;$' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION tdya() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END SELECT CREATE FUNCTION pbzl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;/*;'; SELECT 1; END + U&'$;$' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION tdya() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END SELECT CREATE FUNCTION pbzl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;/*;'; SELECT 1; END + U&'$;$' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 $w$;;select 1;$x$$w$ + E'\\\';'';' SELECT /*/*; $x$;;*/*/ ,;
| 1 --;;--;;;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $w$;;select 1;$x$$w$ + E'\\\';'';' SELECT /*/*; $x$;;*/*/ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $w$;;select 1;$x$$w$ + E'\\\';'';' SELECT /*/*; $x$;;*/*/ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY kes FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY epq FROM stdin;
| a;b
| 1	\N
| \.
|
| , 'select 1;$;$$x$;' col; + /*;select 1;*/ , /*/*;a;ba;b*/*/ AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kes FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kes FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + --;;
|  AS x b'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --;;
|  AS x b'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --;;
|  AS x b'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + CREATE FUNCTION weob() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;end;'; END AS x CREATE FUNCTION _kkq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;--;'; END 1;col 'END ;select 1;\' FROM t WHERE a = '$x$;);(' +;
| SELECT /*/*--;*/*/ col --$x$;--;
|  +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION weob() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;end;'; END AS x CREATE FUNCTION _kkq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;--;'; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION weob() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;end;'; END AS x CREATE FUNCTION _kkq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;--;'; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , U&';;end;END ;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&';;end;END ;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&';;end;END ;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + $hwu$end;$hwu$ SELECT 'a;b$;$select 1;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $hwu$end;$hwu$ SELECT 'a;b$;$select 1;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $hwu$end;$hwu$ SELECT 'a;b$;$select 1;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 th$p$e , $ik$$x$;$ix$$ik$ col /*/*/*; $x$;/ *;*/*/*/ AS x;col U&'/*;a;bEND ;' FROM t WHERE a = ');(;;' ,;
| SELECT X'01;10' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 th$p$e , $ik$$x$;$ix$$ik$ col /*/*/*; $x$;/ *;*/*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 th$p$e , $ik$$x$;$ix$$ik$ col /*/*/*; $x$;/ *;*/*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT U&'end;;;' 1 --);(;;
|  FROM t WHERE a = CREATE FUNCTION ejrq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;select 1;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "1"
| position: 19
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'end;;;' 1 --);(;;
|  FROM t WHERE a = CREATE FUNCTION ejrq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;select 1;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "1"
| position: 19
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'end;;;' 1 --);(;;
|  FROM t WHERE a = CREATE FUNCTION ejrq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;select 1;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "1"
| position: 19
> exec
| col CREATE FUNCTION wwjb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; '; END col;
| , 'END ;end;; ' 1 /*/*select 1;$x$;*/*/ AS x (SELECT 1; SELECT 2) col; -- trailing note;
| 1 '$x$;/*;select 1;\' SELECT; AS x U&';select 1;/*;' col $ggtt$a;bEND ;$ggtx$$ggtt$ , --end;
|  col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION wwjb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; '; END col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION wwjb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; '; END col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) ,;
| SELECT CREATE FUNCTION hlhx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; END +;
| 1 '$x$;' + --a;b$;$;
|  1; FROM t WHERE a = U&'/*;;;' + --$;$; 
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x --; select 1;
|  FROM t WHERE a =;
| COPY gek FROM stdin;
| a;b
| c;d;e
| \.
| + $$;$$ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --; select 1;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --; select 1;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x "end;; " + X'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "end;; " + X'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "end;; " + X'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col E''';\n;'';' col oh$v$n 1 U&';;$;$' +;AS x zr$v$m + --;;$;$
|  FROM t WHERE a = --select 1;end;;;
|  col;
| AS x CREATE FUNCTION ypot() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END , nq$_$_ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col E''';\n;'';' col oh$v$n 1 U&';;$;$' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col E''';\n;'';' col oh$v$n 1 U&';;$;$' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT to$w$c SELECT; -- trailing note;
| , "; end;" 1 $$$x$;; $$ col U&'END ;' col;
| col "--;END ;; " FROM t WHERE a = ";END ;;" +;
< error
| 42601 syntax error at or near ","
| position: 41
> exec
| BEGIN
< ok 0 -
> exec
| SELECT to$w$c SELECT;
< error
| 42703 column "to$w$c" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|  -- trailing note;
| , "; end;" 1 $$$x$;; $$ col U&'END ;' col;
< error
| 42601 syntax error at or near ","
| position: 20
> exec
| col $oz$select 1;$ox$$oz$ FROM t WHERE a = '$;$\' AS x; + /*/*/*--;end;select 1;*/*/*/ + b'01;10' FROM t WHERE a =;
| + CREATE FUNCTION etaz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $oz$select 1;$ox$$oz$ FROM t WHERE a = '$;$\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $oz$select 1;$ox$$oz$ FROM t WHERE a = '$;$\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY jtg FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jtg FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jtg FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x e'\n;\\\\' SELECT CREATE FUNCTION gisr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END + U&';;' col; col --; $x$;select 1;
|  + X'01;10' 1 CREATE FUNCTION chr_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;--;'; SELECT 1; END 1; COPY aol FROM stdin;
| a;b
| a;b
| \.
|
| col 'END ;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x e'\n;\\\\' SELECT CREATE FUNCTION gisr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END + U&';;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x e'\n;\\\\' SELECT CREATE FUNCTION gisr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END + U&';;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col /*/*select 1;a;b$;$*/*/ SELECT;
| SELECT E''';\\'';' 1 e'\';\\\n;' +;
| AS x /*$x$;END ;--;*/ col U&'--;$x$;' SELECT B'01;10' SELECT; -- trailing note;
| COPY izz FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*select 1;a;b$;$*/*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*select 1;a;b$;$*/*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY vyd FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY eru FROM stdin;
| 1	\N
| c;d;e
| \.
|
| , CREATE FUNCTION ufbf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vyd FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vyd FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY acf FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| FROM t WHERE a = U&';;; ' SELECT ----;
|  FROM t WHERE a = 'select 1;select 1;;' 1;
| AS x U&'$;$' +; 1 CREATE FUNCTION mmjl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';$x$;'; END , $o$; a;b$;$$x$$o$ AS x CREATE FUNCTION blhb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$/*;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY acf FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY acf FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , $z$);(a;b$x$$z$ ,; COPY jem FROM stdin;
| c;d;e
| c;d;e
| \.
|
| col uf$_$p col; -- trailing note;
| col CREATE FUNCTION zqyp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$);('; END + /*/*/*; ;;a;b*/*/*/ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $z$);(a;b$x$$z$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $z$);(a;b$x$$z$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x $$; end;$$ col '/*;\' SELECT --; 
|  SELECT;
| FROM t WHERE a = $k$$x$;$x$$k$ SELECT B'01;10' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $$; end;$$ col '/*;\' SELECT --; 
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $$; end;$$ col '/*;\' SELECT --; 
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = ");(END ;" AS x B'01;10' , CREATE FUNCTION fhqm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;a;b'; SELECT 1; END 1;
| , vx$a$k FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ");(END ;" AS x B'01;10' , CREATE FUNCTION fhqm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;a;b'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ");(END ;" AS x B'01;10' , CREATE FUNCTION fhqm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;a;b'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY imf FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY dxi FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY imf FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY imf FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x "END ;$x$;" +;
| COPY f__ FROM stdin;
| c;d;e
| a;b
| \.
|
| , ");(end;$;$" SELECT; -- trailing note;
| AS x $ovqd$);(--;end;$ovqd$ + (SELECT 1; SELECT 2) 1 U&'/*;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "END ;$x$;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "END ;$x$;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 x'01;10' 1;
| AS x ";);(" col U&'select 1;a;b);(' + X'01;10' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 x'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 x'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , '$;$' 1 hu$h$u + "select 1;END ;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , '$;$' 1 hu$h$u + "select 1;END ;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , '$;$' 1 hu$h$u + "select 1;END ;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY uoe FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY sgx FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| + /*/*/*; a;b*/*/*/ AS x CREATE FUNCTION jkao() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $x$;'; SELECT 1; END 1 $zpel$$x$;$x$;select 1;$zpex$$zpel$ +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY uoe FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY uoe FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + $ozk$);($ozx$$ozk$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $ozk$);($ozx$$ozk$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $ozk$);($ozx$$ozk$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 CREATE FUNCTION lcai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END FROM t WHERE a = (SELECT 1; SELECT 2) , --end;$;$
|  AS x;
| COPY yok FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION lcai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END FROM t WHERE a = (SELECT 1; SELECT 2) , --end;$;$
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION lcai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END FROM t WHERE a = (SELECT 1; SELECT 2) , --end;$;$
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = ";end;);(" , U&'--;; ' SELECT e''';\';'';' SELECT; -- trailing note;
| SELECT $$; ; select 1;$$ FROM t WHERE a = zs$l$d col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ";end;);(" , U&'--;; ' SELECT e''';\';'';' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ";end;);(" , U&'--;; ' SELECT e''';\';'';' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + --$x$;/*;
|  SELECT /*/*$;$*/*/ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --$x$;/*;
|  SELECT /*/*$;$*/*/ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --$x$;/*;
|  SELECT /*/*$;$*/*/ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , /*/ *;*/ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/ *;*/ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/ *;*/ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + U&'a;b/*;' FROM t WHERE a =;
| 1 (SELECT 1; SELECT 2) SELECT gs$j$t SELECT;
| 1 /*/*/*$;$$x$;*/*/*/ 1 '/*;\' col;
| AS x --end;
|  , x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'a;b/*;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'a;b/*;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY u_d FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY u_d FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY u_d FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , --/*;a;bend;
|  1 "$;$a;b;;" , (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --/*;a;bend;
|  1 "$;$a;b;;" , (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --/*;a;bend;
|  1 "$;$a;b;;" , (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT CREATE FUNCTION rdhn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; '; SELECT 1; END , CREATE FUNCTION akji() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$$x$;'; END ,;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION rdhn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; '; SELECT 1; END , CREATE FUNCTION akji() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$$x$;'; END ,;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION rdhn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; '; SELECT 1; END , CREATE FUNCTION akji() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$$x$;'; END ,;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| AS x E''';' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x E''';' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x E''';' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY soh FROM stdin;
| c;d;e
| 1	\N
| \.
|  COPY _en FROM stdin;
| a;b
| 1	\N
| \.
|
| AS x um$q$q SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY soh FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY soh FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col --; 
|  1 --; end;
|  FROM t WHERE a =; 1 ';a;ba;b' SELECT "select 1;select 1;" , /*/*/ *;*/*/ 1;
| col (SELECT 1; SELECT 2) , '--;' FROM t WHERE a = (SELECT 1; SELECT 2) +;
| + '/*;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --; 
|  1 --; end;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --; 
|  1 --; end;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 B'01;10' AS x x'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 B'01;10' AS x x'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 B'01;10' AS x x'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT --/*;--;
|  FROM t WHERE a =;
| AS x E'\';\';' + X'01;10' + U&'; ' 1;
| , CREATE FUNCTION hoez() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END FROM t WHERE a = e'\n;' SELECT;
< error
| 42601 syntax error at or near ";"
| position: 34
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --/*;--;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 34
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --/*;--;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 34
> exec
| SELECT /*/*end;*/*/ , $$a;b$$ SELECT;
| COPY mtc FROM stdin;
| 1	\N
| 1	\N
| \.
| + $cx$/*;end;a;b$cx$ col CREATE FUNCTION huu_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;/*;end;'; SELECT 1; END 1 --/*;
|  SELECT;
| FROM t WHERE a = $$/*;$$ + ht$_$r 1;
< error
| 42601 syntax error at or near ","
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*end;*/*/ , $$a;b$$ SELECT;
< error
| 42601 syntax error at or near ","
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*/*end;*/*/ , $$a;b$$ SELECT;
< error
| 42601 syntax error at or near ","
| position: 21
> exec
| + b'01;10' ,;AS x "$;$);();(" col /*/*--;;;--;*/*/ 1 B'01;10' +;
| , ';end;' ,;
| + $dqwo$);(select 1;$dqwx$$dqwo$ + b'01;10' 1 qq$s$b col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + b'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + b'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col CREATE FUNCTION zlgb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$;$'; SELECT 1; END + U&'; ;;' + CREATE FUNCTION vzib() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION zlgb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$;$'; SELECT 1; END + U&'; ;;' + CREATE FUNCTION vzib() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION zlgb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$;$'; SELECT 1; END + U&'; ;;' + CREATE FUNCTION vzib() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + ud$x$z 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ud$x$z 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ud$x$z 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY xri FROM stdin;
| 1	\N
| a;b
| \.
|
| COPY aca FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xri FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xri FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x (SELECT 1; SELECT 2) , '$;$' +;
| FROM t WHERE a = "END ;;;" 1 U&');(END ;a;b' FROM t WHERE a = $m$; $m$ +;
| COPY tdm FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) , '$;$' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) , '$;$' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY kab FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| SELECT ");(;; " ,;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kab FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kab FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT ");(" SELECT;COPY zaz FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| FROM t WHERE a = b'01;10' col;
| 1 ----;a;b
|  ,;
< error
| 42601 syntax error at or near "a"
| position: 42
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ");(" SELECT;
< error
| 42703 column ");(" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zaz FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY pyq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pyq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pyq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , E'\\' + (SELECT 1; SELECT 2) + b'01;10' +;
| , B'01;10' SELECT (SELECT 1; SELECT 2) SELECT (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , E'\\' + (SELECT 1; SELECT 2) + b'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , E'\\' + (SELECT 1; SELECT 2) + b'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , ");(" FROM t WHERE a = /*end;select 1;end;*/ SELECT x'01;10' AS x;
| , /*/*/*/ *;; ; */*/*/ AS x 'select 1;--;\' , $$select 1;$$ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ");(" FROM t WHERE a = /*end;select 1;end;*/ SELECT x'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ");(" FROM t WHERE a = /*end;select 1;end;*/ SELECT x'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x /*$;$*/ FROM t WHERE a = e''';\n;' +;COPY fvl FROM stdin;
| 1	\N
| 1	\N
| \.
| + yt$_$m SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*$;$*/ FROM t WHERE a = e''';\n;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*$;$*/ FROM t WHERE a = e''';\n;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY crv FROM stdin;
| 1	\N
| c;d;e
| \.
|
| + et$v$v , --END ;end;
|  AS x $ggo$a;b; $ggo$ col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY crv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY crv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY toa FROM stdin;
| a;b
| a;b
| \.
|
| col e'\n;'';' AS x '$x$;END ;' FROM t WHERE a = "a;b" 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY toa FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY toa FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT --$;$--;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 18
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --$;$--;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 18
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --$;$--;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 18
> exec
| COPY umu FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY umu FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY umu FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , CREATE FUNCTION tm_k() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x --END ;/*;; 
|  +; -- trailing note;
| COPY aly FROM stdin;
| 1	\N
| 1	\N
| \.
|  AS x ';;select 1;\' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION tm_k() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x --END ;/*;; 
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION tm_k() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x --END ;/*;; 
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = 'end;$x$;' col; SELECT cw$l$b FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = 'end;$x$;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = 'end;$x$;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + jx$v$k col E'\\' FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a =;
| FROM t WHERE a = $mtx$END ;end;);($mtx$ AS x --end;;a;b
|  , --$;$; 
|  AS x;
| COPY xhp FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + jx$v$k col E'\\' FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + jx$v$k col E'\\' FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 CREATE FUNCTION lf_u() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;a;b'; SELECT 1; END SELECT $$;$$ SELECT;
| COPY ihi FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| COPY hft FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION lf_u() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;a;b'; SELECT 1; END SELECT $$;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION lf_u() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;a;b'; SELECT 1; END SELECT $$;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + E'\\'';'';' + --);(/*;;
|  1 $zjtf$$;$end;$zjtf$ +;
| , /*;;;;--;*/ FROM t WHERE a = U&'end;a;b' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\\'';'';' + --);(/*;;
|  1 $zjtf$$;$end;$zjtf$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\\'';'';' + --);(/*;;
|  1 $zjtf$$;$end;$zjtf$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY gzo FROM stdin;
| a;b
| a;b
| \.
|  AS x "--;; ;" AS x;
| col /*/*/*$x$;;;*/*/*/ , B'01;10' +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gzo FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gzo FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col 'select 1;\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'select 1;\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'select 1;\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x /*/*/*$x$;*/*/*/ AS x U&'end;' AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*$x$;*/*/*/ AS x U&'end;' AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*$x$;*/*/*/ AS x U&'end;' AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT E'\n;' AS x "; " AS x "$;$;$;$" 1; -- trailing note;
| col CREATE FUNCTION loey() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;END ;;;'; END AS x CREATE FUNCTION cxqs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; END ;;'; END FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ""; ""
| position: 20
> exec
| BEGIN
< ok 0 -
> exec
| SELECT E'\n;' AS x "; " AS x "$;$;$;$" 1;
< error
| 42601 syntax error at or near ""; ""
| position: 20
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT E'\n;' AS x "; " AS x "$;$;$;$" 1;
< error
| 42601 syntax error at or near ""; ""
| position: 20
> exec
| FROM t WHERE a = B'01;10' SELECT --end;end;a;b
|  col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = B'01;10' SELECT --end;end;a;b
|  col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = B'01;10' SELECT --end;end;a;b
|  col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY jxj FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jxj FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jxj FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ydl FROM stdin;
| 1	\N
| c;d;e
| \.
|
| col --/*;$x$;
|  + uk$a$n SELECT; -- trailing note;
| + _r$e$e AS x 'a;b\' col b'01;10' +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ydl FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ydl FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , b'01;10' SELECT ig$c$e col; -- trailing note;
| COPY yhz FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , b'01;10' SELECT ig$c$e col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , b'01;10' SELECT ig$c$e col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 b'01;10' AS x 'a;ba;b' 1 e'\n;\\'';' ,;
| FROM t WHERE a = CREATE FUNCTION hqbo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; select 1;select 1;'; SELECT 1; END ,;AS x $mcy$END ;;$mcy$ + U&'/*;);($;$' ,;
| COPY azj FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 b'01;10' AS x 'a;ba;b' 1 e'\n;\\'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 b'01;10' AS x 'a;ba;b' 1 e'\n;\\'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY rbu FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| , /*/*; --;*/*/ 1; -- trailing note;
| + (SELECT 1; SELECT 2) 1 /*/*$;$;; */*/ 1 --select 1;/*;$x$;
|  AS x; -- trailing note;
| 1 /*$;$END ;*/ col '$x$;end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY rbu FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY rbu FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT xd$n$m AS x B'01;10' + $$/*;$$ 1; 1 /*/*/*select 1;*/*/*/ + wd$n$h ,;
< error
| 42601 syntax error at or near "B'01;10'"
| position: 20
> exec
| BEGIN
< ok 0 -
> exec
| SELECT xd$n$m AS x B'01;10' + $$/*;$$ 1;
< error
| 42601 syntax error at or near "B'01;10'"
| position: 20
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT xd$n$m AS x B'01;10' + $$/*;$$ 1;
< error
| 42601 syntax error at or near "B'01;10'"
| position: 20
> exec
| , /*/*/*);(*/*/*/ 1 $$);(select 1;$x$;$$ SELECT;
| SELECT /*$x$;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*/*);(*/*/*/ 1 $$);(select 1;$x$;$$ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*/*);(*/*/*/ 1 $$);(select 1;$x$;$$ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) 1 e'\\\n;' col;COPY rmf FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) 1 e'\\\n;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) 1 e'\\\n;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col $__nz$select 1;/*;$__nx$$__nz$ , b'01;10' col B'01;10' ,;1 (SELECT 1; SELECT 2) + x'01;10' 1 U&'--;' FROM t WHERE a =;
| , 'a;b;;--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $__nz$select 1;/*;$__nx$$__nz$ , b'01;10' col B'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $__nz$select 1;/*;$__nx$$__nz$ , b'01;10' col B'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY olg FROM stdin;
| a;b
| c;d;e
| \.
| AS x --; $x$;
|  SELECT /*/*$;$END ;$x$;*/*/ ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY olg FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY olg FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT U&';;/*;select 1;' SELECT e'\\\';' 1 ");($x$;end;" ,;
| 1 b'01;10' , "/*;" FROM t WHERE a = 'select 1;select 1;;;\' ,;
| COPY gmm FROM stdin;
| a;b
| 1	\N
| \.
|
| COPY euc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "e'\\\\\\';'"
| position: 34
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&';;/*;select 1;' SELECT e'\\\';' 1 ");($x$;end;" ,;
< error
| 42601 syntax error at or near "e'\\\\\\';'"
| position: 34
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&';;/*;select 1;' SELECT e'\\\';' 1 ");($x$;end;" ,;
< error
| 42601 syntax error at or near "e'\\\\\\';'"
| position: 34
> exec
| COPY rel FROM stdin;
| a;b
| c;d;e
| \.
|  -- trailing note;
| , cg$r$e , /*/*/*end;*/*/*/ col wf$n$h ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY rel FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY rel FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , fq$w$t +;
| COPY djy FROM stdin;
| a;b
| a;b
| \.
|
| COPY jzu FROM stdin;
| 1	\N
| 1	\N
| \.
|
| FROM t WHERE a = --end;
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , fq$w$t +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , fq$w$t +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 B'01;10' 1; FROM t WHERE a = --select 1;;a;b
|  SELECT tk$t$z + CREATE FUNCTION qwhu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$;$);('; SELECT 1; END +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 B'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 B'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , --;;
|  + ";$x$;$x$;" AS x ';--;select 1;' ,; SELECT CREATE FUNCTION byej() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(--;END ;'; END ,;
| SELECT b'01;10' + /*/*end;*/*/ FROM t WHERE a = --;;END ;
|  1;
| + md$e$i 1 '; );(a;b' + x'01;10' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --;;
|  + ";$x$;$x$;" AS x ';--;select 1;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --;;
|  + ";$x$;$x$;" AS x ';--;select 1;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col $$);($$ FROM t WHERE a = CREATE FUNCTION bjud() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;a;bselect 1;'; END ,; COPY xcc FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| FROM t WHERE a = B'01;10' + /*/*/*);($x$;;*/*/*/ +;
| AS x E'\\\n;' AS x E''';' AS x hy$_$u col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$);($$ FROM t WHERE a = CREATE FUNCTION bjud() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;a;bselect 1;'; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$);($$ FROM t WHERE a = CREATE FUNCTION bjud() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;a;bselect 1;'; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x --/*;; ; 
|  AS x ki$w$b , --);(;;
|  +;
| 1 U&'select 1;;;a;b' AS x CREATE FUNCTION epnd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END 1 B'01;10' SELECT; -- trailing note;
| FROM t WHERE a = '$x$;\' , CREATE FUNCTION zpji() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;END ;END ;'; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --/*;; ; 
|  AS x ki$w$b , --);(;;
|  +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --/*;; ; 
|  AS x ki$w$b , --);(;;
|  +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT lq$w$g SELECT;+ "end;;;" 1 E''';\n;' ,; , --/*;select 1;
|  SELECT (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT lq$w$g SELECT;
< error
| 42703 column "lq$w$g" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| + "end;;;" 1 E''';\n;' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + U&'end;END ;end;' + (SELECT 1; SELECT 2) 1;
| , (SELECT 1; SELECT 2) 1 CREATE FUNCTION ickt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'end;END ;end;' + (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'end;END ;end;' + (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) col $$END ;; --;$$ 1 E'\';' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) col $$END ;; --;$$ 1 E'\';' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) col $$END ;; --;$$ 1 E'\';' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col "a;b$;$; " + "END ;);(" 1 E'\n;\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "a;b$;$; " + "END ;);(" 1 E'\n;\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "a;b$;$; " + "END ;);(" 1 E'\n;\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY xkf FROM stdin;
| a;b
| a;b
| \.
|
| + x'01;10' SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xkf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xkf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY hfu FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hfu FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hfu FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 CREATE FUNCTION tmet() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;;; '; END AS x (SELECT 1; SELECT 2) AS x;AS x CREATE FUNCTION _ekw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END FROM t WHERE a =; -- trailing note;
| SELECT "/*;select 1;$x$;" AS x E'\';\';'';' , ua$b$t 1; , e'\\\n;\n;' AS x vf$q$g FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION tmet() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;;; '; END AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION tmet() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;;; '; END AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , "end;" SELECT E'\\' AS x (SELECT 1; SELECT 2) AS x;
| FROM t WHERE a = B'01;10' , U&'; ' + U&';;;$;$' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "end;" SELECT E'\\' AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "end;" SELECT E'\\' AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col --$x$;;;;
|  SELECT "; " col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --$x$;;;;
|  SELECT "; " col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --$x$;;;;
|  SELECT "; " col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY _cq FROM stdin;
| 1	\N
| 1	\N
| \.
|  COPY abi FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY jzu FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _cq FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _cq FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY ymm FROM stdin;
| c;d;e
| a;b
| \.
|
| FROM t WHERE a = xn$r$q col (SELECT 1; SELECT 2) col; AS x U&'end;;;/*;' FROM t WHERE a = U&'a;b' AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ymm FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ymm FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x X'01;10' SELECT --$;$--;
|  1 CREATE FUNCTION iezg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;'; SELECT 1; END +;
| AS x co$e$w , e''';\';'';' col; -- trailing note;
| FROM t WHERE a = x'01;10' 1 vl$f$_ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x X'01;10' SELECT --$;$--;
|  1 CREATE FUNCTION iezg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x X'01;10' SELECT --$;$--;
|  1 CREATE FUNCTION iezg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT --;;--;;
|  col $$;a;b$$ 1;SELECT X'01;10' SELECT E'\';' ,;
< error
| 42601 syntax error at or near "1"
| position: 31
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --;;--;;
|  col $$;a;b$$ 1;
< error
| 42601 syntax error at or near "1"
| position: 31
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --;;--;;
|  col $$;a;b$$ 1;
< error
| 42601 syntax error at or near "1"
| position: 31
> exec
| , ');(--;a;b' 1 e'\\'';\n;' col $$$x$;;;$$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ');(--;a;b' 1 e'\\'';\n;' col $$$x$;;;$$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ');(--;a;b' 1 e'\\'';\n;' col $$$x$;;;$$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 B'01;10' FROM t WHERE a =;COPY xtl FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY ags FROM stdin;
| a;b
| a;b
| \.
|
| 1 (SELECT 1; SELECT 2) col;COPY ouw FROM stdin;
| a;b
| a;b
| \.
|
| COPY _mc FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ags FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ags FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = "; ;;" 1;AS x (SELECT 1; SELECT 2) FROM t WHERE a = "select 1;end;$x$;" , CREATE FUNCTION ndhr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; END AS x;
| + $$END ;--;$$ col X'01;10' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "; ;;" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "; ;;" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col 'a;b\' FROM t WHERE a = U&'--;;;;' FROM t WHERE a = $_on$END ;$_ox$$_on$ 1; -- trailing note;
| + CREATE FUNCTION hdbg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$;$'; SELECT 1; END col $$;;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'a;b\' FROM t WHERE a = U&'--;;;;' FROM t WHERE a = $_on$END ;$_ox$$_on$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'a;b\' FROM t WHERE a = U&'--;;;;' FROM t WHERE a = $_on$END ;$_ox$$_on$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY nzx FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 /*/*$;$a;bEND ;*/*/ +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY nzx FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY nzx FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY xzd FROM stdin;
| c;d;e
| 1	\N
| \.
| COPY bvk FROM stdin;
| a;b
| 1	\N
| \.
|
| FROM t WHERE a = ----;;
|  SELECT (SELECT 1; SELECT 2) ,; -- trailing note;
| col b'01;10' , e'\n;\n;\\' col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xzd FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xzd FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 ut$r$o AS x;
| FROM t WHERE a = se$y$l ,;
| + (SELECT 1; SELECT 2) AS x e''';\\' 1;
| + /*/*/*$;$; );(*/*/*/ SELECT /*;end;*/ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ut$r$o AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ut$r$o AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 /*/*/*--;END ;--;*/*/*/ col "end;a;bselect 1;" AS x $gxhj$select 1;;;$;$$gxhx$$gxhj$ FROM t WHERE a =;
| , _r$q$w + x'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*--;END ;--;*/*/*/ col "end;a;bselect 1;" AS x $gxhj$select 1;;;$;$$gxhx$$gxhj$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*--;END ;--;*/*/*/ col "end;a;bselect 1;" AS x $gxhj$select 1;;;$;$$gxhx$$gxhj$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY ugu FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ugu FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ugu FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = /*a;b/ *;*/ + (SELECT 1; SELECT 2) SELECT x'01;10' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*a;b/ *;*/ + (SELECT 1; SELECT 2) SELECT x'01;10' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*a;b/ *;*/ + (SELECT 1; SELECT 2) SELECT x'01;10' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT wf$b$_ SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "("
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT wf$b$_ SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "("
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT wf$b$_ SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "("
| position: 22
> exec
| col --);(
|  SELECT U&'END ;' ,;1 oa$b$r ,;
| AS x /*/*end;;*/*/ , CREATE FUNCTION bsxi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --);(
|  SELECT U&'END ;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --);(
|  SELECT U&'END ;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY gut FROM stdin;
| a;b
| 1	\N
| \.
|  SELECT E'\\' SELECT;
| FROM t WHERE a = --$;$/*;select 1;
|  +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gut FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gut FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY iwz FROM stdin;
| c;d;e
| a;b
| \.
|
| SELECT $v_p$);($v_p$ 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY iwz FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY iwz FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = ";" , E'\';'';' SELECT 'select 1;/*;\' SELECT; COPY y_n FROM stdin;
| 1	\N
| a;b
| \.
|
| 1 CREATE FUNCTION hgli() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END FROM t WHERE a = --; select 1;
|  SELECT;
| FROM t WHERE a = "--;$;$" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ";" , E'\';'';' SELECT 'select 1;/*;\' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ";" , E'\';'';' SELECT 'select 1;/*;\' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 ip$p$_ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ip$p$_ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ip$p$_ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY xtw FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| SELECT pv$n$y 1 e'\';\\' col; , te$_$y SELECT '--;$;$\' SELECT "end;END ;end;" FROM t WHERE a =;COPY rin FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xtw FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xtw FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT 'END ;; ' ,;
| SELECT sb$b$b AS x X'01;10' SELECT --END ;$;$
|  ,;
| FROM t WHERE a = "end;/*;;" AS x;
< error
| 42601 syntax error at or near ";"
| position: 19
> exec
| BEGIN
< ok 0 -
> exec
| SELECT 'END ;; ' ,;
< error
| 42601 syntax error at or near ";"
| position: 19
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT 'END ;; ' ,;
< error
| 42601 syntax error at or near ";"
| position: 19
> exec
| SELECT (SELECT 1; SELECT 2) col;COPY vzb FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| FROM t WHERE a = $$end;END ;$;$$$ 1 CREATE FUNCTION tqj_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;'; SELECT 1; END AS x --;$x$;
|  AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$end;END ;$;$$$ 1 CREATE FUNCTION tqj_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;'; SELECT 1; END AS x --;$x$;
|  AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$end;END ;$;$$$ 1 CREATE FUNCTION tqj_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;'; SELECT 1; END AS x --;$x$;
|  AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + CREATE FUNCTION xs_l() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END + 'end;; ' 1 ";;" ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION xs_l() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END + 'end;; ' 1 ";;" ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION xs_l() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END + 'end;; ' 1 ";;" ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x qh$_$_ AS x (SELECT 1; SELECT 2) col;
| FROM t WHERE a = --;--;select 1;
|  FROM t WHERE a =;
| COPY lpj FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x qh$_$_ AS x (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x qh$_$_ AS x (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT $uthl$END ;$uthl$ col ';;$;$' AS x;
< error
| 42601 syntax error at or near "';;$;$'"
| position: 30
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $uthl$END ;$uthl$ col ';;$;$' AS x;
< error
| 42601 syntax error at or near "';;$;$'"
| position: 30
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $uthl$END ;$uthl$ col ';;$;$' AS x;
< error
| 42601 syntax error at or near "';;$;$'"
| position: 30
> exec
| COPY vxd FROM stdin;
| c;d;e
| a;b
| \.
| 1 CREATE FUNCTION kvzo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(a;b$x$;'; SELECT 1; END AS x $$--;--;;$$ + /*/*/*END ;);(*/*/*/ 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vxd FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vxd FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY dad FROM stdin;
| a;b
| a;b
| \.
|  COPY nrk FROM stdin;
| a;b
| c;d;e
| \.
| + CREATE FUNCTION wumg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END FROM t WHERE a = 'select 1;select 1;; \' , "select 1;;" 1;AS x X'01;10' 1 ql$v$t 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dad FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dad FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col 'select 1;;;\' 1 $$$;$$$ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'select 1;;;\' 1 $$$;$$$ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'select 1;;;\' 1 $$$;$$$ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY kva FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| + /*/*/ *;;;*/*/ FROM t WHERE a = /*/*a;b*/*/ SELECT e'\';\\\n;' 1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kva FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kva FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY xat FROM stdin;
| 1	\N
| c;d;e
| \.
|  COPY pby FROM stdin;
| 1	\N
| 1	\N
| \.
|
| COPY ddd FROM stdin;
| a;b
| 1	\N
| \.
|
| + b'01;10' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xat FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xat FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY bqp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bqp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bqp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT '$x$;\' , E''';\';' SELECT U&'select 1;' 1;
| AS x $$/*;$$ 1 U&'; ' 1;
< error
| 42601 syntax error at or near "U&'select 1;'"
| position: 35
> exec
| BEGIN
< ok 0 -
> exec
| SELECT '$x$;\' , E''';\';' SELECT U&'select 1;' 1;
< error
| 42601 syntax error at or near "U&'select 1;'"
| position: 35
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT '$x$;\' , E''';\';' SELECT U&'select 1;' 1;
< error
| 42601 syntax error at or near "U&'select 1;'"
| position: 35
> exec
| COPY ykf FROM stdin;
| 1	\N
| c;d;e
| \.
|
| , CREATE FUNCTION bzth() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; end;'; END col om$k$y SELECT /*/*/*a;b*/*/*/ SELECT;
| , --$;$);(
|  AS x CREATE FUNCTION vgpf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;END ;/*;'; SELECT 1; END FROM t WHERE a = b'01;10' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ykf FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ykf FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col CREATE FUNCTION gtnz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;a;b'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION gtnz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;a;b'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION gtnz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;a;b'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY yrv FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yrv FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yrv FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 --end;$x$;
|  AS x U&'; ;' , (SELECT 1; SELECT 2) 1;
| AS x "END ;; " 1 "end;" AS x;
| FROM t WHERE a = ";;END ;" col /*/*;;a;bselect 1;*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --end;$x$;
|  AS x U&'; ;' , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --end;$x$;
|  AS x U&'; ;' , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + /*; END ;*/ SELECT U&'--;' , b'01;10' AS x;
| SELECT kr$c$j , U&';;;' col "$;$" 1; SELECT "; ;;" FROM t WHERE a = ";" SELECT --select 1;
|  FROM t WHERE a =;
| COPY ont FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*; END ;*/ SELECT U&'--;' , b'01;10' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*; END ;*/ SELECT U&'--;' , b'01;10' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY i_f FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY i_f FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY i_f FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 e'\\'';' SELECT --a;b$x$;
|  col B'01;10' FROM t WHERE a =;
| COPY wmj FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\\'';' SELECT --a;b$x$;
|  col B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\\'';' SELECT --a;b$x$;
|  col B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + yv$u$l AS x;
| COPY ys_ FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + yv$u$l AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + yv$u$l AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col ly$a$y + al$r$_ + (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ly$a$y + al$r$_ + (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ly$a$y + al$r$_ + (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY gtx FROM stdin;
| 1	\N
| c;d;e
| \.
|  FROM t WHERE a = ';END ;end;' + ef$o$h AS x x'01;10' AS x; -- trailing note;
| AS x CREATE FUNCTION _fwd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;bselect 1;'; SELECT 1; END FROM t WHERE a = "$x$;end;select 1;" col;
| AS x CREATE FUNCTION ogbo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;END ;'; SELECT 1; END , ');(;;; ' + ';' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gtx FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gtx FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = --/*;
|  + sc$z$e FROM t WHERE a = ao$z$r 1;
| + $$a;bend;$$ , /*/*/*);(*/*/*/ , $ptkq$end;$ptkq$ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --/*;
|  + sc$z$e FROM t WHERE a = ao$z$r 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --/*;
|  + sc$z$e FROM t WHERE a = ao$z$r 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 --$;$$;$select 1;
|  AS x "; " AS x;
| FROM t WHERE a = "$;$--;" FROM t WHERE a = U&'$x$;;;;' +;
| + af$u$n SELECT 'END ;\' + "--;a;b" SELECT;
| FROM t WHERE a = B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --$;$$;$select 1;
|  AS x "; " AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --$;$$;$select 1;
|  AS x "; " AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY njj FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| + U&'/*;$;$);(' AS x zi$a$d col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY njj FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY njj FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , /*/*;$;$;*/*/ col 'END ;/*;$;$\' +; -- trailing note;
| COPY xgu FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY rvo FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY zcc FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*;$;$;*/*/ col 'END ;/*;$;$\' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*;$;$;*/*/ col 'END ;/*;$;$\' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + U&'select 1;' ,;
| AS x --);($;$
|  FROM t WHERE a =;COPY i_z FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'select 1;' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'select 1;' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY qqs FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY aoc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qqs FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qqs FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT x'01;10' FROM t WHERE a =;
| FROM t WHERE a = /*END ;$x$;a;b*/ col r_$z$a , b'01;10' 1;
| AS x $ourv$END ;;END ;$ourv$ col X'01;10' ,;
| , U&'$x$;/*;--;' SELECT ';;a;b);(' AS x;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| COPY _bl FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| + --END ;--;
|  FROM t WHERE a = '; END ;' AS x; COPY cqu FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 E'\n;\n;' FROM t WHERE a = /*$;$*/ ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _bl FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _bl FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY qlf FROM stdin;
| a;b
| a;b
| \.
|
| 1 e'\\\\' ,;
| + iu$h$_ + U&'; ;$;$' ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qlf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qlf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT "; --;" FROM t WHERE a = --select 1;; 
|  AS x;
| + B'01;10' AS x "$x$;; END ;" + CREATE FUNCTION hows() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;'; END col;
| col "end;--;" 1 '; ;;\' + ii$d$q ,;AS x e'\\'';' +;
< error
| 42601 syntax error at or near "AS"
| position: 48
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "; --;" FROM t WHERE a = --select 1;; 
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 48
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "; --;" FROM t WHERE a = --select 1;; 
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 48
> exec
| FROM t WHERE a = /*end;; ;;*/ 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*end;; ;;*/ 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*end;; ;;*/ 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY uu_ FROM stdin;
| c;d;e
| a;b
| \.
|
| , U&'select 1;' , ");(end;" SELECT X'01;10' AS x;
| COPY yoi FROM stdin;
| a;b
| c;d;e
| \.
| FROM t WHERE a = CREATE FUNCTION ibof() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY uu_ FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY uu_ FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = ";;$x$;--;" col $u$select 1;END ;$x$$u$ AS x "/*;--;a;b" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ";;$x$;--;" col $u$select 1;END ;$x$$u$ AS x "/*;--;a;b" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ";;$x$;--;" col $u$select 1;END ;$x$$u$ AS x "/*;--;a;b" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = U&');(select 1;' ,;
| FROM t WHERE a = /*);(*/ SELECT U&'end;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&');(select 1;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&');(select 1;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col "$;$/*;" , X'01;10' FROM t WHERE a = CREATE FUNCTION acbn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;END ;'; END 1; -- trailing note;
| AS x E''';' FROM t WHERE a = (SELECT 1; SELECT 2) AS x 'a;b\' col;
| COPY lxg FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "$;$/*;" , X'01;10' FROM t WHERE a = CREATE FUNCTION acbn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;END ;'; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "$;$/*;" , X'01;10' FROM t WHERE a = CREATE FUNCTION acbn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;END ;'; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 ux$q$o SELECT (SELECT 1; SELECT 2) AS x CREATE FUNCTION xbuu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END ,;
| col 'a;b;;' FROM t WHERE a = B'01;10' SELECT /*/*; / *;END ;*/*/ ,; -- trailing note;
| + (SELECT 1; SELECT 2) col;
| FROM t WHERE a = "--;);(end;" SELECT --);(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ux$q$o SELECT (SELECT 1; SELECT 2) AS x CREATE FUNCTION xbuu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ux$q$o SELECT (SELECT 1; SELECT 2) AS x CREATE FUNCTION xbuu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + tg$k$d AS x;
| 1 $nm$a;b$nx$$nm$ AS x (SELECT 1; SELECT 2) , '$;$' ,;
| FROM t WHERE a = --$x$;
|  AS x U&'end;--;' AS x B'01;10' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + tg$k$d AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + tg$k$d AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY sry FROM stdin;
| c;d;e
| c;d;e
| \.
| SELECT U&'/*;;;' 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sry FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sry FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 e''';\';' FROM t WHERE a = $$);();($x$;$$ col $$END ;$$ AS x; col B'01;10' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e''';\';' FROM t WHERE a = $$);();($x$;$$ col $$END ;$$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e''';\';' FROM t WHERE a = $$);();($x$;$$ col $$END ;$$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , --select 1;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --select 1;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --select 1;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x E'\\' FROM t WHERE a = "; select 1;;;" col; COPY yvu FROM stdin;
| a;b
| 1	\N
| \.
|
| COPY vzx FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x E'\\' FROM t WHERE a = "; select 1;;;" col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x E'\\' FROM t WHERE a = "; select 1;;;" col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + /*/*select 1;);(*/*/ col U&');(; ' FROM t WHERE a = $b$$x$;a;b;$x$$b$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*select 1;);(*/*/ col U&');(; ' FROM t WHERE a = $b$$x$;a;b;$x$$b$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*select 1;);(*/*/ col U&');(; ' FROM t WHERE a = $b$$x$;a;b;$x$$b$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , (SELECT 1; SELECT 2) + '$x$;END ;' + /*/*; ;*/*/ col; SELECT '; select 1;' AS x ';;$;$END ;' , gc$k$c col; , CREATE FUNCTION fqxb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$x$;'; END col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) + '$x$;END ;' + /*/*; ;*/*/ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) + '$x$;END ;' + /*/*; ;*/*/ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col $nir$end;$nix$$nir$ , (SELECT 1; SELECT 2) + $$a;b);($$ AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $nir$end;$nix$$nir$ , (SELECT 1; SELECT 2) + $$a;b);($$ AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $nir$end;$nix$$nir$ , (SELECT 1; SELECT 2) + $$a;b);($$ AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY mzl FROM stdin;
| a;b
| c;d;e
| \.
|  FROM t WHERE a = $er$$x$;$er$ col E'\';' SELECT; -- trailing note;
| col "/*;; " AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mzl FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mzl FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , '$;$;;$;$' 1;, E'\';' + ----;;; 
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , '$;$;;$;$' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , '$;$;;$;$' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 $hoyb$$x$;;$hoyx$$hoyb$ + --;;; 
|  SELECT E'\';\n;'';' ,; -- trailing note;
| COPY mnx FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $hoyb$$x$;;$hoyx$$hoyb$ + --;;; 
|  SELECT E'\';\n;'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $hoyb$$x$;;$hoyx$$hoyb$ + --;;; 
|  SELECT E'\';\n;'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY sxd FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sxd FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sxd FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x /*; END ;*/ SELECT (SELECT 1; SELECT 2) SELECT --);(
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*; END ;*/ SELECT (SELECT 1; SELECT 2) SELECT --);(
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*; END ;*/ SELECT (SELECT 1; SELECT 2) SELECT --);(
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + $$END ;; $$ col;
| COPY mwt FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| AS x "$;$/*;" SELECT; -- trailing note;
| FROM t WHERE a = /*/*;*/*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $$END ;; $$ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $$END ;; $$ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = /*/*/*end;*/*/*/ 1 /*/*/*;;select 1;*/*/*/ FROM t WHERE a = CREATE FUNCTION o_eg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END FROM t WHERE a =;
| COPY qlb FROM stdin;
| a;b
| c;d;e
| \.
|
| FROM t WHERE a = B'01;10' FROM t WHERE a = /*end;*/ FROM t WHERE a = "end;$x$;$;$" 1;+ CREATE FUNCTION ymfp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END , /*;;$;$$x$;*/ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*/*/*end;*/*/*/ 1 /*/*/*;;select 1;*/*/*/ FROM t WHERE a = CREATE FUNCTION o_eg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*/*/*end;*/*/*/ 1 /*/*/*;;select 1;*/*/*/ FROM t WHERE a = CREATE FUNCTION o_eg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT bt$f$s , e'\n;\';' col B'01;10' col;
< error
| 42601 syntax error at or near "B'01;10'"
| position: 31
> exec
| BEGIN
< ok 0 -
> exec
| SELECT bt$f$s , e'\n;\';' col B'01;10' col;
< error
| 42601 syntax error at or near "B'01;10'"
| position: 31
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT bt$f$s , e'\n;\';' col B'01;10' col;
< error
| 42601 syntax error at or near "B'01;10'"
| position: 31
> exec
| SELECT x'01;10' AS x --;;;
|  ,; -- trailing note;
| COPY idu FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| 1 /*$x$;--;; */ , x'01;10' ,;
< error
| 42601 syntax error at or near ";"
| position: 30
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' AS x --;;;
|  ,;
< error
| 42601 syntax error at or near ";"
| position: 30
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT x'01;10' AS x --;;;
|  ,;
< error
| 42601 syntax error at or near ";"
| position: 30
> exec
| , (SELECT 1; SELECT 2) AS x "a;b;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) AS x "a;b;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) AS x "a;b;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT e'\';' SELECT;
| col (SELECT 1; SELECT 2) col 'select 1;a;ba;b\' col --);(end;$x$;
|  +;
| col jv$s$g , U&'end;select 1;--;' SELECT X'01;10' FROM t WHERE a =;
| AS x CREATE FUNCTION axtm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$);('; END 1 E'\n;' , $dbi$$x$;);();($dbi$ col;
< error
| 42601 syntax error at or near "col"
| position: 23
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\';' SELECT;
< ok 1 -
> exec
|
| col (SELECT 1; SELECT 2) col 'select 1;a;ba;b\' col --);(end;$x$;
|  +;
< error
| 42601 syntax error at or near "col"
| position: 2
> exec
| ROLLBACK
< ok 0 -
> exec
|
| col (SELECT 1; SELECT 2) col 'select 1;a;ba;b\' col --);(end;$x$;
|  +;
< error
| 42601 syntax error at or near "col"
| position: 2
> exec
| AS x ph$p$k AS x CREATE FUNCTION omgj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ph$p$k AS x CREATE FUNCTION omgj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ph$p$k AS x CREATE FUNCTION omgj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x /*/*$x$;--;select 1;*/*/ SELECT U&'$;$' FROM t WHERE a = X'01;10' AS x; -- trailing note;
| COPY ypn FROM stdin;
| a;b
| a;b
| \.
|
| COPY wyf FROM stdin;
| c;d;e
| c;d;e
| \.
|
| , --;
|  1 E''';\\' 1 $i$; a;b$i$ ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*$x$;--;select 1;*/*/ SELECT U&'$;$' FROM t WHERE a = X'01;10' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*$x$;--;select 1;*/*/ SELECT U&'$;$' FROM t WHERE a = X'01;10' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x U&'select 1;' col;
| , /*/*/*end;select 1;a;b*/*/*/ col;
| SELECT (SELECT 1; SELECT 2) AS x;
| AS x $$;;/*;$$ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'select 1;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'select 1;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) SELECT x'01;10' AS x;
| col e'\';' SELECT U&'$;$END ;' + $$;;;);($$ ,; col $$; /*;);($$ +;
| + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) SELECT x'01;10' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) SELECT x'01;10' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT U&');(' ,;
| COPY sbm FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&');(' ,;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&');(' ,;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| col ');(' SELECT rj$i$i FROM t WHERE a = CREATE FUNCTION xayw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END ,; -- trailing note;
| col b'01;10' + $vsc$);($;$;;$vsc$ FROM t WHERE a =;, ");(; ;;" 1;
| COPY guo FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ');(' SELECT rj$i$i FROM t WHERE a = CREATE FUNCTION xayw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ');(' SELECT rj$i$i FROM t WHERE a = CREATE FUNCTION xayw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col 'end;a;bEND ;' 1 (SELECT 1; SELECT 2) FROM t WHERE a = hd$b$r FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'end;a;bEND ;' 1 (SELECT 1; SELECT 2) FROM t WHERE a = hd$b$r FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'end;a;bEND ;' 1 (SELECT 1; SELECT 2) FROM t WHERE a = hd$b$r FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + $eokm$--;; $eokx$$eokm$ FROM t WHERE a = $$);(select 1;);($$ ,;
| FROM t WHERE a = CREATE FUNCTION _oyu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';; /*;'; SELECT 1; END FROM t WHERE a = --select 1;--;end;
|  +;
| SELECT $$$x$;select 1;$$ + b'01;10' 1;
| + "$;$END ;" FROM t WHERE a = CREATE FUNCTION ruxd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;'; SELECT 1; END 1 U&'/*;; ' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $eokm$--;; $eokx$$eokm$ FROM t WHERE a = $$);(select 1;);($$ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $eokm$--;; $eokx$$eokm$ FROM t WHERE a = $$);(select 1;);($$ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY xyv FROM stdin;
| 1	\N
| a;b
| \.
| AS x $$;;END ;$$ +;
| 1 e''';\n;' + CREATE FUNCTION _lu_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;'; END AS x U&'; ;;' +;
| + /*/*/*--;*/*/*/ col x'01;10' , E''';\';' 1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xyv FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xyv FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col x'01;10' 1;
| SELECT $dllt$end;--;$dllt$ 1 CREATE FUNCTION yull() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;);('; SELECT 1; END col;
| + e'\n;' 1 e'\n;\';\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col x'01;10' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col x'01;10' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY zba FROM stdin;
| 1	\N
| c;d;e
| \.
|
| COPY rtm FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| col /*/*;;/ *;end;*/*/ AS x gz$p$b SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zba FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zba FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY syt FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| col X'01;10' ,;
| FROM t WHERE a = CREATE FUNCTION nvsd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;'; SELECT 1; END AS x X'01;10' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY syt FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY syt FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , (SELECT 1; SELECT 2) ,; -- trailing note;
| + U&'/*;' AS x B'01;10' FROM t WHERE a = bh$k$j col;
| FROM t WHERE a = $ior$select 1;;;$ior$ col;
| AS x CREATE FUNCTION jdib() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;);(select 1;'; END FROM t WHERE a = (SELECT 1; SELECT 2) , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 jf$o$l SELECT;
| AS x E''';' , U&'end;' FROM t WHERE a = --;end;
|  AS x; -- trailing note;
| + x'01;10' AS x;
| col (SELECT 1; SELECT 2) AS x xq$u$x col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 jf$o$l SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 jf$o$l SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY gje FROM stdin;
| c;d;e
| a;b
| \.
|
| 1 --;;
|  FROM t WHERE a = '/*;select 1;' SELECT ";" col;
| COPY hci FROM stdin;
| 1	\N
| a;b
| \.
|
| AS x --$;$
|  col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gje FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gje FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY jlf FROM stdin;
| c;d;e
| a;b
| \.
|  COPY lvx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jlf FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jlf FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col $$end;/*;$$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$end;/*;$$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$end;/*;$$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x jx$t$l AS x CREATE FUNCTION hbhu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; /*;a;b'; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x jx$t$l AS x CREATE FUNCTION hbhu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; /*;a;b'; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x jx$t$l AS x CREATE FUNCTION hbhu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; /*;a;b'; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY mn_ FROM stdin;
| c;d;e
| a;b
| \.
|
| 1 ';;);();(' FROM t WHERE a = CREATE FUNCTION oegn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mn_ FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mn_ FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + /*/*/ *;;;$x$;*/*/ 1 /*/*/*$;$*/*/*/ SELECT;
| AS x B'01;10' FROM t WHERE a =; , x'01;10' SELECT;, b'01;10' AS x xm$r$l FROM t WHERE a = 'a;b/*;\' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*/ *;;;$x$;*/*/ 1 /*/*/*$;$*/*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*/ *;;;$x$;*/*/ 1 /*/*/*$;$*/*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col (SELECT 1; SELECT 2) SELECT; + (SELECT 1; SELECT 2) col U&'$;$' FROM t WHERE a = $$select 1;$$ ,;
| 1 (SELECT 1; SELECT 2) FROM t WHERE a = "a;bend;;;" FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a =;
| COPY jtl FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col U&'select 1;$;$' , "/*;" + e''';\\\';' 1;
| FROM t WHERE a = CREATE FUNCTION ccnv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END + /*/*/*;$x$;);(*/*/*/ FROM t WHERE a =;
| FROM t WHERE a = $cxk$;$x$;$cxk$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'select 1;$;$' , "/*;" + e''';\\\';' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'select 1;$;$' , "/*;" + e''';\\\';' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT "a;b" , CREATE FUNCTION phdu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; '; SELECT 1; END col;
| , te$z$l +;
| COPY hgn FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "CREATE"
| position: 16
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "a;b" , CREATE FUNCTION phdu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; '; SELECT 1; END col;
< error
| 42601 syntax error at or near "CREATE"
| position: 16
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "a;b" , CREATE FUNCTION phdu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; '; SELECT 1; END col;
< error
| 42601 syntax error at or near "CREATE"
| position: 16
> exec
| AS x CREATE FUNCTION _rnn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;--;end;'; SELECT 1; END , E'\';' AS x; -- trailing note;
| COPY lka FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION _rnn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;--;end;'; SELECT 1; END , E'\';' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION _rnn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;--;end;'; SELECT 1; END , E'\';' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY rin FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY rin FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY rin FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col 'a;b;select 1;\' col b'01;10' ,; COPY l_m FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'a;b;select 1;\' col b'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'a;b;select 1;\' col b'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x CREATE FUNCTION ckga() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$;$'; SELECT 1; END + U&'a;b$x$;--;' SELECT U&'; ' FROM t WHERE a =;
| SELECT e'\';\\' , --a;b
|  SELECT /*--;*/ col;
| col b'01;10' FROM t WHERE a =;
| COPY snb FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION ckga() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$;$'; SELECT 1; END + U&'a;b$x$;--;' SELECT U&'; ' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION ckga() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$;$'; SELECT 1; END + U&'a;b$x$;--;' SELECT U&'; ' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY dek FROM stdin;
| a;b
| 1	\N
| \.
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dek FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dek FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT CREATE FUNCTION cohl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END SELECT;
| COPY mpx FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION cohl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION cohl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| AS x (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x X'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x X'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x X'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x CREATE FUNCTION ccvh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;;'; END ,;FROM t WHERE a = e'\n;\';' FROM t WHERE a = $$END ;);($$ SELECT --;;
|  FROM t WHERE a =;
| + /*;*/ +;
| COPY myi FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION ccvh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;;'; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION ccvh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;;'; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = $cpy$end;--;$cpy$ SELECT;
| col wt$u$e + ";;; END ;" 1 fn$h$c FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $cpy$end;--;$cpy$ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $cpy$end;--;$cpy$ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = e'\\\\\';' FROM t WHERE a = cf$o$v +;
| COPY rag FROM stdin;
| a;b
| c;d;e
| \.
| 1 (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = e'\\\\\';' FROM t WHERE a = cf$o$v +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = e'\\\\\';' FROM t WHERE a = cf$o$v +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT E'\n;\';\n;' +;
| 1 fh$o$b SELECT --; a;b
|  AS x (SELECT 1; SELECT 2) col;FROM t WHERE a = /*/*/*select 1;select 1;--;*/*/*/ col CREATE FUNCTION jpcr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END SELECT;
| SELECT (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ";"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT E'\n;\';\n;' +;
< error
| 42601 syntax error at or near ";"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT E'\n;\';\n;' +;
< error
| 42601 syntax error at or near ";"
| position: 22
> exec
| SELECT b'01;10' FROM t WHERE a =;
| col 'END ;;;' SELECT ";;" ,;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| BEGIN
< ok 0 -
> exec
| SELECT b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| col /*/*/*a;b--;*/*/*/ FROM t WHERE a = $$; $$ , E''';' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*/*a;b--;*/*/*/ FROM t WHERE a = $$; $$ , E''';' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*/*a;b--;*/*/*/ FROM t WHERE a = $$; $$ , E''';' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY dhx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dhx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dhx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x '$;$);(--;\' 1 e'\n;'';' ,;col X'01;10' FROM t WHERE a = (SELECT 1; SELECT 2) SELECT ");(/*;$;$" FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x '$;$);(--;\' 1 e'\n;'';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x '$;$);(--;\' 1 e'\n;'';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + X'01;10' col B'01;10' SELECT;col U&');($x$;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + X'01;10' col B'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + X'01;10' col B'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + CREATE FUNCTION zhhu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END 1 'select 1;;/*;' col --select 1;;
|  FROM t WHERE a =;
| COPY _ky FROM stdin;
| c;d;e
| 1	\N
| \.
|
| SELECT /*/*/*select 1;*/*/*/ 1 '$;$' FROM t WHERE a = CREATE FUNCTION dert() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION zhhu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END 1 'select 1;;/*;' col --select 1;;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION zhhu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END 1 'select 1;;/*;' col --select 1;;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x U&'$;$' col '$;$' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'$;$' col '$;$' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'$;$' col '$;$' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY f_t FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY f_t FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY f_t FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY xvt FROM stdin;
| 1	\N
| a;b
| \.
|
| SELECT U&'/*;' 1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xvt FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xvt FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT vv$t$b AS x x'01;10' 1; FROM t WHERE a = CREATE FUNCTION qdgg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;;;'; SELECT 1; END col (SELECT 1; SELECT 2) +; -- trailing note;
| FROM t WHERE a = B'01;10' ,;
| 1 CREATE FUNCTION zd_p() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END FROM t WHERE a = b'01;10' + --;;
|  1;
< error
| 42601 syntax error at or near "x'01;10'"
| position: 20
> exec
| BEGIN
< ok 0 -
> exec
| SELECT vv$t$b AS x x'01;10' 1;
< error
| 42601 syntax error at or near "x'01;10'"
| position: 20
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT vv$t$b AS x x'01;10' 1;
< error
| 42601 syntax error at or near "x'01;10'"
| position: 20
> exec
| + e''';' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + e''';' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + e''';' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT CREATE FUNCTION dtvf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$;;;;'; END + U&');(end;$x$;' AS x; -- trailing note;
| , (SELECT 1; SELECT 2) col; -- trailing note;
| FROM t WHERE a = --;$;$$;$
|  SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION dtvf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$;;;;'; END + U&');(end;$x$;' AS x;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION dtvf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$;;;;'; END + U&');(end;$x$;' AS x;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| col dn$q$q AS x E'\';\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col dn$q$q AS x E'\';\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col dn$q$q AS x E'\';\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + CREATE FUNCTION owbj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b/*;'; SELECT 1; END FROM t WHERE a = E''';'';'';' +;
| , U&'$x$;END ;select 1;' 1;
| + (SELECT 1; SELECT 2) , mh$o$b 1 "END ;" col; -- trailing note;
| col ";;; " , U&'select 1;' + /*/*select 1;*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION owbj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b/*;'; SELECT 1; END FROM t WHERE a = E''';'';'';' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION owbj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b/*;'; SELECT 1; END FROM t WHERE a = E''';'';'';' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , y_$p$a col /*/*select 1;*/*/ +; + $$; --;$$ AS x /*/*/*);(select 1;*/*/*/ SELECT X'01;10' 1;+ /*/*/*END ;/ *;end;*/*/*/ ,;
| 1 '$x$;);($;$' , CREATE FUNCTION ghek() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';/*;; '; SELECT 1; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , y_$p$a col /*/*select 1;*/*/ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , y_$p$a col /*/*select 1;*/*/ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , --/*;a;b
|  AS x;
| FROM t WHERE a = ----;;; 
|  SELECT --$x$;
|  col "select 1;$x$;a;b" 1;
| + jh$v$f +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --/*;a;b
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --/*;a;b
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY zia FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zia FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zia FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x $$--;$$ +;
| SELECT "select 1;/*;select 1;" +;AS x (SELECT 1; SELECT 2) 1; SELECT CREATE FUNCTION rurc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;);(;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $$--;$$ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $$--;$$ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT e'\';\\\n;' col;
| COPY hrh FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| FROM t WHERE a = (SELECT 1; SELECT 2) 1 /*; ;;*/ FROM t WHERE a =;
| , ey$a$n AS x CREATE FUNCTION wgwu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; a;b'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "c"
| position: 46
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\';\\\n;' col;
< ok 1 -
> exec
|
| COPY hrh FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 23
> exec
| ROLLBACK
< ok 0 -
> exec
|
| COPY hrh FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 23
> exec
| 1 /*;;$x$;END ;*/ + (SELECT 1; SELECT 2) + $$end;--;$$ col;
| AS x $$--;select 1;;;$$ 1; -- trailing note;
| + "; );(;;" , fu$o$x AS x U&'$x$;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*;;$x$;END ;*/ + (SELECT 1; SELECT 2) + $$end;--;$$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*;;$x$;END ;*/ + (SELECT 1; SELECT 2) + $$end;--;$$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY irs FROM stdin;
| c;d;e
| c;d;e
| \.
|
| SELECT 'select 1;;;\' 1 /*/*);(a;b*/*/ +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY irs FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY irs FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY jzs FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jzs FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jzs FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x nh$a$w + X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x nh$a$w + X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x nh$a$w + X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY gii FROM stdin;
| a;b
| a;b
| \.
|
| FROM t WHERE a = --select 1;
|  SELECT b'01;10' ,;
| FROM t WHERE a = U&';);(' col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gii FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gii FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col U&'/*;END ;' col /*/ *;end;--;*/ col;
| + /*/*--;*/*/ +;
| + oi$i$o 1 /*/*--;end;end;*/*/ ,;
| + --end;);(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'/*;END ;' col /*/ *;end;--;*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'/*;END ;' col /*/ *;end;--;*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , 'a;b--;END ;\' col CREATE FUNCTION isuy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;select 1;'; SELECT 1; END +;
| col U&'/*;a;b' col $kmfs$;;$;$select 1;$kmfx$$kmfs$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , 'a;b--;END ;\' col CREATE FUNCTION isuy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;select 1;'; SELECT 1; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , 'a;b--;END ;\' col CREATE FUNCTION isuy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;select 1;'; SELECT 1; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = CREATE FUNCTION xtqa() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;select 1;select 1;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION xtqa() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;select 1;select 1;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION xtqa() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;select 1;select 1;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY sgv FROM stdin;
| c;d;e
| c;d;e
| \.
|  + X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sgv FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sgv FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + CREATE FUNCTION pwju() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; END FROM t WHERE a = e'\';' AS x '$;$; ' AS x;
| col (SELECT 1; SELECT 2) col;
| COPY uzi FROM stdin;
| c;d;e
| c;d;e
| \.
|
| col uz$r$w col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION pwju() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; END FROM t WHERE a = e'\';' AS x '$;$; ' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION pwju() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; END FROM t WHERE a = e'\';' AS x '$;$; ' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , "select 1;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "select 1;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "select 1;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , --$;$/*;
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --$;$/*;
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --$;$/*;
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 B'01;10' SELECT ';;--;' FROM t WHERE a =; COPY cpr FROM stdin;
| c;d;e
| 1	\N
| \.
|
| AS x $$/*;);(/*;$$ FROM t WHERE a = x'01;10' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 B'01;10' SELECT ';;--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 B'01;10' SELECT ';;--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , U&'select 1;; ' SELECT /*/*/*END ;/ *;*/*/*/ 1 "a;b; " 1; SELECT CREATE FUNCTION agkg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;);('; SELECT 1; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'select 1;; ' SELECT /*/*/*END ;/ *;*/*/*/ 1 "a;b; " 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'select 1;; ' SELECT /*/*/*END ;/ *;*/*/*/ 1 "a;b; " 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY gbn FROM stdin;
| 1	\N
| a;b
| \.
|
| , --END ;--;
|  ,;
| FROM t WHERE a = "--;);($;$" ,;
| col e'\n;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gbn FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gbn FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT (SELECT 1; SELECT 2) 1; COPY dlt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| SELECT B'01;10' 1; COPY vvv FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| COPY vyg FROM stdin;
| c;d;e
| 1	\N
| \.
|
| col $o$; ;);($x$$o$ FROM t WHERE a = x'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT B'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT B'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| SELECT nc$p$o 1 $$/*;;$$ AS x;
| SELECT 'a;bselect 1;\' AS x (SELECT 1; SELECT 2) 1 "end;;" AS x; -- trailing note;
| COPY quf FROM stdin;
| a;b
| 1	\N
| \.
|
| col e'\n;'';' , uo$g$b col;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| BEGIN
< ok 0 -
> exec
| SELECT nc$p$o 1 $$/*;;$$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT nc$p$o 1 $$/*;;$$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| AS x CREATE FUNCTION kwlq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;end;'; END AS x "--;" + U&';$x$;' 1; -- trailing note;
| , '$;$select 1;);(' SELECT;
| COPY wfm FROM stdin;
| a;b
| a;b
| \.
| FROM t WHERE a = "/*;end;end;" 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION kwlq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;end;'; END AS x "--;" + U&';$x$;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION kwlq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;end;'; END AS x "--;" + U&';$x$;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) , e'\';\n;\';' AS x;
| SELECT e'\\\\\n;' col 'a;b);($x$;' ,;
| , CREATE FUNCTION _qoo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;; '; SELECT 1; END FROM t WHERE a = X'01;10' AS x; COPY fcs FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) , e'\';\n;\';' AS x;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) , e'\';\n;\';' AS x;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| COPY owu FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY owu FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY owu FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY nvr FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| COPY wyb FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY nvr FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY nvr FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY pjs FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pjs FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pjs FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY dgz FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dgz FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dgz FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = CREATE FUNCTION hwvu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b;'; SELECT 1; END , (SELECT 1; SELECT 2) +;
| 1 (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION hwvu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b;'; SELECT 1; END , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION hwvu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b;'; SELECT 1; END , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + CREATE FUNCTION wwjs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; $x$;'; END col;
| COPY fmc FROM stdin;
| a;b
| 1	\N
| \.
|
| col U&'/*;);();(' AS x "; " col /*end;; */ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION wwjs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; $x$;'; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION wwjs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; $x$;'; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + e'\n;\n;\n;' 1 'a;b);(' AS x B'01;10' 1; -- trailing note;
| 1 (SELECT 1; SELECT 2) 1 e'\';\n;' +; FROM t WHERE a = /*/*);(;;*/*/ SELECT U&'a;b' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + e'\n;\n;\n;' 1 'a;b);(' AS x B'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + e'\n;\n;\n;' 1 'a;b);(' AS x B'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY cwj FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cwj FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cwj FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x E'\';' col ');(; \' col;
| + "END ;" , "END ;; " 1 CREATE FUNCTION zbuk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$; ; '; END ,; -- trailing note;
| COPY bpz FROM stdin;
| 1	\N
| a;b
| \.
|
| COPY ozb FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x E'\';' col ');(; \' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x E'\';' col ');(; \' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY iil FROM stdin;
| 1	\N
| 1	\N
| \.
|
| COPY mpd FROM stdin;
| a;b
| a;b
| \.
|
| SELECT b'01;10' 1;
| SELECT uo$s$g 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY iil FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY iil FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col $_l$$;$END ;$_l$ 1;
| COPY soe FROM stdin;
| c;d;e
| a;b
| \.
|
| col ';end;;\' 1 $$; $$ , 'select 1;end;' ,; 1 oo$v$l FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $_l$$;$END ;$_l$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $_l$$;$END ;$_l$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , $jeg$select 1;end;select 1;$jeg$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $jeg$select 1;end;select 1;$jeg$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $jeg$select 1;end;select 1;$jeg$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 U&'$x$;);(' AS x E'\\' , U&'select 1;/*;;' AS x;
| COPY ux_ FROM stdin;
| a;b
| a;b
| \.
|  + $$$;$$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'$x$;);(' AS x E'\\' , U&'select 1;/*;;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'$x$;);(' AS x E'\\' , U&'select 1;/*;;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 e'\\\n;' SELECT /*/*/*$x$;;a;b*/*/*/ FROM t WHERE a =;
| col '; --;end;' AS x /*/*/*$;$select 1;select 1;*/*/*/ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\\\n;' SELECT /*/*/*$x$;;a;b*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\\\n;' SELECT /*/*/*$x$;;a;b*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col CREATE FUNCTION ptwt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x B'01;10' col;
| col U&'$;$--;select 1;' + /*/*;;*/*/ FROM t WHERE a = de$b$u 1;
| COPY sqk FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION ptwt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x B'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION ptwt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x B'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) col U&'END ;' AS x;
| + "END ;end;$;$" +;
| COPY hiv FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col U&'END ;' AS x;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col U&'END ;' AS x;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| COPY ikb FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ikb FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ikb FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT $hgsv$--;END ;$hgsv$ FROM t WHERE a = CREATE FUNCTION lapu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');();($x$;'; END col;
| 1 (SELECT 1; SELECT 2) SELECT;
| + b'01;10' , X'01;10' 1 /*a;bend;*/ SELECT; -- trailing note;
| FROM t WHERE a = U&'/*;--;/*;' SELECT U&'; end;' FROM t WHERE a = CREATE FUNCTION xi_j() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;; '; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "CREATE"
| position: 46
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $hgsv$--;END ;$hgsv$ FROM t WHERE a = CREATE FUNCTION lapu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');();($x$;'; END col;
< error
| 42601 syntax error at or near "CREATE"
| position: 46
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $hgsv$--;END ;$hgsv$ FROM t WHERE a = CREATE FUNCTION lapu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');();($x$;'; END col;
< error
| 42601 syntax error at or near "CREATE"
| position: 46
> exec
| COPY zee FROM stdin;
| 1	\N
| c;d;e
| \.
| , U&');(a;b' + (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zee FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zee FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 /*$x$;--;*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*$x$;--;*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*$x$;--;*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = B'01;10' SELECT X'01;10' SELECT U&'; $x$;' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = B'01;10' SELECT X'01;10' SELECT U&'; $x$;' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = B'01;10' SELECT X'01;10' SELECT U&'; $x$;' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x e''';' AS x "a;bselect 1;" 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x e''';' AS x "a;bselect 1;" 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x e''';' AS x "a;bselect 1;" 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + uk$s$f +; SELECT (SELECT 1; SELECT 2) SELECT;
| + /*/*/*--;a;b*/*/*/ + --/*;
|  FROM t WHERE a =;
| col e'\n;\';'';' + '/*;\' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + uk$s$f +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + uk$s$f +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x ec$e$k SELECT (SELECT 1; SELECT 2) FROM t WHERE a = /*/*/*$x$;$x$;*/*/*/ FROM t WHERE a =;
| FROM t WHERE a = U&'--;;;end;' , U&'end;' col 'select 1;\' +;
| FROM t WHERE a = $$);($$ + --END ;; );(
|  AS x CREATE FUNCTION ttgh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;;;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ec$e$k SELECT (SELECT 1; SELECT 2) FROM t WHERE a = /*/*/*$x$;$x$;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ec$e$k SELECT (SELECT 1; SELECT 2) FROM t WHERE a = /*/*/*$x$;$x$;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 /*/*/*;$x$;*/*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*;$x$;*/*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*;$x$;*/*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY oeq FROM stdin;
| a;b
| c;d;e
| \.
|
| 1 /*/*/*$x$;END ;*/*/*/ , (SELECT 1; SELECT 2) col X'01;10' SELECT;1 e'\n;\\\\' +;1 U&'select 1;a;b' SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY oeq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY oeq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY _xd FROM stdin;
| 1	\N
| a;b
| \.
|
| COPY rjy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _xd FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _xd FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 CREATE FUNCTION wqlm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END col X'01;10' SELECT --$;$/*;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION wqlm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END col X'01;10' SELECT --$;$/*;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION wqlm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END col X'01;10' SELECT --$;$/*;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 E''';\n;'';' 1 U&'$;$/*;END ;' 1;
| AS x '/*;' col X'01;10' FROM t WHERE a = (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 E''';\n;'';' 1 U&'$;$/*;END ;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 E''';\n;'';' 1 U&'$;$/*;END ;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , --/*;
|  + ';;' AS x "end;a;b--;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --/*;
|  + ';;' AS x "end;a;b--;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --/*;
|  + ';;' AS x "end;a;b--;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY xlb FROM stdin;
| 1	\N
| 1	\N
| \.
|
| COPY fxv FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xlb FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xlb FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , $$END ;$$ AS x; AS x co$s$s 1 'select 1;' SELECT;
| FROM t WHERE a = U&');(' , 'end;a;b' col CREATE FUNCTION pawr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;);('; SELECT 1; END 1; -- trailing note;
| SELECT /*END ;/ *;/ *;*/ , 'a;b' FROM t WHERE a = "$x$;a;b" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $$END ;$$ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $$END ;$$ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + E'\n;\\\n;' col;
| COPY jwv FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\n;\\\n;' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\n;\\\n;' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY ikg FROM stdin;
| c;d;e
| a;b
| \.
|
| AS x --$;$--;;;
|  FROM t WHERE a = U&'--;select 1;);(' 1 ui$f$n col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ikg FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ikg FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY qql FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qql FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qql FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x /*/*END ;$x$;*/*/ 1 x'01;10' +;
| col (SELECT 1; SELECT 2) col;1 ";" FROM t WHERE a =; SELECT 'END ;select 1;a;b' + (SELECT 1; SELECT 2) SELECT "$x$;;);(" ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*END ;$x$;*/*/ 1 x'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*END ;$x$;*/*/ 1 x'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY sr_ FROM stdin;
| 1	\N
| 1	\N
| \.
|  + (SELECT 1; SELECT 2) AS x;
| FROM t WHERE a = X'01;10' FROM t WHERE a = /*/*END ;;;*/*/ col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sr_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sr_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT /*/*/*/ *;/ *;*/*/*/ 1;
| + e''';' AS x; , b'01;10' AS x --end;
|  SELECT ');(a;b\' col;
< error
| 42601 syntax error at or near "+"
| position: 32
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*/*/ *;/ *;*/*/*/ 1;
< ok 1 -
> exec
|
| + e''';' AS x;
< error
| 42601 syntax error at or near "+"
| position: 2
> exec
| ROLLBACK
< ok 0 -
> exec
|
| + e''';' AS x;
< error
| 42601 syntax error at or near "+"
| position: 2
> exec
| , e''';\';' FROM t WHERE a = "select 1;--;" +;
| col (SELECT 1; SELECT 2) FROM t WHERE a = "/*;select 1;" SELECT '$;$select 1;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e''';\';' FROM t WHERE a = "select 1;--;" +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e''';\';' FROM t WHERE a = "select 1;--;" +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 E'\\' AS x (SELECT 1; SELECT 2) + E'\\\';' 1; FROM t WHERE a = CREATE FUNCTION mlny() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;end;;'; END FROM t WHERE a =;, e''';'';' , hw$n$f +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 E'\\' AS x (SELECT 1; SELECT 2) + E'\\\';' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 E'\\' AS x (SELECT 1; SELECT 2) + E'\\\';' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY pcj FROM stdin;
| a;b
| a;b
| \.
|
| , --END ;
|  +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pcj FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pcj FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY hbl FROM stdin;
| c;d;e
| a;b
| \.
|  1 ');(--;; \' SELECT; -- trailing note;
| + CREATE FUNCTION osoi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;/*;'; SELECT 1; END AS x $$$x$;;$$ + --/*;);(
|  col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hbl FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hbl FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col "select 1;" + X'01;10' SELECT --;a;bselect 1;
|  col;
| SELECT 'END ;end;' 1;+ e''';\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "select 1;" + X'01;10' SELECT --;a;bselect 1;
|  col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "select 1;" + X'01;10' SELECT --;a;bselect 1;
|  col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 B'01;10' col $$; $$ SELECT CREATE FUNCTION yosz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;END ;'; END 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 B'01;10' col $$; $$ SELECT CREATE FUNCTION yosz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;END ;'; END 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 B'01;10' col $$; $$ SELECT CREATE FUNCTION yosz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;END ;'; END 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 e'\\\n;' col U&'--;' col U&'end;' AS x;
| COPY acm FROM stdin;
| c;d;e
| a;b
| \.
|
| col ';;' col U&'select 1;' SELECT e'\\'';'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\\\n;' col U&'--;' col U&'end;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\\\n;' col U&'--;' col U&'end;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 mi$y$f ,;
| AS x (SELECT 1; SELECT 2) SELECT;col /*/*/*; */*/*/ FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 mi$y$f ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 mi$y$f ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col (SELECT 1; SELECT 2) + $$a;b$$ ,; , U&';select 1;' SELECT nt$s$c , E'\n;\n;' AS x;
| AS x E'\n;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) + $$a;b$$ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) + $$a;b$$ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col CREATE FUNCTION ehlq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(--;'; END 1 "$;$" 1;
| , e'\n;' + CREATE FUNCTION jxhh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; END FROM t WHERE a =; -- trailing note;
| , 'a;b;' 1 --;;select 1;end;
|  + U&';;$x$;END ;' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION ehlq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(--;'; END 1 "$;$" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION ehlq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(--;'; END 1 "$;$" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT U&'; ' 1 X'01;10' 1 E''';'';' ,;
| SELECT b'01;10' + U&'; ' ,;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'; ' 1 X'01;10' 1 E''';'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'; ' 1 X'01;10' 1 E''';'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| + b'01;10' SELECT (SELECT 1; SELECT 2) +; -- trailing note;
| SELECT $q$;;$;$;;$q$ 1 $p$$;$END ;end;$x$$p$ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + b'01;10' SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + b'01;10' SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = --select 1;
|  SELECT;
| 1 b'01;10' SELECT;
| col /*/*);(*/*/ , CREATE FUNCTION hqud() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; END SELECT $_xv$end;;;--;$_xv$ AS x;
| , B'01;10' FROM t WHERE a = E'\\\\\';' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --select 1;
|  SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --select 1;
|  SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col co$l$u FROM t WHERE a = '/*;$x$;--;' AS x $$$;$end;$$ SELECT;
| COPY zwz FROM stdin;
| a;b
| c;d;e
| \.
|  AS x ');(' FROM t WHERE a = "$;$" col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col co$l$u FROM t WHERE a = '/*;$x$;--;' AS x $$$;$end;$$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col co$l$u FROM t WHERE a = '/*;$x$;--;' AS x $$$;$end;$$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY ftt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ftt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ftt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT e'\';\n;\\' + $$a;b$$ SELECT;
| , $$a;b;;$$ FROM t WHERE a =; -- trailing note;
| COPY yky FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 38
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\';\n;\\' + $$a;b$$ SELECT;
< error
| 42725 operator is not unique: unknown + unknown
| position: 20
| hint: Could not choose a best candidate operator. You might need to add explicit type casts.
> exec
| ROLLBACK
< ok 0 -
> exec
|
| , $$a;b;;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 2
> exec
| FROM t WHERE a = U&'$;$$x$;--;' ,; -- trailing note;
| , CREATE FUNCTION jbwa() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$;$'; SELECT 1; END +;AS x E'\';'';' AS x CREATE FUNCTION kozi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;END ;'; SELECT 1; END col;1 E'\\\';\';' AS x U&'END ;END ;' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'$;$$x$;--;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'$;$$x$;--;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col CREATE FUNCTION jdxt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;'; SELECT 1; END AS x "end;" SELECT CREATE FUNCTION bfwu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$$;$'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION jdxt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;'; SELECT 1; END AS x "end;" SELECT CREATE FUNCTION bfwu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$$;$'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION jdxt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;'; SELECT 1; END AS x "end;" SELECT CREATE FUNCTION bfwu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$$;$'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col ");(end;end;" col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ");(end;end;" col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ");(end;end;" col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = U&'end;;' AS x;1 E'\';'';'';' SELECT;
| COPY vyj FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'end;;' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'end;;' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT 'END ;' AS x U&'END ;;' SELECT;AS x mp$t$a , B'01;10' ,;
< error
| 42601 syntax error at or near "U&'END ;;'"
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT 'END ;' AS x U&'END ;;' SELECT;
< error
| 42601 syntax error at or near "U&'END ;;'"
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT 'END ;' AS x U&'END ;;' SELECT;
< error
| 42601 syntax error at or near "U&'END ;;'"
| position: 21
> exec
| , $r$/*;END ;; $r$ AS x 'select 1;' ,;
| COPY tfh FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $r$/*;END ;; $r$ AS x 'select 1;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $r$/*;END ;; $r$ AS x 'select 1;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + 'end;);(END ;\' col x'01;10' + $$;select 1;$$ +; COPY zpb FROM stdin;
| a;b
| c;d;e
| \.
|
| 1 b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + 'end;);(END ;\' col x'01;10' + $$;select 1;$$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + 'end;);(END ;\' col x'01;10' + $$;select 1;$$ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = CREATE FUNCTION zlil() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END SELECT '$;$/*;\' ,;
| COPY orl FROM stdin;
| c;d;e
| a;b
| \.
|
| col "--;" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION zlil() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END SELECT '$;$/*;\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION zlil() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END SELECT '$;$/*;\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x U&'; $;$' ,;
| COPY fak FROM stdin;
| c;d;e
| a;b
| \.
|
| 1 U&'$x$;end;' FROM t WHERE a = "--;$x$;" SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'; $;$' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'; $;$' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY bhh FROM stdin;
| a;b
| a;b
| \.
|
| col br$z$h AS x hr$l$_ FROM t WHERE a = "end;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bhh FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bhh FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + ');(' AS x $$--;a;bEND ;$$ 1;
| SELECT /*/*select 1;select 1;;;*/*/ , /*;*/ + wk$z$r +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ');(' AS x $$--;a;bEND ;$$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ');(' AS x $$--;a;bEND ;$$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY oqu FROM stdin;
| c;d;e
| 1	\N
| \.
|
| FROM t WHERE a = $h$;;$h$ SELECT;
| COPY fbu FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY oqu FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY oqu FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = gu$s$e SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = gu$s$e SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = gu$s$e SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + E'\\\n;\';' FROM t WHERE a = CREATE FUNCTION mqus() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END + "/*;; /*;" +; SELECT /*end;a;ba;b*/ SELECT CREATE FUNCTION nbok() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\\\n;\';' FROM t WHERE a = CREATE FUNCTION mqus() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END + "/*;; /*;" +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\\\n;\';' FROM t WHERE a = CREATE FUNCTION mqus() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END + "/*;; /*;" +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a = E'\';\\\\' FROM t WHERE a =;
| COPY eu_ FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY hmh FROM stdin;
| 1	\N
| a;b
| \.
|
| COPY irc FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a = E'\';\\\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a = E'\';\\\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , CREATE FUNCTION mwds() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;/*;$x$;'; SELECT 1; END AS x; -- trailing note;
| + ----;;$;$
|  AS x X'01;10' FROM t WHERE a = ";" AS x;
| 1 '--;$x$;$x$;\' 1 (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION mwds() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;/*;$x$;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION mwds() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;/*;$x$;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 E'\\\n;\n;' SELECT dm$d$t 1;
| SELECT '; ;END ;\' AS x;
| , --/*;;;;
|  AS x x'01;10' + U&';;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 E'\\\n;\n;' SELECT dm$d$t 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 E'\\\n;\n;' SELECT dm$d$t 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 --$;$select 1;
|  col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --$;$select 1;
|  col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --$;$select 1;
|  col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY cqp FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY pld FROM stdin;
| a;b
| a;b
| \.
|
| SELECT cr$i$m , "$;$" , "end;" ,; -- trailing note;
| , xp$w$z FROM t WHERE a = --end;; );(
|  ,;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cqp FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cqp FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY bij FROM stdin;
| c;d;e
| a;b
| \.
|
| col $pj$--;$x$;$pj$ FROM t WHERE a = U&');(;;' + (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bij FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bij FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 --select 1;--;
|  AS x; -- trailing note;
| col $ono$select 1;;select 1;$onx$$ono$ AS x ';a;b$x$;\' + (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --select 1;--;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --select 1;--;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY qcn FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| FROM t WHERE a = (SELECT 1; SELECT 2) SELECT /*--;END ;*/ + E''';\\\n;' +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qcn FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qcn FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col --; 
|  FROM t WHERE a = CREATE FUNCTION lmwh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;$x$;'; SELECT 1; END + (SELECT 1; SELECT 2) SELECT; AS x CREATE FUNCTION u_db() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';select 1;'; END SELECT X'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --; 
|  FROM t WHERE a = CREATE FUNCTION lmwh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;$x$;'; SELECT 1; END + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --; 
|  FROM t WHERE a = CREATE FUNCTION lmwh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;$x$;'; SELECT 1; END + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY _nw FROM stdin;
| c;d;e
| a;b
| \.
|
| + /*/*/*$x$;*/*/*/ AS x $gn$;;END ;$gn$ FROM t WHERE a =; COPY h__ FROM stdin;
| a;b
| a;b
| \.
|
| AS x CREATE FUNCTION sgsd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$x$;'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _nw FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _nw FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 "a;b$x$;--;" +; -- trailing note;
| , (SELECT 1; SELECT 2) FROM t WHERE a = x'01;10' SELECT;+ e'\\' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "a;b$x$;--;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "a;b$x$;--;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , CREATE FUNCTION oedt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END , $xt$--;end;$xt$ SELECT xz$c$m 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION oedt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END , $xt$--;end;$xt$ SELECT xz$c$m 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION oedt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END , $xt$--;end;$xt$ SELECT xz$c$m 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + $vtmm$);(end;;$vtmm$ FROM t WHERE a = --;end;
|  col;
| + --);(
|  col;
| + --$x$;);(
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $vtmm$);(end;;$vtmm$ FROM t WHERE a = --;end;
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $vtmm$);(end;;$vtmm$ FROM t WHERE a = --;end;
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY jjk FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY trb FROM stdin;
| c;d;e
| c;d;e
| \.
|  AS x 'end;' AS x $$--;a;bEND ;$$ ,; , --END ;select 1;;
|  + --END ;;;
|  FROM t WHERE a = (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jjk FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jjk FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x $$a;b$$ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $$a;b$$ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $$a;b$$ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY yjz FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| FROM t WHERE a = B'01;10' +;
| + ');($;$' SELECT pn$f$c +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yjz FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yjz FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 b'01;10' SELECT "$;$/*;;" , CREATE FUNCTION uooh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 b'01;10' SELECT "$;$/*;;" , CREATE FUNCTION uooh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 b'01;10' SELECT "$;$/*;;" , CREATE FUNCTION uooh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x E'\\\';' SELECT;COPY _m_ FROM stdin;
| 1	\N
| c;d;e
| \.
|
| COPY djw FROM stdin;
| a;b
| 1	\N
| \.
| FROM t WHERE a = CREATE FUNCTION jtbi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END , e'\\' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x E'\\\';' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x E'\\\';' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 x'01;10' 1 ');(;' col; col 'a;b' SELECT; -- trailing note;
| FROM t WHERE a = E'\\'';' +;
| col $n$select 1;END ;select 1;$n$ + e'\';\n;\';' SELECT /*end;;*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 x'01;10' 1 ');(;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 x'01;10' 1 ');(;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT ls$m$b SELECT;
| FROM t WHERE a = "$;$" col; -- trailing note;
| SELECT cr$n$y FROM t WHERE a = --$;$;;
|  +;
< error
| 42601 syntax error at or near "FROM"
| position: 23
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ls$m$b SELECT;
< error
| 42703 column "ls$m$b" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| FROM t WHERE a = "$;$" col;
< error
| 42601 syntax error at or near "FROM"
| position: 2
> exec
| FROM t WHERE a = e''';\';' , $ivvj$$x$;/*;$ivvj$ + CREATE FUNCTION t_xq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END ,;SELECT e'\\' , ";END ;$;$" +; , $iioe$--;$iioe$ FROM t WHERE a = B'01;10' SELECT --$x$;/*;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = e''';\';' , $ivvj$$x$;/*;$ivvj$ + CREATE FUNCTION t_xq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = e''';\';' , $ivvj$$x$;/*;$ivvj$ + CREATE FUNCTION t_xq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col '/*;\' col;
| COPY kgp FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '/*;\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '/*;\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x '$;$a;b' , B'01;10' 1 $qzw$select 1;END ;$qzw$ +; -- trailing note;
| col --END ;end;
|  1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x '$;$a;b' , B'01;10' 1 $qzw$select 1;END ;$qzw$ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x '$;$a;b' , B'01;10' 1 $qzw$select 1;END ;$qzw$ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 X'01;10' ,;
| AS x '; END ;' col;
| col (SELECT 1; SELECT 2) , $$$;@;$$$ FROM t WHERE a =;, --; a;b
|  col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 X'01;10' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 X'01;10' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , b'01;10' AS x --; );(
|  SELECT;
| COPY dtb FROM stdin;
| c;d;e
| a;b
| \.
|
| SELECT ');($;$\' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , b'01;10' AS x --; );(
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , b'01;10' AS x --; );(
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = 'select 1;select 1;' + CREATE FUNCTION bjag() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';);(--;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = 'select 1;select 1;' + CREATE FUNCTION bjag() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';);(--;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = 'select 1;select 1;' + CREATE FUNCTION bjag() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';);(--;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x /*/*select 1;*/*/ AS x ah$q$q , CREATE FUNCTION njxs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END col;
| COPY tye FROM stdin;
| a;b
| 1	\N
| \.
|  COPY ilr FROM stdin;
| a;b
| 1	\N
| \.
|
| 1 '/*;;;\' + uu$w$k + ';' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*select 1;*/*/ AS x ah$q$q , CREATE FUNCTION njxs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*select 1;*/*/ AS x ah$q$q , CREATE FUNCTION njxs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x ';' SELECT;
| FROM t WHERE a = E'\n;\';\n;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ';' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ';' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + $uj$);(--;$uj$ SELECT (SELECT 1; SELECT 2) col ';);($;$' FROM t WHERE a =;
| , (SELECT 1; SELECT 2) FROM t WHERE a = x'01;10' , ku$k$v col;, $nmk$;;/*;end;$nmk$ 1 (SELECT 1; SELECT 2) AS x;
| + 'end;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $uj$);(--;$uj$ SELECT (SELECT 1; SELECT 2) col ';);($;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $uj$);(--;$uj$ SELECT (SELECT 1; SELECT 2) col ';);($;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT '; end;\' +;
< error
| 42601 syntax error at or near ";"
| position: 19
> exec
| BEGIN
< ok 0 -
> exec
| SELECT '; end;\' +;
< error
| 42601 syntax error at or near ";"
| position: 19
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT '; end;\' +;
< error
| 42601 syntax error at or near ";"
| position: 19
> exec
| COPY ntj FROM stdin;
| a;b
| 1	\N
| \.
| SELECT "select 1;$x$;" 1;
| + $$end;$$ + (SELECT 1; SELECT 2) AS x U&'$;$END ;' +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ntj FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ntj FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + e'\';\\\';' + e'\';' SELECT e''';'';' col;COPY gjz FROM stdin;
| c;d;e
| 1	\N
| \.
|
| SELECT e''';'';' FROM t WHERE a = U&'--;;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + e'\';\\\';' + e'\';' SELECT e''';'';' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + e'\';\\\';' + e'\';' SELECT e''';'';' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY zin FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zin FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zin FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = "/*;END ;end;" AS x X'01;10' col /*/*/*;;select 1;*/*/*/ ,; -- trailing note;
| SELECT "select 1;; $;$" 1 $$;);(a;b$$ SELECT /*/*/*$;$end;*/*/*/ 1;
| SELECT CREATE FUNCTION xuac() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; END +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "/*;END ;end;" AS x X'01;10' col /*/*/*;;select 1;*/*/*/ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "/*;END ;end;" AS x X'01;10' col /*/*/*;;select 1;*/*/*/ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + /*/*--;$;$*/*/ FROM t WHERE a = E'\';'';' FROM t WHERE a =; COPY aad FROM stdin;
| a;b
| c;d;e
| \.
|
| , U&'END ;' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*--;$;$*/*/ FROM t WHERE a = E'\';'';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*--;$;$*/*/ FROM t WHERE a = E'\';'';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT /*/*/*; --;;*/*/*/ 1 /*/*end;;;*/*/ +;
| COPY faw FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 45
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*/*; --;;*/*/*/ 1 /*/*end;;;*/*/ +;
< error
| 42601 syntax error at or near ";"
| position: 45
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*/*/*; --;;*/*/*/ 1 /*/*end;;;*/*/ +;
< error
| 42601 syntax error at or near ";"
| position: 45
> exec
| 1 --;);(
|  AS x vm$a$z ,; -- trailing note;
| AS x "; " SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --;);(
|  AS x vm$a$z ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --;);(
|  AS x vm$a$z ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 $$END ;$$ 1 uo$m$t 1 U&'END ;$x$;' AS x;
| SELECT --$x$;/*;--;
|  AS x CREATE FUNCTION oclz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $$END ;$$ 1 uo$m$t 1 U&'END ;$x$;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $$END ;$$ 1 uo$m$t 1 U&'END ;$x$;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , /*; $x$;*/ col;
| FROM t WHERE a = "--;select 1;" 1 $dqfd$--;; $dqfx$$dqfd$ , --;;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*; $x$;*/ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*; $x$;*/ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col ----;
|  AS x;
| col 'a;b; ' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ----;
|  AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ----;
|  AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) AS x; -- trailing note;
| AS x /*;;/ *;;;*/ col '$;$' SELECT;
| , ";" FROM t WHERE a = E'\';\';\\' + X'01;10' ,;
| FROM t WHERE a = /*/*/*;a;b*/*/*/ col --select 1;select 1;
|  col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col s_$d$k AS x '--;' FROM t WHERE a =; AS x CREATE FUNCTION sgbz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;a;bend;'; SELECT 1; END SELECT;, (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col s_$d$k AS x '--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col s_$d$k AS x '--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , ');(' FROM t WHERE a = ';$;$end;\' SELECT $ft_q$END ;$ft_q$ 1;
| FROM t WHERE a = "a;bend;; " + X'01;10' col; -- trailing note;
| 1 ");(" FROM t WHERE a = e''';' + U&'$;$/*;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ');(' FROM t WHERE a = ';$;$end;\' SELECT $ft_q$END ;$ft_q$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ');(' FROM t WHERE a = ';$;$end;\' SELECT $ft_q$END ;$ft_q$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + x'01;10' col;COPY rrz FROM stdin;
| a;b
| 1	\N
| \.
| COPY nuc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + x'01;10' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + x'01;10' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 ");(" 1 $$/*;--;END ;$$ , U&'--;' col;1 B'01;10' col;AS x CREATE FUNCTION wbai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;);('; SELECT 1; END , $$a;b$;$$$ ,;
| COPY trp FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ");(" 1 $$/*;--;END ;$$ , U&'--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ");(" 1 $$/*;--;END ;$$ , U&'--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY and FROM stdin;
| a;b
| 1	\N
| \.
| FROM t WHERE a = e''';' +;
| 1 --end;
|  SELECT; -- trailing note;
| col ----;; 
|  1;
< error
| 42601 syntax error at or near "and"
| position: 6
> exec
| BEGIN
< ok 0 -
> exec
| COPY and FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "and"
| position: 6
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY and FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "and"
| position: 6
> exec
| COPY m_c FROM stdin;
| c;d;e
| 1	\N
| \.
|
| , U&'a;ba;b);(' , U&'select 1;' FROM t WHERE a = '; a;b' +;, /*/*a;b*/*/ SELECT ck$v$o SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY m_c FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY m_c FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , sh$t$d col $lkw$/*;$lkx$$lkw$ col U&'end;;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , sh$t$d col $lkw$/*;$lkx$$lkw$ col U&'end;;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , sh$t$d col $lkw$/*;$lkx$$lkw$ col U&'end;;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT ';;;\' AS x;
| AS x x'01;10' 1 '/*;END ;' AS x; COPY fat FROM stdin;
| 1	\N
| a;b
| \.
|
| + "/*;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ';;;\' AS x;
< ok 1 -
> exec
|
| AS x x'01;10' 1 '/*;END ;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 2
> exec
| ROLLBACK
< ok 0 -
> exec
|
| AS x x'01;10' 1 '/*;END ;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 2
> exec
| FROM t WHERE a = --select 1;$;$$x$;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --select 1;$;$$x$;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --select 1;$;$$x$;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col /*/*$x$;*/*/ SELECT 'a;b\' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*$x$;*/*/ SELECT 'a;b\' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*$x$;*/*/ SELECT 'a;b\' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY teq FROM stdin;
| 1	\N
| a;b
| \.
|
| COPY hpi FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY teq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY teq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col /*/*$;$$x$;; */*/ + hu$c$o 1; -- trailing note;
| SELECT e''';' 1 E'\n;\';' , e''';'';' 1;
| FROM t WHERE a = $$;a;b;;$$ ,;
| 1 '; ' , U&'select 1;);(/*;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*$;$$x$;; */*/ + hu$c$o 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*$;$$x$;; */*/ + hu$c$o 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY ehq FROM stdin;
| 1	\N
| a;b
| \.
|
| + --end;end;
|  1;
| COPY wwd FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ehq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ehq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , x'01;10' AS x --$x$;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , x'01;10' AS x --$x$;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , x'01;10' AS x --$x$;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY zbr FROM stdin;
| a;b
| c;d;e
| \.
| col e'\n;\n;\\' AS x kq$c$n SELECT '$;$END ;/*;\' col; , "select 1;end;$x$;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zbr FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zbr FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col U&'$x$;' col;
| COPY oo_ FROM stdin;
| 1	\N
| 1	\N
| \.
|  COPY dmc FROM stdin;
| a;b
| a;b
| \.
| SELECT $f$$;$);($f$ + $$);();($$ , --;;
|  +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'$x$;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'$x$;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT --);(
|  SELECT $pj$$;$$x$;$px$$pj$ , (SELECT 1; SELECT 2) +;
| , (SELECT 1; SELECT 2) + e''';\';\\' + ');(' ,;
| col $$;$$ SELECT;AS x 'select 1;\' FROM t WHERE a = /*/*end;$;$$x$;*/*/ 1;
< error
| 42601 syntax error at or near "SELECT"
| position: 15
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --);(
|  SELECT $pj$$;$$x$;$px$$pj$ , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "SELECT"
| position: 15
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --);(
|  SELECT $pj$$;$$x$;$px$$pj$ , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "SELECT"
| position: 15
> exec
| 1 CREATE FUNCTION qhbw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;;;'; SELECT 1; END FROM t WHERE a = $$$;$);(;$$ +;
| 1 $$end;$$ AS x B'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION qhbw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;;;'; SELECT 1; END FROM t WHERE a = $$$;$);(;$$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION qhbw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;;;'; SELECT 1; END FROM t WHERE a = $$$;$);(;$$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY ayy FROM stdin;
| 1	\N
| c;d;e
| \.
|
| FROM t WHERE a = 'a;b\' AS x CREATE FUNCTION _yod() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;end;'; END SELECT "; " +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ayy FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ayy FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + B'01;10' SELECT /*/*/*;*/*/*/ 1 "$x$;" AS x;FROM t WHERE a = e'\\'';\n;' SELECT "$x$;/*;" FROM t WHERE a = --;
|  AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + B'01;10' SELECT /*/*/*;*/*/*/ 1 "$x$;" AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + B'01;10' SELECT /*/*/*;*/*/*/ 1 "$x$;" AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 qd$u$m +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 qd$u$m +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 qd$u$m +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , B'01;10' 1;
| col "--;end;/*;" ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , B'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , B'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , CREATE FUNCTION nic_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;; '; SELECT 1; END AS x "$;$$x$;" +;SELECT (SELECT 1; SELECT 2) AS x; 1 '$x$;--;;;\' SELECT CREATE FUNCTION mgea() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$;$'; END ,;1 (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION nic_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;; '; SELECT 1; END AS x "$;$$x$;" +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION nic_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;; '; SELECT 1; END AS x "$;$$x$;" +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = /*/*/*select 1;;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*/*/*select 1;;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*/*/*select 1;;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 CREATE FUNCTION dmma() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bEND ;'; SELECT 1; END col U&';;' FROM t WHERE a =;
| col "; );(;" +; 1 'END ;' 1 /*/*/*$x$;*/*/*/ , E''';'';' FROM t WHERE a =;
| COPY _np FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION dmma() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bEND ;'; SELECT 1; END col U&';;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION dmma() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bEND ;'; SELECT 1; END col U&';;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , ");(" FROM t WHERE a = E''';\n;' col;
| FROM t WHERE a = "--;;;" 1 ';;\' 1;
| AS x "/*;select 1;/*;" SELECT b'01;10' SELECT '/*;/*;--;\' SELECT; col $flth$end;$;$$flth$ FROM t WHERE a = bc$e$u SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ");(" FROM t WHERE a = E''';\n;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ");(" FROM t WHERE a = E''';\n;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY rvk FROM stdin;
| 1	\N
| a;b
| \.
|
| , "/*;END ;; " AS x (SELECT 1; SELECT 2) +;
| col (SELECT 1; SELECT 2) SELECT;
| AS x '$;$end;END ;' , --$;$$;$end;
|  SELECT ";;" ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY rvk FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY rvk FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + E'\';'';' SELECT e'\\' col CREATE FUNCTION eeo_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;'; END 1;
| SELECT e'\';\';\';' AS x --;
|  AS x --END ;; --;
|  col;
| AS x CREATE FUNCTION pivo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END + 'a;b' FROM t WHERE a = /*/*/*/ *;END ;*/*/*/ +; -- trailing note;
| 1 e'\';\n;\';' col U&'; ' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\';'';' SELECT e'\\' col CREATE FUNCTION eeo_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;'; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\';'';' SELECT e'\\' col CREATE FUNCTION eeo_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;'; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY cgp FROM stdin;
| 1	\N
| 1	\N
| \.
|
| COPY qnj FROM stdin;
| 1	\N
| c;d;e
| \.
|  , E'\n;\\' FROM t WHERE a = CREATE FUNCTION ckhh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(END ;;'; END FROM t WHERE a = ';' +; -- trailing note;
| COPY ivh FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cgp FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cgp FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY pbv FROM stdin;
| a;b
| a;b
| \.
|
| , --/*;);(;
|  AS x --$;$end;;
|  1 E'\n;' ,;
| , (SELECT 1; SELECT 2) 1 "select 1;select 1;END ;" col CREATE FUNCTION udsx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$x$;'; SELECT 1; END ,;
| + (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pbv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pbv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY kyf FROM stdin;
| c;d;e
| a;b
| \.
|  FROM t WHERE a = '--;select 1;end;' + "$x$;" 1;
| col x'01;10' + /*/*/*; ; */*/*/ 1 /*/*/*select 1;*/*/*/ ,;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kyf FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kyf FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , --;END ;
|  col CREATE FUNCTION cyye() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END , U&'/*;' SELECT;
| COPY gid FROM stdin;
| 1	\N
| a;b
| \.
|
| col (SELECT 1; SELECT 2) SELECT;
| COPY xer FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --;END ;
|  col CREATE FUNCTION cyye() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END , U&'/*;' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --;END ;
|  col CREATE FUNCTION cyye() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END , U&'/*;' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , "; $x$;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "; $x$;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "; $x$;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT b'01;10' col b'01;10' +;
| SELECT (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "b'01;10'"
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT b'01;10' col b'01;10' +;
< error
| 42601 syntax error at or near "b'01;10'"
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT b'01;10' col b'01;10' +;
< error
| 42601 syntax error at or near "b'01;10'"
| position: 21
> exec
| COPY mkl FROM stdin;
| c;d;e
| 1	\N
| \.
|  SELECT zj$_$o 1 CREATE FUNCTION umne() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; SELECT 1; END + e'\n;\\' col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mkl FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mkl FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col $js$--;$jx$$js$ + --end;END ;
|  col (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $js$--;$jx$$js$ + --end;END ;
|  col (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $js$--;$jx$$js$ + --end;END ;
|  col (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT '$;$$;$END ;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 39
> exec
| BEGIN
< ok 0 -
> exec
| SELECT '$;$$;$END ;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 39
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT '$;$$;$END ;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 39
> exec
| 1 /*a;bend;*/ SELECT U&'select 1;;$;$' 1;
| FROM t WHERE a = --/*;
|  col E'\\\';' 1 /*/*; --;*/*/ AS x;
| COPY cpg FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*a;bend;*/ SELECT U&'select 1;;$;$' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*a;bend;*/ SELECT U&'select 1;;$;$' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + /*/*/*a;b;;select 1;*/*/*/ 1;
| SELECT qe$b$u col;
| col CREATE FUNCTION duec() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;END ;$x$;'; END col '/*;/*;' , b'01;10' ,;
| col "select 1;END ;END ;" SELECT $a$$;$);(;;$x$$a$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*/*a;b;;select 1;*/*/*/ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*/*a;b;;select 1;*/*/*/ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY sdm FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 CREATE FUNCTION gwra() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$);('; END SELECT _o$d$s AS x;
| + (SELECT 1; SELECT 2) FROM t WHERE a = U&');(' , b'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sdm FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sdm FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col U&'end;);(' col CREATE FUNCTION lsul() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END FROM t WHERE a = "/*;;--;" 1; -- trailing note;
| 1 $$select 1;$$ 1 U&'END ;END ;select 1;' col 'select 1;;/*;' 1;
| FROM t WHERE a = '$;$select 1;END ;' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'end;);(' col CREATE FUNCTION lsul() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END FROM t WHERE a = "/*;;--;" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'end;);(' col CREATE FUNCTION lsul() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END FROM t WHERE a = "/*;;--;" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY gul FROM stdin;
| a;b
| 1	\N
| \.
| COPY ssa FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY tok FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gul FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gul FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY pva FROM stdin;
| c;d;e
| c;d;e
| \.
|  1 '--;\' FROM t WHERE a = --;);(end;
|  1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pva FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pva FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , "end;END ;end;" FROM t WHERE a = e'\';'';' col; -- trailing note;
| 1 ');(\' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "end;END ;end;" FROM t WHERE a = e'\';'';' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "end;END ;end;" FROM t WHERE a = e'\';'';' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + --end;; 
|  SELECT; -- trailing note;
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = ");();(" +;
| 1 lo$q$o +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --end;; 
|  SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --end;; 
|  SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 CREATE FUNCTION ubyc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; );('; SELECT 1; END SELECT $$;$;$select 1;$$ + U&'select 1;a;b--;' col;
| col /*/*/*;);(a;b*/*/*/ + --select 1;END ;
|  FROM t WHERE a = 'select 1;);(' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION ubyc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; );('; SELECT 1; END SELECT $$;$;$select 1;$$ + U&'select 1;a;b--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION ubyc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; );('; SELECT 1; END SELECT $$;$;$select 1;$$ + U&'select 1;a;b--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY lsg FROM stdin;
| c;d;e
| a;b
| \.
|  -- trailing note;
| col E'\';'';\\' FROM t WHERE a = "end;;" AS x;
| , --$x$;END ;
|  FROM t WHERE a = $$--;$;$$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lsg FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lsg FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x --END ;
|  + $coyf$;END ;$coyx$$coyf$ FROM t WHERE a =;
| AS x CREATE FUNCTION lzlw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;select 1;a;b'; SELECT 1; END col; -- trailing note;
| col '/*;\' AS x CREATE FUNCTION zmiq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --END ;
|  + $coyf$;END ;$coyx$$coyf$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --END ;
|  + $coyf$;END ;$coyx$$coyf$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT --);(
|  col (SELECT 1; SELECT 2) ,;
| SELECT (SELECT 1; SELECT 2) AS x;
| 1 /*a;b--;*/ , "/*;" +;
< error
| 42601 syntax error at or near "SELECT"
| position: 20
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --);(
|  col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "SELECT"
| position: 20
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --);(
|  col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "SELECT"
| position: 20
> exec
| FROM t WHERE a = E''';\n;'';' AS x $$/*;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E''';\n;'';' AS x $$/*;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E''';\n;'';' AS x $$/*;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = rc$n$h SELECT; -- trailing note;
| COPY wbh FROM stdin;
| a;b
| a;b
| \.
|
| 1 $$$x$;a;b; $$ ,; -- trailing note;
| COPY fpt FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = rc$n$h SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = rc$n$h SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 CREATE FUNCTION wrxo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';END ;--;'; SELECT 1; END col "end;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION wrxo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';END ;--;'; SELECT 1; END col "end;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION wrxo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';END ;--;'; SELECT 1; END col "end;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x ne$_$a col;
| SELECT CREATE FUNCTION oemv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END SELECT ';;' AS x bu$_$j ,;
| COPY hqp FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| COPY hsk FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ne$_$a col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ne$_$a col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = $$;;$$ col as$r$c col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$;;$$ col as$r$c col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$;;$$ col as$r$c col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col "--;" FROM t WHERE a = '--;);(' + /*;;*/ AS x;
| COPY egd FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "--;" FROM t WHERE a = '--;);(' + /*;;*/ AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "--;" FROM t WHERE a = '--;);(' + /*;;*/ AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY vyk FROM stdin;
| c;d;e
| 1	\N
| \.
|
| , b'01;10' SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vyk FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vyk FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY eyg FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY uyh FROM stdin;
| c;d;e
| c;d;e
| \.
|  COPY miy FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY eyg FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY eyg FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , (SELECT 1; SELECT 2) + '$x$;' AS x e''';\\'';' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) + '$x$;' AS x e''';\\'';' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) + '$x$;' AS x e''';\\'';' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY gma FROM stdin;
| 1	\N
| 1	\N
| \.
|
| col 'a;bEND ;' , x'01;10' FROM t WHERE a =; -- trailing note;
| SELECT '--;' AS x U&';$;$a;b' 1 E''';' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gma FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gma FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , 'END ;;;/*;\' , /*/*/*$x$;;;;;*/*/*/ 1;
| + E'\';\n;\n;' +;
| COPY zph FROM stdin;
| a;b
| 1	\N
| \.
|
| SELECT $$a;b$x$;$$ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , 'END ;;;/*;\' , /*/*/*$x$;;;;;*/*/*/ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , 'END ;;;/*;\' , /*/*/*$x$;;;;;*/*/*/ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x CREATE FUNCTION sklz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;; '; SELECT 1; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION sklz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;; '; SELECT 1; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION sklz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;; '; SELECT 1; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT e'\n;\n;' col;
< ok 1 -
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\n;\n;' col;
< ok 1 -
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY waq FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| COPY hrv FROM stdin;
| c;d;e
| a;b
| \.
|
| + B'01;10' SELECT "/*;" AS x;FROM t WHERE a = ";a;b" ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY waq FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY waq FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col e'\';\\'';' 1;
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = "; end;);(" SELECT CREATE FUNCTION waot() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e'\';\\'';' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e'\';\\'';' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + 'select 1;a;b$;$' + $qgue$--;; ;;$qgue$ AS x;
| , et$c$o +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + 'select 1;a;b$;$' + $qgue$--;; ;;$qgue$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + 'select 1;a;b$;$' + $qgue$--;; ;;$qgue$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY lzi FROM stdin;
| a;b
| c;d;e
| \.
| + "$;$--;" + CREATE FUNCTION bwkf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END +;
| COPY mdo FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| COPY cgh FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lzi FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lzi FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY mvb FROM stdin;
| 1	\N
| c;d;e
| \.
|  col E'\\'';' FROM t WHERE a = U&'select 1;' SELECT dd$n$j FROM t WHERE a =;1 E''';\\' AS x ');(END ;' FROM t WHERE a =;
| , --;;
|  FROM t WHERE a = CREATE FUNCTION mngh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;$;$'; END , 'end;select 1;end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mvb FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mvb FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x /*/*/ *;END ;$x$;*/*/ + "$;$);(" + (SELECT 1; SELECT 2) FROM t WHERE a =;
| COPY uqm FROM stdin;
| c;d;e
| c;d;e
| \.
|  col $wzs$a;bEND ;$x$;$wzs$ FROM t WHERE a =;
| 1 CREATE FUNCTION eopn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$'; SELECT 1; END col B'01;10' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/ *;END ;$x$;*/*/ + "$;$);(" + (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/ *;END ;$x$;*/*/ + "$;$);(" + (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + "; ; /*;" SELECT;
| AS x --; 
|  SELECT $qr$a;b$;$$qr$ , av$i$m FROM t WHERE a =;
| COPY k_v FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "; ; /*;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "; ; /*;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , "end;$;$; " + '$x$;\' AS x U&');(;;' +;
| col qv$s$z FROM t WHERE a =;
| SELECT $jd$);($x$;$jx$$jd$ 1 U&'; ' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "end;$;$; " + '$x$;\' AS x U&');(;;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "end;$;$; " + '$x$;\' AS x U&');(;;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY obl FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY obl FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY obl FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col $$a;b$x$;$$ , /*/*/*$x$;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$a;b$x$;$$ , /*/*/*$x$;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$a;b$x$;$$ , /*/*/*$x$;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x $bk$--;--;$bx$$bk$ SELECT CREATE FUNCTION smoz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; END ,;
| FROM t WHERE a = e''';' , --a;b$x$;/*;
|  1 "a;b);(--;" ,; , "end;" col /*/*--;*/*/ FROM t WHERE a = e'\';'';'';' col;
| COPY g_l FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $bk$--;--;$bx$$bk$ SELECT CREATE FUNCTION smoz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $bk$--;--;$bx$$bk$ SELECT CREATE FUNCTION smoz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT ';;' SELECT; -- trailing note;
| COPY hoa FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 60
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ';;' SELECT;
< ok 1 -
> exec
|  -- trailing note;
| COPY hoa FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 41
> exec
| ROLLBACK
< ok 0 -
> exec
|  -- trailing note;
| COPY hoa FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 41
> exec
| col CREATE FUNCTION oz_s() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;$x$;'; END 1; -- trailing note;
| COPY enj FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION oz_s() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;$x$;'; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION oz_s() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;$x$;'; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 e'\\'';' col --select 1;/*;select 1;
|  1 CREATE FUNCTION yvsr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;;;'; END col;AS x "/*;/*;" SELECT ug$l$h , /*/*--;;;*/*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\\'';' col --select 1;/*;select 1;
|  1 CREATE FUNCTION yvsr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;;;'; END col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\\'';' col --select 1;/*;select 1;
|  1 CREATE FUNCTION yvsr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;;;'; END col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 /*/*end;*/*/ FROM t WHERE a = CREATE FUNCTION lvif() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;);('; END +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*end;*/*/ FROM t WHERE a = CREATE FUNCTION lvif() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;);('; END +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*end;*/*/ FROM t WHERE a = CREATE FUNCTION lvif() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;);('; END +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT w_$z$l 1 $rint$/*;$rint$ + E''';\n;' 1;
| + e'\n;' AS x /*/*end;; */*/ ,;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| BEGIN
< ok 0 -
> exec
| SELECT w_$z$l 1 $rint$/*;$rint$ + E''';\n;' 1;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT w_$z$l 1 $rint$/*;$rint$ + E''';\n;' 1;
< error
| 42601 syntax error at or near "1"
| position: 15
> exec
| + '$x$;' FROM t WHERE a = 'select 1;$;$;' FROM t WHERE a =; -- trailing note;
| COPY fzh FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + '$x$;' FROM t WHERE a = 'select 1;$;$;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + '$x$;' FROM t WHERE a = 'select 1;$;$;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = ";" col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ";" col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ";" col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , /*/*/*select 1;/ *;$x$;*/*/*/ AS x e'\';'';' AS x (SELECT 1; SELECT 2) FROM t WHERE a =;AS x $fsi$/*;; $fsi$ col --a;ba;b
|  FROM t WHERE a = "$x$;select 1;" AS x; + '; ' SELECT (SELECT 1; SELECT 2) + /*/*/*;;*/*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*/*select 1;/ *;$x$;*/*/*/ AS x e'\';'';' AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*/*select 1;/ *;$x$;*/*/*/ AS x e'\';'';' AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , e'\n;' SELECT $$; );($$ col (SELECT 1; SELECT 2) FROM t WHERE a =; -- trailing note;
| SELECT b'01;10' col X'01;10' , ";);(" col; -- trailing note;
| + $$$;$--;$$ SELECT e''';\n;\\' , CREATE FUNCTION orvp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;end;'; END FROM t WHERE a =;
| , (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\n;' SELECT $$; );($$ col (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\n;' SELECT $$; );($$ col (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY dmu FROM stdin;
| c;d;e
| a;b
| \.
| , '$;$/*;' SELECT mt$f$r AS x U&'a;b' SELECT;
| 1 U&'END ;$x$;' 1 (SELECT 1; SELECT 2) col /*/*);(/ *;*/*/ col;
| FROM t WHERE a = rf$p$w + E'\';' FROM t WHERE a = CREATE FUNCTION aszs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;a;bEND ;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dmu FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dmu FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x pj$k$o +;
| COPY fpv FROM stdin;
| 1	\N
| a;b
| \.
|  col $$; ;$$ 1 /*$x$;--;);(*/ ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x pj$k$o +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x pj$k$o +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT /*a;b*/ SELECT; COPY kvi FROM stdin;
| c;d;e
| c;d;e
| \.
|
| AS x yw$g$k +;
< error
| 42601 syntax error at or near "SELECT"
| position: 16
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*a;b*/ SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 16
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*a;b*/ SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 16
> exec
| AS x E''';\';\\' ,;
| SELECT U&'END ;select 1;--;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x E''';\';\\' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x E''';\';\\' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , U&'a;b' SELECT /*/*/*;/ *;*/*/*/ , 'a;b' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'a;b' SELECT /*/*/*;/ *;*/*/*/ , 'a;b' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'a;b' SELECT /*/*/*;/ *;*/*/*/ , 'a;b' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , $$end;);(;;$$ FROM t WHERE a = 'END ;$x$;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $$end;);(;;$$ FROM t WHERE a = 'END ;$x$;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $$end;);(;;$$ FROM t WHERE a = 'END ;$x$;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col '; ; \' 1;
| , /*/*);(end;END ;*/*/ col;
| FROM t WHERE a = CREATE FUNCTION anro() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;);('; SELECT 1; END ,; FROM t WHERE a = CREATE FUNCTION mjkq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;$;$);('; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '; ; \' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '; ; \' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT '/*;' 1 $$;;$$ AS x;
| SELECT "END ;" + $$end;a;bEND ;$$ +;
| 1 "$;$);($x$;" AS x (SELECT 1; SELECT 2) , x'01;10' FROM t WHERE a =; -- trailing note;
| , as$r$g SELECT E'\';'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 14
> exec
| BEGIN
< ok 0 -
> exec
| SELECT '/*;' 1 $$;;$$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 14
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT '/*;' 1 $$;;$$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 14
> exec
| + E''';' SELECT b'01;10' AS x ");(" ,; COPY iia FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E''';' SELECT b'01;10' AS x ");(" ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E''';' SELECT b'01;10' AS x ");(" ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = "--;/*;a;b" FROM t WHERE a = X'01;10' , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "--;/*;a;b" FROM t WHERE a = X'01;10' , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "--;/*;a;b" FROM t WHERE a = X'01;10' , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY fpj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fpj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fpj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + $$; $$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $$; $$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $$; $$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = /*/*);(select 1;);(*/*/ +;
| col x'01;10' , fd$v$c AS x;
| col e'\n;' , /*/*/*;;/ *;$x$;*/*/*/ + (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = /*/*);(select 1;);(*/*/ +;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = /*/*);(select 1;);(*/*/ +;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| , ea$b$x SELECT; -- trailing note;
| SELECT $gk$;$;$/*;$gk$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ea$b$x SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ea$b$x SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY dkx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dkx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dkx FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY wro FROM stdin;
| c;d;e
| a;b
| \.
| 1 "--;;" SELECT;
| col X'01;10' AS x --; 
|  FROM t WHERE a =; COPY jor FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wro FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wro FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT x'01;10' col (SELECT 1; SELECT 2) 1; -- trailing note;
| SELECT U&'--;$x$;;;' FROM t WHERE a = "/*;;" AS x wy$t$i FROM t WHERE a =;
| SELECT --a;bend;
|  1;
< error
| 42601 syntax error at or near "("
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "("
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT x'01;10' col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "("
| position: 21
> exec
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
| SELECT sr$g$k 1 CREATE FUNCTION oens() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;'; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY cwr FROM stdin;
| c;d;e
| c;d;e
| \.
|
| AS x E'\\' col "--;END ;" col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cwr FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cwr FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY cjr FROM stdin;
| c;d;e
| 1	\N
| \.
|
| 1 '$x$;$x$;; ' + (SELECT 1; SELECT 2) FROM t WHERE a = CREATE FUNCTION iphl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';select 1;'; END ,;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cjr FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cjr FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , ok$k$h , nq$z$h col CREATE FUNCTION kcrc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ok$k$h , nq$z$h col CREATE FUNCTION kcrc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ok$k$h , nq$z$h col CREATE FUNCTION kcrc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY wal FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY gum FROM stdin;
| a;b
| c;d;e
| \.
|
| FROM t WHERE a = X'01;10' col $$end;a;b/*;$$ 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wal FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wal FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY xve FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| + CREATE FUNCTION dvmv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END 1 b'01;10' +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xve FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xve FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY _yh FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _yh FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _yh FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , /*/*; */*/ AS x /*/*/*a;b*/*/*/ + --;$;$
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*; */*/ AS x /*/*/*a;b*/*/*/ + --;$;$
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*; */*/ AS x /*/*/*a;b*/*/*/ + --;$;$
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY yyq FROM stdin;
| a;b
| 1	\N
| \.
|
| 1 --;;
|  AS x ay$z$r +;1 --end;;;a;b
|  SELECT x'01;10' SELECT $$END ;END ;; $$ +;
| , CREATE FUNCTION qase() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;a;b; '; SELECT 1; END SELECT e'\';' SELECT x'01;10' SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yyq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yyq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = $$end;$;$$$ 1;
| FROM t WHERE a = /*/*/*select 1;select 1;a;b*/*/*/ SELECT B'01;10' col;
| 1 U&'end;/*;a;b' SELECT ri$o$t AS x;
| 1 ";a;b" SELECT B'01;10' 1 $$end;$$ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$end;$;$$$ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$end;$;$$$ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 /*--;*/ + --END ;--;
|  FROM t WHERE a =;
| 1 --/*;$x$;
|  + ";;end;; " ,;COPY wrn FROM stdin;
| c;d;e
| c;d;e
| \.
| 1 e_$t$r FROM t WHERE a = U&'end;;;' + ci$l$h 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*--;*/ + --END ;--;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*--;*/ + --END ;--;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT b'01;10' col; AS x CREATE FUNCTION uccp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;/*;$x$;'; END , --end;$x$;
|  col;AS x $$;;;END ;$$ 1; col X'01;10' ,;
< error
| 42601 syntax error at or near "AS"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT b'01;10' col;
< error
| 22P02 ";" is not a valid binary digit
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|  AS x CREATE FUNCTION uccp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;/*;$x$;'; END , --end;$x$;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 2
> exec
| FROM t WHERE a = U&'--;' 1 E''';\\'';' SELECT wf$a$j 1;, "/*;" FROM t WHERE a = CREATE FUNCTION gufl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END 1;
| col zc$w$r col CREATE FUNCTION sbon() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;end;end;'; END + E'\n;\';' +;
| COPY hbo FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'--;' 1 E''';\\'';' SELECT wf$a$j 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'--;' 1 E''';\\'';' SELECT wf$a$j 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + E'\n;\\' SELECT;
| FROM t WHERE a = /*/*/*select 1;select 1;*/*/*/ AS x;SELECT U&'end;$x$;' 1 U&'end;end;' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\n;\\' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\n;\\' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY box FROM stdin;
| a;b
| a;b
| \.
|
| FROM t WHERE a = ";;$x$;;" ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY box FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY box FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , U&'END ;a;b--;' AS x;
| COPY hjr FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'END ;a;b--;' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'END ;a;b--;' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + yy$u$i 1 "a;b);($x$;" , --$x$;--;
|  ,; -- trailing note;
| COPY ntf FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| COPY uhf FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + yy$u$i 1 "a;b);($x$;" , --$x$;--;
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + yy$u$i 1 "a;b);($x$;" , --$x$;--;
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = --a;b);($x$;
|  col e'\n;'';\\' ,;
| , U&'$;$' FROM t WHERE a =; COPY cww FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| COPY ygk FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --a;b);($x$;
|  col e'\n;'';\\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --a;b);($x$;
|  col e'\n;'';\\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY fkz FROM stdin;
| a;b
| c;d;e
| \.
|
| AS x CREATE FUNCTION daot() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$x$;'; END 1 (SELECT 1; SELECT 2) col;
| COPY uij FROM stdin;
| 1	\N
| c;d;e
| \.
| COPY pjo FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fkz FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fkz FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY lwn FROM stdin;
| c;d;e
| a;b
| \.
|  COPY yny FROM stdin;
| 1	\N
| a;b
| \.
|
| 1 CREATE FUNCTION bbjd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END col --; end;
|  AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
| FROM t WHERE a = (SELECT 1; SELECT 2) , E'\\\\\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lwn FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lwn FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x E'\';'';\';' AS x /*/*;;*/*/ FROM t WHERE a = ns$m$b FROM t WHERE a =; , ----;--;--;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x E'\';'';\';' AS x /*/*;;*/*/ FROM t WHERE a = ns$m$b FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x E'\';'';\';' AS x /*/*;;*/*/ FROM t WHERE a = ns$m$b FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY brt FROM stdin;
| a;b
| a;b
| \.
|
| FROM t WHERE a = $t$END ;$t$ col CREATE FUNCTION tpij() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$x$;a;b'; SELECT 1; END col '$x$;/*;' AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY brt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY brt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + E''';\';'';' 1; COPY srw FROM stdin;
| c;d;e
| a;b
| \.
|
| + go$m$w 1 '/*;end;$x$;' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E''';\';'';' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E''';\';'';' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY wrb FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a = x'01;10' SELECT e'\n;\\\\' 1; , /*/*/ *;a;b*/*/ FROM t WHERE a = E'\';\n;'';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wrb FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wrb FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY bak FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bak FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bak FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + e'\n;' col "$;$;" col; COPY wm_ FROM stdin;
| a;b
| a;b
| \.
|
| AS x --end;; 
|  col; , (SELECT 1; SELECT 2) + gq$q$c AS x --$;$
|  SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + e'\n;' col "$;$;" col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + e'\n;' col "$;$;" col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x U&'$x$;select 1;' , $bdl$END ;$x$;/*;$bdl$ SELECT U&'end;select 1;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'$x$;select 1;' , $bdl$END ;$x$;/*;$bdl$ SELECT U&'end;select 1;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'$x$;select 1;' , $bdl$END ;$x$;/*;$bdl$ SELECT U&'end;select 1;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 E'\';\n;' AS x (SELECT 1; SELECT 2) ,;
| col --$x$;$;$END ;
|  col;
| , E''';'';\n;' col;
| 1 /*/*/*; */*/*/ AS x "a;b" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 E'\';\n;' AS x (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 E'\';\n;' AS x (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = pn$q$c + "$x$;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = pn$q$c + "$x$;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = pn$q$c + "$x$;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = /*/*/*;;*/*/*/ AS x --a;ba;b$x$;
|  SELECT $$;;;END ;$$ FROM t WHERE a =;
| 1 ';select 1;);(' + "select 1;; $x$;" col /*/*; );(*/*/ col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*/*/*;;*/*/*/ AS x --a;ba;b$x$;
|  SELECT $$;;;END ;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*/*/*;;*/*/*/ AS x --a;ba;b$x$;
|  SELECT $$;;;END ;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + (SELECT 1; SELECT 2) SELECT; -- trailing note;
| 1 E''';\\\\' col; + U&'a;b;;;' col; COPY exi FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = X'01;10' +;
| + --$;$--;);(
|  FROM t WHERE a = /*/*/*a;b$x$;*/*/*/ FROM t WHERE a =;
| col vx$a$i col U&'/*;end;' 1;
| + U&';);(/*;' AS x --/*;select 1;
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = X'01;10' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = X'01;10' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) +;
| FROM t WHERE a = X'01;10' AS x $$$;$;$$ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY yak FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yak FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yak FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , e'\n;' 1; -- trailing note;
| SELECT X'01;10' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\n;' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\n;' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col /*/*/*$;$*/*/*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*/*$;$*/*/*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*/*$;$*/*/*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT ";;/*;" col;
| COPY qnx FROM stdin;
| a;b
| c;d;e
| \.
|
| + "/*;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 42
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ";;/*;" col;
< error
| 42703 column ";;/*;" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| COPY qnx FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 23
> exec
| COPY yaz FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| AS x _l$p$h 1 U&'$x$;' SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yaz FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yaz FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col py$w$g + td$n$g 1;COPY eii FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 ----;end;end;
|  ,;
| col CREATE FUNCTION pnxv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;$;$/*;'; END col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col py$w$g + td$n$g 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col py$w$g + td$n$g 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY uif FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY ejg FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY sct FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY uif FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY uif FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = x'01;10' 1 U&'$x$;END ;' ,;
| COPY nus FROM stdin;
| a;b
| c;d;e
| \.
|
| col U&'a;b' FROM t WHERE a =; -- trailing note;
| , $_h$a;bend;$_x$$_h$ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' 1 U&'$x$;END ;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' 1 U&'$x$;END ;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT /*/*/*$;$END ;;*/*/*/ , U&'/*;end;' FROM t WHERE a =;
| FROM t WHERE a = e'\';'';' 1; , U&'a;b--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 30
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*/*$;$END ;;*/*/*/ , U&'/*;end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 30
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*/*/*$;$END ;;*/*/*/ , U&'/*;end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 30
> exec
| 1 U&'$x$;end;--;' AS x B'01;10' col --END ;/*;--;
|  ,; -- trailing note;
| SELECT CREATE FUNCTION gctf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END col E'\';\';\\' col 'a;bselect 1;END ;' +;
| 1 /*/*/*);(*/*/*/ + xc$d$k AS x CREATE FUNCTION hdml() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END FROM t WHERE a =;
| + /*/*/*END ;select 1;*/*/*/ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'$x$;end;--;' AS x B'01;10' col --END ;/*;--;
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'$x$;end;--;' AS x B'01;10' col --END ;/*;--;
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT "$;$" + X'01;10' SELECT;
< error
| 42703 column "$;$" does not exist
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "$;$" + X'01;10' SELECT;
< error
| 42703 column "$;$" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _qf FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _qf FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _qf FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col $$);(select 1;$;$$$ FROM t WHERE a = '$x$;$x$;\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$);(select 1;$;$$$ FROM t WHERE a = '$x$;$x$;\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$);(select 1;$;$$$ FROM t WHERE a = '$x$;$x$;\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = B'01;10' + U&'/*;' FROM t WHERE a = "END ;$x$;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = B'01;10' + U&'/*;' FROM t WHERE a = "END ;$x$;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = B'01;10' + U&'/*;' FROM t WHERE a = "END ;$x$;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY uji FROM stdin;
| a;b
| a;b
| \.
|
| SELECT E'\\\n;\n;' +;
| COPY c_m FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY uji FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY uji FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ftm FROM stdin;
| a;b
| c;d;e
| \.
|
| + e'\\\\' ,;COPY _jt FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| FROM t WHERE a = "END ;END ;);(" , e'\';\';\\' col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ftm FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ftm FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , "END ;" , /*/*$;$);(a;b*/*/ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "END ;" , /*/*$;$);(a;b*/*/ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "END ;" , /*/*$;$);(a;b*/*/ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 /*; $;$select 1;*/ FROM t WHERE a = ng$k$b SELECT;
| COPY hfq FROM stdin;
| 1	\N
| c;d;e
| \.
|
| , "$;$;" FROM t WHERE a =;
| COPY oad FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*; $;$select 1;*/ FROM t WHERE a = ng$k$b SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*; $;$select 1;*/ FROM t WHERE a = ng$k$b SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = e'\\' FROM t WHERE a = E'\n;\';' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = e'\\' FROM t WHERE a = E'\n;\';' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = e'\\' FROM t WHERE a = E'\n;\';' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = av$x$q AS x e''';' col U&'end;end;select 1;' col;
| + E'\n;\\' AS x (SELECT 1; SELECT 2) , E'\';' 1; -- trailing note;
| FROM t WHERE a = U&'END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = av$x$q AS x e''';' col U&'end;end;select 1;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = av$x$q AS x e''';' col U&'end;end;select 1;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , b'01;10' AS x CREATE FUNCTION xern() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;END ;'; SELECT 1; END +;
| COPY zet FROM stdin;
| c;d;e
| a;b
| \.
|  FROM t WHERE a = CREATE FUNCTION mwwd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';/*;'; SELECT 1; END + CREATE FUNCTION leys() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;select 1;$x$;'; SELECT 1; END , zv$_$d FROM t WHERE a =; col (SELECT 1; SELECT 2) AS x td$c$e AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , b'01;10' AS x CREATE FUNCTION xern() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;END ;'; SELECT 1; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , b'01;10' AS x CREATE FUNCTION xern() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;END ;'; SELECT 1; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col ';;\' 1 B'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ';;\' 1 B'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ';;\' 1 B'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 --a;b--;--;
|  FROM t WHERE a = "/*;" AS x CREATE FUNCTION __hx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;;'; END ,;
| , --;END ;; 
|  FROM t WHERE a = e'\';' + "; " AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --a;b--;--;
|  FROM t WHERE a = "/*;" AS x CREATE FUNCTION __hx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;;'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --a;b--;--;
|  FROM t WHERE a = "/*;" AS x CREATE FUNCTION __hx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;;'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY krx FROM stdin;
| a;b
| a;b
| \.
|
| 1 --select 1;
|  , U&'a;ba;b' +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY krx FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY krx FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 e'\n;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\n;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\n;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = e''';' FROM t WHERE a = (SELECT 1; SELECT 2) 1;
| COPY emz FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = e''';' FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = e''';' FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 /*/*/*$x$;END ;;;*/*/*/ + --select 1;; 
|  SELECT x'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*$x$;END ;;;*/*/*/ + --select 1;; 
|  SELECT x'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*$x$;END ;;;*/*/*/ + --select 1;; 
|  SELECT x'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x U&'; END ;' col U&'/*;' 1;
| COPY ony FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'; END ;' col U&'/*;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'; END ;' col U&'/*;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 ';;\' SELECT --END ;/*;
|  ,;COPY one FROM stdin;
| a;b
| a;b
| \.
|
| + B'01;10' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ';;\' SELECT --END ;/*;
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ';;\' SELECT --END ;/*;
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + U&'$;$end;select 1;' FROM t WHERE a = /*a;b$x$;end;*/ AS x U&'--;' 1; -- trailing note;
| FROM t WHERE a = "a;b--;" FROM t WHERE a = _y$n$i FROM t WHERE a =; -- trailing note;
| FROM t WHERE a = mf$t$r SELECT;
| 1 --$x$;$x$;;
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'$;$end;select 1;' FROM t WHERE a = /*a;b$x$;end;*/ AS x U&'--;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'$;$end;select 1;' FROM t WHERE a = /*a;b$x$;end;*/ AS x U&'--;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + ";END ;" , ';\' +;
| AS x U&'a;b; ' FROM t WHERE a = e''';\n;\\' col (SELECT 1; SELECT 2) +;
| 1 "--;" AS x;
| 1 e''';\n;' FROM t WHERE a = $jit$$x$;; $jit$ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ";END ;" , ';\' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ";END ;" , ';\' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , "a;b);(" 1 CREATE FUNCTION ddpc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;'; SELECT 1; END 1 ----;
|  1; FROM t WHERE a = CREATE FUNCTION sbuc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END FROM t WHERE a =; -- trailing note;
| 1 (SELECT 1; SELECT 2) FROM t WHERE a = U&'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "a;b);(" 1 CREATE FUNCTION ddpc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;'; SELECT 1; END 1 ----;
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "a;b);(" 1 CREATE FUNCTION ddpc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;'; SELECT 1; END 1 ----;
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x x'01;10' SELECT;
| SELECT --/*;
|  AS x;
| AS x $$$;@x$;$$ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x x'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x x'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x /*$x$;$;$*/ AS x '$x$;a;b' FROM t WHERE a =; + (SELECT 1; SELECT 2) ,;+ CREATE FUNCTION vhlc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; SELECT 1; END FROM t WHERE a = "select 1;" 1 e''';\n;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*$x$;$;$*/ AS x '$x$;a;b' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*$x$;$;$*/ AS x '$x$;a;b' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 U&'END ;select 1;' FROM t WHERE a =;col /*/*/*$;$*/*/*/ + dk$c$p SELECT U&'$;$' 1;, CREATE FUNCTION pocz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; END 1 /*select 1;*/ ,;
| COPY tx_ FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'END ;select 1;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'END ;select 1;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = 'end;--;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = 'end;--;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = 'end;--;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY aom FROM stdin;
| a;b
| c;d;e
| \.
|  -- trailing note;
| + U&'a;b' 1 ');($x$;end;\' AS x;
| + CREATE FUNCTION dzoq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END FROM t WHERE a =;
| , $$end;$$ SELECT X'01;10' SELECT '; ' +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY aom FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY aom FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 '$x$;a;b--;' SELECT U&';;' + "$x$;end;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '$x$;a;b--;' SELECT U&';;' + "$x$;end;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '$x$;a;b--;' SELECT U&';;' + "$x$;end;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT --;;; --;
|  1 U&');(' SELECT;
| FROM t WHERE a = /*/*/*$;$*/*/*/ AS x $thin$);(--;;$thix$$thin$ SELECT B'01;10' col; -- trailing note;
| col /*/*/ *;end;*/*/ AS x e'\\'';\';' 1 CREATE FUNCTION bsyk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;/*;'; SELECT 1; END SELECT; -- trailing note;
| , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "U&');('"
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --;;; --;
|  1 U&');(' SELECT;
< error
| 42601 syntax error at or near "U&');('"
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --;;; --;
|  1 U&');(' SELECT;
< error
| 42601 syntax error at or near "U&');('"
| position: 21
> exec
| COPY vmb FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vmb FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vmb FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col U&'END ;' ,;COPY ewa FROM stdin;
| a;b
| c;d;e
| \.
| COPY ids FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'END ;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'END ;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , --);(
|  SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --);(
|  SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --);(
|  SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT b'01;10' FROM t WHERE a = xw$l$o ,;
| COPY nzc FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| , /*/*/*/ *;--;a;b*/*/*/ FROM t WHERE a = E'\n;'';'';' FROM t WHERE a =;
| , $$; $;$$$ SELECT;
< error
| 42601 syntax error at or near ","
| position: 41
> exec
| BEGIN
< ok 0 -
> exec
| SELECT b'01;10' FROM t WHERE a = xw$l$o ,;
< error
| 42601 syntax error at or near ","
| position: 41
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT b'01;10' FROM t WHERE a = xw$l$o ,;
< error
| 42601 syntax error at or near ","
| position: 41
> exec
| 1 (SELECT 1; SELECT 2) , --);(;;
|  +;
| AS x b'01;10' 1;
| 1 '$;$' col; SELECT ht$f$h FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) , --);(;;
|  +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) , --);(;;
|  +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , B'01;10' FROM t WHERE a = ";;; " FROM t WHERE a = (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , B'01;10' FROM t WHERE a = ";;; " FROM t WHERE a = (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , B'01;10' FROM t WHERE a = ";;; " FROM t WHERE a = (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY ubf FROM stdin;
| c;d;e
| c;d;e
| \.
|  + $$a;b$$ SELECT ----;select 1;;
|  col --;--;
|  SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ubf FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ubf FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , x'01;10' col $tlgt$END ;; $tlgx$$tlgt$ 1;
| COPY ahz FROM stdin;
| a;b
| a;b
| \.
|
| AS x --a;b;;--;
|  SELECT 'a;b$x$;\' +;1 CREATE FUNCTION vjfp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;; a;b'; SELECT 1; END 1 (SELECT 1; SELECT 2) FROM t WHERE a = "$x$;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , x'01;10' col $tlgt$END ;; $tlgx$$tlgt$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , x'01;10' col $tlgt$END ;; $tlgx$$tlgt$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =; -- trailing note;
| COPY _hd FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 e'\n;\n;'';' + ";" col /*select 1;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 $$a;b$$ 1 CREATE FUNCTION aqwf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; END + ');(' FROM t WHERE a =;
| , (SELECT 1; SELECT 2) col '; \' FROM t WHERE a = U&'$x$;;;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $$a;b$$ 1 CREATE FUNCTION aqwf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; END + ');(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $$a;b$$ 1 CREATE FUNCTION aqwf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; END + ');(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col by$e$t + E'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col by$e$t + E'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col by$e$t + E'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY xyc FROM stdin;
| c;d;e
| c;d;e
| \.
|  1 CREATE FUNCTION mbhs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END col;COPY p_v FROM stdin;
| c;d;e
| a;b
| \.
|  AS x x'01;10' FROM t WHERE a = U&'$;$--;' col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xyc FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xyc FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;, /*/*;--;END ;*/*/ , /*/*--;select 1;*/*/ FROM t WHERE a = CREATE FUNCTION imnt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;;$;$'; END FROM t WHERE a =;COPY bxg FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY qul FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 /*/*$;$/ *;*/*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*$;$/ *;*/*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*$;$/ *;*/*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = e'\\' , /*/*--;select 1;END ;*/*/ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = e'\\' , /*/*--;select 1;END ;*/*/ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = e'\\' , /*/*--;select 1;END ;*/*/ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY ria FROM stdin;
| a;b
| c;d;e
| \.
|
| AS x U&');();(' ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ria FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ria FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY _ig FROM stdin;
| 1	\N
| a;b
| \.
|  COPY acl FROM stdin;
| 1	\N
| a;b
| \.
|
| FROM t WHERE a = x'01;10' ,;
| + jy$b$q SELECT '--;/*;);(\' , $$select 1;;;$$ col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _ig FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _ig FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , (SELECT 1; SELECT 2) + /*/*$x$;select 1;*/*/ col --$x$;--;
|  SELECT;
| COPY _uj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) + /*/*$x$;select 1;*/*/ col --$x$;--;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) + /*/*$x$;select 1;*/*/ col --$x$;--;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT ";;;;END ;" col;
< error
| 42703 column ";;;;END ;" does not exist
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ";;;;END ;" col;
< error
| 42703 column ";;;;END ;" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qdb FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| FROM t WHERE a = 'END ;select 1;);(\' AS x x'01;10' 1 $w$);(; /*;$x$$w$ +;
| col (SELECT 1; SELECT 2) col '$;$' , CREATE FUNCTION yivg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';END ;select 1;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qdb FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qdb FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + /*/*select 1;/ *;; */*/ +;
| FROM t WHERE a = (SELECT 1; SELECT 2) col im$o$y ,;
| COPY scp FROM stdin;
| 1	\N
| c;d;e
| \.
| col X'01;10' , B'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*select 1;/ *;; */*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*select 1;/ *;; */*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) SELECT CREATE FUNCTION tcex() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;; '; END SELECT;
| COPY qrs FROM stdin;
| 1	\N
| a;b
| \.
| AS x zk$s$_ AS x (SELECT 1; SELECT 2) AS x E'\\' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) SELECT CREATE FUNCTION tcex() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;; '; END SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) SELECT CREATE FUNCTION tcex() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;; '; END SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col /*/*--;/ *;*/*/ col ----;END ;
|  , E'\n;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*--;/ *;*/*/ col ----;END ;
|  , E'\n;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*--;/ *;*/*/ col ----;END ;
|  , E'\n;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 CREATE FUNCTION apv_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$; '; SELECT 1; END SELECT 'end;' FROM t WHERE a =;
| col b'01;10' SELECT ca$j$a ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION apv_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$; '; SELECT 1; END SELECT 'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION apv_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$; '; SELECT 1; END SELECT 'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col x'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col x'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col x'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = '--;\' + x'01;10' AS x x'01;10' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = '--;\' + x'01;10' AS x x'01;10' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = '--;\' + x'01;10' AS x x'01;10' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col cr$p$s , '; ; ' FROM t WHERE a =; 1 $$select 1;;;select 1;$$ ,;
| COPY rda FROM stdin;
| a;b
| a;b
| \.
|  + $$a;bEND ;$x$;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col cr$p$s , '; ; ' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col cr$p$s , '; ; ' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , --$;$; );(
|  +;
| 1 B'01;10' SELECT E'\n;' 1 x'01;10' AS x; FROM t WHERE a = e'\n;'';' FROM t WHERE a = X'01;10' 1 X'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --$;$; );(
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --$;$; );(
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY onw FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 "a;b--;" AS x;AS x $$);(/*;);($$ ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY onw FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY onw FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY tnk FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY tnk FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY tnk FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY vid FROM stdin;
| 1	\N
| 1	\N
| \.
|
| col ');($;$' + U&'$;$' FROM t WHERE a =; FROM t WHERE a = (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vid FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vid FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + B'01;10' ,;
| + 'end;;; ' 1;
| SELECT U&';' 1 ");(end;" ,; + ';;' , b'01;10' , /*/*/*END ;*/*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + B'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + B'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col ah$_$q , "end;" 1 ru$g$s +; -- trailing note;
| 1 CREATE FUNCTION evan() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ah$_$q , "end;" 1 ru$g$s +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ah$_$q , "end;" 1 ru$g$s +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT --select 1;; 
|  AS x $$select 1;$$ ,;
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 23
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --select 1;; 
|  AS x $$select 1;$$ ,;
< error
| 42601 syntax error at or near "AS"
| position: 23
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --select 1;; 
|  AS x $$select 1;$$ ,;
< error
| 42601 syntax error at or near "AS"
| position: 23
> exec
| 1 ls$d$s col U&');(a;b' 1 '--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ls$d$s col U&');(a;b' 1 '--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ls$d$s col U&');(a;b' 1 '--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = U&'select 1;$;$);(' + /*/*/*);(*/*/*/ , --; 
|  1;
| , ');(;;$x$;\' + --/*;);(
|  , '$;$$x$;\' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'select 1;$;$);(' + /*/*/*);(*/*/*/ , --; 
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'select 1;$;$);(' + /*/*/*);(*/*/*/ , --; 
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , --; 
|  AS x;SELECT --a;b/*;
|  AS x ");(END ;end;" 1; 1 U&'end;;;--;' FROM t WHERE a = E'\n;\';\\' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --; 
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --; 
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x --END ;
|  FROM t WHERE a = --select 1;select 1;);(
|  1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --END ;
|  FROM t WHERE a = --select 1;select 1;);(
|  1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --END ;
|  FROM t WHERE a = --select 1;select 1;);(
|  1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , (SELECT 1; SELECT 2) col U&';;);(' 1;
| AS x --a;b
|  , E'\';' , $$END ;a;b/*;$$ 1; COPY gks FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| COPY upm FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) col U&';;);(' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) col U&';;);(' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY zpw FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zpw FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zpw FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT ne$s$u +;
| 1 E'\';'';\';' FROM t WHERE a = U&'END ;$x$;; ' + CREATE FUNCTION _zjv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(;$x$;'; END ,;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ne$s$u +;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT ne$s$u +;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| SELECT (SELECT 1; SELECT 2) col;
| COPY ulx FROM stdin;
| a;b
| c;d;e
| \.
|
| 1 X'01;10' 1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| SELECT U&'; $;$; ' col U&'$;$' SELECT; -- trailing note;
| SELECT (SELECT 1; SELECT 2) SELECT U&';;);(select 1;' AS x;FROM t WHERE a = 'end;' , $$/*;$$ AS x;
| FROM t WHERE a = x'01;10' 1 /*/*/*/ *;*/*/*/ , U&'--;;;' +;
< error
| 42601 syntax error at or near "U&'$;$'"
| position: 24
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'; $;$; ' col U&'$;$' SELECT;
< error
| 42601 syntax error at or near "U&'$;$'"
| position: 24
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'; $;$; ' col U&'$;$' SELECT;
< error
| 42601 syntax error at or near "U&'$;$'"
| position: 24
> exec
| SELECT $hh$a;b$x$;$hh$ FROM t WHERE a =;1 U&';;$;$;' ,;FROM t WHERE a = X'01;10' col;
< error
| 42601 syntax error at or near ";"
| position: 40
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $hh$a;b$x$;$hh$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 40
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $hh$a;b$x$;$hh$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 40
> exec
| COPY eiz FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY eiz FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY eiz FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT $xrx$;;END ;@$xrx$ SELECT (SELECT 1; SELECT 2) ,; AS x 'end;; $x$;\' + --;
|  FROM t WHERE a = $oldi$;;);($oldi$ 1; COPY ufh FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "("
| position: 34
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $xrx$;;END ;@$xrx$ SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "("
| position: 34
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $xrx$;;END ;@$xrx$ SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "("
| position: 34
> exec
| COPY cmj FROM stdin;
| 1	\N
| 1	\N
| \.
|
| FROM t WHERE a = '--;select 1;\' ,;col /*$;$select 1;);(*/ +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cmj FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cmj FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 "--;$;$" ,;
| COPY rtq FROM stdin;
| 1	\N
| a;b
| \.
|
| COPY bni FROM stdin;
| a;b
| 1	\N
| \.
| COPY cff FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "--;$;$" ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "--;$;$" ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY gkp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gkp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gkp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x (SELECT 1; SELECT 2) ,;
| , yh$d$d + X'01;10' 1 E'\';\\' FROM t WHERE a =;
| col U&'; ; );(' col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , (SELECT 1; SELECT 2) AS x B'01;10' FROM t WHERE a =;
| COPY ypl FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) AS x B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) AS x B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + U&'END ;' AS x --END ;
|  col --a;b
|  col;
| COPY hzm FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'END ;' AS x --END ;
|  col --a;b
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'END ;' AS x --END ;
|  col --a;b
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT U&'a;b' 1 "/*;" SELECT; -- trailing note;
| + jx$h$b AS x ';' , $$select 1;select 1;$$ ,; -- trailing note;
| col ----;;
|  , 'select 1;\' ,;
< error
| 42601 syntax error at or near "1"
| position: 16
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'a;b' 1 "/*;" SELECT;
< error
| 42601 syntax error at or near "1"
| position: 16
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'a;b' 1 "/*;" SELECT;
< error
| 42601 syntax error at or near "1"
| position: 16
> exec
| + --END ;/*;$x$;
|  FROM t WHERE a = e'\n;' SELECT $teb$select 1;$teb$ 1; -- trailing note;
| SELECT U&'a;b; ' + nz$a$a SELECT B'01;10' +; , /*/*end;);(; */*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --END ;/*;$x$;
|  FROM t WHERE a = e'\n;' SELECT $teb$select 1;$teb$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --END ;/*;$x$;
|  FROM t WHERE a = e'\n;' SELECT $teb$select 1;$teb$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col --);(--;
|  AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --);(--;
|  AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --);(--;
|  AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY _no FROM stdin;
| 1	\N
| 1	\N
| \.
|  + U&'END ;);(' + ";;" AS x;SELECT CREATE FUNCTION _jqh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;$;$'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _no FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _no FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col /*/*;;*/*/ 1;
| COPY ipq FROM stdin;
| c;d;e
| c;d;e
| \.
| FROM t WHERE a = (SELECT 1; SELECT 2) 1 /*/*/*--;--;*/*/*/ FROM t WHERE a =;+ (SELECT 1; SELECT 2) col X'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*;;*/*/ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*;;*/*/ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , ';;$;$\' + $sb$$;$/*;$;$$sx$$sb$ + --END ;/*;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ';;$;$\' + $sb$$;$/*;$;$$sx$$sb$ + --END ;/*;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ';;$;$\' + $sb$$;$/*;$;$$sx$$sb$ + --END ;/*;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , --END ;$;$/*;
|  FROM t WHERE a = e''';\\\';' +;
| + 'END ;/*;' FROM t WHERE a = 'end;' col;COPY vvq FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --END ;$;$/*;
|  FROM t WHERE a = e''';\\\';' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --END ;$;$/*;
|  FROM t WHERE a = e''';\\\';' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = xn$h$j + e'\\\';' ,;
| 1 /*--;*/ AS x "END ;);(" col; FROM t WHERE a = '--;' FROM t WHERE a = E'\\' FROM t WHERE a = U&'select 1;--;' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = xn$h$j + e'\\\';' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = xn$h$j + e'\\\';' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY _el FROM stdin;
| a;b
| 1	\N
| \.
|
| SELECT x'01;10' +; , U&';;select 1;a;b' , X'01;10' , /*/*/*a;b);(end;*/*/*/ ,;
| FROM t WHERE a = X'01;10' AS x /*/*/*end;/ *;*/*/*/ 1 xu$g$x SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _el FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _el FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + E'\';' + x'01;10' , /*;;;;*/ +; -- trailing note;
| , U&'$;$--;;;' AS x (SELECT 1; SELECT 2) AS x;
| AS x CREATE FUNCTION ndfg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;ba;b'; END SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\';' + x'01;10' , /*;;;;*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\';' + x'01;10' , /*;;;;*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY jbd FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jbd FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jbd FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x CREATE FUNCTION rtnd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END ,; -- trailing note;
| SELECT wt$p$q AS x (SELECT 1; SELECT 2) 1 E'\';\n;'';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION rtnd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION rtnd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + $pcp$);($pcp$ + ';\' col /*/*/*END ;*/*/*/ +;1 U&';;;;; ' ,;
| FROM t WHERE a = U&'/*;' +;
| SELECT --$;$
|  SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $pcp$);($pcp$ + ';\' col /*/*/*END ;*/*/*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $pcp$);($pcp$ + ';\' col /*/*/*END ;*/*/*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col ia$b$_ 1;
| , CREATE FUNCTION oovr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; ;;'; END + X'01;10' , E''';' 1;
| COPY dt_ FROM stdin;
| a;b
| 1	\N
| \.
| 1 'select 1;; ; ' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ia$b$_ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ia$b$_ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY gil FROM stdin;
| c;d;e
| a;b
| \.
| FROM t WHERE a = ';END ;--;' 1 ----;select 1;
|  SELECT; col CREATE FUNCTION wrrd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; SELECT 1; END , (SELECT 1; SELECT 2) , dz$t$h FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gil FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gil FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 "select 1;$x$;" FROM t WHERE a =;COPY wzp FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY pbr FROM stdin;
| a;b
| a;b
| \.
|  1 is$f$d AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "select 1;$x$;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "select 1;$x$;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col $$$;$select 1;);($$ FROM t WHERE a = rd$_$q SELECT;
| SELECT "end;" col CREATE FUNCTION krkq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;select 1;$x$;'; END +; -- trailing note;
| AS x --select 1;select 1;
|  + U&'/*;; a;b' SELECT;
| AS x --$;$END ;end;
|  + E'\\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$$;$select 1;);($$ FROM t WHERE a = rd$_$q SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$$;$select 1;);($$ FROM t WHERE a = rd$_$q SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = $$$x$;$$ , $$;;$;$$$ AS x;
| FROM t WHERE a = --$;$$;$; 
|  1 --$;$
|  AS x '$;$);($;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$$x$;$$ , $$;;$;$$$ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$$x$;$$ , $$;;$;$$$ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + /*/*/*);(;;select 1;*/*/*/ 1 /*/*select 1;$;$select 1;*/*/ SELECT;
| COPY hxk FROM stdin;
| 1	\N
| a;b
| \.
|  AS x $$end;END ;$$ SELECT ----;);(;;
|  FROM t WHERE a =; COPY tqd FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*/*);(;;select 1;*/*/*/ 1 /*/*select 1;$;$select 1;*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*/*);(;;select 1;*/*/*/ 1 /*/*select 1;$;$select 1;*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 'select 1;/*;' SELECT ';;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 'select 1;/*;' SELECT ';;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 'select 1;/*;' SELECT ';;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY ms_ FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| SELECT $$;a;b$$ ,; col e'\';\';' + CREATE FUNCTION hqdp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';; '; SELECT 1; END , yu$j$a SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ms_ FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ms_ FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x q_$m$a , 'end;end;a;b\' SELECT;
| + E'\\\';\';' 1 b'01;10' ,;
| col U&'END ;/*;' AS x X'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x q_$m$a , 'end;end;a;b\' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x q_$m$a , 'end;end;a;b\' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY ubl FROM stdin;
| 1	\N
| c;d;e
| \.
|
| , $$; select 1;$$ + "; ;;select 1;" , $$$x$;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ubl FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ubl FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + U&'--;;;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'--;;;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'--;;;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col (SELECT 1; SELECT 2) SELECT vw$i$w col --;end;
|  SELECT;
| SELECT (SELECT 1; SELECT 2) + ----;; ;;
|  col;
| + ");(;" , /*$x$;select 1;*/ 1 (SELECT 1; SELECT 2) ,; 1 'a;b' col yo$d$p SELECT U&'--;/*;$;$' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) SELECT vw$i$w col --;end;
|  SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) SELECT vw$i$w col --;end;
|  SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) +;AS x "END ;--;$x$;" , U&'--;' ,;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| COPY skg FROM stdin;
| c;d;e
| 1	\N
| \.
|
| + E'\\' SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY skg FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY skg FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + X'01;10' + CREATE FUNCTION knhi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $;$'; END AS x E'\';' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + X'01;10' + CREATE FUNCTION knhi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $;$'; END AS x E'\';' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + X'01;10' + CREATE FUNCTION knhi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $;$'; END AS x E'\';' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + --);(;
|  col; COPY wwp FROM stdin;
| a;b
| 1	\N
| \.
| AS x U&';;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --);(;
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --);(;
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = "$;$);(a;b" 1 'END ;' 1;
| 1 $$--;a;b$$ col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "$;$);(a;b" 1 'END ;' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "$;$);(a;b" 1 'END ;' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY ajq FROM stdin;
| c;d;e
| c;d;e
| \.
|
| 1 $px$END ;;;@$px$ 1 e'\';' FROM t WHERE a = ');($x$;);(' 1;
| + CREATE FUNCTION _gmt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;/*;$x$;'; END + --$x$;END ;
|  1 ';' ,;AS x $$$;$end;$$ SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ajq FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ajq FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , b'01;10' ,;FROM t WHERE a = od$q$e + U&'select 1;;;' +;
| COPY kbz FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , b'01;10' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , b'01;10' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , --select 1;
|  FROM t WHERE a = CREATE FUNCTION yfkc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END col;
| COPY qqg FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --select 1;
|  FROM t WHERE a = CREATE FUNCTION yfkc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --select 1;
|  FROM t WHERE a = CREATE FUNCTION yfkc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x 'end;;;end;' , U&'$x$;);(a;b' col CREATE FUNCTION xmhb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;/*;'; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x 'end;;;end;' , U&'$x$;);(a;b' col CREATE FUNCTION xmhb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;/*;'; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x 'end;;;end;' , U&'$x$;);(a;b' col CREATE FUNCTION xmhb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;/*;'; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT --;;a;b$;$
|  SELECT CREATE FUNCTION rweu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; END + X'01;10' 1; SELECT /*/*--;*/*/ , U&'$;$' 1;FROM t WHERE a = hr$u$r col --/*;; ;
|  , B'01;10' AS x;
< error
| 42601 syntax error at or near "SELECT"
| position: 20
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --;;a;b$;$
|  SELECT CREATE FUNCTION rweu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; END + X'01;10' 1;
< error
| 42601 syntax error at or near "SELECT"
| position: 20
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --;;a;b$;$
|  SELECT CREATE FUNCTION rweu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; END + X'01;10' 1;
< error
| 42601 syntax error at or near "SELECT"
| position: 20
> exec
| SELECT E'\';\\' FROM t WHERE a =;+ B'01;10' col;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| BEGIN
< ok 0 -
> exec
| SELECT E'\';\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT E'\';\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| COPY dnx FROM stdin;
| a;b
| c;d;e
| \.
|  COPY aah FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY csm FROM stdin;
| 1	\N
| c;d;e
| \.
| COPY udm FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dnx FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dnx FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , e'\';\';\n;' col --;;
|  ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\';\';\n;' col --;;
|  ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\';\';\n;' col --;;
|  ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col /*/*/*);(*/*/*/ , $$$x$;;;$;$$$ SELECT /*$;$*/ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*/*);(*/*/*/ , $$$x$;;;$;$$$ SELECT /*$;$*/ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*/*);(*/*/*/ , $$$x$;;;$;$$$ SELECT /*$;$*/ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col CREATE FUNCTION opbk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;; $x$;'; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION opbk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;; $x$;'; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION opbk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;; $x$;'; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 "$x$;$x$;" col --END ;
|  SELECT '$;$;' +;
| AS x $dgcj$/*;/*;; $dgcx$$dgcj$ + $$; select 1;$;$$$ SELECT; -- trailing note;
| col B'01;10' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "$x$;$x$;" col --END ;
|  SELECT '$;$;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "$x$;$x$;" col --END ;
|  SELECT '$;$;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY ggw FROM stdin;
| c;d;e
| 1	\N
| \.
|
| + (SELECT 1; SELECT 2) col /*/*/*select 1;$x$;*/*/*/ col (SELECT 1; SELECT 2) ,;
| col U&'end;$;$' SELECT ');($;$--;\' FROM t WHERE a = "END ;a;bEND ;" col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ggw FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ggw FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = U&'; ' col;
| AS x e'\n;' SELECT; 1 tc$j$s + fu$y$b ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'; ' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'; ' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY yvo FROM stdin;
| a;b
| 1	\N
| \.
|  COPY yu_ FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yvo FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yvo FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , ";;; ; " FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ";;; ; " FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ";;; ; " FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col $t$$;$$t$ SELECT;
| col --);(END ;--;
|  col "end;end;a;b" 1; AS x "/*;;$x$;" , tg$n$w AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $t$$;$$t$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $t$$;$$t$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT $nkf$);(a;bend;$nkf$ col;
< ok 1 -
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $nkf$);(a;bend;$nkf$ col;
< ok 1 -
> exec
| ROLLBACK
< ok 0 -
> exec
| , x'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , x'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , x'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY whf FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| col --;
|  +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY whf FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY whf FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = $pj$a;bEND ;;;$pj$ SELECT th$b$f FROM t WHERE a = /*;);(*/ 1;+ e'\\'';' + b'01;10' SELECT U&');(' +;
| FROM t WHERE a = "--;--;);(" ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $pj$a;bEND ;;;$pj$ SELECT th$b$f FROM t WHERE a = /*;);(*/ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $pj$a;bEND ;;;$pj$ SELECT th$b$f FROM t WHERE a = /*;);(*/ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 wu$i$g col /*$x$;; a;b*/ ,; col CREATE FUNCTION ___u() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;end;$x$;'; SELECT 1; END FROM t WHERE a =; -- trailing note;
| , (SELECT 1; SELECT 2) AS x U&';' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 wu$i$g col /*$x$;; a;b*/ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 wu$i$g col /*$x$;; a;b*/ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 --a;b;;
|  SELECT;
| 1 b'01;10' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --a;b;;
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --a;b;;
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + $$--;END ;--;$$ + (SELECT 1; SELECT 2) 1 --$x$;$x$;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $$--;END ;--;$$ + (SELECT 1; SELECT 2) 1 --$x$;$x$;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $$--;END ;--;$$ + (SELECT 1; SELECT 2) 1 --$x$;$x$;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col /*/*);(; $;$*/*/ col CREATE FUNCTION zyb_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;'; SELECT 1; END ,;
| SELECT CREATE FUNCTION bvyj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$$;$'; END 1 "--;end;$x$;" 1;
| 1 ul$r$m SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*);(; $;$*/*/ col CREATE FUNCTION zyb_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*);(; $;$*/*/ col CREATE FUNCTION zyb_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY vuz FROM stdin;
| 1	\N
| 1	\N
| \.
|  , U&'a;bselect 1;' + $$;$$ , $mnt$a;b$mnx$$mnt$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vuz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vuz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY q_r FROM stdin;
| 1	\N
| 1	\N
| \.
|  AS x X'01;10' FROM t WHERE a = ';;;--;\' + CREATE FUNCTION atea() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;b'; END FROM t WHERE a =; -- trailing note;
| , CREATE FUNCTION wcqm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(--;'; END col /*END ;; ;*/ FROM t WHERE a = rw$m$w +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY q_r FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY q_r FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , /*/*--;; ; */*/ col (SELECT 1; SELECT 2) FROM t WHERE a = X'01;10' +;col '; );(' + $$end;$$ +;
| FROM t WHERE a = e''';' FROM t WHERE a =;
| 1 _u$w$j + x'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*--;; ; */*/ col (SELECT 1; SELECT 2) FROM t WHERE a = X'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*--;; ; */*/ col (SELECT 1; SELECT 2) FROM t WHERE a = X'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT /*/*/*;;a;b; */*/*/ SELECT CREATE FUNCTION sugo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';--;'; END , "a;b$;$select 1;" SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 28
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*/*;;a;b; */*/*/ SELECT CREATE FUNCTION sugo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';--;'; END , "a;b$;$select 1;" SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 28
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*/*/*;;a;b; */*/*/ SELECT CREATE FUNCTION sugo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';--;'; END , "a;b$;$select 1;" SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 28
> exec
| + ';;;END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ';;;END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ';;;END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = CREATE FUNCTION i_mb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b;;'; SELECT 1; END SELECT;
| COPY mzg FROM stdin;
| a;b
| a;b
| \.
|
| SELECT E'\\' + CREATE FUNCTION ydyv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;ba;bselect 1;'; END AS x; -- trailing note;
| COPY xmb FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION i_mb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b;;'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION i_mb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b;;'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT x'01;10' FROM t WHERE a = "end;" AS x; -- trailing note;
| FROM t WHERE a = --a;b
|  FROM t WHERE a =; -- trailing note;
| FROM t WHERE a = --/*;;;
|  1 /*/*/*$x$;*/*/*/ ,;
< error
| 42601 syntax error at or near "AS"
| position: 41
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' FROM t WHERE a = "end;" AS x;
< error
| 42601 syntax error at or near "AS"
| position: 41
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT x'01;10' FROM t WHERE a = "end;" AS x;
< error
| 42601 syntax error at or near "AS"
| position: 41
> exec
| 1 $_pkk$;;/*;$_pkk$ FROM t WHERE a =;
| + e''';\';' AS x U&'; select 1;' col; -- trailing note;
| COPY oc_ FROM stdin;
| a;b
| a;b
| \.
| , "end;" + /*/*/*$x$;; $x$;*/*/*/ SELECT b'01;10' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $_pkk$;;/*;$_pkk$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $_pkk$;;/*;$_pkk$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY yio FROM stdin;
| 1	\N
| a;b
| \.
|
| + qg$j$p FROM t WHERE a = E''';\\'';' +; -- trailing note;
| FROM t WHERE a = av$h$o SELECT --end;/*;
|  SELECT --/*;;a;b
|  1;
| col X'01;10' col "--;end;select 1;" 1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yio FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yio FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , "select 1;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "select 1;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "select 1;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x x'01;10' AS x B'01;10' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x x'01;10' AS x B'01;10' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x x'01;10' AS x B'01;10' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , CREATE FUNCTION gfwl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END SELECT --END ;$x$;select 1;
|  SELECT;
| COPY byh FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION gfwl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END SELECT --END ;$x$;select 1;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION gfwl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END SELECT --END ;$x$;select 1;
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY sl_ FROM stdin;
| a;b
| a;b
| \.
|
| SELECT $$/*;$$ col B'01;10' FROM t WHERE a = B'01;10' +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sl_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sl_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , $$$;@x$;$$ col U&');(; );(' FROM t WHERE a = fi$d$o AS x; col $$$x$;;;END ;$$ SELECT B'01;10' + kb$z$a SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $$$;@x$;$$ col U&');(; );(' FROM t WHERE a = fi$d$o AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $$$;@x$;$$ col U&');(; );(' FROM t WHERE a = fi$d$o AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 "$;$);($x$;" col rq$g$q SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "$;$);($x$;" col rq$g$q SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "$;$);($x$;" col rq$g$q SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col ";;--;; " SELECT "$;$" AS x;
| , CREATE FUNCTION xxzp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;END ;'; SELECT 1; END , wj$b$t FROM t WHERE a =; -- trailing note;
| + B'01;10' SELECT U&'end;end;' AS x mu$_$l col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ";;--;; " SELECT "$;$" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ";;--;; " SELECT "$;$" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY slq FROM stdin;
| 1	\N
| a;b
| \.
|
| 1 /*/*/*;;; */*/*/ SELECT;
| COPY bhj FROM stdin;
| a;b
| c;d;e
| \.
| SELECT e'\n;\';\n;' , U&'select 1;select 1;/*;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY slq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY slq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 e'\\' , X'01;10' col 'END ;' 1; -- trailing note;
| + /*/*/*;*/*/*/ + 'END ;/*;a;b' AS x; COPY nhc FROM stdin;
| c;d;e
| 1	\N
| \.
|
| + (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\\' , X'01;10' col 'END ;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\\' , X'01;10' col 'END ;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x $v$--;$x$$v$ , ----;
|  col; FROM t WHERE a = $wdn$END ;$wdn$ FROM t WHERE a = $$;$;$--;$$ 1 U&'end;$;$' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $v$--;$x$$v$ , ----;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $v$--;$x$$v$ , ----;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = ";;select 1;" FROM t WHERE a = $g$a;b;$g$ FROM t WHERE a = CREATE FUNCTION zqpb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;'; END +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ";;select 1;" FROM t WHERE a = $g$a;b;$g$ FROM t WHERE a = CREATE FUNCTION zqpb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;'; END +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ";;select 1;" FROM t WHERE a = $g$a;b;$g$ FROM t WHERE a = CREATE FUNCTION zqpb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;'; END +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT $$a;b$$ FROM t WHERE a = U&'select 1;a;b' AS x 'a;b;end;\' SELECT; -- trailing note;
| col "; " col;col $$select 1;$$ AS x ';;; end;' ,;
< error
| 42601 syntax error at or near "AS"
| position: 50
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$a;b$$ FROM t WHERE a = U&'select 1;a;b' AS x 'a;b;end;\' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 50
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$a;b$$ FROM t WHERE a = U&'select 1;a;b' AS x 'a;b;end;\' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 50
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) SELECT --a;b$x$;; 
|  +; + U&'$;$);(END ;' FROM t WHERE a =; -- trailing note;
| + $jba$a;bend;$;$$jba$ 1 ";;" FROM t WHERE a = '$x$;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) SELECT --a;b$x$;; 
|  +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) SELECT --a;b$x$;; 
|  +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY cmq FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY jdf FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cmq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cmq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x U&'a;b' col --$;$--;;;
|  + _j$v$e SELECT;
| 1 '; ;;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'a;b' col --$;$--;;;
|  + _j$v$e SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'a;b' col --$;$--;;;
|  + _j$v$e SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , $$/*;end;$$ +;
| SELECT ';; ' col; -- trailing note;
| + $yrmb$$x$;$yrmx$$yrmb$ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $$/*;end;$$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $$/*;end;$$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x /*/*/*a;b*/*/*/ + (SELECT 1; SELECT 2) 1;
| + $uthw$select 1;a;b$x$;$uthx$$uthw$ AS x;
| + E''';\\' , ky$y$t FROM t WHERE a = U&'; ' SELECT;
| AS x 'a;b$x$;a;b' col e'\n;\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*a;b*/*/*/ + (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*a;b*/*/*/ + (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + /*select 1;;end;*/ FROM t WHERE a =;SELECT CREATE FUNCTION jhsd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END col lo$y$x , x'01;10' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*select 1;;end;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*select 1;;end;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col "END ;a;b" + '$x$;' FROM t WHERE a = '/*;; );(\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "END ;a;b" + '$x$;' FROM t WHERE a = '/*;; );(\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "END ;a;b" + '$x$;' FROM t WHERE a = '/*;; );(\' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 /*/*/*;;$;$*/*/*/ SELECT $zqkr$;$zqkx$$zqkr$ col;
| , "; $;$" +;
| COPY _we FROM stdin;
| 1	\N
| a;b
| \.
|  -- trailing note;
| , (SELECT 1; SELECT 2) SELECT $xyk$/*;a;b$xyx$$xyk$ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*;;$;$*/*/*/ SELECT $zqkr$;$zqkx$$zqkr$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*;;$;$*/*/*/ SELECT $zqkr$;$zqkx$$zqkr$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT U&';;);(; ' , (SELECT 1; SELECT 2) FROM t WHERE a = E'\n;\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 31
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&';;);(; ' , (SELECT 1; SELECT 2) FROM t WHERE a = E'\n;\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 31
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&';;);(; ' , (SELECT 1; SELECT 2) FROM t WHERE a = E'\n;\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 31
> exec
| AS x ');(;; ' FROM t WHERE a =;AS x e'\\' + --;$;$/*;
|  +; , (SELECT 1; SELECT 2) FROM t WHERE a = X'01;10' 1;
| COPY hfj FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ');(;; ' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ');(;; ' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col pu$b$r AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col pu$b$r AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col pu$b$r AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT e'\\' SELECT U&'select 1;$x$;' + /*/*/*;;; */*/*/ ,;
| COPY cez FROM stdin;
| a;b
| a;b
| \.
|
| 1 ";" + ");(a;b$x$;" AS x; -- trailing note;
| COPY oro FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "U&'select 1;$x$;'"
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\\' SELECT U&'select 1;$x$;' + /*/*/*;;; */*/*/ ,;
< error
| 42601 syntax error at or near "U&'select 1;$x$;'"
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT e'\\' SELECT U&'select 1;$x$;' + /*/*/*;;; */*/*/ ,;
< error
| 42601 syntax error at or near "U&'select 1;$x$;'"
| position: 21
> exec
| FROM t WHERE a = "select 1;--;; " AS x $v$--;$v$ SELECT;
| COPY asg FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "select 1;--;; " AS x $v$--;$v$ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "select 1;--;; " AS x $v$--;$v$ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT e'\\\\' SELECT;
| col CREATE FUNCTION ywvg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;'; SELECT 1; END AS x E'\';' AS x;col 'select 1;' +;
< error
| 42601 syntax error at or near "col"
| position: 24
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\\\\' SELECT;
< ok 1 -
> exec
|
| col CREATE FUNCTION ywvg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;'; SELECT 1; END AS x E'\';' AS x;
< error
| 42601 syntax error at or near "col"
| position: 2
> exec
| ROLLBACK
< ok 0 -
> exec
|
| col CREATE FUNCTION ywvg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;'; SELECT 1; END AS x E'\';' AS x;
< error
| 42601 syntax error at or near "col"
| position: 2
> exec
| , --select 1;;a;b
|  , x'01;10' 1 $$$;$END ;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --select 1;;a;b
|  , x'01;10' 1 $$$;$END ;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --select 1;;a;b
|  , x'01;10' 1 $$$;$END ;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col --a;b--;$x$;
|  FROM t WHERE a = vx$n$f col x'01;10' SELECT;
| , $veef$a;b$veef$ FROM t WHERE a =; -- trailing note;
| COPY prg FROM stdin;
| 1	\N
| 1	\N
| \.
| COPY pvi FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --a;b--;$x$;
|  FROM t WHERE a = vx$n$f col x'01;10' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --a;b--;$x$;
|  FROM t WHERE a = vx$n$f col x'01;10' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY sya FROM stdin;
| a;b
| 1	\N
| \.
|
| COPY khk FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sya FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sya FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x U&';$x$;' SELECT CREATE FUNCTION jtbo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;$;$'; END + --;select 1;);(
|  FROM t WHERE a =;
| col X'01;10' 1;
| col nm$k$a SELECT U&'END ;;;);(' FROM t WHERE a = vq$i$i AS x; -- trailing note;
| AS x $$);($;$$$ col ";a;b$x$;" AS x U&';;;;select 1;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&';$x$;' SELECT CREATE FUNCTION jtbo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;$;$'; END + --;select 1;);(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&';$x$;' SELECT CREATE FUNCTION jtbo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;$;$'; END + --;select 1;);(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col (SELECT 1; SELECT 2) SELECT x'01;10' AS x;
| AS x E'\\'';' SELECT U&'; $x$;);(' 1; FROM t WHERE a = E'\n;\n;' + 'a;b\' FROM t WHERE a = U&'--;;;END ;' SELECT;
| SELECT b'01;10' FROM t WHERE a = "a;b" +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) SELECT x'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) SELECT x'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + $lfw$select 1;END ;;$lfx$$lfw$ AS x 'select 1;end;\' AS x ";END ;" col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $lfw$select 1;END ;;$lfx$$lfw$ AS x 'select 1;end;\' AS x ";END ;" col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $lfw$select 1;END ;;$lfx$$lfw$ AS x 'select 1;end;\' AS x ";END ;" col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY bxm FROM stdin;
| a;b
| 1	\N
| \.
|  FROM t WHERE a = CREATE FUNCTION ndzi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END AS x ----;; $x$;
|  +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bxm FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bxm FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col CREATE FUNCTION wkiz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;a;b'; END + U&');(select 1;/*;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION wkiz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;a;b'; END + U&');(select 1;/*;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION wkiz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;a;b'; END + U&');(select 1;/*;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY ehe FROM stdin;
| 1	\N
| a;b
| \.
|
| COPY rah FROM stdin;
| c;d;e
| a;b
| \.
|  1 dn$u$u SELECT;
| SELECT $$end;$;$; $$ col "$;$$;$" SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ehe FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ehe FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col e'\n;' AS x ";" + '; END ;/*;\' 1;
| + e''';' col $bahc$a;b; $bahx$$bahc$ FROM t WHERE a =;
| + --$;$$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e'\n;' AS x ";" + '; END ;/*;\' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e'\n;' AS x ";" + '; END ;/*;\' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x $$END ;$x$;$$ AS x;COPY lkc FROM stdin;
| a;b
| c;d;e
| \.
|  1 '--;end;' col $s$;;);(; $s$ SELECT de$a$j col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $$END ;$x$;$$ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $$END ;$x$;$$ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col CREATE FUNCTION rcuc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;;'; SELECT 1; END 1; -- trailing note;
| SELECT --/*;/*;$;$
|  col;
| SELECT $$$;@;$end;$$ 1 CREATE FUNCTION keq_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; ;'; END + (SELECT 1; SELECT 2) SELECT;
| AS x CREATE FUNCTION tzjd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;;;'; END col X'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION rcuc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION rcuc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY owq FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY owq FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY owq FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT e'\\\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\\\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT e'\\\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 33
> exec
| col (SELECT 1; SELECT 2) 1 U&'/*;;;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) 1 U&'/*;;;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) 1 U&'/*;;;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 ');(;);(\' col 'end;\' , ";" ,;
| COPY hpg FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| COPY tuh FROM stdin;
| c;d;e
| c;d;e
| \.
|
| AS x (SELECT 1; SELECT 2) , ');(/*;' SELECT CREATE FUNCTION ltot() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;--;--;'; END SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ');(;);(\' col 'end;\' , ";" ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ');(;);(\' col 'end;\' , ";" ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x "END ;" FROM t WHERE a =;
| COPY pmk FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "END ;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "END ;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY gyg FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gyg FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gyg FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY shj FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 U&'--;--;' +; -- trailing note;
| FROM t WHERE a = e'\\\';' col; COPY zfy FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY shj FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY shj FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT E'\\\';' FROM t WHERE a = $$end;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 43
> exec
| BEGIN
< ok 0 -
> exec
| SELECT E'\\\';' FROM t WHERE a = $$end;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 43
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT E'\\\';' FROM t WHERE a = $$end;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 43
> exec
| COPY sal FROM stdin;
| 1	\N
| a;b
| \.
| + E''';\n;\\' +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sal FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sal FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x U&';a;b;;' 1 E'\n;\n;' col; + /*/*select 1;);(; */*/ ,;
| + CREATE FUNCTION njbm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END SELECT /*/*/ *;);(;;*/*/ col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&';a;b;;' 1 E'\n;\n;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&';a;b;;' 1 E'\n;\n;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 $nduj$END ;$ndux$$nduj$ AS x $$; /*;$$ 1 CREATE FUNCTION tgft() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;END ;'; SELECT 1; END FROM t WHERE a =;
| FROM t WHERE a = CREATE FUNCTION skmj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b--;END ;'; END FROM t WHERE a =;
| FROM t WHERE a = $$--;;select 1;$$ , ";$x$;" FROM t WHERE a =;
| col e'\\\n;\';' col /*/*/*a;b$x$;*/*/*/ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $nduj$END ;$ndux$$nduj$ AS x $$; /*;$$ 1 CREATE FUNCTION tgft() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;END ;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $nduj$END ;$ndux$$nduj$ AS x $$; /*;$$ 1 CREATE FUNCTION tgft() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;END ;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + CREATE FUNCTION vmlp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END col /*/*/*;;*/*/*/ +; -- trailing note;
| 1 /*/*/*select 1;end;*/*/*/ SELECT e'\n;\n;' SELECT;
| 1 $$$;$a;bselect 1;$$ 1 "--;--;" AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION vmlp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END col /*/*/*;;*/*/*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION vmlp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END col /*/*/*;;*/*/*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT b'01;10' AS x;
< error
| 22P02 ";" is not a valid binary digit
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT b'01;10' AS x;
< error
| 22P02 ";" is not a valid binary digit
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\n;'';\n;' col "--;" , b'01;10' FROM t WHERE a =;FROM t WHERE a = U&'END ;;;' ,;
| COPY pyk FROM stdin;
| c;d;e
| a;b
| \.
| AS x CREATE FUNCTION h_cq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;select 1;'; SELECT 1; END , $emdi$);(a;b; $emdx$$emdi$ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E'\n;'';\n;' col "--;" , b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\n;'';\n;' col "--;" , b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = ';;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ';;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ';;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY ltw FROM stdin;
| a;b
| c;d;e
| \.
|  -- trailing note;
| 1 B'01;10' 1 (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ltw FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ltw FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY _ng FROM stdin;
| 1	\N
| a;b
| \.
|
| COPY zzq FROM stdin;
| c;d;e
| c;d;e
| \.
|
| AS x --END ;/*;a;b
|  ,;
| COPY suz FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _ng FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _ng FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY aon FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| AS x CREATE FUNCTION kdih() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$);(select 1;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY aon FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY aon FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT "$x$;" SELECT E''';' 1;
| FROM t WHERE a = --/*;;;/*;
|  +; 1 $$END ;$;$$$ SELECT;
| + '/*;/*;;' ,;
< error
| 42601 syntax error at or near "E''';'"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "$x$;" SELECT E''';' 1;
< error
| 42601 syntax error at or near "E''';'"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "$x$;" SELECT E''';' 1;
< error
| 42601 syntax error at or near "E''';'"
| position: 22
> exec
| + "END ;select 1;a;b" col ty$j$y ,; + ----;
|  SELECT e'\n;\';' AS x wy$s$i ,;
| + /*end;;;$;$*/ 1;COPY wzm FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "END ;select 1;a;b" col ty$j$y ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "END ;select 1;a;b" col ty$j$y ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , --$;$END ;
|  1;
| SELECT '$;$' +;
| AS x '$;$end;\' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --$;$END ;
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --$;$END ;
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = --end;
|  SELECT;
| col --a;b
|  + U&'$x$;END ;$;$' FROM t WHERE a = B'01;10' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --end;
|  SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --end;
|  SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) 1;
| SELECT CREATE FUNCTION nglz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;);('; SELECT 1; END + (SELECT 1; SELECT 2) +;
| + CREATE FUNCTION xqdf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END 1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| col 'END ;end;' FROM t WHERE a =; -- trailing note;
| SELECT $$; ;; $$ SELECT '/*;end;' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'END ;end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'END ;end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY cbk FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY xqu FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cbk FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cbk FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col U&'/*;' 1 ft$i$c SELECT /*/*end;END ;*/*/ AS x; SELECT CREATE FUNCTION ithh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END 1 'select 1;;; \' + x'01;10' AS x;
| COPY fez FROM stdin;
| 1	\N
| a;b
| \.
|  COPY bco FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'/*;' 1 ft$i$c SELECT /*/*end;END ;*/*/ AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'/*;' 1 ft$i$c SELECT /*/*end;END ;*/*/ AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 rw$f$u , U&'; --;' SELECT;col X'01;10' AS x /*;;;select 1;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 rw$f$u , U&'; --;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 rw$f$u , U&'; --;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT CREATE FUNCTION ek_b() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$'; END FROM t WHERE a = --END ;--;
|  AS x;
| COPY eyw FROM stdin;
| 1	\N
| 1	\N
| \.
|
| , '$x$;/*;\' + ';;' SELECT; -- trailing note;
| COPY nac FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION ek_b() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$'; END FROM t WHERE a = --END ;--;
|  AS x;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION ek_b() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$'; END FROM t WHERE a = --END ;--;
|  AS x;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| + ');(END ;; ' AS x CREATE FUNCTION gojm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;$x$;select 1;'; SELECT 1; END 1 U&'a;b;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ');(END ;; ' AS x CREATE FUNCTION gojm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;$x$;select 1;'; SELECT 1; END 1 U&'a;b;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ');(END ;; ' AS x CREATE FUNCTION gojm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;$x$;select 1;'; SELECT 1; END 1 U&'a;b;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY the FROM stdin;
| c;d;e
| a;b
| \.
|
| , x'01;10' SELECT /*);(--;$x$;*/ AS x; -- trailing note;
| col U&'$x$;' 1 "$;$end;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY the FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY the FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT rv$o$g FROM t WHERE a = X'01;10' SELECT; , "select 1;;" 1 --$;$
|  FROM t WHERE a = CREATE FUNCTION bbll() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;;'; END FROM t WHERE a =; -- trailing note;
| COPY lmt FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| , B'01;10' ,;
< error
| 42601 syntax error at or near "SELECT"
| position: 41
> exec
| BEGIN
< ok 0 -
> exec
| SELECT rv$o$g FROM t WHERE a = X'01;10' SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 41
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT rv$o$g FROM t WHERE a = X'01;10' SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 41
> exec
| 1 /*/*--;*/*/ SELECT CREATE FUNCTION hgth() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';; '; SELECT 1; END FROM t WHERE a =;
| AS x E'\\\';' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*--;*/*/ SELECT CREATE FUNCTION hgth() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';; '; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*--;*/*/ SELECT CREATE FUNCTION hgth() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';; '; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY trq FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| 1 U&');(;;end;' 1 $$$x$;$;$$$ + '$x$;);(' FROM t WHERE a =;SELECT $$;a;b; $$ SELECT "$;$END ;; " AS x ';;--;end;' AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY trq FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY trq FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = --$;$$;$; 
|  FROM t WHERE a = $cse$--;;$cse$ SELECT;
| COPY unh FROM stdin;
| 1	\N
| 1	\N
| \.
|
| SELECT ");(a;b$;$" + (SELECT 1; SELECT 2) SELECT U&';$x$;' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --$;$$;$; 
|  FROM t WHERE a = $cse$--;;$cse$ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --$;$$;$; 
|  FROM t WHERE a = $cse$--;;$cse$ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + e'\\' FROM t WHERE a =;
| , --$x$;END ;
|  +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + e'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + e'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = U&';;--;select 1;' FROM t WHERE a = /*;$;$*/ col;
| , "; " AS x at$l$g 1 (SELECT 1; SELECT 2) 1; COPY xan FROM stdin;
| 1	\N
| c;d;e
| \.
|
| FROM t WHERE a = E'\n;' AS x $$$x$;$$ + $lve$;;;end;$lve$ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&';;--;select 1;' FROM t WHERE a = /*;$;$*/ col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&';;--;select 1;' FROM t WHERE a = /*;$;$*/ col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT $$;;/*;$$ +;
< error
| 42601 syntax error at or near ";"
| position: 19
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$;;/*;$$ +;
< error
| 42601 syntax error at or near ";"
| position: 19
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$;;/*;$$ +;
< error
| 42601 syntax error at or near ";"
| position: 19
> exec
| + "/*;" col CREATE FUNCTION zbwf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';);('; SELECT 1; END + --;$x$;
|  SELECT;
| FROM t WHERE a = CREATE FUNCTION ppmj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; END , "select 1;END ;" AS x;
| + 'end;$x$;\' FROM t WHERE a = $$a;b;$$ 1;
| COPY pjv FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "/*;" col CREATE FUNCTION zbwf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';);('; SELECT 1; END + --;$x$;
|  SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "/*;" col CREATE FUNCTION zbwf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';);('; SELECT 1; END + --;$x$;
|  SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col 'a;b);(a;b' + ';; a;b' ,;COPY hof FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'a;b);(a;b' + ';; a;b' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'a;b);(a;b' + ';; a;b' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY cab FROM stdin;
| c;d;e
| 1	\N
| \.
|
| 1 E''';' , '; ;;' +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cab FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cab FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY efo FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY efo FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY efo FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = "end;/*;" FROM t WHERE a = B'01;10' SELECT b'01;10' FROM t WHERE a =; -- trailing note;
| , /*/*/ *;);(;*/*/ AS x;
| 1 ";" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "end;/*;" FROM t WHERE a = B'01;10' SELECT b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "end;/*;" FROM t WHERE a = B'01;10' SELECT b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col "$;$;a;b" +; -- trailing note;
| AS x "end;--;" col b'01;10' 1; col gx$s$i 1 nt$i$k + bi$i$s +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "$;$;a;b" +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "$;$;a;b" +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = X'01;10' SELECT CREATE FUNCTION jyqi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = X'01;10' SELECT CREATE FUNCTION jyqi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = X'01;10' SELECT CREATE FUNCTION jyqi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x U&'END ;' col x'01;10' SELECT; + (SELECT 1; SELECT 2) FROM t WHERE a = (SELECT 1; SELECT 2) col CREATE FUNCTION ogyp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END FROM t WHERE a =; COPY v_r FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'END ;' col x'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'END ;' col x'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 $ab$; ; select 1;$ax$$ab$ col;
| COPY tvj FROM stdin;
| a;b
| 1	\N
| \.
|
| FROM t WHERE a = /*/*end;);(a;b*/*/ SELECT E''';' FROM t WHERE a = $$;; $x$;$$ col;COPY nez FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $ab$; ; select 1;$ax$$ab$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $ab$; ; select 1;$ax$$ab$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + X'01;10' SELECT CREATE FUNCTION nfgl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; SELECT 1; END col; -- trailing note;
| 1 X'01;10' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + X'01;10' SELECT CREATE FUNCTION nfgl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + X'01;10' SELECT CREATE FUNCTION nfgl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col U&'select 1;;END ;' FROM t WHERE a =;
| AS x ";;select 1;$;$" col (SELECT 1; SELECT 2) col; -- trailing note;
| SELECT "--;$x$;" col;
| col CREATE FUNCTION ller() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(select 1;$;$'; END +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'select 1;;END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'select 1;;END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 U&'$x$;END ;$;$' SELECT 'end;; \' col U&'$x$;' ,; -- trailing note;
| AS x '--;end;' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'$x$;END ;$;$' SELECT 'end;; \' col U&'$x$;' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'$x$;END ;$;$' SELECT 'end;; \' col U&'$x$;' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 "end;END ;$;$" 1 $n$END ;$x$$n$ AS x;COPY kpw FROM stdin;
| c;d;e
| 1	\N
| \.
| col U&'end;' AS x;
| COPY zxx FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "end;END ;$;$" 1 $n$END ;$x$$n$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "end;END ;$;$" 1 $n$END ;$x$$n$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY _hv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _hv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _hv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT CREATE FUNCTION uxt_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;END ;$x$;'; END + $$/*;$$ +;
| + CREATE FUNCTION b_pw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;end;'; END , CREATE FUNCTION btou() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END 1;
| 1 U&';' SELECT --);($;$
|  col $f$--;END ;$x$$f$ ,;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION uxt_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;END ;$x$;'; END + $$/*;$$ +;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION uxt_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;END ;$x$;'; END + $$/*;$$ +;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| , /*/*/*; $x$;;;*/*/*/ FROM t WHERE a = CREATE FUNCTION cleg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bend;'; SELECT 1; END +; col --select 1;
|  1 "; ; " ,;AS x "END ;" col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*/*; $x$;;;*/*/*/ FROM t WHERE a = CREATE FUNCTION cleg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bend;'; SELECT 1; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*/*; $x$;;;*/*/*/ FROM t WHERE a = CREATE FUNCTION cleg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bend;'; SELECT 1; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT U&'--;' FROM t WHERE a = zp$w$n 1;
| COPY pdp FROM stdin;
| c;d;e
| a;b
| \.
|
| SELECT U&'--;; ;;' , E'\\'';\';' FROM t WHERE a = --;);(
|  SELECT;
| 1 ----;
|  SELECT e'\\\';\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 40
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'--;' FROM t WHERE a = zp$w$n 1;
< error
| 42601 syntax error at or near "1"
| position: 40
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'--;' FROM t WHERE a = zp$w$n 1;
< error
| 42601 syntax error at or near "1"
| position: 40
> exec
| col $$/*;$$ +; -- trailing note;
| FROM t WHERE a = CREATE FUNCTION rxkm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;;/*;'; SELECT 1; END FROM t WHERE a = (SELECT 1; SELECT 2) SELECT E'\';' FROM t WHERE a =; -- trailing note;
| FROM t WHERE a = U&');($x$;);(' SELECT (SELECT 1; SELECT 2) col $osz$select 1;$osz$ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$/*;$$ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$/*;$$ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT "$;$" , U&');(a;b' 1 x'01;10' col;
| COPY t_i FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 27
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "$;$" , U&');(a;b' 1 x'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 27
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "$;$" , U&');(a;b' 1 x'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 27
> exec
| COPY uiz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY uiz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY uiz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col /*/ *;$;$; */ + ";" FROM t WHERE a =;
| + '$;$);(; ' AS x ';;a;bend;' col --END ;/*;
|  ,;
| 1 --$x$;
|  +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/ *;$;$; */ + ";" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/ *;$;$; */ + ";" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 ay$d$x col oh$b$n + /*/ *;*/ ,; 1 "/*;end;$x$;" + (SELECT 1; SELECT 2) , --/*;end;
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ay$d$x col oh$b$n + /*/ *;*/ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ay$d$x col oh$b$n + /*/ *;*/ ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x "$;$" AS x --);(/*;
|  FROM t WHERE a = km$x$_ SELECT;
| col ----;a;b; 
|  col jj$s$p , '/*;;;; ' SELECT;
| + U&'end;a;b' col;
| SELECT $$; $$ FROM t WHERE a = U&'; ' AS x X'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "$;$" AS x --);(/*;
|  FROM t WHERE a = km$x$_ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "$;$" AS x --);(/*;
|  FROM t WHERE a = km$x$_ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + $$;$$ AS x hb$h$l SELECT; 1 b'01;10' ,; -- trailing note;
| SELECT ----;$x$;
|  , e'\';\\\n;' SELECT CREATE FUNCTION ckpx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; /*;END ;'; END AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $$;$$ AS x hb$h$l SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $$;$$ AS x hb$h$l SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) AS x;
| FROM t WHERE a = U&'--;--;' SELECT kq$_$h +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY bgr FROM stdin;
| 1	\N
| 1	\N
| \.
|
| + nl$f$w 1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bgr FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bgr FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 "--;;" col ');($x$;end;' col '; ' 1;COPY oeo FROM stdin;
| 1	\N
| 1	\N
| \.
| col E'\n;' col $$end;;;$$ AS x;
| SELECT $$$;$;/*;$$ 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "--;;" col ');($x$;end;' col '; ' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "--;;" col ');($x$;end;' col '; ' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x --END ;$;$
|  + /*/*select 1;a;b*/*/ , --select 1;
|  AS x; col CREATE FUNCTION topn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;end;'; SELECT 1; END SELECT (SELECT 1; SELECT 2) SELECT --; $;$
|  1;
| + ");(" FROM t WHERE a = /*/*select 1;;;*/*/ 1 X'01;10' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --END ;$;$
|  + /*/*select 1;a;b*/*/ , --select 1;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --END ;$;$
|  + /*/*select 1;a;b*/*/ , --select 1;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x ";;a;b);(" FROM t WHERE a = fo$a$p AS x U&'a;b);(' SELECT;
| , x'01;10' 1 X'01;10' FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a =;
| FROM t WHERE a = e'\n;' FROM t WHERE a = "end;" ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ";;a;b);(" FROM t WHERE a = fo$a$p AS x U&'a;b);(' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ";;a;b);(" FROM t WHERE a = fo$a$p AS x U&'a;b);(' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = "a;b/*;a;b" SELECT $$$x$;a;ba;b$$ FROM t WHERE a = "END ;" AS x;
| FROM t WHERE a = CREATE FUNCTION zbum() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);(;;'; SELECT 1; END FROM t WHERE a = U&'select 1;' SELECT wc$e$w +;
| col $rh$; $rh$ ,; -- trailing note;
| COPY i__ FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "a;b/*;a;b" SELECT $$$x$;a;ba;b$$ FROM t WHERE a = "END ;" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "a;b/*;a;b" SELECT $$$x$;a;ba;b$$ FROM t WHERE a = "END ;" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT CREATE FUNCTION zqot() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; /*;'; SELECT 1; END SELECT;
| 1 CREATE FUNCTION bddl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;b/*;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION zqot() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; /*;'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION zqot() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; /*;'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| COPY kfk FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kfk FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kfk FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col ';end;' FROM t WHERE a =;
| , "select 1;select 1;/*;" ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ';end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ';end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT "/*;" SELECT /*select 1;*/ ,;
< error
| 42601 syntax error at or near ";"
| position: 36
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "/*;" SELECT /*select 1;*/ ,;
< error
| 42601 syntax error at or near ";"
| position: 36
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "/*;" SELECT /*select 1;*/ ,;
< error
| 42601 syntax error at or near ";"
| position: 36
> exec
| FROM t WHERE a = E'\n;\n;'';' SELECT CREATE FUNCTION werb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END AS x;
| COPY cc_ FROM stdin;
| 1	\N
| a;b
| \.
|
| SELECT cp$o$u FROM t WHERE a = --end;;end;
|  FROM t WHERE a = "--;/*;" ,;COPY zqf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E'\n;\n;'';' SELECT CREATE FUNCTION werb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\n;\n;'';' SELECT CREATE FUNCTION werb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY lqd FROM stdin;
| a;b
| 1	\N
| \.
|
| 1 '$;$' , km$w$s FROM t WHERE a = (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lqd FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lqd FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , $$;$$ 1 E'\\\\' col;
| , CREATE FUNCTION gmqm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END 1;
| AS x $$a;bselect 1;$$ +; -- trailing note;
| COPY zcb FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $$;$$ 1 E'\\\\' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $$;$$ 1 E'\\\\' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY bsh FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bsh FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bsh FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , ");(select 1;$x$;" col CREATE FUNCTION xsui() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;);(--;'; SELECT 1; END + "/*;;;--;" AS x;
| 1 ee$i$i + wa$a$t , --a;bEND ;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ");(select 1;$x$;" col CREATE FUNCTION xsui() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;);(--;'; SELECT 1; END + "/*;;;--;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ");(select 1;$x$;" col CREATE FUNCTION xsui() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;);(--;'; SELECT 1; END + "/*;;;--;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 xz$i$j col b'01;10' FROM t WHERE a =;SELECT $$;;;$;$$$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 xz$i$j col b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 xz$i$j col b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + --;END ;
|  AS x $$--;select 1;$$ SELECT U&';;' col; -- trailing note;
| COPY _ft FROM stdin;
| c;d;e
| c;d;e
| \.
|
| SELECT /*/*/*);(select 1;*/*/*/ 1 e''';\\' 1; COPY dte FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --;END ;
|  AS x $$--;select 1;$$ SELECT U&';;' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --;END ;
|  AS x $$--;select 1;$$ SELECT U&';;' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT U&'--;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 35
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'--;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 35
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'--;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 35
> exec
| , ----;
|  + gd$h$k FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ----;
|  + gd$h$k FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ----;
|  + gd$h$k FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , /*;;end;*/ AS x; -- trailing note;
| COPY pzp FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*;;end;*/ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*;;end;*/ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 /*/*/*;--;a;b*/*/*/ col; -- trailing note;
| col X'01;10' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*;--;a;b*/*/*/ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*;--;a;b*/*/*/ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , X'01;10' + uj$c$g 1 (SELECT 1; SELECT 2) +;
| + e'\\' SELECT (SELECT 1; SELECT 2) col;col U&';' SELECT X'01;10' col 'END ;$x$;--;\' 1;
| AS x $b_ie$END ;$b_ie$ 1 --);(
|  FROM t WHERE a = ");(" FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , X'01;10' + uj$c$g 1 (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , X'01;10' + uj$c$g 1 (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x CREATE FUNCTION tymz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;/*;);('; SELECT 1; END col;
| FROM t WHERE a = ');(' , --;;; 
|  FROM t WHERE a =;
| SELECT 'select 1;end;a;b' SELECT X'01;10' SELECT mt$m$p FROM t WHERE a =; COPY qhr FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION tymz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;/*;);('; SELECT 1; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION tymz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;/*;);('; SELECT 1; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY sfr FROM stdin;
| 1	\N
| 1	\N
| \.
|
| , --; 
|  1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sfr FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sfr FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY kbf FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kbf FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kbf FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 /*/*/*;;end;;;*/*/*/ col _z$_$e AS x;
| SELECT --a;b
|  SELECT E'\';\n;' SELECT b'01;10' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*;;end;;;*/*/*/ col _z$_$e AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*;;end;;;*/*/*/ col _z$_$e AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col '$x$;' FROM t WHERE a =; -- trailing note;
| + $$select 1;$$ + (SELECT 1; SELECT 2) 1; SELECT $$--;;$;$$$ FROM t WHERE a = 'select 1;' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '$x$;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '$x$;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + E'\\' FROM t WHERE a =; FROM t WHERE a = E''';' + --;;
|  ,; -- trailing note;
| + "a;b" FROM t WHERE a = hm$v$q FROM t WHERE a = --select 1;
|  AS x;
| 1 zy$b$i , /*/*/*$x$;*/*/*/ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x /*/*/*);(*/*/*/ FROM t WHERE a = CREATE FUNCTION vbiz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*);(*/*/*/ FROM t WHERE a = CREATE FUNCTION vbiz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*);(*/*/*/ FROM t WHERE a = CREATE FUNCTION vbiz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY eqh FROM stdin;
| c;d;e
| 1	\N
| \.
|
| SELECT x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY eqh FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY eqh FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , $ei$$;$$ei$ 1;
| , "/*;);(" AS x $$;;;;end;$$ col; -- trailing note;
| 1 $$$x$;end;$$ ,;
| col (SELECT 1; SELECT 2) col "$x$;" 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $ei$$;$$ei$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $ei$$;$$ei$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , U&'select 1;' 1 ag$i$n col 'END ;$x$;select 1;' AS x;
| AS x ";$x$;; " +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'select 1;' 1 ag$i$n col 'END ;$x$;select 1;' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'select 1;' 1 ag$i$n col 'END ;$x$;select 1;' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY jzt FROM stdin;
| a;b
| 1	\N
| \.
|
| AS x U&'a;ba;b' col --/*;--;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jzt FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jzt FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col $$select 1;$$ 1; -- trailing note;
| col mj$f$q AS x $kqzk$select 1;;;$kqzk$ 1;
| , --end;);(
|  col e''';\n;'';' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$select 1;$$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$select 1;$$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , $jcn$);(/*;$jcn$ 1;
| 1 U&'select 1;/*;' AS x;1 (SELECT 1; SELECT 2) 1;+ '$x$;--;; ' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $jcn$);(/*;$jcn$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $jcn$);(/*;$jcn$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT U&'/*;' ,; -- trailing note;
| SELECT ';;; \' AS x /*/*a;b; select 1;*/*/ col hz$s$m col;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'/*;' ,;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'/*;' ,;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| col U&'; a;b' 1 U&');(' FROM t WHERE a = --$;$
|  ,;
| + $pd$);();($x$;$pd$ , '; ' 1 B'01;10' AS x;
| SELECT "/*;--;" AS x U&'END ;' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'; a;b' 1 U&');(' FROM t WHERE a = --$;$
|  ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'; a;b' 1 U&');(' FROM t WHERE a = --$;$
|  ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY ghd FROM stdin;
| 1	\N
| 1	\N
| \.
|
| + x'01;10' SELECT --;;select 1;);(
|  col e'\';\n;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ghd FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ghd FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , /*--;*/ col --;;a;b
|  1 ';;\' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*--;*/ col --;;a;b
|  1 ';;\' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*--;*/ col --;;a;b
|  1 ';;\' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT e'\\\';' col ua$v$j +;
| 1 e'\n;'';' col;
< error
| 42601 syntax error at or near "ua$v$j"
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\\\';' col ua$v$j +;
< error
| 42601 syntax error at or near "ua$v$j"
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT e'\\\';' col ua$v$j +;
< error
| 42601 syntax error at or near "ua$v$j"
| position: 21
> exec
| + ';;);(\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ';;);(\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ';;);(\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , --end;end;);(
|  SELECT; COPY wyk FROM stdin;
| c;d;e
| c;d;e
| \.
|
| col e'\n;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --end;end;);(
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --end;end;);(
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col CREATE FUNCTION w_ue() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';a;b; '; END , (SELECT 1; SELECT 2) , CREATE FUNCTION gaj_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;; '; SELECT 1; END FROM t WHERE a =;
| SELECT CREATE FUNCTION _gfg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; a;bend;'; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION w_ue() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';a;b; '; END , (SELECT 1; SELECT 2) , CREATE FUNCTION gaj_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;; '; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION w_ue() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';a;b; '; END , (SELECT 1; SELECT 2) , CREATE FUNCTION gaj_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;; '; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = zm$j$h 1 CREATE FUNCTION jami() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;'; SELECT 1; END , (SELECT 1; SELECT 2) FROM t WHERE a =; -- trailing note;
| 1 ");(" FROM t WHERE a = /*/*/*;);(*/*/*/ + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = zm$j$h 1 CREATE FUNCTION jami() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;'; SELECT 1; END , (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = zm$j$h 1 CREATE FUNCTION jami() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;'; SELECT 1; END , (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + CREATE FUNCTION dlm_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; /*;'; SELECT 1; END 1 U&');(a;b' FROM t WHERE a = $_d$$x$;END ;$_d$ SELECT; -- trailing note;
| COPY mqd FROM stdin;
| c;d;e
| 1	\N
| \.
|  + X'01;10' +; -- trailing note;
| COPY lfe FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION dlm_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; /*;'; SELECT 1; END 1 U&');(a;b' FROM t WHERE a = $_d$$x$;END ;$_d$ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION dlm_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; /*;'; SELECT 1; END 1 U&');(a;b' FROM t WHERE a = $_d$$x$;END ;$_d$ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , U&'$x$;; ' AS x;
| COPY doj FROM stdin;
| a;b
| 1	\N
| \.
|
| col /*/*/*END ;*/*/*/ , b'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'$x$;; ' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'$x$;; ' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x E''';\';' 1 (SELECT 1; SELECT 2) + $$end;;;/*;$$ AS x;, $$select 1;end;$$ AS x;
| FROM t WHERE a = CREATE FUNCTION gjzs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;--;'; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x E''';\';' 1 (SELECT 1; SELECT 2) + $$end;;;/*;$$ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x E''';\';' 1 (SELECT 1; SELECT 2) + $$end;;;/*;$$ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT $$; ; $$ col (SELECT 1; SELECT 2) col /*/*; */*/ AS x;
< error
| 42601 syntax error at or near "("
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$; ; $$ col (SELECT 1; SELECT 2) col /*/*; */*/ AS x;
< error
| 42601 syntax error at or near "("
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$; ; $$ col (SELECT 1; SELECT 2) col /*/*; */*/ AS x;
< error
| 42601 syntax error at or near "("
| position: 21
> exec
| COPY rld FROM stdin;
| 1	\N
| c;d;e
| \.
|
| , lw$_$v FROM t WHERE a = x'01;10' +;
| + (SELECT 1; SELECT 2) FROM t WHERE a = /*/*/ *;$x$;*/*/ ,;
| FROM t WHERE a = '$x$;end;$;$' col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY rld FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY rld FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , E'\n;\\' AS x /*/*/*);(*/*/*/ + --;a;b);(
|  +; COPY rfv FROM stdin;
| a;b
| a;b
| \.
|
| COPY lji FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , E'\n;\\' AS x /*/*/*);(*/*/*/ + --;a;b);(
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , E'\n;\\' AS x /*/*/*);(*/*/*/ + --;a;b);(
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY yh_ FROM stdin;
| a;b
| 1	\N
| \.
|
| , CREATE FUNCTION hofg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END , CREATE FUNCTION ndas() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; END SELECT CREATE FUNCTION frkg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;;;END ;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yh_ FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yh_ FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col /*/*/*END ;; */*/*/ 1 --;);(end;
|  AS x /*/*select 1;; */*/ 1;, pv$x$a ,;
| col e'\';' 1 --a;b; );(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*/*END ;; */*/*/ 1 --;);(end;
|  AS x /*/*select 1;; */*/ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*/*END ;; */*/*/ 1 --;);(end;
|  AS x /*/*select 1;; */*/ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = --a;b/*;; 
|  , X'01;10' col --$x$;$x$;
|  ,; 1 li$m$o AS x (SELECT 1; SELECT 2) FROM t WHERE a = "--;" ,;
| 1 "select 1;" ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --a;b/*;; 
|  , X'01;10' col --$x$;$x$;
|  ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --a;b/*;; 
|  , X'01;10' col --$x$;$x$;
|  ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 CREATE FUNCTION kbgu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';; --;'; END ,; -- trailing note;
| AS x ';$;$' + $$);($x$;$$ AS x ";end;a;b" ,;
| , $pmq$);($pmx$$pmq$ AS x e'\\'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION kbgu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';; --;'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION kbgu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';; --;'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = E''';' , (SELECT 1; SELECT 2) 1;
| + $$;;a;b);($$ + U&'select 1;end;end;' AS x '/*;' SELECT;
| COPY tws FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E''';' , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E''';' , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + CREATE FUNCTION d_bo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;end;'; END AS x /*/*--;end;*/*/ SELECT "END ;);(; " AS x; -- trailing note;
| + "$;$; /*;" SELECT "--;" FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION d_bo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;end;'; END AS x /*/*--;end;*/*/ SELECT "END ;);(; " AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION d_bo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;end;'; END AS x /*/*--;end;*/*/ SELECT "END ;);(; " AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x /*);(*/ FROM t WHERE a = $opgw$end;$opgx$$opgw$ FROM t WHERE a = U&';END ;' 1;
| COPY o_e FROM stdin;
| c;d;e
| 1	\N
| \.
| 1 'select 1;--;$x$;' SELECT $zwdc$END ;--;--;$zwdc$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*);(*/ FROM t WHERE a = $opgw$end;$opgx$$opgw$ FROM t WHERE a = U&';END ;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*);(*/ FROM t WHERE a = $opgw$end;$opgx$$opgw$ FROM t WHERE a = U&';END ;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) SELECT CREATE FUNCTION lera() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END + --);($x$;
|  FROM t WHERE a =;
| SELECT X'01;10' 1 "; " 1 /*/*/*/ *;$x$;*/*/*/ col;
| , pk$i$_ 1;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) SELECT CREATE FUNCTION lera() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END + --);($x$;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) SELECT CREATE FUNCTION lera() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END + --);($x$;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| COPY aaq FROM stdin;
| a;b
| 1	\N
| \.
|
| AS x (SELECT 1; SELECT 2) , 'end;\' AS x; col U&'$;$/*;;' SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY aaq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY aaq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 U&';; $;$' SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&';; $;$' SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&';; $;$' SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , yq$v$p +;
| COPY obv FROM stdin;
| a;b
| c;d;e
| \.
|
| , X'01;10' 1 b'01;10' AS x "$;$$x$;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , yq$v$p +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , yq$v$p +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col U&'end;/*;;;' col; + ux$p$c + '$x$;' SELECT X'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'end;/*;;;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'end;/*;;;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT CREATE FUNCTION dmdw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;'; END 1 CREATE FUNCTION cbvq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;a;b'; SELECT 1; END 1 U&');(' col;
| COPY kpf FROM stdin;
| a;b
| a;b
| \.
|
| + wy$j$c ,;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION dmdw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;'; END 1 CREATE FUNCTION cbvq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;a;b'; SELECT 1; END 1 U&');(' col;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION dmdw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;'; END 1 CREATE FUNCTION cbvq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;a;b'; SELECT 1; END 1 U&');(' col;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| FROM t WHERE a = ';$;$--;\' col;
| , x'01;10' SELECT U&';' col;
| 1 x'01;10' AS x "select 1;" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ';$;$--;\' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ';$;$--;\' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col ww$q$k AS x; FROM t WHERE a = 'a;bselect 1;; \' FROM t WHERE a = md$z$j +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ww$q$k AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ww$q$k AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY emp FROM stdin;
| a;b
| a;b
| \.
|
| FROM t WHERE a = /*END ;;;/ *;*/ + (SELECT 1; SELECT 2) +;1 e''';\n;\\' SELECT ----;;
|  col; col CREATE FUNCTION zhzp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;--;'; END , CREATE FUNCTION vpuc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;select 1;'; END AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY emp FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY emp FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = ');(\' FROM t WHERE a = x'01;10' SELECT $$;;);();($$ +; -- trailing note;
| COPY _zo FROM stdin;
| c;d;e
| c;d;e
| \.
|  FROM t WHERE a = --$x$;
|  , X'01;10' , '; );(\' ,;
| 1 CREATE FUNCTION kqhe() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;a;b'; SELECT 1; END SELECT 'select 1;a;b;' col e'\n;\\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ');(\' FROM t WHERE a = x'01;10' SELECT $$;;);();($$ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ');(\' FROM t WHERE a = x'01;10' SELECT $$;;);();($$ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col '$;$/*;;' AS x CREATE FUNCTION wmty() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $;$END ;'; SELECT 1; END col /*/*/*a;bselect 1;--;*/*/*/ ,; 1 /*/*/*END ;*/*/*/ FROM t WHERE a = e'\\\n;' AS x;
| COPY pfz FROM stdin;
| c;d;e
| a;b
| \.
|  FROM t WHERE a = ----;a;b
|  AS x U&'--;end;end;' AS x /*/*$;$*/*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '$;$/*;;' AS x CREATE FUNCTION wmty() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $;$END ;'; SELECT 1; END col /*/*/*a;bselect 1;--;*/*/*/ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '$;$/*;;' AS x CREATE FUNCTION wmty() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $;$END ;'; SELECT 1; END col /*/*/*a;bselect 1;--;*/*/*/ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = $cbi$select 1;);(/*;$cbx$$cbi$ AS x; -- trailing note;
| SELECT U&'$;$$x$;END ;' , U&');(; ' +;
| COPY ize FROM stdin;
| 1	\N
| a;b
| \.
|  -- trailing note;
| COPY br_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = $cbi$select 1;);(/*;$cbx$$cbi$ AS x;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a = $cbi$select 1;);(/*;$cbx$$cbi$ AS x;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| SELECT b'01;10' FROM t WHERE a = U&'$x$;a;b' 1;
| COPY pop FROM stdin;
| c;d;e
| 1	\N
| \.
| COPY blv FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY l_l FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 46
> exec
| BEGIN
< ok 0 -
> exec
| SELECT b'01;10' FROM t WHERE a = U&'$x$;a;b' 1;
< error
| 42601 syntax error at or near "1"
| position: 46
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT b'01;10' FROM t WHERE a = U&'$x$;a;b' 1;
< error
| 42601 syntax error at or near "1"
| position: 46
> exec
| FROM t WHERE a = $l_f$$x$;select 1;END ;$l_f$ ,;
| col E'\n;\\' SELECT;
| , (SELECT 1; SELECT 2) FROM t WHERE a = U&'END ;/*;;;' col --end;$;$$;$
|  ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $l_f$$x$;select 1;END ;$l_f$ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $l_f$$x$;select 1;END ;$l_f$ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , (SELECT 1; SELECT 2) SELECT E''';\\' 1 (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) SELECT E''';\\' 1 (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) SELECT E''';\\' 1 (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY juq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY juq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY juq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 'end;$;$' 1 e''';\n;\\' SELECT --select 1;a;b$x$;
|  AS x;
| + CREATE FUNCTION vprs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;;'; END col nn$z$y col /*;;;*/ SELECT;
| , "; end;" AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 'end;$;$' 1 e''';\n;\\' SELECT --select 1;a;b$x$;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 'end;$;$' 1 e''';\n;\\' SELECT --select 1;a;b$x$;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , (SELECT 1; SELECT 2) 1 ";;END ;" ,;
| 1 "$x$;--;; " col; -- trailing note;
| + /*/*/*END ;; / *;*/*/*/ SELECT;
| FROM t WHERE a = /*/*/*;*/*/*/ FROM t WHERE a = sx$j$e + /*/*/*/ *;;$;$*/*/*/ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) 1 ";;END ;" ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) 1 ";;END ;" ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col CREATE FUNCTION txxw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;'; END 1;col --);();(--;
|  1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION txxw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;'; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION txxw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;'; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = U&'--;/*;' SELECT x'01;10' SELECT E''';\n;\';' FROM t WHERE a =; -- trailing note;
| 1 $y$select 1;end;; $y$ + rc$g$e + (SELECT 1; SELECT 2) FROM t WHERE a =; 1 px$m$c FROM t WHERE a = --end;/*;
|  col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'--;/*;' SELECT x'01;10' SELECT E''';\n;\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'--;/*;' SELECT x'01;10' SELECT E''';\n;\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x 'a;bEND ;;;' +;
| + (SELECT 1; SELECT 2) SELECT CREATE FUNCTION pova() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;end;'; SELECT 1; END col; -- trailing note;
| SELECT 'end;; $x$;' col CREATE FUNCTION kdyc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END ,;
| col E'\\' AS x gq$q$n SELECT --END ;;;a;b
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x 'a;bEND ;;;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x 'a;bEND ;;;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY osk FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| SELECT --;
|  SELECT E'\n;\\'';' 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY osk FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY osk FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY lge FROM stdin;
| a;b
| a;b
| \.
|
| COPY pcp FROM stdin;
| c;d;e
| 1	\N
| \.
|
| col U&'a;bEND ;END ;' , /*/ *;end;END ;*/ AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lge FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lge FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = e'\';'';' AS x (SELECT 1; SELECT 2) AS x 'end;);(--;' AS x;
| AS x --END ;$;$
|  col;
| + '$x$;a;b--;' col X'01;10' AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = e'\';'';' AS x (SELECT 1; SELECT 2) AS x 'end;);(--;' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = e'\';'';' AS x (SELECT 1; SELECT 2) AS x 'end;);(--;' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT ');($x$;' SELECT CREATE FUNCTION aido() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END 1; + U&'$x$;END ;' col;
| COPY taq FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "CREATE"
| position: 25
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ');($x$;' SELECT CREATE FUNCTION aido() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END 1;
< error
| 42601 syntax error at or near "CREATE"
| position: 25
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT ');($x$;' SELECT CREATE FUNCTION aido() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END 1;
< error
| 42601 syntax error at or near "CREATE"
| position: 25
> exec
| + /*/*a;bend;*/*/ FROM t WHERE a = U&';;;;' +;
| AS x 'select 1;$x$;' col zc$l$g + E'\n;\\' 1;
| , ";" , $d$;);();($d$ , ----;
|  +;
| col '/*;);(;' + e'\\\n;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*a;bend;*/*/ FROM t WHERE a = U&';;;;' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*a;bend;*/*/ FROM t WHERE a = U&';;;;' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col /*/ *;*/ AS x CREATE FUNCTION zscr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END ,;
| COPY rqp FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/ *;*/ AS x CREATE FUNCTION zscr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/ *;*/ AS x CREATE FUNCTION zscr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = e''';' 1 E'\\\n;\';' SELECT;
| 1 "$;$" FROM t WHERE a = dz$a$_ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = e''';' 1 E'\\\n;\';' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = e''';' 1 E'\\\n;\';' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY jor FROM stdin;
| 1	\N
| a;b
| \.
|
| FROM t WHERE a = --$;$;;);(
|  col; SELECT e'\\' ,;
| COPY opq FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jor FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jor FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col (SELECT 1; SELECT 2) col "--;; /*;" FROM t WHERE a = E''';'';\\' col;
| 1 /*/*END ;a;b*/*/ SELECT; -- trailing note;
| col /*/*/*select 1;select 1;*/*/*/ AS x /*/*/*select 1;;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) col "--;; /*;" FROM t WHERE a = E''';'';\\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) col "--;; /*;" FROM t WHERE a = E''';'';\\' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x ");(" , if$n$j SELECT;
| COPY ljg FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ");(" , if$n$j SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ");(" , if$n$j SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + U&'; $;$' FROM t WHERE a =; , /*;;$;$select 1;*/ AS x CREATE FUNCTION wfth() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;end;'; SELECT 1; END AS x;
| COPY oil FROM stdin;
| 1	\N
| a;b
| \.
|
| COPY hjz FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'; $;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'; $;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , U&');(; ;' 1; + /*;*/ 1 --$;$);(
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&');(; ;' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&');(; ;' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col E'\\'';\n;' SELECT $$END ;$$ 1 U&'select 1;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col E'\\'';\n;' SELECT $$END ;$$ 1 U&'select 1;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col E'\\'';\n;' SELECT $$END ;$$ 1 U&'select 1;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x --/*;; ; 
|  AS x $$;;$$ AS x E''';\\' 1; -- trailing note;
| , U&'a;b;a;b' AS x lq$w$s , uh$p$r AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --/*;; ; 
|  AS x $$;;$$ AS x E''';\\' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --/*;; ; 
|  AS x $$;;$$ AS x E''';\\' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col e'\\'';' 1 --$x$;END ;
|  FROM t WHERE a = /*/ *;*/ col;
| 1 B'01;10' 1; , '$x$;END ;--;\' FROM t WHERE a = x'01;10' 1 X'01;10' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e'\\'';' 1 --$x$;END ;
|  FROM t WHERE a = /*/ *;*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e'\\'';' 1 --$x$;END ;
|  FROM t WHERE a = /*/ *;*/ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col /*/*end;select 1;*/*/ col U&'a;bEND ;' FROM t WHERE a = $$; $$ col; -- trailing note;
| FROM t WHERE a = /*/*;;*/*/ col x'01;10' SELECT "select 1;" 1; -- trailing note;
| + CREATE FUNCTION zszn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;$x$;'; END SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*end;select 1;*/*/ col U&'a;bEND ;' FROM t WHERE a = $$; $$ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*end;select 1;*/*/ col U&'a;bEND ;' FROM t WHERE a = $$; $$ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 '$;$;a;b\' +;
| , hb$t$q 1; FROM t WHERE a = E''';\\' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '$;$;a;b\' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '$;$;a;b\' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col ';$x$;$x$;' FROM t WHERE a = nu$k$u SELECT $$end;--;$$ FROM t WHERE a =;
| AS x 'a;b\' , e''';\';' col B'01;10' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ';$x$;$x$;' FROM t WHERE a = nu$k$u SELECT $$end;--;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ';$x$;$x$;' FROM t WHERE a = nu$k$u SELECT $$end;--;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT E'\n;' AS x /*/*;;;*/*/ 1; COPY bhf FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 32
> exec
| BEGIN
< ok 0 -
> exec
| SELECT E'\n;' AS x /*/*;;;*/*/ 1;
< error
| 42601 syntax error at or near "1"
| position: 32
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT E'\n;' AS x /*/*;;;*/*/ 1;
< error
| 42601 syntax error at or near "1"
| position: 32
> exec
| FROM t WHERE a = '$;$end;);(' SELECT "END ;" +;
| , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = '$;$end;);(' SELECT "END ;" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = '$;$end;);(' SELECT "END ;" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , X'01;10' col $$select 1;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , X'01;10' col $$select 1;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , X'01;10' col $$select 1;$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col (SELECT 1; SELECT 2) FROM t WHERE a = "$;$; " FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) FROM t WHERE a = "$;$; " FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) FROM t WHERE a = "$;$; " FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY meq FROM stdin;
| a;b
| a;b
| \.
|
| col (SELECT 1; SELECT 2) col;
| COPY azc FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| SELECT '/*;);(END ;' + '--;a;bEND ;' 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY meq FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY meq FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = "a;b$x$;" , --; 
|  AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "a;b$x$;" , --; 
|  AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "a;b$x$;" , --; 
|  AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY qyg FROM stdin;
| a;b
| c;d;e
| \.
|
| 1 '$x$;select 1;;;\' col $g$END ;$g$ SELECT;
| COPY zvn FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qyg FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qyg FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , '; ' col qm$o$h ,;
| COPY wvu FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , '; ' col qm$o$h ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , '; ' col qm$o$h ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , "$x$;END ;" FROM t WHERE a = ";a;b--;" , $uwno$$x$;$uwnx$$uwno$ 1; -- trailing note;
| AS x $xo_$);(end;end;$xo_$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "$x$;END ;" FROM t WHERE a = ";a;b--;" , $uwno$$x$;$uwnx$$uwno$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "$x$;END ;" FROM t WHERE a = ";a;b--;" , $uwno$$x$;$uwnx$$uwno$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT E'\';\';' 1 CREATE FUNCTION osba() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$select 1;'; END 1; -- trailing note;
| SELECT --;;
|  FROM t WHERE a = CREATE FUNCTION ffcy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$;$;;'; SELECT 1; END 1;
| 1 U&';;' 1 'select 1;--;\' SELECT ");();($x$;" col;
| 1 e'\';\';' SELECT "select 1;" col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 18
> exec
| BEGIN
< ok 0 -
> exec
| SELECT E'\';\';' 1 CREATE FUNCTION osba() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$select 1;'; END 1;
< error
| 42601 syntax error at or near "1"
| position: 18
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT E'\';\';' 1 CREATE FUNCTION osba() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$select 1;'; END 1;
< error
| 42601 syntax error at or near "1"
| position: 18
> exec
| COPY syx FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY syx FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY syx FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , U&'/*;' FROM t WHERE a = B'01;10' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'/*;' FROM t WHERE a = B'01;10' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'/*;' FROM t WHERE a = B'01;10' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x U&';;;; ' , B'01;10' +; -- trailing note;
| , n_$f$e AS x /*/*/*end;);($;$*/*/*/ + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&';;;; ' , B'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&';;;; ' , B'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x B'01;10' + b'01;10' FROM t WHERE a = "$;$" col; -- trailing note;
| , '$x$;select 1;end;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x B'01;10' + b'01;10' FROM t WHERE a = "$;$" col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x B'01;10' + b'01;10' FROM t WHERE a = "$;$" col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 uf$s$p col; -- trailing note;
| , dr$r$n col '$x$;--;\' , ";" ,;1 E''';'';' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 uf$s$p col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 uf$s$p col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = jt$o$b + U&'end;;/*;' col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = jt$o$b + U&'end;;/*;' col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = jt$o$b + U&'end;;/*;' col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , vd$v$l AS x; SELECT U&'a;b' + E'\';' FROM t WHERE a = $$end;;;$$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , vd$v$l AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , vd$v$l AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY gew FROM stdin;
| a;b
| c;d;e
| \.
|  COPY ohh FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gew FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gew FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT se$y$b FROM t WHERE a =; col U&'$;$end;' col "; ; a;b" 1 fq$h$z +;COPY smp FROM stdin;
| c;d;e
| c;d;e
| \.
|
| , --;;; 
|  SELECT b_$v$d SELECT;
< error
| 42601 syntax error at or near ";"
| position: 31
> exec
| BEGIN
< ok 0 -
> exec
| SELECT se$y$b FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 31
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT se$y$b FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 31
> exec
| SELECT '--;end;' 1 --a;b$;$select 1;
|  1 "; $;$$;$" AS x; AS x b'01;10' 1 --);(select 1;
|  1;
| col /*);(select 1;--;*/ 1 CREATE FUNCTION yxxk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;--;'; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 18
> exec
| BEGIN
< ok 0 -
> exec
| SELECT '--;end;' 1 --a;b$;$select 1;
|  1 "; $;$$;$" AS x;
< error
| 42601 syntax error at or near "1"
| position: 18
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT '--;end;' 1 --a;b$;$select 1;
|  1 "; $;$$;$" AS x;
< error
| 42601 syntax error at or near "1"
| position: 18
> exec
| FROM t WHERE a = "END ;" col e'\';' +;AS x ");(;;" SELECT B'01;10' FROM t WHERE a =;
| SELECT /*/*--;$x$;*/*/ FROM t WHERE a = e'\';\n;\\' 1;
| COPY ohg FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "END ;" col e'\';' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "END ;" col e'\';' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 e''';\';' AS x CREATE FUNCTION vjfv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; );(select 1;'; SELECT 1; END 1 X'01;10' SELECT; AS x _l$a$i + th$i$k , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e''';\';' AS x CREATE FUNCTION vjfv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; );(select 1;'; SELECT 1; END 1 X'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e''';\';' AS x CREATE FUNCTION vjfv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; );(select 1;'; SELECT 1; END 1 X'01;10' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + (SELECT 1; SELECT 2) SELECT CREATE FUNCTION kgoz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;;;;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) SELECT CREATE FUNCTION kgoz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;;;;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) SELECT CREATE FUNCTION kgoz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;;;;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY gqq FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gqq FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gqq FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT x'01;10' FROM t WHERE a = --END ;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 43
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' FROM t WHERE a = --END ;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 43
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT x'01;10' FROM t WHERE a = --END ;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 43
> exec
| SELECT x'01;10' , 'end;$x$;select 1;' SELECT;
| COPY yhp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 68
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' , 'end;$x$;select 1;' SELECT;
< error
| 22P02 ";" is not a valid hexadecimal digit
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| COPY yhp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 23
> exec
| COPY sef FROM stdin;
| a;b
| 1	\N
| \.
|  COPY wjd FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sef FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sef FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 $zxk$; /*;end;$zxk$ AS x "select 1;" , --END ;);(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $zxk$; /*;end;$zxk$ AS x "select 1;" , --END ;);(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $zxk$; /*;end;$zxk$ AS x "select 1;" , --END ;);(
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , ';;$x$;' FROM t WHERE a = CREATE FUNCTION smmj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$end;'; END 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ';;$x$;' FROM t WHERE a = CREATE FUNCTION smmj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$end;'; END 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ';;$x$;' FROM t WHERE a = CREATE FUNCTION smmj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$end;'; END 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY jes FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jes FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jes FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + E'\';\';\';' 1 e'\';' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\';\';\';' 1 e'\';' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\';\';\';' 1 e'\';' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT e''';\n;\n;' col /*/*a;b*/*/ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 37
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e''';\n;\n;' col /*/*a;b*/*/ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 37
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT e''';\n;\n;' col /*/*a;b*/*/ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 37
> exec
| FROM t WHERE a = U&'END ;END ;' FROM t WHERE a = '$;$$;$/*;' 1 U&'/*;' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'END ;END ;' FROM t WHERE a = '$;$$;$/*;' 1 U&'/*;' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'END ;END ;' FROM t WHERE a = '$;$$;$/*;' 1 U&'/*;' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT ym$i$y FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 31
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ym$i$y FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 31
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT ym$i$y FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 31
> exec
| col --select 1;/*;select 1;
|  AS x X'01;10' , $$;END ;$$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --select 1;/*;select 1;
|  AS x X'01;10' , $$;END ;$$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --select 1;/*;select 1;
|  AS x X'01;10' , $$;END ;$$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x 'end;' col e'\';' SELECT; -- trailing note;
| COPY oio FROM stdin;
| c;d;e
| c;d;e
| \.
|
| , b'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x 'end;' col e'\';' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x 'end;' col e'\';' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY zn_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
| SELECT U&'/*;a;b--;' 1 "; $;$select 1;" AS x "$;$select 1;" +;
| 1 b'01;10' SELECT --a;b$;$
|  AS x 'END ;' ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zn_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zn_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 '$;$' ,;
| col X'01;10' FROM t WHERE a = "end;" , /*/*/*--;*/*/*/ +;
| FROM t WHERE a = x'01;10' AS x /*/*; */*/ FROM t WHERE a =;
| COPY niz FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '$;$' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '$;$' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + 'end;' FROM t WHERE a =;
| COPY hbn FROM stdin;
| c;d;e
| 1	\N
| \.
|  col --END ;
|  SELECT X'01;10' +; -- trailing note;
| AS x CREATE FUNCTION xlfd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; ; $;$'; END , e''';\\\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + 'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + 'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + $chwy$--;;;$;$$chwx$$chwy$ + '; end;' + $$;;END ;$$ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $chwy$--;;;$;$$chwx$$chwy$ + '; end;' + $$;;END ;$$ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $chwy$--;;;$;$$chwx$$chwy$ + '; end;' + $$;;END ;$$ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = ');(END ;END ;' 1;col ');();(' , --select 1;end;end;
|  FROM t WHERE a =; -- trailing note;
| COPY gdj FROM stdin;
| 1	\N
| 1	\N
| \.
|
| SELECT B'01;10' AS x E'\\'';\n;' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ');(END ;END ;' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ');(END ;END ;' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY iig FROM stdin;
| a;b
| c;d;e
| \.
|  -- trailing note;
| col --select 1;
|  FROM t WHERE a = CREATE FUNCTION cnnc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$x$;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY iig FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY iig FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 'end;$x$;' SELECT "--;" col; , rp$k$x col ';' SELECT b'01;10' AS x; -- trailing note;
| COPY kai FROM stdin;
| a;b
| 1	\N
| \.
|
| FROM t WHERE a = (SELECT 1; SELECT 2) , $$--;$$ AS x 'select 1;a;b' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 'end;$x$;' SELECT "--;" col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 'end;$x$;' SELECT "--;" col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT --END ;a;b$;$
|  FROM t WHERE a = --end;$x$;;;
|  +;
| FROM t WHERE a = X'01;10' col ");(; " ,;
< error
| 42601 syntax error at or near ";"
| position: 55
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --END ;a;b$;$
|  FROM t WHERE a = --end;$x$;;;
|  +;
< error
| 42601 syntax error at or near ";"
| position: 55
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --END ;a;b$;$
|  FROM t WHERE a = --end;$x$;;;
|  +;
< error
| 42601 syntax error at or near ";"
| position: 55
> exec
| + "; /*;" FROM t WHERE a = $oul$$;$);($x$;$oul$ + E'\n;\\\';' FROM t WHERE a =;COPY mma FROM stdin;
| a;b
| 1	\N
| \.
| 1 --select 1;a;b
|  , b'01;10' SELECT;
| 1 '$x$;\' 1 B'01;10' 1 E'\';'';\\' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "; /*;" FROM t WHERE a = $oul$$;$);($x$;$oul$ + E'\n;\\\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "; /*;" FROM t WHERE a = $oul$$;$);($x$;$oul$ + E'\n;\\\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + --);(
|  FROM t WHERE a = b'01;10' 1 ----;
|  FROM t WHERE a =; 1 --;; 
|  SELECT;
| COPY c_d FROM stdin;
| 1	\N
| c;d;e
| \.
|
| AS x U&'$;$);(;' , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --);(
|  FROM t WHERE a = b'01;10' 1 ----;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --);(
|  FROM t WHERE a = b'01;10' 1 ----;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col /*END ;;; */ , U&'$;$$;$a;b' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*END ;;; */ , U&'$;$$;$a;b' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*END ;;; */ , U&'$;$$;$a;b' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + E'\\\\\';' col;
| , B'01;10' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\\\\\';' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\\\\\';' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + ';;' AS x;
| COPY cyn FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ';;' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ';;' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY gdo FROM stdin;
| c;d;e
| c;d;e
| \.
|
| col x'01;10' SELECT CREATE FUNCTION hlrz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';);(END ;'; END ,;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gdo FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gdo FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + e'\n;'';\\' FROM t WHERE a = /*/*--;);(*/*/ SELECT "/*;end;" ,;
| FROM t WHERE a = e'\\\n;' col;
| AS x h_$e$o 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + e'\n;'';\\' FROM t WHERE a = /*/*--;);(*/*/ SELECT "/*;end;" ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + e'\n;'';\\' FROM t WHERE a = /*/*--;);(*/*/ SELECT "/*;end;" ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , lr$c$u col X'01;10' col /*/ *;$;$);(*/ AS x;
| 1 U&');($x$;' AS x CREATE FUNCTION jjhr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$;$'; SELECT 1; END col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , lr$c$u col X'01;10' col /*/ *;$;$);(*/ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , lr$c$u col X'01;10' col /*/ *;$;$);(*/ AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY qpw FROM stdin;
| 1	\N
| c;d;e
| \.
|
| AS x x'01;10' col; -- trailing note;
| COPY qpt FROM stdin;
| c;d;e
| a;b
| \.
| COPY iki FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qpw FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qpw FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY ucl FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ucl FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ucl FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY _nr FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY fdk FROM stdin;
| a;b
| c;d;e
| \.
|  -- trailing note;
| COPY nih FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _nr FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _nr FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ghm FROM stdin;
| a;b
| a;b
| \.
|
| + $$select 1;$;$$$ FROM t WHERE a = --/*;a;b
|  col;
| COPY qzf FROM stdin;
| 1	\N
| c;d;e
| \.
|
| COPY rqv FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ghm FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ghm FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT "END ;a;b--;" SELECT;
| + x'01;10' SELECT jl$h$v 1 U&';$;$;;' +;
| COPY qmn FROM stdin;
| c;d;e
| a;b
| \.
|
| , U&'/*;' SELECT CREATE FUNCTION vswc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;);($x$;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "+"
| position: 30
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "END ;a;b--;" SELECT;
< error
| 42703 column "END ;a;b--;" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| + x'01;10' SELECT jl$h$v 1 U&';$;$;;' +;
< error
| 42601 syntax error at or near "+"
| position: 2
> exec
| , /*/*/*$;$/ *;*/*/*/ 1 U&');(' col /*/*;$;$END ;*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*/*$;$/ *;*/*/*/ 1 U&');(' col /*/*;$;$END ;*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*/*$;$/ *;*/*/*/ 1 U&');(' col /*/*;$;$END ;*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + CREATE FUNCTION odqn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);(; '; SELECT 1; END FROM t WHERE a =;
| 1 (SELECT 1; SELECT 2) col $$/*;$$ FROM t WHERE a = --; $x$;
|  1;
| SELECT ";; " + --;;/*;--;
|  , CREATE FUNCTION nwoi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);();('; END AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION odqn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);(; '; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION odqn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);(; '; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 $grc$$x$;END ;end;$grc$ AS x CREATE FUNCTION elyn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;'; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $grc$$x$;END ;end;$grc$ AS x CREATE FUNCTION elyn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;'; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $grc$$x$;END ;end;$grc$ AS x CREATE FUNCTION elyn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;;;'; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY oma FROM stdin;
| a;b
| a;b
| \.
|
| , /*;;--;*/ AS x --; a;b
|  col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY oma FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY oma FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 --end;;;
|  SELECT X'01;10' FROM t WHERE a = U&';a;b' +; COPY usc FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --end;;;
|  SELECT X'01;10' FROM t WHERE a = U&';a;b' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --end;;;
|  SELECT X'01;10' FROM t WHERE a = U&';a;b' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) SELECT; -- trailing note;
| COPY tpc FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| SELECT le$s$_ AS x --$;$
|  FROM t WHERE a = --$;$
|  SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 51
> exec
| BEGIN
< ok 0 -
> exec
| SELECT le$s$_ AS x --$;$
|  FROM t WHERE a = --$;$
|  SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 51
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT le$s$_ AS x --$;$
|  FROM t WHERE a = --$;$
|  SELECT;
< error
| 42601 syntax error at or near "SELECT"
| position: 51
> exec
| COPY dvj FROM stdin;
| c;d;e
| c;d;e
| \.
|
| + U&';;--;' AS x (SELECT 1; SELECT 2) SELECT e'\n;' ,;
| 1 "select 1;);(--;" AS x /*a;b*/ col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dvj FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dvj FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT CREATE FUNCTION cqy_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; END col $$end;select 1;; $$ col (SELECT 1; SELECT 2) +;
| 1 B'01;10' + lg$s$j col (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION cqy_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; END col $$end;select 1;; $$ col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION cqy_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; END col $$end;select 1;; $$ col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| AS x CREATE FUNCTION luod() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b);(;'; SELECT 1; END 1;col CREATE FUNCTION xowv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;);('; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION luod() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b);(;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION luod() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b);(;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY vdw FROM stdin;
| a;b
| c;d;e
| \.
|
| + E'\\' ,;
| AS x (SELECT 1; SELECT 2) , (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vdw FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vdw FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col U&'; ; --;' FROM t WHERE a = CREATE FUNCTION dgtr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; a;b'; SELECT 1; END ,;
| , xy$p$a AS x CREATE FUNCTION bcha() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;end;END ;'; SELECT 1; END SELECT "end;" SELECT;
| SELECT ';$x$;\' + --END ;
|  col; COPY unt FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'; ; --;' FROM t WHERE a = CREATE FUNCTION dgtr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; a;b'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'; ; --;' FROM t WHERE a = CREATE FUNCTION dgtr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; a;b'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY jpc FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jpc FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jpc FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = "end;; select 1;" 1 /*select 1;end;; */ col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "end;; select 1;" 1 /*select 1;end;; */ col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "end;; select 1;" 1 /*select 1;end;; */ col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , ';' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ';' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ';' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col (SELECT 1; SELECT 2) +;
| FROM t WHERE a = CREATE FUNCTION jmai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END FROM t WHERE a = "$;$$x$;" AS x;
| 1 $k$$x$;$x$;$k$ 1 CREATE FUNCTION wntm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;END ;'; SELECT 1; END FROM t WHERE a =;
| COPY lza FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x wd$t$n AS x CREATE FUNCTION iyhe() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END ,;
| FROM t WHERE a = "; select 1;);(" FROM t WHERE a =; -- trailing note;
| , --/*;END ;
|  SELECT lr$f$m 1 /*$x$;$x$;*/ +; AS x /*/*/*;;select 1;;*/*/*/ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x wd$t$n AS x CREATE FUNCTION iyhe() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x wd$t$n AS x CREATE FUNCTION iyhe() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col '/*;' 1;col E'\';\';\\' + E''';'';'';' col vn$i$o AS x;SELECT /*/*/*END ;$x$;*/*/*/ + /*/*);(a;b*/*/ SELECT;col /*/*;;*/*/ + 'select 1;;\' SELECT /*/*a;b*/*/ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '/*;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '/*;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY us_ FROM stdin;
| a;b
| c;d;e
| \.
|
| col CREATE FUNCTION znmd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';END ;'; END SELECT E'\';' FROM t WHERE a =;SELECT --a;b
|  + CREATE FUNCTION tvbj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;end;$x$;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY us_ FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY us_ FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ctp FROM stdin;
| c;d;e
| 1	\N
| \.
|
| col (SELECT 1; SELECT 2) FROM t WHERE a = $$);(select 1;END ;$$ +; AS x CREATE FUNCTION ddhj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT;
| , U&';' SELECT U&'end;select 1;' +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ctp FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ctp FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x "select 1;select 1;--;" col;
| 1 CREATE FUNCTION raxf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; /*;'; SELECT 1; END col (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "select 1;select 1;--;" col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "select 1;select 1;--;" col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col e'\\\n;' + CREATE FUNCTION mcwv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;'; END 1;
| COPY llu FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e'\\\n;' + CREATE FUNCTION mcwv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;'; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e'\\\n;' + CREATE FUNCTION mcwv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;select 1;'; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col (SELECT 1; SELECT 2) SELECT nd$m$p ,;
| FROM t WHERE a = '$;$' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) SELECT nd$m$p ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) SELECT nd$m$p ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 E'\n;\n;' AS x e'\\' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 E'\n;\n;' AS x e'\\' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 E'\n;\n;' AS x e'\\' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x --$x$;select 1;select 1;
|  FROM t WHERE a = /*/*END ;*/*/ +;
| col (SELECT 1; SELECT 2) col E'\n;'';' + $kdjr$$x$;select 1;;$kdjr$ ,;
| COPY sps FROM stdin;
| c;d;e
| a;b
| \.
|
| col CREATE FUNCTION umzf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;a;bselect 1;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --$x$;select 1;select 1;
|  FROM t WHERE a = /*/*END ;*/*/ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --$x$;select 1;select 1;
|  FROM t WHERE a = /*/*END ;*/*/ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 e''';' + --;;--;
|  col E'\\\n;\\' SELECT;
| SELECT CREATE FUNCTION pacf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;END ;;;'; SELECT 1; END + (SELECT 1; SELECT 2) 1;
| COPY erd FROM stdin;
| 1	\N
| a;b
| \.
|
| col "$x$;" + --);(/*;
|  +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e''';' + --;;--;
|  col E'\\\n;\\' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e''';' + --;;--;
|  col E'\\\n;\\' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) AS x E'\';' AS x "$x$;" +;
| FROM t WHERE a = "$x$;" AS x (SELECT 1; SELECT 2) FROM t WHERE a = "$;$);();(" 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) AS x E'\';' AS x "$x$;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) AS x E'\';' AS x "$x$;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY fon FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| COPY _pk FROM stdin;
| 1	\N
| 1	\N
| \.
| AS x 'END ;/*;select 1;' 1 $qzx$--;--;a;b@$qzx$ AS x /*a;b);(*/ ,;
| COPY ikf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fon FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fon FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x /*/*/*;;a;b$;$*/*/*/ FROM t WHERE a = /*end;*/ SELECT e'\';' 1;
| col $wy$END ;; $wy$ col CREATE FUNCTION ujoc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; END + (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*;;a;b$;$*/*/*/ FROM t WHERE a = /*end;*/ SELECT e'\';' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*;;a;b$;$*/*/*/ FROM t WHERE a = /*end;*/ SELECT e'\';' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT U&'end;END ;' FROM t WHERE a = $x$--;end;@;@$x$ FROM t WHERE a = nd$m$b +;
| COPY iop FROM stdin;
| a;b
| a;b
| \.
|  1 ue$_$m SELECT --);(end;;
|  SELECT CREATE FUNCTION kbbn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END 1;
| SELECT CREATE FUNCTION pbtb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;END ;;;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 56
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'end;END ;' FROM t WHERE a = $x$--;end;@;@$x$ FROM t WHERE a = nd$m$b +;
< error
| 42601 syntax error at or near "FROM"
| position: 56
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'end;END ;' FROM t WHERE a = $x$--;end;@;@$x$ FROM t WHERE a = nd$m$b +;
< error
| 42601 syntax error at or near "FROM"
| position: 56
> exec
| + U&'--;' col /*/*/*END ;*/*/*/ col;
| SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'--;' col /*/*/*END ;*/*/*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'--;' col /*/*/*END ;*/*/*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col (SELECT 1; SELECT 2) , rm$h$s , --/*;--;);(
|  +;
| COPY sxz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) , rm$h$s , --/*;--;);(
|  +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) , rm$h$s , --/*;--;);(
|  +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col '$x$;' + ");(select 1;END ;" col (SELECT 1; SELECT 2) +; -- trailing note;
| COPY rvb FROM stdin;
| a;b
| c;d;e
| \.
|
| FROM t WHERE a = B'01;10' + '--;--;END ;\' SELECT rp$z$k ,;
| , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '$x$;' + ");(select 1;END ;" col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '$x$;' + ");(select 1;END ;" col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , e'\';\n;\';' SELECT U&';' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\';\n;\';' SELECT U&';' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\';\n;\';' SELECT U&';' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY ael FROM stdin;
| a;b
| 1	\N
| \.
|
| 1 $saq$; );($saq$ FROM t WHERE a =;
| + "a;b);(" +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ael FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ael FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = U&'$;$a;b' 1 $$end;$$ , (SELECT 1; SELECT 2) 1;COPY ohf FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'$;$a;b' 1 $$end;$$ , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'$;$a;b' 1 $$end;$$ , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x "$x$;" FROM t WHERE a = e''';\n;\';' col;
| COPY cwk FROM stdin;
| 1	\N
| a;b
| \.
|  FROM t WHERE a = $$--;;;; $$ col --end;$;$
|  + --a;b;--;
|  col;COPY zyr FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "$x$;" FROM t WHERE a = e''';\n;\';' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "$x$;" FROM t WHERE a = e''';\n;\';' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT tp$h$i , e'\n;'';\\' AS x;
| SELECT E''';\';\n;' ,;
< error
| 42601 syntax error at or near ";"
| position: 56
> exec
| BEGIN
< ok 0 -
> exec
| SELECT tp$h$i , e'\n;'';\\' AS x;
< error
| 42703 column "tp$h$i" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| SELECT E''';\';\n;' ,;
< error
| 42601 syntax error at or near ";"
| position: 23
> exec
| + --;;$x$;
|  SELECT --$;$end;
|  col; -- trailing note;
| + CREATE FUNCTION aubp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END +;
| + X'01;10' 1;
| col '; ' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --;;$x$;
|  SELECT --$;$end;
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --;;$x$;
|  SELECT --$;$end;
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY mba FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mba FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mba FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col E'\';'';' SELECT;SELECT (SELECT 1; SELECT 2) FROM t WHERE a =;AS x E'\\\\'';' FROM t WHERE a = '$x$;' , e''';' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col E'\';'';' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col E'\';'';' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT --;;--;
|  col; -- trailing note;
| COPY qtb FROM stdin;
| a;b
| 1	\N
| \.
|
| AS x '--;--;\' col;col $qorh$--;;;$qorh$ AS x;
< error
| 42601 syntax error at or near "a"
| position: 61
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --;;--;
|  col;
< error
| 42703 column "col" does not exist
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
|  -- trailing note;
| COPY qtb FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 41
> exec
| COPY idv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY idv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY idv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , /*/*/*;;end;*/*/*/ AS x E'\n;' FROM t WHERE a =; -- trailing note;
| col ----;
|  , /*/*/ *;$;$;;*/*/ col /*/*/*$;$*/*/*/ FROM t WHERE a =;
| , e'\';' FROM t WHERE a = /*/*/*;*/*/*/ AS x;
| SELECT "END ;$;$;;" AS x '/*;' AS x _b$t$f ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*/*;;end;*/*/*/ AS x E'\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*/*;;end;*/*/*/ AS x E'\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + "a;b;;" , gw$p$k +;
| FROM t WHERE a = U&'end;);(END ;' FROM t WHERE a = e'\';'';' , b'01;10' ,;
| FROM t WHERE a = zg$w$a SELECT x'01;10' , b'01;10' SELECT; -- trailing note;
| , (SELECT 1; SELECT 2) + "$;$end;);(" col jh$x$x AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "a;b;;" , gw$p$k +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "a;b;;" , gw$p$k +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) col; -- trailing note;
| + x'01;10' FROM t WHERE a = (SELECT 1; SELECT 2) 1 ne$w$u AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT ';;;select 1;\' SELECT; , b'01;10' , $$);(--;END ;$$ AS x;
< error
| 42601 syntax error at or near ","
| position: 32
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ';;;select 1;\' SELECT;
< ok 1 -
> exec
|  , b'01;10' , $$);(--;END ;$$ AS x;
< error
| 42601 syntax error at or near ","
| position: 2
> exec
| ROLLBACK
< ok 0 -
> exec
|  , b'01;10' , $$);(--;END ;$$ AS x;
< error
| 42601 syntax error at or near ","
| position: 2
> exec
| , ----;--;end;
|  col;
| col '; ' 1 CREATE FUNCTION sqdv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END SELECT '; ;;select 1;\' col; -- trailing note;
| 1 E'\\'';\n;' AS x CREATE FUNCTION wjap() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ----;--;end;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ----;--;end;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , e'\';' AS x --; a;b
|  , $$/*;end;$$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\';' AS x --; a;b
|  , $$/*;end;$$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\';' AS x --; a;b
|  , $$/*;end;$$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY nkd FROM stdin;
| 1	\N
| 1	\N
| \.
|
| , '; /*;end;' SELECT $$a;b;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY nkd FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY nkd FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY gsz FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gsz FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gsz FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , CREATE FUNCTION rcke() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;END ;'; SELECT 1; END 1 x'01;10' + /*/*end;END ;$x$;*/*/ SELECT;
| + $eg$a;b/*;$eg$ col ';' AS x (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION rcke() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;END ;'; SELECT 1; END 1 x'01;10' + /*/*end;END ;$x$;*/*/ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION rcke() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;END ;'; SELECT 1; END 1 x'01;10' + /*/*end;END ;$x$;*/*/ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY vud FROM stdin;
| c;d;e
| c;d;e
| \.
|
| SELECT "end;);(" col (SELECT 1; SELECT 2) SELECT;
| SELECT ');(/*;; \' 1 CREATE FUNCTION wttt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$--;;'; END col;
| 1 (SELECT 1; SELECT 2) + /*/*/*$;$$x$;END ;*/*/*/ SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vud FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vud FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x U&'select 1;$;$' SELECT;
| COPY me_ FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| COPY dyt FROM stdin;
| 1	\N
| 1	\N
| \.
|  col oh$b$u FROM t WHERE a = $ai$$;$a;b$ai$ FROM t WHERE a = B'01;10' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'select 1;$;$' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'select 1;$;$' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , b'01;10' AS x;
| AS x (SELECT 1; SELECT 2) 1;
| COPY udg FROM stdin;
| 1	\N
| c;d;e
| \.
|
| , /*/*end;end;*/*/ col ');(/*;$x$;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , b'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , b'01;10' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY gst FROM stdin;
| 1	\N
| a;b
| \.
|
| + U&'END ;;' 1;
| AS x --select 1;
|  SELECT ';);(\' col au$j$i +;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gst FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gst FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 $cwjx$end;$x$;a;b@$cwjx$ FROM t WHERE a = (SELECT 1; SELECT 2) AS x '--;' ,;
| col tc$g$v FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $cwjx$end;$x$;a;b@$cwjx$ FROM t WHERE a = (SELECT 1; SELECT 2) AS x '--;' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $cwjx$end;$x$;a;b@$cwjx$ FROM t WHERE a = (SELECT 1; SELECT 2) AS x '--;' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 'END ;' col b'01;10' FROM t WHERE a =;
| COPY qox FROM stdin;
| a;b
| c;d;e
| \.
|  COPY pty FROM stdin;
| c;d;e
| 1	\N
| \.
|
| AS x e'\n;\n;\n;' col $$; $$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 'END ;' col b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 'END ;' col b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , CREATE FUNCTION bbft() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$/*;;'; SELECT 1; END 1 (SELECT 1; SELECT 2) AS x; AS x rc$_$e +;
| SELECT '; /*;' 1 ci$f$k +; 1 $z$end;a;bselect 1;$z$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION bbft() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$/*;;'; SELECT 1; END 1 (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION bbft() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$/*;;'; SELECT 1; END 1 (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x --end;select 1;;
|  ,; COPY nzo FROM stdin;
| 1	\N
| 1	\N
| \.
|
| col ';a;b; \' 1 /*/*$x$;$x$;select 1;*/*/ col $$END ;$x$;$$ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --end;select 1;;
|  ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --end;select 1;;
|  ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT U&';;; ; ' SELECT $dm$;$;$a;b$dm$ ,; -- trailing note;
| COPY _dc FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "$dm$;$;$a;b$dm$"
| position: 26
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&';;; ; ' SELECT $dm$;$;$a;b$dm$ ,;
< error
| 42601 syntax error at or near "$dm$;$;$a;b$dm$"
| position: 26
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&';;; ; ' SELECT $dm$;$;$a;b$dm$ ,;
< error
| 42601 syntax error at or near "$dm$;$;$a;b$dm$"
| position: 26
> exec
| COPY lox FROM stdin;
| c;d;e
| a;b
| \.
| col --a;b
|  col x'01;10' 1 --/*;select 1;;
|  SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lox FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lox FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x B'01;10' AS x;
| , --);(
|  FROM t WHERE a = CREATE FUNCTION rodc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; END 1;
| FROM t WHERE a = CREATE FUNCTION k_uy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; ;;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x B'01;10' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x B'01;10' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) + '--;$;$' AS x e'\n;'';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) + '--;$;$' AS x e'\n;'';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) + '--;$;$' AS x e'\n;'';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , X'01;10' , b'01;10' , /*; $;$;;*/ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , X'01;10' , b'01;10' , /*; $;$;;*/ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , X'01;10' , b'01;10' , /*; $;$;;*/ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + CREATE FUNCTION ptre() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; ;'; END AS x E'\';\n;' + (SELECT 1; SELECT 2) SELECT;
| COPY xuq FROM stdin;
| a;b
| 1	\N
| \.
|
| AS x $$select 1;);(select 1;$$ col;
| COPY wyv FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION ptre() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; ;'; END AS x E'\';\n;' + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION ptre() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; ;'; END AS x E'\';\n;' + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT sa$n$a SELECT ";;" +; -- trailing note;
| 1 'a;b; );(' , $$END ;$$ col; -- trailing note;
| AS x /*/*/*select 1;);(*/*/*/ AS x; -- trailing note;
| + e'\n;\\\n;' + /*/*/*/ *;; a;b*/*/*/ FROM t WHERE a = E'\\' col;
< error
| 42601 syntax error at or near "";;""
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT sa$n$a SELECT ";;" +;
< error
| 42601 syntax error at or near "";;""
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT sa$n$a SELECT ";;" +;
< error
| 42601 syntax error at or near "";;""
| position: 22
> exec
| FROM t WHERE a = E'\n;\';\';' col CREATE FUNCTION fxls() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b--;'; END AS x (SELECT 1; SELECT 2) +;
| COPY lmh FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E'\n;\';\';' col CREATE FUNCTION fxls() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b--;'; END AS x (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\n;\';\';' col CREATE FUNCTION fxls() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b--;'; END AS x (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col E'\';\\\';' FROM t WHERE a =;COPY sib FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col E'\';\\\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col E'\';\\\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + pb$m$_ + /*/*end;*/*/ SELECT /*a;b; */ FROM t WHERE a =;
| + '/*;/*;\' col; AS x $vqf$;;end;;;$vqx$$vqf$ col "end;--;" col;
| , (SELECT 1; SELECT 2) col CREATE FUNCTION ylyq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b;;'; END SELECT B'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + pb$m$_ + /*/*end;*/*/ SELECT /*a;b; */ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + pb$m$_ + /*/*end;*/*/ SELECT /*a;b; */ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = $lpl$end;$lpx$$lpl$ ,;
| COPY bb_ FROM stdin;
| c;d;e
| 1	\N
| \.
|
| AS x E'\n;\';\\' AS x E'\';'';\n;' col eo$b$y AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $lpl$end;$lpx$$lpl$ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $lpl$end;$lpx$$lpl$ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 x'01;10' +;
| , ';END ;' +;
| col CREATE FUNCTION xdsx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END FROM t WHERE a = "$x$;$x$;end;" +;
| FROM t WHERE a = e'\\'';\';' , U&'END ;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 x'01;10' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 x'01;10' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY zvi FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| col /*/*a;b*/*/ SELECT e'\n;' 1 B'01;10' AS x;
| SELECT --$;$
|  , CREATE FUNCTION bkxx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;/*;end;'; SELECT 1; END AS x $$$;$end;$;$$$ SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zvi FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zvi FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , $$select 1;$$ col;AS x "end;" col;AS x E'\\\\\n;' col;
| SELECT U&'$x$;' + b'01;10' SELECT --/*;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $$select 1;$$ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $$select 1;$$ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x "$x$;" 1 'a;b$;$' col;
| col --/*;
|  +;
| COPY sfr FROM stdin;
| c;d;e
| c;d;e
| \.
|
| 1 /*;*/ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "$x$;" 1 'a;b$;$' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "$x$;" 1 'a;b$;$' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY dtm FROM stdin;
| a;b
| a;b
| \.
|
| 1 kc$c$_ + CREATE FUNCTION owce() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$select 1;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dtm FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dtm FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY zhj FROM stdin;
| a;b
| c;d;e
| \.
|  -- trailing note;
| COPY mql FROM stdin;
| c;d;e
| 1	\N
| \.
| col U&'$x$;select 1;;' col /*/*/*$x$;$;$*/*/*/ +;
| AS x U&'; $x$;select 1;' AS x CREATE FUNCTION alox() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; --;'; END FROM t WHERE a = E'\\\\\';' ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zhj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zhj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = CREATE FUNCTION kauu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;;'; END ,; -- trailing note;
| FROM t WHERE a = U&'/*;;;' FROM t WHERE a =;
| FROM t WHERE a = U&'select 1;' SELECT b'01;10' FROM t WHERE a =;
| SELECT ';;\' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION kauu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;;'; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION kauu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;;'; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + ";;$x$;" SELECT;
| COPY mhr FROM stdin;
| c;d;e
| a;b
| \.
|
| FROM t WHERE a = rn$x$l col b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ";;$x$;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ";;$x$;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT 'select 1;$x$;' ,;
< error
| 42601 syntax error at or near ";"
| position: 25
> exec
| BEGIN
< ok 0 -
> exec
| SELECT 'select 1;$x$;' ,;
< error
| 42601 syntax error at or near ";"
| position: 25
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT 'select 1;$x$;' ,;
< error
| 42601 syntax error at or near ";"
| position: 25
> exec
| FROM t WHERE a = cu$d$e SELECT e'\';' FROM t WHERE a = '$;$select 1;' col;
| 1 U&');();(;;' AS x B'01;10' +;
| COPY ogo FROM stdin;
| 1	\N
| c;d;e
| \.
|  1 $$a;b;;a;b$$ col ps$l$e col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = cu$d$e SELECT e'\';' FROM t WHERE a = '$;$select 1;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = cu$d$e SELECT e'\';' FROM t WHERE a = '$;$select 1;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , /*/*; END ;a;b*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*; END ;a;b*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*; END ;a;b*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT x'01;10' FROM t WHERE a = --;;--;
|  + b'01;10' col; AS x --; /*;end;
|  , "end;END ;" col /*/*);(*/*/ +;
| col 'select 1;' col;
< error
| 42601 syntax error at or near "col"
| position: 54
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' FROM t WHERE a = --;;--;
|  + b'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 54
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT x'01;10' FROM t WHERE a = --;;--;
|  + b'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 54
> exec
| col U&'$x$;$;$' FROM t WHERE a = U&'select 1;$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'$x$;$;$' FROM t WHERE a = U&'select 1;$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'$x$;$;$' FROM t WHERE a = U&'select 1;$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY hed FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY grd FROM stdin;
| 1	\N
| a;b
| \.
|  COPY btp FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hed FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hed FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + b'01;10' SELECT;
| COPY ioy FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + b'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + b'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY jnd FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| col /*/*/*$;$end;;;*/*/*/ FROM t WHERE a = by$v$u ,; -- trailing note;
| FROM t WHERE a = --select 1;/*;
|  +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jnd FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jnd FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , $$/*;;$$ AS x $gfns$$;$/*;;$gfnx$$gfns$ SELECT; -- trailing note;
| AS x (SELECT 1; SELECT 2) 1 --end;$x$;
|  col $$select 1;$$ ,;
| COPY gng FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $$/*;;$$ AS x $gfns$$;$/*;;$gfnx$$gfns$ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $$/*;;$$ AS x $gfns$$;$/*;;$gfnx$$gfns$ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + --;;; $x$;
|  AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --;;; $x$;
|  AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --;;; $x$;
|  AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY beh FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY sau FROM stdin;
| a;b
| 1	\N
| \.
|
| col (SELECT 1; SELECT 2) col b'01;10' FROM t WHERE a = "--;END ;;;" col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY beh FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY beh FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x hd$z$t + x'01;10' 1;col E'\';\';' , ";" SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x hd$z$t + x'01;10' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x hd$z$t + x'01;10' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + $$/*;$$ AS x e'\';' ,;COPY ysq FROM stdin;
| a;b
| c;d;e
| \.
|  1 U&'/*;;;' col U&'; ' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $$/*;$$ AS x e'\';' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $$/*;$$ AS x e'\';' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x U&'/*;$;$--;' , ----;end;
|  SELECT;+ --$x$;
|  + /*/*/*; ;;*/*/*/ col; -- trailing note;
| col CREATE FUNCTION vaqu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;'; SELECT 1; END + U&'end;' , "/*;" +;COPY fsl FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'/*;$;$--;' , ----;end;
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'/*;$;$--;' , ----;end;
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x E'\';' FROM t WHERE a = "select 1;);(select 1;" SELECT '$x$;\' SELECT;
| + "; " SELECT /*/*--;*/*/ 1 "END ;$;$" SELECT; FROM t WHERE a = /*END ;*/ col U&');($;$' ,;
| SELECT x'01;10' 1 /*/*/*;;;);(*/*/*/ 1 '--;' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x E'\';' FROM t WHERE a = "select 1;);(select 1;" SELECT '$x$;\' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x E'\';' FROM t WHERE a = "select 1;);(select 1;" SELECT '$x$;\' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 CREATE FUNCTION urak() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;; '; SELECT 1; END col '$;$' FROM t WHERE a = /*/*;;;*/*/ FROM t WHERE a =; -- trailing note;
| 1 "$x$;);(/*;" FROM t WHERE a = CREATE FUNCTION ixqu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;'; END col; -- trailing note;
| SELECT /*;;; */ + CREATE FUNCTION dwxq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(select 1;$x$;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION urak() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;; '; SELECT 1; END col '$;$' FROM t WHERE a = /*/*;;;*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION urak() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;; '; SELECT 1; END col '$;$' FROM t WHERE a = /*/*;;;*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + tv$r$p SELECT /*/*);(*/*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + tv$r$p SELECT /*/*);(*/*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + tv$r$p SELECT /*/*);(*/*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x e''';'';\n;' AS x "select 1;" + om$q$o 1;
| COPY hjv FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x e''';'';\n;' AS x "select 1;" + om$q$o 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x e''';'';\n;' AS x "select 1;" + om$q$o 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY g_o FROM stdin;
| 1	\N
| c;d;e
| \.
|
| SELECT --$x$;
|  1;
| COPY wfb FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY g_o FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY g_o FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT X'01;10' AS x E'\\' SELECT;
< error
| 42601 syntax error at or near "E'\\\\'"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT X'01;10' AS x E'\\' SELECT;
< error
| 42601 syntax error at or near "E'\\\\'"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT X'01;10' AS x E'\\' SELECT;
< error
| 42601 syntax error at or near "E'\\\\'"
| position: 22
> exec
| SELECT U&';;);(END ;' +; col B'01;10' 1 '--;\' SELECT b'01;10' SELECT;
| SELECT --$x$;
|  AS x; -- trailing note;
| AS x (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ";"
| position: 24
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&';;);(END ;' +;
< error
| 42601 syntax error at or near ";"
| position: 24
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&';;);(END ;' +;
< error
| 42601 syntax error at or near ";"
| position: 24
> exec
| COPY mcq FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mcq FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mcq FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY dzs FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dzs FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dzs FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = "$x$;$;$;;" SELECT $$a;b--;$$ , '--;' FROM t WHERE a =;COPY xrx FROM stdin;
| a;b
| 1	\N
| \.
|
| col --/*;$x$;a;b
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "$x$;$;$;;" SELECT $$a;b--;$$ , '--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "$x$;$;$;;" SELECT $$a;b--;$$ , '--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + --$;$
|  SELECT "--;" FROM t WHERE a =; col --/*;
|  + (SELECT 1; SELECT 2) 1;SELECT --a;b
|  1;
| + "end;;;; " 1 $a$end;--;$x$$a$ + km$v$f col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --$;$
|  SELECT "--;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --$;$
|  SELECT "--;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , ix$_$r col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ix$_$r col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ix$_$r col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , '/*;;' , $$);(END ;; $$ , "END ;$x$;a;b" +;
| col /*/*/*; $;$select 1;*/*/*/ FROM t WHERE a = rt$q$x + (SELECT 1; SELECT 2) col;+ x'01;10' col; -- trailing note;
| AS x B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , '/*;;' , $$);(END ;; $$ , "END ;$x$;a;b" +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , '/*;;' , $$);(END ;; $$ , "END ;$x$;a;b" +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = rd$j$v FROM t WHERE a = (SELECT 1; SELECT 2) col /*/*end;END ;select 1;*/*/ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = rd$j$v FROM t WHERE a = (SELECT 1; SELECT 2) col /*/*end;END ;select 1;*/*/ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = rd$j$v FROM t WHERE a = (SELECT 1; SELECT 2) col /*/*end;END ;select 1;*/*/ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY eil FROM stdin;
| c;d;e
| a;b
| \.
|
| col $$END ;$$ 1 /*/*--;select 1;--;*/*/ , CREATE FUNCTION henb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END ,;AS x ir$l$b 1 ----;
|  ,;
| COPY xth FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY eil FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY eil FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , x'01;10' col;
| COPY fpp FROM stdin;
| 1	\N
| 1	\N
| \.
|
| col x'01;10' FROM t WHERE a =;1 X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , x'01;10' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , x'01;10' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT CREATE FUNCTION dfem() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT 'END ;end;END ;' AS x;
| FROM t WHERE a = U&'a;b;);(' AS x CREATE FUNCTION zrga() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;b'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION dfem() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT 'END ;end;END ;' AS x;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION dfem() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT 'END ;end;END ;' AS x;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| COPY uon FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY kmf FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY uon FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY uon FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) 1 oj$g$b 1;
| col "end;;;/*;" col (SELECT 1; SELECT 2) 1 'END ;--;$;$' AS x;
| , $$END ;; $;$$$ col e''';' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) 1 oj$g$b 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) 1 oj$g$b 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + X'01;10' AS x; -- trailing note;
| FROM t WHERE a = qc$z$i AS x ck$p$k AS x; -- trailing note;
| col --END ;;select 1;
|  AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + X'01;10' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + X'01;10' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + (SELECT 1; SELECT 2) SELECT CREATE FUNCTION qatw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END 1 /*/*/ *;END ;;;*/*/ SELECT;COPY n_o FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) SELECT CREATE FUNCTION qatw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END 1 /*/*/ *;END ;;;*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) SELECT CREATE FUNCTION qatw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END 1 /*/*/ *;END ;;;*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x ";;END ;" AS x (SELECT 1; SELECT 2) SELECT (SELECT 1; SELECT 2) 1;
| FROM t WHERE a = ----;--;;
|  ,; -- trailing note;
| , $u$;; $u$ , '$x$;' ,;, b'01;10' AS x '; END ;' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ";;END ;" AS x (SELECT 1; SELECT 2) SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ";;END ;" AS x (SELECT 1; SELECT 2) SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col $$$;$/*;a;b$$ AS x --END ;select 1;
|  + X'01;10' AS x;
| + --$x$;$x$;;
|  SELECT x'01;10' , '$x$;$x$;' +;
| 1 --$;$
|  SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$$;$/*;a;b$$ AS x --END ;select 1;
|  + X'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$$;$/*;a;b$$ AS x --END ;select 1;
|  + X'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 ";" FROM t WHERE a = ";" FROM t WHERE a = --$x$;;;
|  1;
| FROM t WHERE a = U&';;$;$' , o_$a$q + (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ";" FROM t WHERE a = ";" FROM t WHERE a = --$x$;;;
|  1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ";" FROM t WHERE a = ";" FROM t WHERE a = --$x$;;;
|  1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , E'\';\';' SELECT kd$p$g AS x ';' AS x; + --$;$a;b
|  FROM t WHERE a = E'\\'';\\' , $$$x$;a;bEND ;$$ 1;COPY arl FROM stdin;
| a;b
| a;b
| \.
|
| 1 X'01;10' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , E'\';\';' SELECT kd$p$g AS x ';' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , E'\';\';' SELECT kd$p$g AS x ';' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 --;
|  AS x --end;
|  SELECT /*a;b);($;$*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --;
|  AS x --end;
|  SELECT /*a;b);($;$*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --;
|  AS x --end;
|  SELECT /*a;b);($;$*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT CREATE FUNCTION yqpo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT CREATE FUNCTION xydj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END FROM t WHERE a =;
| COPY wft FROM stdin;
| c;d;e
| 1	\N
| \.
|
| FROM t WHERE a = b'01;10' SELECT;
| + $$;a;bEND ;$$ + '--;/*;/*;\' col CREATE FUNCTION zwhd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;select 1;'; END 1;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION yqpo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT CREATE FUNCTION xydj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION yqpo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT CREATE FUNCTION xydj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| col /*/*/ *;$x$;$;$*/*/ SELECT;, 'end;);(' FROM t WHERE a = (SELECT 1; SELECT 2) col; -- trailing note;
| COPY lic FROM stdin;
| a;b
| c;d;e
| \.
|
| SELECT CREATE FUNCTION btco() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; SELECT 1; END SELECT b'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*/ *;$x$;$;$*/*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*/ *;$x$;$;$*/*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col ----;$;$
|  SELECT ks$n$s FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ----;$;$
|  SELECT ks$n$s FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ----;$;$
|  SELECT ks$n$s FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = x'01;10' AS x;
| col ');(\' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = $$/*;$;$$$ col $$a;b$$ FROM t WHERE a =; -- trailing note;
| col /*$;$/ *;*/ FROM t WHERE a = CREATE FUNCTION nuhm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;);();('; SELECT 1; END ,;
| SELECT "a;b" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$/*;$;$$$ col $$a;b$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$/*;$;$$$ col $$a;b$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col U&';;' col "; --;END ;" 1;
| col "a;bEND ;" FROM t WHERE a =;
| COPY rhx FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&';;' col "; --;END ;" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&';;' col "; --;END ;" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY xuq FROM stdin;
| a;b
| c;d;e
| \.
|
| col E'\\'';' ,; -- trailing note;
| + $xj$;;/*;end;$xx$$xj$ FROM t WHERE a = ma$d$q SELECT e''';' SELECT;
| , ra$p$e SELECT E''';\';'';' 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xuq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xuq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , X'01;10' SELECT U&');(' AS x U&'$;$' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , X'01;10' SELECT U&');(' AS x U&'$;$' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , X'01;10' SELECT U&');(' AS x U&'$;$' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col /*/*a;b*/*/ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*a;b*/*/ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*a;b*/*/ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY agb FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY agb FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY agb FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 wi$_$w +; -- trailing note;
| col --;select 1;
|  AS x;
| AS x /*/*a;b*/*/ col "$;$$x$;$;$" SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 wi$_$w +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 wi$_$w +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = "$x$;$x$;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "$x$;$x$;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "$x$;$x$;" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY xtc FROM stdin;
| a;b
| 1	\N
| \.
|  AS x "; $x$;/*;" col ");(" +;
| col 'select 1;\' 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xtc FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xtc FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ahh FROM stdin;
| a;b
| a;b
| \.
| COPY lze FROM stdin;
| 1	\N
| 1	\N
| \.
| COPY mjp FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ahh FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ahh FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 --; 
|  FROM t WHERE a =; -- trailing note;
| COPY mqe FROM stdin;
| 1	\N
| c;d;e
| \.
|
| SELECT dl$b$x , --END ;$;$;;
|  col;
| COPY sej FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --; 
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --; 
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , E''';'';'';' + E'\\\';\';' 1 X'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , E''';'';'';' + E'\\\';\';' 1 X'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , E''';'';'';' + E'\\\';\';' 1 X'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY vam FROM stdin;
| a;b
| c;d;e
| \.
|  FROM t WHERE a = $$/*;$$ + /*select 1;;;*/ SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vam FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vam FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT $$;;$$ AS x; SELECT "a;b; " FROM t WHERE a = --a;b; $;$
|  , '$x$;$x$;\' col; FROM t WHERE a = (SELECT 1; SELECT 2) col /*/*);(*/*/ 1;
| 1 CREATE FUNCTION rxup() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 65
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$;;$$ AS x;
< ok 1 -
> exec
|  SELECT "a;b; " FROM t WHERE a = --a;b; $;$
|  , '$x$;$x$;\' col;
< error
| 42601 syntax error at or near ","
| position: 46
> exec
| ROLLBACK
< ok 0 -
> exec
|  SELECT "a;b; " FROM t WHERE a = --a;b; $;$
|  , '$x$;$x$;\' col;
< error
| 42601 syntax error at or near ","
| position: 46
> exec
| COPY kiy FROM stdin;
| 1	\N
| 1	\N
| \.
|
| col ";;" , CREATE FUNCTION qync() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;end;--;'; SELECT 1; END SELECT; -- trailing note;
| COPY mug FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kiy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kiy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY viz FROM stdin;
| 1	\N
| c;d;e
| \.
| FROM t WHERE a = CREATE FUNCTION amwc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;a;b;'; END 1 E'\n;' , ----;
|  col;
| 1 "END ;" SELECT (SELECT 1; SELECT 2) + B'01;10' +;
| col U&'--;$;$' col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY viz FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY viz FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY aoo FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY fvt FROM stdin;
| a;b
| c;d;e
| \.
|
| FROM t WHERE a = /*/*select 1;/ *;*/*/ + $u$$x$;);(/*;$u$ AS x X'01;10' ,; 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY aoo FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY aoo FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT iu$i$i col E'\';' +; -- trailing note;
| col '/*;' col ');($x$;' SELECT $k_j_$);(;$k_jx$$k_j_$ ,;
< error
| 42601 syntax error at or near "E'\\';'"
| position: 19
> exec
| BEGIN
< ok 0 -
> exec
| SELECT iu$i$i col E'\';' +;
< error
| 42601 syntax error at or near "E'\\';'"
| position: 19
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT iu$i$i col E'\';' +;
< error
| 42601 syntax error at or near "E'\\';'"
| position: 19
> exec
| AS x U&');(' , jd$h$k FROM t WHERE a = '/*;; ;\' AS x;
| FROM t WHERE a = "$;$a;bselect 1;" AS x '--;end;\' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&');(' , jd$h$k FROM t WHERE a = '/*;; ;\' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&');(' , jd$h$k FROM t WHERE a = '/*;; ;\' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x '; select 1;' FROM t WHERE a = b'01;10' ,;
| 1 (SELECT 1; SELECT 2) SELECT '/*;/*;select 1;' col CREATE FUNCTION wxfj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;);('; SELECT 1; END col; COPY orz FROM stdin;
| c;d;e
| a;b
| \.
|
| SELECT E'\';\\' FROM t WHERE a = uq$h$x , e''';'';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x '; select 1;' FROM t WHERE a = b'01;10' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x '; select 1;' FROM t WHERE a = b'01;10' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) , /*END ;; ;*/ FROM t WHERE a = "END ;/*;" SELECT;
| + "$;$; ; " AS x;
| 1 x'01;10' col b'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) , /*END ;; ;*/ FROM t WHERE a = "END ;/*;" SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) , /*END ;; ;*/ FROM t WHERE a = "END ;/*;" SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY oax FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY oax FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY oax FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 '; --;\' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '; --;\' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '; --;\' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT "$;$/*;;;" AS x --;; 
|  1 B'01;10' FROM t WHERE a =;
| FROM t WHERE a = E'\';\n;\';' 1; -- trailing note;
| COPY rwc FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 31
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "$;$/*;;;" AS x --;; 
|  1 B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 31
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "$;$/*;;;" AS x --;; 
|  1 B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 31
> exec
| , ----;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ----;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ----;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY ziq FROM stdin;
| a;b
| 1	\N
| \.
|
| , "; select 1;$;$" FROM t WHERE a = $$$x$;a;b;;$$ SELECT ");(select 1;select 1;" 1;col "end;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ziq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ziq FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = ';' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ';' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ';' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY tff FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY tff FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY tff FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , pe$v$g col e''';' ,; -- trailing note;
| FROM t WHERE a = E''';' AS x;
| SELECT (SELECT 1; SELECT 2) + /*/*a;b*/*/ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , pe$v$g col e''';' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , pe$v$g col e''';' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x ";;" col CREATE FUNCTION wquv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; );(/*;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ";;" col CREATE FUNCTION wquv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; );(/*;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ";;" col CREATE FUNCTION wquv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; );(/*;'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY apx FROM stdin;
| 1	\N
| c;d;e
| \.
|
| COPY vaj FROM stdin;
| 1	\N
| c;d;e
| \.
| COPY oql FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY apx FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY apx FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , E'\n;\\'';' , X'01;10' 1 CREATE FUNCTION uzlb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END SELECT;
| , $v$select 1;$v$ AS x;
| COPY ily FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , E'\n;\\'';' , X'01;10' 1 CREATE FUNCTION uzlb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , E'\n;\\'';' , X'01;10' 1 CREATE FUNCTION uzlb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col X'01;10' , b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col X'01;10' , b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col X'01;10' , b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 /*/*;;;*/*/ 1; + $$$;$$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*;;;*/*/ 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*;;;*/*/ 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY hjr FROM stdin;
| a;b
| a;b
| \.
|
| , '/*;' col --select 1;;
|  SELECT E'\\\\' col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hjr FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hjr FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = $rj$$x$;$;$select 1;$rj$ SELECT (SELECT 1; SELECT 2) SELECT;
| , $$end;/*;$$ FROM t WHERE a =;
| 1 ds$h$w SELECT "$;$select 1;select 1;" +; -- trailing note;
| , e'\\\';' , /*--;end;END ;*/ col CREATE FUNCTION _hpm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $rj$$x$;$;$select 1;$rj$ SELECT (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $rj$$x$;$;$select 1;$rj$ SELECT (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 x'01;10' col;1 CREATE FUNCTION rlzi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b;;end;'; SELECT 1; END 1;
| COPY cou FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 x'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 x'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col "end;" AS x --);(
|  FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "end;" AS x --);(
|  FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "end;" AS x --);(
|  FROM t WHERE a = (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) SELECT $lrb$);(select 1;; $lrb$ ,;
| COPY ucw FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) SELECT $lrb$);(select 1;; $lrb$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) SELECT $lrb$);(select 1;; $lrb$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY pxf FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pxf FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pxf FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = e'\\' FROM t WHERE a = (SELECT 1; SELECT 2) SELECT;
| SELECT '$;$' FROM t WHERE a =;
| col e'\n;\';\n;' 1 '$x$;' FROM t WHERE a = ';;$;$);(\' ,;
| AS x CREATE FUNCTION vqfd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END col U&'--;$x$;$;$' AS x $$);($$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = e'\\' FROM t WHERE a = (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = e'\\' FROM t WHERE a = (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT --;
|  1 /*/*/*--;; */*/*/ +;COPY csx FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY jln FROM stdin;
| 1	\N
| 1	\N
| \.
|
| + ';;$;$$x$;' AS x E'\n;\';' AS x;
< error
| 42601 syntax error at or near ";"
| position: 34
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --;
|  1 /*/*/*--;; */*/*/ +;
< error
| 42601 syntax error at or near ";"
| position: 34
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --;
|  1 /*/*/*--;; */*/*/ +;
< error
| 42601 syntax error at or near ";"
| position: 34
> exec
| + '$;$);(end;\' + (SELECT 1; SELECT 2) , eo$y$v col;
| SELECT (SELECT 1; SELECT 2) AS x;AS x (SELECT 1; SELECT 2) +; + $dda$select 1;--;$dda$ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + '$;$);(end;\' + (SELECT 1; SELECT 2) , eo$y$v col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + '$;$);(end;\' + (SELECT 1; SELECT 2) , eo$y$v col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , "--;$;$" FROM t WHERE a = CREATE FUNCTION lmgv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;;'; SELECT 1; END col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "--;$;$" FROM t WHERE a = CREATE FUNCTION lmgv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;;'; SELECT 1; END col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "--;$;$" FROM t WHERE a = CREATE FUNCTION lmgv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;;'; SELECT 1; END col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , X'01;10' + CREATE FUNCTION flpf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;'; END 1 CREATE FUNCTION gfja() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x; AS x --$;$
|  ,;
| , U&');(select 1;' FROM t WHERE a = /*);(select 1;END ;*/ AS x CREATE FUNCTION fiwg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , X'01;10' + CREATE FUNCTION flpf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;'; END 1 CREATE FUNCTION gfja() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , X'01;10' + CREATE FUNCTION flpf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;'; END 1 CREATE FUNCTION gfja() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + "$;$" col e'\\' FROM t WHERE a =;
| col /*/*/*;;$;$a;b*/*/*/ , $ikyu$end;select 1;$ikyx$$ikyu$ + /*/*;;;;select 1;*/*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "$;$" col e'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "$;$" col e'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , CREATE FUNCTION jxpr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;--;'; END 1 U&'a;b;;;' FROM t WHERE a =;
| FROM t WHERE a = '; ;' , '; ' FROM t WHERE a =;
| COPY nxb FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION jxpr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;--;'; END 1 U&'a;b;;;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION jxpr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;--;'; END 1 U&'a;b;;;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col B'01;10' FROM t WHERE a =;
| SELECT e'\';' , b'01;10' col ";;);(" ,;
| col /*/*/*--;$x$;*/*/*/ , /*/*;;);(*/*/ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 ----;;;
|  FROM t WHERE a = e'\n;\';'';' ,; -- trailing note;
| col k_$m$i 1 "; " AS x B'01;10' ,;AS x n_$g$a col 'select 1;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ----;;;
|  FROM t WHERE a = e'\n;\';'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ----;;;
|  FROM t WHERE a = e'\n;\';'';' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT "/*;);(" col (SELECT 1; SELECT 2) 1;
| COPY f_p FROM stdin;
| 1	\N
| 1	\N
| \.
|
| AS x ");(end;$;$" col "--;; $;$" + /*/*/*END ;*/*/*/ ,; + U&';;a;b' col (SELECT 1; SELECT 2) col $dub$END ;a;b$dub$ AS x;
< error
| 42601 syntax error at or near "("
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "/*;);(" col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "("
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "/*;);(" col (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "("
| position: 21
> exec
| COPY c_x FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY c_x FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY c_x FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY b_a FROM stdin;
| a;b
| 1	\N
| \.
|  , /*/*/*$;$*/*/*/ FROM t WHERE a = E''';' 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY b_a FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY b_a FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY axg FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY hxn FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY axg FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY axg FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = $gs$/*;/*;$gs$ AS x;, U&'a;bEND ;/*;' , mo$g$e 1 X'01;10' FROM t WHERE a =;
| + b_$a$c AS x $$;;);();($$ , U&'$x$;$;$' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $gs$/*;/*;$gs$ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $gs$/*;/*;$gs$ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = b'01;10' 1 CREATE FUNCTION pfiv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bselect 1;'; SELECT 1; END SELECT X'01;10' ,;
| SELECT $$select 1;/*;$$ + $$;$$ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = b'01;10' 1 CREATE FUNCTION pfiv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bselect 1;'; SELECT 1; END SELECT X'01;10' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = b'01;10' 1 CREATE FUNCTION pfiv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bselect 1;'; SELECT 1; END SELECT X'01;10' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) ,;col (SELECT 1; SELECT 2) 1 U&';;$;$' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , E'\\\\\n;' +;FROM t WHERE a = $to$select 1;$to$ SELECT CREATE FUNCTION summ() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$END ;'; END , "--;end;" col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , E'\\\\\n;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , E'\\\\\n;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT "a;ba;b" 1 (SELECT 1; SELECT 2) SELECT;
| COPY ens FROM stdin;
| c;d;e
| c;d;e
| \.
|
| FROM t WHERE a = ';;' , x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "a;ba;b" 1 (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "a;ba;b" 1 (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| SELECT y_$h$k +;col no$x$d +;
| 1 E''';\';'';' , /*/*/*end;END ;;;*/*/*/ 1;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| BEGIN
< ok 0 -
> exec
| SELECT y_$h$k +;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT y_$h$k +;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| FROM t WHERE a = /*end;$x$;*/ +;
| COPY ruq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*end;$x$;*/ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*end;$x$;*/ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x /*/*/*/ *;; */*/*/ +; , '$x$;end;);(\' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*/ *;; */*/*/ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*/ *;; */*/*/ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , --end;select 1;; 
|  1 --$;$
|  col uc$z$a SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --end;select 1;; 
|  1 --$;$
|  col uc$z$a SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --end;select 1;; 
|  1 --$;$
|  col uc$z$a SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 ----;/*;/*;
|  SELECT "end;END ;$;$" +; col U&'--;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ----;/*;/*;
|  SELECT "end;END ;$;$" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ----;/*;/*;
|  SELECT "end;END ;$;$" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , $wfr$;select 1;$wfr$ AS x --select 1;$x$;;;
|  col --; 
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $wfr$;select 1;$wfr$ AS x --select 1;$x$;;;
|  col --; 
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $wfr$;select 1;$wfr$ AS x --select 1;$x$;;;
|  col --; 
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) + /*$x$;$x$;*/ + U&'end;' +; -- trailing note;
| col E'\';' +;
| COPY vhr FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY wbq FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) + /*$x$;$x$;*/ + U&'end;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) + /*$x$;$x$;*/ + U&'end;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col /*/*$x$;a;b;*/*/ , e'\';\\\n;' FROM t WHERE a = "$x$;--;" FROM t WHERE a =;
| + $$select 1;$x$;$$ AS x;
| SELECT E'\\'';' FROM t WHERE a = '$x$;select 1;' 1;SELECT --$x$;--;select 1;
|  , ----;
|  , qs$n$x col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*$x$;a;b;*/*/ , e'\';\\\n;' FROM t WHERE a = "$x$;--;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*$x$;a;b;*/*/ , e'\';\\\n;' FROM t WHERE a = "$x$;--;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , --a;b
|  + U&'select 1;' 1;
| FROM t WHERE a = CREATE FUNCTION nrku() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --a;b
|  + U&'select 1;' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --a;b
|  + U&'select 1;' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = --);(;a;b
|  + ");(" FROM t WHERE a = "$x$;;END ;" 1;SELECT (SELECT 1; SELECT 2) AS x; -- trailing note;
| SELECT '$;$' AS x U&');(' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --);(;a;b
|  + ");(" FROM t WHERE a = "$x$;;END ;" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --);(;a;b
|  + ");(" FROM t WHERE a = "$x$;;END ;" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + ft$u$l , /*/*$;$a;b*/*/ 1;1 B'01;10' + U&'--;' FROM t WHERE a =;
| COPY kxi FROM stdin;
| a;b
| a;b
| \.
|
| AS x /*/*$;$*/*/ , (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ft$u$l , /*/*$;$a;b*/*/ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ft$u$l , /*/*$;$a;b*/*/ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT "/*;" + /*/*/*$x$;END ;select 1;*/*/*/ 1;
| + --; $;$; 
|  SELECT B'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 50
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "/*;" + /*/*/*$x$;END ;select 1;*/*/*/ 1;
< error
| 42703 column "/*;" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| + --; $;$; 
|  SELECT B'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 2
> exec
| FROM t WHERE a = oq$c$l SELECT;
| + $$a;b$$ ,;AS x jd$d$g + $jdoz$end;$x$;$jdoz$ FROM t WHERE a = U&';);(--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = oq$c$l SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = oq$c$l SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x si$f$z col (SELECT 1; SELECT 2) AS x;
| SELECT ');(' FROM t WHERE a = X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x si$f$z col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x si$f$z col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + CREATE FUNCTION qnrl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END , B'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION qnrl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END , B'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION qnrl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END , B'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY hva FROM stdin;
| 1	\N
| a;b
| \.
|
| 1 U&'select 1;' col e'\';'';'';' , '$;$$x$;; \' ,; -- trailing note;
| SELECT gb$w$m ,; COPY pcx FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hva FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hva FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY ysi FROM stdin;
| c;d;e
| c;d;e
| \.
| col CREATE FUNCTION tbqo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;$;$'; END SELECT /*$x$;/ *;a;b*/ + U&'$;$END ;/*;' 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ysi FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ysi FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 U&'end;' 1 CREATE FUNCTION yxfj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;b'; END 1;
| AS x --;;--;;;
|  FROM t WHERE a = CREATE FUNCTION on_r() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END AS x u_$j$f AS x;
| col CREATE FUNCTION xnec() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;select 1;END ;'; SELECT 1; END ,;1 CREATE FUNCTION _vsj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;; --;'; END col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'end;' 1 CREATE FUNCTION yxfj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;b'; END 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'end;' 1 CREATE FUNCTION yxfj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;b'; END 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x CREATE FUNCTION wvnp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b--;'; SELECT 1; END 1 ");(/*;" +;FROM t WHERE a = B'01;10' 1 fp$i$o AS x --select 1;);(
|  1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION wvnp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b--;'; SELECT 1; END 1 ");(/*;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION wvnp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b--;'; SELECT 1; END 1 ");(/*;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x X'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x X'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x X'01;10' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x X'01;10' AS x --$;$/*;
|  SELECT --$;$a;b
|  ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x X'01;10' AS x --$;$/*;
|  SELECT --$;$a;b
|  ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x X'01;10' AS x --$;$/*;
|  SELECT --$;$a;b
|  ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT --a;b
|  col U&';;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 42
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --a;b
|  col U&';;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 42
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --a;b
|  col U&';;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 42
> exec
| COPY kvw FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kvw FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kvw FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT "$x$;);(; " AS x e'\';\\' + '$x$;' col;
| + x'01;10' col;
< error
| 42601 syntax error at or near "e'\\';\\\\'"
| position: 25
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "$x$;);(; " AS x e'\';\\' + '$x$;' col;
< error
| 42601 syntax error at or near "e'\\';\\\\'"
| position: 25
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "$x$;);(; " AS x e'\';\\' + '$x$;' col;
< error
| 42601 syntax error at or near "e'\\';\\\\'"
| position: 25
> exec
| COPY asj FROM stdin;
| 1	\N
| a;b
| \.
|
| col /*);(*/ col kg$s$d +;
| + U&'$x$;END ;select 1;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY asj FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY asj FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , (SELECT 1; SELECT 2) col e'\';\\\';' col;
| SELECT (SELECT 1; SELECT 2) , $$END ;$$ AS x --end;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) col e'\';\\\';' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) col e'\';\\\';' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY bmp FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| COPY okw FROM stdin;
| a;b
| 1	\N
| \.
|
| SELECT oc$k$p AS x (SELECT 1; SELECT 2) SELECT e'\';' SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bmp FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bmp FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col '/*;END ;' +;
| col ";;" FROM t WHERE a = '/*;end;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '/*;END ;' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '/*;END ;' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT U&'END ;$x$;);(' SELECT; -- trailing note;
| COPY kf_ FROM stdin;
| a;b
| a;b
| \.
|
| , (SELECT 1; SELECT 2) , B'01;10' col (SELECT 1; SELECT 2) +;
| 1 ir$_$v , CREATE FUNCTION jxgc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "a"
| position: 72
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'END ;$x$;);(' SELECT;
< ok 1 -
> exec
|  -- trailing note;
| COPY kf_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 41
> exec
| ROLLBACK
< ok 0 -
> exec
|  -- trailing note;
| COPY kf_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 41
> exec
| col 'select 1;;;END ;' col;col $_r$);($_x$$_r$ FROM t WHERE a = /*/*end;--;$x$;*/*/ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'select 1;;;END ;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'select 1;;;END ;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY exm FROM stdin;
| a;b
| 1	\N
| \.
|  AS x '--;' + e''';\';\\' ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY exm FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY exm FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col $$--;; $$ FROM t WHERE a = --end;
|  +;
| , b_$n$_ 1; -- trailing note;
| + "end;" col _z$r$z SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$--;; $$ FROM t WHERE a = --end;
|  +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$--;; $$ FROM t WHERE a = --end;
|  +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY quz FROM stdin;
| 1	\N
| 1	\N
| \.
|  1 /*END ;*/ FROM t WHERE a = (SELECT 1; SELECT 2) AS x;col E'\\' 1; -- trailing note;
| AS x /*/*; */*/ AS x B'01;10' 1 "end;" col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY quz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY quz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , x'01;10' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , x'01;10' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , x'01;10' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x /*/*/*/ *;; end;*/*/*/ 1;
| FROM t WHERE a = b'01;10' ,;
| SELECT U&';;' 1 B'01;10' col;
| 1 'a;bselect 1;; ' col 'end;$;$;;\' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*/ *;; end;*/*/*/ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*/ *;; end;*/*/*/ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY iyn FROM stdin;
| 1	\N
| a;b
| \.
|
| , $ls$$;$;end;$lx$$ls$ col /*/*/*/ *;--;*/*/*/ 1 b'01;10' AS x;COPY uwk FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| + b'01;10' SELECT (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY iyn FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY iyn FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 "a;b$;$" AS x;
| 1 U&'$;$end;' FROM t WHERE a =;
| COPY sk_ FROM stdin;
| a;b
| a;b
| \.
|
| COPY roc FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "a;b$;$" AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "a;b$;$" AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) + U&'--;' SELECT CREATE FUNCTION mxif() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;a;b$x$;'; END +;
| AS x $mfi$$x$;$mfi$ SELECT; -- trailing note;
| col /*$;$; */ 1 B'01;10' AS x; -- trailing note;
| , ");(" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) + U&'--;' SELECT CREATE FUNCTION mxif() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;a;b$x$;'; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) + U&'--;' SELECT CREATE FUNCTION mxif() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;a;b$x$;'; END +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , zk$m$i col ko$s$g ,;col xg$m$b +; COPY bso FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , zk$m$i col ko$s$g ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , zk$m$i col ko$s$g ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , '--;' + CREATE FUNCTION xsqy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; ; '; END 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , '--;' + CREATE FUNCTION xsqy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; ; '; END 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , '--;' + CREATE FUNCTION xsqy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; ; '; END 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , CREATE FUNCTION jnqi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END FROM t WHERE a = U&'; ' , B'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION jnqi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END FROM t WHERE a = U&'; ' , B'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION jnqi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END FROM t WHERE a = U&'; ' , B'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY _fn FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _fn FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _fn FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY kyn FROM stdin;
| c;d;e
| 1	\N
| \.
|
| col i_$v$s FROM t WHERE a = --select 1;; 
|  AS x /*END ;a;b*/ SELECT;
| COPY tjt FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY qkb FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kyn FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kyn FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x /*/*/*;;a;b*/*/*/ AS x ";$x$;" AS x "END ;" FROM t WHERE a =;
| COPY pgf FROM stdin;
| c;d;e
| 1	\N
| \.
|
| SELECT U&'a;b' col CREATE FUNCTION njbx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END FROM t WHERE a =; + "/*;" 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*;;a;b*/*/*/ AS x ";$x$;" AS x "END ;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*;;a;b*/*/*/ AS x ";$x$;" AS x "END ;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x CREATE FUNCTION hesn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';--;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION hesn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';--;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION hesn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';--;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY wlk FROM stdin;
| 1	\N
| c;d;e
| \.
|  COPY hqi FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| + --;;;
|  + U&'--;' AS x;
| COPY tqk FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wlk FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wlk FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY dz_ FROM stdin;
| c;d;e
| 1	\N
| \.
| col --);(--;
|  SELECT CREATE FUNCTION dqlp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$END ;'; END SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dz_ FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dz_ FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col $$END ;;;$$ , --$;$a;b$;$
|  1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$END ;;;$$ , --$;$a;b$;$
|  1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$END ;;;$$ , --$;$a;b$;$
|  1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT U&';/*;END ;' col;
| 1 $$--;$$ + CREATE FUNCTION trhx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(/*;'; END +; -- trailing note;
| , --$;$);(end;
|  AS x mu$g$q col;
< error
| 42601 syntax error at or near "1"
| position: 27
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&';/*;END ;' col;
< ok 1 -
> exec
|
| 1 $$--;$$ + CREATE FUNCTION trhx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(/*;'; END +;
< error
| 42601 syntax error at or near "1"
| position: 2
> exec
| ROLLBACK
< ok 0 -
> exec
|
| 1 $$--;$$ + CREATE FUNCTION trhx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(/*;'; END +;
< error
| 42601 syntax error at or near "1"
| position: 2
> exec
| COPY smk FROM stdin;
| a;b
| a;b
| \.
|
| FROM t WHERE a = "select 1;;END ;" 1;
| COPY tel FROM stdin;
| 1	\N
| c;d;e
| \.
|
| COPY zxp FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY smk FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY smk FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 '; ; end;\' FROM t WHERE a = CREATE FUNCTION zbtl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;select 1;END ;'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '; ; end;\' FROM t WHERE a = CREATE FUNCTION zbtl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;select 1;END ;'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '; ; end;\' FROM t WHERE a = CREATE FUNCTION zbtl() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;select 1;END ;'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x 'select 1;' SELECT $hcrb$/*;);(end;$hcrx$$hcrb$ AS x CREATE FUNCTION furu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x 'select 1;' SELECT $hcrb$/*;);(end;$hcrx$$hcrb$ AS x CREATE FUNCTION furu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x 'select 1;' SELECT $hcrb$/*;);(end;$hcrx$$hcrb$ AS x CREATE FUNCTION furu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; END SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col U&'; ' + $x$end;$x$ + $fy$select 1;$x$;;;$fx$$fy$ +;
| AS x --END ;$;$
|  , "a;b;;$x$;" AS x; SELECT U&'/*;);(/*;' AS x hg$c$_ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'; ' + $x$end;$x$ + $fy$select 1;$x$;;;$fx$$fy$ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'; ' + $x$end;$x$ + $fy$select 1;$x$;;;$fx$$fy$ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col yz$d$v ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col yz$d$v ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col yz$d$v ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT (SELECT 1; SELECT 2) +;
| , e''';' col;
| SELECT (SELECT 1; SELECT 2) col ";;END ;--;" col --select 1;
|  col; , "a;b/*;select 1;" 1 'select 1;\' col "; ;" FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| AS x yp$p$y SELECT CREATE FUNCTION hcop() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x yp$p$y SELECT CREATE FUNCTION hcop() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x yp$p$y SELECT CREATE FUNCTION hcop() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY cdt FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| FROM t WHERE a = "select 1;" , ');(/*;' + (SELECT 1; SELECT 2) AS x; FROM t WHERE a = fj$n$i + /*a;bEND ;;*/ col qf$v$u col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cdt FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cdt FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
| COPY ugb FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY glc FROM stdin;
| c;d;e
| a;b
| \.
|
| FROM t WHERE a = CREATE FUNCTION nliz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;'; END AS x E'\\'';' + "; " ,;
| , $$;select 1;end;$$ FROM t WHERE a = $$);();($$ FROM t WHERE a = B'01;10' AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY glc FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY glc FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , e'\\' SELECT U&'a;b/*;/*;' ,;SELECT '--;' + (SELECT 1; SELECT 2) AS x;
| + $i$a;b$i$ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\\' SELECT U&'a;b/*;/*;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\\' SELECT U&'a;b/*;/*;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT '; --;; ' col U&');(--;' + /*/*/*$x$;;a;b*/*/*/ col;
| COPY x_f FROM stdin;
| 1	\N
| a;b
| \.
|  + CREATE FUNCTION witq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END AS x b'01;10' SELECT;
| FROM t WHERE a = "--;" col;
< error
| 42601 syntax error at or near "U&');(--;'"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT '; --;; ' col U&');(--;' + /*/*/*$x$;;a;b*/*/*/ col;
< error
| 42601 syntax error at or near "U&');(--;'"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT '; --;; ' col U&');(--;' + /*/*/*$x$;;a;b*/*/*/ col;
< error
| 42601 syntax error at or near "U&');(--;'"
| position: 22
> exec
| COPY wka FROM stdin;
| 1	\N
| c;d;e
| \.
|
| , CREATE FUNCTION yufc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b$;$'; SELECT 1; END , E''';' 1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wka FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wka FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 '$;$);(' col B'01;10' FROM t WHERE a =; -- trailing note;
| col (SELECT 1; SELECT 2) , aj$q$w 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '$;$);(' col B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '$;$);(' col B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY hby FROM stdin;
| c;d;e
| 1	\N
| \.
| SELECT U&'select 1;/*;$x$;' SELECT x'01;10' +;+ "select 1;$x$;" col ';select 1;$;$' FROM t WHERE a = $nag$--;END ;$nag$ col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hby FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hby FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col '--;end;' ,;
| COPY cxc FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY _mi FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col '--;end;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col '--;end;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY nhj FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY khd FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY nhj FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY nhj FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x --select 1;end;
|  SELECT;
| col --$;$
|  ,;+ e''';\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --select 1;end;
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --select 1;end;
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 --END ;
|  SELECT U&'a;b' AS x $a$$x$;$x$$a$ +; -- trailing note;
| 1 --END ;);(select 1;
|  1; -- trailing note;
| AS x CREATE FUNCTION dxju() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$;;a;b'; END + '$;$\' SELECT ';;);(END ;' +;
| + E''';\\\\' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --END ;
|  SELECT U&'a;b' AS x $a$$x$;$x$$a$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --END ;
|  SELECT U&'a;b' AS x $a$$x$;$x$$a$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY bu_ FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| col --end;
|  +; -- trailing note;
| + ";" , $vq$END ;a;b;$vq$ col; AS x (SELECT 1; SELECT 2) col U&'a;bselect 1;' ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bu_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bu_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , b'01;10' SELECT "END ;select 1;end;" FROM t WHERE a = vo$m$v 1;+ nz$l$_ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , b'01;10' SELECT "END ;select 1;end;" FROM t WHERE a = vo$m$v 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , b'01;10' SELECT "END ;select 1;end;" FROM t WHERE a = vo$m$v 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 "end;" col e'\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "end;" col e'\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "end;" col e'\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + "end;$;$/*;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "end;$;$/*;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "end;$;$/*;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col "$;$select 1;;;" 1 $hzbp$/*;$hzbp$ col; -- trailing note;
| FROM t WHERE a = /*/*END ;*/*/ SELECT;
| COPY mpx FROM stdin;
| 1	\N
| 1	\N
| \.
| 1 'END ;\' AS x 'select 1;' AS x $i_mh$;$i_mx$$i_mh$ 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "$;$select 1;;;" 1 $hzbp$/*;$hzbp$ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "$;$select 1;;;" 1 $hzbp$/*;$hzbp$ col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col (SELECT 1; SELECT 2) AS x B'01;10' FROM t WHERE a = gp$n$g AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) AS x B'01;10' FROM t WHERE a = gp$n$g AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) AS x B'01;10' FROM t WHERE a = gp$n$g AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x jy$f$x AS x; -- trailing note;
| AS x --$;$--;/*;
|  + (SELECT 1; SELECT 2) , U&'/*;' FROM t WHERE a =; COPY bui FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x jy$f$x AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x jy$f$x AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + "a;b" SELECT b'01;10' + 'a;b' SELECT;
| + (SELECT 1; SELECT 2) 1 CREATE FUNCTION ldna() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; END FROM t WHERE a =;col lt$g$z SELECT CREATE FUNCTION sacd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$; '; END col '$;$--;$;$\' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "a;b" SELECT b'01;10' + 'a;b' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "a;b" SELECT b'01;10' + 'a;b' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + b'01;10' , /*/*/*/ *;;a;b*/*/*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + b'01;10' , /*/*/*/ *;;a;b*/*/*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + b'01;10' , /*/*/*/ *;;a;b*/*/*/ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = U&');(' col (SELECT 1; SELECT 2) SELECT; SELECT ----;
|  , zy$f$m FROM t WHERE a = CREATE FUNCTION trky() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$;$'; SELECT 1; END +;
| SELECT --select 1;);(;;
|  FROM t WHERE a = $n$);(;$;$$x$$n$ 1 ');();(\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&');(' col (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&');(' col (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + e''';' AS x; COPY _ul FROM stdin;
| c;d;e
| c;d;e
| \.
|  FROM t WHERE a = x'01;10' 1; COPY epl FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + e''';' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + e''';' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = $$$x$;$$ , $kq$; ;; $kx$$kq$ 1 ';$;$--;\' +;
| SELECT $av_l$--;END ;$av_l$ AS x '--;;;\' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$$x$;$$ , $kq$; ;; $kx$$kq$ 1 ';$;$--;\' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$$x$;$$ , $kq$; ;; $kx$$kq$ 1 ';$;$--;\' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY won FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| FROM t WHERE a = 'a;b$;$);(\' SELECT (SELECT 1; SELECT 2) +;
| + --$x$;
|  , CREATE FUNCTION osuf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($x$;'; END col /*/*end;end;*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY won FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY won FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 CREATE FUNCTION kdab() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;'; SELECT 1; END FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
| COPY qhx FROM stdin;
| c;d;e
| a;b
| \.
|  -- trailing note;
| FROM t WHERE a = /*/*select 1;; $;$*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION kdab() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;'; SELECT 1; END FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION kdab() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;'; SELECT 1; END FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = b'01;10' 1 "a;b" SELECT (SELECT 1; SELECT 2) FROM t WHERE a =;
| COPY xhy FROM stdin;
| c;d;e
| a;b
| \.
|  -- trailing note;
| , (SELECT 1; SELECT 2) + ');($x$;; ' col ';;$;$' FROM t WHERE a =;
| COPY _gb FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = b'01;10' 1 "a;b" SELECT (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = b'01;10' 1 "a;b" SELECT (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + '; END ;\' FROM t WHERE a = CREATE FUNCTION mogg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;'; END +;
| COPY wfw FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + '; END ;\' FROM t WHERE a = CREATE FUNCTION mogg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;'; END +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + '; END ;\' FROM t WHERE a = CREATE FUNCTION mogg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;'; END +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col "/*;/*;$x$;" + U&'; ; select 1;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "/*;/*;$x$;" + U&'; ; select 1;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "/*;/*;$x$;" + U&'; ; select 1;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + ----;a;b);(
|  ,; SELECT b'01;10' AS x;
| FROM t WHERE a = $$$x$;$$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ----;a;b);(
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ----;a;b);(
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col (SELECT 1; SELECT 2) col b'01;10' ,; -- trailing note;
| SELECT "end;" SELECT;
| FROM t WHERE a = E'\';'';\';' , ----;
|  AS x;COPY ptj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) col b'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) col b'01;10' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , /*;;--;*/ SELECT --end;$x$;--;
|  col;
| , "--;$;$--;" + ';$x$;;;' , ";;);(" SELECT; -- trailing note;
| COPY vwn FROM stdin;
| a;b
| c;d;e
| \.
| SELECT X'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*;;--;*/ SELECT --end;$x$;--;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*;;--;*/ SELECT --end;$x$;--;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , rg$g$r , --/*;
|  1;+ --; select 1;
|  + dg$k$a + E'\';\\\\' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , rg$g$r , --/*;
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , rg$g$r , --/*;
|  1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col $$end;--;end;$$ +;
| SELECT CREATE FUNCTION vtwo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;end;$;$'; END SELECT;
| AS x CREATE FUNCTION gdiz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(/*;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$end;--;end;$$ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$end;--;end;$$ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col xo$t$k AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col xo$t$k AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col xo$t$k AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x U&'END ;' 1 --END ;end;
|  AS x (SELECT 1; SELECT 2) col;
| FROM t WHERE a = U&'; END ;$;$' + --a;bEND ;);(
|  , (SELECT 1; SELECT 2) 1; AS x U&'$x$;; ;;' SELECT /*/*select 1;;;*/*/ , --end;;;
|  SELECT;
| col U&';select 1;select 1;' AS x CREATE FUNCTION t_ux() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END AS x /*/*END ;$;$);(*/*/ col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'END ;' 1 --END ;end;
|  AS x (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'END ;' 1 --END ;end;
|  AS x (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY csl FROM stdin;
| a;b
| c;d;e
| \.
|
| FROM t WHERE a = yr$l$t col (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY csl FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY csl FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , U&'select 1;END ;end;' col go$e$q 1 gj$m$c +;
| + (SELECT 1; SELECT 2) col '/*;END ;$x$;\' AS x;
| FROM t WHERE a = /*/*/*a;b*/*/*/ 1 (SELECT 1; SELECT 2) col X'01;10' +;
| + /*/*;;*/*/ AS x $qd$);(;$qd$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'select 1;END ;end;' col go$e$q 1 gj$m$c +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'select 1;END ;end;' col go$e$q 1 gj$m$c +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x X'01;10' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x X'01;10' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x X'01;10' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + /*/*; END ;;;*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*; END ;;;*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*; END ;;;*/*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , CREATE FUNCTION bw_d() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;);(; '; SELECT 1; END SELECT E'\n;'';' SELECT CREATE FUNCTION vtdi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END 1;
| AS x E'\\'';' col (SELECT 1; SELECT 2) + U&');($;$;;' +;
| 1 /*/*/*;/ *;*/*/*/ + (SELECT 1; SELECT 2) + ';\' SELECT;
| + "select 1;a;b);(" col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION bw_d() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;);(; '; SELECT 1; END SELECT E'\n;'';' SELECT CREATE FUNCTION vtdi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION bw_d() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;);(; '; SELECT 1; END SELECT E'\n;'';' SELECT CREATE FUNCTION vtdi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT x'01;10' AS x (SELECT 1; SELECT 2) col e'\';' AS x;
< error
| 42601 syntax error at or near "("
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' AS x (SELECT 1; SELECT 2) col e'\';' AS x;
< error
| 42601 syntax error at or near "("
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT x'01;10' AS x (SELECT 1; SELECT 2) col e'\';' AS x;
< error
| 42601 syntax error at or near "("
| position: 22
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = U&'; );(' col E''';\\' +;
| COPY drw FROM stdin;
| a;b
| 1	\N
| \.
|
| FROM t WHERE a = U&'$x$;; ' AS x /*a;b);(;*/ col ";/*;" AS x;
| FROM t WHERE a = '/*;\' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = U&'; );(' col E''';\\' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = U&'; );(' col E''';\\' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x $yj$$x$;$yj$ +;
| + (SELECT 1; SELECT 2) AS x pw$h$d , x'01;10' ,; -- trailing note;
| , $$/*;a;bend;$$ col;col B'01;10' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $yj$$x$;$yj$ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $yj$$x$;$yj$ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x CREATE FUNCTION ujpf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;--;'; SELECT 1; END + --; select 1;;
|  col;
| AS x ";$;$" 1; -- trailing note;
| SELECT 'END ;' FROM t WHERE a = E''';' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION ujpf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;--;'; SELECT 1; END + --; select 1;;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION ujpf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;--;'; SELECT 1; END + --; select 1;;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + $cqye$select 1;a;b/*;$cqyx$$cqye$ 1;
| 1 ak$b$l , U&'$;$select 1;' SELECT; -- trailing note;
| COPY zjl FROM stdin;
| a;b
| 1	\N
| \.
|
| col --select 1;
|  1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $cqye$select 1;a;b/*;$cqyx$$cqye$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $cqye$select 1;a;b/*;$cqyx$$cqye$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) col tl$i$q + "; " 1; , (SELECT 1; SELECT 2) , b'01;10' FROM t WHERE a = U&';select 1;end;' FROM t WHERE a =;
| FROM t WHERE a = CREATE FUNCTION j_cy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b);('; END SELECT '--;END ;end;' + CREATE FUNCTION calk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;;'; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) col tl$i$q + "; " 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) col tl$i$q + "; " 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY zgp FROM stdin;
| c;d;e
| a;b
| \.
|  -- trailing note;
| COPY byg FROM stdin;
| a;b
| c;d;e
| \.
|
| SELECT ';;/*;select 1;\' 1;
| FROM t WHERE a = $umk$;;$umx$$umk$ col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zgp FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zgp FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 --/*;
|  , --/*;);(END ;
|  + U&'$;$$;$' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --/*;
|  , --/*;);(END ;
|  + U&'$;$$;$' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --/*;
|  , --/*;);(END ;
|  + U&'$;$$;$' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col X'01;10' , (SELECT 1; SELECT 2) 1 CREATE FUNCTION qbig() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END SELECT;
| AS x cm$n$b AS x CREATE FUNCTION mdca() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END , ";$x$;" ,; -- trailing note;
| col 'a;b' AS x U&'select 1;' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col X'01;10' , (SELECT 1; SELECT 2) 1 CREATE FUNCTION qbig() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col X'01;10' , (SELECT 1; SELECT 2) 1 CREATE FUNCTION qbig() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY wbm FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wbm FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wbm FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + --;;
|  FROM t WHERE a = $pvv_$;$pvv_$ FROM t WHERE a = U&'--;select 1;$x$;' 1; -- trailing note;
| AS x (SELECT 1; SELECT 2) 1 B'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --;;
|  FROM t WHERE a = $pvv_$;$pvv_$ FROM t WHERE a = U&'--;select 1;$x$;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --;;
|  FROM t WHERE a = $pvv_$;$pvv_$ FROM t WHERE a = U&'--;select 1;$x$;' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + U&'select 1;);();(' + (SELECT 1; SELECT 2) AS x CREATE FUNCTION dkdn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$; ;'; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'select 1;);();(' + (SELECT 1; SELECT 2) AS x CREATE FUNCTION dkdn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$; ;'; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'select 1;);();(' + (SELECT 1; SELECT 2) AS x CREATE FUNCTION dkdn() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$; ;'; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x /*;*/ , --);(
|  1;
| FROM t WHERE a = ----;
|  1 X'01;10' + "end;END ;--;" AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*;*/ , --);(
|  1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*;*/ , --);(
|  1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY der FROM stdin;
| 1	\N
| 1	\N
| \.
|
| , U&';;;;' AS x; -- trailing note;
| col $a$; --;$x$;$a$ FROM t WHERE a = _m$_$e col --);(/*;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY der FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY der FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = --/*;--;end;
|  ,; -- trailing note;
| COPY opm FROM stdin;
| 1	\N
| 1	\N
| \.
|
| COPY fgf FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --/*;--;end;
|  ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --/*;--;end;
|  ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col e'\';'';' ,;
| COPY pob FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e'\';'';' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e'\';'';' ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY dip FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| , $gd$);(select 1;$gx$$gd$ FROM t WHERE a = CREATE FUNCTION xagv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;/*;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dip FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dip FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 E'\\\';' , --END ;;
|  + rx$f$f +;
| col E'\\\\\';' +; , B'01;10' col e'\\' SELECT E'\\\';' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 E'\\\';' , --END ;;
|  + rx$f$f +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 E'\\\';' , --END ;;
|  + rx$f$f +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) 1 dz$p$u SELECT '$x$;select 1;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) 1 dz$p$u SELECT '$x$;select 1;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) 1 dz$p$u SELECT '$x$;select 1;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col e''';' AS x u_$d$j SELECT bw$p$h +;AS x x'01;10' AS x E''';\\\n;' + $$select 1;$$ 1;
| + --;;/*;
|  , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e''';' AS x u_$d$j SELECT bw$p$h +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e''';' AS x u_$d$j SELECT bw$p$h +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY rie FROM stdin;
| 1	\N
| c;d;e
| \.
| COPY uje FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY rie FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY rie FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x b'01;10' SELECT $dfiq$$;$$dfix$$dfiq$ FROM t WHERE a = --$x$;--;
|  col; + ----;
|  col qh$t$g AS x e''';'';'';' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x b'01;10' SELECT $dfiq$$;$$dfix$$dfiq$ FROM t WHERE a = --$x$;--;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x b'01;10' SELECT $dfiq$$;$$dfix$$dfiq$ FROM t WHERE a = --$x$;--;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY gps FROM stdin;
| c;d;e
| 1	\N
| \.
|
| AS x /*$;$*/ 1 $$; ;;$$ ,;
| FROM t WHERE a = (SELECT 1; SELECT 2) 1 /*$;$$;$*/ +;
| AS x '; --;$x$;' , /*$;$end;*/ SELECT x'01;10' SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gps FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gps FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 /*/*a;b/ *;*/*/ 1;
| col /*/*;;select 1;;*/*/ 1 'a;b/*;END ;' col e'\';\n;'';' col;+ E'\';\';' , B'01;10' 1 --a;b
|  FROM t WHERE a =;FROM t WHERE a = /*/*/*a;b*/*/*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*a;b/ *;*/*/ 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*a;b/ *;*/*/ 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + CREATE FUNCTION yldv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';$;$'; END 1 /*$;$--;*/ FROM t WHERE a = x'01;10' 1;
| COPY zqf FROM stdin;
| a;b
| a;b
| \.
|
| COPY waf FROM stdin;
| a;b
| a;b
| \.
|
| SELECT /*/*/*a;ba;b;;*/*/*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION yldv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';$;$'; END 1 /*$;$--;*/ FROM t WHERE a = x'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION yldv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';$;$'; END 1 /*$;$--;*/ FROM t WHERE a = x'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT yf$a$a SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "("
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT yf$a$a SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "("
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT yf$a$a SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "("
| position: 22
> exec
| AS x e''';\\' FROM t WHERE a =;
| 1 (SELECT 1; SELECT 2) +; -- trailing note;
| 1 X'01;10' col --/*;
|  1 B'01;10' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x e''';\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x e''';\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + /*/*$;$*/*/ SELECT (SELECT 1; SELECT 2) FROM t WHERE a = "$;$);(" col;
| AS x --end;end;a;b
|  col /*/*END ;*/*/ SELECT /*$x$;*/ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*$;$*/*/ SELECT (SELECT 1; SELECT 2) FROM t WHERE a = "$;$);(" col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*$;$*/*/ SELECT (SELECT 1; SELECT 2) FROM t WHERE a = "$;$);(" col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 U&'; ; select 1;' FROM t WHERE a = CREATE FUNCTION qfkx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$'; END col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'; ; select 1;' FROM t WHERE a = CREATE FUNCTION qfkx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$'; END col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'; ; select 1;' FROM t WHERE a = CREATE FUNCTION qfkx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$;$'; END col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY dct FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dct FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dct FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY hvq FROM stdin;
| 1	\N
| a;b
| \.
|
| 1 --$x$;--;
|  AS x;
| + e'\n;\';'';' + --; $;$end;
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hvq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hvq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col (SELECT 1; SELECT 2) + CREATE FUNCTION oqrt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$'; SELECT 1; END SELECT;
| FROM t WHERE a = $vopp$$x$;end;; $vopp$ SELECT; COPY y_a FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) + CREATE FUNCTION oqrt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) + CREATE FUNCTION oqrt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x --;END ;end;
|  FROM t WHERE a = U&'$;$/*;select 1;' , rp$v$g AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --;END ;end;
|  FROM t WHERE a = U&'$;$/*;select 1;' , rp$v$g AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --;END ;end;
|  FROM t WHERE a = U&'$;$/*;select 1;' , rp$v$g AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 B'01;10' AS x; AS x $h_z$);(/*;$h_x$$h_z$ , X'01;10' AS x ';;; ' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 B'01;10' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 B'01;10' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + b'01;10' FROM t WHERE a =;
| + "$;$" FROM t WHERE a = e'\\'';' ,;
| , 'select 1;' , /*/*/*/ *;;;*/*/*/ , (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col E'\\'';\';' , (SELECT 1; SELECT 2) SELECT U&'$;$' +; 1 B'01;10' , CREATE FUNCTION _enb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END FROM t WHERE a =;
| , --select 1;a;b
|  FROM t WHERE a = "/*;" AS x;, (SELECT 1; SELECT 2) FROM t WHERE a = X'01;10' , nb$m$_ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col E'\\'';\';' , (SELECT 1; SELECT 2) SELECT U&'$;$' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col E'\\'';\';' , (SELECT 1; SELECT 2) SELECT U&'$;$' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 '; a;bend;' +;
| COPY jsu FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| , '$x$;;' col x'01;10' + X'01;10' ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '; a;bend;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '; a;bend;' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT X'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT X'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT X'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 17
> exec
| COPY vto FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY vto FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY vto FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = U&'a;b' col kx$b$s SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'a;b' col kx$b$s SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'a;b' col kx$b$s SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY kjr FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kjr FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kjr FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + --$x$;
|  col;
| COPY qpv FROM stdin;
| a;b
| c;d;e
| \.
| col 'end;/*;END ;' FROM t WHERE a =;
| col --);(select 1;end;
|  col /*/*a;b*/*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --$x$;
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --$x$;
|  col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col fi$p$e SELECT $$$x$;end;$;$$$ ,;
| , U&';' col CREATE FUNCTION jwae() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bEND ;);('; END col;1 e'\\\\' 1 --; 
|  1 'end;a;b; ' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col fi$p$e SELECT $$$x$;end;$;$$$ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col fi$p$e SELECT $$$x$;end;$;$$$ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY wzr FROM stdin;
| 1	\N
| a;b
| \.
|
| + CREATE FUNCTION _flw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;end;end;'; SELECT 1; END +;
| col CREATE FUNCTION qrev() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$x$;;;'; SELECT 1; END FROM t WHERE a = /*/*end;;END ;*/*/ SELECT E'\n;'';' +;
| COPY cqi FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wzr FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wzr FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY gen FROM stdin;
| c;d;e
| 1	\N
| \.
|
| AS x U&');(END ;' 1 U&';;; ' FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gen FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gen FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = ik$v$g AS x U&';;' col U&'$;$' FROM t WHERE a =;FROM t WHERE a = --;;;;
|  FROM t WHERE a = lm$c$g ,;
| COPY mnh FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ik$v$g AS x U&';;' col U&'$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ik$v$g AS x U&';;' col U&'$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY grb FROM stdin;
| c;d;e
| a;b
| \.
|  AS x ');(\' AS x ';/*;;' 1 U&'$;$' FROM t WHERE a =;COPY sxe FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY grb FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY grb FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x --;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + --;/*;; 
|  + $$;$x$;$x$;$$ 1 x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --;/*;; 
|  + $$;$x$;$x$;$$ 1 x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --;/*;; 
|  + $$;$x$;$x$;$$ 1 x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 'end;\' 1 U&'select 1;/*;select 1;' col U&'end;);(; ' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 'end;\' 1 U&'select 1;/*;select 1;' col U&'end;);(; ' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 'end;\' 1 U&'select 1;/*;select 1;' col U&'end;);(; ' +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) SELECT;
| COPY mqc FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 gj$t$z col;
| FROM t WHERE a = U&'end;$;$--;' SELECT U&'--;;' SELECT bz$c$f ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY zzr FROM stdin;
| c;d;e
| a;b
| \.
| , CREATE FUNCTION vmgd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;end;'; END AS x (SELECT 1; SELECT 2) col $clbv$/*;$;$; $clbv$ SELECT;
| FROM t WHERE a = et$_$l col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zzr FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zzr FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY fzg FROM stdin;
| c;d;e
| a;b
| \.
|
| SELECT lh$a$g col $vwqb$/*;$vwqb$ 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fzg FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fzg FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x "END ;" SELECT "select 1;" +;1 ";--;;;" AS x b'01;10' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "END ;" SELECT "select 1;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "END ;" SELECT "select 1;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY ial FROM stdin;
| c;d;e
| c;d;e
| \.
|  + E'\n;'';' FROM t WHERE a = '$x$;' FROM t WHERE a =;SELECT '--;END ;;' 1 $rq$$;$);(;;$rx$$rq$ SELECT b'01;10' AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ial FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ial FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x --end;
|  AS x; , ";;;;; " AS x;
| + '/*;a;b' , U&'; ' col;
| FROM t WHERE a = --/*;
|  SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --end;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --end;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + (SELECT 1; SELECT 2) AS x;
| COPY mhb FROM stdin;
| c;d;e
| a;b
| \.
|
| 1 U&'end;' FROM t WHERE a = "; $x$;" col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x U&';;end;' col km$r$x + --END ;$;$
|  1; COPY zgw FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| 1 --;;
|  col;col b'01;10' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&';;end;' col km$r$x + --END ;$;$
|  1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&';;end;' col km$r$x + --END ;$;$
|  1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col "a;bend;$;$" col;
| 1 'select 1;/*;$x$;' FROM t WHERE a = '/*;\' FROM t WHERE a = '$x$;;;\' 1; FROM t WHERE a = b'01;10' + X'01;10' +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "a;bend;$;$" col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "a;bend;$;$" col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + U&'$;$' , U&'; ;;' + CREATE FUNCTION zklu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;END ;'; END SELECT;
| SELECT $$a;b$;$a;b$$ AS x xb$c$i col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'$;$' , U&'; ;;' + CREATE FUNCTION zklu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;END ;'; END SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'$;$' , U&'; ;;' + CREATE FUNCTION zklu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;END ;'; END SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + CREATE FUNCTION wdut() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END AS x;
| + '--;);(a;b\' col;
| 1 '$x$;--;$;$' AS x CREATE FUNCTION awey() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION wdut() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION wdut() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col U&'$;$;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'$;$;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'$;$;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY gej FROM stdin;
| 1	\N
| 1	\N
| \.
|
| AS x '$x$;' ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gej FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gej FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 mq$s$z FROM t WHERE a = '--;;;select 1;' SELECT sx$t$s FROM t WHERE a =; -- trailing note;
| , ');(\' + x'01;10' +;
| , X'01;10' + X'01;10' FROM t WHERE a = "/*;END ;" 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 mq$s$z FROM t WHERE a = '--;;;select 1;' SELECT sx$t$s FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 mq$s$z FROM t WHERE a = '--;;;select 1;' SELECT sx$t$s FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , --$x$;
|  col;
| 1 /*/*/*;;;select 1;*/*/*/ SELECT bf$s$g AS x; COPY eqq FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --$x$;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --$x$;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x 'end;' + (SELECT 1; SELECT 2) ,;COPY toh FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x 'end;' + (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x 'end;' + (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY e_g FROM stdin;
| c;d;e
| 1	\N
| \.
|
| col ";" SELECT nq$i$s col;
| COPY ufx FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY e_g FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY e_g FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + X'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , e'\n;'';\';' SELECT (SELECT 1; SELECT 2) 1;
| , /*/*/ *;END ;*/*/ col "--;; " + b'01;10' 1;
| col "end;/*;" , --/*;
|  SELECT "; ; $;$" col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\n;'';\';' SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\n;'';\';' SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY pg_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pg_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pg_ FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 U&'--;;;--;' , U&'a;b; end;' SELECT;
| COPY lpl FROM stdin;
| c;d;e
| a;b
| \.
|  -- trailing note;
| 1 re$a$e , CREATE FUNCTION uh_e() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END AS x B'01;10' 1;
| SELECT E'\n;\\\';' SELECT E''';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'--;;;--;' , U&'a;b; end;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'--;;;--;' , U&'a;b; end;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x js$a$f ,;
| 1 U&'$;$/*;' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x js$a$f ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x js$a$f ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x "END ;);(" , X'01;10' AS x CREATE FUNCTION tleg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';a;b'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "END ;);(" , X'01;10' AS x CREATE FUNCTION tleg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';a;b'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "END ;);(" , X'01;10' AS x CREATE FUNCTION tleg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';a;b'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY zti FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zti FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zti FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 /*/*/*select 1;*/*/*/ FROM t WHERE a = U&'--;' col; FROM t WHERE a = --$;$
|  +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*select 1;*/*/*/ FROM t WHERE a = U&'--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*select 1;*/*/*/ FROM t WHERE a = U&'--;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col e'\\\\'';' SELECT (SELECT 1; SELECT 2) 1;
| COPY maa FROM stdin;
| c;d;e
| 1	\N
| \.
|
| + X'01;10' 1 /*/*END ;;;*/*/ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e'\\\\'';' SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e'\\\\'';' SELECT (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY xhz FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xhz FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xhz FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = '; ' , --;;
|  +;AS x '/*;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = '; ' , --;;
|  +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = '; ' , --;;
|  +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT du$t$m AS x CREATE FUNCTION kvai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;a;b'; END AS x ';;$;$' SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 20
> exec
| BEGIN
< ok 0 -
> exec
| SELECT du$t$m AS x CREATE FUNCTION kvai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;a;b'; END AS x ';;$;$' SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 20
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT du$t$m AS x CREATE FUNCTION kvai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;a;b'; END AS x ';;$;$' SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 20
> exec
| FROM t WHERE a = $u$; $;$$x$$u$ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $u$; $;$$x$$u$ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $u$; $;$$x$$u$ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x U&'select 1;$;$; ' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'select 1;$;$; ' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'select 1;$;$; ' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col --a;b;
|  SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --a;b;
|  SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --a;b;
|  SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT $st_b$a;b; $st_b$ SELECT;
< ok 1 -
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $st_b$a;b; $st_b$ SELECT;
< ok 1 -
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*select 1;END ;*/ 1 U&');(END ;' + (SELECT 1; SELECT 2) 1; -- trailing note;
| 1 U&';;end;' 1; + "a;b" AS x; -- trailing note;
| SELECT /*/*/*; end;;;*/*/*/ SELECT E'\\\';'';' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*select 1;END ;*/ 1 U&');(END ;' + (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*select 1;END ;*/ 1 U&');(END ;' + (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY fwz FROM stdin;
| 1	\N
| 1	\N
| \.
|
| FROM t WHERE a = wm$k$a ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fwz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fwz FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = E'\n;' SELECT;
| FROM t WHERE a = b'01;10' FROM t WHERE a = (SELECT 1; SELECT 2) +;
| AS x E''';' FROM t WHERE a =;
| + CREATE FUNCTION ioen() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;; $;$'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E'\n;' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\n;' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 $$$;$end;a;b$$ AS x U&'END ;--;' 1 "a;bend;$x$;" 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $$$;$end;a;b$$ AS x U&'END ;--;' 1 "a;bend;$x$;" 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $$$;$end;a;b$$ AS x U&'END ;--;' 1 "a;bend;$x$;" 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col (SELECT 1; SELECT 2) + e'\n;' col;
| FROM t WHERE a = CREATE FUNCTION cqtb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;a;bEND ;'; END + $hv$$x$;$hv$ +;SELECT xr$a$c FROM t WHERE a = e'\n;'';' +;AS x /*/*);(--;*/*/ , /*/*$x$;;/ *;*/*/ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) + e'\n;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) + e'\n;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , $$$;$$$ col;
| , $$; );($$ col $$a;bend;a;b$$ SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $$$;$$$ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $$$;$$$ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x --$;$
|  SELECT; -- trailing note;
| AS x B'01;10' 1;
| col (SELECT 1; SELECT 2) + --);(/*;select 1;
|  1;
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --$;$
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --$;$
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY mow FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mow FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mow FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + ve$x$w SELECT;
| + "end;;;" col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ve$x$w SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ve$x$w SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col --$;$;;;
|  SELECT --$;$select 1;
|  FROM t WHERE a =;
| , U&'end;;select 1;' col e'\n;\n;\n;' SELECT $nbta$; $nbtx$$nbta$ FROM t WHERE a =;
| AS x '$;$end;\' FROM t WHERE a =;
| COPY o_d FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --$;$;;;
|  SELECT --$;$select 1;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --$;$;;;
|  SELECT --$;$select 1;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY see FROM stdin;
| 1	\N
| 1	\N
| \.
|
| AS x "END ;" , ";" SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY see FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY see FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 ';;' 1;
| col /*/*a;bend;*/*/ +; col sb$u$v FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a =;COPY cqj FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ';;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ';;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY rzq FROM stdin;
| 1	\N
| c;d;e
| \.
|  + 'select 1;; --;' AS x U&'select 1;' AS x --; 
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY rzq FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY rzq FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col E''';\\\';' , CREATE FUNCTION _vzr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END + $$;;; $$ SELECT;
| COPY sow FROM stdin;
| 1	\N
| 1	\N
| \.
|
| AS x x'01;10' AS x lq$j$w , iv$q$a ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col E''';\\\';' , CREATE FUNCTION _vzr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END + $$;;; $$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col E''';\\\';' , CREATE FUNCTION _vzr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; SELECT 1; END + $$;;; $$ SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x /*/*);($;$select 1;*/*/ col CREATE FUNCTION mmee() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END AS x rc$b$x +; + ');(END ;\' + --;;$x$;
|  SELECT /*/ *;*/ 1;+ '; END ;--;\' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*);($;$select 1;*/*/ col CREATE FUNCTION mmee() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END AS x rc$b$x +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*);($;$select 1;*/*/ col CREATE FUNCTION mmee() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END AS x rc$b$x +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + "$;$" 1 hs$z$o AS x U&'a;b;' SELECT;
| FROM t WHERE a = e'\\\n;\';' 1 B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "$;$" 1 hs$z$o AS x U&'a;b;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "$;$" 1 hs$z$o AS x U&'a;b;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY _qn FROM stdin;
| a;b
| c;d;e
| \.
| + E'\n;'';' FROM t WHERE a = /*/*/*; ;);(*/*/*/ AS x '$;$' +;
| COPY mmd FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _qn FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _qn FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT $$/*;;;$$ SELECT ");(end;" SELECT CREATE FUNCTION mwtw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);('; END +;
| SELECT X'01;10' FROM t WHERE a =;
| col we$p$y , CREATE FUNCTION joji() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END col U&'--;;' AS x;
| 1 (SELECT 1; SELECT 2) , (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "");(end;""
| position: 25
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$/*;;;$$ SELECT ");(end;" SELECT CREATE FUNCTION mwtw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);('; END +;
< error
| 42601 syntax error at or near "");(end;""
| position: 25
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$/*;;;$$ SELECT ");(end;" SELECT CREATE FUNCTION mwtw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);('; END +;
< error
| 42601 syntax error at or near "");(end;""
| position: 25
> exec
| , 'select 1;;' SELECT $$select 1;$;$a;b$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , 'select 1;;' SELECT $$select 1;$;$a;b$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , 'select 1;;' SELECT $$select 1;$;$a;b$$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = x'01;10' AS x '$x$;$x$;a;b\' + ');(--;select 1;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' AS x '$x$;$x$;a;b\' + ');(--;select 1;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' AS x '$x$;$x$;a;b\' + ');(--;select 1;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 $$a;b--;$$ AS x CREATE FUNCTION irce() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END , (SELECT 1; SELECT 2) 1; -- trailing note;
| , (SELECT 1; SELECT 2) SELECT x'01;10' + CREATE FUNCTION ogxw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; END ,;
| FROM t WHERE a = --END ;; 
|  1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $$a;b--;$$ AS x CREATE FUNCTION irce() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $$a;b--;$$ AS x CREATE FUNCTION irce() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; END , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + dx$b$m AS x 'end;/*;\' FROM t WHERE a = U&');(; a;b' AS x; -- trailing note;
| FROM t WHERE a = /*--;;);(*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + dx$b$m AS x 'end;/*;\' FROM t WHERE a = U&');(; a;b' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + dx$b$m AS x 'end;/*;\' FROM t WHERE a = U&');(; a;b' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = x'01;10' SELECT /*/*select 1;a;b*/*/ SELECT;
| COPY qaa FROM stdin;
| c;d;e
| a;b
| \.
|
| SELECT --;end;select 1;
|  FROM t WHERE a = wg$i$a +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' SELECT /*/*select 1;a;b*/*/ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' SELECT /*/*select 1;a;b*/*/ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY yzn FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY xft FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yzn FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yzn FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = $r$;;;$x$$r$ FROM t WHERE a = CREATE FUNCTION fwlb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;b);('; END col;
| SELECT U&'$;$' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $r$;;;$x$$r$ FROM t WHERE a = CREATE FUNCTION fwlb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;b);('; END col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $r$;;;$x$$r$ FROM t WHERE a = CREATE FUNCTION fwlb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;b);('; END col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY ozu FROM stdin;
| 1	\N
| 1	\N
| \.
|  col --end;/*;
|  SELECT;
| SELECT U&';;/*;END ;' + CREATE FUNCTION egrb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;;'; SELECT 1; END 1 CREATE FUNCTION bbda() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;b$;$'; SELECT 1; END ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ozu FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ozu FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT CREATE FUNCTION j_gh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END SELECT --END ;
|  SELECT U&'$;$' SELECT; -- trailing note;
| SELECT x'01;10' AS x /*/*/*a;b);($;$*/*/*/ AS x;
| , ';;' col 'end;end;;;' FROM t WHERE a = /*/*/*--;select 1;--;*/*/*/ AS x;FROM t WHERE a = os$g$f + $jp$a;b; $jp$ SELECT /*/*/*);(*/*/*/ 1;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION j_gh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END SELECT --END ;
|  SELECT U&'$;$' SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION j_gh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END SELECT --END ;
|  SELECT U&'$;$' SELECT;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| 1 e'\n;\';\';' FROM t WHERE a = /*/*/*a;b/ *;*/*/*/ + ");(" FROM t WHERE a =;COPY wak FROM stdin;
| c;d;e
| c;d;e
| \.
|
| AS x "a;bselect 1;" FROM t WHERE a = e'\\' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\n;\';\';' FROM t WHERE a = /*/*/*a;b/ *;*/*/*/ + ");(" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\n;\';\';' FROM t WHERE a = /*/*/*a;b/ *;*/*/*/ + ");(" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) col CREATE FUNCTION cqai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;'; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) col CREATE FUNCTION cqai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;'; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) col CREATE FUNCTION cqai() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(end;'; END ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY zxq FROM stdin;
| c;d;e
| c;d;e
| \.
|
| 1 x'01;10' FROM t WHERE a = "end;" , x'01;10' 1; SELECT U&'end;END ;' , e''';\';' 1;COPY baq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zxq FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zxq FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x ";;" FROM t WHERE a = "$;$;;;" , --);();(
|  AS x;
| SELECT $lb$/*;;$lx$$lb$ AS x --a;b; ;
|  AS x (SELECT 1; SELECT 2) col;
| 1 --end;--;select 1;
|  col --$x$;--;
|  , '; \' AS x;
| COPY crt FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ";;" FROM t WHERE a = "$;$;;;" , --);();(
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ";;" FROM t WHERE a = "$;$;;;" , --);();(
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 /*/*/*end;);(*/*/*/ SELECT CREATE FUNCTION bqrq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;/*;'; SELECT 1; END +;FROM t WHERE a = $kcfr$--;);($kcfx$$kcfr$ FROM t WHERE a = X'01;10' SELECT "/*;;" SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*end;);(*/*/*/ SELECT CREATE FUNCTION bqrq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;/*;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*end;);(*/*/*/ SELECT CREATE FUNCTION bqrq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;/*;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT ja$q$t SELECT '$x$;);(\' ,;COPY zux FROM stdin;
| a;b
| a;b
| \.
|
| + "$;$$;$" +;
< error
| 42601 syntax error at or near "'$x$;);(\\'"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ja$q$t SELECT '$x$;);(\' ,;
< error
| 42601 syntax error at or near "'$x$;);(\\'"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT ja$q$t SELECT '$x$;);(\' ,;
< error
| 42601 syntax error at or near "'$x$;);(\\'"
| position: 22
> exec
| SELECT ";;" col U&'end;$x$;' 1;
| + "$x$;; a;b" + /*/*/ *;);(*/*/ col "/*;$x$;" AS x;
| 1 U&'--;; ' , ";" col B'01;10' 1; COPY sz_ FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "U&'end;$x$;'"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ";;" col U&'end;$x$;' 1;
< error
| 42601 syntax error at or near "U&'end;$x$;'"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT ";;" col U&'end;$x$;' 1;
< error
| 42601 syntax error at or near "U&'end;$x$;'"
| position: 17
> exec
| 1 rj$i$m SELECT X'01;10' , U&';;$x$;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 rj$i$m SELECT X'01;10' , U&';;$x$;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 rj$i$m SELECT X'01;10' , U&';;$x$;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x $$);(END ;$$ FROM t WHERE a = B'01;10' AS x CREATE FUNCTION qttk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$);('; END col; COPY yol FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $$);(END ;$$ FROM t WHERE a = B'01;10' AS x CREATE FUNCTION qttk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$);('; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $$);(END ;$$ FROM t WHERE a = B'01;10' AS x CREATE FUNCTION qttk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$);('; END col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY cqu FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY uvv FROM stdin;
| c;d;e
| a;b
| \.
|
| FROM t WHERE a = U&';;end;' AS x /*/*/*end;;*/*/*/ 1 ot$s$v SELECT;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cqu FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cqu FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x $kdmo$$;$$x$;$kdmx$$kdmo$ , e'\';\\' , "END ;);(" SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $kdmo$$;$$x$;$kdmx$$kdmo$ , e'\';\\' , "END ;);(" SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $kdmo$$;$$x$;$kdmx$$kdmo$ , e'\';\\' , "END ;);(" SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , CREATE FUNCTION wvsj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;/*;'; SELECT 1; END 1 "a;b" AS x; SELECT (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION wvsj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;/*;'; SELECT 1; END 1 "a;b" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION wvsj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;/*;'; SELECT 1; END 1 "a;b" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 e'\\\\\';' + e''';' + --/*;
|  AS x; -- trailing note;
| + $$;;$$ col ----;
|  SELECT;
| AS x B'01;10' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 e'\\\\\';' + e''';' + --/*;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 e'\\\\\';' + e''';' + --/*;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , U&'$x$;' FROM t WHERE a =;
| , /*/*/*$x$;END ;*/*/*/ 1 /*/*);(--;$;$*/*/ 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , U&'$x$;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , U&'$x$;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = $a$/*;end;$a$ SELECT $$a;b$$ AS x x'01;10' FROM t WHERE a =; + --END ;;;
|  AS x (SELECT 1; SELECT 2) SELECT;
| , $$;$$ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $a$/*;end;$a$ SELECT $$a;b$$ AS x x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $a$/*;end;$a$ SELECT $$a;b$$ AS x x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , hg$m$x SELECT CREATE FUNCTION ungi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;/*;'; END , CREATE FUNCTION xtaz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;'; END SELECT; -- trailing note;
| AS x "--;$;$;;" FROM t WHERE a = U&'end;' FROM t WHERE a = --;;);(;
|  col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , hg$m$x SELECT CREATE FUNCTION ungi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;/*;'; END , CREATE FUNCTION xtaz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;'; END SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , hg$m$x SELECT CREATE FUNCTION ungi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;/*;'; END , CREATE FUNCTION xtaz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;$x$;'; END SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT /*/*select 1;;*/*/ +; AS x /*);(*/ + /*/*end;*/*/ FROM t WHERE a = "select 1;$;$" 1;
< error
| 42601 syntax error at or near ";"
| position: 28
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*select 1;;*/*/ +;
< error
| 42601 syntax error at or near ";"
| position: 28
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*/*select 1;;*/*/ +;
< error
| 42601 syntax error at or near ";"
| position: 28
> exec
| col b'01;10' FROM t WHERE a = X'01;10' col;
| FROM t WHERE a = /*;*/ +;
| col ';;--;;;' FROM t WHERE a = $dd_j$$x$;$;$$dd_x$$dd_j$ FROM t WHERE a = "--;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col b'01;10' FROM t WHERE a = X'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col b'01;10' FROM t WHERE a = X'01;10' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY jgj FROM stdin;
| 1	\N
| a;b
| \.
|  -- trailing note;
| + hn$c$t + "select 1;--;END ;" SELECT;
| 1 'END ;' + "; " SELECT "; " SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jgj FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jgj FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , $myl$);($;$end;$myl$ SELECT U&'; /*;' + "select 1;" +;
| + $j$select 1;$x$;$x$$j$ +;
| COPY eol FROM stdin;
| a;b
| c;d;e
| \.
|
| + --END ;;;select 1;
|  SELECT CREATE FUNCTION tutg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; end;/*;'; END ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $myl$);($;$end;$myl$ SELECT U&'; /*;' + "select 1;" +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $myl$);($;$end;$myl$ SELECT U&'; /*;' + "select 1;" +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + U&';' col /*);(*/ AS x;1 E'\';\';'';' col;
| SELECT E'\';\n;\';' FROM t WHERE a = /*/ *;$x$;/ *;*/ FROM t WHERE a = $m$;$;$$x$$m$ AS x;
| COPY fqx FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&';' col /*);(*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&';' col /*);(*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY aak FROM stdin;
| c;d;e
| 1	\N
| \.
|  SELECT /*/*/*; */*/*/ col;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY aak FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY aak FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = nh$b$y 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = nh$b$y 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = nh$b$y 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) col $$--;$$ AS x --;
|  AS x; col tf$n$v ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) col $$--;$$ AS x --;
|  AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) col $$--;$$ AS x --;
|  AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY agd FROM stdin;
| a;b
| 1	\N
| \.
| + B'01;10' FROM t WHERE a = (SELECT 1; SELECT 2) col '/*;end;$;$' 1;COPY lrw FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY agd FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY agd FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + $dmnu$;;);(END ;$dmnu$ AS x; SELECT $i$;;; $i$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $dmnu$;;);(END ;$dmnu$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $dmnu$;;);(END ;$dmnu$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + 'a;bEND ;$x$;' FROM t WHERE a = X'01;10' 1;, $sv$$;$$sx$$sv$ AS x vr$b$p col;
| COPY _zw FROM stdin;
| 1	\N
| 1	\N
| \.
|
| COPY hjb FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + 'a;bEND ;$x$;' FROM t WHERE a = X'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + 'a;bEND ;$x$;' FROM t WHERE a = X'01;10' 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + CREATE FUNCTION fado() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;END ;; '; SELECT 1; END SELECT $newf$);($newf$ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION fado() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;END ;; '; SELECT 1; END SELECT $newf$);($newf$ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION fado() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;END ;; '; SELECT 1; END SELECT $newf$);($newf$ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , --$;$a;bEND ;
|  , CREATE FUNCTION qyvs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;END ;END ;'; END , "--;;" col; -- trailing note;
| COPY tws FROM stdin;
| c;d;e
| 1	\N
| \.
|
| 1 U&'$;$' 1 "--;" FROM t WHERE a = x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --$;$a;bEND ;
|  , CREATE FUNCTION qyvs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;END ;END ;'; END , "--;;" col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --$;$a;bEND ;
|  , CREATE FUNCTION qyvs() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;END ;END ;'; END , "--;;" col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT ';;;;\' 1 U&'$;$select 1;--;' FROM t WHERE a =; -- trailing note;
| FROM t WHERE a = e'\n;\';\\' + CREATE FUNCTION seza() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;--;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 16
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ';;;;\' 1 U&'$;$select 1;--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 16
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT ';;;;\' 1 U&'$;$select 1;--;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 16
> exec
| AS x --END ;
|  1 e'\';'';\n;' AS x --);($x$;END ;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --END ;
|  1 e'\';'';\n;' AS x --);($x$;END ;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --END ;
|  1 e'\';'';\n;' AS x --);($x$;END ;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY mrd FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mrd FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mrd FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , e'\';\n;' col;
| , (SELECT 1; SELECT 2) FROM t WHERE a = eu$l$h SELECT U&'a;b' col;, /*/*/*a;bselect 1;$;$*/*/*/ AS x U&';; ' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e'\';\n;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e'\';\n;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x CREATE FUNCTION oubm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;--;'; END FROM t WHERE a = B'01;10' SELECT ';;;;\' 1; SELECT ');(' , CREATE FUNCTION prdc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END ,; -- trailing note;
| , /*;;*/ 1 $ng$a;b/*;$ng$ + --end;);(
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION oubm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;--;'; END FROM t WHERE a = B'01;10' SELECT ';;;;\' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION oubm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;--;'; END FROM t WHERE a = B'01;10' SELECT ';;;;\' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + b'01;10' +;
| COPY ons FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + b'01;10' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + b'01;10' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 ----;end;END ;
|  AS x; 1 U&';;' col (SELECT 1; SELECT 2) +;
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = /*/*; */*/ 1 CREATE FUNCTION ayxo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ----;end;END ;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ----;end;END ;
|  AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY eo_ FROM stdin;
| a;b
| 1	\N
| \.
|
| , "/*;a;bselect 1;" AS x oj$k$w AS x; 1 CREATE FUNCTION ciwm() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);('; SELECT 1; END + U&'$;$END ;' SELECT;
| COPY niz FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY eo_ FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY eo_ FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY hon FROM stdin;
| 1	\N
| 1	\N
| \.
|  + 'a;b);(; ' AS x /*/*$;$);(*/*/ SELECT X'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hon FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hon FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT $tj$select 1;$tj$ col (SELECT 1; SELECT 2) , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "("
| position: 30
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $tj$select 1;$tj$ col (SELECT 1; SELECT 2) , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "("
| position: 30
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $tj$select 1;$tj$ col (SELECT 1; SELECT 2) , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "("
| position: 30
> exec
| FROM t WHERE a = U&'a;b;;END ;' +;
| COPY rds FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'a;b;;END ;' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'a;b;;END ;' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + CREATE FUNCTION qbli() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;a;b'; END col is$_$v col /*);(*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION qbli() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;a;b'; END col is$_$v col /*);(*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION qbli() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;a;b'; END col is$_$v col /*);(*/ +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x "a;b);(a;b" FROM t WHERE a = /*/*/*;;$;$*/*/*/ col (SELECT 1; SELECT 2) ,;, $$;;end;;$$ ,;FROM t WHERE a = E'\n;'';' col U&'$x$;$x$;' +;
| FROM t WHERE a = $$; end;$x$;$$ , /*/*$;$;;*/*/ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "a;b);(a;b" FROM t WHERE a = /*/*/*;;$;$*/*/*/ col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "a;b);(a;b" FROM t WHERE a = /*/*/*;;$;$*/*/*/ col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + (SELECT 1; SELECT 2) ,;
| col U&'end;' , $ef$;$x$;$ex$$ef$ AS x;1 "--;);(a;b" FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY jdq FROM stdin;
| 1	\N
| 1	\N
| \.
| , bd$r$q col CREATE FUNCTION ulpo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;a;bend;'; SELECT 1; END SELECT ";;;" FROM t WHERE a =;
| AS x (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jdq FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jdq FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x U&'select 1;' +;
| 1 $$;$x$;$$ 1 B'01;10' 1;
| col CREATE FUNCTION tktc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;b);('; END SELECT;
| 1 /*/*; ; ;*/*/ + ';;;;;' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'select 1;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'select 1;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = cc$g$k 1 "a;b" 1 --END ;select 1;end;
|  1; -- trailing note;
| col e''';\';' FROM t WHERE a = "END ;/*;end;" AS x ie$r$m FROM t WHERE a =;, (SELECT 1; SELECT 2) AS x;
| FROM t WHERE a = CREATE FUNCTION qhjp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;end;$;$'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = cc$g$k 1 "a;b" 1 --END ;select 1;end;
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = cc$g$k 1 "a;b" 1 --END ;select 1;end;
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col CREATE FUNCTION ndts() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION ndts() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION ndts() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT x'01;10' SELECT "$;$/*;--;" FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
| COPY pty FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| FROM t WHERE a = ');(; ' SELECT $lk$select 1;/*;END ;$lx$$lk$ SELECT;
< error
| 42601 syntax error at or near ""$;$/*;--;""
| position: 24
> exec
| BEGIN
< ok 0 -
> exec
| SELECT x'01;10' SELECT "$;$/*;--;" FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ""$;$/*;--;""
| position: 24
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT x'01;10' SELECT "$;$/*;--;" FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ""$;$/*;--;""
| position: 24
> exec
| 1 CREATE FUNCTION rcjd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; end;; '; END , (SELECT 1; SELECT 2) +;
| + $$END ;$;$$$ +;
| + ';;\' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION rcjd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; end;; '; END , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION rcjd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; end;; '; END , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 '$x$;/*;/*;' col; 1 $qwhs$--;);($qwhs$ AS x; -- trailing note;
| + md$z$l , vo$j$z SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '$x$;/*;/*;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '$x$;/*;/*;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY flu FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY sap FROM stdin;
| c;d;e
| c;d;e
| \.
|
| col "end;;end;" FROM t WHERE a = $$a;b; $;$$$ , ';;$;$select 1;' +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY flu FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY flu FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY zyc FROM stdin;
| a;b
| a;b
| \.
| SELECT b'01;10' SELECT /*/*/*select 1;;;end;*/*/*/ col;
| , e''';' + --;;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zyc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zyc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col E''';' AS x; -- trailing note;
| SELECT x'01;10' AS x;
| + /*/*select 1;*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col E''';' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col E''';' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = E'\\'';' col;
| FROM t WHERE a = e'\n;'';' , "--;" ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E'\\'';' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\\'';' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY too FROM stdin;
| a;b
| c;d;e
| \.
|
| 1 ----;
|  AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY too FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY too FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x --select 1;a;b; 
|  SELECT; 1 U&';' 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --select 1;a;b; 
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --select 1;a;b; 
|  SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT U&';;;;' AS x;COPY k_r FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 43
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&';;;;' AS x;
< ok 1 -
> exec
| COPY k_r FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY k_r FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col X'01;10' AS x; SELECT /*/*/*END ;*/*/*/ 1;
| SELECT --$x$;$x$;
|  , $$select 1;; /*;$$ 1; -- trailing note;
| SELECT kt$b$s col x'01;10' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col X'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col X'01;10' AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY gym FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 CREATE FUNCTION przd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gym FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gym FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 '$;$$x$;;' 1 e'\\\';\';' col; AS x in$o$e SELECT /*/*/*end;END ;;;*/*/*/ FROM t WHERE a = --;/*;--;
|  1;
| + rd$e$s col;COPY nbe FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '$;$$x$;;' 1 e'\\\';\';' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '$;$$x$;;' 1 e'\\\';\';' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = ss$h$v AS x; + "a;bEND ;/*;" 1;
| + "--;select 1;END ;" AS x e'\';\';\\' FROM t WHERE a =;
| COPY tvk FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ss$h$v AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ss$h$v AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x e'\n;' SELECT;
| SELECT /*/*;;END ;;*/*/ + (SELECT 1; SELECT 2) 1;AS x b'01;10' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x e'\n;' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x e'\n;' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = ";" , U&'$x$;--;/*;' FROM t WHERE a = U&'END ;);(;;' ,;
| AS x U&';;' FROM t WHERE a = ih$n$t FROM t WHERE a =; -- trailing note;
| , U&'; ' ,; -- trailing note;
| COPY bum FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ";" , U&'$x$;--;/*;' FROM t WHERE a = U&'END ;);(;;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ";" , U&'$x$;--;/*;' FROM t WHERE a = U&'END ;);(;;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = --END ;);(;
|  1; COPY dfl FROM stdin;
| a;b
| 1	\N
| \.
|  -- trailing note;
| SELECT "a;b" SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --END ;);(;
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --END ;);(;
|  1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + CREATE FUNCTION apyp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bselect 1;/*;'; END AS x B'01;10' + $$);($;$$$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION apyp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bselect 1;/*;'; END AS x B'01;10' + $$);($;$$$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION apyp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bselect 1;/*;'; END AS x B'01;10' + $$);($;$$$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , "END ;; " + ';end;$x$;' + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "END ;; " + ';end;$x$;' + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "END ;; " + ';end;$x$;' + (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| + U&'select 1;' + X'01;10' AS x (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'select 1;' + X'01;10' AS x (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'select 1;' + X'01;10' AS x (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY boe FROM stdin;
| c;d;e
| 1	\N
| \.
|
| 1 '/*;' 1;
| SELECT U&'$;$' , ";;);(" SELECT (SELECT 1; SELECT 2) SELECT;
| 1 --$x$;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY boe FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY boe FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x U&'END ;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'END ;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'END ;' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + CREATE FUNCTION l_cv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;bend;'; SELECT 1; END , "END ;select 1;;" SELECT;
| + U&'; select 1;end;' AS x;
| , (SELECT 1; SELECT 2) 1 'a;b\' col;1 b'01;10' SELECT b'01;10' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION l_cv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;bend;'; SELECT 1; END , "END ;select 1;;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION l_cv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;a;bend;'; SELECT 1; END , "END ;select 1;;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT /*/*;;$;$*/*/ AS x "/*;end;;;" , /*$x$;$x$;end;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*/*;;$;$*/*/ AS x "/*;end;;;" , /*$x$;$x$;end;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*/*;;$;$*/*/ AS x "/*;end;;;" , /*$x$;$x$;end;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 22
> exec
| 1 $y_r$END ;$y_r$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $y_r$END ;$y_r$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $y_r$END ;$y_r$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + CREATE FUNCTION erjq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END AS x X'01;10' +;
| COPY bxq FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY qar FROM stdin;
| 1	\N
| 1	\N
| \.
| FROM t WHERE a = U&'; ;;a;b' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION erjq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END AS x X'01;10' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION erjq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END AS x X'01;10' +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT CREATE FUNCTION mlnf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END + CREATE FUNCTION vgrk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bEND ;$;$'; SELECT 1; END FROM t WHERE a =;
| + b'01;10' col;
| 1 (SELECT 1; SELECT 2) col "$;$a;b" col cm$u$f +;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION mlnf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END + CREATE FUNCTION vgrk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bEND ;$;$'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION mlnf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;'; SELECT 1; END + CREATE FUNCTION vgrk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bEND ;$;$'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| SELECT ku$f$_ +; -- trailing note;
| , e'\\' + E'\\\n;' +;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ku$f$_ +;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT ku$f$_ +;
< error
| 42601 syntax error at or near ";"
| position: 16
> exec
| + CREATE FUNCTION _zct() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$'; END +;+ U&'END ;select 1;' col ";" ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION _zct() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$'; END +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION _zct() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$'; END +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY yrd FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY giw FROM stdin;
| a;b
| c;d;e
| \.
|  -- trailing note;
| + ');(end;' 1 $axw$;;; $axx$$axw$ 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yrd FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yrd FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , "--;);(/*;" FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "--;);(/*;" FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "--;);(/*;" FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = --END ;/*;--;
|  FROM t WHERE a = pp$z$h +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --END ;/*;--;
|  FROM t WHERE a = pp$z$h +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --END ;/*;--;
|  FROM t WHERE a = pp$z$h +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , CREATE FUNCTION gjsk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;ba;b'; SELECT 1; END FROM t WHERE a = pr$v$q , ay$o$d FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION gjsk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;ba;b'; SELECT 1; END FROM t WHERE a = pr$v$q , ay$o$d FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION gjsk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;ba;b'; SELECT 1; END FROM t WHERE a = pr$v$q , ay$o$d FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , 'end;; \' +;
| COPY kkv FROM stdin;
| c;d;e
| c;d;e
| \.
|
| SELECT --end;
|  FROM t WHERE a = "a;b" 1 E''';'';\n;' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , 'end;; \' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , 'end;; \' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY dok FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dok FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dok FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , (SELECT 1; SELECT 2) +;, x'01;10' FROM t WHERE a = U&'$;$' col $iqpy$/*;$iqpx$$iqpy$ FROM t WHERE a =;
| FROM t WHERE a = jm$q$y 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = /*/*/*;end;*/*/*/ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*/*/*;end;*/*/*/ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*/*/*;end;*/*/*/ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 /*/ *;a;b*/ col (SELECT 1; SELECT 2) , U&'$x$;END ;;' col; SELECT B'01;10' col "--;" AS x ----;--;END ;
|  1;
| AS x (SELECT 1; SELECT 2) col ps$a$k + x'01;10' SELECT;
| + B'01;10' FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/ *;a;b*/ col (SELECT 1; SELECT 2) , U&'$x$;END ;;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/ *;a;b*/ col (SELECT 1; SELECT 2) , U&'$x$;END ;;' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| , ----;END ;
|  1 CREATE FUNCTION urvb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;a;bselect 1;'; SELECT 1; END SELECT "end;/*;" ,; + ----;--;;
|  AS x;
| AS x "select 1;END ;;" SELECT $$/*;;;$$ + b'01;10' SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ----;END ;
|  1 CREATE FUNCTION urvb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;a;bselect 1;'; SELECT 1; END SELECT "end;/*;" ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ----;END ;
|  1 CREATE FUNCTION urvb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;a;bselect 1;'; SELECT 1; END SELECT "end;/*;" ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT /*a;b$x$;select 1;*/ AS x 'select 1;' SELECT (SELECT 1; SELECT 2) +;
| col b'01;10' + --);(
|  SELECT X'01;10' ,;
< error
| 42601 syntax error at or near "AS"
| position: 29
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*a;b$x$;select 1;*/ AS x 'select 1;' SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "AS"
| position: 29
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT /*a;b$x$;select 1;*/ AS x 'select 1;' SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "AS"
| position: 29
> exec
| 1 $$--;end;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 $$--;end;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 $$--;end;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x ar$_$j AS x vo$k$b col; -- trailing note;
| SELECT --END ;;
|  ,;
| 1 "/*;);(" 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ar$_$j AS x vo$k$b col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ar$_$j AS x vo$k$b col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , CREATE FUNCTION kdyr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END + x'01;10' SELECT x'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION kdyr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END + x'01;10' SELECT x'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION kdyr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; SELECT 1; END + x'01;10' SELECT x'01;10' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x --END ;
|  FROM t WHERE a = E''';\';' col CREATE FUNCTION rahk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;end;$x$;'; END AS x;
| , U&');(' AS x /*/*$x$;; ; */*/ + x'01;10' ,;
| SELECT CREATE FUNCTION umqf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END FROM t WHERE a = --;;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --END ;
|  FROM t WHERE a = E''';\';' col CREATE FUNCTION rahk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;end;$x$;'; END AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --END ;
|  FROM t WHERE a = E''';\';' col CREATE FUNCTION rahk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;end;$x$;'; END AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , B'01;10' 1; FROM t WHERE a = x'01;10' AS x e''';\n;' + $$END ;end;$$ SELECT;
| COPY udu FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , B'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , B'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col e'\\\n;'';' SELECT U&'end;END ;;;' col;
| SELECT "end;);(" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e'\\\n;'';' SELECT U&'end;END ;;;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e'\\\n;'';' SELECT U&'end;END ;;;' col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY tbn FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| , "end;;;" FROM t WHERE a = $$;;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY tbn FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY tbn FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY hkq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hkq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hkq FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 CREATE FUNCTION xagf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;;'; END FROM t WHERE a =;
| AS x "--;--;END ;" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION xagf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION xagf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$$x$;;'; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY qai FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY qai FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY qai FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| + "; select 1;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "; select 1;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "; select 1;" SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY _rx FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY aoy FROM stdin;
| a;b
| c;d;e
| \.
|
| AS x $tqsa$END ;a;bEND ;$tqsa$ AS x U&'; a;bend;' SELECT /*/*a;b*/*/ FROM t WHERE a =;
| COPY uql FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _rx FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _rx FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col U&'$x$;;;' + '--;' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'$x$;;;' + '--;' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'$x$;;;' + '--;' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY oxe FROM stdin;
| 1	\N
| a;b
| \.
|  -- trailing note;
| AS x CREATE FUNCTION yxld() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');('; SELECT 1; END +; -- trailing note;
| AS x --$x$;select 1;
|  col x'01;10' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY oxe FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY oxe FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY hzn FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hzn FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hzn FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY _cv FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _cv FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _cv FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;1 (SELECT 1; SELECT 2) SELECT U&'select 1;' FROM t WHERE a =; -- trailing note;
| SELECT --;$x$;
|  AS x;SELECT /*/*/*;*/*/*/ , U&'$x$;;;' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = ym$g$m ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = ym$g$m ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = ym$g$m ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col e''';' , b'01;10' FROM t WHERE a = --END ;;;select 1;
|  1; col ";;; " , /*/*; */*/ AS x $ck$;;;$ck$ +;
| COPY nfc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e''';' , b'01;10' FROM t WHERE a = --END ;;;select 1;
|  1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e''';' , b'01;10' FROM t WHERE a = --END ;;;select 1;
|  1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = $$; $$ SELECT /*/ *;*/ , "$;$$;$" AS x;AS x X'01;10' 1 /*/*END ;end;*/*/ col E''';\n;' 1;
| SELECT yv$p$k FROM t WHERE a = --select 1;a;b
|  FROM t WHERE a = $$;a;ba;b$$ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$; $$ SELECT /*/ *;*/ , "$;$$;$" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$; $$ SELECT /*/ *;*/ , "$;$$;$" AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 B'01;10' col $$);($$ , --;select 1;
|  SELECT;
| COPY rse FROM stdin;
| c;d;e
| 1	\N
| \.
| 1 CREATE FUNCTION oogy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;--;'; END + "$x$;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 B'01;10' col $$);($$ , --;select 1;
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 B'01;10' col $$);($$ , --;select 1;
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY cpw FROM stdin;
| a;b
| a;b
| \.
|
| SELECT _g$c$u ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cpw FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cpw FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 --; 
|  col --);($;$
|  AS x "; " AS x; FROM t WHERE a = CREATE FUNCTION pr_j() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END + CREATE FUNCTION iuiv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;;;--;'; END 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --; 
|  col --);($;$
|  AS x "; " AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --; 
|  col --);($;$
|  AS x "; " AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col $opy$; $opx$$opy$ + ");();(END ;" SELECT; AS x (SELECT 1; SELECT 2) AS x;
| SELECT (SELECT 1; SELECT 2) SELECT "$x$;" SELECT (SELECT 1; SELECT 2) +; AS x CREATE FUNCTION lmov() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; SELECT 1; END col (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $opy$; $opx$$opy$ + ");();(END ;" SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $opy$; $opx$$opy$ + ");();(END ;" SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + ph$u$m col "end;);(" SELECT $nn$select 1;;;$nx$$nn$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ph$u$m col "end;);(" SELECT $nn$select 1;;;$nx$$nn$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ph$u$m col "end;);(" SELECT $nn$select 1;;;$nx$$nn$ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + CREATE FUNCTION yduu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; $;$'; SELECT 1; END FROM t WHERE a = x'01;10' SELECT '--;end;' SELECT; FROM t WHERE a = X'01;10' SELECT U&';; ' 1 X'01;10' FROM t WHERE a =;
| SELECT $$END ;--;$$ 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + CREATE FUNCTION yduu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; $;$'; SELECT 1; END FROM t WHERE a = x'01;10' SELECT '--;end;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + CREATE FUNCTION yduu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;; $;$'; SELECT 1; END FROM t WHERE a = x'01;10' SELECT '--;end;' SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x uf$u$c +;
| COPY _zi FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x uf$u$c +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x uf$u$c +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| col (SELECT 1; SELECT 2) AS x;
| FROM t WHERE a = E'\';\';\\' AS x "--;$;$END ;" SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY w_e FROM stdin;
| c;d;e
| c;d;e
| \.
|
| SELECT CREATE FUNCTION ylua() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END 1;
| , 'a;b' 1; + --/*;; 
|  1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY w_e FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY w_e FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY _km FROM stdin;
| c;d;e
| c;d;e
| \.
|  + ';' SELECT U&'; a;b' ,;FROM t WHERE a = CREATE FUNCTION pots() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;END ;END ;'; END AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _km FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _km FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT --END ;$x$;
|  AS x;
| AS x nh$p$m FROM t WHERE a = "select 1;end;--;" + CREATE FUNCTION qwyx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';);('; END AS x;
< error
| 42601 syntax error at or near "AS"
| position: 21
> exec
| BEGIN
< ok 0 -
> exec
| SELECT --END ;$x$;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 21
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT --END ;$x$;
|  AS x;
< error
| 42601 syntax error at or near "AS"
| position: 21
> exec
| col U&'/*;;;$;$' , /*/*/ *;end;*/*/ , --;end;select 1;
|  FROM t WHERE a =; FROM t WHERE a = (SELECT 1; SELECT 2) 1 (SELECT 1; SELECT 2) SELECT; -- trailing note;
| + --a;b
|  col; 1 (SELECT 1; SELECT 2) AS x E'\';\n;'';' + --);($x$;
|  ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'/*;;;$;$' , /*/*/ *;end;*/*/ , --;end;select 1;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'/*;;;$;$' , /*/*/ *;end;*/*/ , --;end;select 1;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + /*/*/*;;$;$*/*/*/ AS x /*);(*/ AS x;col 'end;--;\' ,; COPY gxh FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*/*;;$;$*/*/*/ AS x /*);(*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*/*;;$;$*/*/*/ AS x /*);(*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col $clj$; $x$;; $clj$ FROM t WHERE a =; AS x "END ;$x$;a;b" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $clj$; $x$;; $clj$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $clj$; $x$;; $clj$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + '; \' 1 /*; a;b*/ col (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + '; \' 1 /*; a;b*/ col (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + '; \' 1 /*; a;b*/ col (SELECT 1; SELECT 2) FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY rwd FROM stdin;
| c;d;e
| a;b
| \.
|
| , CREATE FUNCTION efnh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';; '; END col "--;" + (SELECT 1; SELECT 2) ,;
| SELECT $$end;$$ , "; " 1 /*/*END ;*/*/ 1;
| SELECT U&');(/*;);(' 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY rwd FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY rwd FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY dfw FROM stdin;
| a;b
| 1	\N
| \.
|
| AS x 'end;;;;' , 'end;--;' AS x B'01;10' SELECT;
| COPY umf FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY dfw FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY dfw FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , "/*;" FROM t WHERE a = e'\n;' FROM t WHERE a = --/*;;
|  ,;
| FROM t WHERE a = $$end;; $$ AS x;
| COPY xec FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| COPY ilq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "/*;" FROM t WHERE a = e'\n;' FROM t WHERE a = --/*;;
|  ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "/*;" FROM t WHERE a = e'\n;' FROM t WHERE a = --/*;;
|  ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = --a;b$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = --a;b$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) FROM t WHERE a = --a;b$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 gd$g$b AS x E'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 gd$g$b AS x E'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 gd$g$b AS x E'\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY ark FROM stdin;
| c;d;e
| a;b
| \.
| COPY zhz FROM stdin;
| 1	\N
| a;b
| \.
| col U&'end;' AS x /*/*$x$;*/*/ + $$);($$ ,;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ark FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ark FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 U&'END ;$x$;END ;' AS x /*/*/ *;*/*/ AS x;
| COPY huo FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'END ;$x$;END ;' AS x /*/*/ *;*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'END ;$x$;END ;' AS x /*/*/ *;*/*/ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = o_$k$v 1;
| 1 rn$f$q col U&'; ' FROM t WHERE a = "END ;END ;" col; -- trailing note;
| col '; ' col --;end;
|  ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = o_$k$v 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = o_$k$v 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY e_u FROM stdin;
| a;b
| c;d;e
| \.
| SELECT x'01;10' , ';;end;' 1 (SELECT 1; SELECT 2) FROM t WHERE a =;
| 1 /*/*;;$;$);(*/*/ AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY e_u FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY e_u FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY p_x FROM stdin;
| 1	\N
| a;b
| \.
|
| col e'\';\\' SELECT (SELECT 1; SELECT 2) , U&'a;bend;;;' col; FROM t WHERE a = --/*;
|  col;
| , cu$r$f FROM t WHERE a = X'01;10' col $$$;$$$ col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY p_x FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY p_x FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + /*END ;);(*/ col CREATE FUNCTION mjfe() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END col --/*;
|  SELECT;
| COPY gfv FROM stdin;
| a;b
| 1	\N
| \.
|  + CREATE FUNCTION ws_i() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;select 1;; '; END SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*END ;);(*/ col CREATE FUNCTION mjfe() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END col --/*;
|  SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*END ;);(*/ col CREATE FUNCTION mjfe() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END col --/*;
|  SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col U&'/*;$x$;end;' AS x U&'$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'/*;$x$;end;' AS x U&'$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'/*;$x$;end;' AS x U&'$;$' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , $flq$a;b--;$flq$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $flq$a;b--;$flq$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $flq$a;b--;$flq$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , '; );(--;' +;
| AS x b'01;10' SELECT $$;END ;$$ SELECT CREATE FUNCTION nbqi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;--;'; SELECT 1; END col;
| COPY kve FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , '; );(--;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , '; );(--;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = tu$s$z , ";;" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = tu$s$z , ";;" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = tu$s$z , ";;" +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = --;;$x$;
|  ,;
| 1 --end;;;
|  FROM t WHERE a = X'01;10' +;
| 1 'end;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --;;$x$;
|  ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --;;$x$;
|  ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x ';;--;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ';;--;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ';;--;' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x --;end;
|  AS x e''';' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x --;end;
|  AS x e''';' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x --;end;
|  AS x e''';' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + /*--;a;b; */ , ");(" FROM t WHERE a = "/*;" +;
| COPY bbe FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*--;a;b; */ , ");(" FROM t WHERE a = "/*;" +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*--;a;b; */ , ");(" FROM t WHERE a = "/*;" +;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + "END ;--;a;b" , CREATE FUNCTION yquw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;'; END 1 $zsd$);(; $zsx$$zsd$ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "END ;--;a;b" , CREATE FUNCTION yquw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;'; END 1 $zsd$);(; $zsx$$zsd$ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "END ;--;a;b" , CREATE FUNCTION yquw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;;'; END 1 $zsd$);(; $zsx$$zsd$ col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x x'01;10' SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x x'01;10' SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x x'01;10' SELECT (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY lpa FROM stdin;
| 1	\N
| c;d;e
| \.
|
| AS x (SELECT 1; SELECT 2) SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lpa FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lpa FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , CREATE FUNCTION dqxb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; '; SELECT 1; END , $ld$/*;$ld$ , ----;$;$a;b
|  ,;
| 1 U&';' FROM t WHERE a = X'01;10' 1 U&'/*;; ' ,;
| COPY hgu FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION dqxb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; '; SELECT 1; END , $ld$/*;$ld$ , ----;$;$a;b
|  ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION dqxb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;; '; SELECT 1; END , $ld$/*;$ld$ , ----;$;$a;b
|  ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = --; 
|  col $jeb$end;/*;$jex$$jeb$ 1 /*/*/*$x$;);(*/*/*/ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = --; 
|  col $jeb$end;/*;$jex$$jeb$ 1 /*/*/*$x$;);(*/*/*/ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = --; 
|  col $jeb$end;/*;$jex$$jeb$ 1 /*/*/*$x$;);(*/*/*/ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = '$x$;;;;;' col "end;select 1;" , (SELECT 1; SELECT 2) AS x; -- trailing note;
| COPY hmt FROM stdin;
| 1	\N
| c;d;e
| \.
|
| COPY bov FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = '$x$;;;;;' col "end;select 1;" , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = '$x$;;;;;' col "end;select 1;" , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY sul FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sul FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sul FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT B'01;10' ,;
| AS x 'END ;;;' 1 $eb$select 1;a;b$eb$ col "/*;);(" 1;
| 1 (SELECT 1; SELECT 2) FROM t WHERE a = /*);($;$$x$;*/ col e'\n;\n;' +;
< error
| 42601 syntax error at or near ";"
| position: 18
> exec
| BEGIN
< ok 0 -
> exec
| SELECT B'01;10' ,;
< error
| 42601 syntax error at or near ";"
| position: 18
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT B'01;10' ,;
< error
| 42601 syntax error at or near ";"
| position: 18
> exec
| COPY yob FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yob FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yob FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , $ve$$;$$vx$$ve$ FROM t WHERE a = 'END ;;/*;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $ve$$;$$vx$$ve$ FROM t WHERE a = 'END ;;/*;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $ve$$;$$vx$$ve$ FROM t WHERE a = 'END ;;/*;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY gmt FROM stdin;
| a;b
| 1	\N
| \.
|
| FROM t WHERE a = '$x$;' col E'\\' SELECT X'01;10' ,; , $sgu_$/*;$;$$sgux$$sgu_$ AS x (SELECT 1; SELECT 2) + X'01;10' AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY gmt FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY gmt FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY zg_ FROM stdin;
| 1	\N
| a;b
| \.
|
| 1 (SELECT 1; SELECT 2) 1; -- trailing note;
| + B'01;10' + CREATE FUNCTION _xvw() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;/*;'; END SELECT "--;" +; -- trailing note;
| COPY qss FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY zg_ FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY zg_ FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT $$;;/*;$$ FROM t WHERE a =; -- trailing note;
| COPY qww FROM stdin;
| a;b
| c;d;e
| \.
| + x'01;10' AS x (SELECT 1; SELECT 2) 1;
| 1 U&'/*;;;end;' SELECT $$END ;$$ + --select 1;
|  col;
< error
| 42601 syntax error at or near ";"
| position: 34
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$;;/*;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 34
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$;;/*;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 34
> exec
| col CREATE FUNCTION xxoa() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;/*;'; END AS x;
| COPY ses FROM stdin;
| 1	\N
| c;d;e
| \.
|
| AS x "$;$select 1;END ;" SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION xxoa() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;/*;'; END AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION xxoa() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';end;/*;'; END AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x yl$j$w col $$; $$ 1 "--;a;b--;" +;1 X'01;10' FROM t WHERE a = "--;end;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x yl$j$w col $$; $$ 1 "--;a;b--;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x yl$j$w col $$; $$ 1 "--;a;b--;" +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , '$x$;;\' + sc$h$a AS x x'01;10' 1;
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x; 1 E'\\\n;\\' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , '$x$;;\' + sc$h$a AS x x'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , '$x$;;\' + sc$h$a AS x x'01;10' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY elh FROM stdin;
| c;d;e
| c;d;e
| \.
|
| AS x (SELECT 1; SELECT 2) 1;
| 1 oz$c$w col CREATE FUNCTION xesq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;bend;a;b'; END 1 "$;$" +; -- trailing note;
| COPY sm_ FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY elh FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY elh FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY _ew FROM stdin;
| a;b
| a;b
| \.
|
| col E'\n;' AS x; -- trailing note;
| col /*/*select 1;select 1;*/*/ AS x;
| + ';; \' FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _ew FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _ew FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| AS x --);(--;/*;
|  SELECT U&'--;a;b);(' 1 $$$;$/*;$$ AS x; -- trailing note;
| 1 (SELECT 1; SELECT 2) FROM t WHERE a =;AS x e'\n;\\' 1 vo$o$r 1 CREATE FUNCTION ural() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --);(--;/*;
|  SELECT U&'--;a;b);(' 1 $$$;$/*;$$ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --);(--;/*;
|  SELECT U&'--;a;b);(' 1 $$$;$/*;$$ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = CREATE FUNCTION pjpo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;/*;;;'; SELECT 1; END +; -- trailing note;
| col /*/*$;$*/*/ col;SELECT E''';\\' + '; a;b' 1 /*/*/*/ *;END ;select 1;*/*/*/ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION pjpo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;/*;;;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION pjpo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;/*;;;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , "; END ;" , e'\';' SELECT E'\\' 1;
| COPY lix FROM stdin;
| a;b
| a;b
| \.
|
| SELECT U&';;;;END ;' AS x (SELECT 1; SELECT 2) AS x /*/*/*/ *;*/*/*/ FROM t WHERE a =;
| col U&'; ' , E'\n;\n;' +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "; END ;" , e'\';' SELECT E'\\' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "; END ;" , e'\';' SELECT E'\\' 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x /*/*/*END ;*/*/*/ FROM t WHERE a = iw$f$u SELECT;
| SELECT /*/*/*END ;*/*/*/ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*END ;*/*/*/ FROM t WHERE a = iw$f$u SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*END ;*/*/*/ FROM t WHERE a = iw$f$u SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY udc FROM stdin;
| a;b
| a;b
| \.
|
| + '$;$' , (SELECT 1; SELECT 2) 1 mr$h$n ,;, --END ;--;/*;
|  SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY udc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY udc FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ifj FROM stdin;
| 1	\N
| a;b
| \.
|
| SELECT (SELECT 1; SELECT 2) SELECT /*/*/*);(;;*/*/*/ AS x; -- trailing note;
| col E'\n;\\\';' col $m$end;END ;$x$$m$ SELECT ek$u$f col; -- trailing note;
| SELECT /*/*/*;; a;b*/*/*/ col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ifj FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ifj FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY xms FROM stdin;
| c;d;e
| 1	\N
| \.
|
| , (SELECT 1; SELECT 2) AS x x'01;10' , U&'select 1;' FROM t WHERE a =;
| FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY xms FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY xms FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 U&'end;' + (SELECT 1; SELECT 2) + $zh$a;ba;bselect 1;$zh$ col;
| SELECT --$;$END ;
|  1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 U&'end;' + (SELECT 1; SELECT 2) + $zh$a;ba;bselect 1;$zh$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 U&'end;' + (SELECT 1; SELECT 2) + $zh$a;ba;bselect 1;$zh$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + $$$x$;;; $$ AS x ');(select 1;);(' + E'\n;' AS x;
| , --select 1;$;$a;b
|  SELECT;
| col /*/*;*/*/ + --$x$;);(
|  AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + $$$x$;;; $$ AS x ');(select 1;);(' + E'\n;' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + $$$x$;;; $$ AS x ');(select 1;);(' + E'\n;' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x --; 
|  AS x ";a;b" FROM t WHERE a = /*/*select 1;$x$;;;*/*/ SELECT;
| SELECT (SELECT 1; SELECT 2) AS x;+ "; --;);(" FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --; 
|  AS x ";a;b" FROM t WHERE a = /*/*select 1;$x$;;;*/*/ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --; 
|  AS x ";a;b" FROM t WHERE a = /*/*select 1;$x$;;;*/*/ SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x gs$y$k AS x; AS x ');(;;a;b' col;COPY zay FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x gs$y$k AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x gs$y$k AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = x'01;10' + (SELECT 1; SELECT 2) col "a;b" ,;
| FROM t WHERE a = --$x$;
|  AS x (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' + (SELECT 1; SELECT 2) col "a;b" ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' + (SELECT 1; SELECT 2) col "a;b" ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + /*/*/*end;select 1;*/*/*/ + CREATE FUNCTION uudz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b/*;;'; SELECT 1; END 1;
| , --select 1;
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + /*/*/*end;select 1;*/*/*/ + CREATE FUNCTION uudz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b/*;;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + /*/*/*end;select 1;*/*/*/ + CREATE FUNCTION uudz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b/*;;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col e'\\' SELECT;FROM t WHERE a = CREATE FUNCTION hmug() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;b'; SELECT 1; END 1 /*--;/ *;);(*/ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col e'\\' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col e'\\' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT "--;);(; " , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ";"
| position: 30
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "--;);(; " , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ";"
| position: 30
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT "--;);(; " , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near ";"
| position: 30
> exec
| SELECT $$;$$ AS x --$x$;
|  , x'01;10' FROM t WHERE a =; AS x (SELECT 1; SELECT 2) 1;
| COPY aot FROM stdin;
| a;b
| c;d;e
| \.
|
| + $$select 1;$$ col CREATE FUNCTION lnse() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;--;'; SELECT 1; END SELECT --END ;a;b/*;
|  1;
< error
| 42601 syntax error at or near ";"
| position: 54
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$;$$ AS x --$x$;
|  , x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 54
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$;$$ AS x --$x$;
|  , x'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 54
> exec
| 1 --a;b
|  , "END ;; " FROM t WHERE a =; -- trailing note;
| , --;;$x$;
|  col;
| SELECT /*/*/*; END ;*/*/*/ +; -- trailing note;
| AS x 'end;' , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --a;b
|  , "END ;; " FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --a;b
|  , "END ;; " FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col CREATE FUNCTION atjx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END SELECT /*END ;*/ col /*);(--;*/ +;COPY dxy FROM stdin;
| 1	\N
| a;b
| \.
| AS x ib$v$r AS x ee$w$p +; SELECT CREATE FUNCTION f_nx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;);(END ;'; END col $$;; $$ , /*END ;*/ AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION atjx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END SELECT /*END ;*/ col /*);(--;*/ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION atjx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; '; END SELECT /*END ;*/ col /*);(--;*/ +;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY lne FROM stdin;
| 1	\N
| a;b
| \.
|  -- trailing note;
| col "select 1;select 1;a;b" AS x;
| AS x "--;a;b" col ");($x$;--;" ,;
| col ");($x$;" col ');(a;b\' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lne FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lne FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + E'\n;\n;\n;' AS x ');(a;b;;' col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\n;\n;\n;' AS x ');(a;b;;' col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\n;\n;\n;' AS x ');(a;b;;' col (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| FROM t WHERE a = E'\';'';\\' FROM t WHERE a = CREATE FUNCTION qalc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;bend;'; SELECT 1; END col li$_$n col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E'\';'';\\' FROM t WHERE a = CREATE FUNCTION qalc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;bend;'; SELECT 1; END col li$_$n col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\';'';\\' FROM t WHERE a = CREATE FUNCTION qalc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$a;bend;'; SELECT 1; END col li$_$n col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| SELECT U&'; --;$x$;' col U&'; $x$;end;' col uq$r$k AS x; -- trailing note;
| COPY o_b FROM stdin;
| c;d;e
| 1	\N
| \.
|
| col $hi$$x$;);(--;$hi$ , U&'$;$' 1 $$END ;; ;$$ 1;
< error
| 42601 syntax error at or near "U&'; $x$;end;'"
| position: 26
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'; --;$x$;' col U&'; $x$;end;' col uq$r$k AS x;
< error
| 42601 syntax error at or near "U&'; $x$;end;'"
| position: 26
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'; --;$x$;' col U&'; $x$;end;' col uq$r$k AS x;
< error
| 42601 syntax error at or near "U&'; $x$;end;'"
| position: 26
> exec
| AS x b'01;10' col $xm_n$;/*;$xm_n$ ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x b'01;10' col $xm_n$;/*;$xm_n$ ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x b'01;10' col $xm_n$;/*;$xm_n$ ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , e''';'';' col CREATE FUNCTION fpij() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;a;b'; SELECT 1; END + ----;a;b; 
|  SELECT; -- trailing note;
| , /*/*;select 1;*/*/ 1 --a;b/*;END ;
|  , --end;$x$;
|  SELECT; SELECT --$x$;/*;);(
|  + x'01;10' FROM t WHERE a = $$;END ;$$ SELECT;1 'a;b$;$;;\' , U&'END ;;;' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , e''';'';' col CREATE FUNCTION fpij() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;a;b'; SELECT 1; END + ----;a;b; 
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , e''';'';' col CREATE FUNCTION fpij() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;a;b'; SELECT 1; END + ----;a;b; 
|  SELECT;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY ruu FROM stdin;
| a;b
| c;d;e
| \.
|
| AS x 'select 1;' AS x B'01;10' col x'01;10' +; -- trailing note;
| COPY bhq FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ruu FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ruu FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col b'01;10' SELECT cs$b$k FROM t WHERE a = (SELECT 1; SELECT 2) AS x;, CREATE FUNCTION uorx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;$x$;'; END ,;
| COPY mwo FROM stdin;
| c;d;e
| 1	\N
| \.
|
| , fy$v$t col xv$i$r col ");(;;" ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col b'01;10' SELECT cs$b$k FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col b'01;10' SELECT cs$b$k FROM t WHERE a = (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT CREATE FUNCTION kypz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END + (SELECT 1; SELECT 2) + 'select 1;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION kypz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END + (SELECT 1; SELECT 2) + 'select 1;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION kypz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END + (SELECT 1; SELECT 2) + 'select 1;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| FROM t WHERE a = '/*;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = '/*;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = '/*;' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = e'\n;\';' +;AS x b'01;10' AS x;
| COPY apo FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = e'\n;\';' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = e'\n;\';' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 --select 1;
|  AS x E'\';\n;' SELECT;
| col --; 
|  col xn$a$o +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 --select 1;
|  AS x E'\';\n;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 --select 1;
|  AS x E'\';\n;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| SELECT dm$g$r FROM t WHERE a = $jjm$;;END ;$jjx$$jjm$ 1;COPY rkb FROM stdin;
| 1	\N
| a;b
| \.
|  FROM t WHERE a = CREATE FUNCTION _kb_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$END ;end;'; END ,;
< error
| 42601 syntax error at or near "1"
| position: 55
> exec
| BEGIN
< ok 0 -
> exec
| SELECT dm$g$r FROM t WHERE a = $jjm$;;END ;$jjx$$jjm$ 1;
< error
| 42601 syntax error at or near "1"
| position: 55
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT dm$g$r FROM t WHERE a = $jjm$;;END ;$jjx$$jjm$ 1;
< error
| 42601 syntax error at or near "1"
| position: 55
> exec
| SELECT CREATE FUNCTION zlre() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;'; SELECT 1; END , "select 1;/*;" SELECT CREATE FUNCTION cjyh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END +;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION zlre() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;'; SELECT 1; END , "select 1;/*;" SELECT CREATE FUNCTION cjyh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END +;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION zlre() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;select 1;'; SELECT 1; END , "select 1;/*;" SELECT CREATE FUNCTION cjyh() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$'; END +;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| 1 /*/*/ *;;;*/*/ SELECT "; " FROM t WHERE a = vk$o$i FROM t WHERE a =;
| AS x "/*;" , (SELECT 1; SELECT 2) + (SELECT 1; SELECT 2) 1;
| + $wyi$--;$wyi$ 1 /*/*/*end;;;);(*/*/*/ +; FROM t WHERE a = ly$o$m 1 CREATE FUNCTION kimu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;/*;'; SELECT 1; END FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/ *;;;*/*/ SELECT "; " FROM t WHERE a = vk$o$i FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/ *;;;*/*/ SELECT "; " FROM t WHERE a = vk$o$i FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY kxm FROM stdin;
| c;d;e
| a;b
| \.
| , /*/*select 1;select 1;$;$*/*/ + ";;$;$--;" FROM t WHERE a =;
| + --; ; a;b
|  FROM t WHERE a = $b$$x$;$b$ FROM t WHERE a = CREATE FUNCTION ujkz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; END col;
| SELECT $bcmw$);(end;$bcmx$$bcmw$ AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kxm FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kxm FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| col "a;b" FROM t WHERE a =; -- trailing note;
| AS x U&'END ;' +; , "a;bend;" col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "a;b" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "a;b" FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| , CREATE FUNCTION mhqq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';select 1;select 1;'; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION mhqq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';select 1;select 1;'; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION mhqq() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';select 1;select 1;'; END +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , ");(" 1 CREATE FUNCTION txiu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END FROM t WHERE a = U&';;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ");(" 1 CREATE FUNCTION txiu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END FROM t WHERE a = U&';;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ");(" 1 CREATE FUNCTION txiu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END FROM t WHERE a = U&';;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT "; " + gu$f$x col;
< error
| 42703 column "; " does not exist
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT "; " + gu$f$x col;
< error
| 42703 column "; " does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yvw FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY yvw FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY yvw FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY auo FROM stdin;
| 1	\N
| a;b
| \.
|
| 1 $wk_m$select 1;$wk_x$$wk_m$ 1;
| , /*/*/*; );(*/*/*/ col ';;select 1;' AS x;
| FROM t WHERE a = '$;$$x$;\' SELECT ");(select 1;" col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY auo FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY auo FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x /*end;*/ SELECT E'\\\n;'';' SELECT $$--;;;$$ 1;1 'a;b/*;' FROM t WHERE a = '--;\' AS x; AS x ay$r$n + /*/*end;end;*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*end;*/ SELECT E'\\\n;'';' SELECT $$--;;;$$ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*end;*/ SELECT E'\\\n;'';' SELECT $$--;;;$$ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY tlm FROM stdin;
| c;d;e
| 1	\N
| \.
| COPY _is FROM stdin;
| a;b
| a;b
| \.
|
| 1 ----;end;);(
|  FROM t WHERE a = $$;;);(; $$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY tlm FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY tlm FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , --$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , --$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , --$;$
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 /*/*/*end;*/*/*/ FROM t WHERE a =;
| col 'select 1;\' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*/*/*end;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*/*/*end;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 "end;a;ba;b" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "end;a;ba;b" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "end;a;ba;b" FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x "select 1;" col X'01;10' FROM t WHERE a = $$);($$ +;
| FROM t WHERE a = --END ;/*;
|  AS x;col --);(end;
|  FROM t WHERE a =; -- trailing note;
| AS x $g$end;end;$g$ , --END ;select 1;;
|  AS x $$--;$x$;;;$$ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x "select 1;" col X'01;10' FROM t WHERE a = $$);($$ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x "select 1;" col X'01;10' FROM t WHERE a = $$);($$ +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = x'01;10' SELECT $$; $;$;;$$ FROM t WHERE a =; -- trailing note;
| , /*; */ , yt$l$p , CREATE FUNCTION wesd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; END col;FROM t WHERE a = U&'END ;$;$' col; SELECT --/*;/*;
|  , "$x$;a;b" 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' SELECT $$; $;$;;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = x'01;10' SELECT $$; $;$;;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col $$;;--;;$$ SELECT E'\';\\' FROM t WHERE a =; -- trailing note;
| COPY qyv FROM stdin;
| 1	\N
| c;d;e
| \.
|
| COPY _vb FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$;;--;;$$ SELECT E'\';\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$;;--;;$$ SELECT E'\';\\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| 1 ');(/*;);(' SELECT;
| SELECT /*;*/ AS x U&';' AS x /*/*/*select 1;--;/ *;*/*/*/ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 ');(/*;);(' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 ');(/*;);(' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x U&');(' SELECT pe$c$v SELECT E'\';\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&');(' SELECT pe$c$v SELECT E'\';\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&');(' SELECT pe$c$v SELECT E'\';\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = X'01;10' col ha$b$b , hh$o$j FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = X'01;10' col ha$b$b , hh$o$j FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = X'01;10' col ha$b$b , hh$o$j FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| col --end;
|  SELECT;
| FROM t WHERE a = gt$x$c SELECT ');($x$;end;' SELECT $$);(a;b$$ ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col --end;
|  SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col --end;
|  SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| FROM t WHERE a = "END ;select 1;; " FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "END ;select 1;; " FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "END ;select 1;; " FROM t WHERE a = (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = /*/*$;$*/*/ col;
| + ck$_$m 1 CREATE FUNCTION pmuo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;);(/*;'; END + U&'END ;' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*/*$;$*/*/ col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*/*$;$*/*/ col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + (SELECT 1; SELECT 2) AS x;
| COPY wqk FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY sh_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
| + /*end;;*/ + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY sh_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY sh_ FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| 1 /*;;END ;*/ SELECT --; 
|  ,;
| FROM t WHERE a = "select 1;$;$; " col e'\';\';' 1 $jlnm$--;/*;$jlnm$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 /*;;END ;*/ SELECT --; 
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 /*;;END ;*/ SELECT --; 
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY wrs FROM stdin;
| 1	\N
| a;b
| \.
|  -- trailing note;
| col $$$x$;;;$$ AS x; col /*END ;*/ 1 $$;;--;$;$$$ + U&';;END ;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY wrs FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY wrs FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| FROM t WHERE a = fh$r$f 1 ";" SELECT e'\\\\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = fh$r$f 1 ";" SELECT e'\\\\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = fh$r$f 1 ";" SELECT e'\\\\' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY psv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY psv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY psv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| , /*end;*/ +;
| AS x ek$p$_ +; AS x '$;$' AS x B'01;10' 1; -- trailing note;
| SELECT $_p$a;b$_p$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*end;*/ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*end;*/ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| 1 da$i$t 1 $$;;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 da$i$t 1 $$;;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 da$i$t 1 $$;;$$ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY lzt FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lzt FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lzt FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| FROM t WHERE a = 'select 1;END ;;' FROM t WHERE a =;
| SELECT U&');($x$;' + me$z$v FROM t WHERE a = "$;$" col;
| SELECT X'01;10' SELECT E'\';\n;' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = 'select 1;END ;;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = 'select 1;END ;;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY pjt FROM stdin;
| 1	\N
| c;d;e
| \.
| col (SELECT 1; SELECT 2) , CREATE FUNCTION mhmv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;select 1;;;'; END SELECT --;end;; 
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY pjt FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY pjt FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x U&';;--;' AS x; -- trailing note;
| col E'\';\';\n;' 1;
| 1 'select 1;END ;\' 1 ';;);(' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&';;--;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&';;--;' AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = $$);(;;$x$;$$ FROM t WHERE a = /*/*/*; end;select 1;*/*/*/ +;
| FROM t WHERE a = b'01;10' + /*/*/*end;end;*/*/*/ FROM t WHERE a =; col X'01;10' FROM t WHERE a = /*/*/*end;a;b);(*/*/*/ ,; -- trailing note;
| col e'\';\';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $$);(;;$x$;$$ FROM t WHERE a = /*/*/*; end;select 1;*/*/*/ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $$);(;;$x$;$$ FROM t WHERE a = /*/*/*; end;select 1;*/*/*/ +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + (SELECT 1; SELECT 2) col vx$j$v col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) col vx$j$v col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) col vx$j$v col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x U&'a;b;' SELECT;
| + $ntk$;$ntx$$ntk$ AS x --; ;;END ;
|  1 '$;$' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'a;b;' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'a;b;' SELECT;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + "$;$$;$" AS x 'end;' + CREATE FUNCTION wbxi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;);('; SELECT 1; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + "$;$$;$" AS x 'end;' + CREATE FUNCTION wbxi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;);('; SELECT 1; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + "$;$$;$" AS x 'end;' + CREATE FUNCTION wbxi() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;);('; SELECT 1; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| SELECT U&'select 1;--;; ' FROM t WHERE a = U&'/*;' 1;
< error
| 42601 syntax error at or near "1"
| position: 52
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&'select 1;--;; ' FROM t WHERE a = U&'/*;' 1;
< error
| 42601 syntax error at or near "1"
| position: 52
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&'select 1;--;; ' FROM t WHERE a = U&'/*;' 1;
< error
| 42601 syntax error at or near "1"
| position: 52
> exec
| FROM t WHERE a = "a;b;;" , (SELECT 1; SELECT 2) +; col sq$w$l +;, b'01;10' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "a;b;;" , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "a;b;;" , (SELECT 1; SELECT 2) +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY odk FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY odk FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY odk FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT $yv$--;);($yv$ + cl$d$k , x'01;10' SELECT;
| , "$x$;" 1;
| COPY qwl FROM stdin;
| a;b
| c;d;e
| \.
|
| COPY ncz FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near ","
| position: 51
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $yv$--;);($yv$ + cl$d$k , x'01;10' SELECT;
< error
| 42703 column "cl$d$k" does not exist
| position: 25
> exec
| ROLLBACK
< ok 0 -
> exec
|
| , "$x$;" 1;
< error
| 42601 syntax error at or near ","
| position: 2
> exec
| COPY jaf FROM stdin;
| 1	\N
| 1	\N
| \.
| SELECT E'\n;'';' FROM t WHERE a = --$x$;select 1;; 
|  AS x; -- trailing note;
| COPY nbn FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY jaf FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY jaf FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| , $eu$/*;;$eu$ ,;+ $$;$$ col X'01;10' , CREATE FUNCTION hkup() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;--;'; SELECT 1; END AS x;
| 1 CREATE FUNCTION fefy() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;end;'; END 1 $$;a;b$$ col e'\';\n;' 1;
| col /*$x$;;;);(*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $eu$/*;;$eu$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $eu$/*;;$eu$ ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , CREATE FUNCTION bcxr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END FROM t WHERE a = $n$;$n$ +; -- trailing note;
| , (SELECT 1; SELECT 2) col /*/*/*);(/ *;*/*/*/ AS x E'\\' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , CREATE FUNCTION bcxr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END FROM t WHERE a = $n$;$n$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , CREATE FUNCTION bcxr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;'; END FROM t WHERE a = $n$;$n$ +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| , /*/*END ;/ *;$;$*/*/ , '; $;$select 1;\' + b'01;10' col; COPY ayt FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*END ;/ *;$;$*/*/ , '; $;$select 1;\' + b'01;10' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*END ;/ *;$;$*/*/ , '; $;$select 1;\' + b'01;10' col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY nly FROM stdin;
| c;d;e
| c;d;e
| \.
|  COPY gcm FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY nly FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY nly FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY quf FROM stdin;
| a;b
| c;d;e
| \.
|
| FROM t WHERE a = bo$w$q SELECT x'01;10' col;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY quf FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY quf FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT ni$w$i + ----;END ;
|  1;
| , wp$u$i col ----;;;
|  AS x;
| FROM t WHERE a = B'01;10' , U&'END ;' + jx$_$b ,;
| SELECT "a;b; $x$;" SELECT (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near ","
| position: 32
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ni$w$i + ----;END ;
|  1;
< error
| 42703 column "ni$w$i" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| , wp$u$i col ----;;;
|  AS x;
< error
| 42601 syntax error at or near ","
| position: 2
> exec
| AS x /*/*$;$END ;);(*/*/ AS x /*/*/*);(*/*/*/ col;
| col E''';\';'';' SELECT E'\\' AS x; col "end;a;b$x$;" FROM t WHERE a =;
| COPY jwr FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*$;$END ;);(*/*/ AS x /*/*/*);(*/*/*/ col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*$;$END ;);(*/*/ AS x /*/*/*);(*/*/*/ col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY oiv FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY oiv FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY oiv FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| + --;;/*;
|  ,;
| AS x CREATE FUNCTION hqkg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;a;b;;'; END col U&'select 1;);(' AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + --;;/*;
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + --;;/*;
|  ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| col /*/*/*--;*/*/*/ FROM t WHERE a = (SELECT 1; SELECT 2) ,;
| COPY pod FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col /*/*/*--;*/*/*/ FROM t WHERE a = (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col /*/*/*--;*/*/*/ FROM t WHERE a = (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col 'END ;;--;' 1; -- trailing note;
| FROM t WHERE a = (SELECT 1; SELECT 2) FROM t WHERE a = --;
|  FROM t WHERE a =;
| AS x e''';'';\';' , /*/*;;$x$;;;*/*/ SELECT;FROM t WHERE a = bv$e$g FROM t WHERE a = E'\';' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col 'END ;;--;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col 'END ;;--;' 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| AS x ro$x$d FROM t WHERE a = $$;;/*;; $$ AS x;
| COPY rli FROM stdin;
| a;b
| 1	\N
| \.
| , b'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x ro$x$d FROM t WHERE a = $$;;/*;; $$ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x ro$x$d FROM t WHERE a = $$;;/*;; $$ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| + b'01;10' + $$END ;;;select 1;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + b'01;10' + $$END ;;;select 1;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + b'01;10' + $$END ;;;select 1;$$ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY agk FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY vov FROM stdin;
| 1	\N
| 1	\N
| \.
|
| + --$x$;end;a;b
|  + --$;$
|  FROM t WHERE a = E'\\\';'';' +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY agk FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY agk FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY hgv FROM stdin;
| a;b
| a;b
| \.
|
| COPY szp FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hgv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hgv FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| 1 nb$d$r SELECT E'\n;\';\n;' AS x; -- trailing note;
| COPY hnj FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 nb$d$r SELECT E'\n;\';\n;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 nb$d$r SELECT E'\n;\';\n;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = CREATE FUNCTION cbxb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;$x$;'; END FROM t WHERE a = (SELECT 1; SELECT 2) + CREATE FUNCTION jnov() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;$x$;select 1;'; SELECT 1; END AS x;
| FROM t WHERE a = __$t$e col U&'; select 1;' SELECT /*/*/*$x$;*/*/*/ FROM t WHERE a =;
| SELECT E'\';' SELECT ';;; \' FROM t WHERE a = --$x$;select 1;
|  FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION cbxb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;$x$;'; END FROM t WHERE a = (SELECT 1; SELECT 2) + CREATE FUNCTION jnov() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;$x$;select 1;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION cbxb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;$x$;'; END FROM t WHERE a = (SELECT 1; SELECT 2) + CREATE FUNCTION jnov() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;$x$;select 1;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = U&'; $;$$;$' +;
| + $$);(end;$$ , x'01;10' 1; COPY rbc FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = U&'; $;$$;$' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = U&'; $;$$;$' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x $$);($$ ,; COPY mhp FROM stdin;
| 1	\N
| 1	\N
| \.
| col (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $$);($$ ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $$);($$ ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| FROM t WHERE a = "$;$;;" col e'\\' +;
| COPY xnu FROM stdin;
| c;d;e
| 1	\N
| \.
|  -- trailing note;
| + CREATE FUNCTION gtwc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$);('; SELECT 1; END , /*/*$x$;$;$*/*/ col U&'a;ba;b' ,;
| COPY wgo FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = "$;$;;" col e'\\' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = "$;$;;" col e'\\' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , /*/*$x$;select 1;*/*/ 1 /*;;*/ AS x ");();(;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , /*/*$x$;select 1;*/*/ 1 /*;;*/ AS x ");();(;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , /*/*$x$;select 1;*/*/ 1 /*;;*/ AS x ");();(;" AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT /*; */ 1 /*/*/*);($x$;--;*/*/*/ AS x; -- trailing note;
| COPY tog FROM stdin;
| a;b
| 1	\N
| \.
|
| FROM t WHERE a = CREATE FUNCTION qpom() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;);('; END , wu$b$a +;
| + B'01;10' 1 (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "a"
| position: 85
> exec
| BEGIN
< ok 0 -
> exec
| SELECT /*; */ 1 /*/*/*);($x$;--;*/*/*/ AS x;
< ok 1 -
> exec
|  -- trailing note;
| COPY tog FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 41
> exec
| ROLLBACK
< ok 0 -
> exec
|  -- trailing note;
| COPY tog FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 41
> exec
| FROM t WHERE a = CREATE FUNCTION hytg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$x$;$x$;'; SELECT 1; END FROM t WHERE a = /*;;;$;$*/ ,;
| COPY dbb FROM stdin;
| c;d;e
| a;b
| \.
|  AS x U&'end;/*;;;' ,;
| COPY wwu FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION hytg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$x$;$x$;'; SELECT 1; END FROM t WHERE a = /*;;;$;$*/ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = CREATE FUNCTION hytg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$x$;$x$;'; SELECT 1; END FROM t WHERE a = /*;;;$;$*/ ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x CREATE FUNCTION hbhj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END , X'01;10' , vv$m$u +;
| 1 CREATE FUNCTION qbak() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$x$;'; SELECT 1; END ,;
| 1 'END ;\' SELECT "a;b;;;" FROM t WHERE a = /*);($x$;*/ col;
| col U&'END ;;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION hbhj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END , X'01;10' , vv$m$u +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION hbhj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END , X'01;10' , vv$m$u +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT e'\\\n;\n;' + /*/*);(a;b*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 37
> exec
| BEGIN
< ok 0 -
> exec
| SELECT e'\\\n;\n;' + /*/*);(a;b*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 37
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT e'\\\n;\n;' + /*/*);(a;b*/*/ ,;
< error
| 42601 syntax error at or near ","
| position: 37
> exec
| + '--;;' AS x --select 1;;;
|  1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + '--;;' AS x --select 1;;;
|  1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + '--;;' AS x --select 1;;;
|  1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| 1 "a;b" +;
| AS x CREATE FUNCTION omaj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;END ;'; SELECT 1; END ,;
| COPY ynv FROM stdin;
| 1	\N
| a;b
| \.
|
| FROM t WHERE a = --select 1;
|  + x'01;10' , b'01;10' 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "a;b" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "a;b" +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = /*$x$;*/ AS x;
| 1 CREATE FUNCTION tbcz() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;;'; SELECT 1; END SELECT $$select 1;; $$ FROM t WHERE a =;
| SELECT ";$x$;" , CREATE FUNCTION kp_h() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;$x$;/*;'; SELECT 1; END col; COPY pzz FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*$x$;*/ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*$x$;*/ AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| + x'01;10' ,;
| COPY snt FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + x'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + x'01;10' ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , (SELECT 1; SELECT 2) AS x (SELECT 1; SELECT 2) 1;
| SELECT (SELECT 1; SELECT 2) SELECT CREATE FUNCTION _icb() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; SELECT 1; END AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) AS x (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , (SELECT 1; SELECT 2) AS x (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| COPY agk FROM stdin;
| a;b
| a;b
| \.
|
| + CREATE FUNCTION ugib() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; END +;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY agk FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY agk FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT U&';' AS x X'01;10' SELECT ox$t$y SELECT; -- trailing note;
| col X'01;10' FROM t WHERE a = e'\';'';\';' + ');(END ;\' +;
| FROM t WHERE a = E''';'';' AS x x'01;10' col E'\';' col;
| COPY duj FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "X'01;10'"
| position: 19
> exec
| BEGIN
< ok 0 -
> exec
| SELECT U&';' AS x X'01;10' SELECT ox$t$y SELECT;
< error
| 42601 syntax error at or near "X'01;10'"
| position: 19
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT U&';' AS x X'01;10' SELECT ox$t$y SELECT;
< error
| 42601 syntax error at or near "X'01;10'"
| position: 19
> exec
| , ');($x$;' FROM t WHERE a = CREATE FUNCTION goit() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;; '; END 1 (SELECT 1; SELECT 2) AS x;
| COPY ldy FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , ');($x$;' FROM t WHERE a = CREATE FUNCTION goit() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;; '; END 1 (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , ');($x$;' FROM t WHERE a = CREATE FUNCTION goit() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;; '; END 1 (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x bp$_$y ,; -- trailing note;
| AS x 'a;b' FROM t WHERE a =; -- trailing note;
| COPY fxn FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x bp$_$y ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x bp$_$y ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY srj FROM stdin;
| 1	\N
| c;d;e
| \.
| FROM t WHERE a = --;; ;
|  AS x;
| col B'01;10' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY srj FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY srj FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| AS x b'01;10' AS x E''';' + (SELECT 1; SELECT 2) col;+ (SELECT 1; SELECT 2) FROM t WHERE a = CREATE FUNCTION gzvj() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;$;$select 1;'; SELECT 1; END , 'END ;;END ;' col; 1 --;;$x$;$;$
|  col;
| COPY ebe FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x b'01;10' AS x E''';' + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x b'01;10' AS x E''';' + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| SELECT CREATE FUNCTION pvld() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END 1 /*/*/*a;b;;--;*/*/*/ col $$select 1;;$$ col;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| BEGIN
< ok 0 -
> exec
| SELECT CREATE FUNCTION pvld() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END 1 /*/*/*a;b;;--;*/*/*/ col $$select 1;;$$ col;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT CREATE FUNCTION pvld() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;'; SELECT 1; END 1 /*/*/*a;b;;--;*/*/*/ col $$select 1;;$$ col;
< error
| 42601 syntax error at or near "CREATE"
| position: 8
> exec
| FROM t WHERE a = E'\\\\' , /*/*select 1;$x$;$;$*/*/ 1;
| AS x --;;);(;
|  1 CREATE FUNCTION vncx() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';'; SELECT 1; END SELECT U&');(' col;COPY csq FROM stdin;
| a;b
| a;b
| \.
|
| AS x e'\';' SELECT e'\\\';' 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = E'\\\\' , /*/*select 1;$x$;$;$*/*/ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = E'\\\\' , /*/*select 1;$x$;$;$*/*/ 1;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| AS x U&'a;b; ' + E'\';\\'';' ,;, '$x$;' 1;
| , --);(
|  , (SELECT 1; SELECT 2) 1 --end;--;$x$;
|  ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&'a;b; ' + E'\';\\'';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&'a;b; ' + E'\';\\'';' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY fpe FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| + U&';;; ' SELECT --);(END ;$;$
|  ,;COPY oij FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY fpe FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY fpe FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY tgs FROM stdin;
| a;b
| a;b
| \.
|
| col CREATE FUNCTION seaa() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$x$;'; END 1;
| , E'\\\\' 1 $$select 1;);($x$;$$ SELECT CREATE FUNCTION jvok() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;ba;b; '; SELECT 1; END AS x;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY tgs FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY tgs FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = /*/*select 1;*/*/ AS x E''';' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*/*select 1;*/*/ AS x E''';' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*/*select 1;*/*/ AS x E''';' col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 CREATE FUNCTION a_mo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;'; SELECT 1; END col 'end;$;$' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION a_mo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;'; SELECT 1; END col 'end;$;$' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION a_mo() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$select 1;'; SELECT 1; END col 'end;$;$' col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) col --;;$;$
|  SELECT;
| AS x CREATE FUNCTION revk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$;$'; END AS x aa$p$l SELECT CREATE FUNCTION _dik() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;'; SELECT 1; END SELECT;
| 1 $$a;b$$ , hp$l$g 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) col --;;$;$
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) col --;;$;$
|  SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY uju FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY uju FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY uju FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col $$end;; --;$$ + E'\';' AS x e'\n;'';' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col $$end;; --;$$ + E'\';' AS x e'\n;'';' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col $$end;; --;$$ + E'\';' AS x e'\n;'';' SELECT;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| + (SELECT 1; SELECT 2) AS x CREATE FUNCTION _sep() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$END ;END ;'; SELECT 1; END col; -- trailing note;
| FROM t WHERE a = $p_$/*;end;END ;$p_$ FROM t WHERE a = ow$a$u AS x (SELECT 1; SELECT 2) col;FROM t WHERE a = /*select 1;*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) AS x CREATE FUNCTION _sep() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$END ;END ;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + (SELECT 1; SELECT 2) AS x CREATE FUNCTION _sep() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$END ;END ;'; SELECT 1; END col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY ubc FROM stdin;
| c;d;e
| 1	\N
| \.
|
| COPY hvo FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| col "END ;select 1;END ;" col (SELECT 1; SELECT 2) +;AS x qn$z$d AS x 'select 1;' +;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ubc FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ubc FROM stdin;
| c;d;e
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| , sn$s$t AS x e'\';\';\n;' AS x U&';;' ,;
| + /*END ;;;*/ AS x X'01;10' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , sn$s$t AS x e'\';\';\n;' AS x U&';;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , sn$s$t AS x e'\';\';\n;' AS x U&';;' ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| col U&'select 1;$x$;select 1;' AS x me$z$t FROM t WHERE a =;
| + U&'a;b;;' SELECT ----;--;);(
|  SELECT "; " FROM t WHERE a =;
| COPY jbh FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col U&'select 1;$x$;select 1;' AS x me$z$t FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col U&'select 1;$x$;select 1;' AS x me$z$t FROM t WHERE a =;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY mbv FROM stdin;
| 1	\N
| 1	\N
| \.
|
| col U&'$x$;);(;;' AS x X'01;10' FROM t WHERE a =;
| SELECT $$END ;;$$ SELECT E'\\' col U&'select 1;);(' +;
| COPY _tf FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY mbv FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY mbv FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + ak$k$x + /*/*/*$x$;);(*/*/*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + ak$k$x + /*/*/*$x$;);(*/*/*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + ak$k$x + /*/*/*$x$;);(*/*/*/ AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| AS x $e$END ;;;$x$$e$ 1 $w$$x$;/*;$x$$w$ 1;
| SELECT ");(" col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $e$END ;;;$x$$e$ 1 $w$$x$;/*;$x$$w$ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $e$END ;;;$x$$e$ 1 $w$$x$;/*;$x$$w$ 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY ebu FROM stdin;
| 1	\N
| 1	\N
| \.
|  COPY qlc FROM stdin;
| c;d;e
| c;d;e
| \.
|
| AS x "--;" FROM t WHERE a = --;;--;
|  col; , CREATE FUNCTION pfem() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '$;$);(;'; SELECT 1; END SELECT;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ebu FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ebu FROM stdin;
| 1	\N
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| COPY eaf FROM stdin;
| a;b
| a;b
| \.
|
| SELECT b'01;10' col $$END ;end;/*;$$ , CREATE FUNCTION qqfk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b'; END 1;
| 1 B'01;10' AS x;
| SELECT '$x$;' FROM t WHERE a = U&'/*;end;; ' ,;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY eaf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY eaf FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| FROM t WHERE a = /*/*$;$*/*/ SELECT;
| SELECT e'\n;' SELECT U&'/*;; end;' 1; -- trailing note;
| COPY qcy FROM stdin;
| c;d;e
| a;b
| \.
|
| FROM t WHERE a = --a;b);(
|  , E'\n;' +;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = /*/*$;$*/*/ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = /*/*$;$*/*/ SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| , "$x$;" + U&'end;;;' FROM t WHERE a = --;;
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , "$x$;" + U&'end;;;' FROM t WHERE a = --;;
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , "$x$;" + U&'end;;;' FROM t WHERE a = --;;
|  +;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| SELECT $lmel$select 1;select 1;;;$lmex$$lmel$ FROM t WHERE a =; COPY oei FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| , /*end;$x$;);(*/ AS x b'01;10' +;
| SELECT (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near ";"
| position: 63
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $lmel$select 1;select 1;;;$lmex$$lmel$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 63
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $lmel$select 1;select 1;;;$lmex$$lmel$ FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 63
> exec
| COPY bwj FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| COPY rmj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY bwj FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY bwj FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x U&';' SELECT --select 1;$;$$;$
|  , --a;b; ;
|  +; , E'\n;'';' FROM t WHERE a = --;;
|  , E''';\\' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x U&';' SELECT --select 1;$;$$;$
|  , --a;b; ;
|  +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x U&';' SELECT --select 1;$;$$;$
|  , --a;b; ;
|  +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY _iz FROM stdin;
| c;d;e
| c;d;e
| \.
|
| FROM t WHERE a = --; $;$; 
|  AS x U&'$;$$x$;' AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY _iz FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY _iz FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| SELECT E'\\' + (SELECT 1; SELECT 2) FROM t WHERE a = /*end;$x$;/ *;*/ FROM t WHERE a =;
| AS x $lu$$x$;);($lx$$lu$ col;COPY nwv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near ";"
| position: 25
> exec
| BEGIN
< ok 0 -
> exec
| SELECT E'\\' + (SELECT 1; SELECT 2) FROM t WHERE a = /*end;$x$;/ *;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 25
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT E'\\' + (SELECT 1; SELECT 2) FROM t WHERE a = /*end;$x$;/ *;*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near ";"
| position: 25
> exec
| 1 (SELECT 1; SELECT 2) , $dpju$$x$;$dpjx$$dpju$ +;
| FROM t WHERE a = "end;END ;" , $hqfw$select 1;a;b$x$;$hqfx$$hqfw$ 1 /*; ; $x$;*/ FROM t WHERE a =; -- trailing note;
| FROM t WHERE a = X'01;10' AS x --end;
|  +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) , $dpju$$x$;$dpjx$$dpju$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) , $dpju$$x$;$dpjx$$dpju$ +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x (SELECT 1; SELECT 2) , CREATE FUNCTION reik() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;;;'; SELECT 1; END 1; -- trailing note;
| SELECT (SELECT 1; SELECT 2) FROM t WHERE a =;
| 1 "/*;$x$;" + --; $;$
|  col; -- trailing note;
| AS x $$select 1;$;@;$$$ , rz$s$y + mr$d$l +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) , CREATE FUNCTION reik() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;;;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x (SELECT 1; SELECT 2) , CREATE FUNCTION reik() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; --;;;'; SELECT 1; END 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY akf FROM stdin;
| a;b
| c;d;e
| \.
|  AS x (SELECT 1; SELECT 2) , x'01;10' + "--;" ,;
| col U&'$;$' AS x;1 '$;$--;' FROM t WHERE a = E'\n;\n;\n;' FROM t WHERE a =;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY akf FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY akf FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| col ----;
|  , (SELECT 1; SELECT 2) 1;+ qy$e$u col rw$t$a FROM t WHERE a = $llc$--;$x$;$llc$ FROM t WHERE a =;
| , (SELECT 1; SELECT 2) FROM t WHERE a = "; END ;" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col ----;
|  , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col ----;
|  , (SELECT 1; SELECT 2) 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col CREATE FUNCTION pkgk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $x$;'; END AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col CREATE FUNCTION pkgk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $x$;'; END AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col CREATE FUNCTION pkgk() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; $x$;'; END AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| col "select 1;);(;;" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "select 1;);(;;" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "select 1;);(;;" 1;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT yy$b$p , E'\n;' 1 ");(select 1;" 1;1 e'\';'';' + "$;$--;$;$" ,;
< error
| 42601 syntax error at or near "1"
| position: 24
> exec
| BEGIN
< ok 0 -
> exec
| SELECT yy$b$p , E'\n;' 1 ");(select 1;" 1;
< error
| 42601 syntax error at or near "1"
| position: 24
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT yy$b$p , E'\n;' 1 ");(select 1;" 1;
< error
| 42601 syntax error at or near "1"
| position: 24
> exec
| COPY eqw FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY eqw FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY eqw FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ktv FROM stdin;
| 1	\N
| c;d;e
| \.
|  -- trailing note;
| SELECT ej$r$h FROM t WHERE a = $$);($$ FROM t WHERE a =;
| COPY egr FROM stdin;
| 1	\N
| a;b
| \.
|
| + x'01;10' col U&';$;$' , U&'END ;' AS x;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ktv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ktv FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| col "select 1;" col;, U&'$x$;);(' 1 /*;END ;*/ AS x; COPY sqh FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col "select 1;" col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col "select 1;" col;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| COPY hef FROM stdin;
| c;d;e
| c;d;e
| \.
|
| COPY orq FROM stdin;
| 1	\N
| a;b
| \.
|  -- trailing note;
| 1 U&'end;--;end;' + (SELECT 1; SELECT 2) 1 b'01;10' 1;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY hef FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY hef FROM stdin;
| c;d;e
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| 1 CREATE FUNCTION qszc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; select 1;'; SELECT 1; END +;
| COPY fyg FROM stdin;
| c;d;e
| a;b
| \.
|  COPY cut FROM stdin;
| 1	\N
| 1	\N
| \.
|  -- trailing note;
| col U&';' , 'a;b);(select 1;' FROM t WHERE a = /*/*/*--;$;$*/*/*/ 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION qszc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; select 1;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION qszc() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '; select 1;'; SELECT 1; END +;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| AS x --END ;$x$;$x$;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x --END ;$x$;$x$;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x --END ;$x$;$x$;
|  col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 '/*;;\' FROM t WHERE a =;
| AS x --$x$;/*;
|  col 'a;bselect 1;\' SELECT ');(;;' SELECT;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 '/*;;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 '/*;;\' FROM t WHERE a =;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| COPY kip FROM stdin;
| a;b
| a;b
| \.
|
| FROM t WHERE a = X'01;10' , X'01;10' +;
| SELECT CREATE FUNCTION qxgu() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ';;$;$'; SELECT 1; END col --select 1;end;
|  1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY kip FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY kip FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY ghh FROM stdin;
| c;d;e
| a;b
| \.
|  -- trailing note;
| + $$$;$$$ AS x U&'/*;--;; ' 1 B'01;10' col; -- trailing note;
| FROM t WHERE a = oi$r$y AS x --; ;
|  col;COPY bdo FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ghh FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ghh FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| AS x CREATE FUNCTION imlp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; --;'; END + (SELECT 1; SELECT 2) col; -- trailing note;
| col oa$m$c + ";;;;);(" FROM t WHERE a = x'01;10' ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x CREATE FUNCTION imlp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; --;'; END + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x CREATE FUNCTION imlp() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'a;b; --;'; END + (SELECT 1; SELECT 2) col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| , 'END ;);(END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , 'END ;);(END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , 'END ;);(END ;' FROM t WHERE a =;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x b'01;10' , b'01;10' col; -- trailing note;
| 1 E''';' AS x U&'$x$;);(' 1 ");(;;;;" col; COPY dna FROM stdin;
| a;b
| 1	\N
| \.
|  col ma$y$b , (SELECT 1; SELECT 2) FROM t WHERE a = E'\n;'';' +;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x b'01;10' , b'01;10' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x b'01;10' , b'01;10' col;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| AS x e'\\'';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x e'\\'';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x e'\\'';' FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| COPY ddp FROM stdin;
| 1	\N
| a;b
| \.
| , CREATE FUNCTION se_q() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '/*;$x$;'; END AS x B'01;10' 1; -- trailing note;
| COPY nfi FROM stdin;
| c;d;e
| a;b
| \.
|  col --;
|  1 U&';;' ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ddp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ddp FROM stdin;
| 1	\N
| a;b
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| SELECT $y$--;;;$x$$y$ + $$; $$ AS x;
| AS x U&');();(' FROM t WHERE a = '--;a;b\' , tu$d$g AS x; COPY hks FROM stdin;
| a;b
| 1	\N
| \.
|  AS x CREATE FUNCTION zzob() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'END ;END ;'; SELECT 1; END 1 "a;b" FROM t WHERE a = kh$p$q ,;
< error
| 42601 syntax error at or near "AS"
| position: 38
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $y$--;;;$x$$y$ + $$; $$ AS x;
< error
| 42725 operator is not unique: unknown + unknown
| position: 23
| hint: Could not choose a best candidate operator. You might need to add explicit type casts.
> exec
| ROLLBACK
< ok 0 -
> exec
|
| AS x U&');();(' FROM t WHERE a = '--;a;b\' , tu$d$g AS x;
< error
| 42601 syntax error at or near "AS"
| position: 2
> exec
| + X'01;10' 1 "; " 1 CREATE FUNCTION lxz_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;/*;'; END 1;
| 1 $$select 1;--;;;$$ 1 CREATE FUNCTION oitf() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');($;$/*;'; SELECT 1; END AS x;
| 1 U&'; a;b/*;' SELECT U&'select 1;end;; ' , e''';'';\n;' col;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + X'01;10' 1 "; " 1 CREATE FUNCTION lxz_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;/*;'; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + X'01;10' 1 "; " 1 CREATE FUNCTION lxz_() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;/*;'; END 1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , $p$select 1;$x$$p$ col CREATE FUNCTION cdqd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; '; END ,;
| COPY gan FROM stdin;
| c;d;e
| c;d;e
| \.
|  -- trailing note;
| SELECT B'01;10' AS x e'\\\n;\\' AS x;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , $p$select 1;$x$$p$ col CREATE FUNCTION cdqd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; '; END ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , $p$select 1;$x$$p$ col CREATE FUNCTION cdqd() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT ');(; '; END ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| AS x $$;;$$ + (SELECT 1; SELECT 2) ,;
| + /*/*/*end;end;*/*/*/ FROM t WHERE a = --/*;--;$x$;
|  , --$x$;$;$
|  SELECT;
| COPY bxm FROM stdin;
| a;b
| 1	\N
| \.
|
| col xz$z$y 1;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x $$;;$$ + (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x $$;;$$ + (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 CREATE FUNCTION wodt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$;$; '; END 1 "$;$;;;;" 1; COPY fhl FROM stdin;
| 1	\N
| c;d;e
| \.
|
| 1 /*/*/*a;b--;);(*/*/*/ col;AS x --end;
|  AS x $egqb$$;$/*;--;$egqx$$egqb$ AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 CREATE FUNCTION wodt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$;$; '; END 1 "$;$;;;;" 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 CREATE FUNCTION wodt() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT '--;$;$; '; END 1 "$;$;;;;" 1;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| FROM t WHERE a = e'\\'';'';' ,;
| + "END ;;);(" SELECT 'END ;;\' SELECT;
| FROM t WHERE a = (SELECT 1; SELECT 2) + '$x$;; ' + "select 1;" FROM t WHERE a =;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = e'\\'';'';' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = e'\\'';'';' ,;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| FROM t WHERE a = $jbii$END ;END ;;$jbix$$jbii$ col $yhm$/*;$;$$yhx$$yhm$ col;
| COPY aqz FROM stdin;
| c;d;e
| a;b
| \.
|  col $$);($x$;a;b$$ ,; -- trailing note;
| + (SELECT 1; SELECT 2) 1 /*/*$;$a;b*/*/ + ko$w$m SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = $jbii$END ;END ;;$jbix$$jbii$ col $yhm$/*;$;$$yhx$$yhm$ col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = $jbii$END ;END ;;$jbix$$jbii$ col $yhm$/*;$;$$yhx$$yhm$ col;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| 1 (SELECT 1; SELECT 2) 1 --END ;END ;; 
|  , $k$$;$$x$;$x$$k$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) 1 --END ;END ;; 
|  , $k$$;$$x$;$x$$k$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 (SELECT 1; SELECT 2) 1 --END ;END ;; 
|  , $k$$;$$x$;$x$$k$ col;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| col X'01;10' FROM t WHERE a = ");(;END ;" AS x;
| COPY jem FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| col X'01;10' FROM t WHERE a = ");(;END ;" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| col X'01;10' FROM t WHERE a = ");(;END ;" AS x;
< error
| 42601 syntax error at or near "col"
| position: 1
> exec
| SELECT $$select 1;$;$/*;$$ 1 e'\';' + /*/*--;*/*/ SELECT;1 pw$c$u +;
< error
| 42601 syntax error at or near "1"
| position: 28
> exec
| BEGIN
< ok 0 -
> exec
| SELECT $$select 1;$;$/*;$$ 1 e'\';' + /*/*--;*/*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 28
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT $$select 1;$;$/*;$$ 1 e'\';' + /*/*--;*/*/ SELECT;
< error
| 42601 syntax error at or near "1"
| position: 28
> exec
| COPY noy FROM stdin;
| a;b
| 1	\N
| \.
|
| SELECT U&'$;$;END ;' 1;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY noy FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY noy FROM stdin;
| a;b
| 1	\N
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT (SELECT 1; SELECT 2) ,;
| FROM t WHERE a = ';;--;\' col;
| + --; $x$;
|  AS x;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| BEGIN
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| ROLLBACK
< ok 0 -
> exec
| SELECT (SELECT 1; SELECT 2) ,;
< error
| 42601 syntax error at or near ";"
| position: 17
> exec
| COPY ryb FROM stdin;
| c;d;e
| a;b
| \.
|
| COPY gea FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ryb FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ryb FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| COPY ake FROM stdin;
| a;b
| c;d;e
| \.
|
| col e'\';\';' SELECT;
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY ake FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY ake FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| SELECT ");(" SELECT;
| SELECT ----;a;b
|  SELECT b'01;10' 1 --; /*;a;b
|  +;
| AS x d_$y$p , e'\\\n;' 1;
< error
| 42601 syntax error at or near "SELECT"
| position: 39
> exec
| BEGIN
< ok 0 -
> exec
| SELECT ");(" SELECT;
< error
| 42703 column ");(" does not exist
| position: 8
> exec
| ROLLBACK
< ok 0 -
> exec
|
| SELECT ----;a;b
|  SELECT b'01;10' 1 --; /*;a;b
|  +;
< error
| 42601 syntax error at or near "SELECT"
| position: 19
> exec
| AS x /*/*/*a;b$;$select 1;*/*/*/ FROM t WHERE a =;
| , E'\';\\' AS x --a;bselect 1;$;$
|  AS x U&';--;$;$' SELECT;, 'a;bend;' FROM t WHERE a =; AS x /*/*/*$x$;a;b$x$;*/*/*/ , CREATE FUNCTION x_sr() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;'; SELECT 1; END , /*/*/*select 1;*/*/*/ AS x;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| AS x /*/*/*a;b$;$select 1;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| AS x /*/*/*a;b$;$select 1;*/*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "AS"
| position: 1
> exec
| 1 "select 1;" AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| 1 "select 1;" AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| 1 "select 1;" AS x;
< error
| 42601 syntax error at or near "1"
| position: 1
> exec
| + pz$k$n AS x /*/ *;END ;select 1;*/ , (SELECT 1; SELECT 2) AS x;
| + ';select 1;END ;' SELECT;
| AS x '--;' , ne$_$b , $umuc$;$umuc$ SELECT;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + pz$k$n AS x /*/ *;END ;select 1;*/ , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + pz$k$n AS x /*/ *;END ;select 1;*/ , (SELECT 1; SELECT 2) AS x;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| COPY lwl FROM stdin;
| 1	\N
| c;d;e
| \.
| 1 CREATE FUNCTION wuvv() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'end;;;'; END FROM t WHERE a = $oih$/*;/*;);($oih$ AS x;
| col --END ;
|  ,;
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lwl FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lwl FROM stdin;
| 1	\N
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "1"
| position: 22
> exec
| + E'\\\\\n;' , B'01;10' col /*/*end;;);(*/*/ FROM t WHERE a =; FROM t WHERE a = ";;/*;" SELECT U&'/*;' 1; AS x $$/*;$;$a;b$$ ,;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + E'\\\\\n;' , B'01;10' col /*/*end;;);(*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + E'\\\\\n;' , B'01;10' col /*/*end;;);(*/*/ FROM t WHERE a =;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| + U&'a;b;;end;' col --$;$end;);(
|  1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| + U&'a;b;;end;' col --$;$end;);(
|  1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| + U&'a;b;;end;' col --$;$end;);(
|  1;
< error
| 42601 syntax error at or near "+"
| position: 1
> exec
| , X'01;10' SELECT --select 1;
|  ,; SELECT U&'select 1;' 1;
| FROM t WHERE a = 'a;b' SELECT $sj$; $sx$$sj$ col;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| , X'01;10' SELECT --select 1;
|  ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| , X'01;10' SELECT --select 1;
|  ,;
< error
| 42601 syntax error at or near ","
| position: 1
> exec
| FROM t WHERE a = sq$s$q 1 U&'$;$' , e'\n;'';\\' SELECT;
| + bs$a$n AS x /*/*/*select 1;; end;*/*/*/ FROM t WHERE a =;+ CREATE FUNCTION etkg() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 'select 1;'; END SELECT $_be$/*;;;$_be$ + U&'/*;$;$' AS x;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| BEGIN
< ok 0 -
> exec
| FROM t WHERE a = sq$s$q 1 U&'$;$' , e'\n;'';\\' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| ROLLBACK
< ok 0 -
> exec
| FROM t WHERE a = sq$s$q 1 U&'$;$' , e'\n;'';\\' SELECT;
< error
| 42601 syntax error at or near "FROM"
| position: 1
> exec
| COPY lnw FROM stdin;
| a;b
| a;b
| \.
|  -- trailing note;
| COPY zyy FROM stdin;
| 1	\N
| 1	\N
| \.
|
| COPY haj FROM stdin;
| a;b
| c;d;e
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY lnw FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY lnw FROM stdin;
| a;b
| a;b
| \.
|
< error
| 42601 syntax error at or near "a"
| position: 22
> exec
| COPY cut FROM stdin;
| c;d;e
| a;b
| \.
| , U&'END ;/*;; ' + e'\\\';' col '$x$;\' SELECT;
| FROM t WHERE a = --$x$;;;
|  col;, (SELECT 1; SELECT 2) + E'\\\n;' SELECT $quxq$END ;;;$quxq$ AS x;
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| BEGIN
< ok 0 -
> exec
| COPY cut FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
> exec
| ROLLBACK
< ok 0 -
> exec
| COPY cut FROM stdin;
| c;d;e
| a;b
| \.
|
< error
| 42601 syntax error at or near "c"
| position: 22
