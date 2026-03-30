package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestStringLiteral(t *testing.T) {
	node := mustParse(t, `db.c.find({k: "hello"})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "find" {
		t.Fatalf("expected find, got %s", stmt.Method)
	}
	doc := stmt.Args[0].(*ast.Document)
	val := doc.Pairs[0].Value.(*ast.StringLiteral)
	if val.Value != "hello" {
		t.Errorf("expected 'hello', got %q", val.Value)
	}
}

func TestNumberLiteralInt(t *testing.T) {
	node := mustParse(t, `db.c.find({k: 42})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	val := doc.Pairs[0].Value.(*ast.NumberLiteral)
	if val.Value != "42" || val.IsFloat {
		t.Errorf("expected int 42, got %q isFloat=%v", val.Value, val.IsFloat)
	}
}

func TestNumberLiteralFloat(t *testing.T) {
	node := mustParse(t, `db.c.find({k: 3.14})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	val := doc.Pairs[0].Value.(*ast.NumberLiteral)
	if val.Value != "3.14" || !val.IsFloat {
		t.Errorf("expected float 3.14, got %q isFloat=%v", val.Value, val.IsFloat)
	}
}

func TestNumberLiteralExponent(t *testing.T) {
	node := mustParse(t, `db.c.find({k: 1e10})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	val := doc.Pairs[0].Value.(*ast.NumberLiteral)
	if val.Value != "1e10" || !val.IsFloat {
		t.Errorf("expected float 1e10, got %q isFloat=%v", val.Value, val.IsFloat)
	}
}

func TestNumberLiteralNegative(t *testing.T) {
	node := mustParse(t, `db.c.find({k: -7})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	val := doc.Pairs[0].Value.(*ast.NumberLiteral)
	if val.Value != "-7" || val.IsFloat {
		t.Errorf("expected int -7, got %q isFloat=%v", val.Value, val.IsFloat)
	}
}

func TestBoolLiterals(t *testing.T) {
	node := mustParse(t, `db.c.find({a: true, b: false})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	bTrue := doc.Pairs[0].Value.(*ast.BoolLiteral)
	bFalse := doc.Pairs[1].Value.(*ast.BoolLiteral)
	if !bTrue.Value {
		t.Error("expected true")
	}
	if bFalse.Value {
		t.Error("expected false")
	}
}

func TestNullLiteral(t *testing.T) {
	node := mustParse(t, `db.c.find({k: null})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	_, ok := doc.Pairs[0].Value.(*ast.NullLiteral)
	if !ok {
		t.Error("expected NullLiteral")
	}
}

func TestRegexLiteral(t *testing.T) {
	node := mustParse(t, `db.c.find({k: /abc/i})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	re := doc.Pairs[0].Value.(*ast.RegexLiteral)
	if re.Pattern != "abc" || re.Flags != "i" {
		t.Errorf("expected /abc/i, got /%s/%s", re.Pattern, re.Flags)
	}
}

func TestRegexLiteralNoFlags(t *testing.T) {
	node := mustParse(t, `db.c.find({k: /^test$/})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	re := doc.Pairs[0].Value.(*ast.RegexLiteral)
	if re.Pattern != "^test$" || re.Flags != "" {
		t.Errorf("expected /^test$/, got /%s/%s", re.Pattern, re.Flags)
	}
}

func TestHelperObjectId(t *testing.T) {
	node := mustParse(t, `db.c.find({_id: ObjectId("507f1f77bcf86cd799439011")})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	hc := doc.Pairs[0].Value.(*ast.HelperCall)
	if hc.Name != "ObjectId" {
		t.Errorf("expected ObjectId, got %s", hc.Name)
	}
	if len(hc.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(hc.Args))
	}
	s := hc.Args[0].(*ast.StringLiteral)
	if s.Value != "507f1f77bcf86cd799439011" {
		t.Errorf("wrong ObjectId arg: %q", s.Value)
	}
}

func TestHelperISODate(t *testing.T) {
	node := mustParse(t, `db.c.find({d: ISODate("2024-01-01")})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	hc := doc.Pairs[0].Value.(*ast.HelperCall)
	if hc.Name != "ISODate" {
		t.Errorf("expected ISODate, got %s", hc.Name)
	}
}

func TestHelperNumberLong(t *testing.T) {
	node := mustParse(t, `db.c.find({n: NumberLong(123)})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	hc := doc.Pairs[0].Value.(*ast.HelperCall)
	if hc.Name != "NumberLong" {
		t.Errorf("expected NumberLong, got %s", hc.Name)
	}
}

func TestHelperTimestamp(t *testing.T) {
	node := mustParse(t, `db.c.find({t: Timestamp(1, 0)})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	hc := doc.Pairs[0].Value.(*ast.HelperCall)
	if hc.Name != "Timestamp" || len(hc.Args) != 2 {
		t.Errorf("expected Timestamp with 2 args, got %s with %d args", hc.Name, len(hc.Args))
	}
}

func TestHelperUUID(t *testing.T) {
	node := mustParse(t, `db.c.find({u: UUID("abc")})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	hc := doc.Pairs[0].Value.(*ast.HelperCall)
	if hc.Name != "UUID" {
		t.Errorf("expected UUID, got %s", hc.Name)
	}
}

func TestHelperBinData(t *testing.T) {
	node := mustParse(t, `db.c.find({b: BinData(0, "data")})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	hc := doc.Pairs[0].Value.(*ast.HelperCall)
	if hc.Name != "BinData" || len(hc.Args) != 2 {
		t.Errorf("expected BinData with 2 args, got %s with %d", hc.Name, len(hc.Args))
	}
}

func TestHelperNoArgs(t *testing.T) {
	node := mustParse(t, `db.c.find({_id: ObjectId()})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	hc := doc.Pairs[0].Value.(*ast.HelperCall)
	if hc.Name != "ObjectId" || len(hc.Args) != 0 {
		t.Errorf("expected ObjectId() with 0 args, got %s with %d", hc.Name, len(hc.Args))
	}
}

func TestArrayLiteral(t *testing.T) {
	node := mustParse(t, `db.c.find({k: [1, "two", true]})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	arr := doc.Pairs[0].Value.(*ast.Array)
	if len(arr.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arr.Elements))
	}
	_ = arr.Elements[0].(*ast.NumberLiteral)
	_ = arr.Elements[1].(*ast.StringLiteral)
	_ = arr.Elements[2].(*ast.BoolLiteral)
}

func TestEmptyArray(t *testing.T) {
	node := mustParse(t, `db.c.find({k: []})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	arr := doc.Pairs[0].Value.(*ast.Array)
	if len(arr.Elements) != 0 {
		t.Errorf("expected 0 elements, got %d", len(arr.Elements))
	}
}

func TestArrayTrailingComma(t *testing.T) {
	node := mustParse(t, `db.c.find({k: [1, 2,]})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	arr := doc.Pairs[0].Value.(*ast.Array)
	if len(arr.Elements) != 2 {
		t.Errorf("expected 2 elements, got %d", len(arr.Elements))
	}
}

func TestNestedArray(t *testing.T) {
	node := mustParse(t, `db.c.find({k: [[1, 2], [3]]})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	arr := doc.Pairs[0].Value.(*ast.Array)
	if len(arr.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arr.Elements))
	}
	inner := arr.Elements[0].(*ast.Array)
	if len(inner.Elements) != 2 {
		t.Errorf("expected 2 inner elements, got %d", len(inner.Elements))
	}
}

func TestNewKeywordError(t *testing.T) {
	mustFail(t, `db.c.find({k: new Date()})`)
}
