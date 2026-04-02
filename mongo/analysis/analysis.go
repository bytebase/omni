// Package analysis extracts structural information from parsed MongoDB statements.
package analysis

import "github.com/bytebase/omni/mongo/ast"

// Operation classifies a MongoDB statement by its analytical behavior.
type Operation int

const (
	OpUnknown   Operation = iota
	OpFind                // find
	OpFindOne             // findOne
	OpAggregate           // aggregate
	OpCount               // countDocuments, estimatedDocumentCount, count
	OpDistinct            // distinct
	OpRead                // getIndexes, stats, storageSize, etc.
	OpWrite               // insertOne, updateMany, deleteOne, etc.
	OpAdmin               // createIndex, drop, renameCollection, etc.
	OpInfo                // show dbs, getCollectionNames, serverStatus, etc.
	OpExplain             // .explain() wrapped
)

// IsRead returns true for read operations (find, findOne, aggregate, count, distinct, read).
func (op Operation) IsRead() bool {
	switch op {
	case OpFind, OpFindOne, OpAggregate, OpCount, OpDistinct, OpRead:
		return true
	}
	return false
}

// IsWrite returns true for write operations.
func (op Operation) IsWrite() bool { return op == OpWrite }

// IsAdmin returns true for administrative operations.
func (op Operation) IsAdmin() bool { return op == OpAdmin }

// IsInfo returns true for informational/metadata operations.
func (op Operation) IsInfo() bool { return op == OpInfo }

// String returns a human-readable name for the operation.
func (op Operation) String() string {
	switch op {
	case OpFind:
		return "find"
	case OpFindOne:
		return "findOne"
	case OpAggregate:
		return "aggregate"
	case OpCount:
		return "count"
	case OpDistinct:
		return "distinct"
	case OpRead:
		return "read"
	case OpWrite:
		return "write"
	case OpAdmin:
		return "admin"
	case OpInfo:
		return "info"
	case OpExplain:
		return "explain"
	default:
		return "unknown"
	}
}

// StatementAnalysis is the result of analyzing a parsed MongoDB statement.
type StatementAnalysis struct {
	Operation        Operation
	MethodName       string   // original method name (e.g. "countDocuments" when Operation is OpCount)
	Collection       string   // target collection; empty for db-level/show commands
	PredicateFields  []string // sorted dot-paths from query filter document
	PipelineStages   []string // aggregate only: stage names in order
	ShapePreserving  bool     // true when all pipeline stages preserve document structure
	UnsupportedStage string   // first non-shape-preserving, non-join stage
	Joins            []JoinInfo
}

// JoinInfo records a join extracted from a $lookup or $graphLookup pipeline stage.
type JoinInfo struct {
	Collection string // the "from" field
	AsField    string // the "as" field
}

// Analyze extracts structural information from a single parsed MongoDB AST node.
// Returns nil for unrecognizable nodes.
func Analyze(node ast.Node) *StatementAnalysis {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.CollectionStatement:
		return analyzeCollection(n)
	case *ast.DatabaseStatement:
		return analyzeDatabase(n)
	case *ast.ShowCommand:
		return &StatementAnalysis{Operation: OpInfo, MethodName: "show"}
	case *ast.BulkStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: "bulkOp", Collection: n.Collection}
	case *ast.RsStatement:
		return analyzeRs(n)
	case *ast.ShStatement:
		return analyzeSh(n)
	case *ast.EncryptionStatement:
		return analyzeEncryption(n)
	case *ast.PlanCacheStatement:
		return analyzePlanCache(n)
	case *ast.SpStatement:
		return analyzeSp(n)
	case *ast.ConnectionStatement:
		// Connection constructors (e.g. connect, Mongo) are not read/info/admin;
		// default to OpWrite to match existing Bytebase behavior for unclassified statements.
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.Constructor}
	case *ast.NativeFunctionCall:
		// Native shell functions (e.g. load, sleep) are not read/info/admin;
		// default to OpWrite to match existing Bytebase behavior for unclassified statements.
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.Name}
	}
	return nil
}

func analyzeDatabase(s *ast.DatabaseStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.Method}
	switch s.Method {
	case "dropDatabase", "createCollection", "createView":
		a.Operation = OpAdmin
	case "getCollectionNames", "getCollectionInfos", "serverStatus",
		"serverBuildInfo", "version", "hostInfo", "getName",
		"listCommands", "stats":
		a.Operation = OpInfo
	case "runCommand", "adminCommand":
		a.Operation = classifyCommandArgs(s.Args)
	case "getSiblingDB", "getMongo":
		// These return a database handle for further chaining; they are not
		// read/info/admin so default to OpWrite per existing Bytebase behavior.
		a.Operation = OpWrite
	default:
		a.Operation = OpInfo
	}
	return a
}

