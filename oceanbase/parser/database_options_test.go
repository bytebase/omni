package parser

import "testing"

// The option list CREATE DATABASE and ALTER DATABASE read, and the names its
// CHARACTER SET, CHARSET and COLLATE take. This is where OceanBase departs
// from the MySQL grammar this tree forked, in four places:
//
//   - BINARY names a charset but not a collation, so `CHARSET = binary` parses
//     and `COLLATE = binary` is a 1064, while the same word quoted -- in
//     single quotes or in backquotes -- is an ordinary ident_or_text name and
//     parses.
//   - CHARSET is not reserved, yet `ALTER DATABASE CHARSET = 'utf8'` reaches
//     the nameless arm rather than naming a database called CHARSET. MySQL
//     8.0, MariaDB 11.8.8 and TiDB v8.5.5 all answer 1064 there.
//   - ENCRYPTION is not an option, on either statement.
//   - READ ONLY takes no value, where MySQL's reads {DEFAULT | 0 | 1}, and it
//     reads on CREATE DATABASE as well. READ WRITE is its other half, which
//     MySQL has on neither statement.
//
// The engine column is what 5.7.25-OceanBase_CE-v4.3.5.6 answered on
// 2026-11-04, read with scripts/oceanbase-oracle.sh. OceanBase needs about
// 8 GB and CI has no room for it, so TestDatabaseOptionsOceanBaseOracle reads
// that column back from a live server only where the DSN is set; this test
// holds the parser to the record unconditionally, which is the gate CI runs.
type databaseOptionCase struct {
	sql string
	// engine is whether OceanBase CE 4.3.5.6 parses the statement.
	engine bool
	// pending marks a statement OceanBase accepts that this parser does not
	// read yet. Implementing it means clearing the flag in the same change.
	pending bool
}

