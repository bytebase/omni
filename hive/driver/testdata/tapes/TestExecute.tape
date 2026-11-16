-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestExecute default#1
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
| USE `omni_execute`
< ok - -
> exec
|
| 		CREATE TABLE exec_t (id INT, note STRING) STORED AS PARQUET
< ok - -
> exec
|
| 		INSERT INTO exec_t VALUES (1, 'a'), (2, 'b')
< ok - -
> exec
| USE `omni_execute`
< ok - -
> query
| SELECT id, note FROM exec_t
< result
< column "id" "INT"
< column "note" "STRING"
< row
: int64 1
: string
| a
< row
: int64 2
: string
| b
