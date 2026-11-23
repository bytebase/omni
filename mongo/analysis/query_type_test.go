package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mongo/parser"
)

func TestClassifyStatement(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      QueryType
	}{
		// Select commands.
		{name: "find", statement: `db.users.find({})`, want: Select},
		{name: "findOne", statement: `db.users.findOne({name: "John"})`, want: Select},
		{name: "count", statement: `db.users.count({})`, want: Select},
		{name: "countDocuments", statement: `db.users.countDocuments({})`, want: Select},
		{name: "estimatedDocumentCount", statement: `db.users.estimatedDocumentCount()`, want: Select},
		{name: "distinct", statement: `db.users.distinct("name")`, want: Select},
		{name: "aggregate without $out/$merge", statement: `db.users.aggregate([{$match: {}}])`, want: Select},
		{name: "aggregate with $project read-only", statement: `db.users.aggregate([{$project: {name: 1}}])`, want: Select},
		{name: "find with trailing semicolon", statement: `db.users.find({});`, want: Select},

		// DML commands.
		{name: "insertOne", statement: `db.users.insertOne({name: "John"})`, want: DML},
		{name: "insertMany", statement: `db.users.insertMany([{a: 1}, {a: 2}])`, want: DML},
		{name: "updateOne", statement: `db.users.updateOne({}, {$set: {a: 1}})`, want: DML},
		{name: "updateMany", statement: `db.users.updateMany({}, {$set: {a: 1}})`, want: DML},
		{name: "deleteOne", statement: `db.users.deleteOne({})`, want: DML},
		{name: "deleteMany", statement: `db.users.deleteMany({})`, want: DML},
		{name: "replaceOne", statement: `db.users.replaceOne({a: 1}, {a: 2})`, want: DML},
		{name: "findOneAndUpdate", statement: `db.users.findOneAndUpdate({}, {$set: {a: 1}})`, want: DML},
		{name: "findOneAndReplace", statement: `db.users.findOneAndReplace({a: 1}, {a: 2})`, want: DML},
		{name: "findOneAndDelete", statement: `db.users.findOneAndDelete({a: 1})`, want: DML},
		{name: "bulkWrite", statement: `db.users.bulkWrite([{insertOne: {document: {a: 1}}}])`, want: DML},

		// Aggregate with $out / $merge (including after non-shape-preserving stages) and dynamic pipelines.
		{name: "aggregate with $out", statement: `db.users.aggregate([{$match: {}}, {$out: "output"}])`, want: DML},
		{name: "aggregate with $merge", statement: `db.users.aggregate([{$match: {}}, {$merge: {into: "output"}}])`, want: DML},
		{name: "aggregate with $project then $out", statement: `db.users.aggregate([{$project: {name: 1}}, {$out: "output"}])`, want: DML},
		{name: "aggregate with $group then $merge", statement: `db.users.aggregate([{$group: {_id: "$k"}}, {$merge: {into: "output"}}])`, want: DML},
		{name: "aggregate with pipeline $lookup then $out", statement: `db.users.aggregate([{$lookup: {from: "orders", pipeline: [{$match: {}}], as: "o"}}, {$out: "output"}])`, want: DML},
		{name: "aggregate with dynamic pipeline variable", statement: `db.users.aggregate(pipelineVar)`, want: DML},
		{name: "aggregate with dynamic stage after $project", statement: `db.users.aggregate([{$project: {name: 1}}, dynamicStage])`, want: DML},
		{name: "aggregate with multi-key stage doc", statement: `db.users.aggregate([{$match: {}, $out: "output"}])`, want: DML},
		{name: "aggregate with identity-escaped $out", statement: `db.users.aggregate([{"$\out": "output"}])`, want: DML},
		{name: "aggregate with braced unicode $out", statement: `db.users.aggregate([{"$\u{6f}ut": "output"}])`, want: DML},
		{name: "aggregate with octal-escaped $out", statement: `db.users.aggregate([{"$\157ut": "output"}])`, want: DML},
		{name: "aggregate with line-continuation $out", statement: "db.users.aggregate([{\"$\\\nout\": \"output\"}])", want: DML},
		{name: "aggregate with identity-escaped $merge", statement: `db.users.aggregate([{"$\merge": {into: "output"}}])`, want: DML},

		// DDL commands.
		{name: "createCollection", statement: `db.createCollection("test")`, want: DDL},
		{name: "createView", statement: `db.createView("v", "users", [])`, want: DDL},
		{name: "dropDatabase", statement: `db.dropDatabase()`, want: DDL},
		{name: "drop collection", statement: `db.users.drop()`, want: DDL},
		{name: "createIndex", statement: `db.users.createIndex({name: 1})`, want: DDL},
		{name: "createIndexes", statement: `db.users.createIndexes([{key: {name: 1}}])`, want: DDL},
		{name: "dropIndex", statement: `db.users.dropIndex("name_1")`, want: DDL},
		{name: "dropIndexes", statement: `db.users.dropIndexes()`, want: DDL},
		{name: "renameCollection", statement: `db.users.renameCollection("newName")`, want: DDL},
		{name: "reIndex", statement: `db.users.reIndex()`, want: DDL},

		// Info commands (SelectInfoSchema).
		{name: "show dbs", statement: `show dbs`, want: SelectInfoSchema},
		{name: "show collections", statement: `show collections`, want: SelectInfoSchema},
		{name: "getCollectionNames", statement: `db.getCollectionNames()`, want: SelectInfoSchema},
		{name: "getCollectionInfos", statement: `db.getCollectionInfos()`, want: SelectInfoSchema},
		{name: "db.stats()", statement: `db.stats()`, want: SelectInfoSchema},
		{name: "db.getName()", statement: `db.getName()`, want: SelectInfoSchema},
		{name: "collection stats", statement: `db.users.stats()`, want: SelectInfoSchema},
		{name: "serverStatus", statement: `db.serverStatus()`, want: SelectInfoSchema},
		{name: "serverBuildInfo", statement: `db.serverBuildInfo()`, want: SelectInfoSchema},
		{name: "db.version()", statement: `db.version()`, want: SelectInfoSchema},
		{name: "hostInfo", statement: `db.hostInfo()`, want: SelectInfoSchema},
		{name: "listCommands", statement: `db.listCommands()`, want: SelectInfoSchema},
		{name: "getIndexes", statement: `db.users.getIndexes()`, want: SelectInfoSchema},
		{name: "dataSize", statement: `db.users.dataSize()`, want: SelectInfoSchema},
		{name: "storageSize", statement: `db.users.storageSize()`, want: SelectInfoSchema},
		{name: "totalIndexSize", statement: `db.users.totalIndexSize()`, want: SelectInfoSchema},
		{name: "totalSize", statement: `db.users.totalSize()`, want: SelectInfoSchema},
		{name: "isCapped", statement: `db.users.isCapped()`, want: DML},
		{name: "validate", statement: `db.users.validate()`, want: SelectInfoSchema},
		{name: "validate full true", statement: `db.users.validate({full: true})`, want: SelectInfoSchema},
		{name: "validate legacy bool", statement: `db.users.validate(true)`, want: SelectInfoSchema},
		{name: "validate repair false", statement: `db.logs.validate({repair: false})`, want: SelectInfoSchema},
		{name: "validate repair true", statement: `db.logs.validate({repair: true})`, want: DML},
		{name: "validate identity-escaped repair true", statement: `db.logs.validate({"rep\air": true})`, want: DML},
		{name: "validate octal-escaped repair true", statement: `db.logs.validate({"rep\141ir": true})`, want: DML},
		{name: "validate repair truthy number", statement: `db.logs.validate({repair: 1})`, want: DML},
		{name: "validate repair variable", statement: `db.logs.validate({repair: doRepair})`, want: DML},
		{name: "validate options variable", statement: `db.logs.validate(opts)`, want: DML},
		{name: "validate extra argument", statement: `db.logs.validate({full: true}, extra)`, want: DML},
		{name: "latencyStats", statement: `db.users.latencyStats()`, want: SelectInfoSchema},
		{name: "getShardDistribution", statement: `db.users.getShardDistribution()`, want: SelectInfoSchema},

		// Safe BSON constructors vs mutating/unrecognized calls inside arguments.
		{name: "find with ObjectId and ISODate", statement: `db.users.find({_id: ObjectId("507f1f77bcf86cd799439011"), ts: ISODate("2024-01-01")})`, want: Select},
		{name: "find with static BSON factory", statement: `db.users.find({bin: Binary.createFromBase64("AA=="), oid: ObjectId.createFromHexString("507f1f77bcf86cd799439011")})`, want: Select},
		{name: "find with EJSON.stringify safe forms", statement: `db.users.find({a: EJSON.stringify({k: 1}), b: EJSON.stringify({k: 1}, null, 2), c: EJSON.stringify({k: 1}, ["k"]), d: EJSON.stringify({k: 1}, {relaxed: true})})`, want: Select},
		{name: "find with safe JS primitive constants", statement: `db.users.find({a: undefined, b: NaN, c: Infinity})`, want: Select},
		{name: "find with EJSON.stringify replacer callback", statement: `db.users.find({x: EJSON.stringify({a: 1}, mutatingCallback)})`, want: DML},
		{name: "find with EJSON.stringify toJSON callback", statement: `db.users.find({x: EJSON.stringify({toJSON: mutatingCallback})})`, want: DML},
		{name: "find with Date.parse toString callback", statement: `db.users.find({x: Date.parse({toString: mutatingCallback})})`, want: DML},
		{name: "find with RegExp toString callback", statement: `db.users.find({x: RegExp({toString: mutatingCallback})})`, want: DML},
		{name: "find with toBSON callback", statement: `db.users.find({x: {toBSON: mutatingCallback}})`, want: DML},
		{name: "runCommand with escaped aggregate and $out", statement: `db.runCommand({"agg\regate": "users", "pip\eline": [{"$\out": "output"}], cursor: {}})`, want: DML},
		{name: "find with nested deleteMany in filter", statement: `db.users.find({deleted: db.getCollection("audit").deleteMany({})})`, want: DML},
		{name: "find with nested mutation inside ObjectId arg", statement: `db.users.find({_id: ObjectId(db.dropDatabase())})`, want: DML},
		{name: "find with mutation hidden before static factory", statement: `db.users.find({_id: ObjectId(db.dropDatabase()).createFromHexString("507f1f77bcf86cd799439011")})`, want: DML},
		{name: "find with unrecognized helper call", statement: `db.users.find({x: customHelper()})`, want: DML},
		{name: "find with mutation in cursor modifier", statement: `db.users.find({}).sort({a: db.getCollection("users").drop()})`, want: DML},
		{name: "explain with mutation in filter", statement: `db.users.explain().find({x: db.getCollection("audit").deleteMany({})})`, want: DML},
		{name: "explain suffix with mutation in filter", statement: `db.users.find({x: db.getCollection("audit").deleteMany({})}).explain()`, want: DML},
		{name: "aggregate with mutation in stage", statement: `db.users.aggregate([{$match: {x: db.getCollection("audit").deleteMany({})}}])`, want: DML},
		{name: "runCommand find with mutation in filter", statement: `db.runCommand({find: "users", filter: {x: db.getCollection("audit").deleteMany({})}})`, want: DML},
		{name: "db.stats with mutation in arg", statement: `db.stats({scale: db.dropDatabase()})`, want: DML},
		{name: "db.stats with chained call", statement: `db.stats().drop()`, want: DML},
		{name: "rs.status with mutation in arg", statement: `rs.status(db.dropDatabase())`, want: DML},
		{name: "sh.status with mutation in arg", statement: `sh.status(db.dropDatabase())`, want: DML},

		// Safe cursor chains vs callback-taking or post-terminal chained methods.
		{name: "find with safe cursor chain and toArray", statement: `db.users.find({}).sort({a: 1}).limit(10).toArray()`, want: Select},
		{name: "aggregate with safe toArray terminal", statement: `db.users.aggregate([{$match: {}}]).toArray()`, want: Select},
		{name: "distinct with explain suffix", statement: `db.users.distinct("name").explain()`, want: Explain},
		{name: "find with forEach callback", statement: `db.users.find({}).forEach(mutatingCallback)`, want: DML},
		{name: "find with map callback", statement: `db.users.find({}).map(mutatingCallback)`, want: DML},
		{name: "find with toArray then array sort callback", statement: `db.users.find({}).toArray().sort(eval)`, want: DML},
		{name: "find with toArray then forEach", statement: `db.users.find({}).toArray().forEach(eval)`, want: DML},
		{name: "distinct with forEach eval", statement: `db.users.distinct("script").forEach(eval)`, want: DML},
		{name: "distinct with array sort callback", statement: `db.users.distinct("script").sort(eval)`, want: DML},
		{name: "aggregate with forEach callback", statement: `db.users.aggregate([{$match: {}}]).forEach(eval)`, want: DML},
		{name: "explain prefix find with forEach callback", statement: `db.users.explain().find({}).forEach(eval)`, want: DML},
		{name: "find with explain then chained call", statement: `db.users.find({}).explain().forEach(eval)`, want: DML},
		{name: "findOne with chained method", statement: `db.users.findOne({}).forEach(eval)`, want: DML},

		// Explain.
		{name: "find with explain suffix", statement: `db.users.find({}).explain()`, want: Explain},
		{name: "find with sort and explain suffix", statement: `db.users.find({}).sort({a: 1}).explain()`, want: Explain},
		{name: "explain prefix with find", statement: `db.users.explain().find({})`, want: Explain},
		{name: "aggregate read-only with explain suffix", statement: `db.users.aggregate([{$match: {}}]).explain()`, want: Explain},
		{name: "explain prefix with read-only aggregate", statement: `db.users.explain().aggregate([{$match: {}}])`, want: Explain},
		{name: "insertOne with explain suffix retains DML", statement: `db.users.insertOne({name: "x"}).explain()`, want: DML},
		{name: "drop with explain suffix retains DDL", statement: `db.users.drop().explain()`, want: DDL},
		{name: "aggregate $out with explain suffix retains DML", statement: `db.users.aggregate([{$out: "output"}]).explain()`, want: DML},
		{name: "explain prefix with deleteOne retains DML", statement: `db.users.explain().deleteOne({})`, want: DML},
		{name: "explain prefix with drop retains DDL", statement: `db.users.explain().drop()`, want: DDL},
		{name: "explain prefix with aggregate $merge retains DML", statement: `db.users.explain().aggregate([{$merge: {into: "output"}}])`, want: DML},
		{name: "unknown method with explain suffix retains DML", statement: `db.users.someUnknownMethod().explain()`, want: DML},

		// Bulk operations.
		{name: "initializeOrderedBulkOp", statement: `db.users.initializeOrderedBulkOp().insert({a: 1}).execute()`, want: DML},

		// RS statements.
		{name: "rs.status()", statement: `rs.status()`, want: SelectInfoSchema},
		{name: "rs.conf()", statement: `rs.conf()`, want: SelectInfoSchema},
		{name: "rs.config()", statement: `rs.config()`, want: SelectInfoSchema},
		{name: "rs.printReplicationInfo()", statement: `rs.printReplicationInfo()`, want: SelectInfoSchema},
		{name: "rs.printSecondaryReplicationInfo()", statement: `rs.printSecondaryReplicationInfo()`, want: SelectInfoSchema},
		{name: "rs.initiate()", statement: `rs.initiate()`, want: DML},
		{name: "rs.help()", statement: `rs.help()`, want: DML},

		// SH statements.
		{name: "sh.status()", statement: `sh.status()`, want: SelectInfoSchema},
		{name: "sh.getBalancerState()", statement: `sh.getBalancerState()`, want: SelectInfoSchema},
		{name: "sh.isBalancerRunning()", statement: `sh.isBalancerRunning()`, want: SelectInfoSchema},
		{name: "sh.addShard()", statement: `sh.addShard("host:27017")`, want: DML},
		{name: "sh.help()", statement: `sh.help()`, want: DML},

		// runCommand and adminCommand.
		{name: "runCommand find", statement: `db.runCommand({find: "users"})`, want: Select},
		{name: "runCommand aggregate read-only", statement: `db.runCommand({aggregate: "users", pipeline: [{$match: {}}, {$project: {name: 1}}], cursor: {}})`, want: Select},
		{name: "runCommand aggregate $out", statement: `db.runCommand({aggregate: "users", pipeline: [{$out: "output"}], cursor: {}})`, want: DML},
		{name: "runCommand aggregate $project then $merge", statement: `db.runCommand({aggregate: "users", pipeline: [{$project: {name: 1}}, {$merge: {into: "output"}}], cursor: {}})`, want: DML},
		{name: "runCommand aggregate dynamic pipeline", statement: `db.runCommand({aggregate: "users", pipeline: myPipeline, cursor: {}})`, want: DML},
		{name: "runCommand aggregate missing pipeline", statement: `db.runCommand({aggregate: "users", cursor: {}})`, want: DML},
		{name: "runCommand insert", statement: `db.runCommand({insert: "users", documents: [{a: 1}]})`, want: DML},
		{name: "runCommand create", statement: `db.runCommand({create: "newColl"})`, want: DDL},
		{name: "runCommand serverStatus", statement: `db.runCommand({serverStatus: 1})`, want: SelectInfoSchema},
		{name: "adminCommand find", statement: `db.adminCommand({find: "users"})`, want: Select},
		{name: "adminCommand aggregate read-only", statement: `db.adminCommand({aggregate: 1, pipeline: [{$currentOp: {}}], cursor: {}})`, want: Select},
		{name: "adminCommand aggregate $out", statement: `db.adminCommand({aggregate: "users", pipeline: [{$out: "output"}], cursor: {}})`, want: DML},
		{name: "adminCommand aggregate $merge", statement: `db.adminCommand({aggregate: "users", pipeline: [{$merge: {into: "output"}}], cursor: {}})`, want: DML},
		{name: "adminCommand insert", statement: `db.adminCommand({insert: "users", documents: [{a: 1}]})`, want: DML},
		{name: "adminCommand create", statement: `db.adminCommand({create: "newColl"})`, want: DDL},
		{name: "adminCommand serverStatus", statement: `db.adminCommand({serverStatus: 1})`, want: SelectInfoSchema},

		// Unknown methods default to DML.
		{name: "unknown collection method", statement: `db.users.someUnknownMethod()`, want: DML},
		{name: "unknown db method", statement: `db.someUnknownDbMethod()`, want: DML},

		// Empty input defaults to DML with no error.
		{name: "empty input", statement: ``, want: DML},
		{name: "whitespace only", statement: "   \n\t  ", want: DML},
		{name: "semicolon only", statement: `;`, want: DML},
		{name: "comment only", statement: `// just a comment`, want: DML},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ClassifyStatement(tc.statement)
			if err != nil {
				t.Fatalf("ClassifyStatement(%q) unexpected error: %v", tc.statement, err)
			}
			if got != tc.want {
				t.Fatalf("ClassifyStatement(%q) = %v, want %v", tc.statement, got, tc.want)
			}
		})
	}
}

