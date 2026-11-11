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

// numberTokenForm is a statement with %s where a number goes, and whether MariaDB 11.8
// parses it with each of numberLiterals there, in order: Y or N.
type numberTokenForm struct {
	sql     string
	accepts string
}

// numberTokenForms are the positions held to MariaDB 11.8 for every literal. The lexer
// gives a run of digits past 64 bits the token of a decimal, so a position answers
// for 18446744073709551616 as it does for 1.5 unless it holds the number to a range.
var numberTokenForms = []numberTokenForm{
	// An integer of up to 64 bits.
	{`SELECT NOW(%s)`, "YYYYYYNNYNN"},
	{`SELECT CURTIME(%s)`, "YYYYYYNNYNN"},
	{`SELECT SYSDATE(%s)`, "YYYYYYNNYNN"},
	{`SELECT UTC_TIMESTAMP(%s)`, "YYYYYYNNYNN"},
	{`SELECT CURRENT_TIMESTAMP(%s)`, "YYYYYYNNYNN"},
	{`DO WEIGHT_STRING(a AS CHAR(%s))`, "YYYYYYNNYNN"},
	{`DO WEIGHT_STRING(a AS BINARY(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a DATETIME DEFAULT NOW(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a DATETIME ON UPDATE NOW(%s))`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = %s)`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = %s)`, "YYYYYYNNYNN"},
	{`CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = %s)`, "YYYYYYNNYNN"},
	{`ALTER TABLE t ADD PARTITION PARTITIONS %s`, "YYYYYYNNYNN"},
	{`ALTER TABLE t COALESCE PARTITION %s`, "YYYYYYNNYNN"},
	{`SELECT 1 LIMIT %s`, "YYYYYYNNYNN"},
	{`SELECT 1 LIMIT %s, 1`, "YYYYYYNNYNN"},
	{`SELECT 1 LIMIT 1 OFFSET %s`, "YYYYYYNNYNN"},
	{`DELETE FROM t LIMIT %s`, "YYYYYYNNYNN"},
	{`UPDATE t SET a = 1 LIMIT %s`, "YYYYYYNNYNN"},
	{`HANDLER t READ FIRST LIMIT %s`, "YYYYYYNNYNN"},
	{`SHOW WARNINGS LIMIT %s`, "YYYYYYNNYNN"},
	{`SHOW BINLOG EVENTS LIMIT %s`, "YYYYYYNNYNN"},
	{`SHOW PROFILE LIMIT %s`, "YYYYYYNNYNN"},

	// Any number.
	{`SELECT %s`, "YYYYYYYYYYY"},
	{`SELECT -%s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT DEFAULT %s)`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT DEFAULT -%s)`, "YYYYYYYYYYY"},
	{`ALTER TABLE t ALTER COLUMN a SET DEFAULT %s`, "YYYYYYYYYYY"},
	{`ALTER TABLE t ALTER COLUMN a SET DEFAULT -%s`, "YYYYYYYYYYY"},
	{`CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE %s)`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) AUTO_INCREMENT = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) AVG_ROW_LENGTH = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) CHECKSUM = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) TABLE_CHECKSUM = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) DELAY_KEY_WRITE = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) KEY_BLOCK_SIZE = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) MAX_ROWS = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) MIN_ROWS = %s`, "YYYYYYYYYYY"},
	{`ALTER TABLE t AUTO_INCREMENT = %s`, "YYYYYYYYYYY"},
	{`CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (%s))`, "YYYYYYYYYYY"},
	{`SHOW BINLOG EVENTS FROM %s`, "YYYYYYYYYYY"},
	{`SHOW RELAYLOG EVENTS FROM %s`, "YYYYYYYYYYY"},
	{`KILL %s`, "YYYYYYYYYYY"},
	{`KILL QUERY %s`, "YYYYYYYYYYY"},
	{`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR %s BEGIN END; END`, "YYYYYYYYYYY"},
	{`CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR %s; END`, "YYYYYYYYYYY"},
	{`SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = %s`, "YYYYYYYYYYY"},
	{`GET DIAGNOSTICS CONDITION %s @a = MESSAGE_TEXT`, "YYYYYYYYYYY"},
	{`CHANGE MASTER TO IGNORE_SERVER_IDS = (%s)`, "YYYYYYYYYYY"},
	{`CHANGE MASTER TO MASTER_PORT = %s`, "YYYYYYYYYYY"},
	{`CHANGE MASTER TO MASTER_LOG_POS = %s`, "YYYYYYYYYYY"},
	{`CHANGE MASTER TO RELAY_LOG_POS = %s`, "YYYYYYYYYYY"},
	{`CHANGE MASTER TO MASTER_CONNECT_RETRY = %s`, "YYYYYYYYYYY"},
	{`CHANGE MASTER TO MASTER_SSL = %s`, "YYYYYYYYYYY"},
	{`START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = %s`, "YYYYYYYYYYY"},
	{`START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = %s`, "YYYYYYYYYYY"},
	{`RESET MASTER TO %s`, "YYYYYYYYYYY"},
	{`CREATE USER u WITH MAX_QUERIES_PER_HOUR %s`, "YYYYYYYYYYY"},
	{`CREATE USER u WITH MAX_UPDATES_PER_HOUR %s`, "YYYYYYYYYYY"},
	{`CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR %s`, "YYYYYYYYYYY"},
	{`ALTER USER u WITH MAX_QUERIES_PER_HOUR %s`, "YYYYYYYYYYY"},
	{`SELECT SUM(a) OVER (ROWS %s PRECEDING) FROM t`, "YYYYYYYYYYY"},
	{`SELECT SUM(a) OVER (ROWS BETWEEN %s PRECEDING AND %s FOLLOWING) FROM t`, "YYYYYYYYYYY"},
	{`SELECT INTERVAL %s DAY + a`, "YYYYYYYYYYY"},
	{`CREATE EVENT e ON SCHEDULE EVERY %s DAY DO SELECT 1`, "YYYYYYYYYYY"},
	{`SELECT a FROM t ORDER BY %s`, "YYYYYYYYYYY"},
	{`SELECT CHAR(%s USING utf8mb4)`, "YYYYYYYYYYY"},
	{`SELECT /*+ MAX_EXECUTION_TIME(%s) */ 1`, "YYYYYYYYYYY"},
	{`SET @a = %s`, "YYYYYYYYYYY"},
	{`SELECT a FROM t WHERE a IN (%s)`, "YYYYYYYYYYY"},
	{`INSERT INTO t VALUES (%s)`, "YYYYYYYYYYY"},
	{`CALL p(%s)`, "YYYYYYYYYYY"},
	{`SELECT a FROM t FOR UPDATE WAIT %s`, "YYYYYYYYYYY"},
	{`CHANGE MASTER TO IGNORE_SERVER_IDS = (1, %s)`, "YYYYYYYYYYY"},

	// An integer of up to 31 bits.
	{`CREATE INDEX i ON t (a(%s))`, "YYNNNNNNYNN"},
	{`CREATE TABLE t (a CHAR(9), KEY (a(%s)))`, "YYNNNNNNYNN"},
	{`LOAD DATA INFILE 'f' INTO TABLE t IGNORE %s LINES`, "YYNNNNNNYNN"},
	{`LOAD DATA INFILE 'f' INTO TABLE t IGNORE %s ROWS`, "YYNNNNNNYNN"},
	{`SHOW PROFILE FOR QUERY %s`, "YYNNNNNNYNN"},
	{`CREATE USER u PASSWORD EXPIRE INTERVAL %s DAY`, "YYNNNNNNYNN"},
	{`CREATE USER u WITH MAX_USER_CONNECTIONS %s`, "YYNNNNNNYNN"},

	// An integer or a decimal.
	{`CREATE TABLE t (a YEAR(%s))`, "YYYYYYYYYYN"},

	// Held to a range of its own while the statement parses.
	{`CREATE TABLE t (a DECIMAL(%s,0))`, "YNNNNNNNYNN"},
	{`CREATE TABLE t (a DECIMAL(10,%s))`, "YNNNNNNNYNN"},
	{`CREATE TABLE t (a FLOAT(%s,0))`, "YNNNNNNNYNN"},
	{`CREATE TABLE t (a FLOAT(10,%s))`, "YNNNNNNNYNN"},
	{`CREATE TABLE t (a DOUBLE(%s,0))`, "YNNNNNNNYNN"},
	{`CREATE TABLE t (a DOUBLE(10,%s))`, "YNNNNNNNYNN"},
	{`CREATE TABLE t (a REAL(10,%s))`, "YNNNNNNNYNN"},
	{`SELECT CAST(a AS DECIMAL(%s,0))`, "YNNNNNNNYNN"},
	{`SELECT CAST(a AS DECIMAL(10,%s))`, "YNNNNNNNYNN"},
	{`CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = %s (a)`, "YNNNNNNNYNN"},
	{`CREATE TABLE t (a CHAR(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a BINARY(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a BIT(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a INT(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a BIGINT(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a NCHAR(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a DECIMAL(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a NUMERIC(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a FIXED(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a FLOAT(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a TIME(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a DATETIME(%s))`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a TIMESTAMP(%s))`, "YNNNNNNNYYN"},
	{`SELECT CAST(a AS DECIMAL(%s))`, "YNNNNNNNYYN"},
	{`SELECT CAST(a AS TIME(%s))`, "YNNNNNNNYYN"},
	{`SELECT CAST(a AS DATETIME(%s))`, "YNNNNNNNYYN"},
	{`SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(%s) PATH '$')) jt`, "YNNNNNNNYYN"},
	{`CREATE FUNCTION f(a CHAR(%s)) RETURNS INT RETURN 1`, "YNNNNNNNYYN"},
	{`CREATE TABLE t (a INT) PACK_KEYS = %s`, "YNNNNNNNYYY"},
	{`CREATE TABLE t (a INT) STATS_AUTO_RECALC = %s`, "YNNNNNNNYYY"},
	{`CREATE TABLE t (a INT) STATS_PERSISTENT = %s`, "YNNNNNNNYYY"},
	{`CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = %s`, "YNNNNNNNYYY"},
	{`CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = %s`, "YNNNNNNNYYY"},
	{`CHANGE MASTER TO MASTER_DELAY = %s`, "YYNNNNNNYYY"},
	{`CREATE TABLE t (a VARCHAR(%s))`, "YYYNNNNNYYN"},
	{`CREATE TABLE t (a VARBINARY(%s))`, "YYYNNNNNYYN"},
	{`CREATE TABLE t (a BLOB(%s))`, "YYYNNNNNYYN"},
	{`CREATE TABLE t (a TEXT(%s))`, "YYYNNNNNYYN"},
	{`CREATE TABLE t (a NATIONAL VARCHAR(%s))`, "YYYNNNNNYYN"},
	{`CREATE TABLE t (a CHAR VARYING(%s))`, "YYYNNNNNYYN"},
	{`SELECT CAST(a AS CHAR(%s))`, "YYYNNNNNYYN"},
	{`SELECT CAST(a AS NCHAR(%s))`, "YYYNNNNNYYN"},
	{`SELECT CAST(a AS BINARY(%s))`, "YYYNNNNNYYN"},
	{`SELECT CONVERT(a, CHAR(%s))`, "YYYNNNNNYYN"},
}

// digitRunForms are the positions held to MariaDB 11.8 for the runs of digits alone, the
// first nine of numberLiterals.
var digitRunForms = []numberTokenForm{
	// Any number.
	{`EXPLAIN FOR CONNECTION %s`, "YYYYYYYYY"},
}

// numberTokenCases are single statements MariaDB 11.8 was asked about: the edges of the
// ranges above, and the decimals a position reads by their whole part.
var numberTokenCases = []struct {
	sql    string
	accept bool
}{
	{`CREATE TABLE t (a POINT SRID 0x10)`, false},
	{`CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 4294967`, true},
	{`CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 4294967.5`, false},
	{`CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 4294968`, false},
	{`CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 0`, true},
	{`CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 0.0001`, true},
	{`CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 99999999999999999999.5`, false},
	{`CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 1e30`, false},
	{`CHANGE MASTER TO MASTER_DELAY = 2147483647.9`, true},
	{`CHANGE MASTER TO MASTER_DELAY = 2147483648.0`, false},
	{`CHANGE MASTER TO MASTER_DELAY = 1e30`, true},
	{`CHANGE MASTER TO MASTER_DELAY = 0x80000000`, false},
	{`RESET MASTER TO 0`, true},
	{`RESET MASTER TO 2000000000`, true},
	{`RESET MASTER TO 2000000001`, true},
	{`RESET MASTER TO 2147483646`, true},
	{`CREATE TABLE t (a CHAR(0x10))`, false},
	{`CREATE TABLE t (a CHAR(1.))`, true},
	{`CREATE TABLE t (a CHAR(.5))`, false},
	{`CREATE TABLE t (a CHAR(99999999999999999999.5))`, false},
	{`CREATE TABLE t (a DECIMAL(1.5))`, true},
	{`CREATE TABLE t (a FLOAT(1.5))`, true},
	{`CREATE TABLE t (a DECIMAL(99999999999999999999.5))`, false},
	{`SELECT CAST(a AS CHAR(4294967295))`, true},
	{`SELECT CAST(a AS CHAR(4294967296))`, false},
	{`SELECT CAST(a AS CHAR(9223372036854775806))`, false},
	{`SELECT CAST(a AS CHAR(18446744073709551615.5))`, false},
	{`SELECT CAST(a AS CHAR(99999999999999999999.5))`, false},
	{`CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 18446744073709551615K`, false},
	{`CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 99999999999999999999K`, false},
	{`CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 4M`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 99999999999999999999K`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 'x'`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 4K`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 4K`, false},
	{`CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1.5)`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0 BEGIN END; END`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0.5 BEGIN END; END`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 4294967296 BEGIN END; END`, true},
	{`SELECT 1 LIMIT 0x10`, false},
	{`CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483647K`, false},
	{`CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483648K`, false},
	{`CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483648G`, false},
	{`CREATE TABLE t (a INT) AUTOEXTEND_SIZE = K`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483647M`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483648M`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE 4M`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE 16384`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP 1`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE 4M MAX_SIZE 1G EXTENT_SIZE 1M`, false},
	{`ALTER TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M`, false},
	{`CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE 4M UNDO_BUFFER_SIZE 4M`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0; END`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0.9; END`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1.9; END`, true},
	{`CREATE TABLE t (a DECIMAL(2147483648))`, false},
	{`CREATE TABLE t (a DECIMAL(1.5,0))`, false},
	{`CREATE TABLE t (a DECIMAL(10,1.5))`, false},
	{`CREATE TABLE t (a TIME(1.5))`, true},
	{`CREATE TABLE t (a CHAR(1.5) CHARACTER SET utf8mb4)`, true},
	{`CREATE FUNCTION f(a CHAR(1.5)) RETURNS VARCHAR(2.5) RETURN a`, true},
	{`SELECT 1 FROM t WHERE a = 99999999999999999999 OR a = -99999999999999999999`, true},
	{`SELECT 18446744073709551616e0, 18446744073709551616.0, 0x10, 99999999999999999999abc FROM t`, true},
	{`CREATE TABLE t (a DECIMAL(.5))`, false},
	{`CREATE TABLE t (a TIME(.5))`, false},
	{`CREATE TABLE t (a VARCHAR(.5))`, false},
	{`CREATE TABLE t (a INT(.5))`, false},
	{`SELECT CAST(a AS CHAR(.5))`, false},
	{`SELECT CAST(a AS DECIMAL(.5))`, false},
	{`SELECT CAST(a AS CHAR(1.))`, true},
	{`CREATE TABLE t (a TIMESTAMP(1.5))`, true},
	{`CREATE TABLE t (a DECIMAL(1.5,0))`, false},
	{`CREATE TABLE t (a FLOAT(1.5,0))`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR .5 BEGIN END; END`, false},
	{`CREATE USER u WITH MAX_USER_CONNECTIONS 2147483647`, true},
	{`CREATE USER u PASSWORD EXPIRE INTERVAL 0 DAY`, false},
	{`RESET MASTER TO .5`, true},
	{`RESET MASTER TO 0`, true},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 'x'`, false},
	{`CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384 INITIAL_SIZE 10`, false},
	{`CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967295`, true},
	{`CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967296`, true},
	{`CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65535`, true},
	{`CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65536`, true},
	{`CHANGE MASTER TO MASTER_DELAY = 2147483647`, true},
	{`CHANGE MASTER TO MASTER_DELAY = 2147483648`, false},
	{`CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1.5 BEGIN END; END`, true},
	{`CREATE TABLE t (a VARCHAR(1e1))`, false},
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
				t.Errorf("Parse(%q) accepts=%v, MariaDB 11.8 accepts=%v: %v", c.sql, got, c.accept, err)
			}
		})
	}
}
