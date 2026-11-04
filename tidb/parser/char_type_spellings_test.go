package parser

import "testing"

// MySQL spells CHAR and VARCHAR several ways and TiDB carries all of them:
// CHARACTER is CHAR, VARYING after either makes it VARCHAR, VARCHARACTER is
// VARCHAR in one word, and LONG takes the VARCHAR spellings as MEDIUMTEXT.
// CHAR SET is CHARACTER SET in two words.
// https://dev.mysql.com/doc/refman/8.0/en/char.html
//
// The rule that decides the rejections below is that every spelling of VARCHAR
// carries a length and every spelling of CHAR may leave it out. Verified on
// TiDB v8.5.5, which is also what the corpus in testdata records.

// charTypeSpellingsAccepted and charTypeSpellingsRejected are checked against
// the parser without a container below, and against a live TiDB in
// TestCharTypeSpellingsMatchTiDB, so the verdict recorded here is the engine's
// and not a memory of it.
var charTypeSpellingsAccepted = []string{
	`CREATE TABLE t (a CHARACTER)`,
	`CREATE TABLE t (a CHARACTER(10))`,
	`CREATE TABLE t (a CHAR VARYING(50))`,
	`CREATE TABLE t (a CHARACTER VARYING(50))`,
	`CREATE TABLE t (a VARCHARACTER(50))`,
	`CREATE TABLE t (a CHARACTER VARYING(0))`,
	`CREATE TABLE t (a CHARACTER VARYING(10) CHARACTER SET utf8)`,
	`CREATE TABLE t (a VARCHARACTER(50) CHARSET utf8)`,
	// NATIONAL and NCHAR take the same spellings.
	`CREATE TABLE t (a NATIONAL CHARACTER)`,
	`CREATE TABLE t (a NATIONAL CHARACTER VARYING(50))`,
	`CREATE TABLE t (a NATIONAL CHAR VARYING(50))`,
	`CREATE TABLE t (a NATIONAL VARCHARACTER(50))`,
	`CREATE TABLE t (a NCHAR VARCHAR(50))`,
	`CREATE TABLE t (a NCHAR VARYING(50))`,
	`CREATE TABLE t (a NCHAR VARCHARACTER(50))`,
	// LONG, with no length on any spelling.
	`CREATE TABLE t (a LONG VARCHARACTER)`,
	`CREATE TABLE t (a LONG CHAR VARYING)`,
	`CREATE TABLE t (a LONG CHARACTER VARYING)`,
	`CREATE TABLE t (a LONG CHARACTER VARYING CHARACTER SET utf8)`,
	// CHAR SET as the charset clause, after LONG and after a named type.
	`CREATE TABLE t (a LONG CHAR SET utf8mb4)`,
	`CREATE TABLE t (a CHAR CHAR SET utf8mb4)`,
	// The same spellings in a column change.
	`ALTER TABLE t MODIFY a CHARACTER VARYING(10)`,
	`ALTER TABLE t CHANGE COLUMN a a LONG CHAR VARYING BINARY CHARSET utf8 FIRST`,
}

var charTypeSpellingsRejected = []string{
	// A VARCHAR spelling with no length.
	`CREATE TABLE t (a VARCHAR)`,
	`CREATE TABLE t (a VARCHARACTER)`,
	`CREATE TABLE t (a CHAR VARYING)`,
	`CREATE TABLE t (a CHARACTER VARYING)`,
	`CREATE TABLE t (a NVARCHAR)`,
	`CREATE TABLE t (a NATIONAL VARCHAR)`,
	`CREATE TABLE t (a NCHAR VARYING)`,
	`CREATE TABLE t (a VARBINARY)`,
	// LONG takes no length with any spelling.
	`CREATE TABLE t (a LONG CHAR VARYING(10))`,
	`CREATE TABLE t (a LONG VARCHARACTER(10))`,
	// CHARACTER and CHAR only start a charset clause when SET follows, so
	// neither stands as a trailing modifier on its own.
	`CREATE TABLE t (a LONG CHARACTER)`,
	`CREATE TABLE t (a TEXT CHARACTER)`,
	// ...and neither stands in for the type: the column has none here.
	`CREATE TABLE t (a CHARACTER SET utf8)`,
	`CREATE TABLE t (a CHAR SET utf8)`,
}

func TestCharTypeSpellingsAccepted(t *testing.T) {
	for _, sql := range charTypeSpellingsAccepted {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

func TestCharTypeSpellingsRejected(t *testing.T) {
	for _, sql := range charTypeSpellingsRejected {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but TiDB v8.5.5 rejects it", sql)
			}
		})
	}
}

// TestCharTypeSpellingsMatchTiDB holds the two tables above to the engine
// itself. The corpus in testdata already covers the spellings it carries; this
// covers the rest, and above all the rejections, which are a rule about a
// missing length rather than a spelling TiDB has never heard of.
func TestCharTypeSpellingsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for want, cases := range map[bool][]string{
		true:  charTypeSpellingsAccepted,
		false: charTypeSpellingsRejected,
	} {
		for _, sql := range cases {
			t.Run(sql, func(t *testing.T) {
				accepts, err := tc.prepareSyntax(sql)
				if err != nil {
					t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
				}
				if accepts != want {
					t.Errorf("TiDB %ss %q, this file records %s",
						verdictName(accepts), sql, verdictName(want))
				}
			})
		}
	}
}
