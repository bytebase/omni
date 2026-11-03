package parser

import (
	"testing"

	ast "github.com/bytebase/omni/mariadb/ast"
)

// indexIgnorabilityCases are the index visibility forms, each with MariaDB
// 11.8.8's verdict. MariaDB hides an index from the optimizer with IGNORED and
// shows it again with NOT IGNORED, on any index definition and in
// ALTER {INDEX | KEY} [IF EXISTS] name; MySQL's INVISIBLE is a syntax error,
// and VISIBLE is taken on a definition as a no-op but not by ALTER INDEX.
// A comment versioned 50700 to 99999 is skipped unread, so MySQL's
// /*!80000 INVISIBLE */ stands. accept is the engine's verdict, held by
// TestIndexIgnorabilityOracle; 4174 (a primary key cannot be ignored) and a
// missing table are not syntax errors and count as accepted.
var indexIgnorabilityCases = []struct {
	sql    string
	accept bool
}{
	// IGNORED and NOT IGNORED on every kind of index definition.
	{"CREATE TABLE t (a INT, KEY i (a) IGNORED)", true},
	{"CREATE TABLE t (a INT, KEY i (a) NOT IGNORED)", true},
	{"CREATE TABLE t (a INT, INDEX (a) IGNORED)", true},
	{"CREATE TABLE t (a INT, UNIQUE KEY u (a) COMMENT 'x' IGNORED)", true},
	{"CREATE TABLE t (a INT, KEY i (a) IGNORED COMMENT 'x')", true},
	{"CREATE TABLE t (a INT, KEY i USING HASH (a) IGNORED KEY_BLOCK_SIZE=8)", true},
	{"CREATE TABLE t (a INT, KEY i (a) USING BTREE IGNORED)", true},
	{"CREATE TABLE t (a INT, b TEXT, FULLTEXT KEY f (b) IGNORED)", true},
	{"CREATE TABLE t (g GEOMETRY NOT NULL, SPATIAL KEY s (g) IGNORED)", true},
	{"CREATE TABLE t (a INT, CONSTRAINT c UNIQUE (a) IGNORED)", true},
	{"CREATE TABLE t (a INT, KEY i (a) IGNORED NOT IGNORED)", true},
	{"CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) NOT IGNORED)", true},
	{"CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) IGNORED)", true}, // 4174
	{"CREATE TABLE t (ignored INT, KEY ignored (ignored) IGNORED)", true},
	{"CREATE INDEX i ON t (a) IGNORED", true},
	{"CREATE INDEX i ON t (a) NOT IGNORED ALGORITHM=INPLACE", true},
	{"CREATE UNIQUE INDEX i ON t (a) COMMENT 'u' IGNORED LOCK=NONE", true},
	{"ALTER TABLE t ADD INDEX i (a) IGNORED", true},
	{"ALTER TABLE t ADD CONSTRAINT u UNIQUE (a) NOT IGNORED", true},
	{"ALTER TABLE t ADD PRIMARY KEY (a) IGNORED", true}, // 4174 once t exists
	// VISIBLE on a definition, MariaDB's no-op for MySQL 8.0 compatibility.
	{"CREATE TABLE t (a INT, KEY i (a) VISIBLE)", true},
	{"CREATE TABLE t (a INT, KEY i (a) IGNORED VISIBLE)", true},
	{"CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) VISIBLE)", true},
	{"CREATE INDEX i ON t (a) VISIBLE", true},
	{"ALTER TABLE t ADD INDEX i (a) VISIBLE", true},
	// ALTER {INDEX | KEY} [IF EXISTS] name [NOT] IGNORED.
	{"ALTER TABLE t ALTER INDEX i IGNORED", true},
	{"ALTER TABLE t ALTER INDEX i NOT IGNORED", true},
	{"ALTER TABLE t ALTER KEY i IGNORED", true},
	{"ALTER TABLE t ALTER INDEX IF EXISTS i NOT IGNORED", true},
	{"ALTER TABLE t ALTER KEY IF EXISTS i IGNORED", true},
	{"ALTER TABLE t ALTER INDEX `PRIMARY` IGNORED", true},
	{"ALTER TABLE t ALTER INDEX ignored NOT IGNORED", true},
	{"ALTER TABLE t ALTER INDEX i IGNORED, ADD INDEX j (a) IGNORED", true},
	// MySQL's comment-wrapped INVISIBLE, which MariaDB skips by its version.
	{"CREATE TABLE t (a INT, KEY i (a) /*!80000 INVISIBLE */)", true},
	{"CREATE TABLE t (a INT, KEY i (a) /*!50700 INVISIBLE */)", true},
	{"CREATE TABLE t (a INT, KEY i (a) /*!99999 INVISIBLE */)", true},
	{"CREATE INDEX i ON t (a) /*!80000 INVISIBLE */", true},
	{"ALTER TABLE t ADD INDEX i (a) /*!80000 INVISIBLE */", true},
	// Column-level INVISIBLE is MariaDB's invisible column and stays.
	{"CREATE TABLE t (a INT INVISIBLE, b INT)", true},
	{"ALTER TABLE t ADD COLUMN e INT INVISIBLE", true},

	// MySQL's INVISIBLE, read by MariaDB, is a syntax error on any index.
	{"CREATE TABLE t (a INT, KEY i (a) INVISIBLE)", false},
	{"CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) INVISIBLE)", false},
	{"CREATE INDEX i ON t (a) INVISIBLE", false},
	{"ALTER TABLE t ADD INDEX i (a) INVISIBLE", false},
	{"CREATE TABLE t (a INT, KEY i (a) /*!INVISIBLE */)", false},
	{"CREATE TABLE t (a INT, KEY i (a) /*!50699 INVISIBLE */)", false},
	{"CREATE TABLE t (a INT, KEY i (a) /*!100000 INVISIBLE */)", false},
	// ALTER INDEX takes only the ignorability, and requires it.
	{"ALTER TABLE t ALTER INDEX i INVISIBLE", false},
	{"ALTER TABLE t ALTER INDEX i VISIBLE", false},
	{"ALTER TABLE t ALTER INDEX i", false},
	{"ALTER TABLE t ALTER INDEX i NOT", false},
	{"ALTER TABLE t ALTER INDEX i IGNORED IGNORED", false},
	{"ALTER TABLE t ALTER INDEX i `IGNORED`", false},
	{"ALTER TABLE t ALTER INDEX i /*!80000 INVISIBLE */", false},
	{"ALTER TABLE t ALTER INDEX primary IGNORED", false},
	{"ALTER TABLE t ALTER COLUMN a IGNORED", false},
	// IGNORED is a word only unquoted, NOT only before it, and not a column attribute.
	{"CREATE TABLE t (a INT, KEY i (a) `IGNORED`)", false},
	{"CREATE TABLE t (a INT, KEY i (a) NOT)", false},
	{"CREATE TABLE t (a INT, KEY i (a) NOT VISIBLE)", false},
	{"CREATE TABLE t (a INT, KEY i (a) IGNORE)", false},
	{"CREATE TABLE t (a INT, KEY i (a) COMMENT '/*!80000' INVISIBLE)", false},
	{"CREATE TABLE t (a INT UNIQUE IGNORED)", false},
	{"CREATE TABLE t (a INT PRIMARY KEY IGNORED)", false},
}

