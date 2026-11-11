package parser

import (
	"math"
	"strings"
	"testing"
)

// numberLiterals are the numbers written at each position below. The first nine are
// runs of digits: 1, the widest that fits 31, 63 and 64 bits and the one past each,
// twenty nines, and a 1 behind leading zeros, which are not part of its width. The
// last two are a decimal and a number with an exponent.
var numberLiterals = []string{
	"1",
	"2147483647",
	"2147483648",
	"9223372036854775807",
	"9223372036854775808",
	"18446744073709551615",
	"18446744073709551616",
	"99999999999999999999",
	"000000000000000000000001",
	"1.5",
	"1e1",
}

// numberTokenForm is a statement with %s where a number goes, and whether OceanBase 4.3.5
// parses it with each of numberLiterals there, in order: Y or N.
type numberTokenForm struct {
	sql     string
	accepts string
}

// numberTokenForms are the positions held to OceanBase 4.3.5 for every literal. The lexer
// gives a run of digits past 64 bits the token of a decimal, so a position answers
// for 18446744073709551616 as it does for 1.5 unless it holds the number to a range.
var numberTokenForms = []numberTokenForm{
	// An integer of up to 64 bits.
	{`CREATE TABLE t (a BIT(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a BIGINT(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a YEAR(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a DECIMAL(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a DECIMAL(%s,0))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a DECIMAL(10,%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a NUMERIC(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a FIXED(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a FLOAT(%s,0))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a FLOAT(10,%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a DOUBLE(%s,0))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a TIME(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a DATETIME(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a TIMESTAMP(%s))`, "YYYYYYNNYNN"},
	{`SELECT CAST(a AS DECIMAL(%s))`, "YYYYYYNNYNN"},
	{`SELECT CAST(a AS DECIMAL(%s,0))`, "YYYYYYNNYNN"},
	{`SELECT CAST(a AS DECIMAL(10,%s))`, "YYYYYYNNYNN"},
	{`SELECT CAST(a AS FLOAT(%s))`, "YYYYYYNNYNN"},
	{`SELECT CAST(a AS TIME(%s))`, "YYYYYYNNYNN"},
	{`SELECT CAST(a AS DATETIME(%s))`, "YYYYYYNNYNN"},
	{`SELECT NOW(%s)`, "YYYYYYNNYNN"},
	{`SELECT CURTIME(%s)`, "YYYYYYNNYNN"},
	{`SELECT SYSDATE(%s)`, "YYYYYYNNYNN"},
	{`SELECT UTC_TIMESTAMP(%s)`, "YYYYYYNNYNN"},
	{`SELECT CURRENT_TIMESTAMP(%s)`, "YYYYYYNNYNN"},
	{`DO WEIGHT_STRING(a LEVEL %s)`, "YYYYYYNNYNN"},
	{`DO WEIGHT_STRING(a LEVEL 1 - %s)`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a POINT SRID %s)`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a DATETIME DEFAULT NOW(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a DATETIME ON UPDATE NOW(%s))`, "YYYYYYNNYNN"},
	{`CREATE INDEX i ON t (a(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a CHAR(9), KEY (a(%s)))`, "YYYYYYNNYNN"},
	{`CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE %s)`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) KEY_BLOCK_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) MAX_ROWS = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) MIN_ROWS = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) PACK_KEYS = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) STATS_AUTO_RECALC = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) STATS_PERSISTENT = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS %s`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS %s (PARTITION p0 VALUES LESS THAN (10))`, "YYYYYYNNYNN"},
	{`SELECT 1 LIMIT %s`, "YYYYYYNNYNN"},
	{`SELECT 1 LIMIT %s, 1`, "YYYYYYNNYNN"},
	{`SELECT 1 LIMIT 1 OFFSET %s`, "YYYYYYNNYNN"},
	{`DELETE FROM t LIMIT %s`, "YYYYYYNNYNN"},
	{`UPDATE t SET a = 1 LIMIT %s`, "YYYYYYNNYNN"},
	{`LOAD DATA INFILE 'f' INTO TABLE t IGNORE %s LINES`, "YYYYYYNNYNN"},
	{`LOAD DATA INFILE 'f' INTO TABLE t IGNORE %s ROWS`, "YYYYYYNNYNN"},
	{`CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR %s`, "YYYYYYNNYNN"},
	{`CREATE USER u WITH MAX_USER_CONNECTIONS %s`, "YYYYYYNNYNN"},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = %s`, "YYYYYYNNYNN"},
	{`CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = %s`, "YYYYYYNNYNN"},
	{`CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = %s`, "YYYYYYNNYNN"},
	{`ALTER TABLESPACE ts INITIAL_SIZE = %s`, "YYYYYYNNYNN"},
	{`ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = %s`, "YYYYYYNNYNN"},

	// Any number.
	{`SELECT JSON_VALUE(a, '$' DEFAULT %s ON EMPTY)`, "YYYYYYYYYYY"},
	{`SELECT JSON_VALUE(a, '$' DEFAULT -%s ON EMPTY)`, "YYYYYYYYYYY"},
	{`SELECT %s`, "YYYYYYYYYYY"},
	{`SELECT -%s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT DEFAULT %s)`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT DEFAULT -%s)`, "YYYYYYYYYYY"},
	{`ALTER TABLE t ALTER COLUMN a SET DEFAULT %s`, "YYYYYYYYYYY"},
	{`ALTER TABLE t ALTER COLUMN a SET DEFAULT -%s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) AUTO_INCREMENT = %s`, "YYYYYYYYYYY"},
	{`ALTER TABLE t AUTO_INCREMENT = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (%s))`, "YYYYYYYYYYY"},
	{`KILL %s`, "YYYYYYYYYYY"},
	{`KILL QUERY %s`, "YYYYYYYYYYY"},
	{`SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = %s`, "YYYYYYYYYYY"},
	{`SELECT SUM(a) OVER (ROWS %s PRECEDING) FROM t`, "YYYYYYYYYYY"},
	{`SELECT SUM(a) OVER (ROWS BETWEEN %s PRECEDING AND %s FOLLOWING) FROM t`, "YYYYYYYYYYY"},
	{`SELECT INTERVAL %s DAY + a`, "YYYYYYYYYYY"},
	{`SELECT a FROM t ORDER BY %s`, "YYYYYYYYYYY"},
	{`SELECT CHAR(%s USING utf8mb4)`, "YYYYYYYYYYY"},
	{`SELECT /*+ MAX_EXECUTION_TIME(%s) */ 1`, "YYYYYYYYYYY"},
	{`SET @a = %s`, "YYYYYYYYYYY"},
	{`SELECT a FROM t WHERE a IN (%s)`, "YYYYYYYYYYY"},
	{`INSERT INTO t VALUES (%s)`, "YYYYYYYYYYY"},
	{`CALL p(%s)`, "YYYYYYYYYYY"},
	{`SELECT a FROM t FOR UPDATE WAIT %s`, "YYYYYYYYYYY"},

	// An integer of up to 63 bits.
	{`DO WEIGHT_STRING(a AS CHAR(%s))`, "YYYYNNNNYNN"},
	{`DO WEIGHT_STRING(a AS BINARY(%s))`, "YYYYNNNNYNN"},
	{`CREATE TABLE t (a INT) AVG_ROW_LENGTH = %s`, "YYYYNNNNYNN"},
	{`CREATE TABLE t (a INT) CHECKSUM = %s`, "YYYYNNNNYNN"},
	{`CREATE TABLE t (a INT) DELAY_KEY_WRITE = %s`, "YYYYNNNNYNN"},
	{`HANDLER t READ FIRST LIMIT %s`, "YYYYNNNNYNN"},

	// Held to a range of its own while the statement parses.
	{`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH %s BUCKETS`, "YNNNNNNNYNN"},
}

