package parser

import "testing"

// MySQL's keyword-choice rules stand for each of their forms at every site
// that names them (mysql/parser's TestKeywordChoiceSitesMatchMySQL): CHAR SET
// for CHARACTER SET, TABLES for TABLE, NO_WAIT beside WAIT. MariaDB's
// sql_yacc.yy began as MySQL's, so this asks MariaDB 11.8.8 the same
// statements and holds omni to its answers.
func TestKeywordChoiceSitesMatchMariaDB(t *testing.T) {
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
	o := startMariaDB(t)
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			engine := o.prepareAccepts(t, sql)
			omni, err := parsesAsOneStatement(sql)
			if omni != engine {
				t.Errorf("omni accepts=%v, MariaDB accepts=%v for %q: %v", omni, engine, sql, err)
			}
		})
	}
}
