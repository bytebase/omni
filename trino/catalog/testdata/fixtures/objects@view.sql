CREATE TABLE t (b integer NOT NULL, a integer COMMENT 'a', d varchar(10)) COMMENT 't';
CREATE VIEW a AS SELECT b, a, d FROM t;
