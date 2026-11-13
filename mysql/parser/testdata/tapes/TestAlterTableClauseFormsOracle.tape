-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestAlterTableClauseFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| ALTER TABLE db.t RENAME to db1.t1
< ok 0 0
> exec ^
: string
| ALTER TABLE db.t RENAME db1.t1
< ok 0 0
> exec ^
: string
| ALTER TABLE db.t RENAME as db1.t1
< ok 0 0
> exec ^
: string
| ALTER TABLE db.t RENAME = db1.t1
< ok 0 0
> exec ^
: string
| ALTER TABLE t RENAME = t1
< ok 0 0
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO b
< ok 0 0
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
< ok 0 0
> exec ^
: string
| ALTER TABLE t CONVERT TO CHAR SET UTF8
< ok 0 0
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET UTF8
< ok 0 0
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET DEFAULT
< ok 0 0
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET DEFAULT COLLATE utf8_bin
< error
| 1253 42000 COLLATION 'utf8mb3_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| ALTER TABLE t CONVERT TO utf8
< error
| 1064 42000 …near 'utf8' at line 1
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET
< error
| 1064 42000 …near '' at line 1
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
< ok 0 0
> exec ^
: string
| ALTER TABLE t UNION ()
< ok 0 0
> exec ^
: string
| ALTER TABLE t UNION = (d.t1, d.t2)
< ok 0 0
> exec ^
: string
| ALTER TABLE t_n UNION ( ), KEY_BLOCK_SIZE = 1
< ok 0 0
> exec ^
: string
| ALTER TABLE t UNION
< error
| 1064 42000 …near '' at line 1
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
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ALGORITHM = DEFAULT,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ALGORITHM = DEFAULT, LOCK = NONE
< ok 0 0
> exec ^
: string
| alter table d_n.t_n add constraint foreign key ident (ident(1)) references d_n.t_n match full on delete set null
< ok 0 0
> exec ^
: string
| alter table t_n add constraint ident foreign key (ident,ident(1)) references t_n match full on update set null on delete restrict
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH FULL
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH SIMPLE
< ok 0 0
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
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD ( CHECK ( true ) )
< ok 0 0
> exec ^
: string
| alter table t add column (b int, constraint c unique key (b))
< ok 0 0
> exec ^
: string
| ALTER TABLE `t` ADD COLUMN (`a1` INT, PRIMARY KEY (`a1`))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD (a1 int, a2 int, PRIMARY KEY (a1), UNIQUE (a2))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN (INDEX (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN (FOREIGN KEY (a) REFERENCES u(b))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a INT, b INT)
< ok 0 0
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
< ok 0 0
> exec ^
: string
| alter table t partition by list FIELDS(a) (PARTITION p0 VALUES IN (5, 10, 15))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) PARTITION BY RANGE FIELDS(a) (PARTITION x VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec ^
: string
| ALTER TABLE t PARTITION BY KEY FIELDS(a)
< error
| 1064 42000 …near 'FIELDS(a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE ENGINE text)
< error
| 1286 42000 Unknown storage engine 'text'
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE ENGINE = text)
< error
| 1286 42000 Unknown storage engine 'text'
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) MAX_ROWS 12 NODEGROUP 2)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s MIN_ROWS 11))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE DISK)
< error
| 1064 42000 …near 'DISK)' at line 1
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
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ALTER CHECK ident
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c CHECK (a > 1)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT CHECK (a > 1)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c CHECK (a > 1) NOT ENFORCED
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CONSTRAINT c CHECK (a > 1))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT x'41'
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT _utf8'x'
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT DATE '2020-01-01'
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT DATE '2020-01-01'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT (1+1)
< ok 0 0
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
