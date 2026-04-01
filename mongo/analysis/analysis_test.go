package analysis_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/mongo"
	"github.com/bytebase/omni/mongo/analysis"
)

func mustAnalyze(t *testing.T, input string) *analysis.StatementAnalysis {
	t.Helper()
	stmts, err := mongo.Parse(input)
	if err != nil {
		t.Fatalf("Parse(%q): %v", input, err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Parse(%q): got %d statements, want 1", input, len(stmts))
	}
	return analysis.Analyze(stmts[0].AST)
}

func assertPredicateFields(t *testing.T, sa *analysis.StatementAnalysis, want []string) {
	t.Helper()
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	slices.Sort(want)
	if !slices.Equal(sa.PredicateFields, want) {
		t.Errorf("PredicateFields = %v, want %v", sa.PredicateFields, want)
	}
}

func TestOperationIsRead(t *testing.T) {
	tests := []struct {
		op   analysis.Operation
		want bool
	}{
		{analysis.OpFind, true},
		{analysis.OpFindOne, true},
		{analysis.OpAggregate, true},
		{analysis.OpCount, true},
		{analysis.OpDistinct, true},
		{analysis.OpRead, true},
		{analysis.OpWrite, false},
		{analysis.OpAdmin, false},
		{analysis.OpInfo, false},
		{analysis.OpExplain, false},
		{analysis.OpUnknown, false},
	}
	for _, tc := range tests {
		if got := tc.op.IsRead(); got != tc.want {
			t.Errorf("%v.IsRead() = %v, want %v", tc.op, got, tc.want)
		}
	}
}

func TestOperationString(t *testing.T) {
	tests := []struct {
		op   analysis.Operation
		want string
	}{
		{analysis.OpFind, "find"},
		{analysis.OpFindOne, "findOne"},
		{analysis.OpAggregate, "aggregate"},
		{analysis.OpCount, "count"},
		{analysis.OpDistinct, "distinct"},
		{analysis.OpRead, "read"},
		{analysis.OpWrite, "write"},
		{analysis.OpAdmin, "admin"},
		{analysis.OpInfo, "info"},
		{analysis.OpExplain, "explain"},
		{analysis.OpUnknown, "unknown"},
	}
	for _, tc := range tests {
		if got := tc.op.String(); got != tc.want {
			t.Errorf("%v.String() = %q, want %q", tc.op, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Collection statement tests
// ---------------------------------------------------------------------------

func TestAnalyzeCollectionFind(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find()`)
	if sa.Operation != analysis.OpFind {
		t.Errorf("Operation = %v, want OpFind", sa.Operation)
	}
	if sa.Collection != "users" {
		t.Errorf("Collection = %q, want %q", sa.Collection, "users")
	}
	if sa.MethodName != "find" {
		t.Errorf("MethodName = %q, want %q", sa.MethodName, "find")
	}
}

func TestAnalyzeCollectionFindOne(t *testing.T) {
	sa := mustAnalyze(t, `db.users.findOne()`)
	if sa.Operation != analysis.OpFindOne {
		t.Errorf("Operation = %v, want OpFindOne", sa.Operation)
	}
}

func TestAnalyzeCollectionAggregate(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate()`)
	if sa.Operation != analysis.OpAggregate {
		t.Errorf("Operation = %v, want OpAggregate", sa.Operation)
	}
}

func TestAnalyzeCollectionWrite(t *testing.T) {
	methods := []string{
		"insertOne", "insertMany", "updateOne", "updateMany",
		"deleteOne", "deleteMany", "replaceOne",
	}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			sa := mustAnalyze(t, "db.users."+method+"({})")
			if sa.Operation != analysis.OpWrite {
				t.Errorf("%s: Operation = %v, want OpWrite", method, sa.Operation)
			}
			if sa.Collection != "users" {
				t.Errorf("%s: Collection = %q, want %q", method, sa.Collection, "users")
			}
		})
	}
}

func TestAnalyzeCollectionAdmin(t *testing.T) {
	methods := []string{"createIndex", "dropIndex", "drop"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			sa := mustAnalyze(t, "db.users."+method+"({})")
			if sa.Operation != analysis.OpAdmin {
				t.Errorf("%s: Operation = %v, want OpAdmin", method, sa.Operation)
			}
		})
	}
}

