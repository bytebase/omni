package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestEmptyDocument(t *testing.T) {
	node := mustParse(t, `db.c.find({})`)
	stmt := node.(*ast.CollectionStatement)
	doc := stmt.Args[0].(*ast.Document)
	if len(doc.Pairs) != 0 {
		t.Errorf("expected 0 pairs, got %d", len(doc.Pairs))
	}
}

func TestUnquotedKeys(t *testing.T) {
	node := mustParse(t, `db.c.find({name: "alice", age: 30})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	if len(doc.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(doc.Pairs))
	}
	if doc.Pairs[0].Key != "name" {
		t.Errorf("expected key 'name', got %q", doc.Pairs[0].Key)
	}
	if doc.Pairs[1].Key != "age" {
		t.Errorf("expected key 'age', got %q", doc.Pairs[1].Key)
	}
}

func TestQuotedKeys(t *testing.T) {
	node := mustParse(t, `db.c.find({"first name": "alice", 'last-name': "smith"})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	if len(doc.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(doc.Pairs))
	}
	if doc.Pairs[0].Key != "first name" {
		t.Errorf("expected key 'first name', got %q", doc.Pairs[0].Key)
	}
	if doc.Pairs[1].Key != "last-name" {
		t.Errorf("expected key 'last-name', got %q", doc.Pairs[1].Key)
	}
}

func TestDollarKeys(t *testing.T) {
	node := mustParse(t, `db.c.find({$gt: 5, $in: [1, 2]})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	if len(doc.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(doc.Pairs))
	}
	if doc.Pairs[0].Key != "$gt" {
		t.Errorf("expected key '$gt', got %q", doc.Pairs[0].Key)
	}
	if doc.Pairs[1].Key != "$in" {
		t.Errorf("expected key '$in', got %q", doc.Pairs[1].Key)
	}
}

func TestTrailingComma(t *testing.T) {
	node := mustParse(t, `db.c.find({a: 1, b: 2,})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	if len(doc.Pairs) != 2 {
		t.Errorf("expected 2 pairs, got %d", len(doc.Pairs))
	}
}

func TestNestedDocument(t *testing.T) {
	node := mustParse(t, `db.c.find({a: {b: {c: 1}}})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	if len(doc.Pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(doc.Pairs))
	}
	inner := doc.Pairs[0].Value.(*ast.Document)
	if len(inner.Pairs) != 1 {
		t.Fatalf("expected 1 inner pair, got %d", len(inner.Pairs))
	}
	deepest := inner.Pairs[0].Value.(*ast.Document)
	if len(deepest.Pairs) != 1 {
		t.Fatalf("expected 1 deepest pair, got %d", len(deepest.Pairs))
	}
	num := deepest.Pairs[0].Value.(*ast.NumberLiteral)
	if num.Value != "1" {
		t.Errorf("expected 1, got %s", num.Value)
	}
}

func TestKeywordAsKey(t *testing.T) {
	// Keywords like "find", "sort", "limit" can appear as document keys
	node := mustParse(t, `db.c.find({sort: 1, limit: 10, find: true})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	if len(doc.Pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(doc.Pairs))
	}
	if doc.Pairs[0].Key != "sort" {
		t.Errorf("expected key 'sort', got %q", doc.Pairs[0].Key)
	}
	if doc.Pairs[1].Key != "limit" {
		t.Errorf("expected key 'limit', got %q", doc.Pairs[1].Key)
	}
	if doc.Pairs[2].Key != "find" {
		t.Errorf("expected key 'find', got %q", doc.Pairs[2].Key)
	}
}

func TestMixedKeys(t *testing.T) {
	node := mustParse(t, `db.c.insertOne({_id: 1, "name": "test", $set: true})`)
	stmt := node.(*ast.CollectionStatement)
	doc := stmt.Args[0].(*ast.Document)
	if len(doc.Pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(doc.Pairs))
	}
	if doc.Pairs[0].Key != "_id" {
		t.Errorf("expected '_id', got %q", doc.Pairs[0].Key)
	}
	if doc.Pairs[1].Key != "name" {
		t.Errorf("expected 'name', got %q", doc.Pairs[1].Key)
	}
	if doc.Pairs[2].Key != "$set" {
		t.Errorf("expected '$set', got %q", doc.Pairs[2].Key)
	}
}

func TestDocumentWithArrayValue(t *testing.T) {
	node := mustParse(t, `db.c.find({tags: ["a", "b"], score: {$gt: 5}})`)
	doc := node.(*ast.CollectionStatement).Args[0].(*ast.Document)
	if len(doc.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(doc.Pairs))
	}
	arr := doc.Pairs[0].Value.(*ast.Array)
	if len(arr.Elements) != 2 {
		t.Errorf("expected 2 array elements, got %d", len(arr.Elements))
	}
	inner := doc.Pairs[1].Value.(*ast.Document)
	if inner.Pairs[0].Key != "$gt" {
		t.Errorf("expected '$gt', got %q", inner.Pairs[0].Key)
	}
}
