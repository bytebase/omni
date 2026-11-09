package parser

import "testing"

// The option list CREATE DATABASE and ALTER DATABASE read, and the names its
// CHARACTER SET, CHARSET and COLLATE take. The name is MySQL's ident_or_text
// plus the reserved word BINARY -- an identifier, a non-reserved keyword, a
// backquoted identifier or a string literal in either quote -- so an integer
// and every other reserved keyword are a 1064.
//
// The rejections are the point of this file. The corpus in tidb/parser pins the
// forms pingcap's own parser test writes, and it writes none of these: a bare
// COLLATE, a conflicting charset, a comma between options, an option list left
// empty, a DEFAULT with nothing after it, and the options each engine does not
// have. Every case here was measured against TiDB v8.5.5, and the expectation
// is what that server answered, so an image upgrade that changes a verdict
// fails here instead of passing silently.

// databaseOptionCase is a statement and whether the engine's parser accepts it.
type databaseOptionCase struct {
	sql    string
	accept bool
}

var databaseOptionCases = []databaseOptionCase{
	// The charset and collation name in every shape the engines read.
	{"ALTER DATABASE d CHARACTER SET = 'utf8'", true},
	{"ALTER DATABASE d CHARACTER SET 'utf8'", true},
	{`ALTER DATABASE d CHARACTER SET = "utf8"`, true},
	{"ALTER DATABASE d CHARACTER SET = `utf8`", true},
	{"ALTER DATABASE d CHARACTER SET utf8", true},
	{"ALTER DATABASE d DEFAULT CHARACTER SET = 'utf8'", true},
	{"ALTER SCHEMA d DEFAULT CHARSET = 'UTF8'", true},
	{"ALTER DATABASE d CHARSET = binary", true},
	{"ALTER DATABASE d CHARSET binary", true},
	{"ALTER DATABASE d CHARACTER SET binary", true},
	{"ALTER DATABASE d COLLATE = 'utf8_bin'", true},
	{"ALTER DATABASE d COLLATE 'utf8_bin'", true},
	{"ALTER DATABASE d COLLATE = `utf8_bin`", true},
	{"ALTER DATABASE d DEFAULT COLLATE = 'utf8_bin'", true},
	// A name the server does not know is a 1115 when the statement runs, not a
	// 1064 when it parses.
	{"ALTER DATABASE d CHARSET = 'nosuchcharset'", true},
	// Several options at once, with and without DEFAULT, and conflicting
	// charsets, which the server settles after parsing rather than refusing.
	{"ALTER DATABASE d CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'", true},
	{"ALTER DATABASE d CHARACTER SET = 'utf8' CHARACTER SET = 'utf8mb4'", true},
	{"ALTER DATABASE d DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8' COLLATE = 'utf8mb4_bin'", true},
	// The nameless arm, which the reserved words COLLATE and DEFAULT reach.
	{"ALTER DATABASE COLLATE = 'utf8_bin'", true},
	{"ALTER SCHEMA COLLATE = 'utf8_bin'", true},
	{"ALTER DATABASE DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'", true},
	// An empty quoted name is a name: the server parses it and answers 1102.
	{"ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'", true},
	// CREATE DATABASE reads the same names and needs no option at all.
	{"CREATE DATABASE d", true},
	{"CREATE DATABASE d CHARACTER SET = 'utf8'", true},
	{"CREATE DATABASE d COLLATE = 'utf8_bin'", true},
	{"CREATE DATABASE d CHARSET = binary", true},
	{"CREATE SCHEMA d DEFAULT COLLATE = 'UTF8_BiN'", true},
	{"CREATE DATABASE d DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8mb4_bin'", true},

	// BINARY is a collation name as well as a charset name.
	{"ALTER DATABASE d COLLATE = binary", true},
	{"ALTER DATABASE d COLLATE binary", true},
	{"ALTER DATABASE d CHARSET=binary COLLATE = binary", true},

	// ENCRYPTION, whose value is a quoted Y or N and is required.
	{"ALTER DATABASE d ENCRYPTION = 'N'", true},
	{"ALTER DATABASE d ENCRYPTION 'N'", true},
	{"ALTER DATABASE d DEFAULT ENCRYPTION = 'N'", true},
	{"CREATE DATABASE d ENCRYPTION = 'N'", true},
	{"CREATE DATABASE d DEFAULT ENCRYPTION = 'Y'", true},
	{"ALTER DATABASE d ENCRYPTION", false},
	{"ALTER DATABASE d ENCRYPTION = N", false},
	// ENCRYPTION is not reserved either, so the nameless arm does not reach it.
	{"ALTER DATABASE ENCRYPTION = 'N'", false},

	// TiDB v8.5.5 has no READ ONLY option, in any spelling, on either
	// statement.
	{"ALTER DATABASE d READ ONLY = 1", false},
	{"ALTER DATABASE d READ ONLY = 0", false},
	{"ALTER DATABASE d READ ONLY = DEFAULT", false},
	{"ALTER DATABASE d READ ONLY 1", false},
	{"ALTER DATABASE d READ ONLY", false},
	{"ALTER DATABASE READ ONLY = 1", false},
	{"ALTER DATABASE d READ", false},
	{"ALTER DATABASE d READ WRITE", false},
	{"ALTER DATABASE d DEFAULT READ ONLY = 1", false},
	{"ALTER DATABASE d CHARSET = utf8 READ ONLY", false},

	// The name is required wherever an option takes one.
	{"ALTER DATABASE d COLLATE", false},
	{"ALTER DATABASE d COLLATE =", false},
	{"ALTER DATABASE d CHARSET =", false},
	{"ALTER DATABASE d CHARACTER SET", false},
	{"ALTER DATABASE d CHARACTER SET =", false},
	{"ALTER DATABASE d CHARACTER SET = 'utf8' COLLATE", false},
	// An integer is no name, and neither is a reserved keyword other than
	// BINARY.
	{"ALTER DATABASE d CHARSET = 1", false},
	{"ALTER DATABASE d COLLATE = 1", false},
	{"ALTER DATABASE d CHARSET = select", false},
	{"ALTER DATABASE d COLLATE = select", false},
	{"CREATE DATABASE d CHARSET = 1", false},
	{"CREATE DATABASE d COLLATE = select", false},
	// The option list carries no commas, and no comma stands after the name.
	{"ALTER DATABASE d CHARSET = 'utf8', COLLATE = 'utf8_bin'", false},
	{"ALTER DATABASE d, CHARSET = 'utf8'", false},
	{"CREATE DATABASE d CHARSET = 'utf8', COLLATE = 'utf8_bin'", false},
	// A completed option is not followed by a bare token: no DatabaseOption
	// begins with an identifier, so a trailing word after a valid collation
	// name is a 1064. omni and the engine agree here.
	//
	// The one exception is a string naming a collation the server does not know:
	// `ALTER DATABASE COLLATE 's' zzz` is accepted by TiDB v8.5.5 (and recorded
	// as a gap at parser_y_negatives_corpus.tsv CollationName#1/insert). The
	// CollationName action calls GetCollationByName, and on an unknown name it
	// returns 1, which fires the grammar's error recovery and swallows the
	// trailing tokens; a known name parses strictly and rejects them. omni is a
	// syntactic parser with no collation catalog, so it cannot accept the
	// unknown-name form without also accepting the known-name form the engine
	// rejects -- that is why the string case is not pinned here as agreement.
	{"ALTER DATABASE COLLATE binary zzz", false},
	// ALTER DATABASE needs one option; CREATE DATABASE is the one that does not.
	{"ALTER DATABASE", false},
	{"ALTER DATABASE d", false},
	{"ALTER SCHEMA", false},
	{"ALTER SCHEMA d", false},
	// DEFAULT is a prefix, not an option, in front of one and after one.
	{"ALTER DATABASE d DEFAULT", false},
	{"ALTER DATABASE d DEFAULT = 'utf8'", false},
	{"ALTER DATABASE d CHARSET = utf8 DEFAULT", false},
	// CHARACTER starts an option only when SET follows it.
	{"ALTER DATABASE d CHARSET = utf8 CHARACTER", false},
	{"CREATE DATABASE d CHARSET = utf8 CHARACTER", false},
	// A reserved word is no database name.
	{"ALTER DATABASE binary CHARSET = 'utf8'", false},
	// READ ONLY is no CREATE DATABASE option on any of these engines.
	{"CREATE DATABASE d READ ONLY = 1", false},
	{"CREATE DATABASE d READ ONLY 0", false},
	{"CREATE DATABASE d READ ONLY = DEFAULT", false},

	// CHARSET is not reserved, so it reads as the database name and the
	// statement then runs out of options.
	{"ALTER DATABASE CHARSET = 'utf8'", false},
	{"ALTER DATABASE CHARSET utf8", false},

	// TiDB's own two options. PLACEMENT POLICY SET DEFAULT is an arm of its
	// own, where SET is not the policy name, and the DEFAULT prefix reaches
	// PLACEMENT POLICY but not SET TIFLASH REPLICA. PLACEMENT is not reserved
	// and yet the nameless arm reaches it, which no other non-reserved keyword
	// here does.
	{"ALTER DATABASE d PLACEMENT POLICY SET DEFAULT", true},
	{"ALTER DATABASE d CHARSET = utf8 PLACEMENT POLICY SET DEFAULT", true},
	{"ALTER DATABASE d DEFAULT PLACEMENT POLICY = 'p'", true},
	{"ALTER DATABASE PLACEMENT POLICY = 'p'", true},
	{"ALTER DATABASE d SET TIFLASH REPLICA 1", true},
	{"ALTER DATABASE SET TIFLASH REPLICA 1", true},
	{"ALTER DATABASE d PLACEMENT POLICY SET", false},
	{"ALTER DATABASE d PLACEMENT POLICY SET x", false},
	{"ALTER DATABASE d PLACEMENT", false},
	{"ALTER DATABASE d DEFAULT SET TIFLASH REPLICA 1", false},
}

// TestDatabaseOptionsOracle asserts omni agrees with the pinned TiDB on
// databaseOptionCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it: 1064 is the parser's rejection,
// and any other server error means the text parsed and failed later.
func TestDatabaseOptionsOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for i, dbc := range databaseOptionCases {
		t.Run(dbc.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = databaseOptionCases[i-1].sql
			}
			accepts, _, err := tc.syntaxVerdict(t, dbc.sql, previous)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != dbc.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbCorpusVersion, accepts, dbc.accept, dbc.sql)
			}
			_, perr := Parse(dbc.sql)
			if omniAccepts := perr == nil; omniAccepts != dbc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbCorpusVersion, dbc.accept, dbc.sql, perr)
			}
		})
	}
}
