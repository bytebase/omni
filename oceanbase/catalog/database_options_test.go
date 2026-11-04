package catalog

import "testing"

// The charset and collation a database option sets, for the names the parser
// began reading with this change: a string literal in either quote, a
// backquoted identifier, and BINARY. A form that parses and then leaves the
// schema's charset alone is worse than one that is rejected, because Bytebase
// would review later migrations against a schema that never existed.
func TestAlterDatabaseCharsetOptions(t *testing.T) {
	for _, tc := range []struct {
		name          string
		alter         string
		wantCharset   string
		wantCollation string
	}{
		{
			name:          "CHARACTER SET takes a string literal",
			alter:         "ALTER DATABASE testdb CHARACTER SET = 'latin1'",
			wantCharset:   "latin1",
			wantCollation: "latin1_swedish_ci",
		},
		{
			name:          "CHARSET takes a double-quoted literal",
			alter:         `ALTER DATABASE testdb CHARSET = "ascii"`,
			wantCharset:   "ascii",
			wantCollation: "ascii_general_ci",
		},
		{
			name:          "CHARACTER SET takes a backquoted identifier",
			alter:         "ALTER DATABASE testdb CHARACTER SET = `latin1`",
			wantCharset:   "latin1",
			wantCollation: "latin1_swedish_ci",
		},
		{
			name:          "COLLATE takes a string literal and leaves the charset",
			alter:         "ALTER DATABASE testdb COLLATE = 'utf8mb4_unicode_ci'",
			wantCharset:   "utf8mb4",
			wantCollation: "utf8mb4_unicode_ci",
		},
		{
			name:          "CHARACTER SET and COLLATE together",
			alter:         "ALTER DATABASE testdb CHARACTER SET = 'latin1' COLLATE = 'latin1_bin'",
			wantCharset:   "latin1",
			wantCollation: "latin1_bin",
		},
		{
			name:          "DEFAULT CHARSET and DEFAULT COLLATE",
			alter:         "ALTER SCHEMA testdb DEFAULT CHARSET = 'latin1' DEFAULT COLLATE = 'latin1_general_ci'",
			wantCharset:   "latin1",
			wantCollation: "latin1_general_ci",
		},
		{
			// The nameless arm applies to the database in use, which wtSetup
			// selects.
			name:          "no database name applies to the current one",
			alter:         "ALTER DATABASE COLLATE = 'utf8mb4_bin'",
			wantCharset:   "utf8mb4",
			wantCollation: "utf8mb4_bin",
		},
		{
			// OceanBase's collation name carries no BINARY arm, so the
			// parser refuses `COLLATE = binary`; quoted it is an ordinary
			// name, and CHARSET takes the bare word.
			name:          "CHARSET takes the reserved word BINARY",
			alter:         "ALTER DATABASE testdb CHARSET = binary",
			wantCharset:   "binary",
			wantCollation: "binary",
		},
		{
			name:          "COLLATE takes a quoted binary",
			alter:         "ALTER DATABASE testdb COLLATE = 'binary'",
			wantCharset:   "utf8mb4",
			wantCollation: "binary",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, tc.alter)

			db := c.GetDatabase("testdb")
			if db == nil {
				t.Fatal("database testdb is gone")
			}
			if db.Charset != tc.wantCharset {
				t.Errorf("charset = %q, want %q", db.Charset, tc.wantCharset)
			}
			if db.Collation != tc.wantCollation {
				t.Errorf("collation = %q, want %q", db.Collation, tc.wantCollation)
			}
		})
	}
}

// READ ONLY and READ WRITE are OceanBase's own database options, which the
// parser began reading with this change. The catalog holds no read-only flag,
// so what matters is that Exec takes them without an error and leaves the
// schema's charset and collation where they were: a statement that failed here
// would stop a migration omni is only reviewing.
func TestDatabaseReadOnlyOptionsAreNoOps(t *testing.T) {
	for _, sql := range []string{
		"ALTER DATABASE testdb READ ONLY",
		"ALTER DATABASE testdb READ WRITE",
		"ALTER DATABASE testdb CHARSET = 'latin1' READ ONLY",
		"CREATE DATABASE other READ ONLY",
		"CREATE DATABASE other READ WRITE",
	} {
		t.Run(sql, func(t *testing.T) {
			c := wtSetup(t)
			before := *c.GetDatabase("testdb")
			wtExec(t, c, sql)

			db := c.GetDatabase("testdb")
			if db == nil {
				t.Fatal("database testdb is gone")
			}
			wantCharset, wantCollation := before.Charset, before.Collation
			if sql == "ALTER DATABASE testdb CHARSET = 'latin1' READ ONLY" {
				wantCharset, wantCollation = "latin1", "latin1_swedish_ci"
			}
			if db.Charset != wantCharset {
				t.Errorf("charset = %q, want %q", db.Charset, wantCharset)
			}
			if db.Collation != wantCollation {
				t.Errorf("collation = %q, want %q", db.Collation, wantCollation)
			}
		})
	}
}
