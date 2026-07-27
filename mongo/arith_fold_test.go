package mongo_test

import (
	"strings"
	"testing"

	mongo "github.com/bytebase/omni/mongo"
	"github.com/bytebase/omni/mongo/ast"
)

// numberIn digs the value out of the first statement's first document field.
func foldedField(t *testing.T, js, field string) ast.Node {
	t.Helper()
	stmts, err := mongo.Parse(js)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	var found ast.Node
	var scan func(n ast.Node)
	scan = func(n ast.Node) {
		if found != nil || n == nil {
			return
		}
		switch v := n.(type) {
		case *ast.Document:
			for _, kv := range v.Pairs {
				if kv.Key == field {
					found = kv.Value
					return
				}
				scan(kv.Value)
			}
		case *ast.Array:
			for _, el := range v.Elements {
				scan(el)
			}
		case *ast.CollectionStatement:
			for _, a := range v.Args {
				scan(a)
			}
		}
	}
	scan(stmts[0].AST)
	if found == nil {
		t.Fatalf("field %s not found", field)
	}
	return found
}

// TestArithmeticConstantFolding pins constant-expression folding with
// JavaScript number semantics. Value AND text-form (BSON type carrier) are
// asserted: mongosh-verified on MongoDB 7 — integral results within int32
// arrive as BSON int, everything else as double; % follows JS remainder
// signs; division is float division.
func TestArithmeticConstantFolding(t *testing.T) {
	num := func(t *testing.T, js string) *ast.NumberLiteral {
		t.Helper()
		n, ok := foldedField(t, js, "x").(*ast.NumberLiteral)
		if !ok {
			t.Fatalf("expected folded NumberLiteral")
		}
		return n
	}

	for _, tc := range []struct {
		name    string
		js      string
		value   string
		isFloat bool
	}{
		{"customer ttl", `db.c.find({x: 90 * 24 * 60 * 60});`, "7776000", false},
		{"int division stays int form", `db.c.find({x: 6 / 2});`, "3", false},
		{"float division", `db.c.find({x: 6 / 4});`, "1.5", true},
		{"js remainder sign", `db.c.find({x: -7 % 3});`, "-1", false},
		{"float remainder", `db.c.find({x: 5.5 % 2});`, "1.5", true},
		{"parens precedence", `db.c.find({x: (1 + 2) * 3});`, "9", false},
		{"unary in expression", `db.c.find({x: 3 * -2});`, "-6", false},
		{"adjacent signed number is subtraction", `db.c.find({x: 5 -2});`, "3", false},
		{"int32 boundary stays int", `db.c.find({x: 2147483647 + 0});`, "2147483647", false},
		{"beyond int32 becomes double form", `db.c.find({x: 2147483647 + 1});`, "2.147483648e+09", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			n := num(t, tc.js)
			if n.Value != tc.value || n.IsFloat != tc.isFloat {
				t.Fatalf("folded to Value=%q IsFloat=%v, want Value=%q IsFloat=%v", n.Value, n.IsFloat, tc.value, tc.isFloat)
			}
		})
	}

	t.Run("string concatenation", func(t *testing.T) {
		s, ok := foldedField(t, `db.c.find({x: "a" + "b"});`, "x").(*ast.StringLiteral)
		if !ok || s.Value != "ab" {
			t.Fatalf("expected folded StringLiteral \"ab\", got %#v", s)
		}
	})

	t.Run("customer statement end to end", func(t *testing.T) {
		js := `db.cs_customer_frequency.createIndex({ trans_date: 1 }, { expireAfterSeconds: 90 * 24 * 60 * 60, name: "idx2" });`
		n, ok := foldedField(t, js, "expireAfterSeconds").(*ast.NumberLiteral)
		if !ok || n.Value != "7776000" || n.IsFloat {
			t.Fatalf("expireAfterSeconds folded wrong: %#v", n)
		}
	})
}

// TestArithmeticFoldingRejections pins the closed-whitelist boundary: only
// pure-literal subtrees fold; everything else fails, in the same shapes as
// before the change where applicable.
func TestArithmeticFoldingRejections(t *testing.T) {
	for _, tc := range []struct {
		name    string
		js      string
		wantErr string
	}{
		{"division by zero", `db.c.find({x: 1 / 0});`, "not a finite number"},
		{"zero by zero", `db.c.find({x: 0 / 0});`, "not a finite number"},
		{"number plus string", `db.c.find({x: 1 + "a"});`, "syntax error"},
		{"string times number", `db.c.find({x: "a" * 2});`, "syntax error"},
		{"identifier in arithmetic", `db.c.find({x: y * 2});`, "syntax error"},
		{"call in arithmetic", `db.c.find({x: ISODate("2024-01-01") + 1});`, "syntax error"},
		{"new still unsupported", `db.c.insertOne({d: new Date()});`, `"new" keyword is not supported`},
		{"method call form still rejected", `db.c.find({x: Date.now()});`, "expected ("},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := mongo.Parse(tc.js)
			if err == nil {
				t.Fatal("expected parse error")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErr)
			}
		})
	}
}
