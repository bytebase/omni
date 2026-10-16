package parser

import (
	"errors"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TestParseSequenceDDL pins the AST of CREATE, ALTER, and DROP SEQUENCE.
// Acceptance itself is checked against TiDB in TestSequenceTiDBOracle.
func TestParseSequenceDDL(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			"CREATE SEQUENCE s",
			"{CREATE_SEQUENCE :loc 0 :name {TABLEREF :loc 16 :name s}}",
		},
		{
			"CREATE SEQUENCE IF NOT EXISTS db.s INCREMENT BY -2 START WITH 10 MINVALUE -100 MAXVALUE 1000 CACHE 20 CYCLE",
			"{CREATE_SEQUENCE :loc 0 :if_not_exists true :name {TABLEREF :loc 30 :schema db :name s}" +
				" :options {SEQUENCE_OPTION :increment_by -2} {SEQUENCE_OPTION :start_with 10}" +
				" {SEQUENCE_OPTION :minvalue -100} {SEQUENCE_OPTION :maxvalue 1000}" +
				" {SEQUENCE_OPTION :cache 20} {SEQUENCE_OPTION :cycle}}",
		},
		{
			// '=' in place of BY / WITH, and an explicit '+' sign.
			"CREATE SEQUENCE s INCREMENT = 2 START = +10 MINVALUE = 1 MAXVALUE = 9 CACHE = 3 NOCYCLE",
			"{CREATE_SEQUENCE :loc 0 :name {TABLEREF :loc 16 :name s}" +
				" :options {SEQUENCE_OPTION :increment_by 2} {SEQUENCE_OPTION :start_with 10}" +
				" {SEQUENCE_OPTION :minvalue 1} {SEQUENCE_OPTION :maxvalue 9}" +
				" {SEQUENCE_OPTION :cache 3} {SEQUENCE_OPTION :nocycle}}",
		},
		{
			"CREATE SEQUENCE s INCREMENT 2 START 10 NOMINVALUE NOMAXVALUE NOCACHE NO CYCLE",
			"{CREATE_SEQUENCE :loc 0 :name {TABLEREF :loc 16 :name s}" +
				" :options {SEQUENCE_OPTION :increment_by 2} {SEQUENCE_OPTION :start_with 10}" +
				" {SEQUENCE_OPTION :nominvalue} {SEQUENCE_OPTION :nomaxvalue}" +
				" {SEQUENCE_OPTION :nocache} {SEQUENCE_OPTION :nocycle}}",
		},
		{
			// Options keep source order and repeats; TiDB applies them in order.
			"CREATE SEQUENCE s NO MINVALUE NO MAXVALUE NO CACHE START WITH 1 START WITH 2",
			"{CREATE_SEQUENCE :loc 0 :name {TABLEREF :loc 16 :name s}" +
				" :options {SEQUENCE_OPTION :nominvalue} {SEQUENCE_OPTION :nomaxvalue}" +
				" {SEQUENCE_OPTION :nocache} {SEQUENCE_OPTION :start_with 1} {SEQUENCE_OPTION :start_with 2}}",
		},
		{
			// SignedNum bounds: '-' reaches -2^63, an unsigned value 2^63-1.
			"CREATE SEQUENCE s MINVALUE -9223372036854775808 MAXVALUE 9223372036854775807",
			"{CREATE_SEQUENCE :loc 0 :name {TABLEREF :loc 16 :name s}" +
				" :options {SEQUENCE_OPTION :minvalue -9223372036854775808} {SEQUENCE_OPTION :maxvalue 9223372036854775807}}",
		},
		{
			// Table options follow the sequence options, with or without commas.
			"CREATE SEQUENCE s START 1 COMMENT = 'ids', ENGINE = InnoDB COMMENT 'b'",
			"{CREATE_SEQUENCE :loc 0 :name {TABLEREF :loc 16 :name s}" +
				" :options {SEQUENCE_OPTION :start_with 1}" +
				" :table_options {TABLE_OPT :loc 26 :name COMMENT :val ids}" +
				" {TABLE_OPT :loc 43 :name ENGINE :val InnoDB} {TABLE_OPT :loc 59 :name COMMENT :val b}}",
		},
		{
			// The grammar takes every table option; TiDB's DDL rejects this
			// one with 8227, not a syntax error.
			"CREATE SEQUENCE s AUTO_INCREMENT = 5",
			"{CREATE_SEQUENCE :loc 0 :name {TABLEREF :loc 16 :name s}" +
				" :table_options {TABLE_OPT :loc 18 :name AUTO_INCREMENT :val 5}}",
		},
		{
			"ALTER SEQUENCE IF EXISTS s RESTART",
			"{ALTER_SEQUENCE :loc 0 :if_exists true :name {TABLEREF :loc 25 :name s}" +
				" :options {SEQUENCE_OPTION :restart}}",
		},
		{
			"ALTER SEQUENCE s RESTART WITH 5 RESTART = -6 RESTART 7 RESTART +8 INCREMENT BY 3 NOCACHE",
			"{ALTER_SEQUENCE :loc 0 :name {TABLEREF :loc 15 :name s}" +
				" :options {SEQUENCE_OPTION :restart_with 5} {SEQUENCE_OPTION :restart_with -6}" +
				" {SEQUENCE_OPTION :restart_with 7} {SEQUENCE_OPTION :restart_with 8}" +
				" {SEQUENCE_OPTION :increment_by 3} {SEQUENCE_OPTION :nocache}}",
		},
		{
			"ALTER SEQUENCE s START WITH 3 RESTART NO CYCLE",
			"{ALTER_SEQUENCE :loc 0 :name {TABLEREF :loc 15 :name s}" +
				" :options {SEQUENCE_OPTION :start_with 3} {SEQUENCE_OPTION :restart} {SEQUENCE_OPTION :nocycle}}",
		},
		{
			"DROP SEQUENCE s",
			"{DROP_SEQUENCE :loc 0 :sequences {TABLEREF :loc 14 :name s}}",
		},
		{
			"DROP SEQUENCE IF EXISTS a, db.b, `c d`",
			"{DROP_SEQUENCE :loc 0 :if_exists true :sequences {TABLEREF :loc 24 :name a}" +
				" {TABLEREF :loc 27 :schema db :name b} {TABLEREF :loc 33 :name c d}}",
		},
		{
			// The sequence keywords are unreserved, so they name objects.
			"CREATE SEQUENCE sequence",
			"{CREATE_SEQUENCE :loc 0 :name {TABLEREF :loc 16 :name sequence}}",
		},
		{
			"DROP SEQUENCE nocache, cycle",
			"{DROP_SEQUENCE :loc 0 :sequences {TABLEREF :loc 14 :name nocache} {TABLEREF :loc 23 :name cycle}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if len(list.Items) != 1 {
				t.Fatalf("got %d statements, want 1", len(list.Items))
			}
			if got := nodes.NodeToString(list.Items[0]); got != tt.want {
				t.Errorf("AST mismatch\n got: %s\nwant: %s", got, tt.want)
			}
		})
	}
}

