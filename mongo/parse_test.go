package mongo_test

import (
	"errors"
	"testing"

	mongo "github.com/bytebase/omni/mongo"
	"github.com/bytebase/omni/mongo/parser"
)

// TestParseStrict pins the strict contract of Parse: any statement that fails
// to parse fails the whole input, so callers never silently lose statements
// (BYT-9950: a dropped statement made a migration appear to skip a command).
func TestParseStrict(t *testing.T) {
	t.Run("mixed valid and invalid fails", func(t *testing.T) {
		input := "db.users.find();\nthis is not mongosh\ndb.orders.find();"
		_, err := mongo.Parse(input)
		if err == nil {
			t.Fatal("expected error for input containing an invalid statement")
		}
		var pe *parser.ParseError
		if !errors.As(err, &pe) {
			t.Fatalf("expected *parser.ParseError, got %T: %v", err, err)
		}
		if pe.Line != 2 {
			t.Errorf("expected error on line 2, got %d", pe.Line)
		}
	})

	t.Run("best effort recovers partial results", func(t *testing.T) {
		input := "db.users.find();\nthis is not mongosh\ndb.orders.find();"
		result := mongo.ParseBestEffort(input)
		if len(result.Statements) != 2 {
			t.Fatalf("expected 2 statements, got %d", len(result.Statements))
		}
		if len(result.Errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(result.Errors))
		}
	})

	t.Run("all valid succeeds", func(t *testing.T) {
		input := "db.users.find();\ndb.orders.countDocuments();"
		stmts, err := mongo.Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(stmts) != 2 {
			t.Fatalf("expected 2 statements, got %d", len(stmts))
		}
	})
}

// TestParseCreateIndexArithmetic pins the BYT-9950 statement: arithmetic
// expressions are not supported, and the strict Parse must surface the error
// instead of silently dropping the statement.
func TestParseCreateIndexArithmetic(t *testing.T) {
	input := `db.cs_customer_frequency.createIndex(
  { trans_date: 1 },
  { expireAfterSeconds: 90 * 24 * 60 * 60, name: "cs_customer_frequency_idx2" }
);`
	_, err := mongo.Parse(input)
	if err == nil {
		t.Fatal("expected parse error for arithmetic expression")
	}
	var pe *parser.ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *parser.ParseError, got %T: %v", err, err)
	}
}
