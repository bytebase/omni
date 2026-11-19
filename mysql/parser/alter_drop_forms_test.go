package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mysql/ast"
)

// TestAlterTableDropForms covers the DROP arms of MySQL 8.0's
// alter_list_item: RESTRICT or CASCADE after a column (opt_restrict, which the
// server ignores), and KEY after PRIMARY and FOREIGN. The command spans its
// last token, so a trailing RESTRICT or CASCADE is inside its Loc.
// Container-verified by TestAlterTableDropFormsOracle.
func TestAlterTableDropForms(t *testing.T) {
	cases := []struct {
		sql  string
		typ  ast.AlterTableCmdType
		name string
	}{
		{"ALTER TABLE t DROP COLUMN m RESTRICT", ast.ATDropColumn, "m"},
		{"ALTER TABLE t DROP COLUMN m CASCADE", ast.ATDropColumn, "m"},
		{"ALTER TABLE t DROP m RESTRICT", ast.ATDropColumn, "m"},
		{"ALTER TABLE t DROP k CASCADE", ast.ATDropColumn, "k"},
		{"ALTER TABLE t DROP COLUMN m /* c */ restrict", ast.ATDropColumn, "m"},
		{"ALTER TABLE t DROP COLUMN `restrict`", ast.ATDropColumn, "restrict"},
		{"ALTER TABLE t DROP COLUMN m", ast.ATDropColumn, "m"},
		{"ALTER TABLE t DROP PRIMARY KEY", ast.ATDropConstraint, "PRIMARY"},
		{"ALTER TABLE t DROP FOREIGN KEY f", ast.ATDropConstraint, "f"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseAlterDropFormsStmt(t, tc.sql)
			if len(stmt.Commands) != 1 {
				t.Fatalf("got %d commands, want 1", len(stmt.Commands))
			}
			cmd := stmt.Commands[0]
			if cmd.Type != tc.typ || cmd.Name != tc.name || cmd.IfExists {
				t.Errorf("got Type=%v Name=%q IfExists=%v, want Type=%v Name=%q IfExists=false",
					cmd.Type, cmd.Name, cmd.IfExists, tc.typ, tc.name)
			}
			wantLoc := ast.Loc{Start: strings.Index(tc.sql, "DROP"), End: len(tc.sql)}
			if cmd.Loc != wantLoc {
				t.Errorf("cmd Loc = %+v, want %+v", cmd.Loc, wantLoc)
			}
		})
	}

	t.Run("several commands", func(t *testing.T) {
		sql := "ALTER TABLE t DROP COLUMN m RESTRICT, DROP k CASCADE, ADD COLUMN z INT"
		stmt := parseAlterDropFormsStmt(t, sql)
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

// TestAlterTableDropFormsReject pins the DROP forms MySQL 8.0 rejects with
// 1064, container-verified by TestAlterTableDropFormsOracle. IF EXISTS on a
// column is MariaDB syntax, and only a column takes RESTRICT or CASCADE.
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
		{"ALTER TABLE t DROP COLUMN IF EXISTS m", "IF"},
		{"ALTER TABLE t DROP IF EXISTS m", "IF"},
		{"ALTER TABLE t DROP COLUMN IF EXISTS m RESTRICT", "IF"},
		{"ALTER TABLE t DROP COLUMN m RESTRICT CASCADE", "CASCADE"},
		{"ALTER TABLE t DROP COLUMN m RESTRICT RESTRICT", "RESTRICT"},
		{"ALTER TABLE t DROP COLUMN restrict", "restrict"},
		{"ALTER TABLE t DROP cascade", "cascade"},
		{"ALTER TABLE t DROP INDEX i RESTRICT", "RESTRICT"},
		{"ALTER TABLE t DROP KEY i CASCADE", "CASCADE"},
		{"ALTER TABLE t DROP FOREIGN KEY f RESTRICT", "RESTRICT"},
		{"ALTER TABLE t DROP CONSTRAINT c CASCADE", "CASCADE"},
		{"ALTER TABLE t DROP CHECK c RESTRICT", "RESTRICT"},
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

// TestAlterTableConstraintKeyword keeps the keyword a DROP or ALTER action
// wrote on the command, and in the AST's text form, which is what the parser's
// compare tests hold a statement's tree to. DROP CHECK and DROP CONSTRAINT,
// and ALTER CHECK and ALTER CONSTRAINT, are one action in the grammar that
// mysql:8.0.46 answers with different codes when the constraint is not there
// (3821 for CHECK, 3940 for CONSTRAINT, read on the pinned server on
// 2026-11-20), so a tree without the keyword cannot answer what the statement
// asked. The catalog reads the field; TestSimulationCausesOracle holds it to
// the server.
func TestAlterTableConstraintKeyword(t *testing.T) {
	cases := []struct {
		sql  string
		want ast.ConstraintKeyword
		word string
	}{
		{"ALTER TABLE t DROP CHECK c", ast.ConstraintKeywordCheck, "CHECK"},
		{"ALTER TABLE t DROP CONSTRAINT c", ast.ConstraintKeywordConstraint, "CONSTRAINT"},
		{"ALTER TABLE t DROP PRIMARY KEY", ast.ConstraintKeywordPrimaryKey, "PRIMARY KEY"},
		{"ALTER TABLE t DROP FOREIGN KEY f", ast.ConstraintKeywordForeignKey, "FOREIGN KEY"},
		{"ALTER TABLE t ALTER CHECK c ENFORCED", ast.ConstraintKeywordCheck, "CHECK"},
		{"ALTER TABLE t ALTER CONSTRAINT c NOT ENFORCED", ast.ConstraintKeywordConstraint, "CONSTRAINT"},
		// The keywords of their own command types are not this field's.
		{"ALTER TABLE t DROP COLUMN c", ast.ConstraintKeywordNone, ""},
		{"ALTER TABLE t DROP INDEX i", ast.ConstraintKeywordNone, ""},
	}
	dumps := map[string]string{}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseAlterDropFormsStmt(t, tc.sql)
			if len(stmt.Commands) != 1 {
				t.Fatalf("got %d commands, want 1", len(stmt.Commands))
			}
			cmd := stmt.Commands[0]
			if cmd.ConstraintKeyword != tc.want {
				t.Fatalf("%s: ConstraintKeyword = %v, want %v", tc.sql, cmd.ConstraintKeyword, tc.want)
			}
			dump := ast.NodeToString(parseList(t, tc.sql))
			if tc.word == "" {
				if strings.Contains(dump, ":constraint_keyword") {
					t.Errorf("%s: tree holds a keyword it wrote none of: %s", tc.sql, dump)
				}
				dumps[tc.sql] = dump
				return
			}
			if want := " :constraint_keyword " + tc.word; !strings.Contains(dump, want) {
				t.Errorf("%s: tree lost %q: %s", tc.sql, want, dump)
			}
			// The text form is a fixed point: the same statement parses to
			// the same tree, keyword included.
			if again := ast.NodeToString(parseList(t, tc.sql)); again != dump {
				t.Errorf("%s: second parse dumped %s, first %s", tc.sql, again, dump)
			}
			dumps[tc.sql] = dump
		})
	}
	// The two spellings of one action keep different trees, which is the whole
	// point of the field.
	if dumps["ALTER TABLE t DROP CHECK c"] == dumps["ALTER TABLE t DROP CONSTRAINT c"] {
		t.Errorf("DROP CHECK and DROP CONSTRAINT of the same name dump one tree: %s", dumps["ALTER TABLE t DROP CHECK c"])
	}
}

func parseList(t *testing.T, sql string) *ast.List {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	return list
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

func parseAlterDropFormsStmt(t *testing.T, sql string) *ast.AlterTableStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if list.Len() != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", sql, list.Len())
	}
	stmt, ok := list.Items[0].(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("Parse(%q) = %T, want *ast.AlterTableStmt", sql, list.Items[0])
	}
	return stmt
}
