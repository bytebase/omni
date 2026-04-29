package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseTypeNumber tests NUMBER, NUMBER(p), NUMBER(p,s).
func TestParseTypeNumber(t *testing.T) {
	tests := []struct {
		input    string
		names    []string
		modCount int
	}{
		{"NUMBER", []string{"NUMBER"}, 0},
		{"NUMBER(10)", []string{"NUMBER"}, 1},
		{"NUMBER(10,2)", []string{"NUMBER"}, 2},
		{"INTEGER", []string{"INTEGER"}, 0},
		{"SMALLINT", []string{"SMALLINT"}, 0},
		{"FLOAT", []string{"FLOAT"}, 0},
		{"FLOAT(53)", []string{"FLOAT"}, 1},
		{"DECIMAL(10,2)", []string{"DECIMAL"}, 2},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			tn, parseErr1 := p.parseTypeName()
			if parseErr1 != nil {
				t.Fatalf("parse: %v", parseErr1)
			}
			if tn == nil {
				t.Fatal("expected non-nil TypeName")
			}
			if tn.Names.Len() != len(tc.names) {
				t.Fatalf("expected %d name parts, got %d", len(tc.names), tn.Names.Len())
			}
			for i, n := range tc.names {
				s := tn.Names.Items[i].(*ast.String).Str
				if s != n {
					t.Errorf("name[%d]: expected %q, got %q", i, n, s)
				}
			}
			if tn.TypeMods.Len() != tc.modCount {
				t.Errorf("expected %d type mods, got %d", tc.modCount, tn.TypeMods.Len())
			}
		})
	}
}

// TestParseTypeChar tests CHAR, VARCHAR, VARCHAR2, NCHAR, NVARCHAR2.
func TestParseTypeChar(t *testing.T) {
	tests := []struct {
		input    string
		name     string
		modCount int
	}{
		{"CHAR", "CHAR", 0},
		{"CHAR(100)", "CHAR", 1},
		{"VARCHAR2(255)", "VARCHAR2", 1},
		{"VARCHAR(100)", "VARCHAR", 1},
		{"NCHAR(50)", "NCHAR", 1},
		{"NVARCHAR2(200)", "NVARCHAR2", 1},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			tn, parseErr2 := p.parseTypeName()
			if parseErr2 != nil {
				t.Fatalf("parse: %v", parseErr2)
			}
			if tn == nil {
				t.Fatal("expected non-nil TypeName")
			}
			s := tn.Names.Items[0].(*ast.String).Str
			if s != tc.name {
				t.Errorf("expected name %q, got %q", tc.name, s)
			}
			if tn.TypeMods.Len() != tc.modCount {
				t.Errorf("expected %d type mods, got %d", tc.modCount, tn.TypeMods.Len())
			}
		})
	}
}

// TestParseTypeLOB tests CLOB, BLOB, NCLOB.
func TestParseTypeLOB(t *testing.T) {
	for _, input := range []string{"CLOB", "BLOB", "NCLOB"} {
		t.Run(input, func(t *testing.T) {
			p := newTestParser(input)
			tn, parseErr3 := p.parseTypeName()
			if parseErr3 != nil {
				t.Fatalf("parse: %v", parseErr3)
			}
			if tn == nil {
				t.Fatal("expected non-nil TypeName")
			}
			s := tn.Names.Items[0].(*ast.String).Str
			if s != input {
				t.Errorf("expected name %q, got %q", input, s)
			}
		})
	}
}

// TestParseTypeDatetime tests DATE, TIMESTAMP, TIMESTAMP WITH TIME ZONE.
func TestParseTypeDatetime(t *testing.T) {
	tests := []struct {
		input    string
		names    []string
		modCount int
	}{
		{"DATE", []string{"DATE"}, 0},
		{"TIMESTAMP", []string{"TIMESTAMP"}, 0},
		{"TIMESTAMP(6)", []string{"TIMESTAMP"}, 1},
		{"TIMESTAMP WITH TIME ZONE", []string{"TIMESTAMP", "WITH", "TIME", "ZONE"}, 0},
		{"TIMESTAMP(6) WITH LOCAL TIME ZONE", []string{"TIMESTAMP", "WITH", "LOCAL", "TIME", "ZONE"}, 1},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			tn, parseErr4 := p.parseTypeName()
			if parseErr4 != nil {
				t.Fatalf("parse: %v", parseErr4)
			}
			if tn == nil {
				t.Fatal("expected non-nil TypeName")
			}
			if tn.Names.Len() != len(tc.names) {
				t.Fatalf("expected %d name parts, got %d", len(tc.names), tn.Names.Len())
			}
			for i, n := range tc.names {
				s := tn.Names.Items[i].(*ast.String).Str
				if s != n {
					t.Errorf("name[%d]: expected %q, got %q", i, n, s)
				}
			}
			if tn.TypeMods.Len() != tc.modCount {
				t.Errorf("expected %d type mods, got %d", tc.modCount, tn.TypeMods.Len())
			}
		})
	}
}

