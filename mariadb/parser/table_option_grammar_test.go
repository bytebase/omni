package parser

import "testing"

// TestMariaDBTableAndDatabaseOptionGrammar holds MariaDB's own table and
// database options to sql_yacc.yy, as MariaDB 11.8.8 answers them.
func TestMariaDBTableAndDatabaseOptionGrammar(t *testing.T) {
	for _, sql := range []string{
		"ALTER TABLE t PAGE_CHECKSUM 1",
		"ALTER TABLE t PAGE_CHECKSUM = +1.5",
		"ALTER TABLE t TRANSACTIONAL DEFAULT",
		"ALTER TABLE t SEQUENCE = 0x1",
		"ALTER TABLE t READ_ONLY = 's'",
		"ALTER TABLE t WRAPPER := x",
		"ALTER TABLE t page_compressed = 0x1",
		"ALTER TABLE t WITH SYSTEM VERSIONING",
		"CREATE TABLE t SELECT 1 a WITH SYSTEM VERSIONING",
		"ALTER ONLINE IGNORE ONLINE TABLE t ADD c INT",
		"ALTER DATABASE COMMENT 'c'",
		"ALTER DATABASE COMMENT = 'c' CHARACTER SET utf8mb4",
		"ALTER DATABASE d COMMENT 'c'",
		"ALTER DATABASE comment COMMENT 'c'",
		"CREATE DATABASE d COMMENT 'c'",
		"ALTER DATABASE d UPGRADE DATA DIRECTORY NAME",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"ALTER TABLE t PAGE_CHECKSUM 'x'",
		"ALTER TABLE t TRANSACTIONAL",
		"ALTER TABLE t READ_ONLY 's'",
		"ALTER TABLE t page_compressed = 1.5",
		"ALTER TABLE t page_compressed = NULL",
		"ALTER ONLINE VIEW v AS SELECT 1",
		"ALTER ONLINE",
		"ALTER DATABASE COMMENT",
		"ALTER DATABASE d UPGRADE DATA DIRECTORY",
		"ALTER DATABASE UPGRADE DATA DIRECTORY NAME",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): want an error", sql)
		}
	}
}
