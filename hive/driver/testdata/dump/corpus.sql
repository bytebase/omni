-- The schema the dump tests read: an external table, a name ending in "__", two
-- tables Hive translated to external and then made managed, one of them
-- partitioned, a view with a comment and a materialized view over it, beside
-- two databases that read each other's views.
CREATE EXTERNAL TABLE omni_dump.ext_users (id INT COMMENT 'the id', name STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS TEXTFILE LOCATION '/tmp/omni/dump_ext_users'

-- A name ending in "__", which Hive 4 has no index to give one to: every name
-- SHOW TABLES answers is a table's.
CREATE TABLE omni_dump.events__ (id INT) STORED AS PARQUET

-- Another database, with a view named as this one's a_top is and one reading
-- this database's z_base.
DROP DATABASE IF EXISTS omni_other CASCADE

CREATE DATABASE omni_other

CREATE TABLE omni_other.t (id INT) STORED AS PARQUET

CREATE VIEW omni_other.a_top AS SELECT id FROM omni_other.t

-- A view of this database reading the other database's a_top — a name a_top has
-- here too, which is not a read of this database's — and holding its reader's
-- qualified name in a literal, which is not a read either: counting it would
-- close a cycle with a_top's and leave the order to names.
CREATE VIEW omni_dump.z_base AS SELECT id, '`omni_dump`.`a_top`' AS tag FROM omni_other.a_top

-- A view of the other database reading z_base, and this database's a_top reading
-- that view: creating a_top resolves the other database's view, so z_base has to
-- be there first even though it is not a_top's stated parent. The backquoted
-- alias holds a quote.
CREATE VIEW omni_other.middle AS SELECT id FROM omni_dump.z_base

CREATE VIEW omni_dump.a_top AS SELECT id AS `it's` FROM omni_other.middle

-- A table Hive translated to external and then made managed: the only managed
-- table the pinned server can make, and one a dump restores as managed by
-- stating the ALTER beside the CREATE.
CREATE TABLE omni_dump.orders (id BIGINT, amount DECIMAL(10,2), note STRING COMMENT 'a note') STORED AS PARQUET

ALTER TABLE omni_dump.orders SET TBLPROPERTIES('EXTERNAL'='FALSE')

INSERT INTO omni_dump.orders VALUES (1, 1.5, 'n1')

-- A partitioned managed table with one partition and two rows, whose partition
-- is a name and a location rather than a clause SHOW CREATE TABLE states.
CREATE TABLE omni_dump.events (id INT, body STRING) PARTITIONED BY (dt STRING) STORED AS PARQUET

ALTER TABLE omni_dump.events SET TBLPROPERTIES('EXTERNAL'='FALSE')

ALTER TABLE omni_dump.events ADD PARTITION (dt='2024-01-01')

-- A partition that lives outside the table's own directory, which its name does
-- not state and only the describe of the partition does.
ALTER TABLE omni_dump.events ADD PARTITION (dt='2024-01-02') LOCATION 'file:/tmp/omni/dump_events_dt_2024-01-02'

INSERT INTO omni_dump.events PARTITION (dt='2024-01-01') VALUES (1, 'x'), (2, 'y')

-- A view whose comment SHOW CREATE TABLE leaves out, and the quote in it, which
-- a replay escapes rather than doubling.
CREATE VIEW omni_dump.z_users COMMENT 'the users\' rows' AS SELECT id, name FROM omni_dump.ext_users

-- A materialized view over that view, named so that name order alone would place
-- it before the view it resolves.
CREATE MATERIALIZED VIEW omni_dump.a_users_mv DISABLE REWRITE COMMENT 'the users materialized' AS SELECT id FROM omni_dump.z_users