// TestParseSequenceRejects pins where each rejected statement fails.
// TestSequenceTiDBOracle confirms TiDB rejects these as syntax errors.
func TestParseSequenceRejects(t *testing.T) {
	tests := []struct {
		sql string
		pos int // byte offset of the ParseError
	}{
		{"CREATE OR REPLACE SEQUENCE s", 18},
		{"CREATE TEMPORARY SEQUENCE s", 17},
		{"CREATE SEQUENCE IF EXISTS s", 19},
		{"CREATE SEQUENCE s AS BIGINT", 18},
		{"CREATE SEQUENCE s RESTART", 18},
		{"CREATE SEQUENCE s CACHE", 23},
		{"CREATE SEQUENCE s MAXVALUE", 26},
		{"CREATE SEQUENCE s NO", 20},
		{"CREATE SEQUENCE s NO START", 21},
		{"CREATE SEQUENCE s START WITH 1.5", 29},
		{"CREATE SEQUENCE s START WITH '1'", 29},
		{"CREATE SEQUENCE s START WITH (1)", 29},
		{"CREATE SEQUENCE s START WITH 0x10", 29},
		{"CREATE SEQUENCE s START WITH = 1", 29},
		{"CREATE SEQUENCE s INCREMENT BY = 1", 31},
		{"CREATE SEQUENCE s START 9223372036854775808", 24},
		{"CREATE SEQUENCE s MINVALUE -9223372036854775809", 28},
		{"CREATE SEQUENCE s MINVALUE -18446744073709551616", 28},
		{"CREATE SEQUENCE s COMMENT 'x',", 30},
		{"ALTER SEQUENCE s", 16},
		{"ALTER SEQUENCE IF NOT EXISTS s RESTART", 18},
		{"ALTER SEQUENCE s RESTART WITH", 29},
		{"ALTER SEQUENCE s RESTART WITH = 5", 30},
		{"DROP SEQUENCE", 13},
		{"DROP SEQUENCE s,", 16},
		{"DROP TEMPORARY SEQUENCE s", 15},
		{"DROP SEQUENCE IF NOT EXISTS s", 17},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Fatal("Parse accepted, want a syntax error")
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("error %T is not a *ParseError: %v", err, err)
			}
			if pe.Position != tt.pos {
				t.Errorf("error position = %d, want %d (%v)", pe.Position, tt.pos, err)
			}
		})
	}
}

