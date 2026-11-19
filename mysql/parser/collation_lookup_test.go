package parser

import (
	"errors"
	"testing"
)

// TestCollationLookupCode holds the code the parser records for a collation name
// MySQL resolves while it parses: 1273, ER_UNKNOWN_COLLATION, for a name no
// server of the pin has. The server raises it at the name, so nothing after it
// is read and `ALTER TABLE a COLLATE utf8mb4 zzz` is 1273 rather than the 1064
// the trailing zzz is. `utf8mb4` is the case the catalog meets: it is a
// character set, and no collation, so `ALTER TABLE a COLLATE utf8mb4` is 1273
// before the 1146 the missing table a would answer.
//
// Measured on mysql:8.0.46: every statement here answers its code, and so does
// the same statement with a syntax error appended. The verdict stays a
// rejection, which is what compat.tsv records: its judge counts a lookup error
// as parsed, so these rows are accept/accept there and the catalog answers the
// code (see mysql/catalog/simulationOmni).
func TestCollationLookupCode(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"ALTER TABLE a COLLATE utf8mb4", 1273},
		{"ALTER TABLE a COLLATE utf8mb4 zzz", 1273},
		{"ALTER TABLE a ADD a INT COLLATE utf8mb4", 1273},
		{"ALTER TABLE a ADD a INT COLLATE utf8mb4 zzz", 1273},
		{"ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4", 1273},
		{"CREATE TABLE a (b VARCHAR(3) COLLATE utf8mb4)", 1273},
		{"CREATE TABLE a (b INT) COLLATE utf8mb4", 1273},
		{"CREATE TABLE a (b INT) COLLATE = utf8mb4_nosuch", 1273},
		{"CREATE DATABASE d COLLATE utf8mb4", 1273},
		{"ALTER DATABASE d COLLATE BINLOG", 1273},
		{"SELECT 1 COLLATE utf8mb4", 1273}, {"SET NAMES utf8mb4 COLLATE utf8mb4", 1273},

		// A name the server has, or one the grammar refuses before the name is
		// looked up, is the trailing zzz's 1064.
		{"ALTER TABLE a COLLATE utf8mb4_bin zzz", 0},
		{"ALTER TABLE a COLLATE UTF8MB4_BIN zzz", 0},
		{"ALTER TABLE a ADD a INT COLLATE BINARY zzz", 0},
		{"ALTER TABLE a ADD a INT COLLATE DEFAULT zzz", 0},
		{"CREATE TABLE a (b INT) COLLATE utf8mb4_bin zzz", 0},
		{"CREATE DATABASE d COLLATE utf8mb4_bin zzz", 0},
		{"SELECT 1 COLLATE BINARY", 0},
		{"SET NAMES utf8mb4 COLLATE utf8mb4_bin zzz", 0},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("%s parsed, want the code %d", tc.sql, tc.code)
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("%s: %v is no ParseError", tc.sql, err)
			}
			if pe.Code != tc.code {
				t.Fatalf("%s answered the code %d (%s), want %d", tc.sql, pe.Code, pe.Message, tc.code)
			}
			if tc.code != 0 && pe.Message != "Unknown collation: '"+collationOf(tc.sql)+"'" {
				t.Fatalf("%s answered %q, want the server's message", tc.sql, pe.Message)
			}
		})
	}
}

// collationOf returns the name after the statement's first COLLATE, with the
// optional '=' dropped, which is the name the test's message must quote.
func collationOf(sql string) string {
	for i := 0; i+7 < len(sql); i++ {
		if sql[i:i+7] != "COLLATE" {
			continue
		}
		rest := sql[i+7:]
		for len(rest) > 0 && (rest[0] == ' ' || rest[0] == '=') {
			rest = rest[1:]
		}
		name := ""
		for len(rest) > 0 && rest[0] != ' ' && rest[0] != ')' {
			name += string(rest[0])
			rest = rest[1:]
		}
		return name
	}
	return ""
}
