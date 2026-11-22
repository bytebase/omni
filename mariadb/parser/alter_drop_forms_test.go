package parser

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestAlterTableDropForms covers RESTRICT or CASCADE after a column, which
// MariaDB's alter_list_item allows on DROP [COLUMN] only (opt_restrict, a
// no-op for ALTER TABLE), and KEY after PRIMARY and FOREIGN. The command spans
// its last token, so a trailing RESTRICT or CASCADE is inside its Loc.
// Container-verified by TestMariaDBAlterDropFormsOracle.
func TestAlterTableDropForms(t *testing.T) {
	cases := []struct {
		sql      string
		typ      ast.AlterTableCmdType
		name     string
		ifExists bool
	}{
		{"ALTER TABLE t DROP COLUMN m RESTRICT", ast.ATDropColumn, "m", false},
		{"ALTER TABLE t DROP COLUMN m CASCADE", ast.ATDropColumn, "m", false},
		{"ALTER TABLE t DROP m RESTRICT", ast.ATDropColumn, "m", false},
		{"ALTER TABLE t DROP k CASCADE", ast.ATDropColumn, "k", false},
		{"ALTER TABLE t DROP COLUMN m /* c */ restrict", ast.ATDropColumn, "m", false},
		{"ALTER TABLE t DROP COLUMN `restrict`", ast.ATDropColumn, "restrict", false},
		{"ALTER TABLE t DROP COLUMN IF EXISTS m RESTRICT", ast.ATDropColumn, "m", true},
		{"ALTER TABLE t DROP IF EXISTS m CASCADE", ast.ATDropColumn, "m", true},
		{"ALTER TABLE t DROP PRIMARY KEY", ast.ATDropConstraint, "PRIMARY", false},
		{"ALTER TABLE t DROP FOREIGN KEY f", ast.ATDropConstraint, "f", false},
		{"ALTER TABLE t DROP FOREIGN KEY IF EXISTS f", ast.ATDropConstraint, "f", true},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseSeqStmt[*ast.AlterTableStmt](t, tc.sql)
			if len(stmt.Commands) != 1 {
				t.Fatalf("got %d commands, want 1", len(stmt.Commands))
			}
			cmd := stmt.Commands[0]
			if cmd.Type != tc.typ || cmd.Name != tc.name || cmd.IfExists != tc.ifExists {
				t.Errorf("got Type=%v Name=%q IfExists=%v, want Type=%v Name=%q IfExists=%v",
					cmd.Type, cmd.Name, cmd.IfExists, tc.typ, tc.name, tc.ifExists)
			}
			wantLoc := ast.Loc{Start: strings.Index(tc.sql, "DROP"), End: len(tc.sql)}
			if cmd.Loc != wantLoc {
				t.Errorf("cmd Loc = %+v, want %+v", cmd.Loc, wantLoc)
			}
		})
	}

	t.Run("several commands", func(t *testing.T) {
		sql := "ALTER TABLE t DROP COLUMN m RESTRICT, DROP k CASCADE, ADD COLUMN z INT"
		stmt := parseSeqStmt[*ast.AlterTableStmt](t, sql)
		if len(stmt.Commands) != 3 {
			t.Fatalf("got %d commands, want 3", len(stmt.Commands))
		}
		for i, want := range []string{"DROP COLUMN m RESTRICT", "DROP k CASCADE", "ADD COLUMN z INT"} {
			start := strings.Index(sql, want)
			wantLoc := ast.Loc{Start: start, End: start + len(want)}
			if got := stmt.Commands[i].Loc; got != wantLoc {
				t.Errorf("command %d Loc = %+v, want %+v (%q)", i, got, wantLoc, want)
			}
		}
	})
}

// TestAlterTableDropFormsReject pins the DROP forms MariaDB 11.8 rejects with
// 1064, container-verified by TestMariaDBAlterDropFormsOracle. Only a column
// takes RESTRICT or CASCADE.
func TestAlterTableDropFormsReject(t *testing.T) {
	cases := []struct {
		sql string
		at  string // the last occurrence of this text is where the error points; "" for end of input
	}{
		{"ALTER TABLE t DROP PRIMARY", ""},
		{"ALTER TABLE t DROP PRIMARY, ADD PRIMARY KEY (id)", ", ADD"},
		{"ALTER TABLE t DROP PRIMARY INDEX", "INDEX"},
		{"ALTER TABLE t DROP FOREIGN f", "f"},
		{"ALTER TABLE t DROP FOREIGN IF EXISTS f", "IF"},
		{"ALTER TABLE t DROP COLUMN m RESTRICT CASCADE", "CASCADE"},
		{"ALTER TABLE t DROP COLUMN m RESTRICT RESTRICT", "RESTRICT"},
		{"ALTER TABLE t DROP COLUMN restrict", "restrict"},
		{"ALTER TABLE t DROP cascade", "cascade"},
		{"ALTER TABLE t DROP INDEX i RESTRICT", "RESTRICT"},
		{"ALTER TABLE t DROP KEY IF EXISTS i CASCADE", "CASCADE"},
		{"ALTER TABLE t DROP FOREIGN KEY f RESTRICT", "RESTRICT"},
		{"ALTER TABLE t DROP CONSTRAINT c CASCADE", "CASCADE"},
		{"ALTER TABLE t DROP CONSTRAINT PRIMARY KEY RESTRICT", "RESTRICT"},
		{"ALTER TABLE t DROP PRIMARY KEY CASCADE", "CASCADE"},
		{"ALTER TABLE pt DROP PARTITION p0 RESTRICT", "RESTRICT"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %v, want a *ParseError", tc.sql, err)
			}
			want := len(tc.sql)
			if tc.at != "" {
				want = strings.LastIndex(tc.sql, tc.at)
			}
			if pe.Position != want {
				t.Errorf("error Position = %d, want %d: %v", pe.Position, want, err)
			}
		})
	}
}

