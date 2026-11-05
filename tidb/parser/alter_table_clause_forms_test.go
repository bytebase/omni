package parser

import "testing"

// The clause set ALTER TABLE reads, measured against pingcap/tidb:v8.5.5 on
// 2026-11-04 and recorded here as that server answered.
//
// The rejections are the point of this file. The corpus in
// testdata/pingcap_parser_corpus.tsv pins the forms pingcap's own parser tests
// write, and it is overwhelmingly accept rows: it holds almost none of the
// shapes a reader of a comma-separated clause list gets wrong -- a clause
// twice, a clause after the statement should have ended, a clause that belongs
// only to CREATE TABLE, a trailing comma, a list with nothing in it. Those are
// here.
//
// The last group is the words this parser reads without registering them as
// keywords. ATTRIBUTES, STATS_OPTIONS, TOPN, SAMPLES, SAMPLERATE, CMSKETCH,
// DEPTH and WIDTH are read as bare words precisely so they stay usable as
// names; the cases assert that, because a registration would have taken them
// away from every site that reads a name.
//
// Ref: AlterTableStmt, AlterTableSpec and AlterTableSpecSingleOpt in TiDB
// v8.5.5 pkg/parser/parser.y; alter_table_stmt in MySQL 8.0.46
// sql/sql_yacc.yy.

// alterTableClauseCase is a statement and whether the engine's parser accepts
// it.
type alterTableClauseCase struct {
	sql    string
	accept bool
}