// TestSequenceLoc checks that statement and option Locs end at their last
// token and stay absolute in a multi-statement script.
func TestSequenceLoc(t *testing.T) {
	sql := "  CREATE SEQUENCE s CACHE 5  ;\nALTER SEQUENCE b RESTART WITH 3 ;\nDROP SEQUENCE a, b ;"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list.Items) != 3 {
		t.Fatalf("got %d statements, want 3", len(list.Items))
	}

	create := list.Items[0].(*nodes.CreateSequenceStmt)
	checkLoc(t, sql, "create", create.Loc, "CREATE SEQUENCE s CACHE 5")
	checkLoc(t, sql, "create option", create.Options[0].Loc, "CACHE 5")

	alter := list.Items[1].(*nodes.AlterSequenceStmt)
	checkLoc(t, sql, "alter", alter.Loc, "ALTER SEQUENCE b RESTART WITH 3")
	checkLoc(t, sql, "alter option", alter.Options[0].Loc, "RESTART WITH 3")

	drop := list.Items[2].(*nodes.DropSequenceStmt)
	checkLoc(t, sql, "drop", drop.Loc, "DROP SEQUENCE a, b")

	if v := CheckLocations(t, sql); len(v) > 0 {
		t.Errorf("Loc violations: %v", v)
	}
}

func checkLoc(t *testing.T, sql, what string, loc nodes.Loc, want string) {
	t.Helper()
	if loc.Start < 0 || loc.End > len(sql) || loc.Start >= loc.End {
		t.Errorf("%s: invalid Loc %+v", what, loc)
		return
	}
	if got := sql[loc.Start:loc.End]; got != want {
		t.Errorf("%s: Loc %+v covers %q, want %q", what, loc, got, want)
	}
}

// TestSequenceKeywordsAsIdentifiers guards the new keyword tokens: TiDB
// keeps all of them unreserved (parser.y:7077-7085).
func TestSequenceKeywordsAsIdentifiers(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE sequence (increment INT, minvalue INT, nominvalue INT, nomaxvalue INT, nocache INT, cycle INT, nocycle INT)",
		"SELECT increment, minvalue, nominvalue, nomaxvalue, nocache, cycle, nocycle FROM sequence",
		"SELECT s.cycle AS increment FROM sequence AS s",
		"CREATE SEQUENCE nocache",
		"ALTER SEQUENCE minvalue RESTART",
		"DROP SEQUENCE increment, nocycle",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

// TestSequenceTruncationNoPanic parses every prefix of sequence DDL; each
// must return a result or a ParseError, never panic.
func TestSequenceTruncationNoPanic(t *testing.T) {
	for _, sql := range []string{
		"CREATE SEQUENCE IF NOT EXISTS db.s INCREMENT BY -2 START WITH +10 NO MINVALUE MAXVALUE = 9 NOCACHE CYCLE COMMENT 'x', ENGINE = InnoDB",
		"ALTER SEQUENCE IF EXISTS s RESTART WITH -5 NO CYCLE RESTART",
		"DROP SEQUENCE IF EXISTS a, db.b",
	} {
		for i := 0; i <= len(sql); i++ {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Parse(%q) panicked: %v", sql[:i], r)
					}
				}()
				_, err := Parse(sql[:i])
				var pe *ParseError
				if err != nil && !errors.As(err, &pe) {
					t.Errorf("Parse(%q) returned %T, want *ParseError", sql[:i], err)
				}
			}()
		}
	}
}
