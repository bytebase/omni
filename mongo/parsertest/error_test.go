package parsertest

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mongo/parser"
)

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty method", "db.users."},
		{"missing parens", "db.users.find"},
		{"unclosed brace", `db.users.find({name: "alice")`},
		{"unclosed bracket", `db.users.find([1, 2, 3)`},
		{"unknown top-level", "foobar"},
		{"invalid show target", "show something"},
		{"missing collection", "db..find()"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.Parse(tc.input)
			if err == nil {
				t.Fatalf("expected parse error for input: %s", tc.input)
			}
			var pe *parser.ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("expected *parser.ParseError, got %T: %v", err, err)
			}
			if pe.Line < 1 {
				t.Errorf("expected Line >= 1, got %d", pe.Line)
			}
			if pe.Column < 1 {
				t.Errorf("expected Column >= 1, got %d", pe.Column)
			}
			if pe.Message == "" {
				t.Error("expected non-empty error message")
			}
		})
	}
}

func TestNewKeywordErrorMessage(t *testing.T) {
	input := `db.users.find({_id: new ObjectId("abc")})`
	_, err := parser.Parse(input)
	if err == nil {
		t.Fatal("expected parse error for 'new' keyword")
	}
	var pe *parser.ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *parser.ParseError, got %T: %v", err, err)
	}
	if !strings.Contains(pe.Message, "new") {
		t.Errorf("expected error message to mention 'new', got: %s", pe.Message)
	}
	// "new" is at column 21 (0-indexed: 20)
	if pe.Line != 1 {
		t.Errorf("expected line 1, got %d", pe.Line)
	}
}

func TestErrorPosition(t *testing.T) {
	// The error should point to the problematic token.
	input := "show something"
	_, err := parser.Parse(input)
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *parser.ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *parser.ParseError, got %T", err)
	}
	// Error should be on line 1
	if pe.Line != 1 {
		t.Errorf("expected line 1, got %d", pe.Line)
	}
}

func TestErrorOnSecondLine(t *testing.T) {
	input := "show dbs\nfoobar"
	_, err := parser.Parse(input)
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *parser.ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *parser.ParseError, got %T", err)
	}
	if pe.Line != 2 {
		t.Errorf("expected line 2, got %d", pe.Line)
	}
	if pe.Column != 1 {
		t.Errorf("expected column 1, got %d", pe.Column)
	}
}

func TestEmptyInput(t *testing.T) {
	nodes, err := parser.Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(nodes))
	}
}

func TestOnlySemicolons(t *testing.T) {
	nodes, err := parser.Parse(";;;")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(nodes))
	}
}