// TestIndexIgnorability holds the parser to MariaDB's verdict on each form.
func TestIndexIgnorability(t *testing.T) {
	for _, tc := range indexIgnorabilityCases {
		_, err := Parse(tc.sql)
		if got := err == nil; got != tc.accept {
			t.Errorf("%s: parser accepts = %v, MariaDB 11.8.8 accepts = %v (err: %v)", tc.sql, got, tc.accept, err)
		}
	}
}

// TestIndexIgnorabilityOracle asks MariaDB 11.8.8 about every case (container
// test). Only 1064 is a syntax error.
func TestIndexIgnorabilityOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, tc := range indexIgnorabilityCases {
		v, code, msg := o.classify(tc.sql)
		if v == vErrored {
			t.Fatalf("%s: MariaDB gave no verdict: %s", tc.sql, msg)
		}
		if got := v == vAccepted; got != tc.accept {
			t.Errorf("%s: MariaDB accepts = %v (%d %s), case says %v", tc.sql, got, code, msg, tc.accept)
		}
	}
}

// TestIndexIgnorabilityAST pins what each form leaves in the tree: the option
// name a definition carries, and the command type, name and IF EXISTS of an
// ALTER INDEX.
func TestIndexIgnorabilityAST(t *testing.T) {
	options := func(t *testing.T, sql string) []string {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
		var opts []*ast.IndexOption
		switch n := list.Items[0].(type) {
		case *ast.CreateTableStmt:
			opts = n.Constraints[0].IndexOptions
		case *ast.CreateIndexStmt:
			opts = n.Options
		case *ast.AlterTableStmt:
			opts = n.Commands[0].Constraint.IndexOptions
		default:
			t.Fatalf("%s: unexpected %T", sql, n)
		}
		names := make([]string, len(opts))
		for i, o := range opts {
			names[i] = o.Name
		}
		return names
	}
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"CREATE TABLE t (a INT, KEY i (a) COMMENT 'x' IGNORED)", []string{"COMMENT", "IGNORED"}},
		{"CREATE TABLE t (a INT, KEY i (a) NOT IGNORED VISIBLE)", []string{"NOT IGNORED", "VISIBLE"}},
		{"CREATE INDEX i ON t (a) IGNORED", []string{"IGNORED"}},
		{"ALTER TABLE t ADD UNIQUE u (a) NOT IGNORED", []string{"NOT IGNORED"}},
		// MariaDB skips /*!80000 ... */ unread, so the index has no option.
		{"CREATE TABLE t (a INT, KEY i (a) /*!80000 INVISIBLE */)", nil},
	} {
		got := options(t, tc.sql)
		if len(got) != len(tc.want) {
			t.Errorf("%s: options %v, want %v", tc.sql, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("%s: options %v, want %v", tc.sql, got, tc.want)
				break
			}
		}
	}

	for _, tc := range []struct {
		sql      string
		typ      ast.AlterTableCmdType
		name     string
		ifExists bool
	}{
		{"ALTER TABLE t ALTER INDEX i IGNORED", ast.ATAlterIndexInvisible, "i", false},
		{"ALTER TABLE t ALTER KEY i NOT IGNORED", ast.ATAlterIndexVisible, "i", false},
		{"ALTER TABLE t ALTER INDEX IF EXISTS `PRIMARY` IGNORED", ast.ATAlterIndexInvisible, "PRIMARY", true},
	} {
		stmt := parseAlterTable(t, tc.sql)
		cmd := stmt.Commands[0]
		if cmd.Type != tc.typ || cmd.Name != tc.name || cmd.IfExists != tc.ifExists {
			t.Errorf("%s: type %d name %q if exists %v, want %d %q %v", tc.sql, cmd.Type, cmd.Name, cmd.IfExists, tc.typ, tc.name, tc.ifExists)
		}
	}
}