// TestParseTypeInterval tests INTERVAL types.
func TestParseTypeInterval(t *testing.T) {
	tests := []struct {
		input string
		names []string
	}{
		{"INTERVAL YEAR TO MONTH", []string{"INTERVAL", "YEAR", "TO", "MONTH"}},
		{"INTERVAL DAY TO SECOND", []string{"INTERVAL", "DAY", "TO", "SECOND"}},
		{"INTERVAL YEAR(4) TO MONTH", []string{"INTERVAL", "YEAR", "TO", "MONTH"}},
		{"INTERVAL DAY(2) TO SECOND(6)", []string{"INTERVAL", "DAY", "TO", "SECOND"}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			tn, parseErr5 := p.parseTypeName()
			if parseErr5 != nil {
				t.Fatalf("parse: %v", parseErr5)
			}
			if tn == nil {
				t.Fatal("expected non-nil TypeName")
			}
			if tn.Names.Len() != len(tc.names) {
				t.Fatalf("expected %d name parts, got %d", len(tc.names), tn.Names.Len())
			}
			for i, n := range tc.names {
				s := tn.Names.Items[i].(*ast.String).Str
				if s != n {
					t.Errorf("name[%d]: expected %q, got %q", i, n, s)
				}
			}
		})
	}
}

// TestParseTypeRowid tests ROWID type.
func TestParseTypeRowid(t *testing.T) {
	p := newTestParser("ROWID")
	tn, parseErr1 := p.parseTypeName()
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	if tn == nil {
		t.Fatal("expected non-nil TypeName")
	}
	s := tn.Names.Items[0].(*ast.String).Str
	if s != "ROWID" {
		t.Errorf("expected ROWID, got %q", s)
	}
}

// TestParseTypeRAW tests RAW(n).
func TestParseTypeRAW(t *testing.T) {
	p := newTestParser("RAW(2000)")
	tn, parseErr2 := p.parseTypeName()
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	if tn == nil {
		t.Fatal("expected non-nil TypeName")
	}
	s := tn.Names.Items[0].(*ast.String).Str
	if s != "RAW" {
		t.Errorf("expected RAW, got %q", s)
	}
	if tn.TypeMods.Len() != 1 {
		t.Errorf("expected 1 type mod, got %d", tn.TypeMods.Len())
	}
}

// TestParseTypeLong tests LONG and LONG RAW.
func TestParseTypeLong(t *testing.T) {
	tests := []struct {
		input string
		names []string
	}{
		{"LONG", []string{"LONG"}},
		{"LONG RAW", []string{"LONG", "RAW"}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			tn, parseErr6 := p.parseTypeName()
			if parseErr6 != nil {
				t.Fatalf("parse: %v", parseErr6)
			}
			if tn == nil {
				t.Fatal("expected non-nil TypeName")
			}
			if tn.Names.Len() != len(tc.names) {
				t.Fatalf("expected %d name parts, got %d", len(tc.names), tn.Names.Len())
			}
		})
	}
}

// TestParseTypePctType tests variable%TYPE.
func TestParseTypePctType(t *testing.T) {
	p := newTestParser("employees.salary%TYPE")
	tn, parseErr3 := p.parseTypeName()
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	if tn == nil {
		t.Fatal("expected non-nil TypeName")
	}
	if !tn.IsPercType {
		t.Error("expected IsPercType=true")
	}
}

// TestParseTypePctRowtype tests cursor%ROWTYPE.
func TestParseTypePctRowtype(t *testing.T) {
	p := newTestParser("employees%ROWTYPE")
	tn, parseErr4 := p.parseTypeName()
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	if tn == nil {
		t.Fatal("expected non-nil TypeName")
	}
	if !tn.IsPercRowtype {
		t.Error("expected IsPercRowtype=true")
	}
}

// TestParseTypeUserDefined tests user-defined type names.
func TestParseTypeUserDefined(t *testing.T) {
	p := newTestParser("my_schema.my_type")
	tn, parseErr5 := p.parseTypeName()
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	if tn == nil {
		t.Fatal("expected non-nil TypeName")
	}
	if tn.Names.Len() != 2 {
		t.Fatalf("expected 2 name parts, got %d", tn.Names.Len())
	}
	s0 := tn.Names.Items[0].(*ast.String).Str
	s1 := tn.Names.Items[1].(*ast.String).Str
	if s0 != "MY_SCHEMA" {
		t.Errorf("expected MY_SCHEMA, got %q", s0)
	}
	if s1 != "MY_TYPE" {
		t.Errorf("expected MY_TYPE, got %q", s1)
	}
}

// TestParseTypeLoc tests that location is recorded on type names.
func TestParseTypeLoc(t *testing.T) {
	p := newTestParser("NUMBER(10,2)")
	tn, parseErr6 := p.parseTypeName()
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	if tn.Loc.Start != 0 {
		t.Errorf("expected Start=0, got %d", tn.Loc.Start)
	}
	if tn.Loc.End <= tn.Loc.Start {
		t.Errorf("expected End > Start, got End=%d Start=%d", tn.Loc.End, tn.Loc.Start)
	}
}

// TestParseTypeCharWithByteSemantic tests CHAR(n BYTE) and CHAR(n CHAR).
func TestParseTypeCharWithByteSemantic(t *testing.T) {
	tests := []struct {
		input    string
		modCount int
	}{
		{"CHAR(100 BYTE)", 2},
		{"VARCHAR2(255 CHAR)", 2},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			p := newTestParser(tc.input)
			tn, parseErr7 := p.parseTypeName()
			if parseErr7 != nil {
				t.Fatalf("parse: %v", parseErr7)
			}
			if tn == nil {
				t.Fatal("expected non-nil TypeName")
			}
			if tn.TypeMods.Len() != tc.modCount {
				t.Errorf("expected %d type mods, got %d", tc.modCount, tn.TypeMods.Len())
			}
		})
	}
}