var alterTableClauseCases = []alterTableClauseCase{
	// CACHE and NOCACHE: no value, and a comma between them.
	{"ALTER TABLE t CACHE", true},
	{"ALTER TABLE t NOCACHE", true},
	{"ALTER TABLE t CACHE, NOCACHE", true},
	{"ALTER TABLE t CACHE NOCACHE", false},
	{"ALTER TABLE t CACHE = 1", false},

	// READ ONLY and READ WRITE, which chain only with a comma.
	{"ALTER TABLE t READ ONLY", true},
	{"ALTER TABLE t READ WRITE", true},
	{"ALTER TABLE t READ ONLY, READ WRITE", true},
	{"ALTER TABLE t READ", false},
	{"ALTER TABLE t READ ONLY READ WRITE", false},
	{"ALTER TABLE t READ ONLY ENGINE = InnoDB", false},
	{"ALTER TABLE t ENGINE = InnoDB READ ONLY", false},

	// COMPACT is a whole arm of ALTER TABLE, so nothing joins it.
	{"alter table t compact", true},
	{"alter table t compact tiflash replica", true},
	{"alter table t compact partition p1,p2", true},
	{"alter table t compact partition p1,p2 tiflash replica", true},
	{"ALTER TABLE t COMPACT PARTITION", false},
	{"ALTER TABLE t COMPACT TIFLASH", false},
	{"ALTER TABLE t COMPACT REPLICA", false},
	{"ALTER TABLE t COMPACT, CACHE", false},
	{"ALTER TABLE t CACHE, COMPACT", false},

	// ATTRIBUTES and STATS_OPTIONS: a string or DEFAULT, never a number, and
	// an action of their own rather than a table option, so neither chains
	// with a space.
	{"ALTER TABLE t ATTRIBUTES='str'", true},
	{"ALTER TABLE t ATTRIBUTES 'str1,str2'", true},
	{"ALTER TABLE t ATTRIBUTES=DeFaUlT", true},
	{"ALTER TABLE t STATS_OPTIONS=\"str1,str2\"", true},
	{"ALTER TABLE t STATS_OPTIONS=DEFAULT", true},
	{"ALTER TABLE t STATS_OPTIONS = 'a', ATTRIBUTES = 'b'", true},
	{"ALTER TABLE t ATTRIBUTES", false},
	{"ALTER TABLE t ATTRIBUTES = 1", false},
	{"ALTER TABLE t ATTRIBUTES 1", false},
	{"ALTER TABLE t STATS_OPTIONS", false},
	{"ALTER TABLE t STATS_OPTIONS = 1", false},
	{"ALTER TABLE t ATTRIBUTES = 'a' ATTRIBUTES = 'b'", false},
	{"ALTER TABLE t ATTRIBUTES = 'a' ENGINE = x", false},
	{"ALTER TABLE t ENGINE = x ATTRIBUTES = 'a'", false},
	// STATS_OPTIONS is an ALTER TABLE action and not a CREATE TABLE option.
	{"CREATE TABLE t (a INT) STATS_OPTIONS = 'str'", false},

	// PARTITION name ... closes the statement: no comma before it, nothing
	// after it, and the options are the partition definition's own.
	{"ALTER TABLE t PARTITION p", true},
	{"ALTER TABLE t PARTITION p ATTRIBUTES='str'", true},
	{"ALTER TABLE t PARTITION p ATTRIBUTES=DEFAULT", true},
	{"ALTER TABLE t PARTITION p ENGINE = x MAX_ROWS 5", true},
	{"ALTER TABLE t PARTITION p PLACEMENT POLICY = x", true},
	{"ALTER TABLE t PARTITION p STORAGE ENGINE = x", true},
	{"ALTER TABLE t CACHE PARTITION p ATTRIBUTES = 'a'", true},
	{"ALTER TABLE t ADD COLUMN a INT PARTITION p ENGINE = x", true},
	{"ALTER TABLE t PARTITION", false},
	{"ALTER TABLE t PARTITION p ENGINE = x, MAX_ROWS 5", false},
	{"ALTER TABLE t PARTITION p STORAGE DISK", false},
	{"ALTER TABLE t PARTITION p ATTRIBUTES = 'a', CACHE", false},
	{"ALTER TABLE t CACHE, PARTITION p ATTRIBUTES = 'a'", false},
	{"ALTER TABLE t PARTITION p ATTRIBUTES = 'a', ATTRIBUTES = 'b'", false},
	{"ALTER TABLE t ADD COLUMN a INT, PARTITION p ENGINE = x", false},

	// The table options ALTER TABLE shares with CREATE TABLE, chained with
	// spaces and commas in one statement, and the empty UNION list.
	{"ALTER TABLE t_n UNION ( ), KEY_BLOCK_SIZE = 1", true},
	{"ALTER TABLE t UNION = (d.t1, d.t2)", true},
	{"ALTER TABLE d_n.t_n LOCK DEFAULT , UNION = ( t_n , d_n.t_n ) REMOVE PARTITIONING", true},
	{"ALTER TABLE d_n.t_n ALGORITHM = DEFAULT , MAX_ROWS 10, UNION ( d_n.t_n ) , ROW_FORMAT REDUNDANT, STATS_PERSISTENT = DEFAULT", true},
	{"ALTER TABLE t ROW_FORMAT REDUNDANT STATS_PERSISTENT = DEFAULT KEY_BLOCK_SIZE 1", true},
	{"ALTER TABLE t UNION", false},
	{"ALTER TABLE t UNION = (t1,)", false},
	{"ALTER TABLE t UNION = (,)", false},
	{"ALTER TABLE t LOCK DEFAULT, UNION = (t1), REMOVE PARTITIONING", false},
	{"ALTER TABLE t ENGINE = InnoDB,", false},
	{"ALTER TABLE t ADD COLUMN a INT,", false},
	{"ALTER TABLE t ALGORITHM = DEFAULT,", false},
	{"ALTER TABLE t ALGORITHM = DEFAULT ALGORITHM = DEFAULT", false},
	{"ALTER TABLE t ALGORITHM = DEFAULT LOCK = NONE", false},
	{"ALTER TABLE t ALGORITHM = DEFAULT, LOCK = NONE", true},
	// TTL_ENABLE is a table option whose value the grammar's own action
	// checks, so a word that is neither ON nor OFF is a parse error.
	{"ALTER TABLE t TTL_ENABLE = 'ON'", true},
	{"ALTER TABLE t TTL_ENABLE = 'off'", true},
	{"ALTER TABLE t TTL_ENABLE = 'test_case'", false},
	{"alter table t /*T![ttl] TTL_ENABLE = 'test_case' */", false},

	// FORCE on its own, and the two options it may precede.
	{"ALTER TABLE t FORCE", true},
	{"ALTER TABLE t FORCE AUTO_INCREMENT 3", true},
	{"ALTER TABLE t FORCE AUTO_INCREMENT = 3", true},
	{"alter table t force auto_random_base = 50", true},
	{"alter table t auto_increment 30, force auto_random_base 40", true},
	{"ALTER TABLE t FORCE, AUTO_INCREMENT 3", true},
	{"ALTER TABLE t FORCE ENGINE = InnoDB", false},
	{"ALTER TABLE t FORCE AUTO_ID_CACHE = 3", false},
	{"ALTER TABLE t FORCE AUTO_INCREMENT", false},

	// REMOVE PARTITIONING, and the partitioning clause's FIELDS spelling of
	// COLUMNS, which MySQL's own lexer maps to the same token.
	{"alter table t partition by range FIELDS(a) (partition x values less than maxvalue)", true},
	{"alter table t partition by list FIELDS(a) (PARTITION p0 VALUES IN (5, 10, 15))", true},
	{"alter table t partition by range FIELDS(a,b,c) (partition p1 values less than (1,1,1))", true},
	{"ALTER TABLE t PARTITION BY KEY FIELDS(a)", false},

	// UPDATE INDEXES, where GLOBAL or LOCAL is required on every entry.
	{"alter table t partition by range(a) update indexes (a local)", false},
	{"alter table t add column a int partition by hash(a) update indexes (idx_a global)", true},
	{"alter table t add column a int partition by hash(a) update indexes (idx_a global, idx_b local)", true},
	{"ALTER TABLE t PARTITION BY HASH(a) UPDATE INDEXES ()", false},
	{"ALTER TABLE t PARTITION BY HASH(a) UPDATE INDEXES (i)", false},
	{"ALTER TABLE t PARTITION BY HASH(a) UPDATE INDEXES (i GLOBAL,)", false},

	// The interval-partition actions. LAST and FIRST are entries of the
	// action list; SPLIT MAXVALUE and MERGE FIRST close the statement.
	{"ALTER TABLE t LAST PARTITION LESS THAN (1000)", true},
	{"ALTER TABLE t LAST PARTITION LESS THAN (1000) NO_WRITE_TO_BINLOG", true},
	{"ALTER TABLE t FIRST PARTITION LESS THAN (1000)", true},
	{"ALTER TABLE t FIRST PARTITION LESS THAN (1000) IF EXISTS", true},
	{"ALTER TABLE t split MAXVALUE PARTITION LESS THAN (1000)", true},
	{"ALTER TABLE t merge first PARTITION LESS THAN (1000)", true},
	{"ALTER TABLE t ADD COLUMN a INT, LAST PARTITION LESS THAN (1)", true},
	{"ALTER TABLE t LAST PARTITION LESS THAN (1), ADD COLUMN a INT", true},
	{"ALTER TABLE t LAST PARTITION LESS THAN 1000", false},
	{"ALTER TABLE t LAST PARTITION LESS THAN (MAXVALUE)", false},
	{"ALTER TABLE t SPLIT PARTITION LESS THAN (1000)", false},
	{"ALTER TABLE t MERGE PARTITION LESS THAN (1000)", false},
	{"ALTER TABLE t ADD COLUMN a INT, SPLIT MAXVALUE PARTITION LESS THAN (1)", false},
	{"ALTER TABLE t SPLIT MAXVALUE PARTITION LESS THAN (1), ADD COLUMN a INT", false},

	// ANALYZE PARTITION, which is a whole arm of ALTER TABLE in TiDB: a plain
	// name list with no ALL, and nothing beside it.
	{"ALTER TABLE t ANALYZE PARTITION p", true},
	{"ALTER TABLE t ANALYZE PARTITION p, q", true},
	{"alter table t analyze partition a with 4 buckets", true},
	{"alter table t analyze partition a index b", true},
	{"alter table t analyze partition a index b with 4 buckets", true},
	{"ALTER TABLE t ANALYZE PARTITION p INDEX", true},
	{"ALTER TABLE t ANALYZE PARTITION p INDEX i, j WITH 4 BUCKETS, 3 TOPN", true},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 2 CMSKETCH DEPTH", true},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 2 CMSKETCH WIDTH", true},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 2 SAMPLES", true},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 0.5 SAMPLERATE", true},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 2 SAMPLERATE", true},
	{"ALTER TABLE t ANALYZE PARTITION ALL", false},
	{"ALTER TABLE t ANALYZE PARTITION p, ADD COLUMN a INT", false},
	{"ALTER TABLE t ADD COLUMN a INT, ANALYZE PARTITION p", false},
	{"ALTER TABLE t ANALYZE PARTITION p WITH", false},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 4", false},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 4 BUCKETS,", false},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 2 CMSKETCH", false},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 0.5 BUCKETS", false},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 4 BUCKETS 3 TOPN", false},

	// IF NOT EXISTS and IF EXISTS, on the actions that have them and not on
	// the actions that do not.
	{"ALTER TABLE t ADD COLUMN IF NOT EXISTS (a SMALLINT UNSIGNED)", true},
	{"ALTER TABLE t ADD COLUMN IF NOT EXISTS a SMALLINT UNSIGNED AFTER b", true},
	{"ALTER TABLE t ADD IF NOT EXISTS a INT", true},
	{"ALTER TABLE employees ADD PARTITION IF NOT EXISTS ( PARTITION P2 VALUES LESS THAN MAXVALUE)", true},
	{"ALTER TABLE employees ADD PARTITION IF NOT EXISTS PARTITIONS 5", true},
	{"ALTER TABLE t ADD PARTITION IF NOT EXISTS LOCAL PARTITIONS 5", true},
	{"alter table t drop partition if exists p1, p2", true},
	{"ALTER TABLE t MODIFY COLUMN IF EXISTS a varchar(255)", true},
	{"ALTER TABLE t CHANGE COLUMN IF EXISTS a b varchar(255)", true},
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS (a) USING BTREE COMMENT 'a'", true},
	{"ALTER TABLE t ADD KEY IF NOT EXISTS (a) USING HASH COMMENT 'a'", true},
	{"ALTER TABLE t DROP INDEX IF EXISTS a", true},
	{"ALTER TABLE t DROP KEY IF EXISTS a", true},
	{"ALTER TABLE t ADD CONSTRAINT fk_t2_id FOREIGN KEY IF NOT EXISTS (t2_id) REFERENCES t(id)", true},
	{"ALTER TABLE t ADD COLUMN IF EXISTS a INT", false},
	{"ALTER TABLE t MODIFY COLUMN IF NOT EXISTS a INT", false},
	{"ALTER TABLE t ADD UNIQUE IF NOT EXISTS (a)", false},
	{"ALTER TABLE t ADD UNIQUE KEY IF NOT EXISTS (a)", false},
	{"ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a)", false},
	{"ALTER TABLE t ADD FULLTEXT IF NOT EXISTS (a)", false},
	{"ALTER TABLE t ADD CONSTRAINT c CHECK IF NOT EXISTS (a > 1)", false},
	{"ALTER TABLE t DROP CONSTRAINT IF EXISTS c", false},
	{"ALTER TABLE t DROP CHECK IF EXISTS c", false},
	{"ALTER TABLE t DROP FOREIGN KEY IF EXISTS c", false},
	{"ALTER TABLE t DROP PRIMARY KEY IF EXISTS", false},

	// GLOBAL, LOCAL, CLUSTERED and NONCLUSTERED are index options, so they
	// appear in any order, more than once, and on every key kind -- on
	// CREATE TABLE and CREATE INDEX too.
	{"alter table t add index (a) local", true},
	{"alter table t add index (a) global", true},
	{"alter table t add unique (a) local", true},
	{"alter table t add unique key (a) global", true},
	{"alter table t add primary key (a) global", true},
	{"ALTER TABLE t ADD UNIQUE INDEX (a) LOCAL", true},
	{"ALTER TABLE t ADD INDEX (a) GLOBAL LOCAL", true},
	{"ALTER TABLE t ADD INDEX (a) COMMENT 'c' GLOBAL", true},
	{"ALTER TABLE t ADD INDEX (a) GLOBAL COMMENT 'c'", true},
	{"ALTER TABLE t ADD FULLTEXT (a) GLOBAL", true},
	{"ALTER TABLE t ADD INDEX (a) NONCLUSTERED", true},
	{"ALTER TABLE t ADD PRIMARY KEY (a) CLUSTERED GLOBAL", true},
	{"ALTER TABLE t ADD PRIMARY KEY (a) GLOBAL CLUSTERED", true},
	{"CREATE TABLE t (a int, INDEX (a) GLOBAL)", true},
	{"CREATE INDEX i ON t (a) GLOBAL", true},
	// On a column option it is one of the two at most, and after CLUSTERED.
	{"ALTER TABLE t ADD a INT PRIMARY KEY GLOBAL", true},
	{"ALTER TABLE t ADD a INT PRIMARY KEY CLUSTERED GLOBAL", true},
	{"ALTER TABLE t ADD a INT UNIQUE GLOBAL", true},
	{"ALTER TABLE t ADD a INT UNIQUE KEY LOCAL", true},
	{"ALTER TABLE t ADD a INT KEY GLOBAL", true},
	{"ALTER TABLE t ADD a INT PRIMARY KEY GLOBAL CLUSTERED", false},
	{"ALTER TABLE t ADD a INT UNIQUE GLOBAL LOCAL", false},

	// VECTOR INDEX, which takes no CONSTRAINT symbol and is spelled INDEX.
	{"ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a, b))) USING HNSW COMMENT 'a'", true},
	{"ALTER TABLE t ADD VECTOR INDEX ((lower(a))) USING HNSW COMMENT 'a'", true},
	{"ALTER TABLE t ADD VECTOR INDEX (a, (VEC_COSINE_DISTANCE(a))) USING HNSW COMMENT 'a'", false},
	{"ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a))) USING HYPO COMMENT 'a'", false},
	{"ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a))) COMMENT 'a'", true},
	{"ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS ((VEC_COSINE_DISTANCE(a))) USING HNSW COMMENT 'a'", true},
	{"ALTER TABLE t ADD VECTOR INDEX idx USING HNSW ((lower(a)))", true},
	{"ALTER TABLE t ADD VECTOR INDEX ((lower(a))) GLOBAL", true},
	{"ALTER TABLE t ADD VECTOR KEY ((lower(a)))", false},
	{"ALTER TABLE t ADD VECTOR INDEX ()", false},
	{"ALTER TABLE t ADD CONSTRAINT c VECTOR INDEX ((lower(a)))", false},
	// HYPO and HNSW are index type names, and unquoted ones only.
	{"CREATE INDEX i ON t (a) USING HYPO", true},
	{"CREATE INDEX i ON t (a) USING `HYPO`", false},

	// RENAME, whose destination is a table name and whose separator may be an
	// equals sign; RENAME COLUMN and RENAME INDEX need their TO.
	{"ALTER TABLE db.t RENAME to db1.t1", true},
	{"ALTER TABLE db.t RENAME db1.t1", true},
	{"ALTER TABLE db.t RENAME = db1.t1", true},
	{"ALTER TABLE db.t RENAME as db1.t1", true},
	{"ALTER TABLE t RENAME = t1", true},
	{"ALTER TABLE t RENAME COLUMN a b", false},
	{"ALTER TABLE t RENAME INDEX a b", false},

	// CONVERT TO, in all three spellings of the keyword, with DEFAULT as a
	// charset name and with no equals sign before the collation.
	{"ALTER TABLE t CONVERT TO CHARSET UTF8", true},
	{"ALTER TABLE t CONVERT TO CHARSET UTF8 COLLATE UTF8_BIN", true},
	{"alter table d_n.t_n convert to charset default", true},
	{"alter table d_n.t_n convert to char set default", true},
	{"ALTER TABLE t CONVERT TO CHARACTER SET DEFAULT COLLATE utf8_bin", true},
	{"ALTER TABLE t CONVERT TO utf8", false},
	{"ALTER TABLE t CONVERT TO CHARSET", false},
	{"ALTER TABLE t CONVERT CHARSET utf8", false},
	{"ALTER TABLE t CONVERT TO CHARSET utf8 COLLATE = utf8_bin", false},

	// ADD [COLUMN] ( ... ) reads a table element list, so constraints come
	// with the columns; the list may not be empty and takes no trailing
	// comma.
	{"ALTER TABLE t ADD ( CHECK ( true ) )", true},
	{"ALTER TABLE t ADD ( CONSTRAINT CHECK ( true ) )", true},
	{"alter table t add column (b int, constraint c unique key (b))", true},
	{"ALTER TABLE t ADD COLUMN ( CONSTRAINT CHECK ( true ) ENFORCED , CHECK ( true ) )", true},
	{"ALTER TABLE t ADD COLUMN (a1 int, a2 int, CONSTRAINT b1 CHECK (a1>0), CONSTRAINT b2 CHECK (a2<10))", true},
	{"ALTER TABLE `t` ADD COLUMN (`a1` INT, PRIMARY KEY (`a1`))", true},
	{"ALTER TABLE t ADD (a1 int, a2 int, PRIMARY KEY (a1), UNIQUE (a2))", true},
	{"ALTER TABLE t ADD COLUMN (INDEX (a))", true},
	{"ALTER TABLE t ADD COLUMN (FOREIGN KEY (a) REFERENCES u(b))", true},
	{"ALTER TABLE t ADD COLUMN ()", false},
	{"ALTER TABLE t ADD COLUMN (a INT,)", false},

	// A column's CONSTRAINT [symbol] CHECK, where CHECK is the only word that
	// may follow the symbol.
	{"ALTER TABLE t ADD a1 int CONSTRAINT ident CHECK ( a1 > 1 ) REFERENCES b ON DELETE CASCADE ON UPDATE CASCADE", true},
	{"ALTER TABLE t ADD COLUMN a DATE CONSTRAINT CHECK ( a > 0 ) FIRST", true},
	{"ALTER TABLE t ADD a2 int CONSTRAINT ident CHECK (a2 > 1) ENFORCED", true},
	{"ALTER TABLE t ADD a2 int CONSTRAINT ident CHECK (a2 > 1) NOT ENFORCED", true},
	{"CREATE TABLE t (a INT CONSTRAINT c CHECK (a > 1))", true},
	{"ALTER TABLE t ADD a INT CONSTRAINT c NOT NULL", false},
	{"ALTER TABLE t ADD a INT CONSTRAINT c", false},
	// A column with no type is a syntax error however the options read.
	{"ALTER TABLE t ADD a CONSTRAINT c CHECK (a>1)", false},

	// A foreign key's column list is a key-part list, and its MATCH takes
	// one of three words.
	{"alter table d_n.t_n add constraint foreign key ident (ident(1)) references d_n.t_n match full on delete set null", true},
	{"alter table t_n add constraint ident foreign key (ident,ident(1)) references t_n match full on update set null on delete restrict", true},
	{"alter table d_n.t_n add foreign key ident (ident, ident(1) asc) references t_n match partial on delete cascade remove partitioning", true},
	{"alter table d_n.t_n add constraint foreign key (ident asc) references d_n.t_n match simple on update cascade on delete cascade", true},
	{"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH foo", false},
	// A functional key part in a foreign key's list parses on TiDB, whose
	// own grammar reads the same IndexPartSpecificationList an index takes;
	// v8.5.5 then crashes building the key (1105, a nil dereference), which
	// is not a parse error. MySQL 8.0.46, MariaDB 11.8.8 and OceanBase CE
	// 4.3.5.6 all answer 1064, and their forks reject it.
	{"ALTER TABLE t ADD FOREIGN KEY ((a + 1)) REFERENCES u(b)", true},
	{"ALTER TABLE t ADD CONSTRAINT c FOREIGN KEY ((a+1)) REFERENCES u(b)", true},
	{"CREATE TABLE t (a INT, FOREIGN KEY ((a + 1)) REFERENCES u(b))", true},
	// IF NOT EXISTS on a vector index, which the catalog honours too.
	{"ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS idx ((lower(a)))", true},

	// A partition definition's own options, which are not a table's: STORAGE
	// ENGINE is one of them and STORAGE {DISK | MEMORY} is not, and the
	// catch-all list partition is spelled two ways.
	{"ALTER TABLE d_n.t_n WITHOUT VALIDATION , ADD PARTITION ( PARTITION ident VALUES LESS THAN ( MAXVALUE ) STORAGE ENGINE text_string MAX_ROWS 12 )", true},
	{"ALTER TABLE d_n.t_n WITH VALIDATION , ADD PARTITION NO_WRITE_TO_BINLOG (PARTITION ident VALUES LESS THAN MAXVALUE STORAGE ENGINE = text_string, PARTITION ident VALUES LESS THAN ( MAXVALUE ) (SUBPARTITION text_string MIN_ROWS 11))", true},
	{"ALTER TABLE d_n.t_n WITHOUT VALIDATION , ADD PARTITION ( PARTITION ident VALUES IN ( DEFAULT ) STORAGE ENGINE text_string MAX_ROWS 12 )", true},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (DEFAULT, 1))", true},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (1, DEFAULT))", true},
	{"ALTER TABLE t ADD PARTITION (PARTITION p DEFAULT)", true},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) PLACEMENT POLICY x)", true},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE DISK)", false},
	{"ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) TTL_ENABLE = 'ON')", false},

	// Placement policy on one partition, with and without the version
	// comment TiDB writes it inside.
	{"alter table m partition t placement policy=\"ww\"", true},
	{"alter table m partition t /*T![placement] placement policy=\"ww\" */", true},

	// SET DEFAULT takes a signed literal, or any expression in parentheses.
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT 1", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT -1", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT 'x'", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT TRUE", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT (CURRENT_TIMESTAMP)", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT (1+1)", true},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT CURRENT_TIMESTAMP", false},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT NOW()", false},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT 1+1", false},
	{"ALTER TABLE t ALTER COLUMN a SET DEFAULT -'x'", false},

	// ALTER {CHECK | CONSTRAINT} needs its enforcement.
	{"ALTER TABLE t ALTER CONSTRAINT ident ENFORCED", true},
	{"ALTER TABLE t ALTER CHECK ident NOT ENFORCED", true},
	{"ALTER TABLE t_n ALTER CONSTRAINT ident", false},
	{"ALTER TABLE t ALTER CHECK ident", false},

	// The words this parser reads without registering them, which stay
	// usable as table, column, index and alias names.
	{"ALTER TABLE attributes ADD COLUMN a INT", true},
	{"ALTER TABLE t ADD COLUMN attributes INT", true},
	{"ALTER TABLE t ADD COLUMN stats_options INT", true},
	{"ALTER TABLE t ADD COLUMN vector INT", true},
	{"ALTER TABLE t ADD COLUMN compact INT", true},
	{"ALTER TABLE t DROP COLUMN cache", true},
	{"ALTER TABLE t DROP COLUMN nocache", true},
	{"ALTER TABLE t DROP COLUMN last", true},
	{"ALTER TABLE t DROP COLUMN global", true},
	{"ALTER TABLE t RENAME COLUMN a TO attributes", true},
	{"ALTER TABLE t RENAME TO attributes", true},
	{"ALTER TABLE t ADD INDEX i (attributes)", true},
	{"ALTER TABLE t ADD INDEX attributes (a)", true},
	{"ALTER TABLE t ADD INDEX cmsketch (a)", true},
	{"ALTER TABLE t ADD INDEX i (global, local)", true},
	{"ALTER TABLE t MODIFY COLUMN attributes INT", true},
	{"ALTER TABLE t CHANGE COLUMN a attributes INT", true},
	{"CREATE TABLE attributes (a INT)", true},
	{"CREATE TABLE t (attributes INT, stats_options INT, topn INT, samplerate INT, cmsketch INT, depth INT, width INT, samples INT)", true},
	{"CREATE INDEX attributes ON t (a)", true},
	{"SELECT * FROM t AS attributes", true},
	{"SELECT attributes.a FROM t attributes", true},
	{"DELETE FROM t attributes WHERE attributes.id = 1", true},
	{"UPDATE t attributes SET attributes.a = 1", true},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 4 buckets", true},
	{"ALTER TABLE t ANALYZE PARTITION p WITH 4 topn", true},
	// The words this file's clauses dispatch on are names too, in every
	// position that reads one: an action word only begins an action at the
	// start of one.
	{"ALTER TABLE t ADD INDEX i (first, last, merge, split)", true},
	{"ALTER TABLE t ADD COLUMN merge INT", true},
	{"ALTER TABLE t ADD COLUMN split INT", true},
	{"ALTER TABLE t ADD COLUMN global INT", true},
	{"ALTER TABLE t ADD COLUMN local INT", true},
	{"ALTER TABLE t ADD COLUMN clustered INT", true},
	{"ALTER TABLE t ADD COLUMN nonclustered INT", true},
	{"ALTER TABLE t ADD INDEX global (a)", true},
	{"ALTER TABLE t DROP INDEX local", true},
	{"ALTER TABLE t MODIFY COLUMN compact INT", true},
	{"ALTER TABLE t CHANGE COLUMN a vector INT", true},
	{"ALTER TABLE t RENAME COLUMN a TO last", true},
	{"ALTER TABLE merge ADD COLUMN a INT", true},
	{"ALTER TABLE t ADD PARTITION (PARTITION last VALUES LESS THAN (1))", true},
	{"ALTER TABLE t DROP PARTITION first, last", true},
	{"ALTER TABLE t TRUNCATE PARTITION cache", true},
	{"ALTER TABLE t ANALYZE PARTITION buckets", true},
	{"ALTER TABLE t ANALYZE PARTITION p INDEX topn", true},
	{"ALTER TABLE t COMPACT PARTITION tiflash", true},
	{"ALTER TABLE t PARTITION attributes ENGINE = x", true},
	{"ALTER TABLE t PARTITION p PLACEMENT POLICY = attributes", true},
	{"ALTER TABLE t ALTER COLUMN global SET DEFAULT 1", true},
	{"ALTER TABLE t ALTER INDEX global VISIBLE", true},
	{"ALTER TABLE t ADD CONSTRAINT global CHECK (a > 1)", true},
	{"ALTER TABLE t DROP CONSTRAINT vector", true},
	{"SELECT first, last, merge, split, global, local FROM t", true},
	// A backquoted spelling is a name, never the clause.
	{"ALTER TABLE t `ATTRIBUTES` = 'a'", false},
	{"ALTER TABLE t `STATS_OPTIONS` = 'a'", false},

	// A word left over after a clause ends the statement is a syntax error,
	// not a clause the reader stops before. A reader that returns on an
	// unknown token instead of failing leaves the word unconsumed, and a
	// clause whose list may be empty -- the partition options, the COMPACT
	// modifiers, the analyze options -- is where that happens.
	{"ALTER TABLE t PARTITION p GARBAGE", false},
	{"ALTER TABLE t COMPACT GARBAGE", false},
	{"ALTER TABLE t ANALYZE PARTITION p GARBAGE", false},
	{"ALTER TABLE t CACHE GARBAGE", false},
	{"ALTER TABLE t READ ONLY GARBAGE", false},
	{"ALTER TABLE t ATTRIBUTES = 'a' GARBAGE", false},
	{"ALTER TABLE t SPLIT MAXVALUE PARTITION LESS THAN (1) GARBAGE", false},
	{"ALTER TABLE t LAST PARTITION LESS THAN (1) GARBAGE", false},
	{"ALTER TABLE t ADD VECTOR INDEX ((lower(a))) GARBAGE", false},
	{"ALTER TABLE t PARTITION BY HASH(a) UPDATE INDEXES (i GLOBAL) GARBAGE", false},
	{"ALTER TABLE t ADD COLUMN (a INT, PRIMARY KEY (a)) GARBAGE", false},
	{"ALTER TABLE t RENAME = d.t1 GARBAGE", false},
	{"ALTER TABLE t CONVERT TO CHARSET utf8 GARBAGE", false},
	{"ALTER TABLE t UNION = () GARBAGE", false},
}

// TestAlterTableClauseForms holds this parser to alterTableClauseCases without
// a server. It needs no container, so it is the gate a change to the ALTER
// TABLE clause readers runs first.
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
