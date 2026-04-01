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

// ---------------------------------------------------------------------------
// Masking parity integration tests
// ---------------------------------------------------------------------------

func TestAnalyzeMaskingParity(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		wantOp          analysis.Operation
		wantMethod      string
		wantCollection  string
		wantPredicates  []string
		wantShapePres   bool
		wantUnsupported string
		wantJoins       []analysis.JoinInfo
		wantNil         bool
	}{
		{
			name:           "find with predicate",
			input:          `db.users.find({email: "alice@example.com"})`,
			wantOp:         analysis.OpFind,
			wantMethod:     "find",
			wantCollection: "users",
			wantPredicates: []string{"email"},
		},
		{
			name:           "findOne with dot path key",
			input:          `db.users.findOne({"contact.phone": "123"})`,
			wantOp:         analysis.OpFindOne,
			wantMethod:     "findOne",
			wantCollection: "users",
			wantPredicates: []string{"contact.phone"},
		},
		{
			name:           "getCollection access",
			input:          `db.getCollection("users").find({email: "x"})`,
			wantOp:         analysis.OpFind,
			wantMethod:     "find",
			wantCollection: "users",
			wantPredicates: []string{"email"},
		},
		{
			name: "find with logical operators",
			input: `db.getCollection("users").find({
				$or: [
					{email: "a@example.com"},
					{contact: {phone: "123"}}
				],
				$and: [
					{"profile.ssn": "111"},
					{name: "alice"}
				],
				$nor: [{status: "inactive"}]
			})`,
			wantOp:         analysis.OpFind,
			wantMethod:     "find",
			wantCollection: "users",
			wantPredicates: []string{"contact", "contact.phone", "email", "name", "profile.ssn", "status"},
		},
		{
			name:           "aggregate shape-preserving",
			input:          `db.users.aggregate([{$match: {name: "alice"}}])`,
			wantOp:         analysis.OpAggregate,
			wantMethod:     "aggregate",
			wantCollection: "users",
			wantPredicates: []string{"name"},
			wantShapePres:  true,
		},
		{
			name:           "aggregate multiple shape-preserving stages",
			input:          `db.users.aggregate([{$match: {status: "active"}}, {$sort: {name: 1}}, {$limit: 10}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantPredicates: []string{"status"},
			wantShapePres:  true,
		},
		{
			name:           "aggregate match with logical operators",
			input:          `db.users.aggregate([{$match: {$or: [{age: {$gt: 18}}, {name: "alice"}]}}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantPredicates: []string{"age", "name"},
			wantShapePres:  true,
		},
		{
			name:           "aggregate addFields and unset",
			input:          `db.users.aggregate([{$addFields: {fullName: "test"}}, {$unset: "ssn"}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantShapePres:  true,
		},
		{
			name:           "aggregate empty pipeline",
			input:          `db.users.aggregate([])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantShapePres:  true,
		},
		{
			name:           "aggregate no arguments",
			input:          `db.users.aggregate()`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantShapePres:  true,
		},
		{
			name:            "aggregate $group unsupported",
			input:           `db.users.aggregate([{$group: {_id: "$status"}}])`,
			wantOp:          analysis.OpAggregate,
			wantCollection:  "users",
			wantUnsupported: "$group",
		},
		{
			name:            "aggregate $project unsupported",
			input:           `db.users.aggregate([{$match: {name: "alice"}}, {$project: {name: 1}}])`,
			wantOp:          analysis.OpAggregate,
			wantCollection:  "users",
			wantUnsupported: "$project",
		},
		{
			name:           "aggregate $lookup simple",
			input:          `db.users.aggregate([{$lookup: {from: "orders", localField: "_id", foreignField: "userId", as: "orders"}}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantShapePres:  true,
			wantJoins:      []analysis.JoinInfo{{Collection: "orders", AsField: "orders"}},
		},
		{
			name:            "aggregate $lookup pipeline form unsupported",
			input:           `db.users.aggregate([{$lookup: {from: "orders", pipeline: [{$match: {status: "active"}}], as: "orders"}}])`,
			wantOp:          analysis.OpAggregate,
			wantCollection:  "users",
			wantUnsupported: "$lookup",
		},
		{
			name:           "aggregate $graphLookup",
			input:          `db.users.aggregate([{$graphLookup: {from: "employees", startWith: "$reportsTo", connectFromField: "reportsTo", connectToField: "name", as: "reportingHierarchy"}}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantShapePres:  true,
			wantJoins:      []analysis.JoinInfo{{Collection: "employees", AsField: "reportingHierarchy"}},
		},
		{
			name:            "aggregate $out unsupported",
			input:           `db.users.aggregate([{$match: {status: "active"}}, {$out: "activeUsers"}])`,
			wantOp:          analysis.OpAggregate,
			wantCollection:  "users",
			wantUnsupported: "$out",
		},
		{
			name:           "aggregate $unwind shape-preserving",
			input:          `db.users.aggregate([{$unwind: "$tags"}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantShapePres:  true,
		},
		{
			name:           "aggregate $match and $unwind",
			input:          `db.users.aggregate([{$match: {status: "active"}}, {$unwind: "$tags"}])`,
			wantOp:         analysis.OpAggregate,
			wantCollection: "users",
			wantPredicates: []string{"status"},
			wantShapePres:  true,
		},
		{
			name:            "aggregate $replaceRoot unsupported",
			input:           `db.users.aggregate([{$replaceRoot: {newRoot: "$contact"}}])`,
			wantOp:          analysis.OpAggregate,
			wantCollection:  "users",
			wantUnsupported: "$replaceRoot",
		},
		{
			name:            "aggregate $count unsupported",
			input:           `db.users.aggregate([{$count: "total"}])`,
			wantOp:          analysis.OpAggregate,
			wantCollection:  "users",
			wantUnsupported: "$count",
		},
		{
			name:           "countDocuments unsupported read",
			input:          `db.users.countDocuments({})`,
			wantOp:         analysis.OpCount,
			wantMethod:     "countDocuments",
			wantCollection: "users",
		},
		{
			name:           "distinct unsupported read",
			input:          `db.users.distinct("name")`,
			wantOp:         analysis.OpDistinct,
			wantMethod:     "distinct",
			wantCollection: "users",
		},
		{
			name:           "write method returns non-nil",
			input:          `db.users.insertOne({name: "alice"})`,
			wantOp:         analysis.OpWrite,
			wantCollection: "users",
		},
		{
			name:   "show dbs",
			input:  `show dbs`,
			wantOp: analysis.OpInfo,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sa := mustAnalyze(t, tc.input)
			if tc.wantNil {
				if sa != nil {
					t.Fatalf("got %+v, want nil", sa)
				}
				return
			}
			if sa == nil {
				t.Fatal("Analyze returned nil")
			}
			if sa.Operation != tc.wantOp {
				t.Errorf("Operation = %v, want %v", sa.Operation, tc.wantOp)
			}
			if tc.wantMethod != "" && sa.MethodName != tc.wantMethod {
				t.Errorf("MethodName = %q, want %q", sa.MethodName, tc.wantMethod)
			}
			if tc.wantCollection != "" && sa.Collection != tc.wantCollection {
				t.Errorf("Collection = %q, want %q", sa.Collection, tc.wantCollection)
			}
			if tc.wantPredicates != nil {
				slices.Sort(tc.wantPredicates)
				if !slices.Equal(sa.PredicateFields, tc.wantPredicates) {
					t.Errorf("PredicateFields = %v, want %v", sa.PredicateFields, tc.wantPredicates)
				}
			}
			if sa.ShapePreserving != tc.wantShapePres {
				t.Errorf("ShapePreserving = %v, want %v", sa.ShapePreserving, tc.wantShapePres)
			}
			if sa.UnsupportedStage != tc.wantUnsupported {
				t.Errorf("UnsupportedStage = %q, want %q", sa.UnsupportedStage, tc.wantUnsupported)
			}
			if tc.wantJoins != nil {
				if !slices.Equal(sa.Joins, tc.wantJoins) {
					t.Errorf("Joins = %v, want %v", sa.Joins, tc.wantJoins)
				}
			}
		})
	}
}
