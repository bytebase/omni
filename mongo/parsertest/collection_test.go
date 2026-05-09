package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestBasicFind(t *testing.T) {
	node := mustParse(t, `db.users.find({name: "alice"})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Collection != "users" {
		t.Errorf("expected collection 'users', got %q", stmt.Collection)
	}
	if stmt.Method != "find" {
		t.Errorf("expected method 'find', got %q", stmt.Method)
	}
	if stmt.AccessMethod != "dot" {
		t.Errorf("expected access method 'dot', got %q", stmt.AccessMethod)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestFindWithProjection(t *testing.T) {
	node := mustParse(t, `db.users.find({}, {name: 1, _id: 0})`)
	stmt := node.(*ast.CollectionStatement)
	if len(stmt.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(stmt.Args))
	}
}

func TestFindWithCursorChain(t *testing.T) {
	node := mustParse(t, `db.users.find({}).sort({name: 1}).limit(10).skip(5)`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "find" {
		t.Errorf("expected 'find', got %q", stmt.Method)
	}
	if len(stmt.CursorMethods) != 3 {
		t.Fatalf("expected 3 cursor methods, got %d", len(stmt.CursorMethods))
	}
	if stmt.CursorMethods[0].Method != "sort" {
		t.Errorf("expected 'sort', got %q", stmt.CursorMethods[0].Method)
	}
	if stmt.CursorMethods[1].Method != "limit" {
		t.Errorf("expected 'limit', got %q", stmt.CursorMethods[1].Method)
	}
	if stmt.CursorMethods[2].Method != "skip" {
		t.Errorf("expected 'skip', got %q", stmt.CursorMethods[2].Method)
	}
}

func TestFindOne(t *testing.T) {
	node := mustParse(t, `db.users.findOne({_id: 1})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "findOne" {
		t.Errorf("expected 'findOne', got %q", stmt.Method)
	}
}

func TestInsertOne(t *testing.T) {
	node := mustParse(t, `db.users.insertOne({name: "bob", age: 25})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "insertOne" {
		t.Errorf("expected 'insertOne', got %q", stmt.Method)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestInsertMany(t *testing.T) {
	node := mustParse(t, `db.users.insertMany([{name: "a"}, {name: "b"}])`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "insertMany" {
		t.Errorf("expected 'insertMany', got %q", stmt.Method)
	}
}

func TestUpdateOne(t *testing.T) {
	node := mustParse(t, `db.users.updateOne({name: "alice"}, {$set: {age: 30}})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "updateOne" {
		t.Errorf("expected 'updateOne', got %q", stmt.Method)
	}
	if len(stmt.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(stmt.Args))
	}
}

func TestDeleteOne(t *testing.T) {
	node := mustParse(t, `db.users.deleteOne({name: "alice"})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "deleteOne" {
		t.Errorf("expected 'deleteOne', got %q", stmt.Method)
	}
}

func TestDeleteMany(t *testing.T) {
	node := mustParse(t, `db.users.deleteMany({status: "inactive"})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "deleteMany" {
		t.Errorf("expected 'deleteMany', got %q", stmt.Method)
	}
}

func TestAggregate(t *testing.T) {
	node := mustParse(t, `db.orders.aggregate([{$match: {status: "A"}}, {$group: {_id: "$cust_id", total: {$sum: "$amount"}}}])`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "aggregate" {
		t.Errorf("expected 'aggregate', got %q", stmt.Method)
	}
}

func TestCreateIndex(t *testing.T) {
	node := mustParse(t, `db.users.createIndex({name: 1})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "createIndex" {
		t.Errorf("expected 'createIndex', got %q", stmt.Method)
	}
}

func TestDropIndex(t *testing.T) {
	node := mustParse(t, `db.users.dropIndex("name_1")`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "dropIndex" {
		t.Errorf("expected 'dropIndex', got %q", stmt.Method)
	}
}

func TestGetIndexes(t *testing.T) {
	node := mustParse(t, `db.users.getIndexes()`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "getIndexes" {
		t.Errorf("expected 'getIndexes', got %q", stmt.Method)
	}
	if len(stmt.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(stmt.Args))
	}
}

func TestExplainPrefix(t *testing.T) {
	node := mustParse(t, `db.users.explain().find({name: "alice"})`)
	stmt := node.(*ast.CollectionStatement)
	if !stmt.Explain {
		t.Error("expected explain to be true")
	}
	if stmt.Method != "find" {
		t.Errorf("expected 'find', got %q", stmt.Method)
	}
}

func TestExplainWithArg(t *testing.T) {
	node := mustParse(t, `db.users.explain("executionStats").find({})`)
	stmt := node.(*ast.CollectionStatement)
	if !stmt.Explain {
		t.Error("expected explain to be true")
	}
	if len(stmt.ExplainArgs) != 1 {
		t.Fatalf("expected 1 explain arg, got %d", len(stmt.ExplainArgs))
	}
	s := stmt.ExplainArgs[0].(*ast.StringLiteral)
	if s.Value != "executionStats" {
		t.Errorf("expected 'executionStats', got %q", s.Value)
	}
	if stmt.Method != "find" {
		t.Errorf("expected 'find', got %q", stmt.Method)
	}
}

func TestGetCollectionAccess(t *testing.T) {
	node := mustParse(t, `db.getCollection("my-coll").find({})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Collection != "my-coll" {
		t.Errorf("expected 'my-coll', got %q", stmt.Collection)
	}
	if stmt.AccessMethod != "getCollection" {
		t.Errorf("expected 'getCollection', got %q", stmt.AccessMethod)
	}
	if stmt.Method != "find" {
		t.Errorf("expected 'find', got %q", stmt.Method)
	}
}

func TestBracketAccess(t *testing.T) {
	node := mustParse(t, `db["my-coll"].find({})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Collection != "my-coll" {
		t.Errorf("expected 'my-coll', got %q", stmt.Collection)
	}
	if stmt.AccessMethod != "bracket" {
		t.Errorf("expected 'bracket', got %q", stmt.AccessMethod)
	}
}

func TestBracketAccessSingleQuote(t *testing.T) {
	node := mustParse(t, `db['my-coll'].find({})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Collection != "my-coll" {
		t.Errorf("expected 'my-coll', got %q", stmt.Collection)
	}
}

func TestCountDocuments(t *testing.T) {
	node := mustParse(t, `db.users.countDocuments({active: true})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "countDocuments" {
		t.Errorf("expected 'countDocuments', got %q", stmt.Method)
	}
}

func TestDistinct(t *testing.T) {
	node := mustParse(t, `db.users.distinct("name", {active: true})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "distinct" {
		t.Errorf("expected 'distinct', got %q", stmt.Method)
	}
	if len(stmt.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(stmt.Args))
	}
}

func TestFindCursorToArray(t *testing.T) {
	node := mustParse(t, `db.users.find({}).toArray()`)
	stmt := node.(*ast.CollectionStatement)
	if len(stmt.CursorMethods) != 1 {
		t.Fatalf("expected 1 cursor method, got %d", len(stmt.CursorMethods))
	}
	if stmt.CursorMethods[0].Method != "toArray" {
		t.Errorf("expected 'toArray', got %q", stmt.CursorMethods[0].Method)
	}
}

func TestReplaceOne(t *testing.T) {
	node := mustParse(t, `db.users.replaceOne({_id: 1}, {name: "new"})`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "replaceOne" {
		t.Errorf("expected 'replaceOne', got %q", stmt.Method)
	}
}

func TestDrop(t *testing.T) {
	node := mustParse(t, `db.users.drop()`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "drop" {
		t.Errorf("expected 'drop', got %q", stmt.Method)
	}
}

func TestDatabaseMethod(t *testing.T) {
	node := mustParse(t, `db.getCollectionNames()`)
	stmt := node.(*ast.DatabaseStatement)
	if stmt.Method != "getCollectionNames" {
		t.Errorf("expected 'getCollectionNames', got %q", stmt.Method)
	}
}

func TestDatabaseDropDatabase(t *testing.T) {
	node := mustParse(t, `db.dropDatabase()`)
	stmt := node.(*ast.DatabaseStatement)
	if stmt.Method != "dropDatabase" {
		t.Errorf("expected 'dropDatabase', got %q", stmt.Method)
	}
}

func TestDatabaseRunCommand(t *testing.T) {
	node := mustParse(t, `db.runCommand({ping: 1})`)
	stmt := node.(*ast.DatabaseStatement)
	if stmt.Method != "runCommand" {
		t.Errorf("expected 'runCommand', got %q", stmt.Method)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestDatabaseCreateCollection(t *testing.T) {
	node := mustParse(t, `db.createCollection("test")`)
	stmt := node.(*ast.DatabaseStatement)
	if stmt.Method != "createCollection" {
		t.Errorf("expected 'createCollection', got %q", stmt.Method)
	}
}

func TestDatabaseGetSiblingDB(t *testing.T) {
	node := mustParse(t, `db.getSiblingDB("admin")`)
	stmt := node.(*ast.DatabaseStatement)
	if stmt.Method != "getSiblingDB" {
		t.Errorf("expected 'getSiblingDB', got %q", stmt.Method)
	}
}

func TestCollectionStatementLoc(t *testing.T) {
	node := mustParse(t, `db.c.find({})`)
	assertLoc(t, node, 0, 13)
}

func TestMultipleStatements(t *testing.T) {
	nodes := mustParseN(t, `db.a.find({}); db.b.insertOne({x: 1})`, 2)
	s1 := nodes[0].(*ast.CollectionStatement)
	s2 := nodes[1].(*ast.CollectionStatement)
	if s1.Collection != "a" || s1.Method != "find" {
		t.Errorf("statement 1: expected a.find, got %s.%s", s1.Collection, s1.Method)
	}
	if s2.Collection != "b" || s2.Method != "insertOne" {
		t.Errorf("statement 2: expected b.insertOne, got %s.%s", s2.Collection, s2.Method)
	}
}

func TestTrailingCommaInArguments(t *testing.T) {
	// Single trailing comma after the only argument — JS / mongosh accept
	// this and gomongoFallback events show it in the wild as
	// db.processed_files.find({...},).sort({...}).
	node := mustParse(t, `db.users.find({name: "alice"},)`)
	stmt := node.(*ast.CollectionStatement)
	if stmt.Method != "find" {
		t.Errorf("expected 'find', got %q", stmt.Method)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}

	// Trailing comma after the second argument.
	node = mustParse(t, `db.users.find({}, {name: 1},)`)
	stmt = node.(*ast.CollectionStatement)
	if len(stmt.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(stmt.Args))
	}

	// Trailing comma combined with a chained cursor method — exact shape
	// from the gomongoFallback events.
	node = mustParse(t, `db.processed_files.find({_id: "x"},).sort({created_at: -1})`)
	stmt = node.(*ast.CollectionStatement)
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
	if len(stmt.CursorMethods) != 1 || stmt.CursorMethods[0].Method != "sort" {
		t.Errorf("expected one .sort() cursor method, got %+v", stmt.CursorMethods)
	}

	// Cursor-method argument list also tolerates a trailing comma.
	node = mustParse(t, `db.users.find().limit(10,)`)
	stmt = node.(*ast.CollectionStatement)
	if len(stmt.CursorMethods) != 1 {
		t.Fatalf("expected 1 cursor method, got %d", len(stmt.CursorMethods))
	}
}

func TestArgumentListErrors(t *testing.T) {
	// A bare comma (no argument before it) is still an error — we only
	// accept ONE trailing comma after a real argument.
	mustFail(t, `db.users.find(,)`)
	// A leading comma is an error.
	mustFail(t, `db.users.find(, {a: 1})`)
	// Two commas in a row are an error.
	mustFail(t, `db.users.find({a: 1},,)`)
}
