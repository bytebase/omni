package parser

import (
	"testing"
)

// MySQL's keyword-choice rules stand for each of their forms at every site
// that names them (mysql/parser's TestKeywordChoiceSitesMatchMySQL): CHAR SET
// for CHARACTER SET, TABLES for TABLE, NO_WAIT beside WAIT. This asks the
// pinned TiDB the same statements, through PREPARE as the compat judge does,
// and holds omni to its answers.
func TestKeywordChoiceSitesMatchTiDB(t *testing.T) {
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
		// TiDB's INTO OUTFILE takes no character set at all (parser.y:10333).
		`SELECT 1 INTO OUTFILE 'f' CHARACTER SET utf8mb4`,
		`SELECT 1 INTO OUTFILE 'f' CHARSET utf8mb4`,
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
	tc := startTiDB(t)
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatalf("connect to TiDB: %v", err)
	}
	defer conn.Close()
	prepare := func(text string) error {
		if _, err := conn.ExecContext(tc.ctx, "SET @omni_probe = ?", text); err != nil {
			t.Fatalf("binding the statement: %v", err)
		}
		_, err := conn.ExecContext(tc.ctx, "PREPARE omni_probe_stmt FROM @omni_probe")
		return err
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			engine, err := compatJudgeVerdict(sql, prepare)
			if err != nil {
				t.Fatalf("no verdict for %q: %v", sql, err)
			}
			if omni := omniVerdict(sql); omni != engine {
				t.Errorf("omni %s, TiDB %s for %q", omni, engine, sql)
			}
		})
	}
}
