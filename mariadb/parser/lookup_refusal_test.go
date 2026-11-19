package parser

import (
	"errors"
	"testing"
)

// TestLookupRefusalCode holds the code the parser records for the two names MariaDB
// looks up while it parses: a collation of another character set than the one a type
// settled (1253) and a FULLTEXT key's parser plugin, which the pinned server does not
// ship (1128). The server raises the code at the name, so nothing after it is read and
// `COLLATE utf8mb4_bin zzz` is 1253 rather than the 1064 the trailing zzz is. Measured
// on mariadb:11.8.8: every statement here answers its code through PREPARE, and one
// that appends a syntax error to it answers the same code again.
//
// The verdict stays a rejection, which is what compat.tsv records: its judge counts a
// lookup error as parsed, so these rows are accept/accept or a known gap there, and
// the catalog answers the code (see catalog.validateDDL).
func TestLookupRefusalCode(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin zzz", 1253},
		{"ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY zzz", 1253},
		{"ALTER TABLE a ADD a CHAR BYTE COLLATE utf8mb4_bin zzz", 1253},
		{"ALTER TABLE a ADD a CHAR UNICODE COLLATE utf8mb4_bin zzz", 1253},
		{"ALTER TABLE a ADD a NCHAR COLLATE utf8mb4_bin zzz", 1253},
		{"ALTER TABLE a ADD a NVARCHAR(3) COLLATE utf8mb4_bin zzz", 1253},
		{"ALTER TABLE a ADD a NATIONAL CHAR COLLATE utf8mb4_bin zzz", 1253},
		{"ALTER TABLE a ADD a ENUM('x') ASCII COLLATE utf8mb4_bin zzz", 1253},
		{"CREATE TABLE z (a CHAR ASCII COLLATE utf8mb4_bin) zzz", 1253},
		{"ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a zzz", 1128},
		{"CREATE TABLE z (a TEXT, FULLTEXT KEY k (a) WITH PARSER p) zzz", 1128},
		{"CREATE TABLE z (a TEXT, FULLTEXT(a) WITH PARSER p) zzz", 1128},
		{"CREATE FULLTEXT INDEX i ON a(a) WITH PARSER p zzz", 1128},

		// The name is one the type's own character set has, or the type settled
		// none, so nothing is refused and the trailing zzz is the 1064.
		{"ALTER TABLE a ADD a CHAR CHARACTER SET ascii COLLATE ascii_bin zzz", 0},
		{"ALTER TABLE a ADD a CHAR CHARACTER SET utf8 COLLATE utf8mb3_bin zzz", 0},
		{"ALTER TABLE a ADD a CHAR COLLATE utf8mb4_bin zzz", 0},
		{"ALTER TABLE a ADD a CHAR ASCII COLLATE DEFAULT zzz", 0},
		{"ALTER TABLE a ADD a INT COLLATE utf8mb4_bin zzz", 0},
		// A non-FULLTEXT index has no parser at all: WITH PARSER after it is the
		// server's 1064, not a lookup.
		{"CREATE INDEX i ON a(a) WITH PARSER p", 0},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if err == nil {
				if tc.code != 0 {
					t.Fatalf("%s parsed, want the code %d", tc.sql, tc.code)
				}
				return
			}
			if !errors.As(err, &pe) {
				t.Fatalf("%s: %v is no ParseError", tc.sql, err)
			}
			if pe.Code != tc.code {
				t.Fatalf("%s answered the code %d (%s), want %d", tc.sql, pe.Code, pe.Message, tc.code)
			}
		})
	}
}
