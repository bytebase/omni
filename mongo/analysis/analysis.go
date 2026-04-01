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
		return &StatementAnalysis{Operation: OpInfo, MethodName: n.Method}
	case *ast.ShowCommand:
		return &StatementAnalysis{Operation: OpInfo, MethodName: "show"}
	case *ast.BulkStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: "bulkOp", Collection: n.Collection}
	case *ast.RsStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.MethodName}
	case *ast.ShStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.MethodName}
	case *ast.EncryptionStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.Target}
	case *ast.PlanCacheStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: "getPlanCache", Collection: n.Collection}
	case *ast.SpStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.MethodName}
	case *ast.ConnectionStatement:
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.Constructor}
	case *ast.NativeFunctionCall:
		return &StatementAnalysis{Operation: OpWrite, MethodName: n.Name}
	}
	return nil
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