func TestClassifyStatementErrors(t *testing.T) {
	t.Run("incomplete find preserves parse error", func(t *testing.T) {
		got, err := ClassifyStatement(`db.users.find({`)
		if got != DML {
			t.Fatalf("ClassifyStatement = %v, want %v", got, DML)
		}
		var parseErr *parser.ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("expected *parser.ParseError, got %T (%v)", err, err)
		}
		if parseErr.Position < 0 {
			t.Fatalf("expected non-negative ParseError.Position, got %d", parseErr.Position)
		}
	})

	t.Run("unparseable input preserves parse error", func(t *testing.T) {
		got, err := ClassifyStatement(`@@@###`)
		if got != DML {
			t.Fatalf("ClassifyStatement = %v, want %v", got, DML)
		}
		var parseErr *parser.ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("expected *parser.ParseError, got %T (%v)", err, err)
		}
	})

	t.Run("multi-statement script returns error", func(t *testing.T) {
		got, err := ClassifyStatement(`db.users.find({}); db.users.drop()`)
		if got != DML {
			t.Fatalf("ClassifyStatement = %v, want %v", got, DML)
		}
		if err == nil {
			t.Fatalf("expected error for multi-statement input")
		}
	})
}

func TestQueryTypeOf(t *testing.T) {
	// Ensure QueryType is an alias for omni/analysis.QueryType.
	checkShared := func(qt analysis.QueryType) {
		if qt != analysis.Select {
			t.Fatalf("expected alias equality with analysis.Select, got %v", qt)
		}
	}
	checkShared(Select)

	tests := []struct {
		name string
		in   *StatementAnalysis
		want QueryType
	}{
		{name: "nil analysis", in: nil, want: DML},
		{name: "OpUnknown", in: &StatementAnalysis{Operation: OpUnknown}, want: DML},
		{name: "OpFind", in: &StatementAnalysis{Operation: OpFind, MethodName: "find"}, want: Select},
		{name: "OpFindOne", in: &StatementAnalysis{Operation: OpFindOne, MethodName: "findOne"}, want: Select},
		{name: "OpCount", in: &StatementAnalysis{Operation: OpCount, MethodName: "countDocuments"}, want: Select},
		{name: "OpDistinct", in: &StatementAnalysis{Operation: OpDistinct, MethodName: "distinct"}, want: Select},
		{name: "OpAggregate read", in: &StatementAnalysis{Operation: OpAggregate, PipelineStages: []string{"$match", "$project"}}, want: Select},
		{name: "OpAggregate $out", in: &StatementAnalysis{Operation: OpAggregate, PipelineStages: []string{"$match", "$out"}}, want: DML},
		{name: "OpAggregate $project then $out", in: &StatementAnalysis{Operation: OpAggregate, UnsupportedStage: "$project", PipelineStages: []string{"$project", "$out"}}, want: DML},
		{name: "OpAggregate $merge", in: &StatementAnalysis{Operation: OpAggregate, PipelineStages: []string{"$merge"}}, want: DML},
		{name: "OpAggregate unknown pipeline", in: &StatementAnalysis{Operation: OpAggregate, UnsupportedStage: "unknown"}, want: DML},
		{name: "OpAggregate unknown stage", in: &StatementAnalysis{Operation: OpAggregate, PipelineStages: []string{"$project", "unknown"}}, want: DML},
		{name: "OpRead collection metadata", in: &StatementAnalysis{Operation: OpRead, MethodName: "getIndexes"}, want: SelectInfoSchema},
		{name: "OpRead runCommand", in: &StatementAnalysis{Operation: OpRead, MethodName: "runCommand"}, want: Select},
		{name: "OpRead runCommand with $out stage", in: &StatementAnalysis{Operation: OpRead, MethodName: "runCommand", PipelineStages: []string{"$out"}}, want: DML},
		{name: "OpRead adminCommand", in: &StatementAnalysis{Operation: OpRead, MethodName: "adminCommand"}, want: Select},
		{name: "OpRead adminCommand with $merge stage", in: &StatementAnalysis{Operation: OpRead, MethodName: "adminCommand", PipelineStages: []string{"$merge"}}, want: DML},
		{name: "OpInfo show", in: &StatementAnalysis{Operation: OpInfo, MethodName: "show"}, want: SelectInfoSchema},
		{name: "OpInfo runCommand", in: &StatementAnalysis{Operation: OpInfo, MethodName: "runCommand"}, want: SelectInfoSchema},
		{name: "OpInfo adminCommand", in: &StatementAnalysis{Operation: OpInfo, MethodName: "adminCommand"}, want: SelectInfoSchema},
		{name: "OpInfo unknown method", in: &StatementAnalysis{Operation: OpInfo, MethodName: "customDbMethod"}, want: DML},
		{name: "OpWrite", in: &StatementAnalysis{Operation: OpWrite, MethodName: "insertOne"}, want: DML},
		{name: "OpAdmin", in: &StatementAnalysis{Operation: OpAdmin, MethodName: "createIndex"}, want: DDL},
		{name: "OpExplain", in: &StatementAnalysis{Operation: OpExplain, MethodName: "find"}, want: Explain},
		{name: "out of range operation", in: &StatementAnalysis{Operation: Operation(999)}, want: DML},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := QueryTypeOf(tc.in); got != tc.want {
				t.Fatalf("QueryTypeOf(%+v) = %v, want %v", tc.in, got, tc.want)
			}
			if got := tc.in.QueryType(); got != tc.want {
				t.Fatalf("(%+v).QueryType() = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
