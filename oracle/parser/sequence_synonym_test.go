package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCreateSequence(t *testing.T) {
	p := newTestParser("SEQUENCE emp_seq")
	stmt, parseErr1 := p.parseCreateSequenceStmt(0)
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	if stmt == nil {
		t.Fatal("expected CreateSequenceStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "EMP_SEQ" {
		t.Errorf("expected sequence name EMP_SEQ, got %v", stmt.Name)
	}
}

func TestParseCreateSequenceWithOptions(t *testing.T) {
	p := newTestParser("SEQUENCE hr.emp_seq INCREMENT BY 1 START WITH 100 MAXVALUE 999999 MINVALUE 1 CYCLE CACHE 20 ORDER")
	stmt, parseErr2 := p.parseCreateSequenceStmt(0)
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	if stmt.Name == nil || stmt.Name.Schema != "HR" || stmt.Name.Name != "EMP_SEQ" {
		t.Errorf("expected HR.EMP_SEQ, got %v", stmt.Name)
	}
	if stmt.IncrementBy == nil {
		t.Error("expected non-nil IncrementBy")
	} else {
		n, ok := stmt.IncrementBy.(*ast.NumberLiteral)
		if !ok || n.Ival != 1 {
			t.Errorf("expected IncrementBy=1, got %v", stmt.IncrementBy)
		}
	}
	if stmt.StartWith == nil {
		t.Error("expected non-nil StartWith")
	} else {
		n := stmt.StartWith.(*ast.NumberLiteral)
		if n.Ival != 100 {
			t.Errorf("expected StartWith=100, got %v", n.Ival)
		}
	}
	if stmt.MaxValue == nil {
		t.Error("expected non-nil MaxValue")
	}
	if stmt.MinValue == nil {
		t.Error("expected non-nil MinValue")
	}
	if !stmt.Cycle {
		t.Error("expected Cycle to be true")
	}
	if stmt.Cache == nil {
		t.Error("expected non-nil Cache")
	}
	if !stmt.Order {
		t.Error("expected Order to be true")
	}
}

func TestParseCreateSequenceNoOptions(t *testing.T) {
	p := newTestParser("SEQUENCE s NOMAXVALUE NOMINVALUE NOCYCLE NOCACHE NOORDER")
	stmt, parseErr3 := p.parseCreateSequenceStmt(0)
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	if !stmt.NoMaxValue {
		t.Error("expected NoMaxValue to be true")
	}
	if !stmt.NoMinValue {
		t.Error("expected NoMinValue to be true")
	}
	if !stmt.NoCycle {
		t.Error("expected NoCycle to be true")
	}
	if !stmt.NoCache {
		t.Error("expected NoCache to be true")
	}
	if !stmt.NoOrder {
		t.Error("expected NoOrder to be true")
	}
}

func TestParseCreateSequenceLoc(t *testing.T) {
	p := newTestParser("SEQUENCE s")
	stmt, parseErr4 := p.parseCreateSequenceStmt(0)
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}

func TestParseCreateSynonym(t *testing.T) {
	p := newTestParser("SYNONYM emp FOR hr.employees")
	stmt, parseErr5 := p.parseCreateSynonymStmt(0, false, false)
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	if stmt == nil {
		t.Fatal("expected CreateSynonymStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "EMP" {
		t.Errorf("expected synonym name EMP, got %v", stmt.Name)
	}
	if stmt.Target == nil || stmt.Target.Schema != "HR" || stmt.Target.Name != "EMPLOYEES" {
		t.Errorf("expected target HR.EMPLOYEES, got %v", stmt.Target)
	}
}

func TestParseCreatePublicSynonym(t *testing.T) {
	p := newTestParser("SYNONYM emp FOR hr.employees")
	stmt, parseErr6 := p.parseCreateSynonymStmt(0, false, true)
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	if !stmt.Public {
		t.Error("expected Public to be true")
	}
}

func TestParseCreateOrReplaceSynonym(t *testing.T) {
	p := newTestParser("SYNONYM emp FOR hr.employees")
	stmt, parseErr7 := p.parseCreateSynonymStmt(0, true, false)
	if parseErr7 != nil {
		t.Fatalf("parse: %v", parseErr7)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
}

func TestParseCreateSynonymLoc(t *testing.T) {
	p := newTestParser("SYNONYM s FOR t")
	stmt, parseErr8 := p.parseCreateSynonymStmt(0, false, false)
	if parseErr8 != nil {
		t.Fatalf("parse: %v", parseErr8)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}

func TestParseCreateDatabaseLink(t *testing.T) {
	p := newTestParser("DATABASE LINK remote_db CONNECT TO admin IDENTIFIED BY secret USING 'prod_server'")
	stmt, parseErr9 := p.parseCreateDatabaseLinkStmt(0, false)
	if parseErr9 != nil {
		t.Fatalf("parse: %v", parseErr9)
	}
	if stmt == nil {
		t.Fatal("expected CreateDatabaseLinkStmt, got nil")
	}
	if stmt.Name != "REMOTE_DB" {
		t.Errorf("expected link name REMOTE_DB, got %q", stmt.Name)
	}
	if stmt.ConnectTo != "ADMIN" {
		t.Errorf("expected ConnectTo ADMIN, got %q", stmt.ConnectTo)
	}
	if stmt.Identified != "SECRET" {
		t.Errorf("expected Identified SECRET, got %q", stmt.Identified)
	}
	if stmt.Using != "prod_server" {
		t.Errorf("expected Using prod_server, got %q", stmt.Using)
	}
}

func TestParseCreatePublicDatabaseLink(t *testing.T) {
	p := newTestParser("DATABASE LINK remote_db CONNECT TO admin IDENTIFIED BY pass USING 'srv'")
	stmt, parseErr10 := p.parseCreateDatabaseLinkStmt(0, true)
	if parseErr10 != nil {
		t.Fatalf("parse: %v", parseErr10)
	}
	if !stmt.Public {
		t.Error("expected Public to be true")
	}
}

func TestParseCreateDatabaseLinkLoc(t *testing.T) {
	p := newTestParser("DATABASE LINK db CONNECT TO u IDENTIFIED BY p USING 's'")
	stmt, parseErr11 := p.parseCreateDatabaseLinkStmt(0, false)
	if parseErr11 != nil {
		t.Fatalf("parse: %v", parseErr11)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}

// TestParseCreateSequenceOptions covers every CREATE SEQUENCE option group in
// any order, as Oracle 23ai accepts them.
func TestParseCreateSequenceOptions(t *testing.T) {
	tests := []struct {
		sql  string
		want func(*ast.CreateSequenceStmt) bool
	}{
		{
			sql: "CREATE SEQUENCE app.seq_search\n    START WITH 1 INCREMENT BY 1 CACHE 20 NOORDER NOCYCLE NOKEEP NOSCALE GLOBAL",
			want: func(s *ast.CreateSequenceStmt) bool {
				return s.NoOrder && s.NoCycle && s.NoKeep && s.NoScale && s.Global && s.Cache != nil
			},
		},
		{
			sql: "CREATE SEQUENCE s KEEP SCALE EXTEND SESSION",
			want: func(s *ast.CreateSequenceStmt) bool {
				return s.Keep && s.Scale && s.ScaleExtend && s.Session
			},
		},
		{
			sql:  "CREATE SEQUENCE s SCALE NOEXTEND",
			want: func(s *ast.CreateSequenceStmt) bool { return s.Scale && s.ScaleNoExtend },
		},
		{
			sql:  "CREATE SEQUENCE s SCALE START WITH 5",
			want: func(s *ast.CreateSequenceStmt) bool { return s.Scale && !s.ScaleExtend && s.StartWith != nil },
		},
		// Oracle parses SHARD and NOSHARD, then raises ORA-02511 outside a
		// sharded database.
		{
			sql:  "CREATE SEQUENCE s SHARD EXTEND",
			want: func(s *ast.CreateSequenceStmt) bool { return s.Shard && s.ShardExtend },
		},
		{
			sql:  "CREATE SEQUENCE s NOSHARD",
			want: func(s *ast.CreateSequenceStmt) bool { return s.NoShard },
		},
		// One EXTEND or NOEXTEND, after either clause, covers SCALE and SHARD.
		{
			sql:  "CREATE SEQUENCE s SCALE SHARD EXTEND",
			want: func(s *ast.CreateSequenceStmt) bool { return s.Scale && s.Shard && s.ShardExtend && !s.ScaleExtend },
		},
		{
			sql:  "CREATE SEQUENCE s SCALE NOEXTEND SHARD",
			want: func(s *ast.CreateSequenceStmt) bool { return s.ScaleNoExtend && s.Shard && !s.ShardNoExtend },
		},
		{
			sql:  "CREATE SEQUENCE IF NOT EXISTS s SHARING = NONE START WITH 1",
			want: func(s *ast.CreateSequenceStmt) bool { return s.IfNotExists && s.Sharing == "NONE" },
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateSequenceStmt)
			if !tt.want(stmt) {
				t.Fatalf("unexpected options: %s", ast.NodeToString(stmt))
			}
			if violations := CheckLocations(t, tt.sql); len(violations) > 0 {
				t.Fatalf("Loc violations: %v", violations)
			}
		})
	}
}

// TestParseCreateSequenceRejects lists forms Oracle 23ai rejects: a repeated
// or conflicting option group (ORA-02278 through ORA-02285, ORA-03178,
// ORA-41415, ORA-64600), SHARING or IF NOT EXISTS out of place (ORA-03049),
// and KEEP with a value (ORA-03047).
func TestParseCreateSequenceRejects(t *testing.T) {
	tests := []string{
		"CREATE SEQUENCE s INCREMENT BY 1 INCREMENT BY 2",
		"CREATE SEQUENCE s START WITH 1 START WITH 2",
		"CREATE SEQUENCE s MAXVALUE 5 NOMAXVALUE",
		"CREATE SEQUENCE s MINVALUE 1 MINVALUE 2",
		"CREATE SEQUENCE s CYCLE NOCYCLE",
		"CREATE SEQUENCE s CACHE 5 NOCACHE",
		"CREATE SEQUENCE s ORDER NOORDER",
		"CREATE SEQUENCE s KEEP NOKEEP",
		"CREATE SEQUENCE s NOSCALE SCALE",
		"CREATE SEQUENCE s SESSION GLOBAL",
		"CREATE SEQUENCE s NOSHARD SHARD",
		"CREATE SEQUENCE s START WITH 1 SHARING = NONE",
		"CREATE SEQUENCE s SHARING = ALL",
		"CREATE SEQUENCE s IF NOT EXISTS",
		"CREATE SEQUENCE s KEEP 5",
		// One EXTEND or NOEXTEND covers both SCALE and SHARD; writing it for
		// each is a "duplicate or conflicting EXTEND clause" parsing error per
		// the ALTER SEQUENCE reference. A non-sharded Oracle raises ORA-02511
		// at SHARD first, so these rest on the documentation.
		"CREATE SEQUENCE s SCALE EXTEND SHARD NOEXTEND",
		"CREATE SEQUENCE s SCALE EXTEND SHARD EXTEND",
		"CREATE SEQUENCE s SHARD NOEXTEND SCALE NOEXTEND",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
}