func TestAnalyzeCollectionRead(t *testing.T) {
	methods := []string{"getIndexes", "stats", "storageSize"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			sa := mustAnalyze(t, "db.users."+method+"()")
			if sa.Operation != analysis.OpRead {
				t.Errorf("%s: Operation = %v, want OpRead", method, sa.Operation)
			}
		})
	}
}

func TestAnalyzeCollectionCount(t *testing.T) {
	methods := []string{"countDocuments", "estimatedDocumentCount", "count"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			sa := mustAnalyze(t, "db.users."+method+"()")
			if sa.Operation != analysis.OpCount {
				t.Errorf("%s: Operation = %v, want OpCount", method, sa.Operation)
			}
		})
	}
}

func TestAnalyzeCollectionDistinct(t *testing.T) {
	sa := mustAnalyze(t, `db.users.distinct("name")`)
	if sa.Operation != analysis.OpDistinct {
		t.Errorf("Operation = %v, want OpDistinct", sa.Operation)
	}
	if sa.MethodName != "distinct" {
		t.Errorf("MethodName = %q, want %q", sa.MethodName, "distinct")
	}
}

func TestAnalyzeCollectionExplain(t *testing.T) {
	sa := mustAnalyze(t, `db.users.explain().find()`)
	if sa.Operation != analysis.OpExplain {
		t.Errorf("Operation = %v, want OpExplain", sa.Operation)
	}
	if sa.Collection != "users" {
		t.Errorf("Collection = %q, want %q", sa.Collection, "users")
	}
}

func TestAnalyzeCollectionGetCollection(t *testing.T) {
	sa := mustAnalyze(t, `db.getCollection("users").find()`)
	if sa.Operation != analysis.OpFind {
		t.Errorf("Operation = %v, want OpFind", sa.Operation)
	}
	if sa.Collection != "users" {
		t.Errorf("Collection = %q, want %q", sa.Collection, "users")
	}
}

// ---------------------------------------------------------------------------
// Predicate field extraction tests
// ---------------------------------------------------------------------------

func TestPredicateSimpleField(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find({email: "alice@example.com"})`)
	assertPredicateFields(t, sa, []string{"email"})
}

func TestPredicateNestedDocument(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find({contact: {phone: "123"}})`)
	assertPredicateFields(t, sa, []string{"contact", "contact.phone"})
}

func TestPredicateDotPathKey(t *testing.T) {
	sa := mustAnalyze(t, `db.users.findOne({"contact.phone": "123"})`)
	assertPredicateFields(t, sa, []string{"contact.phone"})
}

func TestPredicateLogicalOperators(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find({
		$or: [{name: "alice"}, {email: "alice@example.com"}],
		$and: [{status: "active"}, {$nor: [{contact: {phone: "123"}}, {"profile.ssn": "456"}]}]
	})`)
	assertPredicateFields(t, sa, []string{"contact", "contact.phone", "email", "name", "profile.ssn", "status"})
}

func TestPredicateComparisonOperatorsSkipped(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find({age: {$gt: 18, $lt: 65}})`)
	assertPredicateFields(t, sa, []string{"age"})
}