var databaseOptionCases = []databaseOptionCase{
	// The charset and collation name in every shape the engine reads.
	{sql: "ALTER DATABASE obdbopt CHARACTER SET = 'utf8'", engine: true},
	{sql: "ALTER DATABASE obdbopt CHARACTER SET 'utf8'", engine: true},
	{sql: `ALTER DATABASE obdbopt CHARACTER SET = "utf8"`, engine: true},
	{sql: "ALTER DATABASE obdbopt CHARACTER SET = `utf8`", engine: true},
	{sql: "ALTER DATABASE obdbopt CHARACTER SET utf8", engine: true},
	{sql: "ALTER DATABASE obdbopt DEFAULT CHARACTER SET = 'utf8'", engine: true},
	{sql: "ALTER SCHEMA obdbopt DEFAULT CHARSET = 'UTF8'", engine: true},
	{sql: "ALTER DATABASE obdbopt COLLATE = 'utf8_bin'", engine: true},
	{sql: "ALTER DATABASE obdbopt COLLATE 'utf8_bin'", engine: true},
	{sql: "ALTER DATABASE obdbopt DEFAULT COLLATE = 'utf8_bin'", engine: true},
	// A name the server does not know is a 1115, not a syntax error.
	{sql: "ALTER DATABASE obdbopt CHARSET = 'nosuchcharset'", engine: true},
	// Several options at once, and a conflicting charset, which the engine
	// settles after parsing.
	{sql: "ALTER DATABASE obdbopt CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'", engine: true},
	{sql: "ALTER DATABASE obdbopt CHARACTER SET = 'utf8' CHARACTER SET = 'utf8mb4'", engine: true},
	{sql: "ALTER DATABASE obdbopt DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8' COLLATE = 'utf8mb4_bin'", engine: true},
	// The nameless arm, which CHARSET reaches here and nowhere else.
	{sql: "ALTER DATABASE CHARACTER SET = 'utf8'", engine: true},
	{sql: "ALTER DATABASE COLLATE = 'utf8_bin'", engine: true},
	{sql: "ALTER SCHEMA COLLATE = 'utf8_bin'", engine: true},
	{sql: "ALTER DATABASE CHARSET = 'utf8'", engine: true},
	{sql: "ALTER DATABASE CHARSET utf8", engine: true},
	{sql: "ALTER DATABASE DEFAULT CHARSET = 'utf8'", engine: true},
	// An empty quoted name is a name: the engine parses it and answers 1102.
	{sql: "ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'", engine: true},

	// BINARY is a charset name and not a collation name.
	{sql: "ALTER DATABASE obdbopt CHARSET = binary", engine: true},
	{sql: "ALTER DATABASE obdbopt CHARSET binary", engine: true},
	{sql: "ALTER DATABASE obdbopt CHARACTER SET binary", engine: true},
	{sql: "ALTER DATABASE obdbopt COLLATE = `binary`", engine: true},
	{sql: "ALTER DATABASE obdbopt COLLATE = 'binary'", engine: true},
	{sql: "ALTER DATABASE obdbopt COLLATE = binary"},
	{sql: "ALTER DATABASE obdbopt COLLATE binary"},
	{sql: "ALTER DATABASE obdbopt CHARSET=binary COLLATE = binary"},

	// READ ONLY and READ WRITE, which take no value, read on CREATE DATABASE
	// as well, and read beside a charset option on either side.
	{sql: "ALTER DATABASE obdbopt READ ONLY", engine: true},
	{sql: "ALTER DATABASE obdbopt CHARSET = utf8 READ ONLY", engine: true},
	{sql: "ALTER DATABASE obdbopt READ ONLY CHARSET = utf8", engine: true},
	{sql: "ALTER DATABASE obdbopt READ WRITE", engine: true},
	{sql: "ALTER DATABASE obdbopt READ ONLY = 1"},
	{sql: "ALTER DATABASE obdbopt READ ONLY = 0"},
	{sql: "ALTER DATABASE obdbopt READ ONLY = DEFAULT"},
	{sql: "ALTER DATABASE obdbopt READ ONLY 1"},
	{sql: "ALTER DATABASE READ ONLY = 1"},
	{sql: "ALTER DATABASE obdbopt READ"},
	{sql: "ALTER DATABASE obdbopt DEFAULT READ ONLY = 1"},
	{sql: "ALTER DATABASE obdbopt READ ONLY, CHARSET = utf8"},

	// ENCRYPTION is not an option here, on either statement.
	{sql: "ALTER DATABASE obdbopt ENCRYPTION = 'N'"},
	{sql: "ALTER DATABASE obdbopt ENCRYPTION 'N'"},
	{sql: "ALTER DATABASE obdbopt DEFAULT ENCRYPTION = 'N'"},
	{sql: "ALTER DATABASE obdbopt ENCRYPTION"},
	{sql: "ALTER DATABASE ENCRYPTION = 'N'"},

	// The name is required wherever an option takes one.
	{sql: "ALTER DATABASE obdbopt COLLATE"},
	{sql: "ALTER DATABASE obdbopt COLLATE ="},
	{sql: "ALTER DATABASE obdbopt CHARSET ="},
	{sql: "ALTER DATABASE obdbopt CHARACTER SET"},
	{sql: "ALTER DATABASE obdbopt CHARACTER SET ="},
	{sql: "ALTER DATABASE obdbopt CHARACTER SET = 'utf8' COLLATE"},
	// An integer is no name, and neither is a reserved keyword other than BINARY.
	{sql: "ALTER DATABASE obdbopt CHARSET = 1"},
	{sql: "ALTER DATABASE obdbopt COLLATE = 1"},
	{sql: "ALTER DATABASE obdbopt CHARSET = select"},
	{sql: "ALTER DATABASE obdbopt COLLATE = select"},
	// DEFAULT is no value either, which is where MariaDB differs.
	{sql: "ALTER DATABASE obdbopt CHARACTER SET = DEFAULT"},
	{sql: "ALTER DATABASE obdbopt COLLATE = DEFAULT"},
	// The option list carries no commas, and no comma stands after the name.
	{sql: "ALTER DATABASE obdbopt CHARSET = 'utf8', COLLATE = 'utf8_bin'"},
	{sql: "ALTER DATABASE obdbopt, CHARSET = 'utf8'"},
	// ALTER DATABASE needs one option; CREATE DATABASE is the one that does not.
	{sql: "ALTER DATABASE"},
	{sql: "ALTER DATABASE obdbopt"},
	{sql: "ALTER SCHEMA"},
	{sql: "ALTER SCHEMA obdbopt"},
	// DEFAULT is a prefix, not an option, in front of one and after one.
	{sql: "ALTER DATABASE obdbopt DEFAULT"},
	{sql: "ALTER DATABASE obdbopt DEFAULT = 'utf8'"},
	{sql: "ALTER DATABASE obdbopt CHARSET = utf8 DEFAULT"},
	// CHARACTER starts an option only when SET follows it.
	{sql: "ALTER DATABASE obdbopt CHARSET = utf8 CHARACTER"},
	// A reserved word is no database name.
	{sql: "ALTER DATABASE binary CHARSET = 'utf8'"},
	// MariaDB's COMMENT and UPGRADE DATA DIRECTORY NAME, and TiDB's placement
	// options, are not this engine's.
	{sql: "ALTER DATABASE obdbopt COMMENT = 'x'"},
	{sql: "ALTER DATABASE obdbopt COMMENT 'x'"},
	{sql: "ALTER DATABASE COMMENT = 'x'"},
	{sql: "ALTER DATABASE obdbopt UPGRADE DATA DIRECTORY NAME"},
	{sql: "ALTER DATABASE PLACEMENT POLICY = 'p'"},

	// CREATE DATABASE reads the same names and needs no option at all.
	{sql: "CREATE DATABASE obdbopt_c", engine: true},
	{sql: "CREATE DATABASE obdbopt_c CHARACTER SET = 'utf8'", engine: true},
	{sql: "CREATE DATABASE obdbopt_c COLLATE = 'utf8_bin'", engine: true},
	{sql: "CREATE DATABASE obdbopt_c CHARSET = binary", engine: true},
	{sql: "CREATE SCHEMA obdbopt_c DEFAULT COLLATE = 'UTF8_BiN'", engine: true},
	{sql: "CREATE DATABASE obdbopt_c DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8mb4_bin'", engine: true},
	{sql: "CREATE DATABASE obdbopt_c CHARSET = 1"},
	{sql: "CREATE DATABASE obdbopt_c COLLATE = select"},
	{sql: "CREATE DATABASE obdbopt_c CHARSET = 'utf8', COLLATE = 'utf8_bin'"},
	{sql: "CREATE DATABASE obdbopt_c CHARSET = utf8 CHARACTER"},
	{sql: "CREATE DATABASE obdbopt_c READ ONLY", engine: true},
	{sql: "CREATE DATABASE obdbopt_c READ WRITE", engine: true},
	{sql: "CREATE DATABASE obdbopt_c READ ONLY = 1"},
	{sql: "CREATE DATABASE obdbopt_c READ ONLY 0"},
	{sql: "CREATE DATABASE obdbopt_c READ ONLY = DEFAULT"},
	{sql: "CREATE DATABASE obdbopt_c ENCRYPTION = 'N'"},
	{sql: "CREATE DATABASE obdbopt_c DEFAULT ENCRYPTION = 'Y'"},
	{sql: "CREATE DATABASE obdbopt_c COMMENT = 'hello'"},
}

// TestDatabaseOptions holds the parser to the recorded OceanBase verdicts. It
// needs no server, so it is the gate CI runs: the parser must read exactly what
// CE 4.3.5.6 parses, minus the pending rows.
func TestDatabaseOptions(t *testing.T) {
	for _, c := range databaseOptionCases {
		t.Run(c.sql, func(t *testing.T) {
			want := c.engine && !c.pending
			_, err := Parse(c.sql)
			if got := err == nil; got != want {
				t.Errorf("parser accepts=%v, want %v for %q (OceanBase accepts=%v, pending=%v): %v",
					got, want, c.sql, c.engine, c.pending, err)
			}
		})
	}
}
