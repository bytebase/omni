package parser

import (
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestSequenceTiDBOracle runs sequence DDL against TiDB v8.5.5 in order on
// one connection and pins the result TiDB returns for each statement. omni
// must accept exactly the statements TiDB does not reject with 1064
// (ER_PARSE_ERROR). The other error numbers are DDL checks that run after a
// successful parse, which omni leaves to the engine:
//
//   - 1050: the sequence already exists.
//   - 1146: ALTER SEQUENCE IF EXISTS names a missing sequence.
//   - 4136: option values conflict (INCREMENT 0, MINVALUE at the int64 floor).
//   - 4139: DROP SEQUENCE names a missing sequence.
//   - 8227: a table option other than COMMENT or ENGINE.
func TestSequenceTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	// USE is session state: a dedicated connection keeps every statement on
	// the database it selects, and the cases depend on each other's DDL.
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatalf("acquire connection: %v", err)
	}
	defer conn.Close()

	conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_seq")
	if _, err := conn.ExecContext(tc.ctx, "CREATE DATABASE omni_seq"); err != nil {
		t.Fatalf("create database: %v", err)
	}
	defer conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_seq")
	if _, err := conn.ExecContext(tc.ctx, "USE omni_seq"); err != nil {
		t.Fatalf("use database: %v", err)
	}

	cases := []struct {
		sql     string
		tidbErr uint16 // MySQL error number TiDB returns; 0 = success
	}{
		// CREATE SEQUENCE.
		{"CREATE SEQUENCE seq_min", 0},
		{"CREATE SEQUENCE IF NOT EXISTS seq_min", 0},
		{"CREATE SEQUENCE seq_min", 1050},
		{"CREATE SEQUENCE omni_seq.seq_qualified", 0},
		{"CREATE SEQUENCE `seq quoted`", 0},
		{"CREATE SEQUENCE seq_by INCREMENT BY 2 START WITH 10 MINVALUE 1 MAXVALUE 1000 CACHE 20 CYCLE", 0},
		{"CREATE SEQUENCE seq_eq INCREMENT = 2 START = 10 MINVALUE = 1 MAXVALUE = 1000 CACHE = 20 NOCYCLE", 0},
		{"CREATE SEQUENCE seq_bare INCREMENT 2 START 10 MINVALUE 1 MAXVALUE 1000 CACHE 20 NO CYCLE", 0},
		{"CREATE SEQUENCE seq_no1 NOMINVALUE NOMAXVALUE NOCACHE NOCYCLE", 0},
		{"CREATE SEQUENCE seq_no2 NO MINVALUE NO MAXVALUE NO CACHE NO CYCLE", 0},
		{"CREATE SEQUENCE seq_neg INCREMENT BY -1 START WITH -1 MINVALUE -100 MAXVALUE -1", 0},
		{"CREATE SEQUENCE seq_plus INCREMENT BY +1 START WITH +5", 0},
		{"CREATE SEQUENCE seq_repeat START WITH 1 START WITH 2", 0},
		{"CREATE SEQUENCE seq_comment COMMENT 'ids'", 0},
		{"CREATE SEQUENCE seq_tblopts START 1 COMMENT = 'ids', ENGINE = InnoDB", 0},
		{"CREATE SEQUENCE seq_engine ENGINE InnoDB COMMENT 'a' COMMENT 'b'", 0},
		{"CREATE SEQUENCE seq_ai AUTO_INCREMENT = 5", 8227},
		{"CREATE SEQUENCE seq_cs DEFAULT CHARSET = utf8mb4", 8227},
		{"CREATE SEQUENCE seq_zero INCREMENT BY 0", 4136},
		{"CREATE SEQUENCE seq_floor MINVALUE -9223372036854775808", 4136},
		{"CREATE OR REPLACE SEQUENCE seq_x", 1064},
		{"CREATE TEMPORARY SEQUENCE seq_x", 1064},
		{"CREATE SEQUENCE IF EXISTS seq_x", 1064},
		{"CREATE SEQUENCE seq_x AS BIGINT", 1064},
		{"CREATE SEQUENCE seq_x RESTART", 1064},
		{"CREATE SEQUENCE seq_x CACHE", 1064},
		{"CREATE SEQUENCE seq_x MAXVALUE", 1064},
		{"CREATE SEQUENCE seq_x NO", 1064},
		{"CREATE SEQUENCE seq_x NO START", 1064},
		{"CREATE SEQUENCE seq_x START WITH 1.5", 1064},
		{"CREATE SEQUENCE seq_x START WITH (1)", 1064},
		{"CREATE SEQUENCE seq_x START WITH 1+1", 1064},
		{"CREATE SEQUENCE seq_x START WITH '1'", 1064},
		{"CREATE SEQUENCE seq_x START WITH 0x10", 1064},
		{"CREATE SEQUENCE seq_x START WITH 9223372036854775808", 1064},
		{"CREATE SEQUENCE seq_x MINVALUE -9223372036854775809", 1064},
		{"CREATE SEQUENCE seq_x INCREMENT BY = 1", 1064},
		{"CREATE SEQUENCE seq_x START WITH = 1", 1064},
		{"CREATE SEQUENCE seq_x COMMENT 'x' START 1", 1064},
		{"CREATE SEQUENCE seq_x COMMENT 'x',", 1064},
		{"CREATE SEQUENCE seq_x , COMMENT 'x'", 1064},
		{"CREATE SEQUENCE seq_x START TRANSACTION", 1064},
		{"CREATE SEQUENCE seq_x PARTITION BY HASH(id)", 1064},

		// ALTER SEQUENCE.
		{"ALTER SEQUENCE seq_min RESTART", 0},
		{"ALTER SEQUENCE IF EXISTS seq_min RESTART WITH 5", 0},
		{"ALTER SEQUENCE IF EXISTS seq_missing RESTART", 1146},
		{"ALTER SEQUENCE seq_min RESTART = 5", 0},
		{"ALTER SEQUENCE seq_min RESTART 5", 0},
		{"ALTER SEQUENCE seq_neg RESTART -5", 0},
		{"ALTER SEQUENCE seq_min RESTART RESTART WITH 3", 0},
		{"ALTER SEQUENCE seq_min INCREMENT BY 3 CACHE 100 NOCYCLE", 0},
		{"ALTER SEQUENCE seq_min START WITH 3 RESTART", 0},
		{"ALTER SEQUENCE seq_min NO MINVALUE NO MAXVALUE NO CACHE NO CYCLE", 0},
		{"ALTER SEQUENCE seq_min", 1064},
		{"ALTER SEQUENCE IF NOT EXISTS seq_min RESTART", 1064},
		{"ALTER SEQUENCE seq_min RESTART WITH", 1064},
		{"ALTER SEQUENCE seq_min RESTART WITH = 5", 1064},
		{"ALTER SEQUENCE seq_min RESTART, CACHE 5", 1064},
		{"ALTER SEQUENCE seq_min COMMENT 'x'", 1064},
		{"ALTER SEQUENCE seq_min ENGINE = InnoDB", 1064},
		{"ALTER SEQUENCE seq_min AS INT", 1064},

		// DROP SEQUENCE.
		{"DROP SEQUENCE seq_by", 0},
		{"DROP SEQUENCE IF EXISTS seq_by", 0},
		{"DROP SEQUENCE seq_eq, seq_bare, omni_seq.seq_no1", 0},
		{"DROP SEQUENCE IF EXISTS seq_no2, seq_missing", 0},
		{"DROP SEQUENCE seq_missing", 4139},
		{"DROP SEQUENCE", 1064},
		{"DROP SEQUENCE seq_min,", 1064},
		{"DROP TEMPORARY SEQUENCE seq_min", 1064},
		{"DROP SEQUENCE seq_min CASCADE", 1064},
		{"DROP SEQUENCE seq_min RESTRICT", 1064},
		{"DROP SEQUENCE IF NOT EXISTS seq_min", 1064},

		// The sequence keywords are unreserved upstream.
		{"CREATE SEQUENCE sequence", 0},
		{"CREATE SEQUENCE nocache", 0},
		{"ALTER SEQUENCE nocache RESTART", 0},
		{"DROP SEQUENCE sequence, nocache", 0},
		{"CREATE TABLE cycle (increment INT, minvalue INT, nominvalue INT, nomaxvalue INT, nocache INT, nocycle INT, sequence INT)", 0},
		{"SELECT increment, minvalue, nominvalue, nomaxvalue, nocache, nocycle, sequence FROM cycle", 0},
	}

	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			_, tidbErr := conn.ExecContext(tc.ctx, c.sql)
			var got uint16
			if tidbErr != nil {
				var myErr *gomysql.MySQLError
				if !errors.As(tidbErr, &myErr) {
					t.Fatalf("TiDB returned a non-MySQL error: %v", tidbErr)
				}
				got = myErr.Number
			}
			if got != c.tidbErr {
				t.Errorf("TiDB error number = %d, want %d (err=%v)", got, c.tidbErr, tidbErr)
			}

			_, omniErr := Parse(c.sql)
			if wantAccept := c.tidbErr != 1064; (omniErr == nil) != wantAccept {
				t.Errorf("omni accepts=%v (err=%v), TiDB error number %d", omniErr == nil, omniErr, got)
			}
		})
	}
}
