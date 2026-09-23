CREATE EXTERNAL TABLE t1 OPTIONS(a = b, c = "def");
CREATE PRIVATE EXTERNAL TABLE a.b.c.t1 OPTIONS(size = 10);
CREATE OR REPLACE TEMP EXTERNAL TABLE IF NOT EXISTS T OPTIONS();
CREATE EXTERNAL TABLE `function` OPTIONS();
CREATE EXTERNAL TABLE `function` OPTIONS();
CREATE PRIVATE EXTERNAL TABLE t
(
  x int64
) OPTIONS();
CREATE OR REPLACE TEMP EXTERNAL TABLE IF NOT EXISTS t
(
  x int64
) OPTIONS();
CREATE EXTERNAL TABLE t1
(
  element1 bool,
  element2 string
) OPTIONS(a = '{"jsonkey": "jsonvalue"}', c = "def");
CREATE EXTERNAL TABLE t1
(
  x int64 NOT NULL
) OPTIONS(a = b);
CREATE EXTERNAL TABLE projectid.datasetid.tablename
(
  x int64,
  y string,
  z string,
  PRIMARY KEY(x, y ASC, z DESC)
) OPTIONS(a = b);
CREATE EXTERNAL TABLE projectid.datasetid.tablename WITH PARTITION COLUMNS OPTIONS();
CREATE EXTERNAL TABLE projectid.datasetid.tablename
(
  x int64
) WITH PARTITION COLUMNS(
  y int64,
  z string
) OPTIONS();
CREATE EXTERNAL TABLE projectid.datasetid.tablename WITH PARTITION COLUMNS(
  y int64,
  PRIMARY KEY(y)
) OPTIONS();
CREATE EXTERNAL TABLE projectid.datasetid.tablename
(
  x int64,
  PRIMARY KEY(x)
) WITH PARTITION COLUMNS(
  CHECK(x > 0) ENFORCED
) OPTIONS();
CREATE EXTERNAL TABLE projectid.datasetid.tablename
(
  x int64
) WITH CONNECTION DEFAULT OPTIONS(x = foo);
CREATE EXTERNAL TABLE projectid.datasetid.tablename
(
  x int64
) WITH CONNECTION conn_1 OPTIONS(x = foo);
CREATE EXTERNAL TABLE projectid.datasetid.tablename
(
  x int64
) WITH CONNECTION proj.loc.conn_id OPTIONS(x = foo);
CREATE EXTERNAL TABLE projectid.datasetid.tablename
(
  x int64
) WITH PARTITION COLUMNS WITH CONNECTION `project.country-region-5.connection_1` OPTIONS(x = bar);
CREATE EXTERNAL TABLE t
(
  x int64
) WITH PARTITION COLUMNS(
  y int64,
  PRIMARY KEY(y DESC)
) WITH CONNECTION a.b.c OPTIONS(x = baz);
CREATE EXTERNAL TABLE t1
(
  a int32,
  PRIMARY KEY(a DESC)
) OPTIONS();