// digitRunForms are the positions held to OceanBase 4.3.5 for the runs of digits alone, the
// first nine of numberLiterals.
var digitRunForms = []numberTokenForm{
	// Any number.
	{`CREATE TABLE t (a CHAR(%s))`, "YYYYYYYYY"},
	{`CREATE TABLE t (a VARCHAR(%s))`, "YYYYYYYYY"},
	{`CREATE TABLE t (a BINARY(%s))`, "YYYYYYYYY"},
	{`CREATE TABLE t (a VARBINARY(%s))`, "YYYYYYYYY"},
	{`CREATE TABLE t (a BLOB(%s))`, "YYYYYYYYY"},
	{`CREATE TABLE t (a TEXT(%s))`, "YYYYYYYYY"},
	{`CREATE TABLE t (a NCHAR(%s))`, "YYYYYYYYY"},
	{`CREATE TABLE t (a NATIONAL VARCHAR(%s))`, "YYYYYYYYY"},
	{`CREATE TABLE t (a CHAR VARYING(%s))`, "YYYYYYYYY"},
	{`CREATE TABLE t (a FLOAT(%s))`, "YYYYYYYYY"},
	{`SELECT CAST(a AS CHAR(%s))`, "YYYYYYYYY"},
	{`SELECT CAST(a AS NCHAR(%s))`, "YYYYYYYYY"},
	{`SELECT CAST(a AS NATIONAL CHAR(%s))`, "YYYYYYYYY"},
	{`SELECT CAST(a AS BINARY(%s))`, "YYYYYYYYY"},
	{`SELECT CONVERT(a, CHAR(%s))`, "YYYYYYYYY"},
	{`SELECT JSON_VALUE(a, '$' RETURNING CHAR(%s))`, "YYYYYYYYY"},
	{`SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(%s) PATH '$')) jt`, "YYYYYYYYY"},
	{`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR %s BEGIN END; END`, "YYYYYYYYY"},
	{`CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR %s; END`, "YYYYYYYYY"},
	{`CREATE FUNCTION f(a CHAR(%s)) RETURNS INT RETURN 1`, "YYYYYYYYY"},
}

