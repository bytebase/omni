package parser

import "testing"

// The ALTER TABLE clause set, measured against OceanBase CE 4.3.5.6 on
// 2026-11-04 and recorded here as that server answered. The forms came in with the family
// closed in tidb/parser, where a recorded corpus measures the gap; this fork
// has no such corpus, so the table is the measurement.
//
// The rejections are the point of the file: a clause read twice, a clause
// after the statement should have ended, a trailing comma, an empty list, and
// the clauses OceanBase does not have.
//
// Ref: alter_table_stmt, opt_part_option, signed_literal and
// opt_match_clause in MySQL 8.0.46 sql/sql_yacc.yy.

// alterTableClauseCase is a statement and whether the engine's parser accepts
// it.
type alterTableClauseCase struct {
	sql    string
	accept bool
}

var alterTableClauseCases = []alterTableClauseCase{
	{"ALTER TABLE db.t RENAME to db1.t1", true},
	{"ALTER TABLE db.t RENAME db1.t1", true},
	{"ALTER TABLE db.t RENAME as db1.t1", true},
	{"ALTER TABLE db.t RENAME = db1.t1", false},
	{"ALTER TABLE t RENAME = t1", false},
	{"ALTER TABLE t RENAME COLUMN a TO b", true},
	{"ALTER TABLE t RENAME COLUMN a b", false},
	{"ALTER TABLE t RENAME INDEX a b", false},
	{"ALTER TABLE t CONVERT TO CHARACTER SET utf8", true},
	{"ALTER TABLE t CONVERT TO CHAR SET UTF8", true},
	{"ALTER TABLE t CONVERT TO CHARSET UTF8", false},
	{"ALTER TABLE t CONVERT TO CHARACTER SET DEFAULT", false},
	{"ALTER TABLE t CONVERT TO CHARACTER SET DEFAULT COLLATE utf8_bin", false},
	{"ALTER TABLE t CONVERT TO utf8", false},
	{"ALTER TABLE t CONVERT TO CHARSET", false},
	{"ALTER TABLE t CONVERT CHARSET utf8", false},
	{"ALTER TABLE t CONVERT TO CHARACTER SET utf8 COLLATE = utf8_bin", false},
	{"ALTER TABLE t UNION = ()", true},
	{"ALTER TABLE t UNION ()", true},
	{"ALTER TABLE t UNION = (d.t1, d.t2)", true},
	{"ALTER TABLE t_n UNION ( ), KEY_BLOCK_SIZE = 1", true},
	{"ALTER TABLE t UNION", false},
	{"ALTER TABLE t UNION = (t1,)", false},
	{"ALTER TABLE t UNION = (,)", false},
	{"ALTER TABLE t ENGINE = InnoDB,", false},
	{"ALTER TABLE t ADD COLUMN a INT,", false},
	{"ALTER TABLE t ALGORITHM = DEFAULT,", false},
	{"ALTER TABLE t ALGORITHM = DEFAULT, LOCK = NONE", true},
	{"alter table d_n.t_n add constraint foreign key ident (ident(1)) references d_n.t_n match full on delete set null", false},
	{"alter table t_n add constraint ident foreign key (ident,ident(1)) references t_n match full on update set null on delete restrict", false},
	{"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH FULL", true},
	{"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH SIMPLE", true},
	{"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH foo", false},
	// A foreign key's list is key_list, not the key_list_with_expression an
	// index takes, so a functional key part is a syntax error even where
	// `ADD INDEX ((a + 1))` runs.
	{"ALTER TABLE t ADD FOREIGN KEY ((a + 1)) REFERENCES u(b)", false},
	{"ALTER TABLE t ADD CONSTRAINT c FOREIGN KEY ((a+1)) REFERENCES u(b)", false},
	{"CREATE TABLE t (a INT, FOREIGN KEY ((a + 1)) REFERENCES u(b))", false},
	{"ALTER TABLE t ADD INDEX ((a + 1))", true},
	{"ALTER TABLE t ADD ( CHECK ( true ) )", false},
	{"alter table t add column (b int, constraint c unique key (b))", false},
	{"ALTER TABLE `t` ADD COLUMN (`a1` INT, PRIMARY KEY (`a1`))", false},
	{"ALTER TABLE t ADD (a1 int, a2 int, PRIMARY KEY (a1), UNIQUE (a2))", false},
	{"ALTER TABLE t ADD COLUMN (INDEX (a))", false},
	{"ALTER TABLE t ADD COLUMN (FOREIGN KEY (a) REFERENCES u(b))", false},
	{"ALTER TABLE t ADD COLUMN (a INT, b INT)", true},
	{"ALTER TABLE t ADD COLUMN ()", false},
	{"ALTER TABLE t ADD COLUMN (a INT,)", false},
	{"alter table t partition by range FIELDS(a) (partition x values less than maxvalue)", false},
	{"alter table t partition by list FIELDS(a) (PARTITION p0 VALUES IN (5, 10, 15))", false},
	{"CREATE TABLE t (a int) PARTITION BY RANGE FIELDS(a) (PARTITION x VALUES LESS THAN MAXVALUE)", false},
	{"ALTER TABLE t PARTITION BY KEY FIELDS(a)", false},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE ENGINE text)", false},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE ENGINE = text)", false},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) MAX_ROWS 12 NODEGROUP 2)", false},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s MIN_ROWS 11))", false},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE DISK)", false},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) CHARSET utf8)", false},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) AUTO_INCREMENT = 1)", false},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) ROW_FORMAT = DYNAMIC)", false},
	{"ALTER TABLE t ALTER CONSTRAINT ident", false},
	{"ALTER TABLE t ALTER CHECK ident", false},
	{"ALTER TABLE t ADD a INT CONSTRAINT c CHECK (a > 1)", true},
	{"ALTER TABLE t ADD a INT CONSTRAINT CHECK (a > 1)", true},
	{"ALTER TABLE t ADD a INT CONSTRAINT c CHECK (a > 1) NOT ENFORCED", true},
	{"CREATE TABLE t (a INT CONSTRAINT c CHECK (a > 1))", true},
	{"ALTER TABLE t ADD a INT CONSTRAINT c NOT NULL", false},
	{"ALTER TABLE t ADD a INT CONSTRAINT c", false},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT -1", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT x'41'", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT _utf8'x'", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT DATE '2020-01-01'", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT (1+1)", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT CURRENT_TIMESTAMP", false},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT NOW()", false},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT 1+1", false},
}

// TestAlterTableClauseForms holds this parser to alterTableClauseCases without
// a server.
func TestAlterTableClauseForms(t *testing.T) {
	for _, c := range alterTableClauseCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v (err=%v)", c.sql, accepts, c.accept, err)
			}
		})
	}
}
