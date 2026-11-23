package parser

import "testing"

// A keyword choice the grammar names as a rule of its own stands for each of
// its forms at every site that names it: character_set is CHAR SET (CHAR_SYM
// is CHAR or CHARACTER in lex.h) or CHARSET, table_or_tables is TABLE or
// TABLES, ts_option_wait is WAIT or NO_WAIT. The generated compat rows reach
// some of these sites; this asks MySQL 8.0 about the rest, and about the
// spellings next to them it refuses.
func TestKeywordChoiceSitesMatchMySQL(t *testing.T) {
	cases := []string{
		`CREATE DATABASE d CHAR SET utf8mb4`,
		`CREATE DATABASE d DEFAULT CHAR SET = utf8mb4`,
		`ALTER DATABASE d CHAR SET utf8mb4`,
		`ALTER DATABASE d CHAR utf8mb4`,
		`CREATE TABLE t (a INT) CHAR SET utf8mb4`,
		`CREATE TABLE t (a INT) DEFAULT CHAR SET = utf8mb4`,
		`CREATE TABLE t (a INT) CHAR`,
		`ALTER TABLE t CHAR SET utf8mb4`,
		`LOAD DATA INFILE 'f' INTO TABLE t CHAR SET utf8mb4`,
		`LOAD DATA INFILE 'f' INTO TABLE t CHARACTER utf8mb4`,
		`SELECT 1 INTO OUTFILE 'f' CHAR SET utf8mb4`,
		`SELECT 1 INTO OUTFILE 'f' CHARACTER utf8mb4`,
		`SHOW CHAR SET`,
		`DROP TABLES t`,
		`DROP TEMPORARY TABLES IF EXISTS t, u`,
		`ANALYZE TABLES t`,
		`ANALYZE LOCAL TABLES t`,
		`CHECK TABLES t`,
		`OPTIMIZE TABLES t`,
		`REPAIR TABLES t`,
		`CREATE LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT`,
		`CREATE LOGFILE GROUP g ADD UNDOFILE 'u' WAIT NO_WAIT`,
		`ALTER LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT`,
	}
	o := startParserOracle(t)
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertOracleMatch(t, o, sql)
		})
	}
}
