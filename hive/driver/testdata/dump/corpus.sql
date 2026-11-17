-- The schema the dump tests read. The external table and the two views are what
-- TestDumpRestores restores; the managed table and the partitioned one carry
-- what a dump leaves out, so each line of theirs ends with the gap marker.
CREATE EXTERNAL TABLE omni_dump.ext_users (id INT COMMENT 'the id', name STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS TEXTFILE LOCATION '/tmp/omni/dump_ext_users'

CREATE VIEW omni_dump.z_base AS SELECT id FROM omni_dump.ext_users

-- Only the order of the views, not their names, can replay this one.
CREATE VIEW omni_dump.a_top AS SELECT id FROM omni_dump.z_base

-- A table Hive translated to external and then made managed: the only managed
-- table the pinned server can make, and the one a dump cannot restore as
-- managed.
CREATE TABLE omni_dump.orders (id BIGINT, amount DECIMAL(10,2), note STRING COMMENT 'a note') STORED AS PARQUET -- not restorable: TestDumpGaps

ALTER TABLE omni_dump.orders SET TBLPROPERTIES('EXTERNAL'='FALSE') -- not restorable: TestDumpGaps

INSERT INTO omni_dump.orders VALUES (1, 1.5, 'n1') -- not restorable: TestDumpGaps

-- A partitioned managed table with one partition and two rows.
CREATE TABLE omni_dump.events (id INT, body STRING) PARTITIONED BY (dt STRING) STORED AS PARQUET -- not restorable: TestDumpGaps

ALTER TABLE omni_dump.events SET TBLPROPERTIES('EXTERNAL'='FALSE') -- not restorable: TestDumpGaps

ALTER TABLE omni_dump.events ADD PARTITION (dt='2024-01-01') -- not restorable: TestDumpGaps

INSERT INTO omni_dump.events PARTITION (dt='2024-01-01') VALUES (1, 'x'), (2, 'y') -- not restorable: TestDumpGaps