func classifyCommandArgs(args []ast.Node) Operation {
	if len(args) == 0 {
		return OpWrite
	}
	doc, ok := args[0].(*ast.Document)
	if !ok || len(doc.Pairs) == 0 {
		return OpWrite
	}
	switch doc.Pairs[0].Key {
	case "find", "aggregate", "count", "distinct":
		return OpRead
	case "serverStatus", "listCollections", "listIndexes", "listDatabases",
		"collStats", "dbStats", "hostInfo", "buildInfo", "connectionStatus":
		return OpInfo
	case "create", "drop", "createIndexes", "dropIndexes", "renameCollection", "collMod":
		return OpAdmin
	default:
		return OpWrite
	}
}

func analyzeRs(s *ast.RsStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.MethodName}
	switch s.MethodName {
	case "status", "conf", "config", "printReplicationInfo", "printSecondaryReplicationInfo":
		a.Operation = OpInfo
	default:
		a.Operation = OpWrite
	}
	return a
}

func analyzeSh(s *ast.ShStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.MethodName}
	switch s.MethodName {
	case "status", "getBalancerState", "isBalancerRunning":
		a.Operation = OpInfo
	default:
		a.Operation = OpWrite
	}
	return a
}

func analyzeEncryption(s *ast.EncryptionStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.Target}
	if len(s.ChainedMethods) > 0 {
		last := s.ChainedMethods[len(s.ChainedMethods)-1]
		a.MethodName = last.Name
		switch last.Name {
		case "getKey", "getKeyByAltName", "getKeys", "decrypt", "encrypt", "encryptExpression":
			a.Operation = OpInfo
		default:
			a.Operation = OpWrite
		}
	} else {
		a.Operation = OpWrite
	}
	return a
}

func analyzePlanCache(s *ast.PlanCacheStatement) *StatementAnalysis {
	a := &StatementAnalysis{
		MethodName: "getPlanCache",
		Collection: s.Collection,
	}
	if len(s.ChainedMethods) > 0 {
		last := s.ChainedMethods[len(s.ChainedMethods)-1]
		a.MethodName = last.Name
		switch last.Name {
		case "list", "help":
			a.Operation = OpInfo
		default:
			a.Operation = OpWrite
		}
	} else {
		a.Operation = OpWrite
	}
	return a
}

func analyzeSp(s *ast.SpStatement) *StatementAnalysis {
	a := &StatementAnalysis{MethodName: s.MethodName}
	if s.SubMethod != "" {
		switch s.SubMethod {
		case "stats", "sample":
			a.Operation = OpInfo
		default:
			a.Operation = OpWrite
		}
	} else {
		switch s.MethodName {
		case "listConnections", "listStreamProcessors":
			a.Operation = OpInfo
		default:
			a.Operation = OpWrite
		}
	}
	return a
}

func analyzeCollection(s *ast.CollectionStatement) *StatementAnalysis {
	a := &StatementAnalysis{
		MethodName: s.Method,
		Collection: s.Collection,
	}

	if s.Explain {
		a.Operation = OpExplain
		return a
	}

	switch s.Method {
	case "find":
		a.Operation = OpFind
		a.PredicateFields = extractPredicateFields(s.Args)
	case "findOne":
		a.Operation = OpFindOne
		a.PredicateFields = extractPredicateFields(s.Args)
	case "aggregate":
		a.Operation = OpAggregate
		analyzePipelineInto(s.Args, a)
	case "countDocuments", "estimatedDocumentCount", "count":
		a.Operation = OpCount
	case "distinct":
		a.Operation = OpDistinct
	case "insertOne", "insertMany", "updateOne", "updateMany", "deleteOne", "deleteMany",
		"replaceOne", "findOneAndUpdate", "findOneAndReplace", "findOneAndDelete",
		"bulkWrite", "mapReduce":
		a.Operation = OpWrite
	case "createIndex", "createIndexes", "dropIndex", "dropIndexes", "drop",
		"renameCollection", "reIndex":
		a.Operation = OpAdmin
	case "getIndexes", "stats", "storageSize", "totalIndexSize", "totalSize",
		"dataSize", "validate", "latencyStats", "getShardDistribution":
		a.Operation = OpRead
	default:
		a.Operation = OpUnknown
	}

	return a
}
