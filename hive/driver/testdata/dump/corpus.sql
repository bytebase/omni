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

-- A partition that overrides its table's file format and SerDe properties, which
-- ADD PARTITION does not carry and only ALTER TABLE ... PARTITION (...) restores.
ALTER TABLE omni_dump.events ADD PARTITION (dt='2024-01-03')

ALTER TABLE omni_dump.events PARTITION (dt='2024-01-03') SET FILEFORMAT ORC

ALTER TABLE omni_dump.events PARTITION (dt='2024-01-03') SET SERDEPROPERTIES ('orc.compress'='ZLIB')

INSERT INTO omni_dump.events PARTITION (dt='2024-01-01') VALUES (1, 'x'), (2, 'y')

-- An external partitioned table whose partitions live at both default and custom
-- locations, with one partition unsetting a table SerDe property.
CREATE EXTERNAL TABLE omni_dump.ext_events (id INT) PARTITIONED BY (dt STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY '|' STORED AS TEXTFILE LOCATION '/tmp/omni/dump_ext_events'

ALTER TABLE omni_dump.ext_events ADD PARTITION (dt='2024-01-01')

ALTER TABLE omni_dump.ext_events ADD PARTITION (dt='2024-01-02') LOCATION 'file:/tmp/omni/dump_ext_events_dt_2024-01-02'

ALTER TABLE omni_dump.ext_events PARTITION (dt='2024-01-02') UNSET SERDEPROPERTIES ('field.delim')

-- Tables with named PRIMARY KEY, UNIQUE, and FOREIGN KEY constraints where the
-- child sorts before the parent: SHOW CREATE TABLE appends ALTER TABLE lines
-- without a terminating semicolon on CREATE TABLE, and the foreign key must run
-- after the parent table and its primary key exist.
CREATE TABLE omni_dump.z_parent (id INT, code STRING, CONSTRAINT pk_z_parent PRIMARY KEY (id) DISABLE NOVALIDATE RELY, CONSTRAINT uq_z_parent UNIQUE (code) DISABLE NOVALIDATE RELY) STORED AS PARQUET

CREATE TABLE omni_dump.a_child (id INT, parent_id INT, CONSTRAINT fk_a_child_parent FOREIGN KEY (parent_id) REFERENCES omni_dump.z_parent(id) DISABLE NOVALIDATE RELY) STORED AS PARQUET

-- A view with column comments, a view comment, and custom TBLPROPERTIES, all of
-- which SHOW CREATE TABLE omits.
CREATE VIEW omni_dump.v_docs (id COMMENT 'the id', label COMMENT 'it\'s a \\label') COMMENT 'documented view' TBLPROPERTIES ('team'='core') AS SELECT id, name FROM omni_dump.ext_users

-- A view whose comment SHOW CREATE TABLE leaves out, and the quote in it, which
-- a replay escapes rather than doubling. Its query names a column with an alias
-- SHOW CREATE TABLE has to quote: the statement the server states is the one a
-- replay takes, whatever the query could not be restated from.
CREATE VIEW omni_dump.z_users COMMENT 'the users\' rows' AS SELECT id AS 1bad, name FROM omni_dump.ext_users

-- A materialized view over that view, named so that name order alone would place
-- it before the view it resolves.
CREATE MATERIALIZED VIEW omni_dump.a_users_mv DISABLE REWRITE COMMENT 'the users materialized' AS SELECT `1bad`
FROM omni_dump.z_users

-- A partitioned materialized view: the statement states a partition as a table
-- does, with the column's type, and a materialized view takes it without.
CREATE MATERIALIZED VIEW omni_dump.p_events_mv DISABLE REWRITE COMMENT 'says PARTITIONED BY (old)' PARTITIONED ON (dt) AS SELECT id, dt FROM omni_dump.events