func TestPredicateNoArgs(t *testing.T) {
	sa := mustAnalyze(t, `db.users.find()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.PredicateFields != nil {
		t.Errorf("PredicateFields = %v, want nil", sa.PredicateFields)
	}
}

func TestPredicateFindOneFields(t *testing.T) {
	sa := mustAnalyze(t, `db.users.findOne({name: "alice", age: 30})`)
	assertPredicateFields(t, sa, []string{"age", "name"})
}

// ---------------------------------------------------------------------------
// Pipeline analysis tests
// ---------------------------------------------------------------------------

func TestPipelineShapePreserving(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([
		{$match: {name: "alice"}},
		{$sort: {name: 1}},
		{$limit: 10}
	])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if !sa.ShapePreserving {
		t.Error("ShapePreserving = false, want true")
	}
	want := []string{"$match", "$sort", "$limit"}
	if !slices.Equal(sa.PipelineStages, want) {
		t.Errorf("PipelineStages = %v, want %v", sa.PipelineStages, want)
	}
	assertPredicateFields(t, sa, []string{"name"})
}

func TestPipelineUnsupportedGroup(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$group: {_id: "$name"}}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.ShapePreserving {
		t.Error("ShapePreserving = true, want false")
	}
	if sa.UnsupportedStage != "$group" {
		t.Errorf("UnsupportedStage = %q, want %q", sa.UnsupportedStage, "$group")
	}
}

func TestPipelineUnsupportedProject(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$match: {status: "active"}}, {$project: {name: 1}}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.UnsupportedStage != "$project" {
		t.Errorf("UnsupportedStage = %q, want %q", sa.UnsupportedStage, "$project")
	}
}

func TestPipelineEmpty(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if !sa.ShapePreserving {
		t.Error("ShapePreserving = false, want true")
	}
	if sa.UnsupportedStage != "" {
		t.Errorf("UnsupportedStage = %q, want empty", sa.UnsupportedStage)
	}
}

func TestPipelineNoArgs(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate()`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.Operation != analysis.OpAggregate {
		t.Errorf("Operation = %v, want OpAggregate", sa.Operation)
	}
}

func TestPipelineLookupSimple(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([
		{$lookup: {from: "orders", localField: "id", foreignField: "userId", as: "orders"}}
	])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.UnsupportedStage != "" {
		t.Errorf("UnsupportedStage = %q, want empty", sa.UnsupportedStage)
	}
	if len(sa.Joins) != 1 {
		t.Fatalf("len(Joins) = %d, want 1", len(sa.Joins))
	}
	if sa.Joins[0].Collection != "orders" {
		t.Errorf("Joins[0].Collection = %q, want %q", sa.Joins[0].Collection, "orders")
	}
	if sa.Joins[0].AsField != "orders" {
		t.Errorf("Joins[0].AsField = %q, want %q", sa.Joins[0].AsField, "orders")
	}
}

func TestPipelineLookupPipelineForm(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([
		{$lookup: {from: "orders", let: {uid: "$id"}, pipeline: [{$match: {}}], as: "orders"}}
	])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if sa.UnsupportedStage != "$lookup" {
		t.Errorf("UnsupportedStage = %q, want %q", sa.UnsupportedStage, "$lookup")
	}
}

func TestPipelineGraphLookup(t *testing.T) {
	sa := mustAnalyze(t, `db.employees.aggregate([
		{$graphLookup: {from: "employees", startWith: "$managerId", connectFromField: "managerId", connectToField: "_id", as: "hierarchy"}}
	])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if len(sa.Joins) != 1 {
		t.Fatalf("len(Joins) = %d, want 1", len(sa.Joins))
	}
	if sa.Joins[0].Collection != "employees" {
		t.Errorf("Joins[0].Collection = %q, want %q", sa.Joins[0].Collection, "employees")
	}
	if sa.Joins[0].AsField != "hierarchy" {
		t.Errorf("Joins[0].AsField = %q, want %q", sa.Joins[0].AsField, "hierarchy")
	}
}

func TestPipelineUnwind(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$unwind: "$tags"}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if !sa.ShapePreserving {
		t.Error("ShapePreserving = false, want true")
	}
}

func TestPipelineMatchWithUnwind(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([{$match: {status: "active"}}, {$unwind: "$tags"}])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	if !sa.ShapePreserving {
		t.Error("ShapePreserving = false, want true")
	}
	assertPredicateFields(t, sa, []string{"status"})
}

func TestPipelineMatchLogicalOps(t *testing.T) {
	sa := mustAnalyze(t, `db.users.aggregate([
		{$match: {$or: [{name: "alice"}, {age: 30}]}}
	])`)
	if sa == nil {
		t.Fatal("Analyze returned nil")
	}
	assertPredicateFields(t, sa, []string{"age", "name"})
}

// ---------------------------------------------------------------------------
// Non-collection statement tests
// ---------------------------------------------------------------------------

func TestAnalyzeShowDbs(t *testing.T) {
	sa := mustAnalyze(t, `show dbs`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("got %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeShowCollections(t *testing.T) {
	sa := mustAnalyze(t, `show collections`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("got %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeDatabaseDropDatabase(t *testing.T) {
	sa := mustAnalyze(t, `db.dropDatabase()`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpAdmin {
		t.Errorf("got %v, want OpAdmin", sa.Operation)
	}
}

func TestAnalyzeDatabaseCreateCollection(t *testing.T) {
	sa := mustAnalyze(t, `db.createCollection("test")`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpAdmin {
		t.Errorf("got %v, want OpAdmin", sa.Operation)
	}
}

func TestAnalyzeDatabaseGetCollectionNames(t *testing.T) {
	sa := mustAnalyze(t, `db.getCollectionNames()`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("got %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeDatabaseServerStatus(t *testing.T) {
	sa := mustAnalyze(t, `db.serverStatus()`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("got %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeRsStatus(t *testing.T) {
	sa := mustAnalyze(t, `rs.status()`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("got %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeRsInitiate(t *testing.T) {
	sa := mustAnalyze(t, `rs.initiate()`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("got %v, want OpWrite", sa.Operation)
	}
}

func TestAnalyzeShStatus(t *testing.T) {
	sa := mustAnalyze(t, `sh.status()`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("got %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeShAddShard(t *testing.T) {
	sa := mustAnalyze(t, `sh.addShard("shard1/host:port")`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("got %v, want OpWrite", sa.Operation)
	}
}

func TestAnalyzeBulkStatement(t *testing.T) {
	sa := mustAnalyze(t, `db.users.initializeOrderedBulkOp().insert({name: "alice"}).execute()`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("got %v, want OpWrite", sa.Operation)
	}
	if sa.Collection != "users" {
		t.Errorf("Collection = %q, want %q", sa.Collection, "users")
	}
}

func TestAnalyzeNativeFunction(t *testing.T) {
	sa := mustAnalyze(t, `sleep(1000)`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("got %v, want OpWrite", sa.Operation)
	}
}

func TestAnalyzeNil(t *testing.T) {
	sa := analysis.Analyze(nil)
	if sa != nil {
		t.Errorf("got %v, want nil", sa)
	}
}

func TestAnalyzeRunCommandFind(t *testing.T) {
	sa := mustAnalyze(t, `db.runCommand({find: "users", filter: {status: "active"}})`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpRead {
		t.Errorf("got %v, want OpRead", sa.Operation)
	}
}

func TestAnalyzeRunCommandCreate(t *testing.T) {
	sa := mustAnalyze(t, `db.runCommand({create: "newCollection"})`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpAdmin {
		t.Errorf("got %v, want OpAdmin", sa.Operation)
	}
}

func TestAnalyzeRunCommandInsert(t *testing.T) {
	sa := mustAnalyze(t, `db.runCommand({insert: "users", documents: [{name: "alice"}]})`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpWrite {
		t.Errorf("got %v, want OpWrite", sa.Operation)
	}
}

func TestAnalyzeRunCommandServerStatus(t *testing.T) {
	sa := mustAnalyze(t, `db.runCommand({serverStatus: 1})`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("got %v, want OpInfo", sa.Operation)
	}
}

func TestAnalyzeAdminCommand(t *testing.T) {
	sa := mustAnalyze(t, `db.adminCommand({listDatabases: 1})`)
	if sa == nil {
		t.Fatal("nil")
	}
	if sa.Operation != analysis.OpInfo {
		t.Errorf("got %v, want OpInfo", sa.Operation)
	}
}