// numberTokenCases are single statements OceanBase 4.3.5 was asked about: the edges of the
// ranges above, and the decimals a position reads by their whole part.
var numberTokenCases = []struct {
	sql    string
	accept bool
}{
	{`CREATE TABLE t (a POINT SRID 4294967295)`, true},
	{`CREATE TABLE t (a POINT SRID 4294967296)`, true},
	{`CREATE TABLE t (a POINT SRID 0x10)`, false},
	{`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 0 BUCKETS`, false},
	{`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1024 BUCKETS`, true},
	{`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1025 BUCKETS`, false},
	{`CHANGE MASTER TO MASTER_DELAY = 0x7fffffff`, false},
	{`CHANGE MASTER TO MASTER_DELAY = 0x80000000`, false},
	{`RESET MASTER TO 0x10`, false},
	{`CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 0x10`, false},
	{`CREATE TABLE t (a CHAR(0x10))`, false},
	{`CREATE TABLE t (a DECIMAL(1.5))`, false},
	{`CREATE TABLE t (a DECIMAL(99999999999999999999.5))`, false},
	{`SELECT CAST(a AS CHAR(4294967295))`, true},
	{`SELECT CAST(a AS CHAR(4294967296))`, true},
	{`SELECT CAST(a AS CHAR(9223372036854775806))`, true},
	{`SELECT CAST(a AS CHAR(18446744073709551615.5))`, true},
	{`SELECT CAST(a AS CHAR(99999999999999999999.5))`, true},
	{`CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1.5)`, false},
	{`CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 0x10`, false},
	{`CREATE USER u WITH MAX_QUERIES_PER_HOUR 0x10`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0 BEGIN END; END`, true},
	{`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 4294967296 BEGIN END; END`, true},
	{`SELECT 1 LIMIT 0x10`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384`, true},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP 1`, true},
	{`CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 0x10)`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0; END`, true},
	{`CREATE TABLE t (a DECIMAL(2147483648))`, true},
	{`CREATE TABLE t (a DECIMAL(1.5,0))`, false},
	{`CREATE TABLE t (a DECIMAL(10,1.5))`, false},
	{`CREATE TABLE t (a TIME(1.5))`, false},
	{`SELECT 1 FROM t WHERE a = 99999999999999999999 OR a = -99999999999999999999`, true},
	{`SELECT 18446744073709551616e0, 18446744073709551616.0, 0x10, 99999999999999999999abc FROM t`, true},
	{`CREATE TABLE t (a YEAR(.5))`, false},
	{`CREATE TABLE t (a DECIMAL(.5))`, false},
	{`CREATE TABLE t (a TIME(.5))`, false},
	{`CREATE TABLE t (a INT(.5))`, false},
	{`SELECT CAST(a AS CHAR(.5))`, true},
	{`SELECT CAST(a AS DECIMAL(.5))`, false},
	{`SELECT CAST(a AS CHAR(1.))`, true},
	{`CREATE TABLE t (a TIMESTAMP(1.5))`, false},
	{`CREATE TABLE t (a DECIMAL(1.5,0))`, false},
	{`CREATE TABLE t (a FLOAT(1.5,0))`, false},
	{`CREATE USER u WITH MAX_USER_CONNECTIONS 2147483647`, true},
	{`EXPLAIN FOR CONNECTION .5`, false},
	{`RESET MASTER TO .5`, false},
	{`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 254 BUCKETS`, true},
	{`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2048 BUCKETS`, false},
	{`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2049 BUCKETS`, false},
	{`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 65535 BUCKETS`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384 INITIAL_SIZE 10`, true},
	{`CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967295`, true},
	{`CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967296`, true},
	{`CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65535`, true},
	{`CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65536`, true},
	{`CREATE TABLE t (a DECIMAL(1e1))`, false},
	{`CREATE TABLE t (a CHAR(9), KEY (a(1.5)))`, false},
}

// numberTokenStatements writes each literal into each form.
func numberTokenStatements() []struct {
	sql    string
	accept bool
} {
	statements := numberTokenCases
	for _, forms := range [][]numberTokenForm{numberTokenForms, digitRunForms} {
		for _, form := range forms {
			for i, verdict := range form.accepts {
				statements = append(statements, struct {
					sql    string
					accept bool
				}{strings.ReplaceAll(form.sql, "%s", numberLiterals[i]), verdict == 'Y'})
			}
		}
	}
	return statements
}

// TestNumberTokenWidth holds the lexer to the one place the width is decided: a run
// of digits is an integer token while it fits 64 unsigned bits and the decimal token
// past that, and every number carries its whole part.
func TestNumberTokenWidth(t *testing.T) {
	for _, c := range []struct {
		text  string
		token int
		whole int64
	}{
		{"18446744073709551615", tokICONST, math.MaxInt64},
		{"18446744073709551616", tokFCONST, math.MaxInt64},
		{"000000000000000000000001", tokICONST, 1},
		{"2147483648", tokICONST, 2147483648},
		{"1.9", tokFCONST, 1},
		{"1e2", tokFCONST, 1},
		{".5", tokFCONST, 0},
	} {
		tok := NewLexer(c.text).NextToken()
		if tok.Type != c.token || tok.Ival != c.whole {
			t.Errorf("%s lexes as token %d with whole part %d, want token %d with %d", c.text, tok.Type, tok.Ival, c.token, c.whole)
		}
	}
}

func TestNumberTokenForms(t *testing.T) {
	for _, c := range numberTokenStatements() {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if got := err == nil; got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, OceanBase 4.3.5 accepts=%v: %v", c.sql, got, c.accept, err)
			}
		})
	}
}