// TestCollectAlterTableDropKey checks that completion offers KEY after DROP
// PRIMARY and DROP FOREIGN, and constraint_ref only after DROP FOREIGN KEY.
func TestCollectAlterTableDropKey(t *testing.T) {
	for _, sql := range []string{"ALTER TABLE t1 DROP PRIMARY ", "ALTER TABLE t1 DROP FOREIGN "} {
		t.Run(sql, func(t *testing.T) {
			cs := Collect(sql, len(sql))
			if cs == nil || !cs.HasToken(kwKEY) {
				t.Errorf("expected KEY after %q", sql)
			}
			if cs != nil && cs.HasRule("constraint_ref") {
				t.Errorf("constraint_ref offered before KEY after %q", sql)
			}
		})
	}
	sql := "ALTER TABLE t1 DROP FOREIGN KEY "
	if cs := Collect(sql, len(sql)); cs == nil || !cs.HasRule("constraint_ref") {
		t.Errorf("expected constraint_ref after %q", sql)
	}
}

// TestAlterTableQualifiedFieldIdent covers 2-part and 3-part field_ident
// qualifiers in ALTER TABLE and CREATE TABLE, including 3-part names on an
// unqualified table whose database qualifier is recorded in FieldDBs for
// session-aware callers (such as the catalog) to check against the current DB.
func TestAlterTableQualifiedFieldIdent(t *testing.T) {
	acceptCases := []struct {
		sql      string
		name     string
		ifExists bool
		fieldDBs []string
	}{
		{"ALTER TABLE a DROP a.a.a", "a", false, []string{"a"}},
		{"ALTER TABLE a DROP COLUMN a.a.a", "a", false, []string{"a"}},
		{"ALTER TABLE a DROP IF EXISTS a.a.a", "a", true, []string{"a"}},
		{"ALTER TABLE a DROP COLUMN IF EXISTS a.a.a RESTRICT", "a", true, []string{"a"}},
		{"ALTER TABLE a DROP COLUMN IF EXISTS a.a.a CASCADE", "a", true, []string{"a"}},
		{"ALTER TABLE a DROP d.a.c", "c", false, []string{"d"}},
		{"ALTER TABLE a DROP `d`.`a`.`c`", "c", false, []string{"d"}},
		{"ALTER TABLE a DROP a.c", "c", false, nil},
		{"ALTER TABLE a DROP .c", "c", false, nil},
		{"ALTER TABLE d.a DROP d.a.c", "c", false, nil},
		{"ALTER TABLE d.a DROP a.c", "c", false, nil},
	}
	for _, tc := range acceptCases {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseSeqStmt[*ast.AlterTableStmt](t, tc.sql)
			if len(stmt.Commands) != 1 {
				t.Fatalf("got %d commands, want 1", len(stmt.Commands))
			}
			cmd := stmt.Commands[0]
			if cmd.Type != ast.ATDropColumn || cmd.Name != tc.name || cmd.IfExists != tc.ifExists {
				t.Errorf("got Type=%v Name=%q IfExists=%v, want ATDropColumn Name=%q IfExists=%v",
					cmd.Type, cmd.Name, cmd.IfExists, tc.name, tc.ifExists)
			}
			if !slices.Equal(stmt.FieldDBs, tc.fieldDBs) {
				t.Errorf("FieldDBs = %v, want %v", stmt.FieldDBs, tc.fieldDBs)
			}
			wantLoc := ast.Loc{Start: strings.Index(tc.sql, "DROP"), End: len(tc.sql)}
			if cmd.Loc != wantLoc {
				t.Errorf("cmd Loc = %+v, want %+v", cmd.Loc, wantLoc)
			}
		})
	}

	t.Run("CREATE TABLE 3-part field_ident", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateTableStmt](t, "CREATE TABLE a (d.a.c INT)")
		if len(stmt.Columns) != 1 || stmt.Columns[0].Name != "c" {
			t.Fatalf("got columns %+v, want single column c", stmt.Columns)
		}
		if !slices.Equal(stmt.FieldDBs, []string{"d"}) {
			t.Errorf("FieldDBs = %v, want [d]", stmt.FieldDBs)
		}
	})

	t.Run("multi-command FieldDBs order", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.AlterTableStmt](t, "ALTER TABLE a DROP d1.a.c1, DROP d2.a.c2")
		if !slices.Equal(stmt.FieldDBs, []string{"d1", "d2"}) {
			t.Errorf("FieldDBs = %v, want [d1 d2]", stmt.FieldDBs)
		}
	})

	rejectCases := []string{
		"ALTER TABLE a DROP a.a.a.a",
		"ALTER TABLE a DROP COLUMN a.a.a.a",
		"ALTER TABLE a DROP a.a.a zzz",
		"ALTER TABLE a DROP a.b.a",
		"ALTER TABLE a DROP b.a",
		"ALTER TABLE d.a DROP x.a.c",
		"ALTER TABLE d.a DROP d.b.c",
		"ALTER TABLE a DROP ``.a.c",
		"ALTER TABLE a DROP `d `.a.c",
		"ALTER TABLE a DROP a.a.",
		"ALTER TABLE a DROP a..a",
		"CREATE TABLE a (d.b.c INT)",
		"CREATE TABLE d.a (x.a.c INT)",
	}
	for _, sql := range rejectCases {
		t.Run("reject/"+sql, func(t *testing.T) {
			ParseExpectError(t, sql)
		})
	}
}
