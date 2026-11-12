-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBaseFulltextParserGrammarOracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS obft
< ok 1 0
> exec
| USE obft
< ok 0 0
> exec
| DROP TABLE IF EXISTS obft_t
< ok 0 0
> exec
| CREATE TABLE obft_t (id INT NOT NULL, a INT, c VARCHAR(100))
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384 COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ngram2 PARSER_PROPERTIES=(min_ngram_size=2,max_ngram_size=3) BLOCK_SIZE 16384)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode="smart") BLOCK_SIZE 16384)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=4,max_token_size=20) BLOCK_SIZE 16384 /*!80000 INVISIBLE */)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(max_token_size=50, min_token_size=2))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c, d) WITH PARSER space PARSER_PROPERTIES (min_token_size=2 max_token_size=80))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES = (min_token_size = 3))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(MIN_TOKEN_SIZE=3, Max_Token_Size=40))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3, min_token_size=6))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=5) PARSER_PROPERTIES=(max_token_size=40))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(ngram_token_size=3) WITH PARSER ngram)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode=''max'' ''_word''))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(ik_mode="max_word", stopword_table=''none''))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ik PARSER_PROPERTIES=(stopword_table=''none'', dict_table=''none'', quantifier_table=''none''))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER `ngram`)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER NGRAM)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ngram WITH PARSER space)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) COMMENT ''x'' WITH PARSER ngram LOCAL)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER space PARSER_PROPERTIES=(ngram_token_size=3))'
< error
| 1235 0A000 space config not supported
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER nosuchparser)'
< error
| 1128 HY000 Function nosuchparser is not defined
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk USING BTREE (c))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) USING BTREE)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY USING BTREE fk (c))'
< error
| 1064 42000 …near 'fk (c))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk USING HASH (c))'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk USING (c))'
< error
| 1064 42000 …near '(c))' at line 1
> exec
| PREPARE obft_probe FROM 'ALTER TABLE obft_t ADD FULLTEXT KEY fk USING (c)'
< error
| 1064 42000 …near '(c)' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=())'
< error
| 1064 42000 …near '))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(foo=1))'
< error
| 1064 42000 …near 'foo=1))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(`min_token_size`=3))'
< error
| 1064 42000 …near '`min_token_size`=3))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size 3))'
< error
| 1064 42000 …near '3))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=''3''))'
< error
| 1064 42000 …near ''3'))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=-3))'
< error
| 1064 42000 …near '-3))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=+3))'
< error
| 1064 42000 …near '+3))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=3.5))'
< error
| 1064 42000 …near '3.5))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=0x3))'
< error
| 1064 42000 …near '0x3))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(stopword_table=3))'
< error
| 1064 42000 …near '3))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(ik_mode=smart))'
< error
| 1064 42000 …near 'smart))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(ik_mode=N''smart''))'
< error
| 1064 42000 …near 'N'smart'))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=3,))'
< error
| 1064 42000 …near '))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=3 , , max_token_size=40))'
< error
| 1064 42000 …near ', max_token_size=40))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES min_token_size=3)'
< error
| 1064 42000 …near 'min_token_size=3)' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=3) PARSER_PROPERTIES=())'
< error
| 1064 42000 …near '))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), d VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ''ngram'')'
< error
| 1064 42000 …near ''ngram')' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE FULLTEXT INDEX fk ON obft_t (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=4)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE FULLTEXT INDEX fk ON obft_t (c) PARSER_PROPERTIES=(min_token_size=3) LOCAL'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE FULLTEXT INDEX fk ON obft_t (c) PARSER_PROPERTIES=(ngram_token_size=3) WITH PARSER ngram LOCAL'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'ALTER TABLE obft_t ADD FULLTEXT KEY fk (c) WITH PARSER beng PARSER_PROPERTIES=(min_token_size=2)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'ALTER TABLE obft_t ADD FULLTEXT INDEX fk (c) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=3) COMMENT ''x'''
< ok 0 0
> exec
| PREPARE obft_probe FROM 'ALTER TABLE obft_t ADD FULLTEXT KEY fk (c) PARSER_PROPERTIES=(min_token_size=3)'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'ALTER TABLE obft_t ADD FULLTEXT fk (c) WITH PARSER ngram, ADD COLUMN zz INT'
< ok 0 0
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) WITH PARSER ngram)'
< error
| 1064 42000 …near 'PARSER ngram)' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY k (a) WITH PARSER ngram)'
< error
| 1064 42000 …near 'PARSER ngram)' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL KEY k (g) WITH PARSER ngram)'
< error
| 1064 42000 …near 'PARSER ngram)' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE INDEX k ON obft_t (a) WITH PARSER ngram'
< error
| 1064 42000 …near 'PARSER ngram' at line 1
> exec
| PREPARE obft_probe FROM 'ALTER TABLE obft_t ADD KEY k (a) WITH PARSER ngram'
< error
| 1064 42000 …near 'PARSER ngram' at line 1
> exec
| PREPARE obft_probe FROM 'ALTER TABLE obft_t ADD PRIMARY KEY (id) WITH PARSER ngram'
< error
| 1064 42000 …near 'PARSER ngram' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) PARSER_PROPERTIES=(min_token_size=3))'
< error
| 1064 42000 …near 'PARSER_PROPERTIES=(min_token_size=3))' at line 1
> exec
| PREPARE obft_probe FROM 'CREATE INDEX k ON obft_t (a) PARSER_PROPERTIES=(min_token_size=3)'
< error
| 1064 42000 …near 'PARSER_PROPERTIES=(min_token_size=3)' at line 1
> exec
| PREPARE obft_probe FROM 'ALTER TABLE obft_t ADD UNIQUE KEY k (a) PARSER_PROPERTIES=(min_token_size=3)'
< error
| 1064 42000 …near 'PARSER_PROPERTIES=(min_token_size=3)' at line 1
> exec
| DROP DATABASE IF EXISTS obft
< ok 1 0
