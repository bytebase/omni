package parser

import (
	"errors"
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
