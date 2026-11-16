-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestQuery default#1
> exec
| CREATE EXTERNAL TABLE ext_users (id INT COMMENT 'the id', name STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS TEXTFILE LOCATION '/tmp/omni/ext_users'
< ok - -
> exec
| CREATE TABLE orders (id BIGINT, amount DECIMAL(10,2), note STRING COMMENT 'a note') STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE orders SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| INSERT INTO orders VALUES (1, 1.5, 'n1')
< ok - -
> exec
| CREATE TABLE events (id INT, body STRING) PARTITIONED BY (dt STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE events SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| ALTER TABLE events ADD PARTITION (dt='2024-01-01')
< ok - -
> exec
| INSERT INTO events PARTITION (dt='2024-01-01') VALUES (1, 'x'), (2, 'y')
< ok - -
> exec
| CREATE VIEW v_users AS SELECT id, name FROM ext_users
< ok - -
> exec
| USE `omni_query`
< ok - -
> query
| SELECT CAST(1 AS TINYINT) AS tiny, CAST(2 AS SMALLINT) AS small, CAST(3 AS INT) AS i
< result
< column "tiny" "TINYINT"
< column "small" "SMALLINT"
< column "i" "INT"
< row
: int64 1
: int64 2
: int64 3
> query
| SELECT CAST(4 AS BIGINT) AS big
< result
< column "big" "BIGINT"
< row
: int64 4
> query
| SELECT CAST(1.5 AS DOUBLE) AS dbl, CAST(2.5 AS FLOAT) AS flt
< result
< column "dbl" "DOUBLE"
< column "flt" "FLOAT"
< row
: float64 1.5
: float64 2.5
> query
| SELECT CAST('x' AS STRING) AS s, true AS b, NULL AS n
< result
< column "s" "STRING"
< column "b" "BOOLEAN"
< column "n" "NULL"
< row
: string
| x
: bool true
: null
> query
| SELECT CAST('abc' AS BINARY) AS bin
< result
< column "bin" "BINARY"
< row
: bytes
| abc
> query
| SELECT id, amount, note FROM orders
< result
< column "id" "BIGINT"
< column "amount" "DECIMAL"
< column "note" "STRING"
< row
: int64 1
: string
| 1.50
: string
| n1
