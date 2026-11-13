-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestAlterTableClauseFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| ALTER TABLE db.t RENAME to db1.t1
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| ALTER TABLE db.t RENAME db1.t1
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| ALTER TABLE db.t RENAME as db1.t1
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| ALTER TABLE db.t RENAME = db1.t1
< error
| 1064 42000 …near '= db1.t1' at line 1
> exec ^
: string
| ALTER TABLE t RENAME = t1
< error
| 1064 42000 …near '= t1' at line 1
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO b
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a b
< error
| 1064 42000 …near 'b' at line 1
> exec ^
: string
| ALTER TABLE t RENAME INDEX a b
< error
| 1064 42000 …near 'b' at line 1
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET utf8
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t CONVERT TO CHAR SET UTF8
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET UTF8
< error
| 1064 42000 …near 'CHARSET UTF8' at line 1
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET DEFAULT COLLATE utf8_bin
< error
| 1064 42000 …near 'DEFAULT COLLATE utf8_bin' at line 1
> exec ^
: string
| ALTER TABLE t CONVERT TO utf8
< error
| 1064 42000 …near 'utf8' at line 1
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET
< error
| 1064 42000 …near 'CHARSET' at line 1
> exec ^
: string
| ALTER TABLE t CONVERT CHARSET utf8
< error
| 1064 42000 …near 'CHARSET utf8' at line 1
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET utf8 COLLATE = utf8_bin
< error
| 1064 42000 …near '= utf8_bin' at line 1
> exec ^
: string
| ALTER TABLE t UNION = ()
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t UNION ()
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t UNION = (d.t1, d.t2)
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t_n UNION ( ), KEY_BLOCK_SIZE = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t UNION
< error
| 1064 42000 …near 'UNION' at line 1
> exec ^
: string
| ALTER TABLE t UNION = (t1,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| ALTER TABLE t UNION = (,)
< error
| 1064 42000 …near ',)' at line 1
> exec ^
: string
| ALTER TABLE t ENGINE = InnoDB,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| ALTER TABLE t ALGORITHM = DEFAULT,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| ALTER TABLE t ALGORITHM = DEFAULT, LOCK = NONE
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table d_n.t_n add constraint foreign key ident (ident(1)) references d_n.t_n match full on delete set null
< error
| 1064 42000 …near '(1)) references d_n.t_n match full on delete set null' at line 1
> exec ^
: string
| alter table t_n add constraint ident foreign key (ident,ident(1)) references t_n match full on update set null on delete restrict
< error
| 1064 42000 …near '(1)) references t_n match full on update set null on delete restrict' at line 1
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH FULL
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH SIMPLE
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY ((a + 1)) REFERENCES u(b)
< error
| 1064 42000 …near '(a + 1)) REFERENCES u(b)' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c FOREIGN KEY ((a+1)) REFERENCES u(b)
< error
| 1064 42000 …near '(a+1)) REFERENCES u(b)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, FOREIGN KEY ((a + 1)) REFERENCES u(b))
< error
| 1064 42000 …near '(a + 1)) REFERENCES u(b))' at line 1
> exec ^
: string
| ALTER TABLE t ADD INDEX ((a + 1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ADD ( CHECK ( true ) )
< error
| 1064 42000 …near 'CHECK ( true ) )' at line 1
> exec ^
: string
| alter table t add column (b int, constraint c unique key (b))
< error
| 1064 42000 …near 'constraint c unique key (b))' at line 1
> exec ^
: string
| ALTER TABLE `t` ADD COLUMN (`a1` INT, PRIMARY KEY (`a1`))
< error
| 1064 42000 …near 'PRIMARY KEY (`a1`))' at line 1
> exec ^
: string
| ALTER TABLE t ADD (a1 int, a2 int, PRIMARY KEY (a1), UNIQUE (a2))
< error
| 1064 42000 …near 'PRIMARY KEY (a1), UNIQUE (a2))' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN (INDEX (a))
< error
| 1064 42000 …near 'INDEX (a))' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN (FOREIGN KEY (a) REFERENCES u(b))
< error
| 1064 42000 …near 'FOREIGN KEY (a) REFERENCES u(b))' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a INT, b INT)
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ADD COLUMN ()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a INT,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| alter table t partition by range FIELDS(a) (partition x values less than maxvalue)
< error
| 1064 42000 …near 'FIELDS(a) (partition x values less than maxvalue)' at line 1
> exec ^
: string
| alter table t partition by list FIELDS(a) (PARTITION p0 VALUES IN (5, 10, 15))
< error
| 1064 42000 …near 'FIELDS(a) (PARTITION p0 VALUES IN (5, 10, 15))' at line 1
> exec ^
: string
| CREATE TABLE t (a int) PARTITION BY RANGE FIELDS(a) (PARTITION x VALUES LESS THAN MAXVALUE)
< error
| 1064 42000 …near 'FIELDS(a) (PARTITION x VALUES LESS THAN MAXVALUE)' at line 1
> exec ^
: string
| ALTER TABLE t PARTITION BY KEY FIELDS(a)
< error
| 1064 42000 …near 'FIELDS(a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE ENGINE text)
< error
| 1064 42000 …near 'STORAGE ENGINE text)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE ENGINE = text)
< error
| 1064 42000 …near 'STORAGE ENGINE = text)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) MAX_ROWS 12 NODEGROUP 2)
< error
| 1064 42000 …near 'MAX_ROWS 12 NODEGROUP 2)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s MIN_ROWS 11))
< error
| 1064 42000 …near 'MIN_ROWS 11))' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE DISK)
< error
| 1064 42000 …near 'STORAGE DISK)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) CHARSET utf8)
< error
| 1064 42000 …near 'CHARSET utf8)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) AUTO_INCREMENT = 1)
< error
| 1064 42000 …near 'AUTO_INCREMENT = 1)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) ROW_FORMAT = DYNAMIC)
< error
| 1064 42000 …near 'ROW_FORMAT = DYNAMIC)' at line 1
> exec ^
: string
| ALTER TABLE t ALTER CONSTRAINT ident
< error
| 1064 42000 …near 'ident' at line 1
> exec ^
: string
| ALTER TABLE t ALTER CHECK ident
< error
| 1064 42000 …near 'ident' at line 1
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c CHECK (a > 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT CHECK (a > 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c CHECK (a > 1) NOT ENFORCED
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT CONSTRAINT c CHECK (a > 1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c
< error
| 1064 42000 …near 'c' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT x'41'
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT _utf8'x'
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT DATE '2020-01-01'
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT (1+1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT CURRENT_TIMESTAMP
< error
| 1064 42000 …near 'CURRENT_TIMESTAMP' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NOW()
< error
| 1064 42000 …near 'NOW()' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1+1
< error
| 1064 42000 …near '+1' at line 1
